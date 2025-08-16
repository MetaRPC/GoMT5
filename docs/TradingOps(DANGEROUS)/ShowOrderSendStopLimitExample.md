# Placing a STOP‑LIMIT Order (Example)

> **Request:** place a **Stop‑Limit** pending order (buy or sell) with trigger and limit prices.

---

### Code Example

```go
// High-level helper (prints a human‑readable result):
svc.ShowOrderSendStopLimitExample(ctx, "EURUSD", true /*isBuy*/, 1.09100 /*trigger*/, 1.09120 /*limit*/)

// What it does internally (simplified):
volume := 0.10
slippage := int32(10)         // points
comment := "SLimit from service"
magic := int32(98765)
exp := timestamppb.New(time.Now().Add(24 * time.Hour))

res, err := svc.account.OrderSendStopLimit(
    ctx,
    "EURUSD", // symbol
    true,      // isBuy (true=BUY_STOP_LIMIT, false=SELL_STOP_LIMIT)
    volume,
    1.09100,   // trigger price (stop)
    1.09120,   // limit price
    &slippage, // max slippage (points) for the *limit* execution
    nil,       // sl (absolute price) or nil
    nil,       // tp (absolute price) or nil
    &comment,  // user/EA comment
    &magic,    // EA magic number
    exp,       // expiration (optional; required if using time-specified policies)
)
if err != nil {
    log.Printf("❌ OrderSendStopLimit error: %v", err)
    return
}

if ord := res.GetOrder(); ord != 0 {
    fmt.Printf("✅ STOP_LIMIT placed. Order:%d Trigger:%.5f Limit:%.5f Code:%d\n",
        ord, 1.09100, 1.09120, res.GetReturnedCode())
    return
}
if deal := res.GetDeal(); deal != 0 {
    fmt.Printf("✅ STOP_LIMIT executed immediately. Deal:%d Price:%.5f Code:%d\n",
        deal, res.GetPrice(), res.GetReturnedCode())
    return
}
fmt.Printf("⚠️ STOP_LIMIT response without ticket | Price: %.5f | Code: %d\n",
    res.GetPrice(), res.GetReturnedCode())
```

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowOrderSendStopLimitExample(ctx context.Context, symbol string, isBuy bool, trigger, limit float64)
```

---

## 🔽 Input

High-level helper takes `symbol`, `isBuy`, `trigger`, `limit`; the underlying RPC accepts the full set below.

### Underlying `OrderSendStopLimit` parameters

| Parameter    | Type                         | Required | Description                                             |
| ------------ | ---------------------------- | -------- | ------------------------------------------------------- |
| `ctx`        | `context.Context`            | yes      | Timeout/cancel control.                                 |
| `symbol`     | `string`                     | yes      | Symbol name.                                            |
| `isBuy`      | `bool`                       | yes      | `true` → BUY\_STOP\_LIMIT, `false` → SELL\_STOP\_LIMIT. |
| `volume`     | `float64`                    | yes      | Lots (respect `VolumeMin/Max/Step`).                    |
| `trigger`    | `float64`                    | yes      | **Stop** price that arms the order.                     |
| `limit`      | `float64`                    | yes      | **Limit** price to execute after trigger.               |
| `slippage`   | `*int32`                     | no       | Max slippage (points) for the limit execution.          |
| `sl`         | `*float64`                   | no       | Stop Loss **absolute** price or `nil`.                  |
| `tp`         | `*float64`                   | no       | Take Profit **absolute** price or `nil`.                |
| `comment`    | `*string`                    | no       | User/EA comment.                                        |
| `magic`      | `*int32`                     | no       | EA magic number.                                        |
| `expiration` | `*google.protobuf.Timestamp` | no       | Expiration; set if you need time‑limited pending.       |

---

## ⬆️ Output

RPC returns **`OrderSendData`** (same container as `OrderSend`). Key fields:

| Field                     | Type     | Description                                          |
| ------------------------- | -------- | ---------------------------------------------------- |
| `Order`                   | `uint64` | Pending order ticket (if placed).                    |
| `Deal`                    | `uint64` | Deal ticket (if triggered and executed immediately). |
| `Price`                   | `double` | Server price value from processing.                  |
| `ReturnedCode`            | `uint32` | Numeric result code.                                 |
| `ReturnedStringCode`      | `string` | Short string code.                                   |
| `ReturnedCodeDescription` | `string` | Human‑readable description.                          |

---

## Direction & Types

Although the helper uses a boolean, internally this maps to **TF order types**:

* `isBuy=true`  → `BUY_STOP_LIMIT`
* `isBuy=false` → `SELL_STOP_LIMIT`

These correspond to the usual `ENUM_ORDER_TYPE_TF` values used in checks and to `TMT5_ENUM_ORDER_TYPE` used in low‑level send calls.

---

## 🎯 Purpose

* Place a pending order that triggers at one price and must execute at (or better than) a second price.
* Useful for breakouts with controlled entry price.

---

## 🧩 Notes & Tips

* **Trigger vs Limit:** the order arms at `trigger`, then submits a *limit* at `limit`. For **BUY\_STOP\_LIMIT**, typically `limit ≥ trigger`; for **SELL\_STOP\_LIMIT**, typically `limit ≤ trigger`. Brokers can reject invalid relations.
* **Immediate execution:** In fast markets, the trigger may fire instantly and your limit may fill immediately — that’s the “Deal” branch above.
* **SL/TP are absolute prices** (not offsets). Compute from pips/ticks using `Digits`.
* **Slippage is in points.** Map carefully if your UI works in pips (5‑digit FX: 1 pip = 10 points).
* Set `expiration` when you want the pending to auto‑cancel after some time; otherwise time‑in‑force defaults apply per broker.
