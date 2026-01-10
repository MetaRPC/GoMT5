# 💰📆 Get Profit This Month (`GetProfitThisMonth`)

> **Sugar method:** Returns total realized profit/loss from this month's closed positions (1st day to now).

**API Information:**

* **Method:** `sugar.GetProfitThisMonth()`
* **Timeout:** 5 seconds
* **Returns:** Total profit as float64

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetProfitThisMonth() (float64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates month range) |

| Output | Type | Description |
|--------|------|-------------|
| `float64` | float | Total realized profit/loss from this month |
| `error` | `error` | Error if calculation failed |

---

## 💬 Just the Essentials

* **What it is:** Get this month's total profit/loss from closed positions (1st day 00:00 to now).
* **Why you need it:** Monthly performance tracking, goal monitoring, month-over-month comparison.
* **Sanity check:** Returns 0 if no closed positions this month. Month starts on 1st.

---

## 🎯 When to Use

✅ **Monthly targets** - Check if monthly goal met

✅ **Performance tracking** - Monitor month-to-date profit

✅ **Dashboard display** - Show monthly performance

✅ **Reporting** - Generate monthly profit summaries

---

## 🔗 Usage Examples

### 1) Basic usage - check this month's profit

```go
profit, err := sugar.GetProfitThisMonth()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

monthName := time.Now().Format("January")

fmt.Printf("%s profit: $%.2f\n", monthName, profit)

if profit > 0 {
    fmt.Println("✅ Profitable month so far")
} else if profit < 0 {
    fmt.Println("❌ Losing month so far")
} else {
    fmt.Println("➡️  Breakeven")
}
```

---

### 2) Monthly target tracking

```go
monthlyTarget := 2000.0 // $2000 per month

profit, _ := sugar.GetProfitThisMonth()

progress := (profit / monthlyTarget) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MONTHLY TARGET PROGRESS          ║")
fmt.Println("╚═══════════════════════════════════════╝")

monthName := time.Now().Format("January 2006")
fmt.Printf("Month:      %s\n\n", monthName)

fmt.Printf("Target:     $%.2f\n", monthlyTarget)
fmt.Printf("Current:    $%.2f\n", profit)
fmt.Printf("Progress:   %.1f%%\n\n", progress)

// Calculate if on track
dayOfMonth := time.Now().Day()
daysInMonth := time.Date(time.Now().Year(), time.Now().Month()+1, 0,
    0, 0, 0, 0, time.Now().Location()).Day()

monthProgress := float64(dayOfMonth) / float64(daysInMonth) * 100

if profit >= monthlyTarget {
    fmt.Printf("🎯 Target achieved! (+$%.2f)\n", profit-monthlyTarget)
} else {
    remaining := monthlyTarget - profit
    fmt.Printf("Remaining:  $%.2f\n", remaining)

    projected := (profit / float64(dayOfMonth)) * float64(daysInMonth)
    fmt.Printf("Projected:  $%.2f\n\n", projected)

    if projected >= monthlyTarget {
        fmt.Println("✅ On track to meet target")
    } else {
        shortfall := monthlyTarget - projected
        fmt.Printf("⚠️  Projected shortfall: $%.2f\n", shortfall)
    }
}
```

---

### 3) Month-over-month comparison

```go
thisMonthProfit, _ := sugar.GetProfitThisMonth()

// Get last month
now := time.Now()
firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)

lastMonthDeals, _ := sugar.GetDealsDateRange(
    firstOfLastMonth,
    firstOfThisMonth.Add(-time.Second),
)

lastMonthProfit := 0.0
for _, deal := range lastMonthDeals {
    lastMonthProfit += deal.Profit
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MONTH-OVER-MONTH COMPARISON      ║")
fmt.Println("╚═══════════════════════════════════════╝")

lastMonthName := firstOfLastMonth.Format("January")
thisMonthName := now.Format("January")

fmt.Printf("%s:  $%.2f\n", lastMonthName, lastMonthProfit)
fmt.Printf("%s:  $%.2f\n", thisMonthName, thisMonthProfit)

diff := thisMonthProfit - lastMonthProfit

if diff > 0 {
    if lastMonthProfit != 0 {
        pct := (diff / lastMonthProfit) * 100
        fmt.Printf("\nImprovement: +$%.2f (%.1f%%) 📈\n", diff, pct)
    } else {
        fmt.Printf("\nImprovement: +$%.2f 📈\n", diff)
    }
} else if diff < 0 {
    if lastMonthProfit != 0 {
        pct := (-diff / lastMonthProfit) * 100
        fmt.Printf("\nDecline: $%.2f (-%.1f%%) 📉\n", -diff, pct)
    } else {
        fmt.Printf("\nDecline: $%.2f 📉\n", -diff)
    }
} else {
    fmt.Println("\nNo change ➡️")
}
```

---

### 4) Daily average this month

```go
profit, _ := sugar.GetProfitThisMonth()

now := time.Now()
daysElapsed := now.Day()
daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()

avgPerDay := profit / float64(daysElapsed)
projectedMonthEnd := avgPerDay * float64(daysInMonth)

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MONTHLY PERFORMANCE ANALYSIS     ║")
fmt.Println("╚═══════════════════════════════════════╝")

monthName := now.Format("January 2006")
fmt.Printf("Month:              %s\n\n", monthName)

fmt.Printf("Days elapsed:       %d / %d\n", daysElapsed, daysInMonth)
fmt.Printf("Month-to-date:      $%.2f\n", profit)
fmt.Printf("Daily average:      $%.2f\n", avgPerDay)
fmt.Printf("Month projection:   $%.2f\n", projectedMonthEnd)
```

---

### 5) Monthly profit limit check

```go
maxMonthlyLoss := -1000.0 // Stop if lose $1000 in a month

profit, _ := sugar.GetProfitThisMonth()

monthName := time.Now().Format("January")

fmt.Printf("%s profit: $%.2f\n", monthName, profit)

if profit <= maxMonthlyLoss {
    fmt.Println("🛑 MONTHLY LOSS LIMIT REACHED")
    fmt.Println("   STOP TRADING AND REVIEW STRATEGY")
    return
}

if profit < maxMonthlyLoss*0.8 {
    fmt.Printf("⚠️  Warning: Approaching monthly loss limit ($%.2f)\n",
        maxMonthlyLoss)
}
```

---

### 6) Monthly dashboard

```go
func ShowMonthlyDashboard(sugar *mt5.MT5Sugar) {
    monthProfit, _ := sugar.GetProfitThisMonth()
    weekProfit, _ := sugar.GetProfitThisWeek()
    todayProfit, _ := sugar.GetProfitToday()

    balance, _ := sugar.GetBalance()
    equity, _ := sugar.GetEquity()

    now := time.Now()
    monthName := now.Format("January 2006")

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      MONTHLY PERFORMANCE DASHBOARD    ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Month:         %s\n\n", monthName)

    fmt.Printf("Balance:       $%.2f\n", balance)
    fmt.Printf("Equity:        $%.2f\n\n", equity)

    fmt.Printf("This month:    $%.2f\n", monthProfit)
    fmt.Printf("This week:     $%.2f\n", weekProfit)
    fmt.Printf("Today:         $%.2f\n\n", todayProfit)

    if monthProfit > 0 {
        monthlyROI := (monthProfit / balance) * 100
        fmt.Printf("Monthly ROI:   +%.2f%%\n", monthlyROI)
    }
}

// Usage:
ShowMonthlyDashboard(sugar)
```

---

### 7) Last 6 months tracking

```go
func ShowLast6Months(sugar *mt5.MT5Sugar) {
    fmt.Println("Last 6 Months Performance:")
    fmt.Println("─────────────────────────────────────────")

    now := time.Now()

    for i := 5; i >= 0; i-- {
        monthStart := time.Date(now.Year(), now.Month()-time.Month(i), 1,
            0, 0, 0, 0, now.Location())

        monthEnd := monthStart.AddDate(0, 1, 0)

        if monthEnd.After(now) {
            monthEnd = now
        } else {
            monthEnd = monthEnd.Add(-time.Second)
        }

        deals, _ := sugar.GetDealsDateRange(monthStart, monthEnd)

        profit := 0.0
        for _, deal := range deals {
            profit += deal.Profit
        }

        monthName := monthStart.Format("Jan 2006")
        status := "✅"
        if profit < 0 {
            status = "❌"
        }

        if i == 0 {
            fmt.Printf("%s %-12s: $%.2f (current)\n", status, monthName, profit)
        } else {
            fmt.Printf("%s %-12s: $%.2f\n", status, monthName, profit)
        }
    }
}

// Usage:
ShowLast6Months(sugar)
```

---

### 8) Monthly milestone tracker

```go
monthlyMilestones := []struct {
    Amount  float64
    Message string
}{
    {500, "🟢 $500 milestone"},
    {1000, "✅ $1,000 milestone"},
    {2000, "🌟 $2,000 milestone - monthly target!"},
    {3000, "🎉 $3,000 milestone - excellent!"},
    {5000, "🏆 $5,000 milestone - outstanding!"},
}

profit, _ := sugar.GetProfitThisMonth()

monthName := time.Now().Format("January")

fmt.Printf("%s profit: $%.2f\n\n", monthName, profit)

fmt.Println("Milestones:")
for _, milestone := range monthlyMilestones {
    if profit >= milestone.Amount {
        fmt.Printf("✅ %s\n", milestone.Message)
    } else {
        remaining := milestone.Amount - profit
        fmt.Printf("⬜ $%.2f to %s\n", remaining, milestone.Message)
    }
}
```

---

### 9) Monthly consistency rating

```go
profit, _ := sugar.GetProfitThisMonth()

// Get last 6 months for comparison
now := time.Now()
monthlyProfits := []float64{}

for i := 1; i <= 6; i++ {
    monthStart := time.Date(now.Year(), now.Month()-time.Month(i), 1,
        0, 0, 0, 0, now.Location())
    monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

    deals, _ := sugar.GetDealsDateRange(monthStart, monthEnd)

    monthProfit := 0.0
    for _, deal := range deals {
        monthProfit += deal.Profit
    }

    monthlyProfits = append(monthlyProfits, monthProfit)
}

// Calculate average
avgMonthlyProfit := 0.0
for _, p := range monthlyProfits {
    avgMonthlyProfit += p
}
avgMonthlyProfit /= float64(len(monthlyProfits))

fmt.Println("Monthly Consistency Analysis:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("This month:   $%.2f\n", profit)
fmt.Printf("6-month avg:  $%.2f\n\n", avgMonthlyProfit)

if profit > avgMonthlyProfit*1.2 {
    fmt.Println("🌟 Exceptional month - well above average!")
} else if profit > avgMonthlyProfit*0.8 {
    fmt.Println("✅ Consistent with recent months")
} else {
    fmt.Println("⚠️  Below recent average - review needed")
}
```

---

### 10) Advanced monthly profit manager

```go
type MonthlyProfitManager struct {
    Target          float64
    MaxLoss         float64
    MinWeekly       float64
}

func NewMonthlyProfitManager(target, maxLoss, minWeekly float64) *MonthlyProfitManager {
    return &MonthlyProfitManager{
        Target:    target,
        MaxLoss:   maxLoss,
        MinWeekly: minWeekly,
    }
}

func (mpm *MonthlyProfitManager) Assess(sugar *mt5.MT5Sugar) {
    monthProfit, _ := sugar.GetProfitThisMonth()
    weekProfit, _ := sugar.GetProfitThisWeek()
    todayProfit, _ := sugar.GetProfitToday()

    now := time.Now()
    dayOfMonth := now.Day()
    daysInMonth := time.Date(now.Year(), now.Month()+1, 0,
        0, 0, 0, 0, now.Location()).Day()

    monthName := now.Format("January 2006")

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      MONTHLY PROFIT ASSESSMENT        ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Month:             %s\n\n", monthName)

    fmt.Printf("Monthly target:    $%.2f\n", mpm.Target)
    fmt.Printf("Current profit:    $%.2f\n", monthProfit)
    fmt.Printf("This week:         $%.2f\n", weekProfit)
    fmt.Printf("Today:             $%.2f\n\n", todayProfit)

    // Progress
    progress := (monthProfit / mpm.Target) * 100
    fmt.Printf("Progress:          %.1f%%\n\n", progress)

    // Assessment
    fmt.Println("Assessment:")

    if monthProfit >= mpm.Target {
        fmt.Println("  🎯 Monthly target achieved!")
    } else {
        remaining := mpm.Target - monthProfit
        daysLeft := daysInMonth - dayOfMonth

        if daysLeft > 0 {
            neededPerDay := remaining / float64(daysLeft)
            fmt.Printf("  📊 $%.2f remaining ($%.2f/day needed)\n",
                remaining, neededPerDay)
        }
    }

    if monthProfit <= mpm.MaxLoss {
        fmt.Println("  🛑 Monthly loss limit reached")
    }

    if weekProfit < mpm.MinWeekly {
        fmt.Printf("  ⚠️  This week below minimum ($%.2f)\n", mpm.MinWeekly)
    }

    // Daily/weekly averages
    avgDaily := monthProfit / float64(dayOfMonth)
    fmt.Printf("\nDaily average:     $%.2f\n", avgDaily)

    // Estimate weeks
    weeksInMonth := float64(daysInMonth) / 7.0
    avgWeekly := monthProfit / (float64(dayOfMonth) / 7.0)
    fmt.Printf("Weekly average:    $%.2f\n", avgWeekly)

    if avgWeekly >= mpm.MinWeekly {
        fmt.Println("✅ Good weekly pace")
    } else {
        fmt.Println("⚠️  Below target weekly pace")
    }

    // Projection
    projected := avgDaily * float64(daysInMonth)
    fmt.Printf("\nProjected month:   $%.2f\n", projected)

    if projected >= mpm.Target {
        fmt.Println("✅ On track for target")
    } else {
        shortfall := mpm.Target - projected
        fmt.Printf("⚠️  Projected shortfall: $%.2f\n", shortfall)
    }
}

// Usage:
manager := NewMonthlyProfitManager(2000, -1000, 400)
manager.Assess(sugar)
```

---

## 🔗 Related Methods

**🍬 Other time periods:**

* `GetProfitToday()` - Today's profit total
* `GetProfitThisWeek()` - This week's profit total

**🍬 Detailed data:**

* `GetDealsThisMonth()` - Get full month's deal information

---

## ⚠️ Common Pitfalls

### 1) Confusing with floating profit

```go
// ❌ WRONG - this is CLOSED positions only
monthProfit, _ := sugar.GetProfitThisMonth()
// Does NOT include currently OPEN positions!

// ✅ CORRECT - for total including open
closedProfit, _ := sugar.GetProfitThisMonth()
floatingProfit, _ := sugar.GetProfit()
totalProfit := closedProfit + floatingProfit
```

### 2) Early month = incomplete data

```go
// ❌ WRONG - comparing Day 5 to full month
thisMonthProfit, _ := sugar.GetProfitThisMonth()  // Day 5
lastMonthProfit := 2000.0                         // Full 30 days
// Not fair comparison!

// ✅ CORRECT - adjust for days
dayOfMonth := time.Now().Day()
dailyAvg := thisMonthProfit / float64(dayOfMonth)
fmt.Printf("Daily average: $%.2f\n", dailyAvg)
```

### 3) Not handling month boundaries

```go
// ❌ WRONG - month changes at midnight
// Calling at 23:59 on last day of month

// ✅ CORRECT - aware of month transitions
monthProfit, _ := sugar.GetProfitThisMonth()
monthName := time.Now().Format("January")
fmt.Printf("%s profit: $%.2f\n", monthName, monthProfit)
```

---

## 💎 Pro Tips

1. **Month starts 1st** - Always starts from 1st day 00:00:00

2. **Closed only** - Returns profit from CLOSED positions (not open)

3. **Zero OK** - Returns 0 if no closed positions (not error)

4. **Partial data** - Mid-month calls return incomplete month

5. **Server time** - Based on MT5 server timezone

---

## 📊 Month Range

```
Month range calculation:
now = time.Now()
startOfMonth = Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

Returns profit from 1st day 00:00 to now
```

---

**See also:** [`GetDealsThisMonth.md`](GetDealsThisMonth.md), [`GetProfitToday.md`](GetProfitToday.md), [`GetProfitThisWeek.md`](GetProfitThisWeek.md)
