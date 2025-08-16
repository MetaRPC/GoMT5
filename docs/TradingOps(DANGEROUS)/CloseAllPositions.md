# 🧹 Closing All Open Positions

> **Request:** close **all** open positions for the account (market). Use with caution.

---

### Code Example

```go
// High-level helper (prints a status line):
svc.ShowCloseAllPositions(ctx)

// Internals (simplified):
if err := svc.account.CloseAllPositions(ctx); err != nil {
    log.Printf("❌ CloseAllPositions error: %v", err)
    return
}
fmt.Println("✅ All positions closed (or none existed).")
```

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowCloseAllPositions(ctx context.Context)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description             |
| --------- | ----------------- | -------- | ----------------------- |
| `ctx`     | `context.Context` | yes      | Timeout/cancel control. |

---

## ⬆️ Output

* On success: `✅ All positions closed (or none existed).`
* On error: `❌ CloseAllPositions error: <err>`

> Under the hood, the implementation may iterate over current positions and issue a close for each. Behavior (e.g., batching) can vary by broker/API.

---

## 🎯 Purpose

* Emergency/Panic close (risk off).
* End‑of‑day cleanup.
* Reset state before strategy change or redeploy.

---

## ⚠️ Notes & Tips

* **Danger zone:** this command will attempt to close *every* open position. Double‑confirm in UI or logs before calling from automation.
* **Permissions & trading hours:** broker may reject closes outside session or for symbols in `CLOSE_ONLY` mode.
* **Slippage control:** mass‑closing under volatility can result in unexpected fills; consider pre‑hedging or staged exits if precision is critical.
* Works best in **netting** accounts. In **hedging**, expect ticket‑by‑ticket closes.
