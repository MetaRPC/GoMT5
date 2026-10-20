/*══════════════════════════════════════════════════════════════════════════════
 ORCHESTRATOR: PositionScaler (Pyramiding & Averaging)

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
   Advanced position sizing orchestrator that implements three sophisticated
   strategies: PYRAMIDING (add to winners), AVERAGING DOWN (add to losers),
   and SCALE OUT (gradually exit). Maximizes profits in trending markets or
   recovers from drawdowns using intelligent scaling algorithms.

 STRATEGY MODES:

   1️⃣ PYRAMIDING MODE (Recommended for trending markets):
      • Adds to WINNING positions as they move in your favor
      • Increases exposure when trend is confirmed by price action
      • Uses smaller lot sizes for each additional scale-in (risk management)
      • Locks in profits at each level while letting winners run
      • Best for strong trending markets (EURUSD trending, XAUUSD rallies)

   2️⃣ AVERAGING DOWN MODE (High risk - use with caution):
      • Adds to LOSING positions at better prices
      • Reduces average entry price to recover faster
      • VERY RISKY - requires strict maximum loss limits
      • Only suitable in strong support/resistance areas
      • NOT recommended for beginners

   3️⃣ SCALE OUT MODE (Profit protection):
      • Gradually exits positions at predefined profit levels
      • Takes partial profits while letting remainder run
      • Reduces risk exposure over time
      • Secures gains while maintaining some market exposure
      • Best for volatile markets or profit protection

 

   ⚠️ WARNING: Averaging down is DANGEROUS if price continues falling!
   Max loss limit (MaxLossToAverage) prevents catastrophic losses.

 KEY PARAMETERS:
   • Mode: Pyramiding | AveragingDown | ScaleOut
   • TriggerDistance: Points moved to trigger next scale (default: 200)
   • MinProfitToScale: Min profit before scaling (pyramiding, default: 100)
   • MaxLossToAverage: Max loss before stopping (averaging, default: 500)
   • InitialLotSize: BUY base position opened for a listed symbol without a
     position, once per run (default: 0.10 lots; 0 = manage existing only)
   • ScaleLotSize: Size for each scale-in (default: 0.05 lots)
   • MaxScales: Maximum number of scale-ins (default: 3)
   • TotalMaxLotSize: Maximum total position size (default: 1.0 lots)
   • StopLossPerScale: SL distance for scale-ins (default: 150 pts)
   • ScaleOutPercent: Part of the largest position closed per level (ScaleOut)

   Distances are in POINTS of the symbol: the point size is read from the
   symbol (EURUSD 0.00001, USDJPY 0.001, XAUUSD 0.01). Volumes are rounded
   down to the symbol's volume step.

 USE CASES:

   ✅ PYRAMIDING - Use when:
   • Strong trend confirmed (EURUSD breaking resistance, gold rally)
   • You want to maximize profit in trending conditions
   • Risk-reward ratio improves as trend continues
   • Example: NFP data confirms USD strength → pyramid EURUSD shorts

   ⚠️ AVERAGING DOWN - Use when:
   • ONLY in strong support/resistance with high confidence
   • Small position sizes with strict loss limits
   • Mean-reversion expected (range-bound markets)
   • Example: EURUSD at 1.0500 support, average down if dips to 1.0480

   ✅ SCALE OUT - Use when:
   • Volatile markets where securing partial profits is wise
   • Large position that you want to derisk gradually
   • Profit target reached but still bullish/bearish
   • Example: Gold +$50, take 33% profit, let rest run

 COMMAND-LINE USAGE:
   cd examples/demos

   go run main.go 12
   go run main.go scaler


 CONFIGURATION:
   ⚙️ All parameters configured in main.go → RunOrchestrator_PositionScaler()
   📍 See end of this file for detailed configuration examples and documentation

══════════════════════════════════════════════════════════════════════════════*/

package orchestrators

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// ══════════════════════════════════════════════════════════════════════════════
// CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════

// ScalingMode defines the position scaling strategy.
type ScalingMode int

const (
	Pyramiding   ScalingMode = iota // Add to winning positions
	AveragingDown                    // Add to losing positions
	ScaleOut                         // Gradually exit positions
)

// PositionScalerConfig holds position scaling parameters.
// Distances are in points of the symbol (the point size is read from the symbol).
type PositionScalerConfig struct {
	// Strategy
	Mode ScalingMode // Scaling mode

	// Trigger Rules
	TriggerDistance  float64 // Points moved to trigger next scale (from the base entry, then from the last scale)
	MinProfitToScale float64 // Minimum profit before scaling (pyramiding)
	MaxLossToAverage float64 // Maximum loss before stopping averaging

	// Position Sizing
	InitialLotSize  float64 // BUY base position opened for a listed symbol without one (0 = manage existing only)
	ScaleLotSize    float64 // Size for each scale-in
	MaxScales       int     // Maximum number of scale-ins (or scale-outs)
	ReducePerScale  float64 // Reduce lot size by this fraction each scale (0-1)

	// Risk Management
	TotalMaxLotSize float64 // Maximum total size of the symbol's positions in the base direction
	ScaleOutPercent float64 // Fraction of the largest position to close at each level (scale-out)
	StopLossPerScale float64 // SL distance for each scale-in position (0 = no SL)

	// Operational
	Symbols       []string      // Symbols to manage (empty = all)
	CheckInterval time.Duration // How often to check for scaling opportunities (0 = 5s)
}

// DefaultPositionScalerConfig returns sensible defaults for pyramiding.
func DefaultPositionScalerConfig(symbol string) PositionScalerConfig {
	return PositionScalerConfig{
		Mode:             Pyramiding,
		TriggerDistance:  200,
		MinProfitToScale: 100,
		MaxLossToAverage: 500,
		InitialLotSize:   0.10,
		ScaleLotSize:     0.05,
		MaxScales:        3,
		ReducePerScale:   0.0,
		TotalMaxLotSize:  1.0,
		ScaleOutPercent:  0.33,
		StopLossPerScale: 150,
		Symbols:          []string{symbol},
		CheckInterval:    5 * time.Second,
	}
}

// validate applies the default CheckInterval and rejects values that cannot work.
func (c *PositionScalerConfig) validate() error {
	if c.CheckInterval <= 0 {
		c.CheckInterval = 5 * time.Second
	}
	if c.TriggerDistance <= 0 {
		return fmt.Errorf("TriggerDistance must be > 0 points, got %v", c.TriggerDistance)
	}
	if c.MaxScales <= 0 {
		return fmt.Errorf("MaxScales must be > 0, got %d", c.MaxScales)
	}
	if c.InitialLotSize < 0 || c.MinProfitToScale < 0 || c.StopLossPerScale < 0 {
		return fmt.Errorf("InitialLotSize, MinProfitToScale and StopLossPerScale must not be negative")
	}
	if c.ReducePerScale < 0 || c.ReducePerScale >= 1 {
		return fmt.Errorf("ReducePerScale must be in [0, 1), got %v", c.ReducePerScale)
	}
	switch c.Mode {
	case Pyramiding, AveragingDown:
		if c.ScaleLotSize <= 0 || c.TotalMaxLotSize <= 0 {
			return fmt.Errorf("ScaleLotSize and TotalMaxLotSize must be > 0")
		}
		if c.Mode == AveragingDown && c.MaxLossToAverage <= 0 {
			return fmt.Errorf("MaxLossToAverage must be > 0 points in AveragingDown mode")
		}
	case ScaleOut:
		if c.ScaleOutPercent <= 0 || c.ScaleOutPercent > 1 {
			return fmt.Errorf("ScaleOutPercent must be in (0, 1], got %v", c.ScaleOutPercent)
		}
	default:
		return fmt.Errorf("unknown scaling mode %d", c.Mode)
	}
	return nil
}

// ══════════════════════════════════════════════════════════════════════════════
// POSITION SCALER IMPLEMENTATION
// ══════════════════════════════════════════════════════════════════════════════

// PositionScaler manages position scaling strategies.
type PositionScaler struct {
	*BaseOrchestrator
	sugar  *mt5.MT5Sugar
	config PositionScalerConfig

	// Position Tracking (mu guards it: GetTrackedGroups may run on another goroutine)
	mu            sync.Mutex
	trackedGroups map[string]*PositionGroup // Symbol -> position group
	symbolSpecs   map[string]symbolSpec     // Symbol -> point, digits, volume limits
	baseOpened    map[string]bool           // Symbols this run already opened a base position for
}

// PositionGroup tracks a group of related positions for one symbol.
type PositionGroup struct {
	Symbol         string
	IsBuy          bool
	BaseTicket     uint64    // Original position ticket (the oldest position of the symbol)
	ScaleTickets   []uint64  // Tickets of scale-in positions (or of the positions scaled out)
	BasePrice      float64   // Entry price of base position
	TotalLotSize   float64   // Total size of the symbol's positions in the base direction
	ScaleCount     int       // Number of scales executed
	LastScalePrice float64   // Price of last scale (0 = none yet)
	LastScaleTime  time.Time // When last scaled
}

// NewPositionScaler creates a new position scaling orchestrator.
func NewPositionScaler(sugar *mt5.MT5Sugar, config PositionScalerConfig) *PositionScaler {
	return &PositionScaler{
		BaseOrchestrator: NewBaseOrchestrator("Position Scaler"),
		sugar:            sugar,
		config:           config,
		trackedGroups:    make(map[string]*PositionGroup),
		symbolSpecs:      make(map[string]symbolSpec),
		baseOpened:       make(map[string]bool),
	}
}

// Start validates the configuration and begins position scaling operations.
func (p *PositionScaler) Start() error {
	if p.IsRunning() {
		return fmt.Errorf("position scaler already running")
	}

	// Validate configuration (applies the default CheckInterval)
	if err := p.config.validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create context, mark as started, stop on Ctrl+C
	p.begin()

	// Start monitoring loop
	p.goLoop(p.monitorLoop)

	return nil
}

// Stop gracefully stops position scaling (positions stay open).
// It waits until a check that is in progress has finished.
func (p *PositionScaler) Stop() error {
	if !p.IsRunning() {
		return fmt.Errorf("position scaler not running")
	}

	// Cancel context and wait for the monitoring loop
	p.stopLoops()

	// Clear tracked groups
	p.mu.Lock()
	p.trackedGroups = make(map[string]*PositionGroup)
	p.mu.Unlock()

	// Mark as stopped
	p.MarkStopped()

	return nil
}

// monitorLoop continuously monitors positions for scaling opportunities.
func (p *PositionScaler) monitorLoop() {
	ticker := time.NewTicker(p.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.GetContext().Done():
			return
		case <-ticker.C:
			p.checkScalingOpportunities()
		}
	}
}

// checkScalingOpportunities looks for positions to scale.
func (p *PositionScaler) checkScalingOpportunities() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Get all open positions
	positions, err := p.sugar.GetOpenPositions()
	if err != nil {
		p.IncrementError(fmt.Sprintf("failed to get positions: %v", err))
		return
	}

	// Open the base position of a listed symbol that has none (once per run)
	if p.openBasePositions(positions) {
		return // the next check picks the new position up
	}

	// Update tracked groups with current positions
	p.updateTrackedGroups(positions)

	// Check each group for scaling opportunities
	for _, group := range p.trackedGroups {
		if p.shouldScale(group) {
			if err := p.executeScale(group); err != nil {
				p.IncrementError(fmt.Sprintf("scale failed for %s: %v", group.Symbol, err))
			} else {
				p.IncrementSuccess()
			}
		}
	}

	// Update metrics
	p.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.CurrentPositions = len(p.trackedGroups)
	})
}

// openBasePositions opens a BUY base position of InitialLotSize for every listed symbol without a position.
// The example has no entry signal, so the base is always a BUY; it is opened at most once per symbol and run.
// Nothing is opened in ScaleOut mode, with InitialLotSize 0, or when Symbols is empty (= all symbols).
func (p *PositionScaler) openBasePositions(positions []*pb.PositionInfo) bool {
	if p.config.Mode == ScaleOut || p.config.InitialLotSize <= 0 {
		return false
	}
	opened := false
	for _, symbol := range p.config.Symbols {
		if p.baseOpened[symbol] || hasPositionOn(positions, symbol) {
			continue
		}
		p.baseOpened[symbol] = true // one attempt per run (a closed market must not cause a retry every check)

		spec, err := p.getSymbolSpec(symbol)
		if err != nil {
			p.IncrementError(fmt.Sprintf("base position for %s: %v", symbol, err))
			continue
		}
		lotSize := spec.volume(p.config.InitialLotSize)
		if lotSize == 0 {
			p.IncrementError(fmt.Sprintf("base position for %s: InitialLotSize %.2f is below the symbol's minimum volume %.2f",
				symbol, p.config.InitialLotSize, spec.VolumeMin))
			continue
		}
		ticket, err := p.sugar.BuyMarket(symbol, lotSize)
		if err != nil {
			p.IncrementError(fmt.Sprintf("base position for %s: %v", symbol, err))
			continue
		}
		opened = true
		p.IncrementSuccess()
		fmt.Printf("\n  [BASE] Opened BUY %s %.2f lots (ticket #%d)\n", symbol, lotSize, ticket)
		p.UpdateMetrics(func(m *OrchestratorMetrics) {
			m.TotalTrades++
			m.LastOperation = fmt.Sprintf("[BASE] Opened BUY %s %.2f lots (#%d)", symbol, lotSize, ticket)
		})
	}
	return opened
}

// hasPositionOn reports whether positions contain one on symbol.
func hasPositionOn(positions []*pb.PositionInfo, symbol string) bool {
	for _, pos := range positions {
		if pos.Symbol == symbol {
			return true
		}
	}
	return false
}

// updateTrackedGroups updates position groups based on current positions.
// The oldest position (lowest ticket) of a symbol is the base and sets the direction; the group's
// size counts the symbol's positions in that direction only.
func (p *PositionScaler) updateTrackedGroups(positions []*pb.PositionInfo) {
	// Base position per tracked symbol
	bases := make(map[string]*pb.PositionInfo)
	for _, pos := range positions {
		// Skip if symbol not in our list
		if !p.isSymbolTracked(pos.Symbol) {
			continue
		}
		if base, ok := bases[pos.Symbol]; !ok || pos.Ticket < base.Ticket {
			bases[pos.Symbol] = pos
		}
	}

	for symbol, base := range bases {
		isBuy := base.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY

		// Check if we're tracking this group (a new base starts a new group)
		group, exists := p.trackedGroups[symbol]
		if !exists || group.BaseTicket != base.Ticket {
			group = &PositionGroup{
				Symbol:       symbol,
				IsBuy:        isBuy,
				BaseTicket:   base.Ticket,
				BasePrice:    base.PriceOpen,
				ScaleTickets: make([]uint64, 0),
				ScaleCount:   0,
			}
			p.trackedGroups[symbol] = group
		}

		// Update total lot size (positions in the base direction)
		totalSize := 0.0
		for _, pos := range positions {
			if pos.Symbol == symbol && (pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY) == isBuy {
				totalSize += pos.Volume
			}
		}
		group.TotalLotSize = totalSize
	}

	// Remove groups for symbols no longer in positions
	for symbol := range p.trackedGroups {
		if _, ok := bases[symbol]; !ok {
			delete(p.trackedGroups, symbol)
		}
	}
}

// shouldScale determines if a position group should be scaled.
func (p *PositionScaler) shouldScale(group *PositionGroup) bool {
	// Check if we've hit max scales
	if group.ScaleCount >= p.config.MaxScales {
		return false
	}

	// Get symbol point value
	spec, err := p.getSymbolSpec(group.Symbol)
	if err != nil {
		p.IncrementError(fmt.Sprintf("failed to get symbol point for %s: %v", group.Symbol, err))
		return false
	}
	point := spec.Point

	// Check if total position size would exceed maximum (scale-in modes)
	if p.config.Mode != ScaleOut {
		nextScaleSize := spec.volume(p.calculateNextScaleSize(group))
		if nextScaleSize == 0 || group.TotalLotSize+nextScaleSize > p.config.TotalMaxLotSize+1e-9 {
			return false
		}
	}

	// Get current price
	priceInfo, err := p.sugar.GetPriceInfo(group.Symbol)
	if err != nil {
		return false
	}

	// Calculate profit/loss in points (from the base entry) and the move since the last scale
	var profitPoints float64
	var currentPrice float64

	if group.IsBuy {
		currentPrice = priceInfo.Bid
		profitPoints = (currentPrice - group.BasePrice) / point
	} else {
		currentPrice = priceInfo.Ask
		profitPoints = (group.BasePrice - currentPrice) / point
	}

	// Favourable move since the last scale (no scale yet: since the base entry)
	moveSinceLast := profitPoints
	if group.ScaleCount > 0 {
		if group.IsBuy {
			moveSinceLast = (currentPrice - group.LastScalePrice) / point
		} else {
			moveSinceLast = (group.LastScalePrice - currentPrice) / point
		}
	}

	// Check scaling conditions based on mode
	switch p.config.Mode {
	case Pyramiding:
		// For pyramiding: position must be in profit and the price must have moved
		// TriggerDistance in our favour since the last scale (first scale: since the base entry)
		return profitPoints >= p.config.MinProfitToScale && moveSinceLast >= p.config.TriggerDistance

	case AveragingDown:
		// For averaging: position must be in loss but not too much
		if profitPoints > 0 || -profitPoints > p.config.MaxLossToAverage {
			return false
		}
		// ... and the price must have moved TriggerDistance against us since the last average
		return -moveSinceLast >= p.config.TriggerDistance

	case ScaleOut:
		// For scale-out: each level is TriggerDistance beyond the previous one
		// (first level: TriggerDistance from the base entry)
		return moveSinceLast >= p.config.TriggerDistance
	}

	return false
}

// executeScale executes a scaling operation.
func (p *PositionScaler) executeScale(group *PositionGroup) error {
	spec, err := p.getSymbolSpec(group.Symbol)
	if err != nil {
		return err
	}
	point := spec.Point

	if p.config.Mode == ScaleOut {
		// Scale out: close partial position
		return p.executeScaleOut(group, spec)
	}

	// Scale in: open additional position
	lotSize := spec.volume(p.calculateNextScaleSize(group))
	if lotSize == 0 {
		return fmt.Errorf("scale size is below the symbol's minimum volume %.2f", spec.VolumeMin)
	}

	// Get current price for SL calculation
	priceInfo, err := p.sugar.GetPriceInfo(group.Symbol)
	if err != nil {
		return err
	}

	// SL distance: StopLossPerScale, at least the broker's stops level (0 = no SL)
	slDistance := 0.0
	if p.config.StopLossPerScale > 0 {
		slDistance = p.config.StopLossPerScale
		if float64(spec.StopLevel) > slDistance {
			slDistance = float64(spec.StopLevel)
		}
	}

	var ticket uint64
	if group.IsBuy {
		// Open additional buy position
		sl := 0.0
		if slDistance > 0 {
			sl = spec.price(priceInfo.Bid - slDistance*point)
		}
		tp := 0.0 // No TP for scale-ins

		ticket, err = p.sugar.BuyMarketWithSLTP(group.Symbol, lotSize, sl, tp)
		if err != nil {
			return err
		}

		group.LastScalePrice = priceInfo.Bid
	} else {
		// Open additional sell position
		sl := 0.0
		if slDistance > 0 {
			sl = spec.price(priceInfo.Ask + slDistance*point)
		}
		tp := 0.0

		ticket, err = p.sugar.SellMarketWithSLTP(group.Symbol, lotSize, sl, tp)
		if err != nil {
			return err
		}

		group.LastScalePrice = priceInfo.Ask
	}

	// Update group
	group.ScaleTickets = append(group.ScaleTickets, ticket)
	group.ScaleCount++
	group.LastScaleTime = time.Now()
	group.TotalLotSize += lotSize

	p.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.TotalTrades++
		m.LastOperation = fmt.Sprintf("[SCALE #%d/%d] Opened %s position #%d → +%.2f lots (Total: %.2f)",
			group.ScaleCount, p.config.MaxScales, group.Symbol, ticket, lotSize, group.TotalLotSize)
	})

	return nil
}

// executeScaleOut closes ScaleOutPercent of the largest position of the group's direction.
func (p *PositionScaler) executeScaleOut(group *PositionGroup, spec symbolSpec) error {
	positions, err := p.sugar.GetPositionsBySymbol(group.Symbol)
	if err != nil {
		return err
	}

	// Find the largest position in the group's direction
	var largestPos *pb.PositionInfo
	for _, pos := range positions {
		if (pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY) != group.IsBuy {
			continue
		}
		if largestPos == nil || pos.Volume > largestPos.Volume {
			largestPos = pos
		}
	}
	if largestPos == nil {
		return fmt.Errorf("no positions found for %s", group.Symbol)
	}

	// Close a percentage of it (rounded down to the volume step; the whole position when
	// the rest would be below the minimum volume)
	closeVolume := spec.volume(largestPos.Volume * p.config.ScaleOutPercent)
	if closeVolume == 0 || largestPos.Volume-closeVolume < spec.VolumeMin-1e-9 {
		closeVolume = largestPos.Volume
	}

	if closeVolume >= largestPos.Volume-1e-9 {
		err = p.sugar.ClosePosition(largestPos.Ticket)
	} else {
		err = p.sugar.ClosePositionPartial(largestPos.Ticket, closeVolume)
	}
	if err != nil {
		return err
	}

	// The next level is TriggerDistance beyond this price
	priceInfo, err := p.sugar.GetPriceInfo(group.Symbol)
	if err == nil {
		if group.IsBuy {
			group.LastScalePrice = priceInfo.Bid
		} else {
			group.LastScalePrice = priceInfo.Ask
		}
	} else {
		group.LastScalePrice = largestPos.PriceCurrent
	}

	group.ScaleTickets = append(group.ScaleTickets, largestPos.Ticket)
	group.ScaleCount++
	group.LastScaleTime = time.Now()
	group.TotalLotSize -= closeVolume

	p.UpdateMetrics(func(m *OrchestratorMetrics) {
		m.TotalTrades++
		m.LastOperation = fmt.Sprintf("[SCALE-OUT #%d/%d] Closed %s position #%d → -%.2f lots (Remaining: %.2f)",
			group.ScaleCount, p.config.MaxScales, group.Symbol, largestPos.Ticket, closeVolume, group.TotalLotSize)
	})

	return nil
}

// calculateNextScaleSize calculates the lot size for next scale-in (before rounding to the volume step).
func (p *PositionScaler) calculateNextScaleSize(group *PositionGroup) float64 {
	baseSize := p.config.ScaleLotSize

	// Apply reduction per scale if configured
	if p.config.ReducePerScale > 0 {
		reduction := 1.0 - (p.config.ReducePerScale * float64(group.ScaleCount))
		baseSize = baseSize * reduction
	}

	return baseSize
}

// getSymbolSpec gets (and caches) the symbol's point size, digits and volume limits
// with sugar.GetSymbolInfo: a point is 0.00001 for EURUSD, 0.001 for USDJPY, 0.01 for XAUUSD.
func (p *PositionScaler) getSymbolSpec(symbol string) (symbolSpec, error) {
	if spec, exists := p.symbolSpecs[symbol]; exists {
		return spec, nil
	}

	spec, err := loadSymbolSpec(p.sugar, symbol)
	if err != nil {
		return symbolSpec{}, err
	}
	p.symbolSpecs[symbol] = spec

	return spec, nil
}

// isSymbolTracked checks if symbol is in tracked list.
func (p *PositionScaler) isSymbolTracked(symbol string) bool {
	if len(p.config.Symbols) == 0 {
		return true
	}

	for _, s := range p.config.Symbols {
		if s == symbol {
			return true
		}
	}
	return false
}

// GetTrackedGroups returns a snapshot of the currently tracked position groups.
func (p *PositionScaler) GetTrackedGroups() map[string]*PositionGroup {
	p.mu.Lock()
	defer p.mu.Unlock()

	groups := make(map[string]*PositionGroup, len(p.trackedGroups))
	for symbol, group := range p.trackedGroups {
		g := *group
		g.ScaleTickets = append([]uint64(nil), group.ScaleTickets...)
		groups[symbol] = &g
	}
	return groups
}

/*══════════════════════════════════════════════════════════════════════════════
 DETAILED CONFIGURATION GUIDE
══════════════════════════════════════════════════════════════════════════════

 ⚙️ PARAMETER CONFIGURATION IS LOCATED IN main.go

 WHY THIS SEPARATION EXISTS:
 • 12_position_scaler.go = STRATEGY ENGINE (orchestrator logic, algorithm)
 • main.go → RunOrchestrator_PositionScaler() = RUNTIME CONFIGURATION (parameters)

 THIS SEPARATION IS NEEDED FOR:
 1️⃣ Code Reusability
    → Same orchestrator can run Pyramiding, Averaging, or ScaleOut
    → Switch strategies by changing config.Mode in main.go

 2️⃣ Quick Testing
    → Want tighter scaling? Change TriggerDistance in main.go
    → Want more conservative? Increase MinProfitToScale in main.go
    → Core algorithm remains untouched

 3️⃣ User Examples
    → main.go shows HOW to configure each scaling mode
    → All parameters visible with explanations

 4️⃣ Centralized Entry Point
    → All strategies launch through main.go
    → Single entry point: go run main.go scaler

 📍 WHERE TO CONFIGURE PARAMETERS:
 main.go → func RunOrchestrator_PositionScaler() (CONFIGURATION block)

 CONFIGURATION CODE IN main.go:

 func RunOrchestrator_PositionScaler() error {
     // ... connection code ...

     // ╔════════════════════════════════════════════════════════════╗
     // ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
     // ╚════════════════════════════════════════════════════════════╝
     orchConfig := orchestrators.PositionScalerConfig{
         Mode:             orchestrators.Pyramiding, // ← Choose mode
         TriggerDistance:  200,                      // ← Scale every 200pts
         MinProfitToScale: 100,                      // ← Min profit to start
         InitialLotSize:   0.10,                     // ← BUY base if none
         ScaleLotSize:     0.05,                     // ← Scale size
         MaxScales:        3,                        // ← Max 3 scales
         TotalMaxLotSize:  1.0,                      // ← Total max
         StopLossPerScale: 150,                      // ← SL per scale
         Symbols:          []string{cfg.TestSymbol},
         CheckInterval:    5 * time.Second,
     }

     scaler := orchestrators.NewPositionScaler(sugar, orchConfig)
     scaler.Start()
     // ... runs for 5 minutes ...
     scaler.Stop()
 }

 💡 EXAMPLE: Different Strategy Configurations

 // Option 1: Conservative Pyramiding (default in main.go)
 Mode:             Pyramiding,
 TriggerDistance:  200,     // Wide spacing = fewer scales
 MinProfitToScale: 100,     // Wait for confirmed profit
 MaxScales:        3,       // Limited scales = controlled risk

 // Option 2: Aggressive Pyramiding (modify in main.go)
 Mode:             Pyramiding,
 TriggerDistance:  100,     // ← Tighter spacing = more scales
 MinProfitToScale: 50,      // ← Scale sooner
 MaxScales:        5,       // ← More aggressive

 // Option 3: Averaging Down Recovery (modify in main.go)
 Mode:             AveragingDown,
 TriggerDistance:  200,     // Average every 200pts down
 MaxLossToAverage: 500,     // ← CRITICAL: Stop at -500pts
 MaxScales:        2,       // ← Limit risk exposure

 // Option 4: Scale Out Profit Taking (modify in main.go)
 Mode:             ScaleOut,
 TriggerDistance:  300,     // Take profits at +300/+600/+900pts
 ScaleOutPercent:  0.33,    // ← Close 33% of the largest position
 MaxScales:        3,       // Exit in 3 stages

 📝 IMPORTANT:
 • To change parameters → edit main.go, NOT this file
 • This file (12_position_scaler.go) contains only ORCHESTRATOR LOGIC
 • main.go contains CONFIGURATION for specific runs
 • Look for the section: ORCHESTRATOR RUNNERS

VISUAL EXAMPLE - PYRAMIDING MODE:

   Initial Entry: 0.10 lots @ 1.10000 (Entry signal)
   ────────────────────────────────────────────────────────────

   Time 0:  Price=1.10000  Position: 0.10 lots  Profit=0      [BASE POSITION]
   Time 1:  Price=1.10200  Position: 0.10 lots  Profit=+200pts [Waiting for trigger]
   Time 2:  Price=1.10200  Position: 0.15 lots  Profit=+200pts [✓ SCALE #1: +0.05 lots]
            → Trigger reached (+200pts from base)
            → Add 0.05 lots @ 1.10200
            → Total exposure: 0.15 lots

   Time 3:  Price=1.10350  Position: 0.15 lots  Profit=+350pts [Waiting for next trigger]
   Time 4:  Price=1.10400  Position: 0.20 lots  Profit=+400pts [✓ SCALE #2: +0.05 lots]
            → Trigger reached (+200pts from Scale #1)
            → Add 0.05 lots @ 1.10400
            → Total exposure: 0.20 lots

   Time 5:  Price=1.10550  Position: 0.20 lots  Profit=+550pts [Waiting for next trigger]
   Time 6:  Price=1.10600  Position: 0.25 lots  Profit=+600pts [✓ SCALE #3: +0.05 lots]
            → Trigger reached (+200pts from Scale #2)
            → Add 0.05 lots @ 1.10600
            → Total exposure: 0.25 lots (MAX REACHED)

   Time 7:  Price=1.11000  Position: 0.25 lots  Profit=+1000pts [HOLDING - Max scales reached]
            → No more scaling allowed (MaxScales=3 reached)
            → Total position: 0.25 lots across 4 entries
            → Average entry: ~1.10200

   RESULT:
   Without pyramiding: 0.10 lots × 1000pts = $100 profit
   With pyramiding:    0.25 lots weighted avg = ~$200 profit (100% increase!)

 VISUAL EXAMPLE - AVERAGING DOWN MODE (Risky):

   Initial Entry: 0.10 lots @ 1.10000 (Prediction: price will rise)
   ────────────────────────────────────────────────────────────

   Time 0:  Price=1.10000  Position: 0.10 lots  P/L=0         [BASE POSITION]
   Time 1:  Price=1.09800  Position: 0.10 lots  P/L=-200pts   [Waiting for trigger]
   Time 2:  Price=1.09750  Position: 0.15 lots  P/L=-250pts   [✓ AVERAGE #1: +0.05 @ 1.09750]
            → Price moved -200pts from base
            → Add 0.05 lots @ better price (1.09750)
            → New avg entry: 1.09916

   Time 3:  Price=1.09500  Position: 0.15 lots  P/L=-450pts   [Waiting for next trigger]
   Time 4:  Price=1.09450  Position: 0.20 lots  P/L=-500pts   [✓ AVERAGE #2: +0.05 @ 1.09450]
            → Price moved -200pts from Average #1
            → Add 0.05 lots @ 1.09450
            → New avg entry: 1.09775

   Time 5:  Price=1.09900  Position: 0.20 lots  P/L=+125pts   [✓ RECOVERY!]
            → Price recovered to 1.09900
            → Now in profit due to better avg entry (1.09775 vs 1.10000)
            → Without averaging: still -100pts loss
            → With averaging: +125pts profit

══════════════════════════════════════════════════════════════════════════════*/
