# 🔒 Close Position (`ClosePosition`)

> **Sugar method:** Closes an open position completely by ticket number.

**API Information:**

* **Method:** `sugar.ClosePosition(ticket)`
* **Timeout:** 10 seconds
* **Returns:** Error if close failed

---

## 📋 Method Signature

```go
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `ticket` | `uint64` | Position ticket number to close |

| Output | Type | Description |
|--------|------|-------------|
| `error` | `error` | `nil` if successful, error if failed |

---

## 💬 Just the Essentials

* **What it is:** Closes entire position immediately at current market price.
* **Why you need it:** Exit trades manually, cut losses, take profits.
* **Sanity check:** Position must exist and be open. Returns error if position not found.

---

## 🎯 When to Use

✅ **Manual exit** - Close position when you want to exit

✅ **Stop loss hit** - Manually close after SL triggered

✅ **Take profit** - Lock in gains

✅ **Emergency close** - Get out of bad trade quickly

---

## 🔗 Usage Examples

### 1) Basic usage - close by ticket

```go
ticket := uint64(12345) // Position ticket from opening

err := sugar.ClosePosition(ticket)
if err != nil {
    fmt.Printf("Failed to close: %v\n", err)
    return
}

fmt.Printf("✅ Position #%d closed successfully\n", ticket)
```

---

### 2) Close and verify

```go
ticket := uint64(12345)

fmt.Printf("Closing position #%d...\n", ticket)

err := sugar.ClosePosition(ticket)
if err != nil {
    fmt.Printf("❌ Close failed: %v\n", err)
    return
}

fmt.Println("✅ Close successful")

// Verify position no longer exists
time.Sleep(1 * time.Second)
_, err = sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Println("✅ Confirmed - position no longer exists")
} else {
    fmt.Println("⚠️  Warning - position still exists")
}
```

---

### 3) Close with profit check

```go
ticket := uint64(12345)

// Check profit before closing
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

fmt.Printf("Position #%d status:\n", ticket)
fmt.Printf("  Symbol: %s\n", pos.Symbol)
fmt.Printf("  Volume: %.2f\n", pos.Volume)
fmt.Printf("  Profit: $%.2f\n", pos.Profit)

// Close position
err = sugar.ClosePosition(ticket)
if err != nil {
    fmt.Printf("Close failed: %v\n", err)
    return
}

if pos.Profit > 0 {
    fmt.Printf("✅ Profit locked: $%.2f\n", pos.Profit)
} else {
    fmt.Printf("✅ Loss cut: $%.2f\n", pos.Profit)
}
```

---

### 4) Close at target profit

```go
ticket := uint64(12345)
targetProfit := 100.0 // $100

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

fmt.Printf("Monitoring position #%d for $%.2f profit...\n", ticket, targetProfit)

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        fmt.Println("Position no longer exists")
        return
    }

    fmt.Printf("Current profit: $%.2f\n", pos.Profit)

    if pos.Profit >= targetProfit {
        fmt.Printf("🎯 Target reached! Closing...\n")
        sugar.ClosePosition(ticket)
        fmt.Printf("✅ Position closed at $%.2f profit\n", pos.Profit)
        return
    }
}
```

---

### 5) Close on stop loss breach

```go
ticket := uint64(12345)
maxLoss := -50.0 // Stop at $50 loss

ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()

fmt.Printf("Monitoring position #%d with $%.2f max loss...\n", ticket, maxLoss)

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        return
    }

    if pos.Profit <= maxLoss {
        fmt.Printf("🛑 Stop loss hit: $%.2f\n", pos.Profit)
        err := sugar.ClosePosition(ticket)
        if err != nil {
            fmt.Printf("❌ Emergency close failed: %v\n", err)
        } else {
            fmt.Println("✅ Position closed - loss limited")
        }
        return
    }

    fmt.Printf("Current P/L: $%.2f (SL: $%.2f)\n", pos.Profit, maxLoss)
}
```

---

### 6) Close all positions for a symbol

```go
symbol := "EURUSD"

positions, _ := sugar.GetPositionsBySymbol(symbol)

fmt.Printf("Closing %d %s positions...\n", len(positions), symbol)

closedCount := 0
for _, pos := range positions {
    err := sugar.ClosePosition(pos.Ticket)
    if err != nil {
        fmt.Printf("❌ Failed to close #%d: %v\n", pos.Ticket, err)
        continue
    }

    closedCount++
    fmt.Printf("✅ Closed #%d (%.2f lots, $%.2f P/L)\n",
        pos.Ticket, pos.Volume, pos.Profit)
}

fmt.Printf("\n%d/%d positions closed\n", closedCount, len(positions))
```

---

### 7) Close oldest position first

```go
positions, _ := sugar.GetOpenPositions()

if len(positions) == 0 {
    fmt.Println("No positions to close")
    return
}

// Find oldest position
var oldest *PositionInfo
for _, pos := range positions {
    if oldest == nil || pos.TimeOpen.Before(oldest.TimeOpen) {
        oldest = pos
    }
}

fmt.Printf("Closing oldest position #%d (opened: %s)\n",
    oldest.Ticket, oldest.TimeOpen.Format("2006-01-02 15:04:05"))

err := sugar.ClosePosition(oldest.Ticket)
if err != nil {
    fmt.Printf("Close failed: %v\n", err)
} else {
    fmt.Printf("✅ Closed - held for %v\n", time.Since(oldest.TimeOpen).Round(time.Minute))
}
```

---

### 8) Close losing positions only

```go
positions, _ := sugar.GetOpenPositions()

fmt.Println("Closing losing positions...")

for _, pos := range positions {
    if pos.Profit < 0 {
        fmt.Printf("Closing #%d: %s (loss: $%.2f)\n",
            pos.Ticket, pos.Symbol, pos.Profit)

        err := sugar.ClosePosition(pos.Ticket)
        if err != nil {
            fmt.Printf("  ❌ Failed: %v\n", err)
        } else {
            fmt.Printf("  ✅ Closed\n")
        }
    }
}
```

---

### 9) Close with retry logic

```go
func ClosePositionWithRetry(sugar *mt5.MT5Sugar, ticket uint64, maxRetries int) error {
    for attempt := 1; attempt <= maxRetries; attempt++ {
        err := sugar.ClosePosition(ticket)
        if err == nil {
            fmt.Printf("✅ Position #%d closed (attempt %d)\n", ticket, attempt)
            return nil
        }

        fmt.Printf("❌ Attempt %d failed: %v\n", attempt, err)

        if attempt < maxRetries {
            fmt.Printf("   Retrying in 2 seconds...\n")
            time.Sleep(2 * time.Second)
        }
    }

    return fmt.Errorf("failed to close after %d attempts", maxRetries)
}

// Usage:
err := ClosePositionWithRetry(sugar, 12345, 3)
```

---

### 10) Advanced close with analysis

```go
func ClosePositionWithAnalysis(sugar *mt5.MT5Sugar, ticket uint64) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      CLOSING POSITION                 ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Get position details before closing
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        return fmt.Errorf("position not found: %w", err)
    }

    fmt.Printf("Position #%d:\n", ticket)
    fmt.Printf("  Symbol:    %s\n", pos.Symbol)
    fmt.Printf("  Type:      %s\n", pos.Type)
    fmt.Printf("  Volume:    %.2f lots\n", pos.Volume)
    fmt.Printf("  Open:      %.5f\n", pos.PriceOpen)
    fmt.Printf("  Current:   %.5f\n", pos.PriceCurrent)
    fmt.Printf("  Profit:    $%.2f\n", pos.Profit)
    fmt.Printf("  Duration:  %v\n", time.Since(pos.TimeOpen).Round(time.Minute))

    // Calculate metrics
    var pips float64
    if pos.Type == "BUY" {
        pips = (pos.PriceCurrent - pos.PriceOpen) / 0.00001
    } else {
        pips = (pos.PriceOpen - pos.PriceCurrent) / 0.00001
    }

    fmt.Printf("  Pips:      %.1f\n", pips)

    // Close position
    fmt.Println("\nClosing...")
    err = sugar.ClosePosition(ticket)
    if err != nil {
        return fmt.Errorf("close failed: %w", err)
    }

    fmt.Println("\n✅ Position closed successfully")
    fmt.Printf("   Final P/L: $%.2f (%.1f pips)\n", pos.Profit, pips)

    return nil
}

// Usage:
ClosePositionWithAnalysis(sugar, 12345)
```

---

## 🔗 Related Methods

**🍬 Other close methods:**

* `ClosePositionPartial()` - Close only part of position
* `CloseAllPositions()` - Close all positions at once

**🍬 Position info:**

* `GetPositionByTicket()` - Get position details before closing
* `GetOpenPositions()` - List all positions

---

## ⚠️ Common Pitfalls

### 1) Not checking if position exists

```go
// ❌ WRONG - blindly closing
sugar.ClosePosition(12345) // Might not exist!

// ✅ CORRECT - check first
pos, err := sugar.GetPositionByTicket(12345)
if err != nil {
    fmt.Println("Position doesn't exist")
    return
}
sugar.ClosePosition(12345)
```

### 2) Ignoring close errors

```go
// ❌ WRONG - ignoring errors
sugar.ClosePosition(ticket)
fmt.Println("Closed!") // Maybe not!

// ✅ CORRECT - handle errors
err := sugar.ClosePosition(ticket)
if err != nil {
    fmt.Printf("Failed: %v\n", err)
    return
}
```

### 3) Closing during high spread

```go
// ❌ WRONG - close regardless of spread
sugar.ClosePosition(ticket)

// ✅ BETTER - check spread first
spread, _ := sugar.GetSpread("EURUSD")
if spread > 20 {
    fmt.Printf("⚠️  High spread: %.0f - wait for better conditions?\n", spread)
}
// Then decide whether to close
```

---

## 💎 Pro Tips

1. **Check profit first** - Know what you're closing

2. **Verify close** - Check position no longer exists after

3. **Handle errors** - Market might be closed, position might not exist

4. **Be patient** - Sometimes close takes a few seconds to execute

5. **Consider spread** - Closing during wide spread = worse exit price

---

**See also:** [`ClosePositionPartial.md`](ClosePositionPartial.md), [`CloseAllPositions.md`](CloseAllPositions.md), [`GetPositionByTicket.md`](GetPositionByTicket.md)
