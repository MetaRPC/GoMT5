# 🔐 Closing an Order by Ticket (Example)

> **Request:** close an existing order/position by `ticket` (market close by default). Includes full and partial close patterns.

---

### Code Example — full close (market) ✅

```go
// High-level helper (prints a status line):
svc.ShowOrderCloseExample(ctx, 123456789)

// Internals (simplified):
// Passing price=nil, volume=nil → let server close at market with full remaining volume.
res, err := svc.account.OrderClose(ctx, 123456789, nil /*price*/, nil /*volume*/)
if err != nil {
    log.Printf("❌ OrderClose error: %v", err)
    return
}
fmt.Printf("✅ Order closed. CloseMode: %s | Code: %d (%s/%s)\n",
    res.GetCloseMode().String(),
    res.GetReturnedCode(),
    res.GetReturnedStringCode(),
    res.GetReturnedCodeDescription(),
)
```

---

### Code Example — partial close (specify volume) 🧩

```go
// Close only part of the position, e.g. 0.05 lots
partial := 0.05
res, err := svc.account.OrderClose(ctx, 123456789, nil /*market price*/, &partial)
if err != nil {
    log.Printf("❌ OrderClose(partial) error: %v", err)
    return
}
fmt.Printf("✅ Partial close sent | Volume closed: %.2f | Code: %d\n", partial, res.GetReturnedCode())
```

> If your broker/account type doesn’t allow partial closes (e.g., netting mode), the server may reject or perform a reduce‑by operation differently.

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowOrderCloseExample(ctx context.Context, ticket uint64)
```

---

## 🔽 Input

Underlying RPC `OrderClose` accepts:

| Parameter | Type              | Required | Description                                             |
| --------- | ----------------- | -------- | ------------------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Timeout/cancel control.                                 |
| `ticket`  | `uint64`          | yes      | Ticket of the order/position to close.                  |
| `price`   | `*float64`        | no       | Close price. `nil` → market close at current Bid/Ask.   |
| `volume`  | `*float64`        | no       | Volume to close in lots. `nil` → full remaining volume. |

---

## ⬆️ Output

Returns **`OrderSendData`** with fields used for diagnostics:

| Field                     | Type     | Description                                   |
| ------------------------- | -------- | --------------------------------------------- |
| `CloseMode`               | `enum`   | Server’s close mode (string via `.String()`). |
| `ReturnedCode`            | `uint32` | Numeric result code.                          |
| `ReturnedStringCode`      | `string` | Short code.                                   |
| `ReturnedCodeDescription` | `string` | Human‑readable description.                   |
| `Price`                   | `double` | Server close price (if provided).             |
| `Volume`                  | `double` | Volume processed by the close operation.      |

---

## 🎯 Purpose

* Exit an existing trade entirely or partially.
* Programmatic risk management (e.g., scale‑out on targets or risk rules).

---

## ⚠️ Notes & Tips

* **Market vs specified price:** `price=nil` is simplest (market). If you pass a price, ensure it makes sense for the side; otherwise expect rejects.
* **Slippage control:** this API variant doesn’t take a slippage parameter; if you require tight control, consider alternative close flows (e.g., reverse order or modify then close depending on broker).
* **Netting vs Hedging:** partial close semantics vary. In **netting** accounts, closing reduces net position; in **hedging**, closes affect a specific ticket.
* **Permissions & trading hours:** broker may reject closes outside session or for symbols in `TRADE_MODE_CLOSE_ONLY`.
