# 🔓 Close Position Partial (`ClosePositionPartial`)

> **Sugar method:** Closes part of a position, leaving the rest open.

**API Information:**

* **Method:** `sugar.ClosePositionPartial(ticket, volume)`
* **Timeout:** 30 seconds
* **Returns:** Error if close failed
* **Request:** `pb.OrderCloseRequest{Ticket: ticket, Volume: volume}` sent via `MT5Service.CloseOrder` (MT5 `OrderClose`)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `ticket` | `uint64` | Position ticket number |
| `volume` | `float64` | Volume to close (must be less than position volume) |

| Output | Type | Description |
|--------|------|-------------|
| `error` | `error` | `nil` when the server returns code 10009 (done); `"ClosePositionPartial failed: ..."` if the call fails; `"partial close rejected, code: N (meaning)"` for any other return code (the meaning comes from `helpers.GetRetCodeMessage`) |

---

## 💬 Just the Essentials

* **What it is:** Close part of position, keep rest running.
* **Why you need it:** Lock in partial profits while letting winners run.
* **Sanity check:** Volume must be less than position size. Position stays open with reduced volume.
* **No normalization:** the volume is sent exactly as given - it is not rounded to the symbol's volume step and not checked against the position volume. Round it yourself (e.g. `pos.Volume / 2` of a 0.03-lot position is 0.015, which most symbols reject); an invalid volume comes back as an error.

---

## 🎯 When to Use

✅ **Scale out** - Take partial profits at targets

✅ **Risk reduction** - Close half, move SL to breakeven

✅ **Pyramid exit** - Exit scaled positions gradually

✅ **Lock profits** - Secure gains while staying in trade

---

## 🔗 Usage Examples

### 1) Basic usage - close half position

```go
ticket := uint64(12345)

// Get current position
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}
fmt.Printf("Position: %.2f lots\n", pos.Volume)

// Close half (rounded down to a 0.01 volume step - the method does not round)
halfVolume := math.Floor(pos.Volume/2/0.01) * 0.01
err = sugar.ClosePositionPartial(ticket, halfVolume)
if err != nil {
    fmt.Printf("Partial close failed: %v\n", err)
    return
}

fmt.Printf("✅ Closed %.2f lots\n", halfVolume)
fmt.Printf("   Remaining: %.2f lots\n", pos.Volume-halfVolume)
```

---

### 2) Take profits at first target

```go
ticket := uint64(12345)
firstTargetProfit := 50.0 // $50

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

alreadyScaled := false

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil { // "position with ticket N not found" once it is closed
        fmt.Printf("Stopped: %v\n", err)
        return
    }

    if !alreadyScaled && pos.Profit >= firstTargetProfit {
        fmt.Printf("🎯 First target reached: $%.2f\n", pos.Profit)

        // Close 50% of position
        halfVolume := pos.Volume / 2
        err := sugar.ClosePositionPartial(ticket, halfVolume)
        if err != nil {
            fmt.Printf("Failed to scale out: %v\n", err)
            continue
        }

        fmt.Printf("✅ Scaled out 50%% at $%.2f profit\n", pos.Profit)
        fmt.Printf("   Remaining: %.2f lots\n", halfVolume)
        alreadyScaled = true
    }

    fmt.Printf("Current P/L: $%.2f\n", pos.Profit)
}
```

---

### 3) Three-target strategy

```go
ticket := uint64(12345)

// Close 1/3 at each target
targets := []float64{50.0, 100.0, 150.0}
closedTargets := 0

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    return
}
originalVolume := pos.Volume
volumePerTarget := originalVolume / 3.0 // round to the symbol's volume step in real code

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil { // position fully closed (or query failed)
        return
    }

    if closedTargets < len(targets) && pos.Profit >= targets[closedTargets] {
        targetNum := closedTargets + 1

        fmt.Printf("🎯 Target %d reached: $%.2f\n", targetNum, pos.Profit)

        err := sugar.ClosePositionPartial(ticket, volumePerTarget)
        if err != nil {
            fmt.Printf("Failed to close: %v\n", err)
            continue
        }

        fmt.Printf("✅ Closed 1/3 at target %d\n", targetNum)
        closedTargets++

        if closedTargets == len(targets) {
            fmt.Println("All targets hit - position fully closed")
            return
        }
    }
}
```

---

### 4) Close half, move SL to breakeven

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

fmt.Printf("Position #%d:\n", ticket)
fmt.Printf("  Volume: %.2f\n", pos.Volume)
fmt.Printf("  Entry:  %.5f\n", pos.PriceOpen)
fmt.Printf("  Profit: $%.2f\n", pos.Profit)

// Close half (rounded down to a 0.01 volume step)
halfVolume := math.Floor(pos.Volume/2/0.01) * 0.01
err = sugar.ClosePositionPartial(ticket, halfVolume)
if err != nil {
    fmt.Printf("Partial close failed: %v\n", err)
    return
}

fmt.Printf("✅ Closed half position (%.2f lots)\n", halfVolume)

// Move SL to breakeven
err = sugar.ModifyPositionSLTP(ticket, pos.PriceOpen, pos.TakeProfit)
if err != nil {
    fmt.Printf("⚠️  Failed to move SL: %v\n", err)
} else {
    fmt.Printf("✅ SL moved to breakeven: %.5f\n", pos.PriceOpen)
}

fmt.Println("\n🎯 Now trading risk-free!")
```

---

### 5) Progressive scaling out

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    return
}
originalVolume := pos.Volume

// Close 25% at each profit milestone
milestones := []float64{30.0, 60.0, 100.0}
closePercents := []float64{0.25, 0.33, 0.50} // 25%, 33%, 50% of remaining

closedMilestones := 0

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil { // position closed (or query failed)
        return
    }

    if closedMilestones < len(milestones) && pos.Profit >= milestones[closedMilestones] {
        percent := closePercents[closedMilestones]
        closeVolume := math.Floor(pos.Volume*percent/0.01) * 0.01 // 0.01 volume step

        fmt.Printf("🎯 Milestone %d: $%.2f\n", closedMilestones+1, pos.Profit)
        fmt.Printf("   Closing %.0f%% (%.2f lots)\n", percent*100, closeVolume)

        err := sugar.ClosePositionPartial(ticket, closeVolume)
        if err != nil {
            fmt.Printf("Failed: %v\n", err)
            continue
        }

        closedMilestones++
        newVolume := pos.Volume - closeVolume

        fmt.Printf("✅ Scaled out - remaining: %.2f lots (%.0f%% of original)\n",
            newVolume, (newVolume/originalVolume)*100)
    }
}
```

---

### 6) Partial close on pip target

```go
ticket := uint64(12345)
targetPips := 50.0

ticker := time.NewTicker(3 * time.Second)
defer ticker.Stop()

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil { // position closed (or query failed)
        return
    }

    var currentPips float64 // (EURUSD: 1 pip = 0.0001 = 10 points)
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        currentPips = (pos.PriceCurrent - pos.PriceOpen) / 0.0001
    } else {
        currentPips = (pos.PriceOpen - pos.PriceCurrent) / 0.0001
    }

    fmt.Printf("Current: %.1f pips\n", currentPips)

    if currentPips >= targetPips {
        fmt.Printf("🎯 Target reached: %.1f pips\n", currentPips)

        // Close 2/3 of position (rounded down to a 0.01 volume step)
        closeVolume := math.Floor(pos.Volume*0.66/0.01) * 0.01
        err := sugar.ClosePositionPartial(ticket, closeVolume)
        if err != nil {
            fmt.Printf("Failed: %v\n", err)
            continue
        }

        fmt.Printf("✅ Closed 66%% at %.1f pips\n", currentPips)
        return
    }
}
```

---

### 7) Reduce risk on drawdown

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    return
}
originalVolume := pos.Volume

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil { // position closed (or query failed)
        return
    }

    // If losing more than $30, reduce position size
    if pos.Profit < -30.0 && pos.Volume == originalVolume {
        // Close 50% to reduce risk (rounded down to a 0.01 volume step)
        closeVolume := math.Floor(pos.Volume/2/0.01) * 0.01

        fmt.Printf("⚠️  Drawdown: $%.2f - reducing position\n", pos.Profit)

        err := sugar.ClosePositionPartial(ticket, closeVolume)
        if err != nil {
            fmt.Printf("Failed to reduce: %v\n", err)
            continue
        }

        fmt.Printf("✅ Reduced position by 50%%\n")
        fmt.Printf("   New size: %.2f lots\n", pos.Volume-closeVolume)
    }

    fmt.Printf("P/L: $%.2f, Volume: %.2f\n", pos.Profit, pos.Volume)
}
```

---

### 8) Close specific lot amount

```go
ticket := uint64(12345)
lotsToClose := 0.05 // Close exactly 0.05 lots

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

fmt.Printf("Position: %.2f lots\n", pos.Volume)

if lotsToClose >= pos.Volume {
    fmt.Printf("❌ Cannot close %.2f - only %.2f available\n", lotsToClose, pos.Volume)
    return
}

err = sugar.ClosePositionPartial(ticket, lotsToClose)
if err != nil {
    fmt.Printf("Close failed: %v\n", err)
    return
}

fmt.Printf("✅ Closed %.2f lots\n", lotsToClose)
fmt.Printf("   Remaining: %.2f lots\n", pos.Volume-lotsToClose)
```

---

### 9) Trailing partial close

```go
ticket := uint64(12345)
highestProfit := 0.0
partialClosed := false

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil { // position closed (or query failed)
        return
    }

    if pos.Profit > highestProfit {
        highestProfit = pos.Profit
    }

    // If profit drops 30% from highest, close half
    profitDrop := highestProfit - pos.Profit
    dropPercent := (profitDrop / highestProfit) * 100

    fmt.Printf("Profit: $%.2f (High: $%.2f, Drop: %.1f%%)\n",
        pos.Profit, highestProfit, dropPercent)

    if !partialClosed && dropPercent >= 30.0 && highestProfit > 50.0 {
        fmt.Printf("⚠️  Profit dropped 30%% from peak - scaling out\n")

        halfVolume := math.Floor(pos.Volume/2/0.01) * 0.01 // 0.01 volume step
        err := sugar.ClosePositionPartial(ticket, halfVolume)
        if err != nil {
            fmt.Printf("Failed: %v\n", err)
            continue
        }

        fmt.Printf("✅ Closed half position at $%.2f\n", pos.Profit)
        partialClosed = true
    }
}
```

---

### 10) Advanced partial close function

```go
func ScaleOutPosition(sugar *mt5.MT5Sugar, ticket uint64, closePercent float64, reason string) error {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        return fmt.Errorf("position not found: %w", err)
    }

    // ClosePositionPartial sends the volume as given - round it to the volume step yourself
    info, err := sugar.GetSymbolInfo(pos.Symbol)
    if err != nil {
        return err
    }
    closeVolume := math.Floor(pos.Volume*(closePercent/100.0)/info.VolumeStep) * info.VolumeStep

    // Validate close volume
    if closeVolume >= pos.Volume {
        return fmt.Errorf("close volume %.2f >= position volume %.2f",
            closeVolume, pos.Volume)
    }

    if closeVolume < info.VolumeMin {
        return fmt.Errorf("close volume too small: %.2f", closeVolume)
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      PARTIAL CLOSE                    ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Reason:    %s\n", reason)
    fmt.Printf("Position:  #%d\n", ticket)
    fmt.Printf("Current:   %.2f lots\n", pos.Volume)
    fmt.Printf("Closing:   %.2f lots (%.0f%%)\n", closeVolume, closePercent)
    fmt.Printf("Remaining: %.2f lots\n\n", pos.Volume-closeVolume)

    err = sugar.ClosePositionPartial(ticket, closeVolume)
    if err != nil {
        return fmt.Errorf("partial close failed: %w", err)
    }

    fmt.Println("✅ Partial close successful")
    return nil
}

// Usage:
ScaleOutPosition(sugar, 12345, 50.0, "First profit target reached")
```

---

## 🔗 Related Methods

**🍬 Other close methods:**

* `ClosePosition()` - Close entire position
* `CloseAllPositions()` - Close all positions

**🍬 Position management:**

* `ModifyPositionSLTP()` - Adjust SL/TP after partial close
* `GetPositionByTicket()` - Check remaining volume

---

## ⚠️ Common Pitfalls

### 1) Closing more than position size

```go
// ❌ WRONG - trying to close more than available
pos, _ := sugar.GetPositionByTicket(ticket)
sugar.ClosePositionPartial(ticket, pos.Volume+0.1) // Will fail!

// ✅ CORRECT - validate first
if closeVolume >= pos.Volume {
    fmt.Println("Cannot close more than position size")
    return
}
```

### 2) Not updating SL after partial close

```go
// ❌ WRONG - SL stays same after reducing position
sugar.ClosePositionPartial(ticket, halfVolume)
// Risk is now different but SL unchanged!

// ✅ CORRECT - adjust SL for new position size
sugar.ClosePositionPartial(ticket, halfVolume)
sugar.ModifyPositionSLTP(ticket, newSL, tp)
```

### 3) Forgetting minimum volume requirements

```go
// ❌ WRONG - might be below broker minimum
sugar.ClosePositionPartial(ticket, 0.001) // Too small!

// ✅ CORRECT - check minimum volume and step (the method does not round)
info, _ := sugar.GetSymbolInfo(symbol)
closeVolume = math.Floor(closeVolume/info.VolumeStep) * info.VolumeStep
if closeVolume < info.VolumeMin {
    fmt.Println("Volume too small")
    return
}
```

---

## 💎 Pro Tips

1. **Scale out gradually** - Don't close too much at once

2. **Move SL to breakeven** - After taking partial profits

3. **Let winners run** - Keep part of position for big moves

4. **3-target strategy** - Close 1/3 at each profit level

5. **Check minimum volume** - Ensure remaining position meets broker minimums

---

**See also:** [`ClosePosition.md`](ClosePosition.md), [`ModifyPositionSLTP.md`](ModifyPositionSLTP.md), [`GetPositionByTicket.md`](GetPositionByTicket.md)
