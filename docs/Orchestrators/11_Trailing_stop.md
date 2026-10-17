# TrailingStopManager - Automated Profit Protection

## Description

**🧩 TrailingStopManager** is an automated orchestrator designed to **PROTECT PROFITS** on open positions. It continuously monitors all positions and automatically moves stop-loss orders as price moves favorably, locking in gains without manual intervention. Perfect for trending markets where you want to let winners run while protecting accumulated profits.

**Strategy Principle**: Once a position reaches a specified profit threshold, the manager activates a trailing stop mechanism:

- **BUY positions**: Stop loss moves UPWARD only as price rises (never down)
- **SELL positions**: Stop loss moves DOWNWARD only as price falls (never up)
- Maintains fixed distance from current price (e.g., 200 points behind)
- Locks in profits automatically during reversals
- Works 24/7 managing multiple positions simultaneously
- No manual chart watching required

**File**: `examples/demos/orchestrators/11_trailing_stop.go`

> **IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY**
>
> This orchestrator is a demonstration showing how GoMT5 methods function and combine into automated strategies. **NOT a production-ready trading strategy!** Always test on demo accounts first.

---

## ⚙️ Architecture

```
TRAILING STOP MANAGER ORCHESTRATOR
    |
MT5Sugar Instance
    |
    |  |  |
GetOpenPositions  GetPriceInfo  ModifyPositionSL
(monitoring)      (tracking)     (adjusting)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetOpenPositions`, `GetPriceInfo`, `ModifyPositionSL`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo` for position data structures

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `TrailingDistance` | float64 | `200` | Distance in points to trail behind current price |
| `ActivationProfit` | float64 | `300` | Profit in points needed to activate trailing |
| `UpdateInterval` | time.Duration | `2s` | How often to check and update positions |
| `Symbols` | []string | `[]` | Symbols to manage (empty = all symbols) |
| `MinDistance` | float64 | `100` | Declared "minimum distance from current price", but **not read by the current code** (no check uses it) |
| `StepSize` | float64 | `50` | Minimum step size for SL adjustments (avoid excessive mods) |

> **Point size:** all distances are in points, and the manager does not query the symbol's point size: `getSymbolPoint()` caches a hard-coded `0.00001` (5-digit symbol) for every symbol. The distances are therefore only correct for 5-digit symbols such as EURUSD.

### Configuration Example

```go
config := orchestrators.TrailingStopConfig{
    TrailingDistance: 200,           // Tail 200 points behind price
    ActivationProfit: 300,           // Activate after +300 points profit
    UpdateInterval:   2 * time.Second, // Check every 2 seconds
    Symbols:          []string{},    // Manage all symbols
    MinDistance:      100,           // Min 100pts from current price
    StepSize:         50,            // Adjust in 50pt increments
}

tsManager := orchestrators.NewTrailingStopManager(sugar, config)
```

---

## How to Run

You can execute this orchestrator using several command variations:

```bash
# Navigate to examples directory
cd examples/demos

# Option 1: By number
go run main.go 11

# Option 2: By keyword
go run main.go trailing
```

Both commands will launch the **TrailingStopManager** orchestrator with the configuration specified in `main.go -> RunOrchestrator_TrailingStop()`. The runner connects with `mt5.NewMT5Sugar(...)` + `QuickConnect(...)`, runs the manager for 3 minutes (180 s) with a live progress line, then calls `Stop()` and prints the final metrics.

---

## Algorithm

### 🟢 Flowchart

```
START
  |
Initialize TrailingStopManager
  |
Create empty tracker map for positions
  |
MONITORING LOOP (every UpdateInterval):
  | Get all open positions
  |
  FOR each position:
    | Check if symbol is tracked (if Symbols list specified)
    | Get or create position tracker
    | Get current price (Bid for BUY, Ask for SELL)
    | Calculate profit in points
    |
    | Trailing not active yet and profit >= ActivationProfit?
      -> YES: activate trailing (stays active for this position)
    | Trailing active?
      -> NO                       -> YES (same check, no extra cycle)
      Skip (not active)           Calculate new SL:
                                     | BUY:  newSL = currentPrice - TrailingDistance*point
                                     | SELL: newSL = currentPrice + TrailingDistance*point
                                     |
                                  Is new SL better than current SL (or current SL = 0)?
                                     | BUY:  newSL > currentSL?
                                     | SELL: newSL < currentSL?
                                     |
                                  YES -> Move >= StepSize*point (or current SL = 0)?
                                           YES -> ModifyPositionSL
                                           NO  -> Skip
                                  NO  -> Skip (don't move SL unfavorably)
  |
Clean up trackers for closed positions
  |
REPEAT until stopped
  |
ON STOP:
  | Cancel monitoring loop
  | Clear all trackers
  |
END
```

### 🛠️ Step-by-step Description

### 1. Initialization

> **Source**: `examples/demos/orchestrators/11_trailing_stop.go`  
> **Lines**: `169–179`


```go
func NewTrailingStopManager(sugar *mt5.MT5Sugar, config TrailingStopConfig) *TrailingStopManager {
    return &TrailingStopManager{
        BaseOrchestrator: NewBaseOrchestrator("Trailing Stop Manager"),
        sugar:            sugar,
        config:           config,
        trackedPositions: make(map[uint64]*positionTracker),
        symbolDigits:     make(map[string]int),
        symbolPoints:     make(map[string]float64),
    }
}
```

**What happens**:

- Creates new manager instance with provided configuration
- Initializes empty tracker maps for positions and symbol parameters
- Each position gets a dedicated tracker storing its state (its first-seen SL is taken from `pos.StopLoss`)
- The symbol maps are a cache filled by `getSymbolPoint()`, which stores a hard-coded point of `0.00001` (5 digits) for every symbol - it does not query the symbol

#### 2. Starting the Manager

> Lines `182–198`  


```go
func (t *TrailingStopManager) Start() error {
    if t.IsRunning() {
        return fmt.Errorf("trailing stop manager already running")
    }

    // Create context
    ctx, cancel := context.WithCancel(context.Background())
    t.SetContext(ctx, cancel)

    // Mark as started
    t.MarkStarted()

    // Start monitoring loop
    go t.monitorLoop()

    return nil
}
```

**What happens**:

- Checks if already running (prevents duplicate starts)
- Creates cancellable context for graceful shutdown
- Marks orchestrator as started in metrics
- Launches monitoring loop in separate goroutine
- Returns immediately (non-blocking)

#### 3. Monitoring Loop

> Lines `218–231`  


```go
func (t *TrailingStopManager) monitorLoop() {
    ticker := time.NewTicker(t.config.UpdateInterval)
    defer ticker.Stop()

    for {
        select {
        case <-t.GetContext().Done():
            return
        case <-ticker.C:
            t.updateAllTrailingStops()
        }
    }
}
```

**Monitoring actions every UpdateInterval**:

1. **Ticker triggers**: Every 2 seconds by default
2. **Context check**: Monitors for stop signal
3. **Update all**: Calls `updateAllTrailingStops()` to process positions
4. **Graceful exit**: Stops cleanly when context is cancelled

---

## Visual Example

### Scenario: BUY EURUSD Position with Trailing

**Configuration**:
- Entry: 1.10000
- TrailingDistance: 200 points
- ActivationProfit: 300 points
- StepSize: 50 points

```
TIME 0: Position Opened

Price:        1.10000 (Entry)
SL:           0 (no stop loss yet)
Profit:       0 points
Status:       Waiting for activation (need +300pts)

TIME 1: Price Rises to 1.10300 (Activation + First SL Update)

Price:        1.10300
Profit:       +300 points
Status:       TRAILING ACTIVATED!
Calculation:  newSL = 1.10300 - 0.00200 = 1.10100  (same check)
SL:           0 -> 1.10100  UPDATED
Locked:       +100 points guaranteed

TIME 2: Price at 1.10400

Price:        1.10400
Calculation:  newSL = 1.10400 - 0.00200 = 1.10200
Step check:   1.10200 - 1.10100 = 100 points >= StepSize 50
SL:           1.10100 -> 1.10200  UPDATED
Profit:       +400 points
Locked:       +200 points guaranteed

TIME 3: Price at 1.10500

Price:        1.10500
Calculation:  newSL = 1.10500 - 0.00200 = 1.10300
SL:           1.10200 -> 1.10300  UPDATED
Profit:       +500 points
Locked:       +300 points guaranteed

TIME 4: Price Drops to 1.10450

Price:        1.10450
Calculation:  newSL = 1.10450 - 0.00200 = 1.10250
Check:        1.10250 > 1.10300? NO
SL:           1.10300 -> unchanged (SL only moves UP for BUY)
Profit:       +450 points (floating)
Locked:       +300 points guaranteed

TIME 5: Price Reverses to 1.10300

Price:        1.10300
Event:        STOP LOSS TRIGGERED @ 1.10300
Position:     CLOSED
Final Profit: +300 points = $3 (for 0.01 lots EURUSD)
Result:       Profit protected from reversal!

WITHOUT TRAILING STOP:

Price reached 1.10500 (+500pts)
Then reversed to 1.10000 (entry)
Final result: $0 profit (or worse if continues down)

WITH TRAILING STOP:

Price reached 1.10500 (+500pts)
Trailing activated, SL moved to 1.10300
Price reversed, stopped out at 1.10300
Final result: +300pts = $3 profit LOCKED IN (0.01 lots)
```

---

## 💡 Used MT5Sugar Methods

### 1. GetOpenPositions

```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all currently open positions across all symbols.

**Usage in TrailingStopManager**:
- Called every UpdateInterval (default 2 seconds)
- Returns full list of positions with: ticket, symbol, type, open price, current SL/TP, profit
- Used to identify which positions need trailing stop management

**Returns**: Slice of position info or error

### 2. GetPriceInfo

```go
func (s *MT5Sugar) GetPriceInfo(symbol string) (*PriceInfo, error)
```

**Purpose**: Gets current Bid/Ask prices for a symbol.

**Usage in TrailingStopManager**:
- Called for each position to get current price
- BUY positions use Bid price (exit price)
- SELL positions use Ask price (exit price)
- Used to calculate profit in points and new trailing SL

**Returns**: `mt5.PriceInfo{Symbol, Bid, Ask, SpreadPips, Time}` (3 s timeout)

### 3. ModifyPositionSL

```go
func (s *MT5Sugar) ModifyPositionSL(ticket uint64, sl float64) error
```

**Purpose**: Modifies only the stop loss for an existing position: it sends an `OrderModifyRequest` with only `Ticket` and `StopLoss` set (`TakeProfit` is not sent). 10 s timeout; error unless the return code is 10009.

**Usage in TrailingStopManager**:
- Called when new calculated SL is more favorable than current
- Updates stop loss on MT5 server
- Does not send a take profit value
- Core method for implementing trailing logic

**Parameters**:
- `ticket`: Position ticket number
- `sl`: New stop loss price (sent as given, not rounded)

---

## Risk Management

### Protection Benefits

**TrailingStopManager provides**:

1. **Profit Protection**:
   - Locks in gains automatically
   - Prevents giving back large profits
   - No manual monitoring required

2. **Psychological Benefit**:
   - Removes emotion from profit-taking
   - Enforces discipline on winning trades
   - Lets winners run without fear

3. **Risk Reduction**:
   - Converts open profit to guaranteed minimum
   - Reduces risk as position becomes profitable
   - Provides downside protection on winning trades

### Risk Considerations

1. **Premature Exits**:
   - Too tight TrailingDistance -> stopped out on normal volatility
   - **Solution**: Match TrailingDistance to symbol's ATR
   - Example: EURUSD 200pts (note: the point size is hard-coded to 0.00001, so non-5-digit symbols such as XAUUSD need a code change first)

2. **Late Activation**:
   - Too high ActivationProfit -> miss protection on smaller winners
   - **Solution**: Set based on average winning trade size
   - Example: Scalping 100pts, Swing 300pts

3. **Excessive Modifications**:
   - Too small StepSize -> many server calls, slippage
   - **Solution**: Set StepSize to 25-50% of TrailingDistance
   - Example: TrailingDistance 200pts -> StepSize 50pts

4. **Gap Risk**:
   - Stop loss can be triggered at worse price during gaps
   - Trailing doesn't protect against weekend gaps
   - **Solution**: Close positions before high-risk events

---

## Configuration Profiles

### Conservative (Wider Safety)

```go
TrailingDistance: 200,   // Wide buffer to avoid premature exits
ActivationProfit: 300,   // Wait for solid profit before activating
UpdateInterval:   2 * time.Second,
```

### Aggressive (Tighter Trailing)

```go
TrailingDistance: 100,   // Closer trailing = higher risk/reward
ActivationProfit: 150,   // Activate sooner
UpdateInterval:   1 * time.Second,        // More frequent checks
```

### Scalping (Very Tight)

```go
TrailingDistance: 50,    // Very tight trailing for quick exits
ActivationProfit: 100,   // Activate on small profits
UpdateInterval:   500 * time.Millisecond, // Near real-time monitoring
```

---

## When to Use Trailing Stops

### Suitable Conditions

- **Trending markets**: Strong directional movement (uptrend for BUY, downtrend for SELL)
- **Breakout trades**: After level breaks, let winners run
- **Momentum strategies**: When price moves quickly in your favor
- **Swing trading**: Multi-hour or multi-day positions
- **Hands-off management**: When you can't watch charts constantly
- **Multiple positions**: Managing several trades simultaneously
- **Profit protection**: Want to lock in gains on winning trades

### Unsuitable Conditions

- **Range-bound markets**: Sideways price action will trigger stops prematurely
- **High volatility**: Wide swings can stop you out on normal fluctuations
- **News events**: Spikes can trigger stops before reversal
- **Scalping with fixed targets**: Better to use fixed TP for quick exits
- **Mean reversion strategies**: Trailing fights the core strategy
- **Very tight stops**: Sub-20 pip stops get triggered by noise

---

## Configuration Location

**Runtime configuration is located in**: `examples/demos/main.go -> RunOrchestrator_TrailingStop()`

This separation exists because:

1. **Code Reusability**: Same orchestrator runs with different parameters
2. **Quick Testing**: Change numbers without modifying strategy logic
3. **User Examples**: Shows how to properly configure the orchestrator
4. **Centralized Entry Point**: Single place for all orchestrator launches

