# 🎯 Calculate SL/TP Prices (`CalculateSLTP`)

> **Sugar method:** Converts pip distances to actual SL/TP prices - handles the math for you.

**API Information:**

* **Method:** `sugar.CalculateSLTP(symbol, direction, entryPrice, stopLossPips, takeProfitPips)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `GetSymbolInfo()`, `GetAsk()`, `GetBid()`
* **Timeout:** 3 seconds
* **Returns:** SL price, TP price

---

## 📋 Method Signature

```go
func (s *MT5Sugar) CalculateSLTP(
    symbol string,
    direction string,
    entryPrice float64,
    stopLossPips float64,
    takeProfitPips float64,
) (float64, float64, error)
```

---

## 🔽 Input

| Parameter | Type | Description |
|-----------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| `direction` | `string` | "BUY" or "SELL" |
| `entryPrice` | `float64` | Entry price (use `0` for current market price) |
| `stopLossPips` | `float64` | Stop Loss distance in **points** (e.g., 50) |
| `takeProfitPips` | `float64` | Take Profit distance in **points** (e.g., 100) |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `sl` | `float64` | Stop Loss price |
| `tp` | `float64` | Take Profit price |
| `error` | `error` | Error if calculation fails |

---

## 💬 Just the Essentials

* **What it is:** Converts pip distances into actual SL/TP prices.
* **Why you need it:** You think in pips ("50 pip stop"), but MT5 needs exact prices.
* **Sanity check:** Handles BUY vs SELL direction automatically - BUY SL below entry, SELL SL above entry.

---

## 🎯 When to Use

✅ **Convert pips to prices** - "50 pip SL" → actual price

✅ **Plan trades** - Calculate SL/TP before opening

✅ **Modify positions** - Adjust SL/TP by pip distance

✅ **Risk/Reward planning** - Work with pip-based R:R ratios

---

## 🔢 How It Works

```
BUY direction:
- Entry:  1.08500
- SL:     Entry - (50 pips × point size) = 1.08000 (below entry)
- TP:     Entry + (100 pips × point size) = 1.09000 (above entry)

SELL direction:
- Entry:  1.08500
- SL:     Entry + (50 pips × point size) = 1.09000 (above entry)
- TP:     Entry - (100 pips × point size) = 1.08000 (below entry)

Point size = symbol.Point (e.g., 0.00001 for EURUSD with 5 digits)
```

---

## 🔗 Usage Examples

### 1) Basic usage - calculate SL/TP for BUY

```go
symbol := "EURUSD"
direction := "BUY"
entryPrice := 0.0      // 0 = use current market price
stopLoss := 50.0       // 50 pips
takeProfit := 100.0    // 100 pips

sl, tp, err := sugar.CalculateSLTP(symbol, direction, entryPrice, stopLoss, takeProfit)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Direction:   %s\n", direction)
fmt.Printf("Stop Loss:   %.5f (-%0.f pips)\n", sl, stopLoss)
fmt.Printf("Take Profit: %.5f (+%.0f pips)\n", tp, takeProfit)

// Output example:
// Direction:   BUY
// Stop Loss:   1.08000 (-50 pips)
// Take Profit: 1.09000 (+100 pips)
```

---

### 2) Calculate for SELL position

```go
sl, tp, err := sugar.CalculateSLTP("EURUSD", "SELL", 0, 50, 100)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Get current price for reference
bid, _ := sugar.GetBid("EURUSD")

fmt.Printf("SELL @ %.5f\n", bid)
fmt.Printf("SL:    %.5f (+50 pips above entry)\n", sl)
fmt.Printf("TP:    %.5f (-100 pips below entry)\n", tp)
```

---

### 3) Use specific entry price (not current market)

```go
// Plan a trade with specific entry
symbol := "EURUSD"
plannedEntry := 1.08500
stopLoss := 50.0
takeProfit := 100.0

sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", plannedEntry, stopLoss, takeProfit)

fmt.Printf("Planned trade:\n")
fmt.Printf("Entry:       %.5f\n", plannedEntry)
fmt.Printf("Stop Loss:   %.5f (distance: %.0f pips)\n", sl, stopLoss)
fmt.Printf("Take Profit: %.5f (distance: %.0f pips)\n", tp, takeProfit)
fmt.Printf("Risk/Reward: 1:%.1f\n", takeProfit/stopLoss)

// Output:
// Planned trade:
// Entry:       1.08500
// Stop Loss:   1.08000 (distance: 50 pips)
// Take Profit: 1.09000 (distance: 100 pips)
// Risk/Reward: 1:2.0
```

---

### 4) Complete trading workflow with SL/TP calculation

```go
symbol := "EURUSD"
direction := "BUY"
volume := 0.1
stopLoss := 50.0
takeProfit := 100.0

// Step 1: Calculate SL/TP prices
sl, tp, err := sugar.CalculateSLTP(symbol, direction, 0, stopLoss, takeProfit)
if err != nil {
    fmt.Printf("SL/TP calculation failed: %v\n", err)
    return
}

// Step 2: Validate position
canOpen, reason, err := sugar.CanOpenPosition(symbol, volume)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
    return
}

if !canOpen {
    fmt.Printf("Cannot open: %s\n", reason)
    return
}

// Step 3: Open position
var ticket uint64
if direction == "BUY" {
    ticket, err = sugar.BuyMarketWithSLTP(symbol, volume, sl, tp)
} else {
    ticket, err = sugar.SellMarketWithSLTP(symbol, volume, sl, tp)
}

if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Position #%d opened\n", ticket)
fmt.Printf("   SL: %.5f, TP: %.5f\n", sl, tp)
```

---

### 5) Different R:R ratios

```go
symbol := "EURUSD"
direction := "BUY"
stopLoss := 50.0

// Test different R:R ratios
ratios := []float64{1.0, 1.5, 2.0, 3.0}

fmt.Printf("SL/TP calculations for different R:R ratios (SL=%0.f pips):\n", stopLoss)
fmt.Println("─────────────────────────────────────────")

for _, ratio := range ratios {
    takeProfit := stopLoss * ratio

    sl, tp, _ := sugar.CalculateSLTP(symbol, direction, 0, stopLoss, takeProfit)

    fmt.Printf("R:R 1:%.1f → SL: %.5f, TP: %.5f (%.0f pips)\n",
        ratio, sl, tp, takeProfit)
}

// Output example:
// SL/TP calculations for different R:R ratios (SL=50 pips):
// ─────────────────────────────────────────
// R:R 1:1.0 → SL: 1.08000, TP: 1.08500 (50 pips)
// R:R 1:1.5 → SL: 1.08000, TP: 1.08750 (75 pips)
// R:R 1:2.0 → SL: 1.08000, TP: 1.09000 (100 pips)
// R:R 1:3.0 → SL: 1.08000, TP: 1.09500 (150 pips)
```

---

### 6) Show SL/TP distances from current price

```go
func ShowSLTPDistances(sugar *mt5.MT5Sugar, symbol, direction string, slPips, tpPips float64) {
    sl, tp, _ := sugar.CalculateSLTP(symbol, direction, 0, slPips, tpPips)

    // Get current price
    var currentPrice float64
    if direction == "BUY" {
        currentPrice, _ = sugar.GetAsk(symbol)
    } else {
        currentPrice, _ = sugar.GetBid(symbol)
    }

    // Calculate distances
    info, _ := sugar.GetSymbolInfo(symbol)
    slDistance := math.Abs(currentPrice-sl) / info.Point
    tpDistance := math.Abs(currentPrice-tp) / info.Point

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      SL/TP CALCULATION                ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Symbol:    %s\n", symbol)
    fmt.Printf("Direction: %s\n\n", direction)

    fmt.Printf("Entry:     %.5f (current market)\n\n", currentPrice)

    fmt.Printf("Stop Loss:\n")
    fmt.Printf("  Price:   %.5f\n", sl)
    fmt.Printf("  Distance: %.0f pips\n", slDistance)
    fmt.Printf("  Direction: ")
    if sl < currentPrice {
        fmt.Println("Below entry ⬇️")
    } else {
        fmt.Println("Above entry ⬆️")
    }

    fmt.Printf("\nTake Profit:\n")
    fmt.Printf("  Price:   %.5f\n", tp)
    fmt.Printf("  Distance: %.0f pips\n", tpDistance)
    fmt.Printf("  Direction: ")
    if tp > currentPrice {
        fmt.Println("Above entry ⬆️")
    } else {
        fmt.Println("Below entry ⬇️")
    }

    fmt.Printf("\nRisk/Reward: 1:%.1f\n", tpPips/slPips)
}

// Usage:
ShowSLTPDistances(sugar, "EURUSD", "BUY", 50, 100)
```

---

### 7) Asymmetric SL/TP (different distances)

```go
symbol := "GBPUSD"
direction := "SELL"

// Tight SL, wide TP (aggressive)
tightSL := 30.0
wideTP := 150.0

sl1, tp1, _ := sugar.CalculateSLTP(symbol, direction, 0, tightSL, wideTP)

// Wide SL, tight TP (conservative)
wideSL := 100.0
tightTP := 50.0

sl2, tp2, _ := sugar.CalculateSLTP(symbol, direction, 0, wideSL, tightTP)

fmt.Println("Trade plan comparison:")
fmt.Println()

fmt.Printf("Aggressive (tight SL, wide TP):\n")
fmt.Printf("  SL: %.5f (%.0f pips)\n", sl1, tightSL)
fmt.Printf("  TP: %.5f (%.0f pips)\n", tp1, wideTP)
fmt.Printf("  R:R: 1:%.1f\n", wideTP/tightSL)
fmt.Println()

fmt.Printf("Conservative (wide SL, tight TP):\n")
fmt.Printf("  SL: %.5f (%.0f pips)\n", sl2, wideSL)
fmt.Printf("  TP: %.5f (%.0f pips)\n", tp2, tightTP)
fmt.Printf("  R:R: 1:%.1f\n", tightTP/wideSL)
```

---

### 8) Modify existing position SL/TP

```go
ticket := uint64(123456)

// Get current position
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

// Calculate new SL/TP from current price
newSLPips := 30.0   // Tighter stop
newTPPips := 150.0  // Wider target

direction := "BUY"
if pos.Type == 1 { // SELL
    direction = "SELL"
}

// Use current position price as entry
newSL, newTP, _ := sugar.CalculateSLTP(
    pos.Symbol,
    direction,
    pos.OpenPrice,
    newSLPips,
    newTPPips,
)

fmt.Printf("Modifying position #%d:\n", ticket)
fmt.Printf("Current SL: %.5f → New SL: %.5f\n", pos.StopLoss, newSL)
fmt.Printf("Current TP: %.5f → New TP: %.5f\n", pos.TakeProfit, newTP)

// Modify position
err = sugar.ModifyPosition(ticket, newSL, newTP)
if err != nil {
    fmt.Printf("Modification failed: %v\n", err)
} else {
    fmt.Println("✅ Position modified")
}
```

---

### 9) Multi-symbol SL/TP calculator

```go
type TradeSetup struct {
    Symbol     string
    Direction  string
    SLPips     float64
    TPPips     float64
}

func CalculateMultipleSetups(sugar *mt5.MT5Sugar, setups []TradeSetup) {
    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          TRADE SETUPS CALCULATOR                      ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")

    for i, setup := range setups {
        sl, tp, err := sugar.CalculateSLTP(
            setup.Symbol,
            setup.Direction,
            0,
            setup.SLPips,
            setup.TPPips,
        )

        if err != nil {
            fmt.Printf("%d. %s %s: Error - %v\n",
                i+1, setup.Direction, setup.Symbol, err)
            continue
        }

        // Get current price
        var currentPrice float64
        if setup.Direction == "BUY" {
            currentPrice, _ = sugar.GetAsk(setup.Symbol)
        } else {
            currentPrice, _ = sugar.GetBid(setup.Symbol)
        }

        rr := setup.TPPips / setup.SLPips

        fmt.Printf("\n%d. %s %s @ %.5f\n",
            i+1, setup.Direction, setup.Symbol, currentPrice)
        fmt.Printf("   SL: %.5f (%.0f pips)\n", sl, setup.SLPips)
        fmt.Printf("   TP: %.5f (%.0f pips)\n", tp, setup.TPPips)
        fmt.Printf("   R:R: 1:%.1f\n", rr)
    }
}

// Usage:
setups := []TradeSetup{
    {"EURUSD", "BUY", 50, 100},
    {"GBPUSD", "SELL", 60, 120},
    {"USDJPY", "BUY", 40, 120},
}

CalculateMultipleSetups(sugar, setups)
```

---

### 10) Advanced SL/TP helper

```go
type SLTPHelper struct {
    sugar *mt5.MT5Sugar
}

func NewSLTPHelper(sugar *mt5.MT5Sugar) *SLTPHelper {
    return &SLTPHelper{sugar: sugar}
}

func (h *SLTPHelper) CalculateWithRR(
    symbol, direction string,
    entryPrice, slPips, rrRatio float64,
) (float64, float64, error) {
    // Calculate TP based on R:R ratio
    tpPips := slPips * rrRatio

    return h.sugar.CalculateSLTP(symbol, direction, entryPrice, slPips, tpPips)
}

func (h *SLTPHelper) CalculateTrailingStop(
    symbol string,
    direction string,
    entryPrice float64,
    currentPrice float64,
    trailingPips float64,
) (float64, error) {
    info, err := h.sugar.GetSymbolInfo(symbol)
    if err != nil {
        return 0, err
    }

    var newSL float64

    if direction == "BUY" {
        // For BUY, trailing stop moves up
        newSL = currentPrice - (trailingPips * info.Point)
    } else {
        // For SELL, trailing stop moves down
        newSL = currentPrice + (trailingPips * info.Point)
    }

    return newSL, nil
}

func (h *SLTPHelper) CalculateBreakEven(
    symbol, direction string,
    entryPrice float64,
    bufferPips float64,
) (float64, error) {
    info, err := h.sugar.GetSymbolInfo(symbol)
    if err != nil {
        return 0, err
    }

    var breakEvenSL float64

    if direction == "BUY" {
        // Move SL above entry by buffer
        breakEvenSL = entryPrice + (bufferPips * info.Point)
    } else {
        // Move SL below entry by buffer
        breakEvenSL = entryPrice - (bufferPips * info.Point)
    }

    return breakEvenSL, nil
}

func (h *SLTPHelper) ShowCalculation(
    symbol, direction string,
    slPips, tpPips float64,
) {
    sl, tp, _ := h.sugar.CalculateSLTP(symbol, direction, 0, slPips, tpPips)

    var entry float64
    if direction == "BUY" {
        entry, _ = h.sugar.GetAsk(symbol)
    } else {
        entry, _ = h.sugar.GetBid(symbol)
    }

    // Calculate monetary risk (1 lot)
    pipValue, _ := h.calculatePipValue(symbol)
    riskAmount := slPips * pipValue
    rewardAmount := tpPips * pipValue

    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          SL/TP CALCULATION DETAILS                    ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")
    fmt.Printf("Symbol:         %s\n", symbol)
    fmt.Printf("Direction:      %s\n\n", direction)

    fmt.Printf("Entry price:    %.5f\n", entry)
    fmt.Printf("Stop Loss:      %.5f (%.0f pips)\n", sl, slPips)
    fmt.Printf("Take Profit:    %.5f (%.0f pips)\n\n", tp, tpPips)

    fmt.Printf("Risk/Reward:    1:%.1f\n", tpPips/slPips)
    fmt.Printf("Risk (1 lot):   $%.2f\n", riskAmount)
    fmt.Printf("Reward (1 lot): $%.2f\n", rewardAmount)
}

func (h *SLTPHelper) calculatePipValue(symbol string) (float64, error) {
    info, err := h.sugar.GetSymbolInfo(symbol)
    if err != nil {
        return 0, err
    }

    // Simplified pip value (point × 10 × contract size)
    return info.Point * 10 * info.ContractSize, nil
}

// Usage:
helper := NewSLTPHelper(sugar)

// Calculate with R:R ratio
sl, tp, _ := helper.CalculateWithRR("EURUSD", "BUY", 0, 50, 2.0)
fmt.Printf("SL: %.5f, TP: %.5f (1:2 R:R)\n", sl, tp)

// Calculate trailing stop
newSL, _ := helper.CalculateTrailingStop("EURUSD", "BUY", 1.08000, 1.08500, 30)
fmt.Printf("Trailing SL: %.5f (30 pips from current)\n", newSL)

// Calculate break-even
breakEvenSL, _ := helper.CalculateBreakEven("EURUSD", "BUY", 1.08000, 10)
fmt.Printf("Break-even SL: %.5f (+10 pip buffer)\n", breakEvenSL)

// Show detailed calculation
helper.ShowCalculation("EURUSD", "BUY", 50, 100)
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `GetSymbolInfo()` - Get point size for calculations
* `GetAsk()` - Get current ask price (for BUY)
* `GetBid()` - Get current bid price (for SELL)

**🍬 Complementary sugar methods:**

* `BuyMarketWithPips()` - Uses this method internally ⭐
* `SellMarketWithPips()` - Uses this method internally ⭐
* `BuyMarketWithSLTP()` - Opens BUY with exact prices
* `SellMarketWithSLTP()` - Opens SELL with exact prices
* `ModifyPosition()` - Modify SL/TP after opening

**Recommended usage:**
```go
// Calculate SL/TP
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, 50, 100)

// Use in trade
ticket, _ := sugar.BuyMarketWithSLTP("EURUSD", 0.1, sl, tp)

// OR simpler - use BuyMarketWithPips (does both steps)
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
```

---

## ⚠️ Common Pitfalls

### 1) Wrong direction

```go
// ❌ WRONG - using wrong direction
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, 50, 100)
sugar.SellMarket("EURUSD", 0.1) // Selling with BUY SL/TP!

// ✅ CORRECT - match direction
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "SELL", 0, 50, 100)
sugar.SellMarketWithSLTP("EURUSD", 0.1, sl, tp)
```

### 2) Confusing pips and price

```go
// ❌ WRONG - passing price as pips
stopLossPrice := 1.08000
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, stopLossPrice, 100)

// ✅ CORRECT - use pip distance
stopLossPips := 50.0
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, stopLossPips, 100)
```

### 3) Not using result

```go
// ❌ WRONG - calculating but not using
sugar.CalculateSLTP("EURUSD", "BUY", 0, 50, 100)
sugar.BuyMarket("EURUSD", 0.1) // No SL/TP!

// ✅ CORRECT - use the calculated values
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, 50, 100)
sugar.BuyMarketWithSLTP("EURUSD", 0.1, sl, tp)
```

### 4) Using 0 entry for pending orders

```go
// ❌ WRONG - using 0 entry for limit order
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", 0, 50, 100)
// Will use current price, not your limit price!

// ✅ CORRECT - specify exact entry price
limitPrice := 1.08000
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", limitPrice, 50, 100)
```

### 5) Ignoring errors

```go
// ❌ WRONG - ignoring errors
sl, tp, _ := sugar.CalculateSLTP("INVALID", "BUY", 0, 50, 100)
// sl and tp will be 0!

// ✅ CORRECT - check errors
sl, tp, err := sugar.CalculateSLTP("EURUSD", "BUY", 0, 50, 100)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
```

---

## 💎 Pro Tips

1. **Use 0 for entry** - Automatically uses current market price

2. **Direction matters** - BUY: SL below, TP above / SELL: opposite

3. **Points not pips** - Parameter is in points (for 5-digit: 50 points = 5 pips)

4. **R:R ratios** - TP pips / SL pips = Risk/Reward ratio

5. **Simpler alternative** - Use `BuyMarketWithPips()` / `SellMarketWithPips()` instead

6. **Planning tool** - Great for pre-trade calculations and analysis

7. **Works for any symbol** - Handles different point sizes automatically

---

## 📊 BUY vs SELL

```
BUY positions:
Entry:  1.08500
SL:     1.08000  (below entry ⬇️)
TP:     1.09000  (above entry ⬆️)

SELL positions:
Entry:  1.08500
SL:     1.09000  (above entry ⬆️)
TP:     1.08000  (below entry ⬇️)

Remember:
- SL protects against adverse moves
- TP captures profit in favorable direction
```

---

**See also:** [`BuyMarketWithPips.md`](BuyMarketWithPips.md), [`SellMarketWithPips.md`](SellMarketWithPips.md)
