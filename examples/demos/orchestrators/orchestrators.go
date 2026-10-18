/*══════════════════════════════════════════════════════════════════════════════
 FILE: orchestrators.go - COMMON TYPES AND INTERFACES FOR ORCHESTRATORS

 PURPOSE:
   Defines common interfaces, types, and utilities shared by all orchestrators.
   Orchestrators are autonomous trading systems that execute complex strategies.

 ORCHESTRATOR PATTERN:
   Each orchestrator is an independent component that:
   - Runs autonomously in background
   - Monitors market conditions
   - Executes trading logic
   - Tracks metrics and performance
   - Can be started/stopped/monitored

 AVAILABLE ORCHESTRATORS:
   1. Grid Trader          - Automated grid trading strategy
   2. Trailing Stop Manager - Dynamic trailing stop for all positions
   3. Risk Manager         - Account-level risk control and limits
   4. Portfolio Rebalancer - Multi-symbol portfolio balancing
   5. Position Scaler      - Pyramiding and position averaging

══════════════════════════════════════════════════════════════════════════════*/

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                    📚 BEGINNER'S GUIDE - SIMPLE EXPLANATION                 ║
╚══════════════════════════════════════════════════════════════════════════════╝

WHAT IS AN ORCHESTRATOR?
────────────────────────
Think of orchestrators as TRADING ASSISTANTS that work for you automatically.
Instead of you sitting at the computer watching charts 24/7 and clicking buttons,
an orchestrator does specific trading tasks automatically in the background.

Real-World Analogy:
  Imagine you hire 5 different assistants for your trading business:
  • Assistant #1 (Grid Trader): Places buy/sell orders in a grid pattern
  • Assistant #2 (Trailing Stop): Moves your stop losses as price moves in profit
  • Assistant #3 (Risk Manager): Watches your account and prevents over-trading
  • Assistant #4 (Portfolio Rebalancer): Keeps your multi-symbol portfolio balanced
  • Assistant #5 (Position Scaler): Adds to winning positions systematically

  Each assistant has ONE specific job and does it automatically without your input.

HOW DO ORCHESTRATORS WORK?
───────────────────────────
1. YOU START IT: Run the orchestrator (e.g., "go run main.go 13")
2. IT RUNS IN BACKGROUND: Monitors your account and market continuously
3. IT EXECUTES LOGIC: When conditions are met, it takes action (places orders,
   modifies positions, closes trades, etc.)
4. IT TRACKS PERFORMANCE: Records what it did, success/failure, profit/loss
5. YOU STOP IT: When done, orchestrator stops gracefully

Example Flow:
  ┌─────────────────────────────────────────────────────────────────┐
  │  START → Monitor → Condition Met? → Execute Action → Repeat     │
  │            ↑                                           ↓        │
  │            └───────────────── Loop ────────────────────┘        │
  └─────────────────────────────────────────────────────────────────┘

WHY USE ORCHESTRATORS?
──────────────────────
✅ AUTOMATION: No need to watch charts 24/7
✅ CONSISTENCY: Executes strategy exactly as programmed, no emotions
✅ SPEED: Reacts to market changes faster than manual trading
✅ MULTI-TASKING: Run multiple orchestrators simultaneously
✅ TRACKING: Automatically tracks all metrics and performance

IMPORTANT - ORCHESTRATORS VS TRADING BOTS:
───────────────────────────────────────────
Some orchestrators are FULL TRADING BOTS (open/close positions):
  • Grid Trader (opens buy/sell orders in grid)
  • Position Scaler (adds to positions)

Others are MANAGEMENT TOOLS (manage existing positions):
  • Trailing Stop Manager (adjusts stops on your positions)
  • Risk Manager (monitors account limits)
  • Portfolio Rebalancer (balances existing portfolio)

⚠️ IMPORTANT: Check each orchestrator's documentation to understand what it does!

AVAILABLE ORCHESTRATORS - SIMPLE DESCRIPTIONS:
───────────────────────────────────────────────

1. 📊 GRID TRADER (#13)
   What: Opens buy/sell orders at fixed price intervals (grid pattern)
   When: Works in ranging markets (price bounces between levels)
   Example: EURUSD between 1.0800-1.1000, grid every 20 pips
   Risk: Can open MANY positions if price trends strongly

2. 🎯 TRAILING STOP MANAGER (#11)
   What: Automatically moves stop losses as price moves in your favor
   When: You have open positions and want to lock in profits
   Example: Position +100 pips profit, stop moves up automatically
   Risk: Low - only modifies existing positions

3. 🛡️ RISK MANAGER (#14)
   What: Monitors account and prevents over-trading/over-risking
   When: Running 24/7 to protect your account from excessive losses
   Example: Stops trading if daily loss exceeds $500
   Risk: Low - protective tool

4. ⚖️ PORTFOLIO REBALANCER (#15)
   What: Maintains target allocation percentages across multiple symbols
   When: You trade multiple pairs and want balanced exposure
   Example: Keep 25% each in EURUSD, GBPUSD, USDJPY, AUDUSD
   Risk: Medium - adjusts position sizes, does NOT open initial positions

5. 📈 POSITION SCALER (#12)
   What: Adds to winning positions (pyramiding) or averages losing positions
   When: You want to scale into trends or average down
   Example: Position +50 pips? Add another 0.1 lot
   Risk: Medium/High - increases position size

HOW TO RUN ORCHESTRATORS:
──────────────────────────

Step 1: Navigate to examples/demos/ directory
Step 2: Run one of these commands:

   go run main.go 13              ← Run Grid Trader by number
   go run main.go grid            ← Run Grid Trader by short name
   go run main.go trailing        ← Run Trailing Stop Manager
   go run main.go risk            ← Run Risk Manager
   go run main.go rebalancer      ← Run Portfolio Rebalancer
   go run main.go scaler          ← Run Position Scaler

Step 3: Watch console output for status updates
Step 4: Press Ctrl+C to stop (graceful shutdown)

SAFETY TIPS FOR BEGINNERS:
───────────────────────────
1. ⚠️ ALWAYS TEST ON DEMO ACCOUNT FIRST! Never run on live account initially.

2. 📖 READ THE DOCUMENTATION: Each orchestrator has detailed docs in its file.
   Check the configuration section to understand what each parameter does.

3. 💰 START SMALL: Use small position sizes (0.01 lot) when testing.

4. 📊 UNDERSTAND THE STRATEGY: Know what the orchestrator does BEFORE running it.
   Some orchestrators open new positions, others only manage existing ones.

5. 🛑 KNOW HOW TO STOP: Ctrl+C stops the orchestrator gracefully.

6. 📈 MONITOR INITIALLY: Watch the console output for first 10-15 minutes to
   ensure orchestrator behaves as expected.

7. 💸 SET RISK LIMITS: Configure max position size, max positions, stop loss
   in the orchestrator's configuration.

CONFIGURATION:
──────────────
Each orchestrator has its own configuration section in main.go.
Look for comments like:

   // ═══════════════════════════════════════════════════════════════
   // ORCHESTRATOR #13: GRID TRADER
   // ═══════════════════════════════════════════════════════════════

Inside each section, you'll find all configurable parameters with explanations.

NEXT STEPS:
───────────
1. Choose an orchestrator from the list above
2. Open its file (e.g., 13_grid_trader.go) and read the full documentation
3. Open main.go and review the configuration for that orchestrator
4. Adjust configuration if needed (symbol, lot size, parameters)
5. Run on DEMO account first: go run main.go [number]
6. Monitor and learn how it works
7. Once comfortable, you can run on live account (at your own risk)

REMEMBER:
─────────
Orchestrators are POWERFUL TOOLS but also carry RISK. Always understand what
an orchestrator does before running it. Test thoroughly on demo accounts.
Never risk more than you can afford to lose.

Happy Trading! 🚀

╔══════════════════════════════════════════════════════════════════════════════╗
║  For detailed documentation on each orchestrator, see the individual files:  ║
║  • 12_position_scaler.go                                                     ║
║  • 13_grid_trader.go                                                         ║
║  • 14_risk_manager.go                                                        ║
║  • 15_portfolio_rebalancer.go                                                ║
╚══════════════════════════════════════════════════════════════════════════════╝

*/

package orchestrators

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CORE INTERFACES
// ══════════════════════════════════════════════════════════════════════════════

// Orchestrator is the base interface that all orchestrators must implement.
// It provides lifecycle management and status monitoring.
type Orchestrator interface {
	// Start begins the orchestrator's operation.
	// Returns error if already running or startup fails.
	Start() error

	// Stop gracefully stops the orchestrator.
	// Waits for current operations to complete.
	Stop() error

	// GetStatus returns current operational status.
	GetStatus() OrchestratorStatus

	// GetMetrics returns performance metrics.
	GetMetrics() OrchestratorMetrics

	// IsRunning returns true if orchestrator is currently active.
	IsRunning() bool
}

// ══════════════════════════════════════════════════════════════════════════════
// STATUS AND METRICS
// ══════════════════════════════════════════════════════════════════════════════

// OrchestratorStatus represents current operational state of an orchestrator.
type OrchestratorStatus struct {
	Name         string        // Orchestrator name
	IsRunning    bool          // Currently running
	StartTime    time.Time     // When started
	LastUpdate   time.Time     // Last activity timestamp
	ErrorCount   int           // Total errors encountered
	SuccessCount int           // Successful operations
	LastError    string        // Last error message (if any)
	Uptime       time.Duration // Time since start
}

// OrchestratorMetrics tracks performance and trading statistics.
type OrchestratorMetrics struct {
	// Trading Stats
	TotalTrades      int     // Total number of trades executed
	WinningTrades    int     // Number of profitable trades
	LosingTrades     int     // Number of losing trades
	BreakevenTrades  int     // Trades closed at breakeven

	// Financial Metrics
	TotalProfit      float64 // Total realized profit
	TotalLoss        float64 // Total realized loss
	NetProfit        float64 // Net profit (profit - loss)
	MaxDrawdown      float64 // Maximum drawdown experienced
	CurrentDrawdown  float64 // Current drawdown

	// Position Stats
	CurrentPositions int     // Currently open positions
	MaxPositions     int     // Maximum concurrent positions
	AvgPositionSize  float64 // Average position size

	// Performance Metrics
	WinRate          float64 // Win rate percentage
	ProfitFactor     float64 // Gross profit / gross loss
	AvgWin           float64 // Average winning trade
	AvgLoss          float64 // Average losing trade

	// Operational Metrics
	OperationsTotal  int     // Total operations performed
	OperationsFailed int     // Failed operations
	LastOperation    string  // Description of last operation
}

// CalculateWinRate calculates the win rate percentage.
func (m *OrchestratorMetrics) CalculateWinRate() {
	if m.TotalTrades > 0 {
		m.WinRate = (float64(m.WinningTrades) / float64(m.TotalTrades)) * 100
	}
}

// CalculateProfitFactor calculates the profit factor.
func (m *OrchestratorMetrics) CalculateProfitFactor() {
	if m.TotalLoss > 0 {
		m.ProfitFactor = m.TotalProfit / m.TotalLoss
	} else if m.TotalProfit > 0 {
		m.ProfitFactor = 999.99 // Infinite profit factor
	}
}

// CalculateAverages calculates average win and loss.
func (m *OrchestratorMetrics) CalculateAverages() {
	if m.WinningTrades > 0 {
		m.AvgWin = m.TotalProfit / float64(m.WinningTrades)
	}
	if m.LosingTrades > 0 {
		m.AvgLoss = m.TotalLoss / float64(m.LosingTrades)
	}
}

// UpdateMetrics recalculates all derived metrics.
func (m *OrchestratorMetrics) UpdateMetrics() {
	m.NetProfit = m.TotalProfit - m.TotalLoss
	m.CalculateWinRate()
	m.CalculateProfitFactor()
	m.CalculateAverages()
}

// ══════════════════════════════════════════════════════════════════════════════
// BASE ORCHESTRATOR IMPLEMENTATION
// ══════════════════════════════════════════════════════════════════════════════

// BaseOrchestrator provides common functionality for all orchestrators.
// Embed this in your orchestrator implementation.
type BaseOrchestrator struct {
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	running     bool
	status      OrchestratorStatus
	metrics     OrchestratorMetrics
	updateChan  chan struct{}
	loops       sync.WaitGroup // goroutines started with goLoop; stopLoops waits for them
	stopTime    time.Time      // when MarkStopped ran (Uptime after Stop)
}

// NewBaseOrchestrator creates a new base orchestrator with given name.
func NewBaseOrchestrator(name string) *BaseOrchestrator {
	return &BaseOrchestrator{
		status: OrchestratorStatus{
			Name:      name,
			IsRunning: false,
		},
		metrics: OrchestratorMetrics{},
		updateChan: make(chan struct{}, 1),
	}
}

// IsRunning returns true if orchestrator is running.
func (b *BaseOrchestrator) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

// GetStatus returns current status (thread-safe).
func (b *BaseOrchestrator) GetStatus() OrchestratorStatus {
	b.mu.RLock()
	defer b.mu.RUnlock()

	status := b.status
	if b.running {
		status.Uptime = time.Since(b.status.StartTime)
	} else if !b.stopTime.IsZero() {
		status.Uptime = b.stopTime.Sub(b.status.StartTime) // final uptime after Stop()
	}
	return status
}

// GetMetrics returns current metrics (thread-safe).
func (b *BaseOrchestrator) GetMetrics() OrchestratorMetrics {
	b.mu.RLock()
	defer b.mu.RUnlock()

	metrics := b.metrics
	metrics.UpdateMetrics()
	return metrics
}

// UpdateStatus updates the status (thread-safe).
func (b *BaseOrchestrator) UpdateStatus(updateFunc func(*OrchestratorStatus)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	updateFunc(&b.status)
	b.status.LastUpdate = time.Now()

	// Signal update
	select {
	case b.updateChan <- struct{}{}:
	default:
	}
}

// UpdateMetrics updates the metrics (thread-safe).
func (b *BaseOrchestrator) UpdateMetrics(updateFunc func(*OrchestratorMetrics)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	updateFunc(&b.metrics)
	b.metrics.UpdateMetrics()
}

// IncrementError increments error counter.
func (b *BaseOrchestrator) IncrementError(errMsg string) {
	b.UpdateStatus(func(s *OrchestratorStatus) {
		s.ErrorCount++
		s.LastError = errMsg
	})
}

// IncrementSuccess increments success counter.
func (b *BaseOrchestrator) IncrementSuccess() {
	b.UpdateStatus(func(s *OrchestratorStatus) {
		s.SuccessCount++
	})
}

// MarkStarted marks orchestrator as started.
func (b *BaseOrchestrator) MarkStarted() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.running = true
	b.status.IsRunning = true
	b.status.StartTime = time.Now()
	b.status.LastUpdate = time.Now()
}

// MarkStopped marks orchestrator as stopped.
func (b *BaseOrchestrator) MarkStopped() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.running = false
	b.status.IsRunning = false
	b.stopTime = time.Now()
}

// GetContext returns the orchestrator's context.
func (b *BaseOrchestrator) GetContext() context.Context {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ctx
}

// SetContext sets a new context for the orchestrator.
func (b *BaseOrchestrator) SetContext(ctx context.Context, cancel context.CancelFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ctx = ctx
	b.cancel = cancel
}

// CancelContext cancels the orchestrator's context.
func (b *BaseOrchestrator) CancelContext() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cancel != nil {
		b.cancel()
	}
}

// begin creates the orchestrator's context, marks it started and makes Ctrl+C (or SIGTERM) cancel that context.
// A program that waits on GetContext() (like the demos in main.go) then calls Stop() and disconnects,
// instead of the process being killed with the terminal still connected and the work half done.
// The handler is removed after the first interrupt, so a second Ctrl+C ends the program the default way.
func (b *BaseOrchestrator) begin() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	b.SetContext(ctx, cancel)
	b.MarkStarted()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	b.goLoop(func() {
		defer signal.Stop(interrupt)
		select {
		case <-interrupt:
			fmt.Printf("\n  🛑 Interrupt received - stopping %s...\n", b.GetStatus().Name)
			b.CancelContext()
		case <-ctx.Done():
		}
	})
	return ctx
}

// goLoop runs fn in a goroutine that stopLoops waits for.
func (b *BaseOrchestrator) goLoop(fn func()) {
	b.loops.Add(1)
	go func() {
		defer b.loops.Done()
		fn()
	}()
}

// stopLoops cancels the context and waits until every goroutine started with goLoop has returned,
// so Stop() never cleans up while a check is still placing or closing orders.
func (b *BaseOrchestrator) stopLoops() {
	b.CancelContext()
	b.loops.Wait()
}

// ══════════════════════════════════════════════════════════════════════════════
// SYMBOL SPECIFICATION
// ══════════════════════════════════════════════════════════════════════════════

// symbolSpec holds what the strategies need from a symbol's specification. All distances in the
// configurations are in POINTS, and a point is symbol-specific: 0.00001 for EURUSD, 0.001 for USDJPY,
// 0.01 for XAUUSD - so it is read from the symbol, never assumed.
type symbolSpec struct {
	Point      float64 // SYMBOL_POINT
	Digits     int     // price digits (prices are rounded to them before they are sent)
	VolumeMin  float64 // minimum volume
	VolumeMax  float64 // maximum volume
	VolumeStep float64 // volume step
	StopLevel  int32   // SYMBOL_TRADE_STOPS_LEVEL in points (minimum SL/TP distance from the price)
	Contract   float64 // contract size (units of 1 lot)
}

// loadSymbolSpec reads a symbol's specification with sugar.GetSymbolInfo.
func loadSymbolSpec(sugar *mt5.MT5Sugar, symbol string) (symbolSpec, error) {
	info, err := sugar.GetSymbolInfo(symbol)
	if err != nil {
		return symbolSpec{}, fmt.Errorf("symbol %s: %w", symbol, err)
	}
	if info.Point <= 0 {
		return symbolSpec{}, fmt.Errorf("symbol %s: invalid point size %v", symbol, info.Point)
	}
	return symbolSpec{
		Point:      info.Point,
		Digits:     int(info.Digits),
		VolumeMin:  info.VolumeMin,
		VolumeMax:  info.VolumeMax,
		VolumeStep: info.VolumeStep,
		StopLevel:  info.StopLevel,
		Contract:   info.ContractSize,
	}, nil
}

// price rounds a price to the symbol's digits.
func (s symbolSpec) price(p float64) float64 {
	return RoundToDigits(p, s.Digits)
}

// volume rounds a volume DOWN to the symbol's volume step and caps it at the maximum.
// It returns 0 when the result is below the symbol's minimum volume (the trade is not possible).
func (s symbolSpec) volume(v float64) float64 {
	step := s.VolumeStep
	if step <= 0 {
		step = 0.01
	}
	if s.VolumeMax > 0 && v > s.VolumeMax {
		v = s.VolumeMax
	}
	v = math.Floor(v/step+1e-9) * step
	v = math.Round(v*1e8) / 1e8 // remove float noise (0.30000000000000004)
	if v <= 0 || v < s.VolumeMin {
		return 0
	}
	return v
}

// ══════════════════════════════════════════════════════════════════════════════
// UTILITY FUNCTIONS
// ══════════════════════════════════════════════════════════════════════════════

// FormatDuration formats a duration in human-readable form.
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// CalculatePoints converts price difference to points for a symbol.
// For EURUSD with 5 digits: 1 point = 0.00001
func CalculatePoints(priceDiff float64, digits int) float64 {
	multiplier := 1.0
	for i := 0; i < digits; i++ {
		multiplier *= 10.0
	}
	return priceDiff * multiplier
}

// RoundToDigits rounds a price to specified number of digits.
func RoundToDigits(price float64, digits int) float64 {
	multiplier := 1.0
	for i := 0; i < digits; i++ {
		multiplier *= 10.0
	}
	return math.Round(price*multiplier) / multiplier
}
