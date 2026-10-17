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
   - Partially closes a position once its profit (measured from the base entry) reaches `TriggerDistance`
   - Takes partial profits while letting remainder run
   - Reduces risk exposure over time
   - Secures gains while maintaining market exposure
   - Best for volatile markets or profit protection

> **The scaler does not open the base position.** It only works on positions that are already open (opened by you or another program). Positions are grouped per symbol: the first position seen for a symbol becomes the group's base (`BaseTicket`, `BasePrice`, direction), and all profit/loss checks are measured from that base entry.

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
    |  |  |  |
GetOpenPositions  GetPriceInfo  Buy/SellMarketWithSLTP  GetPositionsBySymbol + ClosePositionPartial
(monitoring)      (tracking)    (scaling in)            (scaling out)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetOpenPositions`, `GetPriceInfo`, `BuyMarketWithSLTP`, `SellMarketWithSLTP`, `GetPositionsBySymbol`, `ClosePositionPartial`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo` for position data structures

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Mode` | ScalingMode | `Pyramiding` | Scaling strategy (Pyramiding/AveragingDown/ScaleOut) |
| `TriggerDistance` | float64 | `200` | Points moved to trigger next scale |
| `MinProfitToScale` | float64 | `100` | Min profit before scaling (pyramiding mode) |
| `MaxLossToAverage` | float64 | `500` | Max loss before stopping averaging (critical safety) |
| `InitialLotSize` | float64 | `0.10` | Declared "initial position size", but **not read by the current code** (the scaler never opens the base position) |
| `ScaleLotSize` | float64 | `0.05` | Size for each scale-in operation |
| `MaxScales` | int | `3` | Maximum number of scale operations (scale-ins, or partial closes in ScaleOut mode) |
| `ReducePerScale` | float64 | `0.0` | If > 0: next size = `ScaleLotSize × (1 − ReducePerScale × ScaleCount)` |
| `TotalMaxLotSize` | float64 | `1.0` | Maximum total position size across all scales |
| `ScaleOutPercent` | float64 | `0.33` | Fraction of the position to close at each scale-out (volume is not rounded to the volume step) |
| `StopLossPerScale` | float64 | `150` | SL distance in points for each scale-in |
| `Symbols` | []string | `[]string{symbol}` | Symbols to manage (empty = all symbols) |
| `CheckInterval` | time.Duration | `5s` | How often to check for scaling opportunities |

Defaults are those of `orchestrators.DefaultPositionScalerConfig(symbol string)`.

> **Point size:** all distances are in points, and the scaler does not query the symbol's point size: `getSymbolPoint()` caches a hard-coded `0.00001` (5-digit symbol) for every symbol. Distances are therefore only correct for 5-digit symbols such as EURUSD.

### Configuration Example

```go
// Pyramiding Configuration
config := orchestrators.PositionScalerConfig{
    Mode:             orchestrators.Pyramiding,
    TriggerDistance:  200,           // Scale every 200 points
    MinProfitToScale: 100,           // Activate after +100 points profit
    InitialLotSize:   0.10,          // Not read by the current code
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

Both commands will launch the **PositionScaler** orchestrator with the configuration specified in `main.go -> RunOrchestrator_PositionScaler()` (Pyramiding mode on `cfg.TestSymbol`). The runner runs the scaler for 5 minutes (300 s), then calls `Stop()` and prints the final metrics. Open a position on that symbol yourself - the scaler only adds to existing positions.

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
MONITORING LOOP (every CheckInterval):
  | Get all open positions
  |
  UPDATE GROUPS (one PositionGroup per symbol):
    | Skip positions whose symbol is not tracked (if Symbols list specified)
    | First position seen for a symbol -> new group (base ticket, base price, direction)
    | Recalculate total lot size for each symbol
    | Remove groups whose symbol has no open position any more
  |
  FOR each group: Should we scale this group?
      |
      | Check 1: Have we reached MaxScales? -> YES: Skip
      | Check 2: Would total exceed TotalMaxLotSize? -> YES: Skip
      | Check 3: Get current price; profit in points vs BASE price
      |          (BUY: Bid - BasePrice, SELL: BasePrice - Ask)
      |
      | MODE-SPECIFIC CHECKS:
      |
      | PYRAMIDING MODE:
      |   | Is profit >= MinProfitToScale? -> NO: Skip
      |   | First scale: profit >= TriggerDistance? -> NO: Skip
      |   | Later scales: moved TriggerDistance from last scale price? -> NO: Skip
      |   | YES -> SCALE IN (open additional position)
      |
      | AVERAGING DOWN MODE:
      |   | Is position in loss? -> NO: Skip
      |   | Is loss <= MaxLossToAverage? -> NO: Skip (safety)
      |   | First scale: loss >= TriggerDistance? -> NO: Skip
      |   | Later scales: moved TriggerDistance against last scale price? -> NO: Skip
      |   | YES -> SCALE IN (open additional position at better price)
      |
      | SCALE OUT MODE:
      |   | Is profit (from base) >= TriggerDistance? -> NO: Skip
      |   | YES -> SCALE OUT (close ScaleOutPercent of position)
      |
  FOR SCALE IN operation:
    | Calculate next scale lot size
    | Get current price
    | Calculate stop loss for scale-in (BUY: Bid - StopLossPerScale*point,
    |                                   SELL: Ask + StopLossPerScale*point), TP = 0
    | BuyMarketWithSLTP / SellMarketWithSLTP
    | LastScalePrice = Ask (BUY) / Bid (SELL)
    | Add ticket to PositionGroup.ScaleTickets
    | Increment PositionGroup.ScaleCount
    | Update PositionGroup.TotalLotSize
    | Update metrics
  |
  FOR SCALE OUT operation:
    | GetPositionsBySymbol, take the FIRST position returned
    | Calculate close volume (its volume × ScaleOutPercent, not rounded)
    | ClosePositionPartial
    | Update PositionGroup.ScaleCount
    | Update PositionGroup.TotalLotSize
    | Update metrics
  |
REPEAT until stopped
  |
ON STOP:
  | Cancel monitoring loop
  | Clear all position groups
  |
END
```

### Step-by-step Description

#### 1. Initialization
 
> Lines `194–202`  
> `examples/demos/orchestrators/12_position_scaler.go`


```go
func NewPositionScaler(sugar *mt5.MT5Sugar, config PositionScalerConfig) *PositionScaler {
    return &PositionScaler{
        BaseOrchestrator: NewBaseOrchestrator("Position Scaler"),
        sugar:            sugar,
        config:           config,
        trackedGroups:    make(map[string]*PositionGroup),
        symbolPoints:     make(map[string]float64),
    }
}
```

**What happens**:

- Creates new scaler instance with provided configuration
- Initializes empty position group tracker (maps symbol to PositionGroup)
- Each symbol gets one PositionGroup tracking all related positions
- The point cache is filled by `getSymbolPoint()` with a hard-coded `0.00001` for every symbol (the symbol is not queried)

**PositionGroup structure**:
```go
type PositionGroup struct {
    Symbol         string    // Symbol being traded
    IsBuy          bool      // Direction (true=BUY, false=SELL)
    BaseTicket     uint64    // Original position ticket
    ScaleTickets   []uint64  // Tickets of scale-in positions
    BasePrice      float64   // Entry price of base position
    TotalLotSize   float64   // Total position size across all scales
    ScaleCount     int       // Number of scales executed
    LastScalePrice float64   // Price of last scale-in
    LastScaleTime  time.Time // When last scaled
}
```

#### 2. Starting the Scaler


```go
func (p *PositionScaler) Start() error {
    if p.IsRunning() {
        return fmt.Errorf("position scaler already running")
    }

    // Create context
    ctx, cancel := context.WithCancel(context.Background())
    p.SetContext(ctx, cancel)

    // Mark as started
    p.MarkStarted()

    // Start monitoring loop
    go p.monitorLoop()

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
3. **Check opportunities**: Calls `checkScalingOpportunities()` to process positions
4. **Graceful exit**: Stops cleanly when context is cancelled

#### 4. Scaling Decision Logic


`shouldScale()` (lines 333-437) first returns `false` when `ScaleCount >= MaxScales` or when `TotalLotSize + nextScaleSize > TotalMaxLotSize`, then computes `profitPoints` from the group's base price (BUY: `(Bid - BasePrice) / point`, SELL: `(BasePrice - Ask) / point`) and switches on the mode:

**For Pyramiding Mode** (lines 371-396):
```go
	case Pyramiding:
		// For pyramiding: position must be in profit
		if profitPoints < p.config.MinProfitToScale {
			return false
		}

		// Check if price moved enough from last scale
		if group.ScaleCount > 0 {
			priceMove := 0.0
			if group.IsBuy {
				priceMove = (currentPrice - group.LastScalePrice) / point
			} else {
				priceMove = (group.LastScalePrice - currentPrice) / point
			}

			if priceMove < p.config.TriggerDistance {
				return false
			}
		} else {
			// First scale: check distance from base
			if profitPoints < p.config.TriggerDistance {
				return false
			}
		}

		return true
```

**For Averaging Down Mode** (lines 398-426):
```go
	case AveragingDown:
		// For averaging: position must be in loss but not too much
		if profitPoints > 0 {
			return false
		}

		if -profitPoints > p.config.MaxLossToAverage {
			return false
		}

		// Check if price moved enough from last average
		if group.ScaleCount > 0 {
			priceMove := 0.0
			if group.IsBuy {
				priceMove = (group.LastScalePrice - currentPrice) / point
			} else {
				priceMove = (currentPrice - group.LastScalePrice) / point
			}

			if priceMove < p.config.TriggerDistance {
				return false
			}
		} else {
			if -profitPoints < p.config.TriggerDistance {
				return false
			}
		}

		return true
```

**For Scale Out Mode** (lines 428-433):
```go
	case ScaleOut:
		// For scale-out: check if profit target reached
		if profitPoints < p.config.TriggerDistance {
			return false
		}
		return true
```

> Note: `MaxLossToAverage` is measured from the base entry only, and it only stops further averaging - it does not close anything. In ScaleOut mode the profit is always measured from the base price and no "last scale-out price" is kept, so once the profit reaches `TriggerDistance` the scaler partially closes again on every check until `MaxScales` is reached.

---

## 💡 Visual Examples

### Scenario 1: PYRAMIDING MODE (Recommended)

**Configuration**:

- Entry: 1.10000 BUY EURUSD 0.10 lots
- TriggerDistance: 200 points
- MinProfitToScale: 100 points
- ScaleLotSize: 0.05 lots
- MaxScales: 3
- Prices ignore the spread: a BUY scale price is recorded as the Ask while the move is checked on the Bid, so each later trigger needs TriggerDistance plus the spread

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

- Entry: 1.10000 BUY EURUSD 0.40 lots (opened by you)
- TriggerDistance: 300 points
- ScaleOutPercent: 0.50 (50% - chosen so every close volume is a multiple of 0.01)
- MaxScales: 3
- CheckInterval: 5s

```
TIME 0: Large Position Opened

Price:        1.10000 (Entry)
Position:     0.40 lots
Profit:       0 points
Status:       Waiting for profit target (+300pts from entry)

TIME 1: Check at 1.10300 (FIRST SCALE OUT)

Price:        1.10300
Profit:       +300 points (from base entry)
Action:       SCALE OUT! Close 50% of the first position
Close:        0.40 × 0.50 = 0.20 lots CLOSED
Locked Profit: 0.20 lots × 300pts = +$60 SECURED
Remaining:    0.20 lots still open

TIME 2: Next check, 5 s later, price 1.10310 (SECOND SCALE OUT)

Price:        1.10310
Profit:       +310 points from base - still >= TriggerDistance
Action:       SCALE OUT again (no new +300pts move is required:
              the profit is always measured from the base entry)
Close:        0.20 × 0.50 = 0.10 lots CLOSED
Locked Profit: 0.10 lots × 310pts = +$31 SECURED
Remaining:    0.10 lots still open

TIME 3: Next check, price 1.10305 (THIRD SCALE OUT)

Price:        1.10305
Close:        0.10 × 0.50 = 0.05 lots CLOSED
Locked Profit: 0.05 lots × 305pts = ~+$15 SECURED
Remaining:    0.05 lots still open
Total Locked: $60 + $31 + $15 = ~$106 secured
Status:       MAX SCALES REACHED (3/3) - no more scale-outs

RESULT:
  - Three partial closes on three consecutive checks (about 15 s apart
    in total), all at roughly +300pts
  - The scale-outs are NOT spaced TriggerDistance apart: to take profits
    at +300 / +600 / +900 you would have to change the code (for example
    track a last scale-out price, as the scale-in modes do)
  - 0.05 lots stay open for further gains
  - With ScaleOutPercent 0.33 and 0.30 lots the close volume would be
    0.099 lots: it is sent unrounded and the broker can reject it
```

---

## Used MT5Sugar Methods

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
- Sets stop loss based on StopLossPerScale parameter (BUY: Bid - StopLossPerScale*point, SELL: Ask + StopLossPerScale*point)
- TP is always 0 for scale-ins
- Volume and prices are sent as given (MT5Sugar does not round them)

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
- The scaler partially closes the FIRST position in the returned slice (`positions[0]`) - it does not search for the largest one
- Returns the open positions filtered by symbol (nil slice when none)

**Returns**: Slice of positions for the symbol

### 5. ClosePositionPartial

```go
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error
```

**Purpose**: Closes a partial amount of a position.

**Usage in PositionScaler**:

- Core method for Scale Out mode
- Closes `ScaleOutPercent` of the first position's volume (the volume is sent as calculated, not rounded to the symbol's volume step - a value like 0.30 × 0.33 = 0.099 can be rejected by the broker)
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
   - Example: EURUSD 200pts (the point size is hard-coded to 0.00001, so non-5-digit symbols need a code change first)

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
TriggerDistance:  200,     // Start once profit from entry >= 200pts (then one close per check)
ScaleOutPercent:  0.25,    // Close 25% of the first position each time
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

