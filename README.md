# GoMT5 SDK

> Official MetaTrader 5 Go SDK for high-performance trading automation and market data streaming via gRPC.

[![Docs](https://img.shields.io/badge/docs-GoMT5-blue.svg)](https://metarpc.github.io/GoMT5/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

📄 **Full Documentation**: [https://metarpc.github.io/GoMT5/](https://metarpc.github.io/GoMT5/)

---

## 📦 Installation

```bash
go get github.com/MetaRPC/GoMT5
```

---

## 🏃 How to Run Examples

Clone the repository and run the examples out-of-the-box:

```bash
git clone https://github.com/MetaRPC/GoMT5.git
cd GoMT5/examples/demos
go mod download

# 1. Run with default TRIAL key:
go run main.go lowlevel01

# 2. Or pass your MetaRPC API key directly as an argument:
go run main.go lowlevel01 your_api_key_here

# 3. Or use the MRPC_API_KEY environment variable:
export MRPC_API_KEY="your_api_key_here"        # Windows CMD: set MRPC_API_KEY=your_api_key_here
go run main.go lowlevel01                      # Windows PowerShell: $env:MRPC_API_KEY="your_api_key_here"
```

---

## 🔑 API Key & Authentication

Connecting to MetaRPC production endpoints (`mt5.mrpc.pro:443`) requires an API key:

1. **Sign Up**: Create an account for free at [https://mrpc.pro/signup](https://mrpc.pro/signup).
2. **Generate API Key**: In your MetaRPC Portal dashboard at [https://mrpc.pro/my](https://mrpc.pro/my), go to **API Keys** to generate and copy your personal API token.
3. **Configure Connection**: Pass your API key / token along with the server address (`mt5.mrpc.pro:443`) in your connection settings.

---

## 🆔 Automatic Account ID & Authentication

MetaRPC endpoints route calls using a terminal session identifier (`id`):
- When connecting via `Connect` / `ConnectEx`, the server automatically generates and returns a session GUID (`terminalInstanceGuid`).
- The SDK automatically captures this session ID and attaches both the `id` and `APIKey` headers on all subsequent calls (`AccountSummary`, `OrderSend`, streaming, etc.) — no manual `GetId` or `curl` calls required.
- Pass your API key directly to the Account/Client constructor or via the `MRPC_API_KEY` environment variable.

## 🌐 Production Endpoints

| Environment | Host | Port | Protocol |
| :--- | :--- | :--- | :--- |
| **Production** | `mt5.mrpc.pro` | `443` | TLS / gRPC |
| **Direct API UI (Swagger)** | `https://mt5.mrpc.pro/apiui` | `443` | HTTPS |
| **Portal Dashboard** | `https://mrpc.pro/my` | `443` | HTTPS |
| **Registration / API Key** | `https://mrpc.pro/signup` | `443` | HTTPS |

---

## 📄 Documentation & Guides

Explore comprehensive documentation at [https://metarpc.github.io/GoMT5/](https://metarpc.github.io/GoMT5/):
- 🚀 **Quick Start & First Project**
- 🔑 **Authentication & API Keys**
- 📡 **Live Market Data & gRPC Streaming**
- 💼 **Account Management & Order Execution**
- 📊 **Return Codes & Error Handling Reference**
