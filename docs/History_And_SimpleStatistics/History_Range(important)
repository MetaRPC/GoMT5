# 🕒 History Range (`from`/`to`) — How to Use

> **Purpose:** correctly build and pass `time.Time` ranges to history/statistics methods (orders & deals). Includes ready‑to‑copy helpers, timezone tips, and gotchas.

---

## TL;DR

```go
// Last 7 days
from := time.Now().AddDate(0, 0, -7)
to   := time.Now()

// Use with history APIs:
svc.ShowOrdersHistory(ctx)           // (helper internally builds a range)
svc.ShowDealsCount(ctx, from, to, "")
```

* `from` — start of interval
* `to`   — end of interval
* Prefer **UTC** or a **fixed timezone** for consistency; convert explicitly when needed.

---

## Where ranges are required ✅

These high‑level helpers and/or underlying RPCs need a time range:

* `ShowOrdersHistory(ctx)` → uses a range internally (e.g., last 7 days)
* `ShowDealsCount(ctx, from, to, symbol)` → **you pass** `from`/`to`

> Methods by ticket (`ShowOrderByTicket`, `ShowDealByTicket`) do **not** need a range.

---

## Patterns & Recipes 🍳

### 1) Last N days

```go
func LastNDays(n int) (time.Time, time.Time) {
    to := time.Now().UTC()
    from := to.AddDate(0, 0, -n)
    return from, to
}
from, to := LastNDays(7)
```

### 2) Current day (local timezone)

```go
loc := time.Local // or time.FixedZone(...)
now := time.Now().In(loc)
start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
end   := start.Add(24*time.Hour - time.Nanosecond)
from, to := start, end
```

### 3) This week (Mon..Sun) in a specific TZ

```go
loc, _ := time.LoadLocation("Europe/London")
now := time.Now().In(loc)
offset := (int(now.Weekday()) + 6) % 7 // Monday=0
monday := time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, loc).AddDate(0,0,-offset)
from := monday
to   := monday.AddDate(0,0,7).Add(-time.Nanosecond)
```

### 4) Month to date (UTC)

```go
now := time.Now().UTC()
start := time.Date(now.Year(), now.Month(), 1, 0,0,0,0, time.UTC)
from, to := start, now
```

---

## Passing to methods 🧩

### `ShowDealsCount`

```go
from := time.Now().AddDate(0,0,-7)
to   := time.Now()
// all symbols
svc.ShowDealsCount(ctx, from, to, "")
// single symbol
svc.ShowDealsCount(ctx, from, to, "EURUSD")
```

### `ShowOrdersHistory`

The helper already sets a default range (last 7 days). To customize, make your own wrapper that calls the underlying RPC with your `from`/`to`.

```go
func ShowOrdersHistoryRange(ctx context.Context, svc *MT5Service, from, to time.Time) {
    sortMode := pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_CLOSE_TIME_DESC
    data, err := svc.account.OrdersHistory(ctx, sortMode, &from, &to, nil, nil)
    if err != nil {
        log.Printf("❌ OrdersHistory error: %v", err)
        return
    }
    for _, item := range data.GetHistoryData() {
        o := item.GetHistoryOrder()
        if o == nil { continue }
        fmt.Printf("[%s] #%d %s vol=%.2f open=%.5f last=%.5f closed=%s\n",
            o.GetType().String(), o.GetTicket(), o.GetSymbol(), o.GetVolumeInitial(),
            o.GetPriceOpen(), o.GetPriceCurrent(), o.GetDoneTime().AsTime().Format("2006-01-02 15:04:05"))
    }
}
```

---

## Timezone strategy 🌍

* **Prefer UTC** for program logic: `time.Now().UTC()` — avoids DST surprises.
* Convert for display: `t.In(userTZ)`.
* If your broker/server follows a known TZ (e.g., EET for many FX servers), build ranges in that TZ for reporting cutoff times.

```go
loc, _ := time.LoadLocation("Europe/Athens") // example server TZ
from := time.Date(2025, 8, 1, 0,0,0,0, loc)
to   := time.Date(2025, 8, 31,23,59,59, 0, loc)
```

---

## Inclusivity, edges & gaps ⚠️

* Some backends treat `to` as an exclusive bound. To be safe, bump `to` slightly:

  ```go
  toSafe := to.Add(1 * time.Second)
  ```
* Ensure `from.Before(to)`; if not, swap or adjust.
* Very large ranges → big payloads / pagination. Prefer smaller windows (7–30 days) for interactive tooling.
* Clock skew: if machines have wrong clocks, you may miss the newest items — sync time (NTP).

---

## Validation helpers ✅

```go
func ClampRange(from, to time.Time) (time.Time, time.Time) {
    if to.Before(from) { from, to = to, from }
    return from, to
}

func EndOfDay(t time.Time) time.Time {
    loc := t.Location()
    return time.Date(t.Year(), t.Month(), t.Day(), 23,59,59, 999_999_999, loc)
}
```

---

## FAQ ❓

**Q:** Do I need milliseconds?
**A:** Seconds are typically enough; keep `time.Time` precision, the PB layer handles it.

**Q:** Should I always use UTC?
**A:** Use UTC for logic; convert for presentation/reporting. If server accounting days close in a specific TZ, build ranges in that TZ.

**Q:** Why is my 7‑day query empty?
**A:** Check account/activity, symbol filter, trading days, and ensure `to` isn’t in the past (e.g., stale clock).

---

## Copy‑paste snippets 📋

```go
// Last 7 days UTC
from := time.Now().UTC().AddDate(0,0,-7)
to   := time.Now().UTC()

// Today in local TZ
now := time.Now()
from = time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, now.Location())
to   = EndOfDay(now)

// Month‑to‑date UTC
now = time.Now().UTC()
from = time.Date(now.Year(), now.Month(), 1, 0,0,0,0, time.UTC)
to   = now
```
