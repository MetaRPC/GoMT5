# 🚦 Beginner Run Guide for GoMT5

This guide is designed for **beginners** who already have their `config.json` set up with login, password, server, and default symbol. It explains **what to uncomment** in `examples/main.go` and what happens when you run the code.

---

## ⚡ How to Run

Open your terminal in the project root:

```bash
cd examples
go run .
```

This will use your `examples/main.go` file as the entry point.

---

## 🧪 Safe First Steps

These operations **do not change anything** on your account. They are safe to test.

In `examples/main.go`, uncomment one or more of the following:

```go
// Show account summary
svc.ShowAccountSummary(ctx)

// Show all available symbols
svc.ShowAllSymbols(ctx)

// Get live quote for EURUSD
svc.ShowQuote(ctx, "EURUSD")

// Stream continuous quotes
svc.StreamQuotes(ctx)
```

Run the app again and watch the output.

---

## 📊 Getting Data

Useful methods to retrieve account or market information:

```go
// Print history of trades (last 7 days by default)
svc.ShowOrdersHistory(ctx)

// Print currently open orders
svc.ShowOpenedOrders(ctx)

// Show full symbol parameters (like contract size, digits, etc.)
svc.ShowSymbolParams(ctx, "EURUSD")
```

---

## ⚠️ Trading Operations (Be Careful)

These will actually place, modify, or close orders (even on demo). Always test on a demo account first.

```go
// Open a sample Buy order
svc.ShowOrderSendExample(ctx, "EURUSD")

// Modify order SL/TP (requires a real ticket)
svc.ShowOrderModifyExample(ctx, 123456789)

// Close order by ticket
svc.ShowOrderCloseExample(ctx, 123456789)
```

---

## 🎬 Combo Scenarios

### 1. Open → Check → Close

```go
// Step 1: Open Buy Market order
svc.ShowOrderSendExample(ctx, "EURUSD")

// Step 2: List opened orders
svc.ShowOpenedOrders(ctx)

// Step 3: Close by ticket (replace with real ticket number from output)
svc.ShowOrderCloseExample(ctx, 123456789)
```

---

### 2. Place Pending → Wait → Delete

```go
// Step 1: Place pending order (Buy Limit)
svc.ShowOrderSendExample(ctx, "EURUSD") // adjust inside method to pending type if needed

// Step 2: Check active pending orders
svc.ShowOpenedOrders(ctx)

// Step 3: Delete the pending order by ticket
svc.ShowOrderDeleteExample(ctx, 123456789)
```

---

### 3. Market Data Dashboard

```go
// Show account summary
svc.ShowAccountSummary(ctx)

// Stream real-time quotes
svc.StreamQuotes(ctx)

// Print live profits of open orders
svc.StreamOpenedOrderProfits(ctx)
```

---

## 🧠 Tips for Beginners

* Always start with **safe methods** before trying trading functions.
* Replace `123456789` with the **real order ticket number** you see in the output.
* Keep `config.json` secure — it contains your login details.
* Use **demo accounts** until you are fully confident.

---

This way, even a beginner can follow step by step, uncomment the right methods, and see immediate results without being lost.

---

## 🔗 Combo Scenario A — Place Buy Limit → monitor → delete (real!)

**Goal:** create a pending order with expiration, watch it in streams for \~30s, then delete it by ticket.

```go
// 1) Place Buy Limit (auto-cancel in 24h)
exp := timestamppb.New(time.Now().Add(24 * time.Hour))
svc.PlaceBuyLimit(ctx, selectedSymbol, 0.10, 1.07500, nil, nil, exp)

// 2) Watch tickets/profits for ~30s (internally stop by timeout)
svc.StreamOpenedOrderTickets(ctx)
svc.StreamOpenedOrderProfits(ctx)

// 3) Get the order ticket from console output (or from your terminal log)
//    Then delete it:
// svc.ShowOrderDeleteExample(ctx, YOUR_ORDER_TICKET)
```

**Tip:** To make finding the ticket easier, use a unique comment when placing orders (see trading helpers with `comment`).

---

## 🔗 Combo Scenario B — Market Buy → show position → modify SL/TP → close (real!)

**Goal:** open a position at market, ensure it exists, optionally adjust, then close.

```go
// 1) Execute market BUY (⚠️ real trade)
// svc.BuyMarket(ctx, selectedSymbol, 0.10, nil, nil)

// 2) Confirm in console
svc.ShowPositions(ctx)            // should show an open position for selectedSymbol
svc.ShowHasOpenPosition(ctx, selectedSymbol)

// 3) Modify position SL/TP (needs position ticket) — optional
//    If you know the ticket, call:
// svc.ShowPositionModify(ctx, /*ticket=*/ 123456789, /*newSL=*/ nil, /*newTP=*/ nil)
//    (Without a helper that returns the ticket programmatically, take it from ShowPositions output.)

// 4) Close the position by symbol
// svc.ShowPositionClose(ctx, selectedSymbol)
```

**Note:** If you prefer closing the *order* itself, use `ShowOrderCloseExample(ctx, ticket)` with the exact order ticket.

---

## 🔗 Combo Scenario C — Intraday pending workflow (place → auto-expire)

**Goal:** create orders that auto-cancel by end of day if not filled.

```go
// Expire at local midnight
endOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 0, 0, time.Local)
exp := timestamppb.New(endOfDay)

// Two brackets: BUY_STOP above market and SELL_STOP below
svc.PlaceBuyStop(ctx,  selectedSymbol, 0.10, /*trigger=*/ 1.09200, nil, nil, exp)
svc.PlaceSellStop(ctx, selectedSymbol, 0.10, /*trigger=*/ 1.07800, nil, nil, exp)

// Optional: watch for ~30s
svc.StreamTradeUpdates(ctx)
```

**Cleanup:** If needed, delete by ticket later via `ShowOrderDeleteExample`.

---

## 🔗 Combo Scenario D — Quick read-only dashboard (safe)

**Goal:** one-shot snapshot for monitoring.

```go
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, selectedSymbol)
svc.ShowOpenedOrders(ctx)
svc.ShowPositions(ctx)
from := time.Now().AddDate(0, 0, -7); to := time.Now()
svc.ShowDealsCount(ctx, from, to, "")
```

---

## 🧷 Copy‑Paste Starters

* Need expiration doc? See **SetOrderExpiration.md** in `docs/TradingOps(DANGEROUS)/`.
* Need time range patterns? See **History\_Range(important).md** in `docs/History_And_SimpleStatistics/`.

These combos keep changes minimal and mirror the helpers you already have, so a newcomer can un‑comment a block, run it, and understand the flow step by step. ✅
