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
   • Opens 2 test positions (BUY 0.02 / SELL 0.01) and one more BUY 0.01 later
   • Demonstrates modification operations
   • Closes all positions at the end
   • Recommended to run on DEMO account first
══════════════════════════════════════════════════════════════════════════════*/

package sugar

import (
	"fmt"
	"math"
	"strings"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
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
	sugar.GetAccount().ApiKey = cfg.ApiKey
	// Disconnect however the demo ends (normal end, Fatal, panic, Ctrl+C)
	defer helpers.DisconnectOnExit(sugar)()

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

	// BID and ASK from ONE tick: two calls (GetBid, then GetAsk) can read two different ticks
	quote, err := sugar.GetPriceInfo(testSymbol)
	helpers.Fatal(err, "GetPriceInfo failed")
	bid, ask := quote.Bid, quote.Ask
	symbolInfo, err := sugar.GetSymbolInfo(testSymbol)
	helpers.Fatal(err, "GetSymbolInfo failed")
	point := symbolInfo.Point // SL/TP distances below are in points (200 points = 20 pips on 5-digit EURUSD)
	scale := math.Pow(10, float64(symbolInfo.Digits))
	price := func(p float64) float64 { return math.Round(p*scale) / scale } // round to the symbol's digits

	// failed prints an error like helpers.PrintShortError and counts it for the summary
	failures := 0
	failed := func(err error, context string) bool {
		if helpers.PrintShortError(err, context) {
			failures++
			return true
		}
		return false
	}

	// Open 2 positions for testing (the BUY with 2x the minimum volume, so that half of it
	// can be closed in 3.1 - a partial close must itself be a valid volume)
	ticket1, err1 := sugar.BuyMarket(testSymbol, 2*minVolume)
	ticket2, err2 := sugar.SellMarket(testSymbol, minVolume)

	if err1 != nil || err2 != nil {
		fmt.Println("  ⚠️  Failed to open test positions")
		failed(err1, "BUY error")
		failed(err2, "SELL error")
		// Do not leave the one that did open behind
		for _, opened := range []struct {
			ticket uint64
			err    error
		}{{ticket1, err1}, {ticket2, err2}} {
			if opened.err == nil {
				failed(sugar.ClosePosition(opened.ticket), "Failed to close position")
			}
		}
		return
	}

	fmt.Printf("  ✓ Opened BUY position: #%d (%.2f lots)\n", ticket1, 2*minVolume)
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
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.1. HasOpenPosition() - Check if %s positions exist\n", testSymbol)

	hasPos, err := sugar.HasOpenPosition(testSymbol)
	if !failed(err, "HasOpenPosition failed") {
		fmt.Printf("  Has open %s positions: %v\n", testSymbol, hasPos)
	}

	// ══════════════════════════════════════════════════════════════
	// 1.2. CountOpenPositions()
	//      Count total number of open positions.
	//      Chain: Sugar → Service.GetPositionsTotal() → Account.PositionsTotal() → gRPC
	//      Returns: int - total count of open positions.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.2. CountOpenPositions() - Count total open positions")

	count, err := sugar.CountOpenPositions()
	if !failed(err, "CountOpenPositions failed") {
		fmt.Printf("  Total open positions: %d\n", count)
	}

	// ══════════════════════════════════════════════════════════════
	// 1.3. GetOpenPositions()
	//      Get all open positions with full details.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: []*pb.PositionInfo - slice of all open positions.
	//      SAFE operation - read-only query.
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.3. GetOpenPositions() - Get all open positions")

	positions, err := sugar.GetOpenPositions()
	if !failed(err, "GetOpenPositions failed") {
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
	//      Returns: *pb.PositionInfo, or an error if no open position has that ticket.
	//      SAFE operation - read-only query.
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.4. GetPositionByTicket() - Get specific position #%d\n", ticket1)

	pos, err := sugar.GetPositionByTicket(ticket1)
	if failed(err, "GetPositionByTicket failed") {
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
	//      Returns: []*pb.PositionInfo - positions for the symbol.
	//      SAFE operation - read-only query.
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.5. GetPositionsBySymbol() - Get positions for %s\n", testSymbol)

	symbolPositions, err := sugar.GetPositionsBySymbol(testSymbol)
	if !failed(err, "GetPositionsBySymbol failed") {
		fmt.Printf("  Found %d position(s) for %s\n", len(symbolPositions), testSymbol)
	}

	// ══════════════════════════════════════════════════════════════
	// 1.6. GetTotalProfit()
	//      Get total floating P/L from all positions.
	//      Chain: Sugar → Service.GetOpenedOrders() → Account → gRPC
	//      Returns: float64 - sum of all position profits.
	//      SAFE operation - read-only query.
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.6. GetTotalProfit() - Total floating P&L")

	totalProfit, err := sugar.GetTotalProfit()
	if !failed(err, "GetTotalProfit failed") {
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
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.7. GetProfitBySymbol() - Profit for %s\n", testSymbol)

	symbolProfit, err := sugar.GetProfitBySymbol(testSymbol)
	if !failed(err, "GetProfitBySymbol failed") {
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
	//      Modify Stop Loss level for existing position (the Take Profit is not sent and stays).
	//      Chain: Sugar → Service.ModifyOrder() → Account.OrderModify() → gRPC
	//      Returns: error if modification fails.
	//      DANGEROUS operation - modifies real position!
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	newSL := price(bid - 300*point) // 300 points below BID
	fmt.Printf("\n2.1. ModifyPositionSL() - Set Stop Loss to %.5f (BID -300 points)\n", newSL)

	err = sugar.ModifyPositionSL(ticket1, newSL)
	if !failed(err, "ModifyPositionSL failed") {
		fmt.Printf("  ✓ Stop Loss updated for position #%d\n", ticket1)
	}

	time.Sleep(500 * time.Millisecond)

	// ══════════════════════════════════════════════════════════════
	// 2.2. ModifyPositionTP()
	//      Modify Take Profit level for existing position (the Stop Loss is not sent and stays).
	//      Chain: Sugar → Service.ModifyOrder() → Account.OrderModify() → gRPC
	//      Returns: error if modification fails.
	//      DANGEROUS operation - modifies real position!
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	newTP := price(bid + 500*point) // 500 points above BID
	fmt.Printf("\n2.2. ModifyPositionTP() - Set Take Profit to %.5f (BID +500 points)\n", newTP)

	err = sugar.ModifyPositionTP(ticket1, newTP)
	if !failed(err, "ModifyPositionTP failed") {
		fmt.Printf("  ✓ Take Profit updated for position #%d\n", ticket1)
		if p, err := sugar.GetPositionByTicket(ticket1); !failed(err, "GetPositionByTicket failed") {
			fmt.Printf("    Position #%d now: SL %.5f, TP %.5f\n", ticket1, p.StopLoss, p.TakeProfit)
		}
	}

	time.Sleep(500 * time.Millisecond)

	// ══════════════════════════════════════════════════════════════
	// 2.3. ModifyPositionSLTP()
	//      Modify both Stop Loss and Take Profit in one call.
	//      Chain: Sugar → Service.ModifyOrder() → Account.OrderModify() → gRPC
	//      Returns: error if modification fails.
	//      DANGEROUS operation - modifies real position!
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	newSL2 := price(ask + 300*point) // For SELL position: SL 300 points above ASK
	newTP2 := price(ask - 500*point) // TP 500 points below ASK
	fmt.Printf("\n2.3. ModifyPositionSLTP() - Set both SL and TP\n")
	fmt.Printf("  SL: %.5f, TP: %.5f\n", newSL2, newTP2)
	err = sugar.ModifyPositionSLTP(ticket2, newSL2, newTP2)
	if !failed(err, "ModifyPositionSLTP failed") {
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
	//      Close part of position volume (partial close). The closed volume must be
	//      a valid volume itself (>= minimum, multiple of the volume step).
	//      Chain: Sugar → Service.CloseOrder() → Account.OrderClose() → gRPC
	//      Returns: error if close fails.
	//      DANGEROUS operation - closes real position partially!
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.1. ClosePositionPartial() - Close 50%% of position #%d\n", ticket1)

	partialVolume := minVolume // half of the 2x minimum volume opened above
	err = sugar.ClosePositionPartial(ticket1, partialVolume)
	if failed(err, "ClosePositionPartial failed") {
		fmt.Println("  ℹ️  Partial close may not be supported on all brokers")
	} else {
		fmt.Printf("  ✓ Closed %.2f lots from position #%d\n", partialVolume, ticket1)
	}

	time.Sleep(1 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 3.2. ClosePosition()
	//      Close entire position completely.
	//      Chain: Sugar → Service.CloseOrder() → Account.OrderClose() → gRPC
	//      Returns: error if close fails.
	//      DANGEROUS operation - closes real position!
	//      Timeout: 30 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.2. ClosePosition() - Close entire position #%d\n", ticket1)

	err = sugar.ClosePosition(ticket1)
	if !failed(err, "ClosePosition failed") {
		fmt.Printf("  ✓ Position #%d closed completely\n", ticket1)
	}

	time.Sleep(1 * time.Second)

	// Open a new position for symbol-specific close test
	ticket3, err := sugar.BuyMarket(testSymbol, minVolume)
	if !failed(err, "BuyMarket for the CloseAllBySymbol demo failed") {
		fmt.Printf("\n  Opened test position #%d for CloseAllBySymbol demo\n", ticket3)
		time.Sleep(1 * time.Second)

		// ══════════════════════════════════════════════════════════════
		// 3.3. CloseAllBySymbol()
		//      Close all positions for specific symbol.
		//      Chain: Sugar → Service (GetOpenedOrders + CloseOrder loop) → Account
		//      Returns: int - number of closed positions, and an error naming
		//      every position that could not be closed.
		//      DANGEROUS operation - closes all symbol positions!
		//      Timeout: 10 s to list, 30 s per close
		// ══════════════════════════════════════════════════════════════
		fmt.Printf("\n3.3. CloseAllBySymbol() - Close all %s positions\n", testSymbol)
		
		closedCount, err := sugar.CloseAllBySymbol(testSymbol)
		if !failed(err, "CloseAllBySymbol failed") {
			fmt.Printf("  ✓ Closed %d position(s) for %s\n", closedCount, testSymbol)
		}
	}

	time.Sleep(1 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 3.4. CloseAllPositions()
	//      Close ALL open positions (all symbols).
	//      Chain: Sugar → Service (GetOpenedOrders + CloseOrder loop) → Account
	//      Returns: int - number of closed positions, and an error naming
	//      every position that could not be closed.
	//      DANGEROUS operation - closes ALL positions!
	//      Timeout: 10 s to list, 30 s per close
	// ══════════════════════════════════════════════════════════════
	// First, check if any positions remain
	remainingCount, err := sugar.CountOpenPositions()
	failed(err, "CountOpenPositions failed")
	if err != nil || remainingCount > 0 {
		fmt.Printf("\n3.4. CloseAllPositions() - Close all remaining positions\n")
		closedTotal, err := sugar.CloseAllPositions()
		if !failed(err, "CloseAllPositions failed") {
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
	finalCount, err := sugar.CountOpenPositions()
	if failed(err, "CountOpenPositions failed") {
		finalCount = -1
	} else {
		fmt.Printf("\nFinal open positions count: %d\n", finalCount)
	}

	if finalCount == 0 {
		fmt.Println("✅ All test positions successfully closed!")
	} else if finalCount > 0 {
		fmt.Printf("⚠️  %d position(s) still open\n", finalCount)
	}

	// =============================
	// SUMMARY
	// =============================
	fmt.Println("\n" + strings.Repeat("=", 80))
	if failures == 0 {
		fmt.Println("DEMO COMPLETED SUCCESSFULLY!")
	} else {
		fmt.Printf("DEMO FINISHED - %d operation(s) failed (see the errors above)\n", failures)
	}
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nℹ️  MT5Sugar Position Management Features:")
	fmt.Println("  • Position info: GetOpenPositions, GetPositionByTicket, CountOpenPositions")
	fmt.Println("  • Filtering: GetPositionsBySymbol, GetProfitBySymbol")
	fmt.Println("  • Modification: ModifyPositionSL, ModifyPositionTP, ModifyPositionSLTP")
	fmt.Println("  • Closing: ClosePosition, ClosePositionPartial, CloseAllBySymbol, CloseAllPositions")
	fmt.Println("  • Quick checks: HasOpenPosition, GetTotalProfit")
	fmt.Println("\n✅ All position management operations demonstrated!")
	fmt.Println("   Next: Run demo 09 - sugar_history.go for historical data")
}
