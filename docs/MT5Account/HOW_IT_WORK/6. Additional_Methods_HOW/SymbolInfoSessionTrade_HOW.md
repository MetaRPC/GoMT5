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
    if tradeSessionData.From != nil {
        fromTime := tradeSessionData.From.AsTime()
        fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
        fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
    }
    if tradeSessionData.To != nil {
        toTime := tradeSessionData.To.AsTime()
        toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()
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
if tradeSessionData.From != nil {
    fromTime := tradeSessionData.From.AsTime()
    fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
    fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
}
```

Converts `From` to `time.Time` format and calculates the number of seconds from the start of the day.

---

```go
if tradeSessionData.To != nil {
    toTime := tradeSessionData.To.AsTime()
    toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()
    fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
}
```

Same for the session end.

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
    To (seconds from day start):   86399
```

📘 This means trading is available all day — from 00:00:00 to 23:59:59.

---

### ℹ️ Why Values Can Be Zero

If you see the output:

```
From (seconds from day start): 0
To (seconds from day start):   0
```

this is **normal behavior**, not a code error.

Reasons:

1. **24/5 trading** — the broker doesn't set separate intervals.
2. **Broker hasn't filled in the schedule** — often found on demo servers.
3. **Gateway doesn't transmit these fields** — API returns empty Timestamp.
4. **Day or index contains no data** — there's no session on that day with that index.

✅ In such cases, `0` is interpreted as "all day open for trading".

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
> If it returns zeros — it means trading is available all day or the broker hasn't specified a schedule.
