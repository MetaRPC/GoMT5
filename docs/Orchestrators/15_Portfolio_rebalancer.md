# PortfolioRebalancer - Automated Portfolio Management

## Description

**🧩 PortfolioRebalancer** maintains balanced exposure across multiple trading symbols by automatically rebalancing positions to match target allocations. Like a mutual fund manager, it aims at diversification and prevents over-concentration in any single symbol.

**Strategy Principle**: At every `CheckInterval` it compares each symbol's exposure with its target and trades when the deviation exceeds the threshold:

- **Multi-Symbol Trading**: Manages a portfolio of several symbols simultaneously
- **Target Allocations**: Each symbol gets a fraction of `TotalExposure` (e.g., 0.25 each for 4 symbols)
- **Automatic Rebalancing**: When a symbol's exposure deviates more than `RebalanceThreshold` % from its target value, it trades
- **Diversification**: Spreads exposure across multiple instruments

> **How "exposure" is measured (simplified):** a position's value is `Volume × 100000` - lots times a fixed 100,000, regardless of symbol, price, contract size, account currency or position direction. Price moves therefore do **not** change the computed allocation; only position volumes do. `TotalExposure` and the target values are in the same unit (e.g. `TotalExposure = 10000` → 2,500 per symbol at 25% = 0.025 lots).
>
> **It opens positions itself.** A symbol with no position has a deviation of −100% and gets a BUY, so on an account without positions the first check opens a BUY market position on every symbol in `Allocations`.

**File**: `examples/demos/orchestrators/15_portfolio_rebalancer.go`

> **IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY**
>
> This orchestrator is a demonstration showing how GoMT5 methods function and combine into automated portfolio management. **NOT a production-ready trading strategy!** Always test on demo accounts first. Exposure calculation is simplified and may not match your broker's margin requirements.

---

## What is Portfolio Rebalancing?

### Simple Explanation

Imagine you don't trade just one currency pair (like EURUSD), but trade SEVERAL pairs at the same time. This collection is called a **portfolio**.

**Example Portfolio**: Instead of only EURUSD, you trade:

- EURUSD (Euro/Dollar) - 25%
- GBPUSD (British Pound/Dollar) - 25%
- USDJPY (Dollar/Japanese Yen) - 25%
- XAUUSD (Gold) - 25%

### Why Trade Multiple Symbols?

**Diversification**: If EURUSD moves against you, other pairs might profit

**Risk Reduction**: Don't put all eggs in one basket

**More Opportunities**: Different markets move at different times

### Real-World Analogy

Think of 4 jars where you store coins. You want each jar to always have exactly 25% of your total coins. But as you add/remove coins, the balance shifts. **Rebalancing** means moving coins between jars to restore 25% in each.

---

## 🟢 Architecture

```
PORTFOLIO REBALANCER ORCHESTRATOR
    |
MT5Sugar Instance
    |
    |  |  |
GetOpenPositions  BuyMarket  GetPositionsBySymbol + ClosePosition
(monitoring)      (increase) (decrease)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetOpenPositions`, `GetPositionsBySymbol`, `BuyMarket`, `ClosePosition`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo` for position data structures

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Allocations` | map[string]float64 | equal weights of `symbols` | Target fraction per symbol (0.0-1.0); must sum to 1.0 (±0.001) or `Start()` fails |
| `TotalExposure` | float64 | `10000.0` | Total portfolio "value" in Volume × 100000 units (see above) |
| `MaxDeviation` | float64 | `0.20` | Declared, but **not read by the current code** |
| `RebalanceThreshold` | float64 | `10.0` | Rebalance a symbol when \|current − target\| / target × 100 > this value |
| `CheckInterval` | time.Duration | `1h` | How often to check balance |
| `MinPositionSize` | float64 | `0.01` | Lower bound for the BUY lot size (`max(lot, MinPositionSize)`) |
| `UseMarketOrders` | bool | `true` | Declared, but **not read** (BUYs always use `BuyMarket`) |
| `SlippageTolerance` | float64 | `50` | Declared, but **not read by the current code** |
| `MaxTradesPerCycle` | int | `10` | Max symbol adjustments per rebalancing cycle |

Defaults are those of `orchestrators.DefaultPortfolioRebalancerConfig(symbols []string)`. Besides the `Orchestrator` methods, `*PortfolioRebalancer` has `GetCurrentAllocations() []*SymbolAllocation` (runs a fresh analysis), `GetRebalanceCount() int` and `GetLastRebalanceTime() time.Time`.

### Configuration Example

```go
config := orchestrators.PortfolioRebalancerConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.25,  // 25% EURUSD
        "GBPUSD": 0.25,  // 25% GBPUSD
        "USDJPY": 0.25,  // 25% USDJPY
        "XAUUSD": 0.25,  // 25% Gold
    },
    TotalExposure:      10000.0,          // $10k total portfolio
    MaxDeviation:       0.20,             // 20% max deviation
    RebalanceThreshold: 10.0,             // Rebalance at 10% off
    CheckInterval:      30 * time.Minute, // Check every 30 min
    MinPositionSize:    0.01,             // Min lot size
    UseMarketOrders:    true,             // Use market orders
    SlippageTolerance:  50,               // Max 50 points slippage
    MaxTradesPerCycle:  10,               // Max 10 trades per rebalance
}

rebalancer := orchestrators.NewPortfolioRebalancer(sugar, config)
```

---

## How to Run

```bash
# Navigate to examples directory
cd examples/demos

# Option 1: By number
go run main.go 15

# Option 2: By keyword
go run main.go rebalancer
```

Both commands (or `go run main.go portfolio`) launch the **PortfolioRebalancer** orchestrator with configuration from `main.go -> RunOrchestrator_PortfolioRebalancer()`. The runner uses `CheckInterval: 30 * time.Minute` and runs for 1 hour (3600 s), so the first check happens after 30 minutes; every 10 seconds the progress line calls `GetCurrentAllocations()`. Then it calls `Stop()` and prints the final metrics.

> With the runner's `TotalExposure: 10000` each target is 2,500 = 0.025 lots. Lot sizes are sent unrounded (MT5Sugar does not normalize volumes), so a volume that is not a multiple of the symbol's volume step can be rejected by the broker.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Validate Configuration (in Start):
  - Do allocations sum to 100%? (else Start returns an error)
  - Calculate target values per symbol: TotalExposure × allocation
  |
MONITORING LOOP (every CheckInterval):
  |
  Get all open positions
  |
  Calculate current exposure per symbol:
    For each position:
      currentExposure[symbol] += position.Volume × 100000
  |
  Analyze each symbol in Allocations:
    | Current value vs. target value
    | Deviation % = (current - target) / target × 100
    | Action needed: BUY (below target) / SELL (above) / HOLD
  |
  Is rebalancing needed?
    -> If ANY symbol |deviation %| > RebalanceThreshold:
       |
       REBALANCING (at most MaxTradesPerCycle symbols):
         For each symbol that needs adjustment:
           | If ActionRequired == "BUY":
           |   lot = max(|target - current| / 100000, MinPositionSize)
           |   -> BuyMarket(symbol, lot)
           | If ActionRequired == "SELL":
           |   -> ClosePosition(first position of the symbol)
           |      (closes it completely, whatever the excess is)
           | (a failed adjustment is counted as an error and skipped)
           |
         Update rebalanceCount
         Update lastRebalance timestamp
    |
    -> If all within threshold:
       Continue monitoring
  |
REPEAT until stopped
  |
END
```

### 💡 Visual Example

**Portfolio: TotalExposure 40,000 across 4 symbols, 25% each** (values are Volume × 100000, so a target of 10,000 = 0.10 lots)

#### Initial State (Balanced)

```
Target Allocations:
+--------+----------+--------------+
| Symbol | Target % | Target value |
+--------+----------+--------------+
| EURUSD |   25%    | 10,000       |
| GBPUSD |   25%    | 10,000       |
| USDJPY |   25%    | 10,000       |
| XAUUSD |   25%    | 10,000       |
+--------+----------+--------------+
Open: 0.10 lots on each symbol -> BALANCED
```

#### After Volumes Change (Unbalanced)

```
You added 0.02 lots of EURUSD and partly closed USDJPY to 0.08 lots
(price moves alone never change these values)

+--------+-------------+---------+-----------+--------+
| Symbol | Target val. | Current | Deviation | Action |
+--------+-------------+---------+-----------+--------+
| EURUSD | 10,000      | 12,000  | +20%      | SELL   |
| GBPUSD | 10,000      | 10,000  | 0%        | HOLD   |
| USDJPY | 10,000      |  8,000  | -20%      | BUY    |
| XAUUSD | 10,000      | 10,000  | 0%        | HOLD   |
+--------+-------------+---------+-----------+--------+
TRIGGER: |deviation| > RebalanceThreshold (10%)
```

#### After Rebalancing

```
Rebalancing Actions:
1. EURUSD (SELL): ClosePosition(first EURUSD position) - the whole
   position is closed, e.g. the 0.10-lot one -> 0.02 lots remain
2. USDJPY (BUY):  BuyMarket("USDJPY", max(2000/100000, 0.01)) = 0.02 lots

+--------+-------------+---------+-----------+--------+
| Symbol | Target val. | Current | Deviation | Action |
+--------+-------------+---------+-----------+--------+
| EURUSD | 10,000      |  2,000  | -80%      | BUY    |
| GBPUSD | 10,000      | 10,000  | 0%        | HOLD   |
| USDJPY | 10,000      | 10,000  | 0%        | HOLD   |
| XAUUSD | 10,000      | 10,000  | 0%        | HOLD   |
+--------+-------------+---------+-----------+--------+

Next check (CheckInterval later): EURUSD BUY 0.08 lots -> BALANCED
```

Because a SELL closes a whole position instead of the excess only, the portfolio can need more than one cycle to settle.

---

## Used MT5Sugar Methods

### 1. GetOpenPositions

```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all currently open positions.

**Usage in PortfolioRebalancer**:

- Called every CheckInterval (and by every `GetCurrentAllocations()` call)
- Calculates current exposure per symbol (`Volume × 100000`)
- If it fails, the error is counted and no rebalancing happens in that cycle

### 2. GetPositionsBySymbol

```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*pb.PositionInfo, error)
```

**Purpose**: Gets positions for a specific symbol.

**Usage in PortfolioRebalancer**:

- Used when reducing exposure (SELL action)
- The first position of the returned slice (`positions[0]`) is closed
- Returns positions filtered by symbol (an error is ignored: nothing is closed)

### 3. BuyMarket

```go
func (s *MT5Sugar) BuyMarket(symbol string, volume float64) (uint64, error)
```

**Purpose**: Opens a BUY market position at current price.

**Usage in PortfolioRebalancer**:

- Called when increasing exposure (BUY action); `SellMarket` is never used
- Volume = `max(|target - current| / 100000, MinPositionSize)`, sent unrounded
- Returns the order ticket (the rebalancer prints it)

### 4. ClosePosition

```go
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

**Purpose**: Closes specific position by ticket.

**Usage in PortfolioRebalancer**:

- Called when reducing exposure (SELL action)
- Closes one whole position (the first one of the symbol), not just the excess
- Reduces over-allocated symbols

---

## Configuration Profiles

### Conservative Diversified

```go
Allocations: map[string]float64{
    "EURUSD": 0.40,  // 40% major currency
    "GBPUSD": 0.20,  // 20% secondary
    "USDJPY": 0.20,  // 20% yen
    "XAUUSD": 0.20,  // 20% gold hedge
},
TotalExposure:      5000.0,           // Small portfolio (Volume × 100000 units)
RebalanceThreshold: 5.0,              // Tight rebalancing
CheckInterval:      1 * time.Hour,    // Check hourly
```

**Best for**: Beginners, small accounts ($5k-$10k)

### Balanced Equal-Weight

```go
Allocations: map[string]float64{
    "EURUSD": 0.25,  // 25% each
    "GBPUSD": 0.25,
    "USDJPY": 0.25,
    "XAUUSD": 0.25,
},
TotalExposure:      10000.0,          // Medium portfolio
RebalanceThreshold: 10.0,             // Moderate
CheckInterval:      30 * time.Minute, // Every 30 min
```

**Best for**: Intermediate traders, medium accounts ($10k-$50k)

### Aggressive Commodity-Heavy

```go
Allocations: map[string]float64{
    "XAUUSD": 0.50,  // 50% gold (aggressive!)
    "EURUSD": 0.25,  // 25% euro
    "USDJPY": 0.25,  // 25% yen
},
TotalExposure:      20000.0,          // Larger portfolio
RebalanceThreshold: 15.0,             // Wider tolerance
CheckInterval:      1 * time.Hour,    // Hourly
```

**Best for**: Experienced traders, large accounts ($20k+)

---

## Best Practices

### Set Realistic TotalExposure

- `TotalExposure` is not money: it is compared with `Volume × 100000` of your positions
- Example: TotalExposure 40,000 with 4 × 25% = 0.10 lots per symbol
- Leave margin buffer for drawdowns

### Match CheckInterval to Trading Style

- Day trading: 10-30 minutes
- Swing trading: 1-4 hours
- Position trading: 12-24 hours

### Monitor Rebalancing Costs

- Track: Number of rebalances × average spread
- If costs > 2% of TotalExposure -> increase RebalanceThreshold

### Use Uncorrelated Symbols

- EURUSD + GBPUSD are highly correlated (both move together)
- Better: Mix currencies + commodities (EURUSD + XAUUSD)

### Combine with Risk Manager

- PortfolioRebalancer manages allocation
- RiskManager manages total drawdown
- Together = complete risk control

---

## When to Use

### Good Use Cases

- Multi-symbol trading systems (4-10 symbols)
- Diversification strategy
- Automated portfolio management
- Risk control across symbols

### Bad Use Cases

- Single-symbol trading (no portfolio needed)
- High-frequency trading (overhead too high)
- Strong trending markets (closes whole positions of over-weight symbols)
- Very small accounts (<$1,000)

---

## Configuration Location

**Runtime configuration is located in**: `examples/demos/main.go -> RunOrchestrator_PortfolioRebalancer()`

This separation exists because:

1. **Code reusability**: Same rebalancer runs with different symbol allocations

2. **Quick testing**: Change allocations without modifying logic

3. **User examples**: Shows how to configure portfolio

4. **Centralized entry point**: Single place for all orchestrator launches

