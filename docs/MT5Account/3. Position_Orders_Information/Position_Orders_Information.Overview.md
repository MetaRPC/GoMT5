# MT5Account · Positions & Orders Information - Overview

> Open positions, pending orders, historical deals, order history. Use this page to choose the right API for position/order queries.

## 📁 What lives here

### Current State

* **[PositionsTotal](./PositionsTotal.md)** - count of open positions.
* **[OpenedOrders](./OpenedOrders.md)** - full details of all open positions and pending orders.
* **[OpenedOrdersTickets](./OpenedOrdersTickets.md)** - ticket numbers only (lightweight).

### Historical Data

* **[OrderHistory](./OrderHistory.md)** - historical orders within time range (with pagination).
* **[PositionsHistory](./PositionsHistory.md)** - closed positions within time range (with pagination).

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[PositionsTotal - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/PositionsTotal_HOW.md)**
* **[OpenedOrders - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/OpenedOrders_HOW.md)**
* **[OpenedOrdersTickets - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/OpenedOrdersTickets_HOW.md)**
* **[OrderHistory - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/OrderHistory_HOW.md)**
* **[PositionsHistory - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/PositionsHistory_HOW.md)**

---

## 🧭 Plain English

* **PositionsTotal** → quick **count** of open positions.
* **OpenedOrders** → get **full details** of all open positions/orders.
* **OpenedOrdersTickets** → get **ticket numbers only** (faster than full details).
* **OrderHistory** → query **historical orders** (placed, canceled, filled).
* **PositionsHistory** → query **closed positions** (deals).

> Rule of thumb: need **quick count** → `PositionsTotal`; need **full details** → `OpenedOrders`; need **tickets only** → `OpenedOrdersTickets`; need **history** → `OrderHistory` or `PositionsHistory`.

---

## Quick choose

| If you need…                                     | Use                      | Returns                    | Key inputs                               |
| ------------------------------------------------ | ------------------------ | -------------------------- | ---------------------------------------- |
| Count of open positions                          | `PositionsTotal`         | PositionsTotalData (int64) | *(none)*                                 |
| All open positions/orders (full details)         | `OpenedOrders`           | OpenedOrdersData           | Symbol (optional), Group (optional)      |
| Ticket numbers only (lightweight)                | `OpenedOrdersTickets`    | Two arrays of uint64       | Symbol (optional), Group (optional)      |
| Historical orders in time range                  | `OrderHistory`           | OrdersHistoryData          | InputFrom, InputTo, PageNumber, ItemsPerPage |
| Closed positions in time range                   | `PositionsHistory`       | PositionsHistoryData       | PositionOpenTimeFrom, PositionOpenTimeTo, PageNumber, ItemsPerPage |

---

## ❌ Cross‑refs & gotchas

* **Positions vs Orders** - Positions = executed market positions; Orders = pending orders (limits, stops).
* **OpenedOrders returns both** - Returns open positions AND pending orders in separate arrays.
* **Tickets** - Unique identifier for each position/order.
* **Pagination** - OrderHistory and PositionsHistory support PageNumber/ItemsPerPage for large datasets.
* **Time range** - OrderHistory and PositionsHistory use `google.protobuf.Timestamp` (convert from `time.Time` using `timestamppb.New()`).
* **OpenedOrdersTickets** - Much lighter than OpenedOrders if you only need ticket numbers.
* **Historical data limits** - Broker may limit how far back you can query.

---

## 🟢 Minimal snippets

```go
// Count open positions
data, err := account.PositionsTotal(ctx)
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Open positions: %d\n", data.TotalPositions)
```

```go
// Get all open positions/orders
data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Pending orders: %d, Open positions: %d\n",
    len(data.OpenedOrders), len(data.PositionInfos))

for _, pos := range data.PositionInfos {
    fmt.Printf("#%d: %s %.2f lots @ %.5f, P/L: $%.2f\n",
        pos.Ticket,
        pos.Symbol,
        pos.Volume,
        pos.PriceOpen,
        pos.Profit)
}
```

```go
// Get ticket numbers only (lightweight)
data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Pending order tickets: %v\n", data.OpenedOrdersTickets)
fmt.Printf("Position tickets: %v\n", data.OpenedPositionTickets)
```

```go
// Get historical orders (last 7 days)
now := time.Now()
weekAgo := now.Add(-7 * 24 * time.Hour)

data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
    InputFrom:    timestamppb.New(weekAgo),
    InputTo:      timestamppb.New(now),
    PageNumber:   0,
    ItemsPerPage: 100,
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Historical orders: %d\n", len(data.HistoryData))

for _, item := range data.HistoryData {
    if item.HistoryOrder != nil {
        fmt.Printf("#%d: %s @ %.5f\n",
            item.HistoryOrder.Ticket,
            item.HistoryOrder.Symbol,
            item.HistoryOrder.PriceOpen)
    }
}
```

```go
// Get closed positions (today)
now := time.Now()
startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
    PositionOpenTimeFrom: timestamppb.New(startOfDay),
    PositionOpenTimeTo:   timestamppb.New(now),
    PageNumber:           0,
    ItemsPerPage:         100,
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Closed positions today: %d\n", len(data.HistoryPositions))

totalProfit := 0.0
for _, pos := range data.HistoryPositions {
    totalProfit += pos.Profit
    fmt.Printf("#%d: %s %.2f lots → $%.2f\n",
        pos.PositionTicket,
        pos.Symbol,
        pos.Volume,
        pos.Profit)
}

fmt.Printf("Total profit today: $%.2f\n", totalProfit)
```

---

## See also

* **Trading operations:** [OrderSend](../4.%20Trading_Operations/OrderSend.md) - place new orders
* **Order management:** [OrderClose](../4.%20Trading_Operations/OrderClose.md) - close positions
* **Real-time updates:** [OnTrade](../7.%20Streaming_Methods/OnTrade.md) - subscribe to trade events
