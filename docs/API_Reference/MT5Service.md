# MT5Service API Reference

> **Note:** This documentation is auto-generated from [`examples/mt5/MT5Service.go`](../../examples/mt5/MT5Service.go) with enhanced navigation for easier browsing. For those who prefer viewing the complete API in a single-page reference format.

MT5Service is the mid-level API wrapper that bridges MT5Account (low-level gRPC) and MT5Sugar (high-level convenience). It unwraps protobuf structures and provides Go native types while maintaining full control over request parameters and contexts.

## Table of Contents

### 🏗️ Constructor & Utilities
- [NewMT5Service](#newmt5service)
- [GetAccount](#getaccount)

### 💰 Account Information
- [GetAccountSummary](#getaccountsummary)
- [GetAccountDouble](#getaccountdouble)
- [GetAccountInteger](#getaccountinteger)
- [GetAccountString](#getaccountstring)

### 🔢 Symbol Information
- [GetSymbolsTotal](#getsymbolstotal)
- [GetSymbolName](#getsymbolname)
- [SymbolSelect](#symbolselect)
- [SymbolExist](#symbolexist)
- [IsSymbolSynchronized](#issymbolsynchronized)
- [GetSymbolTick](#getsymboltick)
- [GetSymbolDouble](#getsymboldouble)
- [GetSymbolInteger](#getsymbolinteger)
- [GetSymbolString](#getsymbolstring)
- [GetSymbolMarginRate](#getsymbolmarginrate)
- [GetSymbolSessionQuote](#getsymbolsessionquote)
- [GetSymbolSessionTrade](#getsymbolsessiontrade)
- [GetSymbolParamsMany](#getsymbolparamsmany)

### 📈 Market Depth
- [SubscribeMarketDepth](#subscribemarketdepth)
- [GetMarketDepth](#getmarketdepth)
- [UnsubscribeMarketDepth](#unsubscribemarketdepth)

### 📊 Positions & Orders
- [GetPositionsTotal](#getpositionstotal)
- [GetOpenedOrders](#getopenedorders)
- [GetOpenedTickets](#getopenedtickets)

### ⚖️ Order Calculations & Validation
- [CalculateMargin](#calculatemargin)
- [CalculateProfit](#calculateprofit)
- [CheckOrder](#checkorder)

### 📜 History
- [GetOrderHistory](#getorderhistory)
- [GetPositionsHistory](#getpositionshistory)

### 🛒 Trading Operations
- [PlaceOrder](#placeorder)
- [ModifyOrder](#modifyorder)
- [CloseOrder](#closeorder)

### 📡 Streaming Methods
- [StreamTicks](#streamticks)
- [StreamPositionProfits](#streampositionprofits)
- [StreamOpenedTickets](#streamopenedtickets)
- [StreamTradeUpdates](#streamtradeupdates)
- [StreamTransactions](#streamtransactions)

---

## 🏗️ Constructor & Utilities

## NewMT5Service

Creates a new MT5Service wrapping an MT5Account instance. This is the main entry point for using the mid-level API. It wraps an existing low-level MT5Account and provides convenient Go native types instead of raw protobuf structures.

**Signature**
```go
func NewMT5Service(account *MT5Account) *MT5Service
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| account | *MT5Account | MT5Account instance (low-level gRPC client) |

**Returns**

Returns new `*MT5Service` instance wrapping the provided account.

---

## GetAccount

Returns the underlying MT5Account for direct low-level access. Use this when you need to work with raw protobuf structures or access methods not wrapped by MT5Service.

**Signature**
```go
func (s *MT5Service) GetAccount() *MT5Account
```

**Returns**

Returns `*MT5Account` instance used by this service.

---

## 💰 Account Information

## GetAccountSummary

Retrieves all account information in one call. This is the RECOMMENDED method for getting account data - it returns everything you need in a clean Go struct with native types instead of protobuf.

**Signature**
```go
func (s *MT5Service) GetAccountSummary(ctx context.Context) (*AccountSummary, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |

**Returns**

Returns `*AccountSummary` struct with all account properties (Login, Balance, Equity, UserName, Leverage, TradeMode, CompanyName, Currency, ServerTime, UtcTimezoneShift, Credit) using Go native types with automatic time conversion, or error if request fails.

---

## GetAccountDouble

Retrieves a double-type account property by ID. This method returns the value directly as float64 instead of wrapped protobuf structure. More convenient than MT5Account.AccountInfoDouble for single property queries.

**Signature**
```go
func (s *MT5Service) GetAccountDouble(ctx context.Context, propertyID pb.AccountInfoDoublePropertyType) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| propertyID | pb.AccountInfoDoublePropertyType | Property ID enum (ACCOUNT_BALANCE, ACCOUNT_EQUITY, ACCOUNT_MARGIN, ACCOUNT_MARGIN_FREE, ACCOUNT_PROFIT, etc.) |

**Returns**

Returns property value as float64, or error if request fails.

---

## GetAccountInteger

Retrieves an integer-type account property by ID. This method returns the value directly as int64 instead of wrapped protobuf structure. More convenient than MT5Account.AccountInfoInteger for single property queries.

**Signature**
```go
func (s *MT5Service) GetAccountInteger(ctx context.Context, propertyID pb.AccountInfoIntegerPropertyType) (int64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| propertyID | pb.AccountInfoIntegerPropertyType | Property ID enum (ACCOUNT_LOGIN, ACCOUNT_LEVERAGE, ACCOUNT_LIMIT_ORDERS, ACCOUNT_TRADE_MODE, ACCOUNT_MARGIN_SO_MODE, etc.) |

**Returns**

Returns property value as int64, or error if request fails.

---

## GetAccountString

Retrieves a string-type account property by ID. This method returns the value directly as string instead of wrapped protobuf structure. More convenient than MT5Account.AccountInfoString for single property queries.

**Signature**
```go
func (s *MT5Service) GetAccountString(ctx context.Context, propertyID pb.AccountInfoStringPropertyType) (string, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| propertyID | pb.AccountInfoStringPropertyType | Property ID enum (ACCOUNT_CURRENCY, ACCOUNT_COMPANY, ACCOUNT_NAME, ACCOUNT_SERVER) |

**Returns**

Returns property value as string, or error if request fails.

---

## 🔢 Symbol Information

## GetSymbolsTotal

Returns the count of available symbols. This method returns the count directly as int32 instead of wrapped protobuf. More convenient than MT5Account.SymbolsTotal for quick symbol counting.

**Signature**
```go
func (s *MT5Service) GetSymbolsTotal(ctx context.Context, selectedOnly bool) (int32, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| selectedOnly | bool | If true, count only symbols in Market Watch; if false, count all symbols |

**Returns**

Returns number of symbols as int32, or error if request fails.

---

## GetSymbolName

Retrieves symbol name by index position. Use this to iterate through available symbols. Set selectedOnly=true to get symbols from Market Watch only, false for all symbols.

**Signature**
```go
func (s *MT5Service) GetSymbolName(ctx context.Context, index int32, selectedOnly bool) (string, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| index | int32 | Zero-based index of symbol in list |
| selectedOnly | bool | If true, search only Market Watch; if false, search all symbols |

**Returns**

Returns symbol name as string, or error if index out of range or request fails.

---

## SymbolSelect

Adds or removes a symbol from Market Watch window. Returns true if operation successful. Use select_=true to add symbol to Market Watch, false to remove it.

**Signature**
```go
func (s *MT5Service) SymbolSelect(ctx context.Context, symbol string, select_ bool) (bool, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD", "GBPUSD") |
| select_ | bool | If true, add to Market Watch; if false, remove from Market Watch |

**Returns**

Returns true if operation successful, false otherwise, or error if request fails.

---

## SymbolExist

Checks if a symbol exists in the terminal. Returns two booleans: exists (whether symbol exists) and isCustom (whether it's a custom symbol). Use this before working with a symbol to avoid errors.

**Signature**
```go
func (s *MT5Service) SymbolExist(ctx context.Context, symbol string) (bool, bool, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name to check (e.g., "EURUSD") |

**Returns**

Returns `exists` (true if symbol exists), `isCustom` (true if custom symbol), and error if request fails.

---

## IsSymbolSynchronized

Checks if symbol data is synchronized with the trade server. Returns true if symbol is fully synchronized and ready for trading operations. Use this to ensure quotes are up-to-date before placing orders.

**Signature**
```go
func (s *MT5Service) IsSymbolSynchronized(ctx context.Context, symbol string) (bool, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name to check (e.g., "EURUSD") |

**Returns**

Returns true if symbol synchronized, false otherwise, or error if request fails.

---

## GetSymbolTick

Retrieves the last tick for a symbol. Returns SymbolTick struct with Go native types including time.Time for timestamp (no manual Unix conversion needed). More convenient than MT5Account.SymbolInfoTick.

**Signature**
```go
func (s *MT5Service) GetSymbolTick(ctx context.Context, symbol string) (*SymbolTick, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD", "GBPUSD") |

**Returns**

Returns `*SymbolTick` struct with current tick data (Bid, Ask, Last, Volume, Time, Spread), or error if symbol not found or request fails.

---

## GetSymbolDouble

Retrieves a double-type symbol property. Returns float64 value directly instead of wrapped protobuf. Use this for single property queries like Bid, Ask, Point, etc. For multiple properties, use GetSymbolParamsMany instead.

**Signature**
```go
func (s *MT5Service) GetSymbolDouble(ctx context.Context, symbol string, property pb.SymbolInfoDoubleProperty) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD") |
| property | pb.SymbolInfoDoubleProperty | Property ID enum (BID, ASK, POINT, VOLUME_MIN, VOLUME_MAX, VOLUME_STEP, TRADE_TICK_SIZE, etc.) |

**Returns**

Returns property value as float64, or error if symbol not found or request fails.

---

## GetSymbolInteger

Retrieves an integer-type symbol property. Returns int64 value directly instead of wrapped protobuf. Use this for single property queries like Digits, Spread, StopsLevel, etc. For multiple properties, use GetSymbolParamsMany instead.

**Signature**
```go
func (s *MT5Service) GetSymbolInteger(ctx context.Context, symbol string, property pb.SymbolInfoIntegerProperty) (int64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD") |
| property | pb.SymbolInfoIntegerProperty | Property ID enum (DIGITS, SPREAD, STOPS_LEVEL, FREEZE_LEVEL, TRADE_MODE, TRADE_EXECUTION_MODE, etc.) |

**Returns**

Returns property value as int64, or error if symbol not found or request fails.

---

## GetSymbolString

Retrieves a string-type symbol property. Returns string value directly instead of wrapped protobuf. Use this for single property queries like Description, Path, Currency, etc.

**Signature**
```go
func (s *MT5Service) GetSymbolString(ctx context.Context, symbol string, property pb.SymbolInfoStringProperty) (string, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD") |
| property | pb.SymbolInfoStringProperty | Property ID enum (DESCRIPTION, CURRENCY_BASE, CURRENCY_PROFIT, CURRENCY_MARGIN, PATH, etc.) |

**Returns**

Returns property value as string, or error if symbol not found or request fails.

---

## GetSymbolMarginRate

Retrieves margin rates for a symbol and order type. Returns clean SymbolMarginRate struct instead of nested protobuf fields. Use this to understand margin requirements before placing orders.

**Signature**
```go
func (s *MT5Service) GetSymbolMarginRate(ctx context.Context, symbol string, orderType pb.ENUM_ORDER_TYPE) (*SymbolMarginRate, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD") |
| orderType | pb.ENUM_ORDER_TYPE | Order type enum (ORDER_TYPE_BUY, ORDER_TYPE_SELL, ORDER_TYPE_BUY_LIMIT, etc.) |

**Returns**

Returns `*SymbolMarginRate` struct with InitialMarginRate and MaintenanceMarginRate, or error if request fails.

---

## GetSymbolSessionQuote

Retrieves quote session times for a symbol. Returns SessionTime struct with time.Time fields instead of protobuf timestamps (no manual .AsTime() conversions needed). Use this to check when quotes are available.

**Signature**
```go
func (s *MT5Service) GetSymbolSessionQuote(ctx context.Context, symbol string, dayOfWeek pb.DayOfWeek, sessionIndex uint32) (*SessionTime, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD") |
| dayOfWeek | pb.DayOfWeek | Day of week enum (SUNDAY=0 to SATURDAY=6) |
| sessionIndex | uint32 | Session index (0-based, usually 0 for main session) |

**Returns**

Returns `*SessionTime` struct with From and To times, or error if session not found or request fails.

---

## GetSymbolSessionTrade

Retrieves trade session times for a symbol. Returns SessionTime struct with time.Time fields. Similar to GetSymbolSessionQuote but for trading sessions instead of quote sessions.

**Signature**
```go
func (s *MT5Service) GetSymbolSessionTrade(ctx context.Context, symbol string, dayOfWeek pb.DayOfWeek, sessionIndex uint32) (*SessionTime, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (e.g., "EURUSD") |
| dayOfWeek | pb.DayOfWeek | Day of week enum (SUNDAY=0 to SATURDAY=6) |
| sessionIndex | uint32 | Session index (0-based, usually 0 for main session) |

**Returns**

Returns `*SessionTime` struct with From and To times, or error if session not found or request fails.

---

## GetSymbolParamsMany

Retrieves comprehensive parameters for multiple symbols. This is the RECOMMENDED method for getting symbol information - returns all important parameters in one call with clean Go structs. Supports filtering, sorting, and pagination for large symbol lists.

**Signature**
```go
func (s *MT5Service) GetSymbolParamsMany(ctx context.Context, symbolName *string, sortType *pb.AH_SYMBOL_PARAMS_MANY_SORT_TYPE, pageNumber *int32, itemsPerPage *int32) ([]SymbolParams, int32, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbolName | *string | Optional filter by symbol name (nil for all symbols) |
| sortType | *pb.AH_SYMBOL_PARAMS_MANY_SORT_TYPE | Optional sort type enum (nil for default sorting) |
| pageNumber | *int32 | Optional page number for pagination (nil for page 1) |
| itemsPerPage | *int32 | Optional items per page (nil for all items) |

**Returns**

Returns slice of `SymbolParams` structs, total count of symbols matching filter, or error if request fails.

---

## 📈 Market Depth

## SubscribeMarketDepth

Subscribes to Depth of Market (DOM) updates for a symbol. Must be called before GetMarketDepth. Returns true if subscription successful. Use this to start receiving Level 2 market data.

**Signature**
```go
func (s *MT5Service) SubscribeMarketDepth(ctx context.Context, symbol string) (bool, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name to subscribe to (e.g., "EURUSD") |

**Returns**

Returns true if subscription successful, false otherwise, or error if request fails.

---

## GetMarketDepth

Retrieves current Depth of Market (DOM) snapshot for a symbol. Requires prior SubscribeMarketDepth call. Returns slice of BookInfo entries with price levels, volumes, and order types.

**Signature**
```go
func (s *MT5Service) GetMarketDepth(ctx context.Context, symbol string) ([]BookInfo, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name (must be subscribed first) |

**Returns**

Returns slice of `BookInfo` entries with Type (buy/sell), Price, and Volume for each level, or error if not subscribed or request fails.

---

## UnsubscribeMarketDepth

Unsubscribes from Depth of Market updates. Call this to stop receiving DOM updates and free resources. Returns true if unsubscription successful.

**Signature**
```go
func (s *MT5Service) UnsubscribeMarketDepth(ctx context.Context, symbol string) (bool, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| symbol | string | Symbol name to unsubscribe from |

**Returns**

Returns true if unsubscription successful, false otherwise, or error if request fails.

---

## 📊 Positions & Orders

## GetPositionsTotal

Returns the count of open positions. Returns int32 directly instead of wrapped protobuf. More convenient than MT5Account.PositionsTotal for quick position counting.

**Signature**
```go
func (s *MT5Service) GetPositionsTotal(ctx context.Context) (int32, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |

**Returns**

Returns number of open positions as int32, or error if request fails.

---

## GetOpenedOrders

Retrieves all open positions and pending orders. Returns protobuf data directly (no unnecessary mapping). This matches the C# MT5Service approach for maximum compatibility.

**Signature**
```go
func (s *MT5Service) GetOpenedOrders(ctx context.Context, sortMode pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE) (*pb.OpenedOrdersData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| sortMode | pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE | Sort mode enum for ordering results |

**Returns**

Returns `*pb.OpenedOrdersData` containing positions and pending orders with full details, or error if request fails.

---

## GetOpenedTickets

Retrieves ticket numbers of open positions and pending orders. This is a lightweight alternative to GetOpenedOrders when you only need ticket IDs. Returns two slices: position tickets and order tickets.

**Signature**
```go
func (s *MT5Service) GetOpenedTickets(ctx context.Context) ([]int64, []int64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |

**Returns**

Returns `positionTickets` (slice of open position ticket numbers), `orderTickets` (slice of pending order ticket numbers), and error if request fails.

---

## ⚖️ Order Calculations & Validation

## CalculateMargin

Calculates required margin for a potential order. Use this before placing orders to check if you have enough free margin. Returns margin amount as float64 directly.

**Signature**
```go
func (s *MT5Service) CalculateMargin(ctx context.Context, req *pb.OrderCalcMarginRequest) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCalcMarginRequest | Request with Symbol, OrderType, Volume, and OpenPrice |

**Returns**

Returns required margin in account currency as float64, or error if calculation fails.

---

## CalculateProfit

Calculates potential profit for a hypothetical order. Useful for profit/risk calculations before placing actual orders. Returns profit amount as float64 directly.

**Signature**
```go
func (s *MT5Service) CalculateProfit(ctx context.Context, req *pb.OrderCalcProfitRequest) (float64, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCalcProfitRequest | Request with OrderType, Symbol, Volume, OpenPrice, and ClosePrice |

**Returns**

Returns estimated profit in account currency as float64, or error if calculation fails.

---

## CheckOrder

Validates an order before sending it to the broker. Returns clean OrderCheckResult struct instead of nested protobuf. Use this to validate orders before PlaceOrder to avoid rejections.

**Signature**
```go
func (s *MT5Service) CheckOrder(ctx context.Context, req *pb.OrderCheckRequest) (*OrderCheckResult, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCheckRequest | Request with MqlTradeRequest structure containing trading parameters |

**Returns**

Returns `*OrderCheckResult` with validation results including margin requirements and error codes, or error if request fails.

---

## 📜 History

## GetOrderHistory

Retrieves historical orders and deals for a time period with pagination. Returns protobuf data directly (no unnecessary mapping). This matches the C# MT5Service approach.

**Signature**
```go
func (s *MT5Service) GetOrderHistory(ctx context.Context, from time.Time, to time.Time, sortMode pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE, pageNumber int32, itemsPerPage int32) (*pb.OrdersHistoryData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| from | time.Time | Start time of history range |
| to | time.Time | End time of history range |
| sortMode | pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE | Sort mode enum for ordering results |
| pageNumber | int32 | Page number for pagination (1-based) |
| itemsPerPage | int32 | Number of items per page |

**Returns**

Returns `*pb.OrdersHistoryData` containing orders and deals with execution details, or error if request fails.

---

## GetPositionsHistory

Retrieves closed positions with aggregated P&L data. Returns protobuf data directly (no unnecessary mapping). Supports filtering by time range and pagination.

**Signature**
```go
func (s *MT5Service) GetPositionsHistory(ctx context.Context, sortType pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE, from *time.Time, to *time.Time, pageNumber *int32, itemsPerPage *int32) (*pb.PositionsHistoryData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| sortType | pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE | Sort type enum for ordering results |
| from | *time.Time | Optional start time for position open time filter (nil for no filter) |
| to | *time.Time | Optional end time for position open time filter (nil for no filter) |
| pageNumber | *int32 | Optional page number for pagination (nil for page 1) |
| itemsPerPage | *int32 | Optional items per page (nil for all items) |

**Returns**

Returns `*pb.PositionsHistoryData` containing closed positions with profit/loss details, or error if request fails.

---

## 🛒 Trading Operations

## PlaceOrder

Sends a market or pending order to MT5 terminal. Returns clean OrderResult struct with Go native types instead of nested protobuf. Check ReturnedCode for success (10009 means success).

**Signature**
```go
func (s *MT5Service) PlaceOrder(ctx context.Context, req *pb.OrderSendRequest) (*OrderResult, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderSendRequest | Request with Symbol, Operation, Volume, Price, Slippage, StopLoss, TakeProfit, Comment, ExpertId, StopLimitPrice, ExpirationTimeType, and ExpirationTime |

**Returns**

Returns `*OrderResult` struct with execution details (Ticket, Price, Volume, ReturnedCode, Comment), or error if request fails.

---

## ModifyOrder

Modifies an existing order or position (change SL/TP/price). Returns OrderResult with modification details. Check ReturnedCode for success (10009 means success).

**Signature**
```go
func (s *MT5Service) ModifyOrder(ctx context.Context, req *pb.OrderModifyRequest) (*OrderResult, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderModifyRequest | Request with Ticket, Price (for pending orders), StopLoss, TakeProfit, and optional Expiration |

**Returns**

Returns `*OrderResult` struct with modification status and details, or error if request fails.

---

## CloseOrder

Closes a position or deletes a pending order. Returns operation return code (10009 = success). Simpler than PlaceOrder for closing operations.

**Signature**
```go
func (s *MT5Service) CloseOrder(ctx context.Context, req *pb.OrderCloseRequest) (uint32, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCloseRequest | Request with Ticket, Volume, and Slippage |

**Returns**

Returns operation return code as uint32 (10009 = success), or error if request fails.

---

## 📡 Streaming Methods

## StreamTicks

Streams real-time ticks for specified symbols. Returns channel of SymbolTick structs with Go native types instead of protobuf (time.Time fields, no manual conversions). The channels close when streaming stops.

**Signature**
```go
func (s *MT5Service) StreamTicks(ctx context.Context, symbols []string) (<-chan *SymbolTick, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for cancellation (closing ctx stops the stream) |
| symbols | []string | Slice of symbol names to stream (e.g., []string{"EURUSD", "GBPUSD"}) |

**Returns**

Returns read-only channel of `*SymbolTick` structs with tick data, and read-only channel of errors. Always read from both channels in a select statement.

---

## StreamPositionProfits

Streams real-time profit/loss updates for open positions. Provides P&L updates as market moves, new positions opened, positions modified, and positions closed. The channels close when streaming stops.

**Signature**
```go
func (s *MT5Service) StreamPositionProfits(ctx context.Context) (<-chan *pb.OnPositionProfitData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for cancellation (closing ctx stops the stream) |

**Returns**

Returns read-only channel of `*pb.OnPositionProfitData` with profit updates, and read-only channel of errors. Always read from both channels in a select statement.

---

## StreamOpenedTickets

Streams updates to the list of open position and pending order tickets. This is a lightweight alternative to StreamTradeUpdates when you only need ticket IDs. The channels close when streaming stops.

**Signature**
```go
func (s *MT5Service) StreamOpenedTickets(ctx context.Context) (<-chan *pb.OnPositionsAndPendingOrdersTicketsData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for cancellation (closing ctx stops the stream) |

**Returns**

Returns read-only channel of `*pb.OnPositionsAndPendingOrdersTicketsData` with ticket lists, and read-only channel of errors. Always read from both channels in a select statement.

---

## StreamTradeUpdates

Streams trade events (new/disappeared orders and positions, history updates). Provides real-time notifications about orders and positions opened, orders and positions closed/cancelled, and history orders and deals. The channels close when streaming stops.

**Signature**
```go
func (s *MT5Service) StreamTradeUpdates(ctx context.Context) (<-chan *pb.OnTradeData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for cancellation (closing ctx stops the stream) |

**Returns**

Returns read-only channel of `*pb.OnTradeData` with trade events, and read-only channel of errors. Always read from both channels in a select statement.

---

## StreamTransactions

Streams all trade transaction events (most detailed streaming method). This is the most comprehensive trade event stream providing detailed transaction information about order placement, modification, deletion, deal execution, position opening, modification, closing, and all intermediate states.

**Signature**
```go
func (s *MT5Service) StreamTransactions(ctx context.Context) (<-chan *pb.OnTradeTransactionData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for cancellation (closing ctx stops the stream) |

**Returns**

Returns read-only channel of `*pb.OnTradeTransactionData` with transaction details, and read-only channel of errors. Always read from both channels in a select statement.

---
