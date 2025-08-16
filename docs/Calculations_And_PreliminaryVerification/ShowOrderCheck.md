# Validating a Trade Request (OrderCheck)

> **Request:** dry‑run validation of a trade request without placing it.

---

### Code Example

```go
// High-level (prints retcode, comment, and key balances):
svc.ShowOrderCheck(ctx,
    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL, // action
    pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,              // orderType
    "EURUSD",                                            // symbol
    0.10,                                                 // volume (lots)
    0,                                                    // price (0 for market)
    nil, nil,                                             // sl, tp (optional)
    nil,                                                  // deviation (optional)
    nil,                                                  // magic (optional)
    nil,                                                  // expiration (optional)
)

// Low-level (build full request):
req := &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Symbol:    "EURUSD",
        Volume:    0.10,
        Price:     0, // 0 for market; set entry price for pendings
        StopLoss:  0, // 0 = not set; otherwise absolute price
        TakeProfit:0, // 0 = not set; otherwise absolute price
        Deviation: 10, // points (broker settings apply)
        TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
        TypeTime:    pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC,
        Expiration:  nil, // set only with *_SPECIFIED time
        ExpertAdvisorMagicNumber: 123456,
        Comment: "check via Go",
    },
}
res, err := svc.account.OrderCheck(ctx, req)
if err != nil {
    log.Printf("❌ OrderCheck error: %v", err)
    return
}
chk := res.GetMqlTradeCheckResult()
fmt.Printf("retcode=%d comment=%q margin=%.2f free=%.2f profit=%.2f\n",
    chk.GetReturnedCode(), chk.GetComment(), chk.GetMargin(), chk.GetFreeMargin(), chk.GetProfit())
```

---

### Method Signature

```go
func (s *MT5Service) ShowOrderCheck(
    ctx context.Context,
    action pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS,
    orderType pb.ENUM_ORDER_TYPE_TF,
    symbol string,
    volume float64,
    price float64,
    sl, tp *float64,
    deviation *uint64,
    magic *uint64,
    expiration *timestamppb.Timestamp,
)
```

---

## 🔽 Input

### High-level helper parameters

| Parameter    | Type                                 | Required | Description                                                                 |
| ------------ | ------------------------------------ | -------- | --------------------------------------------------------------------------- |
| `ctx`        | `context.Context`                    | yes      | Timeout/cancel control.                                                     |
| `action`     | `pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS` | yes      | What to do (deal/pending/modify/close/etc).                                 |
| `orderType`  | `pb.ENUM_ORDER_TYPE_TF`              | yes      | Order type (BUY/SELL/limits/stops/etc).                                     |
| `symbol`     | `string`                             | yes      | Symbol name.                                                                |
| `volume`     | `float64` (lots)                     | yes      | Trade volume; respect broker `VolumeMin/Max/Step`.                          |
| `price`      | `float64`                            | yes      | Price for the scenario. Use **0 for market**; set entry price for pendings. |
| `sl`         | `*float64`                           | no       | Stop Loss **absolute** price; `nil` → not set.                              |
| `tp`         | `*float64`                           | no       | Take Profit **absolute** price; `nil` → not set.                            |
| `deviation`  | `*uint64`                            | no       | Max slippage (points) for market orders; `nil` → broker default.            |
| `magic`      | `*uint64`                            | no       | Expert Advisor magic number.                                                |
| `expiration` | `*google.protobuf.Timestamp`         | no       | Expiration time; used with `*_SPECIFIED` time‑in‑force.                     |

### Underlying struct: `MrpcMqlTradeRequest`

| Field                      | Type                                   | Notes                                      |
| -------------------------- | -------------------------------------- | ------------------------------------------ |
| `Action`                   | `pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS`   | See enum below.                            |
| `OrderType`                | `pb.ENUM_ORDER_TYPE_TF`                | See enum below.                            |
| `Symbol`                   | `string`                               |                                            |
| `Volume`                   | `double`                               | Lots.                                      |
| `Price`                    | `double`                               | 0 for market.                              |
| `StopLoss`                 | `double`                               | 0 = not set.                               |
| `TakeProfit`               | `double`                               | 0 = not set.                               |
| `Deviation`                | `uint64`                               | Points.                                    |
| `TypeFilling`              | `pb.MRPC_ENUM_ORDER_TYPE_FILLING`      | FOK/IOC.                                   |
| `TypeTime`                 | `pb.MRPC_ENUM_ORDER_TYPE_TIME`         | GTC/Day/Specified.                         |
| `Expiration`               | `google.protobuf.Timestamp` (optional) | Required when `TypeTime` is `*_SPECIFIED`. |
| `ExpertAdvisorMagicNumber` | `uint64`                               |                                            |
| `Comment`                  | `string`                               |                                            |

---

## ⬆️ Output

Returns **`OrderCheckData`** with `MqlTradeCheckResult` inside.

### `MqlTradeCheckResult` (key fields)

| Field              | Type     | Description                                                           |
| ------------------ | -------- | --------------------------------------------------------------------- |
| `ReturnedCode`     | `uint32` | Result/retcode of the check (0 = success; codes are server‑specific). |
| `Comment`          | `string` | Human‑readable comment from server.                                   |
| `BalanceAfterDeal` | `double` | Balance after hypothetical execution.                                 |
| `EquityAfterDeal`  | `double` | Equity after hypothetical execution.                                  |
| `Profit`           | `double` | P/L effect of the scenario.                                           |
| `Margin`           | `double` | Margin that would be used.                                            |
| `FreeMargin`       | `double` | Free margin left.                                                     |
| `MarginLevel`      | `double` | Margin level percentage.                                              |

---

## Enums

### `MRPC_ENUM_TRADE_REQUEST_ACTIONS`

| Code | Name                   | Meaning                          |
| ---: | ---------------------- | -------------------------------- |
|    0 | `TRADE_ACTION_DEAL`    | Market execution (buy/sell).     |
|    1 | `TRADE_ACTION_PENDING` | Place pending order.             |
|    2 | `TRADE_ACTION_SLTP`    | Modify SL/TP for position/order. |
|    3 | `TRADE_ACTION_MODIFY`  | Modify order parameters.         |
|    4 | `TRADE_ACTION_REMOVE`  | Delete pending order.            |

> There might be additional values depending on your `pb` version; always check the source.

### `ENUM_ORDER_TYPE_TF`

(see full table in *OrderCalcMargin* / *OrderCalcProfit* cards).

### `MRPC_ENUM_ORDER_TYPE_FILLING`

| Code | Name                | Meaning              |
| ---: | ------------------- | -------------------- |
|    0 | `ORDER_FILLING_FOK` | Fill‑or‑Kill.        |
|    1 | `ORDER_FILLING_IOC` | Immediate‑or‑Cancel. |

### `MRPC_ENUM_ORDER_TYPE_TIME`

| Code | Name                       | Meaning                         |
| ---: | -------------------------- | ------------------------------- |
|    0 | `ORDER_TIME_GTC`           | Good‑Till‑Cancel.               |
|    1 | `ORDER_TIME_DAY`           | Good‑For‑Day.                   |
|    2 | `ORDER_TIME_SPECIFIED`     | Good‑Till specified timestamp.  |
|    3 | `ORDER_TIME_SPECIFIED_DAY` | Good‑Till end of specified day. |

---

## 🎯 Purpose

* Pre‑flight validation before `OrderSend` to avoid rejects and surprises.
* Estimate balances/margin after the hypothetical deal.
* Validate broker rules (filling policy, time‑in‑force) per symbol.

---

## 🧩 Notes & Tips

* **Price=0** is commonly accepted for market checks; for pendings always provide the intended entry price.
* `StopLoss`/`TakeProfit` are absolute prices (not offsets); compute them from ticks/pips and `Digits`.
* `Deviation` is in **points**; map your UI “pips” carefully (for 5‑digit FX, 1 pip = 10 points).
* If you pass `TypeTime=*_SPECIFIED`, you must set `Expiration`, otherwise the server will reject the request.
* Treat `ReturnedCode` as authoritative; even successful checks don’t guarantee later `OrderSend` success if market moves.
