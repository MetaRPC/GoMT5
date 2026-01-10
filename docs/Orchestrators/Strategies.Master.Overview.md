# Orchestrators & Presets — Complete Overview

## Introduction

 **Orchestrators** are advanced trading strategies that combine multiple Go MT5 methods into sophisticated automated systems.  
Think of them as *trading robots* that handle complex scenarios like trailing stops, position scaling, grid trading, risk management, and portfolio balancing.

🧠 **Presets** are intelligent systems that combine multiple orchestrators and automatically switch between them based on market conditions.


---

## Quick Navigation

| Orchestrator           | Purpose                                    | Jump to Docs |
|------------------------|--------------------------------------------|--------------|
| **TrailingStopManage** | Protect profits by auto-moving stop-loss   | [Details](#orchestrator-11-trailingstopmanage) · [Full Docs](11_Trailing_stop.md) |
| **PositionScale**      | Scale into winning/losing positions        | [Details](#orchestrator-12-positionscale) · [Full Docs](12_Position_scaler.md) |
| **GridTrade**          | Profit from range-bound markets            | [Details](#orchestrator-13-gridtrade) · [Full Docs](13_Grid_trader.md) |
| **RiskManage**         | Monitor and enforce account safety limits  | [Details](#orchestrator-14-riskmanage) · [Full Docs](14_Risk_manager.md) |
| **PortfolioRebalance** | Maintain balanced multi-symbol allocation  | [Details](#orchestrator-15-portfoliorebalance) · [Full Docs](15_Portfolio_rebalancer.md) |
| **Adaptive Preset**    | Auto-switch strategies based on conditions | [Details](#preset-16-adaptiveorchestratorpreset) · [Full Docs](16_AdaptiveOrchestratorPreset.md) |


---

## Architecture Hierarchy

```
MT5 Account (Low-level gRPC)
|
MT5 Service (Go wrappers)
|
MT5 Suga (High-level convenience)
|
ORCHESTRATORS (Automated strategies)
|
+-- TrailingStopManage (Profit protection)
+-- PositionScale (Position sizing)
+-- GridTrade (Range-bound trading)
+-- RiskManage (Account protection)
+-- PortfolioRebalance (Multi-symbol management)
|
PRESETS (Multi-strategy systems)
|
+-- Adaptive Orchestrator Preset (Intelligent mode switching)
```

---

## Quick Comparison Table

| Orchestrator           | Purpose                 | Best For         | Complexity | Risk Level            |
| ---------------------- | ----------------------- | ---------------- | ---------- | --------------------- |
| **TrailingStopManage** | Protect profits         | Trending markets | Low        | Low                   |
| **PositionScale**      | Add to winners/losers   | Active trading   | Medium     | Medium–High           |
| **GridTrade**          | Range-bound profits     | Sideways markets | Medium     | Medium                |
| **RiskManage**         | Account protection      | All conditions   | Low        | Very Low (protective) |
| **PortfolioRebalance** | Multi-symbol balance    | Diversification  | Medium     | Low–Medium            |
| **Adaptive Preset**    | Auto strategy switching | All conditions   | High       | Medium                |

---

## Orchestrator #11: TrailingStopManage

### Overview

Automatically moves stop-loss orders to protect profits as price moves favorably.

### Key Features

* Activates after profit threshold is reached
* Trails stop-loss behind price (never moves unfavorably)
* Handles multiple positions simultaneously
* Works 24/7 without manual intervention

### Configuration Example

```go
config := orchestrators.TrailingStopConfig{
    TrailingDistance: 200, // Tail 200 points behind
    ActivationProfit: 300, // Activate at 300 points profit
    UpdateInterval:   2 * time.Second,
}

manager := orchestrators.NewTrailingStopManage(suga, config)
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

## Orchestrator #12: PositionScale

### Overview

Adds to winning or losing positions based on the configured mode (Pyramiding, Averaging Down, Scale Out).

### Key Features

* **3 Modes**:

  1. **Pyramiding** — add to winners (trend trading)
  2. **Averaging Down** — add to losers (high risk!)
  3. **Scale Out** — gradually exit positions
* Configurable trigger distances
* Maximum scale limits
* Automatic lot size calculation

### Configuration Example

```go
config := orchestrators.PositionScaleConfig{
    Mode:            Pyramiding, // Add to winners
    Symbol:          "EURUSD",
    TriggerDistance: 200,        // Add every 200 points profit
    MaxScales:       3,          // Max 3 additions
    ScaleLotSize:    0.01,
}

scale := orchestrators.NewPositionScale(suga, config)
scale.Start()
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
* **Action**: Opens additional positions
* **Output**: Scaled exposure
* **Cycle Time**: 5 seconds (configurable)

[Full Documentation](12_Position_scaler.md)

---

## Orchestrator #13: GridTrade

### Overview

📊 Places a grid of buy/sell orders at regular intervals, profiting from price oscillations in range-bound markets.

### Key Features

* Configurable grid size and spacing
* Automatic order placement
* Works best in sideways markets
* Optional grid rebuilding after fills

### Configuration Example

```go
config := orchestrators.GridTradeConfig{
    Symbol:        "EURUSD",
    GridSize:      5,   // 5 grid levels
    GridStep:      100, // 100 points spacing
    LotSize:       0.01,
    MaxPositions:  10,
    CheckInterval: 5 * time.Second,
}

gridTrade := orchestrators.NewGridTrade(suga, config)
gridTrade.Start()
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

## Orchestrator #14: RiskManage

### Overview

Monitors account-level risk metrics 24/7 and automatically intervenes when danger limits are exceeded.

### Key Features

* **5 Protection Types**:

  1. Drawdown protection (percentage and absolute)
  2. Daily loss limits
  3. Daily profit targets
  4. Margin level monitoring
  5. Position count limits
* Emergency close of all positions
* Trade blocking when limits are hit
* Daily automatic reset


### Configuration Example

```go
config := orchestrators.RiskManageConfig{
    MaxDrawdownPercent:  10.0,  // 10% max drawdown
    MaxDrawdownAbsolute: 1000.0, // $1000 max loss
    DailyLossLimit:      500.0,  // $500 daily limit
    DailyProfitTarget:   1000.0, // $1000 daily goal
    MinMarginLevel:      150.0,  // 150% minimum margin
    EnableAutoClose:     true,   // Auto-close on breach
    CheckInterval:       5 * time.Second,
}

riskManage := orchestrators.NewRiskManage(suga, config)
riskManage.Start()
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

## Orchestrator #15: PortfolioRebalance

### Overview

Maintains balanced exposure across multiple trading symbols by automatically rebalancing positions to match target allocations.

### Key Features

* Equal-weight or custom allocation percentages
* Automatic rebalancing when deviation exceeds threshold
* Multi-symbol portfolio management
* Diversification enforcement

### Configuration Example

```go
config := orchestrators.PortfolioRebalanceConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.25, // 25% each
        "GBPUSD": 0.25,
        "USDJPY": 0.25,
        "XAUUSD": 0.25,
    },
    TotalExposure:      10000.0,          // $10k portfolio
    RebalanceThreshold: 10.0,             // Rebalance at 10% off
    CheckInterval:      30 * time.Minute,
}

rebalance := orchestrators.NewPortfolioRebalance(suga, config)
rebalance.Start()
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
* **Cycle Time**: 30 minutes (configurable)

[Full Documentation](15_Portfolio_rebalancer.md)

---

## Preset #16: Adaptive Orchestrator Preset

### Overview

Intelligent multi-strategy system that analyzes market conditions every 5 minutes and automatically selects the most appropriate orchestrator to run.

### Key Features

* **5 Market Modes**:

  1. **Grid Mode** — low volatility → GridTrade
  2. **Managed Mode** — medium volatility → TrailingStop + PositionScale
  3. **Protection Mode** — high volatility → RiskManage
  4. **Portfolio Mode** — multi-symbol → PortfolioRebalance
  5. **Trending Mode** — strong trend → PositionScale (pyramiding)
* Automatic mode switching every cycle
* Built-in safety limits
* Continuous operation with daily limits


### Configuration Example

```go
preset := presets.NewAdaptiveOrchestratorPreset(suga)

// Customize
preset.Symbol = "EURUSD"
preset.LowVolatilityThreshold = 25.0  // Grid below 25 pts
preset.HighVolatilityThreshold = 50.0 // Protection above 50 pts
preset.CycleDuration = 5 * time.Minute
preset.MaxDailyLoss = 500.0
preset.MaxDailyProfit = 1000.0

// Run (blocks until stopped with Ctrl+C)
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
| **Strong Uptrend**       | TrailingStop + PositionScale (Pyramiding) | Protect profits, add to winners |
| **Strong Downtrend**     | TrailingStop (short positions)            | Protect short profits           |
| **Sideways / Ranging**   | GridTrade                                 | Profit from oscillations        |
| **Low Volatility**       | GridTrade                                 | Small frequent profits          |
| **High Volatility**      | RiskManage                                | Protection mode                 |
| **Multi-Symbol Active**  | PortfolioRebalance                        | Balance allocation              |
| **Uncertain / Changing** | Adaptive Preset                           | Auto-adapt to conditions        |

---

### By Trading Style

| Trading Style        | Recommended Orchestrators               |
| -------------------- | --------------------------------------- |
| **Scalping**         | TrailingStop (tight parameters)         |
| **Day Trading**      | TrailingStop + RiskManage               |
| **Swing Trading**    | PositionScale (pyramiding) + RiskManage |
| **Position Trading** | PortfolioRebalance + RiskManage         |
| **Grid Trading**     | GridTrade + RiskManage                  |
| **Automated 24/7**   | Adaptive Preset                         |

---

### By Risk Tolerance

| Risk Level          | Orchestrator Combination                    |
| ------------------- | ------------------------------------------- |
| **Conservative**    | RiskManage only (5% drawdown limit)         |
| **Moderate**        | TrailingStop + RiskManage                   |
| **Balanced**        | PortfolioRebalance + RiskManage             |
| **Aggressive**      | PositionScale + TrailingStop + RiskManage   |
| **Very Aggressive** | GridTrade + PositionScale (NOT RECOMMENDED) |

---

## Common Orchestrator Combinations

### Combination 1: Safe Trend Trading

**Orchestrators**: TrailingStop + RiskManage

**Purpose**: Follow trends with profit protection and account safety

```go
// Trailing stop for profit protection
tsConfig := orchestrators.DefaultTrailingStopConfig()
trailingManage := orchestrators.NewTrailingStopManage(suga, tsConfig)
trailingManage.Start()

// Risk management for account protection
riskConfig := orchestrators.DefaultRiskManageConfig()
riskManage := orchestrators.NewRiskManage(suga, riskConfig)
riskManage.Start()
```


**Best For**: Beginners, conservative traders, trending markets

---

### Combination 2: Aggressive Scaling

**Orchestrators**: PositionScale + TrailingStop + RiskManage

**Purpose**: Scale into trends, protect profits, limit risk

```go
// Scale for pyramiding
scaleConfig := orchestrators.DefaultPositionScaleConfig("EURUSD")
scaleConfig.Mode = orchestrators.Pyramiding
scale := orchestrators.NewPositionScale(suga, scaleConfig)
scale.Start()

// Trailing stop for profit protection
tsManage := orchestrators.NewTrailingStopManage(suga, tsConfig)
tsManage.Start()

// Risk management for safety
riskManage := orchestrators.NewRiskManage(suga, riskConfig)
riskManage.Start()
```

**Best For**: Experienced traders, strong trends, large accounts

---

### Combination 3: Portfolio Diversification

**Orchestrators**: PortfolioRebalance + RiskManage

**Purpose**: Maintain balanced multi-symbol portfolio with protection

```go
// Portfolio rebalance
portfolioConfig := orchestrators.PortfolioRebalanceConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.25,
        "GBPUSD": 0.25,
        "USDJPY": 0.25,
        "XAUUSD": 0.25,
    },
    TotalExposure: 10000.0,
}

rebalance := orchestrators.NewPortfolioRebalance(suga, portfolioConfig)
rebalance.Start()

// Risk management for overall protection
riskManage := orchestrators.NewRiskManage(suga, riskConfig)
riskManage.Start()
```

**Best For**: Portfolio managers, diversification seekers, medium–large accounts

---

### Combination 4: Range Trading

**Orchestrators**: GridTrade + RiskManage

**Purpose**: Profit from sideways markets with risk control

```go
// Grid trade
gridConfig := orchestrators.GridTradeConfig{
    Symbol:   "EURUSD",
    GridSize: 5,
    GridStep: 100,
}

gridTrade := orchestrators.NewGridTrade(suga, gridConfig)
gridTrade.Start()

// Risk management (important for grid!)
riskConfig := orchestrators.RiskManageConfig{
    MaxDrawdownPercent: 15.0, // Wider for grid
    EnableAutoClose:    true,
}

riskManage := orchestrators.NewRiskManage(suga, riskConfig)
riskManage.Start()
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

go run main.go 12  # PositionScale
go run main.go scale

go run main.go 13  # GridTrade
go run main.go grid

go run main.go 14  # RiskManage
go run main.go risk

go run main.go 15  # PortfolioRebalance
go run main.go rebalance

# Adaptive preset
go run main.go 16       # Adaptive Preset
go run main.go adaptive
go run main.go preset
```

---

### Programmatic Usage

```go
// Initialize connection
suga := mt5.NewMT5Suga(account)

// Create orchestrator with config
config := orchestrators.DefaultTrailingStopConfig()
orchestrator := orchestrators.NewTrailingStopManage(suga, config)

// Start (non-blocking)
if err := orchestrator.Start(); err != nil {
    log.Fatal(err)
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

**Always Use RiskManage**

* Combine **any** orchestrator with RiskManage
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

* TrailingStop + RiskManage = safe
* GridTrade + PositionScale = dangerous!
* Portfolio + RiskManage = smart diversification

---

### ⚠️ DON'T

**Never Trade Without Risk Management**

* Even on demo, practice with RiskManage
* Learn proper risk control habits
* Safety should be automatic

**Don't Run Conflicting Orchestrators**

* GridTrade + Trending Mode scaling = conflict
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

* Invalid configuration (check error message)
* Symbol not tradable
* Insufficient margin
* Already running

**Solution**:

```go
if err := orchestrator.Start(); err != nil {
    fmt.Printf("Start failed: %v\n", err)
    // Check error details
}
```

```

---

### Problem: Orchestrator Stops Unexpectedly

**Possible Causes**:

* Safety limit reached (check logs)
* Context cancelled
* Fatal error encountered

**Solution**:

```go
status := orchestrator.GetStatus()
fmt.Printf("Status: %+v\n", status)
// Review error count and last operation
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
// Get status
status := orchestrator.GetStatus()

// Display metrics
fmt.Printf("Name: %s\n", status.Name)
fmt.Printf("Running: %v\n", status.IsRunning)
fmt.Printf("SuccessCount: %d\n", status.SuccessCount)
fmt.Printf("ErrorCount: %d\n", status.ErrorCount)
fmt.Printf("TotalOperations: %d\n", status.SuccessCount+status.ErrorCount)
fmt.Printf("LastOperation: %s\n", status.LastOperation)
fmt.Printf("Uptime: %v\n", time.Since(status.StartTime))
```

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

1. **Extend BaseOrchestrator**:

```go
type MyOrchestrator struct {
    *orchestrators.BaseOrchestrator
    suga   *mt5.MT5Suga
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

Enhance orchestrators with technical analysis:

```go
// Get bars for analysis
bars, _ := suga.GetBars("EURUSD", pb.Timeframe_M15, 100)

// Calculate indicators
atr := calculateATR(bars, 14)
ma20 := calculateMA(bars, 20)

// Adjust orchestrator parameters based on analysis
if atr > 50 {
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
    scale.Start()
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
| 11 | TrailingStopManage | Profit protection  | Low        |
| 12 | PositionScale      | Position sizing    | Medium     |
| 13 | GridTrade          | Range trading      | Medium     |
| 14 | RiskManage         | Account protection | Low        |
| 15 | PortfolioRebalance | Portfolio balance  | Medium     |
| 16 | Adaptive Preset    | Multi-strategy     | High       |

---

### Key Takeaways

**Always Use RiskManage**

* Combine with any orchestrator
* Essential for automated trading
* Prevents catastrophic losses

**Start Simple**

* Begin with TrailingStop + RiskManage
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
3. **Start with RiskManage** — Always combine with risk management
4. **Monitor Performance** — Track metrics and adjust parameters
5. **Graduate to Combinations** — Once comfortable, try multi-orchestrator setups
6. **Consider Adaptive Preset** — For hands-off adaptive trading

---

## Complete Documentation Links

* [TrailingStopManage — Full Docs](11_Trailing_stop.md)
* [PositionScale — Full Docs](12_Position_scaler.md)
* [GridTrade — Full Docs](13_Grid_trader.md)
* [RiskManage — Full Docs](14_Risk_manager.md)
* [PortfolioRebalance — Full Docs](15_Portfolio_rebalancer.md)
* [Adaptive Orchestrator Preset — Full Docs](16_AdaptiveOrchestratorPreset.md)


---

**Remember**: Orchestrators are powerful tools, but they require understanding and proper configuration. Always test on demo first, use risk management, and monitor performance actively until you are confident in their behavior.
