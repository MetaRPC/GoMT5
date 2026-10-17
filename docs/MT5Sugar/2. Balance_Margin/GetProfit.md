# 📈 Get Floating Profit/Loss (`GetProfit`)

> **Sugar method:** Returns current floating profit/loss from all open positions.

**API Information:**

* **Method:** `sugar.GetProfit()`
* **Timeout:** 3 seconds
* **Returns:** Floating P/L as `float64`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetProfit() (float64, error)
```

---

## 💬 Just the Essentials

* **What it is:** Total unrealized profit/loss from ALL open positions (not closed deals!).
* **Why you need it:** Monitor performance in real-time, set profit targets, manage risk.
* **Sanity check:** Positive = winning, Negative = losing, Zero = breakeven or no positions.
* **Under the hood:** `service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_PROFIT)`.

---

## 🧮 Formula

```
Floating Profit = Sum of profit from all open positions

Equity = Balance + Floating Profit

Examples:
  3 positions: +$100, +$50, -$30 → Total Profit: +$120
  No open positions → Profit: $0
  All losing: -$50, -$80 → Total Profit: -$130
```

---

## 🔗 Usage Examples

### 1) Basic usage

```go
profit, err := sugar.GetProfit()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Floating P/L: $%.2f\n", profit)

if profit > 0 {
    fmt.Println("✅ Currently winning")
} else if profit < 0 {
    fmt.Println("❌ Currently losing")
} else {
    fmt.Println("➖ Breakeven")
}
```

---

### 2) Close all at profit target

```go
func MonitorProfitTarget(sugar *mt5.MT5Sugar, targetProfit float64) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        profit, _ := sugar.GetProfit()

        fmt.Printf("Current P/L: $%.2f (Target: $%.2f)\n", profit, targetProfit)

        if profit >= targetProfit {
            fmt.Printf("🎯 Profit target reached! Closing all positions...\n")
            sugar.CloseAllPositions()
            fmt.Println("✅ All positions closed")
            return
        }
    }
}

// Run in background
go MonitorProfitTarget(sugar, 500.0) // Close all when profit hits $500
```

---

### 3) Track profit percentage

```go
balance, _ := sugar.GetBalance()
profit, _ := sugar.GetProfit()

profitPercent := (profit / balance) * 100

fmt.Printf("Balance:       $%.2f\n", balance)
fmt.Printf("Floating P/L:  $%.2f (%.2f%%)\n", profit, profitPercent)

if profitPercent > 2 {
    fmt.Println("🎉 Great performance!")
} else if profitPercent < -2 {
    fmt.Println("⚠️  Consider closing losing positions")
}
```

---

### 4) Daily profit monitoring

```go
startBalance, _ := sugar.GetBalance()
startTime := time.Now()

ticker := time.NewTicker(1 * time.Minute)
defer ticker.Stop()

for range ticker.C {
    profit, _ := sugar.GetProfit()
    duration := time.Since(startTime)

    fmt.Printf("[%02d:%02d] Floating P/L: $%.2f\n",
        int(duration.Hours()),
        int(duration.Minutes())%60,
        profit)

    // Close all if daily loss limit reached
    if profit < -200 {
        fmt.Println("🛑 Daily loss limit reached - closing all")
        sugar.CloseAllPositions()
        return
    }
}
```

---

### 5) Real-time P/L dashboard

```go
balance, _ := sugar.GetBalance()
equity, _ := sugar.GetEquity()
profit, _ := sugar.GetProfit()
positions, _ := sugar.CountOpenPositions()

profitPercent := (profit / balance) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║      REAL-TIME PERFORMANCE            ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Open Positions:   %d\n", positions)
fmt.Printf("Balance:          $%.2f\n", balance)
fmt.Printf("Equity:           $%.2f\n", equity)
fmt.Printf("Floating P/L:     $%.2f (%.2f%%)\n", profit, profitPercent)
fmt.Println()

if profit > 0 {
    fmt.Printf("✅ UP by $%.2f from balance\n", profit)
} else if profit < 0 {
    fmt.Printf("❌ DOWN by $%.2f from balance\n", -profit)
} else {
    fmt.Println("➖ At breakeven")
}
```

---

## 🔗 Related Methods

* `GetBalance()` - Account balance
* `GetEquity()` - Balance + Profit
* `GetTotalProfit()` - Sum of `Profit` over the open-position list (GetProfit reads the account's `ACCOUNT_PROFIT` instead)
* `GetProfitBySymbol()` - Profit for specific symbol
* `GetDealsToday()` - Closed positions opened today (realized results)

---

## ⚠️ Common Pitfalls

### 1) Confusing floating vs realized profit

```go
// ❌ WRONG - GetProfit() is FLOATING, not realized
profit, _ := sugar.GetProfit()
fmt.Println("Money in my pocket:", profit) // WRONG! It's not realized yet!

// ✅ CORRECT - for realized profit, use history
realizedToday, _ := sugar.GetProfitToday() // Sum of Profit of closed positions opened today (swap/commission not included)
```

---

## 💎 Pro Tips

1. **Floating = not yours yet** - Only realized profit is real
2. **Monitor vs balance %** - Track profit as % of balance
3. **Set profit targets** - Close all when target reached
4. **Set loss limits** - Close all if loss exceeds limit
5. **Changes every tick** - Profit updates in real-time

---

**See also:** [`GetEquity.md`](GetEquity.md), [`GetProfitToday.md`](../8.%20History_Statistics/GetProfitToday.md)
