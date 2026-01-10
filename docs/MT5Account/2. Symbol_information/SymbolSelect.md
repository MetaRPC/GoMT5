# ✅ Add/Remove Symbol from Market Watch

> **Request:** add trading symbol to Market Watch window or remove it. Controls symbol visibility in terminal.

**API Information:**

* **SDK wrapper:** `MT5Account.SymbolSelect(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolSelect` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolSelect(SymbolSelectRequest) → SymbolSelectReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolSelect(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

## 💬 Just the essentials

* **What it is.** Adds or removes trading symbols from the Market Watch window.
* **Why you need it.** Control which symbols receive real-time quotes and are visible in terminal.
* **Required for quotes.** Symbol must be in Market Watch to receive streaming tick updates.

---

## 🎯 Purpose

Use it to:

* Add symbols to Market Watch before requesting quotes
* Remove unused symbols to reduce data traffic
* Manage symbol visibility dynamically
* Prepare symbols for trading or monitoring

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolSelect - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolSelect_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolSelect adds or removes a symbol from Market Watch window.
// Use Select=true to add, Select=false to remove.
func (a *MT5Account) SymbolSelect(
    ctx context.Context,
    req *pb.SymbolSelectRequest,
) (*pb.SymbolSelectData, error)
```

**Request message:**

```protobuf
SymbolSelectRequest {
  string Symbol = 1;  // Symbol name
  bool Select = 2;    // true = add to Market Watch, false = remove
}
```

**Reply message:**

```protobuf
SymbolSelectReply {
  oneof response {
    SymbolSelectData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                       | Description                                   |
| --------- | -------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`          | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolSelectRequest`  | Request with Symbol name and Select flag      |

**Request fields:**

| Field    | Type     | Description                                       |
| -------- | -------- | ------------------------------------------------- |
| `Symbol` | `string` | Symbol name (e.g., "EURUSD")                      |
| `Select` | `bool`   | `true` = add to Market Watch, `false` = remove    |

---

## ⬆️ Output — `SymbolSelectData`

| Field     | Type   | Go Type | Description                                 |
| --------- | ------ | ------- | ------------------------------------------- |
| `Success` | `bool` | `bool`  | True if operation succeeded                 |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Quote prerequisite:** Symbol must be in Market Watch to receive streaming quotes via OnSymbolTick.
* **Idempotent:** Calling Select=true on already selected symbol is safe (returns success).

---

## 🔗 Usage Examples

### 1) Add symbol to Market Watch

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

    // Add EURUSD to Market Watch
    data, err := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
        Symbol: "EURUSD",
        Select: true,
    })
    if err != nil {
        panic(err)
    }

    if data.Success {
        fmt.Println("EURUSD added to Market Watch")
    }
}
```

### 2) Remove symbol from Market Watch

```go
func RemoveSymbol(account *mt5.MT5Account, symbol string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
        Symbol: symbol,
        Select: false,
    })
    if err != nil {
        return fmt.Errorf("failed to remove symbol: %w", err)
    }

    if !data.Success {
        return fmt.Errorf("symbol removal unsuccessful")
    }

    return nil
}
```

### 3) Ensure symbol is in Market Watch before trading

```go
func PrepareSymbolForTrading(account *mt5.MT5Account, symbol string) error {
    ctx := context.Background()

    // First check if it exists
    existData, err := account.SymbolExist(ctx, &pb.SymbolExistRequest{
        Symbol: symbol,
    })
    if err != nil || !existData.Exist {
        return fmt.Errorf("symbol %s does not exist", symbol)
    }

    // Add to Market Watch
    selectData, err := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
        Symbol: symbol,
        Select: true,
    })
    if err != nil {
        return fmt.Errorf("failed to add to Market Watch: %w", err)
    }

    if !selectData.Success {
        return fmt.Errorf("symbol selection failed")
    }

    fmt.Printf("Symbol %s ready for trading\n", symbol)
    return nil
}
```

### 4) Batch add symbols to Market Watch

```go
func AddSymbolsToMarketWatch(account *mt5.MT5Account, symbols []string) []string {
    ctx := context.Background()
    failed := []string{}

    for _, symbol := range symbols {
        data, err := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
            Symbol: symbol,
            Select: true,
        })

        if err != nil || !data.Success {
            failed = append(failed, symbol)
            fmt.Printf("Failed to add %s\n", symbol)
        } else {
            fmt.Printf("Added %s to Market Watch\n", symbol)
        }
    }

    return failed
}

// Usage:
// symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
// failed := AddSymbolsToMarketWatch(account, symbols)
```

### 5) Clean up Market Watch

```go
func CleanMarketWatch(account *mt5.MT5Account, keepSymbols []string) error {
    ctx := context.Background()

    // Get current Market Watch symbols
    totalData, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Selected: true,
    })
    if err != nil {
        return err
    }

    // Create keep set
    keepSet := make(map[string]bool)
    for _, s := range keepSymbols {
        keepSet[s] = true
    }

    // Remove symbols not in keep list
    for i := int64(0); i < totalData.Total; i++ {
        nameData, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
            Pos:      i,
            Selected: true,
        })
        if err != nil {
            continue
        }

        if !keepSet[nameData.Name] {
            account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
                Symbol: nameData.Name,
                Select: false,
            })
            fmt.Printf("Removed %s from Market Watch\n", nameData.Name)
        }
    }

    return nil
}
```

---

## 🔧 Common Patterns

### Safe symbol selection with retry

```go
func EnsureSymbolSelected(account *mt5.MT5Account, symbol string) error {
    ctx := context.Background()
    maxRetries := 3

    for i := 0; i < maxRetries; i++ {
        data, err := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
            Symbol: symbol,
            Select: true,
        })

        if err == nil && data.Success {
            return nil
        }

        time.Sleep(500 * time.Millisecond)
    }

    return fmt.Errorf("failed to select symbol %s after %d retries", symbol, maxRetries)
}
```

### Toggle symbol visibility

```go
func ToggleSymbol(account *mt5.MT5Account, symbol string, enable bool) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
        Symbol: symbol,
        Select: enable,
    })

    if err != nil {
        return err
    }

    if !data.Success {
        return fmt.Errorf("toggle operation failed")
    }

    action := "removed from"
    if enable {
        action = "added to"
    }
    fmt.Printf("%s %s Market Watch\n", symbol, action)
    return nil
}
```

---

## 📚 See Also

* [SymbolsTotal](./SymbolsTotal.md) - Get total count of available symbols
* [SymbolExist](./SymbolExist.md) - Check if symbol exists
* [SymbolName](./SymbolName.md) - Get symbol name by position
* [OnSymbolTick](../7.%20Streaming_Methods/OnSymbolTick.md) - Stream real-time quotes (requires symbol in Market Watch)
