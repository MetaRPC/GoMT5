# 🔒🔁 Close All Positions (`CloseAllPositions`)

> **Sugar method:** Closes all open positions at once across all symbols.

**API Information:**

* **Method:** `sugar.CloseAllPositions()`
* **Timeout:** 10 seconds per position
* **Returns:** Error if any close failed

---

## 📋 Method Signature

```go
func (s *MT5Sugar) CloseAllPositions() error
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters required |

| Output | Type | Description |
|--------|------|-------------|
| `error` | `error` | `nil` if all closed successfully, error otherwise |

---

## 💬 Just the Essentials

* **What it is:** Emergency "close everything" button - closes all open positions.
* **Why you need it:** Emergency exit, end of trading day, major news event.
* **Sanity check:** Closes ALL positions - no undo! Use with caution.

---

## ⚠️ DANGER ZONE

**This method closes EVERY open position!**

- No confirmation prompt
- No way to undo
- Affects all symbols
- Use carefully in production!

---

## 🎯 When to Use

✅ **Emergency exit** - Market crash, major news, need to exit NOW

✅ **End of day** - Close all positions before weekend

✅ **Account protection** - Margin call prevention

✅ **Strategy reset** - Start fresh, close everything

❌ **NOT for selective closing** - Use `ClosePosition()` or `ClosePositionsBySymbol()` instead

---

## 🔗 Usage Examples

### 1) Basic usage - emergency close

```go
fmt.Println("⚠️  EMERGENCY: Closing all positions!")

err := sugar.CloseAllPositions()
if err != nil {
    fmt.Printf("❌ Failed to close all: %v\n", err)
    return
}

fmt.Println("✅ All positions closed")
```

---

### 2) Close all with confirmation

```go
positions, _ := sugar.GetOpenPositions()

if len(positions) == 0 {
    fmt.Println("No positions to close")
    return
}

fmt.Printf("⚠️  WARNING: About to close %d positions\n", len(positions))

// Show what will be closed
totalProfit := 0.0
for _, pos := range positions {
    fmt.Printf("  - #%d: %s %.2f lots ($%.2f)\n",
        pos.Ticket, pos.Symbol, pos.Volume, pos.Profit)
    totalProfit += pos.Profit
}

fmt.Printf("\nTotal P/L: $%.2f\n", totalProfit)
fmt.Println("\nClosing all positions in 3 seconds...")
time.Sleep(3 * time.Second)

err := sugar.CloseAllPositions()
if err != nil {
    fmt.Printf("Failed: %v\n", err)
} else {
    fmt.Printf("✅ Closed %d positions with $%.2f total P/L\n",
        len(positions), totalProfit)
}
```

---

### 3) End of trading day routine

```go
func EndOfDayClose(sugar *mt5.MT5Sugar) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      END OF DAY ROUTINE               ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Check current positions
    positions, _ := sugar.GetOpenPositions()
    fmt.Printf("Open positions: %d\n", len(positions))

    if len(positions) == 0 {
        fmt.Println("✅ No positions to close - ready for weekend")
        return
    }

    // Show summary
    var totalProfit float64
    buyCount, sellCount := 0, 0

    for _, pos := range positions {
        totalProfit += pos.Profit
        if pos.Type == "BUY" {
            buyCount++
        } else {
            sellCount++
        }
    }

    fmt.Printf("\nSummary:\n")
    fmt.Printf("  BUY positions:  %d\n", buyCount)
    fmt.Printf("  SELL positions: %d\n", sellCount)
    fmt.Printf("  Total P/L:      $%.2f\n", totalProfit)

    // Close all
    fmt.Println("\nClosing all positions...")
    err := sugar.CloseAllPositions()
    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }

    fmt.Println("✅ All positions closed - have a nice weekend!")
}

// Usage: Run at 16:55 on Friday
EndOfDayClose(sugar)
```

---

### 4) Close all on drawdown threshold

```go
maxDrawdown := -500.0 // $500 max loss

ticker := time.NewTicker(10 * time.Second)
defer ticker.Stop()

for range ticker.C {
    equity, _ := sugar.GetEquity()
    profit, _ := sugar.GetProfit()

    fmt.Printf("Equity: $%.2f, Floating P/L: $%.2f\n", equity, profit)

    if profit <= maxDrawdown {
        fmt.Printf("🚨 DRAWDOWN LIMIT REACHED: $%.2f\n", profit)
        fmt.Println("Closing all positions NOW!")

        err := sugar.CloseAllPositions()
        if err != nil {
            fmt.Printf("❌ Failed to close: %v\n", err)
        } else {
            fmt.Println("✅ All positions closed - loss limited")
        }
        return
    }
}
```

---

### 5) Close all before major news

```go
// News event in 5 minutes
newsTime := time.Now().Add(5 * time.Minute)

fmt.Printf("Major news at %s - closing all positions\n",
    newsTime.Format("15:04"))

positions, _ := sugar.GetOpenPositions()

if len(positions) > 0 {
    fmt.Printf("Closing %d positions before news...\n", len(positions))

    err := sugar.CloseAllPositions()
    if err != nil {
        fmt.Printf("❌ Failed: %v\n", err)
        return
    }

    fmt.Println("✅ All positions closed - ready for news event")
} else {
    fmt.Println("No positions to close")
}
```

---

### 6) Close all and show summary

```go
positions, _ := sugar.GetOpenPositions()

if len(positions) == 0 {
    fmt.Println("No positions to close")
    return
}

fmt.Printf("Closing %d positions:\n", len(positions))

// Record positions before closing
type ClosedPosition struct {
    Ticket uint64
    Symbol string
    Type   string
    Volume float64
    Profit float64
}

var closedPositions []ClosedPosition

for _, pos := range positions {
    closedPositions = append(closedPositions, ClosedPosition{
        Ticket: pos.Ticket,
        Symbol: pos.Symbol,
        Type:   pos.Type,
        Volume: pos.Volume,
        Profit: pos.Profit,
    })
}

// Close all
err := sugar.CloseAllPositions()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Show what was closed
fmt.Println("\n✅ Closed positions:")
fmt.Println("─────────────────────────────────────────")

totalProfit := 0.0
for _, pos := range closedPositions {
    fmt.Printf("#%d: %s %s %.2f lots - $%.2f\n",
        pos.Ticket, pos.Symbol, pos.Type, pos.Volume, pos.Profit)
    totalProfit += pos.Profit
}

fmt.Println("─────────────────────────────────────────")
fmt.Printf("Total P/L: $%.2f\n", totalProfit)
```

---

### 7) Close all with error handling

```go
positions, _ := sugar.GetOpenPositions()
totalPositions := len(positions)

if totalPositions == 0 {
    fmt.Println("No positions to close")
    return
}

fmt.Printf("Attempting to close %d positions...\n", totalPositions)

err := sugar.CloseAllPositions()

// Verify all closed
time.Sleep(2 * time.Second)
remainingPositions, _ := sugar.GetOpenPositions()

if len(remainingPositions) == 0 {
    fmt.Printf("✅ Success: All %d positions closed\n", totalPositions)
} else {
    fmt.Printf("⚠️  Warning: %d positions still open\n", len(remainingPositions))
    fmt.Println("Remaining positions:")
    for _, pos := range remainingPositions {
        fmt.Printf("  - #%d: %s %.2f lots\n", pos.Ticket, pos.Symbol, pos.Volume)
    }
}

if err != nil {
    fmt.Printf("Error occurred: %v\n", err)
}
```

---

### 8) Scheduled close all (cron-style)

```go
func ScheduledCloseAll(sugar *mt5.MT5Sugar, closeTime string) {
    fmt.Printf("Scheduled close at %s\n", closeTime)

    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        now := time.Now().Format("15:04")

        if now == closeTime {
            fmt.Printf("🕐 %s - Closing all positions\n", closeTime)

            positions, _ := sugar.GetOpenPositions()
            if len(positions) > 0 {
                err := sugar.CloseAllPositions()
                if err != nil {
                    fmt.Printf("Failed: %v\n", err)
                } else {
                    fmt.Printf("✅ Closed %d positions\n", len(positions))
                }
            } else {
                fmt.Println("No positions to close")
            }
            return
        }

        fmt.Printf("Waiting... (current: %s, target: %s)\n", now, closeTime)
    }
}

// Usage: Close all positions at 16:55
ScheduledCloseAll(sugar, "16:55")
```

---

### 9) Close all on margin level

```go
dangerMarginLevel := 150.0 // 150%

ticker := time.NewTicker(15 * time.Second)
defer ticker.Stop()

for range ticker.C {
    marginLevel, _ := sugar.GetMarginLevel()

    fmt.Printf("Margin Level: %.0f%%\n", marginLevel)

    if marginLevel < dangerMarginLevel {
        fmt.Printf("🚨 DANGER: Margin level %.0f%% < %.0f%%\n",
            marginLevel, dangerMarginLevel)

        positions, _ := sugar.GetOpenPositions()
        fmt.Printf("Closing all %d positions to protect account!\n", len(positions))

        err := sugar.CloseAllPositions()
        if err != nil {
            fmt.Printf("❌ EMERGENCY CLOSE FAILED: %v\n", err)
        } else {
            fmt.Println("✅ All positions closed - margin call avoided")
        }
        return
    }
}
```

---

### 10) Advanced close all with retry

```go
func CloseAllPositionsWithRetry(sugar *mt5.MT5Sugar, maxRetries int) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      CLOSING ALL POSITIONS            ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    positions, _ := sugar.GetOpenPositions()
    if len(positions) == 0 {
        fmt.Println("No positions to close")
        return nil
    }

    initialCount := len(positions)
    fmt.Printf("Positions to close: %d\n\n", initialCount)

    for attempt := 1; attempt <= maxRetries; attempt++ {
        fmt.Printf("Attempt %d/%d...\n", attempt, maxRetries)

        err := sugar.CloseAllPositions()

        // Check how many closed
        time.Sleep(2 * time.Second)
        remainingPositions, _ := sugar.GetOpenPositions()
        closedCount := initialCount - len(remainingPositions)

        fmt.Printf("  Closed: %d/%d\n", closedCount, initialCount)

        if len(remainingPositions) == 0 {
            fmt.Println("\n✅ All positions closed successfully!")
            return nil
        }

        if attempt < maxRetries {
            fmt.Printf("  %d positions still open - retrying...\n\n", len(remainingPositions))
            time.Sleep(2 * time.Second)
        }
    }

    // Final check
    finalPositions, _ := sugar.GetOpenPositions()
    return fmt.Errorf("failed to close all positions - %d still open", len(finalPositions))
}

// Usage:
err := CloseAllPositionsWithRetry(sugar, 3)
```

---

## 🔗 Related Methods

**🍬 Selective close methods:**

* `ClosePosition()` - Close single position
* `ClosePositionsBySymbol()` - Close all positions for specific symbol

**🍬 Position info:**

* `GetOpenPositions()` - See what will be closed
* `GetProfit()` - Check total P/L before closing

---

## ⚠️ Common Pitfalls

### 1) Not checking what you're closing

```go
// ❌ WRONG - blindly closing everything
sugar.CloseAllPositions()

// ✅ CORRECT - check first
positions, _ := sugar.GetOpenPositions()
fmt.Printf("About to close %d positions\n", len(positions))
// Show details, then close
```

### 2) No verification after close

```go
// ❌ WRONG - assuming success
sugar.CloseAllPositions()
fmt.Println("Done!")

// ✅ CORRECT - verify
sugar.CloseAllPositions()
time.Sleep(2 * time.Second)
remaining, _ := sugar.GetOpenPositions()
if len(remaining) > 0 {
    fmt.Printf("⚠️  %d positions still open!\n", len(remaining))
}
```

### 3) Using in loops/triggers without protection

```go
// ❌ WRONG - might trigger multiple times
if someCondition {
    sugar.CloseAllPositions() // Could run every tick!
}

// ✅ CORRECT - add flag
var alreadyClosed bool
if someCondition && !alreadyClosed {
    sugar.CloseAllPositions()
    alreadyClosed = true
}
```

---

## 💎 Pro Tips

1. **Always preview** - Check what you're closing first

2. **Record P/L** - Save profit/loss before closing

3. **Verify closure** - Check no positions remain after

4. **Use sparingly** - This is an emergency tool

5. **Add confirmation** - Especially in production systems

---

## 🚨 Production Warning

**In production:**

- Add multiple confirmations before closing all
- Log all positions before closing
- Implement rate limiting (don't allow rapid repeated calls)
- Consider user confirmation for manual systems
- Have rollback plan if needed

---

**See also:** [`ClosePosition.md`](ClosePosition.md), [`GetOpenPositions.md`](GetOpenPositions.md)
