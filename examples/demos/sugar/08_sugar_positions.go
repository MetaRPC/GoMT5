/*══════════════════════════════════════════════════════════════════════════════
 FILE: 08_sugar_positions.go - HIGH-LEVEL SUGAR API POSITION MANAGEMENT DEMO

 PURPOSE:
   Demonstrates MT5Sugar position management: querying, modifying, and closing.
   SIMPLIFIED position operations without complex parameter handling!

 🎯 WHO SHOULD USE THIS:
   • Managing open positions programmatically
   • Building position monitoring tools
   • Learning position modification (SL/TP)
   • Quick position closing operations

 📚 WHAT THIS DEMO COVERS (3 Categories):

   1. POSITION INFO METHODS (7 methods)
      • GetOpenPositions() - All open positions
      • GetPositionByTicket(ticket) - Specific position
      • GetPositionsBySymbol(symbol) - Positions for symbol
      • HasOpenPosition(symbol) - Check if positions exist
      • CountOpenPositions() - Count total positions
      • GetTotalProfit() - Total floating P&L
      • GetProfitBySymbol(symbol) - Profit for symbol

   2. POSITION MODIFICATION (3 methods)
      • ModifyPositionSL(ticket, sl) - Change Stop Loss
      • ModifyPositionTP(ticket, tp) - Change Take Profit
      • ModifyPositionSLTP(ticket, sl, tp) - Change both SL/TP

   3. POSITION CLOSING (4 methods)
      • ClosePosition(ticket) - Close full position
      • ClosePositionPartial(ticket, volume) - Partial close
      • CloseAllPositions() - Close all positions
      • CloseAllBySymbol(symbol) - Close all for symbol

 🔄 API LEVELS:
   HIGH-LEVEL (Sugar) - THIS FILE: Simple position operations
   MID-LEVEL (Service): More control, custom close parameters
   LOW-LEVEL (Account): Full control, manual request construction

 🚀 HOW TO RUN:
   cd examples/demos
   go run main.go 8    (or select [8] from menu)

 ⚠️  WARNING: This demo OPENS and CLOSES test positions!
   • Opens 2 test positions (BUY/SELL)
   • Demonstrates modification operations
   • Closes all positions at the end
   • Recommended to run on DEMO account first
══════════════════════════════════════════════════════════════════════════════*/

package sugar

import (
	"fmt"
	"strings"
	"time"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	"github.com/MetaRPC/GoMT5/mt5"
)

// RunSugarPositionsDemo demonstrates MT5Sugar position management:
// - Position Management (ClosePosition, ModifyPosition, CloseAll, etc.)
// - Position Info (GetOpenPositions, GetPositionByTicket, HasOpenPosition, etc.)
func RunSugarPositionsDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("MT5 SUGAR API - POSITION MANAGEMENT DEMO")
	fmt.Println(strings.Repeat("=", 80))

	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// Connect
	fmt.Println("\n📡 Connecting to MT5...")
	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	helpers.Fatal(err, "Failed to create Sugar instance")

	err = sugar.QuickConnect(cfg.MtCluster)
	helpers.Fatal(err, "Connection failed")
	fmt.Println("  ✓ Connected!")

	testSymbol := cfg.TestSymbol
	minVolume := 0.01

	// =============================
	// PREPARATION: Open test positions
	// =============================
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("PREPARATION: Opening test positions")
	fmt.Println(strings.Repeat("-", 80))

	bid, _ := sugar.GetBid(testSymbol)
	ask, _ := sugar.GetAsk(testSymbol)

	// Open 2 positions for testing
	ticket1, err1 := sugar.BuyMarket(testSymbol, minVolume)
	ticket2, err2 := sugar.SellMarket(testSymbol, minVolume)

	if err1 != nil || err2 != nil {
		fmt.Println("  ⚠️  Failed to open test positions")
		helpers.PrintShortError(err1, "BUY error")
		helpers.PrintShortError(err2, "SELL error")
		return
	}

	fmt.Printf("  ✓ Opened BUY position: #%d\n", ticket1)
	fmt.Printf("  ✓ Opened SELL position: #%d\n", ticket2)
	time.Sleep(1 * time.Second) // Let positions register

	// ══════════════════════════════════════════════════════════════════════════════
	// 1. POSITION INFO METHODS
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("1. POSITION INFO METHODS")
	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 1.1. HasOpenPosition()
	//      Check if any open positions exist for symbol.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: bool - true if positions exist.
	//      SAFE operation - read-only query.
	//      Timeout: 5 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.1. HasOpenPosition() - Check if %s positions exist\n", testSymbol)

	hasPos, err := sugar.HasOpenPosition(testSymbol)
	if !helpers.PrintShortError(err, "HasOpenPosition failed") {
		fmt.Printf("  Has open %s positions: %v\n", testSymbol, hasPos)
	}

	// ══════════════════════════════════════════════════════════════
	// 1.2. CountOpenPositions()
	//      Count total number of open positions.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: int - total count of open positions.
	//      SAFE operation - read-only query.
	//      Timeout: 5 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.2. CountOpenPositions() - Count total open positions")

	count, err := sugar.CountOpenPositions()
	if !helpers.PrintShortError(err, "CountOpenPositions failed") {
		fmt.Printf("  Total open positions: %d\n", count)
	}

	// ══════════════════════════════════════════════════════════════
	// 1.3. GetOpenPositions()
	//      Get all open positions with full details.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: []*Position - slice of all open positions.
	//      SAFE operation - read-only query.
	//      Timeout: 5 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.3. GetOpenPositions() - Get all open positions")

	positions, err := sugar.GetOpenPositions()
	if !helpers.PrintShortError(err, "GetOpenPositions failed") {
		fmt.Printf("  Found %d position(s):\n", len(positions))
		for i, pos := range positions {
			posType := "BUY"
			if pos.Type == pb.BMT5_ENUM_POSITION_TYPE_BMT5_POSITION_TYPE_SELL {
				posType = "SELL"
			}
			fmt.Printf("    %d. Ticket #%d: %s %.2f lots, Profit: %.2f\n",
				i+1, pos.Ticket, posType, pos.Volume, pos.Profit)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 1.4. GetPositionByTicket()
	//      Get specific position by ticket number.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: *Position - position details or nil if not found.
	//      SAFE operation - read-only query.
	//      Timeout: 5 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.4. GetPositionByTicket() - Get specific position #%d\n", ticket1)

	pos, err := sugar.GetPositionByTicket(ticket1)
	if helpers.PrintShortError(err, "GetPositionByTicket failed") {
		// Error already printed
	} else if pos == nil {
		fmt.Println("  ⚠️  Position not found")
	} else {
		fmt.Printf("  ✓ Found position:\n")
		fmt.Printf("    Symbol: %s\n", pos.Symbol)
		fmt.Printf("    Volume: %.2f\n", pos.Volume)
		fmt.Printf("    Price:  %.5f\n", pos.PriceOpen)
		fmt.Printf("    Profit: %.2f\n", pos.Profit)
	}

	// ══════════════════════════════════════════════════════════════
	// 1.5. GetPositionsBySymbol()
	//      Get all positions for specific symbol.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: []*Position - positions for the symbol.
	//      SAFE operation - read-only query.
	//      Timeout: 5 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.5. GetPositionsBySymbol() - Get positions for %s\n", testSymbol)

	symbolPositions, err := sugar.GetPositionsBySymbol(testSymbol)
	if !helpers.PrintShortError(err, "GetPositionsBySymbol failed") {
		fmt.Printf("  Found %d position(s) for %s\n", len(symbolPositions), testSymbol)
	}

	// ══════════════════════════════════════════════════════════════
	// 1.6. GetTotalProfit()
	//      Get total floating P/L from all positions.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: float64 - sum of all position profits.
	//      SAFE operation - read-only query.
	//      Timeout: 5 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.6. GetTotalProfit() - Total floating P&L")

	totalProfit, err := sugar.GetTotalProfit()
	if !helpers.PrintShortError(err, "GetTotalProfit failed") {
		if totalProfit >= 0 {
			fmt.Printf("  Total profit: +%.2f\n", totalProfit)
		} else {
			fmt.Printf("  Total profit: %.2f\n", totalProfit)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 1.7. GetProfitBySymbol()
	//      Get floating P/L for specific symbol.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: float64 - sum of profits for the symbol.
	//      SAFE operation - read-only query.
	//      Timeout: 5 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.7. GetProfitBySymbol() - Profit for %s\n", testSymbol)

	symbolProfit, err := sugar.GetProfitBySymbol(testSymbol)
	if !helpers.PrintShortError(err, "GetProfitBySymbol failed") {
		if symbolProfit >= 0 {
			fmt.Printf("  %s profit: +%.2f\n", testSymbol, symbolProfit)
		} else {
			fmt.Printf("  %s profit: %.2f\n", testSymbol, symbolProfit)
		}
	}

	// ══════════════════════════════════════════════════════════════════════════════
	// 2. POSITION MODIFICATION
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("2. POSITION MODIFICATION METHODS")
	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 2.1. ModifyPositionSL()
	//      Modify Stop Loss level for existing position.
	//      Chain: Sugar → Service.ModifyOrder() → Account.OrderSend() → gRPC
	//      Returns: error if modification fails.
	//      DANGEROUS operation - modifies real position!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	newSL := bid - 0.0030 // 30 pips below
	fmt.Printf("\n2.1. ModifyPositionSL() - Set Stop Loss to %.5f\n", newSL)

	err = sugar.ModifyPositionSL(ticket1, newSL)
	if !helpers.PrintShortError(err, "ModifyPositionSL failed") {
		fmt.Printf("  ✓ Stop Loss updated for position #%d\n", ticket1)
	}

	time.Sleep(500 * time.Millisecond)

	// ══════════════════════════════════════════════════════════════
	// 2.2. ModifyPositionTP()
	//      Modify Take Profit level for existing position.
	//      Chain: Sugar → Service.ModifyOrder() → Account.OrderSend() → gRPC
	//      Returns: error if modification fails.
	//      DANGEROUS operation - modifies real position!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	newTP := bid + 0.0050 // 50 pips above
	fmt.Printf("\n2.2. ModifyPositionTP() - Set Take Profit to %.5f\n", newTP)

	err = sugar.ModifyPositionTP(ticket1, newTP)
	if !helpers.PrintShortError(err, "ModifyPositionTP failed") {
		fmt.Printf("  ✓ Take Profit updated for position #%d\n", ticket1)
	}

	time.Sleep(500 * time.Millisecond)

	// ══════════════════════════════════════════════════════════════
	// 2.3. ModifyPositionSLTP()
	//      Modify both Stop Loss and Take Profit in one call.
	//      Chain: Sugar → Service.ModifyOrder() → Account.OrderSend() → gRPC
	//      Returns: error if modification fails.
	//      DANGEROUS operation - modifies real position!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	newSL2 := ask + 0.0030  // For SELL position
	newTP2 := ask - 0.0050
	fmt.Printf("\n2.3. ModifyPositionSLTP() - Set both SL and TP\n")
	fmt.Printf("  SL: %.5f, TP: %.5f\n", newSL2, newTP2)
	err = sugar.ModifyPositionSLTP(ticket2, newSL2, newTP2)
	if !helpers.PrintShortError(err, "ModifyPositionSLTP failed") {
		fmt.Printf("  ✓ SL/TP updated for position #%d\n", ticket2)
	}

	time.Sleep(500 * time.Millisecond)

	// ══════════════════════════════════════════════════════════════════════════════
	// 3. POSITION CLOSING
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("3. POSITION CLOSING METHODS")
	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 3.1. ClosePositionPartial()
	//      Close part of position volume (partial close).
	//      Chain: Sugar → Service.CloseOrder() → Account.OrderSend() → gRPC
	//      Returns: error if close fails.
	//      DANGEROUS operation - closes real position partially!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.1. ClosePositionPartial() - Close 50%% of position #%d\n", ticket1)

	partialVolume := minVolume / 2
	err = sugar.ClosePositionPartial(ticket1, partialVolume)
	if helpers.PrintShortError(err, "ClosePositionPartial failed") {
		fmt.Println("  ℹ️  Partial close may not be supported on all brokers")
	} else {
		fmt.Printf("  ✓ Closed %.3f lots from position #%d\n", partialVolume, ticket1)
	}

	time.Sleep(1 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 3.2. ClosePosition()
	//      Close entire position completely.
	//      Chain: Sugar → Service.CloseOrder() → Account.OrderSend() → gRPC
	//      Returns: error if close fails.
	//      DANGEROUS operation - closes real position!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.2. ClosePosition() - Close entire position #%d\n", ticket1)

	err = sugar.ClosePosition(ticket1)
	if !helpers.PrintShortError(err, "ClosePosition failed") {
		fmt.Printf("  ✓ Position #%d closed completely\n", ticket1)
	}

	time.Sleep(1 * time.Second)

	// Open a new position for symbol-specific close test
	ticket3, err := sugar.BuyMarket(testSymbol, minVolume)
	if err == nil {
		fmt.Printf("\n  Opened test position #%d for CloseAllBySymbol demo\n", ticket3)
		time.Sleep(1 * time.Second)

		// ══════════════════════════════════════════════════════════════
		// 3.3. CloseAllBySymbol()
		//      Close all positions for specific symbol.
		//      Chain: Sugar → Service (GetOpenedOrders + CloseOrder loop) → Account
		//      Returns: int - number of closed positions.
		//      DANGEROUS operation - closes all symbol positions!
		//      Timeout: 30 seconds
		// ══════════════════════════════════════════════════════════════
		fmt.Printf("\n3.3. CloseAllBySymbol() - Close all %s positions\n", testSymbol)
		
		closedCount, err := sugar.CloseAllBySymbol(testSymbol)
		if !helpers.PrintShortError(err, "CloseAllBySymbol failed") {
			fmt.Printf("  ✓ Closed %d position(s) for %s\n", closedCount, testSymbol)
		}
	}

	time.Sleep(1 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 3.4. CloseAllPositions()
	//      Close ALL open positions (all symbols).
	//      Chain: Sugar → Service (GetOpenedOrders + CloseOrder loop) → Account
	//      Returns: int - number of closed positions.
	//      DANGEROUS operation - closes ALL positions!
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	// First, check if any positions remain
	remainingCount, _ := sugar.CountOpenPositions()
	if remainingCount > 0 {
		fmt.Printf("\n3.4. CloseAllPositions() - Close all remaining positions\n")
		closedTotal, err := sugar.CloseAllPositions()
		if !helpers.PrintShortError(err, "CloseAllPositions failed") {
			fmt.Printf("  ✓ Closed %d position(s)\n", closedTotal)
		}
	} else {
		fmt.Println("\n3.4. CloseAllPositions() - No positions to close")
		fmt.Println("  ℹ️  All positions already closed")
	}

	// =============================
	// FINAL CHECK
	// =============================
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("FINAL STATUS CHECK")
	fmt.Println(strings.Repeat("-", 80))

	time.Sleep(1 * time.Second)
	finalCount, _ := sugar.CountOpenPositions()
	fmt.Printf("\nFinal open positions count: %d\n", finalCount)

	if finalCount == 0 {
		fmt.Println("✅ All test positions successfully closed!")
	} else {
		fmt.Printf("⚠️  %d position(s) still open\n", finalCount)
	}

	// =============================
	// SUMMARY
	// =============================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("DEMO COMPLETED SUCCESSFULLY!")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nℹ️  MT5Sugar Position Management Features:")
	fmt.Println("  • Position info: GetOpenPositions, GetPositionByTicket, CountOpenPositions")
	fmt.Println("  • Filtering: GetPositionsBySymbol, GetProfitBySymbol")
	fmt.Println("  • Modification: ModifyPositionSL, ModifyPositionTP, ModifyPositionSLTP")
	fmt.Println("  • Closing: ClosePosition, ClosePositionPartial, CloseAllBySymbol, CloseAllPositions")
	fmt.Println("  • Quick checks: HasOpenPosition, GetTotalProfit")
	fmt.Println("\n✅ All position management operations demonstrated!")
	fmt.Println("   Next: Run 08_sugar_history.go for historical data")
}
