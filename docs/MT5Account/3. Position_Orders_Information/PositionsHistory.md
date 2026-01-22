# ✅ Get Closed Positions History

> **Request:** retrieve history of closed positions within a specified time range with profit/loss, swap, and commission information. Supports pagination.

**API Information:**

* **Low-level API:** `MT5Account.PositionsHistory(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `PositionsHistory` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `PositionsHistory(PositionsHistoryRequest) → PositionsHistoryReply`
* **Low‑level client (generated):** `AccountHelperClient.PositionsHistory(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieves closed positions within a time range with full P&L breakdown.
* **Why you need it.** Analyze trading performance, calculate statistics, review profitability.
* **Complete P&L.** Includes net profit, swap, commission, entry/exit prices.

---

## 🎯 Purpose

Use it to:

* Retrieve closed positions from history
* Analyze trading performance with P&L
* Generate performance reports
* Calculate win rate and statistics
* Track position modifications and outcomes

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [PositionsHistory - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/PositionsHistory_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// PositionsHistory retrieves closed positions with profit/loss information.
// Supports pagination and symbol filtering.
func (a *MT5Account) PositionsHistory(
    ctx context.Context,
    req *pb.PositionsHistoryRequest,
) (*pb.PositionsHistoryData, error)
```

**Request message:**

```protobuf
PositionsHistoryRequest {
  AH_ENUM_POSITIONS_HISTORY_SORT_TYPE sort_type = 1;           // Sort mode (optional)
  google.protobuf.Timestamp position_open_time_from = 2;       // Start date (optional)
  google.protobuf.Timestamp position_open_time_to = 3;         // End date (optional)
  int32 page_number = 4;                                        // Page number (optional)
  int32 items_per_page = 5;                                     // Items per page (optional)
}
```

**Reply message:**

```protobuf
PositionsHistoryReply {
  oneof response {
    PositionsHistoryData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                             | Description                                   |
| --------- | -------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                | Context for deadline/timeout and cancellation |
| `req`     | `*pb.PositionsHistoryRequest`    | Request with date range and pagination        |

**Request fields:**

| Field                   | Type                                        | Description                                                           |
| ----------------------- | ------------------------------------------- | --------------------------------------------------------------------- |
| `SortType`              | `AH_ENUM_POSITIONS_HISTORY_SORT_TYPE`       | Sort mode (0=open time asc, 1=open time desc, 2=ticket asc, 3=ticket desc) |
| `PositionOpenTimeFrom`  | `google.protobuf.Timestamp`                 | Start date (optional)                                                 |
| `PositionOpenTimeTo`    | `google.protobuf.Timestamp`                 | End date (optional)                                                   |
| `PageNumber`            | `int32`                                     | Page number for pagination (optional, default 0)                      |
| `ItemsPerPage`          | `int32`                                     | Number of items per page (optional, default 0 = all)                  |

---

## ⬆️ Output — `PositionsHistoryData`

| Field              | Type                     | Go Type                       | Description                              |
| ------------------ | ------------------------ | ----------------------------- | ---------------------------------------- |
| `HistoryPositions` | `PositionHistoryInfo[]`  | `[]*pb.PositionHistoryInfo`   | Array of closed positions with P&L data  |

**PositionHistoryInfo main fields:**

- `Index` (int32) - Record index
- `PositionTicket` (uint64) - Position ticket number
- `OrderType` (AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE) - Order type (BUY=0, SELL=1, etc.)

- `OpenTime`, `CloseTime` (google.protobuf.Timestamp) - Open and close timestamps

- `Volume` (double) - Position volume
- `OpenPrice`, `ClosePrice` (double) - Open and close prices
- `StopLoss`, `TakeProfit` (double) - SL/TP levels
- `Profit` (double) - Net profit/loss
- `Commission`, `Fee`, `Swap` (double) - Trading costs
- `Symbol` (string) - Trading symbol
- `Comment` (string) - Position comment
- `Magic` (int64) - Magic number

---

> **💡 Enum Usage Note:** The tables show simplified constant names for readability.
> In Go code, use full names with the enum type prefix.
>
> Format: `pb.<ENUM_TYPE>_<CONSTANT_NAME>`
>
> Example: `pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC`

### 📘 Enum: AH_ENUM_POSITIONS_HISTORY_SORT_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `AH_POSITION_OPEN_TIME_ASC` | Sort by position open time (ascending) |
| 1 | `AH_POSITION_OPEN_TIME_DESC` | Sort by position open time (descending) |
| 2 | `AH_POSITION_TICKET_ASC` | Sort by position ticket (ascending) |
| 3 | `AH_POSITION_TICKET_DESC` | Sort by position ticket (descending) |

### 📘 Enum: AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `AH_ORDER_TYPE_BUY` | Market Buy order |
| 1 | `AH_ORDER_TYPE_SELL` | Market Sell order |
| 2 | `AH_ORDER_TYPE_BUY_LIMIT` | Buy Limit pending order |
| 3 | `AH_ORDER_TYPE_SELL_LIMIT` | Sell Limit pending order |
| 4 | `AH_ORDER_TYPE_BUY_STOP` | Buy Stop pending order |
| 5 | `AH_ORDER_TYPE_SELL_STOP` | Sell Stop pending order |
| 6 | `AH_ORDER_TYPE_BUY_STOP_LIMIT` | Buy Stop Limit (pending Buy Limit order at StopLimit price) |
| 7 | `AH_ORDER_TYPE_SELL_STOP_LIMIT` | Sell Stop Limit (pending Sell Limit order at StopLimit price) |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `15s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Time range:** PositionOpenTimeFrom and PositionOpenTimeTo are `google.protobuf.Timestamp` (use `timestamppb.New()` to convert).
* **Pagination:** Use PageNumber and ItemsPerPage for large datasets.
* **Sort mode:** SortType allows sorting by open time or ticket number (ascending/descending).

---

## 🔗 Usage Examples

### 1) Get last month's closed positions

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

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    now := time.Now()
    monthAgo := now.Add(-30 * 24 * time.Hour)

    data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
        PositionOpenTimeFrom: timestamppb.New(monthAgo),
        PositionOpenTimeTo:   timestamppb.New(now),
        PageNumber:           0,
        ItemsPerPage:         100,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Closed positions last 30 days: %d\n", len(data.HistoryPositions))

    var totalProfit float64
    for _, pos := range data.HistoryPositions {
        totalProfit += pos.Profit
    }
    fmt.Printf("Total profit: %.2f\n", totalProfit)
}
```

### 2) Calculate win rate

```go
func CalculateWinRate(account *mt5.MT5Account, days int) (float64, error) {
    ctx := context.Background()

    now := time.Now()
    fromDate := now.Add(-time.Duration(days) * 24 * time.Hour)

    data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
        PositionOpenTimeFrom: timestamppb.New(fromDate),
        PositionOpenTimeTo:   timestamppb.New(now),
    })
    if err != nil {
        return 0, err
    }

    if len(data.HistoryPositions) == 0 {
        return 0, nil
    }

    wins := 0
    for _, pos := range data.HistoryPositions {
        if pos.Profit > 0 {
            wins++
        }
    }

    winRate := float64(wins) / float64(len(data.HistoryPositions)) * 100.0
    return winRate, nil
}

// Usage:
// winRate, _ := CalculateWinRate(account, 30)
// fmt.Printf("Win rate (30 days): %.2f%%\n", winRate)
```

### 3) Get detailed P&L breakdown

```go
type PnLBreakdown struct {
    GrossProfit   float64
    GrossLoss     float64
    NetProfit     float64
    TotalSwap     float64
    TotalCommission float64
    Trades        int
    WinningTrades int
    LosingTrades  int
}

func GetPnLBreakdown(account *mt5.MT5Account, fromDate, toDate time.Time) (*PnLBreakdown, error) {
    ctx := context.Background()

    data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
        PositionOpenTimeFrom: timestamppb.New(fromDate),
        PositionOpenTimeTo:   timestamppb.New(toDate),
    })
    if err != nil {
        return nil, err
    }

    breakdown := &PnLBreakdown{
        Trades: len(data.HistoryPositions),
    }

    for _, pos := range data.HistoryPositions {
        breakdown.NetProfit += pos.Profit
        breakdown.TotalSwap += pos.Swap
        breakdown.TotalCommission += pos.Commission

        if pos.Profit > 0 {
            breakdown.GrossProfit += pos.Profit
            breakdown.WinningTrades++
        } else if pos.Profit < 0 {
            breakdown.GrossLoss += pos.Profit
            breakdown.LosingTrades++
        }
    }

    return breakdown, nil
}

// Usage:
// fromDate := time.Now().Add(-30 * 24 * time.Hour)
// toDate := time.Now()
// breakdown, _ := GetPnLBreakdown(account, fromDate, toDate)
// fmt.Printf("Net: %.2f, Gross Profit: %.2f, Gross Loss: %.2f\n",
//     breakdown.NetProfit, breakdown.GrossProfit, breakdown.GrossLoss)
```

### 4) Find best and worst trades

```go
func FindBestWorstTrades(account *mt5.MT5Account, days int) {
    ctx := context.Background()

    now := time.Now()
    fromDate := now.Add(-time.Duration(days) * 24 * time.Hour)

    data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
        PositionOpenTimeFrom: timestamppb.New(fromDate),
        PositionOpenTimeTo:   timestamppb.New(now),
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if len(data.HistoryPositions) == 0 {
        fmt.Println("No closed positions found")
        return
    }

    var bestTrade, worstTrade *pb.PositionHistoryInfo
    for _, pos := range data.HistoryPositions {
        if bestTrade == nil || pos.Profit > bestTrade.Profit {
            bestTrade = pos
        }
        if worstTrade == nil || pos.Profit < worstTrade.Profit {
            worstTrade = pos
        }
    }

    fmt.Printf("Best trade: %s, Profit: %.2f\n", bestTrade.Symbol, bestTrade.Profit)
    fmt.Printf("Worst trade: %s, Profit: %.2f\n", worstTrade.Symbol, worstTrade.Profit)
}
```

### 5) Calculate profitability by symbol

```go
func GetProfitBySymbol(account *mt5.MT5Account, days int) map[string]float64 {
    ctx := context.Background()

    now := time.Now()
    fromDate := now.Add(-time.Duration(days) * 24 * time.Hour)

    data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
        PositionOpenTimeFrom: timestamppb.New(fromDate),
        PositionOpenTimeTo:   timestamppb.New(now),
    })
    if err != nil {
        return nil
    }

    profitBySymbol := make(map[string]float64)
    for _, pos := range data.HistoryPositions {
        profitBySymbol[pos.Symbol] += pos.Profit
    }

    return profitBySymbol
}

// Usage:
// profits := GetProfitBySymbol(account, 30)
// for symbol, profit := range profits {
//     fmt.Printf("%s: %.2f\n", symbol, profit)
// }
```

---

## 🔧 Common Patterns

### Get today's closed positions

```go
func GetTodayClosedPositions(account *mt5.MT5Account) ([]*pb.PositionHistoryInfo, error) {
    ctx := context.Background()

    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

    data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
        PositionOpenTimeFrom: timestamppb.New(startOfDay),
        PositionOpenTimeTo:   timestamppb.New(now),
    })
    if err != nil {
        return nil, err
    }

    return data.HistoryPositions, nil
}
```

### Calculate average profit per trade

```go
func CalculateAverageProfit(account *mt5.MT5Account, days int) (float64, error) {
    ctx := context.Background()

    now := time.Now()
    fromDate := now.Add(-time.Duration(days) * 24 * time.Hour)

    data, err := account.PositionsHistory(ctx, &pb.PositionsHistoryRequest{
        PositionOpenTimeFrom: timestamppb.New(fromDate),
        PositionOpenTimeTo:   timestamppb.New(now),
    })
    if err != nil {
        return 0, err
    }

    if len(data.HistoryPositions) == 0 {
        return 0, nil
    }

    var totalProfit float64
    for _, pos := range data.HistoryPositions {
        totalProfit += pos.Profit
    }

    return totalProfit / float64(len(data.HistoryPositions)), nil
}
```

---

## 📚 See Also

* [OrderHistory](./OrderHistory.md) - Get historical orders
* [OpenedOrders](./OpenedOrders.md) - Get current open positions
* [OnTrade](../7.%20Streaming_Methods/OnTrade.md) - Stream trade events
* [AccountSummary](../1.%20Account_information/AccountSummary.md) - Current account statistics
