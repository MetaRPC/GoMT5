# Listing Open Order Tickets (IDs Only)

> **Request:** fetch lightweight lists of open ticket IDs (positions + pending orders).

---

### Code Example

```go
// High-level (prints IDs):
svc.ShowOpenedOrderTickets(ctx)

// Low-level (work with IDs directly):
data, err := svc.account.OpenedOrdersTickets(ctx)
if err != nil {
    log.Printf("❌ OpenedOrdersTickets error: %v", err)
    return
}
positionIDs := data.GetPositionTickets()       // []uint64
pendingIDs  := data.GetPendingOrderTickets()   // []uint64
fmt.Println("Positions:", positionIDs)
fmt.Println("Pending:", pendingIDs)
```

---

### Method Signature

```go
func (s *MT5Service) ShowOpenedOrderTickets(ctx context.Context)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC. |

---

## ⬆️ Output

Returns **`OpenedOrdersTicketsData`** with two slices of ticket IDs.

| Field                 | Type       | Description                              |
| --------------------- | ---------- | ---------------------------------------- |
| `PositionTickets`     | `[]uint64` | Tickets of currently open **positions**. |
| `PendingOrderTickets` | `[]uint64` | Tickets of active **pending orders**.    |

> This call is intentionally lightweight compared to `OpenedOrders`, useful when you only need identifiers.

---

## 🎯 Purpose

* Quickly obtain lists of IDs to drive follow-up calls (e.g., `HistoryOrderByTicket`, `OrderModify`, `OrderClose`).
* Build compact UIs or health checks without downloading full order payloads.

---

## 🧩 Notes & Tips

* An empty list is normal when there are no active positions/pending orders.
* Combine with `OpenedOrders` if you need full details beyond IDs.
* Tickets are `uint64` — store/print them without lossy conversions.
