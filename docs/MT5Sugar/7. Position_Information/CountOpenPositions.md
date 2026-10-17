# 🔢 Count Open Positions (`CountOpenPositions`)

> **Sugar method:** Returns the exact count of open positions (fastest way to get count).

**API Information:**

* **Method:** `sugar.CountOpenPositions()`
* **Timeout:** 3 seconds
* **Returns:** Integer count

---

## 📋 Method Signature

```go
func (s *MT5Sugar) CountOpenPositions() (int, error)
```

It calls `MT5Service.GetPositionsTotal(ctx)` (MT5 `PositionsTotal`) and converts the `int32` result to `int`.
Pending orders are not counted.

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters required |

| Output | Type | Description |
|--------|------|-------------|
| `int` | integer | Number of open positions |
| `error` | `error` | `"CountOpenPositions failed: ..."` if the call failed (count is 0 then) |

---

## 💬 Just the Essentials

* **What it is:** Get count of open positions without fetching any position data.
* **Why you need it:** Fastest way to count - more efficient than `len(GetOpenPositions())`.
* **Sanity check:** Returns 0 if no positions. Most efficient counting method.

---

## 🎯 When to Use

✅ **Quick count** - Fast integer answer without fetching data

✅ **Validation** - Check position count limits

✅ **Dashboards** - Display count efficiently

✅ **Monitoring** - Track position count over time

---

## 🔗 Usage Examples

### 1) Basic usage - get count

```go
count, err := sugar.CountOpenPositions()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Open positions: %d\n", count)

if count == 0 {
    fmt.Println("No positions - ready to trade")
} else {
    fmt.Printf("Currently managing %d positions\n", count)
}
```

---

### 2) Position limit enforcement

```go
maxPositions := 10

count, _ := sugar.CountOpenPositions()

fmt.Printf("Position count: %d / %d max\n", count, maxPositions)

if count >= maxPositions {
    fmt.Println("❌ Position limit reached - cannot open new trades")
    return
}

remaining := maxPositions - count
fmt.Printf("✅ Can open %d more positions\n", remaining)
```

---

### 3) Before opening position

```go
maxSimultaneous := 5

// Check current count
count, _ := sugar.CountOpenPositions()

if count >= maxSimultaneous {
    fmt.Printf("⚠️  Already at max (%d positions)\n", maxSimultaneous)
    fmt.Println("   Wait for position to close before opening new")
    return
}

fmt.Printf("Current: %d/%d - opening new position...\n", count, maxSimultaneous)
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)
fmt.Printf("✅ Opened #%d (now %d/%d)\n", ticket, count+1, maxSimultaneous)
```

---

### 4) Monitor position count changes

```go
func MonitorPositionCount(sugar *mt5.MT5Sugar, duration time.Duration) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    timeout := time.After(duration)

    previousCount := -1

    for {
        select {
        case <-timeout:
            return

        case <-ticker.C:
            count, _ := sugar.CountOpenPositions()

            if previousCount == -1 {
                fmt.Printf("Initial count: %d positions\n", count)
            } else if count > previousCount {
                fmt.Printf("📈 Count increased: %d → %d (+%d)\n",
                    previousCount, count, count-previousCount)
            } else if count < previousCount {
                fmt.Printf("📉 Count decreased: %d → %d (-%d)\n",
                    previousCount, count, previousCount-count)
            } else {
                fmt.Printf("➡️  Count unchanged: %d\n", count)
            }

            previousCount = count
        }
    }
}

// Usage: Monitor for 5 minutes
MonitorPositionCount(sugar, 5*time.Minute)
```

---

### 5) Dashboard counter

```go
func ShowDashboard(sugar *mt5.MT5Sugar) {
    count, _ := sugar.CountOpenPositions()
    balance, _ := sugar.GetBalance()
    equity, _ := sugar.GetEquity()
    profit, _ := sugar.GetProfit()

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║         ACCOUNT DASHBOARD             ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Positions:  %d\n", count)
    fmt.Printf("Balance:    $%.2f\n", balance)
    fmt.Printf("Equity:     $%.2f\n", equity)
    fmt.Printf("Profit:     $%.2f\n", profit)

    if count > 0 {
        avgProfitPerPosition := profit / float64(count)
        fmt.Printf("Avg P/L:    $%.2f per position\n", avgProfitPerPosition)
    }
}

// Usage:
ticker := time.NewTicker(5 * time.Second)
for range ticker.C {
    ShowDashboard(sugar)
    fmt.Println()
}
```

---

### 6) Position count alert

```go
warnThreshold := 7
maxThreshold := 10

ticker := time.NewTicker(30 * time.Second)
defer ticker.Stop()

for range ticker.C {
    count, _ := sugar.CountOpenPositions()

    fmt.Printf("[%s] Positions: %d ",
        time.Now().Format("15:04:05"), count)

    if count >= maxThreshold {
        fmt.Println("🔴 MAX REACHED")
    } else if count >= warnThreshold {
        fmt.Println("🟡 WARNING")
    } else {
        fmt.Println("🟢 OK")
    }
}
```

---

### 7) Compare position count with target

```go
targetPositions := 3

count, _ := sugar.CountOpenPositions()

fmt.Printf("Target: %d positions\n", targetPositions)
fmt.Printf("Current: %d positions\n", count)

diff := count - targetPositions

if diff == 0 {
    fmt.Println("✅ At target")
} else if diff > 0 {
    fmt.Printf("⚠️  %d over target\n", diff)
} else {
    fmt.Printf("📉 %d below target\n", -diff)
}
```

---

### 8) Wait for all positions to close

```go
fmt.Println("Waiting for all positions to close...")

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

timeout := time.After(10 * time.Minute)

for {
    select {
    case <-timeout:
        count, _ := sugar.CountOpenPositions()
        fmt.Printf("⏰ Timeout - still %d positions open\n", count)
        return

    case <-ticker.C:
        count, _ := sugar.CountOpenPositions()

        if count == 0 {
            fmt.Println("✅ All positions closed!")
            return
        }

        fmt.Printf("%d positions remaining...\n", count)
    }
}
```

---

### 9) Position count statistics

```go
func TrackPositionCountStats(sugar *mt5.MT5Sugar, duration time.Duration) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    timeout := time.After(duration)

    minCount := 999999
    maxCount := 0
    totalCount := 0
    samples := 0

    for {
        select {
        case <-timeout:
            avgCount := float64(totalCount) / float64(samples)

            fmt.Println("\n╔═══════════════════════════════════════╗")
            fmt.Println("║    POSITION COUNT STATISTICS          ║")
            fmt.Println("╚═══════════════════════════════════════╝")
            fmt.Printf("Duration:   %v\n", duration)
            fmt.Printf("Samples:    %d\n", samples)
            fmt.Printf("Minimum:    %d positions\n", minCount)
            fmt.Printf("Maximum:    %d positions\n", maxCount)
            fmt.Printf("Average:    %.1f positions\n", avgCount)
            return

        case <-ticker.C:
            count, _ := sugar.CountOpenPositions()

            if count < minCount {
                minCount = count
            }
            if count > maxCount {
                maxCount = count
            }

            totalCount += count
            samples++

            fmt.Printf("Count: %d (min: %d, max: %d, avg: %.1f)\n",
                count, minCount, maxCount, float64(totalCount)/float64(samples))
        }
    }
}

// Usage: Track for 1 hour
TrackPositionCountStats(sugar, 1*time.Hour)
```

---

### 10) Advanced position count manager

```go
type PositionCountManager struct {
    MinPositions int
    MaxPositions int
    TargetPositions int
}

func (pcm *PositionCountManager) GetStatus(sugar *mt5.MT5Sugar) string {
    count, err := sugar.CountOpenPositions()
    if err != nil {
        return "ERROR"
    }

    if count < pcm.MinPositions {
        return "UNDEREXPOSED"
    } else if count > pcm.MaxPositions {
        return "OVEREXPOSED"
    } else if count == pcm.TargetPositions {
        return "OPTIMAL"
    } else {
        return "NORMAL"
    }
}

func (pcm *PositionCountManager) ShowReport(sugar *mt5.MT5Sugar) {
    count, _ := sugar.CountOpenPositions()
    status := pcm.GetStatus(sugar)

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║    POSITION COUNT MANAGER             ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Current:  %d positions\n", count)
    fmt.Printf("Target:   %d positions\n", pcm.TargetPositions)
    fmt.Printf("Range:    %d - %d\n", pcm.MinPositions, pcm.MaxPositions)
    fmt.Printf("Status:   %s\n\n", status)

    // Show recommendations
    switch status {
    case "UNDEREXPOSED":
        needed := pcm.MinPositions - count
        fmt.Printf("💡 Recommendation: Open %d more positions\n", needed)

    case "OVEREXPOSED":
        excess := count - pcm.MaxPositions
        fmt.Printf("⚠️  Warning: Close %d positions\n", excess)

    case "OPTIMAL":
        fmt.Println("✅ Perfect - at target")

    case "NORMAL":
        fmt.Println("✅ Within acceptable range")
    }
}

func (pcm *PositionCountManager) CanOpenNew(sugar *mt5.MT5Sugar) bool {
    count, _ := sugar.CountOpenPositions()
    return count < pcm.MaxPositions
}

// Usage:
manager := &PositionCountManager{
    MinPositions:    2,
    TargetPositions: 5,
    MaxPositions:    10,
}

manager.ShowReport(sugar)

if manager.CanOpenNew(sugar) {
    fmt.Println("\n✅ Approved to open new position")
} else {
    fmt.Println("\n❌ Cannot open - at maximum")
}
```

---

## 🔗 Related Methods

**🍬 Symbol-specific:**

* There is no per-symbol count method - use `len(sugar.GetPositionsBySymbol(symbol))`

**🍬 Related checks:**

* `HasOpenPosition(symbol)` - Boolean check for one symbol (fetches the full position list)
* `sugar.GetService().GetOpenedTickets(ctx)` - Ticket numbers only (position tickets, pending order tickets)
* `GetOpenPositions()` - Full data

---

## ⚠️ Common Pitfalls

### 1) Using GetOpenPositions() when you only need count

```go
// ❌ WRONG - fetching all data just to count
positions, _ := sugar.GetOpenPositions()
count := len(positions) // Wasteful!

// ✅ CORRECT - dedicated count method
count, _ := sugar.CountOpenPositions() // Much faster!
```

### 2) Not handling errors

```go
// ❌ WRONG - ignoring errors
count, _ := sugar.CountOpenPositions()
if count >= 10 {
    // Might be wrong if error occurred!
}

// ✅ CORRECT - check errors
count, err := sugar.CountOpenPositions()
if err != nil {
    fmt.Printf("Count failed: %v\n", err)
    return
}
```

### 3) Assuming count means specific positions exist

```go
// ❌ WRONG - count doesn't tell you WHICH positions
count, _ := sugar.CountOpenPositions()
if count > 0 {
    // Can't assume anything about the positions themselves
}

// ✅ CORRECT - fetch positions if you need details
count, _ := sugar.CountOpenPositions()
if count > 0 {
    positions, _ := sugar.GetOpenPositions()
    // Now work with actual positions
}
```

---

## 💎 Pro Tips

1. **Fastest count** - Use this instead of `len(GetOpenPositions())`

2. **Limits** - Perfect for enforcing position count limits

3. **Monitoring** - Track count changes over time

4. **Dashboards** - Display count without overhead

5. **Per symbol** - Use `HasOpenPosition(symbol)` for a per-symbol boolean, this for the total count

---

## 🔍 Payload Comparison

From smallest to largest response:
```
1. CountOpenPositions()                  - Integer count ⭐ YOU ARE HERE (PositionsTotal)
2. GetService().GetOpenedTickets(ctx)    - Ticket numbers only
3. GetOpenPositions()                    - Full position data
   HasOpenPosition(symbol),
   GetPositionsBySymbol(symbol)          - call GetOpenPositions() and filter in Go
```

Use `CountOpenPositions()` when you:

- Need exact count
- Don't need position details
- Want to enforce limits
- Track count over time

---

**See also:** [`HasOpenPosition.md`](HasOpenPosition.md), [`GetOpenPositions.md`](../6.%20Position_Management/GetOpenPositions.md)
