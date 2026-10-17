# GoMT5 Glossary

> Project-specific terms and concepts. This glossary covers GoMT5 architecture, components, and trading automation terminology used throughout the codebase.

---

## 🏗️ Architectural Terms

### Three-Tier Architecture

Core design pattern of GoMT5 with three abstraction levels:

- **Tier 1 (MT5Account):** Low-level gRPC communication with MT5 terminal
- **Tier 2 (MT5Service):** Wrapper methods with simplified signatures
- **Tier 3 (MT5Sugar):** High-level convenience methods (one-liners, SL/TP helpers, risk-based sizing)

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

**Location:** `package/Helpers/MT5Account.go` (2400+ lines)

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
- No normalization (values are sent as you pass them; you control precision)
- **36 methods**: account, symbols, positions & orders, market depth, trading, 5 streams
- Typed result structs (`AccountSummary`, `SymbolTick`, `OrderResult`, ...)

**When to use:** Need Go types but full control over request parameters, building custom strategies.

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

- Volumes and prices are sent exactly as given (no rounding - see [Normalization](#normalization))
- Risk management helpers
- Batch operations (CloseAllPositions, CloseAllBySymbol)
- Distance-based SL/TP helpers (BuyMarketWithPips, CalculateSLTP, etc.)
- One-line operations with built-in timeouts
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

// Risk-based position sizing (2% risk, SL 50 points away)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, _ = sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

---

## 🔧 Technical Concepts

### Normalization
Adjusting trading parameters to the broker's requirements before sending them.

**GoMT5 does NOT normalize automatically.** MT5Account, MT5Service and MT5Sugar send volumes and prices exactly as you pass them.
A volume that is not a multiple of the symbol's volume step (or is outside its min/max) is refused by the broker with
return code 10014 (invalid volume); a price with too many decimals or too close to the market can be refused with 10015/10016.

**The only exception:** `sugar.CalculatePositionSize()` returns a volume already rounded **down** to `VolumeStep` and clamped to `VolumeMin`/`VolumeMax`, then capped by free margin (an error if that leaves it below `VolumeMin`).

**What you should normalize yourself:**

- **Volumes:** to the symbol's volume step (e.g., 0.01 lot) within min/max
- **Prices / SL / TP:** to the symbol's digits (e.g., 5 decimal places for EURUSD)

**How (using `GetSymbolInfo`):**
```go
info, err := sugar.GetSymbolInfo("EURUSD") // Digits, Point, VolumeMin, VolumeMax, VolumeStep, StopLevel, ...
if err != nil {
    return
}

// Volume: round down to the step, keep within limits
volume := math.Floor(0.0234/info.VolumeStep) * info.VolumeStep // 0.02
volume = math.Max(info.VolumeMin, math.Min(volume, info.VolumeMax))

// Price: round to the symbol's digits
factor := math.Pow(10, float64(info.Digits))
price := math.Round(1.09876543*factor) / factor // 1.09877

ticket, err := sugar.BuyLimit("EURUSD", volume, price)
```

---

### Risk-Based Volume Calculation
Calculate position size based on a percentage of the balance you are willing to risk rather than a fixed lot size.

**Formula (as implemented in `CalculatePositionSize`):**
`volume = (balance × riskPercent / 100) / (stopLossPips × ContractSize × Point)`

**Parameters:**

- `riskPercent` - Percentage of the balance to risk (e.g., 2.0 = 2%)
- `stopLossPips` - Distance to SL, multiplied by the symbol's `Point` (so it is measured in **points**, e.g. 50 = 5 pips on EURUSD)
- Result: rounded down to the volume step, clamped to min/max, and capped so that it needs at most 80% of the free margin

**Methods (MT5Sugar):**
```go
// Calculate volume for 2% risk with SL 50 points away
volume, err := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// Buy with risk-based calculation
buyTicket, err := sugar.BuyMarketWithPips("EURUSD", volume, 50, 100)

// Sell with risk-based calculation
sellTicket, err := sugar.SellMarketWithPips("EURUSD", volume, 50, 100)
```

**Use case:** Consistent risk management across all trades.

---

### Points vs Pips
**Point:** Smallest price movement for a symbol (1 tick).

**Pip:** Traditional forex unit (0.0001 for most pairs).

**Relationship:**

- **5-digit quotes:** 1 pip = 10 points (EURUSD: 1.10000 → 1.10010 = 1 pip)
- **3-digit JPY quotes:** 1 pip = 10 points (USDJPY: 110.000 → 110.010 = 1 pip)
- **4-digit / 2-digit quotes:** 1 pip = 1 point

**In GoMT5:** distances are measured in **points**. The Sugar helpers whose parameters are named `...Pips`
(`CalculateSLTP`, `BuyMarketWithPips`, `SellMarketWithPips`, `CalculatePositionSize`) multiply them by the symbol's `Point`,
so they are points too: on 5-digit EURUSD, `50` means 50 points = 5 pips. Stop level (`GetMinStopLevel`, `SymbolInfo.StopLevel`) and `SymbolInfo.Spread` are also in points.

**Conversion:**
```go
point, _ := service.GetSymbolDouble(ctx, "EURUSD", pb.SymbolInfoDoubleProperty_SYMBOL_POINT) // 0.00001
points := 150.0
pips := points / 10.0         // For 5-digit pairs: 15 pips
distance := points * point    // Price distance: 0.00150
```

**Why points?** Universal across all instruments (forex, metals, indices, crypto).

---

### Pip-Based Methods
Convenience methods that take SL/TP as distances from the entry price instead of absolute prices.
The distances are multiplied by the symbol's `Point` (see [Points vs Pips](#points-vs-pips)).

**Methods (MT5Sugar):**

```go
// Market orders with SL/TP as distances
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.01, 200, 300)
// SL 200 points (20 pips) below entry, TP 300 points (30 pips) above

ticket, _ = sugar.SellMarketWithPips("GBPUSD", 0.01, 200, 300)

// Pending orders: there are no pip-based pending methods.
// Compute the levels with CalculateSLTP, then place the order with absolute prices:
ask, _ := sugar.GetAsk("EURUSD")
point := 0.00001                        // or sugar.GetSymbolInfo("EURUSD") → info.Point
entry := ask - 200*point                // 20 pips below Ask
sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", entry, 150, 300)
ticket, _ = sugar.BuyLimitWithSLTP("EURUSD", 0.01, entry, sl, tp)

// SellLimitWithSLTP works the same way with "SELL" (entry above Bid).
// BuyStop / SellStop take only (symbol, volume, price) - no SL/TP.
```

**Benefits:**

- No manual SL/TP price calculations for market orders
- Think in distances instead of absolute prices
- Fewer errors than absolute prices

**Note:** the computed prices are not rounded to the symbol's digits (see [Normalization](#normalization)).

---

### Trailing Stop
Dynamic Stop Loss that follows price in profit direction.

**How it works:**

1. Position opens with initial SL
2. When profit reaches threshold (e.g., +40 pips)
3. SL moves to breakeven or better
4. SL continues to follow price at fixed distance
5. Locks in profit as price moves favorably

**Implementation example (BUY position):**
```go
position, _ := sugar.GetPositionByTicket(ticket)
point, _ := service.GetSymbolDouble(ctx, position.Symbol, pb.SymbolInfoDoubleProperty_SYMBOL_POINT)
bid, _ := sugar.GetBid(position.Symbol)

activation := 400.0       // points of profit before trailing starts (40 pips on 5-digit)
trailingDistance := 200.0 // distance in points (20 pips * 10 for 5-digit)

if bid-position.PriceOpen >= activation*point {
    newSL := bid - trailingDistance*point

    if newSL > position.StopLoss {
        err := sugar.ModifyPositionSL(ticket, newSL)
        if err != nil {
            fmt.Println("SL update failed:", err)
        }
    }
}
```

A complete trailing-stop manager is in `examples/demos/orchestrators/11_trailing_stop.go` (menu item 11).

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

**Methods (absolute price - Sugar):**
```go
ticket, _ := sugar.BuyLimit(symbol, volume, price)
ticket, _ = sugar.SellLimit(symbol, volume, price)
ticket, _ = sugar.BuyStop(symbol, volume, price)
ticket, _ = sugar.SellStop(symbol, volume, price)

// Limit orders with SL/TP
ticket, _ = sugar.BuyLimitWithSLTP(symbol, volume, price, sl, tp)
ticket, _ = sugar.SellLimitWithSLTP(symbol, volume, price, sl, tp)
```

**SL/TP from distances (Sugar):**
```go
sl, tp, _ := sugar.CalculateSLTP(symbol, "BUY", price, 150, 300) // distances in points
ticket, _ := sugar.BuyLimitWithSLTP(symbol, volume, price, sl, tp)
```

**Full control (Service):** `service.PlaceOrder(ctx, &pb.OrderSendRequest{...})` with `Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT` (or `_SELL_LIMIT`, `_BUY_STOP`, `_SELL_STOP`, `_BUY_STOP_LIMIT`, `_SELL_STOP_LIMIT`).

**Delete a pending order:** `service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: ticket})`.

---

## 🔌 gRPC and Protocol Terms

### gRPC
High-performance RPC (Remote Procedure Call) framework using HTTP/2.

**In GoMT5:**

- MT5Account tier sends gRPC requests to the MetaRPC gateway (`mt5.mrpc.pro:443`), which runs the MT5 terminal
- Request/Response pattern for all operations
- Context-based timeout management

**Connection setup:**
```go
// import mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
// The API key is read from MRPC_API_KEY (or set account.ApiKey); ConnectEx sets the terminal id.
account, err := mt5.NewMT5Account(591129415, "password", "mt5.mrpc.pro:443", uuid.Nil)
if err != nil {
    log.Fatal(err)
}
defer account.DisconnectAndClose()

_, err = account.ConnectEx(ctx, &pb.ConnectExRequest{
    User:          591129415,
    Password:      "password",
    MtClusterName: "MetaQuotes-Demo", // broker server name as shown in MetaTrader
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
result, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})

// Context with cancellation (for streams)
streamCtx, stopStream := context.WithCancel(context.Background())
defer stopStream()
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
req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
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

- **Always read from BOTH channels** (the data channel is unbuffered - an unread one blocks the stream goroutine; the error channel tells you why the stream ended)
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

**Always check `ReturnedCode` in trading operations.** A request the terminal refuses often does not come back as a
`ReturnedCode` at all but as an error: an `*mt5.ApiError` whose `MqlErrorTradeIntCode()` carries the code
(e.g. 10014 for an invalid volume, enum `pb.MqlErrorTradeCode`). Check both.

**Helper functions:**
```go
import mt5 "github.com/MetaRPC/GoMT5/package/Helpers"

result, err := account.OrderSend(ctx, req)
var apiErr *mt5.ApiError
if errors.As(err, &apiErr) {
    fmt.Println("Refused:", apiErr.MqlErrorTradeIntCode(), apiErr.MqlErrorTradeDescription())
    return
}

if mt5.IsRetCodeSuccess(result.ReturnedCode) {
    fmt.Println("Trade successful!")
}

if mt5.IsRetCodeRetryable(result.ReturnedCode) {
    // Retry with exponential backoff
}

message := mt5.GetRetCodeMessage(result.ReturnedCode)
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
| **Spread** | Difference between Ask and Bid (broker's commission) | Ask 1.10002 - Bid 1.10000 = 2 points (0.2 pip) spread |
| **Long / Buy** | Opening a BUY position (profit when price goes up) | Buy EURUSD at 1.10000, close at 1.10050 = +50 pips profit |
| **Short / Sell** | Opening a SELL position (profit when price goes down) | Sell EURUSD at 1.10000, close at 1.09950 = +50 pips profit |
| **Position** | Open trade (BUY or SELL) that can result in profit or loss | BUY 0.1 lot EURUSD at 1.10000 (currently in profit +$50) |
| **Order** | Command to open a position | Market Order (immediate), Pending Order (at specific price) |
| **Market Order** | Order that executes immediately at current price | BUY EURUSD at current Ask price |
| **Pending Order** | Order that triggers automatically when price reaches level | BUY LIMIT at 1.09950 (when price drops to that level) |
| **Lot / Volume** | Position size (1 lot = 100,000 units of base currency) | 0.01 lot = 1,000 units (micro lot) |
| **Stop Loss (SL)** | Price level for automatic loss-limiting exit | SL at 1.09950 limits loss to 50 pips |
| **Take Profit (TP)** | Price level for automatic profit-taking exit | TP at 1.10050 locks in 50 pips profit |
| **Pip** | Standard unit of price movement (0.0001 for most forex) | EURUSD: 1.10000 → 1.10010 = 1 pip |
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
| **Slippage** | Difference between expected and executed price | Requested 1.10000, executed at 1.10003 = 3 points (0.3 pip) slippage |

**Critical for API usage:**

- `GetBid()` / `GetAsk()` - Current prices
- `GetBalance()` / `GetEquity()` / `GetMargin()` - Account state
- `BuyMarket()` / `SellMarket()` - Open positions
- `ClosePosition()` - Close positions
- `ModifyPositionSLTP(ticket, sl, tp)` - Modify SL/TP

---

### Risk Amount
Dollar amount you're willing to risk if Stop Loss hits.

**Example:** $50 risk per trade means if SL hits, you lose exactly $50.

**Used in:**

```go
// Calculate volume for 2% risk with SL 500 points (50 pips) away
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)

// Open with the same SL distance (TP 1000 points)
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 500, 1000)
// Risks about 2% of balance (less if the volume was capped by free margin or rounded down to the step)
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
position, _ := sugar.GetPositionByTicket(ticket)
entryPrice := position.PriceOpen
breakevenThreshold := 400.0 // points (40 pips on 5-digit)

point, _ := service.GetSymbolDouble(ctx, position.Symbol, pb.SymbolInfoDoubleProperty_SYMBOL_POINT)
bid, _ := sugar.GetBid(position.Symbol)

if bid-entryPrice >= breakevenThreshold*point { // BUY position
    err := sugar.ModifyPositionSL(ticket, entryPrice)
    if err != nil {
        fmt.Println("Breakeven failed:", err)
    }
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
gridStep := 200.0 // points (20 pips on 5-digit)
levels := 5

info, _ := sugar.GetSymbolInfo("EURUSD") // Bid, Ask, Point, Digits
for i := 1; i <= levels; i++ {
    offset := float64(i) * gridStep * info.Point

    buyPrice := info.Ask - offset
    sl, tp, _ := sugar.CalculateSLTP("EURUSD", "BUY", buyPrice, 150, 300)
    sugar.BuyLimitWithSLTP("EURUSD", 0.01, buyPrice, sl, tp)

    sellPrice := info.Bid + offset
    sl, tp, _ = sugar.CalculateSLTP("EURUSD", "SELL", sellPrice, 150, 300)
    sugar.SellLimitWithSLTP("EURUSD", 0.01, sellPrice, sl, tp)
}
```

A complete grid strategy is in `examples/demos/orchestrators/13_grid_trader.go` (menu item 13).

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
// Scalp with 10 pip SL, 15 pip TP (distances in points: 100 and 150 on 5-digit EURUSD)
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.01, 100, 150)
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
info, _ := sugar.GetSymbolInfo("EURUSD")
offset := 300 * info.Point // 300 points = 30 pips

buyTicket, _ := sugar.BuyStop("EURUSD", 0.01, info.Ask+offset)
sellTicket, _ := sugar.SellStop("EURUSD", 0.01, info.Bid-offset)

// ... wait for news ...

// Cancel the unfilled order (closing a pending order deletes it)
service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: sellTicket})
_ = buyTicket
```

**⚠️ Risk:** Slippage during news can trigger both orders.

---

### Volume Limits
Broker restrictions on position size.

**Get via:**
```go
info, err := sugar.GetSymbolInfo("EURUSD")
minVol := info.VolumeMin    // e.g., 0.01
maxVol := info.VolumeMax    // e.g., 100.0
stepVol := info.VolumeStep  // e.g., 0.01
```

The same values are available one by one through `service.GetSymbolDouble(ctx, "EURUSD", pb.SymbolInfoDoubleProperty_SYMBOL_VOLUME_MIN)` (`_VOLUME_MAX`, `_VOLUME_STEP`).

**Used by:** `sugar.CalculatePositionSize()` (rounds down to the step and clamps to min/max). All other methods send the volume as given, so validate it yourself.

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

- Single entry point for all demos (`examples/demos/main.go`)
- Interactive menu when started without arguments
- Command routing by number (1-18) or alias (`service`, `sugar06`, `grid`, `inspect`, ...)
- Optional API key as `--api-key <key>` (or a second argument)

**How it works:**
```
cd examples/demos
go run main.go 4
    ↓
main.go routes to
    ↓
service/04_service_demo.go
    ↓
service.RunService04()
```

---

### config Package
Configuration loader with support for JSON files and environment variables.

**Methods:**

```go
import "github.com/MetaRPC/GoMT5/examples/demos/config"

// Load from config.json or env vars
cfg, err := config.LoadConfig()

// Use configuration
fmt.Printf("Host: %s\n", cfg.Host)
fmt.Printf("Cluster: %s\n", cfg.MtCluster)
```

**Priority:**
1. `config/config.json` (also looked up as `examples/demos/config/config.json` and `config.json`)
2. Environment variables (fallback, when `MT5_USER`, `MT5_PASSWORD` and `MT5_HOST` are set)
3. Built-in defaults (`mt5.mrpc.pro:443`, `MetaQuotes-Demo`, `EURUSD`, `0.01`)

If the resulting login is 0, the password is empty or the cluster is `FxPro-MT5 Demo`, `LoadConfig()` opens a fresh demo account on
`MetaQuotes-Demo` and uses it. The API key comes from `api_key`, `MRPC_API_KEY` or `--api-key` (default `TRIAL`).

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
  "mt_cluster": "MetaQuotes-Demo",
  "api_key": "",
  "test_symbol": "EURUSD",
  "test_volume": 0.01
}
```

Leave `user` at 0 (as shipped) to have a demo account opened automatically.

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
MT5_GRPC_SERVER="mt5.mrpc.pro:443"   # default: MT5_HOST:MT5_PORT
MT5_CLUSTER="MetaQuotes-Demo"
MT5_TEST_SYMBOL="EURUSD"
MT5_TEST_VOLUME=0.01
MRPC_API_KEY="your-api-key"          # default: TRIAL
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
- Must keep reading the data channel (or cancel) to prevent leaks, and read the error channel to learn why a stream stopped

**Example:**
```go
// Streaming spawns goroutine
ctx, cancel := context.WithCancel(context.Background())
dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}})

// Background processing
go func() {
    for data := range dataChan {
        fmt.Println(data) // process the tick
    }
}()
go func() {
    for err := range errChan {
        fmt.Println("stream error:", err)
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
// Close positions (pending orders are not touched)
count, _ := sugar.CloseAllPositions()            // All positions
count, _ = sugar.CloseAllBySymbol("EURUSD")      // By symbol

// Delete pending orders - there is no batch method; loop over them (MT5Service)
orders, _ := service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
for _, o := range orders.OpenedOrders {
    service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: o.Ticket})
}
```

**Use case:** Emergency exits, end-of-session cleanup, strategy resets.

---

### History Queries
Retrieving past orders and positions for analysis.

**Methods (MT5Account):**
```go
// Get order history with pagination
from := timestamppb.New(time.Now().AddDate(0, 0, -7))
to := timestamppb.Now()
req := &pb.OrderHistoryRequest{
    InputFrom:     from,
    InputTo:       to,
    InputSortMode: pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_CLOSE_TIME_DESC,
    PageNumber:    0,
    ItemsPerPage:  100,
}
history, err := account.OrderHistory(ctx, req)

// Get position history (paging fields are optional pointers)
page, perPage := int32(0), int32(100)
posReq := &pb.PositionsHistoryRequest{
    SortType:             pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_DESC,
    PositionOpenTimeFrom: from,
    PositionOpenTimeTo:   to,
    PageNumber:           &page,
    ItemsPerPage:         &perPage,
}
positions, err := account.PositionsHistory(ctx, posReq)
```

The same queries with Go types: `service.GetOrderHistory(...)`, `service.GetPositionsHistory(...)`, or simply `sugar.GetDealsToday()` / `GetDealsThisWeek()` / `GetDealsDateRange(from, to)`.

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
    case data, ok := <-dataChan:
        if !ok {
            return // stream ended
        }
        fmt.Println(data) // Process data
    case err, ok := <-errChan:
        if !ok {
            return
        }
        fmt.Println(err) // Handle error
    case <-ctx.Done():
        return  // Cleanup
    }
}
```

**Why important:** Read from multiple channels simultaneously.

---

### Composition
Go's way of combining types: a struct holds another one as a field.

**In GoMT5:**
```go
type MT5Sugar struct {
    service *MT5Service // wraps MT5Service (unexported field, not embedded)
    // ...
}

type MT5Service struct {
    account *helpers.MT5Account // wraps MT5Account
}
```

**Result:** Sugar does not inherit the Service methods. Reach the lower tiers explicitly with
`sugar.GetService()` (→ `*MT5Service`) and `sugar.GetAccount()` / `service.GetAccount()` (→ `*MT5Account`).


"Trade safely, code cleanly, and may your contexts always cancel gracefully."

---
