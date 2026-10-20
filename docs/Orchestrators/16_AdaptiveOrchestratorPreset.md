
# 🧩 Adaptive Orchestrator Preset — Intelligent Multi-Strategy System

## Description

**Adaptive Orchestrator Preset** is an intelligent trading system that automatically selects and executes the most appropriate orchestrator based on real-time market condition analysis.

It combines **all five orchestrators** into a single adaptive system that dynamically switches strategies every cycle.

### Strategy Principle

Continuously analyzes market conditions and adapts the trading approach:

* **Market Analysis** — calculates volatility, detects trends, and counts active positions
* **Strategy Selection** — chooses the best orchestrator for current conditions
* **Dynamic Execution** — runs the selected orchestrator in 5-minute cycles (plus a 30-second pause between cycles)
* **Continuous Adaptation** — repeats analysis and switches strategies as markets change
* **Safety Controls** — total-loss limit and daily loss/profit limits checked after every cycle

**File:** `examples/demos/presets/16_AdaptiveOrchestratorPreset.go` (package `presets`, import `"github.com/MetaRPC/GoMT5/examples/demos/presets"` — part of the repository's `examples` module, available by cloning the repository)

> **IMPORTANT DISCLAIMER — EDUCATIONAL EXAMPLE ONLY**
>
> This preset demonstrates advanced orchestrator composition and adaptive strategy selection.
> **NOT production-ready!**
> Market analysis is simplified for demonstration purposes. Always test extensively on demo accounts and customize for your specific market knowledge and risk tolerance.

---

## What Is an Adaptive Preset?

### Simple Explanation

Instead of running one fixed strategy, the **Adaptive Orchestrator Preset** acts like an intelligent trading manager that:

1. Looks at the market at the start of every cycle (5-minute cycle + 30-second pause)
2. Decides which strategy fits best (grid, trailing, risk protection, etc.)
3. Runs that strategy for one cycle
4. Repeats with continuous market analysis

**Analogy:** Like a professional trader who switches between different strategies:

* **Quiet market** (low volatility) → Grid trading
* **Normal market** → Trailing stops with position scaling
* **Volatile market** → Risk protection mode
* **Trending market** → Pyramiding into winning positions
* **Multiple symbols** → Portfolio balancing


### Dependencies

* **MT5Sugar** — high-level convenience API (`GetBalance`, `GetPriceInfo`, `GetOpenPositions`)
* **All 5 Orchestrators** — GridTrader, TrailingStopManager, PositionScaler, RiskManager, PortfolioRebalancer (package `orchestrators`)
* **Helpers** — `helpers.WaitWithProgressBarAndCallback` progress bar (package `examples/demos/helpers`)

---

## 💡 Market Modes and Strategy Selection

### Mode 1: Grid Mode (Low Volatility)

**Condition:** Volatility < `LowVolatilityThreshold` (25 points)

**Why:** Range-bound markets are perfect for grid strategies.

**Orchestrator Used:** GridTrader

**Configuration** (from `executeGridMode`):

```go
config := orchestrators.GridTraderConfig{
    Symbol:        p.Symbol, // "EURUSD" by default
    GridSize:      5,        // 5 grid levels each side
    GridStep:      100,      // 100 points spacing
    LotSize:       0.01,     // Micro-lots
    MaxPositions:  10,
    TakeProfit:    0,
    StopLoss:      0,
    CheckInterval: 5 * time.Second,
    RebuildOnFill: false,
}
```

**Expected Behavior:**

* Places 5 BUY LIMIT and 5 SELL LIMIT pending orders around the mid price
* With `TakeProfit: 0` each order's TP is one grid step away; `StopLoss: 0` = no SL (see [GridTrader](13_Grid_trader.md))
* Works best in sideways markets
* When the cycle ends, `Stop()` deletes the grid's pending orders; filled positions stay open

---

### Mode 2: Managed Mode (Medium Volatility)

**Condition:** 25–50 points volatility and no strong trend (|TrendStrength| ≤ 0.7)

**Why:** Normal conditions suit position management with scaling.

**Orchestrators Used:**

1. TrailingStopManager (profit protection)
2. PositionScaler (pyramiding)

**Configuration** (from `executeManagedMode`):

```go
// Trailing Stop: defaults (TrailingDistance 200, ActivationProfit 300, UpdateInterval 2s, all symbols)
tsConfig := orchestrators.DefaultTrailingStopConfig()
tsManager := orchestrators.NewTrailingStopManager(p.sugar, tsConfig)

// Position Scaler: defaults for the symbol (TriggerDistance 200, MinProfitToScale 100,
// InitialLotSize 0.10, ScaleLotSize 0.05, MaxScales 3, TotalMaxLotSize 1.0, StopLossPerScale 150, CheckInterval 5s)
scalerConfig := orchestrators.DefaultPositionScalerConfig(p.Symbol)
scalerConfig.Mode = orchestrators.Pyramiding
scaler := orchestrators.NewPositionScaler(p.sugar, scalerConfig)
```

**Expected Behavior:**

* Protects profits with trailing stops
* Opens a 0.10-lot BUY base position on the symbol if it has none (`InitialLotSize`, once per scaler run) and scales into it once it is in profit
* Balanced approach for normal markets

---

### Mode 3: Protection Mode (High Volatility)

**Condition:** Volatility > `HighVolatilityThreshold` (50 points)

**Why:** Volatile markets require strict risk control.

**Orchestrator Used:** RiskManager

**Configuration** (from `executeProtectionMode`):

```go
config := orchestrators.DefaultRiskManagerConfig() // CheckInterval 5s, $1000 abs. drawdown, $500 daily loss, ...
config.MaxDrawdownPercent = 5.0                    // 5% max drawdown
config.EnableAutoClose = true                      // Auto-close on breach
```

**Expected Behavior:**

* Monitors risk limits every 5 seconds during the cycle
* Closes ALL positions if drawdown reaches 5% (or $1000) from the peak, or the daily loss reaches $500; closes the most losing position while the margin level is below 150%
* The peak/daily baseline is the equity when this cycle's RiskManager starts
* Defensive mode to prevent blow-ups

---

### Mode 4: Portfolio Mode (Multi-Symbol)

**Condition:** `EnablePortfolioMode` is true (default `false`) **and** two or more different symbols have open positions

**Why:** Multiple symbols require balanced exposure.

**Orchestrator Used:** PortfolioRebalancer

**Configuration** (from `executePortfolioMode`):

```go
// PortfolioSymbols defaults to EURUSD, GBPUSD, USDJPY, XAUUSD -> 0.25 each,
// TotalExposure 10000, RebalanceThreshold 10, MaxDeviation 0.20
config := orchestrators.DefaultPortfolioRebalancerConfig(p.PortfolioSymbols)
// The default 1-hour check interval would never come round in a cycle: the rebalancer checks
// when it starts and then every fifth of the cycle.
config.CheckInterval = p.CycleDuration / 5
rebalancer := orchestrators.NewPortfolioRebalancer(p.sugar, config)
```


**Expected Behavior**:

* Checks when it starts and then every minute (`CycleDuration / 5`), and rebalances the open positions when a symbol deviates more than 10% from its target
* Resizes existing positions only: symbols of `PortfolioSymbols` without a position are left alone (see [PortfolioRebalancer](15_Portfolio_rebalancer.md))

---

### Mode 5: TRENDING MODE (Strong Directional Movement)

**Condition**: Volatility between the thresholds and |TrendStrength| > 0.7 (the "trend" is a placeholder derived from the spread, see Critical Warnings)

**Why**: Strong trends are perfect for pyramiding strategies

**Orchestrator Used**: PositionScaler (pyramiding mode)

**Configuration** (from `executeTrendingMode`):

```go
config := orchestrators.DefaultPositionScalerConfig(p.Symbol)
config.Mode = orchestrators.Pyramiding
config.TriggerDistance = 200 // 200 points trigger
config.MaxScales = 3         // Max 3 scales
```

**Expected Behavior**:

* Adds to winning positions
* Rides strong trends
* Maximizes profit in directional markets

---

## Configuration Parameters

| Parameter                    | Type          | Default    | Description                    |
| ---------------------------- | ------------- | ---------- | ------------------------------ |
| `Symbol`                     | string        | `"EURUSD"` | Primary trading symbol         |
| `BaseRiskAmount`             | float64       | `20.0`     | Only used for the total-loss stop: stop when total P/L < −5 × BaseRiskAmount |
| `LowVolatilityThreshold`     | float64       | `25.0`     | Max points for grid mode       |
| `HighVolatilityThreshold`    | float64       | `50.0`     | Min points for protection mode |
| `CycleDuration`              | time.Duration | `5m`       | Duration of one cycle (0 = 5 m) |
| `EnablePortfolioMode`        | bool          | `false`    | Enable multi-symbol management |
| `PortfolioSymbols`           | []string      | 4 symbols  | Symbols for portfolio mode     |
| `MaxConcurrentOrchestrators` | int           | `2`        | Most orchestrators started per cycle; a further one is not started (0 = no limit) |
| `MaxDailyLoss`               | float64       | `500.0`    | Stop when balance − day start balance < −MaxDailyLoss |
| `MaxDailyProfit`             | float64       | `1000.0`   | Stop when balance − day start balance > MaxDailyProfit |
| `MaxCycles`                  | int           | `0`        | Stop after this many cycles (0 = until Ctrl+C / SIGTERM or a safety stop; the demo runner sets 2) |

"Daily" limits are measured from the balance at `Execute()` start, and from the balance at the first safety check of each new day (local date) after that. The fields are exported, so set them on the preset before calling `Execute()`.

### Configuration Example

Imports: `presets "github.com/MetaRPC/GoMT5/examples/demos/presets"` and `mt5 "github.com/MetaRPC/GoMT5/examples/mt5"` (both in the repository's `examples` module).

```go
sugar, err := mt5.NewMT5Sugar(12345678, "your_password", "mt5.mrpc.pro:443")
if err != nil {
    return err
}
defer sugar.Disconnect() // also when QuickConnect fails
if err := sugar.QuickConnect("YourBroker-MT5 Demo"); err != nil {
    return err
}

preset := presets.NewAdaptiveOrchestratorPreset(sugar)

// Customize parameters
preset.Symbol = "EURUSD"
preset.BaseRiskAmount = 20.0
preset.LowVolatilityThreshold = 25.0
preset.HighVolatilityThreshold = 50.0
preset.CycleDuration = 5 * time.Minute
preset.EnablePortfolioMode = false
preset.MaxDailyLoss = 500.0
preset.MaxDailyProfit = 1000.0
preset.MaxCycles = 0 // 0 = until Ctrl+C; e.g. 2 = stop after two cycles

// Execute: blocks until Ctrl+C / SIGTERM, a safety stop or MaxCycles cycles, returns the total P/L
totalProfit, err := preset.Execute()
```

---

## How to Run

You can execute this preset using several command variations:

```bash
# Navigate to examples directory
cd examples/demos

# Option 1: By number
go run main.go 16

# Option 2: By keyword
go run main.go adaptive

# Option 3: By category
go run main.go preset
```

All commands launch `main.go -> RunOrchestrator_AdaptivePreset()`, which connects with `mt5.NewMT5Sugar(...)` + `QuickConnect(...)` and runs the preset with its default configuration and `MaxCycles = 2` (`presets.NewAdaptiveOrchestratorPreset(sugar)`, `preset.MaxCycles = 2`, `Execute()`).
The preset runs two cycles (about 11 minutes) - so it also ends when nobody can press `Ctrl+C`, e.g. with stdin closed - unless you press `Ctrl+C` or a safety limit stops it earlier; either way `Execute()` stops the running orchestrators, shows the final results and returns, the runner prints `📊 Final Result: $<total P/L>` and disconnects the terminal (the runner disconnects however the run ends, also on an error or Ctrl+C).

---

## 🟢 Algorithm

### Complete Flow

```
START
|
Initialize Preset
- Get starting balance
- Set thresholds
- Setup signal handling
|
INFINITE CYCLE LOOP:
|
+-- CYCLE #N START
|
+-- STEP 1: ANALYZE MARKET
|   - Get price info for primary symbol
|   - Calculate volatility (spread × 10)
|   - Count open positions and symbols
|   - Estimate trend strength
|
| SELECT MODE (priority order):
|   1. Portfolio mode? (EnablePortfolioMode and 2+ symbols)
|   2. High volatility? (>50 pts) → Protection Mode
|   3. Low volatility? (<25 pts) → Grid Mode
|   4. Strong trend? (>70%) → Trending Mode
|   5. Otherwise → Managed Mode
|
+-- STEP 2: STOP PREVIOUS ORCHESTRATORS
|   - Stop all orchestrators from last cycle
|   - Clear active orchestrators list
|
+-- STEP 3: EXECUTE SELECTED MODE
|   - Grid Mode → Start GridTrader
|   - Managed Mode → Start TrailingStop + PositionScaler
|   - Protection Mode → Start RiskManager
|   - Portfolio Mode → Start PortfolioRebalancer
|   - Trending Mode → Start PositionScaler (pyramiding)
|
+-- STEP 4: MONITOR CYCLE
|   - Wait for CycleDuration (5 minutes)
|   - Show progress bar
|   - Display orchestrator status every 10 seconds
|
+-- STEP 5: CALCULATE RESULTS
|   - Get end balance
|   - Calculate cycle profit/loss
|   - Update total profit
|
+-- STEP 6: CHECK SAFETY LIMITS
|   - Total loss > $100? → STOP
|   - Daily loss > $500? → STOP
|   - Daily profit > $1000? → STOP (success)
|   - MaxCycles cycles done (if > 0)? → STOP
|
+-- STEP 7: PAUSE
|   - Wait 30 seconds before next cycle
|
+-- REPEAT CYCLE
|
ON SAFETY LIMIT, MaxCycles or Ctrl+C / SIGTERM:
- The cycle loop ends (Ctrl+C also ends the running cycle early)
- Stop all active orchestrators
- Show final results
- Return total profit
|
END
```

#### 1. Entry Point: Execute

> Lines `269–323`


```go
func (p *AdaptiveOrchestratorPreset) Execute() (float64, error) {
	printHeader()
	if p.CycleDuration <= 0 {
		p.CycleDuration = 5 * time.Minute
	}
	fmt.Printf("   Getting starting balance...\n")

	// Initialize
	var err error
	p.initialBalance, err = p.loadedBalance()
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}
	p.dailyStartBalance = p.initialBalance
	p.dailyStartDate = time.Now()

	fmt.Printf("  💰 Starting balance: $%.2f\n", p.initialBalance)
	fmt.Printf("  📊 Primary symbol: %s\n", p.Symbol)
	fmt.Printf("  🎯 Base risk: $%.2f\n", p.BaseRiskAmount)
	fmt.Printf("  📈 Volatility thresholds: Low < %.0f pts < Medium < %.0f pts < High\n\n",
		p.LowVolatilityThreshold, p.HighVolatilityThreshold)

	// Setup context and signal handling
	p.ctx, p.cancel = context.WithCancel(context.Background())
	defer p.cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// Main cycle loop
	cyclesDone := make(chan struct{})
	go func() {
		defer close(cyclesDone)
		p.runCycles()
	}()

	// Wait for interrupt, or for the cycles to end by themselves (safety stop or MaxCycles)
	select {
	case <-sigChan:
		fmt.Println("\n\n🛑 Stopping preset...")
	case <-cyclesDone:
		fmt.Println("\n\n🛑 Cycles ended - stopping preset...")
	}
	p.cancel()
	<-cyclesDone // the current cycle returns before its orchestrators are stopped

	// Stop all active orchestrators
	p.stopAllOrchestrators()

	// Show final results
	p.showFinalResults()

	return p.totalProfit, nil
}
```

**What happens**:

* Gets starting balance for profit calculation (fails if the balance cannot be read). All balance reads go through `loadedBalance()`, which turns a balance of 0 into the error `account data not loaded yet (balance 0)`: a terminal that the server has just started again answers 0 for a few seconds, and taken as real that would be a cycle "losing" the whole balance and a total-loss stop
* Creates context for graceful shutdown
* Sets up Ctrl+C / SIGTERM handling
* Launches main cycle loop in a goroutine
* Blocks until the user stops execution (Ctrl+C / SIGTERM) **or** the cycle loop ends by itself after a safety limit or `MaxCycles` cycles
* Cancels the context and waits for the cycle loop to return, so no cycle starts orchestrators while they are being stopped
* Stops the orchestrators of the last cycle, shows the final results and returns `(totalProfit, nil)` - the caller can then disconnect

---

#### 2. Main Cycle Loop
  
> Lines `326–360`  


```go
func (p *AdaptiveOrchestratorPreset) runCycles() {
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			p.cycleNumber++
			p.executeCycle()

			// Interrupted during the cycle
			if p.ctx.Err() != nil {
				return
			}

			// Check safety limits (Execute stops everything when the cycles end)
			if p.checkSafetyLimits() {
				return
			}

			// Planned number of cycles done
			if p.MaxCycles > 0 && p.cycleNumber >= p.MaxCycles {
				fmt.Printf("
  ✓ %d cycle(s) done (MaxCycles)
", p.cycleNumber)
				return
			}

			// Wait before next cycle
			fmt.Printf("\n⏸️  Pausing 30 seconds before next cycle...\n")
			select {
			case <-p.ctx.Done():
				return
			case <-time.After(30 * time.Second):
			}
		}
	}
}
```

**What happens**:

* Loop running cycles (endless with `MaxCycles = 0`)
* Increments cycle counter
* Executes one trading cycle
* Returns at once when the context was cancelled during the cycle (Ctrl+C)
* Checks safety limits after each cycle
* Ends the loop if limits are breached or `MaxCycles` cycles are done; `Execute()` sees the loop end and stops everything
* Pauses 30 seconds between cycles (interruptible by the context)

---

#### 3. Execute Single Cycle
  
> Lines `363–420`  


```go
func (p *AdaptiveOrchestratorPreset) executeCycle() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  CYCLE #%-3d                                               ║\n", p.cycleNumber)
	fmt.Printf("╚════════════════════════════════════════════════════════════╝\n\n")

	// Step 1: Analyze market conditions
	condition, err := p.analyzeMarketConditions()
	if err != nil {
		fmt.Printf("  ✗ Market analysis failed: %v\n", err)
		return
	}

	fmt.Printf("  🔍 Market Analysis Complete:\n")
	fmt.Printf("     Mode: %s\n", condition.Mode)
	fmt.Printf("     Volatility: %.1f points\n", condition.VolatilityPoints)
	fmt.Printf("     Trend: %.2f (%.0f%% strength)\n", condition.TrendStrength, math.Abs(condition.TrendStrength)*100)
	fmt.Printf("     Reason: %s\n", condition.Reason)
	fmt.Printf("     Active Positions: %d\n", condition.ActivePositions)
	fmt.Printf("     Active Symbols: %d\n\n", condition.ActiveSymbols)

	// Step 2: Stop any running orchestrators from previous cycle
	p.stopAllOrchestrators()

	// Step 3: Select and execute appropriate orchestrator(s)
	cycleStartBalance, startErr := p.loadedBalance()

	switch condition.Mode {
	case GridMode:
		p.executeGridMode(condition)
	case ManagedMode:
		p.executeManagedMode(condition)
	case ProtectionMode:
		p.executeProtectionMode(condition)
	case PortfolioMode:
		p.executePortfolioMode(condition)
	case TrendingMode:
		p.executeTrendingMode(condition)
	default:
		fmt.Printf("  ⚠️  Unknown market mode, skipping cycle\n")
	}

	// Monitor for cycle duration
	p.monitorCycle()

	// Calculate cycle profit (not counted when a balance could not be read)
	cycleEndBalance, endErr := p.loadedBalance()
	if startErr != nil || endErr != nil {
		fmt.Printf("\n  ⚠️  Cycle #%d result unknown: balance not available (%v %v)\n", p.cycleNumber, startErr, endErr)
		return
	}
	cycleProfit := cycleEndBalance - cycleStartBalance
	p.totalProfit += cycleProfit

	fmt.Printf("\n  📊 Cycle #%d Result:\n", p.cycleNumber)
	fmt.Printf("     Profit: $%.2f\n", cycleProfit)
	fmt.Printf("     Total P/L: $%.2f\n", p.totalProfit)
	fmt.Printf("     Current Balance: $%.2f\n", cycleEndBalance)
}
```

**What happens**:

* Analyzes current market conditions
* Stops any orchestrators from the previous cycle
* Records starting balance
* Selects and executes the appropriate mode
* Monitors execution for the cycle duration (5 minutes)
* Calculates profit for this cycle (balance difference = realized P/L only); if a balance cannot be read (or reads 0) the cycle result is skipped instead of being counted as a huge loss
* Updates total profit
* If the market analysis fails, the cycle ends at once (no orchestrators are stopped or started) and the loop moves on to the safety check and pause

---

#### 4. Market Analysis
  
> Lines `427–480`  


```go
func (p *AdaptiveOrchestratorPreset) analyzeMarketConditions() (*MarketCondition, error) {
	// Get current market data
	priceInfo, err := p.sugar.GetPriceInfo(p.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get price info: %w", err)
	}

	// Calculate volatility from spread
	spreadPoints := priceInfo.SpreadPips
	estimatedVolatility := spreadPoints * 10 // Rough estimate

	// Get active positions
	positions, err := p.sugar.GetOpenPositions()
	if err != nil {
		positions = make([]*pb.PositionInfo, 0)
	}

	// Count unique symbols
	symbolMap := make(map[string]bool)
	for _, pos := range positions {
		symbolMap[pos.Symbol] = true
	}

	condition := &MarketCondition{
		VolatilityPoints: estimatedVolatility,
		ActivePositions:  len(positions),
		ActiveSymbols:    len(symbolMap),
	}

	// Simplified trend detection (would use bars in production)
	// For demo, use random-like value based on spread
	condition.TrendStrength = math.Mod(spreadPoints, 2.0) - 1.0 // -1.0 to 1.0

	// Determine market mode (PRIORITY ORDER: Portfolio → Protection → Grid → Trending → Managed)
	if p.EnablePortfolioMode && len(symbolMap) >= 2 {
		condition.Mode = PortfolioMode
		condition.Reason = fmt.Sprintf("Multi-symbol portfolio (%d symbols active)", len(symbolMap))
	} else if estimatedVolatility > p.HighVolatilityThreshold {
		condition.Mode = ProtectionMode
		condition.Reason = fmt.Sprintf("High volatility (%.1f pts) - risk protection needed", estimatedVolatility)
	} else if estimatedVolatility < p.LowVolatilityThreshold {
		// LOW VOLATILITY = GRID MODE (checked BEFORE trending to prioritize range trading)
		condition.Mode = GridMode
		condition.Reason = fmt.Sprintf("Low volatility (%.1f pts) - range-bound market", estimatedVolatility)
	} else if math.Abs(condition.TrendStrength) > 0.7 {
		condition.Mode = TrendingMode
		condition.Reason = fmt.Sprintf("Strong trend detected (strength: %.0f%%)", math.Abs(condition.TrendStrength)*100)
	} else {
		condition.Mode = ManagedMode
		condition.Reason = fmt.Sprintf("Medium volatility (%.1f pts) - normal conditions", estimatedVolatility)
	}

	return condition, nil
}
```

**What happens**:

* Fetches current price info using `GetPriceInfo()`
* Calculates volatility as `spread × 10`, where the spread is `PriceInfo.SpreadPips` = (Ask − Bid) / point, i.e. in **points** (simplified estimate)
* Gets all open positions (an error is treated as "no positions")
* Counts unique symbols for portfolio detection
* Derives a placeholder `TrendStrength` = `math.Mod(spread, 2) − 1`
* Applies mode selection logic with priority:

  1. **Portfolio** if `EnablePortfolioMode` and 2+ symbols active

  2. **Protection** if volatility > 50 pts

  3. **Grid** if volatility < 25 pts

  4. **Trending** if trend strength > 70%

  5. **Managed** otherwise (default)

---

#### 5. Mode Execution Example: Grid Mode
  
> Lines `487–504`, `563–576`  


```go
func (p *AdaptiveOrchestratorPreset) executeGridMode(condition *MarketCondition) {
	fmt.Printf("  🎯 Executing: GRID TRADER ORCHESTRATOR\n\n")

	config := orchestrators.GridTraderConfig{
		Symbol:        p.Symbol,
		GridSize:      5,
		GridStep:      100,
		LotSize:       0.01,
		MaxPositions:  10,
		TakeProfit:    0,
		StopLoss:      0,
		CheckInterval: 5 * time.Second,
		RebuildOnFill: false,
	}

	gridTrader := orchestrators.NewGridTrader(p.sugar, config)
	p.startOrchestrator(gridTrader, "Grid Trader")
}

// startOrchestrator starts an orchestrator unless MaxConcurrentOrchestrators are already
// running (0 = no limit) and reports whether it runs.
func (p *AdaptiveOrchestratorPreset) startOrchestrator(orch orchestrators.Orchestrator, name string) bool {
	if p.MaxConcurrentOrchestrators > 0 && len(p.activeOrchestrators) >= p.MaxConcurrentOrchestrators {
		fmt.Printf("  ⏭️  %s not started: %d orchestrator(s) running (MaxConcurrentOrchestrators = %d)\n",
			name, len(p.activeOrchestrators), p.MaxConcurrentOrchestrators)
		return false
	}
	if err := orch.Start(); err != nil {
		fmt.Printf("  ✗ Failed to start %s: %v\n", name, err)
		return false
	}
	p.activeOrchestrators = append(p.activeOrchestrators, orch)
	fmt.Printf("  ✓ %s started\n", name)
	return true
}
```

**What happens**:

* Creates configuration for GridTrader
* Instantiates a new GridTrader orchestrator
* `startOrchestrator()` enforces `MaxConcurrentOrchestrators`, starts the orchestrator (its `Start()` validates the configuration) and adds it to the active orchestrators list
* Same pattern is used for the other modes (ManagedMode starts the scaler only if the trailing stop manager started)

---

#### 6. Cycle Monitoring
 
> Lines `583–601`  


```go
func (p *AdaptiveOrchestratorPreset) monitorCycle() {
	totalSeconds := int(p.CycleDuration.Seconds())

	// Use progress bar with callback to show orchestrator status
	helpers.WaitWithProgressBarAndCallback(
		totalSeconds,
		fmt.Sprintf("Cycle #%d monitoring", p.cycleNumber),
		10*time.Second, // Update orchestrator status every 10 seconds
		func() bool {
			// Show detailed orchestrator status during callback
			fmt.Println() // New line before status
			p.showOrchestratorStatus()
			return true // Continue waiting
		},
		p.ctx,
	)

	fmt.Printf("\n  ✓ Cycle #%d completed\n", p.cycleNumber)
}
```

**What happens**:

* Waits for the cycle duration (default: 5 minutes)
* Shows a progress bar
* Every 10 seconds displays orchestrator status
* Shows operations count, success, and errors
* User can see live performance

---

#### 7. Safety Limits Check
 
> Lines `644–675`  


```go
func (p *AdaptiveOrchestratorPreset) checkSafetyLimits() bool {
	// Check max loss
	if p.totalProfit < -p.BaseRiskAmount*5 {
		fmt.Printf("\n  🛑 STOP: Total loss exceeds $%.2f\n", p.BaseRiskAmount*5)
		return true
	}

	// Check daily loss (the day starts again at local midnight)
	currentBalance, err := p.loadedBalance()
	if err != nil {
		fmt.Printf("\n  ⚠️  Safety check skipped: %v\n", err)
		return false
	}
	if y, m, d := time.Now().Date(); !(y == p.dailyStartDate.Year() && m == p.dailyStartDate.Month() && d == p.dailyStartDate.Day()) {
		p.dailyStartBalance = currentBalance
		p.dailyStartDate = time.Now()
	}
	dailyProfit := currentBalance - p.dailyStartBalance

	if dailyProfit < -p.MaxDailyLoss {
		fmt.Printf("\n  🛑 STOP: Daily loss limit ($%.2f) reached\n", p.MaxDailyLoss)
		return true
	}

	// Check daily profit target
	if dailyProfit > p.MaxDailyProfit {
		fmt.Printf("\n  ✓ SUCCESS: Daily profit target ($%.2f) reached!\n", p.MaxDailyProfit)
		return true
	}

	return false
}
```

**What happens**:

* **Total loss check**: Stops if loss > $100 (5 × $20 base risk)
* **Daily loss check**: Stops if the loss since the start of the day (or of `Execute()`) > $500; the baseline moves to the current balance on a new local date
* **Daily profit check**: Stops if the profit since the start of the day > $1000 (success)
* A balance read error skips the check (it is not mistaken for a loss)
* Returns `true` if any limit is breached; `runCycles` then ends, and `Execute()` stops the orchestrators, shows the final results and returns

---

#### 8. Stop All Orchestrators
 
> Lines `620–627`  


```go
func (p *AdaptiveOrchestratorPreset) stopAllOrchestrators() {
    for _, orch := range p.activeOrchestrators {
        if orch.IsRunning() {
            orch.Stop()
        }
    }
    p.activeOrchestrators = make([]orchestrators.Orchestrator, 0)
}
```

**What happens**:

* Iterates through all active orchestrators
* Stops each one if still running
* Clears the orchestrators list
* Called before starting a new cycle, and by `Execute()` when the preset stops
* Each `Stop()` waits for the orchestrator's check in progress (a GridTrader also deletes its pending orders)
* Ensures a clean slate for the next mode

---

## Mode Selection Priority

The preset uses **priority-based mode selection** to handle conflicting conditions:

### Priority Order (Highest to Lowest)

1. **Portfolio Mode** — if `EnablePortfolioMode = true` AND 2+ symbols have open positions

2. **Protection Mode** — if volatility > 50 points

3. **Grid Mode** — if volatility < 25 points

4. **Trending Mode** — if trend strength > 70%

5. **Managed Mode** — default fallback for medium volatility

### Why This Order?

* **Portfolio First** — multi-symbol management overrides single-symbol strategies
* **Protection Second** — safety takes priority in high volatility
* **Grid Third** — low-volatility grid comes before trend detection
* **Trending Fourth** — only if volatility is medium but strong trend exists
* **Managed Last** — balanced approach when no special conditions apply

---

## Volatility Calculation

### Simplified Method (Demo)

```go
spreadPoints := priceInfo.SpreadPips
estimatedVolatility := spreadPoints * 10
```

**Current Implementation**:

* Uses the current spread (in points) as a volatility proxy
* Multiplies it by 10
* **Simplified for demonstration purposes**

### Production Alternatives

For real trading, calculate volatility using:

1. **ATR (Average True Range)** — 14-period ATR from bars
2. **Historical Volatility** — standard deviation of price changes
3. **VWAP Deviation** — distance from volume-weighted average price

GoMT5's `MT5Account` / `MT5Service` / `MT5Sugar` have **no bar-history method**, so bars would have to come from another source. A measurement that works with the existing API is the range of sampled prices:

**Example with sampled prices**:

```go
// Sample the Bid once per second for one minute
high, low := 0.0, 0.0
for i := 0; i < 60; i++ {
    info, err := sugar.GetPriceInfo("EURUSD")
    if err == nil {
        if high == 0 || info.Bid > high {
            high = info.Bid
        }
        if low == 0 || info.Bid < low {
            low = info.Bid
        }
    }
    time.Sleep(time.Second)
}

// Range in points (the point size comes from the symbol: EURUSD 0.00001)
info, err := sugar.GetSymbolInfo("EURUSD")
if err != nil {
    return err
}
rangePoints := (high - low) / info.Point
fmt.Printf("1-minute range: %.0f points\n", rangePoints)
```

---

## Safety Features

### 1. Stop-Loss Protection

**Total Loss Limit**: Stops if loss > 5 × base risk

```go
if p.totalProfit < -p.BaseRiskAmount*5 {
    // Stop preset
}
```

**Default**: $100 total loss limit ($20 × 5)

---

### 2. Daily Loss Limit

**Maximum Daily Loss**: Stops if daily loss exceeds configured limit

```go
dailyProfit := currentBalance - dailyStartBalance
if dailyProfit < -p.MaxDailyLoss {
    // Stop preset
}
```

**Default**: $500 daily loss limit

---

### 3. Daily Profit Target

**Profit Goal**: Stops when daily profit target is reached

```go
if dailyProfit > p.MaxDailyProfit {
    // Stop preset (success!)
}
```

**Default**: $1000 daily profit target

---

### 4. Maximum Concurrent Orchestrators

**Orchestrator Limit**: at most `MaxConcurrentOrchestrators` orchestrators are started per cycle (default 2; 0 = no limit). A mode that wants more starts the first ones and reports the rest as not started.

```go
preset.MaxConcurrentOrchestrators = 2 // Managed Mode needs 2
```

**Number per mode**:

* Grid Mode: 1 orchestrator
* Managed Mode: 2 orchestrators
* Other modes: 1 orchestrator

---

### 5. Graceful Shutdown

**Signal Handling**: Catches Ctrl+C for clean exit - and a safety stop ends the preset the same way

```go
select {
case <-sigChan: // Ctrl+C / SIGTERM
case <-cyclesDone: // safety stop
}
p.cancel()
<-cyclesDone             // the running cycle returns
p.stopAllOrchestrators() // Clean shutdown
```

`Execute()` then returns, and the runner's deferred `sugar.Disconnect()` stops the terminal.

---

## ⚙️ Critical Warnings

### 1. Simplified Volatility Calculation

* Current implementation: `spread × 10`
* **Not accurate for production trading**
* Use ATR, historical volatility, or proper indicators
* Test calculations on demo accounts first

---

### 2. Trend Detection Is Placeholder

* Current: based on spread modulo
* **Not real trend analysis**
* Production: use moving averages, MACD, trend lines
* Implement proper technical analysis

---

### 3. Continuous Operation Risks

* Preset runs in an infinite loop until stopped
* Can accumulate losses if market analysis fails
* Monitor actively during first runs
* Set conservative safety limits

---

### 4. Mode Switching Overhead

* Stops and starts orchestrators every cycle
* Stopping a GridTrader deletes its pending orders; a new RiskManager starts with a new peak/daily baseline each cycle; a new PositionScaler may open a base position again if the symbol has none (other `Stop()` calls do not close positions)
* Consider longer cycle durations (10–15 minutes)
* Balance adaptation speed vs. stability

---

### 5. Multiple Orchestrators Complexity

* Managed Mode runs 2 orchestrators simultaneously
* Ensure they don’t conflict (e.g., both modifying the same position)
* The current implementation does not coordinate them: the TrailingStopManager moves the SL of every open position (including the PositionScaler's scale-ins, which get their own SL)

---


## Best Practices

### Start with Demo Account

* Run the preset on a demo account for a minimum of 1 week
* Observe mode switches and orchestrator behavior
* Verify volatility and trend detection accuracy
* Adjust thresholds based on your symbols

---

### Customize Thresholds for Your Market

* Default thresholds (25 pts, 50 pts) are tuned for EURUSD
* Volatile pairs (e.g. GBPJPY) need higher thresholds
* Stable pairs (e.g. EURCHF) need lower thresholds
* Always test on your specific symbols

---

### Monitor First 10 Cycles Manually

* Watch mode selection logic closely
* Verify the appropriate orchestrator is chosen
* Check if modes switch too frequently
* Adjust `CycleDuration` if needed

---

### Set Conservative Safety Limits

* Start with low `BaseRiskAmount` ($10–$20)
* Set tight `MaxDailyLoss` (1–2% of balance)
* Use realistic `MaxDailyProfit` (2–3% of balance)
* Tighten limits until you trust the system

---

### Improve Market Analysis

* Replace spread-based volatility with ATR
* Implement proper trend detection (MA crossovers)
* Add time-of-day filters (avoid low-liquidity hours)
* Consider fundamental factors (news events)

---

### Log All Mode Switches

* Record why each mode was selected
* Track mode effectiveness (profit per mode)
* Identify which modes work best
* Disable underperforming modes

---

### Regular Performance Review

* Check `GetCycleNumber()` and `GetTotalProfit()`
* Calculate average profit per cycle
* Identify loss-making modes
* Optimize orchestrator configurations

---

## 🔧 When to Use This Preset

### Ideal Use Cases

**Set-and-Forget Trading**

* Want automated strategy that adapts to markets
* Don’t want to manually switch strategies
* Available to monitor but not actively trade

**Multi-Strategy Testing**

* Want to compare different orchestrators
* Test which strategies work in which conditions
* Learn orchestrator behavior

**24/7 Automated Trading**

* Running unattended trading system
* Need adaptive approach for changing markets
* Built-in safety stops for protection

**Learning Advanced Patterns**

* Understanding orchestrator composition
* Learning market condition detection
* Building your own adaptive systems

---

### Avoid Using When

**High-Frequency Trading**

* 5-minute cycles are too slow for HFT
* Mode switching overhead too high
* Use single specialized orchestrator instead

**Very Specific Strategy**

* You already have one proven strategy
* No need for adaptation
* Use that single orchestrator directly

**Small Account (< $500)**

* Preset complexity not worth it for small capital
* Safety limits may stop prematurely
* Use simpler single-orchestrator approach

**During Major News Events**

* Market analysis may fail in extreme volatility
* Unpredictable mode switching
* Disable preset or use manual trading

---

## Monitoring and Control

### Checking Status During Execution

The preset displays live status every 10 seconds during monitoring:

```
  [15:04:05] Active orchestrators: 2
     - Trailing Stop Manager: 5 ops (4 success, 1 errors)
     - Position Scaler: 3 ops (3 success, 0 errors)
```

(The names come from each orchestrator's `GetStatus().Name`.)

---

### Accessing Metrics Programmatically

```go
// Get current cycle number
cycleNum := preset.GetCycleNumber()

// Get total profit / loss
totalProfit := preset.GetTotalProfit()

// Get active orchestrators count
activeCount := preset.GetActiveOrchestratorCount()

// Get orchestrator details
active := preset.GetActiveOrchestrators()
for _, orch := range active {
    status := orch.GetStatus()
    fmt.Printf("%s: %d operations\n", status.Name, status.SuccessCount)
}
```

---

### 💡 Final Results Display

When the preset stops (Ctrl+C or safety limit), it shows a summary:

```
+============================================================+
|  FINAL RESULTS                                            |
+============================================================+
  Initial Balance: $10000.00
  Final Balance: $10350.00
  Total Profit/Loss: $350.00
  Cycles Completed: 12
  Average Per Cycle: $29.17
+============================================================+

(`Average Per Cycle` is left out when no cycle has completed.)
```

---

## Troubleshooting

### Problem: Preset Never Enters Grid Mode

**Solution**: Check volatility threshold

* Current: `LowVolatilityThreshold = 25.0`
* Volatility = spread in points × 10, so a spread of 2.5 points or more never selects Grid Mode — increase the threshold
* Try: `preset.LowVolatilityThreshold = 40.0`

---

### Problem: Switches Modes Every Cycle (Unstable)

**Solution**: Widen volatility bands

* Create a buffer zone between thresholds
* Example: Low < 20 pts, High > 60 pts (40 pt gap)
* Or increase `CycleDuration` to 10–15 minutes

---

### Problem: Always in Protection Mode

**Solution**: Increase high volatility threshold

* Current: `HighVolatilityThreshold = 50.0`
* Increase to 80–100 for volatile pairs
* Example: `preset.HighVolatilityThreshold = 80.0`

---

### Problem: Hit Max Daily Loss Quickly

**Solution**: Reduce risk or increase limit

* Lower the orchestrators' lot sizes in the `execute...Mode` functions
* Or increase `MaxDailyLoss` from $500 to $1000 (`BaseRiskAmount` only sets the total-loss stop at 5 × its value)
* Verify orchestrator configurations are conservative

---

### Problem: Orchestrator Errors During Execution

**Solution**: Check individual orchestrator logs

* Review error messages in status display
* Verify the symbol is tradable
* Check account has sufficient margin
* Test individual orchestrators separately first

---

## Improving the Preset

### Custom Market Analysis

Replace the simplified volatility calculation. GoMT5 has no bar-history method, so a replacement has to work from prices you collect yourself (or from an external data source). Sketch of a change inside `analyzeMarketConditions()`, using a price history that you keep (for example filled from `GetPriceInfo` samples or `GetService().StreamTicks`):

```go
// p.bids []float64 would be a new field you fill with sampled Bid prices
if len(p.bids) >= 2 {
    high, low := p.bids[0], p.bids[0]
    for _, b := range p.bids {
        high = math.Max(high, b)
        low = math.Min(low, b)
    }
    condition.VolatilityPoints = (high - low) / p.point // p.point: a field you fill from sugar.GetSymbolInfo(p.Symbol).Point

    // Simple trend: last price vs first price, scaled to -1..1
    first, last := p.bids[0], p.bids[len(p.bids)-1]
    if high > low {
        condition.TrendStrength = (last - first) / (high - low)
    }
}

// ...rest of mode selection logic
```

---

### Adding Time-of-Day Filters

Sketch of code you would add to your copy of `16_AdaptiveOrchestratorPreset.go` (package `presets`; none of this exists in the shipped preset):

```go
// Don’t trade during low-liquidity hours
func (p *AdaptiveOrchestratorPreset) isTradingHours() bool {
    now := time.Now().UTC()
    hour := now.Hour()

    // Avoid Asian session (22:00–02:00 UTC)
    if hour >= 22 || hour < 2 {
        return false
    }

    return true
}

// In executeCycle():
if !p.isTradingHours() {
    fmt.Println("Outside trading hours, skipping cycle")
    return
}
```

---

### Mode-Specific Optimization

Sketch (not in the shipped code; `MarketMode` and `GridMode` are the preset's own types, so this belongs inside package `presets`):

```go
// Track performance per mode
type ModePerformance struct {
    Mode          MarketMode
    TotalCycles   int
    TotalProfit   float64
    SuccessRate   float64
}

// Disable underperforming modes
if modePerf[GridMode].SuccessRate < 0.4 {
    // Skip Grid Mode for this session
    continue
}
```


---

## Integration with Other Systems

### With External Signals

Sketch for your copy of the preset (`fetchSignalFromAPI` stands for your own signal source):

```go
// Add signal integration
type ExternalSignal struct {
    Symbol    string
    Direction string // "BUY" or "SELL"
    Strength  float64
}

func (p *AdaptiveOrchestratorPreset) checkExternalSignals() *ExternalSignal {
    // Check your signal provider
    signal := fetchSignalFromAPI()
    return signal
}

// Modify mode selection
if signal := p.checkExternalSignals(); signal != nil {
    if signal.Direction == "BUY" && signal.Strength > 0.8 {
        condition.Mode = TrendingMode
        condition.Reason = "Strong external buy signal"
    }
}
```

---

### With Custom Risk Management

To use stricter limits, change the configuration lines of `executeProtectionMode` in your copy of the file. The shipped code is:

```go
	config := orchestrators.DefaultRiskManagerConfig()
	config.MaxDrawdownPercent = 5.0
	config.EnableAutoClose = true
```

A stricter variant (starting from the defaults keeps the other limits; in a `RiskManagerConfig` literal a missing limit is 0 = off, and a missing `CheckInterval` becomes 5 s):

```go
	config := orchestrators.DefaultRiskManagerConfig()
	config.MaxDrawdownPercent = 3.0 // Stricter limit
	config.DailyLossLimit = 200.0   // Lower daily limit
	config.MinMarginLevel = 200.0   // Higher margin safety
	config.EnableAutoClose = true
```

---

## Summary

**Adaptive Orchestrator Preset** is an intelligent multi-strategy trading system:

* **Adaptive Selection** — automatically chooses the best orchestrator based on market conditions
* **Five Market Modes** — Grid, Managed, Protection, Portfolio, Trending
* **Continuous Operation** — runs in cycles with automatic adaptation
* **Built-in Safety** — total-loss stop and daily loss/profit limits checked after each cycle
* **Highly Customizable** — thresholds and configurations are adjustable

**Remember**:

* Market analysis is simplified (improve for production use)
* Test extensively on demo accounts
* Monitor first runs closely
* Customize thresholds for your symbols
* Start with conservative limits

**The goal**: Show how to build an adaptive trading system that selects a strategy for current market conditions — a starting point to improve, not a *set-and-forget* system.
