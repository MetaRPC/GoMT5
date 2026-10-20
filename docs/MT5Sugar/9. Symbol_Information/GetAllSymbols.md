# 📋🔍 Get All Symbols (`GetAllSymbols`)

> **Sugar method:** Returns a list of all available trading symbols on the broker.

**API Information:**

* **Method:** `sugar.GetAllSymbols()`
* **Timeout:** 60 seconds (for all pages together)
* **Returns:** Slice of symbol names

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetAllSymbols() ([]string, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters required |

| Output | Type | Description |
|--------|------|-------------|
| `[]string` | slice | List of all symbol names |
| `error` | `error` | `"GetAllSymbols failed (page N): ..."` if a page query failed |

---

## 💬 Just the Essentials

* **What it is:** Get a complete list of all symbols available on your broker.
* **Why you need it:** Symbol discovery, building watchlists, symbol validation.
* **Sanity check:** Returns symbol names only (not full info). Use GetSymbolInfo() for details.
* **How it works:** `SymbolParamsMany` answers one page at a time - a request without paging returns only the **first page** (50 names on the MetaQuotes-Demo test server). So the method reads pages of 1000 (`PageNumber` / `ItemsPerPage`) until it has `SymbolsTotal` names (or a page comes back empty); only the `Name` of each entry is kept. One 60-second timeout covers all pages (12,293 symbols on MetaQuotes-Demo took about 10 s).

```go
func (s *MT5Sugar) GetAllSymbols() ([]string, error) {
    ctx, cancel := context.WithTimeout(s.ctx, 60*time.Second)
    defer cancel()

    const perPage = int32(1000)
    var symbols []string
    for page := int32(1); ; page++ {
        pageNumber, itemsPerPage := page, perPage
        params, total, err := s.service.GetSymbolParamsMany(ctx, nil, nil, &pageNumber, &itemsPerPage)
        if err != nil {
            return nil, fmt.Errorf("GetAllSymbols failed (page %d): %w", page, err)
        }
        for _, p := range params {
            symbols = append(symbols, p.Name)
        }
        if len(params) == 0 || int32(len(symbols)) >= total {
            break
        }
    }

    return symbols, nil
}
```

---

## 🎯 When to Use

✅ **Symbol discovery** - Find available instruments

✅ **Watchlist building** - Create custom symbol lists

✅ **Symbol validation** - Check if symbol exists

✅ **Market scanner** - Iterate through all symbols

---

## 🔗 Usage Examples

### 1) Basic usage - list all symbols

```go
symbols, err := sugar.GetAllSymbols()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Total symbols: %d\n\n", len(symbols))

// Show first 20
for i := 0; i < 20 && i < len(symbols); i++ {
    fmt.Printf("%d. %s\n", i+1, symbols[i])
}

if len(symbols) > 20 {
    fmt.Printf("... and %d more\n", len(symbols)-20)
}
```

---

### 2) Filter by symbol pattern

```go
symbols, _ := sugar.GetAllSymbols()

// Find all USD pairs
usdPairs := []string{}
for _, symbol := range symbols {
    if strings.Contains(symbol, "USD") {
        usdPairs = append(usdPairs, symbol)
    }
}

fmt.Printf("USD pairs: %d\n", len(usdPairs))
for _, pair := range usdPairs {
    fmt.Println("  " + pair)
}
```

---

### 3) Group symbols by type

```go
symbols, _ := sugar.GetAllSymbols()

forex := []string{}
metals := []string{}
indices := []string{}
crypto := []string{}
other := []string{}

for _, symbol := range symbols {
    // Simple classification based on common patterns
    if strings.HasPrefix(symbol, "XAU") || strings.HasPrefix(symbol, "XAG") {
        metals = append(metals, symbol)
    } else if strings.Contains(symbol, "BTC") || strings.Contains(symbol, "ETH") {
        crypto = append(crypto, symbol)
    } else if strings.HasSuffix(symbol, "USD") || strings.HasSuffix(symbol, "EUR") {
        forex = append(forex, symbol)
    } else if len(symbol) <= 6 && !strings.Contains(symbol, ".") {
        forex = append(forex, symbol)
    } else {
        other = append(other, symbol)
    }
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      SYMBOLS BY CATEGORY              ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Forex:   %d symbols\n", len(forex))
fmt.Printf("Metals:  %d symbols\n", len(metals))
fmt.Printf("Crypto:  %d symbols\n", len(crypto))
fmt.Printf("Indices: %d symbols\n", len(indices))
fmt.Printf("Other:   %d symbols\n", len(other))
```

---

### 4) Check if symbol exists

```go
func SymbolExists(sugar *mt5.MT5Sugar, searchSymbol string) bool {
    symbols, err := sugar.GetAllSymbols()
    if err != nil {
        return false
    }

    for _, symbol := range symbols {
        if symbol == searchSymbol {
            return true
        }
    }

    return false
}

// Usage:
if SymbolExists(sugar, "EURUSD") {
    fmt.Println("✅ EURUSD is available")
} else {
    fmt.Println("❌ EURUSD not found")
}
```

---

### 5) Build custom watchlist

```go
// Define symbols of interest
watchlist := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"}

allSymbols, _ := sugar.GetAllSymbols()

// Convert to map for faster lookup
symbolMap := make(map[string]bool)
for _, symbol := range allSymbols {
    symbolMap[symbol] = true
}

fmt.Println("Watchlist Validation:")
fmt.Println("─────────────────────────────────────────")

availableSymbols := []string{}

for _, symbol := range watchlist {
    if symbolMap[symbol] {
        fmt.Printf("✅ %s - Available\n", symbol)
        availableSymbols = append(availableSymbols, symbol)
    } else {
        fmt.Printf("❌ %s - Not found\n", symbol)
    }
}

fmt.Printf("\n%d/%d symbols available\n",
    len(availableSymbols), len(watchlist))
```

---

### 6) Find symbols by suffix

```go
func FindSymbolsBySuffix(symbols []string, suffix string) []string {
    matches := []string{}

    for _, symbol := range symbols {
        if strings.HasSuffix(symbol, suffix) {
            matches = append(matches, symbol)
        }
    }

    return matches
}

// Usage:
allSymbols, _ := sugar.GetAllSymbols()

// Find all pairs vs USD
usdPairs := FindSymbolsBySuffix(allSymbols, "USD")
fmt.Printf("Pairs vs USD: %d\n", len(usdPairs))

// Find all pairs vs EUR
eurPairs := FindSymbolsBySuffix(allSymbols, "EUR")
fmt.Printf("Pairs vs EUR: %d\n", len(eurPairs))
```

---

### 7) Symbol market scanner

```go
func ScanMarketForSymbols(sugar *mt5.MT5Sugar, pattern string) {
    symbols, _ := sugar.GetAllSymbols()

    fmt.Printf("Scanning for symbols matching '%s'...\n\n", pattern)

    matches := []string{}

    for _, symbol := range symbols {
        if strings.Contains(strings.ToUpper(symbol), strings.ToUpper(pattern)) {
            matches = append(matches, symbol)
        }
    }

    if len(matches) == 0 {
        fmt.Println("No matches found")
        return
    }

    fmt.Printf("Found %d matches:\n", len(matches))
    for i, symbol := range matches {
        fmt.Printf("%d. %s\n", i+1, symbol)
    }
}

// Usage:
ScanMarketForSymbols(sugar, "GBP")  // Find all GBP pairs
ScanMarketForSymbols(sugar, "XAU")  // Find gold symbols
```

---

### 8) Export symbols to file

```go
func ExportSymbolsToFile(symbols []string, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    for _, symbol := range symbols {
        fmt.Fprintln(writer, symbol)
    }

    return nil
}

// Usage:
symbols, _ := sugar.GetAllSymbols()
err := ExportSymbolsToFile(symbols, "all_symbols.txt")
if err != nil {
    fmt.Printf("Export failed: %v\n", err)
} else {
    fmt.Printf("✅ Exported %d symbols to all_symbols.txt\n", len(symbols))
}
```

---

### 9) Symbol statistics

```go
func AnalyzeSymbolList(symbols []string) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      SYMBOL STATISTICS                ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Total symbols:    %d\n\n", len(symbols))

    // Length distribution
    lengthCount := make(map[int]int)
    for _, symbol := range symbols {
        lengthCount[len(symbol)]++
    }

    fmt.Println("Symbol length distribution:")
    for length := 1; length <= 20; length++ {
        if count := lengthCount[length]; count > 0 {
            fmt.Printf("  %d chars:  %d symbols\n", length, count)
        }
    }

    // Common prefixes
    prefixCount := make(map[string]int)
    for _, symbol := range symbols {
        if len(symbol) >= 3 {
            prefix := symbol[:3]
            prefixCount[prefix]++
        }
    }

    fmt.Println("\nMost common prefixes:")
    type prefixPair struct {
        prefix string
        count  int
    }

    pairs := []prefixPair{}
    for prefix, count := range prefixCount {
        pairs = append(pairs, prefixPair{prefix, count})
    }

    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].count > pairs[j].count
    })

    for i := 0; i < 5 && i < len(pairs); i++ {
        fmt.Printf("  %s: %d symbols\n", pairs[i].prefix, pairs[i].count)
    }
}

// Usage:
symbols, _ := sugar.GetAllSymbols()
AnalyzeSymbolList(symbols)
```

---

### 10) Advanced symbol manager

```go
type SymbolManager struct {
    AllSymbols    []string
    SymbolMap     map[string]bool
    Favorites     []string
}

func NewSymbolManager(sugar *mt5.MT5Sugar) (*SymbolManager, error) {
    symbols, err := sugar.GetAllSymbols()
    if err != nil {
        return nil, err
    }

    symbolMap := make(map[string]bool)
    for _, symbol := range symbols {
        symbolMap[symbol] = true
    }

    return &SymbolManager{
        AllSymbols: symbols,
        SymbolMap:  symbolMap,
        Favorites:  []string{},
    }, nil
}

func (sm *SymbolManager) Exists(symbol string) bool {
    return sm.SymbolMap[symbol]
}

func (sm *SymbolManager) AddFavorite(symbol string) error {
    if !sm.Exists(symbol) {
        return fmt.Errorf("symbol %s not available", symbol)
    }

    // Check if already in favorites
    for _, fav := range sm.Favorites {
        if fav == symbol {
            return fmt.Errorf("symbol %s already in favorites", symbol)
        }
    }

    sm.Favorites = append(sm.Favorites, symbol)
    return nil
}

func (sm *SymbolManager) SearchByPattern(pattern string) []string {
    matches := []string{}

    for _, symbol := range sm.AllSymbols {
        if strings.Contains(strings.ToUpper(symbol), strings.ToUpper(pattern)) {
            matches = append(matches, symbol)
        }
    }

    return matches
}

func (sm *SymbolManager) GetForexPairs() []string {
    forex := []string{}

    for _, symbol := range sm.AllSymbols {
        // Simple heuristic: 6 chars, all letters
        if len(symbol) == 6 {
            allLetters := true
            for _, c := range symbol {
                if !unicode.IsLetter(c) {
                    allLetters = false
                    break
                }
            }
            if allLetters {
                forex = append(forex, symbol)
            }
        }
    }

    return forex
}

func (sm *SymbolManager) ShowReport() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      SYMBOL MANAGER REPORT            ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Total symbols:     %d\n", len(sm.AllSymbols))
    fmt.Printf("Favorites:         %d\n", len(sm.Favorites))

    forexPairs := sm.GetForexPairs()
    fmt.Printf("Forex pairs:       %d\n\n", len(forexPairs))

    if len(sm.Favorites) > 0 {
        fmt.Println("Favorite Symbols:")
        for i, symbol := range sm.Favorites {
            fmt.Printf("  %d. %s\n", i+1, symbol)
        }
    }
}

// Usage:
manager, _ := NewSymbolManager(sugar)

// Add favorites
manager.AddFavorite("EURUSD")
manager.AddFavorite("GBPUSD")
manager.AddFavorite("XAUUSD")

// Search
gbpSymbols := manager.SearchByPattern("GBP")
fmt.Printf("Found %d GBP symbols\n", len(gbpSymbols))

// Report
manager.ShowReport()

// Check availability
if manager.Exists("BTCUSD") {
    fmt.Println("✅ BTCUSD available")
}
```

---

## 🔗 Related Methods

**🍬 Symbol details:**

* `GetSymbolInfo()` - Get full symbol information
* `IsSymbolAvailable()` - Check if specific symbol exists and is synchronized

**🍬 Symbol properties:**

* `GetMinStopLevel()` - Get stop level
* `GetSymbolDigits()` - Get decimal precision

---

## ⚠️ Common Pitfalls

### 1) Assuming all symbols are tradeable

```go
// ❌ WRONG - symbol in list doesn't mean it's tradeable
symbols, _ := sugar.GetAllSymbols()
// Some might be for display only!

// ✅ BETTER - at least verify the symbol exists and is synchronized
// (IsSymbolAvailable does NOT check the symbol's trade mode / trading permissions)
for _, symbol := range symbols {
    available, _ := sugar.IsSymbolAvailable(symbol)
    if available {
        // Symbol exists and has synchronized data
    }
}
```

### 2) Not handling large lists

```go
// ❌ WRONG - might be thousands of symbols (12,000+ on brokers with stocks)
symbols, _ := sugar.GetAllSymbols()
for _, symbol := range symbols {
    // Doing expensive operation per symbol
    info, _ := sugar.GetSymbolInfo(symbol) // Slow!
}

// ✅ CORRECT - filter first, then get details
symbols, _ := sugar.GetAllSymbols()
usdPairs := []string{}
for _, s := range symbols {
    if strings.HasSuffix(s, "USD") {
        usdPairs = append(usdPairs, s)
    }
}
for _, symbol := range usdPairs {
    info, _ := sugar.GetSymbolInfo(symbol)
}
```

### 3) Case-insensitive matching

```go
// ❌ WRONG - case matters!
if symbol == "eurusd" {

// ✅ CORRECT - case-insensitive comparison
if strings.EqualFold(symbol, "eurusd") {
```

---

## 💎 Pro Tips

1. **Cache the list** - Symbols rarely change, cache to avoid repeated calls

2. **Filter before details** - Get symbol list, filter, then fetch details

3. **Case handling** - Symbol names are case-sensitive, use exact names

4. **Availability** - Symbol in list doesn't guarantee tradeability

5. **Large lists** - Brokers may have thousands of symbols (all pages are read, 1000 per request), handle efficiently

---

## 📊 Common Symbol Patterns

```
Forex pairs (6 chars):
- EURUSD, GBPUSD, USDJPY, etc.
- Pattern: [A-Z]{6}

Metals:
- XAUUSD (Gold vs USD)
- XAGUSD (Silver vs USD)
- Pattern: XA[GU]USD

Indices:
- US30, NAS100, SPX500
- Pattern: [A-Z]+[0-9]+

Crypto:
- BTCUSD, ETHUSD
- Pattern: .*BTC.*, .*ETH.*
```

---

**See also:** [`GetSymbolInfo.md`](GetSymbolInfo.md), [`IsSymbolAvailable.md`](IsSymbolAvailable.md)
