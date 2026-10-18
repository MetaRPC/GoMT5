# RiskManager - Account Protection System

## Description

**🧩 RiskManager** is an automated account guardian that monitors risk metrics 24/7 and automatically intervenes when danger limits are exceeded. It **PREVENTS ACCOUNT BLOW-UPS** by enforcing strict risk rules: drawdown limits, daily loss caps, margin monitoring, and position limits. The "safety net" that saves you from catastrophic losses.

**Strategy Principle**: Continuously monitors critical account parameters and takes protective action when thresholds are breached:

- **Drawdown Protection**: Monitors both percentage and absolute drawdown from the peak equity; after an emergency close the peak is reset to the current equity
- **Daily Loss Limit**: Enforces maximum daily loss caps (prevents revenge trading)
- **Daily Profit Target**: Blocks trading for the rest of the day when the daily profit goal is reached
- **Margin Safety**: Closes the most losing position while the margin level is below the minimum, and blocks trading while the margin usage is above the maximum
- **Position Limits**: Logs warnings when position counts or per-symbol counts exceed the limits, and reduces a position larger than `MaxPositionSize`
- **Emergency Actions**: Automatically closes positions (`EnableAutoClose`) and blocks trading (`EnableTradeBlocking`) when limits hit

> **How "blocking" works.** The risk manager cannot intercept orders sent by other programs. While trading is blocked (`IsTradingBlocked()` is true), it closes every position that was not open when the block started - so a trade opened anyway is undone at the next check. Your own trading code should still check `IsTradingBlocked()` before opening positions.

> **One event per breach.** Every limit is checked on every check, but a breach logs its event (and takes its one-time action) once, when the limit is first exceeded; it can fire again only after the value went back inside the limit. A limit of **0 switches it off**.

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
GetEquity  GetMargin  GetMarginLevel  GetOpenPositions
(drawdown, (margin    (safety)        (limits, blocking)
 daily)     usage)
    |
CloseAllPositions / ClosePosition / ClosePositionPartial
(emergency intervention, blocked trades, oversized positions)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetBalance`, `GetEquity`, `GetMargin`, `GetMarginLevel`, `GetOpenPositions`, `GetSymbolInfo`, `CloseAllPositions`, `ClosePosition`, `ClosePositionPartial`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo` for position data structures

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `MaxDrawdownPercent` | float64 | `10.0` | Maximum drawdown from the peak equity, in % (0 = off) |
| `MaxDrawdownAbsolute` | float64 | `1000.0` | Maximum drawdown from the peak equity, in account currency (0 = off) |
| `DailyLossLimit` | float64 | `500.0` | Maximum loss today (equity - equity at the start of the day; 0 = off) |
| `DailyProfitTarget` | float64 | `1000.0` | Trading is blocked for the rest of the day after this profit (0 = off) |
| `MinMarginLevel` | float64 | `150.0` | Minimum margin level % (equity / margin × 100; 0 = off) |
| `MaxMarginUsed` | float64 | `80.0` | Maximum margin usage % (margin / equity × 100): above it trading is blocked until it is back below (0 = off) |
| `MaxOpenPositions` | int | `20` | Maximum concurrent positions - a WARNING event (0 = off) |
| `MaxSymbolExposure` | int | `5` | Maximum positions per symbol - a WARNING event (0 = off) |
| `MaxPositionSize` | float64 | `1.0` | Maximum lots per position: a WARNING event and, with `EnableAutoClose`, the excess volume is closed (0 = off) |
| `CheckInterval` | time.Duration | `5s` | How often to check risk metrics (0 = 5 s) |
| `EnableAutoClose` | bool | `true` | Automatically close positions on breach |
| `EnableTradeBlocking` | bool | `true` | Block trading on the daily limits and the margin usage limit (new positions are closed while blocked) |

Defaults are those of `orchestrators.DefaultRiskManagerConfig()`. `Start()` validates the configuration: a `CheckInterval` of 0 becomes 5 s, a negative limit returns an `invalid configuration` error.

Besides `Start`/`Stop`/`GetStatus`/`GetMetrics`/`IsRunning`, the `*RiskManager` has `GetRiskEvents() []RiskEvent` (a copy), `IsTradingBlocked() bool`, `GetTodayProfit() float64`, `GetPeakBalance() float64` (the peak **equity** the drawdown is measured from) and `GetDailyStartBalance() float64` (the equity at the start of the day). They are safe to call while the manager runs.

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

Both commands will launch the **RiskManager** orchestrator with the configuration specified in `main.go -> RunOrchestrator_RiskManager()`. The runner monitors for 15 minutes (900 s), showing the largest drawdown so far, today's profit, the number of risk events and the blocked flag, then calls `Stop()` and prints the final metrics; the runner then disconnects the terminal - it does so however the run ends, also on an error or Ctrl+C. Each orchestrator also cancels its own context on Ctrl+C / SIGTERM, so a program of yours that waits on `GetContext()` can call `Stop()` and disconnect.

> **Warning:** the runner uses `EnableAutoClose: true` - a drawdown or daily-loss breach closes ALL open positions of the account, and `EnableTradeBlocking: true` closes positions opened while trading is blocked.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Start(): validate config (CheckInterval 0 -> 5s, negative limits -> error)
  |
Initialize baseline values:
  - Starting balance
  - Peak = current equity
  - Daily start = current equity
  - Trading blocked = false, no active breaches
  |
MONITORING LOOP (every CheckInterval):
  |
  Check if new day (local date) -> daily start = equity, daily block lifted
  |
  Get current account state:
    - Equity (error -> count error, skip this check)
    - Margin, Margin Level (error -> margin checks skipped)
  |
  Update peak if equity higher
  |
  Calculate metrics:
    - Current drawdown = peak - equity
    - Drawdown % = (drawdown / peak) × 100   (metrics: MaxDrawdown = largest %)
    - Today profit = equity - daily start
  |
  CHECK RISK LIMITS (a limit of 0 is off; each breach logs ONE event):
    |
    1. DRAWDOWN LIMITS:
       | Newly drawdown % >= MaxDrawdownPercent and/or drawdown $ >= MaxDrawdownAbsolute:
       |   -> Log CRITICAL event(s)
       |   -> If EnableAutoClose: CLOSE ALL POSITIONS once, then peak = current equity
    |
    2. DAILY LIMITS (once per day each):
       | Today loss >= DailyLossLimit:
       |   -> Log CRITICAL event
       |   -> If EnableAutoClose: CLOSE ALL POSITIONS
       |   -> If EnableTradeBlocking: block trading until tomorrow
       |
       | Today profit >= DailyProfitTarget:
       |   -> Log INFO event
       |   -> If EnableTradeBlocking: block trading until tomorrow
    |
    3. MARGIN LIMITS:
       | 0 < margin level < MinMarginLevel:
       |   -> If EnableAutoClose: log CRITICAL event + close most losing position (every check while low)
       |   -> else: log CRITICAL event once
       | margin / equity × 100 > MaxMarginUsed:
       |   -> Log WARNING event once
       |   -> If EnableTradeBlocking: trading blocked while above the limit
    |
    4. POSITION LIMITS (GetOpenPositions, after the actions above):
       | open positions > MaxOpenPositions -> WARNING event
       | positions of a symbol > MaxSymbolExposure -> WARNING event
       | position volume > MaxPositionSize -> WARNING event;
       |   if EnableAutoClose: close the excess (ClosePositionPartial)
    |
    5. TRADE BLOCK (EnableTradeBlocking):
       | First check under the block: remember the open positions
       | Later checks: close every other position (TRADE_BLOCKED event)
  |
  Update status metrics
  |
REPEAT until stopped (Stop() or Ctrl+C)
  |
ON STOP:
  | Cancel monitoring loop and wait for a check in progress
  | (main.go then prints the final metrics)
  |
END
```

### 🛠️ Step-by-step Description

#### 1. Initialization
 
> Lines `219–243`  


```go
func (r *RiskManager) initialize() error {
	// Get starting balance and equity
	balance, err := r.sugar.GetBalance()
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}
	equity, err := r.sugar.GetEquity()
	if err != nil {
		return fmt.Errorf("failed to get equity: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.startingBalance = balance
	r.peakBalance = equity
	r.dailyStartBalance = equity
	// ... trading unblocked, no active breaches ...
	return nil
}
```

**What happens**:

- Gets the current balance (kept as the starting balance) and equity
- The peak and the daily start are the current **equity** (drawdown and today's P/L both include floating P/L)
- Initializes trading as unblocked and clears the breach memory
- `Start()` fails if the balance or equity cannot be read

#### 2. Monitoring Loop
 
> Lines `246–258`  


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
4. **Graceful exit**: Stops cleanly when context is cancelled (`Stop()` or Ctrl+C); `Stop()` waits for a check in progress

#### 3. Daily Reset Check
  
> Lines `652–675`  

`checkDailyReset()` compares the local date (year, month, day of the machine running the program) with the date of the last reset. On a new day it sets the daily start to the current equity, lifts the daily block and allows the daily loss / profit target events again. Today's profit starts again from zero.

#### 4. One Event per Breach
 
> Lines `322–332`  


```go
func (r *RiskManager) breach(key string, exceeded bool) bool {
	if !exceeded {
		delete(r.activeBreaches, key)
		return false
	}
	if r.activeBreaches[key] {
		return false
	}
	r.activeBreaches[key] = true
	return true
}
```

`breach()` is true only on the first check that finds a limit exceeded; while it stays exceeded nothing is logged again, and once the value is back inside the limit the breach is cleared.

#### 5. Drawdown Protection
 
> Lines `337–378`  


```go
func (r *RiskManager) checkDrawdownLimit(drawdownPercent, drawdownAbsolute float64) {
	// ... action text ...
	var reasons []string

	// Check percentage drawdown
	if r.breach("MAX_DRAWDOWN_PERCENT", r.config.MaxDrawdownPercent > 0 && drawdownPercent >= r.config.MaxDrawdownPercent) {
		r.logRiskEvent("MAX_DRAWDOWN_PERCENT", "CRITICAL", /* ... */)
		reasons = append(reasons, "Maximum drawdown exceeded")
	}

	// Check absolute drawdown
	if r.breach("MAX_DRAWDOWN_ABSOLUTE", r.config.MaxDrawdownAbsolute > 0 && drawdownAbsolute >= r.config.MaxDrawdownAbsolute) {
		r.logRiskEvent("MAX_DRAWDOWN_ABSOLUTE", "CRITICAL", /* ... */)
		reasons = append(reasons, "Maximum absolute drawdown exceeded")
	}

	if len(reasons) > 0 && r.config.EnableAutoClose {
		r.closeAllPositionsEmergency(strings.Join(reasons, ", "))
		r.resetPeak()
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
- **A limit of 0 is off** (with `>=` a 0 limit would otherwise be "breached" by a drawdown of 0 on every check)
- When both are hit in one check, both events are logged and `CloseAllPositions` runs once
- After the emergency close `resetPeak()` makes the current equity the new peak, so the next breach is measured from there instead of repeating on every check
- Without `EnableAutoClose` the breach is only logged (once); it can fire again after the drawdown went back below the limit
- A drawdown breach does **not** block trading

#### 6. Daily Limits Protection
 
> Lines `381–420`  

- Today's profit is `equity - dailyStartBalance` (the daily start is the equity at the start of the day), so floating P/L counts
- **Daily loss limit** (once per day): logs a CRITICAL event, closes all positions (if `EnableAutoClose`) and blocks trading until the next day (if `EnableTradeBlocking`)
- **Daily profit target** (once per day): logs an INFO event and blocks trading until the next day (if `EnableTradeBlocking`); existing positions stay open

#### 7. Margin Safety
  
> Lines `423–453`  

- **Margin Level** = (Equity / Margin) × 100, read with `GetMarginLevel()`; 0 (no margin used) is ignored
- While the level is below `MinMarginLevel` (strictly `<`) and `EnableAutoClose` is on, the position with the most negative `Profit` is closed on every check (each close is logged) until the level is safe again; without auto-close the breach is logged once
- **Margin usage** = Margin / Equity × 100. Above `MaxMarginUsed` a WARNING event is logged once and, with `EnableTradeBlocking`, trading is blocked until the usage is back below the limit
- **Example**: MinMarginLevel = 150%
  - Equity = $10,000, Margin = $6,667
  - Margin Level = 150% (borderline, not below); margin usage 66.7% (below 80%)
  - If it drops to 149%: close worst position

#### 8. Position Limits
  
> Lines `456–538`  

- `MaxOpenPositions` and `MaxSymbolExposure`: WARNING events only (one per breach)
- `MaxPositionSize`: a WARNING event per oversized position; with `EnableAutoClose` the excess is closed - rounded UP to the symbol's volume step (from `GetSymbolInfo`), so the rest is at most the limit (the whole position when the rest would be below the minimum volume)

#### 9. Trade Blocking
  
> Lines `550–583`  

```go
func (r *RiskManager) enforceTradeBlock(positions []*pb.PositionInfo) {
	if !r.IsTradingBlocked() || !r.config.EnableTradeBlocking {
		r.allowedTickets = nil
		return
	}

	// First check under the block: remember what is open now
	if r.allowedTickets == nil {
		// ... every open ticket is allowed ...
		return
	}

	for _, pos := range positions {
		if r.allowedTickets[pos.Ticket] {
			continue
		}
		// ... ClosePosition + TRADE_BLOCKED event ...
	}
}
```

Trading is blocked by the daily loss limit or the daily profit target (until the next day) and by `MaxMarginUsed` (while above it). Positions that were open when the block started are left alone; every position opened later is closed at the next check.

#### 10. Emergency Close All
  
> Lines `601–619`  

- Calls `sugar.CloseAllPositions()`: closes every open position of the account (pending orders are not touched)
- `CloseAllPositions` keeps going when a close fails and returns an error naming the failed tickets; the risk manager counts it as an error and still logs how many positions were closed
- Cannot be undone
- Triggered by (only when `EnableAutoClose` is true):
  - Maximum drawdown exceeded (% or $)
  - Daily loss limit exceeded
- Logs an `EMERGENCY_CLOSE` event and updates the metrics

---

## Visual Examples

### Scenario 1: Drawdown Protection in Action

**Initial State**:

- Account equity: $10,000
- MaxDrawdownPercent: 10%
- MaxDrawdownAbsolute: $1,000

```
TIME 0: Trading Starts

Equity:       $10,000
Peak:         $10,000
Drawdown:     $0 (0%)
Positions:    0
Status:       Monitoring

TIME 1: Good Trading Day

Equity:       $11,500
Peak:         $11,500 (new peak!)
Drawdown:     $0 (0%)
Positions:    3
Status:       Monitoring - all good

TIME 2: Market Turns Against Positions

Equity:       $10,800
Peak:         $11,500 (unchanged - peak is max)
Drawdown:     $700 (6.1%)
Positions:    5
Status:       Monitoring - within limits

TIME 3: Drawdown Approaching Limit

Equity:       $10,550
Peak:         $11,500
Drawdown:     $950 (8.3%)
Status:       Within limits (no event; there is no "approaching" warning)

TIME 4: DRAWDOWN LIMIT EXCEEDED

Equity:       $10,500
Peak:         $11,500
Drawdown:     $1,000 (8.7%)
TRIGGER:      MaxDrawdownAbsolute HIT ($1,000 >= $1,000)
              (the 10% limit would only be reached at $1,150)
ACTION:       EMERGENCY CLOSE ALL 5 POSITIONS (CloseAllPositions)
Result:       Positions closed, equity ~$10,500
Peak reset:   $10,500 - the next breach is measured from here
              (another $1,000 loss from $10,500)

WITHOUT RISK MANAGER:
- Positions could continue losing
- Drawdown could reach 20-30% or more
- Account could blow up completely

WITH RISK MANAGER:
- Loss limited to about $1,000 from peak (whichever limit hits first)
- $1,000 loss vs potential $3,000+ loss
- Account survives to trade another day
```

### Scenario 2: Daily Loss Limit

**Configuration**:

- DailyLossLimit: $500
- EnableAutoClose: true
- EnableTradeBlocking: true

```
TIME 0: Start of Day (Midnight Reset)

Daily Start Equity:  $10,000
Current Equity:      $10,000
Today Profit:        $0
Trading Blocked:     false

TIME 1: Morning Trades (09:00)

Current Equity:      $9,800
Today Profit:        -$200
Status:              Within limits

TIME 2: Bad Trading Session (11:00)

Current Equity:      $9,600
Today Profit:        -$400
Status:              Within limits (no "approaching" warning)

TIME 3: DAILY LIMIT HIT (12:30)

Current Equity:      $9,500
Today Profit:        -$500
TRIGGER:             DailyLossLimit HIT!
ACTION:              1. CLOSE ALL POSITIONS
                     2. BLOCK TRADING until tomorrow
Status:              IsTradingBlocked() == true for rest of day

TIME 4: Attempt to Trade (14:00)

Status:              IsTradingBlocked() == true
Result:              A strategy that checks IsTradingBlocked() skips the
                     trade. A position opened anyway is closed by the risk
                     manager at its next check (TRADE_BLOCKED event)

TIME 5: Next Day (Midnight Reset)

Daily Start Equity:  $9,500 (updated to current)
Today Profit:        $0 (reset)
Trading Blocked:     false (unblocked)
Status:              New day, trading allowed again

Note: today's profit is equity-based (floating P/L included). The
closing trades themselves can push the day's loss beyond $500.

BENEFIT:
- Prevented emotional "revenge trading"
- Limited the daily loss to about $500
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

Equity:       $7,450
Margin Used:  $5,000
Margin Level: 149%
TRIGGER:      MarginLevel < MinMarginLevel
ACTION:       Close most losing position

Position Details:
- Position 1: -$200
- Position 2: +$50
- Position 3: -$500 <- WORST LOSS
- Position 4: +$100
- Position 5: -$150

ACTION TAKEN: Close Position #3 (worst loser)

TIME 4: After Closing Worst Position

Equity:       $7,450
Margin Used:  $4,000 (reduced)
Margin Level: 186% (improved!)
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

- Called every CheckInterval (3 s timeout)
- Used to calculate the drawdown from the peak and today's P/L (equity - equity at the start of the day)
- Read at start (peak, daily start), after an emergency close (new peak) and at the daily reset
- Core metric; if it fails, the whole check is skipped

**Returns**: Current equity or error

### 2. GetBalance

```go
func (s *MT5Sugar) GetBalance() (float64, error)
```

**Purpose**: Gets current account balance (closed positions only).

**Usage in RiskManager**:

- Read once in `initialize()` (the starting balance); `Start()` fails if it cannot be read

**Returns**: Current balance or error

### 3. GetMarginLevel / GetMargin

```go
func (s *MT5Sugar) GetMarginLevel() (float64, error)
func (s *MT5Sugar) GetMargin() (float64, error)
```

**Purpose**: Gets the current margin level percentage / the margin in use.

**Usage in RiskManager**:

- Margin level: risk of a margin call; below `MinMarginLevel` the most losing position is closed (a value of 0 is ignored)
- Margin / equity × 100 = margin usage; above `MaxMarginUsed` trading is blocked
- If either cannot be read, the margin checks of that check are skipped (counted as an error)

**Returns**: Margin level percentage / margin, or error

### 4. GetOpenPositions

```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all currently open positions.

**Usage in RiskManager**:

- Counts total positions vs MaxOpenPositions
- Counts per-symbol exposure, checks each position's volume against MaxPositionSize
- Finds positions opened while trading is blocked
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
- Closes positions only (pending orders are not touched); it keeps going when a close fails

**Returns**: Number of positions closed, and an error when listing failed or a close failed (the error names the failed tickets)

### 6. ClosePosition

```go
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

**Purpose**: Closes a specific position by ticket number.

**Usage in RiskManager**:

- Closes worst losing position when margin low
- Closes positions opened while trading is blocked
- Frees up margin without closing all positions
- More surgical than CloseAllPositions

**Parameters**:
- `ticket`: Position ticket to close

**Returns**: Error if failed

### 7. ClosePositionPartial and GetSymbolInfo

```go
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error)
```

**Usage in RiskManager**: with `EnableAutoClose`, a position larger than `MaxPositionSize` is reduced with `ClosePositionPartial`; the excess is rounded up to the symbol's `VolumeStep` (read once per symbol with `GetSymbolInfo`).

---

## Risk Profiles

### Conservative Profile (Beginners)

```go
config := orchestrators.RiskManagerConfig{
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
config := orchestrators.RiskManagerConfig{
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
config := orchestrators.RiskManagerConfig{
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

9. **MARGIN_USAGE**: Margin usage above `MaxMarginUsed`

10. **POSITION_SIZE**: A position larger than `MaxPositionSize`

11. **TRADE_BLOCKED**: A position opened while trading is blocked (closed again)

Risk events are not counted as errors in the orchestrator status: the status counts errors (failed API calls) and one success per completed check. Each event is also printed to the console when it is logged.

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
  - Value and Limit
  - ActionTaken (what the manager did, e.g. `"Closed all positions"`, `"trading blocked until tomorrow"`, `"Closed 0.50 lots"`, or `"Logged"`)
- One event per breach: an exceeded limit is logged when it is first exceeded, not on every check

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

