# MT5Account · Trading Operations - Overview

> Order execution, position management, margin calculations, trade validation. Use this page to choose the right API for trading operations.

## 📁 What lives here

### Order Execution & Management

* **[OrderSend](./OrderSend.md)** - place market or pending orders.
* **[OrderModify](./OrderModify.md)** - modify SL/TP or order parameters.
* **[OrderClose](./OrderClose.md)** - close positions (full or partial).

### Pre-Trade Calculations

* **[OrderCalcMargin](./OrderCalcMargin.md)** - calculate margin required for trade.
* **[OrderCalcProfit](./OrderCalcProfit.md)** - calculate potential profit for price movement.
* **[OrderCheck](./OrderCheck.md)** - validate trade request before execution.

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[OrderSend - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderSend_HOW.md)**
* **[OrderModify - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderModify_HOW.md)**
* **[OrderClose - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderClose_HOW.md)**
* **[OrderCalcMargin - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderCalcMargin_HOW.md)**
* **[OrderCalcProfit - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderCalcProfit_HOW.md)**
* **[OrderCheck - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderCheck_HOW.md)**

---

## 🧭 Plain English

* **OrderSend** → **place orders** (market BUY/SELL, limit, stop orders).
* **OrderModify** → **change SL/TP** or modify pending order price.
* **OrderClose** → **close positions** completely or partially.
* **OrderCalcMargin** → **calculate margin** needed before trading.
* **OrderCalcProfit** → **estimate profit** for potential trade.
* **OrderCheck** → **validate order** before sending (dry-run).

> Rule of thumb: **always use** `OrderCheck` before `OrderSend`; use `OrderCalcMargin` to verify sufficient margin; use `OrderModify` for SL/TP changes.

---

## Quick choose

| If you need…                                     | Use                      | Returns                    | Key inputs                          |
| ------------------------------------------------ | ------------------------ | -------------------------- | ----------------------------------- |
| Place market/pending order                       | `OrderSend`              | OrderSendResult            | Symbol, type, volume, price, SL, TP |
| Modify SL/TP or order price                      | `OrderModify`            | OrderModifyResult          | Ticket, new SL, new TP, new price   |
| Close position (full/partial)                    | `OrderClose`             | OrderCloseResult           | Ticket, volume                      |
| Calculate required margin                        | `OrderCalcMargin`        | float64                    | Symbol, type, volume, price         |
| Calculate potential profit                       | `OrderCalcProfit`        | float64                    | Symbol, type, volume, open, close   |
| Validate order before execution                  | `OrderCheck`             | OrderCheckResult           | Symbol, type, volume, price, SL, TP |

---

## ❌ Cross‑refs & gotchas

* **Return codes** - Check `ReturnedCode` field: 10009 = success, others = error.
* **OrderCheck first** - ALWAYS call OrderCheck before OrderSend to avoid rejections.
* **Margin calculation** - Use OrderCalcMargin to verify sufficient margin before trading.
* **Order types** - BUY, SELL (market), BUY_LIMIT, SELL_LIMIT, BUY_STOP, SELL_STOP (pending).
* **Volume** - Must be multiple of VOLUME_STEP and between VOLUME_MIN and VOLUME_MAX.
* **SL/TP validation** - Must respect STOP_LEVEL (minimum distance from current price).
* **Partial close** - OrderClose can close partial volume (e.g., close 0.5 of 1.0 lot position).
* **Comment** field in result - Contains detailed error message if order fails.

---

## 🟢 Minimal snippets

```go
// Place BUY market order
result, err := account.OrderSend(ctx, &pb.OrderSendRequest{
    Symbol:   "EURUSD",
    Type:     pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:   0.1,
    StopLoss: 1.08000,
    TakeProfit: 1.09000,
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

if result.ReturnedCode == 10009 {
    fmt.Printf("✅ Order placed: #%d\n", result.Order)
} else {
    fmt.Printf("❌ Order failed: %s (code %d)\n", result.Comment, result.ReturnedCode)
}
```

```go
// Validate order before sending
checkResult, _ := account.OrderCheck(ctx, &pb.OrderCheckRequest{
    Symbol:   "EURUSD",
    Type:     pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:   0.1,
    StopLoss: 1.08000,
    TakeProfit: 1.09000,
})

if checkResult.ReturnedCode == 10009 {
    fmt.Println("✅ Order validation passed")
    // Now safe to send
    result, _ := account.OrderSend(ctx, &pb.OrderSendRequest{...})
} else {
    fmt.Printf("❌ Validation failed: %s\n", checkResult.Comment)
}
```

```go
// Calculate margin before trading
margin, _ := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
    Symbol:    "EURUSD",
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:    1.0,
    OpenPrice: 1.08500,
})

fmt.Printf("Required margin: $%.2f\n", margin)

// Check if we have enough
accountSummary, _ := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
freeMargin := accountSummary.AccountFreeMargin

if margin <= freeMargin {
    fmt.Println("✅ Sufficient margin")
} else {
    fmt.Printf("❌ Insufficient margin (need $%.2f, have $%.2f)\n", margin, freeMargin)
}
```

```go
// Modify SL/TP of existing position
result, _ := account.OrderModify(ctx, &pb.OrderModifyRequest{
    Ticket:     123456,
    StopLoss:   1.08200, // New SL
    TakeProfit: 1.09500, // New TP
})

if result.ReturnedCode == 10009 {
    fmt.Println("✅ SL/TP modified")
} else {
    fmt.Printf("❌ Modification failed: %s\n", result.Comment)
}
```

```go
// Close position
result, _ := account.OrderClose(ctx, &pb.OrderCloseRequest{
    Ticket: 123456,
    Volume: 0.0, // 0 = close full position
})

if result.ReturnedCode == 10009 {
    fmt.Printf("✅ Position closed at %.5f\n", result.ClosePrice)
    fmt.Printf("   Profit: $%.2f\n", result.Profit)
} else {
    fmt.Printf("❌ Close failed: %s\n", result.Comment)
}
```

```go
// Calculate potential profit
profit, _ := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
    Symbol:     "EURUSD",
    OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:     0.1,
    OpenPrice:  1.08500,
    ClosePrice: 1.09000, // +50 pips
})

fmt.Printf("Potential profit for 50 pip move: $%.2f\n", profit)
```

```go
// Complete professional trading workflow
symbol := "EURUSD"
volume := 0.1
stopLoss := 1.08000
takeProfit := 1.09000

// Step 1: Calculate margin
margin, _ := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
    Symbol:    symbol,
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:    volume,
})

// Step 2: Check free margin
summary, _ := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
if margin > summary.AccountFreeMargin {
    log.Fatal("Insufficient margin")
}

// Step 3: Validate order
checkResult, _ := account.OrderCheck(ctx, &pb.OrderCheckRequest{
    Symbol:     symbol,
    Type:       pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:     volume,
    StopLoss:   stopLoss,
    TakeProfit: takeProfit,
})

if checkResult.ReturnedCode != 10009 {
    log.Fatalf("Validation failed: %s", checkResult.Comment)
}

// Step 4: Send order
result, _ := account.OrderSend(ctx, &pb.OrderSendRequest{
    Symbol:     symbol,
    Type:       pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:     volume,
    StopLoss:   stopLoss,
    TakeProfit: takeProfit,
})

if result.ReturnedCode == 10009 {
    fmt.Printf("✅ Order placed: #%d\n", result.Order)
} else {
    log.Fatalf("Order failed: %s", result.Comment)
}
```

---

## See also

* **Account info:** [AccountSummary](../1.%20Account_information/AccountSummary.md) - check margin before trading
* **Position info:** [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - get current positions
* **Symbol info:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get current prices
* **Real-time updates:** [OnTrade](../7.%20Streaming_Methods/OnTrade.md) - monitor trade events
