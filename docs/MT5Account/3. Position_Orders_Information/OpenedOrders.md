# ✅ Get All Open Positions and Orders

> **Request:** get complete information about all open positions and pending orders, including profit/loss, prices, volumes, timestamps and other details.

**API Information:**

* **Low-level API:** `MT5Account.OpenedOrders(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `OpenedOrders` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `OpenedOrders(OpenedOrdersRequest) → OpenedOrdersReply`
* **Low‑level client (generated):** `AccountHelperClient.OpenedOrders(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Returns complete details of all currently open positions and pending orders.
* **Why you need it.** Monitor P&L, manage positions, analyze trading activity.
* **Comprehensive data.** Includes all order parameters, current prices, and profit calculations.

---

## 🎯 Purpose

Use it to:

* Monitor current P&L of open positions
* List all active trading positions
* Check order parameters (SL/TP/prices)
* Calculate total exposure
* Implement position management logic

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OpenedOrders - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/OpenedOrders_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// OpenedOrders retrieves all currently opened orders and positions with full details.
// Returns comprehensive information including profit/loss, prices, volumes, and timestamps.
func (a *MT5Account) OpenedOrders(
    ctx context.Context,
    req *pb.OpenedOrdersRequest,
) (*pb.OpenedOrdersData, error)
```

**Request message:**

```protobuf
OpenedOrdersRequest {
  BMT5_ENUM_OPENED_ORDER_SORT_TYPE InputSortMode = 1;  // Sort order for results
}
```

**Reply message:**

```protobuf
OpenedOrdersReply {
  oneof response {
    OpenedOrdersData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                         | Description                                   |
| --------- | ---------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`            | Context for deadline/timeout and cancellation |
| `req`     | `*pb.OpenedOrdersRequest`    | Request with sort mode parameter              |

**Request fields:**

| Field          | Type                                 | Description                                        |
| -------------- | ------------------------------------ | -------------------------------------------------- |
| `InputSortMode` | `BMT5_ENUM_OPENED_ORDER_SORT_TYPE` | Sort order for returned orders and positions       |

---

## ⬆️ Output — `OpenedOrdersData`

| Field           | Type                  | Go Type                    | Description                          |
| --------------- | --------------------- | -------------------------- | ------------------------------------ |
| `OpenedOrders`  | `OpenedOrderInfo[]`   | `[]*pb.OpenedOrderInfo`    | Array of pending orders              |
| `PositionInfos` | `PositionInfo[]`      | `[]*pb.PositionInfo`       | Array of open positions              |

**PositionInfo structure includes:**

- `Index` - Position index
- `Ticket` - Position ticket number
- `OpenTime` - Position open time (timestamp)
- `Volume` - Position volume
- `PriceOpen` - Open price
- `StopLoss` - Stop loss level
- `TakeProfit` - Take profit level
- `PriceCurrent` - Current market price
- `Swap` - Swap charges
- `Profit` - Current profit/loss
- `LastUpdateTime` - Last update timestamp
- `Type` - Position type (BMT5_ENUM_POSITION_TYPE)
- `MagicNumber` - Magic number
- `Identifier` - Position identifier
- `Reason` - Position open reason (BMT5_ENUM_POSITION_REASON)
- `Symbol` - Trading symbol
- `Comment` - Position comment
- `ExternalId` - External identifier
- `PositionCommission` - Position commission
- `AccountLogin` - Account login

**OpenedOrderInfo structure (pending orders):**

- `Index` - Order index
- `Ticket` - Order ticket number
- `PriceCurrent` - Current market price
- `PriceOpen` - Order open price
- `StopLimit` - Stop limit price (for stop-limit orders)
- `StopLoss` - Stop loss level
- `TakeProfit` - Take profit level
- `VolumeCurrent` - Current volume (for partially filled orders)
- `VolumeInitial` - Initial volume
- `MagicNumber` - Magic number
- `Reason` - Order creation reason
- `Type` - Order type (BMT5_ENUM_ORDER_TYPE)
- `State` - Order state (BMT5_ENUM_ORDER_STATE)
- `TimeExpiration` - Expiration time (timestamp)
- `TimeSetup` - Setup time (timestamp)
- `TimeDone` - Done time (timestamp)
- `TypeFilling` - Order filling type (BMT5_ENUM_ORDER_TYPE_FILLING)
- `TypeTime` - Order time type (BMT5_ENUM_ORDER_TYPE_TIME)
- `PositionId` - Position ID
- `PositionById` - Position by ID
- `Symbol` - Trading symbol
- `ExternalId` - External identifier
- `Comment` - Order comment
- `AccountLogin` - Account login

---

> **💡 Enum Usage Note:** The tables show simplified constant names for readability.
> In Go code, use full names with the enum type prefix.
>
> Format: `pb.<ENUM_TYPE>_<CONSTANT_NAME>`
>
> Example: `pb.BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY`

### 📘 Enum: BMT5_ENUM_ORDER_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_ORDER_TYPE_BUY` | Market Buy order (long position) |
| 1 | `BMT5_ORDER_TYPE_SELL` | Market Sell order (short position) |
| 2 | `BMT5_ORDER_TYPE_BUY_LIMIT` | Buy Limit pending order |
| 3 | `BMT5_ORDER_TYPE_SELL_LIMIT` | Sell Limit pending order |
| 4 | `BMT5_ORDER_TYPE_BUY_STOP` | Buy Stop pending order |
| 5 | `BMT5_ORDER_TYPE_SELL_STOP` | Sell Stop pending order |
| 6 | `BMT5_ORDER_TYPE_BUY_STOP_LIMIT` | Buy Stop Limit (pending Buy Limit order at StopLimit price) |
| 7 | `BMT5_ORDER_TYPE_SELL_STOP_LIMIT` | Sell Stop Limit (pending Sell Limit order at StopLimit price) |
| 8 | `BMT5_ORDER_TYPE_CLOSE_BY` | Order to close a position by an opposite one |

**Note:** For open positions (PositionInfo), use BMT5_ENUM_POSITION_TYPE instead. For pending orders (OpenedOrderInfo), all types 0-8 can be present.

### 📘 Enum: BMT5_ENUM_POSITION_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_POSITION_TYPE_BUY` | Long position (Buy) |
| 1 | `BMT5_POSITION_TYPE_SELL` | Short position (Sell) |

**Used in:** PositionInfo.Type

### 📘 Enum: BMT5_ENUM_POSITION_REASON

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_POSITION_REASON_CLIENT` | Position opened from desktop terminal |
| 1 | `BMT5_POSITION_REASON_MOBILE` | Position opened from mobile application |
| 2 | `BMT5_POSITION_REASON_WEB` | Position opened from web platform |
| 3 | `BMT5_POSITION_REASON_EXPERT` | Position opened by Expert Advisor or script |
| 4 | `ORDER_REASON_SL` | Position closed by Stop Loss |
| 5 | `ORDER_REASON_TP` | Position closed by Take Profit |
| 6 | `ORDER_REASON_SO` | Position closed by Stop Out |

**Used in:** PositionInfo.Reason

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
| 7 | `BMT5_ORDER_STATE_REQUEST_ADD` | Order is being registered (placing to the trading system) |
| 8 | `BMT5_ORDER_STATE_REQUEST_MODIFY` | Order is being modified (changing its parameters) |
| 9 | `BMT5_ORDER_STATE_REQUEST_CANCEL` | Order is being deleted (deleting from the trading system) |

**Used in:** OpenedOrderInfo.State

### 📘 Enum: BMT5_ENUM_ORDER_TYPE_FILLING

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_ORDER_FILLING_FOK` | Fill or Kill - order must be filled completely or not at all |
| 1 | `BMT5_ORDER_FILLING_IOC` | Immediate or Cancel - fill available volume, cancel the rest |
| 2 | `BMT5_ORDER_FILLING_RETURN` | Return - order is placed with the broker for execution |
| 3 | `BMT5_ORDER_FILLING_BOC` | Book or Cancel - order must be placed as a passive order (limit) |

**Used in:** OpenedOrderInfo.TypeFilling

### 📘 Enum: BMT5_ENUM_ORDER_TYPE_TIME

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_ORDER_TIME_GTC` | Good Till Cancelled - order stays until explicitly cancelled |
| 1 | `BMT5_ORDER_TIME_DAY` | Good Till Day - order valid until end of trading day |
| 2 | `BMT5_ORDER_TIME_SPECIFIED` | Good Till Specified - order valid until specified date/time |
| 3 | `BMT5_ORDER_TIME_SPECIFIED_DAY` | Good Till Specified Day - order valid until end of specified day |

**Used in:** OpenedOrderInfo.TypeTime

### 📘 Enum: BMT5_ENUM_OPENED_ORDER_SORT_TYPE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC` | Sort by open time (ascending) |
| 1 | `BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_DESC` | Sort by open time (descending) |
| 2 | `BMT5_OPENED_ORDER_SORT_BY_ORDER_TICKET_ID_ASC` | Sort by order ticket ID (ascending) |
| 3 | `BMT5_OPENED_ORDER_SORT_BY_ORDER_TICKET_ID_DESC` | Sort by order ticket ID (descending) |


---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `10s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Sorting:** Use InputSortMode to control the order of returned results (by open time or ticket ID, ascending or descending).
* **Performance:** For large accounts, consider using OpenedOrdersTickets for lightweight checks.

---

## 🔗 Usage Examples

### 1) Get all open positions

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

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC,
    })
    if err != nil {
        panic(err)
    }

    totalCount := len(data.OpenedOrders) + len(data.PositionInfos)
    fmt.Printf("Total items: %d (Pending: %d, Positions: %d)\n\n",
        totalCount, len(data.OpenedOrders), len(data.PositionInfos))

    // Display positions
    for _, pos := range data.PositionInfos {
        fmt.Printf("Position Ticket: %d\n", pos.Ticket)
        fmt.Printf("  Symbol: %s\n", pos.Symbol)
        fmt.Printf("  Type: %d\n", pos.Type)
        fmt.Printf("  Volume: %.2f\n", pos.Volume)
        fmt.Printf("  Open Price: %.5f\n", pos.PriceOpen)
        fmt.Printf("  Current Price: %.5f\n", pos.PriceCurrent)
        fmt.Printf("  Profit: %.2f\n", pos.Profit)
        fmt.Printf("  SL: %.5f, TP: %.5f\n\n", pos.StopLoss, pos.TakeProfit)
    }
}
```

### 2) Get positions sorted by ticket ID

```go
func GetPositionsSortedByTicket(account *mt5.MT5Account) ([]*pb.PositionInfo, error) {
    ctx := context.Background()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_ORDER_TICKET_ID_ASC,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get positions: %w", err)
    }

    return data.PositionInfos, nil
}

// Usage:
// positions, _ := GetPositionsSortedByTicket(account)
// fmt.Printf("Positions sorted by ticket: %d\n", len(positions))
```

### 3) Calculate total profit

```go
func GetTotalProfit(account *mt5.MT5Account) (float64, error) {
    ctx := context.Background()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC,
    })
    if err != nil {
        return 0, err
    }

    var totalProfit float64
    for _, pos := range data.PositionInfos {
        totalProfit += pos.Profit
    }

    return totalProfit, nil
}

// Usage:
// profit, _ := GetTotalProfit(account)
// fmt.Printf("Total P&L: %.2f\n", profit)
```

### 4) Find losing positions

```go
func GetLosingPositions(account *mt5.MT5Account) ([]*pb.PositionInfo, error) {
    ctx := context.Background()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC,
    })
    if err != nil {
        return nil, err
    }

    losingPositions := []*pb.PositionInfo{}
    for _, pos := range data.PositionInfos {
        if pos.Profit < 0 {
            losingPositions = append(losingPositions, pos)
        }
    }

    return losingPositions, nil
}

// Usage:
// losers, _ := GetLosingPositions(account)
// for _, pos := range losers {
//     fmt.Printf("Losing position %d: %.2f\n", pos.Ticket, pos.Profit)
// }
```

### 5) Monitor position details

```go
func DisplayPositionDetails(account *mt5.MT5Account) {
    ctx := context.Background()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    var totalProfit, totalSwap, totalCommission float64
    var buyVolume, sellVolume float64

    fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
    fmt.Printf("║ Open Positions Summary                                         ║\n")
    fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")

    for _, pos := range data.PositionInfos {
        fmt.Printf("║ #%-10d %-10s Vol:%-6.2f P&L:%-10.2f ║\n",
            pos.Ticket, pos.Symbol, pos.Volume, pos.Profit)

        totalProfit += pos.Profit
        totalSwap += pos.Swap
        totalCommission += pos.PositionCommission

        if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
            buyVolume += pos.Volume
        } else if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_SELL {
            sellVolume += pos.Volume
        }
    }

    fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
    fmt.Printf("║ Total Profit:     %10.2f                                ║\n", totalProfit)
    fmt.Printf("║ Total Swap:       %10.2f                                ║\n", totalSwap)
    fmt.Printf("║ Total Commission: %10.2f                                ║\n", totalCommission)
    fmt.Printf("║ Buy Volume:       %10.2f                                ║\n", buyVolume)
    fmt.Printf("║ Sell Volume:      %10.2f                                ║\n", sellVolume)
    fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}
```

### 6) Find positions by magic number

```go
func GetPositionsByMagic(account *mt5.MT5Account, magic int64) ([]*pb.PositionInfo, error) {
    ctx := context.Background()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC,
    })
    if err != nil {
        return nil, err
    }

    positions := []*pb.PositionInfo{}
    for _, pos := range data.PositionInfos {
        if pos.MagicNumber == magic {
            positions = append(positions, pos)
        }
    }

    return positions, nil
}
```

---

## 🔧 Common Patterns

### Check if symbol has open positions

```go
func HasPositionsForSymbol(account *mt5.MT5Account, symbol string) (bool, error) {
    ctx := context.Background()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC,
    })
    if err != nil {
        return false, err
    }

    // Filter by symbol
    for _, pos := range data.PositionInfos {
        if pos.Symbol == symbol {
            return true, nil
        }
    }

    return false, nil
}
```

### Calculate exposure by symbol

```go
func GetSymbolExposure(account *mt5.MT5Account, symbol string) (buyVolume, sellVolume float64, err error) {
    ctx := context.Background()

    data, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC,
    })
    if err != nil {
        return 0, 0, err
    }

    for _, pos := range data.PositionInfos {
        if pos.Symbol != symbol {
            continue
        }
        if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
            buyVolume += pos.Volume
        } else if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_SELL {
            sellVolume += pos.Volume
        }
    }

    return buyVolume, sellVolume, nil
}
```

---

## 📚 See Also

* [PositionsTotal](./PositionsTotal.md) - Get quick position count
* [OpenedOrdersTickets](./OpenedOrdersTickets.md) - Get only ticket numbers (lightweight)
* [OnPositionProfit](../7.%20Streaming_Methods/OnPositionProfit.md) - Stream real-time P&L updates
* [OrderClose](../4.%20Trading_Operations/OrderClose.md) - Close positions
