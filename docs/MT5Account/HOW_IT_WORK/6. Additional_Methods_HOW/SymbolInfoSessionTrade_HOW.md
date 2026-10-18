### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`SymbolInfoSessionTrade()`** method is used to retrieve **trading session times** — periods when trading operations are allowed for a symbol.
>
> Unlike `SymbolInfoSessionQuote()` (which shows when quotes are available), this method shows **when trading is allowed** — that is, when the server accepts orders and executes trades.


---

## 🧩 Code example

```go
tradeSessionReq := &pb.SymbolInfoSessionTradeRequest{
    Symbol:       cfg.TestSymbol,
    DayOfWeek:    pb.DayOfWeek_MONDAY,
    SessionIndex: 0,
}
tradeSessionData, err := account.SymbolInfoSessionTrade(ctx, tradeSessionReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoSessionTrade failed")
} else {
    fmt.Printf("  Monday trade session #0:\n")
    // From and To are *timestamppb.Timestamp holding a time of day (broker server time); their date is not meaningful.
    // Count the seconds from the start of From's day: a session that runs to midnight ends at 24:00 = 86400
    // (To is then the next day's 00:00, which Hour()/Minute() alone would show as 0).
    if tradeSessionData.From != nil && tradeSessionData.To != nil {
        dayStart := tradeSessionData.From.AsTime().Truncate(24 * time.Hour)
        fromSeconds := int64(tradeSessionData.From.AsTime().Sub(dayStart).Seconds())
        toSeconds := int64(tradeSessionData.To.AsTime().Sub(dayStart).Seconds())
        fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
        fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
    }
}
```


### 🟢 Detailed Code Explanation

```go
tradeSessionReq := &pb.SymbolInfoSessionTradeRequest{
    Symbol:       cfg.TestSymbol,
    DayOfWeek:    pb.DayOfWeek_MONDAY,
    SessionIndex: 0,
}
```

A request is created with parameters:

* **`Symbol`** — trading instrument (e.g., `EURUSD`);
* **`DayOfWeek`** — day of the week for which the schedule is needed;
* **`SessionIndex`** — sequential session number for that day.

---

```go
tradeSessionData, err := account.SymbolInfoSessionTrade(ctx, tradeSessionReq)
```

A request is sent to the MetaTrader server via gRPC. The response returns a structure with `From` and `To` fields — trading session start and end.

---

```go
if tradeSessionData.From != nil && tradeSessionData.To != nil {
    dayStart := tradeSessionData.From.AsTime().Truncate(24 * time.Hour)
    fromSeconds := int64(tradeSessionData.From.AsTime().Sub(dayStart).Seconds())
    toSeconds := int64(tradeSessionData.To.AsTime().Sub(dayStart).Seconds())
    fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
    fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
}
```

`From` and `To` carry a time of day in the broker server's time zone (seconds from 00:00); their date part has no
meaning. Both are converted to `time.Time` and counted from the start of `From`'s day. A session that runs to
midnight ends at **24:00**: `To` is then the next day's 00:00, so it gives 86400. Taking only
`Hour()`/`Minute()`/`Second()` of `To` would print 0 for that end.

---


### 📘 What is `SessionIndex`

The `SessionIndex` parameter **is not related to the day of the week**, it indicates the **session number within the selected day**.

For example, if a symbol has several trading intervals during the day:

| DayOfWeek | SessionIndex | From  | To    |
| --------- | ------------ | ----- | ----- |
| MONDAY    | 0            | 09:00 | 12:00 |
| MONDAY    | 1            | 13:00 | 17:00 |
| TUESDAY   | 0            | 09:00 | 17:00 |

That is:

* `DayOfWeek` — day of the week;
* `SessionIndex` — sequential session number on that day.

If an instrument has only one trading session per day (e.g., Forex 24 hours), then `SessionIndex = 0`.

---

## 📦 What the Server Returns

```protobuf
message SymbolInfoSessionTradeData {
  google.protobuf.Timestamp From = 1; // Trading session start time
  google.protobuf.Timestamp To = 2;   // Trading session end time
}
```

---

## 💡 Example Output

```
Monday trade session #0:
    From (seconds from day start): 0
    To (seconds from day start):   86400
```

📘 This means trading is available all day — from 00:00 to 24:00 (EURUSD on MetaQuotes-Demo).

---

### ℹ️ When There Is No Session

A day (or a `SessionIndex`) without a session does not answer with zeros: the call returns an error,
`ERR_MARKET_SESSION_INDEX`. To read every session of a day, ask for `SessionIndex` 0, 1, 2, … until that error.

---

### 🧠 What It's Used For

The `SymbolInfoSessionTrade()` method is used:

* to **check trading hours** before opening orders;
* to **avoid trading outside sessions**;
* for **visualizing market operation** (charts, calendars);
* in **backtests and simulations** to account for trading schedules.

---

### 💬 In Simple Terms

> `SymbolInfoSessionTrade()` shows **what hours and days trading is allowed** for a symbol.
> `From` 0 and `To` 86400 mean trading all day; a day without a session answers with `ERR_MARKET_SESSION_INDEX`.
