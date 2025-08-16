# Counting Deals in a Range

> **Request:** get the number of executed deals within a specific time range.

---

### Code Example

```go
from := time.Now().AddDate(0, 0, -7) // 7 days ago
to   := time.Now()                   // now

svc.ShowDealsCount(ctx, from, to, "")
```

---

### Method Signature

```go
func (s *MT5Service) ShowDealsCount(ctx context.Context, from, to time.Time, symbol string)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                                |
| --------- | ----------------- | -------- | ---------------------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Timeout / cancellation control.                            |
| `from`    | `time.Time`       | yes      | Start of time interval. See: **History Range Note**.       |
| `to`      | `time.Time`       | yes      | End of time interval. See: **History Range Note**.         |
| `symbol`  | `string`          | no       | Symbol filter (e.g. `"EURUSD"`), empty string → all deals. |

---

## ⬆️ Output

Prints a single integer count of all deals found in the given range (optionally filtered by symbol).

| Field   | Type  | Description               |
| ------- | ----- | ------------------------- |
| `Count` | `int` | Number of executed deals. |

---

## 🎯 Purpose

* Quick **activity metric**: how many trades were executed in the selected period.
* Useful for **statistics**, **reports**, and **trading frequency analysis**.
* When combined with `ShowOrdersHistory`, gives both **orders placed** and **deals executed**.

---

## 🧩 Notes & Tips

* Always ensure `from < to`, otherwise you’ll get empty results or errors.
* Use empty `symbol` when you want total activity across all instruments.
* Works well for periodic reporting (daily, weekly, monthly trade counts).
* 🔗 For `from/to` handling see: **History_Range(important)**.
