# ✅ MarketBookRelease

> **Request:** Unsubscribe from Depth of Market (DOM / Level II) data and free server resources.

**API Information:**

* **SDK wrapper:** `MT5Account.MarketBookRelease(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `MarketBookRelease` (defined in `mt5-term-api-market-info.proto`)

### RPC

```go
// MarketBookRelease unsubscribes from Depth of Market (DOM) updates.
//
// Use this method to stop receiving Level 2 market data and free resources.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: MarketBookReleaseRequest with Symbol name
//
// Returns MarketBookReleaseData with unsubscription status.
func (a *MT5Account) MarketBookRelease(
    ctx context.Context,
    req *pb.MarketBookReleaseRequest,
) (*pb.MarketBookReleaseData, error)
```

**Request message:** `MarketBookReleaseRequest`

**Reply message:** `MarketBookReleaseReply { oneof response { MarketBookReleaseData data = 1; Error error = 2; } }`

---

## 🔽 Input

| Parameter | Type                           | Description                                             |
| --------- | ------------------------------ | ------------------------------------------------------- |
| `ctx`     | `context.Context`              | Context for deadline/timeout and cancellation           |
| `req`     | `*pb.MarketBookReleaseRequest` | Request with Symbol name                                |

**MarketBookReleaseRequest fields:**

| Field    | Type     | Description                                   |
| -------- | -------- | --------------------------------------------- |
| `Symbol` | `string` | Trading instrument name (e.g., "EURUSD")      |

---

## ⬆️ Output — `MarketBookReleaseData`

| Field     | Type   | Description                                                    |
| --------- | ------ | -------------------------------------------------------------- |
| `Success` | `bool` | `true` if unsubscribed successfully, `false` otherwise         |

---

## 💬 Just the essentials

* **What it is.** Cancels subscription to Level II market data for a symbol.
* **Why you need it.** Free server resources and reduce network traffic when market depth data is no longer needed.
* **Clean up.** Always call this after `MarketBookAdd` to prevent resource leaks.

---

## 🎯 Purpose

Use it to:

* Stop receiving market depth updates
* Free server and client resources
* Clean up after market depth analysis
* Implement proper resource management
* Reduce unnecessary network traffic

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [MarketBookRelease - How it works](../HOW_IT_WORK/5.%20Market_Depth(DOM)_HOW/MarketBookRelease_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `5s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Always cleanup:** Call this in `defer` after successful `MarketBookAdd` to ensure cleanup.
* **No error on double release:** Calling this multiple times for same symbol is safe (no-op).

---

## 🔗 Usage Examples

### 1) Basic unsubscribe

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
    account := mt5.NewMT5Account("user", "password", "server:443")

    err := account.Connect()
    if err != nil {
        panic(err)
    }
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Unsubscribe from market depth
    result, err := account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
        Symbol: "EURUSD",
    })
    if err != nil {
        panic(err)
    }

    if result.Success {
        fmt.Println("Successfully unsubscribed from EURUSD market depth")
    }
}
```

### 2) Subscribe with defer cleanup

```go
func AnalyzeMarketDepth(account *mt5.MT5Account, symbol string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Subscribe
    subResult, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: symbol,
    })
    if err != nil {
        return err
    }

    if !subResult.Success {
        return fmt.Errorf("subscription failed")
    }

    // Ensure cleanup with defer
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel()

        result, err := account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
            Symbol: symbol,
        })
        if err != nil {
            fmt.Printf("Warning: cleanup error - %v\n", err)
        } else if result.Success {
            fmt.Printf("Cleaned up %s subscription\n", symbol)
        }
    }()

    // Your market depth analysis here
    bookData, err := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{
        Symbol: symbol,
    })
    if err != nil {
        return err
    }

    fmt.Printf("Analyzed %d order book levels\n", len(bookData.Book))
    return nil
}
```

### 3) Release multiple symbols

```go
func ReleaseMultipleSymbols(account *mt5.MT5Account, symbols []string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    for _, symbol := range symbols {
        result, err := account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
            Symbol: symbol,
        })
        if err != nil {
            fmt.Printf("%s: Error - %v\n", symbol, err)
            continue
        }

        if result.Success {
            fmt.Printf("%s: Unsubscribed\n", symbol)
        } else {
            fmt.Printf("%s: Not subscribed\n", symbol)
        }
    }
}

// Usage:
// ReleaseMultipleSymbols(account, []string{"EURUSD", "GBPUSD", "USDJPY"})
```

### 4) Safe cleanup pattern

```go
func SafeMarketBookCleanup(account *mt5.MT5Account, symbol string) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
        Symbol: symbol,
    })

    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Printf("%s: Cleanup timeout (may already be released)\n", symbol)
        } else {
            fmt.Printf("%s: Cleanup error - %v\n", symbol, err)
        }
        return
    }

    if result.Success {
        fmt.Printf("%s: Successfully released\n", symbol)
    } else {
        fmt.Printf("%s: Was not subscribed\n", symbol)
    }
}
```

---

## 🔧 Common Patterns

### Resource management wrapper

```go
type MarketBookSession struct {
    account *mt5.MT5Account
    symbol  string
    active  bool
}

func NewMarketBookSession(account *mt5.MT5Account, symbol string) (*MarketBookSession, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
        Symbol: symbol,
    })
    if err != nil {
        return nil, err
    }

    if !result.Success {
        return nil, fmt.Errorf("subscription failed for %s", symbol)
    }

    return &MarketBookSession{
        account: account,
        symbol:  symbol,
        active:  true,
    }, nil
}

func (s *MarketBookSession) Close() error {
    if !s.active {
        return nil
    }

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    result, err := s.account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
        Symbol: s.symbol,
    })
    if err != nil {
        return err
    }

    s.active = false

    if !result.Success {
        return fmt.Errorf("unsubscribe failed")
    }

    return nil
}

// Usage:
// session, err := NewMarketBookSession(account, "EURUSD")
// if err != nil {
//     panic(err)
// }
// defer session.Close()
```

---

## 📚 See Also

* [MarketBookAdd](./MarketBookAdd.md) - Subscribe to market depth updates
* [MarketBookGet](./MarketBookGet.md) - Retrieve current market depth snapshot
