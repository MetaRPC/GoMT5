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

**OnTradeTransactionData structure:**

| Field                   | Type                             | Description                                    |
| ----------------------- | -------------------------------- | ---------------------------------------------- |
| `Type`                  | `MT5_SUB_ENUM_EVENT_GROUP_TYPE`  | Event group type (always TradeTransaction)     |
| `TradeTransaction`      | `*MqlTradeTransaction`           | Transaction details (order, deal info)         |
| `TradeRequest`          | `*MqlTradeRequest`               | Original trade request (may be nil)            |
| `TradeResult`           | `*MqlTradeResult`                | Trade operation result (may be nil)            |
| `TerminalInstanceGuidId`| `string`                         | Terminal instance ID                           |
| `AccountInfo`           | `*OnEventAccountInfo`            | Account information snapshot                   |

**MqlTradeTransaction - key fields:**

| Field                   | Type                                | Description                        |
| ----------------------- | ----------------------------------- | ---------------------------------- |
| `Type`                  | `SUB_ENUM_TRADE_TRANSACTION_TYPE`   | Transaction type (order add/remove, deal add, etc.) |
| `OrderType`             | `SUB_ENUM_ORDER_TYPE`               | Order type (BUY, SELL, BUY_LIMIT, etc.) |
| `OrderState`            | `SUB_ENUM_ORDER_STATE`              | Order state (STARTED, PLACED, FILLED, etc.) |
| `DealType`              | `SUB_ENUM_DEAL_TYPE`                | Deal type (BUY, SELL, BALANCE, etc.) |
| `OrderTimeType`         | `SUB_ENUM_ORDER_TYPE_TIME`          | Order lifetime (GTC, DAY, SPECIFIED, etc.) |
| `DealTicket`            | `uint64`                            | Deal ticket number                 |
| `OrderTicket`           | `uint64`                            | Order ticket number                |
| `Symbol`                | `string`                            | Trading symbol                     |
| `Price`                 | `float64`                           | Transaction price                  |
| `Volume`                | `float64`                           | Transaction volume                 |
| `PriceStopLoss`         | `float64`                           | Stop Loss level                    |
| `PriceTakeProfit`       | `float64`                           | Take Profit level                  |

**MqlTradeRequest - key fields:**

| Field                   | Type                                | Description                        |
| ----------------------- | ----------------------------------- | ---------------------------------- |
| `TradeOperationType`    | `SUB_ENUM_TRADE_REQUEST_ACTIONS`    | Trade action (DEAL, PENDING, SLTP, etc.) |
| `OrderType`             | `SUB_ENUM_ORDER_TYPE`               | Order type (BUY, SELL, BUY_LIMIT, etc.) |
| `OrderTypeFilling`      | `SUB_ENUM_ORDER_TYPE_FILLING`       | Filling mode (FOK, IOC, RETURN, etc.) |
| `TypeTime`              | `SUB_ENUM_ORDER_TYPE_TIME`          | Order lifetime (GTC, DAY, SPECIFIED, etc.) |
| `Symbol`                | `string`                            | Trade symbol                       |
| `RequestedDealVolumeLots`| `float64`                          | Requested volume in lots           |
| `Price`                 | `float64`                           | Price                              |
| `StopLoss`              | `float64`                           | Stop Loss level                    |
| `TakeProfit`            | `float64`                           | Take Profit level                  |

**MqlTradeResult - key fields:**

| Field                   | Type                                | Description                        |
| ----------------------- | ----------------------------------- | ---------------------------------- |
| `TradeReturnCode`       | `MqlErrorTradeCode`                 | Operation return code (DONE, REJECT, etc.) |
| `DealTicket`            | `uint64`                            | Deal ticket if executed            |
| `OrderTicket`           | `uint64`                            | Order ticket if placed             |
| `DealVolume`            | `float64`                           | Actual deal volume                 |
| `DealPrice`             | `float64`                           | Actual deal price                  |
| `BrokerCommentToOperation`| `string`                          | Broker comment                     |

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

### 1) Basic transaction monitoring with ENUM checks

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

                // Check event type (MT5_SUB_ENUM_EVENT_GROUP_TYPE)
                if event.Type == pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_TradeTransaction {
                    if tx := event.TradeTransaction; tx != nil {
                        // Check transaction type (SUB_ENUM_TRADE_TRANSACTION_TYPE)
                        var txType string
                        switch tx.Type {
                        case pb.SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_ORDER_ADD:
                            txType = "ORDER_ADD"
                        case pb.SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_ORDER_UPDATE:
                            txType = "ORDER_UPDATE"
                        case pb.SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_ORDER_DELETE:
                            txType = "ORDER_DELETE"
                        case pb.SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_DEAL_ADD:
                            txType = "DEAL_ADD"
                        default:
                            txType = fmt.Sprintf("OTHER(%d)", tx.Type)
                        }

                        fmt.Printf("[%s] %s: Symbol=%s, Price=%.5f, Volume=%.2f\n",
                            time.Now().Format("15:04:05"),
                            txType,
                            tx.Symbol,
                            tx.Price,
                            tx.Volume)
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

### 3) Deal execution notifier with ENUM checks

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

                if tx := event.TradeTransaction; tx != nil {
                    // Check if this is a deal execution (SUB_ENUM_TRADE_TRANSACTION_TYPE)
                    if tx.Type == pb.SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_DEAL_ADD {
                        // Determine deal type (SUB_ENUM_DEAL_TYPE)
                        var dealTypeStr string
                        switch tx.DealType {
                        case pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY:
                            dealTypeStr = "BUY"
                        case pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_SELL:
                            dealTypeStr = "SELL"
                        case pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BALANCE:
                            dealTypeStr = "BALANCE"
                        case pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_CREDIT:
                            dealTypeStr = "CREDIT"
                        default:
                            dealTypeStr = fmt.Sprintf("OTHER(%d)", tx.DealType)
                        }

                        // Determine order type (SUB_ENUM_ORDER_TYPE)
                        var orderTypeStr string
                        switch tx.OrderType {
                        case pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY:
                            orderTypeStr = "BUY"
                        case pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL:
                            orderTypeStr = "SELL"
                        case pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_LIMIT:
                            orderTypeStr = "BUY_LIMIT"
                        case pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_LIMIT:
                            orderTypeStr = "SELL_LIMIT"
                        case pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_STOP:
                            orderTypeStr = "BUY_STOP"
                        case pb.SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_STOP:
                            orderTypeStr = "SELL_STOP"
                        default:
                            orderTypeStr = fmt.Sprintf("OTHER(%d)", tx.OrderType)
                        }

                        fmt.Printf("\n💰 DEAL EXECUTED\n")
                        fmt.Printf("  Deal Type: %s\n", dealTypeStr)
                        fmt.Printf("  Order Type: %s\n", orderTypeStr)
                        fmt.Printf("  Deal: %d\n", tx.DealTicket)
                        fmt.Printf("  Order: %d\n", tx.OrderTicket)
                        fmt.Printf("  Symbol: %s\n", tx.Symbol)
                        fmt.Printf("  Price: %.5f\n", tx.Price)
                        fmt.Printf("  Volume: %.2f\n", tx.Volume)
                        fmt.Printf("  Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
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

### 5) Order state change monitor with ENUMs

```go
func MonitorOrderStateChanges(account *mt5.MT5Account) {
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

                if tx := event.TradeTransaction; tx != nil {
                    if tx.OrderTicket > 0 {
                        // Check order state using ENUM (SUB_ENUM_ORDER_STATE)
                        var stateName string
                        switch tx.OrderState {
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_STARTED:
                            stateName = "STARTED"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_PLACED:
                            stateName = "PLACED"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_CANCELED:
                            stateName = "CANCELED"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_PARTIAL:
                            stateName = "PARTIAL"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_FILLED:
                            stateName = "FILLED"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REJECTED:
                            stateName = "REJECTED"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_EXPIRED:
                            stateName = "EXPIRED"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_ADD:
                            stateName = "REQUEST_ADD"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_MODIFY:
                            stateName = "REQUEST_MODIFY"
                        case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_CANCEL:
                            stateName = "REQUEST_CANCEL"
                        default:
                            stateName = fmt.Sprintf("UNKNOWN(%d)", tx.OrderState)
                        }

                        // Check order time type (SUB_ENUM_ORDER_TYPE_TIME)
                        var timeType string
                        switch tx.OrderTimeType {
                        case pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_GTC:
                            timeType = "GTC"
                        case pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_DAY:
                            timeType = "DAY"
                        case pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED:
                            timeType = "SPECIFIED"
                        case pb.SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED_DAY:
                            timeType = "SPECIFIED_DAY"
                        default:
                            timeType = fmt.Sprintf("UNKNOWN(%d)", tx.OrderTimeType)
                        }

                        fmt.Printf("[%s] Order %d: State=%s, TimeType=%s\n",
                            time.Now().Format("15:04:05"),
                            tx.OrderTicket,
                            stateName,
                            timeType)
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

    select {}
}
```

---

## 📚 See Also

* [OnTrade](./OnTrade.md) - Stream high-level trade events
* [OnPositionsAndPendingOrdersTickets](./OnPositionsAndPendingOrdersTickets.md) - Stream ticket changes
* [OrderHistory](../3.%20Position_Orders_Information/OrderHistory.md) - Get historical orders
