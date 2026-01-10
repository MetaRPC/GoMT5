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
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT5/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// RiskManagerConfig holds risk management parameters.
type RiskManagerConfig struct {
	// Drawdown Limits
	MaxDrawdownPercent  float64 // Maximum account drawdown percentage
	MaxDrawdownAbsolute float64 // Maximum absolute drawdown amount

	// Daily Limits
	DailyLossLimit   float64 // Maximum daily loss allowed
	DailyProfitTarget float64 // Stop trading after hitting this profit

	// Margin Limits
	MinMarginLevel float64 // Minimum margin level percentage
	MaxMarginUsed  float64 // Maximum margin usage percentage

	// Position Limits
	MaxOpenPositions  int     // Maximum concurrent positions
	MaxSymbolExposure int     // Maximum positions per symbol
	MaxPositionSize   float64 // Maximum lot size per position

	// Operational
	CheckInterval      time.Duration // How often to check risk
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

// ══════════════════════════════════════════════════════════════════════════════
// RISK MANAGER IMPLEMENTATION
// ══════════════════════════════════════════════════════════════════════════════

// RiskManager monitors and enforces account-level risk limits.
type RiskManager struct {
	*BaseOrchestrator
	sugar  *mt5.MT5Sugar
	config RiskManagerConfig

	// State
	startingBalance   float64
	peakBalance       float64
	dailyStartBalance float64
	todayProfit       float64
	tradingBlocked    bool
	lastResetDate     time.Time

	// Risk Events
	riskEvents []RiskEvent
}

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
	}
}

// Start begins risk monitoring.
func (r *RiskManager) Start() error {
	if r.IsRunning() {
		return fmt.Errorf("risk manager already running")
	}

	// Initialize baseline values
	if err := r.initialize(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	r.SetContext(ctx, cancel)

	// Mark as started
	r.MarkStarted()

	// Start monitoring loop
	go r.monitorLoop()

	return nil
}

// Stop gracefully stops risk monitoring.
func (r *RiskManager) Stop() error {
	if !r.IsRunning() {
		return fmt.Errorf("risk manager not running")
	}

	// Cancel context
	r.CancelContext()

	// Mark as stopped
	r.MarkStopped()

	return nil
}

// initialize sets up baseline values.
func (r *RiskManager) initialize() error {
	// Get starting balance
	balance, err := r.sugar.GetBalance()
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	r.startingBalance = balance
	r.peakBalance = balance
	r.dailyStartBalance = balance
	r.tradingBlocked = false

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

	balance, _ := r.sugar.GetBalance()
	marginLevel, _ := r.sugar.GetMarginLevel()
	_, _ = r.sugar.GetProfit() // Get profit (for future use)

	// Update peak balance
	if equity > r.peakBalance {
		r.peakBalance = equity
	}

	// Calculate current drawdown
	currentDrawdown := r.peakBalance - equity
	drawdownPercent := (currentDrawdown / r.peakBalance) * 100

	// Calculate today's profit
	r.todayProfit = balance - r.dailyStartBalance

	// Update metrics
	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.CurrentDrawdown = currentDrawdown
		m.MaxDrawdown = drawdownPercent
	})

	// Check each risk limit
	r.checkDrawdownLimit(drawdownPercent, currentDrawdown)
	r.checkDailyLimits()
	r.checkMarginLimits(marginLevel)
	r.checkPositionLimits()

	// Update status
	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("Monitoring: Equity=%.2f, DD=%.1f%%, Margin=%.0f%%",
			equity, drawdownPercent, marginLevel)
	})
}

// checkDrawdownLimit monitors maximum drawdown.
func (r *RiskManager) checkDrawdownLimit(drawdownPercent, drawdownAbsolute float64) {
	// Check percentage drawdown
	if drawdownPercent >= r.config.MaxDrawdownPercent {
		r.logRiskEvent("MAX_DRAWDOWN_PERCENT", "CRITICAL",
			fmt.Sprintf("Drawdown %.1f%% exceeds limit %.1f%%", drawdownPercent, r.config.MaxDrawdownPercent),
			drawdownPercent, r.config.MaxDrawdownPercent)

		if r.config.EnableAutoClose {
			r.closeAllPositionsEmergency("Maximum drawdown exceeded")
		}
	}

	// Check absolute drawdown
	if drawdownAbsolute >= r.config.MaxDrawdownAbsolute {
		r.logRiskEvent("MAX_DRAWDOWN_ABSOLUTE", "CRITICAL",
			fmt.Sprintf("Drawdown $%.2f exceeds limit $%.2f", drawdownAbsolute, r.config.MaxDrawdownAbsolute),
			drawdownAbsolute, r.config.MaxDrawdownAbsolute)

		if r.config.EnableAutoClose {
			r.closeAllPositionsEmergency("Maximum absolute drawdown exceeded")
		}
	}
}

// checkDailyLimits monitors daily profit/loss limits.
func (r *RiskManager) checkDailyLimits() {
	// Check daily loss limit
	if r.todayProfit < 0 && (-r.todayProfit) >= r.config.DailyLossLimit {
		r.logRiskEvent("DAILY_LOSS_LIMIT", "CRITICAL",
			fmt.Sprintf("Daily loss $%.2f exceeds limit $%.2f", -r.todayProfit, r.config.DailyLossLimit),
			-r.todayProfit, r.config.DailyLossLimit)

		if r.config.EnableAutoClose {
			r.closeAllPositionsEmergency("Daily loss limit exceeded")
		}

		if r.config.EnableTradeBlocking {
			r.tradingBlocked = true
		}
	}

	// Check daily profit target
	if r.config.DailyProfitTarget > 0 && r.todayProfit >= r.config.DailyProfitTarget {
		r.logRiskEvent("DAILY_PROFIT_TARGET", "INFO",
			fmt.Sprintf("Daily profit target $%.2f reached", r.config.DailyProfitTarget),
			r.todayProfit, r.config.DailyProfitTarget)

		if r.config.EnableTradeBlocking {
			r.tradingBlocked = true
			r.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.LastOperation = "Daily profit target reached - trading blocked"
			})
		}
	}
}

// checkMarginLimits monitors margin levels.
func (r *RiskManager) checkMarginLimits(marginLevel float64) {
	// Check minimum margin level
	if marginLevel > 0 && marginLevel < r.config.MinMarginLevel {
		r.logRiskEvent("LOW_MARGIN_LEVEL", "CRITICAL",
			fmt.Sprintf("Margin level %.0f%% below minimum %.0f%%", marginLevel, r.config.MinMarginLevel),
			marginLevel, r.config.MinMarginLevel)

		if r.config.EnableAutoClose {
			// Close most losing position to free margin
			r.closeMostLosingPosition("Low margin level")
		}
	}
}

// checkPositionLimits monitors position count and exposure limits.
func (r *RiskManager) checkPositionLimits() {
	positions, err := r.sugar.GetOpenPositions()
	if err != nil {
		return
	}

	// Check maximum open positions
	if len(positions) > r.config.MaxOpenPositions {
		r.logRiskEvent("MAX_POSITIONS", "WARNING",
			fmt.Sprintf("Open positions %d exceeds limit %d", len(positions), r.config.MaxOpenPositions),
			float64(len(positions)), float64(r.config.MaxOpenPositions))
	}

	// Check per-symbol exposure
	symbolCounts := make(map[string]int)
	for _, pos := range positions {
		symbolCounts[pos.Symbol]++
	}

	for symbol, count := range symbolCounts {
		if count > r.config.MaxSymbolExposure {
			r.logRiskEvent("SYMBOL_EXPOSURE", "WARNING",
				fmt.Sprintf("Symbol %s has %d positions, exceeds limit %d", symbol, count, r.config.MaxSymbolExposure),
				float64(count), float64(r.config.MaxSymbolExposure))
		}
	}
}

// closeAllPositionsEmergency closes all positions immediately.
func (r *RiskManager) closeAllPositionsEmergency(reason string) {
	closed, err := r.sugar.CloseAllPositions()
	if err != nil {
		r.IncrementError(fmt.Sprintf("emergency close failed: %v", err))
		return
	}

	r.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("EMERGENCY: Closed %d positions - %s", closed, reason)
		m.OperationsTotal += closed
	})

	r.logRiskEvent("EMERGENCY_CLOSE", "CRITICAL",
		fmt.Sprintf("Closed all %d positions: %s", closed, reason),
		float64(closed), 0)
}

// closeMostLosingPosition closes the position with largest loss.
func (r *RiskManager) closeMostLosingPosition(reason string) {
	positions, err := r.sugar.GetOpenPositions()
	if err != nil {
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
		if err := r.sugar.ClosePosition(mostLosingTicket); err == nil {
			r.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.LastOperation = fmt.Sprintf("Closed losing position #%d: %s", mostLosingTicket, reason)
			})
		}
	}
}

// checkDailyReset resets daily counters at start of new day.
func (r *RiskManager) checkDailyReset() {
	now := time.Now()
	if now.Day() != r.lastResetDate.Day() {
		balance, err := r.sugar.GetBalance()
		if err == nil {
			r.dailyStartBalance = balance
			r.tradingBlocked = false
			r.lastResetDate = now

			r.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.LastOperation = "Daily reset performed"
			})
		}
	}
}

// logRiskEvent logs a risk event.
func (r *RiskManager) logRiskEvent(eventType, severity, description string, value, limit float64) {
	event := RiskEvent{
		Timestamp:   time.Now(),
		EventType:   eventType,
		Severity:    severity,
		Description: description,
		Value:       value,
		Limit:       limit,
		ActionTaken: "Logged",
	}

	r.riskEvents = append(r.riskEvents, event)

	// Keep only last 100 events
	if len(r.riskEvents) > 100 {
		r.riskEvents = r.riskEvents[len(r.riskEvents)-100:]
	}

	r.IncrementError(description)
}

// GetRiskEvents returns recent risk events.
func (r *RiskManager) GetRiskEvents() []RiskEvent {
	return r.riskEvents
}

// IsTradingBlocked returns whether trading is currently blocked.
func (r *RiskManager) IsTradingBlocked() bool {
	return r.tradingBlocked
}

// GetTodayProfit returns today's profit/loss.
func (r *RiskManager) GetTodayProfit() float64 {
	return r.todayProfit
}

// GetPeakBalance returns the peak balance since start.
func (r *RiskManager) GetPeakBalance() float64 {
	return r.peakBalance
}

// GetDailyStartBalance returns the balance at start of today.
func (r *RiskManager) GetDailyStartBalance() float64 {
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
main.go → func RunOrchestrator_RiskManager() (lines 633-646)

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
  Maximum account drawdown from peak balance as percentage
  Example: 10.0 = allow 10% loss from highest balance reached
  Trigger: If equity drops 10% from peak → EMERGENCY CLOSE ALL
  Tip: Conservative = 5%, Moderate = 10%, Aggressive = 20%

• MaxDrawdownAbsolute (float64)
  Maximum account drawdown in absolute currency amount
  Example: 1000.0 = allow $1000 loss from peak
  Trigger: If equity drops $1000 from peak → EMERGENCY CLOSE ALL
  Tip: Set to amount you're willing to lose in worst case

• DailyLossLimit (float64)
  Maximum loss allowed in one trading day
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
  Trigger: Margin drops below 150% → CLOSE MOST LOSING POSITION
  Tip: Prevents margin call from broker
  Formula: Margin Level = (Equity / Margin) × 100

• MaxMarginUsed (float64)
  Maximum margin usage as percentage of available margin
  Example: 80.0 = use maximum 80% of margin
  Currently logged only (future: could block new trades)
  Tip: Leave buffer for unexpected market moves

• MaxOpenPositions (int)
  Maximum number of concurrent open positions
  Example: 20 = allow maximum 20 positions at once
  Trigger: Exceeds 20 → WARNING logged
  Tip: Prevents over-exposure across account

• MaxSymbolExposure (int)
  Maximum positions per single symbol
  Example: 5 = max 5 positions on EURUSD
  Trigger: Exceeds 5 → WARNING logged
  Tip: Prevents concentration risk on one pair

• MaxPositionSize (float64)
  Maximum lot size for any single position
  Example: 1.0 = no position larger than 1.0 lots
  Currently logged only (future: could enforce on new trades)
  Tip: Limits single-trade risk

• CheckInterval (time.Duration)
  How often to check risk metrics
  Example: 5 * time.Second = check every 5 seconds
  Range: 1s = very frequent, 60s = relaxed
  Tip: Match to trading style (scalping = 1-3s, swing = 10-30s)

• EnableAutoClose (bool)
  Whether to automatically close positions on limit breach
  true = auto-close positions when limits hit
  false = only log events, no automatic closing
  Tip: Use true for automated protection, false for manual control

• EnableTradeBlocking (bool)
  Whether to block new trades when limits hit
  true = prevent new positions after daily limits
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

SEVERITY LEVELS:
• INFO     → Informational (profit target reached)
• WARNING  → Potential issue (position count high)
• CRITICAL → Immediate action taken (drawdown exceeded, emergency close)

EVENT LOGGING:
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

4. TRADE BLOCKING IS PERMANENT FOR THE DAY
   → Once trading is blocked, it stays blocked until daily reset
   → Only way to unblock: wait for midnight or restart orchestrator
   → Plan your daily limits carefully

5. RISK MANAGER DOESN'T PREVENT EXTERNAL TRADES
   → Only monitors positions, doesn't control external systems
   → If you trade manually or with other bots, risk manager can't block those
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
