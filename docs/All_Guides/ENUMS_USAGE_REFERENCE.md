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

| Method | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|--------------|
| **AccountSummary** | Get all account data in one call (RECOMMENDED) | - | `MrpcEnumAccountTradeMode` - Account type (DEMO/CONTEST/REAL) |
| **AccountInfoDouble** | Get double properties (Balance, Equity, Margin, etc.) | `AccountInfoDoublePropertyType` - Property selector (14 values: ACCOUNT_BALANCE, ACCOUNT_EQUITY, ACCOUNT_MARGIN, ACCOUNT_MARGIN_FREE, ACCOUNT_MARGIN_LEVEL, ACCOUNT_PROFIT, ACCOUNT_CREDIT, ACCOUNT_MARGIN_SO_CALL, ACCOUNT_MARGIN_SO_SO, ACCOUNT_MARGIN_INITIAL, ACCOUNT_MARGIN_MAINTENANCE, ACCOUNT_ASSETS, ACCOUNT_LIABILITIES, ACCOUNT_COMMISSION_BLOCKED) | - |
| **AccountInfoInteger** | Get integer properties (Login, Leverage, etc.) | `AccountInfoIntegerPropertyType` - Property selector (11 values: ACCOUNT_LOGIN, ACCOUNT_TRADE_MODE, ACCOUNT_LEVERAGE, ACCOUNT_LIMIT_ORDERS, ACCOUNT_MARGIN_SO_MODE, ACCOUNT_TRADE_ALLOWED, ACCOUNT_TRADE_EXPERT, ACCOUNT_MARGIN_MODE, ACCOUNT_CURRENCY_DIGITS, ACCOUNT_FIFO_CLOSE, ACCOUNT_HEDGE_ALLOWED) | ⚠️ Returns `int64` values that may represent ENUMs (e.g., ACCOUNT_TRADE_MODE: 0=DEMO, 1=CONTEST, 2=REAL) |
| **AccountInfoString** | Get string properties (Currency, Company, etc.) | `AccountInfoStringPropertyType` - Property selector (4 values: ACCOUNT_NAME, ACCOUNT_SERVER, ACCOUNT_CURRENCY, ACCOUNT_COMPANY) | - |

---

## 2. Symbol Information (14 methods)

### 7 methods use ENUMs

| Method | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|--------------|
| **SymbolsTotal** | Count total/selected symbols | - | - |
| **SymbolExist** | Check if symbol exists | - | - |
| **SymbolName** | Get symbol name by index | - | - |
| **SymbolSelect** | Add/remove symbol from Market Watch | - | - |
| **SymbolIsSynchronized** | Check sync status with server | - | - |
| **SymbolInfoDouble** | Get double properties (Bid, Ask, Point, Volume, etc.) | `SymbolInfoDoubleProperty` - Property selector (SYMBOL_BID, SYMBOL_ASK, SYMBOL_POINT, SYMBOL_VOLUME_MIN, SYMBOL_VOLUME_MAX, etc.) | - |
| **SymbolInfoInteger** | Get integer properties (Digits, Spread, Stops Level) | `SymbolInfoIntegerProperty` - Property selector | ⚠️ Returns `int64` values that may represent ENUMs: BMT5_ENUM_SYMBOL_TRADE_MODE, BMT5_ENUM_SYMBOL_TRADE_EXECUTION, BMT5_ENUM_SYMBOL_CALC_MODE, BMT5_ENUM_SYMBOL_SWAP_MODE, BMT5_ENUM_SYMBOL_ORDER_GTC_MODE, BMT5_ENUM_SYMBOL_OPTION_RIGHT, BMT5_ENUM_SYMBOL_CHART_MODE, BMT5_ENUM_SYMBOL_SECTOR, BMT5_ENUM_SYMBOL_INDUSTRY |
| **SymbolInfoString** | Get string properties (Description, Base/Profit Currency) | `SymbolInfoStringProperty` - Property selector (SYMBOL_BASIS, SYMBOL_CATEGORY, SYMBOL_COUNTRY, SYMBOL_CURRENCY_BASE, SYMBOL_CURRENCY_PROFIT, SYMBOL_DESCRIPTION, etc.) | - |
| **SymbolInfoMarginRate** | Get margin requirements for order types | `ENUM_ORDER_TYPE` - Order type for margin calculation (ORDER_TYPE_BUY, ORDER_TYPE_SELL, ORDER_TYPE_BUY_LIMIT, ORDER_TYPE_SELL_LIMIT, ORDER_TYPE_BUY_STOP, ORDER_TYPE_SELL_STOP, etc.) | - |
| **SymbolInfoTick** | Get last tick data with timestamp | - | - |
| **SymbolInfoSessionQuote** | Get quote session times | `DayOfWeek` - Day of week (SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY) | - |
| **SymbolInfoSessionTrade** | Get trade session times | `DayOfWeek` - Day of week (SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY) | - |
| **SymbolParamsMany** | Get detailed parameters for multiple symbols | `AH_SYMBOL_PARAMS_MANY_SORT_TYPE` - Sort mode (AH_PARAMS_MANY_SORT_TYPE_SYMBOL_NAME_ASC, AH_PARAMS_MANY_SORT_TYPE_SYMBOL_NAME_DESC, AH_PARAMS_MANY_SORT_TYPE_MQL_INDEX_ASC, AH_PARAMS_MANY_SORT_TYPE_MQL_INDEX_DESC) | - |
| **TickValueWithSize** | Get tick value and size information for symbols | - | - |

---

## 3. Positions & Orders Information (5 methods)

### 3 methods use ENUMs

| Method | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|--------------|
| **PositionsTotal** | Count open positions | - | - |
| **OpenedOrders** | Get all opened orders & positions with full details | `BMT5_ENUM_OPENED_ORDER_SORT_TYPE` - Sort mode | **In PositionInfo:** BMT5_ENUM_POSITION_TYPE (Position type: BUY/SELL), BMT5_ENUM_POSITION_REASON (Open reason)<br>**In OpenedOrderInfo:** BMT5_ENUM_ORDER_TYPE (Order type), BMT5_ENUM_ORDER_STATE (Order state), BMT5_ENUM_ORDER_TYPE_TIME (Order lifetime), BMT5_ENUM_ORDER_TYPE_FILLING (Fill mode) |
| **OpenedOrdersTickets** | Get only ticket numbers (lightweight) | - | - |
| **OrderHistory** | Get historical orders with pagination | `BMT5_ENUM_ORDER_HISTORY_SORT_TYPE` - Sort mode | **In OrderHistoryData:** BMT5_ENUM_ORDER_STATE, BMT5_ENUM_ORDER_TYPE, BMT5_ENUM_ORDER_TYPE_TIME, BMT5_ENUM_ORDER_TYPE_FILLING<br>**In DealHistoryData:** BMT5_ENUM_DEAL_ENTRY_TYPE, BMT5_ENUM_DEAL_TYPE, BMT5_ENUM_DEAL_REASON |
| **PositionsHistory** | Get historical positions with P&L | `AH_ENUM_POSITIONS_HISTORY_SORT_TYPE` - Sort mode (AH_SORT_BY_OPEN_TIME_ASC, AH_SORT_BY_OPEN_TIME_DESC, AH_SORT_BY_CLOSE_TIME_ASC, AH_SORT_BY_CLOSE_TIME_DESC, AH_SORT_BY_PROFIT_ASC, AH_SORT_BY_PROFIT_DESC) | `AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE` - Order type in PositionHistoryInfo (AH_ORDER_TYPE_BUY, AH_ORDER_TYPE_SELL, AH_ORDER_TYPE_BUY_LIMIT, AH_ORDER_TYPE_SELL_LIMIT, AH_ORDER_TYPE_BUY_STOP, AH_ORDER_TYPE_SELL_STOP, AH_ORDER_TYPE_BUY_STOP_LIMIT, AH_ORDER_TYPE_SELL_STOP_LIMIT) |

---

## 4. Market Depth / DOM (3 methods)

### 1 method uses ENUMs

| Method | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|--------------|
| **MarketBookAdd** | Subscribe to Depth of Market updates | - | - |
| **MarketBookRelease** | Unsubscribe from DOM | - | - |
| **MarketBookGet** | Get current market depth snapshot | - | `BookType` - Order type in order book (in MrpcMqlBookInfo.Type field): BOOK_TYPE_SELL (sell order/Offer), BOOK_TYPE_BUY (buy order/Bid), BOOK_TYPE_SELL_MARKET (market sell), BOOK_TYPE_BUY_MARKET (market buy) |

---

## 5. Trading Operations (6 methods)

### ✅ All methods use ENUMs

| Method | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|--------------|
| **OrderSend** | Send market or pending order | `TMT5_ENUM_ORDER_TYPE` - Order type (TMT5_ORDER_TYPE_BUY, TMT5_ORDER_TYPE_SELL, TMT5_ORDER_TYPE_BUY_LIMIT, TMT5_ORDER_TYPE_SELL_LIMIT, TMT5_ORDER_TYPE_BUY_STOP, TMT5_ORDER_TYPE_SELL_STOP, TMT5_ORDER_TYPE_BUY_STOP_LIMIT, TMT5_ORDER_TYPE_SELL_STOP_LIMIT, TMT5_ORDER_TYPE_CLOSE_BY)<br>`TMT5_ENUM_ORDER_TYPE_TIME` - Order lifetime (TMT5_ORDER_TIME_GTC, TMT5_ORDER_TIME_DAY, TMT5_ORDER_TIME_SPECIFIED, TMT5_ORDER_TIME_SPECIFIED_DAY) | - |
| **OrderModify** | Modify existing order parameters | `TMT5_ENUM_ORDER_TYPE_TIME` - Order lifetime (used in ExpirationTimeType field) | - |
| **OrderClose** | Close market or pending order | - | `MRPC_ORDER_CLOSE_MODE` - Close mode: MRPC_MARKET_ORDER_CLOSE (close market position), MRPC_MARKET_ORDER_PARTIAL_CLOSE (partial close), MRPC_PENDING_ORDER_REMOVE (remove pending order) |
| **OrderCheck** | Validate order before sending | `MRPC_ENUM_TRADE_REQUEST_ACTIONS` - Trade operation action<br>`ENUM_ORDER_TYPE_TF` - Order type (for TradeFunction)<br>`MRPC_ENUM_ORDER_TYPE_FILLING` - Fill mode<br>`MRPC_ENUM_ORDER_TYPE_TIME` - Order lifetime | - |
| **OrderCalcMargin** | Calculate required margin | `ENUM_ORDER_TYPE_TF` - Order type (for TradeFunction) | - |
| **OrderCalcProfit** | Calculate potential profit/loss | `ENUM_ORDER_TYPE_TF` - Order type (for TradeFunction) | - |

---

## 6. Streaming Methods (5 methods)

### 3 methods use ENUMs

| Method | Description | Input ENUMs | Output ENUMs |
|--------|-------------|-------------|--------------|
| **OnSymbolTick** | Stream tick data (Bid/Ask updates) | - | - |
| **OnTrade** | Stream trade events | - | `MT5_SUB_ENUM_EVENT_GROUP_TYPE` - Event type (in OnTradeData.Type field) |
| **OnPositionProfit** | Stream position P&L updates | - | `MT5_SUB_ENUM_EVENT_GROUP_TYPE` - Event type (in OnPositionProfitData.Type field) |
| **OnPositionsAndPendingOrdersTickets** | Stream ticket changes | - | - |
| **OnTradeTransaction** | Stream trade transaction events | - | **9 unique ENUM types, 11 field uses:**<br>1. `MT5_SUB_ENUM_EVENT_GROUP_TYPE` (OnTradeTransactionData.Type)<br>2. `SUB_ENUM_TRADE_TRANSACTION_TYPE` (MqlTradeTransaction.Type)<br>3. `SUB_ENUM_ORDER_TYPE` (used 2×: MqlTradeTransaction.OrderType, MqlTradeRequest.OrderType)<br>4. `SUB_ENUM_ORDER_STATE` (MqlTradeTransaction.OrderState)<br>5. `SUB_ENUM_DEAL_TYPE` (MqlTradeTransaction.DealType)<br>6. `SUB_ENUM_ORDER_TYPE_TIME` (used 2×: MqlTradeTransaction.OrderTimeType, MqlTradeRequest.TypeTime)<br>7. `SUB_ENUM_TRADE_REQUEST_ACTIONS` (MqlTradeRequest.TradeOperationType)<br>8. `SUB_ENUM_ORDER_TYPE_FILLING` (MqlTradeRequest.OrderTypeFilling)<br>9. `MqlErrorTradeCode` (MqlTradeResult.TradeReturnCode) |

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
    SortType: pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_SORT_BY_PROFIT_DESC,
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
// OnTrade - Output ENUM
dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})
go func() {
    for event := range dataChan {
        // Check event type
        switch event.Type {
        case pb.MT5_SUB_ENUM_EVENT_GROUP_TYPE_OnTrade:
            fmt.Println("Trade event received")
            // Process trade data
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
- Output ENUM: only AH_ENUM_POSITIONS_HISTORY_ORDER_TYPE

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

**6. OnTradeTransaction - most complex:**

- Uses **11 fields with ENUMs** (9 unique types) in four nested structures:
  - `OnTradeTransactionData.Type` - MT5_SUB_ENUM_EVENT_GROUP_TYPE
  - `MqlTradeTransaction` - 5 ENUMs (Type, OrderType, OrderState, DealType, OrderTimeType)
  - `MqlTradeRequest` - 4 ENUMs (TradeOperationType, OrderType, OrderTypeFilling, TypeTime)
  - `MqlTradeResult.TradeReturnCode` - MqlErrorTradeCode
- Some ENUM types are reused in different structures (SUB_ENUM_ORDER_TYPE, SUB_ENUM_ORDER_TYPE_TIME)

**7. Methods with Multiple Nested ENUMs:**

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
