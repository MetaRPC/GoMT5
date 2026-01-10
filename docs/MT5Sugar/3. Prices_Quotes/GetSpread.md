# 📏 Get Spread (`GetSpread`)

> **Sugar method:** Returns current spread in points (difference between ASK and BID).

**API Information:**

* **Method:** `sugar.GetSpread(symbol)`
* **Timeout:** 3 seconds
* **Returns:** Spread in points as `float64`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetSpread(symbol string) (float64, error)
```

---

## 💬 Just the Essentials

* **What it is:** Spread = ASK - BID in points (broker's fee for each trade).
* **Why you need it:** Check trading costs before opening positions. High spread = expensive trade.
* **Sanity check:** Lower spread = better. Typical EURUSD spread: 5-20 points.

---

## 🧮 Formula

```
Spread = ASK - BID (in points)

For 5-digit broker (EURUSD):
  ASK: 1.08460, BID: 1.08450
  Spread = (1.08460 - 1.08450) / 0.00001 = 10 points

Trading cost = Spread × Point × ContractSize × Volume
```

---

## 🔗 Usage Examples

### 1) Basic usage

```go
spread, err := sugar.GetSpread("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("EURUSD spread: %.0f points\n", spread)
// Output: EURUSD spread: 10 points
```

---

### 2) Check spread before trading

```go
symbol := "EURUSD"
maxAcceptableSpread := 15.0

spread, _ := sugar.GetSpread(symbol)

fmt.Printf("%s current spread: %.0f points\n", symbol, spread)

if spread > maxAcceptableSpread {
    fmt.Printf("❌ Spread too high (%.0f > %.0f)\n", spread, maxAcceptableSpread)
    fmt.Println("   Waiting for better conditions...")
} else {
    fmt.Printf("✅ Spread acceptable (%.0f ≤ %.0f)\n", spread, maxAcceptableSpread)
    sugar.BuyMarket(symbol, 0.1)
}
```

---

### 3) Compare spreads across symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD", "BTCUSD"}

fmt.Println("Current spreads:")
fmt.Println("─────────────────────────────")

for _, symbol := range symbols {
    spread, err := sugar.GetSpread(symbol)
    if err != nil {
        continue
    }

    // Categorize spread
    status := ""
    if spread < 10 {
        status = "🟢 Excellent"
    } else if spread < 20 {
        status = "🟡 Good"
    } else if spread < 50 {
        status = "🟠 High"
    } else {
        status = "🔴 Very High"
    }

    fmt.Printf("%-8s  %3.0f points  %s\n", symbol, spread, status)
}

// Output:
// Current spreads:
// ─────────────────────────────
// EURUSD    10 points  🟢 Excellent
// GBPUSD    15 points  🟡 Good
// USDJPY    12 points  🟡 Good
// XAUUSD    45 points  🟠 High
// BTCUSD   120 points  🔴 Very High
```

---

### 4) Calculate trading cost

```go
symbol := "EURUSD"
volume := 1.0 // 1 lot

spread, _ := sugar.GetSpread(symbol)
info, _ := sugar.GetSymbolInfo(symbol)

// Calculate cost in account currency
spreadCost := spread * info.Point * info.ContractSize * volume

fmt.Printf("Symbol: %s\n", symbol)
fmt.Printf("Volume: %.2f lots\n", volume)
fmt.Printf("Spread: %.0f points\n", spread)
fmt.Printf("Trading cost: $%.2f\n", spreadCost)

// Output:
// Symbol: EURUSD
// Volume: 1.00 lots
// Spread: 10 points
// Trading cost: $10.00
```

---

### 5) Monitor spread changes

```go
symbol := "EURUSD"

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

fmt.Printf("Monitoring %s spread:\n\n", symbol)

previousSpread := 0.0

for i := 0; i < 20; i++ {
    <-ticker.C

    spread, _ := sugar.GetSpread(symbol)

    if previousSpread > 0 {
        change := spread - previousSpread

        if change > 0 {
            fmt.Printf("%.0f points (↑ widened by %.0f)\n", spread, change)
        } else if change < 0 {
            fmt.Printf("%.0f points (↓ narrowed by %.0f)\n", spread, -change)
        } else {
            fmt.Printf("%.0f points (→ no change)\n", spread)
        }
    } else {
        fmt.Printf("%.0f points (initial)\n", spread)
    }

    previousSpread = spread
}
```

---

### 6) Find best time to trade (lowest spread)

```go
symbol := "EURUSD"
checkDuration := 5 * time.Minute
checkInterval := 10 * time.Second

minSpread := 99999.0
maxSpread := 0.0
totalSpread := 0.0
samples := 0

ticker := time.NewTicker(checkInterval)
defer ticker.Stop()

timeout := time.After(checkDuration)

for {
    select {
    case <-timeout:
        avgSpread := totalSpread / float64(samples)

        fmt.Println("\n═══════════════════════════════════")
        fmt.Printf("Spread analysis for %s:\n", symbol)
        fmt.Printf("Duration: %v\n", checkDuration)
        fmt.Printf("Minimum:  %.0f points\n", minSpread)
        fmt.Printf("Maximum:  %.0f points\n", maxSpread)
        fmt.Printf("Average:  %.0f points\n", avgSpread)
        fmt.Println("═══════════════════════════════════")
        return

    case <-ticker.C:
        spread, _ := sugar.GetSpread(symbol)

        if spread < minSpread {
            minSpread = spread
        }
        if spread > maxSpread {
            maxSpread = spread
        }

        totalSpread += spread
        samples++

        fmt.Printf("Spread: %.0f points (min: %.0f, max: %.0f)\n",
            spread, minSpread, maxSpread)
    }
}
```

---

### 7) Spread warning system

```go
func MonitorSpreadWarnings(sugar *mt5.MT5Sugar, symbol string, warningLevel float64) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        spread, _ := sugar.GetSpread(symbol)

        if spread > warningLevel {
            fmt.Printf("⚠️  WARNING: %s spread is %.0f points (threshold: %.0f)\n",
                symbol, spread, warningLevel)
        } else {
            fmt.Printf("✅ %s spread: %.0f points (OK)\n", symbol, spread)
        }
    }
}

// Run in background
go MonitorSpreadWarnings(sugar, "EURUSD", 20.0)
```

---

### 8) Validate spread before scalping

```go
// Scalping requires tight spreads
func CanScalp(sugar *mt5.MT5Sugar, symbol string) bool {
    spread, err := sugar.GetSpread(symbol)
    if err != nil {
        return false
    }

    maxScalpingSpread := 10.0 // Scalpers need very tight spreads

    if spread > maxScalpingSpread {
        fmt.Printf("❌ Spread %.0f too high for scalping (max: %.0f)\n",
            spread, maxScalpingSpread)
        return false
    }

    fmt.Printf("✅ Spread %.0f acceptable for scalping\n", spread)
    return true
}

if CanScalp(sugar, "EURUSD") {
    // Open quick in/out trades
    sugar.BuyMarket("EURUSD", 0.1)
}
```

---

### 9) Spread dashboard for multiple symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "AUDUSD", "USDCAD"}

ticker := time.NewTicker(10 * time.Second)
defer ticker.Stop()

for range ticker.C {
    fmt.Print("\033[H\033[2J") // Clear screen
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║        SPREAD MONITOR                 ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Time: %s\n\n", time.Now().Format("15:04:05"))

    fmt.Printf("%-10s  %-10s  %-10s\n", "Symbol", "Spread", "Status")
    fmt.Println("─────────────────────────────────────────")

    for _, symbol := range symbols {
        spread, _ := sugar.GetSpread(symbol)

        status := ""
        if spread < 10 {
            status = "🟢 Tight"
        } else if spread < 20 {
            status = "🟡 Normal"
        } else {
            status = "🔴 Wide"
        }

        fmt.Printf("%-10s  %-10.0f  %-10s\n", symbol, spread, status)
    }
}
```

---

### 10) Historical spread tracking

```go
type SpreadRecord struct {
    Time   time.Time
    Spread float64
}

func TrackSpreads(sugar *mt5.MT5Sugar, symbol string, duration time.Duration) {
    records := []SpreadRecord{}

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    timeout := time.After(duration)

    for {
        select {
        case <-timeout:
            // Analyze collected data
            if len(records) == 0 {
                return
            }

            var total float64
            min := records[0].Spread
            max := records[0].Spread

            for _, r := range records {
                total += r.Spread
                if r.Spread < min {
                    min = r.Spread
                }
                if r.Spread > max {
                    max = r.Spread
                }
            }

            avg := total / float64(len(records))

            fmt.Println("\n═══════ SPREAD ANALYSIS ═══════")
            fmt.Printf("Symbol:   %s\n", symbol)
            fmt.Printf("Samples:  %d\n", len(records))
            fmt.Printf("Min:      %.0f points\n", min)
            fmt.Printf("Max:      %.0f points\n", max)
            fmt.Printf("Average:  %.0f points\n", avg)
            return

        case t := <-ticker.C:
            spread, _ := sugar.GetSpread(symbol)
            records = append(records, SpreadRecord{
                Time:   t,
                Spread: spread,
            })
            fmt.Printf("%s: %.0f points\n", t.Format("15:04:05"), spread)
        }
    }
}

// Track for 1 hour
TrackSpreads(sugar, "EURUSD", 1*time.Hour)
```

---

## 🔗 Related Methods

* `GetBid()` - BID price
* `GetAsk()` - ASK price
* `GetPriceInfo()` - Get BID, ASK, and spread at once ⭐

---

## ⚠️ Common Pitfalls

### 1) Ignoring spread when calculating profit

```go
// ❌ WRONG - forgetting spread cost
entry := 1.08460
target := 1.08560
profit := target - entry // 100 pips
// Real profit is less due to spread!

// ✅ CORRECT - account for spread
spread, _ := sugar.GetSpread("EURUSD")
info, _ := sugar.GetSymbolInfo("EURUSD")
spreadInPrice := spread * info.Point
realProfit := (target - entry) - spreadInPrice
```

---

## 💎 Pro Tips

1. **Lower is better** - Always prefer low spread symbols/times
2. **Volatile times = wide spreads** - News events widen spreads
3. **Major pairs = tight spreads** - EUR/USD, GBP/USD have lowest
4. **Check before scalping** - Scalpers need spreads < 10 points
5. **Use GetPriceInfo()** - More efficient than calculating manually

---

**See also:** [`GetBid.md`](GetBid.md), [`GetAsk.md`](GetAsk.md), [`GetPriceInfo.md`](GetPriceInfo.md)
