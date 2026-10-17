# 📅💼 Get Deals Today (`GetDealsToday`)

> **Sugar method:** Returns closed positions ("deals") that were **opened** today (00:00 to now).

**API Information:**

* **Method:** `sugar.GetDealsToday()`
* **Timeout:** 5 seconds
* **Returns:** Slice of position history info, sorted by open time (oldest first)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetDealsToday() ([]*pb.PositionHistoryInfo, error)
```

How it works: `now := time.Now()`, `startOfDay` = today 00:00 on your machine's clock, then
`MT5Service.GetPositionsHistory(ctx, pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC, &startOfDay, &now, nil, nil)`.
The range is sent as `PositionOpenTimeFrom` / `PositionOpenTimeTo`, so it filters by the position's **open** time:
a position opened yesterday and closed today is not included. No paging parameters are sent.

**Time zone note:** the boundaries come from your local clock and are sent as absolute instants in the request fields that
the server reads as broker server time. If your machine, UTC and the broker server use different time zones, the window is
shifted accordingly - use `GetDealsDateRange()` when you need exact boundaries.

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters (auto-calculates today's range) |

| Output | Type | Description |
|--------|------|-------------|
| `[]*pb.PositionHistoryInfo` | slice | Closed positions opened today |
| `error` | `error` | `"GetPositionsHistory failed: ..."` if retrieval failed |

---

## 💬 Just the Essentials

* **What it is:** Get closed trades opened today (00:00 to current time, local clock).
* **Why you need it:** Daily performance tracking, today's profit/loss analysis, trade review.
* **Sanity check:** Returns an empty (possibly `nil`) slice if there are no deals - check `len(deals)`. Automatically calculates today's time range.

---

## 🎯 When to Use

✅ **Daily reports** - Generate today's trading summary

✅ **Performance tracking** - Monitor today's profit/loss

✅ **Trade journal** - Review today's completed trades

✅ **Real-time updates** - Check latest closed positions

---

## 🔗 Usage Examples

### 1) Basic usage - show today's deals

```go
deals, err := sugar.GetDealsToday()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if len(deals) == 0 {
    fmt.Println("No deals today")
    return
}

fmt.Printf("Today's deals: %d\n\n", len(deals))

for i, deal := range deals {
    fmt.Printf("%d. Ticket #%d: %s %.2f lots\n",
        i+1, deal.PositionTicket, deal.Symbol, deal.Volume)
    fmt.Printf("   Profit: $%.2f\n", deal.Profit)
}
```

---

### 2) Calculate today's profit

```go
deals, _ := sugar.GetDealsToday()

if len(deals) == 0 {
    fmt.Println("No trades today")
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

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      TODAY'S PERFORMANCE              ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Total trades: %d\n", len(deals))
fmt.Printf("Winners:      %d (%.1f%%)\n", winCount,
    float64(winCount)/float64(len(deals))*100)
fmt.Printf("Losers:       %d (%.1f%%)\n", lossCount,
    float64(lossCount)/float64(len(deals))*100)
fmt.Printf("Total P/L:    $%.2f\n", totalProfit)

if totalProfit > 0 {
    fmt.Println("✅ Profitable day")
} else {
    fmt.Println("❌ Losing day")
}
```

---

### 3) Today's trading summary by symbol

```go
deals, _ := sugar.GetDealsToday()

if len(deals) == 0 {
    fmt.Println("No deals today")
    return
}

// Group by symbol
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

fmt.Println("Today's Performance by Symbol:")
fmt.Println("─────────────────────────────────────────")

for symbol, stats := range symbolStats {
    fmt.Printf("%-8s: %d trades, $%.2f\n",
        symbol, stats.Count, stats.Profit)
}
```

---

### 4) Find largest win/loss today

```go
deals, _ := sugar.GetDealsToday()

if len(deals) == 0 {
    fmt.Println("No deals today")
    return
}

var largestWin *pb.PositionHistoryInfo
var largestLoss *pb.PositionHistoryInfo

for _, deal := range deals {
    if largestWin == nil || deal.Profit > largestWin.Profit {
        largestWin = deal
    }

    if largestLoss == nil || deal.Profit < largestLoss.Profit {
        largestLoss = deal
    }
}

fmt.Println("Today's Extremes:")
fmt.Println("─────────────────────────────────────────")

if largestWin != nil {
    fmt.Printf("🏆 Largest win:\n")
    fmt.Printf("   #%d %s %.2f lots: $%.2f\n",
        largestWin.PositionTicket, largestWin.Symbol,
        largestWin.Volume, largestWin.Profit)
}

if largestLoss != nil && largestLoss.Profit < 0 {
    fmt.Printf("\n💔 Largest loss:\n")
    fmt.Printf("   #%d %s %.2f lots: $%.2f\n",
        largestLoss.PositionTicket, largestLoss.Symbol,
        largestLoss.Volume, largestLoss.Profit)
}
```

---

### 5) Track hourly trading activity

```go
deals, _ := sugar.GetDealsToday()

if len(deals) == 0 {
    fmt.Println("No deals today")
    return
}

// Group by close hour (CloseTime is broker server time)
hourlyActivity := make(map[int]int)

for _, deal := range deals {
    hour := deal.CloseTime.AsTime().Hour()
    hourlyActivity[hour]++
}

fmt.Println("Hourly Trading Activity:")
fmt.Println("─────────────────────────────────────────")

for hour := 0; hour < 24; hour++ {
    count := hourlyActivity[hour]
    if count > 0 {
        fmt.Printf("%02d:00 - %02d:59: %d trades\n",
            hour, hour, count)
    }
}
```

---

### 6) Check if breakeven or profitable today

```go
deals, _ := sugar.GetDealsToday()

if len(deals) == 0 {
    fmt.Println("No trades today - day just started?")
    return
}

totalProfit := 0.0
for _, deal := range deals {
    totalProfit += deal.Profit
}

fmt.Printf("Today: %d trades, $%.2f\n", len(deals), totalProfit)

switch {
case totalProfit > 100:
    fmt.Println("🎉 Great day! Target exceeded")
case totalProfit > 0:
    fmt.Println("✅ Profitable day")
case totalProfit == 0:
    fmt.Println("➡️  Breakeven")
case totalProfit > -100:
    fmt.Println("⚠️  Small loss - recoverable")
default:
    fmt.Println("🛑 Significant loss - stop trading")
}
```

---

### 7) Real-time deal monitor

```go
func MonitorTodaysDeals(sugar *mt5.MT5Sugar, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    previousCount := 0

    for range ticker.C {
        deals, _ := sugar.GetDealsToday()
        currentCount := len(deals)

        if currentCount > previousCount {
            // New deals closed
            newDeals := currentCount - previousCount

            // Last element = most recently OPENED of today's closed positions
            // (the list is sorted by open time, not close time)
            if len(deals) > 0 {
                latest := deals[len(deals)-1]

                fmt.Printf("[%s] %d new deal(s) closed!\n",
                    time.Now().Format("15:04:05"), newDeals)
                fmt.Printf("   #%d %s: $%.2f\n",
                    latest.PositionTicket, latest.Symbol, latest.Profit)
            }
        }

        previousCount = currentCount
    }
}

// Usage: Monitor every 10 seconds
go MonitorTodaysDeals(sugar, 10*time.Second)
```

---

### 8) Daily performance grade

```go
deals, _ := sugar.GetDealsToday()

if len(deals) == 0 {
    fmt.Println("No trades today - grade: N/A")
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
fmt.Println("║      TODAY'S PERFORMANCE GRADE        ║")
fmt.Println("╚═══════════════════════════════════════╝")

fmt.Printf("Trades:    %d\n", len(deals))
fmt.Printf("Win rate:  %.1f%%\n", winRate)
fmt.Printf("P/L:       $%.2f\n\n", totalProfit)

// Grade based on win rate and profitability
grade := ""
if totalProfit > 0 && winRate >= 70 {
    grade = "🌟 A+ (Excellent)"
} else if totalProfit > 0 && winRate >= 60 {
    grade = "⭐ A (Very Good)"
} else if totalProfit > 0 && winRate >= 50 {
    grade = "✅ B (Good)"
} else if totalProfit >= 0 {
    grade = "🟡 C (Fair)"
} else if totalProfit > -100 {
    grade = "⚠️  D (Poor)"
} else {
    grade = "❌ F (Fail)"
}

fmt.Printf("Grade:     %s\n", grade)
```

---

### 9) Export today's deals to CSV

```go
deals, _ := sugar.GetDealsToday()

if len(deals) == 0 {
    fmt.Println("No deals to export")
    return
}

filename := fmt.Sprintf("deals_%s.csv",
    time.Now().Format("2006-01-02"))

file, _ := os.Create(filename)
defer file.Close()

writer := csv.NewWriter(file)
defer writer.Flush()

// Header
writer.Write([]string{
    "Ticket", "Symbol", "Volume", "Profit",
    "Open Time", "Close Time", "Duration"})

// Data
for _, deal := range deals {
    // OpenTime/CloseTime are *timestamppb.Timestamp in broker server time
    openTime := deal.OpenTime.AsTime()
    closeTime := deal.CloseTime.AsTime()
    duration := closeTime.Sub(openTime)

    writer.Write([]string{
        fmt.Sprintf("%d", deal.PositionTicket),
        deal.Symbol,
        fmt.Sprintf("%.2f", deal.Volume),
        fmt.Sprintf("%.2f", deal.Profit),
        openTime.Format("15:04:05"),
        closeTime.Format("15:04:05"),
        duration.String(),
    })
}

fmt.Printf("✅ Exported %d deals to %s\n", len(deals), filename)
```

---

### 10) Advanced daily statistics manager

```go
type DailyDealStats struct {
    TotalDeals    int
    WinningDeals  int
    LosingDeals   int
    TotalProfit   float64
    LargestWin    float64
    LargestLoss   float64
    AvgProfit     float64
    WinRate       float64
}

func GetDailyDealStats(sugar *mt5.MT5Sugar) (*DailyDealStats, error) {
    deals, err := sugar.GetDealsToday()
    if err != nil {
        return nil, err
    }

    stats := &DailyDealStats{
        TotalDeals: len(deals),
    }

    if len(deals) == 0 {
        return stats, nil
    }

    stats.LargestWin = 0
    stats.LargestLoss = 0

    for _, deal := range deals {
        stats.TotalProfit += deal.Profit

        if deal.Profit > 0 {
            stats.WinningDeals++
            if deal.Profit > stats.LargestWin {
                stats.LargestWin = deal.Profit
            }
        } else if deal.Profit < 0 {
            stats.LosingDeals++
            if deal.Profit < stats.LargestLoss {
                stats.LargestLoss = deal.Profit
            }
        }
    }

    stats.AvgProfit = stats.TotalProfit / float64(len(deals))
    stats.WinRate = float64(stats.WinningDeals) / float64(len(deals)) * 100

    return stats, nil
}

func (s *DailyDealStats) Print() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      DAILY DEAL STATISTICS            ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    if s.TotalDeals == 0 {
        fmt.Println("No deals today")
        return
    }

    fmt.Printf("Total Deals:     %d\n", s.TotalDeals)
    fmt.Printf("Winning:         %d (%.1f%%)\n",
        s.WinningDeals, s.WinRate)
    fmt.Printf("Losing:          %d (%.1f%%)\n",
        s.LosingDeals, 100-s.WinRate)
    fmt.Printf("\n")
    fmt.Printf("Total Profit:    $%.2f\n", s.TotalProfit)
    fmt.Printf("Average P/L:     $%.2f\n", s.AvgProfit)
    fmt.Printf("Largest Win:     $%.2f\n", s.LargestWin)
    fmt.Printf("Largest Loss:    $%.2f\n", s.LargestLoss)
}

// Usage:
stats, _ := GetDailyDealStats(sugar)
stats.Print()
```

---

## 🔗 Related Methods

**🍬 Other time periods:**

* `GetDealsYesterday()` - Yesterday's deals
* `GetDealsThisWeek()` - This week's deals
* `GetDealsThisMonth()` - This month's deals
* `GetDealsDateRange()` - Custom date range

**🍬 Profit calculations:**

* `GetProfitToday()` - Just today's profit total (sum of `Profit`)
* `GetDailyStats()` - Win/loss statistics built from this method

---

## ⚠️ Common Pitfalls

### 1) Assuming deals exist

```go
// ❌ WRONG - might panic if no deals
deals, _ := sugar.GetDealsToday()
firstDeal := deals[0] // Panic if empty!

// ✅ CORRECT - check length first
deals, _ := sugar.GetDealsToday()
if len(deals) > 0 {
    firstDeal := deals[0]
}
```

### 2) Not considering timezone

```go
// ❌ WRONG - assuming "today" is the broker's trading day
// The window is computed from YOUR machine's clock (time.Now())
deals, _ := sugar.GetDealsToday()

// ✅ CORRECT - pass explicit boundaries when the day must match the broker's day
// (the server reads the request times as broker server time)
from := time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC)  // broker 00:00
to := time.Date(2026, 10, 16, 23, 59, 59, 0, time.UTC) // broker 23:59:59
deals, _ := sugar.GetDealsDateRange(from, to)
```

### 2b) Expecting positions opened before today

```go
// ❌ WRONG - a position opened yesterday and closed today is NOT returned:
// the range filters on the position's OPEN time
deals, _ := sugar.GetDealsToday()

// ✅ CORRECT - widen the range and filter on CloseTime yourself
deals, _ := sugar.GetDealsThisWeek()
```

### 3) Treating as real-time feed

```go
// ❌ WRONG - this is closed positions only!
deals, _ := sugar.GetDealsToday()
// This doesn't include currently OPEN positions

// ✅ CORRECT - for open positions use different method
openPositions, _ := sugar.GetOpenPositions() // Current positions
closedDeals, _ := sugar.GetDealsToday()      // Closed today
```

---

## 💎 Pro Tips

1. **Empty slice OK** - Returns an empty slice if no deals (not an error)

2. **Closed only** - Only includes closed positions, not currently open

3. **Local clock** - "Today" is 00:00 to now on your machine's clock, matched against the position's open time

4. **Performance** - Uses 5-second timeout

5. **Auto-range** - Automatically calculates today's range (00:00 to now)

---

## 📊 Deal Structure

Each deal (`*pb.PositionHistoryInfo`) is one closed position:
```
Index          - Index in the result
PositionTicket - Position ticket number
OrderType      - pb.AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE (e.g. ..._AH_ORDER_TYPE_BUY / ..._AH_ORDER_TYPE_SELL)
Symbol         - Trading symbol (e.g., "EURUSD")
Volume         - Trade volume in lots
OpenPrice      - Entry price
ClosePrice     - Exit price
StopLoss       - Stop Loss
TakeProfit     - Take Profit
Profit         - Realized profit/loss
Commission     - Trading commission
Fee            - Fees
Swap           - Swap charges
MarketValue    - Market value
OpenTime       - Open time (*timestamppb.Timestamp, broker server time); OpenTimeUtc = same moment in UTC
CloseTime      - Close time (*timestamppb.Timestamp, broker server time); CloseTimeUtc = same moment in UTC
Comment        - Comment
Magic          - Magic number
```

---

**See also:** [`GetDealsYesterday.md`](GetDealsYesterday.md), [`GetDealsThisWeek.md`](GetDealsThisWeek.md), [`GetProfitToday.md`](GetProfitToday.md), [`GetDealsDateRange.md`](GetDealsDateRange.md)
