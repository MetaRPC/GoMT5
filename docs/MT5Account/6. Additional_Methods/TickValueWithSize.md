# ✅ Get Tick Value and Size Information

> **Request:** get comprehensive tick value and size data for multiple symbols at once.

**API Information:**

* **Low-level API:** `MT5Account.TickValueWithSize(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `TickValueWithSize` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `TickValueWithSize(TickValueWithSizeRequest) -> TickValueWithSizeReply`
* **Low-level client (generated):** `AccountHelperClient.TickValueWithSize(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Returns tick value and size data for symbols including profit/loss calculations.
* **Why you need it.** Calculate position value, P&L, risk management, lot size calculations.
* **Key data.** Provides TradeTickValue, TradeTickValueProfit, TradeTickValueLoss, TradeTickSize, TradeContractSize.

---

## 🎯 Purpose

Use it to:

* Calculate exact position value in account currency
* Determine profit/loss per tick for risk management
* Calculate required margin and position sizes
* Build accurate P&L calculators
* Understand tick size for precise price movements
* Get contract size for position value calculations

---

```go
package mt5

type MT5Account struct {
    // ...
}

// TickValueWithSize returns tick value and size information for multiple symbols.
// This method retrieves comprehensive tick value and size data including profit/loss calculations.
func (a *MT5Account) TickValueWithSize(
    ctx context.Context,
    req *pb.TickValueWithSizeRequest,
) (*pb.TickValueWithSizeData, error)
```

**Request message:**

```protobuf
TickValueWithSizeRequest {
  repeated string SymbolNames = 1;  // Array of trading symbols
}
```

---

## 🔽 Input

| Parameter | Type                                | Description                                   |
| --------- | ----------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                   | Context for deadline/timeout and cancellation |
| `req`     | `*pb.TickValueWithSizeRequest`      | Request with array of symbol names            |

**Request fields:**

| Field         | Type       | Description                                   |
| ------------- | ---------- | --------------------------------------------- |
| `SymbolNames` | `[]string` | Array of trading symbols (e.g., ["EURUSD", "GBPUSD"]) |

---

## ⬆️ Output — TickValueWithSizeData

| Field                  | Type                    | Description                                    |
| ---------------------- | ----------------------- | ---------------------------------------------- |
| `SymbolTickSizeInfos`  | `[]*TickSizeSymbol`     | Array of tick value/size info for each symbol |

### TickSizeSymbol structure:

| Field                  | Type      | Description                                           |
| ---------------------- | --------- | ----------------------------------------------------- |
| `Index`                | `int32`   | Symbol index in request array                         |
| `Name`                 | `string`  | Symbol name                                           |
| `TradeTickValue`       | `float64` | Standard tick value in account currency               |
| `TradeTickValueProfit` | `float64` | Tick value for profitable positions                   |
| `TradeTickValueLoss`   | `float64` | Tick value for losing positions                       |
| `TradeTickSize`        | `float64` | Minimum price change (tick size)                      |
| `TradeContractSize`    | `float64` | Contract size for the symbol                          |

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**-> [TickValueWithSize - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/TickValueWithSize_HOW.md)**

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `5s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Multiple symbols:** Request can include multiple symbols in a single call for efficiency.
* **Account currency:** All tick values are returned in your account currency.
* **Profit vs Loss:** Some symbols have different tick values for profitable vs losing positions (asymmetric pricing).
* **Contract size:** Use this with tick value to calculate position value: `positionValue = lots * contractSize * price`.
* **Risk calculation:** `riskPerTick = lots * tickValue`, multiply by expected price movement in ticks.

---

## 🔗 Usage Examples

### 1) Get tick value for single symbol

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
    "github.com/google/uuid"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    data, err := account.TickValueWithSize(ctx, &pb.TickValueWithSizeRequest{
        SymbolNames: []string{"EURUSD"},
    })
    if err != nil {
        panic(err)
    }

    for _, info := range data.SymbolTickSizeInfos {
        fmt.Printf("Symbol: %s\n", info.Name)
        fmt.Printf("  Tick Value:        %.5f\n", info.TradeTickValue)
        fmt.Printf("  Tick Size:         %.5f\n", info.TradeTickSize)
        fmt.Printf("  Contract Size:     %.2f\n", info.TradeContractSize)
    }
}
```

### 2) Calculate position risk for multiple symbols

```go
func CalculatePositionRisk(account *mt5.MT5Account, symbols []string, lots float64) {
    ctx := context.Background()

    data, err := account.TickValueWithSize(ctx, &pb.TickValueWithSizeRequest{
        SymbolNames: symbols,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Risk per tick for %.2f lots:\n", lots)
    for _, info := range data.SymbolTickSizeInfos {
        riskPerTick := lots * info.TradeTickValue
        fmt.Printf("  %s: %.2f %s per tick\n", info.Name, riskPerTick, "account currency")
    }
}

// Usage:
// CalculatePositionRisk(account, []string{"EURUSD", "GBPUSD", "USDJPY"}, 0.1)
```

### 3) Calculate position value

```go
func CalculatePositionValue(account *mt5.MT5Account, symbol string, lots float64, price float64) (float64, error) {
    ctx := context.Background()

    data, err := account.TickValueWithSize(ctx, &pb.TickValueWithSizeRequest{
        SymbolNames: []string{symbol},
    })
    if err != nil {
        return 0, err
    }

    if len(data.SymbolTickSizeInfos) == 0 {
        return 0, fmt.Errorf("no data for symbol %s", symbol)
    }

    info := data.SymbolTickSizeInfos[0]
    positionValue := lots * info.TradeContractSize * price

    fmt.Printf("Position value for %s:\n", symbol)
    fmt.Printf("  Lots:          %.2f\n", lots)
    fmt.Printf("  Contract Size: %.2f\n", info.TradeContractSize)
    fmt.Printf("  Price:         %.5f\n", price)
    fmt.Printf("  Total Value:   %.2f\n", positionValue)

    return positionValue, nil
}
```

### 4) Calculate P&L per pip

```go
func CalculateProfitPerPip(account *mt5.MT5Account, symbol string, lots float64) {
    ctx := context.Background()

    data, err := account.TickValueWithSize(ctx, &pb.TickValueWithSizeRequest{
        SymbolNames: []string{symbol},
    })
    if err != nil {
        panic(err)
    }

    info := data.SymbolTickSizeInfos[0]

    // Calculate ticks per pip (usually 10 for 5-digit quotes)
    ticksPerPip := 0.0001 / info.TradeTickSize
    profitPerPip := lots * info.TradeTickValueProfit * ticksPerPip
    lossPerPip := lots * info.TradeTickValueLoss * ticksPerPip

    fmt.Printf("%s (%.2f lots):\n", symbol, lots)
    fmt.Printf("  Tick Size:        %.5f\n", info.TradeTickSize)
    fmt.Printf("  Ticks per pip:    %.0f\n", ticksPerPip)
    fmt.Printf("  Profit per pip:   %.2f\n", profitPerPip)
    fmt.Printf("  Loss per pip:     %.2f\n", lossPerPip)
}
```

### 5) Build position risk calculator

```go
func RiskCalculator(account *mt5.MT5Account, symbol string, accountBalance float64, riskPercent float64, stopLossPips float64) {
    ctx := context.Background()

    data, err := account.TickValueWithSize(ctx, &pb.TickValueWithSizeRequest{
        SymbolNames: []string{symbol},
    })
    if err != nil {
        panic(err)
    }

    info := data.SymbolTickSizeInfos[0]

    // Calculate maximum risk in account currency
    maxRisk := accountBalance * (riskPercent / 100.0)

    // Calculate value per pip for 1 lot
    ticksPerPip := 0.0001 / info.TradeTickSize
    valuePerPipOneLot := info.TradeTickValueLoss * ticksPerPip

    // Calculate position size
    lots := maxRisk / (stopLossPips * valuePerPipOneLot)

    fmt.Printf("Risk Calculator for %s:\n", symbol)
    fmt.Printf("  Account Balance:     %.2f\n", accountBalance)
    fmt.Printf("  Risk Percent:        %.2f%%\n", riskPercent)
    fmt.Printf("  Max Risk Amount:     %.2f\n", maxRisk)
    fmt.Printf("  Stop Loss:           %.1f pips\n", stopLossPips)
    fmt.Printf("  Value per pip/lot:   %.2f\n", valuePerPipOneLot)
    fmt.Printf("  Recommended Lots:    %.2f\n", lots)
}

// Usage:
// RiskCalculator(account, "EURUSD", 10000, 2.0, 50) // Risk 2% with 50 pip stop
```

---

## 📚 See Also

* [SymbolInfoDouble](../2.%20Symbol_information/SymbolInfoDouble.md) - Get individual symbol properties
* [SymbolParamsMany](./SymbolParamsMany.md) - Get comprehensive symbol parameters
* [OrderCalcMargin](../4.%20Trading_Operations/OrderCalcMargin.md) - Calculate required margin
* [OrderCalcProfit](../4.%20Trading_Operations/OrderCalcProfit.md) - Calculate potential profit
