# 💰 Get Account Balance (`GetBalance`)

> **Sugar method:** Returns current account balance in one line.

**API Information:**

* **Method:** `sugar.GetBalance()`
* **Timeout:** 3 seconds
* **Returns:** Current balance as `float64`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetBalance() (float64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| *None* | - | No parameters required |

| Output | Type | Description |
|--------|------|-------------|
| `balance` | `float64` | Current account balance |
| `error` | `error` | Error if query fails |

---

## 💬 Just the Essentials

* **What it is:** Gets your account balance - the money you have before considering open positions.
* **Why you need it:** Check available funds, calculate risk amounts, validate before trading.
* **Sanity check:** Balance ≤ Equity (equity includes floating P/L).

---

## 🎯 When to Use

✅ **Before trading** - Check if you have enough funds

✅ **Risk calculation** - Calculate position size based on balance

✅ **Monitoring** - Track account value

✅ **Reporting** - Generate account reports

✅ **Validation** - Ensure sufficient funds for trading

---

## 🔗 Usage Examples

### 1) Basic usage

```go
balance, err := sugar.GetBalance()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Balance: %.2f\n", balance)
// Output: Balance: 10000.00
```

---

### 2) Check before trading

```go
balance, _ := sugar.GetBalance()

if balance < 1000 {
    fmt.Println("❌ Insufficient balance for trading")
    return
}

fmt.Println("✅ Balance sufficient - proceeding with trade")
```

---

### 3) Calculate risk amount

```go
balance, _ := sugar.GetBalance()
riskPercent := 2.0

riskAmount := balance * riskPercent / 100.0

fmt.Printf("Balance:     $%.2f\n", balance)
fmt.Printf("Risk (2%%):   $%.2f\n", riskAmount)

// Better: use CalculatePositionSize()
lotSize, _ := sugar.CalculatePositionSize("EURUSD", riskPercent, 50)
fmt.Printf("Lot size:    %.2f\n", lotSize)
```

---

### 4) Track balance changes

```go
initialBalance, _ := sugar.GetBalance()
fmt.Printf("Starting balance: $%.2f\n", initialBalance)

// ... Trading operations ...
time.Sleep(1 * time.Hour)

currentBalance, _ := sugar.GetBalance()
change := currentBalance - initialBalance
changePercent := (change / initialBalance) * 100

fmt.Printf("\nFinal balance: $%.2f\n", currentBalance)
fmt.Printf("Change:        $%.2f (%.2f%%)\n", change, changePercent)

// Output:
// Starting balance: $10000.00
// Final balance: $10250.00
// Change:        $250.00 (2.50%)
```

---

### 5) Account dashboard

```go
balance, _ := sugar.GetBalance()
equity, _ := sugar.GetEquity()
margin, _ := sugar.GetMargin()
freeMargin, _ := sugar.GetFreeMargin()
profit, _ := sugar.GetProfit()

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║       ACCOUNT DASHBOARD               ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Balance:       $%.2f\n", balance)
fmt.Printf("Equity:        $%.2f\n", equity)
fmt.Printf("Profit/Loss:   $%.2f\n", profit)
fmt.Printf("Margin Used:   $%.2f\n", margin)
fmt.Printf("Free Margin:   $%.2f\n", freeMargin)
```

---

### 6) Validate minimum balance

```go
func ValidateMinimumBalance(sugar *mt5.MT5Sugar, minBalance float64) error {
    balance, err := sugar.GetBalance()
    if err != nil {
        return fmt.Errorf("failed to get balance: %w", err)
    }

    if balance < minBalance {
        return fmt.Errorf("balance $%.2f below minimum $%.2f", balance, minBalance)
    }

    return nil
}

// Usage:
err := ValidateMinimumBalance(sugar, 1000.0)
if err != nil {
    fmt.Printf("❌ %v\n", err)
    return
}

fmt.Println("✅ Balance check passed")
```

---

## 🔗 Related Methods

**🍬 Other balance methods:**

* `GetEquity()` - Balance + floating P/L
* `GetMargin()` - Used margin
* `GetFreeMargin()` - Available margin for trading
* `GetProfit()` - Current floating profit/loss
* `GetAccountInfo()` - Get all account data at once

**💡 Recommended pattern:**
```go
// Instead of calling GetBalance, GetEquity, etc separately:
accountInfo, _ := sugar.GetAccountInfo()
// Now you have: Balance, Equity, Margin, FreeMargin, Profit, etc.
```

---

## ⚠️ Common Pitfalls

### 1) Confusing Balance vs Equity

```go
// ❌ WRONG - using Balance when you need Equity
balance, _ := sugar.GetBalance() // Doesn't include open positions!

// ✅ CORRECT - use Equity for total account value
equity, _ := sugar.GetEquity() // Includes floating P/L
```

### 2) Not checking for errors

```go
// ❌ WRONG - ignoring errors
balance, _ := sugar.GetBalance()
fmt.Printf("Balance: $%.2f\n", balance) // Might be 0 if error!

// ✅ CORRECT - check errors
balance, err := sugar.GetBalance()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
fmt.Printf("Balance: $%.2f\n", balance)
```

---

## 💎 Pro Tips

1. **Use GetAccountInfo()** - More efficient than calling multiple methods

2. **Balance for risk calculation** - Use balance (not equity) for position sizing

3. **Check before trading** - Always verify sufficient balance

4. **Track changes** - Monitor balance to measure performance

5. **Equity is reality** - Balance is historical, Equity is current

---

**See also:** [`GetEquity.md`](GetEquity.md), [`GetAccountInfo.md`](../12.%20Account_Information/GetAccountInfo.md)
