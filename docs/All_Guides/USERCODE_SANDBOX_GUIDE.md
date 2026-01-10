# User Code Sandbox Guide

## What is this?

This is **your sandbox** for writing custom MT5 trading code in Go. Connection setup is already done - you only need to add your trading logic!

## Quick Start

1. **Case already active in main.go**
   - Open `examples/demos/main.go`
   - Case "18" for user code sandbox is ready to use

2. **Open `usercode/18_usercode.go`**
   - Write your code or uncomment examples

3. **Run:**
   ```bash
   cd examples/demos
   go run main.go 18
   ```

That's it! Your code will execute with full MT5 connection.

## How to use

### Option 1: Uncomment examples

The file contains 3 ready-to-use examples:

```go
// Example 1: Get account balance (Sugar - easiest)
// balance, err := sugar.GetBalance()
// if err != nil {
//     fmt.Printf("Error: %v\n", err)
//     return
// }
// fmt.Printf("Balance: %.2f\n", balance)
```

Just remove `//` to activate!

### Option 2: Write your own code

Add your logic between the markers:

```go
// ═══════════════════════════════════════════════════════════
// YOUR CODE HERE
// ═══════════════════════════════════════════════════════════

// Your trading strategy here...

// TODO: Your code here
```

## Available commands

Run your code with:

```bash
# From examples/demos directory
go run main.go 18
go run main.go usercode
go run main.go user
go run main.go sandbox
```

## What's already configured

✓ **Connection** - MT5 terminal connected via gRPC

✓ **Configuration** - Loaded from `config/config.json` or environment variables

✓ **MT5Account** - Low-level gRPC client (variable: `account`)

✓ **MT5Service** - Mid-level wrapper (variable: `service`)

✓ **MT5Sugar** - High-level Sugar API (variable: `sugar`)

✓ **Context** - For timeouts and cancellation (variable: `ctx`)

## Can I mix API levels?

**Yes!** You can use all three levels in the same file:

```go
// Low-level (direct gRPC protobuf)
req := &pb.AccountSummaryRequest{}
reply, err := account.AccountSummary(ctx, req)

// Mid-level (MT5Service)
accountInfo, err := service.GetAccountSummary(ctx)

// High-level (MT5Sugar)
balance, err := sugar.GetBalance()
```

All three variables are available simultaneously:

- `account` - for full control (gRPC protobuf)
- `service` - for convenient methods (Service layer)
- `sugar` - for simplest usage (Sugar API)

## Quick reference

### Get account information

```go
// Sugar API (easiest)
balance, err := sugar.GetBalance()
fmt.Printf("Balance: %.2f\n", balance)

// Service API (more details)
accountInfo, err := service.GetAccountSummary(ctx)
fmt.Printf("Balance: %.2f %s\n", accountInfo.Balance, accountInfo.Currency)
fmt.Printf("Equity:  %.2f\n", accountInfo.Equity)

// Account API (full control)
req := &pb.AccountSummaryRequest{}
reply, err := account.AccountSummary(ctx, req)
fmt.Printf("Balance: %.2f\n", reply.AccountBalance)
```

### Get current price

```go
// Sugar API
bid, _ := sugar.GetBid("EURUSD")
ask, _ := sugar.GetAsk("EURUSD")
fmt.Printf("EURUSD: Bid=%.5f, Ask=%.5f\n", bid, ask)

// Service API
tick, err := service.GetSymbolTick(ctx, "EURUSD")
fmt.Printf("EURUSD Bid: %.5f, Ask: %.5f\n", tick.Bid, tick.Ask)
```

### Open market order

```go
// Sugar API (easiest)
ticket, err := sugar.BuyMarket("EURUSD", 0.01)
if err != nil {
    fmt.Printf("Order failed: %v\n", err)
} else {
    fmt.Printf("Order opened: #%d\n", ticket)
}

// With SL/TP in pips
ticket, err := sugar.BuyMarketWithPips("EURUSD", 0.01, 20, 30)

// Calculate position size based on risk (2% with 50 pip SL)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 50)
ticket, err := sugar.BuyMarketWithPips("EURUSD", lotSize, 50, 100)
```

### Place pending order

```go
// Buy Limit 20 pips below current Ask
ticket, err := sugar.BuyLimitPips("EURUSD", 0.01, 20, 15, 30)

// Buy Stop 20 pips above current Ask
ticket, err := sugar.BuyStopPips("EURUSD", 0.01, 20, 15, 30)

// Sell Limit with absolute price
ticket, err := sugar.SellLimit("EURUSD", 0.01, 1.1050, 1.1070, 1.1030)
```

### Get open positions

```go
// Sugar API
positions, err := sugar.GetOpenPositions()
fmt.Printf("Open positions: %d\n", len(positions))
for _, pos := range positions {
    fmt.Printf("  #%d %s %.2f lot, Profit: %.2f\n",
        pos.Ticket, pos.Symbol, pos.Volume, pos.Profit)
}

// Service API with filters
positions, err := service.GetOpenedOrders(ctx, "EURUSD", "")
```

### Close positions

```go
// Close specific position
err := sugar.ClosePosition(ticket)

// Close all positions for symbol
err := sugar.CloseAllPositions("EURUSD")

// Close all positions (all symbols)
err := sugar.CloseAllPositions("")
```

### Modify position

```go
// Modify SL/TP
err := sugar.ModifyPosition(ticket, 1.0850, 1.0950)

// Set trailing stop
err := sugar.SetTrailingStop(ticket, 20)  // 20 pips trailing
```

### Calculate position size

```go
// Calculate position size based on risk
volume, err := sugar.CalculatePositionSize("EURUSD", 2.0, 20)
fmt.Printf("Volume for 2%% risk with 20 pip SL: %.2f lot\n", volume)

// Then use calculated volume
ticket, err := sugar.BuyMarket("EURUSD", volume)
```

## Return Codes (RetCodes)

**Always check RetCode after trading operations!**

```go
result, err := service.BuyMarket(ctx, "EURUSD", 0.01, 0, 0, 0)
if err != nil {
    fmt.Printf("gRPC error: %v\n", err)
    return
}

// Check RetCode
if result.ReturnedCode == 10009 {  // Success for market orders
    fmt.Printf("✅ Order opened: #%d\n", result.Order)
} else if result.ReturnedCode == 10008 {  // Success for pending orders
    fmt.Printf("✅ Pending order placed: #%d\n", result.Order)
} else {
    fmt.Printf("❌ Order failed: %s (code %d)\n",
        result.ReturnedCodeDescription, result.ReturnedCode)
}
```

**Common RetCodes:**

- `10009` - Market order executed successfully
- `10008` - Pending order placed successfully
- `10019` - Insufficient money (insufficient margin)
- `10016` - Invalid stops (SL/TP too close to price)
- `10018` - Market closed

Full list: [RETURN_CODES_REFERENCE.md](RETURN_CODES_REFERENCE.md)

## Error handling

GoMT5 has two types of errors you must check:

### 1. gRPC transport errors
```go
result, err := service.BuyMarket(ctx, "EURUSD", 0.01, 0, 0, 0)
if err != nil {
    fmt.Printf("Connection or network error: %v\n", err)
    return
}
```

### 2. Trading operation errors (RetCode)

```go
if result.ReturnedCode != 10009 {
    fmt.Printf("Trade rejected: %s\n", result.ReturnedCodeDescription)
}
```

**Always check BOTH!**

## Documentation

- [MT5Account Master Overview](../MT5Account/MT5Account.Master.Overview.md) - Complete API reference (40+ methods)
- [MT5Service API Overview](../MT5Service/MT5Service.Overview.md) - Mid-level wrapper (50+ methods)
- [MT5Sugar API Overview](../MT5Sugar/MT5Sugar.API_Overview.md) - High-level Sugar API (62 methods)
- [Protobuf Inspector Guide](PROTOBUF_INSPECTOR_GUIDE.md) - Interactive type explorer

## Configuration

### Method 1: config.json (Recommended)

Create `examples/demos/config/config.json`:

```json
{
  "user": 591129415,
  "password": "YourPassword",
  "host": "mt5.mrpc.pro",
  "port": 443,
  "grpc_server": "mt5.mrpc.pro:443",
  "mt_cluster": "FxPro-MT5 Demo",
  "test_symbol": "EURUSD",
  "test_volume": 0.01
}
```

### Method 2: Environment variables

```bash
# Linux/Mac
export MT5_USER=591129415
export MT5_PASSWORD="YourPassword"
export MT5_HOST="mt5.mrpc.pro"
export MT5_CLUSTER="FxPro-MT5 Demo"

# Windows PowerShell
$env:MT5_USER="591129415"
$env:MT5_PASSWORD="YourPassword"
$env:MT5_HOST="mt5.mrpc.pro"
```

See [config.go](../../examples/demos/config/config.go) for details.

## Tips

1. **Start simple** - Uncomment one example at a time
2. **Use Sugar API** - Methods like `BuyMarketWithPips()` and `CalculatePositionSize()` are easier than low-level

3. **Check RetCode** - Always validate trading operations (10009 = success)
4. **Test on demo** - Make sure you're using demo account first!
5. **Read documentation** - [RETURN_CODES_REFERENCE.md](RETURN_CODES_REFERENCE.md) explains all error codes

6. **Use context** - The `ctx` variable is already configured with timeout
7. **Explore types** - Use `go run main.go inspect` to explore protobuf types
8. **Check examples** - See `examples/demos/helpers/` for 17+ working examples

## Common mistakes

❌ **Forgot to check err**
```go
ticket, _ := sugar.BuyMarket("EURUSD", 0.01)  // DON'T ignore errors!
```

✅ **Always check errors**
```go
ticket, err := sugar.BuyMarket("EURUSD", 0.01)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
```

❌ **Not checking RetCode**
```go
result, _ := service.BuyMarket(ctx, "EURUSD", 0.01, 0, 0, 0)
// Assuming it's success!
```

✅ **Check RetCode for trading operations**
```go
result, err := service.BuyMarket(ctx, "EURUSD", 0.01, 0, 0, 0)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
if result.ReturnedCode != 10009 {
    fmt.Printf("Trade failed: %s\n", result.ReturnedCodeDescription)
    return
}
```

## Getting help

- **Protobuf types**: Run `go run main.go inspect` to explore API types
- **Error codes**: See [RETURN_CODES_REFERENCE.md](RETURN_CODES_REFERENCE.md)
- **Examples**: Check `examples/demos/helpers/` for working code
- **API documentation**: See docs folder for complete reference

## Example: Complete trading strategy

Here's a complete example demonstrating proper error handling:

```go
func UserCodeMain() error {
    // Setup
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cfg, err := config.LoadConfig()
    if err != nil {
        return fmt.Errorf("config error: %w", err)
    }

    account, err := mt5.NewMT5Account(ctx, cfg.GrpcServer)
    if err != nil {
        return fmt.Errorf("connection error: %w", err)
    }

    err = account.ConnectEx(ctx, &pb.ConnectExRequest{
        Uuid:      cfg.User,
        Password:  cfg.Password,
        MtCluster: cfg.Cluster,
    })
    if err != nil {
        return fmt.Errorf("login error: %w", err)
    }

    service := mt5service.NewMT5Service(account)
    sugar := mt5sugar.NewMT5Sugar(service)

    // Get account information
    balance, err := sugar.GetBalance()
    if err != nil {
        return fmt.Errorf("balance error: %w", err)
    }
    fmt.Printf("Balance: %.2f\n", balance)

    // Calculate position size based on risk (1% with 50 pip SL)
    volume, err := sugar.CalculatePositionSize("EURUSD", 1.0, 50)
    if err != nil {
        return fmt.Errorf("position size error: %w", err)
    }

    // Open position with calculated volume
    ticket, err := sugar.BuyMarketWithPips("EURUSD", volume, 50, 100)
    if err != nil {
        return fmt.Errorf("order error: %w", err)
    }

    fmt.Printf("✅ Position opened: #%d with %.2f lot\n", ticket, volume)

    return nil
}
```

**Happy trading! 🚀**
