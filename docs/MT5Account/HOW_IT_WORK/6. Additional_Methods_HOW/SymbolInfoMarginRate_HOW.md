### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`SymbolInfoMarginRate()`** method is used to retrieve **margin requirements for a symbol** — that is, how much funds are needed to open and maintain a position for a specific instrument.
>
> It allows you to find out two key parameters:
>
> * **InitialMarginRate** — coefficient for opening a position;
> * **MaintenanceMarginRate** — coefficient for maintaining a position.


---

## 🧩 Code example

```go
marginRateReq := &pb.SymbolInfoMarginRateRequest{
    Symbol:    cfg.TestSymbol,
    OrderType: pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY, // Direct protobuf enum
}
marginRateData, err := account.SymbolInfoMarginRate(ctx, marginRateReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoMarginRate failed")
} else {
    fmt.Printf("  Initial margin rate (BUY):     %.2f\n", marginRateData.InitialMarginRate)
    fmt.Printf("  Maintenance margin rate (BUY): %.2f\n", marginRateData.MaintenanceMarginRate)
}
```

### 🟢 Detailed Code Explanation

```go
marginRateReq := &pb.SymbolInfoMarginRateRequest{
    Symbol:    cfg.TestSymbol,
    OrderType: pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY,
}
```

A `SymbolInfoMarginRateRequest` request is created with parameters:

* `Symbol` — trading instrument (e.g., `EURUSD`);
* `OrderType` — order type (`BUY`, `SELL`, `BUY_LIMIT`, etc.).
  Different order types can have different margin coefficients.

---

```go
marginRateData, err := account.SymbolInfoMarginRate(ctx, marginRateReq)
```

The `SymbolInfoMarginRate()` method is called via gRPC.
The result is stored in `marginRateData`.

---

```go
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoMarginRate failed")
} else {
    fmt.Printf("  Initial margin rate (BUY):     %.2f\n", marginRateData.InitialMarginRate)
    fmt.Printf("  Maintenance margin rate (BUY): %.2f\n", marginRateData.MaintenanceMarginRate)
}
```

If there are no errors, both values are displayed:

* `InitialMarginRate` — initial margin coefficient (e.g., 0.02 → 2%);
* `MaintenanceMarginRate` — maintenance margin coefficient (e.g., 0.01 → 1%).

---

## 📦 What the Server Returns

After the call, the method returns a protobuf structure:

```protobuf
message SymbolInfoMarginRateData {
  double InitialMarginRate = 1;      // Initial margin for opening a position
  double MaintenanceMarginRate = 2;  // Maintenance margin for holding a position
}
```

---

## 💡 Example of Value Interpretation

| Field                          | Value | Meaning                                                  |
| ------------------------------ | ----- | -------------------------------------------------------- |
| `InitialMarginRate = 0.02`     | 2%    | For a trade volume of 100,000 EURUSD, 2,000 EUR collateral is needed |
| `MaintenanceMarginRate = 0.01` | 1%    | To hold the position, 1,000 EUR is required              |

---

### 🧠 What It's Used For

The `SymbolInfoMarginRate()` method is used:

* to **calculate required margin** before opening a position;
* to **check risks and available funds** in trading algorithms;
* when **modeling trades** to assess the impact of leverage;
* in **demo examples** where you need to show real trading parameters of a symbol.

---

### 💬 In Simple Terms

> `SymbolInfoMarginRate()` tells you **how much collateral is needed** to open and hold a trade for a symbol.
> It's used for risk calculations, limit checks, and choosing the optimal position size.
