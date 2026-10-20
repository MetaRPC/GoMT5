# ✅ Get Symbol Quote Session Times

> **Request:** get start and end times of quote session for symbol on specified day of week.

**API Information:**

* **Low-level API:** `MT5Account.SymbolInfoSessionQuote(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `SymbolInfoSessionQuote` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `SymbolInfoSessionQuote(SymbolInfoSessionQuoteRequest) → SymbolInfoSessionQuoteReply`
* **Low‑level client (generated):** `MarketInfoClient.SymbolInfoSessionQuote(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Returns when quotes are available for a symbol on specific days.
* **Why you need it.** Check market hours, validate trading times, schedule operations.
* **Time format.** Returns `From`/`To` timestamps: their time of day (broker server time) is the session boundary, the date part is to be ignored.

---

## 🎯 Purpose

Use it to:

* Check when quotes are available for symbols
* Validate if market is open before requesting prices
* Schedule quote monitoring operations
* Build trading session calendars
* Determine quote availability windows

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolInfoSessionQuote retrieves quote session times for a symbol.
// Returns From and To timestamps whose time of day is the session start/end (date part ignored).
func (a *MT5Account) SymbolInfoSessionQuote(
    ctx context.Context,
    req *pb.SymbolInfoSessionQuoteRequest,
) (*pb.SymbolInfoSessionQuoteData, error)
```

**Request message:**

```protobuf
SymbolInfoSessionQuoteRequest {
  string Symbol = 1;       // Trading symbol
  DayOfWeek day_of_week = 2;  // Day of week (0=Sunday to 6=Saturday)
  uint32 session_index = 3;   // Session index (usually 0)
}
```

---

## 🔽 Input

| Parameter | Type                                   | Description                                   |
| --------- | -------------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                      | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolInfoSessionQuoteRequest`    | Request with Symbol, DayOfWeek, SessionIndex  |

**Request fields:**

| Field          | Type     | Description                                   |
| -------------- | -------- | --------------------------------------------- |
| `Symbol`       | `string` | Trading symbol (e.g., "EURUSD")               |
| `DayOfWeek`    | `pb.DayOfWeek` | Day of week (see enum below)            |
| `SessionIndex` | `uint32` | Session index (usually 0 for main session)    |

**DayOfWeek enum values:**

| Value | Protobuf Enum | Description |
|-------|---------------|-------------|
| `0` | `pb.DayOfWeek_SUNDAY` | Sunday |
| `1` | `pb.DayOfWeek_MONDAY` | Monday |
| `2` | `pb.DayOfWeek_TUESDAY` | Tuesday |
| `3` | `pb.DayOfWeek_WEDNESDAY` | Wednesday |
| `4` | `pb.DayOfWeek_THURSDAY` | Thursday |
| `5` | `pb.DayOfWeek_FRIDAY` | Friday |
| `6` | `pb.DayOfWeek_SATURDAY` | Saturday |

**Usage:** Can use either numeric value or protobuf enum constant

---

## ⬆️ Output — `SymbolInfoSessionQuoteData`

| Field  | Type                        | Go Type                    | Description                                    |
| ------ | --------------------------- | -------------------------- | ---------------------------------------------- |
| `From` | `google.protobuf.Timestamp` | `*timestamppb.Timestamp`   | Session start time (use `.AsTime()` to convert)|
| `To`   | `google.protobuf.Timestamp` | `*timestamppb.Timestamp`   | Session end time (use `.AsTime()` to convert)  |

**Note:** To get seconds from day start, convert using:
```go
fromTime := data.From.AsTime()
fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
```

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolInfoSessionQuote - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoSessionQuote_HOW.md)**

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `5s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Time format:** `From`/`To` are timestamps; only their time of day (broker server time) matters. Convert to seconds from midnight as shown above (00:00:00 = 0, 23:59:59 = 86399).
* **Multiple sessions:** Some symbols have multiple quote sessions per day (SessionIndex 0, 1, etc).
* **Weekend:** Most symbols have no quote sessions on Saturday and Sunday.

---

## 🔗 Usage Examples

### 1) Get Monday quote session

```go
package main

import (
    "github.com/google/uuid"
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
)

func main() {
    account, err := mt5.NewMT5Account(591129415, "password", "mt5.mrpc.pro:443", uuid.Nil)
    if err != nil {
        panic(err)
    }
    defer account.Close()

    // Connect first
    _, err = account.ConnectEx(context.Background(), &pb.ConnectExRequest{User: 591129415, Password: "password", MtClusterName: "MetaQuotes-Demo"})
    if err != nil {
        panic(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    data, err := account.SymbolInfoSessionQuote(ctx, &pb.SymbolInfoSessionQuoteRequest{
        Symbol:       "EURUSD",
        DayOfWeek:    pb.DayOfWeek_MONDAY,
        SessionIndex: 0,
    })
    if err != nil {
        panic(err)
    }

    fromTime := data.From.AsTime()
    toTime := data.To.AsTime()

    fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
    toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()

    fmt.Printf("EURUSD Monday quote session: %02d:%02d - %02d:%02d\n",
        fromSeconds/3600, (fromSeconds%3600)/60, toSeconds/3600, (toSeconds%3600)/60)
}
```

### 2) Get week schedule

```go
func GetWeeklyQuoteSchedule(account *mt5.MT5Account, symbol string) {
    ctx := context.Background()

    days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

    fmt.Printf("Quote sessions for %s:\n", symbol)
    for day := 0; day <= 6; day++ {
        data, err := account.SymbolInfoSessionQuote(ctx, &pb.SymbolInfoSessionQuoteRequest{
            Symbol:       symbol,
            DayOfWeek:    pb.DayOfWeek(day),
            SessionIndex: 0,
        })
        if err != nil {
            fmt.Printf("  %s: Not available\n", days[day])
            continue
        }

        fromStr := formatTime(data.From)
        toStr := formatTime(data.To)
        fmt.Printf("  %s: %s - %s\n", days[day], fromStr, toStr)
    }
}

func formatTime(ts *timestamppb.Timestamp) string {
    t := ts.AsTime()
    seconds := t.Hour()*3600 + t.Minute()*60 + t.Second()
    hours := seconds / 3600
    mins := (seconds % 3600) / 60
    return fmt.Sprintf("%02d:%02d", hours, mins)
}
```

### 3) Check if quotes available now

```go
func AreQuotesAvailableNow(account *mt5.MT5Account, symbol string) (bool, error) {
    ctx := context.Background()

    // Session times are broker server time: compare with the server clock, not the local one
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        return false, err
    }
    now := summary.ServerTime.AsTime() // broker wall clock
    dayOfWeek := pb.DayOfWeek(now.Weekday())
    currentSeconds := int64(now.Hour()*3600 + now.Minute()*60 + now.Second())

    data, err := account.SymbolInfoSessionQuote(ctx, &pb.SymbolInfoSessionQuoteRequest{
        Symbol:       symbol,
        DayOfWeek:    dayOfWeek,
        SessionIndex: 0,
    })
    if err != nil {
        return false, err
    }

    fromTime := data.From.AsTime()
    toTime := data.To.AsTime()
    fromSeconds := int64(fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second())
    toSeconds := int64(toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second())

    available := currentSeconds >= fromSeconds && currentSeconds <= toSeconds

    fmt.Printf("%s quotes available now: %v\n", symbol, available)
    return available, nil
}
```

---

## 📚 See Also

* [SymbolInfoSessionTrade](./SymbolInfoSessionTrade.md) - Get trading session times
* [SymbolIsSynchronized](../2.%20Symbol_information/SymbolIsSynchronized.md) - Check quote sync
