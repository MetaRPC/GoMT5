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
  string Symbol = 1;  // Optional symbol filter (empty = all symbols)
}
```

---

## 🔽 Input

| Parameter | Type                           | Description                                   |
| --------- | ------------------------------ | --------------------------------------------- |
| `ctx`     | `context.Context`              | Context for cancellation (cancel to stop stream) |
| `req`     | `*pb.OnPositionProfitRequest`  | Request with optional Symbol filter           |

---

## ⬆️ Output — Channels

| Channel      | Type                              | Description                              |
| ------------ | --------------------------------- | ---------------------------------------- |
| Data Channel | `<-chan *pb.OnPositionProfitData` | Receives P&L updates                     |
| Error Channel| `<-chan error`                    | Receives errors (closed on ctx cancel)   |

**OnPositionProfitData fields:**

| Field          | Type     | Go Type   | Description                    |
| -------------- | -------- | --------- | ------------------------------ |
| `Ticket`       | `uint64` | `uint64`  | Position ticket number         |
| `Symbol`       | `string` | `string`  | Trading symbol                 |
| `Profit`       | `double` | `float64` | Current profit/loss            |
| `CurrentPrice` | `double` | `float64` | Current market price           |

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

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, streams run indefinitely until cancelled.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Channel buffering:** Data channel is buffered (default 100), error channel is buffered (default 10).
* **Goroutine required:** You MUST consume the channels in a separate goroutine to avoid blocking.
* **Symbol filter:** Use Symbol parameter to filter specific symbol, or empty string for all positions.
* **Price-driven:** Updates trigger on price changes that affect position profit.
* **Real-time:** Very low latency updates for position monitoring.

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
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{
        Symbol: "", // All symbols
    })

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }
                fmt.Printf("[%s] Position %d (%s): P&L=%.2f, Price=%.5f\n",
                    time.Now().Format("15:04:05"),
                    update.Ticket,
                    update.Symbol,
                    update.Profit,
                    update.CurrentPrice)

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

    positionProfits := make(map[uint64]float64) // ticket -> profit

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{})

    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }
                positionProfits[update.Ticket] = update.Profit

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

### 3) Auto stop-loss on profit target

```go
func AutoTakeProfitMonitor(account *mt5.MT5Account, targetProfit float64) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{})

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }

                if update.Profit >= targetProfit {
                    fmt.Printf("\n🎯 Target profit reached for position %d: %.2f\n",
                        update.Ticket, update.Profit)

                    // Close position
                    _, err := account.OrderClose(context.Background(), &pb.OrderCloseRequest{
                        Ticket:    update.Ticket,
                        Volume:    0, // Close all
                        Deviation: 20,
                        Comment:   "Auto TP",
                    })

                    if err != nil {
                        fmt.Printf("Failed to close position: %v\n", err)
                    } else {
                        fmt.Printf("Position %d closed at profit %.2f\n",
                            update.Ticket, update.Profit)
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

### 4) Loss limit protection

```go
func LossLimitProtection(account *mt5.MT5Account, maxLossPerPosition float64) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{})

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }

                if update.Profit < -maxLossPerPosition {
                    fmt.Printf("\n⚠️  Loss limit exceeded for position %d: %.2f\n",
                        update.Ticket, update.Profit)

                    // Emergency close
                    account.OrderClose(context.Background(), &pb.OrderCloseRequest{
                        Ticket:    update.Ticket,
                        Volume:    0,
                        Deviation: 50, // Allow higher slippage for emergency
                        Comment:   "Emergency stop",
                    })

                    fmt.Printf("Emergency close executed for position %d\n", update.Ticket)
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

    select {}
}

// Usage:
// LossLimitProtection(account, 50.0) // Close if loss > $50
```

### 5) P&L statistics tracker

```go
type PnLStats struct {
    MaxProfit  float64
    MaxLoss    float64
    UpdateCount int64
}

func TrackPnLStats(account *mt5.MT5Account, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    stats := make(map[uint64]*PnLStats) // ticket -> stats

    dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{})

    for {
        select {
        case update := <-dataChan:
            if update == nil {
                break
            }

            if _, exists := stats[update.Ticket]; !exists {
                stats[update.Ticket] = &PnLStats{
                    MaxProfit: update.Profit,
                    MaxLoss:   update.Profit,
                }
            }

            s := stats[update.Ticket]
            s.UpdateCount++

            if update.Profit > s.MaxProfit {
                s.MaxProfit = update.Profit
            }
            if update.Profit < s.MaxLoss {
                s.MaxLoss = update.Profit
            }

        case err := <-errChan:
            if err != nil {
                fmt.Printf("Error: %v\n", err)
            }

        case <-ctx.Done():
            fmt.Println("\nP&L Statistics:")
            for ticket, s := range stats {
                fmt.Printf("  Position %d:\n", ticket)
                fmt.Printf("    Max profit: %.2f\n", s.MaxProfit)
                fmt.Printf("    Max loss: %.2f\n", s.MaxLoss)
                fmt.Printf("    Updates: %d\n", s.UpdateCount)
            }
            return
        }
    }
}

// Usage:
// TrackPnLStats(account, 10*time.Minute)
```

---

## 📚 See Also

* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get full position details
* [OrderClose](../4.%20Trading_Operations/OrderClose.md) - Close positions
* [OnSymbolTick](./OnSymbolTick.md) - Stream price updates
