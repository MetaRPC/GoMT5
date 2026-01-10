# 📅📆 Get Deals This Month (`GetDealsThisMonth`)

> **Sugar method:** Returns all closed positions (deals) from this month (1st day 00:00 to now).

**API Information:**

* **Method:** `sugar.GetDealsThisMonth()`
* **Timeout:** 5 seconds
* **Returns:** Slice of position history info

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetDealsThisMonth() ([]*pb.PositionHistoryInfo, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates month range) |

| Output | Type | Description |
|--------|------|-------------|
| `[]*pb.PositionHistoryInfo` | slice | All closed positions from this month |
| `error` | `error` | Error if retrieval failed |

---

## 💬 Just the Essentials

* **What it is:** Get all closed trades from 1st day of month 00:00 to current time.
* **Why you need it:** Monthly performance tracking, strategy evaluation, profit/loss reports.
* **Sanity check:** Returns empty slice if no deals this month. Auto-calculates month range.

---

## 🎯 When to Use

✅ **Monthly reports** - Generate month-to-date summaries

✅ **Performance tracking** - Monitor monthly profit/loss

✅ **Strategy analysis** - Review month's trading performance

✅ **Goal tracking** - Check progress toward monthly targets

---

## 🔗 Usage Examples

### 1) Basic usage - show this month's deals

```go
deals, err := sugar.GetDealsThisMonth()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if len(deals) == 0 {
    fmt.Println("No deals this month")
    return
}

fmt.Printf("This month's deals: %d\n\n", len(deals))

totalProfit := 0.0
for i, deal := range deals {
    fmt.Printf("%d. #%d %s: $%.2f\n",
        i+1, deal.Ticket, deal.Symbol, deal.Profit)
    totalProfit += deal.Profit
}

fmt.Printf("\nMonth total: $%.2f\n", totalProfit)
```

---

### 2) Monthly performance summary

```go
deals, _ := sugar.GetDealsThisMonth()

if len(deals) == 0 {
    fmt.Println("No trading this month")
    return
}

totalProfit := 0.0
winCount := 0
lossCount := 0

for _, deal := range deals {
    totalProfit += deal.Profit
    if deal.Profit > 0 {
        winCount++
    } else if deal.Profit < 0 {
        lossCount++
    }
}

winRate := float64(winCount) / float64(len(deals)) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      THIS MONTH'S PERFORMANCE         ║")
fmt.Println("╚═══════════════════════════════════════╝")

monthName := time.Now().Format("January 2006")
fmt.Printf("Month:         %s\n\n", monthName)

fmt.Printf("Total trades:  %d\n", len(deals))
fmt.Printf("Winners:       %d (%.1f%%)\n", winCount, winRate)
fmt.Printf("Losers:        %d (%.1f%%)\n", lossCount, 100-winRate)
fmt.Printf("Total P/L:     $%.2f\n", totalProfit)

if totalProfit > 0 {
    fmt.Println("\n✅ Profitable month so far")
} else {
    fmt.Println("\n❌ Losing month so far")
}
```

---

### 3) Check monthly target progress

```go
monthlyTarget := 2000.0 // $2000 per month

deals, _ := sugar.GetDealsThisMonth()

if len(deals) == 0 {
    fmt.Printf("Monthly target: $%.2f\n", monthlyTarget)
    fmt.Println("No trades yet this month")
    return
}

actualProfit := 0.0
for _, deal := range deals {
    actualProfit += deal.Profit
}

progress := (actualProfit / monthlyTarget) * 100

// Calculate days into month
now := time.Now()
dayOfMonth := now.Day()
daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MONTHLY TARGET PROGRESS          ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Target:       $%.2f\n", monthlyTarget)
fmt.Printf("Current:      $%.2f\n", actualProfit)
fmt.Printf("Progress:     %.1f%%\n\n", progress)

fmt.Printf("Day %d of %d (%.0f%% of month)\n",
    dayOfMonth, daysInMonth,
    float64(dayOfMonth)/float64(daysInMonth)*100)

remaining := monthlyTarget - actualProfit

if actualProfit >= monthlyTarget {
    fmt.Printf("\n🎯 Target achieved! (+$%.2f)\n",
        actualProfit-monthlyTarget)
} else {
    fmt.Printf("\n📊 $%.2f remaining\n", remaining)

    // Project if on track
    projected := (actualProfit / float64(dayOfMonth)) * float64(daysInMonth)
    fmt.Printf("Projected month-end: $%.2f\n", projected)

    if projected >= monthlyTarget {
        fmt.Println("✅ On track to meet target")
    } else {
        fmt.Println("⚠️  Behind pace")
    }
}
```

---

### 4) Weekly breakdown of this month

```go
deals, _ := sugar.GetDealsThisMonth()

if len(deals) == 0 {
    fmt.Println("No deals this month")
    return
}

// Group by week number
weeklyStats := make(map[int]struct {
    Count  int
    Profit float64
})

for _, deal := range deals {
    _, week := deal.TimeClose.ISOWeek()
    stats := weeklyStats[week]
    stats.Count++
    stats.Profit += deal.Profit
    weeklyStats[week] = stats
}

fmt.Println("This Month - Weekly Breakdown:")
fmt.Println("─────────────────────────────────────────")

// Get current week range
_, currentWeek := time.Now().ISOWeek()

// Print weeks in order
for week := currentWeek - 3; week <= currentWeek; week++ {
    stats, exists := weeklyStats[week]
    if !exists {
        continue
    }

    status := "✅"
    if stats.Profit < 0 {
        status = "❌"
    }

    fmt.Printf("%s Week %d: %d trades, $%.2f\n",
        status, week, stats.Count, stats.Profit)
}
```

---

### 5) Top 5 most profitable trades this month

```go
deals, _ := sugar.GetDealsThisMonth()

if len(deals) == 0 {
    fmt.Println("No deals this month")
    return
}

// Sort by profit (descending)
sort.Slice(deals, func(i, j int) bool {
    return deals[i].Profit > deals[j].Profit
})

fmt.Println("Top 5 Most Profitable Trades This Month:")
fmt.Println("─────────────────────────────────────────")

count := 5
if len(deals) < 5 {
    count = len(deals)
}

for i := 0; i < count; i++ {
    deal := deals[i]
    fmt.Printf("%d. #%d %s %.2f lots: $%.2f\n",
        i+1, deal.Ticket, deal.Symbol,
        deal.Volume, deal.Profit)
}
```

---

### 6) Symbol performance this month

```go
deals, _ := sugar.GetDealsThisMonth()

if len(deals) == 0 {
    fmt.Println("No deals this month")
    return
}

symbolStats := make(map[string]struct {
    Count  int
    Profit float64
    Wins   int
})

for _, deal := range deals {
    stats := symbolStats[deal.Symbol]
    stats.Count++
    stats.Profit += deal.Profit
    if deal.Profit > 0 {
        stats.Wins++
    }
    symbolStats[deal.Symbol] = stats
}

fmt.Println("This Month - Symbol Performance:")
fmt.Println("─────────────────────────────────────────")

for symbol, stats := range symbolStats {
    winRate := float64(stats.Wins) / float64(stats.Count) * 100
    avgProfit := stats.Profit / float64(stats.Count)

    status := "✅"
    if stats.Profit < 0 {
        status = "❌"
    }

    fmt.Printf("%s %-8s: %d trades (%.0f%% win), $%.2f (avg: $%.2f)\n",
        status, symbol, stats.Count, winRate,
        stats.Profit, avgProfit)
}
```

---

### 7) Month-over-month comparison

```go
thisMonthDeals, _ := sugar.GetDealsThisMonth()

// Get last month's range
now := time.Now()
firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)

lastMonthDeals, _ := sugar.GetDealsDateRange(
    firstOfLastMonth,
    firstOfThisMonth.Add(-time.Second),
)

thisMonthProfit := 0.0
for _, deal := range thisMonthDeals {
    thisMonthProfit += deal.Profit
}

lastMonthProfit := 0.0
for _, deal := range lastMonthDeals {
    lastMonthProfit += deal.Profit
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MONTH-OVER-MONTH COMPARISON      ║")
fmt.Println("╚═══════════════════════════════════════╝")

lastMonthName := firstOfLastMonth.Format("January")
thisMonthName := now.Format("January")

fmt.Printf("%s:  %d trades, $%.2f\n",
    lastMonthName, len(lastMonthDeals), lastMonthProfit)
fmt.Printf("%s:  %d trades, $%.2f\n",
    thisMonthName, len(thisMonthDeals), thisMonthProfit)

diff := thisMonthProfit - lastMonthProfit

fmt.Printf("\nDifference: $%.2f ", diff)
if diff > 0 {
    if lastMonthProfit != 0 {
        pct := (diff / lastMonthProfit) * 100
        fmt.Printf("📈 (+%.1f%%)\n", pct)
    } else {
        fmt.Println("📈")
    }
} else if diff < 0 {
    if lastMonthProfit != 0 {
        pct := (-diff / lastMonthProfit) * 100
        fmt.Printf("📉 (-%.1f%%)\n", pct)
    } else {
        fmt.Println("📉")
    }
} else {
    fmt.Println("➡️  (No change)")
}
```

---

### 8) Daily average this month

```go
deals, _ := sugar.GetDealsThisMonth()

if len(deals) == 0 {
    fmt.Println("No deals this month")
    return
}

totalProfit := 0.0
for _, deal := range deals {
    totalProfit += deal.Profit
}

now := time.Now()
daysElapsed := now.Day()
daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()

avgPerDay := totalProfit / float64(daysElapsed)
projectedMonthEnd := avgPerDay * float64(daysInMonth)

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      MONTHLY DAILY AVERAGES           ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Days elapsed:       %d / %d\n", daysElapsed, daysInMonth)
fmt.Printf("Total profit:       $%.2f\n", totalProfit)
fmt.Printf("Avg per day:        $%.2f\n", avgPerDay)
fmt.Printf("Projected end:      $%.2f\n", projectedMonthEnd)
```

---

### 9) Risk/reward analysis this month

```go
deals, _ := sugar.GetDealsThisMonth()

if len(deals) == 0 {
    fmt.Println("No deals this month")
    return
}

var totalWins float64
var totalLosses float64
winCount := 0
lossCount := 0

for _, deal := range deals {
    if deal.Profit > 0 {
        totalWins += deal.Profit
        winCount++
    } else if deal.Profit < 0 {
        totalLosses += deal.Profit
        lossCount++
    }
}

avgWin := 0.0
if winCount > 0 {
    avgWin = totalWins / float64(winCount)
}

avgLoss := 0.0
if lossCount > 0 {
    avgLoss = totalLosses / float64(lossCount)
}

riskRewardRatio := 0.0
if avgLoss != 0 {
    riskRewardRatio = avgWin / -avgLoss
}

profitFactor := 0.0
if totalLosses != 0 {
    profitFactor = totalWins / -totalLosses
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      RISK/REWARD ANALYSIS             ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Avg Win:           $%.2f\n", avgWin)
fmt.Printf("Avg Loss:          $%.2f\n", avgLoss)
fmt.Printf("Risk/Reward:       %.2f:1\n", riskRewardRatio)
fmt.Printf("Profit Factor:     %.2f\n\n", profitFactor)

if profitFactor >= 2.0 {
    fmt.Println("✅ Excellent risk management")
} else if profitFactor >= 1.5 {
    fmt.Println("✅ Good risk management")
} else if profitFactor >= 1.0 {
    fmt.Println("🟡 Acceptable")
} else {
    fmt.Println("⚠️  Needs improvement")
}
```

---

### 10) Advanced monthly statistics analyzer

```go
type MonthlyStats struct {
    MonthName       string
    TotalDeals      int
    TotalProfit     float64
    WinRate         float64
    ProfitFactor    float64
    BestWeek        int
    BestWeekProfit  float64
    DaysTraded      int
    AvgDailyProfit  float64
    LargestWin      float64
    LargestLoss     float64
}

func AnalyzeMonthlyPerformance(sugar *mt5.MT5Sugar) (*MonthlyStats, error) {
    deals, err := sugar.GetDealsThisMonth()
    if err != nil {
        return nil, err
    }

    now := time.Now()
    stats := &MonthlyStats{
        MonthName:  now.Format("January 2006"),
        TotalDeals: len(deals),
    }

    if len(deals) == 0 {
        return stats, nil
    }

    winCount := 0
    var totalWins, totalLosses float64
    weeklyProfit := make(map[int]float64)
    tradingDays := make(map[string]bool)

    for _, deal := range deals {
        stats.TotalProfit += deal.Profit
        tradingDays[deal.TimeClose.Format("2006-01-02")] = true

        if deal.Profit > 0 {
            winCount++
            totalWins += deal.Profit

            if deal.Profit > stats.LargestWin {
                stats.LargestWin = deal.Profit
            }
        } else if deal.Profit < 0 {
            totalLosses += deal.Profit

            if deal.Profit < stats.LargestLoss {
                stats.LargestLoss = deal.Profit
            }
        }

        _, week := deal.TimeClose.ISOWeek()
        weeklyProfit[week] += deal.Profit
    }

    stats.WinRate = float64(winCount) / float64(len(deals)) * 100
    stats.DaysTraded = len(tradingDays)

    if stats.DaysTraded > 0 {
        stats.AvgDailyProfit = stats.TotalProfit / float64(stats.DaysTraded)
    }

    if totalLosses != 0 {
        stats.ProfitFactor = totalWins / -totalLosses
    }

    // Find best week
    for week, profit := range weeklyProfit {
        if profit > stats.BestWeekProfit {
            stats.BestWeekProfit = profit
            stats.BestWeek = week
        }
    }

    return stats, nil
}

func (s *MonthlyStats) Print() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      MONTHLY STATISTICS ANALYSIS      ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("Month:              %s\n\n", s.MonthName)

    if s.TotalDeals == 0 {
        fmt.Println("No trading activity this month")
        return
    }

    fmt.Printf("Total Deals:        %d\n", s.TotalDeals)
    fmt.Printf("Total Profit:       $%.2f\n", s.TotalProfit)
    fmt.Printf("Win Rate:           %.1f%%\n", s.WinRate)
    fmt.Printf("Profit Factor:      %.2f\n\n", s.ProfitFactor)

    fmt.Printf("Days Traded:        %d\n", s.DaysTraded)
    fmt.Printf("Avg Daily Profit:   $%.2f\n\n", s.AvgDailyProfit)

    fmt.Printf("Largest Win:        $%.2f\n", s.LargestWin)
    fmt.Printf("Largest Loss:       $%.2f\n\n", s.LargestLoss)

    fmt.Printf("Best Week:          Week %d ($%.2f)\n", s.BestWeek, s.BestWeekProfit)

    fmt.Println("\nMonth Assessment:")
    if s.TotalProfit > 0 && s.WinRate >= 60 && s.ProfitFactor >= 1.5 {
        fmt.Println("🌟 Excellent month")
    } else if s.TotalProfit > 0 && s.WinRate >= 50 {
        fmt.Println("✅ Good month")
    } else if s.TotalProfit >= 0 {
        fmt.Println("🟡 Breakeven month")
    } else {
        fmt.Println("❌ Losing month - review strategy")
    }
}

// Usage:
stats, _ := AnalyzeMonthlyPerformance(sugar)
stats.Print()
```

---

## 🔗 Related Methods

**🍬 Other time periods:**

* `GetDealsToday()` - Today's deals
* `GetDealsYesterday()` - Yesterday's deals
* `GetDealsThisWeek()` - This week's deals
* `GetDealsDateRange()` - Custom date range

**🍬 Profit calculations:**

* `GetProfitThisMonth()` - Just this month's profit total

---

## ⚠️ Common Pitfalls

### 1) Early month = incomplete data

```go
// ❌ WRONG - comparing incomplete month to complete month
thisMonthProfit, _ := sugar.GetProfitThisMonth()  // Day 5
lastMonthProfit := 2000.0                         // Full 30 days
// Not fair comparison!

// ✅ CORRECT - adjust for days
daysPassed := time.Now().Day()
dailyAvg := thisMonthProfit / float64(daysPassed)
fmt.Printf("Daily average: $%.2f\n", dailyAvg)
```

### 2) Not handling timezone

```go
// ❌ WRONG - "this month" is server timezone

// ✅ CORRECT - aware that it's server's month
deals, _ := sugar.GetDealsThisMonth()
fmt.Println("Server's current month")
```

### 3) Assuming deals exist

```go
// ❌ WRONG - might panic
deals, _ := sugar.GetDealsThisMonth()
avgProfit := totalProfit / len(deals) // Panic if empty!

// ✅ CORRECT - check length
if len(deals) > 0 {
    avgProfit := totalProfit / float64(len(deals))
}
```

---

## 💎 Pro Tips

1. **Month starts 1st** - Always starts from 1st day 00:00:00

2. **Incomplete data** - Mid-month calls return partial month

3. **Performance** - 5-second timeout, suitable for regular calls

4. **Empty OK** - Returns empty slice if no trades (not error)

5. **Auto-range** - Automatically calculates from 1st to now

---

## 📊 Month Range Calculation

```
Month range calculation:
now = time.Now()
startOfMonth = Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

Returns all deals from startOfMonth to now
```

---

**See also:** [`GetDealsToday.md`](GetDealsToday.md), [`GetDealsThisWeek.md`](GetDealsThisWeek.md), [`GetProfitThisMonth.md`](GetProfitThisMonth.md), [`GetDealsDateRange.md`](GetDealsDateRange.md)
