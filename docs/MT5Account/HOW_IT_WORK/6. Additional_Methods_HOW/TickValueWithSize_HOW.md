### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`TickValueWithSize()`** method is used to retrieve **tick value and size information** for trading symbols - essential data for calculating position value, profit/loss per tick, and risk management.
>
> This method returns comprehensive tick data including the standard tick value, separate values for profitable and losing positions, tick size (minimum price change), and contract size.

---

## 🧩 Code example

```go
tickValueReq := &pb.TickValueWithSizeRequest{
    SymbolNames: []string{cfg.TestSymbol}, // Array of symbol names
}
tickValueData, err := account.TickValueWithSize(ctx, tickValueReq)
if err != nil {
    helpers.PrintShortError(err, "TickValueWithSize failed")
} else {
    fmt.Printf("  Retrieved tick value/size data for %d symbols:\n", len(tickValueData.SymbolTickSizeInfos))
    for _, info := range tickValueData.SymbolTickSizeInfos {
        fmt.Printf("\n  Symbol: %s (Index: %d)\n", info.Name, info.Index)
        fmt.Printf("    Trade Tick Value:        %.5f\n", info.TradeTickValue)
        fmt.Printf("    Trade Tick Value Profit: %.5f\n", info.TradeTickValueProfit)
        fmt.Printf("    Trade Tick Value Loss:   %.5f\n", info.TradeTickValueLoss)
        fmt.Printf("    Trade Tick Size:         %.5f\n", info.TradeTickSize)
        fmt.Printf("    Trade Contract Size:     %.2f\n", info.TradeContractSize)
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
tickValueReq := &pb.TickValueWithSizeRequest{
    SymbolNames: []string{cfg.TestSymbol}, // Array of symbol names
}
```

A request is created with parameters:

* **`SymbolNames`** - array of trading instrument names (e.g., `["EURUSD", "GBPUSD"]`);
* You can request data for multiple symbols in a single call for efficiency.

---

```go
tickValueData, err := account.TickValueWithSize(ctx, tickValueReq)
```

A call is made to the MetaTrader server via gRPC. The response contains an array of `TickSizeSymbol` structures with detailed tick and contract information for each requested symbol.

---

```go
fmt.Printf("  Retrieved tick value/size data for %d symbols:\n", len(tickValueData.SymbolTickSizeInfos))
```

The method returns a slice `SymbolTickSizeInfos` containing information for all requested symbols. We print the count of symbols for which data was retrieved.

---

```go
for _, info := range tickValueData.SymbolTickSizeInfos {
    fmt.Printf("\n  Symbol: %s (Index: %d)\n", info.Name, info.Index)
```

Iterate through each symbol's tick information:

* **`Name`** - symbol name (e.g., "EURUSD")
* **`Index`** - symbol's position in the request array

---

```go
fmt.Printf("    Trade Tick Value:        %.5f\n", info.TradeTickValue)
fmt.Printf("    Trade Tick Value Profit: %.5f\n", info.TradeTickValueProfit)
fmt.Printf("    Trade Tick Value Loss:   %.5f\n", info.TradeTickValueLoss)
```

Print tick values in account currency:

* **`TradeTickValue`** - standard tick value (cost of one tick movement)
* **`TradeTickValueProfit`** - tick value for profitable positions (may differ due to asymmetric pricing)
* **`TradeTickValueLoss`** - tick value for losing positions

---

```go
fmt.Printf("    Trade Tick Size:         %.5f\n", info.TradeTickSize)
fmt.Printf("    Trade Contract Size:     %.2f\n", info.TradeContractSize)
```

Print additional parameters:

* **`TradeTickSize`** - minimum price change (tick size), typically 0.00001 for 5-digit currency pairs
* **`TradeContractSize`** - contract size for the symbol (e.g., 100,000 for standard forex lots)

---

## 📦 What the Server Returns

```protobuf
message TickValueWithSizeData {
  repeated TickSizeSymbol symbol_tick_size_infos = 1;
}

message TickSizeSymbol {
  int32 Index = 1;                    // Symbol index in request array
  string Name = 2;                    // Symbol name
  double TradeTickValue = 3;          // Standard tick value in account currency
  double TradeTickValueProfit = 4;    // Tick value for profitable positions
  double TradeTickValueLoss = 5;      // Tick value for losing positions
  double TradeTickSize = 6;           // Minimum price change (tick size)
  double TradeContractSize = 7;       // Contract size for the symbol
}
```

---

## 💡 Example Output

```
Retrieved tick value/size data for 1 symbols:

  Symbol: EURUSD (Index: 0)
    Trade Tick Value:        0.84107
    Trade Tick Value Profit: 0.84107
    Trade Tick Value Loss:   0.84109
    Trade Tick Size:         0.00001
    Trade Contract Size:     100000.00
```

This means:
* One tick movement (0.00001) for EURUSD costs 0.84107 EUR in your account currency
* Profitable and losing positions have nearly identical tick values (difference: 0.00002 EUR)
* The contract size is 100,000 units (standard lot)

---

###  Why Tick Value Matters

The tick value is critical for:

1. **Position Value Calculation**
   ```
   Position Value = Lots * Contract Size * Current Price
   ```

2. **Risk Per Tick**
   ```
   Risk Per Tick = Lots * Tick Value
   ```

3. **Profit/Loss Calculation**
   ```
   P&L = Price Movement (in ticks) * Lots * Tick Value
   ```

4. **Stop Loss in Currency**
   ```
   Stop Loss Amount = Stop Distance (pips) * Ticks Per Pip * Lots * Tick Value
   ```

---

### Asymmetric Tick Values

Some symbols have different tick values for profitable vs losing positions. For example:

```
Symbol: XAUUSD (Gold)
  Trade Tick Value Profit: 1.00
  Trade Tick Value Loss:   1.02
```

This 2% difference occurs because:
* Brokers may apply different conversion rates for buy/sell positions
* Cross-currency calculations can introduce asymmetry
* Some instruments have built-in spread in tick value calculation

Always use `TradeTickValueProfit` for buy positions and `TradeTickValueLoss` for sell positions to ensure accurate P&L calculations.

---

### 🧠 What It's Used For

The `TickValueWithSize()` method is used:

* to **calculate exact position value** in account currency;
* for **risk management** - determining profit/loss per tick movement;
* when **sizing positions** based on risk percentage;
* in **P&L calculators** and trading panels;
* for **lot size calculations** based on stop loss distance;
* in **portfolio risk analysis** across multiple symbols.

---

### 💬 In Simple Terms

> `TickValueWithSize()` tells you **how much money you make or lose per tick** for each symbol.
> It provides all the data needed to calculate position sizes, risk amounts, and profit/loss accurately in your account currency.
