# MT5Service - Mid-Level API Overview

> Convenient wrapper over MT5Account, returning clean Go types instead of protobuf

**🧩 API Layer:** **MID-LEVEL** - balance between low-level (MT5Account) and high-level (MT5Sugar)

---

## ⚠️ Important Note About Code Structure

**MT5Service and MT5Account have very similar code structures** - the difference is primarily in **return types**, not in code logic or complexity.

```go
// Low-Level (MT5Account):
req := &pb.AccountInfoDoubleRequest{PropertyId: ACCOUNT_BALANCE}
data, _ := account.AccountInfoDouble(ctx, req)
balance := data.GetRequestedValue()  // ← unpacking

// Mid-Level (MT5Service):
balance, _ := service.GetAccountDouble(ctx, ACCOUNT_BALANCE)  // ← direct value
```

**Key difference**: MT5Service **automatically unpacks** protobuf `Data` wrappers and **converts** types (`Timestamp` → `time.Time`), but the underlying gRPC calls and logic remain identical.

**This is a matter of coding style preference:**

- Use **MT5Account** if you prefer explicit control over protobuf structures
- Use **MT5Service** if you prefer cleaner Go types and less boilerplate

Both approaches are equally valid - choose what feels more natural for your coding style.

---

## 🎯 Why MT5Service Exists

### Low-Level API (MT5Account) Problem

MT5Account is a direct gRPC client working with protobuf structures:

```go
// Low-level (MT5Account) - verbose:
req := &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
}
data, err := account.AccountInfoDouble(ctx, req)
balance := data.GetRequestedValue()  // ← unpacking!

// Another example:
tickData, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: "EURUSD"})
tick := tickData  // protobuf structure
t := time.Unix(tick.Time, 0)  // ← manual time conversion
```

**Problems:**

- ❌ Need to create protobuf requests manually
- ❌ Responses wrapped in Data structures (`.GetRequestedValue()`, `.Total`)
- ❌ Manual `Timestamp` → `time.Time` conversion
- ❌ Code looks like protobuf, not like Go

---

### Solution: MT5Service

MT5Service wraps MT5Account and returns clean Go types:

```go
// Mid-level (MT5Service) - concise:
balance, err := service.GetAccountDouble(ctx,
    pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE)  // ✅ float64 directly!

// Another example:
tick, _ := service.GetSymbolTick(ctx, "EURUSD")
// tick.Time is already time.Time!
fmt.Printf("Price at %s: %.5f\n", tick.Time.Format("15:04:05"), tick.Bid)
```

**Advantages:**

- ✅ Returns clean Go types (`float64`, `int64`, `string`, `*AccountSummary`)
- ✅ Automatically creates request structures
- ✅ Automatically unpacks Data wrappers
- ✅ Converts `Timestamp` → `time.Time` automatically
- ✅ Code reads like standard Go

---

## 🏗️ Architecture Layers

```
┌─────────────────────────────────────────────
│  HIGH-LEVEL: MT5Sugar                         ← Ready-made patterns (future)
│  (ready strategies, patterns)               
└─────────────────────────────────────────────
                    ↓
┌─────────────────────────────────────────────
│  MID-LEVEL: MT5Service ← YOU ARE HERE         ← Clean Go types
│  (Go types, convenient methods)             
└─────────────────────────────────────────────
                    ↓
┌─────────────────────────────────────────────
│  LOW-LEVEL: MT5Account                        ← Protobuf gRPC
│  (protobuf Request/Data, direct gRPC)       
└─────────────────────────────────────────────
                    ↓
┌─────────────────────────────────────────────
│  MT5 Terminal (via gRPC server)             
└─────────────────────────────────────────────
```

**When to use each layer:**

- **MT5Account (Low-Level)** - when you need full control over protobuf
- **MT5Service (Mid-Level)** - **RECOMMENDED** for most tasks
- **MT5Sugar (High-Level)** - for ready-made patterns

---

## 📦 Creating MT5Service

```go
import (
    pb "github.com/MetaRPC/GoMT5/package"
    "your-project/mt5"
)

// 1. Create MT5Account (low-level)
account, err := mt5.NewMT5Account(
    login,
    password,
    host,
    port,
    grpcServer,
    mtCluster,
)
if err != nil {
    return err
}

// 2. Wrap in MT5Service (mid-level)
service := mt5.NewMT5Service(account)

// 3. Use convenient methods
balance, _ := service.GetAccountDouble(ctx,
    pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE)
fmt.Printf("Balance: %.2f\n", balance)
```

**Access to low-level:**
```go
// If you need direct access to MT5Account:
account := service.GetAccount()
// Now you can call low-level methods
```

---

## 📚 All 37 MT5Service Methods

MT5Service contains **37 methods**, divided into **6 categories**:

### 1️⃣ Account Methods (4 methods)

Getting account information without protobuf ceremonies.

| Method | Returns | Description |
|--------|---------|-------------|
| `GetAccountSummary(ctx)` | `*AccountSummary` | **RECOMMENDED** - all account information |
| `GetAccountDouble(ctx, propertyID)` | `float64` | Single double property (Balance, Equity, Margin) |
| `GetAccountInteger(ctx, propertyID)` | `int64` | Single integer property (Login, Leverage) |
| `GetAccountString(ctx, propertyID)` | `string` | Single string property (Currency, Company) |

📖 **[Full documentation: Account Methods](./1.%20Account_Methods.md)**

---

### 2️⃣ Symbol Methods (13 methods)

Working with symbols, parameters, quotes.

| Method | Returns | Description |
|--------|---------|-------------|
| `GetSymbolsTotal(ctx, selectedOnly)` | `int32` | Number of available symbols |
| `SymbolExist(ctx, symbol)` | `(bool, bool)` | Check symbol existence |
| `GetSymbolName(ctx, index, selectedOnly)` | `string` | Symbol name by index |
| `SymbolSelect(ctx, symbol, select_)` | `bool` | Add/remove symbol in Market Watch |
| `IsSymbolSynchronized(ctx, symbol)` | `bool` | Check synchronization |
| `GetSymbolDouble(ctx, symbol, property)` | `float64` | Single double property (Bid, Ask, Point) |
| `GetSymbolInteger(ctx, symbol, property)` | `int64` | Single integer property (Digits, Spread) |
| `GetSymbolString(ctx, symbol, property)` | `string` | Single string property (Description) |
| `GetSymbolMarginRate(ctx, symbol, orderType)` | `*SymbolMarginRate` | Margin rates |
| `GetSymbolTick(ctx, symbol)` | `*SymbolTick` | Last tick (with time.Time) |
| `GetSymbolSessionQuote(ctx, symbol, day, idx)` | `*SessionTime` | Quote session time |
| `GetSymbolSessionTrade(ctx, symbol, day, idx)` | `*SessionTime` | Trade session time |
| `GetSymbolParamsMany(ctx, name, sort, page, perPage)` | `[]SymbolParams` | **RECOMMENDED** - parameters of multiple symbols |

📖 **[Full documentation: Symbol Methods](./2.%20Symbol_Methods.md)**

---

### 3️⃣ Position & Orders Methods (5 methods)

Getting information about open positions, pending orders, history.

| Method | Returns | Description |
|--------|---------|-------------|
| `GetPositionsTotal(ctx)` | `int32` | Number of open positions |
| `GetOpenedOrders(ctx, sortMode)` | `*pb.OpenedOrdersData` | Open positions and pending orders |
| `GetOpenedTickets(ctx)` | `([]int64, []int64)` | **Lightweight** - only ticket numbers |
| `GetOrderHistory(ctx, from, to, sort, page, perPage)` | `*pb.OrdersHistoryData` | Orders and deals history |
| `GetPositionsHistory(ctx, sort, from, to, page, perPage)` | `*pb.PositionsHistoryData` | Closed positions history with P&L |

📖 **[Full documentation: Position & Orders Methods](./3.%20Position_Orders_Methods.md)**

---

### 4️⃣ Market Depth Methods (3 methods)

Working with Depth of Market (DOM) / "order book".

| Method | Returns | Description |
|--------|---------|-------------|
| `SubscribeMarketDepth(ctx, symbol)` | `bool` | Subscribe to DOM updates |
| `GetMarketDepth(ctx, symbol)` | `[]BookInfo` | Current DOM snapshot |
| `UnsubscribeMarketDepth(ctx, symbol)` | `bool` | Unsubscribe from DOM |

📖 **[Full documentation: Market Depth Methods](./5.%20MarketDepth_Methods.md)**

---

### 5️⃣ Trading Methods (6 methods)

Trading operations: placing, modifying, closing orders, calculations.

| Method | Returns | Description |
|--------|---------|-------------|
| `PlaceOrder(ctx, req)` | `*OrderResult` | Send market/pending order |
| `ModifyOrder(ctx, req)` | `*OrderResult` | Modify order/position (SL/TP) |
| `CloseOrder(ctx, req)` | `uint32` | Close position/delete order |
| `CheckOrder(ctx, req)` | `*OrderCheckResult` | Validate order |
| `CalculateMargin(ctx, req)` | `float64` | Calculate required margin |
| `CalculateProfit(ctx, req)` | `float64` | Calculate potential profit |

📖 **[Full documentation: Trading Methods](./4.%20Trading_Methods.md)**

---

### 6️⃣ Streaming Methods (5 methods)

Real-time data streams: ticks, trade events, profits.

| Method | Returns | Description |
|--------|---------|-------------|
| `StreamTicks(ctx, symbols)` | `(<-chan *SymbolTick, <-chan error)` | Real-time ticks (with time.Time) |
| `StreamTradeUpdates(ctx)` | `(<-chan *pb.OnTradeData, <-chan error)` | Trade events (new/closed positions) |
| `StreamPositionProfits(ctx)` | `(<-chan *pb.OnPositionProfitData, <-chan error)` | Real-time P&L updates |
| `StreamOpenedTickets(ctx)` | `(<-chan *pb.OnPositionsAndPendingOrdersTicketsData, <-chan error)` | **Lightweight** - only tickets |
| `StreamTransactions(ctx)` | `(<-chan *pb.OnTradeTransactionData, <-chan error)` | **Detailed** - all transactions |

📖 **[Full documentation: Streaming Methods](./6.%20Streaming_Methods.md)**

---

## 📦 DTO Structures (Data Transfer Objects)

MT5Service uses clean Go structures instead of protobuf:

### AccountSummary

```go
type AccountSummary struct {
    Login                   int64                        // Account number
    Balance                 float64                      // Balance
    Equity                  float64                      // Equity (Balance + Floating P&L)
    UserName                string                       // Client name
    Leverage                int64                        // Leverage (100 = 1:100)
    TradeMode               pb.MrpcEnumAccountTradeMode  // Account type (demo/real)
    CompanyName             string                       // Broker name
    Currency                string                       // Deposit currency
    ServerTime              *time.Time                   // Server time (already time.Time!)
    UtcTimezoneShiftMinutes int64                        // UTC timezone shift
    Credit                  float64                      // Credit
}
```

### SymbolParams

```go
type SymbolParams struct {
    Name                 string  // Symbol name
    Bid                  float64 // Current Bid price
    Ask                  float64 // Current Ask price
    Last                 float64 // Last deal price
    Point                float64 // Point size
    Digits               int32   // Decimal places
    Spread               int32   // Spread in points
    VolumeMin            float64 // Min volume
    VolumeMax            float64 // Max volume
    VolumeStep           float64 // Volume step
    TradeTickSize        float64 // Trade tick size
    TradeTickValue       float64 // Trade tick value
    TradeContractSize    float64 // Contract size
    SwapLong             float64 // Swap for long
    SwapShort            float64 // Swap for short
    MarginInitial        float64 // Initial margin
    MarginMaintenance    float64 // Maintenance margin
}
```

### SymbolTick

```go
type SymbolTick struct {
    Time       time.Time // Tick time (already time.Time!)
    Bid        float64   // Bid price
    Ask        float64   // Ask price
    Last       float64   // Last deal price
    Volume     uint64    // Volume
    TimeMS     int64     // Time in milliseconds
    Flags      uint32    // Tick flags
    VolumeReal float64   // Volume with decimals
}
```

### OrderResult

```go
type OrderResult struct {
    ReturnedCode    uint32  // Return code (10009 = success)
    Deal            uint64  // Deal ticket
    Order           uint64  // Order ticket
    Volume          float64 // Executed volume
    Price           float64 // Execution price
    Bid             float64 // Current Bid
    Ask             float64 // Current Ask
    Comment         string  // Broker comment
    RequestID       uint32  // Request ID
    RetCodeExternal int32   // External return code
}
```

### OrderCheckResult

```go
type OrderCheckResult struct {
    ReturnedCode uint32  // Validation code (0 = success)
    Balance      float64 // Balance after deal
    Equity       float64 // Equity after deal
    Profit       float64 // Profit
    Margin       float64 // Required margin
    MarginFree   float64 // Free margin after
    MarginLevel  float64 // Margin level after (%)
    Comment      string  // Error description
}
```

### BookInfo

```go
type BookInfo struct {
    Type       pb.BookType // BID or ASK
    Price      float64     // Price level
    Volume     int64       // Volume (integer)
    VolumeReal float64     // Volume (decimal)
}
```

---

## 💡 Usage Examples

### Example 1: Getting Balance and Equity

```go
// ❌ WAS (Low-Level MT5Account):
balanceReq := &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
}
balanceData, _ := account.AccountInfoDouble(ctx, balanceReq)
balance := balanceData.GetRequestedValue()

equityReq := &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_EQUITY,
}
equityData, _ := account.AccountInfoDouble(ctx, equityReq)
equity := equityData.GetRequestedValue()

// ✅ BECAME (Mid-Level MT5Service):
balance, _ := service.GetAccountDouble(ctx,
    pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE)
equity, _ := service.GetAccountDouble(ctx,
    pb.AccountInfoDoublePropertyType_ACCOUNT_EQUITY)

// ✅✅ EVEN BETTER - use GetAccountSummary:
summary, _ := service.GetAccountSummary(ctx)
fmt.Printf("Balance: %.2f, Equity: %.2f\n",
    summary.Balance, summary.Equity)
```

---

### Example 2: Getting Current Price

```go
// ❌ WAS (Low-Level):
tickData, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
    Symbol: "EURUSD",
})
bid := tickData.Bid
ask := tickData.Ask
t := time.Unix(tickData.Time, 0)  // ← manual conversion

// ✅ BECAME (Mid-Level):
tick, _ := service.GetSymbolTick(ctx, "EURUSD")
fmt.Printf("Time: %s, Bid: %.5f, Ask: %.5f\n",
    tick.Time.Format("15:04:05"),  // ← already time.Time!
    tick.Bid, tick.Ask)
```

---

### Example 3: Getting Symbol Parameters

```go
// ❌ WAS (Low-Level) - 10+ calls for all properties:
bidData, _ := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
    Symbol: "EURUSD",
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_BID,
})
bid := bidData.Value

askData, _ := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
    Symbol: "EURUSD",
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_ASK,
})
ask := askData.Value
// ... 8 more calls for other properties

// ✅ BECAME (Mid-Level) - 1 call for all properties:
symbolName := "EURUSD"
symbols, _, _ := service.GetSymbolParamsMany(ctx, &symbolName, nil, nil, nil)
if len(symbols) > 0 {
    s := symbols[0]
    fmt.Printf("Symbol: %s\n", s.Name)
    fmt.Printf("Bid: %.5f, Ask: %.5f\n", s.Bid, s.Ask)
    fmt.Printf("Digits: %d, Spread: %d\n", s.Digits, s.Spread)
    fmt.Printf("Volume: %.2f - %.2f\n", s.VolumeMin, s.VolumeMax)
}
```

---

### Example 4: Validating and Placing Order

```go
// ✅ Mid-Level - clear and safe code
func SafePlaceOrder(service *mt5.MT5Service) error {
    ctx := context.Background()

    // 1. Validate order
    checkReq := &pb.OrderCheckRequest{
        Symbol: "EURUSD",
        Action: pb.ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
        Type:   pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY,
        Volume: 1.0,
    }

    checkResult, err := service.CheckOrder(ctx, checkReq)
    if err != nil {
        return err
    }

    if checkResult.ReturnedCode != 0 {
        return fmt.Errorf("invalid order: %s", checkResult.Comment)
    }

    fmt.Printf("Order valid:\n")
    fmt.Printf("  Margin required: %.2f\n", checkResult.Margin)
    fmt.Printf("  Free margin after: %.2f\n", checkResult.MarginFree)

    // 2. Place order
    orderReq := &pb.OrderSendRequest{
        Symbol:    "EURUSD",
        Action:    pb.ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
        Type:      pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY,
        Volume:    1.0,
        Price:     0,
        Deviation: 10,
    }

    result, err := service.PlaceOrder(ctx, orderReq)
    if err != nil {
        return err
    }

    if result.ReturnedCode != 10009 {
        return fmt.Errorf("order rejected: %s", result.Comment)
    }

    fmt.Printf("Order placed:\n")
    fmt.Printf("  Order: %d, Deal: %d\n", result.Order, result.Deal)
    fmt.Printf("  Price: %.5f\n", result.Price)

    return nil
}
```

---

### Example 5: Real-time Price and Profit Monitoring

```go
// ✅ Mid-Level - simple streaming
func MonitorPriceAndProfit(service *mt5.MT5Service) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Real-time ticks
    tickCh, tickErrCh := service.StreamTicks(ctx, []string{"EURUSD"})

    // Real-time profit
    profitCh, profitErrCh := service.StreamPositionProfits(ctx)

    for {
        select {
        case tick, ok := <-tickCh:
            if !ok {
                return
            }
            fmt.Printf("[%s] EURUSD: %.5f\n",
                tick.Time.Format("15:04:05"), tick.Bid)

        case data, ok := <-profitCh:
            if !ok {
                return
            }
            totalProfit := 0.0
            for _, pos := range data.OpenedPositions {
                totalProfit += pos.Profit
            }
            fmt.Printf("📊 Total P&L: %.2f\n", totalProfit)

        case err := <-tickErrCh:
            fmt.Printf("❌ Tick error: %v\n", err)
            return

        case err := <-profitErrCh:
            fmt.Printf("❌ Profit error: %v\n", err)
            return
        }
    }
}
```

---

## 📊 Low-Level vs Mid-Level Comparison

| Aspect | Low-Level (MT5Account) | Mid-Level (MT5Service) |
|--------|----------------------|---------------------|
| **Return type** | Protobuf Data structures | Clean Go types |
| **Request creation** | Manual | Automatic |
| **Data unpacking** | `.GetRequestedValue()`, `.Total` | Direct |
| **Time conversion** | Manual `.AsTime()` | Automatic |
| **Code amount** | More (verbose) | **Less (concise)** |
| **Readability** | Protobuf style | **Go style** |
| **Control** | **Full** | Medium |
| **Recommended for** | Specific tasks | **Most tasks** |

**Conclusion:** Use MT5Service for 90% of tasks, MT5Account - for specific cases.

---

## 💡 Usage Recommendations

###  When to Use MT5Service (RECOMMENDED)

- ✅ Getting account, symbol, position data
- ✅ Trading operations (PlaceOrder, ModifyOrder)
- ✅ Real-time monitoring (StreamTicks, StreamPositionProfits)
- ✅ Calculations (CalculateMargin, CalculateProfit)
- ✅ Working with DOM (Market Depth)
- ✅ Trading history
- ✅ 90% of all tasks

### ⚠️ When to Use MT5Account (Low-Level)

- Need full control over protobuf
- Specific handling of Data structures
- Performance optimization (rare)
- Access to methods not wrapped in MT5Service

**Access to Low-Level:**

```go
account := service.GetAccount()
// Now you can use low-level methods
```

---

## 🔧 Best Practices

### 1. Use "RECOMMENDED" Methods

```go
// ✅ GOOD - use GetAccountSummary for multiple properties
summary, _ := service.GetAccountSummary(ctx)
fmt.Printf("Balance: %.2f, Equity: %.2f\n", summary.Balance, summary.Equity)

// ❌ BAD - don't make 2 calls for 2 properties
balance, _ := service.GetAccountDouble(ctx, ACCOUNT_BALANCE)
equity, _ := service.GetAccountDouble(ctx, ACCOUNT_EQUITY)
```

```go
// ✅ GOOD - GetSymbolParamsMany for all parameters
symbolName := "EURUSD"
symbols, _, _ := service.GetSymbolParamsMany(ctx, &symbolName, nil, nil, nil)
s := symbols[0]  // All properties in one structure

// ❌ BAD - don't make 10 calls GetSymbolDouble/Integer
bid, _ := service.GetSymbolDouble(ctx, "EURUSD", SYMBOL_BID)
ask, _ := service.GetSymbolDouble(ctx, "EURUSD", SYMBOL_ASK)
// ... 8 more calls
```

---

### 2. Handle Errors

```go
// ✅ GOOD - always check errors
balance, err := service.GetAccountDouble(ctx, ACCOUNT_BALANCE)
if err != nil {
    return fmt.Errorf("failed to get balance: %w", err)
}

// ❌ BAD - don't ignore errors
balance, _ := service.GetAccountDouble(ctx, ACCOUNT_BALANCE)
```

---

### 3. Use Context with Timeout

```go
// ✅ GOOD - use timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

balance, err := service.GetAccountDouble(ctx, ACCOUNT_BALANCE)

// ❌ BAD - context.Background() without timeout
ctx := context.Background()
```

---

### 4. Check Return Codes in Trading

```go
// ✅ GOOD - check ReturnedCode
result, _ := service.PlaceOrder(ctx, orderReq)
if result.ReturnedCode != 10009 {  // TRADE_RETCODE_DONE
    return fmt.Errorf("order failed: %s", result.Comment)
}

// ❌ BAD - assume success
result, _ := service.PlaceOrder(ctx, orderReq)
fmt.Printf("Order: %d\n", result.Order)  // may be 0!
```

---

### 5. Use Defer for Cleanup

```go
// ✅ GOOD - defer for unsubscribe
success, _ := service.SubscribeMarketDepth(ctx, "EURUSD")
if success {
    defer service.UnsubscribeMarketDepth(ctx, "EURUSD")
}

// Work with DOM...

// ❌ BAD - forget to unsubscribe (resource leak)
service.SubscribeMarketDepth(ctx, "EURUSD")
// ... work ...
// forgot UnsubscribeMarketDepth!
```

---

## 📚 Related Sections

### Documentation by Method Categories:

- 📖 [Account Methods](./1.%20Account_Methods.md) - account operations (4 methods)
- 📖 [Symbol Methods](./2.%20Symbol_Methods.md) - symbol operations (13 methods)
- 📖 [Position & Orders Methods](./3.%20Position_Orders_Methods.md) - positions and orders (5 methods)
- 📖 [Market Depth Methods](./5.%20MarketDepth_Methods.md) - order book (3 methods)
- 📖 [Trading Methods](./4.%20Trading_Methods.md) - trading operations (6 methods)
- 📖 [Streaming Methods](./6.%20Streaming_Methods.md) - real-time streams (5 methods)

---

## 🎯 Summary

**MT5Service solves the main task** - removes protobuf ceremonies and returns clean Go:

### Was (Low-Level):
```go
req := &pb.AccountInfoDoubleRequest{PropertyId: ACCOUNT_BALANCE}
data, _ := account.AccountInfoDouble(ctx, req)
balance := data.GetRequestedValue()

tickData, _ := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{Symbol: "EURUSD"})
t := time.Unix(tickData.Time, 0)
```

### Became (Mid-Level):
```go
balance, _ := service.GetAccountDouble(ctx, ACCOUNT_BALANCE)

tick, _ := service.GetSymbolTick(ctx, "EURUSD")
// tick.Time is already time.Time!
```

### MT5Service Advantages:

- ✅ **Clean Go types** (`float64`, `int64`, `string`, `time.Time`)
- ✅ **Less code** (30-70% shorter)
- ✅ **Automatic unpacking** of Data structures
- ✅ **Automatic conversion** of time
- ✅ **Convenient DTOs** (`AccountSummary`, `SymbolParams`, `OrderResult`)
- ✅ **Code reads** like standard Go

---

## 🎉 Good Luck!

You now have all the tools to work efficiently with the MT5 trading platform through clean Go code. Whether you choose MT5Service for convenience or MT5Account for control, both approaches will serve you well.

**Happy trading and may your algorithms be profitable!** 📈✨
