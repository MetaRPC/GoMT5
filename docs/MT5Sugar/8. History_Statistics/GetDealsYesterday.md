# 📅📋 Get Deals Yesterday (`GetDealsYesterday`)

> **Sugar method:** Returns all closed positions (deals) from yesterday (full day 00:00 to 23:59:59).

**API Information:**

* **Method:** `sugar.GetDealsYesterday()`
* **Timeout:** 5 seconds
* **Returns:** Slice of position history info

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetDealsYesterday() ([]*pb.PositionHistoryInfo, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates yesterday's range) |

| Output | Type | Description |
|--------|------|-------------|
| `[]*pb.PositionHistoryInfo` | slice | All closed positions from yesterday |
| `error` | `error` | Error if retrieval failed |

---

## 💬 Just the Essentials

* **What it is:** Get all closed trades from yesterday's full trading day.
* **Why you need it:** Historical analysis, compare with today, daily reports.
* **Sanity check:** Returns empty slice if no deals yesterday. Auto-calculates full day range.

---

## 🎯 When to Use

✅ **Daily comparisons** - Compare yesterday vs today

✅ **Historical analysis** - Review previous day's performance

✅ **Morning reports** - Check yesterday's results at day start

✅ **Trend analysis** - Track daily performance over time

---

## 🔗 Usage Examples

### 1) Basic usage - show yesterday's deals

```go
deals, err := sugar.GetDealsYesterday()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if len(deals) == 0 {
    fmt.Println("No deals yesterday")
    return
}

fmt.Printf("Yesterday's deals: %d\n\n", len(deals))

totalProfit := 0.0
for i, deal := range deals {
    fmt.Printf("%d. #%d %s: $%.2f\n",
        i+1, deal.Ticket, deal.Symbol, deal.Profit)
    totalProfit += deal.Profit
}

fmt.Printf("\nTotal: $%.2f\n", totalProfit)
```

---

### 2) Compare yesterday vs today

```go
yesterdayDeals, _ := sugar.GetDealsYesterday()
todayDeals, _ := sugar.GetDealsToday()

yesterdayProfit := 0.0
for _, deal := range yesterdayDeals {
    yesterdayProfit += deal.Profit
}

todayProfit := 0.0
for _, deal := range todayDeals {
    todayProfit += deal.Profit
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      YESTERDAY VS TODAY               ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Yesterday: %d trades, $%.2f\n",
    len(yesterdayDeals), yesterdayProfit)
fmt.Printf("Today:     %d trades, $%.2f\n",
    len(todayDeals), todayProfit)

diff := todayProfit - yesterdayProfit

fmt.Printf("\nDifference: $%.2f ", diff)
if diff > 0 {
    fmt.Printf("📈 (Better today)\n")
} else if diff < 0 {
    fmt.Printf("📉 (Worse today)\n")
} else {
    fmt.Println("➡️  (Same)")
}
```

---

### 3) Yesterday's win rate

```go
deals, _ := sugar.GetDealsYesterday()

if len(deals) == 0 {
    fmt.Println("No deals yesterday")
    return
}

winCount := 0
lossCount := 0
totalProfit := 0.0

for _, deal := range deals {
    totalProfit += deal.Profit

    if deal.Profit > 0 {
        winCount++
    } else if deal.Profit < 0 {
        lossCount++
    }
}

winRate := float64(winCount) / float64(len(deals)) * 100

fmt.Println("Yesterday's Performance:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Total trades: %d\n", len(deals))
fmt.Printf("Winners:      %d (%.1f%%)\n", winCount, winRate)
fmt.Printf("Losers:       %d (%.1f%%)\n", lossCount, 100-winRate)
fmt.Printf("Total P/L:    $%.2f\n", totalProfit)
```

---

### 4) Morning report - yesterday's summary

```go
func ShowMorningReport(sugar *mt5.MT5Sugar) {
    deals, _ := sugar.GetDealsYesterday()

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      MORNING REPORT - YESTERDAY       ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    if len(deals) == 0 {
        fmt.Println("No trading activity yesterday")
        return
    }

    totalProfit := 0.0
    symbolProfit := make(map[string]float64)

    for _, deal := range deals {
        totalProfit += deal.Profit
        symbolProfit[deal.Symbol] += deal.Profit
    }

    fmt.Printf("Total Trades:  %d\n", len(deals))
    fmt.Printf("Total P/L:     $%.2f\n\n", totalProfit)

    fmt.Println("By Symbol:")
    for symbol, profit := range symbolProfit {
        fmt.Printf("  %-8s: $%.2f\n", symbol, profit)
    }

    if totalProfit > 0 {
        fmt.Println("\n✅ Yesterday was profitable")
    } else {
        fmt.Println("\n❌ Yesterday was a loss")
    }
}

// Usage: Call this in your morning routine
ShowMorningReport(sugar)
```

---

### 5) Find best/worst trades yesterday

```go
deals, _ := sugar.GetDealsYesterday()

if len(deals) == 0 {
    fmt.Println("No deals yesterday")
    return
}

var bestTrade *pb.PositionHistoryInfo
var worstTrade *pb.PositionHistoryInfo

for _, deal := range deals {
    if bestTrade == nil || deal.Profit > bestTrade.Profit {
        bestTrade = deal
    }
    if worstTrade == nil || deal.Profit < worstTrade.Profit {
        worstTrade = deal
    }
}

fmt.Println("Yesterday's Best & Worst:")
fmt.Println("─────────────────────────────────────────")

fmt.Printf("🏆 Best:  #%d %s: $%.2f\n",
    bestTrade.Ticket, bestTrade.Symbol, bestTrade.Profit)

fmt.Printf("💔 Worst: #%d %s: $%.2f\n",
    worstTrade.Ticket, worstTrade.Symbol, worstTrade.Profit)

diff := bestTrade.Profit - worstTrade.Profit
fmt.Printf("\nRange: $%.2f\n", diff)
```

---

### 6) Yesterday's trading hours analysis

```go
deals, _ := sugar.GetDealsYesterday()

if len(deals) == 0 {
    fmt.Println("No deals yesterday")
    return
}

// Analyze when most trades closed
hourlyCount := make(map[int]int)
hourlyProfit := make(map[int]float64)

for _, deal := range deals {
    hour := deal.TimeClose.Hour()
    hourlyCount[hour]++
    hourlyProfit[hour] += deal.Profit
}

fmt.Println("Yesterday's Trading Hours:")
fmt.Println("─────────────────────────────────────────")

for hour := 0; hour < 24; hour++ {
    count := hourlyCount[hour]
    if count > 0 {
        profit := hourlyProfit[hour]
        fmt.Printf("%02d:00-%02d:59: %d trades, $%.2f\n",
            hour, hour, count, profit)
    }
}
```

---

### 7) Check if yesterday met target

```go
targetProfit := 100.0 // $100 daily target

deals, _ := sugar.GetDealsYesterday()

if len(deals) == 0 {
    fmt.Println("No trading yesterday")
    return
}

actualProfit := 0.0
for _, deal := range deals {
    actualProfit += deal.Profit
}

fmt.Printf("Yesterday's Performance:\n")
fmt.Printf("Target:  $%.2f\n", targetProfit)
fmt.Printf("Actual:  $%.2f\n", actualProfit)

if actualProfit >= targetProfit {
    pct := (actualProfit / targetProfit) * 100
    fmt.Printf("✅ Target achieved (%.0f%%)\n", pct)
} else {
    shortfall := targetProfit - actualProfit
    fmt.Printf("❌ Missed target by $%.2f\n", shortfall)
}
```

---

### 8) Yesterday vs last 7 days average

```go
yesterdayDeals, _ := sugar.GetDealsYesterday()

yesterdayProfit := 0.0
for _, deal := range yesterdayDeals {
    yesterdayProfit += deal.Profit
}

// Get last 7 days
to := time.Now().AddDate(0, 0, -1) // End at yesterday
from := to.AddDate(0, 0, -7)       // 7 days before

weekDeals, _ := sugar.GetDealsDateRange(from, to)

weekProfit := 0.0
for _, deal := range weekDeals {
    weekProfit += deal.Profit
}

avgProfit := weekProfit / 7.0

fmt.Println("Yesterday vs Weekly Average:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Yesterday:       $%.2f\n", yesterdayProfit)
fmt.Printf("7-day average:   $%.2f\n", avgProfit)

diff := yesterdayProfit - avgProfit
if diff > 0 {
    fmt.Printf("Above average by $%.2f 📈\n", diff)
} else if diff < 0 {
    fmt.Printf("Below average by $%.2f 📉\n", -diff)
} else {
    fmt.Println("Exactly average ➡️")
}
```

---

### 9) Symbol diversity check (yesterday)

```go
deals, _ := sugar.GetDealsYesterday()

if len(deals) == 0 {
    fmt.Println("No deals yesterday")
    return
}

symbolCount := make(map[string]int)

for _, deal := range deals {
    symbolCount[deal.Symbol]++
}

fmt.Println("Yesterday's Symbol Diversity:")
fmt.Println("─────────────────────────────────────────")
fmt.Printf("Total symbols traded: %d\n\n", len(symbolCount))

for symbol, count := range symbolCount {
    pct := float64(count) / float64(len(deals)) * 100
    fmt.Printf("%-8s: %d trades (%.1f%%)\n",
        symbol, count, pct)
}

// Diversity assessment
if len(symbolCount) == 1 {
    fmt.Println("\n⚠️  Warning: Trading only 1 symbol")
} else if len(symbolCount) >= 3 {
    fmt.Println("\n✅ Good diversification")
}
```

---

### 10) Advanced yesterday performance analyzer

```go
type YesterdayPerformance struct {
    TotalDeals     int
    TotalProfit    float64
    WinRate        float64
    BestTrade      float64
    WorstTrade     float64
    AvgWin         float64
    AvgLoss        float64
    ProfitFactor   float64
    SymbolsTraded  int
}

func AnalyzeYesterday(sugar *mt5.MT5Sugar) (*YesterdayPerformance, error) {
    deals, err := sugar.GetDealsYesterday()
    if err != nil {
        return nil, err
    }

    perf := &YesterdayPerformance{
        TotalDeals: len(deals),
    }

    if len(deals) == 0 {
        return perf, nil
    }

    winCount := 0
    winSum := 0.0
    lossSum := 0.0
    symbols := make(map[string]bool)

    for _, deal := range deals {
        perf.TotalProfit += deal.Profit
        symbols[deal.Symbol] = true

        if deal.Profit > 0 {
            winCount++
            winSum += deal.Profit
            if deal.Profit > perf.BestTrade {
                perf.BestTrade = deal.Profit
            }
        } else if deal.Profit < 0 {
            lossSum += deal.Profit
            if deal.Profit < perf.WorstTrade {
                perf.WorstTrade = deal.Profit
            }
        }
    }

    perf.WinRate = float64(winCount) / float64(len(deals)) * 100
    perf.SymbolsTraded = len(symbols)

    if winCount > 0 {
        perf.AvgWin = winSum / float64(winCount)
    }

    lossCount := len(deals) - winCount
    if lossCount > 0 {
        perf.AvgLoss = lossSum / float64(lossCount)
    }

    if lossSum != 0 {
        perf.ProfitFactor = winSum / -lossSum
    }

    return perf, nil
}

func (p *YesterdayPerformance) Print() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║    YESTERDAY'S PERFORMANCE ANALYSIS   ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    if p.TotalDeals == 0 {
        fmt.Println("No trading activity yesterday")
        return
    }

    fmt.Printf("Total Deals:     %d\n", p.TotalDeals)
    fmt.Printf("Total Profit:    $%.2f\n", p.TotalProfit)
    fmt.Printf("Win Rate:        %.1f%%\n", p.WinRate)
    fmt.Printf("Best Trade:      $%.2f\n", p.BestTrade)
    fmt.Printf("Worst Trade:     $%.2f\n", p.WorstTrade)
    fmt.Printf("Avg Win:         $%.2f\n", p.AvgWin)
    fmt.Printf("Avg Loss:        $%.2f\n", p.AvgLoss)
    fmt.Printf("Profit Factor:   %.2f\n", p.ProfitFactor)
    fmt.Printf("Symbols Traded:  %d\n", p.SymbolsTraded)

    fmt.Println("\nAssessment:")
    if p.TotalProfit > 0 && p.WinRate >= 60 && p.ProfitFactor >= 1.5 {
        fmt.Println("🌟 Excellent performance")
    } else if p.TotalProfit > 0 && p.WinRate >= 50 {
        fmt.Println("✅ Good performance")
    } else if p.TotalProfit >= 0 {
        fmt.Println("🟡 Breakeven - could improve")
    } else {
        fmt.Println("❌ Needs improvement")
    }
}

// Usage:
perf, _ := AnalyzeYesterday(sugar)
perf.Print()
```

---

## 🔗 Related Methods

**🍬 Other time periods:**

* `GetDealsToday()` - Today's deals (00:00 to now)
* `GetDealsThisWeek()` - This week's deals
* `GetDealsThisMonth()` - This month's deals
* `GetDealsDateRange()` - Custom date range

**🍬 Profit calculations:**

* `GetProfitToday()` - Today's profit total
* `GetProfitThisWeek()` - This week's profit

---

## ⚠️ Common Pitfalls

### 1) Calling too early in the day

```go
// ❌ WRONG - calling at 00:01 AM
// Yesterday's data might not be finalized yet
deals, _ := sugar.GetDealsYesterday()

// ✅ CORRECT - call after market has been open for a while
// Or in morning report after server has processed
time.Sleep(1 * time.Hour) // Wait until reasonable hour
deals, _ := sugar.GetDealsYesterday()
```

### 2) Timezone confusion

```go
// ❌ WRONG - "yesterday" is server timezone
// Your local "yesterday" might be different
deals, _ := sugar.GetDealsYesterday()

// ✅ CORRECT - aware of timezone
// Yesterday = server's yesterday (full 24h period)
deals, _ := sugar.GetDealsYesterday()
fmt.Println("Server's yesterday (full day)")
```

### 3) Not handling empty result

```go
// ❌ WRONG - assuming deals exist
deals, _ := sugar.GetDealsYesterday()
firstDeal := deals[0] // Panic if no deals!

// ✅ CORRECT - check length
if len(deals) > 0 {
    firstDeal := deals[0]
}
```

---

## 💎 Pro Tips

1. **Full day** - Returns complete 24-hour period (00:00 to 23:59:59)

2. **Morning reports** - Perfect for generating morning summary reports

3. **Comparison** - Easy to compare with today using `GetDealsToday()`

4. **Empty OK** - Returns empty slice if no trading (not an error)

5. **Performance** - 5-second timeout, suitable for frequent calls

---

## 📊 Yesterday's Date Range

```
Yesterday's range calculation:
now = time.Now()
yesterday = now.AddDate(0, 0, -1)
startOfYesterday = yesterday at 00:00:00
endOfYesterday = yesterday at 23:59:59

Returns all deals in this complete 24-hour window
```

---

**See also:** [`GetDealsToday.md`](GetDealsToday.md), [`GetDealsThisWeek.md`](GetDealsThisWeek.md), [`GetDealsDateRange.md`](GetDealsDateRange.md)
