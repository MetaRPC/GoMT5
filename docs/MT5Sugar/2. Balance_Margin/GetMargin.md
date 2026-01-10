# 🔒 Get Used Margin (`GetMargin`)

> **Sugar method:** Returns margin currently used by open positions.

**API Information:**

* **Method:** `sugar.GetMargin()`
* **Timeout:** 3 seconds
* **Returns:** Used margin as `float64`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetMargin() (float64, error)
```

---

## 💬 Just the Essentials

* **What it is:** Money "locked" as collateral for your open positions.
* **Why you need it:** Calculate how much margin you're using, prevent margin calls.
* **Sanity check:** No open positions = Margin is 0. More positions = more margin used.

---

## 🧮 Formula

```
Margin = Sum of margin required for all open positions

Margin Level = (Equity / Margin) × 100%

If Margin Level < 100% → Margin Call risk!
```

---

## 🔗 Usage Examples

### 1) Basic usage

```go
margin, err := sugar.GetMargin()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Used Margin: $%.2f\n", margin)
```

---

### 2) Check margin usage

```go
equity, _ := sugar.GetEquity()
margin, _ := sugar.GetMargin()
freeMargin, _ := sugar.GetFreeMargin()

fmt.Printf("Equity:       $%.2f\n", equity)
fmt.Printf("Used Margin:  $%.2f (%.1f%%)\n", margin, (margin/equity)*100)
fmt.Printf("Free Margin:  $%.2f (%.1f%%)\n", freeMargin, (freeMargin/equity)*100)

// Output:
// Equity:       $10000.00
// Used Margin:  $500.00 (5.0%)
// Free Margin:  $9500.00 (95.0%)
```

---

### 3) Prevent over-leverage

```go
func CanOpenNewPosition(sugar *mt5.MT5Sugar, requiredMargin float64) bool {
    freeMargin, _ := sugar.GetFreeMargin()

    safetyBuffer := 0.2 // Keep 20% buffer
    availableMargin := freeMargin * (1 - safetyBuffer)

    if requiredMargin > availableMargin {
        fmt.Printf("❌ Cannot open: Need $%.2f, Have $%.2f (with buffer)\n",
            requiredMargin, availableMargin)
        return false
    }

    return true
}
```

---

### 4) Margin usage dashboard

```go
equity, _ := sugar.GetEquity()
margin, _ := sugar.GetMargin()
marginLevel, _ := sugar.GetMarginLevel()
positions, _ := sugar.CountOpenPositions()

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║          MARGIN STATUS                ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Open Positions:  %d\n", positions)
fmt.Printf("Equity:          $%.2f\n", equity)
fmt.Printf("Margin Used:     $%.2f\n", margin)
fmt.Printf("Margin Level:    %.2f%%\n", marginLevel)

if marginLevel < 200 {
    fmt.Println("⚠️  WARNING: Low margin level!")
} else {
    fmt.Println("✅ Margin level healthy")
}
```

---

## 🔗 Related Methods

* `GetFreeMargin()` - Available margin
* `GetMarginLevel()` - Margin level percentage
* `GetEquity()` - Total account value
* `CalculateRequiredMargin()` - Margin needed for a position

---

**See also:** [`GetFreeMargin.md`](GetFreeMargin.md), [`GetMarginLevel.md`](GetMarginLevel.md)
