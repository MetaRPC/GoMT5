# ✅ Can Open Position (`CanOpenPosition`)

> **Sugar method:** Comprehensive validation check - verifies if you can open a position with specified volume BEFORE attempting to trade.

**API Information:**

* **Method:** `sugar.CanOpenPosition(symbol, volume)`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `IsSymbolAvailable()`, `GetSymbolInfo()`, `service.GetSymbolTick()`, `service.CalculateMargin()`, `GetFreeMargin()`
* **Timeout:** 5 seconds (for the tick and margin queries; `IsSymbolAvailable`, `GetSymbolInfo` and `GetFreeMargin` use their own 3 s / 5 s / 3 s timeouts)
* **Returns:** Three values (bool, string, error)

---

## 📋 Method Signature

```go
func (s *MT5Sugar) CanOpenPosition(symbol string, volume float64) (bool, string, error)
```

---

## 🔽 Input

| Parameter | Type | Description |
|-----------|------|-------------|
| `symbol` | `string` | Trading symbol (e.g., "EURUSD", "GBPUSD") |
| `volume` | `float64` | Desired lot size (e.g., 0.1, 1.0) |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `can` | `bool` | `true` if position can be opened, `false` if not |
| `reason` | `string` | Explanation if can't open (empty string if can open) |
| `error` | `error` | Error if validation check itself failed |

**Return patterns:**
```go
// Success - can open
(true, "", nil)

// Blocked - validation failed (not an error, just can't open)
(false, "insufficient margin: need 500.00, have 300.00", nil)
(false, "volume 0.15 not a multiple of step 0.01", nil)
(false, "symbol INVALID is not available", nil)

// Error - validation check failed
(false, "", fmt.Errorf("connection timeout"))
```

---

## 💬 Just the Essentials

* **What it is:** Pre-flight check before trading - validates symbol, volume and margin.
* **Why you need it:** **Call this before trading** to catch the most common rejection causes early.
* **Sanity check:** Checks 5 things: symbol availability, volume minimum, volume maximum, volume step, required margin vs free margin. It does NOT check trade mode, market hours, SL/TP or stop levels - the broker can still reject the order.

---

## 🎯 Validation Checks Performed

This method validates:

1. ✅ **Symbol availability** - Does the symbol exist and is it synchronized? (`IsSymbolAvailable`; trade mode is not checked)

2. ✅ **Volume minimum** - Is volume >= VolumeMin?

3. ✅ **Volume maximum** - Is volume <= VolumeMax?

4. ✅ **Volume step** - Is volume a multiple of VolumeStep? (tolerance 0.0001)

5. ✅ **Free margin** - Is the required margin <= free margin? The margin is calculated with `OrderCalcMargin` for a **BUY** at the current **Ask**, whatever direction you intend to trade (no safety buffer).

**If ANY check fails → returns (false, reason, nil)**

**If a query fails (symbol check, symbol info, tick, margin calculation, free margin) → returns (false, "", err)**

---

## 🔗 Usage Examples

### 1) Basic usage - validate before trading

```go
symbol := "EURUSD"
volume := 0.1

canOpen, reason, err := sugar.CanOpenPosition(symbol, volume)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
    return
}

if !canOpen {
    fmt.Printf("❌ Cannot open position: %s\n", reason)
    return
}

fmt.Println("✅ Validation passed - opening position...")
ticket, _ := sugar.BuyMarket(symbol, volume)
fmt.Printf("Position #%d opened\n", ticket)
```

---

### 2) Complete trading workflow with validation

```go
symbol := "EURUSD"
riskPercent := 2.0
stopLoss := 500.0    // points
takeProfit := 1000.0 // points

// Step 1: Calculate position size
lotSize, err := sugar.CalculatePositionSize(symbol, riskPercent, stopLoss)
if err != nil {
    fmt.Printf("Position size calculation failed: %v\n", err)
    return
}

// Step 2: Validate BEFORE trading
canOpen, reason, err := sugar.CanOpenPosition(symbol, lotSize)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
    return
}

if !canOpen {
    fmt.Printf("❌ Cannot open position: %s\n", reason)
    fmt.Printf("   Symbol:   %s\n", symbol)
    fmt.Printf("   Lot size: %.2f\n", lotSize)
    return
}

// Step 3: All checks passed - safe to trade
fmt.Println("✅ All validations passed")
ticket, err := sugar.BuyMarketWithPips(symbol, lotSize, stopLoss, takeProfit)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Position #%d opened successfully\n", ticket)
fmt.Printf("   Lot size: %.2f\n", lotSize)
fmt.Printf("   SL: %.0f points, TP: %.0f points\n", stopLoss, takeProfit)
```

---

### 3) Handle different failure reasons

```go
func OpenPositionSafely(sugar *mt5.MT5Sugar, symbol string, volume float64) {
    canOpen, reason, err := sugar.CanOpenPosition(symbol, volume)

    if err != nil {
        fmt.Printf("🔴 Validation check failed: %v\n", err)
        return
    }

    if !canOpen {
        // Parse the reason and provide helpful feedback
        if strings.Contains(reason, "margin") {
            fmt.Printf("💰 Insufficient margin: %s\n", reason)
            fmt.Println("   → Reduce position size or close other positions")
        } else if strings.Contains(reason, "minimum") {
            fmt.Printf("📉 Volume too small: %s\n", reason)
            info, _ := sugar.GetSymbolInfo(symbol)
            fmt.Printf("   → Minimum lot size: %.2f\n", info.VolumeMin)
        } else if strings.Contains(reason, "maximum") {
            fmt.Printf("📈 Volume too large: %s\n", reason)
            info, _ := sugar.GetSymbolInfo(symbol)
            fmt.Printf("   → Maximum lot size: %.2f\n", info.VolumeMax)
        } else if strings.Contains(reason, "step") {
            fmt.Printf("⚙️  Invalid volume step: %s\n", reason)
            info, _ := sugar.GetSymbolInfo(symbol)
            fmt.Printf("   → Volume must be multiple of %.2f\n", info.VolumeStep)
        } else if strings.Contains(reason, "available") {
            fmt.Printf("❌ Symbol unavailable: %s\n", reason)
            fmt.Println("   → Check symbol name (or wait until the symbol is synchronized)")
        } else {
            fmt.Printf("⚠️  Cannot open: %s\n", reason)
        }
        return
    }

    // Validation passed - open position
    fmt.Println("✅ Opening position...")
    ticket, _ := sugar.BuyMarket(symbol, volume)
    fmt.Printf("   Position #%d opened\n", ticket)
}

// Usage:
OpenPositionSafely(sugar, "EURUSD", 0.1)
```

---

### 4) Validate multiple positions before batch trading

```go
type TradeRequest struct {
    Symbol string
    Volume float64
}

func ValidateMultiplePositions(
    sugar *mt5.MT5Sugar,
    requests []TradeRequest,
) []TradeRequest {
    validRequests := []TradeRequest{}

    fmt.Println("Validating positions...")
    fmt.Println("─────────────────────────────────────────")

    for i, req := range requests {
        canOpen, reason, err := sugar.CanOpenPosition(req.Symbol, req.Volume)

        if err != nil {
            fmt.Printf("%d. %s %.2f lots - ⚠️  Error: %v\n",
                i+1, req.Symbol, req.Volume, err)
            continue
        }

        if !canOpen {
            fmt.Printf("%d. %s %.2f lots - ❌ %s\n",
                i+1, req.Symbol, req.Volume, reason)
            continue
        }

        fmt.Printf("%d. %s %.2f lots - ✅ Valid\n",
            i+1, req.Symbol, req.Volume)
        validRequests = append(validRequests, req)
    }

    fmt.Println("─────────────────────────────────────────")
    fmt.Printf("Valid: %d/%d positions\n", len(validRequests), len(requests))

    return validRequests
}

// Usage:
requests := []TradeRequest{
    {"EURUSD", 0.1},
    {"GBPUSD", 0.2},
    {"USDJPY", 0.15},
    {"INVALID", 0.1},
}

validRequests := ValidateMultiplePositions(sugar, requests)

// Open only valid positions
for _, req := range validRequests {
    ticket, _ := sugar.BuyMarket(req.Symbol, req.Volume)
    fmt.Printf("Opened %s: #%d\n", req.Symbol, ticket)
}
```

---

### 5) Pre-trade checklist display

```go
func ShowPreTradeChecklist(
    sugar *mt5.MT5Sugar,
    symbol string,
    volume float64,
) bool {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      PRE-TRADE VALIDATION CHECK       ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Symbol:     %s\n", symbol)
    fmt.Printf("Lot size:   %.2f\n\n", volume)

    // Check symbol availability
    available, _ := sugar.IsSymbolAvailable(symbol)
    if available {
        fmt.Println("✅ Symbol available")
    } else {
        fmt.Println("❌ Symbol not available")
        return false
    }

    // Get symbol info
    info, _ := sugar.GetSymbolInfo(symbol)

    // Check volume minimum
    if volume >= info.VolumeMin {
        fmt.Printf("✅ Volume >= minimum (%.2f)\n", info.VolumeMin)
    } else {
        fmt.Printf("❌ Volume < minimum (%.2f)\n", info.VolumeMin)
        return false
    }

    // Check volume maximum
    if volume <= info.VolumeMax {
        fmt.Printf("✅ Volume <= maximum (%.2f)\n", info.VolumeMax)
    } else {
        fmt.Printf("❌ Volume > maximum (%.2f)\n", info.VolumeMax)
        return false
    }

    // Check volume step
    steps := volume / info.VolumeStep
    if steps == float64(int(steps)) {
        fmt.Printf("✅ Volume is multiple of step (%.2f)\n", info.VolumeStep)
    } else {
        fmt.Printf("❌ Volume not multiple of step (%.2f)\n", info.VolumeStep)
        return false
    }

    // Check margin
    requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, volume)
    freeMargin, _ := sugar.GetFreeMargin()

    fmt.Printf("   Required margin: $%.2f\n", requiredMargin)
    fmt.Printf("   Free margin:     $%.2f\n", freeMargin)

    if requiredMargin <= freeMargin {
        fmt.Println("✅ Sufficient margin")
    } else {
        fmt.Println("❌ Insufficient margin")
        return false
    }

    fmt.Println("\n✅ ALL CHECKS PASSED - Ready to trade!")
    return true
}

// Usage:
if ShowPreTradeChecklist(sugar, "EURUSD", 0.1) {
    ticket, _ := sugar.BuyMarket("EURUSD", 0.1)
    fmt.Printf("\nPosition opened: #%d\n", ticket)
}
```

---

### 6) Retry with adjusted volume

```go
func OpenPositionWithAdjustment(
    sugar *mt5.MT5Sugar,
    symbol string,
    initialVolume float64,
) (uint64, error) {
    volume := initialVolume

    // Try with initial volume
    canOpen, reason, err := sugar.CanOpenPosition(symbol, volume)
    if err != nil {
        return 0, err
    }

    if canOpen {
        fmt.Printf("✅ Opening %.2f lots...\n", volume)
        return sugar.BuyMarket(symbol, volume)
    }

    // If insufficient margin, try with reduced volume
    if strings.Contains(reason, "margin") {
        fmt.Printf("⚠️  Insufficient margin for %.2f lots\n", volume)

        // Get max capacity
        maxLots, _ := sugar.GetMaxLotSize(symbol)
        if maxLots == 0 {
            return 0, fmt.Errorf("no margin available")
        }

        // Try with max capacity
        volume = maxLots
        fmt.Printf("   Trying with reduced volume: %.2f lots\n", volume)

        canOpen, reason, err = sugar.CanOpenPosition(symbol, volume)
        if err != nil {
            return 0, err
        }

        if canOpen {
            return sugar.BuyMarket(symbol, volume)
        }
    }

    return 0, fmt.Errorf("cannot open position: %s", reason)
}

// Usage:
ticket, err := OpenPositionWithAdjustment(sugar, "EURUSD", 1.0)
if err != nil {
    fmt.Printf("Failed: %v\n", err)
} else {
    fmt.Printf("Success: #%d\n", ticket)
}
```

---

### 7) Comprehensive validation wrapper

```go
type ValidationResult struct {
    CanOpen       bool
    Reason        string
    RequiredMargin float64
    FreeMargin     float64
    MarginRatio    float64
}

func ValidatePositionDetailed(
    sugar *mt5.MT5Sugar,
    symbol string,
    volume float64,
) (*ValidationResult, error) {
    result := &ValidationResult{}

    // Run validation
    canOpen, reason, err := sugar.CanOpenPosition(symbol, volume)
    if err != nil {
        return nil, err
    }

    result.CanOpen = canOpen
    result.Reason = reason

    // Get margin details
    requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, volume)
    freeMargin, _ := sugar.GetFreeMargin()

    result.RequiredMargin = requiredMargin
    result.FreeMargin = freeMargin

    if freeMargin > 0 {
        result.MarginRatio = (requiredMargin / freeMargin) * 100
    }

    return result, nil
}

func ShowValidationReport(result *ValidationResult, symbol string, volume float64) {
    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          POSITION VALIDATION REPORT                   ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")
    fmt.Printf("Symbol:            %s\n", symbol)
    fmt.Printf("Desired volume:    %.2f lots\n\n", volume)

    fmt.Printf("Required margin:   $%.2f\n", result.RequiredMargin)
    fmt.Printf("Available margin:  $%.2f\n", result.FreeMargin)
    fmt.Printf("Margin usage:      %.1f%%\n\n", result.MarginRatio)

    if result.CanOpen {
        fmt.Println("✅ VALIDATION PASSED")
        fmt.Println("   Position can be opened")

        if result.MarginRatio > 70 {
            fmt.Println("\n⚠️  Warning: High margin usage (>70%)")
            fmt.Println("   Consider reducing position size")
        }
    } else {
        fmt.Println("❌ VALIDATION FAILED")
        fmt.Printf("   Reason: %s\n", result.Reason)
    }
}

// Usage:
result, err := ValidatePositionDetailed(sugar, "EURUSD", 0.5)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
    return
}

ShowValidationReport(result, "EURUSD", 0.5)

if result.CanOpen {
    ticket, _ := sugar.BuyMarket("EURUSD", 0.5)
    fmt.Printf("\nPosition opened: #%d\n", ticket)
}
```

---

### 8) Margin safety validator

```go
func ValidateWithMarginSafety(
    sugar *mt5.MT5Sugar,
    symbol string,
    volume float64,
    maxMarginPercent float64,
) (bool, string) {
    // First, standard validation
    canOpen, reason, err := sugar.CanOpenPosition(symbol, volume)
    if err != nil {
        return false, fmt.Sprintf("validation error: %v", err)
    }

    if !canOpen {
        return false, reason
    }

    // Additional check: margin usage threshold
    requiredMargin, _ := sugar.CalculateRequiredMargin(symbol, volume)
    freeMargin, _ := sugar.GetFreeMargin()

    marginUsage := (requiredMargin / freeMargin) * 100

    if marginUsage > maxMarginPercent {
        return false, fmt.Sprintf(
            "margin usage %.1f%% exceeds limit %.1f%%",
            marginUsage, maxMarginPercent,
        )
    }

    return true, ""
}

// Usage with 50% margin limit:
canOpen, reason := ValidateWithMarginSafety(sugar, "EURUSD", 0.5, 50.0)
if !canOpen {
    fmt.Printf("❌ Cannot open: %s\n", reason)
} else {
    fmt.Println("✅ Validation passed with margin safety check")
    ticket, _ := sugar.BuyMarket("EURUSD", 0.5)
    fmt.Printf("Position opened: #%d\n", ticket)
}
```

---

### 9) Batch validation with statistics

```go
type TradeRequest struct {
    Symbol string
    Volume float64
}

func BatchValidateAndReport(sugar *mt5.MT5Sugar, requests []TradeRequest) {
    passCount := 0
    failCount := 0
    errorCount := 0

    failureReasons := make(map[string]int)

    fmt.Println("╔═══════════════════════════════════════════════════════╗")
    fmt.Println("║          BATCH VALIDATION REPORT                      ║")
    fmt.Println("╚═══════════════════════════════════════════════════════╝")
    fmt.Printf("Total requests: %d\n\n", len(requests))

    for i, req := range requests {
        canOpen, reason, err := sugar.CanOpenPosition(req.Symbol, req.Volume)

        status := ""
        if err != nil {
            status = fmt.Sprintf("⚠️  Error: %v", err)
            errorCount++
        } else if canOpen {
            status = "✅ Pass"
            passCount++
        } else {
            status = fmt.Sprintf("❌ %s", reason)
            failCount++

            // Track failure reason
            failureReasons[reason]++
        }

        fmt.Printf("%d. %s %.2f lots - %s\n",
            i+1, req.Symbol, req.Volume, status)
    }

    fmt.Println("\n─────────────────────────────────────────────────────────")
    fmt.Printf("✅ Passed:  %d (%.1f%%)\n",
        passCount, float64(passCount)/float64(len(requests))*100)
    fmt.Printf("❌ Failed:  %d (%.1f%%)\n",
        failCount, float64(failCount)/float64(len(requests))*100)
    fmt.Printf("⚠️  Errors: %d\n", errorCount)

    if len(failureReasons) > 0 {
        fmt.Println("\nFailure breakdown:")
        for reason, count := range failureReasons {
            fmt.Printf("  • %s: %d times\n", reason, count)
        }
    }
}

// Usage:
requests := []TradeRequest{
    {"EURUSD", 0.1},
    {"GBPUSD", 0.2},
    {"USDJPY", 100.0}, // Too large
    {"INVALID", 0.1},  // Invalid symbol
    {"XAUUSD", 0.01},
}

BatchValidateAndReport(sugar, requests)
```

---

### 10) Advanced position validator

```go
type PositionValidator struct {
    sugar                *mt5.MT5Sugar
    maxMarginUsagePercent float64
    maxPositions         int
}

func NewPositionValidator(
    sugar *mt5.MT5Sugar,
    maxMarginPercent float64,
    maxPositions int,
) *PositionValidator {
    return &PositionValidator{
        sugar:                sugar,
        maxMarginUsagePercent: maxMarginPercent,
        maxPositions:         maxPositions,
    }
}

func (pv *PositionValidator) Validate(
    symbol string,
    volume float64,
) (bool, string, error) {
    // Check 1: Standard validation
    canOpen, reason, err := pv.sugar.CanOpenPosition(symbol, volume)
    if err != nil {
        return false, "", err
    }

    if !canOpen {
        return false, reason, nil
    }

    // Check 2: Maximum positions limit
    posCount, _ := pv.sugar.CountOpenPositions()
    if posCount >= pv.maxPositions {
        return false, fmt.Sprintf(
            "maximum positions reached (%d/%d)",
            posCount, pv.maxPositions,
        ), nil
    }

    // Check 3: Margin usage threshold
    requiredMargin, _ := pv.sugar.CalculateRequiredMargin(symbol, volume)
    freeMargin, _ := pv.sugar.GetFreeMargin()

    marginUsage := (requiredMargin / freeMargin) * 100

    if marginUsage > pv.maxMarginUsagePercent {
        return false, fmt.Sprintf(
            "margin usage %.1f%% exceeds limit %.1f%%",
            marginUsage, pv.maxMarginUsagePercent,
        ), nil
    }

    // Check 4: Symbol not already in portfolio (optional rule)
    openPositions, _ := pv.sugar.GetOpenPositions()
    for _, pos := range openPositions {
        if pos.Symbol == symbol {
            return false, fmt.Sprintf(
                "position already exists for %s (ticket #%d)",
                symbol, pos.Ticket,
            ), nil
        }
    }

    return true, "", nil
}

func (pv *PositionValidator) ValidateAndExecute(
    symbol string,
    volume float64,
    direction string,
) (uint64, error) {
    canOpen, reason, err := pv.Validate(symbol, volume)
    if err != nil {
        return 0, fmt.Errorf("validation error: %w", err)
    }

    if !canOpen {
        return 0, fmt.Errorf("validation failed: %s", reason)
    }

    // Execute trade
    if direction == "BUY" {
        return pv.sugar.BuyMarket(symbol, volume)
    } else {
        return pv.sugar.SellMarket(symbol, volume)
    }
}

func (pv *PositionValidator) ShowRules() {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║      POSITION VALIDATOR RULES         ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Max margin usage:  %.1f%%\n", pv.maxMarginUsagePercent)
    fmt.Printf("Max positions:     %d\n", pv.maxPositions)
    fmt.Println("Symbol uniqueness: Enforced")
    fmt.Println("Broker limits:     Enforced")
}

// Usage:
validator := NewPositionValidator(sugar, 50.0, 5)
validator.ShowRules()

// Validate and execute
ticket, err := validator.ValidateAndExecute("EURUSD", 0.1, "BUY")
if err != nil {
    fmt.Printf("Failed: %v\n", err)
} else {
    fmt.Printf("✅ Position opened: #%d\n", ticket)
}
```

---

## 🔗 Related Methods

**📦 Methods used internally:**

* `IsSymbolAvailable()` - Check symbol exists and is synchronized
* `GetSymbolInfo()` - Get volume limits
* `service.GetSymbolTick()` - Get current Ask
* `service.CalculateMargin()` - Calculate required margin (BUY at Ask)
* `GetFreeMargin()` - Get available margin

**🍬 Complementary sugar methods:**

* `CalculatePositionSize()` - Calculate risk-based lot size ⭐
* `GetMaxLotSize()` - Get maximum tradeable volume
* `CalculateRequiredMargin()` - Get margin needed
* `BuyMarket() / SellMarket()` - Execute trades after validation

**Recommended workflow:**
```go
// 1. Calculate lot size
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500) // 500-point SL

// 2. ALWAYS validate before trading
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", lotSize)
if !canOpen {
    fmt.Println("Cannot trade:", reason)
    return
}

// 3. Trade
ticket, _ := sugar.BuyMarket("EURUSD", lotSize)
```

---

## ⚠️ Common Pitfalls

### 1) Not calling before trading

```go
// ❌ WRONG - trading without validation
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)
// Order might get rejected!

// ✅ CORRECT - validate first
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 0.1)
if !canOpen {
    fmt.Println("Cannot trade:", reason)
    return
}
ticket, _ := sugar.BuyMarket("EURUSD", 0.1)
```

### 2) Ignoring the reason when false

```go
// ❌ WRONG - not checking why it failed
canOpen, _, _ := sugar.CanOpenPosition("EURUSD", 0.1)
if !canOpen {
    fmt.Println("Can't trade") // Not helpful!
}

// ✅ CORRECT - use the reason
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 0.1)
if !canOpen {
    fmt.Printf("Can't trade: %s\n", reason) // Informative!
}
```

### 3) Confusing error vs blocked

```go
// ❌ WRONG - treating blocked as error
canOpen, reason, err := sugar.CanOpenPosition("EURUSD", 0.1)
if err != nil || !canOpen {
    // These are different things!
}

// ✅ CORRECT - handle separately
canOpen, reason, err := sugar.CanOpenPosition("EURUSD", 0.1)
if err != nil {
    // Actual error - validation check failed
    fmt.Printf("Error: %v\n", err)
    return
}
if !canOpen {
    // Not an error - just can't open (with reason why)
    fmt.Printf("Blocked: %s\n", reason)
    return
}
```

### 4) Not checking error return

```go
// ❌ WRONG - ignoring error
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 0.1)

// ✅ CORRECT - check error
canOpen, reason, err := sugar.CanOpenPosition("EURUSD", 0.1)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
    return
}
```

### 5) Validating but not using result

```go
// ❌ WRONG - validating but ignoring result
sugar.CanOpenPosition("EURUSD", 0.1)
sugar.BuyMarket("EURUSD", 0.1) // Might still fail!

// ✅ CORRECT - use validation result
canOpen, reason, _ := sugar.CanOpenPosition("EURUSD", 0.1)
if !canOpen {
    fmt.Println("Cannot trade:", reason)
    return
}
sugar.BuyMarket("EURUSD", 0.1)
```

---

## 💎 Pro Tips

1. **ALWAYS call this** - Make it mandatory before every trade

2. **Check all three returns** - (bool, string, error) all matter

3. **Use the reason** - Provides specific feedback on what's wrong

4. **False ≠ Error** - `false` means "can't open" (not a system error)

5. **Validate calculated sizes** - Even CalculatePositionSize() results should be validated

6. **No trade execution** - Read-only, but it makes several server queries (symbol exist/sync, symbol params, stop level, tick, margin, free margin)

7. **Multi-check** - Validates 5 different things in one call

---

## 📊 Return Value Meanings

```
Three return values:

(true, "", nil)
→ ✅ All 5 checks passed (the broker can still reject for other reasons)

(false, "reason", nil)
→ ❌ Can't open (not an error, just blocked)
→ "reason" explains what's wrong

(false, "", error)
→ 🔴 Validation check itself failed
→ System error, not a trading validation issue

Common reasons when false:
- "insufficient margin: need X, have Y"
- "volume X below minimum Y"
- "volume X exceeds maximum Y"
- "volume X not a multiple of step Y"
- "symbol X is not available"
```

---

**See also:** [`CalculatePositionSize.md`](CalculatePositionSize.md), [`GetMaxLotSize.md`](GetMaxLotSize.md), [`CalculateRequiredMargin.md`](CalculateRequiredMargin.md)
