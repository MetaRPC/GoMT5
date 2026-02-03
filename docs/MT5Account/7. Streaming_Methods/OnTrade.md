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

### OnTradeData Structure

```
OnTradeData
├── Type (MT5_SUB_ENUM_EVENT_GROUP_TYPE) - Event type identifier
├── EventData (OnTadeEventData) - Trade events data
├── AccountInfo (*OnEventAccountInfo) - Account snapshot
└── TerminalInstanceGuidId (string) - Terminal instance ID
```

**Fields:**

| Field          | Type     | Go Type   | Description                    |
| -------------- | -------- | --------- | ------------------------------ |
| `Type` | `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (enum) | `int32` | Event type (always `OrderUpdate` = 1 for OnTrade) |
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
* **ENUM constant names:** Use full names like `pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY` - no shortcuts available.

---

## 🧱 ENUMs Overview

OnTrade uses **11 ENUMs in total**:

- **1 Direct ENUM** - returned by the method itself in `OnTradeData.Type`
- **10 Nested ENUMs** - used in nested data structures within `OnTadeEventData`

### Direct Method ENUM

| ENUM | Field | Always Returns | Description |
|------|-------|----------------|-------------|
| `MT5_SUB_ENUM_EVENT_GROUP_TYPE` | `Type` | `OrderUpdate` (1) | Event type identifier (shared across streaming methods) |

**Important Notes:**

- This ENUM has 2 values: `OrderProfit` (0) and `OrderUpdate` (1)
- **OnTrade always returns `OrderUpdate` (1)**
- `OrderProfit` (0) is used by [OnPositionProfit](./OnPositionProfit.md) method
- This shared ENUM allows distinguishing between different streaming event types

### Nested ENUMs (10 total)

These ENUMs belong to nested structures inside `OnTadeEventData`:

| Category | ENUMs | Used In |
|----------|-------|---------|
| **Positions** (2) | `SUB_ENUM_POSITION_TYPE`<br>`SUB_ENUM_POSITION_REASON` | OnTradePositionInfo |
| **Orders** (5) | `SUB_ENUM_ORDER_TYPE`<br>`SUB_ENUM_ORDER_STATE`<br>`SUB_ENUM_ORDER_TYPE_TIME`<br>`SUB_ENUM_ORDER_TYPE_FILLING`<br>`SUB_ENUM_ORDER_REASON` | OnTradeOrderInfo |
| **Deals** (3) | `SUB_ENUM_DEAL_TYPE`<br>`SUB_ENUM_DEAL_ENTRY`<br>`SUB_ENUM_DEAL_REASON` | OnTradeHistoryDealInfo |

---

## 📊 Data Structure & ENUM Usage

### Complete Structure Tree

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

### Field Order in OnTadeEventData

**Important:** Fields appear in this exact order in the protobuf structure:

1. **Orders** → NewOrders, DisappearedOrders, StateChangedOrders
2. **History Orders** → NewHistoryOrders, DisappearedHistoryOrders, UpdatedHistoryOrders
3. **History Deals** → NewHistoryDeals, DisappearedHistoryDeals, UpdatedHistoryDeals
4. **Positions** → NewPositions, DisappearedPositions, UpdatedPositions

---

## 📚 ENUM Reference

### How to Use ENUMs in Go

**Important:** Protobuf ENUMs in Go require **full constant names**:

```go
// Structure: pb.[ENUM_TYPE]_[VALUE_NAME]
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_CLIENT
```

**Why so long?**

- Standard protobuf naming for Go
- Ensures uniqueness across the package
- No shorter alternatives available
- IDE autocomplete works perfectly with these names

**Below are all available ENUMs with their values and usage examples.**

---

### 1. Event Group Type (MT5_SUB_ENUM_EVENT_GROUP_TYPE)

**Used in:** `OnTradeData.Type` field

| Name           | Value | Description                          | Used By Method |
| -------------- | ----- | ------------------------------------ | -------------- |
| `OrderProfit`  | 0     | Order profit event                   | OnPositionProfit |
| `OrderUpdate`  | 1     | Order update event                   | **OnTrade** ✓ |

```go
pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit  // = 0
pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate  // = 1

// Example: Check event type
tradeStream, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})
for event := range tradeStream {
    // For OnTrade, this will ALWAYS be OrderUpdate
    if event.Type == pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate {
        fmt.Println("Trade event received")
    }
}
```

---

### 2. Position Type (SUB_ENUM_POSITION_TYPE)

**Used in:** `OnTradePositionInfo.Type` field

| Name                       | Value | Description                          |
| -------------------------- | ----- | ------------------------------------ |
| `SUB_POSITION_TYPE_BUY`    | 0     | Buy position                         |
| `SUB_POSITION_TYPE_SELL`   | 1     | Sell position                        |

```go
pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_BUY   // = 0
pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_SELL  // = 1

// Example:
for _, pos := range event.EventData.NewPositions {
    if pos.Type == pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_BUY {
        fmt.Printf("New BUY position: %s\n", pos.Symbol)
    }
}
```

---

### 3. Position Reason (SUB_ENUM_POSITION_REASON)

**Used in:** `OnTradePositionInfo.Reason` field

| Name                           | Value | Description                          |
| ------------------------------ | ----- | ------------------------------------ |
| `SUB_POSITION_REASON_CLIENT`   | 0     | Position opened from desktop terminal|
| `SUB_POSITION_REASON_MOBILE`   | 2     | Position opened from mobile app      |
| `SUB_POSITION_REASON_WEB`      | 3     | Position opened from web terminal    |
| `SUB_POSITION_REASON_EXPERT`   | 4     | Position opened by Expert Advisor    |

```go
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_CLIENT  // = 0
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_MOBILE  // = 2
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_WEB     // = 3
pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_EXPERT  // = 4

// Example:
for _, pos := range event.EventData.NewPositions {
    if pos.Reason == pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_EXPERT {
        fmt.Printf("Position opened by EA: %s\n", pos.Symbol)
    }
}
```

---

### 4. Order Type (SUB_ENUM_ORDER_TYPE)

**Used in:** `OnTradeOrderInfo.OrderType`, `OnTradeHistoryOrderInfo.OrderType` fields

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

// Example:
for _, order := range event.EventData.NewOrders {
    if order.OrderType == pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_LIMIT {
        fmt.Printf("New BUY LIMIT order: %s @ %.5f\n", order.Symbol, order.PriceOpen)
    }
}
```

---

### 5. Order State (SUB_ENUM_ORDER_STATE)

**Used in:** `OnTradeOrderInfo.State`, `OnTradeHistoryOrderInfo.State` fields

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

// Example:
for _, order := range event.EventData.StateChangedOrders {
    if order.CurrentOrder.State == pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_FILLED {
        fmt.Printf("Order FILLED: Ticket %d\n", order.CurrentOrder.Ticket)
    }
}
```

---

### 6. Order Time Type (SUB_ENUM_ORDER_TYPE_TIME)

**Used in:** `OnTradeOrderInfo.TimeType`, `OnTradeHistoryOrderInfo.TypeTime` fields

| Name                           | Value | Description                          |
| ------------------------------ | ----- | ------------------------------------ |
| `SUB_ORDER_TIME_GTC`           | 0     | Good till cancel                     |
| `SUB_ORDER_TIME_DAY`           | 1     | Good till current trading day        |
| `SUB_ORDER_TIME_SPECIFIED`     | 2     | Good till specified date             |
| `SUB_ORDER_TIME_SPECIFIED_DAY` | 3     | Good till specified day              |

```go
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_GTC           // = 0
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_DAY           // = 1
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED     // = 2
pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED_DAY // = 3

// Example:
for _, order := range event.EventData.NewOrders {
    if order.TimeType == pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_DAY {
        fmt.Printf("Day order placed: %s\n", order.Symbol)
    }
}
```

---

### 7. Order Filling Type (SUB_ENUM_ORDER_TYPE_FILLING)

**Used in:** `OnTradeOrderInfo.OrderTypeFilling`, `OnTradeHistoryOrderInfo.OrderTypeFilling` fields

| Name                         | Value | Description                          |
| ---------------------------- | ----- | ------------------------------------ |
| `SUB_ORDER_FILLING_FOK`      | 0     | Fill or Kill                         |
| `SUB_ORDER_FILLING_IOC`      | 1     | Immediate or Cancel                  |
| `SUB_ORDER_FILLING_BOC`      | 2     | Book or Cancel                       |
| `SUB_ORDER_FILLING_RETURN`   | 3     | Return                               |

```go
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_FOK    // = 0
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_IOC    // = 1
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_BOC    // = 2
pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_RETURN // = 3

// Example:
for _, order := range event.EventData.NewOrders {
    if order.OrderTypeFilling == pb.SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_FOK {
        fmt.Printf("FOK order: %s\n", order.Symbol)
    }
}
```

---

### 8. Order Reason (SUB_ENUM_ORDER_REASON)

**Used in:** `OnTradeOrderInfo.OrderReason` field

| Name                         | Value | Description                          |
| ---------------------------- | ----- | ------------------------------------ |
| `SUB_ORDER_REASON_CLIENT`    | 0     | Order placed from desktop terminal   |
| `SUB_ORDER_REASON_MOBILE`    | 2     | Order placed from mobile app         |
| `SUB_ORDER_REASON_WEB`       | 3     | Order placed from web terminal       |
| `SUB_ORDER_REASON_EXPERT`    | 4     | Order placed by Expert Advisor       |
| `SUB_ORDER_REASON_SL`        | 5     | Order triggered by Stop Loss         |
| `SUB_ORDER_REASON_TP`        | 6     | Order triggered by Take Profit       |
| `SUB_ORDER_REASON_SO`        | 7     | Order triggered by Stop Out          |

```go
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_CLIENT // = 0
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_MOBILE // = 2
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_WEB    // = 3
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_EXPERT // = 4
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_SL     // = 5
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_TP     // = 6
pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_SO     // = 7

// Example:
for _, order := range event.EventData.NewOrders {
    if order.OrderReason == pb.SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_SL {
        fmt.Printf("Stop Loss triggered: Ticket %d\n", order.Ticket)
    }
}
```

---

### 9. Deal Type (SUB_ENUM_DEAL_TYPE)

**Used in:** `OnTradeHistoryDealInfo.Type` field

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

// Example:
for _, deal := range event.EventData.NewHistoryDeals {
    if deal.Type == pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY {
        fmt.Printf("BUY deal executed: Ticket %d @ %.5f\n", deal.Ticket, deal.Price)
    }
}
```

---

### 10. Deal Entry Type (SUB_ENUM_DEAL_ENTRY)

**Used in:** `OnTradeHistoryDealInfo.Entry` field

| Name                       | Value | Description                          |
| -------------------------- | ----- | ------------------------------------ |
| `SUB_DEAL_ENTRY_IN`        | 0     | Entry into market                    |
| `SUB_DEAL_ENTRY_OUT`       | 1     | Exit from market                     |
| `SUB_DEAL_ENTRY_INOUT`     | 2     | Reverse                              |
| `SUB_DEAL_ENTRY_OUT_BY`    | 3     | Close by opposite position           |

```go
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_IN     // = 0
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_OUT    // = 1
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_INOUT  // = 2
pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_OUT_BY // = 3

// Example:
for _, deal := range event.EventData.NewHistoryDeals {
    if deal.Entry == pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_IN {
        fmt.Printf("Market entry: %s\n", deal.Symbol)
    }
}
```

---

### 11. Deal Reason (SUB_ENUM_DEAL_REASON)

**Used in:** `OnTradeHistoryDealInfo.Reason`, `OnTradeHistoryOrderInfo.Reason` fields

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

// Example:
for _, deal := range event.EventData.NewHistoryDeals {
    if deal.Reason == pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_TP {
        fmt.Printf("Take Profit hit: Position #%d @ %.5f\n", deal.DealPositionId, deal.Price)
    }
}
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

### 6) Complete event processing example

```go
func CompleteTradeEventProcessor(account *mt5.MT5Account) {
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

                // Process all order events
                for _, order := range event.EventData.NewOrders {
                    fmt.Printf("[NEW ORDER] Ticket: %d, Symbol: %s, Type: %d, State: %d\n",
                        order.Ticket, order.Symbol, order.OrderType, order.State)
                }

                for _, order := range event.EventData.DisappearedOrders {
                    fmt.Printf("[DISAPPEARED ORDER] Ticket: %d, Symbol: %s\n",
                        order.Ticket, order.Symbol)
                }

                for _, orderChange := range event.EventData.StateChangedOrders {
                    fmt.Printf("[ORDER STATE CHANGE] Ticket: %d, State: %d → %d\n",
                        orderChange.CurrentOrder.Ticket,
                        orderChange.PreviousOrder.State,
                        orderChange.CurrentOrder.State)

                    // Check if order was filled
                    if orderChange.CurrentOrder.State == pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_FILLED {
                        fmt.Printf("  ✓ Order FILLED: %s\n", orderChange.CurrentOrder.Symbol)
                    }
                }

                // Process history orders
                for _, histOrder := range event.EventData.NewHistoryOrders {
                    fmt.Printf("[NEW HISTORY ORDER] Ticket: %d, Symbol: %s, State: %d\n",
                        histOrder.Ticket, histOrder.Symbol, histOrder.State)
                }

                for _, histOrder := range event.EventData.DisappearedHistoryOrders {
                    fmt.Printf("[DISAPPEARED HISTORY ORDER] Ticket: %d\n", histOrder.Ticket)
                }

                for _, histUpdate := range event.EventData.UpdatedHistoryOrders {
                    fmt.Printf("[UPDATED HISTORY ORDER] Ticket: %d\n",
                        histUpdate.CurrentHistoryOrder.Ticket)
                }

                // Process history deals
                for _, deal := range event.EventData.NewHistoryDeals {
                    dealType := "UNKNOWN"
                    if deal.Type == pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY {
                        dealType = "BUY"
                    } else if deal.Type == pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_SELL {
                        dealType = "SELL"
                    }

                    entryType := "UNKNOWN"
                    if deal.Entry == pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_IN {
                        entryType = "ENTRY"
                    } else if deal.Entry == pb.SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_OUT {
                        entryType = "EXIT"
                    }

                    fmt.Printf("[NEW DEAL] Ticket: %d, Symbol: %s, Type: %s, Entry: %s, Price: %.5f, Volume: %.2f\n",
                        deal.Ticket, deal.Symbol, dealType, entryType, deal.Price, deal.Volume)

                    // Check if deal was triggered by TP/SL
                    if deal.Reason == pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_TP {
                        fmt.Printf("  ✓ Take Profit triggered!\n")
                    } else if deal.Reason == pb.SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_SL {
                        fmt.Printf("  ✗ Stop Loss triggered!\n")
                    }
                }

                for _, deal := range event.EventData.DisappearedHistoryDeals {
                    fmt.Printf("[DISAPPEARED DEAL] Ticket: %d\n", deal.Ticket)
                }

                for _, dealUpdate := range event.EventData.UpdatedHistoryDeals {
                    fmt.Printf("[UPDATED DEAL] Ticket: %d\n",
                        dealUpdate.CurrentHistoryDeal.Ticket)
                }

                // Process positions
                for _, pos := range event.EventData.NewPositions {
                    posType := "BUY"
                    if pos.Type == pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_SELL {
                        posType = "SELL"
                    }

                    fmt.Printf("[NEW POSITION] Ticket: %d, Symbol: %s, Type: %s, Volume: %.2f, Price: %.5f\n",
                        pos.Ticket, pos.Symbol, posType, pos.Volume, pos.PriceOpen)

                    // Check if position was opened by EA
                    if pos.Reason == pb.SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_EXPERT {
                        fmt.Printf("  ⚡ Opened by Expert Advisor\n")
                    }
                }

                for _, pos := range event.EventData.DisappearedPositions {
                    fmt.Printf("[CLOSED POSITION] Ticket: %d, Symbol: %s, Profit: %.2f\n",
                        pos.Ticket, pos.Symbol, pos.Profit)
                }

                for _, posUpdate := range event.EventData.UpdatedPositions {
                    fmt.Printf("[UPDATED POSITION] Ticket: %d, Symbol: %s\n",
                        posUpdate.CurrentPosition.Ticket,
                        posUpdate.CurrentPosition.Symbol)
                    fmt.Printf("  Volume: %.2f → %.2f\n",
                        posUpdate.PreviousPosition.Volume,
                        posUpdate.CurrentPosition.Volume)
                    fmt.Printf("  Profit: %.2f → %.2f\n",
                        posUpdate.PreviousPosition.Profit,
                        posUpdate.CurrentPosition.Profit)
                }

                // Display account info
                if event.AccountInfo != nil {
                    fmt.Printf("[ACCOUNT] Balance: %.2f, Equity: %.2f, Margin: %.2f, Free Margin: %.2f\n",
                        event.AccountInfo.Balance,
                        event.AccountInfo.Equity,
                        event.AccountInfo.Margin,
                        event.AccountInfo.FreeMargin)
                }

                fmt.Println("─────────────────────────────────────")

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Stream error: %v\n", err)
                    return
                }

            case <-ctx.Done():
                fmt.Println("Event processor stopped")
                return
            }
        }
    }()

    // Run indefinitely
    select {}
}
```

---

## 📚 See Also

* [OnTradeTransaction](./OnTradeTransaction.md) - Detailed transaction-level events
* [OnPositionsAndPendingOrdersTickets](./OnPositionsAndPendingOrdersTickets.md) - Stream ticket changes
* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get current positions
