# ✅ Get Total Count of Available Symbols

> **Request:** get the total number of symbols available in MetaTrader 5 terminal. You can get the count of all symbols or only those in Market Watch.

**API Information:**

* **Low-level API:** `MT5Account.SymbolsTotal(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolsTotal` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolsTotal(SymbolsTotalRequest) → SymbolsTotalReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolsTotal(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Returns the count of available symbols in MT5 terminal.
* **Why you need it.** Iterate through all symbols, validate Market Watch setup, or build symbol lists.
* **Mode flag.** Use `Mode: true` to count only visible symbols in Market Watch, `false` for all available symbols.

---

## 🎯 Purpose

Use it to:

* Get total symbol count before iterating with `SymbolName()`
* Validate Market Watch configuration
* Build dynamic symbol selection interfaces
* Monitor available trading instruments

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolsTotal - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolsTotal_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolsTotal returns the total number of available symbols.
// Use Mode=true to count only symbols in Market Watch, false for all symbols.
func (a *MT5Account) SymbolsTotal(
    ctx context.Context,
    req *pb.SymbolsTotalRequest,
) (*pb.SymbolsTotalData, error)
```

**Request message:**

```protobuf
SymbolsTotalRequest {
  bool mode = 1;  // true = Market Watch only, false = all symbols
}
```

**Reply message:**

```protobuf
SymbolsTotalReply {
  oneof response {
    SymbolsTotalData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                        | Description                                                     |
| --------- | --------------------------- | --------------------------------------------------------------- |
| `ctx`     | `context.Context`           | Context for deadline/timeout and cancellation                   |
| `req`     | `*pb.SymbolsTotalRequest`   | Request with `Mode` flag                                        |

**Request fields:**

| Field      | Type   | Description                                                     |
| ---------- | ------ | --------------------------------------------------------------- |
| `Mode`     | `bool` | `true` = count only Market Watch symbols, `false` = all symbols |

---

## ⬆️ Output — `SymbolsTotalData`

| Field   | Type    | Go Type | Description                           |
| ------- | ------- | ------- | ------------------------------------- |
| `Total` | `int32` | `int32` | Total number of symbols available     |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Market Watch vs All:** Mode=true returns only symbols visible in Market Watch window, Mode=false returns all symbols available in terminal.

---

## 🔗 Usage Examples

### 1) Get total symbols in Market Watch

```go
package main

import (
    "github.com/google/uuid"
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
)

func main() {
    account, err := mt5.NewMT5Account(591129415, "password", "mt5.mrpc.pro:443", uuid.Nil)
    if err != nil {
        panic(err)
    }
    defer account.Close()

    // Connect first
    _, err = account.ConnectEx(context.Background(), &pb.ConnectExRequest{User: 591129415, Password: "password", MtClusterName: "MetaQuotes-Demo"})
    if err != nil {
        panic(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Count only symbols in Market Watch
    data, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Mode: true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Symbols in Market Watch: %d\n", data.Total)
}
```

### 2) Get total available symbols (all)

```go
func GetAllSymbolsCount(account *mt5.MT5Account) (int32, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Count all available symbols
    data, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Mode: false,
    })
    if err != nil {
        return 0, fmt.Errorf("failed to get symbols total: %w", err)
    }

    return data.Total, nil
}

// Usage:
// total, _ := GetAllSymbolsCount(account)
// fmt.Printf("Total available symbols: %d\n", total)
```

### 3) Compare Market Watch vs All symbols

```go
func CompareSymbolCounts(account *mt5.MT5Account) {
    ctx := context.Background()

    // Get Market Watch count
    selected, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: true})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Get all symbols count
    all, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: false})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Market Watch: %d symbols\n", selected.Total)
    fmt.Printf("All available: %d symbols\n", all.Total)
    fmt.Printf("Not in Market Watch: %d symbols\n", all.Total-selected.Total)
}
```

### 4) Validate before iteration

```go
func ListAllSymbolNames(account *mt5.MT5Account) ([]string, error) {
    ctx := context.Background()

    // First, get total count
    totalData, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: false})
    if err != nil {
        return nil, err
    }

    symbols := make([]string, 0, totalData.Total)

    // Iterate through all symbols
    for i := int32(0); i < totalData.Total; i++ {
        nameData, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
            Index:    i,
            Selected: false,
        })
        if err != nil {
            continue
        }
        symbols = append(symbols, nameData.Name)
    }

    return symbols, nil
}
```

### 5) Monitor Market Watch changes

```go
func MonitorMarketWatch(account *mt5.MT5Account, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    var lastCount int32 = -1

    for range ticker.C {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        data, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: true})
        cancel()

        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        if lastCount != -1 && data.Total != lastCount {
            fmt.Printf("[%s] Market Watch changed: %d -> %d symbols\n",
                time.Now().Format("15:04:05"),
                lastCount,
                data.Total,
            )
        }

        lastCount = data.Total
    }
}

// Usage:
// MonitorMarketWatch(account, 5*time.Second)
```

---

## 🔧 Common Patterns

### Pre-allocate slice for iteration

```go
func GetMarketWatchSymbols(account *mt5.MT5Account) ([]string, error) {
    ctx := context.Background()

    // Get count first to pre-allocate
    total, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: true})
    if err != nil {
        return nil, err
    }

    // Pre-allocate with exact capacity
    symbols := make([]string, 0, total.Total)

    for i := int32(0); i < total.Total; i++ {
        nameData, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
            Index:    i,
            Selected: true,
        })
        if err != nil {
            continue
        }
        symbols = append(symbols, nameData.Name)
    }

    return symbols, nil
}
```

### Validate symbol availability

```go
func HasMinimumSymbols(account *mt5.MT5Account, minRequired int32) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    data, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: false})
    if err != nil {
        return false, err
    }

    return data.Total >= minRequired, nil
}
```

---

## 📚 See Also

* [SymbolName](./SymbolName.md) - Get symbol name by position index
* [SymbolExist](./SymbolExist.md) - Check if symbol exists
* [SymbolSelect](./SymbolSelect.md) - Add/remove symbols from Market Watch
* [AccountSummary](../1.%20Account_information/AccountSummary.md) - Account summary information
