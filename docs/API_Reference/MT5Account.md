# MT5Account API Reference

> **Note:** Complete API reference for package/Helpers/MT5Account.go with enhanced navigation. Single-page format for easy browsing.

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
- [DisconnectAndClose](#disconnectandclose)
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
- [TickValueWithSize](#tickvaluewithsize)

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

Returns `*MT5Account` instance with established gRPC connection, or error if connection fails. The MT5 terminal itself is not connected yet - call `ConnectEx` (or `Connect`) next.

**Example**
```go
import (
    pb "github.com/MetaRPC/GoMT5/package"
    helpers "github.com/MetaRPC/GoMT5/package/Helpers"
    "github.com/google/uuid"
)

// Create account with a random UUID (the package name is mt5, imported here as helpers)
account, err := helpers.NewMT5Account(12345678, "password", "", uuid.New())
if err != nil {
    log.Fatal(err)
}
defer account.Close()

_, err = account.ConnectEx(ctx, &pb.ConnectExRequest{
    User:          12345678,
    Password:      "password",
    MtClusterName: "YourBroker-MT5 Demo",
})
```

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
| req | *pb.ConnectRequest | Request with User, Password, Host, Port, and optional ExpertsToAdd, TimeoutSeconds, Name, Expiration |

**Returns**

Returns `*pb.ConnectData` with TerminalInstanceGuid and TerminalType, or error on failure.

---

## ConnectEx

Establishes connection to MT5 terminal with extended parameters: the MT5 cluster (broker server) name, connection timeout, terminal name/expiration and Expert Advisors to add.

**Signature**
```go
func (a *MT5Account) ConnectEx(ctx context.Context, req *pb.ConnectExRequest) (*pb.ConnectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.ConnectExRequest | Request with User, Password, MtClusterName, and optional ExpertsToAdd, TimeoutSeconds, Name, Expiration |

**Returns**

Returns `*pb.ConnectData` with TerminalInstanceGuid and TerminalType, or error on failure.

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
| req | *pb.ConnectProxyRequest | Request with User, Password, Host, Port, ProxyUser, ProxyPassword, ProxyHost, ProxyPort, ProxyType, and optional ExpertsToAdd, TimeoutSeconds, Name, Expiration |

**Returns**

Returns `*pb.ConnectProxyData` with UniqueIdentifier and TerminalType, or error on failure.

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
| req | *pb.DisconnectRequest | Request with optional Reason and Delete (true = delete the terminal instance instead of stopping it) |

**Returns**

Returns `*pb.DisconnectData` with UniqueIdentifier and FullLifeTimeSeconds, or error on failure.

---

## Reconnect

Re-establishes connection to MT5 terminal. This method recreates the terminal session without changing connection parameters. (It is not called automatically: on TERMINAL_INSTANCE_NOT_FOUND / network errors ExecuteWithReconnect only retries the original call after a backoff delay.)

**Signature**
```go
func (a *MT5Account) Reconnect(ctx context.Context, req *pb.ReconnectRequest) (*pb.ReconnectData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.ReconnectRequest | Request with optional ForceReconnection |

**Returns**

Returns `*pb.ReconnectData` with TerminalWasRecreated, TerminalType, TerminalInstanceGuid, OldTerminalFullLiveTimeSeconds and OldLogFiles, or error on failure.

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

Returns `*pb.CheckConnectData` with UniqueIdentifier and HealthCheck (`HealthCheck.IsAlive` is the connection status flag), or error on failure.

---

## IsConnected

Returns true if the account has an active gRPC connection and a session id (TerminalInstanceGuid or Id). It is a local check - it does not call the server (use CheckConnect for that).

**Signature**
```go
func (a *MT5Account) IsConnected() bool
```

**Returns**

Returns true if gRPC connection is active, false otherwise.

---

## DisconnectAndClose

Sends a Disconnect RPC (5-second timeout) and then closes the gRPC connection. If `deleteOnDisconnect` is true, the terminal instance is deleted instead of stopped.

**Signature**
```go
func (a *MT5Account) DisconnectAndClose(deleteOnDisconnect ...bool) error
```

**Returns**

Returns the Disconnect error if any, otherwise the Close error; nil on success.

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

Returns `*pb.AccountSummaryData` with AccountLogin, AccountBalance, AccountEquity, AccountUserName, AccountLeverage, AccountTradeMode, AccountCompanyName, AccountCurrency, ServerTime (+ ServerTimeUtc), UtcTimezoneServerTimeShiftMinutes, AccountCredit and ServerUtcOffsetSeconds.

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
| req | *pb.AccountInfoDoubleRequest | Request with PropertyId (ACCOUNT_BALANCE, ACCOUNT_EQUITY, ACCOUNT_MARGIN, ACCOUNT_MARGIN_FREE, ACCOUNT_PROFIT, etc) |

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
| req | *pb.AccountInfoIntegerRequest | Request with PropertyId (ACCOUNT_LOGIN, ACCOUNT_LEVERAGE, ACCOUNT_LIMIT_ORDERS, ACCOUNT_TRADE_MODE, ACCOUNT_MARGIN_SO_MODE, etc) |

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
| req | *pb.AccountInfoStringRequest | Request with PropertyId (ACCOUNT_NAME, ACCOUNT_SERVER, ACCOUNT_CURRENCY, ACCOUNT_COMPANY) |

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
| req | *pb.SymbolsTotalRequest | Request with Mode flag (true for Market Watch symbols only, false for all symbols) |

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
| req | *pb.SymbolNameRequest | Request with Index (zero-based) and Selected flag (true for Market Watch, false for all symbols) |

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
| req | *pb.SymbolExistRequest | Request with Name (symbol name) |

**Returns**

Returns `*pb.SymbolExistData` with Exists flag and IsCustom flag indicating if it's a custom symbol.

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

Returns `*pb.SymbolIsSynchronizedData` with Synchronized flag.

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

Returns `*pb.MrpcMqlTick` with Time (Unix seconds), Bid, Ask, Last, Volume, TimeMsc, Flags and VolumeReal (plus TimeUtc / TimeMscUtc).

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
| req | *pb.SymbolInfoDoubleRequest | Request with Symbol name and Type (SYMBOL_BID, SYMBOL_ASK, SYMBOL_POINT, SYMBOL_VOLUME_MIN, SYMBOL_VOLUME_MAX, SYMBOL_VOLUME_STEP, SYMBOL_TRADE_TICK_SIZE, etc) |

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
| req | *pb.SymbolInfoIntegerRequest | Request with Symbol name and Type (SYMBOL_DIGITS, SYMBOL_SPREAD, SYMBOL_TRADE_STOPS_LEVEL, SYMBOL_TRADE_FREEZE_LEVEL, SYMBOL_TRADE_MODE, SYMBOL_TRADE_EXEMODE, etc) |

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
| req | *pb.SymbolInfoStringRequest | Request with Symbol name and Type (SYMBOL_DESCRIPTION, SYMBOL_CURRENCY_BASE, SYMBOL_CURRENCY_PROFIT, SYMBOL_CURRENCY_MARGIN, SYMBOL_PATH, etc) |

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

Returns `*pb.SymbolInfoSessionQuoteData` with session From and To (Timestamps; only the time of day is meaningful, the date part should be ignored).

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

Returns `*pb.SymbolInfoSessionTradeData` with session From and To (Timestamps; only the time of day is meaningful, the date part should be ignored).

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
| req | *pb.SymbolParamsManyRequest | Request with optional SymbolName filter, SortType, PageNumber and ItemsPerPage (nil SymbolName = all symbols) |

**Returns**

Returns `*pb.SymbolParamsManyData` with SymbolInfos (array of `*pb.SymbolParameters` containing Bid, Ask, Digits, Spread, VolumeMin, VolumeMax, VolumeStep, TradeContractSize, Point, margins, and other trading parameters for each symbol), SymbolsTotal, PageNumber and ItemsPerPage.

---

## TickValueWithSize

Retrieves tick value and size information for multiple symbols. Use this method to get tick values for profit/loss calculations and contract size data.

**Signature**
```go
func (a *MT5Account) TickValueWithSize(ctx context.Context, req *pb.TickValueWithSizeRequest) (*pb.TickValueWithSizeData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.TickValueWithSizeRequest | Request with array of SymbolNames |

**Returns**

Returns `*pb.TickValueWithSizeData` with SymbolTickSizeInfos (array of TickSizeSymbol) containing tick values, tick size, and contract size for each symbol.

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

Returns `*pb.PositionsTotalData` with TotalPositions count of open positions.

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
| req | *pb.OpenedOrdersRequest | Request with InputSortMode (sort type: by open time or by ticket ID, ascending or descending) |

**Returns**

Returns `*pb.OpenedOrdersData` with arrays OpenedOrders (pending orders, `*pb.OpenedOrderInfo`) and PositionInfos (open positions, `*pb.PositionInfo`) containing full details.

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

Returns `*pb.OpenedOrdersTicketsData` with arrays OpenedOrdersTickets and OpenedPositionTickets ([]int64).

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
| req | *pb.OrderCheckRequest | Request with MqlTradeRequest (`*pb.MrpcMqlTradeRequest`: Action, Symbol, Volume, Price, OrderType, StopLoss, TakeProfit and other order parameters) |

**Returns**

Returns `*pb.OrderCheckData` with MqlTradeCheckResult (`*pb.MrpcMqlTradeCheckResult`: ReturnedCode, BalanceAfterDeal, EquityAfterDeal, Profit, Margin, FreeMargin, MarginLevel, Comment).

---

## 📜 History

## OrderHistory

Retrieves historical orders and deals within a specified time range. Use this method to analyze past order activity with pagination support for large datasets.

**Signature**
```go
func (a *MT5Account) OrderHistory(ctx context.Context, req *pb.OrderHistoryRequest) (*pb.OrdersHistoryData, error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control |
| req | *pb.OrderHistoryRequest | Request with InputFrom, InputTo (Timestamps in broker server time, or InputFromUtc / InputToUtc), InputSortMode, and paging: PageNumber + ItemsPerPage, or PageToken + PageSize |

**Returns**

Returns `*pb.OrdersHistoryData` with HistoryData (rows holding HistoryOrder and/or HistoryDeal), ArrayTotal, PageNumber, ItemsPerPage and NextPageToken.

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
| req | *pb.PositionsHistoryRequest | Request with SortType, optional PositionOpenTimeFrom / PositionOpenTimeTo (Timestamps; the filter is on the position OPEN time; *Utc variants available), and paging: PageNumber + ItemsPerPage, or PageToken + PageSize |

**Returns**

Returns `*pb.PositionsHistoryData` with HistoryPositions (array of `*pb.PositionHistoryInfo`: PositionTicket, Symbol, OrderType, OpenTime, CloseTime, Volume, OpenPrice, ClosePrice, StopLoss, TakeProfit, Commission, Fee, Swap, Profit, ...), TotalCount and NextPageToken.

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

Returns `*pb.MarketBookAddData` with OpenedSuccessfully flag.

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

Returns `*pb.MarketBookGetData` with MqlBookInfos (array of `*pb.MrpcMqlBookInfo` entries containing Type, Price, Volume, and VolumeReal for each price level in the order book).

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

Returns `*pb.MarketBookReleaseData` with ClosedSuccessfully flag.

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

Returns `*pb.OrderSendData` with ReturnedCode (+ ReturnedStringCode, ReturnedCodeDescription), Deal, Order, Volume, Price, Bid, Ask, Comment, RequestId and RetCodeExternal.

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
| req | *pb.OrderModifyRequest | Request with Ticket, and optional (pointer) StopLoss, TakeProfit, Price and StopLimit (pending orders only), ExpirationTimeType / ExpirationTime |

**Returns**

Returns `*pb.OrderModifyData` with ReturnedCode (+ ReturnedStringCode, ReturnedCodeDescription), Deal, Order, Volume, Price, Bid, Ask, Comment, RequestId and RetCodeExternal.

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

Streams real-time tick data for one or more symbols. This method provides continuous price updates (Bid/Ask) as they arrive from the server. The stream automatically reconnects on connection loss.

**Signature**
```go
func (a *MT5Account) OnSymbolTick(ctx context.Context, req *pb.OnSymbolTickRequest) (<-chan *pb.OnSymbolTickData, <-chan error)
```

**Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Context for timeout and cancellation control (cancel to stop streaming) |
| req | *pb.OnSymbolTickRequest | Request with SymbolNames ([]string) |

**Returns**

Returns two channels: data channel receives `*pb.OnSymbolTickData` whose SymbolTick (`*pb.MrpcSubscriptionMqlTick`) has Symbol, Time, Bid, Ask, Last, Volume, TimeMsc, Flags and VolumeReal for each tick; error channel receives errors if stream fails (both channels closed on context cancellation).

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
| req | *pb.OnPositionProfitRequest | Request with TimerPeriodMilliseconds and IgnoreEmptyData |

**Returns**

Returns two channels: data channel receives `*pb.OnPositionProfitData` with NewPositions, UpdatedPositions and DeletedPositions (each entry: Ticket, Profit, PositionSymbol) plus AccountInfo; error channel receives errors if stream fails (both channels closed on context cancellation).

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
| req | *pb.OnPositionsAndPendingOrdersTicketsRequest | Request with TimerPeriodMilliseconds |

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

Returns two channels: data channel receives `*pb.OnTradeData` with EventData (new/disappeared/updated orders, history orders, deals and positions) and AccountInfo; error channel receives errors if stream fails (both channels closed on context cancellation).

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
