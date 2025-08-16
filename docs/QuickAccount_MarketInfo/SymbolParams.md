# Getting Symbol Parameters

> **Request:** fetch detailed symbol metadata and live attributes.

---

### Code Example

```go
// High-level (prints selected attributes):
svc.ShowSymbolParams(ctx, "EURUSD")

// Low-level (access full struct):
info, err := svc.account.SymbolParams(ctx, "EURUSD")
if err != nil {
    log.Printf("❌ SymbolParams error: %v", err)
    return
}
fmt.Println("Symbol:", info.GetName())
fmt.Println("Description:", info.GetSymDescription())
fmt.Printf("Digits: %d\n", info.GetDigits())
fmt.Printf("Volume Min/Max/Step: %.2f / %.2f / %.2f\n", info.GetVolumeMin(), info.GetVolumeMax(), info.GetVolumeStep())
fmt.Println("Trade Mode:", info.GetTradeMode())
fmt.Printf("Currencies: base=%s profit=%s margin=%s\n", info.GetCurrencyBase(), info.GetCurrencyProfit(), info.GetCurrencyMargin())
```

---

### Method Signature

```go
func (s *MT5Service) ShowSymbolParams(ctx context.Context, symbol string)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC. |
| `symbol`  | `string`          | yes      | Trading symbol name (e.g., `"EURUSD"`).    |

---

## ⬆️ Output

Returns a **`SymbolParameters`** struct with extensive fields; `ShowSymbolParams` prints a subset.

### Commonly used fields

| Field            | Type                          | Description                                 |
| ---------------- | ----------------------------- | ------------------------------------------- |
| `Name`           | `string`                      | Symbol name.                                |
| `SymDescription` | `string`                      | Human-readable description.                 |
| `Digits`         | `int32`                       | Number of digits (price precision).         |
| `VolumeMin`      | `float64`                     | Minimum lot size.                           |
| `VolumeMax`      | `float64`                     | Maximum lot size.                           |
| `VolumeStep`     | `float64`                     | Lot step.                                   |
| `TradeMode`      | `BMT5_ENUM_SYMBOL_TRADE_MODE` | Trading availability mode (see enum below). |
| `CurrencyBase`   | `string`                      | Base currency.                              |
| `CurrencyProfit` | `string`                      | Profit currency.                            |
| `CurrencyMargin` | `string`                      | Margin currency.                            |

---

### Enum: `BMT5_ENUM_SYMBOL_TRADE_MODE`

| Value | Name                               | Meaning                            |
| ----: | ---------------------------------- | ---------------------------------- |
|     0 | `BMT5_SYMBOL_TRADE_MODE_DISABLED`  | Trade disabled for this symbol.    |
|     1 | `BMT5_SYMBOL_TRADE_MODE_LONGONLY`  | Only long positions allowed.       |
|     2 | `BMT5_SYMBOL_TRADE_MODE_SHORTONLY` | Only short positions allowed.      |
|     3 | `BMT5_SYMBOL_TRADE_MODE_CLOSEONLY` | Only closing positions is allowed. |
|     4 | `BMT5_SYMBOL_TRADE_MODE_FULL`      | No restrictions (full trading).    |

> Enum source: generated `pb` (account-helper).

---

## 📦 Full field inventory (SymbolParameters)

Below is the complete list of fields exposed by `SymbolParameters` (names and Go types). Use getters with the same names prefixed by `Get` in your code (e.g., `GetDigits()`).

| Field            | Go Type                       | JSON tag                    |
| ---------------- | ----------------------------- | --------------------------- |
| `Name`           | `string`                      | `name,omitempty`            |
| `Bid`            | `float64`                     | `bid,omitempty`             |
| `BidHigh`        | `float64`                     | `bid_high,omitempty`        |
| `BidLow`         | `float64`                     | `bid_low,omitempty`         |
| `Ask`            | `float64`                     | `ask,omitempty`             |
| `AskHigh`        | `float64`                     | `ask_high,omitempty`        |
| `AskLow`         | `float64`                     | `ask_low,omitempty`         |
| `Last`           | `float64`                     | `last,omitempty`            |
| `LastHigh`       | `float64`                     | `last_high,omitempty`       |
| `LastLow`        | `float64`                     | `last_low,omitempty`        |
| `VolumeMin`      | `float64`                     | `volume_min,omitempty`      |
| `VolumeMax`      | `float64`                     | `volume_max,omitempty`      |
| `VolumeStep`     | `float64`                     | `volume_step,omitempty`     |
| `Digits`         | `int32`                       | `digits,omitempty`          |
| `TradeMode`      | `BMT5_ENUM_SYMBOL_TRADE_MODE` | `trade_mode,omitempty`      |
| `CurrencyBase`   | `string`                      | `currency_base,omitempty`   |
| `CurrencyProfit` | `string`                      | `currency_profit,omitempty` |
| `CurrencyMargin` | `string`                      | `currency_margin,omitempty` |
| `SymDescription` | `string`                      | `sym_description,omitempty` |
The struct contains 111 fields in total (pricing, sessions, swaps, margin/contract, sector/industry, options, etc.).

---

## 🎯 Purpose

* Introspect tradability, precision, and lot constraints before placing orders.
* Show instrument metadata in dashboards.
* Determine currencies used for P/L and margin calculations.

---

## 🧩 Notes & Tips

* Some fields are dynamic (Bid/Ask/Last, session stats) and reflect current tick; others are static metadata (digits, contract size).
* `TradeMode` reflects broker-side permissions and may differ by account type.
* If you query many symbols at once, consider `SymbolParamsMany` for pagination.
