# ✅ OrderCalcMargin

> **Request:** Calculate required margin for an order before placing it.

**API Information:**

* **SDK wrapper:** `MT5Account.OrderCalcMargin(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.TradeFunctions`
* **Proto definition:** `OrderCalcMargin` (defined in `mt5-term-api-trade-functions.proto`)

## 💬 Just the essentials

* **What it is.** Pre-calculates margin requirement before order placement without actually placing the order.
* **Why you need it.** Verify sufficient account margin before trading to avoid "Not enough money" rejections.
* **Risk management.** Essential for position sizing and account safety - calculate maximum safe lot size based on available margin.

---

## 🎯 Purpose

Use it to:

* Check margin requirements before placing orders
* Calculate maximum position size based on available free margin
* Validate trading parameters in risk management systems
* Display margin requirements in trading UI/dashboard
* Prevent order rejections due to insufficient margin

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OrderCalcMargin - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderCalcMargin_HOW.md)**

---

### RPC

```go
// OrderCalcMargin calculates required margin for an order.
//
// Use this method to determine how much margin will be required before placing an order.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderCalcMarginRequest with Action, Symbol, Volume, and Price
//
// Returns OrderCalcMarginData with Margin value in account currency.
func (a *MT5Account) OrderCalcMargin(
    ctx context.Context,
    req *pb.OrderCalcMarginRequest,
) (*pb.OrderCalcMarginData, error)
```

**Request message:** `OrderCalcMarginRequest`
**Reply message:** `OrderCalcMarginReply { oneof response { OrderCalcMarginData data = 1; Error error = 2; } }`

---

## 🔽 Input

| Parameter | Type                          | Description                                             |
| --------- | ----------------------------- | ------------------------------------------------------- |
| `ctx`     | `context.Context`             | Context for deadline/timeout and cancellation           |
| `req`     | `*pb.OrderCalcMarginRequest`  | Request with Action, Symbol, Volume, and Price          |

**OrderCalcMarginRequest fields:**

| Field       | Type                   | Required | Description                                                       |
| ----------- | ---------------------- | -------- | ----------------------------------------------------------------- |
| `Symbol`    | `string`               | ✅       | Trading instrument name (e.g., "EURUSD")                          |
| `OrderType` | `ENUM_ORDER_TYPE_TF`   | ✅       | Order type: ORDER_TYPE_TF_BUY (0), ORDER_TYPE_TF_SELL (1), etc    |
| `Volume`    | `double`               | ✅       | Trade volume in lots                                              |
| `OpenPrice` | `double`               | ✅       | Order price (use current market price for market orders)          |

---

## ⬆️ Output — `OrderCalcMarginData`

| Field    | Type     | Description                                              |
| -------- | -------- | -------------------------------------------------------- |
| `Margin` | `double` | Required margin in account currency                      |

---

### 📘 Enum: ENUM_ORDER_TYPE_TF

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `ORDER_TYPE_TF_BUY` | Market Buy order |
| 1 | `ORDER_TYPE_TF_SELL` | Market Sell order |
| 2 | `ORDER_TYPE_TF_BUY_LIMIT` | Buy Limit pending order |
| 3 | `ORDER_TYPE_TF_SELL_LIMIT` | Sell Limit pending order |
| 4 | `ORDER_TYPE_TF_BUY_STOP` | Buy Stop pending order |
| 5 | `ORDER_TYPE_TF_SELL_STOP` | Sell Stop pending order |
| 6 | `ORDER_TYPE_TF_BUY_STOP_LIMIT` | Upon reaching the order price, a pending Buy Limit order is placed at the StopLimit price |
| 7 | `ORDER_TYPE_TF_SELL_STOP_LIMIT` | Upon reaching the order price, a pending Sell Limit order is placed at the StopLimit price |
| 8 | `ORDER_TYPE_TF_CLOSE_BY` | Order to close a position by an opposite one |

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `5s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Price parameter:** For market orders, use current bid/ask price. For pending orders, use the intended order price.
* **Hedging vs Netting:** Margin calculation depends on account margin calculation mode (hedging or netting).
* **Leverage consideration:** Result accounts for account leverage and symbol-specific margin requirements.

---

## 🔗 Usage Examples

### 1) Basic margin calculation for market order

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
    account := mt5.NewMT5Account("user", "password", "server:443")

    err := account.Connect()
    if err != nil {
        panic(err)
    }
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Calculate margin for buying 0.1 lots of EURUSD at current price
    marginData, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    "EURUSD",
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    0.1,
        OpenPrice: 1.10000, // Current ask price
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Required margin: %.2f USD\n", marginData.Margin)
    // Output: Required margin: 110.00 USD
}
```

### 2) Check if account has enough margin

```go
func CanOpenPosition(account *mt5.MT5Account, symbol string, volume float64, price float64) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Get account summary to check free margin
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        return false, err
    }

    freeMargin := summary.AccountEquity - summary.AccountBalance // Simplified

    // Calculate required margin
    marginData, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    volume,
        OpenPrice: price,
    })
    if err != nil {
        return false, err
    }

    fmt.Printf("Free margin: %.2f\n", freeMargin)
    fmt.Printf("Required margin: %.2f\n", marginData.Margin)

    return freeMargin >= marginData.Margin, nil
}
```

### 3) Calculate maximum safe lot size

```go
func MaxSafeLotSize(account *mt5.MT5Account, symbol string, price float64, maxMarginUsagePercent float64) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Get account equity
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        return 0, err
    }

    maxAllowedMargin := summary.AccountEquity * (maxMarginUsagePercent / 100.0)

    // Test with 1.0 lot to get margin per lot
    testMargin, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    1.0,
        OpenPrice: price,
    })
    if err != nil {
        return 0, err
    }

    marginPerLot := testMargin.Margin
    maxLots := maxAllowedMargin / marginPerLot

    fmt.Printf("Account equity: %.2f\n", summary.AccountEquity)
    fmt.Printf("Max allowed margin (%.0f%%): %.2f\n", maxMarginUsagePercent, maxAllowedMargin)
    fmt.Printf("Margin per 1.0 lot: %.2f\n", marginPerLot)
    fmt.Printf("Max safe lot size: %.2f\n", maxLots)

    return maxLots, nil
}

// Usage:
// maxLots, err := MaxSafeLotSize(account, "EURUSD", 1.10000, 20.0) // Use max 20% of equity
```

### 4) Compare margin for buy vs sell

```go
func CompareMarginBuySell(account *mt5.MT5Account, symbol string, volume float64, bidPrice, askPrice float64) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Calculate margin for BUY
    buyMargin, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    volume,
        OpenPrice: askPrice,
    })
    if err != nil {
        fmt.Printf("Buy margin error: %v\n", err)
        return
    }

    // Calculate margin for SELL
    sellMargin, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL,
        Volume:    volume,
        OpenPrice: bidPrice,
    })
    if err != nil {
        fmt.Printf("Sell margin error: %v\n", err)
        return
    }

    fmt.Printf("Symbol: %s, Volume: %.2f\n", symbol, volume)
    fmt.Printf("BUY margin (at %.5f): %.2f\n", askPrice, buyMargin.Margin)
    fmt.Printf("SELL margin (at %.5f): %.2f\n", bidPrice, sellMargin.Margin)
    fmt.Printf("Difference: %.2f\n", buyMargin.Margin-sellMargin.Margin)
}
```

### 5) Batch margin calculation for multiple symbols

```go
func CalculateMarginForPortfolio(account *mt5.MT5Account, positions []struct {
    Symbol string
    Volume float64
    Price  float64
}) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    totalMargin := 0.0

    for _, pos := range positions {
        marginData, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
            Symbol:    pos.Symbol,
            OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
            Volume:    pos.Volume,
            OpenPrice: pos.Price,
        })
        if err != nil {
            fmt.Printf("Error for %s: %v\n", pos.Symbol, err)
            continue
        }

        fmt.Printf("%s (%.2f lots): %.2f\n", pos.Symbol, pos.Volume, marginData.Margin)
        totalMargin += marginData.Margin
    }

    fmt.Printf("\nTotal required margin: %.2f\n", totalMargin)
}

// Usage:
// CalculateMarginForPortfolio(account, []struct {
//     Symbol string
//     Volume float64
//     Price  float64
// }{
//     {"EURUSD", 0.1, 1.10000},
//     {"GBPUSD", 0.05, 1.27000},
//     {"USDJPY", 0.08, 150.500},
// })
```

### 6) Margin calculation with error handling

```go
func SafeMarginCheck(account *mt5.MT5Account, symbol string, volume float64, price float64) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    marginData, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    volume,
        OpenPrice: price,
    })

    if err != nil {
        // Handle specific errors
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Request timeout - server may be slow")
        } else {
            fmt.Printf("Margin calculation failed: %v\n", err)
        }
        return
    }

    // Validate result
    if marginData.Margin <= 0 {
        fmt.Printf("Warning: Invalid margin value: %.2f\n", marginData.Margin)
        return
    }

    fmt.Printf("Required margin: %.2f\n", marginData.Margin)

    // Get account info for comparison
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        fmt.Printf("Failed to get account info: %v\n", err)
        return
    }

    marginLevel := (summary.AccountEquity / marginData.Margin) * 100.0
    fmt.Printf("Margin level after order: %.2f%%\n", marginLevel)

    if marginLevel < 200 {
        fmt.Println("Warning: Low margin level!")
    }
}
```

---

## 🔧 Common Patterns

### Pre-trade validation

```go
func ValidateTradeMargin(account *mt5.MT5Account, symbol string, volume float64, price float64, minMarginLevel float64) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Get required margin
    marginData, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    volume,
        OpenPrice: price,
    })
    if err != nil {
        return fmt.Errorf("margin calculation failed: %w", err)
    }

    // Get account state
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        return fmt.Errorf("account summary failed: %w", err)
    }

    // Calculate margin level after trade
    newMarginLevel := (summary.AccountEquity / marginData.Margin) * 100.0

    if newMarginLevel < minMarginLevel {
        return fmt.Errorf("insufficient margin: level would be %.2f%% (min: %.2f%%)",
            newMarginLevel, minMarginLevel)
    }

    return nil
}
```

### Dynamic lot size calculator

```go
func CalculateLotSizeByRisk(account *mt5.MT5Account, symbol string, price float64, riskPercent float64) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        return 0, err
    }

    maxRiskAmount := summary.AccountEquity * (riskPercent / 100.0)

    // Test with small volume to get margin requirement
    testMargin, err := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
        Symbol:    symbol,
        OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:    0.01,
        OpenPrice: price,
    })
    if err != nil {
        return 0, err
    }

    marginPer001Lot := testMargin.Margin
    lots := (maxRiskAmount / marginPer001Lot) * 0.01

    return lots, nil
}
```

---

## 📚 See Also

* [OrderCalcProfit](./OrderCalcProfit.md) - Calculate potential profit for a trade
* [OrderCheck](./OrderCheck.md) - Validate complete order before sending
* [OrderSend](./OrderSend.md) - Place market or pending order
* [AccountSummary](../1.%20Account_information/AccountSummary.md) - Get account balance and equity
* [SymbolInfoMarginRate](../6.%20Additional_Methods/SymbolInfoMarginRate.md) - Get margin requirements for order types
