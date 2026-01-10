# 🔴📉 Sell Stop Order (`SellStop`)

> **Sugar method:** Places pending SELL order that executes when price **breaks BELOW** specified level.

**API Information:**

* **Method:** `sugar.SellStop(symbol, volume, price)`
* **Timeout:** 10 seconds
* **Returns:** Pending order ticket number

---

## 📋 Method Signature

```go
func (s *MT5Sugar) SellStop(symbol string, volume, price float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD") |
| `volume` | `float64` | Lot size (e.g., 0.1) |
| `price` | `float64` | Entry price (must be **BELOW** current BID) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Pending order ticket number |
| `error` | `error` | Error if order rejected |

---

## 💬 Just the Essentials

* **What it is:** Pending SELL order that activates when price **breaks BELOW** support.
* **Why you need it:** Sell on breakdowns - when price confirms downward momentum.
* **Sanity check:** Entry price must be **BELOW** current BID. Order activates on breakdown.

---

## 📊 SELL STOP Logic

```
Current BID: 1.08500
SELL STOP:   1.08000 (500 pips BELOW current price)

When price falls to 1.08000 → Order executes automatically
```

**Use case:** Breakdown trading / Sell on support break / Momentum trading

---

## 🎯 When to Use

✅ **Breakdown trading** - Sell when price breaks support

✅ **Momentum trading** - Enter when downtrend confirmed

✅ **Range breakdowns** - Sell when price exits consolidation downward

✅ **Follow strong trends** - Join established downward momentum

❌ **NOT for rallies** - Use `SellLimit()` for selling rallies

---

## 🔗 Usage Examples

### 1) Basic usage - sell on breakdown

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
supportLevel := 1.08000 // Known support

// Place sell stop below support
ticket, err := sugar.SellStop(symbol, 0.1, supportLevel)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ SELL STOP order placed\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Entry:  %.5f (current: %.5f)\n", supportLevel, currentBid)
fmt.Printf("   Will execute when price breaks %.0f pips lower\n",
    (currentBid-supportLevel)/0.00001)
```

---

### 2) Range breakdown trading

```go
symbol := "EURUSD"

// Price consolidating between 1.08000 and 1.08500
rangeLow := 1.08000
rangeLow -= 0.00010 // Place stop 10 pips below range low

ticket, _ := sugar.SellStop(symbol, 0.1, rangeLow)

fmt.Printf("SELL STOP set for range breakdown\n")
fmt.Printf("Range low: %.5f\n", 1.08000)
fmt.Printf("Entry:     %.5f (10 pips below)\n", rangeLow)
fmt.Printf("Strategy:  Sell if price breaks out of range downward\n")
```

---

### 3) Support breakdown with multiple levels

```go
symbol := "EURUSD"
currentBid, _ := sugar.GetBid(symbol)
volume := 0.05

// Set sell stops at different support breakdown levels
supportLevels := []float64{
    1.08000, // First support
    1.07500, // Second support
    1.07000, // Third support
}

fmt.Println("Placing multiple SELL STOP orders:")

for i, level := range supportLevels {
    // Only place stops below current price
    if level >= currentBid {
        fmt.Printf("Level %d: %.5f - SKIPPED (above current price)\n", i+1, level)
        continue
    }

    ticket, err := sugar.SellStop(symbol, volume, level)
    if err != nil {
        fmt.Printf("Level %d failed: %v\n", i+1, err)
        continue
    }

    pipsAway := (currentBid - level) / 0.00001
    fmt.Printf("✅ Level %d: %.5f (%.0f pips away) - Ticket #%d\n",
        i+1, level, pipsAway, ticket)
}
```

---

### 4) Sell stop with validation

```go
symbol := "EURUSD"
volume := 0.1
breakdownPrice := 1.08000

// Get current price
currentBid, _ := sugar.GetBid(symbol)

// Validate entry price is below current
if breakdownPrice >= currentBid {
    fmt.Printf("❌ Error: Entry %.5f must be below BID %.5f\n",
        breakdownPrice, currentBid)
    return
}

// Check minimum stop level
info, _ := sugar.GetSymbolInfo(symbol)
stopLevel, _ := sugar.GetMinStopLevel(symbol)
minDistance := float64(stopLevel) * info.Point

if (currentBid - breakdownPrice) < minDistance {
    fmt.Printf("❌ Entry too close (min: %.5f)\n", minDistance)
    return
}

// Place order
ticket, _ := sugar.SellStop(symbol, volume, breakdownPrice)
fmt.Printf("✅ SELL STOP placed at %.5f\n", breakdownPrice)
```

---

### 5) Sell stop with SL/TP (ATR-based)

```go
symbol := "EURUSD"
breakdownPrice := 1.08000

// Place sell stop
ticket, _ := sugar.SellStop(symbol, 0.1, breakdownPrice)

// Calculate SL/TP based on ATR or fixed pips
// For breakdown: SL above breakdown level, TP 2x the distance
info, _ := sugar.GetSymbolInfo(symbol)
slDistance := 30.0 // 30 pips SL
tpDistance := 60.0 // 60 pips TP (2:1 R:R)

sl := breakdownPrice + (slDistance * info.Point)
tp := breakdownPrice - (tpDistance * info.Point)

// Set SL/TP on pending order
sugar.ModifyPositionSLTP(ticket, sl, tp)

fmt.Printf("SELL STOP with SL/TP:\n")
fmt.Printf("  Entry: %.5f (breakdown)\n", breakdownPrice)
fmt.Printf("  SL:    %.5f (+30 pips)\n", sl)
fmt.Printf("  TP:    %.5f (-60 pips)\n", tp)
fmt.Printf("  R:R:   1:2\n")
```

---

### 6) Trend continuation sell stop

```go
symbol := "EURUSD"

// Strong downtrend - previous low at 1.08000
previousLow := 1.08000
entryPrice := previousLow - 0.00010 // 10 pips below previous low

ticket, _ := sugar.SellStop(symbol, 0.1, entryPrice)

fmt.Printf("Trend continuation SELL STOP\n")
fmt.Printf("Previous low: %.5f\n", previousLow)
fmt.Printf("Entry:        %.5f (-10 pips)\n", entryPrice)
fmt.Printf("Strategy:     Sell if downtrend continues\n")
```

---

### 7) News breakout strategy (OCO)

```go
// Before major news release - price at 1.08500
symbol := "EURUSD"
currentPrice := 1.08500

// Place sell stop for downside breakout
sellStopPrice := currentPrice - 0.00030 // 30 pips below
sellTicket, _ := sugar.SellStop(symbol, 0.1, sellStopPrice)

// Also place buy stop for upside breakout
buyStopPrice := currentPrice + 0.00030 // 30 pips above
buyTicket, _ := sugar.BuyStop(symbol, 0.1, buyStopPrice)

fmt.Printf("News breakout strategy (OCO):\n")
fmt.Printf("  Current:    %.5f\n", currentPrice)
fmt.Printf("  SELL STOP:  %.5f (ticket #%d)\n", sellStopPrice, sellTicket)
fmt.Printf("  BUY STOP:   %.5f (ticket #%d)\n", buyStopPrice, buyTicket)
fmt.Printf("  Strategy:   Whichever direction breaks first\n")
fmt.Printf("  Important:  Cancel other order when one fills!\n")
```

---

### 8) Trailing sell stop (manual adjustment)

```go
symbol := "EURUSD"
supportLevel := 1.08000

// Initial sell stop
ticket, _ := sugar.SellStop(symbol, 0.1, supportLevel)
fmt.Printf("Initial SELL STOP: %.5f - Ticket #%d\n", supportLevel, ticket)

// Monitor price and adjust stop level
ticker := time.NewTicker(1 * time.Minute)
defer ticker.Stop()

for i := 0; i < 10; i++ {
    <-ticker.C

    currentBid, _ := sugar.GetBid(symbol)

    // If price moves down, move sell stop down (trailing)
    newStopLevel := currentBid - 0.00050 // Always 50 pips below

    if newStopLevel < supportLevel {
        // Modify pending order price
        fmt.Printf("Adjusting SELL STOP from %.5f to %.5f\n",
            supportLevel, newStopLevel)

        // Use OrderModify to change pending order price
        // sugar.ModifyOrder(ticket, newStopLevel, sl, tp)

        supportLevel = newStopLevel
    } else {
        fmt.Printf("Current: %.5f - Stop remains at %.5f\n",
            currentBid, supportLevel)
    }
}
```

---

### 9) Pattern breakdown - double bottom

```go
symbol := "EURUSD"

// Double bottom pattern - two lows at 1.08000
doubleBottomLevel := 1.08000

// Place sell stop below double bottom (invalidation = strong breakdown)
entryPrice := doubleBottomLevel - 0.00015 // 15 pips below

ticket, _ := sugar.SellStop(symbol, 0.1, entryPrice)

fmt.Printf("Double bottom breakdown strategy:\n")
fmt.Printf("Double bottom: %.5f\n", doubleBottomLevel)
fmt.Printf("SELL STOP:     %.5f (-15 pips)\n", entryPrice)
fmt.Printf("Logic:         If double bottom breaks = strong bearish signal\n")
```

---

### 10) Advanced breakdown with monitoring

```go
func PlaceSellStopWithMonitoring(sugar *mt5.MT5Sugar, symbol string, breakdownLevel float64) {
    currentBid, _ := sugar.GetBid(symbol)

    fmt.Printf("╔═══════════════════════════════════════╗\n")
    fmt.Printf("║       BREAKDOWN SELL STOP             ║\n")
    fmt.Printf("╚═══════════════════════════════════════╝\n")
    fmt.Printf("Symbol:       %s\n", symbol)
    fmt.Printf("Current BID:  %.5f\n", currentBid)
    fmt.Printf("Breakdown:    %.5f\n", breakdownLevel)
    fmt.Printf("Distance:     %.0f pips\n", (currentBid-breakdownLevel)/0.00001)

    // Place sell stop
    ticket, err := sugar.SellStop(symbol, 0.1, breakdownLevel)
    if err != nil {
        fmt.Printf("❌ Order failed: %v\n", err)
        return
    }

    fmt.Printf("✅ Order placed: #%d\n", ticket)
    fmt.Println("\nMonitoring for breakdown...")

    // Monitor until filled or timeout
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    timeout := time.After(2 * time.Hour)

    highestBid := currentBid

    for {
        select {
        case <-timeout:
            fmt.Println("⏰ Timeout - breakdown didn't occur in 2 hours")
            return

        case <-ticker.C:
            bid, _ := sugar.GetBid(symbol)
            distance := bid - breakdownLevel
            distancePips := distance / 0.00001

            if bid > highestBid {
                highestBid = bid
            }

            // Check if filled
            _, err := sugar.GetPositionByTicket(ticket)
            if err == nil {
                fmt.Printf("\n📉 BREAKDOWN! Order filled at %.5f!\n", bid)
                fmt.Printf("   High before breakdown: %.5f\n", highestBid)
                return
            }

            if distancePips < 10 {
                fmt.Printf("⚠️  BID: %.5f (%.0f pips from breakdown - CLOSE!)\n",
                    bid, distancePips)
            } else {
                fmt.Printf("BID: %.5f (%.0f pips from breakdown)\n",
                    bid, distancePips)
            }
        }
    }
}

// Usage:
PlaceSellStopWithMonitoring(sugar, "EURUSD", 1.08000)
```

---

## 🔗 Related Methods

**🍬 Other stop orders:**

* `BuyStop()` - BUY when price breaks UP
* `SellStopWithSLTP()` - SELL STOP with SL/TP

**🍬 Limit orders (for rallies):**

* `SellLimit()` - SELL when price rises to resistance

**🍬 Market orders:**

* `SellMarket()` - SELL immediately at current price

---

## ⚠️ Common Pitfalls

### 1) Setting entry price above current price

```go
// ❌ WRONG - Sell Stop must be BELOW current price
currentBid := 1.08500
sugar.SellStop("EURUSD", 0.1, 1.09000) // ABOVE! Will be rejected!

// ✅ CORRECT - Below current price
sugar.SellStop("EURUSD", 0.1, 1.08000) // Below BID
```

### 2) Confusing SellStop with SellLimit

```go
// SELL STOP = wait for price to FALL then sell (breakdown)
sugar.SellStop("EURUSD", 0.1, 1.08000) // Sell when falls to 1.08000

// SELL LIMIT = wait for price to RISE then sell (rally)
sugar.SellLimit("EURUSD", 0.1, 1.09000) // Sell when rises to 1.09000
```

### 3) Chasing breakdowns too close to current price

```go
// ❌ WRONG - too close to current price (might trigger immediately)
currentBid := 1.08005
sugar.SellStop("EURUSD", 0.1, 1.08000) // Only 5 pips away!

// ✅ CORRECT - proper distance for confirmed breakdown
sugar.SellStop("EURUSD", 0.1, 1.07500) // Clear breakdown level
```

---

## 💎 Pro Tips

1. **SELL STOP = sell on weakness** - Wait for breakdown confirmation

2. **Set below support** - Use previous lows, round numbers

3. **Add buffer** - Place 10-20 pips below support for confirmation

4. **Use SL above breakdown** - If price reverses, exit quickly

5. **Remember OCO** - If using both buy/sell stops, cancel one when other fills

---

**See also:** [`BuyStop.md`](BuyStop.md), [`SellLimit.md`](SellLimit.md)
