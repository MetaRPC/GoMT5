# 🔌 Quick Connect to MT5 (`QuickConnect`)

> **Sugar method:** Easiest way to connect to MT5 Terminal - just provide cluster name and you're connected!

**API Information:**

* **Method:** `sugar.QuickConnect(clusterName)`
* **Package:** `mt5` (MT5Sugar), import `mt5 "github.com/MetaRPC/GoMT5/examples/mt5"` - it lives in the repository's
  *examples* module (`github.com/MetaRPC/GoMT5/examples`), so you get it by cloning the repository, not by `go get` of the SDK package
* **Underlying calls:** `account.ConnectEx()` with `pb.ConnectExRequest{User, Password, MtClusterName}`
* **Timeout:** 30 seconds

---

## 📋 Method Signature

```go
func (s *MT5Sugar) QuickConnect(clusterName string) error
```

---

## 🔽 Input

| Parameter | Type | Description |
|-----------|------|-------------|
| `clusterName` | `string` | MT5 cluster name (e.g., "FxPro-MT5 Demo", "FxPro-MT5 Live") |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| `error` | `error` | `nil` if connected successfully, error otherwise |

---

## 💬 Just the Essentials

* **What it is:** Simplest connection method - connects using just cluster name (login and password were given to `NewMT5Sugar`).
* **Why you need it:** You don't need to know the broker's MT5 server host/port - just the cluster name from MT5 terminal.
* **Sanity check:** After successful connect, `IsConnected()` returns `true`.

---

## 🎯 Purpose

Use it when you want the simplest connection:

* **Quick connection** - One parameter, one call
* **Auto-discovery** - The cluster name is resolved server-side; no MT5 server host/port needed
* **Beginner-friendly** - No need to configure the broker's host:port (the gRPC gateway address is the third
  argument of `NewMT5Sugar`; `""` means `"mt5.mrpc.pro:443"`)
* **Demo/Live switching** - Just change cluster name

---

## 🧩 Notes & Tips

* **Cluster name** - Same as in MT5 Terminal (File → Open Account)
* **Case sensitive** - Must match exactly
* **Demo vs Live** - Usually ends with "Demo" or "Live"
* **First call** - Must be called before any trading operations
* **One-time** - Only need to connect once per session
* **Check connection** - Use `IsConnected()` to verify

---

## 🔧 Under the Hood

```go
func (s *MT5Sugar) QuickConnect(clusterName string) error {
    ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
    defer cancel()

    req := &pb.ConnectExRequest{
        User:          s.user,
        Password:      s.password,
        MtClusterName: clusterName,
    }

    _, err := s.GetAccount().ConnectEx(ctx, req)
    return err
}
```

Only `User`, `Password` and `MtClusterName` are set - no base chart symbol, no `TimeoutSeconds`, no experts.
On success `ConnectEx` stores the returned terminal instance id in the account, and every later call uses it.

**What it improves:**

* ✅ **Simple API** - One parameter instead of multiple
* ✅ **No config needed** - Don't need to know server details
* ✅ **Automatic lookup** - Cluster → server resolution

---

## 📊 Low-Level Alternative

**WITHOUT sugar:**
```go
// helpers "github.com/MetaRPC/GoMT5/package/Helpers", pb "github.com/MetaRPC/GoMT5/package", "github.com/google/uuid"
account, _ := helpers.NewMT5Account(12345678, "your_password", "mt5.mrpc.pro:443", uuid.New())

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

_, err := account.ConnectEx(ctx, &pb.ConnectExRequest{
    User:          12345678,
    Password:      "your_password",
    MtClusterName: "YourBroker-MT5 Demo",
})
```

**WITH sugar:**
```go
sugar, _ := mt5.NewMT5Sugar(12345678, "your_password", "mt5.mrpc.pro:443")
err := sugar.QuickConnect("YourBroker-MT5 Demo")
```

**Benefits:**

* ✅ **Shorter and clearer**
* ✅ **Same simplicity**
* ✅ **Built into Sugar instance**

---

## 🔗 Usage Examples

### 1) Basic connection

```go
package main

import (
    "fmt"

    mt5 "github.com/MetaRPC/GoMT5/examples/mt5"
)

func main() {
    // Create Sugar instance
    sugar, err := mt5.NewMT5Sugar(
        591129415,
        "your_password",
        "mt5.mrpc.pro:443",
    )
    if err != nil {
        fmt.Printf("Failed to create Sugar: %v\n", err)
        return
    }

    // Connect to MT5
    err = sugar.QuickConnect("FxPro-MT5 Demo")
    if err != nil {
        fmt.Printf("Connection failed: %v\n", err)
        return
    }

    fmt.Println("✅ Connected successfully!")

    // Now you can use any Sugar methods
    balance, _ := sugar.GetBalance()
    fmt.Printf("Balance: %.2f\n", balance)
}
```

---

### 2) Connection with verification

```go
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")

// Connect
err := sugar.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("❌ Connection failed: %v\n", err)
    return
}

// Verify connection
connected := sugar.IsConnected()

if connected {
    fmt.Println("✅ Connected and verified!")
} else {
    fmt.Println("❌ Connection failed!")
}
```

---

### 3) Demo vs Live switching

```go
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")

// Connect to DEMO
err := sugar.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("Demo connection failed: %v\n", err)
    return
}

fmt.Println("✅ Connected to DEMO account")

// ... Later, to switch to LIVE: the cluster name is the only thing QuickConnect takes,
// login and password come from NewMT5Sugar - so a live account with its own login
// needs its own instance:
// sugarLive, _ := mt5.NewMT5Sugar(liveLogin, "live_password", "mt5.mrpc.pro:443")
// sugarLive.QuickConnect("FxPro-MT5 Live")
```

---

### 4) Connection retry logic

```go
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")

maxRetries := 3
clusterName := "FxPro-MT5 Demo"

for i := 0; i < maxRetries; i++ {
    fmt.Printf("Connection attempt %d/%d...\n", i+1, maxRetries)

    err := sugar.QuickConnect(clusterName)
    if err == nil {
        fmt.Println("✅ Connected successfully!")
        break
    }

    fmt.Printf("❌ Attempt %d failed: %v\n", i+1, err)

    if i < maxRetries-1 {
        fmt.Println("Retrying in 5 seconds...")
        time.Sleep(5 * time.Second)
    } else {
        fmt.Println("❌ All connection attempts failed")
        return
    }
}
```

---

### 5) Connection with timeout

```go
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")

// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Connect (connection has built-in timeout, but we add extra safety)
done := make(chan error, 1)

go func() {
    done <- sugar.QuickConnect("FxPro-MT5 Demo")
}()

select {
case err := <-done:
    if err != nil {
        fmt.Printf("Connection failed: %v\n", err)
        return
    }
    fmt.Println("✅ Connected successfully!")

case <-ctx.Done():
    fmt.Println("❌ Connection timeout")
    return
}
```

---

### 6) Get account info after connection

```go
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")

err := sugar.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("Connection failed: %v\n", err)
    return
}

// Get account information
accountInfo, err := sugar.GetAccountInfo()
if err != nil {
    fmt.Printf("Failed to get account info: %v\n", err)
    return
}

fmt.Println("╔═══════════════════════════════════════╗")
fmt.Println("║       ACCOUNT INFORMATION             ║")
fmt.Println("╚═══════════════════════════════════════╝")
fmt.Printf("Login:      %d\n", accountInfo.Login)
fmt.Printf("Balance:    %.2f %s\n", accountInfo.Balance, accountInfo.Currency)
fmt.Printf("Equity:     %.2f %s\n", accountInfo.Equity, accountInfo.Currency)
fmt.Printf("Leverage:   1:%d\n", accountInfo.Leverage)
fmt.Printf("Company:    %s\n", accountInfo.Company)
```

---

### 7) Connection from environment variables

```go
import (
    "os"
    "strconv"
)

// Get credentials from environment
loginStr := os.Getenv("MT5_USER")
password := os.Getenv("MT5_PASSWORD")
server := os.Getenv("MT5_GRPC_SERVER")
cluster := os.Getenv("MT5_MT_CLUSTER")

// Validate
if loginStr == "" || password == "" || server == "" || cluster == "" {
    fmt.Println("❌ Missing environment variables")
    return
}

// NewMT5Sugar takes the login as uint64
login, err := strconv.ParseUint(loginStr, 10, 64)
if err != nil {
    fmt.Printf("❌ MT5_USER is not a number: %v\n", err)
    return
}

// Connect
sugar, err := mt5.NewMT5Sugar(login, password, server)
if err != nil {
    fmt.Printf("Failed to create Sugar: %v\n", err)
    return
}

err = sugar.QuickConnect(cluster)
if err != nil {
    fmt.Printf("Connection failed: %v\n", err)
    return
}

fmt.Printf("✅ Connected to %s\n", cluster)
```

---

### 8) Connection status monitoring

```go
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")

// Initial connection
err := sugar.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("Connection failed: %v\n", err)
    return
}

fmt.Println("✅ Initial connection successful")

// Monitor connection in background
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        connected := sugar.IsConnected()
        if connected {
            fmt.Println("📡 Connection: OK")
        } else {
            fmt.Println("⚠️  Connection: LOST - Attempting reconnect...")
            err := sugar.QuickConnect("FxPro-MT5 Demo")
            if err == nil {
                fmt.Println("✅ Reconnected!")
            } else {
                fmt.Printf("❌ Reconnect failed: %v\n", err)
            }
        }
    }
}()

// Main program continues...
time.Sleep(5 * time.Minute)
```

---

### 9) Multiple account handling

```go
// Demo account
sugarDemo, _ := mt5.NewMT5Sugar(591129415, "demo_pass", "mt5.mrpc.pro:443")
err := sugarDemo.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("Demo connection failed: %v\n", err)
} else {
    fmt.Println("✅ Demo account connected")
}

// Live account (different credentials)
sugarLive, _ := mt5.NewMT5Sugar(123456789, "live_pass", "mt5.mrpc.pro:443")
err = sugarLive.QuickConnect("FxPro-MT5 Live")
if err != nil {
    fmt.Printf("Live connection failed: %v\n", err)
} else {
    fmt.Println("✅ Live account connected")
}

// Now you can use both accounts
demoBalance, _ := sugarDemo.GetBalance()
liveBalance, _ := sugarLive.GetBalance()

fmt.Printf("Demo balance: %.2f\n", demoBalance)
fmt.Printf("Live balance: %.2f\n", liveBalance)
```

---

### 10) Connection helper function

```go
func ConnectToMT5(login uint64, password, server, cluster string) (*mt5.MT5Sugar, error) {
    // Create Sugar instance
    sugar, err := mt5.NewMT5Sugar(login, password, server)
    if err != nil {
        return nil, fmt.Errorf("failed to create Sugar: %w", err)
    }

    // Connect
    err = sugar.QuickConnect(cluster)
    if err != nil {
        return nil, fmt.Errorf("connection failed: %w", err)
    }

    // Verify connection
    connected := sugar.IsConnected()

    if !connected {
        return nil, fmt.Errorf("connection not established")
    }

    return sugar, nil
}

// Usage:
sugar, err := ConnectToMT5(
    591129415,
    "password",
    "mt5.mrpc.pro:443",
    "FxPro-MT5 Demo",
)
if err != nil {
    fmt.Printf("❌ %v\n", err)
    return
}

fmt.Println("✅ Connected and ready to trade!")
```

---

## 🔗 Related Methods

**🍬 Other connection methods:**

* `IsConnected()` - Check if connected
* `Ping()` - Verify connection health

**📖 Next steps after connection:**

* `GetBalance()` - Check account balance
* `GetAccountInfo()` - Get complete account details
* `GetAllSymbols()` - List available symbols

---

## ⚠️ Common Pitfalls

### 1) Wrong cluster name

```go
// ❌ WRONG - incorrect cluster name
err := sugar.QuickConnect("FxPro Demo") // Missing "-MT5"

// ✅ CORRECT - exact cluster name from MT5
err := sugar.QuickConnect("FxPro-MT5 Demo")
```

### 2) Not checking connection error

```go
// ❌ WRONG - ignoring error
sugar.QuickConnect("FxPro-MT5 Demo")
balance, _ := sugar.GetBalance() // Will fail if not connected!

// ✅ CORRECT - check error
err := sugar.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("Connection failed: %v\n", err)
    return
}
balance, _ := sugar.GetBalance()
```

### 3) Connecting multiple times

```go
// ❌ WRONG - redundant connections
sugar.QuickConnect("FxPro-MT5 Demo")
sugar.QuickConnect("FxPro-MT5 Demo") // Already connected!

// ✅ CORRECT - connect once
err := sugar.QuickConnect("FxPro-MT5 Demo")
if err != nil {
    fmt.Printf("Connection failed: %v\n", err)
    return
}
// Now use Sugar methods without reconnecting
```

### 4) Using wrong credentials

```go
// ❌ WRONG - login in QuickConnect (it's already in NewMT5Sugar)
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")
sugar.QuickConnect("591129415") // WRONG! This is cluster, not login

// ✅ CORRECT - cluster name
sugar, _ := mt5.NewMT5Sugar(591129415, "password", "mt5.mrpc.pro:443")
sugar.QuickConnect("FxPro-MT5 Demo")
```

---

## 💎 Pro Tips

1. **Always check errors** - Connection can fail for many reasons
2. **Use environment variables** - Don't hardcode credentials
3. **Verify after connect** - Use `IsConnected()` to double-check
4. **One connection per session** - Don't reconnect unless necessary
5. **Demo first** - Always test on demo before live
6. **Cluster name** - Copy exact name from MT5 Terminal

---

**See also:** [`IsConnected.md`](IsConnected.md), [`Ping.md`](Ping.md)
