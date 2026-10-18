# 📈 Buy at Market with SL/TP in Pips (`BuyMarketWithPips`)

> **Sugar method:** Opens a BUY position with Stop Loss and Take Profit specified as **distances in points** (not price! - and despite the method name, not pips). More intuitive than price-based methods - you specify risk/reward as distances and it calculates exact prices automatically.

**API Information:**

* **Method:** `sugar.BuyMarketWithPips(symbol, volume, stopLossPips, takeProfitPips)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `CalculateSLTP()`, `BuyMarketWithSLTP()`
* **Timeout:** 30 seconds for the order (`BuyMarketWithSLTP`); before that `CalculateSLTP` queries `GetSymbolInfo` (10 s) and `GetAsk` (3 s)

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
| `stopLossPips` | `float64` | Stop Loss distance in **POINTS from entry** despite the name (e.g., 500 = 50 pips on 5-digit EURUSD) |
| `takeProfitPips` | `float64` | Take Profit distance in **POINTS from entry** despite the name (e.g., 1000 = 100 pips on 5-digit EURUSD) |

**Important:** SL/TP are distances in **points** (multiplied by the symbol's `Point`: 0.00001 on 5-digit EURUSD, 0.001 on 3-digit USDJPY), not price and not pips! On 5- and 3-digit symbols 10 points = 1 pip, so pass 500 for a 50-pip stop.

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Order ticket (`result.Order`); for a market order MT5 uses it as the position ticket (use for tracking/closing) |
| `error` | `error` | The `CalculateSLTP` error (symbol/price query failed), `order rejected, code: N, comment: ...` if the return code is not 10009 (TRADE_RETCODE_DONE), or `BuyMarketWithSLTP failed: ...` |

---

## 💬 Just the Essentials

* **What it is:** Opens BUY at market with SL/TP specified as point distances instead of price.
* **Why you need it:** Much easier than calculating SL/TP prices manually. Just specify "500 points stop, 1000 points profit" and you're done.
* **Sanity check:** If EURUSD (5 digits) Ask is 1.08500, SL=500, TP=1000 (points) → SL will be at 1.08000, TP at 1.09500.

---

## 🎯 Purpose

Use it for intuitive market orders:

* **Think in distances, not prices** - More natural for traders
* **Consistent R:R ratios** - Easy to maintain 1:2 or 1:3 R:R
* **Quick market entries** - Specify SL/TP in one line
* **Auto price calculation** - No manual math needed
* **Works across all symbols** - Uses each symbol's `Point` size

---

## 🔢 How It Works

```
Entry Price = Current ASK
Stop Loss   = ASK - (stopLossPips × point)
Take Profit = ASK + (takeProfitPips × point)

Where:
  point = symbol's minimal price change (e.g., 0.00001 for EURUSD)
SL and TP are rounded to the symbol's digits.
```

**Example:**
```
EURUSD ASK = 1.08500
SL = 500 points → 1.08500 - (500 × 0.00001) = 1.08500 - 0.00500 = 1.08000
TP = 1000 points → 1.08500 + (1000 × 0.00001) = 1.08500 + 0.01000 = 1.09500

BUY at 1.08500, SL at 1.08000, TP at 1.09500
```

---

## 🧩 Notes & Tips

* **Points, not price!** - 500 means "500 points" (= 50 pips on 5-digit EURUSD), not "price 500.00"
* **For BUY** - SL is BELOW entry, TP is ABOVE entry
* **Current market price** - Uses current ASK for entry
* **Auto-calculation** - Handles different symbol digits automatically
* **Do NOT pass 0 to skip SL/TP** - a 0 distance puts SL/TP exactly at the entry (Ask) price, which the broker will normally reject; use `BuyMarket()` to open without SL/TP
* **Rounding** - SL/TP prices are rounded to the symbol's `Digits`; the volume is sent as given (not normalized to `VolumeStep`)
* **Order timeout 30 s** - An order that times out on the client side may still have been executed: check `GetOpenPositions()` before retrying
* **Market execution** - Order fills at current market price
* **Slippage possible** - Actual entry might differ slightly; SL/TP are calculated from the Ask read before the order is sent

---

## 🔧 Under the Hood

```go
func (s *MT5Sugar) BuyMarketWithPips(symbol string, volume, stopLossPips, takeProfitPips float64) (uint64, error) {
	// Calculate SL/TP prices
	sl, tp, err := s.CalculateSLTP(symbol, "BUY", 0, stopLossPips, takeProfitPips)
	if err != nil {
		return 0, err
	}

	// Use existing BuyMarketWithSLTP
	return s.BuyMarketWithSLTP(symbol, volume, sl, tp)
}
```

**What it improves:**

* ✅ **Distance-based input** - Think in points, not prices
* ✅ **Auto price calculation** - No manual math
* ✅ **Symbol-aware** - Uses the symbol's `Point` (5-digit vs 3-digit quotes)
* ✅ **One method call** - Instead of three separate steps

---

## 📊 Low-Level Alternative

**WITHOUT sugar (manual calculation):**
```go
// Get symbol info
info, _ := sugar.GetSymbolInfo("EURUSD")

// Get current ASK
ask, _ := sugar.GetAsk("EURUSD")

// Calculate SL/TP manually (and round to the symbol's digits)
scale := math.Pow(10, float64(info.Digits))
sl := math.Round((ask-500*info.Point)*scale) / scale  // 500 points below
tp := math.Round((ask+1000*info.Point)*scale) / scale // 1000 points above

// Open position
ticket, _ := sugar.BuyMarketWithSLTP("EURUSD", 0.1, sl, tp)
```

**WITH sugar:**
```go
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)
```

**Benefits:**

* ✅ **One line vs 5 lines**
* ✅ **No point calculation needed**
* ✅ **No symbol info lookup**
* ✅ **Clearer intent** - "500 and 1000 points" is obvious

---

## 🔗 Usage Examples

### 1) Basic market BUY with 1:2 R:R

```go
// Open BUY: 0.1 lots, 500 point SL, 1000 point TP (1:2 R:R)
ticket, err := sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY order #%d opened\n", ticket)
fmt.Printf("   SL: 500 points, TP: 1000 points (1:2 R:R)\n")
```

---

### 2) Complete risk-managed trade

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLossPips := 500.0
takeProfitPips := 1000.0

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
fmt.Printf("   Stop Loss:   %.0f points\n", stopLossPips)
fmt.Printf("   Take Profit: %.0f points\n", takeProfitPips)
```

---

### 3) Different risk/reward ratios

```go
symbol := "EURUSD"
volume := 0.1

// Conservative: 1:1.5 R:R
ticket1, _ := sugar.BuyMarketWithPips(symbol, volume, 500, 750)
fmt.Println("Conservative: 500 point SL, 750 point TP (1:1.5)")

// Moderate: 1:2 R:R (most common)
ticket2, _ := sugar.BuyMarketWithPips(symbol, volume, 500, 1000)
fmt.Println("Moderate: 500 point SL, 1000 point TP (1:2)")

// Aggressive: 1:3 R:R
ticket3, _ := sugar.BuyMarketWithPips(symbol, volume, 500, 1500)
fmt.Println("Aggressive: 500 point SL, 1500 point TP (1:3)")

// Output:
// Conservative: 500 point SL, 750 point TP (1:1.5)
// Moderate: 500 point SL, 1000 point TP (1:2)
// Aggressive: 500 point SL, 1500 point TP (1:3)
```

---

### 4) Tight stop vs wide stop

```go
symbol := "GBPUSD"
volume := 0.1

// Scalping: tight stop, tight target
ticket1, _ := sugar.BuyMarketWithPips(symbol, volume, 150, 300)
fmt.Println("Scalp: 150 point SL, 300 point TP")

// Day trading: normal stop
ticket2, _ := sugar.BuyMarketWithPips(symbol, volume, 500, 1000)
fmt.Println("Day: 500 point SL, 1000 point TP")

// Swing trading: wide stop
ticket3, _ := sugar.BuyMarketWithPips(symbol, volume, 1500, 3000)
fmt.Println("Swing: 1500 point SL, 3000 point TP")
```

---

### 5) Show entry, SL, TP prices

```go
symbol := "EURUSD"
volume := 0.1
slPips := 500.0
tpPips := 1000.0

// Get current price
ask, _ := sugar.GetAsk(symbol)

// Calculate what SL/TP prices will be
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", ask, slPips, tpPips)

fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("  BUY ORDER PREVIEW\n")
fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("Symbol:       %s\n", symbol)
fmt.Printf("Entry (ASK):  %.5f\n", ask)
fmt.Printf("Stop Loss:    %.5f (-%.0f points)\n", sl, slPips)
fmt.Printf("Take Profit:  %.5f (+%.0f points)\n", tp, tpPips)
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
// Stop Loss:    1.08000 (-500 points)
// Take Profit:  1.09500 (+1000 points)
// Risk/Reward:  1:2.0
// ═══════════════════════════════════════
// ✅ Order #44444 opened
```

---

### 6) Multiple positions with same parameters

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
volume := 0.1
slPips := 500.0
tpPips := 1000.0

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

### 7) Open without SL/TP, then set SL

```go
// Open without SL/TP: use BuyMarket.
// (BuyMarketWithPips(..., 0, 0) would NOT do this - it would send SL = TP = Ask.)
ticket, err := sugar.BuyMarket("EURUSD", 0.1)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Println("⚠️  Order opened WITHOUT Stop Loss or Take Profit")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Println("   🚨 DANGEROUS - Manually set SL immediately!")

// Immediately set stop loss
sl, _, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, 500, 0)
sugar.ModifyPositionSL(ticket, sl)
fmt.Println("✅ Stop Loss set to 500 points")
```

---

### 8) Move SL to breakeven after profit

```go
symbol := "EURUSD"
volume := 0.1
slPips := 500.0
tpPips := 1000.0

// Open position
ticket, _ := sugar.BuyMarketWithPips(symbol, volume, slPips, tpPips)
fmt.Printf("Position #%d opened\n", ticket)

// Monitor and trail
time.Sleep(10 * time.Second) // Wait for price to move

// Get current position (an error if it is no longer open, e.g. SL/TP hit)
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not open: %v\n", err)
    return
}

if pos.Profit > 0 {
    // Price moved in our favor - move SL to breakeven
    fmt.Println("✅ Position in profit - moving SL to breakeven")

    // Calculate breakeven price
    info, _ := sugar.GetSymbolInfo(symbol)
    breakeven := pos.PriceOpen + (50 * info.Point) // BE + 50 points

    sugar.ModifyPositionSL(ticket, breakeven)
    fmt.Printf("   SL moved to %.5f (breakeven + 50 points)\n", breakeven)
}
```

---

### 9) Partial close at TP1, let rest run to TP2

```go
symbol := "EURUSD"
volume := 0.2  // Total position size
slPips := 500.0
tp1Pips := 500.0  // First target
tp2Pips := 1500.0 // Final target

// Open position
ticket, _ := sugar.BuyMarketWithPips(symbol, volume, slPips, tp1Pips)
fmt.Printf("Position #%d opened: %.2f lots\n", ticket, volume)
fmt.Printf("   SL: %.0f points\n", slPips)
fmt.Printf("   TP1: %.0f points (will close 50%%)\n", tp1Pips)
fmt.Printf("   TP2: %.0f points (remaining 50%%)\n", tp2Pips)

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
slPips := 3000.0  // Gold needs wider stops
tpPips := 6000.0

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
* `CalculateSLTP()` - Calculate SL/TP prices from point distances
* `BuyMarketWithSLTP()` - BUY with prices instead of distances
* `ModifyPositionSLTP()` - Change SL/TP after opening

**📦 Methods used internally:**

* `CalculateSLTP()` - Converts point distances to prices
* `BuyMarketWithSLTP()` - Executes the actual order

**Recommended workflow:**
```go
// 1. Calculate position size
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)

// 2. Validate
canOpen, _, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen { return }

// 3. Trade
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 500, 1000)
```

---

## ⚠️ Common Pitfalls

### 1) Confusing distance with price

```go
// ❌ WRONG - passing price as stop loss
slPrice := 1.08000
sugar.BuyMarketWithPips("EURUSD", 0.1, slPrice, 1000) // WRONG!

// ❌ WRONG - passing pips: 50 is read as 50 POINTS (= 5 pips)
sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)

// ✅ CORRECT - stop loss distance in points (500 points = 50 pips on EURUSD)
slPips := 500.0
sugar.BuyMarketWithPips("EURUSD", 0.1, slPips, 1000)
```

### 2) Using same distances for different symbols

```go
// ❌ WRONG - 500 points SL works for EURUSD but not XAUUSD
sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)  // OK
sugar.BuyMarketWithPips("XAUUSD", 0.1, 500, 1000)  // Too tight!

// ✅ CORRECT - adjust distances based on symbol volatility
sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)    // Forex: 500 points
sugar.BuyMarketWithPips("XAUUSD", 0.1, 3000, 6000)   // Gold: 3000 points
```

### 3) Not using risk-based position sizing

```go
// ❌ WRONG - fixed lot size (ignores balance)
sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)

// ✅ CORRECT - calculate lot size based on risk
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
sugar.BuyMarketWithPips("EURUSD", lotSize, 500, 1000)
```

### 4) Not validating before trading

```go
// ❌ WRONG - trading without validation
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 10.0, 500, 1000) // Might fail!

// ✅ CORRECT - validate first
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 10.0)
if !canOpen {
    fmt.Println("Cannot trade:", reason)
    return
}
sugar.BuyMarketWithPips("EURUSD", 10.0, 500, 1000)
```

### 5) Ignoring errors

```go
// ❌ WRONG - ignoring errors
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)
fmt.Printf("Order #%d opened\n", ticket) // Might be 0!

// ✅ CORRECT - check errors
ticket, err := sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}
fmt.Printf("Order #%d opened\n", ticket)
```

### 6) Unrealistic stops for the symbol

```go
// ❌ WRONG - 50 point SL on volatile pair
sugar.BuyMarketWithPips("GBPJPY", 0.1, 50, 100) // Will get stopped immediately!

// ✅ CORRECT - use appropriate stop distance
sugar.BuyMarketWithPips("GBPJPY", 0.1, 800, 1600) // Realistic for volatility
```

---

## 💎 Pro Tips

1. **Use 1:2 R:R minimum** - Never trade with less than 1:2 risk/reward
2. **Adjust distances to symbol** - Volatile pairs need wider stops
3. **Always calculate lot size** - Use `CalculatePositionSize()` first
4. **Validate before trading** - Use `CanOpenPosition()`
5. **Monitor spread** - High spread = wider stops needed
6. **Think in distances, not price** - This method makes it natural (distances are in points)
7. **Set realistic stops** - Don't use super tight stops that always hit

---

**See also:** [`SellMarketWithPips.md`](SellMarketWithPips.md), [`CalculatePositionSize.md`](../10.%20Risk_Management/CalculatePositionSize.md), [`CanOpenPosition.md`](../10.%20Risk_Management/CanOpenPosition.md)
