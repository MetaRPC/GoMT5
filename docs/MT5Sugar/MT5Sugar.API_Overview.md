# MT5Sugar - Complete API Reference

> **High-level convenience methods** for MetaTrader 5 trading automation in Go

**MT5Sugar** extends `MT5Service` with convenient methods for:

- 🎯 Risk-based position sizing (percentage risk, pip-based SL/TP)
- 📊 Automatic volume and price normalization
- 🔄 Bulk operations (close all positions, manage multiple trades)
- 📈 Historical analysis and statistics
- ✅ Pre-trade validation and margin checking
- 📉 Real-time P&L monitoring and position tracking
- 🔌 Simplified connection management

---

## Navigation by Category

### [01] 🔌 CONNECTION MANAGEMENT

**Connect to MT5 terminal and check connection status**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `QuickConnect()` | Connect using cluster name (RECOMMENDED) | [→ Docs](1.%20Connection/QuickConnect.md) |
| `IsConnected()` | Check if connected to MT5 | [→ Docs](1.%20Connection/IsConnected.md) |
| `Ping()` | Verify connection is alive | [→ Docs](1.%20Connection/Ping.md) |

**Helper methods:**

- `GetService()` - access underlying MT5Service
- `GetAccount()` - access underlying MT5Account

---

### [02] 📊 SYMBOL OPERATIONS

**Working with symbols, prices, and market data**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `GetBid()` | Current bid price for symbol | [→ Docs](3.%20Prices_Quotes/GetBid.md) |
| `GetAsk()` | Current ask price for symbol | [→ Docs](3.%20Prices_Quotes/GetAsk.md) |
| `GetSpread()` | Current spread in points | [→ Docs](3.%20Prices_Quotes/GetSpread.md) |
| `GetPriceInfo()` | Complete price snapshot (bid, ask, spread, time) | [→ Docs](3.%20Prices_Quotes/GetPriceInfo.md) |
| `WaitForPrice()` | Wait for valid price with timeout | [→ Docs](3.%20Prices_Quotes/WaitForPrice.md) |

**Data structures:**

- `PriceInfo` - complete price data structure

---

### [03] 💼 POSITION MANAGEMENT

**Manage existing positions (modify, close, track)**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `ClosePosition()` | Close position completely by ticket | [→ Docs](6.%20Position_Management/ClosePosition.md) |
| `ClosePositionPartial()` | Partial close with volume | [→ Docs](6.%20Position_Management/ClosePositionPartial.md) |
| `CloseAllPositions()` | Close ALL open positions | [→ Docs](6.%20Position_Management/CloseAllPositions.md) |
| `ModifyPositionSLTP()` | Modify both SL and TP | [→ Docs](6.%20Position_Management/ModifyPositionSLTP.md) |

---

### [04] 💹 TRADING OPERATIONS

**Place market and pending orders**

#### Market Orders (Instant Execution)

| Method | Description | Documentation |
|--------|-------------|---------------|
| `BuyMarket()` | Simple BUY at market price | [→ Docs](4.%20Simple_Trading/BuyMarket.md) |
| `SellMarket()` | Simple SELL at market price | [→ Docs](4.%20Simple_Trading/SellMarket.md) |
| `BuyMarketWithSLTP()` | BUY with SL/TP prices | [→ Docs](5.%20Trading_SLTP/BuyMarketWithSLTP.md) |
| `SellMarketWithSLTP()` | SELL with SL/TP prices | [→ Docs](5.%20Trading_SLTP/SellMarketWithSLTP.md) |

#### Pending Orders (Limit/Stop)

| Method | Description | Documentation |
|--------|-------------|---------------|
| `BuyLimit()` | Buy Limit (buy below current price) | [→ Docs](4.%20Simple_Trading/BuyLimit.md) |
| `SellLimit()` | Sell Limit (sell above current price) | [→ Docs](4.%20Simple_Trading/SellLimit.md) |
| `BuyStop()` | Buy Stop (buy above current price) | [→ Docs](4.%20Simple_Trading/BuyStop.md) |
| `SellStop()` | Sell Stop (sell below current price) | [→ Docs](4.%20Simple_Trading/SellStop.md) |
| `BuyLimitWithSLTP()` | Buy Limit with SL/TP | [→ Docs](5.%20Trading_SLTP/BuyLimitWithSLTP.md) |
| `SellLimitWithSLTP()` | Sell Limit with SL/TP | [→ Docs](5.%20Trading_SLTP/SellLimitWithSLTP.md) |

---

### [05] 🏦 ACCOUNT INFORMATION

**Quick access to account balance and margin metrics**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `GetBalance()` | Account balance (realized profit only) | [→ Docs](2.%20Balance_Margin/GetBalance.md) |
| `GetEquity()` | Current equity (balance + floating P/L) | [→ Docs](2.%20Balance_Margin/GetEquity.md) |
| `GetMargin()` | Used margin | [→ Docs](2.%20Balance_Margin/GetMargin.md) |
| `GetFreeMargin()` | Available margin for new positions | [→ Docs](2.%20Balance_Margin/GetFreeMargin.md) |
| `GetMarginLevel()` | Margin level % (Equity/Margin × 100) | [→ Docs](2.%20Balance_Margin/GetMarginLevel.md) |
| `GetProfit()` | Total floating profit/loss | [→ Docs](2.%20Balance_Margin/GetProfit.md) |

---

### [06] 📍 POSITION INFORMATION

**Query and analyze open positions**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `GetOpenPositions()` | Get ALL open positions | [→ Docs](6.%20Position_Management/GetOpenPositions.md) |
| `GetPositionByTicket()` | Get specific position by ticket | [→ Docs](6.%20Position_Management/GetPositionByTicket.md) |
| `GetPositionsBySymbol()` | Get positions for specific symbol | [→ Docs](6.%20Position_Management/GetPositionsBySymbol.md) |
| `HasOpenPosition()` | Check if symbol has open positions | [→ Docs](7.%20Position_Information/HasOpenPosition.md) |
| `CountOpenPositions()` | Total number of open positions | [→ Docs](7.%20Position_Information/CountOpenPositions.md) |

---

### [07] 📜 HISTORY & STATISTICS

**Retrieve trading history and performance analytics**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `GetDealsToday()` | All closed deals from today (00:00 to now) | [→ Docs](8.%20History_Statistics/GetDealsToday.md) |
| `GetDealsYesterday()` | Yesterday's closed deals | [→ Docs](8.%20History_Statistics/GetDealsYesterday.md) |
| `GetDealsThisWeek()` | This week's deals (Monday to now) | [→ Docs](8.%20History_Statistics/GetDealsThisWeek.md) |
| `GetDealsThisMonth()` | This month's deals (1st to now) | [→ Docs](8.%20History_Statistics/GetDealsThisMonth.md) |
| `GetDealsDateRange()` | Custom date range deals | [→ Docs](8.%20History_Statistics/GetDealsDateRange.md) |
| `GetProfitToday()` | Total realized profit today | [→ Docs](8.%20History_Statistics/GetProfitToday.md) |
| `GetProfitThisWeek()` | Total realized profit this week | [→ Docs](8.%20History_Statistics/GetProfitThisWeek.md) |
| `GetProfitThisMonth()` | Total realized profit this month | [→ Docs](8.%20History_Statistics/GetProfitThisMonth.md) |

---

### [08] 🔍 SYMBOL INFORMATION

**Symbol properties and trading conditions**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `GetSymbolInfo()` | Complete symbol information (all properties) | [→ Docs](9.%20Symbol_Information/GetSymbolInfo.md) |
| `GetAllSymbols()` | List of all available symbols | [→ Docs](9.%20Symbol_Information/GetAllSymbols.md) |
| `IsSymbolAvailable()` | Check if symbol exists and is tradable | [→ Docs](9.%20Symbol_Information/IsSymbolAvailable.md) |
| `GetMinStopLevel()` | Minimum distance for SL/TP in points | [→ Docs](9.%20Symbol_Information/GetMinStopLevel.md) |
| `GetSymbolDigits()` | Decimal precision for symbol | [→ Docs](9.%20Symbol_Information/GetSymbolDigits.md) |

**Data structures:**

- `SymbolInfo` - complete symbol data structure

---

### [09] ⚖️ RISK MANAGEMENT

**Position sizing, margin calculations, and pre-trade validation**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `CalculatePositionSize()` | **PRIMARY METHOD** - auto-calculate lot size by risk % | [→ Docs](10.%20Risk_Management/CalculatePositionSize.md) |
| `GetMaxLotSize()` | Maximum tradeable volume based on free margin | [→ Docs](10.%20Risk_Management/GetMaxLotSize.md) |
| `CanOpenPosition()` | Comprehensive validation before trading | [→ Docs](10.%20Risk_Management/CanOpenPosition.md) |
| `CalculateRequiredMargin()` | Margin needed for specific position size | [→ Docs](10.%20Risk_Management/CalculateRequiredMargin.md) |

**Example:**
```go
// Auto-calculate lot size to risk 2% with 50 pip SL
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
```

---

### [10] 🎯 TRADING HELPERS

**Pip-based trading - the intuitive way to trade**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `CalculateSLTP()` | Convert pip distances to SL/TP prices | [→ Docs](11.%20Trading_Helpers/CalculateSLTP.md) |
| `BuyMarketWithPips()` | **MOST USED** - BUY with SL/TP in pips ⭐ | [→ Docs](11.%20Trading_Helpers/BuyMarketWithPips.md) |
| `SellMarketWithPips()` | **MOST USED** - SELL with SL/TP in pips ⭐ | [→ Docs](11.%20Trading_Helpers/SellMarketWithPips.md) |

**Example:**
```go
// Buy EURUSD: 50 pip SL, 100 pip TP (1:2 R:R)
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
```

---

### [11] 💰 ACCOUNT INFORMATION (Advanced)

**Complete account snapshot and daily statistics**

| Method | Description | Documentation |
|--------|-------------|---------------|
| `GetAccountInfo()` | Complete account snapshot (one call) | [→ Docs](12.%20Account_Information/GetAccountInfo.md) |
| `GetDailyStats()` | Today's trading statistics (win rate, profit, etc.) | [→ Docs](8.%20History_Statistics/GetDailyStats.md) |

**Data structures:**

- `AccountInfo` - complete account data (login, balance, equity, margin, etc.) [→ Docs](12.%20Account_Information/AccountInfo.md)
- `DailyStats` - daily performance metrics (win rate, total deals, best/worst trades) [→ Docs](12.%20Account_Information/DailyStats.md)

---

## 🎯 Common Use Cases

### 1. Simple market order
```go
// Buy 0.10 lots EURUSD at market price
ticket, err := sugar.BuyMarket("EURUSD", 0.1)
if err != nil {
    log.Fatalf("Trade failed: %v", err)
}
fmt.Printf("Position opened: #%d\n", ticket)
```

### 2. Market order with pips (intuitive!)
```go
// Buy EURUSD: 50 pip SL, 100 pip TP
ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
fmt.Printf("✅ BUY #%d opened (1:2 R:R)\n", ticket)
```

### 3. Risk-based position sizing
```go
symbol := "EURUSD"
riskPercent := 2.0      // Risk 2% of balance
stopLossPips := 50.0    // 50 pip stop loss

// Calculate lot size automatically
lotSize, _ := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)

// Validate before trading
canOpen, reason, _ := sugar.CanOpenPosition(symbol, lotSize)
if !canOpen {
    fmt.Printf("Cannot trade: %s\n", reason)
    return
}

// Trade with calculated size
ticket, _ := sugar.BuyMarketWithPips(symbol, lotSize, 50, 100)
fmt.Printf("Position #%d opened (risking %.1f%%)\n", ticket, riskPercent)
```

### 4. Complete trading workflow
```go
// Step 1: Check connection
if !sugar.IsConnected() {
    sugar.QuickConnect("FxPro-MT5 Demo")
}

// Step 2: Check symbol availability
available, _ := sugar.IsSymbolAvailable("EURUSD")
if !available {
    log.Fatal("Symbol not available")
}

// Step 3: Calculate position size (risk 2%)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

// Step 4: Validate
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    log.Fatalf("Cannot trade: %s", reason)
}

// Step 5: Trade
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
fmt.Printf("✅ Position #%d opened\n", ticket)
```

### 5. Monitor positions and close on drawdown
```go
// Check total floating P/L
totalProfit, _ := sugar.GetTotalProfit()
fmt.Printf("Floating P/L: $%.2f\n", totalProfit)

// Emergency close on drawdown
if totalProfit < -500 {
    fmt.Println("🚨 Drawdown $500 - closing all positions!")
    closedCount, _ := sugar.CloseAllPositions()
    fmt.Printf("Closed %d positions\n", closedCount)
}
```

### 6. Daily performance report
```go
// Get account snapshot
info, _ := sugar.GetAccountInfo()
fmt.Printf("Account: #%d, Balance: $%.2f, Equity: $%.2f\n",
    info.Login, info.Balance, info.Equity)

// Get today's stats
stats, _ := sugar.GetDailyStats()
fmt.Printf("\nToday's Performance:\n")
fmt.Printf("  Trades: %d\n", stats.TotalDeals)
fmt.Printf("  Win rate: %.1f%%\n", stats.WinRate)
fmt.Printf("  Total profit: $%.2f\n", stats.TotalProfit)
fmt.Printf("  Best trade: $%.2f\n", stats.BestDeal)
fmt.Printf("  Worst trade: $%.2f\n", stats.WorstDeal)
```

### 7. Symbol-specific analysis
```go
symbol := "EURUSD"

// Check if we have open positions
hasPosition, _ := sugar.HasOpenPosition(symbol)
fmt.Printf("%s has open position: %v\n", symbol, hasPosition)

// Get P/L for this symbol
profit, _ := sugar.GetProfitBySymbol(symbol)
fmt.Printf("%s floating P/L: $%.2f\n", symbol, profit)

// Close all positions if losing
if profit < -100 {
    closed, _ := sugar.CloseAllPositions()
    fmt.Printf("Closed %d positions\n", closed)
}
```

### 8. Pending order with validation
```go
symbol := "EURUSD"
volume := 0.1
limitPrice := 1.08500

// Get current price to ensure limit is valid
priceInfo, _ := sugar.GetPriceInfo(symbol)
fmt.Printf("Current bid: %.5f, ask: %.5f\n", priceInfo.Bid, priceInfo.Ask)

// Validate margin
requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, volume)
freeMargin, _ := sugar.GetFreeMargin()

if requiredMargin > freeMargin {
    log.Fatal("Insufficient margin")
}

// Place Buy Limit below current price
if limitPrice < priceInfo.Ask {
    ticket, _ := sugar.BuyLimitWithSLTP(symbol, volume, limitPrice, 1.08000, 1.09000)
    fmt.Printf("Buy Limit placed: #%d\n", ticket)
}
```

### 9. Weekly performance tracking
```go
// Get this week's deals
dealsThisWeek, _ := sugar.GetDealsThisWeek()
fmt.Printf("Trades this week: %d\n", len(dealsThisWeek))

// Calculate metrics
totalProfit := 0.0
winners := 0

for _, deal := range dealsThisWeek {
    totalProfit += deal.Profit
    if deal.Profit > 0 {
        winners++
    }
}

winRate := 0.0
if len(dealsThisWeek) > 0 {
    winRate = (float64(winners) / float64(len(dealsThisWeek))) * 100
}

fmt.Printf("Win rate: %.1f%%\n", winRate)
fmt.Printf("Total profit: $%.2f\n", totalProfit)
```

### 10. Multi-symbol portfolio management
```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

fmt.Println("PORTFOLIO STATUS")
fmt.Println("================")

totalFloating := 0.0

for _, symbol := range symbols {
    count := 0
    profit := 0.0

    positions, _ := sugar.GetPositionsBySymbol(symbol)
    count = len(positions)

    if count > 0 {
        profit, _ = sugar.GetProfitBySymbol(symbol)
        totalFloating += profit
    }

    fmt.Printf("%s: %d positions, P/L: $%.2f\n", symbol, count, profit)
}

fmt.Printf("\nTotal floating P/L: $%.2f\n", totalFloating)
```

---

## 🔗 Related Documentation

- **[MT5Account Documentation](../API_Reference/MT5Account.md)** - Low-level gRPC/Proto methods
- **[MT5Service Documentation](../API_Reference/MT5Service.md)** - Mid-level service layer
- **MT5Sugar** (this document) - High-level convenience methods

---

## 📦 Architecture

```
MT5Sugar (High-level) ← You are here
    ↓ Simplified trading, risk management, analytics
    ↓
MT5Service (Mid-level)
    ↓ Type conversion, timeout management
    ↓
MT5Account (Low-level)
    ↓ Direct gRPC/Protobuf calls
    ↓
MetaTrader 5 Terminal
```

**Layer Comparison:**

| Feature | MT5Account | MT5Service | MT5Sugar |
|---------|------------|------------|----------|
| Complexity | Low-level Proto | Mid-level typed | High-level convenience |
| Learning curve | Steep | Moderate | Gentle |
| Verbosity | High | Medium | Low |
| Risk management | Manual | Manual | Built-in |
| Position sizing | Manual | Manual | Automatic |
| SL/TP | Prices only | Prices only | Pips or prices |
| Use case | Custom wrappers | Standard apps | Trading bots |

---

## 📝 Conventions

- All methods return `(result, error)` - always check errors
- Prices are always **absolute** (e.g., 1.08500), not relative
- Volumes are always in **lots** (e.g., 0.1), not currency units
- **Points** - minimum price increment (for 5-digit: 0.00001)
- **Pips** - standard trader unit (for 5-digit: 1 pip = 10 points = 0.0001)
- SL/TP in "pips" methods use **points** (for compatibility with broker systems)
- Timeouts are built-in (typically 5 seconds) - no need to specify
- All times are MT5 server time (not local time)

---

## 🎓 Best Practices

### Risk Management

```go
// ALWAYS use CalculatePositionSize for risk-based trading
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)  // ✅ Good

// Don't use fixed lot sizes (ignores account size)
lotSize := 0.1  // ❌ Bad - doesn't scale with account
```

### Pre-Trade Validation
```go
// ALWAYS validate before trading
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    log.Printf("Cannot trade: %s", reason)
    return
}

// Then trade
ticket, _ := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

### Error Handling
```go
// ALWAYS check errors
ticket, err := sugar.BuyMarket("EURUSD", 0.1)
if err != nil {
    log.Printf("Trade failed: %v", err)
    return
}
fmt.Printf("Success: #%d\n", ticket)
```

### Use Pip-Based Methods
```go
// Prefer pip-based methods (more intuitive)
sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)  // ✅ Easy to understand

// Instead of price-based (error-prone)
sugar.BuyMarketWithSLTP("EURUSD", 0.1, 1.08450, 1.08950)  // ❌ Manual calculation
```

### Connection Management
```go
// Check connection before trading
if !sugar.IsConnected() {
    if err := sugar.QuickConnect("FxPro-MT5 Demo"); err != nil {
        log.Fatal(err)
    }
}

// Verify with ping
if err := sugar.Ping(); err != nil {
    log.Fatal("Connection lost")
}
```

---

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/gomt5/mt5"
)

func main() {
    // 1. Create MT5Sugar instance
    sugar := mt5.NewMT5Sugar()

    // 2. Connect
    if err := sugar.QuickConnect("FxPro-MT5 Demo"); err != nil {
        log.Fatal(err)
    }

    // 3. Calculate position size (risk 2%)
    lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)

    // 4. Validate
    canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
    if !canOpen {
        log.Fatalf("Cannot trade: %s", reason)
    }

    // 5. Trade with pips
    ticket, err := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("✅ Position #%d opened!\n", ticket)
    fmt.Printf("   Risk: 2%% | SL: 50 pips | TP: 100 pips | R:R 1:2\n")
}
```

---

🎉 **You're ready to build professional trading bots with MT5Sugar!**

> **Pro tip:** Start with `BuyMarketWithPips()`, `CalculatePositionSize()`, and `CanOpenPosition()` - these three methods cover 80% of trading needs!
