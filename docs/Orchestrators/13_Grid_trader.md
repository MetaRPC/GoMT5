# GridTrader - Grid Trading Orchestrator

## Description

**🧩 GridTrader** is an automated grid trading orchestrator designed for **RANGE-BOUND and SIDEWAYS markets**. It places a "fishing net" of pending BUY and SELL orders at fixed price intervals, profiting from oscillations as price bounces between levels. Perfect for choppy, non-trending markets where price moves back and forth.

**Strategy Principle**: A symmetrical grid of pending orders is created above AND below current price:
- **SELL LIMIT orders** placed above current price (profit when price rises then falls back)
- **BUY LIMIT orders** placed below current price (profit when price falls then rises back)
- With `TakeProfit > 0` (or `StopLoss > 0`) each order gets a TP; with `TakeProfit = 0` that TP is one grid step away. With the default config (both 0) the orders are placed **without TP/SL**
- Grid is rebuilt around the new mid price when price drifts more than 2 grid steps from the grid center
- Captures small profits from price oscillations while it runs

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
BuyLimit(WithSLTP)  SellLimit(WithSLTP)  GetPositionsBySymbol  GetPriceInfo
(below)             (above)              (monitoring)          (grid center)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetPriceInfo`, `BuyLimit`, `BuyLimitWithSLTP`, `SellLimit`, `SellLimitWithSLTP`, `GetPositionsBySymbol`)
- **MT5Service**: Mid-level API for cleanup operations (`CloseOrder`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking
- **protobuf types**: `pb.OrderCloseRequest` for order management

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Symbol` | string | `symbol` argument | Trading instrument (e.g., "EURUSD") |
| `GridSize` | int | `5` | Number of grid levels above AND below current price |
| `GridStep` | float64 | `100` | Distance between levels in points (100pts=10pips for EURUSD) |
| `LotSize` | float64 | `0.01` | Volume (lots) for each order in the grid |
| `MaxPositions` | int | `10` | Maximum concurrent open positions allowed |
| `TakeProfit` | float64 | `0` | Take profit in points. 0 = GridStep is used as TP, **but only when `StopLoss > 0`**: with both 0 the orders are sent with plain `BuyLimit`/`SellLimit`, i.e. without TP |
| `StopLoss` | float64 | `0` | Stop loss in points (0=no SL) |
| `CheckInterval` | time.Duration | `5s` | How often to check and update grid |
| `RebuildOnFill` | bool | `false` | Declared "rebuild entire grid when order fills", but **not read by the current code** |

Defaults are those of `orchestrators.DefaultGridTraderConfig(symbol string)`.

> **Point size and prices:** the grid uses a hard-coded point of `0.00001` (`initializeSymbol()` does not query the symbol), so distances are only correct for 5-digit symbols such as EURUSD. Grid levels are built around the Bid/Ask midpoint and are **not rounded** to the symbol's digits - MT5Sugar sends prices exactly as given.

### Configuration Example

```go
config := orchestrators.GridTraderConfig{
    Symbol:        "EURUSD",
    GridSize:      5,                // 5 levels each side (10 total orders)
    GridStep:      100,              // 100 points = 10 pips spacing
    LotSize:       0.01,             // 0.01 lots per order
    MaxPositions:  10,               // Max 10 concurrent positions
    TakeProfit:    0,                // 0 + StopLoss 0 = orders WITHOUT TP/SL
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

Both commands will launch the **GridTrader** orchestrator with the configuration specified in `main.go -> RunOrchestrator_Grid()` (symbol `cfg.TestSymbol`, `TakeProfit: 0`, `StopLoss: 0` - so the pending orders are placed without TP/SL). The runner runs the grid for 10 minutes (600 s), then calls `Stop()` and prints the final metrics.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Initialize symbol parameters (digits = 5, point = 0.00001, hard-coded)
  |
Delete previously tracked pending orders (none on first build)
  |
Get current price (Bid/Ask midpoint) = grid center
  |
Calculate grid levels:
  Levels ABOVE: center + (1×step), (2×step), ..., (GridSize×step)
  Levels BELOW: center - (1×step), (2×step), ..., (GridSize×step)
  |
FOR each level ABOVE current price:
  | Place SELL LIMIT at level
  | TP/SL only if TakeProfit > 0 or StopLoss > 0
  |   (TP = level - TakeProfit, or level - GridStep when TakeProfit = 0)
  |
FOR each level BELOW current price:
  | Place BUY LIMIT at level
  | TP/SL only if TakeProfit > 0 or StopLoss > 0
  |   (TP = level + TakeProfit, or level + GridStep when TakeProfit = 0)
  |
MONITORING LOOP (every CheckInterval):
  | Get positions count for the symbol
  | If count >= MaxPositions -> skip the rest of this check
  | Get current price
  | If mid price moved > 2×GridStep from center -> rebuild grid
  | Track total P/L of open positions
  |
ON STOP:
  | Delete all tracked pending orders (Service.CloseOrder)
  | Open positions are NOT closed
  |
END
```

### Step-by-step Description

#### 1. Initialization

> Lines `272–286`  


```go
func (g *GridTrader) initializeSymbol() error {
	// Get current price
	priceInfo, err := g.sugar.GetPriceInfo(g.config.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get price: %w", err)
	}
	g.currentPrice = priceInfo.Bid

	// Calculate point value from price info
	// For EURUSD with 5 digits: point = 0.00001
	g.digits = 5 // Default, could query from symbol info
	g.point = 0.00001

	return nil
}
```

**What happens**:

- Checks that a price is available for the symbol (the Bid is stored, but `buildGrid()` replaces it with the Bid/Ask midpoint right after)
- Digits and point are **hard-coded** to 5 and 0.00001 for every symbol - the symbol is not queried
- So the grid distances are only correct for 5-digit symbols such as EURUSD

#### 2. Building the Grid

> Lines `289–331`  


```go
func (g *GridTrader) buildGrid() error {
	// Clean existing orders first
	g.cleanupOrders()

	// Get current price
	priceInfo, err := g.sugar.GetPriceInfo(g.config.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get price: %w", err)
	}
	g.currentPrice = (priceInfo.Bid + priceInfo.Ask) / 2

	// Calculate grid levels
	g.gridLevels = make([]float64, 0)
	gridStepPrice := g.config.GridStep * g.point

	// Build levels above and below current price
	for i := 1; i <= g.config.GridSize; i++ {
		levelAbove := g.currentPrice + float64(i)*gridStepPrice
		levelBelow := g.currentPrice - float64(i)*gridStepPrice
		g.gridLevels = append(g.gridLevels, levelAbove, levelBelow)
	}

	// Place orders at each grid level
	for _, level := range g.gridLevels {
		if level > g.currentPrice {
			// Place SELL LIMIT above current price
			if err := g.placeSellLimit(level); err != nil {
				g.IncrementError(fmt.Sprintf("failed to place sell limit: %v", err))
			}
		} else {
			// Place BUY LIMIT below current price
			if err := g.placeBuyLimit(level); err != nil {
				g.IncrementError(fmt.Sprintf("failed to place buy limit: %v", err))
			}
		}
	}

	g.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("Built grid with %d levels", len(g.gridLevels))
	})

	return nil
}
```

**Critical Details**:

- Grid center = (Bid + Ask) / 2 for balanced placement
- Each level is `GridStep × point` distance from center (levels are not rounded to the symbol's digits)
- SELL LIMIT orders go **ABOVE** current price
- BUY LIMIT orders go **BELOW** current price
- Old grid is cleaned up before rebuilding
- A failed order is only counted as an error; the remaining levels are still placed

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
	// Calculate TP/SL if configured
	tp := 0.0
	sl := 0.0

	if g.config.TakeProfit > 0 {
		tp = price + g.config.TakeProfit*g.point
	} else {
		tp = price + g.config.GridStep*g.point
	}

	if g.config.StopLoss > 0 {
		sl = price - g.config.StopLoss*g.point
	}

	var ticket uint64
	var err error

	if g.config.StopLoss > 0 || g.config.TakeProfit > 0 {
		ticket, err = g.sugar.BuyLimitWithSLTP(g.config.Symbol, g.config.LotSize, price, sl, tp)
	} else {
		ticket, err = g.sugar.BuyLimit(g.config.Symbol, g.config.LotSize, price)
	}

	if err != nil {
		return err
	}

	g.activeOrders = append(g.activeOrders, ticket)
	g.IncrementSuccess()

	// Log with numbering
	orderNum := len(g.activeOrders)
	totalOrders := g.config.GridSize * 2 // GridSize levels above + GridSize levels below

	// Print to console immediately
	fmt.Printf("  [GRID #%d/%d] ✅ Placed BUY LIMIT @ %.5f (ticket #%d)\n",
		orderNum, totalOrders, price, ticket)

	g.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("[GRID #%d/%d] Placed BUY LIMIT @ %.5f (ticket #%d)",
			orderNum, totalOrders, price, ticket)
	})

	return nil
}
```

**Key Points**:

- **TP/SL are only sent when `TakeProfit > 0` or `StopLoss > 0`** (`BuyLimitWithSLTP`); otherwise plain `BuyLimit` is used and the order has **no TP and no SL**, even though a TP value is computed
- When TP is sent: `price + TakeProfit` points, or `price + GridStep` points when `TakeProfit = 0`
- Example (`StopLoss > 0` or `TakeProfit = 100`): BUY fills at 1.09900, TP is at 1.10000
- SL is optional (default 0 = no SL)

**Example BUY order placement** (`TakeProfit = 100`, `StopLoss = 0`):
```
Price: 1.09900 (BUY LIMIT)
TP:    1.10000 (price + 100pts)
SL:    none (StopLoss = 0)
```

#### 4. Placing SELL LIMIT Orders

> Lines `382–427`  


Same structure as `placeBuyLimit`, mirrored (excerpt of the TP/SL part):

```go
	if g.config.TakeProfit > 0 {
		tp = price - g.config.TakeProfit*g.point
	} else {
		tp = price - g.config.GridStep*g.point
	}

	if g.config.StopLoss > 0 {
		sl = price + g.config.StopLoss*g.point
	}

	var ticket uint64
	var err error

	if g.config.StopLoss > 0 || g.config.TakeProfit > 0 {
		ticket, err = g.sugar.SellLimitWithSLTP(g.config.Symbol, g.config.LotSize, price, sl, tp)
	} else {
		ticket, err = g.sugar.SellLimit(g.config.Symbol, g.config.LotSize, price)
	}
```

**Key Points**:

- Again, TP/SL are only sent when `TakeProfit > 0` or `StopLoss > 0`
- When sent, TP is `price - TakeProfit` points (or `price - GridStep` points when `TakeProfit = 0`)
- Example: when SELL fills at 1.10100, TP is at 1.10000

**Example SELL order placement** (`TakeProfit = 100`, `StopLoss = 0`):
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

1. **Get current positions**: Count open positions for the symbol (`GetPositionsBySymbol`)
2. **Check max positions**: If `CurrentPositions >= MaxPositions`, the rest of this check is skipped (existing pending orders stay in place)
3. **Get current price**: Fetch latest Bid/Ask
4. **Check price drift**: If price moved > 2×GridStep from center, rebuild grid
5. **Track P/L**: Sum profit/loss of all open positions
6. **Update metrics**: Store in orchestrator metrics for reporting

**Grid rebuild trigger**:

```go
	// Check if price moved significantly from grid center
	priceDiff := ((priceInfo.Bid + priceInfo.Ask) / 2) - g.currentPrice
	gridStepPrice := g.config.GridStep * g.point

	// Rebuild grid if price moved more than 2 grid steps
	if priceDiff > 2*gridStepPrice || priceDiff < -2*gridStepPrice {
		if err := g.buildGrid(); err != nil {
			g.IncrementError(fmt.Sprintf("failed to rebuild grid: %v", err))
		} else {
			g.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.LastOperation = "Grid rebuilt due to price movement"
			})
		}
	}
```

#### 6. Cleanup on Stop
 
> Lines `504–536`  


```go
func (g *GridTrader) cleanupOrders() {
	// Cancel all active pending orders using Service.CloseOrder
	ctx := context.Background()
	service := g.sugar.GetService()

	totalOrders := len(g.activeOrders)
	if totalOrders > 0 {
		fmt.Printf("\n  🧹 Cleaning up %d pending orders...\n", totalOrders)
	}

	deletedCount := 0
	for i, ticket := range g.activeOrders {
		// CloseOrder can delete pending orders (not just close positions)
		req := &pb.OrderCloseRequest{
			Ticket: ticket,
		}
		_, err := service.CloseOrder(ctx, req)
		if err != nil {
			g.IncrementError(fmt.Sprintf("failed to delete order #%d: %v", ticket, err))
			fmt.Printf("  [CLEANUP %d/%d] ❌ Failed to delete order #%d: %v\n", i+1, totalOrders, ticket, err)
		} else {
			deletedCount++
			fmt.Printf("  [CLEANUP %d/%d] ✅ Deleted order #%d\n", i+1, totalOrders, ticket)
		}
	}

	// Clear tracking list
	g.activeOrders = make([]uint64, 0)

	if totalOrders > 0 {
		fmt.Printf("  ✓ Cleanup complete: %d/%d orders deleted\n\n", deletedCount, totalOrders)
	}
}
```

**What happens**:

- An `OrderClose` request is sent for every ticket in `activeOrders` (MT5 OrderClose deletes a pending order or closes a position)
- The list is not pruned when an order fills, so tickets of already-filled orders are sent too
- `CloseOrder` returns `(uint32 retcode, error)`; the retcode is ignored here - an order counts as deleted whenever no error is returned
- The call uses `context.Background()` (no timeout)
- Open positions are not closed explicitly by `Stop()`

---

## Grid Visualization

### Scenario: GridSize=5, GridStep=100, CurrentPrice=1.10000, TakeProfit=100

(The TP column applies only when `TakeProfit > 0` or `StopLoss > 0`; with the default config the orders have no TP.)

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
Profit:    +100 points = +$1 (for 0.01 lots EURUSD)
Action:    No rebuild (price is not > 2 grid steps from the center);
           the filled BUY LIMIT level stays empty, 9 pending orders remain
Total P/L: +$1
```

**Time 3: Price Rises to 1.10200**
```
Event:     SELL LIMIT #1 @ 1.10100 and SELL LIMIT #2 @ 1.10200 TRIGGERED
Positions: SHORT 0.01 @ 1.10100 (TP 1.10000), SHORT 0.01 @ 1.10200 (TP 1.10100)
Status:    2 positions open, 7 pending orders remain
           (exactly 2 grid steps from the center: not > 2, so no rebuild)
```

**Time 4: Price Falls to 1.10100**
```
Event:     SHORT @ 1.10200 CLOSED by TP @ 1.10100
Profit:    +100 points = +$1
Status:    SHORT @ 1.10100 still open (its TP is 1.10000)
Total P/L: +$2 realized
```

**Time 5: Multiple Oscillations**
```
Filled levels are not re-placed (there is no refill logic); the grid
is only rebuilt - all tracked orders deleted and 10 new ones placed
around the new mid price - when price drifts > 2 grid steps (200pts)
from the center, e.g. at 1.10210.
Each trade closed by TP earns +100pts = +$1 at 0.01 lots.
```

**Summary**:

- Range-bound market (1.09500 - 1.10500): each grid level can be filled once until the next rebuild
- Without grid: Hard to catch manually, many small movements missed
- With grid (and a TP configured): automatic profit on each bounce to a TP

---

## Used MT5Sugar Methods

### 1. BuyLimit / BuyLimitWithSLTP

```go
func (s *MT5Sugar) BuyLimit(symbol string, volume, price float64) (uint64, error)

func (s *MT5Sugar) BuyLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

**Purpose**: Places pending Buy Limit order at specified price.

**Parameters in GridTrader**:

- `symbol`: Trading instrument (from config.Symbol)
- `volume`: `config.LotSize` (e.g., 0.01)
- `price`: Calculated grid level BELOW current price
- `sl`: Stop loss price (0 if disabled)
- `tp`: Take profit price (price + TakeProfit, or price + GridStep when TakeProfit = 0)
- `BuyLimitWithSLTP` is used only when `TakeProfit > 0` or `StopLoss > 0`; otherwise `BuyLimit` (no TP/SL)
- Volume and price are sent exactly as given (no rounding)

**Returns**: Order ticket number or error ("order rejected, code: ..." unless the return code is 10009)

### 2. SellLimit / SellLimitWithSLTP

```go
func (s *MT5Sugar) SellLimit(symbol string, volume, price float64) (uint64, error)

func (s *MT5Sugar) SellLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

**Purpose**: Places pending Sell Limit order at specified price.

**Parameters in GridTrader**:

- `price`: Calculated grid level ABOVE current price
- `tp`: Take profit price (price - TakeProfit, or price - GridStep when TakeProfit = 0); sent only when `TakeProfit > 0` or `StopLoss > 0`

### 3. GetPositionsBySymbol

```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all open positions for a symbol.

**Usage in GridTrader**:

- Called every CheckInterval in monitoring loop
- Counts positions to enforce MaxPositions limit
- Calculates total P/L across all positions

### 4. GetPriceInfo

```go
func (s *MT5Sugar) GetPriceInfo(symbol string) (*PriceInfo, error)
```

**Purpose**: Gets current Bid/Ask prices for a symbol.

**Usage in GridTrader**:

- Gets current price for grid center calculation
- Monitors price drift to trigger grid rebuilds

### 5. Service.CloseOrder (via GetService)

```go
func (s *MT5Service) CloseOrder(ctx context.Context, req *pb.OrderCloseRequest) (uint32, error)
```

Call in GridTrader (`cleanupOrders`):

```go
service := g.sugar.GetService()
_, err := service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: ticket})
```

**Purpose**: Closes position or deletes pending order. Returns the MT5 return code (10009 = done) or an error.

**Usage in GridTrader**:

- Cleanup operation when stopping orchestrator
- Deletes all pending orders in `activeOrders` list
- Ensures no hanging orders after exit

---

## Risk Management

### Maximum Risk Calculation

```
Risk per level = LotSize × StopLoss × PointValue

For EURUSD (PointValue ~$1 per point per 1.0 lot):
- LotSize = 0.01
- StopLoss = 500 points (if configured)
- Risk per level = 0.01 × 500 × $1 = $5

Total maximum risk (if all orders trigger):
- GridSize = 5 -> 10 orders total
- Maximum risk = 10 × $5 = $50
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

   - Limits concurrent exposure: at the limit the check is skipped (no rebuild), but pending orders already placed stay active and can still fill
   - Reduces runaway position accumulation
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


