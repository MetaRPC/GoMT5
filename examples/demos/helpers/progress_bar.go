/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/demos/helpers/progress_bar.go
 PURPOSE:
   Utility package for displaying console progress bars during time-based waits.
   Used by orchestrators to visualize countdowns and waiting periods.

 USAGE EXAMPLE:
   // Wait 60 seconds with progress bar
   WaitWithProgressBar(
       60,                           // totalSeconds
       "Waiting for news event",     // message
       ctx,                          // context for cancellation
   )

   Output:
   Waiting for news event: [█████████░░░░░░░░░░] 45% (27s / 60s) - 33s remaining

══════════════════════════════════════════════════════════════════════════════*/

package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	// BarWidth defines the width of progress bar in characters
	BarWidth = 20

	// FilledChar is the character used for filled portion
	FilledChar = "█"

	// EmptyChar is the character used for empty portion
	EmptyChar = "░"
)

// WaitWithProgressBar displays a progress bar while waiting for a specified duration.
// Updates every second with visual progress indicator.
//
// Parameters:
//   - totalSeconds: Total seconds to wait
//   - message: Message to display before progress bar
//   - ctx: Context for cancellation
//
// Example:
//
//	ctx := context.Background()
//	WaitWithProgressBar(60, "Countdown to trade", ctx)
func WaitWithProgressBar(totalSeconds int, message string, ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	elapsed := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Println() // Move to next line
			return

		case <-ticker.C:
			// Display progress
			displayProgress(elapsed, totalSeconds, message)

			elapsed++
			if elapsed > totalSeconds {
				fmt.Println() // Move to next line after completion
				return
			}
		}
	}
}

// CountdownWithoutBar displays a simple countdown without progress bar.
//
// Parameters:
//   - totalSeconds: Total seconds to count down
//   - message: Message to display before countdown
//   - ctx: Context for cancellation
//
// Example:
//
//	CountdownWithoutBar(30, "Starting in", ctx)
func CountdownWithoutBar(totalSeconds int, message string, ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := totalSeconds

	for remaining >= 0 {
		select {
		case <-ctx.Done():
			fmt.Println()
			return

		case <-ticker.C:
			// Clear line and write countdown
			output := fmt.Sprintf("  %s: %ds remaining...", message, remaining)
			clearLine()
			fmt.Print(output)

			remaining--
			if remaining < 0 {
				fmt.Println()
				return
			}
		}
	}
}

// MonitorWithProgressBar displays a monitoring progress bar for operations with unknown duration.
// Updates periodically and shows elapsed time.
//
// Parameters:
//   - maxSeconds: Maximum seconds to monitor
//   - message: Message to display
//   - updateInterval: Update interval (e.g., 5*time.Second)
//   - ctx: Context for cancellation
//
// Example:
//
//	MonitorWithProgressBar(180, "Running strategy", 5*time.Second, ctx)
func MonitorWithProgressBar(maxSeconds int, message string, updateInterval time.Duration, ctx context.Context) {
	startTime := time.Now()
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return

		case <-ticker.C:
			elapsed := int(time.Since(startTime).Seconds())

			if elapsed >= maxSeconds {
				fmt.Println()
				return
			}

			remaining := maxSeconds - elapsed

			// Build output
			output := fmt.Sprintf("  %s: %ds elapsed / %ds max - %ds remaining",
				message, elapsed, maxSeconds, remaining)

			clearLine()
			fmt.Print(output)
		}
	}
}

// displayProgress displays the progress bar with percentage and time info
func displayProgress(elapsed, total int, message string) {
	// Calculate progress
	progress := float64(elapsed) / float64(total)
	filledWidth := int(progress * float64(BarWidth))
	emptyWidth := BarWidth - filledWidth

	// Build progress bar
	bar := strings.Repeat(FilledChar, filledWidth) + strings.Repeat(EmptyChar, emptyWidth)
	percent := int(progress * 100)
	remaining := total - elapsed

	// Build output string
	output := fmt.Sprintf("  %s: [%s] %d%% (%ds / %ds) - %ds remaining",
		message, bar, percent, elapsed, total, remaining)

	// Clear line and write progress
	clearLine()
	fmt.Print(output)
}

// clearLine clears the current console line
func clearLine() {
	// ANSI escape code to clear line: \r moves cursor to start, spaces overwrite
	const consoleWidth = 120 // Safe default width
	clearStr := "\r" + strings.Repeat(" ", consoleWidth) + "\r"
	fmt.Print(clearStr)
}

// WaitWithProgressBarAndCallback is an advanced version that calls a callback during wait.
// Useful for checking conditions while waiting.
//
// Parameters:
//   - totalSeconds: Total seconds to wait
//   - message: Message to display
//   - interval: How often to call callback
//   - callback: Function to call at each interval (return false to stop early)
//   - ctx: Context for cancellation
//
// Example:
//
//	WaitWithProgressBarAndCallback(60, "Monitoring", 5*time.Second, func() bool {
//	    // Check some condition
//	    if conditionMet {
//	        return false // Stop waiting early
//	    }
//	    return true // Continue waiting
//	}, ctx)
func WaitWithProgressBarAndCallback(
	totalSeconds int,
	message string,
	interval time.Duration,
	callback func() bool,
	ctx context.Context) {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	callbackTicker := time.NewTicker(interval)
	defer callbackTicker.Stop()

	elapsed := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return

		case <-callbackTicker.C:
			// Call callback
			if callback != nil && !callback() {
				fmt.Println() // Condition met, stop early
				return
			}

		case <-ticker.C:
			// Display progress
			displayProgress(elapsed, totalSeconds, message)

			elapsed++
			if elapsed > totalSeconds {
				fmt.Println()
				return
			}
		}
	}
}

// SpinnerWait displays a spinning animation while waiting.
// Useful for indeterminate waits.
//
// Parameters:
//   - message: Message to display
//   - ctx: Context for cancellation
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	SpinnerWait("Connecting to server", ctx)
func SpinnerWait(message string, ctx context.Context) {
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	i := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return

		case <-ticker.C:
			output := fmt.Sprintf("  %s %s", spinnerChars[i%len(spinnerChars)], message)
			clearLine()
			fmt.Print(output)
			i++
		}
	}
}
