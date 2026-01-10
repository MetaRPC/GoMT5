# ✅ Stream Position and Order Ticket Changes

> **Request:** subscribe to real-time notifications about changes in open position and pending order ticket lists. Receive updates when positions open/close or orders are placed/removed.

**API Information:**

* **SDK wrapper:** `MT5Account.OnPositionsAndPendingOrdersTickets(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.SubscriptionService`
* **Proto definition:** `OnPositionsAndPendingOrdersTickets` (defined in `mt5-term-api-subscription.proto`)

### RPC

* **Service:** `mt5_term_api.SubscriptionService`
* **Method:** `OnPositionsAndPendingOrdersTickets(OnPositionsAndPendingOrdersTicketsRequest) → stream OnPositionsAndPendingOrdersTicketsReply`
* **Low‑level client (generated):** `SubscriptionServiceClient.OnPositionsAndPendingOrdersTickets(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

```go
package mt5

type MT5Account struct {
    // ...
}

// OnPositionsAndPendingOrdersTickets streams changes in open positions and pending orders.
// Returns two channels: data channel and error channel.
func (a *MT5Account) OnPositionsAndPendingOrdersTickets(
    ctx context.Context,
    req *pb.OnPositionsAndPendingOrdersTicketsRequest,
) (<-chan *pb.OnPositionsAndPendingOrdersTicketsData, <-chan error)
```

---

## 🔽 Input

| Parameter | Type                                               | Description                                   |
| --------- | -------------------------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                                  | Context for cancellation (cancel to stop stream) |
| `req`     | `*pb.OnPositionsAndPendingOrdersTicketsRequest`    | Request (empty structure)                     |

---

## ⬆️ Output — Channels

| Channel      | Type                                                   | Description                              |
| ------------ | ------------------------------------------------------ | ---------------------------------------- |
| Data Channel | `<-chan *pb.OnPositionsAndPendingOrdersTicketsData`    | Receives ticket list updates             |
| Error Channel| `<-chan error`                                         | Receives errors (closed on ctx cancel)   |

**OnPositionsAndPendingOrdersTicketsData fields:**

| Field                  | Type       | Go Type      | Description                              |
| ---------------------- | ---------- | ------------ | ---------------------------------------- |
| `PositionTickets`      | `uint64[]` | `[]uint64`   | Array of open position ticket numbers    |
| `PendingOrderTickets`  | `uint64[]` | `[]uint64`   | Array of pending order ticket numbers    |

---

## 💬 Just the essentials

* **What it is.** Real-time stream of changes to position and pending order ticket lists.
* **Why you need it.** Monitor position/order changes, detect new trades, track closures.
* **Lightweight.** Only sends ticket numbers, not full order details.

---

## 🎯 Purpose

Use it to:

* Detect new positions and pending orders immediately
* Monitor position and order closures in real-time
* Track ticket list changes efficiently
* Implement lightweight position monitoring
* Build real-time order tracking systems
* Detect trading activity changes

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OnPositionsAndPendingOrdersTickets - How it works](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnPositionsAndPendingOrdersTickets_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, streams run indefinitely until cancelled.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Channel buffering:** Data channel is buffered (default 100), error channel is buffered (default 10).
* **Goroutine required:** You MUST consume the channels in a separate goroutine to avoid blocking.
* **Lightweight:** Only transmits ticket numbers, not full order details.
* **Empty request:** OnPositionsAndPendingOrdersTicketsRequest is an empty structure.
* **Change detection:** Compare previous ticket list with new list to detect changes.

---

## 🔗 Usage Examples

### 1) Monitor position changes

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
    "github.com/MetaRPC/GoMT5/mt5"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    dataChan, errChan := account.OnPositionsAndPendingOrdersTickets(ctx,
        &pb.OnPositionsAndPendingOrdersTicketsRequest{})

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }
                fmt.Printf("[%s] Open positions: %d, Pending orders: %d\n",
                    time.Now().Format("15:04:05"),
                    len(update.PositionTickets),
                    len(update.PendingOrderTickets))

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

### 2) Detect new positions

```go
func DetectNewPositions(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var lastPositions map[uint64]bool

    dataChan, errChan := account.OnPositionsAndPendingOrdersTickets(ctx,
        &pb.OnPositionsAndPendingOrdersTicketsRequest{})

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }

                currentPositions := make(map[uint64]bool)
                for _, ticket := range update.PositionTickets {
                    currentPositions[ticket] = true

                    // Check if this is a new position
                    if lastPositions != nil && !lastPositions[ticket] {
                        fmt.Printf("🆕 New position opened: %d\n", ticket)
                    }
                }

                // Check for closed positions
                if lastPositions != nil {
                    for ticket := range lastPositions {
                        if !currentPositions[ticket] {
                            fmt.Printf("🔒 Position closed: %d\n", ticket)
                        }
                    }
                }

                lastPositions = currentPositions

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
```

### 3) Track pending order changes

```go
func TrackPendingOrders(account *mt5.MT5Account, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    var lastPending []uint64

    dataChan, errChan := account.OnPositionsAndPendingOrdersTickets(ctx,
        &pb.OnPositionsAndPendingOrdersTicketsRequest{})

    for {
        select {
        case update := <-dataChan:
            if update == nil {
                break
            }

            // Find new pending orders
            newOrders := findNewTickets(lastPending, update.PendingOrderTickets)
            for _, ticket := range newOrders {
                fmt.Printf("➕ Pending order placed: %d\n", ticket)
            }

            // Find removed pending orders (executed or cancelled)
            removed := findNewTickets(update.PendingOrderTickets, lastPending)
            for _, ticket := range removed {
                fmt.Printf("➖ Pending order removed: %d\n", ticket)
            }

            lastPending = update.PendingOrderTickets

        case err := <-errChan:
            if err != nil {
                fmt.Printf("Error: %v\n", err)
            }

        case <-ctx.Done():
            return
        }
    }
}

func findNewTickets(old, new []uint64) []uint64 {
    oldSet := make(map[uint64]bool)
    for _, t := range old {
        oldSet[t] = true
    }

    result := []uint64{}
    for _, t := range new {
        if !oldSet[t] {
            result = append(result, t)
        }
    }
    return result
}
```

### 4) Position counter with alerts

```go
func MonitorPositionCount(account *mt5.MT5Account, maxPositions int) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnPositionsAndPendingOrdersTickets(ctx,
        &pb.OnPositionsAndPendingOrdersTicketsRequest{})

    go func() {
        for {
            select {
            case update := <-dataChan:
                if update == nil {
                    return
                }

                posCount := len(update.PositionTickets)

                if posCount >= maxPositions {
                    fmt.Printf("⚠️  WARNING: %d positions open (max: %d)\n",
                        posCount, maxPositions)
                } else {
                    fmt.Printf("Positions: %d/%d\n", posCount, maxPositions)
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
// MonitorPositionCount(account, 10) // Alert when >= 10 positions
```

### 5) Trading activity logger

```go
type ActivityLog struct {
    Timestamp        time.Time
    PositionCount    int
    PendingCount     int
    NewPositions     []uint64
    ClosedPositions  []uint64
}

func LogTradingActivity(account *mt5.MT5Account, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    logs := []ActivityLog{}
    var lastPositions map[uint64]bool

    dataChan, errChan := account.OnPositionsAndPendingOrdersTickets(ctx,
        &pb.OnPositionsAndPendingOrdersTicketsRequest{})

    for {
        select {
        case update := <-dataChan:
            if update == nil {
                break
            }

            log := ActivityLog{
                Timestamp:     time.Now(),
                PositionCount: len(update.PositionTickets),
                PendingCount:  len(update.PendingOrderTickets),
            }

            // Track position changes
            currentPos := make(map[uint64]bool)
            for _, ticket := range update.PositionTickets {
                currentPos[ticket] = true
                if lastPositions != nil && !lastPositions[ticket] {
                    log.NewPositions = append(log.NewPositions, ticket)
                }
            }

            if lastPositions != nil {
                for ticket := range lastPositions {
                    if !currentPos[ticket] {
                        log.ClosedPositions = append(log.ClosedPositions, ticket)
                    }
                }
            }

            logs = append(logs, log)
            lastPositions = currentPos

        case err := <-errChan:
            if err != nil {
                fmt.Printf("Error: %v\n", err)
            }

        case <-ctx.Done():
            // Print summary
            fmt.Println("\nTrading Activity Summary:")
            for _, log := range logs {
                if len(log.NewPositions) > 0 || len(log.ClosedPositions) > 0 {
                    fmt.Printf("[%s] Positions: %d, Pending: %d\n",
                        log.Timestamp.Format("15:04:05"),
                        log.PositionCount,
                        log.PendingCount)

                    if len(log.NewPositions) > 0 {
                        fmt.Printf("  New: %v\n", log.NewPositions)
                    }
                    if len(log.ClosedPositions) > 0 {
                        fmt.Printf("  Closed: %v\n", log.ClosedPositions)
                    }
                }
            }
            return
        }
    }
}

// Usage:
// LogTradingActivity(account, 30*time.Minute)
```

---

## 📚 See Also

* [OpenedOrdersTickets](../3.%20Position_Orders_Information/OpenedOrdersTickets.md) - Get current ticket snapshot
* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get full position details
* [OnTrade](./OnTrade.md) - Stream all trade events
