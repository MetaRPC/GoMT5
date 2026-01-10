# GoMT5 Glossary

> Project-specific terms and concepts. This glossary covers GoMT5 architecture, components, and trading automation terminology used throughout the codebase.

---

## 🏗️ Architectural Terms

### Three-Tier Architecture

Core design pattern of GoMT5 with three abstraction levels:

- **Tier 1 (MT5Account):** Low-level gRPC communication with MT5 terminal
- **Tier 2 (MT5Service):** Wrapper methods with simplified signatures
- **Tier 3 (MT5Sugar):** High-level convenience methods with auto-normalization

---

### MT5Account
**Tier 1 - Low-Level API**

Base layer providing direct access to MT5 terminal via gRPC protocol.

**Key characteristics:**

- Raw gRPC calls to MT5 terminal
- Built-in connection resilience with automatic reconnection
- Works with protobuf Request/Response objects
- Full control, maximum complexity
- **40+ methods** for all MT5 operations
- Context-based timeout management
- Two-channel streaming (dataChan, errChan)

**When to use:** Custom integrations, need proto-level control, building custom wrappers.

**Location:** `examples/mt5/MT5Account.go` (2300+ lines)

**Documentation:** [MT5Account.Master.Overview.md](../MT5Account/MT5Account.Master.Overview.md)

**Key patterns:**
```go
// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Protobuf request/response
req := &pb.AccountSummaryRequest{}
reply, err := account.AccountSummary(ctx, req)
```

---

### MT5Service
**Tier 2 - Wrappers API**

Middle layer providing simplified method signatures without proto complexity.

**Key characteristics:**

- Direct data return (no proto objects in signatures)
- Type conversion (proto → Go primitives/structs)
- Simplified method names
- No auto-normalization (you control precision)
- **50+ methods** for common scenarios
- Helper functions for common patterns

**When to use:** Need wrappers but not auto-normalization, building custom strategies.

**Location:** `examples/mt5/MT5Service.go` (1160+ lines)

**Documentation:** [MT5Service.Overview.md](../MT5Service/MT5Service.Overview.md)

**Example:**
```go
// Service level - clean Go types
accountInfo, err := service.GetAccountSummary(ctx)
fmt.Printf("Balance: %.2f\n", accountInfo.Balance)
```

---

### MT5Sugar
**Tier 3 - Convenience API** ⭐

High-level API for common trading operations.

**Key characteristics:**

- Auto-normalization of volumes and prices
- Risk management helpers
- Batch operations (CloseAllPositions, etc.)
- Pip-based helpers (BuyMarketWithPips, etc.)
- One-line operations
- Simplest API, handles edge cases
- **Best starting point** for 95% of use cases

**When to use:** Rapid development, prototyping, simple strategies.

**Location:** `examples/mt5/MT5Sugar.go` (2200+ lines)

**Documentation:** [MT5Sugar.API_Overview.md](../MT5Sugar/MT5Sugar.API_Overview.md)

**Examples:**
```go
// One line to get balance
balance, _ := sugar.GetBalance()

// One line to open position
ticket, _ := sugar.BuyMarket("EURUSD", 0.01)

// Risk-based position sizing (2% with 50 pip SL)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

---

## 🔧 Technical Concepts

### Auto-Normalization
Automatic adjustment of trading parameters to broker requirements.

**What gets normalized:**

- **Volumes:** Rounded to broker's volume step (e.g., 0.01 lot)
- **Prices:** Rounded to symbol's tick size/digits (e.g., 5 decimal places for EURUSD)
- **Stop Loss / Take Profit:** Adjusted to symbol precision

**Where:** Only MT5Sugar tier (MT5Service does not auto-normalize)

**Example:**
```go
// You pass: volume=0.0234, price=1.09876543
// Sugar normalizes to: volume=0.02, price=1.09877
ticket, _ := sugar.BuyMarket("EURUSD", 0.0234)
```

**Methods:**
```go
// Normalize price to symbol digits
normalizedPrice := sugar.NormalizePrice("EURUSD", 1.09876543)

// Normalize volume to broker step
normalizedVolume := sugar.NormalizeVolume("EURUSD", 0.0234)
```

---

### Risk-Based Volume Calculation
Calculate position size based on dollar risk rather than fixed lot size.

**Formula:** `volume = riskAmount / (stopLossPips × pipValue)`

**Parameters:**

- `riskAmount` - Dollar amount you're willing to risk (e.g., $50)
- `stopLossPips` - Distance to SL in pips (e.g., 20 pips)
- Result: Lot size that risks exactly $50 if SL hits

**Methods (MT5Sugar):**
```go
// Calculate volume for given risk
volume, err := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
// 2% risk with 50 pip SL

// Buy with risk-based calculation
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, err := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)

// Sell with risk-based calculation
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, err := sugar.SellMarketWithPips("EURUSD", lotSize, 50, 100)
```

**Use case:** Consistent risk management across all trades.

---

### Points vs Pips
**Point:** Smallest price movement for a symbol (1 tick).

**Pip:** Traditional forex unit (0.0001 for most pairs).

**Relationship:**

- **5-digit brokers:** 1 pip = 10 points (EURUSD: 1.10000 → 1.10010 = 1 pip)
- **3-digit brokers:** 1 pip = 1 point (USDJPY: 110.00 → 110.01 = 1 pip)

**In GoMT5:** All APIs use **points** for consistency.

**Conversion:**
```go
point, _ := service.GetPoint(ctx, "EURUSD")  // 0.00001
pips := points / 10.0  // For 5-digit pairs
```

**Why points?** Universal across all instruments (forex, metals, indices, crypto).

---

### Pip-Based Methods
Convenience methods that work with pips instead of points for forex traders.

**Methods (MT5Sugar):**

```go
// Market orders with SL/TP in pips
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.01, 20, 30)
// SL 20 pips, TP 30 pips

ticket, _ := sugar.SellMarketWithPips("GBPUSD", 0.01, 20, 30)

// Pending orders in pips
ticket, _ := sugar.BuyLimitPips("EURUSD", 0.01, 20, 15, 30)
// 20 pips below Ask, SL 15 pips, TP 30 pips

ticket, _ := sugar.BuyStopPips("EURUSD", 0.01, 20, 15, 30)
// 20 pips above Ask

ticket, _ := sugar.SellLimitPips("EURUSD", 0.01, 20, 15, 30)
// 20 pips above Bid

ticket, _ := sugar.SellStopPips("EURUSD", 0.01, 20, 15, 30)
// 20 pips below Bid
```

**Benefits:**

- No manual price calculations needed
- Think in pips (familiar for forex traders)
- Automatic normalization
- Fewer errors than absolute prices

---

### Trailing Stop
Dynamic Stop Loss that follows price in profit direction.

**How it works:**

1. Position opens with initial SL
2. When profit reaches threshold (e.g., +40 pips)
3. SL moves to breakeven or better
4. SL continues to follow price at fixed distance
5. Locks in profit as price moves favorably

**Implementation example:**
```go
currentProfit := position.Profit
trailingThreshold := 40.0 // pips
trailingDistance := 200.0  // distance in points (20 pips * 10 for 5-digit)

if currentProfit >= trailingThreshold*pipValue {
    point, _ := service.GetPoint(ctx, symbol)
    newSL := currentPrice - trailingDistance*point

    if newSL > currentSL {
        err := sugar.ModifyPosition(ticket, newSL, tp)
    }
}
```

**Use case:** Locking in profit on trending markets.

---

### Hedging
Opening opposite position to lock in current profit/loss level.

**How it works:**

1. Main position open (e.g., BUY EURUSD 0.1 lot)
2. Price moves against you (-50 pips)
3. Hedge triggered: SELL EURUSD 0.1 lot
4. Net position = 0 (loss locked at -50 pips)

**Purpose:**

- Lock losses instead of closing at stop-loss
- Protect position during volatility/news
- Wait for better opportunity to exit

**Implementation:**
```go
// Main position
buyTicket, _ := sugar.BuyMarket("EURUSD", 0.1)

// ... price moves against you ...

// Hedge
sellTicket, _ := sugar.SellMarket("EURUSD", 0.1)
// Net position = 0, loss locked
```

**⚠️ Note:** Not all brokers/regulations allow hedging. US brokers typically do not support hedging.

---

### Pending Order
Order that executes automatically when price reaches specified level.

**Types:**

- **BUY LIMIT:** Buy at price BELOW current (expect pullback, then up)
- **SELL LIMIT:** Sell at price ABOVE current (expect rally, then down)
- **BUY STOP:** Buy at price ABOVE current (breakout up)
- **SELL STOP:** Sell at price BELOW current (breakout down)

**Methods (absolute price):**
```go
ticket, _ := service.BuyLimit(ctx, symbol, volume, price, sl, tp, 0)
ticket, _ := service.SellLimit(ctx, symbol, volume, price, sl, tp, 0)
ticket, _ := service.BuyStop(ctx, symbol, volume, price, sl, tp, 0)
ticket, _ := service.SellStop(ctx, symbol, volume, price, sl, tp, 0)
```

**Methods (pip-based - Sugar):**
```go
ticket, _ := sugar.BuyLimitPips(symbol, volume, offsetPips, slPips, tpPips)
ticket, _ := sugar.SellStopPips(symbol, volume, offsetPips, slPips, tpPips)
```

---

## 🔌 gRPC and Protocol Terms

### gRPC
High-performance RPC (Remote Procedure Call) framework using HTTP/2.

**In GoMT5:**

- MT5Account tier sends gRPC requests to MT5 terminal
- Terminal runs gRPC server (configured via gateway)
- Request/Response pattern for all operations
- Context-based timeout management

**Connection setup:**
```go
account, err := mt5.NewMT5Account(ctx, "mt5.mrpc.pro:443")
if err != nil {
    log.Fatal(err)
}

err = account.ConnectEx(ctx, &pb.ConnectExRequest{
    Uuid:      591129415,
    Password:  "password",
    MtCluster: "FxPro-MT5 Demo",
})
```


### Context Pattern
Standard Go way to manage request lifetime, deadlines, and cancellation.

**All methods accept context.Context:**
```go
// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Use in API calls
result, err := account.AccountSummary(ctx, req)

// Context with cancellation (for streams)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

**Why context?**

- Timeout management
- Cancellation propagation
- Cleanup guarantees with defer
- Standard Go pattern

**Documentation:** [Go Context Package](https://pkg.go.dev/context)

---

### Channel Pattern
Go's built-in mechanism for goroutine communication, used for streaming.

**Streaming returns two channels:**
```go
dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data, ok := <-dataChan:
        if !ok {
            return
        }
        // Process tick data

    case err, ok := <-errChan:
        if !ok {
            return
        }
        // Handle error

    case <-ctx.Done():
        return
    }
}
```

**Critical rules:**

- **Always read from BOTH channels** (prevents goroutine leaks)
- **Always use context cancellation** (cleanup)
- **Always defer cancel()** (ensures cleanup)

**Documentation:** [gRPC Stream Management](GRPC_STREAM_MANAGEMENT.md)

---

### Return Codes
Protobuf return codes indicating operation success/failure.

**Common codes:**

- **10009** = Success (TRADE_RETCODE_DONE)
- **10008** = Pending order placed (TRADE_RETCODE_PLACED)
- **10004** = Requote
- **10006** = Request rejected
- **10013** = Invalid request
- **10014** = Invalid volume
- **10015** = Invalid price
- **10016** = Invalid stops
- **10018** = Market closed
- **10019** = Not enough money
- **10031** = No connection to trade server

**Always check `ReturnedCode` in trading operations.**

**Helper functions:**
```go
import "your-project/examples/errors"

if errors.IsRetCodeSuccess(result.ReturnedCode) {
    fmt.Println("Trade successful!")
}

if errors.IsRetCodeRetryable(result.ReturnedCode) {
    // Retry with exponential backoff
}

message := errors.GetRetCodeMessage(result.ReturnedCode)
```

**Documentation:** [Return Codes Reference](RETURN_CODES_REFERENCE.md)

---

## 🎓 Trading Terms (Project Context)

### Basic Trading Terms

Core trading concepts that every trader must understand before working with MT5.

| Term | Description | Example / Details |
|------|-------------|-------------------|
| **Bid** | Price at which broker **buys** from you (you **sell** at Bid) | EURUSD Bid: 1.10000 |
| **Ask** | Price at which broker **sells** to you (you **buy** at Ask) | EURUSD Ask: 1.10002 |
| **Spread** | Difference between Ask and Bid (broker's commission) | Ask 1.10002 - Bid 1.10000 = 2 pips spread |
| **Long / Buy** | Opening a BUY position (profit when price goes up) | Buy EURUSD at 1.10000, close at 1.10050 = +50 pips profit |
| **Short / Sell** | Opening a SELL position (profit when price goes down) | Sell EURUSD at 1.10000, close at 1.09950 = +50 pips profit |
| **Position** | Open trade (BUY or SELL) that can result in profit or loss | BUY 0.1 lot EURUSD at 1.10000 (currently in profit +$50) |
| **Order** | Command to open a position | Market Order (immediate), Pending Order (at specific price) |
| **Market Order** | Order that executes immediately at current price | BUY EURUSD at current Ask price |
| **Pending Order** | Order that triggers automatically when price reaches level | BUY LIMIT at 1.09950 (when price drops to that level) |
| **Lot / Volume** | Position size (1 lot = 100,000 units of base currency) | 0.01 lot = 1,000 units (micro lot) |
| **Stop Loss (SL)** | Price level for automatic loss-limiting exit | SL at 1.09950 limits loss to 50 pips |
| **Take Profit (TP)** | Price level for automatic profit-taking exit | TP at 1.10050 locks in 50 pips profit |
| **Pip** | Standard unit of price movement (0.0001 for most forex) | EURUSD: 1.10000 → 1.10001 = 1 pip |
| **Point** | Smallest price movement (1 tick) | For 5-digit broker: 1 pip = 10 points |
| **Leverage** | Borrowed capital from broker (e.g., 1:100) | With 1:100, you control $100,000 with $1,000 margin |
| **Margin** | Funds locked as collateral for open positions | Open 1 lot EURUSD with 1:100 leverage = $1,000 margin |
| **Free Margin** | Available funds for opening new positions | Balance $10,000 - Margin $1,000 = Free Margin $9,000 |
| **Balance** | Account funds without considering open positions | Initial deposit + closed trades profit/loss |
| **Equity** | Current account value (Balance ± floating profit/loss) | Balance $10,000 + Floating profit $500 = Equity $10,500 |
| **Floating P/L** | Unrealized profit/loss of open positions | Position currently in +$50 profit (not yet closed) |
| **Margin Level** | Ratio of Equity to Margin (%) | Equity $10,500 / Margin $1,000 = 1050% |
| **Margin Call** | Warning when Margin Level drops below threshold | Broker warns at 100% Margin Level |
| **Stop Out** | Forced closure of positions when Margin Level critically low | Broker closes positions at 50% Margin Level |
| **Drawdown** | Decline from peak balance/equity | Peak $11,000 → Current $10,000 = $1,000 drawdown (9.09%) |
| **Slippage** | Difference between expected and executed price | Requested 1.10000, executed at 1.10003 = 3 pips slippage |

**Critical for API usage:**

- `GetBid()` / `GetAsk()` - Current prices
- `GetBalance()` / `GetEquity()` / `GetMargin()` - Account state
- `BuyMarket()` / `SellMarket()` - Open positions
- `ClosePosition()` - Close positions
- `ModifyPosition(ticket, sl, tp)` - Modify SL/TP

---

### Risk Amount
Dollar amount you're willing to risk if Stop Loss hits.

**Example:** $50 risk per trade means if SL hits, you lose exactly $50.

**Used in:**

```go
// Calculate volume for 2% risk with 50 pip SL
volume, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// Open with risk-based calculation
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
// Risk exactly 2% of balance
```

---

### Breakeven
Moving Stop Loss to entry price to eliminate risk.

**Example:**

- Entry: BUY at 1.10000
- Price rises to 1.10050 (+50 pips profit)
- Move SL from 1.09950 to 1.10000 (breakeven)
- Now risk-free: worst case = breakeven

**Implementation:**
```go
entryPrice := position.PriceOpen
breakevenThreshold := 40.0 // pips

if currentProfit >= breakevenThreshold*pipValue {
    err := sugar.ModifyPosition(ticket, entryPrice, tp)
}
```

---

### Grid Trading
Strategy placing multiple pending orders at equal intervals above and below current price.

**How it works:**

1. Place Buy Limit orders below current price (e.g., every 20 pips)
2. Place Sell Limit orders above current price (e.g., every 20 pips)
3. As price oscillates, orders trigger and close at TP
4. Works best in ranging markets

**Example:**
```
Price: 1.10000
Buy Limits:  1.09980 (-20 pips), 1.09960 (-40 pips), 1.09940 (-60 pips)
Sell Limits: 1.10020 (+20 pips), 1.10040 (+40 pips), 1.10060 (+60 pips)
```

**Implementation:**
```go
gridStep := 20.0 // pips
levels := 5

for i := 1; i <= levels; i++ {
    sugar.BuyLimitPips("EURUSD", 0.01, float64(i)*gridStep, 15, 30)
    sugar.SellLimitPips("EURUSD", 0.01, float64(i)*gridStep, 15, 30)
}
```

**⚠️ Risk:** High volatility can trigger many orders, increasing exposure.

---

### Scalping
Fast trading strategy with small profits (5-25 pips) and tight stops (10-20 pips).

**Characteristics:**

- Very short hold times (seconds-minutes)
- High win rate, small profit per trade
- Risk:Reward typically 1:1 to 1:2
- Requires low spreads

**Implementation:**
```go
// Scalp with 10 pip SL, 15 pip TP
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.01, 10, 15)
```

---

### News Trading
Placing orders on both sides of price before major news events.

**How it works:**

1. Before news: Place Buy Stop above + Sell Stop below
2. News releases: Price breaks in one direction
3. One order triggers, cancel the other
4. Capture volatility spike

**Example:**
```
Current price: 1.10000
Buy Stop:  1.10030 (+30 pips)
Sell Stop: 1.09970 (-30 pips)
News → Price spikes to 1.10050 → Buy Stop triggers → Cancel Sell Stop
```

**Implementation:**
```go
buyTicket, _ := sugar.BuyStopPips("EURUSD", 0.01, 30, 20, 40)
sellTicket, _ := sugar.SellStopPips("EURUSD", 0.01, 30, 20, 40)

// ... wait for news ...

// Cancel unfilled order
sugar.DeleteOrder(sellTicket)
```

**⚠️ Risk:** Slippage during news can trigger both orders.

---

### Volume Limits
Broker restrictions on position size.

**Get via:**
```go
limits, err := service.GetVolumeLimits(ctx, "EURUSD")
minVol := limits.MinVolume    // e.g., 0.01
maxVol := limits.MaxVolume    // e.g., 100.0
stepVol := limits.VolumeStep  // e.g., 0.01
```

**Used for:** Auto-normalization in MT5Sugar.

**Validation:**
```go
if volume < minVol {
    volume = minVol
}
if volume > maxVol {
    volume = maxVol
}

// Round to step
volume = math.Round(volume/stepVol) * stepVol
```

---

### main.go
Main entry point that routes commands to corresponding examples.

**Key characteristics:**

- Single entry point for all executable code
- Command routing by number or alias
- Helpful error messages for unknown commands
- Clean shutdown handling

**How it works:**
```
go run main.go 4
    ↓
main.go routes to
    ↓
helpers/04_market_orders.go
    ↓
RunMarketOrdersDemo()
```

---

### config Package
Configuration loader with support for JSON files and environment variables.

**Methods:**

```go
import "your-project/examples/demos/config"

// Load from config.json or env vars
cfg, err := config.LoadConfig()

// Use configuration
fmt.Printf("Host: %s\n", cfg.Host)
fmt.Printf("Cluster: %s\n", cfg.Cluster)
```

**Priority:**
1. `config/config.json` (if exists)
2. Environment variables (fallback)
3. Error if nothing found

**Location:** `examples/demos/config/config.go`

**Documentation:** Detailed comments in config.go header

---

### Protobuf Inspector
Interactive tool for exploring MT5 API types.

**Features:**

- Explore 267 protobuf types
- Search types and fields
- View enum values (67 enums)
- Understand message structures

**Commands:**

- `list` - Show all types
- `<TypeName>` - Explore type
- `search <text>` - Find types
- `field <name>` - Find types with field
- `enum <name>` - Show enum values

**How to run:**
```bash
cd examples/demos
go run main.go inspect
```

**Location:** `examples/demos/helpers/17_protobuf_inspector.go`

**Documentation:** [Protobuf Inspector Guide](PROTOBUF_INSPECTOR_GUIDE.md)

---

## ⚙️ Configuration Terms

### config.json
Configuration file for MT5 connection settings.
Simply put, where you enter your account credentials to run demonstration examples.

**Location:** `examples/demos/config/config.json`

**Format:**
```json
{
  "user": 591129415,
  "password": "YourPassword",
  "host": "mt5.mrpc.pro",
  "port": 443,
  "grpc_server": "mt5.mrpc.pro:443",
  "mt_cluster": "FxPro-MT5 Demo",
  "test_symbol": "EURUSD",
  "test_volume": 0.01
}
```

---


### Environment Variables
Alternative to config.json for configuration.

**Variables:**

```bash
# Required
MT5_USER=591129415
MT5_PASSWORD="YourPassword"
MT5_HOST="mt5.mrpc.pro"

# Optional (have defaults)
MT5_PORT=443
MT5_CLUSTER="FxPro-MT5 Demo"
MT5_TEST_SYMBOL="EURUSD"
MT5_TEST_VOLUME=0.01
```

**Set (Linux/Mac):**
```bash
export MT5_USER=591129415
export MT5_PASSWORD="YourPassword"
```

**Set (Windows PowerShell):**
```powershell
$env:MT5_USER="591129415"
$env:MT5_PASSWORD="YourPassword"
```

**Priority:** config.json takes precedence over env vars.

---

## 🔗 Cross-Component Terms

### Goroutine
Lightweight execution thread in Go.

**In GoMT5:**

- Streaming methods spawn background goroutines
- Must cancel context to stop goroutines
- Must read from both channels to prevent leaks

**Example:**
```go
// Streaming spawns goroutine
dataChan, errChan := account.OnSymbolTick(ctx, req)

// Background processing
go func() {
    for data := range dataChan {
        processData(data)
    }
}()

// Stop goroutine
cancel()  // Triggers ctx.Done()
```

**Critical:** Always `defer cancel()` to prevent goroutine leaks.

**Documentation:** [gRPC Stream Management](GRPC_STREAM_MANAGEMENT.md)

---

### Batch Operations
Executing action across multiple positions/orders at once.

**Methods (MT5Sugar):**
```go
// Close positions
count, _ := sugar.CloseAllPositions("")          // All positions
count, _ := sugar.CloseAllPositions("EURUSD")    // By symbol

// Delete orders
count, _ := sugar.DeleteAllOrders("")            // All pending orders
count, _ := sugar.DeleteAllOrders("EURUSD")      // By symbol
```

**Use case:** Emergency exits, end-of-session cleanup, strategy resets.

---

### History Queries
Retrieving past orders and positions for analysis.

**Methods (MT5Account):**
```go
// Get order history with pagination
req := &pb.OrderHistoryRequest{
    FromDate: fromTimestamp,
    ToDate:   toTimestamp,
    Offset:   0,
    Limit:    100,
}
history, err := account.OrderHistory(ctx, req)

// Get position history
req := &pb.PositionsHistoryRequest{
    FromDate: fromTimestamp,
    ToDate:   toTimestamp,
    Offset:   0,
    Limit:    100,
}
positions, err := account.PositionsHistory(ctx, req)
```

**Use case:** Performance analysis, trade logs, strategy validation.

---

### Educational Project

GoMT5 examples and helpers are educational materials and API demonstrations, NOT production trading systems.

**Implications:**

- ✅ Study code and patterns
- ✅ Modify for your needs
- ✅ Test on demo accounts
- ✅ Use as templates for building your strategies
- ❌ Don't use as-is with real money
- ❌ Don't expect production-grade risk management
- ❌ Don't expect proper error handling for all edge cases

**Remember:** These are educational examples, not battle-tested production systems.

---


## 📚 See Also

- **[MT5Account Master Overview](../MT5Account/MT5Account.Master.Overview.md)** - Complete low-level API reference
- **[MT5Service API Overview](../MT5Service/MT5Service.Overview.md)** - Mid-level wrappers API
- **[MT5Sugar API Overview](../MT5Sugar/MT5Sugar.API_Overview.md)** - High-level Sugar API
- **[gRPC Stream Management](GRPC_STREAM_MANAGEMENT.md)** - Guide to streaming subscriptions
- **[Return Codes Reference](RETURN_CODES_REFERENCE.md)** - Complete return codes reference
- **[User Code Sandbox Guide](USERCODE_SANDBOX_GUIDE.md)** - How to write custom code
- **[Getting Started](GETTING_STARTED.md)** - Setup and quick start guide

---

## Go-Specific Terms

### defer Statement
Go's cleanup mechanism that executes when function returns.

**Critical pattern:**
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // ✅ ALWAYS cleanup
```

**Why important:** Guarantees goroutine cleanup and resource release.

---

### select Statement
Go's multiplexer for channel operations.

**Use in streaming:**
```go
for {
    select {
    case data := <-dataChan:
        // Process data
    case err := <-errChan:
        // Handle error
    case <-ctx.Done():
        return  // Cleanup
    }
}
```

**Why important:** Read from multiple channels simultaneously.

---

### Interface Composition
Go's way of combining types.

**In GoMT5:**
```go
type MT5Sugar struct {
    *MT5Service  // Embeds MT5Service
}

type MT5Service struct {
    account *MT5Account  // Uses MT5Account
}
```

**Result:** Sugar has all Service methods + its own methods.


"Trade safely, code cleanly, and may your contexts always cancel gracefully."

---
