# 📅📊 Get Deals This Week (`GetDealsThisWeek`)

> **Sugar method:** Returns closed positions ("deals") that were **opened** this week (Monday 00:00 to now).

**API Information:**

* **Method:** `sugar.GetDealsThisWeek()`
* **Timeout:** 5 seconds
* **Returns:** Slice of position history info, sorted by open time (oldest first)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetDealsThisWeek() ([]*pb.PositionHistoryInfo, error)
```

How it works: Monday 00:00 of the current week (your machine's clock) and `time.Now()` are passed to
`MT5Service.GetPositionsHistory(ctx, pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC, &startOfWeek, &now, nil, nil)`.
The range is sent as `PositionOpenTimeFrom` / `PositionOpenTimeTo`, so it filters by the position's **open** time:
a position opened last week and closed this week is not included. No paging parameters are sent.

**Time zone note:** the boundaries come from your local clock and are sent as absolute instants in the request fields that
the server reads as broker server time. If your machine, UTC and the broker server use different time zones, the window is
shifted accordingly - use `GetDealsDateRange()` when you need exact boundaries.

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates week range) |

| Output | Type | Description |
|--------|------|-------------|
| `[]*pb.PositionHistoryInfo` | slice | Closed positions opened this week |
| `error` | `error` | `"GetPositionsHistory failed: ..."` if retrieval failed |

---

## 💬 Just the Essentials

* **What it is:** Get closed trades opened between Monday 00:00 and the current time.
* **Why you need it:** Weekly performance tracking, strategy analysis, week-over-week comparison.
* **Sanity check:** Returns an empty (possibly `nil`) slice if there are no deals - check `len(deals)`. Week starts on Monday (on Sunday the range covers Monday..Sunday).

---

## 🎯 When to Use

✅ **Weekly reports** - Generate week-to-date summaries

✅ **Performance tracking** - Monitor weekly profit/loss

✅ **Strategy analysis** - Review week's trading decisions

✅ **Goal tracking** - Check progress toward weekly targets

---

## 🔗 Usage Examples

### 1) Basic usage - show this week's deals

```go
deals, err := sugar.GetDealsThisWeek()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if len(deals) == 0 {
    fmt.Println("No deals this week")
    return
}

fmt.Printf("This week's deals: %d\n\n", len(deals))

totalProfit := 0.0
for i, deal := range deals {
    fmt.Printf("%d. #%d %s: $%.2f\n",
        i+1, deal.PositionTicket, deal.Symbol, deal.Profit)
    totalProfit += deal.Profit
}

fmt.Printf("\nWeek total: $%.2f\n", totalProfit)
```

---

### 2) Weekly performance summary

```go
deals, _ := sugar.GetDealsThisWeek()

if len(deals) == 0 {
    fmt.Println("No trading this week")
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
fmt.Println("║      THIS WEEK'S PERFORMANCE          ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Total trades:  %d\n", len(deals))
fmt.Printf("Winners:       %d (%.1f%%)\n", winCount, winRate)
fmt.Printf("Losers:        %d (%.1f%%)\n", lossCount, 100-winRate)
fmt.Printf("Total P/L:     $%.2f\n", totalProfit)

if totalProfit > 0 {
    fmt.Println("\n✅ Profitable week so far")
} else {
    fmt.Println("\n❌ Losing week so far")
}
```

---

### 3) Daily breakdown of this week

```go
deals, _ := sugar.GetDealsThisWeek()

if len(deals) == 0 {
    fmt.Println("No deals this week")
    return
}

// Group by close day (CloseTime is broker server time)
dailyDeals := make(map[string][]float64)

for _, deal := range deals {
    dayKey := deal.CloseTime.AsTime().Format("Mon 01/02")
    dailyDeals[dayKey] = append(dailyDeals[dayKey], deal.Profit)
}

fmt.Println("This Week - Daily Breakdown:")
fmt.Println("─────────────────────────────────────────")

// Get days in order (Monday to Sunday)
now := time.Now()
weekday := int(now.Weekday())
if weekday == 0 {
    weekday = 7 // Sunday becomes 7
}

for i := 0; i < weekday; i++ {
    day := now.AddDate(0, 0, -(weekday - 1 - i))
    dayKey := day.Format("Mon 01/02")

    profits, exists := dailyDeals[dayKey]
    if !exists {
        fmt.Printf("%s: No trades\n", dayKey)
        continue
    }

    dayProfit := 0.0
    for _, p := range profits {
        dayProfit += p
    }

    status := "✅"
    if dayProfit < 0 {
        status = "❌"
    }

    fmt.Printf("%s %s: %d trades, $%.2f\n",
        status, dayKey, len(profits), dayProfit)
}
```

---

### 4) Check weekly target progress

```go
weeklyTarget := 500.0 // $500 per week

deals, _ := sugar.GetDealsThisWeek()

if len(deals) == 0 {
    fmt.Printf("Weekly target: $%.2f\n", weeklyTarget)
    fmt.Println("No trades yet this week")
    return
}

actualProfit := 0.0
for _, deal := range deals {
    actualProfit += deal.Profit
}

progress := (actualProfit / weeklyTarget) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      WEEKLY TARGET PROGRESS           ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Target:    $%.2f\n", weeklyTarget)
fmt.Printf("Current:   $%.2f\n", actualProfit)
fmt.Printf("Progress:  %.1f%%\n\n", progress)

remaining := weeklyTarget - actualProfit

if actualProfit >= weeklyTarget {
    fmt.Printf("🎯 Target achieved! (+$%.2f)\n",
        actualProfit-weeklyTarget)
} else if progress >= 80 {
    fmt.Printf("🟡 Almost there! ($%.2f remaining)\n", remaining)
} else {
    fmt.Printf("📊 In progress ($%.2f remaining)\n", remaining)
}
```

---

### 5) Best trading day this week

```go
deals, _ := sugar.GetDealsThisWeek()

if len(deals) == 0 {
    fmt.Println("No deals this week")
    return
}

// Group by day and calculate daily profit
dailyProfit := make(map[string]float64)

for _, deal := range deals {
    dayKey := deal.CloseTime.AsTime().Format("Monday 01/02")
    dailyProfit[dayKey] += deal.Profit
}

var bestDay string
var bestProfit float64 = -999999

for day, profit := range dailyProfit {
    if profit > bestProfit {
        bestProfit = profit
        bestDay = day
    }
}

fmt.Println("This Week's Best Day:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Day:     %s\n", bestDay)
fmt.Printf("Profit:  $%.2f\n", bestProfit)
```

---

### 6) Symbol performance this week

```go
deals, _ := sugar.GetDealsThisWeek()

if len(deals) == 0 {
    fmt.Println("No deals this week")
    return
}

symbolStats := make(map[string]struct {
    Count  int
    Profit float64
})

for _, deal := range deals {
    stats := symbolStats[deal.Symbol]
    stats.Count++
    stats.Profit += deal.Profit
    symbolStats[deal.Symbol] = stats
}

fmt.Println("This Week - Symbol Performance:")
fmt.Println("─────────────────────────────────────────")

for symbol, stats := range symbolStats {
    avgProfit := stats.Profit / float64(stats.Count)

    status := "✅"
    if stats.Profit < 0 {
        status = "❌"
    }

    fmt.Printf("%s %-8s: %d trades, $%.2f (avg: $%.2f)\n",
        status, symbol, stats.Count, stats.Profit, avgProfit)
}
```

---

### 7) Week-over-week comparison

```go
thisWeekDeals, _ := sugar.GetDealsThisWeek()

// Get last week's range
now := time.Now()
weekday := int(now.Weekday())
if weekday == 0 {
    weekday = 7
}

// Last week: previous Monday 00:00 to Sunday 23:59:59 (local clock)
thisMonday := now.AddDate(0, 0, -(weekday - 1))
thisMonday = time.Date(thisMonday.Year(), thisMonday.Month(), thisMonday.Day(), 0, 0, 0, 0, now.Location())
lastMonday := thisMonday.AddDate(0, 0, -7)
lastSunday := thisMonday.Add(-time.Second)

lastWeekDeals, _ := sugar.GetDealsDateRange(lastMonday, lastSunday)

thisWeekProfit := 0.0
for _, deal := range thisWeekDeals {
    thisWeekProfit += deal.Profit
}

lastWeekProfit := 0.0
for _, deal := range lastWeekDeals {
    lastWeekProfit += deal.Profit
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      WEEK-OVER-WEEK COMPARISON        ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Last week:  %d trades, $%.2f\n",
    len(lastWeekDeals), lastWeekProfit)
fmt.Printf("This week:  %d trades, $%.2f\n",
    len(thisWeekDeals), thisWeekProfit)

diff := thisWeekProfit - lastWeekProfit

fmt.Printf("\nDifference: $%.2f ", diff)
if diff > 0 {
    pct := (diff / lastWeekProfit) * 100
    fmt.Printf("📈 (+%.1f%%)\n", pct)
} else if diff < 0 {
    pct := (-diff / lastWeekProfit) * 100
    fmt.Printf("📉 (-%.1f%%)\n", pct)
} else {
    fmt.Println("➡️  (No change)")
}
```

---

### 8) Weekly consistency check

```go
deals, _ := sugar.GetDealsThisWeek()

if len(deals) == 0 {
    fmt.Println("No deals this week")
    return
}

// Group by day
dailyProfit := make(map[string]float64)

for _, deal := range deals {
    dayKey := deal.CloseTime.AsTime().Format("Mon")
    dailyProfit[dayKey] += deal.Profit
}

positiveDays := 0
for _, profit := range dailyProfit {
    if profit > 0 {
        positiveDays++
    }
}

consistency := float64(positiveDays) / float64(len(dailyProfit)) * 100

fmt.Println("Weekly Consistency:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Trading days:     %d\n", len(dailyProfit))
fmt.Printf("Profitable days:  %d (%.1f%%)\n",
    positiveDays, consistency)

if consistency >= 80 {
    fmt.Println("✅ Excellent consistency")
} else if consistency >= 60 {
    fmt.Println("🟡 Good consistency")
} else {
    fmt.Println("⚠️  Needs improvement")
}
```

---

### 9) Week projection (estimate end-of-week)

```go
deals, _ := sugar.GetDealsThisWeek()

if len(deals) == 0 {
    fmt.Println("No deals this week yet")
    return
}

currentProfit := 0.0
for _, deal := range deals {
    currentProfit += deal.Profit
}

now := time.Now()
weekday := int(now.Weekday())
if weekday == 0 {
    weekday = 7
}

// Days elapsed in week (including today)
daysElapsed := weekday

// Project to end of week (5 trading days)
projectedProfit := (currentProfit / float64(daysElapsed)) * 5.0

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      WEEKLY PROJECTION                ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Days elapsed:      %d\n", daysElapsed)
fmt.Printf("Current profit:    $%.2f\n", currentProfit)
fmt.Printf("Daily average:     $%.2f\n",
    currentProfit/float64(daysElapsed))
fmt.Printf("Week projection:   $%.2f\n", projectedProfit)

fmt.Println("\n⚠️  Note: Projection assumes consistent daily performance")
```

---

### 10) Advanced weekly statistics analyzer

```go
type WeeklyStats struct {
    TotalDeals      int
    TotalProfit     float64
    WinRate         float64
    BestDay         string
    BestDayProfit   float64
    WorstDay        string
    WorstDayProfit  float64
    MostTradedDay   string
    DaysTraded      int
    AvgDailyProfit  float64
}

func AnalyzeWeeklyPerformance(sugar *mt5.MT5Sugar) (*WeeklyStats, error) {
    deals, err := sugar.GetDealsThisWeek()
    if err != nil {
        return nil, err
    }

    stats := &WeeklyStats{
        TotalDeals: len(deals),
    }

    if len(deals) == 0 {
        return stats, nil
    }

    winCount := 0
    dailyProfit := make(map[string]float64)
    dailyCount := make(map[string]int)

    for _, deal := range deals {
        stats.TotalProfit += deal.Profit

        if deal.Profit > 0 {
            winCount++
        }

        dayKey := deal.CloseTime.AsTime().Format("Mon 01/02")
        dailyProfit[dayKey] += deal.Profit
        dailyCount[dayKey]++
    }

    stats.WinRate = float64(winCount) / float64(len(deals)) * 100
    stats.DaysTraded = len(dailyProfit)

    if stats.DaysTraded > 0 {
        stats.AvgDailyProfit = stats.TotalProfit / float64(stats.DaysTraded)
    }

    // Find best/worst days
    maxCount := 0
    stats.BestDayProfit = -999999
    stats.WorstDayProfit = 999999

    for day, profit := range dailyProfit {
        if profit > stats.BestDayProfit {
            stats.BestDayProfit = profit
            stats.BestDay = day
        }

        if profit < stats.WorstDayProfit {
            stats.WorstDayProfit = profit
            stats.WorstDay = day
        }

        if dailyCount[day] > maxCount {
            maxCount = dailyCount[day]
            stats.MostTradedDay = day
        }
    }

    return stats, nil
}

func (s *WeeklyStats) Print() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      WEEKLY STATISTICS ANALYSIS       ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    if s.TotalDeals == 0 {
        fmt.Println("No trading activity this week")
        return
    }

    fmt.Printf("Total Deals:        %d\n", s.TotalDeals)
    fmt.Printf("Total Profit:       $%.2f\n", s.TotalProfit)
    fmt.Printf("Win Rate:           %.1f%%\n", s.WinRate)
    fmt.Printf("Days Traded:        %d\n", s.DaysTraded)
    fmt.Printf("Avg Daily Profit:   $%.2f\n\n", s.AvgDailyProfit)

    fmt.Printf("Best Day:           %s ($%.2f)\n",
        s.BestDay, s.BestDayProfit)
    fmt.Printf("Worst Day:          %s ($%.2f)\n",
        s.WorstDay, s.WorstDayProfit)
    fmt.Printf("Most Traded Day:    %s\n", s.MostTradedDay)

    fmt.Println("\nWeek Assessment:")
    if s.TotalProfit > 0 && s.WinRate >= 60 {
        fmt.Println("🌟 Excellent week")
    } else if s.TotalProfit > 0 && s.WinRate >= 50 {
        fmt.Println("✅ Good week")
    } else if s.TotalProfit >= 0 {
        fmt.Println("🟡 Breakeven week")
    } else {
        fmt.Println("❌ Losing week - review strategy")
    }
}

// Usage:
stats, _ := AnalyzeWeeklyPerformance(sugar)
stats.Print()
```

---

## 🔗 Related Methods

**🍬 Other time periods:**

* `GetDealsToday()` - Today's deals
* `GetDealsYesterday()` - Yesterday's deals
* `GetDealsThisMonth()` - This month's deals
* `GetDealsDateRange()` - Custom date range

**🍬 Profit calculations:**

* `GetProfitThisWeek()` - Just this week's profit total

---

## ⚠️ Common Pitfalls

### 1) Assuming week starts on Sunday

```go
// ❌ WRONG - GetDealsThisWeek starts the week on Monday
// Your local culture might use Sunday

// ✅ CORRECT - aware that week starts Monday
deals, _ := sugar.GetDealsThisWeek()
// Returns closed positions opened from Monday 00:00 (local clock) to now
```

### 1b) Expecting positions opened last week

```go
// ❌ WRONG - a position opened last Friday and closed on Monday is NOT returned
// (the range filters on the position's OPEN time)
deals, _ := sugar.GetDealsThisWeek()

// ✅ CORRECT - widen the range and filter on CloseTime yourself
deals, _ := sugar.GetDealsDateRange(time.Now().AddDate(0, 0, -14), time.Now())
```

### 2) Comparing incomplete week to complete week

```go
// ❌ WRONG - comparing mid-week to full week
thisWeekProfit, _ := sugar.GetProfitThisWeek()  // Mon-Wed
lastWeekProfit := 500.0                         // Mon-Fri
// Not fair comparison!

// ✅ CORRECT - note the difference
fmt.Println("This week (incomplete) vs last week (complete)")
```

### 3) Empty slice on Sunday/Monday

```go
// ❌ WRONG - not handling early week
deals, _ := sugar.GetDealsThisWeek()
avgProfit := totalProfit / float64(len(deals)) // NaN if empty!

// ✅ CORRECT - check for empty
if len(deals) > 0 {
    avgProfit := totalProfit / float64(len(deals))
}
```

---

## 💎 Pro Tips

1. **Week starts Monday** - GetDealsThisWeek uses Monday as first day of week

2. **Incomplete data** - Mid-week calls return partial week data

3. **Performance** - 5-second timeout

4. **Empty OK** - Returns empty slice if no trades (not error)

5. **Auto-range** - Automatically calculates Monday 00:00 to now

---

## 📊 Week Range Calculation

```
Week range calculation:
now = time.Now()
weekday = int(now.Weekday())
if weekday == 0 { weekday = 7 }  // Sunday becomes 7

startOfWeek = now - (weekday - 1) days at 00:00:00 (now.Location())

Returns closed positions whose OPEN time is between startOfWeek and now
(PositionOpenTimeFrom / PositionOpenTimeTo), sorted by open time ascending
```

---

**See also:** [`GetDealsToday.md`](GetDealsToday.md), [`GetDealsThisMonth.md`](GetDealsThisMonth.md), [`GetProfitThisWeek.md`](GetProfitThisWeek.md), [`GetDealsDateRange.md`](GetDealsDateRange.md)
