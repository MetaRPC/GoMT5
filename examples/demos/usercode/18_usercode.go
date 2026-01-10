// ══════════════════════════════════════════════════════════════════════════════
// FILE: 18_usercode.go - USER CODE SANDBOX
// ══════════════════════════════════════════════════════════════════════════════
//
// 🎯 WHAT IS THIS?
//   Your personal sandbox for writing and testing MT5 trading code.
//   Write once, run instantly with: go run main.go 18
//
//  HOW TO ACTIVATE:
//   1. Open main.go and uncomment line ~100:
//      "github.com/MetaRPC/GoMT5/examples/demos/usercode"
//
//   2. Uncomment lines ~342-344 in main.go:
//      case "18", "usercode", "user", "sandbox", "custom":
//          usercode.RunUserCode()
//          return false, nil
//
//   3. Run: cd examples/demos && go run main.go 18
//  
// 🔑 THREE API LEVELS AVAILABLE:
//   • account (MT5Account) - Low-level gRPC protobuf (full control)
//   • service (MT5Service) - Mid-level native Go types (recommended)
//   • sugar   (MT5Sugar)   - High-level convenience (62+ methods)
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
	"github.com/MetaRPC/GoMT5/mt5"
	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
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
//   account.PositionsGet(ctx, &pb.Empty{})
//   account.SymbolInfoTickRequest(ctx, &pb.SymbolInfoTickRequestRequest{Symbol: "EURUSD"})
//   account.OrderSend(ctx, &pb.OrderSendRequest{Request: &pb.MqlTradeRequest{...}})
//
// LEVEL 2: MT5Service (Mid-level Go types)
// ─────────────────────────────────────────────────────────────────────────────
//   service.GetAccountSummary(ctx)
//   service.GetOpenPositions(ctx)
//   service.GetSymbolInfo(ctx, "EURUSD")
//   service.PlaceMarketOrder(ctx, "EURUSD", 0.01, orderType, sl, tp)
//
// LEVEL 3: MT5Sugar (High-level one-liners)
// ─────────────────────────────────────────────────────────────────────────────
//   sugar.GetBalance() / GetEquity() / GetMargin()
//   sugar.GetBid("EURUSD") / GetAsk("EURUSD") / GetSpread("EURUSD")
//   sugar.BuyMarket("EURUSD", 0.01)
//   sugar.BuyMarketWithPips("EURUSD", 0.01, slPips, tpPips)
//   sugar.GetOpenPositions() / ClosePosition(ticket) / CloseAllPositions()
//   sugar.CalculatePositionSize("EURUSD", riskPercent, slPips)
//
// ══════════════════════════════════════════════════════════════════════════════
