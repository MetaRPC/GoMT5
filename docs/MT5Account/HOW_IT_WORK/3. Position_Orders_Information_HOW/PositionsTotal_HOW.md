### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`PositionsTotal()`** method returns the **total number of open positions** on the trading account.
> This is one of the simplest and most frequently used queries when working with trading strategies.
>
> It does not return details about positions — only their count.
> If you need to get a list with detailed information (symbol, ticket, volume, etc.), use the `OpenedOrders()` or `OpenedOrdersTickets()` methods.


---

## 🧩 Code example

```go
positionsTotalData, err := account.PositionsTotal(ctx)
if err != nil {
    helpers.PrintShortError(err, "PositionsTotal failed")
} else {
    // Direct field access: TotalPositions
    fmt.Printf("  Total open positions:          %d\n", positionsTotalData.TotalPositions)
}
```

---

### 🟢 Detailed Code Explanation

```go
positionsTotalData, err := account.PositionsTotal(ctx)
```

A gRPC request is sent to the MetaTrader server.
In response, the gateway returns a `PositionsTotalResponse` object containing the `TotalPositions` field — the number of open positions on the account.

**Note:** This method does not require any parameters — it simply returns the count of currently open positions.

---

```go
if err != nil {
    helpers.PrintShortError(err, "PositionsTotal failed")
}
```

If an error occurs during the request execution (e.g., no connection or unauthorized user),
it is handled through the `PrintShortError()` helper.

---

```go
fmt.Printf("  Total open positions:          %d\n", positionsTotalData.TotalPositions)
```

If there are no errors, the number of open positions is printed to the console.
For example:

```
Total open positions:          3
```

---

## 📦 What the Server Returns

```protobuf
message PositionsTotalData {
  // Total number of positions
  int32 total_positions = 1;
}
```

---

##  Example Output

```
Total open positions:          2
```

📘 This means that **2 active positions** are currently open on the account (e.g., EURUSD and GBPUSD).

---

### 🧠 What It's Used For

The `PositionsTotal()` method is used:

* to **check for open trades** before opening a new one;
* in **risk management algorithms** that limit the number of active positions;
* during **trading system initialization** to determine the current account state;
* in **monitoring dashboards** to display the total number of trades.

---

### 💬 In Simple Terms

> `PositionsTotal()` is a quick way to find out **how many trades are currently open** on the account.
> The method does not return trade details, only their count — convenient for checks and conditions like:
> *"if no positions, open a new one"*.
