# Sending a Simple Market/Pending Order (Example)

> **Request:** demonstrate a basic `OrderSend` call and how to read its result.

---

### Code Example

```go
// High-level helper (prints a human-readable result):
svc.ShowOrderSendExample(ctx, "EURUSD")

// What it does internally (simplified):
volume := 0.10
slippage := int32(5)                 // points
comment := "Go order test"
magic := int32(123456)

res, err := svc.account.OrderSend(
    ctx,
    "EURUSD",                                        // symbol
    pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,      // market BUY
    volume,                                           // lots
    nil,                                              // price: nil for market (server uses current)
    &slippage,                                        // max slippage (points)
    nil,                                              // SL (absolute price) or nil
    nil,                                              // TP (absolute price) or nil
    &comment,                                         // comment
    &magic,                                           // EA magic number
    nil,                                              // expiration (for pending with specified time)
)
if err != nil {
    log.Printf("❌ OrderSend error: %v", err)
    return
}

// Read result (order placed or deal executed):
order := res.GetOrder() // pending order ticket (if placed as pending)
deal  := res.GetDeal()  // deal ticket (if executed immediately)
price := res.GetPrice()
vol   := res.GetVolume()
code  := res.GetReturnedCode()
if deal != 0 {
    fmt.Printf("✅ Market executed! Deal:%d Price:%.5f Volume:%.2f Code:%d\n", deal, price, vol, code)
} else if order != 0 {
    fmt.Printf("✅ Pending placed! Order:%d Price:%.5f Volume:%.2f Code:%d\n", order, price, vol, code)
} else {
    fmt.Printf("⚠️ No ticket returned | Price: %.5f | Volume: %.2f | Code: %d\n", price, vol, code)
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowOrderSendExample(ctx context.Context, symbol string)
```

---

## 🔽 Input

High-level helper takes only `ctx` and `symbol`. The underlying RPC accepts the full set below:

### Underlying `OrderSend` parameters

| Parameter    | Type                         | Required | Description                                        |
| ------------ | ---------------------------- | -------- | -------------------------------------------------- |
| `ctx`        | `context.Context`            | yes      | Timeout/cancel control.                            |
| `symbol`     | `string`                     | yes      | Symbol name.                                       |
| `orderType`  | `pb.TMT5_ENUM_ORDER_TYPE`    | yes      | Concrete TMT5 order type (see enum).               |
| `volume`     | `float64`                    | yes      | Lots (respect `VolumeMin/Max/Step`).               |
| `price`      | `*float64`                   | no       | `nil` for market; set absolute price for pendings. |
| `slippage`   | `*int32`                     | no       | Max slippage (points).                             |
| `sl`         | `*float64`                   | no       | Stop Loss **absolute price** or `nil`.             |
| `tp`         | `*float64`                   | no       | Take Profit **absolute price** or `nil`.           |
| `comment`    | `*string`                    | no       | User/EA comment.                                   |
| `magic`      | `*int32`                     | no       | EA magic number.                                   |
| `expiration` | `*google.protobuf.Timestamp` | no       | Expiry time for time‑specified pendings.           |

---

## ⬆️ Output

RPC returns **`OrderSendData`**. Typical fields used:

| Field                     | Type     | Description                                     |
| ------------------------- | -------- | ----------------------------------------------- |
| `Order`                   | `uint64` | Pending order ticket (if a pending was placed). |
| `Deal`                    | `uint64` | Deal ticket (if executed immediately).          |
| `Price`                   | `double` | Execution/placement price returned by server.   |
| `Volume`                  | `double` | Executed/placed volume.                         |
| `ReturnedCode`            | `uint32` | Numeric result code.                            |
| `ReturnedStringCode`      | `string` | Short string representation of the code.        |
| `ReturnedCodeDescription` | `string` | Human‑readable description of the code.         |

---

## Enum: `TMT5_ENUM_ORDER_TYPE`

| Code | Name                              | Meaning                 |
| ---: | --------------------------------- | ----------------------- |
|    0 | `TMT5_ORDER_TYPE_BUY`             | Market Buy              |
|    1 | `TMT5_ORDER_TYPE_SELL`            | Market Sell             |
|    2 | `TMT5_ORDER_TYPE_BUY_LIMIT`       | Pending Buy Limit       |
|    3 | `TMT5_ORDER_TYPE_SELL_LIMIT`      | Pending Sell Limit      |
|    4 | `TMT5_ORDER_TYPE_BUY_STOP`        | Pending Buy Stop        |
|    5 | `TMT5_ORDER_TYPE_SELL_STOP`       | Pending Sell Stop       |
|    6 | `TMT5_ORDER_TYPE_BUY_STOP_LIMIT`  | Pending Buy Stop‑Limit  |
|    7 | `TMT5_ORDER_TYPE_SELL_STOP_LIMIT` | Pending Sell Stop‑Limit |

> These are the concrete order types expected by `OrderSend`. For time‑in‑force and filling policy, see `OrderCheck` card.

---

## 🎯 Purpose

* Minimal example of placing a market/pending order from Go.
* Demonstrates reading both **Deal** and **Order** workflow outcomes.
* A solid template to copy and adjust (symbol, volume, SL/TP, etc.).

---

## 🧩 Notes & Tips

* For **market** orders pass `price=nil`; for **pending** pass your absolute entry price.
* SL/TP are **absolute prices**, not offsets — derive from pips/ticks using `Digits`.
* `slippage` is in points; make sure you map your UI “pips” correctly (5‑digit FX: 1 pip = 10 points).
* Always check `ReturnedCode`/`ReturnedStringCode` for diagnostics even if you got a ticket.
