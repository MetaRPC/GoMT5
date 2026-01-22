# ✅ Stream Real-Time Symbol Ticks

> **Request:** subscribe to real-time tick updates for a symbol. Receive Bid/Ask price updates through Go channels.

**API Information:**

* **Low-level API:** `MT5Account.OnSymbolTick(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.SubscriptionService`
* **Proto definition:** `OnSymbolTick` (defined in `mt5-term-api-subscription.proto`)

### RPC

* **Service:** `mt5_term_api.SubscriptionService`
* **Method:** `OnSymbolTick(OnSymbolTickRequest) → stream OnSymbolTickReply`
* **Low‑level client (generated):** `SubscriptionServiceClient.OnSymbolTick(ctx, request, opts...)`

```go
package mt5

type MT5Account struct {
    // ...
}

// OnSymbolTick streams real-time tick data for a symbol.
// Returns two channels: data channel and error channel.
func (a *MT5Account) OnSymbolTick(
    ctx context.Context,
    req *pb.OnSymbolTickRequest,
) (<-chan *pb.OnSymbolTickData, <-chan error)
```

**Request message:**

```protobuf
OnSymbolTickRequest {
  string Symbol = 1;  // Symbol name to stream
}
```

---

## 🔽 Input

| Parameter | Type                        | Description                                   |
| --------- | --------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`           | Context for cancellation (cancel to stop stream) |
| `req`     | `*pb.OnSymbolTickRequest`   | Request with Symbol name                      |

---

## ⬆️ Output — Channels

| Channel     | Type                           | Description                              |
| ----------- | ------------------------------ | ---------------------------------------- |
| Data Channel| `<-chan *pb.OnSymbolTickData`  | Receives tick updates                    |
| Error Channel| `<-chan error`                | Receives errors (closed on ctx cancel)   |

**OnSymbolTickData fields:**

| Field    | Type                        | Go Type                  | Description                        |
| -------- | --------------------------- | ------------------------ | ---------------------------------- |
| `Bid`    | `double`                    | `float64`                | Current bid price                  |
| `Ask`    | `double`                    | `float64`                | Current ask price                  |
| `Last`   | `double`                    | `float64`                | Last deal price                    |
| `Volume` | `uint64`                    | `uint64`                 | Volume                             |
| `Time`   | `google.protobuf.Timestamp` | `*timestamppb.Timestamp` | Tick timestamp                     |
| `Flags`  | `uint32`                    | `uint32`                 | Tick flags (see MT5 documentation) |

---

## 💬 Just the essentials

* **What it is.** Real-time streaming of price updates via Go channels.
* **Why you need it.** Monitor live prices, implement tick-based strategies, real-time dashboards.
* **Goroutines required.** Must consume channels in goroutines to avoid blocking.

---

## 🎯 Purpose

Use it to:

* Stream real-time tick data for symbols
* Monitor live bid/ask prices and spreads
* Implement tick-based trading strategies
* Build real-time price dashboards
* Track price movements and volatility
* Set up price alerts and notifications

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OnSymbolTick - How it works](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnSymbolTick_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, streams run indefinitely until cancelled.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Channel buffering:** Data channel is buffered (default 100), error channel is buffered (default 10).
* **Goroutine required:** You MUST consume the channels in a separate goroutine to avoid blocking.
* **Context cancellation:** Use `context.WithCancel()` or `context.WithTimeout()` to stop streaming.
* **Channel closure:** Both channels close when context is cancelled or stream ends.
* **Symbol must exist:** Use `SymbolSelect` to add symbol to Market Watch before streaming.

---

## 🔗 Usage Examples

### 1) Basic streaming with goroutine

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

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
        Symbol: "EURUSD",
    })

    // Process in goroutine
    go func() {
        for {
            select {
            case tick := <-dataChan:
                if tick == nil {
                    return
                }
                fmt.Printf("[%s] EURUSD: Bid=%.5f, Ask=%.5f, Spread=%.5f\n",
                    time.Now().Format("15:04:05"),
                    tick.Bid, tick.Ask, tick.Ask-tick.Bid)

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Stream error: %v\n", err)
                    return
                }
            }
        }
    }()

    // Keep main running
    <-ctx.Done()
    fmt.Println("Stream stopped")
}
```

### 2) Streaming with context cancellation

```go
func StreamWithCancellation(account *mt5.MT5Account) {
    ctx, cancel := context.WithCancel(context.Background())

    dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
        Symbol: "EURUSD",
    })

    // Start consuming
    go func() {
        for {
            select {
            case tick := <-dataChan:
                if tick == nil {
                    fmt.Println("Data channel closed")
                    return
                }
                fmt.Printf("Tick: Bid=%.5f, Ask=%.5f\n", tick.Bid, tick.Ask)

            case err := <-errChan:
                if err != nil {
                    fmt.Printf("Error: %v\n", err)
                    return
                }

            case <-ctx.Done():
                fmt.Println("Context cancelled")
                return
            }
        }
    }()

    // Cancel after 10 seconds
    time.Sleep(10 * time.Second)
    cancel()
    fmt.Println("Cancellation triggered")
}
```

### 3) Multiple symbol streaming

```go
func StreamMultipleSymbols(account *mt5.MT5Account, symbols []string) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    for _, symbol := range symbols {
        dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
            Symbol: symbol,
        })

        // Launch goroutine for each symbol
        go func(sym string, data <-chan *pb.OnSymbolTickData, errs <-chan error) {
            for {
                select {
                case tick := <-data:
                    if tick == nil {
                        return
                    }
                    fmt.Printf("[%s] Bid=%.5f, Ask=%.5f\n", sym, tick.Bid, tick.Ask)

                case err := <-errs:
                    if err != nil {
                        fmt.Printf("[%s] Error: %v\n", sym, err)
                        return
                    }

                case <-ctx.Done():
                    return
                }
            }
        }(symbol, dataChan, errChan)
    }

    // Run for 30 seconds
    time.Sleep(30 * time.Second)
}

// Usage:
// StreamMultipleSymbols(account, []string{"EURUSD", "GBPUSD", "USDJPY"})
```

### 4) Tick aggregation and statistics

```go
type TickStats struct {
    Symbol    string
    Count     int64
    LastBid   float64
    LastAsk   float64
    MinBid    float64
    MaxBid    float64
}

func AggregateTickStats(account *mt5.MT5Account, symbol string, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    stats := &TickStats{
        Symbol: symbol,
        MinBid: 999999.0,
    }

    dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
        Symbol: symbol,
    })

    for {
        select {
        case tick := <-dataChan:
            if tick == nil {
                break
            }
            stats.Count++
            stats.LastBid = tick.Bid
            stats.LastAsk = tick.Ask

            if tick.Bid < stats.MinBid {
                stats.MinBid = tick.Bid
            }
            if tick.Bid > stats.MaxBid {
                stats.MaxBid = tick.Bid
            }

        case err := <-errChan:
            if err != nil {
                fmt.Printf("Error: %v\n", err)
            }

        case <-ctx.Done():
            fmt.Printf("\n%s Statistics:\n", stats.Symbol)
            fmt.Printf("  Ticks received: %d\n", stats.Count)
            fmt.Printf("  Last Bid/Ask: %.5f / %.5f\n", stats.LastBid, stats.LastAsk)
            fmt.Printf("  Bid range: %.5f - %.5f\n", stats.MinBid, stats.MaxBid)
            return
        }
    }
}

// Usage:
// AggregateTickStats(account, "EURUSD", 60*time.Second)
```

### 5) Price alert system

```go
func PriceAlertMonitor(account *mt5.MT5Account, symbol string, targetPrice float64, above bool) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
        Symbol: symbol,
    })

    fmt.Printf("Monitoring %s for price %s %.5f\n",
        symbol,
        map[bool]string{true: "above", false: "below"}[above],
        targetPrice)

    go func() {
        for {
            select {
            case tick := <-dataChan:
                if tick == nil {
                    return
                }

                triggered := false
                if above && tick.Ask >= targetPrice {
                    triggered = true
                } else if !above && tick.Bid <= targetPrice {
                    triggered = true
                }

                if triggered {
                    fmt.Printf("\n🔔 ALERT: %s price reached %.5f (Bid=%.5f, Ask=%.5f)\n",
                        symbol, targetPrice, tick.Bid, tick.Ask)
                    cancel() // Stop monitoring
                    return
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

    <-ctx.Done()
}

// Usage:
// PriceAlertMonitor(account, "EURUSD", 1.11000, true) // Alert when above 1.11000
```

### 6) Real-time spread monitoring

```go
func MonitorSpread(account *mt5.MT5Account, symbol string, maxSpreadPips float64) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
        Symbol: symbol,
    })

    go func() {
        for {
            select {
            case tick := <-dataChan:
                if tick == nil {
                    return
                }

                spread := tick.Ask - tick.Bid
                spreadPips := spread * 10000 // For 5-digit quotes

                if spreadPips > maxSpreadPips {
                    fmt.Printf("⚠️  High spread alert: %.1f pips (Bid=%.5f, Ask=%.5f)\n",
                        spreadPips, tick.Bid, tick.Ask)
                } else {
                    fmt.Printf("[%s] Spread: %.1f pips\n",
                        time.Now().Format("15:04:05"), spreadPips)
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

    time.Sleep(60 * time.Second)
    cancel()
}

// Usage:
// MonitorSpread(account, "EURUSD", 2.0) // Alert if spread > 2 pips
```

---

## 🔧 Common Patterns

### Stream with graceful shutdown

```go
func StreamWithShutdown(account *mt5.MT5Account, symbol string) {
    ctx, cancel := context.WithCancel(context.Background())

    dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
        Symbol: symbol,
    })

    done := make(chan struct{})

    go func() {
        defer close(done)
        for {
            select {
            case tick := <-dataChan:
                if tick == nil {
                    return
                }
                // Process tick...

            case err := <-errChan:
                if err != nil {
                    return
                }

            case <-ctx.Done():
                fmt.Println("Shutting down gracefully...")
                return
            }
        }
    }()

    // Wait for signal (Ctrl+C, etc)
    // ...then cancel
    cancel()
    <-done // Wait for goroutine to finish
}
```

---

## 📚 See Also

* [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - Get single tick snapshot
* [SymbolSelect](../2.%20Symbol_information/SymbolSelect.md) - Add symbol to Market Watch
* [OnTrade](./OnTrade.md) - Stream trade events
