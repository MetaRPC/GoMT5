# 🔴📈 Sell Limit Order (`SellLimit`)

> **Sugar method:** Places pending SELL order that executes when price **rises** to specified level.

**API Information:**

* **Method:** `sugar.SellLimit(symbol, volume, price)`
* **Timeout:** 10 seconds
* **Returns:** Pending order ticket number (`result.Order`)
* **Request:** `service.PlaceOrder(ctx, &pb.OrderSendRequest{Symbol, Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_LIMIT, Volume, Price: &price})` -
  no SL/TP, expiration, comment or magic number; price and volume are sent exactly as given (no rounding, no check
  against the current price - a wrong price is rejected by the trade server)
* **Success:** only return code `10009` (TRADE_RETCODE_DONE); any other code gives the error
  `order rejected, code: N, comment: ...`, a failed call gives `SellLimit failed: ...`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) SellLimit(symbol string, volume, price float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `price` | `float64` | Entry price (must be **ABOVE** current BID) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Pending order ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** Pending SELL order that activates when price **rises** to your target.
* **Why you need it:** Sell at a better price (resistance level) instead of current market.
* **Sanity check:** Entry price must be **ABOVE** current BID. Order waits until price reaches your level.

---

## 📊 SELL LIMIT Logic

```
Current BID: 1.08500
SELL LIMIT:  1.09000 (500 points ABOVE current price)

When price rises to 1.09000 → Order executes automatically
```

**Use case:** Sell at resistance / Sell the rally / Enter on price rise

---

## 🎯 When to Use

✅ **Resistance levels** - Sell when price reaches resistance

✅ **Rally exhaustion** - Enter on price retracement up

✅ **Sell the rally** - Wait for price to rise before selling

✅ **Better entry** - Get better price than current market

❌ **NOT for breakdowns** - Use `SellStop()` for breakdowns

---

## 🔗 Usage Examples

### 1) Basic usage - sell at resistance

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
resistanceLevel := currentBid + 0.00050 // 50 points above (5-digit symbol)

ticket, err := sugar.SellLimit(symbol, 0.1, resistanceLevel)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ SELL LIMIT order placed\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  %.5f (current: %.5f)\n", resistanceLevel, currentBid)
fmt.Printf("   Will execute when price rises %.0f points\n",
    (resistanceLevel-currentBid)/0.00001)
```

---

### 2) Sell on rally exhaustion

```go
symbol := "EURUSD"

// Current downtrend, wait for rally to fade
currentBid, _ := sugar.GetBid(symbol)
rallyLevel := currentBid + 0.00030 // Wait for 30-point rally

ticket, _ := sugar.SellLimit(symbol, 0.1, rallyLevel)

fmt.Printf("SELL LIMIT set for rally exhaustion entry\n")
fmt.Printf("Current: %.5f\n", currentBid)
fmt.Printf("Entry:   %.5f (rally level)\n", rallyLevel)
```

---

### 3) Multiple sell limits at different levels

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
volume := 0.1

// Set limits at different resistance levels
levels := []float64{
    currentBid + 0.00020, // +20 points
    currentBid + 0.00050, // +50 points
    currentBid + 0.00100, // +100 points
}

fmt.Println("Placing multiple SELL LIMIT orders:")

for i, level := range levels {
    ticket, err := sugar.SellLimit(symbol, volume, level)
    if err != nil {
        fmt.Printf("Level %d failed: %v\n", i+1, err)
        continue
    }

    pointsAway := (level - currentBid) / 0.00001
    fmt.Printf("✅ Level %d: %.5f (%.0f points away) - Ticket #%d\n",
        i+1, level, pointsAway, ticket)
}

// Output (current BID 1.08500):
// ✅ Level 1: 1.08520 (20 points away) - Ticket #12345
// ✅ Level 2: 1.08550 (50 points away) - Ticket #12346
// ✅ Level 3: 1.08600 (100 points away) - Ticket #12347
```

---

### 4) Sell limit with validation

```go
symbol := "EURUSD"
volume := 0.1

// Get current price
currentBid, _ := sugar.GetBid(symbol)

// Calculate entry price (resistance level)
info, _ := sugar.GetSymbolInfo(symbol)
entryPrice := currentBid + (50 * info.Point) // 50 points above

// Validate entry price is above current
if entryPrice <= currentBid {
    fmt.Printf("❌ Error: Entry %.5f must be above BID %.5f\n",
        entryPrice, currentBid)
    return
}

// Check minimum stop level
stopLevel, _ := sugar.GetMinStopLevel(symbol)
minDistance := float64(stopLevel) * info.Point

if (entryPrice - currentBid) < minDistance {
    fmt.Printf("❌ Entry too close (min: %.5f)\n", minDistance)
    return
}

// Place order
ticket, _ := sugar.SellLimit(symbol, volume, entryPrice)
fmt.Printf("✅ SELL LIMIT placed at %.5f\n", entryPrice)
```

---

### 5) Sell limit with SL/TP (using separate method)

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
entryPrice := currentBid + 0.00050 // 50 points above

// Calculate SL/TP from entry price (distances in points)
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", entryPrice, 30, 60)

// SellLimit itself sends no SL/TP; the ModifyPosition* methods are meant for
// open positions. Put SL/TP into the pending order itself:
ticket, _ := sugar.SellLimitWithSLTP(symbol, 0.1, entryPrice, sl, tp)

fmt.Printf("SELL LIMIT with SL/TP: Ticket #%d\n", ticket)
fmt.Printf("  Entry: %.5f\n", entryPrice)
fmt.Printf("  SL:    %.5f (+30 points)\n", sl)
fmt.Printf("  TP:    %.5f (-60 points)\n", tp)
```

---

### 6) Scale in with sell limits

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
totalVolume := 0.3
entries := 3

volumePerEntry := totalVolume / float64(entries)

fmt.Printf("Scaling in with %d SELL LIMIT orders\n", entries)

for i := 1; i <= entries; i++ {
    // Each level 20 points apart
    entryPrice := currentBid + (float64(i) * 0.00020)

    ticket, err := sugar.SellLimit(symbol, volumePerEntry, entryPrice)
    if err != nil {
        continue
    }

    fmt.Printf("Entry %d: %.2f lots at %.5f - Ticket #%d\n",
        i, volumePerEntry, entryPrice, ticket)
}
```

---

### 7) Cancel unfilled sell limit

```go
symbol := "EURUSD"
entryPrice := 1.09000

ticket, _ := sugar.SellLimit(symbol, 0.1, entryPrice)
fmt.Printf("SELL LIMIT placed: Ticket #%d\n", ticket)

// Wait 5 minutes
fmt.Println("Waiting 5 minutes...")
time.Sleep(5 * time.Minute)

// Check if filled (the position opened by the order has the order's ticket)
pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    // Not filled - cancel it. Sugar has no DeleteOrder: delete the pending
    // order with the Service's CloseOrder (MT5 OrderClose deletes a pending order)
    fmt.Println("Order not filled - canceling...")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

### 8) Sell limit at Fibonacci retracement

```go
symbol := "EURUSD"

// Previous swing high and low
swingHigh := 1.09000
swingLow := 1.08000
range_ := swingHigh - swingLow

// 61.8% Fibonacci retracement (common resistance in downtrend)
fib618 := swingLow + (range_ * 0.618)

// Place sell limit at Fib level
ticket, _ := sugar.SellLimit(symbol, 0.1, fib618)

fmt.Printf("SELL LIMIT at Fibonacci 61.8%%\n")
fmt.Printf("Swing High: %.5f\n", swingHigh)
fmt.Printf("Swing Low:  %.5f\n", swingLow)
fmt.Printf("Entry:      %.5f (Fib 61.8%%)\n", fib618)
fmt.Printf("Ticket:     #%d\n", ticket)
```

---

### 9) Grid trading with sell limits

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)

gridLevels := 5
gridSpacing := 0.00020 // 20 points
volume := 0.01

fmt.Printf("Grid trading: %d SELL LIMIT orders\n", gridLevels)

for i := 1; i <= gridLevels; i++ {
    entryPrice := currentBid + (float64(i) * gridSpacing)

    ticket, err := sugar.SellLimit(symbol, volume, entryPrice)
    if err != nil {
        fmt.Printf("Grid %d failed: %v\n", i, err)
        continue
    }

    pointsAway := (entryPrice - currentBid) / 0.00001
    fmt.Printf("Grid %d: %.5f (%.0f points) - #%d\n",
        i, entryPrice, pointsAway, ticket)
}
```

---

### 10) Advanced sell limit with monitoring

```go
func PlaceSellLimitWithMonitoring(sugar *mt5.MT5Sugar, symbol string, pointsAbove float64) {
    // Get current price
    currentBid, _ := sugar.GetBid(symbol)
    info, _ := sugar.GetSymbolInfo(symbol)

    // Calculate entry
    entryPrice := currentBid + (pointsAbove * info.Point)

    fmt.Printf("╔═══════════════════════════════════════╗\n")
    fmt.Printf("║       SELL LIMIT ORDER                ║\n")
    fmt.Printf("╚═══════════════════════════════════════╝\n")
    fmt.Printf("Symbol:       %s\n", symbol)
    fmt.Printf("Current BID:  %.5f\n", currentBid)
    fmt.Printf("Entry Price:  %.5f\n", entryPrice)
    fmt.Printf("Distance:     %.0f points\n", pointsAbove)

    // Place order
    ticket, err := sugar.SellLimit(symbol, 0.1, entryPrice)
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
            bid, _ := sugar.GetBid(symbol)
            distance := entryPrice - bid
            distancePoints := distance / info.Point

            // Check if filled
            _, err := sugar.GetPositionByTicket(ticket)
            if err == nil {
                fmt.Printf("\n✅ ORDER FILLED at %.5f!\n", bid)
                return
            }

            fmt.Printf("BID: %.5f (%.0f points away from entry)\n",
                bid, distancePoints)
        }
    }
}

// Usage:
PlaceSellLimitWithMonitoring(sugar, "EURUSD", 50)
```

---

## 🔗 Related Methods

**🍬 Other limit orders:**

* `BuyLimit()` - BUY when price goes DOWN
* `SellLimitWithSLTP()` - SELL LIMIT with SL/TP

**🍬 Stop orders (for breakdowns):**

* `SellStop()` - SELL when price breaks DOWN

**🍬 Market orders:**

* `SellMarket()` - SELL immediately at current price

---

## ⚠️ Common Pitfalls

### 1) Setting entry price below current price

```go
// ❌ WRONG - Sell Limit must be ABOVE current price
currentBid := 1.08500
sugar.SellLimit("EURUSD", 0.1, 1.08000) // BELOW! Will be rejected!

// ✅ CORRECT - Above current price
sugar.SellLimit("EURUSD", 0.1, 1.09000) // Above BID
```

### 2) Confusing SellLimit with SellStop

```go
// SELL LIMIT = wait for price to RISE then sell
sugar.SellLimit("EURUSD", 0.1, 1.09000) // Sell when rises to 1.09000

// SELL STOP = wait for price to FALL then sell (breakdown)
sugar.SellStop("EURUSD", 0.1, 1.08000) // Sell when falls to 1.08000
```

---

## 💎 Pro Tips

1. **SELL LIMIT = sell more expensive** - Wait for price to rise

2. **Set at resistance levels** - Use technical analysis

3. **Multiple levels** - Scale in at different prices

4. **Monitor fills** - Check if order was executed

5. **Use with SL/TP** - `SellLimitWithSLTP()` puts SL/TP into the pending order

---

**See also:** [`SellLimitWithSLTP.md`](../5.%20Trading_SLTP/SellLimitWithSLTP.md), [`SellStop.md`](SellStop.md), [`BuyLimit.md`](BuyLimit.md)
