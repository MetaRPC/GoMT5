# 🔎 Lookup Historical **Order** by Ticket

> **Request:** fetch details of a **historical order** by its unique `ticket` (no time range needed).

---

### Code Example

```go
// High-level helper (prints a formatted summary):
svc.ShowOrderByTicket(ctx, 123456789)

// Internals (simplified):
o, err := svc.account.HistoryOrderByTicket(ctx, 123456789)
if err != nil {
    log.Printf("❌ HistoryOrderByTicket error: %v", err)
    return
}
if o == nil {
    fmt.Printf("⚠️ Historical order %d not found\n", 123456789)
    return
}

vol  := o.GetVolumeInitial()
open := o.GetPriceOpen()
last := o.GetPriceCurrent()
fmt.Printf("📜 Order #%d | %s | VolumeInitial: %.2f | PriceOpen: %.5f | LastPrice: %.5f",
    o.GetTicket(), o.GetSymbol(), vol, open, last)
if ts := o.GetDoneTime(); ts != nil {
    fmt.Printf(" | Done: %s", ts.AsTime().Format("2006-01-02 15:04:05"))
}
fmt.Println()
```

---

### Method Signature

```go
func (s *MT5Service) ShowOrderByTicket(ctx context.Context, ticket uint64)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                             |
| --------- | ----------------- | -------- | --------------------------------------- |
| `ctx`     | `context.Context` | yes      | Timeout/cancel control.                 |
| `ticket`  | `uint64`          | yes      | Order ticket to look up in **history**. |

> No `from/to` range is required for lookup by ticket.

---

## ⬆️ Output

Prints selected fields from the returned historical order:

| Field           | Type     | Description                              |
| --------------- | -------- | ---------------------------------------- |
| `Ticket`        | `uint64` | Order ticket ID.                         |
| `Symbol`        | `string` | Instrument.                              |
| `Type`          | `enum`   | Order type (Buy/Sell/Limit/Stop/…).      |
| `VolumeInitial` | `double` | Requested volume.                        |
| `PriceOpen`     | `double` | Requested/open price.                    |
| `PriceCurrent`  | `double` | Last price stored for the order.         |
| `TimeSetup`     | `time`   | When order was placed.                   |
| `TimeDone`      | `time`   | When order was closed/cancelled/expired. |

---

## 🎯 Purpose

* Pinpoint an exact order by its ticket to audit behavior and lifecycle.
* Useful in support tickets, logs correlation, and reconciliation tasks.

---

## 🧩 Notes & Tips

* **Orders vs Deals:** this returns **order** metadata; if you need the executed fill(s), use `ShowDealByTicket`.
* If the ticket is very recent and still **open**, it may not appear in history. Check `OpenedOrders` first if needed.
* Timestamps are server timestamps; convert with `t.In(userTZ)` for display.
