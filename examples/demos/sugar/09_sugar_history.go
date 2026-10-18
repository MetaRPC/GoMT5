/*══════════════════════════════════════════════════════════════════════════════
 FILE: 09_sugar_history.go - HIGH-LEVEL SUGAR API HISTORY & PROFITS DEMO

 PURPOSE:
   Demonstrates MT5Sugar historical data methods: deals and profit calculations.
   QUICK access to trading history with pre-defined time ranges!

 🎯 WHO SHOULD USE THIS:
   • Daily/weekly/monthly performance tracking
   • Building trading statistics and reports
   • Analyzing historical performance
   • Quick profit/loss calculations

 📚 WHAT THIS DEMO COVERS (2 Categories):

   1. HISTORICAL DEALS METHODS (5 methods)
      "Deals" are CLOSED POSITIONS selected by their OPEN time. Day/week/month
      use the local calendar of this computer; the window is sent to the server
      as UTC, so the broker's time zone does not shift it.
      • GetDealsToday() - Closed positions opened today (local midnight to now)
      • GetDealsYesterday() - ... opened yesterday
      • GetDealsThisWeek() - ... opened this week (Monday 00:00 to now)
      • GetDealsThisMonth() - ... opened this month
      • GetDealsDateRange(from, to) - ... opened in a custom range

   2. PROFIT CALCULATION METHODS (3 methods)
      Sum of the positions' Profit field (commission, swap and fee not included)
      • GetProfitToday() - Today's total profit
      • GetProfitThisWeek() - This week's profit
      • GetProfitThisMonth() - This month's profit

 💡 BONUS FEATURES DEMONSTRATED:
   • Win/Loss ratio calculation
   • Trade statistics (count, volume, profitability)
   • Sample deals display with formatting

 🔄 API LEVELS:
   HIGH-LEVEL (Sugar) - THIS FILE: Pre-defined time ranges, automatic calculations
   MID-LEVEL (Service): Custom time ranges, manual filtering
   LOW-LEVEL (Account): Full control, protobuf timestamps

 🚀 HOW TO RUN:
   cd examples/demos
   go run main.go 9    (or select [9] from menu)

 ⚠️  WARNING: This demo CREATES TEST HISTORY!
   • Opens and closes 2 test positions (BUY/SELL)
   • Generates sample deals for demonstration
   • Recommended to run on DEMO account first
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

// RunSugarHistoryDemo demonstrates MT5Sugar history methods:
// - GetDealsToday, GetDealsYesterday, GetDealsThisWeek, GetDealsThisMonth
// - GetDealsDateRange
// - GetProfitToday, GetProfitThisWeek, GetProfitThisMonth
func RunSugarHistoryDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("MT5 SUGAR API - HISTORY DEMO (Historical Deals & Profits)")
	fmt.Println(strings.Repeat("=", 80))

	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// Connect
	fmt.Println("\n📡 Connecting to MT5...")
	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	helpers.Fatal(err, "Failed to create Sugar instance")
	sugar.GetAccount().ApiKey = cfg.ApiKey
	// Disconnect however the demo ends (normal end, Fatal, panic, Ctrl+C)
	defer helpers.DisconnectOnExit(sugar)()

	err = sugar.QuickConnect(cfg.MtCluster)
	helpers.Fatal(err, "Connection failed")
	fmt.Println("  ✓ Connected!")

	testSymbol := cfg.TestSymbol
	minVolume := 0.01

	// =============================
	// PREPARATION: Create history
	// =============================
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("PREPARATION: Creating test trade history")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\nExecuting test trades to generate history...")

	closedTrades := 0

	// Trade 1: Quick BUY/CLOSE
	ticket1, err1 := sugar.BuyMarket(testSymbol, minVolume)
	if !helpers.PrintShortError(err1, "Failed to open BUY") {
		fmt.Printf("  ✓ Opened BUY #%d\n", ticket1)
		time.Sleep(2 * time.Second)
		closeErr := sugar.ClosePosition(ticket1)
		if !helpers.PrintShortError(closeErr, "Failed to close position") {
			fmt.Printf("  ✓ Closed BUY #%d\n", ticket1)
			closedTrades++
		}
	}

	time.Sleep(1 * time.Second)

	// Trade 2: Quick SELL/CLOSE
	ticket2, err2 := sugar.SellMarket(testSymbol, minVolume)
	if !helpers.PrintShortError(err2, "Failed to open SELL") {
		fmt.Printf("  ✓ Opened SELL #%d\n", ticket2)
		time.Sleep(2 * time.Second)
		closeErr := sugar.ClosePosition(ticket2)
		if !helpers.PrintShortError(closeErr, "Failed to close position") {
			fmt.Printf("  ✓ Closed SELL #%d\n", ticket2)
			closedTrades++
		}
	}

	if closedTrades > 0 {
		fmt.Printf("\n  ✓ Test history created (%d closed position(s))\n", closedTrades)
	} else {
		fmt.Println("\n  ⚠️  No test history created - the queries below show the existing history only")
	}
	time.Sleep(2 * time.Second) // Let server register deals

	// ══════════════════════════════════════════════════════════════════════════════
	// 1. HISTORICAL DEALS METHODS
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("1. HISTORICAL DEALS METHODS")
	fmt.Println(strings.Repeat("-", 80))

	// ══════════════════════════════════════════════════════════════
	// 1.1. GetDealsToday()
	//      Get the closed positions opened today (local midnight to now).
	//      Chain: Sugar → Service.GetPositionsHistory() (UTC window) → Account → gRPC
	//      Returns: []*PositionHistoryInfo - slice of today's closed positions.
	//      SAFE operation - read-only query.
	//      Timeout: 15 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.1. GetDealsToday() - Closed positions opened today (local midnight to now)")

	dealsToday, err := sugar.GetDealsToday()
	if !helpers.PrintShortError(err, "GetDealsToday failed") {
		fmt.Printf("  Found %d closed position(s) opened today:\n", len(dealsToday))
		for i, deal := range dealsToday {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", len(dealsToday)-5)
				break
			}
			fmt.Printf("    %d. Position #%d: %.2f lots, Profit: %.2f, opened %s (broker time) = %s UTC\n",
				i+1, deal.PositionTicket, deal.Volume, deal.Profit,
				deal.OpenTime.AsTime().Format("2006-01-02 15:04:05"),
				deal.OpenTimeUtc.AsTime().Format("2006-01-02 15:04:05"))
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 1.2. GetDealsYesterday()
	//      Get the closed positions opened yesterday (full local day: 00:00 to 23:59:59).
	//      Chain: Sugar → Service.GetPositionsHistory() (UTC window) → Account → gRPC
	//      Returns: []*PositionHistoryInfo - slice of yesterday's closed positions.
	//      SAFE operation - read-only query.
	//      Timeout: 15 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.2. GetDealsYesterday() - Closed positions opened yesterday")

	dealsYesterday, err := sugar.GetDealsYesterday()
	if !helpers.PrintShortError(err, "GetDealsYesterday failed") {
		if len(dealsYesterday) == 0 {
			fmt.Println("  None opened yesterday")
		} else {
			fmt.Printf("  Found %d closed position(s) opened yesterday\n", len(dealsYesterday))
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 1.3. GetDealsThisWeek()
	//      Get the closed positions opened this week (local Monday 00:00 to now).
	//      Chain: Sugar → Service.GetPositionsHistory() (UTC window) → Account → gRPC
	//      Returns: []*PositionHistoryInfo - slice of this week's closed positions.
	//      SAFE operation - read-only query.
	//      Timeout: 15 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.3. GetDealsThisWeek() - Closed positions opened this week")

	dealsWeek, err := sugar.GetDealsThisWeek()
	if !helpers.PrintShortError(err, "GetDealsThisWeek failed") {
		fmt.Printf("  Found %d closed position(s) opened this week\n", len(dealsWeek))
	}

	// ══════════════════════════════════════════════════════════════
	// 1.4. GetDealsThisMonth()
	//      Get the closed positions opened this month (the 1st, 00:00 local to now).
	//      Chain: Sugar → Service.GetPositionsHistory() (UTC window) → Account → gRPC
	//      Returns: []*PositionHistoryInfo - slice of this month's closed positions.
	//      SAFE operation - read-only query.
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.4. GetDealsThisMonth() - Closed positions opened this month")

	dealsMonth, err := sugar.GetDealsThisMonth()
	if !helpers.PrintShortError(err, "GetDealsThisMonth failed") {
		fmt.Printf("  Found %d closed position(s) opened this month\n", len(dealsMonth))
	}

	// ══════════════════════════════════════════════════════════════
	// 1.5. GetDealsDateRange()
	//      Get the closed positions opened within a custom range (from - to: real moments,
	//      sent to the server as UTC).
	//      Chain: Sugar → Service.GetPositionsHistory() (UTC window) → Account → gRPC
	//      Returns: []*PositionHistoryInfo - slice of closed positions in the range.
	//      SAFE operation - read-only query.
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.5. GetDealsDateRange() - Custom range (last 7 days)")

	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	dealsRange, err := sugar.GetDealsDateRange(weekAgo, now)
	if !helpers.PrintShortError(err, "GetDealsDateRange failed") {
		fmt.Printf("  Found %d closed position(s) opened in the last 7 days\n", len(dealsRange))
		fmt.Printf("  Range: %s to %s UTC\n",
			weekAgo.UTC().Format("2006-01-02 15:04"),
			now.UTC().Format("2006-01-02 15:04"))
	}

	// ══════════════════════════════════════════════════════════════════════════════
	// 2. PROFIT CALCULATION METHODS
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("2. PROFIT CALCULATION METHODS")
	fmt.Println(strings.Repeat("-", 80))

	// ══════════════════════════════════════════════════════════════
	// 2.1. GetProfitToday()
	//      Sum of Profit of the closed positions opened today (no commission/swap/fee).
	//      Chain: Sugar → GetDealsToday() → Service → Account → gRPC
	//      Returns: float64 - sum of profit of today's closed positions.
	//      SAFE operation - read-only query.
	//      Timeout: 15 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.1. GetProfitToday() - Total profit of today's closed positions")

	profitToday, err := sugar.GetProfitToday()
	if !helpers.PrintShortError(err, "GetProfitToday failed") {
		if profitToday >= 0 {
			fmt.Printf("  Today's profit: +%.2f\n", profitToday)
		} else {
			fmt.Printf("  Today's profit: %.2f\n", profitToday)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 2.2. GetProfitThisWeek()
	//      Sum of Profit of the closed positions opened this week (no commission/swap/fee).
	//      Chain: Sugar → GetDealsThisWeek() → Service → Account → gRPC
	//      Returns: float64 - sum of profit of this week's closed positions.
	//      SAFE operation - read-only query.
	//      Timeout: 15 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.2. GetProfitThisWeek() - Total profit from this week")

	profitWeek, err := sugar.GetProfitThisWeek()
	if !helpers.PrintShortError(err, "GetProfitThisWeek failed") {
		if profitWeek >= 0 {
			fmt.Printf("  This week's profit: +%.2f\n", profitWeek)
		} else {
			fmt.Printf("  This week's profit: %.2f\n", profitWeek)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 2.3. GetProfitThisMonth()
	//      Sum of Profit of the closed positions opened this month (no commission/swap/fee).
	//      Chain: Sugar → GetDealsThisMonth() → Service → Account → gRPC
	//      Returns: float64 - sum of profit of this month's closed positions.
	//      SAFE operation - read-only query.
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.3. GetProfitThisMonth() - Total profit from this month")
	
	profitMonth, err := sugar.GetProfitThisMonth()
	if !helpers.PrintShortError(err, "GetProfitThisMonth failed") {
		if profitMonth >= 0 {
			fmt.Printf("  This month's profit: +%.2f\n", profitMonth)
		} else {
			fmt.Printf("  This month's profit: %.2f\n", profitMonth)
		}
	}

	// =============================
	// 3. DETAILED DEALS ANALYSIS
	// =============================
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("3. DETAILED DEALS ANALYSIS (Today)")
	fmt.Println(strings.Repeat("-", 80))

	if len(dealsToday) > 0 {
		fmt.Printf("\nAnalyzing %d closed position(s) opened today:\n", len(dealsToday))

		var totalProfit float64
		var totalVolume float64
		winCount := 0
		lossCount := 0

		for _, deal := range dealsToday {
			totalProfit += deal.Profit
			totalVolume += deal.Volume

			if deal.Profit > 0 {
				winCount++
			} else if deal.Profit < 0 {
				lossCount++
			}
		}

		fmt.Println("\n📊 Statistics:")
		fmt.Printf("  Positions:    %d\n", len(dealsToday))
		fmt.Printf("  Winning:      %d (%.1f%%)\n", winCount, float64(winCount)/float64(len(dealsToday))*100)
		fmt.Printf("  Losing:       %d (%.1f%%)\n", lossCount, float64(lossCount)/float64(len(dealsToday))*100)
		fmt.Printf("  Total volume: %.2f lots\n", totalVolume)

		if totalProfit >= 0 {
			fmt.Printf("  Profit:       +%.2f ✅ (without commission/swap/fee)\n", totalProfit)
		} else {
			fmt.Printf("  Profit:       %.2f ❌ (without commission/swap/fee)\n", totalProfit)
		}

		// Show sample deals
		fmt.Println("\n📋 Sample closed positions (max 10):")
		for i, deal := range dealsToday {
			if i >= 10 {
				fmt.Printf("  ... and %d more\n", len(dealsToday)-10)
				break
			}

			profitStr := fmt.Sprintf("%.2f", deal.Profit)
			if deal.Profit > 0 {
				profitStr = "+" + profitStr
			}

			fmt.Printf("  %2d. Position #%d | %.2f lots | %s\n",
				i+1, deal.PositionTicket, deal.Volume, profitStr)
		}
	} else {
		fmt.Println("\nℹ️  No closed positions opened today to analyze")
		fmt.Println("   (The test trades may take a few seconds to appear in history)")
	}

	// =============================
	// SUMMARY
	// =============================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("DEMO COMPLETED SUCCESSFULLY!")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nℹ️  MT5Sugar History Features:")
	fmt.Println("  • Time-based queries: GetDealsToday, GetDealsYesterday, GetDealsThisWeek, GetDealsThisMonth")
	fmt.Println("  • Custom ranges: GetDealsDateRange(from, to)")
	fmt.Println("  • Profit calculations: GetProfitToday, GetProfitThisWeek, GetProfitThisMonth")
	fmt.Println("  • Returns: []*pb.PositionHistoryInfo (one entry per closed position)")
	fmt.Println("\n✅ All history operations demonstrated!")
	fmt.Println("\n📚 SUGAR API COMPLETE!")
	fmt.Println("   You've now seen all 4 demo categories:")
	fmt.Println("   • 06_sugar_basics.go - Connection, Balance, Prices (read-only)")
	fmt.Println("   • 07_sugar_trading.go - Market & Pending Orders")
	fmt.Println("   • 08_sugar_positions.go - Position Management & Info")
	fmt.Println("   • 09_sugar_history.go - Historical Deals & Profits ✅")
}
