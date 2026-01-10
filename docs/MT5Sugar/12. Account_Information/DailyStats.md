# MT5Sugar · DailyStats Structure

> **Daily trading statistics structure** returned by GetDailyStats(). Contains comprehensive performance metrics for today's trading session.

## 📋 Structure Definition

```go
type DailyStats struct {
    TotalDeals   int       // Total closed deals today
    WinningDeals int       // Number of profitable deals
    LosingDeals  int       // Number of losing deals
    WinRate      float64   // Win rate % (0-100)
    TotalProfit  float64   // Total realized profit/loss today
    BestDeal     float64   // Largest profitable deal
    WorstDeal    float64   // Largest losing deal
}
```

## 📊 Field Details

### TotalDeals (int)

- **What:** Total number of closed positions today (00:00 to now)
- **Example:** `15`
- **Includes:** Both winning and losing trades
- **Excludes:** Open positions (not yet closed)

### WinningDeals (int)

- **What:** Number of profitable closed positions
- **Example:** `9`
- **Criteria:** Profit > 0 (excludes zero-profit trades)
- **Use:** Calculate win rate, performance tracking

### LosingDeals (int)

- **What:** Number of losing closed positions
- **Example:** `6`
- **Criteria:** Profit < 0
- **Note:** TotalDeals = WinningDeals + LosingDeals (+ zero-profit deals)

### WinRate (float64)

- **What:** Percentage of winning trades
- **Formula:** (WinningDeals / TotalDeals) × 100
- **Range:** 0.0 to 100.0
- **Example:** `60.0` (means 60% win rate)
- **Important:**
  - Above 60% = Excellent
  - 50-60% = Good
  - Below 40% = Concerning

### TotalProfit (float64)

- **What:** Sum of all realized profits/losses today
- **Formula:** Sum of all deal profits
- **Example:** `350.75` (positive = profit, negative = loss)
- **Currency:** Account base currency
- **Excludes:** Floating P/L from open positions

### BestDeal (float64)

- **What:** Largest single profitable trade today
- **Example:** `125.50`
- **Use:** Track maximum gain, outlier detection
- **Note:** Always positive (or 0 if no winners)

### WorstDeal (float64)

- **What:** Largest single losing trade today
- **Example:** `-87.25`
- **Use:** Risk management, maximum loss tracking
- **Note:** Always negative (or 0 if no losers)

## 🎯 Common Calculations

### Average Win
```go
avgWin := 0.0
if stats.WinningDeals > 0 {
    // Approximation (need deal details for exact)
    avgWin = stats.BestDeal / float64(stats.WinningDeals)
}
```

### Average Loss
```go
avgLoss := 0.0
if stats.LosingDeals > 0 {
    avgLoss = stats.WorstDeal / float64(stats.LosingDeals)
}
```

### Risk:Reward Ratio
```go
riskRewardRatio := 0.0
if avgLoss != 0 {
    riskRewardRatio = avgWin / -avgLoss
}
// Values: 1:2 = 2.0, 1:1.5 = 1.5, etc.
```

### Profit Per Trade
```go
profitPerTrade := 0.0
if stats.TotalDeals > 0 {
    profitPerTrade = stats.TotalProfit / float64(stats.TotalDeals)
}
```

### Loss Rate
```go
lossRate := 100.0 - stats.WinRate
```

### Profit Factor
```go
// Requires individual deal data
profitFactor := totalWins / -totalLosses
// Above 2.0 = excellent, 1.5-2.0 = good, below 1.0 = losing
```

## 🟢 Usage Examples

### Example 1: Basic stats display

```go
stats, _ := sugar.GetDailyStats()

fmt.Println("📊 TODAY'S PERFORMANCE")
fmt.Println("=====================")
fmt.Printf("Total trades:   %d\n", stats.TotalDeals)
fmt.Printf("Winners:        %d\n", stats.WinningDeals)
fmt.Printf("Losers:         %d\n", stats.LosingDeals)
fmt.Printf("Win rate:       %.1f%%\n", stats.WinRate)
fmt.Printf("Total profit:   $%.2f\n", stats.TotalProfit)
fmt.Printf("Best trade:     $%.2f\n", stats.BestDeal)
fmt.Printf("Worst trade:    $%.2f\n", stats.WorstDeal)
```

### Example 2: Performance rating

```go
stats, _ := sugar.GetDailyStats()

rating := "Unknown"
if stats.TotalDeals >= 5 {
    if stats.WinRate >= 70 && stats.TotalProfit > 0 {
        rating = "⭐⭐⭐ Exceptional"
    } else if stats.WinRate >= 60 && stats.TotalProfit > 0 {
        rating = "⭐⭐ Excellent"
    } else if stats.WinRate >= 50 && stats.TotalProfit > 0 {
        rating = "⭐ Good"
    } else if stats.TotalProfit > 0 {
        rating = "👍 Profitable"
    } else if stats.TotalProfit > -500 {
        rating = "📉 Minor loss"
    } else {
        rating = "🔴 Significant loss"
    }
} else {
    rating = "⏳ Insufficient data (need 5+ trades)"
}

fmt.Printf("Today's rating: %s\n", rating)
```

### Example 3: Risk analysis

```go
stats, _ := sugar.GetDailyStats()

fmt.Println("⚠️ RISK ANALYSIS")
fmt.Println("================")

// Worst trade vs total profit ratio
if stats.TotalProfit > 0 {
    worstImpact := (-stats.WorstDeal / stats.TotalProfit) * 100
    fmt.Printf("Worst trade impact: %.1f%% of profit\n", worstImpact)

    if worstImpact > 50 {
        fmt.Println("🔴 Single loss took >50% of profit!")
    }
}

// Best vs worst comparison
if stats.WorstDeal != 0 {
    ratio := stats.BestDeal / -stats.WorstDeal
    fmt.Printf("Best:Worst ratio: 1:%.2f\n", ratio)

    if ratio < 1.0 {
        fmt.Println("⚠️ Largest loss exceeds largest win")
    }
}
```

### Example 4: Trading consistency

```go
stats, _ := sugar.GetDailyStats()

if stats.TotalDeals == 0 {
    fmt.Println("No trades today")
    return
}

profitPerTrade := stats.TotalProfit / float64(stats.TotalDeals)
consistency := "Unknown"

if profitPerTrade > 50 {
    consistency = "🔥 Highly profitable per trade"
} else if profitPerTrade > 20 {
    consistency = "✅ Good average profit"
} else if profitPerTrade > 0 {
    consistency = "👍 Slightly profitable"
} else if profitPerTrade > -20 {
    consistency = "📉 Small average loss"
} else {
    consistency = "🔴 Large average loss"
}

fmt.Printf("Avg profit/trade: $%.2f (%s)\n", profitPerTrade, consistency)
```

### Example 5: Win/loss distribution

```go
stats, _ := sugar.GetDailyStats()

if stats.TotalDeals == 0 {
    return
}

// Visual representation
winPercent := int(stats.WinRate)
lossPercent := 100 - winPercent

fmt.Println("Win/Loss Distribution:")
fmt.Printf("[")
for i := 0; i < winPercent/2; i++ {
    fmt.Print("=")
}
fmt.Print("|")
for i := 0; i < lossPercent/2; i++ {
    fmt.Print("-")
}
fmt.Printf("]\n")
fmt.Printf("Wins: %d%% | Losses: %d%%\n", winPercent, lossPercent)
```

### Example 6: Alert triggers

```go
func checkAlerts(stats *mt5.DailyStats) {
    // Alert 1: Low win rate
    if stats.TotalDeals >= 10 && stats.WinRate < 40 {
        fmt.Println("🚨 ALERT: Win rate below 40%")
    }

    // Alert 2: Large loss
    if stats.WorstDeal < -500 {
        fmt.Printf("🚨 ALERT: Large single loss $%.2f\n", stats.WorstDeal)
    }

    // Alert 3: Daily loss limit
    if stats.TotalProfit < -1000 {
        fmt.Printf("🚨 ALERT: Daily loss $%.2f exceeds limit\n", stats.TotalProfit)
    }

    // Alert 4: Consecutive losses
    if stats.LosingDeals >= 5 && stats.WinningDeals == 0 {
        fmt.Println("🚨 ALERT: 5+ consecutive losses")
    }

    // Alert 5: Excellent performance
    if stats.TotalDeals >= 10 && stats.WinRate >= 70 && stats.TotalProfit > 500 {
        fmt.Println("🎉 EXCELLENT: Win rate >70% with >$500 profit!")
    }
}
```

### Example 7: JSON export

```go
stats, _ := sugar.GetDailyStats()

// Pretty print JSON
jsonData, _ := json.MarshalIndent(stats, "", "  ")
fmt.Println(string(jsonData))

// Output:
// {
//   "TotalDeals": 15,
//   "WinningDeals": 9,
//   "LosingDeals": 6,
//   "WinRate": 60.0,
//   "TotalProfit": 350.75,
//   "BestDeal": 125.50,
//   "WorstDeal": -87.25
// }
```

### Example 8: CSV export

```go
func exportStatsCSV(stats *mt5.DailyStats) {
    date := time.Now().Format("2006-01-02")

    csvLine := fmt.Sprintf("%s,%d,%d,%d,%.2f,%.2f,%.2f,%.2f\n",
        date,
        stats.TotalDeals,
        stats.WinningDeals,
        stats.LosingDeals,
        stats.WinRate,
        stats.TotalProfit,
        stats.BestDeal,
        stats.WorstDeal,
    )

    // Append to CSV file
    f, _ := os.OpenFile("daily_stats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()

    f.WriteString(csvLine)
    fmt.Println("Stats exported to daily_stats.csv")
}
```

### Example 9: Telegram notification

```go
func sendTelegramReport(stats *mt5.DailyStats) string {
    emoji := "📊"
    if stats.TotalProfit > 500 {
        emoji = "🚀"
    } else if stats.TotalProfit < -500 {
        emoji = "🔴"
    }

    message := fmt.Sprintf(`
%s *Daily Trading Report*

📈 Performance:
• Trades: %d
• Win rate: %.1f%%
• P/L: $%.2f

🏆 Best trade: $%.2f
📉 Worst trade: $%.2f

Status: %s
`,
        emoji,
        stats.TotalDeals,
        stats.WinRate,
        stats.TotalProfit,
        stats.BestDeal,
        stats.WorstDeal,
        getStatus(stats),
    )

    return message
}

func getStatus(stats *mt5.DailyStats) string {
    if stats.TotalProfit > 0 && stats.WinRate >= 60 {
        return "Excellent ✅"
    } else if stats.TotalProfit > 0 {
        return "Profitable 👍"
    } else {
        return "Needs improvement 📉"
    }
}
```

### Example 10: Historical tracking

```go
type StatsHistory struct {
    Date  string
    Stats *mt5.DailyStats
}

var history []StatsHistory

func recordDailyStats(stats *mt5.DailyStats) {
    record := StatsHistory{
        Date:  time.Now().Format("2006-01-02"),
        Stats: stats,
    }

    history = append(history, record)

    // Calculate weekly averages (last 7 days)
    if len(history) >= 7 {
        last7 := history[len(history)-7:]

        totalTrades := 0
        totalProfit := 0.0
        totalWinRate := 0.0

        for _, rec := range last7 {
            totalTrades += rec.Stats.TotalDeals
            totalProfit += rec.Stats.TotalProfit
            totalWinRate += rec.Stats.WinRate
        }

        fmt.Println("\n📅 WEEKLY SUMMARY (Last 7 days)")
        fmt.Printf("Total trades:  %d\n", totalTrades)
        fmt.Printf("Total profit:  $%.2f\n", totalProfit)
        fmt.Printf("Avg win rate:  %.1f%%\n", totalWinRate/7)
    }
}
```

## 🔗 Related

- **[AccountInfo](./AccountInfo.md)** - account information structure
- **[GetDealsToday](../8.%20History_Statistics/GetDealsToday.md)** - raw deal data

## ⚠️ Important Notes

1. **Today's range** - 00:00 to current time (server time, not local)

2. **Closed positions only** - Open positions not included

3. **Win rate** - Based on closed deals, not tick-by-tick P/L

4. **Best/Worst** - Absolute values, not averages

5. **Zero-profit deals** - Excluded from win rate calculation

## 💡 Quick Reference

| Field | Range | Good value | Bad value |
|-------|-------|------------|-----------|
| `WinRate` | 0-100% | >60% | <40% |
| `TotalProfit` | Any | >$0 | <-$500 |
| `BestDeal` | ≥0 | Large | Small |
| `WorstDeal` | ≤0 | Small | Large |
| `TotalDeals` | ≥0 | 5-20 | >50 (overtrading) |

## 🎓 Performance Benchmarks

### Excellent Day

- Win rate: >70%
- Total profit: >$500
- Profit per trade: >$50
- Best:Worst ratio: >2:1

### Good Day

- Win rate: 55-70%
- Total profit: >$200
- Profit per trade: >$20
- Best:Worst ratio: >1.5:1

### Concerning Day

- Win rate: <40%
- Total profit: <-$500
- Profit per trade: <-$20
- Best:Worst ratio: <1:1

---

**Summary:** DailyStats structure provides comprehensive daily trading metrics. Use WinRate and TotalProfit for performance assessment, BestDeal/WorstDeal for risk analysis. Essential for daily reports and trading journals.
