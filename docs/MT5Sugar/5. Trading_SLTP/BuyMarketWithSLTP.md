# 🟢🛡️ Buy Market with SL/TP (`BuyMarketWithSLTP`)

> **Sugar method:** Opens BUY position at market with Stop Loss and Take Profit in one call.

**API Information:**

* **Method:** `sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)`
* **Timeout:** 10 seconds
* **Returns:** Position ticket number

---

## 📋 Method Signature

```go
func (s *MT5Sugar) BuyMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `sl` | `float64` | Stop Loss price (0 = no SL) |
| `tp` | `float64` | Take Profit price (0 = no TP) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Position ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** Opens BUY position with SL/TP set immediately - safer than `BuyMarket()`.
* **Why you need it:** Complete risk management in one call, no separate modification needed.
* **Sanity check:** For BUY: SL must be **below** entry, TP must be **above** entry.

---

## 🎯 When to Use

✅ **Production trading** - Always use SL/TP in real trading

✅ **Risk-managed entries** - Set protection from the start

✅ **Quick entries** - One call instead of two (open + modify)

✅ **Automated trading** - Ensure every position has risk limits

❌ **NOT recommended:** Setting SL=0 and TP=0 (use `BuyMarket()` instead)

---

## 🔗 Usage Examples

### 1) Basic usage - fixed SL/TP prices

```go
symbol := "EURUSD"
volume := 0.1
currentAsk := 1.08500
sl := 1.08000  // 50 pips below entry
tp := 1.09500  // 100 pips above entry

ticket, err := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY position opened with protection\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  ~%.5f\n", currentAsk)
fmt.Printf("   SL:     %.5f (-50 pips)\n", sl)
fmt.Printf("   TP:     %.5f (+100 pips)\n", tp)
```

---

### 2) Calculate SL/TP from current price

```go
symbol := "EURUSD"
volume := 0.1

// Get current price
ask, _ := sugar.GetAsk(symbol)

// Calculate SL/TP
info, _ := sugar.GetSymbolInfo(symbol)
sl := ask - (50 * info.Point)  // 50 pips SL
tp := ask + (100 * info.Point) // 100 pips TP

ticket, _ := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ BUY with calculated SL/TP\n")
fmt.Printf("   Entry: %.5f\n", ask)
fmt.Printf("   SL:    %.5f (50 pips)\n", sl)
fmt.Printf("   TP:    %.5f (100 pips)\n", tp)
fmt.Printf("   R:R:   1:2\n")
```

---

### 3) Using CalculateSLTP helper

```go
symbol := "EURUSD"
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Use helper to calculate SL/TP prices
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, slPips, tpPips)

ticket, _ := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ BUY using CalculateSLTP helper\n")
fmt.Printf("   SL: %.5f (%.0f pips)\n", sl, slPips)
fmt.Printf("   TP: %.5f (%.0f pips)\n", tp, tpPips)
```

---

### 4) Risk-managed entry with position sizing

```go
symbol := "EURUSD"
riskPercent := 2.0
slPips := 50.0
tpPips := 100.0

// Calculate lot size based on risk
lotSize, _ := sugar.CalculatePositionSize(symbol, riskPercent, slPips)

// Calculate SL/TP prices
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, slPips, tpPips)

// Open with full risk management
ticket, err := sugar.BuyMarketWithSLTP(symbol, lotSize, sl, tp)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Risk-managed BUY position\n")
fmt.Printf("   Risk:   %.1f%% of balance\n", riskPercent)
fmt.Printf("   Size:   %.2f lots\n", lotSize)
fmt.Printf("   SL:     %.0f pips\n", slPips)
fmt.Printf("   TP:     %.0f pips\n", tpPips)
fmt.Printf("   Ticket: #%d\n", ticket)
```

---

### 5) 2:1 Risk/Reward ratio

```go
symbol := "EURUSD"
volume := 0.1
slPips := 30.0
tpPips := 60.0 // 2:1 R:R

sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, slPips, tpPips)
ticket, _ := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ BUY with 2:1 Risk/Reward\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Risk:   %.0f pips\n", slPips)
fmt.Printf("   Reward: %.0f pips\n", tpPips)
fmt.Printf("   R:R:    1:%.0f\n", tpPips/slPips)
```

---

### 6) ATR-based SL/TP

```go
symbol := "EURUSD"
volume := 0.1

// Assume we calculated ATR = 80 pips
atrValue := 80.0

// Use 1.5x ATR for SL, 3x ATR for TP
slPips := 1.5 * atrValue  // 120 pips
tpPips := 3.0 * atrValue  // 240 pips

sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, slPips, tpPips)
ticket, _ := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ ATR-based BUY position\n")
fmt.Printf("   ATR:    %.0f pips\n", atrValue)
fmt.Printf("   SL:     %.0f pips (1.5x ATR)\n", slPips)
fmt.Printf("   TP:     %.0f pips (3.0x ATR)\n", tpPips)
fmt.Printf("   Ticket: #%d\n", ticket)
```

---

### 7) Support/Resistance based SL/TP

```go
symbol := "EURUSD"
volume := 0.1

// Technical levels
supportLevel := 1.08000
resistanceLevel := 1.09000

// Get current entry
ask, _ := sugar.GetAsk(symbol)

// Set SL below support, TP at resistance
info, _ := sugar.GetSymbolInfo(symbol)
sl := supportLevel - (10 * info.Point) // 10 pips buffer below support
tp := resistanceLevel

ticket, _ := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ BUY with technical levels\n")
fmt.Printf("   Entry:      %.5f\n", ask)
fmt.Printf("   Support:    %.5f\n", supportLevel)
fmt.Printf("   SL:         %.5f (below support)\n", sl)
fmt.Printf("   Resistance: %.5f\n", resistanceLevel)
fmt.Printf("   TP:         %.5f (at resistance)\n", tp)
```

---

### 8) Only Stop Loss (no TP)

```go
symbol := "EURUSD"
volume := 0.1
slPips := 50.0

sl, _, _ := sugar.CalculateSLTP(symbol, "BUY", 0, slPips, 0)
tp := 0.0 // No take profit

ticket, _ := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ BUY with SL only (manual TP management)\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   SL:     %.5f (%.0f pips)\n", sl, slPips)
fmt.Printf("   TP:     None - will manage manually\n")
```

---

### 9) Validation before opening

```go
symbol := "EURUSD"
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Pre-check 1: Can we open?
canOpen, reason, _ := sugar.CanOpenPosition(symbol, volume)
if !canOpen {
    fmt.Printf("❌ Cannot open: %s\n", reason)
    return
}

// Pre-check 2: Spread acceptable?
spread, _ := sugar.GetSpread(symbol)
if spread > 15 {
    fmt.Printf("❌ Spread too high: %.0f points\n", spread)
    return
}

// All checks passed - open position
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, slPips, tpPips)
ticket, _ := sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ Position opened after validation\n")
fmt.Printf("   Ticket: #%d\n", ticket)
```

---

### 10) Complete trading function

```go
func OpenBuyPosition(sugar *mt5.MT5Sugar, symbol string, riskPercent, slPips, tpPips float64) (uint64, error) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      OPENING BUY POSITION             ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Step 1: Validate symbol
    available, _ := sugar.IsSymbolAvailable(symbol)
    if !available {
        return 0, fmt.Errorf("symbol %s not available", symbol)
    }
    fmt.Printf("✅ Symbol: %s\n", symbol)

    // Step 2: Check spread
    spread, _ := sugar.GetSpread(symbol)
    if spread > 20 {
        return 0, fmt.Errorf("spread too high: %.0f", spread)
    }
    fmt.Printf("✅ Spread: %.0f points\n", spread)

    // Step 3: Calculate position size
    lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, slPips)
    if err != nil {
        return 0, fmt.Errorf("position size calc failed: %w", err)
    }
    fmt.Printf("✅ Lot size: %.2f (%.1f%% risk)\n", lotSize, riskPercent)

    // Step 4: Check margin
    canOpen, reason, _ := sugar.CanOpenPosition(symbol, lotSize)
    if !canOpen {
        return 0, fmt.Errorf("insufficient margin: %s", reason)
    }
    fmt.Printf("✅ Margin: OK\n")

    // Step 5: Calculate SL/TP
    sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, slPips, tpPips)
    fmt.Printf("✅ SL: %.5f (%.0f pips)\n", sl, slPips)
    fmt.Printf("✅ TP: %.5f (%.0f pips)\n", tp, tpPips)

    // Step 6: Open position
    ticket, err := sugar.BuyMarketWithSLTP(symbol, lotSize, sl, tp)
    if err != nil {
        return 0, fmt.Errorf("order failed: %w", err)
    }

    fmt.Printf("\n🎯 Position opened successfully!\n")
    fmt.Printf("   Ticket: #%d\n", ticket)
    return ticket, nil
}

// Usage:
ticket, err := OpenBuyPosition(sugar, "EURUSD", 2.0, 50, 100)
```

---

## 🔗 Related Methods

**🍬 Alternative (using pips instead of prices):**

* `BuyMarketWithPips()` - More intuitive, specify SL/TP in pips ⭐ **RECOMMENDED**

**🍬 Other market orders:**

* `BuyMarket()` - No SL/TP (NOT recommended for production)
* `SellMarketWithSLTP()` - SELL with SL/TP

**🍬 Helper methods:**

* `CalculateSLTP()` - Calculate SL/TP prices from pips
* `CalculatePositionSize()` - Risk-based position sizing

---

## ⚠️ Common Pitfalls

### 1) Wrong SL/TP direction

```go
// ❌ WRONG - SL must be BELOW entry for BUY
ask := 1.08500
sugar.BuyMarketWithSLTP("EURUSD", 0.1, 1.09000, 1.10000) // SL above entry!

// ✅ CORRECT - SL below, TP above
sugar.BuyMarketWithSLTP("EURUSD", 0.1, 1.08000, 1.09000)
```

### 2) Not validating SL/TP distances

```go
// ❌ WRONG - SL too close (might be rejected)
ask := 1.08500
sl := 1.08499 // Only 1 pip!
sugar.BuyMarketWithSLTP("EURUSD", 0.1, sl, 1.09000)

// ✅ CORRECT - check minimum stop level
stopLevel, _ := sugar.GetMinStopLevel("EURUSD")
// Ensure SL is at least stopLevel pips away
```

### 3) Using prices instead of pips

```go
// ❌ HARDER - calculate prices manually
ask, _ := sugar.GetAsk("EURUSD")
info, _ := sugar.GetSymbolInfo("EURUSD")
sl := ask - (50 * info.Point)
tp := ask + (100 * info.Point)
sugar.BuyMarketWithSLTP("EURUSD", 0.1, sl, tp)

// ✅ EASIER - use BuyMarketWithPips
sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
```

---

## 💎 Pro Tips

1. **Use `BuyMarketWithPips()`** - More intuitive than calculating prices

2. **Always set SL** - Never trade without stop loss in production

3. **2:1 minimum R:R** - TP should be at least 2x SL distance

4. **Validate before opening** - Check spread, margin, symbol availability

5. **Calculate position size** - Use `CalculatePositionSize()` for risk management

---

**See also:** [`BuyMarketWithPips.md`](../11. Trading_Helpers/BuyMarketWithPips.md), [`CalculateSLTP.md`](../11. Trading_Helpers/CalculateSLTP.md), [`SellMarketWithSLTP.md`](SellMarketWithSLTP.md)
