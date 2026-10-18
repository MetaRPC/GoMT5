// ══════════════════════════════════════════════════════════════════════════════
// FILE: 18_usercode.go - USER CODE SANDBOX
// ══════════════════════════════════════════════════════════════════════════════
//
// 🎯 WHAT IS THIS?
//   Your personal sandbox for writing and testing MT5 trading code.
//   Write once, run instantly with: go run main.go 18
//
//  HOW TO RUN (nothing to enable - case 18 is wired in main.go):
//   cd examples/demos && go run main.go 18
//   (aliases: usercode, user, sandbox, custom)
//   Put your code below "YOUR CODE HERE"; the terminal is disconnected when
//   RunUserCode returns, also after helpers.Fatal, a panic or Ctrl+C.
//  
// 🔑 THREE API LEVELS AVAILABLE:
//   • account (MT5Account) - Low-level gRPC protobuf (full control)
//   • service (MT5Service) - Mid-level native Go types (recommended)
//   • sugar   (MT5Sugar)   - High-level convenience (60+ methods)
//
//   You can mix all 3 levels in the same code!
//
// ══════════════════════════════════════════════════════════════════════════════

package usercode

import (
	"context"
	"fmt"
	"strings"

	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
	pb "github.com/MetaRPC/GoMT5/package"
)

// RunUserCode - Your sandbox function
func RunUserCode() {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║           USER CODE SANDBOX - Ready to Code!          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// INITIALIZATION - All 3 API levels ready
	// ═══════════════════════════════════════════════════════════
	cfg, err := config.LoadConfig()
	helpers.Fatal(err, "Failed to load config")

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	helpers.Fatal(err, "Failed to create Sugar")
	sugar.GetAccount().ApiKey = cfg.ApiKey
	// Disconnect however the sandbox ends (normal end, Fatal, panic, Ctrl+C)
	defer helpers.DisconnectOnExit(sugar)()

	err = sugar.QuickConnect(cfg.MtCluster)
	helpers.Fatal(err, "Failed to connect")

	service := sugar.GetService() // Mid-level
	account := sugar.GetAccount() // Low-level
	ctx := context.Background()

	// Prevent unused warnings (pb is used in commented Example 3)
	_, _, _, _ = account, service, ctx, pb.ErrorType_UNDEFINED

	fmt.Println("✅ Connected! All 3 API levels initialized:")
	fmt.Println("   • account (MT5Account) - Low-level gRPC")
	fmt.Println("   • service (MT5Service) - Mid-level Go types")
	fmt.Println("   • sugar   (MT5Sugar)   - High-level convenience")
	fmt.Println()
	fmt.Println(strings.Repeat("═", 55))
	fmt.Println("YOUR CODE STARTS HERE ↓")
	fmt.Println(strings.Repeat("═", 55))
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// QUICK START EXAMPLES (Uncomment to try)
	// ═══════════════════════════════════════════════════════════

	// Example 1: Get account balance (Sugar - easiest)
	// ─────────────────────────────────────────────────────────
	// balance, err := sugar.GetBalance()
	// if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	// }
	// fmt.Printf("Balance: %.2f\n", balance)

	// Example 2: Get account info (Service - more control)
	// ─────────────────────────────────────────────────────────
	// accountInfo, err := service.GetAccountSummary(ctx)
	// if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	// }
	// fmt.Printf("Balance: %.2f %s\n", accountInfo.Balance, accountInfo.Currency)
	// fmt.Printf("Equity:  %.2f\n", accountInfo.Equity)

	// Example 3: Get account via protobuf (Account - full control)
	// ─────────────────────────────────────────────────────────
	// req := &pb.AccountSummaryRequest{}
	// reply, err := account.AccountSummary(ctx, req)
	// if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	// }
	// fmt.Printf("Balance (protobuf): %.2f\n", reply.AccountBalance)

	// ═══════════════════════════════════════════════════════════
	// YOUR CODE HERE
	// ═══════════════════════════════════════════════════════════

	// TODO: Your code here

	fmt.Println("\n" + strings.Repeat("═", 55))
	fmt.Println("User code completed!")
	fmt.Println(strings.Repeat("═", 55) + "\n")
}

// ══════════════════════════════════════════════════════════════════════════════
// QUICK REFERENCE
// ══════════════════════════════════════════════════════════════════════════════
//
// LEVEL 1: MT5Account (Low-level protobuf)
// ─────────────────────────────────────────────────────────────────────────────
//   account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
//   account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
//       InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC})
//   account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: "EURUSD"})
//   account.OrderSend(ctx, &pb.OrderSendRequest{Symbol: "EURUSD",
//       Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY, Volume: 0.01})
//
// LEVEL 2: MT5Service (Mid-level Go types)
// ─────────────────────────────────────────────────────────────────────────────
//   service.GetAccountSummary(ctx)
//   service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
//   service.GetSymbolTick(ctx, "EURUSD")
//   service.GetSymbolParamsMany(ctx, &symbolName, nil, nil, nil)
//   service.PlaceOrder(ctx, &pb.OrderSendRequest{Symbol: "EURUSD",
//       Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY, Volume: 0.01})
//
// LEVEL 3: MT5Sugar (High-level one-liners)
// ─────────────────────────────────────────────────────────────────────────────
//   sugar.GetBalance() / GetEquity() / GetMargin()
//   sugar.GetBid("EURUSD") / GetAsk("EURUSD") / GetSpread("EURUSD")   (spread in points)
//   sugar.BuyMarket("EURUSD", 0.01)
//   sugar.BuyMarketWithPips("EURUSD", 0.01, slPoints, tpPoints)       (distances in POINTS)
//   sugar.GetOpenPositions() / ClosePosition(ticket) / CloseAllPositions()
//   sugar.CalculatePositionSize("EURUSD", riskPercent, slPoints)      (distance in POINTS)
//
// ══════════════════════════════════════════════════════════════════════════════
