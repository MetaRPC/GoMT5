# MT5Sugar API Reference

> **Note:** Complete API reference for examples/mt5/MT5Sugar.go with enhanced navigation. Single-page format for easy browsing.

MT5Sugar is the high-level API wrapper providing ultra-simple one-liner methods for all common MT5 operations. It automatically handles contexts, timeouts, and provides smart defaults for all parameters.

## Table of Contents

### 📡 Connection & Health
- [NewMT5Sugar](#newmt5sugar)
- [QuickConnect](#quickconnect)
- [IsConnected](#isconnected)
- [Ping](#ping)
- [GetAccount](#getaccount)
- [GetService](#getservice)
- [Disconnect](#disconnect)

### 💰 Account Information
- [GetAccountInfo](#getaccountinfo)
- [GetBalance](#getbalance)
- [GetEquity](#getequity)
- [GetMargin](#getmargin)
- [GetFreeMargin](#getfreemargin)
- [GetMarginLevel](#getmarginlevel)
- [GetProfit](#getprofit)
- [GetDailyStats](#getdailystats)

### 📊 Position Queries
- [GetOpenPositions](#getopenpositions)
- [GetPositionByTicket](#getpositionbyticket)
- [GetPositionsBySymbol](#getpositionsbysymbol)
- [CountOpenPositions](#countopenpositions)
- [HasOpenPosition](#hasopenposition)
- [GetProfitBySymbol](#getprofitbysymbol)
- [GetTotalProfit](#gettotalprofit)

### 📈 Price Information
- [GetBid](#getbid)
- [GetAsk](#getask)
- [GetSpread](#getspread)
- [GetPriceInfo](#getpriceinfo)
- [WaitForPrice](#waitforprice)

### 🔢 Symbol Information
- [GetAllSymbols](#getallsymbols)
- [GetSymbolInfo](#getsymbolinfo)
- [GetSymbolDigits](#getsymboldigits)
- [GetMinStopLevel](#getminstoplevel)
- [IsSymbolAvailable](#issymbolavailable)

### 🛒 Market Orders
- [BuyMarket](#buymarket)
- [BuyMarketWithSLTP](#buymarketwithsltp)
- [BuyMarketWithPips](#buymarketwithpips)
- [SellMarket](#sellmarket)
- [SellMarketWithSLTP](#sellmarketwithsltp)
- [SellMarketWithPips](#sellmarketwithpips)

### ⏳ Pending Orders
- [BuyLimit](#buylimit)
- [BuyLimitWithSLTP](#buylimitwithsltp)
- [BuyStop](#buystop)
- [SellLimit](#selllimit)
- [SellLimitWithSLTP](#selllimitwithsltp)
- [SellStop](#sellstop)

### ✏️ Position Management
- [ClosePosition](#closeposition)
- [ClosePositionPartial](#closepositionpartial)
- [CloseAllPositions](#closeallpositions)
- [CloseAllBySymbol](#closeallbysymbol)
- [ModifyPositionSL](#modifypositionsl)
- [ModifyPositionTP](#modifypositiontp)
- [ModifyPositionSLTP](#modifypositionsltp)

### 📜 History & Statistics
- [GetDealsToday](#getdealstoday)
- [GetDealsYesterday](#getdealsyesterday)
- [GetDealsThisWeek](#getdealsthisweek)
- [GetDealsThisMonth](#getdealsthismonth)
- [GetDealsDateRange](#getdealsdaterange)
- [GetProfitToday](#getprofittoday)
- [GetProfitThisWeek](#getprofitthisweek)
- [GetProfitThisMonth](#getprofitthismonth)

### ⚖️ Risk Management
- [CalculatePositionSize](#calculatepositionsize)
- [CalculateRequiredMargin](#calculaterequiredmargin)
- [CalculateSLTP](#calculatesltp)
- [CanOpenPosition](#canopenposition)
- [GetMaxLotSize](#getmaxlotsize)

---

## 📡 Connection & Health

## NewMT5Sugar

Creates a new MT5Sugar instance with the provided credentials. This is the main entry point for using the Sugar API. It initializes both the low-level Account and mid-level Service layers automatically.

**Signature**
```go
func NewMT5Sugar(user uint64, password string, grpcServer string) (*MT5Sugar, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| user | uint64 | MT5 account login number |
| password | string | MT5 account password |
| grpcServer | string | gRPC server address (host:port, e.g., "mt5.mrpc.pro:443"; empty string = "mt5.mrpc.pro:443") |

**Returns**

Returns `*MT5Sugar` instance ready for connection, or error if initialization fails.

---

## QuickConnect

Connects to MT5 terminal using cluster name (RECOMMENDED). This is the easiest connection method - just provide your broker's cluster name. Sends `ConnectEx` with the login and password given to NewMT5Sugar and `MtClusterName: clusterName` (nothing else), with a 30-second timeout.

**Signature**
```go
func (s *MT5Sugar) QuickConnect(clusterName string) error
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| clusterName | string | MT5 cluster identifier (e.g., "FxPro-MT5 Demo", "ICMarkets-Live02") |

**Returns**

Returns error if connection fails, nil on success.

---

## IsConnected

Checks if the connection to MT5 terminal is alive. This is a quick boolean check using 3-second timeout. Returns false if connection is dead or health check times out. Does not return errors.

**Signature**
```go
func (s *MT5Sugar) IsConnected() bool
```

**Returns**

Returns true if connected and alive, false otherwise.

---

## Ping

Verifies the connection health to MT5 terminal with detailed error reporting. Unlike IsConnected(), this method returns an error explaining why connection failed. Uses 3-second timeout. Useful for debugging connection issues.

**Signature**
```go
func (s *MT5Sugar) Ping() error
```

**Returns**

Returns error with details if ping fails or connection is dead, nil if healthy.

---

## GetAccount

Returns the underlying MT5Account instance for low-level operations. Use this when you need direct access to protobuf structures or maximum control over request parameters. (To close the connection you can also call `Disconnect()` on the Sugar.)

**Signature**
```go
func (s *MT5Sugar) GetAccount() *helpers.MT5Account
```

**Returns**

Returns `*helpers.MT5Account` (package `github.com/MetaRPC/GoMT5/package/Helpers`) instance used by the underlying Service layer.

---

## GetService

Returns the underlying MT5Service instance for operations that require more control than Sugar API provides. Use this when you need access to mid-level API features like custom timeouts or advanced parameters.

**Signature**
```go
func (s *MT5Sugar) GetService() *MT5Service
```

**Returns**

Returns `*MT5Service` instance used by this Sugar wrapper.

---

## Disconnect

Gracefully disconnects from the MT5 terminal and closes the gRPC connection (via `MT5Service.Disconnect` -> `MT5Account.DisconnectAndClose`). If `deleteOnDisconnect` is true, the terminal instance is deleted instead of stopped.

**Signature**
```go
func (s *MT5Sugar) Disconnect(deleteOnDisconnect ...bool) error
```

**Returns**

Returns error if the Disconnect RPC or closing the connection fails, nil on success (also nil for a nil Sugar).

---

## 💰 Account Information

## GetAccountInfo

Retrieves complete account information in one call. This is more efficient than calling individual Get* methods. Perfect for account monitoring dashboards or trading reports. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetAccountInfo() (*AccountInfo, error)
```

**Returns**

Returns `*AccountInfo` structure with all account data, or error if query fails.

---

## GetBalance

Returns the current account balance (deposit amount). This is the initial deposit plus/minus closed position profits/losses, not affected by floating profit. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetBalance() (float64, error)
```

**Returns**

Returns current balance as float64, or error if query fails.

---

## GetEquity

Returns the current account equity (balance + floating profit). Equity = Balance + Profit from open positions. This is the real-time value of your account including unrealized gains/losses. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetEquity() (float64, error)
```

**Returns**

Returns current equity as float64, or error if query fails.

---

## GetMargin

Returns the amount of margin currently used by open positions. This is the collateral locked by the broker for your active trades. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetMargin() (float64, error)
```

**Returns**

Returns used margin as float64, or error if query fails.

---

## GetFreeMargin

Returns the amount of margin available for new positions. Free Margin = Equity - Used Margin. This is how much you can use for new trades. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetFreeMargin() (float64, error)
```

**Returns**

Returns free margin as float64, or error if query fails.

---

## GetMarginLevel

Returns the margin level percentage. Margin Level = (Equity / Used Margin) * 100. Values below 100% indicate danger of margin call. Returns 0 if no positions are open. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetMarginLevel() (float64, error)
```

**Returns**

Returns margin level percentage as float64, or error if query fails.

---

## GetProfit

Returns the total floating profit/loss from all open positions. This is the unrealized profit that's not yet added to balance. Positive values mean profit, negative mean loss. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetProfit() (float64, error)
```

**Returns**

Returns total floating P/L as float64, or error if query fails.

---

## 📊 Position Queries

## GetOpenPositions

Returns all currently open positions as protobuf PositionInfo structures. The positions are sorted by open time (oldest first). Each PositionInfo contains full details: ticket, symbol, type, volume, open price, current profit, SL/TP, etc. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error)
```

**Returns**

Returns slice of `*pb.PositionInfo` with all open positions, or error if query fails.

---

## GetPositionByTicket

Finds and returns a specific position by its ticket number. This is useful when you need detailed information about a position you opened earlier. If no open position has this ticket (e.g. it was closed) it returns an error "position with ticket N not found". Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetPositionByTicket(ticket uint64) (*pb.PositionInfo, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ticket | uint64 | Position ticket number to search for |

**Returns**

Returns `*pb.PositionInfo` for the position, or error if not found or query fails.

---

## GetPositionsBySymbol

Returns all open positions for a specific trading symbol. Filters positions by symbol name. Useful for monitoring exposure to a single currency pair or asset. Returns an empty (nil) slice if no positions found. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*pb.PositionInfo, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol to filter by (e.g., "EURUSD", "XAUUSD") |

**Returns**

Returns slice of `*pb.PositionInfo` for the symbol, or error if query fails.

---

## CountOpenPositions

Returns the total number of currently open positions. This is more efficient than len(GetOpenPositions()) as it queries the count directly from MT5 without retrieving full position details. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) CountOpenPositions() (int, error)
```

**Returns**

Returns total number of open positions (int), or error if query fails.

---

## HasOpenPosition

Checks if there are any open positions for a specific symbol. This is a quick boolean wrapper around GetPositionsBySymbol (same request) for when you only need to know if positions exist. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) HasOpenPosition(symbol string) (bool, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol to check (e.g., "EURUSD", "GBPUSD") |

**Returns**

Returns true if at least one position exists, false otherwise, or error if query fails.

---

## GetProfitBySymbol

Calculates and returns total floating profit/loss for a specific symbol. This sums up profit from all positions matching the symbol. Useful for tracking per-symbol performance. Returns 0 if no positions for symbol. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetProfitBySymbol(symbol string) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol to calculate profit for (e.g., "EURUSD") |

**Returns**

Returns total profit/loss for symbol as float64, or error if query fails.

---

## GetTotalProfit

Calculates and returns total floating profit/loss from all open positions. This sums up the profit field from all positions. Positive value means total profit, negative means total loss. Returns 0 if no positions open. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetTotalProfit() (float64, error)
```

**Returns**

Returns total profit/loss as float64, or error if query fails.

---

## 📈 Price Information

## GetBid

Returns the current BID price for the specified symbol. BID is the price at which you can SELL. This is the real-time market price. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetBid(symbol string) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

**Returns**

Returns current BID price as float64, or error if symbol not found or query fails.

---

## GetAsk

Returns the current ASK price for the specified symbol. ASK is the price at which you can BUY. This is the real-time market price. The spread is ASK - BID. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetAsk(symbol string) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

**Returns**

Returns current ASK price as float64, or error if symbol not found or query fails.

---

## GetSpread

Returns the current spread in points for the specified symbol. Spread is the difference between ASK and BID in points (not price units). For EURUSD with 5 digits, 1 point = 0.00001. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetSpread(symbol string) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

**Returns**

Returns current spread in points as float64, or error if symbol not found.

---

## GetPriceInfo

Returns complete price information for the specified symbol. This is a convenience method that retrieves BID, ASK, spread, and timestamp all in one call. More efficient than calling individual methods. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetPriceInfo(symbol string) (*PriceInfo, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

**Returns**

Returns `*PriceInfo` structure with all price data, or error if symbol not found.

---

## WaitForPrice

Waits for a price update for the specified symbol with timeout. This method polls for valid price data (BID > 0 and ASK > 0) until timeout expires. Useful for waiting for market to open or for first price tick after connection.

**Signature**
```go
func (s *MT5Sugar) WaitForPrice(symbol string, timeout time.Duration) (*PriceInfo, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol to wait for (e.g., "EURUSD") |
| timeout | time.Duration | Maximum time to wait (e.g., 5*time.Second) |

**Returns**

Returns `*PriceInfo` with valid price data, or error if timeout expires.

---

## 🔢 Symbol Information

## GetAllSymbols

Retrieves a list of all available trading symbols. This returns symbol names only (not full info). Useful for discovering available instruments or building symbol selection menus. Uses 15-second timeout (longer than single symbol queries due to potentially large number of symbols).

**Signature**
```go
func (s *MT5Sugar) GetAllSymbols() ([]string, error)
```

**Returns**

Returns slice of symbol names ([]string), or error if query fails.

---

## GetSymbolInfo

Retrieves comprehensive information about a symbol in one call. This is more efficient than calling individual methods for each property. Perfect for validation before placing orders. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD") |

**Returns**

Returns `*SymbolInfo` structure with all important symbol parameters, or error if symbol not found.

---

## GetSymbolDigits

Returns the number of decimal places for the symbol price. For example, EURUSD typically has 5 digits (1.08123), gold might have 2. This is essential for proper price formatting and calculations. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetSymbolDigits(symbol string) (int32, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |

**Returns**

Returns number of decimal places (int32), or error if symbol not found.

---

## GetMinStopLevel

Returns the minimum allowed distance for Stop Loss/Take Profit in points. This is broker-enforced minimum distance from current price to SL/TP. If 0, there's no minimum (market execution). Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetMinStopLevel(symbol string) (int64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |

**Returns**

Returns minimum stop level in points (int64), or error if symbol not found.

---

## IsSymbolAvailable

Checks if a symbol exists and its data is synchronized: returns true only if SymbolExist reports it exists and IsSymbolSynchronized reports it synchronized. It does not check the symbol's trade mode or trading permissions. Uses 3-second timeout.

**Signature**
```go
func (s *MT5Sugar) IsSymbolAvailable(symbol string) (bool, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Symbol name to check (e.g., "EURUSD") |

**Returns**

Returns true if symbol exists and is synchronized, false otherwise, or error if query fails.

---

## 🛒 Market Orders

## BuyMarket

Opens a BUY position at current market price (instant execution). This is the simplest way to open a long position. Order executes immediately at best available ASK price. Uses 10-second timeout for order execution.

**Signature**
```go
func (s *MT5Sugar) BuyMarket(symbol string, volume float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01 = micro lot, 0.1 = mini lot, 1.0 = standard lot) |

**Returns**

Returns position ticket number (uint64), or error if order rejected or fails.

---

## BuyMarketWithSLTP

Opens a BUY position with Stop Loss and Take Profit. This is the recommended way to open positions with risk management built-in. Order executes immediately at market price with SL/TP set. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) BuyMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| sl | float64 | Stop Loss price (must be BELOW entry price for BUY) |
| tp | float64 | Take Profit price (must be ABOVE entry price for BUY) |

**Returns**

Returns position ticket number (uint64), or error if order rejected.

---

## BuyMarketWithPips

Opens a BUY position with SL/TP specified as distances (not prices). It calls `CalculateSLTP(symbol, "BUY", 0, stopLossPips, takeProfitPips)` (distances multiplied by the symbol's Point, from the current Ask) and then BuyMarketWithSLTP (10-second timeout).

**Signature**
```go
func (s *MT5Sugar) BuyMarketWithPips(symbol string, volume, stopLossPips, takeProfitPips float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |
| volume | float64 | Lot size (e.g., 0.1) |
| stopLossPips | float64 | Stop Loss distance from entry in points (× Point; e.g., 50) |
| takeProfitPips | float64 | Take Profit distance from entry in points (× Point; e.g., 100) |

**Returns**

Returns position ticket number (uint64), or error if order rejected.

---

## SellMarket

Opens a SELL position at current market price (instant execution). This is the simplest way to open a short position. Order executes immediately at best available BID price. Uses 10-second timeout for order execution.

**Signature**
```go
func (s *MT5Sugar) SellMarket(symbol string, volume float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01 = micro lot, 0.1 = mini lot, 1.0 = standard lot) |

**Returns**

Returns position ticket number (uint64), or error if order rejected or fails.

---

## SellMarketWithSLTP

Opens a SELL position with Stop Loss and Take Profit. This is the recommended way to open short positions with risk management built-in. Order executes immediately at market price with SL/TP set. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) SellMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| sl | float64 | Stop Loss price (must be ABOVE entry price for SELL) |
| tp | float64 | Take Profit price (must be BELOW entry price for SELL) |

**Returns**

Returns position ticket number (uint64), or error if order rejected.

---

## SellMarketWithPips

Opens a SELL position with SL/TP specified as distances (not prices). It calls `CalculateSLTP(symbol, "SELL", 0, stopLossPips, takeProfitPips)` (distances multiplied by the symbol's Point, from the current Bid) and then SellMarketWithSLTP (10-second timeout).

**Signature**
```go
func (s *MT5Sugar) SellMarketWithPips(symbol string, volume, stopLossPips, takeProfitPips float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |
| volume | float64 | Lot size (e.g., 0.1) |
| stopLossPips | float64 | Stop Loss distance from entry in points (× Point; e.g., 50) |
| takeProfitPips | float64 | Take Profit distance from entry in points (× Point; e.g., 100) |

**Returns**

Returns position ticket number (uint64), or error if order rejected.

---

## ⏳ Pending Orders

## BuyLimit

Places a pending BUY LIMIT order (executes when price drops to specified level). Buy Limit is used to buy at a lower price than current market. Order remains pending until price reaches the specified level or order is cancelled. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) BuyLimit(symbol string, volume, price float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| price | float64 | Entry price (must be BELOW current ASK for Buy Limit) |

**Returns**

Returns pending order ticket number (uint64), or error if order rejected.

---

## BuyLimitWithSLTP

Places a BUY LIMIT order with Stop Loss and Take Profit. Combines pending order functionality with automatic risk management. Order remains pending until price reaches entry level. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) BuyLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| price | float64 | Entry price (must be BELOW current ASK) |
| sl | float64 | Stop Loss price (must be BELOW entry price) |
| tp | float64 | Take Profit price (must be ABOVE entry price) |

**Returns**

Returns pending order ticket number (uint64), or error if order rejected.

---

## BuyStop

Places a pending BUY STOP order (executes when price rises to specified level). Buy Stop is used to buy at a higher price than current market (breakout trading). Order remains pending until price reaches level or order is cancelled. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) BuyStop(symbol string, volume, price float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| price | float64 | Entry price (must be ABOVE current ASK for Buy Stop) |

**Returns**

Returns pending order ticket number (uint64), or error if order rejected.

---

## SellLimit

Places a pending SELL LIMIT order (executes when price rises to specified level). Sell Limit is used to sell at a higher price than current market. Order remains pending until price reaches the specified level or order is cancelled. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) SellLimit(symbol string, volume, price float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| price | float64 | Entry price (must be ABOVE current BID for Sell Limit) |

**Returns**

Returns pending order ticket number (uint64), or error if order rejected.

---

## SellLimitWithSLTP

Places a SELL LIMIT order with Stop Loss and Take Profit. Combines pending order functionality with automatic risk management. Order remains pending until price reaches entry level. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) SellLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| price | float64 | Entry price (must be ABOVE current BID) |
| sl | float64 | Stop Loss price (must be ABOVE entry price) |
| tp | float64 | Take Profit price (must be BELOW entry price) |

**Returns**

Returns pending order ticket number (uint64), or error if order rejected.

---

## SellStop

Places a pending SELL STOP order (executes when price drops to specified level). Sell Stop is used to sell at a lower price than current market (breakout trading). Order remains pending until price reaches level or order is cancelled. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) SellStop(symbol string, volume, price float64) (uint64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| volume | float64 | Lot size (e.g., 0.01, 0.1, 1.0) |
| price | float64 | Entry price (must be BELOW current BID for Sell Stop) |

**Returns**

Returns pending order ticket number (uint64), or error if order rejected.

---

## ✏️ Position Management

## ClosePosition

Closes a position completely by ticket number. This is the simplest way to close an open position. Closes at current market price (BID for long positions, ASK for short positions). Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) ClosePosition(ticket uint64) error
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ticket | uint64 | Position ticket number to close |

**Returns**

Returns error if close fails or position not found, nil on success.

---

## ClosePositionPartial

Closes a specified volume of a position (partial close). This allows you to take partial profit or reduce exposure while keeping position open. Not all brokers support partial closes. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ticket | uint64 | Position ticket number |
| volume | float64 | Volume to close (must be less than position volume) |

**Returns**

Returns error if close fails, volume invalid, or broker doesn't support partial close.

---

## CloseAllPositions

Closes all currently open positions across all symbols (pending orders are not touched). Iterates through all positions and attempts to close each one. Continues even if some closes fail (failures are not reported). Returns count of successfully closed positions. Uses 30-second timeout for the whole loop.

**Signature**
```go
func (s *MT5Sugar) CloseAllPositions() (int, error)
```

**Returns**

Returns number of positions successfully closed (int); error only if the position list cannot be read.

---

## CloseAllBySymbol

Closes all open positions of one symbol. Same behaviour as CloseAllPositions, filtered by `symbol`. Uses 30-second timeout for the whole loop.

**Signature**
```go
func (s *MT5Sugar) CloseAllBySymbol(symbol string) (int, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol to close all positions for (e.g., "EURUSD") |

**Returns**

Returns number of positions successfully closed (int); error only if the position list cannot be read.

---

## ModifyPositionSL

Modifies only the Stop Loss of an open position (sends `OrderModify` with Ticket and StopLoss). Use 0 to remove the Stop Loss (if the broker allows). Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) ModifyPositionSL(ticket uint64, sl float64) error
```

**Returns**

Returns error if modification rejected (return code other than 10009) or fails, nil on success.

---

## ModifyPositionTP

Modifies only the Take Profit of an open position (sends `OrderModify` with Ticket and TakeProfit). Use 0 to remove the Take Profit (if the broker allows). Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) ModifyPositionTP(ticket uint64, tp float64) error
```

**Returns**

Returns error if modification rejected (return code other than 10009) or fails, nil on success.

---

## ModifyPositionSLTP

Modifies both Stop Loss and Take Profit in one operation. This is the recommended way to update risk management levels. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) ModifyPositionSLTP(ticket uint64, sl, tp float64) error
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ticket | uint64 | Position ticket number |
| sl | float64 | New Stop Loss price (must be valid for position direction) |
| tp | float64 | New Take Profit price (must be valid for position direction) |

**Returns**

Returns error if modification rejected or fails, nil on success.

---

## 📜 History & Statistics

## GetDealsToday

Returns closed positions (deals) whose open time is today (00:00 to now, local clock; the request filters on position open time). Automatically calculates today's date range. Each deal contains full information: ticket, symbol, volume, profit, open/close times, etc. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetDealsToday() ([]*pb.PositionHistoryInfo, error)
```

**Returns**

Returns slice of `*pb.PositionHistoryInfo` with today's deals, or error if query fails.

---

## GetDealsYesterday

Returns closed positions (deals) whose open time is yesterday (full day; filter on position open time). Automatically calculates yesterday's date range (00:00 to 23:59:59). Useful for analyzing previous day's performance. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetDealsYesterday() ([]*pb.PositionHistoryInfo, error)
```

**Returns**

Returns slice of `*pb.PositionHistoryInfo` with yesterday's deals, or error if query fails.

---

## GetDealsThisWeek

Returns closed positions (deals) whose open time is in this week (filter on position open time). Week starts on Monday. Automatically calculates start of week (Monday 00:00) to current time. Useful for weekly performance tracking. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetDealsThisWeek() ([]*pb.PositionHistoryInfo, error)
```

**Returns**

Returns slice of `*pb.PositionHistoryInfo` with this week's deals, or error if query fails.

---

## GetDealsThisMonth

Returns closed positions (deals) whose open time is in this month (filter on position open time). Automatically calculates start of month (1st day 00:00) to current time. Useful for monthly performance tracking and reports. Uses 30-second timeout (longer than day/week queries due to potentially large data volume).

**Signature**
```go
func (s *MT5Sugar) GetDealsThisMonth() ([]*pb.PositionHistoryInfo, error)
```

**Returns**

Returns slice of `*pb.PositionHistoryInfo` with this month's deals, or error if query fails.

---

## GetDealsDateRange

Returns closed positions (deals) whose open time is within a custom date range. You specify exact start and end times. Useful for custom period analysis, backtesting, or generating reports for specific time frames. Uses 30-second timeout (longer to accommodate large date ranges with many deals).

**Signature**
```go
func (s *MT5Sugar) GetDealsDateRange(from, to time.Time) ([]*pb.PositionHistoryInfo, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| from | time.Time | Start date/time for the range (inclusive) |
| to | time.Time | End date/time for the range (inclusive) |

**Returns**

Returns slice of `*pb.PositionHistoryInfo` with deals in range, or error if query fails.

---

## GetProfitToday

Calculates and returns total realized profit/loss from today's closed positions. This sums the Profit field of GetDealsToday() (swap and commission are not included). Positive means net profit, negative means net loss. Returns 0 if no deals today. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetProfitToday() (float64, error)
```

**Returns**

Returns total profit/loss from today's deals as float64, or error if query fails.

---

## GetProfitThisWeek

Calculates and returns total realized profit/loss from this week's deals. This sums the Profit field of GetDealsThisWeek() (swap and commission are not included). Positive means net profit, negative means net loss. Returns 0 if no deals this week. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetProfitThisWeek() (float64, error)
```

**Returns**

Returns total profit/loss from this week's deals as float64, or error if query fails.

---

## GetProfitThisMonth

Calculates and returns total realized profit/loss from this month's deals. This sums the Profit field of GetDealsThisMonth() (swap and commission are not included). Positive means net profit, negative means net loss. Returns 0 if no deals this month. Uses 30-second timeout (via GetDealsThisMonth).

**Signature**
```go
func (s *MT5Sugar) GetProfitThisMonth() (float64, error)
```

**Returns**

Returns total profit/loss from this month's deals as float64, or error if query fails.

---

## GetDailyStats

Calculates trading statistics for today (00:00 to now). This analyzes all closed positions from today and provides performance metrics. Perfect for daily reports and performance tracking. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetDailyStats() (*DailyStats, error)
```

**Returns**

Returns `*DailyStats` structure with today's performance, or error if query fails.

---

## ⚖️ Risk Management

## CalculatePositionSize

Calculates the optimal lot size based on risk percentage with automatic margin limit protection. This is THE MOST IMPORTANT risk management tool - automatically calculates position size considering BOTH risk and margin limits.

**Algorithm:**
1. Calculate size based on risk: (Balance * RiskPercent / 100) / (StopLossPips * PipValue), where PipValue = ContractSize * Point (no tick-value or currency conversion); rounded down to VolumeStep and clamped to VolumeMin..VolumeMax
2. Calculate max size based on free margin (80% of free margin / margin for 1 lot BUY at Bid)
3. Return MINIMUM of the two; error if the result is below VolumeMin. Uses 10-second timeout.

**Signature**
```go
func (s *MT5Sugar) CalculatePositionSize(symbol string, riskPercent, stopLossPips float64) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |
| riskPercent | float64 | Percentage of balance to risk (e.g., 2.0 = 2%) |
| stopLossPips | float64 | Stop Loss distance in points (not price!) |

**Returns**

Returns recommended lot size (float64), or error if calculation fails or insufficient margin.

---

## CalculateRequiredMargin

Calculates how much margin is required to open a position. Calls OrderCalcMargin for a BUY of `volume` at the current Ask (always BUY). This helps you plan your trades and manage account exposure. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) CalculateRequiredMargin(symbol string, volume float64) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |
| volume | float64 | Desired lot size (e.g., 0.1) |

**Returns**

Returns required margin amount (float64), or error if calculation fails.

---

## CalculateSLTP

Calculates Stop Loss and Take Profit prices from entry price and pip distances. This converts pip distances to actual prices based on symbol specifications. Distances are multiplied by the symbol's Point; the results are not rounded. No own timeout: it uses GetSymbolInfo (5 s) and, when entryPrice is 0, GetAsk/GetBid (3 s).

**Signature**
```go
func (s *MT5Sugar) CalculateSLTP(symbol, direction string, entryPrice, stopLossPips, takeProfitPips float64) (float64, float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |
| direction | string | "BUY" (exact, case-sensitive); any other value is treated as SELL |
| entryPrice | float64 | Entry price (use 0 for current market price: Ask for BUY, Bid otherwise) |
| stopLossPips | float64 | Distance to SL in points (e.g., 50) |
| takeProfitPips | float64 | Distance to TP in points (e.g., 100) |

**Returns**

Returns `sl` (Stop Loss price), `tp` (Take Profit price), and error if calculation fails.

---

## CanOpenPosition

Checks if it's possible to open a position with specified volume. Checks, in order: IsSymbolAvailable, volume min/max, volume step, and required margin (OrderCalcMargin for a BUY at Ask) against free margin. Call this before placing an order to catch these rejections early. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) CanOpenPosition(symbol string, volume float64) (bool, string, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |
| volume | float64 | Desired lot size (e.g., 0.1) |

**Returns**

Returns `can` (true if position can be opened), `reason` (explanation if can't open, empty if can), and error if check failed.

---

## GetMaxLotSize

Calculates the maximum lot size you can open with current free margin. This helps prevent margin calls by showing your maximum trading capacity. Computed as 80% of free margin / margin for 1 lot (BUY at Ask), rounded down to VolumeStep and capped at VolumeMax. Uses 5-second timeout.

**Signature**
```go
func (s *MT5Sugar) GetMaxLotSize(symbol string) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| symbol | string | Trading symbol (e.g., "EURUSD") |

**Returns**

Returns maximum safe lot size (float64); 0 with nil error if the free margin is not enough for VolumeMin; error if calculation fails.

---
