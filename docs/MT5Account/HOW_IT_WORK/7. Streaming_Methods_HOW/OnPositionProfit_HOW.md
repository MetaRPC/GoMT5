### Example from file: `examples/demos/lowlevel/03_streaming_methods.go`

The **`OnPositionProfit()`** method is used to receive **profit and loss (P&L)** updates for all open positions in real-time.

This stream is particularly useful for:

* displaying current profit on a monitoring dashboard;
* analyzing position dynamics;
* automated risk management.

> The stream is active only when there are open positions and prices change. If positions are absent — there will be no updates.

---

## 🧩 Code example

```go
profitReq := &pb.OnPositionProfitRequest{}

profitChan, profitErrChan := account.OnPositionProfit(ctx, profitReq)

fmt.Printf("Streaming position P&L updates (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
fmt.Println("  NOTE: This stream sends events only when positions exist and prices change")

eventCount = 0
timeout = time.After(MAX_SECONDS * time.Second)

streamProfit:
for {
	select {
	case profitData, ok := <-profitChan:
		if !ok {
			fmt.Println("  Stream closed by server")
			break streamProfit
		}
		eventCount++
		totalPositions := len(profitData.NewPositions) + len(profitData.UpdatedPositions) + len(profitData.DeletedPositions)
		fmt.Printf("  Event #%d: Type=%v Total=%d (New=%d Updated=%d Deleted=%d)\n",
			eventCount, profitData.Type, totalPositions,
			len(profitData.NewPositions),
			len(profitData.UpdatedPositions),
			len(profitData.DeletedPositions))

		// Show first 3 updated positions
		maxShow := 3
		if len(profitData.UpdatedPositions) > 0 {
			if len(profitData.UpdatedPositions) < maxShow {
				maxShow = len(profitData.UpdatedPositions)
			}
			for i := 0; i < maxShow; i++ {
				pos := profitData.UpdatedPositions[i]
				fmt.Printf("    Updated Position: Ticket=%d Symbol=%s Profit=%.2f\n",
					pos.Ticket,
					pos.PositionSymbol,
					pos.Profit)
			}
		}

		if eventCount >= MAX_EVENTS {
			fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
			break streamProfit
		}

	case err := <-profitErrChan:
		if err != nil {
			helpers.PrintShortError(err, "Stream error")
			break streamProfit
		}

	case <-timeout:
		fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
		break streamProfit
	}
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Subscribe to Stream

```go
profitReq := &pb.OnPositionProfitRequest{}
```

An empty request is created. The subscription covers all positions on the account.

---

### 2️. Start the Stream

```go
profitChan, profitErrChan := account.OnPositionProfit(ctx, profitReq)
```

The method returns two channels:

* `profitChan` — position data stream;
* `profitErrChan` — connection error stream.

---

### 3️. Informational Message

```go
fmt.Println("  NOTE: This stream sends events only when positions exist and prices change")
```

Explains that the stream is active **only when there are positions** and quotes change.

---

### 4️. Main Event Reception Loop

```go
for {
	select {
	case profitData, ok := <-profitChan:
```

Uses `select` to simultaneously listen to three channels: events, errors, and timer.

---

### 5️. Process Received Data

```go
totalPositions := len(profitData.NewPositions) + len(profitData.UpdatedPositions) + len(profitData.DeletedPositions)
fmt.Printf("  Event #%d: Type=%v Total=%d (New=%d Updated=%d Deleted=%d)\n", ...)
```

Each event shows how many positions are:

* **new** (NewPositions),
* **updated** (UpdatedPositions),
* **deleted** (DeletedPositions).

---

### 6️. Display Updated Positions

```go
for i := 0; i < maxShow; i++ {
	pos := profitData.UpdatedPositions[i]
	fmt.Printf("    Updated Position: Ticket=%d Symbol=%s Profit=%.2f\n", ...)
}
```

Shows several of the first updated positions with their tickets, symbols, and profit.

---

### 7️. Stream Termination

The stream closes by event limit (`MAX_EVENTS`), timeout (`MAX_SECONDS`), or connection error.

---

## 📦 What the Server Returns

```protobuf
message OnPositionProfitData {
  ENUM_STREAM_EVENT_TYPE Type = 1;
  repeated PositionInfo NewPositions = 2;
  repeated PositionInfo UpdatedPositions = 3;
  repeated PositionInfo DeletedPositions = 4;
}
```

---

## 💡 Example Output

```
Streaming position P&L updates (max 5 events or 10 seconds)...
  NOTE: This stream sends events only when positions exist and prices change
  Event #1: Type=UPDATE Total=1 (New=0 Updated=1 Deleted=0)
    Updated Position: Ticket=123456 Symbol=EURUSD Profit=12.45
  Event #2: Type=UPDATE Total=1 (New=0 Updated=1 Deleted=0)
    Updated Position: Ticket=123456 Symbol=EURUSD Profit=14.10
  ✓ Received 5 events, stopping stream
```

---

## 🧠 What It's Used For

The `OnPositionProfit()` method is used for:

* tracking P&L on positions in real-time;
* displaying profit dynamics in the interface;
* calculating total account results without additional requests;
* building trading panels and analytical tools.

---

## 💬 In Simple Terms

> `OnPositionProfit()` is a stream that reports:
> "EURUSD position is now +12.45, was +10.30 a minute ago".
> It constantly updates profit for all positions, without repeated server requests.
