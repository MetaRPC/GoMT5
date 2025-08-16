# History & Simple Statistics — Overview

This section groups together methods for retrieving **historical data** and **basic statistics**.
It is useful for building account analytics, trade logs, and monitoring dashboards.

---

## 📂 Methods in this Section

* [ShowOrdersHistory.md](ShowOrdersHistory.md)
  Get all historical orders in a selected time range.

* [ShowDealsCount.md](ShowDealsCount.md)
  Retrieve the number of executed deals within a time period.

* [OrderByTicket.md](OrderByTicket.md)
  Fetch a specific historical order by its ticket ID.

* [DealByTicket.md](DealByTicket.md)
  Fetch a specific historical deal by its ticket ID.

* [History\_Range(important).md](History_Range%28important%29.md)
  Detailed explanation of how to define time ranges (`from` / `to`) in Go.

---

## 🕒 Time Range Example

All history methods use a **time range** for filtering.
Example: last 7 days.

```go
from := time.Now().AddDate(0, 0, -7) // 7 days ago
to   := time.Now()                   // current moment

svc.ShowOrdersHistory(ctx, from, to)
svc.ShowDealsCount(ctx, from, to, "")
```

* `from` → lower bound (inclusive).
* `to` → upper bound (exclusive).
* Both values are converted into **Unix timestamps** internally.

---

## ✅ Best Practices

1. Always define `from` and `to` — otherwise you risk requesting a huge history.
2. Use small ranges when testing (e.g., 1–3 days).
3. For **performance**: prefer `DealsCount` when you only need statistics.
4. For **details**: use `OrderByTicket` or `DealByTicket` for precise lookups.
5. Check broker’s server limits — some servers restrict history depth (e.g., 1 year max).

---

## 🎯 Purpose

The methods in this block allow you to:

* Build trade history reports.
* Generate statistics (win rate, trade count, average profit).
* Backtest simple strategies using real historical executions.
* Diagnose trading activity and validate orders/deals by ticket ID.

---

👉 Use this overview as a **map**, and jump into each `.md` file for full method details.
