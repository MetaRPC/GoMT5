# ✅ Get Symbol Margin Rates by Order Type

> **Request:** get initial and maintenance margin rate coefficients for different order types (Buy, Sell, Limit, Stop) for specified symbol.

**API Information:**

* **Low-level API:** `MT5Account.SymbolInfoMarginRate(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolInfoMarginRate` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolInfoMarginRate(SymbolInfoMarginRateRequest) → SymbolInfoMarginRateReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolInfoMarginRate(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Returns margin rate coefficients for different order types.
* **Why you need it.** Calculate precise margin requirements, understand leverage.
* **Two types.** InitialMarginRate for opening, MaintenanceMarginRate for maintaining positions.

---

## 🎯 Purpose

Use it to:

* Calculate margin requirements for different order types
* Understand leverage impact on positions
* Compare margin rates between buy and sell orders
* Validate available margin before placing orders

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolInfoMarginRate retrieves margin requirements for different order types.
// Returns InitialMarginRate and MaintenanceMarginRate values.
func (a *MT5Account) SymbolInfoMarginRate(
    ctx context.Context,
    req *pb.SymbolInfoMarginRateRequest,
) (*pb.SymbolInfoMarginRateData, error)
```

**Request message:**

```protobuf
SymbolInfoMarginRateRequest {
  string Symbol = 1;     // Trading symbol
  int32 OrderType = 2;   // Order type (BUY, SELL, etc)
}
```

---

## 🔽 Input

| Parameter | Type                                 | Description                                   |
| --------- | ------------------------------------ | --------------------------------------------- |
| `ctx`     | `context.Context`                    | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolInfoMarginRateRequest`    | Request with Symbol and OrderType             |

**Request fields:**

| Field       | Type     | Description                                |
| ----------- | -------- | ------------------------------------------ |
| `Symbol`    | `string` | Trading symbol (e.g., "EURUSD")            |
| `OrderType` | `int32`  | Order type (see enum below)                |

**OrderType enum values:**

| Value | Protobuf Enum | Description |
|-------|---------------|-------------|
| `0` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY` | Market buy order |
| `1` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_SELL` | Market sell order |
| `2` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY_LIMIT` | Buy limit pending order |
| `3` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_SELL_LIMIT` | Sell limit pending order |
| `4` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY_STOP` | Buy stop pending order |
| `5` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_SELL_STOP` | Sell stop pending order |
| `6` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY_STOP_LIMIT` | Buy stop limit pending order |
| `7` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_SELL_STOP_LIMIT` | Sell stop limit pending order |
| `8` | `pb.ENUM_ORDER_TYPE_ORDER_TYPE_CLOSE_BY` | Close by opposite position |

**Usage:** Can use either numeric value or protobuf enum constant

---

## ⬆️ Output — `SymbolInfoMarginRateData`

| Field                   | Type     | Go Type   | Description                             |
| ----------------------- | -------- | --------- | --------------------------------------- |
| `InitialMarginRate`     | `double` | `float64` | Initial margin coefficient              |
| `MaintenanceMarginRate` | `double` | `float64` | Maintenance margin coefficient          |

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolInfoMarginRate - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoMarginRate_HOW.md)**

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Margin rates vary:** Different order types have different margin requirements.
* **Initial vs Maintenance:** Initial margin is required to open position, maintenance margin is required to keep it open.
* **Broker specific:** Margin rates are determined by broker and can vary significantly.

---

## 🔗 Usage Examples

### 1) Get margin rates for buy order

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

    data, err := account.SymbolInfoMarginRate(ctx, &pb.SymbolInfoMarginRateRequest{
        Symbol:    "EURUSD",
        OrderType: 0, // BUY
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("EURUSD Buy margin rates:\n")
    fmt.Printf("  Initial: %.4f\n", data.InitialMarginRate)
    fmt.Printf("  Maintenance: %.4f\n", data.MaintenanceMarginRate)
}
```

### 2) Compare rates for buy vs sell

```go
func CompareMarginRates(account *mt5.MT5Account, symbol string) {
    ctx := context.Background()

    // Buy rates
    buyData, _ := account.SymbolInfoMarginRate(ctx, &pb.SymbolInfoMarginRateRequest{
        Symbol:    symbol,
        OrderType: 0, // BUY
    })

    // Sell rates
    sellData, _ := account.SymbolInfoMarginRate(ctx, &pb.SymbolInfoMarginRateRequest{
        Symbol:    symbol,
        OrderType: 1, // SELL
    })

    fmt.Printf("%s Margin Rates:\n", symbol)
    fmt.Printf("  Buy:  Initial=%.4f, Maintenance=%.4f\n",
        buyData.InitialMarginRate, buyData.MaintenanceMarginRate)
    fmt.Printf("  Sell: Initial=%.4f, Maintenance=%.4f\n",
        sellData.InitialMarginRate, sellData.MaintenanceMarginRate)
}
```

### 3) Get rates for all order types

```go
func GetAllMarginRates(account *mt5.MT5Account, symbol string) {
    ctx := context.Background()

    orderTypes := map[int32]string{
        0: "BUY",
        1: "SELL",
        2: "BUY_LIMIT",
        3: "SELL_LIMIT",
        4: "BUY_STOP",
        5: "SELL_STOP",
        6: "BUY_STOP_LIMIT",
        7: "SELL_STOP_LIMIT",
        8: "CLOSE_BY",
    }

    fmt.Printf("Margin rates for %s:\n", symbol)
    for orderType, name := range orderTypes {
        data, err := account.SymbolInfoMarginRate(ctx, &pb.SymbolInfoMarginRateRequest{
            Symbol:    symbol,
            OrderType: orderType,
        })
        if err != nil {
            continue
        }

        fmt.Printf("  %s: Initial=%.4f, Maintenance=%.4f\n",
            name, data.InitialMarginRate, data.MaintenanceMarginRate)
    }
}
```

---

## 📚 See Also

* [OrderCalcMargin](../4.%20Trading_Operations/OrderCalcMargin.md) - Calculate required margin
* [SymbolInfoDouble](../2.%20Symbol_information/SymbolInfoDouble.md) - Get symbol properties
* [AccountSummary](../1.%20Account_information/AccountSummary.md) - Get account leverage
