# Getting Tick Value / Tick Size / Contract Size

> **Request:** fetch monetary tick value, tick size (price step), and contract size for one or more symbols.

---

### Code Example

```go
// High-level (prints selected fields):
svc.ShowTickValues(ctx, []string{"EURUSD", "GBPUSD"})

// Low-level (work with the payload directly):
data, err := svc.account.TickValueWithSize(ctx, []string{"EURUSD", "GBPUSD"})
if err != nil {
    log.Printf("❌ TickValueWithSize error: %v", err)
    return
}
for _, info := range data.GetSymbolTickSizeInfos() {
    fmt.Printf("%s | TickValue=%.5f | TickSize=%.5f | ContractSize=%.2f\n",
        info.GetName(), info.GetTradeTickValue(), info.GetTradeTickSize(), info.GetTradeContractSize())
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowTickValues(ctx context.Context, symbols []string)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                        |
| --------- | ----------------- | -------- | -------------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC.         |
| `symbols` | `[]string`        | yes      | List of symbol names, e.g., `{"EURUSD","GBPUSD"}`. |

**Rules & notes**

* Symbols should be visible/enabled at the terminal/broker.
* Empty slice likely results in an error or empty payload.

---

## ⬆️ Output

Returns a struct containing an array of **`SymbolTickSizeInfo`** objects.
`ShowTickValues` prints the following fields per symbol:

| Field               | Type     | Description                                                              |
| ------------------- | -------- | ------------------------------------------------------------------------ |
| `Name`              | `string` | Symbol name.                                                             |
| `TradeTickValue`    | `double` | Monetary value of **one price tick** in account currency.                |
| `TradeTickSize`     | `double` | **Price step** (minimum price increment) for the symbol.                 |
| `TradeContractSize` | `double` | Contract size used for calculations (e.g., 100,000 for many FX symbols). |

> Exact field names are taken from generated getters (`GetName`, `GetTradeTickValue`, `GetTradeTickSize`, `GetTradeContractSize`).

---

## 🎯 Purpose

* Convert price movements into money (e.g., P/L per tick for position sizing).
* Validate lot sizing and risk before sending orders.
* Display instrument math in dashboards.

---

## 🧩 Notes & Tips

* `TradeTickValue` is expressed in **account currency**; ensure this matches what you expect when mixing instruments.
* Some brokers define non-standard contract sizes (indices, commodities) — always read from this API, don’t hardcode.
* Pair with `SymbolParams` (digits, volume step) to compute P/L per pip and validate order quantities.
