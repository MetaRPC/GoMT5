# 📋 Get Open Positions (`GetOpenPositions`)

> **Sugar method:** Retrieves all currently open positions across all symbols.

**API Information:**

* **Method:** `sugar.GetOpenPositions()`
* **Timeout:** 3 seconds
* **Returns:** Slice of `*PositionInfo` structures

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetOpenPositions() ([]*PositionInfo, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters required |

| Output | Type | Description |
|--------|------|-------------|
| `[]*PositionInfo` | slice | All open positions |
| `error` | `error` | Error if retrieval failed |

---

## 💬 Just the Essentials

* **What it is:** Get complete list of all open positions.
* **Why you need it:** Dashboard, P/L tracking, risk monitoring, batch operations.
* **Sanity check:** Returns empty slice if no positions. Never returns nil.

---

## 🎯 When to Use

✅ **Dashboard** - Show all active trades

✅ **Total P/L** - Calculate combined profit/loss

✅ **Risk check** - See total exposure

✅ **Batch operations** - Close/modify multiple positions

---

## 🔗 Usage Examples

### 1) Basic usage - list all positions

```go
positions, err := sugar.GetOpenPositions()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if len(positions) == 0 {
    fmt.Println("No open positions")
    return
}

fmt.Printf("Open positions: %d\n\n", len(positions))

for _, pos := range positions {
    fmt.Printf("#%d: %s %s %.2f lots - $%.2f\n",
        pos.Ticket, pos.Symbol, pos.Type, pos.Volume, pos.Profit)
}
```

---

### 2) Calculate total P/L

```go
positions, _ := sugar.GetOpenPositions()

if len(positions) == 0 {
    fmt.Println("No positions")
    return
}

totalProfit := 0.0
totalCommission := 0.0
totalSwap := 0.0

for _, pos := range positions {
    totalProfit += pos.Profit
    totalCommission += pos.Commission
    totalSwap += pos.Swap
}

netProfit := totalProfit + totalCommission + totalSwap

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║       TOTAL P/L SUMMARY               ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Positions:   %d\n", len(positions))
fmt.Printf("Profit:      $%.2f\n", totalProfit)
fmt.Printf("Commission:  $%.2f\n", totalCommission)
fmt.Printf("Swap:        $%.2f\n", totalSwap)
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Net P/L:     $%.2f\n", netProfit)
```

---

### 3) Group by symbol

```go
positions, _ := sugar.GetOpenPositions()

// Group positions by symbol
symbolMap := make(map[string][]*PositionInfo)

for _, pos := range positions {
    symbolMap[pos.Symbol] = append(symbolMap[pos.Symbol], pos)
}

fmt.Println("Positions by symbol:")
fmt.Println("─────────────────────────────────────────")

for symbol, posList := range symbolMap {
    totalVolume := 0.0
    totalProfit := 0.0

    for _, pos := range posList {
        totalVolume += pos.Volume
        totalProfit += pos.Profit
    }

    fmt.Printf("%-8s: %d positions, %.2f lots, $%.2f\n",
        symbol, len(posList), totalVolume, totalProfit)
}
```

---

### 4) Count BUY vs SELL

```go
positions, _ := sugar.GetOpenPositions()

buyCount := 0
sellCount := 0
buyProfit := 0.0
sellProfit := 0.0

for _, pos := range positions {
    if pos.Type == "BUY" {
        buyCount++
        buyProfit += pos.Profit
    } else {
        sellCount++
        sellProfit += pos.Profit
    }
}

fmt.Println("Position breakdown:")
fmt.Printf("  BUY:  %d positions ($%.2f)\n", buyCount, buyProfit)
fmt.Printf("  SELL: %d positions ($%.2f)\n", sellCount, sellProfit)
```

---

### 5) Find largest winning/losing position

```go
positions, _ := sugar.GetOpenPositions()

if len(positions) == 0 {
    fmt.Println("No positions")
    return
}

var largestWin, largestLoss *PositionInfo

for _, pos := range positions {
    if pos.Profit > 0 {
        if largestWin == nil || pos.Profit > largestWin.Profit {
            largestWin = pos
        }
    } else {
        if largestLoss == nil || pos.Profit < largestLoss.Profit {
            largestLoss = pos
        }
    }
}

if largestWin != nil {
    fmt.Printf("🏆 Largest winner: #%d %s - $%.2f\n",
        largestWin.Ticket, largestWin.Symbol, largestWin.Profit)
}

if largestLoss != nil {
    fmt.Printf("📉 Largest loser: #%d %s - $%.2f\n",
        largestLoss.Ticket, largestLoss.Symbol, largestLoss.Profit)
}
```

---

### 6) Positions without SL/TP

```go
positions, _ := sugar.GetOpenPositions()

unprotected := []*PositionInfo{}

for _, pos := range positions {
    if pos.StopLoss == 0 || pos.TakeProfit == 0 {
        unprotected = append(unprotected, pos)
    }
}

if len(unprotected) > 0 {
    fmt.Printf("⚠️  %d positions without full protection:\n", len(unprotected))
    for _, pos := range unprotected {
        fmt.Printf("  #%d %s: ", pos.Ticket, pos.Symbol)
        if pos.StopLoss == 0 {
            fmt.Print("No SL ")
        }
        if pos.TakeProfit == 0 {
            fmt.Print("No TP")
        }
        fmt.Println()
    }
} else {
    fmt.Println("✅ All positions have SL and TP")
}
```

---

### 7) Real-time dashboard

```go
func PositionsDashboard(sugar *mt5.MT5Sugar) {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        positions, _ := sugar.GetOpenPositions()

        // Clear screen
        fmt.Print("\033[H\033[2J")

        fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
        fmt.Println("║                    POSITIONS DASHBOARD                        ║")
        fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
        fmt.Printf("Updated: %s\n\n", time.Now().Format("15:04:05"))

        if len(positions) == 0 {
            fmt.Println("No open positions")
            continue
        }

        fmt.Printf("%-8s %-8s %-6s %-10s %-10s %-10s\n",
            "Ticket", "Symbol", "Type", "Lots", "Entry", "P/L")
        fmt.Println("─────────────────────────────────────────────────────────────────")

        totalPL := 0.0
        for _, pos := range positions {
            fmt.Printf("%-8d %-8s %-6s %-10.2f %-10.5f %-10.2f\n",
                pos.Ticket, pos.Symbol, pos.Type,
                pos.Volume, pos.PriceOpen, pos.Profit)
            totalPL += pos.Profit
        }

        fmt.Println("─────────────────────────────────────────────────────────────────")
        fmt.Printf("Total positions: %d | Total P/L: $%.2f\n", len(positions), totalPL)
    }
}

// Usage:
PositionsDashboard(sugar)
```

---

### 8) Calculate total risk

```go
positions, _ := sugar.GetOpenPositions()

totalRisk := 0.0

for _, pos := range positions {
    if pos.StopLoss > 0 {
        // Calculate potential loss if SL hit
        var riskPips float64
        if pos.Type == "BUY" {
            riskPips = (pos.PriceOpen - pos.StopLoss) / 0.00001
        } else {
            riskPips = (pos.StopLoss - pos.PriceOpen) / 0.00001
        }

        // Estimate risk in dollars (simplified)
        info, _ := sugar.GetSymbolInfo(pos.Symbol)
        riskPerLot := riskPips * info.Point * info.ContractSize
        positionRisk := riskPerLot * pos.Volume

        totalRisk += positionRisk
    }
}

balance, _ := sugar.GetBalance()
riskPercent := (totalRisk / balance) * 100

fmt.Printf("Total account risk:\n")
fmt.Printf("  Balance:     $%.2f\n", balance)
fmt.Printf("  Total risk:  $%.2f\n", totalRisk)
fmt.Printf("  Risk %%:      %.2f%%\n", riskPercent)

if riskPercent > 10 {
    fmt.Println("  ⚠️  Risk over 10%!")
}
```

---

### 9) Export positions to CSV format

```go
positions, _ := sugar.GetOpenPositions()

if len(positions) == 0 {
    fmt.Println("No positions to export")
    return
}

fmt.Println("Ticket,Symbol,Type,Volume,Entry,Current,SL,TP,Profit")

for _, pos := range positions {
    fmt.Printf("%d,%s,%s,%.2f,%.5f,%.5f,%.5f,%.5f,%.2f\n",
        pos.Ticket, pos.Symbol, pos.Type, pos.Volume,
        pos.PriceOpen, pos.PriceCurrent,
        pos.StopLoss, pos.TakeProfit, pos.Profit)
}
```

---

### 10) Advanced position analysis

```go
func AnalyzePositions(sugar *mt5.MT5Sugar) {
    positions, _ := sugar.GetOpenPositions()

    if len(positions) == 0 {
        fmt.Println("No positions to analyze")
        return
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      POSITION ANALYSIS                ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Statistics
    totalPositions := len(positions)
    totalVolume := 0.0
    totalProfit := 0.0
    winningCount := 0
    losingCount := 0
    oldestTime := positions[0].TimeOpen

    for _, pos := range positions {
        totalVolume += pos.Volume
        totalProfit += pos.Profit

        if pos.Profit > 0 {
            winningCount++
        } else if pos.Profit < 0 {
            losingCount++
        }

        if pos.TimeOpen.Before(oldestTime) {
            oldestTime = pos.TimeOpen
        }
    }

    winRate := 0.0
    if totalPositions > 0 {
        winRate = (float64(winningCount) / float64(totalPositions)) * 100
    }

    avgProfit := totalProfit / float64(totalPositions)

    fmt.Printf("Total positions: %d\n", totalPositions)
    fmt.Printf("Total volume:    %.2f lots\n", totalVolume)
    fmt.Printf("Total P/L:       $%.2f\n", totalProfit)
    fmt.Printf("Avg P/L:         $%.2f\n", avgProfit)
    fmt.Println()
    fmt.Printf("Winning:         %d (%.1f%%)\n", winningCount, winRate)
    fmt.Printf("Losing:          %d\n", losingCount)
    fmt.Printf("Breakeven:       %d\n", totalPositions-winningCount-losingCount)
    fmt.Println()
    fmt.Printf("Oldest position: %v ago\n",
        time.Since(oldestTime).Round(time.Minute))

    // Symbol exposure
    symbolCount := make(map[string]int)
    for _, pos := range positions {
        symbolCount[pos.Symbol]++
    }

    fmt.Println("\nSymbol exposure:")
    for symbol, count := range symbolCount {
        percent := (float64(count) / float64(totalPositions)) * 100
        fmt.Printf("  %-8s: %d (%.1f%%)\n", symbol, count, percent)
    }
}

// Usage:
AnalyzePositions(sugar)
```

---

## 🔗 Related Methods

**🍬 Related position info:**

* `GetPositionByTicket()` - Get specific position
* `GetPositionsBySymbol()` - Filter by symbol

**🍬 Batch operations:**

* `CloseAllPositions()` - Close all positions
* `GetProfit()` - Total P/L (faster than summing manually)

---

## ⚠️ Common Pitfalls

### 1) Not checking for empty slice

```go
// ❌ WRONG - assuming positions exist
positions, _ := sugar.GetOpenPositions()
firstPos := positions[0] // Might panic!

// ✅ CORRECT - check length
positions, _ := sugar.GetOpenPositions()
if len(positions) > 0 {
    firstPos := positions[0]
}
```

### 2) Modifying while iterating

```go
// ❌ WRONG - modifying positions during iteration
for _, pos := range positions {
    sugar.ClosePosition(pos.Ticket) // Slice changes during loop!
}

// ✅ CORRECT - collect tickets first
tickets := []uint64{}
for _, pos := range positions {
    tickets = append(tickets, pos.Ticket)
}
for _, ticket := range tickets {
    sugar.ClosePosition(ticket)
}
```

### 3) Not handling errors

```go
// ❌ WRONG - ignoring errors
positions, _ := sugar.GetOpenPositions()

// ✅ CORRECT - check errors
positions, err := sugar.GetOpenPositions()
if err != nil {
    fmt.Printf("Failed to get positions: %v\n", err)
    return
}
```

---

## 💎 Pro Tips

1. **Cache results** - Don't call repeatedly in tight loops

2. **Filter efficiently** - Use GetPositionsBySymbol() for specific symbols

3. **Check length** - Always verify slice isn't empty

4. **Real-time updates** - Refresh every 2-5 seconds for dashboards

5. **Batch operations** - Collect tickets before modifying/closing

---

**See also:** [`GetPositionByTicket.md`](GetPositionByTicket.md), [`GetPositionsBySymbol.md`](GetPositionsBySymbol.md), [`CloseAllPositions.md`](CloseAllPositions.md)
