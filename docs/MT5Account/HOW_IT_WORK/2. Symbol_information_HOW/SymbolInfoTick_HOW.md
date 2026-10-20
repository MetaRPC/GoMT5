### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolInfoTick()`** method is used to get **the last market tick** — that is, the current quote for a symbol, containing Bid, Ask, Last prices, volume, and time of the last update.
>
> This is one of the key methods for getting "live" market data.


---

## 🧩 Code Example

```go
tickReq := &pb.SymbolInfoTickRequest{
    Symbol: cfg.TestSymbol,
}
tickData, err := account.SymbolInfoTick(ctx, tickReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoTick failed")
} else {
    fmt.Printf("  Last tick for %s:\n", cfg.TestSymbol)
    fmt.Printf("    Bid:                         %.5f\n", tickData.Bid)
    fmt.Printf("    Ask:                         %.5f\n", tickData.Ask)
    fmt.Printf("    Last:                        %.5f\n", tickData.Last)
    fmt.Printf("    Volume:                      %d\n", tickData.Volume)

    // Time is int64 (Unix seconds), NOT protobuf Timestamp. It counts in the broker server's time zone:
    // format it as UTC to show the broker's clock (time.Unix alone would convert it to this machine's zone).
    // TimeUtc is the same moment in real UTC.
    tickTime := time.Unix(tickData.Time, 0).UTC()
    fmt.Printf("    Time:                        %s (broker server time)\n", tickTime.Format("2006-01-02 15:04:05"))
    fmt.Printf("    Time (UTC):                  %s\n", time.Unix(tickData.TimeUtc, 0).UTC().Format("2006-01-02 15:04:05"))
}
```

---

### 🟢 Detailed Code Explanation

```go
tickReq := &pb.SymbolInfoTickRequest{
    Symbol: cfg.TestSymbol,
}
```

Creates a request with one parameter `Symbol` — the name of the instrument (e.g., `EURUSD`) for which to get the last tick.

---

```go
tickData, err := account.SymbolInfoTick(ctx, tickReq)
```

Executes a gRPC request to the gateway. Returns a structure with the last update for the specified symbol: Bid, Ask, Last prices, volume, and time.

---

```go
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoTick failed")
} else {
    fmt.Printf("  Last tick for %s:\n", cfg.TestSymbol)
```

Checks for errors. If the request is successful — starts printing symbol data.

---

```go
fmt.Printf("    Bid:                         %.5f\n", tickData.Bid)
fmt.Printf("    Ask:                         %.5f\n", tickData.Ask)
fmt.Printf("    Last:                        %.5f\n", tickData.Last)
fmt.Printf("    Volume:                      %d\n", tickData.Volume)
```

Prints prices and volume:

* **Bid** — price at which you can sell (buyer's quote);
* **Ask** — price at which you can buy (seller's quote);
* **Last** — last trade price (relevant for stocks, futures);
* **Volume** — volume of the last tick or trade.

---

```go
tickTime := time.Unix(tickData.Time, 0).UTC()
fmt.Printf("    Time:                        %s (broker server time)\n", tickTime.Format("2006-01-02 15:04:05"))
fmt.Printf("    Time (UTC):                  %s\n", time.Unix(tickData.TimeUtc, 0).UTC().Format("2006-01-02 15:04:05"))
```

The `Time` field is stored as `int64` — Unix-style seconds counted in the **broker server's time zone**.
`time.Unix()` turns it into a `time.Time`; `.UTC()` keeps that clock as it is (without it, `Format` would shift the
value into the local time zone of the machine running the program and print a time that is neither broker time
nor UTC). `TimeUtc` is the same moment in real UTC.

> Format `2006-01-02 15:04:05` is the standard date template in Go.

---

## 📦 What the server returns

```protobuf
message MrpcMqlTick {
  int64 Time = 1;    // Unix time of last update (broker server time)
  int64 TimeUtc = 9; // Time in UTC (the same moment)
  double Bid = 2;    // Current Bid price
  double Ask = 3;    // Current Ask price
  double Last = 4;   // Last trade price
  uint64 Volume = 5; // Volume of last trade
  // Also available: TimeMsc, TimeMscUtc, Flags, VolumeReal
}
```

---

## 💡 Example output

```
Last tick for EURUSD:
    Bid:                         1.08540
    Ask:                         1.08560
    Last:                        1.08545
    Volume:                      150
    Time:                        2026-01-04 22:15:03 (broker server time)
    Time (UTC):                  2026-01-04 20:15:03
```

---

### 🧠 What it's used for

The `SymbolInfoTick()` method is used:

* to **get the latest quote** before trading actions;
* when **updating interfaces and analytical panels** (e.g., displaying current prices);
* for **logging and market monitoring**;
* in **demo examples**, to visually show current market data.

---

### 💬 In simple terms

> `SymbolInfoTick()` returns the last tick for a symbol:
> current Bid, Ask, Last prices, volume, and update time.
> It's a quick way to find out the current market state for the required instrument.
