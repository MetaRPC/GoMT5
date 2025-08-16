# Sending a Market Order (Sell Example)

> **Request:** demonstrate how to place a Sell Market order and parse the result.

---

### Code Example

```go
// High-level helper:
svc.SellMarket(ctx, "EURUSD", 0.10, nil, nil)

// Internally:
volume := 0.10
slippage := int32(20)         // allowed slippage in points
comment := "Sell market order"
magic := int32(12345)

res, err := svc.account.OrderSend(
    ctx,
    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL, // direct market deal
    pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL,             // Sell
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
    log.Printf("❌ SellMarket error: %v", err)
    return
}

// Read result
deal := res.GetDeal()
price := res.GetPrice()
code  := res.GetReturnedCode()
if deal != 0 {
    fmt.Printf("✅ Sell executed! Deal:%d Volume:%.2f Price:%.5f Code:%d\n", deal, volume, price, code)
} else {
    fmt.Printf("⚠️ No deal ticket returned | Price: %.5f | Code: %d\n", price, code)
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) SellMarket(ctx context.Context, symbol string, volume float64, sl *float64, tp *float64)
```

---

## 🔽 Input

Helper accepts `ctx`, `symbol`, `volume`, optional `sl`, and `tp`. The underlying `OrderSend` call has full parameter set:

| Parameter    | Type                              | Required | Description                            |
| ------------ | --------------------------------- | -------- | -------------------------------------- |
| `ctx`        | `context.Context`                 | yes      | Timeout/cancel control.                |
| `action`     | `MRPC_ENUM_TRADE_REQUEST_ACTIONS` | yes      | `TRADE_ACTION_DEAL` for market trades. |
| `orderType`  | `ENUM_ORDER_TYPE_TF`              | yes      | `ORDER_TYPE_TF_SELL` for Sell.         |
| `symbol`     | `string`                          | yes      | Symbol name.                           |
| `volume`     | `float64`                         | yes      | Lots to sell.                          |
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

* Execute a direct **Sell market order**.
* Straightforward way to enter short exposure at current market price.
* Common for testing connectivity, permissions, and latency.

---

## 🧩 Notes & Tips

* Use `sl`/`tp` to attach protective orders; they must be **absolute prices**.
* Ensure your SL price is *above* market for sells, TP is *below* (basic sanity check to avoid broker rejects).
* Slippage is in points; align with your broker’s pricing (e.g., 1 pip = 10 points on 5-digit FX).
* Always validate `ReturnedCode`/logs for possible rejections (off quotes, trade disabled, etc.).
