package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	// Create new account session object
	account, err := NewMT4Account(501401178, "v8gctta", "", uuid.Nil)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}

	// Connect to terminal by host/port
	err = account.ConnectByServerName(ctx, "RoboForex-Demo", "EURUSD", true, 30)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Get account summary
	summary, err := account.AccountSummary(ctx)
	if err != nil {
		log.Fatalf("Failed to get account summary: %v", err)
	}
	fmt.Printf("Account balance: %.2f\n", summary.GetAccountBalance())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start streaming ticks for these symbols
	tickCh, errCh := account.OnSymbolTick(ctx, []string{"EURUSD", "GBPUSD"})

	// Process data and errors as they arrive
	for {
		select {
		case tick, ok := <-tickCh:
			if !ok {
				fmt.Println("Tick channel closed (stream ended).")
				return
			}
			// Print the tick (customize fields as needed)
			fmt.Printf("Tick: Symbol=%s, Bid=%.5f, Ask=%.5f",
				tick.GetSymbolTick(), tick.SymbolTick.GetBid(), tick.SymbolTick.GetAsk())
		case err, ok := <-errCh:
			if !ok {
				fmt.Println("Error channel closed.")
				return
			}
			log.Printf("Stream error: %v\n", err)
			return
		case <-time.After(30 * time.Second):
			fmt.Println("No ticks received in 30 seconds. Cancelling...")
			cancel()
			return
		}
	}
}
