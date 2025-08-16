# Placing a Pending Buy Limit Order

> **Request:** place a **BUY LIMIT** pending order at a specified price (optionally with SL/TP and expiration).

---

### Code Example

```go
// High-level helper (prints a human-readable result):
svc.PlaceBuyLimit(ctx, "EURUSD", 0.10, 1.07500, nil, nil, timestamppb.New(time.Now().Add(24*time.Hour)))

// Internals (simplified):
vol := 0.10
price := 1.07500
slip := int32(10)                       // points
comment := "BuyLimit"
magic32 := int32(123456)
exp := timestamppb.New(time.Now().Add(24*time.Hour))

// 1) Optional dry-run check (the helper does it for you):
checkReq := &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_LIMIT,
        Symbol:    "EURUSD",
        Volume:    vol,
        Price:     price, // absolute entry price for pending
        Deviation: 10,    // points
        TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
        TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
            if exp != nil { return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED }
            return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
        }(),
        Expiration:               exp,
        ExpertAdvisorMagicNumber: 123456,
        Comment:                  "PlaceBuyLimit helper",
    },
}
if chk, err := svc.account.OrderCheck(ctx, checkReq); err == nil {
    if r := chk.GetMqlTradeCheckResult(); r != nil {
        fmt.Printf("ℹ️ Check BUY_LIMIT: code=%d comment=%q\n", r.GetReturnedCode(), r.GetComment())
    }
}

// 2) Actual send:
res, err := svc.account.OrderSend(
    ctx,
    "EURUSD",
    pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT,
    vol,
    &price, &slip,
    nil, nil,         // sl, tp (absolute prices) — set if needed
    &comment, &magic32,
    exp,
)
if err != nil {
    log.Printf("❌ OrderSend(BUY_LIMIT): %v", err)
    return
}
// Read result
if ord := res.GetOrder(); ord != 0 {
    fmt.Printf("✅ BUY_LIMIT placed: order=%d @ %.5f\n", ord, res.GetPrice())
} else if deal := res.GetDeal(); deal != 0 {
    fmt.Printf("✅ BUY_LIMIT executed immediately: deal=%d @ %.5f\n", deal, res.GetPrice())
} else {
    fmt.Printf("⚠️ BUY_LIMIT sent @ %.5f (no ticket in response)\n", res.GetPrice())
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) PlaceBuyLimit(ctx context.Context, symbol string, volume, price float64, sl, tp *float64, exp *timestamppb.Timestamp)
```

---

## 🔽 Input

| Parameter | Type                         | Required | Description                                                   |
| --------- | ---------------------------- | -------- | ------------------------------------------------------------- |
| `ctx`     | `context.Context`            | yes      | Timeout/cancel control.                                       |
| `symbol`  | `string`                     | yes      | Symbol name.                                                  |
| `volume`  | `float64`                    | yes      | Lots (respect `VolumeMin/Max/Step`).                          |
| `price`   | `float64`                    | yes      | **Absolute entry price** for the pending order.               |
| `sl`      | `*float64`                   | no       | Stop Loss **absolute** price.                                 |
| `tp`      | `*float64`                   | no       | Take Profit **absolute** price.                               |
| `exp`     | `*google.protobuf.Timestamp` | no       | Expiration time; if set, time-in-force becomes `*_SPECIFIED`. |

---

## ⬆️ Output

Returns **`OrderSendData`** (same as other send calls).

| Field                     | Type     | Description                                        |
| ------------------------- | -------- | -------------------------------------------------- |
| `Order`                   | `uint64` | Pending order ticket (expected for a limit).       |
| `Deal`                    | `uint64` | Deal ticket if order triggered & filled instantly. |
| `Price`                   | `double` | Price acknowledged by server.                      |
| `ReturnedCode`            | `uint32` | Numeric result code.                               |
| `ReturnedStringCode`      | `string` | Short code.                                        |
| `ReturnedCodeDescription` | `string` | Human description.                                 |

---

## Enums used

* `MRPC_ENUM_TRADE_REQUEST_ACTIONS`: `TRADE_ACTION_PENDING` (check stage).
* `ENUM_ORDER_TYPE_TF`: `ORDER_TYPE_TF_BUY_LIMIT` (check stage).
* `TMT5_ENUM_ORDER_TYPE`: `TMT5_ORDER_TYPE_BUY_LIMIT` (send stage).
* `MRPC_ENUM_ORDER_TYPE_FILLING`: typically `ORDER_FILLING_FOK`.
* `MRPC_ENUM_ORDER_TYPE_TIME`: `ORDER_TIME_GTC` or `ORDER_TIME_SPECIFIED` when `exp` is provided.

---

## 🎯 Purpose

* Enter long only when price dips to your specified level.
* Attach SL/TP immediately and control **lifetime** via `exp`.

---

## 🧩 Notes & Tips

* Ensure `price` is **below current market** for BUY\_LIMIT, otherwise broker will reject or execute immediately (depending on rules).
* `sl`/`tp` are absolute prices; derive from pips/ticks using `Digits` and `TickValueWithSize`.
* If you pass `exp`, make sure server time zone differences are considered; always use `timestamppb.New(timeX)` to avoid serialization issues.
* The helper performs an `OrderCheck` first (prints code/comment) — use that output to troubleshoot rejects.
