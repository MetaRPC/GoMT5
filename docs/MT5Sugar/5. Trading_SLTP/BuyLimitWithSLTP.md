# 🟢📉🛡️ Buy Limit with SL/TP (`BuyLimitWithSLTP`)

> **Sugar method:** Places pending BUY LIMIT order with Stop Loss and Take Profit set from the start.

**API Information:**

* **Method:** `sugar.BuyLimitWithSLTP(symbol, volume, price, sl, tp)`
* **Timeout:** 10 seconds
* **Returns:** Pending order ticket number (`result.Order`)
* **Request:** `service.PlaceOrder(ctx, &pb.OrderSendRequest{Symbol, Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT, Volume, Price: &price, StopLoss: &sl, TakeProfit: &tp})` -
  no expiration, comment or magic number; price, volume, SL and TP are sent exactly as given (no rounding, no check
  against the current price or of the stop distances - the trade server rejects invalid values)
* **Success:** only return code `10009` (TRADE_RETCODE_DONE); any other code gives the error
  `order rejected, code: N, comment: ...`, a failed call gives `BuyLimitWithSLTP failed: ...`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) BuyLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `price` | `float64` | Entry price (must be **BELOW** current ASK) |
| `sl` | `float64` | Stop Loss price (0 = no SL) |
| `tp` | `float64` | Take Profit price (0 = no TP) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Pending order ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** BUY LIMIT pending order with SL/TP - complete setup in one call.
* **Why you need it:** Set everything (entry, SL, TP) and forget - activates automatically.
* **Sanity check:** Entry below ASK, SL below entry, TP above entry.

---

## 📊 Logic

```
Current ASK: 1.08500
Entry:       1.08000 (LIMIT - 50 pips below)
SL:          1.07500 (50 pips below entry)
TP:          1.09000 (100 pips above entry)

When price drops to 1.08000 → Order fills with SL/TP already set
```

---

## 🎯 When to Use

✅ **Support level entries** - Buy at support with protection

✅ **Pullback trading** - Wait for dip, enter with risk management

✅ **Set and forget** - Don't need to monitor for fill

✅ **Automated strategies** - Complete order setup in advance

---

## 🔗 Usage Examples

### 1) Basic usage - support with SL/TP

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)

entryPrice := 1.08000 // Support level
sl := 1.07500         // 50 pips below entry
tp := 1.09000         // 100 pips above entry

ticket, err := sugar.BuyLimitWithSLTP(symbol, 0.1, entryPrice, sl, tp)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY LIMIT with SL/TP placed\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  %.5f (current: %.5f)\n", entryPrice, currentAsk)
fmt.Printf("   SL:     %.5f (-50 pips from entry)\n", sl)
fmt.Printf("   TP:     %.5f (+100 pips from entry)\n", tp)
```

---

### 2) Calculate SL/TP from entry price

```go
symbol := "EURUSD"
entryPrice := 1.08000
volume := 0.1

info, _ := sugar.GetSymbolInfo(symbol)

// Calculate SL/TP relative to entry
sl := entryPrice - (50 * info.Point)  // 50 points SL
tp := entryPrice + (100 * info.Point) // 100 points TP

ticket, _ := sugar.BuyLimitWithSLTP(symbol, volume, entryPrice, sl, tp)

fmt.Printf("✅ BUY LIMIT with calculated SL/TP\n")
fmt.Printf("   Entry: %.5f\n", entryPrice)
fmt.Printf("   SL:    %.5f (50 points)\n", sl)
fmt.Printf("   TP:    %.5f (100 points)\n", tp)
fmt.Printf("   R:R:   1:2\n")
```

---

### 3) Using CalculateSLTP helper

```go
symbol := "EURUSD"
entryPrice := 1.08000
volume := 0.1
slPips := 50.0  // distances in POINTS (the Sugar "Pips" parameters are points)
tpPips := 100.0

// Calculate SL/TP from entry price
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", entryPrice, slPips, tpPips)

ticket, _ := sugar.BuyLimitWithSLTP(symbol, volume, entryPrice, sl, tp)

fmt.Printf("✅ BUY LIMIT using helper\n")
fmt.Printf("   Entry: %.5f\n", entryPrice)
fmt.Printf("   SL:    %.5f (%.0f points)\n", sl, slPips)
fmt.Printf("   TP:    %.5f (%.0f points)\n", tp, tpPips)
```

---

### 4) Multiple buy limits at support zones

```go
symbol := "EURUSD"
volume := 0.05

// Three support levels
supports := []struct {
    entry float64
    sl    float64
    tp    float64
}{
    {1.08000, 1.07500, 1.09000}, // Strong support
    {1.07500, 1.07000, 1.08500}, // Medium support
    {1.07000, 1.06500, 1.08000}, // Weak support
}

fmt.Println("Placing BUY LIMITS at support zones:")

for i, level := range supports {
    ticket, err := sugar.BuyLimitWithSLTP(symbol, volume, level.entry, level.sl, level.tp)
    if err != nil {
        fmt.Printf("Level %d failed: %v\n", i+1, err)
        continue
    }

    fmt.Printf("✅ Level %d: Entry %.5f, SL %.5f, TP %.5f - Ticket #%d\n",
        i+1, level.entry, level.sl, level.tp, ticket)
}
```

---

### 5) Fibonacci retracement with tight SL

```go
symbol := "EURUSD"

// Previous swing
swingHigh := 1.09000
swingLow := 1.08000
range_ := swingHigh - swingLow

// 61.8% Fibonacci retracement
fib618 := swingHigh - (range_ * 0.618)

info, _ := sugar.GetSymbolInfo(symbol)

// Tight SL below Fib level, TP at swing high
sl := fib618 - (30 * info.Point)  // 30 points SL
tp := swingHigh                    // Target swing high

ticket, _ := sugar.BuyLimitWithSLTP(symbol, 0.1, fib618, sl, tp)

fmt.Printf("✅ Fibonacci BUY LIMIT\n")
fmt.Printf("   Entry: %.5f (Fib 61.8%%)\n", fib618)
fmt.Printf("   SL:    %.5f (30 points)\n", sl)
fmt.Printf("   TP:    %.5f (swing high)\n", tp)
```

---

### 6) Risk-based position sizing with limit order

```go
symbol := "EURUSD"
entryPrice := 1.08000
riskPercent := 2.0
slPips := 50.0
tpPips := 100.0

// Calculate lot size based on risk
lotSize, _ := sugar.CalculatePositionSize(symbol, riskPercent, slPips)

// Calculate SL/TP from entry
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", entryPrice, slPips, tpPips)

ticket, _ := sugar.BuyLimitWithSLTP(symbol, lotSize, entryPrice, sl, tp)

fmt.Printf("✅ Risk-managed BUY LIMIT\n")
fmt.Printf("   Entry:  %.5f\n", entryPrice)
fmt.Printf("   Risk:   %.1f%% of balance\n", riskPercent)
fmt.Printf("   Size:   %.2f lots\n", lotSize)
fmt.Printf("   SL:     %.0f points\n", slPips)
fmt.Printf("   TP:     %.0f points\n", tpPips)
fmt.Printf("   Ticket: #%d\n", ticket)
```

---

### 7) Scale in strategy with SL/TP

```go
symbol := "EURUSD"
totalVolume := 0.3
entries := 3

volumePerEntry := totalVolume / float64(entries)
info, _ := sugar.GetSymbolInfo(symbol)

fmt.Printf("Scaling in with %d BUY LIMIT orders:\n", entries)

baseEntry := 1.08000
baseSL := 1.07500
baseTP := 1.09000

for i := 0; i < entries; i++ {
    // Each entry 20 points lower
    entry := baseEntry - (float64(i) * 20 * info.Point)
    sl := baseSL - (float64(i) * 20 * info.Point)
    tp := baseTP // Same TP for all

    ticket, err := sugar.BuyLimitWithSLTP(symbol, volumePerEntry, entry, sl, tp)
    if err != nil {
        continue
    }

    fmt.Printf("Entry %d: %.2f lots at %.5f (SL %.5f, TP %.5f) - #%d\n",
        i+1, volumePerEntry, entry, sl, tp, ticket)
}
```

---

### 8) Pending order with validation

```go
symbol := "EURUSD"
entryPrice := 1.08000
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Validate entry is below current ASK
currentAsk, _ := sugar.GetAsk(symbol)
if entryPrice >= currentAsk {
    fmt.Printf("❌ Entry %.5f must be below ASK %.5f\n", entryPrice, currentAsk)
    return
}

// Check minimum stop level
stopLevel, _ := sugar.GetMinStopLevel(symbol)
info, _ := sugar.GetSymbolInfo(symbol)
minDistance := float64(stopLevel) * info.Point

if (currentAsk - entryPrice) < minDistance {
    fmt.Printf("❌ Entry too close to current price\n")
    return
}

// Calculate SL/TP
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", entryPrice, slPips, tpPips)

// Validate SL/TP distances
if (entryPrice - sl) < minDistance {
    fmt.Printf("❌ SL too close to entry\n")
    return
}

ticket, _ := sugar.BuyLimitWithSLTP(symbol, volume, entryPrice, sl, tp)
fmt.Printf("✅ Validated BUY LIMIT placed - Ticket #%d\n", ticket)
```

---

### 9) Weekend gap strategy

```go
// Friday close: 1.08500
// Expecting gap down Monday morning
symbol := "EURUSD"

// Set buy limit below Friday close
entryPrice := 1.08000 // 50 pips gap
volume := 0.1

info, _ := sugar.GetSymbolInfo(symbol)

// Tight SL (gap fill expected)
sl := entryPrice - (30 * info.Point)
tp := 1.08500 // Target Friday close

ticket, _ := sugar.BuyLimitWithSLTP(symbol, volume, entryPrice, sl, tp)

fmt.Printf("Weekend gap strategy:\n")
fmt.Printf("   Friday close: 1.08500\n")
fmt.Printf("   Entry:        %.5f (gap down)\n", entryPrice)
fmt.Printf("   SL:           %.5f (30 points)\n", sl)
fmt.Printf("   TP:           %.5f (gap fill)\n", tp)
```

---

### 10) Complete limit order function

```go
// slPips/tpPips are distances in POINTS (as CalculatePositionSize/CalculateSLTP expect)
func PlaceBuyLimitWithProtection(sugar *mt5.MT5Sugar, symbol string, entryPrice, riskPercent, slPips, tpPips float64) (uint64, error) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║    BUY LIMIT WITH SL/TP               ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Validate entry price
    currentAsk, _ := sugar.GetAsk(symbol)
    if entryPrice >= currentAsk {
        return 0, fmt.Errorf("entry %.5f must be below ASK %.5f", entryPrice, currentAsk)
    }
    fmt.Printf("✅ Entry: %.5f (%.0f points below current)\n",
        entryPrice, (currentAsk-entryPrice)/0.00001)

    // Calculate position size
    lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, slPips)
    if err != nil {
        return 0, err
    }
    fmt.Printf("✅ Size: %.2f lots (%.1f%% risk)\n", lotSize, riskPercent)

    // Calculate SL/TP
    sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", entryPrice, slPips, tpPips)
    fmt.Printf("✅ SL: %.5f (%.0f points)\n", sl, slPips)
    fmt.Printf("✅ TP: %.5f (%.0f points)\n", tp, tpPips)

    // Place order
    ticket, err := sugar.BuyLimitWithSLTP(symbol, lotSize, entryPrice, sl, tp)
    if err != nil {
        return 0, err
    }

    fmt.Printf("\n🎯 BUY LIMIT placed successfully!\n")
    fmt.Printf("   Ticket: #%d\n", ticket)
    return ticket, nil
}

// Usage:
ticket, _ := PlaceBuyLimitWithProtection(sugar, "EURUSD", 1.08000, 2.0, 50, 100)
```

---

## 🔗 Related Methods

**🍬 Other pending orders with SL/TP:**

* `SellLimitWithSLTP()` - SELL LIMIT with protection
* BUY STOP with protection - no Sugar method; use `GetService().PlaceOrder()` with `Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP`, `Price`, `StopLoss`, `TakeProfit`

**🍬 Market orders with SL/TP:**

* `BuyMarketWithSLTP()` - Immediate BUY with protection

**🍬 Helper methods:**

* `CalculateSLTP()` - Calculate SL/TP from distances in points
* `BuyLimit()` - Without SL/TP

---

## ⚠️ Common Pitfalls

### 1) Wrong SL/TP direction

```go
// ❌ WRONG - SL must be BELOW entry for BUY
entry := 1.08000
sugar.BuyLimitWithSLTP("EURUSD", 0.1, entry, 1.08500, 1.09000) // SL above!

// ✅ CORRECT - SL below, TP above
sugar.BuyLimitWithSLTP("EURUSD", 0.1, entry, 1.07500, 1.09000)
```

### 2) Calculating SL/TP from current price instead of entry

```go
// ❌ WRONG - SL/TP calculated from current price
point := 0.00001 // EURUSD point (info.Point from GetSymbolInfo)
ask, _ := sugar.GetAsk("EURUSD")
entry := 1.08000
sl := ask - (50 * point) // Wrong! Should be from entry!

// ✅ CORRECT - calculate from entry price
sl = entry - (50 * point)
tp := entry + (100 * point)
```

---

## 💎 Pro Tips

1. **Calculate from entry price** - Not current price

2. **Validate entry level** - Must be below current ASK

3. **Check min stop level** - Both entry and SL must respect broker limits

4. **Use for support trading** - Perfect for buying dips

5. **Monitor for fills** - Check if order was executed

---

**See also:** [`BuyLimit.md`](../4.%20Simple_Trading/BuyLimit.md), [`SellLimitWithSLTP.md`](SellLimitWithSLTP.md), [`CalculateSLTP.md`](../11.%20Trading_Helpers/CalculateSLTP.md)
