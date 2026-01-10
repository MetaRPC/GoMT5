### Example from file: `examples/demos/lowlevel/03_streaming_methods.go`

> The **`OnSymbolTick()`** method is used to subscribe to a **real-time tick data stream** (Bid/Ask).
> It opens a continuous gRPC stream between the client and the MetaTrader server, sending updates on every price change.
>
> This is analogous to the `OnTick()` function in MQL5, but implemented as a persistent asynchronous stream.


---

## 🧩 Code example

```go
tickReq := &pb.OnSymbolTickRequest{
    Symbol: cfg.TestSymbol, // Single symbol to stream
}

tickChan, tickErrChan := account.OnSymbolTick(ctx, tickReq)

fmt.Printf("Streaming %s tick data (max %d events or %d seconds)...\n", cfg.TestSymbol, MAX_EVENTS, MAX_SECONDS)

eventCount := 0
timeout := time.After(MAX_SECONDS * time.Second)

streamTick:
for {
    select {
    case tickData, ok := <-tickChan:
        if !ok {
            fmt.Println("  Stream closed by server")
            break streamTick
        }
        eventCount++
        fmt.Printf("  Event #%d: Bid=%.5f Ask=%.5f Spread=%.5f\n",
            eventCount,
            tickData.Bid,
            tickData.Ask,
            tickData.Ask-tickData.Bid)

        if eventCount >= MAX_EVENTS {
            fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
            break streamTick
        }

    case err := <-tickErrChan:
        if err != nil {
            helpers.PrintShortError(err, "Stream error")
            break streamTick
        }

    case <-timeout:
        fmt.Printf("  ⏱ Timeout after %d seconds, stopping stream\n", MAX_SECONDS)
        break streamTick
    }
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Subscribe to Tick Data Stream

```go
tickReq := &pb.OnSymbolTickRequest{
    Symbol: cfg.TestSymbol,
}
```

Create a request to stream ticks for the symbol `cfg.TestSymbol` (e.g., `EURUSD`).

> 🔸 Note: To stream multiple symbols, create separate streams for each symbol.

---

### 2️. Start the Stream

```go
tickChan, tickErrChan := account.OnSymbolTick(ctx, tickReq)
```

The method returns two channels:

* **tickChan** — data channel for tick updates;
* **tickErrChan** — error channel if the stream is interrupted.

> The stream works asynchronously — new ticks arrive continuously until timeout or event limit is reached.

---

### 3️. Set Up Limits

```go
eventCount := 0
timeout := time.After(MAX_SECONDS * time.Second)
```

* **MAX_EVENTS** — how many ticks to receive (e.g., 10);
* **MAX_SECONDS** — time limit (e.g., 5 seconds).

> This is a safeguard to prevent the example from running indefinitely.

---

### 4️. Main Data Reception Loop

```go
for {
    select {
    case tickData, ok := <-tickChan:
        ...
    case err := <-tickErrChan:
        ...
    case <-timeout:
        ...
    }
}
```

Uses the `select` construct to listen to three sources simultaneously:

1. **tickChan** — new tick arrived;
2. **tickErrChan** — stream error;
3. **timeout** — time expired.

---

### 5️. Process Incoming Tick

```go
fmt.Printf("  Event #%d: Bid=%.5f Ask=%.5f Spread=%.5f\n",
    eventCount,
    tickData.Bid,
    tickData.Ask,
    tickData.Ask-tickData.Bid)
```

Each tick contains:

* **Bid** — buy price;
* **Ask** — sell price;
* **Spread** — difference between Ask and Bid;
* **Time** — tick timestamp;
* **Volume** — tick volume;
* **Last** — last deal price.

> This data is needed for strategy calculations, scalping, spread analysis, and market visualization.

---

### 6️. Stream Termination by Limit or Timeout

```go
if eventCount >= MAX_EVENTS {
    fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
    break streamTick
}
```

The stream stops either by event count (`MAX_EVENTS`) or when `timeout` expires.

---

## 📦 What the Server Returns

```protobuf
message OnSymbolTickData {
  double Bid = 1;     // Current bid price
  double Ask = 2;     // Current ask price
  double Last = 3;    // Last deal price
  uint64 Volume = 4;  // Volume
  google.protobuf.Timestamp Time = 5; // Tick timestamp
  uint32 Flags = 6;   // Tick flags
}
```

---

## 💡 Example Output

```
Streaming EURUSD tick data (max 5 events or 10 seconds)...
  Event #1: Bid=1.08542 Ask=1.08558 Spread=0.00016
  Event #2: Bid=1.08540 Ask=1.08555 Spread=0.00015
  Event #3: Bid=1.08539 Ask=1.08554 Spread=0.00015
  ✓ Received 3 events, stopping stream
```

---

## 🧠 What It's Used For

The `OnSymbolTick()` method is used for:

* receiving real-time quotes;
* building streaming strategies (scalping, market making);
* monitoring broker latency and market activity;
* updating visualization or statistics in real-time.

---

## 💬 In Simple Terms

> `OnSymbolTick()` is a **live price subscription**.
> The gateway server sends you every Bid/Ask change, and you can process it instantly.
> For example, calculate spread, build charts, or react to market changes — all in real-time.
