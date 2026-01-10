# 📊📈 Get Daily Stats (`GetDailyStats`)

> **Sugar method:** Returns comprehensive trading statistics for today (all-in-one daily performance).

**API Information:**

* **Method:** `sugar.GetDailyStats()`
* **Timeout:** 5 seconds
* **Returns:** DailyStats structure with complete statistics

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetDailyStats() (*DailyStats, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates today) |

| Output | Type | Description |
|--------|------|-------------|
| `*DailyStats` | struct | Complete daily statistics structure |
| `error` | `error` | Error if calculation failed |

---

## 📊 DailyStats Structure

```go
type DailyStats struct {
    TotalDeals   int     // Total closed deals today
    WinningDeals int     // Number of profitable deals
    LosingDeals  int     // Number of losing deals
    WinRate      float64 // Win rate percentage (0-100)
    TotalProfit  float64 // Total realized P/L today
    BestDeal     float64 // Largest profitable deal
    WorstDeal    float64 // Largest losing deal
}
```

---

## 💬 Just the Essentials

* **What it is:** Get complete trading statistics for today in one convenient structure.
* **Why you need it:** Daily performance reports, dashboard displays, comprehensive analysis.
* **Sanity check:** Returns empty stats if no deals today (not error). All metrics calculated automatically.

---

## 🎯 When to Use

✅ **Daily reports** - Generate comprehensive daily summary

✅ **Dashboards** - Display all key metrics at once

✅ **Performance analysis** - Detailed today's performance review

✅ **Trading journal** - Record complete daily statistics

---

## 🔗 Usage Examples

### 1) Basic usage - show today's stats

```go
stats, err := sugar.GetDailyStats()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if stats.TotalDeals == 0 {
    fmt.Println("No trades today")
    return
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      TODAY'S TRADING STATISTICS       ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Total Deals:     %d\n", stats.TotalDeals)
fmt.Printf("Winners:         %d\n", stats.WinningDeals)
fmt.Printf("Losers:          %d\n", stats.LosingDeals)
fmt.Printf("Win Rate:        %.1f%%\n\n", stats.WinRate)

fmt.Printf("Total Profit:    $%.2f\n", stats.TotalProfit)
fmt.Printf("Best Trade:      $%.2f\n", stats.BestDeal)
fmt.Printf("Worst Trade:     $%.2f\n", stats.WorstDeal)
```

---

### 2) Daily report with assessment

```go
stats, _ := sugar.GetDailyStats()

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      DAILY PERFORMANCE REPORT         ║")
fmt.Println("╚═══════════════════════════════════════╝")

if stats.TotalDeals == 0 {
    fmt.Println("No trading activity today")
    return
}

// Trading volume
fmt.Printf("Trading Activity:\n")
fmt.Printf("  Total trades:    %d\n", stats.TotalDeals)
fmt.Printf("  Winners:         %d (%.1f%%)\n",
    stats.WinningDeals, stats.WinRate)
fmt.Printf("  Losers:          %d (%.1f%%)\n\n",
    stats.LosingDeals, 100-stats.WinRate)

// Profitability
fmt.Printf("Profitability:\n")
fmt.Printf("  Total P/L:       $%.2f\n", stats.TotalProfit)

if stats.TotalDeals > 0 {
    avgProfit := stats.TotalProfit / float64(stats.TotalDeals)
    fmt.Printf("  Avg per trade:   $%.2f\n\n", avgProfit)
}

// Best/Worst
fmt.Printf("Extremes:\n")
fmt.Printf("  Best trade:      $%.2f\n", stats.BestDeal)
fmt.Printf("  Worst trade:     $%.2f\n\n", stats.WorstDeal)

// Assessment
fmt.Println("Assessment:")
if stats.TotalProfit > 0 && stats.WinRate >= 60 {
    fmt.Println("  🌟 Excellent day!")
} else if stats.TotalProfit > 0 && stats.WinRate >= 50 {
    fmt.Println("  ✅ Good day")
} else if stats.TotalProfit >= 0 {
    fmt.Println("  🟡 Breakeven - could improve")
} else if stats.WinRate >= 50 {
    fmt.Println("  ⚠️  Losing despite good win rate")
} else {
    fmt.Println("  ❌ Poor day - review trades")
}
```

---

### 3) Compare with target metrics

```go
stats, _ := sugar.GetDailyStats()

// Target metrics
targetTrades := 10
targetWinRate := 60.0
targetProfit := 100.0

fmt.Println("Target vs Actual:")
fmt.Println("─────────────────────────────────────────")

// Trades
if stats.TotalDeals >= targetTrades {
    fmt.Printf("✅ Trades: %d / %d\n", stats.TotalDeals, targetTrades)
} else {
    fmt.Printf("❌ Trades: %d / %d (need %d more)\n",
        stats.TotalDeals, targetTrades, targetTrades-stats.TotalDeals)
}

// Win rate
if stats.WinRate >= targetWinRate {
    fmt.Printf("✅ Win rate: %.1f%% / %.1f%%\n", stats.WinRate, targetWinRate)
} else {
    fmt.Printf("❌ Win rate: %.1f%% / %.1f%%\n", stats.WinRate, targetWinRate)
}

// Profit
if stats.TotalProfit >= targetProfit {
    fmt.Printf("✅ Profit: $%.2f / $%.2f\n", stats.TotalProfit, targetProfit)
} else {
    fmt.Printf("❌ Profit: $%.2f / $%.2f\n", stats.TotalProfit, targetProfit)
}
```

---

### 4) Real-time statistics monitor

```go
func MonitorDailyStats(sugar *mt5.MT5Sugar, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    previousDeals := 0

    for range ticker.C {
        stats, _ := sugar.GetDailyStats()

        if stats.TotalDeals > previousDeals {
            fmt.Printf("\n[%s] New trade closed!\n",
                time.Now().Format("15:04:05"))

            fmt.Printf("  Today's stats: %d trades, %.1f%% win, $%.2f\n",
                stats.TotalDeals, stats.WinRate, stats.TotalProfit)

            previousDeals = stats.TotalDeals
        }
    }
}

// Usage: Monitor every 30 seconds
go MonitorDailyStats(sugar, 30*time.Second)
```

---

### 5) Risk/reward ratio analysis

```go
stats, _ := sugar.GetDailyStats()

if stats.TotalDeals == 0 {
    fmt.Println("No trades today")
    return
}

// Calculate average win/loss
avgWin := 0.0
if stats.WinningDeals > 0 {
    // Estimate avg win
    totalWins := stats.TotalProfit - (stats.WorstDeal * float64(stats.LosingDeals))
    avgWin = totalWins / float64(stats.WinningDeals)
}

avgLoss := 0.0
if stats.LosingDeals > 0 {
    // Estimate avg loss
    totalLosses := stats.WorstDeal * float64(stats.LosingDeals)
    avgLoss = totalLosses / float64(stats.LosingDeals)
}

riskRewardRatio := 0.0
if avgLoss != 0 {
    riskRewardRatio = avgWin / -avgLoss
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      RISK/REWARD ANALYSIS             ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Average Win:      $%.2f\n", avgWin)
fmt.Printf("Average Loss:     $%.2f\n", avgLoss)
fmt.Printf("Risk/Reward:      %.2f:1\n\n", riskRewardRatio)

if riskRewardRatio >= 2.0 {
    fmt.Println("✅ Excellent risk management")
} else if riskRewardRatio >= 1.5 {
    fmt.Println("✅ Good risk/reward")
} else if riskRewardRatio >= 1.0 {
    fmt.Println("🟡 Acceptable")
} else {
    fmt.Println("⚠️  Poor risk/reward - review trades")
}
```

---

### 6) Performance grade calculator

```go
func GradeDailyPerformance(stats *DailyStats) string {
    if stats.TotalDeals == 0 {
        return "N/A - No trades"
    }

    score := 0.0

    // Profitability (50 points)
    if stats.TotalProfit > 0 {
        score += 50
    }

    // Win rate (30 points)
    if stats.WinRate >= 70 {
        score += 30
    } else if stats.WinRate >= 60 {
        score += 25
    } else if stats.WinRate >= 50 {
        score += 20
    } else if stats.WinRate >= 40 {
        score += 10
    }

    // Trading volume (20 points)
    if stats.TotalDeals >= 10 {
        score += 20
    } else if stats.TotalDeals >= 5 {
        score += 10
    }

    // Assign grade
    if score >= 90 {
        return "🌟 A+ (Excellent)"
    } else if score >= 80 {
        return "⭐ A (Very Good)"
    } else if score >= 70 {
        return "✅ B (Good)"
    } else if score >= 60 {
        return "🟡 C (Fair)"
    } else if score >= 50 {
        return "⚠️  D (Poor)"
    } else {
        return "❌ F (Fail)"
    }
}

// Usage:
stats, _ := sugar.GetDailyStats()
grade := GradeDailyPerformance(stats)
fmt.Printf("Today's Grade: %s\n", grade)
```

---

### 7) Export to CSV for journal

```go
func ExportDailyStatsToCSV(sugar *mt5.MT5Sugar) error {
    stats, err := sugar.GetDailyStats()
    if err != nil {
        return err
    }

    filename := fmt.Sprintf("daily_stats_%s.csv",
        time.Now().Format("2006-01-02"))

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Check if file is new (write header)
    fileInfo, _ := file.Stat()
    if fileInfo.Size() == 0 {
        writer.Write([]string{
            "Date", "Total Deals", "Winning Deals", "Losing Deals",
            "Win Rate", "Total Profit", "Best Deal", "Worst Deal"})
    }

    // Write stats
    writer.Write([]string{
        time.Now().Format("2006-01-02"),
        fmt.Sprintf("%d", stats.TotalDeals),
        fmt.Sprintf("%d", stats.WinningDeals),
        fmt.Sprintf("%d", stats.LosingDeals),
        fmt.Sprintf("%.2f", stats.WinRate),
        fmt.Sprintf("%.2f", stats.TotalProfit),
        fmt.Sprintf("%.2f", stats.BestDeal),
        fmt.Sprintf("%.2f", stats.WorstDeal),
    })

    fmt.Printf("✅ Exported stats to %s\n", filename)
    return nil
}

// Usage:
ExportDailyStatsToCSV(sugar)
```

---

### 8) Consistency checker

```go
func CheckDailyConsistency(sugar *mt5.MT5Sugar) {
    stats, _ := sugar.GetDailyStats()

    if stats.TotalDeals == 0 {
        fmt.Println("No trades today")
        return
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      CONSISTENCY CHECK                ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    issues := []string{}

    // Check win rate consistency
    if stats.WinRate < 40 {
        issues = append(issues, "Win rate too low (<40%)")
    }

    // Check if best/worst are balanced
    if stats.BestDeal > 0 && stats.WorstDeal < 0 {
        ratio := stats.BestDeal / -stats.WorstDeal
        if ratio < 1.0 {
            issues = append(issues, "Best win smaller than worst loss")
        }
    }

    // Check trading volume
    if stats.TotalDeals < 5 {
        issues = append(issues, "Low trading volume (<5 trades)")
    } else if stats.TotalDeals > 50 {
        issues = append(issues, "Very high volume (>50 trades) - overtrading?")
    }

    // Check profitability
    if stats.TotalProfit < 0 {
        issues = append(issues, "Negative day")
    }

    // Report
    if len(issues) == 0 {
        fmt.Println("✅ All consistency checks passed")
    } else {
        fmt.Println("⚠️  Issues found:")
        for _, issue := range issues {
            fmt.Printf("  - %s\n", issue)
        }
    }
}

// Usage:
CheckDailyConsistency(sugar)
```

---

### 9) Streak tracker

```go
type StreakTracker struct {
    WinStreak  int
    LossStreak int
}

func (st *StreakTracker) Update(stats *DailyStats) {
    if stats.TotalProfit > 0 {
        st.WinStreak++
        st.LossStreak = 0
    } else if stats.TotalProfit < 0 {
        st.LossStreak++
        st.WinStreak = 0
    }
}

func (st *StreakTracker) Show() {
    fmt.Println("Current Streak:")
    if st.WinStreak > 0 {
        fmt.Printf("  🔥 %d winning days\n", st.WinStreak)
        if st.WinStreak >= 5 {
            fmt.Println("  🌟 Excellent streak!")
        }
    } else if st.LossStreak > 0 {
        fmt.Printf("  ❄️  %d losing days\n", st.LossStreak)
        if st.LossStreak >= 3 {
            fmt.Println("  ⚠️  Consider reviewing strategy")
        }
    } else {
        fmt.Println("  ➡️  No active streak")
    }
}

// Usage:
tracker := &StreakTracker{}

stats, _ := sugar.GetDailyStats()
tracker.Update(stats)
tracker.Show()
```

---

### 10) Advanced daily statistics analyzer

```go
type DailyStatsAnalyzer struct {
    Stats       *DailyStats
    TargetStats *DailyStats
}

func NewDailyStatsAnalyzer(sugar *mt5.MT5Sugar, targetProfit float64, targetWinRate float64) (*DailyStatsAnalyzer, error) {
    stats, err := sugar.GetDailyStats()
    if err != nil {
        return nil, err
    }

    return &DailyStatsAnalyzer{
        Stats: stats,
        TargetStats: &DailyStats{
            TotalProfit: targetProfit,
            WinRate:     targetWinRate,
        },
    }, nil
}

func (dsa *DailyStatsAnalyzer) GenerateReport() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      COMPREHENSIVE DAILY ANALYSIS     ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    if dsa.Stats.TotalDeals == 0 {
        fmt.Println("No trading activity today")
        return
    }

    // Overview
    fmt.Println("Trading Overview:")
    fmt.Printf("  Total trades:    %d\n", dsa.Stats.TotalDeals)
    fmt.Printf("  Win rate:        %.1f%%\n", dsa.Stats.WinRate)
    fmt.Printf("  Total P/L:       $%.2f\n\n", dsa.Stats.TotalProfit)

    // Performance vs Target
    fmt.Println("Target Comparison:")

    profitDiff := dsa.Stats.TotalProfit - dsa.TargetStats.TotalProfit
    if profitDiff >= 0 {
        fmt.Printf("  ✅ Profit: $%.2f over target\n", profitDiff)
    } else {
        fmt.Printf("  ❌ Profit: $%.2f under target\n", -profitDiff)
    }

    winRateDiff := dsa.Stats.WinRate - dsa.TargetStats.WinRate
    if winRateDiff >= 0 {
        fmt.Printf("  ✅ Win rate: %.1f%% over target\n\n", winRateDiff)
    } else {
        fmt.Printf("  ❌ Win rate: %.1f%% under target\n\n", -winRateDiff)
    }

    // Trade quality
    fmt.Println("Trade Quality:")

    avgProfit := dsa.Stats.TotalProfit / float64(dsa.Stats.TotalDeals)
    fmt.Printf("  Avg per trade:   $%.2f\n", avgProfit)

    if dsa.Stats.BestDeal > 0 && dsa.Stats.WorstDeal < 0 {
        ratio := dsa.Stats.BestDeal / -dsa.Stats.WorstDeal
        fmt.Printf("  Best/Worst:      %.2f:1\n", ratio)
    }

    range_ := dsa.Stats.BestDeal - dsa.Stats.WorstDeal
    fmt.Printf("  P/L range:       $%.2f\n\n", range_)

    // Recommendations
    fmt.Println("Recommendations:")

    if dsa.Stats.WinRate < 50 {
        fmt.Println("  ⚠️  Win rate below 50% - review entry criteria")
    }

    if avgProfit < 0 {
        fmt.Println("  ⚠️  Negative average - review risk management")
    }

    if dsa.Stats.TotalDeals < 5 {
        fmt.Println("  📊 Low volume - consider more opportunities")
    } else if dsa.Stats.TotalDeals > 30 {
        fmt.Println("  ⚠️  High volume - risk of overtrading")
    }

    if dsa.Stats.TotalProfit >= dsa.TargetStats.TotalProfit &&
       dsa.Stats.WinRate >= dsa.TargetStats.WinRate {
        fmt.Println("  🌟 All targets met - excellent day!")
    }
}

// Usage:
analyzer, _ := NewDailyStatsAnalyzer(sugar, 100.0, 60.0)
analyzer.GenerateReport()
```

---

## 🔗 Related Methods

**🍬 Individual stats:**

* `GetDealsToday()` - Get full deal information
* `GetProfitToday()` - Just profit total

**🍬 Other periods:**

* Create similar stats for week/month using GetDeals* methods

---

## ⚠️ Common Pitfalls

### 1) Not checking for zero trades

```go
// ❌ WRONG - division by zero!
stats, _ := sugar.GetDailyStats()
avgProfit := stats.TotalProfit / stats.TotalDeals

// ✅ CORRECT - check first
if stats.TotalDeals > 0 {
    avgProfit := stats.TotalProfit / float64(stats.TotalDeals)
}
```

### 2) Confusing stats with floating profit

```go
// ❌ WRONG - DailyStats is CLOSED positions only
stats, _ := sugar.GetDailyStats()
// Does NOT include open positions!

// ✅ CORRECT - get open profit separately
floatingProfit, _ := sugar.GetProfit()
totalProfit := stats.TotalProfit + floatingProfit
```

### 3) Assuming BestDeal/WorstDeal are set

```go
// ❌ WRONG - might be zero if no wins/losses
stats, _ := sugar.GetDailyStats()
ratio := stats.BestDeal / stats.WorstDeal // Panic if WorstDeal = 0!

// ✅ CORRECT - check first
if stats.BestDeal > 0 && stats.WorstDeal < 0 {
    ratio := stats.BestDeal / -stats.WorstDeal
}
```

---

## 💎 Pro Tips

1. **All-in-one** - Single call gets complete daily picture

2. **Calculated fields** - Win rate automatically calculated

3. **Empty OK** - Returns empty stats if no trades (not error)

4. **Today only** - Stats from 00:00 server time to now

5. **Efficient** - More efficient than manually calculating from deals

---

## 📊 Field Relationships

```
Total relationships:
- TotalDeals = WinningDeals + LosingDeals + BreakevenDeals
- WinRate = (WinningDeals / TotalDeals) * 100
- TotalProfit = Sum of all deal profits

Quality checks:
- If WinRate > 50% but TotalProfit < 0: Wins too small
- If WinRate < 50% but TotalProfit > 0: Wins are large
- BestDeal / |WorstDeal| > 1.5: Good risk management
```

---

**See also:** [`GetDealsToday.md`](GetDealsToday.md), [`GetProfitToday.md`](GetProfitToday.md)
