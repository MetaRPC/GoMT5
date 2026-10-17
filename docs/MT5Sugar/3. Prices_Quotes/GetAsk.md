# 📈 Get ASK Price (`GetAsk`)

> **Sugar method:** Returns current ASK price for a symbol (price at which you can BUY).

**API Information:**

* **Method:** `sugar.GetAsk(symbol)`
* **Timeout:** 3 seconds
* **Returns:** ASK price as `float64`
* **Underlying call:** `service.GetSymbolTick(ctx, symbol)` (`SymbolInfoTick`), returns `tick.Ask`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) GetAsk(symbol string) (float64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `ask` | `float64` | Current ASK price |
| `error` | `error` | Error if symbol not found |

---

## 💬 Just the Essentials

* **What it is:** Current ASK price - the price you pay when **BUYING**.
* **Why you need it:** Know exact entry price for BUY orders, calculate SL/TP levels.
* **Sanity check:** ASK ≥ BID (equal only on a zero spread).

---

## 🎯 When to Use

✅ **Before BUY orders** - Know exact entry price

✅ **SL/TP calculation** - Calculate levels based on current price

✅ **Price monitoring** - Track market movement

✅ **Spread checking** - Compare with BID to get spread

---

## 🔗 Usage Examples

### 1) Basic usage

```go
ask, err := sugar.GetAsk("EURUSD")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("EURUSD ASK: %.5f\n", ask)
// Output: EURUSD ASK: 1.08460
```

---

### 2) Check price before buying

```go
symbol := "EURUSD"

// Get current ASK
ask, _ := sugar.GetAsk(symbol)

fmt.Printf("Current %s ASK: %.5f\n", symbol, ask)
fmt.Println("This is the price you'll pay if you BUY now")

// Place BUY order
ticket, _ := sugar.BuyMarket(symbol, 0.1)
fmt.Printf("BUY order #%d placed at ~%.5f\n", ticket, ask)
```

---

### 3) Calculate SL/TP for BUY order

```go
symbol := "EURUSD"
ask, _ := sugar.GetAsk(symbol)

// Get symbol info for point size
info, _ := sugar.GetSymbolInfo(symbol)

// Calculate SL/TP (BUY: SL below, TP above)
sl := ask - (50 * info.Point)  // 50 points below entry
tp := ask + (100 * info.Point) // 100 points above entry

fmt.Printf("Entry (ASK):    %.5f\n", ask)
fmt.Printf("Stop Loss:      %.5f (-50 points)\n", sl)
fmt.Printf("Take Profit:    %.5f (+100 points)\n", tp)

// Same result with the CalculateSLTP helper (its distances are points too)
sl2, tp2, _ := sugar.CalculateSLTP(symbol, "BUY", ask, 50, 100)
fmt.Printf("CalculateSLTP:  SL %.5f, TP %.5f\n", sl2, tp2)
```

---

### 4) Show BID/ASK spread

```go
symbol := "EURUSD"

bid, _ := sugar.GetBid(symbol)
ask, _ := sugar.GetAsk(symbol)
info, _ := sugar.GetSymbolInfo(symbol)

spreadInPrice := ask - bid
spreadInPoints := spreadInPrice / info.Point

fmt.Printf("Symbol: %s\n", symbol)
fmt.Printf("BID:    %.5f\n", bid)
fmt.Printf("ASK:    %.5f\n", ask)
fmt.Printf("Spread: %.5f (%.0f points)\n", spreadInPrice, spreadInPoints)

// Output:
// Symbol: EURUSD
// BID:    1.08450
// ASK:    1.08460
// Spread: 0.00010 (10 points)
```

---

### 5) Wait for good price

```go
symbol := "EURUSD"
maxAcceptableAsk := 1.08500

fmt.Printf("Waiting for %s ASK below %.5f...\n", symbol, maxAcceptableAsk)

for {
    ask, _ := sugar.GetAsk(symbol)

    if ask <= maxAcceptableAsk {
        fmt.Printf("✅ Good price! ASK: %.5f\n", ask)
        // Place order
        ticket, _ := sugar.BuyMarket(symbol, 0.1)
        fmt.Printf("Order #%d placed\n", ticket)
        break
    }

    fmt.Printf("Waiting... Current: %.5f (%.5f above target)\n",
        ask, ask-maxAcceptableAsk)
    time.Sleep(5 * time.Second)
}
```

---

### 6) Compare ASK across brokers/accounts

```go
// If you have multiple accounts
accounts := []*mt5.MT5Sugar{sugarDemo, sugarLive}
accountNames := []string{"Demo", "Live"}

symbol := "EURUSD"

fmt.Printf("ASK price comparison for %s:\n", symbol)
fmt.Println("─────────────────────────────")

for i, sugar := range accounts {
    ask, _ := sugar.GetAsk(symbol)
    fmt.Printf("%s:  %.5f\n", accountNames[i], ask)
}

// Output:
// ASK price comparison for EURUSD:
// ─────────────────────────────
// Demo:  1.08460
// Live:  1.08462
```

---

### 7) Entry price preview

```go
symbol := "EURUSD"
volume := 0.1

ask, _ := sugar.GetAsk(symbol)
info, _ := sugar.GetSymbolInfo(symbol)

// Calculate position value
positionValue := volume * info.ContractSize * ask

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║         BUY ORDER PREVIEW             ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Symbol:         %s\n", symbol)
fmt.Printf("Volume:         %.2f lots\n", volume)
fmt.Printf("Entry (ASK):    %.5f\n", ask)
fmt.Printf("Position Value: $%.2f\n", positionValue)
```

---

### 8) Check spread before trading

```go
symbol := "EURUSD"
maxAcceptableSpreadPoints := 15.0

bid, _ := sugar.GetBid(symbol)
ask, _ := sugar.GetAsk(symbol)
info, _ := sugar.GetSymbolInfo(symbol)

spreadPoints := (ask - bid) / info.Point

fmt.Printf("%s Spread: %.0f points\n", symbol, spreadPoints)

if spreadPoints > maxAcceptableSpreadPoints {
    fmt.Printf("❌ Spread too high (%.0f > %.0f points)\n",
        spreadPoints, maxAcceptableSpreadPoints)
    fmt.Println("   Not trading right now")
} else {
    fmt.Println("✅ Spread acceptable - safe to trade")
    sugar.BuyMarket(symbol, 0.1)
}
```

---

### 9) Price change monitor

```go
symbol := "EURUSD"
checkInterval := 10 * time.Second

previousAsk := 0.0
ticker := time.NewTicker(checkInterval)
defer ticker.Stop()

for range ticker.C {
    ask, _ := sugar.GetAsk(symbol)

    if previousAsk > 0 {
        change := ask - previousAsk
        changePercent := (change / previousAsk) * 100

        fmt.Printf("%s ASK: %.5f (", symbol, ask)

        if change > 0 {
            fmt.Printf("↑ +%.5f, +%.4f%%)\n", change, changePercent)
        } else if change < 0 {
            fmt.Printf("↓ %.5f, %.4f%%)\n", change, changePercent)
        } else {
            fmt.Printf("→ no change)\n")
        }
    } else {
        fmt.Printf("%s ASK: %.5f (initial)\n", symbol, ask)
    }

    previousAsk = ask
}
```

---

### 10) Multi-symbol monitoring

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"}

ticker := time.NewTicker(3 * time.Second)
defer ticker.Stop()

for range ticker.C {
    fmt.Print("\033[H\033[2J") // Clear screen
    fmt.Println("╔═══════════════════════════════════════════════╗")
    fmt.Println("║          LIVE ASK PRICES (BUY)                ║")
    fmt.Println("╚═══════════════════════════════════════════════╝")
    fmt.Printf("Updated: %s\n\n", time.Now().Format("15:04:05"))

    fmt.Printf("%-10s  %-12s  %-12s  %s\n", "Symbol", "BID", "ASK", "Spread")
    fmt.Println("─────────────────────────────────────────────────")

    for _, symbol := range symbols {
        bid, _ := sugar.GetBid(symbol)
        ask, _ := sugar.GetAsk(symbol)
        info, _ := sugar.GetSymbolInfo(symbol)

        spreadPoints := (ask - bid) / info.Point

        fmt.Printf("%-10s  %-12.5f  %-12.5f  %.0f points\n",
            symbol, bid, ask, spreadPoints)
    }
}
```

---

## 🔗 Related Methods

* `GetBid()` - Current BID price (for selling)
* `GetSpread()` - Spread in points
* `GetPriceInfo()` - Get BID, ASK, and spread at once ⭐
* `BuyMarket()` - Buy at current ASK price

---

## ⚠️ Common Pitfalls

### 1) Using ASK for SELL orders

```go
// ❌ WRONG - ASK is for buying, not selling!
ask, _ := sugar.GetAsk("EURUSD")
sugar.SellMarket("EURUSD", 0.1) // Will execute at BID, not ASK!

// ✅ CORRECT - use GetBid for SELL orders
bid, _ := sugar.GetBid("EURUSD")
fmt.Printf("SELL will execute at: %.5f (BID)\n", bid)
```

### 2) Forgetting about spread

```go
// ❌ WRONG - not accounting for spread
ask, _ := sugar.GetAsk("EURUSD")
// Entry: 1.08460 (ASK)
// Immediate P/L: negative due to spread!

// ✅ CORRECT - understand spread impact
bid, _ := sugar.GetBid("EURUSD")
ask, _ := sugar.GetAsk("EURUSD")
spread := ask - bid
fmt.Printf("Entry cost includes %.5f spread\n", spread)
```

---

## 💎 Pro Tips

1. **ASK = BUY price** - Remember: ASK is always for buying
2. **Use GetPriceInfo()** - Get BID + ASK + spread (in points) from one method
3. **ASK ≥ BID** - ASK is higher by the spread (you pay the spread)
4. **Check spread first** - High spread = bad time to trade
5. **Price volatility** - ASK changes constantly during market hours

---

**See also:** [`GetBid.md`](GetBid.md), [`GetPriceInfo.md`](GetPriceInfo.md), [`GetSpread.md`](GetSpread.md)
