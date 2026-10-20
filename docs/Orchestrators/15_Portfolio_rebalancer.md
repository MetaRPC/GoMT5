# PortfolioRebalancer - Automated Portfolio Management

## Description

**🧩 PortfolioRebalancer** maintains balanced exposure across multiple trading symbols by automatically resizing your positions to match target allocations. Like a mutual fund manager, it aims at diversification and prevents over-concentration in any single symbol.

**Strategy Principle**: It checks once when it starts and then every `CheckInterval`, compares each symbol's exposure with its target and trades when the deviation exceeds the threshold:

- **Multi-Symbol Trading**: Manages a portfolio of several symbols simultaneously
- **Target Allocations**: Each symbol gets a fraction of `TotalExposure` (e.g., 0.25 each for 4 symbols)
- **Automatic Rebalancing**: When a symbol's exposure deviates more than `RebalanceThreshold` % from its target value, it resizes the position - at the regular check, or at once (checked every minute) when the deviation exceeds `MaxDeviation`
- **Diversification**: Spreads exposure across multiple instruments

> **How "exposure" is measured:** a symbol's exposure is the **net notional value** of its positions in the account currency: `Volume × value of 1 lot`, BUY positive, SELL negative. The value of 1 lot is `ContractSize × price` when the symbol's profit currency is the account currency (EURUSD, XAUUSD on a USD account), `ContractSize` when its base currency is the account currency (USDJPY on a USD account), and `price × TickValue / TickSize` otherwise. So `TotalExposure = 10000` means $10,000 of notional value - 2,500 per symbol at 25%, about 0.02 lots of EURUSD.
>
> **It does not open positions.** The rebalancer resizes the positions you have and keeps their direction: an increase is a trade in the position's direction, a reduction a partial close. A symbol without a position (or with long and short volume that cancel out) stays at HOLD - open its initial position yourself.

**File**: `examples/demos/orchestrators/15_portfolio_rebalancer.go`

> **IMPORTANT DISCLAIMER - EDUCATIONAL EXAMPLE ONLY**
>
> This orchestrator is a demonstration showing how GoMT5 methods function and combine into automated portfolio management. **NOT a production-ready trading strategy!** Always test on demo accounts first. Exposure is notional value; it says nothing about your broker's margin requirements.

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
    |  |  |  |
GetOpenPositions  GetSymbolInfo (+ Service      Service.PlaceOrder          GetPositionsBySymbol +
(monitoring)      symbol/account properties)   or BuyLimit / SellLimit    ClosePositionPartial / ClosePosition
                  (valuation)                  (increase)                 (decrease)
```

### Dependencies

- **MT5Sugar**: High-level convenience API (`GetOpenPositions`, `GetSymbolInfo`, `GetPriceInfo`, `GetPositionsBySymbol`, `BuyLimit`, `SellLimit`, `ClosePositionPartial`, `ClosePosition`)
- **MT5Service**: `GetAccountString` (account currency), `GetSymbolString` / `GetSymbolDouble` (base/profit currency, tick value/size), `PlaceOrder` (market orders with a maximum deviation), `GetOpenedTickets` + `CloseOrder` (unfilled limit orders)
- **BaseOrchestrator**: Foundation pattern with metrics tracking and lifecycle management
- **protobuf types**: `pb.PositionInfo`, `pb.OrderSendRequest`, `pb.OrderCloseRequest`

---

## Configuration Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Allocations` | map[string]float64 | equal weights of `symbols` | Target fraction per symbol (each in (0, 1]); must sum to 1.0 (±0.001) or `Start()` fails |
| `TotalExposure` | float64 | `10000.0` | Total notional value (account currency) the positions should add up to (must be > 0) |
| `MaxDeviation` | float64 | `0.20` | Deviation as a fraction (0.20 = 20%) that triggers a rebalance **at once**, between the regular checks (checked every minute; 0 = off; must not be below `RebalanceThreshold`) |
| `RebalanceThreshold` | float64 | `10.0` | Rebalance a symbol at the regular check when \|current − target\| / target × 100 > this value |
| `CheckInterval` | time.Duration | `1h` | How often to check balance (0 = 1 h); the first check runs at start |
| `MinPositionSize` | float64 | `0.01` | Smallest trade in lots: smaller adjustments are skipped |
| `UseMarketOrders` | bool | `true` | `true`: market orders; `false`: limit orders `SlippageTolerance` points better than the price, deleted at the next check if unfilled (reductions are always partial closes) |
| `SlippageTolerance` | float64 | `50` | Points: maximum deviation of market orders (`OrderSendRequest.Slippage`), or the price improvement of limit orders (must be > 0 for limit orders) |
| `MaxTradesPerCycle` | int | `10` | Max trades per rebalancing cycle (0 = 2 × number of symbols) |

Defaults are those of `orchestrators.DefaultPortfolioRebalancerConfig(symbols []string)`. Besides the `Orchestrator` methods, `*PortfolioRebalancer` has `GetCurrentAllocations() []*SymbolAllocation` (the result of the last check - empty before the first one), `GetRebalanceCount() int` and `GetLastRebalanceTime() time.Time`; they are safe to call while it runs.

### Configuration Example

```go
config := orchestrators.PortfolioRebalancerConfig{
    Allocations: map[string]float64{
        "EURUSD": 0.25,  // 25% EURUSD
        "GBPUSD": 0.25,  // 25% GBPUSD
        "USDJPY": 0.25,  // 25% USDJPY
        "XAUUSD": 0.25,  // 25% Gold
    },
    TotalExposure:      10000.0,          // $10k total notional value
    MaxDeviation:       0.20,             // 20% off -> rebalance at once
    RebalanceThreshold: 10.0,             // Rebalance at 10% off
    CheckInterval:      30 * time.Minute, // Check every 30 min
    MinPositionSize:    0.01,             // Smallest trade
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

Both commands (or `go run main.go portfolio`) launch the **PortfolioRebalancer** orchestrator with configuration from `main.go -> RunOrchestrator_PortfolioRebalancer()`. The runner uses `CheckInterval: 30 * time.Minute` and runs for 1 hour (3600 s): the rebalancer checks at start, after 30 and 60 minutes, and every minute it rebalances at once if a symbol is more than 20% (`MaxDeviation`) off target. Every 10 seconds the progress line shows `GetCurrentAllocations()`. Then it calls `Stop()` and prints the final metrics; the runner then disconnects the terminal - it does so however the run ends, also on an error or Ctrl+C. Each orchestrator also cancels its own context on Ctrl+C / SIGTERM, so a program of yours that waits on `GetContext()` can call `Stop()` and disconnect.

> The rebalancer only resizes positions that exist: on an account without positions on the four symbols it reports every symbol as HOLD and trades nothing. Trade volumes are rounded to the nearest volume step of the symbol (a reduction rounded down could leave the symbol further from its target: in a run with 0.10 lot EURUSD and a $2,500 target, 0.0777 lots rounded down to 0.07 left +35% - above `MaxDeviation`, so the minute check fired every minute with no trade possible; 0.08 leaves -10%).

---

## 🟢 Algorithm

### Flowchart

```
START
  |
Validate Configuration (in Start):
  - Do allocations sum to 100%? TotalExposure > 0? (else Start returns an error)
  - Defaults: CheckInterval 0 -> 1h, MaxTradesPerCycle 0 -> 2 × symbols
  - Read the account currency
  - Calculate target values per symbol: TotalExposure × allocation
  |
FIRST CHECK at once, then:
MONITORING LOOP (every CheckInterval; every minute a MaxDeviation check):
  |
  Delete our limit orders of the last cycle that did not fill
  |
  Get all open positions
  |
  Calculate the net exposure per symbol (account currency):
    For each position of a symbol in Allocations:
      value = Volume × value of 1 lot at the position's current price
      netExposure[symbol] += value (BUY) / -= value (SELL)
  |
  Analyze each symbol in Allocations:
    | No net position -> HOLD (the rebalancer does not open positions)
    | |net exposure| vs. target value
    | Deviation % = (current - target) / target × 100
    | Action needed (in the position's direction): increase / reduce / HOLD
  |
  Is rebalancing needed?
    -> Regular check: ANY symbol |deviation %| > RebalanceThreshold
       MaxDeviation check: ANY symbol |deviation %| > MaxDeviation × 100
       |
       REBALANCING (at most MaxTradesPerCycle trades):
         For each symbol that needs adjustment:
           | lots = |target - current| / value of 1 lot, rounded to the nearest volume step
           | lots < MinPositionSize -> skip
           | Increase -> market order in the position's direction
           |             (max deviation SlippageTolerance), or a limit order
           |             SlippageTolerance points better (UseMarketOrders = false)
           | Reduce   -> partial close of the positions in that direction,
           |             largest first (whole position if the rest would be
           |             below the minimum volume)
           | (a failed adjustment is counted as an error and skipped)
           |
         Update rebalanceCount
         Update lastRebalance timestamp
    |
    -> If all within threshold:
       Continue monitoring
  |
REPEAT until stopped (Stop() or Ctrl+C)
  |
ON STOP: wait for a check in progress, delete unfilled limit orders
  |
END
```

### 💡 Visual Example

**Portfolio: USD account, TotalExposure $40,000 across 4 symbols, 25% each** (target $10,000 each; at EURUSD 1.1000 one lot is $110,000, so $10,000 is about 0.09 lots)

#### Initial State (Balanced)

```
Target Allocations:
+--------+----------+--------------+
| Symbol | Target % | Target value |
+--------+----------+--------------+
| EURUSD |   25%    | $10,000      |
| GBPUSD |   25%    | $10,000      |
| USDJPY |   25%    | $10,000      |
| XAUUSD |   25%    | $10,000      |
+--------+----------+--------------+
Open (all long): positions worth about $10,000 on each symbol -> BALANCED
```

#### After Price and Volume Changes (Unbalanced)

```
You added EURUSD volume, USDJPY was partly closed
(price moves change the values too: exposure is volume × contract × price)

+--------+-------------+---------+-----------+--------+
| Symbol | Target val. | Current | Deviation | Action |
+--------+-------------+---------+-----------+--------+
| EURUSD | $10,000     | $12,100 | +21%      | SELL   |
| GBPUSD | $10,000     | $10,200 | +2%       | HOLD   |
| USDJPY | $10,000     |  $8,000 | -20%      | BUY    |
| XAUUSD | $10,000     |  $9,800 | -2%       | HOLD   |
+--------+-------------+---------+-----------+--------+
TRIGGER: |deviation| > RebalanceThreshold (10%) - and EURUSD > MaxDeviation
         (20%), so the minute check rebalances without waiting for the
         next CheckInterval
```

#### After Rebalancing

```
Rebalancing Actions:
1. EURUSD (SELL = reduce the long): $2,100 / $110,000 per lot = 0.019
   -> rounded to 0.02 lot, closed from the largest EURUSD position
2. USDJPY (BUY = increase the long): $2,000 / $100,000 per lot = 0.02 lots
   -> market BUY 0.02 USDJPY

+--------+-------------+---------+-----------+--------+
| Symbol | Target val. | Current | Deviation | Action |
+--------+-------------+---------+-----------+--------+
| EURUSD | $10,000     |  $9,900 | -1%       | HOLD   |
| GBPUSD | $10,000     | $10,200 | +2%       | HOLD   |
| USDJPY | $10,000     | $10,000 | 0%        | HOLD   |
| XAUUSD | $10,000     |  $9,800 | -2%       | HOLD   |
+--------+-------------+---------+-----------+--------+
```

The volume step limits how closely small targets can be matched: with $10,000 per symbol one EURUSD volume step (0.01 lot) is about 11% of the target.

---

## Used MT5Sugar Methods

### 1. GetOpenPositions

```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Purpose**: Retrieves all currently open positions.

**Usage in PortfolioRebalancer**:

- Called at every check
- Calculates the net exposure per symbol (`Volume × value of 1 lot` at the position's `PriceCurrent`, BUY +, SELL −)
- If it fails, the error is counted and no rebalancing happens in that check

### 2. GetSymbolInfo

```go
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error)
```

**Usage in PortfolioRebalancer**: read once per symbol - `ContractSize` for the value of 1 lot, `VolumeMin`/`VolumeStep`/`VolumeMax` to round trade volumes, `Point`/`Digits` for limit prices. With the base/profit currency (`SYMBOL_CURRENCY_BASE`, `SYMBOL_CURRENCY_PROFIT`) and `SYMBOL_TRADE_TICK_VALUE`/`SYMBOL_TRADE_TICK_SIZE` from the service, it converts the value into the account currency.

### 3. GetPositionsBySymbol

```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*pb.PositionInfo, error)
```

**Purpose**: Gets positions for a specific symbol.

**Usage in PortfolioRebalancer**:

- Used when reducing exposure
- The positions in the net direction are closed from the largest down, only as much volume as the reduction needs

### 4. Service.PlaceOrder / BuyLimit / SellLimit

```go
func (s *MT5Service) PlaceOrder(ctx context.Context, req *pb.OrderSendRequest) (*OrderResult, error)
func (s *MT5Sugar) BuyLimit(symbol string, volume, price float64) (uint64, error)
func (s *MT5Sugar) SellLimit(symbol string, volume, price float64) (uint64, error)
```

**Usage in PortfolioRebalancer**:

- Increasing exposure with `UseMarketOrders: true`: a market order (BUY for a long, SELL for a short position) through `PlaceOrder` with `Slippage = SlippageTolerance` points; anything but return code 10009 is an error
- With `UseMarketOrders: false`: `BuyLimit` at Ask − `SlippageTolerance` points or `SellLimit` at Bid + `SlippageTolerance` points; an order that has not filled by the next check (or by `Stop()`) is deleted

### 5. ClosePositionPartial / ClosePosition

```go
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

**Usage in PortfolioRebalancer**:

- Reducing exposure: closes the needed volume (rounded to the nearest volume step) from the largest positions in the net direction
- A position is closed completely only when the rest would be below the minimum volume
- Each close counts as one trade toward `MaxTradesPerCycle`

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
TotalExposure:      5000.0,           // Small portfolio ($5k notional)
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

- `TotalExposure` is notional value in the account currency, not margin
- Example: TotalExposure $40,000 with 4 × 25% = $10,000 per symbol, about 0.09 lots of EURUSD at 1.1000
- Small targets are matched only to the volume step: keep each target well above the value of one volume step
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
- Strong trending markets (rebalancing trims the winners)
- Very small accounts (<$1,000)

---

## Configuration Location

**Runtime configuration is located in**: `examples/demos/main.go -> RunOrchestrator_PortfolioRebalancer()`

This separation exists because:

1. **Code reusability**: Same rebalancer runs with different symbol allocations

2. **Quick testing**: Change allocations without modifying logic

3. **User examples**: Shows how to configure portfolio

4. **Centralized entry point**: Single place for all orchestrator launches

