/*══════════════════════════════════════════════════════════════════════════════
 FILE: 07_sugar_trading.go - HIGH-LEVEL SUGAR API TRADING DEMO

 PURPOSE:
   Demonstrates MT5Sugar trading operations: market and pending orders.
   ULTRA-SIMPLE one-liner trading methods for quick order execution!

 🎯 WHO SHOULD USE THIS:
   • Quick trading scripts and bots
   • Simple order execution without complexity
   • Learning basic trading operations
   • Prototyping trading strategies

 📚 WHAT THIS DEMO COVERS (3 Categories):

   1. SIMPLE MARKET TRADING (2 methods)
      • BuyMarket(symbol, volume) - Instant BUY
      • SellMarket(symbol, volume) - Instant SELL

   2. PENDING ORDERS (4 methods)
      • BuyLimit(symbol, volume, price) - BUY when price drops
      • SellLimit(symbol, volume, price) - SELL when price rises
      • BuyStop(symbol, volume, price) - BUY when price rises
      • SellStop(symbol, volume, price) - SELL when price drops

   3. TRADING WITH SL/TP (4 methods)
      • BuyMarketWithSLTP() - Market BUY with Stop Loss & Take Profit
      • SellMarketWithSLTP() - Market SELL with SL/TP
      • BuyLimitWithSLTP() - Limit BUY with SL/TP
      • SellLimitWithSLTP() - Limit SELL with SL/TP

 🔄 API LEVELS:
   HIGH-LEVEL (Sugar) - THIS FILE: One-liner trading, automatic defaults
   MID-LEVEL (Service): More parameters, custom settings
   LOW-LEVEL (Account): Full control, manual order construction

 🚀 HOW TO RUN:
   cd examples/demos
   go run main.go 7    (or select [7] from menu)

 ⚠️  WARNING: This demo EXECUTES REAL TRADES on your account!
   • Uses minimum lot size (0.01)
   • Immediately closes market positions
   • Cancels pending orders at the end
   • Recommended to run on DEMO account first
══════════════════════════════════════════════════════════════════════════════*/

package sugar

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

// RunSugarTradingDemo demonstrates MT5Sugar trading methods:
// - Simple Trading (BuyMarket, SellMarket, BuyLimit, SellLimit, BuyStop, SellStop)
// - Trading with SL/TP (BuyMarketWithSLTP, SellMarketWithSLTP, etc.)
func RunSugarTradingDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("MT5 SUGAR API - TRADING DEMO (Market & Pending Orders)")
	fmt.Println(strings.Repeat("=", 80))

	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// Connect
	fmt.Println("\n📡 Connecting to MT5...")
	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	helpers.Fatal(err, "Failed to create Sugar instance")
	sugar.GetAccount().ApiKey = cfg.ApiKey

	err = sugar.QuickConnect(cfg.MtCluster)
	helpers.Fatal(err, "Connection failed")
	defer sugar.Disconnect()
	fmt.Println("  ✓ Connected!")

	ctx := context.Background()

	testSymbol := cfg.TestSymbol
	minVolume := 0.01 // Minimum lot size

	// Get current price
	bid, _ := sugar.GetBid(testSymbol)
	ask, _ := sugar.GetAsk(testSymbol)
	fmt.Printf("\nℹ️  Current %s prices: BID=%.5f, ASK=%.5f\n", testSymbol, bid, ask)

	// ══════════════════════════════════════════════════════════════════════════════
	// 1. SIMPLE MARKET TRADING
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("1. SIMPLE MARKET TRADING")
	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 1.1. BuyMarket()
	//      Open BUY position at current market ASK price.
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of opened position.
	//      DANGEROUS operation - executes real trade!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.1. BuyMarket() - Open BUY position %.2f lots on %s\n", minVolume, testSymbol)

	buyTicket, err := sugar.BuyMarket(testSymbol, minVolume)
	if !helpers.PrintShortError(err, "BuyMarket failed") {
		fmt.Printf("  ✓ BUY order opened! Ticket: %d\n", buyTicket)

		// Close immediately
		time.Sleep(1 * time.Second)
		fmt.Printf("  Closing position #%d...\n", buyTicket)
		closeErr := sugar.ClosePosition(buyTicket)
		if !helpers.PrintShortError(closeErr, "Failed to close position") {
			fmt.Println("  ✓ Position closed")
		}
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 1.2. SellMarket()
	//      Open SELL position at current market BID price.
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of opened position.
	//      DANGEROUS operation - executes real trade!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n1.2. SellMarket() - Open SELL position %.2f lots on %s\n", minVolume, testSymbol)
	
	sellTicket, err := sugar.SellMarket(testSymbol, minVolume)
	if !helpers.PrintShortError(err, "SellMarket failed") {
		fmt.Printf("  ✓ SELL order opened! Ticket: %d\n", sellTicket)

		// Close immediately
		time.Sleep(1 * time.Second)
		fmt.Printf("  Closing position #%d...\n", sellTicket)
		closeErr := sugar.ClosePosition(sellTicket)
		if !helpers.PrintShortError(closeErr, "Failed to close position") {
			fmt.Println("  ✓ Position closed")
		}
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════════════════════
	// 2. PENDING ORDERS
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("2. PENDING ORDERS (Limit & Stop)")
	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 2.1. BuyLimit()
	//      Place BUY LIMIT pending order (executes when price drops).
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of pending order.
	//      DANGEROUS operation - places real pending order!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	buyLimitPrice := bid - 0.0020 // 20 pips below current bid
	fmt.Printf("\n2.1. BuyLimit() - Place BUY LIMIT at %.5f\n", buyLimitPrice)
	buyLimitTicket, err := sugar.BuyLimit(testSymbol, minVolume, buyLimitPrice)
	if !helpers.PrintShortError(err, "BuyLimit failed") {
		fmt.Printf("  ✓ BUY LIMIT placed! Ticket: %d\n", buyLimitTicket)
		fmt.Printf("  Will execute if price drops to %.5f\n", buyLimitPrice)
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 2.2. SellLimit()
	//      Place SELL LIMIT pending order (executes when price rises).
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of pending order.
	//      DANGEROUS operation - places real pending order!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	sellLimitPrice := ask + 0.0020 // 20 pips above current ask
	fmt.Printf("\n2.2. SellLimit() - Place SELL LIMIT at %.5f\n", sellLimitPrice)
	sellLimitTicket, err := sugar.SellLimit(testSymbol, minVolume, sellLimitPrice)
	if !helpers.PrintShortError(err, "SellLimit failed") {
		fmt.Printf("  ✓ SELL LIMIT placed! Ticket: %d\n", sellLimitTicket)
		fmt.Printf("  Will execute if price rises to %.5f\n", sellLimitPrice)
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 2.3. BuyStop()
	//      Place BUY STOP pending order (executes when price rises).
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of pending order.
	//      DANGEROUS operation - places real pending order!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	buyStopPrice := ask + 0.0020 // 20 pips above current ask
	fmt.Printf("\n2.3. BuyStop() - Place BUY STOP at %.5f\n", buyStopPrice)
	buyStopTicket, err := sugar.BuyStop(testSymbol, minVolume, buyStopPrice)
	if !helpers.PrintShortError(err, "BuyStop failed") {
		fmt.Printf("  ✓ BUY STOP placed! Ticket: %d\n", buyStopTicket)
		fmt.Printf("  Will execute if price rises to %.5f\n", buyStopPrice)
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 2.4. SellStop()
	//      Place SELL STOP pending order (executes when price drops).
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of pending order.
	//      DANGEROUS operation - places real pending order!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	sellStopPrice := bid - 0.0020 // 20 pips below current bid
	fmt.Printf("\n2.4. SellStop() - Place SELL STOP at %.5f\n", sellStopPrice)
	sellStopTicket, err := sugar.SellStop(testSymbol, minVolume, sellStopPrice)
	if !helpers.PrintShortError(err, "SellStop failed") {
		fmt.Printf("  ✓ SELL STOP placed! Ticket: %d\n", sellStopTicket)
		fmt.Printf("  Will execute if price drops to %.5f\n", sellStopPrice)
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════════════════════
	// 3. TRADING WITH SL/TP
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("3. TRADING WITH SL/TP (Stop Loss & Take Profit)")
	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 3.1. BuyMarketWithSLTP()
	//      Open BUY position with Stop Loss and Take Profit.
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of opened position.
	//      DANGEROUS operation - executes real trade with SL/TP!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	stopLoss := bid - 0.0030    // 30 pips SL
	takeProfit := bid + 0.0050  // 50 pips TP
	fmt.Printf("\n3.1. BuyMarketWithSLTP() - BUY with SL/TP\n")
	fmt.Printf("  Entry: %.5f, SL: %.5f (-30 pips), TP: %.5f (+50 pips)\n", ask, stopLoss, takeProfit)
	buySlTpTicket, err := sugar.BuyMarketWithSLTP(testSymbol, minVolume, stopLoss, takeProfit)
	if !helpers.PrintShortError(err, "BuyMarketWithSLTP failed") {
		fmt.Printf("  ✓ BUY with SL/TP opened! Ticket: %d\n", buySlTpTicket)

		// Close immediately
		time.Sleep(1 * time.Second)
		fmt.Printf("  Closing position #%d...\n", buySlTpTicket)
		closeErr := sugar.ClosePosition(buySlTpTicket)
		if !helpers.PrintShortError(closeErr, "Failed to close position") {
			fmt.Println("  ✓ Position closed")
		}
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 3.2. SellMarketWithSLTP()
	//      Open SELL position with Stop Loss and Take Profit.
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of opened position.
	//      DANGEROUS operation - executes real trade with SL/TP!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	stopLossSell := ask + 0.0030    // 30 pips SL
	takeProfitSell := ask - 0.0050  // 50 pips TP
	fmt.Printf("\n3.2. SellMarketWithSLTP() - SELL with SL/TP\n")
	fmt.Printf("  Entry: %.5f, SL: %.5f (+30 pips), TP: %.5f (-50 pips)\n", bid, stopLossSell, takeProfitSell)
	sellSlTpTicket, err := sugar.SellMarketWithSLTP(testSymbol, minVolume, stopLossSell, takeProfitSell)
	if !helpers.PrintShortError(err, "SellMarketWithSLTP failed") {
		fmt.Printf("  ✓ SELL with SL/TP opened! Ticket: %d\n", sellSlTpTicket)

		// Close immediately
		time.Sleep(1 * time.Second)
		fmt.Printf("  Closing position #%d...\n", sellSlTpTicket)
		closeErr := sugar.ClosePosition(sellSlTpTicket)
		if !helpers.PrintShortError(closeErr, "Failed to close position") {
			fmt.Println("  ✓ Position closed")
		}
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 3.3. BuyLimitWithSLTP()
	//      Place BUY LIMIT pending order with SL/TP.
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of pending order.
	//      DANGEROUS operation - places real pending order with SL/TP!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	buyLimitSlTpPrice := bid - 0.0020
	stopLossBuyLimit := buyLimitSlTpPrice - 0.0030
	takeProfitBuyLimit := buyLimitSlTpPrice + 0.0050
	fmt.Printf("\n3.3. BuyLimitWithSLTP() - BUY LIMIT with SL/TP\n")
	fmt.Printf("  Entry: %.5f, SL: %.5f, TP: %.5f\n", buyLimitSlTpPrice, stopLossBuyLimit, takeProfitBuyLimit)
	buyLimitSlTpTicket, err := sugar.BuyLimitWithSLTP(testSymbol, minVolume, buyLimitSlTpPrice, stopLossBuyLimit, takeProfitBuyLimit)
	if !helpers.PrintShortError(err, "BuyLimitWithSLTP failed") {
		fmt.Printf("  ✓ BUY LIMIT with SL/TP placed! Ticket: %d\n", buyLimitSlTpTicket)
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// ══════════════════════════════════════════════════════════════
	// 3.4. SellLimitWithSLTP()
	//      Place SELL LIMIT pending order with SL/TP.
	//      Chain: Sugar → Service.PlaceOrder() → Account.OrderSend() → gRPC
	//      Returns: ticket number (uint64) of pending order.
	//      DANGEROUS operation - places real pending order with SL/TP!
	//      Timeout: 10 seconds
	// ══════════════════════════════════════════════════════════════
	sellLimitSlTpPrice := ask + 0.0020
	stopLossSellLimit := sellLimitSlTpPrice + 0.0030
	takeProfitSellLimit := sellLimitSlTpPrice - 0.0050
	fmt.Printf("\n3.4. SellLimitWithSLTP() - SELL LIMIT with SL/TP\n")
	fmt.Printf("  Entry: %.5f, SL: %.5f, TP: %.5f\n", sellLimitSlTpPrice, stopLossSellLimit, takeProfitSellLimit)
	sellLimitSlTpTicket, err := sugar.SellLimitWithSLTP(testSymbol, minVolume, sellLimitSlTpPrice, stopLossSellLimit, takeProfitSellLimit)
	if !helpers.PrintShortError(err, "SellLimitWithSLTP failed") {
		fmt.Printf("  ✓ SELL LIMIT with SL/TP placed! Ticket: %d\n", sellLimitSlTpTicket)
	}

	// Wait for server to process
	time.Sleep(3 * time.Second)

	// =============================
	// 4. CLEANUP PENDING ORDERS
	// =============================
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("4. CLEANUP - Canceling all pending orders")
	fmt.Println(strings.Repeat("-", 80))

	// Get all opened orders (positions + pending)
	openedData, err := sugar.GetService().GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if !helpers.PrintShortError(err, "Failed to get opened orders") {
		// Cancel pending orders
		cancelCount := 0
		for _, pending := range openedData.OpenedOrders {
			ticket := pending.Ticket
			fmt.Printf("  Canceling pending order #%d...\n", ticket)

			// Use Service's CloseOrder to delete pending order
			closeReq := &pb.OrderCloseRequest{
				Ticket: ticket,
			}
			retCode, closeErr := sugar.GetService().CloseOrder(ctx, closeReq)
			if !helpers.PrintShortError(closeErr, "Failed to cancel order") {
				if retCode == 10009 {
					fmt.Printf("    ✓ Canceled\n")
					cancelCount++
				} else {
					fmt.Printf("    ⚠️  Return code: %d\n", retCode)
				}
			}
		}

		if cancelCount > 0 {
			fmt.Printf("\n  ✓ Canceled %d pending order(s)\n", cancelCount)
		} else {
			fmt.Println("\n  ℹ️  No pending orders to cancel")
		}
	}

	// =============================
	// SUMMARY
	// =============================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("DEMO COMPLETED SUCCESSFULLY!")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nℹ️  MT5Sugar Trading Features:")
	fmt.Println("  • Simple market orders: BuyMarket/SellMarket (3 params)")
	fmt.Println("  • Pending orders: BuyLimit/SellLimit/BuyStop/SellStop")
	fmt.Println("  • SL/TP variants: BuyMarketWithSLTP, SellMarketWithSLTP")
	fmt.Println("  • Direct ticket returns: ticket, err := sugar.BuyMarket(...)")
	fmt.Println("\n✅ All trading operations demonstrated!")
	fmt.Println("   Next: Run demo 08 - sugar_positions.go for position management")
}
