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
    // From and To are *timestamppb.Timestamp holding a time of day (broker server time); their date is not meaningful.
    // Count the seconds from the start of From's day: a session that runs to midnight ends at 24:00 = 86400
    // (To is then the next day's 00:00, which Hour()/Minute() alone would show as 0).
    if quoteSessionData.From != nil && quoteSessionData.To != nil {
        dayStart := quoteSessionData.From.AsTime().Truncate(24 * time.Hour)
        fromSeconds := int64(quoteSessionData.From.AsTime().Sub(dayStart).Seconds())
        toSeconds := int64(quoteSessionData.To.AsTime().Sub(dayStart).Seconds())
        fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
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
if quoteSessionData.From != nil && quoteSessionData.To != nil {
    dayStart := quoteSessionData.From.AsTime().Truncate(24 * time.Hour)
    fromSeconds := int64(quoteSessionData.From.AsTime().Sub(dayStart).Seconds())
    toSeconds := int64(quoteSessionData.To.AsTime().Sub(dayStart).Seconds())
    fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
    fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
}
```

`From` and `To` carry a time of day in the broker server's time zone (seconds from 00:00); their date part has no
meaning. Both are converted to `time.Time` and counted from the start of `From`'s day. A session that runs to
midnight ends at **24:00**: `To` is then the next day's 00:00 (MetaQuotes-Demo EURUSD: `From` = 1970-01-01 00:00,
`To` = 1970-01-02 00:00), so `To` gives 86400. Taking only `Hour()`/`Minute()`/`Second()` of `To` would print 0 for
that end.

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
    To (seconds from day start):   86400
```

📘 This means quotes are available all day — from 00:00 to 24:00 (EURUSD on MetaQuotes-Demo).

---

### ℹ️ When There Is No Session

A day (or a `SessionIndex`) without a session does not answer with zeros: the call returns an error,
`ERR_MARKET_SESSION_INDEX` (EURUSD on MetaQuotes-Demo: `SUNDAY` #0 or `MONDAY` #1). To read every session of a day,
ask for `SessionIndex` 0, 1, 2, … until that error.

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
> `From` 0 and `To` 86400 mean quotes all day; a day without a session answers with `ERR_MARKET_SESSION_INDEX`.
