### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolInfoDouble()`** method is used to get **a single numeric property of a symbol**.
> It requests a specific parameter specified in the `Type` field and returns its value in `Value`.


---

## 1️⃣ Getting current Bid price

```go
bidReq := &pb.SymbolInfoDoubleRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_BID, // Current Bid price
}
bidData, err := account.SymbolInfoDouble(ctx, bidReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoDouble(BID) failed")
} else {
    fmt.Printf("  Bid price (SYMBOL_BID):        %.5f\n", bidData.Value)
}
```

### 🟢 Detailed Code Explanation

```go
bidReq := &pb.SymbolInfoDoubleRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_BID,
}
```

Creates a `SymbolInfoDoubleRequest` request with two fields:

* `Symbol` — instrument name, e.g. `EURUSD`, taken from `appsettings.json` → `test_symbol`.
* `Type` — parameter type we want to get. In this case it's `SYMBOL_BID` — current **Bid price**.

---

```go
bidData, err := account.SymbolInfoDouble(ctx, bidReq)
```

Calls the `SymbolInfoDouble()` method via gRPC.

* `ctx` — execution context (manages timeout, cancellation, etc.);
* `bidReq` — prepared request;
* `bidData` — server response with the `Value` field where the required number is stored.

---

```go
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoDouble(BID) failed")
} else {
    fmt.Printf("  Bid price (SYMBOL_BID): %.5f\n", bidData.Value)
}
```

Checks for errors. If everything is successful — prints the value `bidData.Value`.
Format `%.5f` outputs the number with 5 decimal places, as is customary for quotes.

Example result:

```
Bid price (SYMBOL_BID): 1.08540
```

📘 **SYMBOL_BID** — is the current bid price (Bid).
It shows at what price you can sell the instrument at the moment.

---

## 2️⃣ Getting current Ask price

```go
askReq := &pb.SymbolInfoDoubleRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_ASK, // Current Ask price
}
askData, err := account.SymbolInfoDouble(ctx, askReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoDouble(ASK) failed")
} else {
    fmt.Printf("  Ask price (SYMBOL_ASK):        %.5f\n", askData.Value)
}
```

📘 **SYMBOL_ASK** — current ask price (Ask).
Shows at what price you can buy the instrument.
Together with Bid, it forms the spread — the difference between buy and sell prices.

---

## 3️⃣ Getting point size (Point)

```go
pointReq := &pb.SymbolInfoDoubleRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_POINT, // Point size
}
pointData, err := account.SymbolInfoDouble(ctx, pointReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoDouble(POINT) failed")
} else {
    fmt.Printf("  Point size (SYMBOL_POINT):     %.5f\n", pointData.Value)
}
```

📘 **SYMBOL_POINT** — minimum possible price change of the instrument (tick size).
For currency pairs, this value determines quote precision.
Examples: `0.00001` for EURUSD, `0.001` for USDJPY.

---

## 4️⃣ Getting minimum trade volume

```go
volumeMinReq := &pb.SymbolInfoDoubleRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_VOLUME_MIN, // Minimum volume for trading
}
volumeMinData, err := account.SymbolInfoDouble(ctx, volumeMinReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoDouble(VOLUME_MIN) failed")
} else {
    fmt.Printf("  Min volume (VOLUME_MIN):       %.2f\n", volumeMinData.Value)
}
```

📘 **SYMBOL_VOLUME_MIN** — minimum allowed trading volume in lots.
This parameter is set by the broker and determines whether trading is possible, for example, from `0.01` lot.

---

### 📦 What the server returns

After each `SymbolInfoDouble()` call, the gateway returns a simple protobuf structure with one field:

```protobuf
message SymbolInfoDoubleData {
  double Value = 1; // Value of the requested property
}
```

The `Value` field contains the numeric value of the parameter specified in the request (`Type`).
Usually this is a quote, price step, volume, or other market indicator.

---

Below is a list of the most commonly used properties that can be requested via `SymbolInfoDouble()`:

### 📘 All PropertyId enum values (SymbolInfoDoubleProperty)

| Constant | Value | Description |
|----------|-------|-------------|
| `SYMBOL_BID` | 0 | Current Bid price |
| `SYMBOL_BIDHIGH` | 1 | Highest Bid of the day |
| `SYMBOL_BIDLOW` | 2 | Lowest Bid of the day |
| `SYMBOL_ASK` | 3 | Current Ask price |
| `SYMBOL_ASKHIGH` | 4 | Highest Ask of the day |
| `SYMBOL_ASKLOW` | 5 | Lowest Ask of the day |
| `SYMBOL_LAST` | 6 | Last deal price |
| `SYMBOL_LASTHIGH` | 7 | Highest Last price of the day |
| `SYMBOL_LASTLOW` | 8 | Lowest Last price of the day |
| `SYMBOL_VOLUME_REAL` | 9 | Last deal volume |
| `SYMBOL_VOLUMEHIGH_REAL` | 10 | Maximum volume of the day |
| `SYMBOL_VOLUMELOW_REAL` | 11 | Minimum volume of the day |
| `SYMBOL_OPTION_STRIKE` | 12 | Option strike price |
| `SYMBOL_POINT` | 13 | Point size (smallest price change) |
| `SYMBOL_TRADE_TICK_VALUE` | 14 | Tick value in account currency |
| `SYMBOL_TRADE_TICK_VALUE_PROFIT` | 15 | Tick value for profit position |
| `SYMBOL_TRADE_TICK_VALUE_LOSS` | 16 | Tick value for loss position |
| `SYMBOL_TRADE_TICK_SIZE` | 17 | Minimum price change |
| `SYMBOL_TRADE_CONTRACT_SIZE` | 18 | Contract size (lot size in base currency) |
| `SYMBOL_TRADE_ACCRUED_INTEREST` | 19 | Accrued interest |
| `SYMBOL_TRADE_FACE_VALUE` | 20 | Face value (for bonds) |
| `SYMBOL_TRADE_LIQUIDITY_RATE` | 21 | Liquidity rate |
| `SYMBOL_VOLUME_MIN` | 22 | Minimum volume for trading |
| `SYMBOL_VOLUME_MAX` | 23 | Maximum volume for trading |
| `SYMBOL_VOLUME_STEP` | 24 | Volume step (lot size increment) |
| `SYMBOL_VOLUME_LIMIT` | 25 | Maximum total volume of open positions |
| `SYMBOL_SWAP_LONG` | 26 | Swap for long positions |
| `SYMBOL_SWAP_SHORT` | 27 | Swap for short positions |
| `SYMBOL_SWAP_SUNDAY` | 28 | Triple swap day - Sunday |
| `SYMBOL_SWAP_MONDAY` | 29 | Triple swap day - Monday |
| `SYMBOL_SWAP_TUESDAY` | 30 | Triple swap day - Tuesday |
| `SYMBOL_SWAP_WEDNESDAY` | 31 | Triple swap day - Wednesday |
| `SYMBOL_SWAP_THURSDAY` | 32 | Triple swap day - Thursday |
| `SYMBOL_SWAP_FRIDAY` | 33 | Triple swap day - Friday |
| `SYMBOL_SWAP_SATURDAY` | 34 | Triple swap day - Saturday |
| `SYMBOL_MARGIN_INITIAL` | 35 | Initial margin |
| `SYMBOL_MARGIN_MAINTENANCE` | 36 | Maintenance margin |
| `SYMBOL_SESSION_VOLUME` | 37 | Summary volume of the current session |
| `SYMBOL_SESSION_TURNOVER` | 38 | Summary turnover of the current session |
| `SYMBOL_SESSION_INTEREST` | 39 | Summary open interest |
| `SYMBOL_SESSION_BUY_ORDERS_VOLUME` | 40 | Current volume of Buy orders |
| `SYMBOL_SESSION_SELL_ORDERS_VOLUME` | 41 | Current volume of Sell orders |
| `SYMBOL_SESSION_OPEN` | 42 | Open price of the current session |
| `SYMBOL_SESSION_CLOSE` | 43 | Close price of the current session |
| `SYMBOL_SESSION_AW` | 44 | Average weighted price of the current session |
| `SYMBOL_SESSION_PRICE_SETTLEMENT` | 45 | Settlement price of the current session |
| `SYMBOL_SESSION_PRICE_LIMIT_MIN` | 46 | Minimum price of the current session |
| `SYMBOL_SESSION_PRICE_LIMIT_MAX` | 47 | Maximum price of the current session |
| `SYMBOL_MARGIN_HEDGED` | 48 | Hedged margin |
| `SYMBOL_PRICE_CHANGE` | 49 | Change of price in % |
| `SYMBOL_PRICE_VOLATILITY` | 50 | Price volatility in % |
| `SYMBOL_PRICE_THEORETICAL` | 51 | Theoretical option price |
| `SYMBOL_PRICE_DELTA` | 52 | Option/warrant delta |
| `SYMBOL_PRICE_THETA` | 53 | Option/warrant theta |
| `SYMBOL_PRICE_GAMMA` | 54 | Option/warrant gamma |
| `SYMBOL_PRICE_VEGA` | 55 | Option/warrant vega |
| `SYMBOL_PRICE_RHO` | 56 | Option/warrant rho |
| `SYMBOL_PRICE_OMEGA` | 57 | Option/warrant omega |
| `SYMBOL_PRICE_SENSITIVITY` | 58 | Option/warrant sensitivity |
| `SYMBOL_COUNT` | 59 | Total number of properties (not used) |

---

### 🧠 What it's used for

The `SymbolInfoDouble()` method is used:

* to **read market parameters** without needing to request the entire `SymbolInfo` object;
* in **algorithms and strategies** where you need to work with specific numeric values — for example, spread, price step, or tick value;
* for **validating trading conditions** before an operation — check that volume is not less than `VolumeMin`, and the step matches `VolumeStep`;
* for **logging and demo examples**, when you need to output individual symbol parameters to the console.

---

### 💬 In simple terms

> `SymbolInfoDouble()` allows you to quickly get any numeric property of a symbol.
> Each request returns one value — Bid, Ask, Point, volume, etc.
> It's a lightweight and fast call, often used when working with code to read instrument parameters.
