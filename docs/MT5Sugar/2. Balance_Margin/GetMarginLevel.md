# 📊 Get Margin Level (`GetMarginLevel`)

> **Sugar method:** Returns margin level as percentage (Equity/Margin × 100%).

**API Information:**

* **Method:** `sugar.GetMarginLevel()`
* **Timeout:** 3 seconds
* **Returns:** Margin level as `float64` (percentage)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetMarginLevel() (float64, error)
```

---

## 💬 Just the Essentials

* **What it is:** Ratio of Equity to Margin, expressed as percentage.
* **Why you need it:** Monitor margin call risk. Below 100% = margin call territory!
* **Sanity check:** Higher is better. 1000% = safe, 100% = dangerous, <100% = margin call.
* **No margin used = 0:** With no open positions the value is `0` (not infinity) - treat `0` as "no margin used",
  not as danger.
* **Under the hood:** `service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL)` - the
  value comes from the terminal as is.

---

## 🧮 Formula

```
Margin Level = (Equity / Margin) × 100%

Examples:
  Equity: $10,000, Margin: $500 → 2000% (very safe)
  Equity: $10,000, Margin: $5,000 → 200% (safe)
  Equity: $5,000, Margin: $5,000 → 100% (margin call risk!)
  Equity: $4,500, Margin: $5,000 → 90% (MARGIN CALL!)
```

---

## ⚠️ Margin Level Zones

| Level | Status | Description |
|-------|--------|-------------|
| **> 1000%** | 🟢 Excellent | Very safe, plenty of margin |
| **500-1000%** | 🟢 Good | Safe margin usage |
| **200-500%** | 🟡 Moderate | Acceptable but monitor |
| **100-200%** | 🟠 Warning | High risk - reduce positions |
| **< 100%** | 🔴 Danger | MARGIN CALL TERRITORY! |

---

## 🔗 Usage Examples

### 1) Basic usage

```go
marginLevel, err := sugar.GetMarginLevel()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Margin Level: %.2f%%\n", marginLevel)
```

---

### 2) Margin level monitoring

```go
func MonitorMarginLevel(sugar *mt5.MT5Sugar) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        marginLevel, err := sugar.GetMarginLevel()
        if err != nil {
            fmt.Printf("Margin level unavailable: %v\n", err)
            continue
        }

        if marginLevel == 0 {
            fmt.Println("No margin used (no open positions)")
            continue
        }

        fmt.Printf("Margin Level: %.2f%% ", marginLevel)

        if marginLevel > 500 {
            fmt.Println("🟢 SAFE")
        } else if marginLevel > 200 {
            fmt.Println("🟡 MODERATE")
        } else if marginLevel > 100 {
            fmt.Println("🟠 WARNING")
        } else {
            fmt.Println("🔴 DANGER - MARGIN CALL RISK!")
            if _, err := sugar.CloseAllPositions(); err != nil { // Emergency close
                fmt.Printf("Some positions not closed: %v\n", err)
            }
        }
    }
}

go MonitorMarginLevel(sugar)
```

---

### 3) Pre-trade margin check

```go
func SafeToTrade(sugar *mt5.MT5Sugar, minMarginLevel float64) bool {
    marginLevel, err := sugar.GetMarginLevel()
    if err != nil {
        return false
    }

    if marginLevel == 0 {
        return true // no margin used yet
    }

    if marginLevel < minMarginLevel {
        fmt.Printf("❌ Margin level %.2f%% below minimum %.2f%%\n",
            marginLevel, minMarginLevel)
        return false
    }

    return true
}

// Before opening new position
if !SafeToTrade(sugar, 300.0) {
    fmt.Println("Not safe to open new positions")
    return
}

sugar.BuyMarket("EURUSD", 0.1)
```

---

### 4) Complete margin dashboard

```go
equity, _ := sugar.GetEquity()
margin, _ := sugar.GetMargin()
freeMargin, _ := sugar.GetFreeMargin()
marginLevel, _ := sugar.GetMarginLevel()
positions, _ := sugar.CountOpenPositions()

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║        MARGIN LEVEL REPORT            ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Open Positions:   %d\n", positions)
fmt.Printf("Equity:           $%.2f\n", equity)
fmt.Printf("Margin Used:      $%.2f\n", margin)
fmt.Printf("Free Margin:      $%.2f\n", freeMargin)
fmt.Printf("Margin Level:     %.2f%%\n", marginLevel)
fmt.Println()

if margin == 0 {
    fmt.Println("Status: 🟢 No margin used")
} else if marginLevel >= 500 {
    fmt.Println("Status: 🟢 EXCELLENT - Very safe")
} else if marginLevel >= 200 {
    fmt.Println("Status: 🟡 MODERATE - Monitor closely")
} else if marginLevel >= 100 {
    fmt.Println("Status: 🟠 WARNING - Reduce positions!")
} else {
    fmt.Println("Status: 🔴 CRITICAL - Margin call risk!")
}
```

---

## 🔗 Related Methods

* `GetEquity()` - Total account value
* `GetMargin()` - Used margin
* `GetFreeMargin()` - Available margin
* `CloseAllPositions()` - Emergency close when margin low

---

**See also:** [`GetEquity.md`](GetEquity.md), [`GetMargin.md`](GetMargin.md)
