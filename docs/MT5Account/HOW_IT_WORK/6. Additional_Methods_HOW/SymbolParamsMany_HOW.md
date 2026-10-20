### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`SymbolParamsMany()`** method is used to retrieve **detailed parameters for one symbol or for all symbols at once**.
> It allows you to get a wide set of properties — prices, volumes, steps, spread, number of digits, and much more — in a single request.
>
> Unlike the `SymbolInfoDouble`, `SymbolInfoInteger`, and `SymbolInfoString` methods, which return a single value,
> `SymbolParamsMany()` makes a bulk request and returns a list of `SymbolParameters` structures (`SymbolInfos`).


---

## 🧩 Code example

```go
// SymbolParamsMany has no array of symbol names: SymbolName filters one symbol (nil = all symbols)
symbolFilter := cfg.TestSymbol
paramsManyReq := &pb.SymbolParamsManyRequest{
    SymbolName:   &symbolFilter, // Filter by symbol name (optional pointer to string)
    SortType:     nil,           // Optional sort type
    PageNumber:   nil,           // Optional page number
    ItemsPerPage: nil,           // Optional items per page
}
paramsManyData, err := account.SymbolParamsMany(ctx, paramsManyReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolParamsMany failed")
} else {
    fmt.Printf("  Retrieved parameters for %d symbols:\n", len(paramsManyData.SymbolInfos))
    // Show only first 3 symbols
    maxShow := 3
    if len(paramsManyData.SymbolInfos) < maxShow {
        maxShow = len(paramsManyData.SymbolInfos)
    }
    for i := 0; i < maxShow; i++ {
        params := paramsManyData.SymbolInfos[i]
        fmt.Printf("\n  Symbol #%d: %s\n", i+1, params.Name)
        fmt.Printf("    Bid:                         %.5f\n", params.Bid)
        fmt.Printf("    Ask:                         %.5f\n", params.Ask)
        fmt.Printf("    Digits:                      %d\n", params.Digits)
        fmt.Printf("    Spread:                      %d points\n", params.Spread)
        fmt.Printf("    Volume Min:                  %.2f\n", params.VolumeMin)
        fmt.Printf("    Volume Max:                  %.2f\n", params.VolumeMax)
        fmt.Printf("    Volume Step:                 %.2f\n", params.VolumeStep)
        fmt.Printf("    Contract Size:               %.2f\n", params.TradeContractSize)
        fmt.Printf("    Point:                       %.5f\n", params.Point)
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
symbolFilter := cfg.TestSymbol
paramsManyReq := &pb.SymbolParamsManyRequest{
    SymbolName: &symbolFilter,
}
```

A `SymbolParamsManyRequest` request structure is created.
The optional `SymbolName` field (a `*string`) limits the result to one symbol; leave it `nil` to get all symbols. There is no list-of-symbols field: to get several specific symbols, request all and filter on the client side. `SortType`, `PageNumber` and `ItemsPerPage` (all optional pointers) control sorting and paging.

---

```go
paramsManyData, err := account.SymbolParamsMany(ctx, paramsManyReq)
```

A gRPC request is sent to the MetaTrader server.
The response contains a `SymbolInfos` list — an array of `SymbolParameters` structures, where each describes one symbol with a complete set of parameters.

---

```go
fmt.Printf("  Retrieved parameters for %d symbols:\n", len(paramsManyData.SymbolInfos))
```

Displays the number of symbols returned.

---

```go
maxShow := 3
if len(paramsManyData.SymbolInfos) < maxShow {
    maxShow = len(paramsManyData.SymbolInfos)
}
```

Limits the number of displayed results to 3 — for clarity in the example.

---

```go
for i := 0; i < maxShow; i++ {
    params := paramsManyData.SymbolInfos[i]
    fmt.Printf("\n  Symbol #%d: %s\n", i+1, params.Name)
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
fmt.Printf("    Contract Size:               %.2f\n", params.TradeContractSize)
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
  repeated SymbolParameters symbol_infos = 1; // List of symbol parameters
  int32 symbols_total = 2;
  optional int32 page_number = 3;
  optional int32 items_per_page = 4;
}

message SymbolParameters {
  string name = 1;
  double bid = 2;
  double ask = 5;
  double point = 15;
  double trade_contract_size = 20;
  double volume_min = 24;
  double volume_max = 25;
  double volume_step = 26;
  int32 digits = 78;
  int32 spread = 80;
  // ... 116 fields in total (currency, description, trading status, margins, etc)
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

> `SymbolParamsMany()` allows you to get **all main parameters of one symbol or of all symbols in one request**.
> Convenient for loading market data, filtering, and analyzing a large number of instruments.
