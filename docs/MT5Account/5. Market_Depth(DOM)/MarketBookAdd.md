# ✅ MarketBookAdd

> **Request:** Subscribe to Depth of Market (DOM / Level II) data for a symbol.

**API Information:**

* **SDK wrapper:** `MT5Account.MarketBookAdd(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `MarketBookAdd` (defined in `mt5-term-api-market-info.proto`)

### RPC

```go
// MarketBookAdd subscribes to Depth of Market (DOM) updates for a symbol.
//
// Use this method to start receiving Level 2 market data with bid/ask prices
// and volumes at different price levels.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: MarketBookAddRequest with Symbol name
//
// Returns MarketBookAddData with subscription status.
func (a *MT5Account) MarketBookAdd(
    ctx context.Context,
    req *pb.MarketBookAddRequest,
) (*pb.MarketBookAddData, error)
```

**Request message:** `MarketBookAddRequest`

**Reply message:** `MarketBookAddReply { oneof response { MarketBookAddData data = 1; Error error = 2; } }`

---

## 🔽 Input

| Parameter | Type                        | Description                                             |
| --------- | --------------------------- | ------------------------------------------------------- |
| `ctx`     | `context.Context`           | Context for deadline/timeout and cancellation           |
| `req`     | `*pb.MarketBookAddRequest`  | Request with Symbol name                                |

**MarketBookAddRequest fields:**

| Field    | Type     | Description                                   |
| -------- | -------- | --------------------------------------------- |
| `Symbol` | `string` | Trading instrument name (e.g., "EURUSD")      |

---

## ⬆️ Output — `MarketBookAddData`

| Field     | Type   | Description                                              |
| --------- | ------ | -------------------------------------------------------- |
| `Success` | `bool` | `true` if subscription successful, `false` otherwise     |

---

## 💬 Just the essentials

* **What it is.** Subscribes to Level II market data (order book depth) for a specific symbol.
* **Why you need it.** Access detailed market depth showing buy/sell orders at various price levels for advanced trading analysis.
* **Pre-requisite.** Must call this before using `MarketBookGet` to retrieve order book snapshots.

---

## 🎯 Purpose

Use it to:

* Enable Level II market data for a symbol
* Prepare for order book analysis and scalping strategies
* Monitor market liquidity at different price levels
* Analyze support/resistance levels based on order volume
* Implement market microstructure analysis

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [MarketBookAdd - How it works](../HOW_IT_WORK/5.%20Market_Depth(DOM)_HOW/MarketBookAdd_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `5s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Broker support:** Not all brokers provide Level II data. Check with your broker first.
* **Resource cleanup:** Always call `MarketBookRelease` when done to free server resources.
* **Subscription persists:** Subscription remains active until `MarketBookRelease` is called or connection is closed.

---

## 🔗 Usage Examples

### 1) Basic subscription

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
    account := mt5.NewMT5Account("user", "password", "server:443")

    err := account.Connect()
    if err != nil {
        panic(err)
    }
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Subscribe to market depth for EURUSD
    result, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: "EURUSD",
    })
    if err != nil {
        panic(err)
    }

    if result.Success {
        fmt.Println("Successfully subscribed to EURUSD market depth")
    } else {
        fmt.Println("Failed to subscribe to market depth")
    }
}
```

### 2) Subscribe and retrieve market book

```go
func SubscribeAndGetMarketBook(account *mt5.MT5Account, symbol string) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // First, subscribe
    subResult, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: symbol,
    })
    if err != nil {
        fmt.Printf("Subscription error: %v\n", err)
        return
    }

    if !subResult.Success {
        fmt.Println("Subscription failed")
        return
    }

    fmt.Printf("Subscribed to %s market depth\n", symbol)

    // Now retrieve the order book
    bookData, err := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{
        Symbol: symbol,
    })
    if err != nil {
        fmt.Printf("MarketBookGet error: %v\n", err)
        return
    }

    fmt.Printf("Order book has %d levels\n", len(bookData.Book))
    for i, entry := range bookData.Book {
        bookType := "SELL"
        if entry.Type == pb.BookType_BOOK_TYPE_BUY {
            bookType = "BUY"
        }
        fmt.Printf("  Level %d: %s @ %.5f (Volume: %.2f)\n",
            i+1, bookType, entry.Price, entry.VolumeDouble)
    }

    // Clean up
    account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
        Symbol: symbol,
    })
}
```

### 3) Subscribe to multiple symbols

```go
func SubscribeToMultipleSymbols(account *mt5.MT5Account, symbols []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    successCount := 0

    for _, symbol := range symbols {
        result, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
            Symbol: symbol,
        })
        if err != nil {
            fmt.Printf("%s: Error - %v\n", symbol, err)
            continue
        }

        if result.Success {
            fmt.Printf("%s: Subscribed\n", symbol)
            successCount++
        } else {
            fmt.Printf("%s: Failed\n", symbol)
        }
    }

    fmt.Printf("\nSuccessfully subscribed to %d out of %d symbols\n",
        successCount, len(symbols))
}

// Usage:
// SubscribeToMultipleSymbols(account, []string{"EURUSD", "GBPUSD", "USDJPY"})
```

### 4) Subscription with error handling

```go
func SafeMarketBookSubscribe(account *mt5.MT5Account, symbol string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: symbol,
    })

    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("subscription timeout for %s", symbol)
        }
        return fmt.Errorf("subscription failed: %w", err)
    }

    if !result.Success {
        return fmt.Errorf("broker rejected subscription for %s (may not support Level II data)", symbol)
    }

    fmt.Printf("Successfully subscribed to %s market depth\n", symbol)
    return nil
}
```

### 5) Subscribe with automatic cleanup

```go
func WithMarketBookSubscription(account *mt5.MT5Account, symbol string, fn func() error) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Subscribe
    result, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: symbol,
    })
    if err != nil {
        return fmt.Errorf("subscription failed: %w", err)
    }

    if !result.Success {
        return fmt.Errorf("subscription rejected for %s", symbol)
    }

    fmt.Printf("Subscribed to %s\n", symbol)

    // Ensure cleanup
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel()

        account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
            Symbol: symbol,
        })
        fmt.Printf("Unsubscribed from %s\n", symbol)
    }()

    // Execute user function
    return fn()
}

// Usage:
// err := WithMarketBookSubscription(account, "EURUSD", func() error {
//     // Your market book analysis code here
//     return nil
// })
```

### 6) Check subscription status

```go
func CheckMarketBookSupport(account *mt5.MT5Account, symbol string) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Try to subscribe
    result, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: symbol,
    })

    if err != nil {
        fmt.Printf("%s: Error checking support - %v\n", symbol, err)
        return false
    }

    if !result.Success {
        fmt.Printf("%s: Level II data not supported by broker\n", symbol)
        return false
    }

    fmt.Printf("%s: Level II data available\n", symbol)

    // Clean up immediately
    account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
        Symbol: symbol,
    })

    return true
}
```

---

## 🔧 Common Patterns

### Subscribe and monitor

```go
func MonitorMarketDepth(account *mt5.MT5Account, symbol string, duration time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Subscribe
    result, _ := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: symbol,
    })

    if !result.Success {
        fmt.Println("Subscription failed")
        return
    }

    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel()
        account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
            Symbol: symbol,
        })
    }()

    // Monitor for specified duration
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    timeout := time.After(duration)

    for {
        select {
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
            bookData, err := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{
                Symbol: symbol,
            })
            cancel()

            if err == nil && len(bookData.Book) > 0 {
                fmt.Printf("[%s] Best bid: %.5f, Best ask: %.5f\n",
                    time.Now().Format("15:04:05"),
                    bookData.Book[0].Price,
                    bookData.Book[len(bookData.Book)-1].Price)
            }

        case <-timeout:
            return
        }
    }
}
```

---

## 📚 See Also

* [MarketBookRelease](./MarketBookRelease.md) - Unsubscribe from market depth updates
* [MarketBookGet](./MarketBookGet.md) - Retrieve current market depth snapshot
