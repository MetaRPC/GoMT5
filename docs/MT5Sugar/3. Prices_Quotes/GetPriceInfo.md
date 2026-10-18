# 📊 Get Complete Price Information (`GetPriceInfo`)

> **Sugar method:** Returns BID, ASK, spread, and timestamp in one method call.

**API Information:**

* **Method:** `sugar.GetPriceInfo(symbol)`
* **Timeout:** 3 seconds (shared by both requests)
* **Returns:** `*PriceInfo` structure
* **Underlying calls:** `service.GetSymbolTick(ctx, symbol)` + `service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_POINT)` - two requests

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetPriceInfo(symbol string) (*PriceInfo, error)
```

```go
// PriceInfo structure
type PriceInfo struct {
    Symbol     string    // Symbol name
    Bid        float64   // Current BID price
    Ask        float64   // Current ASK price
    SpreadPips float64   // Spread in POINTS (despite the name)
    Time       time.Time // Time of the last tick in broker server time
    TimeUtc    time.Time // The same moment in UTC
}
```

**How the fields are filled:**

* `Bid`, `Ask` - from the same tick (`SymbolInfoTick`).
* `SpreadPips` - despite the name, the spread in **points**: `(Ask - Bid) / point`, where `point` is the symbol's `SYMBOL_POINT`
  (EURUSD with 5 digits: `0.00010 / 0.00001 = 10` points = 1 pip). It is computed from the tick, not read from
  `SYMBOL_SPREAD` like `GetSpread()`, so the two can differ slightly. Rounded to 0.1 point (float noise:
  3.9999999 → 4). A symbol whose `SYMBOL_POINT` is <= 0 gives an error.
* `Time` - the tick time in **broker server time**: MT5 writes the broker's wall clock as a Unix timestamp, and
  `GetSymbolTick` keeps it as a UTC-located `time.Time` (`time.Unix(tick.Time, 0).UTC()`), so
  `priceInfo.Time.Format(...)` prints the broker's clock. Do not use it for age checks - `time.Since(priceInfo.Time)`
  is off by the broker's UTC offset.
* `TimeUtc` - the same moment in real UTC (`time.Unix(tick.TimeUtc, 0).UTC()`). Use it for `time.Since(...)`.

---

## 💬 Just the Essentials

* **What it is:** Get BID, ASK, spread, and time from ONE method instead of three separate methods.
* **Why you need it:** 2 requests instead of 3 (GetBid + GetAsk + GetSpread), and BID/ASK come from the same tick.
* **Sanity check:** priceInfo.Ask ≥ priceInfo.Bid, SpreadPips = (ask - bid) / point.

---

## 🎯 When to Use

✅ **Price snapshots** - Need complete price picture at once

✅ **Performance** - Fewer API calls (2 instead of 3)

✅ **Logging/reporting** - Capture full price state with timestamp

✅ **Display dashboards** - Show all price data together

---

## 🔗 Usage Examples

### 1) Basic usage

```go
priceInfo, err := sugar.GetPriceInfo("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Symbol:  %s\n", priceInfo.Symbol)
fmt.Printf("BID:     %.5f\n", priceInfo.Bid)
fmt.Printf("ASK:     %.5f\n", priceInfo.Ask)
fmt.Printf("Spread:  %.0f points\n", priceInfo.SpreadPips)
fmt.Printf("Time:    %s\n", priceInfo.Time.Format("15:04:05"))

// Output:
// Symbol:  EURUSD
// BID:     1.08450
// ASK:     1.08460
// Spread:  10 points
// Time:    14:23:15
```

---

### 2) Price snapshot for logging

```go
func LogPriceSnapshot(sugar *mt5.MT5Sugar, symbol string) {
    priceInfo, err := sugar.GetPriceInfo(symbol)
    if err != nil {
        fmt.Printf("Error getting price: %v\n", err)
        return
    }

    // Log to file/database
    logEntry := fmt.Sprintf("[%s] %s: BID=%.5f ASK=%.5f SPREAD=%.0f",
        priceInfo.Time.Format("2006-01-02 15:04:05"),
        priceInfo.Symbol,
        priceInfo.Bid,
        priceInfo.Ask,
        priceInfo.SpreadPips)

    fmt.Println(logEntry)
    // In production: write to file or database
}

// Usage:
LogPriceSnapshot(sugar, "EURUSD")
```

---

### 3) Comparison: GetPriceInfo vs individual calls

```go
symbol := "EURUSD"

// ❌ LESS EFFICIENT - 3 separate API calls (and BID/ASK may come from different ticks)
bid, _ := sugar.GetBid(symbol)
ask, _ := sugar.GetAsk(symbol)
spread, _ := sugar.GetSpread(symbol)
fmt.Printf("BID: %.5f, ASK: %.5f, Spread: %.0f\n", bid, ask, spread)

// ✅ MORE EFFICIENT - 2 API calls (tick + point size)
priceInfo, _ := sugar.GetPriceInfo(symbol)
fmt.Printf("BID: %.5f, ASK: %.5f, Spread: %.0f\n",
    priceInfo.Bid, priceInfo.Ask, priceInfo.SpreadPips)
```

---

### 4) Real-time price dashboard

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()

for range ticker.C {
    fmt.Print("\033[H\033[2J") // Clear screen
    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║              REAL-TIME PRICE DASHBOARD                ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")
    fmt.Printf("Updated: %s\n\n", time.Now().Format("15:04:05"))

    fmt.Printf("%-10s  %-12s  %-12s  %-10s\n",
        "Symbol", "BID", "ASK", "Spread")
    fmt.Println("─────────────────────────────────────────────────────────")

    for _, symbol := range symbols {
        priceInfo, err := sugar.GetPriceInfo(symbol)
        if err != nil {
            continue
        }

        fmt.Printf("%-10s  %-12.5f  %-12.5f  %-10.0f\n",
            priceInfo.Symbol,
            priceInfo.Bid,
            priceInfo.Ask,
            priceInfo.SpreadPips)
    }
}
```

---

### 5) Trade signal with price snapshot

```go
func CheckTradingSignal(sugar *mt5.MT5Sugar, symbol string) {
    priceInfo, _ := sugar.GetPriceInfo(symbol)

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║       TRADE SIGNAL CHECK              ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Symbol:    %s\n", priceInfo.Symbol)
    fmt.Printf("Time:      %s\n", priceInfo.Time.Format("15:04:05"))
    fmt.Printf("BID:       %.5f\n", priceInfo.Bid)
    fmt.Printf("ASK:       %.5f\n", priceInfo.Ask)
    fmt.Printf("Spread:    %.0f points\n", priceInfo.SpreadPips)
    fmt.Println()

    // Check conditions
    if priceInfo.SpreadPips > 15 {
        fmt.Println("❌ Spread too high - NO TRADE")
        return
    }

    if priceInfo.Bid < 1.08000 {
        fmt.Println("✅ BUY SIGNAL - Price at support")
        // Place buy order
    } else if priceInfo.Bid > 1.09000 {
        fmt.Println("✅ SELL SIGNAL - Price at resistance")
        // Place sell order
    } else {
        fmt.Println("⏸️  NO SIGNAL - Price in range")
    }
}
```

---

### 6) Price history recorder

```go
type PriceSnapshot struct {
    Time   time.Time
    Bid    float64
    Ask    float64
    Spread float64
}

func RecordPriceHistory(sugar *mt5.MT5Sugar, symbol string, duration time.Duration) {
    history := []PriceSnapshot{}

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    timeout := time.After(duration)

    for {
        select {
        case <-timeout:
            // Save history
            fmt.Printf("\nRecorded %d price snapshots for %s\n",
                len(history), symbol)

            // Show summary
            if len(history) > 0 {
                first := history[0]
                last := history[len(history)-1]

                fmt.Printf("First: %.5f (BID) at %s\n",
                    first.Bid, first.Time.Format("15:04:05"))
                fmt.Printf("Last:  %.5f (BID) at %s\n",
                    last.Bid, last.Time.Format("15:04:05"))
                fmt.Printf("Change: %.5f (%.2f%%)\n",
                    last.Bid-first.Bid,
                    ((last.Bid-first.Bid)/first.Bid)*100)
            }
            return

        case <-ticker.C:
            priceInfo, _ := sugar.GetPriceInfo(symbol)

            snapshot := PriceSnapshot{
                Time:   priceInfo.Time,
                Bid:    priceInfo.Bid,
                Ask:    priceInfo.Ask,
                Spread: priceInfo.SpreadPips,
            }

            history = append(history, snapshot)
            fmt.Printf("%s: BID=%.5f SPREAD=%.0f\n",
                priceInfo.Time.Format("15:04:05"),
                priceInfo.Bid,
                priceInfo.SpreadPips)
        }
    }
}

// Record for 5 minutes
RecordPriceHistory(sugar, "EURUSD", 5*time.Minute)
```

---

### 7) Multi-symbol price comparison

```go
symbols := []string{"EURUSD", "GBPUSD", "AUDUSD", "NZDUSD"}

fmt.Println("USD Pairs Price Comparison:")
fmt.Println("═══════════════════════════════════════════════")
fmt.Printf("%-10s  %-12s  %-12s  %-10s  %s\n",
    "Pair", "BID", "ASK", "Spread", "Time")
fmt.Println("─────────────────────────────────────────────────")

for _, symbol := range symbols {
    priceInfo, err := sugar.GetPriceInfo(symbol)
    if err != nil {
        continue
    }

    fmt.Printf("%-10s  %-12.5f  %-12.5f  %-10.0f  %s\n",
        priceInfo.Symbol,
        priceInfo.Bid,
        priceInfo.Ask,
        priceInfo.SpreadPips,
        priceInfo.Time.Format("15:04:05"))
}
```

---

### 8) Price alert with full info

```go
func PriceAlertWithInfo(sugar *mt5.MT5Sugar, symbol string, targetBid float64) {
    fmt.Printf("🔔 Alert set: %s BID reaches %.5f\n", symbol, targetBid)

    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        priceInfo, _ := sugar.GetPriceInfo(symbol)

        if priceInfo.Bid <= targetBid {
            fmt.Println("\n🚨 PRICE ALERT TRIGGERED!")
            fmt.Println("═══════════════════════════════════════")
            fmt.Printf("Symbol:    %s\n", priceInfo.Symbol)
            fmt.Printf("Target:    %.5f\n", targetBid)
            fmt.Printf("Current:   %.5f (BID)\n", priceInfo.Bid)
            fmt.Printf("ASK:       %.5f\n", priceInfo.Ask)
            fmt.Printf("Spread:    %.0f points\n", priceInfo.SpreadPips)
            fmt.Printf("Time:      %s\n", priceInfo.Time.Format("15:04:05"))
            return
        }

        fmt.Printf("%s BID: %.5f (target: %.5f, %.5f away)\n",
            symbol, priceInfo.Bid, targetBid, priceInfo.Bid-targetBid)
    }
}
```

---

### 9) Check if price is fresh

```go
// PriceInfo.Time is broker server time, so time.Since(priceInfo.Time) is off by the
// broker's UTC offset. PriceInfo.TimeUtc is the same moment in UTC - use that for an age check.
func IsPriceFresh(sugar *mt5.MT5Sugar, symbol string, maxAge time.Duration) bool {
    priceInfo, err := sugar.GetPriceInfo(symbol)
    if err != nil {
        fmt.Printf("⚠️  Cannot read price: %v\n", err)
        return false
    }

    age := time.Since(priceInfo.TimeUtc)

    if age > maxAge {
        fmt.Printf("⚠️  WARNING: Price data is %v old (max: %v)\n", age, maxAge)
        return false
    }

    fmt.Printf("✅ Price data is fresh (%v old)\n", age)
    return true
}

// Usage:
if IsPriceFresh(sugar, "EURUSD", 5*time.Second) {
    // Use price for trading
    sugar.BuyMarket("EURUSD", 0.1)
} else {
    // Price too old, wait for update
    fmt.Println("Waiting for fresh price...")
}
```

---

### 10) Complete trading decision

```go
func MakeTradingDecision(sugar *mt5.MT5Sugar, symbol string) {
    priceInfo, err := sugar.GetPriceInfo(symbol)
    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }

    fmt.Println("╔═══════════════════════════════════════════════╗")
    fmt.Println("║         TRADING DECISION ANALYSIS             ║")
    fmt.Println("╚═══════════════════════════════════════════════╝")
    fmt.Printf("Symbol:      %s\n", priceInfo.Symbol)
    fmt.Printf("Server Time: %s\n", priceInfo.Time.Format("2006-01-02 15:04:05"))
    fmt.Printf("BID:         %.5f\n", priceInfo.Bid)
    fmt.Printf("ASK:         %.5f\n", priceInfo.Ask)
    fmt.Printf("Spread:      %.0f points\n", priceInfo.SpreadPips)
    fmt.Println()

    // Check 1: Spread
    if priceInfo.SpreadPips > 15 {
        fmt.Println("❌ REJECT: Spread too high")
        return
    }
    fmt.Printf("✅ Spread OK (%.0f < 15 points)\n", priceInfo.SpreadPips)

    // Check 2: Price freshness (uses TimeUtc - see example 9)
    age := time.Since(priceInfo.TimeUtc)
    if age > 10*time.Second {
        fmt.Printf("❌ REJECT: Price too old (%v)\n", age)
        return
    }
    fmt.Printf("✅ Price fresh (%.1f seconds old)\n", age.Seconds())

    // Check 3: Market hours (broker server hour: Time holds server wall-clock time)
    hour := priceInfo.Time.Hour()
    if hour < 8 || hour > 17 {
        fmt.Printf("❌ REJECT: Outside trading hours (%d:00)\n", hour)
        return
    }
    fmt.Printf("✅ Market hours OK (%d:00)\n", hour)

    // All checks passed
    fmt.Println("\n🟢 ALL CHECKS PASSED - READY TO TRADE")

    // Calculate position size
    lotSize, _ := sugar.CalculatePositionSize(symbol, 2.0, 50)
    fmt.Printf("Recommended lot size: %.2f (2%% risk, 50-point SL)\n", lotSize)
}
```

---

## 🔗 Related Methods

* `GetBid()` - Get only BID price
* `GetAsk()` - Get only ASK price
* `GetSpread()` - Get only spread
* `WaitForPrice()` - Wait for valid price with timeout

---

## 💎 Pro Tips

1. **Use this instead of individual calls** - Fewer requests (2 vs 3), one consistent tick
2. **Check timestamp** - Verify price is fresh before trading with `TimeUtc` (`Time` is broker server time)
3. **Store snapshots** - Keep price history for analysis
4. **Spread validation** - Always check SpreadPips before trading
5. **Perfect for dashboards** - Get all price data at once

---

**See also:** [`GetBid.md`](GetBid.md), [`GetAsk.md`](GetAsk.md), [`WaitForPrice.md`](WaitForPrice.md)
