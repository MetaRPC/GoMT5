### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolInfoInteger()`** method is used to get **integer properties of a symbol** — parameters that describe market or trading characteristics of the instrument and are expressed in `int32` format.
>
> It works similarly to `SymbolInfoDouble()`, but returns integers: number of decimal places, current spread, trading restrictions, etc.


---

## 1️⃣ Getting number of decimal places (Digits)

```go
digitsReq := &pb.SymbolInfoIntegerRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS, // Number of decimal places
}
digitsData, err := account.SymbolInfoInteger(ctx, digitsReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoInteger(DIGITS) failed")
} else {
    fmt.Printf("  Digits (SYMBOL_DIGITS):        %d\n", digitsData.Value)
}
```

### 🟢 Detailed Code Explanation

```go
digitsReq := &pb.SymbolInfoIntegerRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS,
}
```

Creates a `SymbolInfoIntegerRequest` request with two main fields:

* `Symbol` — instrument name, e.g. `EURUSD` (taken from the `appsettings.json` file, field `test_symbol`);
* `Type` — property type to get. In this case it's `SYMBOL_DIGITS` — number of decimal places in the quote.

---

```go
digitsData, err := account.SymbolInfoInteger(ctx, digitsReq)
```

Calls the `SymbolInfoInteger()` method via gRPC.

* `ctx` — request context (manages execution time);
* `digitsReq` — request structure;
* `digitsData` — response containing the `Value` field (integer).

---

```go
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoInteger(DIGITS) failed")
} else {
    fmt.Printf("  Digits (SYMBOL_DIGITS):        %d\n", digitsData.Value)
}
```

If there's no error, the number of decimal places (`Value`) is printed.
This value determines the precision with which prices are displayed:

* 5 for EURUSD (1.08540)
* 3 for USDJPY (145.123)

---

## 2️⃣ Getting current spread (Spread)

```go
spreadReq := &pb.SymbolInfoIntegerRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoIntegerProperty_SYMBOL_SPREAD,
}
spreadData, err := account.SymbolInfoInteger(ctx, spreadReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoInteger(SPREAD) failed")
} else {
    fmt.Printf("  Spread (SYMBOL_SPREAD):        %d points\n", spreadData.Value)
}
```

📘 **SYMBOL_SPREAD** — current spread expressed in points.
It shows the difference between `Ask` and `Bid` prices.
For example, if `Bid = 1.08540`, `Ask = 1.08560`, and point size is `0.00001`, then the spread equals 20 points.

---

## 3️⃣ Getting stops level (Stops Level)

```go
stopsLevelReq := &pb.SymbolInfoIntegerRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoIntegerProperty_SYMBOL_TRADE_STOPS_LEVEL,
}
stopsLevelData, err := account.SymbolInfoInteger(ctx, stopsLevelReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoInteger(STOPS_LEVEL) failed")
} else {
    fmt.Printf("  Stops level (STOPS_LEVEL):     %d points\n", stopsLevelData.Value)
}
```

📘 **SYMBOL_TRADE_STOPS_LEVEL** — minimum distance (in points) at which orders of type `Stop Loss` or `Take Profit` can be placed from the current price.
For example, if `Stops Level = 10` and price is 1.08540 — stops cannot be set closer than 10 points.

---

### 📦 What the server returns

After each `SymbolInfoInteger()` call, the gateway returns a protobuf structure:

```protobuf
message SymbolInfoIntegerData {
  int64 Value = 1; // Integer value of the property
}
```

The `Value` field contains the numeric value of the selected parameter (`Type`).

---

### 📘 All PropertyId enum values (SymbolInfoIntegerProperty)

| Constant | Value | Description |
|----------|-------|-------------|
| `SYMBOL_SUBSCRIPTION_DELAY` | 0 | Subscription delay in milliseconds |
| `SYMBOL_SECTOR` | 1 | Economic sector (enum SYMBOL_SECTOR) |
| `SYMBOL_INDUSTRY` | 2 | Industry (enum SYMBOL_INDUSTRY) |
| `SYMBOL_CUSTOM` | 3 | Custom symbol flag (0=standard, 1=custom) |
| `SYMBOL_BACKGROUND_COLOR` | 4 | Background color in Market Watch |
| `SYMBOL_CHART_MODE` | 5 | Price chart mode (SYMBOL_CHART_MODE enum) |
| `SYMBOL_EXIST` | 6 | Symbol exists (0=no, 1=yes) |
| `SYMBOL_SELECT` | 7 | Symbol selected in Market Watch (0=no, 1=yes) |
| `SYMBOL_VISIBLE` | 8 | Symbol visible in Market Watch (0=no, 1=yes) |
| `SYMBOL_SESSION_DEALS` | 9 | Number of deals in current session |
| `SYMBOL_SESSION_BUY_ORDERS` | 10 | Number of Buy orders in current session |
| `SYMBOL_SESSION_SELL_ORDERS` | 11 | Number of Sell orders in current session |
| `SYMBOL_VOLUME` | 12 | Last deal volume |
| `SYMBOL_VOLUMEHIGH` | 13 | Maximum volume of the day |
| `SYMBOL_VOLUMELOW` | 14 | Minimum volume of the day |
| `SYMBOL_TIME` | 15 | Last quote time (seconds since 1970.01.01) |
| `SYMBOL_TIME_MSC` | 16 | Last quote time in milliseconds |
| `SYMBOL_DIGITS` | 17 | Number of decimal places |
| `SYMBOL_SPREAD_FLOAT` | 18 | Floating spread flag (0=fixed, 1=floating) |
| `SYMBOL_SPREAD` | 19 | Current spread in points |
| `SYMBOL_TICKS_BOOKDEPTH` | 20 | Maximum depth of market (DOM) |
| `SYMBOL_TRADE_CALC_MODE` | 21 | Margin calculation mode (SYMBOL_CALC_MODE enum) |
| `SYMBOL_TRADE_MODE` | 22 | Trade execution mode (SYMBOL_TRADE_MODE enum) |
| `SYMBOL_START_TIME` | 23 | Date of symbol trade beginning |
| `SYMBOL_EXPIRATION_TIME` | 24 | Date of symbol trade end |
| `SYMBOL_TRADE_STOPS_LEVEL` | 25 | Minimum distance for SL/TP in points |
| `SYMBOL_TRADE_FREEZE_LEVEL` | 26 | Order freeze level in points |
| `SYMBOL_TRADE_EXEMODE` | 27 | Order execution type (SYMBOL_ORDER_EXEMODE enum) |
| `SYMBOL_SWAP_MODE` | 28 | Swap calculation mode (SYMBOL_SWAP_MODE enum) |
| `SYMBOL_SWAP_ROLLOVER3DAYS` | 29 | Day of week for triple swap (SYMBOL_SWAP_ROLLOVER enum) |
| `SYMBOL_MARGIN_HEDGED_USE_LEG` | 30 | Calculating hedged margin using larger leg |
| `SYMBOL_EXPIRATION_MODE` | 31 | Allowed order expiration modes |
| `SYMBOL_FILLING_MODE` | 32 | Allowed order filling modes |
| `SYMBOL_ORDER_MODE` | 33 | Allowed order types |
| `SYMBOL_ORDER_GTC_MODE` | 34 | StopLoss/TakeProfit mode (ORDER_GTC_MODE enum) |
| `SYMBOL_OPTION_MODE` | 35 | Option type (SYMBOL_OPTION_MODE enum) |
| `SYMBOL_OPTION_RIGHT` | 36 | Option right (Call/Put) |

---

### 🧠 What it's used for

The `SymbolInfoInteger()` method is used:

* for **analyzing trading conditions** (spread, number of decimal places, stop levels);
* for **checking symbol settings** before making a trade;
* for **debugging the gateway and demo examples**, when you need to output numeric symbol parameters to the console;
* when **building trading strategies** that require precise work with points and price steps.

---

### 💬 In simple terms

> `SymbolInfoInteger()` returns **integer instrument parameters** — such as number of decimal places, spread, or minimum distance to stops.
> It's a quick way to get important symbol settings without requesting the entire data structure.
