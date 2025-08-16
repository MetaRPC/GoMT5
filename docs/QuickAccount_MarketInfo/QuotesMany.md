# Getting Multiple Quotes

> **Request:** fetch latest ticks for several symbols at once.

---

### Code Example

```go
// High-level (prints formatted quotes):
svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD"})

// Low-level (handle the payload yourself):
symbols := []string{"EURUSD", "GBPUSD"}
qs, err := svc.account.QuoteMany(ctx, symbols)
if err != nil {
    log.Printf("❌ QuoteMany error: %v", err)
    return
}
for _, q := range qs {
    if st := q.GetSymbolTick(); st != nil {
        fmt.Printf("%s | Bid=%.5f | Ask=%.5f | Time=%s\n",
            st.GetSymbol(), st.GetBid(), st.GetAsk(), st.GetTime().AsTime().Format("2006-01-02 15:04:05"))
    }
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowQuotesMany(ctx context.Context, symbols []string)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                            |
| --------- | ----------------- | -------- | ------------------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC.             |
| `symbols` | `[]string`        | yes      | List of trading symbols (e.g., `{"EURUSD","GBPUSD"}`). |

**Rules & notes**

* Empty list → likely returns an error (don’t call with `nil`/empty).
* Ensure each symbol is visible/allowed by broker to avoid empty ticks.

---

## ⬆️ Output

Returns a slice where each element holds a tick snapshot (`MrpcSubscriptionMqlTick`) for a symbol.
`ShowQuotesMany` prints selected fields from each tick.

**Tick fields (per item):**

| Field        | Type                        | Description                                         |
| ------------ | --------------------------- | --------------------------------------------------- |
| `Symbol`     | `string`                    | Symbol name (e.g., `"EURUSD"`).                     |
| `Bid`        | `double`                    | Current Bid price.                                  |
| `Ask`        | `double`                    | Current Ask price.                                  |
| `Last`       | `double`                    | Price of the last deal (Last).                      |
| `Volume`     | `uint64`                    | Volume at `Last`.                                   |
| `Time`       | `google.protobuf.Timestamp` | Last update time (server).                          |
| `TimeMsc`    | `int64`                     | Last update time in milliseconds.                   |
| `Flags`      | `uint32`                    | Tick flags (bitmask; terminal-dependent semantics). |
| `VolumeReal` | `double`                    | Precise volume at `Last`.                           |

> The exact container type is the same as in `Quote`: it exposes `GetSymbolTick()` returning `MrpcSubscriptionMqlTick`.

---

## 🎯 Purpose

* Fetch several snapshots in one call for dashboards/CLI.
* Pre-load pricing for a watchlist before validation or order checks.
* Reduce call overhead versus multiple single-symbol `Quote` calls.

---

## 🧩 Notes & Tips

* If some symbols are unknown or hidden, those items may be missing or empty — validate `GetSymbolTick() != nil`.
* For continuous updates across a basket, consider streaming (`OnSymbolTick`) instead of polling.
* Keep your list short (dozens, not hundreds) to avoid broker-side throttling or slow responses.
