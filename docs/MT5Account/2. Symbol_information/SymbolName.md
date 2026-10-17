# ✅ Get Symbol Name by Index

> **Request:** get trading symbol name by its position in the list. Allows iteration through all available symbols or only symbols in Market Watch.

**API Information:**

* **Low-level API:** `MT5Account.SymbolName(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolName` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolName(SymbolNameRequest) → SymbolNameReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolName(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieves symbol name by its position index in the symbols list.
* **Why you need it.** Iterate through all available symbols without knowing their names in advance.
* **Zero-based indexing.** Position starts at 0 and goes up to SymbolsTotal-1.

---

## 🎯 Purpose

Use it to:

* Iterate through all available symbols
* Build dynamic symbol lists
* Discover available trading instruments
* Create symbol selection interfaces

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolName - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolName_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolName returns the name of a symbol by its position in the list.
// Use Index as zero-based index and Selected to choose between Market Watch or all symbols.
func (a *MT5Account) SymbolName(
    ctx context.Context,
    req *pb.SymbolNameRequest,
) (*pb.SymbolNameData, error)
```

**Request message:**

```protobuf
SymbolNameRequest {
  int32 index = 1;     // Zero-based position index
  bool selected = 2;   // true = Market Watch only, false = all symbols
}
```

**Reply message:**

```protobuf
SymbolNameReply {
  oneof response {
    SymbolNameData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                     | Description                                   |
| --------- | ------------------------ | --------------------------------------------- |
| `ctx`     | `context.Context`        | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolNameRequest`  | Request with Index and Selected flag          |

**Request fields:**

| Field      | Type    | Description                                                  |
| ---------- | ------- | ------------------------------------------------------------ |
| `Index`    | `int32` | Zero-based position index (0 to SymbolsTotal-1)              |
| `Selected` | `bool`  | `true` = Market Watch symbols only, `false` = all symbols    |

---

## ⬆️ Output — `SymbolNameData`

| Field  | Type     | Go Type  | Description                               |
| ------ | -------- | -------- | ----------------------------------------- |
| `Name` | `string` | `string` | Symbol name at the specified position     |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Index bounds:** Position must be within [0, SymbolsTotal-1]. Out of bounds returns error.
* **Performance:** For getting all symbols, combine with SymbolsTotal to pre-allocate slice.

---

## 🔗 Usage Examples

### 1) Get first symbol name

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

    // Get first symbol from Market Watch
    data, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
        Index:    0,
        Selected: true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("First Market Watch symbol: %s\n", data.Name)
}
```

### 2) List all Market Watch symbols

```go
func GetMarketWatchSymbols(account *mt5.MT5Account) ([]string, error) {
    ctx := context.Background()

    // Get total count first
    totalData, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Mode: true,
    })
    if err != nil {
        return nil, err
    }

    // Pre-allocate slice
    symbols := make([]string, 0, totalData.Total)

    // Iterate through all symbols
    for i := int32(0); i < totalData.Total; i++ {
        nameData, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
            Index:    i,
            Selected: true,
        })
        if err != nil {
            continue // Skip errors
        }
        symbols = append(symbols, nameData.Name)
    }

    return symbols, nil
}
```

### 3) Get all available symbols (not just Market Watch)

```go
func GetAllSymbols(account *mt5.MT5Account) ([]string, error) {
    ctx := context.Background()

    // Get total count of all symbols
    totalData, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Mode: false,
    })
    if err != nil {
        return nil, err
    }

    symbols := make([]string, 0, totalData.Total)

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

### 4) Filter symbols by prefix

```go
func GetSymbolsByPrefix(account *mt5.MT5Account, prefix string) ([]string, error) {
    ctx := context.Background()

    totalData, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Mode: false,
    })
    if err != nil {
        return nil, err
    }

    symbols := []string{}

    for i := int32(0); i < totalData.Total; i++ {
        nameData, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
            Index:    i,
            Selected: false,
        })
        if err != nil {
            continue
        }

        if strings.HasPrefix(nameData.Name, prefix) {
            symbols = append(symbols, nameData.Name)
        }
    }

    return symbols, nil
}

// Usage:
// eurSymbols, _ := GetSymbolsByPrefix(account, "EUR")
// fmt.Printf("EUR symbols: %v\n", eurSymbols)
```

### 5) Get symbol at specific position with validation

```go
func GetSymbolAtPosition(account *mt5.MT5Account, position int32, marketWatchOnly bool) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Validate position bounds
    totalData, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Mode: marketWatchOnly,
    })
    if err != nil {
        return "", fmt.Errorf("failed to get total: %w", err)
    }

    if position < 0 || position >= totalData.Total {
        return "", fmt.Errorf("position %d out of bounds [0, %d)", position, totalData.Total)
    }

    // Get symbol name
    nameData, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
        Index:    position,
        Selected: marketWatchOnly,
    })
    if err != nil {
        return "", fmt.Errorf("failed to get symbol name: %w", err)
    }

    return nameData.Name, nil
}
```

---

## 🔧 Common Patterns

### Efficient iteration with pre-allocation

```go
func IterateAllSymbols(account *mt5.MT5Account) []string {
    ctx := context.Background()

    // Get count first
    total, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: false})
    if err != nil {
        return nil
    }

    // Pre-allocate exact capacity
    symbols := make([]string, 0, total.Total)

    // Fast iteration
    for i := int32(0); i < total.Total; i++ {
        name, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
            Index:    i,
            Selected: false,
        })
        if err == nil {
            symbols = append(symbols, name.Name)
        }
    }

    return symbols
}
```

### Parallel symbol name fetching

```go
func GetSymbolsParallel(account *mt5.MT5Account) []string {
    ctx := context.Background()

    total, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{Mode: true})
    if err != nil {
        return nil
    }

    symbols := make([]string, total.Total)
    var wg sync.WaitGroup

    for i := int32(0); i < total.Total; i++ {
        wg.Add(1)
        go func(pos int32) {
            defer wg.Done()
            name, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
                Index:    pos,
                Selected: true,
            })
            if err == nil {
                symbols[pos] = name.Name
            }
        }(i)
    }

    wg.Wait()
    return symbols
}
```

---

## 📚 See Also

* [SymbolsTotal](./SymbolsTotal.md) - Get total count of available symbols
* [SymbolExist](./SymbolExist.md) - Check if symbol exists
* [SymbolSelect](./SymbolSelect.md) - Add/remove symbols from Market Watch
* [SymbolInfoTick](./SymbolInfoTick.md) - Get last tick data for symbol
