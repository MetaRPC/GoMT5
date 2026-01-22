# ✅ Stream Detailed Trade Transactions

> **Request:** subscribe to real-time notifications about all trading state changes. Receive detailed low-level events: order state changes, deal execution, position modifications.

**API Information:**

* **Low-level API:** `MT5Account.OnTradeTransaction(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.SubscriptionService`
* **Proto definition:** `OnTradeTransaction` (defined in `mt5-term-api-subscription.proto`)

### RPC

* **Service:** `mt5_term_api.SubscriptionService`
* **Method:** `OnTradeTransaction(OnTradeTransactionRequest) → stream OnTradeTransactionReply`
* **Low‑level client (generated):** `SubscriptionServiceClient.OnTradeTransaction(ctx, request, opts...)`

```go
package mt5

type MT5Account struct {
    // ...
}

// OnTradeTransaction streams detailed trade transaction events.
// Returns two channels: data channel and error channel.
func (a *MT5Account) OnTradeTransaction(
    ctx context.Context,
    req *pb.OnTradeTransactionRequest,
) (<-chan *pb.OnTradeTransactionData, <-chan error)
```

---

## 🔽 Input

| Parameter | Type                              | Description                                   |
| --------- | --------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                 | Context for cancellation (cancel to stop stream) |
| `req`     | `*pb.OnTradeTransactionRequest`   | Request (empty structure)                     |

---

## ⬆️ Output — Channels

| Channel      | Type                                | Description                              |
| ------------ | ----------------------------------- | ---------------------------------------- |
| Data Channel | `<-chan *pb.OnTradeTransactionData` | Receives transaction event updates       |
| Error Channel| `<-chan error`                      | Receives errors (closed on ctx cancel)   |

**OnTradeTransactionData contains MqlTradeTransaction with key fields:**

| Field         | Type     | Go Type   | Description                        |
| ------------- | -------- | --------- | ---------------------------------- |
| `Type`        | `int32`  | `int32`   | Transaction type (order, deal)     |
| `OrderState`  | `int32`  | `int32`   | Order state                        |
| `DealTicket`  | `uint64` | `uint64`  | Deal ticket number                 |
| `OrderTicket` | `uint64` | `uint64`  | Order ticket number                |
| `Symbol`      | `string` | `string`  | Trading symbol                     |
| `Price`       | `double` | `float64` | Transaction price                  |
| `Volume`      | `double` | `float64` | Transaction volume                 |

**Plus many more fields** - see MT5 MqlTradeTransaction documentation for complete structure.

---

## 💬 Just the essentials

* **What it is.** Low-level stream of all trading state changes and transactions.
* **Why you need it.** Detailed trade monitoring, audit trails, transaction logging.
* **Most detailed.** Provides finest-grained trading event information.

---

## 🎯 Purpose

Use it to:

* Monitor all trading transactions at low level
* Build complete trade audit trails
* Track order state changes in detail
* Implement transaction logging
* Detect deal execution immediately
* Analyze trading operations granularly

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OnTradeTransaction - How it works](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnTradeTransaction_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, streams run indefinitely until cancelled.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Channel buffering:** Data channel is buffered (default 100), error channel is buffered (default 10).
* **Goroutine required:** You MUST consume the channels in a separate goroutine to avoid blocking.
* **Most detailed:** Provides low-level transaction details.
* **Empty request:** OnTradeTransactionRequest is an empty structure.
* **High frequency:** Can generate many events during active trading.

---

## 🔗 Usage Examples

### 1) Basic transaction monitoring

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

    dataChan, errChan := account.OnTradeTransaction(ctx, &pb.OnTradeTransactionRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }
                fmt.Printf("[%s] Transaction: Type=%d, Symbol=%s, Price=%.5f\n",
                    time.Now().Format("15:04:05"),
                    event.Transaction.Type,
                    event.Transaction.Symbol,
                    event.Transaction.Price)

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Error: %v\n", err)
                    return
                }
            }
        }
    }()

    time.Sleep(60 * time.Second)
}
```

### 2) Transaction logger with file output

```go
func TransactionLogger(account *mt5.MT5Account, logFile string) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Failed to open log file: %v\n", err)
        return
    }
    defer file.Close()

    dataChan, errChan := account.OnTradeTransaction(ctx, &pb.OnTradeTransactionRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }

                logEntry := fmt.Sprintf("[%s] Type=%d Order=%d Deal=%d Symbol=%s Price=%.5f Vol=%.2f\n",
                    time.Now().Format("2006-01-02 15:04:05"),
                    event.Transaction.Type,
                    event.Transaction.OrderTicket,
                    event.Transaction.DealTicket,
                    event.Transaction.Symbol,
                    event.Transaction.Price,
                    event.Transaction.Volume)

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

    select {}
}

// Usage:
// TransactionLogger(account, "transactions.log")
```

### 3) Deal execution notifier

```go
func DealExecutionNotifier(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTradeTransaction(ctx, &pb.OnTradeTransactionRequest{})

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }

                // Check if this is a deal execution
                if event.Transaction.DealTicket > 0 {
                    fmt.Printf("\n💰 DEAL EXECUTED\n")
                    fmt.Printf("  Deal: %d\n", event.Transaction.DealTicket)
                    fmt.Printf("  Order: %d\n", event.Transaction.OrderTicket)
                    fmt.Printf("  Symbol: %s\n", event.Transaction.Symbol)
                    fmt.Printf("  Price: %.5f\n", event.Transaction.Price)
                    fmt.Printf("  Volume: %.2f\n", event.Transaction.Volume)
                    fmt.Printf("  Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
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
```

### 4) Transaction statistics tracker

```go
type TransactionStats struct {
    TotalTransactions int64
    DealCount         int64
    OrderCount        int64
    SymbolCounts      map[string]int64
    StartTime         time.Time
}

func TrackTransactionStats(account *mt5.MT5Account, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    stats := &TransactionStats{
        SymbolCounts: make(map[string]int64),
        StartTime:    time.Now(),
    }

    dataChan, errChan := account.OnTradeTransaction(ctx, &pb.OnTradeTransactionRequest{})

    for {
        select {
        case event := <-dataChan:
            if event == nil {
                break
            }

            stats.TotalTransactions++

            if event.Transaction.DealTicket > 0 {
                stats.DealCount++
            }
            if event.Transaction.OrderTicket > 0 {
                stats.OrderCount++
            }

            if event.Transaction.Symbol != "" {
                stats.SymbolCounts[event.Transaction.Symbol]++
            }

        case err := <-errChan:
            if err != nil {
                fmt.Printf("Error: %v\n", err)
            }

        case <-ctx.Done():
            elapsed := time.Since(stats.StartTime)

            fmt.Println("\n╔═══════════════════════════════════════╗")
            fmt.Println("║   Transaction Statistics Summary      ║")
            fmt.Println("╠═══════════════════════════════════════╣")
            fmt.Printf("║ Duration:        %-20v ║\n", elapsed.Round(time.Second))
            fmt.Printf("║ Total:           %-20d ║\n", stats.TotalTransactions)
            fmt.Printf("║ Deals:           %-20d ║\n", stats.DealCount)
            fmt.Printf("║ Orders:          %-20d ║\n", stats.OrderCount)
            fmt.Println("╠═══════════════════════════════════════╣")
            fmt.Println("║ Activity by Symbol:                   ║")

            for symbol, count := range stats.SymbolCounts {
                fmt.Printf("║   %-10s   %-20d ║\n", symbol, count)
            }

            fmt.Println("╚═══════════════════════════════════════╝")
            return
        }
    }
}

// Usage:
// TrackTransactionStats(account, 10*time.Minute)
```

### 5) Order state change monitor

```go
func MonitorOrderStateChanges(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnTradeTransaction(ctx, &pb.OnTradeTransactionRequest{})

    orderStates := map[int32]string{
        0: "STARTED",
        1: "PLACED",
        2: "CANCELED",
        3: "PARTIAL",
        4: "FILLED",
        5: "REJECTED",
        6: "EXPIRED",
        7: "REQUEST_ADD",
        8: "REQUEST_MODIFY",
        9: "REQUEST_CANCEL",
    }

    go func() {
        for {
            select {
            case event := <-dataChan:
                if event == nil {
                    return
                }

                if event.Transaction.OrderTicket > 0 {
                    stateName := orderStates[event.Transaction.OrderState]
                    if stateName == "" {
                        stateName = fmt.Sprintf("UNKNOWN(%d)", event.Transaction.OrderState)
                    }

                    fmt.Printf("[%s] Order %d state: %s\n",
                        time.Now().Format("15:04:05"),
                        event.Transaction.OrderTicket,
                        stateName)
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
```

---

## 📚 See Also

* [OnTrade](./OnTrade.md) - Stream high-level trade events
* [OnPositionsAndPendingOrdersTickets](./OnPositionsAndPendingOrdersTickets.md) - Stream ticket changes
* [OrderHistory](../3.%20Position_Orders_Information/OrderHistory.md) - Get historical orders
