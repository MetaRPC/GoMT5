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
| Place market/pending order                       | `OrderSend`              | OrderSendData              | Symbol, Operation, Volume, Price, StopLoss, TakeProfit |
| Modify SL/TP or order price                      | `OrderModify`            | OrderModifyData            | Ticket, new SL, new TP, new price   |
| Close position (full/partial)                    | `OrderClose`             | OrderCloseData             | Ticket, Volume, Slippage            |
| Calculate required margin                        | `OrderCalcMargin`        | OrderCalcMarginData (`Margin` float64) | Symbol, OrderType, Volume, OpenPrice |
| Calculate potential profit                       | `OrderCalcProfit`        | OrderCalcProfitData (`Profit` float64) | Symbol, OrderType, Volume, OpenPrice, ClosePrice |
| Validate order before execution                  | `OrderCheck`             | OrderCheckData (`MqlTradeCheckResult`) | MqlTradeRequest (Action, Symbol, OrderType, Volume, Price, SL, TP) |

---

## ❌ Cross‑refs & gotchas

* **Return codes** - Check `ReturnedCode` field: 10009 = success for OrderSend/OrderModify/OrderClose; OrderCheck reports 0 when the check passed.
* **OrderCheck first** - ALWAYS call OrderCheck before OrderSend to avoid rejections.
* **Margin calculation** - Use OrderCalcMargin to verify sufficient margin before trading.
* **Order types** - BUY, SELL (market), BUY_LIMIT, SELL_LIMIT, BUY_STOP, SELL_STOP (pending).
* **Volume** - Must be multiple of VOLUME_STEP and between VOLUME_MIN and VOLUME_MAX.
* **SL/TP validation** - Must respect SYMBOL_TRADE_STOPS_LEVEL (minimum distance from current price).
* **Partial close** - OrderClose can close partial volume (e.g., close 0.5 of 1.0 lot position).
* **Result description** - `ReturnedCodeDescription` (OrderSend/OrderModify/OrderClose) and `MqlTradeCheckResult.Comment` (OrderCheck) describe the result code.

---

## 🟢 Minimal snippets

```go
// Place BUY market order (optional fields are pointers)
sl, tp := 1.08000, 1.09000
result, err := account.OrderSend(ctx, &pb.OrderSendRequest{
    Symbol:     "EURUSD",
    Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
    Volume:     0.1,
    StopLoss:   &sl,
    TakeProfit: &tp,
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

if result.ReturnedCode == 10009 {
    fmt.Printf("✅ Order placed: #%d\n", result.Order)
} else {
    fmt.Printf("❌ Order failed: %s (code %d)\n", result.ReturnedCodeDescription, result.ReturnedCode)
}
```

```go
// Validate order before sending (OrderCheck takes a nested MrpcMqlTradeRequest)
tick, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: "EURUSD"})
checkData, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
        Symbol:     "EURUSD",
        OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:     0.1,
        Price:      tick.Ask,
        StopLoss:   1.08000,
        TakeProfit: 1.09000,
    },
})

if err == nil && checkData.MqlTradeCheckResult.ReturnedCode == 0 {
    fmt.Println("✅ Order validation passed")
    // Now safe to send the order with OrderSend
} else if err != nil {
    fmt.Printf("❌ Validation failed: %v\n", err)
} else {
    fmt.Printf("❌ Validation failed: %s\n", checkData.MqlTradeCheckResult.Comment)
}
```

```go
// Calculate margin before trading
marginData, _ := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
    Symbol:    "EURUSD",
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:    1.0,
    OpenPrice: 1.08500,
})
margin := marginData.Margin

fmt.Printf("Required margin: $%.2f\n", margin)

// Check if we have enough (AccountSummary has no free margin: use AccountInfoDouble)
freeData, _ := account.AccountInfoDouble(ctx, &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE,
})
freeMargin := freeData.RequestedValue

if margin <= freeMargin {
    fmt.Println("✅ Sufficient margin")
} else {
    fmt.Printf("❌ Insufficient margin (need $%.2f, have $%.2f)\n", margin, freeMargin)
}
```

```go
// Modify SL/TP of existing position (optional fields are pointers)
newSL, newTP := 1.08200, 1.09500
result, _ := account.OrderModify(ctx, &pb.OrderModifyRequest{
    Ticket:     123456,
    StopLoss:   &newSL, // New SL
    TakeProfit: &newTP, // New TP
})

if result.ReturnedCode == 10009 {
    fmt.Println("✅ SL/TP modified")
} else {
    fmt.Printf("❌ Modification failed: %s\n", result.ReturnedCodeDescription)
}
```

```go
// Close position
result, _ := account.OrderClose(ctx, &pb.OrderCloseRequest{
    Ticket:   123456,
    Volume:   0.0, // 0 = close full position
    Slippage: 20,
})

// OrderCloseData has no close price or profit: only the result code and the close mode
if result.ReturnedCode == 10009 {
    fmt.Printf("✅ Position closed (mode: %v)\n", result.CloseMode)
} else {
    fmt.Printf("❌ Close failed: %s\n", result.ReturnedCodeDescription)
}
```

```go
// Calculate potential profit
profitData, _ := account.OrderCalcProfit(ctx, &pb.OrderCalcProfitRequest{
    Symbol:     "EURUSD",
    OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:     0.1,
    OpenPrice:  1.08500,
    ClosePrice: 1.09000, // +50 pips
})

fmt.Printf("Potential profit for 50 pip move: $%.2f\n", profitData.Profit)
```

```go
// Complete professional trading workflow
symbol := "EURUSD"
volume := 0.1
stopLoss := 1.08000
takeProfit := 1.09000

tick, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: symbol})

// Step 1: Calculate margin
marginData, _ := account.OrderCalcMargin(ctx, &pb.OrderCalcMarginRequest{
    Symbol:    symbol,
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:    volume,
    OpenPrice: tick.Ask,
})

// Step 2: Check free margin
freeData, _ := account.AccountInfoDouble(ctx, &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE,
})
if marginData.Margin > freeData.RequestedValue {
    log.Fatal("Insufficient margin")
}

// Step 3: Validate order
checkData, err := account.OrderCheck(ctx, &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
        Symbol:     symbol,
        OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Volume:     volume,
        Price:      tick.Ask,
        StopLoss:   stopLoss,
        TakeProfit: takeProfit,
    },
})

if err != nil {
    log.Fatalf("Validation failed: %v", err)
}
if checkData.MqlTradeCheckResult.ReturnedCode != 0 {
    log.Fatalf("Validation failed: %s", checkData.MqlTradeCheckResult.Comment)
}

// Step 4: Send order
result, _ := account.OrderSend(ctx, &pb.OrderSendRequest{
    Symbol:     symbol,
    Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
    Volume:     volume,
    StopLoss:   &stopLoss,
    TakeProfit: &takeProfit,
})

if result.ReturnedCode == 10009 {
    fmt.Printf("✅ Order placed: #%d\n", result.Order)
} else {
    log.Fatalf("Order failed: %s", result.ReturnedCodeDescription)
}
```

---

## See also

* **Account info:** [AccountSummary](../1.%20Account_information/AccountSummary.md) - check margin before trading
* **Position info:** [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - get current positions
* **Symbol info:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get current prices
* **Real-time updates:** [OnTrade](../7.%20Streaming_Methods/OnTrade.md) - monitor trade events
