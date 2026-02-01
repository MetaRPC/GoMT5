# ✅ Stream Real-Time Trade Events

> **Request:** subscribe to real-time notifications for all trading operations: order placement, modification, execution, and cancellation.

**API Information:**

* **Low-level API:** `MT5Account.OnTrade(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.SubscriptionService`
* **Proto definition:** `OnTrade` (defined in `mt5-term-api-subscription.proto`)

### RPC

* **Service:** `mt5_term_api.SubscriptionService`
* **Method:** `OnTrade(OnTradeRequest) → stream OnTradeReply`
* **Low‑level client (generated):** `SubscriptionServiceClient.OnTrade(ctx, request, opts...)`

```go
package mt5

type MT5Account struct {
    // ...
}

// OnTrade streams trade events in real-time.
// Returns two channels: data channel and error channel.
func (a *MT5Account) OnTrade(
    ctx context.Context,
    req *pb.OnTradeRequest,
) (<-chan *pb.OnTradeData, <-chan error)
```

---

## 🔽 Input

| Parameter | Type                  | Description                                   |
| --------- | --------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`     | Context for cancellation (cancel to stop stream) |
| `req`     | `*pb.OnTradeRequest`  | Request (empty structure)                     |

---

## ⬆️ Output — Channels

| Channel      | Type                      | Description                              |
| ------------ | ------------------------- | ---------------------------------------- |
| Data Channel | `<-chan *pb.OnTradeData`  | Receives trade event updates             |
| Error Channel| `<-chan error`            | Receives errors (closed on ctx cancel)   |

**OnTradeData fields:**

| Field          | Type     | Go Type   | Description                    |
| -------------- | -------- | --------- | ------------------------------ |
| `Type` | `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (enum) | `int32` | Event type (always OnTrade) - **ENUM!** |
| `EventData` | `OnTadeEventData` | `*OnTadeEventData` | Trade event data (positions, orders, deals) |
| `AccountInfo` | `OnEventAccountInfo` | `*OnEventAccountInfo` | Account snapshot (balance, equity, etc.) |
| `TerminalInstanceGuidId` | `string` | `string` | Terminal instance ID |

---

## 💬 Just the essentials

* **What it is.** Real-time stream of all trading operations.
* **Why you need it.** Monitor order execution, track trading activity, implement notifications.
* **Event-driven.** Notifies on order placement, modification, execution, cancellation.

---

## 🎯 Purpose

Use it to:

* Monitor all trading activity in real-time
* Track order placement, modification, and execution
* Implement trade notifications and alerts
* Log trading operations
* Build trading dashboards
* Detect order fills and rejections immediately

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OnTrade - How it works](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnTrade_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, streams run indefinitely until cancelled.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Channel buffering:** Data channel is unbuffered, error channel is buffered (size 1).
* **Goroutine required:** You MUST consume the channels in a separate goroutine to avoid blocking.
* **Context cancellation:** Use `context.WithCancel()` or `context.WithTimeout()` to stop streaming.
* **All events:** Receives events for ALL trading operations (positions, orders, deals).
* **Empty request:** OnTradeRequest is an empty structure (no parameters needed).
* **ENUM type:** Always check the `Type` field (MT5_SUB_ENUM_EVENT_GROUP_TYPE) to identify event type.

---

## 🧱 ENUMs used in OnTrade

### Output ENUM

| ENUM Type | Field Name | Purpose | Values |
|-----------|------------|---------|--------|
| `MT5_SUB_ENUM_EVENT_GROUP_TYPE` | `Type` | Indicates the event type | `OrderUpdate` (1) - Trade event notification |

**Usage in code:**

```go
tradeStream, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

for event := range tradeStream {
    // Check event type (always OrderUpdate for this stream)
    switch event.Type {
    case pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate:
        fmt.Println("Trade event received")
        // Process trade event data
    }
}
```

**Note:** The `Type` field will always be `OrderUpdate` for OnTrade stream. This ENUM is shared across all streaming methods (OnTrade, OnPositionProfit, OnTradeTransaction) to maintain consistent event identification.

---

### Structure and ENUMs Usage

**OnTradeData Structure:**

```
OnTradeData
├── Type (MT5_SUB_ENUM_EVENT_GROUP_TYPE) - Event group type
├── EventData (OnTadeEventData)
│   ├── NewOrders ([]*OnTradeOrderInfo)
│   │   └── Uses: SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE,
│   │              SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING,
│   │              SUB_ENUM_ORDER_REASON
│   ├── DisappearedOrders ([]*OnTradeOrderInfo)
│   │   └── Uses: SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE,
│   │              SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING,
│   │              SUB_ENUM_ORDER_REASON
│   ├── StateChangedOrders ([]*OnTradeOrderStateChange)
│   │   └── Contains: PreviousOrder, CurrentOrder (OnTradeOrderInfo)
│   │   └── Uses: SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE,
│   │              SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING,
│   │              SUB_ENUM_ORDER_REASON
│   ├── NewHistoryOrders ([]*OnTradeHistoryOrderInfo)
│   │   └── Uses: SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE,
│   │              SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING,
│   │              SUB_ENUM_DEAL_REASON
│   ├── DisappearedHistoryOrders ([]*OnTradeHistoryOrderInfo)
│   │   └── Uses: SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE,
│   │              SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING,
│   │              SUB_ENUM_DEAL_REASON
│   ├── UpdatedHistoryOrders ([]*OnTradeHistoryOrderUpdate)
│   │   └── Contains: PreviousHistoryOrder, CurrentHistoryOrder (OnTradeHistoryOrderInfo)
│   │   └── Uses: SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE,
│   │              SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING,
│   │              SUB_ENUM_DEAL_REASON
│   ├── NewHistoryDeals ([]*OnTradeHistoryDealInfo)
│   │   └── Uses: SUB_ENUM_DEAL_TYPE, SUB_ENUM_DEAL_ENTRY, SUB_ENUM_DEAL_REASON
│   ├── DisappearedHistoryDeals ([]*OnTradeHistoryDealInfo)
│   │   └── Uses: SUB_ENUM_DEAL_TYPE, SUB_ENUM_DEAL_ENTRY, SUB_ENUM_DEAL_REASON
│   ├── UpdatedHistoryDeals ([]*OnTradeHistoryDealUpdate)
│   │   └── Contains: PreviousHistoryDeal, CurrentHistoryDeal (OnTradeHistoryDealInfo)
│   │   └── Uses: SUB_ENUM_DEAL_TYPE, SUB_ENUM_DEAL_ENTRY, SUB_ENUM_DEAL_REASON
│   ├── NewPositions ([]*OnTradePositionInfo)
│   │   └── Uses: SUB_ENUM_POSITION_TYPE, SUB_ENUM_POSITION_REASON
│   ├── DisappearedPositions ([]*OnTradePositionInfo)
│   │   └── Uses: SUB_ENUM_POSITION_TYPE, SUB_ENUM_POSITION_REASON
│   └── UpdatedPositions ([]*OnTradePositionUpdate)
│       └── Contains: PreviousPosition, CurrentPosition (OnTradePositionInfo)
│       └── Uses: SUB_ENUM_POSITION_TYPE, SUB_ENUM_POSITION_REASON
├── AccountInfo (*OnEventAccountInfo)
└── TerminalInstanceGuidId (string)
```

---

### Event Group Type Enum (MT5_SUB_ENUM_EVENT_GROUP_TYPE)

Used in `OnTradeData.Type` field:

| Name           | Value | Description                          |
| -------------- | ----- | ------------------------------------ |
| `OrderProfit`  | 0     | Order profit event                   |
| `OrderUpdate`  | 1     | Order update event                   |

**Usage:**
```go
pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit  // = 0
pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate  // = 1
```

---

### ENUMs Used in Nested Structures

The following enums are used by the nested message types inside `OnTadeEventData`:

| Nested Structure              | Used in Fields                                    | ENUMs Used                                                                                                      |
| ----------------------------- | ------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `OnTradePositionInfo`         | `NewPositions`, `DisappearedPositions`            | `SUB_ENUM_POSITION_TYPE`, `SUB_ENUM_POSITION_REASON`                                                           |
| `OnTradePositionUpdate`       | `UpdatedPositions`                                | `SUB_ENUM_POSITION_TYPE`, `SUB_ENUM_POSITION_REASON`                                                           |
| `OnTradeOrderInfo`            | `NewOrders`, `DisappearedOrders`                  | `SUB_ENUM_ORDER_TYPE`, `SUB_ENUM_ORDER_STATE`, `SUB_ENUM_ORDER_TYPE_TIME`, `SUB_ENUM_ORDER_TYPE_FILLING`, `SUB_ENUM_ORDER_REASON` |
| `OnTradeOrderStateChange`     | `StateChangedOrders`                              | `SUB_ENUM_ORDER_TYPE`, `SUB_ENUM_ORDER_STATE`, `SUB_ENUM_ORDER_TYPE_TIME`, `SUB_ENUM_ORDER_TYPE_FILLING`, `SUB_ENUM_ORDER_REASON` |
| `OnTradeHistoryOrderInfo`     | `NewHistoryOrders`, `DisappearedHistoryOrders`    | `SUB_ENUM_ORDER_TYPE`, `SUB_ENUM_ORDER_STATE`, `SUB_ENUM_ORDER_TYPE_TIME`, `SUB_ENUM_ORDER_TYPE_FILLING`, `SUB_ENUM_DEAL_REASON`  |
| `OnTradeHistoryOrderUpdate`   | `UpdatedHistoryOrders`                            | `SUB_ENUM_ORDER_TYPE`, `SUB_ENUM_ORDER_STATE`, `SUB_ENUM_ORDER_TYPE_TIME`, `SUB_ENUM_ORDER_TYPE_FILLING`, `SUB_ENUM_DEAL_REASON`  |
| `OnTradeHistoryDealInfo`      | `NewHistoryDeals`, `DisappearedHistoryDeals`      | `SUB_ENUM_DEAL_TYPE`, `SUB_ENUM_DEAL_ENTRY`, `SUB_ENUM_DEAL_REASON`                                            |
| `OnTradeHistoryDealUpdate`    | `UpdatedHistoryDeals`                             | `SUB_ENUM_DEAL_TYPE`, `SUB_ENUM_DEAL_ENTRY`, `SUB_ENUM_DEAL_REASON`                                            |

---

### Position Type Enum (SUB_ENUM_POSITION_TYPE)

| Name                       | Value | Description                          |
| -------------------------- | ----- | ------------------------------------ |
| `SUB_POSITION_TYPE_BUY`    | 0     | Buy position                         |
| `SUB_POSITION_TYPE_SELL`   | 1     | Sell position                        |

**Usage:**
```go
pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_BUY   // = 0
pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_SELL  // = 1
```

---

### Position Reason Enum (SUB_ENUM_POSITION_REASON)

| Name                           | Value | Description                          |
| ------------------------------ | ----- | ------------------------------------ |
| `SUB_POSITION_REASON_CLIENT`   | 0     | Position opened from desktop terminal|
| `SUB_POSITION_REASON_MOBILE`   | 2     | Position opened from mobile app      |
| `SUB_POSITION_REASON_WEB`      | 3     | Position opened from web terminal    |
| `SUB_POSITION_REASON_EXPERT`   | 4     | Position opened by Expert Advisor    |

**Usage:**
```go
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_CLIENT  // = 0
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_MOBILE  // = 2
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_WEB     // = 3
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_EXPERT  // = 4
```

---

### Order Type Enum (SUB_ENUM_ORDER_TYPE)

| Name                              | Value | Description                          |
| --------------------------------- | ----- | ------------------------------------ |
| `SUB_ORDER_TYPE_BUY`              | 0     | Market buy order                     |
| `SUB_ORDER_TYPE_SELL`             | 1     | Market sell order                    |
| `SUB_ORDER_TYPE_BUY_LIMIT`        | 2     | Buy limit pending order              |
| `SUB_ORDER_TYPE_SELL_LIMIT`       | 3     | Sell limit pending order             |
| `SUB_ORDER_TYPE_BUY_STOP`         | 4     | Buy stop pending order               |
| `SUB_ORDER_TYPE_SELL_STOP`        | 5     | Sell stop pending order              |
| `SUB_ORDER_TYPE_BUY_STOP_LIMIT`   | 6     | Buy stop limit pending order         |
| `SUB_ORDER_TYPE_SELL_STOP_LIMIT`  | 7     | Sell stop limit pending order        |
| `SUB_ORDER_TYPE_CLOSE_BY`         | 8     | Close by opposite position           |

**Usage:**
```go
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY             // = 0
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL            // = 1
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_LIMIT       // = 2
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_LIMIT      // = 3
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_STOP        // = 4
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_STOP       // = 5
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_STOP_LIMIT  // = 6
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_STOP_LIMIT // = 7
pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_CLOSE_BY        // = 8
```

---

### Order State Enum (SUB_ENUM_ORDER_STATE)

| Name                              | Value | Description                          |
| --------------------------------- | ----- | ------------------------------------ |
| `SUB_ORDER_STATE_STARTED`         | 0     | Order checked, but not yet accepted  |
| `SUB_ORDER_STATE_PLACED`          | 1     | Order accepted                       |
| `SUB_ORDER_STATE_CANCELED`        | 2     | Order canceled by client             |
| `SUB_ORDER_STATE_PARTIAL`         | 3     | Order partially executed             |
| `SUB_ORDER_STATE_FILLED`          | 4     | Order fully executed                 |
| `SUB_ORDER_STATE_REJECTED`        | 5     | Order rejected                       |
| `SUB_ORDER_STATE_EXPIRED`         | 6     | Order expired                        |
| `SUB_ORDER_STATE_REQUEST_ADD`     | 7     | Order being registered               |
| `SUB_ORDER_STATE_REQUEST_MODIFY`  | 8     | Order being modified                 |
| `SUB_ORDER_STATE_REQUEST_CANCEL`  | 9     | Order being deleted                  |

**Usage:**
```go
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_STARTED        // = 0
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_PLACED         // = 1
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_CANCELED       // = 2
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_PARTIAL        // = 3
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_FILLED         // = 4
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REJECTED       // = 5
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_EXPIRED        // = 6
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_ADD    // = 7
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_MODIFY // = 8
pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_CANCEL // = 9
```

---

### Order Time Type Enum (SUB_ENUM_ORDER_TYPE_TIME)

| Name                           | Value | Description                          |
| ------------------------------ | ----- | ------------------------------------ |
| `SUB_ORDER_TIME_GTC`           | 0     | Good till cancel                     |
| `SUB_ORDER_TIME_DAY`           | 1     | Good till current trading day        |
| `SUB_ORDER_TIME_SPECIFIED`     | 2     | Good till specified date             |
| `SUB_ORDER_TIME_SPECIFIED_DAY` | 3     | Good till specified day              |

**Usage:**
```go
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_GTC           // = 0
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_DAY           // = 1
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED     // = 2
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED_DAY // = 3
```

---

### Order Filling Type Enum (SUB_ENUM_ORDER_TYPE_FILLING)

| Name                         | Value | Description                          |
| ---------------------------- | ----- | ------------------------------------ |
| `SUB_ORDER_FILLING_FOK`      | 0     | Fill or Kill                         |
| `SUB_ORDER_FILLING_IOC`      | 1     | Immediate or Cancel                  |
| `SUB_ORDER_FILLING_BOC`      | 2     | Book or Cancel                       |
| `SUB_ORDER_FILLING_RETURN`   | 3     | Return                               |

**Usage:**
```go
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_FOK    // = 0
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_IOC    // = 1
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_BOC    // = 2
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_RETURN // = 3
```

---

### Order Reason Enum (SUB_ENUM_ORDER_REASON)

| Name                         | Value | Description                          |
| ---------------------------- | ----- | ------------------------------------ |
| `SUB_ORDER_REASON_CLIENT`    | 0     | Order placed from desktop terminal   |
| `SUB_ORDER_REASON_MOBILE`    | 2     | Order placed from mobile app         |
| `SUB_ORDER_REASON_WEB`       | 3     | Order placed from web terminal       |
| `SUB_ORDER_REASON_EXPERT`    | 4     | Order placed by Expert Advisor       |
| `SUB_ORDER_REASON_SL`        | 5     | Order triggered by Stop Loss         |
| `SUB_ORDER_REASON_TP`        | 6     | Order triggered by Take Profit       |
| `SUB_ORDER_REASON_SO`        | 7     | Order triggered by Stop Out          |

**Usage:**
```go
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_CLIENT // = 0
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_MOBILE // = 2
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_WEB    // = 3
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_EXPERT // = 4
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_SL     // = 5
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_TP     // = 6
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_SO     // = 7
```

---

### Deal Type Enum (SUB_ENUM_DEAL_TYPE)

| Name                                  | Value | Description                          |
| ------------------------------------- | ----- | ------------------------------------ |
| `SUB_DEAL_TYPE_BUY`                   | 0     | Buy deal                             |
| `SUB_DEAL_TYPE_SELL`                  | 1     | Sell deal                            |
| `SUB_DEAL_TYPE_BALANCE`               | 2     | Balance operation                    |
| `SUB_DEAL_TYPE_CREDIT`                | 3     | Credit operation                     |
| `SUB_DEAL_TYPE_CHARGE`                | 4     | Additional charge                    |
| `SUB_DEAL_TYPE_CORRECTION`            | 5     | Correction                           |
| `SUB_DEAL_TYPE_BONUS`                 | 6     | Bonus                                |
| `SUB_DEAL_TYPE_COMMISSION`            | 7     | Additional commission                |
| `SUB_DEAL_TYPE_COMMISSION_DAILY`      | 8     | Daily commission                     |
| `SUB_DEAL_TYPE_COMMISSION_MONTHLY`    | 9     | Monthly commission                   |
| `SUB_DEAL_TYPE_COMMISSION_AGENT_DAILY`   | 10    | Daily agent commission            |
| `SUB_DEAL_TYPE_COMMISSION_AGENT_MONTHLY` | 11    | Monthly agent commission          |
| `SUB_DEAL_TYPE_INTEREST`              | 12    | Interest rate                        |
| `SUB_DEAL_TYPE_BUY_CANCELED`          | 13    | Canceled buy deal                    |
| `SUB_DEAL_TYPE_SELL_CANCELED`         | 14    | Canceled sell deal                   |
| `SUB_DEAL_DIVIDEND`                   | 15    | Dividend operations                  |
| `SUB_DEAL_DIVIDEND_FRANKED`           | 16    | Franked (non-taxable) dividend       |
| `SUB_DEAL_TAX`                        | 17    | Tax charges                          |

**Usage:**
```go
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY                      // = 0
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_SELL                     // = 1
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BALANCE                  // = 2
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_CREDIT                   // = 3
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_CHARGE                   // = 4
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_CORRECTION               // = 5
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BONUS                    // = 6
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION               // = 7
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_DAILY         // = 8
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_MONTHLY       // = 9
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_AGENT_DAILY   // = 10
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_AGENT_MONTHLY // = 11
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_INTEREST                 // = 12
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY_CANCELED             // = 13
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_SELL_CANCELED            // = 14
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_DIVIDEND                      // = 15
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_DIVIDEND_FRANKED              // = 16
pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TAX                           // = 17
```

---

### Deal Entry Type Enum (SUB_ENUM_DEAL_ENTRY)

| Name                       | Value | Description                          |
| -------------------------- | ----- | ------------------------------------ |
| `SUB_DEAL_ENTRY_IN`        | 0     | Entry into market                    |
| `SUB_DEAL_ENTRY_OUT`       | 1     | Exit from market                     |
| `SUB_DEAL_ENTRY_INOUT`     | 2     | Reverse                              |
| `SUB_DEAL_ENTRY_OUT_BY`    | 3     | Close by opposite position           |

**Usage:**
```go
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_IN     // = 0
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_OUT    // = 1
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_INOUT  // = 2
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_OUT_BY // = 3
```

---

### Deal Reason Enum (SUB_ENUM_DEAL_REASON)

| Name                                 | Value | Description                          |
| ------------------------------------ | ----- | ------------------------------------ |
| `SUB_DEAL_REASON_CLIENT`             | 0     | Deal from desktop terminal           |
| `SUB_DEAL_REASON_MOBILE`             | 1     | Deal from mobile app                 |
| `SUB_DEAL_REASON_WEB`                | 2     | Deal from web terminal               |
| `SUB_DEAL_REASON_EXPERT`             | 3     | Deal by Expert Advisor               |
| `SUB_DEAL_REASON_SL`                 | 4     | Deal by Stop Loss                    |
| `SUB_DEAL_REASON_TP`                 | 5     | Deal by Take Profit                  |
| `SUB_DEAL_REASON_SO`                 | 6     | Deal by Stop Out                     |
| `SUB_DEAL_REASON_ROLLOVER`           | 7     | Deal by rollover                     |
| `SUB_DEAL_REASON_VMARGIN`            | 8     | Deal by variation margin             |
| `SUB_DEAL_REASON_SPLIT`              | 9     | Deal by split                        |
| `SUB_DEAL_REASON_CORPORATE_ACTION`   | 10    | Deal by corporate action             |

**Usage:**
```go
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_CLIENT           // = 0
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_MOBILE           // = 1
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_WEB              // = 2
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_EXPERT           // = 3
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_SL               // = 4
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_TP               // = 5
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_SO               // = 6
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_ROLLOVER         // = 7
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_VMARGIN          // = 8
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_SPLIT            // = 9
pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_CORPORATE_ACTION // = 10
```

---

## 🔗 Usage Examples

### 1) Basic trade event monitoring

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

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }
                fmt.Printf("[%s] Trade event received\n",
                    time.Now().Format("15:04:05"))

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Stream error: %v\n", err)
                    return
                }
            }
        }
    }()

    // Keep running
    time.Sleep(60 * time.Second)
}
```

### 2) Trade notification system

```go
func TradeNotificationSystem(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    fmt.Println("Trade stream closed")
                    return
                }

                fmt.Printf("\n🔔 Trade Event Notification\n")
                fmt.Printf("  Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
                // Display event details based on event type

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("❌ Stream error: %v\n", err)
                    return
                }

            case <-ctx.Done():
                fmt.Println("Notification system stopped")
                return
            }
        }
    }()

    // Run indefinitely
    select {}
}
```

### 3) Trade logger with file output

```go
func TradeLogger(account *mt5.MT5Account, logFile string) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Failed to open log file: %v\n", err)
        return
    }
    defer file.Close()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }

                logEntry := fmt.Sprintf("[%s] Trade event\n",
                    time.Now().Format("2006-01-02 15:04:05"))
                file.WriteString(logEntry)
                fmt.Print(logEntry)

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Error: %v\n", err)
                    return
                }

            case <-ctx.Done():
                return
            }
        }
    }()

    // Run for extended period
    time.Sleep(24 * time.Hour)
}
```

### 4) Trade event counter

```go
type TradeEventStats struct {
    TotalEvents int64
    StartTime   time.Time
}

func CountTradeEvents(account *mt5.MT5Account, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    stats := &TradeEventStats{
        StartTime: time.Now(),
    }

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    for {
        select {
        case event := <-dataChan:
            if event == nil {
                break
            }
            stats.TotalEvents++
            fmt.Printf("\r Trade events: %d", stats.TotalEvents)

        case err := <-errChan:
            if err != nil {
                fmt.Printf("\nError: %v\n", err)
            }

        case <-ctx.Done():
            elapsed := time.Since(stats.StartTime)
            fmt.Printf("\n\nStatistics:\n")
            fmt.Printf("  Duration: %v\n", elapsed)
            fmt.Printf("  Total events: %d\n", stats.TotalEvents)
            fmt.Printf("  Events/minute: %.2f\n",
                float64(stats.TotalEvents)/(elapsed.Minutes()))
            return
        }
    }
}

// Usage:
// CountTradeEvents(account, 5*time.Minute)
```

### 5) Real-time trade alert system

```go
func RealTimeTradeAlerts(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }

                // Send alert (email, Telegram, etc)
                fmt.Printf("🚨 TRADE ALERT: %s\n", time.Now().Format("15:04:05"))

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Alert system error: %v\n", err)
                    return
                }

            case <-ctx.Done():
                fmt.Println("Alert system shutdown")
                return
            }
        }
    }()

    // Keep running
    select {}
}
```

---

## 📚 See Also

* [OnTradeTransaction](./OnTradeTransaction.md) - Detailed transaction-level events
* [OnPositionsAndPendingOrdersTickets](./OnPositionsAndPendingOrdersTickets.md) - Stream ticket changes
* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get current positions
