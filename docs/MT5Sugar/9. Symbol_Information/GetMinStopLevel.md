# 📏🛑 Get Min Stop Level (`GetMinStopLevel`)

> **Sugar method:** Returns the minimum allowed distance for Stop Loss/Take Profit in points.

**API Information:**

* **Method:** `sugar.GetMinStopLevel(symbol)`
* **Timeout:** 3 seconds
* **Returns:** Minimum stop level in points

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetMinStopLevel(symbol string) (int64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `int64` | integer | Minimum stop level in points |
| `error` | `error` | Error if symbol not found |

---

## 💬 Just the Essentials

* **What it is:** Minimum distance from current price for SL/TP orders (in points).
* **Why you need it:** Validate SL/TP distances, avoid order rejections, calculate valid levels.
* **Sanity check:** Returns 0 for symbols with no restriction. Measured in points, not pips!

---

## 🎯 When to Use

✅ **Before setting SL/TP** - Validate stop distances

✅ **Order validation** - Check if SL/TP is valid

✅ **Strategy design** - Know minimum stop requirements

✅ **Error prevention** - Avoid "stops too close" errors

---

## 🔗 Usage Examples

### 1) Basic usage - check stop level

```go
stopLevel, err := sugar.GetMinStopLevel("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if stopLevel == 0 {
    fmt.Println("No stop level restriction")
} else {
    fmt.Printf("Minimum stop level: %d points\n", stopLevel)

    // Convert to pips (for 5-digit broker)
    stopLevelPips := float64(stopLevel) / 10.0
    fmt.Printf("                   %.1f pips\n", stopLevelPips)
}
```

---

### 2) Validate SL distance before placing order

```go
symbol := "EURUSD"
volume := 0.1
stopLossPips := 20.0

// Get minimum stop level
minStopLevel, _ := sugar.GetMinStopLevel(symbol)

// Convert our SL pips to points (5-digit broker)
ourStopPoints := stopLossPips * 10

if ourStopPoints < float64(minStopLevel) {
    fmt.Printf("❌ SL too close!\n")
    fmt.Printf("   Requested: %.0f points (%.1f pips)\n",
        ourStopPoints, stopLossPips)
    fmt.Printf("   Minimum:   %d points (%.1f pips)\n",
        minStopLevel, float64(minStopLevel)/10.0)

    // Adjust to minimum
    stopLossPips = float64(minStopLevel) / 10.0
    fmt.Printf("   Adjusted to: %.1f pips\n", stopLossPips)
}

fmt.Printf("✅ SL distance valid: %.1f pips\n", stopLossPips)
```

---

### 3) Calculate valid SL/TP prices

```go
func CalculateValidSLTP(sugar *mt5.MT5Sugar, symbol string, direction string, slPips, tpPips float64) (float64, float64, error) {
    // Get current prices and stop level
    info, err := sugar.GetSymbolInfo(symbol)
    if err != nil {
        return 0, 0, err
    }

    minStopLevel, _ := sugar.GetMinStopLevel(symbol)

    // Convert pips to points (5-digit)
    slPoints := slPips * 10
    tpPoints := tpPips * 10

    // Validate against minimum
    if slPoints < float64(minStopLevel) {
        slPoints = float64(minStopLevel)
    }
    if tpPoints < float64(minStopLevel) {
        tpPoints = float64(minStopLevel)
    }

    // Convert to price distance
    slDistance := slPoints * info.Point
    tpDistance := tpPoints * info.Point

    // Calculate SL/TP prices
    var sl, tp float64

    if direction == "BUY" {
        sl = info.Bid - slDistance
        tp = info.Bid + tpDistance
    } else { // SELL
        sl = info.Ask + slDistance
        tp = info.Ask - tpDistance
    }

    return sl, tp, nil
}

// Usage:
sl, tp, _ := CalculateValidSLTP(sugar, "EURUSD", "BUY", 20, 40)
fmt.Printf("Valid SL: %.5f\n", sl)
fmt.Printf("Valid TP: %.5f\n", tp)
```

---

### 4) Compare stop levels across symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"}

fmt.Println("Stop Level Comparison:")
fmt.Println("─────────────────────────────────────────")

for _, symbol := range symbols {
    stopLevel, err := sugar.GetMinStopLevel(symbol)
    if err != nil {
        fmt.Printf("%s: Error - %v\n", symbol, err)
        continue
    }

    stopLevelPips := float64(stopLevel) / 10.0

    status := "✅"
    if stopLevel > 50 {
        status = "⚠️"
    }

    fmt.Printf("%s %-8s: %3d points (%.1f pips)\n",
        status, symbol, stopLevel, stopLevelPips)
}
```

---

### 5) Check if scalping is allowed

```go
func CanScalp(sugar *mt5.MT5Sugar, symbol string, maxStopPips float64) (bool, string) {
    stopLevel, err := sugar.GetMinStopLevel(symbol)
    if err != nil {
        return false, fmt.Sprintf("error: %v", err)
    }

    stopLevelPips := float64(stopLevel) / 10.0

    if stopLevelPips > maxStopPips {
        return false, fmt.Sprintf("min stop %.1f pips exceeds max %.1f pips",
            stopLevelPips, maxStopPips)
    }

    return true, fmt.Sprintf("can use stops as tight as %.1f pips", stopLevelPips)
}

// Usage:
canScalp, message := CanScalp(sugar, "EURUSD", 10.0)
if canScalp {
    fmt.Printf("✅ %s is scalping-friendly: %s\n", "EURUSD", message)
} else {
    fmt.Printf("❌ %s not suitable for scalping: %s\n", "EURUSD", message)
}
```

---

### 6) SL/TP distance validator

```go
type StopValidator struct {
    Symbol       string
    MinStopLevel int64
    Point        float64
}

func NewStopValidator(sugar *mt5.MT5Sugar, symbol string) (*StopValidator, error) {
    minStop, err := sugar.GetMinStopLevel(symbol)
    if err != nil {
        return nil, err
    }

    info, err := sugar.GetSymbolInfo(symbol)
    if err != nil {
        return nil, err
    }

    return &StopValidator{
        Symbol:       symbol,
        MinStopLevel: minStop,
        Point:        info.Point,
    }, nil
}

func (sv *StopValidator) ValidateSLDistance(currentPrice, slPrice float64) (bool, string) {
    distance := math.Abs(currentPrice - slPrice)
    distancePoints := distance / sv.Point

    minDistance := float64(sv.MinStopLevel) * sv.Point

    if distance < minDistance {
        return false, fmt.Sprintf("SL too close: %.0f points (min: %d)",
            distancePoints, sv.MinStopLevel)
    }

    return true, fmt.Sprintf("valid: %.0f points", distancePoints)
}

func (sv *StopValidator) AdjustSLToMin(currentPrice float64, direction string) float64 {
    minDistance := float64(sv.MinStopLevel) * sv.Point

    if direction == "BUY" {
        return currentPrice - minDistance
    } else { // SELL
        return currentPrice + minDistance
    }
}

// Usage:
validator, _ := NewStopValidator(sugar, "EURUSD")

// Validate specific SL
currentPrice := 1.10000
proposedSL := 1.09950

valid, message := validator.ValidateSLDistance(currentPrice, proposedSL)
if valid {
    fmt.Printf("✅ %s\n", message)
} else {
    fmt.Printf("❌ %s\n", message)

    // Get minimum valid SL
    minSL := validator.AdjustSLToMin(currentPrice, "BUY")
    fmt.Printf("   Minimum SL: %.5f\n", minSL)
}
```

---

### 7) Display stop level in multiple formats

```go
stopLevel, _ := sugar.GetMinStopLevel("EURUSD")
info, _ := sugar.GetSymbolInfo("EURUSD")

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      STOP LEVEL INFORMATION           ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Symbol:        %s\n\n", "EURUSD")

fmt.Printf("Minimum Stop Level:\n")
fmt.Printf("  Points:      %d\n", stopLevel)
fmt.Printf("  Pips:        %.1f (5-digit)\n", float64(stopLevel)/10.0)
fmt.Printf("  Price:       %.5f\n", float64(stopLevel)*info.Point)

// Calculate examples
bid := info.Bid
minSLBuy := bid - (float64(stopLevel) * info.Point)
minTPBuy := bid + (float64(stopLevel) * info.Point)

fmt.Printf("\nExample (BUY at %.5f):\n", bid)
fmt.Printf("  Min SL:      %.5f\n", minSLBuy)
fmt.Printf("  Min TP:      %.5f\n", minTPBuy)
```

---

### 8) Strategy compatibility checker

```go
type TradingStrategy struct {
    Name           string
    TypicalSLPips  float64
    TypicalTPPips  float64
}

func (ts *TradingStrategy) IsCompatible(sugar *mt5.MT5Sugar, symbol string) (bool, string) {
    stopLevel, err := sugar.GetMinStopLevel(symbol)
    if err != nil {
        return false, fmt.Sprintf("error: %v", err)
    }

    minStopPips := float64(stopLevel) / 10.0

    if ts.TypicalSLPips < minStopPips {
        return false, fmt.Sprintf("SL %.1f pips < min %.1f pips",
            ts.TypicalSLPips, minStopPips)
    }

    if ts.TypicalTPPips < minStopPips {
        return false, fmt.Sprintf("TP %.1f pips < min %.1f pips",
            ts.TypicalTPPips, minStopPips)
    }

    return true, "compatible"
}

// Usage:
scalpStrategy := &TradingStrategy{
    Name:          "Scalping",
    TypicalSLPips: 5.0,
    TypicalTPPips: 10.0,
}

swingStrategy := &TradingStrategy{
    Name:          "Swing Trading",
    TypicalSLPips: 50.0,
    TypicalTPPips: 100.0,
}

symbols := []string{"EURUSD", "GBPUSD", "XAUUSD"}

for _, symbol := range symbols {
    compatible, message := scalpStrategy.IsCompatible(sugar, symbol)

    if compatible {
        fmt.Printf("✅ %s: %s compatible with %s\n",
            symbol, scalpStrategy.Name, message)
    } else {
        fmt.Printf("❌ %s: %s not compatible - %s\n",
            symbol, scalpStrategy.Name, message)
    }
}
```

---

### 9) Pending order distance calculator

```go
func CalculatePendingOrderDistance(sugar *mt5.MT5Sugar, symbol string) {
    stopLevel, _ := sugar.GetMinStopLevel(symbol)
    info, _ := sugar.GetSymbolInfo(symbol)

    minDistance := float64(stopLevel) * info.Point

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║   PENDING ORDER DISTANCE RULES        ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Symbol:              %s\n", symbol)
    fmt.Printf("Current Bid:         %.5f\n", info.Bid)
    fmt.Printf("Current Ask:         %.5f\n\n", info.Ask)

    fmt.Printf("Minimum Distance:    %d points (%.5f)\n\n",
        stopLevel, minDistance)

    fmt.Println("Pending Order Limits:")
    fmt.Println("─────────────────────────────────────────")

    // Buy Stop (above Ask)
    buyStopMin := info.Ask + minDistance
    fmt.Printf("Buy Stop:   ≥ %.5f\n", buyStopMin)

    // Buy Limit (below Bid)
    buyLimitMax := info.Bid - minDistance
    fmt.Printf("Buy Limit:  ≤ %.5f\n", buyLimitMax)

    // Sell Stop (below Bid)
    sellStopMax := info.Bid - minDistance
    fmt.Printf("Sell Stop:  ≤ %.5f\n", sellStopMax)

    // Sell Limit (above Ask)
    sellLimitMin := info.Ask + minDistance
    fmt.Printf("Sell Limit: ≥ %.5f\n", sellLimitMin)
}

// Usage:
CalculatePendingOrderDistance(sugar, "EURUSD")
```

---

### 10) Advanced stop level manager

```go
type StopLevelManager struct {
    cache     map[string]int64
    cacheTime map[string]time.Time
    cacheTTL  time.Duration
}

func NewStopLevelManager() *StopLevelManager {
    return &StopLevelManager{
        cache:     make(map[string]int64),
        cacheTime: make(map[string]time.Time),
        cacheTTL:  10 * time.Minute, // Stop levels rarely change
    }
}

func (slm *StopLevelManager) GetStopLevel(sugar *mt5.MT5Sugar, symbol string) (int64, error) {
    // Check cache
    if cachedTime, exists := slm.cacheTime[symbol]; exists {
        if time.Since(cachedTime) < slm.cacheTTL {
            return slm.cache[symbol], nil
        }
    }

    // Fetch from MT5
    stopLevel, err := sugar.GetMinStopLevel(symbol)
    if err != nil {
        return 0, err
    }

    // Update cache
    slm.cache[symbol] = stopLevel
    slm.cacheTime[symbol] = time.Now()

    return stopLevel, nil
}

func (slm *StopLevelManager) GetStopLevelPips(sugar *mt5.MT5Sugar, symbol string) (float64, error) {
    stopLevel, err := slm.GetStopLevel(sugar, symbol)
    if err != nil {
        return 0, err
    }

    // Convert to pips (assumes 5-digit)
    return float64(stopLevel) / 10.0, nil
}

func (slm *StopLevelManager) ValidateStopDistance(
    sugar *mt5.MT5Sugar,
    symbol string,
    distancePips float64,
) (bool, error) {
    minStopPips, err := slm.GetStopLevelPips(sugar, symbol)
    if err != nil {
        return false, err
    }

    return distancePips >= minStopPips, nil
}

func (slm *StopLevelManager) GenerateReport(sugar *mt5.MT5Sugar, symbols []string) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      STOP LEVEL REPORT                ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    for _, symbol := range symbols {
        stopPips, err := slm.GetStopLevelPips(sugar, symbol)
        if err != nil {
            fmt.Printf("%-8s: Error - %v\n", symbol, err)
            continue
        }

        rating := "✅ Low"
        if stopPips > 50 {
            rating = "⚠️  High"
        } else if stopPips > 20 {
            rating = "🟡 Medium"
        }

        fmt.Printf("%-8s: %5.1f pips  %s\n", symbol, stopPips, rating)
    }
}

// Usage:
manager := NewStopLevelManager()

// Get stop level (cached)
stopLevel, _ := manager.GetStopLevel(sugar, "EURUSD")
fmt.Printf("EURUSD stop level: %d points\n", stopLevel)

// Validate distance
valid, _ := manager.ValidateStopDistance(sugar, "EURUSD", 20.0)
if valid {
    fmt.Println("✅ 20 pip stop is valid")
}

// Generate report
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"}
manager.GenerateReport(sugar, symbols)
```

---

## 🔗 Related Methods

**🍬 Symbol information:**

* `GetSymbolInfo()` - Gets StopLevel as part of full info
* `GetSymbolDigits()` - Get decimal precision

**🍬 Price methods:**

* `GetBid()` / `GetAsk()` - Current prices for SL/TP calculation

---

## ⚠️ Common Pitfalls

### 1) Confusing points and pips

```go
// ❌ WRONG - stopLevel is in POINTS, not pips
stopLevel, _ := sugar.GetMinStopLevel("EURUSD")
// For 5-digit: 20 points = 2 pips!

// ✅ CORRECT - convert to pips
stopLevelPips := float64(stopLevel) / 10.0 // For 5-digit broker
```

### 2) Not validating SL/TP before placing order

```go
// ❌ WRONG - might be rejected
sugar.BuyMarketWithSLTP("EURUSD", 0.1, sl, tp)

// ✅ CORRECT - validate first
stopLevel, _ := sugar.GetMinStopLevel("EURUSD")
// Check SL/TP distances are >= stopLevel points
```

### 3) Assuming zero means no restriction

```go
// ❌ WRONG - zero might mean error
stopLevel, _ := sugar.GetMinStopLevel("EURUSD")
if stopLevel == 0 {
    // Could be error OR no restriction
}

// ✅ CORRECT - check error
stopLevel, err := sugar.GetMinStopLevel("EURUSD")
if err != nil {
    // Error
} else if stopLevel == 0 {
    // No restriction
}
```

---

## 💎 Pro Tips

1. **Points not pips** - Return value is in POINTS (÷10 for pips on 5-digit)

2. **Cache safe** - Stop levels rarely change, safe to cache

3. **Zero valid** - Zero means no stop level restriction

4. **Both SL and TP** - Applies to both SL and TP distances

5. **Pending orders** - Also applies to pending order distance from current price

---

## 📊 Stop Level in Different Units

```
For EURUSD (5-digit broker):
StopLevel = 20 points

Conversions:
- Points:  20
- Pips:    2.0 (÷10)
- Price:   0.00020 (×Point size)

For USDJPY (3-digit broker):
StopLevel = 20 points

Conversions:
- Points:  20
- Pips:    2.0 (÷10)
- Price:   0.020 (×Point size)
```

---

**See also:** [`GetSymbolInfo.md`](GetSymbolInfo.md), [`GetSymbolDigits.md`](GetSymbolDigits.md)
