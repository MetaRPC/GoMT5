# 💰 Calculate Required Margin (`CalculateRequiredMargin`)

> **Sugar method:** Calculates how much margin is required to open a position of specified size.

**API Information:**

* **Method:** `sugar.CalculateRequiredMargin(symbol, volume)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `service.GetSymbolTick()`, `service.CalculateMargin()`
* **Timeout:** 5 seconds
* **Considers:** Leverage, symbol specifications, current price

---

## 📋 Method Signature

```go
func (s *MT5Sugar) CalculateRequiredMargin(symbol string, volume float64) (float64, error)
```

---

## 🔽 Input

| Parameter | Type | Description |
|-----------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1, 1.0, 5.0) |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `margin` | `float64` | Required margin amount in account currency |
| `error` | `error` | Error if calculation fails |

---

## 💬 Just the Essentials

* **What it is:** Calculates margin needed to open a position BEFORE actually trading.
* **Why you need it:** Plan your trades, check if you have enough margin, manage account exposure.
* **Sanity check:** Considers leverage - higher leverage = less margin needed.

---

## 🎯 When to Use

✅ **Before trading** - Check if you have enough margin

✅ **Position planning** - Calculate total margin for multiple positions

✅ **Risk assessment** - Understand margin exposure

✅ **Account management** - Track margin utilization

---

## 🔢 How Margin Works

```
Margin calculation depends on:

1. Symbol contract size (e.g., EURUSD = 100,000 units)
2. Current market price
3. Lot size
4. Account leverage (e.g., 1:100)

Formula (simplified):
Margin = (ContractSize × Lots × Price) / Leverage

Example (EURUSD, 1:100 leverage):
- Contract size: 100,000 EUR
- Lot size: 1.0
- Price: 1.10000
- Leverage: 100

Margin = (100,000 × 1.0 × 1.10000) / 100 = $1,100
```

---

## 🔗 Usage Examples

### 1) Basic usage - calculate margin needed

```go
symbol := "EURUSD"
volume := 1.0

margin, err := sugar.CalculateRequiredMargin(symbol, volume)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Symbol:          %s\n", symbol)
fmt.Printf("Lot size:        %.2f\n", volume)
fmt.Printf("Required margin: $%.2f\n", margin)
```

---

### 2) Check if you have enough margin

```go
symbol := "EURUSD"
volume := 0.5

// Calculate required margin
requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, volume)

// Get current free margin
freeMargin, _ := sugar.GetFreeMargin()

fmt.Printf("Required margin: $%.2f\n", requiredMargin)
fmt.Printf("Free margin:     $%.2f\n", freeMargin)

if requiredMargin <= freeMargin {
    fmt.Println("✅ Sufficient margin - can open position")
    ticket, _ := sugar.BuyMarket(symbol, volume)
    fmt.Printf("Position #%d opened\n", ticket)
} else {
    deficit := requiredMargin - freeMargin
    fmt.Printf("❌ Insufficient margin (short $%.2f)\n", deficit)
}
```

---

### 3) Compare margin across different lot sizes

```go
symbol := "EURUSD"
lotSizes := []float64{0.01, 0.1, 0.5, 1.0, 5.0}

fmt.Printf("Margin Requirements for %s:\n", symbol)
fmt.Println("─────────────────────────────────────────")
fmt.Printf("%-10s  %-15s\n", "Lot Size", "Margin Needed")
fmt.Println("─────────────────────────────────────────")

for _, lotSize := range lotSizes {
    margin, err := sugar.CalculateRequiredMargin(symbol, lotSize)
    if err != nil {
        fmt.Printf("%.2f lots  Error: %v\n", lotSize, err)
        continue
    }

    fmt.Printf("%-10.2f  $%-14.2f\n", lotSize, margin)
}

// Output example:
// Margin Requirements for EURUSD:
// ─────────────────────────────────────────
// Lot Size    Margin Needed
// ─────────────────────────────────────────
// 0.01        $11.00
// 0.10        $110.00
// 0.50        $550.00
// 1.00        $1,100.00
// 5.00        $5,500.00
```

---

### 4) Calculate total margin for multiple positions

```go
type PositionPlan struct {
    Symbol string
    Volume float64
}

plans := []PositionPlan{
    {"EURUSD", 0.5},
    {"GBPUSD", 0.3},
    {"USDJPY", 0.2},
    {"XAUUSD", 0.1},
}

fmt.Println("Multi-Position Margin Calculation:")
fmt.Println("─────────────────────────────────────────")

totalMargin := 0.0

for _, plan := range plans {
    margin, err := sugar.CalculateRequiredMargin(plan.Symbol, plan.Volume)
    if err != nil {
        fmt.Printf("%s %.2f lots: Error - %v\n",
            plan.Symbol, plan.Volume, err)
        continue
    }

    totalMargin += margin

    fmt.Printf("%s %.2f lots: $%.2f\n",
        plan.Symbol, plan.Volume, margin)
}

fmt.Println("─────────────────────────────────────────")
fmt.Printf("Total margin needed: $%.2f\n\n", totalMargin)

// Check if we can open all positions
freeMargin, _ := sugar.GetFreeMargin()
fmt.Printf("Free margin:         $%.2f\n", freeMargin)

if totalMargin <= freeMargin {
    fmt.Println("✅ Can open all positions")
} else {
    shortfall := totalMargin - freeMargin
    fmt.Printf("❌ Short $%.2f for all positions\n", shortfall)
}
```

---

### 5) Margin utilization report

```go
symbol := "EURUSD"
volume := 1.0

// Calculate required margin
requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, volume)

// Get account info
balance, _ := sugar.GetBalance()
equity, _ := sugar.GetEquity()
margin, _ := sugar.GetMargin()
freeMargin, _ := sugar.GetFreeMargin()

// Calculate utilization
currentUtilization := (margin / equity) * 100
newUtilization := ((margin + requiredMargin) / equity) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MARGIN UTILIZATION REPORT        ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Balance:          $%.2f\n", balance)
fmt.Printf("Equity:           $%.2f\n", equity)
fmt.Printf("Used margin:      $%.2f\n", margin)
fmt.Printf("Free margin:      $%.2f\n\n", freeMargin)

fmt.Printf("Current usage:    %.1f%%\n", currentUtilization)
fmt.Println()

fmt.Printf("Planned position:\n")
fmt.Printf("  Symbol:         %s\n", symbol)
fmt.Printf("  Volume:         %.2f lots\n", volume)
fmt.Printf("  Margin needed:  $%.2f\n\n", requiredMargin)

fmt.Printf("After opening:\n")
fmt.Printf("  Total margin:   $%.2f\n", margin+requiredMargin)
fmt.Printf("  Utilization:    %.1f%%\n", newUtilization)

if newUtilization > 80 {
    fmt.Println("\n🔴 DANGER: Very high margin usage!")
} else if newUtilization > 50 {
    fmt.Println("\n🟠 WARNING: High margin usage")
} else {
    fmt.Println("\n✅ Safe margin level")
}
```

---

### 6) Find maximum affordable lot size

```go
func FindMaxAffordableLots(sugar *mt5.MT5Sugar, symbol string) (float64, error) {
    freeMargin, err := sugar.GetFreeMargin()
    if err != nil {
        return 0, err
    }

    // Get symbol info for volume step
    info, err := sugar.GetSymbolInfo(symbol)
    if err != nil {
        return 0, err
    }

    // Binary search for max lot size
    // Start with a reasonable estimate
    testVolume := 1.0

    // Find upper bound
    for {
        margin, _ := sugar.CalculateRequiredMargin(symbol, testVolume)
        if margin > freeMargin {
            break
        }
        testVolume *= 2
    }

    // Binary search between 0 and testVolume
    low := 0.0
    high := testVolume
    maxLots := 0.0

    for i := 0; i < 20; i++ { // 20 iterations should be enough
        mid := (low + high) / 2.0

        // Round to volume step
        steps := int(mid / info.VolumeStep)
        mid = float64(steps) * info.VolumeStep

        margin, _ := sugar.CalculateRequiredMargin(symbol, mid)

        if margin <= freeMargin {
            maxLots = mid
            low = mid
        } else {
            high = mid
        }
    }

    return maxLots, nil
}

// Usage:
maxLots, err := FindMaxAffordableLots(sugar, "EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
} else {
    fmt.Printf("Maximum affordable lot size: %.2f\n", maxLots)

    margin, _ := sugar.CalculateRequiredMargin("EURUSD", maxLots)
    freeMargin, _ := sugar.GetFreeMargin()

    fmt.Printf("Would use: $%.2f of $%.2f free margin\n",
        margin, freeMargin)
}
```

---

### 7) Compare margin across symbols

```go
volume := 1.0
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD", "BTCUSD"}

fmt.Printf("Margin comparison (1.0 lot):\n")
fmt.Println("─────────────────────────────────────────")

type SymbolMargin struct {
    Symbol string
    Margin float64
}

margins := []SymbolMargin{}

for _, symbol := range symbols {
    margin, err := sugar.CalculateRequiredMargin(symbol, volume)
    if err != nil {
        fmt.Printf("%s: Error - %v\n", symbol, err)
        continue
    }

    margins = append(margins, SymbolMargin{symbol, margin})
}

// Sort by margin (lowest to highest)
sort.Slice(margins, func(i, j int) bool {
    return margins[i].Margin < margins[j].Margin
})

for i, sm := range margins {
    fmt.Printf("%d. %-8s: $%.2f\n", i+1, sm.Symbol, sm.Margin)
}

// Output example:
// Margin comparison (1.0 lot):
// ─────────────────────────────────────────
// 1. EURUSD:   $1,100.00
// 2. USDJPY:   $1,200.00
// 3. GBPUSD:   $1,300.00
// 4. XAUUSD:   $2,000.00
// 5. BTCUSD:   $5,000.00
```

---

### 8) Margin-aware position sizer

```go
func CalculatePositionSizeWithMarginCheck(
    sugar *mt5.MT5Sugar,
    symbol string,
    riskPercent float64,
    stopLoss float64,
) (float64, error) {
    // Calculate risk-based size
    riskBasedSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
    if err != nil {
        return 0, err
    }

    // Calculate margin needed
    requiredMargin, err := sugar.CalculateRequiredMargin(symbol, riskBasedSize)
    if err != nil {
        return 0, err
    }

    // Check available margin
    freeMargin, err := sugar.GetFreeMargin()
    if err != nil {
        return 0, err
    }

    fmt.Printf("Risk-based size: %.2f lots\n", riskBasedSize)
    fmt.Printf("Margin needed:   $%.2f\n", requiredMargin)
    fmt.Printf("Free margin:     $%.2f\n", freeMargin)

    // If not enough margin, reduce size
    if requiredMargin > freeMargin {
        // Get max affordable
        maxLots, _ := sugar.GetMaxLotSize(symbol)

        fmt.Printf("⚠️  Insufficient margin - reducing to %.2f lots\n", maxLots)
        return maxLots, nil
    }

    // Check if using too much margin (>50%)
    marginUsage := (requiredMargin / freeMargin) * 100
    if marginUsage > 50 {
        reducedSize := riskBasedSize * 0.5
        fmt.Printf("⚠️  High margin usage (%.1f%%) - reducing to %.2f lots\n",
            marginUsage, reducedSize)
        return reducedSize, nil
    }

    fmt.Printf("✅ Margin OK (%.1f%% usage)\n", marginUsage)
    return riskBasedSize, nil
}

// Usage:
lotSize, err := CalculatePositionSizeWithMarginCheck(
    sugar,
    "EURUSD",
    2.0,  // 2% risk
    50.0, // 50 pip stop
)

if err != nil {
    fmt.Printf("Error: %v\n", err)
} else {
    fmt.Printf("Final lot size: %.2f\n", lotSize)
}
```

---

### 9) Margin stress test

```go
func MarginStressTest(sugar *mt5.MT5Sugar, symbol string) {
    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          MARGIN STRESS TEST                           ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")

    freeMargin, _ := sugar.GetFreeMargin()
    fmt.Printf("Free margin: $%.2f\n\n", freeMargin)

    fmt.Printf("%-10s  %-15s  %-15s  %-10s\n",
        "Lot Size", "Margin Need", "Free After", "Usage %")
    fmt.Println("─────────────────────────────────────────────────────────")

    lotSizes := []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0}

    for _, lotSize := range lotSizes {
        margin, err := sugar.CalculateRequiredMargin(symbol, lotSize)
        if err != nil {
            fmt.Printf("%.2f lots: Error - %v\n", lotSize, err)
            continue
        }

        freeAfter := freeMargin - margin
        usagePercent := (margin / freeMargin) * 100

        status := ""
        if usagePercent > 90 {
            status = " 🔴 DANGER"
        } else if usagePercent > 70 {
            status = " 🟠 HIGH"
        } else if usagePercent > 50 {
            status = " 🟡 MODERATE"
        } else {
            status = " ✅ SAFE"
        }

        fmt.Printf("%-10.2f  $%-14.2f  $%-14.2f  %-10.1f%s\n",
            lotSize, margin, freeAfter, usagePercent, status)

        if freeAfter < 0 {
            break
        }
    }
}

// Usage:
MarginStressTest(sugar, "EURUSD")

// Output example:
// ╔═══════════════════════════════════════════════════════╗
// ║          MARGIN STRESS TEST                           ║
// ╚═══════════════════════════════════════════════════════╝
// Free margin: $10000.00
//
// Lot Size    Margin Need      Free After      Usage %
// ─────────────────────────────────────────────────────────
// 0.10        $110.00          $9,890.00       1.1        ✅ SAFE
// 0.50        $550.00          $9,450.00       5.5        ✅ SAFE
// 1.00        $1,100.00        $8,900.00       11.0       ✅ SAFE
// 2.00        $2,200.00        $7,800.00       22.0       ✅ SAFE
// 5.00        $5,500.00        $4,500.00       55.0       🟡 MODERATE
// 10.00       $11,000.00       $-1,000.00      110.0      🔴 DANGER
```

---

### 10) Advanced margin calculator

```go
type MarginCalculator struct {
    sugar *mt5.MT5Sugar
}

func NewMarginCalculator(sugar *mt5.MT5Sugar) *MarginCalculator {
    return &MarginCalculator{sugar: sugar}
}

func (mc *MarginCalculator) CalculateForPositions(
    positions []struct {
        Symbol string
        Volume float64
    },
) (float64, map[string]float64, error) {
    totalMargin := 0.0
    breakdown := make(map[string]float64)

    for _, pos := range positions {
        margin, err := mc.sugar.CalculateRequiredMargin(pos.Symbol, pos.Volume)
        if err != nil {
            return 0, nil, fmt.Errorf("failed for %s: %w", pos.Symbol, err)
        }

        totalMargin += margin
        breakdown[pos.Symbol] = margin
    }

    return totalMargin, breakdown, nil
}

func (mc *MarginCalculator) GetMarginLeverage(symbol string, volume float64) (float64, error) {
    // Calculate actual leverage being used
    margin, err := mc.sugar.CalculateRequiredMargin(symbol, volume)
    if err != nil {
        return 0, err
    }

    info, err := mc.sugar.GetSymbolInfo(symbol)
    if err != nil {
        return 0, err
    }

    // Position value = ContractSize * Volume * Price
    positionValue := info.ContractSize * volume * info.Ask

    // Leverage = Position Value / Margin
    leverage := positionValue / margin

    return leverage, nil
}

func (mc *MarginCalculator) ShowDetailedReport(symbol string, volume float64) {
    margin, _ := mc.sugar.CalculateRequiredMargin(symbol, volume)
    info, _ := mc.sugar.GetSymbolInfo(symbol)
    leverage, _ := mc.GetMarginLeverage(symbol, volume)

    positionValue := info.ContractSize * volume * info.Ask

    freeMargin, _ := mc.sugar.GetFreeMargin()
    balance, _ := mc.sugar.GetBalance()

    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          MARGIN CALCULATION DETAILS                   ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")
    fmt.Printf("Symbol:           %s\n", symbol)
    fmt.Printf("Volume:           %.2f lots\n\n", volume)

    fmt.Println("Position Details:")
    fmt.Printf("  Contract size:  %.0f\n", info.ContractSize)
    fmt.Printf("  Current price:  %.5f\n", info.Ask)
    fmt.Printf("  Position value: $%.2f\n\n", positionValue)

    fmt.Println("Margin Calculation:")
    fmt.Printf("  Required:       $%.2f\n", margin)
    fmt.Printf("  Leverage:       1:%.0f\n\n", leverage)

    fmt.Println("Account Status:")
    fmt.Printf("  Balance:        $%.2f\n", balance)
    fmt.Printf("  Free margin:    $%.2f\n", freeMargin)

    if margin <= freeMargin {
        marginUsage := (margin / freeMargin) * 100
        fmt.Printf("\n✅ Can open (%.1f%% of free margin)\n", marginUsage)
    } else {
        shortfall := margin - freeMargin
        fmt.Printf("\n❌ Cannot open (short $%.2f)\n", shortfall)
    }
}

// Usage:
calculator := NewMarginCalculator(sugar)

// Show detailed report
calculator.ShowDetailedReport("EURUSD", 1.0)

// Calculate for multiple positions
positions := []struct {
    Symbol string
    Volume float64
}{
    {"EURUSD", 0.5},
    {"GBPUSD", 0.3},
    {"USDJPY", 0.2},
}

totalMargin, breakdown, _ := calculator.CalculateForPositions(positions)

fmt.Printf("\nTotal margin for all positions: $%.2f\n", totalMargin)
fmt.Println("\nBreakdown:")
for symbol, margin := range breakdown {
    fmt.Printf("  %s: $%.2f\n", symbol, margin)
}
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `service.GetSymbolTick()` - Get current price
* `service.CalculateMargin()` - Perform margin calculation

**🍬 Complementary sugar methods:**

* `GetFreeMargin()` - Get available margin ⭐
* `GetMargin()` - Get currently used margin
* `GetMaxLotSize()` - Calculate max affordable volume ⭐
* `CanOpenPosition()` - Validate if position can be opened ⭐
* `CalculatePositionSize()` - Calculate risk-based size

**Recommended workflow:**
```go
// 1. Calculate position size
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// 2. Check margin requirement
margin, _ := sugar.CalculateRequiredMargin("EURUSD", lotSize)
freeMargin, _ := sugar.GetFreeMargin()

// 3. Validate
if margin <= freeMargin {
    canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
    if canOpen {
        ticket, _ := sugar.BuyMarket("EURUSD", lotSize)
    }
}
```

---

## ⚠️ Common Pitfalls

### 1) Not checking against free margin

```go
// ❌ WRONG - calculating but not checking
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 1.0)
sugar.BuyMarket("EURUSD", 1.0) // Might fail!

// ✅ CORRECT - check before trading
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 1.0)
freeMargin, _ := sugar.GetFreeMargin()
if margin <= freeMargin {
    sugar.BuyMarket("EURUSD", 1.0)
}
```

### 2) Confusing with balance

```go
// ❌ WRONG - comparing with balance
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 1.0)
balance, _ := sugar.GetBalance()
if margin <= balance { // WRONG comparison!

// ✅ CORRECT - compare with free margin
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 1.0)
freeMargin, _ := sugar.GetFreeMargin()
if margin <= freeMargin { // CORRECT!
```

### 3) Not accounting for existing positions

```go
// ❌ WRONG - only checking single position
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 1.0)
// But you already have positions using margin!

// ✅ CORRECT - consider total exposure
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 1.0)
currentMargin, _ := sugar.GetMargin()
freeMargin, _ := sugar.GetFreeMargin()
totalAfter := currentMargin + margin
fmt.Printf("Total margin after: $%.2f\n", totalAfter)
```

### 4) Ignoring leverage changes

```go
// ❌ WRONG - assuming fixed margin calculation
// Margin can change if broker adjusts leverage

// ✅ CORRECT - calculate right before trading
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 1.0)
// Fresh calculation with current leverage
```

### 5) Not using CanOpenPosition instead

```go
// ❌ WRONG - manually checking everything
margin, _ := sugar.CalculateRequiredMargin("EURUSD", 0.1)
freeMargin, _ := sugar.GetFreeMargin()
info, _ := sugar.GetSymbolInfo("EURUSD")
// Check volume limits, margin, etc...

// ✅ CORRECT - use CanOpenPosition (does everything)
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 0.1)
if !canOpen {
    fmt.Println(reason)
}
```

---

## 💎 Pro Tips

1. **Always compare with free margin** - Not balance or equity

2. **Calculate right before trading** - Margin requirements can change

3. **Use CanOpenPosition for validation** - It checks margin + other things

4. **Consider existing positions** - Don't look at margins in isolation

5. **Different symbols, different margins** - XAUUSD needs more than EURUSD

6. **Leverage matters** - 1:100 vs 1:500 = big difference in margin

7. **Use for planning** - Great for "what-if" scenarios

---

## 📊 Margin vs Other Checks

```
This method:
✅ Calculates margin requirement
✅ Considers leverage
✅ Uses current price
❌ Doesn't check free margin
❌ Doesn't validate volume limits
❌ Doesn't check symbol availability

CanOpenPosition:
✅ Everything this method does PLUS:
✅ Checks free margin sufficiency
✅ Validates volume limits
✅ Checks symbol availability
✅ Returns detailed reason if blocked

Use this when: You only need the margin amount
Use CanOpenPosition when: You want full validation
```

---

**See also:** [`CanOpenPosition.md`](CanOpenPosition.md), [`GetMaxLotSize.md`](GetMaxLotSize.md)
