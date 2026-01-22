/*══════════════════════════════════════════════════════════════════════════════
 FILE: 05_service_streaming.go - SERVICE STREAMING METHODS DEMO

 PURPOSE:
   Demonstrates MT5Service streaming (real-time) methods for live data updates.
   Shows all 5 streaming methods: Ticks, Trade Updates, Position Profits,
   Opened Tickets, and Transactions.

 🎯 WHO SHOULD USE THIS:
   • Developers building real-time monitoring systems
   • Trading bots that react to market events
   • Applications requiring instant P&L updates
   • Dashboard applications

 📚 WHAT THIS DEMO COVERS (5 Streaming Methods):

   STREAM 1: StreamTicks()
      • Real-time Bid/Ask price updates
      • High-frequency tick data
      • Subscribe to multiple symbols

   STREAM 2: StreamTradeUpdates()
      • New/closed orders and positions
      • History orders and deals
      • Trade execution events

   STREAM 3: StreamPositionProfits()
      • Real-time P&L updates
      • Position profit changes as prices move
      • Position lifecycle (new/updated/deleted)

   STREAM 4: StreamOpenedTickets()
      • Track order/position ticket numbers
      • Lightweight: just ticket IDs
      • Know what's currently open

   STREAM 5: StreamTransactions()
      • Most comprehensive stream
      • All trade transaction details
      • Order/deal/position states

 ⚠️  IMPORTANT - Stream Timing:
   • Start streams BEFORE triggering events you want to catch
   • Streams capture real-time events only (not historical)
   • Add small delay after starting stream (500ms recommended)
   • This demo opens positions AFTER starting streams to demonstrate this

 ⚠️  IMPORTANT - General:
   • Streams are high-frequency (especially ticks)
   • Always use context cancellation or timeout
   • Handle both data and error channels
   • Close streams gracefully to avoid resource leaks

 🚀 HOW TO RUN THIS DEMO:
   cd examples/demos
   go run main.go 5           (or select [5] from menu)
   

══════════════════════════════════════════════════════════════════════════════*/

package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	helpers_mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
	"github.com/google/uuid"
)

const (
	MAX_STREAM_EVENTS  = 5  // Max events per stream
	MAX_STREAM_SECONDS = 10 // Max duration per stream
)

// RunServiceStreaming05 demonstrates all MT5Service streaming methods
func RunServiceStreaming05() error {
	fmt.Println("\n" + repeatString("=", 80))
	fmt.Println("MT5 SERVICE - STREAMING METHODS DEMO")
	fmt.Println(repeatString("=", 80))

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create account and service
	ctx := context.Background()

	account, err := helpers_mt5.NewMT5Account(cfg.User, cfg.Password, cfg.GrpcServer, uuid.New())
	if err != nil {
		return fmt.Errorf("failed to create MT5Account: %w", err)
	}
	defer account.Close()

	service := mt5.NewMT5Service(account)

	// Connect
	fmt.Println("\n📡 Connecting to MT5...")
	// Use context timeout for ConnectEx (replaces old TerminalReadinessWaitingTimeoutSeconds)
	connCtx, connCancel := context.WithTimeout(ctx, 180*time.Second)
	defer connCancel()

	baseSymbol := cfg.TestSymbol
	data, err := account.ConnectEx(connCtx, &pb.ConnectExRequest{
		User:            cfg.User,
		Password:        cfg.Password,
		MtClusterName:   cfg.MtCluster,
		BaseChartSymbol: &baseSymbol,
	})
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	account.Id = uuid.MustParse(data.TerminalInstanceGuid)
	fmt.Printf("  ✓ Connected successfully (Terminal GUID: %s)\n", data.TerminalInstanceGuid)
	fmt.Printf("  Account: %d | Server: %s\n", cfg.User, cfg.MtCluster)

	defer func() {
		fmt.Println("\n📴 Disconnecting...")
		account.Disconnect(context.Background(), &pb.DisconnectRequest{})
		account.Close()
	}()

	fmt.Println("\n" + repeatString("-", 80))
	fmt.Println("STREAM 1: StreamTicks() - Real-time Bid/Ask price updates")
	fmt.Println(repeatString("-", 80))
	// ══════════════════════════════════════════════════════════════
	// STREAM 1: StreamTicks()
	//      Streams real-time Bid/Ask price updates for symbols.
	//      Returns: (<-chan *OnSymbolTickData, <-chan error)
	//      Use for: Price monitoring, tick charts, HFT strategies
	// ══════════════════════════════════════════════════════════════

	streamCtx1, cancel1 := context.WithTimeout(ctx, MAX_STREAM_SECONDS*time.Second)
	defer cancel1()

	tickCh, tickErrCh := service.StreamTicks(streamCtx1, []string{cfg.TestSymbol})

	fmt.Printf("\nStreaming %s tick data (max %d events or %d seconds)...\n",
		cfg.TestSymbol, MAX_STREAM_EVENTS, MAX_STREAM_SECONDS)

	eventCount := 0
	for eventCount < MAX_STREAM_EVENTS {
		select {
		case tick, ok := <-tickCh:
			if !ok {
				fmt.Println("  Stream closed")
				goto stream2
			}
			eventCount++
			fmt.Printf("  Tick #%d: BID=%.5f ASK=%.5f Spread=%.5f Time=%s\n",
				eventCount, tick.Bid, tick.Ask, tick.Ask-tick.Bid,
				tick.Time.Format("15:04:05"))

		case err := <-tickErrCh:
			fmt.Printf("  ❌ Stream error: %v\n", err)
			goto stream2

		case <-streamCtx1.Done():
			fmt.Printf("  ⏱️  Timeout after %d seconds\n", MAX_STREAM_SECONDS)
			goto stream2
		}
	}
	cancel1()

stream2:
	fmt.Println("\n" + repeatString("-", 80))
	fmt.Println("STREAM 2: StreamTradeUpdates() - Trade execution events")
	fmt.Println(repeatString("-", 80))
	// ══════════════════════════════════════════════════════════════
	// STREAM 2: StreamTradeUpdates()
	//      Streams trade events (new/closed orders and positions).
	//      Returns: (<-chan *OnTradeData, <-chan error)
	//      Use for: Monitoring trade execution, order lifecycle
	// ══════════════════════════════════════════════════════════════

	streamCtx2, cancel2 := context.WithTimeout(ctx, MAX_STREAM_SECONDS*time.Second)
	defer cancel2()

	// Start streaming BEFORE opening position
	tradeCh, tradeErrCh := service.StreamTradeUpdates(streamCtx2)

	fmt.Printf("\nStreaming trade events (max %d events or %d seconds)...\n",
		MAX_STREAM_EVENTS, MAX_STREAM_SECONDS)
	fmt.Println("  ℹ️  Stream started, waiting for subscription...")

	// Give stream time to establish connection
	time.Sleep(500 * time.Millisecond)

	// Now open a test position to generate trade events
	fmt.Println("  📤 Opening test position to trigger events...")
	placeReq := &pb.OrderSendRequest{
		Symbol:    cfg.TestSymbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
		Volume:    0.01,
	}

	var testTicket uint64
	placeResult, err := service.PlaceOrder(ctx, placeReq)
	if helpers.PrintShortError(err, "Failed to place test order") {
		// Error already printed
	} else if placeResult.ReturnedCode == 10009 {
		testTicket = placeResult.Order
		fmt.Printf("  ✓ Test position opened: #%d\n", testTicket)
	}

	eventCount = 0
	for eventCount < MAX_STREAM_EVENTS {
		select {
		case data, ok := <-tradeCh:
			if !ok {
				fmt.Println("  Stream closed")
				goto stream3
			}
			eventCount++
			fmt.Printf("  Event #%d: Type=%v\n", eventCount, data.Type)
			if data.EventData != nil {
				fmt.Printf("    Orders: New=%d, Disappeared=%d\n",
					len(data.EventData.NewOrders), len(data.EventData.DisappearedOrders))
				fmt.Printf("    Positions: New=%d, Disappeared=%d\n",
					len(data.EventData.NewPositions), len(data.EventData.DisappearedPositions))
			}

		case err := <-tradeErrCh:
			fmt.Printf("  ❌ Stream error: %v\n", err)
			goto stream3

		case <-streamCtx2.Done():
			fmt.Printf("  ⏱️  Timeout after %d seconds\n", MAX_STREAM_SECONDS)
			goto stream3
		}
	}
	cancel2()

stream3:
	fmt.Println("\n" + repeatString("-", 80))
	fmt.Println("STREAM 3: StreamPositionProfits() - Real-time P&L updates")
	fmt.Println(repeatString("-", 80))
	// ══════════════════════════════════════════════════════════════
	// STREAM 3: StreamPositionProfits()
	//      Streams real-time profit/loss updates for positions.
	//      Returns: (<-chan *OnPositionProfitData, <-chan error)
	//      Use for: P&L monitoring, risk management, live dashboards
	// ══════════════════════════════════════════════════════════════

	streamCtx3, cancel3 := context.WithTimeout(ctx, MAX_STREAM_SECONDS*time.Second)
	defer cancel3()

	profitCh, profitErrCh := service.StreamPositionProfits(streamCtx3)

	fmt.Printf("\nStreaming position profit updates (max %d events or %d seconds)...\n",
		MAX_STREAM_EVENTS, MAX_STREAM_SECONDS)
	fmt.Println("  ℹ️  Profit updates appear when prices change or positions open/close")

	eventCount = 0
	for eventCount < MAX_STREAM_EVENTS {
		select {
		case data, ok := <-profitCh:
			if !ok {
				fmt.Println("  Stream closed")
				goto stream4
			}
			eventCount++
			fmt.Printf("  Event #%d: Update type=%v\n", eventCount, data.Type)

			totalPositions := len(data.NewPositions) + len(data.UpdatedPositions) + len(data.DeletedPositions)
			if totalPositions > 0 {
				if len(data.NewPositions) > 0 {
					fmt.Printf("    New positions: %d\n", len(data.NewPositions))
				}
				if len(data.UpdatedPositions) > 0 {
					fmt.Printf("    Updated positions: %d\n", len(data.UpdatedPositions))
					for i, pos := range data.UpdatedPositions {
						if i >= 3 {
							fmt.Printf("      ... and %d more\n", len(data.UpdatedPositions)-3)
							break
						}
						profitSign := "+"
						if pos.Profit < 0 {
							profitSign = ""
						}
						fmt.Printf("      #%d (%s): Profit=%s%.2f\n",
							pos.Ticket, pos.PositionSymbol, profitSign, pos.Profit)
					}
				}
				if len(data.DeletedPositions) > 0 {
					fmt.Printf("    Deleted positions: %d\n", len(data.DeletedPositions))
				}
			} else {
				fmt.Println("    No position changes")
			}

		case err := <-profitErrCh:
			fmt.Printf("  ❌ Stream error: %v\n", err)
			goto stream4

		case <-streamCtx3.Done():
			fmt.Printf("  ⏱️  Timeout after %d seconds\n", MAX_STREAM_SECONDS)
			goto stream4
		}
	}
	cancel3()

stream4:
	fmt.Println("\n" + repeatString("-", 80))
	fmt.Println("STREAM 4: StreamOpenedTickets() - Track order/position tickets")
	fmt.Println(repeatString("-", 80))
	// ══════════════════════════════════════════════════════════════
	// STREAM 4: StreamOpenedTickets()
	//      Lightweight stream of open position/order ticket IDs.
	//      Returns: (<-chan *OnPositionsAndPendingOrdersTicketsData, <-chan error)
	//      Use for: Quick overview, simple position tracking
	// ══════════════════════════════════════════════════════════════

	streamCtx4, cancel4 := context.WithTimeout(ctx, MAX_STREAM_SECONDS*time.Second)
	defer cancel4()

	ticketsCh, ticketsErrCh := service.StreamOpenedTickets(streamCtx4)

	fmt.Printf("\nStreaming opened tickets (max %d events or %d seconds)...\n",
		MAX_STREAM_EVENTS, MAX_STREAM_SECONDS)
	fmt.Println("  ℹ️  Lightweight stream - only ticket IDs, not full data")

	eventCount = 0
	for eventCount < MAX_STREAM_EVENTS {
		select {
		case data, ok := <-ticketsCh:
			if !ok {
				fmt.Println("  Stream closed")
				goto stream5
			}
			eventCount++
			fmt.Printf("  Event #%d:\n", eventCount)
			fmt.Printf("    Position tickets: %v\n", data.PositionTickets)
			fmt.Printf("    Pending tickets: %v\n", data.PendingOrderTickets)

		case err := <-ticketsErrCh:
			fmt.Printf("  ❌ Stream error: %v\n", err)
			goto stream5

		case <-streamCtx4.Done():
			fmt.Printf("  ⏱️  Timeout after %d seconds\n", MAX_STREAM_SECONDS)
			goto stream5
		}
	}
	cancel4()

stream5:
	fmt.Println("\n" + repeatString("-", 80))
	fmt.Println("STREAM 5: StreamTransactions() - Detailed trade transactions")
	fmt.Println(repeatString("-", 80))
	// ══════════════════════════════════════════════════════════════
	// STREAM 5: StreamTransactions()
	//      Most detailed stream with all trade transaction data.
	//      Returns: (<-chan *OnTradeTransactionData, <-chan error)
	//      Use for: Audit trails, comprehensive analysis, debugging
	// ══════════════════════════════════════════════════════════════

	streamCtx5, cancel5 := context.WithTimeout(ctx, MAX_STREAM_SECONDS*time.Second)
	defer cancel5()

	// Start streaming BEFORE making trades
	txCh, txErrCh := service.StreamTransactions(streamCtx5)

	fmt.Printf("\nStreaming transactions (max %d events or %d seconds)...\n",
		MAX_STREAM_EVENTS, MAX_STREAM_SECONDS)
	fmt.Println("  ℹ️  Stream started, waiting for subscription...")

	// Give stream time to establish connection
	time.Sleep(500 * time.Millisecond)

	// Modify an existing position to generate transaction events
	openedData2, err := service.GetOpenedOrders(ctx,
		pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if err == nil && len(openedData2.PositionInfos) > 0 {
		pos := openedData2.PositionInfos[0]
		fmt.Printf("  📤 Modifying position #%d to trigger transaction events...\n", pos.Ticket)

		// Get current price to calculate SL/TP
		tick2, _ := service.GetSymbolTick(ctx, cfg.TestSymbol)
		newSL := tick2.Bid - 0.001
		newTP := tick2.Bid + 0.002

		modifyReq := &pb.OrderModifyRequest{
			Ticket:     pos.Ticket,
			StopLoss:   &newSL,
			TakeProfit: &newTP,
		}

		_, modifyErr := service.ModifyOrder(ctx, modifyReq)
		if modifyErr == nil {
			fmt.Printf("  ✓ Position modified (SL/TP updated)\n")
		} else {
			fmt.Printf("  ⚠️  Modification may have failed: %v\n", modifyErr)
		}
	} else {
		fmt.Println("  ℹ️  No open positions to modify, events may be limited")
	}

	eventCount = 0
	for eventCount < MAX_STREAM_EVENTS {
		select {
		case data, ok := <-txCh:
			if !ok {
				fmt.Println("  Stream closed")
				goto cleanup
			}
			eventCount++
			fmt.Printf("  Event #%d: Type=%v\n", eventCount, data.Type)
			if data.TradeTransaction != nil {
				tx := data.TradeTransaction
				fmt.Printf("    Transaction: Type=%v, Symbol=%s\n", tx.Type, tx.Symbol)
				if tx.OrderTicket > 0 {
					fmt.Printf("      Order #%d: Type=%v, State=%v\n",
						tx.OrderTicket, tx.OrderType, tx.OrderState)
				}
				if tx.DealTicket > 0 {
					fmt.Printf("      Deal #%d: Volume=%.2f\n",
						tx.DealTicket, tx.Volume)
				}
			}

		case err := <-txErrCh:
			fmt.Printf("  ❌ Stream error: %v\n", err)
			goto cleanup

		case <-streamCtx5.Done():
			fmt.Printf("  ⏱️  Timeout after %d seconds\n", MAX_STREAM_SECONDS)
			goto cleanup
		}
	}
	cancel5()

cleanup:
	fmt.Println("\n" + repeatString("-", 80))
	fmt.Println("CLEANUP: Closing any test positions")
	fmt.Println(repeatString("-", 80))
	// ══════════════════════════════════════════════════════════════
	// CLEANUP
	//      Close any test positions opened during demo.
	//      Ensures account is in clean state after demo.
	// ══════════════════════════════════════════════════════════════

	openedData, err := service.GetOpenedOrders(ctx,
		pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if err == nil && len(openedData.PositionInfos) > 0 {
		for _, pos := range openedData.PositionInfos {
			fmt.Printf("\n  Closing position #%d...\n", pos.Ticket)
			closeReq := &pb.OrderCloseRequest{
				Ticket: pos.Ticket,
			}
			retCode, closeErr := service.CloseOrder(ctx, closeReq)
			if helpers.PrintShortError(closeErr, "Failed to close") {
				// Error already printed
			} else if retCode == 10009 {
				fmt.Println("    ✓ Closed successfully")
			} else {
				fmt.Printf("    ⚠️  Return code: %d\n", retCode)
			}
		}
	} else {
		fmt.Println("\n  No open positions to close")
	}

	fmt.Println("\n" + repeatString("=", 80))
	fmt.Println("DEMO COMPLETED!")
	fmt.Println(repeatString("=", 80))
	// ══════════════════════════════════════════════════════════════
	// SUMMARY
	//      Overview of all 5 streaming methods and best practices.
	// ══════════════════════════════════════════════════════════════

	fmt.Println("\n📊 Streaming Methods Summary:")
	fmt.Println("\n  1️⃣  StreamTicks() - High-frequency price updates")
	fmt.Println("     Use for: Tick charts, price monitoring, HFT strategies")
	fmt.Println("     Frequency: Multiple events per second")

	fmt.Println("\n  2️⃣  StreamTradeUpdates() - Trade event notifications")
	fmt.Println("     Use for: Monitoring order/position changes")
	fmt.Println("     Frequency: When trades execute")

	fmt.Println("\n  3️⃣  StreamPositionProfits() - Real-time P&L tracking")
	fmt.Println("     Use for: Live profit monitoring, risk management")
	fmt.Println("     Frequency: When prices change or positions update")

	fmt.Println("\n  4️⃣  StreamOpenedTickets() - Lightweight ticket tracking")
	fmt.Println("     Use for: Quick overview of open positions/orders")
	fmt.Println("     Frequency: When tickets are added/removed")

	fmt.Println("\n  5️⃣  StreamTransactions() - Comprehensive transaction log")
	fmt.Println("     Use for: Detailed trade analysis, audit trails")
	fmt.Println("     Frequency: All trade-related events")

	fmt.Println("\n💡 Best Practices:")
	fmt.Println("  • Always use context cancellation or timeout")
	fmt.Println("  • Read from both data and error channels (use select)")
	fmt.Println("  • StreamTicks is high-frequency - use wisely")
	fmt.Println("  • For simple P&L: use StreamPositionProfits")
	fmt.Println("  • For detailed analysis: use StreamTransactions")

	fmt.Println("\n⚠️  Important: Stream Timing")
	fmt.Println("  • Start streams BEFORE triggering events you want to catch")
	fmt.Println("  • Streams capture real-time events only (not historical)")
	fmt.Println("  • Add small delay after starting stream (500ms recommended)")

	fmt.Println("\n✅ All 5 streaming methods demonstrated!")

	return nil
}

// Helper function to repeat strings
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
