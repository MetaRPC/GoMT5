### Example from file: `examples/demos/lowlevel/02_trading_operations.go`

> The **`OrderCalcMargin()`** method is used to calculate the required margin before opening an order.
> It allows you to determine in advance whether there are sufficient funds in the account for a trade of a given volume and type.


---

## 🧩 Code example

```go
fmt.Println("\n2.1. OrderCalcMargin() - Calculate required margin")

// Get current prices for the symbol
tickReq := &pb.SymbolInfoTickRequest{
    Symbol: cfg.TestSymbol,
}
tickData, err := account.SymbolInfoTick(ctx, tickReq)
helpers.Fatal(err, "SymbolInfoTick failed")

// Form a margin calculation request for BUY
marginReq := &pb.OrderCalcMarginRequest{
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Symbol:    cfg.TestSymbol,
    Volume:    cfg.TestVolume,
    OpenPrice: tickData.Ask, // use current Ask price
}

marginData, err := account.OrderCalcMargin(ctx, marginReq)
if !helpers.PrintShortError(err, "OrderCalcMargin failed") {
    fmt.Printf("  Symbol:                        %s\n", cfg.TestSymbol)
    fmt.Printf("  Action:                        BUY\n")
    fmt.Printf("  Volume:                        %.2f lots\n", cfg.TestVolume)
    fmt.Printf("  Price:                         %.5f\n", tickData.Ask)
    fmt.Printf("  Required Margin:               %.2f\n", marginData.Margin)
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Get Current Quotes

```go
tickReq := &pb.SymbolInfoTickRequest{ Symbol: cfg.TestSymbol }
tickData, err := account.SymbolInfoTick(ctx, tickReq)
```

Before calculating, you need to know the actual price. For a **BUY** trade, the **Ask** price is used.
`cfg.TestSymbol` — symbol from `config.json`.

---

### 2️. Form the Request

```go
marginReq := &pb.OrderCalcMarginRequest{
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Symbol:    cfg.TestSymbol,
    Volume:    cfg.TestVolume,
    OpenPrice: tickData.Ask,
}
```

* **OrderType** — order type (BUY / SELL);
* **Symbol** — instrument (e.g., `EURUSD`);
* **Volume** — trade volume in lots;
* **OpenPrice** — opening price (Ask for BUY, Bid for SELL).

---

### 3️. Perform the Calculation

```go
marginData, err := account.OrderCalcMargin(ctx, marginReq)
```

A gRPC request is sent to the gateway. MetaTrader calculates how much margin will be required for a trade with the specified parameters.

---

### 4️. Process and Display the Result

```go
fmt.Printf("  Required Margin: %.2f\n", marginData.Margin)
```

Outputs the calculated margin amount in the account currency.

---

## 📦 What the Server Returns

```protobuf
message OrderCalcMarginData {
  double Margin = 1; // Calculated margin for the specified order
}
```

---

## 💡 Example Output

```
Symbol:                        EURUSD
Action:                        BUY
Volume:                        0.10 lots
Price:                         1.08525
Required Margin:               108.52
```

---

## 🧠 What It's Used For

The `OrderCalcMargin()` method is used for:

* checking if there is enough free margin for a new trade;
* dynamically adjusting trade volume to balance;
* assessing deposit load and risk level;
* demonstrative calculations without opening orders.

---

## 💬 In Simple Terms

> `OrderCalcMargin()` is a calculator that tells you:
> "How much money the broker will freeze if you open a BUY/SELL trade with this volume and price now."
> It doesn't create a trade — it only calculates **how much margin will be required**.
