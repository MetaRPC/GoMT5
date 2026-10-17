# 🟢📈 Buy Stop Order (`BuyStop`)

> **Sugar method:** Places pending BUY order that executes when price **breaks ABOVE** specified level.

**API Information:**

* **Method:** `sugar.BuyStop(symbol, volume, price)`
* **Timeout:** 10 seconds
* **Returns:** Pending order ticket number (`result.Order`)
* **Request:** `service.PlaceOrder(ctx, &pb.OrderSendRequest{Symbol, Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP, Volume, Price: &price})` -
  no SL/TP, expiration, comment or magic number; price and volume are sent exactly as given (no rounding, no check
  against the current price - a wrong price is rejected by the trade server)
* **Success:** only return code `10009` (TRADE_RETCODE_DONE); any other code gives the error
  `order rejected, code: N, comment: ...`, a failed call gives `BuyStop failed: ...`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) BuyStop(symbol string, volume, price float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `price` | `float64` | Entry price (must be **ABOVE** current ASK) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Pending order ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** Pending BUY order that activates when price **breaks ABOVE** resistance.
* **Why you need it:** Buy on breakouts - when price confirms upward momentum.
* **Sanity check:** Entry price must be **ABOVE** current ASK. Order activates on breakout.

---

## 📊 BUY STOP Logic

```
Current ASK: 1.08500
BUY STOP:    1.09000 (500 points ABOVE current price)

When price rises to 1.09000 → Order executes automatically
```

**Use case:** Breakout trading / Buy on resistance break / Momentum trading

---

## 🎯 When to Use

✅ **Breakout trading** - Buy when price breaks resistance

✅ **Momentum trading** - Enter when uptrend confirmed

✅ **Range breakouts** - Buy when price exits consolidation upward

✅ **Follow strong trends** - Join established upward momentum

❌ **NOT for pullbacks** - Use `BuyLimit()` for buying dips

---

## 🔗 Usage Examples

### 1) Basic usage - buy on breakout

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)
resistanceLevel := 1.09000 // Known resistance

// Place buy stop above resistance
ticket, err := sugar.BuyStop(symbol, 0.1, resistanceLevel)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY STOP order placed\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  %.5f (current: %.5f)\n", resistanceLevel, currentAsk)
fmt.Printf("   Will execute when price breaks %.0f points higher\n",
    (resistanceLevel-currentAsk)/0.00001)
```

---

### 2) Range breakout trading

```go
symbol := "EURUSD"

// Price consolidating between 1.08000 and 1.08500
rangeHigh := 1.08500
rangeHigh += 0.00010 // Place stop 10 points above range high

ticket, _ := sugar.BuyStop(symbol, 0.1, rangeHigh)

fmt.Printf("BUY STOP set for range breakout\n")
fmt.Printf("Range high: %.5f\n", 1.08500)
fmt.Printf("Entry:      %.5f (10 points above)\n", rangeHigh)
fmt.Printf("Strategy:   Buy if price breaks out of range\n")
```

---

### 3) Momentum breakout with multiple levels

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)
volume := 0.05

// Set buy stops at different breakout levels
breakoutLevels := []float64{
    1.08500, // First resistance
    1.09000, // Second resistance
    1.09500, // Third resistance
}

fmt.Println("Placing multiple BUY STOP orders:")

for i, level := range breakoutLevels {
    // Only place stops above current price
    if level <= currentAsk {
        fmt.Printf("Level %d: %.5f - SKIPPED (below current price)\n", i+1, level)
        continue
    }

    ticket, err := sugar.BuyStop(symbol, volume, level)
    if err != nil {
        fmt.Printf("Level %d failed: %v\n", i+1, err)
        continue
    }

    pointsAway := (level - currentAsk) / 0.00001
    fmt.Printf("✅ Level %d: %.5f (%.0f points away) - Ticket #%d\n",
        i+1, level, pointsAway, ticket)
}
```

---

### 4) Buy stop with validation

```go
symbol := "EURUSD"
volume := 0.1
breakoutPrice := 1.09000

// Get current price
currentAsk, _ := sugar.GetAsk(symbol)

// Validate entry price is above current
if breakoutPrice <= currentAsk {
    fmt.Printf("❌ Error: Entry %.5f must be above ASK %.5f\n",
        breakoutPrice, currentAsk)
    return
}

// Check minimum stop level
info, _ := sugar.GetSymbolInfo(symbol)
stopLevel, _ := sugar.GetMinStopLevel(symbol)
minDistance := float64(stopLevel) * info.Point

if (breakoutPrice - currentAsk) < minDistance {
    fmt.Printf("❌ Entry too close (min: %.5f)\n", minDistance)
    return
}

// Place order
ticket, _ := sugar.BuyStop(symbol, volume, breakoutPrice)
fmt.Printf("✅ BUY STOP placed at %.5f\n", breakoutPrice)
```

---

### 5) Buy stop with SL/TP (ATR-based)

```go
symbol := "EURUSD"
breakoutPrice := 1.09000

// Calculate SL/TP based on ATR or fixed points
// For breakout: SL below breakout level, TP 2x the distance
info, _ := sugar.GetSymbolInfo(symbol)
slDistance := 30.0 // 30 points SL
tpDistance := 60.0 // 60 points TP (2:1 R:R)

sl := breakoutPrice - (slDistance * info.Point)
tp := breakoutPrice + (tpDistance * info.Point)

// BuyStop sends no SL/TP and Sugar has no BuyStopWithSLTP - place the
// BUY STOP with SL/TP in one request through the Service layer
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := sugar.GetService().PlaceOrder(ctx, &pb.OrderSendRequest{
    Symbol:     symbol,
    Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP,
    Volume:     0.1,
    Price:      &breakoutPrice,
    StopLoss:   &sl,
    TakeProfit: &tp,
})
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}
if result.ReturnedCode != 10009 {
    fmt.Printf("Order rejected: code %d, %s\n", result.ReturnedCode, result.Comment)
    return
}

fmt.Printf("BUY STOP with SL/TP: Ticket #%d\n", result.Order)
fmt.Printf("  Entry: %.5f (breakout)\n", breakoutPrice)
fmt.Printf("  SL:    %.5f (-30 points)\n", sl)
fmt.Printf("  TP:    %.5f (+60 points)\n", tp)
fmt.Printf("  R:R:   1:2\n")
```

---

### 6) Trend continuation buy stop

```go
symbol := "EURUSD"

// Strong uptrend - previous high at 1.09000
previousHigh := 1.09000
entryPrice := previousHigh + 0.00010 // 10 points above previous high

ticket, _ := sugar.BuyStop(symbol, 0.1, entryPrice)

fmt.Printf("Trend continuation BUY STOP\n")
fmt.Printf("Previous high: %.5f\n", previousHigh)
fmt.Printf("Entry:         %.5f (+10 points)\n", entryPrice)
fmt.Printf("Strategy:      Buy if uptrend continues\n")
```

---

### 7) News breakout strategy

```go
// Before major news release - price at 1.08500
symbol := "EURUSD"
currentPrice := 1.08500

// Place buy stop for upside breakout
buyStopPrice := currentPrice + 0.00030 // 30 points above
buyTicket, _ := sugar.BuyStop(symbol, 0.1, buyStopPrice)

// Also place sell stop for downside breakout
sellStopPrice := currentPrice - 0.00030 // 30 points below
sellTicket, _ := sugar.SellStop(symbol, 0.1, sellStopPrice)

fmt.Printf("News breakout strategy:\n")
fmt.Printf("  Current:   %.5f\n", currentPrice)
fmt.Printf("  BUY STOP:  %.5f (ticket #%d)\n", buyStopPrice, buyTicket)
fmt.Printf("  SELL STOP: %.5f (ticket #%d)\n", sellStopPrice, sellTicket)
fmt.Printf("  Strategy:  OCO - whichever direction breaks first\n")
fmt.Printf("  Remember:  Cancel other order when one fills!\n")
// Cancel = sugar.GetService().CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: otherTicket})
```

---

### 8) Trailing buy stop (manual adjustment)

```go
symbol := "EURUSD"
resistanceLevel := 1.09000

// Initial buy stop
ticket, _ := sugar.BuyStop(symbol, 0.1, resistanceLevel)
fmt.Printf("Initial BUY STOP: %.5f - Ticket #%d\n", resistanceLevel, ticket)

// Monitor price and adjust stop level
ticker := time.NewTicker(1 * time.Minute)
defer ticker.Stop()

for i := 0; i < 10; i++ {
    <-ticker.C

    currentAsk, _ := sugar.GetAsk(symbol)

    // If price moves up, move buy stop up (trailing)
    newStopLevel := currentAsk + 0.00050 // Always 50 points above

    if newStopLevel > resistanceLevel {
        // Modify pending order price
        fmt.Printf("Adjusting BUY STOP from %.5f to %.5f\n",
            resistanceLevel, newStopLevel)

        // Sugar has no ModifyOrder - change the pending order price through the Service
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        result, err := sugar.GetService().ModifyOrder(ctx, &pb.OrderModifyRequest{
            Ticket: ticket,
            Price:  &newStopLevel,
        })
        cancel()
        if err != nil || result.ReturnedCode != 10009 {
            fmt.Println("   Modify failed - keeping old level")
            continue
        }

        resistanceLevel = newStopLevel
    } else {
        fmt.Printf("Current: %.5f - Stop remains at %.5f\n",
            currentAsk, resistanceLevel)
    }
}
```

---

### 9) Pattern breakout - double top

```go
symbol := "EURUSD"

// Double top pattern - two peaks at 1.09000
doubleTopLevel := 1.09000

// Place buy stop above double top (invalidation = strong breakout)
entryPrice := doubleTopLevel + 0.00015 // 15 points above

ticket, _ := sugar.BuyStop(symbol, 0.1, entryPrice)

fmt.Printf("Double top breakout strategy:\n")
fmt.Printf("Double top:  %.5f\n", doubleTopLevel)
fmt.Printf("BUY STOP:    %.5f (+15 points)\n", entryPrice)
fmt.Printf("Logic:       If double top breaks = strong bullish signal\n")
```

---

### 10) Advanced breakout with volume confirmation

```go
func PlaceBuyStopWithMonitoring(sugar *mt5.MT5Sugar, symbol string, breakoutLevel float64) {
    currentAsk, _ := sugar.GetAsk(symbol)

    fmt.Printf("╔═══════════════════════════════════════╗\n")
    fmt.Printf("║       BREAKOUT BUY STOP               ║\n")
    fmt.Printf("╚═══════════════════════════════════════╝\n")
    fmt.Printf("Symbol:       %s\n", symbol)
    fmt.Printf("Current ASK:  %.5f\n", currentAsk)
    fmt.Printf("Breakout:     %.5f\n", breakoutLevel)
    fmt.Printf("Distance:     %.0f points\n", (breakoutLevel-currentAsk)/0.00001)

    // Place buy stop
    ticket, err := sugar.BuyStop(symbol, 0.1, breakoutLevel)
    if err != nil {
        fmt.Printf("❌ Order failed: %v\n", err)
        return
    }

    fmt.Printf("✅ Order placed: #%d\n", ticket)
    fmt.Println("\nMonitoring for breakout...")

    // Monitor until filled or timeout
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    timeout := time.After(2 * time.Hour)

    for {
        select {
        case <-timeout:
            fmt.Println("⏰ Timeout - breakout didn't occur in 2 hours")
            return

        case <-ticker.C:
            ask, _ := sugar.GetAsk(symbol)
            distance := breakoutLevel - ask
            distancePoints := distance / 0.00001

            // Check if filled
            _, err := sugar.GetPositionByTicket(ticket)
            if err == nil {
                fmt.Printf("\n🚀 BREAKOUT! Order filled at %.5f!\n", ask)
                return
            }

            if distancePoints < 10 {
                fmt.Printf("⚠️  ASK: %.5f (%.0f points from breakout - CLOSE!)\n",
                    ask, distancePoints)
            } else {
                fmt.Printf("ASK: %.5f (%.0f points from breakout)\n",
                    ask, distancePoints)
            }
        }
    }
}

// Usage:
PlaceBuyStopWithMonitoring(sugar, "EURUSD", 1.09000)
```

---

## 🔗 Related Methods

**🍬 Other stop orders:**

* `SellStop()` - SELL when price breaks DOWN
* BUY STOP with SL/TP - no Sugar method; use `GetService().PlaceOrder()` with `StopLoss`/`TakeProfit` (example 5)

**🍬 Limit orders (for pullbacks):**

* `BuyLimit()` - BUY when price drops to support

**🍬 Market orders:**

* `BuyMarket()` - BUY immediately at current price

---

## ⚠️ Common Pitfalls

### 1) Setting entry price below current price

```go
// ❌ WRONG - Buy Stop must be ABOVE current price
currentAsk := 1.08500
sugar.BuyStop("EURUSD", 0.1, 1.08000) // BELOW! Will be rejected!

// ✅ CORRECT - Above current price
sugar.BuyStop("EURUSD", 0.1, 1.09000) // Above ASK
```

### 2) Confusing BuyStop with BuyLimit

```go
// BUY STOP = wait for price to RISE then buy (breakout)
sugar.BuyStop("EURUSD", 0.1, 1.09000) // Buy when rises to 1.09000

// BUY LIMIT = wait for price to DROP then buy (pullback)
sugar.BuyLimit("EURUSD", 0.1, 1.08000) // Buy when drops to 1.08000
```

### 3) Chasing breakouts too close to current price

```go
// ❌ WRONG - too close to current price (might trigger immediately)
currentAsk := 1.08495
sugar.BuyStop("EURUSD", 0.1, 1.08500) // Only 5 points away!

// ✅ CORRECT - proper distance for confirmed breakout
sugar.BuyStop("EURUSD", 0.1, 1.09000) // Clear breakout level
```

---

## 💎 Pro Tips

1. **BUY STOP = buy on strength** - Wait for breakout confirmation

2. **Set above resistance** - Use previous highs, round numbers

3. **Add buffer** - Place a buffer (e.g. 10-20 points) above resistance for confirmation

4. **Use SL below breakout** - If price reverses, exit quickly

5. **Remember OCO** - If using both buy/sell stops, cancel one when other fills

---

**See also:** [`SellStop.md`](SellStop.md), [`BuyLimit.md`](BuyLimit.md)
