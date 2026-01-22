# ✅ Place Market and Pending Orders

> **Request:** primary method for opening positions and placing pending orders. Supports all order types: market buy/sell, buy/sell limit, buy/sell stop, and stop-limit.

**API Information:**

* **Low-level API:** `MT5Account.OrderSend(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.TradingHelper`
* **Proto definition:** `OrderSend` (defined in `mt5-term-api-trading-helper.proto`)

### RPC

* **Service:** `mt5_term_api.TradingHelper`
* **Method:** `OrderSend(OrderSendRequest) → OrderSendReply`
* **Low‑level client (generated):** `TradingHelperClient.OrderSend(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Main trading method for placing market and pending orders.
* **Why you need it.** Open positions, place limit/stop orders, enter trades.
* **All order types.** Supports market orders, limits, stops, and stop-limits.

---

## 🎯 Purpose

Use it to:

* Open market positions (Buy/Sell)
* Place pending orders (Limit/Stop)
* Enter trades with SL/TP
* Execute trading strategies
* Manage order parameters

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OrderSend - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderSend_HOW.md)**

---


```go
package mt5

type MT5Account struct {
    // ...
}

// OrderSend places a market or pending order.
// Returns OrderSendData with ticket number, execution price, and MqlTradeResult.
func (a *MT5Account) OrderSend(
    ctx context.Context,
    req *pb.OrderSendRequest,
) (*pb.OrderSendData, error)
```

**Request message:**

```protobuf
OrderSendRequest {
  string symbol = 1;                                            // Symbol name (REQUIRED)
  TMT5_ENUM_ORDER_TYPE operation = 2;                           // Type of order (REQUIRED)
  double volume = 3;                                            // Requested volume for a deal in lots (REQUIRED)
  optional double price = 4;                                    // Price
  optional uint64 slippage = 5;                                 // Maximal possible deviation from the requested price
  optional double stop_loss = 6;                                // Stop Loss level of the order
  optional double take_profit = 7;                              // Take Profit level of the order
  optional string comment = 8;                                  // Order comment
  optional uint64 expert_id = 9;                                // Expert Advisor ID (magic number)
  optional double stop_limit_price = 10;                        // StopLimit level of the order
  optional TMT5_ENUM_ORDER_TYPE_TIME expiration_time_type = 11; // Order expiration type
  optional google.protobuf.Timestamp expiration_time = 12;      // Order expiration time (for ORDER_TIME_SPECIFIED type)
}
```

**Reply message:**

```protobuf
OrderSendReply {
  oneof response {
    OrderSendData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                      | Description                                   |
| --------- | ------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`         | Context for deadline/timeout and cancellation |
| `req`     | `*pb.OrderSendRequest`    | Request with order parameters                 |

**Request fields:**

| Field                 | Type                                | Required | Description                                          |
| --------------------- | ----------------------------------- | -------- | ---------------------------------------------------- |
| `Symbol`              | `string`                            | ✅       | Trading symbol (e.g., "EURUSD")                      |
| `Operation`           | `TMT5_ENUM_ORDER_TYPE`              | ✅       | Order type (BUY=0, SELL=1, BUY_LIMIT=2, etc)         |
| `Volume`              | `double`                            | ✅       | Order volume in lots                                 |
| `Price`               | `double`                            | optional | Order price (required for pending orders)            |
| `Slippage`            | `uint64`                            | optional | Maximum price deviation (slippage) in points         |
| `StopLoss`            | `double`                            | optional | Stop loss price level                                |
| `TakeProfit`          | `double`                            | optional | Take profit price level                              |
| `Comment`             | `string`                            | optional | Order comment                                        |
| `ExpertId`            | `uint64`                            | optional | Expert Advisor ID (magic number)                     |
| `StopLimitPrice`      | `double`                            | optional | StopLimit level (for STOP_LIMIT orders)              |
| `ExpirationTimeType`  | `TMT5_ENUM_ORDER_TYPE_TIME`         | optional | Order expiration type (GTC, DAY, SPECIFIED, etc)     |
| `ExpirationTime`      | `google.protobuf.Timestamp`         | optional | Expiration time (for SPECIFIED type)                 |

---

## ⬆️ Output — `OrderSendData`

| Field                      | Type     | Description                                           |
| -------------------------- | -------- | ----------------------------------------------------- |
| `ReturnedCode`             | `uint32` | Operation return code (10009 = success)               |
| `Deal`                     | `uint64` | Deal ticket, if it is performed                       |
| `Order`                    | `uint64` | Order ticket, if it is placed                         |
| `Volume`                   | `double` | Deal volume, confirmed by broker                      |
| `Price`                    | `double` | Deal price, confirmed by broker                       |
| `Bid`                      | `double` | Current Bid price                                     |
| `Ask`                      | `double` | Current Ask price                                     |
| `Comment`                  | `string` | Broker comment to operation                           |
| `RequestId`                | `uint32` | Request ID set by terminal during dispatch            |
| `RetCodeExternal`          | `int32`  | Return code of external trading system                |
| `ReturnedStringCode`       | `string` | String representation of return code                  |
| `ReturnedCodeDescription`  | `string` | Description of return code                            |

---

> **💡 Enum Usage Note:** The tables show simplified constant names for readability.
> In Go code, use full names with the enum type prefix.
>
> Format: `pb.<ENUM_TYPE>_<CONSTANT_NAME>`
>
> Example: `pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY`

### 📘 Enum: TMT5_ENUM_ORDER_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `TMT5_ORDER_TYPE_BUY` | Market Buy order |
| 1 | `TMT5_ORDER_TYPE_SELL` | Market Sell order |
| 2 | `TMT5_ORDER_TYPE_BUY_LIMIT` | Buy Limit pending order |
| 3 | `TMT5_ORDER_TYPE_SELL_LIMIT` | Sell Limit pending order |
| 4 | `TMT5_ORDER_TYPE_BUY_STOP` | Buy Stop pending order |
| 5 | `TMT5_ORDER_TYPE_SELL_STOP` | Sell Stop pending order |
| 6 | `TMT5_ORDER_TYPE_BUY_STOP_LIMIT` | Upon reaching the order price, a pending Buy Limit order is placed at the StopLimit price |
| 7 | `TMT5_ORDER_TYPE_SELL_STOP_LIMIT` | Upon reaching the order price, a pending Sell Limit order is placed at the StopLimit price |
| 8 | `TMT5_ORDER_TYPE_CLOSE_BY` | Order to close a position by an opposite one |

### 📘 Enum: TMT5_ENUM_ORDER_TYPE_TIME

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `TMT5_ORDER_TIME_GTC` | Good Till Cancelled (order valid until explicitly cancelled) |
| 1 | `TMT5_ORDER_TIME_DAY` | Good Till Day (order valid until end of trading day) |
| 2 | `TMT5_ORDER_TIME_SPECIFIED` | Good Till Specified Time (order valid until ExpirationTime) |
| 3 | `TMT5_ORDER_TIME_SPECIFIED_DAY` | Good Till Specified Day (order valid until end of specified day) |

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `30s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Use OrderCheck:** Validate orders before sending with OrderCheck method.
* **Slippage:** Deviation parameter controls maximum acceptable slippage.

---

## 🔗 Usage Examples

### 1) Place market buy order

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Get current price
    tick, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: "EURUSD"})

    sl := tick.Ask - 0.0050  // 50 pips SL
    tp := tick.Ask + 0.0100  // 100 pips TP
    slippage := uint64(20)
    expertId := uint64(123456)
    comment := "Test buy order"

    data, err := account.OrderSend(ctx, &pb.OrderSendRequest{
        Symbol:     "EURUSD",
        Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
        Volume:     0.01,
        Price:      &tick.Ask,
        StopLoss:   &sl,
        TakeProfit: &tp,
        Slippage:   &slippage,
        ExpertId:   &expertId,
        Comment:    &comment,
    })
    if err != nil {
        panic(err)
    }

    if data.ReturnedCode == 10009 {
        fmt.Printf("Order placed! Ticket: %d, Price: %.5f\n", data.Order, data.Price)
    } else {
        fmt.Printf("Order failed: %s\n", data.ReturnedCodeDescription)
    }
}
```

### 2) Place market sell order

```go
func PlaceSellOrder(account *mt5.MT5Account, symbol string, volume float64) (uint64, error) {
    ctx := context.Background()

    // Get current price
    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})
    if err != nil {
        return 0, err
    }

    sl := tick.Bid + 0.0050
    tp := tick.Bid - 0.0100
    slippage := uint64(20)
    comment := "Sell order"

    data, err := account.OrderSend(ctx, &pb.OrderSendRequest{
        Symbol:     symbol,
        Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL,
        Volume:     volume,
        Price:      &tick.Bid,
        StopLoss:   &sl,
        TakeProfit: &tp,
        Slippage:   &slippage,
        Comment:    &comment,
    })
    if err != nil {
        return 0, fmt.Errorf("failed to place sell order: %w", err)
    }

    return data.Order, nil
}
```

### 3) Place buy limit order

```go
func PlaceBuyLimit(account *mt5.MT5Account, symbol string, volume, price, sl, tp float64) (uint64, error) {
    ctx := context.Background()

    expertId := uint64(123456)
    comment := "Buy limit order"

    data, err := account.OrderSend(ctx, &pb.OrderSendRequest{
        Symbol:     symbol,
        Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT,
        Volume:     volume,
        Price:      &price,
        StopLoss:   &sl,
        TakeProfit: &tp,
        ExpertId:   &expertId,
        Comment:    &comment,
    })
    if err != nil {
        return 0, err
    }

    fmt.Printf("Buy limit placed at %.5f, ticket: %d\n", price, data.Order)
    return data.Order, nil
}
```

### 4) Place sell stop order

```go
func PlaceSellStop(account *mt5.MT5Account, symbol string, volume, stopPrice float64) (uint64, error) {
    ctx := context.Background()

    sl := stopPrice + 0.0050
    tp := stopPrice - 0.0100
    comment := "Sell stop order"

    data, err := account.OrderSend(ctx, &pb.OrderSendRequest{
        Symbol:     symbol,
        Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP,
        Volume:     volume,
        Price:      &stopPrice,
        StopLoss:   &sl,
        TakeProfit: &tp,
        Comment:    &comment,
    })
    if err != nil {
        return 0, err
    }

    return data.Order, nil
}
```

### 5) Order with validation

```go
func PlaceMarketOrder(account *mt5.MT5Account, symbol string, volume float64, isBuy bool) (uint64, error) {
    ctx := context.Background()

    // Get current price
    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})
    if err != nil {
        return 0, err
    }

    var operation pb.TMT5_ENUM_ORDER_TYPE
    var price, sl, tp float64

    if isBuy {
        operation = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY
        price = tick.Ask
        sl = tick.Ask - 0.0050
        tp = tick.Ask + 0.0100
    } else {
        operation = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL
        price = tick.Bid
        sl = tick.Bid + 0.0050
        tp = tick.Bid - 0.0100
    }

    slippage := uint64(20)

    // Send order
    data, err := account.OrderSend(ctx, &pb.OrderSendRequest{
        Symbol:     symbol,
        Operation:  operation,
        Volume:     volume,
        Price:      &price,
        StopLoss:   &sl,
        TakeProfit: &tp,
        Slippage:   &slippage,
    })
    if err != nil {
        return 0, err
    }

    if data.ReturnedCode == 10009 {
        fmt.Printf("Order placed: %d at %.5f\n", data.Order, data.Price)
        return data.Order, nil
    }

    return 0, fmt.Errorf("order failed: %s", data.ReturnedCodeDescription)
}
```

---

## 🔧 Common Patterns

### Market order with retry

```go
func PlaceMarketOrderWithRetry(account *mt5.MT5Account, symbol string, volume float64, isBuy bool, maxRetries int) (uint64, error) {
    ctx := context.Background()

    slippage := uint64(50)

    for i := 0; i < maxRetries; i++ {
        tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})
        if err != nil {
            continue
        }

        var operation pb.TMT5_ENUM_ORDER_TYPE
        var price float64

        if isBuy {
            operation = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY
            price = tick.Ask
        } else {
            operation = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL
            price = tick.Bid
        }

        data, err := account.OrderSend(ctx, &pb.OrderSendRequest{
            Symbol:    symbol,
            Operation: operation,
            Volume:    volume,
            Price:     &price,
            Slippage:  &slippage,
        })

        if err == nil && data.ReturnedCode == 10009 {
            return data.Order, nil
        }

        time.Sleep(500 * time.Millisecond)
    }

    return 0, fmt.Errorf("failed after %d retries", maxRetries)
}
```

---

## 📚 See Also

* [OrderCheck](./OrderCheck.md) - Validate order before sending
* [OrderModify](./OrderModify.md) - Modify existing orders
* [OrderClose](./OrderClose.md) - Close positions
* [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - Get current prices
