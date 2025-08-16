#  Placing a Pending Sell Limit Order

> **Request:** place a **SELL LIMIT** pending order at a specified price (optionally with SL/TP and expiration).

---

### Code Example

```go
// High-level helper (prints a human-readable result):
svc.PlaceSellLimit(ctx, "EURUSD", 0.10, 1.09500, nil, nil, timestamppb.New(time.Now().Add(24*time.Hour)))

// Internals (simplified):
vol := 0.10
price := 1.09500
slip := int32(10)                       // points
comment := "SellLimit"
magic32 := int32(123456)
exp := timestamppb.New(time.Now().Add(24*time.Hour))

// 1) Optional dry-run check (the helper does it for you):
checkReq := &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
        OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_LIMIT,
        Symbol:     "EURUSD",
        Volume:     vol,
        Price:      price, // absolute entry price for pending
        Deviation:  10,    // points
        TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
        TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
            if exp != nil { return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED }
            return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
        }(),
        Expiration:               exp,
        ExpertAdvisorMagicNumber: 123456,
        Comment:                  "PlaceSellLimit helper",
    },
}
if chk, err := svc.account.OrderCheck(ctx, checkReq); err == nil {
    if r := chk.GetMqlTradeCheckResult(); r != nil {
        fmt.Printf("ℹ️ Check SELL_LIMIT: code=%d comment=%q\n", r.GetReturnedCode(), r.GetComment())
    }
}

// 2) Actual send:
res, err := svc.account.OrderSend(
    ctx,
    "EURUSD",
    pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_LIMIT,
    vol,
    &price, &slip,
    nil, nil,         // sl, tp (absolute prices) — set if needed
    &comment, &magic32,
    exp,
)
if err != nil {
    log.Printf("❌ OrderSend(SELL_LIMIT): %v", err)
    return
}
// Read result
if ord := res.GetOrder(); ord != 0 {
    fmt.Printf("✅ SELL_LIMIT placed: order=%d @ %.5f\n", ord, res.GetPrice())
} else if deal := res.GetDeal(); deal != 0 {
    fmt.Printf("✅ SELL_LIMIT executed immediately: deal=%d @ %.5f\n", deal, res.GetPrice())
} else {
    fmt.Printf("⚠️ SELL_LIMIT sent @ %.5f (no ticket in response)\n", res.GetPrice())
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) PlaceSellLimit(ctx context.Context, symbol string, volume, price float64, sl, tp *float64, exp *timestamppb.Timestamp)
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
* `ENUM_ORDER_TYPE_TF`: `ORDER_TYPE_TF_SELL_LIMIT` (check stage).
* `TMT5_ENUM_ORDER_TYPE`: `TMT5_ORDER_TYPE_SELL_LIMIT` (send stage).
* `MRPC_ENUM_ORDER_TYPE_FILLING`: typically `ORDER_FILLING_FOK`.
* `MRPC_ENUM_ORDER_TYPE_TIME`: `ORDER_TIME_GTC` or `ORDER_TIME_SPECIFIED` when `exp` is provided.

---

## 🎯 Purpose

* Enter short only when price bounces up to your specified level.
* Attach SL/TP immediately and control **lifetime** via `exp`.

---

## 🧩 Notes & Tips

* Ensure `price` is **above current market** for SELL\_LIMIT, otherwise broker may reject or fill immediately depending on rules.
* For a short, sanity-check: **SL above price**, **TP below price**.
* `sl`/`tp` are absolute prices; derive from pips/ticks using `Digits` and `TickValueWithSize`.
* If you pass `exp`, mind server vs local timezone; build the time via `timestamppb.New(timeX)`.
* The helper performs an `OrderCheck` first (prints code/comment) — use that output to troubleshoot rejects.
