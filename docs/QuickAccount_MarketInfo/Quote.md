# Getting a Quote

> **Request:** latest tick for a symbol (Bid/Ask/Last & timestamp).
> Fetch a single quote snapshot for one symbol.

---

### Code Example

```go
// High-level (prints formatted quote):
svc.ShowQuote(ctx, "EURUSD")

// Low-level (for custom handling):
q, err := svc.account.Quote(ctx, "EURUSD")
if err != nil {
    log.Printf("❌ Quote error: %v", err)
    return
}
if st := q.GetSymbolTick(); st != nil {
    fmt.Printf("%s | Bid=%.5f | Ask=%.5f | Time=%s\n",
        st.GetSymbol(), st.GetBid(), st.GetAsk(), st.GetTime().AsTime().Format("2006-01-02 15:04:05"))
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowQuote(ctx context.Context, symbol string)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC. |
| `symbol`  | `string`          | yes      | Trading symbol name (e.g., `"EURUSD"`).    |

---

## ⬆️ Output

`ShowQuote` prints selected fields from **`MrpcSubscriptionMqlTick`** (tick snapshot):

| Field        | Type                        | Description                                                      |
| ------------ | --------------------------- | ---------------------------------------------------------------- |
| `Time`       | `google.protobuf.Timestamp` | Time of the last prices update (server time).                    |
| `Bid`        | `double`                    | Current Bid price.                                               |
| `Ask`        | `double`                    | Current Ask price.                                               |
| `Last`       | `double`                    | Price of the last deal (Last).                                   |
| `Volume`     | `uint64`                    | Volume for the current Last price.                               |
| `TimeMsc`    | `int64`                     | Time of the last update in milliseconds.                         |
| `Flags`      | `uint32`                    | Tick flags (bitmask; specific semantics are terminal-dependent). |
| `VolumeReal` | `double`                    | More precise volume for the current Last price.                  |
| `Symbol`     | `string`                    | Symbol name (e.g., `"EURUSD"`).                                  |

> Structure name and fields come from generated Go types (`mt5-term-api-subscriptions.pb.go`: `MrpcSubscriptionMqlTick`).

---

## 🎯 Purpose

* Grab a *current* snapshot for UI/CLI display.
* Validate symbol availability & connectivity.
* Pre-check pricing before `OrderCheck` / `OrderSend`.

---

## 🧩 Notes & Tips

* Ensure the symbol is visible/selected before requesting quotes; otherwise the terminal may return empty data or error.
* The timestamp is provided both as `Time` (seconds) and `TimeMsc` (milliseconds). Prefer `Time` and format via `AsTime()`.
* For continuous updates, use streaming (`OnSymbolTick`) instead of repeated single `Quote` calls.
