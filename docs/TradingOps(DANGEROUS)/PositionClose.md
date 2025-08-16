# 🔒 Closing an Open Position by Symbol

> **Request:** find an open position by `symbol` and close it at market.

---

### Code Example

```go
// High-level helper (prints a status line):
svc.ShowPositionClose(ctx, "EURUSD")

// Internals (simplified):
p, err := svc.account.PositionGet(ctx, "EURUSD")
if err != nil {
    log.Printf("❌ PositionGet error: %v", err)
    return
}
if p == nil || p.GetTicket() == 0 {
    fmt.Printf("⚠️ No position found for symbol %s\n", "EURUSD")
    return
}

// Close the whole position at market
_, err = svc.account.PositionClose(ctx, p)
if err != nil {
    log.Printf("❌ PositionClose error: %v", err)
    return
}
fmt.Printf("✅ Position closed: Ticket %d (%s)\n", p.GetTicket(), p.GetSymbol())
```

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowPositionClose(ctx context.Context, symbol string)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                             |
| --------- | ----------------- | -------- | ------------------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Timeout/cancel control.                                 |
| `symbol`  | `string`          | yes      | Symbol to look up the open position (e.g., `"EURUSD"`). |

---

## ⬆️ Output

The helper prints one of the following messages:

* `✅ Position closed: Ticket <id> (<symbol>)`
* `⚠️ No position found for symbol <symbol>`
* `❌ PositionGet/PositionClose error: <err>`

> Under the hood, `PositionClose` returns an implementation-specific result; success is signaled by `err == nil`.

---

## 🎯 Purpose

* Programmatically exit an **entire** open position for a given symbol at market.
* Handy for panic‑close buttons, risk rules, or end‑of‑session cleanup.

---

## 🧩 Notes & Tips

* This helper closes **full volume** for the symbol’s current position.

  * For **partial close**, use `OrderClose(ticket, price, &volume)` (see `ShowOrderCloseExample`) or broker‑specific reduce flows.
* **Netting vs Hedging:**

  * In **netting** mode you have at most one position per symbol → this helper fits perfectly.
  * In **hedging** mode there can be multiple tickets per symbol; this helper closes the one returned by `PositionGet` (typically the aggregate or most recent depending on API semantics). If you need per‑ticket control, prefer `OrderClose` by `ticket`.
* **Trading hours & permissions:** close can be rejected outside session or when symbol is in `CLOSE_ONLY` mode.
* Combine with `ShowHasOpenPosition` to avoid attempting a close when nothing is open.
