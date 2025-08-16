# Getting Started with MetaTrader 5 in Go

Welcome to the **MetaRPC MT5 Go Documentation** — your guide to building integrations with **MetaTrader 5** using **Go** and **gRPC**.

This documentation will help you:

* 📘 Explore all available **account, trading, and market methods**
* 💡 Learn from **Go usage examples** with context and timeout control
* 🔁 Work with **real-time streaming** of quotes, orders, positions, and deals
* ⚙️ Understand all **input/output types**, such as `OrderSendData`, `PositionData`, `QuoteData`, and enums like `ENUM_ORDER_TYPE_TF` or `MRPC_ENUM_TRADE_REQUEST_ACTIONS`
* 🟢 A link to a convenient way to view information https://metarpc.github.io/GoMT5/

---

## 📚 Main Sections

### Quick Account & Market Info

* [Account Summary](QuickAccount_MarketInfo/AccountSummary.md)
* [Quote](QuickAccount_MarketInfo/Quote.md)
* [QuotesMany](QuickAccount_MarketInfo/QuotesMany.md)
* [Symbol Params](QuickAccount_MarketInfo/SymbolParams.md)
* [Tick Values](QuickAccount_MarketInfo/TickValues.md)
* [All Symbols](QuickAccount_MarketInfo/AllSymbols.md)

### Opened State Snapshot

* [Opened Orders](Opened_State_Snapshot/OpenedOrders.md)
* [Opened Order Tickets](Opened_State_Snapshot/OpenedOrderTickets.md)
* [Show Positions](Opened_State_Snapshot/ShowPositions.md)
* [Has Open Position](Opened_State_Snapshot/HasOpenPosition.md)

### Calculations & Safety Checks

* [Order Calc Margin](Calculations_And_PreliminaryVerification/OrderCalcMargin.md)
* [Order Calc Profit](Calculations_And_PreliminaryVerification/OrderCalcProfit.md)
* [Order Check](Calculations_And_PreliminaryVerification/ShowOrderCheck.md)

### Trading Operations ⚠️

* **Section index:** TradingOps (DANGEROUS)
* Popular entries:

  * [Buy Market](TradingOps%28DANGEROUS%29/BuyMarket.md)
  * [Sell Market](TradingOps%28DANGEROUS%29/SellMarket.md)
  * [Place Buy Limit](TradingOps%28DANGEROUS%29/PlaceBuyLimit.md)
  * [Place Sell Limit](TradingOps%28DANGEROUS%29/PlaceSellLimit.md)
  * [Place Buy Stop](TradingOps%28DANGEROUS%29/PlaceBuyStop.md)
  * [Place Sell Stop](TradingOps%28DANGEROUS%29/PlaceSellStop.md)
  * [Place Stop Limit](TradingOps%28DANGEROUS%29/PlaceStopLimit.md)
  * [Order Send Example](TradingOps%28DANGEROUS%29/ShowOrderSendExample.md)
  * [Order Send StopLimit Example](TradingOps%28DANGEROUS%29/ShowOrderSendStopLimitExample.md)
  * [Order Modify Example](TradingOps%28DANGEROUS%29/OrderModifyExample.md)
  * [Order Close Example](TradingOps%28DANGEROUS%29/OrderCloseExample.md)
  * [Order Delete Example](TradingOps%28DANGEROUS%29/OrderDeleteExample.md)
  * [Position Modify](TradingOps%28DANGEROUS%29/PositionModify.md)
  * [Position Close](TradingOps%28DANGEROUS%29/PositionClose.md)
  * [Close All Positions](TradingOps%28DANGEROUS%29/CloseAllPositions.md)
  * [Set Order Expiration](TradingOps%28DANGEROUS%29/SetOrderExpiration.md)

### History & Simple Statistics

* **Section overview:** [History & Stats — Overview](History_And_SimpleStatistics/HistoryAndStats_Overview.md)
* [Show Orders History](History_And_SimpleStatistics/ShowOrdersHistory.md)
* [Show Deals Count](History_And_SimpleStatistics/ShowDealsCount.md)
* [Order By Ticket](History_And_SimpleStatistics/OrderByTicket.md)
* [Deal By Ticket](History_And_SimpleStatistics/DealByTicket.md)
* (Time range guide is linked inside the Overview)

---

## 🚀 Quick Start

1. **Configure your `config.json`** with MT5 credentials and connection details.
2. Initialize an `MT5Account`, wrap it in `MT5Service`.
3. Run examples from `main.go` or call the `Show*` helpers.

```go
ctx := context.Background()
svc := mt5.NewMT5Service(account)

// Example: quick account & quote
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, "EURUSD")
```

---

## 🛠 Requirements

* Go 1.21+
* gRPC-Go
* Protobuf Go bindings (fetched from remote repo via `go.mod`)
* VS Code / GoLand / LiteIDE

---

## 🧭 Navigation Tips

* Sections above link **directly** to the markdown files in your repo.
* The **TradingOps** path uses encoded parentheses: `TradingOps%28DANGEROUS%29`.
* History methods rely on `from/to` ranges — see the Overview inside that section.
