# 💰📅 Get Profit Today (`GetProfitToday`)

> **Sugar method:** Returns total realized profit/loss from today's closed positions.

**API Information:**

* **Method:** `sugar.GetProfitToday()`
* **Timeout:** 5 seconds
* **Returns:** Total profit as float64

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetProfitToday() (float64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates today) |

| Output | Type | Description |
|--------|------|-------------|
| `float64` | float | Total realized profit/loss from today |
| `error` | `error` | Error if calculation failed |

---

## 💬 Just the Essentials

* **What it is:** Get today's total profit/loss from closed positions (00:00 to now).
* **Why you need it:** Quick profit check, daily targets, performance dashboards.
* **Sanity check:** Returns 0 if no closed positions today. Positive = profit, negative = loss.

---

## 🎯 When to Use

✅ **Quick profit check** - Fast single number answer

✅ **Daily targets** - Check if target met

✅ **Dashboard display** - Show today's performance

✅ **Performance monitoring** - Track profit in real-time

---

## 🔗 Usage Examples

### 1) Basic usage - check today's profit

```go
profit, err := sugar.GetProfitToday()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Today's profit: $%.2f\n", profit)

if profit > 0 {
    fmt.Println("✅ Profitable day")
} else if profit < 0 {
    fmt.Println("❌ Losing day")
} else {
    fmt.Println("➡️  Breakeven")
}
```

---

### 2) Daily target tracking

```go
dailyTarget := 100.0 // $100 per day

profit, _ := sugar.GetProfitToday()

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      DAILY TARGET PROGRESS            ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Target:   $%.2f\n", dailyTarget)
fmt.Printf("Current:  $%.2f\n", profit)

progress := (profit / dailyTarget) * 100

if profit >= dailyTarget {
    fmt.Printf("🎯 Target achieved! (%.0f%%)\n", progress)
} else if progress >= 80 {
    remaining := dailyTarget - profit
    fmt.Printf("🟡 Close! ($%.2f remaining)\n", remaining)
} else {
    remaining := dailyTarget - profit
    fmt.Printf("📊 In progress ($%.2f remaining)\n", remaining)
}
```

---

### 3) Real-time profit monitor

```go
func MonitorDailyProfit(sugar *mt5.MT5Sugar, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        profit, _ := sugar.GetProfitToday()

        timeStr := time.Now().Format("15:04:05")

        status := "➡️"
        if profit > 0 {
            status = "✅"
        } else if profit < 0 {
            status = "❌"
        }

        fmt.Printf("[%s] Today's P/L: $%.2f %s\n",
            timeStr, profit, status)
    }
}

// Usage: Monitor every 30 seconds
go MonitorDailyProfit(sugar, 30*time.Second)
```

---

### 4) Profit limit check (stop trading rule)

```go
maxDailyLoss := -200.0 // Stop if lose $200

profit, _ := sugar.GetProfitToday()

fmt.Printf("Today's profit: $%.2f\n", profit)

if profit <= maxDailyLoss {
    fmt.Println("🛑 DAILY LOSS LIMIT REACHED")
    fmt.Println("   STOP TRADING FOR TODAY")
    return
}

if profit < maxDailyLoss*0.8 {
    fmt.Printf("⚠️  Warning: Approaching loss limit ($%.2f)\n",
        maxDailyLoss)
}
```

---

### 5) Hourly profit check

```go
hourlyTarget := 10.0 // $10 per hour

currentHour := time.Now().Hour()

profit, _ := sugar.GetProfitToday()

// Simple hourly pace calculation
hoursElapsed := float64(currentHour)
if hoursElapsed == 0 {
    hoursElapsed = 1 // Avoid division by zero
}

hourlyRate := profit / hoursElapsed

fmt.Printf("Today's Performance:\n")
fmt.Printf("  Total:       $%.2f\n", profit)
fmt.Printf("  Hours:       %d\n", currentHour)
fmt.Printf("  Hourly rate: $%.2f\n", hourlyRate)

if hourlyRate >= hourlyTarget {
    fmt.Println("✅ Above target pace")
} else {
    fmt.Println("⚠️  Below target pace")
}
```

---

### 6) Compare with yesterday

```go
todayProfit, _ := sugar.GetProfitToday()

// Get yesterday's deals and sum profit
yesterdayDeals, _ := sugar.GetDealsYesterday()
yesterdayProfit := 0.0
for _, deal := range yesterdayDeals {
    yesterdayProfit += deal.Profit
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      TODAY VS YESTERDAY               ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Yesterday: $%.2f\n", yesterdayProfit)
fmt.Printf("Today:     $%.2f\n", todayProfit)

diff := todayProfit - yesterdayProfit

if diff > 0 {
    fmt.Printf("Improvement: +$%.2f 📈\n", diff)
} else if diff < 0 {
    fmt.Printf("Decline: $%.2f 📉\n", -diff)
} else {
    fmt.Println("Same ➡️")
}
```

---

### 7) Dashboard widget

```go
func ShowProfitDashboard(sugar *mt5.MT5Sugar) {
    todayProfit, _ := sugar.GetProfitToday()
    balance, _ := sugar.GetBalance()
    equity, _ := sugar.GetEquity()

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║         PROFIT DASHBOARD              ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Balance:       $%.2f\n", balance)
    fmt.Printf("Equity:        $%.2f\n", equity)
    fmt.Printf("Today's P/L:   $%.2f\n", todayProfit)

    if todayProfit > 0 {
        pct := (todayProfit / balance) * 100
        fmt.Printf("Daily ROI:     +%.2f%%\n", pct)
    }
}

// Usage: Call periodically
ticker := time.NewTicker(1 * time.Minute)
for range ticker.C {
    ShowProfitDashboard(sugar)
    fmt.Println()
}
```

---

### 8) Profit rate alert system

```go
profitThresholds := []struct {
    Amount  float64
    Message string
}{
    {50, "🟢 $50 milestone"},
    {100, "✅ $100 milestone - daily target!"},
    {200, "🌟 $200 milestone - excellent!"},
    {300, "🎉 $300 milestone - outstanding!"},
}

previousProfit := 0.0

ticker := time.NewTicker(1 * time.Minute)
defer ticker.Stop()

for range ticker.C {
    currentProfit, _ := sugar.GetProfitToday()

    // Check if crossed any threshold
    for _, threshold := range profitThresholds {
        if previousProfit < threshold.Amount && currentProfit >= threshold.Amount {
            fmt.Printf("[%s] %s\n",
                time.Now().Format("15:04"), threshold.Message)
        }
    }

    previousProfit = currentProfit
}
```

---

### 9) Profit consistency checker

```go
func CheckProfitConsistency(sugar *mt5.MT5Sugar) {
    // Get today's profit
    todayProfit, _ := sugar.GetProfitToday()

    // Get last 5 days
    profitHistory := []float64{}

    for i := 1; i <= 5; i++ {
        from := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
        to := from.Add(24 * time.Hour)

        deals, _ := sugar.GetDealsDateRange(from, to)

        dayProfit := 0.0
        for _, deal := range deals {
            dayProfit += deal.Profit
        }

        profitHistory = append(profitHistory, dayProfit)
    }

    // Calculate average
    avgProfit := 0.0
    for _, p := range profitHistory {
        avgProfit += p
    }
    avgProfit /= float64(len(profitHistory))

    fmt.Println("Profit Consistency Check:")
    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Today:      $%.2f\n", todayProfit)
    fmt.Printf("5-day avg:  $%.2f\n", avgProfit)

    if todayProfit > avgProfit*1.2 {
        fmt.Println("🌟 Above average performance today!")
    } else if todayProfit < avgProfit*0.8 {
        fmt.Println("⚠️  Below average performance")
    } else {
        fmt.Println("✅ Consistent with recent performance")
    }
}

// Usage:
CheckProfitConsistency(sugar)
```

---

### 10) Advanced daily profit tracker

```go
type DailyProfitTracker struct {
    Target          float64
    MaxLoss         float64
    LastUpdate      time.Time
    LastProfit      float64
}

func NewDailyProfitTracker(target, maxLoss float64) *DailyProfitTracker {
    return &DailyProfitTracker{
        Target:  target,
        MaxLoss: maxLoss,
    }
}

func (dpt *DailyProfitTracker) Update(sugar *mt5.MT5Sugar) (bool, string) {
    currentProfit, _ := sugar.GetProfitToday()
    now := time.Now()

    // Check if profit changed
    if currentProfit != dpt.LastProfit {
        dpt.LastProfit = currentProfit
        dpt.LastUpdate = now
    }

    // Check target achieved
    if currentProfit >= dpt.Target {
        return false, fmt.Sprintf("🎯 Target achieved: $%.2f", currentProfit)
    }

    // Check max loss
    if currentProfit <= dpt.MaxLoss {
        return false, fmt.Sprintf("🛑 Stop trading: Max loss reached ($%.2f)", currentProfit)
    }

    // Continue trading
    return true, fmt.Sprintf("Current: $%.2f / $%.2f target", currentProfit, dpt.Target)
}

func (dpt *DailyProfitTracker) ShowStatus(sugar *mt5.MT5Sugar) {
    currentProfit, _ := sugar.GetProfitToday()

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      DAILY PROFIT TRACKER             ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Target:        $%.2f\n", dpt.Target)
    fmt.Printf("Max Loss:      $%.2f\n", dpt.MaxLoss)
    fmt.Printf("Current:       $%.2f\n", currentProfit)

    progress := (currentProfit / dpt.Target) * 100
    fmt.Printf("Progress:      %.1f%%\n\n", progress)

    remaining := dpt.Target - currentProfit
    riskUsed := (currentProfit / dpt.MaxLoss) * 100

    fmt.Printf("To target:     $%.2f\n", remaining)
    fmt.Printf("Risk used:     %.1f%%\n", riskUsed)
}

// Usage:
tracker := NewDailyProfitTracker(100, -200)

ticker := time.NewTicker(30 * time.Second)
for range ticker.C {
    canContinue, message := tracker.Update(sugar)

    fmt.Printf("[%s] %s\n", time.Now().Format("15:04"), message)

    if !canContinue {
        tracker.ShowStatus(sugar)
        break
    }
}
```

---

## 🔗 Related Methods

**🍬 Other time periods:**

* `GetProfitThisWeek()` - This week's profit total
* `GetProfitThisMonth()` - This month's profit total

**🍬 Detailed data:**

* `GetDealsToday()` - Get full deal information
* `GetProfit()` - Get floating profit from OPEN positions

---

## ⚠️ Common Pitfalls

### 1) Confusing with floating profit

```go
// ❌ WRONG - this is CLOSED positions only
todayProfit, _ := sugar.GetProfitToday()
// Does NOT include currently OPEN positions!

// ✅ CORRECT - for total including open positions
closedProfit, _ := sugar.GetProfitToday()  // Closed today
floatingProfit, _ := sugar.GetProfit()     // Open positions
totalProfit := closedProfit + floatingProfit
```

### 2) Not handling zero

```go
// ❌ WRONG - assuming non-zero means trades
profit, _ := sugar.GetProfitToday()
if profit != 0 {
    // Might be no trades, just checking if profitable
}

// ✅ CORRECT - check if there were deals
deals, _ := sugar.GetDealsToday()
if len(deals) > 0 {
    profit, _ := sugar.GetProfitToday()
    fmt.Printf("Profit from %d trades: $%.2f\n", len(deals), profit)
}
```

### 3) Timezone assumptions

```go
// ❌ WRONG - "today" is server's today

// ✅ CORRECT - aware that it's server time
profit, _ := sugar.GetProfitToday()
fmt.Println("Today's profit (server time): $%.2f", profit)
```

---

## 💎 Pro Tips

1. **Closed only** - Returns profit from CLOSED positions (not open)

2. **Zero OK** - Returns 0 if no closed positions (not an error)

3. **Fast** - Single number result, faster than `GetDealsToday()`

4. **Real-time safe** - Can call frequently for dashboards

5. **Server time** - "Today" based on MT5 server timezone

---

## 📊 Profit Sign Convention

```
Positive value: Net profit today
Zero: Breakeven (or no closed positions)
Negative value: Net loss today

Example:
+150.50 = $150.50 profit
0.00 = Breakeven or no trades
-75.25 = $75.25 loss
```

---

**See also:** [`GetDealsToday.md`](GetDealsToday.md), [`GetProfitThisWeek.md`](GetProfitThisWeek.md)
