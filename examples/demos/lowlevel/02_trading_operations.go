/*══════════════════════════════════════════════════════════════════════════════
 FILE: 02_trading_operations.go - LOW-LEVEL TRADING OPERATIONS DEMO

 ⚠️  DANGER: This demo executes REAL TRADING operations!
             Use ONLY on DEMO accounts! Real money at risk otherwise!

 PURPOSE:
   Demonstrates MT5 trading operations: calculations, validation, and execution.
   Shows the complete workflow from order preparation to closing positions.
   CRITICAL: All operations execute REAL trades - verify demo account first!


 📚 WHAT THIS DEMO COVERS (4 Steps):

   STEP 1: CREATE MT5ACCOUNT & CONNECT
      • NewMT5Account() - Initialize account
      • Connect() - Connect to MT5 terminal

   STEP 2: ORDER CALCULATIONS (SAFE - READ-ONLY)
      • OrderCalcMargin() - Calculate required margin for order
      • OrderCalcProfit() - Calculate potential profit/loss

   STEP 3: ORDER VALIDATION (SAFE - READ-ONLY)
      • OrderCheck() - Validate order parameters before sending

   STEP 4: TRADING OPERATIONS (⚠️ DANGEROUS - REAL TRADES!)
      • OrderSend() - Place market BUY order
      • OrderModify() - Add/modify Stop Loss and Take Profit
      • OrderClose() - Close opened position

   FINAL: DISCONNECT
      • Disconnect() - Close connection to MT5

 ⚠️  SAFETY WARNINGS:
   • ALWAYS use DEMO accounts for testing
   • Check account type before running (use capabilities check)
   • Understand each operation before executing
   • Real money accounts require additional safeguards
   • This demo is for EDUCATIONAL purposes only

 💡 RECOMMENDED WORKFLOW:
   1. Verify you're on DEMO account: go run main.go capabilities
   2. Understand calculations and validation (STEP 2-3)
   3. Only then proceed to trading operations (STEP 4)

 🚀 HOW TO RUN THIS DEMO:
   cd examples/demos
   go run main.go 2          (or select [2] from interactive menu)

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
)


func RunTrading02() error {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("MT5 LOWLEVEL DEMO 02: TRADING OPERATIONS (DANGEROUS)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("⚠️  WARNING: This demo uses REAL TRADING operations!")
	fmt.Println("    Make sure you are using a DEMO account!")
	fmt.Println()

	// ═══════════════════════════════════════════════════════════════════════
	// LOAD CONFIGURATION
	// ═══════════════════════════════════════════════════════════════════════
	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 1: CREATE MT5ACCOUNT INSTANCE & CONNECT
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("STEP 1: Creating MT5Account instance and connecting...")
	fmt.Println("───────────────────────────────────────────────────────────")

	account, err := mt5.NewMT5Account(cfg.User, cfg.Password, cfg.GrpcServer, uuid.New())
	helpers.Fatal(err, "Failed to create MT5Account")
	if cfg.ApiKey != "" {
		account.ApiKey = cfg.ApiKey
	}
	defer account.DisconnectAndClose()
	fmt.Printf("✓ MT5Account created (UUID: %s)\n", account.Id)

	// Create cancellable context for proper cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ConnectEx - Connect to MT5 cluster
	connectExReq := &pb.ConnectExRequest{
		User:          cfg.User,
		Password:      cfg.Password,
		MtClusterName: cfg.MtCluster,
	}

	// Use context timeout for ConnectEx (replaces old TerminalReadinessWaitingTimeoutSeconds)
	connectCtx, connectCancel := context.WithTimeout(ctx, 180*time.Second)
	defer connectCancel()

	connectData, err := account.ConnectEx(connectCtx, connectExReq)
	helpers.Fatal(err, "ConnectEx failed")

	// CRITICAL: Update account GUID with the one returned by server
	account.TerminalInstanceGuid = connectData.TerminalInstanceGuid
	account.Id = mt5.ParseGuidSafe(connectData.TerminalInstanceGuid)
	fmt.Printf("✓ Connected (Terminal GUID: %s)\n", connectData.TerminalInstanceGuid)

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 2: ORDER CALCULATIONS (SAFE - READ-ONLY)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 2: Order Calculations (Safe - Read-Only)")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 2.1. ORDER CALC MARGIN
	//      Calculate margin required to open an order.
	//      Returns: Required margin in account currency.
	//      SAFE operation - no trades executed.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.1. OrderCalcMargin() - Calculate required margin")

	// Get current price first
	tickReq := &pb.SymbolInfoTickRequest{
		Symbol: cfg.TestSymbol,
	}
	tickData, err := account.SymbolInfoTick(ctx, tickReq)
	helpers.Fatal(err, "SymbolInfoTick failed")

	// Calculate margin for a BUY order
	marginReq := &pb.OrderCalcMarginRequest{
		OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Symbol:    cfg.TestSymbol,
		Volume:    cfg.TestVolume,
		OpenPrice: tickData.Ask, // Use current Ask price for BUY
	}
	marginData, err := account.OrderCalcMargin(ctx, marginReq)
	if !helpers.PrintShortError(err, "OrderCalcMargin failed") {
		// Direct field access: Margin
		fmt.Printf("  Symbol:                        %s\n", cfg.TestSymbol)
		fmt.Printf("  Action:                        BUY\n")
		fmt.Printf("  Volume:                        %.2f lots\n", cfg.TestVolume)
		fmt.Printf("  Price:                         %.5f\n", tickData.Ask)
		fmt.Printf("  Required Margin:               %.2f\n", marginData.Margin)
	}

	// ══════════════════════════════════════════════════════════════
	// 2.2. ORDER CALC PROFIT
	//      Calculate potential profit/loss for an order.
	//      Returns: P&L in account currency for given entry/exit.
	//      SAFE operation - no trades executed.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.2. OrderCalcProfit() - Calculate potential profit/loss")

	// Calculate profit if we BUY at Ask and SELL at Bid (immediate loss due to spread)
	profitReq := &pb.OrderCalcProfitRequest{
		OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Symbol:     cfg.TestSymbol,
		Volume:     cfg.TestVolume,
		OpenPrice:  tickData.Ask, // Entry price
		ClosePrice: tickData.Bid, // Exit price (immediate close = spread loss)
	}
	profitData, err := account.OrderCalcProfit(ctx, profitReq)
	if !helpers.PrintShortError(err, "OrderCalcProfit failed") {
		// Direct field access: Profit
		fmt.Printf("  Symbol:                        %s\n", cfg.TestSymbol)
		fmt.Printf("  Action:                        BUY\n")
		fmt.Printf("  Volume:                        %.2f lots\n", cfg.TestVolume)
		fmt.Printf("  Price Open (Ask):              %.5f\n", tickData.Ask)
		fmt.Printf("  Price Close (Bid):             %.5f\n", tickData.Bid)
		fmt.Printf("  Potential Profit/Loss:         %.2f (spread loss)\n", profitData.Profit)
	}

	// Calculate profit with 10 pips profit target
	pipSize := 0.0001 // For EURUSD, 1 pip = 0.0001
	targetPrice := tickData.Ask + (10 * pipSize)
	profitTargetReq := &pb.OrderCalcProfitRequest{
		OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Symbol:     cfg.TestSymbol,
		Volume:     cfg.TestVolume,
		OpenPrice:  tickData.Ask,
		ClosePrice: targetPrice,
	}
	profitTargetData, err := account.OrderCalcProfit(ctx, profitTargetReq)
	if !helpers.PrintShortError(err, "OrderCalcProfit (target) failed") {
		fmt.Printf("\n  If price moves +10 pips to %.5f:\n", targetPrice)
		fmt.Printf("  Potential Profit:              %.2f\n", profitTargetData.Profit)
	}

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 3: ORDER VALIDATION (SAFE - READ-ONLY)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 3: Order Validation (Safe - Read-Only)")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 3.1. ORDER CHECK
	//      Validate order parameters BEFORE sending to broker.
	//      Returns: Simulated balance/equity/margin after deal.
	//      SAFE operation - no trades executed (dry-run).
	//      ⚠️  NOTE: Often fails on DEMO accounts (broker limitation)
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n3.1. OrderCheck() - Validate order parameters")

	// Create order check request with IOC filling (most compatible)
	orderCheckReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:      pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
			Symbol:      cfg.TestSymbol,
			Volume:      cfg.TestVolume,
			OrderType:   pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
			Price:       tickData.Ask,
			StopLoss:    0.0,
			TakeProfit:  0.0,
			Deviation:   10,
			TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_IOC, // IOC - most compatible
			TypeTime:    pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC,
			Comment:     "OrderCheck validation",
		},
	}

	checkData, err := account.OrderCheck(ctx, orderCheckReq)
	if err != nil {
		// OrderCheck failed - this is EXPECTED on demo accounts
		fmt.Println("  ❌ OrderCheck FAILED (expected on demo accounts)")
		fmt.Printf("     Error: %v\n", err)

		fmt.Println("\n  ℹ️  This is a known limitation:")
		fmt.Println("     • DEMO accounts: OrderCheck often not supported by broker")
		fmt.Println("     • gRPC gateway: May have null-check issues with Timestamp")
		fmt.Println("     • Workaround: Use OrderCalcMargin() for validation (step 2.1)")
		fmt.Println("     • OrderSend() will work despite OrderCheck failure")
		fmt.Println()
	} else if checkData != nil && checkData.MqlTradeCheckResult != nil {
		// OrderCheck succeeded!
		result := checkData.MqlTradeCheckResult
		fmt.Println("  ✅ OrderCheck SUCCESS!")
		fmt.Printf("     Return Code:        %d\n", result.ReturnedCode)
		fmt.Printf("     Comment:            %s\n", result.Comment)
		fmt.Printf("     Required Margin:    %.2f\n", result.Margin)
		fmt.Printf("     Balance After:      %.2f\n", result.BalanceAfterDeal)
		fmt.Printf("     Equity After:       %.2f\n", result.EquityAfterDeal)
		fmt.Printf("     Free Margin After:  %.2f\n", result.FreeMargin)
		fmt.Printf("     Margin Level:       %.2f%%\n\n", result.MarginLevel)

		if result.ReturnedCode == 0 {
			fmt.Println("  ✓ Order validation PASSED - safe to proceed with OrderSend()")
		} else {
			fmt.Printf("  ⚠️  Validation returned code %d - check before trading\n", result.ReturnedCode)
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 4: TRADING OPERATIONS (DANGEROUS - REAL TRADES!)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\n⚠️  STEP 4: Trading Operations (DANGEROUS - Real Trades!)")
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println("  The following operations will place REAL orders!")
	fmt.Println("  Make sure you are on a DEMO account!")
	fmt.Println()
	// ══════════════════════════════════════════════════════════════
	// 4.1. ORDER SEND (⚠️ REAL TRADE!)
	//      Place a market BUY order on broker.
	//      Returns: Order ticket, deal ticket, execution price.
	//      DANGEROUS - Executes REAL trade with REAL money!
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n4.1. OrderSend() - Place market BUY order")

	slippage := uint64(10)
	orderSendReq := &pb.OrderSendRequest{
		Symbol:    cfg.TestSymbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
		Volume:    cfg.TestVolume,
		Price:     &tickData.Ask,
		Slippage:  &slippage,
		// StopLoss and TakeProfit can be set here or later via OrderModify
	}

	sendData, err := account.OrderSend(ctx, orderSendReq)
	if !helpers.PrintShortError(err, "OrderSend failed") {
		// Direct field access to OrderSendData
		fmt.Printf("  Order sent result:\n")
		fmt.Printf("    Return Code:                 %d\n", sendData.ReturnedCode)
		fmt.Printf("    Deal Ticket:                 %d\n", sendData.Deal)
		fmt.Printf("    Order Ticket:                %d\n", sendData.Order)
		fmt.Printf("    Volume:                      %.2f\n", sendData.Volume)
		fmt.Printf("    Execution Price:             %.5f\n", sendData.Price)
		// Bid/Ask are optional fields, often 0 in response - skip if not provided
		if sendData.Bid > 0 && sendData.Ask > 0 {
			fmt.Printf("    Market Bid:                  %.5f\n", sendData.Bid)
			fmt.Printf("    Market Ask:                  %.5f\n", sendData.Ask)
		}
		fmt.Printf("    Comment:                     %s\n", sendData.Comment)

		// Check if order executed successfully using helper
		// This demonstrates proper ReturnedCode validation
		if helpers.CheckRetCode(sendData.ReturnedCode, "Order execution") {
			orderTicket := sendData.Order

			// ══════════════════════════════════════════════════════════════
			// 4.2. ORDER MODIFY (⚠️ MODIFIES REAL POSITION!)
			//      Modify opened position - add/change SL/TP levels.
			//      Returns: Modification result.
			//      DANGEROUS - Modifies REAL open position!
			// ══════════════════════════════════════════════════════════════
			fmt.Println("\n4.2. OrderModify() - Add Stop Loss and Take Profit")

			// Calculate SL/TP levels (10 pips SL, 20 pips TP)
			stopLoss := sendData.Price - (10 * pipSize)
			takeProfit := sendData.Price + (20 * pipSize)

			modifyReq := &pb.OrderModifyRequest{
				Ticket:     orderTicket,
				StopLoss:   &stopLoss,
				TakeProfit: &takeProfit,
			}

			modifyData, err := account.OrderModify(ctx, modifyReq)
			if !helpers.PrintShortError(err, "OrderModify failed") {
				fmt.Printf("  Order modify result:\n")
				fmt.Printf("    Return Code:                 %d\n", modifyData.ReturnedCode)
				fmt.Printf("    Order Ticket:                %d\n", modifyData.Order)
				fmt.Printf("    Stop Loss:                   %.5f\n", stopLoss)
				fmt.Printf("    Take Profit:                 %.5f\n", takeProfit)
				fmt.Printf("    Comment:                     %s\n", modifyData.Comment)

				// Check modification result using helper
				helpers.CheckRetCode(modifyData.ReturnedCode, "Position modification")
			}

			// ══════════════════════════════════════════════════════════════
			// 4.3. ORDER CLOSE (⚠️ CLOSES REAL POSITION!)
			//      Close opened position at market price.
			//      Returns: Close result with final P&L.
			//      DANGEROUS - Closes REAL position, realizes profit/loss!
			// ══════════════════════════════════════════════════════════════
			fmt.Println("\n4.3. OrderClose() - Close the position")

			closeReq := &pb.OrderCloseRequest{
				Ticket:   orderTicket,
				Volume:   cfg.TestVolume,
				Slippage: 10,
			}

			closeData, err := account.OrderClose(ctx, closeReq)
			if !helpers.PrintShortError(err, "OrderClose failed") {
				fmt.Printf("  Order close result:\n")
				fmt.Printf("    Return Code:                 %d (%s)\n", closeData.ReturnedCode, closeData.ReturnedStringCode)
				fmt.Printf("    Description:                 %s\n", closeData.ReturnedCodeDescription)
				fmt.Printf("    Close Mode:                  %v\n", closeData.CloseMode)

				// Check close result using helper
				helpers.CheckRetCode(closeData.ReturnedCode, "Position close")
			}
		} else {
			fmt.Printf("    ✗ Order execution FAILED - check return code and comment\n")
		}
	}

	// ═══════════════════════════════════════════════════════════════════════
	// FINAL: DISCONNECT
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nFinal: Disconnecting...")
	fmt.Println("───────────────────────────────────────────────────────────")

	// Disconnect and close connection
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer disconnectCancel()
	account.Disconnect(disconnectCtx, &pb.DisconnectRequest{})
	account.Close()
	fmt.Println("✓ Connection closed")

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("✓ DEMO COMPLETED SUCCESSFULLY")
	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}
