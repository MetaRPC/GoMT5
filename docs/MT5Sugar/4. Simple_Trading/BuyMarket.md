# 🟢 Buy at Market Price (`BuyMarket`)

> **Sugar method:** Opens a BUY position instantly at current market ASK price.

**API Information:**

* **Method:** `sugar.BuyMarket(symbol, volume)`
* **Timeout:** 10 seconds
* **Returns:** Ticket number - `result.Order`, the order ticket; for a market order MT5 uses the same number as the
  position ticket
* **Request:** `service.PlaceOrder(ctx, &pb.OrderSendRequest{Symbol, Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY, Volume})` -
  no price, SL/TP, slippage, comment or magic number; the volume is sent exactly as given (not rounded to the volume step)
* **Success:** only return code `10009` (TRADE_RETCODE_DONE); any other code gives the error
  `order rejected, code: N, comment: ...`, a failed call gives `BuyMarket failed: ...`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) BuyMarket(symbol string, volume float64) (uint64, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| `volume` | `float64` | Lot size (e.g., 0.01, 0.1, 1.0) |

| Output | Type | Description |
|--------|------|-------------|
| `ticket` | `uint64` | Order ticket = position ticket of the new position (for tracking/closing); 0 on error |
| `error` | `error` | Error if order rejected or execution failed |

---

## 💬 Just the Essentials

* **What it is:** Instantly opens a BUY position at current market ASK price.
* **Why you need it:** Simplest way to go LONG on a symbol.
* **Sanity check:** Order executes immediately - no waiting, no SL/TP set.

---

## ⚠️ IMPORTANT

**This method does NOT set Stop Loss or Take Profit!**

- Use `BuyMarketWithSLTP()` for trades with risk management
- Use `BuyMarketWithPips()` for trades with SL/TP as distances (in points, despite the name)
- Or manually add SL/TP after: `ModifyPositionSLTP(ticket, sl, tp)`

---

## 🎯 When to Use

✅ **Quick entries** - When you want instant market entry

✅ **Manual management** - You'll set SL/TP manually later

✅ **Testing/learning** - Simple way to test trading

❌ **NOT for production** - Always use SL/TP in real trading!

---

## 🔗 Usage Examples

### 1) Basic usage

```go
ticket, err := sugar.BuyMarket("EURUSD", 0.1)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY position opened, ticket #%d\n", ticket)
```

---

### 2) With error handling

```go
symbol := "EURUSD"
volume := 0.1

ticket, err := sugar.BuyMarket(symbol, volume)
if err != nil {
    fmt.Printf("❌ BUY failed: %v\n", err)
    return
}

fmt.Printf("✅ BUY %s %.2f lots\n", symbol, volume)
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("⚠️  No SL/TP set - add manually!\n")
```

---

### 3) Add SL/TP immediately after opening

```go
symbol := "EURUSD"
volume := 0.1

// Open position
ticket, err := sugar.BuyMarket(symbol, volume)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("Position #%d opened\n", ticket)

// Calculate SL/TP (50 points SL, 100 points TP)
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, 50, 100)

// Set SL/TP
err = sugar.ModifyPositionSLTP(ticket, sl, tp)
if err != nil {
    fmt.Printf("⚠️  Failed to set SL/TP: %v\n", err)
} else {
    fmt.Printf("✅ SL/TP set: %.5f / %.5f\n", sl, tp)
}
```

---

### 4) Check price before buying

```go
symbol := "EURUSD"
maxPrice := 1.09000

// Check current ASK
ask, _ := sugar.GetAsk(symbol)

if ask > maxPrice {
    fmt.Printf("❌ Price too high: %.5f > %.5f\n", ask, maxPrice)
    return
}

fmt.Printf("✅ Good price: %.5f\n", ask)

// Place order
ticket, _ := sugar.BuyMarket(symbol, 0.1)
fmt.Printf("BUY order #%d placed at ~%.5f\n", ticket, ask)
```

---

### 5) Risk-managed entry

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLossPips := 50.0    // distance in POINTS (the Sugar "Pips" parameters are points)
takeProfitPips := 100.0 // points

// Step 1: Calculate position size based on risk
lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)
if err != nil {
    fmt.Printf("Position size calc failed: %v\n", err)
    return
}

// Step 2: Validate we can open
canOpen, reason, _ := sugar.CanOpenPosition(symbol, lotSize)
if !canOpen {
    fmt.Printf("Cannot open: %s\n", reason)
    return
}

// Step 3: Open position
ticket, err := sugar.BuyMarket(symbol, lotSize)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

// Step 4: Set SL/TP
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", 0, stopLossPips, takeProfitPips)
sugar.ModifyPositionSLTP(ticket, sl, tp)

fmt.Printf("✅ Complete setup!\n")
fmt.Printf("   Ticket: #%d\n", ticket)
fmt.Printf("   Size: %.2f lots (%.1f%% risk)\n", lotSize, riskPercent)
fmt.Printf("   SL: %.0f points, TP: %.0f points\n", stopLossPips, takeProfitPips)
```

---

### 6) Open multiple positions

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
volume := 0.1

fmt.Println("Opening BUY positions:")

for _, symbol := range symbols {
    ticket, err := sugar.BuyMarket(symbol, volume)
    if err != nil {
        fmt.Printf("❌ %s: Failed - %v\n", symbol, err)
        continue
    }

    fmt.Printf("✅ %s: Ticket #%d\n", symbol, ticket)
}
```

---

### 7) Track position after opening

```go
symbol := "EURUSD"
volume := 0.1

// Open position
ticket, err := sugar.BuyMarket(symbol, volume)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("Position #%d opened\n", ticket)

// Monitor position
time.Sleep(5 * time.Second)

pos, err := sugar.GetPositionByTicket(ticket)
if err != nil {
    fmt.Printf("Position not found: %v\n", err)
    return
}

fmt.Printf("\nPosition Status:\n")
fmt.Printf("  Symbol:      %s\n", pos.Symbol)
fmt.Printf("  Volume:      %.2f\n", pos.Volume)
fmt.Printf("  Open Price:  %.5f\n", pos.PriceOpen)
fmt.Printf("  Current P/L: $%.2f\n", pos.Profit)
```

---

### 8) Wait for confirmation

```go
symbol := "EURUSD"
volume := 0.1

fmt.Printf("Placing BUY order for %s...\n", symbol)

ticket, err := sugar.BuyMarket(symbol, volume)
if err != nil {
    fmt.Printf("❌ Order REJECTED: %v\n", err)
    return
}

fmt.Printf("✅ Order FILLED\n")
fmt.Printf("   Ticket:  #%d\n", ticket)

// Verify position exists
time.Sleep(1 * time.Second)
positions, _ := sugar.GetOpenPositions()

found := false
for _, pos := range positions {
    if pos.Ticket == ticket {
        found = true
        fmt.Printf("   Entry:   %.5f\n", pos.PriceOpen)
        fmt.Printf("   Volume:  %.2f lots\n", pos.Volume)
        break
    }
}

if !found {
    fmt.Println("⚠️  WARNING: Position not found in open positions!")
}
```

---

### 9) Scale into position

```go
symbol := "EURUSD"
totalLots := 0.3
entries := 3

lotPerEntry := totalLots / float64(entries)

fmt.Printf("Scaling into %s with %d entries of %.2f lots each\n",
    symbol, entries, lotPerEntry)

tickets := []uint64{}

for i := 1; i <= entries; i++ {
    fmt.Printf("\nEntry %d/%d: ", i, entries)

    ticket, err := sugar.BuyMarket(symbol, lotPerEntry)
    if err != nil {
        fmt.Printf("❌ Failed - %v\n", err)
        continue
    }

    tickets = append(tickets, ticket)
    fmt.Printf("✅ Ticket #%d\n", ticket)

    // Wait between entries
    if i < entries {
        time.Sleep(5 * time.Second)
    }
}

fmt.Printf("\nOpened %d/%d positions\n", len(tickets), entries)
fmt.Printf("Total tickets: %v\n", tickets)
```

---

### 10) Advanced order with pre-checks

```go
func PlaceBuyOrder(sugar *mt5.MT5Sugar, symbol string, volume float64) (uint64, error) {
    // Pre-check 1: Connection
    if !sugar.IsConnected() {
        return 0, fmt.Errorf("not connected to MT5")
    }

    // Pre-check 2: Symbol availability
    available, _ := sugar.IsSymbolAvailable(symbol)
    if !available {
        return 0, fmt.Errorf("symbol %s not available", symbol)
    }

    // Pre-check 3: Spread check
    spread, _ := sugar.GetSpread(symbol)
    if spread > 20 {
        return 0, fmt.Errorf("spread too high: %.0f points", spread)
    }

    // Pre-check 4: Margin check
    canOpen, reason, _ := sugar.CanOpenPosition(symbol, volume)
    if !canOpen {
        return 0, fmt.Errorf("cannot open: %s", reason)
    }

    // All checks passed - place order
    fmt.Printf("✅ All pre-checks passed\n")
    fmt.Printf("   Symbol:  %s\n", symbol)
    fmt.Printf("   Volume:  %.2f lots\n", volume)
    fmt.Printf("   Spread:  %.0f points\n", spread)

    ticket, err := sugar.BuyMarket(symbol, volume)
    if err != nil {
        return 0, fmt.Errorf("order execution failed: %w", err)
    }

    fmt.Printf("✅ Order filled - Ticket #%d\n", ticket)
    return ticket, nil
}

// Usage:
ticket, err := PlaceBuyOrder(sugar, "EURUSD", 0.1)
```

---

## 🔗 Related Methods

**🍬 Better alternatives with SL/TP:**

* `BuyMarketWithSLTP()` - BUY with SL/TP prices
* `BuyMarketWithPips()` - BUY with SL/TP distances in points ⭐ **RECOMMENDED**

**🍬 Other market orders:**

* `SellMarket()` - SELL at market

**🍬 Position management:**

* `ModifyPositionSLTP()` - Add SL/TP after opening
* `ClosePosition()` - Close the position

---

## ⚠️ Common Pitfalls

### 1) Trading without Stop Loss

```go
// ❌ DANGEROUS - No stop loss!
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)
// Position has no protection!

// ✅ CORRECT - Use version with SL/TP
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
```

### 2) Not checking for errors

```go
// ❌ WRONG - ignoring errors
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)
fmt.Printf("Ticket: %d\n", ticket) // Might be 0!

// ✅ CORRECT - check errors
ticket, err := sugar.BuyMarket("EURUSD", 0.1)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}
```

### 3) Using fixed lot size

```go
// ❌ WRONG - fixed lot size (no risk management)
sugar.BuyMarket("EURUSD", 0.1)

// ✅ CORRECT - calculate based on risk
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
sugar.BuyMarket("EURUSD", lotSize)
```

### 4) Not validating before trading

```go
// ❌ WRONG - trading without checks
sugar.BuyMarket("EURUSD", 10.0) // Might fail!

// ✅ CORRECT - validate first
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 10.0)
if !canOpen {
    fmt.Println("Cannot trade:", reason)
    return
}
sugar.BuyMarket("EURUSD", 10.0)
```

---

## 💎 Pro Tips

1. **Always use SL/TP** - Use `BuyMarketWithPips()` instead for safer trading

2. **Calculate position size** - Use `CalculatePositionSize()` for risk management

3. **Validate first** - Use `CanOpenPosition()` before placing order

4. **Check spread** - High spread = bad entry

5. **For production** - Never use this method without immediately adding SL/TP

---

## 🚨 Production Recommendation

**DON'T use `BuyMarket()` in production!**

Instead, use the safer alternatives:

```go
// ✅ RECOMMENDED for production
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)

// OR with full risk management
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

---

**See also:** [`BuyMarketWithPips.md`](../11.%20Trading_Helpers/BuyMarketWithPips.md), [`BuyMarketWithSLTP.md`](../5.%20Trading_SLTP/BuyMarketWithSLTP.md), [`SellMarket.md`](SellMarket.md)
