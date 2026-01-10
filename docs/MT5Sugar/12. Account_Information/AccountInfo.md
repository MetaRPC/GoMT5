# MT5Sugar · AccountInfo Structure

> **Complete account information structure** returned by GetAccountInfo(). Contains all essential account metrics in one convenient package.

## 📋 Structure Definition

```go
type AccountInfo struct {
    Login       int64     // Account login number
    Balance     float64   // Account balance
    Equity      float64   // Current equity (balance + floating P/L)
    Margin      float64   // Used margin
    FreeMargin  float64   // Free margin available
    MarginLevel float64   // Margin level % ((Equity/Margin)*100)
    Profit      float64   // Total floating profit/loss
    Currency    string    // Account currency (USD, EUR, etc.)
    Leverage    int64     // Account leverage (100 = 1:100)
    Company     string    // Broker company name
}
```

## 📊 Field Details

### Login (int64)

- **What:** MT5 account number
- **Example:** `591129415`
- **Use:** Account identification, logging, reports

### Balance (float64)

- **What:** Realized balance (closed positions only)
- **Formula:** Initial deposit + realized profits - realized losses
- **Example:** `10000.00`
- **Important:** Does NOT include floating P/L from open positions

### Equity (float64)

- **What:** Real-time account value
- **Formula:** Balance + Profit (floating P/L)
- **Example:** `10250.75`
- **Important:** Changes with every price tick if you have open positions

### Margin (float64)

- **What:** Total margin used by open positions
- **Formula:** Sum of all position margins
- **Example:** `500.00`
- **Important:** Higher margin = less available for new trades

### FreeMargin (float64)

- **What:** Margin available for new positions
- **Formula:** Equity - Margin
- **Example:** `9750.75`
- **Important:** New positions reduce this value

### MarginLevel (float64)

- **What:** Margin safety indicator (percentage)
- **Formula:** (Equity / Margin) × 100
- **Example:** `2050.15` (means 2050.15%)
- **Critical levels:**
  - Below 100% = Margin call risk
  - Below 50% = Stop out (broker closes positions)
  - Above 1000% = Very safe

### Profit (float64)

- **What:** Total floating profit/loss from ALL open positions
- **Formula:** Sum of all open position P/L
- **Example:** `250.75` (positive = profit, negative = loss)
- **Important:** Changes in real-time with price movements

### Currency (string)

- **What:** Account base currency
- **Example:** `"USD"`, `"EUR"`, `"GBP"`
- **Use:** Profit calculations, reporting

### Leverage (int64)

- **What:** Account leverage multiplier
- **Example:** `100` (means 1:100 leverage)
- **Common values:** 50, 100, 200, 500
- **Impact:** Higher leverage = less margin needed

### Company (string)

- **What:** Broker/company name
- **Example:** `"FxPro Financial Services Ltd"`
- **Use:** Identification, multi-account management

## 🎯 Common Calculations

### Drawdown (%)

```go
drawdown := ((info.Equity - info.Balance) / info.Balance) * 100
// Negative = losing, Positive = winning
```

### Margin Usage (%)
```go
marginUsage := (info.Margin / info.Equity) * 100
// Higher = more risk
```

### Available Capacity
```go
capacity := info.FreeMargin / info.Margin
// How many times current position size you can still open
```

### Daily Return (%)
```go
// Requires yesterday's balance
dailyReturn := ((info.Balance - yesterdayBalance) / yesterdayBalance) * 100
```

## 🟢 Usage Examples

### Example 1: Basic account display

```go
info, _ := sugar.GetAccountInfo()

fmt.Printf("Account: #%d (%s)\n", info.Login, info.Company)
fmt.Printf("Currency: %s, Leverage: 1:%d\n\n", info.Currency, info.Leverage)

fmt.Printf("Balance:  $%.2f\n", info.Balance)
fmt.Printf("Equity:   $%.2f\n", info.Equity)
fmt.Printf("Profit:   $%.2f\n\n", info.Profit)

fmt.Printf("Margin:   $%.2f\n", info.Margin)
fmt.Printf("Free:     $%.2f\n", info.FreeMargin)
fmt.Printf("Level:    %.2f%%\n", info.MarginLevel)
```

### Example 2: Risk assessment

```go
info, _ := sugar.GetAccountInfo()

// Check margin level
if info.MarginLevel < 100 {
    fmt.Println("🔴 CRITICAL: Margin call risk!")
} else if info.MarginLevel < 200 {
    fmt.Println("🟠 WARNING: Low margin level")
} else if info.MarginLevel < 500 {
    fmt.Println("🟡 CAUTION: Moderate margin level")
} else {
    fmt.Println("✅ Healthy margin level")
}

// Check margin usage
marginUsage := (info.Margin / info.Equity) * 100
fmt.Printf("Margin usage: %.1f%%\n", marginUsage)

if marginUsage > 50 {
    fmt.Println("⚠️ High margin usage - reduce exposure")
}
```

### Example 3: Drawdown monitoring

```go
info, _ := sugar.GetAccountInfo()

drawdown := ((info.Equity - info.Balance) / info.Balance) * 100

fmt.Printf("Current drawdown: %.2f%%\n", drawdown)

if drawdown < -10 {
    fmt.Println("🔴 ALERT: Drawdown exceeds 10%!")
} else if drawdown < -5 {
    fmt.Println("⚠️ WARNING: Significant drawdown")
} else if drawdown > 0 {
    fmt.Printf("✅ In profit: +%.2f%%\n", drawdown)
}
```

### Example 4: Position sizing calculator

```go
info, _ := sugar.GetAccountInfo()

// Use 1% risk per trade
riskPercent := 1.0
riskAmount := info.Balance * (riskPercent / 100)

fmt.Printf("Account balance: $%.2f\n", info.Balance)
fmt.Printf("Risk per trade (1%%): $%.2f\n", riskAmount)

// Check available margin
availableRisk := info.FreeMargin * 0.2 // Use only 20% of free margin
fmt.Printf("Max safe risk: $%.2f\n", availableRisk)

if riskAmount > availableRisk {
    fmt.Println("⚠️ Reduce position size - insufficient margin")
}
```

### Example 5: Multi-account comparison

```go
accounts := []*mt5.AccountInfo{account1, account2, account3}

fmt.Println("ACCOUNT COMPARISON")
fmt.Println("==================")

for i, info := range accounts {
    fmt.Printf("\nAccount %d: #%d\n", i+1, info.Login)
    fmt.Printf("  Balance:  $%.2f\n", info.Balance)
    fmt.Printf("  Equity:   $%.2f\n", info.Equity)
    fmt.Printf("  Profit:   $%.2f\n", info.Profit)
    fmt.Printf("  Margin:   %.2f%%\n", info.MarginLevel)
}
```

### Example 6: JSON export

```go
info, _ := sugar.GetAccountInfo()

// Convert to JSON
jsonData, _ := json.MarshalIndent(info, "", "  ")
fmt.Println(string(jsonData))

// Output:
// {
//   "Login": 591129415,
//   "Balance": 10000.00,
//   "Equity": 10250.75,
//   "Margin": 500.00,
//   ...
// }
```

### Example 7: Database record

```go
type AccountSnapshot struct {
    Timestamp time.Time
    Account   *mt5.AccountInfo
}

func saveSnapshot(info *mt5.AccountInfo) {
    snapshot := AccountSnapshot{
        Timestamp: time.Now(),
        Account:   info,
    }

    // Save to database
    db.Insert("snapshots", snapshot)

    fmt.Printf("Snapshot saved: %s\n", snapshot.Timestamp)
}
```

### Example 8: Alert system

```go
func checkAccountHealth(info *mt5.AccountInfo) []string {
    var alerts []string

    // Margin level check
    if info.MarginLevel < 100 {
        alerts = append(alerts, "CRITICAL: Margin level below 100%")
    }

    // Drawdown check
    drawdown := ((info.Equity - info.Balance) / info.Balance) * 100
    if drawdown < -15 {
        alerts = append(alerts, fmt.Sprintf("CRITICAL: Drawdown %.1f%%", drawdown))
    }

    // Large position check
    marginUsage := (info.Margin / info.Equity) * 100
    if marginUsage > 70 {
        alerts = append(alerts, "WARNING: High margin usage")
    }

    // Floating loss check
    if info.Profit < -1000 {
        alerts = append(alerts, fmt.Sprintf("WARNING: Large floating loss $%.2f", info.Profit))
    }

    return alerts
}

// Usage
info, _ := sugar.GetAccountInfo()
alerts := checkAccountHealth(info)

if len(alerts) > 0 {
    fmt.Println("🚨 ALERTS:")
    for _, alert := range alerts {
        fmt.Printf("  - %s\n", alert)
    }
} else {
    fmt.Println("✅ All systems normal")
}
```

### Example 9: Daily report

```go
func generateDailyReport(info *mt5.AccountInfo, stats *mt5.DailyStats) {
    fmt.Println("╔════════════════════════════════════════╗")
    fmt.Println("║        DAILY TRADING REPORT            ║")
    fmt.Println("╚════════════════════════════════════════╝")

    fmt.Printf("\n📊 ACCOUNT: #%d (%s)\n", info.Login, info.Company)
    fmt.Printf("Currency: %s, Leverage: 1:%d\n\n", info.Currency, info.Leverage)

    fmt.Println("💰 BALANCE SHEET")
    fmt.Printf("Balance:      $%.2f\n", info.Balance)
    fmt.Printf("Equity:       $%.2f\n", info.Equity)
    fmt.Printf("Floating P/L: $%.2f\n", info.Profit)

    fmt.Println("\n📈 MARGIN STATUS")
    fmt.Printf("Used margin:  $%.2f\n", info.Margin)
    fmt.Printf("Free margin:  $%.2f\n", info.FreeMargin)
    fmt.Printf("Margin level: %.2f%%\n", info.MarginLevel)

    fmt.Println("\n📊 TODAY'S STATS")
    fmt.Printf("Trades:       %d\n", stats.TotalDeals)
    fmt.Printf("Win rate:     %.1f%%\n", stats.WinRate)
    fmt.Printf("Total profit: $%.2f\n", stats.TotalProfit)
}
```

### Example 10: Risk calculator

```go
func calculateMaxPosition(info *mt5.AccountInfo, symbol string, riskPercent float64) {
    // Risk amount
    riskAmount := info.Balance * (riskPercent / 100)

    // Available margin check
    usableMargin := info.FreeMargin * 0.5 // Use only 50%

    fmt.Printf("Account Balance: $%.2f\n", info.Balance)
    fmt.Printf("Risk Amount (%%.1f%%): $%.2f\n", riskPercent, riskAmount)
    fmt.Printf("Usable Margin: $%.2f\n", usableMargin)

    // Calculate position size
    lotSize, _ := sugar.CalculatePositionSize(symbol, riskPercent, 50)

    // Check if affordable
    requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, lotSize)

    if requiredMargin > usableMargin {
        fmt.Println("❌ Position size exceeds available margin")
        lotSize = lotSize * (usableMargin / requiredMargin)
        fmt.Printf("Adjusted lot size: %.2f\n", lotSize)
    } else {
        fmt.Printf("✅ Safe lot size: %.2f\n", lotSize)
    }
}
```

## 🔗 Related

- **[GetAccountInfo](./GetAccountInfo.md)** - method that returns this structure
- **[DailyStats](./DailyStats.md)** - trading statistics structure


## ⚠️ Important Notes

1. **Balance vs Equity** - Balance = realized only, Equity = real-time value

2. **Margin Level** - Below 100% is danger zone!

3. **Floating Profit** - Changes every tick if positions are open

4. **Free Margin** - Required for opening new positions

5. **Currency** - All monetary values in this currency

## 💡 Quick Reference

| Field | Meaning | When to use |
|-------|---------|-------------|
| `Balance` | Realized money | Performance tracking |
| `Equity` | Real-time value | Current account worth |
| `Margin` | Used margin | Risk monitoring |
| `FreeMargin` | Available margin | Position sizing |
| `MarginLevel` | Safety % | Risk alerts |
| `Profit` | Floating P/L | Real-time monitoring |

---

**Summary:** AccountInfo structure provides complete account snapshot. Use Balance for historical tracking, Equity for current value, and MarginLevel for risk monitoring. Essential for risk management and account monitoring.
