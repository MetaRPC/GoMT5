# RiskManager - Account Protection System

## Description

**🧩 RiskManager** is an automated account guardian that monitors risk metrics 24/7 and automatically intervenes when danger limits are exceeded. It **PREVENTS ACCOUNT BLOW-UPS** by enforcing strict risk rules: drawdown limits, daily loss caps, margin monitoring, and position limits. The "safety net" that saves you from catastrophic losses.

**Strategy Principle**: Continuously monitors critical account parameters and takes protective action when thresholds are breached:

- **Drawdown Protection**: Monitors both percentage and absolute drawdown from peak balance
- **Daily Loss Limit**: Enforces maximum daily loss caps (prevents revenge trading)
- **Daily Profit Target**: Locks in good days by stopping trading at profit goal
- **Margin Safety**: Monitors margin levels and prevents margin calls
- **Position Limits**: Controls maximum position counts and per-symbol exposure
- **Emergency Actions**: Automatically closes positions or blocks new trades when limits hit

**File**: `examples/demos/orchestrators/14_risk_manager.go`

> **IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY**
>
> This orchestrator is a demonstration showing how GoMT5 methods function and combine into automated risk management systems. **NOT a production-ready trading strategy!** Always test on demo accounts first and adjust limits to your specific risk tolerance.

---

## 🟢 Architecture

```
RISK MANAGER ORCHESTRATOR
    |
MT5Sugar Instance
    |
    |  |  |  |
GetEquity  GetBalance  GetMarginLevel  GetOpenPositions
(monitoring) (tracking)    (safety)        (limits)
    |
CloseAllPositions / ClosePosition
(emergency intervention)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetEquity`, `GetBalance`, `GetMarginLevel`, `GetOpenPositions`, `CloseAllPositions`, `ClosePosition`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo` for position data structures

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `MaxDrawdownPercent` | float64 | `10.0` | Maximum account drawdown percentage from peak |
| `MaxDrawdownAbsolute` | float64 | `1000.0` | Maximum absolute drawdown amount |
| `DailyLossLimit` | float64 | `500.0` | Maximum daily loss allowed |
| `DailyProfitTarget` | float64 | `1000.0` | Stop trading after hitting this profit |
| `MinMarginLevel` | float64 | `150.0` | Minimum margin level percentage |
| `MaxMarginUsed` | float64 | `80.0` | Maximum margin usage percentage |
| `MaxOpenPositions` | int | `20` | Maximum concurrent positions |
| `MaxSymbolExposure` | int | `5` | Maximum positions per symbol |
| `MaxPositionSize` | float64 | `1.0` | Maximum lot size per position |
| `CheckInterval` | time.Duration | `5s` | How often to check risk metrics |
| `EnableAutoClose` | bool | `true` | Automatically close positions on breach |
| `EnableTradeBlocking` | bool | `true` | Block new trades when limits hit |

### Configuration Example

```go
config := orchestrators.RiskManagerConfig{
    MaxDrawdownPercent:  10.0,              // 10% max drawdown from peak
    MaxDrawdownAbsolute: 1000.0,            // $1000 max drawdown
    DailyLossLimit:      500.0,             // $500 daily loss limit
    DailyProfitTarget:   1000.0,            // $1000 daily profit target
    MinMarginLevel:      150.0,             // 150% min margin level
    MaxMarginUsed:       80.0,              // 80% max margin usage
    MaxOpenPositions:    20,                // Max 20 positions
    MaxSymbolExposure:   5,                 // Max 5 per symbol
    MaxPositionSize:     1.0,               // Max 1.0 lot
    CheckInterval:       5 * time.Second,   // Check every 5 seconds
    EnableAutoClose:     true,              // Auto-close on breach
    EnableTradeBlocking: true,              // Block trades on breach
}

riskManager := orchestrators.NewRiskManager(sugar, config)
```

---

## How to Run

You can execute this orchestrator using several command variations:

```bash
# Navigate to examples directory
cd examples/demos

# Option 1: By number
go run main.go 14

# Option 2: By keyword
go run main.go risk
```

Both commands will launch the **RiskManager** orchestrator with the configuration specified in `main.go -> RunOrchestrator_RiskManager()`.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Initialize baseline values:
  - Starting balance
  - Peak balance
  - Daily start balance
  - Trading blocked = false
  |
MONITORING LOOP (every CheckInterval):
  |
  Check if new day -> Reset daily counters
  |
  Get current account state:
    - Equity
    - Balance
    - Margin Level
    - Open positions
  |
  Update peak balance if equity higher
  |
  Calculate metrics:
    - Current drawdown = peak - equity
    - Drawdown % = (drawdown / peak) × 100
    - Today profit = balance - daily start balance
  |
  CHECK RISK LIMITS:
    |
    1. DRAWDOWN LIMITS:
       | If drawdown % >= MaxDrawdownPercent:
       |   -> Log CRITICAL event
       |   -> If EnableAutoClose: CLOSE ALL POSITIONS
       |
       | If drawdown $ >= MaxDrawdownAbsolute:
       |   -> Log CRITICAL event
       |   -> If EnableAutoClose: CLOSE ALL POSITIONS
    |
    2. DAILY LIMITS:
       | If today loss >= DailyLossLimit:
       |   -> Log CRITICAL event
       |   -> If EnableAutoClose: CLOSE ALL POSITIONS
       |   -> If EnableTradeBlocking: Block trading
       |
       | If today profit >= DailyProfitTarget:
       |   -> Log INFO event
       |   -> If EnableTradeBlocking: Block trading
    |
    3. MARGIN LIMITS:
       | If margin level < MinMarginLevel:
       |   -> Log CRITICAL event
       |   -> If EnableAutoClose: Close most losing position
    |
    4. POSITION LIMITS:
       | If open positions > MaxOpenPositions:
       |   -> Log WARNING event
       |
       | For each symbol:
       |   If positions > MaxSymbolExposure:
       |     -> Log WARNING event
  |
  Update status metrics
  |
REPEAT until stopped
  |
ON STOP:
  | Cancel monitoring loop
  | Print risk event summary
  |
END
```

### 🛠️ Step-by-step Description

#### 1. Initialization
 
> Lines `188–201`  


```go
func (r *RiskManager) initialize() error {
    // Get starting balance
    balance, err := r.sugar.GetBalance()
    if err != nil {
        return fmt.Errorf("failed to get balance: %w", err)
    }

    r.startingBalance = balance
    r.peakBalance = balance
    r.dailyStartBalance = balance
    r.tradingBlocked = false

    return nil
}
```

**What happens**:

- Gets current account balance as baseline
- Sets peak balance (tracks highest equity reached)
- Sets daily start balance (for daily P/L calculations)
- Initializes trading as unblocked
- All subsequent drawdown is measured from this peak

#### 2. Monitoring Loop
 
> Lines `204–216`  


```go
func (r *RiskManager) monitorLoop() {
    ticker := time.NewTicker(r.config.CheckInterval)
    defer ticker.Stop()

    for {
        select {
        case <-r.GetContext().Done():
            return
        case <-ticker.C:
            r.checkRiskLimits()
        }
    }
}
```

**Monitoring actions every CheckInterval**:

1. **Ticker triggers**: Every 5 seconds by default
2. **Context check**: Monitors for stop signal
3. **Risk checks**: Calls `checkRiskLimits()` to evaluate all parameters
4. **Graceful exit**: Stops cleanly when context is cancelled

#### 3. Daily Reset Check
  
> Lines `411–425`  


```go
func (r *RiskManager) checkDailyReset() {
    now := time.Now()
    if now.Day() != r.lastResetDate.Day() {
        balance, err := r.sugar.GetBalance()
        if err == nil {
            r.dailyStartBalance = balance
            r.tradingBlocked = false
            r.lastResetDate = now
        }
    }
}
```

**What happens**:

- Checks if current day differs from last reset day
- If new day: resets daily start balance to current balance
- Unblocks trading (allows new trading day to begin)
- Resets today's profit to zero (new calculations start)

#### 4. Drawdown Protection
 
> Lines `266–288`  


```go
func (r *RiskManager) checkDrawdownLimit(drawdownPercent, drawdownAbsolute float64) {
    // Check percentage drawdown
    if drawdownPercent >= r.config.MaxDrawdownPercent {
        r.logRiskEvent("MAX_DRAWDOWN_PERCENT", "CRITICAL",
            fmt.Sprintf("Drawdown %.1f%% exceeds limit %.1f%%", drawdownPercent, r.config.MaxDrawdownPercent),
            drawdownPercent, r.config.MaxDrawdownPercent)

        if r.config.EnableAutoClose {
            r.closeAllPositionsEmergency("Maximum drawdown exceeded")
        }
    }

    // Check absolute drawdown
    if drawdownAbsolute >= r.config.MaxDrawdownAbsolute {
        r.logRiskEvent("MAX_DRAWDOWN_ABSOLUTE", "CRITICAL",
            fmt.Sprintf("Drawdown $%.2f exceeds limit $%.2f", drawdownAbsolute, r.config.MaxDrawdownAbsolute),
            drawdownAbsolute, r.config.MaxDrawdownAbsolute)

        if r.config.EnableAutoClose {
            r.closeAllPositionsEmergency("Maximum absolute drawdown exceeded")
        }
    }
}
```

**Critical protection**:

- **Percentage drawdown**: Triggers when equity drops X% from peak
  - Example: Peak = $10,000, MaxDrawdownPercent = 10%
  - Trigger at equity = $9,000 (10% drop from peak)
- **Absolute drawdown**: Triggers when equity drops $X from peak
  - Example: Peak = $10,000, MaxDrawdownAbsolute = $1000
  - Trigger at equity = $9,000 ($1000 drop from peak)
- **Both are checked**: Whichever threshold is hit first triggers emergency action

#### 5. Daily Limits Protection
 
> Lines `291–320`  


```go
func (r *RiskManager) checkDailyLimits() {
    // Check daily loss limit
    if r.todayProfit < 0 && (-r.todayProfit) >= r.config.DailyLossLimit {
        r.logRiskEvent("DAILY_LOSS_LIMIT", "CRITICAL",
            fmt.Sprintf("Daily loss $%.2f exceeds limit $%.2f", -r.todayProfit, r.config.DailyLossLimit),
            -r.todayProfit, r.config.DailyLossLimit)

        if r.config.EnableAutoClose {
            r.closeAllPositionsEmergency("Daily loss limit exceeded")
        }

        if r.config.EnableTradeBlocking {
            r.tradingBlocked = true
        }
    }

    // Check daily profit target
    if r.config.DailyProfitTarget > 0 && r.todayProfit >= r.config.DailyProfitTarget {
        r.logRiskEvent("DAILY_PROFIT_TARGET", "INFO",
            fmt.Sprintf("Daily profit target $%.2f reached", r.config.DailyProfitTarget),
            r.todayProfit, r.config.DailyProfitTarget)

        if r.config.EnableTradeBlocking {
            r.tradingBlocked = true
        }
    }
}
```

**Daily protections**:

- **Daily loss limit**: Stops trading after losing too much today
  - Prevents "revenge trading" after bad day
  - Closes all positions and blocks new trades
- **Daily profit target**: Stops trading after hitting profit goal
  - Locks in good days
  - Only blocks new trades (keeps existing positions open)

#### 6. Margin Safety
  
> Lines `323–335`  


```go
func (r *RiskManager) checkMarginLimits(marginLevel float64) {
    // Check minimum margin level
    if marginLevel > 0 && marginLevel < r.config.MinMarginLevel {
        r.logRiskEvent("LOW_MARGIN_LEVEL", "CRITICAL",
            fmt.Sprintf("Margin level %.0f%% below minimum %.0f%%", marginLevel, r.config.MinMarginLevel),
            marginLevel, r.config.MinMarginLevel)

        if r.config.EnableAutoClose {
            // Close most losing position to free margin
            r.closeMostLosingPosition("Low margin level")
        }
    }
}
```

**Margin protection**:

- **Margin Level** = (Equity / Margin Used) × 100
- When margin level drops below threshold:
  - Closes most losing position to free up margin
  - Prevents broker margin call (typically at 50-100%)
- **Example**: MinMarginLevel = 150%
  - Equity = $10,000, Margin = $6,667
  - Margin Level = 150% (borderline)
  - If drops to 149%: close worst position

#### 7. Emergency Close All
  
> Lines `367–382`  


```go
func (r *RiskManager) closeAllPositionsEmergency(reason string) {
    closed, err := r.sugar.CloseAllPositions()
    if err != nil {
        r.IncrementError(fmt.Sprintf("emergency close failed: %v", err))
        return
    }

    r.UpdateMetrics(func(m *OrchestratorMetrics) {
        m.LastOperation = fmt.Sprintf("EMERGENCY: Closed %d positions - %s", closed, reason)
        m.OperationsTotal += closed
    })

    r.logRiskEvent("EMERGENCY_CLOSE", "CRITICAL",
        fmt.Sprintf("Closed all %d positions: %s", closed, reason),
        float64(closed), 0)
}
```

**Emergency action**:

- Immediately closes ALL open positions
- Cannot be undone
- Triggered by:
  - Maximum drawdown exceeded (% or $)
  - Daily loss limit exceeded
- Logs critical event for analysis
- Updates metrics with emergency action

---

## Visual Examples

### Scenario 1: Drawdown Protection in Action

**Initial State**:

- Account balance: $10,000
- MaxDrawdownPercent: 10%
- MaxDrawdownAbsolute: $1,000

```
TIME 0: Trading Starts

Balance:      $10,000
Peak Balance: $10,000
Drawdown:     $0 (0%)
Positions:    0
Status:       Monitoring

TIME 1: Good Trading Day

Balance:      $11,500
Peak Balance: $11,500 (new peak!)
Drawdown:     $0 (0%)
Positions:    3
Status:       Monitoring - all good

TIME 2: Market Turns Against Positions

Balance:      $10,800
Peak Balance: $11,500 (unchanged - peak is max)
Drawdown:     $700 (6.1%)
Positions:    5
Status:       Monitoring - within limits

TIME 3: Drawdown Approaching Limit

Balance:      $10,400
Peak Balance: $11,500
Drawdown:     $1,100 (9.6%)
Status:       WARNING - approaching 10% limit

TIME 4: DRAWDOWN LIMIT EXCEEDED

Balance:      $10,350
Peak Balance: $11,500
Drawdown:     $1,150 (10.0%)
TRIGGER:      MaxDrawdownPercent HIT!
ACTION:       EMERGENCY CLOSE ALL 5 POSITIONS
Result:       Trading stopped, account protected at $10,350
Prevented:    Further losses from uncontrolled drawdown

WITHOUT RISK MANAGER:
- Positions could continue losing
- Drawdown could reach 20-30% or more
- Account could blow up completely

WITH RISK MANAGER:
- Loss limited to 10% from peak
- $1,150 loss vs potential $3,000+ loss
- Account survives to trade another day
```

### Scenario 2: Daily Loss Limit

**Configuration**:

- DailyLossLimit: $500
- EnableAutoClose: true
- EnableTradeBlocking: true

```
TIME 0: Start of Day (Midnight Reset)

Daily Start Balance: $10,000
Current Balance:     $10,000
Today Profit:        $0
Trading Blocked:     false

TIME 1: Morning Trades (09:00)

Current Balance:     $9,800
Today Profit:        -$200
Status:              Within limits

TIME 2: Bad Trading Session (11:00)

Current Balance:     $9,600
Today Profit:        -$400
Status:              WARNING - approaching limit

TIME 3: DAILY LIMIT HIT (12:30)

Current Balance:     $9,500
Today Profit:        -$500
TRIGGER:             DailyLossLimit HIT!
ACTION:              1. CLOSE ALL POSITIONS
                     2. BLOCK NEW TRADES
Status:              Trading blocked for rest of day

TIME 4: Attempt to Trade (14:00)

Status:              Trading blocked - limit reached
Result:              Cannot open new positions
                     Prevents revenge trading

TIME 5: Next Day (Midnight Reset)

Daily Start Balance: $9,500 (updated to current)
Today Profit:        $0 (reset)
Trading Blocked:     false (unblocked)
Status:              New day, trading allowed again

BENEFIT:
- Prevented emotional "revenge trading"
- Limited daily loss to exactly $500
- Forced break to reassess strategy
- Protected remaining capital
```

### Scenario 3: Margin Level Protection

**Configuration**:

- MinMarginLevel: 150%
- EnableAutoClose: true

```
TIME 0: Normal Trading

Equity:       $10,000
Margin Used:  $5,000
Margin Level: 200% (healthy)
Positions:    5 positions
Status:       All systems normal

TIME 1: Positions Move Against

Equity:       $9,000
Margin Used:  $5,000
Margin Level: 180%
Status:       Still safe

TIME 2: Approaching Danger

Equity:       $8,000
Margin Used:  $5,000
Margin Level: 160%
Status:       Getting close to limit

TIME 3: MARGIN LEVEL LOW

Equity:       $7,500
Margin Used:  $5,000
Margin Level: 150%
TRIGGER:      MarginLevel = MinMarginLevel
ACTION:       Close most losing position

Position Details:
- Position 1: -$200
- Position 2: +$50
- Position 3: -$500 <- WORST LOSS
- Position 4: +$100
- Position 5: -$150

ACTION TAKEN: Close Position #3 (worst loser)

TIME 4: After Closing Worst Position

Equity:       $7,500
Margin Used:  $4,000 (reduced)
Margin Level: 187.5% (improved!)
Positions:    4 remaining
Status:       Margin level safe again

RESULT:
- Prevented broker margin call
- Closed only worst position
- Other positions remain open
- Can continue trading
```

---

## Used MT5Sugar Methods

### 1. GetEquity

```go
func (s *MT5Sugar) GetEquity() (float64, error)
```

**Purpose**: Gets current account equity (balance + floating P/L).

**Usage in RiskManager**:

- Called every CheckInterval
- Used to calculate current drawdown from peak
- Core metric for drawdown protection

**Returns**: Current equity or error

### 2. GetBalance

```go
func (s *MT5Sugar) GetBalance() (float64, error)
```

**Purpose**: Gets current account balance (closed positions only).

**Usage in RiskManager**:

- Used for daily profit/loss calculations
- Compared to dailyStartBalance to get today's performance
- Reset baseline at start of each day

**Returns**: Current balance or error

### 3. GetMarginLevel

```go
func (s *MT5Sugar) GetMarginLevel() (float64, error)
```

**Purpose**: Gets current margin level percentage.

**Usage in RiskManager**:

- Monitors risk of margin call
- Triggers protective action when too low
- Formula: (Equity / Margin) × 100

**Returns**: Margin level percentage or error

### 4. GetOpenPositions

```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all currently open positions.

**Usage in RiskManager**:

- Counts total positions vs MaxOpenPositions
- Counts per-symbol exposure
- Identifies worst losing position for margin protection

**Returns**: Slice of position info or error

### 5. CloseAllPositions

```go
func (s *MT5Sugar) CloseAllPositions() (int, error)
```

**Purpose**: Closes all open positions immediately.

**Usage in RiskManager**:

- Emergency action when critical limits breached
- Used for drawdown limits and daily loss limits
- Cannot be undone - final protective measure

**Returns**: Number of positions closed or error

### 6. ClosePosition

```go
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

**Purpose**: Closes a specific position by ticket number.

**Usage in RiskManager**:

- Closes worst losing position when margin low
- Frees up margin without closing all positions
- More surgical than CloseAllPositions

**Parameters**:
- `ticket`: Position ticket to close

**Returns**: Error if failed

---

## Risk Profiles

### Conservative Profile (Beginners)

```go
RiskManagerConfig{
    MaxDrawdownPercent:  5.0,              // Very tight
    MaxDrawdownAbsolute: 500.0,            // Small absolute limit
    DailyLossLimit:      200.0,            // Stop early on losses
    DailyProfitTarget:   300.0,            // Lock in profits early
    MinMarginLevel:      200.0,            // High margin safety
    MaxMarginUsed:       50.0,             // Conservative margin usage
    MaxOpenPositions:    10,               // Limited positions
    MaxSymbolExposure:   2,                // Minimal per-symbol
    MaxPositionSize:     0.5,              // Small sizes
    CheckInterval:       3 * time.Second,  // Frequent checks
    EnableAutoClose:     true,
    EnableTradeBlocking: true,
}
```

**Best for**:

- New traders learning risk management
- Small accounts ($1,000-$5,000)
- Demo account testing
- High-risk strategies needing tight control

### Moderate Profile (Default)

```go
RiskManagerConfig{
    MaxDrawdownPercent:  10.0,             // Moderate tolerance
    MaxDrawdownAbsolute: 1000.0,           // Reasonable limit
    DailyLossLimit:      500.0,            // Balanced daily limit
    DailyProfitTarget:   1000.0,           // Good profit target
    MinMarginLevel:      150.0,            // Standard margin safety
    MaxMarginUsed:       80.0,             // Healthy margin usage
    MaxOpenPositions:    20,               // Moderate count
    MaxSymbolExposure:   5,                // Reasonable per-symbol
    MaxPositionSize:     1.0,              // Standard lot size
    CheckInterval:       5 * time.Second,  // Regular checks
    EnableAutoClose:     true,
    EnableTradeBlocking: true,
}
```

**Best for**:

- Experienced traders
- Medium accounts ($5,000-$20,000)
- Swing trading strategies
- Multi-pair trading

### Aggressive Profile (Experienced)

```go
RiskManagerConfig{
    MaxDrawdownPercent:  20.0,             // Wider tolerance
    MaxDrawdownAbsolute: 5000.0,           // Large absolute limit
    DailyLossLimit:      2000.0,           // Higher daily tolerance
    DailyProfitTarget:   5000.0,           // Ambitious target
    MinMarginLevel:      100.0,            // Minimal margin safety
    MaxMarginUsed:       90.0,             // Aggressive margin usage
    MaxOpenPositions:    50,               // Many positions
    MaxSymbolExposure:   10,               // High per-symbol
    MaxPositionSize:     5.0,              // Large sizes
    CheckInterval:       10 * time.Second, // Less frequent checks
    EnableAutoClose:     true,
    EnableTradeBlocking: false,            // No trade blocking
}
```

**Best for**:

- Professional traders only
- Large accounts ($20,000+)
- Grid/martingale strategies
- When you can monitor actively

**WARNING**: Only use aggressive profile if you fully understand the risks!

---

## Risk Event Types

### Event Types

1. **MAX_DRAWDOWN_PERCENT**: Drawdown percentage exceeded

2. **MAX_DRAWDOWN_ABSOLUTE**: Drawdown dollar amount exceeded

3. **DAILY_LOSS_LIMIT**: Daily loss limit hit

4. **DAILY_PROFIT_TARGET**: Daily profit target reached

5. **LOW_MARGIN_LEVEL**: Margin level too low

6. **MAX_POSITIONS**: Too many open positions

7. **SYMBOL_EXPOSURE**: Too many positions on one symbol

8. **EMERGENCY_CLOSE**: Positions closed by risk manager

### Severity Levels

- **INFO**: Informational (profit target reached)
- **WARNING**: Potential issue (position count high)
- **CRITICAL**: Immediate action taken (drawdown exceeded, emergency close)

### Event Logging

- Last 100 events kept in memory
- Accessible via `GetRiskEvents()`
- Each event contains:
  - Timestamp
  - EventType
  - Severity
  - Description
  - Current value vs limit
  - Action taken

---

## Best Practices

### Set Limits Based on Account Size

- **Drawdown**: 5-10% of account
- **Daily loss**: 2-5% of account
- **Position size**: 1-2% risk per trade

### Test on Demo First

- Run risk manager on demo for 1 week minimum
- Ensure limits trigger correctly
- Verify emergency close works as expected

### Monitor Risk Events

- Check `GetRiskEvents()` regularly
- Look for patterns (same limit hit repeatedly?)
- Adjust limits if needed

### Combine with Trading Strategy

- Risk manager complements your strategy, doesn't replace it
- Your strategy should have its own stop losses
- Risk manager is "last line of defense"

### Use Daily Profit Target

- Lock in good days by stopping at target
- Prevents "giving it all back" syndrome
- Set realistic targets (don't be too greedy)

---

## Configuration Location

**Runtime configuration is located in**: `examples/demos/main.go -> RunOrchestrator_RiskManager()`

This separation exists because:

1. **Code reusability**: Same risk engine runs with different limits

2. **Quick testing**: Change limits without modifying protection logic

3. **User examples**: Shows how to configure risk management

4. **Centralized entry point**: Single place for all orchestrator launches

