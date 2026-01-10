# MT5Account - Master Overview

> One page to **orient fast**: what lives where, how to choose the right API, and jump links to every **overview** and **method spec** in this docs set.

---

## 🚦 Start here - Section Overviews

* **[Account\_Information - Overview](./1.%20Account_information/Account_Information.Overview.md)**
  Account balance/equity/margin/leverage, complete snapshot or single properties.

* **[Symbol\_Information - Overview](./2.%20Symbol_information/Symbol_Information.Overview.md)**
  Quotes, symbol properties, trading rules, Market Watch management.

* **[Position\_Orders\_Information - Overview](./3.%20Position_Orders_Information/Position_Orders_Information.Overview.md)**
  Open positions, pending orders, historical deals, order history.

* **[Trading\_Operations - Overview](./4.%20Trading_Operations/Trading_Operations.Overview.md)**
  Order execution, position management, margin calculations, trade validation.

* **[Market\_Depth\_DOM - Overview](./5.%20Market_Depth(DOM)/Market_Depth.Overview.md)**
  Level II quotes, order book data, market depth subscription.

* **[Additional\_Methods - Overview](./6.%20Additional_Methods/Additional_Methods.Overview.md)**
  Advanced symbol info, trading sessions, margin rates, batch operations.

* **[Streaming\_Methods - Overview](./7.%20Streaming_Methods/Streaming_Methods.Overview.md)**
  Real-time streams: ticks, trades, profit updates, transaction log.

---

## 🧭 How to pick an API

| If you need…                   | Go to…                      | Typical calls                                                                 |
| ------------------------------ | --------------------------- | ----------------------------------------------------------------------------- |
| Account snapshot               | Account\_Information        | `AccountSummary`, `AccountInfoDouble`, `AccountInfoInteger`                   |
| Quotes & symbol properties     | Symbol\_Information         | `SymbolInfoTick`, `SymbolInfoDouble`, `SymbolsTotal`                          |
| Current positions & orders     | Position\_Orders\_Information | `PositionsTotal`, `OpenedOrders`, `OpenedOrdersTickets`                     |
| Historical trades              | Position\_Orders\_Information | `OrderHistory`, `PositionsHistory`                                          |
| Level II / Order book          | Market\_Depth\_DOM          | `MarketBookAdd`, `MarketBookGet`, `MarketBookRelease`                         |
| Trading operations             | Trading\_Operations         | `OrderSend`, `OrderModify`, `OrderClose`                                      |
| Pre-trade calculations         | Trading\_Operations         | `OrderCalcMargin`, `OrderCheck`                                               |
| Advanced symbol data           | Additional\_Methods         | `SymbolParamsMany`, `SymbolInfoSessionTrade`                                  |
| Real-time updates              | Streaming\_Methods          | `OnSymbolTick`, `OnTrade`, `OnPositionProfit`                                 |

---

## 🔌 Usage pattern (gRPC protocol)

Every method follows gRPC client-server pattern:

```go
// Create client connection
conn, err := grpc.Dial("localhost:8002", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Create MT5Account client
client := pb.NewAccountHelperClient(conn)

// Make request with context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Call method
result, err := client.AccountSummary(ctx, &pb.AccountSummaryRequest{})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

// Use result
fmt.Printf("Balance: $%.2f\n", result.AccountBalance)
```

---

Every method follows the same shape:

* **Proto Service/Method:** `Service.Method(Request) → Reply`
* **Go client:** `client.Method(ctx, &Request{...})`
* **Reply structure:** Contains either `Data` payload or `Error`
* **Return codes:** Trading operations return status codes (10009 = success)

**Timestamps:** Use `google.protobuf.Timestamp` (convert with `timestamppb.New(time)` and `.AsTime()`).

**Streaming methods:** Return `stream` object, use `stream.Recv()` in loop to receive updates.

---

# 📚 Full Index · All Method Specs

---

## 📄 Account Information

* **Overview:** [Account\_Information.Overview.md](./1.%20Account_information/Account_Information.Overview.md)

### Complete Snapshot

* [AccountSummary.md](./1.%20Account_information/AccountSummary.md) - All account info at once (balance, equity, margin, etc.)

### Individual Properties

* [AccountInfoDouble.md](./1.%20Account_information/AccountInfoDouble.md) - Single double value (balance, equity, margin, profit, etc.)
* [AccountInfoInteger.md](./1.%20Account_information/AccountInfoInteger.md) - Single integer value (login, leverage, limit orders, etc.)
* [AccountInfoString.md](./1.%20Account_information/AccountInfoString.md) - Single string value (name, server, currency, company)

---

## 📊 Symbol Information

* **Overview:** [Symbol\_Information.Overview.md](./2.%20Symbol_information/Symbol_Information.Overview.md)

### Current Quotes

* [SymbolInfoTick.md](./2.%20Symbol_information/SymbolInfoTick.md) - Current quote for symbol (bid, ask, last, volume, time)

### Symbol Inventory & Management

* [SymbolsTotal.md](./2.%20Symbol_information/SymbolsTotal.md) - Count of available symbols
* [SymbolName.md](./2.%20Symbol_information/SymbolName.md) - Get symbol name by index
* [SymbolSelect.md](./2.%20Symbol_information/SymbolSelect.md) - Enable/disable symbol in Market Watch
* [SymbolExist.md](./2.%20Symbol_information/SymbolExist.md) - Check if symbol exists
* [SymbolIsSynchronized.md](./2.%20Symbol_information/SymbolIsSynchronized.md) - Check symbol data sync status

### Symbol Properties

* [SymbolInfoDouble.md](./2.%20Symbol_information/SymbolInfoDouble.md) - Single double property (bid, ask, point, volume min/max, etc.)
* [SymbolInfoInteger.md](./2.%20Symbol_information/SymbolInfoInteger.md) - Single integer property (digits, spread, stops level, etc.)
* [SymbolInfoString.md](./2.%20Symbol_information/SymbolInfoString.md) - Single string property (description, currency, path)

---

## 📦 Positions & Orders

* **Overview:** [Position\_Orders\_Information.Overview.md](./3.%20Position_Orders_Information/Position_Orders_Information.Overview.md)

### Current State

* [PositionsTotal.md](./3.%20Position_Orders_Information/PositionsTotal.md) - Count of open positions
* [OpenedOrders.md](./3.%20Position_Orders_Information/OpenedOrders.md) - Full details of all open positions and pending orders
* [OpenedOrdersTickets.md](./3.%20Position_Orders_Information/OpenedOrdersTickets.md) - Ticket numbers only (lightweight)

### Historical Data

* [OrderHistory.md](./3.%20Position_Orders_Information/OrderHistory.md) - Historical orders within time range (with pagination)
* [PositionsHistory.md](./3.%20Position_Orders_Information/PositionsHistory.md) - Closed positions within time range (with pagination)

---

## 🛠 Trading Actions

* **Overview:** [Trading\_Operations.Overview.md](./4.%20Trading_Operations/Trading_Operations.Overview.md)

### Order Execution & Management

* [OrderSend.md](./4.%20Trading_Operations/OrderSend.md) - Place market or pending orders
* [OrderModify.md](./4.%20Trading_Operations/OrderModify.md) - Modify SL/TP or order parameters
* [OrderClose.md](./4.%20Trading_Operations/OrderClose.md) - Close positions (full or partial)

### Pre-Trade Calculations

* [OrderCalcMargin.md](./4.%20Trading_Operations/OrderCalcMargin.md) - Calculate margin required for trade
* [OrderCalcProfit.md](./4.%20Trading_Operations/OrderCalcProfit.md) - Calculate potential profit
* [OrderCheck.md](./4.%20Trading_Operations/OrderCheck.md) - Validate trade request before execution

---

## 📈 Market Depth (DOM)

* **Overview:** [Market\_Depth.Overview.md](./5.%20Market_Depth(DOM)/Market_Depth.Overview.md)

### Level II Quotes

* [MarketBookAdd.md](./5.%20Market_Depth(DOM)/MarketBookAdd.md) - Subscribe to Market Depth for symbol
* [MarketBookGet.md](./5.%20Market_Depth(DOM)/MarketBookGet.md) - Get current order book data
* [MarketBookRelease.md](./5.%20Market_Depth(DOM)/MarketBookRelease.md) - Unsubscribe from Market Depth

---

## 🔧 Additional Methods

* **Overview:** [Additional\_Methods.Overview.md](./6.%20Additional_Methods/Additional_Methods.Overview.md)

### Advanced Symbol Information

* [SymbolInfoMarginRate.md](./6.%20Additional_Methods/SymbolInfoMarginRate.md) - Margin rates for symbol and order type
* [SymbolInfoSessionQuote.md](./6.%20Additional_Methods/SymbolInfoSessionQuote.md) - Quote session times
* [SymbolInfoSessionTrade.md](./6.%20Additional_Methods/SymbolInfoSessionTrade.md) - Trade session times
* [SymbolParamsMany.md](./6.%20Additional_Methods/SymbolParamsMany.md) - Detailed parameters for multiple symbols (112 fields!)

---

## 📡 Subscriptions (Streaming)

* **Overview:** [Streaming\_Methods.Overview.md](./7.%20Streaming_Methods/Streaming_Methods.Overview.md)

### Real-Time Price Updates

* [OnSymbolTick.md](./7.%20Streaming_Methods/OnSymbolTick.md) - Real-time tick stream for symbols

### Trading Events

* [OnTrade.md](./7.%20Streaming_Methods/OnTrade.md) - Position/order changes (opened, closed, modified)
* [OnTradeTransaction.md](./7.%20Streaming_Methods/OnTradeTransaction.md) - Detailed transaction log (complete audit trail)

### Position Monitoring

* [OnPositionProfit.md](./7.%20Streaming_Methods/OnPositionProfit.md) - Periodic profit/loss updates
* [OnPositionsAndPendingOrdersTickets.md](./7.%20Streaming_Methods/OnPositionsAndPendingOrdersTickets.md) - Periodic ticket lists (lightweight)

---

## 🎯 Quick Navigation by Use Case

| I want to... | Use this method |
|-------------|-----------------|
| **ACCOUNT INFORMATION** |
| Get complete account snapshot | [AccountSummary](./1.%20Account_information/AccountSummary.md) |
| Get account balance | [AccountInfoDouble](./1.%20Account_information/AccountInfoDouble.md) (BALANCE) |
| Get account equity | [AccountInfoDouble](./1.%20Account_information/AccountInfoDouble.md) (EQUITY) |
| Get account leverage | [AccountInfoInteger](./1.%20Account_information/AccountInfoInteger.md) (LEVERAGE) |
| Get account currency | [AccountInfoString](./1.%20Account_information/AccountInfoString.md) (CURRENCY) |
| **SYMBOL INFORMATION** |
| Get current price for symbol | [SymbolInfoTick](./2.%20Symbol_information/SymbolInfoTick.md) |
| List all available symbols | [SymbolsTotal](./2.%20Symbol_information/SymbolsTotal.md) + [SymbolName](./2.%20Symbol_information/SymbolName.md) |
| Add symbol to Market Watch | [SymbolSelect](./2.%20Symbol_information/SymbolSelect.md) (true) |
| Get symbol digits (decimal places) | [SymbolInfoInteger](./2.%20Symbol_information/SymbolInfoInteger.md) (DIGITS) |
| Get point size for symbol | [SymbolInfoDouble](./2.%20Symbol_information/SymbolInfoDouble.md) (POINT) |
| Get complete symbol data (batch) | [SymbolParamsMany](./6.%20Additional_Methods/SymbolParamsMany.md) |
| **POSITIONS & ORDERS** |
| Count open positions | [PositionsTotal](./3.%20Position_Orders_Information/PositionsTotal.md) |
| Get all open positions (full details) | [OpenedOrders](./3.%20Position_Orders_Information/OpenedOrders.md) |
| Get position ticket numbers only | [OpenedOrdersTickets](./3.%20Position_Orders_Information/OpenedOrdersTickets.md) |
| Get historical orders | [OrderHistory](./3.%20Position_Orders_Information/OrderHistory.md) |
| Get historical deals/trades | [PositionsHistory](./3.%20Position_Orders_Information/PositionsHistory.md) |
| **MARKET DEPTH** |
| Subscribe to Level II quotes | [MarketBookAdd](./5.%20Market_Depth(DOM)/MarketBookAdd.md) |
| Get order book data | [MarketBookGet](./5.%20Market_Depth(DOM)/MarketBookGet.md) |
| Unsubscribe from Level II | [MarketBookRelease](./5.%20Market_Depth(DOM)/MarketBookRelease.md) |
| **TRADING OPERATIONS** |
| Open market BUY position | [OrderSend](./4.%20Trading_Operations/OrderSend.md) (type=BUY) |
| Open market SELL position | [OrderSend](./4.%20Trading_Operations/OrderSend.md) (type=SELL) |
| Place BUY LIMIT order | [OrderSend](./4.%20Trading_Operations/OrderSend.md) (type=BUY_LIMIT) |
| Place SELL LIMIT order | [OrderSend](./4.%20Trading_Operations/OrderSend.md) (type=SELL_LIMIT) |
| Place BUY STOP order | [OrderSend](./4.%20Trading_Operations/OrderSend.md) (type=BUY_STOP) |
| Place SELL STOP order | [OrderSend](./4.%20Trading_Operations/OrderSend.md) (type=SELL_STOP) |
| Modify SL/TP of position | [OrderModify](./4.%20Trading_Operations/OrderModify.md) |
| Close a position | [OrderClose](./4.%20Trading_Operations/OrderClose.md) |
| Calculate margin before trade | [OrderCalcMargin](./4.%20Trading_Operations/OrderCalcMargin.md) |
| Calculate potential profit | [OrderCalcProfit](./4.%20Trading_Operations/OrderCalcProfit.md) |
| Validate trade before execution | [OrderCheck](./4.%20Trading_Operations/OrderCheck.md) |
| **REAL-TIME SUBSCRIPTIONS** |
| Stream live prices | [OnSymbolTick](./7.%20Streaming_Methods/OnSymbolTick.md) |
| Monitor trade events | [OnTrade](./7.%20Streaming_Methods/OnTrade.md) |
| Track profit changes | [OnPositionProfit](./7.%20Streaming_Methods/OnPositionProfit.md) |
| Monitor ticket changes | [OnPositionsAndPendingOrdersTickets](./7.%20Streaming_Methods/OnPositionsAndPendingOrdersTickets.md) |
| Detailed transaction log | [OnTradeTransaction](./7.%20Streaming_Methods/OnTradeTransaction.md) |

---

## 🏗️ API Architecture

### Layer 1: MT5Account (Low-Level) - You are here!

**What:** Direct proto/gRPC communication with MT5 terminal.

**When to use:**
- Need full control over protocol
- Building custom wrappers
- Proto-level integration required

**Characteristics:**
- Works with proto Request/Response objects
- Raw gRPC method calls
- Complete access to all MT5 functions
- Highest complexity

**Location:** Core gRPC client (generated from proto files)

**Documentation:** This folder (you are here!)

---

### Layer 2: MT5Service

**What:** Simplified wrapper methods without proto complexity.

**When to use:**
- Want simplified API but not auto-normalization
- Building custom convenience layers
- Need direct data returns

**Characteristics:**
- Simple method signatures
- Type conversions (proto → Go primitives)
- No proto objects in return values
- No auto-normalization

**Location:** `examples/mt5/MT5Service.go`

**Documentation:** [MT5Service.Overview.md](../MT5Service/MT5Service.Overview.md) *(if exists)*

---

### Layer 3: MT5Sugar

**What:** High-level convenience API with ~60 smart methods.

**When to use:**
- Most trading scenarios (95% of cases)
- Want auto-normalization
- Need risk management helpers
- Building strategies quickly

**Characteristics:**
- Auto-normalization of volumes/prices
- Risk-based position sizing
- Batch operations
- Smart helpers

**Location:** `examples/mt5/MT5Sugar.go`

**Documentation:** [MT5Sugar.API_Overview.md](../MT5Sugar/MT5Sugar.API_Overview.md)

---

## 🎓 Learning Path

**Recommended sequence:** Start from foundation (MT5Account) → Build up to convenience layers (MT5Service → MT5Sugar)

### Step 1: Master the Foundation (MT5Account) - You are here!

**Why first:** MT5Account is the foundation - everything else is built on top of it. Understanding the protocol level gives you complete control.

```
1. Read: This documentation folder (docs/MT5Account/)
2. Study: Proto definitions and gRPC communication
3. Understand: Request/Response patterns
4. Learn: Return codes and error handling
5. Practice: Low-level method calls
```

**Goal:** Deep understanding of MT5 protocol and terminal communication.

**Documentation:** [MT5Account.Master.Overview.md](./MT5Account.Master.Overview.md) (you are here!)

---

### Step 2: Understand Wrappers (MT5Service)

**Why second:** Once you know the foundation, see how MT5Service simplifies it by wrapping proto objects.

```
1. Study: How MT5Service wraps MT5Account methods
2. Compare: Wrapper vs low-level implementations
3. Learn: Type conversions (proto → Go primitives)
4. Practice: Simplified method calls without proto objects
```

**Goal:** Learn to build clean API wrappers on top of complex protocols.

---

### Step 3: Use Convenience Layer (MT5Sugar)

**Why last:** With foundation + wrappers understood, appreciate how MT5Sugar adds auto-normalization and smart helpers.

```
1. Study: MT5Sugar convenience methods (~60 methods)
2. Learn: Auto-normalization, risk management, batch operations
3. Use: High-level methods for rapid strategy development
4. Build: Trading strategies using Sugar API
```

**Goal:** Rapid strategy development with production-ready convenience methods.

**Documentation:** [MT5Sugar.API_Overview.md](../MT5Sugar/MT5Sugar.API_Overview.md)

---

**Summary:** MT5Account (foundation) → MT5Service (wrappers) → MT5Sugar (convenience)

---

## 💡 Key Concepts

### Proto Return Codes

* **10009** = Success / DONE
* **10004** = Requote
* **10006** = Request rejected
* **10013** = Invalid request
* **10014** = Invalid volume
* **10015** = Invalid price
* **10016** = Invalid stops
* **10018** = Market closed
* **10019** = Not enough money
* **10031** = No connection with trade server

Always check `ReturnedCode` field in trading operations.



---

## ⚠️ Important Notes

* **Demo account first:** Always test on demo before live trading.
* **Check return codes:** Every trading operation returns status code (10009 = success).
* **Validate parameters:** Use `OrderCheck()` before `OrderSend()`.
* **Handle errors:** Network/protocol errors can occur.
* **Context management:** Always use context with timeout for requests.
* **Stream cleanup:** Cancel context to stop streams properly.
* **UTC timestamps:** All times are in UTC, not local time.
* **Broker limitations:** Not all brokers support all features (DOM, hedging, etc.).


---

"Trade safe, code clean, and may your gRPC streams always flow smoothly."
