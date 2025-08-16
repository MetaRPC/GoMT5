# ShowOrdersHistory

> **Request:** retrieve and display the list of past orders within a specified time range.

---

### Code Example

```go
from := time.Now().AddDate(0, 0, -7) // 7 days ago
to   := time.Now()                   // now

svc.ShowOrdersHistory(ctx, from, to)
```

Internally this calls the account method to fetch all closed/cancelled/expired orders that fall into the requested time interval.

---

### Method Signature

```go
func (s *MT5Service) ShowOrdersHistory(ctx context.Context, from time.Time, to time.Time)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                      |
| --------- | ----------------- | -------- | -------------------------------- |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancel.         |
| `from`    | `time.Time`       | yes      | Start of time range (inclusive). |
| `to`      | `time.Time`       | yes      | End of time range (exclusive).   |

📎 For practical usage and examples on how to construct `from` and `to` (week, today, month-to-date, etc.) see 👉 **[History Range — How To Use](68a0bfff30648191b2b6511a2691647b)**

---

## ⬆️ Output

Prints the retrieved orders to console. Each order record may include:

| Field           | Type   | Description                                     |
| --------------- | ------ | ----------------------------------------------- |
| `Order`         | uint64 | Ticket ID.                                      |
| `Symbol`        | string | Instrument traded.                              |
| `VolumeInitial` | double | Requested volume.                               |
| `VolumeCurrent` | double | Remaining volume (if partially filled).         |
| `PriceOpen`     | double | Requested/open price.                           |
| `PriceCurrent`  | double | Current market price at close.                  |
| `Type`          | enum   | Order type (Buy, Sell, Limit, Stop, etc.).      |
| `State`         | enum   | Final order state (filled, cancelled, expired). |
| `Reason`        | enum   | Why the order was closed/cancelled.             |
| `TimeSetup`     | time   | When order was placed.                          |
| `TimeDone`      | time   | When order was closed/expired.                  |

---

## 🎯 Purpose

* Inspect what trades/orders were placed and how they evolved.
* Debug past executions (why order was cancelled or rejected).
* Collect statistics (counts, volumes, durations).

---

## 🧩 Notes & Tips

* MT5 differentiates **Orders** (instructions) from **Deals** (executed fills). Orders may result in zero, one, or many deals.
* History can be large → always restrict `from/to` range.
* Use `ShowDealsCount` or per-ticket lookups (`ShowOrderByTicket`, `ShowDealByTicket`) for more targeted queries.
