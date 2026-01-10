# 💎 Get Account Equity (`GetEquity`)

> **Sugar method:** Returns current account equity (balance + floating P/L).

**API Information:**

* **Method:** `sugar.GetEquity()`
* **Timeout:** 3 seconds
* **Returns:** Current equity as `float64`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetEquity() (float64, error)
```

---

## 🔽 Input / ⬆️ Output

| Output | Type | Description |
|--------|------|-------------|
| `equity` | `float64` | Current account equity (Balance + Floating P/L) |
| `error` | `error` | Error if query fails |

---

## 💬 Just the Essentials

* **What it is:** Your REAL account value = Balance + Floating Profit/Loss from open positions.
* **Why you need it:** See your actual current worth, calculate margin level, monitor risk.
* **Sanity check:** Equity = Balance + Profit. If no open positions, Equity = Balance.

---

## 🧮 Formula

```
Equity = Balance + Floating Profit/Loss

Examples:
  Balance: $10,000, Open P/L: +$200  → Equity: $10,200
  Balance: $10,000, Open P/L: -$150  → Equity: $9,850
  Balance: $10,000, No positions     → Equity: $10,000
```

---

## 🎯 When to Use

✅ **Risk monitoring** - Check if equity is dropping dangerously

✅ **Margin level calculation** - Equity / Margin × 100%

✅ **Stop-loss for account** - Close all if equity drops below X

✅ **Performance tracking** - Real-time account value

✅ **Margin call prevention** - Monitor equity vs margin

---

## 🔗 Usage Examples

### 1) Basic usage

```go
equity, err := sugar.GetEquity()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Equity: $%.2f\n", equity)
// Output: Equity: $10250.00
```

---

### 2) Compare Balance vs Equity

```go
balance, _ := sugar.GetBalance()
equity, _ := sugar.GetEquity()
profit, _ := sugar.GetProfit()

fmt.Printf("Balance:  $%.2f\n", balance)
fmt.Printf("Equity:   $%.2f\n", equity)
fmt.Printf("Profit:   $%.2f\n", profit)
fmt.Printf("Verify:   %.2f + %.2f = %.2f ✅\n", balance, profit, equity)

// Output:
// Balance:  $10000.00
// Equity:   $10250.00
// Profit:   $250.00
// Verify:   10000.00 + 250.00 = 10250.00 ✅
```

---

### 3) Monitor drawdown

```go
initialEquity, _ := sugar.GetEquity()
maxDrawdownPercent := 10.0 // Stop if lose 10%

for {
    currentEquity, _ := sugar.GetEquity()
    drawdown := (initialEquity - currentEquity) / initialEquity * 100

    fmt.Printf("Equity: $%.2f (Drawdown: %.2f%%)\n", currentEquity, drawdown)

    if drawdown > maxDrawdownPercent {
        fmt.Printf("🚨 Max drawdown reached! Closing all positions...\n")
        sugar.CloseAllPositions()
        break
    }

    time.Sleep(10 * time.Second)
}
```

---

### 4) Calculate margin level

```go
equity, _ := sugar.GetEquity()
margin, _ := sugar.GetMargin()

if margin == 0 {
    fmt.Println("No positions open - margin level: N/A")
} else {
    marginLevel := (equity / margin) * 100
    fmt.Printf("Equity:        $%.2f\n", equity)
    fmt.Printf("Margin:        $%.2f\n", margin)
    fmt.Printf("Margin Level:  %.2f%%\n", marginLevel)

    if marginLevel < 100 {
        fmt.Println("⚠️  WARNING: Margin level below 100%!")
    }
}

// Output:
// Equity:        $10200.00
// Margin:        $500.00
// Margin Level:  2040.00%
```

---

### 5) Equity-based stop loss

```go
func MonitorEquityStopLoss(sugar *mt5.MT5Sugar, stopLossEquity float64) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        equity, _ := sugar.GetEquity()

        if equity < stopLossEquity {
            fmt.Printf("🚨 Equity $%.2f below stop loss $%.2f\n",
                equity, stopLossEquity)
            fmt.Println("   Closing all positions...")

            sugar.CloseAllPositions()
            fmt.Println("   All positions closed!")
            return
        }
    }
}

// Run in background
go MonitorEquityStopLoss(sugar, 9000.0) // Stop if equity drops below $9000
```

---

### 6) Performance dashboard

```go
balance, _ := sugar.GetBalance()
equity, _ := sugar.GetEquity()
profit, _ := sugar.GetProfit()
profitPercent := (profit / balance) * 100

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║         ACCOUNT PERFORMANCE           ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Balance:       $%.2f\n", balance)
fmt.Printf("Equity:        $%.2f\n", equity)
fmt.Printf("Floating P/L:  $%.2f (%.2f%%)\n", profit, profitPercent)

if profit > 0 {
    fmt.Println("Status:        ✅ PROFIT")
} else if profit < 0 {
    fmt.Println("Status:        ❌ LOSS")
} else {
    fmt.Println("Status:        ➖ BREAKEVEN")
}
```

---

## 🔗 Related Methods

* `GetBalance()` - Account balance (without open positions)
* `GetProfit()` - Floating profit/loss
* `GetMargin()` - Used margin
* `GetMarginLevel()` - Margin level percentage
* `GetAccountInfo()` - All account data at once ⭐

---

## ⚠️ Common Pitfalls

### 1) Using Equity for position sizing

```go
// ❌ WRONG - don't use Equity for risk calculation
equity, _ := sugar.GetEquity()
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50) // Uses Balance internally ✅

// Equity changes with every tick - use Balance for consistent risk!
```

### 2) Not monitoring equity in live trading

```go
// ❌ WRONG - only checking balance
balance, _ := sugar.GetBalance()
// Missing: equity might be much lower due to open losses!

// ✅ CORRECT - monitor equity for real account value
equity, _ := sugar.GetEquity()
```

---

## 💎 Pro Tips

1. **Equity is reality** - This is your real account value RIGHT NOW
2. **Monitor for margin calls** - If Equity/Margin < 100%, danger!
3. **Use for drawdown limits** - Stop trading if equity drops X%
4. **Balance for position sizing** - But equity for risk monitoring
5. **Equity < Balance = losing** - Floating losses on open positions

---

**See also:** [`GetBalance.md`](GetBalance.md), [`GetProfit.md`](GetProfit.md), [`GetMarginLevel.md`](GetMarginLevel.md)
