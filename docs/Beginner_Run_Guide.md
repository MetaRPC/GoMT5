# 🚦 Beginner Run Guide — “What do I run and where?”

This page is for a newcomer who already has **config filled** and wants to **start doing things** with MT5 from code.

> You run everything from **`examples/`**:
>
> ```bash
> cd examples
> go run .
> ```
>
> All snippets below go **into `examples/main.go`**, inside `main()` **after** you create `svc := mt5.NewMT5Service(account)` and choose `selectedSymbol`.

---

## ✅ Scenario 0 — Minimal sanity check (safe)

**Goal:** make sure terminal is connected and symbol works.

```go
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, selectedSymbol)
```

**Expected:** balance/equity + live bid/ask in console.

---

## 🧭 Scenario 1 — Snapshot of current state (safe)

**Goal:** see what is open now.

```go
svc.ShowOpenedOrders(ctx)
svc.ShowOpenedOrderTickets(ctx)
svc.ShowPositions(ctx)
svc.ShowHasOpenPosition(ctx, selectedSymbol)
```

**Expected:** lists of orders, tickets, positions, or “no items”.

---

## 🧮 Scenario 2 — Calculations before trading (safe)

**Goal:** estimate margin/profit and validate a future trade.

```go
// 2.1 Margin for a potential BUY
svc.ShowOrderCalcMargin(ctx, selectedSymbol, pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, 0.10, 0)

// 2.2 Profit estimation if price moves from 1.08000 to 1.08350
svc.ShowOrderCalcProfit(ctx, selectedSymbol, pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, 0.10, 1.08000, 1.08350)

// 2.3 Dry‑run check (no trade yet): market BUY 0.10 lots
svc.ShowOrderCheck(
    ctx,
    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL, // market action
    pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,              // BUY
    selectedSymbol,
    0.10,
    0,   // price=0 → market on most servers
    nil, // sl
    nil, // tp
    nil, // deviation
    nil, // magic
    nil, // expiration
)
```

**Expected:** margin/profit numbers and `OrderCheck` retcode/comment.

---

## 🧪 Scenario 3 — Place a *pending* order with expiration (real!)

**Goal:** safely test order placement without instant execution.

> ⚠️ This is a **real trading action** (even on demo). Keep volume small.

```go
exp := timestamppb.New(time.Now().Add(24 * time.Hour)) // auto-cancel in 24h

// Example: Buy Limit at 1.07500, 0.10 lots
svc.PlaceBuyLimit(ctx, selectedSymbol, 0.10, 1.07500, nil, nil, exp)

// See result in console, then you may delete it:
// svc.ShowOrderDeleteExample(ctx, YOUR_ORDER_TICKET)
```

**Notes:** for Sell/Stop/StopLimit use corresponding helpers (`PlaceSellLimit`, `PlaceBuyStop`, `PlaceStopLimit`, …).

---

## ⚡ Scenario 4 — Market trade (real!)

**Goal:** execute now at market price.

> ⚠️ Real trade! Uncomment only when you are **ready**.

```go
// svc.BuyMarket(ctx, selectedSymbol, 0.10, nil, nil)
// svc.SellMarket(ctx, selectedSymbol, 0.10, nil, nil)
```

**Next steps:** check `svc.ShowPositions(ctx)` and close/modify if needed:

```go
// svc.ShowOrderCloseExample(ctx, 123456789)
// svc.ShowPositionClose(ctx, selectedSymbol)
```

---

## 📊 Scenario 5 — History & simple stats (safe)

**Goal:** read what happened earlier.

```go
from := time.Now().AddDate(0, 0, -7)
to   := time.Now()

svc.ShowOrdersHistory(ctx)
svc.ShowDealsCount(ctx, from, to, "") // all symbols
// svc.ShowOrderByTicket(ctx, 123456789)
// svc.ShowDealByTicket(ctx, 987654321)
```

---

## 📡 Scenario 6 — Streaming (safe to read)

**Goal:** watch live updates for a short period.

```go
svc.StreamQuotes(ctx)               // live ticks
// svc.StreamOpenedOrderProfits(ctx) // P/L per open order
// svc.StreamOpenedOrderTickets(ctx) // open tickets
// svc.StreamTradeUpdates(ctx)       // trade events
```

**Tip:** each stream stops by timeout inside the helper (≈30s) or on error.

---

## 🔧 Where exactly to paste?

Open `examples/main.go` and find the block where you already have:

```go
svc := mt5.NewMT5Service(account)
selectedSymbol := cfg.DefaultSymbol // after EnsureSymbolVisible(...)
```

Paste any scenario **right below** this block. You can run scenarios one by one, or combine several — they will execute in order.

---

## 🆘 Common pitfalls

* **“no Go files in …”** — You ran from the wrong folder. Use: `cd examples && go run .`
* **Proxy/internet issues** — set `MT5_PROXY` env var if you’re behind a firewall.
* **Symbol not found** — we try several aliases; still failing? Ask your broker for exact symbol name.
* **Trade rejected** — check `AccountSummary` (free margin), then run `ShowOrderCheck` to see the reason code.

---

## 🧠 Quick copy block (everything safe)

```go
// Minimal safe pack
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, selectedSymbol)
svc.ShowOpenedOrders(ctx)
svc.ShowPositions(ctx)
from := time.Now().AddDate(0, 0, -7); to := time.Now()
svc.ShowDealsCount(ctx, from, to, "")
```

That’s it — uncomment gradually, watch the console, and move from safe checks to real actions when you’re ready. 🚀
