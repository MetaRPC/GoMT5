# 💸 Get Free Margin (`GetFreeMargin`)

> **Sugar method:** Returns margin available for opening new positions.

**API Information:**

* **Method:** `sugar.GetFreeMargin()`
* **Timeout:** 3 seconds
* **Returns:** Free margin as `float64`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetFreeMargin() (float64, error)
```

---

## 💬 Just the Essentials

* **What it is:** Money available for opening NEW positions (Equity - Used Margin).
* **Why you need it:** Check if you can open more positions, prevent margin calls.
* **Sanity check:** Free Margin = Equity - Used Margin.

---

## 🧮 Formula

```
Free Margin = Equity - Used Margin

Example:
  Equity: $10,000
  Margin: $500
  Free Margin: $9,500
```

---

## 🔗 Usage Examples

### 1) Basic usage

```go
freeMargin, err := sugar.GetFreeMargin()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Free Margin: $%.2f\n", freeMargin)
```

---

### 2) Check before opening position

```go
symbol := "EURUSD"
volume := 0.1

// Calculate required margin (always computed for a BUY at the current Ask)
requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, volume)
freeMargin, _ := sugar.GetFreeMargin()

fmt.Printf("Required: $%.2f\n", requiredMargin)
fmt.Printf("Available: $%.2f\n", freeMargin)

if requiredMargin > freeMargin {
    fmt.Println("❌ Insufficient margin!")
    return
}

// Safe to open
ticket, _ := sugar.BuyMarket(symbol, volume)
fmt.Printf("✅ Position #%d opened\n", ticket)
```

---

### 3) Calculate max tradeable lots

```go
func GetMaxTradableLots(sugar *mt5.MT5Sugar, symbol string) float64 {
    // GetMaxLotSize already uses 80% of free margin for safety:
    // the largest volume whose margin (BUY at Ask, calculated by the server for that
    // volume) fits in free margin × 0.8, rounded down to the volume step and capped at
    // the symbol's max volume.
    // It returns 0 when even the minimum volume does not fit.
    maxLots, err := sugar.GetMaxLotSize(symbol)
    if err != nil {
        return 0
    }

    return maxLots
}

maxLots := GetMaxTradableLots(sugar, "EURUSD")
fmt.Printf("Max tradeable: %.2f lots\n", maxLots)
```

---

### 4) Margin availability check

```go
freeMargin, _ := sugar.GetFreeMargin()
equity, _ := sugar.GetEquity()

availabilityPercent := (freeMargin / equity) * 100

fmt.Printf("Free Margin:  $%.2f (%.1f%% of equity)\n",
    freeMargin, availabilityPercent)

if availabilityPercent > 80 {
    fmt.Println("✅ Plenty of margin available")
} else if availabilityPercent > 50 {
    fmt.Println("⚠️  Moderate margin usage")
} else if availabilityPercent > 20 {
    fmt.Println("🚨 High margin usage - be careful!")
} else {
    fmt.Println("🔴 CRITICAL: Very low free margin!")
}
```

---

## 🔗 Related Methods

* `GetMargin()` - Used margin
* `GetMarginLevel()` - Margin level percentage
* `GetMaxLotSize()` - Max lots you can trade
* `CanOpenPosition()` - Validate before trading ⭐

---

**See also:** [`GetMargin.md`](GetMargin.md), [`CanOpenPosition.md`](../10.%20Risk_Management/CanOpenPosition.md)
