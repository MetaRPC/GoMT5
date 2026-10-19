/*══════════════════════════════════════════════════════════════════════════════
 ORCHESTRATOR: RiskManager (Account Protection System)

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
   Account guardian that monitors risk metrics 24/7 and automatically intervenes
   when danger limits are exceeded. PREVENTS ACCOUNT BLOW-UPS by enforcing
   strict risk rules: drawdown limits, daily loss caps, margin monitoring.
   The "safety net" that saves you from catastrophic losses.

 STRATEGY:
   • Continuously monitors: Equity, Balance, Drawdown, Margin Level, Daily P/L
   • Enforces hard limits on all risk parameters
   • EMERGENCY CLOSE ALL when critical thresholds breached
   • Blocks trading when daily limits hit (prevents revenge trading)
   • Logs all risk events for post-analysis

 KEY PROTECTIONS:
   1️⃣ Drawdown Protection   2️⃣ Daily Loss Limit   3️⃣ Margin Safety
   4️⃣ Position Limits       5️⃣ Daily Profit Target

 COMMAND-LINE USAGE:
   cd examples/demos

   go run main.go 14
   go run main.go risk

 CONFIGURATION:
   ⚙️ All parameters configured in main.go → RunOrchestrator_RiskManager()
   📍 See end of this file for detailed configuration examples and documentation

══════════════════════════════════════════════════════════════════════════════*/

package orchestrators

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// RiskManagerConfig holds risk management parameters.
// Every limit is optional: 0 switches it off.
type RiskManagerConfig struct {
	// Drawdown Limits (from the peak equity)
	MaxDrawdownPercent  float64 // Maximum account drawdown percentage
	MaxDrawdownAbsolute float64 // Maximum absolute drawdown amount

	// Daily Limits (today's P/L = equity - equity at the start of the day)
	DailyLossLimit   float64 // Maximum daily loss allowed
	DailyProfitTarget float64 // Stop trading after hitting this profit

	// Margin Limits
	MinMarginLevel float64 // Minimum margin level percentage (equity / margin × 100)
	MaxMarginUsed  float64 // Maximum margin usage percentage (margin / equity × 100)

	// Position Limits
	MaxOpenPositions  int     // Maximum concurrent positions
	MaxSymbolExposure int     // Maximum positions per symbol
	MaxPositionSize   float64 // Maximum lot size per position

	// Operational
	CheckInterval      time.Duration // How often to check risk (0 = 5s)
	EnableAutoClose    bool          // Automatically close positions
	EnableTradeBlocking bool         // Block new trades when limits hit
}

// DefaultRiskManagerConfig returns conservative default settings.
func DefaultRiskManagerConfig() RiskManagerConfig {
	return RiskManagerConfig{
		MaxDrawdownPercent:  10.0,
		MaxDrawdownAbsolute: 1000.0,
		DailyLossLimit:      500.0,
		DailyProfitTarget:   1000.0,
		MinMarginLevel:      150.0,
		MaxMarginUsed:       80.0,
		MaxOpenPositions:    20,
		MaxSymbolExposure:   5,
		MaxPositionSize:     1.0,
		CheckInterval:       5 * time.Second,
		EnableAutoClose:     true,
		EnableTradeBlocking: true,
	}
}

// validate applies the default CheckInterval and rejects negative limits (0 = limit off).
func (c *RiskManagerConfig) validate() error {
	if c.CheckInterval <= 0 {
		c.CheckInterval = 5 * time.Second
	}
	if c.MaxDrawdownPercent < 0 || c.MaxDrawdownAbsolute < 0 || c.DailyLossLimit < 0 || c.DailyProfitTarget < 0 ||
		c.MinMarginLevel < 0 || c.MaxMarginUsed < 0 || c.MaxOpenPositions < 0 || c.MaxSymbolExposure < 0 ||
		c.MaxPositionSize < 0 {
		return fmt.Errorf("limits must not be negative (0 switches a limit off)")
	}
	return nil
}

// ══════════════════════════════════════════════════════════════════════════════
// RISK MANAGER IMPLEMENTATION
// ══════════════════════════════════════════════════════════════════════════════

// RiskManager monitors and enforces account-level risk limits.
type RiskManager struct {
	*BaseOrchestrator
	sugar  *mt5.MT5Sugar
	config RiskManagerConfig

	// State - written by the monitoring goroutine under mu, read by the getters under mu
	mu                sync.Mutex
	startingBalance   float64
	peakBalance       float64 // peak EQUITY since start (reset to the equity after an emergency close)
	dailyStartBalance float64 // equity at the start of the day
	todayProfit       float64 // equity - dailyStartBalance
	tradingBlocked    bool
	lastResetDate     time.Time

	// Used by the monitoring goroutine only
	dailyBlocked   bool                  // daily loss limit or profit target hit: blocked until the next day
	marginBlocked  bool                  // margin usage above MaxMarginUsed: blocked while it is
	dailyLossHit   bool                  // daily loss limit handled today
	dailyTargetHit bool                  // daily profit target handled today
	allowedTickets map[uint64]bool       // positions open when trading got blocked (nil = not blocked)
	activeBreaches map[string]bool       // limits currently exceeded (an event is logged once per breach)
	symbolSpecs    map[string]symbolSpec // volume step per symbol (MaxPositionSize reduction)

	// Risk Events
	riskEvents []RiskEvent
}

// errAccountNotLoaded: the equity read 0. That is not a real reading: a terminal that has just been
// started again by the server (after a TERMINAL_INSTANCE_NOT_FOUND, for example) answers 0 for the
// account values - and an empty position list - until it has loaded the account. Taken as real, it
// was a 100% drawdown: emergency close, trading blocked for the day, and the positions that showed
// up a few seconds later were closed as "opened while trading is blocked" (seen on a test server).
var errAccountNotLoaded = errors.New("account data not loaded yet (equity 0)")

// RiskEvent records a risk limit breach.
type RiskEvent struct {
	Timestamp   time.Time
	EventType   string
	Severity    string
	Description string
	Value       float64
	Limit       float64
	ActionTaken string
}

// NewRiskManager creates a new risk management orchestrator.
func NewRiskManager(sugar *mt5.MT5Sugar, config RiskManagerConfig) *RiskManager {
	return &RiskManager{
		BaseOrchestrator: NewBaseOrchestrator("Risk Manager"),
		sugar:            sugar,
		config:           config,
		riskEvents:       make([]RiskEvent, 0),
		lastResetDate:    time.Now(),
		activeBreaches:   make(map[string]bool),
		symbolSpecs:      make(map[string]symbolSpec),
	}
}

// Start validates the configuration and begins risk monitoring.
func (r *RiskManager) Start() error {
	if r.IsRunning() {
		return fmt.Errorf("risk manager already running")
	}

	// Validate configuration (applies the default CheckInterval)
	if err := r.config.validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize baseline values
	if err := r.initialize(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Create context, mark as started, stop on Ctrl+C
	r.begin()

	// Start monitoring loop
	r.goLoop(r.monitorLoop)

	return nil
}

// Stop gracefully stops risk monitoring.
// It waits until a check that is in progress has finished.
func (r *RiskManager) Stop() error {
	if !r.IsRunning() {
		return fmt.Errorf("risk manager not running")
	}

	// Cancel context and wait for the monitoring loop
	r.stopLoops()

	// Mark as stopped
	r.MarkStopped()

	return nil
}

// initialize sets up baseline values.
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
	if equity <= 0 {
		// A 0 baseline would turn the first real reading into a huge daily profit
		return errAccountNotLoaded
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.startingBalance = balance
	r.peakBalance = equity
	r.dailyStartBalance = equity
	r.todayProfit = 0
	r.tradingBlocked = false
	r.lastResetDate = time.Now()
	r.dailyBlocked, r.marginBlocked, r.dailyLossHit, r.dailyTargetHit = false, false, false, false
	r.allowedTickets = nil
	r.activeBreaches = make(map[string]bool)

	return nil
}

// monitorLoop continuously monitors risk metrics.
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

// checkRiskLimits checks all risk parameters and takes action if needed.
func (r *RiskManager) checkRiskLimits() {
	// Reset daily tracking if new day
	r.checkDailyReset()

	// Get current account state
	equity, err := r.sugar.GetEquity()
	if err != nil {
		r.IncrementError(fmt.Sprintf("failed to get equity: %v", err))
		return
	}
	if equity <= 0 {
		// Not a real reading (see errAccountNotLoaded): no limit is checked on it
		r.IncrementError(errAccountNotLoaded.Error() + " - check skipped")
		return
	}
	margin, marginErr := r.sugar.GetMargin()
	marginLevel, levelErr := r.sugar.GetMarginLevel()

	// Update peak equity, drawdown and today's P/L
	r.mu.Lock()
	if equity > r.peakBalance {
		r.peakBalance = equity
	}
	currentDrawdown := r.peakBalance - equity
	drawdownPercent := 0.0
	if r.peakBalance > 0 {
		drawdownPercent = (currentDrawdown / r.peakBalance) * 100
	}
	r.todayProfit = equity - r.dailyStartBalance
	r.mu.Unlock()

	// Update metrics (MaxDrawdown = the largest drawdown % seen)
	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.CurrentDrawdown = currentDrawdown
		m.MaxDrawdown = math.Max(m.MaxDrawdown, drawdownPercent)
	})

	// Check each risk limit
	r.checkDrawdownLimit(drawdownPercent, currentDrawdown)
	r.checkDailyLimits()
	if marginErr == nil && levelErr == nil {
		r.checkMarginLimits(marginLevel, margin, equity)
	} else {
		r.IncrementError(fmt.Sprintf("failed to get margin: %v %v", marginErr, levelErr))
	}

	// Positions after the actions above
	positions, err := r.sugar.GetOpenPositions()
	if err != nil {
		r.IncrementError(fmt.Sprintf("failed to get positions: %v", err))
	} else {
		r.checkPositionLimits(positions)
		r.enforceTradeBlock(positions)
	}

	// Update status
	r.IncrementSuccess()
	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("Monitoring: Equity=%.2f, DD=%.1f%%, Margin=%.0f%%",
			equity, drawdownPercent, marginLevel)
	})
}

// breach reports whether a limit has just been exceeded: true on the first check that finds
// exceeded=true, false while it stays exceeded, and the breach clears when exceeded is false.
// This way one breach logs one event instead of one per check.
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

// checkDrawdownLimit monitors maximum drawdown (a limit of 0 is off).
// After an emergency close the peak is reset to the current equity, so the next breach
// is measured from there instead of repeating on every check.
func (r *RiskManager) checkDrawdownLimit(drawdownPercent, drawdownAbsolute float64) {
	action := "Logged"
	if r.config.EnableAutoClose {
		action = "Closed all positions"
	}
	var reasons []string

	// Check percentage drawdown
	if r.breach("MAX_DRAWDOWN_PERCENT", r.config.MaxDrawdownPercent > 0 && drawdownPercent >= r.config.MaxDrawdownPercent) {
		r.logRiskEvent("MAX_DRAWDOWN_PERCENT", "CRITICAL",
			fmt.Sprintf("Drawdown %.1f%% exceeds limit %.1f%%", drawdownPercent, r.config.MaxDrawdownPercent),
			drawdownPercent, r.config.MaxDrawdownPercent, action)
		reasons = append(reasons, "Maximum drawdown exceeded")
	}

	// Check absolute drawdown
	if r.breach("MAX_DRAWDOWN_ABSOLUTE", r.config.MaxDrawdownAbsolute > 0 && drawdownAbsolute >= r.config.MaxDrawdownAbsolute) {
		r.logRiskEvent("MAX_DRAWDOWN_ABSOLUTE", "CRITICAL",
			fmt.Sprintf("Drawdown $%.2f exceeds limit $%.2f", drawdownAbsolute, r.config.MaxDrawdownAbsolute),
			drawdownAbsolute, r.config.MaxDrawdownAbsolute, action)
		reasons = append(reasons, "Maximum absolute drawdown exceeded")
	}

	if len(reasons) > 0 && r.config.EnableAutoClose {
		r.closeAllPositionsEmergency(strings.Join(reasons, ", "))
		r.resetPeak()
	}
}

// resetPeak makes the current equity the new peak and clears the drawdown breaches.
func (r *RiskManager) resetPeak() {
	equity, err := r.sugar.GetEquity()
	if err != nil {
		r.IncrementError(fmt.Sprintf("failed to get equity: %v", err))
		return
	}
	if equity <= 0 {
		r.IncrementError(errAccountNotLoaded.Error() + " - peak not reset")
		return
	}
	r.mu.Lock()
	r.peakBalance = equity
	r.mu.Unlock()
	delete(r.activeBreaches, "MAX_DRAWDOWN_PERCENT")
	delete(r.activeBreaches, "MAX_DRAWDOWN_ABSOLUTE")
}

// checkDailyLimits monitors daily profit/loss limits (once per day each; a limit of 0 is off).
func (r *RiskManager) checkDailyLimits() {
	r.mu.Lock()
	todayProfit := r.todayProfit
	r.mu.Unlock()

	// Check daily loss limit
	if !r.dailyLossHit && r.config.DailyLossLimit > 0 && -todayProfit >= r.config.DailyLossLimit {
		r.dailyLossHit = true
		action := r.actionText(r.config.EnableAutoClose, "closed all positions", r.config.EnableTradeBlocking, "trading blocked until tomorrow")
		r.logRiskEvent("DAILY_LOSS_LIMIT", "CRITICAL",
			fmt.Sprintf("Daily loss $%.2f exceeds limit $%.2f", -todayProfit, r.config.DailyLossLimit),
			-todayProfit, r.config.DailyLossLimit, action)

		if r.config.EnableAutoClose {
			r.closeAllPositionsEmergency("Daily loss limit exceeded")
		}

		if r.config.EnableTradeBlocking {
			r.dailyBlocked = true
			r.updateBlocked()
		}
	}

	// Check daily profit target
	if !r.dailyTargetHit && r.config.DailyProfitTarget > 0 && todayProfit >= r.config.DailyProfitTarget {
		r.dailyTargetHit = true
		r.logRiskEvent("DAILY_PROFIT_TARGET", "INFO",
			fmt.Sprintf("Daily profit target $%.2f reached", r.config.DailyProfitTarget),
			todayProfit, r.config.DailyProfitTarget,
			r.actionText(false, "", r.config.EnableTradeBlocking, "trading blocked until tomorrow"))

		if r.config.EnableTradeBlocking {
			r.dailyBlocked = true
			r.updateBlocked()
			r.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.LastOperation = "Daily profit target reached - trading blocked"
			})
		}
	}
}

// checkMarginLimits monitors margin levels (a limit of 0 is off).
func (r *RiskManager) checkMarginLimits(marginLevel, margin, equity float64) {
	// Check minimum margin level: while it is too low, one losing position is closed per check
	lowMargin := r.config.MinMarginLevel > 0 && marginLevel > 0 && marginLevel < r.config.MinMarginLevel
	if r.breach("LOW_MARGIN_LEVEL", lowMargin) || (lowMargin && r.config.EnableAutoClose) {
		r.logRiskEvent("LOW_MARGIN_LEVEL", "CRITICAL",
			fmt.Sprintf("Margin level %.0f%% below minimum %.0f%%", marginLevel, r.config.MinMarginLevel),
			marginLevel, r.config.MinMarginLevel, r.actionText(r.config.EnableAutoClose, "closed the most losing position", false, ""))

		if r.config.EnableAutoClose {
			// Close most losing position to free margin
			r.closeMostLosingPosition("Low margin level")
		}
	}

	// Check maximum margin usage: new trades are blocked while it is too high
	usedPercent := 0.0
	if equity > 0 {
		usedPercent = margin / equity * 100
	}
	tooMuchMargin := r.config.MaxMarginUsed > 0 && usedPercent > r.config.MaxMarginUsed
	if r.breach("MARGIN_USAGE", tooMuchMargin) {
		r.logRiskEvent("MARGIN_USAGE", "WARNING",
			fmt.Sprintf("Margin usage %.1f%% exceeds limit %.1f%%", usedPercent, r.config.MaxMarginUsed),
			usedPercent, r.config.MaxMarginUsed,
			r.actionText(false, "", r.config.EnableTradeBlocking, "trading blocked while above the limit"))
	}
	if r.marginBlocked != (tooMuchMargin && r.config.EnableTradeBlocking) {
		r.marginBlocked = tooMuchMargin && r.config.EnableTradeBlocking
		r.updateBlocked()
	}
}

// checkPositionLimits monitors position count, exposure and size limits (a limit of 0 is off).
func (r *RiskManager) checkPositionLimits(positions []*pb.PositionInfo) {
	// Check maximum open positions
	if r.breach("MAX_POSITIONS", r.config.MaxOpenPositions > 0 && len(positions) > r.config.MaxOpenPositions) {
		r.logRiskEvent("MAX_POSITIONS", "WARNING",
			fmt.Sprintf("Open positions %d exceeds limit %d", len(positions), r.config.MaxOpenPositions),
			float64(len(positions)), float64(r.config.MaxOpenPositions), "Logged")
	}

	// Check per-symbol exposure
	symbolCounts := make(map[string]int)
	for _, pos := range positions {
		symbolCounts[pos.Symbol]++
	}
	seen := make(map[string]bool)

	for symbol, count := range symbolCounts {
		key := "SYMBOL_EXPOSURE:" + symbol
		seen[key] = true
		if r.breach(key, r.config.MaxSymbolExposure > 0 && count > r.config.MaxSymbolExposure) {
			r.logRiskEvent("SYMBOL_EXPOSURE", "WARNING",
				fmt.Sprintf("Symbol %s has %d positions, exceeds limit %d", symbol, count, r.config.MaxSymbolExposure),
				float64(count), float64(r.config.MaxSymbolExposure), "Logged")
		}
	}

	// Check maximum position size: with auto-close the excess volume is closed
	for _, pos := range positions {
		key := fmt.Sprintf("POSITION_SIZE:%d", pos.Ticket)
		seen[key] = true
		if r.breach(key, r.config.MaxPositionSize > 0 && pos.Volume > r.config.MaxPositionSize+1e-9) {
			action := "Logged"
			if r.config.EnableAutoClose {
				action = r.reducePosition(pos)
			}
			r.logRiskEvent("POSITION_SIZE", "WARNING",
				fmt.Sprintf("Position #%d %s %.2f lots exceeds limit %.2f", pos.Ticket, pos.Symbol, pos.Volume, r.config.MaxPositionSize),
				pos.Volume, r.config.MaxPositionSize, action)
		}
	}

	// Forget breaches of symbols and positions that are gone
	for key := range r.activeBreaches {
		if (strings.HasPrefix(key, "SYMBOL_EXPOSURE:") || strings.HasPrefix(key, "POSITION_SIZE:")) && !seen[key] {
			delete(r.activeBreaches, key)
		}
	}
}

// reducePosition closes the part of a position above MaxPositionSize (rounded up to the volume step).
func (r *RiskManager) reducePosition(pos *pb.PositionInfo) string {
	spec, ok := r.symbolSpecs[pos.Symbol]
	if !ok {
		var err error
		if spec, err = loadSymbolSpec(r.sugar, pos.Symbol); err != nil {
			r.IncrementError(fmt.Sprintf("failed to reduce #%d: %v", pos.Ticket, err))
			return "Reduce failed: " + err.Error()
		}
		r.symbolSpecs[pos.Symbol] = spec
	}
	step := spec.VolumeStep
	if step <= 0 {
		step = 0.01
	}
	excess := math.Ceil((pos.Volume-r.config.MaxPositionSize)/step-1e-9) * step
	excess = math.Round(excess*1e8) / 1e8

	var err error
	if excess >= pos.Volume-1e-9 || pos.Volume-excess < spec.VolumeMin-1e-9 {
		excess = pos.Volume
		err = r.sugar.ClosePosition(pos.Ticket)
	} else {
		err = r.sugar.ClosePositionPartial(pos.Ticket, excess)
	}
	if err != nil {
		r.IncrementError(fmt.Sprintf("failed to reduce #%d: %v", pos.Ticket, err))
		return "Reduce failed: " + err.Error()
	}
	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.OperationsTotal++
		m.LastOperation = fmt.Sprintf("Reduced #%d by %.2f lots (MaxPositionSize %.2f)", pos.Ticket, excess, r.config.MaxPositionSize)
	})
	return fmt.Sprintf("Closed %.2f lots", excess)
}

// updateBlocked sets tradingBlocked from the daily and the margin block.
func (r *RiskManager) updateBlocked() {
	r.mu.Lock()
	r.tradingBlocked = r.dailyBlocked || r.marginBlocked
	r.mu.Unlock()
}

// enforceTradeBlock makes "blocked" mean something: while trading is blocked, every position that
// was not open when the block started is closed again. The risk manager cannot intercept orders
// sent by other programs, so it undoes them; positions open at the block start are left alone.
func (r *RiskManager) enforceTradeBlock(positions []*pb.PositionInfo) {
	if !r.IsTradingBlocked() || !r.config.EnableTradeBlocking {
		r.allowedTickets = nil
		return
	}

	// First check under the block: remember what is open now
	if r.allowedTickets == nil {
		r.allowedTickets = make(map[uint64]bool, len(positions))
		for _, pos := range positions {
			r.allowedTickets[pos.Ticket] = true
		}
		return
	}

	for _, pos := range positions {
		if r.allowedTickets[pos.Ticket] {
			continue
		}
		action := "Closed (trading is blocked)"
		if err := r.sugar.ClosePosition(pos.Ticket); err != nil {
			r.IncrementError(fmt.Sprintf("failed to close blocked trade #%d: %v", pos.Ticket, err))
			action = "Close failed: " + err.Error()
		} else {
			r.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.OperationsTotal++
				m.LastOperation = fmt.Sprintf("Closed #%d: opened while trading is blocked", pos.Ticket)
			})
		}
		r.logRiskEvent("TRADE_BLOCKED", "WARNING",
			fmt.Sprintf("Position #%d %s opened while trading is blocked", pos.Ticket, pos.Symbol),
			pos.Volume, 0, action)
	}
}

// actionText describes what a breach does with the given switches.
func (r *RiskManager) actionText(closeOn bool, closeText string, blockOn bool, blockText string) string {
	var parts []string
	if closeOn {
		parts = append(parts, closeText)
	}
	if blockOn {
		parts = append(parts, blockText)
	}
	if len(parts) == 0 {
		return "Logged"
	}
	return strings.Join(parts, ", ")
}

// closeAllPositionsEmergency closes all positions immediately.
func (r *RiskManager) closeAllPositionsEmergency(reason string) {
	closed, err := r.sugar.CloseAllPositions()
	if err != nil {
		r.IncrementError(fmt.Sprintf("emergency close: %v", err))
	}

	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("EMERGENCY: Closed %d positions - %s", closed, reason)
		m.OperationsTotal += closed
	})

	action := fmt.Sprintf("Closed %d positions", closed)
	if err != nil {
		action += " (some failed: " + err.Error() + ")"
	}
	r.logRiskEvent("EMERGENCY_CLOSE", "CRITICAL",
		fmt.Sprintf("Closed all %d positions: %s", closed, reason),
		float64(closed), 0, action)
}

// closeMostLosingPosition closes the position with largest loss.
func (r *RiskManager) closeMostLosingPosition(reason string) {
	positions, err := r.sugar.GetOpenPositions()
	if err != nil {
		r.IncrementError(fmt.Sprintf("failed to get positions: %v", err))
		return
	}

	var mostLosingTicket uint64
	mostLoss := 0.0

	for _, pos := range positions {
		if pos.Profit < mostLoss {
			mostLoss = pos.Profit
			mostLosingTicket = pos.Ticket
		}
	}

	if mostLosingTicket > 0 {
		if err := r.sugar.ClosePosition(mostLosingTicket); err != nil {
			r.IncrementError(fmt.Sprintf("failed to close #%d: %v", mostLosingTicket, err))
		} else {
			r.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.OperationsTotal++
				m.LastOperation = fmt.Sprintf("Closed losing position #%d: %s", mostLosingTicket, reason)
			})
		}
	}
}

// checkDailyReset resets daily counters at start of new day (local date).
func (r *RiskManager) checkDailyReset() {
	now := time.Now()
	y1, m1, d1 := now.Date()
	y0, m0, d0 := r.lastResetDate.Date()
	if y1 == y0 && m1 == m0 && d1 == d0 {
		return
	}
	equity, err := r.sugar.GetEquity()
	if err != nil {
		r.IncrementError(fmt.Sprintf("daily reset: failed to get equity: %v", err))
		return
	}
	if equity <= 0 {
		r.IncrementError("daily reset: " + errAccountNotLoaded.Error()) // retried on the next check
		return
	}

	r.mu.Lock()
	r.dailyStartBalance = equity
	r.lastResetDate = now
	r.mu.Unlock()
	r.dailyBlocked, r.dailyLossHit, r.dailyTargetHit = false, false, false
	r.updateBlocked()

	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = "Daily reset performed"
	})
}

// logRiskEvent logs a risk event (risk events are not counted as errors).
func (r *RiskManager) logRiskEvent(eventType, severity, description string, value, limit float64, action string) {
	event := RiskEvent{
		Timestamp:   time.Now(),
		EventType:   eventType,
		Severity:    severity,
		Description: description,
		Value:       value,
		Limit:       limit,
		ActionTaken: action,
	}

	r.mu.Lock()
	r.riskEvents = append(r.riskEvents, event)

	// Keep only last 100 events
	if len(r.riskEvents) > 100 {
		r.riskEvents = r.riskEvents[len(r.riskEvents)-100:]
	}
	r.mu.Unlock()

	fmt.Printf("\n  ⚠️  [%s] %s: %s → %s\n", severity, eventType, description, action)
}

// GetRiskEvents returns a copy of the recent risk events.
func (r *RiskManager) GetRiskEvents() []RiskEvent {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]RiskEvent(nil), r.riskEvents...)
}

// IsTradingBlocked returns whether trading is currently blocked.
func (r *RiskManager) IsTradingBlocked() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tradingBlocked
}

// GetTodayProfit returns today's profit/loss (equity - equity at the start of the day).
func (r *RiskManager) GetTodayProfit() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.todayProfit
}

// GetPeakBalance returns the peak equity the drawdown is measured from
// (since start, or since the last emergency close).
func (r *RiskManager) GetPeakBalance() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.peakBalance
}

// GetDailyStartBalance returns the equity at start of today.
func (r *RiskManager) GetDailyStartBalance() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.dailyStartBalance
}

/*══════════════════════════════════════════════════════════════════════════════

 ██████╗ ██╗███████╗██╗  ██╗    ███╗   ███╗ █████╗ ███╗   ██╗ █████╗  ██████╗ ███████╗██████╗
 ██╔══██╗██║██╔════╝██║ ██╔╝    ████╗ ████║██╔══██╗████╗  ██║██╔══██╗██╔════╝ ██╔════╝██╔══██╗
 ██████╔╝██║███████╗█████╔╝     ██╔████╔██║███████║██╔██╗ ██║███████║██║  ███╗█████╗  ██████╔╝
 ██╔══██╗██║╚════██║██╔═██╗     ██║╚██╔╝██║██╔══██║██║╚██╗██║██╔══██║██║   ██║██╔══╝  ██╔══██╗
 ██║  ██║██║███████║██║  ██╗    ██║ ╚═╝ ██║██║  ██║██║ ╚████║██║  ██║╚██████╔╝███████╗██║  ██║
 ╚═╝  ╚═╝╚═╝╚══════╝╚═╝  ╚═╝    ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝

  DETAILED CONFIGURATION GUIDE
  Located at end of file to keep header clean and focused

══════════════════════════════════════════════════════════════════════════════*/

/*
PROGRAMMATIC USAGE & CONFIGURATION

⚙️ PARAMETER CONFIGURATION IS LOCATED IN main.go

WHY THIS SEPARATION EXISTS:
• 14_risk_manager.go = RISK ENGINE (monitoring logic, protection algorithms)
• main.go → RunOrchestrator_RiskManager() = RUNTIME CONFIGURATION (limits)

THIS SEPARATION IS NEEDED FOR:
1️⃣ Code Reusability
   → Same risk engine can run with different risk limits
   → No need to modify risk logic to change limits

2️⃣ Quick Testing
   → Want stricter limits? Change numbers in main.go
   → Want different thresholds? Again, only change main.go
   → Core protection logic remains untouched

3️⃣ User Examples
   → main.go shows HOW to configure risk management
   → All available parameters and their default values are visible

4️⃣ Centralized Entry Point
   → All orchestrators launch through main.go
   → Single entry point: go run main.go risk → RunOrchestrator_RiskManager()

📍 WHERE TO CONFIGURE PARAMETERS:
main.go → func RunOrchestrator_RiskManager() (CONFIGURATION block)

CONFIGURATION CODE IN main.go:

func RunOrchestrator_RiskManager() error {
    // ... connection code ...

    // ╔════════════════════════════════════════════════════════════╗
    // ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
    // ╚════════════════════════════════════════════════════════════╝
    orchConfig := orchestrators.RiskManagerConfig{
        MaxDrawdownPercent:  10.0,             // ← 10% max drawdown
        MaxDrawdownAbsolute: 1000.0,           // ← $1000 max drawdown
        DailyLossLimit:      500.0,            // ← $500 daily loss limit
        DailyProfitTarget:   1000.0,           // ← $1000 daily profit target
        MinMarginLevel:      150.0,            // ← 150% min margin level
        MaxMarginUsed:       80.0,             // ← 80% max margin usage
        MaxOpenPositions:    20,               // ← Max 20 positions
        MaxSymbolExposure:   5,                // ← Max 5 per symbol
        MaxPositionSize:     1.0,              // ← Max 1.0 lot
        CheckInterval:       5 * time.Second,  // ← Check every 5 seconds
        EnableAutoClose:     true,             // ← Auto-close on breach
        EnableTradeBlocking: true,             // ← Block trades on breach
    }

    riskManager := orchestrators.NewRiskManager(sugar, orchConfig)
    riskManager.Start()
    // ... runs for 15 minutes ...
    riskManager.Stop()
}

💡 EXAMPLE CONFIGURATIONS FOR DIFFERENT RISK PROFILES

╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 1: CONSERVATIVE (LOW RISK - BEGINNERS)                           ║
╚═══════════════════════════════════════════════════════════════════════════╝

RiskManagerConfig{
    MaxDrawdownPercent:  5.0,              // ← Very tight drawdown limit
    MaxDrawdownAbsolute: 500.0,            // ← Small absolute limit
    DailyLossLimit:      200.0,            // ← Stop early on losses
    DailyProfitTarget:   300.0,            // ← Lock in profits early
    MinMarginLevel:      200.0,            // ← High margin safety
    MaxMarginUsed:       50.0,             // ← Conservative margin usage
    MaxOpenPositions:    10,               // ← Limited positions
    MaxSymbolExposure:   2,                // ← Minimal per-symbol exposure
    MaxPositionSize:     0.5,              // ← Small position sizes
    CheckInterval:       3 * time.Second,  // ← Frequent checks
    EnableAutoClose:     true,
    EnableTradeBlocking: true,
}

BEST FOR:
• New traders learning risk management
• Small account sizes ($1,000-$5,000)
• Demo account testing
• High-risk strategies that need tight control
• Scalping or high-frequency trading

EXPECTED BEHAVIOR:
• Quick intervention on losses
• Early profit taking
• Very limited exposure
• Multiple daily resets if limits hit


╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 2: MODERATE (BALANCED RISK - DEFAULT)                            ║
╚═══════════════════════════════════════════════════════════════════════════╝

RiskManagerConfig{
    MaxDrawdownPercent:  10.0,             // ← Moderate drawdown tolerance
    MaxDrawdownAbsolute: 1000.0,           // ← Reasonable absolute limit
    DailyLossLimit:      500.0,            // ← Balanced daily limit
    DailyProfitTarget:   1000.0,           // ← Good profit target
    MinMarginLevel:      150.0,            // ← Standard margin safety
    MaxMarginUsed:       80.0,             // ← Healthy margin usage
    MaxOpenPositions:    20,               // ← Moderate position count
    MaxSymbolExposure:   5,                // ← Reasonable per-symbol
    MaxPositionSize:     1.0,              // ← Standard lot size
    CheckInterval:       5 * time.Second,  // ← Regular checks
    EnableAutoClose:     true,
    EnableTradeBlocking: true,
}

BEST FOR:
• Experienced traders
• Medium account sizes ($5,000-$20,000)
• Swing trading strategies
• Multi-pair trading
• Automated trading systems

EXPECTED BEHAVIOR:
• Balanced protection and freedom
• Allows decent drawdown before intervention
• Moderate position limits
• Good for most trading styles


╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 3: AGGRESSIVE (HIGH RISK - EXPERIENCED)                          ║
╚═══════════════════════════════════════════════════════════════════════════╝

RiskManagerConfig{
    MaxDrawdownPercent:  20.0,             // ← Wider drawdown tolerance
    MaxDrawdownAbsolute: 5000.0,           // ← Large absolute limit
    DailyLossLimit:      2000.0,           // ← Higher daily loss tolerance
    DailyProfitTarget:   5000.0,           // ← Ambitious profit target
    MinMarginLevel:      100.0,            // ← Minimal margin safety
    MaxMarginUsed:       90.0,             // ← Aggressive margin usage
    MaxOpenPositions:    50,               // ← Many concurrent positions
    MaxSymbolExposure:   10,               // ← High per-symbol exposure
    MaxPositionSize:     5.0,              // ← Large position sizes
    CheckInterval:       10 * time.Second, // ← Less frequent checks
    EnableAutoClose:     true,
    EnableTradeBlocking: false,            // ← No trade blocking (risky!)
}

BEST FOR:
• Professional traders only
• Large account sizes ($20,000+)
• Grid trading or martingale strategies
• Portfolio management
• When you can monitor actively

EXPECTED BEHAVIOR:
• Maximum trading freedom
• Late intervention (higher risk)
• Many positions allowed
• Suitable for aggressive strategies

⚠️ WARNING: Only use aggressive profile if you fully understand the risks!


╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 4: MONITORING ONLY (NO AUTO-CLOSE)                               ║
╚═══════════════════════════════════════════════════════════════════════════╝

RiskManagerConfig{
    MaxDrawdownPercent:  15.0,
    MaxDrawdownAbsolute: 2000.0,
    DailyLossLimit:      1000.0,
    DailyProfitTarget:   2000.0,
    MinMarginLevel:      120.0,
    MaxMarginUsed:       85.0,
    MaxOpenPositions:    30,
    MaxSymbolExposure:   7,
    MaxPositionSize:     2.0,
    CheckInterval:       5 * time.Second,
    EnableAutoClose:     false,            // ← Only log, don't close
    EnableTradeBlocking: true,             // ← Block new trades only
}

BEST FOR:
• Manual traders who want alerts
• Testing risk thresholds
• Learning your risk tolerance
• When you prefer manual intervention
• Backtesting and analysis

EXPECTED BEHAVIOR:
• Logs all risk events
• Blocks new trades when limits hit
• Does NOT close existing positions
• You maintain full control


╔═══════════════════════════════════════════════════════════════════════════╗
║ PARAMETER DETAILED EXPLANATIONS                                          ║
╚═══════════════════════════════════════════════════════════════════════════╝

• MaxDrawdownPercent (float64)
  Maximum account drawdown from the peak equity as percentage (0 = off)
  Example: 10.0 = allow 10% loss from the highest equity reached
  Trigger: If equity drops 10% from peak → EMERGENCY CLOSE ALL (once);
           the peak is then reset to the current equity
  Tip: Conservative = 5%, Moderate = 10%, Aggressive = 20%

• MaxDrawdownAbsolute (float64)
  Maximum account drawdown in absolute currency amount (0 = off)
  Example: 1000.0 = allow $1000 loss from peak
  Trigger: If equity drops $1000 from peak → EMERGENCY CLOSE ALL
  Tip: Set to amount you're willing to lose in worst case

• DailyLossLimit (float64)
  Maximum loss allowed in one trading day (equity - equity at day start; 0 = off)
  Example: 500.0 = stop trading after losing $500 today
  Trigger: Daily loss hits $500 → CLOSE ALL + BLOCK TRADING
  Tip: Prevents "revenge trading" after bad day
  Reset: Automatically resets at midnight (new trading day)

• DailyProfitTarget (float64)
  Profit target for the day - stops trading when reached
  Example: 1000.0 = stop after making $1000 today
  Trigger: Daily profit hits $1000 → BLOCK TRADING (positions remain open)
  Tip: Locks in good days, prevents giving back profits
  Use 0 to disable daily profit target

• MinMarginLevel (float64)
  Minimum acceptable margin level percentage
  Example: 150.0 = require at least 150% margin level
  Trigger: Margin drops below 150% → CLOSE MOST LOSING POSITION (every check while low)
  Tip: Prevents margin call from broker
  Formula: Margin Level = (Equity / Margin) × 100

• MaxMarginUsed (float64)
  Maximum margin usage: margin / equity × 100 (0 = off)
  Example: 80.0 = margin may use at most 80% of the equity
  Trigger: Above 80% → WARNING + BLOCK TRADING while above (EnableTradeBlocking)
  Tip: Leave buffer for unexpected market moves

• MaxOpenPositions (int)
  Maximum number of concurrent open positions
  Example: 20 = allow maximum 20 positions at once
  Trigger: Exceeds 20 → WARNING logged (once per breach; 0 = off)
  Tip: Prevents over-exposure across account

• MaxSymbolExposure (int)
  Maximum positions per single symbol
  Example: 5 = max 5 positions on EURUSD
  Trigger: Exceeds 5 → WARNING logged (once per breach; 0 = off)
  Tip: Prevents concentration risk on one pair

• MaxPositionSize (float64)
  Maximum lot size for any single position
  Example: 1.0 = no position larger than 1.0 lots (0 = off)
  Trigger: Larger → WARNING; with EnableAutoClose the excess is closed
  Tip: Limits single-trade risk

• CheckInterval (time.Duration)
  How often to check risk metrics
  Example: 5 * time.Second = check every 5 seconds (0 = 5s)
  Range: 1s = very frequent, 60s = relaxed
  Tip: Match to trading style (scalping = 1-3s, swing = 10-30s)

• EnableAutoClose (bool)
  Whether to automatically close positions on limit breach
  true = auto-close positions when limits hit
  false = only log events, no automatic closing
  Tip: Use true for automated protection, false for manual control

• EnableTradeBlocking (bool)
  Whether to block new trades when limits hit
  true = after the daily limits (until midnight) and while margin usage is
         above MaxMarginUsed, positions opened later are closed again
  false = allow trading even after limits (risky!)
  Tip: Always enable unless you have specific reason not to


╔═══════════════════════════════════════════════════════════════════════════╗
║ HOW RISK EVENTS WORK                                                     ║
╚═══════════════════════════════════════════════════════════════════════════╝

RISK EVENT TYPES:
1. MAX_DRAWDOWN_PERCENT    → Drawdown % exceeded
2. MAX_DRAWDOWN_ABSOLUTE   → Drawdown $ exceeded
3. DAILY_LOSS_LIMIT        → Daily loss limit hit
4. DAILY_PROFIT_TARGET     → Daily profit target reached
5. LOW_MARGIN_LEVEL        → Margin level too low
6. MAX_POSITIONS           → Too many open positions
7. SYMBOL_EXPOSURE         → Too many positions on one symbol
8. EMERGENCY_CLOSE         → Positions closed by risk manager
9. MARGIN_USAGE            → Margin usage above MaxMarginUsed
10. POSITION_SIZE          → Position larger than MaxPositionSize
11. TRADE_BLOCKED          → Position opened while trading is blocked (closed)

SEVERITY LEVELS:
• INFO     → Informational (profit target reached)
• WARNING  → Potential issue (position count high)
• CRITICAL → Immediate action taken (drawdown exceeded, emergency close)

EVENT LOGGING:
• One event per breach (logged when a limit is first exceeded)
• Last 100 events kept in memory
• Accessible via GetRiskEvents()
• Each event contains:
  - Timestamp
  - EventType
  - Severity
  - Description
  - Current value vs limit
  - Action taken


╔═══════════════════════════════════════════════════════════════════════════╗
║ RISK WARNINGS & BEST PRACTICES                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝

⚠️ CRITICAL WARNINGS:

1. EMERGENCY CLOSE ALL IS FINAL
   → When drawdown limits hit, ALL positions close immediately
   → Cannot be undone
   → Make sure your limits are reasonable for your strategy

2. DAILY LIMITS RESET AT MIDNIGHT
   → System time based (ensure server time is correct)
   → May reset mid-session depending on timezone
   → Plan around daily reset times

3. MARGIN LEVEL CALCULATIONS
   → Different brokers calculate margin differently
   → Test your MinMarginLevel on demo first
   → Some brokers have their own margin call levels (50-100%)

4. TRADE BLOCKING BY THE DAILY LIMITS LASTS FOR THE DAY
   → Once the daily limits block trading, it stays blocked until daily reset
   → Only way to unblock: wait for midnight or restart orchestrator
   → The margin usage block ends when the usage is back below MaxMarginUsed
   → Plan your daily limits carefully

5. RISK MANAGER CAN'T REJECT EXTERNAL ORDERS
   → It only sees positions: while trading is blocked it closes positions
     opened after the block started (at its next check)
   → Orders still reach the broker; check IsTradingBlocked() in your code
   → Best used with single trading system


💡 BEST PRACTICES:

✅ Test on DEMO First
   → Run risk manager on demo account for 1 week minimum
   → Ensure limits trigger correctly
   → Verify emergency close works as expected

✅ Set Limits Based on Account Size
   → Drawdown: 5-10% of account
   → Daily loss: 2-5% of account
   → Position size: 1-2% risk per trade

✅ Monitor Risk Events
   → Check GetRiskEvents() regularly
   → Look for patterns (same limit hit repeatedly?)
   → Adjust limits if needed

✅ Combine with Trading Strategy
   → Risk manager complements your strategy, doesn't replace it
   → Your strategy should have its own stop losses
   → Risk manager is "last line of defense"

✅ Use Trade Blocking Wisely
   → Blocks help prevent revenge trading
   → Consider disabling for grid/martingale strategies
   → Re-enable after reviewing why limit was hit

✅ Daily Profit Target Strategy
   → Lock in good days by stopping at target
   → Prevents "giving it all back" syndrome
   → Set realistic targets (don't be too greedy)

✅ Regular Monitoring
   → Even with risk manager, check account regularly
   → Review why limits were hit
   → Adjust strategy or limits accordingly

✅ Understand Your Risk Tolerance
   → Conservative: 5% drawdown, $200 daily loss
   → Moderate: 10% drawdown, $500 daily loss
   → Aggressive: 20% drawdown, $2000 daily loss

═══════════════════════════════════════════════════════════════════════════*/
