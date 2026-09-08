/*══════════════════════════════════════════════════════════════════════════════
 FILE: 01_general_operations.go - LOW-LEVEL MT5 API INFORMATION DEMO

 PURPOSE:
   Comprehensive demonstration of MT5 information retrieval methods via MT5Account.
   This is a REFERENCE GUIDE for account, symbol, position, and market data queries
   WITHOUT trading operations (see 02_trading_operations.go for trading examples).
   

 📚 WHAT THIS DEMO COVERS (6 Steps):

   STEP 1: CREATE MT5ACCOUNT INSTANCE
      • NewMT5Account() - Initialize MT5 account object with credentials

   STEP 2: CONNECTION TO MT5 SERVER
      • ConnectEx() - Connect to MT5 cluster via gRPC proxy (RECOMMENDED)
      • CheckConnect() - Verify connection status

   STEP 3: ACCOUNT INFORMATION METHODS
      • AccountSummary() - Get all account data in one call (RECOMMENDED)
      • AccountInfoDouble() - Individual double properties (Balance, Equity, Margin, etc.)
      • AccountInfoInteger() - Integer properties (Login, Leverage, etc.)
      • AccountInfoString() - String properties (Currency, Company, etc.)

   STEP 4: SYMBOL INFORMATION & OPERATIONS
      • SymbolsTotal() - Count total/selected symbols
      • SymbolExist() - Check if symbol exists
      • SymbolName() - Get symbol name by index from Market Watch
      • SymbolSelect() - Add/remove symbol from Market Watch
      • SymbolIsSynchronized() - Check sync status with server
      • SymbolInfoDouble() - Bid, Ask, Point, Volume Min/Max/Step
      • SymbolInfoInteger() - Digits, Spread, Stops Level
      • SymbolInfoString() - Description, Base/Profit Currency
      • SymbolInfoMarginRate() - Get margin requirements
      • SymbolInfoTick() - Get last tick data with timestamp
      • SymbolInfoSessionQuote() - Quote session times
      • SymbolInfoSessionTrade() - Trade session times
      • SymbolParamsMany() - Detailed parameters for multiple symbols
      • TickValueWithSize() - Tick value and size info

   STEP 5: POSITIONS & ORDERS INFORMATION
      • PositionsTotal() - Count open positions
      • OpenedOrders() - Get all opened orders & positions
      • OpenedOrdersTickets() - Get only ticket numbers (lightweight)
      • OrderHistory() - Historical orders with pagination
      • PositionsHistory() - Historical positions with P&L

   STEP 6: MARKET DEPTH (DOM - Depth of Market)
      • MarketBookAdd() - Subscribe to DOM updates
      • MarketBookGet() - Get current market depth snapshot
      • MarketBookRelease() - Unsubscribe from DOM
      ⚠️  Note: DOM typically NOT available for forex pairs on demo accounts

 🔄 API LEVELS COMPARISON:

   LOW-LEVEL (MT5Account) - THIS FILE:
   ✓ Direct gRPC/protobuf calls
   ✓ Maximum control and flexibility
   ✓ See exactly what MT5 API returns
   ✗ More verbose code
   ✗ Manual protobuf structure handling

   MID-LEVEL (MT5Service):
   ✓ Wrapper with native Go types (time.Time, float64)
   ✓ 30-50% code reduction
   ✓ Cleaner error handling
   ✗ Less control over exact API calls

   HIGH-LEVEL (MT5Sugar):
   ✓ One-liner operations with smart defaults
   ✓ Ultra-simple API for common tasks
   ✓ Best for quick prototyping
   ✗ Least flexibility

 🚀 HOW TO RUN THIS DEMO:
   cd examples/demos
   go run main.go 1          (or select [1] from interactive menu)

══════════════════════════════════════════════════════════════════════════════*/

package lowlevel

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RunGeneral01() error {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("MT5 LOWLEVEL DEMO 01: ACCOUNT INFORMATION")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// ═══════════════════════════════════════════════════════════════════════
	// LOAD CONFIGURATION
	// ═══════════════════════════════════════════════════════════════════════
	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 1: CREATE MT5ACCOUNT INSTANCE
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("STEP 1: Creating MT5Account instance...")
	fmt.Println("───────────────────────────────────────────────────────────")

	account, err := mt5.NewMT5Account(cfg.User, cfg.Password, cfg.GrpcServer, uuid.New())
	helpers.Fatal(err, "Failed to create MT5Account")

	fmt.Printf("✓ MT5Account created (UUID: %s)\n", account.Id)
	defer account.Close()

	// Create cancellable context for proper cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 2: CONNECTION TO MT5 SERVER
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 2: Connection to MT5 Server")
	fmt.Println("───────────────────────────────────────────────────────────")

	baseSymbol := cfg.TestSymbol
	connectExReq := &pb.ConnectExRequest{
		User:            cfg.User,
		Password:        cfg.Password,
		MtClusterName:   cfg.MtCluster,
		BaseChartSymbol: &baseSymbol,
	}

	fmt.Printf("  User:          %d\n", cfg.User)
	fmt.Printf("  Cluster:       %s\n", cfg.MtCluster)
	fmt.Printf("  Base Symbol:   %s\n", baseSymbol)
	fmt.Printf("  Context Timeout: 180 seconds\n")
	fmt.Println()

	// Use context timeout for ConnectEx (replaces old TerminalReadinessWaitingTimeoutSeconds)
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer connectCancel()

	connectData, err := account.ConnectEx(connectCtx, connectExReq)
	helpers.Fatal(err, "ConnectEx failed")

	// CRITICAL: Update account GUID with the one returned by server
	account.Id = uuid.MustParse(connectData.TerminalInstanceGuid)
	fmt.Printf("✓ Connected successfully\n")
	fmt.Printf("  Terminal GUID: %s\n", connectData.TerminalInstanceGuid)

	// CheckConnect - Verify connection
	fmt.Println("\n2.2. CheckConnect() - Verify connection status")
	checkReq := &pb.CheckConnectRequest{}
	checkData, err := account.CheckConnect(ctx, checkReq)
	helpers.Fatal(err, "CheckConnect failed")

	// Access nested field properly
	if checkData.HealthCheck != nil {
		fmt.Printf("✓ Connection alive: %v\n", checkData.HealthCheck.IsAlive)
	} else {
		fmt.Println("⚠ HealthCheck data not available")
	}

	// #region Account Info

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 3: ACCOUNT INFORMATION METHODS
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 3: Account Information Methods")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 3.1. ACCOUNT SUMMARY
	//      Get complete account information in ONE call.
	//      Returns: Balance, Equity, Margin, Leverage, Currency,
	//               Server Time, Trade Mode, and more.
	//      RECOMMENDED method for retrieving account data.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.1. AccountSummary() - Get all account data in one call")

	summaryReq := &pb.AccountSummaryRequest{}
	// Use timeout context to prevent hanging
	summaryCtx, summaryCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer summaryCancel()

	summaryData, err := account.AccountSummary(summaryCtx, summaryReq)
	helpers.Fatal(err, "AccountSummary failed")

	fmt.Println("\nAccount Summary (direct protobuf field access):")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("  Login:               %d\n", summaryData.AccountLogin)
	fmt.Printf("  UserName:            %s\n", summaryData.AccountUserName)
	fmt.Printf("  Company:             %s\n", summaryData.AccountCompanyName)
	fmt.Printf("  Currency:            %s\n", summaryData.AccountCurrency)
	fmt.Printf("  Balance:             %.2f\n", summaryData.AccountBalance)
	fmt.Printf("  Equity:              %.2f\n", summaryData.AccountEquity)
	fmt.Printf("  Credit:              %.2f\n", summaryData.AccountCredit)
	fmt.Printf("  Leverage:            1:%d\n", summaryData.AccountLeverage)
	fmt.Printf("  Trade Mode:          %v\n", summaryData.AccountTradeMode)

	// ServerTime is a protobuf Timestamp - need to convert
	if summaryData.ServerTime != nil {
		serverTime := summaryData.ServerTime.AsTime()
		fmt.Printf("  Server Time:         %s\n", serverTime.Format("2006-01-02 15:04:05"))
	}

// UTC Timezone Shift: server time offset from UTC in minutes
// For example: 120 minutes = UTC+2 (the server is 2 hours ahead of UTC)
	fmt.Printf("  UTC Timezone Shift:  %d minutes (UTC%+.1f)\n",
		summaryData.UtcTimezoneServerTimeShiftMinutes,
		float64(summaryData.UtcTimezoneServerTimeShiftMinutes)/60.0)

	// ══════════════════════════════════════════════════════════════
	// 3.2. ACCOUNT INFO DOUBLE
	//      Get individual DOUBLE properties from account.
	//      Examples: Balance, Equity, Margin, Profit, Credit.
	//      NOTE: For complete data, use AccountSummary() instead.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.2. AccountInfoDouble() - Get specific double property (example: Balance)")

	balanceReq := &pb.AccountInfoDoubleRequest{
		PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
	}
	balanceData, err := account.AccountInfoDouble(ctx, balanceReq)
	if err != nil {
		helpers.PrintShortError(err, "AccountInfoDouble(BALANCE) failed")
	} else {
		fmt.Printf("  Balance:                       %.2f\n", balanceData.GetRequestedValue())
	}

	// ══════════════════════════════════════════════════════════════
	// 3.3. ACCOUNT INFO INTEGER
	//      Get individual INTEGER properties from account.
	//      Examples: Login, Leverage, Trade Mode, Limit Orders.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.3. AccountInfoInteger() - Get specific integer property (example: Leverage)")

	leverageReq := &pb.AccountInfoIntegerRequest{
		PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE,
	}
	leverageData, err := account.AccountInfoInteger(ctx, leverageReq)
	if err != nil {
		helpers.PrintShortError(err, "AccountInfoInteger(LEVERAGE) failed")
	} else {
		fmt.Printf("  Leverage:                 1:%d\n", leverageData.GetRequestedValue())
	}

	// ══════════════════════════════════════════════════════════════
	// 3.4. ACCOUNT INFO STRING
	//      Get individual STRING properties from account.
	//      Examples: Name, Company, Server, Currency.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.4. AccountInfoString() - Get specific string property (example: Company)")

	companyReq := &pb.AccountInfoStringRequest{
		PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY,
	}
	companyData, err := account.AccountInfoString(ctx, companyReq)
	if err != nil {
		helpers.PrintShortError(err, "AccountInfoString(COMPANY) failed")
	} else {
		fmt.Printf("  Company:                  %s\n", companyData.GetRequestedValue())
	}

	// #endregion Account Info

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 4: SYMBOL INFORMATION & OPERATIONS
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 4: Symbol Information & Operations")
	fmt.Println("───────────────────────────────────────────────────────────")

	//#region Symbol Info

	// ══════════════════════════════════════════════════════════════
	// 4.1. SYMBOLS TOTAL
	//      Count available symbols or symbols in Market Watch.
	//      Mode: false = all symbols, true = Market Watch only.
	//      Useful for iterating through symbols.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.1. SymbolsTotal() - Count symbols")

	// Count all available symbols
	allSymbolsReq := &pb.SymbolsTotalRequest{
		Mode: false, // false = all symbols, true = Market Watch only
	}
	allSymbolsData, err := account.SymbolsTotal(ctx, allSymbolsReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolsTotal(all) failed")
	} else {
		fmt.Printf("  Total available symbols:       %d\n", allSymbolsData.Total)
	}

	// Count symbols in Market Watch only
	selectedSymbolsReq := &pb.SymbolsTotalRequest{
		Mode: true, // true = Market Watch only
	}
	selectedSymbolsData, err := account.SymbolsTotal(ctx, selectedSymbolsReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolsTotal(selected) failed")
	} else {
		fmt.Printf("  Symbols in Market Watch:       %d\n", selectedSymbolsData.Total)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.2. SYMBOL EXIST
	//      Check if a symbol exists in terminal.
	//      Returns: Exists (bool), IsCustom (bool).
	//      Useful before querying symbol properties.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.2. SymbolExist() - Check if symbol exists")

	existReq := &pb.SymbolExistRequest{
		Name: cfg.TestSymbol, // Direct protobuf field: Name (not Symbol)
	}
	existData, err := account.SymbolExist(ctx, existReq)
	if !helpers.PrintShortError(err, fmt.Sprintf("SymbolExist(%s) failed", cfg.TestSymbol)) {
		// Direct field access: Exists (not Exist)
		fmt.Printf("  Symbol '%s' exists:      %v\n", cfg.TestSymbol, existData.Exists)
		fmt.Printf("  Is custom symbol:         %v\n", existData.IsCustom)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.3. SYMBOL NAME
	//      Get symbol name by index from Market Watch.
	//      Index starts at 0. Use with SymbolsTotal().
	//      Returns symbol name as string.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.3. SymbolName() - Get symbol name by index (from Market Watch)")

	// Use actual count from SymbolsTotal
	var symbolsToShow int32 = 3
	if selectedSymbolsData != nil && selectedSymbolsData.Total < symbolsToShow {
		symbolsToShow = selectedSymbolsData.Total
	}

	if symbolsToShow == 0 {
		fmt.Println("  No symbols in Market Watch")
	} else {
		fmt.Printf("  Showing first %d symbols from Market Watch:\n", symbolsToShow)
		for i := int32(0); i < symbolsToShow; i++ {
			nameReq := &pb.SymbolNameRequest{
				Index:    i,    // Direct protobuf field: Index (not Pos)
				Selected: true, // true = Market Watch, false = all symbols
			}
			nameData, err := account.SymbolName(ctx, nameReq)
			if !helpers.PrintShortError(err, fmt.Sprintf("SymbolName(pos=%d) failed", i)) {
				if nameData.Name != "" {
					fmt.Printf("    [%d] %s\n", i, nameData.Name)
				} else {
					fmt.Printf("    [%d] (empty - no symbol at this position)\n", i)
				}
			}
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 4.4. SYMBOL SELECT
	//      Add or remove symbol from Market Watch.
	//      Select: true = add, false = remove.
	//      Must be selected to receive quotes and trade.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.4. SymbolSelect() - Add/remove symbol from Market Watch")

	selectReq := &pb.SymbolSelectRequest{
		Symbol: cfg.TestSymbol,
		Select: true, // true = add to Market Watch, false = remove
	}
	selectData, err := account.SymbolSelect(ctx, selectReq)
	if !helpers.PrintShortError(err, fmt.Sprintf("SymbolSelect(%s) failed", cfg.TestSymbol)) {
		fmt.Printf("  Symbol '%s' added to Market Watch: %v\n", cfg.TestSymbol, selectData.Success)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.5. SYMBOL IS SYNCHRONIZED
	//      Check if symbol data is synchronized with server.
	//      Returns: IsSynchronized (bool).
	//      Ensures symbol has latest quotes before trading.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.5. SymbolIsSynchronized() - Check sync status with server")

	syncReq := &pb.SymbolIsSynchronizedRequest{
		Symbol: cfg.TestSymbol,
	}
	syncData, err := account.SymbolIsSynchronized(ctx, syncReq)
	if !helpers.PrintShortError(err, fmt.Sprintf("SymbolIsSynchronized(%s) failed", cfg.TestSymbol)) {
		// Direct field access: Synchronized (not IsSynchronized)
		fmt.Printf("  Symbol '%s' synchronized:  %v\n", cfg.TestSymbol, syncData.Synchronized)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.6. SYMBOL INFO DOUBLE
	//      Get individual DOUBLE properties of symbol.
	//      Examples: Bid, Ask, Point, Volume Min/Max/Step, Spread.
	//      Returns single property value as double.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.6. SymbolInfoDouble() - Get double properties")

	bidReq := &pb.SymbolInfoDoubleRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoDoubleProperty_SYMBOL_BID, // Direct field: Type (not PropertyId)
	}
	bidData, err := account.SymbolInfoDouble(ctx, bidReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoDouble(BID) failed")
	} else {
		// Direct field access: Value (not GetRequestedValue())
		fmt.Printf("  Bid price (SYMBOL_BID):        %.5f\n", bidData.Value)
	}

	askReq := &pb.SymbolInfoDoubleRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoDoubleProperty_SYMBOL_ASK,
	}
	askData, err := account.SymbolInfoDouble(ctx, askReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoDouble(ASK) failed")
	} else {
		fmt.Printf("  Ask price (SYMBOL_ASK):        %.5f\n", askData.Value)
	}

	pointReq := &pb.SymbolInfoDoubleRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoDoubleProperty_SYMBOL_POINT,
	}
	pointData, err := account.SymbolInfoDouble(ctx, pointReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoDouble(POINT) failed")
	} else {
		fmt.Printf("  Point size (SYMBOL_POINT):     %.5f\n", pointData.Value)
	}

	volumeMinReq := &pb.SymbolInfoDoubleRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoDoubleProperty_SYMBOL_VOLUME_MIN,
	}
	volumeMinData, err := account.SymbolInfoDouble(ctx, volumeMinReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoDouble(VOLUME_MIN) failed")
	} else {
		fmt.Printf("  Min volume (VOLUME_MIN):       %.2f\n", volumeMinData.Value)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.7. SYMBOL INFO INTEGER
	//      Get individual INTEGER properties of symbol.
	//      Examples: Digits, Spread, Stops Level, Trade Mode.
	//      Returns single property value as integer.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.7. SymbolInfoInteger() - Get integer properties")

	digitsReq := &pb.SymbolInfoIntegerRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS, // Direct field: Type (not PropertyId)
	}
	digitsData, err := account.SymbolInfoInteger(ctx, digitsReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoInteger(DIGITS) failed")
	} else {
		// Direct field access: Value (not GetRequestedValue())
		fmt.Printf("  Digits (SYMBOL_DIGITS):        %d\n", digitsData.Value)
	}

	spreadReq := &pb.SymbolInfoIntegerRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoIntegerProperty_SYMBOL_SPREAD,
	}
	spreadData, err := account.SymbolInfoInteger(ctx, spreadReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoInteger(SPREAD) failed")
	} else {
		fmt.Printf("  Spread (SYMBOL_SPREAD):        %d points\n", spreadData.Value)
	}

	stopsLevelReq := &pb.SymbolInfoIntegerRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoIntegerProperty_SYMBOL_TRADE_STOPS_LEVEL,
	}
	stopsLevelData, err := account.SymbolInfoInteger(ctx, stopsLevelReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoInteger(STOPS_LEVEL) failed")
	} else {
		fmt.Printf("  Stops level (STOPS_LEVEL):     %d points\n", stopsLevelData.Value)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.8. SYMBOL INFO STRING
	//      Get individual STRING properties of symbol.
	//      Examples: Description, Base Currency, Profit Currency.
	//      Returns single property value as string.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.8. SymbolInfoString() - Get string properties")

	descReq := &pb.SymbolInfoStringRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoStringProperty_SYMBOL_DESCRIPTION, // Direct field: Type (not PropertyId)
	}
	descData, err := account.SymbolInfoString(ctx, descReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoString(DESCRIPTION) failed")
	} else {
		// Direct field access: Value (not GetRequestedValue())
		fmt.Printf("  Description:                   %s\n", descData.Value)
	}

	baseReq := &pb.SymbolInfoStringRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_BASE,
	}
	baseData, err := account.SymbolInfoString(ctx, baseReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoString(CURRENCY_BASE) failed")
	} else {
		fmt.Printf("  Base currency:                 %s\n", baseData.Value)
	}

	profitCurrReq := &pb.SymbolInfoStringRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_PROFIT,
	}
	profitCurrData, err := account.SymbolInfoString(ctx, profitCurrReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoString(CURRENCY_PROFIT) failed")
	} else {
		fmt.Printf("  Profit currency:               %s\n", profitCurrData.Value)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.9. SYMBOL INFO MARGIN RATE
	//      Get margin requirements for symbol.
	//      Returns: InitialMarginRate, MaintenanceMarginRate.
	//      Used for calculating required margin.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.9. SymbolInfoMarginRate() - Get margin requirements")

	marginRateReq := &pb.SymbolInfoMarginRateRequest{
		Symbol:    cfg.TestSymbol,
		OrderType: pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY, // Direct protobuf enum
	}
	marginRateData, err := account.SymbolInfoMarginRate(ctx, marginRateReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoMarginRate failed")
	} else {
		// Direct field access to InitialMarginRate and MaintenanceMarginRate
		fmt.Printf("  Initial margin rate (BUY):     %.2f\n", marginRateData.InitialMarginRate)
		fmt.Printf("  Maintenance margin rate (BUY): %.2f\n", marginRateData.MaintenanceMarginRate)
	}

	// ══════════════════════════════════════════════════════════════
	// 4.10. SYMBOL INFO TICK
	//      Get last tick (price update) for symbol.
	//      Returns: Time, Bid, Ask, Last, Volume, Flags.
	//      Most recent price data from server.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.10. SymbolInfoTick() - Get last tick data with timestamp")

	tickReq := &pb.SymbolInfoTickRequest{
		Symbol: cfg.TestSymbol,
	}
	tickData, err := account.SymbolInfoTick(ctx, tickReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoTick failed")
	} else {
		fmt.Printf("  Last tick for %s:\n", cfg.TestSymbol)
		// Direct field access to Bid, Ask, Last, Volume
		fmt.Printf("    Bid:                         %.5f\n", tickData.Bid)
		fmt.Printf("    Ask:                         %.5f\n", tickData.Ask)
		fmt.Printf("    Last:                        %.5f\n", tickData.Last)
		fmt.Printf("    Volume:                      %d\n", tickData.Volume)

		// Time is int64 (Unix timestamp), NOT protobuf Timestamp
		tickTime := time.Unix(tickData.Time, 0)
		fmt.Printf("    Time:                        %s\n", tickTime.Format("2006-01-02 15:04:05"))
	}

	// ══════════════════════════════════════════════════════════════
	// 4.11. SYMBOL INFO SESSION QUOTE
	//      Get quote session trading hours.
	//      Returns: From (start time), To (end time) in seconds.
	//      When symbol quotes are available.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.11. SymbolInfoSessionQuote() - Get quote session times")

	quoteSessionReq := &pb.SymbolInfoSessionQuoteRequest{
		Symbol:       cfg.TestSymbol,
		DayOfWeek:    pb.DayOfWeek_MONDAY,
		SessionIndex: 0,
	}
	quoteSessionData, err := account.SymbolInfoSessionQuote(ctx, quoteSessionReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoSessionQuote failed")
	} else {
		fmt.Printf("  Monday quote session #0:\n")
		// From and To are *timestamppb.Timestamp - convert to seconds from day start
		if quoteSessionData.From != nil {
			fromTime := quoteSessionData.From.AsTime()
			fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
			fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
		}
		if quoteSessionData.To != nil {
			toTime := quoteSessionData.To.AsTime()
			toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()
			fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 4.12. SYMBOL INFO SESSION TRADE
	//      Get trading session hours.
	//      Returns: From (start time), To (end time) in seconds.
	//      When symbol can be traded (buy/sell).
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.12. SymbolInfoSessionTrade() - Get trade session times")

	tradeSessionReq := &pb.SymbolInfoSessionTradeRequest{
		Symbol:       cfg.TestSymbol,
		DayOfWeek:    pb.DayOfWeek_MONDAY,
		SessionIndex: 0,
	}
	tradeSessionData, err := account.SymbolInfoSessionTrade(ctx, tradeSessionReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoSessionTrade failed")
	} else {
		fmt.Printf("  Monday trade session #0:\n")
		// From and To are *timestamppb.Timestamp - convert to seconds from day start
		if tradeSessionData.From != nil {
			fromTime := tradeSessionData.From.AsTime()
			fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
			fmt.Printf("    From (seconds from day start): %d\n", fromSeconds)
		}
		if tradeSessionData.To != nil {
			toTime := tradeSessionData.To.AsTime()
			toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()
			fmt.Printf("    To (seconds from day start):   %d\n", toSeconds)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 4.13. SYMBOL PARAMS MANY
	//      Get comprehensive parameters for multiple symbols at once.
	//      Returns: Bid, Ask, Digits, Spread, Point, Volume specs, etc.
	//      RECOMMENDED for getting complete symbol data efficiently.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.13. SymbolParamsMany() - Get detailed parameters for multiple symbols")

	// Note: SymbolParamsMany uses pagination/filtering, not array of symbol names
	// For simplicity, getting first page of symbols matching test symbol
	symbolFilter := cfg.TestSymbol
	paramsManyReq := &pb.SymbolParamsManyRequest{
		SymbolName:   &symbolFilter, // Filter by symbol name (optional pointer to string)
		SortType:     nil,           // Optional sort type
		PageNumber:   nil,           // Optional page number
		ItemsPerPage: nil,           // Optional items per page
	}
	paramsManyData, err := account.SymbolParamsMany(ctx, paramsManyReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolParamsMany failed")
	} else {
		// Direct field access: SymbolInfos (not SymbolParams)
		fmt.Printf("  Retrieved parameters for %d symbols matching '%s':\n", len(paramsManyData.SymbolInfos), cfg.TestSymbol)
		// Show first 3 symbols only
		maxShow := 3
		if len(paramsManyData.SymbolInfos) < maxShow {
			maxShow = len(paramsManyData.SymbolInfos)
		}
		for i := 0; i < maxShow; i++ {
			info := paramsManyData.SymbolInfos[i]
			fmt.Printf("\n  Symbol #%d: %s\n", i+1, info.Name)
			fmt.Printf("    Bid:                         %.5f\n", info.Bid)
			fmt.Printf("    Ask:                         %.5f\n", info.Ask)
			fmt.Printf("    Digits:                      %d\n", info.Digits)
			fmt.Printf("    Spread:                      %d points\n", info.Spread)
			fmt.Printf("    Volume Min:                  %.2f\n", info.VolumeMin)
			fmt.Printf("    Volume Max:                  %.2f\n", info.VolumeMax)
			fmt.Printf("    Volume Step:                 %.2f\n", info.VolumeStep)
			fmt.Printf("    Contract Size:               %.2f\n", info.TradeContractSize)
			fmt.Printf("    Point:                       %.5f\n", info.Point)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 4.14. TICK VALUE WITH SIZE
	//      Get tick value and tick size information for symbols.
	//      Returns: TradeTickValue, TradeTickSize, TradeContractSize, etc.
	//      Useful for calculating position value and P&L accurately.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.14. SymbolInfoDouble() - Get tick value and size info")

	tickValReq := &pb.SymbolInfoDoubleRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE,
	}
	tickValData, err := account.SymbolInfoDouble(ctx, tickValReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoDouble(SYMBOL_TRADE_TICK_VALUE) failed")
	} else {
		fmt.Printf("    Trade Tick Value:        %.5f\n", tickValData.Value)
	}

	tickSizeReq := &pb.SymbolInfoDoubleRequest{
		Symbol: cfg.TestSymbol,
		Type:   pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_SIZE,
	}
	tickSizeData, err := account.SymbolInfoDouble(ctx, tickSizeReq)
	if err != nil {
		helpers.PrintShortError(err, "SymbolInfoDouble(SYMBOL_TRADE_TICK_SIZE) failed")
	} else {
		fmt.Printf("    Trade Tick Size:         %.5f\n", tickSizeData.Value)
	}

	//#endregion Symbol Info

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 5: POSITIONS & ORDERS INFORMATION
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 5: Positions & Orders Information")
	fmt.Println("───────────────────────────────────────────────────────────")

	//#region Positions Info

	// ══════════════════════════════════════════════════════════════
	// 5.1. POSITIONS TOTAL
	//      Count currently open positions.
	//      Returns: Total number of open positions.
	//      Lightweight check for open trades.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n5.1. PositionsTotal() - Count open positions")

	// Note: PositionsTotal() doesn't take a request parameter
	positionsTotalData, err := account.PositionsTotal(ctx)
	if err != nil {
		helpers.PrintShortError(err, "PositionsTotal failed")
	} else {
		// Direct field access: TotalPositions
		fmt.Printf("  Total open positions:          %d\n", positionsTotalData.TotalPositions)
	}

	// ══════════════════════════════════════════════════════════════
	// 5.2. OPENED ORDERS
	//      Get all currently opened orders and positions with details.
	//      Returns: Full position data (ticket, symbol, volume, profit, etc).
	//      Complete snapshot of open trades.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n5.2. OpenedOrders() - Get all opened orders & positions")

	openedOrdersReq := &pb.OpenedOrdersRequest{}
	openedOrdersData, err := account.OpenedOrders(ctx, openedOrdersReq)
	if err != nil {
		helpers.PrintShortError(err, "OpenedOrders failed")
	} else {
		// Direct field access: OpenedOrders and PositionInfos
		totalOrders := len(openedOrdersData.OpenedOrders) + len(openedOrdersData.PositionInfos)
		fmt.Printf("  Total opened orders/positions: %d\n", totalOrders)
		fmt.Printf("    Pending orders:              %d\n", len(openedOrdersData.OpenedOrders))
		fmt.Printf("    Open positions:              %d\n", len(openedOrdersData.PositionInfos))

		// Show first 2 positions if any exist
		maxShow := 2
		if len(openedOrdersData.PositionInfos) < maxShow {
			maxShow = len(openedOrdersData.PositionInfos)
		}
		for i := 0; i < maxShow; i++ {
			pos := openedOrdersData.PositionInfos[i]
			fmt.Printf("\n  Position #%d:\n", i+1)
			fmt.Printf("    Ticket:                      %d\n", pos.Ticket)
			fmt.Printf("    Symbol:                      %s\n", pos.Symbol)
			fmt.Printf("    Type:                        %v\n", pos.Type)
			fmt.Printf("    Volume:                      %.2f\n", pos.Volume)
			fmt.Printf("    Price Open:                  %.5f\n", pos.PriceOpen)
			fmt.Printf("    Current Price:               %.5f\n", pos.PriceCurrent)
			fmt.Printf("    Profit:                      %.2f\n", pos.Profit)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 5.3. OPENED ORDERS TICKETS
	//      Get only ticket numbers of opened orders (lightweight).
	//      Returns: Array of ticket IDs only.
	//      Fast check without full position details.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n5.3. OpenedOrdersTickets() - Get ticket numbers only")

	ticketsReq := &pb.OpenedOrdersTicketsRequest{}
	ticketsData, err := account.OpenedOrdersTickets(ctx, ticketsReq)
	if err != nil {
		helpers.PrintShortError(err, "OpenedOrdersTickets failed")
	} else {
		// Direct field access: OpenedOrdersTickets and OpenedPositionTickets
		totalTickets := len(ticketsData.OpenedOrdersTickets) + len(ticketsData.OpenedPositionTickets)
		fmt.Printf("  Total tickets:                 %d\n", totalTickets)
		fmt.Printf("    Pending order tickets:       %d\n", len(ticketsData.OpenedOrdersTickets))
		fmt.Printf("    Position tickets:            %d\n", len(ticketsData.OpenedPositionTickets))

		// Show first few position tickets if any exist
		if len(ticketsData.OpenedPositionTickets) > 0 {
			fmt.Printf("  First position tickets:        ")
			maxShow := 5
			if len(ticketsData.OpenedPositionTickets) < maxShow {
				maxShow = len(ticketsData.OpenedPositionTickets)
			}
			for i := 0; i < maxShow; i++ {
				fmt.Printf("%d ", ticketsData.OpenedPositionTickets[i])
			}
			fmt.Println()
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 5.4. ORDER HISTORY
	//      Get historical orders with time range and pagination.
	//      Returns: Closed orders with timestamps and details.
	//      Use for analyzing past trading activity.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n5.4. OrderHistory() - Get historical orders")

	// Set time range (last 30 days) - must use timestamppb.Timestamp
	now := time.Now()
	fromTimestamp := timestamppb.New(now.AddDate(0, 0, -30))
	toTimestamp := timestamppb.New(now)

	orderHistoryReq := &pb.OrderHistoryRequest{
		InputFrom: fromTimestamp,
		InputTo:   toTimestamp,
	}
	orderHistoryData, err := account.OrderHistory(ctx, orderHistoryReq)
	if err != nil {
		helpers.PrintShortError(err, "OrderHistory failed")
	} else {
		// Direct field access: HistoryData (array of *HistoryData)
		fmt.Printf("  Historical orders (last 30d):  %d\n", len(orderHistoryData.HistoryData))

		// Show first 2 orders if any exist
		maxShow := 2
		if len(orderHistoryData.HistoryData) < maxShow {
			maxShow = len(orderHistoryData.HistoryData)
		}
		for i := 0; i < maxShow; i++ {
			item := orderHistoryData.HistoryData[i]
			// HistoryData contains HistoryOrder *OrderHistoryData
			if item.HistoryOrder != nil {
				order := item.HistoryOrder
				fmt.Printf("\n  Order #%d:\n", i+1)
				fmt.Printf("    Ticket:                      %d\n", order.Ticket)
				fmt.Printf("    Symbol:                      %s\n", order.Symbol)
				fmt.Printf("    Type:                        %v\n", order.Type)
				fmt.Printf("    Volume:                      %.2f\n", order.VolumeInitial)
				fmt.Printf("    Price Open:                  %.5f\n", order.PriceOpen)
			}
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 5.5. POSITIONS HISTORY
	//      Get historical positions with P&L and details.
	//      Returns: Closed positions with profit/loss data.
	//      Analyze trading performance and results.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n5.5. PositionsHistory() - Get historical positions with P&L")

	positionsHistoryReq := &pb.PositionsHistoryRequest{
		PositionOpenTimeFrom: fromTimestamp,
		PositionOpenTimeTo:   toTimestamp,
	}
	positionsHistoryData, err := account.PositionsHistory(ctx, positionsHistoryReq)
	if err != nil {
		helpers.PrintShortError(err, "PositionsHistory failed")
	} else {
		// Direct field access: HistoryPositions (array of *PositionHistoryInfo)
		fmt.Printf("  Historical positions (last 30d): %d\n", len(positionsHistoryData.HistoryPositions))

		// Show first 2 positions if any exist
		maxShow := 2
		if len(positionsHistoryData.HistoryPositions) < maxShow {
			maxShow = len(positionsHistoryData.HistoryPositions)
		}
		for i := 0; i < maxShow; i++ {
			pos := positionsHistoryData.HistoryPositions[i]
			fmt.Printf("\n  Position #%d:\n", i+1)
			fmt.Printf("    Position Ticket:             %d\n", pos.PositionTicket)
			fmt.Printf("    Symbol:                      %s\n", pos.Symbol)
			fmt.Printf("    Order Type:                  %v\n", pos.OrderType)
			fmt.Printf("    Volume:                      %.2f\n", pos.Volume)
			fmt.Printf("    Profit:                      %.2f\n", pos.Profit)
		}
	}

	//#endregion Positions Info

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 6: MARKET DEPTH / DOM (Depth of Market) - Level 2 Data
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 6: Market Depth / DOM (Level 2 Data)")
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("⚠️  IMPORTANT: Market Book (DOM) Limitations")
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println("Market Book (Depth of Market) is typically NOT AVAILABLE for:")
	fmt.Println("  • Forex currency pairs (like EURUSD)")
	fmt.Println("  • Demo accounts")
	fmt.Println()
	fmt.Println("Why? Forex is a DECENTRALIZED market:")
	fmt.Println("  ✗ No central exchange (unlike stocks/futures)")
	fmt.Println("  ✗ No unified order book")
	fmt.Println("  ✗ Each broker only sees their own client orders")
	fmt.Println()
	fmt.Println("Market Book IS available for:")
	fmt.Println("  ✓ Exchange-traded instruments (Stocks, Futures, Options)")
	fmt.Println("  ✓ Some brokers on LIVE accounts (limited data)")
	fmt.Println()
	fmt.Println("Expected behavior for this demo:")
	fmt.Println("  → Subscription may succeed (MarketBookAdd)")
	fmt.Println("  → But data retrieval will timeout (MarketBookGet)")
	fmt.Println("  → This is NORMAL and EXPECTED")
	fmt.Println()
	fmt.Println("For detailed capability check, use: go run main.go capabilities")
	fmt.Println()

	//#region Market Book (DOM)

	// ══════════════════════════════════════════════════════════════
	// 6.1. MARKET BOOK ADD
	//      Subscribe to Market Book (Depth of Market) updates.
	//      Returns: OpenedSuccessfully (bool).
	//      NOTE: Subscription may succeed even if no data available.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("6.1. MarketBookAdd() - Subscribe to DOM")
	fmt.Println("    Attempting subscription (may succeed even if no data)...")

	// Use short timeout
	bookAddCtx, bookAddCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer bookAddCancel()

	marketBookAddReq := &pb.MarketBookAddRequest{
		Symbol: cfg.TestSymbol,
	}
	marketBookAddData, err := account.MarketBookAdd(bookAddCtx, marketBookAddReq)
	if err != nil {
		fmt.Printf("  ❌ Subscription failed: %v\n", err)
		fmt.Println("     → Broker does not support DOM for this symbol")
	} else {
		if marketBookAddData.OpenedSuccessfully {
			fmt.Printf("  ✓ Subscription accepted for '%s'\n", cfg.TestSymbol)
			fmt.Println("     → This doesn't mean data will be available")
		} else {
			fmt.Println("  ❌ Subscription rejected by broker")
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 6.2. MARKET BOOK GET
	//      Get current Market Book (Level 2) snapshot.
	//      Returns: Buy/Sell orders with prices and volumes.
	//      Shows order book depth (if available).
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n6.2. MarketBookGet() - Get DOM snapshot")
	fmt.Println("    Testing if data is actually available...")

	// Use very short timeout - we know it will likely fail
	bookGetCtx, bookGetCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer bookGetCancel()

	marketBookGetReq := &pb.MarketBookGetRequest{
		Symbol: cfg.TestSymbol,
	}
	marketBookGetData, err := account.MarketBookGet(bookGetCtx, marketBookGetReq)
	if err != nil {
		fmt.Println("  ❌ No data available (timeout)")
		fmt.Println("     → This is EXPECTED for forex pairs on demo accounts")
		fmt.Println("     → Not an error - just a limitation of forex market structure")
	} else {
		if len(marketBookGetData.MqlBookInfos) > 0 {
			fmt.Printf("  ✅ SUCCESS! Received %d price levels\n", len(marketBookGetData.MqlBookInfos))
			fmt.Println("     → This is UNUSUAL for forex - broker provides limited DOM")
			fmt.Println()

			// Show first 5 levels
			maxShow := 5
			if len(marketBookGetData.MqlBookInfos) < maxShow {
				maxShow = len(marketBookGetData.MqlBookInfos)
			}
			fmt.Println("  First few levels:")
			for i := 0; i < maxShow; i++ {
				level := marketBookGetData.MqlBookInfos[i]
				bookType := "SELL"
				if level.Type == pb.BookType_BOOK_TYPE_BUY {
					bookType = "BUY"
				}
				fmt.Printf("    [%d] %s  Price: %.5f  Volume: %d\n",
					i+1, bookType, level.Price, level.Volume)
			}
		} else {
			fmt.Println("  ⚠️  Call succeeded but no data returned")
			fmt.Println("     → Broker has no DOM data to show")
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 6.3. MARKET BOOK RELEASE
	//      Unsubscribe from Market Book updates.
	//      Returns: ClosedSuccessfully (bool).
	//      Clean up DOM subscription when done.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n6.3. MarketBookRelease() - Unsubscribe from DOM")
	fmt.Println("    Cleaning up subscription...")

	// Use short timeout
	bookReleaseCtx, bookReleaseCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer bookReleaseCancel()

	marketBookReleaseReq := &pb.MarketBookReleaseRequest{
		Symbol: cfg.TestSymbol,
	}
	marketBookReleaseData, err := account.MarketBookRelease(bookReleaseCtx, marketBookReleaseReq)
	if err != nil {
		fmt.Printf("  ⚠️  Unsubscribe failed: %v\n", err)
		fmt.Println("     → Not critical - connection will close anyway")
	} else {
		if marketBookReleaseData.ClosedSuccessfully {
			fmt.Println("  ✓ Unsubscribed successfully")
		} else {
			fmt.Println("  ⚠️  Unsubscribe reported unsuccessful")
		}
	}

	//#endregion Market Book (DOM)

	// ══════════════════════════════════════════════════════════════════════════════════
	// UNDERSTANDING gRPC/PROTOBUF REQUEST PATTERN
	// ══════════════════════════════════════════════════════════════════════════════════
	//
	// WHY DO WE ALWAYS CREATE REQUEST OBJECTS?
	//
	// You may have noticed this pattern throughout this demo:
	//
	//   summaryReq := &pb.AccountSummaryRequest{}      // ← Create request object
	//   summaryData, err := account.AccountSummary(ctx, summaryReq)
	//
	// Even when the request is EMPTY ({}), we still create it. Why?
	//
	// ANSWER: This is the standard gRPC/protobuf pattern. ALL gRPC methods have
	// the same signature, requiring a request object - even if it has no fields.
	//
	//   func MethodName(ctx context.Context, req *RequestType) (*ResponseType, error)
	//                                              ^
	//                                        Always required!
	//
	// ──────────────────────────────────────────────────────────────────────────────────
	// EXAMPLE 1: Empty Request (no parameters needed)
	// ──────────────────────────────────────────────────────────────────────────────────
	//
	//   summaryReq := &pb.AccountSummaryRequest{}  // ← Empty braces {} = no parameters
	//   summaryData, err := account.AccountSummary(ctx, summaryReq)
	//
	// AccountSummary() doesn't need parameters - just return ALL account data.
	//
	// ──────────────────────────────────────────────────────────────────────────────────
	// EXAMPLE 2: Request WITH Parameters
	// ──────────────────────────────────────────────────────────────────────────────────
	//
	//   balanceReq := &pb.AccountInfoDoubleRequest{
	//       PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,  // ← Specify property
	//   }
	//   balanceData, err := account.AccountInfoDouble(ctx, balanceReq)
	//
	// AccountInfoDouble() NEEDS to know WHICH property you want (Balance? Equity? Margin?).
	//
	// ──────────────────────────────────────────────────────────────────────────────────
	// EXAMPLE 3: Request with Multiple Parameters
	// ──────────────────────────────────────────────────────────────────────────────────
	//
	//   tickReq := &pb.SymbolInfoTickRequest{
	//       Symbol: "EURUSD",  // ← Specify symbol
	//   }
	//   tickData, err := account.SymbolInfoTick(ctx, tickReq)
	//
	// ──────────────────────────────────────────────────────────────────────────────────
	// KEY TAKEAWAY
	// ──────────────────────────────────────────────────────────────────────────────────
	//
	// This is a UNIFIED approach - all methods work the same way:
	//   1. Create request object (even if empty)
	//   2. Pass context + request to method
	//   3. Get response + error back
	//
	// Benefits:
	//   ✓ Consistent API across ALL methods
	//   ✓ Easy to add parameters later without breaking compatibility
	//   ✓ Type-safe - compiler checks request structure
	//   ✓ Clear what data each method needs
	//
	// ══════════════════════════════════════════════════════════════════════════════════

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("✓ DEMO COMPLETED SUCCESSFULLY")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}
