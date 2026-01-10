# ⏳ Wait For Valid Price (`WaitForPrice`)

> **Sugar method:** Waits until symbol has valid price data (useful after connection or market open).

**API Information:**

* **Method:** `sugar.WaitForPrice(symbol, timeout)`
* **Timeout:** Custom (you specify)
* **Returns:** `*PriceInfo` when valid price received

---

## 📋 Method Signature

```go
func (s *MT5Sugar) WaitForPrice(symbol string, timeout time.Duration) (*PriceInfo, error)
```

---

## 💬 Just the Essentials

* **What it is:** Blocks until symbol has valid price (BID > 0 and ASK > 0) or timeout expires.
* **Why you need it:** After connection or during market open, prices might not be available yet.
* **Sanity check:** Returns first valid PriceInfo or error on timeout.

---

## 🎯 When to Use

✅ **After connection** - Wait for first tick after connecting

✅ **Market open** - Wait for market to start publishing prices

✅ **Symbol switch** - Ensure new symbol is ready

✅ **Weekend/gap handling** - Wait for market to resume

---

## 🔗 Usage Examples

### 1) Basic usage - wait after connection

```go
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")
sugar.QuickConnect("FxPro-MT5 Demo")

// Wait for price to be available
symbol := "EURUSD"
priceInfo, err := sugar.WaitForPrice(symbol, 10*time.Second)
if err != nil {
    fmt.Printf("Timeout waiting for price: %v\n", err)
    return
}

fmt.Printf("✅ %s price ready!\n", symbol)
fmt.Printf("   BID: %.5f\n", priceInfo.Bid)
fmt.Printf("   ASK: %.5f\n", priceInfo.Ask)
```

---

### 2) Robust connection sequence

```go
func ConnectAndWaitForPrice(login uint64, password, server, cluster, symbol string) error {
    // Step 1: Create Sugar
    sugar, err := mt5.NewMT5Sugar(login, password, server)
    if err != nil {
        return fmt.Errorf("failed to create Sugar: %w", err)
    }

    // Step 2: Connect
    fmt.Println("🔌 Connecting to MT5...")
    err = sugar.QuickConnect(cluster)
    if err != nil {
        return fmt.Errorf("connection failed: %w", err)
    }

    // Step 3: Wait for price
    fmt.Printf("⏳ Waiting for %s price...\n", symbol)
    priceInfo, err := sugar.WaitForPrice(symbol, 15*time.Second)
    if err != nil {
        return fmt.Errorf("price timeout: %w", err)
    }

    fmt.Printf("✅ Ready to trade!\n")
    fmt.Printf("   %s BID: %.5f ASK: %.5f\n",
        symbol, priceInfo.Bid, priceInfo.Ask)

    return nil
}

// Usage:
err := ConnectAndWaitForPrice(
    591129415,
    "password",
    "mt5.mrpc.pro:443",
    "FxPro-MT5 Demo",
    "EURUSD")
```

---

### 3) Wait for multiple symbols

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

fmt.Println("Waiting for prices on all symbols...")

for _, symbol := range symbols {
    priceInfo, err := sugar.WaitForPrice(symbol, 10*time.Second)
    if err != nil {
        fmt.Printf("❌ %s: Timeout - %v\n", symbol, err)
        continue
    }

    fmt.Printf("✅ %s: BID=%.5f ASK=%.5f\n",
        symbol, priceInfo.Bid, priceInfo.Ask)
}
```

---

### 4) Market open detector

```go
func WaitForMarketOpen(sugar *mt5.MT5Sugar, symbol string) {
    fmt.Printf("⏰ Waiting for %s market to open...\n", symbol)

    startTime := time.Now()

    priceInfo, err := sugar.WaitForPrice(symbol, 2*time.Hour)
    if err != nil {
        fmt.Printf("❌ Market didn't open within 2 hours\n")
        return
    }

    elapsed := time.Since(startTime)

    fmt.Printf("✅ Market OPEN!\n")
    fmt.Printf("   Waited: %v\n", elapsed.Round(time.Second))
    fmt.Printf("   First price: BID=%.5f ASK=%.5f\n",
        priceInfo.Bid, priceInfo.Ask)
    fmt.Printf("   Server time: %s\n", priceInfo.Time.Format("15:04:05"))
}

// Usage on Monday morning
WaitForMarketOpen(sugar, "EURUSD")
```

---

### 5) Retry with increasing timeout

```go
func WaitForPriceWithRetry(sugar *mt5.MT5Sugar, symbol string, maxAttempts int) (*PriceInfo, error) {
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        timeout := time.Duration(attempt) * 5 * time.Second

        fmt.Printf("Attempt %d/%d: Waiting %v for %s price...\n",
            attempt, maxAttempts, timeout, symbol)

        priceInfo, err := sugar.WaitForPrice(symbol, timeout)
        if err == nil {
            fmt.Printf("✅ Got price on attempt %d\n", attempt)
            return priceInfo, nil
        }

        fmt.Printf("❌ Attempt %d failed: %v\n", attempt, err)

        if attempt < maxAttempts {
            fmt.Println("   Retrying...")
            time.Sleep(2 * time.Second)
        }
    }

    return nil, fmt.Errorf("all %d attempts failed", maxAttempts)
}

// Usage:
priceInfo, err := WaitForPriceWithRetry(sugar, "EURUSD", 3)
```

---

### 6) Weekend gap handler

```go
func HandleWeekendGap(sugar *mt5.MT5Sugar, symbol string) {
    fmt.Printf("📅 Weekend ended - waiting for %s to resume...\n", symbol)

    // Long timeout for weekend → Monday transition
    priceInfo, err := sugar.WaitForPrice(symbol, 30*time.Minute)
    if err != nil {
        fmt.Printf("❌ Market didn't resume: %v\n", err)
        return
    }

    // Check for gap
    // You'd need to compare with Friday's close price
    fmt.Printf("✅ Market resumed!\n")
    fmt.Printf("   Monday open: BID=%.5f\n", priceInfo.Bid)
    fmt.Printf("   Server time: %s\n", priceInfo.Time.Format("Mon 15:04:05"))

    // Could add gap detection logic here
}
```

---

### 7) Trading bot startup sequence

```go
func StartTradingBot(sugar *mt5.MT5Sugar, symbols []string) error {
    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║       TRADING BOT STARTUP             ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    // Check connection
    fmt.Println("\n1. Checking connection...")
    if !sugar.IsConnected() {
        return fmt.Errorf("not connected to MT5")
    }
    fmt.Println("   ✅ Connected")

    // Wait for prices on all symbols
    fmt.Println("\n2. Waiting for price feeds...")
    for i, symbol := range symbols {
        fmt.Printf("   [%d/%d] %s... ", i+1, len(symbols), symbol)

        priceInfo, err := sugar.WaitForPrice(symbol, 15*time.Second)
        if err != nil {
            fmt.Printf("❌ TIMEOUT\n")
            return fmt.Errorf("%s price not available", symbol)
        }

        fmt.Printf("✅ %.5f\n", priceInfo.Bid)
    }

    fmt.Println("\n✅ BOT READY TO TRADE")
    return nil
}

// Usage:
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
err := StartTradingBot(sugar, symbols)
if err != nil {
    fmt.Printf("Startup failed: %v\n", err)
}
```

---

### 8) Price feed validator

```go
func ValidatePriceFeed(sugar *mt5.MT5Sugar, symbol string) bool {
    fmt.Printf("Validating %s price feed...\n", symbol)

    // Try to get price within reasonable time
    priceInfo, err := sugar.WaitForPrice(symbol, 5*time.Second)
    if err != nil {
        fmt.Printf("❌ Price feed INVALID: %v\n", err)
        return false
    }

    // Validate price values
    if priceInfo.Bid <= 0 || priceInfo.Ask <= 0 {
        fmt.Printf("❌ Price feed INVALID: Invalid price values\n")
        return false
    }

    if priceInfo.Ask <= priceInfo.Bid {
        fmt.Printf("❌ Price feed INVALID: ASK <= BID\n")
        return false
    }

    // Check timestamp freshness
    age := time.Since(priceInfo.Time)
    if age > 10*time.Second {
        fmt.Printf("❌ Price feed STALE: %v old\n", age)
        return false
    }

    fmt.Printf("✅ Price feed VALID\n")
    fmt.Printf("   BID: %.5f, ASK: %.5f\n", priceInfo.Bid, priceInfo.Ask)
    fmt.Printf("   Age: %.1f seconds\n", age.Seconds())
    return true
}
```

---

### 9) Emergency reconnection

```go
func EmergencyReconnect(sugar *mt5.MT5Sugar, cluster, symbol string) error {
    fmt.Println("🚨 Emergency reconnection initiated...")

    // Reconnect
    err := sugar.QuickConnect(cluster)
    if err != nil {
        return fmt.Errorf("reconnect failed: %w", err)
    }

    fmt.Println("✅ Reconnected - waiting for price feed...")

    // Wait for price to confirm connection is working
    _, err = sugar.WaitForPrice(symbol, 30*time.Second)
    if err != nil {
        return fmt.Errorf("price feed not restored: %w", err)
    }

    fmt.Println("✅ Connection fully restored!")
    return nil
}

// Usage during connection loss
if !sugar.IsConnected() {
    err := EmergencyReconnect(sugar, "FxPro-MT5 Demo", "EURUSD")
    if err != nil {
        fmt.Printf("Recovery failed: %v\n", err)
    }
}
```

---

### 10) Symbol availability checker

```go
func CheckSymbolAvailability(sugar *mt5.MT5Sugar, symbols []string) {
    fmt.Println("Checking symbol availability:")
    fmt.Println("─────────────────────────────────────")

    available := []string{}
    unavailable := []string{}

    for _, symbol := range symbols {
        fmt.Printf("%-10s ... ", symbol)

        priceInfo, err := sugar.WaitForPrice(symbol, 3*time.Second)
        if err != nil {
            fmt.Println("❌ NOT AVAILABLE")
            unavailable = append(unavailable, symbol)
        } else {
            fmt.Printf("✅ BID=%.5f\n", priceInfo.Bid)
            available = append(available, symbol)
        }
    }

    fmt.Println("\n═══════════════════════════════════════")
    fmt.Printf("Available:   %d symbols\n", len(available))
    fmt.Printf("Unavailable: %d symbols\n", len(unavailable))

    if len(unavailable) > 0 {
        fmt.Println("\nUnavailable symbols:")
        for _, s := range unavailable {
            fmt.Printf("  - %s\n", s)
        }
    }
}

// Usage:
symbols := []string{"EURUSD", "GBPUSD", "XAUUSD", "INVALID", "BTCUSD"}
CheckSymbolAvailability(sugar, symbols)
```

---

## 🔗 Related Methods

* `GetPriceInfo()` - Get price immediately (no waiting)
* `IsSymbolAvailable()` - Check if symbol exists
* `QuickConnect()` - Connect to MT5

---

## ⚠️ Common Pitfalls

### 1) Timeout too short

```go
// ❌ WRONG - might timeout during slow connection
priceInfo, _ := sugar.WaitForPrice("EURUSD", 1*time.Second)

// ✅ CORRECT - reasonable timeout
priceInfo, _ := sugar.WaitForPrice("EURUSD", 10*time.Second)
```

### 2) Not checking for error

```go
// ❌ WRONG - ignoring timeout error
priceInfo, _ := sugar.WaitForPrice("EURUSD", 5*time.Second)
fmt.Printf("Price: %.5f\n", priceInfo.Bid) // Might panic if timeout!

// ✅ CORRECT - handle timeout
priceInfo, err := sugar.WaitForPrice("EURUSD", 5*time.Second)
if err != nil {
    fmt.Printf("Timeout: %v\n", err)
    return
}
```

---

## 💎 Pro Tips

1. **Use after connection** - Always wait for price before trading
2. **Reasonable timeout** - 10-15 seconds is usually enough
3. **Multiple symbols** - Check all symbols before starting bot
4. **Weekend handling** - Use longer timeout (30min) on Monday
5. **Fallback plan** - Have reconnection logic if timeout occurs

---

**See also:** [`GetPriceInfo.md`](GetPriceInfo.md), [`QuickConnect.md`](../1.%20Connection/QuickConnect.md)
