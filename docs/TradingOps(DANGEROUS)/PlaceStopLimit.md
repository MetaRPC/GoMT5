# 🎯 Placing a Stop‑Limit Pending Order (Buy/Sell)

> **Request:** place a **STOP‑LIMIT** order with a trigger (stop) price and execution (limit) price. Works for both directions via `isBuy`.

---

### Code Example

```go
// High-level helper (prints a human-readable result):
svc.PlaceStopLimit(ctx, "EURUSD", true /*isBuy*/, 0.10, 1.09100 /*trigger*/, 1.09120 /*limit*/, nil, nil, timestamppb.New(time.Now().Add(24*time.Hour)))

// Internals (simplified):
vol := 0.10
trigger := 1.09100
limit := 1.09120
slip := int32(10)
comment := "StopLimit"
magic32 := int32(123456)
exp := timestamppb.New(time.Now().Add(24*time.Hour))

// 1) Optional dry-run check (the helper does it for you):
orderTypeTF := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_STOP_LIMIT
if !isBuy { orderTypeTF = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_STOP_LIMIT }
checkReq := &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
        OrderType:  orderTypeTF,
        Symbol:     "EURUSD",
        Volume:     vol,
        Price:      trigger, // stop (trigger)
        Deviation:  10,      // points for the limit fill
        TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
        TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
            if exp != nil { return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED }
            return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
        }(),
        Expiration:               exp,
        ExpertAdvisorMagicNumber: 123456,
        Comment:                  "PlaceStopLimit helper",
    },
}
if chk, err := svc.account.OrderCheck(ctx, checkReq); err == nil {
    if r := chk.GetMqlTradeCheckResult(); r != nil {
        fmt.Printf("ℹ️ Check STOP_LIMIT: code=%d comment=%q\n", r.GetReturnedCode(), r.GetComment())
    }
}

// 2) Actual send (OrderSendEx: needs both trigger & limit):
res, err := svc.account.OrderSendEx(
    ctx,
    "EURUSD",
    func() pb.TMT5_ENUM_ORDER_TYPE {
        if isBuy { return pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP_LIMIT }
        return pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP_LIMIT
    }(),
    vol,
    &trigger, &slip,
    nil, nil,            // sl, tp
    &comment, &magic32,
    exp,
    &limit,              // the extra limit price specific to Stop‑Limit
)
if err != nil {
    log.Printf("❌ OrderSendEx(STOP_LIMIT): %v", err)
    return
}
if ord := res.GetOrder(); ord != 0 {
    fmt.Printf("✅ STOP_LIMIT placed: order=%d | trigger=%.5f | limit=%.5f\n", ord, trigger, limit)
} else if deal := res.GetDeal(); deal != 0 {
    fmt.Printf("✅ STOP_LIMIT executed immediately: deal=%d | price=%.5f\n", deal, res.GetPrice())
} else {
    fmt.Printf("⚠️ STOP_LIMIT sent | price=%.5f (no ticket in response)\n", res.GetPrice())
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) PlaceStopLimit(ctx context.Context, symbol string, isBuy bool, volume, trigger, limit float64, sl, tp *float64, exp *timestamppb.Timestamp)
```

---

## 🔽 Input

| Parameter | Type                         | Required | Description                                                    |
| --------- | ---------------------------- | -------- | -------------------------------------------------------------- |
| `ctx`     | `context.Context`            | yes      | Timeout/cancel control.                                        |
| `symbol`  | `string`                     | yes      | Symbol name.                                                   |
| `isBuy`   | `bool`                       | yes      | Direction: `true` → Buy Stop‑Limit, `false` → Sell Stop‑Limit. |
| `volume`  | `float64`                    | yes      | Lots (respect `VolumeMin/Max/Step`).                           |
| `trigger` | `float64`                    | yes      | **Stop price** that arms the order.                            |
| `limit`   | `float64`                    | yes      | **Limit price** to execute after trigger.                      |
| `sl`      | `*float64`                   | no       | Stop Loss **absolute** price.                                  |
| `tp`      | `*float64`                   | no       | Take Profit **absolute** price.                                |
| `exp`     | `*google.protobuf.Timestamp` | no       | Expiration time; if set, time-in-force becomes `*_SPECIFIED`.  |

---

## ⬆️ Output

Returns **`OrderSendData`**.

| Field                     | Type     | Description                                          |
| ------------------------- | -------- | ---------------------------------------------------- |
| `Order`                   | `uint64` | Pending order ticket (expected for Stop‑Limit).      |
| `Deal`                    | `uint64` | Deal ticket if order triggered & filled immediately. |
| `Price`                   | `double` | Server price (acknowledged).                         |
| `ReturnedCode`            | `uint32` | Numeric result code.                                 |
| `ReturnedStringCode`      | `string` | Short code.                                          |
| `ReturnedCodeDescription` | `string` | Human description.                                   |

---

## Enums used

* `MRPC_ENUM_TRADE_REQUEST_ACTIONS`: `TRADE_ACTION_PENDING` (check stage).
* `ENUM_ORDER_TYPE_TF`: `*_STOP_LIMIT` (check stage; depends on direction).
* `TMT5_ENUM_ORDER_TYPE`: `*_STOP_LIMIT` (send stage; depends on direction).
* `MRPC_ENUM_ORDER_TYPE_FILLING`: typically `ORDER_FILLING_FOK`.
* `MRPC_ENUM_ORDER_TYPE_TIME`: `ORDER_TIME_GTC` or `ORDER_TIME_SPECIFIED` when `exp` is provided.

---

## 🎯 Purpose

* Combine breakout logic with price control: trigger at one level, execute as a limit at another.
* Useful to avoid poor fills during fast moves.

---

## 🧩 Notes & Tips

* **Direction rules:** for **Buy Stop‑Limit** usually `limit ≥ trigger`; for **Sell Stop‑Limit** usually `limit ≤ trigger`. Brokers enforce constraints.
* If market gaps through your trigger and limit, fill may not happen — that’s by design of limit. Use plain Stop if you must be filled.
* SL/TP are absolute prices; compute from pips/ticks using symbol `Digits`.
* Always inspect `ReturnedCode` & server comment to troubleshoot rejects.
