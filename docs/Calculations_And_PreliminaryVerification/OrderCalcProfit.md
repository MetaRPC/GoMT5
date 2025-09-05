# Calculating Profit for a Hypothetical Trade

> **Request:** compute P/L for a given order type, volume, and price change.

---

### Code Example

```go
// High-level (prints computed P/L in account currency):
svc.ShowOrderCalcProfit(ctx, "EURUSD", pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, 0.10, 1.08000, 1.08350)

// Low-level (manual request/response):
req := &pb.OrderCalcProfitRequest{
    OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, // BUY/SELL or a pending type
    Symbol:     "EURUSD",
    Volume:     0.10,        // lots
    OpenPrice:  1.08000,     // entry price (or current for market what-if)
    ClosePrice: 1.08350,     // exit price used for P/L calc
}
res, err := svc.account.OrderCalcProfit(ctx, req)
if err != nil {
    log.Printf("❌ OrderCalcProfit error: %v", err)
    return
}
fmt.Printf("💰 Profit calc: %.2f\n", res.GetProfit())
```

---

### Method Signature

```go
func (s *MT5Service) ShowOrderCalcProfit(ctx context.Context, symbol string, orderType pb.ENUM_ORDER_TYPE_TF, volume float64, openPrice, closePrice float64)
```

---

## 🔽 Input

| Parameter    | Type                    | Required | Description                                              |
| ------------ | ----------------------- | -------- | -------------------------------------------------------- |
| `ctx`        | `context.Context`       | yes      | Controls timeout/cancellation for the RPC.               |
| `symbol`     | `string`                | yes      | Trading symbol (e.g., `"EURUSD"`).                       |
| `orderType`  | `pb.ENUM_ORDER_TYPE_TF` | yes      | Order kind (BUY/SELL or pending) — see enum below.       |
| `volume`     | `float64`               | yes      | Trade volume in **lots** (respect `VolumeMin/Max/Step`). |
| `openPrice`  | `float64`               | yes      | Entry price used for the what‑if scenario.               |
| `closePrice` | `float64`               | yes      | Exit price used for the what‑if scenario.                |

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

---

## ⬆️ Output

Returns **`OrderCalcProfitData`**.

| Field    | Type      | Description                                            |
| -------- | --------- | ------------------------------------------------------ |
| `Profit` | `float64` | P/L in **account currency** for the provided scenario. |

---

## 🎯 Purpose

* What‑if analysis for risk management and take‑profit/stop‑loss planning.
* Validate strategy math without sending any orders.
* Back‑of‑the‑envelope checks for dashboards and bots.

---

## 🧩 Notes & Tips

* The result depends on contract size, tick value, and digits — pair with `SymbolParams` and `TickValueWithSize`.
* For BUY positions: Profit grows as `closePrice` > `openPrice`; for SELL — on the contrary.
* Consider spreads/commissions/slippage separately; this call focuses on price P/L, not full trade cost.
