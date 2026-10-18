package helpers

/*══════════════════════════════════════════════════════════════════════════════
 FILE: helpers/error_helper.go - ERROR HANDLING HELPERS

 PURPOSE:
   Provides centralized error handling utilities for MT5 API demos.
   Simplifies error checking and formatting across all demo files.

 KEY FUNCTIONS:
   • Fatal() - Print error, run the OnFatal hooks (disconnect) and exit (for critical errors)
   • OnFatal() / DisconnectOnExit() - Make sure the terminal is disconnected however a demo ends
   • DisconnectSugar() - Disconnect and print the result
   • PrintIfError() - Print error and continue (for non-critical errors)
   • PrintSuccess() - Print success message
   • PrintWarning() - Print warning message
   • FormatApiError() - Format ApiError with full details

 USAGE EXAMPLE:

   // Critical error - exits program
   summaryData, err := account.AccountSummary(ctx, req)
   helpers.Fatal(err, "Failed to get account summary")

   // Non-critical error - continues execution
   marginData, err := account.SymbolInfoMarginRate(ctx, req)
   if !helpers.PrintIfError(err, "Margin rate not available (broker doesn't support)") {
       fmt.Printf("Rate: %.2f\n", marginData.InitialMarginRate)
   }

══════════════════════════════════════════════════════════════════════════════*/

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"

	mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
)

// ══════════════════════════════════════════════════════════════════════════════
// CRITICAL ERROR HANDLING - EXITS PROGRAM
// ══════════════════════════════════════════════════════════════════════════════

// exitHooks run (last registered first) before Fatal - or Ctrl+C - ends the program.
var (
	exitHooksMu     sync.Mutex
	exitHooks       []*func()
	interruptHooked sync.Once
	// gracefulInterrupt: the running code (an orchestrator or preset) handles the first Ctrl+C itself
	// (it stops, then its deferred disconnect runs); only a second Ctrl+C runs the exit hooks and exits.
	gracefulInterrupt atomic.Bool
)

// LetCallerHandleInterrupt makes the first Ctrl+C the running code's own business: an orchestrator or preset
// that stops gracefully on Ctrl+C (deleting its pending orders, ...) must not be raced by the immediate
// disconnect-and-exit of the exit hooks. A second Ctrl+C still runs the hooks and exits. The returned function
// switches it off again (use it with defer).
func LetCallerHandleInterrupt() (restore func()) {
	gracefulInterrupt.Store(true)
	return func() { gracefulInterrupt.Store(false) }
}

// OnFatal registers a function that runs before Fatal exits the program (and on Ctrl+C).
// os.Exit skips deferred calls, so cleanup such as a disconnect has to be registered here.
// The returned function unregisters it. Most demos use DisconnectOnExit instead.
func OnFatal(hook func()) (unregister func()) {
	h := &hook
	exitHooksMu.Lock()
	exitHooks = append(exitHooks, h)
	exitHooksMu.Unlock()

	interruptHooked.Do(func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		go func() {
			first := true
			for range sig {
				if first && gracefulInterrupt.Load() {
					first = false
					continue // the orchestrator stops by itself; its deferred disconnect follows
				}
				fmt.Println("\nInterrupted - cleaning up...")
				runExitHooks()
				os.Exit(130)
			}
		}()
	})

	return func() {
		exitHooksMu.Lock()
		defer exitHooksMu.Unlock()
		for i, x := range exitHooks {
			if x == h {
				exitHooks = append(exitHooks[:i], exitHooks[i+1:]...)
				return
			}
		}
	}
}

// runExitHooks runs and clears the registered exit hooks, last registered first.
func runExitHooks() {
	exitHooksMu.Lock()
	hooks := exitHooks
	exitHooks = nil
	exitHooksMu.Unlock()
	for i := len(hooks) - 1; i >= 0; i-- {
		(*hooks[i])()
	}
}

// DisconnectOnExit makes sure the demo's terminal is disconnected however the demo ends: normally
// (the returned function, for defer), through Fatal, a panic or Ctrl+C. Call it right after creating
// the Sugar/Service - before connecting - because a failed connect may still have started the terminal:
//
//   sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
//   helpers.Fatal(err, "NewMT5Sugar failed")
//   defer helpers.DisconnectOnExit(sugar)()
func DisconnectOnExit(d Disconnecter) func() {
	var once sync.Once
	disconnect := func() { once.Do(func() { DisconnectSugar(d) }) }
	unregister := OnFatal(disconnect)
	return func() {
		unregister()
		disconnect()
	}
}

// Fatal prints detailed error information and exits the program.
// Use this for critical errors where continuing is impossible.
// Functions registered with OnFatal (e.g. a disconnect) run before the exit.
//
// Example:
//   cfg, err := config.LoadConfig()
//   helpers.Fatal(err, "Failed to load configuration")
func Fatal(err error, context string) {
	if err == nil {
		return
	}

	// Check if this is an ApiError from MT5 server
	var msg string
	var apiErr *mt5.ApiError
	if errors.As(err, &apiErr) {
		msg = fmt.Sprintf("❌ FATAL: %s\n%s", context, FormatApiError(apiErr))
	} else {
		// Regular error
		msg = fmt.Sprintf("❌ FATAL: %s: %v", context, err)
	}
	log.Print(msg)

	runExitHooks()
	os.Exit(1)
}

// Disconnecter is anything with the MT5Sugar/MT5Service Disconnect method.
type Disconnecter interface {
	Disconnect(deleteOnDisconnect ...bool) error
}

// DisconnectSugar stops the demo's terminal on the server and closes the connection, printing the result.
// Safe to call more than once (later calls do nothing). Use it with defer and OnFatal so every demo
// run ends with its terminal disconnected - also when the demo fails.
func DisconnectSugar(d Disconnecter) {
	if d == nil {
		return
	}
	if err := d.Disconnect(); err != nil {
		PrintWarning(fmt.Sprintf("Disconnect failed: %v", err))
		return
	}
	fmt.Println("🔌 Disconnected (terminal stopped)")
}

// ══════════════════════════════════════════════════════════════════════════════
// NON-CRITICAL ERROR HANDLING - CONTINUES EXECUTION
// ══════════════════════════════════════════════════════════════════════════════

// PrintIfError prints error details and returns true if error exists.
// Use this for non-critical operations where the program should continue.
//
// Returns:
//   true if error occurred (use with if statement to skip success path)
//   false if no error (continue with normal flow)
//
// Example:
//   marginData, err := account.SymbolInfoMarginRate(ctx, req)
//   if !helpers.PrintIfError(err, "Margin rate not available") {
//       fmt.Printf("Rate: %.2f\n", marginData.InitialMarginRate)
//   }
func PrintIfError(err error, context string) bool {
	if err == nil {
		return false
	}

	// Check if this is an ApiError from MT5 server
	var apiErr *mt5.ApiError
	if errors.As(err, &apiErr) {
		fmt.Printf("  ⚠️  %s\n", context)

		// Print trade error if present (most specific)
		if apiErr.MqlErrorTradeIntCode() != 0 {
			fmt.Printf("      Trade Error: %s (%d)\n",
				apiErr.MqlErrorTradeDescription(), apiErr.MqlErrorTradeIntCode())
			fmt.Printf("      Code: %s\n", apiErr.MqlErrorTradeCode())
		} else if apiErr.MqlErrorIntCode() != 0 {
			// Print MQL error if present
			fmt.Printf("      MQL Error: %s (%d)\n",
				apiErr.MqlErrorDescription(), apiErr.MqlErrorIntCode())
			fmt.Printf("      Code: %s\n", apiErr.MqlErrorCode())
		} else {
			// Generic API error
			fmt.Printf("      API Error: %s\n", apiErr.ErrorCode())
		}

		return true
	}

	// Regular error - simple format
	fmt.Printf("  ⚠️  %s: %v\n", context, err)
	return true
}

// PrintShortError prints a brief error message (clean and readable).
// Use when you want minimal output without protobuf dumps.
//
// Example:
//   helpers.PrintShortError(err, "Balance retrieval failed")
func PrintShortError(err error, context string) bool {
	if err == nil {
		return false
	}

	var apiErr *mt5.ApiError
	if errors.As(err, &apiErr) {
		// For trade errors, show the trade error description (most user-friendly)
		if apiErr.MqlErrorTradeIntCode() != 0 {
			fmt.Printf("  ❌ %s: %s (%s)\n",
				context,
				apiErr.MqlErrorTradeDescription(),
				apiErr.MqlErrorTradeCode())
		} else if apiErr.MqlErrorIntCode() != 0 {
			// For MQL errors, show MQL error description
			fmt.Printf("  ❌ %s: %s (%s)\n",
				context,
				apiErr.MqlErrorDescription(),
				apiErr.MqlErrorCode())
		} else {
			// Generic API error
			fmt.Printf("  ❌ %s: %s\n", context, apiErr.ErrorCode())
		}
	} else {
		fmt.Printf("  ❌ %s: %v\n", context, err)
	}

	return true
}

// ══════════════════════════════════════════════════════════════════════════════
// NOTE: PrintSuccess, PrintError, PrintWarning are defined in connection.go
// ══════════════════════════════════════════════════════════════════════════════

// ══════════════════════════════════════════════════════════════════════════════
// DETAILED ERROR FORMATTING
// ══════════════════════════════════════════════════════════════════════════════

// FormatApiError returns a detailed multi-line string representation of ApiError.
// Use when you need full error details for logging or debugging.
//
// Example:
//   if apiErr, ok := err.(*mt5.ApiError); ok {
//       log.Println(helpers.FormatApiError(apiErr))
//   }
func FormatApiError(apiErr *mt5.ApiError) string {
	if apiErr == nil {
		return "ApiError{nil}"
	}

	var result string

	// Trade error details (most specific and user-friendly)
	if apiErr.MqlErrorTradeIntCode() != 0 {
		result = fmt.Sprintf("Trade Error: %s (%d)\n",
			apiErr.MqlErrorTradeDescription(),
			apiErr.MqlErrorTradeIntCode())
		result += fmt.Sprintf("Code: %s\n", apiErr.MqlErrorTradeCode())
	} else if apiErr.MqlErrorIntCode() != 0 {
		// MQL error details
		result = fmt.Sprintf("MQL Error: %s (%d)\n",
			apiErr.MqlErrorDescription(),
			apiErr.MqlErrorIntCode())
		result += fmt.Sprintf("Code: %s\n", apiErr.MqlErrorCode())
	} else {
		// Generic API error
		result = fmt.Sprintf("API Error: %s\n", apiErr.ErrorCode())
	}

	// Command info (if available)
	if apiErr.CommandTypeName() != "" {
		result += fmt.Sprintf("Command: %s\n", apiErr.CommandTypeName())
	}

	return result
}

// ══════════════════════════════════════════════════════════════════════════════
// TRADE RETURN CODE HELPERS
// ══════════════════════════════════════════════════════════════════════════════

// CheckRetCode checks trade operation return code and prints appropriate message.
// Returns true if operation succeeded (RetCode == 10009).
//
// Example:
//   placeData, err := account.OrderSend(ctx, req)
//   helpers.Fatal(err, "OrderSend failed")
//
//   if helpers.CheckRetCode(placeData.ReturnedCode, "Order placement") {
//       fmt.Printf("Order ticket: %d\n", placeData.OrderTicket)
//   }
func CheckRetCode(retCode uint32, operation string) bool {
	if mt5.IsRetCodeSuccess(retCode) {
		fmt.Printf("  ✓ %s successful (RetCode: %d)\n", operation, retCode)
		return true
	}

	fmt.Printf("  ❌ %s failed (RetCode: %d)\n", operation, retCode)
	fmt.Printf("     %s\n", mt5.GetRetCodeMessage(retCode))

	// Provide helpful hints for common errors
	switch retCode {
	case mt5.TradeRetCodeNoMoney:
		fmt.Println("     💡 Hint: Check account margin - insufficient funds")
	case mt5.TradeRetCodeInvalidStops:
		fmt.Println("     💡 Hint: SL/TP too close to market price - check SYMBOL_TRADE_STOPS_LEVEL")
	case mt5.TradeRetCodeInvalidVolume:
		fmt.Println("     💡 Hint: Check SYMBOL_VOLUME_MIN, SYMBOL_VOLUME_MAX, SYMBOL_VOLUME_STEP")
	case mt5.TradeRetCodeMarketClosed:
		fmt.Println("     💡 Hint: Market is closed - check trading hours")
	case mt5.TradeRetCodeRequote, mt5.TradeRetCodePriceChanged:
		fmt.Println("     💡 Hint: Price changed - retry with updated price")
	}

	return false
}

// PrintRetCodeWarning prints a warning for non-critical trade return codes.
// Use when operation partially succeeded or requires retry.
//
// Example:
//   if mt5.IsRetCodeRequote(retCode) {
//       helpers.PrintRetCodeWarning(retCode, "Price changed - retrying...")
//   }
func PrintRetCodeWarning(retCode uint32, context string) {
	fmt.Printf("  ⚠️  %s (RetCode: %d)\n", context, retCode)
	fmt.Printf("     %s\n", mt5.GetRetCodeMessage(retCode))
}

// ══════════════════════════════════════════════════════════════════════════════
// CONNECTION ERROR HANDLING - ErrNotConnected
// ══════════════════════════════════════════════════════════════════════════════

// CheckConnection checks if error is ErrNotConnected and prints appropriate message.
// Returns true if connection error was detected.
//
// This function demonstrates proper usage of errors.Is() with sentinel errors.
//
// Example:
//   data, err := account.OrderSend(ctx, req)
//   if helpers.CheckConnection(err) {
//       // Handle reconnection logic here
//       return fmt.Errorf("connection lost, please reconnect")
//   }
func CheckConnection(err error) bool {
	if errors.Is(err, mt5.ErrNotConnected) {
		fmt.Println("  ⚠️  Connection Error: Not connected to MT5 terminal")
		fmt.Println("     Please ensure you called Connect() or ConnectEx() first")
		return true
	}
	return false
}
