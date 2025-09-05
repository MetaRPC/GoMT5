# Listing Open Positions

> **Request:** fetch snapshot of all currently open positions (netted per symbol on MT5).

---

### Code Example

```go
// High-level (prints selected fields per position):
svc.ShowPositions(ctx)

// Low-level (full access to the slice):
infos, err := svc.account.PositionsGet(ctx)
if err != nil {
    log.Printf("❌ PositionsGet error: %v", err)
    return
}
if len(infos) == 0 {
    fmt.Println("📭 No open positions.")
    return
}
for _, p := range infos {
    fmt.Printf("Pos #%d | %s | Volume=%.2f | Open=%.5f | Profit=%.2f\n",
        p.GetTicket(), p.GetSymbol(), p.GetVolume(), p.GetPriceOpen(), p.GetProfit())
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowPositions(ctx context.Context)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC. |

---

## ⬆️ Output

Returns a slice of **`PositionInfo`** objects.
`ShowPositions` prints a subset of fields per item.

### `PositionInfo` (key fields)

| Field          | Type                        | Description                             |
| -------------- | --------------------------- | --------------------------------------- |
| `Ticket`       | `uint64`                    | Position ticket ID.                     |
| `Symbol`       | `string`                    | Symbol name.                            |
| `Type`         | `BMT5_ENUM_POSITION_TYPE`   | Position direction (see enum below).    |
| `Volume`       | `double`                    | Net volume (lots).                      |
| `PriceOpen`    | `double`                    | Open price.                             |
| `PriceCurrent` | `double`                    | Current market price.                   |
| `Swap`         | `double`                    | Accrued swap.                           |
| `Commission`   | `double`                    | Accrued commissions.                    |
| `Profit`       | `double`                    | Floating P/L.                           |
| `StopLoss`     | `double`                    | Stop Loss price.                        |
| `TakeProfit`   | `double`                    | Take Profit price.                      |
| `TimeUpdate`   | `google.protobuf.Timestamp` | Last update time.                       |
| `MagicNumber`  | `int64`                     | EA magic number (if any).               |
| `Identifier`   | `uint64`                    | Position identifier (may match Ticket). |

---

### Enum: `BMT5_ENUM_POSITION_TYPE`

| Value | Name                      | Meaning        |
| ----: | ------------------------- | -------------- |
|     0 | `BMT5_POSITION_TYPE_BUY`  | Long position  |
|     1 | `BMT5_POSITION_TYPE_SELL` | Short position |

---

## 🎯 Purpose

* Show current exposure per symbol for dashboards and risk views.
* Drive position management (modify SL/TP, close, partial close in netting accounts via orders).
* Validate that expected positions exist before running automation.

---

## 🧩 Notes & Tips

* MT5 (netting) keeps **one position per symbol**; multiple fills aggregate into a single net position.
* Use `PositionGet(ctx, symbol)` when you need exactly one symbol; `PositionsGet` returns all.
* For closing or modifying, you’ll need either the **position ticket** (for `PositionModify` / `PositionClose`) or to place an opposite **market order** depending on broker settings.
