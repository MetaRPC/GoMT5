# 🔴 Sell at Market Price (`SellMarket`)

> **Sugar method:** Opens a SELL position instantly at current market BID price.

**API Information:**

* **Method:** `sugar.SellMarket(symbol, volume)`
* **Timeout:** 30 seconds
* **Returns:** Ticket number - `result.Order`, the order ticket; for a market order MT5 uses the same number as the
  position ticket
* **Request:** `service.PlaceOrder(ctx, &pb.OrderSendRequest{Symbol, Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL, Volume})` -
  no price, SL/TP, slippage, comment or magic number; the volume is sent exactly as given (not rounded to the volume step)
* **Success:** only return code `10009` (TRADE_RETCODE_DONE); any other code gives the error
  `order rejected, code: N, comment: ...`, a failed call gives `SellMarket failed: ...`

---

## 📋 Method Signature

```go
func (s *MT5Sugar) SellMarket(symbol string, volume float64) (uint64, error)
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

* **What it is:** Instantly opens a SELL position at current market BID price.
* **Why you need it:** Simplest way to go SHORT on a symbol.
* **Sanity check:** Order executes immediately - no waiting, no SL/TP set.

---

## ⚠️ IMPORTANT

**This method does NOT set Stop Loss or Take Profit!**

- Use `SellMarketWithSLTP()` for trades with risk management
- Use `SellMarketWithPips()` for trades with SL/TP as distances (in points, despite the name)
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
ticket, err := sugar.SellMarket("EURUSD", 0.1)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ SELL position opened, ticket #%d\n", ticket)
```

---

### 2) With immediate SL/TP

```go
symbol := "EURUSD"
volume := 0.1

// Open SELL position
ticket, err := sugar.SellMarket(symbol, volume)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("Position #%d opened\n", ticket)

// SELL: SL above entry, TP below entry (50 / 100 points from the current BID)
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, 50, 100)

// Set SL/TP
sugar.ModifyPositionSLTP(ticket, sl, tp)
fmt.Printf("✅ SL/TP set: %.5f / %.5f\n", sl, tp)
```

---

### 3) Sell at resistance level

```go
symbol := "EURUSD"
resistance := 1.09000

bid, _ := sugar.GetBid(symbol)

if bid < resistance {
    fmt.Printf("❌ Price %.5f below resistance %.5f\n", bid, resistance)
    return
}

fmt.Printf("✅ At resistance: %.5f\n", bid)

// Sell at resistance
ticket, _ := sugar.SellMarket(symbol, 0.1)
fmt.Printf("SELL order #%d placed\n", ticket)
```

---

### 4) Risk-managed SELL

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLossPips := 50.0    // distance in POINTS (the Sugar "Pips" parameters are points)
takeProfitPips := 100.0 // points

// Calculate position size
lotSize, _ := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)

// Validate
canOpen, reason, _ := sugar.CanOpenPosition(symbol, lotSize)
if !canOpen {
    fmt.Printf("Cannot open: %s\n", reason)
    return
}

// Open SELL
ticket, _ := sugar.SellMarket(symbol, lotSize)

// Set SL/TP
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, stopLossPips, takeProfitPips)
sugar.ModifyPositionSLTP(ticket, sl, tp)

fmt.Printf("✅ SELL setup complete\n")
fmt.Printf("   Ticket: #%d, Size: %.2f lots\n", ticket, lotSize)
```

---

### 5) Hedge existing BUY position

```go
symbol := "EURUSD"

// GetPositionsBySymbol returns BUY and SELL positions - keep only the BUYs
positions, _ := sugar.GetPositionsBySymbol(symbol)

var totalBuyVolume float64
for _, pos := range positions {
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        totalBuyVolume += pos.Volume
    }
}

if totalBuyVolume == 0 {
    fmt.Println("No BUY positions to hedge")
    return
}

// Note: a separate opposite position needs a hedging account;
// on a netting account this SELL reduces/closes the BUY instead.

fmt.Printf("Total BUY volume: %.2f lots\n", totalBuyVolume)

// Open SELL to hedge
ticket, _ := sugar.SellMarket(symbol, totalBuyVolume)

fmt.Printf("✅ Hedge SELL opened: Ticket #%d\n", ticket)
fmt.Printf("   Volume: %.2f lots (matches BUY positions)\n", totalBuyVolume)
```

---

### 6) Pair trading - BUY one, SELL another

```go
// Buy EURUSD, Sell GBPUSD (pair trade)
ticketBuy, _ := sugar.BuyMarket("EURUSD", 0.1)
ticketSell, _ := sugar.SellMarket("GBPUSD", 0.1)

fmt.Println("Pair trade opened:")
fmt.Printf("  BUY EURUSD:  Ticket #%d\n", ticketBuy)
fmt.Printf("  SELL GBPUSD: Ticket #%d\n", ticketSell)
```

---

### 7) Sell on breakdown

```go
symbol := "EURUSD"
supportLevel := 1.08000

fmt.Printf("Monitoring %s for breakdown below %.5f...\n", symbol, supportLevel)

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for range ticker.C {
    bid, _ := sugar.GetBid(symbol)

    if bid < supportLevel {
        fmt.Printf("🔴 BREAKDOWN! BID %.5f broke below %.5f\n", bid, supportLevel)

        // Sell on breakdown
        ticket, _ := sugar.SellMarket(symbol, 0.1)
        fmt.Printf("✅ SELL order #%d placed\n", ticket)
        return
    }

    fmt.Printf("BID: %.5f (%.5f above support)\n", bid, bid-supportLevel)
}
```

---

### 8) Average into short position

```go
symbol := "EURUSD"
entries := 3
lotPerEntry := 0.1

fmt.Printf("Averaging into SELL position with %d entries\n", entries)

tickets := []uint64{}

for i := 1; i <= entries; i++ {
    bid, _ := sugar.GetBid(symbol)

    ticket, err := sugar.SellMarket(symbol, lotPerEntry)
    if err != nil {
        fmt.Printf("Entry %d failed: %v\n", i, err)
        continue
    }

    tickets = append(tickets, ticket)
    fmt.Printf("Entry %d: Ticket #%d at %.5f\n", i, ticket, bid)

    if i < entries {
        time.Sleep(10 * time.Second) // Wait before next entry
    }
}

fmt.Printf("\nTotal: %d SELL positions opened\n", len(tickets))
```

---

### 9) Quick scalp SELL

```go
symbol := "EURUSD"
volume := 0.1

// Check spread (scalping needs tight spread)
spread, _ := sugar.GetSpread(symbol)
if spread > 10 {
    fmt.Printf("Spread too high for scalping: %.0f points\n", spread)
    return
}

// Quick SELL entry
ticket, _ := sugar.SellMarket(symbol, volume)
fmt.Printf("Scalp SELL opened: #%d\n", ticket)

// Set tight SL/TP (15 / 10 points)
sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, 15, 10)
sugar.ModifyPositionSLTP(ticket, sl, tp)

fmt.Println("Tight SL/TP set for quick scalp")
```

---

### 10) Complete SELL workflow

```go
func OpenSellPosition(sugar *mt5.MT5Sugar, symbol string) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║       SELL ORDER WORKFLOW             ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Step 1: Check connection
    if !sugar.IsConnected() {
        return fmt.Errorf("not connected")
    }
    fmt.Println("✅ Connected")

    // Step 2: Check spread
    spread, _ := sugar.GetSpread(symbol)
    fmt.Printf("✅ Spread: %.0f points\n", spread)
    if spread > 20 {
        return fmt.Errorf("spread too high: %.0f", spread)
    }

    // Step 3: Get current price
    bid, _ := sugar.GetBid(symbol)
    fmt.Printf("✅ Current BID: %.5f\n", bid)

    // Step 4: Calculate position size (2% risk, 50-point SL)
    lotSize, _ := sugar.CalculatePositionSize(symbol, 2.0, 50)
    fmt.Printf("✅ Position size: %.2f lots\n", lotSize)

    // Step 5: Validate
    canOpen, reason, _ := sugar.CanOpenPosition(symbol, lotSize)
    if !canOpen {
        return fmt.Errorf("validation failed: %s", reason)
    }
    fmt.Println("✅ Validation passed")

    // Step 6: Open SELL
    ticket, err := sugar.SellMarket(symbol, lotSize)
    if err != nil {
        return fmt.Errorf("order failed: %w", err)
    }
    fmt.Printf("✅ SELL opened: Ticket #%d\n", ticket)

    // Step 7: Set SL/TP
    sl, tp, _ := sugar.CalculateSLTP(symbol, "SELL", 0, 50, 100)
    sugar.ModifyPositionSLTP(ticket, sl, tp)
    fmt.Printf("✅ SL/TP: %.5f / %.5f\n", sl, tp)

    fmt.Println("\n✅ SELL POSITION COMPLETE")
    return nil
}

// Usage:
err := OpenSellPosition(sugar, "EURUSD")
```

---

## 🔗 Related Methods

**🍬 Better alternatives with SL/TP:**

* `SellMarketWithSLTP()` - SELL with SL/TP prices
* `SellMarketWithPips()` - SELL with SL/TP distances in points ⭐ **RECOMMENDED**

**🍬 Other market orders:**

* `BuyMarket()` - BUY at market

**🍬 Position management:**

* `ModifyPositionSLTP()` - Add SL/TP after opening
* `ClosePosition()` - Close the position

---

## ⚠️ Common Pitfalls

### 1) Confusing SELL SL/TP direction

```go
// ❌ WRONG - SL/TP in wrong direction for SELL
bid := 1.08500
point := 0.00001
sl := bid - 50*point  // WRONG! SL should be ABOVE for SELL
tp := bid + 100*point // WRONG! TP should be BELOW for SELL

// ✅ CORRECT - use CalculateSLTP (distances in points)
sl, tp, _ = sugar.CalculateSLTP("EURUSD", "SELL", bid, 50, 100)
// SELL: SL is ABOVE entry, TP is BELOW entry
```

### 2) Trading without Stop Loss

```go
// ❌ DANGEROUS - No stop loss!
sugar.SellMarket("EURUSD", 0.1)

// ✅ CORRECT - Use version with SL/TP
sugar.SellMarketWithPips("EURUSD", 0.1, 50, 100)
```

---

## 💎 Pro Tips

1. **SELL = short = profit when price drops**

2. **SL above entry, TP below** - Remember SELL direction

3. **Use `SellMarketWithPips()`** - Much safer than bare `SellMarket()`

4. **Calculate lot size** - Always use risk-based sizing

5. **Validate before selling** - Check margin and spread

---

## 🚨 Production Recommendation

**DON'T use `SellMarket()` in production!**

Instead:
```go
// ✅ RECOMMENDED
ticket, _ := sugar.SellMarketWithPips("EURUSD", 0.1, 50, 100)

// OR with full risk management
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, _ := sugar.SellMarketWithPips("EURUSD", lotSize, 50, 100)
```

---

**See also:** [`SellMarketWithPips.md`](../11.%20Trading_Helpers/SellMarketWithPips.md), [`SellMarketWithSLTP.md`](../5.%20Trading_SLTP/SellMarketWithSLTP.md), [`BuyMarket.md`](BuyMarket.md)
