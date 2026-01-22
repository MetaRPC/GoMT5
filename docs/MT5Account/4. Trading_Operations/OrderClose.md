# ✅ Close Positions and Delete Orders

> **Request:** close an open position at current market price or delete a pending order. Supports partial position closing.

**API Information:**

* **Low-level API:** `MT5Account.OrderClose(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.TradingHelper`
* **Proto definition:** `OrderClose` (defined in `mt5-term-api-trading-helper.proto`)

### RPC

* **Service:** `mt5_term_api.TradingHelper`
* **Method:** `OrderClose(OrderCloseRequest) → OrderCloseReply`
* **Low‑level client (generated):** `TradingHelperClient.OrderClose(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Closes market positions or cancels pending orders.
* **Why you need it.** Exit trades, take profits, cut losses, cancel pending orders.
* **Partial closing.** Supports closing part of a position via Volume parameter.

---

## 🎯 Purpose

Use it to:

* Close market positions at current price
* Cancel pending orders
* Partially close positions
* Exit trades programmatically
* Implement exit strategies

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OrderClose - How it works](../HOW_IT_WORK/4.%20Trading_Operations_HOW/OrderClose_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// OrderClose closes an existing market position or deletes a pending order.
// Supports partial closing via Volume parameter.
func (a *MT5Account) OrderClose(
    ctx context.Context,
    req *pb.OrderCloseRequest,
) (*pb.OrderCloseData, error)
```

**Request message:**

```protobuf
OrderCloseRequest {
  uint64 ticket = 1;      // Order ticket to close (REQUIRED)
  double volume = 2;      // Volume to close (0 = close all)
  int32 slippage = 5;     // Maximum price deviation (slippage) in points
}
```

---

## 🔽 Input

| Parameter | Type                      | Description                                   |
| --------- | ------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`         | Context for deadline/timeout and cancellation |
| `req`     | `*pb.OrderCloseRequest`   | Request with ticket and close parameters      |

**Request fields:**

| Field      | Type     | Required | Description                                      |
| ---------- | -------- | -------- | ------------------------------------------------ |
| `Ticket`   | `uint64` | ✅       | Order ticket number to close                     |
| `Volume`   | `double` | ✅       | Volume to close (0 = close all)                  |
| `Slippage` | `int32`  | ✅       | Maximum price deviation (slippage) in points     |

---

## ⬆️ Output — `OrderCloseData`

| Field                      | Type                    | Description                                           |
| -------------------------- | ----------------------- | ----------------------------------------------------- |
| `ReturnedCode`             | `uint32`                | Operation return code (10009 = success)               |
| `ReturnedStringCode`       | `string`                | String representation of return code                  |
| `ReturnedCodeDescription`  | `string`                | Description of return code                            |
| `CloseMode`                | `MRPC_ORDER_CLOSE_MODE` | Close mode (see enum below)                           |

---

### 📘 Enum: MRPC_ORDER_CLOSE_MODE

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `MRPC_MARKET_ORDER_CLOSE` | Full market position close |
| 1 | `MRPC_MARKET_ORDER_PARTIAL_CLOSE` | Partial market position close |
| 2 | `MRPC_PENDING_ORDER_REMOVE` | Pending order removal |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `30s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Full close:** Set Volume to 0 or full position volume to close entire position.
* **Pending orders:** Deletes pending orders (no market execution).

---

## 🔗 Usage Examples

### 1) Close position completely

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    ticket := uint64(12345678)

    data, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
        Ticket:   ticket,
        Volume:   0,  // Close full position
        Slippage: 20,
    })
    if err != nil {
        panic(err)
    }

    if data.ReturnedCode == 10009 {
        fmt.Printf("Position %d closed successfully\n", ticket)
    } else {
        fmt.Printf("Close failed: %s\n", data.ReturnedCodeDescription)
    }
}
```

### 2) Partial position close

```go
func PartialClose(account *mt5.MT5Account, ticket uint64, closeVolume float64) error {
    ctx := context.Background()

    data, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
        Ticket:   ticket,
        Volume:   closeVolume,
        Slippage: 20,
    })
    if err != nil {
        return fmt.Errorf("failed to close: %w", err)
    }

    if data.ReturnedCode != 10009 {
        return fmt.Errorf("close unsuccessful: %s", data.ReturnedCodeDescription)
    }

    fmt.Printf("Partially closed position %d (%.2f lots)\n", ticket, closeVolume)
    return nil
}

// Usage:
// PartialClose(account, 12345678, 0.05) // Close 0.05 lots
```

### 3) Close all positions

```go
func CloseAllPositions(account *mt5.MT5Account) error {
    ctx := context.Background()

    // Get all open positions
    positions, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{})
    if err != nil {
        return err
    }

    for _, order := range positions.Orders {
        data, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
            Ticket:   order.Ticket,
            Volume:   0, // Full close
            Slippage: 50,
        })
        if err != nil {
            fmt.Printf("Failed to close %d: %v\n", order.Ticket, err)
            continue
        }
        if data.ReturnedCode == 10009 {
            fmt.Printf("Closed position %d\n", order.Ticket)
        } else {
            fmt.Printf("Failed to close %d: %s\n", order.Ticket, data.ReturnedCodeDescription)
        }
        time.Sleep(100 * time.Millisecond) // Small delay between closes
    }

    return nil
}
```

### 4) Close losing positions only

```go
func CloseLosingPositions(account *mt5.MT5Account) error {
    ctx := context.Background()

    positions, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{})
    if err != nil {
        return err
    }

    for _, order := range positions.Orders {
        if order.Profit < 0 { // Only losing positions
            data, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
                Ticket:   order.Ticket,
                Volume:   0,
                Slippage: 30,
            })
            if err != nil {
                fmt.Printf("Failed to close losing position %d: %v\n", order.Ticket, err)
                continue
            }
            if data.ReturnedCode == 10009 {
                fmt.Printf("Closed losing position %d (P&L: %.2f)\n", order.Ticket, order.Profit)
            }
        }
    }

    return nil
}
```

### 5) Close positions for specific symbol

```go
func CloseSymbolPositions(account *mt5.MT5Account, symbol string) error {
    ctx := context.Background()

    positions, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{
        Symbol: symbol,
    })
    if err != nil {
        return err
    }

    closedCount := 0
    for _, order := range positions.Orders {
        data, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
            Ticket:   order.Ticket,
            Volume:   0,
            Slippage: 30,
        })
        if err != nil {
            fmt.Printf("Failed to close %d: %v\n", order.Ticket, err)
            continue
        }
        if data.ReturnedCode == 10009 {
            closedCount++
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Printf("Closed %d %s positions\n", closedCount, symbol)
    return nil
}
```

### 6) Close with confirmation

```go
func ClosePositionWithConfirmation(account *mt5.MT5Account, ticket uint64) error {
    ctx := context.Background()

    // Get position details first
    positions, err := account.OpenedOrders(ctx, &pb.OpenedOrdersRequest{})
    if err != nil {
        return err
    }

    var order *pb.Order
    for _, o := range positions.Orders {
        if o.Ticket == ticket {
            order = o
            break
        }
    }

    if order == nil {
        return fmt.Errorf("position %d not found", ticket)
    }

    fmt.Printf("Closing position %d: %s, Volume: %.2f, P&L: %.2f\n",
        order.Ticket, order.Symbol, order.Volume, order.Profit)

    // Close position
    data, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
        Ticket:   ticket,
        Volume:   0,
        Slippage: 20,
    })
    if err != nil {
        return err
    }

    if data.ReturnedCode != 10009 {
        return fmt.Errorf("close unsuccessful: %s", data.ReturnedCodeDescription)
    }

    fmt.Printf("Position closed successfully\n")

    // Wait and verify it's closed
    time.Sleep(1 * time.Second)
    isOpen, _ := IsTicketOpen(account, ticket)
    if isOpen {
        return fmt.Errorf("position %d still open after close", ticket)
    }

    fmt.Println("Close confirmed")
    return nil
}

func IsTicketOpen(account *mt5.MT5Account, ticket uint64) (bool, error) {
    ctx := context.Background()
    data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})
    if err != nil {
        return false, err
    }

    for _, t := range data.Tickets {
        if t == ticket {
            return true, nil
        }
    }
    return false, nil
}
```

---

## 🔧 Common Patterns

### Close with retry

```go
func CloseWithRetry(account *mt5.MT5Account, ticket uint64, maxRetries int) error {
    ctx := context.Background()

    for i := 0; i < maxRetries; i++ {
        data, err := account.OrderClose(ctx, &pb.OrderCloseRequest{
            Ticket:   ticket,
            Volume:   0,
            Slippage: 50,
        })

        if err == nil && data.ReturnedCode == 10009 {
            fmt.Printf("Position %d closed successfully\n", ticket)
            return nil
        }

        if err != nil {
            fmt.Printf("Close attempt %d failed: %v\n", i+1, err)
        } else {
            fmt.Printf("Close attempt %d failed: %s\n", i+1, data.ReturnedCodeDescription)
        }
        time.Sleep(500 * time.Millisecond)
    }

    return fmt.Errorf("failed to close after %d retries", maxRetries)
}
```

---

## 📚 See Also

* [OrderSend](./OrderSend.md) - Place new orders
* [OrderModify](./OrderModify.md) - Modify existing orders
* [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - Get open positions
* [OpenedOrdersTickets](../3.%20Position_Orders_Information/OpenedOrdersTickets.md) - Get position tickets
