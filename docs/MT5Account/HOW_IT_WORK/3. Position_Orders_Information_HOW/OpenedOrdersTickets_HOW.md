### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`OpenedOrdersTickets()`** method is used to get **only the ticket numbers** of all active orders and positions on the account.
> This is a lightweight version of the `OpenedOrders()` method, which returns a minimal data set without extra information about prices, volumes, and symbols.
>
> This call is useful when you simply need to know **which trades exist**, and then optionally get details for each ticket through separate requests.


---

## 🧩 Code example

```go
fmt.Println("\n5.3. OpenedOrdersTickets() - Get ticket numbers only")

ticketsReq := &pb.OpenedOrdersTicketsRequest{}
ticketsData, err := account.OpenedOrdersTickets(ctx, ticketsReq)
if err != nil {
    helpers.PrintShortError(err, "OpenedOrdersTickets failed")
} else {
    // Direct field access: OpenedOrdersTickets and OpenedPositionTickets
    totalTickets := len(ticketsData.OpenedOrdersTickets) + len(ticketsData.OpenedPositionTickets)
    fmt.Printf("  Total tickets:                 %d\n", totalTickets)
    fmt.Printf("    Pending order tickets:       %d\n", len(ticketsData.OpenedOrdersTickets))
    fmt.Printf("    Position tickets:            %d\n", len(ticketsData.OpenedPositionTickets))

    // Show first few position tickets if any exist
    if len(ticketsData.OpenedPositionTickets) > 0 {
        fmt.Printf("  First position tickets:        ")
        maxShow := 5
        if len(ticketsData.OpenedPositionTickets) < maxShow {
            maxShow = len(ticketsData.OpenedPositionTickets)
        }
        for i := 0; i < maxShow; i++ {
            fmt.Printf("%d ", ticketsData.OpenedPositionTickets[i])
        }
        fmt.Println()
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
ticketsReq := &pb.OpenedOrdersTicketsRequest{}
```

An empty request is created. The method does not require parameters, as it simply requests **all active position and order tickets**.

---

```go
ticketsData, err := account.OpenedOrdersTickets(ctx, ticketsReq)
```

A gRPC request `OpenedOrdersTickets()` is sent to the MetaTrader server.
In response, a structure is returned with two lists:

* **`OpenedOrdersTickets`** — tickets of all pending orders;
* **`OpenedPositionTickets`** — tickets of all open positions.

---

```go
totalTickets := len(ticketsData.OpenedOrdersTickets) + len(ticketsData.OpenedPositionTickets)
fmt.Printf("  Total tickets:                 %d\n", totalTickets)
fmt.Printf("    Pending order tickets:       %d\n", len(ticketsData.OpenedOrdersTickets))
fmt.Printf("    Position tickets:            %d\n", len(ticketsData.OpenedPositionTickets))
```

The total number of tickets is calculated and their distribution by type is displayed:

* **Pending orders** — pending orders that have not yet been executed;
* **Positions** — open market positions.

Example output:

```
Total tickets:                 4
  Pending order tickets:       1
  Position tickets:            3
```

---

```go
if len(ticketsData.OpenedPositionTickets) > 0 {
    fmt.Printf("  First position tickets:        ")
    maxShow := 5
    if len(ticketsData.OpenedPositionTickets) < maxShow {
        maxShow = len(ticketsData.OpenedPositionTickets)
    }
    for i := 0; i < maxShow; i++ {
        fmt.Printf("%d ", ticketsData.OpenedPositionTickets[i])
    }
    fmt.Println()
}
```

If there are active positions on the account, the first few tickets are displayed (up to 5 by default).
This is convenient for quick visual control.

Example output:

```
First position tickets:        1234567 1234568 1234569
```

---

## 📦 What the Server Returns

```protobuf
message OpenedOrdersTicketsData {
  repeated int64 opened_orders_tickets = 1;     // Pending order tickets
  repeated int64 opened_position_tickets = 2;   // Open position tickets
}
```

> **📌 Note:** For a complete description of all structure fields, see the main documentation file [OpenedOrdersTickets.md](../../3.%20Position_Orders_Information/OpenedOrdersTickets.md).

---

## 💡 Example Output

```
Total tickets:                 3
  Pending order tickets:       0
  Position tickets:            3
  First position tickets:      5674321 5674322 5674323
```

📘 This means that 3 active positions are currently open on the account, and their tickets can be used for further queries (`PositionGet`, `OrderGet`, etc.).

---

### 🧠 What It's Used For

The `OpenedOrdersTickets()` method is used:

* for **quickly getting a list of identifiers** of active orders and positions;
* to **optimize performance** when full data is not needed;
* in **asynchronous processing**, where tickets are then processed by separate threads;
* in **indicators or alert systems** where only the presence of positions matters;
* during **robot synchronization with the server** to reconcile tickets without loading full information.

---

### 💬 In Simple Terms

> `OpenedOrdersTickets()` is a lightweight way to find out **which trades are currently open**,
> without requesting extra data. Returns only tickets — identifiers of positions and orders.
> Fast, simple, and ideal for background checks or mass synchronization.
