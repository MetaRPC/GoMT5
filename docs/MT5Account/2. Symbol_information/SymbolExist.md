# ✅ Check Symbol Existence

> **Request:** check if specified symbol exists in MetaTrader 5 terminal. Returns existence flag and custom symbol indicator.

**API Information:**

* **SDK wrapper:** `MT5Account.SymbolExist(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolExist` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolExist(SymbolExistRequest) → SymbolExistReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolExist(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

---

## 💬 Just the essentials

* **What it is.** Checks if a trading symbol exists in the MT5 terminal.
* **Why you need it.** Validate symbol names before requesting data or placing orders.
* **Custom symbols.** IsCustom flag indicates whether the symbol is user-created or standard broker symbol.

---

## 🎯 Purpose

Use it to:

* Validate symbol names before trading operations
* Detect typos in symbol names
* Identify custom symbols in the terminal
* Build symbol selection and validation logic

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolExist - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolExist_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolExist checks if a symbol with specified name exists.
// Returns Exist flag and IsCustom flag indicating if it's a custom symbol.
func (a *MT5Account) SymbolExist(
    ctx context.Context,
    req *pb.SymbolExistRequest,
) (*pb.SymbolExistData, error)
```

**Request message:**

```protobuf
SymbolExistRequest {
  string Symbol = 1;  // Symbol name to check (e.g., "EURUSD")
}
```

**Reply message:**

```protobuf
SymbolExistReply {
  oneof response {
    SymbolExistData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                        | Description                                   |
| --------- | --------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`           | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolExistRequest`    | Request with Symbol name                      |

**Request fields:**

| Field    | Type     | Description                               |
| -------- | -------- | ----------------------------------------- |
| `Symbol` | `string` | Symbol name to check (e.g., "EURUSD")     |

---

## ⬆️ Output — `SymbolExistData`

| Field      | Type   | Go Type | Description                                  |
| ---------- | ------ | ------- | -------------------------------------------- |
| `Exists`   | `bool` | `bool`  | True if symbol exists in terminal            |
| `IsCustom` | `bool` | `bool`  | True if it's a custom (user-defined) symbol  |


---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Case sensitivity:** Symbol names are typically case-insensitive but should match broker convention.

---

## 🔗 Usage Examples

### 1) Check if symbol exists

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
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    data, err := account.SymbolExist(ctx, &pb.SymbolExistRequest{
        Symbol: "EURUSD",
    })
    if err != nil {
        panic(err)
    }

    if data.Exists {
        fmt.Printf("Symbol EURUSD exists (custom: %v)\n", data.IsCustom)
    } else {
        fmt.Println("Symbol EURUSD not found")
    }
}
```

### 2) Validate symbol before trading

```go
func ValidateAndTrade(account *mt5.MT5Account, symbol string, volume float64) error {
    ctx := context.Background()

    // First, check if symbol exists
    existData, err := account.SymbolExist(ctx, &pb.SymbolExistRequest{
        Symbol: symbol,
    })
    if err != nil {
        return fmt.Errorf("failed to check symbol: %w", err)
    }

    if !existData.Exists {
        return fmt.Errorf("symbol %s does not exist", symbol)
    }

    // Proceed with trading
    fmt.Printf("Symbol %s validated, proceeding with order...\n", symbol)
    return nil
}
```

### 3) Check multiple symbols

```go
func CheckSymbols(account *mt5.MT5Account, symbols []string) map[string]bool {
    ctx := context.Background()
    results := make(map[string]bool)

    for _, symbol := range symbols {
        data, err := account.SymbolExist(ctx, &pb.SymbolExistRequest{
            Symbol: symbol,
        })
        if err != nil {
            results[symbol] = false
            continue
        }
        results[symbol] = data.Exists
    }

    return results
}

// Usage:
// symbols := []string{"EURUSD", "GBPUSD", "INVALID", "BTCUSD"}
// results := CheckSymbols(account, symbols)
// for symbol, exists := range results {
//     fmt.Printf("%s: %v\n", symbol, exists)
// }
```

### 4) Filter custom symbols

```go
func GetCustomSymbols(account *mt5.MT5Account, symbols []string) []string {
    ctx := context.Background()
    customSymbols := []string{}

    for _, symbol := range symbols {
        data, err := account.SymbolExist(ctx, &pb.SymbolExistRequest{
            Symbol: symbol,
        })
        if err != nil || !data.Exists {
            continue
        }

        if data.IsCustom {
            customSymbols = append(customSymbols, symbol)
        }
    }

    return customSymbols
}
```

### 5) Symbol existence guard

```go
type SymbolValidator struct {
    account *mt5.MT5Account
    cache   map[string]bool
}

func (v *SymbolValidator) IsValid(symbol string) (bool, error) {
    // Check cache first
    if exists, ok := v.cache[symbol]; ok {
        return exists, nil
    }

    // Query MT5
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    data, err := v.account.SymbolExist(ctx, &pb.SymbolExistRequest{
        Symbol: symbol,
    })
    if err != nil {
        return false, err
    }

    // Cache result
    v.cache[symbol] = data.Exists
    return data.Exists, nil
}
```

---

## 🔧 Common Patterns

### Pre-validation before operations

```go
func SafeSymbolOperation(account *mt5.MT5Account, symbol string) error {
    ctx := context.Background()

    // Validate symbol exists
    existData, err := account.SymbolExist(ctx, &pb.SymbolExistRequest{
        Symbol: symbol,
    })
    if err != nil {
        return fmt.Errorf("validation error: %w", err)
    }

    if !existData.Exists {
        return fmt.Errorf("symbol %s not available", symbol)
    }

    // Proceed with actual operation
    fmt.Printf("Operating on %s (custom: %v)\n", symbol, existData.IsCustom)
    return nil
}
```

### Symbol existence checker with timeout

```go
func QuickSymbolCheck(account *mt5.MT5Account, symbol string) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    data, err := account.SymbolExist(ctx, &pb.SymbolExistRequest{
        Symbol: symbol,
    })

    return err == nil && data.Exists
}
```

---

## 📚 See Also

* [SymbolsTotal](./SymbolsTotal.md) - Get total count of available symbols
* [SymbolName](./SymbolName.md) - Get symbol name by position index
* [SymbolSelect](./SymbolSelect.md) - Add/remove symbols from Market Watch
* [SymbolInfoTick](./SymbolInfoTick.md) - Get last tick data for symbol
