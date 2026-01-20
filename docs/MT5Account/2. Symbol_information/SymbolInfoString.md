# ✅ Get Symbol String Property

> **Request:** get string property of trading symbol, such as description, base currency, profit currency and other text parameters.

**API Information:**

* **SDK wrapper:** `MT5Account.SymbolInfoString(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolInfoString` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolInfoString(SymbolInfoStringRequest) → SymbolInfoStringReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolInfoString(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

## 💬 Just the essentials

* **What it is.** Retrieves text (string) properties of a trading symbol.
* **Why you need it.** Get symbol description, currency information, and metadata.
* **Property ID required.** Must specify which property to retrieve using PropertyId enum.

---

## 🎯 Purpose

Use it to:

* Get symbol full description/name
* Identify base and profit currencies
* Build symbol information displays
* Categorize symbols by currency pairs
* Validate currency compatibility

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolInfoString - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoString_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolInfoString retrieves a string-type symbol property.
// Use PropertyId to specify which property (DESCRIPTION, CURRENCY_BASE, etc).
func (a *MT5Account) SymbolInfoString(
    ctx context.Context,
    req *pb.SymbolInfoStringRequest,
) (*pb.SymbolInfoStringData, error)
```

**Request message:**

```protobuf
SymbolInfoStringRequest {
  string Symbol = 1;       // Symbol name
  int32 PropertyId = 2;    // Property identifier (DESCRIPTION, CURRENCY_BASE, etc)
}
```

**Reply message:**

```protobuf
SymbolInfoStringReply {
  oneof response {
    SymbolInfoStringData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                            | Description                                   |
| --------- | ------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`               | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolInfoStringRequest`   | Request with Symbol and PropertyId            |

**Request fields:**

| Field        | Type     | Description                                                        |
| ------------ | -------- | ------------------------------------------------------------------ |
| `Symbol`     | `string` | Symbol name (e.g., "EURUSD")                                       |
| `PropertyId` | `int32`  | Property identifier from SymbolInfoStringProperty enum             |

**PropertyId enum values (SymbolInfoStringProperty):**

| Constant | Value | Description |
|----------|-------|-------------|
| `SYMBOL_BASIS` | 0 | Underlying asset of derivative |
| `SYMBOL_CATEGORY` | 1 | Symbol category |
| `SYMBOL_COUNTRY` | 2 | Country the symbol belongs to |
| `SYMBOL_SECTOR_NAME` | 3 | Economic sector name |
| `SYMBOL_INDUSTRY_NAME` | 4 | Industry name |
| `SYMBOL_CURRENCY_BASE` | 5 | Base currency (first currency in pair) |
| `SYMBOL_CURRENCY_PROFIT` | 6 | Profit currency (second currency in pair) |
| `SYMBOL_CURRENCY_MARGIN` | 7 | Margin calculation currency |
| `SYMBOL_BANK` | 8 | Source of last quote (bank/exchange) |
| `SYMBOL_DESCRIPTION` | 9 | Symbol description/full name |
| `SYMBOL_EXCHANGE` | 10 | Exchange the symbol is traded on |
| `SYMBOL_FORMULA` | 11 | Formula for custom symbol price calculation |
| `SYMBOL_ISIN` | 12 | International Securities Identification Number |
| `SYMBOL_PAGE` | 13 | Web page URL with symbol information |
| `SYMBOL_PATH` | 14 | Symbol path in Market Watch tree |

---

## ⬆️ Output — `SymbolInfoStringData`

| Field   | Type     | Go Type  | Description                           |
| ------- | -------- | -------- | ------------------------------------- |
| `Value` | `string` | `string` | The requested string property value   |

---



## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Localization:** Symbol descriptions may vary based on broker's language settings.
* **Currency codes:** Base and profit currencies are returned as 3-letter ISO codes (e.g., "USD", "EUR").

---

## 🔗 Usage Examples

### 1) Get symbol description

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
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    data, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     "EURUSD",
        PropertyId: 9, // SYMBOL_DESCRIPTION
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("EURUSD Description: %s\n", data.Value)
}
```

### 2) Get base currency

```go
func GetBaseCurrency(account *mt5.MT5Account, symbol string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 5, // SYMBOL_CURRENCY_BASE
    })
    if err != nil {
        return "", fmt.Errorf("failed to get base currency: %w", err)
    }

    return data.Value, nil
}

// Usage:
// baseCurrency, _ := GetBaseCurrency(account, "EURUSD")
// fmt.Printf("Base currency: %s\n", baseCurrency) // "EUR"
```

### 3) Get profit currency

```go
func GetProfitCurrency(account *mt5.MT5Account, symbol string) (string, error) {
    ctx := context.Background()

    data, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 6, // SYMBOL_CURRENCY_PROFIT
    })
    if err != nil {
        return "", err
    }

    return data.Value, nil
}

// Usage:
// profitCurrency, _ := GetProfitCurrency(account, "EURUSD")
// fmt.Printf("Profit currency: %s\n", profitCurrency) // "USD"
```

### 4) Get complete symbol info

```go
type SymbolInfo struct {
    Symbol          string
    Description     string
    BaseCurrency    string
    ProfitCurrency  string
    MarginCurrency  string
}

func GetSymbolInfo(account *mt5.MT5Account, symbol string) (*SymbolInfo, error) {
    ctx := context.Background()

    // Get description
    desc, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 9, // SYMBOL_DESCRIPTION
    })
    if err != nil {
        return nil, err
    }

    // Get base currency
    base, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 5, // SYMBOL_CURRENCY_BASE
    })
    if err != nil {
        return nil, err
    }

    // Get profit currency
    profit, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 6, // SYMBOL_CURRENCY_PROFIT
    })
    if err != nil {
        return nil, err
    }

    // Get margin currency
    margin, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 7, // SYMBOL_CURRENCY_MARGIN
    })
    if err != nil {
        return nil, err
    }

    return &SymbolInfo{
        Symbol:         symbol,
        Description:    desc.Value,
        BaseCurrency:   base.Value,
        ProfitCurrency: profit.Value,
        MarginCurrency: margin.Value,
    }, nil
}
```

### 5) Filter symbols by base currency

```go
func GetSymbolsByBaseCurrency(account *mt5.MT5Account, baseCurrency string) ([]string, error) {
    ctx := context.Background()

    // Get all symbols
    totalData, err := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
        Selected: false,
    })
    if err != nil {
        return nil, err
    }

    symbols := []string{}

    for i := int64(0); i < totalData.Total; i++ {
        // Get symbol name
        nameData, err := account.SymbolName(ctx, &pb.SymbolNameRequest{
            Pos:      i,
            Selected: false,
        })
        if err != nil {
            continue
        }

        // Get base currency
        currData, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
            Symbol:     nameData.Name,
            PropertyId: 5, // SYMBOL_CURRENCY_BASE
        })
        if err != nil {
            continue
        }

        if currData.Value == baseCurrency {
            symbols = append(symbols, nameData.Name)
        }
    }

    return symbols, nil
}

// Usage:
// eurSymbols, _ := GetSymbolsByBaseCurrency(account, "EUR")
// fmt.Printf("EUR-based symbols: %v\n", eurSymbols)
```

### 6) Display symbol details

```go
func DisplaySymbolDetails(account *mt5.MT5Account, symbol string) {
    ctx := context.Background()

    // Get description
    desc, _ := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 9, // SYMBOL_DESCRIPTION
    })

    // Get currencies
    base, _ := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 5, // SYMBOL_CURRENCY_BASE
    })

    profit, _ := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 6, // SYMBOL_CURRENCY_PROFIT
    })

    // Get path
    path, _ := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 14, // SYMBOL_PATH
    })

    fmt.Printf("Symbol: %s\n", symbol)
    fmt.Printf("Description: %s\n", desc.Value)
    fmt.Printf("Base Currency: %s\n", base.Value)
    fmt.Printf("Profit Currency: %s\n", profit.Value)
    fmt.Printf("Path: %s\n", path.Value)
}
```

---

## 🔧 Common Patterns

### Build currency pair display

```go
func FormatCurrencyPair(account *mt5.MT5Account, symbol string) string {
    ctx := context.Background()

    base, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 5, // SYMBOL_CURRENCY_BASE
    })
    if err != nil {
        return symbol
    }

    profit, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 6, // SYMBOL_CURRENCY_PROFIT
    })
    if err != nil {
        return symbol
    }

    return fmt.Sprintf("%s/%s", base.Value, profit.Value)
}

// Usage:
// pair := FormatCurrencyPair(account, "EURUSD")
// fmt.Println(pair) // "EUR/USD"
```

### Validate currency compatibility

```go
func IsAccountCurrencyCompatible(account *mt5.MT5Account, symbol string, accountCurrency string) (bool, error) {
    ctx := context.Background()

    profitCurrency, err := account.SymbolInfoString(ctx, &pb.SymbolInfoStringRequest{
        Symbol:     symbol,
        PropertyId: 6, // SYMBOL_CURRENCY_PROFIT
    })
    if err != nil {
        return false, err
    }

    return profitCurrency.Value == accountCurrency, nil
}
```

---

## 📚 See Also

* [SymbolInfoDouble](./SymbolInfoDouble.md) - Get double symbol properties
* [SymbolInfoInteger](./SymbolInfoInteger.md) - Get integer symbol properties
* [SymbolParamsMany](../6.%20Additional_Methods/SymbolParamsMany.md) - Get all symbol parameters at once
* [AccountSummary](../1.%20Account_information/AccountSummary.md) - Get account currency information
