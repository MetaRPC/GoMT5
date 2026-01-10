# 🔴🛡️ Sell Market with SL/TP (`SellMarketWithSLTP`)

> **Sugar method:** Opens SELL position at market with Stop Loss and Take Profit in one call.

**API Information:**

* **Method:** `sugar.SellMarketWithSLTP(symbol, volume, sl, tp)`
* **Timeout:** 10 seconds
* **Returns:** Position ticket number

---

## 📋 Method Signature

```go
func (s *MT5Sugar) SellMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error)
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

* **What it is:** Opens SELL position with SL/TP set immediately - safer than `SellMarket()`.
* **Why you need it:** Complete risk management in one call, no separate modification needed.
* **Sanity check:** For SELL: SL must be **above** entry, TP must be **below** entry.

---

## 🎯 When to Use

✅ **Production trading** - Always use SL/TP in real trading

✅ **Risk-managed entries** - Set protection from the start

✅ **Quick entries** - One call instead of two (open + modify)

✅ **Automated trading** - Ensure every position has risk limits

❌ **NOT recommended:** Setting SL=0 and TP=0 (use `SellMarket()` instead)

---

## 🔗 Usage Examples

### 1) Basic usage - fixed SL/TP prices

```go
symbol := "EURUSD"
volume := 0.1
currentBid := 1.08500
sl := 1.09000  // 50 pips above entry
tp := 1.07500  // 100 pips below entry

ticket, err := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ SELL position opened with protection\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  ~%.5f\n", currentBid)
fmt.Printf("   SL:     %.5f (+50 pips)\n", sl)
fmt.Printf("   TP:     %.5f (-100 pips)\n", tp)
```

---

### 2) Calculate SL/TP from current price

```go
symbol := "EURUSD"
volume := 0.1

// Get current price
bid, _ := sugar.GetBid(symbol)

// Calculate SL/TP
info, _ := sugar.GetSymbolInfo(symbol)
sl := bid + (50 * info.Point)  // 50 pips SL (above entry)
tp := bid - (100 * info.Point) // 100 pips TP (below entry)

ticket, _ := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ SELL with calculated SL/TP\n")
fmt.Printf("   Entry: %.5f\n", bid)
fmt.Printf("   SL:    %.5f (50 pips above)\n", sl)
fmt.Printf("   TP:    %.5f (100 pips below)\n", tp)
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
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, slPips, tpPips)

ticket, _ := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ SELL using CalculateSLTP helper\n")
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
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, slPips, tpPips)

// Open with full risk management
ticket, err := sugar.SellMarketWithSLTP(symbol, lotSize, sl, tp)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Risk-managed SELL position\n")
fmt.Printf("   Risk:   %.1f%% of balance\n", riskPercent)
fmt.Printf("   Size:   %.2f lots\n", lotSize)
fmt.Printf("   SL:     %.0f pips\n", slPips)
fmt.Printf("   TP:     %.0f pips\n", tpPips)
fmt.Printf("   Ticket: #%d\n", ticket)
```

---

### 5) 3:1 Risk/Reward ratio

```go
symbol := "EURUSD"
volume := 0.1
slPips := 30.0
tpPips := 90.0 // 3:1 R:R

sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, slPips, tpPips)
ticket, _ := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ SELL with 3:1 Risk/Reward\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Risk:   %.0f pips\n", slPips)
fmt.Printf("   Reward: %.0f pips\n", tpPips)
fmt.Printf("   R:R:    1:%.0f\n", tpPips/slPips)
```

---

### 6) Resistance/Support based SL/TP

```go
symbol := "EURUSD"
volume := 0.1

// Technical levels
resistanceLevel := 1.09000
supportLevel := 1.08000

// Get current entry
bid, _ := sugar.GetBid(symbol)

// Set SL above resistance, TP at support
info, _ := sugar.GetSymbolInfo(symbol)
sl := resistanceLevel + (10 * info.Point) // 10 pips buffer above resistance
tp := supportLevel

ticket, _ := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ SELL with technical levels\n")
fmt.Printf("   Entry:      %.5f\n", bid)
fmt.Printf("   Resistance: %.5f\n", resistanceLevel)
fmt.Printf("   SL:         %.5f (above resistance)\n", sl)
fmt.Printf("   Support:    %.5f\n", supportLevel)
fmt.Printf("   TP:         %.5f (at support)\n", tp)
```

---

### 7) Trend following sell

```go
symbol := "EURUSD"
volume := 0.1

// Strong downtrend - previous low at 1.08000
targetLevel := 1.08000

bid, _ := sugar.GetBid(symbol)
info, _ := sugar.GetSymbolInfo(symbol)

// Tight SL above recent swing high
sl := bid + (40 * info.Point)  // 40 pips SL
tp := targetLevel              // Target previous low

ticket, _ := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ Trend following SELL\n")
fmt.Printf("   Entry:  %.5f\n", bid)
fmt.Printf("   SL:     %.5f (above swing)\n", sl)
fmt.Printf("   Target: %.5f (previous low)\n", tp)
```

---

### 8) Breakdown sell with wide SL

```go
symbol := "EURUSD"
volume := 0.05 // Smaller size for wider SL

bid, _ := sugar.GetBid(symbol)
info, _ := sugar.GetSymbolInfo(symbol)

// Wide SL for breakdown volatility
slPips := 100.0
tpPips := 200.0

sl := bid + (slPips * info.Point)
tp := bid - (tpPips * info.Point)

ticket, _ := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ Breakdown SELL with wide SL\n")
fmt.Printf("   Volume: %.2f (reduced for wider SL)\n", volume)
fmt.Printf("   SL:     %.0f pips\n", slPips)
fmt.Printf("   TP:     %.0f pips\n", tpPips)
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

// Pre-check 3: Check price direction
bid, _ := sugar.GetBid(symbol)
fmt.Printf("Current BID: %.5f\n", bid)

// All checks passed - open position
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, slPips, tpPips)
ticket, _ := sugar.SellMarketWithSLTP(symbol, volume, sl, tp)

fmt.Printf("✅ Position opened after validation\n")
fmt.Printf("   Ticket: #%d\n", ticket)
```

---

### 10) Complete trading function

```go
func OpenSellPosition(sugar *mt5.MT5Sugar, symbol string, riskPercent, slPips, tpPips float64) (uint64, error) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      OPENING SELL POSITION            ║")
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
    sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, slPips, tpPips)
    fmt.Printf("✅ SL: %.5f (%.0f pips)\n", sl, slPips)
    fmt.Printf("✅ TP: %.5f (%.0f pips)\n", tp, tpPips)

    // Step 6: Open position
    ticket, err := sugar.SellMarketWithSLTP(symbol, lotSize, sl, tp)
    if err != nil {
        return 0, fmt.Errorf("order failed: %w", err)
    }

    fmt.Printf("\n🎯 Position opened successfully!\n")
    fmt.Printf("   Ticket: #%d\n", ticket)
    return ticket, nil
}

// Usage:
ticket, err := OpenSellPosition(sugar, "EURUSD", 2.0, 50, 100)
```

---

## 🔗 Related Methods

**🍬 Alternative (using pips instead of prices):**

* `SellMarketWithPips()` - More intuitive, specify SL/TP in pips ⭐ **RECOMMENDED**

**🍬 Other market orders:**

* `SellMarket()` - No SL/TP (NOT recommended for production)
* `BuyMarketWithSLTP()` - BUY with SL/TP

**🍬 Helper methods:**

* `CalculateSLTP()` - Calculate SL/TP prices from pips
* `CalculatePositionSize()` - Risk-based position sizing

---

## ⚠️ Common Pitfalls

### 1) Wrong SL/TP direction

```go
// ❌ WRONG - SL must be ABOVE entry for SELL
bid := 1.08500
sugar.SellMarketWithSLTP("EURUSD", 0.1, 1.08000, 1.07000) // SL below entry!

// ✅ CORRECT - SL above, TP below
sugar.SellMarketWithSLTP("EURUSD", 0.1, 1.09000, 1.08000)
```

### 2) Confusing BID and ASK

```go
// ❌ WRONG - using ASK for SELL (SELL uses BID price)
ask, _ := sugar.GetAsk("EURUSD")
sl := ask + (50 * point)

// ✅ CORRECT - use BID for SELL calculations
bid, _ := sugar.GetBid("EURUSD")
sl := bid + (50 * point)
```

### 3) Using prices instead of pips

```go
// ❌ HARDER - calculate prices manually
bid, _ := sugar.GetBid("EURUSD")
info, _ := sugar.GetSymbolInfo("EURUSD")
sl := bid + (50 * info.Point)
tp := bid - (100 * info.Point)
sugar.SellMarketWithSLTP("EURUSD", 0.1, sl, tp)

// ✅ EASIER - use SellMarketWithPips
sugar.SellMarketWithPips("EURUSD", 0.1, 50, 100)
```

---

## 💎 Pro Tips

1. **Use `SellMarketWithPips()`** - More intuitive than calculating prices

2. **Always set SL** - Never trade without stop loss in production

3. **2:1 minimum R:R** - TP should be at least 2x SL distance

4. **Validate before opening** - Check spread, margin, symbol availability

5. **Calculate position size** - Use `CalculatePositionSize()` for risk management

---

**See also:** [`SellMarketWithPips.md`](../11. Trading_Helpers/SellMarketWithPips.md), [`CalculateSLTP.md`](../11. Trading_Helpers/CalculateSLTP.md), [`BuyMarketWithSLTP.md`](BuyMarketWithSLTP.md)
