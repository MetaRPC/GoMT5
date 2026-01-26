### Example from file: `examples/demos/lowlevel/02_trading_operations.go`


> The **`OrderSend()`** method executes a **real trading operation** — sends an order to the broker through the gateway and returns the processing result.
> While `OrderCheck()` only validates parameter correctness, `OrderSend()` is the actual **order execution**.


---

## 🧱 Primary Purpose

`OrderSend()` is used for:

* opening market trades (`BUY`, `SELL`);
* placing pending orders (`BUY_LIMIT`, `SELL_STOP`, etc.);
* executing trading strategies through the gateway in real time.

After execution, the broker returns a result — the order was **executed**, **rejected**, or **queued**.

---

## 🧩 Code example

```go
fmt.Println("\n4.1. OrderSend() - Place market BUY order")

slippage := uint64(10)
orderSendReq := &pb.OrderSendRequest{
    Symbol:    cfg.TestSymbol,
    Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
    Volume:    cfg.TestVolume,
    Price:     &tickData.Ask,
    Slippage:  &slippage,
}

sendData, err := account.OrderSend(ctx, orderSendReq)
if !helpers.PrintShortError(err, "OrderSend failed") {
    fmt.Printf("  Order sent result:\n")
    fmt.Printf("    Return Code:                 %d\n", sendData.ReturnedCode)
    fmt.Printf("    Deal Ticket:                 %d\n", sendData.Deal)
    fmt.Printf("    Order Ticket:                %d\n", sendData.Order)
    fmt.Printf("    Volume:                      %.2f\n", sendData.Volume)
    fmt.Printf("    Execution Price:             %.5f\n", sendData.Price)

    if sendData.Bid > 0 && sendData.Ask > 0 {
        fmt.Printf("    Market Bid:                  %.5f\n", sendData.Bid)
        fmt.Printf("    Market Ask:                  %.5f\n", sendData.Ask)
    }

    fmt.Printf("    Comment:                     %s\n", sendData.Comment)

    // Check if order executed successfully using helper from errors.go
    if mt5.IsRetCodeSuccess(sendData.ReturnedCode) {
        fmt.Printf("    ✓ Order EXECUTED successfully!\n")
        orderTicket := sendData.Order
        _ = orderTicket // example placeholder
    } else {
        fmt.Printf("    ❌ Order failed: %s (code: %d)\n",
            mt5.GetRetCodeMessage(sendData.ReturnedCode),
            sendData.ReturnedCode)
    }
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Define Allowable Price Deviation

```go
slippage := uint64(10)
```

`slippage` — maximum allowable price deviation (in points) between the moment of sending and trade execution.

> For example: if you want to buy at 1.0850, but the price moved to 1.0851 — with `slippage = 10` the order will still be accepted.

---

### 2️. Form the Request

```go
orderSendReq := &pb.OrderSendRequest{ ... }
```

Set order parameters:

* **Symbol** — instrument (`cfg.TestSymbol` from configuration);
* **Operation** — order type (`BUY` or `SELL`);
* **Volume** — trade volume in lots (`cfg.TestVolume`);
* **Price** — current Ask price for buying;
* **Slippage** — maximum allowable price deviation.

> You can also specify StopLoss, TakeProfit or Comment — either here or later through `OrderModify()`.

---

### 3️. Send Order to Server

```go
sendData, err := account.OrderSend(ctx, orderSendReq)
```

A gRPC call is made. The MetaTrader server validates parameters and executes the order (or rejects it).

---

### 4️. Parse the Response

```go
fmt.Printf("    Return Code: %d\n", sendData.ReturnedCode)
fmt.Printf("    Deal Ticket: %d\n", sendData.Deal)
fmt.Printf("    Order Ticket: %d\n", sendData.Order)
fmt.Printf("    Execution Price: %.5f\n", sendData.Price)
```

Key response fields:

* **ReturnedCode** — operation result (`10009` = successfully executed);
* **Deal** — deal number, if the order was already executed;
* **Order** — order number, if the order was created;
* **Volume**, **Price** — execution volume and price.

---

### 5️. Display Additional Data (Bid/Ask)

```go
if sendData.Bid > 0 && sendData.Ask > 0 {
    fmt.Printf("    Market Bid: %.5f\n", sendData.Bid)
    fmt.Printf("    Market Ask: %.5f\n", sendData.Ask)
}
```

MetaTrader may return current market prices at the moment of execution. If the broker does not provide this data, they are zero.

---

### 6️. Verify Execution Success

```go
// Check if order executed successfully using helper from errors.go
if mt5.IsRetCodeSuccess(sendData.ReturnedCode) {
    fmt.Printf("    ✓ Order EXECUTED successfully!\n")
    orderTicket := sendData.Order
    _ = orderTicket // example placeholder
} else {
    fmt.Printf("    ❌ Order failed: %s (code: %d)\n",
        mt5.GetRetCodeMessage(sendData.ReturnedCode),
        sendData.ReturnedCode)
}
```

**Why use `mt5.IsRetCodeSuccess()` instead of checking `== 10009` manually?**

1. **Readability**: `IsRetCodeSuccess()` is self-documenting - you immediately understand what's being checked
2. **Maintainability**: If MT5 adds more success codes in the future, we update one place instead of all manual checks
3. **Error Messages**: `mt5.GetRetCodeMessage()` provides human-readable descriptions for all 40+ return codes
4. **Type Safety**: The helper function is defined in `package/Helpers/errors.go` with proper type checking

Code `10009` (`TRADE_RETCODE_DONE`) means the order was executed successfully and the position is open. The helper functions are defined in [errors.go](../../../package/Helpers/errors.go) and provide a cleaner, more maintainable way to validate trade operations.

---

## 📦 What the Server Returns

```protobuf
message OrderSendData {
  uint32 ReturnedCode = 1;
  uint64 Deal = 2;
  uint64 Order = 3;
  double Volume = 4;
  double Price = 5;
  double Bid = 6;
  double Ask = 7;
  string Comment = 8;
}
```

---

## 💡 Example Terminal Output

```
Order sent result:
    Return Code:                 10009
    Deal Ticket:                 128013245
    Order Ticket:                128013245
    Volume:                      0.10
    Execution Price:             1.08542
    Market Bid:                  1.08540
    Market Ask:                  1.08542
    Comment:                     Request executed successfully
    ✓ Order EXECUTED successfully!
```

---

## 🧠 What It's Used For

`OrderSend()` is used for:

* opening market and pending orders;
* executing strategies through the gateway;
* testing trading algorithms;
* verifying correct order execution.

---

## 💬 In Simple Terms

> `OrderSend()` is the "real Make Trade button".
> It sends the order to the broker, receives a response, and reports
> whether the order was **executed**, **rejected**, or **awaiting execution**.
