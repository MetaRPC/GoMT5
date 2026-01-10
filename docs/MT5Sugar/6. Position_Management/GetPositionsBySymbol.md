# 🎯 Get Positions By Symbol (`GetPositionsBySymbol`)

> **Sugar method:** Retrieves all open positions for a specific trading symbol.

**API Information:**

* **Method:** `sugar.GetPositionsBySymbol(symbol)`
* **Timeout:** 3 seconds
* **Returns:** Slice of `*PositionInfo` structures

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*PositionInfo, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `[]*PositionInfo` | slice | Positions for specified symbol |
| `error` | `error` | Error if retrieval failed |

---

## 💬 Just the Essentials

* **What it is:** Filter open positions to show only one symbol.
* **Why you need it:** Symbol-specific operations, hedging, symbol P/L tracking.
* **Sanity check:** Returns empty slice if no positions for that symbol.

---

## 🎯 When to Use

✅ **Symbol analysis** - Check exposure on specific pair

✅ **Hedging** - See existing positions before opening opposite

✅ **Symbol P/L** - Calculate profit for one symbol

✅ **Batch close** - Close all positions for specific symbol

---

## 🔗 Usage Examples

### 1) Basic usage - get symbol positions

```go
symbol := "EURUSD"

positions, err := sugar.GetPositionsBySymbol(symbol)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if len(positions) == 0 {
    fmt.Printf("No open %s positions\n", symbol)
    return
}

fmt.Printf("%s positions: %d\n\n", symbol, len(positions))

for _, pos := range positions {
    fmt.Printf("#%d: %s %.2f lots - $%.2f\n",
        pos.Ticket, pos.Type, pos.Volume, pos.Profit)
}
```

---

### 2) Calculate symbol P/L

```go
symbol := "EURUSD"

positions, _ := sugar.GetPositionsBySymbol(symbol)

if len(positions) == 0 {
    fmt.Printf("No %s positions\n", symbol)
    return
}

totalProfit := 0.0
totalVolume := 0.0

for _, pos := range positions {
    totalProfit += pos.Profit
    totalVolume += pos.Volume
}

fmt.Printf("%s Summary:\n", symbol)
fmt.Printf("  Positions: %d\n", len(positions))
fmt.Printf("  Volume:    %.2f lots\n", totalVolume)
fmt.Printf("  Total P/L: $%.2f\n", totalProfit)
```

---

### 3) Check net position (BUY vs SELL)

```go
symbol := "EURUSD"

positions, _ := sugar.GetPositionsBySymbol(symbol)

buyVolume := 0.0
sellVolume := 0.0

for _, pos := range positions {
    if pos.Type == "BUY" {
        buyVolume += pos.Volume
    } else {
        sellVolume += pos.Volume
    }
}

netVolume := buyVolume - sellVolume

fmt.Printf("%s Net Position:\n", symbol)
fmt.Printf("  BUY:  %.2f lots\n", buyVolume)
fmt.Printf("  SELL: %.2f lots\n", sellVolume)
fmt.Printf("  NET:  %.2f lots ", netVolume)

if netVolume > 0 {
    fmt.Println("(NET LONG)")
} else if netVolume < 0 {
    fmt.Println("(NET SHORT)")
} else {
    fmt.Println("(NEUTRAL)")
}
```

---

### 4) Close all positions for symbol

```go
symbol := "EURUSD"

positions, _ := sugar.GetPositionsBySymbol(symbol)

if len(positions) == 0 {
    fmt.Printf("No %s positions to close\n", symbol)
    return
}

fmt.Printf("Closing %d %s positions...\n", len(positions), symbol)

closedCount := 0
totalProfit := 0.0

for _, pos := range positions {
    totalProfit += pos.Profit

    err := sugar.ClosePosition(pos.Ticket)
    if err != nil {
        fmt.Printf("❌ Failed to close #%d: %v\n", pos.Ticket, err)
        continue
    }

    closedCount++
    fmt.Printf("✅ Closed #%d ($%.2f)\n", pos.Ticket, pos.Profit)
}

fmt.Printf("\nClosed %d/%d positions\n", closedCount, len(positions))
fmt.Printf("Total P/L: $%.2f\n", totalProfit)
```

---

### 5) Check if hedged

```go
symbol := "EURUSD"

positions, _ := sugar.GetPositionsBySymbol(symbol)

if len(positions) == 0 {
    fmt.Printf("No %s positions\n", symbol)
    return
}

hasBuy := false
hasSell := false

for _, pos := range positions {
    if pos.Type == "BUY" {
        hasBuy = true
    } else {
        hasSell = true
    }
}

if hasBuy && hasSell {
    fmt.Printf("✅ %s is HEDGED (both BUY and SELL positions)\n", symbol)
} else if hasBuy {
    fmt.Printf("📈 %s: Only LONG positions\n", symbol)
} else {
    fmt.Printf("📉 %s: Only SHORT positions\n", symbol)
}
```

---

### 6) Average entry price

```go
symbol := "EURUSD"

positions, _ := sugar.GetPositionsBySymbol(symbol)

if len(positions) == 0 {
    fmt.Printf("No %s positions\n", symbol)
    return
}

// Separate BUY and SELL
buyPositions := []*PositionInfo{}
sellPositions := []*PositionInfo{}

for _, pos := range positions {
    if pos.Type == "BUY" {
        buyPositions = append(buyPositions, pos)
    } else {
        sellPositions = append(sellPositions, pos)
    }
}

// Calculate weighted average for BUY
if len(buyPositions) > 0 {
    totalWeightedPrice := 0.0
    totalVolume := 0.0

    for _, pos := range buyPositions {
        totalWeightedPrice += pos.PriceOpen * pos.Volume
        totalVolume += pos.Volume
    }

    avgEntry := totalWeightedPrice / totalVolume
    fmt.Printf("BUY average entry: %.5f (%.2f lots)\n", avgEntry, totalVolume)
}

// Calculate weighted average for SELL
if len(sellPositions) > 0 {
    totalWeightedPrice := 0.0
    totalVolume := 0.0

    for _, pos := range sellPositions {
        totalWeightedPrice += pos.PriceOpen * pos.Volume
        totalVolume += pos.Volume
    }

    avgEntry := totalWeightedPrice / totalVolume
    fmt.Printf("SELL average entry: %.5f (%.2f lots)\n", avgEntry, totalVolume)
}
```

---

### 7) Find oldest and newest positions

```go
symbol := "EURUSD"

positions, _ := sugar.GetPositionsBySymbol(symbol)

if len(positions) == 0 {
    fmt.Printf("No %s positions\n", symbol)
    return
}

oldest := positions[0]
newest := positions[0]

for _, pos := range positions {
    if pos.TimeOpen.Before(oldest.TimeOpen) {
        oldest = pos
    }
    if pos.TimeOpen.After(newest.TimeOpen) {
        newest = pos
    }
}

fmt.Printf("%s Position Age:\n", symbol)
fmt.Printf("Oldest: #%d opened %v ago ($%.2f)\n",
    oldest.Ticket, time.Since(oldest.TimeOpen).Round(time.Minute), oldest.Profit)
fmt.Printf("Newest: #%d opened %v ago ($%.2f)\n",
    newest.Ticket, time.Since(newest.TimeOpen).Round(time.Minute), newest.Profit)
```

---

### 8) Symbol exposure check before trading

```go
symbol := "EURUSD"
maxPositions := 3
maxVolume := 0.5

positions, _ := sugar.GetPositionsBySymbol(symbol)

currentCount := len(positions)
currentVolume := 0.0

for _, pos := range positions {
    currentVolume += pos.Volume
}

fmt.Printf("%s Exposure Check:\n", symbol)
fmt.Printf("  Current positions: %d / %d max\n", currentCount, maxPositions)
fmt.Printf("  Current volume:    %.2f / %.2f max\n", currentVolume, maxVolume)

canTrade := true
newVolume := 0.1

if currentCount >= maxPositions {
    fmt.Println("  ❌ Max positions reached")
    canTrade = false
}

if currentVolume+newVolume > maxVolume {
    fmt.Println("  ❌ Max volume would be exceeded")
    canTrade = false
}

if canTrade {
    fmt.Printf("  ✅ Can open %.2f lots\n", newVolume)
} else {
    fmt.Println("  ⛔ Cannot open new position")
}
```

---

### 9) Monitor symbol positions dashboard

```go
func MonitorSymbol(sugar *mt5.MT5Sugar, symbol string) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        positions, _ := sugar.GetPositionsBySymbol(symbol)

        // Clear screen
        fmt.Print("\033[H\033[2J")

        fmt.Println("╔═══════════════════════════════════════╗")
        fmt.Printf("║      %s POSITIONS", symbol)
        for i := 0; i < 23-len(symbol); i++ {
            fmt.Print(" ")
        }
        fmt.Println("║")
        fmt.Println("╚═══════════════════════════════════════╝")
        fmt.Printf("Updated: %s\n\n", time.Now().Format("15:04:05"))

        if len(positions) == 0 {
            fmt.Printf("No %s positions\n", symbol)
            continue
        }

        buyCount, sellCount := 0, 0
        buyVolume, sellVolume := 0.0, 0.0
        buyProfit, sellProfit := 0.0, 0.0

        for _, pos := range positions {
            if pos.Type == "BUY" {
                buyCount++
                buyVolume += pos.Volume
                buyProfit += pos.Profit
            } else {
                sellCount++
                sellVolume += pos.Volume
                sellProfit += pos.Profit
            }
        }

        fmt.Printf("BUY:  %d positions, %.2f lots, $%.2f\n",
            buyCount, buyVolume, buyProfit)
        fmt.Printf("SELL: %d positions, %.2f lots, $%.2f\n",
            sellCount, sellVolume, sellProfit)
        fmt.Println("─────────────────────────────────────────")
        fmt.Printf("NET:  %.2f lots, $%.2f\n",
            buyVolume-sellVolume, buyProfit+sellProfit)

        fmt.Println("\nPositions:")
        for _, pos := range positions {
            fmt.Printf("  #%d %s %.2f lots @ %.5f → $%.2f\n",
                pos.Ticket, pos.Type, pos.Volume, pos.PriceOpen, pos.Profit)
        }
    }
}

// Usage:
MonitorSymbol(sugar, "EURUSD")
```

---

### 10) Advanced symbol position analysis

```go
func AnalyzeSymbolPositions(sugar *mt5.MT5Sugar, symbol string) {
    positions, _ := sugar.GetPositionsBySymbol(symbol)

    if len(positions) == 0 {
        fmt.Printf("No %s positions to analyze\n", symbol)
        return
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Printf("║      %s ANALYSIS", symbol)
    for i := 0; i < 22-len(symbol); i++ {
        fmt.Print(" ")
    }
    fmt.Println("║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Separate by type
    buyPos := []*PositionInfo{}
    sellPos := []*PositionInfo{}

    for _, pos := range positions {
        if pos.Type == "BUY" {
            buyPos = append(buyPos, pos)
        } else {
            sellPos = append(sellPos, pos)
        }
    }

    // BUY analysis
    if len(buyPos) > 0 {
        buyVolume := 0.0
        buyProfit := 0.0
        buyWeightedEntry := 0.0

        for _, pos := range buyPos {
            buyVolume += pos.Volume
            buyProfit += pos.Profit
            buyWeightedEntry += pos.PriceOpen * pos.Volume
        }

        avgBuyEntry := buyWeightedEntry / buyVolume

        fmt.Println("\n📈 LONG Positions:")
        fmt.Printf("   Count:       %d\n", len(buyPos))
        fmt.Printf("   Volume:      %.2f lots\n", buyVolume)
        fmt.Printf("   Avg Entry:   %.5f\n", avgBuyEntry)
        fmt.Printf("   Total P/L:   $%.2f\n", buyProfit)
    }

    // SELL analysis
    if len(sellPos) > 0 {
        sellVolume := 0.0
        sellProfit := 0.0
        sellWeightedEntry := 0.0

        for _, pos := range sellPos {
            sellVolume += pos.Volume
            sellProfit += pos.Profit
            sellWeightedEntry += pos.PriceOpen * pos.Volume
        }

        avgSellEntry := sellWeightedEntry / sellVolume

        fmt.Println("\n📉 SHORT Positions:")
        fmt.Printf("   Count:       %d\n", len(sellPos))
        fmt.Printf("   Volume:      %.2f lots\n", sellVolume)
        fmt.Printf("   Avg Entry:   %.5f\n", avgSellEntry)
        fmt.Printf("   Total P/L:   $%.2f\n", sellProfit)
    }

    // Net position
    netVolume := 0.0
    netProfit := 0.0
    for _, pos := range positions {
        if pos.Type == "BUY" {
            netVolume += pos.Volume
        } else {
            netVolume -= pos.Volume
        }
        netProfit += pos.Profit
    }

    fmt.Println("\n═══════════════════════════════════════")
    fmt.Printf("Net Position:    %.2f lots ", netVolume)
    if netVolume > 0 {
        fmt.Println("(LONG)")
    } else if netVolume < 0 {
        fmt.Println("(SHORT)")
    } else {
        fmt.Println("(NEUTRAL)")
    }
    fmt.Printf("Combined P/L:    $%.2f\n", netProfit)
}

// Usage:
AnalyzeSymbolPositions(sugar, "EURUSD")
```

---

## 🔗 Related Methods

**🍬 Related position info:**

* `GetOpenPositions()` - Get all positions (all symbols)
* `GetPositionByTicket()` - Get specific position

**🍬 Symbol-specific operations:**

* `CloseAllPositions()` - Close all (use with filter for symbol)

---

## ⚠️ Common Pitfalls

### 1) Wrong symbol case

```go
// ❌ WRONG - case sensitive!
positions, _ := sugar.GetPositionsBySymbol("eurusd") // Won't find "EURUSD"

// ✅ CORRECT - match MT5 symbol exactly
positions, _ := sugar.GetPositionsBySymbol("EURUSD")
```

### 2) Not checking empty result

```go
// ❌ WRONG - assuming positions exist
positions, _ := sugar.GetPositionsBySymbol("EURUSD")
firstPos := positions[0] // Might panic!

// ✅ CORRECT - check length
if len(positions) > 0 {
    firstPos := positions[0]
}
```

### 3) Using instead of GetOpenPositions + filter

```go
// ❌ INEFFICIENT - if you need all symbols anyway
eurusdPos, _ := sugar.GetPositionsBySymbol("EURUSD")
gbpusdPos, _ := sugar.GetPositionsBySymbol("GBPUSD")
// Multiple API calls!

// ✅ BETTER - one call, filter in code
allPositions, _ := sugar.GetOpenPositions()
// Filter yourself if needed
```

---

## 💎 Pro Tips

1. **Symbol case** - Use exact MT5 symbol name (usually uppercase)

2. **Net position** - Calculate BUY - SELL to see exposure

3. **Before trading** - Check existing positions to avoid over-exposure

4. **Hedging** - Useful to see both sides of hedged positions

5. **Batch operations** - Collect tickets before closing/modifying

---

**See also:** [`GetOpenPositions.md`](GetOpenPositions.md), [`GetPositionByTicket.md`](GetPositionByTicket.md), [`ClosePosition.md`](ClosePosition.md)
