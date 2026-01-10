# 📉 Get BID Price (`GetBid`)

> **Sugar method:** Returns current BID price for a symbol (price at which you can SELL).

**API Information:**

* **Method:** `sugar.GetBid(symbol)`
* **Timeout:** 3 seconds
* **Returns:** BID price as `float64`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetBid(symbol string) (float64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `bid` | `float64` | Current BID price |
| `error` | `error` | Error if symbol not found |

---

## 💬 Just the Essentials

* **What it is:** Current BID price - the price you get when **SELLING**.
* **Why you need it:** Check entry price for SELL orders, calculate SL/TP levels.
* **Sanity check:** BID < ASK (always lower than ask price).

---

## 🎯 When to Use

✅ **Before SELL orders** - Know exact entry price

✅ **SL/TP calculation** - Calculate levels based on current price

✅ **Price monitoring** - Track market movement

✅ **Spread checking** - Compare with ASK to get spread

---

## 🔗 Usage Examples

### 1) Basic usage

```go
bid, err := sugar.GetBid("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("EURUSD BID: %.5f\n", bid)
// Output: EURUSD BID: 1.08450
```

---

### 2) Check price before selling

```go
symbol := "EURUSD"

// Get current BID
bid, _ := sugar.GetBid(symbol)

fmt.Printf("Current %s BID: %.5f\n", symbol, bid)
fmt.Println("This is the price you'll get if you SELL now")

// Place SELL order
ticket, _ := sugar.SellMarket(symbol, 0.1)
fmt.Printf("SELL order #%d placed at ~%.5f\n", ticket, bid)
```

---

### 3) Calculate SL/TP for SELL order

```go
symbol := "EURUSD"
bid, _ := sugar.GetBid(symbol)

// Get symbol info for point size
info, _ := sugar.GetSymbolInfo(symbol)

// Calculate SL/TP (SELL: SL above, TP below)
sl := bid + (50 * info.Point)  // 50 pips above entry
tp := bid - (100 * info.Point) // 100 pips below entry

fmt.Printf("Entry (BID):    %.5f\n", bid)
fmt.Printf("Stop Loss:      %.5f (+50 pips)\n", sl)
fmt.Printf("Take Profit:    %.5f (-100 pips)\n", tp)

// Better: use CalculateSLTP helper
sl2, tp2, _ := sugar.CalculateSLTP(symbol, "SELL", bid, 50, 100)
```

---

### 4) Monitor price movement

```go
symbol := "EURUSD"
previousBid := 0.0

ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for i := 0; i < 10; i++ {
    <-ticker.C

    bid, _ := sugar.GetBid(symbol)

    if previousBid > 0 {
        change := bid - previousBid
        fmt.Printf("%s BID: %.5f (change: %+.5f)\n", symbol, bid, change)
    } else {
        fmt.Printf("%s BID: %.5f\n", symbol, bid)
    }

    previousBid = bid
}
```

---

### 5) Compare multiple symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"}

fmt.Println("Current BID prices:")
fmt.Println("─────────────────────────────")

for _, symbol := range symbols {
    bid, err := sugar.GetBid(symbol)
    if err != nil {
        fmt.Printf("%s: Error - %v\n", symbol, err)
        continue
    }

    fmt.Printf("%-8s  %.5f\n", symbol, bid)
}

// Output:
// Current BID prices:
// ─────────────────────────────
// EURUSD    1.08450
// GBPUSD    1.26320
// USDJPY    149.850
// XAUUSD    2045.50
```

---

### 6) Check if price reached target

```go
symbol := "EURUSD"
targetPrice := 1.08000

fmt.Printf("Waiting for %s BID to reach %.5f...\n", symbol, targetPrice)

for {
    bid, _ := sugar.GetBid(symbol)

    if bid <= targetPrice {
        fmt.Printf("✅ Target reached! BID: %.5f\n", bid)
        break
    }

    fmt.Printf("Current: %.5f (%.5f away)\n", bid, bid-targetPrice)
    time.Sleep(2 * time.Second)
}
```

---

### 7) Calculate pip distance

```go
symbol := "EURUSD"
referencePrice := 1.08500

currentBid, _ := sugar.GetBid(symbol)
info, _ := sugar.GetSymbolInfo(symbol)

// Calculate distance in pips
pipDistance := (currentBid - referencePrice) / info.Point

fmt.Printf("Reference:  %.5f\n", referencePrice)
fmt.Printf("Current:    %.5f\n", currentBid)
fmt.Printf("Distance:   %.0f pips\n", pipDistance)

if pipDistance > 0 {
    fmt.Println("Price moved UP")
} else {
    fmt.Println("Price moved DOWN")
}
```

---

### 8) Show BID/ASK spread

```go
symbol := "EURUSD"

bid, _ := sugar.GetBid(symbol)
ask, _ := sugar.GetAsk(symbol)

// Better: use GetPriceInfo() to get both at once
priceInfo, _ := sugar.GetPriceInfo(symbol)

fmt.Printf("Symbol: %s\n", symbol)
fmt.Printf("BID:    %.5f\n", bid)
fmt.Printf("ASK:    %.5f\n", ask)
fmt.Printf("Spread: %.1f pips\n", priceInfo.SpreadPips)
```

---

### 9) Price alert system

```go
func PriceAlert(sugar *mt5.MT5Sugar, symbol string, alertPrice float64, direction string) {
    fmt.Printf("🔔 Alert set: %s BID %s %.5f\n", symbol, direction, alertPrice)

    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        bid, _ := sugar.GetBid(symbol)

        triggered := false
        if direction == "below" && bid <= alertPrice {
            triggered = true
        } else if direction == "above" && bid >= alertPrice {
            triggered = true
        }

        if triggered {
            fmt.Printf("🚨 ALERT! %s BID is %.5f (%s %.5f)\n",
                symbol, bid, direction, alertPrice)
            return
        }
    }
}

// Usage:
go PriceAlert(sugar, "EURUSD", 1.08000, "below")
```

---

### 10) Real-time price dashboard

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()

for range ticker.C {
    fmt.Print("\033[H\033[2J") // Clear screen
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║       REAL-TIME BID PRICES            ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Time: %s\n\n", time.Now().Format("15:04:05"))

    for _, symbol := range symbols {
        bid, _ := sugar.GetBid(symbol)
        fmt.Printf("%-8s  %.5f\n", symbol, bid)
    }

    time.Sleep(2 * time.Second)
}
```

---

## 🔗 Related Methods

* `GetAsk()` - Current ASK price (for buying)
* `GetSpread()` - Spread in points
* `GetPriceInfo()` - Get BID, ASK, and spread at once ⭐
* `SellMarket()` - Sell at current BID price

---

## ⚠️ Common Pitfalls

### 1) Using BID for BUY orders

```go
// ❌ WRONG - BID is for selling, not buying!
bid, _ := sugar.GetBid("EURUSD")
sugar.BuyMarket("EURUSD", 0.1) // Will execute at ASK, not BID!

// ✅ CORRECT - use GetAsk for BUY orders
ask, _ := sugar.GetAsk("EURUSD")
fmt.Printf("BUY will execute at: %.5f (ASK)\n", ask)
```

### 2) Not checking for errors

```go
// ❌ WRONG - ignoring errors
bid, _ := sugar.GetBid("INVALID_SYMBOL")
fmt.Printf("Price: %.5f\n", bid) // Will be 0!

// ✅ CORRECT - check errors
bid, err := sugar.GetBid("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
```

---

## 💎 Pro Tips

1. **BID = SELL price** - Remember: BID is always for selling
2. **Use GetPriceInfo()** - More efficient than calling GetBid + GetAsk separately
3. **BID < ASK** - BID is always lower (spread is the difference)
4. **Price changes** - BID updates with every tick, call frequently
5. **For monitoring** - Use GetPriceInfo() to get complete price snapshot

---

**See also:** [`GetAsk.md`](GetAsk.md), [`GetPriceInfo.md`](GetPriceInfo.md), [`GetSpread.md`](GetSpread.md)
