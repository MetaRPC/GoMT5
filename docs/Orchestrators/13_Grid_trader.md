# GridTrader - Grid Trading Orchestrator

## Description

**🧩 GridTrader** is an automated grid trading orchestrator designed for **RANGE-BOUND and SIDEWAYS markets**. It places a "fishing net" of pending BUY and SELL orders at fixed price intervals, profiting from oscillations as price bounces between levels. Perfect for choppy, non-trending markets where price moves back and forth.

**Strategy Principle**: A symmetrical grid of pending orders is created above AND below current price:
- **SELL LIMIT orders** placed above current price (profit when price rises then falls back)
- **BUY LIMIT orders** placed below current price (profit when price falls then rises back)
- Each level acts as a profit target for the opposite direction
- Grid automatically rebuilds when price moves significantly
- Works 24/7 capturing small profits from price oscillations

**File**: `examples/demos/orchestrators/13_grid_trader.go`

> **IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY**
>
> This orchestrator is a demonstration showing how GoMT5 methods function and combine into automated strategies. **NOT a production-ready trading strategy!** Always test on demo accounts first.

---

## Architecture

```
GRID TRADER ORCHESTRATOR
    |
MT5Sugar Instance
    |
    |  |  |
BuyLimit  SellLimit  GetPositions
(below)   (above)    (monitoring)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`BuyLimit`, `SellLimit`, `GetPositionsBySymbol`)
- **MT5Service**: Mid-level API for cleanup operations (`CloseOrder`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking
- **protobuf types**: `pb.OrderCloseRequest` for order management

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Symbol` | string | *required* | Trading instrument (e.g., "EURUSD") |
| `GridSize` | int | `5` | Number of grid levels above AND below current price |
| `GridStep` | float64 | `100` | Distance between levels in points (100pts=10pips for EURUSD) |
| `LotSize` | float64 | `0.01` | Volume (lots) for each order in the grid |
| `MaxPositions` | int | `10` | Maximum concurrent open positions allowed |
| `TakeProfit` | float64 | `0` | Take profit in points (0=use GridStep as TP) |
| `StopLoss` | float64 | `0` | Stop loss in points (0=no SL) |
| `CheckInterval` | time.Duration | `5s` | How often to check and update grid |
| `RebuildOnFill` | bool | `false` | Whether to rebuild entire grid when order fills |

### Configuration Example

```go
config := orchestrators.GridTraderConfig{
    Symbol:        "EURUSD",
    GridSize:      5,                // 5 levels each side (10 total orders)
    GridStep:      100,              // 100 points = 10 pips spacing
    LotSize:       0.01,             // 0.01 lots per order
    MaxPositions:  10,               // Max 10 concurrent positions
    TakeProfit:    0,                // Use grid step as TP
    StopLoss:      0,                // No SL by default
    CheckInterval: 5 * time.Second,  // Check every 5 seconds
    RebuildOnFill: false,            // Don't rebuild on each fill
}

gridTrader := orchestrators.NewGridTrader(sugar, config)
```

---

## How to Run

You can execute this orchestrator using several command variations:

```bash
# Navigate to examples directory
cd examples/demos

# Option 1: By number
go run main.go 13

# Option 2: By keyword
go run main.go grid
```

Both commands will launch the **GridTrader** orchestrator with the configuration specified in `main.go -> RunOrchestrator_Grid()`.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Initialize symbol parameters (digits, point value)
  |
Get current price (Bid/Ask midpoint)
  |
Calculate grid levels:
  Levels ABOVE: currentPrice + (1×step), (2×step), ..., (GridSize×step)
  Levels BELOW: currentPrice - (1×step), (2×step), ..., (GridSize×step)
  |
FOR each level ABOVE current price:
  | Place SELL LIMIT at level
  | Set TP at (level - GridStep)
  |
FOR each level BELOW current price:
  | Place BUY LIMIT at level
  | Set TP at (level + GridStep)
  |
MONITORING LOOP (every CheckInterval):
  | Get current positions count
  | Check if max positions reached -> pause if yes
  | Get current price
  | If price moved > 2×GridStep from center -> rebuild grid
  | Track total P/L of open positions
  |
ON STOP:
  | Cancel all pending orders
  | Close all positions (if configured)
  |
END
```

### Step-by-step Description

#### 1. Initialization

> Lines `272–286`  


```go
func (g *GridTrader) initializeSymbol() error {
    priceInfo, err := g.sugar.GetPriceInfo(g.config.Symbol)
    if err != nil {
        return fmt.Errorf("failed to get price: %w", err)
    }
    g.currentPrice = priceInfo.Bid

    g.digits = 5         // Symbol decimal digits
    g.point = 0.00001    // Point value for EURUSD

    return nil
}
```

**What happens**:

- Current Bid price is obtained for grid center calculation
- Symbol parameters (digits, point value) are initialized
- For EURUSD: 5 digits, point = 0.00001

#### 2. Building the Grid

> Lines `289–331`  


```go
func (g *GridTrader) buildGrid() error {
    g.cleanupOrders()  // Remove old grid first

    // Get fresh current price
    priceInfo, err := g.sugar.GetPriceInfo(g.config.Symbol)
    g.currentPrice = (priceInfo.Bid + priceInfo.Ask) / 2

    // Calculate grid step in price units
    gridStepPrice := g.config.GridStep * g.point

    // Build levels above and below
    for i := 1; i <= g.config.GridSize; i++ {
        levelAbove := g.currentPrice + float64(i)*gridStepPrice
        levelBelow := g.currentPrice - float64(i)*gridStepPrice

        g.placeSellLimit(levelAbove)  // SELL above
        g.placeBuyLimit(levelBelow)   // BUY below
    }

    return nil
}
```

**Critical Details**:

- Grid center = (Bid + Ask) / 2 for balanced placement
- Each level is `GridStep × point` distance from center
- SELL LIMIT orders go **ABOVE** current price
- BUY LIMIT orders go **BELOW** current price
- Old grid is cleaned up before rebuilding

**Example with GridSize=5, GridStep=100, CurrentPrice=1.10000**:

```
Grid levels calculated:
levelAbove[1] = 1.10000 + (1×0.00100) = 1.10100  SELL LIMIT
levelAbove[2] = 1.10000 + (2×0.00100) = 1.10200  SELL LIMIT
levelAbove[3] = 1.10000 + (3×0.00100) = 1.10300  SELL LIMIT
levelAbove[4] = 1.10000 + (4×0.00100) = 1.10400  SELL LIMIT
levelAbove[5] = 1.10000 + (5×0.00100) = 1.10500  SELL LIMIT

levelBelow[1] = 1.10000 - (1×0.00100) = 1.09900  BUY LIMIT
levelBelow[2] = 1.10000 - (2×0.00100) = 1.09800  BUY LIMIT
levelBelow[3] = 1.10000 - (3×0.00100) = 1.09700  BUY LIMIT
levelBelow[4] = 1.10000 - (4×0.00100) = 1.09600  BUY LIMIT
levelBelow[5] = 1.10000 - (5×0.00100) = 1.09500  BUY LIMIT

Total: 10 pending orders
```

#### 3. Placing BUY LIMIT Orders
 
> Lines `334–379`  


```go
func (g *GridTrader) placeBuyLimit(price float64) error {
    // Calculate TP/SL
    tp := price + g.config.GridStep*g.point  // TP one grid level up
    sl := 0.0
    if g.config.StopLoss > 0 {
        sl = price - g.config.StopLoss*g.point
    }

    // Place order
    ticket, err := g.sugar.BuyLimitWithSLTP(
        g.config.Symbol,
        g.config.LotSize,
        price,
        sl,
        tp,
    )

    g.activeOrders = append(g.activeOrders, ticket)
    return nil
}
```

**Key Points**:

- **TP is set one grid level ABOVE entry** (price + GridStep)
- This means: when BUY fills at 1.09900, TP is at 1.10000
- As price rises back to grid center, position closes with profit
- SL is optional (default 0 = no SL)

**Example BUY order placement**:
```
Price: 1.09900 (BUY LIMIT)
TP:    1.10000 (price + 100pts)
SL:    none (StopLoss = 0)
```

#### 4. Placing SELL LIMIT Orders

> Lines `382–427`  


```go
func (g *GridTrader) placeSellLimit(price float64) error {
    // Calculate TP/SL
    tp := price - g.config.GridStep*g.point  // TP one grid level down
    sl := 0.0
    if g.config.StopLoss > 0 {
        sl = price + g.config.StopLoss*g.point
    }

    // Place order
    ticket, err := g.sugar.SellLimitWithSLTP(
        g.config.Symbol,
        g.config.LotSize,
        price,
        sl,
        tp,
    )

    g.activeOrders = append(g.activeOrders, ticket)
    return nil
}
```

**Key Points**:

- **TP is set one grid level BELOW entry** (price - GridStep)
- This means: when SELL fills at 1.10100, TP is at 1.10000
- As price falls back to grid center, position closes with profit

**Example SELL order placement**:
```
Price: 1.10100 (SELL LIMIT)
TP:    1.10000 (price - 100pts)
SL:    none (StopLoss = 0)
```

#### 5. Monitoring Loop

> Lines `430–501`  


```go
func (g *GridTrader) monitorLoop() {
    ticker := time.NewTicker(g.config.CheckInterval)
    defer ticker.Stop()

    for {
        select {
        case <-g.GetContext().Done():
            return
        case <-ticker.C:
            g.checkAndUpdateGrid()
        }
    }
}
```

**Monitoring actions every CheckInterval**:

1. **Get current positions**: Count open positions for the symbol
2. **Check max positions**: If `CurrentPositions >= MaxPositions`, pause
3. **Get current price**: Fetch latest Bid/Ask
4. **Check price drift**: If price moved > 2×GridStep from center, rebuild grid
5. **Track P/L**: Sum profit/loss of all open positions
6. **Update metrics**: Store in orchestrator metrics for reporting

**Grid rebuild trigger**:

```go
priceDiff := currentPrice - g.currentPrice
gridStepPrice := g.config.GridStep * g.point

// Rebuild if price moved more than 2 grid steps
if priceDiff > 2*gridStepPrice || priceDiff < -2*gridStepPrice {
    g.buildGrid()  // Rebuild centered on new price
}
```

#### 6. Cleanup on Stop
 
> Lines `504–536`  


```go
func (g *GridTrader) cleanupOrders() {
    ctx := context.Background()
    service := g.sugar.GetService()

    for _, ticket := range g.activeOrders {
        req := &pb.OrderCloseRequest{Ticket: ticket}
        service.CloseOrder(ctx, req)  // Delete pending order
    }

    g.activeOrders = make([]uint64, 0)  // Clear list
}
```

**What happens**:

- All pending orders in `activeOrders` list are canceled
- Uses `Service.CloseOrder()` which can delete pending orders
- Grid is fully cleaned up to prevent hanging orders
- Orchestrator exits cleanly

---

## Grid Visualization

### Scenario: GridSize=5, GridStep=100, CurrentPrice=1.10000

```
1.10500  [SELL LIMIT #5]  <- Level +5 (+500pts)      TP: 1.10400
         +500pts from center                          SL: none

1.10400  [SELL LIMIT #4]  <- Level +4 (+400pts)      TP: 1.10300
         +400pts from center                          SL: none

1.10300  [SELL LIMIT #3]  <- Level +3 (+300pts)      TP: 1.10200
         +300pts from center                          SL: none

1.10200  [SELL LIMIT #2]  <- Level +2 (+200pts)      TP: 1.10100
         +200pts from center                          SL: none

1.10100  [SELL LIMIT #1]  <- Level +1 (+100pts)      TP: 1.10000
         +100pts from center                          SL: none

1.10000  >>> CURRENT PRICE <<<  [Grid Center]

1.09900  [BUY LIMIT #1]   <- Level -1 (-100pts)      TP: 1.10000
         -100pts from center                          SL: none

1.09800  [BUY LIMIT #2]   <- Level -2 (-200pts)      TP: 1.09900
         -200pts from center                          SL: none

1.09700  [BUY LIMIT #3]   <- Level -3 (-300pts)      TP: 1.09800
         -300pts from center                          SL: none

1.09600  [BUY LIMIT #4]   <- Level -4 (-400pts)      TP: 1.09700
         -400pts from center                          SL: none

1.09500  [BUY LIMIT #5]   <- Level -5 (-500pts)      TP: 1.09600
         -500pts from center                          SL: none

Total: 10 pending orders (5 BUY + 5 SELL)
```

### Working Scenario: Price Movement & Profit

**Time 0: Grid Placed**
```
Price:     1.10000
Status:    All 10 pending orders active
Positions: 0
P/L:       $0.00
```

**Time 1: Price Drops to 1.09900**
```
Event:     BUY LIMIT #1 @ 1.09900 TRIGGERED
Position:  LONG 0.01 lots @ 1.09900
TP:        1.10000 (next grid level up)
Status:    1 position open, 9 pending orders remain
```

**Time 2: Price Returns to 1.10000**
```
Event:     BUY position CLOSED by TP @ 1.10000
Profit:    +100 points = +$10 (for 0.01 lots)
Action:    Grid rebuilds, new 10 pending orders placed
Total P/L: +$10
```

**Time 3: Price Rises to 1.10200**
```
Event:     SELL LIMIT #2 @ 1.10200 TRIGGERED
Position:  SHORT 0.01 lots @ 1.10200
TP:        1.10100 (next grid level down)
Status:    1 position open, 9 pending orders remain
```

**Time 4: Price Falls to 1.10100**
```
Event:     SELL position CLOSED by TP @ 1.10100
Profit:    +100 points = +$10
Total P/L: +$20
```

**Time 5: Multiple Oscillations**
```
Price oscillates: 1.09900 -> 1.10000 -> 1.09800 -> 1.10100
Result:  4 trades triggered, each with +100pts profit
Total profit: 4 × $10 = $40
```

**Summary**:

- Range-bound market (1.09500 - 1.10500): Grid captures EVERY oscillation
- Without grid: Hard to catch manually, many small movements missed
- With grid: Automatic profit on each bounce
- Typical result: 10-20 small wins per day in choppy markets

---

## Used MT5Sugar Methods

### 1. BuyLimit / BuyLimitWithSLTP

```go
func (s *MT5Sugar) BuyLimit(symbol string, volume float64, price float64) (uint64, error)

func (s *MT5Sugar) BuyLimitWithSLTP(
    symbol string,
    volume float64,
    price float64,
    sl float64,
    tp float64,
) (uint64, error)
```

**Purpose**: Places pending Buy Limit order at specified price.

**Parameters in GridTrader**:

- `symbol`: Trading instrument (from config.Symbol)
- `volume`: `config.LotSize` (e.g., 0.01)
- `price`: Calculated grid level BELOW current price
- `sl`: Stop loss price (0 if disabled)
- `tp`: Take profit price (price + GridStep)

**Returns**: Order ticket number or error

### 2. SellLimit / SellLimitWithSLTP

```go
func (s *MT5Sugar) SellLimit(symbol string, volume float64, price float64) (uint64, error)

func (s *MT5Sugar) SellLimitWithSLTP(
    symbol string,
    volume float64,
    price float64,
    sl float64,
    tp float64,
) (uint64, error)
```

**Purpose**: Places pending Sell Limit order at specified price.

**Parameters in GridTrader**:

- `price`: Calculated grid level ABOVE current price
- `tp`: Take profit price (price - GridStep)

### 3. GetPositionsBySymbol

```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*mt5.Position, error)
```

**Purpose**: Retrieves all open positions for a symbol.

**Usage in GridTrader**:

- Called every CheckInterval in monitoring loop
- Counts positions to enforce MaxPositions limit
- Calculates total P/L across all positions

### 4. GetPriceInfo

```go
func (s *MT5Sugar) GetPriceInfo(symbol string) (*mt5.PriceInfo, error)
```

**Purpose**: Gets current Bid/Ask prices for a symbol.

**Usage in GridTrader**:

- Gets current price for grid center calculation
- Monitors price drift to trigger grid rebuilds

### 5. Service.CloseOrder (via GetService)

```go
service := sugar.GetService()
service.CloseOrder(ctx context.Context, req *pb.OrderCloseRequest) (*pb.OrderCloseData, error)
```

**Purpose**: Closes position or deletes pending order.

**Usage in GridTrader**:

- Cleanup operation when stopping orchestrator
- Deletes all pending orders in `activeOrders` list
- Ensures no hanging orders after exit

---

## Risk Management

### Maximum Risk Calculation

```
Risk per level = LotSize × StopLoss × PointValue

For EURUSD (PointValue ~$10 per 1.0 lot):
- LotSize = 0.01
- StopLoss = 500 points (if configured)
- Risk per level = 0.01 × 500 × $10 = $50

Total maximum risk (if all orders trigger):
- GridSize = 5 -> 10 orders total
- Maximum risk = 10 × $50 = $500
```

### Risk Considerations

1. **No Stop Loss Default**:

   - Default `StopLoss = 0` means **NO stop loss protection**
   - Grid trading relies on price returning to range
   - If range breaks violently, losses can be severe
   - **Recommendation**: Set `StopLoss = 3 × GridStep` for safety

2. **Overexposure Risk**:

   - `GridSize × LotSize` = total volume exposure
   - Example: GridSize=10, LotSize=0.1 -> 1.0 lots total
   - Ensure account has adequate margin for full grid

3. **Trending Market Danger**:

   - Grid designed for RANGES, not TRENDS
   - In strong trend: all orders on one side fill
   - Result: Multiple losing positions, no offsetting wins
   - **Solution**: Only use during range-bound sessions

4. **MaxPositions Safety**:

   - Limits concurrent exposure
   - Prevents runaway position accumulation
   - Set to `GridSize × 2` or higher for buffer

### Recommended Settings by Account Size

**Small Account ($500-$1000)**:
```go
GridSize:      3,    // Limited exposure
GridStep:      100,  // Standard spacing
LotSize:       0.01, // Minimum lots
MaxPositions:  6,    // Safety limit
StopLoss:      300,  // 30 pip SL for protection
```

**Medium Account ($5000-$10000)**:
```go
GridSize:      5,    // Balanced grid
GridStep:      100,
LotSize:       0.02,
MaxPositions:  10,
StopLoss:      200,  // Tighter SL
```

**Large Account ($10000+)**:
```go
GridSize:      10,   // Wide grid
GridStep:      50,   // Tighter spacing
LotSize:       0.05,
MaxPositions:  20,
StopLoss:      0,    // No SL (experienced only)
```

---

## When to Use Grid Trading

### Suitable Conditions

- **Range-bound markets**: Price oscillating within defined levels
- **Low volatility**: Small, predictable price movements
- **Asian session**: Typically calmer trading hours
- **Post-news consolidation**: After big moves, price settles into range
- **Clear support/resistance**: Known price boundaries
- **Sideways forex pairs**: EURUSD at night, consolidating pairs
- **Cryptocurrency weekends**: Lower volatility, tighter ranges

### Unsuitable Conditions

- **Trending markets**: Strong directional movement
- **High-impact news events**: NFP, FOMC, Brexit votes
- **Breakout scenarios**: Price breaking key levels
- **Low liquidity**: Wide spreads, gap risk
- **High volatility**: Violent price swings
- **Pre-news periods**: Uncertainty before announcements

### Best Market Sessions

| Session | Suitability | Notes |
|---------|-------------|-------|
| **Asian (Tokyo)** | Excellent | Low volatility, tight ranges |
| **European (London)** | Moderate | Higher volatility, watch for trends |
| **US (New York)** | Moderate | Can be volatile, use wider grids |
| **Overlap (London+NY)** | Risky | High volatility, trending tendency |
| **Weekend (Crypto)** | Good | Lower volatility in crypto markets |

---

## Configuration Location

**Runtime configuration is located in**: `examples/demos/main.go -> RunOrchestrator_Grid()`

This separation exists because:

1. **Code reusability**: Same orchestrator runs with different parameters

2. **Quick testing**: Change numbers without modifying strategy logic

3. **User examples**: Shows how to properly configure the orchestrator

4. **Centralized entry point**: Single place for all orchestrator launches


