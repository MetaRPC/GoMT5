### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`PositionsHistory()`** method returns the history of **closed positions** for a specified time range.
> This is one of the main ways to get information about profit and loss (P&L) on the account for a period.


---

## 🧩 Code example

```go
fmt.Println("\n5.5. PositionsHistory() - Get historical positions with P&L")

positionsHistoryReq := &pb.PositionsHistoryRequest{
    SortType:                pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
    PositionOpenTimeFromUtc: fromTimestamp, // time.Now() is UTC-based: use the *Utc fields
    PositionOpenTimeToUtc:   toTimestamp,
    PageNumber:              nil,  // nil for all items
    ItemsPerPage:            nil,  // nil for all items
}
positionsHistoryData, err := account.PositionsHistory(ctx, positionsHistoryReq)
if err != nil {
    helpers.PrintShortError(err, "PositionsHistory failed")
} else {
    fmt.Printf("  Historical positions (last 30d): %d\n", len(positionsHistoryData.HistoryPositions))

    maxShow := 2
    if len(positionsHistoryData.HistoryPositions) < maxShow {
        maxShow = len(positionsHistoryData.HistoryPositions)
    }
    for i := 0; i < maxShow; i++ {
        pos := positionsHistoryData.HistoryPositions[i]
        fmt.Printf("\n  Position #%d:\n", i+1)
        fmt.Printf("    Position Ticket:             %d\n", pos.PositionTicket)
        fmt.Printf("    Symbol:                      %s\n", pos.Symbol)
        fmt.Printf("    Order Type:                  %v\n", pos.OrderType)
        fmt.Printf("    Volume:                      %.2f\n", pos.Volume)
        fmt.Printf("    Profit:                      %.2f\n", pos.Profit)
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
positionsHistoryReq := &pb.PositionsHistoryRequest{
    SortType:                pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
    PositionOpenTimeFromUtc: fromTimestamp,
    PositionOpenTimeToUtc:   toTimestamp,
    PageNumber:              nil,
    ItemsPerPage:            nil,
}
```

A request is created with a date range specification for which position history is needed.
Time format is `timestamppb.Timestamp`, as in the `OrderHistory()` method.

**Request Parameters:**

- `SortType` - sort type:
  - `AH_POSITION_OPEN_TIME_ASC` (0) - sort by open time, ascending
  - `AH_POSITION_OPEN_TIME_DESC` (1) - sort by open time, descending
  - `AH_POSITION_TICKET_ASC` (2) - sort by ticket, ascending
  - `AH_POSITION_TICKET_DESC` (3) - sort by ticket, descending
- `PositionOpenTimeFrom` / `PositionOpenTimeFromUtc` - start of the position open time filter (optional); the first is in the broker server's time zone, the second in UTC: send one of them, not both
- `PositionOpenTimeTo` / `PositionOpenTimeToUtc` - end of the position open time filter (optional), same rule
- `PageNumber` - page number for pagination (optional, nil for all items)
- `ItemsPerPage` - number of items per page (optional, nil for all items)

---

```go
positionsHistoryData, err := account.PositionsHistory(ctx, positionsHistoryReq)
```

A gRPC request `PositionsHistory()` is sent to the MetaTrader server.
In response, a list of all positions that were opened and closed in the specified period is returned.

---

```go
fmt.Printf("  Historical positions (last 30d): %d\n", len(positionsHistoryData.HistoryPositions))
```

Outputs the total number of closed positions found for the last 30 days.

---

```go
for i := 0; i < maxShow; i++ {
    pos := positionsHistoryData.HistoryPositions[i]
    fmt.Printf("\n  Position #%d:\n", i+1)
    fmt.Printf("    Position Ticket:             %d\n", pos.PositionTicket)
    fmt.Printf("    Symbol:                      %s\n", pos.Symbol)
    fmt.Printf("    Order Type:                  %v\n", pos.OrderType)
    fmt.Printf("    Volume:                      %.2f\n", pos.Volume)
    fmt.Printf("    Profit:                      %.2f\n", pos.Profit)
}
```

Outputs brief information for each position:

* **PositionTicket** — unique position number;
* **Symbol** — instrument (e.g., `EURUSD`);
* **OrderType** — trade direction (`BUY`, `SELL`, etc.);
* **Volume** — position volume;
* **Profit** — final profit or loss.

---

## 📦 What the Server Returns

```protobuf
message PositionsHistoryData {
  repeated PositionHistoryInfo history_positions = 1;
}

message PositionHistoryInfo {
  int32 index = 1;
  uint64 position_ticket = 2;
  AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE order_type = 3;
  google.protobuf.Timestamp open_time = 4;
  google.protobuf.Timestamp close_time = 5;
  double volume = 6;
  double open_price = 7;
  double close_price = 8;
  double stop_loss = 9;
  double take_profit = 10;
  double market_value = 11;
  double commission = 12;
  double fee = 13;
  double profit = 14;
  double swap = 15;
  string comment = 16;
  string symbol = 17;
  int64 magic = 18;
}
```

>  In addition to the fields shown in the example, the full response also includes open/close prices, commissions, swaps, and trade times.

---

## 💡 Example Output

```
Historical positions (last 30d): 2

  Position #1:
    Position Ticket:             127998001
    Symbol:                      EURUSD
    Order Type:                  ORDER_TYPE_BUY
    Volume:                      0.10
    Profit:                      25.35

  Position #2:
    Position Ticket:             128003420
    Symbol:                      GBPUSD
    Order Type:                  ORDER_TYPE_SELL
    Volume:                      0.05
    Profit:                      -7.60
```

---

## ⚠️ Possible Reasons for Empty Output

If the method returns zero positions, this may be related to:

* absence of closed trades in the last 30 days;
* broker policy that stores history for a limited time;
* use of a demo server (history is often cleared);
* too short date range in the request.

---

### 🧠 What It's Used For

The `PositionsHistory()` method is used:

* for **analyzing trading results** (PnL, win rate, drawdown);
* when **forming reports and analytical dashboards**;
* for **backtesting strategies** to load historical data;
* in **robots** that restore statistics after restart.

---

### 💬 In Simple Terms

> `PositionsHistory()` returns **all closed positions** for a specified period.
> Shows symbol, direction, volume, and profit/loss.
> This is the main way to get P&L history for an account.
