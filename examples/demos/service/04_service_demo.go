/*══════════════════════════════════════════════════════════════════════════════
 FILE: 04_service_demo.go - MID-LEVEL SERVICE API DEMO

 PURPOSE:
   Demonstrates MT5Service mid-level API wrapper for cleaner, more idiomatic Go code.
   Shows how to use MT5Service methods that return native Go types instead of protobuf.
   This is the RECOMMENDED layer for most trading applications.


 📚 WHAT THIS DEMO COVERS (30 Methods in 9 Steps):

   STEP 1: CREATE MT5SERVICE & CONNECT
      • NewMT5Service() - Create service wrapper
      • ConnectEx() - Connect to MT5 cluster

   STEP 2: ACCOUNT INFORMATION (4 methods)
      • GetAccountSummary() - All account data in one struct
      • GetAccountDouble() - Individual double properties
      • GetAccountInteger() - Individual integer properties
      • GetAccountString() - Individual string properties

   STEP 3: SYMBOL INFORMATION (13 methods)
      • GetSymbolTick() - Current market prices
      • GetSymbolParamsMany() - Multiple symbols at once
      • GetSymbolsTotal() - Count available symbols
      • GetSymbolInteger() - Individual integer properties
      • GetSymbolString() - Individual string properties
      • SymbolExist() - Check if symbol exists (returns bool directly)
      • GetSymbolName() - Get symbol name by index
      • GetSymbolDouble() - Get double property (Bid, Ask, Point, etc.)
      • SymbolSelect() - Add/remove symbol from Market Watch
      • IsSymbolSynchronized() - Check symbol data sync status
      • GetSymbolMarginRate() - Get margin rates for order types
      • GetSymbolSessionQuote() - Get quote session times
      • GetSymbolSessionTrade() - Get trading session times

   STEP 4: POSITIONS & ORDERS (4 methods)
      • GetOpenedOrders() - Full position/order details
      • GetOpenedTickets() - Lightweight ticket numbers only
      • GetPositionsTotal() - Count of open positions
      • GetPositionsHistory() - Historical closed positions

   STEP 5: PRE-TRADE CALCULATIONS (2 methods)
      • CalculateMargin() - Required margin for order
      • CalculateProfit() - Potential profit calculation

   STEP 6: ORDER HISTORY (1 method)
      • GetOrderHistory() - Historical orders/deals with separation

   STEP 7: TRADING OPERATIONS (4 methods)
      • CheckOrder() - Pre-validate order before sending
      • PlaceOrder() - Send order to broker (structure demo)
      • ModifyOrder() - Modify order SL/TP (structure demo)
      • CloseOrder() - Close position (structure demo)

   STEP 8: MARKET DEPTH (3 methods)
      • SubscribeMarketDepth() - Subscribe to DOM updates
      • GetMarketDepth() - Get current DOM snapshot
      • UnsubscribeMarketDepth() - Unsubscribe from DOM

   STEP 9: STREAMING (Bonus - 1 method)
      • StreamTicks() - Real-time tick stream

   FINAL: DISCONNECT
      • Disconnect() - Close connection

 ⚡ API LEVELS COMPARISON:

   LOW-LEVEL (MT5Account):
      • Direct protobuf Request/Data structures
      • account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
      • Returns: *pb.AccountSummaryData
      • Manual time conversions: data.ServerTime.AsTime()

   MID-LEVEL (MT5Service) - THIS DEMO:
      • Wrapper over MT5Account
      • service.GetAccountSummary(ctx)
      • Returns: *AccountSummary (clean struct)
      • Auto-converted times: summary.ServerTime (*time.Time)

   HIGH-LEVEL (MT5Sugar) - Coming Soon:
      • Business logic and ready patterns
      • One-liner operations with smart defaults

 💡 WHY USE MT5SERVICE?
      Less code (30-80% reduction)
      Native Go types (time.Time, float64, etc)
      No manual protobuf Request building
      Direct value returns (no .GetRequestedValue())
      Better separation (positions/orders in separate slices)
      Type safety (compiler catches more errors)

 ⚠️  IMPORTANT:
   • MT5Service wraps MT5Account (low-level)
   • You still have access to account.* methods if needed

 🚀 HOW TO RUN THIS DEMO:
   cd examples/demos
   go run main.go 4          (or select [4] from interactive menu)

══════════════════════════════════════════════════════════════════════════════*/

package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	"github.com/MetaRPC/GoMT5/mt5"
	"github.com/google/uuid"
)

const (
	MAX_STREAMING_EVENTS  = 5  // Max tick events to receive
	MAX_STREAMING_SECONDS = 10 // Max streaming duration
)

// RunService04 demonstrates MT5Service mid-level API
func RunService04() error {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("MT5 SERVICE DEMO: Mid-Level API (Recommended)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// ═══════════════════════════════════════════════════════════════════════
	// LOAD CONFIGURATION
	// ═══════════════════════════════════════════════════════════════════════
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// #region CONNECT
	// ═══════════════════════════════════════════════════════════════════════
	// STEP 1: CREATE MT5SERVICE & CONNECT
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\nSTEP 1: Create MT5Service and Connect")
	fmt.Println("───────────────────────────────────────────────────────────")

	// Create low-level MT5Account first
	account, err := mt5.NewMT5Account(cfg.User, cfg.Password, cfg.GrpcServer, uuid.New())
	if err != nil {
		return fmt.Errorf("failed to create MT5Account: %w", err)
	}
	fmt.Printf("✓ MT5Account created (UUID: %s)\n", account.Id)
	defer account.Close()

	// Create MT5Service wrapper
	service := mt5.NewMT5Service(account)
	fmt.Println("✓ MT5Service created (mid-level wrapper)")

	ctx := context.Background()

	// ConnectEx - Connect to MT5 cluster
	baseSymbol := cfg.TestSymbol
	connectExReq := &pb.ConnectExRequest{
		User:            cfg.User,
		Password:        cfg.Password,
		MtClusterName:   cfg.MtCluster,
		BaseChartSymbol: &baseSymbol,
	}

	// Use context timeout for ConnectEx (replaces old TerminalReadinessWaitingTimeoutSeconds)
	connectCtx, connectCancel := context.WithTimeout(ctx, 180*time.Second)
	defer connectCancel()

	connectData, err := account.ConnectEx(connectCtx, connectExReq)
	if err != nil {
		return fmt.Errorf("ConnectEx failed: %w", err)
	}

	account.Id = uuid.MustParse(connectData.TerminalInstanceGuid)
	fmt.Printf("✓ Connected (Terminal GUID: %s)\n", connectData.TerminalInstanceGuid)

	// #endregion
	// #region ACCOUNT INFO

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 2: ACCOUNT INFORMATION
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 2: Account Information Methods")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 2.1. GET ACCOUNT SUMMARY
	//      Get complete account information in ONE call.
	//      Returns: Clean AccountSummary struct with native Go types.
	//      All times already converted to *time.Time.
	//      Shorter field names: Balance (not AccountBalance).
	//      RECOMMENDED for retrieving account data.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.1. GetAccountSummary() - Get all account data")

	summary, err := service.GetAccountSummary(ctx)
	if err != nil {
		return fmt.Errorf("GetAccountSummary failed: %w", err)
	}

	fmt.Printf("  Login:         %d\n", summary.Login)
	fmt.Printf("  User Name:     %s\n", summary.UserName)
	fmt.Printf("  Balance:       %.2f %s\n", summary.Balance, summary.Currency)
	fmt.Printf("  Equity:        %.2f\n", summary.Equity)
	fmt.Printf("  Credit:        %.2f\n", summary.Credit)
	fmt.Printf("  Leverage:      1:%d\n", summary.Leverage)
	fmt.Printf("  Company:       %s\n", summary.CompanyName)
	if summary.ServerTime != nil {
		fmt.Printf("  Server Time:   %s (UTC%+.1f)\n",
			summary.ServerTime.Format("2006-01-02 15:04:05"),
			float64(summary.UtcTimezoneShiftMinutes)/60.0)
	}
	fmt.Printf("  Trade Mode:    %s\n", summary.TradeMode)

	// ══════════════════════════════════════════════════════════════
	// 2.2. GET ACCOUNT DOUBLE
	//      Get individual account property by ID.
	//      Returns: float64 directly (no Data struct extraction).
	//      Use when you need just one specific value.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.2. GetAccountDouble() - Get individual property")

	balance, err := service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE)
	if err != nil {
		return fmt.Errorf("GetAccountDouble failed: %w", err)
	}

	fmt.Printf("  Balance:       %.2f (same as above, direct return)\n", balance)

	// ══════════════════════════════════════════════════════════════
	// 2.3. GET ACCOUNT INTEGER
	//      Get individual account integer property by ID.
	//      Returns: int64 directly (no Data struct extraction).
	//      Use for Leverage, Login, TradeMode, etc.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.3. GetAccountInteger() - Get integer property")

	leverage, err := service.GetAccountInteger(ctx, pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE)
	if err != nil {
		return fmt.Errorf("GetAccountInteger failed: %w", err)
	}

	login, err := service.GetAccountInteger(ctx, pb.AccountInfoIntegerPropertyType_ACCOUNT_LOGIN)
	if err != nil {
		return fmt.Errorf("GetAccountInteger failed: %w", err)
	}

	fmt.Printf("  Leverage:      1:%d (direct int64 return)\n", leverage)
	fmt.Printf("  Login:         %d\n", login)

	// ══════════════════════════════════════════════════════════════
	// 2.4. GET ACCOUNT STRING
	//      Get individual account string property by ID.
	//      Returns: string directly (no Data struct extraction).
	//      Use for Currency, Company, Name, etc.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.4. GetAccountString() - Get string property")

	currency, err := service.GetAccountString(ctx, pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY)
	if err != nil {
		return fmt.Errorf("GetAccountString failed: %w", err)
	}

	company, err := service.GetAccountString(ctx, pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY)
	if err != nil {
		return fmt.Errorf("GetAccountString failed: %w", err)
	}

	fmt.Printf("  Currency:      %s (direct string return)\n", currency)
	fmt.Printf("  Company:       %s\n", company)

	// #endregion
    // #region SYMBOL INFO

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 3: SYMBOL INFORMATION
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 3: Symbol Information Methods")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 3.1. GET SYMBOL TICK
	//      Get current market prices for a symbol.
	//      Returns: Clean SymbolTick struct.
	//      Time field is time.Time (not Unix timestamp).
	//      Easy access to Bid, Ask, Spread.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.1. GetSymbolTick() - Get current market prices")

	tick, err := service.GetSymbolTick(ctx, cfg.TestSymbol)
	if err != nil {
		return fmt.Errorf("GetSymbolTick failed: %w", err)
	}

	spread := tick.Ask - tick.Bid
	fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
	fmt.Printf("  Bid:           %.5f\n", tick.Bid)
	fmt.Printf("  Ask:           %.5f\n", tick.Ask)
	fmt.Printf("  Spread:        %.5f (%.1f points)\n", spread, spread*10000)
	fmt.Printf("  Last:          %.5f\n", tick.Last)
	fmt.Printf("  Volume:        %d\n", tick.Volume)
	fmt.Printf("  Time:          %s\n", tick.Time.Format("2006-01-02 15:04:05"))

	// ══════════════════════════════════════════════════════════════
	// 3.2. GET SYMBOL PARAMS MANY
	//      Get multiple symbols' parameters at once.
	//      Supports pagination (page, perPage).
	//      Returns: Slice of clean SymbolParam structs.
	//      Efficient batch operation.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.2. GetSymbolParamsMany() - Get multiple symbols (showing first 5 alphabetically)")

	page := int32(1)
	perPage := int32(5)
	// Note: nil filter returns ALL symbols (2000+), sorted alphabetically
	// First symbols are often CFD indices (#AUS200, #ChinaA50) which may show zero prices when markets are closed
	symbols, total, err := service.GetSymbolParamsMany(ctx, nil, nil, &page, &perPage)
	if err != nil {
		return fmt.Errorf("GetSymbolParamsMany failed: %w", err)
	}

	fmt.Printf("  Retrieved %d symbols (total available: %d):\n", len(symbols), total)
	for i, sym := range symbols {
		if i >= 3 {
			fmt.Printf("  ... and %d more\n", len(symbols)-3)
			break
		}
		fmt.Printf("    %d. %-10s  Bid: %.5f  Ask: %.5f  Digits: %d  Spread: %d pts\n",
			i+1, sym.Name, sym.Bid, sym.Ask, sym.Digits, sym.Spread)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.3. GET SYMBOLS TOTAL
	//      Get count of available symbols.
	//      Returns: int32 directly (no Data struct).
	//      selectedOnly=true for Market Watch count only.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.3. GetSymbolsTotal() - Count available symbols")

	allSymbolsCount, err := service.GetSymbolsTotal(ctx, false)
	if err != nil {
		return fmt.Errorf("GetSymbolsTotal failed: %w", err)
	}

	watchSymbolsCount, err := service.GetSymbolsTotal(ctx, true)
	if err != nil {
		return fmt.Errorf("GetSymbolsTotal failed: %w", err)
	}

	fmt.Printf("  All Symbols:       %d (total in terminal)\n", allSymbolsCount)
	fmt.Printf("  Market Watch Only: %d (selected symbols)\n", watchSymbolsCount)

	// ══════════════════════════════════════════════════════════════
	// 3.4. GET SYMBOL INTEGER
	//      Get individual symbol integer property.
	//      Returns: int64 directly (no Data struct).
	//      Use for Digits, Spread, TradeMode, etc.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.4. GetSymbolInteger() - Get integer property")

	digits, err := service.GetSymbolInteger(ctx, cfg.TestSymbol, pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS)
	if err != nil {
		return fmt.Errorf("GetSymbolInteger failed: %w", err)
	}

	spreadPoints, err := service.GetSymbolInteger(ctx, cfg.TestSymbol, pb.SymbolInfoIntegerProperty_SYMBOL_SPREAD)
	if err != nil {
		return fmt.Errorf("GetSymbolInteger failed: %w", err)
	}

	fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
	fmt.Printf("  Digits:        %d (decimal places)\n", digits)
	fmt.Printf("  Spread:        %d points\n", spreadPoints)

	// ══════════════════════════════════════════════════════════════
	// 3.5. GET SYMBOL STRING
	//      Get individual symbol string property.
	//      Returns: string directly (no Data struct).
	//      Use for Description, Path, Currency, etc.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.5. GetSymbolString() - Get string property")

	description, err := service.GetSymbolString(ctx, cfg.TestSymbol, pb.SymbolInfoStringProperty_SYMBOL_DESCRIPTION)
	if err != nil {
		return fmt.Errorf("GetSymbolString failed: %w", err)
	}

	baseCurrency, err := service.GetSymbolString(ctx, cfg.TestSymbol, pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_BASE)
	if err != nil {
		return fmt.Errorf("GetSymbolString failed: %w", err)
	}

	profitCurrency, err := service.GetSymbolString(ctx, cfg.TestSymbol, pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_PROFIT)
	if err != nil {
		return fmt.Errorf("GetSymbolString failed: %w", err)
	}

	fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
	fmt.Printf("  Description:   %s\n", description)
	fmt.Printf("  Base Currency: %s\n", baseCurrency)
	fmt.Printf("  Profit Currency: %s\n", profitCurrency)

	// ══════════════════════════════════════════════════════════════
	// 3.6. SYMBOL EXIST
	//      Check if symbol exists in terminal.
	//      Returns: bool directly (no Data struct).
	//      Fast check before working with symbol.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.6. SymbolExist() - Check if symbol exists")

	existsEUR, selectedEUR, err := service.SymbolExist(ctx, "EURUSD")
	if err != nil {
		return fmt.Errorf("SymbolExist failed: %w", err)
	}

	existsFake, selectedFake, err := service.SymbolExist(ctx, "FAKESYMBOL123")
	if err != nil {
		return fmt.Errorf("SymbolExist failed: %w", err)
	}

	fmt.Printf("  EURUSD exists:      %t (selected: %t)\n", existsEUR, selectedEUR)
	fmt.Printf("  FAKESYMBOL123:      %t (selected: %t)\n", existsFake, selectedFake)

	// ══════════════════════════════════════════════════════════════
	// 3.7. GET SYMBOL NAME
	//      Get symbol name by index position.
	//      Returns: string directly (no Data struct).
	//      Useful for iterating all symbols.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.7. GetSymbolName() - Get symbol name by index")

	symbolName, err := service.GetSymbolName(ctx, 0, false) // First symbol, all symbols
	if !helpers.PrintShortError(err, "GetSymbolName failed") {
		fmt.Printf("  First symbol (index 0): %s\n", symbolName)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.8. GET SYMBOL DOUBLE
	//      Get individual symbol double property.
	//      Returns: float64 directly (no Data struct).
	//      Use for Bid, Ask, Point, etc.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.8. GetSymbolDouble() - Get double property")

	bidPrice, err := service.GetSymbolDouble(ctx, cfg.TestSymbol, pb.SymbolInfoDoubleProperty_SYMBOL_BID)
	if err != nil {
		return fmt.Errorf("GetSymbolDouble failed: %w", err)
	}

	askPrice, err := service.GetSymbolDouble(ctx, cfg.TestSymbol, pb.SymbolInfoDoubleProperty_SYMBOL_ASK)
	if err != nil {
		return fmt.Errorf("GetSymbolDouble failed: %w", err)
	}

	fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
	fmt.Printf("  Bid:           %.5f (direct double return)\n", bidPrice)
	fmt.Printf("  Ask:           %.5f\n", askPrice)
	fmt.Printf("  Spread:        %.5f points\n", askPrice-bidPrice)

	// ══════════════════════════════════════════════════════════════
	// 3.9. SYMBOL SELECT
	//      Add/remove symbol from Market Watch.
	//      Returns: bool directly (success/failure).
	//      Useful for managing visible symbols.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.9. SymbolSelect() - Add/remove from Market Watch")

	// Ensure test symbol is in Market Watch
	selected, err := service.SymbolSelect(ctx, cfg.TestSymbol, true)
	if !helpers.PrintShortError(err, "SymbolSelect failed") {
		fmt.Printf("  Symbol %s in Market Watch: %t\n", cfg.TestSymbol, selected)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.10. IS SYMBOL SYNCHRONIZED
	//       Check if symbol data is synchronized.
	//       Returns: bool directly (no Data struct).
	//       Ensures symbol is ready for trading.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.10. IsSymbolSynchronized() - Check sync status")

	isSynced, err := service.IsSymbolSynchronized(ctx, cfg.TestSymbol)
	if !helpers.PrintShortError(err, "IsSymbolSynchronized failed") {
		fmt.Printf("  Symbol %s synchronized: %t\n", cfg.TestSymbol, isSynced)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.11. GET SYMBOL MARGIN RATE
	//       Get margin rates for order types.
	//       Returns: MarginRate struct with direct values.
	//       Shows margin requirements per order type.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.11. GetSymbolMarginRate() - Get margin rates")

	marginRate, err := service.GetSymbolMarginRate(ctx, cfg.TestSymbol, pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY)
	if !helpers.PrintShortError(err, "GetSymbolMarginRate failed") {
		fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
		fmt.Printf("  Order Type:    BUY\n")
		fmt.Printf("  Initial:       %.2f%%\n", marginRate.InitialMarginRate*100)
		fmt.Printf("  Maintenance:   %.2f%%\n", marginRate.MaintenanceMarginRate*100)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.12. GET SYMBOL SESSION QUOTE
	//       Get quote session times.
	//       Returns: SessionInfo with start/end times.
	//       Shows when quotes are available.
	//       Note: For forex, 00:00-00:00 means 24-hour session (normal)
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.12. GetSymbolSessionQuote() - Get quote session")

	quoteSession, err := service.GetSymbolSessionQuote(ctx, cfg.TestSymbol, pb.DayOfWeek_MONDAY, 0)
	if !helpers.PrintShortError(err, "GetSymbolSessionQuote failed") {
		fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
		fmt.Printf("  Day:           Monday, Session 0\n")
		fmt.Printf("  From:          %02d:%02d\n", quoteSession.From.Hour(), quoteSession.From.Minute())
		fmt.Printf("  To:            %02d:%02d\n", quoteSession.To.Hour(), quoteSession.To.Minute())
		fmt.Printf("  (00:00-00:00 = 24-hour forex trading)\n")
	}

	// ══════════════════════════════════════════════════════════════
	// 3.13. GET SYMBOL SESSION TRADE
	//       Get trading session times.
	//       Returns: SessionInfo with start/end times.
	//       Shows when trading is allowed.
	//       Note: For forex, 00:00-00:00 means 24-hour session (normal)
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.13. GetSymbolSessionTrade() - Get trade session")

	tradeSession, err := service.GetSymbolSessionTrade(ctx, cfg.TestSymbol, pb.DayOfWeek_MONDAY, 0)
	if !helpers.PrintShortError(err, "GetSymbolSessionTrade failed") {
		fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
		fmt.Printf("  Day:           Monday, Session 0\n")
		fmt.Printf("  From:          %02d:%02d\n", tradeSession.From.Hour(), tradeSession.From.Minute())
		fmt.Printf("  To:            %02d:%02d\n", tradeSession.To.Hour(), tradeSession.To.Minute())
		fmt.Printf("  (00:00-00:00 = 24-hour forex trading)\n")
	}

	// #endregion
	// #region POSITIONS

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 4: POSITIONS & ORDERS
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 4: Positions & Orders Methods")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 4.1. GET OPENED ORDERS
	//      Get all open positions and pending orders.
	//      Returns: TWO separate slices (positions, orders).
	//      All time fields are time.Time (auto-converted).
	//      Clean Position and PendingOrder structs.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.1. GetOpenedOrders() - Get positions and pending orders")

	openedData, err := service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if err != nil {
		return fmt.Errorf("GetOpenedOrders failed: %w", err)
	}

	fmt.Printf("  Open Positions: %d\n", len(openedData.PositionInfos))
	if len(openedData.PositionInfos) > 0 {
		for i, pos := range openedData.PositionInfos {
			if i >= 3 {
				fmt.Printf("    ... and %d more positions\n", len(openedData.PositionInfos)-3)
				break
			}
			fmt.Printf("    Position #%d: Ticket=%d Symbol=%s Type=%s Volume=%.2f Profit=%.2f OpenTime=%s\n",
				i+1, pos.Ticket, pos.Symbol, pos.Type, pos.Volume, pos.Profit,
				pos.OpenTime.AsTime().Format("2006-01-02 15:04"))
		}
	} else {
		fmt.Println("    (No open positions)")
	}

	fmt.Printf("\n  Pending Orders: %d\n", len(openedData.OpenedOrders))
	if len(openedData.OpenedOrders) > 0 {
		for i, ord := range openedData.OpenedOrders {
			if i >= 3 {
				fmt.Printf("    ... and %d more orders\n", len(openedData.OpenedOrders)-3)
				break
			}
			fmt.Printf("    Order #%d: Ticket=%d Symbol=%s Type=%s Volume=%.2f Price=%.5f SetupTime=%s\n",
				i+1, ord.Ticket, ord.Symbol, ord.Type, ord.VolumeInitial, ord.PriceOpen,
				ord.TimeSetup.AsTime().Format("2006-01-02 15:04"))
		}
	} else {
		fmt.Println("    (No pending orders)")
	}

	// ══════════════════════════════════════════════════════════════
	// 4.2. GET OPENED TICKETS
	//      Get only ticket numbers (lightweight).
	//      Returns: Two slices of ticket IDs (int64).
	//      Much faster when you only need ticket numbers.
	//      Use for quick "what's open?" checks.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.2. GetOpenedTickets() - Get ticket numbers only (lightweight)")

	posTickets, orderTickets, err := service.GetOpenedTickets(ctx)
	if err != nil {
		return fmt.Errorf("GetOpenedTickets failed: %w", err)
	}

	fmt.Printf("  Position Tickets: %v\n", posTickets)
	fmt.Printf("  Order Tickets:    %v\n", orderTickets)

	// ══════════════════════════════════════════════════════════════
	// 4.3. GET POSITIONS TOTAL
	//      Get total number of open positions.
	//      Returns: int32 directly (no Data struct).
	//      Fast count without loading full position data.
	//      Note: 0 is normal if account has no open positions
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.3. GetPositionsTotal() - Get count of open positions")

	totalPositions, err := service.GetPositionsTotal(ctx)
	if !helpers.PrintShortError(err, "GetPositionsTotal failed") {
		fmt.Printf("  Total Open Positions: %d", totalPositions)
		if totalPositions == 0 {
			fmt.Printf(" (no positions currently open)")
		}
		fmt.Println()
	}

	// ══════════════════════════════════════════════════════════════
	// 4.4. GET POSITIONS HISTORY
	//      Get historical positions (closed positions).
	//      Returns: []HistoryPosition with all details.
	//      Supports pagination and time filtering.
	//      All time fields are time.Time (auto-converted).
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.4. GetPositionsHistory() - Get closed positions")

	// Get last 24 hours of history (first 10 positions)
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	pageNum := int32(0)
	pageSize := int32(10)

	posHistData, err := service.GetPositionsHistory(
		ctx,
		pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_DESC,
		&yesterday,
		&now,
		&pageNum,
		&pageSize,
	)
	if !helpers.PrintShortError(err, "GetPositionsHistory failed") {
		fmt.Printf("  History Positions (last 24h): %d\n", len(posHistData.HistoryPositions))
		if len(posHistData.HistoryPositions) > 0 {
			for i, hp := range posHistData.HistoryPositions {
				if i >= 3 {
					fmt.Printf("    ... and %d more positions\n", len(posHistData.HistoryPositions)-3)
					break
				}
				fmt.Printf("    Position #%d: Ticket=%d Symbol=%s Volume=%.2f Profit=%.2f Duration=%s\n",
					i+1, hp.PositionTicket, hp.Symbol, hp.Volume, hp.Profit,
					hp.CloseTime.AsTime().Sub(hp.OpenTime.AsTime()).Round(time.Minute))
			}
		} else {
			fmt.Println("    (No closed positions in last 24 hours)")
		}
	}

    // #endregion
    // #region CALCULATIONS

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 5: PRE-TRADE CALCULATIONS
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 5: Pre-Trade Calculation Methods")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 5.1. CALCULATE MARGIN
	//      Calculate required margin for an order.
	//      Returns: float64 directly (no Data struct).
	//      Use before placing orders to check margin requirements.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n5.1. CalculateMargin() - Required margin for order")

	marginReq := &pb.OrderCalcMarginRequest{
		OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Symbol:    cfg.TestSymbol,
		Volume:    cfg.TestVolume,
		OpenPrice: tick.Ask,
	}

	margin, err := service.CalculateMargin(ctx, marginReq)
	if !helpers.PrintShortError(err, "CalculateMargin failed") {
		fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
		fmt.Printf("  Volume:        %.2f lots\n", cfg.TestVolume)
		fmt.Printf("  Required Margin: %.2f %s\n", margin, summary.Currency)
	}

	// ══════════════════════════════════════════════════════════════
	// 5.2. CALCULATE PROFIT
	//      Calculate potential profit for a trade.
	//      Returns: float64 directly.
	//      Use to estimate P&L before trading.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n5.2. CalculateProfit() - Potential profit calculation")

	// Calculate profit for +10 pips target
	pipSize := 0.0001 // For EURUSD
	targetPrice := tick.Ask + (10 * pipSize)

	profitReq := &pb.OrderCalcProfitRequest{
		OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Symbol:     cfg.TestSymbol,
		Volume:     cfg.TestVolume,
		OpenPrice:  tick.Ask,
		ClosePrice: targetPrice,
	}

	profit, err := service.CalculateProfit(ctx, profitReq)
	if !helpers.PrintShortError(err, "CalculateProfit failed") {
		fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
		fmt.Printf("  Volume:        %.2f lots\n", cfg.TestVolume)
		fmt.Printf("  Entry:         %.5f (Ask)\n", tick.Ask)
		fmt.Printf("  Target:        %.5f (+10 pips)\n", targetPrice)
		fmt.Printf("  Potential Profit: %.2f %s\n", profit, summary.Currency)
	}

    // #endregion
	// #region ORDER HISTORY

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 6: ORDER HISTORY
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 6: Order History Methods")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 6.1. GET ORDER HISTORY
	//      Get historical orders and deals for a time period.
	//      Returns: TWO separate slices (orders, deals).
	//      ADVANTAGE: Orders and Deals separated (not mixed).
	//      All time fields are time.Time (auto-converted).
	//      Supports pagination for large histories.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n6.1. GetOrderHistory() - Historical orders & deals with separation")

	// Get last 7 days of history
	historyFrom := time.Now().AddDate(0, 0, -7)
	historyTo := time.Now()
	historyPage := int32(1)
	historyPerPage := int32(50)

	histData, err := service.GetOrderHistory(ctx,
		historyFrom, historyTo,
		pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_CLOSE_TIME_DESC,
		historyPage, historyPerPage)
	if !helpers.PrintShortError(err, "GetOrderHistory failed") {
		// Separate orders and deals
		var histOrders []*pb.OrderHistoryData
		var histDeals []*pb.DealHistoryData
		for _, item := range histData.HistoryData {
			if item.HistoryOrder != nil {
				histOrders = append(histOrders, item.HistoryOrder)
			}
			if item.HistoryDeal != nil {
				histDeals = append(histDeals, item.HistoryDeal)
			}
		}

		fmt.Printf("  Time Period:   %s to %s\n",
			historyFrom.Format("2006-01-02"), historyTo.Format("2006-01-02"))
		fmt.Printf("  Total Records: %d\n", histData.ArrayTotal)
		fmt.Printf("  Orders:        %d (pending/limit/stop orders)\n", len(histOrders))
		fmt.Printf("  Deals:         %d (actual executions with P&L)\n", len(histDeals))

		if len(histOrders) > 0 {
			fmt.Println("\n  Recent Orders (first 3):")
			for i, ord := range histOrders {
				if i >= 3 {
					break
				}
				fmt.Printf("    Order #%d: Ticket=%d Symbol=%s Type=%s Volume=%.2f Price=%.5f State=%s\n",
					i+1, ord.Ticket, ord.Symbol, ord.Type, ord.VolumeInitial, ord.PriceOpen, ord.State)
			}
		}

		if len(histDeals) > 0 {
			fmt.Println("\n  Recent Deals (first 3):")
			for i, deal := range histDeals {
				if i >= 3 {
					break
				}
				fmt.Printf("    Deal #%d: Ticket=%d Symbol=%s Type=%s Volume=%.2f Price=%.5f Profit=%.2f\n",
					i+1, deal.Ticket, deal.Symbol, deal.Type, deal.Volume, deal.Price, deal.Profit)
			}
		}

		if len(histOrders) == 0 && len(histDeals) == 0 {
			fmt.Println("  (No history in the last 7 days)")
		}
	}

    // #endregion
	// #region TRADING OPERATIONS

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 7: TRADING OPERATIONS (Advanced)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 7: Trading Operations Methods (Advanced)")
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println("NOTE: These methods show API structure.")
	fmt.Println("      CheckOrder runs validation, others show structure only.")

	// ══════════════════════════════════════════════════════════════
	// 7.1. CHECK ORDER (Pre-validation)
	//      Validate an order BEFORE sending to broker.
	//      Returns: Clean OrderCheckResult struct.
	//      ADVANTAGE: Auto-extracts nested MqlTradeCheckResult fields.
	//      Shows what account state will be after execution.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n7.1. CheckOrder() - Pre-validate order")

	sl := tick.Ask - (50 * 0.0001)  // 50 pips SL
	tp := tick.Ask + (100 * 0.0001) // 100 pips TP

	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
			Symbol:     cfg.TestSymbol,
			OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
			Volume:     cfg.TestVolume,
			Price:      tick.Ask,
			StopLoss:   sl,
			TakeProfit: tp,
		},
	}

	checkResult, err := service.CheckOrder(ctx, checkReq)
	if err != nil {
		// CheckOrder often fails on DEMO accounts (broker limitation)
		fmt.Println("  ❌ CheckOrder FAILED (expected on demo accounts)")

		// Truncate error message to avoid terminal clutter
		errMsg := err.Error()
		if len(errMsg) > 150 {
			errMsg = errMsg[:150] + "... [truncated]"
		}
		fmt.Printf("     Error: %s\n", errMsg)

		fmt.Println("\n  ℹ️  Known limitation on DEMO accounts:")
		fmt.Println("     • OrderCheck not supported by many demo brokers")
		fmt.Println("     • Use OrderCalcMargin()/OrderCalcProfit() for validation instead")
		fmt.Println("     • OrderSend() will still work despite this error")
		fmt.Println()
	} else {
		fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
		fmt.Printf("  Type:          BUY\n")
		fmt.Printf("  Volume:        %.2f lots\n", cfg.TestVolume)
		fmt.Printf("  Price:         %.5f\n", tick.Ask)
		fmt.Printf("  Return Code:   %d\n", checkResult.ReturnedCode)
		if checkResult.ReturnedCode == 0 {
			fmt.Printf("  ✓ Order Valid:\n")
			fmt.Printf("    Required Margin:  %.2f %s\n", checkResult.Margin, summary.Currency)
			fmt.Printf("    Balance After:    %.2f %s\n", checkResult.Balance, summary.Currency)
			fmt.Printf("    Equity After:     %.2f %s\n", checkResult.Equity, summary.Currency)
			fmt.Printf("    Free Margin:      %.2f %s\n", checkResult.MarginFree, summary.Currency)
			fmt.Printf("    Margin Level:     %.2f%%\n", checkResult.MarginLevel)
		} else {
			fmt.Printf("  ✗ Order Invalid: %s\n", checkResult.Comment)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 7.2. PLACE ORDER (⚠️ REAL TRADE!)
	//      Send order to broker for execution.
	//      Returns: OrderResult with deal/order tickets.
	//      ADVANTAGE: Auto-converts nested result fields.
	//      DANGEROUS: Places REAL order with minimal volume.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n7.2. PlaceOrder() - Place market BUY order")

	pricePtr := tick.Ask
	commentPtr := "GoMT5 Service Demo"

	placeReq := &pb.OrderSendRequest{
		Symbol:    cfg.TestSymbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
		Volume:    cfg.TestVolume,
		Price:     &pricePtr,
		Comment:   &commentPtr,
		// Note: SL/TP will be added via ModifyOrder (to demonstrate both methods)
	}

	orderResult, err := service.PlaceOrder(ctx, placeReq)
	if !helpers.PrintShortError(err, "PlaceOrder failed") {
		fmt.Printf("  Order sent result:\n")
		fmt.Printf("    Return Code:   %d\n", orderResult.ReturnedCode)
		fmt.Printf("    Deal Ticket:   %d\n", orderResult.Deal)
		fmt.Printf("    Order Ticket:  %d\n", orderResult.Order)
		fmt.Printf("    Volume:        %.2f\n", orderResult.Volume)
		fmt.Printf("    Execution Price: %.5f\n", orderResult.Price)
		// Bid/Ask are optional fields in OrderSend response - skip if not provided
		if orderResult.Bid > 0 && orderResult.Ask > 0 {
			fmt.Printf("    Market Bid:    %.5f\n", orderResult.Bid)
			fmt.Printf("    Market Ask:    %.5f\n", orderResult.Ask)
		}
		fmt.Printf("    Comment:       %s\n", orderResult.Comment)

		if orderResult.ReturnedCode == 10009 {
			fmt.Printf("    ✓ Order EXECUTED successfully!\n")
			orderTicket := orderResult.Order

			// ══════════════════════════════════════════════════════════════
			// 7.3. MODIFY ORDER (⚠️ MODIFIES REAL POSITION!)
			//      Modify existing order or position (change SL/TP/price).
			//      Returns: OrderResult with modification details.
			//      ADVANTAGE: Simplified result extraction.
			//      DANGEROUS: Modifies REAL open position!
			// ══════════════════════════════════════════════════════════════
			fmt.Println("\n7.3. ModifyOrder() - Add Stop Loss and Take Profit")

			// Calculate SL/TP levels (10 pips SL, 20 pips TP)
			newSL := orderResult.Price - (10 * 0.0001)
			newTP := orderResult.Price + (20 * 0.0001)

			modifyReq := &pb.OrderModifyRequest{
				Ticket:     orderTicket,
				StopLoss:   &newSL,
				TakeProfit: &newTP,
			}

			modifyResult, err := service.ModifyOrder(ctx, modifyReq)
			if !helpers.PrintShortError(err, "ModifyOrder failed") {
				fmt.Printf("  Order modify result:\n")
				fmt.Printf("    Return Code:   %d\n", modifyResult.ReturnedCode)
				fmt.Printf("    Order Ticket:  %d\n", modifyResult.Order)
				fmt.Printf("    Stop Loss:     %.5f\n", newSL)
				fmt.Printf("    Take Profit:   %.5f\n", newTP)
				fmt.Printf("    Comment:       %s\n", modifyResult.Comment)

				if modifyResult.ReturnedCode == 10009 {
					fmt.Printf("    ✓ Position MODIFIED successfully!\n")
				}
			}

			// ══════════════════════════════════════════════════════════════
			// 7.4. CLOSE ORDER (⚠️ CLOSES REAL POSITION!)
			//      Close position or delete pending order.
			//      Returns: uint32 return code directly (simplified).
			//      ADVANTAGE: Returns code directly (not nested in struct).
			//      DANGEROUS: Closes REAL position, realizes profit/loss!
			// ══════════════════════════════════════════════════════════════
			fmt.Println("\n7.4. CloseOrder() - Close the position")

			closeReq := &pb.OrderCloseRequest{
				Ticket:   orderTicket,
				Volume:   cfg.TestVolume,
				Slippage: 10,
			}

			retCode, err := service.CloseOrder(ctx, closeReq)
			if !helpers.PrintShortError(err, "CloseOrder failed") {
				fmt.Printf("  Order close result:\n")
				fmt.Printf("    Return Code:   %d\n", retCode)
				fmt.Printf("    Ticket:        %d\n", orderTicket)

				if retCode == 10009 {
					fmt.Printf("    ✓ Position CLOSED successfully!\n")
				} else {
					fmt.Printf("    ⚠️  Close returned code %d\n", retCode)
				}
			}
		} else {
			fmt.Printf("    ⚠️  Order placement returned code %d - %s\n", orderResult.ReturnedCode, orderResult.Comment)
		}
	}

    // #endregion
	// #region DOM

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 8: MARKET DEPTH (DOM)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 8: Market Depth (DOM) Methods")
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println("\nℹ️  IMPORTANT: Market Depth (Order Book) is only available for exchange-traded instruments")
	fmt.Println("   • Works for: Stocks, Futures, Options (exchange symbols)")
	fmt.Println("   • NOT available for: Forex pairs (OTC trading, no centralized order book)")
	fmt.Println("   • Testing with EURUSD will result in timeout/empty data - this is expected!")

	// ══════════════════════════════════════════════════════════════
	// 8.1. SUBSCRIBE MARKET DEPTH
	//      Subscribe to Depth of Market updates.
	//      ADVANTAGE: Takes symbol string (not Request object).
	//      Must call before GetMarketDepth.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n8.1. SubscribeMarketDepth() - Subscribe to DOM updates")

	subscribed, err := service.SubscribeMarketDepth(ctx, cfg.TestSymbol)
	if !helpers.PrintShortError(err, "SubscribeMarketDepth failed") && subscribed {
		fmt.Printf("  ✓ Subscribed to Market Depth for %s\n", cfg.TestSymbol)

		// ══════════════════════════════════════════════════════════════
		// 8.2. GET MARKET DEPTH
		//      Get current DOM snapshot.
		//      ADVANTAGE: Returns clean []BookInfo slice.
		//      Requires prior SubscribeMarketDepth call.
		// ══════════════════════════════════════════════════════════════
		fmt.Println("\n8.2. GetMarketDepth() - Get current DOM snapshot")

		domBooks, err := service.GetMarketDepth(ctx, cfg.TestSymbol)
		if !helpers.PrintShortError(err, "GetMarketDepth failed") {
			fmt.Printf("  Symbol:        %s\n", cfg.TestSymbol)
			fmt.Printf("  Depth Levels:  %d\n", len(domBooks))

			if len(domBooks) > 0 {
				// Show top 3 bids and top 3 asks
				bidCount := 0
				askCount := 0
				fmt.Println("\n  Top Bid Levels:")
				for i, book := range domBooks {
					if book.Type == pb.BookType_BOOK_TYPE_BUY {
						bidCount++
						fmt.Printf("    Bid #%d: Price=%.5f Volume=%.2f\n",
							bidCount, book.Price, book.VolumeReal)
						if bidCount >= 3 {
							break
						}
					}
					if i >= 20 { // Safety limit
						break
					}
				}

				fmt.Println("\n  Top Ask Levels:")
				for i, book := range domBooks {
					if book.Type == pb.BookType_BOOK_TYPE_SELL {
						askCount++
						fmt.Printf("    Ask #%d: Price=%.5f Volume=%.2f\n",
							askCount, book.Price, book.VolumeReal)
						if askCount >= 3 {
							break
						}
					}
					if i >= 20 { // Safety limit
						break
					}
				}

				if bidCount == 0 && askCount == 0 {
					fmt.Println("  (No market depth data available)")
				}
			} else {
				fmt.Println("  (Market depth not available for this symbol)")
			}
		}

		// ══════════════════════════════════════════════════════════════
		// 8.3. UNSUBSCRIBE MARKET DEPTH
		//      Unsubscribe from DOM updates.
		//      ADVANTAGE: Takes symbol string (not Request object).
		//      Always clean up subscriptions.
		// ══════════════════════════════════════════════════════════════
		fmt.Println("\n8.3. UnsubscribeMarketDepth() - Unsubscribe from DOM")

		unsubscribed, err := service.UnsubscribeMarketDepth(ctx, cfg.TestSymbol)
		if !helpers.PrintShortError(err, "UnsubscribeMarketDepth failed") && unsubscribed {
			fmt.Printf("  ✓ Unsubscribed from Market Depth for %s\n", cfg.TestSymbol)
		}
	} else {
		fmt.Printf("  ⚠️  Could not subscribe to Market Depth for %s\n", cfg.TestSymbol)
	}

	// #endregion
    // #region STREAMING

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 9: STREAMING (Bonus)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 9: Streaming Methods (Bonus)")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 9.1. STREAM TICKS
	//      Real-time tick data stream.
	//      Returns: Channel of *SymbolTick structs.
	//      All fields auto-converted (time.Time, etc).
	//      Cleaner than low-level OnSymbolTick.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n9.1. StreamTicks() - Real-time tick stream")

	streamCtx, streamCancel := context.WithTimeout(ctx, MAX_STREAMING_SECONDS*time.Second)
	defer streamCancel()

	tickCh, errCh := service.StreamTicks(streamCtx, []string{cfg.TestSymbol})

	fmt.Printf("Streaming %s ticks (max %d events or %d seconds)...\n", cfg.TestSymbol, MAX_STREAMING_EVENTS, MAX_STREAMING_SECONDS)

	eventCount := 0
streamLoop:
	for {
		select {
		case tick, ok := <-tickCh:
			if !ok {
				break streamLoop
			}
			eventCount++
			spread := tick.Ask - tick.Bid
			fmt.Printf("  Tick #%d: Bid=%.5f Ask=%.5f Spread=%.5f Time=%s\n",
				eventCount, tick.Bid, tick.Ask, spread, tick.Time.Format("15:04:05"))

			if eventCount >= MAX_STREAMING_EVENTS {
				streamCancel()
				break streamLoop
			}

		case err := <-errCh:
			if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				helpers.PrintShortError(err, "Stream error")
			}
			break streamLoop

		case <-streamCtx.Done():
			break streamLoop
		}
	}

	fmt.Printf("\n✓ Received %d tick events\n", eventCount)

    // #endregion

	// ═══════════════════════════════════════════════════════════════════════
	// FINAL: DISCONNECT
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nFINAL: Disconnect")
	fmt.Println("───────────────────────────────────────────────────────────")

	_, err = service.GetAccount().Disconnect(ctx, &pb.DisconnectRequest{})
	if !helpers.PrintShortError(err, "Disconnect failed") {
		fmt.Println("✓ Disconnected successfully")
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("✓ DEMO COMPLETED")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("SUMMARY:")
	fmt.Println("  • Demonstrated 19 MT5Service methods (excluding streams)")
	fmt.Println("  • MT5Service provides cleaner API than low-level MT5Account")
	fmt.Println("  • Native Go types (time.Time, float64, int64, string)")
	fmt.Println("  • 30-80% less code for common operations")
	fmt.Println("  • Better separation (positions/orders/deals in separate slices)")
	fmt.Println("  • Direct value returns (no Data struct extraction)")
	fmt.Println("  • No Request object creation needed for most methods")
	fmt.Println()
	fmt.Println("KEY ADVANTAGES SHOWN:")
	fmt.Println("  ✓ Account: GetAccountInteger/String return direct values")
	fmt.Println("  ✓ Symbol: GetSymbolInteger/String/Total return direct values")
	fmt.Println("  ✓ History: GetOrderHistory separates orders and deals")
	fmt.Println("  ✓ Market Depth: Subscribe/Get/Unsubscribe take string symbol")
	fmt.Println()
	fmt.Println("NEXT STEPS:")
	fmt.Println("  1. Review MT5Service.go source: examples/mt5/MT5Service.go")
	fmt.Println("  2. Compare with lowlevel demos to see the difference")
	fmt.Println("  3. Build your trading application with MT5Service!")

	return nil
}
