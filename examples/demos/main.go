/*══════════════════════════════════════════════════════════════════════════════
 FILE: main.go - MAIN ENTRY POINT FOR ALL EXAMPLES

 ╔═══════════════════════════════════════════════════════════════════════════╗
 ║                        GoMT5 - MetaTrader 5 gRPC Client                   ║
 ╚═══════════════════════════════════════════════════════════════════════════╝

 This is the ENTRY POINT for all demonstration examples.
 Provides an interactive menu to showcase how the MT5 gRPC API works at different
 abstraction levels.

 ╔═══════════════════════════════════════════════════════════════════════════╗
 ║                        LEARNING PROGRESSION                               ║
 ╚═══════════════════════════════════════════════════════════════════════════╝

 Follow this order to understand the architecture from ground up:

 1️⃣ LOW-LEVEL (MT5Account)
    → examples/demos/lowlevel/
    → Direct gRPC/protobuf calls
    → Raw data structures
    → Understand the foundation

 2️⃣ MID-LEVEL (MT5Service) - RECOMMENDED
    → examples/demos/service/
    → Wrapper facade over MT5Account
    → Cleaner API, same functionality
    → Organized method groups

 3️⃣ HIGH-LEVEL (MT5Sugar - Extension Methods)
    → examples/demos/sugar/
    → One-liner operations
    → Risk-based calculations
    → Smart defaults and normalization

 4️⃣ ORCHESTRATORS (Complete Strategies)
    → examples/demos/orchestrators/
    → Automated trade lifecycle
    → Multiple API calls combined
    → Entry → Monitor → Exit logic

 5️⃣ PRESETS (Ready-to-Use Systems)
    → examples/demos/Presets/
    → Full trading systems
    → Adaptive market analysis
    → Multiple orchestrators combined

 ╔═══════════════════════════════════════════════════════════════════════════╗
 ║                              USAGE                                        ║
 ╚═══════════════════════════════════════════════════════════════════════════╝

 Interactive Menu:
   go run main.go                → Shows menu with all examples

 Direct Launch (skip menu):
   go run main.go lowlevel01     → Low-level general operations
   go run main.go service        → Mid-level Service API
   go run main.go sugar06        → High-level Sugar basics
   go run main.go grid           → Grid Trading orchestrator
   go run main.go adaptive       → Adaptive Market Preset

 Available Commands:
   lowlevel01, lowlevel02, lowlevel03, service, service05,
   sugar06, sugar07, sugar08, sugar09,
   grid, trailing, scaler, risk, rebalancer, adaptive

 ╔═══════════════════════════════════════════════════════════════════════════╗
 ║                         PROJECT STRUCTURE                                 ║
 ╚═══════════════════════════════════════════════════════════════════════════╝

 examples/demos/
 ├── main.go                   ← THIS FILE (entry point)
 ├── config/config.json        ← Connection settings
 ├── helpers/                  ← Connection helper
 ├── lowlevel/                 ← Direct gRPC examples [1-3]
 ├── service/                  ← MT5Service examples [4-5]
 ├── sugar/                    ← Sugar API examples [6-9]
 ├── orchestrators/            ← Trading strategies [10-14]
 ├── Presets/                  ← Complete systems [15]
 └── errors/                   ← Error handling utilities

══════════════════════════════════════════════════════════════════════════════*/

package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/MetaRPC/GoMT5/examples/demos/config"
	"github.com/MetaRPC/GoMT5/examples/demos/helpers"
	"github.com/MetaRPC/GoMT5/examples/demos/lowlevel"
	"github.com/MetaRPC/GoMT5/examples/demos/orchestrators"
	"github.com/MetaRPC/GoMT5/examples/demos/presets"
	"github.com/MetaRPC/GoMT5/examples/demos/service"
	"github.com/MetaRPC/GoMT5/examples/demos/sugar"
	"github.com/MetaRPC/GoMT5/examples/demos/usercode"
	"github.com/MetaRPC/GoMT5/mt5"
)

func main() {
	// Главный цикл
	for {
		var command string

		// Получить команду из аргументов или показать меню
		if len(os.Args) > 1 {
			// Режим прямой команды: go run main.go <command>
			command = strings.ToLower(os.Args[1])
		} else {
			// Режим интерактивного меню
			printBanner()
			command = showMenu()
		}

		// Выполнить команду
		exitRequested, err := executeCommand(command)
		if err != nil {
			fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
			fmt.Println("║                    ERROR OCCURRED                          ║")
			fmt.Println("╚════════════════════════════════════════════════════════════╝")
			fmt.Printf("\nError: %v\n", err)

			if len(os.Args) > 1 {
				os.Exit(1)
			}

			fmt.Println("\nPress Enter to return to menu...")
			fmt.Scanln()
			continue
		}

		if exitRequested {
			fmt.Println("\nExiting...")
			return
		}

		// Режим командной строки: выход после одного запуска
		if len(os.Args) > 1 {
			fmt.Println("\n\nPress Enter to exit...")
			fmt.Scanln()
			return
		}

		// Интерактивный режим: спросить продолжить
		fmt.Println()
		fmt.Println("┌──────────────────────────────────────────────────────────────┐")
		fmt.Println("│  [M] Return to Main Menu  |  [0] Exit                        │")
		fmt.Println("└──────────────────────────────────────────────────────────────┘")
		fmt.Print("\nYour choice: ")

		var nextAction string
		fmt.Scanln(&nextAction)
		nextAction = strings.ToLower(strings.TrimSpace(nextAction))
		fmt.Println()

		if nextAction == "0" || nextAction == "exit" || nextAction == "quit" {
			fmt.Println("Exiting...")
			return
		}

		// Продолжить к меню
	}
}

// ═════════════════════════════════════════════════════════════════
// BANNER
// ═════════════════════════════════════════════════════════════════
func printBanner() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                  ║")
	fmt.Println("║              GoMT5 - MetaTrader 5 gRPC Client                    ║")
	fmt.Println("║                                                                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// ═════════════════════════════════════════════════════════════════
// MENU
// ═════════════════════════════════════════════════════════════════
func showMenu() string {
	fmt.Println("┌──────────────────────────────────────────────────────────────────┐")
	fmt.Println("│  LOW-LEVEL EXAMPLES (Direct gRPC)                                │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  [1]  General Operations   → go run main.go lowlevel01           │")
	fmt.Println("│  [2]  Trading Operations   → go run main.go lowlevel02           │")
	fmt.Println("│  [3]  Streaming Methods    → go run main.go lowlevel03           │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  MID-LEVEL (MT5Service Wrapper) - RECOMMENDED                    │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  [4]  Service Examples     → go run main.go service              │")
	fmt.Println("│  [5]  Streaming Methods    → go run main.go service05            │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  HIGH-LEVEL (Sugar API - Simplest)                               │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  [6]  Basics (Connection)  → go run main.go sugar06              │")
	fmt.Println("│  [7]  Trading Operations   → go run main.go sugar07              │")
	fmt.Println("│  [8]  Position Management  → go run main.go sugar08              │")
	fmt.Println("│  [9]  History & Profits    → go run main.go sugar09              │")
	fmt.Println("│  [10] Advanced Features    → go run main.go sugar10              │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  ORCHESTRATORS (Automated Strategies)                            │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  [11] Trailing Stop        → go run main.go trailing             │")
	fmt.Println("│  [12] Position Scaler      → go run main.go scaler               │")
	fmt.Println("│  [13] Grid Trading         → go run main.go grid                 │")
	fmt.Println("│  [14] Risk Manager         → go run main.go risk                 │")
	fmt.Println("│  [15] Portfolio Rebalancer → go run main.go rebalancer           │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  PRESETS & TOOLS                                                 │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  [16] Adaptive Preset      → go run main.go adaptive             │")
	fmt.Println("│  [17] Protobuf Inspector   → go run main.go inspect              │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  USER CODE SANDBOX (Uncomment in main.go to enable)              │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  [18] User Code Sandbox    → go run main.go 18       (DISABLED)  │")
	fmt.Println("│      See examples/demos/usercode/README.md to enable             │")
	fmt.Println("├──────────────────────────────────────────────────────────────────┤")
	fmt.Println("│  [0]  EXIT                                                       │")
	fmt.Println("└──────────────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Print("Enter your choice: ")

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)
	fmt.Println()

	return input
}

// ═════════════════════════════════════════════════════════════════
// #region COMMAND EXECUTOR
// ═════════════════════════════════════════════════════════════════
func executeCommand(command string) (exitRequested bool, err error) {
	command = strings.ToLower(command)

	switch command {

	// ═════════════════════════════════════════════════════════════
	// LOW-LEVEL EXAMPLES
	// ═════════════════════════════════════════════════════════════
	case "1", "lowlevel01":
		return false, lowlevel.RunGeneral01()

	case "2", "lowlevel02":
		return false, lowlevel.RunTrading02()

	case "3", "lowlevel03", "streaming", "stream":
		return false, lowlevel.RunStreaming03()

	// ═════════════════════════════════════════════════════════════
	// MID-LEVEL (SERVICE)
	// ═════════════════════════════════════════════════════════════
	case "4", "service", "mid":
		return false, service.RunService04()

	case "5", "service05", "servicestreaming", "service-streaming":
		return false, service.RunServiceStreaming05()

	// ═════════════════════════════════════════════════════════════
	// HIGH-LEVEL (SUGAR)
	// ═════════════════════════════════════════════════════════════
	case "6", "sugar06", "sugarbasics", "basics":
		sugar.RunSugarBasicsDemo()
		return false, nil

	case "7", "sugar07", "sugartrading", "trading":
		sugar.RunSugarTradingDemo()
		return false, nil

	case "8", "sugar08", "sugarpositions", "positions":
		sugar.RunSugarPositionsDemo()
		return false, nil

	case "9", "sugar09", "sugarhistory", "history":
		sugar.RunSugarHistoryDemo()
		return false, nil

	case "10", "sugar10", "sugaradvanced", "advanced":
		sugar.RunSugarAdvancedDemo()
		return false, nil

	// ═════════════════════════════════════════════════════════════
	// ORCHESTRATORS
	// ═════════════════════════════════════════════════════════════
	case "11", "trailing", "trailingstop":
		return false, RunOrchestrator_TrailingStop()

	case "12", "scaler", "positionscaler":
		return false, RunOrchestrator_PositionScaler()

	case "13", "grid", "gridtrading":
		return false, RunOrchestrator_Grid()

	case "14", "risk", "riskmanager":
		return false, RunOrchestrator_RiskManager()

	case "15", "rebalancer", "portfolio":
		return false, RunOrchestrator_PortfolioRebalancer()

	// ═════════════════════════════════════════════════════════════
	// PRESETS & TOOLS
	// ═════════════════════════════════════════════════════════════
	case "16", "adaptive", "preset":
		return false, RunOrchestrator_AdaptivePreset()

	case "17", "inspect", "inspector", "proto":
		helpers.RunProtobufInspector()
		return false, nil

	// ═════════════════════════════════════════════════════════════
	// USER CODE SANDBOX (DISABLED BY DEFAULT)
	// ═════════════════════════════════════════════════════════════
	// 🎓 YOUR PERSONAL SANDBOX FOR LEARNING & TESTING MT5 CODE
	//
	// INSTRUCTIONS TO ENABLE:
	//   1. Uncomment line ~100 in this file:
	//      // "github.com/MetaRPC/GoMT5/examples/demos/usercode"
	//
	//   2. Uncomment the case block below (lines ~326-328):
	//      // case "18", "usercode", "user", "sandbox", "custom":
	//      // 	usercode.RunUserCode()
	//      // 	return false, nil
	//
	//   3. Run: go run main.go 18
	//
	//   4. Edit your code in: examples/demos/usercode/18_usercode.go
	//      (File has educational examples showing all 3 API levels!)
	//
	// WHAT YOU'LL LEARN:
	//   • MT5Account (low-level) - Raw gRPC protocol
	//   • MT5Service (mid-level) - Native Go types
	//   • MT5Sugar (high-level)  - 62+ convenience methods
	//   • How to mix all 3 levels for maximum power!
	// ═════════════════════════════════════════════════════════════

	case "18", "usercode", "user", "sandbox", "custom":
		usercode.RunUserCode()
		return false, nil

	// ═════════════════════════════════════════════════════════════
	// EXIT
	// ═════════════════════════════════════════════════════════════
	case "0", "x", "exit", "quit":
		return true, nil

	// ═════════════════════════════════════════════════════════════
	// INVALID COMMAND
	// ═════════════════════════════════════════════════════════════
	default:
		fmt.Printf("\n❌ Invalid command: %s\n", command)
		fmt.Println("Valid commands: 1-15, 0 (exit)")
		fmt.Println("Or use keywords:")
		fmt.Println("  Low-level:      lowlevel01, lowlevel02, lowlevel03")
		fmt.Println("  Service:        service, service05")
		fmt.Println("  Sugar:          sugar06, sugar07, sugar08, sugar09")
		fmt.Println("  Orchestrators:  trailing, scaler, grid, risk, rebalancer")
		fmt.Println("  Presets:        adaptive")
		return false, nil
	}
}

// #endregion

// ═════════════════════════════════════════════════════════════════
// #region ORCHESTRATOR RUNNERS
//
// Each runner creates a configured orchestrator with customizable
// parameters.
// You can modify these settings before running the orchestrator.
// ═════════════════════════════════════════════════════════════════

// RunOrchestrator_TrailingStop demonstrates trailing stop manager.
// Automatically adjusts stop loss as position moves into profit.
func RunOrchestrator_TrailingStop() error {
	fmt.Println("\n=== TRAILING STOP MANAGER ===")
	fmt.Println("Automated trailing stop management for all positions")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	if err != nil {
		return fmt.Errorf("failed to create MT5Sugar: %w", err)
	}

	err = sugar.QuickConnect(cfg.MtCluster)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sugar.GetService().GetAccount().Close()

	// ╔════════════════════════════════════════════════════════════╗
	// ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
	// ╚════════════════════════════════════════════════════════════╝
	orchConfig := orchestrators.TrailingStopConfig{
		TrailingDistance: 200,             // Trail 200 points behind
		ActivationProfit: 300,             // Activate after 300 points profit
		UpdateInterval:   2 * time.Second, // Check every 2 seconds
		Symbols:          []string{},      // All symbols (empty = all)
		MinDistance:      100,             // Min 100 points from price
		StepSize:         50,              // Adjust in 50 point steps
	}

	fmt.Println("\n📋 Configuration:")
	fmt.Printf("  Activation Profit:  %d points\n", int(orchConfig.ActivationProfit))
	fmt.Printf("  Trailing Distance:  %d points\n", int(orchConfig.TrailingDistance))
	fmt.Printf("  Update Interval:    %v\n", orchConfig.UpdateInterval)

	tsManager := orchestrators.NewTrailingStopManager(sugar, orchConfig)

	fmt.Println("\n🚀 Starting Trailing Stop Manager...")
	if err := tsManager.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	// Run for 3 minutes with live progress bar and metrics
	fmt.Println("  ✓ Starting monitoring...")
	fmt.Println()

	helpers.WaitWithProgressBarAndCallback(
		180,                      // 3 minutes = 180 seconds
		"Trailing Stop Active",
		2 * time.Second,          // Update callback every 2 seconds (matches UpdateInterval)
		func() bool {
			// Display live metrics during operation
			status := tsManager.GetStatus()
			metrics := tsManager.GetMetrics()

			fmt.Printf("\r  📊 Positions: %d | SL Updates: %d | Success: %d | Errors: %d        ",
				metrics.CurrentPositions,
				status.SuccessCount,
				status.SuccessCount,
				status.ErrorCount)

			return true // Continue monitoring
		},
		tsManager.GetContext(),
	)
	fmt.Println() // New line after progress bar completes

	fmt.Println("\n🛑 Stopping...")
	if err := tsManager.Stop(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	showOrchestratorMetrics(tsManager)
	return nil
}

// RunOrchestrator_PositionScaler demonstrates position pyramiding.
// Adds to winning positions as they move in profit.
func RunOrchestrator_PositionScaler() error {
	fmt.Println("\n=== POSITION SCALER (PYRAMIDING) ===")
	fmt.Println("Automatically scales into profitable positions")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	if err != nil {
		return fmt.Errorf("failed to create MT5Sugar: %w", err)
	}

	err = sugar.QuickConnect(cfg.MtCluster)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sugar.GetService().GetAccount().Close()

	// ╔════════════════════════════════════════════════════════════╗
	// ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
	// ╚════════════════════════════════════════════════════════════╝
	orchConfig := orchestrators.PositionScalerConfig{
		Mode:             orchestrators.Pyramiding, // Pyramiding mode
		TriggerDistance:  200,                      // Scale every 200 points
		MinProfitToScale: 100,                      // Min 100 points profit
		InitialLotSize:   0.10,                     // Start with 0.10 lots
		ScaleLotSize:     0.05,                     // Add 0.05 lots per scale
		MaxScales:        3,                        // Max 3 scale-ins
		TotalMaxLotSize:  1.0,                      // Max 1.0 total lots
		StopLossPerScale: 150,                      // 150 points SL per scale
		Symbols:          []string{cfg.TestSymbol}, // EURUSD
		CheckInterval:    5 * time.Second,
	}

	fmt.Println("\n📋 Configuration:")
	fmt.Printf("  Mode:            Pyramiding\n")
	fmt.Printf("  Symbol:          %s\n", cfg.TestSymbol)
	fmt.Printf("  Initial Size:    %.2f lots\n", orchConfig.InitialLotSize)
	fmt.Printf("  Scale Size:      %.2f lots\n", orchConfig.ScaleLotSize)
	fmt.Printf("  Max Scales:      %d\n", orchConfig.MaxScales)
	fmt.Printf("  Trigger:         %d points profit\n", int(orchConfig.TriggerDistance))

	scaler := orchestrators.NewPositionScaler(sugar, orchConfig)

	fmt.Println("\n🚀 Starting Position Scaler...")
	if err := scaler.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	// Run for 5 minutes with live progress bar and metrics
	fmt.Println("  ✓ Starting monitoring...")
	fmt.Println()

	helpers.WaitWithProgressBarAndCallback(
		300,                      // 5 minutes = 300 seconds
		"Position Scaler Active",
		5 * time.Second,          // Update callback every 5 seconds (matches CheckInterval)
		func() bool {
			// Display live metrics during operation
			status := scaler.GetStatus()
			metrics := scaler.GetMetrics()

			fmt.Printf("\r  📊 Groups: %d | Scales: %d | Trades: %d | Errors: %d        ",
				metrics.CurrentPositions,
				status.SuccessCount,
				metrics.TotalTrades,
				status.ErrorCount)

			return true // Continue monitoring
		},
		scaler.GetContext(),
	)
	fmt.Println() // New line after progress bar completes

	fmt.Println("\n🛑 Stopping...")
	if err := scaler.Stop(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	showOrchestratorMetrics(scaler)
	return nil
}

// RunOrchestrator_Grid demonstrates grid trading strategy.
// Places buy and sell orders at fixed intervals around current price.
func RunOrchestrator_Grid() error {
	fmt.Println("\n=== GRID TRADING ===")
	fmt.Println("Automated grid trading with fixed intervals")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	if err != nil {
		return fmt.Errorf("failed to create MT5Sugar: %w", err)
	}

	err = sugar.QuickConnect(cfg.MtCluster)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sugar.GetService().GetAccount().Close()

	// ╔════════════════════════════════════════════════════════════╗
	// ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
	// ╚════════════════════════════════════════════════════════════╝
	orchConfig := orchestrators.GridTraderConfig{
		Symbol:         cfg.TestSymbol,
		GridSize:       5,               // 5 levels above and below
		GridStep:       100,             // 100 points (10 pips) spacing
		LotSize:        0.01,            // 0.01 lots per order
		MaxPositions:   10,              // Max 10 concurrent positions
		TakeProfit:     0,               // Use grid step as TP
		StopLoss:       0,               // No stop loss
		CheckInterval:  5 * time.Second, // Check every 5 seconds
		RebuildOnFill:  false,           // Don't rebuild grid on fill
	}

	fmt.Println("\n📋 Configuration:")
	fmt.Printf("  Symbol:        %s\n", orchConfig.Symbol)
	fmt.Printf("  Grid Size:     %d levels (each side)\n", orchConfig.GridSize)
	fmt.Printf("  Grid Step:     %d points\n", int(orchConfig.GridStep))
	fmt.Printf("  Lot Size:      %.2f\n", orchConfig.LotSize)
	fmt.Printf("  Max Positions: %d\n", orchConfig.MaxPositions)

	gridTrader := orchestrators.NewGridTrader(sugar, orchConfig)

	fmt.Println("\n🚀 Starting Grid Trader...")
	if err := gridTrader.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	// Run for 10 minutes with live progress bar and metrics
	fmt.Println("  ✓ Starting monitoring...")
	fmt.Println()

	helpers.WaitWithProgressBarAndCallback(
		600,                      // 10 minutes = 600 seconds
		"Grid Trader Active",
		5*time.Second, // Update callback every 5 seconds (matches CheckInterval)
		func() bool {
			// Display live metrics during operation
			status := gridTrader.GetStatus()
			metrics := gridTrader.GetMetrics()

			fmt.Printf("\r  📊 Positions: %d | Orders: %d | Trades: %d | Errors: %d        ",
				metrics.CurrentPositions,
				status.SuccessCount, // Active pending orders count
				metrics.TotalTrades,
				status.ErrorCount)

			return true // Continue monitoring
		},
		gridTrader.GetContext(),
	)
	fmt.Println() // New line after progress bar completes

	fmt.Println("\n🛑 Stopping...")
	if err := gridTrader.Stop(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	showOrchestratorMetrics(gridTrader)
	return nil
}

// RunOrchestrator_RiskManager demonstrates account-level risk management.
// Monitors account metrics and enforces risk limits.
func RunOrchestrator_RiskManager() error {
	fmt.Println("\n=== RISK MANAGER ===")
	fmt.Println("Account-level risk monitoring and enforcement")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	if err != nil {
		return fmt.Errorf("failed to create MT5Sugar: %w", err)
	}

	err = sugar.QuickConnect(cfg.MtCluster)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sugar.GetService().GetAccount().Close()

	// ╔════════════════════════════════════════════════════════════╗
	// ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
	// ╚════════════════════════════════════════════════════════════╝
	orchConfig := orchestrators.RiskManagerConfig{
		MaxDrawdownPercent:  10.0,             // 10% max drawdown
		MaxDrawdownAbsolute: 1000.0,           // $1000 max drawdown
		DailyLossLimit:      500.0,            // $500 daily loss limit
		DailyProfitTarget:   1000.0,           // $1000 daily profit target
		MinMarginLevel:      150.0,            // 150% min margin level
		MaxMarginUsed:       80.0,             // 80% max margin usage
		MaxOpenPositions:    20,               // Max 20 positions
		MaxSymbolExposure:   5,                // Max 5 per symbol
		MaxPositionSize:     1.0,              // Max 1.0 lot
		CheckInterval:       5 * time.Second,  // Check every 5 seconds
		EnableAutoClose:     true,             // Auto-close on breach
		EnableTradeBlocking: true,             // Block trades on breach
	}

	fmt.Println("\n📋 Configuration:")
	fmt.Printf("  Max Drawdown:      %.1f%% or $%.2f\n", orchConfig.MaxDrawdownPercent, orchConfig.MaxDrawdownAbsolute)
	fmt.Printf("  Daily Loss Limit:  $%.2f\n", orchConfig.DailyLossLimit)
	fmt.Printf("  Min Margin Level:  %.1f%%\n", orchConfig.MinMarginLevel)
	fmt.Printf("  Max Positions:     %d\n", orchConfig.MaxOpenPositions)
	fmt.Printf("  Auto-Close:        %v\n", orchConfig.EnableAutoClose)

	riskManager := orchestrators.NewRiskManager(sugar, orchConfig)

	fmt.Println("\n🚀 Starting Risk Manager...")
	if err := riskManager.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	// Run for 15 minutes with live progress bar and risk metrics
	fmt.Println("  ✓ Starting monitoring...")
	fmt.Println()

	helpers.WaitWithProgressBarAndCallback(
		900,                      // 15 minutes = 900 seconds
		"Risk Manager Active",
		5*time.Second, // Update callback every 5 seconds (matches CheckInterval)
		func() bool {
			// Display live risk metrics during operation
			metrics := riskManager.GetMetrics()
			todayProfit := riskManager.GetTodayProfit()
			isBlocked := riskManager.IsTradingBlocked()
			riskEvents := len(riskManager.GetRiskEvents())

			blockedStr := "✅"
			if isBlocked {
				blockedStr = "🔒"
			}

			profitSign := ""
			if todayProfit > 0 {
				profitSign = "+"
			}

			fmt.Printf("\r  🛡️ DD: %.1f%% | Daily: %s%.2f | Events: %d | Trading: %s        ",
				metrics.MaxDrawdown,   // Current drawdown percentage
				profitSign,
				todayProfit,           // Today's profit/loss
				riskEvents,            // Risk events count
				blockedStr)            // Trading status

			return true // Continue monitoring
		},
		riskManager.GetContext(),
	)
	fmt.Println() // New line after progress bar completes

	fmt.Println("\n🛑 Stopping...")
	if err := riskManager.Stop(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	showOrchestratorMetrics(riskManager)
	return nil
}

// RunOrchestrator_PortfolioRebalancer demonstrates portfolio rebalancing.
// Maintains target allocations across multiple symbols.
func RunOrchestrator_PortfolioRebalancer() error {
	fmt.Println("\n=== PORTFOLIO REBALANCER ===")
	fmt.Println("Automatic portfolio rebalancing across symbols")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	if err != nil {
		return fmt.Errorf("failed to create MT5Sugar: %w", err)
	}

	err = sugar.QuickConnect(cfg.MtCluster)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sugar.GetService().GetAccount().Close()

	// ╔════════════════════════════════════════════════════════════╗
	// ║  CONFIGURATION - MODIFY THESE SETTINGS                     ║
	// ╚════════════════════════════════════════════════════════════╝
	orchConfig := orchestrators.PortfolioRebalancerConfig{
		Allocations: map[string]float64{
			"EURUSD": 0.25, // 25%
			"GBPUSD": 0.25, // 25%
			"USDJPY": 0.25, // 25%
			"XAUUSD": 0.25, // 25%
		},
		TotalExposure:      10000.0,          // $10k total exposure
		MaxDeviation:       0.20,             // 20% max deviation
		RebalanceThreshold: 10.0,             // Rebalance at 10% off
		CheckInterval:      30 * time.Minute, // Check every 30 min
		MinPositionSize:    0.01,
		UseMarketOrders:    true,
		SlippageTolerance:  50,
		MaxTradesPerCycle:  10,
	}

	fmt.Println("\n📋 Configuration:")
	fmt.Println("  Portfolio:")
	for symbol, pct := range orchConfig.Allocations {
		fmt.Printf("    %s: %.0f%%\n", symbol, pct*100)
	}
	fmt.Printf("  Total Exposure:    $%.2f\n", orchConfig.TotalExposure)
	fmt.Printf("  Rebalance at:      %.0f%% deviation\n", orchConfig.RebalanceThreshold)

	rebalancer := orchestrators.NewPortfolioRebalancer(sugar, orchConfig)

	fmt.Println("\n🚀 Starting Portfolio Rebalancer...")
	if err := rebalancer.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	// Run for 1 hour with progress bar
	fmt.Println("  ✓ Starting monitoring...")
	fmt.Println()

	helpers.WaitWithProgressBarAndCallback(
		3600,                      // 1 hour = 3600 seconds
		"Portfolio Rebalancer Active",
		10*time.Second, // Update callback every 10 seconds
		func() bool {
			// Display live portfolio metrics
			rebalanceCount := rebalancer.GetRebalanceCount()
			allocations := rebalancer.GetCurrentAllocations()

			// Count symbols needing adjustment
			needsAdjustment := 0
			maxDeviation := 0.0
			for _, alloc := range allocations {
				if alloc.NeedsAdjustment {
					needsAdjustment++
				}
				if math.Abs(alloc.DeviationPercent) > maxDeviation {
					maxDeviation = math.Abs(alloc.DeviationPercent)
				}
			}

			balancedStr := "✅"
			if needsAdjustment > 0 {
				balancedStr = "⚠️"
			}

			fmt.Printf("\r  📊 Rebalances: %d | Max Deviation: %.1f%% | Unbalanced: %d | Status: %s        ",
				rebalanceCount,
				maxDeviation,
				needsAdjustment,
				balancedStr)

			return true // Continue monitoring
		},
		rebalancer.GetContext(),
	)
	fmt.Println() // New line after progress bar completes

	fmt.Println("\n🛑 Stopping...")
	if err := rebalancer.Stop(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	showOrchestratorMetrics(rebalancer)
	return nil
}

// RunOrchestrator_AdaptivePreset demonstrates adaptive multi-strategy system.
// Automatically selects best orchestrator based on market conditions.
func RunOrchestrator_AdaptivePreset() error {
	fmt.Println("\n=== ADAPTIVE ORCHESTRATOR PRESET ===")
	fmt.Println("Intelligent multi-strategy system with market regime detection")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
	if err != nil {
		return fmt.Errorf("failed to create MT5Sugar: %w", err)
	}

	err = sugar.QuickConnect(cfg.MtCluster)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sugar.GetService().GetAccount().Close()

	// Create and execute preset
	preset := presets.NewAdaptiveOrchestratorPreset(sugar)

	fmt.Println("\n🚀 Starting Adaptive Preset...")
	fmt.Println("  ✓ Press Ctrl+C to stop")

	totalProfit, err := preset.Execute()
	if err != nil {
		return fmt.Errorf("preset execution failed: %w", err)
	}

	fmt.Printf("\n\n📊 Final Result: $%.2f\n", totalProfit)
	return nil
}

// ═════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ═════════════════════════════════════════════════════════════════

// showOrchestratorMetrics displays final orchestrator metrics.
func showOrchestratorMetrics(orch orchestrators.Orchestrator) {
	status := orch.GetStatus()
	metrics := orch.GetMetrics()

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("FINAL METRICS")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Printf("\n📊 Orchestrator: %s\n", status.Name)
	fmt.Printf("   Uptime:       %s\n", orchestrators.FormatDuration(status.Uptime))
	fmt.Printf("   Operations:   %d total (%d success, %d failed)\n",
		status.SuccessCount+status.ErrorCount,
		status.SuccessCount,
		status.ErrorCount)

	if metrics.TotalTrades > 0 {
		fmt.Printf("\n💰 Trading Stats:\n")
		fmt.Printf("   Total Trades:     %d\n", metrics.TotalTrades)
		fmt.Printf("   Winning:          %d (%.1f%%)\n", metrics.WinningTrades, metrics.WinRate)
		fmt.Printf("   Losing:           %d\n", metrics.LosingTrades)
		fmt.Printf("   Net Profit:       %.2f\n", metrics.NetProfit)
		fmt.Printf("   Current Positions: %d\n", metrics.CurrentPositions)
	}

	if status.LastError != "" {
		fmt.Printf("\n⚠️  Last Error: %s\n", status.LastError)
	}

	fmt.Println()
}

// #endregion