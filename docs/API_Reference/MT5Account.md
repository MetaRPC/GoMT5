# MT5Account API Reference

> **Note:** This documentation is auto-generated from [`examples/mt5/MT5Account.go`](../../examples/mt5/MT5Account.go) with enhanced navigation for easier browsing. For those who prefer viewing the complete API in a single-page reference format.

MT5Account represents a low-level gRPC client for MetaTrader 5 terminal. All methods accept protobuf Request objects and return protobuf Data objects. This is the foundation layer that directly communicates with the MT5 gRPC server.

## Table of Contents

### 🏗️ Constructor & Connection
- [NewMT5Account](#newmt5account)
- [Connect](#connect)
- [ConnectEx](#connectex)
- [ConnectProxy](#connectproxy)
- [Disconnect](#disconnect)
- [Reconnect](#reconnect)
- [CheckConnect](#checkconnect)
- [IsConnected](#isconnected)
- [Close](#close)

### 💰 Account Information
- [AccountSummary](#accountsummary)
- [AccountInfoDouble](#accountinfodouble)
- [AccountInfoInteger](#accountinfointeger)
- [AccountInfoString](#accountinfostring)

### 🔢 Symbol Information
- [SymbolsTotal](#symbolstotal)
- [SymbolName](#symbolname)
- [SymbolSelect](#symbolselect)
- [SymbolExist](#symbolexist)
- [SymbolIsSynchronized](#symbolissynchronized)
- [SymbolInfoTick](#symbolinfotick)
- [SymbolInfoDouble](#symbolinfodouble)
- [SymbolInfoInteger](#symbolinfointeger)
- [SymbolInfoString](#symbolinfostring)
- [SymbolInfoMarginRate](#symbolinfomarginrate)
- [SymbolInfoSessionQuote](#symbolinfosessionquote)
- [SymbolInfoSessionTrade](#symbolinfosessiontrade)
- [SymbolParamsMany](#symbolparamsmany)

### 📊 Positions & Orders
- [PositionsTotal](#positionstotal)
- [OpenedOrders](#openedorders)
- [OpenedOrdersTickets](#openedorderstickets)
- [OrderCalcMargin](#ordercalcmargin)
- [OrderCalcProfit](#ordercalcprofit)
- [OrderCheck](#ordercheck)

### 📜 History
- [OrderHistory](#orderhistory)
- [PositionsHistory](#positionshistory)

### 📈 Market Depth
- [MarketBookAdd](#marketbookadd)
- [MarketBookGet](#marketbookget)
- [MarketBookRelease](#marketbookrelease)

### 🛒 Trading Operations
- [OrderSend](#ordersend)
- [OrderModify](#ordermodify)
- [OrderClose](#orderclose)

### 📡 Streaming Methods
- [OnSymbolTick](#onsymboltick)
- [OnPositionProfit](#onpositionprofit)
- [OnPositionsAndPendingOrdersTickets](#onpositionsandpendingorderstickets)
- [OnTrade](#ontrade)
- [OnTradeTransaction](#ontradetransaction)

---

## 🏗️ Constructor & Connection

## NewMT5Account

Creates a new MT5Account instance with gRPC connection. Default grpcServer is "mt5.mrpc.pro:443" if empty string is provided. The connection is established with TLS, keepalive, and automatic reconnect configured.

**Signature**
```go
func NewMT5Account(user uint64, password string, grpcServer string, id uuid.UUID) (*MT5Account, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| user | uint64 | MT5 account login number |
| password | string | MT5 account password |
| grpcServer | string | gRPC server address (e.g., "mt5.mrpc.pro:443", empty for default) |
| id | uuid.UUID | Unique identifier for this account instance |

**Returns**

Returns `*MT5Account` instance with established gRPC connection, or error if connection fails.

---

## Connect

Establishes basic connection to MT5 terminal. This is a simplified connection method that uses default settings. For advanced connection configuration, use ConnectEx instead.

**Signature**
```go
func (a *MT5Account) Connect(ctx context.Context, req *pb.ConnectRequest) (*pb.ConnectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.ConnectRequest | Request with User, Password, and optional ServerName |

**Returns**

Returns `*pb.ConnectData` with session UUID and connection status, or error on failure.

---

## ConnectEx

Establishes connection to MT5 terminal with extended parameters. This method provides full control over connection settings including MT5 cluster name for connection, connection timeout settings, base chart symbol selection, and Expert Advisors to add.

**Signature**
```go
func (a *MT5Account) ConnectEx(ctx context.Context, req *pb.ConnectExRequest) (*pb.ConnectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.ConnectExRequest | Request with User, Password, MtClusterName, BaseChartSymbol, TerminalReadinessWaitingTimeoutSeconds, ExpertsToAdd |

**Returns**

Returns `*pb.ConnectData` with session UUID and connection status, or error on failure.

---

## ConnectProxy

Establishes connection to MT5 terminal through proxy server. Use this method when MT5 terminal access requires proxy configuration.

**Signature**
```go
func (a *MT5Account) ConnectProxy(ctx context.Context, req *pb.ConnectProxyRequest) (*pb.ConnectProxyData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.ConnectProxyRequest | Request with proxy host, port, type, credentials and MT5 account credentials |

**Returns**

Returns `*pb.ConnectProxyData` with session UUID and connection status, or error on failure.

---

## Disconnect

Closes the connection to MT5 terminal. This method gracefully terminates the active MT5 session.

**Signature**
```go
func (a *MT5Account) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.DisconnectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.DisconnectRequest | Empty request structure |

**Returns**

Returns `*pb.DisconnectData` with disconnection status, or error on failure.

---

## Reconnect

Re-establishes connection to MT5 terminal. This method recreates the terminal session without changing connection parameters. Used internally by ExecuteWithReconnect on TERMINAL_INSTANCE_NOT_FOUND errors.

**Signature**
```go
func (a *MT5Account) Reconnect(ctx context.Context, req *pb.ReconnectRequest) (*pb.ReconnectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.ReconnectRequest | Empty request structure |

**Returns**

Returns `*pb.ReconnectData` with new session UUID, or error on failure.

---

## CheckConnect

Verifies the current connection status to MT5 terminal. Use this method to ping the terminal and confirm the session is still active.

**Signature**
```go
func (a *MT5Account) CheckConnect(ctx context.Context, req *pb.CheckConnectRequest) (*pb.CheckConnectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.CheckConnectRequest | Empty request structure |

**Returns**

Returns `*pb.CheckConnectData` with connection status flag, or error on failure.

---

## IsConnected

Returns true if the account has an active gRPC connection.

**Signature**
```go
func (a *MT5Account) IsConnected() bool
```

**Returns**

Returns true if gRPC connection is active, false otherwise.

---

## Close

Closes the gRPC connection and cleans up resources.

**Signature**
```go
func (a *MT5Account) Close() error
```

**Returns**

Returns error if cleanup fails, nil on success.

---

## 💰 Account Information

## AccountSummary

Retrieves all account information in one call. This is the recommended method for getting account data as it returns all properties in a single request, avoiding multiple round-trips.

**Signature**
```go
func (a *MT5Account) AccountSummary(ctx context.Context, req *pb.AccountSummaryRequest) (*pb.AccountSummaryData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.AccountSummaryRequest | Empty request structure |

**Returns**

Returns `*pb.AccountSummaryData` with Login, Balance, Equity, UserName, Leverage, TradeMode, CompanyName, Currency, ServerTime, UtcTimezoneShift, and Credit.

---

## AccountInfoDouble

Retrieves a double-type account property. Use this method when you need a specific numeric account property. For multiple properties, use AccountSummary instead.

**Signature**
```go
func (a *MT5Account) AccountInfoDouble(ctx context.Context, req *pb.AccountInfoDoubleRequest) (*pb.AccountInfoDoubleData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.AccountInfoDoubleRequest | Request with property_id (ACCOUNT_BALANCE, ACCOUNT_EQUITY, ACCOUNT_MARGIN, ACCOUNT_MARGIN_FREE, ACCOUNT_PROFIT, etc) |

**Returns**

Returns `*pb.AccountInfoDoubleData` with the requested double value.

---

## AccountInfoInteger

Retrieves an integer-type account property. Use this method when you need a specific integer account property.

**Signature**
```go
func (a *MT5Account) AccountInfoInteger(ctx context.Context, req *pb.AccountInfoIntegerRequest) (*pb.AccountInfoIntegerData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.AccountInfoIntegerRequest | Request with property_id (ACCOUNT_LOGIN, ACCOUNT_LEVERAGE, ACCOUNT_LIMIT_ORDERS, ACCOUNT_TRADE_MODE, ACCOUNT_MARGIN_SO_MODE, etc) |

**Returns**

Returns `*pb.AccountInfoIntegerData` with the requested int64 value.

---

## AccountInfoString

Retrieves a string-type account property. Use this method when you need a specific string account property.

**Signature**
```go
func (a *MT5Account) AccountInfoString(ctx context.Context, req *pb.AccountInfoStringRequest) (*pb.AccountInfoStringData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.AccountInfoStringRequest | Request with property_id (ACCOUNT_NAME, ACCOUNT_SERVER, ACCOUNT_CURRENCY, ACCOUNT_COMPANY) |

**Returns**

Returns `*pb.AccountInfoStringData` with the requested string value.

---

## 🔢 Symbol Information

## SymbolsTotal

Returns the number of available symbols. Use this method to count symbols either in Market Watch or all available symbols.

**Signature**
```go
func (a *MT5Account) SymbolsTotal(ctx context.Context, req *pb.SymbolsTotalRequest) (*pb.SymbolsTotalData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolsTotalRequest | Request with Selected flag (true for Market Watch symbols only, false for all symbols) |

**Returns**

Returns `*pb.SymbolsTotalData` with total count of symbols.

---

## SymbolName

Returns the name of a symbol by its position in the list. Use this method to iterate through available symbols by index.

**Signature**
```go
func (a *MT5Account) SymbolName(ctx context.Context, req *pb.SymbolNameRequest) (*pb.SymbolNameData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolNameRequest | Request with Pos (zero-based index) and Selected flag (true for Market Watch, false for all symbols) |

**Returns**

Returns `*pb.SymbolNameData` with symbol name at the specified position.

---

## SymbolSelect

Adds or removes a symbol from Market Watch window. Use this method to manage which symbols are visible in Market Watch.

**Signature**
```go
func (a *MT5Account) SymbolSelect(ctx context.Context, req *pb.SymbolSelectRequest) (*pb.SymbolSelectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolSelectRequest | Request with Symbol name and Select flag (true to add, false to remove) |

**Returns**

Returns `*pb.SymbolSelectData` with success status of the operation.

---

## SymbolExist

Checks if a symbol with specified name exists. Use this method to verify symbol availability before requesting data or placing orders.

**Signature**
```go
func (a *MT5Account) SymbolExist(ctx context.Context, req *pb.SymbolExistRequest) (*pb.SymbolExistData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolExistRequest | Request with Symbol name |

**Returns**

Returns `*pb.SymbolExistData` with Exist flag and IsCustom flag indicating if it's a custom symbol.

---

## SymbolIsSynchronized

Checks if symbol data is synchronized with trade server. Use this method to ensure symbol quotes are up-to-date before trading operations.

**Signature**
```go
func (a *MT5Account) SymbolIsSynchronized(ctx context.Context, req *pb.SymbolIsSynchronizedRequest) (*pb.SymbolIsSynchronizedData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolIsSynchronizedRequest | Request with Symbol name |

**Returns**

Returns `*pb.SymbolIsSynchronizedData` with IsSynchronized flag.

---

## SymbolInfoTick

Retrieves the last tick data for a symbol. Use this method to get the most recent price update with timestamp.

**Signature**
```go
func (a *MT5Account) SymbolInfoTick(ctx context.Context, req *pb.SymbolInfoTickRequest) (*pb.MrpcMqlTick, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolInfoTickRequest | Request with Symbol name |

**Returns**

Returns `*pb.MrpcMqlTick` with Bid, Ask, Last, Volume, Time, TimeMS, Flags, VolumReal and spread values.

---

## SymbolInfoDouble

Retrieves a double-type symbol property. Use this method to get numeric symbol properties like prices, volumes, and trading parameters.

**Signature**
```go
func (a *MT5Account) SymbolInfoDouble(ctx context.Context, req *pb.SymbolInfoDoubleRequest) (*pb.SymbolInfoDoubleData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolInfoDoubleRequest | Request with Symbol name and PropertyId (BID, ASK, POINT, VOLUME_MIN, VOLUME_MAX, VOLUME_STEP, TRADE_TICK_SIZE, etc) |

**Returns**

Returns `*pb.SymbolInfoDoubleData` with the requested double value.

---

## SymbolInfoInteger

Retrieves an integer-type symbol property. Use this method to get integer symbol properties like digits, spread, and trading restrictions.

**Signature**
```go
func (a *MT5Account) SymbolInfoInteger(ctx context.Context, req *pb.SymbolInfoIntegerRequest) (*pb.SymbolInfoIntegerData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolInfoIntegerRequest | Request with Symbol name and PropertyId (DIGITS, SPREAD, STOPS_LEVEL, FREEZE_LEVEL, TRADE_MODE, TRADE_EXECUTION_MODE, etc) |

**Returns**

Returns `*pb.SymbolInfoIntegerData` with the requested int64 value.

---

## SymbolInfoString

Retrieves a string-type symbol property. Use this method to get text symbol properties like description and currency information.

**Signature**
```go
func (a *MT5Account) SymbolInfoString(ctx context.Context, req *pb.SymbolInfoStringRequest) (*pb.SymbolInfoStringData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolInfoStringRequest | Request with Symbol name and PropertyId (DESCRIPTION, CURRENCY_BASE, CURRENCY_PROFIT, CURRENCY_MARGIN, PATH, etc) |

**Returns**

Returns `*pb.SymbolInfoStringData` with the requested string value.

---

## SymbolInfoMarginRate

Retrieves margin requirements for different order types. Use this method to calculate margin before placing orders.

**Signature**
```go
func (a *MT5Account) SymbolInfoMarginRate(ctx context.Context, req *pb.SymbolInfoMarginRateRequest) (*pb.SymbolInfoMarginRateData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolInfoMarginRateRequest | Request with Symbol name and OrderType (ORDER_TYPE_BUY, ORDER_TYPE_SELL, etc) |

**Returns**

Returns `*pb.SymbolInfoMarginRateData` with InitialMarginRate and MaintenanceMarginRate values.

---

## SymbolInfoSessionQuote

Retrieves quote session times for a symbol. Use this method to check when quotes are available for trading.

**Signature**
```go
func (a *MT5Account) SymbolInfoSessionQuote(ctx context.Context, req *pb.SymbolInfoSessionQuoteRequest) (*pb.SymbolInfoSessionQuoteData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolInfoSessionQuoteRequest | Request with Symbol name, DayOfWeek (SUNDAY=0 to SATURDAY=6) and SessionIndex |

**Returns**

Returns `*pb.SymbolInfoSessionQuoteData` with session From and To times in seconds from day start.

---

## SymbolInfoSessionTrade

Retrieves trade session times for a symbol. Use this method to check when trading operations are allowed.

**Signature**
```go
func (a *MT5Account) SymbolInfoSessionTrade(ctx context.Context, req *pb.SymbolInfoSessionTradeRequest) (*pb.SymbolInfoSessionTradeData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolInfoSessionTradeRequest | Request with Symbol name, DayOfWeek (SUNDAY=0 to SATURDAY=6) and SessionIndex |

**Returns**

Returns `*pb.SymbolInfoSessionTradeData` with session From and To times in seconds from day start.

---

## SymbolParamsMany

Retrieves detailed parameters for multiple symbols in one call. This is the recommended method for getting comprehensive symbol data as it returns all properties for multiple symbols in a single request, avoiding multiple round-trips.

**Signature**
```go
func (a *MT5Account) SymbolParamsMany(ctx context.Context, req *pb.SymbolParamsManyRequest) (*pb.SymbolParamsManyData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.SymbolParamsManyRequest | Request with array of Symbol names |

**Returns**

Returns `*pb.SymbolParamsManyData` with array of SymbolParams containing Bid, Ask, Digits, Spread, VolumeMin, VolumeMax, VolumeStep, ContractSize, Point, margins, and other trading parameters for each symbol.

---

## 📊 Positions & Orders

## PositionsTotal

Returns the number of currently open positions. Use this method for quick check of open positions count without retrieving full details.

**Signature**
```go
func (a *MT5Account) PositionsTotal(ctx context.Context) (*pb.PositionsTotalData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |

**Returns**

Returns `*pb.PositionsTotalData` with Total count of open positions.

---

## OpenedOrders

Retrieves all currently opened orders and positions with full details. This method returns comprehensive information about all active trading positions and pending orders including profit/loss, prices, volumes, and timestamps.

**Signature**
```go
func (a *MT5Account) OpenedOrders(ctx context.Context, req *pb.OpenedOrdersRequest) (*pb.OpenedOrdersData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OpenedOrdersRequest | Request with InputSortMode (sort type: by open time, close time, or ticket ID) |

**Returns**

Returns `*pb.OpenedOrdersData` with arrays of opened_orders (pending orders) and position_infos (open positions) containing full details.

---

## OpenedOrdersTickets

Retrieves only ticket numbers of currently opened orders and positions. This is a lightweight alternative to OpenedOrders when you only need ticket IDs for subsequent operations or monitoring.

**Signature**
```go
func (a *MT5Account) OpenedOrdersTickets(ctx context.Context, req *pb.OpenedOrdersTicketsRequest) (*pb.OpenedOrdersTicketsData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OpenedOrdersTicketsRequest | Empty request structure |

**Returns**

Returns `*pb.OpenedOrdersTicketsData` with arrays of opened_orders_tickets and opened_position_tickets.

---

## OrderCalcMargin

Calculates required margin for an order. Use this method to determine how much margin will be required before placing an order.

**Signature**
```go
func (a *MT5Account) OrderCalcMargin(ctx context.Context, req *pb.OrderCalcMarginRequest) (*pb.OrderCalcMarginData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCalcMarginRequest | Request with Symbol, OrderType, Volume, and OpenPrice |

**Returns**

Returns `*pb.OrderCalcMarginData` with Margin value in account currency.

---

## OrderCalcProfit

Calculates potential profit for a trade. Use this method to estimate profit/loss before placing an order or to calculate current profit at a specified price level.

**Signature**
```go
func (a *MT5Account) OrderCalcProfit(ctx context.Context, req *pb.OrderCalcProfitRequest) (*pb.OrderCalcProfitData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCalcProfitRequest | Request with OrderType, Symbol, Volume, OpenPrice, and ClosePrice |

**Returns**

Returns `*pb.OrderCalcProfitData` with Profit value in account currency.

---

## OrderCheck

Validates an order before sending it to the server. Use this method to pre-validate trading requests without actually placing orders. Useful for checking margin requirements and detecting potential errors.

**Signature**
```go
func (a *MT5Account) OrderCheck(ctx context.Context, req *pb.OrderCheckRequest) (*pb.OrderCheckData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCheckRequest | Request with Action, Symbol, Volume, Price, OrderType, StopLoss, TakeProfit and other order parameters |

**Returns**

Returns `*pb.OrderCheckData` with validation result including margin requirements, estimated profit, and MqlTradeCheckResult structure with validation status and possible error codes.

---

## 📜 History

## OrderHistory

Retrieves historical orders within a specified time range. Use this method to analyze past order activity with pagination support for large datasets.

**Signature**
```go
func (a *MT5Account) OrderHistory(ctx context.Context, req *pb.OrderHistoryRequest) (*pb.OrdersHistoryData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderHistoryRequest | Request with FromDate, ToDate (unix timestamps), optional Symbol filter, Offset and Limit for pagination |

**Returns**

Returns `*pb.OrdersHistoryData` with array of historical Order objects including execution details, prices, volumes, and final status.

---

## PositionsHistory

Retrieves closed positions with profit/loss information. Use this method to analyze trading performance and calculate statistics for closed positions within a time range.

**Signature**
```go
func (a *MT5Account) PositionsHistory(ctx context.Context, req *pb.PositionsHistoryRequest) (*pb.PositionsHistoryData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.PositionsHistoryRequest | Request with FromDate, ToDate (unix timestamps), optional Symbol filter, Offset and Limit for pagination |

**Returns**

Returns `*pb.PositionsHistoryData` with array of closed Position objects including entry/exit prices, volumes, swap, commission, net profit, and close timestamps.

---

## 📈 Market Depth

## MarketBookAdd

Subscribes to Depth of Market (DOM) updates for a symbol. Use this method to start receiving Level 2 market data with bid/ask prices and volumes at different price levels.

**Signature**
```go
func (a *MT5Account) MarketBookAdd(ctx context.Context, req *pb.MarketBookAddRequest) (*pb.MarketBookAddData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.MarketBookAddRequest | Request with Symbol name |

**Returns**

Returns `*pb.MarketBookAddData` with subscription status.

---

## MarketBookGet

Retrieves current market depth snapshot for a symbol. Use this method to get the current order book state with all price levels, volumes, and order types (buy/sell) at each level.

**Signature**
```go
func (a *MT5Account) MarketBookGet(ctx context.Context, req *pb.MarketBookGetRequest) (*pb.MarketBookGetData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.MarketBookGetRequest | Request with Symbol name |

**Returns**

Returns `*pb.MarketBookGetData` with array of BookStruct entries containing Type (buy/sell), Price, Volume, and VolumeDouble for each price level in the order book.

---

## MarketBookRelease

Unsubscribes from Depth of Market (DOM) updates. Use this method to stop receiving Level 2 market data and free resources.

**Signature**
```go
func (a *MT5Account) MarketBookRelease(ctx context.Context, req *pb.MarketBookReleaseRequest) (*pb.MarketBookReleaseData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.MarketBookReleaseRequest | Request with Symbol name |

**Returns**

Returns `*pb.MarketBookReleaseData` with unsubscription status.

---

## 🛒 Trading Operations

## OrderSend

Places a market or pending order. This is the main trading method for opening positions and placing pending orders. Supports all order types: market buy/sell, buy/sell limit, buy/sell stop, and stop-limit.

**Signature**
```go
func (a *MT5Account) OrderSend(ctx context.Context, req *pb.OrderSendRequest) (*pb.OrderSendData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderSendRequest | Request with Symbol, Operation, Volume, Price, Slippage, StopLoss, TakeProfit, Comment, ExpertId, StopLimitPrice, ExpirationTimeType, and ExpirationTime |

**Returns**

Returns `*pb.OrderSendData` with returned code, deal ticket, order ticket, execution price, volume, bid/ask prices, comment, and request ID.

---

## OrderModify

Modifies an existing pending order or position. Use this method to change price levels (entry price for pending orders, StopLoss and TakeProfit for positions and pending orders).

**Signature**
```go
func (a *MT5Account) OrderModify(ctx context.Context, req *pb.OrderModifyRequest) (*pb.OrderModifyData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderModifyRequest | Request with Ticket, Price (for pending orders), StopLoss, TakeProfit, and optional Expiration |

**Returns**

Returns `*pb.OrderModifyData` with modification status and MqlTradeResult structure.

---

## OrderClose

Closes an existing market position or deletes a pending order. Use this method to exit positions at current market price or cancel pending orders.

**Signature**
```go
func (a *MT5Account) OrderClose(ctx context.Context, req *pb.OrderCloseRequest) (*pb.OrderCloseData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderCloseRequest | Request with Ticket, Volume, and Slippage |

**Returns**

Returns `*pb.OrderCloseData` with ReturnedCode, ReturnedStringCode, ReturnedCodeDescription, and CloseMode (market close, partial close, or pending order remove).

---

## 📡 Streaming Methods

## OnSymbolTick

Streams real-time tick data for a symbol. This method provides continuous price updates (Bid/Ask) as they arrive from the server. The stream automatically reconnects on connection loss.

**Signature**
```go
func (a *MT5Account) OnSymbolTick(ctx context.Context, req *pb.OnSymbolTickRequest) (<-chan *pb.OnSymbolTickData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control (cancel to stop streaming) |
| req | *pb.OnSymbolTickRequest | Request with Symbol name |

**Returns**

Returns two channels: data channel receives `*pb.OnSymbolTickData` with Bid, Ask, Last, Volume, Time for each tick; error channel receives errors if stream fails (both channels closed on context cancellation).

---

## OnPositionProfit

Streams real-time profit/loss updates for open positions. This method provides continuous P&L updates as market prices change, useful for monitoring account performance and implementing risk management.

**Signature**
```go
func (a *MT5Account) OnPositionProfit(ctx context.Context, req *pb.OnPositionProfitRequest) (<-chan *pb.OnPositionProfitData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control (cancel to stop streaming) |
| req | *pb.OnPositionProfitRequest | Request with optional Symbol filter |

**Returns**

Returns two channels: data channel receives `*pb.OnPositionProfitData` with Ticket, Symbol, Profit, and current price; error channel receives errors if stream fails (both channels closed on context cancellation).

---

## OnPositionsAndPendingOrdersTickets

Streams changes in open positions and pending orders. This method notifies whenever positions are opened/closed or pending orders are added/removed, providing only ticket numbers for efficient monitoring.

**Signature**
```go
func (a *MT5Account) OnPositionsAndPendingOrdersTickets(ctx context.Context, req *pb.OnPositionsAndPendingOrdersTicketsRequest) (<-chan *pb.OnPositionsAndPendingOrdersTicketsData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control (cancel to stop streaming) |
| req | *pb.OnPositionsAndPendingOrdersTicketsRequest | Empty request structure |

**Returns**

Returns two channels: data channel receives `*pb.OnPositionsAndPendingOrdersTicketsData` with arrays of PositionTickets and PendingOrderTickets; error channel receives errors if stream fails (both channels closed on context cancellation).

---

## OnTrade

Streams trade events in real-time. This method provides notifications about all trading operations: order placement, modification, execution, and cancellation.

**Signature**
```go
func (a *MT5Account) OnTrade(ctx context.Context, req *pb.OnTradeRequest) (<-chan *pb.OnTradeData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control (cancel to stop streaming) |
| req | *pb.OnTradeRequest | Empty request structure |

**Returns**

Returns two channels: data channel receives `*pb.OnTradeData` with trade event details; error channel receives errors if stream fails (both channels closed on context cancellation).

---

## OnTradeTransaction

Streams detailed trade transaction events. This method provides low-level notifications about every change in trading state, including order state changes, deal executions, and position modifications.

**Signature**
```go
func (a *MT5Account) OnTradeTransaction(ctx context.Context, req *pb.OnTradeTransactionRequest) (<-chan *pb.OnTradeTransactionData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control (cancel to stop streaming) |
| req | *pb.OnTradeTransactionRequest | Empty request structure |

**Returns**

Returns two channels: data channel receives `*pb.OnTradeTransactionData` with MqlTradeTransaction containing Type, OrderState, DealTicket, OrderTicket, Symbol, Price, Volume; error channel receives errors if stream fails (both channels closed on context cancellation).

---
