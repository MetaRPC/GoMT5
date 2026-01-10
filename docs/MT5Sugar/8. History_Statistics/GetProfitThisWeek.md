# 💰📊 Get Profit This Week (`GetProfitThisWeek`)

> **Sugar method:** Returns total realized profit/loss from this week's closed positions (Monday to now).

**API Information:**

* **Method:** `sugar.GetProfitThisWeek()`
* **Timeout:** 5 seconds
* **Returns:** Total profit as float64

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetProfitThisWeek() (float64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates week range) |

| Output | Type | Description |
|--------|------|-------------|
| `float64` | float | Total realized profit/loss from this week |
| `error` | `error` | Error if calculation failed |

---

## 💬 Just the Essentials

* **What it is:** Get this week's total profit/loss from closed positions (Monday 00:00 to now).
* **Why you need it:** Weekly performance tracking, target monitoring, week-over-week comparison.
* **Sanity check:** Returns 0 if no closed positions this week. Week starts Monday.

---

## 🎯 When to Use

✅ **Weekly targets** - Check if weekly goal met

✅ **Performance tracking** - Monitor week-to-date profit

✅ **Dashboard display** - Show weekly performance

✅ **Trend analysis** - Track weekly performance over time

---

## 🔗 Usage Examples

### 1) Basic usage - check this week's profit

```go
profit, err := sugar.GetProfitThisWeek()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("This week's profit: $%.2f\n", profit)

if profit > 0 {
    fmt.Println("✅ Profitable week so far")
} else if profit < 0 {
    fmt.Println("❌ Losing week so far")
} else {
    fmt.Println("➡️  Breakeven")
}
```

---

### 2) Weekly target tracking

```go
weeklyTarget := 500.0 // $500 per week

profit, _ := sugar.GetProfitThisWeek()

progress := (profit / weeklyTarget) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      WEEKLY TARGET PROGRESS           ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Target:    $%.2f\n", weeklyTarget)
fmt.Printf("Current:   $%.2f\n", profit)
fmt.Printf("Progress:  %.1f%%\n\n", progress)

remaining := weeklyTarget - profit

if profit >= weeklyTarget {
    fmt.Printf("🎯 Target achieved! (+$%.2f)\n", profit-weeklyTarget)
} else if progress >= 80 {
    fmt.Printf("🟡 Almost there! ($%.2f remaining)\n", remaining)
} else {
    fmt.Printf("📊 In progress ($%.2f remaining)\n", remaining)
}
```

---

### 3) Compare with last week

```go
thisWeekProfit, _ := sugar.GetProfitThisWeek()

// Get last week's range
now := time.Now()
weekday := int(now.Weekday())
if weekday == 0 {
    weekday = 7
}

lastWeekStart := now.AddDate(0, 0, -(weekday + 6))
lastWeekEnd := now.AddDate(0, 0, -(weekday))

lastWeekDeals, _ := sugar.GetDealsDateRange(lastWeekStart, lastWeekEnd)

lastWeekProfit := 0.0
for _, deal := range lastWeekDeals {
    lastWeekProfit += deal.Profit
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      WEEK-OVER-WEEK COMPARISON        ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Last week:  $%.2f\n", lastWeekProfit)
fmt.Printf("This week:  $%.2f\n", thisWeekProfit)

diff := thisWeekProfit - lastWeekProfit

if diff > 0 {
    pct := (diff / lastWeekProfit) * 100
    fmt.Printf("Improvement: +$%.2f (%.1f%%) 📈\n", diff, pct)
} else if diff < 0 {
    pct := (-diff / lastWeekProfit) * 100
    fmt.Printf("Decline: $%.2f (-%.1f%%) 📉\n", -diff, pct)
} else {
    fmt.Println("No change ➡️")
}
```

---

### 4) Daily average this week

```go
profit, _ := sugar.GetProfitThisWeek()

now := time.Now()
weekday := int(now.Weekday())
if weekday == 0 {
    weekday = 7 // Sunday becomes 7
}

daysElapsed := weekday
avgPerDay := profit / float64(daysElapsed)

// Project to end of week (5 trading days)
projectedWeekEnd := avgPerDay * 5

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      WEEKLY PERFORMANCE ANALYSIS      ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Days elapsed:       %d\n", daysElapsed)
fmt.Printf("Week-to-date:       $%.2f\n", profit)
fmt.Printf("Daily average:      $%.2f\n", avgPerDay)
fmt.Printf("Week projection:    $%.2f\n", projectedWeekEnd)
```

---

### 5) Weekly profit limit check

```go
maxWeeklyLoss := -500.0 // Stop if lose $500 in a week

profit, _ := sugar.GetProfitThisWeek()

fmt.Printf("This week's profit: $%.2f\n", profit)

if profit <= maxWeeklyLoss {
    fmt.Println("🛑 WEEKLY LOSS LIMIT REACHED")
    fmt.Println("   REVIEW STRATEGY BEFORE CONTINUING")
    return
}

if profit < maxWeeklyLoss*0.8 {
    fmt.Printf("⚠️  Warning: Approaching weekly loss limit ($%.2f)\n",
        maxWeeklyLoss)
}
```

---

### 6) Weekly progress dashboard

```go
func ShowWeeklyDashboard(sugar *mt5.MT5Sugar) {
    profit, _ := sugar.GetProfitThisWeek()

    // Get account info
    balance, _ := sugar.GetBalance()
    equity, _ := sugar.GetEquity()

    // Get today's profit
    todayProfit, _ := sugar.GetProfitToday()

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      WEEKLY PERFORMANCE DASHBOARD     ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Balance:       $%.2f\n", balance)
    fmt.Printf("Equity:        $%.2f\n\n", equity)

    fmt.Printf("This week:     $%.2f\n", profit)
    fmt.Printf("Today:         $%.2f\n\n", todayProfit)

    if profit > 0 {
        weeklyROI := (profit / balance) * 100
        fmt.Printf("Weekly ROI:    +%.2f%%\n", weeklyROI)
    }
}

// Usage:
ShowWeeklyDashboard(sugar)
```

---

### 7) Multi-week tracking

```go
func ShowLastNWeeks(sugar *mt5.MT5Sugar, n int) {
    fmt.Println("Last Weeks Performance:")
    fmt.Println("─────────────────────────────────────────")

    now := time.Now()

    for i := 0; i < n; i++ {
        // Calculate week range
        offset := i * 7
        weekEnd := now.AddDate(0, 0, -offset)

        weekday := int(weekEnd.Weekday())
        if weekday == 0 {
            weekday = 7
        }

        weekStart := weekEnd.AddDate(0, 0, -(weekday - 1))
        weekStart = time.Date(weekStart.Year(), weekStart.Month(),
            weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

        if i == 0 {
            // Current week (partial)
            deals, _ := sugar.GetDealsThisWeek()
            profit := 0.0
            for _, deal := range deals {
                profit += deal.Profit
            }

            fmt.Printf("This week:   $%.2f (partial)\n", profit)
        } else {
            // Complete weeks
            weekEndTime := weekStart.AddDate(0, 0, 7).Add(-time.Second)
            deals, _ := sugar.GetDealsDateRange(weekStart, weekEndTime)

            profit := 0.0
            for _, deal := range deals {
                profit += deal.Profit
            }

            weekLabel := weekStart.Format("Jan 02")
            status := "✅"
            if profit < 0 {
                status = "❌"
            }

            fmt.Printf("%s Week of %s: $%.2f\n", status, weekLabel, profit)
        }
    }
}

// Usage: Show last 4 weeks
ShowLastNWeeks(sugar, 4)
```

---

### 8) Weekly milestone tracker

```go
weeklyMilestones := []struct {
    Amount  float64
    Message string
}{
    {100, "🟢 $100 milestone"},
    {250, "✅ $250 milestone"},
    {500, "🌟 $500 milestone - weekly target!"},
    {750, "🎉 $750 milestone - excellent week!"},
    {1000, "🏆 $1000 milestone - outstanding!"},
}

profit, _ := sugar.GetProfitThisWeek()

fmt.Printf("This week's profit: $%.2f\n\n", profit)

fmt.Println("Milestones:")
for _, milestone := range weeklyMilestones {
    if profit >= milestone.Amount {
        fmt.Printf("✅ %s\n", milestone.Message)
    } else {
        remaining := milestone.Amount - profit
        fmt.Printf("⬜ $%.2f to %s\n", remaining, milestone.Message)
    }
}
```

---

### 9) Weekly consistency rating

```go
profit, _ := sugar.GetProfitThisWeek()

// Get last 4 weeks for comparison
now := time.Now()
weeklyProfits := []float64{}

for i := 1; i <= 4; i++ {
    offset := i * 7
    weekEnd := now.AddDate(0, 0, -offset)

    weekday := int(weekEnd.Weekday())
    if weekday == 0 {
        weekday = 7
    }

    weekStart := weekEnd.AddDate(0, 0, -(weekday - 1))
    weekStart = time.Date(weekStart.Year(), weekStart.Month(),
        weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

    weekEndTime := weekStart.AddDate(0, 0, 7).Add(-time.Second)
    deals, _ := sugar.GetDealsDateRange(weekStart, weekEndTime)

    weekProfit := 0.0
    for _, deal := range deals {
        weekProfit += deal.Profit
    }

    weeklyProfits = append(weeklyProfits, weekProfit)
}

// Calculate average
avgWeeklyProfit := 0.0
for _, p := range weeklyProfits {
    avgWeeklyProfit += p
}
avgWeeklyProfit /= float64(len(weeklyProfits))

fmt.Println("Weekly Consistency Analysis:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("This week:   $%.2f\n", profit)
fmt.Printf("4-week avg:  $%.2f\n\n", avgWeeklyProfit)

if profit > avgWeeklyProfit*1.2 {
    fmt.Println("🌟 Exceptional week - above average!")
} else if profit > avgWeeklyProfit*0.8 {
    fmt.Println("✅ Consistent with recent weeks")
} else {
    fmt.Println("⚠️  Below recent average")
}
```

---

### 10) Advanced weekly profit manager

```go
type WeeklyProfitManager struct {
    Target      float64
    MaxLoss     float64
    MinDaily    float64
}

func NewWeeklyProfitManager(target, maxLoss, minDaily float64) *WeeklyProfitManager {
    return &WeeklyProfitManager{
        Target:   target,
        MaxLoss:  maxLoss,
        MinDaily: minDaily,
    }
}

func (wpm *WeeklyProfitManager) Assess(sugar *mt5.MT5Sugar) {
    weekProfit, _ := sugar.GetProfitThisWeek()
    todayProfit, _ := sugar.GetProfitToday()

    now := time.Now()
    weekday := int(now.Weekday())
    if weekday == 0 {
        weekday = 7
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      WEEKLY PROFIT ASSESSMENT         ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Weekly target:     $%.2f\n", wpm.Target)
    fmt.Printf("Current profit:    $%.2f\n", weekProfit)
    fmt.Printf("Today's profit:    $%.2f\n\n", todayProfit)

    // Progress
    progress := (weekProfit / wpm.Target) * 100
    fmt.Printf("Progress:          %.1f%%\n\n", progress)

    // Assessment
    fmt.Println("Assessment:")

    if weekProfit >= wpm.Target {
        fmt.Println("  🎯 Weekly target achieved!")
    } else {
        remaining := wpm.Target - weekProfit
        daysLeft := 5 - weekday // Assume 5 trading days
        if daysLeft > 0 {
            neededPerDay := remaining / float64(daysLeft)
            fmt.Printf("  📊 $%.2f remaining ($%.2f/day needed)\n",
                remaining, neededPerDay)
        }
    }

    if weekProfit <= wpm.MaxLoss {
        fmt.Println("  🛑 Weekly loss limit reached")
    }

    if todayProfit < wpm.MinDaily {
        fmt.Printf("  ⚠️  Today below minimum daily ($%.2f)\n", wpm.MinDaily)
    }

    // Daily average
    avgDaily := weekProfit / float64(weekday)
    fmt.Printf("\nDaily average:     $%.2f\n", avgDaily)

    if avgDaily >= wpm.MinDaily {
        fmt.Println("✅ Good daily pace")
    } else {
        fmt.Println("⚠️  Below target pace")
    }
}

// Usage:
manager := NewWeeklyProfitManager(500, -500, 50)
manager.Assess(sugar)
```

---

## 🔗 Related Methods

**🍬 Other time periods:**

* `GetProfitToday()` - Today's profit total
* `GetProfitThisMonth()` - This month's profit total

**🍬 Detailed data:**

* `GetDealsThisWeek()` - Get full week's deal information

---

## ⚠️ Common Pitfalls

### 1) Confusing with floating profit

```go
// ❌ WRONG - this is CLOSED positions only
weekProfit, _ := sugar.GetProfitThisWeek()
// Does NOT include currently OPEN positions!

// ✅ CORRECT - for total including open
closedProfit, _ := sugar.GetProfitThisWeek()
floatingProfit, _ := sugar.GetProfit()
totalProfit := closedProfit + floatingProfit
```

### 2) Week starts Monday

```go
// ❌ WRONG - assuming week starts Sunday

// ✅ CORRECT - MT5 week starts Monday
profit, _ := sugar.GetProfitThisWeek()
fmt.Println("Monday to now profit: $%.2f", profit)
```

### 3) Incomplete week comparisons

```go
// ❌ WRONG - comparing mid-week to full week
thisWeekProfit, _ := sugar.GetProfitThisWeek()  // Mon-Wed
lastWeekProfit := 500.0                         // Full week
// Not fair comparison!

// ✅ CORRECT - note the difference
fmt.Println("This week (partial) vs last week (complete)")
```

---

## 💎 Pro Tips

1. **Week starts Monday** - MT5 uses Monday as first day of week

2. **Closed only** - Returns profit from CLOSED positions (not open)

3. **Zero OK** - Returns 0 if no closed positions (not error)

4. **Partial data** - Mid-week calls return incomplete week

5. **Server time** - Based on MT5 server timezone

---

## 📊 Week Range

```
Week range calculation:
now = time.Now()
weekday = int(now.Weekday())
if weekday == 0 { weekday = 7 }  // Sunday becomes 7

startOfWeek = now - (weekday - 1) days at 00:00:00

Returns profit from Monday 00:00 to now
```

---

**See also:** [`GetDealsThisWeek.md`](GetDealsThisWeek.md), [`GetProfitToday.md`](GetProfitToday.md), [`GetProfitThisMonth.md`](GetProfitThisMonth.md)
