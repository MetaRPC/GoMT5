# GoMT5 Project Map

> Complete reference to project structure. Shows what is located where, what is user-facing vs internal, and how components are connected.

---

## 🗺️ Project Overview

```
GoMT5/
├── 📦 package/ - Independent module (portable)
│   ├── Helpers/MT5Account.go (Layer 1 - Foundation)
│   ├── Helpers/errors.go (Error handling & trade return codes)
│   ├── Proto definitions (*.pb.go)
│   └── gRPC stubs (*_grpc.pb.go)
│
├── 📦 examples/mt5/ - High-level API layers
│   ├── MT5Service.go (Layer 2 - Wrappers)
│   └── MT5Sugar.go (Layer 3 - Convenience)
│
├── 🎯 User Code (Orchestrators, Presets, Examples)
├── 📚 Documentation
└── ⚙️ Configuration and build

External dependencies:
└── 🔌 gRPC & Proto (Go modules)
```

---

## 📦 Core API (Three-layer architecture)

**What:** Three-tier architecture for MT5 trading automation.

**User interaction:** Import and use, but typically don't modify.

```
package/Helpers/
├── MT5Account.go              ← LAYER 1: Low-level gRPC ⭐ FOUNDATION
│   └── Direct gRPC calls to MT5 terminal
│   └── Connection management with retry logic
│   └── Proto Request/Response handling
│   └── Built-in connection resilience
│   └── Independent Go module (portable)
│
└── errors.go                  ← Error handling & trade result codes
    └── ErrNotConnected sentinel error
    └── ApiError wrapper (3-level: API/MQL/Trade)
    └── Trade return code constants & helpers
    └── Centralized error inspection methods

examples/mt5/
├── MT5Service.go              ← LAYER 2: Wrapper methods
│   └── Simplified signatures (no proto objects)
│   └── Type conversion (proto → Go primitives)
│   └── Direct data return
│   └── Typed result structs (AccountSummary, SymbolTick, OrderResult, ...)
│
└── MT5Sugar.go                ← LAYER 3: Convenience layer ⭐
    └── One-line trading (BuyMarket, BuyLimit, BuyMarketWithSLTP, ...)
    └── Risk management (CalculatePositionSize, CanOpenPosition, GetMaxLotSize)
    └── Distance-based SL/TP (BuyMarketWithPips, SellMarketWithPips, CalculateSLTP)
    └── Batch operations (CloseAllPositions, CloseAllBySymbol)
    └── Snapshots (GetAccountInfo, GetSymbolInfo, GetDailyStats)
    └── No automatic normalization - volumes and prices are sent as given

go.mod / go.sum                ← Module dependencies
```

**Architecture flow:**
```
MT5Sugar → uses → MT5Service → uses → MT5Account → gRPC → MT5 Terminal
       ↓                ↓                    ↓
examples/mt5/    examples/mt5/      package/Helpers/
```

**💡 Creating Your Own Project:**

For your own standalone project using GoMT5, you only need the `package` module:

```go
import (
    pb "github.com/MetaRPC/GoMT5/package"           // protobuf messages and enums
    mt5 "github.com/MetaRPC/GoMT5/package/Helpers"  // MT5Account, ApiError, return-code helpers
)
```

The `package` module contains **everything you need to start**:

- ✅ All protobuf definitions (proto-generated code)
- ✅ gRPC stubs and service contracts
- ✅ MT5Account (Layer 1 - Foundation)
- ✅ Independent Go module (can be used without examples/)

This makes the package **portable** and easy to integrate into any Go project!

**User decision:**

- **Building your own app:** Import `package` and use MT5Account directly
- **Learning/Examples:** Use the full GoMT5 repo with all 3 layers (MT5Service/MT5Sugar are `github.com/MetaRPC/GoMT5/examples/mt5` in the clone)
- **95% of demo cases:** Start with `MT5Sugar` (highest level, easiest)
- **Need wrappers:** Move to `MT5Service` (Go types, full control over request parameters)
- **Need raw proto:** Move to `MT5Account` (full control)

**Documentation:**

- [MT5Account API Reference](../API_Reference/MT5Account.md) ⭐ **FOUNDATION OF EVERYTHING**
- [MT5Service API Reference](../API_Reference/MT5Service.md)
- [MT5Sugar API Reference](../API_Reference/MT5Sugar.md)

---

## 🎯 User Code (Your Trading Strategies)

### Orchestrators (examples/demos/orchestrators/)

**What:** Ready-made trading strategy implementations.

```
examples/demos/orchestrators/
├── 13_grid_trader.go             ← Grid trading (sideways markets)
├── 11_trailing_stop.go           ← Trailing stop (price following)
├── 12_position_scaler.go         ← Position scaling
├── 14_risk_manager.go            ← Risk manager
├── 15_portfolio_rebalancer.go    ← Portfolio rebalancing
└── orchestrators.go              ← Base interface/functionality
```

**Purpose:** Educational examples showing complete strategy workflows:

- Entry logic (risk-based volume where applicable)
- Position monitoring with progress bars
- Exit management and cleanup
- Performance tracking (balance, equity, P/L)
- Configurable parameters via config structs (e.g. `TrailingStopConfig`, `GridTraderConfig`)

**How to use:**

1. Study existing orchestrators
2. Copy one as a template
3. Modify for your strategy
4. Test on demo account

**How to run:**
```bash
go run examples/demos/main.go grid         # Grid Trader
go run examples/demos/main.go trailing     # Trailing Stop
go run examples/demos/main.go scaler       # Position Scaler
go run examples/demos/main.go risk         # Risk Manager
go run examples/demos/main.go portfolio    # Portfolio Rebalancer
```

**Documentation:**
- [Grid Trader](../Orchestrators/13_Grid_trader.md)
- [Trailing Stop](../Orchestrators/11_Trailing_stop.md)
- [Position Scaler](../Orchestrators/12_Position_scaler.md)
- [Risk Manager](../Orchestrators/14_Risk_manager.md)
- [Portfolio Rebalancer](../Orchestrators/15_Portfolio_rebalancer.md)

---

### Presets (examples/demos/presets/)

**What:** Combinations of multiple orchestrators with adaptive logic based on market analysis.

**User interaction:** ✅ **Advanced usage** - combine multiple strategies.

```
examples/demos/presets/
└── 16_AdaptiveOrchestratorPreset.go    ← Intelligent multi-strategy
```

**Purpose:** Demonstrate how to:

- Combine multiple orchestrators
- Adaptive decision making (volatility → strategy)
- Market condition analysis (simplified demo)
- Multi-phase trading sessions
- Performance tracking by phases

**How to run:**
```bash
go run examples/demos/main.go preset       # Adaptive Market Preset
go run examples/demos/main.go adaptive     # Same thing
```

---

### Examples (examples/demos/)

**What:** Runnable examples demonstrating API usage at different layers.

**User interaction:** ✅ **Learning materials** - run to understand the API.

```
examples/demos/
├── lowlevel/                          ← MT5Account examples (proto level) ⭐ FOUNDATION
│   ├── 01_general_operations.go       ← General operations (connection, account, symbols)
│   ├── 02_trading_operations.go       ← Trading operations (orders, positions)
│   └── 03_streaming_methods.go        ← Streaming methods (real-time subscriptions)
│
├── service/                           ← MT5Service examples (wrapper level)
│   ├── 04_service_demo.go             ← Service API demonstration
│   └── 05_service_streaming.go        ← Service streaming methods
│
├── sugar/                             ← MT5Sugar examples (convenience level)
│   ├── 06_sugar_basics.go             ← Sugar API basics (balance, prices)
│   ├── 07_sugar_trading.go            ← Trading (market/limit orders)
│   ├── 08_sugar_positions.go          ← Position management
│   ├── 09_sugar_history.go            ← History and statistics
│   └── 10_sugar_advanced.go           ← Advanced Sugar capabilities
│
└── usercode/                          ← User code sandbox
    └── 18_usercode.go                 ← Your custom strategies
```

**How to run:**
```bash
# Low-level examples (MT5Account - FOUNDATION OF EVERYTHING)
go run examples/demos/main.go lowlevel01   # General operations
go run examples/demos/main.go lowlevel02   # Trading operations
go run examples/demos/main.go lowlevel03   # Streaming methods

# Service examples (MT5Service - wrappers)
go run examples/demos/main.go service      # Service API demo
go run examples/demos/main.go service05    # Service streaming methods

# Sugar examples (MT5Sugar - convenience API)
go run examples/demos/main.go sugar06      # Sugar basics
go run examples/demos/main.go sugar07      # Sugar trading
go run examples/demos/main.go sugar08      # Sugar positions
go run examples/demos/main.go sugar09      # Sugar history
go run examples/demos/main.go sugar10      # Advanced Sugar

# UserCode (your code)
go run examples/demos/main.go usercode     # Custom strategies
```

---

### main.go (examples/demos/)

**What:** Main entry point that routes `go run` commands to corresponding examples/orchestrators/presets.

**User interaction:** 📋 **Runner + Documentation** - launches everything.

```
main.go
├── main()                              ← Entry point, parses arguments (--api-key) or shows the menu
├── showMenu()                          ← Interactive menu [1]-[18], [0] exit
├── executeCommand()                    ← Maps numbers and aliases to runners
├── RunOrchestrator_TrailingStop() ...  ← Configure and launch orchestrators (one function each)
├── RunOrchestrator_AdaptivePreset()    ← Launches the preset
└── Documentation in header             ← Command reference
```

**How it works:**

```
go run examples/demos/main.go grid
    ↓
main()  // os.Args[1] = "grid"
    ↓
executeCommand("grid")
    ↓
RunOrchestrator_Grid()
    ↓
orchestrators.NewGridTrader(sugar, config).Start()
```

**Purpose:**

- Single entry point for all runnable code
- Command routing by menu number or alias (grid, trailing, preset, etc.)
- Lists the valid commands for an unknown command

**Available commands:** See header comment in `main.go` for full list.

---

### Helpers (examples/demos/helpers/)

**What:** Utilities for examples and orchestrators.

```
examples/demos/helpers/
├── connection.go                 ← MT5 connection setup
├── error_helper.go               ← Error handling and return codes
├── progress_bar.go               ← Visual progress bars
└── 17_protobuf_inspector.go      ← Protobuf structure inspector (runnable)
```

**Connection (connection.go):**
```go
// Load config.json / env vars, create MT5Account and ConnectEx by server name
account, cfg, err := helpers.CreateAndConnectAccount()
defer helpers.DisconnectAccount(account)
service := mt5.NewMT5Service(account)

// MT5Sugar creates its own account from the credentials
sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
err = sugar.QuickConnect(cfg.MtCluster)
```

**Progress bar (progress_bar.go):**
```go
// Visual countdown during orchestrator operation
helpers.WaitWithProgressBar(60, "Monitoring positions", ctx)
```

**Errors (error_helper.go):**
```go
// Check return codes and handle errors
if !helpers.CheckRetCode(result.ReturnedCode, "Order placement") {
    // prints the code and its description
}
helpers.PrintIfError(err, "Order placement failed") // prints ApiError details when err != nil
```

**Protobuf inspector (17_protobuf_inspector.go):**
```go
// Interactive explorer of protobuf types, fields and enums (menu item 17)
helpers.RunProtobufInspector()
```

---

## 📚 Documentation (docs/)

**What:** Complete API and strategy documentation.

**User interaction:** 📖 **Read first!** Comprehensive reference.

```
docs/                                  (mkdocs.yml - the MkDocs configuration - is in the repo root)
├── index.md                           ← Home page - project introduction
│
├── styles/custom.css                  ← Custom theme (ocean aurora)
├── javascripts/ux.js                  ← Interactive features
│
├── All_Guides/                        ← Guides
│   ├── MT5_For_Beginners.md           ← Demo account registration
│   ├── GETTING_STARTED.md             ← ⭐ Start here! Setup and first steps
│   ├── Your_First_Project.md          ← Your first project
│   ├── GLOSSARY.md                    ← ⭐ Terms and definitions
│   ├── GRPC_STREAM_MANAGEMENT.md      ← Managing streaming subscriptions
│   ├── RETURN_CODES_REFERENCE.md      ← Proto return code reference
│   ├── ENUMS_USAGE_REFERENCE.md       ← Which methods use which enums
│   ├── PROTOBUF_INSPECTOR_GUIDE.md    ← Protobuf inspector tool
│   ├── USERCODE_SANDBOX_GUIDE.md      ← How to write custom strategies
│   └── PROJECT_MAP.md                 ← ⭐ This file - complete structure
│
├── API_Reference/                     ← Concise API documentation
│   │                                     (slightly enhanced for better navigation)
│   ├── MT5Account.md                  ← ⭐ Layer 1 API (foundation of everything) → from package/Helpers/MT5Account.go
│   ├── MT5Service.md                  ← Layer 2 API → from examples/mt5/MT5Service.go
│   └── MT5Sugar.md                    ← Layer 3 API → from examples/mt5/MT5Sugar.go
│
├── MT5Account/                        ← ⭐ FOUNDATION OF EVERYTHING - Detailed Layer 1 documentation
│   ├── MT5Account.Master.Overview.md  ← ⭐ Complete API reference
│   │
│   ├── 1. Account_information/        ← Account methods (~4 files)
│   │   ├── AccountInfoDouble.md       ← Get account double parameters
│   │   ├── AccountSummary.md          ← Complete account summary
│   │   └── ...                        ← And other account methods
│   │   └── 💡 Each example linked with examples/demos/lowlevel
│   │
│   ├── 2. Symbol_information/         ← Symbol/market data methods (~9 files)
│   │   ├── SymbolInfoTick.md          ← Current symbol tick
│   │   ├── SymbolInfoDouble.md        ← Symbol double parameters
│   │   ├── SymbolsTotal.md            ← Total symbols count
│   │   └── ...                        ← And other symbol methods
│   │   └── 💡 Examples in examples/demos/lowlevel
│   │
│   ├── 3. Position_Orders_Information/ ← Position/order methods (~5 files)
│   │   ├── OpenedOrders.md            ← List of open orders
│   │   ├── PositionsTotal.md          ← Total positions count
│   │   └── ...                        ← And other position methods
│   │   └── 💡 Examples in examples/demos/lowlevel
│   │
│   ├── 4. Trading_Operations/         ← Trading operation methods (~6 files)
│   │   ├── OrderSend.md               ← Send order (main method)
│   │   ├── OrderCheck.md              ← Check order before sending
│   │   ├── OrderCalcMargin.md         ← Calculate margin
│   │   ├── OrderCalcProfit.md         ← Calculate profit
│   │   └── ...                        ← And other trading methods
│   │   └── 💡 Examples in examples/demos/lowlevel/02_trading_operations.go
│   │
│   ├── 5. Market_Depth(DOM)/          ← Market depth methods (~3 files)
│   │   ├── MarketBookAdd.md           ← Subscribe to market depth
│   │   ├── MarketBookGet.md           ← Get market depth data
│   │   └── ...                        ← And other DOM methods
│   │
│   ├── 6. Additional_Methods/         ← Additional methods (~5 files)
│   │   ├── SymbolInfoMarginRate.md    ← Symbol margin rates
│   │   ├── SymbolInfoSessionQuote.md  ← Quote trading sessions
│   │   └── ...                        ← And other auxiliary methods
│   │
│   ├── 7. Streaming_Methods/          ← Real-time subscription methods
│   │   ├── OnSymbolTick.md            ← Subscribe to symbol ticks
│   │   ├── OnTrade.md                 ← Subscribe to trade events
│   │   ├── OnPositionProfit.md        ← Subscribe to profit changes
│   │   └── ...                        ← And other streaming methods
│   │   └── 💡 Stream management examples in All_Guides/GRPC_STREAM_MANAGEMENT
│   │
│   └── HOW_IT_WORK/                   ← Step-by-step explanations (*_HOW.md), one folder per group
│
├── MT5Service/                        ← Service level method documentation
│   ├── MT5Service.Overview.md          ← ⭐ Complete Service API reference
│   ├── 1. Account_Methods.md          ← Account helper methods
│   ├── 2. Symbol_Methods.md           ← Symbol helper methods
│   ├── 3. Position_Orders_Methods.md  ← Position/order helper methods
│   ├── 4. Trading_Methods.md          ← Trading helper methods
│   ├── 5. MarketDepth_Methods.md      ← Market depth helper methods
│   └── 6. Streaming_Methods.md        ← Streaming helper methods
│
├── MT5Sugar/                          ← Sugar level method documentation
│   ├── MT5Sugar.API_Overview.md        ← ⭐ Complete Sugar API reference
│   │
│   ├── 1. Connection/                  ← Connection methods (~3 files)
│   │   ├── QuickConnect.md            ← Quick connection
│   │   ├── IsConnected.md             ← Check connection
│   │   └── Ping.md                    ← Connection test
│   │
│   ├── 2. Balance_Margin/              ← Balance and margin (~6 files)
│   │   ├── GetBalance.md              ← Get balance
│   │   ├── GetEquity.md               ← Get equity
│   │   ├── GetFreeMargin.md           ← Free margin
│   │   └── ...                        ← And other balance methods
│   │
│   ├── 3. Prices_Quotes/               ← Prices and quotes (~5 files)
│   │   ├── GetBid.md                  ← Get Bid
│   │   ├── GetAsk.md                  ← Get Ask
│   │   ├── GetSpread.md               ← Get spread
│   │   └── ...                        ← And other price methods
│   │
│   ├── 4. Simple_Trading/              ← Simple trading (~6 files)
│   │   ├── BuyMarket.md               ← Buy at market
│   │   ├── SellMarket.md              ← Sell at market
│   │   ├── BuyLimit.md                ← Buy Limit order
│   │   └── ...                        ← And other simple orders
│   │
│   ├── 5. Trading_SLTP/                ← Trading with SL/TP (~4 files)
│   │   ├── BuyMarketWithSLTP.md       ← Buy with SL/TP
│   │   ├── SellMarketWithSLTP.md      ← Sell with SL/TP
│   │   └── ...                        ← And other orders with SL/TP
│   │
│   ├── 6. Position_Management/         ← Position management (~7 files)
│   │   ├── ClosePosition.md           ← Close position
│   │   ├── CloseAllPositions.md       ← Close all positions
│   │   ├── ModifyPositionSLTP.md      ← Modify SL/TP
│   │   └── ...                        ← And other management methods
│   │
│   ├── 7. Position_Information/        ← Position information (~2 files)
│   │   ├── HasOpenPosition.md         ← Has open position
│   │   └── CountOpenPositions.md      ← Count positions
│   │
│   ├── 8. History_Statistics/          ← History and statistics (~9 files)
│   │   ├── GetDealsToday.md           ← Deals today
│   │   ├── GetProfitThisWeek.md       ← Profit this week
│   │   ├── GetDealsDateRange.md       ← Deals for period
│   │   └── ...                        ← And other history methods
│   │
│   ├── 9. Symbol_Information/          ← Symbol information (~5 files)
│   │   ├── GetSymbolInfo.md           ← Complete symbol information
│   │   ├── GetAllSymbols.md           ← All available symbols
│   │   └── ...                        ← And other symbol methods
│   │
│   ├── 10. Risk_Management/            ← Risk management (~4 files)
│   │   ├── CalculatePositionSize.md   ← Calculate position size
│   │   ├── CanOpenPosition.md         ← Can open position
│   │   └── ...                        ← And other risk methods
│   │
│   ├── 11. Trading_Helpers/            ← Trading helpers (~3 files)
│   │   ├── BuyMarketWithPips.md       ← Buy with SL/TP as distances (points)
│   │   ├── CalculateSLTP.md           ← Calculate SL/TP
│   │   └── ...                        ← And other helpers
│   │
│   └── 12. Account_Information/        ← Account information (~3 files)
│       ├── GetAccountInfo.md          ← Account information in one call
│       ├── AccountInfo.md             ← AccountInfo structure
│       └── DailyStats.md              ← DailyStats structure
│
└── Orchestrators/                     ← Strategy documentation
    ├── 13_Grid_trader.md              ← Grid trading
    ├── 11_Trailing_stop.md            ← Trailing stop
    ├── 12_Position_scaler.md          ← Position scaling
    ├── 14_Risk_manager.md             ← Risk manager
    ├── 15_Portfolio_rebalancer.md     ← Portfolio rebalancing
    ├── 16_AdaptiveOrchestratorPreset.md ← Adaptive preset
    └── Strategies.Master.Overview.md  ← Complete strategy overview
```

**Structure:**

- Each method has its own `.md` file with examples
- Overview files (`*.Master.Overview.md`) provide navigation
- `MT5Account/HOW_IT_WORK/**/*_HOW.md` files explain algorithms step by step
- Links between related methods
- Usage examples in each file

**⭐ Important about MT5Account:**
- **FOUNDATION OF EVERYTHING** - all methods here are the foundation
- Each documentation example is linked with real code in `examples/demos/lowlevel/`
- Understanding MT5Account is critical for effective use of MT5Service and MT5Sugar

---

## 🔌 gRPC & Proto (Go modules)

**What:** Protocol Buffer and gRPC libraries for communication with MT5 terminal.

**User interaction:** 📋 **Reference only** - managed via Go modules.

**Key dependencies:**

- `google.golang.org/grpc` - gRPC client
- `google.golang.org/protobuf` - Protocol Buffers runtime
- `github.com/MetaRPC/GoMT5/package` - MT5 Proto definitions (independent module)

**Package structure:**

```
package/
├── Helpers/
│   ├── MT5Account.go          ← Layer 1 implementation
│   └── errors.go              ← ApiError, ErrNotConnected, trade return codes
├── *.pb.go                    ← Generated protobuf code
├── *_grpc.pb.go               ← Generated gRPC stubs
├── go.mod                     ← Independent module
└── go.sum                     ← Module dependencies
```

**How it works:**

1. `package/` is an independent Go module
2. Contains both proto-generated code and MT5Account implementation
3. Can be imported separately: `github.com/MetaRPC/GoMT5/package`
4. MT5Service and MT5Sugar import from package module
5. All layers use proto-generated types from package

**Proto-generated types:**

- Proto package `mt5_term_api` (Go package imported as `pb`) - Trading API types
- Request/Response message types
- Enum definitions
- Service contracts

**Purpose:**

- Define gRPC service contracts
- Type-safe communication with MT5 terminal
- Used by MT5Account layer
- Hidden by MT5Service and MT5Sugar layers

---

## 📊 Component Interaction Diagram

```
YOUR CODE (User)
  ├─ Orchestrators (strategy implementations)
  ├─ Presets (multi-strategies)
  └─ Examples (learning materials)
                  │
                  │ uses
                  ↓
MT5Sugar (Layer 3 - Convenience)
  ├─ One-line operations
  ├─ Risk management
  ├─ Distance-based SL/TP helpers
  └─ Batch operations
                  │
                  │ uses
                  ↓
MT5Service (Layer 2 - Wrappers)
  ├─ Direct data return
  ├─ Type conversion
  └─ Simplified signatures
                  │
                  │ uses
                  ↓
MT5Account (Layer 1 - Low level) ⭐ FOUNDATION
  📍 Location: package/Helpers/MT5Account.go
  ├─ Proto Request/Response
  ├─ gRPC communication
  ├─ Connection management
  ├─ Auto-reconnection
  └─ Independent Go module (portable)
                  │
                  │ gRPC
                  ↓
MT5 Gateway (mt5term) or MT5 Terminal
  └─ MetaTrader 5 with gRPC server
```

---

## 🔍 File Naming Conventions

### Core API (Multi-location)

**Layer 1 (Foundation):**
- `package/Helpers/MT5Account.go` - Low-level gRPC (independent module)

**Layers 2-3 (High-level wrappers):**
- `examples/mt5/MT5Service.go` - Wrapper methods
- `examples/mt5/MT5Sugar.go` - Convenience API
- `go.mod / go.sum` - Dependencies

### User Code (examples/demos/)
- `NN_name.go` - Numbered demos, matching the menu items (e.g. `orchestrators/13_grid_trader.go`)
- `presets/16_AdaptiveOrchestratorPreset.go` - Multi-strategy preset
- `main.go` - Entry point and command router
- `helpers/*.go` - Utilities (`connection.go`, `error_helper.go`, `progress_bar.go`)

### Documentation (docs/)
- `*.Master.Overview.md` - Complete category overviews
- `*.Overview.md` - Section overviews
- `MethodName.md` - Individual method documentation
- `*_HOW.md` - Algorithm explanations (in `MT5Account/HOW_IT_WORK/`)

---

## 📂 What to Modify vs What to Leave Alone

### ✅ MODIFY (User Code)

**Recommended starting point:**
```
examples/demos/usercode/18_usercode.go  ← ⭐ SANDBOX - start writing your code here!
                                           All 3 API levels already initialized and ready.
                                           Run: go run main.go 18
```

**Other files for modification:**
```
examples/demos/orchestrators/  ← Copy and customize for your strategies
examples/demos/presets/        ← Create your multi-strategies
examples/demos/lowlevel/       ← Add your low-level examples
examples/demos/service/        ← Add your service examples
examples/demos/sugar/          ← Add your sugar examples
examples/demos/usercode/       ← Create your custom files alongside 18_usercode.go
examples/demos/config/         ← Configure for your MT5 terminal/gateway
examples/demos/main.go         ← Add new command routing if needed
README.md                      ← Update with your changes
```

### 📖 READ (Core API)

```
package/Helpers/MT5Account.go  ← Use but don't modify (import and call) ⭐ FOUNDATION
examples/mt5/MT5Service.go     ← Use but don't modify
examples/mt5/MT5Sugar.go       ← Use but don't modify
docs/                          ← Reference documentation
```

### 🔒 LEAVE ALONE (Generated/Build)

```
.vscode/                       ← VS Code settings
go.work.sum                    ← Go workspace (auto-generated)
docs/site/                     ← Built documentation (auto-generated by MkDocs)
docs/styles/                   ← Documentation theme (don't change without understanding)
docs/javascripts/              ← Documentation scripts (don't change without understanding)
```

---

## 🎯 Project Philosophy

**Goal:** Make MT5 trading automation accessible through progressive complexity.

**Three-tier design:**

1. **Low level (MT5Account):** Full control, proto/gRPC ⭐ **FOUNDATION OF EVERYTHING**
2. **Wrappers (MT5Service):** Simplified method calls
3. **Convenience (MT5Sugar):** One-liners, batteries included

**User code:**

- **Orchestrators:** Ready-made strategy templates
- **Presets:** Adaptive multi-strategies
- **Examples:** Learning materials at all levels

---

## 🛠️ Troubleshooting

### Build Issues

```bash
# Clean and rebuild (from examples/demos)
go clean
go build ./...

# Restore modules
go mod download

# Check Go version
go version   # Should be 1.23 or higher (go.mod: go 1.23.6)
```

### Runtime Issues

```
1. Always test on demo account first
2. Check return codes (10009 = success, 10031 = no connection with the trade server) and ApiError.MqlErrorTradeIntCode() for refused requests
3. Monitor console output for errors
4. Use retry logic for intermittent issues
5. Verify broker allows your strategy type (hedging, etc.)
```

---

## 📈 Performance Considerations

### Connection Management
- Single gRPC connection shared across operations
- Built-in auto-reconnection handles temporary failures
- Retry logic with exponential backoff and jitter (0.5s → 1s → 2s → 4s, capped at 5s) for calls; streams restart after ~0.5s

### Rate Limiting
- 3-second delays between order placements (Sugar demos)
- Gateway may enforce additional rate limits
- Adjust delays based on broker requirements

### Resource Usage
- Blocking calls bounded by context timeouts; goroutines and channels for streams
- Context for graceful shutdown
- Proper cleanup in defer blocks

---

## 📝 Best Practices

### Code Organization
```
✅ DO: Separate concerns (analysis, execution, monitoring)
✅ DO: Use context for lifecycle management
✅ DO: Add comprehensive error handling
✅ DO: Document your strategy logic clearly
✅ DO: Use helpers.WaitWithProgressBar for long operations

❌ DON'T: Mix strategy logic with API calls
❌ DON'T: Use time.Sleep (use time.After with context)
❌ DON'T: Ignore return codes
❌ DON'T: Test on live accounts without extensive demo testing
```

### Strategy Development
```
✅ DO: Start with existing orchestrator as template
✅ DO: Test each component separately
✅ DO: Log all trading decisions and outcomes
✅ DO: Use demo accounts for development
✅ DO: Implement proper risk management

❌ DON'T: Over-optimize on limited data
❌ DON'T: Ignore edge cases and failures
❌ DON'T: Use fixed lots without risk calculation
❌ DON'T: Deploy without backtesting and forward testing
```

---

> 💡 **Remember:** This is an educational project. All orchestrators and presets are demonstration examples, not production-ready trading systems. Always test on demo accounts, thoroughly understand the code, and implement proper risk management before considering live trading.

---

"Trade safely, code cleanly, and may your gRPC connections always be stable."
