# PositionScaler - Advanced Position Sizing

## Description

**🧩 PositionScaler** is an advanced position management orchestrator that implements three sophisticated scaling strategies: **PYRAMIDING** (add to winners), **AVERAGING DOWN** (add to losers), and **SCALE OUT** (gradual profit taking). It maximizes profits in trending markets or manages recovery from drawdowns using intelligent scaling algorithms.

**Strategy Modes**:

1. **PYRAMIDING MODE** (Recommended for trending markets):
   - Adds to WINNING positions as they move favorably
   - Increases exposure when trend is confirmed by price action
   - Uses smaller lot sizes for each additional scale-in (risk management)
   - Each scale-in position gets its own stop loss (`StopLossPerScale` points)
   - Best for strong trending markets (EURUSD trending)

2. **AVERAGING DOWN MODE** (High risk - use with caution):
   - Adds to LOSING positions at better prices
   - Reduces average entry price to recover faster
   - VERY RISKY - requires strict maximum loss limits
   - Only suitable in strong support/resistance areas
   - NOT recommended for beginners

3. **SCALE OUT MODE** (Profit protection):
   - Partially closes the largest position at profit levels `TriggerDistance` apart (first level: `TriggerDistance` from the base entry, each next one `TriggerDistance` beyond the previous scale-out)
   - Takes partial profits while letting remainder run
   - Reduces risk exposure over time
   - Secures gains while maintaining market exposure
   - Best for volatile markets or profit protection

> **Base position.** In Pyramiding and AveragingDown mode the scaler opens a **BUY** base position of `InitialLotSize` for each symbol listed in `Symbols` that has no open position - once per run (the example has no entry signal, so the base is always a BUY). With `InitialLotSize: 0`, in ScaleOut mode, or with an empty `Symbols` list (= all symbols) it opens nothing and only works on positions that are already open (opened by you or another program). Positions are grouped per symbol: the oldest position (lowest ticket) of a symbol is the group's base (`BaseTicket`, `BasePrice`, direction), the group's size counts the symbol's positions in that direction, and all profit/loss checks are measured from the base entry.

**File**: `examples/demos/orchestrators/12_position_scaler.go`

> **IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY**
>
> This orchestrator demonstrates how GoMT5 methods function and combine into automated strategies. **NOT a production-ready trading strategy!** Averaging down is especially dangerous if price continues falling. Always test on demo accounts first.

---

## ⚙️ Architecture

```
POSITION SCALER ORCHESTRATOR
    |
MT5Sugar Instance
    |
    |  |  |  |  |
GetOpenPositions  GetSymbolInfo  GetPriceInfo  BuyMarket / Buy/SellMarketWithSLTP     GetPositionsBySymbol + ClosePosition(Partial)
(monitoring)      (point, volume) (tracking)   (base position / scaling in)          (scaling out)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetOpenPositions`, `GetSymbolInfo`, `GetPriceInfo`, `BuyMarket`, `BuyMarketWithSLTP`, `SellMarketWithSLTP`, `GetPositionsBySymbol`, `ClosePosition`, `ClosePositionPartial`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo` for position data structures

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Mode` | ScalingMode | `Pyramiding` | Scaling strategy (Pyramiding/AveragingDown/ScaleOut) |
| `TriggerDistance` | float64 | `200` | Points the price must move for the next scale: from the base entry for the first, from the last scale price for the next ones (must be > 0) |
| `MinProfitToScale` | float64 | `100` | Min profit before scaling (pyramiding mode) |
| `MaxLossToAverage` | float64 | `500` | Max loss before stopping averaging (critical safety; must be > 0 in AveragingDown mode) |
| `InitialLotSize` | float64 | `0.10` | Size of the BUY base position opened for a listed symbol without a position (once per run; 0 = only manage existing positions; not used in ScaleOut mode) |
| `ScaleLotSize` | float64 | `0.05` | Size for each scale-in operation (must be > 0 for the scale-in modes) |
| `MaxScales` | int | `3` | Maximum number of scale operations (scale-ins, or partial closes in ScaleOut mode; must be > 0) |
| `ReducePerScale` | float64 | `0.0` | If > 0: next size = `ScaleLotSize × (1 − ReducePerScale × ScaleCount)` (0 to < 1) |
| `TotalMaxLotSize` | float64 | `1.0` | Maximum total size of the symbol's positions in the base direction (must be > 0 for the scale-in modes) |
| `ScaleOutPercent` | float64 | `0.33` | Fraction of the largest position to close at each scale-out (0 to 1) |
| `StopLossPerScale` | float64 | `150` | SL distance in points for each scale-in (at least the symbol's stops level; 0 = no SL) |
| `Symbols` | []string | `[]string{symbol}` | Symbols to manage (empty = all symbols) |
| `CheckInterval` | time.Duration | `5s` | How often to check for scaling opportunities (0 = 5 s) |

Defaults are those of `orchestrators.DefaultPositionScalerConfig(symbol string)`. `Start()` validates the configuration: a `CheckInterval` of 0 becomes 5 s, values that cannot work (see the table) return an `invalid configuration` error instead of starting.

> **Point size and volumes:** all distances are in points of the symbol. The scaler reads the point size, digits, stops level and volume limits once per symbol with `sugar.GetSymbolInfo(symbol)`, so 200 points are 0.00200 on EURUSD, 0.200 on USDJPY and 2.00 on XAUUSD. Every volume is rounded down to the symbol's volume step; a size below the minimum volume is not traded. SL prices are rounded to the symbol's digits.

### Configuration Example

```go
// Pyramiding Configuration
config := orchestrators.PositionScalerConfig{
    Mode:             orchestrators.Pyramiding,
    TriggerDistance:  200,           // Scale every 200 points
    MinProfitToScale: 100,           // Activate after +100 points profit
    InitialLotSize:   0.10,          // BUY base position if EURUSD has none
    ScaleLotSize:     0.05,          // Add 0.05 lots each scale
    MaxScales:        3,             // Maximum 3 additional positions
    TotalMaxLotSize:  1.0,           // Total max 1.0 lots
    StopLossPerScale: 150,           // 150pt SL for scale-ins
    Symbols:          []string{"EURUSD"},
    CheckInterval:    5 * time.Second,
}

scaler := orchestrators.NewPositionScaler(sugar, config)
```

---

## How to Run

You can execute this orchestrator using several command variations:

```bash
# Navigate to examples directory
cd examples/demos

# Option 1: By number
go run main.go 12

# Option 2: By keyword
go run main.go scaler
```

Both commands will launch the **PositionScaler** orchestrator with the configuration specified in `main.go -> RunOrchestrator_PositionScaler()` (Pyramiding mode on `cfg.TestSymbol`). The runner runs the scaler for 5 minutes (300 s), then calls `Stop()` and prints the final metrics; the runner then disconnects the terminal - it does so however the run ends, also on an error or Ctrl+C. Each orchestrator also cancels its own context on Ctrl+C / SIGTERM, so a program of yours that waits on `GetContext()` can call `Stop()` and disconnect. If the symbol has no position, the first check opens a 0.10-lot BUY base position (`InitialLotSize`); the scaler adds to it once the price has moved `TriggerDistance` in its favour. Positions stay open when the scaler stops.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Initialize PositionScaler
  |
Create empty position group tracker
  |
Start(): validate config (CheckInterval 0 -> 5s; invalid values -> error)
  |
MONITORING LOOP (every CheckInterval):
  | Get all open positions
  |
  BASE POSITION (Pyramiding / AveragingDown, InitialLotSize > 0):
    | Listed symbol without a position and no base opened for it in this run?
    |   -> BuyMarket(symbol, InitialLotSize rounded to the volume step); end this check
  |
  UPDATE GROUPS (one PositionGroup per symbol):
    | Skip positions whose symbol is not tracked (if Symbols list specified)
    | Oldest position (lowest ticket) of a symbol = base -> group (base ticket, base price, direction)
    | A new base ticket starts a new group (ScaleCount 0)
    | Total lot size = the symbol's positions in the base direction
    | Remove groups whose symbol has no open position any more
  |
  FOR each group: Should we scale this group?
      |
      | Check 1: Have we reached MaxScales? -> YES: Skip
      | Check 2 (scale-in modes): would total exceed TotalMaxLotSize,
      |          or is the next size below the minimum volume? -> YES: Skip
      | Check 3: Get current price; profit in points vs BASE price
      |          (BUY: Bid - BasePrice, SELL: BasePrice - Ask), point from GetSymbolInfo
      |          and the favourable move since the last scale (no scale yet: = profit)
      |
      | MODE-SPECIFIC CHECKS:
      |
      | PYRAMIDING MODE:
      |   | Is profit >= MinProfitToScale? -> NO: Skip
      |   | Moved TriggerDistance in our favour since the last scale (or the base entry)? -> NO: Skip
      |   | YES -> SCALE IN (open additional position)
      |
      | AVERAGING DOWN MODE:
      |   | Is position in loss? -> NO: Skip
      |   | Is loss <= MaxLossToAverage? -> NO: Skip (safety)
      |   | Moved TriggerDistance against us since the last scale (or the base entry)? -> NO: Skip
      |   | YES -> SCALE IN (open additional position at better price)
      |
      | SCALE OUT MODE:
      |   | Moved TriggerDistance in our favour since the last scale-out (or the base entry)? -> NO: Skip
      |   | YES -> SCALE OUT (close ScaleOutPercent of the largest position)
      |
  FOR SCALE IN operation:
    | Calculate next scale lot size, round down to the volume step
    | Get current price
    | Calculate stop loss for scale-in (distance = max(StopLossPerScale, stops level);
    |   BUY: Bid - distance*point, SELL: Ask + distance*point; none if StopLossPerScale = 0), TP = 0
    | BuyMarketWithSLTP / SellMarketWithSLTP
    | LastScalePrice = Bid (BUY) / Ask (SELL) - the price the moves are measured with
    | Add ticket to PositionGroup.ScaleTickets
    | Increment PositionGroup.ScaleCount
    | Update PositionGroup.TotalLotSize
    | Update metrics
  |
  FOR SCALE OUT operation:
    | GetPositionsBySymbol, take the LARGEST position in the group's direction
    | Close volume = its volume × ScaleOutPercent, rounded down to the volume step
    |   (the whole position if the rest would be below the minimum volume)
    | ClosePositionPartial (or ClosePosition for the whole position)
    | LastScalePrice = current Bid (BUY) / Ask (SELL): the next level is TriggerDistance beyond
    | Update PositionGroup.ScaleCount
    | Update PositionGroup.TotalLotSize
    | Update metrics
  |
REPEAT until stopped (Stop() or Ctrl+C)
  |
ON STOP:
  | Cancel monitoring loop and wait for a check in progress
  | Clear all position groups (positions stay open)
  |
END
```

### Step-by-step Description

#### 1. Initialization
 
> Lines `232–241`  
> `examples/demos/orchestrators/12_position_scaler.go`


```go
func NewPositionScaler(sugar *mt5.MT5Sugar, config PositionScalerConfig) *PositionScaler {
    return &PositionScaler{
        BaseOrchestrator: NewBaseOrchestrator("Position Scaler"),
        sugar:            sugar,
        config:           config,
        trackedGroups:    make(map[string]*PositionGroup),
        symbolSpecs:      make(map[string]symbolSpec),
        baseOpened:       make(map[string]bool),
    }
}
```

**What happens**:

- Creates new scaler instance with provided configuration
- Initializes empty position group tracker (maps symbol to PositionGroup)
- Each symbol gets one PositionGroup tracking all related positions
- `symbolSpecs` caches point size, digits, stops level and volume limits per symbol, read once with `sugar.GetSymbolInfo(symbol)`
- `baseOpened` remembers for which symbols this run already tried to open the base position

**PositionGroup structure**:
```go
type PositionGroup struct {
    Symbol         string    // Symbol being traded
    IsBuy          bool      // Direction (true=BUY, false=SELL)
    BaseTicket     uint64    // Original position ticket (the oldest position of the symbol)
    ScaleTickets   []uint64  // Tickets of scale-in positions (or of the positions scaled out)
    BasePrice      float64   // Entry price of base position
    TotalLotSize   float64   // Total size of the symbol's positions in the base direction
    ScaleCount     int       // Number of scales executed
    LastScalePrice float64   // Price of last scale (0 = none yet)
    LastScaleTime  time.Time // When last scaled
}
```

#### 2. Starting the Scaler

> Lines `244–262`

```go
func (p *PositionScaler) Start() error {
    if p.IsRunning() {
        return fmt.Errorf("position scaler already running")
    }

    // Validate configuration (applies the default CheckInterval)
    if err := p.config.validate(); err != nil {
        return fmt.Errorf("invalid configuration: %w", err)
    }

    // Create context, mark as started, stop on Ctrl+C
    p.begin()

    // Start monitoring loop
    p.goLoop(p.monitorLoop)

    return nil
}
```

**What happens**:

- Checks if already running (prevents duplicate starts)
- Validates the configuration (`CheckInterval` 0 becomes 5 s - a zero interval would make `time.NewTicker` panic; values that cannot work are an error)
- `begin()` creates the cancellable context, marks the orchestrator as started and cancels the context on Ctrl+C / SIGTERM
- `goLoop()` launches the monitoring loop in a goroutine that `Stop()` waits for
- Returns immediately (non-blocking)

#### 3. Monitoring Loop

> Lines `285–297`

```go
func (p *PositionScaler) monitorLoop() {
    ticker := time.NewTicker(p.config.CheckInterval)
    defer ticker.Stop()

    for {
        select {
        case <-p.GetContext().Done():
            return
        case <-ticker.C:
            p.checkScalingOpportunities()
        }
    }
}
```

**Monitoring actions every CheckInterval**:

1. **Ticker triggers**: Every 5 seconds by default
2. **Context check**: Monitors for stop signal
3. **Check opportunities**: Calls `checkScalingOpportunities()`, which first opens a missing base position (once per symbol and run), then updates the groups and scales them
4. **Graceful exit**: Stops cleanly when context is cancelled (`Stop()` or Ctrl+C); `Stop()` waits for a check in progress

#### 4. Scaling Decision Logic


`shouldScale()` (lines 439-511) first returns `false` when `ScaleCount >= MaxScales`, or (scale-in modes) when the next size rounded to the volume step is 0 or `TotalLotSize + nextScaleSize > TotalMaxLotSize`. It then computes `profitPoints` from the group's base price (BUY: `(Bid - BasePrice) / point`, SELL: `(BasePrice - Ask) / point`, `point` from `GetSymbolInfo`) and the favourable move since the last scale:

```go
	// Favourable move since the last scale (no scale yet: since the base entry)
	moveSinceLast := profitPoints
	if group.ScaleCount > 0 {
		if group.IsBuy {
			moveSinceLast = (currentPrice - group.LastScalePrice) / point
		} else {
			moveSinceLast = (group.LastScalePrice - currentPrice) / point
		}
	}
```

and switches on the mode (lines 490-508):

```go
	switch p.config.Mode {
	case Pyramiding:
		// For pyramiding: position must be in profit and the price must have moved
		// TriggerDistance in our favour since the last scale (first scale: since the base entry)
		return profitPoints >= p.config.MinProfitToScale && moveSinceLast >= p.config.TriggerDistance

	case AveragingDown:
		// For averaging: position must be in loss but not too much
		if profitPoints > 0 || -profitPoints > p.config.MaxLossToAverage {
			return false
		}
		// ... and the price must have moved TriggerDistance against us since the last average
		return -moveSinceLast >= p.config.TriggerDistance

	case ScaleOut:
		// For scale-out: each level is TriggerDistance beyond the previous one
		// (first level: TriggerDistance from the base entry)
		return moveSinceLast >= p.config.TriggerDistance
	}
```

> Note: `MaxLossToAverage` is measured from the base entry only, and it only stops further averaging - it does not close anything. `LastScalePrice` is recorded with the same price the moves are checked with (Bid for BUY groups, Ask for SELL groups), so the spread does not add to the trigger distance. In ScaleOut mode `LastScalePrice` is the price of the last partial close, so the levels are `TriggerDistance` apart (e.g. +300 / +600 / +900).

---

## 💡 Visual Examples

### Scenario 1: PYRAMIDING MODE (Recommended)

**Configuration**:

- Entry: 1.10000 BUY EURUSD 0.10 lots
- TriggerDistance: 200 points
- MinProfitToScale: 100 points
- ScaleLotSize: 0.05 lots
- MaxScales: 3
- Prices in the example are Bid prices: the moves are checked on the Bid and the scale price is recorded as the Bid (a BUY fills at the Ask, one spread higher)

```
TIME 0: Initial Position Opened

Price:        1.10000 (Entry)
Position:     0.10 lots
Profit:       0 points
Status:       Waiting for activation (need +100pts)

TIME 1: Price Rises to 1.10100

Price:        1.10100
Position:     0.10 lots
Profit:       +100 points
Status:       Profit threshold met but waiting for TriggerDistance

TIME 2: Price at 1.10200 (FIRST SCALE)

Price:        1.10200
Profit:       +200 points
Action:       PYRAMIDING! Add 0.05 lots @ 1.10200
New Position: 0.15 lots total
Breakdown:    - Base: 0.10 lots @ 1.10000
              - Scale #1: 0.05 lots @ 1.10200
Status:       Watching for next trigger (+200pts from 1.10200)

TIME 3: Price at 1.10400 (SECOND SCALE)

Price:        1.10400
Profit:       +400 points (base), +200pts (scale #1)
Action:       PYRAMIDING! Add 0.05 lots @ 1.10400
New Position: 0.20 lots total
Breakdown:    - Base: 0.10 lots @ 1.10000 (+400pts)
              - Scale #1: 0.05 lots @ 1.10200 (+200pts)
              - Scale #2: 0.05 lots @ 1.10400 (just opened)
Status:       Watching for next trigger

TIME 4: Price at 1.10600 (THIRD SCALE)

Price:        1.10600
Profit:       +600 points (base), +400pts (scale #1), +200pts (scale #2)
Action:       PYRAMIDING! Add 0.05 lots @ 1.10600
New Position: 0.25 lots total
Breakdown:    - Base: 0.10 lots @ 1.10000 (+600pts)
              - Scale #1: 0.05 lots @ 1.10200 (+400pts)
              - Scale #2: 0.05 lots @ 1.10400 (+200pts)
              - Scale #3: 0.05 lots @ 1.10600 (just opened)
Status:       MAX SCALES REACHED (3/3) - No more pyramiding

TIME 5: Price at 1.11000 (Holding)

Price:        1.11000
Position:     0.25 lots
Total Profit: 0.10×1000 + 0.05×800 + 0.05×600 + 0.05×400 = $190
Without Pyramiding: 0.10×1000 = $100
Improvement:  90% more profit!

RESULT COMPARISON:

Without Pyramiding:
  - Single position: 0.10 lots
  - Profit: 1000 points × 0.10 lots = $100

With Pyramiding:
  - Total positions: 0.25 lots across 4 entries
  - Average entry: ~1.10240 (weighted)
  - Total profit: $190 (90% increase!)
  - Captured more of the trend's momentum
```

### Scenario 2: AVERAGING DOWN MODE (High Risk!)

**Configuration**:

- Entry: 1.10000 BUY EURUSD 0.10 lots
- TriggerDistance: 200 points
- MaxLossToAverage: 500 points (CRITICAL SAFETY)
- ScaleLotSize: 0.05 lots
- MaxScales: 2

```
TIME 0: Initial Position Opened

Price:        1.10000 (Entry - expecting rise)
Position:     0.10 lots
Profit:       0 points
Status:       Waiting for price movement

TIME 1: Price Falls to 1.09800 (First Trigger)

Price:        1.09800
Position:     0.10 lots
Profit:       -200 points (-$20)
Action:       AVERAGING DOWN! Add 0.05 lots @ 1.09800
New Position: 0.15 lots total
New Avg Entry: (1.10000×0.10 + 1.09800×0.05) / 0.15 = 1.09933
Breakdown:    - Base: 0.10 lots @ 1.10000 (-200pts)
              - Average #1: 0.05 lots @ 1.09800 (just opened)
Status:       Average entry improved from 1.10000 to 1.09933

TIME 2: Price Falls Further to 1.09600 (Second Trigger)

Price:        1.09600
Position:     0.15 lots
Loss:         -400 points total (still under MaxLossToAverage)
Action:       AVERAGING DOWN! Add 0.05 lots @ 1.09600
New Position: 0.20 lots total
New Avg Entry: (1.10000×0.10 + 1.09800×0.05 + 1.09600×0.05) / 0.20 = 1.09850
Breakdown:    - Base: 0.10 lots @ 1.10000 (-400pts)
              - Average #1: 0.05 lots @ 1.09800 (-200pts)
              - Average #2: 0.05 lots @ 1.09600 (just opened)
Status:       MAX SCALES REACHED (2/2) - No more averaging

TIME 3: Price Falls to 1.09400 (DANGER ZONE)

Price:        1.09400
Position:     0.20 lots
Loss:         -600 points base, -400 average #1, -200 average #2
Total Loss:   -$90
Status:       Would trigger another average but:
              - MaxScales reached (2/2)
              - Loss would exceed MaxLossToAverage (500pts)
              SAFETY MECHANISM PREVENTS FURTHER AVERAGING!

TIME 4: SCENARIO A - Recovery!

Price:        1.09850 (matches our average entry!)
Position:     0.20 lots
Profit:       $0 (break even)
Status:       WITHOUT averaging: would still be -$15 loss
              WITH averaging: break even achieved

TIME 5: Price Recovers to 1.10000 (Original Entry)

Price:        1.10000
Position:     0.20 lots
Profit:       (1.10000 - 1.09850) × 0.20 lots = 150 points = +$30
Status:       WITHOUT averaging: $0 profit
              WITH averaging: +$30 profit
              Recovery successful!

WARNING - What if price continued falling?

TIME X: Price Falls to 1.09000 (Disaster Scenario)

Price:        1.09000
Position:     0.20 lots (averaging stopped at 2 scales)
Total Loss:   -1000pts base, -800 average #1, -600 average #2 = -$170

Without averaging: 0.10 lots × -1000pts = -$100 loss
With averaging:    0.10×1000 + 0.05×800 + 0.05×600 = -$170 loss
Result:           70% BIGGER LOSS by averaging into losing position!

This is why MaxLossToAverage and MaxScales are CRITICAL safety parameters!
```

### Scenario 3: SCALE OUT MODE (Profit Protection)

**Configuration**:

- Entry: 1.10000 BUY EURUSD 0.40 lots (opened by you - ScaleOut mode opens no base position)
- TriggerDistance: 300 points
- ScaleOutPercent: 0.50 (50%)
- MaxScales: 3
- CheckInterval: 5s

```
TIME 0: Large Position Opened

Price:        1.10000 (Entry)
Position:     0.40 lots
Profit:       0 points
Status:       Waiting for the first level (+300pts from entry)

TIME 1: Check at 1.10300 (FIRST SCALE OUT)

Price:        1.10300
Move:         +300 points from the base entry
Action:       SCALE OUT! Close 50% of the largest position
Close:        0.40 × 0.50 = 0.20 lots CLOSED
Locked Profit: 0.20 lots × 300pts = +$60 SECURED
Remaining:    0.20 lots still open
Next level:   1.10300 + 300pts = 1.10600

TIME 2: Next check, 5 s later, price 1.10310

Price:        1.10310
Move:         +10 points since the last scale-out - below TriggerDistance
Action:       none

TIME 3: Price reaches 1.10600 (SECOND SCALE OUT)

Price:        1.10600
Close:        0.20 × 0.50 = 0.10 lots CLOSED
Locked Profit: 0.10 lots × 600pts = +$60 SECURED
Remaining:    0.10 lots still open
Next level:   1.10900

TIME 4: Price reaches 1.10900 (THIRD SCALE OUT)

Price:        1.10900
Close:        0.10 × 0.50 = 0.05 lots CLOSED
Locked Profit: 0.05 lots × 900pts = +$45 SECURED
Remaining:    0.05 lots still open
Total Locked: $60 + $60 + $45 = $165 secured
Status:       MAX SCALES REACHED (3/3) - no more scale-outs

RESULT:
  - Three partial closes at +300 / +600 / +900 points
  - 0.05 lots stay open for further gains
  - Close volumes are rounded down to the volume step: with
    ScaleOutPercent 0.33 and 0.30 lots the scaler closes 0.09 lots
    (not 0.099); when the rest would be below the minimum volume
    the whole position is closed
```

---

## Used MT5Sugar Methods

### 0. GetSymbolInfo and BuyMarket

```go
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error)
func (s *MT5Sugar) BuyMarket(symbol string, volume float64) (uint64, error)
```

- `GetSymbolInfo` is read once per symbol: `Point` converts the point distances into prices, `Digits` rounds SL prices, `StopLevel` is the minimum SL distance, `VolumeMin`/`VolumeStep`/`VolumeMax` round every volume
- `BuyMarket` opens the BUY base position of `InitialLotSize` for a listed symbol without a position (once per run)

### 1. GetOpenPositions

```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all currently open positions across all symbols.

**Usage in PositionScaler**:

- Called every CheckInterval (default 5 seconds)
- Returns full list of positions with: ticket, symbol, type, open price, volume, profit
- Used to update PositionGroups and identify scaling opportunities

**Returns**: Slice of position info or error

### 2. GetPriceInfo

```go
func (s *MT5Sugar) GetPriceInfo(symbol string) (*PriceInfo, error)
```

**Purpose**: Gets current Bid/Ask prices for a symbol.

**Usage in PositionScaler**:

- Called when evaluating scaling conditions
- BUY positions use Bid price (exit price) for profit calculation
- SELL positions use Ask price (exit price) for profit calculation
- Used to determine if TriggerDistance threshold is met

**Returns**: `mt5.PriceInfo{Symbol, Bid, Ask, SpreadPips, Time}` (3 s timeout)

### 3. BuyMarketWithSLTP / SellMarketWithSLTP

```go
func (s *MT5Sugar) BuyMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error)
func (s *MT5Sugar) SellMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error)
```

**Purpose**: Opens market order with stop loss and take profit.

**Usage in PositionScaler**:

- Called when scaling into position (Pyramiding or Averaging Down)
- Opens additional position with calculated lot size
- Sets stop loss based on StopLossPerScale parameter (BUY: Bid - distance*point, SELL: Ask + distance*point, distance = max(StopLossPerScale, stops level); no SL when StopLossPerScale is 0)
- TP is always 0 for scale-ins
- MT5Sugar sends volume and prices as given, so the scaler rounds them first (volume down to the volume step, SL to the symbol's digits)

**Parameters**:

- `symbol`: Trading symbol
- `volume`: Position size (ScaleLotSize, reduced by ReducePerScale if set)
- `sl`: Stop loss price
- `tp`: Take profit price (0 here)

**Returns**: The order ticket (`result.Order`, used by MT5 as the position ticket for a market order) or an error ("order rejected, code: ..." unless the return code is 10009)

### 4. GetPositionsBySymbol

```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*pb.PositionInfo, error)
```

**Purpose**: Gets all positions for a specific symbol.

**Usage in PositionScaler**:

- Used in Scale Out mode to find positions to close
- The scaler picks the LARGEST position in the group's direction
- Returns the open positions filtered by symbol (nil slice when none)

**Returns**: Slice of positions for the symbol

### 5. ClosePositionPartial / ClosePosition

```go
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

**Purpose**: Closes a partial amount of a position / the whole position.

**Usage in PositionScaler**:

- Core methods for Scale Out mode
- Closes `ScaleOutPercent` of the largest position's volume, rounded down to the symbol's volume step (0.30 × 0.33 = 0.099 becomes 0.09)
- `ClosePosition` is used instead when the rest would be below the minimum volume (or the rounded volume is 0)
- Allows gradual profit taking while maintaining exposure

**Parameters**:

- `ticket`: Position ticket to partially close
- `volume`: Amount in lots to close

---

## Risk Management

### Protection Benefits

**PositionScaler provides**:

1. **Pyramiding Mode**:

   - Maximizes profit in trending markets
   - Adds to winners only (reduced risk)
   - Smaller lot sizes for scales (risk control)
   - Best risk/reward when trend is confirmed

2. **Scale Out Mode**:

   - Locks in partial profits progressively
   - Reduces position risk over time
   - Psychological benefit of secured profits
   - Maintains some exposure for further gains

3. **Position Size Control**:

   - TotalMaxLotSize prevents over-exposure
   - MaxScales limits number of additions
   - ReducePerScale can decrease size per level

### Risk Considerations

1. **Averaging Down Dangers**:

   - Can DOUBLE or TRIPLE losses if trend continues
   - **Critical**: Always set MaxLossToAverage
   - **Critical**: Limit MaxScales to 2-3 maximum
   - Only use at strong support/resistance levels
   - NOT recommended for beginners

2. **Pyramiding Risks**:

   - Too aggressive scaling can create top-heavy position
   - Reversal can eliminate profits quickly
   - **Solution**: Use trailing stops on entire group
   - **Solution**: Set conservative TriggerDistance

3. **Over-Scaling**:

   - Too small TriggerDistance -> too many scales
   - **Solution**: Match TriggerDistance to symbol volatility
   - Example: EURUSD 200pts (= 0.00200; the point size is read from each symbol, so 200 points are 2.00 on XAUUSD)

4. **Capital Requirements**:

   - Scaling requires available margin for additional positions
   - **Solution**: Never use more than 30% of account for scaling strategies
   - **Solution**: Monitor margin levels continuously

5. **Correlation Risk**:

   - Scaling on correlated symbols compounds risk
   - Example: Pyramiding EURUSD and GBPUSD simultaneously
   - **Solution**: Limit scaling to uncorrelated symbols

---

## 🔧 Configuration Profiles

### Conservative Pyramiding

```go
Mode:             Pyramiding,
TriggerDistance:  200,     // Wide spacing = fewer scales
MinProfitToScale: 150,     // Wait for confirmed profit
MaxScales:        2,       // Limited scales = controlled risk
ScaleLotSize:     0.03,    // Smaller additions
TotalMaxLotSize:  0.50,    // Conservative total
StopLossPerScale: 200,     // Wider stops
```

### Aggressive Pyramiding

```go
Mode:             Pyramiding,
TriggerDistance:  100,     // Tighter spacing = more scales
MinProfitToScale: 50,      // Scale sooner
MaxScales:        5,       // More aggressive
ScaleLotSize:     0.05,    // Larger additions
TotalMaxLotSize:  1.50,    // Higher total
StopLossPerScale: 100,     // Tighter stops
```

### Cautious Averaging Down (High Risk!)

```go
Mode:             AveragingDown,
TriggerDistance:  200,     // Average every 200pts down
MaxLossToAverage: 300,     // CRITICAL: Stop early
MaxScales:        2,       // LIMIT exposure
ScaleLotSize:     0.03,    // SMALL additions only
TotalMaxLotSize:  0.30,    // Very conservative
StopLossPerScale: 150,
```

### Gradual Profit Taking

```go
Mode:             ScaleOut,
TriggerDistance:  200,     // Close at +200 / +400 / +600 / +800pts
ScaleOutPercent:  0.25,    // Close 25% of the largest position each time
MaxScales:        4,       // At most 4 partial closes
```

---

## When to Use Position Scaling

### Pyramiding - Suitable Conditions

- **Strong trending markets**: Clear directional movement confirmed by multiple timeframes
- **Breakout trades**: After resistance/support breaks, trend continuation expected
- **Momentum strategies**: When price accelerates in your favor
- **Good risk/reward**: Each scale-in has better R:R than previous
- **Examples**:
  - NFP data confirms USD strength, pyramid EURUSD shorts
  - Gold breaks $2000 resistance, pyramid long positions
  - GBPUSD trending down 500+ pips, add to shorts

### Pyramiding - Unsuitable Conditions

- **Range-bound markets**: Sideways movement will create losing scales
- **Counter-trend positions**: Fighting main trend direction
- **High volatility**: Whipsaws will trigger stops on scales
- **Weak account balance**: Need margin for multiple positions
- **Overextended moves**: Price already moved 80%+ of expected range

### Averaging Down - Suitable Conditions (VERY LIMITED!)

- **ONLY at major support/resistance**: Multi-year levels, round numbers
- **Mean reversion expected**: Statistical probability of bounce
- **Small position sizes**: Never more than 20% of account
- **Strict loss limits**: Always set MaxLossToAverage
- **High confidence setups**: Multiple confirmations required
- **Examples**:
  - EURUSD at 1.0500 major support (tested 5+ times)
  - Gold at $1900 psychological level with RSI oversold
  - Index retracements to key Fibonacci levels

### Averaging Down - AVOID These Conditions

- **Trending markets**: NEVER average against a strong trend
- **Breaking support/resistance**: Price acceleration usually continues
- **News-driven moves**: Fundamental shifts create sustained trends
- **Weak support zones**: Untested or recently broken levels
- **Large positions**: Risk of catastrophic loss
- **Beginners**: Requires experience to manage psychology

### Scale Out - Suitable Conditions

- **Large profitable positions**: Want to secure partial gains
- **Volatile markets**: Protect profits from reversals
- **Approaching resistance**: Profit-taking expected
- **Uncertain outlook**: Bullish but cautious
- **Risk management**: Reducing exposure gradually
- **Examples**:
  - Gold up $50, approaching resistance, scale out 30%
  - EURUSD +300 pips profit, volatility increasing
  - Large swing trade position, take profits in stages

---

## Configuration Location

**Runtime configuration is located in**: `examples/demos/main.go -> RunOrchestrator_PositionScaler()`

This separation exists because:

1. **Code Reusability**: Same orchestrator runs different strategies (Pyramiding/Averaging/ScaleOut)

2. **Quick Testing**: Change mode and parameters without modifying strategy logic

3. **User Examples**: Shows how to properly configure each scaling mode

4. **Centralized Entry Point**: Single place for all orchestrator launches

