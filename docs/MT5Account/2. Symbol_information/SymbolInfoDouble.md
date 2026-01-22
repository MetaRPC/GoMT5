# ✅ Get Symbol Double Property

> **Request:** get numeric (double) property of trading symbol, such as BID, ASK, POINT, volumes, spreads and other trading parameters.

**API Information:**

* **Low-level API:** `MT5Account.SymbolInfoDouble(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolInfoDouble` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolInfoDouble(SymbolInfoDoubleRequest) → SymbolInfoDoubleReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolInfoDouble(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieves numeric (double) properties of a trading symbol.
* **Why you need it.** Get current prices, point size, volume limits, and other trading parameters.
* **Property ID required.** Must specify which property to retrieve using PropertyId enum.

---

## 🎯 Purpose

Use it to:

* Get current Bid/Ask prices
* Retrieve point size for price calculations
* Check minimum/maximum trading volumes
* Obtain tick size and value for order placement
* Calculate trading parameters

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolInfoDouble - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoDouble_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolInfoDouble retrieves a double-type symbol property.
// Use PropertyId to specify which property (BID, ASK, POINT, VOLUME_MIN, etc).
func (a *MT5Account) SymbolInfoDouble(
    ctx context.Context,
    req *pb.SymbolInfoDoubleRequest,
) (*pb.SymbolInfoDoubleData, error)
```

**Request message:**

```protobuf
SymbolInfoDoubleRequest {
  string Symbol = 1;       // Symbol name
  int32 PropertyId = 2;    // Property identifier (BID, ASK, POINT, etc)
}
```

**Reply message:**

```protobuf
SymbolInfoDoubleReply {
  oneof response {
    SymbolInfoDoubleData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                            | Description                                   |
| --------- | ------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`               | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolInfoDoubleRequest`   | Request with Symbol and PropertyId            |

**Request fields:**

| Field        | Type     | Description                                                  |
| ------------ | -------- | ------------------------------------------------------------ |
| `Symbol`     | `string` | Symbol name (e.g., "EURUSD")                                 |
| `PropertyId` | `int32`  | Property identifier from SymbolInfoDoubleProperty enum       |

**PropertyId enum values (SymbolInfoDoubleProperty):**

| Constant | Value | Description |
|----------|-------|-------------|
| `SYMBOL_BID` | 0 | Current Bid price |
| `SYMBOL_BIDHIGH` | 1 | Highest Bid of the day |
| `SYMBOL_BIDLOW` | 2 | Lowest Bid of the day |
| `SYMBOL_ASK` | 3 | Current Ask price |
| `SYMBOL_ASKHIGH` | 4 | Highest Ask of the day |
| `SYMBOL_ASKLOW` | 5 | Lowest Ask of the day |
| `SYMBOL_LAST` | 6 | Last deal price |
| `SYMBOL_LASTHIGH` | 7 | Highest Last price of the day |
| `SYMBOL_LASTLOW` | 8 | Lowest Last price of the day |
| `SYMBOL_VOLUME_REAL` | 9 | Last deal volume |
| `SYMBOL_VOLUMEHIGH_REAL` | 10 | Maximum volume of the day |
| `SYMBOL_VOLUMELOW_REAL` | 11 | Minimum volume of the day |
| `SYMBOL_OPTION_STRIKE` | 12 | Option strike price |
| `SYMBOL_POINT` | 13 | Point size (smallest price change) |
| `SYMBOL_TRADE_TICK_VALUE` | 14 | Tick value in account currency |
| `SYMBOL_TRADE_TICK_VALUE_PROFIT` | 15 | Tick value for profit position |
| `SYMBOL_TRADE_TICK_VALUE_LOSS` | 16 | Tick value for loss position |
| `SYMBOL_TRADE_TICK_SIZE` | 17 | Minimum price change |
| `SYMBOL_TRADE_CONTRACT_SIZE` | 18 | Contract size (lot size in base currency) |
| `SYMBOL_TRADE_ACCRUED_INTEREST` | 19 | Accrued interest |
| `SYMBOL_TRADE_FACE_VALUE` | 20 | Face value (for bonds) |
| `SYMBOL_TRADE_LIQUIDITY_RATE` | 21 | Liquidity rate |
| `SYMBOL_VOLUME_MIN` | 22 | Minimum volume for trading |
| `SYMBOL_VOLUME_MAX` | 23 | Maximum volume for trading |
| `SYMBOL_VOLUME_STEP` | 24 | Volume step (lot size increment) |
| `SYMBOL_VOLUME_LIMIT` | 25 | Maximum total volume of open positions |
| `SYMBOL_SWAP_LONG` | 26 | Swap for long positions |
| `SYMBOL_SWAP_SHORT` | 27 | Swap for short positions |
| `SYMBOL_SWAP_SUNDAY` | 28 | Triple swap day - Sunday |
| `SYMBOL_SWAP_MONDAY` | 29 | Triple swap day - Monday |
| `SYMBOL_SWAP_TUESDAY` | 30 | Triple swap day - Tuesday |
| `SYMBOL_SWAP_WEDNESDAY` | 31 | Triple swap day - Wednesday |
| `SYMBOL_SWAP_THURSDAY` | 32 | Triple swap day - Thursday |
| `SYMBOL_SWAP_FRIDAY` | 33 | Triple swap day - Friday |
| `SYMBOL_SWAP_SATURDAY` | 34 | Triple swap day - Saturday |
| `SYMBOL_MARGIN_INITIAL` | 35 | Initial margin |
| `SYMBOL_MARGIN_MAINTENANCE` | 36 | Maintenance margin |
| `SYMBOL_SESSION_VOLUME` | 37 | Summary volume of the current session |
| `SYMBOL_SESSION_TURNOVER` | 38 | Summary turnover of the current session |
| `SYMBOL_SESSION_INTEREST` | 39 | Summary open interest |
| `SYMBOL_SESSION_BUY_ORDERS_VOLUME` | 40 | Current volume of Buy orders |
| `SYMBOL_SESSION_SELL_ORDERS_VOLUME` | 41 | Current volume of Sell orders |
| `SYMBOL_SESSION_OPEN` | 42 | Open price of the current session |
| `SYMBOL_SESSION_CLOSE` | 43 | Close price of the current session |
| `SYMBOL_SESSION_AW` | 44 | Average weighted price of the current session |
| `SYMBOL_SESSION_PRICE_SETTLEMENT` | 45 | Settlement price of the current session |
| `SYMBOL_SESSION_PRICE_LIMIT_MIN` | 46 | Minimum price of the current session |
| `SYMBOL_SESSION_PRICE_LIMIT_MAX` | 47 | Maximum price of the current session |
| `SYMBOL_MARGIN_HEDGED` | 48 | Hedged margin |
| `SYMBOL_PRICE_CHANGE` | 49 | Change of price in % |
| `SYMBOL_PRICE_VOLATILITY` | 50 | Price volatility in % |
| `SYMBOL_PRICE_THEORETICAL` | 51 | Theoretical option price |
| `SYMBOL_PRICE_DELTA` | 52 | Option/warrant delta |
| `SYMBOL_PRICE_THETA` | 53 | Option/warrant theta |
| `SYMBOL_PRICE_GAMMA` | 54 | Option/warrant gamma |
| `SYMBOL_PRICE_VEGA` | 55 | Option/warrant vega |
| `SYMBOL_PRICE_RHO` | 56 | Option/warrant rho |
| `SYMBOL_PRICE_OMEGA` | 57 | Option/warrant omega |
| `SYMBOL_PRICE_SENSITIVITY` | 58 | Option/warrant sensitivity |
| `SYMBOL_COUNT` | 59 | Total number of properties (not used) |

---

## ⬆️ Output — `SymbolInfoDoubleData`

| Field   | Type     | Go Type  | Description                              |
| ------- | -------- | -------- | ---------------------------------------- |
| `Value` | `double` | `float64` | The requested double property value     |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Synchronization:** Symbol should be synchronized for accurate real-time prices.
* **Performance:** For multiple properties, consider using SymbolParamsMany instead.

---

## 🔗 Usage Examples

### 1) Get current Bid price

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

    data, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     "EURUSD",
        PropertyId: 0, // SYMBOL_BID
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("EURUSD Bid: %.5f\n", data.Value)
}
```

### 2) Get Bid and Ask spread

```go
func GetSpread(account *mt5.MT5Account, symbol string) (float64, error) {
    ctx := context.Background()

    // Get Bid
    bidData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 0, // SYMBOL_BID
    })
    if err != nil {
        return 0, err
    }

    // Get Ask
    askData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 3, // SYMBOL_ASK
    })
    if err != nil {
        return 0, err
    }

    spread := askData.Value - bidData.Value
    return spread, nil
}
```

### 3) Get point size for calculations

```go
func GetPointSize(account *mt5.MT5Account, symbol string) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 13, // SYMBOL_POINT
    })
    if err != nil {
        return 0, fmt.Errorf("failed to get point: %w", err)
    }

    return data.Value, nil
}

// Usage for stop-loss calculation:
// point, _ := GetPointSize(account, "EURUSD")
// slDistance := 50 * point  // 50 pips
```

### 4) Check minimum trading volume

```go
func GetVolumeConstraints(account *mt5.MT5Account, symbol string) (min, max, step float64, err error) {
    ctx := context.Background()

    // Get minimum volume
    minData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 22, // SYMBOL_VOLUME_MIN
    })
    if err != nil {
        return 0, 0, 0, err
    }

    // Get maximum volume
    maxData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 23, // SYMBOL_VOLUME_MAX
    })
    if err != nil {
        return 0, 0, 0, err
    }

    // Get volume step
    stepData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 24, // SYMBOL_VOLUME_STEP
    })
    if err != nil {
        return 0, 0, 0, err
    }

    return minData.Value, maxData.Value, stepData.Value, nil
}
```

### 5) Calculate pip value

```go
func CalculatePipValue(account *mt5.MT5Account, symbol string) (float64, error) {
    ctx := context.Background()

    // Get tick size
    tickSizeData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 17, // SYMBOL_TRADE_TICK_SIZE
    })
    if err != nil {
        return 0, err
    }

    // Get tick value
    tickValueData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 14, // SYMBOL_TRADE_TICK_VALUE
    })
    if err != nil {
        return 0, err
    }

    // Get point
    pointData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: 13, // SYMBOL_POINT
    })
    if err != nil {
        return 0, err
    }

    // Pip value = (tick value / tick size) * point
    pipValue := (tickValueData.Value / tickSizeData.Value) * pointData.Value
    return pipValue, nil
}
```

### 6) Monitor price changes

```go
func MonitorPrice(account *mt5.MT5Account, symbol string, interval time.Duration) {
    ctx := context.Background()
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    var lastBid float64

    for range ticker.C {
        data, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
            Symbol:     symbol,
            PropertyId: 0, // SYMBOL_BID
        })

        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        if lastBid != 0 {
            change := data.Value - lastBid
            fmt.Printf("[%s] %s: %.5f (change: %+.5f)\n",
                time.Now().Format("15:04:05"),
                symbol,
                data.Value,
                change,
            )
        }

        lastBid = data.Value
    }
}
```

---

## 🔧 Common Patterns

### Get current price with validation

```go
func GetCurrentPrice(account *mt5.MT5Account, symbol string, isBid bool) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    propertyId := int32(1) // BID
    if !isBid {
        propertyId = 2 // ASK
    }

    data, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
        Symbol:     symbol,
        PropertyId: propertyId,
    })

    if err != nil {
        return 0, fmt.Errorf("failed to get price: %w", err)
    }

    if data.Value <= 0 {
        return 0, fmt.Errorf("invalid price: %.5f", data.Value)
    }

    return data.Value, nil
}
```

### Validate lot size

```go
func ValidateLotSize(account *mt5.MT5Account, symbol string, volume float64) (float64, error) {
    ctx := context.Background()

    min, max, step, err := GetVolumeConstraints(account, symbol)
    if err != nil {
        return 0, err
    }

    if volume < min {
        return min, nil
    }
    if volume > max {
        return max, nil
    }

    // Round to step
    adjusted := math.Round(volume/step) * step
    return adjusted, nil
}
```

---

## 📚 See Also

* [SymbolInfoInteger](./SymbolInfoInteger.md) - Get integer symbol properties
* [SymbolInfoString](./SymbolInfoString.md) - Get string symbol properties
* [SymbolParamsMany](../6.%20Additional_Methods/SymbolParamsMany.md) - Get all symbol parameters at once
* [SymbolInfoTick](./SymbolInfoTick.md) - Get tick data with prices
