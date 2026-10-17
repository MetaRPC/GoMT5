# ✅ Modify Orders and Positions

> **Request:** modify parameters of existing pending orders or StopLoss/TakeProfit levels of open positions.

**API Information:**

* **Low-level API:** `MT5Account.OrderModify(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.TradingHelper`
* **Proto definition:** `OrderModify` (defined in `mt5-term-api-trading-helper.proto`)

### RPC

* **Service:** `mt5_term_api.TradingHelper`
* **Method:** `OrderModify(OrderModifyRequest) → OrderModifyReply`
* **Low‑level client (generated):** `TradingHelperClient.OrderModify(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Modifies parameters of existing orders and positions.
* **Why you need it.** Adjust SL/TP levels, change pending order prices, manage risk.
* **Flexible.** Works with both pending orders and open positions.

---

## 🎯 Purpose

Use it to:

* Adjust stop-loss and take-profit levels
* Move pending order entry prices
* Implement trailing stops
* Update risk management parameters
* Change order expiration times

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OrderModify - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderModify_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// OrderModify modifies an existing pending order or position.
// Can change entry price (pending orders), StopLoss, TakeProfit, and expiration.
func (a *MT5Account) OrderModify(
    ctx context.Context,
    req *pb.OrderModifyRequest,
) (*pb.OrderModifyData, error)
```

**Request message:**

```protobuf
OrderModifyRequest {
  uint64 ticket = 1;                                        // Pending order or position ticket (REQUIRED)
  optional double stop_loss  = 2;                           // New SL for position or pending order
  optional double take_profit  = 3;                         // New TP for position or pending order
  optional double price = 4;                                // New price ONLY for pending order modification
  optional TMT5_ENUM_ORDER_TYPE_TIME expiration_time_type = 5;  // Expiration type ONLY for pending order
  optional google.protobuf.Timestamp expiration_time = 6;   // Order expiration time (for ORDER_TIME_SPECIFIED), broker server time
  optional google.protobuf.Timestamp expiration_time_utc = 9; // expiration_time in UTC: send either expiration_time or expiration_time_utc, not both
  optional double stop_limit = 8;                           // New stop limit ONLY for pending order modification
}
```

---

## 🔽 Input

| Parameter | Type                       | Description                                   |
| --------- | -------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`          | Context for deadline/timeout and cancellation |
| `req`     | `*pb.OrderModifyRequest`   | Request with ticket and new parameters        |

**Request fields:**

| Field               | Type                                | Required | Description                                       |
| ------------------- | ----------------------------------- | -------- | ------------------------------------------------- |
| `Ticket`            | `uint64`                            | ✅       | Order ticket number to modify                     |
| `StopLoss`          | `*float64`                          | optional | New stop loss level (0 = remove)                  |
| `TakeProfit`        | `*float64`                          | optional | New take profit level (0 = remove)                |
| `Price`             | `*float64`                          | optional | New entry price (for pending orders only)         |
| `ExpirationTimeType`| `*TMT5_ENUM_ORDER_TYPE_TIME`        | optional | Expiration type (GTC, DAY, SPECIFIED, etc)        |
| `ExpirationTime`    | `*timestamppb.Timestamp`            | optional | Expiration time (for SPECIFIED type), broker server time |
| `ExpirationTimeUtc` | `*timestamppb.Timestamp`            | optional | Expiration time in UTC (send this or `ExpirationTime`, not both) |
| `StopLimit`         | `*float64`                          | optional | New stop limit level (for stop-limit orders)      |

---

## ⬆️ Output — `OrderModifyData`

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
> Example: `pb.TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_GTC`

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
* **Remove SL/TP:** Set StopLoss or TakeProfit to 0 to remove them.
* **Price for positions:** Cannot modify entry price of open positions (only pending orders).

---

## 🔗 Usage Examples

### 1) Modify stop-loss and take-profit

```go
package main

import (
    "github.com/google/uuid"
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
)

func main() {
    account, err := mt5.NewMT5Account(591129415, "password", "mt5.mrpc.pro:443", uuid.Nil)
    if err != nil {
        panic(err)
    }
    defer account.Close()

    // Connect first
    _, err = account.ConnectEx(context.Background(), &pb.ConnectExRequest{User: 591129415, Password: "password", MtClusterName: "MetaQuotes-Demo"})
    if err != nil {
        panic(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    ticket := uint64(12345678)
    newSL := 1.10000
    newTP := 1.12000

    data, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
        Ticket:     ticket,
        StopLoss:   &newSL,
        TakeProfit: &newTP,
    })
    if err != nil {
        panic(err)
    }

    // Check modification result using helper from errors.go
    if mt5.IsRetCodeSuccess(data.ReturnedCode) {
        fmt.Printf("Order %d modified: SL=%.5f, TP=%.5f\n", ticket, newSL, newTP)
    } else {
        fmt.Printf("Modification failed: %s (code: %d)\n",
            mt5.GetRetCodeMessage(data.ReturnedCode),
            data.ReturnedCode)
    }
}
```

### 2) Move stop-loss to breakeven

```go
func MoveToBreakeven(account *mt5.MT5Account, ticket uint64, openPrice float64) error {
    ctx := context.Background()

    data, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
        Ticket:   ticket,
        StopLoss: &openPrice,
        // Don't pass TakeProfit to keep existing TP
    })
    if err != nil {
        return fmt.Errorf("failed to move to breakeven: %w", err)
    }

    // Validate modification result using helper from errors.go
    if !mt5.IsRetCodeSuccess(data.ReturnedCode) {
        return fmt.Errorf("modification unsuccessful: %s (code: %d)",
            mt5.GetRetCodeMessage(data.ReturnedCode),
            data.ReturnedCode)
    }

    fmt.Printf("Moved order %d to breakeven\n", ticket)
    return nil
}
```

### 3) Implement trailing stop

```go
func TrailingStop(account *mt5.MT5Account, ticket uint64, trailDistance float64) error {
    ctx := context.Background()

    // Get current position
    positions, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{})
    if err != nil {
        return err
    }

    var pos *pb.PositionInfo
    for _, p := range positions.PositionInfos {
        if p.Ticket == ticket {
            pos = p
            break
        }
    }

    if pos == nil {
        return fmt.Errorf("position not found")
    }

    // Calculate new SL
    var newSL float64
    if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
        newSL = pos.PriceCurrent - trailDistance
        // Only move SL up, never down
        if newSL <= pos.StopLoss {
            return nil // No change needed
        }
    } else { // Sell position
        newSL = pos.PriceCurrent + trailDistance
        // Only move SL down, never up
        if pos.StopLoss > 0 && newSL >= pos.StopLoss {
            return nil // No change needed
        }
    }

    data, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
        Ticket:   ticket,
        StopLoss: &newSL,
    })
    if err != nil {
        return err
    }

    // Check modification result using helper from errors.go
    if mt5.IsRetCodeSuccess(data.ReturnedCode) {
        fmt.Printf("Trailing stop updated for %d: new SL=%.5f\n", ticket, newSL)
    }

    return nil
}
```

### 4) Modify pending order price

```go
func ModifyPendingOrderPrice(account *mt5.MT5Account, ticket uint64, newPrice float64) error {
    ctx := context.Background()

    data, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
        Ticket: ticket,
        Price:  &newPrice,
    })
    if err != nil {
        return err
    }

    if data.ReturnedCode != 10009 {
        return fmt.Errorf("failed to modify pending order price: %s", data.ReturnedCodeDescription)
    }

    fmt.Printf("Pending order %d price changed to %.5f\n", ticket, newPrice)
    return nil
}
```

### 5) Remove stop-loss

```go
func RemoveStopLoss(account *mt5.MT5Account, ticket uint64) error {
    ctx := context.Background()

    zeroSL := 0.0
    data, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
        Ticket:   ticket,
        StopLoss: &zeroSL, // 0 means remove SL
    })
    if err != nil {
        return err
    }

    if data.ReturnedCode == 10009 {
        fmt.Printf("Stop-loss removed from order %d\n", ticket)
    }

    return nil
}
```

---

## 🔧 Common Patterns

### Safe SL/TP modification

```go
func SafeModifySLTP(account *mt5.MT5Account, ticket uint64, newSL, newTP float64) error {
    ctx := context.Background()

    // Validate stop levels first
    stopsLevel, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol: "EURUSD", // Should get symbol from position
        Type:   pb.SymbolInfoIntegerProperty_SYMBOL_TRADE_STOPS_LEVEL,
    })
    if err != nil {
        return err
    }

    // Add validation logic here (stopsLevel.Value = minimum SL/TP distance in points)...
    _ = stopsLevel.Value

    data, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
        Ticket:     ticket,
        StopLoss:   &newSL,
        TakeProfit: &newTP,
    })

    if err != nil {
        return fmt.Errorf("modification failed: %w", err)
    }

    if data.ReturnedCode != 10009 {
        return fmt.Errorf("modification unsuccessful: %s", data.ReturnedCodeDescription)
    }

    return nil
}
```

---

## 📚 See Also

* [OrderSend](./OrderSend.md) - Place new orders
* [OrderClose](./OrderClose.md) - Close positions
* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get current positions
* [SymbolInfoInteger](../2.%20Symbol_information/SymbolInfoInteger.md) - Get stops level
