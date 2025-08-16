# 🧰 Using GoMT5 via CLI (No GUI)

This guide shows how to run **GoMT5** straight from the terminal — no GUI, only code and the console. Perfect for devs, ops, and anyone who prefers scripts over buttons.

---

## 🔧 Requirements

| Tool / Thing           | Purpose                                           |
| ---------------------- | ------------------------------------------------- |
| Go 1.21+               | Build & run the project                           |
| MetaTrader 5           | Terminal with MetaRPC plugin / connector running  |
| `config.json`          | Login credentials, server name, default symbol    |
| Terminal (cmd/PS/Bash) | All operations are executed from the command line |

> Proto note: `.proto` files are **not** stored locally — Go bindings are pulled from the remote repo via `go.mod`.

---

## 📁 Project Structure (your repo)

```bash
GoMT5/
├─ docs/                                 # Documentation (what you're reading)
├─ config/
│  ├─ config.go                           # Loads config.json
│  └─ config.json                         # Login/Server/DefaultSymbol
├─ mt5/
│  ├─ MT5Account.go                       # Low-level account & connection helpers
│  └─ MT5Service.go                       # High-level helpers (Show*/Place*/Buy/Sell)
├─ main.go                                # Entry point to run examples from code
├─ go.work / go.work.sum                  # Workspace files
├─ mkdocs.yml                             # Docs site config (optional)
└─ examples/                              # (optional) extra runnable samples
```

---

## 🧩 Example `config.json`

```json
{
  "Login": 501401178,
  "Password": "v8gctta",
  "Server": "RoboForex-Demo",
  "DefaultSymbol": "EURUSD"
}
```

> You can also set a proxy via env var:
>
> * `MT5_PROXY=http://host:port`
> * `MT5_PROXY=socks5://user:pass@host:port`

---

## 🚀 Run It

From the **repo root**:

```bash
# Windows PowerShell / cmd / Bash — all the same
# make sure you're in the folder where main.go lives

go run .
```

If all is well, you’ll see logs like:

```
using proxy: ... (or "no proxy set")
connect(wait on server ...) to RoboForex-...
terminal is ready
symbol ready: EURUSD
✅ Done.
```

**Typical fix if you get** `no Go files` — you’re not in the folder with `main.go`. `cd` to repo root and re-run `go run .`.

---

## 🧪 What You Can Call (by blocks)

### 1) Quick Account & Market Info

* `svc.ShowAccountSummary(ctx)` — balance/equity/currency
* `svc.ShowQuote(ctx, symbol)` — live bid/ask
* `svc.ShowQuotesMany(ctx, []string{"EURUSD","GBPUSD"})`
* `svc.ShowSymbolParams(ctx, symbol)`
* `svc.ShowTickValues(ctx, []string{...})`
* `svc.ShowAllSymbols(ctx)` *(prints a lot)*

### 2) Opened State Snapshot

* `svc.ShowOpenedOrders(ctx)`
* `svc.ShowOpenedOrderTickets(ctx)`
* `svc.ShowPositions(ctx)`
* `svc.ShowHasOpenPosition(ctx, symbol)`

### 3) Calculations & Pre-Trade Checks (safe)

* `svc.ShowOrderCalcMargin(ctx, symbol, orderType, volume, openPrice)`
* `svc.ShowOrderCalcProfit(ctx, symbol, orderType, volume, open, close)`
* `svc.ShowOrderCheck(ctx, action, orderType, symbol, volume, price, sl, tp, deviation, magic, exp)`

### 4) Trading Ops (⚠️ real actions)

* `svc.BuyMarket(ctx, symbol, volume, sl, tp)`
* `svc.SellMarket(ctx, symbol, volume, sl, tp)`
* `svc.PlaceBuyLimit(ctx,  symbol, volume, price, sl, tp, exp)`
* `svc.PlaceSellLimit(ctx, symbol, volume, price, sl, tp, exp)`
* `svc.PlaceBuyStop(ctx,   symbol, volume, trigger, sl, tp, exp)`
* `svc.PlaceSellStop(ctx,  symbol, volume, trigger, sl, tp, exp)`
* `svc.PlaceStopLimit(ctx, isBuy,  volume, trigger, limit, sl, tp, exp)`
* `svc.ShowOrderModifyExample(ctx, ticket)`
* `svc.ShowOrderCloseExample(ctx,  ticket)`
* `svc.ShowOrderDeleteExample(ctx, ticket)`
* `svc.ShowPositionModify(ctx, ticket, newSL, newTP)`
* `svc.ShowPositionClose(ctx, symbol)`
* `svc.ShowCloseAllPositions(ctx)` *(handle with care!)*

### 5) History & Simple Stats (read-only)

* `svc.ShowOrdersHistory(ctx, from, to)`
* `svc.ShowDealsCount(ctx, from, to, symbol)`
* `svc.ShowOrderByTicket(ctx, orderTicket)`
* `svc.ShowDealByTicket(ctx, dealTicket)`

> Need help with `from/to` ranges? See the doc: **History Range (from/to) — How to Use**.

---

## 🔄 Streaming / Subscriptions

* `svc.StreamQuotes(ctx)` — live ticks per symbol(s)
* `svc.StreamOpenedOrderProfits(ctx)` — real-time P/L per open order
* `svc.StreamOpenedOrderTickets(ctx)` — current open tickets
* `svc.StreamTradeUpdates(ctx)` — trade events (new/updated orders)

Example output:

```
[Tick] EURUSD | Bid: 1.09876 | Ask: 1.09889 | Time: 2025-07-29 18:00:01
```

---

## 🧑‍💻 Enabling a Function in `main.go`

Just call the helpers you need:

```go
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, "EURUSD")
// svc.BuyMarket(ctx, "EURUSD", 0.10, nil, nil) // ⚠️ real trade — keep commented until ready
```

You can chain them — e.g., open a pending order and immediately subscribe to updates.

---

## 🧠 Tips

* Prefer `context.WithTimeout` for network calls.
* Set `MT5_PROXY` if you’re behind a corporate firewall.
* Even on **demo** — trades are real for the broker side. Test carefully.
* Use `timestamppb.New(time.Now().Add(...))` to set **expiration** for pending orders (see *SetOrderExpiration*).

---

## 📎 Quick Example

```go
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, "EURUSD")
svc.ShowOpenedOrders(ctx)
svc.StreamQuotes(ctx)
```

Minimal. Fast. Scriptable. Exactly what CLI folks love. 🧪⚡
