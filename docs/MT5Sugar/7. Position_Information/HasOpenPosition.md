# ✔️ Has Open Position (`HasOpenPosition`)

> **Sugar method:** Checks if at least one open position exists for a specific symbol.

**API Information:**

* **Method:** `sugar.HasOpenPosition(symbol)`
* **Timeout:** 5 seconds (it calls `GetPositionsBySymbol()`, which calls `GetOpenPositions()`)
* **Returns:** Boolean (true if at least one position is open on that symbol)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) HasOpenPosition(symbol string) (bool, error)
```

Implementation: `len(GetPositionsBySymbol(symbol)) > 0`. The full position list is fetched and filtered in Go
(exact, case-sensitive symbol comparison). Pending orders are not considered.

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| `symbol` | `string` | Trading symbol to check (e.g., "EURUSD") |

| Output | Type | Description |
|--------|------|-------------|
| `bool` | boolean | `true` if at least one position is open on `symbol`, `false` if none |
| `error` | `error` | The `GetOpenPositions` error if the check failed (`false` is returned then) |

---

## 💬 Just the Essentials

* **What it is:** Boolean check - "Do I have an open position on this symbol?"
* **Why you need it:** Avoid opening a duplicate position on a symbol, check a symbol before closing it.
* **Sanity check:** Same cost as `GetPositionsBySymbol(symbol)` - it fetches all open positions. For "any position on any symbol" use `CountOpenPositions() > 0`.

---

## 🎯 When to Use

✅ **One position per symbol** - Skip a new entry if the symbol already has a position

✅ **Before trading** - Verify a symbol is flat before starting a strategy on it

✅ **End of day** - Confirm a symbol is closed

✅ **Validation** - Simple yes/no for a symbol

---

## 🔗 Usage Examples

### 1) Basic usage - check one symbol

```go
hasPos, err := sugar.HasOpenPosition("EURUSD")
if err != nil {
    fmt.Printf("Check failed: %v\n", err)
    return
}

if hasPos {
    fmt.Println("✅ You have open EURUSD positions")
} else {
    fmt.Println("❌ No open EURUSD positions")
}
```

---

### 2) Verify clean slate before trading a symbol

```go
hasPos, _ := sugar.HasOpenPosition("EURUSD")

if hasPos {
    fmt.Println("⚠️  WARNING: EURUSD positions already open!")
    fmt.Println("   Close existing positions before starting strategy")
    return
}

fmt.Println("✅ Clean slate - ready to trade EURUSD")
// Start trading strategy
```

---

### 3) End of day validation for a list of symbols

```go
func ValidateEndOfDay(sugar *mt5.MT5Sugar, symbols []string) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      END OF DAY VALIDATION            ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    for _, symbol := range symbols {
        hasPos, err := sugar.HasOpenPosition(symbol)
        if err != nil {
            return fmt.Errorf("validation failed: %w", err)
        }

        if hasPos {
            fmt.Printf("❌ ALERT: %s positions still open!\n", symbol)
            fmt.Println("   Action required: Close them")
            return fmt.Errorf("%s positions still open", symbol)
        }
    }

    fmt.Println("✅ All clear - no open positions on the checked symbols")
    fmt.Println("   Safe to shutdown")
    return nil
}

// Usage: Run before weekend shutdown
err := ValidateEndOfDay(sugar, []string{"EURUSD", "GBPUSD"})
```

---

### 4) Polling for a symbol to close

```go
// Wait for all EURUSD positions to close
fmt.Println("Waiting for EURUSD positions to close...")

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

timeout := time.After(5 * time.Minute)

for {
    select {
    case <-timeout:
        fmt.Println("⏰ Timeout - positions still open")
        return

    case <-ticker.C:
        hasPos, err := sugar.HasOpenPosition("EURUSD")
        if err != nil {
            fmt.Printf("Check failed: %v\n", err)
            continue
        }

        if !hasPos {
            fmt.Println("✅ All EURUSD positions closed!")
            return
        }

        fmt.Println("Positions still open - waiting...")
    }
}
```

---

### 5) One position per symbol

```go
symbol := "EURUSD"
volume := 0.1

// Check if already trading this symbol
hasPos, _ := sugar.HasOpenPosition(symbol)

if hasPos {
    fmt.Printf("⚠️  Already have an active %s position\n", symbol)
    fmt.Println("   Strategy: Only one position per symbol allowed")
    return
}

fmt.Printf("✅ No active %s position - opening trade\n", symbol)
ticket, _ := sugar.BuyMarket(symbol, volume)
fmt.Printf("Position #%d opened\n", ticket)
```

---

### 6) Monitoring for exit signal

```go
// Close a symbol's positions when a condition is met
symbol := "EURUSD"

ticker := time.NewTicker(10 * time.Second)
defer ticker.Stop()

for range ticker.C {
    hasPos, _ := sugar.HasOpenPosition(symbol)

    if !hasPos {
        fmt.Println("No positions - waiting for entry signal...")
        continue
    }

    // Check exit condition (example: time-based, local clock)
    currentHour := time.Now().Hour()

    if currentHour >= 17 { // Close at 5 PM
        fmt.Printf("🕔 17:00 - Closing all %s positions\n", symbol)
        closed, _ := sugar.CloseAllBySymbol(symbol)
        fmt.Printf("Closed %d positions\n", closed)
        return
    }

    fmt.Println("Positions active - monitoring...")
}
```

---

### 7) System health check

```go
func SystemHealthCheck(sugar *mt5.MT5Sugar, symbol string) {
    fmt.Println("System Health Check:")
    fmt.Println("─────────────────────────────")

    // Check connection
    connected := sugar.IsConnected()
    if connected {
        fmt.Println("✅ Connection: OK")
    } else {
        fmt.Println("❌ Connection: FAILED")
    }

    // Check positions on the symbol
    hasPos, err := sugar.HasOpenPosition(symbol)
    if err != nil {
        fmt.Println("❌ Position check: FAILED")
    } else if hasPos {
        fmt.Printf("⚠️  %s positions: OPEN\n", symbol)
    } else {
        fmt.Printf("✅ %s positions: NONE\n", symbol)
    }

    // Check balance
    balance, _ := sugar.GetBalance()
    fmt.Printf("💰 Balance: $%.2f\n", balance)
}

// Usage:
SystemHealthCheck(sugar, "EURUSD")
```

---

### 8) Strategy state machine

```go
type StrategyState int

const (
    StateIdle StrategyState = iota
    StateTrading
    StateClosing
)

func GetStrategyState(sugar *mt5.MT5Sugar, symbol string) StrategyState {
    hasPos, err := sugar.HasOpenPosition(symbol)
    if err != nil {
        return StateIdle
    }

    if hasPos {
        return StateTrading
    }

    return StateIdle
}

// Usage:
state := GetStrategyState(sugar, "EURUSD")

switch state {
case StateIdle:
    fmt.Println("State: IDLE - Looking for entries")
case StateTrading:
    fmt.Println("State: TRADING - Managing positions")
case StateClosing:
    fmt.Println("State: CLOSING - Exiting positions")
}
```

---

### 9) Auto-restart after close

```go
func MonitorAndRestart(sugar *mt5.MT5Sugar, symbol string) {
    fmt.Printf("Monitoring %s positions - will restart when all closed\n", symbol)

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    wasOpen := false

    for range ticker.C {
        hasPos, _ := sugar.HasOpenPosition(symbol)

        if hasPos {
            wasOpen = true
            fmt.Println("Positions open - monitoring...")
        } else if wasOpen {
            // Positions were open, now all closed
            fmt.Println("\n✅ All positions closed!")
            fmt.Println("🔄 Restarting trading strategy...")

            wasOpen = false
            // Restart your strategy here
            // StartTradingStrategy(sugar)
        } else {
            fmt.Println("No positions - waiting for entry signal...")
        }
    }
}
```

---

### 10) Advanced position validator

```go
func ValidatePositionState(sugar *mt5.MT5Sugar, symbol string) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      POSITION STATE VALIDATOR         ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Fetch the symbol's positions once (HasOpenPosition would fetch the same list)
    positions, err := sugar.GetPositionsBySymbol(symbol)
    if err != nil {
        return fmt.Errorf("position check failed: %w", err)
    }

    fmt.Printf("Has open %s positions: %v\n\n", symbol, len(positions) > 0)

    if len(positions) == 0 {
        fmt.Println("✅ VALIDATION PASSED")
        fmt.Println("   No positions to validate")
        return nil
    }

    fmt.Printf("Found %d open positions:\n", len(positions))

    // Validate each position
    issues := 0
    for _, pos := range positions {
        fmt.Printf("\n#%d %s:\n", pos.Ticket, pos.Symbol)

        if pos.StopLoss == 0 {
            fmt.Println("  ❌ No Stop Loss")
            issues++
        } else {
            fmt.Println("  ✅ Stop Loss set")
        }

        if pos.TakeProfit == 0 {
            fmt.Println("  ⚠️  No Take Profit")
        } else {
            fmt.Println("  ✅ Take Profit set")
        }
    }

    if issues > 0 {
        return fmt.Errorf("%d positions have issues", issues)
    }

    fmt.Println("\n✅ VALIDATION PASSED")
    return nil
}

// Usage:
err := ValidatePositionState(sugar, "EURUSD")
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

---

## 🔗 Related Methods

**🍬 Other checks:**
* `CountOpenPositions()` - Total number of open positions (all symbols); `> 0` answers "any position at all?"
* `GetPositionsBySymbol()` - The positions behind this check

**🍬 Full position info:**

* `GetOpenPositions()` - Get all position details
* `sugar.GetService().GetOpenedTickets(ctx)` - Get just ticket numbers (positions and pending orders)

---

## ⚠️ Common Pitfalls

### 1) Using when you need details

```go
// ❌ WRONG - checking, then fetching separately
hasPos, _ := sugar.HasOpenPosition("EURUSD")
if hasPos {
    positions, _ := sugar.GetPositionsBySymbol("EURUSD") // Same list fetched twice!
}

// ✅ CORRECT - just get positions directly
positions, _ := sugar.GetPositionsBySymbol("EURUSD")
if len(positions) > 0 {
    // Work with positions
}
```

### 2) Not handling errors

```go
// ❌ WRONG - ignoring errors
hasPos, _ := sugar.HasOpenPosition("EURUSD")
if !hasPos {
    // Might be a false negative: on error the result is false!
}

// ✅ CORRECT - check errors
hasPos, err := sugar.HasOpenPosition("EURUSD")
if err != nil {
    fmt.Printf("Check failed: %v\n", err)
    return
}
```

### 3) Wrong symbol name

```go
// ❌ WRONG - the comparison is exact and case-sensitive
hasPos, _ := sugar.HasOpenPosition("eurusd") // false even with EURUSD positions

// ✅ CORRECT - use the broker's exact symbol name (including any suffix, e.g. "EURUSD.m")
hasPos, _ := sugar.HasOpenPosition("EURUSD")
```

---

## 💎 Pro Tips

1. **Yes/no per symbol** - Use when you only need a yes/no answer for one symbol

2. **Any symbol** - Use `CountOpenPositions() > 0` to check across all symbols

3. **Validation** - Good for end-of-day or startup checks

4. **State machines** - Use to track trading state per symbol

5. **Combine with other checks** - Connection + positions = system health

---

## 🔍 Performance Note

- It costs the same as `GetOpenPositions()` - the full position list is fetched, then filtered by symbol
- `CountOpenPositions()` returns only a number (`PositionsTotal`), but counts all symbols
- ❌ Don't call it and then fetch the positions anyway - call `GetPositionsBySymbol()` once

---

**See also:** [`CountOpenPositions.md`](CountOpenPositions.md), [`GetPositionsBySymbol.md`](../6.%20Position_Management/GetPositionsBySymbol.md), [`GetOpenPositions.md`](../6.%20Position_Management/GetOpenPositions.md)
