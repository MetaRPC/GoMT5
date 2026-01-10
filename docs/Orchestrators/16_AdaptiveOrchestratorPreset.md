
# 🧩 Adaptive Orchestrator Preset — Intelligent Multi-Strategy System

## Description

**Adaptive Orchestrator Preset** is an intelligent trading system that automatically selects and executes the most appropriate orchestrator based on real-time market condition analysis.

It combines **all five orchestrators** into a single adaptive system that dynamically switches strategies every cycle.

### Strategy Principle

Continuously analyzes market conditions and adapts the trading approach:

* **Market Analysis** — calculates volatility, detects trends, and counts active positions
* **Strategy Selection** — chooses the best orchestrator for current conditions
* **Dynamic Execution** — runs the selected orchestrator in 5-minute cycles
* **Continuous Adaptation** — repeats analysis and switches strategies as markets change
* **Safety Controls** — built-in stop-loss, daily limits, and risk management

**File:** `examples/demos/pesets/16_AdaptiveOrchestratorPreset.go`

> **IMPORTANT DISCLAIMER — EDUCATIONAL EXAMPLE ONLY**
>
> This preset demonstrates advanced orchestrator composition and adaptive strategy selection.
> **NOT production-ready!**
> Market analysis is simplified for demonstration purposes. Always test extensively on demo accounts and customize for your specific market knowledge and risk tolerance.

---

## What Is an Adaptive Preset?

### Simple Explanation

Instead of running one fixed strategy, the **Adaptive Orchestrator Preset** acts like an intelligent trading manager that:

1. Looks at the market every 5 minutes
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

* **MT5 Suga** — high-level convenience API
* **All 5 Orchestrators** — GridTrade, TrailingStopManage, PositionScale, RiskManage, PortfolioRebalance
* **Helpers** — progress bar and utility helpers

---

## 💡 Market Modes and Strategy Selection

### Mode 1: Grid Mode (Low Volatility)

**Condition:** Volatility < 25 points

**Why:** Range-bound markets are perfect for grid strategies.

**Orchestrator Used:** GridTrade

**Configuration:**

```go
GridTradeConfig{
    Symbol:        "EURUSD",
    GridSize:      5,      // 5 grid levels
    GridStep:      100,    // 100 points spacing
    LotSize:       0.01,   // Micro-lots
    MaxPositions:  10,
    CheckInterval: 5 * time.Second,
}
```

**Expected Behavior:**

* Places a grid of buy/sell orders
* Profits from small price oscillations
* Works best in sideways markets

---

### Mode 2: Managed Mode (Medium Volatility)

**Condition:** 25–50 points volatility

**Why:** Normal conditions suit position management with scaling.

**Orchestrators Used:**

1. TrailingStopManage (profit protection)
2. PositionScale (pyramiding)

**Configuration:**

```go
// Trailing Stop
TrailingStopConfig{
    TrailingDistance: 200, // 200 points trail
    ActivationProfit: 300, // Activate at 300 points profit
    UpdateInterval:   2 * time.Second,
}

// Position Scaling
PositionScaleConfig{
    Mode:            Pyramiding, // Add to winners
    TriggerDistance: 200,        // 200 points trigger
    MaxScales:       3,          // Max 3 additions
}
```

**Expected Behavior:**

* Protects profits with trailing stops
* Scales into winning positions
* Balanced approach for normal markets

---

### Mode 3: Protection Mode (High Volatility)

**Condition:** Volatility > 50 points

**Why:** Volatile markets require strict risk control.

**Orchestrator Used:** RiskManage

**Configuration:**

```go
RiskManageConfig{
    MaxDrawdownPercent: 5.0,  // 5% max drawdown
    EnableAutoClose:    true, // Auto-close on breach
    CheckInterval:      5 * time.Second,
}
```

**Expected Behavior:**

* Continuously monitors risk limits
* Closes positions if drawdown exceeds 5%
* Defensive mode to prevent blow-ups

---

### Mode 4: Portfolio Mode (Multi-Symbol)

**Condition:** Two or more different symbols have open positions

**Why:** Multiple symbols require balanced exposure.

**Orchestrator Used:** PortfolioRebalance

**Configuration:**

```go
PortfolioRebalanceConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.25,
        "GBPUSD": 0.25,
        "USDJPY": 0.25,
        "XAUUSD": 0.25,
    },
    TotalExposure:     10000.0,
    RebalanceThreshold: 10.0,
}
```


**Expected Behavior**:

* Maintains balanced allocation
* Rebalances when deviation > 10%
* Diversification management

---

### Mode 5: TRENDING MODE (Strong Directional Movement)

**Condition**: Trend strength > 70%

**Why**: Strong trends are perfect for pyramiding strategies

**Orchestrator Used**: PositionScale (pyramiding mode)

**Configuration**:

```go
PositionScaleConfig{
    Mode:            Pyramiding, // Add to winners
    TriggerDistance: 200,        // 200 points trigger
    MaxScales:       3,          // Max 3 scales
}
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
| `BaseRiskAmount`             | float64       | `20.0`     | Risk per trade in dollars      |
| `LowVolatilityThreshold`     | float64       | `25.0`     | Max points for grid mode       |
| `HighVolatilityThreshold`    | float64       | `50.0`     | Min points for protection mode |
| `CycleDuration`              | time.Duration | `5m`       | Duration of one cycle          |
| `EnablePortfolioMode`        | bool          | `false`    | Enable multi-symbol management |
| `PortfolioSymbols`           | []string      | 4 symbols  | Symbols for portfolio mode     |
| `MaxConcurrentOrchestrators` | int           | `2`        | Max orchestrators at once      |
| `MaxDailyLoss`               | float64       | `500.0`    | Maximum daily loss allowed     |
| `MaxDailyProfit`             | float64       | `1000.0`   | Daily profit target            |

### Configuration Example

```go
preset := presets.NewAdaptiveOrchestratorPreset(suga)

// Customize parameters
preset.Symbol = "EURUSD"
preset.BaseRiskAmount = 20.0
preset.LowVolatilityThreshold = 25.0
preset.HighVolatilityThreshold = 50.0
preset.CycleDuration = 5 * time.Minute
preset.EnablePortfolioMode = false
preset.MaxDailyLoss = 500.0
preset.MaxDailyProfit = 1000.0

// Execute
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

All commands launch the **Adaptive Orchestrator Preset** with default configuration.
To stop execution, press `Ctrl+C`.

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
|   1. Portfolio mode? (2+ symbols)
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
|   - Grid Mode → Start GridTrade
|   - Managed Mode → Start TrailingStop + PositionScale
|   - Protection Mode → Start RiskManage
|   - Portfolio Mode → Start PortfolioRebalance
|   - Trending Mode → Start PositionScale (pyramiding)
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
|
+-- STEP 7: PAUSE
|   - Wait 30 seconds before next cycle
|
+-- REPEAT CYCLE
|
ON STOP (Ctrl+C or safety limit):
- Stop all active orchestrators
- Show final results
- Return total profit
|
END
```


```go
func (p *AdaptiveOrchestratorPreset) Execute() (float64, error) {
    // Get starting balance
    p.initialBalance, err := p.suga.GetBalance()
    p.dailyStartBalance = p.initialBalance

    // Setup context and signal handling
    p.ctx, p.cancel = context.WithCancel(context.Background())
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // Start main cycle loop
    go p.runCycles()

    // Wait for interrupt
    <-sigChan
}
```

**What happens**:

* Gets starting balance for profit calculation
* Creates context for graceful shutdown
* Sets up Ctrl+C handling
* Launches main cycle loop in a goroutine
* Blocks until user stops execution (Ctrl+C)

---

#### 2. Main Cycle Loop
  
> Lines `307–331`  


```go
func (p *AdaptiveOrchestratorPreset) runCycles() {
    for {
        select {
        case <-p.ctx.Done():
            return
        default:
            p.cycleNumber++
            p.executeCycle()

            // Check safety limits
            if p.checkSafetyLimits() {
                p.cancel()
                return
            }

            // Pause 30 seconds
            time.Sleep(30 * time.Second)
        }
    }
}
```

**What happens**:

* Infinite loop running cycles
* Increments cycle counter
* Executes one trading cycle
* Checks safety limits after each cycle
* Stops if limits are breached
* Pauses 30 seconds between cycles

---

#### 3. Execute Single Cycle
  
> Lines `334–387`  


```go
func (p *AdaptiveOrchestratorPreset) executeCycle() {
    // Step 1: Analyze market
    condition, err := p.analyzeMarketConditions()

    // Step 2: Stop previous orchestrators
    p.stopAllOrchestrators()

    // Step 3: Execute appropriate mode
    cycleStartBalance, _ := p.suga.GetBalance()

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
    }

    // Step 4: Monitor cycle
    p.monitorCycle()

    // Step 5: Calculate profit
    cycleEndBalance, _ := p.suga.GetBalance()
    cycleProfit := cycleEndBalance - cycleStartBalance
    p.totalProfit += cycleProfit
}
```

**What happens**:

* Analyzes current market conditions
* Stops any orchestrators from the previous cycle
* Records starting balance
* Selects and executes the appropriate mode
* Monitors execution for the cycle duration (5 minutes)
* Calculates profit for this cycle
* Updates total profit

---

#### 4. Market Analysis
  
> Lines `394–447`  


```go
func (p *AdaptiveOrchestratorPreset) analyzeMarketConditions() (*MarketCondition, error) {
    // Get price info
    priceInfo, err := p.suga.GetPriceInfo(p.Symbol)

    // Calculate volatility from spread
    spreadPoints := priceInfo.SpreadPips
    estimatedVolatility := spreadPoints * 10

    // Get active positions
    positions, _ := p.suga.GetOpenPositions()

    // Count unique symbols
    symbolMap := make(map[string]bool)
    for _, pos := range positions {
        symbolMap[pos.Symbol] = true
    }

    // Determine mode (priority order)
    if p.EnablePortfolioMode && len(symbolMap) >= 2 {
        condition.Mode = PortfolioMode
    } else if estimatedVolatility > p.HighVolatilityThreshold {
        condition.Mode = ProtectionMode
    } else if estimatedVolatility < p.LowVolatilityThreshold {
        condition.Mode = GridMode
    } else if math.Abs(condition.TrendStrength) > 0.7 {
        condition.Mode = TrendingMode
    } else {
        condition.Mode = ManagedMode
    }
}
```

**What happens**:

* Fetches current price info using `GetPriceInfo()`
* Calculates volatility as `spread × 10` (simplified estimate)
* Gets all open positions
* Counts unique symbols for portfolio detection
* Applies mode selection logic with priority:

  1. **Portfolio** if 2+ symbols active

  2. **Protection** if volatility > 50 pts

  3. **Grid** if volatility < 25 pts

  4. **Trending** if trend strength > 70%

  5. **Managed** otherwise (default)

---

#### 5. Mode Execution Example: Grid Mode
  
> Lines `454–477`  


```go
func (p *AdaptiveOrchestratorPreset) executeGridMode(condition *MarketCondition) {
    config := orchestrators.GridTradeConfig{
        Symbol:        p.Symbol,
        GridSize:      5,
        GridStep:      100,
        LotSize:       0.01,
        MaxPositions:  10,
        CheckInterval: 5 * time.Second,
    }

    gridTrade := orchestrators.NewGridTrade(p.suga, config)
    if err := gridTrade.Start(); err != nil {
        return
    }

    p.activeOrchestrators = append(p.activeOrchestrators, gridTrade)
}
```

**What happens**:

* Creates configuration for GridTrade
* Instantiates a new GridTrade orchestrator
* Starts the orchestrator
* Adds it to the active orchestrators list
* Same pattern is used for other modes (ManagedMode, ProtectionMode, etc.)

---

#### 6. Cycle Monitoring
 
> Lines `562–580`  


```go
func (p *AdaptiveOrchestratorPreset) monitorCycle() {
    totalSeconds := int(p.CycleDuration.Seconds())

    helpers.WaitWithProgressBarAndCallback(
        totalSeconds,
        fmt.Sprintf("Cycle #%d monitoring", p.cycleNumber),
        10*time.Second, // Update every 10 seconds
        func() bool {
            p.showOrchestratorStatus()
            return true
        },
        p.ctx,
    )
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
 
> Lines `609–632`  


```go
func (p *AdaptiveOrchestratorPreset) checkSafetyLimits() bool {
    // Check total loss
    if p.totalProfit < -p.BaseRiskAmount*5 {
        fmt.Printf("STOP: Total loss exceeds $%.2f\n", p.BaseRiskAmount*5)
        return true
    }

    // Check daily loss
    currentBalance, _ := p.suga.GetBalance()
    dailyProfit := currentBalance - p.dailyStartBalance

    if dailyProfit < -p.MaxDailyLoss {
        fmt.Printf("STOP: Daily loss limit ($%.2f) reached\n", p.MaxDailyLoss)
        return true
    }

    // Check daily profit target
    if dailyProfit > p.MaxDailyProfit {
        fmt.Printf("SUCCESS: Daily profit target ($%.2f) reached!\n", p.MaxDailyProfit)
        return true
    }

    return false
}
```

**What happens**:

* **Total loss check**: Stops if loss > $100 (5 × $20 base risk)
* **Daily loss check**: Stops if daily loss > $500
* **Daily profit check**: Stops if daily profit > $1000 (success)
* Returns `true` if any limit is breached, stopping the preset

---

#### 8. Stop All Orchestrators
 
> Lines `599–606`  


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
* Called before starting a new cycle
* Ensures a clean slate for the next mode

---

## Mode Selection Priority

The preset uses **priority-based mode selection** to handle conflicting conditions:

### Priority Order (Highest to Lowest)

1. **Portfolio Mode** — if `EnablePortfolioMode = true` AND 2+ symbols are active

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

* Uses spread as a volatility proxy
* Multiplies by 10 for point estimation
* **Simplified for demonstration purposes**

### Production Alternatives

For real trading, calculate volatility using:

1. **ATR (Average True Range)** — 14-period ATR from bars
2. **Historical Volatility** — standard deviation of price changes
3. **VWAP Deviation** — distance from volume-weighted average price

**Example with bars**:

```go
// Get recent bars
bars, _ := suga.GetBars(symbol, timeframe, 100)

// Calculate ATR
atr := calculateATR(bars, 14)

// Use ATR as volatility
if atr < 20 {
    mode = GridMode
} else if atr > 50 {
    mode = ProtectionMode
}
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

**Orchestrator Limit**: Controls maximum orchestrators running at once

```go
MaxConcurrentOrchestrators: 2
```

**Current Usage**:

* Grid Mode: 1 orchestrator
* Managed Mode: 2 orchestrators
* Other modes: 1 orchestrator

---

### 5. Graceful Shutdown

**Signal Handling**: Catches Ctrl+C for clean exit

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

<-sigChan // Wait for Ctrl+C
p.stopAllOrchestrators() // Clean shutdown
```

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
* May close positions prematurely
* Consider longer cycle durations (10–15 minutes)
* Balance adaptation speed vs. stability

---

### 5. Multiple Orchestrators Complexity

* Managed Mode runs 2 orchestrators simultaneously
* Ensure they don’t conflict (e.g., both modifying the same position)
* Current implement


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
- TrailingStopManage: 5 ops (4 success, 1 error)
- PositionScale: 3 ops (3 success, 0 errors)
```

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
orchestrators := preset.GetActiveOrchestrators()
for _, orch := range orchestrators {
    status := orch.GetStatus()
    fmt.Printf("%s: %d operations\n", status.Name, status.SuccessCount)
}
```

---

### 💡 Final Results Display

When the preset stops (Ctrl+C or safety limit), it shows a summary:

```
+============================================================+
| FINAL RESULTS                                              |
+============================================================+
Initial Balance:      $10000.00
Final Balance:        $10350.00
Total Profit / Loss:  $350.00
Cycles Completed:     12
Average Per Cycle:    $29.17
+============================================================+
```

---

## Troubleshooting

### Problem: Preset Never Enters Grid Mode

**Solution**: Check volatility threshold

* Current: `LowVolatilityThreshold = 25.0`
* If your market always has spread > 2.5 pips, increase the threshold
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

* Lower `BaseRiskAmount` from $20 to $10
* Or increase `MaxDailyLoss` from $500 to $1000
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

Replace simplified volatility calculation:

```go
func (p *AdaptiveOrchestratorPreset) analyzeMarketConditions() (*MarketCondition, error) {
    // Get bars for analysis
    bars, _ := p.suga.GetBars(p.Symbol, pb.Timeframe_M15, 100)

    // Calculate ATR
    atr := calculateATR(bars, 14)

    // Calculate moving averages for trend
    ma20 := calculateMA(bars, 20)
    ma50 := calculateMA(bars, 50)

    // Detect trend
    trendStrength := 0.0
    if ma20 > ma50 {
        trendStrength = (ma20 - ma50) / ma50
    } else {
        trendStrength = (ma20 - ma50) / ma50
    }

    // Use ATR for volatility
    condition.VolatilityPoints = atr
    condition.TrendStrength = trendStrength

    // ...rest of mode selection logic
}
```

---

### Adding Time-of-Day Filters

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

```go
// Override default RiskManage configuration
func (p *AdaptiveOrchestratorPreset) executeProtectionMode(condition *MarketCondition) {
    config := orchestrators.RiskManageConfig{
        MaxDrawdownPercent: 3.0,  // Stricter limit
        DailyLossLimit:     200.0, // Lower daily limit
        MinMarginLevel:     200.0, // Higher margin safety
        EnableAutoClose:    true,
    }

    riskManage := orchestrators.NewRiskManage(p.suga, config)
    riskManage.Start()
}
```

---

## Summary

**Adaptive Orchestrator Preset** is an intelligent multi-strategy trading system:

* **Adaptive Selection** — automatically chooses the best orchestrator based on market conditions
* **Five Market Modes** — Grid, Managed, Protection, Portfolio, Trending
* **Continuous Operation** — runs in cycles with automatic adaptation
* **Built-in Safety** — multiple safety stops and daily limits
* **Highly Customizable** — thresholds and configurations are adjustable

**Remember**:

* Market analysis is simplified (improve for production use)
* Test extensively on demo accounts
* Monitor first runs closely
* Customize thresholds for your symbols
* Start with conservative limits

**The goal**: Create a complete adaptive trading system that automatically selects the right strategy for current market conditions, providing true *set-and-forget* automated trading with intelligent risk management.
