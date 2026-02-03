# ENUMs Usage Reference - All Methods

> Complete reference for ENUM usage across all GoMT5 API methods

---

## 📊 Summary Statistics

| Method Group | With ENUMs | Total | Percentage |
|--------------|------------|-------|------------|
| **Account Information** | 4 | 4 | 100% ✅ |
| **Symbol Information** | 7 | 14 | 50% |
| **Positions & Orders** | 3 | 5 | 60% |
| **Market Depth/DOM** | 1 | 3 | 33% |
| **Trading Operations** | 6 | 6 | 100% ✅ |
| **Streaming Methods** | 3 | 5 | 60% |
| **TOTAL** | **24** | **37** | **65%** |

---

## 1. Account Information (4 methods)

### ✅ All methods use ENUMs

| Method | ENUMs Count | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|-------------|--------------|
| **AccountSummary** | **[In: 0, Out: 1]**<br>Total: 1 ENUM | Get all account data in one call (RECOMMENDED) | - | `MrpcEnumAccountTradeMode` (3 values) - Account type (DEMO/CONTEST/REAL) |
| **AccountInfoDouble** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get double properties (Balance, Equity, Margin, etc.) | `AccountInfoDoublePropertyType` (14 values) - Property selector | - |
| **AccountInfoInteger** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get integer properties (Login, Leverage, etc.) | `AccountInfoIntegerPropertyType` (11 values) - Property selector | ⚠️ Returns `int64` values that may represent ENUMs (e.g., ACCOUNT_TRADE_MODE: 0=DEMO, 1=CONTEST, 2=REAL) |
| **AccountInfoString** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get string properties (Currency, Company, etc.) | `AccountInfoStringPropertyType` (4 values) - Property selector | - |

---

## 2. Symbol Information (14 methods)

### 7 methods use ENUMs

| Method | ENUMs Count | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|-------------|--------------|
| **SymbolsTotal** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Count total/selected symbols | - | - |
| **SymbolExist** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Check if symbol exists | - | - |
| **SymbolName** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Get symbol name by index | - | - |
| **SymbolSelect** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Add/remove symbol from Market Watch | - | - |
| **SymbolIsSynchronized** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Check sync status with server | - | - |
| **SymbolInfoDouble** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get double properties (Bid, Ask, Point, Volume, etc.) | `SymbolInfoDoubleProperty` (60 values) - Property selector | - |
| **SymbolInfoInteger** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get integer properties (Digits, Spread, Stops Level) | `SymbolInfoIntegerProperty` (37 values) - Property selector | ⚠️ Returns `int64` values that may represent ENUMs: BMT5_ENUM_SYMBOL_TRADE_MODE, BMT5_ENUM_SYMBOL_TRADE_EXECUTION, BMT5_ENUM_SYMBOL_CALC_MODE, BMT5_ENUM_SYMBOL_SWAP_MODE, BMT5_ENUM_SYMBOL_ORDER_GTC_MODE, BMT5_ENUM_SYMBOL_OPTION_RIGHT, BMT5_ENUM_SYMBOL_CHART_MODE, BMT5_ENUM_SYMBOL_SECTOR, BMT5_ENUM_SYMBOL_INDUSTRY |
| **SymbolInfoString** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get string properties (Description, Base/Profit Currency) | `SymbolInfoStringProperty` (15 values) - Property selector | - |
| **SymbolInfoMarginRate** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get margin requirements for order types | `ENUM_ORDER_TYPE` (9 values) - Order type for margin calculation | - |
| **SymbolInfoTick** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Get last tick data with timestamp | - | - |
| **SymbolInfoSessionQuote** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get quote session times | `DayOfWeek` (7 values) - Day of week | - |
| **SymbolInfoSessionTrade** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get trade session times | `DayOfWeek` (7 values) - Day of week | - |
| **SymbolParamsMany** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Get detailed parameters for multiple symbols | `AH_SYMBOL_PARAMS_MANY_SORT_TYPE` (4 values) - Sort mode | - |
| **TickValueWithSize** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Get tick value and size information for symbols | - | - |

---

## 3. Positions & Orders Information (5 methods)

### 3 methods use ENUMs

| Method | ENUMs Count | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|-------------|--------------|
| **PositionsTotal** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Count open positions | - | - |
| **OpenedOrders** | **[In: 1, Out: 6]**<br>Total: 7 ENUMs | Get all opened orders & positions with full details | `BMT5_ENUM_OPENED_ORDER_SORT_TYPE` (4 values) - Sort mode | **6 ENUMs:**<br>**In PositionInfo:** BMT5_ENUM_POSITION_TYPE, BMT5_ENUM_POSITION_REASON<br>**In OpenedOrderInfo:** BMT5_ENUM_ORDER_TYPE, BMT5_ENUM_ORDER_STATE, BMT5_ENUM_ORDER_TYPE_TIME, BMT5_ENUM_ORDER_TYPE_FILLING |
| **OpenedOrdersTickets** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Get only ticket numbers (lightweight) | - | - |
| **OrderHistory** | **[In: 1, Out: 7]**<br>Total: 8 ENUMs | Get historical orders with pagination | `BMT5_ENUM_ORDER_HISTORY_SORT_TYPE` (6 values) - Sort mode | **7 ENUMs:**<br>**In OrderHistoryData:** BMT5_ENUM_ORDER_STATE, BMT5_ENUM_ORDER_TYPE, BMT5_ENUM_ORDER_TYPE_TIME, BMT5_ENUM_ORDER_TYPE_FILLING<br>**In DealHistoryData:** BMT5_ENUM_DEAL_ENTRY_TYPE, BMT5_ENUM_DEAL_TYPE, BMT5_ENUM_DEAL_REASON |
| **PositionsHistory** | **[In: 1, Out: 1]**<br>Total: 2 ENUMs | Get historical positions with P&L | `AH_ENUM_POSITIONS_HISTORY_SORT_TYPE` (4 values) - Sort mode | `AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE` (9 values) - Order type in PositionHistoryInfo |

---

## 4. Market Depth / DOM (3 methods)

### 1 method uses ENUMs

| Method | ENUMs Count | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|-------------|--------------|
| **MarketBookAdd** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Subscribe to Depth of Market updates | - | - |
| **MarketBookRelease** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Unsubscribe from DOM | - | - |
| **MarketBookGet** | **[In: 0, Out: 1]**<br>Total: 1 ENUM | Get current market depth snapshot | - | `BookType` (4 values) - Order type in order book (BOOK_TYPE_SELL, BOOK_TYPE_BUY, BOOK_TYPE_SELL_MARKET, BOOK_TYPE_BUY_MARKET) |

---

## 5. Trading Operations (6 methods)

### ✅ All methods use ENUMs

| Method | ENUMs Count | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|-------------|--------------|
| **OrderSend** | **[In: 2, Out: 0]**<br>Total: 2 ENUMs | Send market or pending order | `TMT5_ENUM_ORDER_TYPE` (9 values) - Order type<br>`TMT5_ENUM_ORDER_TYPE_TIME` (4 values) - Order lifetime | - |
| **OrderModify** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Modify existing order parameters | `TMT5_ENUM_ORDER_TYPE_TIME` (4 values) - Order lifetime (ExpirationTimeType field) | - |
| **OrderClose** | **[In: 0, Out: 1]**<br>Total: 1 ENUM | Close market or pending order | - | `MRPC_ORDER_CLOSE_MODE` (3 values) - Close mode: MRPC_MARKET_ORDER_CLOSE, MRPC_MARKET_ORDER_PARTIAL_CLOSE, MRPC_PENDING_ORDER_REMOVE |
| **OrderCheck** | **[In: 4, Out: 0]**<br>Total: 4 ENUMs | Validate order before sending | `MRPC_ENUM_TRADE_REQUEST_ACTIONS` (6 values) - Trade action<br>`ENUM_ORDER_TYPE_TF` (9 values) - Order type<br>`MRPC_ENUM_ORDER_TYPE_FILLING` (4 values) - Fill mode<br>`MRPC_ENUM_ORDER_TYPE_TIME` (4 values) - Order lifetime | - |
| **OrderCalcMargin** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Calculate required margin | `ENUM_ORDER_TYPE_TF` (9 values) - Order type | - |
| **OrderCalcProfit** | **[In: 1, Out: 0]**<br>Total: 1 ENUM | Calculate potential profit/loss | `ENUM_ORDER_TYPE_TF` (9 values) - Order type | - |

---

## 6. Streaming Methods (5 methods)

### ✅ 3/5 methods use ENUMs (60%)

| Method | ENUMs Count | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|-------------|--------------|
| **OnSymbolTick** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Stream tick data (Bid/Ask updates) | - | - |
| **OnTrade** | **[In: 0, Out: 11]**<br>Total: 11 ENUMs | Stream trade events (positions, orders, deals) | - | **11 ENUMs (1 direct, 10 nested):**<br>**Direct ENUM (1):**<br>1. `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (OnTradeData.Type) - 2 values<br>**Nested ENUMs (10):**<br>2. `SUB_ENUM_POSITION_TYPE` (OnTradePositionInfo)<br>3. `SUB_ENUM_POSITION_REASON` (OnTradePositionInfo)<br>4. `SUB_ENUM_ORDER_TYPE` (OnTradeOrderInfo, OnTradeHistoryOrderInfo)<br>5. `SUB_ENUM_ORDER_STATE` (OnTradeOrderInfo, OnTradeHistoryOrderInfo)<br>6. `SUB_ENUM_ORDER_TYPE_TIME` (OnTradeOrderInfo, OnTradeHistoryOrderInfo)<br>7. `SUB_ENUM_ORDER_TYPE_FILLING` (OnTradeOrderInfo, OnTradeHistoryOrderInfo)<br>8. `SUB_ENUM_ORDER_REASON` (OnTradeOrderInfo)<br>9. `SUB_ENUM_DEAL_TYPE` (OnTradeHistoryDealInfo)<br>10. `SUB_ENUM_DEAL_ENTRY` (OnTradeHistoryDealInfo)<br>11. `SUB_ENUM_DEAL_REASON` (OnTradeHistoryDealInfo, OnTradeHistoryOrderInfo) |
| **OnPositionProfit** | **[In: 0, Out: 1]**<br>Total: 1 ENUM | Stream position P&L updates | - | `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (2 values) - Event type (OnPositionProfitData.Type) |
| **OnPositionsAndPendingOrdersTickets** | **[In: 0, Out: 0]**<br>Total: 0 ENUMs | Stream ticket changes | - | - |
| **OnTradeTransaction** | **[In: 0, Out: 9]**<br>Total: 9 ENUMs<br>(11 field uses) | Stream trade transaction events | - | **9 ENUMs (1 direct, 8 nested):**<br>**Direct ENUM (1):**<br>1. `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (OnTradeTransactionData.Type)<br>**Nested ENUMs (8):**<br>2. `SUB_ENUM_TRADE_TRANSACTION_TYPE` (MqlTradeTransaction.Type)<br>3. `SUB_ENUM_ORDER_TYPE` (×2: MqlTradeTransaction.OrderType, MqlTradeRequest.OrderType)<br>4. `SUB_ENUM_ORDER_STATE` (MqlTradeTransaction.OrderState)<br>5. `SUB_ENUM_DEAL_TYPE` (MqlTradeTransaction.DealType)<br>6. `SUB_ENUM_ORDER_TYPE_TIME` (×2: MqlTradeTransaction.OrderTimeType, MqlTradeRequest.TypeTime)<br>7. `SUB_ENUM_TRADE_REQUEST_ACTIONS` (MqlTradeRequest.TradeOperationType)<br>8. `SUB_ENUM_ORDER_TYPE_FILLING` (MqlTradeRequest.OrderTypeFilling)<br>9. `MqlErrorTradeCode` (MqlTradeResult.TradeReturnCode) |

---

## 🔧 Code Usage Examples

### 1. Account Information

```go
// AccountSummary - Output ENUM
summary, _ := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
switch summary.AccountTradeMode {
case pb.MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_DEMO:
    fmt.Println("Demo account")
case pb.MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_REAL:
    fmt.Println("Real account")
}

// AccountInfoDouble - Input ENUM
req := &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
}
data, _ := account.AccountInfoDouble(ctx, req)
fmt.Printf("Balance: %.2f\n", data.RequestedValue)
```

### 2. Trading Operations

```go
// OrderSend - Input ENUMs
req := &pb.OrderSendRequest{
    Symbol:    "EURUSD",
    Volume:    0.01,
    OrderType: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
    TypeTime:  pb.TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_GTC,
}
result, _ := account.OrderSend(ctx, req)

// OrderModify - Input ENUM
modifyReq := &pb.OrderModifyRequest{
    Ticket:             ticket,
    StopLoss:           &newSL,
    ExpirationTimeType: pb.TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_DAY,
}
modifyResult, _ := account.OrderModify(ctx, modifyReq)

// OrderClose - Output ENUM
closeResult, _ := account.OrderClose(ctx, closeReq)
switch closeResult.CloseMode {
case pb.MRPC_ORDER_CLOSE_MODE_MRPC_MARKET_ORDER_CLOSE:
    fmt.Println("Market order closed")
case pb.MRPC_ORDER_CLOSE_MODE_MRPC_PENDING_ORDER_REMOVE:
    fmt.Println("Pending order removed")
}
```

### 3. Symbol Information

```go
// SymbolInfoSessionQuote - Input ENUM
req := &pb.SymbolInfoSessionQuoteRequest{
    Symbol:    "EURUSD",
    DayOfWeek: pb.BMT5_ENUM_DAY_OF_WEEK_BMT5_MONDAY,
}
sessionData, _ := account.SymbolInfoSessionQuote(ctx, req)
```

### 4. Positions & Orders

```go
// OpenedOrders - Input ENUM
req := &pb.OpenedOrdersRequest{
    SortType: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_SORT_BY_TICKET_ASC,
}
orders, _ := account.OpenedOrders(ctx, req)

// PositionsHistory - Input ENUM
histReq := &pb.PositionsHistoryRequest{
    SortType: pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_DESC,
}
history, _ := account.PositionsHistory(ctx, histReq)

// Check Output ENUM in result
for _, pos := range history.HistoryPositions {
    switch pos.OrderType {
    case pb.AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_BUY:
        fmt.Println("Buy position")
    case pb.AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE_AH_ORDER_TYPE_SELL:
        fmt.Println("Sell position")
    }
}
```

### 5. Market Depth / DOM

```go
// MarketBookGet - Output ENUM
req := &pb.MarketBookGetRequest{
    Symbol: "EURUSD",
}
bookData, _ := account.MarketBookGet(ctx, req)

// Check order types in the book
for _, bookInfo := range bookData.MqlBookInfos {
    switch bookInfo.Type {
    case pb.BookType_BOOK_TYPE_SELL:
        fmt.Printf("Sell order: Price=%.5f, Volume=%d\n", bookInfo.Price, bookInfo.Volume)
    case pb.BookType_BOOK_TYPE_BUY:
        fmt.Printf("Buy order: Price=%.5f, Volume=%d\n", bookInfo.Price, bookInfo.Volume)
    case pb.BookType_BOOK_TYPE_SELL_MARKET:
        fmt.Printf("Market sell: Price=%.5f, Volume=%d\n", bookInfo.Price, bookInfo.Volume)
    case pb.BookType_BOOK_TYPE_BUY_MARKET:
        fmt.Printf("Market buy: Price=%.5f, Volume=%d\n", bookInfo.Price, bookInfo.Volume)
    }
}
```

### 6. Streaming Methods

```go
// OnTrade - Output ENUMs (11 unique types in nested structures)
dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})
go func() {
    for event := range dataChan {
        // Check top-level event type
        switch event.Type {
        case pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate:
            fmt.Println("Trade event received")

            // Process positions (uses SUB_ENUM_POSITION_TYPE, SUB_ENUM_POSITION_REASON)
            for _, pos := range event.EventData.NewPositions {
                switch pos.Type {
                case pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_BUY:
                    fmt.Printf("New BUY position #%d\n", pos.Ticket)
                case pb.SUB_ENUM_POSITION_TYPE_SUB_POSITION_TYPE_SELL:
                    fmt.Printf("New SELL position #%d\n", pos.Ticket)
                }
            }

            // Process orders (uses SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE, etc.)
            for _, order := range event.EventData.NewOrders {
                switch order.State {
                case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_PLACED:
                    fmt.Printf("Order #%d placed\n", order.Ticket)
                case pb.SUB_ENUM_ORDER_STATE_SUB_ORDER_STATE_FILLED:
                    fmt.Printf("Order #%d filled\n", order.Ticket)
                }
            }

            // Process deals (uses SUB_ENUM_DEAL_TYPE, SUB_ENUM_DEAL_ENTRY, SUB_ENUM_DEAL_REASON)
            for _, deal := range event.EventData.NewHistoryDeals {
                switch deal.Type {
                case pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_BUY:
                    fmt.Printf("Buy deal #%d\n", deal.Ticket)
                case pb.SUB_ENUM_DEAL_TYPE_SUB_DEAL_TYPE_SELL:
                    fmt.Printf("Sell deal #%d\n", deal.Ticket)
                }
            }
        }
    }
}()

// OnPositionProfit - Output ENUM
profitChan, profitErrChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{
    TimerPeriodMilliseconds: 1000,
    IgnoreEmptyData:         true,
})
go func() {
    for event := range profitChan {
        // Check event type (same ENUM as OnTrade)
        switch event.Type {
        case pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit:
            fmt.Printf("Position profit updates: %d new, %d updated, %d deleted\n",
                len(event.NewPositions),
                len(event.UpdatedPositions),
                len(event.DeletedPositions))
        }
    }
}()

// OnTradeTransaction - Output ENUMs (9 unique types, 11 uses)
txChan, txErrChan := account.OnTradeTransaction(ctx, &pb.OnTradeTransactionRequest{})
go func() {
    for tx := range txChan {
        // ENUM #1: MT5_SUB_ENUM_EVENT_GROUP_TYPE
        switch tx.Type {
        case pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_TradeTransaction:
            fmt.Println("Trade transaction event")
        }

        // ENUM #2-6: In MqlTradeTransaction
        if tx.TradeTransaction != nil {
            switch tx.TradeTransaction.Type {
            case pb.SUB_ENUM_TRADE_TRANSACTION_TYPE_TRADE_TRANSACTION_ORDER_ADD:
                fmt.Println("Order added")
            case pb.SUB_ENUM_TRADE_TRANSACTION_TYPE_TRADE_TRANSACTION_DEAL_ADD:
                fmt.Println("Deal added")
            }

            switch tx.TradeTransaction.OrderType {
            case pb.SUB_ENUM_ORDER_TYPE_ORDER_TYPE_BUY:
                fmt.Println("Buy order")
            case pb.SUB_ENUM_ORDER_TYPE_ORDER_TYPE_SELL:
                fmt.Println("Sell order")
            }

            switch tx.TradeTransaction.OrderState {
            case pb.SUB_ENUM_ORDER_STATE_ORDER_STATE_STARTED:
                fmt.Println("Order started")
            case pb.SUB_ENUM_ORDER_STATE_ORDER_STATE_FILLED:
                fmt.Println("Order filled")
            }

            switch tx.TradeTransaction.DealType {
            case pb.SUB_ENUM_DEAL_TYPE_DEAL_TYPE_BUY:
                fmt.Println("Buy deal")
            case pb.SUB_ENUM_DEAL_TYPE_DEAL_TYPE_SELL:
                fmt.Println("Sell deal")
            }

            switch tx.TradeTransaction.OrderTimeType {
            case pb.SUB_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC:
                fmt.Println("Good Till Cancel")
            }
        }

        // ENUM #7-8: In MqlTradeRequest
        if tx.TradeRequest != nil {
            switch tx.TradeRequest.TradeOperationType {
            case pb.SUB_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL:
                fmt.Println("Deal action")
            case pb.SUB_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING:
                fmt.Println("Pending order action")
            }

            switch tx.TradeRequest.OrderTypeFilling {
            case pb.SUB_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK:
                fmt.Println("Fill or Kill")
            case pb.SUB_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_IOC:
                fmt.Println("Immediate or Cancel")
            }
        }

        // ENUM #9: In MqlTradeResult
        if tx.TradeResult != nil {
            switch tx.TradeResult.TradeReturnCode {
            case pb.MqlErrorTradeCode_TRADE_RETCODE_DONE:
                fmt.Println("Trade successful")
            case pb.MqlErrorTradeCode_TRADE_RETCODE_REJECT:
                fmt.Println("Trade rejected")
            }
        }
    }
}()
```

---

## 📝 Important Notes

### 🔴 Full Name Required in Go

In Go, you **MUST** write the full ENUM name with prefix:

✅ **Correct:**
```go
pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE
pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY
pb.MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_DEMO
```

❌ **Incorrect (won't compile):**
```go
ACCOUNT_BALANCE           // error: undefined
ORDER_TYPE_BUY           // error: undefined
ACCOUNT_TRADE_MODE_DEMO  // error: undefined
```

### 🔴 Full Name Format

```
pb.<ENUM_TYPE>_<CONSTANT_NAME>
│   │            │
│   │            └─ Constant name
│   └─ ENUM type
└─ Package prefix
```

### ⚠️ Special Cases

**1. AccountInfoInteger and SymbolInfoInteger:**

- Some properties return values that **represent ENUMs**, but are returned as `int64`
- These are **NOT typed ENUMs**, just numbers!
- You need to manually compare with known values

**2. PositionsHistory:**

- Does NOT return Deal ENUMs (BMT5_ENUM_DEAL_TYPE, BMT5_ENUM_DEAL_REASON, BMT5_ENUM_DEAL_ENTRY_TYPE)
- Returns only POSITION information, not deals
- Input ENUM: AH_ENUM_POSITIONS_HISTORY_SORT_TYPE (4 values: OPEN_TIME_ASC/DESC, TICKET_ASC/DESC)
- Output ENUM: AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE (9 values, includes CLOSE_BY)

**3. OrderModify:**

- Uses Input ENUM for ExpirationTimeType
- Does NOT use ENUMs in response (only ReturnedCode as uint32)

**4. ReturnedCode:**

- ReturnedCode field in all trading operations is **NOT an ENUM**
- It's a plain uint32 operation return code
- See [RETURN_CODES_REFERENCE.md](./RETURN_CODES_REFERENCE.md) for code list

**5. Streaming Methods - OnPositionProfit:**

- OnPositionProfit USES Output ENUM: `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (in Type field of OnPositionProfitData structure)
- This same ENUM is used in OnTrade and OnTradeTransaction
- Represents event type in stream (OrderProfit, OrderPlaced, OrderModified, TradeTransaction, etc.)

**6. OnTrade - most complex streaming method:**

- Uses **11 unique ENUM types** in nested structures (more than OnTradeTransaction!):
  - `OnTradeData.Type` - MT5_SUB_ENUM_EVENT_GROUP_TYPE (top-level event type)
  - `OnTradePositionInfo` (in NewPositions, DisappearedPositions, UpdatedPositions) - 2 ENUMs:
    - SUB_ENUM_POSITION_TYPE, SUB_ENUM_POSITION_REASON
  - `OnTradeOrderInfo` (in NewOrders, DisappearedOrders, StateChangedOrders) - 5 ENUMs:
    - SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE, SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING, SUB_ENUM_ORDER_REASON
  - `OnTradeHistoryOrderInfo` (in NewHistoryOrders, DisappearedHistoryOrders, UpdatedHistoryOrders) - 5 ENUMs:
    - SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_STATE, SUB_ENUM_ORDER_TYPE_TIME, SUB_ENUM_ORDER_TYPE_FILLING, SUB_ENUM_DEAL_REASON (note: DEAL_REASON, not ORDER_REASON!)
  - `OnTradeHistoryDealInfo` (in NewHistoryDeals, DisappearedHistoryDeals, UpdatedHistoryDeals) - 3 ENUMs:
    - SUB_ENUM_DEAL_TYPE, SUB_ENUM_DEAL_ENTRY, SUB_ENUM_DEAL_REASON
- Many ENUM types are reused across different structures
- **Important:** Always check `event.Type == MT5_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate` (not OnTrade!)

**7. OnTradeTransaction - second most complex:**

- Uses **11 fields with ENUMs** (9 unique types) in four nested structures:
  - `OnTradeTransactionData.Type` - MT5_SUB_ENUM_EVENT_GROUP_TYPE
  - `MqlTradeTransaction` - 5 ENUMs (Type, OrderType, OrderState, DealType, OrderTimeType)
  - `MqlTradeRequest` - 4 ENUMs (TradeOperationType, OrderType, OrderTypeFilling, TypeTime)
  - `MqlTradeResult.TradeReturnCode` - MqlErrorTradeCode
- Some ENUM types are reused in different structures (SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_TYPE_TIME)

**8. Methods with Multiple Nested ENUMs:**

- **OnTrade** uses 11 unique ENUM types (most complex!)
- **OnTradeTransaction** uses 9 unique ENUM types (11 field uses)
- **OrderHistory** returns 7 different ENUMs (4 in OrderHistoryData + 3 in DealHistoryData)
- **OpenedOrders** returns 6 different ENUMs (3 in PositionInfo + 3 in OpenedOrderInfo)

---

## 🔗 See Also

### Documentation by Group:

- [Account Information Overview](../MT5Account/1.%20Account_information/Account_Information.Overview.md)
- [Symbol Information Overview](../MT5Account/2.%20Symbol_information/Symbol_Information.Overview.md)
- [Positions & Orders Overview](../MT5Account/3.%20Position_Orders_Information/Position_Orders.Overview.md)
- [Market Depth Overview](../MT5Account/5.%20Market_Depth(DOM)/Market_Depth.Overview.md)
- [Trading Operations Overview](../MT5Account/4.%20Trading_Operations/Trading_Operations.Overview.md)

### Other References:

- [Return Codes Reference](./RETURN_CODES_REFERENCE.md) - trading operation return codes
- [Error Handling Guide](./ERROR_HANDLING_GUIDE.md) - error handling
