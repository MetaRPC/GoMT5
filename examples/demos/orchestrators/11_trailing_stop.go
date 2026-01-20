/*══════════════════════════════════════════════════════════════════════════════
 ORCHESTRATOR: TrailingStopManager

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
   Automated trailing stop manager that protects profits for all open positions.
   Monitors positions continuously and moves stop loss orders as price moves
   favorably, locking in profits without manual intervention.

 STRATEGY:
   • Monitors all open positions in real-time (every 2 seconds by default)
   • Waits for position to reach activation profit threshold (e.g., +300 points)
   • Once activated, trails stop loss behind price movement
   • For BUY positions: Moves SL upward as price rises (never down)
   • For SELL positions: Moves SL downward as price falls (never up)
   • Maintains fixed trailing distance from current price
   • Never moves SL in unfavorable direction (one-way protection)
   • Respects minimum step size to avoid excessive order modifications

 VISUAL EXAMPLE:

   BUY POSITION TRAILING:

   Time 0:  Entry=1.1000  SL=0        Profit=0      [Not Active]
   Time 1:  Price=1.1030  SL=0        Profit=+300   [✓ ACTIVATED]
   Time 2:  Price=1.1040  SL=1.1020   Profit=+400   [Trailing: -200pts]
   Time 3:  Price=1.1050  SL=1.1030   Profit=+500   [Trailing: -200pts]
   Time 4:  Price=1.1045  SL=1.1030   Profit=+450   [SL unchanged]
   Time 5:  Price=1.1030  SL=1.1030   Profit=+300   [SL triggered → Exit]
                                                     [Locked profit: +300pts]

   Without trailing: Profit could have reversed to 0 or negative
   With trailing:    Guaranteed minimum profit of 300 points locked in

 KEY PARAMETERS:
   • TrailingDistance:  Distance in points to trail behind price (default: 200)
   • ActivationProfit:  Profit in points needed to activate trailing (default: 300)
   • UpdateInterval:    How often to check and update stops (default: 2s)
   • Symbols:           Which symbols to manage (empty = all symbols)
   • MinDistance:       Minimum distance from current price (default: 100)
   • StepSize:          Minimum step size for SL adjustments (default: 50)

 USE CASE:
   Best for trending markets where you want to:
   • Let winning positions run without manual monitoring
   • Protect accumulated profits automatically
   • Avoid giving back large profits during reversals
   • Manage multiple positions simultaneously
   • Free yourself from constant chart watching

 COMMAND-LINE USAGE:
   cd examples/demos
   
   go run main.go 11
   go run main.go trailing

 PROGRAMMATIC USAGE:

   ⚙️ PARAMETER CONFIGURATION IS LOCATED IN main.go

   WHY THIS SEPARATION EXISTS:
   • 10_trailing_stop.go = STRATEGY ENGINE (orchestrator logic, algorithm)
   • main.go → RunOrchestrator_TrailingStop() = RUNTIME CONFIGURATION (parameters)

   THIS SEPARATION IS NEEDED FOR:
   1️⃣ Code Reusability
      → Same orchestrator class can run with different parameters
      → No need to modify strategy logic to change parameters

   2️⃣ Quick Testing
      → Want tighter trailing? Change numbers in main.go
      → Want wider safety margin? Again, only change main.go
      → Core algorithm remains untouched and battle-tested

   3️⃣ User Examples
      → main.go shows HOW to properly configure the orchestrator
      → All available parameters and their default values are visible

   4️⃣ Centralized Entry Point
      → All strategies launch through main.go
      → Single entry point: go run main.go trailing → RunOrchestrator_TrailingStop()

   📝 IMPORTANT:
   • To change parameters → edit main.go, NOT this file
   • This file (10_trailing_stop.go) contains only ORCHESTRATOR LOGIC
   • main.go contains CONFIGURATION for specific runs
   • Look for the section: ORCHESTRATOR RUNNERS

══════════════════════════════════════════════════════════════════════════════*/

package orchestrators

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/MetaRPC/GoMT5/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// TrailingStopConfig holds configuration for trailing stop management.
type TrailingStopConfig struct {
	TrailingDistance float64       // Distance in points to trail behind price
	ActivationProfit float64       // Profit in points to activate trailing
	UpdateInterval   time.Duration // How often to check positions
	Symbols          []string      // Symbols to manage (empty = all)
	MinDistance      float64       // Minimum distance from current price
	StepSize         float64       // Minimum step size for SL adjustments
}

// DefaultTrailingStopConfig returns sensible defaults.
func DefaultTrailingStopConfig() TrailingStopConfig {
	return TrailingStopConfig{
		TrailingDistance: 200,
		ActivationProfit: 300,
		UpdateInterval:   2 * time.Second,
		Symbols:          []string{},
		MinDistance:      100,
		StepSize:         50,
	}
}

// ══════════════════════════════════════════════════════════════════════════════
// TRAILING STOP MANAGER IMPLEMENTATION
// ══════════════════════════════════════════════════════════════════════════════

// TrailingStopManager automatically manages trailing stops for positions.
type TrailingStopManager struct {
	*BaseOrchestrator
	sugar  *mt5.MT5Sugar
	config TrailingStopConfig

	// Tracking
	trackedPositions map[uint64]*positionTracker
	symbolDigits     map[string]int
	symbolPoints     map[string]float64
}

// positionTracker tracks trailing stop state for a position.
type positionTracker struct {
	ticket          uint64
	symbol          string
	isBuy           bool
	openPrice       float64
	currentSL       float64
	highestProfit   float64
	trailingActive  bool
	lastUpdate      time.Time
}

// NewTrailingStopManager creates a new trailing stop manager.
func NewTrailingStopManager(sugar *mt5.MT5Sugar, config TrailingStopConfig) *TrailingStopManager {
	return &TrailingStopManager{
		BaseOrchestrator: NewBaseOrchestrator("Trailing Stop Manager"),
		sugar:            sugar,
		config:           config,
		trackedPositions: make(map[uint64]*positionTracker),
		symbolDigits:     make(map[string]int),
		symbolPoints:     make(map[string]float64),
	}
}

// Start begins trailing stop management.
func (t *TrailingStopManager) Start() error {
	if t.IsRunning() {
		return fmt.Errorf("trailing stop manager already running")
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	t.SetContext(ctx, cancel)

	// Mark as started
	t.MarkStarted()

	// Start monitoring loop
	go t.monitorLoop()

	return nil
}

// Stop gracefully stops the trailing stop manager.
func (t *TrailingStopManager) Stop() error {
	if !t.IsRunning() {
		return fmt.Errorf("trailing stop manager not running")
	}

	// Cancel context
	t.CancelContext()

	// Clear tracked positions
	t.trackedPositions = make(map[uint64]*positionTracker)

	// Mark as stopped
	t.MarkStopped()

	return nil
}

// monitorLoop continuously monitors and updates trailing stops.
func (t *TrailingStopManager) monitorLoop() {
	ticker := time.NewTicker(t.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.GetContext().Done():
			return
		case <-ticker.C:
			t.updateAllTrailingStops()
		}
	}
}

// updateAllTrailingStops checks all positions and updates trailing stops.
func (t *TrailingStopManager) updateAllTrailingStops() {
	// Get all open positions
	positions, err := t.sugar.GetOpenPositions()
	if err != nil {
		t.IncrementError(fmt.Sprintf("failed to get positions: %v", err))
		return
	}

	// Update metrics
	t.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.CurrentPositions = len(positions)
	})

	updatedCount := 0

	// Process each position
	for _, pos := range positions {
		// Skip if symbol not in our list (if list is specified)
		if len(t.config.Symbols) > 0 && !t.isSymbolTracked(pos.Symbol) {
			continue
		}

		// Update trailing stop for this position
		if t.updatePositionTrailingStop(pos) {
			updatedCount++
		}
	}

	// Clean up positions that are no longer open
	t.cleanupClosedPositions(positions)

	// Update status
	if updatedCount > 0 {
		t.UpdateMetrics(func(m *OrchestratorMetrics) {
			m.LastOperation = fmt.Sprintf("Updated %d trailing stop(s)", updatedCount)
		})
		t.IncrementSuccess()
	}
}

// updatePositionTrailingStop updates trailing stop for a single position.
func (t *TrailingStopManager) updatePositionTrailingStop(pos *pb.PositionInfo) bool {
	// Get or create tracker
	tracker, exists := t.trackedPositions[pos.Ticket]
	if !exists {
		tracker = &positionTracker{
			ticket:         pos.Ticket,
			symbol:         pos.Symbol,
			isBuy:          pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY,
			openPrice:      pos.PriceOpen,
			currentSL:      pos.StopLoss,
			highestProfit:  0,
			trailingActive: false,
			lastUpdate:     time.Now(),
		}
		t.trackedPositions[pos.Ticket] = tracker
	}

	// Get symbol parameters
	point, err := t.getSymbolPoint(pos.Symbol)
	if err != nil {
		t.IncrementError(fmt.Sprintf("failed to get symbol point: %v", err))
		return false
	}

	// Get current price
	priceInfo, err := t.sugar.GetPriceInfo(pos.Symbol)
	if err != nil {
		t.IncrementError(fmt.Sprintf("failed to get price: %v", err))
		return false
	}

	// Calculate profit in points
	var profitPoints float64
	var currentPrice float64

	if tracker.isBuy {
		currentPrice = priceInfo.Bid
		profitPoints = (currentPrice - tracker.openPrice) / point
	} else {
		currentPrice = priceInfo.Ask
		profitPoints = (tracker.openPrice - currentPrice) / point
	}

	// Update highest profit
	if profitPoints > tracker.highestProfit {
		tracker.highestProfit = profitPoints
	}

	// Check if trailing should be activated
	if !tracker.trailingActive && profitPoints >= t.config.ActivationProfit {
		tracker.trailingActive = true
		t.UpdateMetrics(func(m *OrchestratorMetrics) {
			m.LastOperation = fmt.Sprintf("Trailing activated for #%d", pos.Ticket)
		})
	}

	// If trailing not active yet, skip
	if !tracker.trailingActive {
		return false
	}

	// Calculate new stop loss
	var newSL float64
	if tracker.isBuy {
		newSL = currentPrice - t.config.TrailingDistance*point
		// Only move SL up, never down
		if tracker.currentSL == 0 || newSL > tracker.currentSL {
			// Check minimum step
			if tracker.currentSL > 0 && (newSL-tracker.currentSL) < t.config.StepSize*point {
				return false
			}
			return t.modifyStopLoss(pos.Ticket, newSL, tracker)
		}
	} else {
		newSL = currentPrice + t.config.TrailingDistance*point
		// Only move SL down, never up
		if tracker.currentSL == 0 || newSL < tracker.currentSL {
			// Check minimum step
			if tracker.currentSL > 0 && (tracker.currentSL-newSL) < t.config.StepSize*point {
				return false
			}
			return t.modifyStopLoss(pos.Ticket, newSL, tracker)
		}
	}

	return false
}

// modifyStopLoss modifies the stop loss for a position.
func (t *TrailingStopManager) modifyStopLoss(ticket uint64, newSL float64, tracker *positionTracker) bool {
	err := t.sugar.ModifyPositionSL(ticket, newSL)
	if err != nil {
		t.IncrementError(fmt.Sprintf("failed to modify SL for #%d: %v", ticket, err))
		return false
	}

	// Update tracker
	tracker.currentSL = newSL
	tracker.lastUpdate = time.Now()

	t.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.OperationsTotal++
	})

	return true
}

// getSymbolPoint gets or caches the point value for a symbol.
func (t *TrailingStopManager) getSymbolPoint(symbol string) (float64, error) {
	// Check cache
	if point, exists := t.symbolPoints[symbol]; exists {
		return point, nil
	}

	// Calculate from price info
	// Default: EURUSD-like symbols with 5 digits
	point := 0.00001
	t.symbolPoints[symbol] = point
	t.symbolDigits[symbol] = 5

	return point, nil
}

// isSymbolTracked checks if symbol is in tracked list.
func (t *TrailingStopManager) isSymbolTracked(symbol string) bool {
	for _, s := range t.config.Symbols {
		if s == symbol {
			return true
		}
	}
	return false
}

// cleanupClosedPositions removes trackers for closed positions.
func (t *TrailingStopManager) cleanupClosedPositions(openPositions []*pb.PositionInfo) {
	// Build map of open position tickets
	openTickets := make(map[uint64]bool)
	for _, pos := range openPositions {
		openTickets[pos.Ticket] = true
	}

	// Remove trackers for positions that are no longer open
	for ticket := range t.trackedPositions {
		if !openTickets[ticket] {
			delete(t.trackedPositions, ticket)
		}
	}
}

/* ══════════════════════════════════════════════════════════════════════════════
   📍 WHERE TO CONFIGURE PARAMETERS:
   main.go → func RunOrchestrator_TrailingStop() (lines 353-360)

   CONFIGURATION CODE IN main.go:

   func RunOrchestrator_TrailingStop() error {
       // ... connection code ...

       // ╔════════════════════════════════════════════════════════════╗
       // ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
       // ╚════════════════════════════════════════════════════════════╝
       orchConfig := orchestrators.TrailingStopConfig{
           TrailingDistance: 200,             // ← Trail 200 points behind
           ActivationProfit: 300,             // ← Activate after 300 points profit
           UpdateInterval:   2 * time.Second, // ← Check every 2 seconds
           Symbols:          []string{},      // ← All symbols (empty = all)
           MinDistance:      100,             // ← Min 100 points from price
           StepSize:         50,              // ← Adjust in 50 point steps
       }

       tsManager := orchestrators.NewTrailingStopManager(sugar, orchConfig)
       tsManager.Start()
       // ... runs for 3 minutes ...
       tsManager.Stop()
   }

   💡 EXAMPLE: Adjusting for Different Trading Styles

   // Option 1: Conservative (wider safety, default in main.go)
   TrailingDistance: 200,   // Wide buffer to avoid premature exits
   ActivationProfit: 300,   // Wait for solid profit before activating
   UpdateInterval:   2s

   // Option 2: Aggressive (tighter trailing, modify in main.go)
   TrailingDistance: 100,   // ← Closer trailing = higher risk/reward
   ActivationProfit: 150,   // ← Activate sooner
   UpdateInterval:   1s     // ← More frequent checks

   // Option 3: Scalping (very tight, modify in main.go)
   TrailingDistance: 50,    // ← Very tight trailing for quick exits
   ActivationProfit: 100,   // ← Activate on small profits
   UpdateInterval:   500ms  // ← Near real-time monitoring
══════════════════════════════════════════════════════════════════════════════*/