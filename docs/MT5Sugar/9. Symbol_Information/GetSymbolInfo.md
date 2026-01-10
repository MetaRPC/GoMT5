# 🔍📊 Get Symbol Info (`GetSymbolInfo`)

> **Sugar method:** Returns comprehensive symbol information in one convenient structure.

**API Information:**

* **Method:** `sugar.GetSymbolInfo(symbol)`
* **Timeout:** 5 seconds
* **Returns:** SymbolInfo structure with all parameters

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "XAUUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `*SymbolInfo` | struct | Complete symbol information |
| `error` | `error` | Error if symbol not found |

---

## 📊 SymbolInfo Structure

```go
type SymbolInfo struct {
    Name         string  // Symbol name (e.g., "EURUSD")
    Bid          float64 // Current BID price
    Ask          float64 // Current ASK price
    Digits       int32   // Number of decimal places
    Point        float64 // Point size (minimal price change)
    VolumeMin    float64 // Minimum volume for trading
    VolumeMax    float64 // Maximum volume for trading
    VolumeStep   float64 // Volume step
    Spread       int32   // Current spread in points
    StopLevel    int32   // Minimum stop level in points
    ContractSize float64 // Contract size (for 1 lot)
}
```

---

## 💬 Just the Essentials

* **What it is:** Get all important symbol parameters in one call.
* **Why you need it:** Validation before trading, position sizing, SL/TP calculations.
* **Sanity check:** More efficient than calling individual methods for each property.

---

## 🎯 When to Use

✅ **Before trading** - Validate symbol parameters

✅ **Position sizing** - Calculate lot sizes with constraints

✅ **SL/TP validation** - Check stop level requirements

✅ **Symbol comparison** - Compare trading conditions

---

## 🔗 Usage Examples

### 1) Basic usage - get symbol info

```go
info, err := sugar.GetSymbolInfo("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Symbol: %s\n", info.Name)
fmt.Printf("Bid:    %.5f\n", info.Bid)
fmt.Printf("Ask:    %.5f\n", info.Ask)
fmt.Printf("Spread: %d points\n", info.Spread)
fmt.Printf("Digits: %d\n", info.Digits)
```

---

### 2) Validate symbol before trading

```go
func ValidateSymbol(sugar *mt5.MT5Sugar, symbol string, volume float64) (bool, string) {
    info, err := sugar.GetSymbolInfo(symbol)
    if err != nil {
        return false, fmt.Sprintf("symbol not found: %v", err)
    }

    // Check volume limits
    if volume < info.VolumeMin {
        return false, fmt.Sprintf("volume %.2f below minimum %.2f",
            volume, info.VolumeMin)
    }

    if volume > info.VolumeMax {
        return false, fmt.Sprintf("volume %.2f exceeds maximum %.2f",
            volume, info.VolumeMax)
    }

    // Check volume step
    steps := volume / info.VolumeStep
    if steps != float64(int(steps)) {
        return false, fmt.Sprintf("volume %.2f not multiple of step %.2f",
            volume, info.VolumeStep)
    }

    return true, "validated"
}

// Usage:
valid, message := ValidateSymbol(sugar, "EURUSD", 0.1)
if valid {
    fmt.Println("✅ Symbol validated")
} else {
    fmt.Printf("❌ Validation failed: %s\n", message)
}
```

---

### 3) Calculate SL/TP distances

```go
info, _ := sugar.GetSymbolInfo("EURUSD")

// Calculate SL distance in price
stopLossPips := 20.0
slDistancePoints := stopLossPips * 10 // For 5-digit broker

// Check against stop level
if slDistancePoints < float64(info.StopLevel) {
    fmt.Printf("⚠️  SL too close! Minimum: %d points\n", info.StopLevel)
    slDistancePoints = float64(info.StopLevel)
}

slDistancePrice := slDistancePoints * info.Point

fmt.Printf("Symbol: %s\n", info.Name)
fmt.Printf("Point size: %.5f\n", info.Point)
fmt.Printf("Min stop level: %d points\n", info.StopLevel)
fmt.Printf("SL distance: %.5f price units\n", slDistancePrice)
```

---

### 4) Display trading conditions

```go
info, _ := sugar.GetSymbolInfo("EURUSD")

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      TRADING CONDITIONS               ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Symbol:        %s\n\n", info.Name)

fmt.Printf("Current Prices:\n")
fmt.Printf("  Bid:         %.5f\n", info.Bid)
fmt.Printf("  Ask:         %.5f\n", info.Ask)
fmt.Printf("  Spread:      %d points\n\n", info.Spread)

fmt.Printf("Volume Limits:\n")
fmt.Printf("  Min:         %.2f lots\n", info.VolumeMin)
fmt.Printf("  Max:         %.2f lots\n", info.VolumeMax)
fmt.Printf("  Step:        %.2f lots\n\n", info.VolumeStep)

fmt.Printf("Trading Info:\n")
fmt.Printf("  Digits:      %d\n", info.Digits)
fmt.Printf("  Point:       %.5f\n", info.Point)
fmt.Printf("  Stop level:  %d points\n", info.StopLevel)
fmt.Printf("  Contract:    %.0f\n", info.ContractSize)
```

---

### 5) Compare multiple symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

fmt.Println("Symbol Comparison:")
fmt.Println("─────────────────────────────────────────")

for _, symbol := range symbols {
    info, err := sugar.GetSymbolInfo(symbol)
    if err != nil {
        fmt.Printf("%s: Error - %v\n", symbol, err)
        continue
    }

    spreadCost := float64(info.Spread) * info.Point

    fmt.Printf("%-8s: Spread %d pts (%.5f), Min lot %.2f\n",
        symbol, info.Spread, spreadCost, info.VolumeMin)
}
```

---

### 6) Calculate position value

```go
info, _ := sugar.GetSymbolInfo("EURUSD")

volume := 1.0 // 1 lot

positionValue := volume * info.ContractSize

fmt.Printf("Symbol:         %s\n", info.Name)
fmt.Printf("Volume:         %.2f lots\n", volume)
fmt.Printf("Contract size:  %.0f\n", info.ContractSize)
fmt.Printf("Position value: $%.2f\n", positionValue)

// Pip value (for 1 lot)
pipValue := info.Point * 10 * info.ContractSize
fmt.Printf("Pip value:      $%.2f per lot\n", pipValue)
```

---

### 7) Round volume to valid step

```go
func RoundVolumeToStep(volume float64, info *SymbolInfo) float64 {
    // Round to nearest valid step
    steps := volume / info.VolumeStep
    roundedSteps := math.Round(steps)
    rounded := roundedSteps * info.VolumeStep

    // Clamp to min/max
    if rounded < info.VolumeMin {
        return info.VolumeMin
    }
    if rounded > info.VolumeMax {
        return info.VolumeMax
    }

    return rounded
}

// Usage:
info, _ := sugar.GetSymbolInfo("EURUSD")

requestedVolume := 0.15 // Want 0.15 lots
validVolume := RoundVolumeToStep(requestedVolume, info)

fmt.Printf("Requested: %.2f lots\n", requestedVolume)
fmt.Printf("Valid:     %.2f lots (step: %.2f)\n",
    validVolume, info.VolumeStep)
```

---

### 8) Spread cost calculator

```go
info, _ := sugar.GetSymbolInfo("EURUSD")

volume := 1.0 // 1 lot

// Spread in price units
spreadPrice := float64(info.Spread) * info.Point

// Spread cost for position
spreadCost := spreadPrice * info.ContractSize * volume

fmt.Printf("Symbol:      %s\n", info.Name)
fmt.Printf("Spread:      %d points\n", info.Spread)
fmt.Printf("Volume:      %.2f lots\n", volume)
fmt.Printf("Spread cost: $%.2f\n", spreadCost)

if spreadCost > 10 {
    fmt.Println("⚠️  High spread - consider waiting")
}
```

---

### 9) Symbol quality checker

```go
func AssessSymbolQuality(info *SymbolInfo) string {
    issues := []string{}

    // Check spread
    if info.Spread > 30 {
        issues = append(issues, "High spread")
    }

    // Check minimum volume
    if info.VolumeMin > 0.01 {
        issues = append(issues, "High minimum volume")
    }

    // Check stop level
    if info.StopLevel > 50 {
        issues = append(issues, "Large stop level restriction")
    }

    if len(issues) == 0 {
        return "✅ Excellent trading conditions"
    } else if len(issues) == 1 {
        return fmt.Sprintf("🟡 Good (%s)", issues[0])
    } else {
        return fmt.Sprintf("⚠️  Fair (%s)", strings.Join(issues, ", "))
    }
}

// Usage:
info, _ := sugar.GetSymbolInfo("EURUSD")
assessment := AssessSymbolQuality(info)
fmt.Printf("%s: %s\n", info.Name, assessment)
```

---

### 10) Advanced symbol analyzer

```go
type SymbolAnalyzer struct {
    Info *SymbolInfo
}

func NewSymbolAnalyzer(sugar *mt5.MT5Sugar, symbol string) (*SymbolAnalyzer, error) {
    info, err := sugar.GetSymbolInfo(symbol)
    if err != nil {
        return nil, err
    }

    return &SymbolAnalyzer{Info: info}, nil
}

func (sa *SymbolAnalyzer) CalculatePipValue(volume float64) float64 {
    // For standard pairs (10 pips = 1 point movement)
    return sa.Info.Point * 10 * sa.Info.ContractSize * volume
}

func (sa *SymbolAnalyzer) CalculateSpreadCost(volume float64) float64 {
    spreadPrice := float64(sa.Info.Spread) * sa.Info.Point
    return spreadPrice * sa.Info.ContractSize * volume
}

func (sa *SymbolAnalyzer) ValidateVolume(volume float64) (float64, error) {
    // Check minimum
    if volume < sa.Info.VolumeMin {
        return 0, fmt.Errorf("volume %.2f below minimum %.2f",
            volume, sa.Info.VolumeMin)
    }

    // Check maximum
    if volume > sa.Info.VolumeMax {
        return 0, fmt.Errorf("volume %.2f exceeds maximum %.2f",
            volume, sa.Info.VolumeMax)
    }

    // Round to step
    steps := volume / sa.Info.VolumeStep
    rounded := math.Round(steps) * sa.Info.VolumeStep

    return rounded, nil
}

func (sa *SymbolAnalyzer) CalculateMinSLDistance() float64 {
    // Minimum SL distance in price units
    return float64(sa.Info.StopLevel) * sa.Info.Point
}

func (sa *SymbolAnalyzer) GenerateReport() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      SYMBOL ANALYSIS REPORT           ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Symbol:        %s\n\n", sa.Info.Name)

    // Current market
    fmt.Println("Market Snapshot:")
    fmt.Printf("  Bid:         %.5f\n", sa.Info.Bid)
    fmt.Printf("  Ask:         %.5f\n", sa.Info.Ask)
    fmt.Printf("  Spread:      %d points (%.5f)\n",
        sa.Info.Spread, float64(sa.Info.Spread)*sa.Info.Point)

    // Trading specs
    fmt.Println("\nTrading Specifications:")
    fmt.Printf("  Precision:   %d decimals\n", sa.Info.Digits)
    fmt.Printf("  Point:       %.5f\n", sa.Info.Point)
    fmt.Printf("  Stop level:  %d points (%.5f)\n",
        sa.Info.StopLevel, sa.CalculateMinSLDistance())

    // Volume constraints
    fmt.Println("\nVolume Constraints:")
    fmt.Printf("  Minimum:     %.2f lots\n", sa.Info.VolumeMin)
    fmt.Printf("  Maximum:     %.2f lots\n", sa.Info.VolumeMax)
    fmt.Printf("  Step:        %.2f lots\n", sa.Info.VolumeStep)

    // Contract
    fmt.Println("\nContract Information:")
    fmt.Printf("  Size:        %.0f units\n", sa.Info.ContractSize)
    fmt.Printf("  Pip value:   $%.2f per lot\n",
        sa.CalculatePipValue(1.0))

    // Cost analysis (1 lot)
    fmt.Println("\nCost Analysis (1 lot):")
    spreadCost := sa.CalculateSpreadCost(1.0)
    fmt.Printf("  Spread cost: $%.2f\n", spreadCost)

    // Assessment
    fmt.Println("\nAssessment:")
    if sa.Info.Spread <= 20 {
        fmt.Println("  ✅ Low spread - good for scalping")
    } else if sa.Info.Spread <= 50 {
        fmt.Println("  🟡 Moderate spread - suitable for day trading")
    } else {
        fmt.Println("  ⚠️  High spread - best for swing trading")
    }

    if sa.Info.VolumeMin <= 0.01 {
        fmt.Println("  ✅ Micro lots available - good for small accounts")
    }

    if sa.Info.StopLevel <= 30 {
        fmt.Println("  ✅ Low stop level - flexible SL/TP placement")
    } else {
        fmt.Println("  ⚠️  High stop level - limits tight stops")
    }
}

// Usage:
analyzer, _ := NewSymbolAnalyzer(sugar, "EURUSD")
analyzer.GenerateReport()

// Calculate for specific volume
pipValue := analyzer.CalculatePipValue(0.5)
fmt.Printf("\nPip value for 0.5 lots: $%.2f\n", pipValue)
```

---

## 🔗 Related Methods

**🍬 Individual properties:**

* `GetBid()` - Just BID price
* `GetAsk()` - Just ASK price
* `GetSpread()` - Just spread
* `GetSymbolDigits()` - Just digits
* `GetMinStopLevel()` - Just stop level

**🍬 Symbol availability:**

* `IsSymbolAvailable()` - Check if tradeable
* `GetAllSymbols()` - List all symbols

---

## ⚠️ Common Pitfalls

### 1) Symbol name case sensitivity

```go
// ❌ WRONG - case matters!
info, _ := sugar.GetSymbolInfo("eurusd")

// ✅ CORRECT - use exact broker symbol name
info, _ := sugar.GetSymbolInfo("EURUSD")
```

### 2) Not checking volume step

```go
// ❌ WRONG - ignoring volume step
volume := 0.15 // Might not be valid!

// ✅ CORRECT - round to step
info, _ := sugar.GetSymbolInfo("EURUSD")
steps := math.Round(volume / info.VolumeStep)
validVolume := steps * info.VolumeStep
```

### 3) Confusing points and pips

```go
// ❌ WRONG - points ≠ pips for 5-digit brokers
spread := info.Spread // This is in POINTS
// For EURUSD with 5 digits: 10 points = 1 pip

// ✅ CORRECT - convert if needed
spreadPips := float64(info.Spread) / 10.0 // For 5-digit
```

---

## 💎 Pro Tips

1. **Cache info** - Symbol info changes rarely, can cache for efficiency

2. **All-in-one** - More efficient than multiple individual calls

3. **Validation** - Always validate volumes against VolumeMin/Max/Step

4. **Stop level** - Check StopLevel before placing tight SL/TP

5. **Point size** - Use Point for price calculations, not hardcoded values

---

## 📊 Important Relationships

```
Price calculations:
- Spread in price = Spread points × Point size
- Pip value = Point × 10 × ContractSize (for 5-digit)
- SL distance = (SL pips × 10) × Point size

Volume validation:
- Volume must be ≥ VolumeMin
- Volume must be ≤ VolumeMax
- Volume must be multiple of VolumeStep

Stop validation:
- SL/TP distance must be ≥ StopLevel points
- StopLevel distance in price = StopLevel × Point
```

---

**See also:** [`GetAllSymbols.md`](GetAllSymbols.md), [`IsSymbolAvailable.md`](IsSymbolAvailable.md), [`GetMinStopLevel.md`](GetMinStopLevel.md), [`GetSymbolDigits.md`](GetSymbolDigits.md)
