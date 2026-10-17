# 🏓 Ping MT5 Server (`Ping`)

> **Sugar method:** Asks the server for the terminal's health check (`CheckConnect`) and reports why the connection is not healthy.

**API Information:**

* **Method:** `sugar.Ping()`
* **Package:** `mt5` (MT5Sugar)
* **Underlying calls:** `account.CheckConnect()` - checks `HealthCheck.IsAlive` in the reply
* **Timeout:** 3 seconds

---

## 📋 Method Signature

```go
func (s *MT5Sugar) Ping() error
```

---

## 🔽 Input

*No parameters*

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `error` | `error` | `nil` if ping successful, error if connection problem |

---

## 💬 Just the Essentials

* **What it is:** Sends a ping to MT5 server to test if it's reachable and responsive.
* **Why you need it:** Verify connection health, measure latency, detect network issues.
* **Sanity check:** If ping returns `nil`, connection is healthy.

---

## 🎯 Purpose

Use it for detailed connection testing:

* **Health check** - Verify server is responsive
* **Latency measurement** - Test connection speed
* **Network diagnostics** - Detect connectivity issues
* **Pre-trading validation** - Ensure stable connection
* **Monitoring** - Track connection quality over time

---

## 🧩 Notes & Tips

* **Returns error** - `nil` = success, error = problem
* **3-second timeout** - Fails if no response
* **Same call as IsConnected** - Both send one `CheckConnect`; `Ping` returns the reason as an error instead of `false`
* **Use for diagnostics** - Better than IsConnected for troubleshooting
* **Measure latency** - Time the call to measure round-trip
* **Network issues** - Detects network problems quickly

---

## 🔧 Under the Hood

```go
func (s *MT5Sugar) Ping() error {
    ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
    defer cancel()

    data, err := s.GetAccount().CheckConnect(ctx, &pb.CheckConnectRequest{})
    if err != nil {
        return fmt.Errorf("ping failed: %w", err)
    }

    if data.HealthCheck == nil || !data.HealthCheck.IsAlive {
        return fmt.Errorf("not connected to MT5 terminal")
    }

    return nil
}
```

**What it does:**

* ✅ **Health check** - One `CheckConnect` round trip
* ✅ **3-second timeout** - Doesn't hang forever
* ✅ **Two error kinds** - `"ping failed: ..."` (the call failed or timed out) or `"not connected to MT5 terminal"` (the server answered, terminal not alive)
* ✅ **Clean error** - nil = OK, error = problem

---

## 📊 Low-Level Alternative

**WITHOUT sugar:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

data, err := account.CheckConnect(ctx, &pb.CheckConnectRequest{})
if err == nil && (data.HealthCheck == nil || !data.HealthCheck.IsAlive) {
    err = fmt.Errorf("not connected to MT5 terminal")
}
```

**WITH sugar:**
```go
err := sugar.Ping()
```

**Benefits:**

* ✅ **One line vs three**
* ✅ **Built-in timeout**
* ✅ **Cleaner code**

---

## 🔗 Usage Examples

### 1) Basic ping

```go
err := sugar.Ping()
if err != nil {
    fmt.Printf("❌ Ping failed: %v\n", err)
    return
}

fmt.Println("✅ Ping successful - server is responsive")
```

---

### 2) Measure latency

```go
start := time.Now()
err := sugar.Ping()
latency := time.Since(start)

if err != nil {
    fmt.Printf("❌ Ping failed: %v\n", err)
    return
}

fmt.Printf("✅ Ping successful\n")
fmt.Printf("   Latency: %v\n", latency)

if latency > 500*time.Millisecond {
    fmt.Println("   ⚠️  High latency detected!")
}
```

---

### 3) Connection quality check

```go
func CheckConnectionQuality(sugar *mt5.MT5Sugar, samples int) {
    var totalLatency time.Duration
    failures := 0

    fmt.Printf("Testing connection (%d pings)...\n", samples)

    for i := 0; i < samples; i++ {
        start := time.Now()
        err := sugar.Ping()
        latency := time.Since(start)

        if err != nil {
            failures++
            fmt.Printf("  Ping %d: ❌ Failed\n", i+1)
        } else {
            totalLatency += latency
            fmt.Printf("  Ping %d: ✅ %v\n", i+1, latency)
        }

        time.Sleep(time.Second)
    }

    fmt.Println("─────────────────────────────────")
    fmt.Printf("Results:\n")
    fmt.Printf("  Successful: %d/%d (%.1f%%)\n",
        samples-failures, samples,
        float64(samples-failures)/float64(samples)*100)

    if failures < samples {
        avgLatency := totalLatency / time.Duration(samples-failures)
        fmt.Printf("  Avg latency: %v\n", avgLatency)
    }

    if failures > 0 {
        fmt.Printf("  ⚠️  %d ping(s) failed\n", failures)
    }
}

// Usage:
CheckConnectionQuality(sugar, 5)
```

---

### 4) Pre-trading connection test

```go
func ValidateConnectionForTrading(sugar *mt5.MT5Sugar) error {
    // Test with 3 pings
    failures := 0

    for i := 0; i < 3; i++ {
        err := sugar.Ping()
        if err != nil {
            failures++
        }
        time.Sleep(500 * time.Millisecond)
    }

    if failures >= 2 {
        return fmt.Errorf("connection unstable: %d/3 pings failed", failures)
    }

    return nil
}

// Before trading
err := ValidateConnectionForTrading(sugar)
if err != nil {
    fmt.Printf("❌ %v\n", err)
    fmt.Println("   Trading not recommended")
    return
}

fmt.Println("✅ Connection stable - safe to trade")
```

---

### 5) Continuous monitoring

```go
func MonitorConnectionHealth(sugar *mt5.MT5Sugar, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    consecutiveFailures := 0

    for range ticker.C {
        err := sugar.Ping()

        if err != nil {
            consecutiveFailures++
            fmt.Printf("⚠️  Ping failed (%d consecutive failures): %v\n",
                consecutiveFailures, err)

            if consecutiveFailures >= 3 {
                fmt.Println("🚨 CONNECTION LOST - Too many failures!")
                // Trigger reconnect or alert
            }
        } else {
            if consecutiveFailures > 0 {
                fmt.Println("✅ Connection restored")
            }
            consecutiveFailures = 0
            fmt.Println("📡 Ping OK")
        }
    }
}

// Run in background
go MonitorConnectionHealth(sugar, 30*time.Second)
```

---

### 6) Compare Ping vs IsConnected

```go
fmt.Println("Testing connection methods:")

// Method 1: Ping
start := time.Now()
pingErr := sugar.Ping()
pingTime := time.Since(start)

// Method 2: IsConnected
start = time.Now()
connected := sugar.IsConnected()
connTime := time.Since(start)

fmt.Println("─────────────────────────────────")
fmt.Printf("Ping():        %v (%v)\n", pingErr == nil, pingTime)
fmt.Printf("IsConnected(): %v (%v)\n", connected, connTime)
fmt.Println("─────────────────────────────────")

// Output might be:
// Ping():        true (45ms)
// IsConnected(): true (48ms)
```

---

### 7) Retry logic with ping

```go
func ConnectWithPing(sugar *mt5.MT5Sugar, cluster string, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        fmt.Printf("Connection attempt %d/%d...\n", i+1, maxRetries)

        // Connect
        err := sugar.QuickConnect(cluster)
        if err != nil {
            fmt.Printf("  Connect failed: %v\n", err)
            time.Sleep(3 * time.Second)
            continue
        }

        // Verify with ping
        err = sugar.Ping()
        if err != nil {
            fmt.Printf("  Ping failed: %v\n", err)
            time.Sleep(3 * time.Second)
            continue
        }

        fmt.Println("✅ Connected and verified!")
        return nil
    }

    return fmt.Errorf("all %d connection attempts failed", maxRetries)
}

// Usage:
err := ConnectWithPing(sugar, "FxPro-MT5 Demo", 3)
if err != nil {
    fmt.Printf("❌ %v\n", err)
}
```

---

### 8) Benchmark connection speed

```go
func BenchmarkConnection(sugar *mt5.MT5Sugar, iterations int) {
    latencies := make([]time.Duration, 0, iterations)

    fmt.Printf("Benchmarking connection (%d pings)...\n", iterations)

    for i := 0; i < iterations; i++ {
        start := time.Now()
        err := sugar.Ping()
        latency := time.Since(start)

        if err == nil {
            latencies = append(latencies, latency)
        }

        time.Sleep(200 * time.Millisecond)
    }

    if len(latencies) == 0 {
        fmt.Println("❌ All pings failed")
        return
    }

    // Calculate statistics
    var total time.Duration
    min := latencies[0]
    max := latencies[0]

    for _, lat := range latencies {
        total += lat
        if lat < min {
            min = lat
        }
        if lat > max {
            max = lat
        }
    }

    avg := total / time.Duration(len(latencies))

    fmt.Println("═══════════════════════════════════")
    fmt.Println("  CONNECTION BENCHMARK RESULTS")
    fmt.Println("═══════════════════════════════════")
    fmt.Printf("Pings sent:    %d\n", iterations)
    fmt.Printf("Successful:    %d (%.1f%%)\n",
        len(latencies),
        float64(len(latencies))/float64(iterations)*100)
    fmt.Printf("Min latency:   %v\n", min)
    fmt.Printf("Avg latency:   %v\n", avg)
    fmt.Printf("Max latency:   %v\n", max)
}

// Usage:
BenchmarkConnection(sugar, 20)
```

---

### 9) Alert on connection issues

```go
func MonitorWithAlerts(sugar *mt5.MT5Sugar) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    lastStatus := true

    for range ticker.C {
        err := sugar.Ping()
        currentStatus := (err == nil)

        // Status changed
        if currentStatus != lastStatus {
            if currentStatus {
                fmt.Println("🔔 ALERT: Connection RESTORED")
                // Send notification (email, telegram, etc.)
            } else {
                fmt.Println("🚨 ALERT: Connection LOST")
                // Send urgent notification
            }
        }

        lastStatus = currentStatus
    }
}

// Run in background
go MonitorWithAlerts(sugar)
```

---

### 10) Connection diagnostics tool

```go
func DiagnoseConnection(sugar *mt5.MT5Sugar) {
    fmt.Println("╔═══════════════════════════════════════════════╗")
    fmt.Println("║       CONNECTION DIAGNOSTICS                  ║")
    fmt.Println("╚═══════════════════════════════════════════════╝")
    fmt.Println()

    // Test 1: Ping
    fmt.Println("Test 1: Server Ping")
    start := time.Now()
    pingErr := sugar.Ping()
    pingLatency := time.Since(start)

    if pingErr != nil {
        fmt.Printf("  ❌ FAILED: %v\n", pingErr)
    } else {
        fmt.Printf("  ✅ SUCCESS (Latency: %v)\n", pingLatency)
    }
    fmt.Println()

    // Test 2: IsConnected
    fmt.Println("Test 2: Connection Status")
    connected := sugar.IsConnected()
    if connected {
        fmt.Println("  ✅ CONNECTED")
    } else {
        fmt.Println("  ❌ NOT CONNECTED")
    }
    fmt.Println()

    // Test 3: Get Balance (actual operation)
    fmt.Println("Test 3: Sample Operation (GetBalance)")
    balance, balErr := sugar.GetBalance()
    if balErr != nil {
        fmt.Printf("  ❌ FAILED: %v\n", balErr)
    } else {
        fmt.Printf("  ✅ SUCCESS (Balance: %.2f)\n", balance)
    }
    fmt.Println()

    // Summary
    fmt.Println("═══════════════════════════════════════════════")
    allPassed := (pingErr == nil && connected && balErr == nil)
    if allPassed {
        fmt.Println("✅ ALL TESTS PASSED - Connection is healthy")
    } else {
        fmt.Println("❌ SOME TESTS FAILED - Connection has issues")
    }
}

// Usage:
DiagnoseConnection(sugar)
```

---

## 🔗 Related Methods

**🍬 Connection methods:**

* `QuickConnect()` - Connect to MT5 Terminal
* `IsConnected()` - Check if connected

**📖 Typical diagnostic pattern:**

```go
// 1. Check with IsConnected
connected := sugar.IsConnected()

// 2. If issues, use Ping for details
if !connected {
    err := sugar.Ping()
    fmt.Printf("Ping error: %v\n", err)
}
```

---

## ⚠️ Common Pitfalls

### 1) Confusing nil with success

```go
// ❌ WRONG - checking for non-nil (backwards!)
if sugar.Ping() != nil {
    fmt.Println("Connection OK") // WRONG!
}

// ✅ CORRECT - nil means success
if sugar.Ping() == nil {
    fmt.Println("Connection OK")
}
```

### 2) Not handling timeout

```go
// ❌ WRONG - might block for 3 seconds
sugar.Ping() // Blocked for up to 3 seconds!

// ✅ CORRECT - be aware of timeout
fmt.Println("Pinging...")
err := sugar.Ping() // Max 3 seconds
if err != nil {
    fmt.Println("Ping failed or timed out")
}
```

### 3) Pinging too frequently

```go
// ❌ WRONG - excessive pinging
for {
    sugar.Ping() // Every loop iteration!
    time.Sleep(100 * time.Millisecond)
}

// ✅ CORRECT - reasonable interval
for {
    sugar.Ping()
    time.Sleep(30 * time.Second) // Every 30 seconds
}
```

---

## 💎 Pro Tips

1. **Use for diagnostics** - Better than IsConnected for troubleshooting
2. **Measure latency** - Time the call to measure connection speed
3. **Reasonable intervals** - Don't ping more than once per 10-30 seconds
4. **Before critical ops** - Ping before important trading operations
5. **Combine with monitoring** - Track latency trends over time
6. **Network issues** - Ping failures often indicate network problems

---

**See also:** [`QuickConnect.md`](QuickConnect.md), [`IsConnected.md`](IsConnected.md)
