# ✅ Get Symbol Last Tick

> **Request:** get last tick (quote) data for trading symbol with timestamp and detailed information about prices and volumes.

**API Information:**

* **Low-level API:** `MT5Account.SymbolInfoTick(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolInfoTick` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolInfoTick(SymbolInfoTickRequest) → SymbolInfoTickReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolInfoTick(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieves the most recent tick (quote) data for a symbol.
* **Why you need it.** Get current market prices with timestamp for trading decisions.
* **Complete data.** Returns Bid, Ask, Last price, volumes, and millisecond-precision timestamp.

---

## 🎯 Purpose

Use it to:

* Get current market prices before placing orders
* Monitor price changes with timestamps
* Calculate spreads and market conditions
* Validate quote freshness
* Build price history snapshots

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolInfoTick - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoTick_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolInfoTick retrieves the last tick data for a symbol.
// Returns MrpcMqlTick with Bid, Ask, Last, Volume, Time, and other tick properties.
func (a *MT5Account) SymbolInfoTick(
    ctx context.Context,
    req *pb.SymbolInfoTickRequest,
) (*pb.MrpcMqlTick, error)
```

**Request message:**

```protobuf
SymbolInfoTickRequest {
  string Symbol = 1;  // Symbol name
}
```

**Reply message:**

```protobuf
SymbolInfoTickReply {
  oneof response {
    MrpcMqlTick data = 1;  // Tick data
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                          | Description                                   |
| --------- | ----------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`             | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolInfoTickRequest`   | Request with Symbol name                      |

**Request fields:**

| Field    | Type     | Description                       |
| -------- | -------- | --------------------------------- |
| `Symbol` | `string` | Symbol name (e.g., "EURUSD")      |

---

## ⬆️ Output — `MrpcMqlTick`

| Field        | Type     | Go Type   | Description                                    |
| ------------ | -------- | --------- | ---------------------------------------------- |
| `Time`       | `int64`  | `int64`   | Last quote time (Unix timestamp seconds)       |
| `Bid`        | `double` | `float64` | Current Bid price                              |
| `Ask`        | `double` | `float64` | Current Ask price                              |
| `Last`       | `double` | `float64` | Last deal price                                |
| `Volume`     | `uint64` | `uint64`  | Volume for the last deal                       |
| `TimeMS`     | `int64`  | `int64`   | Last quote time (Unix timestamp milliseconds)  |
| `Flags`      | `uint32` | `uint32`  | Tick flags (bid/ask/last changed indicators)   |
| `VolumeReal` | `double` | `float64` | Volume with extended accuracy                  |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Symbol synchronization:** Symbol should be synchronized for accurate tick data.
* **Millisecond precision:** TimeMS provides millisecond-level timestamp accuracy.

---

## 🔗 Usage Examples

### 1) Get last tick data

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

    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
        Symbol: "EURUSD",
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("EURUSD Tick:\n")
    fmt.Printf("  Bid: %.5f\n", tick.Bid)
    fmt.Printf("  Ask: %.5f\n", tick.Ask)
    fmt.Printf("  Spread: %.5f\n", tick.Ask-tick.Bid)
    fmt.Printf("  Time: %s\n", time.Unix(tick.Time, 0).Format("2006-01-02 15:04:05"))
}
```

### 2) Get current spread

```go
func GetCurrentSpread(account *mt5.MT5Account, symbol string) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
        Symbol: symbol,
    })
    if err != nil {
        return 0, fmt.Errorf("failed to get tick: %w", err)
    }

    spread := tick.Ask - tick.Bid
    return spread, nil
}

// Usage:
// spread, _ := GetCurrentSpread(account, "EURUSD")
// fmt.Printf("Current spread: %.5f\n", spread)
```

### 3) Check quote freshness

```go
func IsQuoteFresh(account *mt5.MT5Account, symbol string, maxAge time.Duration) (bool, error) {
    ctx := context.Background()

    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
        Symbol: symbol,
    })
    if err != nil {
        return false, err
    }

    tickTime := time.Unix(tick.Time, 0)
    age := time.Since(tickTime)

    return age <= maxAge, nil
}

// Usage:
// fresh, _ := IsQuoteFresh(account, "EURUSD", 5*time.Second)
// if fresh {
//     fmt.Println("Quote is fresh")
// }
```

### 4) Get mid price

```go
func GetMidPrice(account *mt5.MT5Account, symbol string) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
        Symbol: symbol,
    })
    if err != nil {
        return 0, err
    }

    midPrice := (tick.Bid + tick.Ask) / 2.0
    return midPrice, nil
}
```

### 5) Monitor price changes

```go
func MonitorPriceChanges(account *mt5.MT5Account, symbol string, interval time.Duration) {
    ctx := context.Background()
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    var lastBid, lastAsk float64

    for range ticker.C {
        tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
            Symbol: symbol,
        })
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        if lastBid != 0 {
            bidChange := tick.Bid - lastBid
            askChange := tick.Ask - lastAsk
            spread := tick.Ask - tick.Bid

            fmt.Printf("[%s] %s: Bid=%.5f (%+.5f), Ask=%.5f (%+.5f), Spread=%.5f\n",
                time.Unix(tick.Time, 0).Format("15:04:05"),
                symbol,
                tick.Bid, bidChange,
                tick.Ask, askChange,
                spread,
            )
        }

        lastBid = tick.Bid
        lastAsk = tick.Ask
    }
}

// Usage:
// MonitorPriceChanges(account, "EURUSD", 1*time.Second)
```

### 6) Get tick with validation

```go
type ValidatedTick struct {
    Tick      *pb.MrpcMqlTick
    Age       time.Duration
    IsFresh   bool
    HasLast   bool
}

func GetValidatedTick(account *mt5.MT5Account, symbol string, maxAge time.Duration) (*ValidatedTick, error) {
    ctx := context.Background()

    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
        Symbol: symbol,
    })
    if err != nil {
        return nil, err
    }

    tickTime := time.Unix(tick.Time, 0)
    age := time.Since(tickTime)

    return &ValidatedTick{
        Tick:    tick,
        Age:     age,
        IsFresh: age <= maxAge,
        HasLast: tick.Last > 0,
    }, nil
}

// Usage:
// vt, _ := GetValidatedTick(account, "EURUSD", 10*time.Second)
// if vt.IsFresh {
//     fmt.Printf("Fresh tick: Bid=%.5f, Age=%v\n", vt.Tick.Bid, vt.Age)
// }
```

---

## 🔧 Common Patterns

### Safe price retrieval with retry

```go
func GetTickWithRetry(account *mt5.MT5Account, symbol string, maxRetries int) (*pb.MrpcMqlTick, error) {
    ctx := context.Background()

    for i := 0; i < maxRetries; i++ {
        tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
            Symbol: symbol,
        })

        if err == nil && tick.Bid > 0 && tick.Ask > 0 {
            return tick, nil
        }

        if i < maxRetries-1 {
            time.Sleep(500 * time.Millisecond)
        }
    }

    return nil, fmt.Errorf("failed to get valid tick after %d retries", maxRetries)
}
```

### Calculate price statistics

```go
func GetPriceStats(account *mt5.MT5Account, symbol string) {
    ctx := context.Background()

    tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
        Symbol: symbol,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    spread := tick.Ask - tick.Bid
    midPrice := (tick.Bid + tick.Ask) / 2.0
    spreadPct := (spread / midPrice) * 100.0

    fmt.Printf("%s Statistics:\n", symbol)
    fmt.Printf("  Bid: %.5f\n", tick.Bid)
    fmt.Printf("  Ask: %.5f\n", tick.Ask)
    fmt.Printf("  Mid: %.5f\n", midPrice)
    fmt.Printf("  Spread: %.5f (%.4f%%)\n", spread, spreadPct)
    fmt.Printf("  Last: %.5f\n", tick.Last)
    fmt.Printf("  Volume: %.2f\n", tick.VolumeReal)
}
```

---

## 📚 See Also

* [SymbolInfoDouble](./SymbolInfoDouble.md) - Get individual price properties
* [OnSymbolTick](../7.%20Streaming_Methods/OnSymbolTick.md) - Stream real-time tick updates
* [SymbolIsSynchronized](./SymbolIsSynchronized.md) - Check quote synchronization
* [OrderSend](../4.%20Trading_Operations/OrderSend.md) - Place orders with current prices
