/*══════════════════════════════════════════════════════════════════════════════
 ORCHESTRATOR: PortfolioRebalancer

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
   Maintains balanced exposure across multiple symbols by automatically
   rebalancing positions to match target allocations.

 STRATEGY:
   • Monitors actual exposure vs. target (every 30 min)
   • Rebalances when deviation exceeds threshold (>10%)
   • Opens/closes positions to restore balance

 HOW TO RUN (Command-Line):
   Step 1: Open terminal/command prompt
   Step 2: Navigate to: examples/demos/
   Step 3: Run one of these commands:

      go run main.go 15           ← Run by number (Portfolio Rebalancer is #15)
      go run main.go rebalancer   ← Run by short name
      go run main.go portfolio    ← Run by full name

   What happens:
   • Connects to your MT5 account
   • Starts monitoring portfolio balance every 30 minutes
   • Automatically rebalances when needed
   • Runs for 1 hour (default), then stops
   • You'll see live updates in console

 CONFIGURATION:
   ⚙️ All parameters configured in main.go → RunOrchestrator_PortfolioRebalancer()
   📍 See end of this file for detailed configuration examples and documentation

══════════════════════════════════════════════════════════════════════════════*/

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                    📚 BEGINNER'S GUIDE - SIMPLE EXPLANATION                  ║
╚══════════════════════════════════════════════════════════════════════════════╝

WHAT IS A PORTFOLIO?
────────────────────
Imagine you don't trade just one currency pair (like EURUSD), but trade SEVERAL
pairs at the same time. This collection of different trading instruments is
called a "portfolio".

Example:
  Instead of trading only EURUSD, you trade:
  • EURUSD (Euro/Dollar)
  • GBPUSD (British Pound/Dollar)
  • USDJPY (Dollar/Japanese Yen)
  • XAUUSD (Gold)

WHY TRADE MULTIPLE SYMBOLS?
────────────────────────────
  ✓ DIVERSIFICATION: If EURUSD moves against you, other pairs might profit
  ✓ RISK REDUCTION: Don't put all eggs in one basket
  ✓ MORE OPPORTUNITIES: Different markets move at different times

WHAT IS REBALANCING?
────────────────────
Rebalancing means keeping your portfolio BALANCED - maintaining specific
percentages for each symbol.

Real-World Analogy:
  Think of 4 jars where you store coins. You want each jar to always have
  exactly 25% of your total coins. But as you add/remove coins, the balance
  shifts. Rebalancing means moving coins between jars to restore 25% in each.

CONCRETE EXAMPLE:
─────────────────
Let's say you have $10,000 and want to split it equally among 4 currency pairs:

  Target Portfolio (what you WANT):
  ┌─────────────────┬─────────┬────────────┐
  │ Symbol          │ Target  │ Amount     │
  ├─────────────────┼─────────┼────────────┤
  │ EURUSD          │ 25%     │ $2,500     │
  │ GBPUSD          │ 25%     │ $2,500     │
  │ USDJPY          │ 25%     │ $2,500     │
  │ XAUUSD (Gold)   │ 25%     │ $2,500     │
  ├─────────────────┼─────────┼────────────┤
  │ TOTAL           │ 100%    │ $10,000    │
  └─────────────────┴─────────┴────────────┘

WHAT HAPPENS WHEN BALANCE SHIFTS?
──────────────────────────────────
After trading for a while, prices change and your actual portfolio might look
like this:

  Actual Portfolio (what you HAVE):
  ┌─────────────────┬─────────┬────────────┬───────────┐
  │ Symbol          │ Target  │ Actual     │ Deviation │
  ├─────────────────┼─────────┼────────────┼───────────┤
  │ EURUSD          │ 25%     │ $3,200     │ 32% ⚠️    │  Too much!
  │ GBPUSD          │ 25%     │ $2,100     │ 21% ⚠️    │  Too little
  │ USDJPY          │ 25%     │ $2,600     │ 26% ✓     │  Close enough
  │ XAUUSD (Gold)   │ 25%     │ $2,100     │ 21% ⚠️    │  Too little
  ├─────────────────┼─────────┼────────────┼───────────┤
  │ TOTAL           │ 100%    │ $10,000    │           │
  └─────────────────┴─────────┴────────────┴───────────┘

HOW THE REBALANCER FIXES THIS:
───────────────────────────────
The orchestrator automatically:

  1. DETECTS imbalance: "EURUSD is 32% but should be 25%"

  2. CALCULATES corrections:
     • EURUSD: Reduce by $700 (from $3,200 to $2,500) → SELL
     • GBPUSD: Increase by $400 (from $2,100 to $2,500) → BUY
     • USDJPY: OK, no action needed
     • XAUUSD: Increase by $400 (from $2,100 to $2,500) → BUY

  3. EXECUTES trades to restore balance

  4. RESULT: Portfolio is balanced again at 25% each!

HOW IT WORKS IN PRACTICE:
──────────────────────────
  • Every 30 minutes, checks if portfolio is balanced
  • If deviation > 10%, triggers rebalancing
  • Automatically opens/closes positions to restore target %
  • Keeps running until you stop it

WHEN TO USE THIS ORCHESTRATOR:
───────────────────────────────
  ✓ You want to trade multiple currency pairs
  ✓ You want automatic diversification
  ✓ You want to maintain specific allocation percentages
  ✓ You have a long-term strategy (not scalping)

WHEN NOT TO USE:
────────────────
  ✗ You only trade one symbol (no portfolio needed)
  ✗ You want manual control over each position
  ✗ You're doing high-frequency scalping
  ✗ You don't want automated rebalancing

KEY SETTINGS:
─────────────
  • Allocations: Which symbols and what % for each (must total 100%)
  • TotalExposure: Total portfolio value ($10,000 in example)
  • RebalanceThreshold: How far off before rebalancing (e.g., 10%)
  • CheckInterval: How often to check (e.g., every 30 minutes)

SAFETY NOTES:
─────────────
  ⚠️ This orchestrator WILL open and close positions automatically
  ⚠️ Make sure your allocations sum to exactly 100%
  ⚠️ Test with small amounts first
  ⚠️ Monitor the first few rebalancing cycles manually

For detailed configuration examples, see the end of this file.

═══════════════════════════════════════════════════════════════════════════════
*/

package orchestrators

import (
	"context"
	"fmt"
	"math"
	"time"

	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// PortfolioRebalancerConfig holds portfolio rebalancing parameters.
type PortfolioRebalancerConfig struct {
	// Portfolio Definition
	Allocations map[string]float64 // Symbol -> target allocation percentage (0.0-1.0)

	// Exposure Limits
	TotalExposure   float64 // Total portfolio value to maintain
	MaxDeviation    float64 // Maximum deviation from target before rebalancing

	// Rebalancing Rules
	RebalanceThreshold float64       // % deviation to trigger rebalance
	CheckInterval      time.Duration // How often to check balance
	MinPositionSize    float64       // Minimum lot size for positions

	// Trading Parameters
	UseMarketOrders   bool    // Use market orders (true) or limit orders (false)
	SlippageTolerance float64 // Maximum slippage in points for limit orders
	MaxTradesPerCycle int     // Max trades per rebalancing cycle
}

// DefaultPortfolioRebalancerConfig returns sensible defaults.
func DefaultPortfolioRebalancerConfig(symbols []string) PortfolioRebalancerConfig {
	// Create equal-weight allocation
	allocations := make(map[string]float64)
	weight := 1.0 / float64(len(symbols))
	for _, symbol := range symbols {
		allocations[symbol] = weight
	}

	return PortfolioRebalancerConfig{
		Allocations:        allocations,
		TotalExposure:      10000.0,
		MaxDeviation:       0.20,
		RebalanceThreshold: 10.0,
		CheckInterval:      1 * time.Hour,
		MinPositionSize:    0.01,
		UseMarketOrders:    true,
		SlippageTolerance:  50,
		MaxTradesPerCycle:  10,
	}
}

// ══════════════════════════════════════════════════════════════════════════════
// PORTFOLIO REBALANCER IMPLEMENTATION
// ══════════════════════════════════════════════════════════════════════════════

// PortfolioRebalancer maintains target allocation across symbols.
type PortfolioRebalancer struct {
	*BaseOrchestrator
	sugar  *mt5.MT5Sugar
	config PortfolioRebalancerConfig

	// Portfolio State
	currentAllocations map[string]float64 // Current exposure per symbol
	targetValues       map[string]float64 // Target $ value per symbol
	lastRebalance      time.Time
	rebalanceCount     int
}

// SymbolAllocation represents current vs target allocation for a symbol.
type SymbolAllocation struct {
	Symbol          string
	TargetPercent   float64
	CurrentPercent  float64
	TargetValue     float64
	CurrentValue    float64
	Deviation       float64
	DeviationPercent float64
	NeedsAdjustment bool
	ActionRequired  string // "BUY", "SELL", "HOLD"
	AdjustmentValue float64
}

// NewPortfolioRebalancer creates a new portfolio rebalancing orchestrator.
func NewPortfolioRebalancer(sugar *mt5.MT5Sugar, config PortfolioRebalancerConfig) *PortfolioRebalancer {
	return &PortfolioRebalancer{
		BaseOrchestrator:   NewBaseOrchestrator("Portfolio Rebalancer"),
		sugar:              sugar,
		config:             config,
		currentAllocations: make(map[string]float64),
		targetValues:       make(map[string]float64),
		lastRebalance:      time.Now(),
	}
}

// Start begins portfolio monitoring and rebalancing.
func (p *PortfolioRebalancer) Start() error {
	if p.IsRunning() {
		return fmt.Errorf("portfolio rebalancer already running")
	}

	// Validate configuration
	if err := p.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Calculate target values
	p.calculateTargetValues()

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	p.SetContext(ctx, cancel)

	// Mark as started
	p.MarkStarted()

	// Start monitoring loop
	go p.monitorLoop()

	return nil
}

// Stop gracefully stops portfolio rebalancing.
func (p *PortfolioRebalancer) Stop() error {
	if !p.IsRunning() {
		return fmt.Errorf("portfolio rebalancer not running")
	}

	// Cancel context
	p.CancelContext()

	// Mark as stopped
	p.MarkStopped()

	return nil
}

// validateConfig ensures configuration is valid.
func (p *PortfolioRebalancer) validateConfig() error {
	// Check allocations sum to 100%
	totalAllocation := 0.0
	for _, allocation := range p.config.Allocations {
		totalAllocation += allocation
	}

	if math.Abs(totalAllocation-1.0) > 0.001 {
		return fmt.Errorf("allocations must sum to 100%%, got %.2f%%", totalAllocation*100)
	}

	return nil
}

// calculateTargetValues calculates target dollar values for each symbol.
func (p *PortfolioRebalancer) calculateTargetValues() {
	for symbol, allocation := range p.config.Allocations {
		p.targetValues[symbol] = p.config.TotalExposure * allocation
	}
}

// monitorLoop continuously monitors and rebalances portfolio.
func (p *PortfolioRebalancer) monitorLoop() {
	ticker := time.NewTicker(p.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.GetContext().Done():
			return
		case <-ticker.C:
			p.checkAndRebalance()
		}
	}
}

// checkAndRebalance analyzes portfolio and rebalances if needed.
func (p *PortfolioRebalancer) checkAndRebalance() {
	// Analyze current state
	allocations := p.analyzePortfolio()

	// Check if rebalancing is needed
	needsRebalance := false
	for _, alloc := range allocations {
		if alloc.NeedsAdjustment {
			needsRebalance = true
			break
		}
	}

	if !needsRebalance {
		p.UpdateMetrics(func(m *OrchestratorMetrics) {
			m.LastOperation = "Portfolio balanced - no action needed"
		})
		return
	}

	// Perform rebalancing
	if err := p.executeRebalance(allocations); err != nil {
		p.IncrementError(fmt.Sprintf("rebalance failed: %v", err))
		return
	}

	p.rebalanceCount++
	p.lastRebalance = time.Now()
	p.IncrementSuccess()
}

// analyzePortfolio calculates current allocations vs targets.
func (p *PortfolioRebalancer) analyzePortfolio() []*SymbolAllocation {
	allocations := make([]*SymbolAllocation, 0)

	// Get all positions
	positions, err := p.sugar.GetOpenPositions()
	if err != nil {
		p.IncrementError(fmt.Sprintf("failed to get positions: %v", err))
		return allocations
	}

	// Calculate current exposure per symbol
	currentExposure := make(map[string]float64)
	totalExposure := 0.0

	for _, pos := range positions {
		// Calculate position value (simplified - should use actual margin/value)
		positionValue := pos.Volume * 100000 // Simplified lot value
		currentExposure[pos.Symbol] += positionValue
		totalExposure += positionValue
	}

	// Analyze each symbol
	for symbol, targetPercent := range p.config.Allocations {
		currentValue := currentExposure[symbol]
		currentPercent := 0.0
		if totalExposure > 0 {
			currentPercent = currentValue / totalExposure
		}

		targetValue := p.targetValues[symbol]
		deviation := currentValue - targetValue
		deviationPercent := 0.0
		if targetValue > 0 {
			deviationPercent = (deviation / targetValue) * 100
		}

		alloc := &SymbolAllocation{
			Symbol:           symbol,
			TargetPercent:    targetPercent * 100,
			CurrentPercent:   currentPercent * 100,
			TargetValue:      targetValue,
			CurrentValue:     currentValue,
			Deviation:        deviation,
			DeviationPercent: deviationPercent,
		}

		// Determine if adjustment needed
		if math.Abs(deviationPercent) > p.config.RebalanceThreshold {
			alloc.NeedsAdjustment = true
			alloc.AdjustmentValue = -deviation

			if deviation < 0 {
				alloc.ActionRequired = "BUY"
			} else {
				alloc.ActionRequired = "SELL"
			}
		} else {
			alloc.ActionRequired = "HOLD"
		}

		allocations = append(allocations, alloc)
	}

	return allocations
}

// executeRebalance performs trades to restore target allocations.
func (p *PortfolioRebalancer) executeRebalance(allocations []*SymbolAllocation) error {
	tradesExecuted := 0

	for _, alloc := range allocations {
		if !alloc.NeedsAdjustment {
			continue
		}

		if tradesExecuted >= p.config.MaxTradesPerCycle {
			p.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.LastOperation = fmt.Sprintf("Max trades per cycle reached (%d)", p.config.MaxTradesPerCycle)
			})
			break
		}

		// Execute adjustment
		if err := p.adjustSymbolExposure(alloc); err != nil {
			p.IncrementError(fmt.Sprintf("failed to adjust %s: %v", alloc.Symbol, err))
			continue
		}

		tradesExecuted++
	}

	p.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("Rebalanced portfolio: %d trades executed", tradesExecuted)
		m.OperationsTotal += tradesExecuted
	})

	return nil
}

// adjustSymbolExposure adjusts exposure for a single symbol.
func (p *PortfolioRebalancer) adjustSymbolExposure(alloc *SymbolAllocation) error {
	// Calculate lot size needed
	// Simplified calculation - should use proper contract size
	lotSize := math.Abs(alloc.AdjustmentValue) / 100000

	// Round to minimum position size
	lotSize = math.Max(lotSize, p.config.MinPositionSize)

	if alloc.ActionRequired == "BUY" {
		// Open buy position
		ticket, err := p.sugar.BuyMarket(alloc.Symbol, lotSize)
		if err != nil {
			return fmt.Errorf("buy failed: %w", err)
		}

		p.UpdateMetrics(func(m *OrchestratorMetrics) {
			m.TotalTrades++
		})

		fmt.Printf("Opened BUY position #%d for %s: %.2f lots\n", ticket, alloc.Symbol, lotSize)

	} else if alloc.ActionRequired == "SELL" {
		// First try to close existing positions
		positions, _ := p.sugar.GetPositionsBySymbol(alloc.Symbol)
		if len(positions) > 0 {
			// Close partial or full position
			if err := p.sugar.ClosePosition(positions[0].Ticket); err != nil {
				return fmt.Errorf("close failed: %w", err)
			}

			p.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.TotalTrades++
			})
		}
	}

	return nil
}

// GetCurrentAllocations returns current portfolio allocation state.
func (p *PortfolioRebalancer) GetCurrentAllocations() []*SymbolAllocation {
	return p.analyzePortfolio()
}

// GetRebalanceCount returns number of times portfolio was rebalanced.
func (p *PortfolioRebalancer) GetRebalanceCount() int {
	return p.rebalanceCount
}

// GetLastRebalanceTime returns when portfolio was last rebalanced.
func (p *PortfolioRebalancer) GetLastRebalanceTime() time.Time {
	return p.lastRebalance
}

/*══════════════════════════════════════════════════════════════════════════════

 ██████╗  ██████╗ ██████╗ ████████╗███████╗ ██████╗ ██╗     ██╗ ██████╗     ██████╗ ███████╗██████╗  █████╗ ██╗      █████╗ ███╗   ██╗ ██████╗███████╗██████╗
 ██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝██╔═══██╗██║     ██║██╔═══██╗    ██╔══██╗██╔════╝██╔══██╗██╔══██╗██║     ██╔══██╗████╗  ██║██╔════╝██╔════╝██╔══██╗
 ██████╔╝██║   ██║██████╔╝   ██║   █████╗  ██║   ██║██║     ██║██║   ██║    ██████╔╝█████╗  ██████╔╝███████║██║     ███████║██╔██╗ ██║██║     █████╗  ██████╔╝
 ██╔═══╝ ██║   ██║██╔══██╗   ██║   ██╔══╝  ██║   ██║██║     ██║██║   ██║    ██╔══██╗██╔══╝  ██╔══██╗██╔══██║██║     ██╔══██║██║╚██╗██║██║     ██╔══╝  ██╔══██╗
 ██║     ╚██████╔╝██║  ██║   ██║   ██║     ╚██████╔╝███████╗██║╚██████╔╝    ██║  ██║███████╗██████╔╝██║  ██║███████╗██║  ██║██║ ╚████║╚██████╗███████╗██║  ██║
 ╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝      ╚═════╝ ╚══════╝╚═╝ ╚═════╝     ╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝╚══════╝╚═╝  ╚═╝

  DETAILED CONFIGURATION GUIDE
  Located at end of file to keep header clean and focused

══════════════════════════════════════════════════════════════════════════════*/

/*
PROGRAMMATIC USAGE & CONFIGURATION

⚙️ PARAMETER CONFIGURATION IS LOCATED IN main.go

WHY THIS SEPARATION EXISTS:
• 15_portfolio_rebalancer.go = REBALANCING ENGINE (orchestrator logic, algorithms)
• main.go → RunOrchestrator_PortfolioRebalancer() = RUNTIME CONFIGURATION (allocations)

THIS SEPARATION IS NEEDED FOR:
1️⃣ Code Reusability
   → Same rebalancer can run with different symbol allocations
   → No need to modify rebalancing logic to change portfolio mix

2️⃣ Quick Testing
   → Want different symbols? Change Allocations in main.go
   → Want stricter rebalancing? Change RebalanceThreshold in main.go
   → Core algorithm remains untouched

3️⃣ User Examples
   → main.go shows HOW to configure portfolio allocations
   → All available parameters and their default values are visible

4️⃣ Centralized Entry Point
   → All orchestrators launch through main.go
   → Single entry point: go run main.go rebalancer → RunOrchestrator_PortfolioRebalancer()

📍 WHERE TO CONFIGURE PARAMETERS:
main.go → func RunOrchestrator_PortfolioRebalancer() (lines 710-776)

CONFIGURATION CODE IN main.go:

func RunOrchestrator_PortfolioRebalancer() error {
    // ... connection code ...

    // ╔════════════════════════════════════════════════════════════╗
    // ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
    // ╚════════════════════════════════════════════════════════════╝
    orchConfig := orchestrators.PortfolioRebalancerConfig{
        Allocations: map[string]float64{
            "EURUSD": 0.25,  // ← 25% EURUSD
            "GBPUSD": 0.25,  // ← 25% GBPUSD
            "USDJPY": 0.25,  // ← 25% USDJPY
            "XAUUSD": 0.25,  // ← 25% XAUUSD (Gold)
        },
        TotalExposure:      10000.0,          // ← $10k total portfolio
        MaxDeviation:       0.20,             // ← 20% max deviation allowed
        RebalanceThreshold: 10.0,             // ← Rebalance at 10% off
        CheckInterval:      30 * time.Minute, // ← Check every 30 min
        MinPositionSize:    0.01,             // ← Minimum lot size
        UseMarketOrders:    true,             // ← Use market orders
        SlippageTolerance:  50,               // ← Max slippage: 50 points
        MaxTradesPerCycle:  10,               // ← Max 10 trades per rebalance
    }

    rebalancer := orchestrators.NewPortfolioRebalancer(sugar, orchConfig)
    rebalancer.Start()
    // ... runs for 1 hour ...
    rebalancer.Stop()
}

💡 EXAMPLE: Different Portfolio Allocations

// Option 1: Equal-Weight Multi-Currency (default in main.go)
Allocations: map[string]float64{
    "EURUSD": 0.25,  // 25% each symbol
    "GBPUSD": 0.25,
    "USDJPY": 0.25,
    "XAUUSD": 0.25,
},
TotalExposure: 10000.0,
RebalanceThreshold: 10.0,  // Rebalance at 10% deviation

// Option 2: Conservative Diversified (modify in main.go)
Allocations: map[string]float64{
    "EURUSD": 0.40,  // ← 40% major currency
    "GBPUSD": 0.20,  // ← 20% secondary currency
    "USDJPY": 0.20,  // ← 20% yen exposure
    "XAUUSD": 0.20,  // ← 20% gold hedge
},
TotalExposure: 10000.0,
RebalanceThreshold: 5.0,   // ← Tighter rebalancing

// Option 3: Aggressive Gold-Heavy (modify in main.go)
Allocations: map[string]float64{
    "XAUUSD": 0.50,  // ← 50% gold (aggressive!)
    "EURUSD": 0.25,  // ← 25% euro
    "USDJPY": 0.25,  // ← 25% yen
},
TotalExposure: 20000.0,   // ← Larger portfolio
RebalanceThreshold: 15.0, // ← Wider tolerance

// Option 4: Multi-Symbol Diversified (modify in main.go)
Allocations: map[string]float64{
    "EURUSD": 0.20,  // ← 20% each
    "GBPUSD": 0.20,
    "USDJPY": 0.20,
    "XAUUSD": 0.20,
    "USDCAD": 0.20,  // ← 5 symbols total
},
TotalExposure: 15000.0,
RebalanceThreshold: 12.0,

📝 IMPORTANT:
• To change parameters → edit main.go, NOT this file
• This file (15_portfolio_rebalancer.go) contains only ORCHESTRATOR LOGIC
• main.go contains CONFIGURATION for specific runs
• Look for the section: ORCHESTRATOR RUNNERS


══════════════════════════════════════════════════════════════════════════════
 WHAT IS PORTFOLIO REBALANCING?
══════════════════════════════════════════════════════════════════════════════

Portfolio rebalancing is an automated strategy that maintains target allocation
percentages across multiple trading symbols. Like a mutual fund manager, it
ensures diversification and prevents over-concentration in any single symbol.

VISUAL EXAMPLE: EQUAL-WEIGHT 4-SYMBOL PORTFOLIO
════════════════════════════════════════════════════════════════════════════

Initial Setup:
──────────────
Portfolio Size: $10,000
Target Allocation: 25% per symbol (4 symbols)

┌────────────────────────────────────────────────────────────────────────┐
│ Symbol  │ Target % │ Target Value │ Current Value │ Deviation │ Action │
├─────────┼──────────┼──────────────┼───────────────┼───────────┼────────┤
│ EURUSD  │   25%    │   $2,500     │    $2,500     │    0%     │  HOLD  │
│ GBPUSD  │   25%    │   $2,500     │    $2,500     │    0%     │  HOLD  │
│ USDJPY  │   25%    │   $2,500     │    $2,500     │    0%     │  HOLD  │
│ XAUUSD  │   25%    │   $2,500     │    $2,500     │    0%     │  HOLD  │
└─────────┴──────────┴──────────────┴───────────────┴───────────┴────────┘
                                                     Portfolio: BALANCED ✅


AFTER 1 HOUR: Market Moves
───────────────────────────
EURUSD rises (+20%), USDJPY falls (-15%), others stable

┌────────────────────────────────────────────────────────────────────────┐
│ Symbol  │ Target % │ Target Value │ Current Value │ Deviation │ Action │
├─────────┼──────────┼──────────────┼───────────────┼───────────┼────────┤
│ EURUSD  │   25%    │   $2,500     │    $3,000     │  +20%     │ SELL ❌│
│ GBPUSD  │   25%    │   $2,500     │    $2,400     │   -4%     │ HOLD ✅│
│ USDJPY  │   25%    │   $2,500     │    $2,125     │  -15%     │ BUY 📈 │
│ XAUUSD  │   25%    │   $2,500     │    $2,475     │   -1%     │ HOLD ✅│
└─────────┴──────────┴──────────────┴───────────────┴───────────┴────────┘
                                                   Portfolio: UNBALANCED ⚠️

TRIGGER: EURUSD deviation (+20%) exceeds RebalanceThreshold (10%)
ACTION:  Rebalancing needed!


REBALANCING EXECUTED:
─────────────────────
1. SELL $500 worth of EURUSD (reduce from $3,000 → $2,500)
2. BUY $375 worth of USDJPY (increase from $2,125 → $2,500)

┌────────────────────────────────────────────────────────────────────────┐
│ Symbol  │ Target % │ Target Value │ Current Value │ Deviation │ Action │
├─────────┼──────────┼──────────────┼───────────────┼───────────┼────────┤
│ EURUSD  │   25%    │   $2,500     │    $2,500     │    0%     │  HOLD  │
│ GBPUSD  │   25%    │   $2,500     │    $2,400     │   -4%     │  HOLD  │
│ USDJPY  │   25%    │   $2,500     │    $2,500     │    0%     │  HOLD  │
│ XAUUSD  │   25%    │   $2,500     │    $2,475     │   -1%     │  HOLD  │
└─────────┴──────────┴──────────────┴───────────────┴───────────┴────────┘
                                                     Portfolio: BALANCED ✅

RESULT: Portfolio restored to target allocations


══════════════════════════════════════════════════════════════════════════════
 HOW PORTFOLIO REBALANCER WORKS (STEP-BY-STEP)
══════════════════════════════════════════════════════════════════════════════

STEP 1: INITIALIZATION (Start)
───────────────────────────────

1. Validates configuration:
   • Do all allocations sum to 100%?
   • Example: 25% + 25% + 25% + 25% = 100% ✅

2. Calculates target dollar values:
   TotalExposure = $10,000
   EURUSD: 25% → $2,500
   GBPUSD: 25% → $2,500
   USDJPY: 25% → $2,500
   XAUUSD: 25% → $2,500

3. Starts monitoring loop (every 30 minutes by default)


STEP 2: MONITORING LOOP (Every CheckInterval)
──────────────────────────────────────────────

Every 30 minutes, the orchestrator:

1. Gets all open positions:
   positions, _ := sugar.GetOpenPositions()

2. Calculates current exposure per symbol:
   For each position:
     currentExposure[symbol] += position.Volume * 100000

3. Analyzes each symbol (analyzePortfolio):
   • Current value vs. target value
   • Deviation percentage
   • Action needed: BUY / SELL / HOLD

4. Checks if rebalancing needed:
   If ANY symbol deviation > RebalanceThreshold:
     → Trigger rebalancing!


STEP 3: REBALANCING EXECUTION (executeRebalance)
─────────────────────────────────────────────────

For each symbol that needs adjustment:

1. Calculate lot size needed:
   lotSize = AdjustmentValue / 100000

2. Execute trade:
   • If ActionRequired == "BUY":
     → sugar.BuyMarket(symbol, lotSize)

   • If ActionRequired == "SELL":
     → sugar.ClosePosition(ticket)

3. Limit: MaxTradesPerCycle = 10 trades per rebalance

4. Update metrics:
   • Increment rebalanceCount
   • Update lastRebalance timestamp


STEP 4: REPEAT
───────────────

Loop continues every CheckInterval (30 min) until Stop() is called.


══════════════════════════════════════════════════════════════════════════════
 PARAMETER DETAILED EXPLANATIONS
══════════════════════════════════════════════════════════════════════════════

• Allocations (map[string]float64)
  Target allocation percentage for each symbol
  Example: map[string]float64{"EURUSD": 0.25, "GBPUSD": 0.25}
  Rules: Must sum to 1.0 (100%)
  Tip: Use equal weights (1.0 / number of symbols) for simplicity

• TotalExposure (float64)
  Total portfolio value to maintain across all symbols
  Example: 10000.0 = $10,000 total exposure
  Calculation: Each symbol gets TotalExposure × Allocation%
  Tip: Set based on account size (e.g., 50% of balance)

• MaxDeviation (float64)
  Maximum allowed deviation from target (as decimal)
  Example: 0.20 = allow up to 20% deviation before intervention
  Currently logged only (future: emergency rebalance trigger)
  Tip: 0.20 (20%) is reasonable, 0.10 (10%) is stricter

• RebalanceThreshold (float64)
  Percentage deviation that triggers rebalancing
  Example: 10.0 = rebalance when symbol is 10% off target
  Trigger: If symbol deviates >10% from target → rebalance
  Tip: Lower = more frequent rebalancing (5-10%), higher = less frequent (15-20%)

• CheckInterval (time.Duration)
  How often to check portfolio balance
  Example: 30 * time.Minute = check every 30 minutes
  Range: 10m = very active, 1h = relaxed, 24h = daily
  Tip: Match to trading style (day trading = 10-30m, swing = 1-4h)

• MinPositionSize (float64)
  Minimum lot size for any position
  Example: 0.01 = 1 micro lot minimum
  Prevents tiny rebalancing trades
  Tip: 0.01 for most brokers, 0.1 for larger accounts

• UseMarketOrders (bool)
  Whether to use market orders (true) or limit orders (false)
  true = immediate execution at market price
  false = limit orders with SlippageTolerance (not yet implemented)
  Tip: Use true for simplicity, false for tighter control

• SlippageTolerance (float64)
  Maximum slippage in points for limit orders
  Example: 50 = allow 50 points (5 pips) slippage
  Currently unused (UseMarketOrders is always true in current implementation)
  Future feature for limit order placement

• MaxTradesPerCycle (int)
  Maximum trades to execute in one rebalancing cycle
  Example: 10 = max 10 trades per rebalance
  Prevents excessive trading costs
  Tip: Set to 2× number of symbols (e.g., 4 symbols → 8 trades)


══════════════════════════════════════════════════════════════════════════════
 EXAMPLE PORTFOLIO CONFIGURATIONS
══════════════════════════════════════════════════════════════════════════════

╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 1: CONSERVATIVE DIVERSIFIED (4 SYMBOLS)                           ║
╚═══════════════════════════════════════════════════════════════════════════╝

PortfolioRebalancerConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.40,  // 40% major currency pair
        "GBPUSD": 0.20,  // 20% secondary currency
        "USDJPY": 0.20,  // 20% yen exposure
        "XAUUSD": 0.20,  // 20% gold as hedge
    },
    TotalExposure:      5000.0,           // Small portfolio
    RebalanceThreshold: 5.0,              // Tight rebalancing
    CheckInterval:      1 * time.Hour,    // Check hourly
    MinPositionSize:    0.01,
    UseMarketOrders:    true,
    MaxTradesPerCycle:  8,
}

BEST FOR:
• Beginners learning portfolio management
• Small accounts ($5,000-$10,000)
• Risk-averse traders
• Long-term position holding

EXPECTED BEHAVIOR:
• Frequent rebalancing (5% threshold)
• Low individual symbol risk
• Gold provides downside protection


╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 2: BALANCED MULTI-CURRENCY (EQUAL WEIGHT)                         ║
╚═══════════════════════════════════════════════════════════════════════════╝

PortfolioRebalancerConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.25,  // 25% each
        "GBPUSD": 0.25,
        "USDJPY": 0.25,
        "XAUUSD": 0.25,
    },
    TotalExposure:      10000.0,          // Medium portfolio
    RebalanceThreshold: 10.0,             // Moderate rebalancing
    CheckInterval:      30 * time.Minute, // Check every 30 min
    MinPositionSize:    0.01,
    UseMarketOrders:    true,
    MaxTradesPerCycle:  10,
}

BEST FOR:
• Intermediate traders
• Medium accounts ($10,000-$50,000)
• Diversified exposure
• Automated management

EXPECTED BEHAVIOR:
• Balanced across all symbols
• Moderate rebalancing frequency
• Good for trending markets


╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 3: AGGRESSIVE COMMODITY-HEAVY                                     ║
╚═══════════════════════════════════════════════════════════════════════════╝

PortfolioRebalancerConfig{
    Allocations: map[string]float64{
        "XAUUSD": 0.50,  // 50% gold (aggressive!)
        "EURUSD": 0.25,  // 25% euro
        "USDJPY": 0.25,  // 25% yen
    },
    TotalExposure:      20000.0,          // Larger portfolio
    RebalanceThreshold: 15.0,             // Wider tolerance
    CheckInterval:      1 * time.Hour,    // Hourly checks
    MinPositionSize:    0.05,             // Larger minimum
    UseMarketOrders:    true,
    MaxTradesPerCycle:  6,
}

BEST FOR:
• Experienced traders
• Large accounts ($20,000+)
• Gold trading specialists
• High volatility tolerance

EXPECTED BEHAVIOR:
• Heavy gold exposure (50%)
• Less frequent rebalancing (15% threshold)
• Higher profit/loss volatility


╔═══════════════════════════════════════════════════════════════════════════╗
║ PROFILE 4: MULTI-SYMBOL DIVERSIFIED (5+ SYMBOLS)                          ║
╚═══════════════════════════════════════════════════════════════════════════╝

PortfolioRebalancerConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.20,   // 20% each
        "GBPUSD": 0.20,
        "USDJPY": 0.20,
        "XAUUSD": 0.20,
        "USDCAD": 0.20,   // 5 symbols total
    },
    TotalExposure:      15000.0,
    RebalanceThreshold: 12.0,
    CheckInterval:      45 * time.Minute,
    MinPositionSize:    0.01,
    UseMarketOrders:    true,
    MaxTradesPerCycle:  12,
}

BEST FOR:
• Portfolio managers
• Diversification seekers
• Medium-large accounts
• Multiple market exposure

EXPECTED BEHAVIOR:
• Maximum diversification
• Lower correlation risk
• Moderate rebalancing frequency


══════════════════════════════════════════════════════════════════════════════
 RISK WARNINGS & BEST PRACTICES
══════════════════════════════════════════════════════════════════════════════

⚠️ CRITICAL WARNINGS:

1. ALLOCATION MUST SUM TO 100%
   → If allocations don't sum to 1.0 (100%), Start() will fail
   → Validation: 0.25 + 0.25 + 0.25 + 0.25 = 1.0 ✅
   → Always verify your allocation percentages

2. OVER-REBALANCING COSTS
   → Frequent rebalancing = more trades = more spreads/commissions
   → Set RebalanceThreshold appropriately (10-15% recommended)
   → Monitor trading costs vs. portfolio performance

3. TRENDING MARKETS
   → Rebalancing cuts winners and adds to losers
   → In strong trends, this can reduce profits
   → Consider wider RebalanceThreshold (15-20%) in trending markets

4. EXPOSURE CALCULATION IS SIMPLIFIED
   → Current implementation: pos.Volume × 100000
   → Doesn't account for leverage, margin, or contract size differences
   → Test on DEMO first to verify calculations for your symbols

5. MARKET HOURS MATTER
   → Rebalancing during low liquidity = wider spreads
   → Avoid rebalancing during Asian session (low volume)
   → Best: European/US session overlap (8am-12pm EST)


💡 BEST PRACTICES:

✅ Start with Equal-Weight Allocation
   → Simplest configuration: 1.0 / number of symbols
   → Example: 4 symbols → 0.25 each
   → Easy to understand and manage

✅ Test on DEMO Account First
   → Run for 1 week on demo
   → Monitor rebalancing frequency
   → Verify exposure calculations
   → Adjust RebalanceThreshold if needed

✅ Set Realistic TotalExposure
   → Don't use full account balance
   → Recommended: 30-50% of account balance
   → Leave buffer for drawdowns and margin

✅ Match CheckInterval to Trading Style
   → Day trading: 10-30 minutes
   → Swing trading: 1-4 hours
   → Position trading: 12-24 hours

✅ Monitor Rebalancing Costs
   → Track: Number of rebalances × average spread
   → If costs > 2% of TotalExposure → increase RebalanceThreshold
   → Balance: Fewer rebalances vs. tighter balance

✅ Use Correlated Symbols Carefully
   → EURUSD + GBPUSD are highly correlated
   → Both may move together → minimal rebalancing benefit
   → Better: Mix currencies + commodities (EURUSD + XAUUSD)

✅ Regular Portfolio Review
   → Check GetCurrentAllocations() daily
   → Look for patterns (one symbol always over-allocated?)
   → Adjust target allocations if market regime changes

✅ Risk Management Integration
   → Combine with RiskManager orchestrator
   → Portfolio Rebalancer manages allocation
   → RiskManager manages total drawdown
   → Together = complete risk control


══════════════════════════════════════════════════════════════════════════════
 WHEN TO USE PORTFOLIO REBALANCER
══════════════════════════════════════════════════════════════════════════════

✅ GOOD USE CASES:

1. Multi-Symbol Trading Systems
   → You trade 4-10 symbols simultaneously
   → Want to maintain balanced exposure
   → Prevents over-concentration in one symbol

2. Diversification Strategy
   → Reduce correlation risk
   → Mix currencies + commodities + indices
   → Smooth out equity curve

3. Automated Portfolio Management
   → Set-and-forget approach
   → Automatic rebalancing without manual intervention
   → Good for hands-off traders

4. Risk Control Across Symbols
   → Limit exposure to any single symbol
   → Prevent one winning trade from dominating portfolio
   → Systematic risk management


❌ BAD USE CASES:

1. Single-Symbol Trading
   → If you only trade one symbol → no need for rebalancer
   → Use other orchestrators (Trailing Stop, Risk Manager)

2. High-Frequency Trading
   → Rebalancing overhead too high for HFT
   → Spreads/commissions eat profits
   → Better: Manual portfolio management

3. Strong Trending Markets
   → Rebalancing cuts winners too early
   → Adds to losers prematurely
   → Better: Let winners run, use Trailing Stop

4. Very Small Accounts
   → TotalExposure < $1,000 → rebalancing costs too high
   → Minimum position sizes may prevent proper allocation
   → Better: Focus on single-symbol strategies


══════════════════════════════════════════════════════════════════════════════
 TROUBLESHOOTING COMMON ISSUES
══════════════════════════════════════════════════════════════════════════════

❌ PROBLEM: "allocations must sum to 100%, got X%"
✅ SOLUTION: Check your Allocations map
   map[string]float64{
       "EURUSD": 0.25,
       "GBPUSD": 0.25,
       "USDJPY": 0.25,
       "XAUUSD": 0.25,  // Total = 1.0 ✅
   }

❌ PROBLEM: Portfolio never rebalances
✅ SOLUTION: Check if deviations exceed RebalanceThreshold
   • If threshold = 15%, but deviations are only 8% → no rebalance
   • Lower RebalanceThreshold (e.g., 5-10%)
   • Or wait for larger market moves

❌ PROBLEM: Too frequent rebalancing
✅ SOLUTION: Increase RebalanceThreshold
   • If rebalancing every 30 minutes → too frequent
   • Increase threshold from 10% to 15-20%
   • Increase CheckInterval (30min → 1 hour)

❌ PROBLEM: Exposure calculations seem wrong
✅ SOLUTION: Verify symbol contract sizes
   • Current calculation: Volume × 100000 (assumes forex)
   • Gold (XAUUSD) may need different calculation
   • Test on demo, verify with manual calculations

❌ PROBLEM: MaxTradesPerCycle reached too often
✅ SOLUTION: Increase MaxTradesPerCycle
   • If you have 5 symbols, set to 10+ trades per cycle
   • Rule of thumb: 2× number of symbols
   • Or reduce number of symbols

═══════════════════════════════════════════════════════════════════════════*/
