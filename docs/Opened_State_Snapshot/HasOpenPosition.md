# Checking If a Symbol Has an Open Position

> **Request:** lightweight boolean check if there is an open position for a given symbol.

---

### Code Example

```go
// High-level (prints the result):
svc.ShowHasOpenPosition(ctx, "EURUSD")

// Low-level (use the boolean directly):
ok, err := svc.account.HasOpenPosition(ctx, "EURUSD")
if err != nil {
    log.Printf("❌ HasOpenPosition error: %v", err)
    return
}
if ok {
    fmt.Println("✅ There is an open position for EURUSD")
} else {
    fmt.Println("ℹ️ No open position for EURUSD")
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowHasOpenPosition(ctx context.Context, symbol string)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC. |
| `symbol`  | `string`          | yes      | Trading symbol name (e.g., `"EURUSD"`).    |

---

## ⬆️ Output

Low-level call returns:

| Field | Type    | Description                                        |
| ----- | ------- | -------------------------------------------------- |
| `ok`  | `bool`  | `true` if there is an open position, else `false`. |
| `err` | `error` | Non-nil if the RPC failed.                         |

High-level helper prints a human-readable line with the result.

---

## 🎯 Purpose

* Fast pre-check before modifying/closing a position.
* Gate logic in bots (e.g., avoid duplicate position opens).
* Lightweight health check per symbol without downloading all positions.

---

## 🧩 Notes & Tips

* Symbol must be valid and (preferably) visible; otherwise the terminal/broker may return an error.
* This is a **light** boolean; if you need details (ticket, volume, prices), call `PositionGet(ctx, symbol)` or `PositionsGet(ctx)`.
* On netting accounts, there is at most one position per symbol. On hedging setups, semantics may differ (but MT5 netting is default for many brokers).
