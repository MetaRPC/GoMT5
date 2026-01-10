# 📉 Sell Market with Pips (`SellMarketWithPips`)

> **Sugar method:** Opens SELL position with SL/TP specified in **pips** - the most intuitive way to short!

**API Information:**

* **Method:** `sugar.SellMarketWithPips(symbol, volume, stopLossPips, takeProfitPips)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `CalculateSLTP()`, `SellMarketWithSLTP()`
* **Timeout:** 10 seconds
* **Returns:** Position ticket number

---

## 📋 Method Signature

```go
func (s *MT5Sugar) SellMarketWithPips(
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
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1, 1.0) |
| `stopLossPips` | `float64` | Stop Loss distance in **points** from entry |
| `takeProfitPips` | `float64` | Take Profit distance in **points** from entry |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Position ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** Opens SELL at market price with SL/TP specified in pips.
* **Why you need it:** **Most intuitive way to short** - think in pips, not prices!
* **Sanity check:** "SELL EURUSD 0.1 lots, 50 pip SL, 100 pip TP" → one method call.

---

## 🎯 When to Use

✅ **Market SELL orders with SL/TP** - Most common shorting scenario

✅ **Think in pips** - Natural risk/reward planning

✅ **Quick short entries** - One-liner to open position

✅ **Risk management** - Specify exact pip-based risk

---

## 🔢 How It Works

```
Step 1: Get current BID price (entry for SELL)
Step 2: Calculate SL = entry + (stopLossPips × point)
Step 3: Calculate TP = entry - (takeProfitPips × point)
Step 4: Open SELL position with calculated SL/TP

Example (EURUSD):
- Current BID: 1.08500
- SL pips: 50
- TP pips: 100

SL = 1.08500 + (50 × 0.00001) = 1.09000  (ABOVE entry for SELL)
TP = 1.08500 - (100 × 0.00001) = 1.07500 (BELOW entry for SELL)

Opens SELL @ 1.08500, SL @ 1.09000, TP @ 1.07500
```

---

## 🔗 Usage Examples

### 1) Basic usage - SELL with SL/TP

```go
symbol := "EURUSD"
volume := 0.1
stopLoss := 50.0   // 50 pips
takeProfit := 100.0 // 100 pips (1:2 R:R)

ticket, err := sugar.SellMarketWithPips(symbol, volume, stopLoss, takeProfit)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ SELL position opened: #%d\n", ticket)
fmt.Printf("   Symbol: %s\n", symbol)
fmt.Printf("   Volume: %.2f lots\n", volume)
fmt.Printf("   SL: %.0f pips, TP: %.0f pips\n", stopLoss, takeProfit)
```

---

### 2) Complete trading workflow with risk management

```go
symbol := "GBPUSD"
riskPercent := 2.0
stopLossPips := 60.0
takeProfitPips := 120.0

// Step 1: Calculate position size based on risk
lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)
if err != nil {
    fmt.Printf("Position size calculation failed: %v\n", err)
    return
}

// Step 2: Validate
canOpen, reason, err := sugar.CanOpenPosition(symbol, lotSize)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
    return
}

if !canOpen {
    fmt.Printf("❌ Cannot open position: %s\n", reason)
    return
}

// Step 3: Open SELL position
fmt.Printf("Opening SELL position:\n")
fmt.Printf("  Symbol:      %s\n", symbol)
fmt.Printf("  Lot size:    %.2f (risking %.1f%%)\n", lotSize, riskPercent)
fmt.Printf("  Stop Loss:   %.0f pips\n", stopLossPips)
fmt.Printf("  Take Profit: %.0f pips\n", takeProfitPips)

ticket, err := sugar.SellMarketWithPips(symbol, lotSize, stopLossPips, takeProfitPips)
if err != nil {
    fmt.Printf("❌ Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Position opened: #%d\n", ticket)

// Step 4: Get entry details
pos, _ := sugar.GetPositionByTicket(ticket)
fmt.Printf("\nPosition details:\n")
fmt.Printf("  Entry:  %.5f\n", pos.OpenPrice)
fmt.Printf("  SL:     %.5f (above entry)\n", pos.StopLoss)
fmt.Printf("  TP:     %.5f (below entry)\n", pos.TakeProfit)
```

---

### 3) Show SL/TP placement for SELL

```go
symbol := "EURUSD"
volume := 0.1
slPips := 50.0
tpPips := 100.0

// Get current price
bid, _ := sugar.GetBid(symbol)

// Calculate what SL/TP prices will be
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", bid, slPips, tpPips)

fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("  SELL ORDER PREVIEW\n")
fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("Symbol:       %s\n", symbol)
fmt.Printf("Entry (BID):  %.5f\n", bid)
fmt.Printf("Stop Loss:    %.5f (+%.0f pips ABOVE)\n", sl, slPips)
fmt.Printf("Take Profit:  %.5f (-%.0f pips BELOW)\n", tp, tpPips)
fmt.Printf("Risk/Reward:  1:%.1f\n", tpPips/slPips)
fmt.Printf("═══════════════════════════════════════\n")

// Open position
ticket, _ := sugar.SellMarketWithPips(symbol, volume, slPips, tpPips)
fmt.Printf("✅ SELL order #%d opened\n", ticket)

// Output:
// ═══════════════════════════════════════
//   SELL ORDER PREVIEW
// ═══════════════════════════════════════
// Symbol:       EURUSD
// Entry (BID):  1.08500
// Stop Loss:    1.09000 (+50 pips ABOVE)
// Take Profit:  1.07500 (-100 pips BELOW)
// Risk/Reward:  1:2.0
// ═══════════════════════════════════════
// ✅ SELL order #12345 opened
```

---

### 4) Different R:R ratios for SELL

```go
symbol := "USDJPY"
volume := 0.15
stopLoss := 40.0

// Conservative (1:1.5)
ticket1, _ := sugar.SellMarketWithPips(symbol, volume, stopLoss, stopLoss*1.5)
fmt.Printf("Conservative: #%d (1:1.5 R:R, %.0f/%.0f pips)\n",
    ticket1, stopLoss, stopLoss*1.5)

// Moderate (1:2)
ticket2, _ := sugar.SellMarketWithPips(symbol, volume, stopLoss, stopLoss*2.0)
fmt.Printf("Moderate:     #%d (1:2 R:R, %.0f/%.0f pips)\n",
    ticket2, stopLoss, stopLoss*2.0)

// Aggressive (1:3)
ticket3, _ := sugar.SellMarketWithPips(symbol, volume, stopLoss, stopLoss*3.0)
fmt.Printf("Aggressive:   #%d (1:3 R:R, %.0f/%.0f pips)\n",
    ticket3, stopLoss, stopLoss*3.0)

// Output:
// Conservative: #123456 (1:1.5 R:R, 40/60 pips)
// Moderate:     #123457 (1:2 R:R, 40/80 pips)
// Aggressive:   #123458 (1:3 R:R, 40/120 pips)
```

---

### 5) Multi-symbol short portfolio

```go
type ShortTrade struct {
    Symbol     string
    Volume     float64
    SLPips     float64
    TPPips     float64
}

shorts := []ShortTrade{
    {"EURUSD", 0.1, 50, 100},
    {"GBPUSD", 0.15, 60, 120},
    {"AUDUSD", 0.2, 45, 90},
}

fmt.Println("Opening SELL positions on multiple pairs:")
fmt.Println("─────────────────────────────────────────")

successCount := 0

for i, trade := range shorts {
    ticket, err := sugar.SellMarketWithPips(
        trade.Symbol,
        trade.Volume,
        trade.SLPips,
        trade.TPPips,
    )

    if err != nil {
        fmt.Printf("%d. %s: ❌ Failed - %v\n", i+1, trade.Symbol, err)
        continue
    }

    successCount++
    fmt.Printf("%d. %s: ✅ #%d (%.2f lots, %.0f/%.0f pips)\n",
        i+1, trade.Symbol, ticket, trade.Volume, trade.SLPips, trade.TPPips)
}

fmt.Println("─────────────────────────────────────────")
fmt.Printf("Opened: %d/%d positions\n", successCount, len(shorts))
```

---

### 6) Trend-following short strategy

```go
func ShortOnDowntrend(
    sugar *mt5.MT5Sugar,
    symbol string,
    riskPercent float64,
) (uint64, error) {
    // Example: SELL when price breaks support
    // (In real trading, you'd have trend detection logic)

    stopLossPips := 50.0
    takeProfitPips := 150.0 // 1:3 R:R for trend trades

    // Calculate position size
    lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)
    if err != nil {
        return 0, err
    }

    // Validate
    canOpen, reason, _ := sugar.CanOpenPosition(symbol, lotSize)
    if !canOpen {
        return 0, fmt.Errorf("cannot open: %s", reason)
    }

    fmt.Printf("TREND-FOLLOWING SHORT:\n")
    fmt.Printf("  Symbol:   %s\n", symbol)
    fmt.Printf("  Strategy: Downtrend continuation\n")
    fmt.Printf("  SL:       %.0f pips\n", stopLossPips)
    fmt.Printf("  TP:       %.0f pips (1:3 R:R)\n", takeProfitPips)
    fmt.Printf("  Volume:   %.2f lots\n\n", lotSize)

    // Open SELL
    ticket, err := sugar.SellMarketWithPips(symbol, lotSize, stopLossPips, takeProfitPips)
    if err != nil {
        return 0, err
    }

    fmt.Printf("✅ SELL position #%d opened\n", ticket)

    return ticket, nil
}

// Usage:
ticket, err := ShortOnDowntrend(sugar, "EURUSD", 2.0)
if err != nil {
    fmt.Printf("Failed: %v\n", err)
}
```

---

### 7) Counter-trend short (tighter stops)

```go
func ShortCounterTrend(
    sugar *mt5.MT5Sugar,
    symbol string,
    volume float64,
) (uint64, error) {
    // Counter-trend = tighter stops, closer targets
    stopLossPips := 30.0  // Tight stop
    takeProfitPips := 45.0 // 1:1.5 R:R (counter-trend is riskier)

    fmt.Printf("COUNTER-TREND SHORT:\n")
    fmt.Printf("  Symbol:   %s\n", symbol)
    fmt.Printf("  Strategy: Mean reversion\n")
    fmt.Printf("  SL:       %.0f pips (tight!)\n", stopLossPips)
    fmt.Printf("  TP:       %.0f pips\n", takeProfitPips)
    fmt.Printf("  Volume:   %.2f lots\n\n", volume)

    ticket, err := sugar.SellMarketWithPips(symbol, volume, stopLossPips, takeProfitPips)
    if err != nil {
        return 0, err
    }

    fmt.Printf("✅ SELL position #%d opened\n", ticket)

    return ticket, nil
}

// Usage:
ticket, err := ShortCounterTrend(sugar, "EURUSD", 0.1)
if err != nil {
    fmt.Printf("Failed: %v\n", err)
}
```

---

### 8) Batch SELL with validation

```go
func OpenMultipleSellPositions(
    sugar *mt5.MT5Sugar,
    trades []ShortTrade,
) []uint64 {
    tickets := []uint64{}

    fmt.Println("Opening SELL positions with validation:")
    fmt.Println("─────────────────────────────────────────")

    for i, trade := range trades {
        // Step 1: Validate
        canOpen, reason, err := sugar.CanOpenPosition(trade.Symbol, trade.Volume)
        if err != nil {
            fmt.Printf("%d. %s: ⚠️  Validation error - %v\n",
                i+1, trade.Symbol, err)
            continue
        }

        if !canOpen {
            fmt.Printf("%d. %s: ❌ Cannot open - %s\n",
                i+1, trade.Symbol, reason)
            continue
        }

        // Step 2: Open SELL position
        ticket, err := sugar.SellMarketWithPips(
            trade.Symbol,
            trade.Volume,
            trade.SLPips,
            trade.TPPips,
        )

        if err != nil {
            fmt.Printf("%d. %s: ❌ Order failed - %v\n",
                i+1, trade.Symbol, err)
            continue
        }

        tickets = append(tickets, ticket)
        fmt.Printf("%d. %s: ✅ Opened #%d (%.2f lots)\n",
            i+1, trade.Symbol, ticket, trade.Volume)
    }

    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Successfully opened: %d/%d positions\n",
        len(tickets), len(trades))

    return tickets
}

// Usage:
trades := []ShortTrade{
    {"EURUSD", 0.1, 50, 100},
    {"GBPUSD", 0.15, 60, 120},
    {"USDJPY", 0.1, 40, 80},
}

tickets := OpenMultipleSellPositions(sugar, trades)
fmt.Printf("\nOpened SELL positions: %v\n", tickets)
```

---

### 9) SELL with trade logging

```go
type TradeLogger struct {
    sugar *mt5.MT5Sugar
}

func (tl *TradeLogger) SellWithLog(
    symbol string,
    volume, slPips, tpPips float64,
    reason string,
) (uint64, error) {
    // Get current price
    bid, _ := tl.sugar.GetBid(symbol)
    balance, _ := tl.sugar.GetBalance()

    // Calculate SL/TP prices
    sl, tp, _ := tl.sugar.CalculateSLTP(symbol, "SELL", 0, slPips, tpPips)

    // Log entry
    timestamp := time.Now().Format("2006-01-02 15:04:05")

    fmt.Println("═══════════════════════════════════════════════════════")
    fmt.Printf("TRADE LOG - %s\n", timestamp)
    fmt.Println("═══════════════════════════════════════════════════════")
    fmt.Printf("Action:         SELL\n")
    fmt.Printf("Symbol:         %s\n", symbol)
    fmt.Printf("Volume:         %.2f lots\n", volume)
    fmt.Printf("Entry (BID):    %.5f\n", bid)
    fmt.Printf("Stop Loss:      %.5f (+%.0f pips above)\n", sl, slPips)
    fmt.Printf("Take Profit:    %.5f (-%.0f pips below)\n", tp, tpPips)
    fmt.Printf("R:R Ratio:      1:%.1f\n", tpPips/slPips)
    fmt.Printf("Balance:        $%.2f\n", balance)
    fmt.Printf("Reason:         %s\n", reason)
    fmt.Println("───────────────────────────────────────────────────────")

    // Open position
    ticket, err := tl.sugar.SellMarketWithPips(symbol, volume, slPips, tpPips)

    if err != nil {
        fmt.Printf("Result:         ❌ FAILED\n")
        fmt.Printf("Error:          %v\n", err)
        fmt.Println("═══════════════════════════════════════════════════════")
        return 0, err
    }

    fmt.Printf("Result:         ✅ SUCCESS\n")
    fmt.Printf("Ticket:         #%d\n", ticket)
    fmt.Println("═══════════════════════════════════════════════════════")

    return ticket, nil
}

// Usage:
logger := &TradeLogger{sugar: sugar}

ticket, err := logger.SellWithLog(
    "EURUSD",
    0.1,
    50,
    100,
    "Bearish engulfing on H4, resistance at 1.0900",
)

if err != nil {
    fmt.Printf("Trade failed: %v\n", err)
}
```

---

### 10) Advanced shorting system

```go
type ShortingSystem struct {
    sugar             *mt5.MT5Sugar
    maxShortPositions int
    maxRiskPercent    float64
}

func NewShortingSystem(
    sugar *mt5.MT5Sugar,
    maxPositions int,
    maxRisk float64,
) *ShortingSystem {
    return &ShortingSystem{
        sugar:             sugar,
        maxShortPositions: maxPositions,
        maxRiskPercent:    maxRisk,
    }
}

func (ss *ShortingSystem) CanShort() (bool, string) {
    // Count existing SELL positions
    positions, _ := ss.sugar.GetOpenPositions()
    shortCount := 0

    for _, pos := range positions {
        if pos.Type == 1 { // SELL
            shortCount++
        }
    }

    if shortCount >= ss.maxShortPositions {
        return false, fmt.Sprintf(
            "max short positions reached (%d/%d)",
            shortCount, ss.maxShortPositions,
        )
    }

    return true, ""
}

func (ss *ShortingSystem) OpenShort(
    symbol string,
    slPips float64,
    tpPips float64,
    reason string,
) (uint64, error) {
    // Check if we can short
    canShort, reason := ss.CanShort()
    if !canShort {
        return 0, fmt.Errorf("cannot short: %s", reason)
    }

    // Calculate position size
    lotSize, err := ss.sugar.CalculatePositionSize(
        symbol,
        ss.maxRiskPercent,
        slPips,
    )
    if err != nil {
        return 0, err
    }

    // Validate
    canOpen, validationReason, err := ss.sugar.CanOpenPosition(symbol, lotSize)
    if err != nil {
        return 0, err
    }

    if !canOpen {
        return 0, fmt.Errorf("validation failed: %s", validationReason)
    }

    // Open SELL position
    fmt.Printf("Opening SHORT position:\n")
    fmt.Printf("  Symbol:   %s\n", symbol)
    fmt.Printf("  Volume:   %.2f lots (%.1f%% risk)\n",
        lotSize, ss.maxRiskPercent)
    fmt.Printf("  SL:       %.0f pips\n", slPips)
    fmt.Printf("  TP:       %.0f pips (1:%.1f R:R)\n", tpPips, tpPips/slPips)
    fmt.Printf("  Reason:   %s\n", reason)

    ticket, err := ss.sugar.SellMarketWithPips(symbol, lotSize, slPips, tpPips)
    if err != nil {
        return 0, err
    }

    fmt.Printf("✅ SHORT position opened: #%d\n\n", ticket)

    return ticket, nil
}

func (ss *ShortingSystem) ShowStatus() {
    positions, _ := ss.sugar.GetOpenPositions()
    balance, _ := ss.sugar.GetBalance()

    shortCount := 0
    for _, pos := range positions {
        if pos.Type == 1 { // SELL
            shortCount++
        }
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      SHORTING SYSTEM STATUS           ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Balance:         $%.2f\n", balance)
    fmt.Printf("Short positions: %d/%d\n", shortCount, ss.maxShortPositions)
    fmt.Printf("Max risk:        %.1f%% per trade\n", ss.maxRiskPercent)

    canShort, reason := ss.CanShort()
    if canShort {
        fmt.Println("\n✅ Ready to short")
    } else {
        fmt.Printf("\n❌ Cannot short: %s\n", reason)
    }
}

// Usage:
shortSystem := NewShortingSystem(
    sugar,
    3,   // Max 3 short positions
    2.0, // Max 2% risk per trade
)

shortSystem.ShowStatus()

// Open short with automatic risk management
ticket, err := shortSystem.OpenShort(
    "EURUSD",
    50,
    100,
    "Breakdown below support, strong selling pressure",
)

if err != nil {
    fmt.Printf("Short failed: %v\n", err)
} else {
    fmt.Printf("Short successful: #%d\n", ticket)
}
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `CalculateSLTP()` - Calculates SL/TP prices from pips
* `SellMarketWithSLTP()` - Opens position with exact prices

**🍬 Complementary sugar methods:**

* `BuyMarketWithPips()` - BUY version of this method ⭐
* `SellMarket()` - SELL without SL/TP
* `SellMarketWithSLTP()` - SELL with exact prices
* `CalculatePositionSize()` - Calculate risk-based lot size ⭐
* `CanOpenPosition()` - Validate before opening ⭐

**Recommended workflow:**
```go
// 1. Calculate lot size
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// 2. Validate
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    fmt.Println("Cannot open:", reason)
    return
}

// 3. Trade with pips (most intuitive!)
ticket, _ := sugar.SellMarketWithPips("EURUSD", lotSize, 50, 100)
```

---

## ⚠️ Common Pitfalls

### 1) Confusing SELL vs BUY SL placement

```go
// ❌ WRONG - thinking SL is below entry for SELL
// For SELL: SL is ABOVE entry, TP is BELOW entry

// ✅ CORRECT - understand SELL mechanics
// SELL @ 1.08500
// SL @ 1.09000 (50 pips ABOVE - protects if price goes UP)
// TP @ 1.07500 (100 pips BELOW - profit when price goes DOWN)
```

### 2) Confusing pips and price

```go
// ❌ WRONG - passing price as pips
stopLossPrice := 1.09000
sugar.SellMarketWithPips("EURUSD", 0.1, stopLossPrice, 100)

// ✅ CORRECT - use pip distance
stopLossPips := 50.0
sugar.SellMarketWithPips("EURUSD", 0.1, stopLossPips, 100)
```

### 3) Not validating before opening

```go
// ❌ WRONG - shorting without validation
sugar.SellMarketWithPips("EURUSD", 10.0, 50, 100) // Might fail!

// ✅ CORRECT - validate first
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 10.0)
if !canOpen {
    fmt.Println("Cannot open:", reason)
    return
}
sugar.SellMarketWithPips("EURUSD", 10.0, 50, 100)
```

### 4) Using fixed lot size

```go
// ❌ WRONG - fixed lot size (ignores account growth)
sugar.SellMarketWithPips("EURUSD", 0.1, 50, 100)

// ✅ CORRECT - dynamic lot size based on risk
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
sugar.SellMarketWithPips("EURUSD", lotSize, 50, 100)
```

### 5) Not checking errors

```go
// ❌ WRONG - ignoring errors
sugar.SellMarketWithPips("INVALID", 0.1, 50, 100)

// ✅ CORRECT - check errors
ticket, err := sugar.SellMarketWithPips("EURUSD", 0.1, 50, 100)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}
```

---

## 💎 Pro Tips

1. **SL ABOVE entry** - For SELL, stop loss is ABOVE (protects if price rises)

2. **TP BELOW entry** - Take profit is BELOW (profit when price falls)

3. **Think in pips** - This is THE method for shorting

4. **Always validate** - Use `CanOpenPosition()` first

5. **Use CalculatePositionSize** - For proper risk management

6. **R:R ratios** - Use at least 1:1.5, ideally 1:2 or better

7. **Points not pips** - Parameter is points (for 5-digit: 50 points = 5 pips)

---

## 📊 SELL: SL/TP Placement

```
SELL Position:
Entry:  1.08500 (current BID)
SL:     1.09000 (ABOVE entry - 50 pips)
TP:     1.07500 (BELOW entry - 100 pips)

Why SL is above:
- You sold (shorted) at 1.08500
- If price goes UP, you lose
- SL at 1.09000 protects you
- If price hits 1.09000 → close at loss

Why TP is below:
- If price goes DOWN, you profit
- TP at 1.07500 locks in profit
- If price hits 1.07500 → close at profit
```

---

**See also:** [`BuyMarketWithPips.md`](BuyMarketWithPips.md), [`CalculatePositionSize.md`](../10.%20Risk_Management/CalculatePositionSize.md), [`CalculateSLTP.md`](CalculateSLTP.md)
