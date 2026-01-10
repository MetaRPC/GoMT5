### Example from file: `examples/demos/lowlevel/03_streaming_methods.go`

> The **`OnPositionsAndPendingOrdersTickets()`** method is used to receive real-time updates of **order and position ticket numbers**.
It reports when orders or positions **appear, disappear, or change**.

> This stream doesn't transmit detailed data — only **ticket number lists**. It's a lightweight way to synchronize the current account state without excessive load.

---

## 🧩 Code example

```go
ticketsReq := &pb.OnPositionsAndPendingOrdersTicketsRequest{}

ticketsChan, ticketsErrChan := account.OnPositionsAndPendingOrdersTickets(ctx, ticketsReq)

fmt.Printf("Streaming ticket changes (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
fmt.Println("  NOTE: This stream sends events when orders/positions open or close")

eventCount = 0
timeout = time.After(MAX_SECONDS * time.Second)

streamTickets:
for {
	select {
	case ticketsData, ok := <-ticketsChan:
		if !ok {
			fmt.Println("  Stream closed by server")
			break streamTickets
		}
		eventCount++
		totalTickets := len(ticketsData.PendingOrderTickets) + len(ticketsData.PositionTickets)
		fmt.Printf("  Event #%d: Total tickets=%d (Pending Orders=%d Positions=%d)\n",
			eventCount,
			totalTickets,
			len(ticketsData.PendingOrderTickets),
			len(ticketsData.PositionTickets))

		if eventCount >= MAX_EVENTS {
			fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
			break streamTickets
		}

	case err := <-ticketsErrChan:
		if err != nil {
			helpers.PrintShortError(err, "Stream error")
			break streamTickets
		}

	case <-timeout:
		fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
		break streamTickets
	}
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Subscribe to Ticket Stream

```go
ticketsReq := &pb.OnPositionsAndPendingOrdersTicketsRequest{}
```

Create a request without parameters — the stream automatically tracks **all tickets** of positions and orders on the current account.

---

### 2️. Start the Stream

```go
ticketsChan, ticketsErrChan := account.OnPositionsAndPendingOrdersTickets(ctx, ticketsReq)
```

The method returns two channels:

* `ticketsChan` — stream of ticket changes;
* `ticketsErrChan` — connection error channel.

---

### 3️. Event Reading Loop

```go
case ticketsData, ok := <-ticketsChan:
```

Each time orders or positions change, a new message arrives. It contains **current ticket lists**.

---

### 4️. Content Analysis

```go
totalTickets := len(ticketsData.PendingOrderTickets) + len(ticketsData.PositionTickets)
fmt.Printf("  Event #%d: Total tickets=%d (Pending Orders=%d Positions=%d)\n", ...)
```

Display the number of tickets:

| Field                  | Description                       |
| --------------------- | --------------------------------- |
| `PendingOrderTickets` | Tickets of all pending orders     |
| `PositionTickets`     | Tickets of all open positions     |

> After each event, the server returns **current state**, not increments — making synchronization easier.

---

### 5️. Stream Termination

The stream closes when:

* event limit is reached (`MAX_EVENTS`),
* time expires (`MAX_SECONDS`),
* connection error occurs.

---

## 📦 What the Server Returns

```protobuf
message OnPositionsAndPendingOrdersTicketsData {
  repeated uint64 PendingOrderTickets = 1;
  repeated uint64 PositionTickets = 2;
}
```

> The server returns only ticket numbers — unique identifiers of orders and positions.

---

## 💡 Example Output

```
Streaming ticket changes (max 5 events or 10 seconds)...
  NOTE: This stream sends events when orders/positions open or close
  Event #1: Total tickets=3 (Pending Orders=1 Positions=2)
  Event #2: Total tickets=2 (Pending Orders=0 Positions=2)
  ✓ Received 5 events, stopping stream
```

---

## 🧠 What It's Used For

The `OnPositionsAndPendingOrdersTickets()` method is used for:

* tracking the appearance and disappearance of orders/positions;
* synchronizing local systems with the trading account;
* building notifications and logs;
* reducing load when monitoring activity.

---

## 💬 In Simple Terms

> `OnPositionsAndPendingOrdersTickets()` is a **lightweight stream** that reports:
> "You currently have 3 tickets: 2 positions and 1 order. One order closed — now 2 tickets".
> It doesn't transmit details, only numbers — ideal for background state tracking.
