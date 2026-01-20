# ✅ Validate Order Before Execution

> **Request:** validate order validity before sending it to the server without actual placement. Returns margin requirements, potential profit, and possible errors.

**API Information:**

* **SDK wrapper:** `MT5Account.OrderCheck(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.TradeFunctions`
* **Proto definition:** `OrderCheck` (defined in `mt5-term-api-trade-functions.proto`)

### RPC

* **Service:** `mt5_term_api.TradeFunctions`
* **Method:** `OrderCheck(OrderCheckRequest) → OrderCheckReply`
* **Low‑level client (generated):** `TradeFunctionsClient.OrderCheck(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

## 💬 Just the essentials

* **What it is.** Pre-validates trading requests without placing actual orders.
* **Why you need it.** Detect errors, check margin, estimate profit before trading.
* **Safe testing.** Test order validity without risking real execution.

---

## 🎯 Purpose

Use it to:

* Validate orders before OrderSend
* Check margin requirements
* Estimate potential profit/loss
* Detect invalid parameters
* Implement pre-trade safety checks

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OrderCheck - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderCheck_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// OrderCheck validates an order before sending it to the server.
// Returns margin requirements, validation status, and possible errors.
func (a *MT5Account) OrderCheck(
    ctx context.Context,
    req *pb.OrderCheckRequest,
) (*pb.OrderCheckData, error)
```

**Request message:**

```protobuf
OrderCheckRequest {
  MrpcMqlTradeRequest mql_trade_request = 1; // Nested trade request structure (REQUIRED)
}

MrpcMqlTradeRequest {
  MRPC_ENUM_TRADE_REQUEST_ACTIONS action = 1;           // Trade operation type
  uint64 expert_advisor_magic_number = 2;               // Expert Advisor ID (magic number)
  uint64 order = 3;                                     // Order ticket (for modify/remove actions)
  string symbol = 4;                                    // Symbol name
  double volume = 5;                                    // Requested volume for a deal in lots
  double price = 6;                                     // Price
  double stop_limit = 7;                                // StopLimit level of the order
  double stop_loss = 8;                                 // Stop Loss level of the order
  double take_profit = 9;                               // Take Profit level of the order
  uint64 deviation = 10;                                // Maximal possible deviation from the requested price
  ENUM_ORDER_TYPE_TF order_type = 11;                   // Type of order
  MRPC_ENUM_ORDER_TYPE_FILLING type_filling = 12;       // Order execution type (FOK, IOC, RETURN)
  MRPC_ENUM_ORDER_TYPE_TIME type_time = 13;             // Order expiration type
  google.protobuf.Timestamp expiration = 14;            // Order expiration time
  string comment = 15;                                  // Order comment
  uint64 position = 16;                                 // Position ticket
  uint64 position_by = 17;                              // Ticket of an opposite position
}
```

---

## 🔽 Input

| Parameter | Type                      | Description                                   |
| --------- | ------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`         | Context for deadline/timeout and cancellation |
| `req`     | `*pb.OrderCheckRequest`   | Request with nested MqlTradeRequest structure |

**Request contains a nested MrpcMqlTradeRequest with these fields:**

| Field                         | Type                                 | Required | Description                                           |
| ----------------------------- | ------------------------------------ | -------- | ----------------------------------------------------- |
| `Action`                      | `MRPC_ENUM_TRADE_REQUEST_ACTIONS`    | ✅       | Trade action: DEAL (0), PENDING (1), SLTP (2), etc    |
| `ExpertAdvisorMagicNumber`    | `uint64`                             | optional | Expert Advisor ID (magic number)                      |
| `Order`                       | `uint64`                             | optional | Order ticket (for modify/remove actions)              |
| `Symbol`                      | `string`                             | ✅       | Trading symbol                                        |
| `Volume`                      | `double`                             | ✅       | Order volume in lots                                  |
| `Price`                       | `double`                             | ✅       | Order price                                           |
| `StopLimit`                   | `double`                             | optional | StopLimit level                                       |
| `StopLoss`                    | `double`                             | optional | Stop Loss level                                       |
| `TakeProfit`                  | `double`                             | optional | Take Profit level                                     |
| `Deviation`                   | `uint64`                             | optional | Maximum price deviation (slippage)                    |
| `OrderType`                   | `ENUM_ORDER_TYPE_TF`                 | ✅       | Order type: BUY (0), SELL (1), BUY_LIMIT (2), etc     |
| `TypeFilling`                 | `MRPC_ENUM_ORDER_TYPE_FILLING`       | optional | Filling type: FOK (0), IOC (1), RETURN (2), BOC (3)   |
| `TypeTime`                    | `MRPC_ENUM_ORDER_TYPE_TIME`          | optional | Time type: GTC (0), DAY (1), SPECIFIED (2)            |
| `Expiration`                  | `google.protobuf.Timestamp`          | optional | Expiration time                                       |
| `Comment`                     | `string`                             | optional | Order comment                                         |
| `Position`                    | `uint64`                             | optional | Position ticket                                       |
| `PositionBy`                  | `uint64`                             | optional | Opposite position ticket                              |

---

## ⬆️ Output — `OrderCheckData`

| Field                  | Type                        | Description                              |
| ---------------------- | --------------------------- | ---------------------------------------- |
| `MqlTradeCheckResult`  | `MrpcMqlTradeCheckResult`   | Nested validation result structure       |

**MrpcMqlTradeCheckResult contains:**

| Field              | Type     | Description                                           |
| ------------------ | -------- | ----------------------------------------------------- |
| `ReturnedCode`     | `uint32` | Return code (10009 = valid, other = error)            |
| `BalanceAfterDeal` | `double` | Balance after the execution of the deal               |
| `EquityAfterDeal`  | `double` | Equity after the execution of the deal                |
| `Profit`           | `double` | Floating profit                                       |
| `Margin`           | `double` | Margin requirements                                   |
| `FreeMargin`       | `double` | Free margin after order                               |
| `MarginLevel`      | `double` | Margin level percentage                               |
| `Comment`          | `string` | Comment to the reply code (error description)         |

---

### 📘 Enum: MRPC_ENUM_TRADE_REQUEST_ACTIONS

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `TRADE_ACTION_DEAL` | Place a market order for immediate execution (BUY/SELL) |
| 1 | `TRADE_ACTION_PENDING` | Place a pending order (BUY_LIMIT, SELL_STOP, etc) |
| 2 | `TRADE_ACTION_SLTP` | Modify StopLoss/TakeProfit of an open position |
| 3 | `TRADE_ACTION_MODIFY` | Modify parameters of a pending order |
| 4 | `TRADE_ACTION_REMOVE` | Delete a pending order |
| 5 | `TRADE_ACTION_CLOSE_BY` | Close a position by an opposite one |

### 📘 Enum: ENUM_ORDER_TYPE_TF

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `ORDER_TYPE_TF_BUY` | Market Buy order |
| 1 | `ORDER_TYPE_TF_SELL` | Market Sell order |
| 2 | `ORDER_TYPE_TF_BUY_LIMIT` | Buy Limit pending order |
| 3 | `ORDER_TYPE_TF_SELL_LIMIT` | Sell Limit pending order |
| 4 | `ORDER_TYPE_TF_BUY_STOP` | Buy Stop pending order |
| 5 | `ORDER_TYPE_TF_SELL_STOP` | Sell Stop pending order |
| 6 | `ORDER_TYPE_TF_BUY_STOP_LIMIT` | Upon reaching the order price, a pending Buy Limit order is placed at the StopLimit price |
| 7 | `ORDER_TYPE_TF_SELL_STOP_LIMIT` | Upon reaching the order price, a pending Sell Limit order is placed at the StopLimit price |
| 8 | `ORDER_TYPE_TF_CLOSE_BY` | Order to close a position by an opposite one |

### 📘 Enum: MRPC_ENUM_ORDER_TYPE_FILLING

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `ORDER_FILLING_FOK` | Fill Or Kill (entire order must be filled immediately or cancelled) |
| 1 | `ORDER_FILLING_IOC` | Immediate Or Cancel (fill available volume immediately, cancel remainder) |
| 2 | `ORDER_FILLING_RETURN` | Return execution (partial fills allowed, remainder stays as pending order) |
| 3 | `ORDER_FILLING_BOC` | Book Or Cancel (order must be placed in the order book or cancelled) |

### 📘 Enum: MRPC_ENUM_ORDER_TYPE_TIME

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `ORDER_TIME_GTC` | Good Till Cancelled (order valid until explicitly cancelled) |
| 1 | `ORDER_TIME_DAY` | Good Till Day (order valid until end of trading day) |
| 2 | `ORDER_TIME_SPECIFIED` | Good Till Specified Time (order valid until Expiration time) |
| 3 | `ORDER_TIME_SPECIFIED_DAY` | Good Till Specified Day (order valid until end of specified day) |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `10s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **No actual trade:** OrderCheck does NOT place real orders.
* **Return codes:** ReturnedCode 10009 = valid order, other values indicate errors.
* **Nested structure:** Request requires wrapping order parameters in MqlTradeRequest field.

---

## 🔗 Usage Examples

### 1) Basic order validation

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/mt5"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    tick, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: "EURUSD"})

    sl := tick.Ask - 0.0050
    tp := tick.Ask + 0.0100

    data, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
        MqlTradeRequest: &pb.MrpcMqlTradeRequest{
            Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
            Symbol:     "EURUSD",
            Volume:     0.01,
            Price:      tick.Ask,
            OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
            StopLoss:   sl,
            TakeProfit: tp,
        },
    })
    if err != nil {
        panic(err)
    }

    result := data.MqlTradeCheckResult
    if result.ReturnedCode == 10009 {
        fmt.Println("Order is valid!")
        fmt.Printf("Required margin: %.2f\n", result.Margin)
        fmt.Printf("Free margin after: %.2f\n", result.FreeMargin)
    } else {
        fmt.Printf("Order invalid: %s\n", result.Comment)
    }
}
```

### 2) Check margin before trading

```go
func CheckMarginBeforeTrade(account *mt5.MT5Account, symbol string, volume float64) (bool, error) {
    ctx := context.Background()

    tick, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})

    data, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
        MqlTradeRequest: &pb.MrpcMqlTradeRequest{
            Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
            Symbol:    symbol,
            Volume:    volume,
            Price:     tick.Ask,
            OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        },
    })
    if err != nil {
        return false, err
    }

    result := data.MqlTradeCheckResult
    if result.ReturnedCode != 10009 {
        return false, fmt.Errorf("order invalid: %s", result.Comment)
    }

    // Check if we have enough margin
    if result.FreeMargin < 100 { // Example: need at least 100 currency units
        return false, fmt.Errorf("insufficient margin: %.2f", result.FreeMargin)
    }

    fmt.Printf("Margin check passed: Required=%.2f, Free after=%.2f\n",
        result.Margin, result.FreeMargin)
    return true, nil
}
```

### 3) Validate multiple lot sizes

```go
func FindMaxLotSize(account *mt5.MT5Account, symbol string) (float64, error) {
    ctx := context.Background()

    tick, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})

    testVolumes := []float64{1.0, 0.5, 0.1, 0.05, 0.01}

    for _, volume := range testVolumes {
        data, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
            MqlTradeRequest: &pb.MrpcMqlTradeRequest{
                Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
                Symbol:    symbol,
                Volume:    volume,
                Price:     tick.Ask,
                OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
            },
        })
        if err != nil {
            continue
        }

        result := data.MqlTradeCheckResult
        if result.ReturnedCode == 10009 && result.FreeMargin > 0 {
            fmt.Printf("Max lot size: %.2f (Required margin: %.2f)\n",
                volume, result.Margin)
            return volume, nil
        }
    }

    return 0, fmt.Errorf("no valid lot size found")
}
```

### 4) Validate and send pattern

```go
func ValidateAndSend(account *mt5.MT5Account, symbol string, volume float64, isBuy bool) (uint64, error) {
    ctx := context.Background()

    tick, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})

    var orderType pb.TMT5_ENUM_ORDER_TYPE
    var orderTypeTF pb.ENUM_ORDER_TYPE_TF
    var price, sl, tp float64

    if isBuy {
        orderType = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY
        orderTypeTF = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
        price = tick.Ask
        sl = tick.Ask - 0.0050
        tp = tick.Ask + 0.0100
    } else {
        orderType = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL
        orderTypeTF = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
        price = tick.Bid
        sl = tick.Bid + 0.0050
        tp = tick.Bid - 0.0100
    }

    // Validate first
    checkData, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
        MqlTradeRequest: &pb.MrpcMqlTradeRequest{
            Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
            Symbol:     symbol,
            Volume:     volume,
            Price:      price,
            OrderType:  orderTypeTF,
            StopLoss:   sl,
            TakeProfit: tp,
        },
    })
    if err != nil {
        return 0, fmt.Errorf("check failed: %w", err)
    }

    result := checkData.MqlTradeCheckResult
    if result.ReturnedCode != 10009 {
        return 0, fmt.Errorf("order invalid: %s", result.Comment)
    }

    // If valid, send order
    slippage := uint64(20)
    sendData, err := account.OrderSend(ctx, &pb.OrderSendRequest{
        Symbol:     symbol,
        Operation:  orderType,
        Volume:     volume,
        Price:      &price,
        StopLoss:   &sl,
        TakeProfit: &tp,
        Slippage:   &slippage,
    })
    if err != nil {
        return 0, err
    }

    if sendData.ReturnedCode == 10009 {
        fmt.Printf("Order validated and sent: %d\n", sendData.Order)
        return sendData.Order, nil
    }

    return 0, fmt.Errorf("order send failed: %s", sendData.ReturnedCodeDescription)
}
```

### 5) Check stop levels validity

```go
func ValidateStopLevels(account *mt5.MT5Account, symbol string, price, sl, tp float64) error {
    ctx := context.Background()

    data, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
        MqlTradeRequest: &pb.MrpcMqlTradeRequest{
            Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
            Symbol:     symbol,
            Volume:     0.01,
            Price:      price,
            OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
            StopLoss:   sl,
            TakeProfit: tp,
        },
    })
    if err != nil {
        return err
    }

    result := data.MqlTradeCheckResult
    if result.ReturnedCode != 10009 {
        return fmt.Errorf("invalid SL/TP levels: %s", result.Comment)
    }

    fmt.Println("Stop levels are valid")
    return nil
}
```

---

## 🔧 Common Patterns

### Safe order placement

```go
func SafeOrderSend(account *mt5.MT5Account, symbol string, operation pb.TMT5_ENUM_ORDER_TYPE, volume float64) (uint64, error) {
    ctx := context.Background()

    // Get current price
    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})
    if err != nil {
        return 0, err
    }

    var price, sl, tp float64
    var orderTypeTF pb.ENUM_ORDER_TYPE_TF

    if operation == pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY {
        orderTypeTF = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
        price = tick.Ask
        sl = tick.Ask - 0.0050
        tp = tick.Ask + 0.0100
    } else {
        orderTypeTF = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
        price = tick.Bid
        sl = tick.Bid + 0.0050
        tp = tick.Bid - 0.0100
    }

    // Validate first
    checkData, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
        MqlTradeRequest: &pb.MrpcMqlTradeRequest{
            Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
            Symbol:     symbol,
            Volume:     volume,
            Price:      price,
            OrderType:  orderTypeTF,
            StopLoss:   sl,
            TakeProfit: tp,
        },
    })
    if err != nil {
        return 0, fmt.Errorf("validation error: %w", err)
    }

    result := checkData.MqlTradeCheckResult
    if result.ReturnedCode != 10009 {
        return 0, fmt.Errorf("invalid order: %s", result.Comment)
    }

    // Send order
    slippage := uint64(20)
    sendData, err := account.OrderSend(ctx, &pb.OrderSendRequest{
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

    if sendData.ReturnedCode == 10009 {
        return sendData.Order, nil
    }

    return 0, fmt.Errorf("order send failed: %s", sendData.ReturnedCodeDescription)
}
```

---

## 📚 See Also

* [OrderSend](./OrderSend.md) - Place orders after validation
* [OrderCalcMargin](./OrderCalcMargin.md) - Calculate required margin
* [OrderCalcProfit](./OrderCalcProfit.md) - Calculate potential profit
* [AccountSummary](../1.%20Account_information/AccountSummary.md) - Check account balance and margin
