# 📅🎯 Get Deals Date Range (`GetDealsDateRange`)

> **Sugar method:** Returns all closed positions (deals) within a custom date range.

**API Information:**

* **Method:** `sugar.GetDealsDateRange(from, to)`
* **Timeout:** 5 seconds
* **Returns:** Slice of position history info

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetDealsDateRange(from, to time.Time) ([]*pb.PositionHistoryInfo, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `from` | `time.Time` | Start date/time (inclusive) |
| `to` | `time.Time` | End date/time (inclusive) |

| Output | Type | Description |
|--------|------|-------------|
| `[]*pb.PositionHistoryInfo` | slice | All closed positions in range |
| `error` | `error` | Error if retrieval failed |

---

## 💬 Just the Essentials

* **What it is:** Get all closed trades within a custom time range you specify.
* **Why you need it:** Custom period analysis, backtesting, specific timeframe reports.
* **Sanity check:** Returns empty slice if no deals in range. Both `from` and `to` are inclusive.

---

## 🎯 When to Use

✅ **Custom reports** - Analyze specific time periods

✅ **Backtesting** - Review performance for exact date ranges

✅ **Comparison** - Compare different periods

✅ **Historical analysis** - Study past performance patterns

---

## 🔗 Usage Examples

### 1) Basic usage - last 7 days

```go
to := time.Now()
from := to.AddDate(0, 0, -7) // 7 days ago

deals, err := sugar.GetDealsDateRange(from, to)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if len(deals) == 0 {
    fmt.Println("No deals in last 7 days")
    return
}

fmt.Printf("Last 7 days: %d deals\n\n", len(deals))

totalProfit := 0.0
for i, deal := range deals {
    fmt.Printf("%d. #%d %s: $%.2f\n",
        i+1, deal.Ticket, deal.Symbol, deal.Profit)
    totalProfit += deal.Profit
}

fmt.Printf("\nTotal: $%.2f\n", totalProfit)
```

---

### 2) Specific month analysis (last month)

```go
now := time.Now()

// First day of last month
firstOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())

// First day of this month (exclusive end)
firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

deals, _ := sugar.GetDealsDateRange(
    firstOfLastMonth,
    firstOfThisMonth.Add(-time.Second),
)

if len(deals) == 0 {
    fmt.Println("No deals last month")
    return
}

totalProfit := 0.0
for _, deal := range deals {
    totalProfit += deal.Profit
}

monthName := firstOfLastMonth.Format("January 2006")

fmt.Printf("%s Performance:\n", monthName)
fmt.Printf("  Trades:  %d\n", len(deals))
fmt.Printf("  Profit:  $%.2f\n", totalProfit)
```

---

### 3) Compare two equal periods

```go
// This week
now := time.Now()
weekday := int(now.Weekday())
if weekday == 0 {
    weekday = 7
}
thisWeekStart := now.AddDate(0, 0, -(weekday - 1))
thisWeekStart = time.Date(thisWeekStart.Year(), thisWeekStart.Month(),
    thisWeekStart.Day(), 0, 0, 0, 0, thisWeekStart.Location())

thisWeekDeals, _ := sugar.GetDealsDateRange(thisWeekStart, now)

// Last week
lastWeekStart := thisWeekStart.AddDate(0, 0, -7)
lastWeekEnd := thisWeekStart.Add(-time.Second)

lastWeekDeals, _ := sugar.GetDealsDateRange(lastWeekStart, lastWeekEnd)

thisWeekProfit := 0.0
for _, deal := range thisWeekDeals {
    thisWeekProfit += deal.Profit
}

lastWeekProfit := 0.0
for _, deal := range lastWeekDeals {
    lastWeekProfit += deal.Profit
}

fmt.Println("Week Comparison:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Last week:  %d trades, $%.2f\n",
    len(lastWeekDeals), lastWeekProfit)
fmt.Printf("This week:  %d trades, $%.2f\n",
    len(thisWeekDeals), thisWeekProfit)

diff := thisWeekProfit - lastWeekProfit
if diff > 0 {
    fmt.Printf("Improvement: +$%.2f 📈\n", diff)
} else if diff < 0 {
    fmt.Printf("Decline: $%.2f 📉\n", diff)
} else {
    fmt.Println("No change ➡️")
}
```

---

### 4) Specific date range (e.g., between two events)

```go
// Example: Performance between two specific dates
from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

deals, _ := sugar.GetDealsDateRange(from, to)

if len(deals) == 0 {
    fmt.Printf("No deals between %s and %s\n",
        from.Format("Jan 2"), to.Format("Jan 2"))
    return
}

totalProfit := 0.0
for _, deal := range deals {
    totalProfit += deal.Profit
}

fmt.Printf("Period: %s - %s\n",
    from.Format("Jan 2, 2006"),
    to.Format("Jan 2, 2006"))
fmt.Printf("Trades:  %d\n", len(deals))
fmt.Printf("Profit:  $%.2f\n", totalProfit)
```

---

### 5) Last N trading days

```go
func GetLastNTradingDays(sugar *mt5.MT5Sugar, n int) ([]*pb.PositionHistoryInfo, error) {
    to := time.Now()
    from := to.AddDate(0, 0, -n)

    return sugar.GetDealsDateRange(from, to)
}

// Usage: Last 30 trading days
deals, _ := GetLastNTradingDays(sugar, 30)

if len(deals) > 0 {
    totalProfit := 0.0
    for _, deal := range deals {
        totalProfit += deal.Profit
    }

    fmt.Printf("Last 30 days:\n")
    fmt.Printf("  Trades: %d\n", len(deals))
    fmt.Printf("  Profit: $%.2f\n", totalProfit)
    fmt.Printf("  Avg/day: $%.2f\n", totalProfit/30.0)
}
```

---

### 6) Quarter analysis (Q1, Q2, etc.)

```go
func GetQuarterDeals(sugar *mt5.MT5Sugar, year, quarter int) ([]*pb.PositionHistoryInfo, error) {
    var startMonth time.Month
    switch quarter {
    case 1:
        startMonth = time.January
    case 2:
        startMonth = time.April
    case 3:
        startMonth = time.July
    case 4:
        startMonth = time.October
    default:
        return nil, fmt.Errorf("invalid quarter: %d", quarter)
    }

    from := time.Date(year, startMonth, 1, 0, 0, 0, 0, time.UTC)
    to := from.AddDate(0, 3, 0).Add(-time.Second) // 3 months later minus 1 second

    return sugar.GetDealsDateRange(from, to)
}

// Usage: Q1 2024
deals, _ := GetQuarterDeals(sugar, 2024, 1)

if len(deals) > 0 {
    totalProfit := 0.0
    for _, deal := range deals {
        totalProfit += deal.Profit
    }

    fmt.Println("Q1 2024 Performance:")
    fmt.Printf("  Trades: %d\n", len(deals))
    fmt.Printf("  Profit: $%.2f\n", totalProfit)
}
```

---

### 7) Year-to-date performance

```go
now := time.Now()
startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

deals, _ := sugar.GetDealsDateRange(startOfYear, now)

if len(deals) == 0 {
    fmt.Println("No deals year-to-date")
    return
}

totalProfit := 0.0
winCount := 0

for _, deal := range deals {
    totalProfit += deal.Profit
    if deal.Profit > 0 {
        winCount++
    }
}

winRate := float64(winCount) / float64(len(deals)) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      YEAR-TO-DATE PERFORMANCE         ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Year:        %d\n", now.Year())
fmt.Printf("Trades:      %d\n", len(deals))
fmt.Printf("Win rate:    %.1f%%\n", winRate)
fmt.Printf("Total P/L:   $%.2f\n", totalProfit)

daysInYear := now.YearDay()
avgPerDay := totalProfit / float64(daysInYear)
fmt.Printf("Avg/day:     $%.2f\n", avgPerDay)
```

---

### 8) Hourly range analysis (specific trading hours)

```go
// Analyze performance during specific hours (e.g., 9 AM - 5 PM)
from := time.Now().AddDate(0, 0, -30) // Last 30 days
to := time.Now()

deals, _ := sugar.GetDealsDateRange(from, to)

if len(deals) == 0 {
    fmt.Println("No deals in range")
    return
}

// Filter by closing hour
tradingHoursDeals := []*pb.PositionHistoryInfo{}
for _, deal := range deals {
    hour := deal.TimeClose.Hour()
    if hour >= 9 && hour < 17 { // 9 AM to 5 PM
        tradingHoursDeals = append(tradingHoursDeals, deal)
    }
}

tradingHoursProfit := 0.0
for _, deal := range tradingHoursDeals {
    tradingHoursProfit += deal.Profit
}

fmt.Println("Trading Hours Analysis (9 AM - 5 PM):")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Total deals:         %d\n", len(deals))
fmt.Printf("In trading hours:    %d (%.1f%%)\n",
    len(tradingHoursDeals),
    float64(len(tradingHoursDeals))/float64(len(deals))*100)
fmt.Printf("Trading hours P/L:   $%.2f\n", tradingHoursProfit)
```

---

### 9) Compare same period from different years

```go
// This year January
thisYearJan := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
thisYearFeb := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

thisYearDeals, _ := sugar.GetDealsDateRange(thisYearJan, thisYearFeb.Add(-time.Second))

// Last year January
lastYearJan := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
lastYearFeb := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)

lastYearDeals, _ := sugar.GetDealsDateRange(lastYearJan, lastYearFeb.Add(-time.Second))

thisYearProfit := 0.0
for _, deal := range thisYearDeals {
    thisYearProfit += deal.Profit
}

lastYearProfit := 0.0
for _, deal := range lastYearDeals {
    lastYearProfit += deal.Profit
}

fmt.Println("January Year-over-Year Comparison:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Jan 2023:  %d trades, $%.2f\n",
    len(lastYearDeals), lastYearProfit)
fmt.Printf("Jan 2024:  %d trades, $%.2f\n",
    len(thisYearDeals), thisYearProfit)

if thisYearProfit > lastYearProfit {
    fmt.Println("📈 Better performance this year")
} else if thisYearProfit < lastYearProfit {
    fmt.Println("📉 Worse performance this year")
} else {
    fmt.Println("➡️  Same performance")
}
```

---

### 10) Advanced date range analyzer

```go
type DateRangeStats struct {
    From            time.Time
    To              time.Time
    TotalDeals      int
    TotalProfit     float64
    WinRate         float64
    BestDay         string
    BestDayProfit   float64
    WorstDay        string
    WorstDayProfit  float64
    DaysWithTrades  int
}

func AnalyzeDateRange(sugar *mt5.MT5Sugar, from, to time.Time) (*DateRangeStats, error) {
    deals, err := sugar.GetDealsDateRange(from, to)
    if err != nil {
        return nil, err
    }

    stats := &DateRangeStats{
        From:       from,
        To:         to,
        TotalDeals: len(deals),
    }

    if len(deals) == 0 {
        return stats, nil
    }

    winCount := 0
    dailyProfit := make(map[string]float64)

    for _, deal := range deals {
        stats.TotalProfit += deal.Profit

        if deal.Profit > 0 {
            winCount++
        }

        dayKey := deal.TimeClose.Format("2006-01-02")
        dailyProfit[dayKey] += deal.Profit
    }

    stats.WinRate = float64(winCount) / float64(len(deals)) * 100
    stats.DaysWithTrades = len(dailyProfit)

    // Find best/worst days
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
    }

    return stats, nil
}

func (s *DateRangeStats) Print() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      DATE RANGE ANALYSIS              ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    fmt.Printf("From:              %s\n", s.From.Format("2006-01-02"))
    fmt.Printf("To:                %s\n\n", s.To.Format("2006-01-02"))

    if s.TotalDeals == 0 {
        fmt.Println("No trading activity in this period")
        return
    }

    duration := s.To.Sub(s.From)
    days := int(duration.Hours() / 24)

    fmt.Printf("Total Deals:       %d\n", s.TotalDeals)
    fmt.Printf("Total Profit:      $%.2f\n", s.TotalProfit)
    fmt.Printf("Win Rate:          %.1f%%\n\n", s.WinRate)

    fmt.Printf("Period Length:     %d days\n", days)
    fmt.Printf("Days with Trades:  %d\n", s.DaysWithTrades)

    if days > 0 {
        fmt.Printf("Avg Profit/Day:    $%.2f\n\n", s.TotalProfit/float64(days))
    }

    fmt.Printf("Best Day:          %s ($%.2f)\n", s.BestDay, s.BestDayProfit)
    fmt.Printf("Worst Day:         %s ($%.2f)\n", s.WorstDay, s.WorstDayProfit)
}

// Usage:
from := time.Now().AddDate(0, -1, 0) // 1 month ago
to := time.Now()

stats, _ := AnalyzeDateRange(sugar, from, to)
stats.Print()
```

---

## 🔗 Related Methods

**🍬 Pre-defined periods:**

* `GetDealsToday()` - Today's deals (00:00 to now)
* `GetDealsYesterday()` - Yesterday's deals (full day)
* `GetDealsThisWeek()` - This week's deals (Monday to now)
* `GetDealsThisMonth()` - This month's deals (1st to now)

**🍬 Profit calculations:**

* `GetProfitToday()` - Profit total methods

---

## ⚠️ Common Pitfalls

### 1) Inverted date range

```go
// ❌ WRONG - from > to
to := time.Now()
from := to.AddDate(0, 0, 7) // 7 days in future!
deals, _ := sugar.GetDealsDateRange(from, to) // Returns nothing

// ✅ CORRECT - from < to
from := to.AddDate(0, 0, -7) // 7 days in past
deals, _ := sugar.GetDealsDateRange(from, to)
```

### 2) Not handling timezone

```go
// ❌ WRONG - mixing timezones
from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.Local) // Different TZ!

// ✅ CORRECT - same timezone
from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
```

### 3) Exclusive vs inclusive range

```go
// ❌ WRONG - assuming exclusive end
to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC) // Excludes Jan 31!

// ✅ CORRECT - include full last day
to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
```

---

## 💎 Pro Tips

1. **Both inclusive** - Both `from` and `to` are inclusive

2. **Timezone aware** - Use same timezone for both parameters

3. **Flexible** - Most flexible history method for custom periods

4. **Performance** - 5-second timeout, suitable for analysis

5. **Empty OK** - Returns empty slice if no trades (not error)

---

## 📊 Common Date Range Patterns

```go
// Last 7 days
from := time.Now().AddDate(0, 0, -7)
to := time.Now()

// Last 30 days
from := time.Now().AddDate(0, 0, -30)
to := time.Now()

// Last month (complete)
now := time.Now()
firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)
from := firstOfLastMonth
to := firstOfThisMonth.Add(-time.Second)

// This year
from := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
to := time.Now()

// Specific month
from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
```

---

**See also:** [`GetDealsToday.md`](GetDealsToday.md), [`GetDealsThisWeek.md`](GetDealsThisWeek.md), [`GetDealsThisMonth.md`](GetDealsThisMonth.md)
