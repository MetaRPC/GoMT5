### Example from file: `examples/demos/lowlevel/02_trading_operations.go`

> The **`OrderClose()`** method is used to **close a previously opened position** by order ticket.
> This is the final step of the trading cycle: after `OrderSend()` (opening) and `OrderModify()` (modifying parameters), `OrderClose()` completes the trade and locks in the result — profit or loss.


---

## 🧩 Code example

```go
fmt.Println("\n4.3. OrderClose() - Close the position")

closeReq := &pb.OrderCloseRequest{
    Ticket:   orderTicket,
    Volume:   cfg.TestVolume,
    Slippage: 10,
}

closeData, err := account.OrderClose(ctx, closeReq)
if !helpers.PrintShortError(err, "OrderClose failed") {
    fmt.Printf("  Order close result:\n")
    fmt.Printf("    Return Code:                 %d (%s)\n", closeData.ReturnedCode, closeData.ReturnedStringCode)
    fmt.Printf("    Description:                 %s\n", closeData.ReturnedCodeDescription)
    fmt.Printf("    Close Mode:                  %v\n", closeData.CloseMode)

    // Check if position closed successfully using helper from errors.go
    if mt5.IsRetCodeSuccess(closeData.ReturnedCode) {
        fmt.Printf("    ✓ Position CLOSED successfully!\n")
    } else {
        fmt.Printf("    ❌ Close failed: %s (code: %d)\n",
            mt5.GetRetCodeMessage(closeData.ReturnedCode),
            closeData.ReturnedCode)
    }
} else {
    fmt.Printf("    ✗ Order execution FAILED - check return code and comment\n")
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Form the Close Request

```go
closeReq := &pb.OrderCloseRequest{
    Ticket:   orderTicket,
    Volume:   cfg.TestVolume,
    Slippage: 10,
}
```

| Field        | Description                                                    |
| ------------ | -------------------------------------------------------------- |
| **Ticket**   | Order number to close. Obtained from `OrderSend()`.            |
| **Volume**   | Position volume to close. Can be partial.                      |
| **Slippage** | Maximum allowable price deviation (in points).                 |

> 💡 Slippage — protection from slippage. If the price moved too far, the broker may reject the close.

---

### 2️. Send Request to Broker

```go
closeData, err := account.OrderClose(ctx, closeReq)
```

The gateway method makes a gRPC call to the MetaTrader server.
If the broker accepts the request — it closes the position at market price or considering the specified `slippage`.

---

### 3️. Process the Response

```go
fmt.Printf("    Return Code: %d (%s)\n", closeData.ReturnedCode, closeData.ReturnedStringCode)
fmt.Printf("    Description: %s\n", closeData.ReturnedCodeDescription)
fmt.Printf("    Close Mode:  %v\n", closeData.CloseMode)
```

**Main response fields:**

| Field                       | Value                                                    |
| --------------------------- | -------------------------------------------------------- |
| **ReturnedCode**            | Result code (`10009` = successful).                      |
| **ReturnedStringCode**      | Text designation (e.g., `TRADE_RETCODE_DONE`).           |
| **ReturnedCodeDescription** | Detailed result description (human-readable).            |
| **CloseMode**               | Close mode (usually `MARKET`).                           |

---

### 4️. Verify Operation Success

```go
// Check if position closed successfully using helper from errors.go
if mt5.IsRetCodeSuccess(closeData.ReturnedCode) {
    fmt.Printf("    ✓ Position CLOSED successfully!\n")
} else {
    fmt.Printf("    ❌ Close failed: %s (code: %d)\n",
        mt5.GetRetCodeMessage(closeData.ReturnedCode),
        closeData.ReturnedCode)
}
```

**Why use `mt5.IsRetCodeSuccess()` instead of checking `== 10009` manually?**

1. **Code Clarity**: Function name `IsRetCodeSuccess()` immediately tells you what is being validated, unlike magic number `10009`
2. **Automatic Error Descriptions**: `mt5.GetRetCodeMessage()` provides ready-to-use error messages for all 40+ return codes
3. **Common Errors Covered**: The helper automatically handles typical close operation errors:
   - `TRADE_RETCODE_MARKET_CLOSED` - Market closed, cannot close position
   - `TRADE_RETCODE_INVALID_VOLUME` - Invalid volume specified
   - `TRADE_RETCODE_REQUOTE` - Price changed, retry needed
   - `TRADE_RETCODE_POSITION_CLOSED` - Position already closed
4. **Maintainability**: Centralized logic in `package/Helpers/errors.go` - one place to update for all projects

Code **`10009` (`TRADE_RETCODE_DONE`)** — standard confirmation of successful trading request execution.

If the code is different, the helper functions automatically provide the rejection reason without manually checking `ReturnedCodeDescription`. This makes error handling consistent across all trading operations.

---

## 📦 What the Server Returns

```protobuf
message OrderCloseData {
  uint32 ReturnedCode = 1;           // Result code (TRADE_RETCODE_*)
  string ReturnedStringCode = 2;     // Text designation
  string ReturnedCodeDescription = 3;// Result description
  ENUM_CLOSE_MODE CloseMode = 4;     // Position close mode
}
```

---

## 💡 Example Output

```
Order close result:
    Return Code:                 10009 (TRADE_RETCODE_DONE)
    Description:                 Request executed successfully
    Close Mode:                  MARKET
    ✓ Position CLOSED successfully!
```

---

## 🧠 What It's Used For

The `OrderClose()` method is used for:

* closing an active position fully or partially;
* completing the trading cycle (open → modify → close);
* testing trading strategies with full order lifecycle;
* locking in profit or loss on a trade.

---

## 💬 In Simple Terms

> `OrderClose()` is the **command to the broker to close a trade**.
> It says: "Close this position at market, if the price hasn't moved more than 10 pips (slippage)".
> If `10009` is returned — everything succeeded: the position is closed, the result is locked in, the trading cycle is complete.
