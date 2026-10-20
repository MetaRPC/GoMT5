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
   • Resizes the existing positions to restore balance (it never opens
     a symbol that has no position)

 HOW TO RUN (Command-Line):
   Step 1: Open terminal/command prompt
   Step 2: Navigate to: examples/demos/
   Step 3: Run one of these commands:

      go run main.go 15           ← Run by number (Portfolio Rebalancer is #15)
      go run main.go rebalancer   ← Run by short name
      go run main.go portfolio    ← Run by full name

   What happens:
   • Connects to your MT5 account
   • Checks the portfolio balance at once, then every 30 minutes
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
       (partial close of the long EURUSD position)
     • GBPUSD: Increase by $400 (from $2,100 to $2,500) → BUY
     • USDJPY: OK, no action needed
     • XAUUSD: Increase by $400 (from $2,100 to $2,500) → BUY
     (all positions long in this example; for a short position an
      increase is a SELL and a reduction a BUY)

  3. EXECUTES trades to restore balance

  4. RESULT: Portfolio is balanced again at 25% each!

HOW IT WORKS IN PRACTICE:
──────────────────────────
  • Every 30 minutes, checks if portfolio is balanced
  • If deviation > 10%, triggers rebalancing
  • Automatically resizes the existing positions to restore target %
    (increase in the position's direction, or partial close)
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
	"sort"
	"sync"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// PortfolioRebalancerConfig holds portfolio rebalancing parameters.
type PortfolioRebalancerConfig struct {
	// Portfolio Definition
	Allocations map[string]float64 // Symbol -> target allocation (0.0-1.0, all together 1.0)

	// Exposure Limits
	TotalExposure   float64 // Total notional value (account currency) the positions should add up to
	MaxDeviation    float64 // Deviation (0.20 = 20%) that triggers a rebalance at once, between the regular checks (0 = off)

	// Rebalancing Rules
	RebalanceThreshold float64       // % deviation to trigger a rebalance at the regular check (10.0 = 10%)
	CheckInterval      time.Duration // How often to check balance (0 = 1h)
	MinPositionSize    float64       // Trades smaller than this (lots) are skipped

	// Trading Parameters
	UseMarketOrders   bool    // true: market orders; false: limit orders SlippageTolerance points better than the price
	SlippageTolerance float64 // Points: max deviation of market orders / price improvement of limit orders
	MaxTradesPerCycle int     // Max trades per rebalancing cycle (0 = 2 × number of symbols)
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
//
// It rebalances the positions you have: a symbol's exposure is the NET notional value of its
// positions (BUY +, SELL -) in the account currency, and the rebalancer resizes it toward
// TotalExposure × allocation in the direction it already has. It never opens a symbol that
// has no position (it would have to guess the direction).
type PortfolioRebalancer struct {
	*BaseOrchestrator
	sugar  *mt5.MT5Sugar
	config PortfolioRebalancerConfig

	// Used by the monitoring goroutine only
	targetValues    map[string]float64          // Target value per symbol
	symbols         map[string]*rebalanceSymbol // Symbol specification and value conversion
	accountCurrency string
	limitOrders     []uint64 // Unfilled limit orders of the last cycle (UseMarketOrders=false)

	// Shared with the getters (mu)
	mu                 sync.Mutex
	currentAllocations []*SymbolAllocation // Result of the last analysis
	lastRebalance      time.Time
	rebalanceCount     int
}

// SymbolAllocation represents current vs target allocation for a symbol.
type SymbolAllocation struct {
	Symbol          string
	TargetPercent   float64
	CurrentPercent  float64
	TargetValue     float64
	CurrentValue    float64 // |net notional value| in the account currency
	Deviation       float64
	DeviationPercent float64
	NeedsAdjustment bool
	ActionRequired  string // "BUY", "SELL", "HOLD"
	AdjustmentValue float64
	IsLong          bool    // Direction of the net position
	ValuePerLot     float64 // Notional value of 1 lot in the account currency
	Price           float64 // Price the values were computed with
}

// rebalanceSymbol is what the rebalancer needs to value and trade a symbol.
type rebalanceSymbol struct {
	spec         symbolSpec
	profitIsAcct bool    // profit (quote) currency = account currency (EURUSD, XAUUSD on a USD account)
	baseIsAcct   bool    // base currency = account currency (USDJPY on a USD account)
	tickValue    float64 // SYMBOL_TRADE_TICK_VALUE (conversion for crosses)
	tickSize     float64 // SYMBOL_TRADE_TICK_SIZE
}

// valuePerLot is the notional value of 1 lot at price, in the account currency.
func (s *rebalanceSymbol) valuePerLot(price float64) float64 {
	switch {
	case s.profitIsAcct:
		return s.spec.Contract * price
	case s.baseIsAcct:
		return s.spec.Contract
	case s.tickSize > 0 && s.tickValue > 0:
		return price * s.tickValue / s.tickSize // contract × price × (account per quote unit)
	}
	return s.spec.Contract * price // no conversion available: quote currency
}

// NewPortfolioRebalancer creates a new portfolio rebalancing orchestrator.
func NewPortfolioRebalancer(sugar *mt5.MT5Sugar, config PortfolioRebalancerConfig) *PortfolioRebalancer {
	return &PortfolioRebalancer{
		BaseOrchestrator:   NewBaseOrchestrator("Portfolio Rebalancer"),
		sugar:              sugar,
		config:             config,
		targetValues:       make(map[string]float64),
		symbols:            make(map[string]*rebalanceSymbol),
		currentAllocations: make([]*SymbolAllocation, 0),
	}
}

// Start validates the configuration and begins portfolio monitoring and rebalancing.
// The first check runs at once, then every CheckInterval.
func (p *PortfolioRebalancer) Start() error {
	if p.IsRunning() {
		return fmt.Errorf("portfolio rebalancer already running")
	}

	// Validate configuration (applies the defaults)
	if err := p.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Account currency (exposure is valued in it)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	currency, err := p.sugar.GetService().GetAccountString(ctx, pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY)
	cancel()
	if err != nil {
		return fmt.Errorf("failed to get account currency: %w", err)
	}
	p.accountCurrency = currency

	// Calculate target values
	p.calculateTargetValues()

	// Create context, mark as started, stop on Ctrl+C
	p.begin()

	// Start monitoring loop
	p.goLoop(p.monitorLoop)

	return nil
}

// Stop gracefully stops portfolio rebalancing and deletes its unfilled limit orders.
// It waits until a check that is in progress has finished.
func (p *PortfolioRebalancer) Stop() error {
	if !p.IsRunning() {
		return fmt.Errorf("portfolio rebalancer not running")
	}

	// Cancel context and wait for the monitoring loop
	p.stopLoops()

	// Delete limit orders that did not fill
	p.deleteUnfilledLimitOrders()

	// Mark as stopped
	p.MarkStopped()

	return nil
}

// validateConfig ensures configuration is valid and applies the defaults.
func (p *PortfolioRebalancer) validateConfig() error {
	if len(p.config.Allocations) == 0 {
		return fmt.Errorf("no allocations")
	}

	// Check allocations sum to 100%
	totalAllocation := 0.0
	for symbol, allocation := range p.config.Allocations {
		if allocation <= 0 || allocation > 1 {
			return fmt.Errorf("allocation of %s must be in (0, 1], got %v", symbol, allocation)
		}
		totalAllocation += allocation
	}

	if math.Abs(totalAllocation-1.0) > 0.001 {
		return fmt.Errorf("allocations must sum to 100%%, got %.2f%%", totalAllocation*100)
	}

	if p.config.TotalExposure <= 0 {
		return fmt.Errorf("TotalExposure must be > 0, got %v", p.config.TotalExposure)
	}
	if p.config.RebalanceThreshold < 0 || p.config.MaxDeviation < 0 || p.config.MinPositionSize < 0 ||
		p.config.SlippageTolerance < 0 || p.config.MaxTradesPerCycle < 0 {
		return fmt.Errorf("RebalanceThreshold, MaxDeviation, MinPositionSize, SlippageTolerance and MaxTradesPerCycle must not be negative")
	}
	if p.config.MaxDeviation > 0 && p.config.MaxDeviation*100 < p.config.RebalanceThreshold {
		return fmt.Errorf("MaxDeviation %.0f%% is below RebalanceThreshold %.0f%%", p.config.MaxDeviation*100, p.config.RebalanceThreshold)
	}
	if !p.config.UseMarketOrders && p.config.SlippageTolerance <= 0 {
		return fmt.Errorf("limit orders (UseMarketOrders=false) need SlippageTolerance > 0 points")
	}

	// Defaults
	if p.config.CheckInterval <= 0 {
		p.config.CheckInterval = 1 * time.Hour
	}
	if p.config.MaxTradesPerCycle == 0 {
		p.config.MaxTradesPerCycle = 2 * len(p.config.Allocations)
	}

	return nil
}

// calculateTargetValues calculates target values (account currency) for each symbol.
func (p *PortfolioRebalancer) calculateTargetValues() {
	for symbol, allocation := range p.config.Allocations {
		p.targetValues[symbol] = p.config.TotalExposure * allocation
	}
}

// watchInterval is how often the MaxDeviation trigger is checked between the regular checks
// (only when CheckInterval is longer).
const watchInterval = time.Minute

// monitorLoop checks at once, then rebalances every CheckInterval; in between it rebalances
// at once when a symbol deviates more than MaxDeviation.
func (p *PortfolioRebalancer) monitorLoop() {
	ticker := time.NewTicker(p.config.CheckInterval)
	defer ticker.Stop()

	var watch <-chan time.Time
	if p.config.MaxDeviation > 0 && p.config.CheckInterval > watchInterval {
		watchTicker := time.NewTicker(watchInterval)
		defer watchTicker.Stop()
		watch = watchTicker.C
	}

	p.checkAndRebalance(false)

	for {
		select {
		case <-p.GetContext().Done():
			return
		case <-ticker.C:
			p.checkAndRebalance(false)
		case <-watch:
			p.checkAndRebalance(true)
		}
	}
}

// checkAndRebalance analyzes portfolio and rebalances if needed. A regular check rebalances when a
// symbol deviates more than RebalanceThreshold; a watch check only when one deviates more than MaxDeviation.
func (p *PortfolioRebalancer) checkAndRebalance(watchOnly bool) {
	if p.GetContext().Err() != nil {
		return
	}

	// Limit orders of the previous cycle that did not fill are replaced
	p.deleteUnfilledLimitOrders()

	// Analyze current state
	allocations, err := p.analyzePortfolio()
	if err != nil {
		p.IncrementError(fmt.Sprintf("analysis failed: %v", err))
		return
	}

	// Check if rebalancing is needed
	needsRebalance := false
	for _, alloc := range allocations {
		if alloc.NeedsAdjustment && (!watchOnly || math.Abs(alloc.DeviationPercent) > p.config.MaxDeviation*100) {
			needsRebalance = true
			break
		}
	}

	if !needsRebalance {
		if !watchOnly {
			p.UpdateMetrics(func(m *OrchestratorMetrics) {
				m.LastOperation = "Portfolio balanced - no action needed"
			})
		}
		return
	}

	if watchOnly {
		fmt.Printf("\n  ⚠️  Deviation above MaxDeviation (%.0f%%): rebalancing now\n", p.config.MaxDeviation*100)
	}

	// Perform rebalancing
	trades := p.executeRebalance(allocations)
	if trades == 0 {
		return
	}

	p.mu.Lock()
	p.rebalanceCount++
	p.lastRebalance = time.Now()
	p.mu.Unlock()
	p.IncrementSuccess()
}

// getSymbol gets (and caches) a symbol's specification and value conversion.
func (p *PortfolioRebalancer) getSymbol(symbol string) (*rebalanceSymbol, error) {
	if s, ok := p.symbols[symbol]; ok {
		return s, nil
	}
	spec, err := loadSymbolSpec(p.sugar, symbol)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	service := p.sugar.GetService()
	base, err := service.GetSymbolString(ctx, symbol, pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_BASE)
	if err != nil {
		return nil, err
	}
	profit, err := service.GetSymbolString(ctx, symbol, pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_PROFIT)
	if err != nil {
		return nil, err
	}
	tickValue, _ := service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE)
	tickSize, _ := service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_SIZE)

	s := &rebalanceSymbol{
		spec:         spec,
		profitIsAcct: profit == p.accountCurrency,
		baseIsAcct:   base == p.accountCurrency,
		tickValue:    tickValue,
		tickSize:     tickSize,
	}
	p.symbols[symbol] = s
	return s, nil
}

// analyzePortfolio calculates current allocations vs targets and stores the result for GetCurrentAllocations.
func (p *PortfolioRebalancer) analyzePortfolio() ([]*SymbolAllocation, error) {
	// Get all positions
	positions, err := p.sugar.GetOpenPositions()
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// Net exposure per symbol: Volume × value of 1 lot, BUY positive, SELL negative
	netExposure := make(map[string]float64)
	prices := make(map[string]float64)
	for _, pos := range positions {
		if _, managed := p.config.Allocations[pos.Symbol]; !managed {
			continue
		}
		s, err := p.getSymbol(pos.Symbol)
		if err != nil {
			return nil, err
		}
		value := pos.Volume * s.valuePerLot(pos.PriceCurrent)
		if pos.Type != pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY {
			value = -value
		}
		netExposure[pos.Symbol] += value
		prices[pos.Symbol] = pos.PriceCurrent
	}

	totalExposure := 0.0
	for _, value := range netExposure {
		totalExposure += math.Abs(value)
	}

	// Analyze each symbol (sorted for a stable order)
	symbols := make([]string, 0, len(p.config.Allocations))
	for symbol := range p.config.Allocations {
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)

	allocations := make([]*SymbolAllocation, 0, len(symbols))
	for _, symbol := range symbols {
		targetPercent := p.config.Allocations[symbol]
		net := netExposure[symbol]
		currentValue := math.Abs(net)
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
			IsLong:           net > 0,
			ActionRequired:   "HOLD",
			Price:            prices[symbol],
		}

		// A symbol without a (net) position is not opened: the direction is unknown
		if currentValue == 0 {
			alloc.DeviationPercent = 0
			allocations = append(allocations, alloc)
			continue
		}
		if s, err := p.getSymbol(symbol); err == nil {
			alloc.ValuePerLot = s.valuePerLot(alloc.Price)
		}

		// Determine if adjustment needed: increase or reduce the position in its own direction
		if math.Abs(deviationPercent) > p.config.RebalanceThreshold {
			alloc.NeedsAdjustment = true
			alloc.AdjustmentValue = -deviation

			increase := deviation < 0
			if increase == alloc.IsLong {
				alloc.ActionRequired = "BUY"
			} else {
				alloc.ActionRequired = "SELL"
			}
		}

		allocations = append(allocations, alloc)
	}

	p.mu.Lock()
	p.currentAllocations = allocations
	p.mu.Unlock()

	return allocations, nil
}

// executeRebalance performs trades to restore target allocations and returns how many it made.
func (p *PortfolioRebalancer) executeRebalance(allocations []*SymbolAllocation) int {
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
		trades, err := p.adjustSymbolExposure(alloc, p.config.MaxTradesPerCycle-tradesExecuted)
		tradesExecuted += trades
		if err != nil {
			p.IncrementError(fmt.Sprintf("failed to adjust %s: %v", alloc.Symbol, err))
		}
	}

	p.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.LastOperation = fmt.Sprintf("Rebalanced portfolio: %d trades executed", tradesExecuted)
		m.OperationsTotal += tradesExecuted
		m.TotalTrades += tradesExecuted
	})

	return tradesExecuted
}

// adjustSymbolExposure resizes one symbol's position toward its target and returns the number of trades.
func (p *PortfolioRebalancer) adjustSymbolExposure(alloc *SymbolAllocation, maxTrades int) (int, error) {
	s, err := p.getSymbol(alloc.Symbol)
	if err != nil {
		return 0, err
	}
	if alloc.ValuePerLot <= 0 {
		return 0, fmt.Errorf("no value per lot")
	}

	// Lots needed for the adjustment, rounded to the NEAREST volume step (+ half a step, then down).
	// Rounded down, a reduction could leave the symbol further from its target than one step more:
	// EURUSD 0.10 lot ($11,225), target $2,500 → 0.0777 lots to close; 0.07 left 0.03 lot (+35%,
	// above MaxDeviation, so the minute check fired again every minute with no trade possible),
	// 0.08 leaves 0.02 lot (-10%).
	step := s.spec.VolumeStep
	if step <= 0 {
		step = 0.01
	}
	lotSize := s.spec.volume(math.Abs(alloc.AdjustmentValue)/alloc.ValuePerLot + step/2)
	if lotSize == 0 || lotSize < p.config.MinPositionSize {
		p.UpdateMetrics(func(m *OrchestratorMetrics) {
			m.LastOperation = fmt.Sprintf("%s: adjustment below MinPositionSize/volume step, skipped", alloc.Symbol)
		})
		return 0, nil
	}

	increase := alloc.Deviation < 0
	if increase {
		ticket, err := p.openTrade(alloc.Symbol, s, alloc.IsLong, lotSize)
		if err != nil {
			return 0, err
		}
		kind := "market"
		if !p.config.UseMarketOrders {
			kind = "limit order"
		}
		fmt.Printf("\n  [REBALANCE] %s %s %.2f lots (%s, #%d): exposure %.0f → target %.0f\n",
			alloc.ActionRequired, alloc.Symbol, lotSize, kind, ticket, alloc.CurrentValue, alloc.TargetValue)
		return 1, nil
	}

	return p.reducePositions(alloc, s, lotSize, maxTrades)
}

// openTrade increases the exposure: a market order with SlippageTolerance as the maximum deviation,
// or (UseMarketOrders=false) a limit order SlippageTolerance points better than the current price.
func (p *PortfolioRebalancer) openTrade(symbol string, s *rebalanceSymbol, buy bool, lots float64) (uint64, error) {
	if !p.config.UseMarketOrders {
		priceInfo, err := p.sugar.GetPriceInfo(symbol)
		if err != nil {
			return 0, err
		}
		offset := p.config.SlippageTolerance * s.spec.Point
		var ticket uint64
		if buy {
			ticket, err = p.sugar.BuyLimit(symbol, lots, s.spec.price(priceInfo.Ask-offset))
		} else {
			ticket, err = p.sugar.SellLimit(symbol, lots, s.spec.price(priceInfo.Bid+offset))
		}
		if err != nil {
			return 0, err
		}
		p.limitOrders = append(p.limitOrders, ticket)
		return ticket, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	operation := pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL
	if buy {
		operation = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY
	}
	slippage := uint64(p.config.SlippageTolerance)
	result, err := p.sugar.GetService().PlaceOrder(ctx, &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: operation,
		Volume:    lots,
		Slippage:  &slippage,
	})
	if err != nil {
		return 0, err
	}
	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}
	return result.Order, nil
}

// reducePositions closes lots from the symbol's positions in the net direction, largest first
// (partially where possible; a remainder below the minimum volume closes the whole position).
func (p *PortfolioRebalancer) reducePositions(alloc *SymbolAllocation, s *rebalanceSymbol, lots float64, maxTrades int) (int, error) {
	positions, err := p.sugar.GetPositionsBySymbol(alloc.Symbol)
	if err != nil {
		return 0, err
	}
	sameSide := make([]*pb.PositionInfo, 0, len(positions))
	for _, pos := range positions {
		if (pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY) == alloc.IsLong {
			sameSide = append(sameSide, pos)
		}
	}
	sort.Slice(sameSide, func(i, j int) bool { return sameSide[i].Volume > sameSide[j].Volume })

	trades := 0
	remaining := lots
	for _, pos := range sameSide {
		if remaining < s.spec.VolumeMin-1e-9 || trades >= maxTrades {
			break
		}
		closeVolume := s.spec.volume(math.Min(remaining, pos.Volume))
		if closeVolume == 0 {
			break
		}
		if pos.Volume-closeVolume < s.spec.VolumeMin-1e-9 {
			closeVolume = pos.Volume
			err = p.sugar.ClosePosition(pos.Ticket)
		} else {
			err = p.sugar.ClosePositionPartial(pos.Ticket, closeVolume)
		}
		if err != nil {
			return trades, fmt.Errorf("close #%d: %w", pos.Ticket, err)
		}
		trades++
		remaining -= closeVolume
		fmt.Printf("\n  [REBALANCE] %s %s: closed %.2f of %.2f lots of #%d: exposure %.0f → target %.0f\n",
			alloc.ActionRequired, alloc.Symbol, closeVolume, pos.Volume, pos.Ticket, alloc.CurrentValue, alloc.TargetValue)
	}

	return trades, nil
}

// deleteUnfilledLimitOrders deletes the rebalancer's limit orders that are still pending
// (a filled order's ticket is the position's ticket, so filled ones are left alone).
func (p *PortfolioRebalancer) deleteUnfilledLimitOrders() {
	if len(p.limitOrders) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	service := p.sugar.GetService()
	_, orderTickets, err := service.GetOpenedTickets(ctx)
	if err != nil {
		p.IncrementError(fmt.Sprintf("failed to get pending orders: %v", err))
		return
	}
	pending := make(map[uint64]bool, len(orderTickets))
	for _, t := range orderTickets {
		pending[uint64(t)] = true
	}
	for _, ticket := range p.limitOrders {
		if !pending[ticket] {
			continue
		}
		if code, err := service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: ticket}); err != nil || code != 10009 {
			p.IncrementError(fmt.Sprintf("failed to delete limit order #%d: code %d %v", ticket, code, err))
		}
	}
	p.limitOrders = nil
}

// GetCurrentAllocations returns the result of the last portfolio analysis
// (empty until the first check has run).
func (p *PortfolioRebalancer) GetCurrentAllocations() []*SymbolAllocation {
	p.mu.Lock()
	defer p.mu.Unlock()

	allocations := make([]*SymbolAllocation, len(p.currentAllocations))
	for i, alloc := range p.currentAllocations {
		a := *alloc
		allocations[i] = &a
	}
	return allocations
}

// GetRebalanceCount returns number of times portfolio was rebalanced.
func (p *PortfolioRebalancer) GetRebalanceCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.rebalanceCount
}

// GetLastRebalanceTime returns when portfolio was last rebalanced (zero time: never).
func (p *PortfolioRebalancer) GetLastRebalanceTime() time.Time {
	p.mu.Lock()
	defer p.mu.Unlock()
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
main.go → func RunOrchestrator_PortfolioRebalancer() (CONFIGURATION block)

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
        MaxDeviation:       0.20,             // ← 20% off → rebalance at once
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
   • Do all allocations sum to 100%? (each in (0, 1])
   • Example: 25% + 25% + 25% + 25% = 100% ✅
   • TotalExposure > 0, no negative values, MaxDeviation ≥ RebalanceThreshold
   • Defaults: CheckInterval 0 → 1 hour, MaxTradesPerCycle 0 → 2 × symbols

2. Reads the account currency (exposure is valued in it) and calculates
   target values:
   TotalExposure = $10,000
   EURUSD: 25% → $2,500
   GBPUSD: 25% → $2,500
   USDJPY: 25% → $2,500
   XAUUSD: 25% → $2,500

3. Starts monitoring loop: first check at once, then every CheckInterval;
   in between, every minute, a rebalance when a symbol deviates more than
   MaxDeviation


STEP 2: MONITORING LOOP (Every CheckInterval)
──────────────────────────────────────────────

1. Gets all open positions:
   positions, _ := sugar.GetOpenPositions()

2. Calculates the NET exposure per symbol in the account currency:
   For each position (BUY +, SELL −):
     value = Volume × value of 1 lot
     value of 1 lot = ContractSize × price            (EURUSD, XAUUSD on USD)
                    = ContractSize                    (USDJPY on USD)
                    = price × TickValue / TickSize    (crosses)

3. Analyzes each symbol (analyzePortfolio):
   • |net exposure| vs. target value
   • Deviation percentage
   • A symbol WITHOUT a position is not opened (direction unknown)
   • Action: increase or reduce the position in its own direction
     (BUY / SELL), or HOLD

4. Checks if rebalancing needed:
   If ANY symbol deviation > RebalanceThreshold:
     → Trigger rebalancing!


STEP 3: REBALANCING EXECUTION (executeRebalance)
─────────────────────────────────────────────────

For each symbol that needs adjustment:

1. Calculate lot size needed:
   lotSize = |AdjustmentValue| / value of 1 lot,
   rounded to the nearest volume step (skipped if below MinPositionSize)

2. Execute trade:
   • Increase: market order in the position's direction with
     SlippageTolerance as max deviation (UseMarketOrders = false: a limit
     order SlippageTolerance points better than the price, deleted at the
     next check if unfilled)

   • Reduce: partial close of the positions in that direction, largest
     first (sugar.ClosePositionPartial / ClosePosition)

3. Limit: MaxTradesPerCycle = 10 trades per rebalance

4. Update metrics:
   • Increment rebalanceCount
   • Update lastRebalance timestamp


STEP 4: REPEAT
───────────────

Loop continues every CheckInterval (30 min) until Stop() or Ctrl+C.


══════════════════════════════════════════════════════════════════════════════
 PARAMETER DETAILED EXPLANATIONS
══════════════════════════════════════════════════════════════════════════════

• Allocations (map[string]float64)
  Target allocation percentage for each symbol
  Example: map[string]float64{"EURUSD": 0.25, "GBPUSD": 0.25}
  Rules: Must sum to 1.0 (100%)
  Tip: Use equal weights (1.0 / number of symbols) for simplicity

• TotalExposure (float64)
  Total notional value (account currency) the positions should add up to
  Example: 10000.0 = $10,000 total exposure
  Calculation: Each symbol gets TotalExposure × Allocation%
  Tip: 1 lot EURUSD ≈ $110,000 notional - $2,500 is about 0.02 lots

• MaxDeviation (float64)
  Deviation (as decimal) that triggers a rebalance AT ONCE, without waiting
  for the next CheckInterval (checked every minute)
  Example: 0.20 = rebalance immediately when a symbol is 20% off target
  0 = off; must not be below RebalanceThreshold
  Tip: 0.20 (20%) is reasonable

• RebalanceThreshold (float64)
  Percentage deviation that triggers rebalancing at the regular check
  Example: 10.0 = rebalance when symbol is 10% off target
  Trigger: If symbol deviates >10% from target → rebalance
  Tip: Lower = more frequent rebalancing (5-10%), higher = less frequent (15-20%)

• CheckInterval (time.Duration)
  How often to check portfolio balance (0 = 1 hour)
  Example: 30 * time.Minute = check every 30 minutes
  Range: 10m = very active, 1h = relaxed, 24h = daily
  Tip: Match to trading style (day trading = 10-30m, swing = 1-4h)

• MinPositionSize (float64)
  Smallest trade (lots) the rebalancer makes
  Example: 0.01 = 1 micro lot minimum
  Smaller adjustments are skipped (prevents tiny rebalancing trades)
  Tip: 0.01 for most brokers, 0.1 for larger accounts

• UseMarketOrders (bool)
  true = market orders (immediate execution, max deviation SlippageTolerance)
  false = limit orders SlippageTolerance points better than the current
          price; unfilled ones are deleted at the next check
  Reductions are always partial closes at market
  Tip: Use true for simplicity, false to avoid paying the spread at once

• SlippageTolerance (float64)
  Points: the maximum deviation of market orders (OrderSendRequest.Slippage),
  or how much better than the price a limit order is placed
  Example: 50 = 50 points (5 pips on EURUSD)
  Must be > 0 for limit orders

• MaxTradesPerCycle (int)
  Maximum trades to execute in one rebalancing cycle (0 = 2 × symbols)
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

4. EXPOSURE IS NOTIONAL VALUE
   → Volume × contract size × price, converted to the account currency
   → Leverage and margin are not part of it
   → The rebalancer never opens a symbol without a position: open the
     initial positions yourself (their direction is kept)

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

❌ PROBLEM: A symbol is never traded
✅ SOLUTION: Open its initial position yourself
   • The rebalancer only resizes existing positions (it keeps their
     direction) - a symbol without a position stays at HOLD
   • Exposure = Volume × contract size × price in the account currency

❌ PROBLEM: MaxTradesPerCycle reached too often
✅ SOLUTION: Increase MaxTradesPerCycle
   • If you have 5 symbols, set to 10+ trades per cycle
   • Rule of thumb: 2× number of symbols
   • Or reduce number of symbols

═══════════════════════════════════════════════════════════════════════════*/
