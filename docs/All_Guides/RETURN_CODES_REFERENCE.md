# MT5 Return Codes (RetCodes) Reference - Go

## What is RetCode?

**RetCode** (Return Code) is a status code that MT5 terminal returns after executing a trading operation. The code indicates success or the reason for failure.

RetCodes **appear only in trading methods** because only trading operations go through the broker and can be rejected for various reasons (insufficient margin, invalid price, market closed, etc.).

---

## 📌 RetCode Source

RetCodes are **standard MQL5 codes** defined in official MetaQuotes documentation:

**MQL5 Documentation**: [Trade Result Codes (ENUM_TRADE_RETCODE)](https://www.mql5.com/en/docs/constants/errorswarnings/enum_trade_return_codes)

These codes are:

- **Unified across all languages** (C#, Python, Java, Node.js, Go, PHP)
- **Unified with MT5 terminal** - returned directly from the trading server
- **Defined in protobuf** - enum `MqlErrorTradeCode` in the API proto contract
- **Available as Go constants** - defined in `package/Helpers/errors.go`

---

## Where is RetCode Used in the API?

RetCode is returned in **all trading operations** that modify positions or orders.

It reaches you in one of two ways - always handle both:

- **In the reply:** `ReturnedCode` of `OrderSendData` / `OrderModifyData` / `OrderCloseData` (10009 = done).
- **As an error:** a request the terminal refuses (e.g. an invalid volume) usually comes back as an `*mt5.ApiError`
  whose `MqlErrorTradeIntCode()` carries the RetCode (e.g. 10014) and `MqlErrorTradeCode()` the matching
  `pb.MqlErrorTradeCode` enum value.

### 1. OrderSend (Opening Orders)

```go
import (
    "errors"
    "fmt"

    mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
)

// Using Sugar API
ticket, err := sugar.BuyMarket("EURUSD", 0.01)
if err != nil {
    // Check if it's an ApiError with RetCode
    var apiErr *mt5.ApiError
    if errors.As(err, &apiErr) {
        fmt.Printf("RetCode: %d\n", apiErr.MqlErrorTradeIntCode())
        fmt.Printf("Description: %s\n", apiErr.MqlErrorTradeDescription())
    } else {
        // Sugar also returns a plain error "order rejected, code: N" when ReturnedCode != 10009
        fmt.Println(err)
    }
}

// Using Account API (low level)
sendData, err := account.OrderSend(ctx, req)
if err != nil {
    fmt.Printf("Error: %v\n", err) // gRPC error or *mt5.ApiError (see above)
    return
}

// IMPORTANT: Always check ReturnedCode!
if sendData.ReturnedCode == mt5.TradeRetCodeDone {
    fmt.Printf("✅ Order opened! Ticket: %d\n", sendData.Order)
} else {
    fmt.Printf("❌ Order failed: %s\n",
        mt5.GetRetCodeMessage(sendData.ReturnedCode))
}
```

### 2. OrderModify (Modifying SL/TP)

```go
import (
    "context"
    "fmt"

    mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
    pb "github.com/MetaRPC/GoMT5/package"
)

sl, tp := 1.0850, 1.0950
modifyData, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
    Ticket:     ticket,
    StopLoss:   &sl, // optional fields are pointers
    TakeProfit: &tp,
})

if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if mt5.IsRetCodeSuccess(modifyData.ReturnedCode) {
    fmt.Println("✅ SL/TP modified successfully")
} else {
    fmt.Printf("❌ Modification failed: %s\n",
        mt5.GetRetCodeMessage(modifyData.ReturnedCode))
}
```

### 3. OrderClose (Closing Positions)

```go
import (
    "context"
    "errors"
    "fmt"

    mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
    pb "github.com/MetaRPC/GoMT5/package"
)

// Using Sugar API
err := sugar.ClosePosition(ticket)
if err != nil {
    var apiErr *mt5.ApiError
    if errors.As(err, &apiErr) {
        if apiErr.MqlErrorTradeIntCode() == 10036 {
            fmt.Println("⚠️ Position already closed")
        }
    }
}

// Using Account API (OrderClose closes a position or deletes a pending order)
closeData, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
    Ticket: ticket,
})
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

if closeData.ReturnedCode == mt5.TradeRetCodeDone {
    fmt.Println("✅ Position closed")
}
```

---

## Why RetCode Only in Trading Methods?

**Informational methods** (getting prices, symbols, balance) **DO NOT return RetCode** because:

- They don't go through the broker
- They cannot be "rejected" - either data exists or there's a gRPC error
- They work with local terminal data

**Trading methods** return RetCode because:

- Request is sent to broker via trading server
- Broker validates: margin, symbol rules, trading hours, limits
- Broker can reject request for dozens of reasons
- Each reason has its own unique code

---

## Complete RetCode List

### ✅ Success Codes

| Code | Go Constant | Description | When Returned |
|------|-------------|-------------|---------------|
| **10009** | `TradeRetCodeDone` | **Request completed successfully** | Market order opened/closed |
| **10010** | `TradeRetCodeDonePartial` | **Partial execution** | Only part of volume executed |
| **10008** | `TradeRetCodePlaced` | **Pending order placed** | Limit/Stop order placed |

### ⚠️ Requote Codes (Retry Recommended)

| Code | Go Constant | Description | Action |
|------|-------------|-------------|--------|
| **10004** | `TradeRetCodeRequote` | **Requote** | Price changed, retry request |
| **10020** | `TradeRetCodePriceChanged` | **Price changed** | Similar to requote, retry |

**Helper function:**
```go
if mt5.IsRetCodeRequote(retCode) {
    // Retry with updated price
}
```

### ❌ Request Validation Errors

| Code | Go Constant | Description | Common Cause |
|------|-------------|-------------|--------------|
| **10013** | `TradeRetCodeInvalidRequest` | **Invalid request** | Incorrect parameters |
| **10014** | `TradeRetCodeInvalidVolume` | **Invalid volume** | Volume < MinVolume or > MaxVolume |
| **10015** | `TradeRetCodeInvalidPrice` | **Invalid price** | Price doesn't match symbol rules |
| **10016** | `TradeRetCodeInvalidStops` | **Invalid stops** | SL/TP too close to price (check StopLevel) |
| **10022** | `TradeRetCodeInvalidExpiration` | **Invalid expiration** | Expiration time incorrect |
| **10030** | `TradeRetCodeInvalidFill` | **Invalid order filling type** | Fill type not allowed |
| **10035** | `TradeRetCodeInvalidOrder` | **Invalid order type** | Order type prohibited for symbol |
| **10038** | `TradeRetCodeInvalidCloseVolume` | **Invalid close volume** | Close volume exceeds position volume |

### 🚫 Trading Restrictions

| Code | Go Constant | Description | Reason |
|------|-------------|-------------|--------|
| **10017** | `TradeRetCodeTradeDisabled` | **Trading disabled** | Trading disabled for symbol |
| **10018** | `TradeRetCodeMarketClosed` | **Market closed** | Outside trading hours |
| **10026** | `TradeRetCodeServerDisablesAt` | **Autotrading disabled by server** | Server disabled auto-trading |
| **10027** | `TradeRetCodeClientDisablesAt` | **Autotrading disabled by client** | Terminal disabled auto-trading |
| **10032** | `TradeRetCodeOnlyReal` | **Only real accounts** | Action unavailable on demo |
| **10042** | `TradeRetCodeLongOnly` | **Only long positions allowed** | Short positions prohibited |
| **10043** | `TradeRetCodeShortOnly` | **Only short positions allowed** | Long positions prohibited |
| **10044** | `TradeRetCodeCloseOnly` | **Only position closing allowed** | Opening new positions prohibited |
| **10045** | `TradeRetCodeFifoClose` | **Position close only by FIFO rule** | FIFO rule enforced |
| **10046** | `TradeRetCodeHedgeProhibited` | **Opposite positions prohibited** | Hedging disabled |

### 💰 Resource Limits

| Code | Go Constant | Description | Solution |
|------|-------------|-------------|----------|
| **10019** | `TradeRetCodeNoMoney` | **Insufficient funds** | Check free margin |
| **10033** | `TradeRetCodeLimitOrders` | **Pending orders limit** | Close some pending orders |
| **10034** | `TradeRetCodeLimitVolume` | **Volume limit** | Reduce position volume |
| **10040** | `TradeRetCodeLimitPositions` | **Positions limit** | Close some positions |

### 🔧 Technical Issues (Retryable)

| Code | Go Constant | Description | Retryable |
|------|-------------|-------------|-----------|
| **10011** | `TradeRetCodeError` | **Processing error** | No |
| **10012** | `TradeRetCodeTimeout` | **Request timeout** | Yes |
| **10021** | `TradeRetCodeNoQuotes` | **No quotes** | Yes |
| **10024** | `TradeRetCodeTooManyRequests` | **Rate limiting** | Yes (with delay) |
| **10028** | `TradeRetCodeLocked` | **Request locked for processing** | Yes |
| **10029** | `TradeRetCodeFrozen` | **Order or position frozen** | Yes |
| **10031** | `TradeRetCodeNoConnection` | **No connection** | Yes |

**Helper function:**
```go
if mt5.IsRetCodeRetryable(retCode) {
    // Safe to retry with exponential backoff
    time.Sleep(time.Second * 2)
    // Retry operation...
}
```

### 🔄 State Management

| Code | Go Constant | Description | Meaning |
|------|-------------|-------------|---------|
| **10023** | `TradeRetCodeOrderChanged` | **Order state changed** | Order already modified/closed |
| **10025** | `TradeRetCodeNoChanges` | **No changes** | New parameters = current parameters |
| **10036** | `TradeRetCodePositionClosed` | **Position closed** | Position doesn't exist |
| **10039** | `TradeRetCodeCloseOrderExist` | **Close order exists** | Cannot close more than position volume |

### 🚨 Rejection Codes

| Code | Go Constant | Description | Reason |
|------|-------------|-------------|--------|
| **10006** | `TradeRetCodeReject` | **Request rejected** | Broker rejected |
| **10007** | `TradeRetCodeCancel` | **Request canceled** | Canceled by trader |
| **10041** | `TradeRetCodeRejectCancel` | **Pending order activation rejected** | Activation rejected and canceled |

---

## Go Helper Functions

The `package/Helpers` (mt5 package) provides helper functions for working with RetCodes:

### IsRetCodeSuccess()

```go
func IsRetCodeSuccess(retCode uint32) bool

// Usage
if mt5.IsRetCodeSuccess(sendData.ReturnedCode) {
    fmt.Println("Trade successful!")
}

// Returns true only for retCode == 10009 (TradeRetCodeDone)
```

### IsRetCodeRequote()

```go
func IsRetCodeRequote(retCode uint32) bool

// Usage
if mt5.IsRetCodeRequote(sendData.ReturnedCode) {
    fmt.Println("Price changed, retrying...")
    // Retry with updated price
}

// Returns true for 10004 or 10020
```

### IsRetCodeRetryable()

```go
func IsRetCodeRetryable(retCode uint32) bool

// Usage
if mt5.IsRetCodeRetryable(sendData.ReturnedCode) {
    // These errors are temporary, retry with delay
    time.Sleep(time.Second * 2)
    // Retry operation
}

// Returns true for: Timeout, NoConnection, Frozen, Locked, TooManyRequests, NoQuotes
```

### GetRetCodeMessage()

```go
func GetRetCodeMessage(retCode uint32) string

// Usage
if sendData.ReturnedCode != mt5.TradeRetCodeDone {
    fmt.Printf("Error: %s\n", mt5.GetRetCodeMessage(sendData.ReturnedCode))
}

// Returns human-readable description for any RetCode
```

---

## Usage Examples

### Example 1: Basic Order Placement with RetCode Check

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"

    mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/google/uuid"
)

func main() {
    // Connection (see GETTING_STARTED.md for connection details; API key from MRPC_API_KEY)
    account, err := mt5.NewMT5Account(12345678, "your_password", "mt5.mrpc.pro:443", uuid.Nil)
    if err != nil {
        fmt.Printf("Setup error: %v\n", err)
        return
    }
    defer account.DisconnectAndClose()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
    defer cancel()

    _, err = account.ConnectEx(ctx, &pb.ConnectExRequest{
        User:          12345678,
        Password:      "your_password",
        MtClusterName: "MetaQuotes-Demo",
    })
    if err != nil {
        fmt.Printf("Connection error: %v\n", err)
        return
    }

    // Place buy market order
    req := &pb.OrderSendRequest{
        Symbol:    "EURUSD",
        Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
        Volume:    0.01,
    }

    sendData, err := account.OrderSend(ctx, req)
    if err != nil {
        var apiErr *mt5.ApiError
        if errors.As(err, &apiErr) {
            fmt.Printf("❌ Order refused: RetCode %d - %s\n",
                apiErr.MqlErrorTradeIntCode(), apiErr.MqlErrorTradeDescription())
            return
        }
        fmt.Printf("gRPC error: %v\n", err)
        return
    }

    // Check RetCode
    if sendData.ReturnedCode == mt5.TradeRetCodeDone {
        fmt.Printf("✅ Order opened successfully!\n")
        fmt.Printf("   Ticket: #%d\n", sendData.Order)
        fmt.Printf("   Volume: %.2f lots\n", sendData.Volume)
        fmt.Printf("   Price: %.5f\n", sendData.Price)
    } else {
        fmt.Printf("❌ Order failed!\n")
        fmt.Printf("   RetCode: %d\n", sendData.ReturnedCode)
        fmt.Printf("   Description: %s\n",
            mt5.GetRetCodeMessage(sendData.ReturnedCode))
    }
}
```

### Example 2: Handling Common Errors (Switch Statement)

```go
import (
    "errors"
    "fmt"
    "time"

    mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
)

ticket, err := sugar.BuyMarket("GBPUSD", 0.5)
if err != nil {
    var apiErr *mt5.ApiError
    if errors.As(err, &apiErr) {
        retCode := uint32(apiErr.MqlErrorTradeIntCode()) // int32 → uint32 to compare with the constants

        switch retCode {
        case mt5.TradeRetCodeNoMoney:
            fmt.Println("⚠️ Insufficient funds!")
            fmt.Println("   Solution: Reduce volume or add margin")

        case mt5.TradeRetCodeMarketClosed:
            fmt.Println("⚠️ Market closed")
            fmt.Println("   Solution: Try during trading hours")

        case mt5.TradeRetCodeInvalidStops:
            fmt.Println("⚠️ SL/TP too close to market price")
            fmt.Println("   Solution: Increase distance (check StopLevel)")

        case mt5.TradeRetCodeInvalidVolume:
            fmt.Println("⚠️ Invalid volume")
            // Get symbol limits
            info, _ := sugar.GetSymbolInfo("GBPUSD")
            fmt.Printf("   Min: %.2f, Max: %.2f, Step: %.2f\n",
                info.VolumeMin, info.VolumeMax, info.VolumeStep)

        case mt5.TradeRetCodeRequote,
             mt5.TradeRetCodePriceChanged:
            fmt.Println("⚠️ Price changed, retrying...")
            time.Sleep(time.Millisecond * 100)
            // Retry order

        default:
            fmt.Printf("❌ Error: %s\n", apiErr.MqlErrorTradeDescription())
        }
    } else {
        fmt.Printf("❌ %v\n", err) // e.g. "order rejected, code: N" from Sugar
    }
} else {
    fmt.Printf("✅ Order #%d opened\n", ticket)
}
```

### Example 3: Retry Logic for Temporary Errors

```go
import (
    "errors"
    "fmt"
    "time"

    sdk "github.com/MetaRPC/GoMT5/examples/mt5"   // MT5Sugar
    mt5 "github.com/MetaRPC/GoMT5/package/Helpers" // ApiError, RetCode helpers
)

func PlaceOrderWithRetry(
    sugar *sdk.MT5Sugar,
    symbol string,
    volume float64,
    maxRetries int,
) (uint64, error) {

    for attempt := 1; attempt <= maxRetries; attempt++ {
        ticket, err := sugar.BuyMarket(symbol, volume)

        // Success
        if err == nil {
            return ticket, nil
        }

        // Check if retryable
        var apiErr *mt5.ApiError
        if errors.As(err, &apiErr) {
            retCode := apiErr.MqlErrorTradeIntCode()

            // Requote - retry immediately
            if mt5.IsRetCodeRequote(uint32(retCode)) {
                fmt.Printf("Requote on attempt %d, retrying...\n", attempt)
                time.Sleep(time.Millisecond * 100)
                continue
            }

            // Retryable error - exponential backoff
            if mt5.IsRetCodeRetryable(uint32(retCode)) {
                waitTime := time.Second * time.Duration(attempt)
                fmt.Printf("Temporary error on attempt %d, waiting %v...\n",
                    attempt, waitTime)
                time.Sleep(waitTime)
                continue
            }

            // Permanent error - stop retrying
            return 0, fmt.Errorf("permanent error: %s",
                apiErr.MqlErrorTradeDescription())
        }

        // Other error (gRPC, network, etc.)
        return 0, err
    }

    return 0, fmt.Errorf("failed after %d attempts", maxRetries)
}

// Usage
ticket, err := PlaceOrderWithRetry(sugar, "EURUSD", 0.01, 3)
if err != nil {
    fmt.Printf("Failed to place order: %v\n", err)
} else {
    fmt.Printf("Order placed: #%d\n", ticket)
}
```

### Example 4: Using helpers.CheckRetCode()

```go
import "github.com/MetaRPC/GoMT5/examples/demos/helpers"

sendData, err := account.OrderSend(ctx, req)
helpers.Fatal(err, "OrderSend failed")

// CheckRetCode prints formatted message and returns true on success
if helpers.CheckRetCode(sendData.ReturnedCode, "Order placement") {
    fmt.Printf("Order ticket: %d\n", sendData.Order)
}

// Output on success:
//   ✓ Order placement successful (RetCode: 10009)

// Output on failure (e.g., market closed):
//   ❌ Order placement failed (RetCode: 10018)
//      Market is closed
//      💡 Hint: Market is closed - check trading hours
```

### Example 5: Modifying SL/TP with Validation

```go
import (
    "context"
    "fmt"

    mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
    pb "github.com/MetaRPC/GoMT5/package"
)

sl, tp := 1.0850, 1.0950
modifyData, err := account.OrderModify(ctx, &pb.OrderModifyRequest{
    Ticket:     123456,
    StopLoss:   &sl,
    TakeProfit: &tp,
})

if err != nil {
    fmt.Printf("Error: %v\n", err) // gRPC error or *mt5.ApiError with the RetCode
    return
}

switch modifyData.ReturnedCode {
case mt5.TradeRetCodeDone:
    fmt.Println("✅ SL/TP updated successfully")

case mt5.TradeRetCodeNoChanges:
    fmt.Println("⚠️ New SL/TP same as current - no action taken")

case mt5.TradeRetCodePositionClosed:
    fmt.Println("⚠️ Position already closed by SL/TP or manually")

case mt5.TradeRetCodeInvalidStops:
    fmt.Println("❌ Invalid SL/TP distance")
    // Get symbol stop level
    stopLevel, _ := sugar.GetMinStopLevel("EURUSD")
    fmt.Printf("   Minimum distance: %d points\n", stopLevel)

default:
    fmt.Printf("❌ Modification failed: %s\n",
        mt5.GetRetCodeMessage(modifyData.ReturnedCode))
}
```

---

## Best Practices

### ✅ DO:

1. **Always check ReturnedCode** after trading operations
   ```go
   if sendData.ReturnedCode == mt5.TradeRetCodeDone {
       // Success
   }
   ```

2. **Use helper functions** for cleaner code
   ```go
   if mt5.IsRetCodeSuccess(retCode) { ... }
   if mt5.IsRetCodeRequote(retCode) { ... }
   ```

3. **Retry requotes** (10004, 10020)
   ```go
   if mt5.IsRetCodeRequote(retCode) {
       // Retry with updated price
   }
   ```

4. **Check margin before trading** to avoid 10019
   ```go
   freeMargin, _ := sugar.GetFreeMargin()
   requiredMargin, _ := sugar.CalculateRequiredMargin("EURUSD", 0.01)
   if freeMargin < requiredMargin {
       fmt.Println("Insufficient margin!")
   }
   ```

### ❌ DON'T:

1. **Don't ignore ReturnedCode** - checking only `err != nil` is NOT enough!
   ```go
   // WRONG:
   sendData, err := account.OrderSend(ctx, req)
   if err == nil {
       // This doesn't mean trade succeeded!
   }

   // CORRECT:
   if err == nil && sendData.ReturnedCode == mt5.TradeRetCodeDone {
       // Now we know trade succeeded
   }
   ```

2. **Don't use magic numbers** - use constants
   ```go
   // WRONG:
   if retCode == 10009 { ... }

   // CORRECT:
   if retCode == mt5.TradeRetCodeDone { ... }
   ```

3. **Don't retry permanent errors** (insufficient margin, market closed)
   ```go
   // Use IsRetCodeRetryable() to check
   ```

4. **Don't assume success = no error**
   ```go
   // Need to check both the error AND RetCode
   if err != nil {
       return err // gRPC/network error or *mt5.ApiError (refused request)
   }
   if sendData.ReturnedCode != 10009 {
       return fmt.Errorf("trade failed: %s",
           mt5.GetRetCodeMessage(sendData.ReturnedCode))
   }
   ```

---

## Quick Reference Table

| Category | RetCodes | Action |
|----------|----------|--------|
| **Success** | 10008, 10009, 10010 | Continue normally |
| **Requote** | 10004, 10020 | Retry immediately |
| **Temporary** | 10012, 10021, 10024, 10028, 10029, 10031 | Retry with delay |
| **Validation** | 10013-10016, 10022, 10030, 10035, 10038 | Fix parameters |
| **Restrictions** | 10017, 10018, 10026, 10027, 10032, 10042-10046 | Check trading conditions |
| **Limits** | 10019, 10033, 10034, 10040 | Reduce volume/positions |
| **State** | 10023, 10025, 10036, 10039 | Check current state |
| **Rejection** | 10006, 10007, 10041 | Check request |
| **Technical** | 10011 | Contact support |

---

## Constants Reference

All RetCode constants are defined in `package/Helpers/errors.go`:

```go
package mt5

const (
    // Success codes
    TradeRetCodeDone        uint32 = 10009 // Request completed successfully
    TradeRetCodeDonePartial uint32 = 10010 // Only part of the request was completed
    TradeRetCodePlaced      uint32 = 10008 // Order placed (pending order activated)

    // Requote codes
    TradeRetCodeRequote      uint32 = 10004 // Requote (price changed, need to retry)
    TradeRetCodePriceChanged uint32 = 10020 // Prices changed (requote)

    // Request rejection codes
    TradeRetCodeReject            uint32 = 10006 // Request rejected
    TradeRetCodeCancel            uint32 = 10007 // Request canceled by trader
    TradeRetCodeInvalidRequest    uint32 = 10013 // Invalid request
    TradeRetCodeInvalidVolume     uint32 = 10014 // Invalid volume in the request
    TradeRetCodeInvalidPrice      uint32 = 10015 // Invalid price in the request
    TradeRetCodeInvalidStops      uint32 = 10016 // Invalid stops in the request (SL/TP too close)
    TradeRetCodeInvalidExpiration uint32 = 10022 // Invalid order expiration date in the request
    TradeRetCodeInvalidFill       uint32 = 10030 // Invalid order filling type
    TradeRetCodeInvalidOrder      uint32 = 10035 // Incorrect or prohibited order type
    TradeRetCodeInvalidCloseVolume uint32 = 10038 // Invalid close volume (exceeds position volume)

    // Trading restriction codes
    TradeRetCodeTradeDisabled    uint32 = 10017 // Trade is disabled
    TradeRetCodeMarketClosed     uint32 = 10018 // Market is closed
    TradeRetCodeServerDisablesAt uint32 = 10026 // Autotrading disabled by server
    TradeRetCodeClientDisablesAt uint32 = 10027 // Autotrading disabled by client terminal
    TradeRetCodeOnlyReal         uint32 = 10032 // Operation is allowed only for live accounts
    TradeRetCodeLongOnly         uint32 = 10042 // Only long positions allowed
    TradeRetCodeShortOnly        uint32 = 10043 // Only short positions allowed
    TradeRetCodeCloseOnly        uint32 = 10044 // Only position close operations allowed
    TradeRetCodeFifoClose        uint32 = 10045 // Position close only by FIFO rule
    TradeRetCodeHedgeProhibited  uint32 = 10046 // Opposite positions on same symbol prohibited (hedging disabled)

    // Resource limit codes
    TradeRetCodeNoMoney        uint32 = 10019 // Not enough money to complete the request (insufficient margin)
    TradeRetCodeLimitOrders    uint32 = 10033 // The number of pending orders has reached the limit
    TradeRetCodeLimitVolume    uint32 = 10034 // The volume of orders and positions for the symbol has reached the limit
    TradeRetCodeLimitPositions uint32 = 10040 // The number of open positions has reached the limit

    // Technical issue codes
    TradeRetCodeError           uint32 = 10011 // Request processing error
    TradeRetCodeTimeout         uint32 = 10012 // Request canceled by timeout
    TradeRetCodeNoQuotes        uint32 = 10021 // No quotes to process the request
    TradeRetCodeTooManyRequests uint32 = 10024 // Too frequent requests
    TradeRetCodeLocked          uint32 = 10028 // Request locked for processing
    TradeRetCodeFrozen          uint32 = 10029 // Order or position frozen
    TradeRetCodeNoConnection    uint32 = 10031 // No connection with the trade server

    // State management codes
    TradeRetCodeOrderChanged    uint32 = 10023 // Order state changed
    TradeRetCodeNoChanges       uint32 = 10025 // No changes in request
    TradeRetCodePositionClosed  uint32 = 10036 // Position with the specified identifier already closed
    TradeRetCodeCloseOrderExist uint32 = 10039 // A close order already exists for a specified position
    TradeRetCodeRejectCancel    uint32 = 10041 // Pending order activation rejected and canceled
)
```
