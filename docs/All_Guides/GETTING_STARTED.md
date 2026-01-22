# Getting Started with GoMT5

> **Welcome to GoMT5** - a comprehensive educational project for learning MT5 trading automation from scratch using Go.

---

## 🚀 Prerequisites and Installation

Before starting work with GoMT5, you need to set up your development environment.

### Step 1: Install Go 1.21+

GoMT5 requires Go version 1.21 or higher.

**Download and install:**

- **Official website:** [Go Downloads](https://go.dev/dl/)
- Choose the installer for your platform

**Verify installation:**

```bash
go version
# Should show: go1.21.x or higher
```

---

### Step 2: Install Code Editor

**Visual Studio Code (Recommended):**

- **Download:** [VS Code](https://code.visualstudio.com/)

- **Extensions to install:**

  - Go (official Go extension from Google)
  - Go Test Explorer (optional)
  - Go Template Support (optional)

**GoLand (Alternative):**

- **Download:** [GoLand by JetBrains](https://www.jetbrains.com/go/)

- Professional IDE for Go with advanced features

**You can also use:**

- Vim/Neovim with vim-go
- Sublime Text
- Any text editor + command line

---

### Step 3: Clone the Repository

Clone the GoMT5 project from GitHub:

```bash
git clone https://github.com/MetaRPC/GoMT5
cd GoMT5
```

**If you don't have Git installed:**

- Download from [git-scm.com](https://git-scm.com/)
- Or download the project as ZIP from GitHub and extract

---

### Step 4: Understanding the Connection Flow

GoMT5 connects to the MT5 terminal via **gRPC gateway**.

**Connection flow:**
```
GoMT5 → gRPC → mt5term Gateway → MT5 Terminal
```

**What is mt5term Gateway?**

- External gateway process that bridges GoMT5 with the MT5 terminal
- Handles connection pooling and session management
- Connection settings are specified in `examples/demos/config/config.json`

### 🗂️ Configuration File: config.json

- Before starting work, let's specify your actual credentials in the configuration file.
- It's better to use a Demo account for learning.

**Option 1: In config.json located at path `examples/demos/config/config.json` there are all fields you need to fill.**


```json
{
  "user": 591129415,
  "password": "IpoHj17tYu67@",
  "host": "mt5.mrpc.pro",
  "port": 443,
  "grpc_server": "mt5.mrpc.pro:443",
  "mt_cluster": "FxPro-MT5 Demo",
  "test_symbol": "EURUSD",
  "test_volume": 0.01
}
```

**Option 2: Or set up environment variables and substitute your credentials.**

```bash
# Linux/Mac
export MT5_USER=591129415
export MT5_PASSWORD="YourPassword"
export MT5_HOST="mt5.mrpc.pro"
export MT5_PORT=443
export MT5_CLUSTER="FxPro-MT5 Demo" (Or any other broker-server)

# Windows PowerShell
$env:MT5_USER="591129415"
$env:MT5_PASSWORD="YourPassword"
$env:MT5_HOST="mt5.mrpc.pro"
```

**Configuration parameters explanation:**

| Parameter | Description | Example |
|-----------|-------------|---------|
| **user** | Your MT5 account login number | `591129415` |
| **password** | Your MT5 account password (master password) | `"YourPassword"` |
| **mt_cluster** | MT5 server name from your broker | `"FxPro-MT5 Demo"` |
| **host** | Gateway host address (provided by MetaRPC) | `"mt5.mrpc.pro"` |
| **port** | Gateway port number | `443` (HTTPS) |
| **grpc_server** | Full gRPC server address (combination of host:port) | `"mt5.mrpc.pro:443"` |
| **test_symbol** | Default trading symbol for examples | `"EURUSD"` |
| **test_volume** | Default volume for testing | `0.01` |

**Important notes:**

- **user, password, mt_cluster** - These are your MT5 account credentials
- **host, port** - Provided by the MetaRPC team
- **grpc_server** - Must match the `host:port` format
- **test_symbol** - Change to your preferred trading symbol
- **test_volume** - Adjust for your account size

---

## 📋 MT5 Account Setup

If you don't have an MT5 demo account yet or need help creating one, refer to the beginner's guide:

👉 **[MT5 for Beginners - Creating Demo Account](MT5_For_Beginners/)**

This guide covers:

- Downloading and installing MT5 terminal
- Step-by-step demo account creation
- Understanding master password and investor password
- Choosing a broker (optional)

---

## 🎯 About the Project

This project is a **demonstration of capabilities** of our team's gateway for reproducing methods and functionality. It is designed to help you build your own trading logic system in the future.

We will guide you through all major aspects - from basic manual trading to a fully customizable algorithmic trading system. This journey will reveal the full potential of your acquired knowledge and fundamental understanding of trading and markets.

**What you will learn:**

- What gRPC methods do and how to use them directly
- How methods can be modified for your needs
- How to optimize your Go code for performance
- How to create convenient input/output systems
- How to effectively track positions by symbols
- How to build intelligent risk management systems
- How to work with Go channels and goroutines for real-time streaming

**All we ask from you:**

> **The desire to learn, learn, and learn again.** In the end, this will lead to significant results and, most importantly, to a solid foundation of knowledge in algorithmic trading with Go.

---

## 🏗️ Project Architecture: Three-Tier System

The project consists of **three interconnected files** in the `examples/mt5/` directory, each building upon the previous one. Understanding this chain is key to mastering GoMT5.

### Tier 1: MT5Account.go - Low-Level gRPC Foundation

**What it is:** Direct gRPC calls to the MT5 terminal - the absolute foundation of everything.

**[📖 MT5Account Master Overview](../MT5Account/MT5Account.Master.Overview.md)**

**File:** `package/Helpers/MT5Account.go` (2300+ lines)

- Raw protocol buffer messages and gRPC communication
- Maximum control and flexibility over each request/response
- **All other tiers use this internally**
- Automatic reconnection with exponential backoff
- **40+ methods**, covering all MT5 operations
- **Best for:** Advanced users who need granular control

**Key features:**

- Context-based timeout management
- Two-channel streaming (dataChan, errChan)
- Built-in retry logic for network failures
- Session management with UUID tracking

**📚 Two-Level Documentation:**

For each MT5Account method, **two types of documentation** are available:

1. **Main reference** - brief method description, parameters, return values
   - Example: `docs/MT5Account/1. Account_information/AccountInfoDouble.md`

2. **HOW_IT_WORK** - detailed step-by-step explanation with live code examples
   - Example: `docs/MT5Account/HOW_IT_WORK/1. Account_information_HOW/AccountInfoDouble_HOW.md`
   - Line-by-line breakdown of actual code from demo files
   - Detailed comments for each step
   - Explanation of protobuf structures and gRPC calls operation

This allows you to quickly find reference OR deeply understand how the method works.

### Tier 2: MT5Service.go - Convenient Wrappers

**What it is:** Wrapper methods that simplify working with MT5Account gRPC calls.

**[📖 MT5Service API Overview](../MT5Service/MT5Service.Overview.md)**

**File:** `examples/mt5/MT5Service.go` (1160+ lines)

- Simplified error handling and response parsing
- Pre-configured common operations
- Go-style return types (no protobuf in signatures)
- Easier to work with than raw gRPC
- **50+ methods** for common trading scenarios
- **Best for:** Most common trading scenarios

**Key features:**
- Clean API without using protobuf
- Helper functions for common patterns
- Structured error handling
- Type-safe operations

### Tier 3: MT5Sugar.go - High-Level Helpers

**What it is:** Syntactic sugar and helper methods for maximum productivity.

**[📖 MT5Sugar API Overview](../MT5Sugar/MT5Sugar.API_Overview.md)**

**File:** `examples/mt5/MT5Sugar.go` (2200+ lines)

- One-line operations for common tasks
- Smart defaults and automatic parameter inference
- Risk-based position size calculation
- SL/TP calculation based on pips
- Most intuitive and beginner-friendly
- **Best for:** Quick prototyping and simple strategies

**Key features:**

- `BuyMarket(symbol, volume)` - one line to open a position
- `CalculatePositionSize(symbol, riskPercent, slPips)` - automatic volume calculation by risk
- `BuyMarketWithPips(symbol, volume, slPips, tpPips)` - open with SL/TP in pips
- `GetBalance()`, `GetEquity()` - instant account information
- `CloseAllPositions()` - emergency exit

---

### 📚 Understanding the Chain

This three-file chain represents the evolution from low-level control to high-level convenience:

```
MT5Sugar.go (easiest, highest abstraction)
    ↓ uses
MT5Service.go (convenient wrappers)
    ↓ uses
MT5Account.go (raw gRPC, foundation)
    ↓ communicates with
MT5 Terminal (via gateway)
```

**Each overview document includes:**

- Detailed method descriptions with parameters
- Return types and error handling patterns
- Usage examples and best practices
- Common patterns and pitfalls to avoid
- Go-specific guidance (goroutines, channels, context)

---

### 🎓 Recommended Learning Paths

**Path A: For Go Developers (Bottom-Up Approach)**

If you have experience with Go and want to understand everything deeply:

1. **Start with MT5Account** - Study the gRPC foundation and context patterns
2. **Move to MT5Service** - Understand convenient wrappers
3. **Finish with MT5Sugar** - Appreciate high-level abstractions

✅ This path gives you full control and deep understanding.

**Path B: For Traders (Top-Down Approach)**

If you're new to trading automation and want quick results:

1. **Start with MT5Sugar** - Easy, intuitive methods for quick trading
2. **Move to MT5Service** - Learn more advanced patterns as needed
3. **Deep dive into MT5Account** - Understand the foundation for full control

✅ This path allows you to start trading quickly, leaving room for growth.

---

## 📂 Demo Examples

You can explore demo files that showcase various aspects of the SDK. These files are organized by complexity level and are located in the `examples/demos/` folder.

**Each file includes inline code comments explaining what each operation does.**

### examples/demos/lowlevel/ (MT5Account - Low Level)

Protobuf/gRPC methods - full control, maximum flexibility:

- **01_general_operations.go** - Information methods (account, symbols, positions, ticks)
- **02_trading_operations.go** - Trading operations (calculations, validation, orders)
- **03_streaming_methods.go** - Streaming methods (real-time ticks, deals, positions)

### examples/demos/service/ (MT5Service - Mid Level)

Go wrappers over protobuf - native types, more convenient to work with:

- **04_service_demo.go** - Service API wrappers (account, symbols, positions, orders)
- **05_service_streaming.go** - Service API streaming methods (ticks, deals, profits)

### examples/demos/sugar/ (MT5Sugar - High Level)

High-level API - maximum simplification, one line of code:

- **06_sugar_basics.go** - Basics: connection, balance, prices
- **07_sugar_trading.go** - Trading: market/pending orders, SL/TP
- **08_sugar_positions.go** - Positions: queries, modification, closing
- **09_sugar_history.go** - History: deals, profit over period
- **10_sugar_advanced.go** - Advanced: risk management, symbol info, helpers

### examples/demos/orchestrators/ (Strategies)

Full-fledged trading strategies - examples of complex automation:

- **11_trailing_stop.go** - Automatic trailing stop manager
- **12_position_scaler.go** - Position scaling
- **13_grid_trader.go** - Grid trading for sideways market
- **14_risk_manager.go** - Automatic risk management
- **15_portfolio_rebalancer.go** - Portfolio rebalancing

### examples/demos/helpers/

Helper utilities and tools:

- **17_protobuf_inspector.go** - Interactive protobuf types explorer
- **connection.go** - Connection helper
- **error_helper.go** - Error handling
- **progress_bar.go** - Progress indicator

### examples/demos/usercode/

Sandbox for your custom code:

- **18_usercode.go** - Template file for your own trading logic

**[📖 User Code Sandbox Guide](USERCODE_SANDBOX_GUIDE.md)**

---

## 🎮 Running Examples with main.go

### What is main.go?

`examples/demos/main.go` is the **single entry point** for all GoMT5 demo examples. It works as a dispatcher that:

1. **Runs any example** by number or readable name (alias)
2. **Manages connection** to MT5 via gRPC (creation, validation, closing)
3. **Loads configuration** from `config/config.json` (credentials, server)
4. **Ensures clean exit** (graceful shutdown, connection closing)
5. **Shows interactive menu** for example selection (if run without parameters)

**Benefits of centralized main.go:**

- No need to copy connection code into each example
- All examples use one configuration
- Easy to switch between examples
- Automatic error handling and cleanup

---

### 🟢 Two Ways to Run

#### Method 1: Interactive Menu (recommended for learning)

Run without parameters - a menu with all available examples will appear:

```bash
cd examples/demos
go run main.go
```

You will see:
```
╔═══════════════════════════════════════════════════════════╗
║              GoMT5 - DEMONSTRATION EXAMPLES               ║
╚═══════════════════════════════════════════════════════════╝

Select an example to run:

[LOW-LEVEL] MT5Account - Direct gRPC/Protobuf
  1  - General operations (account, symbols, positions)
  2  - Trading operations (orders, modify, close)
  3  - Streaming methods (real-time ticks, trades)

[MID-LEVEL] MT5Service - Go Wrappers
  4  - Service API demo
  5  - Service streaming demo

[HIGH-LEVEL] MT5Sugar - Convenience API
  6  - Sugar basics (connection, balance, prices)
  7  - Sugar trading (market orders, pending orders)
  8  - Sugar positions (query, modify, close)
  9  - Sugar history (deals, profit calculations)
  10 - Sugar advanced (risk management, symbol info)

[ORCHESTRATORS] Trading Strategies
  11 - Trailing Stop Manager
  12 - Position Scaler (pyramiding/averaging)
  13 - Grid Trader
  14 - Risk Manager
  15 - Portfolio Rebalancer

[PRESETS] Complete Systems
  16 - Adaptive Market Preset

[TOOLS]
  17 - Protobuf Inspector (explore data structures)
  18 - User Code Sandbox (your custom code)

Enter number or alias (or 'x' to exit):
```

#### Method 2: Direct Run (quick access)

If you know what you want to run, use **number** or **alias**:

```bash
cd examples/demos

# By number
go run main.go 1        # Will run example #1
go run main.go 11       # Will run Trailing Stop
go run main.go 18       # Will run User Sandbox

# By alias (more memorable)
go run main.go lowlevel01       # Same as 1
go run main.go trailing         # Same as 11
go run main.go usercode         # Same as 18
go run main.go inspect          # Protobuf Inspector
```

---

### 📋 Full Command Table

#### LOW-LEVEL (MT5Account) — Direct gRPC Calls

| Number | Aliases | File | Description |
|-------|--------|------|----------|
| **1** | `lowlevel01` | `lowlevel/01_general_operations.go` | Information methods: account, symbols, positions, ticks |
| **2** | `lowlevel02` | `lowlevel/02_trading_operations.go` | Trading operations: calculations, validation, orders |
| **3** | `lowlevel03`, `streaming`, `stream` | `lowlevel/03_streaming_methods.go` | Real-time streaming: ticks, deals, positions |

**When to use**: Learning low-level protobuf structures, maximum control.

---

#### MID-LEVEL (MT5Service) — Go Wrappers

| Number | Aliases | File | Description |
|-------|--------|------|----------|
| **4** | `service`, `mid` | `service/04_service_demo.go` | Service API wrappers: account, symbols, positions, orders |
| **5** | `service05`, `servicestreaming`, `service-streaming` | `service/05_service_streaming.go` | Service streaming methods: ticks, deals, profits |

**When to use**: More convenient work with Go types, less boilerplate code.

---

#### HIGH-LEVEL (MT5Sugar) — Simplified API

| Number | Aliases | File | Description |
|-------|--------|------|----------|
| **6** | `sugar06`, `sugarbasics`, `basics` | `sugar/06_sugar_basics.go` | Basics: connection, balance, prices |
| **7** | `sugar07`, `sugartrading`, `trading` | `sugar/07_sugar_trading.go` | Trading: market/pending orders, SL/TP |
| **8** | `sugar08`, `sugarpositions`, `positions` | `sugar/08_sugar_positions.go` | Positions: queries, modification, closing |
| **9** | `sugar09`, `sugarhistory`, `history` | `sugar/09_sugar_history.go` | History: deals over period, profit calculation |
| **10** | `sugar10`, `sugaradvanced`, `advanced` | `sugar/10_sugar_advanced.go` | Advanced: risk management, symbol info |

**When to use**: Quick prototyping, simple scripts, one line of code.

---

#### ORCHESTRATORS — Trading Strategies

| Number | Aliases | File | Description |
|-------|--------|------|----------|
| **11** | `trailing`, `trailingstop` | `orchestrators/11_trailing_stop.go` | Automatic trailing stop manager |
| **12** | `scaler`, `positionscaler` | `orchestrators/12_position_scaler.go` | Position scaling (pyramiding/averaging) |
| **13** | `grid`, `gridtrading` | `orchestrators/13_grid_trader.go` | Grid trading for sideways market |
| **14** | `risk`, `riskmanager` | `orchestrators/14_risk_manager.go` | Automatic risk management |
| **15** | `rebalancer`, `portfolio` | `orchestrators/15_portfolio_rebalancer.go` | Portfolio rebalancing |

**When to use**: Learning complex strategies, combining multiple API methods.

---

#### PRESETS — Ready-Made Systems

| Number | Aliases | File | Description |
|-------|--------|------|----------|
| **16** | `adaptive`, `preset` | `presets/16_AdaptiveOrchestratorPreset.go` | Adaptive trading system (multiple orchestrators) |

**When to use**: Learning complete trading systems.

---

#### TOOLS — Tools

| Number | Aliases | File | Description |
|-------|--------|------|----------|
| **17** | `inspect`, `inspector`, `proto` | `helpers/17_protobuf_inspector.go` | Interactive protobuf structures explorer |
| **18** | `usercode`, `user`, `sandbox`, `custom` | `usercode/18_usercode.go` | Sandbox for your code |

**When to use**:

- **17**: Exploring protobuf message structure
- **18**: Writing your own trading logic

---

### Initial Setup (mandatory!)

**Before running any examples, you must install dependencies:**

```bash
# 1. Navigate to the examples directory
cd examples/demos

# 2. Install all dependencies (only once)
go mod tidy
```

**What `go mod tidy` does:**

- ✅ Downloads all necessary Go modules
- ✅ Adds missing dependencies to `go.mod`
- ✅ Removes unused packages
- ✅ Updates `go.sum` file with checksums
- ✅ Prepares the project for running

**⚠️ Important**: This command needs to be run only ONCE before the first run. After that, all examples will work.

---

### Usage Examples

#### Initial Learning (with interactive menu)

```bash
cd examples/demos

# Run menu (dependencies already installed above)
go run main.go

# Select an example (enter number and press Enter)
# Start with 1, 2, 3 to understand the basics
```

#### Quick Testing of Specific Functions

```bash
cd examples/demos

# Testing Sugar API (simplest)
go run main.go basics          # Balance, prices
go run main.go trading         # Opening orders
go run main.go positions       # Working with positions

# Testing orchestrators
go run main.go trailing        # Trailing stop
go run main.go grid            # Grid trading
go run main.go scaler          # Pyramiding/averaging

# Data exploration
go run main.go inspect         # View protobuf structures

# Working with your code
go run main.go usercode        # Your sandbox
```

#### Learning Sequence (recommended)

```bash
# 1. Start with low level (understanding basics)
go run main.go 1               # How gRPC calls work
go run main.go 2               # How trading operations execute

# 2. Move to mid level (Go wrappers)
go run main.go service         # More convenient API

# 3. Explore high level (simplification)
go run main.go basics          # Sugar API - one line of code
go run main.go trading         # Simple trading operations
go run main.go advanced        # Risk management

# 4. Explore orchestrators (strategies)
go run main.go trailing        # Ready-made strategy
go run main.go grid            # Complex logic

# 5. Write your code
go run main.go usercode        # Your trading logic
```

---


### How main.go Handles Commands

Inside `main.go` there is a large `switch` block that routes commands:

```go
switch arg {
case "1", "lowlevel01":
    return lowlevel.RunGeneral01()      // → lowlevel/01_general_operations.go

case "11", "trailing", "trailingstop":
    return RunOrchestrator_TrailingStop()  // → orchestrators/11_trailing_stop.go

case "18", "usercode", "user", "sandbox", "custom":
    return usercode.RunUserCode()          // → usercode/18_usercode.go
}
```

Each command calls the corresponding function, which:

1. Gets the connection to MT5 (already created by main.go)
2. Executes the demonstration
3. Returns control to main.go
4. main.go automatically closes the connection

**This means**: You do NOT need to worry about connection/disconnection in each example!

---

## 🔧 Development Dependencies

GoMT5 uses the following key packages:

```go
import (
    // gRPC and protobuf
    "google.golang.org/grpc"
    "google.golang.org/protobuf/types/known/emptypb"

    // Protobuf definitions for MT5
    pb "github.com/MetaRPC/GoMT5/package"
)
```

---

## 🎯 Exploring Advanced Features

After mastering the three-tier chain (MT5Account → MT5Service → MT5Sugar), you can explore advanced features:

### Real-Time Streaming

Learn how to work with real-time data streams in Go:

**[📖 gRPC Stream Management](GRPC_STREAM_MANAGEMENT.md)**

Topics covered in the guide:

- Context-based stream management
- Reading from two channels (data and errors)
- Goroutine management and leak prevention
- Automatic reconnection patterns
- Multiple concurrent streams

### Protobuf Inspector

Interactive tool for exploring MT5 API types:

**[📖 Protobuf Inspector Guide](PROTOBUF_INSPECTOR_GUIDE.md)**

Run with:

```bash
cd examples/demos
go run main.go inspect
```

Features:

- Explore 267 protobuf types
- Search types and fields
- View enum values
- Understand message structures

---

## 🛠️ Building Your Own System

After studying all examples and understanding the architecture, you can start building your own trading system:

**Sandbox location:** `examples/demos/usercode/18_usercode.go`

This file is prepared as a starter template for your custom code. Before you start coding here, make sure you've studied:

1. The three-tier API (MT5Account → MT5Service → MT5Sugar)
2. At least a few helper examples to understand the patterns
3. Error handling patterns from examples
4. Stream management, if you need real-time data

**[📖 User Code Sandbox Guide](USERCODE_SANDBOX_GUIDE.md)**

**Quick start:**

```bash
cd examples/demos
go run main.go 18
# Edit usercode/18_usercode.go and uncomment examples
```

---


### Error Handling

GoMT5 has comprehensive error handling:

**File:** `package/Helpers/errors.go`

Features:

- `ErrNotConnected` - sentinel error for connection problems
- `ApiError` - wraps protobuf errors with Go methods
- Trading operation return codes - constants for all MT5 return codes
- Helper functions: `IsRetCodeSuccess()`, `IsRetCodeRetryable()`

### Return Codes Reference

Understanding return code meanings when executing trading operations:

**[📖 Return Codes Reference](RETURN_CODES_REFERENCE.md)**

Common codes:

- `10009` - Success (market order)
- `10008` - Success (pending order)
- `10019` - Insufficient margin
- `10016` - Invalid stops (SL/TP too close)
- `10018` - Market closed

### gRPC Stream Management

If you're working with real-time streaming (ticks, events):

**[📖 gRPC Stream Management](GRPC_STREAM_MANAGEMENT.md)**

Critical patterns:

- Always use `context.WithCancel()` or `context.WithTimeout()`
- Always `defer cancel()`
- Always read from BOTH channels (data and errors)
- Use `sync.WaitGroup` for multiple streams

---

## 📖 Complete Documentation

### Core API Documentation

- **[MT5Account Master Overview](../MT5Account/MT5Account.Master.Overview.md)** - Complete low-level API reference
- **[MT5Service API Overview](../MT5Service/MT5Service.Overview.md)** - Mid-level wrappers API
- **[MT5Sugar API Overview](../MT5Sugar/MT5Sugar.API_Overview.md)** - High-level Sugar API

### Guides and Tutorials

- **[User Code Sandbox Guide](USERCODE_SANDBOX_GUIDE.md)** - Quick start for custom code
- **[gRPC Stream Management](GRPC_STREAM_MANAGEMENT.md)** - Working with real-time streams
- **[Protobuf Inspector Guide](PROTOBUF_INSPECTOR_GUIDE.md)** - Interactive type explorer
- **[Return Codes Reference](RETURN_CODES_REFERENCE.md)** - Trading operation return codes

### Code Examples

All examples are located in `examples/demos/` numbered 01-17.

---

## 💬 Support and Community

If you encounter problems in the code or have questions about the gateway:

**Support repository:** [https://github.com/Moongoord/MetaRPC-Gateway-Support](https://github.com/Moongoord/MetaRPC-Gateway-Support)

This is a discussion repository where you can:

- Ask questions about the gateway
- Report issues
- Share your experience
- Get help from the community

**Note:** This repository is currently under development and will be available to everyone soon.

---

## 🎓 Conclusion

**The MetaRPC team** strives to create favorable conditions for learning fundamental trading principles and building algorithmic trading systems with Go.

We believe that with diligence and desire to learn, you will be able to master everything - from low-level protocol communication to complex trading systems.

**Your journey starts here:**

1. Set up the environment (above)
2. Create or configure your MT5 demo account ([MT5 for Beginners](MT5_For_Beginners/))
3. Choose a learning path (bottom-up or top-down)
4. Run your first example: `go run main.go 1`
5. Study the code, experiment and create

**Good luck on your algorithmic trading journey with Go!**

> "The foundation of success in algorithmic trading is not only understanding markets, but also understanding the code that interacts with them. Master both, and you will have unlimited possibilities."
>
> — MetaRPC Team

---

## 🚀 Quick Start Checklist

- [ ] Install Go 1.21+ (`go version`)
- [ ] Clone the repository
- [ ] Install VS Code with Go extension
- [ ] Create `config.json` with your MT5 credentials
- [ ] Run `cd examples/demos && go mod tidy`
- [ ] Run first example: `go run main.go 1`
- [ ] Explore the three API tiers
- [ ] Try protobuf inspector: `go run main.go inspect`
- [ ] Read the streaming guide if you need real-time data
- [ ] Start building in `18_usercode.go`

---

### ⚡ Want Results Here and Now?

**Don't want to clone the entire repository and study the architecture?**

**Need a quick start in your own directory?**

If you can't wait to try the gRPC gateway capabilities as soon as possible and see your first working code in 10 minutes:

📖 **[Your First Project - Project from Scratch](Your_First_Project.md)**

**What you'll get:**

- ✅ Create your project from scratch (without cloning the repository)
- ✅ Connect only necessary dependencies
- ✅ Write your first low-level method in 10 minutes
- ✅ Get account balance and see the result
- ✅ Start working with the gateway in your code immediately

**Difference in approaches:**

| This Document (GETTING_STARTED) | Your First Project |
|----------------------------------|-------------------|
| Clone ready-made repository | Create project from scratch |
| Study examples and architecture | Write working code immediately |
| Full SDK immersion | Minimal quick start |
| For deep understanding | For instant result |

💡 **Recommendation:** After a quick start with Your First Project, return to this document to learn the complete architecture and all SDK capabilities.

---


## Go-Specific Tips

**Effective Go practices used in this project:**

- **Context everywhere** - All methods accept `context.Context` for timeout management
- **Error handling** - Always check errors, Go has no exceptions
- **Channels for streams** - Two-channel pattern (data, errors)
- **Goroutines** - Background processing for streams
- **Defer for cleanup** - `defer cancel()` ensures resource cleanup
- **Interfaces** - Clean abstractions between tiers
- **Composition** - Sugar embeds Service, Service uses Account

**Resources for Go beginners:**

- [Tour of Go](https://go.dev/tour/) - Official interactive tutorial
- [Effective Go](https://go.dev/doc/effective_go) - Best practices guide
- [Go by Example](https://gobyexample.com/) - Learning by examples

---
