# 🔴📈🛡️ Sell Limit with SL/TP (`SellLimitWithSLTP`)

> **Sugar method:** Places pending SELL LIMIT order with Stop Loss and Take Profit set from the start.

**API Information:**

* **Method:** `sugar.SellLimitWithSLTP(symbol, volume, price, sl, tp)`
* **Timeout:** 10 seconds
* **Returns:** Pending order ticket number

---

## 📋 Method Signature

```go
func (s *MT5Sugar) SellLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `price` | `float64` | Entry price (must be **ABOVE** current BID) |
| `sl` | `float64` | Stop Loss price (0 = no SL) |
| `tp` | `float64` | Take Profit price (0 = no TP) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Pending order ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** SELL LIMIT pending order with SL/TP - complete setup in one call.
* **Why you need it:** Set everything (entry, SL, TP) and forget - activates automatically.
* **Sanity check:** Entry above BID, SL above entry, TP below entry.

---

## 📊 Logic

```
Current BID: 1.08500
Entry:       1.09000 (LIMIT - 50 pips above)
SL:          1.09500 (50 pips above entry)
TP:          1.08000 (100 pips below entry)

When price rises to 1.09000 → Order fills with SL/TP already set
```

---

## 🎯 When to Use

✅ **Resistance level entries** - Sell at resistance with protection

✅ **Rally exhaustion** - Wait for price rise, enter with risk management

✅ **Set and forget** - Don't need to monitor for fill

✅ **Automated strategies** - Complete order setup in advance

---

## 🔗 Usage Examples

### 1) Basic usage - resistance with SL/TP

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)

entryPrice := 1.09000 // Resistance level
sl := 1.09500         // 50 pips above entry
tp := 1.08000         // 100 pips below entry

ticket, err := sugar.SellLimitWithSLTP(symbol, 0.1, entryPrice, sl, tp)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ SELL LIMIT with SL/TP placed\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  %.5f (current: %.5f)\n", entryPrice, currentBid)
fmt.Printf("   SL:     %.5f (+50 pips from entry)\n", sl)
fmt.Printf("   TP:     %.5f (-100 pips from entry)\n", tp)
```

---

### 2) Calculate SL/TP from entry price

```go
symbol := "EURUSD"
entryPrice := 1.09000
volume := 0.1

info, _ := sugar.GetSymbolInfo(symbol)

// Calculate SL/TP relative to entry
sl := entryPrice + (50 * info.Point)  // 50 pips SL (above entry)
tp := entryPrice - (100 * info.Point) // 100 pips TP (below entry)

ticket, _ := sugar.SellLimitWithSLTP(symbol, volume, entryPrice, sl, tp)

fmt.Printf("✅ SELL LIMIT with calculated SL/TP\n")
fmt.Printf("   Entry: %.5f\n", entryPrice)
fmt.Printf("   SL:    %.5f (50 pips above)\n", sl)
fmt.Printf("   TP:    %.5f (100 pips below)\n", tp)
fmt.Printf("   R:R:   1:2\n")
```

---

### 3) Using CalculateSLTP helper

```go
symbol := "EURUSD"
entryPrice := 1.09000
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Calculate SL/TP from entry price
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", entryPrice, slPips, tpPips)

ticket, _ := sugar.SellLimitWithSLTP(symbol, volume, entryPrice, sl, tp)

fmt.Printf("✅ SELL LIMIT using helper\n")
fmt.Printf("   Entry: %.5f\n", entryPrice)
fmt.Printf("   SL:    %.5f (%.0f pips)\n", sl, slPips)
fmt.Printf("   TP:    %.5f (%.0f pips)\n", tp, tpPips)
```

---

### 4) Multiple sell limits at resistance zones

```go
symbol := "EURUSD"
volume := 0.05

// Three resistance levels
resistances := []struct {
    entry float64
    sl    float64
    tp    float64
}{
    {1.09000, 1.09500, 1.08000}, // Strong resistance
    {1.09500, 1.10000, 1.08500}, // Medium resistance
    {1.10000, 1.10500, 1.09000}, // Weak resistance
}

fmt.Println("Placing SELL LIMITS at resistance zones:")

for i, level := range resistances {
    ticket, err := sugar.SellLimitWithSLTP(symbol, volume, level.entry, level.sl, level.tp)
    if err != nil {
        fmt.Printf("Level %d failed: %v\n", i+1, err)
        continue
    }

    fmt.Printf("✅ Level %d: Entry %.5f, SL %.5f, TP %.5f - Ticket #%d\n",
        i+1, level.entry, level.sl, level.tp, ticket)
}
```

---

### 5) Fibonacci retracement sell

```go
symbol := "EURUSD"

// Previous swing in downtrend
swingHigh := 1.09000
swingLow := 1.08000
range_ := swingHigh - swingLow

// 50% Fibonacci retracement (sell opportunity)
fib50 := swingLow + (range_ * 0.50)

info, _ := sugar.GetSymbolInfo(symbol)

// SL above Fib level, TP at swing low
sl := fib50 + (30 * info.Point)  // 30 pips SL
tp := swingLow                    // Target swing low

ticket, _ := sugar.SellLimitWithSLTP(symbol, 0.1, fib50, sl, tp)

fmt.Printf("✅ Fibonacci SELL LIMIT\n")
fmt.Printf("   Entry: %.5f (Fib 50%%)\n", fib50)
fmt.Printf("   SL:    %.5f (30 pips)\n", sl)
fmt.Printf("   TP:    %.5f (swing low)\n", tp)
```

---

### 6) Risk-based position sizing with sell limit

```go
symbol := "EURUSD"
entryPrice := 1.09000
riskPercent := 2.0
slPips := 50.0
tpPips := 100.0

// Calculate lot size based on risk
lotSize, _ := sugar.CalculatePositionSize(symbol, riskPercent, slPips)

// Calculate SL/TP from entry
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", entryPrice, slPips, tpPips)

ticket, _ := sugar.SellLimitWithSLTP(symbol, lotSize, entryPrice, sl, tp)

fmt.Printf("✅ Risk-managed SELL LIMIT\n")
fmt.Printf("   Entry:  %.5f\n", entryPrice)
fmt.Printf("   Risk:   %.1f%% of balance\n", riskPercent)
fmt.Printf("   Size:   %.2f lots\n", lotSize)
fmt.Printf("   SL:     %.0f pips\n", slPips)
fmt.Printf("   TP:     %.0f pips\n", tpPips)
fmt.Printf("   Ticket: #%d\n", ticket)
```

---

### 7) Round number resistance

```go
symbol := "EURUSD"

// Round number acts as psychological resistance
roundNumber := 1.09000

info, _ := sugar.GetSymbolInfo(symbol)

// Entry at round number
entry := roundNumber
sl := entry + (40 * info.Point)  // 40 pips SL
tp := entry - (80 * info.Point)  // 80 pips TP

ticket, _ := sugar.SellLimitWithSLTP(symbol, 0.1, entry, sl, tp)

fmt.Printf("✅ Round number SELL LIMIT\n")
fmt.Printf("   Entry: %.5f (psychological level)\n", entry)
fmt.Printf("   SL:    %.5f (40 pips)\n", sl)
fmt.Printf("   TP:    %.5f (80 pips, 1:2 R:R)\n", tp)
```

---

### 8) Trendline resistance sell

```go
symbol := "EURUSD"

// Trendline resistance at current time
trendlinePrice := 1.08800
volume := 0.1

info, _ := sugar.GetSymbolInfo(symbol)

// Tight SL above trendline
sl := trendlinePrice + (25 * info.Point)  // 25 pips
tp := trendlinePrice - (75 * info.Point)  // 75 pips (1:3 R:R)

ticket, _ := sugar.SellLimitWithSLTP(symbol, volume, trendlinePrice, sl, tp)

fmt.Printf("✅ Trendline resistance SELL\n")
fmt.Printf("   Entry: %.5f (trendline)\n", trendlinePrice)
fmt.Printf("   SL:    %.5f (25 pips)\n", sl)
fmt.Printf("   TP:    %.5f (75 pips)\n", tp)
fmt.Printf("   R:R:   1:3\n")
```

---

### 9) Pending order with full validation

```go
symbol := "EURUSD"
entryPrice := 1.09000
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Validate entry is above current BID
currentBid, _ := sugar.GetBid(symbol)
if entryPrice <= currentBid {
    fmt.Printf("❌ Entry %.5f must be above BID %.5f\n", entryPrice, currentBid)
    return
}

// Check minimum stop level
stopLevel, _ := sugar.GetMinStopLevel(symbol)
info, _ := sugar.GetSymbolInfo(symbol)
minDistance := float64(stopLevel) * info.Point

if (entryPrice - currentBid) < minDistance {
    fmt.Printf("❌ Entry too close to current price\n")
    return
}

// Calculate SL/TP
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", entryPrice, slPips, tpPips)

// Validate SL/TP distances
if (sl - entryPrice) < minDistance {
    fmt.Printf("❌ SL too close to entry\n")
    return
}

ticket, _ := sugar.SellLimitWithSLTP(symbol, volume, entryPrice, sl, tp)
fmt.Printf("✅ Validated SELL LIMIT placed - Ticket #%d\n", ticket)
```

---

### 10) Complete sell limit function

```go
func PlaceSellLimitWithProtection(
    sugar *mt5.MT5Sugar,
    symbol string,
    entryPrice float64,
    riskPercent float64,
    slPips float64,
    tpPips float64,
) (uint64, error) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║    SELL LIMIT WITH SL/TP              ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Validate entry price
    currentBid, _ := sugar.GetBid(symbol)
    if entryPrice <= currentBid {
        return 0, fmt.Errorf("entry %.5f must be above BID %.5f", entryPrice, currentBid)
    }
    fmt.Printf("✅ Entry: %.5f (%.0f pips above current)\n",
        entryPrice, (entryPrice-currentBid)/0.00001)

    // Calculate position size
    lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, slPips)
    if err != nil {
        return 0, err
    }
    fmt.Printf("✅ Size: %.2f lots (%.1f%% risk)\n", lotSize, riskPercent)

    // Calculate SL/TP
    sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", entryPrice, slPips, tpPips)
    fmt.Printf("✅ SL: %.5f (%.0f pips)\n", sl, slPips)
    fmt.Printf("✅ TP: %.5f (%.0f pips)\n", tp, tpPips)

    // Place order
    ticket, err := sugar.SellLimitWithSLTP(symbol, lotSize, entryPrice, sl, tp)
    if err != nil {
        return 0, err
    }

    fmt.Printf("\n🎯 SELL LIMIT placed successfully!\n")
    fmt.Printf("   Ticket: #%d\n", ticket)
    return ticket, nil
}

// Usage:
ticket, _ := PlaceSellLimitWithProtection(sugar, "EURUSD", 1.09000, 2.0, 50, 100)
```

---

## 🔗 Related Methods

**🍬 Other pending orders with SL/TP:**

* `BuyLimitWithSLTP()` - BUY LIMIT with protection
* `SellStopWithSLTP()` - SELL STOP with protection

**🍬 Market orders with SL/TP:**

* `SellMarketWithSLTP()` - Immediate SELL with protection

**🍬 Helper methods:**

* `CalculateSLTP()` - Calculate SL/TP from pips
* `SellLimit()` - Without SL/TP (add manually later)

---

## ⚠️ Common Pitfalls

### 1) Wrong SL/TP direction

```go
// ❌ WRONG - SL must be ABOVE entry for SELL
entry := 1.09000
sugar.SellLimitWithSLTP("EURUSD", 0.1, entry, 1.08500, 1.08000) // SL below!

// ✅ CORRECT - SL above, TP below
sugar.SellLimitWithSLTP("EURUSD", 0.1, entry, 1.09500, 1.08000)
```

### 2) Calculating SL/TP from current price instead of entry

```go
// ❌ WRONG - SL/TP calculated from current price
bid, _ := sugar.GetBid("EURUSD")
entry := 1.09000
sl := bid + (50 * point) // Wrong! Should be from entry!

// ✅ CORRECT - calculate from entry price
sl := entry + (50 * point)
tp := entry - (100 * point)
```

### 3) Entry price below current BID

```go
// ❌ WRONG - SELL LIMIT must be ABOVE current BID
bid := 1.08500
sugar.SellLimitWithSLTP("EURUSD", 0.1, 1.08000, sl, tp) // Below!

// ✅ CORRECT - entry above current BID
sugar.SellLimitWithSLTP("EURUSD", 0.1, 1.09000, sl, tp)
```

---

## 💎 Pro Tips

1. **Calculate from entry price** - Not current price

2. **Validate entry level** - Must be above current BID

3. **Check min stop level** - Both entry and SL must respect broker limits

4. **Use for resistance trading** - Perfect for selling rallies

5. **Monitor for fills** - Check if order was executed

---

**See also:** [`SellLimit.md`](../4. Simple_Trading/SellLimit.md), [`BuyLimitWithSLTP.md`](BuyLimitWithSLTP.md), [`CalculateSLTP.md`](../11. Trading_Helpers/CalculateSLTP.md)
