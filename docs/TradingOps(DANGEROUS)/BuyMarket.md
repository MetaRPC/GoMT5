# Sending a Market Order (Buy Example)

> **Request:** demonstrate how to place a Buy Market order and parse the result.

---

### Code Example

```go
// High-level helper:
svc.BuyMarket(ctx, "EURUSD", 0.10, nil, nil)

// Internally:
volume := 0.10
slippage := int32(20)         // allowed slippage in points
comment := "Buy market order"
magic := int32(12345)

res, err := svc.account.OrderSend(
    ctx,
    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL, // direct market deal
    pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,              // Buy
    "EURUSD",
    volume,
    0,              // price=0 means "market" on most servers
    &slippage,
    nil,            // StopLoss (absolute price)
    nil,            // TakeProfit (absolute price)
    &comment,
    &magic,
    nil,            // expiration not used for market orders
)
if err != nil {
    log.Printf("❌ BuyMarket error: %v", err)
    return
}

// Read result
deal := res.GetDeal()
price := res.GetPrice()
code  := res.GetReturnedCode()
if deal != 0 {
    fmt.Printf("✅ Buy executed! Deal:%d Volume:%.2f Price:%.5f Code:%d\n", deal, volume, price, code)
} else {
    fmt.Printf("⚠️ No deal ticket returned | Price: %.5f | Code: %d\n", price, code)
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) BuyMarket(ctx context.Context, symbol string, volume float64, sl *float64, tp *float64)
```

---

## 🔽 Input

Helper accepts `ctx`, `symbol`, `volume`, optional `sl`, and `tp`. The underlying `OrderSend` call has full parameter set:

| Parameter    | Type                              | Required | Description                            |
| ------------ | --------------------------------- | -------- | -------------------------------------- |
| `ctx`        | `context.Context`                 | yes      | Timeout/cancel control.                |
| `action`     | `MRPC_ENUM_TRADE_REQUEST_ACTIONS` | yes      | `TRADE_ACTION_DEAL` for market trades. |
| `orderType`  | `ENUM_ORDER_TYPE_TF`              | yes      | `ORDER_TYPE_TF_BUY` for Buy.           |
| `symbol`     | `string`                          | yes      | Symbol name.                           |
| `volume`     | `float64`                         | yes      | Lots to buy.                           |
| `price`      | `float64`                         | no       | `0` → market execution.                |
| `slippage`   | `*int32`                          | no       | Max slippage (points).                 |
| `sl`         | `*float64`                        | no       | Stop Loss absolute price.              |
| `tp`         | `*float64`                        | no       | Take Profit absolute price.            |
| `comment`    | `*string`                         | no       | Order comment.                         |
| `magic`      | `*int32`                          | no       | EA magic number.                       |
| `expiration` | `*google.protobuf.Timestamp`      | no       | Not used for market orders.            |

---

## ⬆️ Output

Returns **`OrderSendData`** with main fields:

| Field                     | Type     | Description                                 |
| ------------------------- | -------- | ------------------------------------------- |
| `Deal`                    | `uint64` | Deal ticket if trade executed.              |
| `Order`                   | `uint64` | Order ticket (rarely used in market deals). |
| `Price`                   | `double` | Execution price.                            |
| `Volume`                  | `double` | Executed volume.                            |
| `ReturnedCode`            | `uint32` | Numeric result code.                        |
| `ReturnedStringCode`      | `string` | Short code.                                 |
| `ReturnedCodeDescription` | `string` | Human description.                          |

---

## 🎯 Purpose

* Execute a direct **Buy market order**.
* Simplest way to enter a trade at current market price.
* Common for testing connectivity and trading logic.

---

## 🧩 Notes & Tips

* Use `sl` and `tp` to attach StopLoss/TakeProfit in the same request.
* Slippage should reflect broker contract size: e.g., `20` points ≈ 2 pips on 5‑digit brokers.
* Execution is immediate or rejected depending on broker’s market model (Market/ECN/STP).
* Always check `ReturnedCode` and logs for broker rejections or errors.
