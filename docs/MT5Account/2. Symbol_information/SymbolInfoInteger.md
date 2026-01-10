# ✅ Get Symbol Integer Property

> **Request:** get integer property of trading symbol, such as decimal places count (digits), spread, stops level and other parameters.

**API Information:**

* **SDK wrapper:** `MT5Account.SymbolInfoInteger(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolInfoInteger` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolInfoInteger(SymbolInfoIntegerRequest) → SymbolInfoIntegerReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolInfoInteger(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

## 💬 Just the essentials

* **What it is.** Retrieves integer properties of a trading symbol.
* **Why you need it.** Get precision (digits), spread, stop levels, and trading mode information.
* **Property ID required.** Must specify which property to retrieve using PropertyId enum.

---

## 🎯 Purpose

Use it to:

* Get symbol decimal precision (digits)
* Check current spread in points
* Retrieve minimum stop-loss/take-profit distance
* Validate trading restrictions
* Determine order execution mode

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolInfoInteger - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoInteger_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolInfoInteger retrieves an integer-type symbol property.
// Use PropertyId to specify which property (DIGITS, SPREAD, STOPS_LEVEL, etc).
func (a *MT5Account) SymbolInfoInteger(
    ctx context.Context,
    req *pb.SymbolInfoIntegerRequest,
) (*pb.SymbolInfoIntegerData, error)
```

**Request message:**

```protobuf
SymbolInfoIntegerRequest {
  string Symbol = 1;       // Symbol name
  int32 PropertyId = 2;    // Property identifier (DIGITS, SPREAD, etc)
}
```

**Reply message:**

```protobuf
SymbolInfoIntegerReply {
  oneof response {
    SymbolInfoIntegerData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                             | Description                                   |
| --------- | -------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolInfoIntegerRequest`   | Request with Symbol and PropertyId            |

**Request fields:**

| Field        | Type     | Description                                                     |
| ------------ | -------- | --------------------------------------------------------------- |
| `Symbol`     | `string` | Symbol name (e.g., "EURUSD")                                    |
| `PropertyId` | `int32`  | Property identifier from SymbolInfoIntegerProperty enum         |

**PropertyId enum values (SymbolInfoIntegerProperty):**

| Constant | Value | Description |
|----------|-------|-------------|
| `SYMBOL_SUBSCRIPTION_DELAY` | 0 | Subscription delay in milliseconds |
| `SYMBOL_SECTOR` | 1 | Economic sector (enum SYMBOL_SECTOR) |
| `SYMBOL_INDUSTRY` | 2 | Industry (enum SYMBOL_INDUSTRY) |
| `SYMBOL_CUSTOM` | 3 | Custom symbol flag (0=standard, 1=custom) |
| `SYMBOL_BACKGROUND_COLOR` | 4 | Background color in Market Watch |
| `SYMBOL_CHART_MODE` | 5 | Price chart mode (SYMBOL_CHART_MODE enum) |
| `SYMBOL_EXIST` | 6 | Symbol exists (0=no, 1=yes) |
| `SYMBOL_SELECT` | 7 | Symbol selected in Market Watch (0=no, 1=yes) |
| `SYMBOL_VISIBLE` | 8 | Symbol visible in Market Watch (0=no, 1=yes) |
| `SYMBOL_SESSION_DEALS` | 9 | Number of deals in current session |
| `SYMBOL_SESSION_BUY_ORDERS` | 10 | Number of Buy orders in current session |
| `SYMBOL_SESSION_SELL_ORDERS` | 11 | Number of Sell orders in current session |
| `SYMBOL_VOLUME` | 12 | Last deal volume |
| `SYMBOL_VOLUMEHIGH` | 13 | Maximum volume of the day |
| `SYMBOL_VOLUMELOW` | 14 | Minimum volume of the day |
| `SYMBOL_TIME` | 15 | Last quote time (seconds since 1970.01.01) |
| `SYMBOL_TIME_MSC` | 16 | Last quote time in milliseconds |
| `SYMBOL_DIGITS` | 17 | Number of decimal places |
| `SYMBOL_SPREAD_FLOAT` | 18 | Floating spread flag (0=fixed, 1=floating) |
| `SYMBOL_SPREAD` | 19 | Current spread in points |
| `SYMBOL_TICKS_BOOKDEPTH` | 20 | Maximum depth of market (DOM) |
| `SYMBOL_TRADE_CALC_MODE` | 21 | Margin calculation mode (SYMBOL_CALC_MODE enum) |
| `SYMBOL_TRADE_MODE` | 22 | Trade execution mode (SYMBOL_TRADE_MODE enum) |
| `SYMBOL_START_TIME` | 23 | Date of symbol trade beginning |
| `SYMBOL_EXPIRATION_TIME` | 24 | Date of symbol trade end |
| `SYMBOL_TRADE_STOPS_LEVEL` | 25 | Minimum distance for SL/TP in points |
| `SYMBOL_TRADE_FREEZE_LEVEL` | 26 | Order freeze level in points |
| `SYMBOL_TRADE_EXEMODE` | 27 | Order execution type (SYMBOL_ORDER_EXEMODE enum) |
| `SYMBOL_SWAP_MODE` | 28 | Swap calculation mode (SYMBOL_SWAP_MODE enum) |
| `SYMBOL_SWAP_ROLLOVER3DAYS` | 29 | Day of week for triple swap (SYMBOL_SWAP_ROLLOVER enum) |
| `SYMBOL_MARGIN_HEDGED_USE_LEG` | 30 | Calculating hedged margin using larger leg |
| `SYMBOL_EXPIRATION_MODE` | 31 | Allowed order expiration modes |
| `SYMBOL_FILLING_MODE` | 32 | Allowed order filling modes |
| `SYMBOL_ORDER_MODE` | 33 | Allowed order types |
| `SYMBOL_ORDER_GTC_MODE` | 34 | StopLoss/TakeProfit mode (ORDER_GTC_MODE enum) |
| `SYMBOL_OPTION_MODE` | 35 | Option type (SYMBOL_OPTION_MODE enum) |
| `SYMBOL_OPTION_RIGHT` | 36 | Option right (Call/Put) |

---

## ⬆️ Output — `SymbolInfoIntegerData`

| Field   | Type    | Go Type | Description                            |
| ------- | ------- | ------- | -------------------------------------- |
| `Value` | `int64` | `int64` | The requested integer property value   |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Spread variability:** Spread can change dynamically during trading sessions.
* **Stop levels:** STOPS_LEVEL is critical for placing orders with SL/TP.

---

## 🔗 Usage Examples

### 1) Get symbol digits (decimal precision)

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

    data, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     "EURUSD",
        PropertyId: 17, // SYMBOL_DIGITS
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("EURUSD Digits: %d\n", data.Value)
}
```

### 2) Get current spread

```go
func GetSpreadPoints(account *mt5.MT5Account, symbol string) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     symbol,
        PropertyId: 18, // SYMBOL_SPREAD
    })
    if err != nil {
        return 0, fmt.Errorf("failed to get spread: %w", err)
    }

    return data.Value, nil
}

// Usage:
// spread, _ := GetSpreadPoints(account, "EURUSD")
// fmt.Printf("Spread: %d points\n", spread)
```

### 3) Check minimum stops level

```go
func GetStopsLevel(account *mt5.MT5Account, symbol string) (int64, error) {
    ctx := context.Background()

    data, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     symbol,
        PropertyId: 25, // SYMBOL_TRADE_STOPS_LEVEL
    })
    if err != nil {
        return 0, err
    }

    return data.Value, nil
}

// Usage for SL/TP validation:
// stopsLevel, _ := GetStopsLevel(account, "EURUSD")
// fmt.Printf("Minimum SL/TP distance: %d points\n", stopsLevel)
```

### 4) Validate stop-loss distance

```go
func ValidateStopLoss(account *mt5.MT5Account, symbol string, currentPrice, slPrice float64) error {
    ctx := context.Background()

    // Get stops level
    stopsData, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     symbol,
        PropertyId: 25, // SYMBOL_TRADE_STOPS_LEVEL
    })
    if err != nil {
        return err
    }

    // Get point size
    pointData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 13, // SYMBOL_POINT
    })
    if err != nil {
        return err
    }

    // Calculate distance in points
    distance := math.Abs(currentPrice-slPrice) / pointData.Value
    minDistance := float64(stopsData.Value)

    if distance < minDistance {
        return fmt.Errorf("SL too close: %.0f points, minimum: %.0f points", distance, minDistance)
    }

    return nil
}
```

### 5) Get trade execution mode

```go
func GetTradeMode(account *mt5.MT5Account, symbol string) (string, error) {
    ctx := context.Background()

    data, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     symbol,
        PropertyId: 22, // SYMBOL_TRADE_MODE
    })
    if err != nil {
        return "", err
    }

    modes := map[int64]string{
        0: "DISABLED",
        1: "LONGONLY",
        2: "SHORTONLY",
        3: "CLOSEONLY",
        4: "FULL",
    }

    mode, ok := modes[data.Value]
    if !ok {
        return "UNKNOWN", nil
    }

    return mode, nil
}
```

### 6) Check if trading is allowed

```go
func IsTradingAllowed(account *mt5.MT5Account, symbol string) (bool, error) {
    ctx := context.Background()

    data, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     symbol,
        PropertyId: 22, // SYMBOL_TRADE_MODE
    })
    if err != nil {
        return false, err
    }

    // 0 = disabled, 1 = long only, 2 = short only, 3 = close only, 4 = full
    return data.Value > 0, nil
}
```

---

## 🔧 Common Patterns

### Format price with correct digits

```go
func FormatPrice(account *mt5.MT5Account, symbol string, price float64) (string, error) {
    ctx := context.Background()

    data, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     symbol,
        PropertyId: 17, // SYMBOL_DIGITS
    })
    if err != nil {
        return "", err
    }

    format := fmt.Sprintf("%%.%df", data.Value)
    return fmt.Sprintf(format, price), nil
}

// Usage:
// formatted, _ := FormatPrice(account, "EURUSD", 1.12345678)
// fmt.Println(formatted) // "1.12346" (if digits=5)
```

### Calculate minimum SL/TP distance

```go
func CalculateMinStopDistance(account *mt5.MT5Account, symbol string) (float64, error) {
    ctx := context.Background()

    // Get stops level in points
    stopsData, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
        Symbol:     symbol,
        PropertyId: 25, // SYMBOL_TRADE_STOPS_LEVEL
    })
    if err != nil {
        return 0, err
    }

    // Get point size
    pointData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 13, // SYMBOL_POINT
    })
    if err != nil {
        return 0, err
    }

    // Distance in price units
    distance := float64(stopsData.Value) * pointData.Value
    return distance, nil
}
```

---

## 📚 See Also

* [SymbolInfoDouble](./SymbolInfoDouble.md) - Get double symbol properties
* [SymbolInfoString](./SymbolInfoString.md) - Get string symbol properties
* [SymbolParamsMany](../6.%20Additional_Methods/SymbolParamsMany.md) - Get all symbol parameters at once
* [OrderCheck](../4.%20Trading_Operations/OrderCheck.md) - Validate order parameters
