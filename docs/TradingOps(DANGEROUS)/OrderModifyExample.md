# ✏️ Modifying an Existing Order (Example)

> **Request:** modify fields of an existing ticket — price (for pendings), SL/TP, and/or expiration.

---

### Code Example

```go
// High-level helper (prints a simple status):
svc.ShowOrderModifyExample(ctx, 123456789)

// Internals (simplified):
// For demo: change only SL/TP (typical for market positions or pendings)
newSL := 1.0500
newTP := 1.0900

res, err := svc.account.OrderModify(
    ctx,
    123456789,   // ticket (pending order or position/order ticket depending on broker semantics)
    nil,         // newPrice: set only for PENDING orders to move entry price
    &newSL,      // Stop Loss (absolute price) or nil to keep
    &newTP,      // Take Profit (absolute price) or nil to keep
    nil,         // expiration: for pendings; nil to keep/remove depending on TIF
)
if err != nil {
    log.Printf("❌ OrderModify error: %v", err)
    return
}
if res != nil {
    fmt.Println("✅ Order successfully modified.")
} else {
    fmt.Println("⚠️ Order was NOT modified.")
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowOrderModifyExample(ctx context.Context, ticket uint64)
```

---

## 🔽 Input

Underlying RPC `OrderModify` accepts:

| Parameter    | Type                         | Required | Description                                                      |
| ------------ | ---------------------------- | -------- | ---------------------------------------------------------------- |
| `ctx`        | `context.Context`            | yes      | Timeout/cancel control.                                          |
| `ticket`     | `uint64`                     | yes      | Ticket to modify (pending order / order).                        |
| `newPrice`   | `*float64`                   | no       | **Pending orders only** — new entry price. `nil` → keep current. |
| `newSL`      | `*float64`                   | no       | Stop Loss **absolute** price. `nil` → keep.                      |
| `newTP`      | `*float64`                   | no       | Take Profit **absolute** price. `nil` → keep.                    |
| `expiration` | `*google.protobuf.Timestamp` | no       | New expiration for pendings (use `timestamppb.New(...)`).        |

> For **market positions**, you typically only change SL/TP. Moving entry price applies to **pending** orders.

---

## ⬆️ Output

Returns an implementation-specific confirmation (non-nil on success). The helper prints one of:

* `✅ Order successfully modified.`
* `⚠️ Order was NOT modified.`

---

## 🎯 Purpose

* Adjust protective levels (SL/TP) after entry.
* Reprice or extend lifetime of **pending** orders.

---

## 🧩 Notes & Tips

* SL/TP are **absolute prices** — derive them from pips using symbol `Digits` and current quote.
* To change expiration, build `exp := timestamppb.New(timeX)` and pass it as the last argument.
* Brokers validate distances to current price; too-tight SL/TP or invalid `newPrice` will be rejected.
* Always log/inspect the broker’s returned code/comment if available for troubleshooting.
