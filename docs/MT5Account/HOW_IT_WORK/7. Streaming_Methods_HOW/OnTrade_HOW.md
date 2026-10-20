### Example from file: `examples/demos/lowlevel/03_streaming_methods.go`

> The **`OnTrade()`** method is used to receive **real-time trading events**.
> It notifies about all trading account changes: order opening and closing, new positions appearing, deals in history, etc.
>
> Unlike `OnSymbolTick()`, which streams quotes, `OnTrade()` transmits **trading events** occurring on the account.


---

## 🧩 Code example

```go
tradeReq := &pb.OnTradeRequest{}

tradeCtx, tradeCancel := context.WithCancel(ctx)
tradeChan, tradeErrChan := account.OnTrade(tradeCtx, tradeReq)

fmt.Printf("Streaming trade events (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
fmt.Println("  NOTE: This stream sends events only when trades occur")

eventCount = 0
timeout = time.After(MAX_SECONDS * time.Second)

streamTrade:
for {
    select {
    case tradeData, ok := <-tradeChan:
        if !ok {
            fmt.Println("  Stream closed by server")
            break streamTrade
        }
        eventCount++
        fmt.Printf("  Event #%d: Type=%v\n", eventCount, tradeData.Type)
        if tradeData.EventData != nil {
            ed := tradeData.EventData
            fmt.Printf("    New Orders: %d, Disappeared Orders: %d\n",
                len(ed.NewOrders), len(ed.DisappearedOrders))
            fmt.Printf("    New Positions: %d, Disappeared Positions: %d\n",
                len(ed.NewPositions), len(ed.DisappearedPositions))
            fmt.Printf("    New History Orders: %d, New History Deals: %d\n",
                len(ed.NewHistoryOrders), len(ed.NewHistoryDeals))
        }

        if eventCount >= MAX_EVENTS {
            fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
            break streamTrade
        }

    case err := <-tradeErrChan:
        if err != nil {
            helpers.PrintShortError(err, "Stream error")
            break streamTrade
        }

    case <-timeout:
        fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
        break streamTrade
    }
}
tradeCancel() // closes this stream on the server
```

---

## 🟢 Detailed Code Explanation

### 1️. Subscribe to Trade Event Stream

```go
tradeReq := &pb.OnTradeRequest{}
```

Create an empty request, as the `OnTrade()` stream doesn't require symbols — it listens to **all trading events on the account**.

---

### 2️. Start the Stream

```go
tradeCtx, tradeCancel := context.WithCancel(ctx)
tradeChan, tradeErrChan := account.OnTrade(tradeCtx, tradeReq)
```

The stream gets its own context, derived from `ctx`: `tradeCancel()` after the loop closes this stream on the server at once. Without it the stream would stay open, unread, until `ctx` itself ends — and every later stream of the demo would pile up next to it.

The method returns two channels:

* **tradeChan** — stream with trading events;
* **tradeErrChan** — error stream (e.g., if the server closes the connection).

> The stream is active only when events occur. If there are no trades — there will be no data.

---

### 3️. Display Informational Message

```go
fmt.Printf("Streaming trade events (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
fmt.Println("  NOTE: This stream sends events only when trades occur")
```

Informs the user that the stream won't be constantly active — events arrive **only when there are trading changes**.

---

### 4️. Event Processing Loop

```go
for {
    select {
    case tradeData, ok := <-tradeChan:
```

Uses the `select` construct to listen to multiple sources simultaneously:

* trading events (`tradeChan`),
* errors (`tradeErrChan`),
* timeout expiration (`timeout`).

---

### 5️. Process Trade Event

```go
fmt.Printf("  Event #%d: Type=%v\n", eventCount, tradeData.Type)
```

Display the event group type (`MT5_SUB_ENUM_EVENT_GROUP_TYPE`). For `OnTrade()` it is always `OrderUpdate` (1);
what actually changed is described by the lists in `EventData` (new/disappeared orders, positions, history orders and deals).

---

### 6️. Detailed Change Information

```go
if tradeData.EventData != nil {
    ed := tradeData.EventData
    fmt.Printf("    New Orders: %d, Disappeared Orders: %d\n",
        len(ed.NewOrders), len(ed.DisappearedOrders))
    fmt.Printf("    New Positions: %d, Disappeared Positions: %d\n",
        len(ed.NewPositions), len(ed.DisappearedPositions))
    fmt.Printf("    New History Orders: %d, New History Deals: %d\n",
        len(ed.NewHistoryOrders), len(ed.NewHistoryDeals))
}
```

Each event contains an **`EventData`** object, which indicates which elements changed:

| Field                   | Description                                        |
| ---------------------- | -------------------------------------------------- |
| `NewOrders`            | New orders that appeared in the system             |
| `DisappearedOrders`    | Orders that disappeared (deleted or executed)      |
| `NewPositions`         | New open positions                                 |
| `DisappearedPositions` | Positions that were closed                         |
| `NewHistoryOrders`     | New orders added to history                        |
| `NewHistoryDeals`      | New deals added to history                         |

`EventData` also has `StateChangedOrders`, `DisappearedHistoryOrders`, `UpdatedHistoryOrders`, `DisappearedHistoryDeals`, `UpdatedHistoryDeals` and `UpdatedPositions`.

> Thus, `OnTrade()` returns a **delta of changes** — allowing you to synchronize local data with the server.

---

### 7️. Stream Termination

The stream closes when the event limit (`MAX_EVENTS`) is reached, timeout expires, or a connection error occurs.

---

## 📦 What the Server Returns

```protobuf
message OnTradeData {
  MT5_SUB_ENUM_EVENT_GROUP_TYPE type = 1;  // Event type (OrderUpdate)
  OnTadeEventData event_data = 2;          // Change details
  OnEventAccountInfo account_info = 3;
  string terminal_instance_guid_id = 4;
}

message OnTadeEventData {
  repeated OnTradeOrderInfo new_orders = 1;
  repeated OnTradeOrderInfo disappeared_orders = 2;
  repeated OnTradeOrderStateChange state_changed_orders = 3;
  repeated OnTradeHistoryOrderInfo new_history_orders = 4;
  repeated OnTradeHistoryOrderInfo disappeared_history_orders = 5;
  repeated OnTradeHistoryOrderUpdate updated_history_orders = 6;
  repeated OnTradeHistoryDealInfo new_history_deals = 7;
  repeated OnTradeHistoryDealInfo disappeared_history_deals = 8;
  repeated OnTradeHistoryDealUpdate updated_history_deals = 9;
  repeated OnTradePositionInfo new_positions = 10;
  repeated OnTradePositionInfo disappeared_positions = 11;
  repeated OnTradePositionUpdate updated_positions = 12;
}
```

---

## 💡 Example Output

```
Streaming trade events (max 3 events or 10 seconds)...
  NOTE: This stream sends events only when trades occur
  Event #1: Type=OrderUpdate
    New Orders: 1, Disappeared Orders: 0
    New Positions: 0, Disappeared Positions: 0
    New History Orders: 0, New History Deals: 0
  ✓ Received 3 events, stopping stream
```

---

## 🧠 What It's Used For

The `OnTrade()` method is used for:

* synchronizing local state with the MetaTrader account;
* monitoring broker activity in real-time;
* implementing reactive strategies — automatic reaction to events;
* building trading journals and activity reports.

---

## 💬 In Simple Terms

> `OnTrade()` is a **stream of trading event notifications**.
> It reports: "A new order opened", "A position closed", "A new deal appeared in history".
> The stream activates only when there are changes — no unnecessary data.
