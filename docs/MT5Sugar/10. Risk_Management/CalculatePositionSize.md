# ⚖️ Calculate Position Size Based on Risk (`CalculatePositionSize`)

> **Sugar method:** THE MOST IMPORTANT risk management tool - calculates a position size from the % of balance you want to risk and the stop-loss distance, then caps it by your free margin.

**API Information:**

* **Method:** `sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `GetBalance()`, `GetSymbolInfo()`, `service.GetSymbolDouble()` (tick value, tick size), `service.GetSymbolTick()`, `GetFreeMargin()`, `service.CalculateMargin()`
* **Timeout:** 10 seconds (for the tick value/size, tick and margin queries; `GetBalance`, `GetSymbolInfo` and `GetFreeMargin` use their own 3 s / 10 s / 3 s timeouts)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) CalculatePositionSize(
    symbol string,
    riskPercent float64,
    stopLossPips float64,
) (float64, error)
```

---

## 🔽 Input

| Parameter | Type | Description |
|-----------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |
| `riskPercent` | `float64` | Percentage of balance to risk (e.g., 2.0 = 2%) |
| `stopLossPips` | `float64` | Stop Loss distance in **POINTS** despite the name (1 point = `SYMBOL_POINT`: 0.00001 on 5-digit EURUSD, so 10 points = 1 pip; 0.001 on 3-digit USDJPY). Not a price! Must be > 0 |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `lotSize` | `float64` | Recommended lot size for `riskPercent` (capped by free margin) |
| `error` | `error` | `stop loss distance must be > 0 points, got X`, an error if a query fails or the symbol is not found, `symbol X has no tick value yet (...)`, `risk X% with a N-point stop allows Y lots, below the minimum volume Z of SYMBOL`, or `insufficient margin: recommended size X is below minimum Y` |

**Returned lot size:**

- Rounded DOWN to the symbol's `VolumeStep` (with a small float tolerance: 0.3 / 0.01 gives 0.30, not 0.29)
- Capped at `VolumeMax`. A risk-based size below `VolumeMin` is an **error** (`risk X% with a N-point stop allows Y lots, below the minimum volume ...`) - it is never raised to `VolumeMin`, which would risk more than `riskPercent`
- Then capped at 80% of free margin ÷ margin for 1 lot (margin of a BUY at the current ASK - the method has no direction; also rounded down to `VolumeStep`)
- If that cap pushes it below `VolumeMin`, an `insufficient margin` error is returned instead
- Ready to use in `BuyMarket()` / `SellMarket()` calls

---

## 💬 Just the Essentials

* **What it is:** Auto-calculates lot size based on how much you want to risk (% of balance).
* **Why you need it:** **NEVER use fixed lot sizes!** This keeps the risk at the same % on every trade regardless of balance.
* **Sanity check:** Balance = $10,000 (USD account), risk = 2%, SL = 500 points (50 pips on a 5-digit EURUSD; 1 point of 1 lot = 1 USD) → risk amount $200 → 0.40 lots (unless the margin cap is lower). Verified on a demo account: balance 100,008.77 USD, 2%, 500 points EURUSD → 4.00 lots.

---

## 🎯 Purpose

Use it for professional risk management:

* **Risk exactly X% per trade** (recommended: 1-2%)
* **Auto-adjust position size** as your balance changes
* **Prevent overleveraging** by calculating safe lot sizes
* **Follow professional money management** rules
* **Scale positions consistently** across different symbols

---

## 🔢 Formula

```
LotSize = (Balance × RiskPercent / 100) / (StopLossPoints × PointValue)

Where:
  PointValue = SYMBOL_TRADE_TICK_VALUE × SYMBOL_POINT / SYMBOL_TRADE_TICK_SIZE
               (value of a 1-point move for 1 lot, in the ACCOUNT currency)

Then:
  - Round DOWN to VolumeStep
  - Below VolumeMin → error (never raised to VolumeMin); above VolumeMax → VolumeMax
  - Cap at MaxLots = (FreeMargin × 0.8) / MarginFor1Lot  (BUY, 1 lot, at ASK - no direction; rounded down to VolumeStep)
  - Error "insufficient margin" if the result is below VolumeMin
```

**Example:**
```
Balance:     $10,000 (USD account)
Risk:        2% ($200)
Stop Loss:   500 points (= 50 pips on 5-digit EURUSD)
PointValue:  tick value 1.00 USD × 0.00001 / tick size 0.00001 = $1 per point per lot (EURUSD)

LotSize = $200 / (500 × $1) = $200 / $500 = 0.40 lots
```

USDJPY on the same USD account at 150.00: point 0.001, tick size 0.001, tick value = 100 JPY = 0.6667 USD per lot
→ PointValue = $0.6667 → $200 / (500 × $0.6667) = 0.60 lots. The tick value is already converted to the
account currency by the terminal, so the risk is the same $200 on every symbol.

---

## 🧩 Notes & Tips

* **Always use this method** - Don't hardcode lot sizes!
* **Risk percent** - Professional traders use 1-2% max per trade
* **Stop loss in POINTS** - Despite the name `stopLossPips`, the value is a distance in points, not pips and not a price (500 = 500 points = 50 pips on a 5-digit symbol). It must be > 0, otherwise an error is returned
* **Symbol matters** - Different symbols have different tick values and point sizes
* **Auto-rounding** - Result is rounded down to the volume step
* **Safety limits** - Result is capped at the broker's max volume and by 80% of free margin (margin of a BUY at ASK; a SELL's margin may differ slightly)
* **Account currency** - The point value comes from `SYMBOL_TRADE_TICK_VALUE`, which the terminal gives in the account currency, so USDJPY, EURGBP etc. are sized correctly on any account. A symbol without a tick value yet (no quote) gives an error
* **Minimum volume** - If the risk-based size is below `VolumeMin`, the method returns an error instead of trading more risk than you asked for (small account, tight risk or a wide stop) - raise the risk %, tighten the stop, or skip the trade

---

## 🔧 Under the Hood

This method performs several steps (from `MT5Sugar.go`):

```go
func (s *MT5Sugar) CalculatePositionSize(symbol string, riskPercent, stopLossPips float64) (float64, error) {
    if stopLossPips <= 0 {
        return 0, fmt.Errorf("stop loss distance must be > 0 points, got %v", stopLossPips)
    }

    ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
    defer cancel()

    // Get account balance
    balance, err := s.GetBalance()
    if err != nil {
        return 0, fmt.Errorf("failed to get balance: %w", err)
    }

    // Get symbol info
    info, err := s.GetSymbolInfo(symbol)
    if err != nil {
        return 0, err
    }

    // Calculate risk amount
    riskAmount := balance * riskPercent / 100.0

    // Value of 1 point for 1 lot in the ACCOUNT currency (contract size * point would be in the quote
    // currency, wrong for e.g. USDJPY on a USD account).
    tickValue, err := s.service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE)
    // ... error check
    tickSize, err := s.service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_SIZE)
    // ... error check
    if tickValue <= 0 || tickSize <= 0 {
        return 0, fmt.Errorf("symbol %s has no tick value yet (tick value %v, tick size %v)", symbol, tickValue, tickSize)
    }
    pointValue := tickValue * info.Point / tickSize

    // Calculate position size based on risk
    positionSize := floorToStep(riskAmount/(stopLossPips*pointValue), info.VolumeStep)

    // Within min/max limits. A size below the minimum volume is NOT raised to it: that would risk more
    // than riskPercent.
    if positionSize < info.VolumeMin {
        return 0, fmt.Errorf("risk %.2f%% with a %v-point stop allows %.4f lots, below the minimum volume %.2f of %s",
            riskPercent, stopLossPips, riskAmount/(stopLossPips*pointValue), info.VolumeMin, symbol)
    }
    if positionSize > info.VolumeMax {
        positionSize = info.VolumeMax
    }

    // CRITICAL: Check margin limit - prevent margin calls
    freeMargin, err := s.GetFreeMargin()
    // ... error check

    tick, err := s.service.GetSymbolTick(ctx, symbol)
    // ... error check

    // Calculate margin required for 1 lot (a BUY opens at ASK)
    marginReq := &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    1.0,
        OpenPrice: tick.Ask,
    }

    marginForOneLot, err := s.service.CalculateMargin(ctx, marginReq)
    // ... error check

    if marginForOneLot > 0 {
        // Calculate max lots with 80% safety buffer, rounded down to the volume step
        maxLots := floorToStep((freeMargin*0.8)/marginForOneLot, info.VolumeStep)

        // Return minimum of risk-based size and margin-limited size
        if positionSize > maxLots {
            positionSize = maxLots
        }
    }

    // Final check - ensure not below minimum
    if positionSize < info.VolumeMin {
        return 0, fmt.Errorf("insufficient margin: recommended size %.2f is below minimum %.2f", positionSize, info.VolumeMin)
    }

    return positionSize, nil
}

// floorToStep rounds a volume down to a multiple of step (with a little tolerance for float noise:
// 0.3/0.01 = 29.9999999 must give 0.30, not 0.29) and trims the float noise of the result.
func floorToStep(volume, step float64) float64 {
    if step <= 0 {
        return volume
    }
    v := math.Floor(volume/step+1e-9) * step
    return math.Round(v*1e8) / 1e8
}
```

**What it improves:**

* ✅ **Automatic calculation** - No manual math needed
* ✅ **Dynamic risk** - Adjusts to balance changes
* ✅ **Symbol-aware** - Handles different contract sizes
* ✅ **Broker-compliant** - Respects volume limits
* ✅ **Ready to trade** - Result can be used directly

---

## 📊 Low-Level Alternative

**WITHOUT sugar (manual calculation):**
```go
symbol := "EURUSD"

// Get balance
balance, _ := service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE)

// Get symbol params (the name filter may return similar names: take the exact one)
params, _, _ := service.GetSymbolParamsMany(ctx, &symbol, nil, nil, nil)
var p mt5.SymbolParams
for _, sp := range params {
    if sp.Name == symbol {
        p = sp
    }
}

// Get current price
tick, _ := service.GetSymbolTick(ctx, symbol)

// Value of 1 point for 1 lot in the ACCOUNT currency
tickValue, _ := service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE)
tickSize, _ := service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_SIZE)
pointValue := tickValue * p.Point / tickSize

// Calculate risk amount
riskAmount := balance * 2.0 / 100.0

// Calculate lot size (500 points stop loss)
lotSize := riskAmount / (500.0 * pointValue)

// Round down to step (small tolerance for float noise)
lotSize = math.Floor(lotSize/p.VolumeStep+1e-9) * p.VolumeStep

// Limits: below the minimum the risk does not allow a trade (raising it would risk more)
if lotSize < p.VolumeMin {
    // risk too small for the minimum volume - do not trade
}
if lotSize > p.VolumeMax {
    lotSize = p.VolumeMax
}

// Cap by 80% of free margin
freeMargin, _ := service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE)
marginForOneLot, _ := service.CalculateMargin(ctx, &pb.OrderCalcMarginRequest{
    Symbol:    symbol,
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:    1.0,
    OpenPrice: tick.Ask, // a BUY opens at ASK
})
if marginForOneLot > 0 {
    maxLots := math.Floor((freeMargin*0.8)/marginForOneLot/p.VolumeStep+1e-9) * p.VolumeStep
    if lotSize > maxLots {
        lotSize = maxLots
    }
}
if lotSize < p.VolumeMin {
    // insufficient margin
}
```

**WITH sugar:**
```go
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500) // 500 points
```

**Benefits:**

* ✅ **One line vs 20+ lines**
* ✅ **No manual formula**
* ✅ **No rounding logic needed**
* ✅ **Automatic error handling**
* ✅ **Built-in validation**

---

## 🔗 Usage Examples

### 1) Basic usage - risk 2% with a 500-point (50-pip) stop

```go
// Risk 2% of balance with a 500-point stop loss (= 50 pips on 5-digit EURUSD)
lotSize, err := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Recommended lot size: %.2f\n", lotSize)
// Output ($10,000 balance): Recommended lot size: 0.40
```

---

### 2) Complete trade setup with risk management

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLossPips := 500.0   // points
takeProfitPips := 1000.0 // points

// Step 1: Calculate position size
lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)
if err != nil {
    fmt.Printf("Position size calculation failed: %v\n", err)
    return
}

// Step 2: Validate we can open this position
canOpen, reason, err := sugar.CanOpenPosition(symbol, lotSize)
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
    return
}
if !canOpen {
    fmt.Printf("Cannot open position: %s\n", reason)
    return
}

// Step 3: Open position
ticket, err := sugar.BuyMarketWithPips(symbol, lotSize, stopLossPips, takeProfitPips)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Order #%d opened with lot size %.2f\n", ticket, lotSize)
fmt.Printf("   Risk: %.1f%% of balance (%.0f points SL)\n", riskPercent, stopLossPips)
```

---

### 3) Conservative vs aggressive risk levels

```go
symbol := "EURUSD"
stopLoss := 500.0 // points

// Conservative (1% risk)
conservativeLots, _ := sugar.CalculatePositionSize(symbol, 1.0, stopLoss)
fmt.Printf("Conservative (1%% risk): %.2f lots\n", conservativeLots)

// Moderate (2% risk)
moderateLots, _ := sugar.CalculatePositionSize(symbol, 2.0, stopLoss)
fmt.Printf("Moderate (2%% risk):     %.2f lots\n", moderateLots)

// Aggressive (5% risk) - NOT RECOMMENDED!
aggressiveLots, _ := sugar.CalculatePositionSize(symbol, 5.0, stopLoss)
fmt.Printf("Aggressive (5%% risk):   %.2f lots (⚠️ HIGH RISK!)\n", aggressiveLots)

// Output ($10,000 balance):
// Conservative (1% risk): 0.20 lots
// Moderate (2% risk):     0.40 lots
// Aggressive (5% risk):   1.00 lots (⚠️ HIGH RISK!)
```

---

### 4) Different stop loss distances

```go
symbol := "EURUSD"
riskPercent := 2.0

// Tight stop (250 points) → larger position size
tightSL, _ := sugar.CalculatePositionSize(symbol, riskPercent, 250)
fmt.Printf("250 point SL: %.2f lots\n", tightSL)

// Normal stop (500 points)
normalSL, _ := sugar.CalculatePositionSize(symbol, riskPercent, 500)
fmt.Printf("500 point SL: %.2f lots\n", normalSL)

// Wide stop (1000 points) → smaller position size
wideSL, _ := sugar.CalculatePositionSize(symbol, riskPercent, 1000)
fmt.Printf("1000 point SL: %.2f lots\n", wideSL)

// Output (for $10,000 balance):
// 250 point SL: 0.80 lots
// 500 point SL: 0.40 lots
// 1000 point SL: 0.20 lots
//
// Note: All target the same $200 risk (2% of $10,000)
```

---

### 5) Multi-symbol portfolio with consistent risk

```go
riskPercent := 2.0
stopLoss := 500.0 // points

symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"}

fmt.Println("Consistent 2% risk across all symbols:")
fmt.Println("─────────────────────────────────────")

for _, symbol := range symbols {
    lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
    if err != nil {
        fmt.Printf("%s: Error - %v\n", symbol, err)
        continue
    }

    fmt.Printf("%s: %.2f lots\n", symbol, lotSize)
}

// Output (example: $10,000 USD account, typical contract specs, USDJPY at 150.00):
// Consistent 2% risk across all symbols:
// ─────────────────────────────────────
// EURUSD: 0.40 lots   (1 point of 1 lot = $1)
// GBPUSD: 0.40 lots   (1 point of 1 lot = $1)
// USDJPY: 0.60 lots   (1 point of 1 lot = 100 JPY = $0.6667)
// XAUUSD: 0.40 lots   (contract 100 oz, point 0.01: 1 point of 1 lot = $1)
//
// Note: the point value comes from SYMBOL_TRADE_TICK_VALUE, which is in the
// account currency, so every symbol risks the same $200. The same 500 points
// is a very different PRICE distance per symbol (0.00500 on EURUSD, 0.500 on
// USDJPY, 5.00 on XAUUSD) - choose the stop per symbol in real trading.
```

---

### 6) Show risk amount in money

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLoss := 500.0 // points

// Get balance
balance, _ := sugar.GetBalance()

// Calculate position size
lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Calculate risk amount
riskAmount := balance * riskPercent / 100.0

fmt.Printf("═══════════════════════════════════\n")
fmt.Printf("  RISK CALCULATION\n")
fmt.Printf("═══════════════════════════════════\n")
fmt.Printf("Symbol:        %s\n", symbol)
fmt.Printf("Balance:       $%.2f\n", balance)
fmt.Printf("Risk:          %.1f%% ($%.2f)\n", riskPercent, riskAmount)
fmt.Printf("Stop Loss:     %.0f points\n", stopLoss)
fmt.Printf("═══════════════════════════════════\n")
fmt.Printf("Lot Size:      %.2f\n", lotSize)
fmt.Printf("═══════════════════════════════════\n")

// Output:
// ═══════════════════════════════════
//   RISK CALCULATION
// ═══════════════════════════════════
// Symbol:        EURUSD
// Balance:       $10000.00
// Risk:          2.0% ($200.00)
// Stop Loss:     500 points
// ═══════════════════════════════════
// Lot Size:      0.40
// ═══════════════════════════════════
```

---

### 7) Account growth simulation

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLoss := 500.0 // points

balances := []float64{1000, 5000, 10000, 50000, 100000}

fmt.Println("Position sizing as account grows:")
fmt.Println("──────────────────────────────────────────────")
fmt.Printf("%-12s  %-12s  %-12s\n", "Balance", "Risk (2%)", "Lot Size")
fmt.Println("──────────────────────────────────────────────")

for _, balance := range balances {
    // Create temporary calculation (in real code you'd have different accounts)
    riskAmount := balance * riskPercent / 100.0

    // For demonstration, calculate approximate lot size
    // In real code, use sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
    approxLots := riskAmount / (stopLoss * 1.0) // EURUSD: 1 point × 1 lot = $1 (ignores margin cap)

    fmt.Printf("$%-11.2f  $%-11.2f  %.2f lots\n",
        balance, riskAmount, approxLots)
}

// Output:
// Position sizing as account grows:
// ──────────────────────────────────────────────
// Balance       Risk (2%)     Lot Size
// ──────────────────────────────────────────────
// $1000.00      $20.00        0.04 lots
// $5000.00      $100.00       0.20 lots
// $10000.00     $200.00       0.40 lots
// $50000.00     $1000.00      2.00 lots
// $100000.00    $2000.00      4.00 lots
```

---

### 8) Validate calculated size before trading

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLoss := 500.0    // points
takeProfit := 1000.0 // points

// Calculate position size
lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
if err != nil {
    fmt.Printf("❌ Calculation failed: %v\n", err)
    return
}

// Validate before trading
canOpen, reason, err := sugar.CanOpenPosition(symbol, lotSize)
if err != nil {
    fmt.Printf("❌ Validation error: %v\n", err)
    return
}

if !canOpen {
    fmt.Printf("❌ Cannot open position: %s\n", reason)
    fmt.Printf("   Calculated lot size: %.2f\n", lotSize)
    return
}

// All checks passed - open position
fmt.Printf("✅ All checks passed\n")
fmt.Printf("   Symbol:     %s\n", symbol)
fmt.Printf("   Lot size:   %.2f\n", lotSize)
fmt.Printf("   Stop Loss:  %.0f points\n", stopLoss)
fmt.Printf("   Take Profit: %.0f points\n", takeProfit)

ticket, _ := sugar.BuyMarketWithPips(symbol, lotSize, stopLoss, takeProfit)
fmt.Printf("   Ticket:     #%d\n", ticket)
```

---

### 9) Handle minimum volume edge case

```go
symbol := "EURUSD"
riskPercent := 0.5  // Very small risk
stopLoss := 1000.0  // Wide stop (points)

lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
if err != nil {
    // e.g. "risk 0.50% with a 1000-point stop allows 0.0040 lots, below the minimum volume 0.01 of EURUSD"
    // (on a $800 account) - the method never raises the size to the minimum, that would risk more
    fmt.Printf("⚠️  %v\n", err)
    fmt.Printf("   Consider:\n")
    fmt.Printf("   - Increasing risk %% (currently %.1f%%)\n", riskPercent)
    fmt.Printf("   - Reducing stop loss (currently %.0f points)\n", stopLoss)
    fmt.Printf("   - Trading on a larger account\n")
    return
}

fmt.Printf("✅ Lot size %.2f is valid\n", lotSize)
```

---

### 10) Build a risk calculator tool

```go
func CalculateTradeRisk(sugar *mt5.MT5Sugar, symbol string, riskPercent, stopLossPips float64) {
    // Get account info
    balance, _ := sugar.GetBalance()

    // Calculate position size
    lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Calculate risk amount
    riskAmount := balance * riskPercent / 100.0

    // Calculate potential profit (if 1:2 R/R)
    potentialProfit := riskAmount * 2.0

    // Get current price
    ask, _ := sugar.GetAsk(symbol)

    // Calculate SL/TP prices (distances in points)
    sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, stopLossPips, stopLossPips*2)

    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║             TRADE RISK CALCULATOR                     ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")
    fmt.Println()
    fmt.Printf("Symbol:           %s\n", symbol)
    fmt.Printf("Current Price:    %.5f\n", ask)
    fmt.Println()
    fmt.Printf("Account Balance:  $%.2f\n", balance)
    fmt.Printf("Risk Percent:     %.1f%%\n", riskPercent)
    fmt.Printf("Risk Amount:      $%.2f\n", riskAmount)
    fmt.Println()
    fmt.Printf("Stop Loss:        %.0f points (%.5f)\n", stopLossPips, sl)
    fmt.Printf("Take Profit:      %.0f points (%.5f)\n", stopLossPips*2, tp)
    fmt.Printf("Risk/Reward:      1:2\n")
    fmt.Println()
    fmt.Printf("Lot Size:         %.2f\n", lotSize)
    fmt.Printf("Potential Loss:   $%.2f (if SL hit)\n", riskAmount)
    fmt.Printf("Potential Profit: $%.2f (if TP hit)\n", potentialProfit)
    fmt.Println()
}

// Usage:
CalculateTradeRisk(sugar, "EURUSD", 2.0, 500)
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `GetBalance()` - Get current account balance
* `GetSymbolInfo()` - Get symbol specifications
* `service.GetSymbolDouble()` - Tick value and tick size (point value in the account currency)
* `service.GetSymbolTick()` - Get current market price (ASK for the margin of 1 lot)
* `GetFreeMargin()`, `service.CalculateMargin()` - Margin cap

**🍬 Complementary sugar methods:**

* `CanOpenPosition()` - Validate position before opening ⭐
* `GetMaxLotSize()` - Calculate maximum tradeable volume
* `CalculateRequiredMargin()` - Calculate margin needed
* `BuyMarketWithPips()` - Open BUY with calculated lot size
* `SellMarketWithPips()` - Open SELL with calculated lot size

**Recommended workflow:**
```go
// 1. Calculate position size (500-point SL)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)

// 2. Validate
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    fmt.Println("Cannot open:", reason)
    return
}

// 3. Trade
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 500, 1000) // SL/TP in points
```

---

## ⚠️ Common Pitfalls

### 1) Using fixed lot sizes instead of calculated

```go
// ❌ WRONG - fixed lot size (ignores balance changes)
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)

// ✅ CORRECT - dynamic lot size based on risk
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
ticket, _ := sugar.BuyMarket("EURUSD", lotSize)
```

### 2) Confusing points, pips and price

```go
// ❌ WRONG - passing price as stop loss
slPrice := 1.08000
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, slPrice) // WRONG!

// ❌ WRONG - passing pips: 50 is read as 50 POINTS (= 5 pips) -> 10x too large a position
lotSize, _ = sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// ✅ CORRECT - stop loss distance in points (500 points = 50 pips on 5-digit EURUSD)
stopLossPips := 500.0
lotSize, _ = sugar.CalculatePositionSize("EURUSD", 2.0, stopLossPips)
```

### 3) Not checking for errors

```go
// ❌ WRONG - ignoring errors
lotSize, _ := sugar.CalculatePositionSize("INVALID", 2.0, 500)
sugar.BuyMarket("INVALID", lotSize) // Will fail!

// ✅ CORRECT - check errors
lotSize, err := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
```

### 4) Risking too much

```go
// ❌ WRONG - risking 10% per trade (very dangerous!)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 10.0, 500)

// ✅ CORRECT - professional risk management (1-2% max)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
```

### 5) Not validating before trading

```go
// ❌ WRONG - trading without validation
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
sugar.BuyMarket("EURUSD", lotSize) // Might fail!

// ✅ CORRECT - validate first
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    fmt.Println("Cannot trade:", reason)
    return
}
sugar.BuyMarket("EURUSD", lotSize)
```

### 6) Calculating manually instead of using Sugar

```go
// ❌ WRONG - manual calculation (error-prone)
balance, _ := sugar.GetBalance()
riskAmount := balance * 0.02
lotSize := riskAmount / (50 * 10) // Pips or points? Point value in which currency? Volume step?

// ✅ CORRECT - use Sugar (handles everything)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
```

---

## 💎 Pro Tips

1. **Always use 1-2% risk** - Never more than 2% per trade
2. **Tighter SL = larger position** - But don't use unrealistic stops
3. **Validate after calculating** - Use `CanOpenPosition()`
4. **Different symbols need different sizes** - Don't copy lot sizes between symbols
5. **Account grows → position size grows** - That's the power of this method!
6. **Check for minimum volume** - Some brokers have high minimums

---

**See also:** [`CanOpenPosition.md`](CanOpenPosition.md), [`BuyMarketWithPips.md`](../11.%20Trading_Helpers/BuyMarketWithPips.md)
