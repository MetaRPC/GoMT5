# 👤 Get Account Info (`GetAccountInfo`)

> **Sugar method:** Gets ALL account information in one call - balance, equity, margin, leverage, and more!

**API Information:**

* **Method:** `sugar.GetAccountInfo()`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `service.GetAccountSummary()`, `GetMargin()`, `GetFreeMargin()`, `GetMarginLevel()`, `GetProfit()`
* **Timeout:** 5 seconds for `GetAccountSummary`; the four extra queries (`GetMargin`, `GetFreeMargin`, `GetMarginLevel`, `GetProfit`) use their own 3-second timeouts
* **Returns:** Complete AccountInfo structure

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetAccountInfo() (*AccountInfo, error)
```

---

## 🔽 Input

*No parameters required*

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `accountInfo` | `*AccountInfo` | Structure with all account data |
| `error` | `error` | `GetAccountInfo failed: ...` if `GetAccountSummary` fails |

**Note:** errors of the four extra queries are ignored - a failed query leaves `Margin`, `FreeMargin`, `MarginLevel` or `Profit` at `0` without returning an error.

### AccountInfo Structure

```go
type AccountInfo struct {
    Login       int64     // Account login number
    Balance     float64   // Account balance
    Equity      float64   // Current equity (balance + floating P/L)
    Margin      float64   // Used margin
    FreeMargin  float64   // Free margin available
    MarginLevel float64   // Margin level percentage
    Profit      float64   // Total floating profit/loss
    Currency    string    // Account currency (e.g., "USD", "EUR")
    Leverage    int64     // Account leverage (e.g., 100 for 1:100)
    Company     string    // Broker company name
}
```

---

## 💬 Just the Essentials

* **What it is:** Gets complete account snapshot in one method call.
* **Why you need it:** Perfect for dashboards, reports, and monitoring.
* **Sanity check:** One method call for you; internally it makes 5 server requests (account summary + margin, free margin, margin level, profit).

---

## 🎯 When to Use

✅ **Account dashboards** - Display complete account status

✅ **Risk monitoring** - Check margin level and exposure

✅ **Daily reports** - Generate account summary

✅ **Trading bots** - Get account state before trading

---

## 🔗 Usage Examples

### 1) Basic usage - display account info

```go
info, err := sugar.GetAccountInfo()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("  ACCOUNT INFORMATION\n")
fmt.Printf("═══════════════════════════════════════\n")
fmt.Printf("Login:        %d\n", info.Login)
fmt.Printf("Company:      %s\n", info.Company)
fmt.Printf("Currency:     %s\n", info.Currency)
fmt.Printf("Leverage:     1:%d\n\n", info.Leverage)

fmt.Printf("Balance:      $%.2f\n", info.Balance)
fmt.Printf("Equity:       $%.2f\n", info.Equity)
fmt.Printf("Profit:       $%.2f\n\n", info.Profit)

fmt.Printf("Margin:       $%.2f\n", info.Margin)
fmt.Printf("Free Margin:  $%.2f\n", info.FreeMargin)
fmt.Printf("Margin Level: %.2f%%\n", info.MarginLevel)
fmt.Printf("═══════════════════════════════════════\n")

// Output example:
// ═══════════════════════════════════════
//   ACCOUNT INFORMATION
// ═══════════════════════════════════════
// Login:        12345678
// Company:      MetaQuotes Software Corp.
// Currency:     USD
// Leverage:     1:100
//
// Balance:      $10,000.00
// Equity:       $10,250.00
// Profit:       $250.00
//
// Margin:       $1,100.00
// Free Margin:  $9,150.00
// Margin Level: 931.82%
// ═══════════════════════════════════════
```

---

### 2) Check account health

```go
info, err := sugar.GetAccountInfo()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Println("Account Health Check:")
fmt.Println("─────────────────────────────────────────")

// Check 1: Margin level
if info.MarginLevel < 100 {
    fmt.Println("🔴 CRITICAL: Margin level below 100%!")
    fmt.Printf("   Current: %.2f%%\n", info.MarginLevel)
} else if info.MarginLevel < 200 {
    fmt.Println("🟠 WARNING: Margin level below 200%")
    fmt.Printf("   Current: %.2f%%\n", info.MarginLevel)
} else {
    fmt.Println("✅ Margin level healthy")
    fmt.Printf("   Current: %.2f%%\n", info.MarginLevel)
}

// Check 2: Drawdown
drawdown := ((info.Equity - info.Balance) / info.Balance) * 100
if drawdown < -10 {
    fmt.Printf("\n🔴 High drawdown: %.2f%%\n", drawdown)
} else if drawdown < -5 {
    fmt.Printf("\n🟠 Moderate drawdown: %.2f%%\n", drawdown)
} else {
    fmt.Printf("\n✅ Drawdown acceptable: %.2f%%\n", drawdown)
}

// Check 3: Free margin
marginUsage := (info.Margin / info.Equity) * 100
fmt.Printf("\nMargin usage: %.1f%%\n", marginUsage)

if marginUsage > 80 {
    fmt.Println("🔴 Very high margin usage!")
} else if marginUsage > 50 {
    fmt.Println("🟠 High margin usage")
} else {
    fmt.Println("✅ Safe margin usage")
}
```

---

### 3) Real-time account monitor

```go
func MonitorAccount(sugar *mt5.MT5Sugar, intervalSeconds int) {
    ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
    defer ticker.Stop()

    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          ACCOUNT MONITOR (Press Ctrl+C to stop)       ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")

    for range ticker.C {
        info, err := sugar.GetAccountInfo()
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        timestamp := time.Now().Format("15:04:05")

        fmt.Printf("\n[%s] Balance: $%.2f | Equity: $%.2f | P/L: $%.2f | Margin: %.1f%%\n",
            timestamp,
            info.Balance,
            info.Equity,
            info.Profit,
            info.MarginLevel,
        )

        // Alert if margin level drops
        if info.MarginLevel < 150 {
            fmt.Printf("⚠️  WARNING: Low margin level! (%.2f%%)\n", info.MarginLevel)
        }

        // Alert if large loss
        if info.Profit < -info.Balance*0.05 { // -5% of balance
            fmt.Printf("🔴 ALERT: Large floating loss! ($%.2f)\n", info.Profit)
        }
    }
}

// Usage:
MonitorAccount(sugar, 5) // Update every 5 seconds
```

---

### 4) Daily account report

```go
func GenerateDailyReport(sugar *mt5.MT5Sugar) {
    info, err := sugar.GetAccountInfo()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Get today's deals
    dealsToday, _ := sugar.GetDealsToday()
    profitToday, _ := sugar.GetProfitToday()
    dailyStats, _ := sugar.GetDailyStats()

    // Get positions
    openPositions, _ := sugar.GetOpenPositions()

    date := time.Now().Format("2006-01-02")

    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Printf("║          DAILY ACCOUNT REPORT - %s            ║\n", date)
    fmt.Println("╚═══════════════════════════════════════════════════════╝")

    fmt.Println("\n📊 ACCOUNT STATUS")
    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Account:      #%d (%s)\n", info.Login, info.Company)
    fmt.Printf("Currency:     %s\n", info.Currency)
    fmt.Printf("Leverage:     1:%d\n\n", info.Leverage)

    fmt.Printf("Balance:      $%10.2f\n", info.Balance)
    fmt.Printf("Equity:       $%10.2f\n", info.Equity)
    fmt.Printf("Margin:       $%10.2f\n", info.Margin)
    fmt.Printf("Free:         $%10.2f\n", info.FreeMargin)
    fmt.Printf("Margin Level: %9.2f%%\n", info.MarginLevel)

    fmt.Println("\n📈 TODAY'S TRADING")
    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Closed deals:    %d\n", len(dealsToday))
    if dailyStats != nil {
        fmt.Printf("Win rate:        %.1f%%\n", dailyStats.WinRate)
        fmt.Printf("Best deal:       $%.2f\n", dailyStats.BestDeal)
        fmt.Printf("Worst deal:      $%.2f\n", dailyStats.WorstDeal)
    }
    fmt.Printf("Total profit:    $%.2f\n", profitToday)

    fmt.Println("\n💼 OPEN POSITIONS")
    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Count:           %d\n", len(openPositions))
    fmt.Printf("Floating P/L:    $%.2f\n", info.Profit)

    // Calculate daily performance
    dailyReturn := (profitToday / info.Balance) * 100

    fmt.Println("\n📊 PERFORMANCE")
    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Daily return:    %.2f%%\n", dailyReturn)

    if info.Profit >= 0 {
        fmt.Printf("Total exposure:  +$%.2f\n", info.Profit+profitToday)
    } else {
        fmt.Printf("Total exposure:  -$%.2f\n", math.Abs(info.Profit+profitToday))
    }

    fmt.Println("═══════════════════════════════════════════════════════")
}

// Usage:
GenerateDailyReport(sugar)
```

---

### 5) Pre-trade risk check

```go
func CheckAccountBeforeTrade(sugar *mt5.MT5Sugar, requiredMargin float64) (bool, string) {
    info, err := sugar.GetAccountInfo()
    if err != nil {
        return false, fmt.Sprintf("error getting account info: %v", err)
    }

    fmt.Println("Pre-Trade Risk Check:")
    fmt.Println("─────────────────────────────────────────")

    // Check 1: Margin level
    if info.MarginLevel < 300 {
        return false, fmt.Sprintf(
            "margin level too low (%.2f%%, need >300%%)",
            info.MarginLevel,
        )
    }
    fmt.Printf("✅ Margin level: %.2f%%\n", info.MarginLevel)

    // Check 2: Free margin
    if requiredMargin > info.FreeMargin {
        return false, fmt.Sprintf(
            "insufficient free margin (need $%.2f, have $%.2f)",
            requiredMargin, info.FreeMargin,
        )
    }
    fmt.Printf("✅ Free margin: $%.2f (need $%.2f)\n",
        info.FreeMargin, requiredMargin)

    // Check 3: Current drawdown
    drawdown := ((info.Equity - info.Balance) / info.Balance) * 100
    if drawdown < -10 {
        return false, fmt.Sprintf(
            "high drawdown (%.2f%%, max -10%%)",
            drawdown,
        )
    }
    fmt.Printf("✅ Drawdown: %.2f%%\n", drawdown)

    fmt.Println("─────────────────────────────────────────")
    fmt.Println("✅ All checks passed - safe to trade")

    return true, ""
}

// Usage before trading:
requiredMargin := 500.0

canTrade, reason := CheckAccountBeforeTrade(sugar, requiredMargin)
if !canTrade {
    fmt.Printf("❌ Cannot trade: %s\n", reason)
    return
}

// Proceed with trade...
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000) // SL/TP in points
fmt.Printf("Trade opened: #%d\n", ticket)
```

---

### 6) Compare account snapshots

```go
type AccountSnapshot struct {
    Timestamp time.Time
    Info      *mt5.AccountInfo
}

func TakeSnapshot(sugar *mt5.MT5Sugar) (*AccountSnapshot, error) {
    info, err := sugar.GetAccountInfo()
    if err != nil {
        return nil, err
    }

    return &AccountSnapshot{
        Timestamp: time.Now(),
        Info:      info,
    }, nil
}

func CompareSnapshots(before, after *AccountSnapshot) {
    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          ACCOUNT COMPARISON                           ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")

    duration := after.Timestamp.Sub(before.Timestamp)

    fmt.Printf("Time period: %s\n\n", duration.Round(time.Second))

    // Balance change
    balanceChange := after.Info.Balance - before.Info.Balance
    balanceChangePercent := (balanceChange / before.Info.Balance) * 100

    fmt.Println("BALANCE:")
    fmt.Printf("  Before: $%.2f\n", before.Info.Balance)
    fmt.Printf("  After:  $%.2f\n", after.Info.Balance)
    if balanceChange >= 0 {
        fmt.Printf("  Change: +$%.2f (+%.2f%%)\n", balanceChange, balanceChangePercent)
    } else {
        fmt.Printf("  Change: -$%.2f (%.2f%%)\n", math.Abs(balanceChange), balanceChangePercent)
    }

    // Equity change
    equityChange := after.Info.Equity - before.Info.Equity
    equityChangePercent := (equityChange / before.Info.Equity) * 100

    fmt.Println("\nEQUITY:")
    fmt.Printf("  Before: $%.2f\n", before.Info.Equity)
    fmt.Printf("  After:  $%.2f\n", after.Info.Equity)
    if equityChange >= 0 {
        fmt.Printf("  Change: +$%.2f (+%.2f%%)\n", equityChange, equityChangePercent)
    } else {
        fmt.Printf("  Change: -$%.2f (%.2f%%)\n", math.Abs(equityChange), equityChangePercent)
    }

    // Margin level
    fmt.Println("\nMARGIN LEVEL:")
    fmt.Printf("  Before: %.2f%%\n", before.Info.MarginLevel)
    fmt.Printf("  After:  %.2f%%\n", after.Info.MarginLevel)

    fmt.Println("═══════════════════════════════════════════════════════")
}

// Usage:
// Take snapshot before trading
snapshot1, _ := TakeSnapshot(sugar)

// ... do some trading ...

time.Sleep(1 * time.Hour)

// Take snapshot after trading
snapshot2, _ := TakeSnapshot(sugar)

// Compare
CompareSnapshots(snapshot1, snapshot2)
```

---

### 7) Account status widget

```go
func ShowAccountWidget(sugar *mt5.MT5Sugar) {
    info, err := sugar.GetAccountInfo()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Determine account health status
    var healthStatus string
    var healthColor string

    if info.MarginLevel < 100 {
        healthStatus = "CRITICAL"
        healthColor = "🔴"
    } else if info.MarginLevel < 200 {
        healthStatus = "WARNING"
        healthColor = "🟠"
    } else if info.MarginLevel < 500 {
        healthStatus = "MODERATE"
        healthColor = "🟡"
    } else {
        healthStatus = "HEALTHY"
        healthColor = "🟢"
    }

    // Calculate metrics
    marginUsage := (info.Margin / info.Equity) * 100
    drawdown := ((info.Equity - info.Balance) / info.Balance) * 100

    fmt.Println("┌─────────────────────────────────────────┐")
    fmt.Printf("│ Account #%d%-23s│\n", info.Login, "")
    fmt.Println("├─────────────────────────────────────────┤")
    fmt.Printf("│ Balance:   $%-28.2f│\n", info.Balance)
    fmt.Printf("│ Equity:    $%-28.2f│\n", info.Equity)
    fmt.Printf("│ P/L:       $%-28.2f│\n", info.Profit)
    fmt.Println("├─────────────────────────────────────────┤")
    fmt.Printf("│ Margin:    %-29.1f%%│\n", marginUsage)
    fmt.Printf("│ Drawdown:  %-29.2f%%│\n", drawdown)
    fmt.Println("├─────────────────────────────────────────┤")
    fmt.Printf("│ Status:    %s %-26s│\n", healthColor, healthStatus)
    fmt.Println("└─────────────────────────────────────────┘")
}

// Usage:
ShowAccountWidget(sugar)

// Output:
// ┌───────────────────────────────────────
// │ Account #12345678                      
// ├───────────────────────────────────────
// │ Balance:   $10,000.00                   
// │ Equity:    $10,250.00                   
// │ P/L:       $250.00                      
// ├───────────────────────────────────────
// │ Margin:    10.7%                       
// │ Drawdown:  2.50%                        
// ├───────────────────────────────────────
// │ Status:    🟢 HEALTHY                  
// └───────────────────────────────────────
```

---

### 8) Export account data to JSON

```go
import "encoding/json"

func ExportAccountToJSON(sugar *mt5.MT5Sugar, filename string) error {
    info, err := sugar.GetAccountInfo()
    if err != nil {
        return err
    }

    // Create export structure with timestamp
    export := struct {
        Timestamp   string            `json:"timestamp"`
        AccountInfo *mt5.AccountInfo  `json:"account_info"`
    }{
        Timestamp:   time.Now().Format(time.RFC3339),
        AccountInfo: info,
    }

    // Marshal to JSON
    data, err := json.MarshalIndent(export, "", "  ")
    if err != nil {
        return err
    }

    // Write to file
    err = os.WriteFile(filename, data, 0644)
    if err != nil {
        return err
    }

    fmt.Printf("✅ Account data exported to %s\n", filename)

    return nil
}

// Usage:
err := ExportAccountToJSON(sugar, "account_snapshot.json")
if err != nil {
    fmt.Printf("Export failed: %v\n", err)
}

// Result: account_snapshot.json
// {
//   "timestamp": "2024-01-15T14:30:00Z",
//   "account_info": {
//     "Login": 12345678,
//     "Balance": 10000.00,
//     "Equity": 10250.00,
//     "Margin": 1100.00,
//     "FreeMargin": 9150.00,
//     "MarginLevel": 931.82,
//     "Profit": 250.00,
//     "Currency": "USD",
//     "Leverage": 100,
//     "Company": "MetaQuotes Software Corp."
//   }
// }
```

---

### 9) Multi-account manager

```go
type AccountManager struct {
    accounts map[string]*mt5.MT5Sugar
}

func NewAccountManager() *AccountManager {
    return &AccountManager{
        accounts: make(map[string]*mt5.MT5Sugar),
    }
}

func (am *AccountManager) AddAccount(name string, sugar *mt5.MT5Sugar) {
    am.accounts[name] = sugar
}

func (am *AccountManager) ShowAllAccounts() {
    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          MULTI-ACCOUNT OVERVIEW                       ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")

    totalBalance := 0.0
    totalEquity := 0.0
    totalProfit := 0.0

    for name, sugar := range am.accounts {
        info, err := sugar.GetAccountInfo()
        if err != nil {
            fmt.Printf("\n%s: ❌ Error - %v\n", name, err)
            continue
        }

        totalBalance += info.Balance
        totalEquity += info.Equity
        totalProfit += info.Profit

        fmt.Printf("\n%s:\n", name)
        fmt.Printf("  Balance:  $%10.2f\n", info.Balance)
        fmt.Printf("  Equity:   $%10.2f\n", info.Equity)
        fmt.Printf("  P/L:      $%10.2f\n", info.Profit)
        fmt.Printf("  Margin:   %9.2f%%\n", info.MarginLevel)
    }

    fmt.Println("\n─────────────────────────────────────────────────────────")
    fmt.Printf("TOTAL ACROSS ALL ACCOUNTS:\n")
    fmt.Printf("  Balance:  $%10.2f\n", totalBalance)
    fmt.Printf("  Equity:   $%10.2f\n", totalEquity)
    fmt.Printf("  P/L:      $%10.2f\n", totalProfit)
    fmt.Println("═══════════════════════════════════════════════════════")
}

// Usage:
manager := NewAccountManager()
manager.AddAccount("Main", sugar1)
manager.AddAccount("Backup", sugar2)
manager.AddAccount("Testing", sugar3)

manager.ShowAllAccounts()
```

---

### 10) Advanced account dashboard

```go
type AccountDashboard struct {
    sugar *mt5.MT5Sugar
}

func NewAccountDashboard(sugar *mt5.MT5Sugar) *AccountDashboard {
    return &AccountDashboard{sugar: sugar}
}

func (ad *AccountDashboard) Render() error {
    // Get account info
    info, err := ad.sugar.GetAccountInfo()
    if err != nil {
        return err
    }

    // Get additional data
    positions, _ := ad.sugar.GetOpenPositions()
    dealsToday, _ := ad.sugar.GetDealsToday()
    dailyStats, _ := ad.sugar.GetDailyStats()

    // Calculate metrics
    marginUsage := (info.Margin / info.Equity) * 100
    drawdown := ((info.Equity - info.Balance) / info.Balance) * 100
    dailyReturn := 0.0
    if dailyStats != nil && len(dealsToday) > 0 {
        dailyProfit, _ := ad.sugar.GetProfitToday()
        dailyReturn = (dailyProfit / info.Balance) * 100
    }

    // Render dashboard
    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          TRADING DASHBOARD                            ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")

    // Account section
    fmt.Println("\n📊 ACCOUNT OVERVIEW")
    fmt.Println("─────────────────────────────────────────────────────────")
    fmt.Printf("Account:     #%d\n", info.Login)
    fmt.Printf("Broker:      %s\n", info.Company)
    fmt.Printf("Currency:    %s\n", info.Currency)
    fmt.Printf("Leverage:    1:%d\n", info.Leverage)

    // Balance section
    fmt.Println("\n💰 BALANCE & EQUITY")
    fmt.Println("─────────────────────────────────────────────────────────")
    fmt.Printf("Balance:     $%10.2f\n", info.Balance)
    fmt.Printf("Equity:      $%10.2f\n", info.Equity)
    fmt.Printf("Floating P/L: $%9.2f", info.Profit)
    if info.Profit >= 0 {
        fmt.Printf(" 🟢\n")
    } else {
        fmt.Printf(" 🔴\n")
    }
    fmt.Printf("Drawdown:    %10.2f%%\n", drawdown)

    // Margin section
    fmt.Println("\n📈 MARGIN USAGE")
    fmt.Println("─────────────────────────────────────────────────────────")
    fmt.Printf("Used:        $%10.2f\n", info.Margin)
    fmt.Printf("Free:        $%10.2f\n", info.FreeMargin)
    fmt.Printf("Level:       %10.2f%%", info.MarginLevel)

    if info.MarginLevel < 100 {
        fmt.Printf(" 🔴\n")
    } else if info.MarginLevel < 200 {
        fmt.Printf(" 🟠\n")
    } else {
        fmt.Printf(" 🟢\n")
    }

    fmt.Printf("Usage:       %10.1f%%\n", marginUsage)

    // Positions section
    fmt.Println("\n💼 OPEN POSITIONS")
    fmt.Println("─────────────────────────────────────────────────────────")
    fmt.Printf("Count:       %d\n", len(positions))

    if len(positions) > 0 {
        buyCount := 0
        sellCount := 0
        for _, pos := range positions {
            if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
                buyCount++
            } else {
                sellCount++
            }
        }
        fmt.Printf("BUY:         %d\n", buyCount)
        fmt.Printf("SELL:        %d\n", sellCount)
    }

    // Today's trading
    if dailyStats != nil {
        fmt.Println("\n📊 TODAY'S PERFORMANCE")
        fmt.Println("─────────────────────────────────────────────────────────")
        fmt.Printf("Trades:      %d\n", dailyStats.TotalDeals)
        fmt.Printf("Win rate:    %.1f%%\n", dailyStats.WinRate)
        fmt.Printf("Profit:      $%.2f", dailyStats.TotalProfit)
        if dailyStats.TotalProfit >= 0 {
            fmt.Printf(" 🟢\n")
        } else {
            fmt.Printf(" 🔴\n")
        }
        fmt.Printf("Return:      %.2f%%\n", dailyReturn)
    }

    fmt.Println("═══════════════════════════════════════════════════════")

    return nil
}

// Usage:
dashboard := NewAccountDashboard(sugar)

// Update dashboard every 10 seconds
ticker := time.NewTicker(10 * time.Second)
defer ticker.Stop()

for range ticker.C {
    // Clear screen (optional)
    fmt.Print("\033[H\033[2J")

    // Render dashboard
    if err := dashboard.Render(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `service.GetAccountSummary()` - Get core account data
* `GetMargin()` - Get used margin
* `GetFreeMargin()` - Get free margin
* `GetMarginLevel()` - Get margin level
* `GetProfit()` - Get floating P/L

**🍬 Complementary sugar methods:**

* `GetDailyStats()` - Get today's trading stats ⭐
* `GetDealsToday()` - Get today's closed positions
* `GetOpenPositions()` - Get currently open positions
* `GetBalance()` - Get balance only
* `GetEquity()` - Get equity only

**Use cases:**

```go
// Complete account snapshot
info, _ := sugar.GetAccountInfo()

// Just need balance
balance, _ := sugar.GetBalance()

// Full dashboard
info, _ := sugar.GetAccountInfo()
stats, _ := sugar.GetDailyStats()
positions, _ := sugar.GetOpenPositions()
```

---

## ⚠️ Common Pitfalls

### 1) Confusing balance and equity

```go
info, _ := sugar.GetAccountInfo()

// ❌ WRONG - using balance to check if can trade
if info.Balance > 1000 {
    // Balance doesn't include floating P/L!
}

// ✅ CORRECT - use equity (includes floating P/L)
if info.Equity > 1000 {
    // Equity = Balance + Floating P/L
}
```

### 2) Not checking margin level

```go
// ❌ WRONG - ignoring margin level
info, _ := sugar.GetAccountInfo()
sugar.BuyMarket("EURUSD", 10.0) // Might cause margin call!

// ✅ CORRECT - check margin level first
info, _ := sugar.GetAccountInfo()
if info.MarginLevel < 300 {
    fmt.Println("Margin level too low!")
    return
}
```

### 3) Using stale data

```go
// ❌ WRONG - using old account info
info, _ := sugar.GetAccountInfo()
time.Sleep(10 * time.Minute)
// info is now stale!
if info.FreeMargin > 1000 {
    sugar.BuyMarket("EURUSD", 1.0)
}

// ✅ CORRECT - get fresh data before trading
info, _ := sugar.GetAccountInfo()
if info.FreeMargin > 1000 {
    sugar.BuyMarket("EURUSD", 1.0)
}
```

### 4) Not handling errors

```go
// ❌ WRONG - ignoring errors
info, _ := sugar.GetAccountInfo()
fmt.Printf("Balance: $%.2f\n", info.Balance) // Might panic!

// ✅ CORRECT - check errors
info, err := sugar.GetAccountInfo()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
fmt.Printf("Balance: $%.2f\n", info.Balance)
```

### 5) Misunderstanding margin level

```go
// ❌ WRONG - thinking higher margin is good
info, _ := sugar.GetAccountInfo()
if info.Margin > 5000 {
    fmt.Println("High margin is good!") // NO!
}

// ✅ CORRECT - higher margin LEVEL % is good
if info.MarginLevel > 500 {
    fmt.Println("Healthy margin level")
}

// Margin = how much you're using (lower is better)
// Margin Level = (Equity/Margin)*100 (higher is better)
```

---

## 💎 Pro Tips

1. **One call, all data** - Convenient snapshot (internally 5 requests, not one)

2. **Check margin level** - Below 100% = margin call risk

3. **Monitor equity, not balance** - Equity includes floating P/L

4. **Refresh frequently** - Account info changes constantly

5. **Use for dashboards** - Perfect for real-time monitoring

6. **Margin Level formula** - (Equity / Margin) × 100

7. **Margin usage** - (Margin / Equity) × 100 = % of equity used

---

## 📊 Account Metrics Explained

```
Balance:
- Your account balance (realized P/L only)
- Doesn't change until you close positions

Equity:
- Balance + Floating P/L
- Changes in real-time with price movements
- Equity = Balance + Profit

Margin:
- How much margin you're using
- Lower is better (less risk)

Free Margin:
- Equity - Margin
- How much you can still use for new positions

Margin Level:
- (Equity / Margin) × 100
- Higher is better
- Below 100% = margin call
- Above 200% = safe
- Above 500% = very safe

Margin Usage:
- (Margin / Equity) × 100
- Lower is better
- Above 80% = very risky
- Below 30% = conservative
```

---

**See also:** [`GetDailyStats.md`](../8.%20History_Statistics/GetDailyStats.md)
