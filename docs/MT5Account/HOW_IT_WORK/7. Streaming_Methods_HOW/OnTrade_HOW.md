### Example from file: `examples/demos/lowlevel/03_streaming_methods.go`

> The **`OnTrade()`** method is used to receive **real-time trading events**.
> It notifies about all trading account changes: order opening and closing, new positions appearing, deals in history, etc.
>
> Unlike `OnSymbolTick()`, which streams quotes, `OnTrade()` transmits **trading events** occurring on the account.


---

## 🧩 Code example

```go
tradeReq := &pb.OnTradeRequest{}

tradeChan, tradeErrChan := account.OnTrade(ctx, tradeReq)

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
tradeChan, tradeErrChan := account.OnTrade(ctx, tradeReq)
```

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

Display the type of trading event.
Examples of possible types:

* `TRADE_EVENT_ORDER_ADDED`
* `TRADE_EVENT_ORDER_DELETED`
* `TRADE_EVENT_POSITION_ADDED`
* `TRADE_EVENT_POSITION_CLOSED`

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

> Thus, `OnTrade()` returns a **delta of changes** — allowing you to synchronize local data with the server.

---

### 7️. Stream Termination

The stream closes when the event limit (`MAX_EVENTS`) is reached, timeout expires, or a connection error occurs.

---

## 📦 What the Server Returns

```protobuf
message OnTradeData {
  ENUM_TRADE_EVENT_TYPE Type = 1; // Event type
  EventData EventData = 2;         // Change details
}

message EventData {
  repeated OrderInfo NewOrders = 1;
  repeated OrderInfo DisappearedOrders = 2;
  repeated PositionInfo NewPositions = 3;
  repeated PositionInfo DisappearedPositions = 4;
  repeated HistoryOrderInfo NewHistoryOrders = 5;
  repeated HistoryDealInfo NewHistoryDeals = 6;
}
```

---

## 💡 Example Output

```
Streaming trade events (max 3 events or 10 seconds)...
  NOTE: This stream sends events only when trades occur
  Event #1: Type=TRADE_EVENT_ORDER_ADDED
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
