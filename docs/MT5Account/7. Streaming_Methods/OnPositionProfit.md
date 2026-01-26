# ✅ Stream Real-Time Position Profit/Loss

> **Request:** subscribe to real-time profit/loss updates for open positions.

**API Information:**

* **Low-level API:** `MT5Account.OnPositionProfit(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.SubscriptionService`
* **Proto definition:** `OnPositionProfit` (defined in `mt5-term-api-subscription.proto`)

### RPC

* **Service:** `mt5_term_api.SubscriptionService`
* **Method:** `OnPositionProfit(OnPositionProfitRequest) → stream OnPositionProfitReply`
* **Low‑level client (generated):** `SubscriptionServiceClient.OnPositionProfit(ctx, request, opts...)`

```go
package mt5

type MT5Account struct {
    // ...
}

// OnPositionProfit streams real-time profit/loss updates for open positions.
// Returns two channels: data channel and error channel.
func (a *MT5Account) OnPositionProfit(
    ctx context.Context,
    req *pb.OnPositionProfitRequest,
) (<-chan *pb.OnPositionProfitData, <-chan error)
```

**Request message:**

```protobuf
OnPositionProfitRequest {
  int32 TimerPeriodMilliseconds = 1;  // Update interval in milliseconds
  bool IgnoreEmptyData = 2;           // Skip empty updates (no changes)
}
```

---

## 🔽 Input

| Parameter | Type                           | Description                                   |
| --------- | ------------------------------ | --------------------------------------------- |
| `ctx`     | `context.Context`              | Context for cancellation (cancel to stop stream) |
| `req`     | `*pb.OnPositionProfitRequest`  | Request with update interval settings          |
| `req.TimerPeriodMilliseconds` | `int32` | Update interval in milliseconds (e.g., 1000 = 1 second) |
| `req.IgnoreEmptyData` | `bool` | If true, skip updates when no positions changed |

---

## ⬆️ Output — Channels

| Channel      | Type                              | Description                              |
| ------------ | --------------------------------- | ---------------------------------------- |
| Data Channel | `<-chan *pb.OnPositionProfitData` | Receives P&L updates                     |
| Error Channel| `<-chan error`                    | Receives errors (closed on ctx cancel)   |

**OnPositionProfitData fields:**

| Field          | Type     | Go Type   | Description                    |
| -------------- | -------- | --------- | ------------------------------ |
| `Type` | `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (enum) | `int32` | Event type (always OrderProfit) - **ENUM!** |
| `NewPositions` | `repeated OnPositionProfitPositionInfo` | `[]*OnPositionProfitPositionInfo` | Newly opened positions |
| `UpdatedPositions` | `repeated OnPositionProfitPositionInfo` | `[]*OnPositionProfitPositionInfo` | Positions with P&L changes |
| `DeletedPositions` | `repeated OnPositionProfitPositionInfo` | `[]*OnPositionProfitPositionInfo` | Closed positions |
| `AccountInfo` | `OnEventAccountInfo` | `*OnEventAccountInfo` | Account snapshot (balance, equity, etc.) |
| `TerminalInstanceGuidId` | `string` | `string` | Terminal instance ID |

**OnPositionProfitPositionInfo fields:**

| Field          | Type     | Go Type   | Description                    |
| -------------- | -------- | --------- | ------------------------------ |
| `Index`        | `int32` | `int32`  | Position index         |
| `Ticket`       | `int64` | `int64`  | Position ticket number         |
| `Profit`       | `double` | `float64` | Current profit/loss            |
| `PositionSymbol` | `string` | `string` | Trading symbol           |

---

## 💬 Just the essentials

* **What it is.** Real-time stream of position profit/loss changes.
* **Why you need it.** Monitor P&L, implement stop-loss logic, track account performance.
* **Dynamic updates.** Updates whenever prices change affecting position profit.

---

## 🎯 Purpose

Use it to:

* Monitor real-time profit/loss for open positions
* Implement dynamic stop-loss and take-profit logic
* Track account performance
* Set up profit/loss alerts
* Build real-time P&L dashboards
* Detect position changes immediately

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OnPositionProfit - How it works](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnPositionProfit_HOW.md)**

---

## 🧱 ENUMs used in OnPositionProfit

### Output ENUM

| ENUM Type | Field Name | Purpose | Values |
|-----------|------------|---------|--------|
| `MT5_SUB_ENUM_EVENT_GROUP_TYPE` | `Type` | Indicates the event type | `OrderProfit` (0) - Position profit/loss update event |

**Usage in code:**

```go
profitStream, _ := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{
    TimerPeriodMilliseconds: 1000,
    IgnoreEmptyData:         true,
})

for {
    event, err := profitStream.Recv()
    if err != nil {
        break
    }

    // Check event type (always OrderProfit for this stream)
    switch event.Type {
    case pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit:
        fmt.Printf("Position profit updates: %d new, %d updated, %d deleted\n",
            len(event.NewPositions),
            len(event.UpdatedPositions),
            len(event.DeletedPositions))
    }
}
```

**Note:** The `Type` field will always be `OrderProfit` for OnPositionProfit stream. This ENUM is shared across all streaming methods (OnTrade, OnPositionProfit, OnTradeTransaction) to maintain consistent event identification.

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, streams run indefinitely until cancelled.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Channel buffering:** Data channel is buffered (default 100), error channel is buffered (default 10).
* **Goroutine required:** You MUST consume the channels in a separate goroutine to avoid blocking.
* **Timer interval:** Use `TimerPeriodMilliseconds` to control update frequency (e.g., 1000 = update every 1 second).
* **Ignore empty data:** Set `IgnoreEmptyData` to true to skip updates when no positions changed.
* **Three arrays:** `NewPositions` (opened), `UpdatedPositions` (profit changed), `DeletedPositions` (closed).
* **Real-time:** Very low latency updates for position monitoring.
* **ENUM type:** Always check the `Type` field (MT5_SUB_ENUM_EVENT_GROUP_TYPE) to identify event type.

---

## 🔗 Usage Examples

### 1) Monitor all positions P&L

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
    "github.com/google/uuid"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    err := account.Connect()
    if err != nil {
        panic(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{
        TimerPeriodMilliseconds: 1000, // Update every 1 second
        IgnoreEmptyData:         true, // Skip empty updates
    })

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }

                // Check event type (ENUM!)
                if update.Type == pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit {
                    // Process new positions
                    for _, pos := range update.NewPositions {
                        fmt.Printf("[NEW] Position %d (%s): P&L=%.2f\n",
                            pos.Ticket, pos.PositionSymbol, pos.Profit)
                    }

                    // Process updated positions
                    for _, pos := range update.UpdatedPositions {
                        fmt.Printf("[UPD] Position %d (%s): P&L=%.2f\n",
                            pos.Ticket, pos.PositionSymbol, pos.Profit)
                    }

                    // Process closed positions
                    for _, pos := range update.DeletedPositions {
                        fmt.Printf("[DEL] Position %d (%s): Final P&L=%.2f\n",
                            pos.Ticket, pos.PositionSymbol, pos.Profit)
                    }
                }

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Error: %v\n", err)
                    return
                }
            }
        }
    }()

    <-ctx.Done()
}
```

### 2) Real-time total P&L dashboard

```go
func RealTimePnLDashboard(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    positionProfits := make(map[int64]float64) // ticket -> profit

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{
        TimerPeriodMilliseconds: 500,  // Update every 500ms
        IgnoreEmptyData:         false, // Get all updates
    })

    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }

                // Check event type
                if update.Type == pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit {
                    // Update profit for new and updated positions
                    for _, pos := range update.NewPositions {
                        positionProfits[pos.Ticket] = pos.Profit
                    }
                    for _, pos := range update.UpdatedPositions {
                        positionProfits[pos.Ticket] = pos.Profit
                    }

                    // Remove closed positions
                    for _, pos := range update.DeletedPositions {
                        delete(positionProfits, pos.Ticket)
                    }
                }

            case <-ticker.C:
                var totalProfit float64
                for _, profit := range positionProfits {
                    totalProfit += profit
                }

                fmt.Printf("\r Total P&L: %.2f USD (Positions: %d)",
                    totalProfit, len(positionProfits))

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("\nError: %v\n", err)
                    return
                }

            case <-ctx.Done():
                return
            }
        }
    }()

    // Run for 5 minutes
    time.Sleep(5 * time.Minute)
}
```

### 3) Auto take profit monitor

```go
func AutoTakeProfitMonitor(account *mt5.MT5Account, targetProfit float64) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{
        TimerPeriodMilliseconds: 1000,
        IgnoreEmptyData:         true,
    })

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }

                // Check only updated positions
                for _, pos := range update.UpdatedPositions {
                    if pos.Profit >= targetProfit {
                        fmt.Printf("\n🎯 Target profit reached for position %d (%s): %.2f\n",
                            pos.Ticket, pos.PositionSymbol, pos.Profit)

                        // Close position
                        _, err := account.OrderClose(context.Background(), &pb.OrderCloseRequest{
                            Ticket:    uint64(pos.Ticket),
                            Volume:    0, // Close all
                            Deviation: 20,
                            Comment:   "Auto TP",
                        })

                        if err != nil {
                            fmt.Printf("Failed to close position: %v\n", err)
                        } else {
                            fmt.Printf("Position %d closed at profit %.2f\n",
                                pos.Ticket, pos.Profit)
                        }
                    }
                }

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

    time.Sleep(24 * time.Hour)
}

// Usage:
// AutoTakeProfitMonitor(account, 100.0) // Close when profit >= $100
```

---

## 📚 See Also

* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get full position details
* [OrderClose](../4.%20Trading_Operations/OrderClose.md) - Close positions
* [OnSymbolTick](./OnSymbolTick.md) - Stream price updates
