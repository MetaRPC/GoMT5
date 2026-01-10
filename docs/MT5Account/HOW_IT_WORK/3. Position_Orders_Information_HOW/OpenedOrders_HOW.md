### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`OpenedOrders()`** method returns all **active orders and open positions** on the trading account.
> This is one of the key methods for monitoring the current state of the account.
>
> It allows you to get at once:
>
> 1. A list of **open market positions** (PositionInfos);
> 2. A list of **pending orders** awaiting execution (OpenedOrders).


---

## 🧩 Code example

```go
fmt.Println("\n5.2. OpenedOrders() - Get all opened orders & positions")

openedOrdersReq := &pb.OpenedOrdersRequest{
    InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_OPENED_ORDER_SORT_BY_OPEN_TIME,
}
openedOrdersData, err := account.OpenedOrders(ctx, openedOrdersReq)
if err != nil {
    helpers.PrintShortError(err, "OpenedOrders failed")
} else {
    // Direct field access: OpenedOrders and PositionInfos
    totalOrders := len(openedOrdersData.OpenedOrders) + len(openedOrdersData.PositionInfos)
    fmt.Printf("  Total opened orders/positions: %d\n", totalOrders)
    fmt.Printf("    Pending orders:              %d\n", len(openedOrdersData.OpenedOrders))
    fmt.Printf("    Open positions:              %d\n", len(openedOrdersData.PositionInfos))

    // Show first 2 positions if any exist
    maxShow := 2
    if len(openedOrdersData.PositionInfos) < maxShow {
        maxShow = len(openedOrdersData.PositionInfos)
    }
    for i := 0; i < maxShow; i++ {
        pos := openedOrdersData.PositionInfos[i]
        fmt.Printf("\n  Position #%d:\n", i+1)
        fmt.Printf("    Ticket:                      %d\n", pos.Ticket)
        fmt.Printf("    Symbol:                      %s\n", pos.Symbol)
        fmt.Printf("    Type:                        %v\n", pos.Type)
        fmt.Printf("    Volume:                      %.2f\n", pos.Volume)
        fmt.Printf("    Price Open:                  %.5f\n", pos.PriceOpen)
        fmt.Printf("    Current Price:               %.5f\n", pos.PriceCurrent)
        fmt.Printf("    Profit:                      %.2f\n", pos.Profit)
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
openedOrdersReq := &pb.OpenedOrdersRequest{
    InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_OPENED_ORDER_SORT_BY_OPEN_TIME,
}
```

A request is created with a sort mode parameter. The method accepts `InputSortMode` to control how results are sorted:

- `OPENED_ORDER_SORT_BY_OPEN_TIME` (0) - sort by open time
- `OPENED_ORDER_SORT_BY_CLOSE_TIME` (1) - sort by close time
- `OPENED_ORDER_SORT_BY_TICKET_ID` (2) - sort by ticket ID

---

```go
openedOrdersData, err := account.OpenedOrders(ctx, openedOrdersReq)
```

A gRPC request `OpenedOrders()` is sent to the MetaTrader server.
The gateway returns a structure with two lists:

* **`OpenedOrders`** — all pending orders (awaiting execution);
* **`PositionInfos`** — all active open positions.

---

```go
totalOrders := len(openedOrdersData.OpenedOrders) + len(openedOrdersData.PositionInfos)
fmt.Printf("  Total opened orders/positions: %d\n", totalOrders)
fmt.Printf("    Pending orders:              %d\n", len(openedOrdersData.OpenedOrders))
fmt.Printf("    Open positions:              %d\n", len(openedOrdersData.PositionInfos))
```

The total number of active elements is calculated: orders + positions.
Then the number of pending orders and active positions is output separately.

Example output:

```
Total opened orders/positions: 3
  Pending orders:              1
  Open positions:              2
```

---

```go
maxShow := 2
if len(openedOrdersData.PositionInfos) < maxShow {
    maxShow = len(openedOrdersData.PositionInfos)
}
```

Limits the number of displayed positions (for the example - no more than two),
to keep the console readable.

---

```go
for i := 0; i < maxShow; i++ {
    pos := openedOrdersData.PositionInfos[i]
    fmt.Printf("\n  Position #%d:\n", i+1)
    fmt.Printf("    Ticket:                      %d\n", pos.Ticket)
    fmt.Printf("    Symbol:                      %s\n", pos.Symbol)
    fmt.Printf("    Type:                        %v\n", pos.Type)
    fmt.Printf("    Volume:                      %.2f\n", pos.Volume)
    fmt.Printf("    Price Open:                  %.5f\n", pos.PriceOpen)
    fmt.Printf("    Current Price:               %.5f\n", pos.PriceCurrent)
    fmt.Printf("    Profit:                      %.2f\n", pos.Profit)
}
```

Iterates through the list of positions and outputs the main parameters of each:

* **Ticket** — unique position identifier;
* **Symbol** — instrument symbol (e.g., `EURUSD`);
* **Type** — direction (BUY / SELL);
* **Volume** — position volume;
* **PriceOpen** — opening price;
* **PriceCurrent** — current market price;
* **Profit** — current floating profit or loss.

---

## 📦 What the Server Returns

```protobuf
message OpenedOrdersData {
  repeated OpenedOrderInfo opened_orders = 1;     // Pending orders
  repeated PositionInfo position_infos = 2;       // Active positions
}

message OpenedOrderInfo {
  uint32 index = 1;
  uint64 ticket = 2;
  double price_current = 3;
  double price_open = 4;
  double stop_limit = 5;
  double stop_loss = 6;
  double take_profit = 7;
  double volume_current = 8;
  double volume_initial = 9;
  int64 magic_number = 10;
  int32 reason = 11;
  BMT5_ENUM_ORDER_TYPE type = 12;
  BMT5_ENUM_ORDER_STATE state = 13;
  google.protobuf.Timestamp time_expiration = 14;
  google.protobuf.Timestamp time_setup = 15;
  google.protobuf.Timestamp time_done = 16;
  BMT5_ENUM_ORDER_TYPE_FILLING type_filling = 17;
  BMT5_ENUM_ORDER_TYPE_TIME type_time = 18;
  int64 position_id = 19;
  int64 position_by_id = 20;
  string symbol = 21;
  string external_id = 22;
  string comment = 23;
  int64 account_login = 24;
}

message PositionInfo {
  uint32 index = 1;
  uint64 ticket = 2;
  google.protobuf.Timestamp open_time = 3;
  double volume = 4;
  double price_open = 5;
  double stop_loss = 6;
  double take_profit = 7;
  double price_current = 8;
  double swap = 9;
  double profit = 10;
  google.protobuf.Timestamp last_update_time = 11;
  BMT5_ENUM_POSITION_TYPE type = 12;
  int64 magic_number = 13;
  int64 identifier = 14;
  BMT5_ENUM_POSITION_REASON reason = 15;
  string symbol = 16;
  string comment = 17;
  string external_id = 18;
  double position_commission = 19;
  int64 account_login = 20;
}
```

---

## 💡 Example Output

```
Total opened orders/positions: 3
  Pending orders:              1
  Open positions:              2

  Position #1:
    Ticket:                      123456789
    Symbol:                      EURUSD
    Type:                        POSITION_TYPE_BUY
    Volume:                      0.10
    Price Open:                  1.08450
    Current Price:               1.08620
    Profit:                      17.00
```

---

### 🧠 What It's Used For

The `OpenedOrders()` method is used:

* to **monitor all active deals and orders** on the account;
* to **restore the state of positions when starting a trading robot**;
* for **analyzing current profit/loss**;
* when **creating custom interfaces** and monitoring dashboards;
* for **logging and reporting**.

---

### 💬 In Simple Terms

> `OpenedOrders()` shows everything that is currently active on the account:
> open positions and pending orders — with tickets, prices, volumes, and profits.
> This is the main method for monitoring and visualizing trading state.
