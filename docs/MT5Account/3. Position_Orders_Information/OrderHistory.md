# ✅ Get Order History

> **Request:** retrieve historical orders within a specified time range with pagination support for large datasets.

**API Information:**

* **Low-level API:** `MT5Account.OrderHistory(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `OrderHistory` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `OrderHistory(OrderHistoryRequest) → OrderHistoryReply`
* **Low‑level client (generated):** `AccountHelperClient.OrderHistory(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieves historical orders within a specified time range.
* **Why you need it.** Analyze trading history, calculate statistics, review past activity.
* **Pagination support.** Use PageNumber and ItemsPerPage for handling large datasets efficiently.

---

## 🎯 Purpose

Use it to:

* Retrieve historical order records
* Analyze trading performance
* Generate trading reports
* Audit order execution history
* Track order modifications

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OrderHistory - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/OrderHistory_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// OrderHistory retrieves historical orders within a specified time range.
// Supports pagination via Offset and Limit parameters.
func (a *MT5Account) OrderHistory(
    ctx context.Context,
    req *pb.OrderHistoryRequest,
) (*pb.OrdersHistoryData, error)
```

**Request message:**

```protobuf
OrderHistoryRequest {
  google.protobuf.Timestamp inputFrom = 1;                   // Start date (required - server time)
  google.protobuf.Timestamp inputTo = 2;                     // End date (required - server time)
  BMT5_ENUM_ORDER_HISTORY_SORT_TYPE inputSortMode = 3;     // Sort mode (0 - by open time, 1 - by close time, 2 - by ticket id)
  int32 pageNumber = 4;                                      // Page number (default 0)
  int32 itemsPerPage = 5;                                    // Items per page (default 0 = all)
}
```

**Reply message:**

```protobuf
OrderHistoryReply {
  oneof response {
    OrdersHistoryData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                         | Description                                   |
| --------- | ---------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`            | Context for deadline/timeout and cancellation |
| `req`     | `*pb.OrderHistoryRequest`    | Request with date range and pagination        |

**Request fields:**

| Field          | Type                                   | Description                                                  |
| -------------- | -------------------------------------- | ------------------------------------------------------------ |
| `InputFrom`    | `google.protobuf.Timestamp`            | Start date (required - server time)                          |
| `InputTo`      | `google.protobuf.Timestamp`            | End date (required - server time)                            |
| `InputSortMode`| `BMT5_ENUM_ORDER_HISTORY_SORT_TYPE`    | Sort mode (optional: 0=open time asc, 1=open time desc, etc)|
| `PageNumber`   | `int32`                                | Page number for pagination (default 0)                       |
| `ItemsPerPage` | `int32`                                | Number of items per page (default 0 = all)                   |

---

## ⬆️ Output — `OrdersHistoryData`

| Field          | Type             | Go Type             | Description                                      |
| -------------- | ---------------- | ------------------- | ------------------------------------------------ |
| `ArrayTotal`   | `int32`          | `int32`             | Total number of orders in query result           |
| `PageNumber`   | `int32`          | `int32`             | Current page number                              |
| `ItemsPerPage` | `int32`          | `int32`             | Number of items per page                         |
| `HistoryData`  | `HistoryData[]`  | `[]*pb.HistoryData` | Array of historical data (orders and deals)      |

**HistoryData structure:**

- `Index` (uint32) - Record index
- `HistoryOrder` (*OrderHistoryData) - Order information (if available)
- `HistoryDeal` (*DealHistoryData) - Deal information (if available)

**OrderHistoryData complete structure:**

- `Ticket` (uint64) - Order ticket number
- `SetupTime` (*timestamppb.Timestamp) - Time when order was placed
- `DoneTime` (*timestamppb.Timestamp) - Time when order was executed/cancelled
- `State` (BMT5_ENUM_ORDER_STATE) - Order state
- `PriceCurrent` (float64) - Current market price
- `PriceOpen` (float64) - Order open price
- `StopLimit` (float64) - Stop limit price
- `StopLoss` (float64) - Stop loss level
- `TakeProfit` (float64) - Take profit level
- `VolumeCurrent` (float64) - Current volume (for partially filled orders)
- `VolumeInitial` (float64) - Initial volume
- `MagicNumber` (int64) - Magic number
- `Type` (BMT5_ENUM_ORDER_TYPE) - Order type
- `TimeExpiration` (*timestamppb.Timestamp) - Expiration time
- `TypeFilling` (BMT5_ENUM_ORDER_TYPE_FILLING) - Order filling type
- `TypeTime` (BMT5_ENUM_ORDER_TYPE_TIME) - Order time type
- `PositionId` (uint64) - Position ID
- `Symbol` (string) - Trading symbol
- `ExternalId` (string) - External identifier
- `Comment` (string) - Order comment
- `AccountLogin` (int64) - Account login

**DealHistoryData complete structure:**

- `Ticket` (uint64) - Deal ticket number
- `Profit` (float64) - Deal profit
- `Commission` (float64) - Commission
- `Fee` (float64) - Additional fee
- `Price` (float64) - Execution price
- `StopLoss` (float64) - Stop loss level
- `TakeProfit` (float64) - Take profit level
- `Swap` (float64) - Swap charges
- `Volume` (float64) - Deal volume
- `EntryType` (BMT5_ENUM_DEAL_ENTRY_TYPE) - Deal entry type
- `Time` (*timestamppb.Timestamp) - Deal execution time
- `Type` (BMT5_ENUM_DEAL_TYPE) - Deal type
- `Reason` (BMT5_ENUM_DEAL_REASON) - Deal execution reason
- `PositionId` (uint64) - Position ID
- `Comment` (string) - Deal comment
- `Symbol` (string) - Trading symbol
- `ExternalId` (string) - External identifier
- `AccountLogin` (int64) - Account login

---

> **💡 Enum Usage Note:** The tables show simplified constant names for readability.
> In Go code, use full names with the enum type prefix.
>
> Format: `pb.<ENUM_TYPE>_<CONSTANT_NAME>`
>
> Example: `pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC`

### 📘 Enum: BMT5_ENUM_ORDER_HISTORY_SORT_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_SORT_BY_OPEN_TIME_ASC` | Sort by open time (ascending) |
| 1 | `BMT5_SORT_BY_OPEN_TIME_DESC` | Sort by open time (descending) |
| 2 | `BMT5_SORT_BY_CLOSE_TIME_ASC` | Sort by close time (ascending) |
| 3 | `BMT5_SORT_BY_CLOSE_TIME_DESC` | Sort by close time (descending) |
| 4 | `BMT5_SORT_BY_ORDER_TICKET_ID_ASC` | Sort by order ticket ID (ascending) |
| 5 | `BMT5_SORT_BY_ORDER_TICKET_ID_DESC` | Sort by order ticket ID (descending) |

### 📘 Enum: BMT5_ENUM_ORDER_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_ORDER_TYPE_BUY` | Market Buy order |
| 1 | `BMT5_ORDER_TYPE_SELL` | Market Sell order |
| 2 | `BMT5_ORDER_TYPE_BUY_LIMIT` | Buy Limit pending order |
| 3 | `BMT5_ORDER_TYPE_SELL_LIMIT` | Sell Limit pending order |
| 4 | `BMT5_ORDER_TYPE_BUY_STOP` | Buy Stop pending order |
| 5 | `BMT5_ORDER_TYPE_SELL_STOP` | Sell Stop pending order |
| 6 | `BMT5_ORDER_TYPE_BUY_STOP_LIMIT` | Buy Stop Limit (pending Buy Limit order at StopLimit price) |
| 7 | `BMT5_ORDER_TYPE_SELL_STOP_LIMIT` | Sell Stop Limit (pending Sell Limit order at StopLimit price) |
| 8 | `BMT5_ORDER_TYPE_CLOSE_BY` | Order to close a position by an opposite one |

### 📘 Enum: BMT5_ENUM_ORDER_STATE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_ORDER_STATE_STARTED` | Order checked, but not yet accepted by broker |
| 1 | `BMT5_ORDER_STATE_PLACED` | Order accepted |
| 2 | `BMT5_ORDER_STATE_CANCELED` | Order canceled by client |
| 3 | `BMT5_ORDER_STATE_PARTIAL` | Order partially executed |
| 4 | `BMT5_ORDER_STATE_FILLED` | Order fully executed |
| 5 | `BMT5_ORDER_STATE_REJECTED` | Order rejected |
| 6 | `BMT5_ORDER_STATE_EXPIRED` | Order expired |
| 7 | `BMT5_ORDER_STATE_REQUEST_ADD` | Order is being registered (placing to trading system) |
| 8 | `BMT5_ORDER_STATE_REQUEST_MODIFY` | Order is being modified (changing parameters) |
| 9 | `BMT5_ORDER_STATE_REQUEST_CANCEL` | Order is being deleted (deleting from trading system) |

**Used in:** OrderHistoryData.State

### 📘 Enum: BMT5_ENUM_ORDER_TYPE_FILLING

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_ORDER_FILLING_FOK` | Fill or Kill - order must be filled completely or not at all |
| 1 | `BMT5_ORDER_FILLING_IOC` | Immediate or Cancel - fill available volume, cancel the rest |
| 2 | `BMT5_ORDER_FILLING_RETURN` | Return - order is placed with the broker for execution |
| 3 | `BMT5_ORDER_FILLING_BOC` | Book or Cancel - order must be placed as a passive order (limit) |

**Used in:** OrderHistoryData.TypeFilling

### 📘 Enum: BMT5_ENUM_ORDER_TYPE_TIME

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_ORDER_TIME_GTC` | Good Till Cancelled - order stays until explicitly cancelled |
| 1 | `BMT5_ORDER_TIME_DAY` | Good Till Day - order valid until end of trading day |
| 2 | `BMT5_ORDER_TIME_SPECIFIED` | Good Till Specified - order valid until specified date/time |
| 3 | `BMT5_ORDER_TIME_SPECIFIED_DAY` | Good Till Specified Day - order valid until end of specified day |

**Used in:** OrderHistoryData.TypeTime

### 📘 Enum: BMT5_ENUM_DEAL_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_DEAL_TYPE_BUY` | Buy deal |
| 1 | `BMT5_DEAL_TYPE_SELL` | Sell deal |
| 2 | `BMT5_DEAL_TYPE_BALANCE` | Balance operation |
| 3 | `BMT5_DEAL_TYPE_CREDIT` | Credit operation |
| 4 | `BMT5_DEAL_TYPE_CHARGE` | Additional charge |
| 5 | `BMT5_DEAL_TYPE_CORRECTION` | Correction |
| 6 | `BMT5_DEAL_TYPE_BONUS` | Bonus |
| 7 | `BMT5_DEAL_TYPE_COMMISSION` | Additional commission |
| 8 | `BMT5_DEAL_TYPE_COMMISSION_DAILY` | Daily commission |
| 9 | `BMT5_DEAL_TYPE_COMMISSION_MONTHLY` | Monthly commission |
| 10 | `BMT5_DEAL_TYPE_COMMISSION_AGENT_DAILY` | Daily agent commission |
| 11 | `BMT5_DEAL_TYPE_COMMISSION_AGENT_MONTHLY` | Monthly agent commission |
| 12 | `BMT5_DEAL_TYPE_INTEREST` | Interest rate |
| 13 | `BMT5_DEAL_TYPE_BUY_CANCELED` | Canceled buy deal |
| 14 | `BMT5_DEAL_TYPE_SELL_CANCELED` | Canceled sell deal |
| 15 | `BMT5_DEAL_DIVIDEND` | Dividend operations |
| 16 | `BMT5_DEAL_DIVIDEND_FRANKED` | Franked (non-taxable) dividend operations |
| 17 | `BMT5_DEAL_TAX` | Tax charges |

**Used in:** DealHistoryData.Type

### 📘 Enum: BMT5_ENUM_DEAL_REASON

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_DEAL_REASON_CLIENT` | Deal executed from desktop terminal |
| 1 | `BMT5_DEAL_REASON_MOBILE` | Deal executed from mobile application |
| 2 | `BMT5_DEAL_REASON_WEB` | Deal executed from web platform |
| 3 | `BMT5_DEAL_REASON_EXPERT` | Deal executed by Expert Advisor or script |
| 4 | `BMT5_DEAL_REASON_SL` | Deal executed by Stop Loss |
| 5 | `BMT5_DEAL_REASON_TP` | Deal executed by Take Profit |
| 6 | `BMT5_DEAL_REASON_SO` | Deal executed by Stop Out |
| 7 | `BMT5_DEAL_REASON_ROLLOVER` | Deal executed due to rollover |
| 8 | `BMT5_DEAL_REASON_VMARGIN` | Deal executed after variation margin |
| 9 | `BMT5_DEAL_REASON_SPLIT` | Deal executed after split |
| 10 | `BMT5_DEAL_REASON_CORPORATE_ACTION` | Deal executed due to corporate action |

**Used in:** DealHistoryData.Reason

### 📘 Enum: BMT5_ENUM_DEAL_ENTRY_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_DEAL_ENTRY_IN` | Entry into the market |
| 1 | `BMT5_DEAL_ENTRY_OUT` | Exit from the market |
| 2 | `BMT5_DEAL_ENTRY_INOUT` | Reverse (close and open opposite position) |
| 3 | `BMT5_DEAL_ENTRY_OUT_BY` | Close position by an opposite one |

**Used in:** DealHistoryData.EntryType

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `15s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Time range:** InputFrom and InputTo are `google.protobuf.Timestamp` (use `timestamppb.New()` to convert from `time.Time`).
* **Pagination:** Use PageNumber and ItemsPerPage to retrieve large datasets in chunks.
* **Sort mode:** InputSortMode allows sorting by open time, close time, or ticket ID.

---

## 🔗 Usage Examples

### 1) Get last 7 days history

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    now := time.Now()
    weekAgo := now.Add(-7 * 24 * time.Hour)

    data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
        InputFrom:     timestamppb.New(weekAgo),
        InputTo:       timestamppb.New(now),
        InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_DESC,
        PageNumber:    0,
        ItemsPerPage:  100,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Orders in last 7 days: %d\n", len(data.HistoryData))
}
```

### 2) Get history with sorting

```go
func GetHistoryWithSorting(account *mt5.MT5Account, days int) ([]*pb.HistoryData, error) {
    ctx := context.Background()

    now := time.Now()
    fromDate := now.Add(-time.Duration(days) * 24 * time.Hour)

    data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
        InputFrom:     timestamppb.New(fromDate),
        InputTo:       timestamppb.New(now),
        InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_CLOSE_TIME_DESC,
        PageNumber:    0,
        ItemsPerPage:  0, // All orders
    })
    if err != nil {
        return nil, err
    }

    return data.HistoryData, nil
}

// Usage:
// orders, _ := GetHistoryWithSorting(account, 30)
// fmt.Printf("Orders last 30 days (sorted by close time): %d\n", len(orders))
```

### 3) Paginated history retrieval

```go
func GetAllHistoryPaginated(account *mt5.MT5Account, fromDate, toDate time.Time, pageSize int32) ([]*pb.HistoryData, error) {
    ctx := context.Background()

    allHistory := []*pb.HistoryData{}
    pageNumber := int32(0)

    for {
        data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
            InputFrom:     timestamppb.New(fromDate),
            InputTo:       timestamppb.New(toDate),
            InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC,
            PageNumber:    pageNumber,
            ItemsPerPage:  pageSize,
        })
        if err != nil {
            return nil, err
        }

        if len(data.HistoryData) == 0 {
            break
        }

        allHistory = append(allHistory, data.HistoryData...)
        pageNumber++

        fmt.Printf("Retrieved %d orders (total: %d)\n", len(data.HistoryData), len(allHistory))

        if int32(len(data.HistoryData)) < pageSize {
            break
        }
    }

    return allHistory, nil
}

// Usage:
// fromDate := time.Now().Add(-30 * 24 * time.Hour)
// toDate := time.Now()
// orders, _ := GetAllHistoryPaginated(account, fromDate, toDate, 100)
```

### 4) Calculate monthly statistics

```go
func GetMonthlyStats(account *mt5.MT5Account, year, month int) {
    ctx := context.Background()

    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

    data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
        InputFrom:     timestamppb.New(startDate),
        InputTo:       timestamppb.New(endDate),
        InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC,
        PageNumber:    0,
        ItemsPerPage:  0,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    var totalProfit float64
    wins := 0
    losses := 0

    for _, item := range data.HistoryData {
        if item.HistoryOrder != nil {
            // Note: Profit may not be available in OrderHistoryData
            // Use DealHistoryData for profit information
        }
        if item.HistoryDeal != nil {
            totalProfit += item.HistoryDeal.Profit
            if item.HistoryDeal.Profit > 0 {
                wins++
            } else if item.HistoryDeal.Profit < 0 {
                losses++
            }
        }
    }

    fmt.Printf("Month %d/%d Statistics:\n", month, year)
    fmt.Printf("  Total Orders: %d\n", len(data.HistoryData))
    fmt.Printf("  Wins: %d\n", wins)
    fmt.Printf("  Losses: %d\n", losses)
    if wins+losses > 0 {
        fmt.Printf("  Win Rate: %.2f%%\n", float64(wins)/float64(wins+losses)*100)
    }
    fmt.Printf("  Total Profit: %.2f\n", totalProfit)
}
```

### 5) Find orders by date range

```go
func FindOrdersByDateRange(account *mt5.MT5Account, startDate, endDate string) ([]*pb.HistoryData, error) {
    ctx := context.Background()

    start, _ := time.Parse("2006-01-02", startDate)
    end, _ := time.Parse("2006-01-02", endDate)
    end = end.Add(24 * time.Hour).Add(-time.Second)

    data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
        InputFrom:     timestamppb.New(start),
        InputTo:       timestamppb.New(end),
        InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC,
        PageNumber:    0,
        ItemsPerPage:  0,
    })
    if err != nil {
        return nil, err
    }

    return data.HistoryData, nil
}

// Usage:
// orders, _ := FindOrdersByDateRange(account, "2024-01-01", "2024-01-31")
```

---

## 🔧 Common Patterns

### Get today's history

```go
func GetTodayHistory(account *mt5.MT5Account) ([]*pb.HistoryData, error) {
    ctx := context.Background()

    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

    data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
        InputFrom:     timestamppb.New(startOfDay),
        InputTo:       timestamppb.New(now),
        InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC,
        PageNumber:    0,
        ItemsPerPage:  0,
    })
    if err != nil {
        return nil, err
    }

    return data.HistoryData, nil
}
```

### Calculate profit for period

```go
func CalculatePeriodProfit(account *mt5.MT5Account, fromDate, toDate time.Time) (float64, error) {
    ctx := context.Background()

    data, err := account.OrderHistory(ctx, &pb.OrderHistoryRequest{
        InputFrom:     timestamppb.New(fromDate),
        InputTo:       timestamppb.New(toDate),
        InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC,
        PageNumber:    0,
        ItemsPerPage:  0,
    })
    if err != nil {
        return 0, err
    }

    var totalProfit float64
    for _, item := range data.HistoryData {
        // Profit is in DealHistoryData, not OrderHistoryData
        if item.HistoryDeal != nil {
            totalProfit += item.HistoryDeal.Profit
        }
    }

    return totalProfit, nil
}
```

---

## 📚 See Also

* [PositionsHistory](./PositionsHistory.md) - Get closed positions with P&L
* [OpenedOrders](./OpenedOrders.md) - Get current open positions
* [OnTrade](../7.%20Streaming_Methods/OnTrade.md) - Stream trade events in real-time
