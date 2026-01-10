### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`SymbolInfoSessionQuote()`** method is used to retrieve **quote session times** — periods when the broker provides quotes for a symbol throughout the week.
>
> Each trading instrument in MetaTrader can have one or more such time windows, for example: morning and evening sessions, or a 24/5 round-the-clock mode.


---

## 🧩 Code example

```go
quoteSessionReq := &pb.SymbolInfoSessionQuoteRequest{
    Symbol:       cfg.TestSymbol,
    DayOfWeek:    pb.DayOfWeek_MONDAY,
    SessionIndex: 0,
}
quoteSessionData, err := account.SymbolInfoSessionQuote(ctx, quoteSessionReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoSessionQuote failed")
} else {
    fmt.Printf("  Monday quote session #0:\n")
    if quoteSessionData.From != nil {
        fromTime := quoteSessionData.From.AsTime()
        fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
        fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
    }
    if quoteSessionData.To != nil {
        toTime := quoteSessionData.To.AsTime()
        toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()
        fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
quoteSessionReq := &pb.SymbolInfoSessionQuoteRequest{
    Symbol:       cfg.TestSymbol,
    DayOfWeek:    pb.DayOfWeek_MONDAY,
    SessionIndex: 0,
}
```

A request is created with parameters:

* **`Symbol`** — trading instrument name (e.g., `EURUSD`);
* **`DayOfWeek`** — day of the week for which the session is requested (`MONDAY`, `TUESDAY`, etc.);
* **`SessionIndex`** — session index for that day (if there are multiple, e.g., morning and evening).

---

```go
quoteSessionData, err := account.SymbolInfoSessionQuote(ctx, quoteSessionReq)
```

A call is made to the MetaTrader server via gRPC. The response contains a structure with `From` and `To` fields — session start and end times.

---

```go
if quoteSessionData.From != nil {
    fromTime := quoteSessionData.From.AsTime()
    fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
    fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
}
```

If the `From` field is filled, it's converted to `time.Time`, then translated to seconds from the start of the day.

---

```go
if quoteSessionData.To != nil {
    toTime := quoteSessionData.To.AsTime()
    toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()
    fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
}
```

Similarly, the session end time is calculated.

---

## 📦 What the Server Returns

```protobuf
message SymbolInfoSessionQuoteData {
  google.protobuf.Timestamp From = 1; // Quote session start time
  google.protobuf.Timestamp To = 2;   // Quote session end time
}
```

---

## 💡 Example Output

```
Monday quote session #0:
    From (seconds from day start): 0
    To (seconds from day start):   86399
```

📘 This means quotes are available all day — from 00:00:00 to 23:59:59.

---

### ℹ️ Why Values Can Be Zero

If the terminal displays:

```
From (seconds from day start): 0
To (seconds from day start):   0
```

this is **not a code error**, but one of the standard MetaTrader API behaviors.

Reasons:

1. **Instrument trades 24/5 non-stop** — the server doesn't set separate intervals.
2. **Broker hasn't configured sessions** — demo servers often return zeros.
3. **Gateway doesn't support these fields** — API returns empty `Timestamp`.
4. **Invalid `SessionIndex` or `DayOfWeek`** — there may be no sessions on that day.

✅ In these cases, `0` is interpreted as "active all day".

---

### 🧠 What It's Used For

The `SymbolInfoSessionQuote()` method is used:

* to **determine active market hours** for an instrument;
* for **restricting trading outside working hours**;
* when building **analytical panels** and trading calendars;
* in **demo examples** where it's important to show how to get quote schedules via API.

---

### 💬 In Simple Terms

> `SymbolInfoSessionQuote()` shows **what hours the symbol receives quotes**.
> Even if it returns 0, it means quotes are available all day (24/5).
