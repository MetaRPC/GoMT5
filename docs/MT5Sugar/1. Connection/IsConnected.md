# ✅ Check Connection Status (`IsConnected`)

> **Sugar method:** Verifies if you're currently connected to MT5 Terminal.

**API Information:**

* **Method:** `sugar.IsConnected()`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `account.CheckConnect()` - reads `HealthCheck.IsAlive` from the reply
* **Timeout:** 3 seconds

---

## 📋 Method Signature

```go
func (s *MT5Sugar) IsConnected() bool
```

---

## 🔽 Input

*No parameters*

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `connected` | `bool` | `true` if connected and healthy, `false` otherwise |

---

## 💬 Just the Essentials

* **What it is:** Quick check if MT5 connection is alive and working.
* **Why you need it:** Verify connection before trading operations, detect disconnects.
* **Sanity check:** After `QuickConnect()` succeeds, this should return `true`.

---

## 🎯 Purpose

Use it to monitor connection health:

* **Pre-flight check** - Verify connection before trading
* **Reconnect logic** - Detect when connection is lost
* **Health monitoring** - Periodic connection checks
* **Error recovery** - Know when to retry operations

---

## 🧩 Notes & Tips

* **Fast check** - 3-second timeout
* **True = ready** - Safe to call trading methods
* **False = disconnected** - Need to reconnect
* **Blocking call** - One `CheckConnect` round trip, waits at most 3 seconds; any error or timeout gives `false`
* **Use before trading** - Good practice to verify connection
* **Periodic checks** - Monitor in long-running bots

---

## 🔧 Under the Hood

```go
func (s *MT5Sugar) IsConnected() bool {
    ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
    defer cancel()

    data, err := s.GetAccount().CheckConnect(ctx, &pb.CheckConnectRequest{})
    if err != nil {
        return false
    }

    return data.HealthCheck != nil && data.HealthCheck.IsAlive
}
```

**What it does:**

* ✅ **Asks the server** - `CheckConnect` returns the terminal's health check, `true` only when `IsAlive`
* ✅ **Quick timeout** - 3 seconds max
* ✅ **Clean result** - True/false, easy to use

---

## 📊 Low-Level Alternative

**WITHOUT sugar:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

data, err := account.CheckConnect(ctx, &pb.CheckConnectRequest{})
connected := err == nil && data.HealthCheck != nil && data.HealthCheck.IsAlive
```

(`account.IsConnected()` on `MT5Account` is a different, local check: it only tells whether the gRPC connection
and a terminal id exist, without asking the server.)

**WITH sugar:**
```go
connected := sugar.IsConnected()
```

**Benefits:**

* ✅ **One line vs three**
* ✅ **Clear intent**
* ✅ **Built-in timeout**

---

## 🔗 Usage Examples

### 1) Basic connection check

```go
connected := sugar.IsConnected()

if connected {
    fmt.Println("✅ Connected to MT5")
} else {
    fmt.Println("❌ Not connected")
}
```

---

### 2) Verify before trading

```go
// Check connection before placing order
connected := sugar.IsConnected()
if !connected {
    fmt.Println("❌ Not connected - cannot place order")
    return
}

// Connection OK - safe to trade
ticket, err := sugar.BuyMarket("EURUSD", 0.1)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
    return
}

fmt.Printf("✅ Order #%d placed\n", ticket)
```

---

### 3) Automatic reconnection

```go
func EnsureConnected(sugar *mt5.MT5Sugar, clusterName string) error {
    connected := sugar.IsConnected()

    if connected {
        return nil // Already connected
    }

    // Connection lost - reconnect
    fmt.Println("⚠️  Connection lost - reconnecting...")
    err := sugar.QuickConnect(clusterName)
    if err != nil {
        return fmt.Errorf("reconnection failed: %w", err)
    }

    fmt.Println("✅ Reconnected successfully")
    return nil
}

// Usage:
err := EnsureConnected(sugar, "FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("Cannot connect: %v\n", err)
    return
}

// Now safe to trade
sugar.BuyMarket("EURUSD", 0.1)
```

---

### 4) Connection monitoring loop

```go
func MonitorConnection(sugar *mt5.MT5Sugar, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        connected := sugar.IsConnected()

        if connected {
            fmt.Println("📡 Connection: OK")
        } else {
            fmt.Println("❌ Connection: LOST")
            // Trigger reconnect logic here
        }
    }
}

// Run in background
go MonitorConnection(sugar, 30*time.Second)
```

---

### 5) Pre-trading validation

```go
func ValidateTrading(sugar *mt5.MT5Sugar) error {
    // Check connection
    connected := sugar.IsConnected()
    if !connected {
        return fmt.Errorf("not connected to MT5")
    }

    // Check balance
    balance, err := sugar.GetBalance()
    if err != nil {
        return fmt.Errorf("cannot get balance: %w", err)
    }

    if balance <= 0 {
        return fmt.Errorf("insufficient balance: %.2f", balance)
    }

    return nil
}

// Before trading
err := ValidateTrading(sugar)
if err != nil {
    fmt.Printf("❌ Cannot trade: %v\n", err)
    return
}

fmt.Println("✅ All checks passed - ready to trade")
```

---

### 6) Retry with connection check

```go
func PlaceOrderWithRetry(sugar *mt5.MT5Sugar, symbol string, volume float64, maxRetries int) (uint64, error) {
    for i := 0; i < maxRetries; i++ {
        // Check connection before each attempt
        connected := sugar.IsConnected()
        if !connected {
            fmt.Printf("Attempt %d: Not connected - reconnecting...\n", i+1)
            sugar.QuickConnect("FxPro-MT5 Demo")
            time.Sleep(2 * time.Second)
            continue
        }

        // Try to place order
        ticket, err := sugar.BuyMarket(symbol, volume)
        if err == nil {
            return ticket, nil // Success
        }

        fmt.Printf("Attempt %d failed: %v\n", i+1, err)
        time.Sleep(time.Second)
    }

    return 0, fmt.Errorf("all %d attempts failed", maxRetries)
}

// Usage:
ticket, err := PlaceOrderWithRetry(sugar, "EURUSD", 0.1, 3)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
} else {
    fmt.Printf("Order #%d placed\n", ticket)
}
```

---

### 7) Status dashboard

```go
func PrintStatus(sugar *mt5.MT5Sugar) {
    connected := sugar.IsConnected()

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║          MT5 STATUS DASHBOARD         ║")
    fmt.Println("╚═══════════════════════════════════════╝")

    if connected {
        fmt.Println("Connection:  ✅ CONNECTED")

        balance, _ := sugar.GetBalance()
        equity, _ := sugar.GetEquity()
        positions, _ := sugar.CountOpenPositions()

        fmt.Printf("Balance:     %.2f\n", balance)
        fmt.Printf("Equity:      %.2f\n", equity)
        fmt.Printf("Positions:   %d\n", positions)
    } else {
        fmt.Println("Connection:  ❌ DISCONNECTED")
        fmt.Println("Balance:     N/A")
        fmt.Println("Equity:      N/A")
        fmt.Println("Positions:   N/A")
    }
}

// Usage:
PrintStatus(sugar)
```

---

### 8) Connection timeout handling

```go
func CheckConnectionWithTimeout(sugar *mt5.MT5Sugar, timeout time.Duration) (bool, error) {
    done := make(chan bool, 1)

    go func() {
        connected := sugar.IsConnected()
        done <- connected
    }()

    select {
    case connected := <-done:
        return connected, nil
    case <-time.After(timeout):
        return false, fmt.Errorf("connection check timeout after %v", timeout)
    }
}

// Usage with custom timeout
connected, err := CheckConnectionWithTimeout(sugar, 5*time.Second)
if err != nil {
    fmt.Printf("Check failed: %v\n", err)
}
```

---

### 9) Bot startup sequence

```go
func StartTradingBot(login uint64, password, server, cluster string) error {
    // Create Sugar
    sugar, err := mt5.NewMT5Sugar(login, password, server)
    if err != nil {
        return fmt.Errorf("failed to create Sugar: %w", err)
    }

    // Connect
    fmt.Println("🔌 Connecting to MT5...")
    err = sugar.QuickConnect(cluster)
    if err != nil {
        return fmt.Errorf("connection failed: %w", err)
    }

    // Verify connection
    fmt.Println("🔍 Verifying connection...")
    connected := sugar.IsConnected()
    if !connected {
        return fmt.Errorf("connection verification failed")
    }

    // Get account info
    fmt.Println("📊 Loading account info...")
    info, err := sugar.GetAccountInfo()
    if err != nil {
        return fmt.Errorf("failed to get account info: %w", err)
    }

    fmt.Println("╔═══════════════════════════════════════╗")
    fmt.Println("║       BOT STARTED SUCCESSFULLY        ║")
    fmt.Println("╚═══════════════════════════════════════╝")
    fmt.Printf("Account:  %d\n", info.Login)
    fmt.Printf("Balance:  %.2f %s\n", info.Balance, info.Currency)
    fmt.Printf("Company:  %s\n", info.Company)

    return nil
}

// Usage:
err := StartTradingBot(591129415, "password", "mt5.mrpc.pro:443", "FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("❌ Startup failed: %v\n", err)
}
```

---

### 10) Connection status with fallback

```go
func GetBalanceOrDefault(sugar *mt5.MT5Sugar) float64 {
    // Check connection first
    connected := sugar.IsConnected()
    if !connected {
        fmt.Println("⚠️  Not connected - returning 0")
        return 0.0
    }

    // Connection OK - get balance
    balance, err := sugar.GetBalance()
    if err != nil {
        fmt.Printf("⚠️  Balance error: %v - returning 0\n", err)
        return 0.0
    }

    return balance
}

// Safe to call even when disconnected
balance := GetBalanceOrDefault(sugar)
fmt.Printf("Balance: %.2f\n", balance)
```

---

## 🔗 Related Methods

**🍬 Connection methods:**

* `QuickConnect()` - Connect to MT5 Terminal
* `Ping()` - Detailed connection health check

**📖 Typical usage pattern:**

```go
// 1. Connect
sugar.QuickConnect("FxPro-MT5 Demo")

// 2. Verify
connected := sugar.IsConnected()
if !connected {
    return
}

// 3. Trade
sugar.BuyMarket("EURUSD", 0.1)
```

---

## ⚠️ Common Pitfalls

### 1) Ignoring the result

```go
// ❌ WRONG - ignoring connection status
sugar.IsConnected()
sugar.BuyMarket("EURUSD", 0.1) // Might fail!

// ✅ CORRECT - check before trading
connected := sugar.IsConnected()
if !connected {
    return
}
sugar.BuyMarket("EURUSD", 0.1)
```

### 2) Not checking after QuickConnect

```go
// ❌ WRONG - assuming connection worked
sugar.QuickConnect("FxPro-MT5 Demo")
// Start trading immediately...

// ✅ CORRECT - verify connection
err := sugar.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    return
}
connected := sugar.IsConnected()
if !connected {
    return
}
// Now safe to trade
```

### 3) Checking too frequently

```go
// ❌ WRONG - checking before every single operation
for i := 0; i < 1000; i++ {
    connected := sugar.IsConnected()
    if !connected { return }
    sugar.GetBalance() // Too many checks!
}

// ✅ CORRECT - check once before loop
connected := sugar.IsConnected()
if !connected { return }
for i := 0; i < 1000; i++ {
    sugar.GetBalance()
}
```

---

## 💎 Pro Tips

1. **Check once, trade many** - Don't check before every operation
2. **Periodic monitoring** - Check every 30-60 seconds in bots
3. **Before critical operations** - Always check before trading
4. **Reconnect on false** - Auto-reconnect when connection lost
5. **Combine with Ping()** - Use Ping() for detailed health checks

---

**See also:** [`QuickConnect.md`](QuickConnect.md), [`Ping.md`](Ping.md)
