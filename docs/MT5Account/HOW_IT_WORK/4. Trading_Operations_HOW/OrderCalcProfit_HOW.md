### Example from file: `examples/demos/lowlevel/02_trading_operations.go`

> The **`OrderCalcProfit()`** method is used to calculate the potential profit or loss for a trade with specified parameters.
> It does not open a trade — it only calculates the result when opening and closing a position at the specified prices.


---

## 🧩 Code example

```go
fmt.Println("\n2.2. OrderCalcProfit() - Calculate potential profit/loss")

// Calculate loss from immediate trade closure (spread)
profitReq := &pb.OrderCalcProfitRequest{
    OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Symbol:     cfg.TestSymbol,
    Volume:     cfg.TestVolume,
    OpenPrice:  tickData.Ask, // entry price
    ClosePrice: tickData.Bid, // exit price - immediate close at Bid
}
profitData, err := account.OrderCalcProfit(ctx, profitReq)
if !helpers.PrintShortError(err, "OrderCalcProfit failed") {
    fmt.Printf("  Symbol:                        %s\n", cfg.TestSymbol)
    fmt.Printf("  Action:                        BUY\n")
    fmt.Printf("  Volume:                        %.2f lots\n", cfg.TestVolume)
    fmt.Printf("  Price Open (Ask):              %.5f\n", tickData.Ask)
    fmt.Printf("  Price Close (Bid):             %.5f\n", tickData.Bid)
    fmt.Printf("  Potential Profit/Loss:         %.2f (spread loss)\n", profitData.Profit)
}

// Calculate profit if price rises by 10 pips.
// A pip is 10 points on 5- and 3-digit symbols (EURUSD: point 0.00001, pip 0.0001; USDJPY: pip 0.01),
// otherwise 1 point: read the symbol's point and digits instead of assuming EURUSD.
pointData, err := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_POINT,
})
if err != nil {
    return fmt.Errorf("SymbolInfoDouble(POINT) failed: %w", err)
}
digitsData, err := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS,
})
if err != nil {
    return fmt.Errorf("SymbolInfoInteger(DIGITS) failed: %w", err)
}
pipSize := pointData.Value
if digitsData.Value == 3 || digitsData.Value == 5 {
    pipSize *= 10
}
targetPrice := tickData.Ask + (10 * pipSize)

profitTargetReq := &pb.OrderCalcProfitRequest{
    OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Symbol:     cfg.TestSymbol,
    Volume:     cfg.TestVolume,
    OpenPrice:  tickData.Ask,
    ClosePrice: targetPrice,
}
profitTargetData, err := account.OrderCalcProfit(ctx, profitTargetReq)
if !helpers.PrintShortError(err, "OrderCalcProfit (target) failed") {
    fmt.Printf("\n  If price moves +10 pips to %.5f:\n", targetPrice)
    fmt.Printf("  Potential Profit:              %.2f\n", profitTargetData.Profit)
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Calculating Immediate Loss Due to Spread

When opening a trade at Ask and closing at Bid, there is an immediate loss equal to the spread size.

```go
OpenPrice: tickData.Ask,
ClosePrice: tickData.Bid,
```

The method returns a negative profit value (`Profit < 0`).

---

### 2️. Calculating Profit from 10 Pips Price Movement

Define the pip size and target price:

```go
pipSize := pointData.Value // SYMBOL_POINT
if digitsData.Value == 3 || digitsData.Value == 5 { // SYMBOL_DIGITS
    pipSize *= 10
}
targetPrice := tickData.Ask + (10 * pipSize)
```

A pip is not a fixed number: on 5- and 3-digit quotes it is 10 points (EURUSD 1.12345: point 0.00001, pip 0.0001;
USDJPY 151.234: pip 0.01), otherwise one point. The demo therefore reads the symbol's `SYMBOL_POINT` and
`SYMBOL_DIGITS` instead of hard-coding the EURUSD value, so `cfg.TestSymbol` can be any symbol.
If the price rises by 10 pips (0.0010 for EURUSD), the result will be positive.
The method calculates how much profit this movement will generate for the specified volume.

---

## 📦 What the Server Returns

```protobuf
message OrderCalcProfitData {
  double Profit = 1; // Calculated profit or loss in account currency
}
```

---

## 💡 Example Output

```
Symbol:                        EURUSD
Action:                        BUY
Volume:                        0.10 lots
Price Open (Ask):              1.08520
Price Close (Bid):             1.08500
Potential Profit/Loss:         -2.00 (spread loss)

If price moves +10 pips to 1.08620:
Potential Profit:              10.00
```

---

## 🧠 What It's Used For

The `OrderCalcProfit()` method is used for:

* estimating potential profit or loss without opening a trade;
* testing strategies and risk management;
* calculating P/L in simulations;
* visualizing theoretical results.

---

## 💬 In Simple Terms

> `OrderCalcProfit()` is a profit calculator.
> It shows **how much you will earn or lose** if you buy at one price and sell at another.
> This is a pure calculation, without actual order execution.
