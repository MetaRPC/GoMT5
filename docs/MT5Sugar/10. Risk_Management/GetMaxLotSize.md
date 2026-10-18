# 📊 Get Maximum Lot Size (`GetMaxLotSize`)

> **Sugar method:** Calculates the maximum lot size you can open based on your current free margin.

**API Information:**

* **Method:** `sugar.GetMaxLotSize(symbol)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `GetFreeMargin()`, `GetSymbolInfo()`, `service.GetSymbolTick()`, `service.CalculateMargin()`
* **Timeout:** 5 seconds (for the tick and margin queries; `GetFreeMargin` and `GetSymbolInfo` use their own 3 s / 10 s timeouts)
* **Safety buffer:** Uses 80% of free margin (20% buffer to prevent margin calls)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetMaxLotSize(symbol string) (float64, error)
```

---

## 🔽 Input

| Parameter | Type | Description |
|-----------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `maxLots` | `float64` | Maximum safe lot size you can trade |
| `error` | `error` | Error if a query fails, `failed to calculate margin for <symbol>` if the margin for 1 lot is 0, or `insufficient margin: 80% of the free margin X does not cover the minimum volume Y of <symbol>` |

**Special cases:**

- Returns an **error** (`insufficient margin: ...`) if you don't have enough margin even for the minimum lot - never `(0, nil)`
- Rounded DOWN to symbol's `VolumeStep` (same `floorToStep` as `CalculatePositionSize`, with a small float tolerance: 0.3 / 0.01 gives 0.30, not 0.29)
- Capped at `VolumeMax` (a result below `VolumeMin` is the error above, it is not raised to `VolumeMin`)
- Margin per lot is calculated for a **BUY** at the current **Ask** - the method has no direction, a SELL's margin may differ slightly

---

## 💬 Just the Essentials

* **What it is:** Calculates how much you can trade based on available margin.
* **Why you need it:** Prevents margin calls, shows your maximum trading capacity.
* **Sanity check:** Uses 80% of free margin (keeps 20% buffer for safety).

---

## 🎯 When to Use

✅ **Check trading capacity** - See max position size before trading

✅ **Prevent margin calls** - Ensure you don't over-leverage

✅ **Multi-position planning** - Calculate if you can open multiple positions

✅ **Account stress testing** - Understand your maximum exposure

---

## 🔢 How It Works

```
Step 1: Get your free margin (e.g., $5,000)
Step 2: Calculate margin required for 1 lot of the symbol (OrderCalcMargin, BUY at Ask)
Step 3: Apply safety buffer (80% of free margin)
Step 4: Calculate: maxLots = (freeMargin × 0.8) / marginPerLot
Step 5: Round down to VolumeStep
Step 6: Below VolumeMin → error "insufficient margin: ..."; above VolumeMax → return VolumeMax

Example:
- Free margin: $5,000
- Margin per lot: $500 (for EURUSD with 1:100 leverage)
- Safety margin: $5,000 × 0.8 = $4,000
- Max lots: $4,000 / $500 = 8.0 lots
```

---

## 🔗 Usage Examples

### 1) Basic usage - check max capacity

```go
maxLots, err := sugar.GetMaxLotSize("EURUSD")
if err != nil {
    // also "insufficient margin: ..." when not even the minimum lot is covered
    fmt.Printf("❌ %v\n", err)
    return
}

fmt.Printf("✅ Maximum lot size: %.2f lots\n", maxLots)
```

---

### 2) Safe position sizing

```go
symbol := "EURUSD"

// Get maximum capacity
maxLots, _ := sugar.GetMaxLotSize(symbol)

// Use only 50% of max capacity for safety
// (round to VolumeStep yourself - MT5Sugar sends the volume exactly as given)
info, _ := sugar.GetSymbolInfo(symbol)
safeLots := float64(int(maxLots*0.5/info.VolumeStep)) * info.VolumeStep

fmt.Printf("Maximum capacity:  %.2f lots\n", maxLots)
fmt.Printf("Safe position:     %.2f lots (50%% of max)\n", safeLots)

// Trade with safe size
if safeLots >= info.VolumeMin {
    ticket, _ := sugar.BuyMarket(symbol, safeLots)
    fmt.Printf("Opened position #%d\n", ticket)
}
```

---

### 3) Multi-position planning

```go
symbol := "EURUSD"
numPositions := 3

// Get max capacity
maxLots, _ := sugar.GetMaxLotSize(symbol)

// Divide among positions
lotsPerPosition := maxLots / float64(numPositions)

// Round to symbol's step
info, _ := sugar.GetSymbolInfo(symbol)
steps := int(lotsPerPosition / info.VolumeStep)
lotsPerPosition = float64(steps) * info.VolumeStep

fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("MULTI-POSITION PLANNING\n")
fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("Max capacity:      %.2f lots\n", maxLots)
fmt.Printf("Desired positions: %d\n", numPositions)
fmt.Printf("Lots per position: %.2f lots\n", lotsPerPosition)
fmt.Printf("═══════════════════════════════════════\n")

if lotsPerPosition >= info.VolumeMin {
    fmt.Println("✅ Can open all positions")
} else {
    fmt.Println("❌ Insufficient margin for desired setup")
}
```

---

### 4) Compare calculated size vs max capacity

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLoss := 500.0 // points

// Get risk-based position size
calculatedLots, _ := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)

// Get maximum capacity
maxLots, _ := sugar.GetMaxLotSize(symbol)

fmt.Printf("Risk-based size: %.2f lots\n", calculatedLots)
fmt.Printf("Max capacity:    %.2f lots\n", maxLots)

if calculatedLots > maxLots {
    fmt.Printf("⚠️  Warning: Risk-based size (%.2f) exceeds capacity (%.2f)!\n",
        calculatedLots, maxLots)
    fmt.Printf("   Using maximum: %.2f lots\n", maxLots)
    calculatedLots = maxLots
} else {
    usagePercent := (calculatedLots / maxLots) * 100
    fmt.Printf("✅ Using %.1f%% of capacity\n", usagePercent)
}
```

---

### 5) Margin utilization report

```go
symbol := "EURUSD"

// Get account info
freeMargin, _ := sugar.GetFreeMargin()
margin, _ := sugar.GetMargin()
balance, _ := sugar.GetBalance()

// Get max lot size
maxLots, _ := sugar.GetMaxLotSize(symbol)

// Calculate margin for max position
requiredForMax, _ := sugar.CalculateRequiredMargin(symbol, maxLots)

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MARGIN UTILIZATION REPORT        ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Balance:           $%.2f\n", balance)
fmt.Printf("Used margin:       $%.2f\n", margin)
fmt.Printf("Free margin:       $%.2f\n", freeMargin)
fmt.Println()
fmt.Printf("Symbol:            %s\n", symbol)
fmt.Printf("Max lot size:      %.2f lots\n", maxLots)
fmt.Printf("Margin for max:    $%.2f (80%% buffer)\n", requiredForMax)
fmt.Println()

utilization := (margin / balance) * 100
fmt.Printf("Current usage:     %.1f%% of balance\n", utilization)

if utilization > 50 {
    fmt.Println("⚠️  High margin usage - reduce positions!")
} else {
    fmt.Println("✅ Healthy margin level")
}
```

---

### 6) Check before opening additional position

```go
func CanOpenAdditionalPosition(sugar *mt5.MT5Sugar, symbol string, desiredLots float64) bool {
    maxLots, err := sugar.GetMaxLotSize(symbol)
    if err != nil {
        fmt.Printf("Error checking capacity: %v\n", err)
        return false
    }

    if desiredLots > maxLots {
        fmt.Printf("❌ Cannot open %.2f lots (max: %.2f)\n", desiredLots, maxLots)
        return false
    }

    // Use 90% of max as threshold for safety
    safeThreshold := maxLots * 0.9

    if desiredLots > safeThreshold {
        fmt.Printf("⚠️  Warning: %.2f lots is %.1f%% of capacity\n",
            desiredLots, (desiredLots/maxLots)*100)
        fmt.Println("   Consider reducing position size")
    }

    fmt.Printf("✅ Can open %.2f lots (max: %.2f)\n", desiredLots, maxLots)
    return true
}

// Usage:
if CanOpenAdditionalPosition(sugar, "EURUSD", 0.5) {
    ticket, _ := sugar.BuyMarket("EURUSD", 0.5)
    fmt.Printf("Position opened: #%d\n", ticket)
}
```

---

### 7) Leverage impact demonstration

```go
symbols := []string{"EURUSD", "GBPUSD", "XAUUSD", "BTCUSD"}

fmt.Println("Maximum Lot Sizes by Symbol:")
fmt.Println("─────────────────────────────────────────")

for _, symbol := range symbols {
    maxLots, err := sugar.GetMaxLotSize(symbol)
    if err != nil {
        fmt.Printf("%s: Error - %v\n", symbol, err)
        continue
    }

    // Get margin for 1 lot
    marginPerLot, _ := sugar.CalculateRequiredMargin(symbol, 1.0)

    fmt.Printf("%-8s: %.2f lots (margin per lot: $%.2f)\n",
        symbol, maxLots, marginPerLot)
}

// Output example:
// Maximum Lot Sizes by Symbol:
// ─────────────────────────────────────────
// EURUSD:   8.00 lots (margin per lot: $500.00)
// GBPUSD:   6.50 lots (margin per lot: $615.38)
// XAUUSD:   4.00 lots (margin per lot: $1000.00)
// BTCUSD:   1.00 lots (margin per lot: $4000.00)
```

---

### 8) Dynamic position sizing with capacity check

```go
func OpenPositionWithCapacityCheck(
    sugar *mt5.MT5Sugar,
    symbol string,
    riskPercent float64,
    stopLoss float64,
) (uint64, error) {
    // Calculate risk-based size
    lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
    if err != nil {
        return 0, fmt.Errorf("position size calculation failed: %w", err)
    }

    // Check against maximum capacity
    maxLots, err := sugar.GetMaxLotSize(symbol)
    if err != nil {
        return 0, fmt.Errorf("max lot size check failed: %w", err)
    }

    // Use smaller of the two
    finalLots := lotSize
    if lotSize > maxLots {
        fmt.Printf("⚠️  Risk-based size %.2f exceeds capacity %.2f\n",
            lotSize, maxLots)
        fmt.Printf("   Reducing to maximum: %.2f lots\n", maxLots)
        finalLots = maxLots
    }

    // Validate
    canOpen, reason, _ := sugar.CanOpenPosition(symbol, finalLots)
    if !canOpen {
        return 0, fmt.Errorf("cannot open position: %s", reason)
    }

    // Open position
    return sugar.BuyMarket(symbol, finalLots)
}

// Usage:
ticket, err := OpenPositionWithCapacityCheck(sugar, "EURUSD", 2.0, 500) // 500-point SL
if err != nil {
    fmt.Printf("Failed: %v\n", err)
} else {
    fmt.Printf("Success: Position #%d\n", ticket)
}
```

---

### 9) Account capacity monitor

```go
func MonitorAccountCapacity(sugar *mt5.MT5Sugar, symbols []string) {
    balance, _ := sugar.GetBalance()
    freeMargin, _ := sugar.GetFreeMargin()

    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          ACCOUNT CAPACITY MONITOR                     ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")
    fmt.Printf("Balance:      $%.2f\n", balance)
    fmt.Printf("Free Margin:  $%.2f\n\n", freeMargin)

    fmt.Printf("%-10s  %-12s  %-15s\n", "Symbol", "Max Lots", "Margin Need")
    fmt.Println("─────────────────────────────────────────────────────────")

    totalCapacity := 0.0

    for _, symbol := range symbols {
        maxLots, err := sugar.GetMaxLotSize(symbol)
        if err != nil {
            fmt.Printf("%-10s  Error: %v\n", symbol, err)
            continue
        }

        marginNeed, _ := sugar.CalculateRequiredMargin(symbol, maxLots)
        totalCapacity += marginNeed

        fmt.Printf("%-10s  %-12.2f  $%-14.2f\n",
            symbol, maxLots, marginNeed)
    }

    fmt.Println("─────────────────────────────────────────────────────────")
    fmt.Printf("Note: Using 80%% safety buffer\n")

    if totalCapacity > freeMargin {
        fmt.Println("\n⚠️  Warning: Cannot open all positions at max simultaneously")
    } else {
        fmt.Println("\n✅ Sufficient margin for all symbols")
    }
}

// Usage:
watchlist := []string{"EURUSD", "GBPUSD", "USDJPY"}
MonitorAccountCapacity(sugar, watchlist)
```

---

### 10) Advanced capacity manager

```go
type CapacityManager struct {
    sugar *mt5.MT5Sugar
}

func NewCapacityManager(sugar *mt5.MT5Sugar) *CapacityManager {
    return &CapacityManager{sugar: sugar}
}

func (cm *CapacityManager) GetSafetyLevel() string {
    freeMargin, _ := cm.sugar.GetFreeMargin()
    balance, _ := cm.sugar.GetBalance()

    ratio := freeMargin / balance

    if ratio > 0.8 {
        return "✅ Excellent - Low risk"
    } else if ratio > 0.5 {
        return "🟡 Good - Moderate usage"
    } else if ratio > 0.3 {
        return "🟠 Caution - High usage"
    } else {
        return "🔴 Danger - Very high usage"
    }
}

func (cm *CapacityManager) CanOpenMultiplePositions(
    positions []struct {
        Symbol string
        Lots   float64
    },
) (bool, string) {
    freeMargin, _ := cm.sugar.GetFreeMargin()

    totalMarginNeeded := 0.0

    for _, pos := range positions {
        marginNeeded, err := cm.sugar.CalculateRequiredMargin(pos.Symbol, pos.Lots)
        if err != nil {
            return false, fmt.Sprintf("failed to calculate margin for %s", pos.Symbol)
        }
        totalMarginNeeded += marginNeeded
    }

    // Use 90% of free margin as threshold
    safeLimit := freeMargin * 0.9

    if totalMarginNeeded > safeLimit {
        return false, fmt.Sprintf(
            "insufficient margin: need $%.2f, safe limit $%.2f",
            totalMarginNeeded, safeLimit,
        )
    }

    return true, ""
}

func (cm *CapacityManager) GetOptimalDistribution(
    symbols []string,
    totalRisk float64,
) map[string]float64 {
    distribution := make(map[string]float64)

    // Get max lots for each symbol
    maxLots := make(map[string]float64)
    totalCapacity := 0.0

    for _, symbol := range symbols {
        max, _ := cm.sugar.GetMaxLotSize(symbol)
        maxLots[symbol] = max

        margin, _ := cm.sugar.CalculateRequiredMargin(symbol, max)
        totalCapacity += margin
    }

    // Distribute proportionally
    for symbol, max := range maxLots {
        margin, _ := cm.sugar.CalculateRequiredMargin(symbol, max)
        proportion := margin / totalCapacity
        distribution[symbol] = max * proportion * totalRisk
    }

    return distribution
}

func (cm *CapacityManager) ShowReport() {
    balance, _ := cm.sugar.GetBalance()
    margin, _ := cm.sugar.GetMargin()
    freeMargin, _ := cm.sugar.GetFreeMargin()

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      CAPACITY MANAGER REPORT          ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Balance:       $%.2f\n", balance)
    fmt.Printf("Used Margin:   $%.2f (%.1f%%)\n",
        margin, (margin/balance)*100)
    fmt.Printf("Free Margin:   $%.2f (%.1f%%)\n",
        freeMargin, (freeMargin/balance)*100)
    fmt.Println()
    fmt.Printf("Safety Level:  %s\n", cm.GetSafetyLevel())
}

// Usage:
manager := NewCapacityManager(sugar)
manager.ShowReport()

// Check multiple positions
positions := []struct {
    Symbol string
    Lots   float64
}{
    {"EURUSD", 0.5},
    {"GBPUSD", 0.3},
    {"USDJPY", 0.2},
}

canOpen, reason := manager.CanOpenMultiplePositions(positions)
if !canOpen {
    fmt.Printf("Cannot open positions: %s\n", reason)
} else {
    fmt.Println("✅ Can open all positions")
}
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `GetFreeMargin()` - Get available margin
* `GetSymbolInfo()` - Get volume limits and step
* `service.CalculateMargin()` - Calculate margin for 1 lot

**🍬 Complementary sugar methods:**

* `CalculatePositionSize()` - Calculate size based on risk % ⭐
* `CanOpenPosition()` - Validate position can be opened ⭐
* `CalculateRequiredMargin()` - Calculate margin for specific volume
* `GetFreeMargin()` - Get current free margin

**Recommended workflow:**
```go
// 1. Calculate risk-based size (500-point SL)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)

// 2. Check against max capacity
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
if lotSize > maxLots {
    lotSize = maxLots
}

// 3. Validate
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    fmt.Println("Cannot open:", reason)
    return
}

// 4. Trade
ticket, _ := sugar.BuyMarket("EURUSD", lotSize)
```

---

## ⚠️ Common Pitfalls

### 1) Using 100% of max capacity

```go
// ❌ WRONG - using full capacity (risky!)
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
sugar.BuyMarket("EURUSD", maxLots) // Dangerous!

// ✅ CORRECT - use percentage of max for safety
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
info, _ := sugar.GetSymbolInfo("EURUSD")
safeLots := float64(int(maxLots*0.5/info.VolumeStep)) * info.VolumeStep // 50% of capacity, rounded down to step
sugar.BuyMarket("EURUSD", safeLots)
```

### 2) Ignoring the error

```go
// ❌ WRONG - ignoring the error (insufficient margin returns 0 AND an error)
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
sugar.BuyMarket("EURUSD", maxLots) // Will fail if maxLots = 0!

// ✅ CORRECT - check the error
maxLots, err := sugar.GetMaxLotSize("EURUSD")
if err != nil {
    fmt.Println("Cannot size the position:", err) // e.g. "insufficient margin: ..."
    return
}
```

### 3) Not considering existing positions

```go
// ❌ WRONG - checking capacity without existing positions context
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
// But you already have open positions using margin!

// ✅ CORRECT - consider total exposure
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
openPositions, _ := sugar.CountOpenPositions()
if openPositions >= 3 {
    fmt.Println("Too many positions already open")
    return
}
```

### 4) Confusing with CalculatePositionSize

```go
// ❌ WRONG - these are different things!
// GetMaxLotSize = what you CAN open (margin-based)
// CalculatePositionSize = what you SHOULD open (risk-based)

maxLots, _ := sugar.GetMaxLotSize("EURUSD") // Maybe 10 lots
// Just because you CAN doesn't mean you SHOULD!

// ✅ CORRECT - use risk-based, check against max
riskLots, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
lotSize := math.Min(riskLots, maxLots)
```

### 5) Not accounting for safety buffer

```go
// ❌ WRONG - thinking it uses 100% of free margin
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
// Actually uses 80% (20% buffer built-in)

// ✅ CORRECT - understand it has safety buffer
maxLots, _ := sugar.GetMaxLotSize("EURUSD")
fmt.Println("Max lots includes 20% safety buffer")
```

---

## 💎 Pro Tips

1. **Safety buffer** - Method uses 80% of free margin automatically

2. **Use percentage** - Never use 100% of returned capacity, use 50-70%

3. **Risk-based is better** - Use `CalculatePositionSize()` for risk %, use this for upper limit

4. **Check before adding** - Check capacity before opening additional positions

5. **Symbol matters** - Different symbols have different margin requirements

6. **Monitor regularly** - Capacity changes as positions are opened/closed

7. **Insufficient margin is an error** - When not even the minimum lot is covered, it returns `0` with an `insufficient margin: ...` error (not `0, nil`)

---

## 📊 Capacity vs Risk Sizing

```
Two different approaches:

GetMaxLotSize():
- "How much CAN I trade?"
- Based on margin availability
- Shows maximum capacity
- Use for upper limit checking

CalculatePositionSize():
- "How much SHOULD I trade?"
- Based on risk percentage
- Follows risk management rules
- Use for actual position sizing

Best practice: Use both together!
lotSize = min(riskBasedSize, maxCapacity * 0.5)
```

---

**See also:** [`CalculatePositionSize.md`](CalculatePositionSize.md), [`CanOpenPosition.md`](CanOpenPosition.md), [`CalculateRequiredMargin.md`](CalculateRequiredMargin.md)
