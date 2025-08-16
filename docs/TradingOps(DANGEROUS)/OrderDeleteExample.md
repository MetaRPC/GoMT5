# 🗑️ Deleting a Pending Order by Ticket (Example)

> **Request:** delete/cancel an existing **pending** order (limit/stop/stop‑limit) by its `ticket`.

---

### Code Example — with pre‑check 🔎

```go
// Suppose you know the ticket, but want to verify it exists and is pending:
func deletePendingIfExists(ctx context.Context, svc *MT5Service, ticket uint64) {
    // 1) Optional: try to fetch historical/active; broker APIs differ.
    //    If you have a direct "get pending by ticket" use it. Here we do a simple check via OpenedOrders.
    data, err := svc.account.OpenedOrders(ctx)
    if err != nil {
        log.Printf("❌ OpenedOrders error: %v", err)
        return
    }
    isPending := false
    for _, o := range data.GetOpenedOrders() {
        if o.GetTicket() == ticket {
            // Consider types: *_LIMIT, *_STOP, *_STOP_LIMIT as pending.
            t := o.GetType()
            if t == pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_LIMIT ||
               t == pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_LIMIT ||
               t == pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_STOP ||
               t == pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_STOP ||
               t == pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_STOP_LIMIT ||
               t == pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_STOP_LIMIT {
                isPending = true
                break
            }
        }
    }
    if !isPending {
        log.Printf("ℹ️ Ticket %d is not a pending order (or not found)", ticket)
        return
    }

    // 2) Delete the pending order
    res, err := svc.account.DeleteOrder(ctx, ticket)
    if err != nil {
        log.Printf("❌ DeleteOrder error: %v", err)
        return
    }
    fmt.Printf("✅ Pending order deleted. CloseMode: %s | Code: %d (%s/%s)\n",
        res.GetCloseMode().String(),
        res.GetReturnedCode(),
        res.GetReturnedStringCode(),
        res.GetReturnedCodeDescription(),
    )
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowOrderDeleteExample(ctx context.Context, ticket uint64)
```

---

## 🔽 Input

Underlying RPC `DeleteOrder` accepts:

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Timeout/cancel control.                    |
| `ticket`  | `uint64`          | yes      | Ticket of the **pending** order to delete. |

---

## ⬆️ Output

Returns **`OrderSendData`** (close/delete result):

| Field                     | Type     | Description                 |
| ------------------------- | -------- | --------------------------- |
| `CloseMode`               | `enum`   | Server’s delete/close mode. |
| `ReturnedCode`            | `uint32` | Numeric result code.        |
| `ReturnedStringCode`      | `string` | Short code.                 |
| `ReturnedCodeDescription` | `string` | Human‑readable description. |

---

## 🎯 Purpose

* Cancel a pending order that is no longer needed (e.g., invalidated setup).

---

## ⚠️ Notes & Tips

* This call is for **pending** orders only. To exit a live position, use `OrderClose`.
* Race conditions: the order may fill between your check and delete. Always handle server codes gracefully.
* Permissions: some symbols can be in `CLOSE_ONLY` or disallow deletes during certain sessions.
