# 🔍 Get Position By Ticket (`GetPositionByTicket`)

> **Sugar method:** Retrieves detailed information about a specific position by ticket number.

**API Information:**

* **Method:** `sugar.GetPositionByTicket(ticket)`
* **Timeout:** 5 seconds (it calls `GetOpenPositions()` and searches the result)
* **Returns:** `*pb.PositionInfo` structure

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetPositionByTicket(ticket uint64) (*pb.PositionInfo, error)
```

`pb.PositionInfo` is the generated protobuf message (`pb "github.com/MetaRPC/GoMT5/package"`). Its fields:

| Field | Type | Meaning |
|-------|------|---------|
| `Index` | `uint32` | Position index |
| `Ticket` | `uint64` | Position ticket |
| `Symbol` | `string` | Trading symbol |
| `Type` | `pb.BMT5_ENUM_POSITION_TYPE` | `pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY` (0) or `..._SELL` (1) |
| `Volume` | `float64` | Position volume (lots) |
| `PriceOpen` | `float64` | Entry price |
| `PriceCurrent` | `float64` | Current price |
| `StopLoss` | `float64` | Stop Loss price (0 if none) |
| `TakeProfit` | `float64` | Take Profit price (0 if none) |
| `Profit` | `float64` | Current profit/loss |
| `Swap` | `float64` | Swap/rollover |
| `PositionCommission` | `float64` | Commission |
| `OpenTime` | `*timestamppb.Timestamp` | Open time, broker server time |
| `OpenTimeUtc` | `*timestamppb.Timestamp` | The same moment in UTC |
| `LastUpdateTime` / `LastUpdateTimeUtc` | `*timestamppb.Timestamp` | Last change |
| `MagicNumber` | `int64` | Magic number |
| `Identifier` | `int64` | Position identifier |
| `Reason` | `pb.BMT5_ENUM_POSITION_REASON` | Why the position was opened |
| `Comment`, `ExternalId` | `string` | Comment, external id |
| `AccountLogin` | `int64` | Account login |

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `ticket` | `uint64` | Position ticket number |

| Output | Type | Description |
|--------|------|-------------|
| `*pb.PositionInfo` | struct | Position details |
| `error` | `error` | `"position with ticket N not found"` if no open position has that ticket, or the `GetOpenPositions` error |

---

## 💬 Just the Essentials

* **What it is:** Get all information about a specific position.
* **Why you need it:** Check profit, check if still open, get entry price, check SL/TP.
* **Sanity check:** Returns an error (never `nil, nil`) if the position doesn't exist (closed or never existed). The pointer is `nil` whenever `err != nil`, so check the error before using it. Only open positions are searched - a pending order ticket is "not found".

---

## 🎯 When to Use

✅ **Check profit** - Monitor specific position P/L

✅ **Verify existence** - Check if position still open

✅ **Get details** - Before modifying or closing

✅ **Track position** - Monitor progress of specific trade

---

## 🔗 Usage Examples

### 1) Basic usage - get position details

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

fmt.Printf("Position #%d:\n", pos.Ticket)
fmt.Printf("  Symbol:   %s\n", pos.Symbol)
fmt.Printf("  Type:     %s\n", pos.Type) // BMT5_POSITION_TYPE_BUY / BMT5_POSITION_TYPE_SELL
fmt.Printf("  Volume:   %.2f lots\n", pos.Volume)
fmt.Printf("  Entry:    %.5f\n", pos.PriceOpen)
fmt.Printf("  Current:  %.5f\n", pos.PriceCurrent)
fmt.Printf("  SL:       %.5f\n", pos.StopLoss)
fmt.Printf("  TP:       %.5f\n", pos.TakeProfit)
fmt.Printf("  Profit:   $%.2f\n", pos.Profit)
```

---

### 2) Check if position still exists

```go
ticket := uint64(12345)

_, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("❌ Position #%d no longer exists (closed or SL/TP hit)\n", ticket)
    return
}

fmt.Printf("✅ Position #%d is still open\n", ticket)
```

---

### 3) Monitor position until target

```go
ticket := uint64(12345)
targetProfit := 100.0

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

fmt.Printf("Monitoring position #%d until $%.2f profit\n", ticket, targetProfit)

for range ticker.C {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        fmt.Println("Position closed")
        return
    }

    fmt.Printf("Profit: $%.2f / $%.2f (%.1f%%)\n",
        pos.Profit, targetProfit, (pos.Profit/targetProfit)*100)

    if pos.Profit >= targetProfit {
        fmt.Printf("🎯 Target reached! Final profit: $%.2f\n", pos.Profit)
        return
    }
}
```

---

### 4) Calculate pips profit

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

// Calculate pips
var pips float64
if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
    pips = (pos.PriceCurrent - pos.PriceOpen) / 0.00001
} else {
    pips = (pos.PriceOpen - pos.PriceCurrent) / 0.00001
}

fmt.Printf("Position #%d:\n", pos.Ticket)
fmt.Printf("  Type:        %s\n", pos.Type)
fmt.Printf("  Entry:       %.5f\n", pos.PriceOpen)
fmt.Printf("  Current:     %.5f\n", pos.PriceCurrent)
fmt.Printf("  Pips:        %.1f\n", pips)
fmt.Printf("  Dollar P/L:  $%.2f\n", pos.Profit)
```

---

### 5) Check time in trade

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

// OpenTimeUtc for durations; OpenTime is the broker server's wall clock
duration := time.Since(pos.OpenTimeUtc.AsTime())

fmt.Printf("Position #%d opened at %s (broker time)\n",
    pos.Ticket, pos.OpenTime.AsTime().Format("2006-01-02 15:04:05"))
fmt.Printf("Time in trade: %v\n", duration.Round(time.Minute))
fmt.Printf("Current P/L: $%.2f\n", pos.Profit)

if duration > 24*time.Hour {
    fmt.Println("⚠️  Position open for more than 24 hours")
}
```

---

### 6) Verify SL/TP are set

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

fmt.Printf("Position #%d protection check:\n", pos.Ticket)

if pos.StopLoss == 0 {
    fmt.Println("  ❌ No Stop Loss - DANGEROUS!")
} else {
    fmt.Printf("  ✅ Stop Loss: %.5f\n", pos.StopLoss)
}

if pos.TakeProfit == 0 {
    fmt.Println("  ⚠️  No Take Profit")
} else {
    fmt.Printf("  ✅ Take Profit: %.5f\n", pos.TakeProfit)
}
```

---

### 7) Display position summary

```go
func ShowPositionSummary(sugar *mt5.MT5Sugar, ticket uint64) {
    pos, err := sugar.GetPositionByTicket(ticket)
    if err != nil {
        fmt.Printf("Position #%d not found\n", ticket)
        return
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║       POSITION SUMMARY                ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Ticket:      #%d\n", pos.Ticket)
    fmt.Printf("Symbol:      %s\n", pos.Symbol)
    fmt.Printf("Type:        %s\n", pos.Type)
    fmt.Printf("Volume:      %.2f lots\n", pos.Volume)
    fmt.Println()
    fmt.Printf("Entry:       %.5f\n", pos.PriceOpen)
    fmt.Printf("Current:     %.5f\n", pos.PriceCurrent)
    fmt.Printf("Stop Loss:   %.5f\n", pos.StopLoss)
    fmt.Printf("Take Profit: %.5f\n", pos.TakeProfit)
    fmt.Println()

    // Calculate pips
    var pips float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        pips = (pos.PriceCurrent - pos.PriceOpen) / 0.00001
    } else {
        pips = (pos.PriceOpen - pos.PriceCurrent) / 0.00001
    }

    fmt.Printf("Pips:        %.1f\n", pips)
    fmt.Printf("Profit:      $%.2f\n", pos.Profit)
    fmt.Printf("Commission:  $%.2f\n", pos.PositionCommission)
    fmt.Printf("Swap:        $%.2f\n", pos.Swap)
    fmt.Printf("Net:         $%.2f\n", pos.Profit+pos.PositionCommission+pos.Swap)
    fmt.Println()
    fmt.Printf("Opened:      %s (broker time)\n", pos.OpenTime.AsTime().Format("2006-01-02 15:04:05"))
    fmt.Printf("Duration:    %v\n", time.Since(pos.OpenTimeUtc.AsTime()).Round(time.Minute))
}

// Usage:
ShowPositionSummary(sugar, 12345)
```

---

### 8) Distance to SL/TP

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

info, _ := sugar.GetSymbolInfo(pos.Symbol)

fmt.Printf("Position #%d distances:\n", pos.Ticket)

if pos.StopLoss > 0 {
    var slDistance float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        slDistance = (pos.PriceCurrent - pos.StopLoss) / info.Point
    } else {
        slDistance = (pos.StopLoss - pos.PriceCurrent) / info.Point
    }
    fmt.Printf("  Distance to SL: %.0f points\n", slDistance)
}

if pos.TakeProfit > 0 {
    var tpDistance float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        tpDistance = (pos.TakeProfit - pos.PriceCurrent) / info.Point
    } else {
        tpDistance = (pos.PriceCurrent - pos.TakeProfit) / info.Point
    }
    fmt.Printf("  Distance to TP: %.0f points\n", tpDistance)
}
```

---

### 9) Position risk/reward ratio

```go
ticket := uint64(12345)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

if pos.StopLoss > 0 && pos.TakeProfit > 0 {
    var riskPips, rewardPips float64

    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        riskPips = (pos.PriceOpen - pos.StopLoss) / 0.00001
        rewardPips = (pos.TakeProfit - pos.PriceOpen) / 0.00001
    } else {
        riskPips = (pos.StopLoss - pos.PriceOpen) / 0.00001
        rewardPips = (pos.PriceOpen - pos.TakeProfit) / 0.00001
    }

    rr := rewardPips / riskPips

    fmt.Printf("Position #%d Risk/Reward:\n", pos.Ticket)
    fmt.Printf("  Risk:   %.0f pips\n", riskPips)
    fmt.Printf("  Reward: %.0f pips\n", rewardPips)
    fmt.Printf("  R:R:    1:%.1f\n", rr)

    if rr < 2.0 {
        fmt.Println("  ⚠️  R:R less than 1:2")
    } else {
        fmt.Println("  ✅ Good R:R ratio")
    }
}
```

---

### 10) Advanced position monitor

```go
func MonitorPosition(sugar *mt5.MT5Sugar, ticket uint64, duration time.Duration) {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║       POSITION MONITOR                ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    timeout := time.After(duration)

    highestProfit := 0.0
    lowestProfit := 0.0

    for {
        select {
        case <-timeout:
            fmt.Println("\n⏰ Monitoring ended")
            return

        case <-ticker.C:
            pos, err := sugar.GetPositionByTicket(ticket)
            if err != nil {
                fmt.Println("\n❌ Position closed")
                return
            }

            // Track high/low
            if pos.Profit > highestProfit {
                highestProfit = pos.Profit
            }
            if pos.Profit < lowestProfit {
                lowestProfit = pos.Profit
            }

            // Calculate pips
            var pips float64
            if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
                pips = (pos.PriceCurrent - pos.PriceOpen) / 0.00001
            } else {
                pips = (pos.PriceOpen - pos.PriceCurrent) / 0.00001
            }

            // Display update
            fmt.Printf("[%s] ", time.Now().Format("15:04:05"))
            fmt.Printf("%.5f | ", pos.PriceCurrent)
            fmt.Printf("%.1f pips | ", pips)
            fmt.Printf("$%.2f ", pos.Profit)

            if pos.Profit > 0 {
                fmt.Printf("📈")
            } else {
                fmt.Printf("📉")
            }

            fmt.Printf(" (High: $%.2f, Low: $%.2f)\n", highestProfit, lowestProfit)
        }
    }
}

// Usage: Monitor for 30 minutes
MonitorPosition(sugar, 12345, 30*time.Minute)
```

---

## 🔗 Related Methods

**🍬 Other position info:**

* `GetOpenPositions()` - Get all positions
* `GetPositionsBySymbol()` - Get positions for specific symbol

**🍬 Position management:**

* `ModifyPositionSLTP()` - Modify position after getting info
* `ClosePosition()` - Close position

---

## ⚠️ Common Pitfalls

### 1) Not checking for errors

```go
// ❌ WRONG - position might not exist
pos, _ := sugar.GetPositionByTicket(ticket)
fmt.Printf("Profit: $%.2f\n", pos.Profit) // nil pointer panic when not found!

// ✅ CORRECT - check error
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Println("Position not found")
    return
}
```

### 2) Using wrong ticket number

```go
// ❌ WRONG - using position from hours ago
oldTicket := uint64(12345)
pos, _ := sugar.GetPositionByTicket(oldTicket) // Might be closed!

// ✅ CORRECT - use fresh ticket from recent open
// (BuyMarket returns the order ticket; MT5 uses it as the position ticket of a market order)
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)
pos, err := sugar.GetPositionByTicket(ticket)
```

### 3) Assuming price direction

```go
// ❌ WRONG - assuming BUY
pips := (pos.PriceCurrent - pos.PriceOpen) / 0.00001

// ✅ CORRECT - check position type
var pips float64
if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
    pips = (pos.PriceCurrent - pos.PriceOpen) / 0.00001
} else {
    pips = (pos.PriceOpen - pos.PriceCurrent) / 0.00001
}
```

---

## 💎 Pro Tips

1. **Check errors** - Position might be closed

2. **Cache ticket** - Store ticket when opening position

3. **Type checking** - Always check BUY vs SELL for calculations

4. **Monitor changes** - Compare values over time

5. **Use for decisions** - Check profit before modifying/closing

---

**See also:** [`GetOpenPositions.md`](GetOpenPositions.md), [`ModifyPositionSLTP.md`](ModifyPositionSLTP.md), [`ClosePosition.md`](ClosePosition.md)
