/*══════════════════════════════════════════════════════════════════════════════
 FILE: 10_sugar_advanced.go - ADVANCED SUGAR API DEMO (RISK MANAGEMENT & MORE)

 PURPOSE:
   Demonstrates NEW advanced MT5Sugar methods for risk management, symbol info,
   trading helpers, and account analytics. Professional risk management made EASY!

 🎯 WHO SHOULD USE THIS:
   • Professional traders using proper risk management
   • Anyone wanting automated position sizing
   • Traders working with multiple symbols
   • Performance tracking and account monitoring

 📚 WHAT THIS DEMO COVERS (4 NEW Categories):

   1. RISK MANAGEMENT METHODS 🔥 MOST IMPORTANT (4 methods)
      • CalculatePositionSize(symbol, risk%, SL_pips) - Auto position sizing
      • GetMaxLotSize(symbol) - Maximum tradeable volume
      • CanOpenPosition(symbol, volume) - Pre-trade validation
      • CalculateRequiredMargin(symbol, volume) - Margin calculation

   2. SYMBOL INFORMATION METHODS (3 methods + struct)
      • GetSymbolInfo(symbol) - Complete symbol data (includes digits, stop level)
      • GetAllSymbols() - List all available symbols
      • IsSymbolAvailable(symbol) - Check symbol tradability

   3. TRADING HELPERS (3 methods)
      • CalculateSLTP(symbol, direction, entry, SL_pips, TP_pips)
      • BuyMarketWithPips(symbol, volume, SL_pips, TP_pips)
      • SellMarketWithPips(symbol, volume, SL_pips, TP_pips)

   4. ACCOUNT INFORMATION (2 methods + 2 structs)
      • GetAccountInfo() - Complete account snapshot
      • GetDailyStats() - Today's trading statistics

 💡 KEY FEATURES DEMONSTRATED:
   • Professional risk management (risk 2% per trade)
   • Automated position sizing based on account balance
   • Pre-trade validation to prevent errors
   • Working with pips instead of raw prices
   • Comprehensive account analytics

 🔄 API LEVELS:
   HIGH-LEVEL (Sugar) - THIS FILE: One-liners, smart defaults, automatic calculations
   MID-LEVEL (Service): More control, native Go types
   LOW-LEVEL (Account): Full control, protobuf structures

 🚀 HOW TO RUN:
   cd examples/demos
   go run main.go 10    (or select [10] from menu)

 ⚠️  SAFE DEMO - READ ONLY:
   • This demo is SAFE - mostly read operations
   • Optional: Can open 1 test position with risk management
   • All trading actions require explicit confirmation
   • Recommended for DEMO account, but safe for LIVE too
══════════════════════════════════════════════════════════════════════════════*/

package sugar

import (
	"fmt"
	"strings"
	"time"

	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// RunSugarAdvancedDemo demonstrates advanced Sugar API features
func RunSugarAdvancedDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("MT5 SUGAR API - ADVANCED FEATURES DEMO (Risk Management)")
	fmt.Println(strings.Repeat("=", 80))

	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// Create Sugar instance
	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	helpers.Fatal(err, "Failed to create Sugar instance")

	// Connect to MT5
	fmt.Println("\n🔌 Connecting to MT5...")
	err = sugar.QuickConnect(cfg.MtCluster)
	helpers.Fatal(err, "Failed to connect to MT5")
	fmt.Println("✅ Connected successfully!")

	printHeader("🎯 ADVANCED SUGAR API DEMO - RISK MANAGEMENT & MORE")

	// ══════════════════════════════════════════════════════════════════════
	// PART 1: ACCOUNT INFORMATION - Get complete account snapshot
	// ══════════════════════════════════════════════════════════════════════
	printSection("1. ACCOUNT INFORMATION - Complete Snapshot")

	accountInfo, err := sugar.GetAccountInfo()
	if err != nil {
		printError("GetAccountInfo", err)
		return
	}

	printSuccess("GetAccountInfo()", "Complete account data retrieved")
	fmt.Printf("    Account Login:    %d\n", accountInfo.Login)
	fmt.Printf("   💰 Balance:          %.2f %s\n", accountInfo.Balance, accountInfo.Currency)
	fmt.Printf("    Equity:           %.2f %s\n", accountInfo.Equity, accountInfo.Currency)
	fmt.Printf("    Margin Used:      %.2f %s\n", accountInfo.Margin, accountInfo.Currency)
	fmt.Printf("    Free Margin:      %.2f %s\n", accountInfo.FreeMargin, accountInfo.Currency)
	fmt.Printf("    Margin Level:     %.2f%%\n", accountInfo.MarginLevel)
	fmt.Printf("    Floating P/L:     %.2f %s\n", accountInfo.Profit, accountInfo.Currency)
	fmt.Printf("     Leverage:         1:%d\n", accountInfo.Leverage)
	fmt.Printf("    Broker:           %s\n", accountInfo.Company)
	fmt.Println()

	// ══════════════════════════════════════════════════════════════════════
	// PART 2: DAILY STATISTICS - Trading performance tracking
	// ══════════════════════════════════════════════════════════════════════
	printSection("2. DAILY STATISTICS - Today's Performance")

	stats, err := sugar.GetDailyStats()
	if err != nil {
		printError("GetDailyStats", err)
	} else {
		printSuccess("GetDailyStats()", "Today's statistics calculated")
		fmt.Printf("    Total Deals:      %d\n", stats.TotalDeals)
		fmt.Printf("    Winning Deals:    %d\n", stats.WinningDeals)
		fmt.Printf("    Losing Deals:     %d\n", stats.LosingDeals)
		fmt.Printf("   📈 Win Rate:         %.1f%%\n", stats.WinRate)
		fmt.Printf("    Total Profit:     %.2f %s\n", stats.TotalProfit, accountInfo.Currency)
		fmt.Printf("    Best Deal:        %.2f %s\n", stats.BestDeal, accountInfo.Currency)
		fmt.Printf("    Worst Deal:       %.2f %s\n", stats.WorstDeal, accountInfo.Currency)
		fmt.Println()
	}

	// ══════════════════════════════════════════════════════════════════════
	// PART 3: SYMBOL INFORMATION - Complete symbol data
	// ══════════════════════════════════════════════════════════════════════
	printSection("3. SYMBOL INFORMATION - Comprehensive Symbol Data")

	symbol := "EURUSD"
	symbolInfo, err := sugar.GetSymbolInfo(symbol)
	if err != nil {
		printError("GetSymbolInfo", err)
		return
	}

	printSuccess(fmt.Sprintf("GetSymbolInfo(\"%s\")", symbol), "Complete symbol data retrieved")
	fmt.Printf("    Symbol:           %s\n", symbolInfo.Name)
	fmt.Printf("    BID:              %.5f\n", symbolInfo.Bid)
	fmt.Printf("    ASK:              %.5f\n", symbolInfo.Ask)
	fmt.Printf("    Spread:           %d points\n", symbolInfo.Spread)
	fmt.Printf("    Digits:           %d\n", symbolInfo.Digits)
	fmt.Printf("    Point:            %.5f\n", symbolInfo.Point)
	fmt.Printf("    Volume Min:       %.2f\n", symbolInfo.VolumeMin)
	fmt.Printf("    Volume Max:       %.2f\n", symbolInfo.VolumeMax)
	fmt.Printf("    Volume Step:      %.2f\n", symbolInfo.VolumeStep)
	fmt.Printf("    Min Stop Level:   %d points\n", symbolInfo.StopLevel)
	fmt.Printf("    Contract Size:    %.0f\n", symbolInfo.ContractSize)
	fmt.Println()

	// Check symbol availability
	available, err := sugar.IsSymbolAvailable(symbol)
	if err != nil {
		printError("IsSymbolAvailable", err)
	} else {
		if available {
			printSuccess(fmt.Sprintf("IsSymbolAvailable(\"%s\")", symbol), "Symbol is available for trading")
		} else {
			printWarning(fmt.Sprintf("IsSymbolAvailable(\"%s\")", symbol), "Symbol is NOT available")
		}
	}

	// Get all symbols
	allSymbols, err := sugar.GetAllSymbols()
	if err != nil {
		printError("GetAllSymbols", err)
	} else {
		printSuccess("GetAllSymbols()", fmt.Sprintf("Found %d symbols", len(allSymbols)))
		fmt.Printf("   First 10 symbols: ")
		for i := 0; i < 10 && i < len(allSymbols); i++ {
			fmt.Printf("%s ", allSymbols[i])
		}
		fmt.Println("...")
		fmt.Println()
	}

	// ══════════════════════════════════════════════════════════════════════
	// PART 4: RISK MANAGEMENT - THE MOST IMPORTANT PART! 🔥
	// ══════════════════════════════════════════════════════════════════════
	printSection("4. RISK MANAGEMENT - Professional Position Sizing 🔥")

	// Example: Risk 2% of balance with 50 pip stop loss
	riskPercent := 2.0
	stopLossPips := 50.0

	lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLossPips)
	if err != nil {
		printError("CalculatePositionSize", err)
	} else {
		printSuccess(
			fmt.Sprintf("CalculatePositionSize(\"%s\", %.1f%%, %.0f pips)", symbol, riskPercent, stopLossPips),
			fmt.Sprintf("Recommended lot size: %.2f", lotSize),
		)
		riskAmount := accountInfo.Balance * riskPercent / 100.0
		fmt.Printf("    Balance:          %.2f %s\n", accountInfo.Balance, accountInfo.Currency)
		fmt.Printf("     Risk Amount:      %.2f %s (%.1f%%)\n", riskAmount, accountInfo.Currency, riskPercent)
		fmt.Printf("   🛑 Stop Loss:        %.0f pips\n", stopLossPips)
		fmt.Printf("    Position Size:    %.2f lots\n", lotSize)
		fmt.Printf("    If SL hits, you lose exactly %.2f %s (%.1f%% of balance)\n",
			riskAmount, accountInfo.Currency, riskPercent)
		fmt.Println()
	}

	// Calculate maximum lot size
	maxLots, err := sugar.GetMaxLotSize(symbol)
	if err != nil {
		printError("GetMaxLotSize", err)
	} else {
		printSuccess(fmt.Sprintf("GetMaxLotSize(\"%s\")", symbol), fmt.Sprintf("Max lots: %.2f", maxLots))
		fmt.Printf("    Free Margin:      %.2f %s\n", accountInfo.FreeMargin, accountInfo.Currency)
		fmt.Printf("    Maximum Lots:     %.2f (with 80%% safety buffer)\n", maxLots)
		fmt.Println()
	}

	// Validate if we can open a position
	testVolume := 0.1
	canOpen, reason, err := sugar.CanOpenPosition(symbol, testVolume)
	if err != nil {
		printError("CanOpenPosition", err)
	} else {
		if canOpen {
			printSuccess(
				fmt.Sprintf("CanOpenPosition(\"%s\", %.2f)", symbol, testVolume),
				"✅ Position CAN be opened - all checks passed",
			)
		} else {
			printWarning(
				fmt.Sprintf("CanOpenPosition(\"%s\", %.2f)", symbol, testVolume),
				fmt.Sprintf("❌ Position CANNOT be opened: %s", reason),
			)
		}
	}

	// Calculate required margin for a position
	requiredMargin, err := sugar.CalculateRequiredMargin(symbol, testVolume)
	if err != nil {
		printError("CalculateRequiredMargin", err)
	} else {
		printSuccess(
			fmt.Sprintf("CalculateRequiredMargin(\"%s\", %.2f)", symbol, testVolume),
			fmt.Sprintf("Required margin: %.2f %s", requiredMargin, accountInfo.Currency),
		)
		fmt.Println()
	}

	// ══════════════════════════════════════════════════════════════════════
	// PART 5: TRADING HELPERS - Working with pips instead of prices
	// ══════════════════════════════════════════════════════════════════════
	printSection("5. TRADING HELPERS - Pip-Based SL/TP Calculation")

	// Calculate SL/TP from pips
	direction := "BUY"
	entryPrice := 0.0 // 0 = use current market price
	slPips := 50.0
	tpPips := 100.0

	sl, tp, err := sugar.CalculateSLTP(symbol, direction, entryPrice, slPips, tpPips)
	if err != nil {
		printError("CalculateSLTP", err)
	} else {
		currentAsk, _ := sugar.GetAsk(symbol)
		printSuccess(
			fmt.Sprintf("CalculateSLTP(\"%s\", \"%s\", market, %.0f, %.0f)", symbol, direction, slPips, tpPips),
			"SL/TP prices calculated from pips",
		)
		fmt.Printf("    Current ASK:      %.5f\n", currentAsk)
		fmt.Printf("    Stop Loss:        %.5f (%.0f pips below entry)\n", sl, slPips)
		fmt.Printf("    Take Profit:      %.5f (%.0f pips above entry)\n", tp, tpPips)
		fmt.Printf("    Risk/Reward:      1:%.1f\n", tpPips/slPips)
		fmt.Println()
	}

	// Demonstrate BuyMarketWithPips and SellMarketWithPips
	printInfo("OPTIONAL TRADING DEMO", "You can test pip-based trading methods:")
	fmt.Printf("   Example 1: sugar.BuyMarketWithPips(\"%s\", 0.01, 50, 100)\n", symbol)
	fmt.Printf("              Opens BUY 0.01 lots, SL=50 pips, TP=100 pips\n")
	fmt.Printf("   \n")
	fmt.Printf("   Example 2: sugar.SellMarketWithPips(\"%s\", 0.01, 50, 100)\n", symbol)
	fmt.Printf("              Opens SELL 0.01 lots, SL=50 pips, TP=100 pips\n")
	fmt.Printf("   \n")
	fmt.Printf("   💡 These methods automatically calculate exact SL/TP prices!\n")
	fmt.Println()

	// ══════════════════════════════════════════════════════════════════════
	// PART 6: PRACTICAL EXAMPLE - Complete risk-managed trade setup
	// ══════════════════════════════════════════════════════════════════════
	printSection("6. PRACTICAL EXAMPLE - Complete Risk-Managed Trade")

	fmt.Println("   🎯 SCENARIO: Open a BUY position with proper risk management")
	fmt.Println()

	tradingSymbol := "EURUSD"
	targetRisk := 2.0      // Risk 2% of balance
	targetSL := 50.0       // 50 pip stop loss
	targetTP := 100.0      // 100 pip take profit (1:2 R/R)

	fmt.Printf("   1️⃣  Symbol:           %s\n", tradingSymbol)
	fmt.Printf("   2️⃣  Risk:             %.1f%% of balance\n", targetRisk)
	fmt.Printf("   3️⃣  Stop Loss:        %.0f pips\n", targetSL)
	fmt.Printf("   4️⃣  Take Profit:      %.0f pips\n", targetTP)
	fmt.Printf("   5️⃣  Risk/Reward:      1:%.1f\n", targetTP/targetSL)
	fmt.Println()

	// Step 1: Check symbol availability
	fmt.Println("   STEP 1: Verify symbol availability...")
	isAvailable, err := sugar.IsSymbolAvailable(tradingSymbol)
	if err != nil || !isAvailable {
		fmt.Println("   ❌ Symbol not available!")
		return
	}
	fmt.Println("   ✅ Symbol is available")

	// Step 2: Calculate optimal position size
	fmt.Println("   STEP 2: Calculate position size based on risk...")
	optimalLots, err := sugar.CalculatePositionSize(tradingSymbol, targetRisk, targetSL)
	if err != nil {
		printError("   Position size calculation failed", err)
		return
	}
	fmt.Printf("   ✅ Optimal lot size: %.2f lots\n", optimalLots)

	// Step 3: Validate we can open this position
	fmt.Println("   STEP 3: Pre-trade validation...")
	canTrade, validationReason, err := sugar.CanOpenPosition(tradingSymbol, optimalLots)
	if err != nil {
		printError("   Validation failed", err)
		return
	}
	if !canTrade {
		fmt.Printf("   ❌ Cannot open position: %s\n", validationReason)
		return
	}
	fmt.Println("   ✅ All validation checks passed")

	// Step 4: Calculate required margin
	fmt.Println("   STEP 4: Calculate margin requirement...")
	margin, err := sugar.CalculateRequiredMargin(tradingSymbol, optimalLots)
	if err != nil {
		printError("   Margin calculation failed", err)
		return
	}
	fmt.Printf("   ✅ Required margin: %.2f %s\n", margin, accountInfo.Currency)

	// Step 5: Show complete trade setup
	fmt.Println("   STEP 5: Complete trade setup ready:")
	fmt.Println()
	fmt.Println("   ╔════════════════════════════════════════════════════════╗")
	fmt.Println("   ║          READY TO EXECUTE - ALL CHECKS PASSED          ║")
	fmt.Println("   ╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("    Symbol:           %s\n", tradingSymbol)
	fmt.Printf("    Direction:        BUY\n")
	fmt.Printf("    Lot Size:         %.2f lots\n", optimalLots)
	fmt.Printf("    Stop Loss:        %.0f pips\n", targetSL)
	fmt.Printf("    Take Profit:      %.0f pips\n", targetTP)
	fmt.Printf("   ⚠️  Risk Amount:      %.2f %s (%.1f%%)\n",
		accountInfo.Balance*targetRisk/100, accountInfo.Currency, targetRisk)
	fmt.Printf("    Potential Reward: %.2f %s (%.1f%%)\n",
		accountInfo.Balance*targetRisk*2/100, accountInfo.Currency, targetRisk*2)
	fmt.Printf("    Margin Required:  %.2f %s\n", margin, accountInfo.Currency)
	fmt.Println()
	fmt.Println("    To execute this trade, use:")
	fmt.Printf("   ticket, _ := sugar.BuyMarketWithPips(\"%s\", %.2f, %.0f, %.0f)\n",
		tradingSymbol, optimalLots, targetSL, targetTP)
	fmt.Println()

	// ══════════════════════════════════════════════════════════════════════
	// SUMMARY
	// ══════════════════════════════════════════════════════════════════════
	printSection("✨ DEMO COMPLETE - Advanced Features Summary")

	fmt.Println("   ✅ Account Information:     GetAccountInfo(), GetDailyStats()")
	fmt.Println("   ✅ Symbol Information:      GetSymbolInfo(), GetAllSymbols(), IsSymbolAvailable()")
	fmt.Println("   ✅ Risk Management:         CalculatePositionSize(), GetMaxLotSize(), CanOpenPosition()")
	fmt.Println("   ✅ Trading Helpers:         CalculateSLTP(), BuyMarketWithPips(), SellMarketWithPips()")
	fmt.Println("   ✅ Practical Example:       Complete risk-managed trade setup demonstrated")
	fmt.Println()
	fmt.Println("    You now have professional risk management tools at your fingertips!")
	fmt.Println("    Always use CalculatePositionSize() to risk a fixed % per trade.")
	fmt.Println("     Never risk more than 1-2% of your account per trade.")
	fmt.Println()
}

// ══════════════════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS FOR PRETTY PRINTING
// ══════════════════════════════════════════════════════════════════════════════

func printHeader(text string) {
	fmt.Println()
	fmt.Println("══════════════════════════════════════════════════════════════════════════════")
	fmt.Println(" " + text)
	fmt.Println("══════════════════════════════════════════════════════════════════════════════")
	fmt.Println()
}

func printSection(title string) {
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")
	fmt.Println(" " + title)
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")
	fmt.Println()
	time.Sleep(300 * time.Millisecond) // Small delay for readability
}

func printSuccess(operation, result string) {
	fmt.Printf("   ✅ %s", operation)
	if result != "" {
		fmt.Printf(" → %s", result)
	}
	fmt.Println()
}

func printError(operation string, err error) {
	fmt.Printf("   ❌ %s: %v\n", operation, err)
	fmt.Println()
}

func printWarning(operation, message string) {
	fmt.Printf("   ⚠️  %s: %s\n", operation, message)
	fmt.Println()
}

func printInfo(title, message string) {
	fmt.Printf("   ℹ️  %s: %s\n", title, message)
}
