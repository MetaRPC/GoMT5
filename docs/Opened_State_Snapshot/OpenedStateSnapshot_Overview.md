# Opened State Snapshot — Overview

This section groups together methods for retrieving the **current trading state**: open orders, their tickets, and active positions. It’s perfect for building dashboards and health checks before trading.

---

## 📂 Methods in this Section

* [OpenedOrders.md](OpenedOrders.md)
  Get a list of **currently open orders** (pending orders and/or active orders depending on server semantics).

* [OpenedOrderTickets.md](OpenedOrderTickets.md)
  Retrieve **ticket IDs** of open orders for later operations (modify/close/delete).

* [ShowPositions.md](ShowPositions.md)
  List **active positions** (symbol, volume, price open, current P/L, etc.).

* [HasOpenPosition.md](HasOpenPosition.md)
  Quick boolean check: **is there any open position** for a given symbol?

---

## 🔎 Quick Snapshot Example (read‑only)

```go
// One-shot snapshot of the current state
svc.ShowOpenedOrders(ctx)
svc.ShowOpenedOrderTickets(ctx)
svc.ShowPositions(ctx)
svc.ShowHasOpenPosition(ctx, selectedSymbol)
```

* Safe to run: these methods **do not execute trades**.
* Use the printed tickets/positions to decide next actions (e.g., modify/close in TradingOps).

---

## ✅ Best Practices

1. Run this block **before** any trading action to understand the current state.
2. Keep output logs — they help to find the right **ticket IDs** for later operations.
3. Combine with **quotes** (from *Quick Account & Market Info*) to see live prices next to positions.
4. For real-time monitoring, use **Streaming** methods (tickets/profits) as a complement.

---

## 🎯 Purpose

* Build a lightweight **status panel** of your account.
* Prepare data (tickets/positions) for subsequent trading operations.
* Perform safe diagnostics without sending any orders.

---

👉 Use this overview as a **map**, and open each `.md` for full details and code samples.
