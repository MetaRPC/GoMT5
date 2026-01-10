# 📈 Buy at Market with SL/TP in Pips (`BuyMarketWithPips`)

> **Sugar method:** Opens a BUY position with Stop Loss and Take Profit specified in **pips** (not price!). More intuitive than price-based methods - you specify risk/reward in pips and it calculates exact prices automatically.

**API Information:**

* **Method:** `sugar.BuyMarketWithPips(symbol, volume, stopLossPips, takeProfitPips)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `CalculateSLTP()`, `BuyMarketWithSLTP()`
* **Timeout:** 10 seconds

---

## 📋 Method Signature

```go
func (s *MT5Sugar) BuyMarketWithPips(
    symbol string,
    volume float64,
    stopLossPips float64,
    takeProfitPips float64,
) (uint64, error)
```

---

## 🔽 Input

| Parameter | Type | Description |
|-----------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1, 0.01, 1.0) |
| `stopLossPips` | `float64` | Stop Loss distance in **pips from entry** (e.g., 50) |
| `takeProfitPips` | `float64` | Take Profit distance in **pips from entry** (e.g., 100) |

**Important:** SL/TP are in **pips** (points), not price!

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Position ticket number (use for tracking/closing) |
| `error` | `error` | Error if order rejected or execution failed |

---

## 💬 Just the Essentials

* **What it is:** Opens BUY at market with SL/TP specified in pips instead of price.
* **Why you need it:** Much easier than calculating SL/TP prices manually. Just specify "50 pips stop, 100 pips profit" and you're done.
* **Sanity check:** If EURUSD is at 1.08500, SL=50, TP=100 → SL will be at 1.08000, TP at 1.09000.

---

## 🎯 Purpose

Use it for intuitive market orders:

* **Think in pips, not prices** - More natural for traders
* **Consistent R:R ratios** - Easy to maintain 1:2 or 1:3 R:R
* **Quick market entries** - Specify SL/TP in one line
* **Auto price calculation** - No manual math needed
* **Works across all symbols** - Handles different pip values automatically

---

## 🔢 How It Works

```
Entry Price = Current ASK
Stop Loss   = ASK - (stopLossPips × point)
Take Profit = ASK + (takeProfitPips × point)

Where:
  point = symbol's minimal price change (e.g., 0.00001 for EURUSD)
```

**Example:**
```
EURUSD ASK = 1.08500
SL = 50 pips → 1.08500 - (50 × 0.00010) = 1.08000
TP = 100 pips → 1.08500 + (100 × 0.00010) = 1.09000

BUY at 1.08500, SL at 1.08000, TP at 1.09000
```

---

## 🧩 Notes & Tips

* **Pips, not price!** - 50 means "50 points" not "price 50.00"
* **For BUY** - SL is BELOW entry, TP is ABOVE entry
* **Current market price** - Uses current ASK for entry
* **Auto-calculation** - Handles different symbol digits automatically
* **Set both to 0** - To open without SL/TP (not recommended!)
* **Market execution** - Order fills at current market price
* **Slippage possible** - Actual entry might differ slightly

---

## 🔧 Under the Hood

```go
func (s *MT5Sugar) BuyMarketWithPips(symbol string, volume, stopLossPips, takeProfitPips float64) (uint64, error) {
    // 1. Calculate exact SL/TP prices from pips
    sl, tp, err := s.CalculateSLTP(symbol, "BUY", 0, stopLossPips, takeProfitPips)
    if err != nil {
        return 0, err
    }

    // 2. Open BUY position with calculated prices
    return s.BuyMarketWithSLTP(symbol, volume, sl, tp)
}
```

**What it improves:**

* ✅ **Pip-based input** - Think in pips, not prices
* ✅ **Auto price calculation** - No manual math
* ✅ **Symbol-aware** - Handles 5-digit vs 3-digit brokers
* ✅ **One method call** - Instead of three separate steps

---

## 📊 Low-Level Alternative

**WITHOUT sugar (manual calculation):**
```go
// Get symbol info
info, _ := sugar.GetSymbolInfo("EURUSD")

// Get current ASK
ask, _ := sugar.GetAsk("EURUSD")

// Calculate SL/TP manually
sl := ask - (50 * info.Point)  // 50 pips below
tp := ask + (100 * info.Point) // 100 pips above

// Open position
ticket, _ := sugar.BuyMarketWithSLTP("EURUSD", 0.1, sl, tp)
```

**WITH sugar:**
```go
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
```

**Benefits:**

* ✅ **One line vs 5 lines**
* ✅ **No point calculation needed**
* ✅ **No symbol info lookup**
* ✅ **Clearer intent** - "50 and 100 pips" is obvious

---

## 🔗 Usage Examples

### 1) Basic market BUY with 1:2 R:R

```go
// Open BUY: 0.1 lots, 50 pip SL, 100 pip TP (1:2 R:R)
ticket, err := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY order #%d opened\n", ticket)
fmt.Printf("   SL: 50 pips, TP: 100 pips (1:2 R:R)\n")
```

---

### 2) Complete risk-managed trade

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLossPips := 50.0
takeProfitPips := 100.0

// Step 1: Calculate position size based on risk
lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)
if err != nil {
    fmt.Printf("Position size calculation failed: %v\n", err)
    return
}

// Step 2: Validate we can open
canOpen, reason, _ := sugar.CanOpenPosition(symbol, lotSize)
if !canOpen {
    fmt.Printf("Cannot open: %s\n", reason)
    return
}

// Step 3: Open position with calculated size
ticket, err := sugar.BuyMarketWithPips(symbol, lotSize, stopLossPips, takeProfitPips)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Trade opened successfully!\n")
fmt.Printf("   Ticket:      #%d\n", ticket)
fmt.Printf("   Lot size:    %.2f (risk: %.1f%%)\n", lotSize, riskPercent)
fmt.Printf("   Stop Loss:   %.0f pips\n", stopLossPips)
fmt.Printf("   Take Profit: %.0f pips\n", takeProfitPips)
```

---

### 3) Different risk/reward ratios

```go
symbol := "EURUSD"
volume := 0.1

// Conservative: 1:1.5 R:R
ticket1, _ := sugar.BuyMarketWithPips(symbol, volume, 50, 75)
fmt.Println("Conservative: 50 pip SL, 75 pip TP (1:1.5)")

// Moderate: 1:2 R:R (most common)
ticket2, _ := sugar.BuyMarketWithPips(symbol, volume, 50, 100)
fmt.Println("Moderate: 50 pip SL, 100 pip TP (1:2)")

// Aggressive: 1:3 R:R
ticket3, _ := sugar.BuyMarketWithPips(symbol, volume, 50, 150)
fmt.Println("Aggressive: 50 pip SL, 150 pip TP (1:3)")

// Output:
// Conservative: 50 pip SL, 75 pip TP (1:1.5)
// Moderate: 50 pip SL, 100 pip TP (1:2)
// Aggressive: 50 pip SL, 150 pip TP (1:3)
```

---

### 4) Tight stop vs wide stop

```go
symbol := "GBPUSD"
volume := 0.1

// Scalping: tight stop, tight target
ticket1, _ := sugar.BuyMarketWithPips(symbol, volume, 15, 30)
fmt.Println("Scalp: 15 pip SL, 30 pip TP")

// Day trading: normal stop
ticket2, _ := sugar.BuyMarketWithPips(symbol, volume, 50, 100)
fmt.Println("Day: 50 pip SL, 100 pip TP")

// Swing trading: wide stop
ticket3, _ := sugar.BuyMarketWithPips(symbol, volume, 150, 300)
fmt.Println("Swing: 150 pip SL, 300 pip TP")
```

---

### 5) Show entry, SL, TP prices

```go
symbol := "EURUSD"
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Get current price
ask, _ := sugar.GetAsk(symbol)

// Calculate what SL/TP prices will be
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", ask, slPips, tpPips)

fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("  BUY ORDER PREVIEW\n")
fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("Symbol:       %s\n", symbol)
fmt.Printf("Entry (ASK):  %.5f\n", ask)
fmt.Printf("Stop Loss:    %.5f (-%​.0f pips)\n", sl, slPips)
fmt.Printf("Take Profit:  %.5f (+%.0f pips)\n", tp, tpPips)
fmt.Printf("Risk/Reward:  1:%.1f\n", tpPips/slPips)
fmt.Printf("═══════════════════════════════════════\n")

// Open position
ticket, _ := sugar.BuyMarketWithPips(symbol, volume, slPips, tpPips)
fmt.Printf("✅ Order #%d opened\n", ticket)

// Output:
// ═══════════════════════════════════════
//   BUY ORDER PREVIEW
// ═══════════════════════════════════════
// Symbol:       EURUSD
// Entry (ASK):  1.08500
// Stop Loss:    1.08000 (-50 pips)
// Take Profit:  1.09000 (+100 pips)
// Risk/Reward:  1:2.0
// ═══════════════════════════════════════
// ✅ Order #44444 opened
```

---

### 6) Multiple positions with same parameters

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
volume := 0.1
slPips := 50.0
tpPips := 100.0

fmt.Println("Opening BUY positions on multiple pairs:")
fmt.Println("─────────────────────────────────────────")

for _, symbol := range symbols {
    ticket, err := sugar.BuyMarketWithPips(symbol, volume, slPips, tpPips)
    if err != nil {
        fmt.Printf("❌ %s: Failed - %v\n", symbol, err)
        continue
    }

    fmt.Printf("✅ %s: Ticket #%d\n", symbol, ticket)
}

// Output:
// Opening BUY positions on multiple pairs:
// ─────────────────────────────────────────
// ✅ EURUSD: Ticket #44444
// ✅ GBPUSD: Ticket #44445
// ✅ USDJPY: Ticket #44446
```

---

### 7) Breakeven entry (no SL/TP)

```go
// Open without SL/TP (pass 0 for both)
ticket, err := sugar.BuyMarketWithPips("EURUSD", 0.1, 0, 0)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Println("⚠️  Order opened WITHOUT Stop Loss or Take Profit")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Println("   🚨 DANGEROUS - Manually set SL immediately!")

// Immediately set stop loss
sl, _, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, 50, 0)
sugar.ModifyPositionSL(ticket, sl)
fmt.Println("✅ Stop Loss set to 50 pips")
```

---

### 8) Trail after X pips profit

```go
symbol := "EURUSD"
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Open position
ticket, _ := sugar.BuyMarketWithPips(symbol, volume, slPips, tpPips)
fmt.Printf("Position #%d opened\n", ticket)

// Monitor and trail
time.Sleep(10 * time.Second) // Wait for price to move

// Get current position
pos, _ := sugar.GetPositionByTicket(ticket)

if pos.Profit > 0 {
    // Price moved in our favor - move SL to breakeven
    fmt.Println("✅ Position in profit - moving SL to breakeven")

    // Calculate breakeven price
    info, _ := sugar.GetSymbolInfo(symbol)
    breakeven := pos.PriceOpen + (5 * info.Point) // BE + 5 pips

    sugar.ModifyPositionSL(ticket, breakeven)
    fmt.Printf("   SL moved to %.5f (breakeven + 5 pips)\n", breakeven)
}
```

---

### 9) Partial close at TP1, let rest run to TP2

```go
symbol := "EURUSD"
volume := 0.2  // Total position size
slPips := 50.0
tp1Pips := 50.0  // First target
tp2Pips := 150.0 // Final target

// Open position
ticket, _ := sugar.BuyMarketWithPips(symbol, volume, slPips, tp1Pips)
fmt.Printf("Position #%d opened: %.2f lots\n", ticket, volume)
fmt.Printf("   SL: %.0f pips\n", slPips)
fmt.Printf("   TP1: %.0f pips (will close 50%%)\n", tp1Pips)
fmt.Printf("   TP2: %.0f pips (remaining 50%%)\n", tp2Pips)

// ... Later, when price reaches TP1 ...
// Close half position
halfVolume := volume / 2
err := sugar.ClosePositionPartial(ticket, halfVolume)
if err == nil {
    fmt.Println("✅ Closed 50% at TP1")

    // Move SL to breakeven on remaining position
    pos, _ := sugar.GetPositionByTicket(ticket)
    sugar.ModifyPositionSL(ticket, pos.PriceOpen)

    // Set new TP to TP2
    _, tp2, _ := sugar.CalculateSLTP(symbol, "BUY", 0, 0, tp2Pips)
    sugar.ModifyPositionTP(ticket, tp2)

    fmt.Println("✅ SL moved to breakeven, TP set to TP2")
}
```

---

### 10) Validate symbol before trading

```go
symbol := "XAUUSD"
volume := 0.1
slPips := 300.0  // Gold needs wider stops
tpPips := 600.0

// Check symbol availability
available, err := sugar.IsSymbolAvailable(symbol)
if err != nil || !available {
    fmt.Printf("❌ %s is not available for trading\n", symbol)
    return
}

// Get symbol info to show specifications
info, _ := sugar.GetSymbolInfo(symbol)
fmt.Printf("Symbol: %s\n", info.Name)
fmt.Printf("  Spread:     %d points\n", info.Spread)
fmt.Printf("  Min volume: %.2f\n", info.VolumeMin)
fmt.Printf("  Max volume: %.2f\n", info.VolumeMax)

// Validate volume
if volume < info.VolumeMin {
    fmt.Printf("❌ Volume %.2f below minimum %.2f\n", volume, info.VolumeMin)
    return
}

// Open position
ticket, err := sugar.BuyMarketWithPips(symbol, volume, slPips, tpPips)
if err != nil {
    fmt.Printf("❌ Order failed: %v\n", err)
    return
}

fmt.Printf("✅ %s BUY #%d opened\n", symbol, ticket)
```

---

## 🔗 Related Methods

**🍬 Complementary sugar methods:**

* `SellMarketWithPips()` - SELL version of this method
* `CalculatePositionSize()` - Calculate lot size based on risk ⭐
* `CanOpenPosition()` - Validate before trading ⭐
* `CalculateSLTP()` - Calculate SL/TP prices from pips
* `BuyMarketWithSLTP()` - BUY with prices instead of pips
* `ModifyPositionSLTP()` - Change SL/TP after opening

**📦 Methods used internally:**

* `CalculateSLTP()` - Converts pips to prices
* `BuyMarketWithSLTP()` - Executes the actual order

**Recommended workflow:**
```go
// 1. Calculate position size
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// 2. Validate
canOpen, _, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen { return }

// 3. Trade
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

---

## ⚠️ Common Pitfalls

### 1) Confusing pips with price

```go
// ❌ WRONG - passing price as stop loss
slPrice := 1.08000
sugar.BuyMarketWithPips("EURUSD", 0.1, slPrice, 100) // WRONG!

// ✅ CORRECT - stop loss in pips
slPips := 50.0
sugar.BuyMarketWithPips("EURUSD", 0.1, slPips, 100)
```

### 2) Using same pip values for different symbols

```go
// ❌ WRONG - 50 pips SL works for EURUSD but not XAUUSD
sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)  // OK
sugar.BuyMarketWithPips("XAUUSD", 0.1, 50, 100)  // Too tight!

// ✅ CORRECT - adjust pips based on symbol volatility
sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)    // Forex: 50 pips
sugar.BuyMarketWithPips("XAUUSD", 0.1, 300, 600)   // Gold: 300 pips
```

### 3) Not using risk-based position sizing

```go
// ❌ WRONG - fixed lot size (ignores balance)
sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)

// ✅ CORRECT - calculate lot size based on risk
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

### 4) Not validating before trading

```go
// ❌ WRONG - trading without validation
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 10.0, 50, 100) // Might fail!

// ✅ CORRECT - validate first
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 10.0)
if !canOpen {
    fmt.Println("Cannot trade:", reason)
    return
}
sugar.BuyMarketWithPips("EURUSD", 10.0, 50, 100)
```

### 5) Ignoring errors

```go
// ❌ WRONG - ignoring errors
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
fmt.Printf("Order #%d opened\n", ticket) // Might be 0!

// ✅ CORRECT - check errors
ticket, err := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}
fmt.Printf("Order #%d opened\n", ticket)
```

### 6) Unrealistic stops for the symbol

```go
// ❌ WRONG - 5 pip SL on volatile pair
sugar.BuyMarketWithPips("GBPJPY", 0.1, 5, 10) // Will get stopped immediately!

// ✅ CORRECT - use appropriate stop distance
sugar.BuyMarketWithPips("GBPJPY", 0.1, 80, 160) // Realistic for volatility
```

---

## 💎 Pro Tips

1. **Use 1:2 R:R minimum** - Never trade with less than 1:2 risk/reward
2. **Adjust pips to symbol** - Volatile pairs need wider stops
3. **Always calculate lot size** - Use `CalculatePositionSize()` first
4. **Validate before trading** - Use `CanOpenPosition()`
5. **Monitor spread** - High spread = wider stops needed
6. **Think in pips, not price** - This method makes it natural
7. **Set realistic stops** - Don't use super tight stops that always hit

---

**See also:** [`SellMarketWithPips.md`](SellMarketWithPips.md), [`CalculatePositionSize.md`](../10.%20Risk_Management/CalculatePositionSize.md), [`CanOpenPosition.md`](../10.%20Risk_Management/CanOpenPosition.md)
