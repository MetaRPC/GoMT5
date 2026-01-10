# ✅ Get Ticket Numbers of Open Positions

> **Request:** get only ticket numbers of open positions and pending orders without full data. Lightweight alternative to OpenedOrders.

**API Information:**

* **SDK wrapper:** `MT5Account.OpenedOrdersTickets(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `OpenedOrdersTickets` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `OpenedOrdersTickets(OpenedOrdersTicketsRequest) → OpenedOrdersTicketsReply`
* **Low‑level client (generated):** `AccountHelperClient.OpenedOrdersTickets(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

## 💬 Just the essentials

* **What it is.** Returns only ticket numbers of open positions without detailed data.
* **Why you need it.** Lightweight check for position tracking, monitoring, or subsequent operations.
* **Performance.** Much faster than OpenedOrders when you don't need full position details.

---

## 🎯 Purpose

Use it to:

* Track which positions are open
* Monitor position changes efficiently
* Build position tracking systems
* Compare ticket lists between requests
* Lightweight position monitoring

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [OpenedOrdersTickets - How it works](../HOW_IT_WORK/3.%20Position_Orders_Information_HOW/OpenedOrdersTickets_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// OpenedOrdersTickets retrieves only ticket numbers of currently opened orders and positions.
// This is a lightweight alternative when you only need ticket IDs.
func (a *MT5Account) OpenedOrdersTickets(
    ctx context.Context,
    req *pb.OpenedOrdersTicketsRequest,
) (*pb.OpenedOrdersTicketsData, error)
```

**Request message:**

```protobuf
OpenedOrdersTicketsRequest {
  string Symbol = 1;  // Optional symbol filter (empty = all symbols)
  string Group = 2;   // Optional group filter
}
```

**Reply message:**

```protobuf
OpenedOrdersTicketsReply {
  oneof response {
    OpenedOrdersTicketsData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                                | Description                                   |
| --------- | ----------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`                   | Context for deadline/timeout and cancellation |
| `req`     | `*pb.OpenedOrdersTicketsRequest`    | Request with optional Symbol and Group filter |

**Request fields:**

| Field    | Type     | Description                                         |
| -------- | -------- | --------------------------------------------------- |
| `Symbol` | `string` | Optional symbol filter (empty string = all symbols) |
| `Group`  | `string` | Optional group filter                               |

---

## ⬆️ Output — `OpenedOrdersTicketsData`

| Field                    | Type       | Go Type    | Description                         |
| ------------------------ | ---------- | ---------- | ----------------------------------- |
| `OpenedOrdersTickets`    | `uint64[]` | `[]uint64` | Array of pending order tickets      |
| `OpenedPositionTickets`  | `uint64[]` | `[]uint64` | Array of open position tickets      |

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `5s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Symbol filter:** Use Symbol field to get tickets for specific instrument only.
* **Performance:** Much faster than OpenedOrders, use when you only need ticket numbers.

---

## 🔗 Usage Examples

### 1) Get all open position tickets

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
    "github.com/MetaRPC/GoMT5/mt5"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{
        Symbol: "", // All symbols
    })
    if err != nil {
        panic(err)
    }

    totalTickets := len(data.OpenedOrdersTickets) + len(data.OpenedPositionTickets)
    fmt.Printf("Total tickets: %d (Pending: %d, Positions: %d)\n",
        totalTickets, len(data.OpenedOrdersTickets), len(data.OpenedPositionTickets))

    fmt.Printf("\nPosition tickets:\n")
    for _, ticket := range data.OpenedPositionTickets {
        fmt.Printf("  - %d\n", ticket)
    }
}
```

### 2) Get tickets for specific symbol

```go
func GetSymbolTickets(account *mt5.MT5Account, symbol string) ([]uint64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{
        Symbol: symbol,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get tickets: %w", err)
    }

    return data.OpenedPositionTickets, nil
}

// Usage:
// tickets, _ := GetSymbolTickets(account, "EURUSD")
// fmt.Printf("EURUSD tickets: %v\n", tickets)
```

### 3) Monitor new positions

```go
func MonitorNewPositions(account *mt5.MT5Account, interval time.Duration) {
    ctx := context.Background()
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    var lastTickets []uint64

    for range ticker.C {
        data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        currentTickets := data.OpenedPositionTickets

        if lastTickets != nil {
            // Find new tickets
            newTickets := findNewTickets(lastTickets, currentTickets)
            for _, ticket := range newTickets {
                fmt.Printf("[%s] New position opened: %d\n",
                    time.Now().Format("15:04:05"),
                    ticket,
                )
            }

            // Find closed tickets
            closedTickets := findNewTickets(currentTickets, lastTickets)
            for _, ticket := range closedTickets {
                fmt.Printf("[%s] Position closed: %d\n",
                    time.Now().Format("15:04:05"),
                    ticket,
                )
            }
        }

        lastTickets = currentTickets
    }
}

func findNewTickets(old, new []uint64) []uint64 {
    oldSet := make(map[uint64]bool)
    for _, t := range old {
        oldSet[t] = true
    }

    newTickets := []uint64{}
    for _, t := range new {
        if !oldSet[t] {
            newTickets = append(newTickets, t)
        }
    }
    return newTickets
}
```

### 4) Check if specific ticket is still open

```go
func IsTicketOpen(account *mt5.MT5Account, ticket uint64) (bool, error) {
    ctx := context.Background()

    data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})
    if err != nil {
        return false, err
    }

    for _, t := range data.OpenedPositionTickets {
        if t == ticket {
            return true, nil
        }
    }

    return false, nil
}

// Usage:
// isOpen, _ := IsTicketOpen(account, 12345678)
// if isOpen {
//     fmt.Println("Position is still open")
// }
```

### 5) Quick position count

```go
func QuickPositionCount(account *mt5.MT5Account, symbol string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{
        Symbol: symbol,
    })
    if err != nil {
        return 0, err
    }

    return len(data.OpenedPositionTickets), nil
}
```

### 6) Compare ticket sets

```go
func CompareTicketLists(account *mt5.MT5Account) {
    ctx := context.Background()

    // Get tickets at two different times
    data1, _ := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})
    time.Sleep(5 * time.Second)
    data2, _ := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})

    added := findNewTickets(data1.OpenedPositionTickets, data2.OpenedPositionTickets)
    removed := findNewTickets(data2.OpenedPositionTickets, data1.OpenedPositionTickets)

    fmt.Printf("Added positions: %v\n", added)
    fmt.Printf("Removed positions: %v\n", removed)
    fmt.Printf("Unchanged: %d positions\n", len(data1.OpenedPositionTickets)-len(removed))
}
```

---

## 🔧 Common Patterns

### Fast position existence check

```go
func HasAnyPositions(account *mt5.MT5Account) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    data, err := account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})
    if err != nil {
        return false
    }

    return len(data.OpenedPositionTickets) > 0
}
```

### Track ticket changes

```go
type TicketMonitor struct {
    account     *mt5.MT5Account
    lastTickets map[uint64]bool
}

func (tm *TicketMonitor) CheckChanges() (added, removed []uint64, err error) {
    ctx := context.Background()

    data, err := tm.account.OpenedOrdersTickets(ctx, &pb.OpenedOrdersTicketsRequest{})
    if err != nil {
        return nil, nil, err
    }

    currentTickets := make(map[uint64]bool)
    for _, t := range data.OpenedPositionTickets {
        currentTickets[t] = true

        // Check if new
        if tm.lastTickets != nil && !tm.lastTickets[t] {
            added = append(added, t)
        }
    }

    // Check for removed
    if tm.lastTickets != nil {
        for t := range tm.lastTickets {
            if !currentTickets[t] {
                removed = append(removed, t)
            }
        }
    }

    tm.lastTickets = currentTickets
    return added, removed, nil
}
```

---

## 📚 See Also

* [OpenedOrders](./OpenedOrders.md) - Get full position details
* [PositionsTotal](./PositionsTotal.md) - Get quick position count
* [OnPositionsAndPendingOrdersTickets](../7.%20Streaming_Methods/OnPositionsAndPendingOrdersTickets.md) - Stream ticket changes
* [OrderClose](../4.%20Trading_Operations/OrderClose.md) - Close positions by ticket
