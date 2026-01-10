# ✔️ Has Open Position (`HasOpenPosition`)

> **Sugar method:** Quickly checks if ANY open position exists (across all symbols).

**API Information:**

* **Method:** `sugar.HasOpenPosition()`
* **Timeout:** 3 seconds
* **Returns:** Boolean (true if any position open)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) HasOpenPosition() (bool, error)
```

---

## 🔽 Input / ⬆️ Output

| Input | Type | Description |
|-------|------|-------------|
| None | - | No parameters required |

| Output | Type | Description |
|--------|------|-------------|
| `bool` | boolean | `true` if any position open, `false` if none |
| `error` | `error` | Error if check failed |

---

## 💬 Just the Essentials

* **What it is:** Fast boolean check - "Do I have any open positions?"
* **Why you need it:** Quick validation before opening new positions or ending session.
* **Sanity check:** More efficient than `len(GetOpenPositions()) > 0`.

---

## 🎯 When to Use

✅ **Quick check** - Fast boolean answer without fetching full position data

✅ **Before trading** - Verify no positions before starting new strategy

✅ **End of day** - Confirm everything closed before shutdown

✅ **Validation** - Check state without overhead of full position retrieval


---

## 🔗 Usage Examples

### 1) Basic usage - quick check

```go
hasPos, err := sugar.HasOpenPosition()
if err != nil {
    fmt.Printf("Check failed: %v\n", err)
    return
}

if hasPos {
    fmt.Println("✅ You have open positions")
} else {
    fmt.Println("❌ No open positions")
}
```

---

### 2) Verify clean slate before trading

```go
hasPos, _ := sugar.HasOpenPosition()

if hasPos {
    fmt.Println("⚠️  WARNING: Positions already open!")
    fmt.Println("   Close existing positions before starting strategy")
    return
}

fmt.Println("✅ Clean slate - ready to trade")
// Start trading strategy
```

---

### 3) End of day validation

```go
func ValidateEndOfDay(sugar *mt5.MT5Sugar) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      END OF DAY VALIDATION            ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    hasPos, err := sugar.HasOpenPosition()
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    if hasPos {
        fmt.Println("❌ ALERT: Positions still open!")
        fmt.Println("   Action required: Close all positions")
        return fmt.Errorf("positions still open")
    }

    fmt.Println("✅ All clear - no open positions")
    fmt.Println("   Safe to shutdown")
    return nil
}

// Usage: Run before weekend shutdown
err := ValidateEndOfDay(sugar)
```

---

### 4) Polling for position close

```go
// Wait for all positions to close
fmt.Println("Waiting for all positions to close...")

ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

timeout := time.After(5 * time.Minute)

for {
    select {
    case <-timeout:
        fmt.Println("⏰ Timeout - positions still open")
        return

    case <-ticker.C:
        hasPos, _ := sugar.HasOpenPosition()

        if !hasPos {
            fmt.Println("✅ All positions closed!")
            return
        }

        fmt.Println("Positions still open - waiting...")
    }
}
```

---

### 5) Before opening first position

```go
symbol := "EURUSD"
volume := 0.1

// Check if already trading
hasPos, _ := sugar.HasOpenPosition()

if hasPos {
    fmt.Println("⚠️  Already have active positions")
    fmt.Println("   Strategy: Only one position at a time allowed")
    return
}

fmt.Println("✅ No active positions - opening first trade")
ticket, _ := sugar.BuyMarket(symbol, volume)
fmt.Printf("Position #%d opened\n", ticket)
```

---

### 6) Monitoring for exit signal

```go
// Close all when specific condition met
ticker := time.NewTicker(10 * time.Second)
defer ticker.Stop()

for range ticker.C {
    hasPos, _ := sugar.HasOpenPosition()

    if !hasPos {
        fmt.Println("No positions - waiting for entry signal...")
        continue
    }

    // Check exit condition (example: time-based)
    currentHour := time.Now().Hour()

    if currentHour >= 17 { // Close at 5 PM
        fmt.Println("🕔 17:00 - Closing all positions")
        sugar.CloseAllPositions()
        return
    }

    fmt.Println("Positions active - monitoring...")
}
```

---

### 7) System health check

```go
func SystemHealthCheck(sugar *mt5.MT5Sugar) {
    fmt.Println("System Health Check:")
    fmt.Println("─────────────────────────────")

    // Check connection
    connected := sugar.IsConnected()
    if connected {
        fmt.Println("✅ Connection: OK")
    } else {
        fmt.Println("❌ Connection: FAILED")
    }

    // Check positions
    hasPos, err := sugar.HasOpenPosition()
    if err != nil {
        fmt.Println("❌ Position check: FAILED")
    } else if hasPos {
        fmt.Println("⚠️  Positions: OPEN")
    } else {
        fmt.Println("✅ Positions: NONE")
    }

    // Check balance
    balance, _ := sugar.GetBalance()
    fmt.Printf("💰 Balance: $%.2f\n", balance)
}

// Usage:
SystemHealthCheck(sugar)
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

func GetStrategyState(sugar *mt5.MT5Sugar) StrategyState {
    hasPos, err := sugar.HasOpenPosition()
    if err != nil {
        return StateIdle
    }

    if hasPos {
        return StateTrading
    }

    return StateIdle
}

// Usage:
state := GetStrategyState(sugar)

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
func MonitorAndRestart(sugar *mt5.MT5Sugar) {
    fmt.Println("Monitoring positions - will restart when all closed")

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    wasOpen := false

    for range ticker.C {
        hasPos, _ := sugar.HasOpenPosition()

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
func ValidatePositionState(sugar *mt5.MT5Sugar) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      POSITION STATE VALIDATOR         ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Quick check first
    hasPos, err := sugar.HasOpenPosition()
    if err != nil {
        return fmt.Errorf("position check failed: %w", err)
    }

    fmt.Printf("Has open positions: %v\n\n", hasPos)

    if !hasPos {
        fmt.Println("✅ VALIDATION PASSED")
        fmt.Println("   No positions to validate")
        return nil
    }

    // If positions exist, get details
    positions, _ := sugar.GetOpenPositions()

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
err := ValidatePositionState(sugar)
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

---

## 🔗 Related Methods

**🍬 More specific checks:**
* `HasOpenPositionBySymbol()` - Check for specific symbol
* `CountOpenPositions()` - Get exact count

**🍬 Full position info:**

* `GetOpenPositions()` - Get all position details (slower)
* `GetPositionTickets()` - Get just ticket numbers

---

## ⚠️ Common Pitfalls

### 1) Using when you need details

```go
// ❌ WRONG - checking, then fetching separately
hasPos, _ := sugar.HasOpenPosition()
if hasPos {
    positions, _ := sugar.GetOpenPositions() // Duplicate API call!
}

// ✅ CORRECT - just get positions directly
positions, _ := sugar.GetOpenPositions()
if len(positions) > 0 {
    // Work with positions
}
```

### 2) Not handling errors

```go
// ❌ WRONG - ignoring errors
hasPos, _ := sugar.HasOpenPosition()
if hasPos {
    // Might be false negative due to error!
}

// ✅ CORRECT - check errors
hasPos, err := sugar.HasOpenPosition()
if err != nil {
    fmt.Printf("Check failed: %v\n", err)
    return
}
```

### 3) Polling too frequently

```go
// ❌ WRONG - checking every second (wasteful)
for {
    hasPos, _ := sugar.HasOpenPosition()
    time.Sleep(1 * time.Second)
}

// ✅ CORRECT - reasonable interval
ticker := time.NewTicker(5 * time.Second)
for range ticker.C {
    hasPos, _ := sugar.HasOpenPosition()
}
```

---

## 💎 Pro Tips

1. **Fast check** - Use when you only need yes/no answer

2. **Before GetOpenPositions()** - Check first to avoid unnecessary API call

3. **Validation** - Perfect for end-of-day or startup checks

4. **State machines** - Use to track trading state

5. **Combine with other checks** - Connection + positions = system health

---

## 🔍 Performance Note

This method is optimized for speed:

- ✅ Faster than `GetOpenPositions()` (doesn't fetch full position data)
- ✅ More efficient than counting positions
- ✅ Perfect for frequent polling/monitoring
- ❌ Don't use if you need position details anyway

---

**See also:** [`CountOpenPositions.md`](CountOpenPositions.md), [`GetOpenPositions.md`](../6.%20Position_Management/GetOpenPositions.md)
