# Listing All Symbols

> **Request:** fetch the full list of symbols available to the connected account/broker.

---

### Code Example

```go
// High-level (prints one per line):
svc.ShowAllSymbols(ctx)

// Low-level (work with the slice yourself):
names, err := svc.account.ShowAllSymbols(ctx)
if err != nil {
    log.Printf("❌ ShowAllSymbols error: %v", err)
    return
}
for _, name := range names {
    fmt.Println(name)
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowAllSymbols(ctx context.Context)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                               |
| --------- | ----------------- | -------- | ----------------------------------------- |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC |

---

## ⬆️ Output

Returns a slice of symbol names and, in the high-level helper, prints them to stdout.

| Item       | Type     | Description                                 |
| ---------- | -------- | ------------------------------------------- |
| `names[i]` | `string` | Symbol name (e.g., `"EURUSD"`, `"XAUUSD"`). |

---

## 🎯 Purpose

* Discover what instruments are tradable on the current server/account.
* Feed watchlists, dropdowns, and auto-complete in UIs/CLIs.
* Validate a symbol before attempting `EnsureSymbolVisible`, `Quote`, or order placement.

---

## 🧩 Notes & Tips

* The list may be large (hundreds/thousands). Consider filtering (prefixes, groups) on your side.
* Not every symbol from this list is guaranteed to be *visible* in Market Watch. Use `EnsureSymbolVisible(ctx, symbol)` before requesting ticks/orders.
* Names are broker-defined; suffixes like `.pro`, `.mini`, `.ecn` are common. Use your own suffix probing when selecting a default symbol.
* If you iterate this list and call per-symbol RPCs, add short per-call timeouts to avoid long runs.
