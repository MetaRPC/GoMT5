# Your First Project in 10 Minutes

> **Practice Before Theory** - create a working MT5 trading project before diving into documentation

---

## Why This Guide?

I want to show you through a simple example how easy it is to use our gRPC gateway for working with MetaTrader 5.

**Before you dive into learning the basics and project fundamentals - let's create your first project.**

We'll download one Go module `package`, which contains:

- ✅ Protobuf definitions of all MT5 methods
- ✅ MT5Account - ready-to-use gRPC client
- ✅ Error handler - ApiError types and return codes
- ✅ Everything necessary to get started

**This is the foundation** for your future algorithmic trading system.

---


> 💡 After you get your first results, proceed to [GETTING_STARTED.md](./GETTING_STARTED.md) for a deep understanding of the SDK architecture.

---

## Step 1: Install Go 1.21 or Higher

If you don't have Go installed yet:

**Download and install:**

- [Go Download](https://go.dev/dl/)

**Verify installation:**

```bash
go version
# Should show: go version go1.21.x or higher
```

---

## Step 2: Create a New Go Project

Open terminal (command line) and execute:

```bash
# Create project folder
mkdir MyMT5Project
cd MyMT5Project

# Initialize Go module
go mod init mymt5project
```

**What happened:**

- ✅ Created `MyMT5Project` folder
- ✅ Created `go.mod` file inside - your project manifest
- ✅ Now you can add dependencies

---

## Step 3: Install the GoMT5 package Module

This is the most important step - installing the **single module** that contains everything you need:

```bash
go get github.com/MetaRPC/GoMT5/package
```

> **📌 Important for beginners:** After running the command, you will **NOT see new folders** in your project.
> This is normal! Go modules are installed in the system cache (`C:\Users\<username>\go\pkg\mod`), not copied to the project.

**How to verify the installation was successful?**

**Method 1:** Open the `go.mod` file and make sure a line with the package appeared:

```go
github.com/MetaRPC/GoMT5/package v0.0.0-XXXXXXXXXXXXXXXX-XXXXXXXXXXXX
```

(Version may differ - this is normal)

**Method 2:** Run the verification command:

```bash
go list -m github.com/MetaRPC/GoMT5/package
```

If you see the package version - **everything is installed correctly!** ✅

---

## Step 4: Create config.json Configuration File

Create a `config.json` file in the project root:

```json
{
  "user": 591129415,
  "password": "YourPassword123",
  "host": "mt5.mrpc.pro",
  "port": 443,
  "grpc_server": "mt5.mrpc.pro:443",
  "mt_cluster": "YourBroker-MT5 Demo",
  "test_symbol": "EURUSD",
  "test_volume": 0.01
}
```

**Parameter explanations:**

| Parameter | Description | Where to Get |
|----------|----------|-----------|
| **user** | Your MT5 account number (login) | In MT5 terminal: Tools → Options → Login |
| **password** | Master password for MT5 account | The one you received during registration |
| **host** | gRPC gateway host | `mt5.mrpc.pro` (public gateway) |
| **port** | gRPC gateway port | `443` (standard HTTPS port) |
| **grpc_server** | Full gRPC gateway address | `mt5.mrpc.pro:443` (host:port) |
| **mt_cluster** | Your broker's cluster name | In MT5 terminal: see server name |
| **test_symbol** | Trading symbol for examples | `EURUSD`, `GBPUSD`, `BTCUSD`, etc. |
| **test_volume** | Volume for test orders | `0.01` (minimum lot) |

**Replace:**

- `user`, `password`, `mt_cluster` - with your MT5 demo account data
- `grpc_server`, `host`, `port` - can be left as is (MetaRPC public gateway)


---

## Step 5: Create Configuration Structure

Create a `config.go` file in the project root:

```go
package main

import (
	"encoding/json"
	"os"
)

// Config contains MT5 connection settings
type Config struct {
	User       uint64  `json:"user"`
	Password   string  `json:"password"`
	Host       string  `json:"host"`
	Port       int     `json:"port"`
	GrpcServer string  `json:"grpc_server"`
	MtCluster  string  `json:"mt_cluster"`
	TestSymbol string  `json:"test_symbol"`
	TestVolume float64 `json:"test_volume"`
}

// LoadConfig loads configuration from JSON file
func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
```

---

## Step 6: Write Code to Connect and Get Account Information

Create a `main.go` file in the project root:

```go
package main

import (
	"context"
	"fmt"
	"log"

	mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
	pb "github.com/MetaRPC/GoMT5/package"
)

func main() {

	// ============================================================================
	// STEP 1: LOAD CONFIGURATION
	// ============================================================================

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("❌ Error loading config.json: %v", err)
	}

	fmt.Println("📋 Connection configuration:")
	fmt.Printf("   User: %d\n", config.User)
	fmt.Printf("   Cluster: %s\n", config.MtCluster)
	fmt.Printf("   gRPC Server: %s\n", config.GrpcServer)
	fmt.Printf("   Symbol: %s\n\n", config.TestSymbol)

	// ============================================================================
	// STEP 2: CREATE MT5ACCOUNT
	// ============================================================================

	fmt.Println("🔌 Creating MT5Account...")

	// Create MT5Account with auto-generated UUID
	account, err := mt5.NewMT5AccountAuto(
		config.User,
		config.Password,
		config.GrpcServer,
	)
	if err != nil {
		log.Fatalf("❌ Error creating MT5Account: %v", err)
	}
	defer account.Close()

	fmt.Printf("✅ MT5Account created (UUID: %s)\n\n", account.Id.String())

	// ============================================================================
	// STEP 3: CONNECT TO MT5
	// ============================================================================

	fmt.Println("🔗 Connecting to MT5 terminal...")

	ctx := context.Background()

	// Prepare connection request
	connectReq := &pb.ConnectExRequest{
		User:          config.User,
		Password:      config.Password,
		MtClusterName: config.MtCluster,
	}

	// Execute connection
	connectData, err := account.ConnectEx(ctx, connectReq)
	if err != nil {
		log.Fatalf("❌ Connection error: %v", err)
	}

	fmt.Printf("✅ Connected successfully!\n")
	fmt.Printf("   Instance ID: %s\n\n", connectData.TerminalInstanceGuid)

	// ============================================================================
	// STEP 4: GET ACCOUNT INFORMATION
	// ============================================================================

	// Create account information request
	accountReq := &pb.AccountSummaryRequest{}

	// Execute request
	accountData, err := account.AccountSummary(ctx, accountReq)
	if err != nil {
		log.Fatalf("❌ Error getting account data: %v", err)
	}

	// ============================================================================
	// STEP 5: OUTPUT RESULTS
	// ============================================================================

	fmt.Println("\n╔════════════════════════════════════════════════════════╗")
	fmt.Println("║              ACCOUNT INFORMATION                       ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("   Login:              %d\n", accountData.AccountLogin)
	fmt.Printf("   Username:           %s\n", accountData.AccountUserName)
	fmt.Printf("   Company:            %s\n", accountData.AccountCompanyName)
	fmt.Printf("   Currency:           %s\n", accountData.AccountCurrency)
	fmt.Println()
	fmt.Printf("💰 Balance:            %.2f %s\n", accountData.AccountBalance, accountData.AccountCurrency)
	fmt.Printf("💎 Equity:             %.2f %s\n", accountData.AccountEquity, accountData.AccountCurrency)
	fmt.Println()
	fmt.Printf("   Credit:             %.2f %s\n", accountData.AccountCredit, accountData.AccountCurrency)
	fmt.Printf("   Leverage:           1:%d\n", accountData.AccountLeverage)
	fmt.Println()

	if accountData.ServerTime != nil {
		serverTime := accountData.ServerTime.AsTime()
		fmt.Printf("   Server time:        %s\n", serverTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   UTC timezone:       %+d minutes\n", accountData.UtcTimezoneServerTimeShiftMinutes)
	}

	fmt.Println()
	fmt.Println("╚════════════════════════════════════════════════════════╝")

	// ============================================================================
	// STEP 6: DISCONNECT FROM MT5
	// ============================================================================

	fmt.Println("\n🔌 Disconnecting from MT5...")

	disconnectReq := &pb.DisconnectRequest{}
	_, err = account.Disconnect(ctx, disconnectReq)
	if err != nil {
		log.Printf("⚠️  Warning on disconnect: %v", err)
	} else {
		fmt.Println("✅ Disconnected successfully!")
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════╗")
	fmt.Println("║   🎉 CONGRATULATIONS! YOUR FIRST PROJECT WORKS! 🎉     ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝\n")
}
```

---

## Step 7: Run the Project

Save all files and execute:

```bash
go mod tidy
go run .
```

**What does `go mod tidy` do?**

This command automatically:

- ✅ Downloads all missing dependencies (gRPC, protobuf, UUID)
- ✅ Removes unused dependencies from `go.mod`
- ✅ Updates `go.sum` file with checksums

Now, when all code files are created, `go mod tidy` will see the real imports and pull in the necessary dependencies.


**Expected result:**

```
╔════════════════════════════════════════════════════════╗
║         GoMT5 - Your First Project with MT5            ║
╚════════════════════════════════════════════════════════╝

📋 Connection configuration:
   User: 591129415
   Cluster: FxPro-MT5 Demo
   gRPC Server: mt5.mrpc.pro:443
   Symbol: EURUSD

🔌 Creating MT5Account...
✅ MT5Account created (UUID: 12345678-90ab-cdef-1234-567890abcdef)

🔗 Connecting to MT5 terminal...
✅ Connected successfully!
   Instance ID: 98765432-10ab-cdef-5678-901234567890


╔════════════════════════════════════════════════════════╗
║              ACCOUNT INFORMATION                       ║
╚════════════════════════════════════════════════════════╝

   Login:              591129415
   Username:           Demo User
   Company:            FxPro Financial Services Ltd
   Currency:           USD

💰 Balance:            10000.00 USD
💎 Equity:             10000.00 USD

   Credit:             0.00 USD
   Leverage:           1:100

   Server time:        2025-01-22 15:30:45
   UTC timezone:       +120 minutes

╚════════════════════════════════════════════════════════╝

🔌 Disconnecting from MT5...
✅ Disconnected successfully!

╔════════════════════════════════════════════════════════╗
║   🎉 CONGRATULATIONS! YOUR FIRST PROJECT WORKS! 🎉     ║
╚════════════════════════════════════════════════════════╝
```

---

## 🎉 Congratulations! You Did It!

You just:

✅ Created a new Go project from scratch
✅ Integrated the **single** Go module `package` for working with MT5
✅ Configured connection settings
✅ Connected to MT5 terminal via gRPC
✅ Got complete account information programmatically

**This was a low-level approach** using `MT5Account` and protobuf directly.

---

## 📁 Your Project Structure

After completing all steps, your project structure should look like this:

```
MyMT5Project/
├── config.json          # MT5 connection configuration
├── config.go            # Load configuration from JSON
├── main.go              # Main application code
├── go.mod               # Go module with dependencies
└── go.sum               # Dependency checksums
```

**Contents of go.mod:**

```go
module mymt5project

go 1.21

require (
	github.com/MetaRPC/GoMT5/package v0.0.0-20260120212705-d4be7827736c // indirect
	github.com/google/uuid v1.6.0 // indirect
	// ... other automatically installed dependencies (grpc, protobuf, etc.)
)
```

---

## 🚀 What's Next?

Now that you have a working project, you can:

### 1. Add More Functionality

**Examples of what you can do:**

#### Get Current Quotes

```go
tickReq := &pb.SymbolInfoTickRequest{
	Symbol: config.TestSymbol,
}

tick, err := account.SymbolInfoTick(ctx, tickReq)
if err == nil {
	fmt.Printf("Bid: %.5f, Ask: %.5f\n", tick.Bid, tick.Ask)
}
```

#### Get All Open Positions

```go
ordersReq := &pb.OpenedOrdersRequest{
	// InputSortMode is optional (default 0 - sort by open time)
}

ordersData, err := account.OpenedOrders(ctx, ordersReq)
if err == nil {
	fmt.Printf("Open positions: %d\n", len(ordersData.PositionInfos))
	for _, pos := range ordersData.PositionInfos {
		fmt.Printf("  #%d %s %.2f lots, Profit: %.2f\n",
			pos.Ticket, pos.Symbol, pos.Volume, pos.Profit)
	}
}
```

**Optional:** If you need sorting, use `InputSortMode`:
```go
ordersReq := &pb.OpenedOrdersRequest{
	InputSortMode: 0, // 0 = by time (ASC), 1 = by time (DESC), 2 = by ticket (ASC), 3 = by ticket (DESC)
}
```

#### Open a Market Order

```go
comment := "GoMT5 Test Order"
orderReq := &pb.OrderSendRequest{
	Symbol:    config.TestSymbol,
	Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
	Volume:    0.01, // 0.01 lot
	Comment:   &comment,
}

orderData, err := account.OrderSend(ctx, orderReq)
if err == nil {
	fmt.Printf("Order opened: Deal #%d, Order #%d\n", orderData.Deal, orderData.Order)
}
```

#### Streaming Data

```go
// Subscribe to real-time ticks (with limit: 5 seconds or 10 events)
tickReq := &pb.OnSymbolTickRequest{
	SymbolNames: []string{config.TestSymbol},
}

dataChan, errChan := account.OnSymbolTick(ctx, tickReq)

// Create context with 5-second timeout
streamCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

eventCount := 0
maxEvents := 10

fmt.Println("🔄 Receiving tick stream (maximum 5 seconds or 10 events)...")

for {
	select {
	case data := <-dataChan:
		tick := data.SymbolTick
		eventCount++
		fmt.Printf("[%d] %s - Bid: %.5f, Ask: %.5f\n",
			eventCount, time.Now().Format("15:04:05"), tick.Bid, tick.Ask)

		if eventCount >= maxEvents {
			fmt.Println("✅ Received maximum number of events")
			return
		}
	case err := <-errChan:
		log.Printf("❌ Stream error: %v", err)
		return
	case <-streamCtx.Done():
		fmt.Printf("✅ Received %d events in 5 seconds\n", eventCount)
		return
	}
}
```

### 2. Learn the Complete SDK Architecture

The GoMT5 repository has **three API layers**:

```
┌─────────────────────────────────────────────┐
│  MT5Sugar (Layer 3) - Convenient API       │
│  examples/mt5/MT5Sugar.go                   │
│  sugar.BuyMarket("EURUSD", 0.01)            │
└─────────────────────────────────────────────┘
              ↓ uses
┌─────────────────────────────────────────────┐
│  MT5Service (Layer 2) - Wrappers            │
│  examples/mt5/MT5Service.go                 │
│  service.GetBalance()                       │
└─────────────────────────────────────────────┘
              ↓ uses
┌─────────────────────────────────────────────┐
│  MT5Account (Layer 1) - Base gRPC ⭐        │
│  package/Helpers/MT5Account.go              │
│  account.AccountSummary(ctx, req)           │
└─────────────────────────────────────────────┘
```

**You just used Layer 1 (MT5Account)** - this is the foundation of everything!

To use layers 2 and 3:

- Clone the GoMT5 repository
- Study [GETTING_STARTED.md](./GETTING_STARTED.md)
- Look at examples in `examples/demos/`

### 3. Study Ready-Made Examples

The GoMT5 repository has many examples:

- `examples/demos/lowlevel/` - examples with MT5Account (what you used)
- `examples/demos/service/` - examples with MT5Service
- `examples/demos/sugar/` - examples with MT5Sugar

### 4. Read Documentation

- [MT5Account API Reference](../API_Reference/MT5Account.md) - ⭐ complete reference for the base level
- [PROJECT_MAP.md](./PROJECT_MAP.md) - project map and architecture
- [GRPC_STREAM_MANAGEMENT.md](./GRPC_STREAM_MANAGEMENT.md) - working with streaming data
- [RETURN_CODES_REFERENCE.md](./RETURN_CODES_REFERENCE.md) - operation return codes

---


### What is the `package` Module?

`package` is an **independent Go module** that contains:

- MT5Account (base gRPC client)
- All protobuf definitions of MT5 API
- gRPC stubs for all methods
- Helper types and structures

This is a **portable module** - you can use it in any Go project!


### How to Work with Environment Variables Instead of config.json?

You can use environment variables:

```go
import (
	"os"
	"strconv"
)

func LoadConfigFromEnv() (*Config, error) {
	userStr := os.Getenv("MT5_USER")
	user, err := strconv.ParseUint(userStr, 10, 64)
	if err != nil {
		return nil, err
	}

	portStr := os.Getenv("MT5_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 443 // default value
	}

	volumeStr := os.Getenv("MT5_TEST_VOLUME")
	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		volume = 0.01 // default value
	}

	return &Config{
		User:       user,
		Password:   os.Getenv("MT5_PASSWORD"),
		Host:       os.Getenv("MT5_HOST"),
		Port:       port,
		GrpcServer: os.Getenv("MT5_GRPC_SERVER"),
		MtCluster:  os.Getenv("MT5_CLUSTER"),
		TestSymbol: os.Getenv("MT5_TEST_SYMBOL"),
		TestVolume: volume,
	}, nil
}
```

**Set variables:**

```bash

# Windows (PowerShell)
$env:MT5_USER="591129415"
$env:MT5_PASSWORD="YourPassword123"
$env:MT5_HOST="mt5.mrpc.pro"
$env:MT5_PORT="443"
$env:MT5_GRPC_SERVER="mt5.mrpc.pro:443"
$env:MT5_CLUSTER="FxPro-MT5 Demo"
$env:MT5_TEST_SYMBOL="EURUSD"
$env:MT5_TEST_VOLUME="0.01"
```

### How to Use Layer 2 (MT5Service) and Layer 3 (MT5Sugar)?

These layers are in the **main GoMT5 repository**:

1. Clone the repository:

   ```bash
   git clone https://github.com/MetaRPC/GoMT5.git
   ```

2. Copy needed files to your project:

   - `examples/mt5/MT5Service.go` (Layer 2)
   - `examples/mt5/MT5Sugar.go` (Layer 3)

3. Use convenient methods:

   ```go
   // Layer 2 - Service
   service := mt5.NewMT5Service(account)
   balance, _ := service.GetBalance()

   // Layer 3 - Sugar
   sugar := mt5.NewMT5Sugar(service)
   ticket, _ := sugar.BuyMarket("EURUSD", 0.01)
   ```

See details in [GETTING_STARTED.md](./GETTING_STARTED.md)

---

## 📝 Summary: What We Did

In this guide, you created a minimalist project that:

1. ✅ **Uses only Go modules** - doesn't require cloning the repository

2. ✅ **Imports the package module** - the only dependency for MT5

3. ✅ **Connects to MT5** via gRPC gateway

4. ✅ **Reads configuration** from `config.json`

5. ✅ **Uses MT5Account** (Layer 1 - base level)

6. ✅ **Gets account information** and outputs to console

**This is the foundation** for any of your MT5 projects in Go.

---


**Good luck developing trading systems! 🚀**

"Trade safely, code cleanly, and may your gRPC connections always be stable."
