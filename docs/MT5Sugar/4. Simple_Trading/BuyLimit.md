# 🟢📉 Buy Limit Order (`BuyLimit`)

> **Sugar method:** Places pending BUY order that executes when price **drops** to specified level.

**API Information:**

* **Method:** `sugar.BuyLimit(symbol, volume, price)`
* **Timeout:** 30 seconds
* **Returns:** Pending order ticket number (`result.Order`)
* **Request:** `service.PlaceOrder(ctx, &pb.OrderSendRequest{Symbol, Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT, Volume, Price: &price})` -
  no SL/TP, expiration, comment or magic number; price and volume are sent exactly as given (no rounding, no check
  against the current price - a wrong price is rejected by the trade server)
* **Success:** only return code `10009` (TRADE_RETCODE_DONE); any other code gives the error
  `order rejected, code: N, comment: ...`, a failed call gives `BuyLimit failed: ...`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) BuyLimit(symbol string, volume, price float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `price` | `float64` | Entry price (must be **BELOW** current ASK) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Pending order ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** Pending BUY order that activates when price **drops** to your target.
* **Why you need it:** Buy at a better price (support level) instead of current market.
* **Sanity check:** Entry price must be **BELOW** current ASK. Order waits until price reaches your level.

---

## 📊 BUY LIMIT Logic

```
Current ASK: 1.08500
BUY LIMIT:   1.08000 (500 points BELOW current price)

When price drops to 1.08000 → Order executes automatically
```

**Use case:** Buy at support / Buy the dip / Enter on pullback

---

## 🎯 When to Use

✅ **Support levels** - Buy when price reaches support

✅ **Pullbacks** - Enter on price retracement

✅ **Buy the dip** - Wait for price to drop before buying

✅ **Better entry** - Get better price than current market

❌ **NOT for breakouts** - Use `BuyStop()` for breakouts

---

## 🔗 Usage Examples

### 1) Basic usage - buy at support

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)
supportLevel := currentAsk - 0.00050 // 50 points below (5-digit symbol)

ticket, err := sugar.BuyLimit(symbol, 0.1, supportLevel)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY LIMIT order placed\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  %.5f (current: %.5f)\n", supportLevel, currentAsk)
fmt.Printf("   Will execute when price drops %.0f points\n",
    (currentAsk-supportLevel)/0.00001)
```

---

### 2) Buy on pullback

```go
symbol := "EURUSD"

// Current uptrend, wait for pullback
currentAsk, _ := sugar.GetAsk(symbol)
pullbackLevel := currentAsk - 0.00030 // Wait for 30-point pullback

ticket, _ := sugar.BuyLimit(symbol, 0.1, pullbackLevel)

fmt.Printf("BUY LIMIT set for pullback entry\n")
fmt.Printf("Current: %.5f\n", currentAsk)
fmt.Printf("Entry:   %.5f (pullback level)\n", pullbackLevel)
```

---

### 3) Multiple buy limits at different levels

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)
volume := 0.1

// Set limits at different support levels
levels := []float64{
    currentAsk - 0.00020, // -20 points
    currentAsk - 0.00050, // -50 points
    currentAsk - 0.00100, // -100 points
}

fmt.Println("Placing multiple BUY LIMIT orders:")

for i, level := range levels {
    ticket, err := sugar.BuyLimit(symbol, volume, level)
    if err != nil {
        fmt.Printf("Level %d failed: %v\n", i+1, err)
        continue
    }

    pointsAway := (currentAsk - level) / 0.00001
    fmt.Printf("✅ Level %d: %.5f (%.0f points away) - Ticket #%d\n",
        i+1, level, pointsAway, ticket)
}

// Output (current ASK 1.08500):
// ✅ Level 1: 1.08480 (20 points away) - Ticket #12345
// ✅ Level 2: 1.08450 (50 points away) - Ticket #12346
// ✅ Level 3: 1.08400 (100 points away) - Ticket #12347
```

---

### 4) Buy limit with validation

```go
symbol := "EURUSD"
volume := 0.1

// Get current price
currentAsk, _ := sugar.GetAsk(symbol)

// Calculate entry price (support level)
info, _ := sugar.GetSymbolInfo(symbol)
entryPrice := currentAsk - (50 * info.Point) // 50 points below

// Validate entry price is below current
if entryPrice >= currentAsk {
    fmt.Printf("❌ Error: Entry %.5f must be below ASK %.5f\n",
        entryPrice, currentAsk)
    return
}

// Check minimum stop level
stopLevel, _ := sugar.GetMinStopLevel(symbol)
minDistance := float64(stopLevel) * info.Point

if (currentAsk - entryPrice) < minDistance {
    fmt.Printf("❌ Entry too close (min: %.5f)\n", minDistance)
    return
}

// Place order
ticket, _ := sugar.BuyLimit(symbol, volume, entryPrice)
fmt.Printf("✅ BUY LIMIT placed at %.5f\n", entryPrice)
```

---

### 5) Buy limit with SL/TP (using separate method)

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)
entryPrice := currentAsk - 0.00050 // 50 points below

// Calculate SL/TP from entry price (distances in points)
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", entryPrice, 30, 60)

// BuyLimit itself sends no SL/TP; the ModifyPosition* methods are meant for
// open positions. Put SL/TP into the pending order itself:
ticket, _ := sugar.BuyLimitWithSLTP(symbol, 0.1, entryPrice, sl, tp)

fmt.Printf("BUY LIMIT with SL/TP: Ticket #%d\n", ticket)
fmt.Printf("  Entry: %.5f\n", entryPrice)
fmt.Printf("  SL:    %.5f (-30 points)\n", sl)
fmt.Printf("  TP:    %.5f (+60 points)\n", tp)
```

---

### 6) Scale in with buy limits

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)
totalVolume := 0.3
entries := 3

volumePerEntry := totalVolume / float64(entries)

fmt.Printf("Scaling in with %d BUY LIMIT orders\n", entries)

for i := 1; i <= entries; i++ {
    // Each level 20 points apart
    entryPrice := currentAsk - (float64(i) * 0.00020)

    ticket, err := sugar.BuyLimit(symbol, volumePerEntry, entryPrice)
    if err != nil {
        continue
    }

    fmt.Printf("Entry %d: %.2f lots at %.5f - Ticket #%d\n",
        i, volumePerEntry, entryPrice, ticket)
}
```

---

### 7) Cancel unfilled buy limit

```go
symbol := "EURUSD"
entryPrice := 1.08000

ticket, _ := sugar.BuyLimit(symbol, 0.1, entryPrice)
fmt.Printf("BUY LIMIT placed: Ticket #%d\n", ticket)

// Wait 5 minutes
fmt.Println("Waiting 5 minutes...")
time.Sleep(5 * time.Minute)

// Check if filled (the position opened by the order has the order's ticket)
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    // Not filled - cancel it. Sugar has no DeleteOrder: delete the pending
    // order with the Service's CloseOrder (MT5 OrderClose deletes a pending order)
    fmt.Println("Order not filled - canceling...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // trade requests: 30 s like Sugar
    defer cancel()
    retCode, err := sugar.GetService().CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: ticket})
    if err != nil || retCode != 10009 {
        fmt.Printf("❌ Cancel failed: code %d, err %v\n", retCode, err)
    } else {
        fmt.Println("✅ Pending order deleted")
    }
} else {
    fmt.Printf("✅ Order filled at %.5f\n", pos.PriceOpen)
}
```

---

### 8) Buy limit at Fibonacci retracement

```go
symbol := "EURUSD"

// Previous swing high and low
swingHigh := 1.09000
swingLow := 1.08000
range_ := swingHigh - swingLow

// 38.2% Fibonacci retracement
fib382 := swingHigh - (range_ * 0.382)

// Place buy limit at Fib level
ticket, _ := sugar.BuyLimit(symbol, 0.1, fib382)

fmt.Printf("BUY LIMIT at Fibonacci 38.2%%\n")
fmt.Printf("Swing High: %.5f\n", swingHigh)
fmt.Printf("Swing Low:  %.5f\n", swingLow)
fmt.Printf("Entry:      %.5f (Fib 38.2%%)\n", fib382)
fmt.Printf("Ticket:     #%d\n", ticket)
```

---

### 9) Grid trading with buy limits

```go
symbol := "EURUSD"
currentAsk, _ := sugar.GetAsk(symbol)

gridLevels := 5
gridSpacing := 0.00020 // 20 points
volume := 0.01

fmt.Printf("Grid trading: %d BUY LIMIT orders\n", gridLevels)

for i := 1; i <= gridLevels; i++ {
    entryPrice := currentAsk - (float64(i) * gridSpacing)

    ticket, err := sugar.BuyLimit(symbol, volume, entryPrice)
    if err != nil {
        fmt.Printf("Grid %d failed: %v\n", i, err)
        continue
    }

    pointsAway := (currentAsk - entryPrice) / 0.00001
    fmt.Printf("Grid %d: %.5f (%.0f points) - #%d\n",
        i, entryPrice, pointsAway, ticket)
}
```

---

### 10) Advanced buy limit with monitoring

```go
func PlaceBuyLimitWithMonitoring(sugar *mt5.MT5Sugar, symbol string, pointsBelow float64) {
    // Get current price
    currentAsk, _ := sugar.GetAsk(symbol)
    info, _ := sugar.GetSymbolInfo(symbol)

    // Calculate entry
    entryPrice := currentAsk - (pointsBelow * info.Point)

    fmt.Printf("╔═══════════════════════════════════════╗\n")
    fmt.Printf("║       BUY LIMIT ORDER                 ║\n")
    fmt.Printf("╚═══════════════════════════════════════╝\n")
    fmt.Printf("Symbol:       %s\n", symbol)
    fmt.Printf("Current ASK:  %.5f\n", currentAsk)
    fmt.Printf("Entry Price:  %.5f\n", entryPrice)
    fmt.Printf("Distance:     %.0f points\n", pointsBelow)

    // Place order
    ticket, err := sugar.BuyLimit(symbol, 0.1, entryPrice)
    if err != nil {
        fmt.Printf("❌ Order failed: %v\n", err)
        return
    }

    fmt.Printf("✅ Order placed: #%d\n", ticket)
    fmt.Println("\nMonitoring price...")

    // Monitor until filled or timeout
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    timeout := time.After(30 * time.Minute)

    for {
        select {
        case <-timeout:
            fmt.Println("⏰ Timeout - order not filled in 30 minutes")
            return

        case <-ticker.C:
            ask, _ := sugar.GetAsk(symbol)
            distance := ask - entryPrice
            distancePoints := distance / info.Point

            // Check if filled
            _, err := sugar.GetPositionByTicket(ticket)
            if err == nil {
                fmt.Printf("\n✅ ORDER FILLED at %.5f!\n", ask)
                return
            }

            fmt.Printf("ASK: %.5f (%.0f points away from entry)\n",
                ask, distancePoints)
        }
    }
}

// Usage:
PlaceBuyLimitWithMonitoring(sugar, "EURUSD", 50)
```

---

## 🔗 Related Methods

**🍬 Other limit orders:**
* `SellLimit()` - SELL when price goes UP
* `BuyLimitWithSLTP()` - BUY LIMIT with SL/TP

**🍬 Stop orders (for breakouts):**
* `BuyStop()` - BUY when price breaks UP

**🍬 Market orders:**
* `BuyMarket()` - BUY immediately at current price

---

## ⚠️ Common Pitfalls

### 1) Setting entry price above current price

```go
// ❌ WRONG - Buy Limit must be BELOW current price
currentAsk := 1.08500
sugar.BuyLimit("EURUSD", 0.1, 1.09000) // ABOVE! Will be rejected!

// ✅ CORRECT - Below current price
sugar.BuyLimit("EURUSD", 0.1, 1.08000) // Below ASK
```

### 2) Confusing BuyLimit with BuyStop

```go
// BUY LIMIT = wait for price to DROP then buy
sugar.BuyLimit("EURUSD", 0.1, 1.08000) // Buy when drops to 1.08000

// BUY STOP = wait for price to RISE then buy (breakout)
sugar.BuyStop("EURUSD", 0.1, 1.09000) // Buy when rises to 1.09000
```

---

## 💎 Pro Tips

1. **BUY LIMIT = buy cheaper** - Wait for price to drop
2. **Set at support levels** - Use technical analysis
3. **Multiple levels** - Scale in at different prices
4. **Monitor fills** - Check if order was executed
5. **Use with SL/TP** - `BuyLimitWithSLTP()` puts SL/TP into the pending order

---

**See also:** [`BuyLimit WithSLTP.md`](../5.%20Trading_SLTP/BuyLimitWithSLTP.md), [`BuyStop.md`](BuyStop.md), [`SellLimit.md`](SellLimit.md)
