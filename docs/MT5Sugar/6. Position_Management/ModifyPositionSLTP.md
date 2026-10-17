# ✏️ Modify Position SL/TP (`ModifyPositionSLTP`)

> **Sugar method:** Modifies Stop Loss and Take Profit of an open position.

**API Information:**

* **Method:** `sugar.ModifyPositionSLTP(ticket, sl, tp)`
* **Timeout:** 10 seconds
* **Returns:** Error if modification failed
* **Request:** `pb.OrderModifyRequest{Ticket: ticket, StopLoss: &sl, TakeProfit: &tp}` sent via `MT5Service.ModifyOrder`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) ModifyPositionSLTP(ticket uint64, sl, tp float64) error
```

Both values are always sent, exactly as given (no rounding to the symbol's digits, no stop-level check).
To change only one of them use `ModifyPositionSL(ticket, sl)` or `ModifyPositionTP(ticket, tp)`, which send only that field.

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `ticket` | `uint64` | Position ticket number |
| `sl` | `float64` | New Stop Loss price (0 = remove SL) |
| `tp` | `float64` | New Take Profit price (0 = remove TP) |

| Output | Type | Description |
|--------|------|-------------|
| `error` | `error` | `nil` when the server returns code 10009 (done); `"ModifyPositionSLTP failed: ..."` if the call fails; `"modify rejected, code: N, comment: ..."` for any other return code |

---

## 💬 Just the Essentials

* **What it is:** Change SL and TP of an existing position.
* **Why you need it:** Trailing stop, breakeven, adjust risk after partial close.
* **Sanity check:** For BUY: SL below entry, TP above. For SELL: SL above entry, TP below.

---

## 🎯 When to Use

✅ **Move to breakeven** - Set SL to entry price once in profit

✅ **Trailing stop** - Follow price with moving SL

✅ **Adjust risk** - Change SL/TP based on market conditions

✅ **Add protection** - Add SL/TP to position that has none

---

## 🔗 Usage Examples

### 1) Basic usage - set new SL/TP

```go
ticket := uint64(12345)
newSL := 1.08000
newTP := 1.09000

err := sugar.ModifyPositionSLTP(ticket, newSL, newTP)
if err != nil {
    fmt.Printf("Modification failed: %v\n", err)
    return
}

fmt.Printf("✅ Position #%d modified\n", ticket)
fmt.Printf("   New SL: %.5f\n", newSL)
fmt.Printf("   New TP: %.5f\n", newTP)
```

---

### 2) Move Stop Loss to breakeven

```go
ticket := uint64(12345)
profitThreshold := 50.0 // $50

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

if pos.Profit >= profitThreshold {
    fmt.Printf("Profit $%.2f reached threshold $%.2f\n", pos.Profit, profitThreshold)
    fmt.Println("Moving SL to breakeven...")

    // Set SL to entry price
    err := sugar.ModifyPositionSLTP(ticket, pos.PriceOpen, pos.TakeProfit)
    if err != nil {
        fmt.Printf("Failed: %v\n", err)
        return
    }

    fmt.Printf("✅ SL moved to breakeven: %.5f\n", pos.PriceOpen)
    fmt.Println("   Now trading risk-free!")
}
```

---

### 3) Trailing stop loss

```go
ticket := uint64(12345)
trailDistance := 50.0 // 50 points

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

info, _ := sugar.GetSymbolInfo("EURUSD")
trailDistancePrice := trailDistance * info.Point

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        fmt.Println("Position closed")
        return
    }

    var newSL float64

    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        // Trail below current price
        newSL = pos.PriceCurrent - trailDistancePrice

        // Only move SL up, never down
        if newSL > pos.StopLoss {
            err := sugar.ModifyPositionSLTP(ticket, newSL, pos.TakeProfit)
            if err != nil {
                fmt.Printf("Trail failed: %v\n", err)
                continue
            }

            fmt.Printf("✅ Trailing SL: %.5f → %.5f (%.0f points)\n",
                pos.StopLoss, newSL, trailDistance)
        }
    } else { // SELL
        // Trail above current price
        newSL = pos.PriceCurrent + trailDistancePrice

        // Only move SL down, never up
        if newSL < pos.StopLoss || pos.StopLoss == 0 {
            err := sugar.ModifyPositionSLTP(ticket, newSL, pos.TakeProfit)
            if err != nil {
                fmt.Printf("Trail failed: %v\n", err)
                continue
            }

            fmt.Printf("✅ Trailing SL: %.5f → %.5f (%.0f points)\n",
                pos.StopLoss, newSL, trailDistance)
        }
    }

    fmt.Printf("Current: %.5f, SL: %.5f, Profit: $%.2f\n",
        pos.PriceCurrent, pos.StopLoss, pos.Profit)
}
```

---

### 4) Add SL/TP to position without protection

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

// Check if position has no SL/TP
if pos.StopLoss == 0 && pos.TakeProfit == 0 {
    fmt.Printf("⚠️  Position #%d has no protection!\n", ticket)
    fmt.Println("Adding SL/TP...")

    // CalculateSLTP takes a direction string: "BUY", anything else is treated as SELL
    direction := "SELL"
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        direction = "BUY"
    }

    // Distances are in points; entry 0 = current Ask (BUY) / Bid (SELL)
    sl, tp, _ := sugar.CalculateSLTP(pos.Symbol, direction, 0, 50, 100)

    err := sugar.ModifyPositionSLTP(ticket, sl, tp)
    if err != nil {
        fmt.Printf("Failed: %v\n", err)
        return
    }

    fmt.Printf("✅ Protection added:\n")
    fmt.Printf("   SL: %.5f (50 points)\n", sl)
    fmt.Printf("   TP: %.5f (100 points)\n", tp)
}
```

---

### 5) Widen Stop Loss (reduce risk after profit)

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

if pos.Profit > 100.0 {
    fmt.Printf("Large profit: $%.2f - widening SL\n", pos.Profit)

    info, _ := sugar.GetSymbolInfo(pos.Symbol)

    var newSL float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        // Move SL closer to entry (widen stop distance)
        newSL = pos.PriceOpen - (20 * info.Point) // 20 points instead of original 50
    } else {
        newSL = pos.PriceOpen + (20 * info.Point)
    }

    err := sugar.ModifyPositionSLTP(ticket, newSL, pos.TakeProfit)
    if err != nil {
        fmt.Printf("Failed: %v\n", err)
        return
    }

    fmt.Printf("✅ SL widened to give position more room\n")
    fmt.Printf("   New SL: %.5f\n", newSL)
}
```

---

### 6) Tighten Take Profit

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

// If 80% to target, tighten TP to lock profit sooner
if pos.TakeProfit > 0 {
    var distanceToTP float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        distanceToTP = pos.TakeProfit - pos.PriceCurrent
        targetDistance := pos.TakeProfit - pos.PriceOpen

        if distanceToTP < (targetDistance * 0.2) { // 80% there
            // Tighten TP to halfway between current and original TP
            newTP := pos.PriceCurrent + (distanceToTP * 0.5)

            err := sugar.ModifyPositionSLTP(ticket, pos.StopLoss, newTP)
            if err != nil {
                fmt.Printf("Failed: %v\n", err)
                return
            }

            fmt.Printf("✅ TP tightened to lock profit\n")
            fmt.Printf("   Old TP: %.5f\n", pos.TakeProfit)
            fmt.Printf("   New TP: %.5f\n", newTP)
        }
    }
}
```

---

### 7) Remove SL/TP completely

```go
ticket := uint64(12345)

// Remove both SL and TP (set to 0)
err := sugar.ModifyPositionSLTP(ticket, 0, 0)
if err != nil {
    fmt.Printf("Failed: %v\n", err)
    return
}

fmt.Printf("✅ Position #%d: SL and TP removed\n", ticket)
fmt.Println("⚠️  WARNING: Position now has no protection!")
```

---

### 8) Step trailing (move in increments)

```go
ticket := uint64(12345)
stepSize := 20.0 // Move SL every 20 points of profit

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    return
}
info, _ := sugar.GetSymbolInfo(pos.Symbol)
stepPrice := stepSize * info.Point

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

lastSLUpdate := 0.0

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil { // position closed (or query failed)
        return
    }

    var profitPips float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        profitPips = (pos.PriceCurrent - pos.PriceOpen) / info.Point
    } else {
        profitPips = (pos.PriceOpen - pos.PriceCurrent) / info.Point
    }

    // Move SL in 20 point steps
    steps := math.Floor(profitPips / stepSize)

    if steps > lastSLUpdate && steps > 0 {
        var newSL float64
        if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
            newSL = pos.PriceOpen + (steps * stepPrice)
        } else {
            newSL = pos.PriceOpen - (steps * stepPrice)
        }

        err := sugar.ModifyPositionSLTP(ticket, newSL, pos.TakeProfit)
        if err == nil {
            fmt.Printf("✅ Step %d: SL moved to %.5f\n", int(steps), newSL)
            lastSLUpdate = steps
        }
    }

    fmt.Printf("Profit: %.1f pips, Steps: %.0f, SL: %.5f\n",
        profitPips, steps, pos.StopLoss)
}
```

---

### 9) Time-based SL tightening

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}
openDuration := time.Since(pos.OpenTimeUtc.AsTime()) // OpenTime is broker time, OpenTimeUtc is UTC

// After 2 hours, tighten SL
if openDuration > 2*time.Hour {
    fmt.Printf("Position open for %v - tightening SL\n", openDuration.Round(time.Minute))

    info, _ := sugar.GetSymbolInfo(pos.Symbol)

    var newSL float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        // Move SL to 20 points below current (was 50 points)
        newSL = pos.PriceCurrent - (20 * info.Point)

        if newSL > pos.StopLoss {
            err := sugar.ModifyPositionSLTP(ticket, newSL, pos.TakeProfit)
            if err != nil {
                fmt.Printf("Failed: %v\n", err)
                return
            }

            fmt.Printf("✅ SL tightened after 2 hours\n")
            fmt.Printf("   New SL: %.5f\n", newSL)
        }
    }
}
```

---

### 10) Advanced modification function

```go
// slPoints/tpPoints are distances in points (CalculateSLTP multiplies them by Point)
func ModifySLTP(sugar *mt5.MT5Sugar, ticket uint64, slPoints float64, tpPoints float64, reason string) error {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        return fmt.Errorf("position not found: %w", err)
    }

    direction := "SELL"
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        direction = "BUY"
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      MODIFYING POSITION               ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Reason:   %s\n", reason)
    fmt.Printf("Position: #%d\n", ticket)
    fmt.Printf("Symbol:   %s %s\n", pos.Symbol, direction)
    fmt.Printf("Entry:    %.5f\n", pos.PriceOpen)
    fmt.Printf("Current:  %.5f\n\n", pos.PriceCurrent)

    // Calculate new SL/TP from the entry price
    sl, tp, err := sugar.CalculateSLTP(pos.Symbol, direction, pos.PriceOpen, slPoints, tpPoints)
    if err != nil {
        return err
    }

    fmt.Printf("Old SL:   %.5f\n", pos.StopLoss)
    fmt.Printf("New SL:   %.5f (%.0f points)\n", sl, slPoints)
    fmt.Printf("Old TP:   %.5f\n", pos.TakeProfit)
    fmt.Printf("New TP:   %.5f (%.0f points)\n\n", tp, tpPoints)

    // Validate new SL is better (tighter or breakeven)
    if direction == "BUY" {
        if sl < pos.StopLoss && pos.StopLoss > 0 {
            return fmt.Errorf("new SL %.5f is worse than current %.5f", sl, pos.StopLoss)
        }
    } else {
        if sl > pos.StopLoss && pos.StopLoss > 0 {
            return fmt.Errorf("new SL %.5f is worse than current %.5f", sl, pos.StopLoss)
        }
    }

    // Apply modification
    err = sugar.ModifyPositionSLTP(ticket, sl, tp)
    if err != nil {
        return fmt.Errorf("modification failed: %w", err)
    }

    fmt.Println("✅ Position modified successfully")
    return nil
}

// Usage:
ModifySLTP(sugar, 12345, 30, 100, "Moving to tighter stop")
```

---

## 🔗 Related Methods

**🍬 Helper methods:**

* `CalculateSLTP()` - Calculate SL/TP prices from distances in points
* `GetPositionByTicket()` - Get current position details
* `ModifyPositionSL()` / `ModifyPositionTP()` - Change only one level

**🍬 Opening with SL/TP:**

* `BuyMarketWithSLTP()` - Open with SL/TP from start
* `SellMarketWithSLTP()` - Open with SL/TP from start

---

## ⚠️ Common Pitfalls

### 1) Wrong SL direction

```go
// ❌ WRONG - SL in wrong direction for BUY
// BUY position, setting SL above entry!
sugar.ModifyPositionSLTP(ticket, 1.09000, tp) // Wrong!

// ✅ CORRECT - SL below entry for BUY
sugar.ModifyPositionSLTP(ticket, 1.08000, tp)
```

### 2) Moving SL in wrong direction

```go
// ❌ WRONG - making SL worse
oldSL := 1.08000
newSL := 1.07500 // Moving SL further from current price!
sugar.ModifyPositionSLTP(ticket, newSL, tp)

// ✅ CORRECT - only improve SL (move closer to profit)
if newSL > oldSL { // For BUY positions
    sugar.ModifyPositionSLTP(ticket, newSL, tp)
}
```

### 3) Not validating minimum stop level

```go
// ❌ WRONG - SL might be too close
sugar.ModifyPositionSLTP(ticket, currentPrice-0.00001, tp)

// ✅ CORRECT - check minimum stop level
stopLevel, _ := sugar.GetMinStopLevel(symbol)
minDistance := float64(stopLevel) * point
// Ensure new SL respects minimum distance
```

---

## 💎 Pro Tips

1. **Breakeven first** - Move to breakeven once the position is comfortably in profit

2. **Trail wisely** - Don't trail too tight (give position room)

3. **Only improve SL** - Never move SL to worse position

4. **Step trailing** - Move in increments, not continuously

5. **Time-based** - Tighten SL after position open for hours

---

**See also:** [`GetPositionByTicket.md`](GetPositionByTicket.md), [`CalculateSLTP.md`](../11. Trading_Helpers/CalculateSLTP.md), [`ClosePosition.md`](ClosePosition.md)
