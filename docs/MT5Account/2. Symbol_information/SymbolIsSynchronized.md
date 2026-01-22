# ✅ Check Symbol Synchronization with Server

> **Request:** check if symbol data is synchronized with trading server. Ensures quote freshness before trading operations.

**API Information:**

* **Low-level API:** `MT5Account.SymbolIsSynchronized(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolIsSynchronized` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolIsSynchronized(SymbolIsSynchronizedRequest) → SymbolIsSynchronizedReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolIsSynchronized(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Checks if symbol quotes are synchronized with the broker's trade server.
* **Why you need it.** Ensure quotes are up-to-date before placing orders or making trading decisions.
* **Trading prerequisite.** Unsynchronized symbols may have stale prices and cause order rejections.

---

## 🎯 Purpose

Use it to:

* Validate symbol data before trading operations
* Detect connection issues to trade server
* Implement trading safety checks
* Monitor symbol quote health

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolIsSynchronized - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolIsSynchronized_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolIsSynchronized checks if symbol data is synchronized with trade server.
// Returns Synchronized flag indicating sync status.
func (a *MT5Account) SymbolIsSynchronized(
    ctx context.Context,
    req *pb.SymbolIsSynchronizedRequest,
) (*pb.SymbolIsSynchronizedData, error)
```

**Request message:**

```protobuf
SymbolIsSynchronizedRequest {
  string Symbol = 1;  // Symbol name to check
}
```

**Reply message:**

```protobuf
SymbolIsSynchronizedReply {
  oneof response {
    SymbolIsSynchronizedData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                                | Description                                   |
| --------- | ----------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                   | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolIsSynchronizedRequest`   | Request with Symbol name                      |

**Request fields:**

| Field    | Type     | Description                     |
| -------- | -------- | ------------------------------- |
| `Symbol` | `string` | Symbol name (e.g., "EURUSD")    |

---

## ⬆️ Output — `SymbolIsSynchronizedData`

| Field          | Type   | Go Type | Description                                       |
| -------------- | ------ | ------- | ------------------------------------------------- |
| `Synchronized` | `bool` | `bool`  | True if symbol data is synced with trade server   |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Market Watch required:** Symbol should be in Market Watch for synchronization to occur.
* **Weekend behavior:** During market close, symbols may not be synchronized.

---

## 🔗 Usage Examples

### 1) Check symbol synchronization

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

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    data, err := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
        Symbol: "EURUSD",
    })
    if err != nil {
        panic(err)
    }

    if data.Synchronized {
        fmt.Println("EURUSD is synchronized and ready for trading")
    } else {
        fmt.Println("EURUSD is NOT synchronized - quotes may be stale")
    }
}
```

### 2) Wait for symbol synchronization

```go
func WaitForSync(account *mt5.MT5Account, symbol string, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            data, err := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
                Symbol: symbol,
            })
            if err != nil {
                return err
            }

            if data.Synchronized {
                fmt.Printf("%s synchronized\n", symbol)
                return nil
            }

        case <-ctx.Done():
            return fmt.Errorf("timeout waiting for %s synchronization", symbol)
        }
    }
}

// Usage:
// err := WaitForSync(account, "EURUSD", 10*time.Second)
```

### 3) Pre-trade synchronization check

```go
func SafeTrade(account *mt5.MT5Account, symbol string, volume float64) error {
    ctx := context.Background()

    // Check synchronization first
    syncData, err := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
        Symbol: symbol,
    })
    if err != nil {
        return fmt.Errorf("failed to check sync: %w", err)
    }

    if !syncData.Synchronized {
        return fmt.Errorf("symbol %s not synchronized - trading aborted", symbol)
    }

    // Proceed with trading
    fmt.Printf("Symbol %s synchronized, proceeding with order...\n", symbol)
    return nil
}
```

### 4) Check multiple symbols synchronization

```go
func CheckMultipleSymbols(account *mt5.MT5Account, symbols []string) map[string]bool {
    ctx := context.Background()
    results := make(map[string]bool)

    for _, symbol := range symbols {
        data, err := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
            Symbol: symbol,
        })

        if err != nil {
            results[symbol] = false
            continue
        }

        results[symbol] = data.Synchronized
    }

    return results
}

// Usage:
// symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
// syncStatus := CheckMultipleSymbols(account, symbols)
// for symbol, synced := range syncStatus {
//     fmt.Printf("%s: %v\n", symbol, synced)
// }
```

### 5) Monitor symbol synchronization

```go
func MonitorSynchronization(account *mt5.MT5Account, symbol string, interval time.Duration) {
    ctx := context.Background()
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    lastStatus := false

    for range ticker.C {
        data, err := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
            Symbol: symbol,
        })

        if err != nil {
            fmt.Printf("[%s] Error checking %s: %v\n", time.Now().Format("15:04:05"), symbol, err)
            continue
        }

        // Report status changes
        if data.Synchronized != lastStatus {
            status := "synchronized"
            if !data.Synchronized {
                status = "NOT synchronized"
            }
            fmt.Printf("[%s] %s is now %s\n", time.Now().Format("15:04:05"), symbol, status)
            lastStatus = data.Synchronized
        }
    }
}

// Usage in goroutine:
// go MonitorSynchronization(account, "EURUSD", 5*time.Second)
```

---

## 🔧 Common Patterns

### Synchronization guard

```go
func EnsureSynchronized(account *mt5.MT5Account, symbol string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
        Symbol: symbol,
    })

    if err != nil {
        return fmt.Errorf("sync check failed: %w", err)
    }

    if !data.Synchronized {
        return fmt.Errorf("symbol %s not synchronized", symbol)
    }

    return nil
}
```

### Retry with synchronization check

```go
func ExecuteWhenSynced(account *mt5.MT5Account, symbol string, fn func() error) error {
    maxRetries := 5
    retryDelay := 1 * time.Second

    for i := 0; i < maxRetries; i++ {
        ctx := context.Background()
        data, err := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
            Symbol: symbol,
        })

        if err == nil && data.Synchronized {
            return fn()
        }

        if i < maxRetries-1 {
            time.Sleep(retryDelay)
        }
    }

    return fmt.Errorf("symbol %s failed to synchronize", symbol)
}
```

---

## 📚 See Also

* [SymbolSelect](./SymbolSelect.md) - Add symbol to Market Watch (required for sync)
* [SymbolExist](./SymbolExist.md) - Check if symbol exists
* [SymbolInfoTick](./SymbolInfoTick.md) - Get last tick data
* [OrderCheck](../4.%20Trading_Operations/OrderCheck.md) - Validate order before sending
