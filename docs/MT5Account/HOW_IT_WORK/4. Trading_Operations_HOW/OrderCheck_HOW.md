### Example from file: `examples/demos/lowlevel/02_trading_operations.go`

> The **`OrderCheck()`** method performs a safe preliminary validation of an order before sending it.
> It allows you to verify that the order parameters are correct and that the account has sufficient funds to open it.


---

## 🧩 Code example

```go
fmt.Println("\n3.1. OrderCheck() - Validate order parameters")

orderCheckReq := &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:      pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
        Symbol:      cfg.TestSymbol,
        Volume:      cfg.TestVolume,
        OrderType:   pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
        Price:       tickData.Ask,
        StopLoss:    0.0,
        TakeProfit:  0.0,
        Deviation:   10,
        TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_IOC, // Immediate-Or-Cancel
        TypeTime:    pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC,
        Comment:     "OrderCheck validation",
    },
}

checkData, err := account.OrderCheck(ctx, orderCheckReq)
if err != nil {
    fmt.Println("  ❌ OrderCheck FAILED")
    fmt.Printf("     Error: %v\n", err)
} else if checkData != nil && checkData.MqlTradeCheckResult != nil {
    result := checkData.MqlTradeCheckResult
    fmt.Println("  ✅ OrderCheck SUCCESS!")
    fmt.Printf("     Return Code:        %d\n", result.ReturnedCode)
    fmt.Printf("     Comment:            %s\n", result.Comment)
    fmt.Printf("     Required Margin:    %.2f\n", result.Margin)
    fmt.Printf("     Balance After:      %.2f\n", result.BalanceAfterDeal)
    fmt.Printf("     Equity After:       %.2f\n", result.EquityAfterDeal)
    fmt.Printf("     Free Margin After:  %.2f\n", result.FreeMargin)
    fmt.Printf("     Margin Level:       %.2f%%\n\n", result.MarginLevel)

    // OrderCheck returns 0 on successful validation (unlike OrderSend which returns 10009)
    if result.ReturnedCode == 0 {
        fmt.Println("  ✓ Order validation PASSED – safe to proceed with OrderSend()")
    } else {
        // Use helper from errors.go to get human-readable error description
        fmt.Printf("  ⚠️  Validation failed: %s (code: %d)\n",
            mt5.GetRetCodeMessage(uint32(result.ReturnedCode)),
            result.ReturnedCode)
    }
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Form the Request Structure

```go
orderCheckReq := &pb.OrderCheckRequest{ MqlTradeRequest: &pb.MrpcMqlTradeRequest{ ... } }
```

A request object is created, containing the **`MqlTradeRequest`** structure — a complete description of the order, as if it were actually being sent to MetaTrader.

Key fields:

* `Action` – trade action type (`TRADE_ACTION_DEAL` — execute immediately);
* `Symbol` – instrument (`cfg.TestSymbol` from `config.json`);
* `Volume` – order volume;
* `OrderType` – direction (`BUY` or `SELL`);
* `Price` – opening price;
* `Deviation` – allowable price deviation in points;
* `TypeFilling` – execution mode (`IOC` – Immediate-Or-Cancel);
* `TypeTime` – order lifetime (`GTC` – Good Till Cancelled).

---

### 2️. Send Validation Request

```go
checkData, err := account.OrderCheck(ctx, orderCheckReq)
```

The gateway passes the request to MetaTrader, which validates the order for correctness and returns the result — **without actually opening a position**.

---

### 3️. Process the Result

If there is an error (`err != nil`) — this means the broker or server does not support `OrderCheck` (common on demo accounts).
If a result is received, the validation fields are displayed:

```go
result := checkData.MqlTradeCheckResult
fmt.Printf("     Required Margin:    %.2f\n", result.Margin)
fmt.Printf("     Free Margin After:  %.2f\n", result.FreeMargin)
fmt.Printf("     Margin Level:       %.2f%%\n", result.MarginLevel)
```

* **Margin** – how much margin will be required for the trade;
* **FreeMargin** – how much free margin will remain after opening;
* **MarginLevel** – margin level as a percentage.

---

### 4️. Interpret the Result

```go
// OrderCheck returns 0 on successful validation (unlike OrderSend which returns 10009)
if result.ReturnedCode == 0 {
    fmt.Println("  ✓ Order validation PASSED – safe to proceed with OrderSend()")
} else {
    // Use helper from errors.go to get human-readable error description
    fmt.Printf("  ⚠️  Validation failed: %s (code: %d)\n",
        mt5.GetRetCodeMessage(uint32(result.ReturnedCode)),
        result.ReturnedCode)
}
```

**Important: OrderCheck uses different success code than trading operations**

1. **OrderCheck success**: `ReturnedCode == 0` means validation passed
2. **Trading operations success**: `ReturnedCode == 10009` (checked with `mt5.IsRetCodeSuccess()`)

**Why use `mt5.GetRetCodeMessage()` for errors?**

1. **Automatic Error Descriptions**: Instead of just showing a number, get human-readable messages for common validation errors:
   - `TRADE_RETCODE_INVALID_VOLUME` - Volume doesn't meet broker requirements
   - `TRADE_RETCODE_NO_MONEY` - Insufficient funds for the trade
   - `TRADE_RETCODE_INVALID_STOPS` - SL/TP levels violate broker rules
   - `TRADE_RETCODE_INVALID_PRICE` - Price out of valid range
2. **Consistency**: Same error handling pattern used across all trading operations
3. **Centralized**: All 40+ error codes defined in `package/Helpers/errors.go`

If `ReturnedCode == 0` → the order passed validation and it's safe to proceed with `OrderSend()`.
Any other value indicates a validation error that should be fixed before attempting to send the order.

---

## 📦 What the Server Returns

```protobuf
message OrderCheckData {
  MrpcMqlTradeCheckResult MqlTradeCheckResult = 1;
}

message MrpcMqlTradeCheckResult {
  int32 ReturnedCode = 1;
  string Comment = 2;
  double Margin = 3;
  double BalanceAfterDeal = 4;
  double EquityAfterDeal = 5;
  double FreeMargin = 6;
  double MarginLevel = 7;
}
```

---

## 💡 Example Output

```
✅ OrderCheck SUCCESS!
     Return Code:        0
     Comment:            Request verified successfully
     Required Margin:    120.50
     Balance After:      9879.50
     Equity After:       9879.50
     Free Margin After:  9759.00
     Margin Level:       1120.75%

✓ Order validation PASSED – safe to proceed with OrderSend()
```

---

## 🧠 What It's Used For

The `OrderCheck()` method is used for:

* validating order correctness **before sending**;
* assessing future account state after the trade;
* protection from errors and broker rejections;
* safe testing of strategies and robots.

---

## 💬 In Simple Terms

> `OrderCheck()` is a "rehearsal" before a trade.
> It checks whether an order with the specified parameters can be opened, and how balance, margin, and equity will change after the trade.
> If everything is correct (code `0`) — you can safely proceed with `OrderSend()`.
