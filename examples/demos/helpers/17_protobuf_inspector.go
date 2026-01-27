/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/demos/helpers/protobuf_inspector.go - INTERACTIVE PROTOBUF TYPES INSPECTOR
 PURPOSE:
   Interactive developer utility to explore MT5 protobuf types, fields,
   enums, and data structures from the MT5 gRPC API.

 🎯 WHAT THIS TOOL DOES:
   • Interactive search for types, fields, and enums
   • Real-time inspection of protobuf message structures
   • Field-level discovery (find which types contain specific fields)
   • Enum value exploration (see all possible values)
   • Type browsing (list all available types)

 📖 HOW TO USE:

   1. START THE INSPECTOR:
      cd examples/demos
      go run main.go inspect

   2. AVAILABLE COMMANDS:

      ┌─────────────────────────────────────────────────────────────────────┐
      │ COMMAND          │ DESCRIPTION                                      │
      ├─────────────────────────────────────────────────────────────────────┤
      │ list             │ List all available protobuf types                │
      │ ls               │ (alias for list)                                 │
      ├─────────────────────────────────────────────────────────────────────┤
      │ <TypeName>       │ Inspect specific type (e.g., "PositionInfo")     │
      │                  │ Shows all fields and their types                 │
      ├─────────────────────────────────────────────────────────────────────┤
      │ search <text>    │ Search for types containing text                 │
      │ find <text>      │ (e.g., "search Order" finds all Order* types)    │
      ├─────────────────────────────────────────────────────────────────────┤
      │ field <name>     │ Find all types containing a specific field       │
      │                  │ (e.g., "field Balance" → AccountInfo)            │
      ├─────────────────────────────────────────────────────────────────────┤
      │ enum <name>      │ Show all values of an enum                       │
      │                  │ (e.g., "enum BMT5_ENUM_ORDER_TYPE")              │
      ├─────────────────────────────────────────────────────────────────────┤
      │ help             │ Show this help message                           │
      │ ?                │ (alias for help)                                 │
      ├─────────────────────────────────────────────────────────────────────┤
      │ exit             │ Exit the inspector                               │
      │ quit             │ (alias for exit)                                 │
      └─────────────────────────────────────────────────────────────────────┘

 💡 PRACTICAL EXAMPLES:

   Example 1: Find out what fields PositionInfo has
   > PositionInfo
   Output:
     Ticket: uint64
     Type: BMT5_ENUM_POSITION_TYPE
     Symbol: string
     ...

   Example 2: Find which type has the "Ticket" field
   > field Ticket
   Output:
     Found in: PositionInfo, OrderInfo, DealInfo, ...

   Example 3: See all ORDER_TYPE values
   > enum BMT5_ENUM_ORDER_TYPE
   Output:
     BMT5_ORDER_TYPE_BUY = 0
     BMT5_ORDER_TYPE_SELL = 1
     BMT5_ORDER_TYPE_BUY_LIMIT = 2
     ...

   Example 4: Find all types related to "Position"
   > search Position
   Output:
     PositionInfo
     PositionsGetRequest
     PositionsGetReply
     ...

   Example 5: List all available types
   > list
   Output:
     [Shows all protobuf types]

 🔍 COMMON USE CASES:

   USE CASE 1: "I'm getting 'field not found' error"
   → Use: field <fieldname>
   → Example: field Equity
   → Result: Shows you the correct field name and which type has it

   USE CASE 2: "What fields does X have?"
   → Use: <TypeName>
   → Example: PositionInfo
   → Result: Lists all fields (Ticket, Type, Symbol, etc.)

   USE CASE 3: "What are valid enum values?"
   → Use: enum <EnumName>
   → Example: enum BMT5_ENUM_ORDER_TYPE
   → Result: Shows all values with their numeric codes

   USE CASE 4: "I need to find types related to positions"
   → Use: search Position
   → Result: Lists all types with "Position" in the name

   USE CASE 5: "I want to browse what's available"
   → Use: list
   → Result: Shows all available types to explore

 📊 STATISTICS:
   • Total Types: 267 (all MT5 gRPC protobuf message types)
   • Total Enums: 67 (with 1400+ enum values)
   • Coverage: 100% of MT5 gRPC API types
   • Note: Type count is dynamic - actual count may vary with API updates

 🔑 IMPORTANT ENUMS (frequently used):
   • BMT5_ENUM_ORDER_TYPE          - Order types (BUY, SELL, LIMIT, STOP, CLOSE_BY)
   • BMT5_ENUM_ORDER_TYPE_FILLING  - Fill policies (FOK, IOC, Return, BOC)
   • BMT5_ENUM_ORDER_TYPE_TIME     - Time in force (GTC, Day, Specified)
   • BMT5_ENUM_DEAL_REASON         - Deal reasons (SL, TP, Expert, Client, Mobile)
   • BMT5_ENUM_DEAL_ENTRY_TYPE     - Deal entry (In, Out, InOut, Out By)
   • BMT5_ENUM_POSITION_TYPE       - Position direction (BUY, SELL)
   • BMT5_ENUM_POSITION_REASON     - Why position opened (Client, Expert, Mobile)
   • MqlErrorCode                  - MQL error codes (211 values!)
   • MqlErrorTradeCode             - Trade operation errors (REQUOTE, REJECT, etc.)
   • SymbolInfoDoubleProperty      - Symbol price properties (BID, ASK, POINT, SWAP)
   • SymbolInfoIntegerProperty     - Symbol integer properties (DIGITS, SPREAD, etc.)
   • MRPC_ENUM_TRADE_REQUEST_ACTIONS - Trade actions (DEAL, PENDING, SLTP, MODIFY)

 ✨ FEATURES:
   • Case-insensitive search       - "search Order" = "search order"
   • Partial field matching         - "field profit" finds Profit, TakeProfit
   • Type categorization            - [Request], [Reply], [Type], [Info]
   • Array indicators               - 📚 icon for repeated/array fields
   • Protobuf field numbers         - Shows field #N for each field
   • Smart error messages           - Suggests alternatives when type not found

 📚 FULL DOCUMENTATION:
   For complete API documentation with usage examples, method specs,
   request/response structures, and enum definitions, see:

   docs/MT5Account/MT5Account.Master.Overview.md

   This master overview provides:
   • All method specifications by category
   • Request/Response message structures
   • Enum definitions with descriptions
   • Usage patterns and code examples
   • Jump links to detailed method docs

 USAGE:
   cd examples/demos
   go run main.go inspect

══════════════════════════════════════════════════════════════════════════════*/

package helpers

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	pb "github.com/MetaRPC/GoMT5/package"
)

// ProtobufInspector provides interactive exploration of MT5 protobuf types
type ProtobufInspector struct {
	types map[string]reflect.Type
	enums map[string][]EnumValue
}

// EnumValue represents an enum constant with its name and value
type EnumValue struct {
	Name  string
	Value int32
}

// NewProtobufInspector creates a new protobuf inspector
func NewProtobufInspector() *ProtobufInspector {
	inspector := &ProtobufInspector{
		types: make(map[string]reflect.Type),
		enums: make(map[string][]EnumValue),
	}
	inspector.discoverTypes()
	return inspector
}

// discoverTypes uses reflection to find all protobuf types
func (pi *ProtobufInspector) discoverTypes() {
	// Use reflection to discover ALL types from pb package
	// We'll use a sample instance to get the package, then scan all exported types

	// Register ALL available protobuf types from pb package (267 types)
	sampleTypes := []interface{}{
		&pb.AccountInfoDoubleData{}, &pb.AccountInfoDoubleReply{}, &pb.AccountInfoDoubleReply_Data{}, &pb.AccountInfoDoubleReply_Error{},
		&pb.AccountInfoDoubleRequest{}, &pb.AccountInfoIntegerData{}, &pb.AccountInfoIntegerReply{}, &pb.AccountInfoIntegerReply_Data{},
		&pb.AccountInfoIntegerReply_Error{}, &pb.AccountInfoIntegerRequest{}, &pb.AccountInfoStringData{}, &pb.AccountInfoStringReply{},
		&pb.AccountInfoStringReply_Data{}, &pb.AccountInfoStringReply_Error{}, &pb.AccountInfoStringRequest{}, &pb.AccountSummaryData{},
		&pb.AccountSummaryReply{}, &pb.AccountSummaryReply_Data{}, &pb.AccountSummaryReply_Error{}, &pb.AccountSummaryRequest{},
		&pb.CheckConnectData{}, &pb.CheckConnectReply{}, &pb.CheckConnectReply_Data{}, &pb.CheckConnectReply_Error{},
		&pb.CheckConnectRequest{}, &pb.CloseChartForSymbolReply{}, &pb.CloseChartForSymbolRequest{}, &pb.ConnectData{},
		&pb.ConnectExReply{}, &pb.ConnectExReply_Data{}, &pb.ConnectExReply_Error{}, &pb.ConnectExRequest{},
		&pb.ConnectProxyData{}, &pb.ConnectProxyReply{}, &pb.ConnectProxyReply_Data{}, &pb.ConnectProxyReply_Error{},
		&pb.ConnectProxyRequest{}, &pb.ConnectReply{}, &pb.ConnectReply_Data{}, &pb.ConnectReply_Error{},
		&pb.ConnectRequest{}, &pb.DealHistoryData{}, &pb.DisconnectData{}, &pb.DisconnectReply{},
		&pb.DisconnectReply_Data{}, &pb.DisconnectReply_Error{}, &pb.DisconnectRequest{}, &pb.EaParam{},
		&pb.Error{}, &pb.ErrorProperty{}, &pb.ExpertAdviser{}, &pb.GetEaParamsData{},
		&pb.GetEaParamsReply{}, &pb.GetEaParamsReply_Data{}, &pb.GetEaParamsReply_Error{}, &pb.GetEaParamsRequest{},
		&pb.HealthCheckReply{}, &pb.HealthCheckRequest{}, &pb.HistoryData{}, &pb.LogFileInfo{},
		&pb.MarketBookAddData{}, &pb.MarketBookAddReply{}, &pb.MarketBookAddReply_Data{}, &pb.MarketBookAddReply_Error{},
		&pb.MarketBookAddRequest{}, &pb.MarketBookGetData{}, &pb.MarketBookGetReply{}, &pb.MarketBookGetReply_Data{},
		&pb.MarketBookGetReply_Error{}, &pb.MarketBookGetRequest{}, &pb.MarketBookReleaseData{}, &pb.MarketBookReleaseReply{},
		&pb.MarketBookReleaseReply_Data{}, &pb.MarketBookReleaseReply_Error{}, &pb.MarketBookReleaseRequest{}, &pb.MqlTradeRequest{},
		&pb.MqlTradeResult{}, &pb.MqlTradeTransaction{}, &pb.MrpcMqlBookInfo{}, &pb.MrpcMqlTick{},
		&pb.MrpcMqlTradeCheckResult{}, &pb.MrpcMqlTradeRequest{}, &pb.MrpcSubscriptionMqlTick{}, &pb.OnEventAccountInfo{},
		&pb.OnPositionProfitData{}, &pb.OnPositionProfitPositionInfo{}, &pb.OnPositionProfitReply{}, &pb.OnPositionProfitReply_Data{},
		&pb.OnPositionProfitReply_Error{}, &pb.OnPositionProfitRequest{}, &pb.OnPositionsAndPendingOrdersTicketsData{}, &pb.OnPositionsAndPendingOrdersTicketsReply{},
		&pb.OnPositionsAndPendingOrdersTicketsReply_Data{}, &pb.OnPositionsAndPendingOrdersTicketsReply_Error{}, &pb.OnPositionsAndPendingOrdersTicketsRequest{}, &pb.OnSymbolTickData{},
		&pb.OnSymbolTickReply{}, &pb.OnSymbolTickReply_Data{}, &pb.OnSymbolTickReply_Error{}, &pb.OnSymbolTickRequest{},
		&pb.OnTadeEventData{}, &pb.OnTradeData{}, &pb.OnTradeHistoryDealInfo{}, &pb.OnTradeHistoryDealUpdate{},
		&pb.OnTradeHistoryOrderInfo{}, &pb.OnTradeHistoryOrderUpdate{}, &pb.OnTradeOrderInfo{}, &pb.OnTradeOrderStateChange{},
		&pb.OnTradePositionInfo{}, &pb.OnTradePositionUpdate{}, &pb.OnTradeReply{}, &pb.OnTradeReply_Data{},
		&pb.OnTradeReply_Error{}, &pb.OnTradeRequest{}, &pb.OnTradeTransactionData{}, &pb.OnTradeTransactionReply{},
		&pb.OnTradeTransactionReply_Data{}, &pb.OnTradeTransactionReply_Error{}, &pb.OnTradeTransactionRequest{}, &pb.OpenChartForSymbolReply{},
		&pb.OpenChartForSymbolRequest{}, &pb.OpenChartWithEaParameter{}, &pb.OpenChartWithEaReply{}, &pb.OpenChartWithEaRequest{},
		&pb.OpenedOrderInfo{}, &pb.OpenedOrdersData{}, &pb.OpenedOrdersReply{}, &pb.OpenedOrdersReply_Data{},
		&pb.OpenedOrdersReply_Error{}, &pb.OpenedOrdersRequest{}, &pb.OpenedOrdersTicketsData{}, &pb.OpenedOrdersTicketsReply{},
		&pb.OpenedOrdersTicketsReply_Data{}, &pb.OpenedOrdersTicketsReply_Error{}, &pb.OpenedOrdersTicketsRequest{}, &pb.OpenTerminalChartWithEaData{},
		&pb.OpenTerminalChartWithEaParameter{}, &pb.OpenTerminalChartWithEaReply{}, &pb.OpenTerminalChartWithEaReply_Data{}, &pb.OpenTerminalChartWithEaReply_Error{},
		&pb.OpenTerminalChartWithEaRequest{}, &pb.OrderCalcMarginData{}, &pb.OrderCalcMarginReply{}, &pb.OrderCalcMarginReply_Data{},
		&pb.OrderCalcMarginReply_Error{}, &pb.OrderCalcMarginRequest{}, &pb.OrderCalcProfitData{}, &pb.OrderCalcProfitReply{},
		&pb.OrderCalcProfitReply_Data{}, &pb.OrderCalcProfitReply_Error{}, &pb.OrderCalcProfitRequest{}, &pb.OrderCheckData{},
		&pb.OrderCheckReply{}, &pb.OrderCheckReply_Data{}, &pb.OrderCheckReply_Error{}, &pb.OrderCheckRequest{},
		&pb.OrderCloseData{}, &pb.OrderCloseReply{}, &pb.OrderCloseReply_Data{}, &pb.OrderCloseReply_Error{},
		&pb.OrderCloseRequest{}, &pb.OrderHistoryData{}, &pb.OrderHistoryReply{}, &pb.OrderHistoryReply_Data{},
		&pb.OrderHistoryReply_Error{}, &pb.OrderHistoryRequest{}, &pb.OrderModifyData{}, &pb.OrderModifyReply{},
		&pb.OrderModifyReply_Data{}, &pb.OrderModifyReply_Error{}, &pb.OrderModifyRequest{}, &pb.OrderSendData{},
		&pb.OrderSendReply{}, &pb.OrderSendReply_Data{}, &pb.OrderSendReply_Error{}, &pb.OrderSendRequest{},
		&pb.OrdersHistoryData{}, &pb.PositionHistoryInfo{}, &pb.PositionInfo{}, &pb.PositionsHistoryData{},
		&pb.PositionsHistoryReply{}, &pb.PositionsHistoryReply_Data{}, &pb.PositionsHistoryReply_Error{}, &pb.PositionsHistoryRequest{},
		&pb.PositionsTotalData{}, &pb.PositionsTotalReply{}, &pb.PositionsTotalReply_Data{}, &pb.PositionsTotalReply_Error{},
		&pb.ReconnectData{}, &pb.ReconnectReply{}, &pb.ReconnectReply_Data{}, &pb.ReconnectReply_Error{},
		&pb.ReconnectRequest{}, &pb.SaveChartTemplateReply{}, &pb.SaveChartTemplateRequest{}, &pb.StopListeningReply{},
		&pb.StopListeningRequest{}, &pb.SymbolExistData{}, &pb.SymbolExistReply{}, &pb.SymbolExistReply_Data{},
		&pb.SymbolExistReply_Error{}, &pb.SymbolExistRequest{}, &pb.SymbolInfoDoubleData{}, &pb.SymbolInfoDoubleReply{},
		&pb.SymbolInfoDoubleReply_Data{}, &pb.SymbolInfoDoubleReply_Error{}, &pb.SymbolInfoDoubleRequest{}, &pb.SymbolInfoIntegerData{},
		&pb.SymbolInfoIntegerReply{}, &pb.SymbolInfoIntegerReply_Data{}, &pb.SymbolInfoIntegerReply_Error{}, &pb.SymbolInfoIntegerRequest{},
		&pb.SymbolInfoMarginRateData{}, &pb.SymbolInfoMarginRateReply{}, &pb.SymbolInfoMarginRateReply_Data{}, &pb.SymbolInfoMarginRateReply_Error{},
		&pb.SymbolInfoMarginRateRequest{}, &pb.SymbolInfoSessionQuoteData{}, &pb.SymbolInfoSessionQuoteReply{}, &pb.SymbolInfoSessionQuoteReply_Data{},
		&pb.SymbolInfoSessionQuoteReply_Error{}, &pb.SymbolInfoSessionQuoteRequest{}, &pb.SymbolInfoSessionTradeData{}, &pb.SymbolInfoSessionTradeReply{},
		&pb.SymbolInfoSessionTradeReply_Data{}, &pb.SymbolInfoSessionTradeReply_Error{}, &pb.SymbolInfoSessionTradeRequest{}, &pb.SymbolInfoStringData{},
		&pb.SymbolInfoStringReply{}, &pb.SymbolInfoStringReply_Data{}, &pb.SymbolInfoStringReply_Error{}, &pb.SymbolInfoStringRequest{},
		&pb.SymbolInfoTickRequest{}, &pb.SymbolInfoTickRequestReply{}, &pb.SymbolInfoTickRequestReply_Data{}, &pb.SymbolInfoTickRequestReply_Error{},
		&pb.SymbolIsSynchronizedData{}, &pb.SymbolIsSynchronizedReply{}, &pb.SymbolIsSynchronizedReply_Data{}, &pb.SymbolIsSynchronizedReply_Error{},
		&pb.SymbolIsSynchronizedRequest{}, &pb.SymbolNameData{}, &pb.SymbolNameReply{}, &pb.SymbolNameReply_Data{},
		&pb.SymbolNameReply_Error{}, &pb.SymbolNameRequest{}, &pb.SymbolParameters{}, &pb.SymbolParamsManyData{},
		&pb.SymbolParamsManyReply{}, &pb.SymbolParamsManyReply_Data{}, &pb.SymbolParamsManyReply_Error{}, &pb.SymbolParamsManyRequest{},
		&pb.SymbolSelectData{}, &pb.SymbolSelectReply{}, &pb.SymbolSelectReply_Data{}, &pb.SymbolSelectReply_Error{},
		&pb.SymbolSelectRequest{}, &pb.SymbolsTotalData{}, &pb.SymbolsTotalReply{}, &pb.SymbolsTotalReply_Data{},
		&pb.SymbolsTotalReply_Error{}, &pb.SymbolsTotalRequest{}, &pb.TerminalHealthCheck{}, &pb.TickSizeSymbol{},
		&pb.TickValueWithSizeData{}, &pb.TickValueWithSizeReply{}, &pb.TickValueWithSizeRequest{},
	}

	// Register all types
	for _, sample := range sampleTypes {
		pi.registerType(reflect.TypeOf(sample))
	}

	// AUTO-GENERATED: Register ALL 67 enums from pb package

	pi.registerEnum("AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE", map[string]int32{
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_BUY": 0,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_BUY_LIMIT": 2,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_BUY_STOP": 4,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_BUY_STOP_LIMIT": 6,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_CLOSE_BY": 8,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_SELL": 1,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_SELL_LIMIT": 3,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_SELL_STOP": 5,
		"AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_SELL_STOP_LIMIT": 7,
	})

	pi.registerEnum("AH_ENUM_POSITIONS_HISTORY_SORT_TYPE", map[string]int32{
		"AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC": 0,
		"AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_DESC": 1,
		"AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_TICKET_ASC": 2,
		"AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_TICKET_DESC": 3,
	})

	pi.registerEnum("AH_SYMBOL_PARAMS_MANY_SORT_TYPE", map[string]int32{
		"AH_SYMBOL_PARAMS_MANY_SORT_TYPE_AH_PARAMS_MANY_SORT_TYPE_MQL_INDEX_ASC": 2,
		"AH_SYMBOL_PARAMS_MANY_SORT_TYPE_AH_PARAMS_MANY_SORT_TYPE_MQL_INDEX_DESC": 3,
		"AH_SYMBOL_PARAMS_MANY_SORT_TYPE_AH_PARAMS_MANY_SORT_TYPE_SYMBOL_NAME_ASC": 0,
		"AH_SYMBOL_PARAMS_MANY_SORT_TYPE_AH_PARAMS_MANY_SORT_TYPE_SYMBOL_NAME_DESC": 1,
	})

	pi.registerEnum("AccountInfoDoublePropertyType", map[string]int32{
		"AccountInfoDoublePropertyType_ACCOUNT_ASSETS": 11,
		"AccountInfoDoublePropertyType_ACCOUNT_BALANCE": 0,
		"AccountInfoDoublePropertyType_ACCOUNT_COMMISSION_BLOCKED": 13,
		"AccountInfoDoublePropertyType_ACCOUNT_CREDIT": 1,
		"AccountInfoDoublePropertyType_ACCOUNT_EQUITY": 3,
		"AccountInfoDoublePropertyType_ACCOUNT_LIABILITIES": 12,
		"AccountInfoDoublePropertyType_ACCOUNT_MARGIN": 4,
		"AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE": 5,
		"AccountInfoDoublePropertyType_ACCOUNT_MARGIN_INITIAL": 9,
		"AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL": 6,
		"AccountInfoDoublePropertyType_ACCOUNT_MARGIN_MAINTENANCE": 10,
		"AccountInfoDoublePropertyType_ACCOUNT_MARGIN_SO_CALL": 7,
		"AccountInfoDoublePropertyType_ACCOUNT_MARGIN_SO_SO": 8,
		"AccountInfoDoublePropertyType_ACCOUNT_PROFIT": 2,
	})

	pi.registerEnum("AccountInfoIntegerPropertyType", map[string]int32{
		"AccountInfoIntegerPropertyType_ACCOUNT_CURRENCY_DIGITS": 8,
		"AccountInfoIntegerPropertyType_ACCOUNT_FIFO_CLOSE": 9,
		"AccountInfoIntegerPropertyType_ACCOUNT_HEDGE_ALLOWED": 10,
		"AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE": 2,
		"AccountInfoIntegerPropertyType_ACCOUNT_LIMIT_ORDERS": 3,
		"AccountInfoIntegerPropertyType_ACCOUNT_LOGIN": 0,
		"AccountInfoIntegerPropertyType_ACCOUNT_MARGIN_MODE": 7,
		"AccountInfoIntegerPropertyType_ACCOUNT_MARGIN_SO_MODE": 4,
		"AccountInfoIntegerPropertyType_ACCOUNT_TRADE_ALLOWED": 5,
		"AccountInfoIntegerPropertyType_ACCOUNT_TRADE_EXPERT": 6,
		"AccountInfoIntegerPropertyType_ACCOUNT_TRADE_MODE": 1,
	})

	pi.registerEnum("AccountInfoStringPropertyType", map[string]int32{
		"AccountInfoStringPropertyType_ACCOUNT_COMPANY": 3,
		"AccountInfoStringPropertyType_ACCOUNT_CURRENCY": 2,
		"AccountInfoStringPropertyType_ACCOUNT_NAME": 0,
		"AccountInfoStringPropertyType_ACCOUNT_SERVER": 1,
	})

	pi.registerEnum("BMT5_ENUM_DAY_OF_WEEK", map[string]int32{
		"BMT5_ENUM_DAY_OF_WEEK_BMT5_FRIDAY": 5,
		"BMT5_ENUM_DAY_OF_WEEK_BMT5_MONDAY": 1,
		"BMT5_ENUM_DAY_OF_WEEK_BMT5_SATURDAY": 6,
		"BMT5_ENUM_DAY_OF_WEEK_BMT5_SUNDAY": 0,
		"BMT5_ENUM_DAY_OF_WEEK_BMT5_THURSDAY": 4,
		"BMT5_ENUM_DAY_OF_WEEK_BMT5_TUESDAY": 2,
		"BMT5_ENUM_DAY_OF_WEEK_BMT5_WEDNESDAY": 3,
	})

	pi.registerEnum("BMT5_ENUM_DEAL_ENTRY_TYPE", map[string]int32{
		"BMT5_ENUM_DEAL_ENTRY_TYPE_BMT5_DEAL_ENTRY_IN": 0,
		"BMT5_ENUM_DEAL_ENTRY_TYPE_BMT5_DEAL_ENTRY_INOUT": 2,
		"BMT5_ENUM_DEAL_ENTRY_TYPE_BMT5_DEAL_ENTRY_OUT": 1,
		"BMT5_ENUM_DEAL_ENTRY_TYPE_BMT5_DEAL_ENTRY_OUT_BY": 3,
	})

	pi.registerEnum("BMT5_ENUM_DEAL_REASON", map[string]int32{
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_CLIENT": 0,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_CORPORATE_ACTION": 10,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_EXPERT": 3,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_MOBILE": 1,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_ROLLOVER": 7,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_SL": 4,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_SO": 6,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_SPLIT": 9,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_TP": 5,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_VMARGIN": 8,
		"BMT5_ENUM_DEAL_REASON_BMT5_DEAL_REASON_WEB": 2,
	})

	pi.registerEnum("BMT5_ENUM_DEAL_TYPE", map[string]int32{
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_DIVIDEND": 15,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_DIVIDEND_FRANKED": 16,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TAX": 17,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_BALANCE": 2,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_BONUS": 6,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_BUY": 0,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_BUY_CANCELED": 13,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_CHARGE": 4,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_COMMISSION": 7,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_COMMISSION_AGENT_DAILY": 10,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_COMMISSION_AGENT_MONTHLY": 11,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_COMMISSION_DAILY": 8,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_COMMISSION_MONTHLY": 9,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_CORRECTION": 5,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_CREDIT": 3,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_INTEREST": 12,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_SELL": 1,
		"BMT5_ENUM_DEAL_TYPE_BMT5_DEAL_TYPE_SELL_CANCELED": 14,
	})

	pi.registerEnum("BMT5_ENUM_OPENED_ORDER_SORT_TYPE", map[string]int32{
		"BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC": 0,
		"BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_DESC": 1,
		"BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_ORDER_TICKET_ID_ASC": 2,
		"BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_ORDER_TICKET_ID_DESC": 3,
	})

	pi.registerEnum("BMT5_ENUM_ORDER_HISTORY_SORT_TYPE", map[string]int32{
		"BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_CLOSE_TIME_ASC": 2,
		"BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_CLOSE_TIME_DESC": 3,
		"BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_ASC": 0,
		"BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_OPEN_TIME_DESC": 1,
		"BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_ORDER_TICKET_ID_ASC": 4,
		"BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_ORDER_TICKET_ID_DESC": 5,
	})

	pi.registerEnum("BMT5_ENUM_ORDER_STATE", map[string]int32{
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_CANCELED": 2,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_EXPIRED": 6,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_FILLED": 4,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_PARTIAL": 3,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_PLACED": 1,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REJECTED": 5,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REQUEST_ADD": 7,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REQUEST_CANCEL": 9,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REQUEST_MODIFY": 8,
		"BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_STARTED": 0,
	})

	pi.registerEnum("BMT5_ENUM_ORDER_TYPE", map[string]int32{
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY": 0,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY_LIMIT": 2,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY_STOP": 4,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY_STOP_LIMIT": 6,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_CLOSE_BY": 8,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL": 1,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL_LIMIT": 3,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL_STOP": 5,
		"BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL_STOP_LIMIT": 7,
	})

	pi.registerEnum("BMT5_ENUM_ORDER_TYPE_FILLING", map[string]int32{
		"BMT5_ENUM_ORDER_TYPE_FILLING_BMT5_ORDER_FILLING_BOC": 3,
		"BMT5_ENUM_ORDER_TYPE_FILLING_BMT5_ORDER_FILLING_FOK": 0,
		"BMT5_ENUM_ORDER_TYPE_FILLING_BMT5_ORDER_FILLING_IOC": 1,
		"BMT5_ENUM_ORDER_TYPE_FILLING_BMT5_ORDER_FILLING_RETURN": 2,
	})

	pi.registerEnum("BMT5_ENUM_ORDER_TYPE_TIME", map[string]int32{
		"BMT5_ENUM_ORDER_TYPE_TIME_BMT5_ORDER_TIME_DAY": 1,
		"BMT5_ENUM_ORDER_TYPE_TIME_BMT5_ORDER_TIME_GTC": 0,
		"BMT5_ENUM_ORDER_TYPE_TIME_BMT5_ORDER_TIME_SPECIFIED": 2,
		"BMT5_ENUM_ORDER_TYPE_TIME_BMT5_ORDER_TIME_SPECIFIED_DAY": 3,
	})

	pi.registerEnum("BMT5_ENUM_POSITION_REASON", map[string]int32{
		"BMT5_ENUM_POSITION_REASON_BMT5_POSITION_REASON_CLIENT": 0,
		"BMT5_ENUM_POSITION_REASON_BMT5_POSITION_REASON_EXPERT": 3,
		"BMT5_ENUM_POSITION_REASON_BMT5_POSITION_REASON_MOBILE": 1,
		"BMT5_ENUM_POSITION_REASON_BMT5_POSITION_REASON_WEB": 2,
		"BMT5_ENUM_POSITION_REASON_ORDER_REASON_SL": 4,
		"BMT5_ENUM_POSITION_REASON_ORDER_REASON_SO": 6,
		"BMT5_ENUM_POSITION_REASON_ORDER_REASON_TP": 5,
	})

	pi.registerEnum("BMT5_ENUM_POSITION_TYPE", map[string]int32{
		"BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_BUY": 0,
		"BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_SELL": 1,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_CALC_MODE", map[string]int32{
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_CFD": 3,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_CFDINDEX": 4,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_CFDLEVERAGE": 5,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_EXCH_BONDS": 9,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_EXCH_BONDS_MOEX": 11,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_EXCH_FUTURES": 7,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_EXCH_FUTURES_FORTS": 8,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_EXCH_STOCKS": 6,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_EXCH_STOCKS_MOEX": 10,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_FOREX": 0,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_FOREX_NO_LEVERAGE": 1,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_FUTURES": 2,
		"BMT5_ENUM_SYMBOL_CALC_MODE_BMT5_SYMBOL_CALC_MODE_SERV_COLLATERAL": 12,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_CHART_MODE", map[string]int32{
		"BMT5_ENUM_SYMBOL_CHART_MODE_BMT5_SYMBOL_CHART_MODE_BID": 0,
		"BMT5_ENUM_SYMBOL_CHART_MODE_BMT5_SYMBOL_CHART_MODE_LAST": 1,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_INDUSTRY", map[string]int32{
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ADVERTISING": 15,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_AEROSPACE_DEFENSE": 95,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_AGRICULTURAL_INPUTS": 1,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_AIRLINES": 96,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_AIRPORTS_SERVICES": 97,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ALUMINIUM": 2,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_APPAREL_MANUFACTURING": 22,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_APPAREL_RETAIL": 23,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ASSETS_MANAGEMENT": 66,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_AUTO_DEALERSHIP": 26,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_AUTO_MANUFACTURERS": 24,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_AUTO_PARTS": 25,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BANKS_DIVERSIFIED": 67,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BANKS_REGIONAL": 68,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BEVERAGES_BREWERS": 45,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BEVERAGES_NON_ALCO": 46,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BEVERAGES_WINERIES": 47,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BIOTECHNOLOGY": 84,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BROADCASTING": 16,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BUILDING_MATERIALS": 3,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BUILDING_PRODUCTS": 98,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_BUSINESS_EQUIPMENT": 99,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CAPITAL_MARKETS": 69,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CHEMICALS": 4,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CLOSE_END_FUND_DEBT": 70,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CLOSE_END_FUND_EQUITY": 71,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CLOSE_END_FUND_FOREIGN": 72,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_COKING_COAL": 5,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_COMMUNICATION_EQUIPMENT": 132,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_COMPUTER_HARDWARE": 133,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CONFECTIONERS": 48,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CONGLOMERATES": 100,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CONSULTING_SERVICES": 101,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CONSUMER_ELECTRONICS": 134,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_COPPER": 6,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_CREDIT_SERVICES": 73,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_DEPARTMENT_STORES": 27,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_DIAGNOSTICS_RESEARCH": 85,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_DISCOUNT_STORES": 49,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_DRUGS_MANUFACTURERS": 86,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_DRUGS_MANUFACTURERS_SPEC": 87,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_EDUCATION_TRAINIG": 50,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ELECTRICAL_EQUIPMENT": 102,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ELECTRONIC_COMPONENTS": 135,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ELECTRONIC_DISTRIBUTION": 136,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ENGINEERING_CONSTRUCTION": 103,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_ENTERTAINMENT": 18,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_EXCHANGE_TRADED_FUND": 65,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FARM_HEAVY_MACHINERY": 104,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FARM_PRODUCTS": 51,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FINANCIAL_CONGLOMERATE": 74,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FINANCIAL_DATA_EXCHANGE": 75,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FOOD_DISTRIBUTION": 52,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FOOTWEAR_ACCESSORIES": 28,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FREIGHT_LOGISTICS": 107,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_FURNISHINGS": 29,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_GAMBLING": 30,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_GAMING_MULTIMEDIA": 17,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_GOLD": 7,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_GROCERY_STORES": 53,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_HEALTHCARE_PLANS": 88,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_HEALTH_INFORMATION": 89,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_HOME_IMPROV_RETAIL": 31,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_HOUSEHOLD_PRODUCTS": 54,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INDUSTRIAL_DISTRIBUTION": 105,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INDUSTRIAL_METALS": 9,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INFRASTRUCTURE_OPERATIONS": 106,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INSURANCE_BROKERS": 76,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INSURANCE_DIVERSIFIED": 77,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INSURANCE_LIFE": 78,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INSURANCE_PROPERTY": 79,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INSURANCE_REINSURANCE": 80,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INSURANCE_SPECIALTY": 81,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INTERNET_CONTENT": 19,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_INTERNET_RETAIL": 32,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_IT_SERVICES": 137,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_LEISURE": 33,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_LODGING": 34,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_LUMBER_WOOD": 8,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_LUXURY_GOODS": 35,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_MARINE_SHIPPING": 108,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_MEDICAL_DEVICES": 91,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_MEDICAL_DISTRIBUTION": 92,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_MEDICAL_FACILITIES": 90,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_MEDICAL_INSTRUMENTS": 93,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_METAL_FABRICATION": 109,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_MORTGAGE_FINANCE": 82,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_OIL_GAS_DRILLING": 57,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_OIL_GAS_EP": 58,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_OIL_GAS_EQUIPMENT": 59,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_OIL_GAS_INTEGRATED": 60,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_OIL_GAS_MIDSTREAM": 61,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_OIL_GAS_REFINING": 62,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_PACKAGED_FOODS": 55,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_PACKAGING_CONTAINERS": 36,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_PAPER": 11,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_PERSONAL_SERVICES": 37,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_PHARM_RETAILERS": 94,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_POLLUTION_CONTROL": 110,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_PRECIOUS_METALS": 10,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_PUBLISHING": 20,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_RAILROADS": 111,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REAL_ESTATE_DEVELOPMENT": 120,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REAL_ESTATE_DIVERSIFIED": 121,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REAL_ESTATE_SERVICES": 122,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_RECREATIONAL_VEHICLES": 38,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_DIVERSIFIED": 123,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_HEALTCARE": 124,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_HOTEL_MOTEL": 125,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_INDUSTRIAL": 126,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_MORTAGE": 127,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_OFFICE": 128,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_RESIDENTAL": 129,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_RETAIL": 130,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_REIT_SPECIALITY": 131,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_RENTAL_LEASING": 112,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_RESIDENT_CONSTRUCTION": 39,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_RESORTS_CASINOS": 40,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_RESTAURANTS": 41,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SCIENTIFIC_INSTRUMENTS": 138,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SECURITY_PROTECTION": 113,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SEMICONDUCTORS": 140,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SEMICONDUCTOR_EQUIPMENT": 139,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SHELL_COMPANIES": 83,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SILVER": 12,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SOFTWARE_APPLICATION": 141,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SOFTWARE_INFRASTRUCTURE": 142,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SOLAR": 143,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SPEALITY_BUSINESS_SERVICES": 114,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SPEALITY_MACHINERY": 115,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SPECIALTY_CHEMICALS": 13,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_SPECIALTY_RETAIL": 42,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_STEEL": 14,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_STUFFING_EMPLOYMENT": 116,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_TELECOM": 21,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_TEXTILE_MANUFACTURING": 43,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_THERMAL_COAL": 63,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_TOBACCO": 56,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_TOOLS_ACCESSORIES": 117,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_TRAVEL_SERVICES": 44,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_TRUCKING": 118,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UNDEFINED": 0,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_URANIUM": 64,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_DIVERSIFIED": 144,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_FIRST": 150,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_LAST": 151,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_POWERPRODUCERS": 145,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_REGULATED_ELECTRIC": 147,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_REGULATED_GAS": 148,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_REGULATED_WATER": 149,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_UTILITIES_RENEWABLE": 146,
		"BMT5_ENUM_SYMBOL_INDUSTRY_BMT5_INDUSTRY_WASTE_MANAGEMENT": 119,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_OPTION_MODE", map[string]int32{
		"BMT5_ENUM_SYMBOL_OPTION_MODE_BMT5_SYMBOL_OPTION_MODE_AMERICAN": 1,
		"BMT5_ENUM_SYMBOL_OPTION_MODE_BMT5_SYMBOL_OPTION_MODE_EUROPEAN": 0,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_OPTION_RIGHT", map[string]int32{
		"BMT5_ENUM_SYMBOL_OPTION_RIGHT_BMT5_SYMBOL_OPTION_RIGHT_CALL": 0,
		"BMT5_ENUM_SYMBOL_OPTION_RIGHT_BMT5_SYMBOL_OPTION_RIGHT_PUT": 1,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_ORDER_GTC_MODE", map[string]int32{
		"BMT5_ENUM_SYMBOL_ORDER_GTC_MODE_BMT5_SYMBOL_ORDERS_DAILY": 1,
		"BMT5_ENUM_SYMBOL_ORDER_GTC_MODE_BMT5_SYMBOL_ORDERS_DAILY_EXCLUDING_STOPS": 2,
		"BMT5_ENUM_SYMBOL_ORDER_GTC_MODE_BMT5_SYMBOL_ORDERS_GTC": 0,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_SECTOR", map[string]int32{
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_BASIC_MATERIALS": 1,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_COMMUNICATION_SERVICES": 2,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_CONSUMER_CYCLICAL": 3,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_CONSUMER_DEFENSIVE": 4,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_CURRENCY": 5,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_CURRENCY_CRYPTO": 6,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_ENERGY": 7,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_FINANCIAL": 8,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_HEALTHCARE": 9,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_INDUSTRIALS": 10,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_REAL_ESTATE": 11,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_TECHNOLOGY": 12,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_UNDEFINED": 0,
		"BMT5_ENUM_SYMBOL_SECTOR_BMT5_SECTOR_UTILITIES": 13,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_SWAP_MODE", map[string]int32{
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_CURRENCY_DEPOSIT": 4,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_CURRENCY_MARGIN": 3,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_CURRENCY_PROFIT": 5,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_CURRENCY_SYMBOL": 2,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_DISABLED": 0,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_INTEREST_CURRENT": 6,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_INTEREST_OPEN": 7,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_POINTS": 1,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_REOPEN_BID": 9,
		"BMT5_ENUM_SYMBOL_SWAP_MODE_BMT5_SYMBOL_SWAP_MODE_REOPEN_CURRENT": 8,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_TRADE_EXECUTION", map[string]int32{
		"BMT5_ENUM_SYMBOL_TRADE_EXECUTION_BMT5_SYMBOL_TRADE_EXECUTION_EXCHANGE": 3,
		"BMT5_ENUM_SYMBOL_TRADE_EXECUTION_BMT5_SYMBOL_TRADE_EXECUTION_INSTANT": 1,
		"BMT5_ENUM_SYMBOL_TRADE_EXECUTION_BMT5_SYMBOL_TRADE_EXECUTION_MARKET": 2,
		"BMT5_ENUM_SYMBOL_TRADE_EXECUTION_BMT5_SYMBOL_TRADE_EXECUTION_REQUEST": 0,
	})

	pi.registerEnum("BMT5_ENUM_SYMBOL_TRADE_MODE", map[string]int32{
		"BMT5_ENUM_SYMBOL_TRADE_MODE_BMT5_SYMBOL_TRADE_MODE_CLOSEONLY": 3,
		"BMT5_ENUM_SYMBOL_TRADE_MODE_BMT5_SYMBOL_TRADE_MODE_DISABLED": 0,
		"BMT5_ENUM_SYMBOL_TRADE_MODE_BMT5_SYMBOL_TRADE_MODE_FULL": 4,
		"BMT5_ENUM_SYMBOL_TRADE_MODE_BMT5_SYMBOL_TRADE_MODE_LONGONLY": 1,
		"BMT5_ENUM_SYMBOL_TRADE_MODE_BMT5_SYMBOL_TRADE_MODE_SHORTONLY": 2,
	})

	pi.registerEnum("BookType", map[string]int32{
		"BookType_BOOK_TYPE_BUY": 1,
		"BookType_BOOK_TYPE_BUY_MARKET": 3,
		"BookType_BOOK_TYPE_SELL": 0,
		"BookType_BOOK_TYPE_SELL_MARKET": 2,
	})

	pi.registerEnum("ChartExpertMode", map[string]int32{
		"ChartExpertMode_ON_OPENED_ORDERS_TICKETS": 3,
		"ChartExpertMode_ON_ORDER_PROFIT": 2,
		"ChartExpertMode_ON_TICK": 0,
		"ChartExpertMode_ON_TRADE": 1,
		"ChartExpertMode_ON_TRADE_TRANSACTION": 4,
	})

	pi.registerEnum("DayOfWeek", map[string]int32{
		"DayOfWeek_FRIDAY": 5,
		"DayOfWeek_MONDAY": 1,
		"DayOfWeek_SATURDAY": 6,
		"DayOfWeek_SUNDAY": 0,
		"DayOfWeek_THURSDAY": 4,
		"DayOfWeek_TUESDAY": 2,
		"DayOfWeek_WEDNESDAY": 3,
	})

	pi.registerEnum("EA_PARAM_TYPE", map[string]int32{
		"EA_PARAM_TYPE_EA_PARAM_TYPE_BOOLEAN": 4,
		"EA_PARAM_TYPE_EA_PARAM_TYPE_DOUBLE": 3,
		"EA_PARAM_TYPE_EA_PARAM_TYPE_INTEGER": 2,
		"EA_PARAM_TYPE_EA_PARAM_TYPE_STRING": 1,
		"EA_PARAM_TYPE_EA_PARAM_TYPE_UNDEFINED": 0,
	})

	pi.registerEnum("ENUM_ORDER_TYPE", map[string]int32{
		"ENUM_ORDER_TYPE_ORDER_TYPE_BUY": 0,
		"ENUM_ORDER_TYPE_ORDER_TYPE_BUY_LIMIT": 2,
		"ENUM_ORDER_TYPE_ORDER_TYPE_BUY_STOP": 4,
		"ENUM_ORDER_TYPE_ORDER_TYPE_BUY_STOP_LIMIT": 6,
		"ENUM_ORDER_TYPE_ORDER_TYPE_CLOSE_BY": 8,
		"ENUM_ORDER_TYPE_ORDER_TYPE_SELL": 1,
		"ENUM_ORDER_TYPE_ORDER_TYPE_SELL_LIMIT": 3,
		"ENUM_ORDER_TYPE_ORDER_TYPE_SELL_STOP": 5,
		"ENUM_ORDER_TYPE_ORDER_TYPE_SELL_STOP_LIMIT": 7,
	})

	pi.registerEnum("ENUM_ORDER_TYPE_TF", map[string]int32{
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY": 0,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_LIMIT": 2,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_STOP": 4,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_STOP_LIMIT": 6,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_CLOSE_BY": 8,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL": 1,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_LIMIT": 3,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_STOP": 5,
		"ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_STOP_LIMIT": 7,
	})

	pi.registerEnum("EnumOpenChartWithEaChatPeriod", map[string]int32{
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_CURRENT": 0,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_D1": 19,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_H1": 12,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_H12": 18,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_H2": 13,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_H3": 14,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_H4": 15,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_H6": 16,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_H8": 17,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M1": 1,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M10": 7,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M12": 8,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M15": 9,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M2": 2,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M20": 10,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M3": 3,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M30": 11,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M4": 4,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M5": 5,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_M6": 6,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_MN1": 21,
		"EnumOpenChartWithEaChatPeriod_EA_CHART_PERIOD_W1": 20,
	})

	pi.registerEnum("EnumOpenChartWithEaParemeterType", map[string]int32{
		"EnumOpenChartWithEaParemeterType_EA_PARAM_BOOL": 2,
		"EnumOpenChartWithEaParemeterType_EA_PARAM_DOUBLE": 4,
		"EnumOpenChartWithEaParemeterType_EA_PARAM_INT": 0,
		"EnumOpenChartWithEaParemeterType_EA_PARAM_LONG": 1,
		"EnumOpenChartWithEaParemeterType_EA_PARAM_STRING": 3,
	})

	pi.registerEnum("EnumOpenTerminalChartWithEaChatPeriod", map[string]int32{
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_CURRENT": 0,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_D1": 19,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_H1": 12,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_H12": 18,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_H2": 13,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_H3": 14,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_H4": 15,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_H6": 16,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_H8": 17,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M1": 1,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M10": 7,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M12": 8,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M15": 9,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M2": 2,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M20": 10,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M3": 3,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M30": 11,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M4": 4,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M5": 5,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_M6": 6,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_MN1": 21,
		"EnumOpenTerminalChartWithEaChatPeriod_MRPC_EA_CHART_PERIOD_W1": 20,
	})

	pi.registerEnum("EnumOpenTerminalChartWithEaParameterType", map[string]int32{
		"EnumOpenTerminalChartWithEaParameterType_MRPC_EA_PARAM_BOOL": 2,
		"EnumOpenTerminalChartWithEaParameterType_MRPC_EA_PARAM_DOUBLE": 4,
		"EnumOpenTerminalChartWithEaParameterType_MRPC_EA_PARAM_INT": 0,
		"EnumOpenTerminalChartWithEaParameterType_MRPC_EA_PARAM_LONG": 1,
		"EnumOpenTerminalChartWithEaParameterType_MRPC_EA_PARAM_STRING": 3,
	})

	pi.registerEnum("ErrorType", map[string]int32{
		"ErrorType_MQL_CUSTOM_EXECUTION": 9,
		"ErrorType_MQL_EXECUTION": 8,
		"ErrorType_MQL_TRADE_EXECUTION": 10,
		"ErrorType_MRPC": 1,
		"ErrorType_MRPC_TIMEOUT": 3,
		"ErrorType_MRPC_VALIDATION": 2,
		"ErrorType_TERMINAL_API_GET_PARAMS": 5,
		"ErrorType_TERMINAL_API_SET_RESULT": 6,
		"ErrorType_TERMINAL_API_TIMEOUT": 4,
		"ErrorType_TERMINAL_API_UNEXPECTED": 7,
		"ErrorType_UNDEFINED": 0,
	})

	pi.registerEnum("LogType", map[string]int32{
		"LogType_Global": 0,
		"LogType_MQL4": 2,
		"LogType_MQL5": 1,
	})

	pi.registerEnum("MRPC_ENUM_ORDER_TYPE_FILLING", map[string]int32{
		"MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_BOC": 3,
		"MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK": 0,
		"MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_IOC": 1,
		"MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_RETURN": 2,
	})

	pi.registerEnum("MRPC_ENUM_ORDER_TYPE_TIME", map[string]int32{
		"MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_DAY": 1,
		"MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC": 0,
		"MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED": 2,
		"MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED_DAY": 3,
	})

	pi.registerEnum("MRPC_ENUM_TRADE_REQUEST_ACTIONS", map[string]int32{
		"MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_CLOSE_BY": 5,
		"MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL": 0,
		"MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_MODIFY": 3,
		"MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING": 1,
		"MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_REMOVE": 4,
		"MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_SLTP": 2,
	})

	pi.registerEnum("MRPC_ORDER_CLOSE_MODE", map[string]int32{
		"MRPC_ORDER_CLOSE_MODE_MRPC_MARKET_ORDER_CLOSE": 0,
		"MRPC_ORDER_CLOSE_MODE_MRPC_MARKET_ORDER_PARTIAL_CLOSE": 1,
		"MRPC_ORDER_CLOSE_MODE_MRPC_PENDING_ORDER_REMOVE": 2,
	})

	pi.registerEnum("MT5_SUB_ENUM_EVENT_GROUP_TYPE", map[string]int32{
		"MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit": 0,
		"MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate": 1,
	})

	pi.registerEnum("MqlErrorCode", map[string]int32{
		"MqlErrorCode_ERR_ACCOUNT_WRONG_PROPERTY": 4701,
		"MqlErrorCode_ERR_ARRAY_BAD_SIZE": 4011,
		"MqlErrorCode_ERR_ARRAY_RESIZE_ERROR": 4007,
		"MqlErrorCode_ERR_BOOKS_CANNOT_ADD": 4901,
		"MqlErrorCode_ERR_BOOKS_CANNOT_DELETE": 4902,
		"MqlErrorCode_ERR_BOOKS_CANNOT_GET": 4903,
		"MqlErrorCode_ERR_BOOKS_CANNOT_SUBSCRIBE": 4904,
		"MqlErrorCode_ERR_BUFFERS_NO_MEMORY": 4601,
		"MqlErrorCode_ERR_BUFFERS_WRONG_INDEX": 4602,
		"MqlErrorCode_ERR_CALENDAR_MORE_DATA": 5400,
		"MqlErrorCode_ERR_CALENDAR_NO_DATA": 5402,
		"MqlErrorCode_ERR_CALENDAR_TIMEOUT": 5401,
		"MqlErrorCode_ERR_CANNOT_DELETE_FILE": 5006,
		"MqlErrorCode_ERR_CANNOT_OPEN_FILE": 5004,
		"MqlErrorCode_ERR_CHART_CANNOT_CHANGE": 4106,
		"MqlErrorCode_ERR_CHART_CANNOT_CREATE_TIMER": 4108,
		"MqlErrorCode_ERR_CHART_CANNOT_OPEN": 4105,
		"MqlErrorCode_ERR_CHART_INDICATOR_CANNOT_ADD": 4114,
		"MqlErrorCode_ERR_CHART_INDICATOR_CANNOT_DEL": 4115,
		"MqlErrorCode_ERR_CHART_INDICATOR_NOT_FOUND": 4116,
		"MqlErrorCode_ERR_CHART_NAVIGATE_FAILED": 4111,
		"MqlErrorCode_ERR_CHART_NOT_FOUND": 4103,
		"MqlErrorCode_ERR_CHART_NO_EXPERT": 4104,
		"MqlErrorCode_ERR_CHART_NO_REPLY": 4102,
		"MqlErrorCode_ERR_CHART_SCREENSHOT_FAILED": 4110,
		"MqlErrorCode_ERR_CHART_TEMPLATE_FAILED": 4112,
		"MqlErrorCode_ERR_CHART_WINDOW_NOT_FOUND": 4113,
		"MqlErrorCode_ERR_CHART_WRONG_ID": 4101,
		"MqlErrorCode_ERR_CHART_WRONG_PARAMETER": 4107,
		"MqlErrorCode_ERR_CHART_WRONG_PROPERTY": 4109,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_ERROR": 5305,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_EXIST": 5304,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_NAME_LONG": 5302,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_PARAMETER_ERROR": 5308,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_PARAMETER_LONG": 5309,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_PATH_LONG": 5303,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_PROPERTY_WRONG": 5307,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_SELECTED": 5306,
		"MqlErrorCode_ERR_CUSTOM_SYMBOL_WRONG_NAME": 5301,
		"MqlErrorCode_ERR_CUSTOM_TICKS_WRONG_ORDER": 5310,
		"MqlErrorCode_ERR_CUSTOM_WRONG_PROPERTY": 4603,
		"MqlErrorCode_ERR_DATABASE_ABORT": 5604,
		"MqlErrorCode_ERR_DATABASE_AUTH": 5623,
		"MqlErrorCode_ERR_DATABASE_BIND_PARAMETERS": 5129,
		"MqlErrorCode_ERR_DATABASE_BUSY": 5605,
		"MqlErrorCode_ERR_DATABASE_CANTOPEN": 5614,
		"MqlErrorCode_ERR_DATABASE_CONNECT": 5123,
		"MqlErrorCode_ERR_DATABASE_CONSTRAINT": 5619,
		"MqlErrorCode_ERR_DATABASE_CORRUPT": 5611,
		"MqlErrorCode_ERR_DATABASE_EMPTY": 5616,
		"MqlErrorCode_ERR_DATABASE_ERROR": 5601,
		"MqlErrorCode_ERR_DATABASE_EXECUTE": 5124,
		"MqlErrorCode_ERR_DATABASE_FORMAT": 5624,
		"MqlErrorCode_ERR_DATABASE_FULL": 5613,
		"MqlErrorCode_ERR_DATABASE_INTERNAL": 5120,
		"MqlErrorCode_ERR_DATABASE_INTERRUPT": 5609,
		"MqlErrorCode_ERR_DATABASE_INVALID_HANDLE": 5121,
		"MqlErrorCode_ERR_DATABASE_IOERR": 5610,
		"MqlErrorCode_ERR_DATABASE_LOCKED": 5606,
		"MqlErrorCode_ERR_DATABASE_LOGIC": 5602,
		"MqlErrorCode_ERR_DATABASE_MISMATCH": 5620,
		"MqlErrorCode_ERR_DATABASE_MISUSE": 5621,
		"MqlErrorCode_ERR_DATABASE_NOLFS": 5622,
		"MqlErrorCode_ERR_DATABASE_NOMEM": 5607,
		"MqlErrorCode_ERR_DATABASE_NOTADB": 5626,
		"MqlErrorCode_ERR_DATABASE_NOTFOUND": 5612,
		"MqlErrorCode_ERR_DATABASE_NOT_READY": 5128,
		"MqlErrorCode_ERR_DATABASE_NO_MORE_DATA": 5126,
		"MqlErrorCode_ERR_DATABASE_PERM": 5603,
		"MqlErrorCode_ERR_DATABASE_PREPARE": 5125,
		"MqlErrorCode_ERR_DATABASE_PROTOCOL": 5615,
		"MqlErrorCode_ERR_DATABASE_RANGE": 5625,
		"MqlErrorCode_ERR_DATABASE_READONLY": 5608,
		"MqlErrorCode_ERR_DATABASE_SCHEMA": 5617,
		"MqlErrorCode_ERR_DATABASE_STEP": 5127,
		"MqlErrorCode_ERR_DATABASE_TOOBIG": 5618,
		"MqlErrorCode_ERR_DATABASE_TOO_MANY_OBJECTS": 5122,
		"MqlErrorCode_ERR_FILE_CACHEBUFFER_ERROR": 5005,
		"MqlErrorCode_ERR_FTP_CHANGEDIR": 4523,
		"MqlErrorCode_ERR_FTP_CONNECT_FAILED": 4522,
		"MqlErrorCode_ERR_FTP_FILE_ERROR": 4521,
		"MqlErrorCode_ERR_FTP_NOLOGIN": 4520,
		"MqlErrorCode_ERR_FTP_NOSERVER": 4519,
		"MqlErrorCode_ERR_FTP_SEND_FAILED": 4514,
		"MqlErrorCode_ERR_FUNCTION_NOT_ALLOWED": 4014,
		"MqlErrorCode_ERR_GLOBALVARIABLE_CANNOTREAD": 4504,
		"MqlErrorCode_ERR_GLOBALVARIABLE_CANNOTWRITE": 4505,
		"MqlErrorCode_ERR_GLOBALVARIABLE_EXISTS": 4502,
		"MqlErrorCode_ERR_GLOBALVARIABLE_NOT_FOUND": 4501,
		"MqlErrorCode_ERR_GLOBALVARIABLE_NOT_MODIFIED": 4503,
		"MqlErrorCode_ERR_HISTORY_BARS_LIMIT": 4404,
		"MqlErrorCode_ERR_HISTORY_LOAD_ERRORS": 4405,
		"MqlErrorCode_ERR_HISTORY_NOT_FOUND": 4401,
		"MqlErrorCode_ERR_HISTORY_SMALL_BUFFER": 4407,
		"MqlErrorCode_ERR_HISTORY_TIMEOUT": 4403,
		"MqlErrorCode_ERR_HISTORY_WRONG_PROPERTY": 4402,
		"MqlErrorCode_ERR_INCOMPATIBLE_ARRAYS": 5050,
		"MqlErrorCode_ERR_INDICATOR_CANNOT_ADD": 4805,
		"MqlErrorCode_ERR_INDICATOR_CANNOT_APPLY": 4804,
		"MqlErrorCode_ERR_INDICATOR_CANNOT_CREATE": 4802,
		"MqlErrorCode_ERR_INDICATOR_DATA_NOT_FOUND": 4806,
		"MqlErrorCode_ERR_INDICATOR_NO_MEMORY": 4803,
		"MqlErrorCode_ERR_INDICATOR_UNKNOWN_SYMBOL": 4801,
		"MqlErrorCode_ERR_INDICATOR_WRONG_HANDLE": 4807,
		"MqlErrorCode_ERR_INDICATOR_WRONG_PARAMETERS": 4808,
		"MqlErrorCode_ERR_INTERNAL_ERROR": 4001,
		"MqlErrorCode_ERR_INVALID_ARRAY": 4006,
		"MqlErrorCode_ERR_INVALID_DATETIME": 4010,
		"MqlErrorCode_ERR_INVALID_FILEHANDLE": 5007,
		"MqlErrorCode_ERR_INVALID_HANDLE": 4024,
		"MqlErrorCode_ERR_INVALID_PARAMETER": 4003,
		"MqlErrorCode_ERR_INVALID_POINTER": 4012,
		"MqlErrorCode_ERR_INVALID_POINTER_TYPE": 4013,
		"MqlErrorCode_ERR_INVALID_TYPE": 4023,
		"MqlErrorCode_ERR_MAIL_SEND_FAILED": 4510,
		"MqlErrorCode_ERR_MARKET_LASTTIME_UNKNOWN": 4304,
		"MqlErrorCode_ERR_MARKET_NOT_SELECTED": 4302,
		"MqlErrorCode_ERR_MARKET_SELECT_ERROR": 4305,
		"MqlErrorCode_ERR_MARKET_SELECT_LIMIT": 4306,
		"MqlErrorCode_ERR_MARKET_SESSION_INDEX": 4307,
		"MqlErrorCode_ERR_MARKET_UNKNOWN_SYMBOL": 4301,
		"MqlErrorCode_ERR_MARKET_WRONG_PROPERTY": 4303,
		"MqlErrorCode_ERR_MATH_OVERFLOW": 4019,
		"MqlErrorCode_ERR_MATRIX_CONTAINS_NAN": 5706,
		"MqlErrorCode_ERR_MATRIX_FUNC_NOT_ALLOWED": 5705,
		"MqlErrorCode_ERR_MATRIX_INCONSISTENT": 5702,
		"MqlErrorCode_ERR_MATRIX_INTERNAL": 5700,
		"MqlErrorCode_ERR_MATRIX_INVALID_SIZE": 5703,
		"MqlErrorCode_ERR_MATRIX_INVALID_TYPE": 5704,
		"MqlErrorCode_ERR_MATRIX_NOT_INITIALIZED": 5701,
		"MqlErrorCode_ERR_MQL5_WRONG_PROPERTY": 4512,
		"MqlErrorCode_ERR_NETSOCKET_CANNOT_CONNECT": 5272,
		"MqlErrorCode_ERR_NETSOCKET_HANDSHAKE_FAILED": 5274,
		"MqlErrorCode_ERR_NETSOCKET_INVALIDHANDLE": 5270,
		"MqlErrorCode_ERR_NETSOCKET_IO_ERROR": 5273,
		"MqlErrorCode_ERR_NETSOCKET_NO_CERTIFICATE": 5275,
		"MqlErrorCode_ERR_NETSOCKET_TOO_MANY_OPENED": 5271,
		"MqlErrorCode_ERR_NOTIFICATION_SEND_FAILED": 4515,
		"MqlErrorCode_ERR_NOTIFICATION_TOO_FREQUENT": 4518,
		"MqlErrorCode_ERR_NOTIFICATION_WRONG_PARAMETER": 4516,
		"MqlErrorCode_ERR_NOTIFICATION_WRONG_SETTINGS": 4517,
		"MqlErrorCode_ERR_NOTINITIALIZED_STRING": 4009,
		"MqlErrorCode_ERR_NOT_CUSTOM_SYMBOL": 5300,
		"MqlErrorCode_ERR_NOT_ENOUGH_MEMORY": 4004,
		"MqlErrorCode_ERR_NO_STRING_DATE": 5030,
		"MqlErrorCode_ERR_NUMBER_ARRAYS_ONLY": 5054,
		"MqlErrorCode_ERR_OBJECT_ERROR": 4201,
		"MqlErrorCode_ERR_OBJECT_GETDATE_FAILED": 4204,
		"MqlErrorCode_ERR_OBJECT_GETVALUE_FAILED": 4205,
		"MqlErrorCode_ERR_OBJECT_NOT_FOUND": 4202,
		"MqlErrorCode_ERR_OBJECT_WRONG_PROPERTY": 4203,
		"MqlErrorCode_ERR_ONNX_INTERNAL": 5800,
		"MqlErrorCode_ERR_ONNX_INVALID_PARAMETER": 5805,
		"MqlErrorCode_ERR_ONNX_INVALID_PARAMETERS_COUNT": 5804,
		"MqlErrorCode_ERR_ONNX_INVALID_PARAMETER_SIZE": 5807,
		"MqlErrorCode_ERR_ONNX_INVALID_PARAMETER_TYPE": 5806,
		"MqlErrorCode_ERR_ONNX_NOT_INITIALIZED": 5801,
		"MqlErrorCode_ERR_ONNX_NOT_SUPPORTED": 5802,
		"MqlErrorCode_ERR_ONNX_RUN_FAILED": 5803,
		"MqlErrorCode_ERR_ONNX_WRONG_DIMENSION": 5808,
		"MqlErrorCode_ERR_OPENCL_BUFFER_CREATE": 5112,
		"MqlErrorCode_ERR_OPENCL_CONTEXT_CREATE": 5103,
		"MqlErrorCode_ERR_OPENCL_EXECUTE": 5109,
		"MqlErrorCode_ERR_OPENCL_INTERNAL": 5101,
		"MqlErrorCode_ERR_OPENCL_INVALID_HANDLE": 5102,
		"MqlErrorCode_ERR_OPENCL_KERNEL_CREATE": 5107,
		"MqlErrorCode_ERR_OPENCL_NOT_SUPPORTED": 5100,
		"MqlErrorCode_ERR_OPENCL_PROGRAM_CREATE": 5105,
		"MqlErrorCode_ERR_OPENCL_QUEUE_CREATE": 5104,
		"MqlErrorCode_ERR_OPENCL_SELECTDEVICE": 5114,
		"MqlErrorCode_ERR_OPENCL_SET_KERNEL_PARAMETER": 5108,
		"MqlErrorCode_ERR_OPENCL_TOO_LONG_KERNEL_NAME": 5106,
		"MqlErrorCode_ERR_OPENCL_TOO_MANY_OBJECTS": 5113,
		"MqlErrorCode_ERR_OPENCL_WRONG_BUFFER_OFFSET": 5111,
		"MqlErrorCode_ERR_OPENCL_WRONG_BUFFER_SIZE": 5110,
		"MqlErrorCode_ERR_PLAY_SOUND_FAILED": 4511,
		"MqlErrorCode_ERR_PROGRAM_STOPPED": 4022,
		"MqlErrorCode_ERR_RESOURCE_NAME_DUPLICATED": 4015,
		"MqlErrorCode_ERR_RESOURCE_NAME_IS_TOO_LONG": 4018,
		"MqlErrorCode_ERR_RESOURCE_NOT_FOUND": 4016,
		"MqlErrorCode_ERR_RESOURCE_UNSUPPORTED_TYPE": 4017,
		"MqlErrorCode_ERR_SLEEP_ERROR": 4020,
		"MqlErrorCode_ERR_SMALL_ARRAY": 5052,
		"MqlErrorCode_ERR_SMALL_ASSERIES_ARRAY": 5051,
		"MqlErrorCode_ERR_STRING_OUT_OF_MEMORY": 5034,
		"MqlErrorCode_ERR_STRING_RESIZE_ERROR": 4008,
		"MqlErrorCode_ERR_STRING_TIME_ERROR": 5033,
		"MqlErrorCode_ERR_STRUCT_WITHOBJECTS_ORCLASS": 4005,
		"MqlErrorCode_ERR_SUCCESS": 0,
		"MqlErrorCode_ERR_TERMINAL_WRONG_PROPERTY": 4513,
		"MqlErrorCode_ERR_TOO_LONG_FILENAME": 5003,
		"MqlErrorCode_ERR_TOO_MANY_FILES": 5001,
		"MqlErrorCode_ERR_TOO_MANY_OBJECTS": 4025,
		"MqlErrorCode_ERR_TRADE_CALC_FAILED": 4758,
		"MqlErrorCode_ERR_TRADE_DEAL_NOT_FOUND": 4755,
		"MqlErrorCode_ERR_TRADE_DISABLED": 4752,
		"MqlErrorCode_ERR_TRADE_ORDER_NOT_FOUND": 4754,
		"MqlErrorCode_ERR_TRADE_POSITION_NOT_FOUND": 4753,
		"MqlErrorCode_ERR_TRADE_SEND_FAILED": 4756,
		"MqlErrorCode_ERR_TRADE_WRONG_PROPERTY": 4751,
		"MqlErrorCode_ERR_USER_ERROR_FIRST": 65536,
		"MqlErrorCode_ERR_WEBREQUEST_CONNECT_FAILED": 5201,
		"MqlErrorCode_ERR_WEBREQUEST_INVALID_ADDRESS": 5200,
		"MqlErrorCode_ERR_WEBREQUEST_REQUEST_FAILED": 5203,
		"MqlErrorCode_ERR_WEBREQUEST_TIMEOUT": 5202,
		"MqlErrorCode_ERR_WRONG_FILEHANDLE": 5008,
		"MqlErrorCode_ERR_WRONG_FILENAME": 5002,
		"MqlErrorCode_ERR_WRONG_INTERNAL_PARAMETER": 4002,
		"MqlErrorCode_ERR_WRONG_STRING_DATE": 5031,
		"MqlErrorCode_ERR_WRONG_STRING_TIME": 5032,
		"MqlErrorCode_ERR_ZEROSIZE_ARRAY": 5053,
	})

	pi.registerEnum("MqlErrorTradeCode", map[string]int32{
		"MqlErrorTradeCode_TRADE_RETCODE_CANCEL": 10007,
		"MqlErrorTradeCode_TRADE_RETCODE_CLIENT_DISABLES_AT": 10027,
		"MqlErrorTradeCode_TRADE_RETCODE_CLOSE_ONLY": 10044,
		"MqlErrorTradeCode_TRADE_RETCODE_CLOSE_ORDER_EXIST": 10039,
		"MqlErrorTradeCode_TRADE_RETCODE_CONNECTION": 10031,
		"MqlErrorTradeCode_TRADE_RETCODE_DONE": 10009,
		"MqlErrorTradeCode_TRADE_RETCODE_DONE_PARTIAL": 10010,
		"MqlErrorTradeCode_TRADE_RETCODE_ERROR": 10011,
		"MqlErrorTradeCode_TRADE_RETCODE_FIFO_CLOSE": 10045,
		"MqlErrorTradeCode_TRADE_RETCODE_FROZEN": 10029,
		"MqlErrorTradeCode_TRADE_RETCODE_HEDGE_PROHIBITED": 10046,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID": 10013,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID_CLOSE_VOLUME": 10038,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID_EXPIRATION": 10022,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID_FILL": 10030,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID_ORDER": 10035,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID_PRICE": 10015,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID_STOPS": 10016,
		"MqlErrorTradeCode_TRADE_RETCODE_INVALID_VOLUME": 10014,
		"MqlErrorTradeCode_TRADE_RETCODE_LIMIT_ORDERS": 10033,
		"MqlErrorTradeCode_TRADE_RETCODE_LIMIT_POSITIONS": 10040,
		"MqlErrorTradeCode_TRADE_RETCODE_LIMIT_VOLUME": 10034,
		"MqlErrorTradeCode_TRADE_RETCODE_LOCKED": 10028,
		"MqlErrorTradeCode_TRADE_RETCODE_LONG_ONLY": 10042,
		"MqlErrorTradeCode_TRADE_RETCODE_MARKET_CLOSED": 10018,
		"MqlErrorTradeCode_TRADE_RETCODE_NO_CHANGES": 10025,
		"MqlErrorTradeCode_TRADE_RETCODE_NO_MONEY": 10019,
		"MqlErrorTradeCode_TRADE_RETCODE_ONLY_REAL": 10032,
		"MqlErrorTradeCode_TRADE_RETCODE_ORDER_CHANGED": 10023,
		"MqlErrorTradeCode_TRADE_RETCODE_PLACED": 10008,
		"MqlErrorTradeCode_TRADE_RETCODE_POSITION_CLOSED": 10036,
		"MqlErrorTradeCode_TRADE_RETCODE_PRICE_CHANGED": 10020,
		"MqlErrorTradeCode_TRADE_RETCODE_PRICE_OFF": 10021,
		"MqlErrorTradeCode_TRADE_RETCODE_REJECT": 10006,
		"MqlErrorTradeCode_TRADE_RETCODE_REJECT_CANCEL": 10041,
		"MqlErrorTradeCode_TRADE_RETCODE_REQUOTE": 10004,
		"MqlErrorTradeCode_TRADE_RETCODE_SERVER_DISABLES_AT": 10026,
		"MqlErrorTradeCode_TRADE_RETCODE_SHORT_ONLY": 10043,
		"MqlErrorTradeCode_TRADE_RETCODE_SUCCESS": 0,
		"MqlErrorTradeCode_TRADE_RETCODE_TIMEOUT": 10012,
		"MqlErrorTradeCode_TRADE_RETCODE_TOO_MANY_REQUESTS": 10024,
		"MqlErrorTradeCode_TRADE_RETCODE_TRADE_DISABLED": 10017,
	})

	pi.registerEnum("MrpcEnumAccountTradeMode", map[string]int32{
		"MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_CONTEST": 1,
		"MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_DEMO": 0,
		"MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_REAL": 2,
	})

	pi.registerEnum("ProxyTypes", map[string]int32{
		"ProxyTypes_Https": 1,
		"ProxyTypes_None": 0,
		"ProxyTypes_Socks4": 2,
		"ProxyTypes_Socks5": 3,
	})

	pi.registerEnum("SUB_ENUM_DEAL_ENTRY", map[string]int32{
		"SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_IN": 0,
		"SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_INOUT": 2,
		"SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_OUT": 1,
		"SUB_ENUM_DEAL_ENTRY_SUB_DEAL_ENTRY_OUT_BY": 3,
	})

	pi.registerEnum("SUB_ENUM_DEAL_REASON", map[string]int32{
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_CLIENT": 0,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_CORPORATE_ACTION": 10,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_EXPERT": 3,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_MOBILE": 1,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_ROLLOVER": 7,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_SL": 4,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_SO": 6,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_SPLIT": 9,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_TP": 5,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_VMARGIN": 8,
		"SUB_ENUM_DEAL_REASON_SUB_DEAL_REASON_WEB": 2,
	})

	pi.registerEnum("SUB_ENUM_DEAL_TYPE", map[string]int32{
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_DIVIDEND": 15,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_DIVIDEND_FRANKED": 16,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TAX": 17,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BALANCE": 2,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BONUS": 6,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY": 0,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY_CANCELED": 13,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_CHARGE": 4,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION": 7,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_AGENT_DAILY": 10,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_AGENT_MONTHLY": 11,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_DAILY": 8,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_COMMISSION_MONTHLY": 9,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_CORRECTION": 5,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_CREDIT": 3,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_INTEREST": 12,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_SELL": 1,
		"SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_SELL_CANCELED": 14,
	})

	pi.registerEnum("SUB_ENUM_ORDER_REASON", map[string]int32{
		"SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_CLIENT": 0,
		"SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_EXPERT": 4,
		"SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_MOBILE": 2,
		"SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_SL": 5,
		"SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_SO": 7,
		"SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_TP": 6,
		"SUB_ENUM_ORDER_REASON_SUB_ORDER_REASON_WEB": 3,
	})

	pi.registerEnum("SUB_ENUM_ORDER_STATE", map[string]int32{
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_CANCELED": 2,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_EXPIRED": 6,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_FILLED": 4,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_PARTIAL": 3,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_PLACED": 1,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REJECTED": 5,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_ADD": 7,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_CANCEL": 9,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_REQUEST_MODIFY": 8,
		"SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_STARTED": 0,
	})

	pi.registerEnum("SUB_ENUM_ORDER_TYPE", map[string]int32{
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY": 0,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_LIMIT": 2,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_STOP": 4,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_BUY_STOP_LIMIT": 6,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_CLOSE_BY": 8,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL": 1,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_LIMIT": 3,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_STOP": 5,
		"SUB_ENUM_ORDER_TYPE_SUB_ORDER_TYPE_SELL_STOP_LIMIT": 7,
	})

	pi.registerEnum("SUB_ENUM_ORDER_TYPE_FILLING", map[string]int32{
		"SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_BOC": 2,
		"SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_FOK": 0,
		"SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_IOC": 1,
		"SUB_ENUM_ORDER_TYPE_FILLING_SUB_ORDER_FILLING_RETURN": 3,
	})

	pi.registerEnum("SUB_ENUM_ORDER_TYPE_TIME", map[string]int32{
		"SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_DAY": 1,
		"SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_GTC": 0,
		"SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED": 2,
		"SUB_ENUM_ORDER_TYPE_TIME_SUB_ORDER_TIME_SPECIFIED_DAY": 3,
	})

	pi.registerEnum("SUB_ENUM_POSITION_REASON", map[string]int32{
		"SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_CLIENT": 0,
		"SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_EXPERT": 4,
		"SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_MOBILE": 2,
		"SUB_ENUM_POSITION_REASON_SUB_POSITION_REASON_WEB": 3,
	})

	pi.registerEnum("SUB_ENUM_POSITION_TYPE", map[string]int32{
		"SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_BUY": 0,
		"SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_SELL": 1,
	})

	pi.registerEnum("SUB_ENUM_TRADE_REQUEST_ACTIONS", map[string]int32{
		"SUB_ENUM_TRADE_REQUEST_ACTIONS_SUB_TRADE_ACTION_CLOSE_BY": 6,
		"SUB_ENUM_TRADE_REQUEST_ACTIONS_SUB_TRADE_ACTION_DEAL": 1,
		"SUB_ENUM_TRADE_REQUEST_ACTIONS_SUB_TRADE_ACTION_MODIFY": 4,
		"SUB_ENUM_TRADE_REQUEST_ACTIONS_SUB_TRADE_ACTION_PENDING": 2,
		"SUB_ENUM_TRADE_REQUEST_ACTIONS_SUB_TRADE_ACTION_REMOVE": 5,
		"SUB_ENUM_TRADE_REQUEST_ACTIONS_SUB_TRADE_ACTION_SLTP": 3,
		"SUB_ENUM_TRADE_REQUEST_ACTIONS_SUB_TRADE_ACTION_UNDEFINED": 0,
	})

	pi.registerEnum("SUB_ENUM_TRADE_TRANSACTION_TYPE", map[string]int32{
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_DEAL_ADD": 3,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_DEAL_DELETE": 5,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_DEAL_UPDATE": 4,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_HISTORY_ADD": 6,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_HISTORY_DELETE": 8,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_HISTORY_UPDATE": 7,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_ORDER_ADD": 0,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_ORDER_DELETE": 2,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_ORDER_UPDATE": 1,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_POSITION": 9,
		"SUB_ENUM_TRADE_TRANSACTION_TYPE_SUB_TRADE_TRANSACTION_REQUEST": 10,
	})

	pi.registerEnum("SymbolInfoDoubleProperty", map[string]int32{
		"SymbolInfoDoubleProperty_SYMBOL_ASK": 3,
		"SymbolInfoDoubleProperty_SYMBOL_ASKHIGH": 4,
		"SymbolInfoDoubleProperty_SYMBOL_ASKLOW": 5,
		"SymbolInfoDoubleProperty_SYMBOL_BID": 0,
		"SymbolInfoDoubleProperty_SYMBOL_BIDHIGH": 1,
		"SymbolInfoDoubleProperty_SYMBOL_BIDLOW": 2,
		"SymbolInfoDoubleProperty_SYMBOL_COUNT": 59,
		"SymbolInfoDoubleProperty_SYMBOL_LAST": 6,
		"SymbolInfoDoubleProperty_SYMBOL_LASTHIGH": 7,
		"SymbolInfoDoubleProperty_SYMBOL_LASTLOW": 8,
		"SymbolInfoDoubleProperty_SYMBOL_MARGIN_HEDGED": 48,
		"SymbolInfoDoubleProperty_SYMBOL_MARGIN_INITIAL": 35,
		"SymbolInfoDoubleProperty_SYMBOL_MARGIN_MAINTENANCE": 36,
		"SymbolInfoDoubleProperty_SYMBOL_OPTION_STRIKE": 12,
		"SymbolInfoDoubleProperty_SYMBOL_POINT": 13,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_CHANGE": 49,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_DELTA": 52,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_GAMMA": 54,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_OMEGA": 57,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_RHO": 56,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_SENSITIVITY": 58,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_THEORETICAL": 51,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_THETA": 53,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_VEGA": 55,
		"SymbolInfoDoubleProperty_SYMBOL_PRICE_VOLATILITY": 50,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_AW": 44,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_BUY_ORDERS_VOLUME": 40,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_CLOSE": 43,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_INTEREST": 39,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_OPEN": 42,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_PRICE_LIMIT_MAX": 47,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_PRICE_LIMIT_MIN": 46,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_PRICE_SETTLEMENT": 45,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_SELL_ORDERS_VOLUME": 41,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_TURNOVER": 38,
		"SymbolInfoDoubleProperty_SYMBOL_SESSION_VOLUME": 37,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_FRIDAY": 33,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_LONG": 26,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_MONDAY": 29,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_SATURDAY": 34,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_SHORT": 27,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_SUNDAY": 28,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_THURSDAY": 32,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_TUESDAY": 30,
		"SymbolInfoDoubleProperty_SYMBOL_SWAP_WEDNESDAY": 31,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_ACCRUED_INTEREST": 19,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_CONTRACT_SIZE": 18,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_FACE_VALUE": 20,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_LIQUIDITY_RATE": 21,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_SIZE": 17,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE": 14,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE_LOSS": 16,
		"SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE_PROFIT": 15,
		"SymbolInfoDoubleProperty_SYMBOL_VOLUMEHIGH_REAL": 10,
		"SymbolInfoDoubleProperty_SYMBOL_VOLUMELOW_REAL": 11,
		"SymbolInfoDoubleProperty_SYMBOL_VOLUME_LIMIT": 25,
		"SymbolInfoDoubleProperty_SYMBOL_VOLUME_MAX": 23,
		"SymbolInfoDoubleProperty_SYMBOL_VOLUME_MIN": 22,
		"SymbolInfoDoubleProperty_SYMBOL_VOLUME_REAL": 9,
		"SymbolInfoDoubleProperty_SYMBOL_VOLUME_STEP": 24,
	})

	pi.registerEnum("SymbolInfoIntegerProperty", map[string]int32{
		"SymbolInfoIntegerProperty_SYMBOL_BACKGROUND_COLOR": 4,
		"SymbolInfoIntegerProperty_SYMBOL_CHART_MODE": 5,
		"SymbolInfoIntegerProperty_SYMBOL_CUSTOM": 3,
		"SymbolInfoIntegerProperty_SYMBOL_DIGITS": 17,
		"SymbolInfoIntegerProperty_SYMBOL_EXIST": 6,
		"SymbolInfoIntegerProperty_SYMBOL_EXPIRATION_MODE": 31,
		"SymbolInfoIntegerProperty_SYMBOL_EXPIRATION_TIME": 24,
		"SymbolInfoIntegerProperty_SYMBOL_FILLING_MODE": 32,
		"SymbolInfoIntegerProperty_SYMBOL_INDUSTRY": 2,
		"SymbolInfoIntegerProperty_SYMBOL_MARGIN_HEDGED_USE_LEG": 30,
		"SymbolInfoIntegerProperty_SYMBOL_OPTION_MODE": 35,
		"SymbolInfoIntegerProperty_SYMBOL_OPTION_RIGHT": 36,
		"SymbolInfoIntegerProperty_SYMBOL_ORDER_GTC_MODE": 34,
		"SymbolInfoIntegerProperty_SYMBOL_ORDER_MODE": 33,
		"SymbolInfoIntegerProperty_SYMBOL_SECTOR": 1,
		"SymbolInfoIntegerProperty_SYMBOL_SELECT": 7,
		"SymbolInfoIntegerProperty_SYMBOL_SESSION_BUY_ORDERS": 10,
		"SymbolInfoIntegerProperty_SYMBOL_SESSION_DEALS": 9,
		"SymbolInfoIntegerProperty_SYMBOL_SESSION_SELL_ORDERS": 11,
		"SymbolInfoIntegerProperty_SYMBOL_SPREAD": 19,
		"SymbolInfoIntegerProperty_SYMBOL_SPREAD_FLOAT": 18,
		"SymbolInfoIntegerProperty_SYMBOL_START_TIME": 23,
		"SymbolInfoIntegerProperty_SYMBOL_SUBSCRIPTION_DELAY": 0,
		"SymbolInfoIntegerProperty_SYMBOL_SWAP_MODE": 28,
		"SymbolInfoIntegerProperty_SYMBOL_SWAP_ROLLOVER3DAYS": 29,
		"SymbolInfoIntegerProperty_SYMBOL_TICKS_BOOKDEPTH": 20,
		"SymbolInfoIntegerProperty_SYMBOL_TIME": 15,
		"SymbolInfoIntegerProperty_SYMBOL_TIME_MSC": 16,
		"SymbolInfoIntegerProperty_SYMBOL_TRADE_CALC_MODE": 21,
		"SymbolInfoIntegerProperty_SYMBOL_TRADE_EXEMODE": 27,
		"SymbolInfoIntegerProperty_SYMBOL_TRADE_FREEZE_LEVEL": 26,
		"SymbolInfoIntegerProperty_SYMBOL_TRADE_MODE": 22,
		"SymbolInfoIntegerProperty_SYMBOL_TRADE_STOPS_LEVEL": 25,
		"SymbolInfoIntegerProperty_SYMBOL_VISIBLE": 8,
		"SymbolInfoIntegerProperty_SYMBOL_VOLUME": 12,
		"SymbolInfoIntegerProperty_SYMBOL_VOLUMEHIGH": 13,
		"SymbolInfoIntegerProperty_SYMBOL_VOLUMELOW": 14,
	})

	pi.registerEnum("SymbolInfoStringProperty", map[string]int32{
		"SymbolInfoStringProperty_SYMBOL_BANK": 8,
		"SymbolInfoStringProperty_SYMBOL_BASIS": 0,
		"SymbolInfoStringProperty_SYMBOL_CATEGORY": 1,
		"SymbolInfoStringProperty_SYMBOL_COUNTRY": 2,
		"SymbolInfoStringProperty_SYMBOL_CURRENCY_BASE": 5,
		"SymbolInfoStringProperty_SYMBOL_CURRENCY_MARGIN": 7,
		"SymbolInfoStringProperty_SYMBOL_CURRENCY_PROFIT": 6,
		"SymbolInfoStringProperty_SYMBOL_DESCRIPTION": 9,
		"SymbolInfoStringProperty_SYMBOL_EXCHANGE": 10,
		"SymbolInfoStringProperty_SYMBOL_FORMULA": 11,
		"SymbolInfoStringProperty_SYMBOL_INDUSTRY_NAME": 4,
		"SymbolInfoStringProperty_SYMBOL_ISIN": 12,
		"SymbolInfoStringProperty_SYMBOL_PAGE": 13,
		"SymbolInfoStringProperty_SYMBOL_PATH": 14,
		"SymbolInfoStringProperty_SYMBOL_SECTOR_NAME": 3,
	})

	pi.registerEnum("TMT5_ENUM_ORDER_TYPE", map[string]int32{
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY": 0,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT": 2,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP": 4,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP_LIMIT": 6,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_CLOSE_BY": 8,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL": 1,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_LIMIT": 3,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP": 5,
		"TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP_LIMIT": 7,
	})

	pi.registerEnum("TMT5_ENUM_ORDER_TYPE_TIME", map[string]int32{
		"TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_DAY": 1,
		"TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_GTC": 0,
		"TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_SPECIFIED": 2,
		"TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_SPECIFIED_DAY": 3,
	})

	pi.registerEnum("TerminalType", map[string]int32{
		"TerminalType_MT4": 0,
		"TerminalType_MT5": 1,
	})

}

func (pi *ProtobufInspector) registerType(t reflect.Type) {
	// Remove pointer if present
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	typeName := t.Name()
	pi.types[typeName] = t
}

func (pi *ProtobufInspector) registerEnum(name string, values map[string]int32) {
	var enumVals []EnumValue
	for k, v := range values {
		enumVals = append(enumVals, EnumValue{Name: k, Value: v})
	}

	// Sort by value
	sort.Slice(enumVals, func(i, j int) bool {
		return enumVals[i].Value < enumVals[j].Value
	})

	pi.enums[name] = enumVals
}

// Run starts the interactive inspector loop
func (pi *ProtobufInspector) Run() {
	pi.printHeader()
	pi.printQuickStart()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 2)
		command := strings.ToLower(parts[0])
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}

		switch command {
		case "exit", "quit", "q":
			fmt.Println("\n👋 Goodbye!")
			return

		case "help", "?":
			pi.printHelp()

		case "list", "ls":
			pi.listAllTypes()

		case "search", "find":
			if arg == "" {
				fmt.Println("❌ Usage: search <text>")
			} else {
				pi.searchTypes(arg)
			}

		case "field":
			if arg == "" {
				fmt.Println("❌ Usage: field <fieldname>")
			} else {
				pi.findField(arg)
			}

		case "enum":
			if arg == "" {
				fmt.Println("❌ Usage: enum <enumname>")
			} else {
				pi.inspectEnum(arg)
			}

		default:
			// Assume it's a type name
			pi.inspectType(input)
		}
	}
}

// ═════════════════════════════════════════════════════════════════
// COMMAND IMPLEMENTATIONS
// ═════════════════════════════════════════════════════════════════

func (pi *ProtobufInspector) listAllTypes() {
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║              ALL PROTOBUF TYPES                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	var names []string
	for name := range pi.types {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Printf("\n📦 Found %d types:\n\n", len(names))

	for _, name := range names {
		category := "[Type]   "
		if strings.HasSuffix(name, "Request") {
			category = "[Request]"
		} else if strings.HasSuffix(name, "Reply") {
			category = "[Reply]  "
		} else if strings.HasSuffix(name, "Info") {
			category = "[Info]   "
		}
		fmt.Printf("  %s %s\n", category, name)
	}

	fmt.Println("\n💡 Type a name to inspect it (e.g., 'PositionInfo')")
}

func (pi *ProtobufInspector) searchTypes(searchTerm string) {
	fmt.Printf("\n🔍 Searching for types containing '%s'...\n\n", searchTerm)

	var found []string
	searchLower := strings.ToLower(searchTerm)

	for name := range pi.types {
		if strings.Contains(strings.ToLower(name), searchLower) {
			found = append(found, name)
		}
	}

	sort.Strings(found)

	if len(found) == 0 {
		fmt.Printf("❌ No types found containing '%s'\n", searchTerm)
		return
	}

	fmt.Printf("✅ Found %d type(s):\n\n", len(found))

	for _, name := range found {
		t := pi.types[name]
		kind := "[Struct]"
		if t.Kind() == reflect.Interface {
			kind = "[Interface]"
		}
		fmt.Printf("  %-10s %s\n", kind, name)
	}

	fmt.Println("\n💡 Type a name to inspect it")
}

func (pi *ProtobufInspector) findField(fieldName string) {
	fmt.Printf("\n🔍 Searching for field '%s'...\n\n", fieldName)

	type FieldMatch struct {
		TypeName  string
		FieldName string
		FieldType string
	}

	var matches []FieldMatch
	fieldLower := strings.ToLower(fieldName)

	for typeName, t := range pi.types {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if strings.Contains(strings.ToLower(field.Name), fieldLower) {
				matches = append(matches, FieldMatch{
					TypeName:  typeName,
					FieldName: field.Name,
					FieldType: pi.getTypeName(field.Type),
				})
			}
		}
	}

	if len(matches) == 0 {
		fmt.Printf("❌ No types found with field containing '%s'\n", fieldName)
		return
	}

	// Group by type
	grouped := make(map[string][]FieldMatch)
	for _, m := range matches {
		grouped[m.TypeName] = append(grouped[m.TypeName], m)
	}

	var typeNames []string
	for typeName := range grouped {
		typeNames = append(typeNames, typeName)
	}
	sort.Strings(typeNames)

	fmt.Printf("✅ Found in %d type(s):\n\n", len(grouped))

	for _, typeName := range typeNames {
		fmt.Printf("📦 %s:\n", typeName)
		for _, match := range grouped[typeName] {
			fmt.Printf("   └─ %s: %s\n", match.FieldName, match.FieldType)
		}
		fmt.Println()
	}
}

func (pi *ProtobufInspector) inspectEnum(enumName string) {
	// Find enum case-insensitively
	var found string
	for name := range pi.enums {
		if strings.EqualFold(name, enumName) {
			found = name
			break
		}
	}

	if found == "" {
		fmt.Printf("\n❌ Enum '%s' not found\n", enumName)
		fmt.Println("💡 Try: enum BMT5_ENUM_ORDER_TYPE")
		return
	}

	values := pi.enums[found]

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ ENUM: %-51s ║\n", found)
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	for _, val := range values {
		fmt.Printf("  %-50s = %d\n", val.Name, val.Value)
	}
}

func (pi *ProtobufInspector) inspectType(typeName string) {
	// Try exact match first
	t, exists := pi.types[typeName]

	// Try case-insensitive search
	if !exists {
		for name, typ := range pi.types {
			if strings.EqualFold(name, typeName) {
				t = typ
				typeName = name
				exists = true
				break
			}
		}
	}

	if !exists {
		fmt.Printf("\n❌ Type '%s' not found\n", typeName)
		fmt.Printf("💡 Try: 'search %s' to find similar types\n", typeName)
		fmt.Println("💡 Or:  'list' to see all available types")
		return
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ TYPE: %-51s ║\n", typeName)
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	if t.NumField() == 0 {
		fmt.Println("  (no public fields)")
		return
	}

	fmt.Printf("📋 Fields (%d):\n\n", t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		typeName := pi.getTypeName(field.Type)
		isRepeated := field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array

		prefix := "  •"
		if isRepeated {
			prefix = "  📚"
		}

		fmt.Printf("%s %-40s : %s\n", prefix, field.Name, typeName)

		// Show protobuf tag if present
		if tag := field.Tag.Get("protobuf"); tag != "" {
			// Parse tag to get field number
			parts := strings.Split(tag, ",")
			if len(parts) > 1 {
				// parts[0] contains type (varint, bytes, fixed32, etc.)
				// parts[1] contains field number
				fieldNum := parts[1]
				fmt.Printf("       (protobuf field #%s)\n", fieldNum)
			}
		}
	}

	fmt.Println("\n💡 To see field values of an enum, use: enum <EnumName>")
}

// ═════════════════════════════════════════════════════════════════
// HELPER METHODS
// ═════════════════════════════════════════════════════════════════

func (pi *ProtobufInspector) getTypeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return "*" + pi.getTypeName(t.Elem())
	case reflect.Slice:
		return "[]" + pi.getTypeName(t.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), pi.getTypeName(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", pi.getTypeName(t.Key()), pi.getTypeName(t.Elem()))
	default:
		// Return simple type name for primitives
		name := t.Name()
		if name == "" {
			name = t.String()
		}
		return name
	}
}

// ═════════════════════════════════════════════════════════════════
// UI HELPERS
// ═════════════════════════════════════════════════════════════════

func (pi *ProtobufInspector) printHeader() {
	fmt.Print("\033[H\033[2J") // Clear screen (ANSI escape code)
	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                    ║")
	fmt.Println("║          🔍 INTERACTIVE PROTOBUF TYPES INSPECTOR 🔍                ║")
	fmt.Println("║                                                                    ║")
	fmt.Println("║         Explore MT5 gRPC API Types, Fields & Enums                 ║")
	fmt.Println("║                                                                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println("\nPackage: github.com/MetaRPC/GoMT5/package")
}

func (pi *ProtobufInspector) printQuickStart() {
	fmt.Println("\n┌────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ QUICK START GUIDE                                              │")
	fmt.Println("├────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ • Type 'help' for full command list                            │")
	fmt.Println("│ • Type 'list' to see all available types                       │")
	fmt.Println("│ • Type a type name to inspect it (e.g., 'PositionInfo')        │")
	fmt.Println("│ • Type 'exit' to quit                                          │")
	fmt.Println("└────────────────────────────────────────────────────────────────┘")
}

func (pi *ProtobufInspector) printHelp() {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        COMMAND REFERENCE                           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("📋 BROWSING COMMANDS:")
	fmt.Println("  list, ls              - List all available protobuf types")
	fmt.Println("  <TypeName>            - Inspect specific type (e.g., PositionInfo)")
	fmt.Println()
	fmt.Println("🔍 SEARCH COMMANDS:")
	fmt.Println("  search <text>         - Find types containing text")
	fmt.Println("  find <text>           - (alias for search)")
	fmt.Println("  field <name>          - Find types with specific field")
	fmt.Println("  enum <name>           - Show enum values")
	fmt.Println()
	fmt.Println("ℹ️  UTILITY COMMANDS:")
	fmt.Println("  help, ?               - Show this help")
	fmt.Println("  exit, quit, q         - Exit inspector")
	fmt.Println()
	fmt.Println("💡 EXAMPLES:")
	fmt.Println("  > list                          # See all types")
	fmt.Println("  > PositionInfo                  # Inspect PositionInfo")
	fmt.Println("  > search Order                  # Find types with 'Order'")
	fmt.Println("  > field Ticket                  # Find types with Ticket field")
	fmt.Println("  > enum BMT5_ENUM_ORDER_TYPE     # Show order type values")
	fmt.Println()
	fmt.Println("🎯 COMMON USE CASES:")
	fmt.Println("  • 'field not found' error  → Use: field <fieldname>")
	fmt.Println("  • Need enum values         → Use: enum <EnumName>")
	fmt.Println("  • Explore available types  → Use: list")
	fmt.Println("  • Find related types       → Use: search <keyword>")
}

// RunProtobufInspector starts the interactive protobuf inspector
func RunProtobufInspector() {
	inspector := NewProtobufInspector()
	inspector.Run()
}
