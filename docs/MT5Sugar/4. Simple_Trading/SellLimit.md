# 🔴📈 Sell Limit Order (`SellLimit`)

> **Sugar method:** Places pending SELL order that executes when price **rises** to specified level.

**API Information:**

* **Method:** `sugar.SellLimit(symbol, volume, price)`
* **Timeout:** 10 seconds
* **Returns:** Pending order ticket number

---

## 📋 Method Signature

```go
func (s *MT5Sugar) SellLimit(symbol string, volume, price float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `price` | `float64` | Entry price (must be **ABOVE** current BID) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Pending order ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** Pending SELL order that activates when price **rises** to your target.
* **Why you need it:** Sell at a better price (resistance level) instead of current market.
* **Sanity check:** Entry price must be **ABOVE** current BID. Order waits until price reaches your level.

---

## 📊 SELL LIMIT Logic

```
Current BID: 1.08500
SELL LIMIT:  1.09000 (500 pips ABOVE current price)

When price rises to 1.09000 → Order executes automatically
```

**Use case:** Sell at resistance / Sell the rally / Enter on price rise

---

## 🎯 When to Use

✅ **Resistance levels** - Sell when price reaches resistance

✅ **Rally exhaustion** - Enter on price retracement up

✅ **Sell the rally** - Wait for price to rise before selling

✅ **Better entry** - Get better price than current market

❌ **NOT for breakdowns** - Use `SellStop()` for breakdowns

---

## 🔗 Usage Examples

### 1) Basic usage - sell at resistance

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
resistanceLevel := currentBid + 0.00050 // 50 pips above

ticket, err := sugar.SellLimit(symbol, 0.1, resistanceLevel)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ SELL LIMIT order placed\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  %.5f (current: %.5f)\n", resistanceLevel, currentBid)
fmt.Printf("   Will execute when price rises %.0f pips\n",
    (resistanceLevel-currentBid)/0.00001)
```

---

### 2) Sell on rally exhaustion

```go
symbol := "EURUSD"

// Current downtrend, wait for rally to fade
currentBid, _ := sugar.GetBid(symbol)
rallyLevel := currentBid + 0.00030 // Wait for 30 pip rally

ticket, _ := sugar.SellLimit(symbol, 0.1, rallyLevel)

fmt.Printf("SELL LIMIT set for rally exhaustion entry\n")
fmt.Printf("Current: %.5f\n", currentBid)
fmt.Printf("Entry:   %.5f (rally level)\n", rallyLevel)
```

---

### 3) Multiple sell limits at different levels

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
volume := 0.1

// Set limits at different resistance levels
levels := []float64{
    currentBid + 0.00020, // +20 pips
    currentBid + 0.00050, // +50 pips
    currentBid + 0.00100, // +100 pips
}

fmt.Println("Placing multiple SELL LIMIT orders:")

for i, level := range levels {
    ticket, err := sugar.SellLimit(symbol, volume, level)
    if err != nil {
        fmt.Printf("Level %d failed: %v\n", i+1, err)
        continue
    }

    pipsAway := (level - currentBid) / 0.00001
    fmt.Printf("✅ Level %d: %.5f (%.0f pips away) - Ticket #%d\n",
        i+1, level, pipsAway, ticket)
}

// Output:
// ✅ Level 1: 1.08700 (20 pips away) - Ticket #12345
// ✅ Level 2: 1.09000 (50 pips away) - Ticket #12346
// ✅ Level 3: 1.09500 (100 pips away) - Ticket #12347
```

---

### 4) Sell limit with validation

```go
symbol := "EURUSD"
volume := 0.1

// Get current price
currentBid, _ := sugar.GetBid(symbol)

// Calculate entry price (resistance level)
info, _ := sugar.GetSymbolInfo(symbol)
entryPrice := currentBid + (50 * info.Point) // 50 pips above

// Validate entry price is above current
if entryPrice <= currentBid {
    fmt.Printf("❌ Error: Entry %.5f must be above BID %.5f\n",
        entryPrice, currentBid)
    return
}

// Check minimum stop level
stopLevel, _ := sugar.GetMinStopLevel(symbol)
minDistance := float64(stopLevel) * info.Point

if (entryPrice - currentBid) < minDistance {
    fmt.Printf("❌ Entry too close (min: %.5f)\n", minDistance)
    return
}

// Place order
ticket, _ := sugar.SellLimit(symbol, volume, entryPrice)
fmt.Printf("✅ SELL LIMIT placed at %.5f\n", entryPrice)
```

---

### 5) Sell limit with SL/TP (using separate method)

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
entryPrice := currentBid + 0.00050 // 50 pips above

// Place sell limit
ticket, _ := sugar.SellLimit(symbol, 0.1, entryPrice)

// Calculate SL/TP from entry price
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", entryPrice, 30, 60)

// Set SL/TP on pending order
sugar.ModifyPositionSLTP(ticket, sl, tp)

fmt.Printf("SELL LIMIT with SL/TP:\n")
fmt.Printf("  Entry: %.5f\n", entryPrice)
fmt.Printf("  SL:    %.5f (+30 pips)\n", sl)
fmt.Printf("  TP:    %.5f (-60 pips)\n", tp)
```

---

### 6) Scale in with sell limits

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
totalVolume := 0.3
entries := 3

volumePerEntry := totalVolume / float64(entries)

fmt.Printf("Scaling in with %d SELL LIMIT orders\n", entries)

for i := 1; i <= entries; i++ {
    // Each level 20 pips apart
    entryPrice := currentBid + (float64(i) * 0.00020)

    ticket, err := sugar.SellLimit(symbol, volumePerEntry, entryPrice)
    if err != nil {
        continue
    }

    fmt.Printf("Entry %d: %.2f lots at %.5f - Ticket #%d\n",
        i, volumePerEntry, entryPrice, ticket)
}
```

---

### 7) Cancel unfilled sell limit

```go
symbol := "EURUSD"
entryPrice := 1.09000

ticket, _ := sugar.SellLimit(symbol, 0.1, entryPrice)
fmt.Printf("SELL LIMIT placed: Ticket #%d\n", ticket)

// Wait 5 minutes
fmt.Println("Waiting 5 minutes...")
time.Sleep(5 * time.Minute)

// Check if filled
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    // Not filled - cancel it
    fmt.Println("Order not filled - canceling...")
    // Note: Use OrderDelete for pending orders
    // sugar.DeleteOrder(ticket) // If you have this method
} else {
    fmt.Printf("✅ Order filled at %.5f\n", pos.PriceOpen)
}
```

---

### 8) Sell limit at Fibonacci retracement

```go
symbol := "EURUSD"

// Previous swing high and low
swingHigh := 1.09000
swingLow := 1.08000
range_ := swingHigh - swingLow

// 61.8% Fibonacci retracement (common resistance in downtrend)
fib618 := swingLow + (range_ * 0.618)

// Place sell limit at Fib level
ticket, _ := sugar.SellLimit(symbol, 0.1, fib618)

fmt.Printf("SELL LIMIT at Fibonacci 61.8%%\n")
fmt.Printf("Swing High: %.5f\n", swingHigh)
fmt.Printf("Swing Low:  %.5f\n", swingLow)
fmt.Printf("Entry:      %.5f (Fib 61.8%%)\n", fib618)
fmt.Printf("Ticket:     #%d\n", ticket)
```

---

### 9) Grid trading with sell limits

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)

gridLevels := 5
gridSpacing := 0.00020 // 20 pips
volume := 0.01

fmt.Printf("Grid trading: %d SELL LIMIT orders\n", gridLevels)

for i := 1; i <= gridLevels; i++ {
    entryPrice := currentBid + (float64(i) * gridSpacing)

    ticket, err := sugar.SellLimit(symbol, volume, entryPrice)
    if err != nil {
        fmt.Printf("Grid %d failed: %v\n", i, err)
        continue
    }

    pipsAway := (entryPrice - currentBid) / 0.00001
    fmt.Printf("Grid %d: %.5f (%.0f pips) - #%d\n",
        i, entryPrice, pipsAway, ticket)
}
```

---

### 10) Advanced sell limit with monitoring

```go
func PlaceSellLimitWithMonitoring(sugar *mt5.MT5Sugar, symbol string, pipsAbove float64) {
    // Get current price
    currentBid, _ := sugar.GetBid(symbol)
    info, _ := sugar.GetSymbolInfo(symbol)

    // Calculate entry
    entryPrice := currentBid + (pipsAbove * info.Point)

    fmt.Printf("╔═══════════════════════════════════════╗\n")
    fmt.Printf("║       SELL LIMIT ORDER                ║\n")
    fmt.Printf("╚═══════════════════════════════════════╝\n")
    fmt.Printf("Symbol:       %s\n", symbol)
    fmt.Printf("Current BID:  %.5f\n", currentBid)
    fmt.Printf("Entry Price:  %.5f\n", entryPrice)
    fmt.Printf("Distance:     %.0f pips\n", pipsAbove)

    // Place order
    ticket, err := sugar.SellLimit(symbol, 0.1, entryPrice)
    if err != nil {
        fmt.Printf("❌ Order failed: %v\n", err)
        return
    }

    fmt.Printf("✅ Order placed: #%d\n", ticket)
    fmt.Println("\nMonitoring price...")

    // Monitor until filled or timeout
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    timeout := time.After(30 * time.Minute)

    for {
        select {
        case <-timeout:
            fmt.Println("⏰ Timeout - order not filled in 30 minutes")
            return

        case <-ticker.C:
            bid, _ := sugar.GetBid(symbol)
            distance := entryPrice - bid
            distancePips := distance / info.Point

            // Check if filled
            _, err := sugar.GetPositionByTicket(ticket)
            if err == nil {
                fmt.Printf("\n✅ ORDER FILLED at %.5f!\n", bid)
                return
            }

            fmt.Printf("BID: %.5f (%.0f pips away from entry)\n",
                bid, distancePips)
        }
    }
}

// Usage:
PlaceSellLimitWithMonitoring(sugar, "EURUSD", 50)
```

---

## 🔗 Related Methods

**🍬 Other limit orders:**

* `BuyLimit()` - BUY when price goes DOWN
* `SellLimitWithSLTP()` - SELL LIMIT with SL/TP

**🍬 Stop orders (for breakdowns):**

* `SellStop()` - SELL when price breaks DOWN

**🍬 Market orders:**

* `SellMarket()` - SELL immediately at current price

---

## ⚠️ Common Pitfalls

### 1) Setting entry price below current price

```go
// ❌ WRONG - Sell Limit must be ABOVE current price
currentBid := 1.08500
sugar.SellLimit("EURUSD", 0.1, 1.08000) // BELOW! Will be rejected!

// ✅ CORRECT - Above current price
sugar.SellLimit("EURUSD", 0.1, 1.09000) // Above BID
```

### 2) Confusing SellLimit with SellStop

```go
// SELL LIMIT = wait for price to RISE then sell
sugar.SellLimit("EURUSD", 0.1, 1.09000) // Sell when rises to 1.09000

// SELL STOP = wait for price to FALL then sell (breakdown)
sugar.SellStop("EURUSD", 0.1, 1.08000) // Sell when falls to 1.08000
```

---

## 💎 Pro Tips

1. **SELL LIMIT = sell more expensive** - Wait for price to rise

2. **Set at resistance levels** - Use technical analysis

3. **Multiple levels** - Scale in at different prices

4. **Monitor fills** - Check if order was executed

5. **Use with SL/TP** - Set risk management on pending orders

---

**See also:** [`SellLimitWithSLTP.md`](../5. Trading_SLTP/SellLimitWithSLTP.md), [`SellStop.md`](SellStop.md), [`BuyLimit.md`](BuyLimit.md)
