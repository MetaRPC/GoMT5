### Example from file: `examples\demos\lowlevel\02_trading_operations.go`

> The **`OrderModify()`** method is used to **modify parameters of an already opened position or order**, such as adding Stop Loss (SL) and Take Profit (TP) levels.
> This call does not open a new trade — it updates an existing one by its ticket.


---

## 🧩 Code example

```go
fmt.Println("\n4.2. OrderModify() - Add Stop Loss and Take Profit")

// Calculate SL/TP levels (10 pips SL, 20 pips TP)
stopLoss := sendData.Price - (10 * pipSize)
takeProfit := sendData.Price + (20 * pipSize)

modifyReq := &pb.OrderModifyRequest{
    Ticket:     orderTicket,
    StopLoss:   &stopLoss,
    TakeProfit: &takeProfit,
}

modifyData, err := account.OrderModify(ctx, modifyReq)
if !helpers.PrintShortError(err, "OrderModify failed") {
    fmt.Printf("  Order modify result:\n")
    fmt.Printf("    Return Code:                 %d\n", modifyData.ReturnedCode)
    fmt.Printf("    Order Ticket:                %d\n", modifyData.Order)
    fmt.Printf("    Stop Loss:                   %.5f\n", stopLoss)
    fmt.Printf("    Take Profit:                 %.5f\n", takeProfit)
    fmt.Printf("    Comment:                     %s\n", modifyData.Comment)

    if modifyData.ReturnedCode == 10009 {
        fmt.Printf("    ✓ Position MODIFIED successfully!\n")
    }
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Calculate SL and TP Levels

```go
stopLoss := sendData.Price - (10 * pipSize)
takeProfit := sendData.Price + (20 * pipSize)
```

After the order is opened (`OrderSend()` returned `sendData.Price`), stop levels are calculated:

* **StopLoss** is set 10 pips below the opening price;
* **TakeProfit** — 20 pips above.

`pipSize` — size of one pip (for EURUSD = 0.0001).

> 💡 This way we programmatically set protection levels and profit targets.

---

### 2️. Form the Modification Request

```go
modifyReq := &pb.OrderModifyRequest{
    Ticket:     orderTicket,
    StopLoss:   &stopLoss,
    TakeProfit: &takeProfit,
}
```

Create the `OrderModifyRequest` structure.

| Field        | Description                                                                    |
| ------------ | ------------------------------------------------------------------------------ |
| `Ticket`     | Order number to modify (obtained from `OrderSend()`).                          |
| `StopLoss`   | Stop loss price. Uses pointer (`&`), as the field is optional.                 |
| `TakeProfit` | Take profit price. Also a pointer.                                             |

> You can pass only one of these fields if you need to change only SL or TP.

---

### 3️. Send Request to Broker

```go
modifyData, err := account.OrderModify(ctx, modifyReq)
```

The gateway method calls `OrderModify()` on the MetaTrader side.
The server validates level correctness (e.g., SL/TP must not be closer than the minimum `StopsLevel` distance).

---

### 4️. Parse the Response

```go
fmt.Printf("    Return Code: %d\n", modifyData.ReturnedCode)
fmt.Printf("    Order Ticket: %d\n", modifyData.Order)
fmt.Printf("    Stop Loss: %.5f\n", stopLoss)
fmt.Printf("    Take Profit: %.5f\n", takeProfit)
fmt.Printf("    Comment: %s\n", modifyData.Comment)
```

Key fields:

* **ReturnedCode** — operation status (`10009` = successfully executed);
* **Order** — ticket of the modified order;
* **Comment** — server response (e.g., "Request executed successfully").

---

### 5️. Verify Operation Success

```go
if modifyData.ReturnedCode == 10009 {
    fmt.Printf("    ✓ Position MODIFIED successfully!\n")
}
```

Code `10009` (`TRADE_RETCODE_DONE`) means the changes were applied.

If the code is different — the broker may have rejected the request (e.g., SL/TP levels too close or order already closed).

---

## 📦 What the Server Returns

```protobuf
message OrderModifyData {
  uint32 ReturnedCode = 1;  // Result code (TRADE_RETCODE_*)
  uint64 Order = 2;         // Order number
  string Comment = 3;       // Broker comment
}
```

---

## 💡 Example Output

```
Order modify result:
    Return Code:                 10009
    Order Ticket:                128013245
    Stop Loss:                   1.08442
    Take Profit:                 1.08742
    Comment:                     Request executed successfully
    ✓ Position MODIFIED successfully!
```

---

## 🧠 What It's Used For

The `OrderModify()` method is used for:

* setting Stop Loss and Take Profit after opening a trade;
* changing existing SL/TP levels;
* managing active positions from code;
* adapting orders to market conditions without closing the position.

---

## 💬 In Simple Terms

> `OrderModify()` is **editing an open position**.
> It allows you to add or change stop loss and take profit for an existing order.
> If everything is correct (code `10009`), the broker applies the changes, and the trade is updated on the server.
