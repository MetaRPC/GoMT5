package mt5

/*
MT5Service - A wrapper over MT5Account, returns Go types instead of protobuf.

Architecture layers:
LOW → MT5Account (protobuf Request/Data, direct gRPC)
MID → MT5Service (Go types, removes Data wrappers)
HIGH → MT5Sugar (business logic, ready-made patterns)

Methods (37 items):

ACCOUNT:
- GetAccountSummary() - all account information
- GetAccountDouble() - double property (Balance, Equity)
- GetAccountInteger() - integer property (Login, Leverage)
- GetAccountString() - string property (Currency, Company)

SYMBOL:
- GetSymbolsTotal() - number of symbols
- SymbolExist() - existence check
- GetSymbolName() - symbol name by index
- SymbolSelect() - add/remove to Market Watch
- IsSymbolSynchronized() - synchronization check
- GetSymbolDouble() - double property (Bid, Ask)
- GetSymbolInteger() - integer property (Digits, Spread)
- GetSymbolString() - string property (Description)
- GetSymbolMarginRate() - margin rates
- GetSymbolTick() - latest tick
- GetSymbolSessionQuote() - quote session time
- GetSymbolSessionTrade() - trading session time
- GetSymbolParamsMany() - parameters of multiple symbols

POSITIONS & ORDERS:
- GetPositionsTotal() - number of open positions
- GetOpenedOrders() - all open orders/positions
- GetOpenedTickets() - ticket numbers only
- GetOrderHistory() - order history
- GetPositionsHistory() - closed positions history

MARKET DEPTH:
- SubscribeMarketDepth() - subscribe to DOM
- UnsubscribeMarketDepth() - unsubscribe from DOM
- GetMarketDepth() - current DOM snapshot

TRADING:
- PlaceOrder() - sending an order
- ModifyOrder() - modifying an order/position
- CloseOrder() - closing a position
- CheckOrder() - preliminary order check
- CalculateMargin() - calculating required margin
- CalculateProfit() - calculating potential profit

STREAMING:
- StreamTicks() - tick stream
- StreamTrades() - trade stream
- StreamPositionProfits() - position profit stream
- StreamTicketChanges() - ticket change stream
- StreamTradeTransactions() - trade transaction stream
*/

import (
	"context"
	"fmt"
	"time"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ══════════════════════════════════════════════════════════════════════════════
// SERVICE
// ══════════════════════════════════════════════════════════════════════════════

// MT5Service provides mid-level API wrapping MT5Account with Go native types.
// This layer unwraps protobuf and provides convenient request builders.
type MT5Service struct {
	account *MT5Account
}

// NewMT5Service creates a new MT5Service wrapping an MT5Account instance.
//
// Parameters:
//   - account: MT5Account instance (low-level gRPC client)
//
// Returns new MT5Service instance.
func NewMT5Service(account *MT5Account) *MT5Service {
	return &MT5Service{
		account: account,
	}
}

// GetAccount returns the underlying MT5Account for direct low-level access.
func (s *MT5Service) GetAccount() *MT5Account {
	return s.account
}

// ══════════════════════════════════════════════════════════════════════════════
// #region DATA TRANSFER OBJECTS (DTOs)
//
// These structures map protobuf responses to clean Go types.
// All DTO types are collected here for easy reference and maintenance.
// ══════════════════════════════════════════════════════════════════════════════

// AccountSummary holds all account information in one convenient struct.
//
// ADVANTAGE: Clean Go struct with native types instead of protobuf AccountSummaryData.
// All important account information in one place with time.Time instead of Timestamp.
type AccountSummary struct {
	Login                   int64                        // Account login number
	Balance                 float64                      // Account balance in deposit currency
	Equity                  float64                      // Account equity (Balance + Floating P&L)
	UserName                string                       // Client name
	Leverage                int64                        // Account leverage (e.g., 100 for 1:100)
	TradeMode               pb.MrpcEnumAccountTradeMode  // Account trade mode (demo/real/contest)
	CompanyName             string                       // Broker company name
	Currency                string                       // Deposit currency (USD, EUR, etc.)
	ServerTime              *time.Time                   // Server time (already converted from protobuf)
	UtcTimezoneShiftMinutes int64                        // UTC timezone shift in minutes
	Credit                  float64                      // Credit facility amount
}

// SymbolMarginRate holds margin rate information for a symbol.
//
// ADVANTAGE: Clean Go struct instead of protobuf SymbolInfoMarginRateData.
type SymbolMarginRate struct {
	InitialMarginRate     float64 // Initial margin rate
	MaintenanceMarginRate float64 // Maintenance margin rate
}

// SymbolTick holds current tick information for a symbol.
//
// ADVANTAGE: Clean Go struct with time.Time instead of protobuf SymbolInfoTickData.
// Time is already converted from Unix timestamp to time.Time.
type SymbolTick struct {
	Time       time.Time // Tick time (converted from Unix timestamp)
	Bid        float64   // Current Bid price
	Ask        float64   // Current Ask price
	Last       float64   // Last deal price
	Volume     uint64    // Tick volume
	TimeMS     int64     // Tick time in milliseconds
	Flags      uint32    // Tick flags
	VolumeReal float64   // Tick volume with decimal precision
}

// SessionTime holds trading session time range.
//
// ADVANTAGE: Clean Go struct with time.Time instead of protobuf Timestamp fields.
type SessionTime struct {
	From time.Time // Session start time (already converted from protobuf)
	To   time.Time // Session end time (already converted from protobuf)
}

// SymbolParams holds comprehensive symbol information.
//
// ADVANTAGE: Clean Go struct with all important symbol parameters.
// Much more convenient than making multiple calls to SymbolInfoDouble/Integer/String.
type SymbolParams struct {
	Name                 string  // Symbol name
	Bid                  float64 // Current Bid price
	Ask                  float64 // Current Ask price
	Last                 float64 // Last deal price
	Point                float64 // Point size (minimal price change)
	Digits               int32   // Number of decimal places
	Spread               int32   // Current spread in points
	VolumeMin            float64 // Minimum volume for trading
	VolumeMax            float64 // Maximum volume for trading
	VolumeStep           float64 // Volume step
	TradeTickSize        float64 // Trade tick size
	TradeTickValue       float64 // Trade tick value
	TradeContractSize    float64 // Contract size
	SwapLong             float64 // Swap for long positions
	SwapShort            float64 // Swap for short positions
	MarginInitial        float64 // Initial margin requirement
	MarginMaintenance    float64 // Maintenance margin requirement
}

// BookInfo holds a single Depth of Market (DOM) price level entry.
// Contains bid/ask price, volume, and type information.
type BookInfo struct {
	Type       pb.BookType // SELL (ask) or BUY (bid)
	Price      float64     // Price level
	Volume     int64       // Volume in lots (integer)
	VolumeReal float64     // Volume with decimal precision
}

// OrderResult holds the result of a trading operation.
//
// ADVANTAGE: Clean Go struct instead of protobuf OrderSendData/OrderModifyData.
// All fields in convenient Go types.
type OrderResult struct {
	ReturnedCode    uint32  // Operation return code (10009 = TRADE_RETCODE_DONE)
	Deal            uint64  // Deal ticket number (if executed)
	Order           uint64  // Order ticket number (if placed)
	Volume          float64 // Executed volume confirmed by broker
	Price           float64 // Execution price confirmed by broker
	Bid             float64 // Current Bid price
	Ask             float64 // Current Ask price
	Comment         string  // Broker comment or error description
	RequestID       uint32  // Request ID set by terminal
	RetCodeExternal int32   // Return code from external trading system
}

// OrderCheckResult holds the result of order pre-validation.
//
// ADVANTAGE: Clean Go struct instead of nested protobuf MqlTradeCheckResult.
// Shows what account state will be after order execution.
type OrderCheckResult struct {
	ReturnedCode uint32  // Validation result code
	Balance      float64 // Account balance after deal execution
	Equity       float64 // Account equity after deal execution
	Profit       float64 // Floating profit after deal
	Margin       float64 // Margin requirements for the order
	MarginFree   float64 // Free margin after deal
	MarginLevel  float64 // Margin level after deal
	Comment      string  // Error description (if validation failed)
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region ACCOUNT INFORMATION
// ══════════════════════════════════════════════════════════════════════════════

// GetAccountSummary retrieves all account information in one call.
//
// RECOMMENDED method for getting account data - returns everything you need.
//
// ADVANTAGE over MT5Account.AccountSummary:
//   - Returns clean AccountSummary struct with Go native types
//   - ServerTime is already *time.Time (no manual .AsTime() conversion)
//   - All fields have clear names (Balance instead of AccountBalance)
//   - No need to navigate protobuf AccountSummaryData
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//
// Returns:
//   - AccountSummary struct with all account properties
//   - Error if request failed
func (s *MT5Service) GetAccountSummary(ctx context.Context) (*AccountSummary, error) {
	req := &pb.AccountSummaryRequest{}

	data, err := s.account.AccountSummary(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetAccountSummary failed: %w", err)
	}

	var serverTime *time.Time
	if data.ServerTime != nil {
		t := data.ServerTime.AsTime()
		serverTime = &t
	}

	return &AccountSummary{
		Login:                   data.AccountLogin,
		Balance:                 data.AccountBalance,
		Equity:                  data.AccountEquity,
		UserName:                data.AccountUserName,
		Leverage:                data.AccountLeverage,
		TradeMode:               data.AccountTradeMode,
		CompanyName:             data.AccountCompanyName,
		Currency:                data.AccountCurrency,
		ServerTime:              serverTime,
		UtcTimezoneShiftMinutes: data.UtcTimezoneServerTimeShiftMinutes,
		Credit:                  data.AccountCredit,
	}, nil
}

// GetAccountDouble retrieves a double-type account property by ID.
//
// ADVANTAGE over MT5Account.AccountInfoDouble:
//   - Returns float64 directly instead of AccountInfoDoubleData
//   - No need to call GetRequestedValue() on response
//   - Cleaner function signature
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - propertyID: Property ID from AccountInfoDoublePropertyType enum
//     (e.g., ACCOUNT_BALANCE, ACCOUNT_EQUITY, ACCOUNT_MARGIN, etc.)
//
// Returns:
//   - Property value as float64
//   - Error if request failed
func (s *MT5Service) GetAccountDouble(ctx context.Context, propertyID pb.AccountInfoDoublePropertyType) (float64, error) {
	req := &pb.AccountInfoDoubleRequest{
		PropertyId: propertyID,
	}

	data, err := s.account.AccountInfoDouble(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("GetAccountDouble failed: %w", err)
	}

	return data.GetRequestedValue(), nil
}

// GetAccountInteger retrieves an integer-type account property by ID.
//
// ADVANTAGE over MT5Account.AccountInfoInteger:
//   - Returns int64 directly instead of AccountInfoIntegerData
//   - No need to call GetRequestedValue() on response
//   - Cleaner function signature
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - propertyID: Property ID from AccountInfoIntegerPropertyType enum
//     (e.g., ACCOUNT_LOGIN, ACCOUNT_LEVERAGE, ACCOUNT_TRADE_MODE, etc.)
//
// Returns:
//   - Property value as int64
//   - Error if request failed
func (s *MT5Service) GetAccountInteger(ctx context.Context, propertyID pb.AccountInfoIntegerPropertyType) (int64, error) {
	req := &pb.AccountInfoIntegerRequest{
		PropertyId: propertyID,
	}

	data, err := s.account.AccountInfoInteger(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("GetAccountInteger failed: %w", err)
	}

	return data.GetRequestedValue(), nil
}

// GetAccountString retrieves a string-type account property by ID.
//
// ADVANTAGE over MT5Account.AccountInfoString:
//   - Returns string directly instead of AccountInfoStringData
//   - No need to call GetRequestedValue() on response
//   - Cleaner function signature
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - propertyID: Property ID from AccountInfoStringPropertyType enum
//     (e.g., ACCOUNT_CURRENCY, ACCOUNT_COMPANY, ACCOUNT_NAME, etc.)
//
// Returns:
//   - Property value as string
//   - Error if request failed
func (s *MT5Service) GetAccountString(ctx context.Context, propertyID pb.AccountInfoStringPropertyType) (string, error) {
	req := &pb.AccountInfoStringRequest{
		PropertyId: propertyID,
	}

	data, err := s.account.AccountInfoString(ctx, req)
	if err != nil {
		return "", fmt.Errorf("GetAccountString failed: %w", err)
	}

	return data.GetRequestedValue(), nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region SYMBOL INFORMATION & OPERATIONS
// ══════════════════════════════════════════════════════════════════════════════

// GetSymbolsTotal returns the count of available symbols.
//
// ADVANTAGE over MT5Account.SymbolsTotal:
//   - Returns int32 directly instead of SymbolsTotalData
//   - Cleaner function signature with bool parameter
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - selectedOnly: If true, count only symbols in Market Watch; if false, count all symbols
//
// Returns:
//   - Number of symbols as int32
//   - Error if request failed
func (s *MT5Service) GetSymbolsTotal(ctx context.Context, selectedOnly bool) (int32, error) {
	req := &pb.SymbolsTotalRequest{
		Mode: selectedOnly,
	}

	data, err := s.account.SymbolsTotal(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("GetSymbolsTotal failed: %w", err)
	}

	return data.Total, nil
}

// SymbolExist checks if a symbol exists in the terminal.
// Returns (exists, isCustom, error). Use this before working with a symbol.
func (s *MT5Service) SymbolExist(ctx context.Context, symbol string) (bool, bool, error) {
	req := &pb.SymbolExistRequest{
		Name: symbol,
	}

	data, err := s.account.SymbolExist(ctx, req)
	if err != nil {
		return false, false, fmt.Errorf("SymbolExist failed: %w", err)
	}

	return data.Exists, data.IsCustom, nil
}

// GetSymbolName retrieves symbol name by index position.
// Use selectedOnly=true to get symbols from Market Watch only.
func (s *MT5Service) GetSymbolName(ctx context.Context, index int32, selectedOnly bool) (string, error) {
	req := &pb.SymbolNameRequest{
		Index:    index,
		Selected: selectedOnly,
	}

	data, err := s.account.SymbolName(ctx, req)
	if err != nil {
		return "", fmt.Errorf("GetSymbolName failed: %w", err)
	}

	return data.Name, nil
}

// SymbolSelect adds or removes a symbol from Market Watch window.
// Returns true if operation successful. Use select_=true to add, false to remove.
func (s *MT5Service) SymbolSelect(ctx context.Context, symbol string, select_ bool) (bool, error) {
	req := &pb.SymbolSelectRequest{
		Symbol: symbol,
		Select: select_,
	}

	data, err := s.account.SymbolSelect(ctx, req)
	if err != nil {
		return false, fmt.Errorf("SymbolSelect failed: %w", err)
	}

	return data.Success, nil
}

// IsSymbolSynchronized checks if symbol data is synchronized with the trade server.
// Returns true if symbol is fully synchronized and ready for trading operations.
func (s *MT5Service) IsSymbolSynchronized(ctx context.Context, symbol string) (bool, error) {
	req := &pb.SymbolIsSynchronizedRequest{
		Symbol: symbol,
	}

	data, err := s.account.SymbolIsSynchronized(ctx, req)
	if err != nil {
		return false, fmt.Errorf("IsSymbolSynchronized failed: %w", err)
	}

	return data.Synchronized, nil
}

// GetSymbolDouble retrieves a double-type symbol property (Bid, Ask, Point, etc.).
// Returns float64 value directly. For multiple properties use GetSymbolParamsMany instead.
func (s *MT5Service) GetSymbolDouble(ctx context.Context, symbol string, property pb.SymbolInfoDoubleProperty) (float64, error) {
	req := &pb.SymbolInfoDoubleRequest{
		Symbol: symbol,
		Type:   property,
	}

	data, err := s.account.SymbolInfoDouble(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("GetSymbolDouble failed: %w", err)
	}

	return data.Value, nil
}

// GetSymbolInteger retrieves an integer-type symbol property (Digits, Spread, etc.).
// Returns int64 value directly. For multiple properties use GetSymbolParamsMany instead.
func (s *MT5Service) GetSymbolInteger(ctx context.Context, symbol string, property pb.SymbolInfoIntegerProperty) (int64, error) {
	req := &pb.SymbolInfoIntegerRequest{
		Symbol: symbol,
		Type:   property,
	}

	data, err := s.account.SymbolInfoInteger(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("GetSymbolInteger failed: %w", err)
	}

	return data.Value, nil
}

// GetSymbolString retrieves a string-type symbol property (Description, Path, etc.).
// Returns string value directly.
func (s *MT5Service) GetSymbolString(ctx context.Context, symbol string, property pb.SymbolInfoStringProperty) (string, error) {
	req := &pb.SymbolInfoStringRequest{
		Symbol: symbol,
		Type:   property,
	}

	data, err := s.account.SymbolInfoString(ctx, req)
	if err != nil {
		return "", fmt.Errorf("GetSymbolString failed: %w", err)
	}

	return data.Value, nil
}

// GetSymbolMarginRate retrieves margin rates for a symbol and order type.
//
// ADVANTAGE over MT5Account.SymbolInfoMarginRate:
//   - Returns clean SymbolMarginRate struct instead of protobuf Data
//   - No need to navigate nested protobuf fields
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - symbol: Symbol name (e.g., "EURUSD")
//   - orderType: Order type from ENUM_ORDER_TYPE enum
//
// Returns:
//   - SymbolMarginRate struct with margin rates
//   - Error if request failed
func (s *MT5Service) GetSymbolMarginRate(ctx context.Context, symbol string, orderType pb.ENUM_ORDER_TYPE) (*SymbolMarginRate, error) {
	req := &pb.SymbolInfoMarginRateRequest{
		Symbol:    symbol,
		OrderType: orderType,
	}

	data, err := s.account.SymbolInfoMarginRate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetSymbolMarginRate failed: %w", err)
	}

	return &SymbolMarginRate{
		InitialMarginRate:     data.InitialMarginRate,
		MaintenanceMarginRate: data.MaintenanceMarginRate,
	}, nil
}

// GetSymbolTick retrieves the last tick for a symbol.
//
// ADVANTAGE over MT5Account.SymbolInfoTick:
//   - Returns SymbolTick struct with time.Time (no manual Unix conversion)
//   - Clean struct instead of protobuf SymbolInfoTickData
//   - All fields in convenient Go types
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - symbol: Symbol name (e.g., "EURUSD")
//
// Returns:
//   - SymbolTick struct with current tick data
//   - Error if request failed
func (s *MT5Service) GetSymbolTick(ctx context.Context, symbol string) (*SymbolTick, error) {
	req := &pb.SymbolInfoTickRequest{
		Symbol: symbol,
	}

	data, err := s.account.SymbolInfoTick(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetSymbolTick failed: %w", err)
	}

	return &SymbolTick{
		Time:       time.Unix(data.Time, 0),
		Bid:        data.Bid,
		Ask:        data.Ask,
		Last:       data.Last,
		Volume:     data.Volume,
		TimeMS:     data.TimeMsc,
		Flags:      data.Flags,
		VolumeReal: data.VolumeReal,
	}, nil
}

// GetSymbolSessionQuote retrieves quote session times for a symbol.
//
// ADVANTAGE over MT5Account.SymbolInfoSessionQuote:
//   - Returns SessionTime struct with time.Time (no manual .AsTime() calls)
//   - Clean struct instead of protobuf SymbolInfoSessionQuoteData
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - symbol: Symbol name (e.g., "EURUSD")
//   - dayOfWeek: Day of week from DayOfWeek enum
//   - sessionIndex: Session index (0-based)
//
// Returns:
//   - SessionTime struct with From/To times
//   - Error if request failed
func (s *MT5Service) GetSymbolSessionQuote(ctx context.Context, symbol string, dayOfWeek pb.DayOfWeek, sessionIndex uint32) (*SessionTime, error) {
	req := &pb.SymbolInfoSessionQuoteRequest{
		Symbol:       symbol,
		DayOfWeek:    dayOfWeek,
		SessionIndex: sessionIndex,
	}

	data, err := s.account.SymbolInfoSessionQuote(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetSymbolSessionQuote failed: %w", err)
	}

	return &SessionTime{
		From: data.From.AsTime(),
		To:   data.To.AsTime(),
	}, nil
}

// GetSymbolSessionTrade retrieves trade session times for a symbol.
// Similar to GetSymbolSessionQuote but for trading sessions instead of quote sessions.
func (s *MT5Service) GetSymbolSessionTrade(ctx context.Context, symbol string, dayOfWeek pb.DayOfWeek, sessionIndex uint32) (*SessionTime, error) {
	req := &pb.SymbolInfoSessionTradeRequest{
		Symbol:       symbol,
		DayOfWeek:    dayOfWeek,
		SessionIndex: sessionIndex,
	}

	data, err := s.account.SymbolInfoSessionTrade(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetSymbolSessionTrade failed: %w", err)
	}

	return &SessionTime{
		From: data.From.AsTime(),
		To:   data.To.AsTime(),
	}, nil
}

// GetSymbolParamsMany retrieves comprehensive parameters for multiple symbols.
//
// RECOMMENDED method for getting symbol information - returns all important params.
//
// ADVANTAGE over MT5Account.SymbolParamsMany:
//   - Returns []SymbolParams slice with clean Go structs
//   - All parameters in one place (no need for multiple SymbolInfo* calls)
//   - Supports pagination for large symbol lists
//   - Returns total count along with symbols
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - symbolName: Optional filter by symbol name (nil for all symbols)
//   - sortType: Optional sort type (nil for default sorting)
//   - pageNumber: Optional page number for pagination (nil for page 1)
//   - itemsPerPage: Optional items per page (nil for all items)
//
// Returns:
//   - Slice of SymbolParams structs
//   - Total number of symbols matching the filter
//   - Error if request failed
func (s *MT5Service) GetSymbolParamsMany(ctx context.Context, symbolName *string, sortType *pb.AH_SYMBOL_PARAMS_MANY_SORT_TYPE, pageNumber *int32, itemsPerPage *int32) ([]SymbolParams, int32, error) {
	req := &pb.SymbolParamsManyRequest{
		SymbolName:   symbolName,
		SortType:     sortType,
		PageNumber:   pageNumber,
		ItemsPerPage: itemsPerPage,
	}

	data, err := s.account.SymbolParamsMany(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("GetSymbolParamsMany failed: %w", err)
	}

	symbols := make([]SymbolParams, len(data.SymbolInfos))
	for i, info := range data.SymbolInfos {
		symbols[i] = SymbolParams{
			Name:              info.Name,
			Bid:               info.Bid,
			Ask:               info.Ask,
			Last:              info.Last,
			Point:             info.Point,
			Digits:            info.Digits,
			Spread:            info.Spread,
			VolumeMin:         info.VolumeMin,
			VolumeMax:         info.VolumeMax,
			VolumeStep:        info.VolumeStep,
			TradeTickSize:     info.TradeTickSize,
			TradeTickValue:    info.TradeTickValue,
			TradeContractSize: info.TradeContractSize,
			SwapLong:          info.SwapLong,
			SwapShort:         info.SwapShort,
			MarginInitial:     info.MarginInitial,
			MarginMaintenance: info.MarginMaintenance,
		}
	}

	return symbols, data.SymbolsTotal, nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region POSITIONS & ORDERS INFORMATION
// ══════════════════════════════════════════════════════════════════════════════

// GetPositionsTotal returns the count of open positions.
//
// ADVANTAGE over MT5Account.PositionsTotal:
//   - Returns int32 directly instead of PositionsTotalData
//   - Automatically extracts TotalPositions field
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//
// Returns:
//   - Number of open positions
//   - Error if request failed
func (s *MT5Service) GetPositionsTotal(ctx context.Context) (int32, error) {
	data, err := s.account.PositionsTotal(ctx)
	if err != nil {
		return 0, fmt.Errorf("GetPositionsTotal failed: %w", err)
	}
	return data.TotalPositions, nil
}

// GetOpenedOrders retrieves all open positions and pending orders.
//
// ADVANTAGE over MT5Account.OpenedOrders:
//   - Returns protobuf data directly (no unnecessary mapping)
//   - Same as C# MT5Service approach
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - sortMode: Sort mode from BMT5_ENUM_OPENED_ORDER_SORT_TYPE enum
//
// Returns:
//   - OpenedOrdersData containing positions and pending orders
//   - Error if request failed
func (s *MT5Service) GetOpenedOrders(ctx context.Context, sortMode pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE) (*pb.OpenedOrdersData, error) {
	req := &pb.OpenedOrdersRequest{
		InputSortMode: sortMode,
	}
	data, err := s.account.OpenedOrders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetOpenedOrders failed: %w", err)
	}
	return data, nil
}

// GetOpenedTickets retrieves ticket numbers of open positions and pending orders.
// Lightweight alternative to GetOpenedOrders - returns only ticket numbers, not full details.
// Returns (positionTickets, orderTickets, error).
func (s *MT5Service) GetOpenedTickets(ctx context.Context) ([]int64, []int64, error) {
	req := &pb.OpenedOrdersTicketsRequest{}

	data, err := s.account.OpenedOrdersTickets(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("GetOpenedTickets failed: %w", err)
	}

	return data.OpenedPositionTickets, data.OpenedOrdersTickets, nil
}

// GetOrderHistory retrieves historical orders and deals for a time period with pagination.
//
// ADVANTAGE over MT5Account.OrderHistory:
//   - Returns protobuf data directly (no unnecessary mapping)
//   - Same as C# MT5Service approach
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - from: Start time of history range
//   - to: End time of history range
//   - sortMode: Sort mode from BMT5_ENUM_ORDER_HISTORY_SORT_TYPE enum
//   - pageNumber: Page number for pagination (1-based)
//   - itemsPerPage: Number of items per page
//
// Returns:
//   - OrdersHistoryData containing orders and deals
//   - Error if request failed
func (s *MT5Service) GetOrderHistory(ctx context.Context, from time.Time, to time.Time, sortMode pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE, pageNumber int32, itemsPerPage int32) (*pb.OrdersHistoryData, error) {
	req := &pb.OrderHistoryRequest{
		InputFrom:     timestamppb.New(from),
		InputTo:       timestamppb.New(to),
		InputSortMode: sortMode,
		PageNumber:    pageNumber,
		ItemsPerPage:  itemsPerPage,
	}
	data, err := s.account.OrderHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetOrderHistory failed: %w", err)
	}
	return data, nil
}

// GetPositionsHistory retrieves closed positions with aggregated P&L data.
//
// ADVANTAGE over MT5Account.PositionsHistory:
//   - Returns protobuf data directly (no unnecessary mapping)
//   - Same as C# MT5Service approach
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - sortType: Sort type from AH_ENUM_POSITIONS_HISTORY_SORT_TYPE enum
//   - from: Optional start time for position open time filter (nil for no filter)
//   - to: Optional end time for position open time filter (nil for no filter)
//   - pageNumber: Optional page number for pagination (nil for page 1)
//   - itemsPerPage: Optional items per page (nil for all items)
//
// Returns:
//   - PositionsHistoryData containing closed positions
//   - Error if request failed
func (s *MT5Service) GetPositionsHistory(ctx context.Context, sortType pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE, from *time.Time, to *time.Time, pageNumber *int32, itemsPerPage *int32) (*pb.PositionsHistoryData, error) {
	req := &pb.PositionsHistoryRequest{
		SortType:               sortType,
		PositionOpenTimeFrom:   nil,
		PositionOpenTimeTo:     nil,
		PageNumber:             pageNumber,
		ItemsPerPage:           itemsPerPage,
	}

	if from != nil {
		req.PositionOpenTimeFrom = timestamppb.New(*from)
	}
	if to != nil {
		req.PositionOpenTimeTo = timestamppb.New(*to)
	}

	data, err := s.account.PositionsHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetPositionsHistory failed: %w", err)
	}
	return data, nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region MARKET DEPTH / DOM
// ══════════════════════════════════════════════════════════════════════════════

// SubscribeMarketDepth subscribes to Depth of Market (DOM) updates for a symbol.
// Must be called before GetMarketDepth. Returns true if subscription successful.
func (s *MT5Service) SubscribeMarketDepth(ctx context.Context, symbol string) (bool, error) {
	req := &pb.MarketBookAddRequest{
		Symbol: symbol,
	}

	data, err := s.account.MarketBookAdd(ctx, req)
	if err != nil {
		return false, fmt.Errorf("SubscribeMarketDepth failed: %w", err)
	}

	return data.OpenedSuccessfully, nil
}

// UnsubscribeMarketDepth unsubscribes from Depth of Market updates.
// Call this to stop receiving DOM updates and free resources.
func (s *MT5Service) UnsubscribeMarketDepth(ctx context.Context, symbol string) (bool, error) {
	req := &pb.MarketBookReleaseRequest{
		Symbol: symbol,
	}

	data, err := s.account.MarketBookRelease(ctx, req)
	if err != nil {
		return false, fmt.Errorf("UnsubscribeMarketDepth failed: %w", err)
	}

	return data.ClosedSuccessfully, nil
}

// GetMarketDepth retrieves current Depth of Market (DOM) snapshot for a symbol.
// Requires prior SubscribeMarketDepth call. Returns slice of BookInfo entries.
func (s *MT5Service) GetMarketDepth(ctx context.Context, symbol string) ([]BookInfo, error) {
	req := &pb.MarketBookGetRequest{
		Symbol: symbol,
	}

	data, err := s.account.MarketBookGet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetMarketDepth failed: %w", err)
	}

	books := make([]BookInfo, len(data.MqlBookInfos))
	for i, b := range data.MqlBookInfos {
		books[i] = BookInfo{
			Type:       b.Type,
			Price:      b.Price,
			Volume:     b.Volume,
			VolumeReal: b.VolumeReal,
		}
	}

	return books, nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region TRADING OPERATIONS
// ══════════════════════════════════════════════════════════════════════════════

// PlaceOrder sends a market or pending order to MT5 terminal.
//
// ADVANTAGE over MT5Account.OrderSend:
//   - Returns clean OrderResult struct instead of OrderSendData
//   - All fields in convenient Go types
//   - Easy to check if order was successful (result.ReturnedCode == 10009)
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - req: Order send request (protobuf OrderSendRequest)
//
// Returns:
//   - OrderResult struct with execution details
//   - Error if request failed
func (s *MT5Service) PlaceOrder(ctx context.Context, req *pb.OrderSendRequest) (*OrderResult, error) {
	data, err := s.account.OrderSend(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("PlaceOrder failed: %w", err)
	}

	return &OrderResult{
		ReturnedCode:    data.ReturnedCode,
		Deal:            data.Deal,
		Order:           data.Order,
		Volume:          data.Volume,
		Price:           data.Price,
		Bid:             data.Bid,
		Ask:             data.Ask,
		Comment:         data.Comment,
		RequestID:       data.RequestId,
		RetCodeExternal: data.RetCodeExternal,
	}, nil
}

// ModifyOrder modifies an existing order or position (change SL/TP/price).
// Returns OrderResult with modification details. Check ReturnedCode for success (10009).
func (s *MT5Service) ModifyOrder(ctx context.Context, req *pb.OrderModifyRequest) (*OrderResult, error) {
	data, err := s.account.OrderModify(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ModifyOrder failed: %w", err)
	}

	return &OrderResult{
		ReturnedCode:    data.ReturnedCode,
		Deal:            data.Deal,
		Order:           data.Order,
		Volume:          data.Volume,
		Price:           data.Price,
		Bid:             data.Bid,
		Ask:             data.Ask,
		Comment:         data.Comment,
		RequestID:       data.RequestId,
		RetCodeExternal: data.RetCodeExternal,
	}, nil
}

// CloseOrder closes a position or deletes a pending order.
// Returns operation return code (10009 = success). Simpler than PlaceOrder for closing.
func (s *MT5Service) CloseOrder(ctx context.Context, req *pb.OrderCloseRequest) (uint32, error) {
	data, err := s.account.OrderClose(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("CloseOrder failed: %w", err)
	}

	return data.ReturnedCode, nil
}

// CheckOrder validates an order before sending it to the broker.
//
// ADVANTAGE over MT5Account.OrderCheck:
//   - Returns clean OrderCheckResult struct
//   - Automatically extracts MqlTradeCheckResult fields
//   - No need to navigate nested protobuf structures
//
// Use this to validate orders before PlaceOrder to avoid rejections.
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - req: Order check request (protobuf OrderCheckRequest)
//
// Returns:
//   - OrderCheckResult with validation results
//   - Error if request failed
func (s *MT5Service) CheckOrder(ctx context.Context, req *pb.OrderCheckRequest) (*OrderCheckResult, error) {
	data, err := s.account.OrderCheck(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CheckOrder failed: %w", err)
	}

	result := data.MqlTradeCheckResult
	return &OrderCheckResult{
		ReturnedCode: result.ReturnedCode,
		Balance:      result.BalanceAfterDeal,
		Equity:       result.EquityAfterDeal,
		Profit:       result.Profit,
		Margin:       result.Margin,
		MarginFree:   result.FreeMargin,
		MarginLevel:  result.MarginLevel,
		Comment:      result.Comment,
	}, nil
}

// CalculateMargin calculates required margin for a potential order.
// Use this before placing orders to check if you have enough free margin.
func (s *MT5Service) CalculateMargin(ctx context.Context, req *pb.OrderCalcMarginRequest) (float64, error) {
	data, err := s.account.OrderCalcMargin(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("CalculateMargin failed: %w", err)
	}

	return data.Margin, nil
}

// CalculateProfit calculates potential profit for a hypothetical order.
// Useful for profit/risk calculations before placing actual orders.
func (s *MT5Service) CalculateProfit(ctx context.Context, req *pb.OrderCalcProfitRequest) (float64, error) {
	data, err := s.account.OrderCalcProfit(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("CalculateProfit failed: %w", err)
	}

	return data.Profit, nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region STREAMING METHODS
// ══════════════════════════════════════════════════════════════════════════════

// StreamTicks streams real-time ticks for specified symbols.
//
// ADVANTAGE over MT5Account.OnSymbolTick:
//   - Returns channel of *SymbolTick structs (clean Go types)
//   - Automatically converts protobuf OnSymbolTickData to SymbolTick
//   - Time fields are time.Time (no manual .AsTime() conversions)
//   - Cleaner API with separate tick and error channels
//
// The returned channels will be closed when streaming stops.
// Always read from both channels in a select statement.
//
// Parameters:
//   - ctx: Context for cancellation (closing ctx stops the stream)
//   - symbols: Slice of symbol names to stream (e.g., []string{"EURUSD", "GBPUSD"})
//
// Returns:
//   - Read-only channel of *SymbolTick structs
//   - Read-only channel of errors
func (s *MT5Service) StreamTicks(ctx context.Context, symbols []string) (<-chan *SymbolTick, <-chan error) {
	req := &pb.OnSymbolTickRequest{
		SymbolNames: symbols,
	}

	dataCh, errCh := s.account.OnSymbolTick(ctx, req)

	tickCh := make(chan *SymbolTick)
	outErrCh := make(chan error, 1)

	go func() {
		defer close(tickCh)
		defer close(outErrCh)

		for {
			select {
			case data, ok := <-dataCh:
				if !ok {
					return
				}
				tick := data.SymbolTick
				tickCh <- &SymbolTick{
					Time:       tick.Time.AsTime(),
					Bid:        tick.Bid,
					Ask:        tick.Ask,
					Last:       tick.Last,
					Volume:     tick.Volume,
					TimeMS:     tick.TimeMsc,
					Flags:      tick.Flags,
					VolumeReal: tick.VolumeReal,
				}
			case err, ok := <-errCh:
				if !ok {
					return
				}
				outErrCh <- err
				return
			case <-ctx.Done():
				outErrCh <- ctx.Err()
				return
			}
		}
	}()

	return tickCh, outErrCh
}

// StreamTradeUpdates streams trade events (new/disappeared orders and positions, history updates).
//
// This method provides real-time notifications about:
//   - New orders and positions opened
//   - Orders and positions closed/cancelled
//   - History orders and deals
//
// The returned channels will be closed when streaming stops.
// Always read from both channels in a select statement.
//
// Parameters:
//   - ctx: Context for cancellation (closing ctx stops the stream)
//
// Returns:
//   - Read-only channel of *pb.OnTradeData (protobuf)
//   - Read-only channel of errors
func (s *MT5Service) StreamTradeUpdates(ctx context.Context) (<-chan *pb.OnTradeData, <-chan error) {
	req := &pb.OnTradeRequest{}
	return s.account.OnTrade(ctx, req)
}

// StreamPositionProfits streams real-time profit/loss updates for open positions.
//
// This method provides real-time P&L updates as prices change:
//   - Position profit updates as market moves
//   - New positions opened
//   - Positions modified (SL/TP changes)
//   - Positions closed
//
// The returned channels will be closed when streaming stops.
// Always read from both channels in a select statement.
//
// Parameters:
//   - ctx: Context for cancellation (closing ctx stops the stream)
//
// Returns:
//   - Read-only channel of *pb.OnPositionProfitData (protobuf)
//   - Read-only channel of errors
func (s *MT5Service) StreamPositionProfits(ctx context.Context) (<-chan *pb.OnPositionProfitData, <-chan error) {
	req := &pb.OnPositionProfitRequest{}
	return s.account.OnPositionProfit(ctx, req)
}

// StreamOpenedTickets streams updates to the list of open position and pending order tickets.
//
// This method provides lightweight notifications about ticket changes:
//   - List of currently open position tickets
//   - List of currently pending order tickets
//   - Updates when tickets are added or removed
//
// This is a lightweight alternative to StreamTradeUpdates when you only need ticket IDs.
//
// The returned channels will be closed when streaming stops.
// Always read from both channels in a select statement.
//
// Parameters:
//   - ctx: Context for cancellation (closing ctx stops the stream)
//
// Returns:
//   - Read-only channel of *pb.OnPositionsAndPendingOrdersTicketsData (protobuf)
//   - Read-only channel of errors
func (s *MT5Service) StreamOpenedTickets(ctx context.Context) (<-chan *pb.OnPositionsAndPendingOrdersTicketsData, <-chan error) {
	req := &pb.OnPositionsAndPendingOrdersTicketsRequest{}
	return s.account.OnPositionsAndPendingOrdersTickets(ctx, req)
}

// StreamTransactions streams all trade transaction events (most detailed streaming method).
//
// This method provides the most comprehensive trade event stream:
//   - Order placement, modification, deletion
//   - Deal execution
//   - Position opening, modification, closing
//   - All intermediate states and changes
//
// This is the most powerful streaming method, providing detailed transaction information.
// Use StreamTradeUpdates for simpler trade monitoring or StreamPositionProfits for P&L tracking.
//
// The returned channels will be closed when streaming stops.
// Always read from both channels in a select statement.
//
// Parameters:
//   - ctx: Context for cancellation (closing ctx stops the stream)
//
// Returns:
//   - Read-only channel of *pb.OnTradeTransactionData (protobuf)
//   - Read-only channel of errors
func (s *MT5Service) StreamTransactions(ctx context.Context) (<-chan *pb.OnTradeTransactionData, <-chan error) {
	req := &pb.OnTradeTransactionRequest{}
	return s.account.OnTradeTransaction(ctx, req)
}
// #endregion
