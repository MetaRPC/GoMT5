# 🔢📊 Get Symbol Digits (`GetSymbolDigits`)

> **Sugar method:** Returns the number of decimal places (precision) for a symbol.

**API Information:**

* **Method:** `sugar.GetSymbolDigits(symbol)`
* **Timeout:** 3 seconds
* **Returns:** Number of decimal digits

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetSymbolDigits(symbol string) (int32, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `int32` | integer | Number of decimal places |
| `error` | `error` | Error if symbol not found |

---

## 💬 Just the Essentials

* **What it is:** Get the decimal precision used for price quotes.
* **Why you need it:** Price formatting, calculations, point size determination.
* **Sanity check:** EURUSD = 5 digits (1.10500), USDJPY = 3 digits (110.500).

---

## 🎯 When to Use

✅ **Price formatting** - Display prices correctly

✅ **Calculations** - Determine point size

✅ **Validation** - Check price precision

✅ **Rounding** - Round prices to valid precision

---

## 🔗 Usage Examples

### 1) Basic usage - get digits

```go
digits, err := sugar.GetSymbolDigits("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("EURUSD uses %d decimal places\n", digits)

// Calculate point size
pointSize := 1.0
for i := int32(0); i < digits; i++ {
    pointSize /= 10.0
}

fmt.Printf("Point size: %.5f\n", pointSize)
```

---

### 2) Format price with correct precision

```go
func FormatPrice(sugar *mt5.MT5Sugar, symbol string, price float64) string {
    digits, err := sugar.GetSymbolDigits(symbol)
    if err != nil {
        return fmt.Sprintf("%.5f", price) // Default
    }

    format := fmt.Sprintf("%%.%df", digits)
    return fmt.Sprintf(format, price)
}

// Usage:
price := 1.105003
formatted := FormatPrice(sugar, "EURUSD", price)
fmt.Printf("Price: %s\n", formatted) // "1.10500"
```

---

### 3) Round price to valid precision

```go
func RoundToDigits(sugar *mt5.MT5Sugar, symbol string, price float64) (float64, error) {
    digits, err := sugar.GetSymbolDigits(symbol)
    if err != nil {
        return 0, err
    }

    multiplier := math.Pow(10, float64(digits))
    rounded := math.Round(price*multiplier) / multiplier

    return rounded, nil
}

// Usage:
price := 1.1050033333
rounded, _ := RoundToDigits(sugar, "EURUSD", price)
fmt.Printf("Original: %.10f\n", price)
fmt.Printf("Rounded:  %.5f\n", rounded) // 1.10500
```

---

### 4) Calculate point size from digits

```go
digits, _ := sugar.GetSymbolDigits("EURUSD")

pointSize := math.Pow(10, -float64(digits))

fmt.Printf("Symbol:      EURUSD\n")
fmt.Printf("Digits:      %d\n", digits)
fmt.Printf("Point size:  %.5f\n", pointSize)

// Calculate pip size (usually 10x point for 5-digit)
pipSize := pointSize * 10
fmt.Printf("Pip size:    %.5f\n", pipSize)
```

---

### 5) Compare precision across symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD", "BTCUSD"}

fmt.Println("Symbol Precision Comparison:")
fmt.Println("─────────────────────────────────────────")

for _, symbol := range symbols {
    digits, err := sugar.GetSymbolDigits(symbol)
    if err != nil {
        fmt.Printf("%-8s: Error - %v\n", symbol, err)
        continue
    }

    pointSize := math.Pow(10, -float64(digits))

    fmt.Printf("%-8s: %d digits (point: %.5f)\n",
        symbol, digits, pointSize)
}
```

---

### 6) Validate price precision

```go
func ValidatePricePrecision(sugar *mt5.MT5Sugar, symbol string, price float64) (bool, string) {
    digits, err := sugar.GetSymbolDigits(symbol)
    if err != nil {
        return false, fmt.Sprintf("error: %v", err)
    }

    // Convert to string to count decimals
    priceStr := fmt.Sprintf("%.10f", price)
    parts := strings.Split(priceStr, ".")

    if len(parts) != 2 {
        return false, "invalid price format"
    }

    // Count actual decimal places (excluding trailing zeros)
    decimals := strings.TrimRight(parts[1], "0")
    actualDigits := len(decimals)

    if actualDigits > int(digits) {
        return false, fmt.Sprintf("price has %d decimals, max is %d",
            actualDigits, digits)
    }

    return true, "valid precision"
}

// Usage:
price := 1.10500
valid, message := ValidatePricePrecision(sugar, "EURUSD", price)
if valid {
    fmt.Printf("✅ Price %f: %s\n", price, message)
} else {
    fmt.Printf("❌ Price %f: %s\n", price, message)
}
```

---

### 7) Display price with symbol's precision

```go
func DisplayPriceQuote(sugar *mt5.MT5Sugar, symbol string) {
    digits, _ := sugar.GetSymbolDigits(symbol)
    bid, _ := sugar.GetBid(symbol)
    ask, _ := sugar.GetAsk(symbol)

    format := fmt.Sprintf("%%.%df", digits)

    fmt.Printf("Symbol: %s (%d digits)\n", symbol, digits)
    fmt.Printf("Bid:    %s\n", fmt.Sprintf(format, bid))
    fmt.Printf("Ask:    %s\n", fmt.Sprintf(format, ask))
}

// Usage:
DisplayPriceQuote(sugar, "EURUSD")
// Output:
// Symbol: EURUSD (5 digits)
// Bid:    1.10500
// Ask:    1.10502
```

---

### 8) Normalize price to symbol's precision

```go
type PriceNormalizer struct {
    Symbol     string
    Digits     int32
    Multiplier float64
}

func NewPriceNormalizer(sugar *mt5.MT5Sugar, symbol string) (*PriceNormalizer, error) {
    digits, err := sugar.GetSymbolDigits(symbol)
    if err != nil {
        return nil, err
    }

    return &PriceNormalizer{
        Symbol:     symbol,
        Digits:     digits,
        Multiplier: math.Pow(10, float64(digits)),
    }, nil
}

func (pn *PriceNormalizer) Normalize(price float64) float64 {
    return math.Round(price*pn.Multiplier) / pn.Multiplier
}

func (pn *PriceNormalizer) Format(price float64) string {
    normalized := pn.Normalize(price)
    format := fmt.Sprintf("%%.%df", pn.Digits)
    return fmt.Sprintf(format, normalized)
}

func (pn *PriceNormalizer) GetPointSize() float64 {
    return 1.0 / pn.Multiplier
}

// Usage:
normalizer, _ := NewPriceNormalizer(sugar, "EURUSD")

price := 1.1050033333
normalized := normalizer.Normalize(price)
formatted := normalizer.Format(price)

fmt.Printf("Original:   %.10f\n", price)
fmt.Printf("Normalized: %.5f\n", normalized)
fmt.Printf("Formatted:  %s\n", formatted)
fmt.Printf("Point size: %.5f\n", normalizer.GetPointSize())
```

---

### 9) Calculate price increments

```go
func CalculatePriceIncrements(sugar *mt5.MT5Sugar, symbol string, basePrice float64, steps int) {
    digits, _ := sugar.GetSymbolDigits(symbol)
    pointSize := math.Pow(10, -float64(digits))

    fmt.Printf("Symbol: %s (%d digits)\n", symbol, digits)
    fmt.Printf("Base price: %.*f\n", digits, basePrice)
    fmt.Printf("Point size: %.5f\n\n", pointSize)

    fmt.Println("Price increments:")
    for i := -steps; i <= steps; i++ {
        price := basePrice + (float64(i) * pointSize)

        marker := " "
        if i == 0 {
            marker = "→"
        }

        fmt.Printf("%s %3d points: %.*f\n",
            marker, i, digits, price)
    }
}

// Usage:
CalculatePriceIncrements(sugar, "EURUSD", 1.10500, 5)
```

---

### 10) Advanced symbol precision manager

```go
type SymbolPrecisionManager struct {
    cache      map[string]int32
    formatters map[string]string
}

func NewSymbolPrecisionManager() *SymbolPrecisionManager {
    return &SymbolPrecisionManager{
        cache:      make(map[string]int32),
        formatters: make(map[string]string),
    }
}

func (spm *SymbolPrecisionManager) GetDigits(sugar *mt5.MT5Sugar, symbol string) (int32, error) {
    // Check cache
    if digits, exists := spm.cache[symbol]; exists {
        return digits, nil
    }

    // Fetch from MT5
    digits, err := sugar.GetSymbolDigits(symbol)
    if err != nil {
        return 0, err
    }

    // Cache
    spm.cache[symbol] = digits
    spm.formatters[symbol] = fmt.Sprintf("%%.%df", digits)

    return digits, nil
}

func (spm *SymbolPrecisionManager) FormatPrice(sugar *mt5.MT5Sugar, symbol string, price float64) (string, error) {
    digits, err := spm.GetDigits(sugar, symbol)
    if err != nil {
        return "", err
    }

    formatter := spm.formatters[symbol]
    return fmt.Sprintf(formatter, price), nil
}

func (spm *SymbolPrecisionManager) NormalizePrice(sugar *mt5.MT5Sugar, symbol string, price float64) (float64, error) {
    digits, err := spm.GetDigits(sugar, symbol)
    if err != nil {
        return 0, err
    }

    multiplier := math.Pow(10, float64(digits))
    return math.Round(price*multiplier) / multiplier, nil
}

func (spm *SymbolPrecisionManager) GetPointSize(sugar *mt5.MT5Sugar, symbol string) (float64, error) {
    digits, err := spm.GetDigits(sugar, symbol)
    if err != nil {
        return 0, err
    }

    return math.Pow(10, -float64(digits)), nil
}

func (spm *SymbolPrecisionManager) CalculatePipValue(
    sugar *mt5.MT5Sugar,
    symbol string,
    volume float64,
) (float64, error) {
    // Get symbol info for contract size
    info, err := sugar.GetSymbolInfo(symbol)
    if err != nil {
        return 0, err
    }

    // Get point size
    pointSize, err := spm.GetPointSize(sugar, symbol)
    if err != nil {
        return 0, err
    }

    // Pip = 10 points (for 5-digit broker)
    pipSize := pointSize * 10

    // Pip value = pip size × contract size × volume
    pipValue := pipSize * info.ContractSize * volume

    return pipValue, nil
}

func (spm *SymbolPrecisionManager) GenerateReport(sugar *mt5.MT5Sugar, symbols []string) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      PRECISION ANALYSIS REPORT        ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    for _, symbol := range symbols {
        digits, err := spm.GetDigits(sugar, symbol)
        if err != nil {
            fmt.Printf("%-8s: Error - %v\n", symbol, err)
            continue
        }

        pointSize, _ := spm.GetPointSize(sugar, symbol)
        pipValue, _ := spm.CalculatePipValue(sugar, symbol, 1.0)

        fmt.Printf("\n%-8s:\n", symbol)
        fmt.Printf("  Digits:     %d\n", digits)
        fmt.Printf("  Point:      %.5f\n", pointSize)
        fmt.Printf("  Pip value:  $%.2f per lot\n", pipValue)

        // Show example price
        bid, _ := sugar.GetBid(symbol)
        formatted, _ := spm.FormatPrice(sugar, symbol, bid)
        fmt.Printf("  Example:    %s\n", formatted)
    }
}

// Usage:
manager := NewSymbolPrecisionManager()

// Format prices (cached for performance)
price1, _ := manager.FormatPrice(sugar, "EURUSD", 1.105003)
price2, _ := manager.FormatPrice(sugar, "USDJPY", 110.5033)

fmt.Printf("EURUSD: %s\n", price1) // 1.10500
fmt.Printf("USDJPY: %s\n", price2) // 110.503

// Normalize price
normalized, _ := manager.NormalizePrice(sugar, "EURUSD", 1.1050033)
fmt.Printf("Normalized: %.5f\n", normalized)

// Calculate pip value
pipValue, _ := manager.CalculatePipValue(sugar, "EURUSD", 1.0)
fmt.Printf("Pip value: $%.2f\n", pipValue)

// Generate report
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"}
manager.GenerateReport(sugar, symbols)
```

---

## 🔗 Related Methods

**🍬 Symbol information:**

* `GetSymbolInfo()` - Gets Digits as part of full info
* `GetMinStopLevel()` - Stop level in points

**🍬 Price methods:**

* `GetBid()` / `GetAsk()` - Get current prices

---

## ⚠️ Common Pitfalls

### 1) Confusing digits with decimals

```go
// ❌ WRONG - digits = total decimal places
digits := 5
// For 1.10500: digits = 5 (all decimals)

// ✅ CORRECT - understand it's decimal precision
digits := 5
pointSize := math.Pow(10, -float64(digits)) // 0.00001
```

### 2) Not rounding prices

```go
// ❌ WRONG - using raw float
price := 1.1050033333
sugar.BuyMarket("EURUSD", 0.1) // Might fail!

// ✅ CORRECT - round to valid precision
digits, _ := sugar.GetSymbolDigits("EURUSD")
multiplier := math.Pow(10, float64(digits))
price = math.Round(price*multiplier) / multiplier
```

### 3) Hardcoding precision

```go
// ❌ WRONG - assuming 5 digits
price := fmt.Sprintf("%.5f", value)

// ✅ CORRECT - use actual digits
digits, _ := sugar.GetSymbolDigits(symbol)
format := fmt.Sprintf("%%.%df", digits)
price := fmt.Sprintf(format, value)
```

---

## 💎 Pro Tips

1. **Cache digits** - Digits never change, safe to cache

2. **Point size** - Calculate as 10^(-digits)

3. **Formatting** - Use digits for display formatting

4. **Validation** - Round prices to symbol's precision

5. **Pip = 10 points** - For 5-digit brokers (EURUSD, GBPUSD)

---

## 📊 Common Digit Values

```
Standard Forex (5 digits):
EURUSD, GBPUSD, etc. = 5 digits
- Example: 1.10500
- Point: 0.00001
- Pip: 0.0001 (10 points)

Yen pairs (3 digits):
USDJPY, EURJPY, etc. = 3 digits
- Example: 110.500
- Point: 0.001
- Pip: 0.01 (10 points)

Gold (2 digits):
XAUUSD = 2 digits
- Example: 1850.50
- Point: 0.01
- Pip: 0.10 (10 points)

Crypto (varies):
BTCUSD = 2 digits
- Example: 45000.50
- Point: 0.01
```

---

**See also:** [`GetSymbolInfo.md`](GetSymbolInfo.md), [`GetMinStopLevel.md`](GetMinStopLevel.md)
