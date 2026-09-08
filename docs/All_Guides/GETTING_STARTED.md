# Getting Started with GoMT5

> **Quick Setup & Overview** - Start automating your trading with GoMT5 (Go for MetaTrader 5) in minutes.

---

## 🎯 What is GoMT5?

**GoMT5** is an industrial-grade, type-safe Go client library for interacting with **MetaTrader 5** terminals via high-performance **gRPC**. It eliminates complex C++ DLL wrappers and provides direct, reliable programmatic trading.

### 🌟 Key Advantages
- 🚀 **High Throughput**: Native gRPC streaming for sub-millisecond price ticks and trade execution.
- 🛡️ **Three-Layer Architecture**: Low-level gRPC (`MT5Account`), typed wrapper methods (`MT5Service`), and high-level convenience (`MT5Sugar`).
- 🔄 **Resilient Connection**: Auto-reconnect, exponential backoff, and transparent channel healing.
- 💼 **Production Ready**: Fully verified across institutional accounts, hedge fund systems, and automated retail bots.

---

## 📦 Installation

Install GoMT5 via your standard Go package manager:

```bash
go get github.com/MetaRPC/GoMT5
```

---

## 🔌 Minimal Connection Example

Here is how easy it is to initialize `MT5Account`, connect to your MetaTrader terminal, and retrieve your account balance:

```
client, err := mt.NewMT5Account(user, password, grpcServer)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err = client.ConnectByServerName(ctx, serverName, "EURUSD")
summary, err := client.AccountSummary(ctx)
fmt.Printf("Balance: %.2f, Equity: %.2f\n", summary.Balance, summary.Equity)
```

---

## 🗺️ Documentation Road Map

To get the most out of GoMT5, follow this suggested reading order:

1. 🚀 **[Your First Project](Your_First_Project.md)** - Build and run a working project in 10 minutes.
2. 🗺️ **[Project Map](PROJECT_MAP.md)** - Understand the 3 architectural layers and interaction flow.
3. 📖 **[Glossary](GLOSSARY.md)** - Essential MetaTrader 5 and algorithmic trading terminology.
4. 🐣 **[MT5 for Beginners](MT5_For_Beginners.md)** - Step-by-step terminal setup and demo account guide.
5. 📡 **[gRPC Streaming](GRPC_STREAM_MANAGEMENT.md)** - Subscribe to ticks, trades, DOM, and position updates.
6. 📊 **[Return Codes](RETURN_CODES_REFERENCE.md)** - Complete reference of broker and gateway return codes.
