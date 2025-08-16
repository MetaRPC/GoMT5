# Quick Account & Market Info — Overview

This section contains **fast-access methods** for retrieving account snapshots and market symbol information.
These functions are useful for building dashboards, performing quick checks, and fetching basic trading data.

---

## 📂 Methods in this Section

* [AccountSummary.md](AccountSummary.md)
  Get a quick summary of the trading account (balance, equity, margin, free margin).

* [AllSymbols.md](AllSymbols.md)
  Retrieve a list of all available trading symbols in the terminal.

* [Quote.md](Quote.md)
  Get a **single latest quote** (bid/ask/last) for a chosen symbol.

* [QuotesMany.md](QuotesMany.md)
  Retrieve quotes for **multiple symbols at once**.

* [SymbolParams.md](SymbolParams.md)
  Fetch detailed parameters for a specific symbol (spread, digits, trade mode, etc.).

* [TickValues.md](TickValues.md)
  Get tick size and tick value for a symbol — useful for pip-value calculations.

---

## ⚡ Example Usage

```go
// Get account summary
summary, _ := svc.AccountSummary(ctx)

// Fetch one symbol quote
quote, _ := svc.Quote(ctx, "EURUSD")

// Fetch multiple quotes
quotes, _ := svc.QuotesMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})

// Get symbol parameters
params, _ := svc.SymbolParams(ctx, "EURUSD")
```

---

## ✅ Best Practices

1. Use **`AccountSummary`** for quick UI dashboards (instead of full account info calls).
2. Call **`AllSymbols`** once and cache results — avoid fetching repeatedly.
3. Use **`QuotesMany`** if you monitor multiple pairs to reduce API calls.
4. For **pip value calculations**, always rely on **`TickValues`**.
5. Combine **`SymbolParams`** with **`Quote`** to validate trade conditions.

---

## 🎯 Purpose

The methods in this block allow you to:

* Build quick **account snapshots** (balance/equity widgets).
* Show **real-time quotes** in dashboards.
* Perform **risk calculations** (pip value, spread, leverage impact).
* Verify **symbol trade settings** before sending orders.

---

👉 Use this overview as a **map**, and jump into each `.md` file for detailed method documentation.
