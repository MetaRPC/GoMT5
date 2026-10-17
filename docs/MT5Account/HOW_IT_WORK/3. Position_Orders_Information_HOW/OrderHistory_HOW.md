### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`OrderHistory()`** method is used to get **historical orders** (executed, cancelled, expired) for a selected time period.
> It does not return open positions — only completed operations.


---

## 🧩 Code example

```go
fmt.Println("\n5.4. OrderHistory() - Get historical orders")

// Set time range (last 30 days) - must use timestamppb.Timestamp
now := time.Now()
fromTimestamp := timestamppb.New(now.AddDate(0, 0, -30))
toTimestamp := timestamppb.New(now)

orderHistoryReq := &pb.OrderHistoryRequest{
    InputFromUtc:  fromTimestamp, // time.Now() is UTC-based: use the *Utc fields
    InputToUtc:    toTimestamp,
    InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC,
    PageNumber:    0,  // 0 = page-number paging not used (all items)
    ItemsPerPage:  0,  // 0 = page-number paging not used (all items)
}
orderHistoryData, err := account.OrderHistory(ctx, orderHistoryReq)
if err != nil {
    helpers.PrintShortError(err, "OrderHistory failed")
} else {
    fmt.Printf("  Historical orders (last 30d):  %d\n", len(orderHistoryData.HistoryData))

    maxShow := 2
    if len(orderHistoryData.HistoryData) < maxShow {
        maxShow = len(orderHistoryData.HistoryData)
    }
    for i := 0; i < maxShow; i++ {
        item := orderHistoryData.HistoryData[i]
        if item.HistoryOrder != nil {
            order := item.HistoryOrder
            fmt.Printf("\n  Order #%d:\n", i+1)
            fmt.Printf("    Ticket:                      %d\n", order.Ticket)
            fmt.Printf("    Symbol:                      %s\n", order.Symbol)
            fmt.Printf("    Type:                        %v\n", order.Type)
            fmt.Printf("    Volume:                      %.2f\n", order.VolumeInitial)
            fmt.Printf("    Price Open:                  %.5f\n", order.PriceOpen)
        }
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
now := time.Now()
fromTimestamp := timestamppb.New(now.AddDate(0, 0, -30))
toTimestamp := timestamppb.New(now)
```

A date range is formed — the last 30 days.
MetaTrader requires fields of type `timestamppb.Timestamp`, so conversion from `time.Time` is used.

---

```go
orderHistoryReq := &pb.OrderHistoryRequest{
    InputFromUtc:  fromTimestamp,
    InputToUtc:    toTimestamp,
    InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC,
    PageNumber:    0,
    ItemsPerPage:  0,
}
```

A request is created to get orders for the period between `InputFromUtc` and `InputToUtc`.

**Request Parameters:**

- `InputFrom` / `InputFromUtc` - start time of history range (required): `InputFrom` is in the broker server's time zone, `InputFromUtc` in UTC; send one of them, not both
- `InputTo` / `InputToUtc` - end time of history range (required), same rule
- `InputSortMode` - sort mode (optional):
  - `BMT5_SORT_BY_OPEN_TIME_ASC` (0) / `BMT5_SORT_BY_OPEN_TIME_DESC` (1) - sort by open time
  - `BMT5_SORT_BY_CLOSE_TIME_ASC` (2) / `BMT5_SORT_BY_CLOSE_TIME_DESC` (3) - sort by close time
  - `BMT5_SORT_BY_ORDER_TICKET_ID_ASC` (4) / `BMT5_SORT_BY_ORDER_TICKET_ID_DESC` (5) - sort by ticket ID
- `PageNumber` - page number for legacy page-number paging (1-based, needs `ItemsPerPage`; 0 = not used, all items)
- `ItemsPerPage` - rows per page for legacy page-number paging (0 = not used)
- `PageToken` / `PageSize` - cursor paging (preferred); cannot be combined with `PageNumber`/`ItemsPerPage`

---

```go
orderHistoryData, err := account.OrderHistory(ctx, orderHistoryReq)
```

A gRPC call `OrderHistory()` is sent.
In response, an array of `HistoryData` structures is received, each describing one completed order.

---

```go
fmt.Printf("  Historical orders (last 30d):  %d\n", len(orderHistoryData.HistoryData))
```

The total number of historical orders found for the specified period is displayed.

---

```go
if item.HistoryOrder != nil {
    order := item.HistoryOrder
    fmt.Printf("\n  Order #%d:\n", i+1)
    fmt.Printf("    Ticket:                      %d\n", order.Ticket)
    fmt.Printf("    Symbol:                      %s\n", order.Symbol)
    fmt.Printf("    Type:                        %v\n", order.Type)
    fmt.Printf("    Volume:                      %.2f\n", order.VolumeInitial)
    fmt.Printf("    Price Open:                  %.5f\n", order.PriceOpen)
}
```

The loop outputs the main fields for each order:

* **Ticket** — unique order number;
* **Symbol** — instrument for which it was placed;
* **Type** — order type (`BUY`, `SELL`, `BUY_LIMIT`, etc.);
* **VolumeInitial** — initial volume;
* **PriceOpen** — order opening price.

---

## 📦 What the Server Returns

```protobuf
message OrdersHistoryData {
  int32 arrayTotal = 1;
  int32 pageNumber = 2;
  int32 itemsPerPage = 3;
  repeated HistoryData history_data = 4;
  string next_page_token = 5;  // cursor paging: token of the next page
}

message HistoryData {
  uint32 index = 1;
  OrderHistoryData history_order = 2;
  DealHistoryData history_deal = 3;
}

message OrderHistoryData {
  uint64 ticket = 1;
  google.protobuf.Timestamp setup_time = 2;
  google.protobuf.Timestamp done_time = 3;
  BMT5_ENUM_ORDER_STATE state = 4;
  double price_current = 5;
  double price_open = 6;
  double stop_limit = 7;
  double stop_loss = 8;
  double take_profit = 9;
  double volume_current = 10;
  double volume_initial = 11;
  int64 magic_number = 12;
  BMT5_ENUM_ORDER_TYPE type = 13;
  google.protobuf.Timestamp time_expiration = 14;
  BMT5_ENUM_ORDER_TYPE_FILLING type_filling = 15;
  BMT5_ENUM_ORDER_TYPE_TIME type_time = 16;
  uint64 position_id = 17;
  string symbol = 18;
  string external_id = 19;
  string comment = 20;
  int64 account_login = 21;
}
```

---

## 💡 Example Output

```
Historical orders (last 30d):  2

  Order #1:
    Ticket:                      128007201
    Symbol:                      EURUSD
    Type:                        BMT5_ORDER_TYPE_BUY
    Volume:                      0.01
    Price Open:                  0.00000

  Order #2:
    Ticket:                      128009723
    Symbol:                      EURUSD
    Type:                        BMT5_ORDER_TYPE_BUY
    Volume:                      0.01
    Price Open:                  0.00000
```

---

## ⚠️ Why Price Open Equals 0.00000

If you see that the **`PriceOpen` = 0.00000** field, this is not a code error. Possible reasons:

1. **Historical orders were not executed (Pending Orders)**
   If an order was cancelled by the user or expired by time, the server stores it in history but without an opening price.

2. **Market orders are filled at the market price**
   A market order may be stored with `PriceOpen` = 0; the actual execution price is in the matching deal (`HistoryDeal.Price`).

3. **Historical data limited by server**
   Some brokers store only basic order attributes, and if an order was closed long ago or was not executed, some fields may be zeroed.

4. **MetaTrader demo server**
   On demo servers (especially test clusters), trade history is sometimes incomplete — `PriceOpen` may not be filled.

---

## 🧠 What It's Used For

The `OrderHistory()` method is used:

* for **analyzing trading** for a specific period (PnL, trade frequency, instruments);
* for **restoring history** after restarting a trading robot;
* for **forming reports** and trade journals;
* in **analytical modules** that build charts or statistics on trades.

---

### 💬 In Simple Terms

> `OrderHistory()` allows you to get **history of completed orders** for a specified period.
> If the opening price is zero, it means the order **was not executed** or **data is incomplete** on the server side.
> This is a normal situation when working with demo or inactive orders.
