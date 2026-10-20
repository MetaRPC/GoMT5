/*══════════════════════════════════════════════════════════════════════════════
 FILE: 03_streaming_methods.go - LOW-LEVEL STREAMING METHODS DEMO

 PURPOSE:
   Demonstrates MT5 real-time streaming (gRPC server-streaming) for live data.
   Shows how to subscribe to market ticks, trade events, position updates, and more.
   All streams are REAL-TIME and event-driven.

 🎯 WHO SHOULD USE THIS:
   • Developers building real-time trading applications
   • High-frequency trading (HFT) systems
   • Live monitoring and alerting systems
   • Anyone needing instant market/trade notifications

 📚 WHAT THIS DEMO COVERS (5 Streaming Methods):

   STEP 1: CREATE MT5ACCOUNT & CONNECT
      • NewMT5Account() - Initialize account
      • ConnectEx() - Connect to MT5 cluster

   STEP 2: STREAM 1 - OnSymbolTick()
      • Real-time Bid/Ask price updates
      • Subscribe to multiple symbols
      • High-frequency tick data

   STEP 3: STREAM 2 - OnTrade()
      • Trade execution events
      • New/disappeared orders and positions
      • History orders and deals

   STEP 4: STREAM 3 - OnPositionProfit()
      • Real-time P&L updates for open positions
      • Position profit changes as prices move
      • Position lifecycle (new/updated/deleted)

   STEP 5: STREAM 4 - OnPositionsAndPendingOrdersTickets()
      • Track order/position ticket numbers
      • Know what's currently open
      • Minimal data (just ticket IDs)

   STEP 6: STREAM 5 - OnTradeTransaction()
      • Low-level trade transaction events
      • Every single trading operation
      • Order placement, execution, modification, deletion

   FINAL: DISCONNECT
      • Close all streams
      • Disconnect from MT5

 ⚡ STREAM CHARACTERISTICS:
   • REAL-TIME: Events arrive as they happen (millisecond latency)
   • EVENT-DRIVEN: No polling needed - server pushes data
   • PERSISTENT: Stay open until you close them or error occurs
   • SAFE: Read-only operations (no trades executed)

 ⚠️  IMPORTANT NOTES:
   • Streams require ACTIVE connection to MT5 terminal
   • Some streams only fire on activity (e.g., OnTrade when trading)
   • OnSymbolTick is the most frequent (multiple times per second)
   • Always close streams gracefully (use timeout or max events)
   • Each stream runs in goroutine - handle errors properly

 💡 STREAMING vs POLLING:
   STREAMING (this demo):  Server pushes → instant updates, efficient
   POLLING (file 01):      Client pulls → delayed, resource-intensive

 🚀 HOW TO RUN THIS DEMO:
   cd examples/demos
   go run main.go 3          (or select [3] from interactive menu)

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

func RunStreaming03() error {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("MT5 LOWLEVEL DEMO 03: STREAMING METHODS")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// ═══════════════════════════════════════════════════════════════════════
	// LOAD CONFIGURATION
	// ═══════════════════════════════════════════════════════════════════════
	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// Stream limits
	const (
		MAX_EVENTS  = 10
		MAX_SECONDS = 5
	)

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
	// Runs however the demo ends: the terminal is disconnected (stopped) and the gRPC connection closed.
	// Errors after the connection below are returned (not helpers.Fatal, which exits at once and skips this).
	defer account.DisconnectAndClose()
	fmt.Printf("✓ MT5Account created (UUID: %s)\n", account.Id)

	// Create cancellable context for proper stream cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensures all streams are properly closed
	// Each stream below also gets its own context, cancelled as soon as its loop ends: that closes the stream on
	// the server at once (otherwise every stream would stay open, unread, until the demo ends).

	// ConnectEx - Connect to MT5 cluster
	connectExReq := &pb.ConnectExRequest{
		User:          cfg.User,
		Password:      cfg.Password,
		MtClusterName: cfg.MtCluster,
	}

	// The server waits up to ConnectExRequest.TimeoutSeconds (default 120 s) for the terminal to log in;
	// the call's own deadline must be longer than that.
	connectCtx, connectCancel := context.WithTimeout(ctx, 180*time.Second)
	defer connectCancel()

	connectData, err := account.ConnectEx(connectCtx, connectExReq)
	if err != nil {
		return fmt.Errorf("ConnectEx failed: %w", err)
	}

	// ConnectEx has stored the terminal GUID in the account (account.TerminalInstanceGuid / account.Id)
	fmt.Printf("✓ Connected (Terminal GUID: %s)\n", connectData.TerminalInstanceGuid)

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 2: STREAM 1 - OnSymbolTick (Tick data stream)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 2: OnSymbolTick() - Stream tick data (Bid/Ask updates)")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 2.1. ON SYMBOL TICK
	//      Subscribe to real-time tick data for specified symbols.
	//      Returns: Symbol, Bid, Ask, Last, Volume, Time (broker) / TimeUtc for each tick (spread = Ask - Bid).
	//      HIGH FREQUENCY: Multiple events per second.
	//      Use cases: Price monitoring, tick charts, HFT strategies.
	// ══════════════════════════════════════════════════════════════

	tickReq := &pb.OnSymbolTickRequest{
		SymbolNames: []string{cfg.TestSymbol}, // Symbols to stream (one or more)
	}

	tickCtx, tickCancel := context.WithCancel(ctx)
	tickChan, tickErrChan := account.OnSymbolTick(tickCtx, tickReq)

	fmt.Printf("Streaming %s tick data (max %d events or %d seconds)...\n", cfg.TestSymbol, MAX_EVENTS, MAX_SECONDS)

	eventCount := 0
	timeout := time.After(MAX_SECONDS * time.Second)

streamTick:
	for {
		select {
		case tickData, ok := <-tickChan:
			if !ok {
				fmt.Println("  Stream closed by server")
				break streamTick
			}
			eventCount++
			// The tick itself is nested: SymbolTick (nil in a message without a tick)
			if tickData.SymbolTick != nil {
				tick := tickData.SymbolTick
				fmt.Printf("  Event #%d: Symbol=%s Bid=%.5f Ask=%.5f Spread=%.5f\n",
					eventCount,
					tick.Symbol,
					tick.Bid,
					tick.Ask,
					tick.Ask-tick.Bid)
			}

			if eventCount >= MAX_EVENTS {
				fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
				break streamTick
			}

		case err := <-tickErrChan:
			if err != nil {
				helpers.PrintShortError(err, "Stream error")
				break streamTick
			}

		case <-timeout:
			fmt.Printf("  ⏱ Timeout after %d seconds, stopping stream\n", MAX_SECONDS)
			break streamTick
		}
	}
	tickCancel() // closes this stream on the server

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 3: STREAM 2 - OnTrade (Trade events stream)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 3: OnTrade() - Stream trade events")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 3.1. ON TRADE
	//      Subscribe to trade-related events.
	//      Returns: New/disappeared orders, positions, history changes.
	//      EVENT-DRIVEN: Fires only when trades occur.
	//      Use cases: Trade monitoring, order tracking, execution alerts.
	// ══════════════════════════════════════════════════════════════

	tradeReq := &pb.OnTradeRequest{}

	tradeCtx, tradeCancel := context.WithCancel(ctx)
	tradeChan, tradeErrChan := account.OnTrade(tradeCtx, tradeReq)

	fmt.Printf("Streaming trade events (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
	fmt.Println("  NOTE: This stream sends events only when trades occur")

	eventCount = 0
	timeout = time.After(MAX_SECONDS * time.Second)

streamTrade:
	for {
		select {
		case tradeData, ok := <-tradeChan:
			if !ok {
				fmt.Println("  Stream closed by server")
				break streamTrade
			}
			eventCount++
			// Direct field access to nested EventData structure
			fmt.Printf("  Event #%d: Type=%v\n", eventCount, tradeData.Type)
			if tradeData.EventData != nil {
				ed := tradeData.EventData
				fmt.Printf("    New Orders: %d, Disappeared Orders: %d\n",
					len(ed.NewOrders), len(ed.DisappearedOrders))
				fmt.Printf("    New Positions: %d, Disappeared Positions: %d\n",
					len(ed.NewPositions), len(ed.DisappearedPositions))
				fmt.Printf("    New History Orders: %d, New History Deals: %d\n",
					len(ed.NewHistoryOrders), len(ed.NewHistoryDeals))
			}

			if eventCount >= MAX_EVENTS {
				fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
				break streamTrade
			}

		case err := <-tradeErrChan:
			if err != nil {
				helpers.PrintShortError(err, "Stream error")
				break streamTrade
			}

		case <-timeout:
			fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
			break streamTrade
		}
	}
	tradeCancel() // closes this stream on the server

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 4: STREAM 3 - OnPositionProfit (Position P&L stream)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 4: OnPositionProfit() - Stream position P&L updates")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 4.1. ON POSITION PROFIT
	//      Subscribe to real-time P&L changes for open positions.
	//      Returns: Profit/loss as market prices change.
	//      PERIODIC: One update every TimerPeriodMilliseconds (empty without open positions).
	//      Use cases: Real-time P&L monitoring, risk management.
	// ══════════════════════════════════════════════════════════════

	profitReq := &pb.OnPositionProfitRequest{
		TimerPeriodMilliseconds: 1000,  // one update per second (0 = the default, 1000 ms)
		IgnoreEmptyData:         false, // false = also send updates with no position in them
	}

	profitCtx, profitCancel := context.WithCancel(ctx)
	profitChan, profitErrChan := account.OnPositionProfit(profitCtx, profitReq)

	fmt.Printf("Streaming position P&L updates (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
	fmt.Println("  NOTE: This stream sends an update every TimerPeriodMilliseconds; with no open position the")
	fmt.Println("        update is empty (set IgnoreEmptyData: true to skip empty updates)")

	eventCount = 0
	timeout = time.After(MAX_SECONDS * time.Second)

streamProfit:
	for {
		select {
		case profitData, ok := <-profitChan:
			if !ok {
				fmt.Println("  Stream closed by server")
				break streamProfit
			}
			eventCount++
			// Direct field access to OnPositionProfitData arrays
			totalPositions := len(profitData.NewPositions) + len(profitData.UpdatedPositions) + len(profitData.DeletedPositions)
			fmt.Printf("  Event #%d: Type=%v Total=%d (New=%d Updated=%d Deleted=%d)\n",
				eventCount, profitData.Type, totalPositions,
				len(profitData.NewPositions),
				len(profitData.UpdatedPositions),
				len(profitData.DeletedPositions))

			// Show first 3 updated positions
			maxShow := 3
			if len(profitData.UpdatedPositions) > 0 {
				if len(profitData.UpdatedPositions) < maxShow {
					maxShow = len(profitData.UpdatedPositions)
				}
				for i := 0; i < maxShow; i++ {
					pos := profitData.UpdatedPositions[i]
					fmt.Printf("    Updated Position: Ticket=%d Symbol=%s Profit=%.2f\n",
						pos.Ticket,
						pos.PositionSymbol,
						pos.Profit)
				}
			}

			if eventCount >= MAX_EVENTS {
				fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
				break streamProfit
			}

		case err := <-profitErrChan:
			if err != nil {
				helpers.PrintShortError(err, "Stream error")
				break streamProfit
			}

		case <-timeout:
			fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
			break streamProfit
		}
	}
	profitCancel() // closes this stream on the server

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 5: STREAM 4 - OnPositionsAndPendingOrdersTickets (Tickets stream)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 5: OnPositionsAndPendingOrdersTickets() - Stream ticket changes")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 5.1. ON POSITIONS AND PENDING ORDERS TICKETS
	//      Subscribe to order/position ticket number changes.
	//      Returns: List of currently open tickets (minimal data).
	//      LIGHTWEIGHT: Only ticket IDs, no detailed info.
	//      Use cases: Quick "what's open?" check, ticket tracking.
	// ══════════════════════════════════════════════════════════════

	ticketsReq := &pb.OnPositionsAndPendingOrdersTicketsRequest{
		TimerPeriodMilliseconds: 1000, // one snapshot per second (0 = the default, 1000 ms)
	}

	ticketsCtx, ticketsCancel := context.WithCancel(ctx)
	ticketsChan, ticketsErrChan := account.OnPositionsAndPendingOrdersTickets(ticketsCtx, ticketsReq)

	fmt.Printf("Streaming ticket changes (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
	fmt.Println("  NOTE: This stream sends the current ticket lists every TimerPeriodMilliseconds,")
	fmt.Println("        whether they changed or not: compare them to the previous ones to see changes")

	eventCount = 0
	timeout = time.After(MAX_SECONDS * time.Second)

streamTickets:
	for {
		select {
		case ticketsData, ok := <-ticketsChan:
			if !ok {
				fmt.Println("  Stream closed by server")
				break streamTickets
			}
			eventCount++
			// Direct field access to OnPositionsAndPendingOrdersTicketsData
			totalTickets := len(ticketsData.PendingOrderTickets) + len(ticketsData.PositionTickets)
			fmt.Printf("  Event #%d: Total tickets=%d (Pending Orders=%d Positions=%d)\n",
				eventCount,
				totalTickets,
				len(ticketsData.PendingOrderTickets),
				len(ticketsData.PositionTickets))

			if eventCount >= MAX_EVENTS {
				fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
				break streamTickets
			}

		case err := <-ticketsErrChan:
			if err != nil {
				helpers.PrintShortError(err, "Stream error")
				break streamTickets
			}

		case <-timeout:
			fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
			break streamTickets
		}
	}
	ticketsCancel() // closes this stream on the server

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 6: STREAM 5 - OnTradeTransaction (Trade transaction events)
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nSTEP 6: OnTradeTransaction() - Stream trade transaction events")
	fmt.Println("───────────────────────────────────────────────────────────")

	// ══════════════════════════════════════════════════════════════
	// 6.1. ON TRADE TRANSACTION
	//      Subscribe to ALL trading transaction events.
	//      Returns: Low-level transaction details (order/deal tickets, type, price).
	//      COMPREHENSIVE: Every order placement, execution, modification, deletion.
	//      Use cases: Audit logs, detailed execution tracking, debugging.
	// ══════════════════════════════════════════════════════════════

	transactionReq := &pb.OnTradeTransactionRequest{}

	transactionCtx, transactionCancel := context.WithCancel(ctx)
	transactionChan, transactionErrChan := account.OnTradeTransaction(transactionCtx, transactionReq)

	fmt.Printf("Streaming trade transactions (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
	fmt.Println("  NOTE: This stream sends events for all trade-related transactions")

	eventCount = 0
	timeout = time.After(MAX_SECONDS * time.Second)

streamTransaction:
	for {
		select {
		case transactionData, ok := <-transactionChan:
			if !ok {
				fmt.Println("  Stream closed by server")
				break streamTransaction
			}
			eventCount++
			// Direct field access to nested TradeTransaction structure
			fmt.Printf("  Event #%d: Type=%v\n", eventCount, transactionData.Type)

			if transactionData.TradeTransaction != nil {
				tx := transactionData.TradeTransaction
				fmt.Printf("    Transaction Type=%v Order=%d Deal=%d\n",
					tx.Type,
					tx.OrderTicket,
					tx.DealTicket)

				if tx.Symbol != "" {
					fmt.Printf("    Symbol=%s Price=%.5f Volume=%.2f\n",
						tx.Symbol,
						tx.Price,
						tx.Volume)
				}
			}

			if eventCount >= MAX_EVENTS {
				fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
				break streamTransaction
			}

		case err := <-transactionErrChan:
			if err != nil {
				helpers.PrintShortError(err, "Stream error")
				break streamTransaction
			}

		case <-timeout:
			fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
			break streamTransaction
		}
	}
	transactionCancel() // closes this stream on the server

	// ═══════════════════════════════════════════════════════════════════════
	// FINAL: DISCONNECT
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n\nFinal: Disconnecting...")
	fmt.Println("───────────────────────────────────────────────────────────")

	// Disconnect (stops the terminal) and close the gRPC connection
	if err := account.DisconnectAndClose(); err != nil {
		helpers.PrintShortError(err, "Disconnect failed")
	} else {
		fmt.Println("✓ Disconnected, connection closed")
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("✓ DEMO COMPLETED SUCCESSFULLY")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("NOTE: Streaming methods are real-time and event-driven.")
	fmt.Println("      Some streams may not receive events if there's no activity.")
	fmt.Println("      For example:")
	fmt.Println("      - OnTrade sends events only when trades occur")
	fmt.Println("      - OnPositionProfit / OnPositionsAndPendingOrdersTickets send a snapshot every timer period")
	fmt.Println("        (empty without open positions/orders)")
	fmt.Println("      - OnTradeTransaction fires on any trading activity")

	return nil
}
