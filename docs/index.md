# Getting Started with MetaTrader 5 in Go

Welcome to the **MetaRPC MT5 Go Documentation** — your guide to building integrations with **MetaTrader 5** using **Go** and **gRPC**.

This documentation will help you:

* 📘 Explore all available **account, trading, and market methods**
* 💡 Learn from **Go usage examples** with context and timeout control
* 🔁 Work with **real-time streaming** of quotes, orders, positions, and deals
* ⚙️ Understand all **input/output types**, such as `OrderSendData`, `PositionData`, `QuoteData`, and enums like `ENUM_ORDER_TYPE_TF` or `MRPC_ENUM_TRADE_REQUEST_ACTIONS`

---

## 📚 Main Sections

* **[Quick Account & Market Info](QuickAccount_MarketInfo/index.md)** — quotes, tick values, trading symbols
* **[Opened State Snapshot](Opened_State_Snapshot/index.md)** — open orders, tickets, and active positions
* **[Calculations & Safety Checks](Calculations_And_PreliminaryVerification/index.md)** — margin, profit, and pre-trade validation
* **[Trading Operations ⚠️](TradingOps%28DANGEROUS%29/index.md)** — sending, modifying, and closing orders
* **[History & Simple Statistics](History_And_SimpleStatistics/index.md)** — orders history, deals, and range-based stats
* **[Streaming](Streaming/index.md)** — subscribe to updates on trades, quotes, and profits

---

## 🚀 Quick Start

1. **Configure your `config.json`** with MT5 credentials and connection details
2. Use the `MT5Account` or `MT5Service` structs to access functionality
3. Run examples via `main.go` or helper files like `Show*.go`

---

## 🛠 Requirements

* Go 1.21+
* gRPC-Go
* Protobuf bindings (imported automatically from remote repo)
* VS Code / GoLand / LiteIDE

---

## 🧭 Navigation Tips

* Each section has its own **index.md** with explanations and links to sub-methods
* Code examples are always in **Go**, with comments in English
* Dangerous operations (like closing all positions) are flagged ⚠️
