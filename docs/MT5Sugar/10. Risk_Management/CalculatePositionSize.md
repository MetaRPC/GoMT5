# ⚖️ Calculate Position Size Based on Risk (`CalculatePositionSize`)

> **Sugar method:** THE MOST IMPORTANT risk management tool - automatically calculates position size so you risk exactly X% of your balance on the trade.

**API Information:**

* **Method:** `sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `GetBalance()`, `GetSymbolInfo()`, `GetSymbolTick()`, `CalculateMargin()`, `GetFreeMargin()`
* **Timeout:** 10 seconds

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
| `stopLossPips` | `float64` | Stop Loss distance in **points** (not price!) |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `lotSize` | `float64` | Recommended lot size to risk exactly `riskPercent` |
| `error` | `error` | Error if calculation fails or symbol not found |

**Returned lot size:**

- Already rounded to symbol's `VolumeStep`
- Clamped between `VolumeMin` and `VolumeMax`
- Ready to use in `BuyMarket()` / `SellMarket()` calls

---

## 💬 Just the Essentials

* **What it is:** Auto-calculates lot size based on how much you want to risk (% of balance).
* **Why you need it:** **NEVER use fixed lot sizes!** This ensures you risk the same % on every trade regardless of balance.
* **Sanity check:** If balance = $10,000, risk = 2%, SL = 50 pips → you'll always risk exactly $200, even as balance grows/shrinks.

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
LotSize = (Balance × RiskPercent / 100) / (StopLossPips × PipValue)

Where:
  PipValue = ContractSize × Point

Then:
  - Round to VolumeStep
  - Clamp between VolumeMin and VolumeMax
```

**Example:**
```
Balance:     $10,000
Risk:        2% ($200)
Stop Loss:   50 pips
Pip Value:   $10 per lot (for EURUSD)

LotSize = $200 / (50 × $10) = $200 / $500 = 0.04 lots
```

---

## 🧩 Notes & Tips

* **Always use this method** - Don't hardcode lot sizes!
* **Risk percent** - Professional traders use 1-2% max per trade
* **Stop loss in pips** - Use points, not price! (50 = 50 points)
* **Symbol matters** - Different symbols have different pip values
* **Auto-rounding** - Result is already rounded to valid volume step
* **Safety limits** - Result is clamped to broker's min/max volume
* **Account currency** - Works with any account currency (USD, EUR, etc.)

---

## 🔧 Under the Hood

This method performs several steps:

```go
// 1. Get current balance
balance, _ := sugar.GetBalance()

// 2. Get symbol information (point size, contract size, volume limits)
info, _ := sugar.GetSymbolInfo(symbol)

// 3. Get current market price
tick, _ := service.GetSymbolTick(ctx, symbol)

// 4. Calculate pip value per lot
pipValue := info.ContractSize * info.Point

// 5. Calculate risk amount in money
riskAmount := balance * riskPercent / 100.0

// 6. Calculate position size
positionSize := riskAmount / (stopLossPips * pipValue)

// 7. Round to volume step
positionSize = round(positionSize / info.VolumeStep) * info.VolumeStep

// 8. Clamp to min/max
positionSize = clamp(positionSize, info.VolumeMin, info.VolumeMax)

return positionSize
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
// Get balance
balance, _ := service.GetAccountDouble(ctx, pb.ENUM_ACCOUNT_INFO_DOUBLE_ACCOUNT_BALANCE)

// Get symbol params
params, _ := service.GetSymbolParamsMany(ctx, &symbol, nil, nil, nil)
p := params[0]

// Get current price
tick, _ := service.GetSymbolTick(ctx, symbol)

// Calculate pip value
pipValue := p.TradeContractSize * p.Point

// Calculate risk amount
riskAmount := balance * 2.0 / 100.0

// Calculate lot size
lotSize := riskAmount / (50.0 * pipValue)

// Round to step
lotSize = float64(int(lotSize / p.VolumeStep)) * p.VolumeStep

// Clamp to limits
if lotSize < p.VolumeMin {
    lotSize = p.VolumeMin
}
if lotSize > p.VolumeMax {
    lotSize = p.VolumeMax
}
```

**WITH sugar:**
```go
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
```

**Benefits:**

* ✅ **One line vs 20+ lines**
* ✅ **No manual formula**
* ✅ **No rounding logic needed**
* ✅ **Automatic error handling**
* ✅ **Built-in validation**

---

## 🔗 Usage Examples

### 1) Basic usage - risk 2% with 50 pip stop

```go
// Risk exactly 2% of balance with 50 pip stop loss
lotSize, err := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Recommended lot size: %.2f\n", lotSize)
// Output: Recommended lot size: 0.04
```

---

### 2) Complete trade setup with risk management

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLossPips := 50.0
takeProfitPips := 100.0

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
fmt.Printf("   Risk: %.1f%% of balance (%.0f pips)\n", riskPercent, stopLossPips)
```

---

### 3) Conservative vs aggressive risk levels

```go
symbol := "EURUSD"
stopLoss := 50.0

// Conservative (1% risk)
conservativeLots, _ := sugar.CalculatePositionSize(symbol, 1.0, stopLoss)
fmt.Printf("Conservative (1%% risk): %.2f lots\n", conservativeLots)

// Moderate (2% risk)
moderateLots, _ := sugar.CalculatePositionSize(symbol, 2.0, stopLoss)
fmt.Printf("Moderate (2%% risk):     %.2f lots\n", moderateLots)

// Aggressive (5% risk) - NOT RECOMMENDED!
aggressiveLots, _ := sugar.CalculatePositionSize(symbol, 5.0, stopLoss)
fmt.Printf("Aggressive (5%% risk):   %.2f lots (⚠️ HIGH RISK!)\n", aggressiveLots)

// Output:
// Conservative (1% risk): 0.02 lots
// Moderate (2% risk):     0.04 lots
// Aggressive (5% risk):   0.10 lots (⚠️ HIGH RISK!)
```

---

### 4) Different stop loss distances

```go
symbol := "EURUSD"
riskPercent := 2.0

// Tight stop (25 pips) → larger position size
tightSL, _ := sugar.CalculatePositionSize(symbol, riskPercent, 25)
fmt.Printf("25 pip SL: %.2f lots\n", tightSL)

// Normal stop (50 pips)
normalSL, _ := sugar.CalculatePositionSize(symbol, riskPercent, 50)
fmt.Printf("50 pip SL: %.2f lots\n", normalSL)

// Wide stop (100 pips) → smaller position size
wideSL, _ := sugar.CalculatePositionSize(symbol, riskPercent, 100)
fmt.Printf("100 pip SL: %.2f lots\n", wideSL)

// Output (for $10,000 balance):
// 25 pip SL: 0.08 lots
// 50 pip SL: 0.04 lots
// 100 pip SL: 0.02 lots
//
// Note: All risk exactly $200 (2% of $10,000)
```

---

### 5) Multi-symbol portfolio with consistent risk

```go
riskPercent := 2.0
stopLoss := 50.0

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

// Output:
// Consistent 2% risk across all symbols:
// ─────────────────────────────────────
// EURUSD: 0.04 lots
// GBPUSD: 0.03 lots
// USDJPY: 0.04 lots
// XAUUSD: 0.20 lots
//
// Each position risks exactly the same $ amount
```

---

### 6) Show risk amount in money

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLoss := 50.0

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
fmt.Printf("Stop Loss:     %.0f pips\n", stopLoss)
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
// Stop Loss:     50 pips
// ═══════════════════════════════════
// Lot Size:      0.04
// ═══════════════════════════════════
```

---

### 7) Account growth simulation

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLoss := 50.0

balances := []float64{1000, 5000, 10000, 50000, 100000}

fmt.Println("Position sizing as account grows:")
fmt.Println("──────────────────────────────────────────────")
fmt.Printf("%-12s  %-12s  %-12s\n", "Balance", "Risk (2%)", "Lot Size")
fmt.Println("──────────────────────────────────────────────")

for _, balance := range balances {
    // Create temporary calculation (in real code you'd have different accounts)
    riskAmount := balance * riskPercent / 100.0

    // For demonstration, calculate approximate lot size
    // In real code, use sugar.CalculatePositionSize()
    approxLots := riskAmount / (stopLoss * 10.0) // Simplified

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
stopLoss := 50.0
takeProfit := 100.0

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
fmt.Printf("   Stop Loss:  %.0f pips\n", stopLoss)
fmt.Printf("   Take Profit: %.0f pips\n", takeProfit)

ticket, _ := sugar.BuyMarketWithPips(symbol, lotSize, stopLoss, takeProfit)
fmt.Printf("   Ticket:     #%d\n", ticket)
```

---

### 9) Handle minimum volume edge case

```go
symbol := "EURUSD"
riskPercent := 0.5  // Very small risk
stopLoss := 100.0   // Wide stop

lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Get symbol info to check minimum
info, _ := sugar.GetSymbolInfo(symbol)

if lotSize <= info.VolumeMin {
    fmt.Printf("⚠️  Calculated lot size (%.2f) is at minimum (%.2f)\n",
        lotSize, info.VolumeMin)
    fmt.Printf("   Consider:\n")
    fmt.Printf("   - Increasing risk %% (currently %.1f%%)\n", riskPercent)
    fmt.Printf("   - Reducing stop loss (currently %.0f pips)\n", stopLoss)
    fmt.Printf("   - Trading on a larger account\n")
} else {
    fmt.Printf("✅ Lot size %.2f is valid (min: %.2f)\n", lotSize, info.VolumeMin)
}
```

---

### 10) Build a risk calculator tool

```go
func CalculateTradeRisk(symbol string, riskPercent, stopLossPips float64) {
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

    // Calculate SL/TP prices
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
    fmt.Printf("Stop Loss:        %.0f pips (%.5f)\n", stopLossPips, sl)
    fmt.Printf("Take Profit:      %.0f pips (%.5f)\n", stopLossPips*2, tp)
    fmt.Printf("Risk/Reward:      1:2\n")
    fmt.Println()
    fmt.Printf("Lot Size:         %.2f\n", lotSize)
    fmt.Printf("Potential Loss:   $%.2f (if SL hit)\n", riskAmount)
    fmt.Printf("Potential Profit: $%.2f (if TP hit)\n", potentialProfit)
    fmt.Println()
}

// Usage:
CalculateTradeRisk("EURUSD", 2.0, 50)
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `GetBalance()` - Get current account balance
* `GetSymbolInfo()` - Get symbol specifications
* `service.GetSymbolTick()` - Get current market price

**🍬 Complementary sugar methods:**

* `CanOpenPosition()` - Validate position before opening ⭐
* `GetMaxLotSize()` - Calculate maximum tradeable volume
* `CalculateRequiredMargin()` - Calculate margin needed
* `BuyMarketWithPips()` - Open BUY with calculated lot size
* `SellMarketWithPips()` - Open SELL with calculated lot size

**Recommended workflow:**
```go
// 1. Calculate position size
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// 2. Validate
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    fmt.Println("Cannot open:", reason)
    return
}

// 3. Trade
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

---

## ⚠️ Common Pitfalls

### 1) Using fixed lot sizes instead of calculated

```go
// ❌ WRONG - fixed lot size (ignores balance changes)
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)

// ✅ CORRECT - dynamic lot size based on risk
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, _ := sugar.BuyMarket("EURUSD", lotSize)
```

### 2) Confusing pips and price

```go
// ❌ WRONG - passing price as stop loss
slPrice := 1.08000
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, slPrice) // WRONG!

// ✅ CORRECT - stop loss in pips (points)
stopLossPips := 50.0
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, stopLossPips)
```

### 3) Not checking for errors

```go
// ❌ WRONG - ignoring errors
lotSize, _ := sugar.CalculatePositionSize("INVALID", 2.0, 50)
sugar.BuyMarket("INVALID", lotSize) // Will fail!

// ✅ CORRECT - check errors
lotSize, err := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
```

### 4) Risking too much

```go
// ❌ WRONG - risking 10% per trade (very dangerous!)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 10.0, 50)

// ✅ CORRECT - professional risk management (1-2% max)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
```

### 5) Not validating before trading

```go
// ❌ WRONG - trading without validation
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
sugar.BuyMarket("EURUSD", lotSize) // Might fail!

// ✅ CORRECT - validate first
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
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
lotSize := riskAmount / (50 * 10) // What about pip value? Volume step?

// ✅ CORRECT - use Sugar (handles everything)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
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
