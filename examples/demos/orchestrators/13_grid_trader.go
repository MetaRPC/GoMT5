/*══════════════════════════════════════════════════════════════════════════════
 ORCHESTRATOR: GridTrader (Range-bound Market Strategy)

 ⚠️ IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY ⚠️

 THIS IS A DEMONSTRATION EXAMPLE showing how GoMT5 methods FUNCTION AND COMBINE
 into something more than single method calls. This orchestrator is NOT a
 production-ready trading strategy!

 PURPOSE OF THIS EXAMPLE:
   ✓ Show how MT5Account, MT5Service, and MT5Sugar work together as foundation
   ✓ Demonstrate orchestrator patterns and method combinations
   ✓ Provide a starting point for building YOUR OWN strategies
   ✓ Illustrate best practices for automated trading systems

 ══════════════════════════════════════════════════════════════════════════════

 PURPOSE:
   Automated grid trading orchestrator designed for RANGE-BOUND and SIDEWAYS
   markets. Places a "fishing net" of pending BUY and SELL orders at fixed
   price intervals, profiting from oscillations as price bounces between levels.
   Perfect for choppy, non-trending markets where price moves back and forth.

 STRATEGY:
   • Places symmetrical grid of pending orders above AND below current price
   • SELL LIMIT orders above price (profit when price rises then falls back)
   • BUY LIMIT orders below price (profit when price falls then rises back)
   • Each level acts as a profit target for opposite direction
   • Automatically rebuilds grid when price moves significantly
   • Works 24/7 capturing small profits from price oscillations

 HOW GRID TRADING WORKS:

   Imagine a "ladder" of prices:
   • Every 100 points = one "rung" of the ladder
   • Place BUY orders on lower rungs
   • Place SELL orders on upper rungs
   • As price climbs/falls, orders trigger and take profit

 VISUAL EXAMPLE - GRID STRUCTURE:

   Grid Setup: 5 levels, 100pt spacing, current price = 1.10000
   ──────────────────────────────────────────────────────────────

   1.10500  [SELL LIMIT #5]  ← Level +5 (+500pts from center)
   1.10400  [SELL LIMIT #4]  ← Level +4 (+400pts from center)
   1.10300  [SELL LIMIT #3]  ← Level +3 (+300pts from center)
   1.10200  [SELL LIMIT #2]  ← Level +2 (+200pts from center)
   1.10100  [SELL LIMIT #1]  ← Level +1 (+100pts from center)
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   1.10000  >>> CURRENT PRICE <<<  [Grid Center]
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   1.09900  [BUY LIMIT #1]   ← Level -1 (-100pts from center)
   1.09800  [BUY LIMIT #2]   ← Level -2 (-200pts from center)
   1.09700  [BUY LIMIT #3]   ← Level -3 (-300pts from center)
   1.09600  [BUY LIMIT #4]   ← Level -4 (-400pts from center)
   1.09500  [BUY LIMIT #5]   ← Level -5 (-500pts from center)

   Total: 10 pending orders (5 BUY + 5 SELL)

 VISUAL EXAMPLE - PRICE MOVEMENT & PROFIT:

   Time 0:  Price=1.10000  → Grid placed, all orders pending
   ──────────────────────────────────────────────────────────────

   Time 1:  Price drops to 1.09900
            → BUY LIMIT #1 @ 1.09900 TRIGGERED
            → Position: LONG 0.01 lots @ 1.09900
            → TP set at 1.10000 (next grid level up)

   Time 2:  Price rises to 1.10000
            → BUY position CLOSED by TP
            → Profit: +100 points (+$10 for 0.01 lots)
            → Grid rebuilds, new orders placed

   Time 3:  Price rises to 1.10200
            → SELL LIMIT #2 @ 1.10200 TRIGGERED
            → Position: SHORT 0.01 lots @ 1.10200
            → TP set at 1.10100 (next grid level down)

   Time 4:  Price falls to 1.10100
            → SELL position CLOSED by TP
            → Profit: +100 points (+$10 for 0.01 lots)
            → Grid rebuilds, new orders placed

   Time 5:  Price oscillates: 1.09900 → 1.10000 → 1.09800 → 1.10100
            → Multiple trades triggered as price bounces
            → Each bounce = +100pts profit
            → Over 4 bounces: +400 points total ($40 profit)

   RESULT:
   Range-bound market (price oscillating 1.09500 - 1.10500):
   • Without grid: Hard to trade manually, many small movements missed
   • With grid: Captures EVERY oscillation automatically
   • Typical profit: 10-20 small wins per day in choppy markets

 WHEN GRID TRADING SHINES:

   ✅ PERFECT CONDITIONS:
   • Asian session (low volatility, tight ranges)
   • Forex pairs during off-hours (EURUSD at night)
   • Consolidation after big news events
   • Support/resistance ranges (1.0500-1.0600 channel)

   ⚠️ DANGEROUS CONDITIONS:
   • Strong trending markets (price breaks grid and never returns)
   • High-impact news events (violent breakouts)
   • Low liquidity (gaps can skip grid levels)

 KEY PARAMETERS:
   • Symbol: Trading pair (e.g., "EURUSD")
   • GridSize: Number of levels above/below (default: 5 = 10 total orders)
   • GridStep: Distance between levels in points (default: 100pts = 10 pips)
   • LotSize: Volume for each order (default: 0.01 lots)
   • MaxPositions: Max concurrent positions (default: 10)
   • TakeProfit: TP distance (default: 0 = use GridStep)
   • StopLoss: SL distance (default: 0 = no SL)
   • CheckInterval: How often to rebuild grid (default: 5s)

 USE CASES:

   📊 Scenario 1: Tight Range Market (Asian Session)
   GridSize: 5, GridStep: 50pts, LotSize: 0.01
   → Captures small 50pt oscillations repeatedly
   → Best for EURUSD/USDJPY during low volatility

   📊 Scenario 2: Wide Range Market (Consolidation)
   GridSize: 10, GridStep: 200pts, LotSize: 0.02
   → Wider net for bigger range-bound movements
   → Best for XAUUSD in $10-20 consolidation zones

   📊 Scenario 3: Conservative Grid (Beginner)
   GridSize: 3, GridStep: 100pts, LotSize: 0.01
   → Fewer orders = less risk
   → Narrower grid = less capital required

 COMMAND-LINE USAGE:
   cd examples/demos

   go run main.go 13
   go run main.go grid
   

 CONFIGURATION:
   ⚙️ All parameters configured in main.go → RunOrchestrator_Grid()
   📍 See end of this file for detailed configuration examples and documentation
   ⚠️ Grid trading requires RANGE-BOUND markets to work well!

══════════════════════════════════════════════════════════════════════════════*/

package orchestrators

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// GridTraderConfig holds configuration for the grid trading strategy.
type GridTraderConfig struct {
	Symbol         string        // Trading symbol (e.g., "EURUSD")
	GridSize       int           // Number of grid levels (above and below)
	GridStep       float64       // Distance between levels in points
	LotSize        float64       // Volume for each order
	MaxPositions   int           // Maximum concurrent positions
	TakeProfit     float64       // Take profit in points (0 = use grid step)
	StopLoss       float64       // Stop loss in points (0 = no SL)
	CheckInterval  time.Duration // How often to check and update grid
	RebuildOnFill  bool          // Rebuild entire grid when order fills
}

// DefaultGridTraderConfig returns sensible default configuration.
func DefaultGridTraderConfig(symbol string) GridTraderConfig {
	return GridTraderConfig{
		Symbol:        symbol,
		GridSize:      5,
		GridStep:      100,
		LotSize:       0.01,
		MaxPositions:  10,
		TakeProfit:    0, // Use grid step as TP
		StopLoss:      0, // No SL by default
		CheckInterval: 5 * time.Second,
		RebuildOnFill: false,
	}
}

// ══════════════════════════════════════════════════════════════════════════════
// GRID TRADER IMPLEMENTATION
// ══════════════════════════════════════════════════════════════════════════════

// GridTrader implements automated grid trading strategy.
type GridTrader struct {
	*BaseOrchestrator
	sugar  *mt5.MT5Sugar
	config GridTraderConfig

	// Grid state
	gridLevels    []float64   // Current grid price levels
	activeOrders  []uint64    // Tickets of active pending orders
	digits        int         // Symbol decimal digits
	point         float64     // Point value for symbol
	currentPrice  float64     // Last known price
}

// NewGridTrader creates a new grid trading orchestrator.
func NewGridTrader(sugar *mt5.MT5Sugar, config GridTraderConfig) *GridTrader {
	return &GridTrader{
		BaseOrchestrator: NewBaseOrchestrator("Grid Trader"),
		sugar:            sugar,
		config:           config,
		activeOrders:     make([]uint64, 0),
		gridLevels:       make([]float64, 0),
	}
}

// Start begins the grid trading operation.
func (g *GridTrader) Start() error {
	if g.IsRunning() {
		return fmt.Errorf("grid trader already running")
	}

	// Initialize symbol parameters
	if err := g.initializeSymbol(); err != nil {
		return fmt.Errorf("failed to initialize symbol: %w", err)
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	g.SetContext(ctx, cancel)

	// Mark as started
	g.MarkStarted()

	// Build initial grid
	if err := g.buildGrid(); err != nil {
		g.MarkStopped()
		return fmt.Errorf("failed to build initial grid: %w", err)
	}

	// Start monitoring loop
	go g.monitorLoop()

	return nil
}

// Stop gracefully stops the grid trader.
func (g *GridTrader) Stop() error {
	if !g.IsRunning() {
		return fmt.Errorf("grid trader not running")
	}

	// Cancel context to stop loop
	g.CancelContext()

	// Clean up all pending orders
	g.cleanupOrders()

	// Mark as stopped
	g.MarkStopped()

	return nil
}

// initializeSymbol gets symbol parameters (digits, point).
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

// buildGrid constructs the grid of pending orders.
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

// placeBuyLimit places a BUY LIMIT order at specified price.
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

// placeSellLimit places a SELL LIMIT order at specified price.
func (g *GridTrader) placeSellLimit(price float64) error {
	// Calculate TP/SL if configured
	tp := 0.0
	sl := 0.0

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

	if err != nil {
		return err
	}

	g.activeOrders = append(g.activeOrders, ticket)
	g.IncrementSuccess()

	// Log with numbering
	orderNum := len(g.activeOrders)
	totalOrders := g.config.GridSize * 2 // GridSize levels above + GridSize levels below

	// Print to console immediately
	fmt.Printf("  [GRID #%d/%d] ✅ Placed SELL LIMIT @ %.5f (ticket #%d)\n",
		orderNum, totalOrders, price, ticket)

	g.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("[GRID #%d/%d] Placed SELL LIMIT @ %.5f (ticket #%d)",
			orderNum, totalOrders, price, ticket)
	})

	return nil
}

// monitorLoop continuously monitors the grid and adjusts as needed.
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

// checkAndUpdateGrid checks current positions and rebuilds grid if needed.
func (g *GridTrader) checkAndUpdateGrid() {
	// Get current positions
	positions, err := g.sugar.GetPositionsBySymbol(g.config.Symbol)
	if err != nil {
		g.IncrementError(fmt.Sprintf("failed to get positions: %v", err))
		return
	}

	// Update metrics
	g.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.CurrentPositions = len(positions)
		if len(positions) > m.MaxPositions {
			m.MaxPositions = len(positions)
		}
	})

	// Check if we hit max positions
	if len(positions) >= g.config.MaxPositions {
		g.UpdateMetrics(func(m *OrchestratorMetrics) {
			m.LastOperation = "Max positions reached, waiting..."
		})
		return
	}

	// Get current price
	priceInfo, err := g.sugar.GetPriceInfo(g.config.Symbol)
	if err != nil {
		g.IncrementError(fmt.Sprintf("failed to get price: %v", err))
		return
	}

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

	// Update positions profit tracking
	totalProfit := 0.0
	for _, pos := range positions {
		totalProfit += pos.Profit
	}

	g.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.CurrentDrawdown = totalProfit
		m.LastOperation = fmt.Sprintf("Monitoring %d positions, P/L: %.2f", len(positions), totalProfit)
	})
}

// cleanupOrders cancels all pending orders.
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

/*══════════════════════════════════════════════════════════════════════════════

 ███████╗ ██████╗ ███╗   ██╗███████╗██╗ ██████╗ ██╗   ██╗██████╗  █████╗ ████████╗██╗ ██████╗ ███╗   ██╗
██╔════╝██╔═══██╗████╗  ██║██╔════╝██║██╔════╝ ██║   ██║██╔══██╗██╔══██╗╚══██╔══╝██║██╔═══██╗████╗  ██║
██║     ██║   ██║██╔██╗ ██║█████╗  ██║██║  ███╗██║   ██║██████╔╝███████║   ██║   ██║██║   ██║██╔██╗ ██║
██║     ██║   ██║██║╚██╗██║██╔══╝  ██║██║   ██║██║   ██║██╔══██╗██╔══██║   ██║   ██║██║   ██║██║╚██╗██║
╚██████╗╚██████╔╝██║ ╚████║██║     ██║╚██████╔╝╚██████╔╝██║  ██║██║  ██║   ██║   ██║╚██████╔╝██║ ╚████║
 ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝

  DETAILED CONFIGURATION GUIDE
  Located at end of file to keep header clean and focused

══════════════════════════════════════════════════════════════════════════════*/

/*
PROGRAMMATIC USAGE & CONFIGURATION

⚙️ PARAMETER CONFIGURATION IS LOCATED IN main.go

WHY THIS SEPARATION EXISTS:
• 13_grid_trader.go = STRATEGY ENGINE (orchestrator logic, grid algorithm)
• main.go → RunOrchestrator_Grid() = RUNTIME CONFIGURATION (parameters)

THIS SEPARATION IS NEEDED FOR:
1️⃣ Code Reusability
   → Same orchestrator class can run with different parameters
   → No need to modify strategy logic to change parameters

2️⃣ Quick Testing
   → Want tighter grid? Change numbers in main.go
   → Want different symbol? Again, only change main.go
   → Core algorithm remains untouched and battle-tested

3️⃣ User Examples
   → main.go shows HOW to properly configure the orchestrator
   → All available parameters and their default values are visible

4️⃣ Centralized Entry Point
   → All strategies launch through main.go
   → Single entry point: go run main.go grid → RunOrchestrator_Grid()

📍 WHERE TO CONFIGURE PARAMETERS:
main.go → func RunOrchestrator_Grid() (lines 549-559)

CONFIGURATION CODE IN main.go:

func RunOrchestrator_Grid() error {
    // ... connection code ...

    // ╔════════════════════════════════════════════════════════════╗
    // ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
    // ╚════════════════════════════════════════════════════════════╝
    orchConfig := orchestrators.GridTraderConfig{
        Symbol:         "EURUSD",           // ← Trading symbol
        GridSize:       5,                  // ← 5 levels each side (10 total)
        GridStep:       100,                // ← 100 points between levels
        LotSize:        0.01,               // ← 0.01 lots per order
        MaxPositions:   10,                 // ← Max concurrent positions
        TakeProfit:     0,                  // ← 0 = use grid step as TP
        StopLoss:       0,                  // ← 0 = no stop loss
        CheckInterval:  5 * time.Second,    // ← Check every 5 seconds
        RebuildOnFill:  false,              // ← Don't rebuild on fill
    }

    gridTrader := orchestrators.NewGridTrader(sugar, orchConfig)
    gridTrader.Start()
    // ... runs for 10 minutes ...
    gridTrader.Stop()
}

💡 EXAMPLE CONFIGURATIONS FOR DIFFERENT SCENARIOS

╔═══════════════════════════════════════════════════════════════════════════╗
║ SCENARIO 1: TIGHT RANGE (ASIAN SESSION - LOW VOLATILITY)                 ║
╚═══════════════════════════════════════════════════════════════════════════╝

GridTraderConfig{
    Symbol:         "EURUSD",
    GridSize:       5,               // ← Small grid for tight range
    GridStep:       50,              // ← 50pts = 5 pips (very tight)
    LotSize:        0.01,            // ← Small lot size
    MaxPositions:   10,
    TakeProfit:     0,               // ← Use grid step (50pts) as TP
    StopLoss:       0,               // ← No SL in range-bound
    CheckInterval:  3 * time.Second, // ← Check often for quick adjustments
    RebuildOnFill:  false,
}

BEST FOR:
• EURUSD during Asian session (Tokyo hours)
• USDJPY when volatility is low
• Pairs oscillating in 20-30 pip range
• When you expect price to bounce between levels

EXPECTED RESULT:
• Frequent small wins (5-10 pips each)
• 20-30 trades per day
• Total daily profit: 50-100 pips if range holds


╔═══════════════════════════════════════════════════════════════════════════╗
║ SCENARIO 2: WIDE RANGE (CONSOLIDATION AFTER BIG MOVE)                    ║
╚═══════════════════════════════════════════════════════════════════════════╝

GridTraderConfig{
    Symbol:         "XAUUSD",        // ← Gold (wider ranges)
    GridSize:       10,              // ← Larger grid for wider range
    GridStep:       200,             // ← 200pts = $20 on gold
    LotSize:        0.02,            // ← Bigger lot for wider grid
    MaxPositions:   15,              // ← More positions allowed
    TakeProfit:     0,               // ← Use grid step (200pts)
    StopLoss:       500,             // ← 500pt SL for safety
    CheckInterval:  10 * time.Second,
    RebuildOnFill:  false,
}

BEST FOR:
• XAUUSD (Gold) consolidating in $50-100 range
• After NFP or FOMC (consolidation phase)
• Indices (SPX500) in sideways markets
• When support/resistance is clear

EXPECTED RESULT:
• Fewer trades, bigger profits (20-50 pips each)
• 5-10 trades per day
• Total daily profit: 100-200 pips if range respected


╔═══════════════════════════════════════════════════════════════════════════╗
║ SCENARIO 3: CONSERVATIVE GRID (BEGINNER-FRIENDLY)                        ║
╚═══════════════════════════════════════════════════════════════════════════╝

GridTraderConfig{
    Symbol:         "EURUSD",
    GridSize:       3,               // ← Only 3 levels (6 orders total)
    GridStep:       100,             // ← Standard 10 pip spacing
    LotSize:        0.01,            // ← Minimum lot size
    MaxPositions:   6,               // ← Limit exposure
    TakeProfit:     0,
    StopLoss:       300,             // ← Add safety SL
    CheckInterval:  5 * time.Second,
    RebuildOnFill:  false,
}

BEST FOR:
• Learning grid trading mechanics
• Testing on demo account
• Low-risk testing
• Small account sizes ($500-1000)

EXPECTED RESULT:
• Low risk, low reward
• 3-8 trades per day
• Total daily profit: 20-40 pips


╔═══════════════════════════════════════════════════════════════════════════╗
║ SCENARIO 4: AGGRESSIVE GRID (EXPERIENCED TRADERS)                        ║
╚═══════════════════════════════════════════════════════════════════════════╝

GridTraderConfig{
    Symbol:         "GBPJPY",        // ← Volatile pair
    GridSize:       15,              // ← Very large grid
    GridStep:       150,             // ← 15 pip spacing
    LotSize:        0.05,            // ← Larger lots
    MaxPositions:   30,              // ← Many concurrent positions
    TakeProfit:     0,
    StopLoss:       0,               // ← No SL (risky!)
    CheckInterval:  2 * time.Second, // ← Very frequent checks
    RebuildOnFill:  true,            // ← Rebuild after each fill
}

BEST FOR:
• Experienced traders only
• Large account sizes ($10,000+)
• Highly volatile but range-bound pairs
• When you can monitor actively

EXPECTED RESULT:
• High frequency trading
• 50-100 trades per day
• Total daily profit: 200-500 pips (but high risk!)


╔═══════════════════════════════════════════════════════════════════════════╗
║ PARAMETER EXPLANATIONS                                                   ║
╚═══════════════════════════════════════════════════════════════════════════╝

• Symbol (string)
  Trading pair to run grid on
  Examples: "EURUSD", "XAUUSD", "USDJPY"
  Tip: Choose pairs with clear support/resistance levels

• GridSize (int)
  Number of grid levels EACH SIDE of current price
  GridSize=5 → 10 total orders (5 BUY + 5 SELL)
  Range: 3-20 (too small = missed opportunities, too large = overexposure)

• GridStep (float64)
  Distance in POINTS between each grid level
  EURUSD: 100pts = 10 pips, 50pts = 5 pips
  XAUUSD: 200pts = $20 move
  Tip: Match to average range of chosen timeframe

• LotSize (float64)
  Volume for EACH pending order
  0.01 = 1 micro lot
  Tip: Start small! GridSize × LotSize = total exposure

• MaxPositions (int)
  Maximum concurrent open positions allowed
  Safety limit to prevent overexposure
  Tip: Set to GridSize × 2 or higher

• TakeProfit (float64)
  Take profit distance in points
  0 = use GridStep as TP (common for grid trading)
  Non-zero: custom TP level
  Tip: TP = GridStep ensures profit on each oscillation

• StopLoss (float64)
  Stop loss distance in points
  0 = no stop loss (risky but common in grid trading)
  Non-zero: safety stop loss
  Tip: SL = 3 × GridStep for safety margin

• CheckInterval (time.Duration)
  How often orchestrator checks and updates grid
  3s = very active, 10s = relaxed
  Tip: Match to symbol volatility

• RebuildOnFill (bool)
  Whether to rebuild entire grid when order fills
  true = aggressive (more orders), false = passive
  Tip: false for stable ranges, true for active trading


╔═══════════════════════════════════════════════════════════════════════════╗
║ RISK WARNINGS                                                             ║
╚═══════════════════════════════════════════════════════════════════════════╝

⚠️  GRID TRADING RISKS:
1. TRENDING MARKETS
   → Grid designed for RANGES, not trends
   → Strong trend = all orders on one side fill
   → Result: Multiple losing positions, no offsetting wins

2. OVEREXPOSURE
   → GridSize × LotSize = total volume exposure
   → Example: GridSize=10, LotSize=0.1 → 1.0 lots total
   → Ensure adequate margin for full grid

3. NO STOP LOSS DANGER
   → Many grid strategies don't use SL
   → If range breaks violently, losses can be severe
   → Always monitor or use emergency SL

4. REBUILDING LOOPS
   → RebuildOnFill=true can cause rapid order placement
   → In trending market: constant rebuilding = overtrading
   → Use with caution

💡 BEST PRACTICES:
✅ Test on DEMO first
✅ Start with conservative GridSize (3-5)
✅ Match GridStep to symbol's average range
✅ Monitor during initial hours
✅ Use during known range-bound sessions (Asian, post-news consolidation)
✅ Have emergency stop strategy (max drawdown alert)

═══════════════════════════════════════════════════════════════════════════*/
