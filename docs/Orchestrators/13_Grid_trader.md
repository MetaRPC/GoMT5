# GridTrader - Grid Trading Orchestrator

## Description

**🧩 GridTrader** is an automated grid trading orchestrator designed for **RANGE-BOUND and SIDEWAYS markets**. It places a "fishing net" of pending BUY and SELL orders at fixed price intervals, profiting from oscillations as price bounces between levels. Perfect for choppy, non-trending markets where price moves back and forth.

**Strategy Principle**: A symmetrical grid of pending orders is created above AND below current price:
- **SELL LIMIT orders** placed above current price (profit when price rises then falls back)
- **BUY LIMIT orders** placed below current price (profit when price falls then rises back)
- Every order gets a take profit: `TakeProfit` points away, or one grid step (`GridStep`) away when `TakeProfit = 0`; an SL only when `StopLoss > 0`
- Grid is rebuilt around the new mid price when price drifts more than 2 grid steps from the grid center (and, with `RebuildOnFill`, whenever one of its orders fills)
- `MaxPositions` is enforced: the grid never places more orders than the room left under it, and when the limit is reached the remaining pending orders are deleted (the grid is rebuilt once positions are below the limit again)
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
    |  |  |  |
GetSymbolInfo  BuyLimitWithSLTP  SellLimitWithSLTP  GetPositionsBySymbol  GetPriceInfo
(point/digits) (below)           (above)            (monitoring)          (grid center)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetSymbolInfo`, `GetPriceInfo`, `BuyLimitWithSLTP`, `SellLimitWithSLTP`, `GetPositionsBySymbol`)
- **MT5Service**: Mid-level API for fill tracking and cleanup (`GetOpenedTickets`, `CloseOrder`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking
- **protobuf types**: `pb.OrderCloseRequest` for order management

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Symbol` | string | `symbol` argument | Trading instrument (e.g., "EURUSD"; must not be empty) |
| `GridSize` | int | `5` | Number of grid levels above AND below current price (must be > 0) |
| `GridStep` | float64 | `100` | Distance between levels in points (100pts=10pips for EURUSD; must be > 0) |
| `LotSize` | float64 | `0.01` | Volume (lots) for each order in the grid (rounded down to the volume step; must be at least the symbol's minimum volume) |
| `MaxPositions` | int | `10` | Maximum concurrent open positions on `Symbol` (0 = `2 × GridSize`) |
| `TakeProfit` | float64 | `0` | Take profit in points; 0 = one grid step (`GridStep`). Every order is sent with its TP |
| `StopLoss` | float64 | `0` | Stop loss in points (0=no SL) |
| `CheckInterval` | time.Duration | `5s` | How often to check and update grid (0 = 5 s) |
| `RebuildOnFill` | bool | `false` | `true`: rebuild the whole grid around the current price at the first check after one of its orders filled |

Defaults are those of `orchestrators.DefaultGridTraderConfig(symbol string)`. `Start()` validates the configuration: `CheckInterval` 0 becomes 5 s, `MaxPositions` 0 becomes `2 × GridSize`, an empty symbol or a non-positive `GridSize`/`GridStep`/`LotSize` returns an `invalid configuration` error.

> **Point size and prices:** `initializeSymbol()` reads the symbol's point size, digits and volume limits with `sugar.GetSymbolInfo(symbol)`, so `GridStep: 100` is 0.00100 on EURUSD, 0.100 on USDJPY and 1.00 on XAUUSD. Grid levels are built around the Bid/Ask midpoint and rounded to the symbol's digits (MT5Sugar sends prices exactly as given, and an unrounded price can be rejected).

### Configuration Example

```go
config := orchestrators.GridTraderConfig{
    Symbol:        "EURUSD",
    GridSize:      5,                // 5 levels each side (10 total orders)
    GridStep:      100,              // 100 points = 10 pips spacing
    LotSize:       0.01,             // 0.01 lots per order
    MaxPositions:  10,               // Max 10 concurrent positions
    TakeProfit:    0,                // TP one grid step from each order
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

Both commands will launch the **GridTrader** orchestrator with the configuration specified in `main.go -> RunOrchestrator_Grid()` (symbol `cfg.TestSymbol`, `TakeProfit: 0` - so each order's TP is one grid step away, `StopLoss: 0` - no SL). The runner runs the grid for 10 minutes (600 s), then calls `Stop()` - which deletes the grid's pending orders - and prints the final metrics; the runner then disconnects the terminal - it does so however the run ends, also on an error or Ctrl+C. Each orchestrator also cancels its own context on Ctrl+C / SIGTERM, so a program of yours that waits on `GetContext()` can call `Stop()` and disconnect. Positions opened by the grid stay open with their TP.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Start(): validate config (CheckInterval 0 -> 5s, MaxPositions 0 -> 2×GridSize)
  |
Initialize symbol parameters (point, digits, volume limits from GetSymbolInfo)
  |
Delete the grid's tracked orders that are still pending (none on first build)
  |
Get current price (Bid/Ask midpoint) = grid center
  |
Room = MaxPositions - positions already open on the symbol
  |
Calculate grid levels (rounded to the symbol's digits), nearest first:
  center + 1×step, center - 1×step, center + 2×step, center - 2×step, ...
  |
FOR each level while orders placed < Room:
  | ABOVE the center: SELL LIMIT, TP = level - TP distance
  | BELOW the center: BUY LIMIT,  TP = level + TP distance
  |   (TP distance = TakeProfit, or GridStep when TakeProfit = 0; SL only if StopLoss > 0)
  |
MONITORING LOOP (every CheckInterval):
  | Get the positions on the symbol and the account's pending order tickets
  | Tracked orders that are no longer pending have FILLED -> count them, stop tracking them
  | If positions >= MaxPositions -> delete the grid's pending orders, pause, skip the rest
  | Get current price
  | Rebuild the grid when:
  |   - it was paused and positions are below MaxPositions again, or
  |   - RebuildOnFill is true and an order filled, or
  |   - the mid price moved > 2×GridStep from the center
  | Track total P/L of open positions
  |
ON STOP (Stop() or Ctrl+C):
  | Cancel the loop and wait for a check in progress
  | Delete the grid's orders that are still pending (Service.CloseOrder)
  | Open positions are NOT closed (each has its TP)
  |
END
```

### Step-by-step Description

#### 1. Initialization

> Lines `296–314`  


```go
func (g *GridTrader) initializeSymbol() error {
	spec, err := loadSymbolSpec(g.sugar, g.config.Symbol)
	if err != nil {
		return err
	}
	if spec.volume(g.config.LotSize) == 0 {
		return fmt.Errorf("LotSize %.2f is below the minimum volume %.2f of %s", g.config.LotSize, spec.VolumeMin, g.config.Symbol)
	}
	g.spec = spec

	// Get current price
	priceInfo, err := g.sugar.GetPriceInfo(g.config.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get price: %w", err)
	}
	g.currentPrice = priceInfo.Bid

	return nil
}
```

**What happens**:

- `loadSymbolSpec()` (in `orchestrators.go`) reads the symbol's point size, digits, stops level and volume limits with `sugar.GetSymbolInfo(symbol)`
- `LotSize` is checked against the symbol's minimum volume (orders are sent with `LotSize` rounded down to the volume step)
- Checks that a price is available for the symbol (the Bid is stored, but `buildGrid()` replaces it with the Bid/Ask midpoint right after)

#### 2. Building the Grid

> Lines `319–372`  


```go
func (g *GridTrader) buildGrid() error {
	// Clean existing orders first
	g.cleanupOrders()
	g.paused = false

	// Get current price
	priceInfo, err := g.sugar.GetPriceInfo(g.config.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get price: %w", err)
	}
	g.currentPrice = (priceInfo.Bid + priceInfo.Ask) / 2

	// Room left under MaxPositions
	positions, err := g.sugar.GetPositionsBySymbol(g.config.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get positions: %w", err)
	}
	room := g.config.MaxPositions - len(positions)

	// Calculate grid levels (nearest first: +1, -1, +2, -2, ...)
	g.gridLevels = make([]float64, 0)
	gridStepPrice := g.config.GridStep * g.spec.Point

	// Build levels above and below current price
	for i := 1; i <= g.config.GridSize; i++ {
		levelAbove := g.spec.price(g.currentPrice + float64(i)*gridStepPrice)
		levelBelow := g.spec.price(g.currentPrice - float64(i)*gridStepPrice)
		g.gridLevels = append(g.gridLevels, levelAbove, levelBelow)
	}

	// Place orders at each grid level
	for _, level := range g.gridLevels {
		if len(g.activeOrders) >= room {
			break
		}
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
		m.LastOperation = fmt.Sprintf("Built grid with %d levels (%d orders placed)", len(g.gridLevels), len(g.activeOrders))
	})

	return nil
}
```

**Critical Details**:

- Grid center = (Bid + Ask) / 2 for balanced placement
- Each level is `GridStep × point` distance from center, rounded to the symbol's digits (`g.spec.price`)
- SELL LIMIT orders go **ABOVE** current price
- BUY LIMIT orders go **BELOW** current price
- Old grid is cleaned up before rebuilding
- The nearest levels are placed first, and never more orders than `MaxPositions` minus the positions already open - so the grid cannot fill beyond `MaxPositions`
- A failed order is only counted as an error; the remaining levels are still placed

**Example with GridSize=5, GridStep=100, CurrentPrice=1.10000 (EURUSD, point 0.00001)**:

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

Total: 10 pending orders (placed in the order +1, -1, +2, -2, ...)
```

On USDJPY (point 0.001) the same `GridStep: 100` gives levels 0.100 apart, on XAUUSD (point 0.01) 1.00 apart.

#### 3. Placing BUY LIMIT Orders
 
> Lines `383–413`  


```go
func (g *GridTrader) placeBuyLimit(price float64) error {
	// TP always (TakeProfit or one grid step), SL only if configured
	tp := g.spec.price(price + g.takeProfitPoints()*g.spec.Point)
	sl := 0.0
	if g.config.StopLoss > 0 {
		sl = g.spec.price(price - g.config.StopLoss*g.spec.Point)
	}

	ticket, err := g.sugar.BuyLimitWithSLTP(g.config.Symbol, g.spec.volume(g.config.LotSize), price, sl, tp)
	if err != nil {
		return err
	}

	g.activeOrders = append(g.activeOrders, ticket)
	g.IncrementSuccess()
	// ... numbered console line and metrics ...
	return nil
}
```

**Key Points**:

- Every order is sent with `BuyLimitWithSLTP` and a TP: `price + TakeProfit` points, or `price + GridStep` points when `TakeProfit = 0` (`takeProfitPoints()`)
- Example (`TakeProfit = 0`, `GridStep = 100`): BUY fills at 1.09900, TP is at 1.10000 - the next grid level up
- SL is optional (default 0 = no SL)
- TP/SL are rounded to the symbol's digits, the volume to the volume step

**Example BUY order placement** (defaults):
```
Price: 1.09900 (BUY LIMIT)
TP:    1.10000 (price + 100pts)
SL:    none (StopLoss = 0)
```

#### 4. Placing SELL LIMIT Orders

> Lines `416–446`  


Same structure as `placeBuyLimit`, mirrored (excerpt of the TP/SL part):

```go
	// TP always (TakeProfit or one grid step), SL only if configured
	tp := g.spec.price(price - g.takeProfitPoints()*g.spec.Point)
	sl := 0.0
	if g.config.StopLoss > 0 {
		sl = g.spec.price(price + g.config.StopLoss*g.spec.Point)
	}

	ticket, err := g.sugar.SellLimitWithSLTP(g.config.Symbol, g.spec.volume(g.config.LotSize), price, sl, tp)
```

**Key Points**:

- TP is `price - TakeProfit` points (or `price - GridStep` points when `TakeProfit = 0`)
- Example: when SELL fills at 1.10100, TP is at 1.10000

**Example SELL order placement** (defaults):
```
Price: 1.10100 (SELL LIMIT)
TP:    1.10000 (price - 100pts)
SL:    none (StopLoss = 0)
```

#### 5. Monitoring Loop

> Lines `449–461` (loop), `480–572` (`checkAndUpdateGrid`)


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
2. **Track fills**: Tracked orders whose tickets are no longer among the account's pending orders (`Service.GetOpenedTickets`) have filled; they are counted in `TotalTrades` and no longer tracked
3. **Check max positions**: If `CurrentPositions >= MaxPositions`, the grid's remaining pending orders are deleted (so no more can fill) and the grid is paused
4. **Get current price**: Fetch latest Bid/Ask
5. **Rebuild**: after a pause once positions are below `MaxPositions` again, on a fill when `RebuildOnFill` is true, or when the price moved > 2×GridStep from the center
6. **Track P/L**: Sum profit/loss of all open positions
7. **Update metrics**: Store in orchestrator metrics for reporting

**Grid rebuild trigger**:

```go
	// Rebuild grid: after a MaxPositions pause, on a fill (RebuildOnFill),
	// or when the price moved more than 2 grid steps
	reason := ""
	switch {
	case g.paused:
		reason = "positions below MaxPositions again"
	case g.config.RebuildOnFill && filled > 0:
		reason = "order filled (RebuildOnFill)"
	case priceDiff > 2*gridStepPrice || priceDiff < -2*gridStepPrice:
		reason = "price movement"
	}
```

#### 6. Cleanup on Stop
 
> Lines `577–623`  


```go
func (g *GridTrader) cleanupOrders() {
	totalOrders := len(g.activeOrders)
	if totalOrders == 0 {
		return
	}

	pending, err := g.pendingTickets()
	if err != nil {
		g.IncrementError(fmt.Sprintf("failed to get pending orders: %v", err))
		return // keep the list: the next cleanup tries again
	}

	// Cancel the pending orders using Service.CloseOrder
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	service := g.sugar.GetService()

	// ...
	for i, ticket := range g.activeOrders {
		if !pending[ticket] {
			// no longer pending (filled): skipped
			continue
		}
		// CloseOrder deletes a pending order (and would close a position)
		req := &pb.OrderCloseRequest{
			Ticket: ticket,
		}
		code, err := service.CloseOrder(ctx, req)
		if err == nil && code != 10009 {
			err = fmt.Errorf("rejected, code: %d", code)
		}
		// ... counted and printed ...
	}

	// Clear tracking list
	g.activeOrders = make([]uint64, 0)
}
```

**What happens**:

- Only orders that are **still pending** are deleted: in MT5 a filled order's ticket is the ticket of the position it opened, and an `OrderClose` for it would close that position
- `CloseOrder` returns `(uint32 retcode, error)`; anything but 10009 counts as a failed delete
- The calls share one 30-second deadline
- `Stop()` cancels the monitoring loop and waits for a check in progress before it cleans up, so the loop never places orders while they are being deleted
- Open positions are not closed by `Stop()`; each has its TP

---

## Grid Visualization

### Scenario: GridSize=5, GridStep=100, CurrentPrice=1.10000, TakeProfit=0 (= one grid step)

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
Action:    No rebuild (price is not > 2 grid steps from the center, RebuildOnFill
           is false); the filled BUY LIMIT level stays empty, 9 pending orders remain
           (with RebuildOnFill = true the grid would have been rebuilt around
           1.09900 at the first check after the fill)
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
Filled levels are not re-placed one by one; the grid is rebuilt -
its pending orders deleted and up to 10 new ones placed around the
new mid price (fewer when positions are open: never more than
MaxPositions - open positions) - when price drifts > 2 grid steps
(200pts) from the center, e.g. at 1.10210, or after a fill when
RebuildOnFill is true.
Each trade closed by TP earns +100pts = +$1 at 0.01 lots.
```

**Summary**:

- Range-bound market (1.09500 - 1.10500): each grid level can be filled once until the next rebuild
- Without grid: Hard to catch manually, many small movements missed
- With grid: automatic profit on each bounce to a TP

---

## Used MT5Sugar Methods

### 1. BuyLimitWithSLTP

```go
func (s *MT5Sugar) BuyLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

**Purpose**: Places pending Buy Limit order at specified price with SL/TP.

**Parameters in GridTrader**:

- `symbol`: Trading instrument (from config.Symbol)
- `volume`: `config.LotSize` rounded down to the symbol's volume step (e.g., 0.01)
- `price`: Calculated grid level BELOW current price
- `sl`: Stop loss price (0 if disabled)
- `tp`: Take profit price (price + TakeProfit, or price + GridStep when TakeProfit = 0)
- MT5Sugar sends volume and prices exactly as given, so the grid rounds them first (prices to the symbol's digits)

**Returns**: Order ticket number or error ("order rejected, code: ..." unless the return code is 10009)

### 2. SellLimitWithSLTP

```go
func (s *MT5Sugar) SellLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

**Purpose**: Places pending Sell Limit order at specified price with SL/TP.

**Parameters in GridTrader**:

- `price`: Calculated grid level ABOVE current price
- `tp`: Take profit price (price - TakeProfit, or price - GridStep when TakeProfit = 0)

### GetSymbolInfo

```go
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error)
```

Read once at start: `Point` (grid step and TP/SL distances), `Digits` (price rounding), `VolumeMin`/`VolumeStep`/`VolumeMax` (lot size).

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

### 5. Service.GetOpenedTickets and Service.CloseOrder (via GetService)

```go
func (s *MT5Service) GetOpenedTickets(ctx context.Context) ([]int64, []int64, error)
func (s *MT5Service) CloseOrder(ctx context.Context, req *pb.OrderCloseRequest) (uint32, error)
```

Calls in GridTrader (`pendingTickets`, `cleanupOrders`):

```go
_, orderTickets, err := g.sugar.GetService().GetOpenedTickets(ctx)
code, err := service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: ticket})
```

**Purpose**: `GetOpenedTickets` returns the tickets of the open positions and of the pending orders; `CloseOrder` closes a position or deletes a pending order and returns the MT5 return code (10009 = done) or an error.

**Usage in GridTrader**:

- Every check: tracked orders that are no longer pending have filled
- Cleanup (rebuild, `MaxPositions` reached, stop): deletes the grid's orders that are still pending - never a ticket that has become a position
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

   - Limits concurrent exposure: a grid never has more pending orders than `MaxPositions` minus the open positions, and at the limit the remaining pending orders are deleted
   - Reduces runaway position accumulation
   - Counts all positions on the symbol (also ones not opened by the grid)
   - 0 = `GridSize × 2` (every level may fill)

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


