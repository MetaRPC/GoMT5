# ShowDealByTicket

> **Request:** retrieve a single **deal record** from history by its deal ticket ID.

---

### Code Example

```go
dealID := uint64(987654321)

err := svc.ShowDealByTicket(ctx, dealID)
if err != nil {
    log.Printf("❌ ShowDealByTicket error: %v", err)
    return
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowDealByTicket(ctx context.Context, ticket uint64) error
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                 |
| --------- | ----------------- | -------- | ------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Context for cancellation / timeout control. |
| `ticket`  | `uint64`          | yes      | Deal ticket ID to fetch.                    |

---

## ⬆️ Output

Prints details of a **deal** (trade execution). Fields (from `DealRecordData`):

| Field        | Type     | Description                            |
| ------------ | -------- | -------------------------------------- |
| `Deal`       | `uint64` | Deal ticket ID.                        |
| `Order`      | `uint64` | Associated order ticket.               |
| `Symbol`     | `string` | Symbol name (e.g., `EURUSD`).          |
| `Volume`     | `double` | Executed volume.                       |
| `Price`      | `double` | Execution price.                       |
| `Commission` | `double` | Applied commission.                    |
| `Swap`       | `double` | Applied swap (overnight fee).          |
| `Profit`     | `double` | Profit or loss from this deal.         |
| `Entry`      | `enum`   | Entry type: `IN`, `OUT`, `INOUT`, etc. |
| `Time`       | `int64`  | Unix timestamp of execution.           |
| `Comment`    | `string` | Broker comment if present.             |

---

## 🎯 Purpose

* Fetch a **single executed deal** from history.
* Use when you know the exact deal ticket and need full info (audit, debugging, reconciliation).
* Complements `ShowOrderByTicket`, but operates on **deals** instead of orders.

---

## 🧩 Notes & Tips

* Deals are immutable: once executed, their record does not change.
* Not all deals have corresponding orders (e.g., balance operations).
* Always check `Profit` and `Commission` fields to compute net effect.
* Time is UTC — convert if you need local reporting.
