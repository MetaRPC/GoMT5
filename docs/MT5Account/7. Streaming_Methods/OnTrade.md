# ✅ Stream Real-Time Trade Events

> **Request:** subscribe to real-time notifications for all trading operations: order placement, modification, execution, and cancellation.

**API Information:**

* **Low-level API:** `MT5Account.OnTrade(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.SubscriptionService`
* **Proto definition:** `OnTrade` (defined in `mt5-term-api-subscription.proto`)

### RPC

* **Service:** `mt5_term_api.SubscriptionService`
* **Method:** `OnTrade(OnTradeRequest) → stream OnTradeReply`
* **Low‑level client (generated):** `SubscriptionServiceClient.OnTrade(ctx, request, opts...)`

```go
package mt5

type MT5Account struct {
    // ...
}

// OnTrade streams trade events in real-time.
// Returns two channels: data channel and error channel.
func (a *MT5Account) OnTrade(
    ctx context.Context,
    req *pb.OnTradeRequest,
) (<-chan *pb.OnTradeData, <-chan error)
```

---

## 🔽 Input

| Parameter | Type                  | Description                                   |
| --------- | --------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`     | Context for cancellation (cancel to stop stream) |
| `req`     | `*pb.OnTradeRequest`  | Request (empty structure)                     |

---

## ⬆️ Output — Channels

| Channel      | Type                      | Description                              |
| ------------ | ------------------------- | ---------------------------------------- |
| Data Channel | `<-chan *pb.OnTradeData`  | Receives trade event updates             |
| Error Channel| `<-chan error`            | Receives errors (closed on ctx cancel)   |

**OnTradeData fields:**

| Field          | Type     | Go Type   | Description                    |
| -------------- | -------- | --------- | ------------------------------ |
| `Type` | `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (enum) | `int32` | Event type (always OnTrade) - **ENUM!** |
| `EventData` | `OnTadeEventData` | `*OnTadeEventData` | Trade event data (positions, orders, deals) |
| `AccountInfo` | `OnEventAccountInfo` | `*OnEventAccountInfo` | Account snapshot (balance, equity, etc.) |
| `TerminalInstanceGuidId` | `string` | `string` | Terminal instance ID |

---

## 💬 Just the essentials

* **What it is.** Real-time stream of all trading operations.
* **Why you need it.** Monitor order execution, track trading activity, implement notifications.
* **Event-driven.** Notifies on order placement, modification, execution, cancellation.

---

## 🎯 Purpose

Use it to:

* Monitor all trading activity in real-time
* Track order placement, modification, and execution
* Implement trade notifications and alerts
* Log trading operations
* Build trading dashboards
* Detect order fills and rejections immediately

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OnTrade - How it works](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnTrade_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, streams run indefinitely until cancelled.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Channel buffering:** Data channel is buffered (default 100), error channel is buffered (default 10).
* **Goroutine required:** You MUST consume the channels in a separate goroutine to avoid blocking.
* **Context cancellation:** Use `context.WithCancel()` or `context.WithTimeout()` to stop streaming.
* **All events:** Receives events for ALL trading operations (positions, orders, deals).
* **Empty request:** OnTradeRequest is an empty structure (no parameters needed).
* **ENUM type:** Always check the `Type` field (MT5_SUB_ENUM_EVENT_GROUP_TYPE) to identify event type.

---

## 🧱 ENUMs used in OnTrade

### Output ENUM

| ENUM Type | Field Name | Purpose | Values |
|-----------|------------|---------|--------|
| `MT5_SUB_ENUM_EVENT_GROUP_TYPE` | `Type` | Indicates the event type | `OnTrade` (1) - Trade event notification |

**Usage in code:**

```go
tradeStream, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

for event := range tradeStream {
    // Check event type (always OnTrade for this stream)
    switch event.Type {
    case pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OnTrade:
        fmt.Println("Trade event received")
        // Process trade event data
    }
}
```

**Note:** The `Type` field will always be `OnTrade` for OnTrade stream. This ENUM is shared across all streaming methods (OnTrade, OnPositionProfit, OnTradeTransaction) to maintain consistent event identification.

---

## 🔗 Usage Examples

### 1) Basic trade event monitoring

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

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }
                fmt.Printf("[%s] Trade event received\n",
                    time.Now().Format("15:04:05"))

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Stream error: %v\n", err)
                    return
                }
            }
        }
    }()

    // Keep running
    time.Sleep(60 * time.Second)
}
```

### 2) Trade notification system

```go
func TradeNotificationSystem(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    fmt.Println("Trade stream closed")
                    return
                }

                fmt.Printf("\n🔔 Trade Event Notification\n")
                fmt.Printf("  Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
                // Display event details based on event type

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("❌ Stream error: %v\n", err)
                    return
                }

            case <-ctx.Done():
                fmt.Println("Notification system stopped")
                return
            }
        }
    }()

    // Run indefinitely
    select {}
}
```

### 3) Trade logger with file output

```go
func TradeLogger(account *mt5.MT5Account, logFile string) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Failed to open log file: %v\n", err)
        return
    }
    defer file.Close()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }

                logEntry := fmt.Sprintf("[%s] Trade event\n",
                    time.Now().Format("2006-01-02 15:04:05"))
                file.WriteString(logEntry)
                fmt.Print(logEntry)

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

    // Run for extended period
    time.Sleep(24 * time.Hour)
}
```

### 4) Trade event counter

```go
type TradeEventStats struct {
    TotalEvents int64
    StartTime   time.Time
}

func CountTradeEvents(account *mt5.MT5Account, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    stats := &TradeEventStats{
        StartTime: time.Now(),
    }

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    for {
        select {
        case event := <-dataChan:
            if event == nil {
                break
            }
            stats.TotalEvents++
            fmt.Printf("\r Trade events: %d", stats.TotalEvents)

        case err := <-errChan:
            if err != nil {
                fmt.Printf("\nError: %v\n", err)
            }

        case <-ctx.Done():
            elapsed := time.Since(stats.StartTime)
            fmt.Printf("\n\nStatistics:\n")
            fmt.Printf("  Duration: %v\n", elapsed)
            fmt.Printf("  Total events: %d\n", stats.TotalEvents)
            fmt.Printf("  Events/minute: %.2f\n",
                float64(stats.TotalEvents)/(elapsed.Minutes()))
            return
        }
    }
}

// Usage:
// CountTradeEvents(account, 5*time.Minute)
```

### 5) Real-time trade alert system

```go
func RealTimeTradeAlerts(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }

                // Send alert (email, Telegram, etc)
                fmt.Printf("🚨 TRADE ALERT: %s\n", time.Now().Format("15:04:05"))

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Alert system error: %v\n", err)
                    return
                }

            case <-ctx.Done():
                fmt.Println("Alert system shutdown")
                return
            }
        }
    }()

    // Keep running
    select {}
}
```

---

## 📚 See Also

* [OnTradeTransaction](./OnTradeTransaction.md) - Detailed transaction-level events
* [OnPositionsAndPendingOrdersTickets](./OnPositionsAndPendingOrdersTickets.md) - Stream ticket changes
* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get current positions
