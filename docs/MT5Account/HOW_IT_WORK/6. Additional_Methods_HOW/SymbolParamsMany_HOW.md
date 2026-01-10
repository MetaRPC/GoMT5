### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`SymbolParamsMany()`** method is used to retrieve **detailed parameters for multiple symbols at once**.
> It allows you to get a wide set of properties — prices, volumes, steps, spread, number of digits, and much more — in a single request.
>
> Unlike the `SymbolInfoDouble`, `SymbolInfoInteger`, and `SymbolInfoString` methods, which return a single value,
> `SymbolParamsMany()` makes a bulk request and returns a full list of `SymbolInfo` structures.


---

## 🧩 Code example

```go
// SymbolParamsMany accepts an array of symbol names
paramsManyReq := &pb.SymbolParamsManyRequest{
    Symbols: []string{cfg.TestSymbol}, // Array of symbol names
}
paramsManyData, err := account.SymbolParamsMany(ctx, paramsManyReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolParamsMany failed")
} else {
    fmt.Printf("  Retrieved parameters for %d symbols:\n", len(paramsManyData.Params))
    // Show only first 3 symbols
    maxShow := 3
    if len(paramsManyData.Params) < maxShow {
        maxShow = len(paramsManyData.Params)
    }
    for i := 0; i < maxShow; i++ {
        params := paramsManyData.Params[i]
        fmt.Printf("\n  Symbol #%d: %s\n", i+1, params.Symbol)
        fmt.Printf("    Bid:                         %.5f\n", params.Bid)
        fmt.Printf("    Ask:                         %.5f\n", params.Ask)
        fmt.Printf("    Digits:                      %d\n", params.Digits)
        fmt.Printf("    Spread:                      %d points\n", params.Spread)
        fmt.Printf("    Volume Min:                  %.2f\n", params.VolumeMin)
        fmt.Printf("    Volume Max:                  %.2f\n", params.VolumeMax)
        fmt.Printf("    Volume Step:                 %.2f\n", params.VolumeStep)
        fmt.Printf("    Contract Size:               %.2f\n", params.ContractSize)
        fmt.Printf("    Point:                       %.5f\n", params.Point)
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
paramsManyReq := &pb.SymbolParamsManyRequest{
    Symbols: []string{cfg.TestSymbol},
}
```

A `SymbolParamsManyRequest` request structure is created.
The `Symbols` field contains an array of symbol names you want to retrieve parameters for (e.g., `[]string{"EURUSD", "GBPUSD", "USDJPY"}`).

---

```go
paramsManyData, err := account.SymbolParamsMany(ctx, paramsManyReq)
```

A gRPC request is sent to the MetaTrader server.
The response contains a `Params` list — an array of structures, where each describes one symbol with a complete set of parameters.

---

```go
fmt.Printf("  Retrieved parameters for %d symbols:\n", len(paramsManyData.Params))
```

Displays the number of symbols returned.

---

```go
maxShow := 3
if len(paramsManyData.Params) < maxShow {
    maxShow = len(paramsManyData.Params)
}
```

Limits the number of displayed results to 3 — for clarity in the example.

---

```go
for i := 0; i < maxShow; i++ {
    params := paramsManyData.Params[i]
    fmt.Printf("\n  Symbol #%d: %s\n", i+1, params.Symbol)
```

The loop goes through each symbol, displaying its name and characteristics.

---

```go
fmt.Printf("    Bid:                         %.5f\n", params.Bid)
fmt.Printf("    Ask:                         %.5f\n", params.Ask)
fmt.Printf("    Digits:                      %d\n", params.Digits)
fmt.Printf("    Spread:                      %d points\n", params.Spread)
fmt.Printf("    Volume Min:                  %.2f\n", params.VolumeMin)
fmt.Printf("    Volume Max:                  %.2f\n", params.VolumeMax)
fmt.Printf("    Volume Step:                 %.2f\n", params.VolumeStep)
fmt.Printf("    Contract Size:               %.2f\n", params.ContractSize)
fmt.Printf("    Point:                       %.5f\n", params.Point)
```

These lines sequentially display the main symbol parameters:

* **Bid / Ask** — buy and sell prices;
* **Digits** — number of decimal places;
* **Spread** — difference between Ask and Bid in points;
* **Volume Min / Max / Step** — minimum, maximum, and step for trade volume;
* **Contract Size** — contract size in base currency (e.g., 100,000 for forex);
* **Point** — minimum price change.

---

## 📦 What the Server Returns

```protobuf
message SymbolParamsManyData {
  repeated SymbolParams Params = 1; // List of symbol parameters
}

message SymbolParams {
  string Symbol = 1;
  double Bid = 2;
  double Ask = 3;
  int32 Digits = 4;
  int32 Spread = 5;
  double VolumeMin = 6;
  double VolumeMax = 7;
  double VolumeStep = 8;
  double ContractSize = 9;
  double Point = 10;
  // ... over 112 other fields (currency, description, trading status, margins, etc)
}
```

---

## 💡 Example Output

```
Retrieved parameters for 1 symbols:

  Symbol #1: EURUSD
    Bid:                         1.08540
    Ask:                         1.08560
    Digits:                      5
    Spread:                      20 points
    Volume Min:                  0.01
    Volume Max:                  100.00
    Volume Step:                 0.01
    Contract Size:               100000.00
    Point:                       0.00010
```

---

### 🧠 What It's Used For

The `SymbolParamsMany()` method is used:

* for **initializing a trading terminal or strategy**, loading all instrument parameters;
* for **mass market monitoring** when you need to immediately get Bid/Ask for many assets;
* for **comparing symbol characteristics** (spread, volume, price step);
* in **demo examples and UI** to visualize a list of instruments with key parameters.

---

### 💬 In Simple Terms

> `SymbolParamsMany()` allows you to get **all main parameters of multiple symbols in one request**.
> Convenient for loading market data, filtering, and analyzing a large number of instruments.
