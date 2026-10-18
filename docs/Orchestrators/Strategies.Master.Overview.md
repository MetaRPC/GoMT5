# Orchestrators & Presets — Complete Overview

## Introduction

 **Orchestrators** are advanced trading strategies that combine multiple Go MT5 methods into sophisticated automated systems.  
Think of them as *trading robots* that handle complex scenarios like trailing stops, position scaling, grid trading, risk management, and portfolio balancing.

🧠 **Presets** are intelligent systems that combine multiple orchestrators and automatically switch between them based on market conditions.

All of them are example code in the repository's `examples` module (get it by cloning the repository): orchestrators in `examples/demos/orchestrators` (import `"github.com/MetaRPC/GoMT5/examples/demos/orchestrators"`), the preset in `examples/demos/presets` (import `"github.com/MetaRPC/GoMT5/examples/demos/presets"`), built on `mt5 "github.com/MetaRPC/GoMT5/examples/mt5"` (MT5Sugar).


---

## Quick Navigation

| Orchestrator           | Purpose                                    | Jump to Docs |
|------------------------|--------------------------------------------|--------------|
| **TrailingStopManager** | Protect profits by auto-moving stop-loss   | [Details](#orchestrator-11-trailingstopmanager) · [Full Docs](11_Trailing_stop.md) |
| **PositionScaler**      | Scale into winning/losing positions        | [Details](#orchestrator-12-positionscaler) · [Full Docs](12_Position_scaler.md) |
| **GridTrader**          | Profit from range-bound markets            | [Details](#orchestrator-13-gridtrader) · [Full Docs](13_Grid_trader.md) |
| **RiskManager**         | Monitor and enforce account safety limits  | [Details](#orchestrator-14-riskmanager) · [Full Docs](14_Risk_manager.md) |
| **PortfolioRebalancer** | Maintain balanced multi-symbol allocation  | [Details](#orchestrator-15-portfoliorebalancer) · [Full Docs](15_Portfolio_rebalancer.md) |
| **Adaptive Preset**    | Auto-switch strategies based on conditions | [Details](#preset-16-adaptive-orchestrator-preset) · [Full Docs](16_AdaptiveOrchestratorPreset.md) |


---

## Architecture Hierarchy

```
MT5 Account (Low-level gRPC)
|
MT5 Service (Go wrappers)
|
MT5 Sugar (High-level convenience)
|
ORCHESTRATORS (Automated strategies)
|
+-- TrailingStopManager (Profit protection)
+-- PositionScaler (Position sizing)
+-- GridTrader (Range-bound trading)
+-- RiskManager (Account protection)
+-- PortfolioRebalancer (Multi-symbol management)
|
PRESETS (Multi-strategy systems)
|
+-- Adaptive Orchestrator Preset (Intelligent mode switching)
```

---

## Quick Comparison Table

| Orchestrator           | Purpose                 | Best For         | Complexity | Risk Level            |
| ---------------------- | ----------------------- | ---------------- | ---------- | --------------------- |
| **TrailingStopManager** | Protect profits         | Trending markets | Low        | Low                   |
| **PositionScaler**      | Add to winners/losers   | Active trading   | Medium     | Medium–High           |
| **GridTrader**          | Range-bound profits     | Sideways markets | Medium     | Medium                |
| **RiskManager**         | Account protection      | All conditions   | Low        | Very Low (protective) |
| **PortfolioRebalancer** | Multi-symbol balance    | Diversification  | Medium     | Low–Medium            |
| **Adaptive Preset**    | Auto strategy switching | All conditions   | High       | Medium                |

---

## Orchestrator #11: TrailingStopManager

### Overview

Automatically moves stop-loss orders to protect profits as price moves favorably.

### Key Features

* Activates after profit threshold is reached
* Trails stop-loss behind price (never moves unfavorably)
* Handles multiple positions simultaneously
* Runs in the background until `Stop()` is called (or Ctrl+C)
* Distances in points of each symbol (point size read with `GetSymbolInfo`); the SL is never closer than `MinDistance` or the broker's stops level

### Configuration Example

```go
config := orchestrators.DefaultTrailingStopConfig()
config.TrailingDistance = 200 // Trail 200 points behind
config.ActivationProfit = 300 // Activate at 300 points profit
config.UpdateInterval = 2 * time.Second

manager := orchestrators.NewTrailingStopManager(sugar, config)
manager.Start()
```

### When to Use

✅ **IDEAL FOR**:

* Trending markets where you want to let winners run
* Protecting accumulated profits automatically
* "Set-and-forget" profit protection

⚠️ **AVOID**:

* Ranging or choppy markets (may exit prematurely)
* When you prefer manual stop management


### Quick Stats

* **Input**: Open positions
* **Action**: Modifies stop-loss levels
* **Output**: Protected profits
* **Cycle Time**: 2 seconds (configurable)

[Full Documentation](11_Trailing_stop.md)


---

## Orchestrator #12: PositionScaler

### Overview

Adds to winning or losing positions based on the configured mode (Pyramiding, Averaging Down, Scale Out).

### Key Features

* **3 Modes**:

  1. **Pyramiding** — add to winners (trend trading)
  2. **Averaging Down** — add to losers (high risk!)
  3. **Scale Out** — gradually exit positions
* Configurable trigger distances
* Maximum scale limits
* Fixed scale lot size, optionally reduced per scale (`ReducePerScale`)
* Opens a BUY base position of `InitialLotSize` for a listed symbol without a position (once per run; 0 = only existing positions); scale-outs are `TriggerDistance` apart and close part of the largest position

### Configuration Example

```go
config := orchestrators.DefaultPositionScalerConfig("EURUSD") // Symbols: []string{"EURUSD"}
config.Mode = orchestrators.Pyramiding // Add to winners
config.TriggerDistance = 200           // Add every 200 points
config.MaxScales = 3                   // Max 3 additions
config.ScaleLotSize = 0.01

scaler := orchestrators.NewPositionScaler(sugar, config)
scaler.Start()
```

### When to Use

**✅ IDEAL FOR**:

* **Pyramiding** — strong trending markets
* **Averaging Down** — range-bound markets (risky!)
* **Scale Out** — taking profits gradually

**⚠️ AVOID**:

* Averaging down in trending markets (disaster!)
* Small accounts (need margin for scaling)

### Quick Stats

* **Input**: Open position + price movement
* **Action**: Opens additional positions (ScaleOut mode: partially closes)
* **Output**: Scaled exposure
* **Cycle Time**: 5 seconds (configurable)

[Full Documentation](12_Position_scaler.md)

---

## Orchestrator #13: GridTrader

### Overview

📊 Places a grid of buy/sell orders at regular intervals, profiting from price oscillations in range-bound markets.

### Key Features

* Configurable grid size and spacing
* Automatic BUY LIMIT / SELL LIMIT placement, every order with a TP (one grid step, or `TakeProfit`), prices rounded to the symbol's digits
* Works best in sideways markets
* Grid is rebuilt when price drifts more than 2 grid steps from the center, and after a fill with `RebuildOnFill`
* `MaxPositions` is enforced (no more pending orders than the room left; at the limit they are deleted)

### Configuration Example

```go
config := orchestrators.DefaultGridTraderConfig("EURUSD")
config.GridSize = 5     // 5 grid levels each side
config.GridStep = 100   // 100 points spacing
config.LotSize = 0.01
config.MaxPositions = 10
config.TakeProfit = 100 // 0 = one grid step
config.CheckInterval = 5 * time.Second

gridTrader := orchestrators.NewGridTrader(sugar, config)
gridTrader.Start()
```

### When to Use

**✅ IDEAL FOR**:

* Range-bound, sideways markets
* Low volatility conditions
* Markets with clear support/resistance

**⚠️ AVOID**:

* Strong trending markets (losses accumulate)
* High volatility (grid gets overrun)
* News events or breakouts

### Quick Stats

* **Input**: Price range
* **Action**: Places grid of orders
* **Output**: Small frequent profits
* **Cycle Time**: 5 seconds (configurable)

[Full Documentation](13_Grid_trader.md)

---

## Orchestrator #14: RiskManager

### Overview

Monitors account-level risk metrics while running and automatically intervenes when danger limits are exceeded.

### Key Features

* **5 Protection Types**:

  1. Drawdown protection (percentage and absolute)
  2. Daily loss limits
  3. Daily profit targets
  4. Margin level monitoring
  5. Position count limits
* Emergency close of all positions (`CloseAllPositions`)
* Trade blocking after the daily limits (and while margin usage > `MaxMarginUsed`): positions opened while blocked are closed again; check `IsTradingBlocked()` in your code
* One event per breach; after an emergency close the drawdown peak is reset; a limit of 0 is off
* Daily automatic reset


### Configuration Example

```go
// Start from the defaults (in a literal, a missing limit is 0 = off)
config := orchestrators.DefaultRiskManagerConfig()
config.MaxDrawdownPercent = 10.0   // 10% max drawdown
config.MaxDrawdownAbsolute = 1000.0 // $1000 max loss
config.DailyLossLimit = 500.0      // $500 daily limit
config.DailyProfitTarget = 1000.0  // $1000 daily goal
config.MinMarginLevel = 150.0      // 150% minimum margin
config.EnableAutoClose = true      // Auto-close on breach
config.CheckInterval = 5 * time.Second

riskManager := orchestrators.NewRiskManager(sugar, config)
riskManager.Start()
```

### When to Use

**✅ IDEAL FOR**:

* All trading (always use risk management!)
* Automated systems running unattended
* Learning proper risk control
* Protecting against blow-ups

**⚠️ AVOID**:

* Never avoid risk management
* (Seriously, always use it)

### Quick Stats

* **Input**: Account metrics (equity, balance, margin)
* **Action**: Monitors, alerts, closes positions if needed
* **Output**: Protected account
* **Cycle Time**: 5 seconds (configurable)

[Full Documentation](14_Risk_manager.md)

---

## Orchestrator #15: PortfolioRebalancer

### Overview

Maintains balanced exposure across multiple trading symbols by automatically rebalancing positions to match target allocations.

### Key Features

* Equal-weight or custom allocation percentages
* Automatic rebalancing when deviation exceeds threshold
* Multi-symbol portfolio management
* Exposure = net notional value in the account currency (volume × contract size × price, BUY +, SELL −); it resizes existing positions in their direction (never opens a symbol without a position) and reduces with partial closes

### Configuration Example

```go
// Equal weights (0.25 each), MinPositionSize 0.01, MaxTradesPerCycle 10, ...
config := orchestrators.DefaultPortfolioRebalancerConfig([]string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"})
config.TotalExposure = 10000.0     // $10,000 notional ($2,500 per symbol)
config.RebalanceThreshold = 10.0   // Rebalance at 10% off
config.CheckInterval = 30 * time.Minute

rebalancer := orchestrators.NewPortfolioRebalancer(sugar, config)
rebalancer.Start()
```

### When to Use

**✅ IDEAL FOR**:

* Multi-symbol trading systems
* Diversification strategies
* Portfolio managers
* Long-term position holding

**⚠️ AVOID**:

* Single-symbol trading
* High-frequency trading
* Very small accounts (< $1,000)

### Quick Stats

* **Input**: Multiple symbol positions
* **Action**: Opens/closes positions to restore balance
* **Output**: Balanced portfolio allocation
* **Cycle Time**: 30 minutes in main.go (default 1 hour, configurable)

[Full Documentation](15_Portfolio_rebalancer.md)

---

## Preset #16: Adaptive Orchestrator Preset

### Overview

Intelligent multi-strategy system that analyzes market conditions every 5 minutes and automatically selects the most appropriate orchestrator to run.

### Key Features

* **5 Market Modes**:

  1. **Grid Mode** — low volatility → GridTrader
  2. **Managed Mode** — medium volatility → TrailingStop + PositionScaler
  3. **Protection Mode** — high volatility → RiskManager
  4. **Portfolio Mode** — multi-symbol (only with `EnablePortfolioMode = true`) → PortfolioRebalancer
  5. **Trending Mode** — strong trend → PositionScaler (pyramiding)
* Automatic mode switching every cycle
* Safety limits (total loss, daily loss/profit) checked after each cycle; a safety stop ends `Execute()` like Ctrl+C
* `MaxConcurrentOrchestrators` caps the orchestrators started per cycle
* Volatility and trend are simplified placeholders (spread-based)


### Configuration Example

```go
preset := presets.NewAdaptiveOrchestratorPreset(sugar)

// Customize
preset.Symbol = "EURUSD"
preset.LowVolatilityThreshold = 25.0  // Grid below 25 pts
preset.HighVolatilityThreshold = 50.0 // Protection above 50 pts
preset.CycleDuration = 5 * time.Minute
preset.MaxDailyLoss = 500.0
preset.MaxDailyProfit = 1000.0

// Run (blocks until Ctrl+C or a safety stop)
totalProfit, _ := preset.Execute()
```

### When to Use

**✅ IDEAL FOR**:

* "Set-and-forget" adaptive trading
* Testing multiple strategies
* 24/7 automated trading
* Learning orchestrator composition

**⚠️ AVOID**:

* Very specific proven strategy
* High-frequency trading
* Small accounts (< $500)
* Major news events

### Quick Stats

* **Input**: Market conditions (volatility, trend, positions)
* **Action**: Selects and runs the appropriate orchestrator
* **Output**: Adaptive multi-strategy performance
* **Cycle Time**: 5 minutes + 30-second pause

[Full Documentation](16_AdaptiveOrchestratorPreset.md)

---

## Orchestrator Selection Guide

### By Market Condition

| Market Condition         | Best Orchestrator                         | Why                             |
| ------------------------ | ----------------------------------------- | ------------------------------- |
| **Strong Uptrend**       | TrailingStop + PositionScaler (Pyramiding) | Protect profits, add to winners |
| **Strong Downtrend**     | TrailingStop (short positions)            | Protect short profits           |
| **Sideways / Ranging**   | GridTrader                                 | Profit from oscillations        |
| **Low Volatility**       | GridTrader                                 | Small frequent profits          |
| **High Volatility**      | RiskManager                                | Protection mode                 |
| **Multi-Symbol Active**  | PortfolioRebalancer                        | Balance allocation              |
| **Uncertain / Changing** | Adaptive Preset                           | Auto-adapt to conditions        |

---

### By Trading Style

| Trading Style        | Recommended Orchestrators               |
| -------------------- | --------------------------------------- |
| **Scalping**         | TrailingStop (tight parameters)         |
| **Day Trading**      | TrailingStop + RiskManager               |
| **Swing Trading**    | PositionScaler (pyramiding) + RiskManager |
| **Position Trading** | PortfolioRebalancer + RiskManager         |
| **Grid Trading**     | GridTrader + RiskManager                  |
| **Automated 24/7**   | Adaptive Preset                         |

---

### By Risk Tolerance

| Risk Level          | Orchestrator Combination                    |
| ------------------- | ------------------------------------------- |
| **Conservative**    | RiskManager only (5% drawdown limit)         |
| **Moderate**        | TrailingStop + RiskManager                   |
| **Balanced**        | PortfolioRebalancer + RiskManager             |
| **Aggressive**      | PositionScaler + TrailingStop + RiskManager   |
| **Very Aggressive** | GridTrader + PositionScaler (NOT RECOMMENDED) |

---

## Common Orchestrator Combinations

### Combination 1: Safe Trend Trading

**Orchestrators**: TrailingStop + RiskManager

**Purpose**: Follow trends with profit protection and account safety

```go
// Trailing stop for profit protection
tsConfig := orchestrators.DefaultTrailingStopConfig()
trailingManager := orchestrators.NewTrailingStopManager(sugar, tsConfig)
trailingManager.Start()

// Risk management for account protection
riskConfig := orchestrators.DefaultRiskManagerConfig()
riskManager := orchestrators.NewRiskManager(sugar, riskConfig)
riskManager.Start()
```


**Best For**: Beginners, conservative traders, trending markets

---

### Combination 2: Aggressive Scaling

**Orchestrators**: PositionScaler + TrailingStop + RiskManager

**Purpose**: Scale into trends, protect profits, limit risk

```go
// Scale for pyramiding
scaleConfig := orchestrators.DefaultPositionScalerConfig("EURUSD")
scaleConfig.Mode = orchestrators.Pyramiding
scaler := orchestrators.NewPositionScaler(sugar, scaleConfig)
scaler.Start()

// Trailing stop for profit protection
tsConfig := orchestrators.DefaultTrailingStopConfig()
tsManager := orchestrators.NewTrailingStopManager(sugar, tsConfig)
tsManager.Start()

// Risk management for safety
riskConfig := orchestrators.DefaultRiskManagerConfig()
riskManager := orchestrators.NewRiskManager(sugar, riskConfig)
riskManager.Start()
```

**Best For**: Experienced traders, strong trends, large accounts

---

### Combination 3: Portfolio Diversification

**Orchestrators**: PortfolioRebalancer + RiskManager

**Purpose**: Maintain balanced multi-symbol portfolio with protection

```go
// Portfolio rebalance (defaults: equal weights, TotalExposure 10000, CheckInterval 1h - first check at start)
portfolioConfig := orchestrators.DefaultPortfolioRebalancerConfig([]string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"})

rebalancer := orchestrators.NewPortfolioRebalancer(sugar, portfolioConfig)
rebalancer.Start()

// Risk management for overall protection
riskConfig := orchestrators.DefaultRiskManagerConfig()
riskManager := orchestrators.NewRiskManager(sugar, riskConfig)
riskManager.Start()
```

**Best For**: Portfolio managers, diversification seekers, medium–large accounts

---

### Combination 4: Range Trading

**Orchestrators**: GridTrader + RiskManager

**Purpose**: Profit from sideways markets with risk control

```go
// Grid trade (defaults: 5 levels, 100 points, 0.01 lots, MaxPositions 10, CheckInterval 5s)
gridConfig := orchestrators.DefaultGridTraderConfig("EURUSD")

gridTrader := orchestrators.NewGridTrader(sugar, gridConfig)
gridTrader.Start()

// Risk management (important for grid!)
riskConfig := orchestrators.DefaultRiskManagerConfig()
riskConfig.MaxDrawdownPercent = 15.0 // Wider for grid
riskConfig.EnableAutoClose = true

riskManager := orchestrators.NewRiskManager(sugar, riskConfig)
riskManager.Start()
```

**Best For**: Range-bound markets, low volatility, experienced traders

---

## Running Orchestrators

### Command-Line Execution

```bash
cd examples/demos

# Individual orchestrators
go run main.go 11  # TrailingStop
go run main.go trailing

go run main.go 12  # PositionScaler
go run main.go scaler

go run main.go 13  # GridTrader
go run main.go grid

go run main.go 14  # RiskManager
go run main.go risk

go run main.go 15  # PortfolioRebalancer
go run main.go rebalancer

# Adaptive preset
go run main.go 16       # Adaptive Preset
go run main.go adaptive
go run main.go preset
```

---

### Programmatic Usage

```go
// Initialize connection
sugar, err := mt5.NewMT5Sugar(12345678, "your_password", "mt5.mrpc.pro:443")
if err != nil {
    log.Fatal(err)
}
// Disconnect also when QuickConnect fails (the server may have started the terminal)
defer sugar.Disconnect()
if err := sugar.QuickConnect("YourBroker-MT5 Demo"); err != nil {
    log.Println(err) // not log.Fatal: it would skip the deferred Disconnect
    return
}

// Create orchestrator with config
config := orchestrators.DefaultTrailingStopConfig()
orchestrator := orchestrators.NewTrailingStopManager(sugar, config)

// Start (non-blocking; an invalid config is an error)
if err := orchestrator.Start(); err != nil {
    log.Println(err)
    return
}

// Monitor status
status := orchestrator.GetStatus()
fmt.Printf("Running: %v, Operations: %d\n",
    status.IsRunning,
    status.SuccessCount+status.ErrorCount)

// Stop when done
orchestrator.Stop()
```


---

## ✅ Best Practices

### DO

**Always Use RiskManager**

* Combine **any** orchestrator with RiskManager
* Set conservative limits initially
* Monitor risk events regularly

**Test on Demo First**

* Run each orchestrator on demo for a minimum of 1 week
* Verify behavior matches expectations
* Adjust parameters based on observations

**Start Conservative**

* Use small position sizes initially
* Set tighter safety limits
* Gradually loosen as you gain confidence

**Monitor Actively at First**

* Watch first 10 cycles closely
* Verify orchestrator behavior
* Check for unexpected actions

**Combine Intelligently**

* TrailingStop + RiskManager = safe
* GridTrader + PositionScaler = dangerous!
* Portfolio + RiskManager = smart diversification

---

### ⚠️ DON'T

**Never Trade Without Risk Management**

* Even on demo, practice with RiskManager
* Learn proper risk control habits
* Safety should be automatic

**Don't Run Conflicting Orchestrators**

* GridTrader + Trending Mode scaling = conflict
* Test combinations on demo first
* Understand how they interact

**Don't Set-and-Forget Initially**

* First week: monitor closely
* Second week: check daily
* Third week+: periodic review

**Don't Use Production Strategies on Live Immediately**

* Test all parameters on demo
* Verify calculations are correct
* Ensure orchestrators behave as expected

**Don't Ignore Warnings**

* All orchestrators have risk warnings
* Read "Critical Warnings" sections
* Understand limitations before using

---

## Troubleshooting Common Issues

### Problem: Orchestrator Doesn't Start

**Possible Causes**:

* Already running (`"... already running"`)
* GridTrader: no price for the symbol (`failed to initialize symbol: ...`)
* RiskManager: balance cannot be read (`failed to initialize: ...`)
* PortfolioRebalancer: allocations do not sum to 100% (`invalid configuration: ...`)

**Solution**:

```go
if err := orchestrator.Start(); err != nil {
    fmt.Printf("Start failed: %v\n", err)
    // Check error details
}
```

---

### Problem: Orchestrator Stops Unexpectedly

**Possible Causes**:

* `Stop()` was called — the runners in `main.go` stop each orchestrator after a fixed time (3, 5, 10, 15 or 60 minutes)
* Errors do **not** stop an orchestrator: they are counted in `ErrorCount` / `LastError` and the loop continues
* An invalid configuration makes `Start()` return an `invalid configuration` error (a zero `CheckInterval` / `UpdateInterval` gets its default)
* Ctrl+C cancels the orchestrator's context (`GetContext()`); the runners in `main.go` disconnect the terminal however the run ends

**Solution**:

```go
status := orchestrator.GetStatus()
fmt.Printf("Status: %+v\n", status)
// Review ErrorCount / LastError, and GetMetrics().LastOperation
```

---

### Problem: Too Many Operations / Trades

**Possible Causes**:

* Check interval too short
* Trigger thresholds too sensitive
* Market too volatile for settings

**Solution**:

* Increase `CheckInterval` (e.g. 5s → 10s)
* Widen trigger thresholds
* Review configuration parameters

---

### Problem: Not Enough Operations

**Possible Causes**:

* Thresholds too strict
* No positions meeting criteria
* Wrong market conditions for strategy

**Solution**:

* Lower activation thresholds
* Verify positions exist
* Check if strategy matches market

---

### Problem: Orchestrators Conflict

**Possible Causes**:

* Both modifying the same positions
* Incompatible strategies (grid + pyramiding)

**Solution**:

* Run only compatible combinations
* Test on demo first
* Use Adaptive Preset instead (switches modes automatically)

---

## Performance Monitoring

### Checking Orchestrator Status

```go
// Get status and metrics
status := orchestrator.GetStatus()
metrics := orchestrator.GetMetrics()

// Display them
fmt.Printf("Name: %s\n", status.Name)
fmt.Printf("Running: %v\n", status.IsRunning)
fmt.Printf("SuccessCount: %d\n", status.SuccessCount)
fmt.Printf("ErrorCount: %d\n", status.ErrorCount)
fmt.Printf("TotalOperations: %d\n", status.SuccessCount+status.ErrorCount)
fmt.Printf("LastError: %s\n", status.LastError)
fmt.Printf("LastOperation: %s\n", metrics.LastOperation)
fmt.Printf("Uptime: %s\n", orchestrators.FormatDuration(status.Uptime))
```

`OrchestratorStatus` fields: `Name, IsRunning, StartTime, LastUpdate, ErrorCount, SuccessCount, LastError, Uptime`. `OrchestratorMetrics` holds trade/profit counters (`TotalTrades`, `NetProfit`, `CurrentPositions`, `WinRate`, ...) and `LastOperation`. Note: the RiskManager counts every risk event as an error.

### Metrics to Track

| Metric                 | What to Monitor         | Good Range             |
| ---------------------- | ----------------------- | ---------------------- |
| **Success Rate**       | SuccessCount / TotalOps | > 90%                  |
| **Error Rate**         | ErrorCount / TotalOps   | < 10%                  |
| **Operations / Hour**  | TotalOps / Hours        | Varies by orchestrator |
| **Profit / Operation** | Total P/L / TotalOps    | Positive average       |

---

## Advanced Topics

### Custom Orchestrators

To build your own orchestrator:

1. **Extend BaseOrchestrator** (embedding `*orchestrators.BaseOrchestrator`, created with `orchestrators.NewBaseOrchestrator(name)`, provides `IsRunning`, `GetStatus`, `GetMetrics`, `MarkStarted/MarkStopped`, `SetContext/GetContext/CancelContext`, `IncrementSuccess/IncrementError`, `UpdateStatus/UpdateMetrics`):

```go
type MyOrchestrator struct {
    *orchestrators.BaseOrchestrator
    sugar  *mt5.MT5Sugar
    config MyConfig
}
```


2. **Implement Required Methods**:

```go
func (m *MyOrchestrator) Start() error {
    // Validation
    // Setup
    // Launch goroutine
    return nil
}

func (m *MyOrchestrator) Stop() error {
    // Cleanup
    return nil
}
```

3. **Add Your Logic**:

```go
func (m *MyOrchestrator) executeStrategy() {
    // Your trading logic here
}
```

---

### Market Analysis Integration

Enhance orchestrators with your own analysis. GoMT5 (`MT5Account` / `MT5Service` / `MT5Sugar`) has no bar-history method, so indicators have to be computed from prices you collect (e.g. `GetPriceInfo` samples or `GetService().StreamTicks`) or from an external source:

```go
// Sample the Bid for one minute and measure its range
high, low := 0.0, 0.0
for i := 0; i < 60; i++ {
    if info, err := sugar.GetPriceInfo("EURUSD"); err == nil {
        if high == 0 || info.Bid > high {
            high = info.Bid
        }
        if low == 0 || info.Bid < low {
            low = info.Bid
        }
    }
    time.Sleep(time.Second)
}
info, err := sugar.GetSymbolInfo("EURUSD")
if err != nil {
    return err
}
rangePoints := (high - low) / info.Point // point size of the symbol

// Adjust orchestrator parameters based on analysis
config := orchestrators.DefaultTrailingStopConfig()
if rangePoints > 50 {
    config.TrailingDistance = 300 // Wider trailing in volatile markets
}
```

---

### External Signal Integration

Connect orchestrators to external signals:

```go
// Fetch signal from your system
signal := fetchTradingSignal()

if signal.Strength > 0.8 {
    // Start aggressive orchestrator
    scaler.Start()
} else {
    // Conservative approach
    trailingStop.Start()
}
```

---

## Summary

### Orchestrators Overview

| #  | Name               | Purpose            | Complexity |
| -- | ------------------ | ------------------ | ---------- |
| 11 | TrailingStopManager | Profit protection  | Low        |
| 12 | PositionScaler      | Position sizing    | Medium     |
| 13 | GridTrader          | Range trading      | Medium     |
| 14 | RiskManager         | Account protection | Low        |
| 15 | PortfolioRebalancer | Portfolio balance  | Medium     |
| 16 | Adaptive Preset    | Multi-strategy     | High       |

---

### Key Takeaways

**Always Use RiskManager**

* Combine with any orchestrator
* Essential for automated trading
* Prevents catastrophic losses

**Start Simple**

* Begin with TrailingStop + RiskManager
* Master one orchestrator before combining
* Test thoroughly on demo

**Match Strategy to Market**

* Grid for ranging markets
* Trailing for trends
* Adaptive for changing conditions

**Monitor and Adjust**

* Review performance regularly
* Optimize parameters based on results
* Don’t set-and-forget initially

**Understand Risks**

* Every orchestrator has warnings
* Read documentation fully
* Test calculations on demo


---

## Next Steps

1. **Read Individual Documentation** — Review detailed docs for orchestrators you want to use
2. **Run on Demo** — Test each orchestrator on a demo account for 1+ week
3. **Start with RiskManager** — Always combine with risk management
4. **Monitor Performance** — Track metrics and adjust parameters
5. **Graduate to Combinations** — Once comfortable, try multi-orchestrator setups
6. **Consider Adaptive Preset** — For hands-off adaptive trading

---

## Complete Documentation Links

* [TrailingStopManager — Full Docs](11_Trailing_stop.md)
* [PositionScaler — Full Docs](12_Position_scaler.md)
* [GridTrader — Full Docs](13_Grid_trader.md)
* [RiskManager — Full Docs](14_Risk_manager.md)
* [PortfolioRebalancer — Full Docs](15_Portfolio_rebalancer.md)
* [Adaptive Orchestrator Preset — Full Docs](16_AdaptiveOrchestratorPreset.md)


---

**Remember**: Orchestrators are powerful tools, but they require understanding and proper configuration. Always test on demo first, use risk management, and monitor performance actively until you are confident in their behavior.
