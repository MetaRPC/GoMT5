# User Code Sandbox Guide

## What is this?

This is **your sandbox** for writing custom MT5 trading code in Go. Connection setup is already done - you only need to add your trading logic!

## Quick Start

1. **Case already active in main.go**
   - Open `examples/demos/main.go`
   - Case "18" for user code sandbox is ready to use (the menu label "(DISABLED)" and the "uncomment" notes in the
     file headers are outdated - nothing needs to be uncommented)

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
go run main.go custom
```

## What's already configured

✓ **Connection** - MT5 terminal connected via gRPC (`sugar.QuickConnect(cfg.MtCluster)`)

✓ **Configuration** - Loaded from `config/config.json` or environment variables (`config.LoadConfig()`; API key from `api_key`, `MRPC_API_KEY` or `--api-key`)

✓ **MT5Account** - Low-level gRPC client (variable: `account`)

✓ **MT5Service** - Mid-level wrapper (variable: `service`)

✓ **MT5Sugar** - High-level Sugar API (variable: `sugar`)

✓ **Context** - `context.Background()` (variable: `ctx`) - no timeout; wrap it with `context.WithTimeout` where you need one

✓ **Imports** - `context`, `fmt`, `strings`, `config`, `helpers` (demo helpers), `mt5` (`examples/mt5`) and `pb` - add more as you need them

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

// With SL/TP as distances in points (200 = 20 pips on 5-digit EURUSD)
ticket, err = sugar.BuyMarketWithPips("EURUSD", 0.01, 200, 300)

// Calculate position size based on risk (2% with SL 500 points away)
lotSize, _ := sugar.CalculatePositionSize("EURUSD", 2.0, 500)
ticket, err = sugar.BuyMarketWithPips("EURUSD", lotSize, 500, 1000)
```

Volumes and prices are sent as given (no rounding) - see [Normalization](GLOSSARY.md#normalization).

### Place pending order

```go
info, _ := sugar.GetSymbolInfo("EURUSD") // Bid, Ask, Point, Digits, ...

// Buy Limit 20 pips (200 points) below current Ask
ticket, err := sugar.BuyLimit("EURUSD", 0.01, info.Ask-200*info.Point)

// Buy Stop 20 pips above current Ask
ticket, err = sugar.BuyStop("EURUSD", 0.01, info.Ask+200*info.Point)

// Sell Limit with absolute price, SL and TP
ticket, err = sugar.SellLimitWithSLTP("EURUSD", 0.01, 1.1050, 1.1070, 1.1030)

// Delete a pending order
_, err = service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: ticket})
```

There are no pip-based pending methods; `sugar.CalculateSLTP()` turns distances into SL/TP prices for `BuyLimitWithSLTP` / `SellLimitWithSLTP`.

### Get open positions

```go
// Sugar API
positions, err := sugar.GetOpenPositions()
fmt.Printf("Open positions: %d\n", len(positions))
for _, pos := range positions {
    fmt.Printf("  #%d %s %.2f lot, Profit: %.2f\n",
        pos.Ticket, pos.Symbol, pos.Volume, pos.Profit)
}

// Sugar API filtered by symbol
eurPositions, err := sugar.GetPositionsBySymbol("EURUSD")

// Service API: positions and pending orders, sorted
opened, err := service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
fmt.Printf("Positions: %d, pending orders: %d\n", len(opened.PositionInfos), len(opened.OpenedOrders))
```

### Close positions

```go
// Close specific position
err := sugar.ClosePosition(ticket)

// Close all positions for symbol (returns how many were closed)
closed, err := sugar.CloseAllBySymbol("EURUSD")

// Close all positions (all symbols)
closed, err = sugar.CloseAllPositions()
```

### Modify position

```go
// Modify SL/TP
err := sugar.ModifyPositionSLTP(ticket, 1.0850, 1.0950)

// Modify only SL or only TP
err = sugar.ModifyPositionSL(ticket, 1.0850)
err = sugar.ModifyPositionTP(ticket, 1.0950)
```

There is no trailing-stop method: move the SL yourself with `ModifyPositionSL`, or run the Trailing Stop
orchestrator (`go run main.go trailing`, source `orchestrators/11_trailing_stop.go`).

### Calculate position size

```go
// Calculate position size based on risk (SL distance in points)
volume, err := sugar.CalculatePositionSize("EURUSD", 2.0, 200)
fmt.Printf("Volume for 2%% risk with 200 point (20 pip) SL: %.2f lot\n", volume)

// Then use calculated volume
ticket, err := sugar.BuyMarket("EURUSD", volume)
```

## Return Codes (RetCodes)

**Always check RetCode after trading operations!**

```go
result, err := service.PlaceOrder(ctx, &pb.OrderSendRequest{
    Symbol:    "EURUSD",
    Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
    Volume:    0.01,
})
if err != nil {
    fmt.Printf("Error: %v\n", err) // gRPC error, or the refused request (ApiError with the RetCode)
    return
}

// Check RetCode
if result.ReturnedCode == 10009 {  // Request completed
    fmt.Printf("✅ Order opened: #%d\n", result.Order)
} else if result.ReturnedCode == 10008 {  // Order placed
    fmt.Printf("✅ Pending order placed: #%d\n", result.Order)
} else {
    fmt.Printf("❌ Order failed: %s (code %d)\n",
        result.Comment, result.ReturnedCode)
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

### 1. Returned errors (transport errors and refused requests)
```go
result, err := service.PlaceOrder(ctx, req)
if err != nil {
    // A network/gRPC error, or an *ApiError from package/Helpers when the terminal refused the
    // request - its MqlErrorTradeIntCode() carries the RetCode (e.g. 10014 invalid volume)
    fmt.Printf("Error: %v\n", err)
    return
}
```

### 2. Trading operation errors (RetCode)

```go
if !helpers.CheckRetCode(result.ReturnedCode, "Order placement") { // prints the code, its description and a hint
    fmt.Printf("Trade rejected: %s\n", result.Comment)
}
```

**Always check BOTH!**

## Documentation

- [MT5Account Master Overview](../MT5Account/MT5Account.Master.Overview.md) - Complete API reference (40+ methods)
- [MT5Service API Overview](../MT5Service/MT5Service.Overview.md) - Mid-level wrapper (36 methods)
- [MT5Sugar API Overview](../MT5Sugar/MT5Sugar.API_Overview.md) - High-level Sugar API (60+ methods)
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
  "mt_cluster": "MetaQuotes-Demo",
  "api_key": "",
  "test_symbol": "EURUSD",
  "test_volume": 0.01
}
```

With `user` 0 or an empty `password` (as shipped) a demo account on `MetaQuotes-Demo` is opened automatically.

### Method 2: Environment variables

```bash
# Linux/Mac
export MT5_USER=591129415
export MT5_PASSWORD="YourPassword"
export MT5_HOST="mt5.mrpc.pro"
export MT5_CLUSTER="MetaQuotes-Demo"
export MRPC_API_KEY="your-api-key"

# Windows PowerShell
$env:MT5_USER="591129415"
$env:MT5_PASSWORD="YourPassword"
$env:MT5_HOST="mt5.mrpc.pro"
```

See `examples/demos/config/config.go` for details.

## Tips

1. **Start simple** - Uncomment one example at a time
2. **Use Sugar API** - Methods like `BuyMarketWithPips()` and `CalculatePositionSize()` are easier than low-level (their "pips" are points - see [Points vs Pips](GLOSSARY.md#points-vs-pips))

3. **Check RetCode** - Always validate trading operations (10009 = success)
4. **Test on demo** - Make sure you're using demo account first!
5. **Read documentation** - [RETURN_CODES_REFERENCE.md](RETURN_CODES_REFERENCE.md) explains all error codes

6. **Use context** - The `ctx` variable is `context.Background()`; add `context.WithTimeout` for long calls (Sugar methods have their own timeouts)
7. **Explore types** - Use `go run main.go inspect` to explore protobuf types
8. **Check examples** - See `examples/demos/lowlevel/`, `service/`, `sugar/`, `orchestrators/` and `presets/` for the 17 working demos (menu items 1-17)

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
result, _ := service.PlaceOrder(ctx, req)
// Assuming it's success!
```

✅ **Check RetCode for trading operations**
```go
result, err := service.PlaceOrder(ctx, req)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
if result.ReturnedCode != 10009 {
    fmt.Printf("Trade failed: %s (code %d)\n", result.Comment, result.ReturnedCode)
    return
}
```

## Getting help

- **Protobuf types**: Run `go run main.go inspect` to explore API types
- **Error codes**: See [RETURN_CODES_REFERENCE.md](RETURN_CODES_REFERENCE.md)
- **Examples**: Check `examples/demos/lowlevel/`, `service/`, `sugar/` and `orchestrators/` for working code
- **API documentation**: See docs folder for complete reference

## Example: Complete trading strategy

Here's a complete example demonstrating proper error handling:

```go
func UserCodeMain() error {
    // Setup (the same steps RunUserCode() performs)
    cfg, err := config.LoadConfig()
    if err != nil {
        return fmt.Errorf("config error: %w", err)
    }

    sugar, err := mt5.NewMT5Sugar(cfg.User, cfg.Password, cfg.GrpcServer)
    if err != nil {
        return fmt.Errorf("setup error: %w", err)
    }
    sugar.GetAccount().ApiKey = cfg.ApiKey

    if err := sugar.QuickConnect(cfg.MtCluster); err != nil {
        return fmt.Errorf("login error: %w", err)
    }
    defer sugar.Disconnect()

    // sugar.GetService() and sugar.GetAccount() give the lower levels if you need them

    // Get account information
    balance, err := sugar.GetBalance()
    if err != nil {
        return fmt.Errorf("balance error: %w", err)
    }
    fmt.Printf("Balance: %.2f\n", balance)

    // Calculate position size based on risk (1% with SL 500 points = 50 pips away)
    volume, err := sugar.CalculatePositionSize("EURUSD", 1.0, 500)
    if err != nil {
        return fmt.Errorf("position size error: %w", err)
    }

    // Open position with calculated volume
    ticket, err := sugar.BuyMarketWithPips("EURUSD", volume, 500, 1000)
    if err != nil {
        return fmt.Errorf("order error: %w", err)
    }

    fmt.Printf("✅ Position opened: #%d with %.2f lot\n", ticket, volume)

    return nil
}
```

**Happy trading! 🚀**
