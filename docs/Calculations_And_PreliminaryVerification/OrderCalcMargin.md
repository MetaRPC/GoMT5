# Calculating Required Margin

> **Request:** compute margin required for a hypothetical trade before sending it.

---

### Code Example

```go
// High-level (prints the margin in account currency):
svc.ShowOrderCalcMargin(ctx, "EURUSD", pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, 0.10, 0)

// Low-level (full control):
req := &pb.OrderCalcMarginRequest{
    Symbol:    "EURUSD",
    OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
    Volume:    0.10,     // lots
    OpenPrice: 0,        // 0 for market; for pending use the intended price
}
res, err := svc.account.OrderCalcMargin(ctx, req)
if err != nil {
    log.Printf("❌ OrderCalcMargin error: %v", err)
    return
}
fmt.Printf("🧮 Margin required: %.2f\n", res.GetMargin())
```

---

### Method Signature

```go
func (s *MT5Service) ShowOrderCalcMargin(ctx context.Context, symbol string, orderType pb.ENUM_ORDER_TYPE_TF, volume float64, openPrice float64)
```

---

## 🔽 Input

| Parameter   | Type                    | Required | Description                                                                |
| ----------- | ----------------------- | -------- | -------------------------------------------------------------------------- |
| `ctx`       | `context.Context`       | yes      | Controls timeout/cancellation for the RPC.                                 |
| `symbol`    | `string`                | yes      | Trading symbol name (e.g., `"EURUSD"`).                                    |
| `orderType` | `pb.ENUM_ORDER_TYPE_TF` | yes      | Order kind (market/pending) — see enum below.                              |
| `volume`    | `float64`               | yes      | Trade volume in **lots** (respect broker `VolumeMin/Max/Step`).            |
| `openPrice` | `float64`               | yes      | Price used for the calc. Use **0 for market**; for pendings pass entry px. |

### Enum: `ENUM_ORDER_TYPE_TF`

| Code | Name                            | Meaning                    |
| ---: | ------------------------------- | -------------------------- |
|    0 | `ORDER_TYPE_TF_BUY`             | Market Buy                 |
|    1 | `ORDER_TYPE_TF_SELL`            | Market Sell                |
|    2 | `ORDER_TYPE_TF_BUY_LIMIT`       | Pending Buy Limit          |
|    3 | `ORDER_TYPE_TF_SELL_LIMIT`      | Pending Sell Limit         |
|    4 | `ORDER_TYPE_TF_BUY_STOP`        | Pending Buy Stop           |
|    5 | `ORDER_TYPE_TF_SELL_STOP`       | Pending Sell Stop          |
|    6 | `ORDER_TYPE_TF_BUY_STOP_LIMIT`  | Pending Buy Stop-Limit     |
|    7 | `ORDER_TYPE_TF_SELL_STOP_LIMIT` | Pending Sell Stop-Limit    |
|    8 | `ORDER_TYPE_TF_CLOSE_BY`        | Close by opposite position |

> Enum names/codes are taken from `mt5-term-api-trade-functions.pb.go` (`ENUM_ORDER_TYPE_TF`).

---

## ⬆️ Output

Returns **`OrderCalcMarginData`**.

| Field    | Type      | Description                                                     |
| -------- | --------- | --------------------------------------------------------------- |
| `Margin` | `float64` | Margin required in **account currency** for the given scenario. |

---

## 🎯 Purpose

* Validate affordability of a trade before `OrderSend`.
* Pre-compute margin usage for risk checks and sizing logic.
* What-if analysis for pending orders at specific prices.

---

## 🧩 Notes & Tips

* The result may depend on symbol settings (contract size, leverage tiers, margin currency). Pair with `SymbolParams`.
* For **market** scenarios pass `openPrice=0` — the server uses current prices.
* For **pending** scenarios pass your intended entry price; margin model can differ from market.
* Margin models can change with volatility or session; treat the result as an estimate, not a guarantee.
