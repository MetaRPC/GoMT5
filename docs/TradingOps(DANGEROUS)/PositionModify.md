# 🛠️ Modifying an Open Position (SL/TP)

> **Request:** update **Stop Loss** and/or **Take Profit** for an existing **open position** by its `ticket`.

---

### Code Example

```go
// High-level helper (prints result):
sl := 1.07500
tp := 1.08500
svc.ShowPositionModify(ctx, 987654321 /*position ticket*/, &sl, &tp)

// Internals (simplified):
ok, err := svc.account.PositionModify(ctx, 987654321, &sl, &tp)
if err != nil {
    log.Printf("❌ PositionModify error: %v", err)
    return
}
if ok {
    fmt.Printf("✅ Position %d modified (SL/TP updated)\n", 987654321)
} else {
    fmt.Printf("⚠️ Position %d was NOT modified\n", 987654321)
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) ShowPositionModify(ctx context.Context, ticket uint64, newSL, newTP *float64)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                                  |
| --------- | ----------------- | -------- | ------------------------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Timeout/cancel control.                                      |
| `ticket`  | `uint64`          | yes      | **Position** ticket to modify.                               |
| `newSL`   | `*float64`        | no       | New **Stop Loss** (absolute price). `nil` → do not change.   |
| `newTP`   | `*float64`        | no       | New **Take Profit** (absolute price). `nil` → do not change. |

> Both `newSL` and `newTP` are **absolute prices**, not offsets.

---

## ⬆️ Output

The helper prints a message. Underlying call returns:

| Field | Type    | Description                              |
| ----: | ------- | ---------------------------------------- |
|  `ok` | `bool`  | `true` if the server applied the change. |
| `err` | `error` | Error if the request failed.             |

---

## 🎯 Purpose

* Move protective levels after entry (e.g., trail SL, set TP once structure forms).
* Enforce risk rules programmatically.

---

## 🧩 Notes & Tips

* Validate distances: brokers enforce **min stops** / freeze levels — too tight SL/TP will be rejected.
* Direction sanity:

  * **Long**: `SL < market`, `TP > market`.
  * **Short**: `SL > market`, `TP < market`.
* If you only need to change one level, pass the other as `nil`.
* To modify a **pending order’s price/expiration**, use `OrderModify` instead — this method is only for **open positions**.
* Combine with `SymbolParams` (`Digits`, min distances) and `Quote` to compute valid absolute prices. 📏
