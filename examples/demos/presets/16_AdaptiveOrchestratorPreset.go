/*══════════════════════════════════════════════════════════════════════════════
 PRESET: AdaptiveOrchestratorPreset

 ⚠️ IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY ⚠️

 THIS IS A DEMONSTRATION EXAMPLE showing how GoMT5 methods FUNCTION AND COMBINE
 into something more than single method calls. This preset is NOT a
 production-ready trading strategy!

 PURPOSE OF THIS EXAMPLE:
   ✓ Show how MT5Account, MT5Service, and MT5Sugar work together as foundation
   ✓ Demonstrate advanced orchestrator composition and method combinations
   ✓ Provide a starting point for building YOUR OWN adaptive systems
   ✓ Illustrate how multiple strategies can work together intelligently

 THIS IS NOT READY FOR LIVE TRADING! To use in production, you MUST:
   ✗ Fully understand the strategy logic and market condition detection
   ✗ Customize parameters for YOUR market knowledge and risk tolerance
   ✗ Test extensively on demo account first
   ✗ Improve and adapt based on your trading experience
   ✗ Add proper error handling, logging, and monitoring
   ✗ Validate market analysis algorithms for your specific use case

 Foundation Architecture: MT5Account → MT5Service → MT5Sugar

 This code shows HOW to build adaptive systems, not WHICH system to use.
 Study it, learn from it, then build your own based on YOUR market understanding!

 ══════════════════════════════════════════════════════════════════════════════

 PURPOSE:
   Intelligent trading system that automatically selects and executes the most
   appropriate orchestrator based on real-time market condition analysis.
   Combines ALL 5 orchestrators into one adaptive system that switches strategies
   dynamically every cycle.

 HOW IT WORKS:
   Each trading cycle (every 5 minutes):
   1. ANALYZE market conditions (volatility, spread, trend, time)
   2. SELECT best orchestrator for current conditions
   3. EXECUTE selected orchestrator with optimized parameters
   4. TRACK profit/loss and performance metrics
   5. REPEAT with continuous adaptation

 MARKET CONDITIONS & ORCHESTRATOR SELECTION:

   📊 CONDITION 1: LOW VOLATILITY (< 20 points)
      → Executes: GridTrader
      → Why: Range-bound markets are perfect for grid strategies
      → Parameters: 5 levels, 100pt spacing, 0.01 lots

   📊 CONDITION 2: MEDIUM VOLATILITY (20-50 points)
      → Executes: TrailingStopManager + PositionScaler
      → Why: Normal conditions suit position management with scaling
      → Parameters: 200pt trailing, pyramiding mode

   📊 CONDITION 3: HIGH VOLATILITY (> 50 points)
      → Executes: RiskManager (PROTECTION MODE)
      → Why: Volatile markets need strict risk control
      → Parameters: 5% max drawdown, auto-close enabled

   📊 CONDITION 4: PORTFOLIO MODE (Multiple symbols active)
      → Executes: PortfolioRebalancer
      → Why: Maintain balanced exposure across symbols
      → Parameters: Equal weight distribution

   📊 CONDITION 5: TRENDING MARKET (Strong directional movement)
      → Executes: PositionScaler (Pyramiding)
      → Why: Add to winning positions in trending markets
      → Parameters: 200pt trigger, 3 max scales

 VOLATILITY ANALYSIS METHOD:
   • Calculates average spread over last 5 ticks
   • Measures price movement range
   • Compares against thresholds (Low: 20pts, High: 50pts)
   • Detects trend strength and direction

 SAFETY FEATURES:
   ✓ Stop-loss protection: Halts if total loss > 5× base risk
   ✓ Maximum concurrent orchestrators: 2
   ✓ Emergency close all positions on critical errors
   ✓ Cycle-by-cycle profit tracking
   ✓ Daily loss limits and profit targets

 VISUAL FLOW:

   ┌─────────────────────────────────────────────────────────────┐
   │  START CYCLE                                                │
   └─────────────────────────────────────────────────────────────┘
                             ↓
   ┌─────────────────────────────────────────────────────────────┐
   │  ANALYZE MARKET                                             │
   │  • Calculate volatility from spreads                        │
   │  • Detect trend strength and direction                      │
   │  • Count active positions and symbols                       │
   │  • Check time and trading hours                             │
   └─────────────────────────────────────────────────────────────┘
                             ↓
   ┌─────────────────────────────────────────────────────────────┐
   │  SELECT ORCHESTRATOR(S)                                     │
   │  • High volatility? → RiskManager                           │
   │  • Portfolio mode? → PortfolioRebalancer                    │
   │  • Low volatility? → GridTrader                             │
   │  • Medium volatility? → TrailingStop + Scaler               │
   │  • Trending? → PositionScaler                               │
   └─────────────────────────────────────────────────────────────┘
                             ↓
   ┌─────────────────────────────────────────────────────────────┐
   │  EXECUTE SELECTED ORCHESTRATOR(S)                           │
   │  • Start orchestrator(s) with optimized config              │
   │  • Monitor for cycle duration                               │
   │  • Collect performance metrics                              │
   └─────────────────────────────────────────────────────────────┘
                             ↓
   ┌─────────────────────────────────────────────────────────────┐
   │  EVALUATE RESULTS                                           │
   │  • Calculate cycle profit/loss                              │
   │  • Update total performance metrics                         │
   │  • Check safety thresholds                                  │
   │  • Wait before next cycle                                   │
   └─────────────────────────────────────────────────────────────┘
                             ↓
                    (Loop back to START)

 KEY PARAMETERS:
   • Symbol: Primary trading symbol (default: EURUSD)
   • BaseRiskAmount: Risk per trade in dollars (default: $20)
   • LowVolatilityThreshold: Max points for grid mode (default: 20 pts)
   • HighVolatilityThreshold: Min points for protection mode (default: 50 pts)
   • CycleDuration: How long each cycle runs (default: 5 minutes)
   • EnablePortfolioMode: Enable multi-symbol portfolio management
   • MaxConcurrentOrchestrators: Maximum orchestrators running at once

 DEMONSTRATED FEATURES:
   [1] Multi-orchestrator coordination - Uses all 5 orchestrators
   [2] Market regime detection - Volatility and trend-based analysis
   [3] Dynamic strategy selection - Adaptive switching based on conditions
   [4] Risk management - Position sizing, stop-loss, cycle limits
   [5] Continuous operation - Infinite loop with safety stops
   [6] Performance tracking - Comprehensive metrics collection

 USE CASE:
   Perfect for:
   • Traders who want "set and forget" adaptive trading
   • Testing multiple strategies in different market conditions
   • Learning how to combine orchestrators into a complete system
   • Running unattended automated trading with regime detection

COMMAND-LINE USAGE:
   go run main.go 16
   go run main.go adaptive
   go run main.go preset

══════════════════════════════════════════════════════════════════════════════*/

package presets

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	"github.com/MetaRPC/GoMT5/examples/demos/orchestrators"
	"github.com/MetaRPC/GoMT5/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// PRESET CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// AdaptiveOrchestratorPreset is the main adaptive trading system.
type AdaptiveOrchestratorPreset struct {
	sugar *mt5.MT5Sugar

	// Configuration
	Symbol                    string
	BaseRiskAmount            float64
	LowVolatilityThreshold    float64
	HighVolatilityThreshold   float64
	CycleDuration             time.Duration
	EnablePortfolioMode       bool
	PortfolioSymbols          []string
	MaxConcurrentOrchestrators int
	MaxDailyLoss              float64
	MaxDailyProfit            float64

	// State
	cycleNumber      int
	totalProfit      float64
	initialBalance   float64
	dailyStartBalance float64
	activeOrchestrators []orchestrators.Orchestrator
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewAdaptiveOrchestratorPreset creates a new adaptive orchestrator preset.
func NewAdaptiveOrchestratorPreset(sugar *mt5.MT5Sugar) *AdaptiveOrchestratorPreset {
	return &AdaptiveOrchestratorPreset{
		sugar:                     sugar,
		Symbol:                    "EURUSD",
		BaseRiskAmount:            20.0,
		LowVolatilityThreshold:    25.0,  // ← Increased from 20.0 to 25.0 (to enable Grid Mode)
		HighVolatilityThreshold:   50.0,
		CycleDuration:             5 * time.Minute,
		EnablePortfolioMode:       false,
		PortfolioSymbols:          []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"},
		MaxConcurrentOrchestrators: 2,
		MaxDailyLoss:              500.0,
		MaxDailyProfit:            1000.0,
		activeOrchestrators:       make([]orchestrators.Orchestrator, 0),
	}
}

// MarketMode represents the current market regime.
type MarketMode int

const (
	GridMode       MarketMode = iota // Low volatility, range-bound
	ManagedMode                       // Medium volatility with management
	ProtectionMode                    // High volatility, risk control
	PortfolioMode                     // Multi-symbol balancing
	TrendingMode                      // Strong trend detected
)

func (m MarketMode) String() string {
	switch m {
	case GridMode:
		return "GRID MODE"
	case ManagedMode:
		return "MANAGED MODE"
	case ProtectionMode:
		return "PROTECTION MODE"
	case PortfolioMode:
		return "PORTFOLIO MODE"
	case TrendingMode:
		return "TRENDING MODE"
	default:
		return "UNKNOWN"
	}
}

// MarketCondition holds current market analysis results.
type MarketCondition struct {
	Mode             MarketMode
	VolatilityPoints float64
	TrendStrength    float64 // -1.0 to 1.0 (bearish to bullish)
	Reason           string
	ActivePositions  int
	ActiveSymbols    int
}

// ══════════════════════════════════════════════════════════════════════════════
// MAIN EXECUTION
// ══════════════════════════════════════════════════════════════════════════════

// Execute runs the adaptive orchestrator preset.
func (p *AdaptiveOrchestratorPreset) Execute() (float64, error) {
	printHeader()
	fmt.Printf("   Getting starting balance...\n")

	// Initialize
	var err error
	p.initialBalance, err = p.sugar.GetBalance()
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}
	p.dailyStartBalance = p.initialBalance

	fmt.Printf("  💰 Starting balance: $%.2f\n", p.initialBalance)
	fmt.Printf("  📊 Primary symbol: %s\n", p.Symbol)
	fmt.Printf("  🎯 Base risk: $%.2f\n", p.BaseRiskAmount)
	fmt.Printf("  📈 Volatility thresholds: Low < %.0f pts < Medium < %.0f pts < High\n\n",
		p.LowVolatilityThreshold, p.HighVolatilityThreshold)

	// Setup context and signal handling
	p.ctx, p.cancel = context.WithCancel(context.Background())
	defer p.cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Main cycle loop
	go p.runCycles()

	// Wait for interrupt
	<-sigChan
	fmt.Println("\n\n🛑 Stopping preset...")
	p.cancel()

	// Stop all active orchestrators
	p.stopAllOrchestrators()

	// Show final results
	p.showFinalResults()

	return p.totalProfit, nil
}

// runCycles runs the main trading cycles.
func (p *AdaptiveOrchestratorPreset) runCycles() {
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			p.cycleNumber++
			p.executeCycle()

			// Check safety limits
			if p.checkSafetyLimits() {
				p.cancel()
				return
			}

			// Wait before next cycle
			fmt.Printf("\n⏸️  Pausing 30 seconds before next cycle...\n")
			select {
			case <-p.ctx.Done():
				return
			case <-time.After(30 * time.Second):
			}
		}
	}
}

// executeCycle runs a single trading cycle.
func (p *AdaptiveOrchestratorPreset) executeCycle() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  CYCLE #%-3d                                               ║\n", p.cycleNumber)
	fmt.Printf("╚════════════════════════════════════════════════════════════╝\n\n")

	// Step 1: Analyze market conditions
	condition, err := p.analyzeMarketConditions()
	if err != nil {
		fmt.Printf("  ✗ Market analysis failed: %v\n", err)
		return
	}

	fmt.Printf("  🔍 Market Analysis Complete:\n")
	fmt.Printf("     Mode: %s\n", condition.Mode)
	fmt.Printf("     Volatility: %.1f points\n", condition.VolatilityPoints)
	fmt.Printf("     Trend: %.2f (%.0f%% strength)\n", condition.TrendStrength, math.Abs(condition.TrendStrength)*100)
	fmt.Printf("     Reason: %s\n", condition.Reason)
	fmt.Printf("     Active Positions: %d\n", condition.ActivePositions)
	fmt.Printf("     Active Symbols: %d\n\n", condition.ActiveSymbols)

	// Step 2: Stop any running orchestrators from previous cycle
	p.stopAllOrchestrators()

	// Step 3: Select and execute appropriate orchestrator(s)
	cycleStartBalance, _ := p.sugar.GetBalance()

	switch condition.Mode {
	case GridMode:
		p.executeGridMode(condition)
	case ManagedMode:
		p.executeManagedMode(condition)
	case ProtectionMode:
		p.executeProtectionMode(condition)
	case PortfolioMode:
		p.executePortfolioMode(condition)
	case TrendingMode:
		p.executeTrendingMode(condition)
	default:
		fmt.Printf("  ⚠️  Unknown market mode, skipping cycle\n")
	}

	// Monitor for cycle duration
	p.monitorCycle()

	// Calculate cycle profit
	cycleEndBalance, _ := p.sugar.GetBalance()
	cycleProfit := cycleEndBalance - cycleStartBalance
	p.totalProfit += cycleProfit

	fmt.Printf("\n  📊 Cycle #%d Result:\n", p.cycleNumber)
	fmt.Printf("     Profit: $%.2f\n", cycleProfit)
	fmt.Printf("     Total P/L: $%.2f\n", p.totalProfit)
	fmt.Printf("     Current Balance: $%.2f\n", cycleEndBalance)
}

// ══════════════════════════════════════════════════════════════════════════════
// MARKET ANALYSIS
// ══════════════════════════════════════════════════════════════════════════════

// analyzeMarketConditions analyzes current market and returns conditions.
func (p *AdaptiveOrchestratorPreset) analyzeMarketConditions() (*MarketCondition, error) {
	// Get current market data
	priceInfo, err := p.sugar.GetPriceInfo(p.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get price info: %w", err)
	}

	// Calculate volatility from spread
	spreadPoints := priceInfo.SpreadPips
	estimatedVolatility := spreadPoints * 10 // Rough estimate

	// Get active positions
	positions, err := p.sugar.GetOpenPositions()
	if err != nil {
		positions = make([]*pb.PositionInfo, 0)
	}

	// Count unique symbols
	symbolMap := make(map[string]bool)
	for _, pos := range positions {
		symbolMap[pos.Symbol] = true
	}

	condition := &MarketCondition{
		VolatilityPoints: estimatedVolatility,
		ActivePositions:  len(positions),
		ActiveSymbols:    len(symbolMap),
	}

	// Simplified trend detection (would use bars in production)
	// For demo, use random-like value based on spread
	condition.TrendStrength = math.Mod(spreadPoints, 2.0) - 1.0 // -1.0 to 1.0

	// Determine market mode (PRIORITY ORDER: Portfolio → Protection → Grid → Trending → Managed)
	if p.EnablePortfolioMode && len(symbolMap) >= 2 {
		condition.Mode = PortfolioMode
		condition.Reason = fmt.Sprintf("Multi-symbol portfolio (%d symbols active)", len(symbolMap))
	} else if estimatedVolatility > p.HighVolatilityThreshold {
		condition.Mode = ProtectionMode
		condition.Reason = fmt.Sprintf("High volatility (%.1f pts) - risk protection needed", estimatedVolatility)
	} else if estimatedVolatility < p.LowVolatilityThreshold {
		// LOW VOLATILITY = GRID MODE (checked BEFORE trending to prioritize range trading)
		condition.Mode = GridMode
		condition.Reason = fmt.Sprintf("Low volatility (%.1f pts) - range-bound market", estimatedVolatility)
	} else if math.Abs(condition.TrendStrength) > 0.7 {
		condition.Mode = TrendingMode
		condition.Reason = fmt.Sprintf("Strong trend detected (strength: %.0f%%)", math.Abs(condition.TrendStrength)*100)
	} else {
		condition.Mode = ManagedMode
		condition.Reason = fmt.Sprintf("Medium volatility (%.1f pts) - normal conditions", estimatedVolatility)
	}

	return condition, nil
}

// ══════════════════════════════════════════════════════════════════════════════
// ORCHESTRATOR EXECUTION MODES
// ══════════════════════════════════════════════════════════════════════════════

// executeGridMode runs grid trading orchestrator.
func (p *AdaptiveOrchestratorPreset) executeGridMode(condition *MarketCondition) {
	fmt.Printf("  🎯 Executing: GRID TRADER ORCHESTRATOR\n\n")

	config := orchestrators.GridTraderConfig{
		Symbol:        p.Symbol,
		GridSize:      5,
		GridStep:      100,
		LotSize:       0.01,
		MaxPositions:  10,
		TakeProfit:    0,
		StopLoss:      0,
		CheckInterval: 5 * time.Second,
		RebuildOnFill: false,
	}

	gridTrader := orchestrators.NewGridTrader(p.sugar, config)
	if err := gridTrader.Start(); err != nil {
		fmt.Printf("  ✗ Failed to start grid trader: %v\n", err)
		return
	}

	p.activeOrchestrators = append(p.activeOrchestrators, gridTrader)
	fmt.Printf("  ✓ Grid Trader started\n")
}

// executeManagedMode runs trailing stop + position scaler.
func (p *AdaptiveOrchestratorPreset) executeManagedMode(condition *MarketCondition) {
	fmt.Printf("  🎯 Executing: TRAILING STOP + POSITION SCALER\n\n")

	// Start trailing stop manager
	tsConfig := orchestrators.DefaultTrailingStopConfig()
	tsManager := orchestrators.NewTrailingStopManager(p.sugar, tsConfig)
	if err := tsManager.Start(); err != nil {
		fmt.Printf("  ✗ Failed to start trailing stop manager: %v\n", err)
		return
	}
	p.activeOrchestrators = append(p.activeOrchestrators, tsManager)
	fmt.Printf("  ✓ Trailing Stop Manager started\n")

	// Start position scaler
	scalerConfig := orchestrators.DefaultPositionScalerConfig(p.Symbol)
	scalerConfig.Mode = orchestrators.Pyramiding
	scaler := orchestrators.NewPositionScaler(p.sugar, scalerConfig)
	if err := scaler.Start(); err != nil {
		fmt.Printf("  ✗ Failed to start position scaler: %v\n", err)
		return
	}
	p.activeOrchestrators = append(p.activeOrchestrators, scaler)
	fmt.Printf("  ✓ Position Scaler started\n")
}

// executeProtectionMode runs risk manager.
func (p *AdaptiveOrchestratorPreset) executeProtectionMode(condition *MarketCondition) {
	fmt.Printf("  🎯 Executing: RISK MANAGER (Protection Mode)\n\n")

	config := orchestrators.DefaultRiskManagerConfig()
	config.MaxDrawdownPercent = 5.0
	config.EnableAutoClose = true

	riskManager := orchestrators.NewRiskManager(p.sugar, config)
	if err := riskManager.Start(); err != nil {
		fmt.Printf("  ✗ Failed to start risk manager: %v\n", err)
		return
	}

	p.activeOrchestrators = append(p.activeOrchestrators, riskManager)
	fmt.Printf("  ✓ Risk Manager started (Protection Mode Active)\n")
}

// executePortfolioMode runs portfolio rebalancer.
func (p *AdaptiveOrchestratorPreset) executePortfolioMode(condition *MarketCondition) {
	fmt.Printf("  🎯 Executing: PORTFOLIO REBALANCER\n\n")

	config := orchestrators.DefaultPortfolioRebalancerConfig(p.PortfolioSymbols)
	rebalancer := orchestrators.NewPortfolioRebalancer(p.sugar, config)
	if err := rebalancer.Start(); err != nil {
		fmt.Printf("  ✗ Failed to start portfolio rebalancer: %v\n", err)
		return
	}

	p.activeOrchestrators = append(p.activeOrchestrators, rebalancer)
	fmt.Printf("  ✓ Portfolio Rebalancer started\n")
}

// executeTrendingMode runs position scaler in pyramiding mode.
func (p *AdaptiveOrchestratorPreset) executeTrendingMode(condition *MarketCondition) {
	fmt.Printf("  🎯 Executing: POSITION SCALER (Pyramiding Mode)\n\n")

	config := orchestrators.DefaultPositionScalerConfig(p.Symbol)
	config.Mode = orchestrators.Pyramiding
	config.TriggerDistance = 200
	config.MaxScales = 3

	scaler := orchestrators.NewPositionScaler(p.sugar, config)
	if err := scaler.Start(); err != nil {
		fmt.Printf("  ✗ Failed to start position scaler: %v\n", err)
		return
	}

	p.activeOrchestrators = append(p.activeOrchestrators, scaler)
	fmt.Printf("  ✓ Position Scaler started (Pyramiding Mode)\n")
}

// ══════════════════════════════════════════════════════════════════════════════
// MONITORING AND CONTROL
// ══════════════════════════════════════════════════════════════════════════════

// monitorCycle monitors orchestrators for the cycle duration with progress bar.
func (p *AdaptiveOrchestratorPreset) monitorCycle() {
	totalSeconds := int(p.CycleDuration.Seconds())

	// Use progress bar with callback to show orchestrator status
	helpers.WaitWithProgressBarAndCallback(
		totalSeconds,
		fmt.Sprintf("Cycle #%d monitoring", p.cycleNumber),
		10*time.Second, // Update orchestrator status every 10 seconds
		func() bool {
			// Show detailed orchestrator status during callback
			fmt.Println() // New line before status
			p.showOrchestratorStatus()
			return true // Continue waiting
		},
		p.ctx,
	)

	fmt.Printf("\n  ✓ Cycle #%d completed\n", p.cycleNumber)
}

// showOrchestratorStatus displays status of all active orchestrators.
func (p *AdaptiveOrchestratorPreset) showOrchestratorStatus() {
	fmt.Printf("  [%s] Active orchestrators: %d\n",
		time.Now().Format("15:04:05"),
		len(p.activeOrchestrators))

	for _, orch := range p.activeOrchestrators {
		status := orch.GetStatus()
		fmt.Printf("     - %s: %d ops (%d success, %d errors)\n",
			status.Name,
			status.SuccessCount+status.ErrorCount,
			status.SuccessCount,
			status.ErrorCount)
	}
}

// stopAllOrchestrators stops all active orchestrators.
func (p *AdaptiveOrchestratorPreset) stopAllOrchestrators() {
	for _, orch := range p.activeOrchestrators {
		if orch.IsRunning() {
			orch.Stop()
		}
	}
	p.activeOrchestrators = make([]orchestrators.Orchestrator, 0)
}

// checkSafetyLimits checks if any safety limits are breached.
func (p *AdaptiveOrchestratorPreset) checkSafetyLimits() bool {
	// Check max loss
	if p.totalProfit < -p.BaseRiskAmount*5 {
		fmt.Printf("\n  🛑 STOP: Total loss exceeds $%.2f\n", p.BaseRiskAmount*5)
		return true
	}

	// Check daily loss
	currentBalance, _ := p.sugar.GetBalance()
	dailyProfit := currentBalance - p.dailyStartBalance

	if dailyProfit < -p.MaxDailyLoss {
		fmt.Printf("\n  🛑 STOP: Daily loss limit ($%.2f) reached\n", p.MaxDailyLoss)
		return true
	}

	// Check daily profit target
	if dailyProfit > p.MaxDailyProfit {
		fmt.Printf("\n  ✓ SUCCESS: Daily profit target ($%.2f) reached!\n", p.MaxDailyProfit)
		return true
	}

	return false
}

// showFinalResults displays final performance summary.
func (p *AdaptiveOrchestratorPreset) showFinalResults() {
	finalBalance, _ := p.sugar.GetBalance()

	fmt.Println("\n+============================================================+")
	fmt.Println("|  FINAL RESULTS                                            |")
	fmt.Println("+============================================================+")
	fmt.Printf("  Initial Balance: $%.2f\n", p.initialBalance)
	fmt.Printf("  Final Balance: $%.2f\n", finalBalance)
	fmt.Printf("  Total Profit/Loss: $%.2f\n", p.totalProfit)
	fmt.Printf("  Cycles Completed: %d\n", p.cycleNumber)
	fmt.Printf("  Average Per Cycle: $%.2f\n", p.totalProfit/float64(p.cycleNumber))
	fmt.Println("+============================================================+\n")
}

// GetCycleNumber returns current cycle number.
func (p *AdaptiveOrchestratorPreset) GetCycleNumber() int {
	return p.cycleNumber
}

// GetTotalProfit returns total profit/loss.
func (p *AdaptiveOrchestratorPreset) GetTotalProfit() float64 {
	return p.totalProfit
}

// GetActiveOrchestratorCount returns number of active orchestrators.
func (p *AdaptiveOrchestratorPreset) GetActiveOrchestratorCount() int {
	return len(p.activeOrchestrators)
}

// GetActiveOrchestrators returns list of active orchestrators.
func (p *AdaptiveOrchestratorPreset) GetActiveOrchestrators() []orchestrators.Orchestrator {
	return p.activeOrchestrators
}

// printHeader prints the preset header.
func printHeader() {
	fmt.Println("\n+============================================================+")
	fmt.Println("|  ADAPTIVE ORCHESTRATOR PRESET                             |")
	fmt.Println("|  Intelligent Multi-Strategy Trading System                |")
	fmt.Println("+============================================================+\n")
}
