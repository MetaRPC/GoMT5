# ✅🔍 Is Symbol Available (`IsSymbolAvailable`)

> **Sugar method:** Checks if a symbol exists and is available for trading.

**API Information:**

* **Method:** `sugar.IsSymbolAvailable(symbol)`
* **Timeout:** 3 seconds
* **Returns:** Boolean availability status

---

## 📋 Method Signature

```go
func (s *MT5Sugar) IsSymbolAvailable(symbol string) (bool, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Symbol name to check (e.g., "EURUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `bool` | boolean | `true` if available, `false` otherwise |
| `error` | `error` | Error if query failed |

---

## 💬 Just the Essentials

* **What it is:** Check if a symbol exists and can be traded.
* **Why you need it:** Validation before trading, symbol verification, error prevention.
* **Sanity check:** Returns `false` for non-existent or disabled symbols (not an error).

---

## 🎯 When to Use

✅ **Before trading** - Validate symbol exists

✅ **User input** - Verify user-entered symbols

✅ **Symbol switching** - Check availability when changing symbols

✅ **Error prevention** - Avoid errors from invalid symbols

---

## 🔗 Usage Examples

### 1) Basic usage - check if symbol available

```go
available, err := sugar.IsSymbolAvailable("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if available {
    fmt.Println("✅ EURUSD is available for trading")
} else {
    fmt.Println("❌ EURUSD is not available")
}
```

---

### 2) Validate before opening position

```go
symbol := "EURUSD"
volume := 0.1

// Check symbol first
available, _ := sugar.IsSymbolAvailable(symbol)
if !available {
    fmt.Printf("❌ Cannot trade %s - symbol not available\n", symbol)
    return
}

// Safe to trade
fmt.Printf("✅ %s available - opening position...\n", symbol)
ticket, err := sugar.BuyMarket(symbol, volume)
if err != nil {
    fmt.Printf("Error: %v\n", err)
} else {
    fmt.Printf("Position #%d opened\n", ticket)
}
```

---

### 3) Validate user input

```go
func ValidateUserSymbol(sugar *mt5.MT5Sugar, userInput string) (string, error) {
    // Normalize input
    symbol := strings.ToUpper(strings.TrimSpace(userInput))

    // Check availability
    available, err := sugar.IsSymbolAvailable(symbol)
    if err != nil {
        return "", fmt.Errorf("validation failed: %w", err)
    }

    if !available {
        return "", fmt.Errorf("symbol %s is not available", symbol)
    }

    return symbol, nil
}

// Usage:
userInput := "eurusd"
validSymbol, err := ValidateUserSymbol(sugar, userInput)
if err != nil {
    fmt.Printf("❌ Invalid symbol: %v\n", err)
} else {
    fmt.Printf("✅ Valid symbol: %s\n", validSymbol)
}
```

---

### 4) Check multiple symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "INVALID"}

fmt.Println("Symbol Availability Check:")
fmt.Println("─────────────────────────────────────────")

availableCount := 0

for _, symbol := range symbols {
    available, _ := sugar.IsSymbolAvailable(symbol)

    if available {
        fmt.Printf("✅ %s\n", symbol)
        availableCount++
    } else {
        fmt.Printf("❌ %s\n", symbol)
    }
}

fmt.Printf("\n%d/%d symbols available\n", availableCount, len(symbols))
```

---

### 5) Build available symbols list

```go
func GetAvailableSymbols(sugar *mt5.MT5Sugar, candidates []string) []string {
    available := []string{}

    for _, symbol := range candidates {
        isAvailable, _ := sugar.IsSymbolAvailable(symbol)
        if isAvailable {
            available = append(available, symbol)
        }
    }

    return available
}

// Usage:
watchlist := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD", "BTCUSD"}
available := GetAvailableSymbols(sugar, watchlist)

fmt.Printf("Available from watchlist: %d/%d\n", len(available), len(watchlist))
for _, symbol := range available {
    fmt.Println("  " + symbol)
}
```

---

### 6) Symbol switcher with validation

```go
func SwitchToSymbol(sugar *mt5.MT5Sugar, newSymbol string) error {
    fmt.Printf("Switching to %s...\n", newSymbol)

    // Check availability
    available, err := sugar.IsSymbolAvailable(newSymbol)
    if err != nil {
        return fmt.Errorf("availability check failed: %w", err)
    }

    if !available {
        return fmt.Errorf("symbol %s is not available", newSymbol)
    }

    // Symbol is available, proceed
    fmt.Printf("✅ %s is available\n", newSymbol)

    // Get symbol info
    info, err := sugar.GetSymbolInfo(newSymbol)
    if err != nil {
        return fmt.Errorf("failed to get symbol info: %w", err)
    }

    fmt.Printf("   Spread: %d points\n", info.Spread)
    fmt.Printf("   Min lot: %.2f\n", info.VolumeMin)

    return nil
}

// Usage:
err := SwitchToSymbol(sugar, "GBPUSD")
if err != nil {
    fmt.Printf("Switch failed: %v\n", err)
}
```

---

### 7) Safe symbol operation wrapper

```go
func WithValidSymbol(sugar *mt5.MT5Sugar, symbol string, operation func(string) error) error {
    // Validate symbol first
    available, err := sugar.IsSymbolAvailable(symbol)
    if err != nil {
        return fmt.Errorf("validation error: %w", err)
    }

    if !available {
        return fmt.Errorf("symbol %s not available", symbol)
    }

    // Execute operation
    return operation(symbol)
}

// Usage:
err := WithValidSymbol(sugar, "EURUSD", func(symbol string) error {
    // This only runs if symbol is valid
    ticket, err := sugar.BuyMarket(symbol, 0.1)
    if err != nil {
        return err
    }

    fmt.Printf("Opened position #%d on %s\n", ticket, symbol)
    return nil
})

if err != nil {
    fmt.Printf("Operation failed: %v\n", err)
}
```

---

### 8) Symbol availability monitor

```go
func MonitorSymbolAvailability(sugar *mt5.MT5Sugar, symbol string, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    wasAvailable := false

    for range ticker.C {
        available, _ := sugar.IsSymbolAvailable(symbol)

        if available != wasAvailable {
            // Status changed
            if available {
                fmt.Printf("[%s] %s: ❌ → ✅ Now available\n",
                    time.Now().Format("15:04:05"), symbol)
            } else {
                fmt.Printf("[%s] %s: ✅ → ❌ No longer available\n",
                    time.Now().Format("15:04:05"), symbol)
            }

            wasAvailable = available
        }
    }
}

// Usage: Monitor EURUSD every 30 seconds
go MonitorSymbolAvailability(sugar, "EURUSD", 30*time.Second)
```

---

### 9) Batch symbol validator

```go
type SymbolValidationResult struct {
    Symbol    string
    Available bool
    Error     error
}

func ValidateSymbolsBatch(sugar *mt5.MT5Sugar, symbols []string) []SymbolValidationResult {
    results := make([]SymbolValidationResult, len(symbols))

    for i, symbol := range symbols {
        available, err := sugar.IsSymbolAvailable(symbol)

        results[i] = SymbolValidationResult{
            Symbol:    symbol,
            Available: available,
            Error:     err,
        }
    }

    return results
}

func ShowValidationResults(results []SymbolValidationResult) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      SYMBOL VALIDATION RESULTS        ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    availableCount := 0
    errorCount := 0

    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("⚠️  %-8s: Error - %v\n", result.Symbol, result.Error)
            errorCount++
        } else if result.Available {
            fmt.Printf("✅ %-8s: Available\n", result.Symbol)
            availableCount++
        } else {
            fmt.Printf("❌ %-8s: Not available\n", result.Symbol)
        }
    }

    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Available:  %d\n", availableCount)
    fmt.Printf("Unavailable: %d\n", len(results)-availableCount-errorCount)
    fmt.Printf("Errors:     %d\n", errorCount)
}

// Usage:
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "INVALID", "BTCUSD"}
results := ValidateSymbolsBatch(sugar, symbols)
ShowValidationResults(results)
```

---

### 10) Advanced symbol availability manager

```go
type SymbolAvailabilityManager struct {
    cache      map[string]bool
    cacheTime  map[string]time.Time
    cacheTTL   time.Duration
}

func NewSymbolAvailabilityManager(ttl time.Duration) *SymbolAvailabilityManager {
    return &SymbolAvailabilityManager{
        cache:     make(map[string]bool),
        cacheTime: make(map[string]time.Time),
        cacheTTL:  ttl,
    }
}

func (sam *SymbolAvailabilityManager) IsAvailable(sugar *mt5.MT5Sugar, symbol string) (bool, error) {
    // Check cache
    if cachedTime, exists := sam.cacheTime[symbol]; exists {
        if time.Since(cachedTime) < sam.cacheTTL {
            // Cache hit
            return sam.cache[symbol], nil
        }
    }

    // Cache miss or expired - query MT5
    available, err := sugar.IsSymbolAvailable(symbol)
    if err != nil {
        return false, err
    }

    // Update cache
    sam.cache[symbol] = available
    sam.cacheTime[symbol] = time.Now()

    return available, nil
}

func (sam *SymbolAvailabilityManager) ClearCache() {
    sam.cache = make(map[string]bool)
    sam.cacheTime = make(map[string]time.Time)
}

func (sam *SymbolAvailabilityManager) GetCacheStats() (int, int) {
    total := len(sam.cache)
    expired := 0

    for symbol, cachedTime := range sam.cacheTime {
        if time.Since(cachedTime) >= sam.cacheTTL {
            expired++
            // Clean up expired
            delete(sam.cache, symbol)
            delete(sam.cacheTime, symbol)
        }
    }

    return total - expired, expired
}

func (sam *SymbolAvailabilityManager) ValidateAndExecute(
    sugar *mt5.MT5Sugar,
    symbol string,
    operation func() error,
) error {
    available, err := sam.IsAvailable(sugar, symbol)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    if !available {
        return fmt.Errorf("symbol %s not available", symbol)
    }

    return operation()
}

// Usage:
manager := NewSymbolAvailabilityManager(5 * time.Minute)

// Check with caching
available, _ := manager.IsAvailable(sugar, "EURUSD")
if available {
    fmt.Println("✅ EURUSD available (may be cached)")
}

// Validate and execute
err := manager.ValidateAndExecute(sugar, "EURUSD", func() error {
    ticket, err := sugar.BuyMarket("EURUSD", 0.1)
    if err != nil {
        return err
    }
    fmt.Printf("Opened position #%d\n", ticket)
    return nil
})

if err != nil {
    fmt.Printf("Operation failed: %v\n", err)
}

// Cache stats
valid, expired := manager.GetCacheStats()
fmt.Printf("Cache: %d valid, %d expired\n", valid, expired)
```

---

## 🔗 Related Methods

**🍬 Symbol information:**

* `GetSymbolInfo()` - Get full symbol details
* `GetAllSymbols()` - List all symbols

**🍬 Symbol properties:**

* `GetMinStopLevel()` - Get stop level
* `GetSymbolDigits()` - Get decimal precision

---

## ⚠️ Common Pitfalls

### 1) Not checking the result

```go
// ❌ WRONG - ignoring return value
sugar.IsSymbolAvailable("EURUSD")
sugar.BuyMarket("EURUSD", 0.1) // Might fail!

// ✅ CORRECT - check before trading
available, _ := sugar.IsSymbolAvailable("EURUSD")
if available {
    sugar.BuyMarket("EURUSD", 0.1)
}
```

### 2) Assuming false means error

```go
// ❌ WRONG - false is not an error
available, _ := sugar.IsSymbolAvailable("INVALID")
if !available {
    // This is expected for invalid symbols!
}

// ✅ CORRECT - check error separately
available, err := sugar.IsSymbolAvailable("INVALID")
if err != nil {
    // Actual error
} else if !available {
    // Symbol just doesn't exist or isn't available
}
```

### 3) Case sensitivity

```go
// ❌ WRONG - case matters
available, _ := sugar.IsSymbolAvailable("eurusd")

// ✅ CORRECT - use exact symbol name
available, _ := sugar.IsSymbolAvailable("EURUSD")
```

---

## 💎 Pro Tips

1. **Fast check** - Lightweight operation, safe to call frequently

2. **Cache results** - Symbol availability rarely changes, can cache

3. **Validate early** - Check before expensive operations

4. **Case matters** - Symbol names are case-sensitive

5. **Not an error** - `false` means unavailable (not necessarily an error)

---

## 📊 Availability States

```
Symbol states:
✅ Available (true)   - Symbol exists and can be traded
❌ Unavailable (false) - Symbol doesn't exist OR
                       - Symbol exists but trading disabled OR
                       - Symbol not synchronized

Common reasons for unavailable:
- Symbol doesn't exist on this broker
- Symbol is disabled by broker
- Symbol not in user's watchlist
- Market is closed for this symbol
- Insufficient permissions
```

---

**See also:** [`GetSymbolInfo.md`](GetSymbolInfo.md), [`GetAllSymbols.md`](GetAllSymbols.md)
