# ✅ Get Parameters for Multiple Symbols

> **Request:** get complete trading parameters for multiple symbols in one request. Returns Bid, Ask, Digits, Spread, volumes, contract size and other parameters.

**API Information:**

* **Low-level API:** `MT5Account.SymbolParamsMany(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `SymbolParamsMany` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `SymbolParamsMany(SymbolParamsManyRequest) → SymbolParamsManyReply`
* **Low‑level client (generated):** `AccountHelperClient.SymbolParamsMany(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieves comprehensive trading parameters for multiple symbols at once.
* **Why you need it.** Efficient bulk data retrieval, avoid multiple round-trips.
* **Complete data.** Returns all essential symbol properties in one call.

---

## 🎯 Purpose

Use it to:

* Get complete trading specifications for multiple symbols efficiently
* Build trading dashboards and watchlists
* Validate trading parameters in bulk
* Compare spreads and conditions across symbols
* Reduce network round-trips with batch queries

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolParamsMany retrieves detailed parameters for multiple symbols in one call.
// This is the recommended method for getting comprehensive symbol data.
func (a *MT5Account) SymbolParamsMany(
    ctx context.Context,
    req *pb.SymbolParamsManyRequest,
) (*pb.SymbolParamsManyData, error)
```

**Request message:**

```protobuf
SymbolParamsManyRequest {
  repeated string Symbols = 1;  // Array of symbol names
}
```

---

## 🔽 Input

| Parameter | Type                            | Description                                   |
| --------- | ------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`               | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolParamsManyRequest`   | Request with array of Symbols                 |

**Request fields:**

| Field     | Type       | Description                        |
| --------- | ---------- | ---------------------------------- |
| `Symbols` | `string[]` | Array of symbol names              |

---

## ⬆️ Output — `SymbolParamsManyData`

| Field    | Type              | Go Type              | Description                              |
| -------- | ----------------- | -------------------- | ---------------------------------------- |
| `Params` | `SymbolParams[]`  | `[]*pb.SymbolParams` | Array of symbol parameters               |

**SymbolParams structure includes:**

- `Symbol` - Symbol name
- `Bid` - Current bid price
- `Ask` - Current ask price
- `Digits` - Decimal precision
- `Spread` - Current spread
- `VolumeMin` - Minimum volume
- `VolumeMax` - Maximum volume
- `VolumeStep` - Volume step
- `ContractSize` - Contract size
- `Point` - Point size
- And many other trading parameters...

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolParamsMany - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolParamsMany_HOW.md)**

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Most efficient:** This is the recommended method for retrieving comprehensive symbol data.
* **Bulk operations:** Can query 100+ symbols in single request.
* **112+ fields:** Each SymbolParams contains over 112 different properties.
* **Performance:** Much faster than calling individual SymbolInfo* methods repeatedly.

---

## 🔗 Usage Examples

### 1) Get params for multiple symbols

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

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: []string{"EURUSD", "GBPUSD", "USDJPY"},
    })
    if err != nil {
        panic(err)
    }

    for _, params := range data.Params {
        fmt.Printf("%s:\n", params.Symbol)
        fmt.Printf("  Bid: %.5f, Ask: %.5f\n", params.Bid, params.Ask)
        fmt.Printf("  Digits: %d, Spread: %d\n", params.Digits, params.Spread)
        fmt.Printf("  Volume: %.2f - %.2f (step: %.2f)\n",
            params.VolumeMin, params.VolumeMax, params.VolumeStep)
        fmt.Println()
    }
}
```

### 2) Build trading dashboard

```go
func BuildTradingDashboard(account *mt5.MT5Account, watchlist []string) {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: watchlist,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println("╔════════════════════════════════════════════════════════════╗")
    fmt.Println("║                   TRADING DASHBOARD                        ║")
    fmt.Println("╠═══════════╦══════════╦══════════╦════════╦════════════════╣")
    fmt.Println("║  Symbol   ║   Bid    ║   Ask    ║ Spread ║ Contract Size  ║")
    fmt.Println("╠═══════════╬══════════╬══════════╬════════╬════════════════╣")

    for _, params := range data.Params {
        fmt.Printf("║ %-9s ║ %8.5f ║ %8.5f ║ %6d ║ %14.0f ║\n",
            params.Symbol, params.Bid, params.Ask, params.Spread, params.ContractSize)
    }

    fmt.Println("╚═══════════╩══════════╩══════════╩════════╩════════════════╝")
}

// Usage:
// BuildTradingDashboard(account, []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"})
```

### 3) Validate trading parameters

```go
func ValidateTradingParams(account *mt5.MT5Account, symbol string, volume float64) error {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: []string{symbol},
    })
    if err != nil {
        return err
    }

    if len(data.Params) == 0 {
        return fmt.Errorf("symbol %s not found", symbol)
    }

    params := data.Params[0]

    // Validate volume
    if volume < params.VolumeMin {
        return fmt.Errorf("volume %.2f below minimum %.2f", volume, params.VolumeMin)
    }
    if volume > params.VolumeMax {
        return fmt.Errorf("volume %.2f above maximum %.2f", volume, params.VolumeMax)
    }

    // Check volume step
    remainder := math.Mod(volume, params.VolumeStep)
    if remainder > 0.0001 {
        return fmt.Errorf("volume %.2f not aligned to step %.2f", volume, params.VolumeStep)
    }

    fmt.Printf("Volume %.2f is valid for %s\n", volume, symbol)
    return nil
}
```

### 4) Compare spreads across symbols

```go
func CompareSpreads(account *mt5.MT5Account, symbols []string) {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: symbols,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println("Spread comparison:")
    for _, params := range data.Params {
        spreadPips := float64(params.Spread) * params.Point * 10
        fmt.Printf("  %s: %d points (%.1f pips)\n",
            params.Symbol, params.Spread, spreadPips)
    }
}

// Usage:
// CompareSpreads(account, []string{"EURUSD", "GBPUSD", "EURGBP"})
```

### 5) Find symbols with low spreads

```go
func FindLowSpreadSymbols(account *mt5.MT5Account, symbols []string, maxSpreadPips float64) []string {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: symbols,
    })
    if err != nil {
        return nil
    }

    lowSpreadSymbols := []string{}

    for _, params := range data.Params {
        spreadPips := float64(params.Spread) * params.Point * 10

        if spreadPips <= maxSpreadPips {
            lowSpreadSymbols = append(lowSpreadSymbols, params.Symbol)
            fmt.Printf("%s: %.1f pips\n", params.Symbol, spreadPips)
        }
    }

    return lowSpreadSymbols
}

// Usage:
// symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "EURGBP"}
// lowSpread := FindLowSpreadSymbols(account, symbols, 2.0)
```

---

## 📚 See Also

* [SymbolInfoDouble](../2.%20Symbol_information/SymbolInfoDouble.md) - Get individual symbol properties
* [SymbolInfoInteger](../2.%20Symbol_information/SymbolInfoInteger.md) - Get integer symbol properties
* [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - Get tick data
