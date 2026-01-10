# PortfolioRebalancer - Automated Portfolio Management

## Description

**🧩 PortfolioRebalancer** maintains balanced exposure across multiple trading symbols by automatically rebalancing positions to match target allocations. Like a mutual fund manager, it ensures diversification and prevents over-concentration in any single symbol.

**Strategy Principle**: Continuously monitors actual exposure vs. target allocations and rebalances when deviation exceeds threshold:

- **Multi-Symbol Trading**: Manages portfolio of 2-10+ symbols simultaneously
- **Target Allocations**: Each symbol gets specific percentage (e.g., 25% each for 4 symbols)
- **Automatic Rebalancing**: When symbol deviates >10% from target, rebalances
- **Diversification**: Reduces risk by spreading exposure across multiple instruments
- **Set-and-Forget**: Runs automatically without manual intervention

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
GetOpenPositions  BuyMarket  ClosePosition
(monitoring)      (increase) (decrease)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetOpenPositions`, `GetPositionsBySymbol`, `BuyMarket`, `SellMarket`, `ClosePosition`)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo` for position data structures

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Allocations` | map[string]float64 | *required* | Target allocation % per symbol (must sum to 1.0) |
| `TotalExposure` | float64 | `10000.0` | Total portfolio value to maintain |
| `MaxDeviation` | float64 | `0.20` | Maximum deviation from target (20%) |
| `RebalanceThreshold` | float64 | `10.0` | % deviation to trigger rebalance |
| `CheckInterval` | time.Duration | `1h` | How often to check balance |
| `MinPositionSize` | float64 | `0.01` | Minimum lot size for positions |
| `UseMarketOrders` | bool | `true` | Use market orders (true) or limit (false) |
| `SlippageTolerance` | float64 | `50` | Max slippage in points for limit orders |
| `MaxTradesPerCycle` | int | `10` | Max trades per rebalancing cycle |

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

Both commands launch the **PortfolioRebalancer** orchestrator with configuration from `main.go -> RunOrchestrator_PortfolioRebalancer()`.

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Validate Configuration:
  - Do allocations sum to 100%?
  - Calculate target dollar values per symbol
  |
MONITORING LOOP (every CheckInterval):
  |
  Get all open positions
  |
  Calculate current exposure per symbol:
    For each position:
      currentExposure[symbol] += position.Volume × 100000
  |
  Analyze each symbol:
    | Current value vs. target value
    | Deviation percentage
    | Action needed: BUY / SELL / HOLD
  |
  Is rebalancing needed?
    -> If ANY symbol deviation > RebalanceThreshold:
       |
       REBALANCING:
         For each symbol that needs adjustment:
           | Calculate lot size needed
           | If ActionRequired == "BUY":
           |   -> Open BUY position
           | If ActionRequired == "SELL":
           |   -> Close existing position
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

**Portfolio: $10,000 across 4 symbols, 25% each**

#### Initial State (Balanced)

```
Target Allocations:
+--------+----------+-------------+
| Symbol | Target % | Target $    |
+--------+----------+-------------+
| EURUSD |   25%    | $2,500      |
| GBPUSD |   25%    | $2,500      |
| USDJPY |   25%    | $2,500      |
| XAUUSD |   25%    | $2,500      |
+--------+----------+-------------+
Portfolio: BALANCED
```

#### After Market Moves (Unbalanced)

```
EURUSD rises +20%, USDJPY falls -15%

+--------+----------+-------------+-----------+-----------+--------+
| Symbol | Target % | Target $    | Current $ | Deviation | Action |
+--------+----------+-------------+-----------+-----------+--------+
| EURUSD |   25%    | $2,500      | $3,000    | +20%      | SELL   |
| GBPUSD |   25%    | $2,500      | $2,400    | -4%       | HOLD   |
| USDJPY |   25%    | $2,500      | $2,125    | -15%      | BUY    |
| XAUUSD |   25%    | $2,500      | $2,475    | -1%       | HOLD   |
+--------+----------+-------------+-----------+-----------+--------+
Portfolio: UNBALANCED

TRIGGER: EURUSD deviation (+20%) exceeds RebalanceThreshold (10%)
ACTION: Rebalancing needed!
```

#### After Rebalancing

```
Rebalancing Actions:
1. SELL $500 of EURUSD ($3,000 -> $2,500)
2. BUY $375 of USDJPY ($2,125 -> $2,500)

+--------+----------+-------------+-----------+-----------+--------+
| Symbol | Target % | Target $    | Current $ | Deviation | Action |
+--------+----------+-------------+-----------+-----------+--------+
| EURUSD |   25%    | $2,500      | $2,500    | 0%        | HOLD   |
| GBPUSD |   25%    | $2,500      | $2,400    | -4%       | HOLD   |
| USDJPY |   25%    | $2,500      | $2,500    | 0%        | HOLD   |
| XAUUSD |   25%    | $2,500      | $2,475    | -1%       | HOLD   |
+--------+----------+-------------+-----------+-----------+--------+
Portfolio: BALANCED
```

---

## Used MT5Sugar Methods

### 1. GetOpenPositions

```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all currently open positions.

**Usage in PortfolioRebalancer**:

- Called every CheckInterval
- Calculates current exposure per symbol
- Returns all positions for analysis

### 2. GetPositionsBySymbol

```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*pb.PositionInfo, error)
```

**Purpose**: Gets positions for a specific symbol.

**Usage in PortfolioRebalancer**:

- Used when reducing exposure (SELL action)
- Identifies positions to close
- Returns positions filtered by symbol

### 3. BuyMarket / SellMarket

```go
func (s *MT5Sugar) BuyMarket(symbol string, lots float64) (uint64, error)
func (s *MT5Sugar) SellMarket(symbol string, lots float64) (uint64, error)
```

**Purpose**: Opens market position at current price.

**Usage in PortfolioRebalancer**:

- Called when increasing exposure (BUY action)
- Opens new positions to restore target allocation
- Returns ticket number

### 4. ClosePosition

```go
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

**Purpose**: Closes specific position by ticket.

**Usage in PortfolioRebalancer**:

- Called when reducing exposure (SELL action)
- Closes positions to restore target allocation
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
TotalExposure:      5000.0,           // Small portfolio
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

- Don't use full account balance
- Recommended: 30-50% of account
- Leave buffer for drawdowns

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
- Strong trending markets (cuts winners early)
- Very small accounts (<$1,000)

---

## Configuration Location

**Runtime configuration is located in**: `examples/demos/main.go -> RunOrchestrator_PortfolioRebalancer()`

This separation exists because:

1. **Code reusability**: Same rebalancer runs with different symbol allocations

2. **Quick testing**: Change allocations without modifying logic

3. **User examples**: Shows how to configure portfolio

4. **Centralized entry point**: Single place for all orchestrator launches

