# ✅ Get Total Count of Open Positions

> **Request:** get the count of currently open positions on the trading account. Quick check for open positions without loading full details.

**API Information:**

* **SDK wrapper:** `MT5Account.PositionsTotal(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.TradeFunctions`
* **Proto definition:** `PositionsTotal` (defined in `mt5-term-api-trade-functions.proto`)

### RPC

* **Service:** `mt5_term_api.TradeFunctions`
* **Method:** `PositionsTotal(Empty) → PositionsTotalReply`
* **Low‑level client (generated):** `TradeFunctionsClient.PositionsTotal(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

## 💬 Just the essentials

* **What it is.** Returns the count of currently open positions on the account.
* **Why you need it.** Quickly check if there are open positions without fetching full position details.
* **Lightweight.** Much faster than OpenedOrders when you only need the count.

---

## 🎯 Purpose

Use it to:

* Check if account has any open positions
* Monitor position count changes
* Validate before opening new positions
* Implement position limits
* Dashboard indicators

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [PositionsTotal - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/PositionsTotal_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// PositionsTotal returns the number of currently open positions.
// This is a lightweight call for quick position count check.
func (a *MT5Account) PositionsTotal(
    ctx context.Context,
) (*pb.PositionsTotalData, error)
```

**Request message:**

```protobuf
Empty {}  // No parameters required
```

**Reply message:**

```protobuf
PositionsTotalReply {
  oneof response {
    PositionsTotalData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type              | Description                                   |
| --------- | ----------------- | --------------------------------------------- |
| `ctx`     | `context.Context` | Context for deadline/timeout and cancellation |

**No request parameters required.**

---

## ⬆️ Output — `PositionsTotalData`

| Field            | Type    | Go Type | Description                       |
| ---------------- | ------- | ------- | --------------------------------- |
| `TotalPositions` | `int64` | `int64` | Total number of open positions    |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Performance:** Use this instead of OpenedOrders when you only need position count.
* **Real-time:** Returns current count at the moment of request.

---

## 🔗 Usage Examples

### 1) Get positions count

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/mt5"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    data, err := account.PositionsTotal(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Open positions: %d\n", data.TotalPositions)
}
```

### 2) Check if positions are open

```go
func HasOpenPositions(account *mt5.MT5Account) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.PositionsTotal(ctx)
    if err != nil {
        return false, fmt.Errorf("failed to get positions total: %w", err)
    }

    return data.TotalPositions > 0, nil
}

// Usage:
// hasPositions, _ := HasOpenPositions(account)
// if hasPositions {
//     fmt.Println("Account has open positions")
// }
```

### 3) Enforce maximum positions limit

```go
func CanOpenNewPosition(account *mt5.MT5Account, maxPositions int64) (bool, error) {
    ctx := context.Background()

    data, err := account.PositionsTotal(ctx)
    if err != nil {
        return false, err
    }

    if data.TotalPositions >= maxPositions {
        return false, fmt.Errorf("maximum positions limit reached: %d/%d", data.TotalPositions, maxPositions)
    }

    return true, nil
}

// Usage before opening position:
// canOpen, err := CanOpenNewPosition(account, 10)
// if !canOpen {
//     fmt.Println("Cannot open new position:", err)
// }
```

### 4) Monitor position count changes

```go
func MonitorPositionCount(account *mt5.MT5Account, interval time.Duration) {
    ctx := context.Background()
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    var lastCount int64 = -1

    for range ticker.C {
        data, err := account.PositionsTotal(ctx)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        if lastCount != -1 && data.TotalPositions != lastCount {
            change := data.TotalPositions - lastCount
            fmt.Printf("[%s] Position count changed: %d -> %d (%+d)\n",
                time.Now().Format("15:04:05"),
                lastCount,
                data.TotalPositions,
                change,
            )
        }

        lastCount = data.TotalPositions
    }
}

// Usage in goroutine:
// go MonitorPositionCount(account, 2*time.Second)
```

### 5) Wait for all positions to close

```go
func WaitForAllPositionsClosed(account *mt5.MT5Account, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            data, err := account.PositionsTotal(ctx)
            if err != nil {
                return err
            }

            if data.TotalPositions == 0 {
                fmt.Println("All positions closed")
                return nil
            }

            fmt.Printf("Waiting for %d positions to close...\n", data.TotalPositions)

        case <-ctx.Done():
            return fmt.Errorf("timeout waiting for positions to close")
        }
    }
}

// Usage:
// err := WaitForAllPositionsClosed(account, 30*time.Second)
```

### 6) Position count statistics

```go
type PositionStats struct {
    CurrentCount int64
    MaxCount     int64
    MinCount     int64
    CheckCount   int64
}

func TrackPositionStats(account *mt5.MT5Account, duration time.Duration, interval time.Duration) *PositionStats {
    ctx := context.Background()
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    deadline := time.Now().Add(duration)
    stats := &PositionStats{
        MinCount: int64(^uint64(0) >> 1), // Max int64
    }

    for time.Now().Before(deadline) {
        <-ticker.C

        data, err := account.PositionsTotal(ctx)
        if err != nil {
            continue
        }

        stats.CurrentCount = data.TotalPositions
        stats.CheckCount++

        if data.TotalPositions > stats.MaxCount {
            stats.MaxCount = data.TotalPositions
        }
        if data.TotalPositions < stats.MinCount {
            stats.MinCount = data.TotalPositions
        }
    }

    return stats
}

// Usage:
// stats := TrackPositionStats(account, 1*time.Minute, 5*time.Second)
// fmt.Printf("Min: %d, Max: %d, Current: %d\n", stats.MinCount, stats.MaxCount, stats.CurrentCount)
```

---

## 🔧 Common Patterns

### Safe position count check

```go
func GetPositionCountSafe(account *mt5.MT5Account) int64 {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    data, err := account.PositionsTotal(ctx)
    if err != nil {
        return 0
    }

    return data.TotalPositions
}
```

### Conditional trading based on position count

```go
func TradeIfBelowLimit(account *mt5.MT5Account, maxPositions int64, tradeFn func() error) error {
    ctx := context.Background()

    data, err := account.PositionsTotal(ctx)
    if err != nil {
        return err
    }

    if data.TotalPositions >= maxPositions {
        return fmt.Errorf("position limit reached: %d/%d", data.TotalPositions, maxPositions)
    }

    return tradeFn()
}
```

---

## 📚 See Also

* [OpenedOrders](./OpenedOrders.md) - Get full details of all open positions
* [OpenedOrdersTickets](./OpenedOrdersTickets.md) - Get only ticket numbers (lightweight)
* [OnPositionsAndPendingOrdersTickets](../7.%20Streaming_Methods/OnPositionsAndPendingOrdersTickets.md) - Stream position changes
* [OrderSend](../4.%20Trading_Operations/OrderSend.md) - Open new positions
