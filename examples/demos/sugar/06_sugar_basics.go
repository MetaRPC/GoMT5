/*══════════════════════════════════════════════════════════════════════════════
 FILE: 06_sugar_basics.go - HIGH-LEVEL SUGAR API BASICS DEMO

 PURPOSE:
   Demonstrates MT5Sugar basics: connection, account balance, and price queries.
   This is the SIMPLEST way to interact with MT5 - perfect for quick scripts!

 🎯 WHO SHOULD USE THIS:
   • Beginners learning MT5 API
   • Quick prototyping and testing
   • Simple monitoring scripts
   • Educational purposes

 📚 WHAT THIS DEMO COVERS (3 Categories):

   1. CONNECTION METHODS (3 methods)
      • NewMT5Sugar() - Create and connect in one call
      • IsConnected() - Check connection status
      • Ping() - Verify connection health

   2. QUICK BALANCE METHODS (6 methods)
      • GetBalance() - Account balance
      • GetEquity() - Current equity
      • GetMargin() - Used margin
      • GetFreeMargin() - Available margin
      • GetMarginLevel() - Margin level %
      • GetProfit() - Floating P&L

   3. PRICES & QUOTES (5 methods)
      • GetBid() - Current BID price
      • GetAsk() - Current ASK price
      • GetSpread() - Spread in pips
      • GetPriceInfo() - Complete price info
      • WaitForPrice() - Wait for price update

 🔄 API LEVELS:
   HIGH-LEVEL (Sugar) - THIS FILE: One-liner operations, smart defaults
   MID-LEVEL (Service): More control, native Go types
   LOW-LEVEL (Account): Maximum flexibility, protobuf structures

 🚀 HOW TO RUN:
   cd examples/demos
   go run main.go 6    (or select [6] from menu)

  NOTE: This demo is READ-ONLY and safe to run on live accounts!
══════════════════════════════════════════════════════════════════════════════*/

package sugar

import (
	"fmt"
	"strings"
	"time"

	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	"github.com/MetaRPC/GoMT5/mt5"
)

func RunSugarBasicsDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("MT5 SUGAR API - BASICS DEMO (Read-Only Operations)")
	fmt.Println(strings.Repeat("=", 80))

	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load configuration")

	// ══════════════════════════════════════════════════════════════════════════════
	// 1. CONNECTION METHODS
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("1. CONNECTION METHODS")
	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 1.1. NewMT5Sugar()
	//      Create Sugar API instance with credentials.
	//      Chain: NewMT5Account() → NewMT5Service() → MT5Sugar
	//      Returns: *MT5Sugar instance ready for connection.
	//      SAFE operation - only initialization, no network calls.
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.1. NewMT5Sugar() - Create Sugar instance")

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	helpers.Fatal(err, "NewMT5Sugar failed")
	fmt.Println("  ✓ Sugar instance created!")

	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n1.2. QuickConnect() - Connect to MT5 terminal")

	err = sugar.QuickConnect(cfg.MtCluster)
	helpers.Fatal(err, "QuickConnect failed")
	fmt.Println("  ✓ Connected to MT5!")


	fmt.Println("\n1.3. IsConnected() - Check connection status")

	if sugar.IsConnected() {
		fmt.Println("  ✓ Status: CONNECTED")
	} else {
		fmt.Println("  ❌ Status: NOT CONNECTED")
	}


	fmt.Println("\n1.4. Ping() - Verify connection health")
	
	err = sugar.Ping()
	if !helpers.PrintShortError(err, "Ping failed") {
		fmt.Println("  ✓ Ping successful - connection is healthy")
	}

	// ══════════════════════════════════════════════════════════════════════════════
	// 2. QUICK BALANCE METHODS
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))

	fmt.Println("2. QUICK BALANCE METHODS")

	fmt.Println(strings.Repeat("=", 80))

	// ══════════════════════════════════════════════════════════════
	// 2.1. GetBalance()
	//      Get current account balance (deposit amount).
	//      Chain: Sugar → Service.GetAccountDouble() → Account → gRPC
	//      Returns: float64 balance in account currency.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.1. GetBalance() - Account balance")

	balance, err := sugar.GetBalance()
	if !helpers.PrintShortError(err, "GetBalance failed") {
		fmt.Printf("  Balance: %.2f\n", balance)
	}

	// ══════════════════════════════════════════════════════════════
	// 2.2. GetEquity()
	//      Get current account equity (balance + floating profit).
	//      Chain: Sugar → Service.GetAccountDouble() → Account → gRPC
	//      Returns: float64 equity = balance + open positions P/L.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.2. GetEquity() - Account equity")

	equity, err := sugar.GetEquity()
	if !helpers.PrintShortError(err, "GetEquity failed") {
		fmt.Printf("  Equity: %.2f\n", equity)
	}

	// ══════════════════════════════════════════════════════════════
	// 2.3. GetMargin()
	//      Get amount of margin currently used by open positions.
	//      Chain: Sugar → Service.GetAccountDouble() → Account → gRPC
	//      Returns: float64 used margin (collateral locked).
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.3. GetMargin() - Used margin")

	margin, err := sugar.GetMargin()
	if !helpers.PrintShortError(err, "GetMargin failed") {
		fmt.Printf("  Used Margin: %.2f\n", margin)
	}

	// ══════════════════════════════════════════════════════════════
	// 2.4. GetFreeMargin()
	//      Get amount of margin available for new positions.
	//      Chain: Sugar → Service.GetAccountDouble() → Account → gRPC
	//      Returns: float64 free margin = equity - used margin.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.4. GetFreeMargin() - Free margin")

	freeMargin, err := sugar.GetFreeMargin()
	if !helpers.PrintShortError(err, "GetFreeMargin failed") {
		fmt.Printf("  Free Margin: %.2f\n", freeMargin)
	}

	// ══════════════════════════════════════════════════════════════
	// 2.5. GetMarginLevel()
	//      Get margin level percentage (equity / margin * 100).
	//      Chain: Sugar → Service.GetAccountDouble() → Account → gRPC
	//      Returns: float64 margin level %. 0 = no positions.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.5. GetMarginLevel() - Margin level %")

	marginLevel, err := sugar.GetMarginLevel()
	if !helpers.PrintShortError(err, "GetMarginLevel failed") {
		if marginLevel == 0 {
			fmt.Println("  Margin Level: ∞ (no open positions)")
		} else {
			fmt.Printf("  Margin Level: %.2f%%\n", marginLevel)
		}
	}

	// ══════════════════════════════════════════════════════════════
	// 2.6. GetProfit()
	//      Get total floating profit/loss from all open positions.
	//      Chain: Sugar → Service.GetAccountDouble() → Account → gRPC
	//      Returns: float64 unrealized P/L. Positive = profit.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n2.6. GetProfit() - Current floating profit/loss")

	profit, err := sugar.GetProfit()
	if !helpers.PrintShortError(err, "GetProfit failed") {
		if profit >= 0 {
			fmt.Printf("  Profit: +%.2f\n", profit)
		} else {
			fmt.Printf("  Profit: %.2f\n", profit)
		}
	}

	// ══════════════════════════════════════════════════════════════════════════════
	// 3. PRICES & QUOTES METHODS
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("3. PRICES & QUOTES METHODS")
	fmt.Println(strings.Repeat("=", 80))

	testSymbol := cfg.TestSymbol

	// ══════════════════════════════════════════════════════════════
	// 3.1. GetBid()
	//      Get current BID price (SELL price) for symbol.
	//      Chain: Sugar → Service.GetSymbolTick() → Account → gRPC
	//      Returns: float64 current BID price.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.1. GetBid() - Current BID price for %s\n", testSymbol)

	bid, err := sugar.GetBid(testSymbol)
	if !helpers.PrintShortError(err, "GetBid failed") {
		fmt.Printf("  BID: %.5f\n", bid)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.2. GetAsk()
	//      Get current ASK price (BUY price) for symbol.
	//      Chain: Sugar → Service.GetSymbolTick() → Account → gRPC
	//      Returns: float64 current ASK price.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.2. GetAsk() - Current ASK price for %s\n", testSymbol)

	ask, err := sugar.GetAsk(testSymbol)
	if !helpers.PrintShortError(err, "GetAsk failed") {
		fmt.Printf("  ASK: %.5f\n", ask)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.3. GetSpread()
	//      Get current spread in points (not pips!).
	//      Chain: Sugar → Service.GetSymbolInteger() → Account → gRPC
	//      Returns: float64 spread in points (1 point = minimal price step).
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.3. GetSpread() - Current spread for %s\n", testSymbol)

	spread, err := sugar.GetSpread(testSymbol)
	if !helpers.PrintShortError(err, "GetSpread failed") {
		fmt.Printf("  Spread: %.0f points\n", spread)
	}

	// ══════════════════════════════════════════════════════════════
	// 3.4. GetPriceInfo()
	//      Get complete price information (BID, ASK, spread, time).
	//      Chain: Sugar → Service (GetSymbolTick + GetSymbolInteger) → Account
	//      Returns: *PriceInfo struct with all price data.
	//      SAFE operation - read-only query.
	//      Timeout: 3 seconds
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.4. GetPriceInfo() - Complete price info for %s\n", testSymbol)

	priceInfo, err := sugar.GetPriceInfo(testSymbol)
	if !helpers.PrintShortError(err, "GetPriceInfo failed") {
		fmt.Printf("  Symbol:     %s\n", testSymbol)
		fmt.Printf("  BID:        %.5f\n", priceInfo.Bid)
		fmt.Printf("  ASK:        %.5f\n", priceInfo.Ask)
		fmt.Printf("  Spread:     %.0f points\n", priceInfo.SpreadPips)
		fmt.Printf("  Time:       %s\n", priceInfo.Time.Format("2006-01-02 15:04:05"))
	}

	// ══════════════════════════════════════════════════════════════
	// 3.5. WaitForPrice()
	//      Wait for valid price update (BID > 0 and ASK > 0).
	//      Chain: Sugar → Service.GetSymbolTick() (polling) → Account → gRPC
	//      Returns: *PriceInfo when valid price received.
	//      SAFE operation - read-only query with polling.
	//      Timeout: Custom (here 3 seconds)
	// ══════════════════════════════════════════════════════════════
	fmt.Printf("\n3.5. WaitForPrice() - Wait for price update (3 sec timeout)\n")

	fmt.Printf("  Waiting for %s price change...\n", testSymbol)
	tick, err := sugar.WaitForPrice(testSymbol, 3*time.Second)
	if helpers.PrintShortError(err, "WaitForPrice timeout or error") {
		fmt.Println("  ℹ️  This is normal if price doesn't change within timeout")
	} else {
		fmt.Printf("  ✓ Price received!\n")
		fmt.Printf("    BID: %.5f, ASK: %.5f\n", tick.Bid, tick.Ask)
	}

	// ══════════════════════════════════════════════════════════════════════════════
	// DEMO SUMMARY
	// ══════════════════════════════════════════════════════════════════════════════
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ DEMO COMPLETED SUCCESSFULLY!")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n📊 METHODS DEMONSTRATED (14 total):")
	fmt.Println("   CONNECTION (4):    NewMT5Sugar, QuickConnect, IsConnected, Ping")
	fmt.Println("   BALANCE (6):       GetBalance, GetEquity, GetMargin, GetFreeMargin,")
	fmt.Println("                      GetMarginLevel, GetProfit")
	fmt.Println("   PRICES (5):        GetBid, GetAsk, GetSpread, GetPriceInfo, WaitForPrice")

	fmt.Println("\n🎯 MT5Sugar API ADVANTAGES:")
	fmt.Println("   ✓ One-liner operations:   sugar.GetBalance() - simple!")
	fmt.Println("   ✓ Direct value returns:   float64, not wrapped in structs")
	fmt.Println("   ✓ No protobuf knowledge:  All native Go types")
	fmt.Println("   ✓ Auto timeouts:          3-30 seconds built-in")
	fmt.Println("   ✓ Smart defaults:         EURUSD base symbol, standard parameters")

	fmt.Println("\n📚 WHAT'S NEXT:")
	fmt.Println("   • Trading operations:     Run demo 07 - sugar_trading.go")
	fmt.Println("   • Position management:    Run demo 08 - sugar_positions.go")
	fmt.Println("   • History & profit:       Run demo 09 - sugar_history.go")

	fmt.Println("\n" + strings.Repeat("=", 80))
}
