# Getting Started with MetaTrader 5 in Go

Welcome to the **MetaRPC MT5 Go Documentation** — your guide to building integrations with **MetaTrader 5** using **Go** and **gRPC**.

This documentation provides everything you need to:

* 📘 Explore all available **account, trading, and market methods**
* 💡 Learn from **Go usage examples** with context and timeout control
* 🔁 Work with **real-time streaming** of quotes, orders, positions, and deals
* ⚙️ Understand all **input/output types**, including `OrderSendData`, `PositionData`, `QuoteData`, and enums like `ENUM_ORDER_TYPE_TF` or `MRPC_ENUM_TRADE_REQUEST_ACTIONS`

---

## 📚 Main Sections

* **Quick Account & Market Info** — quotes, tick values, trading symbols
* **Opened State Snapshot** — open orders, tickets, and active positions
* **Calculations & Safety Checks** — margin, profit, and pre-trade validation
* **Trading Operations** — sending, modifying, and closing orders (⚠️ dangerous section)
* **History & Statistics** — orders history, deals, and range-based stats
* **Streaming** — subscribe to continuous updates on trades, quotes, and profits

---

## 🚀 Quick Start

To get started with Go + MetaTrader 5:

1. **Prepare your `config.json`** with MT5 credentials and connection details.
2. Initialize an `MT5Account` and wrap it in an `MT5Service` to access helper methods.
3. Run provided examples (`Show*.go` methods) or call services directly from `main.go`.

---

## 🛠 Requirements

* Go 1.20+
* gRPC-Go
* Protobuf-generated Go bindings for MT5 `.proto` definitions
* VS Code, GoLand, or LiteIDE for development

---

With this documentation, you can:

* Monitor account health and exposure
* Automate trade operations safely
* Build dashboards for quotes and market data
* Run backtests and analyze history

Ready to trade with MT5? Let’s Go 🟢
