package helpers

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/MetaRPC/GoMT5/examples/demos/config"
	mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
	"github.com/google/uuid"
)

// Centralized connection management for all demo programs
// CreateAndConnectAccount creates and connects an MT5 account
func CreateAndConnectAccount() (*mt5.MT5Account, *config.MT5Config, error) {
	PrintSection("CONNECTION")

	// Loading configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Printf("  User:          %d\n", cfg.User)
	fmt.Printf("  gRPC Server:   %s\n", cfg.GrpcServer)
	fmt.Printf("  MT Cluster:    %s\n", cfg.MtCluster)
	fmt.Printf("  Base Symbol:   %s\n", cfg.TestSymbol)

	// Generating GUID (will be updated after ConnectEx)
	sessionId := uuid.New()
	fmt.Printf("  Generated Session ID: %s\n", sessionId)

	// Creating MT5Account
	account, err := mt5.NewMT5Account(cfg.User, cfg.Password, cfg.GrpcServer, sessionId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create MT5Account: %w", err)
	}

	fmt.Println("\n→ Connecting to MT5 terminal...")
	fmt.Printf("  Method: ConnectEx (via server name)\n")
	fmt.Printf("  Server: %s\n", cfg.MtCluster)

	// Connection
	err = ConnectByServerName(account, cfg.MtCluster, cfg.TestSymbol, 120)
	if err != nil {
		account.Close()
		return nil, nil, fmt.Errorf("connection failed: %w", err)
	}

	fmt.Printf("\033[32m  ✓ Connected successfully!\033[0m\n\n")

	return account, cfg, nil
}

// ConnectByServerName connects to MT5 using the server name
func ConnectByServerName(account *mt5.MT5Account, serverName, baseSymbol string, timeoutSeconds int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds+30)*time.Second)
	defer cancel()

	req := &pb.ConnectExRequest{
		User:            account.User,
		Password:        account.Password,
		MtClusterName:   serverName,
		BaseChartSymbol: &baseSymbol,
	}

	// ConnectEx
	reply, err := account.ConnectEx(ctx, req)
	if err != nil {
		return fmt.Errorf("ConnectEx failed: %w", err)
	}

	// CRITICAL: Updating the GUID with a value from the server
	account.Id = uuid.MustParse(reply.TerminalInstanceGuid)

	// Connection check
	checkCtx, checkCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer checkCancel()

	checkData, err := account.CheckConnect(checkCtx, &pb.CheckConnectRequest{})
	if err != nil {
		return fmt.Errorf("connection verification failed: %w", err)
	}

	if checkData.HealthCheck == nil || !checkData.HealthCheck.IsAlive {
		return fmt.Errorf("terminal is not alive after connection")
	}

	return nil
}

// PrintSection displays the section header
func PrintSection(title string) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ %-64s ║\n", title)
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
}

// PrintSuccess displays a success message (in green)
func PrintSuccess(msg string) {
	fmt.Printf("\033[32m%s\033[0m", msg)
}

// PrintError displays an error message (in red)
func PrintError(msg string) {
	fmt.Printf("\033[31m%s\033[0m\n", msg)
}

// PrintWarning displays a warning (in yellow)
func PrintWarning(msg string) {
	fmt.Printf("\033[33m%s\033[0m\n", msg)
}
