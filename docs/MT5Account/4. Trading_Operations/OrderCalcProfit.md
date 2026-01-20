# ✅ OrderCalcProfit

> **Request:** Calculate potential profit/loss for a trade at specified entry and exit prices.

**API Information:**

* **SDK wrapper:** `MT5Account.OrderCalcProfit(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.TradeFunctions`
* **Proto definition:** `OrderCalcProfit` (defined in `mt5-term-api-trade-functions.proto`)

---

## 💬 Just the essentials

* **What it is.** Pre-calculates profit/loss for a hypothetical trade without placing the order.
* **Why you need it.** Estimate potential profit before trading, calculate risk/reward ratio, validate trading strategy.
* **Risk assessment.** Essential for position sizing and stop-loss/take-profit planning.

---

## 🎯 Purpose

Use it to:

* Calculate potential profit before placing orders
* Estimate profit at different exit price levels (TP levels)
* Calculate risk/reward ratios for trading strategies
* Validate stop-loss and take-profit levels
* Display projected P/L in trading UI
* Backtest strategy profit calculations

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OrderCalcProfit - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderCalcProfit_HOW.md)**

---

### RPC

```go
// OrderCalcProfit calculates potential profit for a trade.
//
// Use this method to estimate profit/loss before placing an order or to calculate
// current profit at a specified price level.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderCalcProfitRequest with Action, Symbol, Volume, PriceOpen, and PriceClose
//
// Returns OrderCalcProfitData with Profit value in account currency.
func (a *MT5Account) OrderCalcProfit(
    ctx context.Context,
    req *pb.OrderCalcProfitRequest,
) (*pb.OrderCalcProfitData, error)
```

**Request message:** `OrderCalcProfitRequest`

**Reply message:** `OrderCalcProfitReply { oneof response { OrderCalcProfitData data = 1; Error error = 2; } }`

---

## 🔽 Input

| Parameter | Type                          | Description                                             |
| --------- | ----------------------------- | ------------------------------------------------------- |
| `ctx`     | `context.Context`             | Context for deadline/timeout and cancellation           |
| `req`     | `*pb.OrderCalcProfitRequest`  | Request with Action, Symbol, Volume, PriceOpen, PriceClose |

**OrderCalcProfitRequest fields:**

| Field        | Type                   | Required | Description                                                  |
| ------------ | ---------------------- | -------- | ------------------------------------------------------------ |
| `OrderType`  | `ENUM_ORDER_TYPE_TF`   | ✅       | Order type: ORDER_TYPE_TF_BUY (0), ORDER_TYPE_TF_SELL (1), etc |
| `Symbol`     | `string`               | ✅       | Trading instrument name (e.g., "EURUSD")                     |
| `Volume`     | `double`               | ✅       | Trade volume in lots                                         |
| `OpenPrice`  | `double`               | ✅       | Entry price (buy at Ask, sell at Bid)                        |
| `ClosePrice` | `double`               | ✅       | Exit price (close buy at Bid, close sell at Ask)             |

---

## ⬆️ Output — `OrderCalcProfitData`

| Field    | Type     | Description                                              |
| -------- | -------- | -------------------------------------------------------- |
| `Profit` | `double` | Calculated profit/loss in account currency (negative for loss) |

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
* **Negative profit:** Result is negative when trade would result in a loss.
* **Currency conversion:** Profit is automatically converted to account currency.
* **Spread consideration:** Remember to account for spread when calculating entry/exit prices.

---

## 🔗 Usage Examples

### 1) Basic profit calculation

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
    account := mt5.NewMT5Account("user", "password", "server:443")

    err := account.Connect()
    if err != nil {
        panic(err)
    }
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Calculate profit for buying 0.1 lots EURUSD
    // Entry: 1.10000, Exit: 1.10500 (50 pips profit)
    profitData, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
        OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Symbol:     "EURUSD",
        Volume:     0.1,
        OpenPrice:  1.10000,
        ClosePrice: 1.10500,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Profit for 50 pips: %.2f USD\n", profitData.Profit)
    // Output: Profit for 50 pips: 50.00 USD
}
```

### 2) Calculate risk/reward ratio

```go
func CalculateRiskReward(account *mt5.MT5Account, symbol string, volume float64, entry, stopLoss, takeProfit float64, isBuy bool) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    orderType := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
    if !isBuy {
        orderType = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
    }

    // Calculate risk (loss if stop-loss hit)
    riskData, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
        OrderType:  orderType,
        Symbol:     symbol,
        Volume:     volume,
        OpenPrice:  entry,
        ClosePrice: stopLoss,
    })
    if err != nil {
        fmt.Printf("Risk calculation error: %v\n", err)
        return
    }

    // Calculate reward (profit if take-profit hit)
    rewardData, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
        OrderType:  orderType,
        Symbol:     symbol,
        Volume:     volume,
        OpenPrice:  entry,
        ClosePrice: takeProfit,
    })
    if err != nil {
        fmt.Printf("Reward calculation error: %v\n", err)
        return
    }

    risk := -riskData.Profit // Make risk positive
    reward := rewardData.Profit
    ratio := reward / risk

    fmt.Printf("Entry: %.5f, SL: %.5f, TP: %.5f\n", entry, stopLoss, takeProfit)
    fmt.Printf("Risk: %.2f, Reward: %.2f\n", risk, reward)
    fmt.Printf("Risk/Reward ratio: 1:%.2f\n", ratio)
}

// Usage:
// CalculateRiskReward(account, "EURUSD", 0.1, 1.10000, 1.09500, 1.11000, true)
// Output:
// Entry: 1.10000, SL: 1.09500, TP: 1.11000
// Risk: 50.00, Reward: 100.00
// Risk/Reward ratio: 1:2.00
```

### 3) Calculate profit at multiple take-profit levels

```go
func CalculateProfitLevels(account *mt5.MT5Account, symbol string, volume float64, entry float64, isBuy bool) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    orderType := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
    if !isBuy {
        orderType = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
    }

    // Calculate profit at different pip levels
    pipLevels := []int{10, 20, 50, 100, 200}
    pipValue := 0.0001 // For EURUSD

    fmt.Printf("Profit levels for %s %.2f lots from %.5f:\n", symbol, volume, entry)

    for _, pips := range pipLevels {
        var exitPrice float64
        if isBuy {
            exitPrice = entry + float64(pips)*pipValue
        } else {
            exitPrice = entry - float64(pips)*pipValue
        }

        profitData, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
            OrderType:  orderType,
            Symbol:     symbol,
            Volume:     volume,
            OpenPrice:  entry,
            ClosePrice: exitPrice,
        })
        if err != nil {
            fmt.Printf("  %d pips: Error - %v\n", pips, err)
            continue
        }

        fmt.Printf("  %d pips (%.5f): %.2f USD\n", pips, exitPrice, profitData.Profit)
    }
}

// Usage:
// CalculateProfitLevels(account, "EURUSD", 0.1, 1.10000, true)
```

### 4) Validate stop-loss placement

```go
func ValidateStopLoss(account *mt5.MT5Account, symbol string, volume float64, entry, stopLoss float64, isBuy bool, maxRiskPercent float64) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Get account balance
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        return false, err
    }

    maxRiskAmount := summary.AccountBalance * (maxRiskPercent / 100.0)

    // Calculate potential loss
    orderType := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
    if !isBuy {
        orderType = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
    }

    lossData, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
        OrderType:  orderType,
        Symbol:     symbol,
        Volume:     volume,
        OpenPrice:  entry,
        ClosePrice: stopLoss,
    })
    if err != nil {
        return false, err
    }

    potentialLoss := -lossData.Profit // Make positive
    riskPercent := (potentialLoss / summary.AccountBalance) * 100.0

    fmt.Printf("Account balance: %.2f\n", summary.AccountBalance)
    fmt.Printf("Max risk (%.1f%%): %.2f\n", maxRiskPercent, maxRiskAmount)
    fmt.Printf("Stop-loss risk: %.2f (%.2f%%)\n", potentialLoss, riskPercent)

    isValid := potentialLoss <= maxRiskAmount

    if !isValid {
        fmt.Printf("WARNING: Stop-loss risk exceeds maximum allowed!\n")
    }

    return isValid, nil
}
```

### 5) Compare profit for different lot sizes

```go
func CompareProfitByLotSize(account *mt5.MT5Account, symbol string, entry, exit float64, isBuy bool) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    orderType := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
    if !isBuy {
        orderType = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
    }

    lotSizes := []float64{0.01, 0.1, 0.5, 1.0, 5.0}

    fmt.Printf("Profit comparison for %s (%.5f -> %.5f):\n", symbol, entry, exit)

    for _, lots := range lotSizes {
        profitData, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
            OrderType:  orderType,
            Symbol:     symbol,
            Volume:     lots,
            OpenPrice:  entry,
            ClosePrice: exit,
        })
        if err != nil {
            fmt.Printf("  %.2f lots: Error - %v\n", lots, err)
            continue
        }

        fmt.Printf("  %.2f lots: %.2f USD\n", lots, profitData.Profit)
    }
}
```

### 6) Calculate break-even price

```go
func CalculateBreakEvenPrice(account *mt5.MT5Account, symbol string, volume float64, entry float64, isBuy bool, commission float64) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    orderType := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
    if !isBuy {
        orderType = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
    }

    // Test different exit prices to find break-even
    pipValue := 0.0001
    testRange := 100 // Test 100 pips range

    for i := 0; i <= testRange; i++ {
        var testPrice float64
        if isBuy {
            testPrice = entry + float64(i)*pipValue
        } else {
            testPrice = entry - float64(i)*pipValue
        }

        profitData, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
            OrderType:  orderType,
            Symbol:     symbol,
            Volume:     volume,
            OpenPrice:  entry,
            ClosePrice: testPrice,
        })
        if err != nil {
            continue
        }

        // Account for commission
        netProfit := profitData.Profit - commission

        if netProfit >= -0.01 && netProfit <= 0.01 { // Within 1 cent of break-even
            fmt.Printf("Break-even price: %.5f (Gross: %.2f, Commission: %.2f, Net: %.2f)\n",
                testPrice, profitData.Profit, commission, netProfit)
            return testPrice, nil
        }
    }

    return 0, fmt.Errorf("break-even price not found in range")
}
```

---

## 🔧 Common Patterns

### Pre-trade profit estimation

```go
func EstimateTradeProfit(account *mt5.MT5Account, symbol string, volume float64, entry, tp, sl float64, isBuy bool) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    orderType := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
    if !isBuy {
        orderType = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
    }

    // Calculate TP profit
    tpProfit, _ := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
        OrderType:  orderType,
        Symbol:     symbol,
        Volume:     volume,
        OpenPrice:  entry,
        ClosePrice: tp,
    })

    // Calculate SL loss
    slProfit, _ := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
        OrderType:  orderType,
        Symbol:     symbol,
        Volume:     volume,
        OpenPrice:  entry,
        ClosePrice: sl,
    })

    fmt.Printf("Trade estimate:\n")
    fmt.Printf("  Max profit (TP): %.2f\n", tpProfit.Profit)
    fmt.Printf("  Max loss (SL): %.2f\n", slProfit.Profit)
    fmt.Printf("  R:R ratio: 1:%.2f\n", tpProfit.Profit/(-slProfit.Profit))
}
```

### Position sizing by risk

```go
func CalculateLotSizeByDollarRisk(account *mt5.MT5Account, symbol string, entry, sl float64, isBuy bool, riskAmount float64) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    orderType := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY
    if !isBuy {
        orderType = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL
    }

    // Calculate loss per 1.0 lot
    testProfit, err := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
        OrderType:  orderType,
        Symbol:     symbol,
        Volume:     1.0,
        OpenPrice:  entry,
        ClosePrice: sl,
    })
    if err != nil {
        return 0, err
    }

    lossPerLot := -testProfit.Profit
    lots := riskAmount / lossPerLot

    fmt.Printf("Loss per 1.0 lot: %.2f\n", lossPerLot)
    fmt.Printf("Desired risk: %.2f\n", riskAmount)
    fmt.Printf("Calculated lot size: %.2f\n", lots)

    return lots, nil
}
```

---

## 📚 See Also

* [OrderCalcMargin](./OrderCalcMargin.md) - Calculate required margin for an order
* [OrderCheck](./OrderCheck.md) - Validate complete order before sending
* [OrderSend](./OrderSend.md) - Place market or pending order
* [AccountSummary](../1.%20Account_information/AccountSummary.md) - Get account balance for risk calculations
