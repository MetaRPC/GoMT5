### Example from file: `examples/demos/lowlevel/03_streaming_methods.go`

> The **`OnSymbolTick()`** method is used to subscribe to a **real-time tick data stream** (Bid/Ask).
> It opens a continuous gRPC stream between the client and the MetaTrader server, sending updates on every price change.
>
> This is analogous to the `OnTick()` function in MQL5, but implemented as a persistent asynchronous stream.


---

## 🧩 Code example

```go
tickReq := &pb.OnSymbolTickRequest{
    SymbolNames: []string{cfg.TestSymbol}, // Symbols to stream (one or more)
}

tickCtx, tickCancel := context.WithCancel(ctx)
tickChan, tickErrChan := account.OnSymbolTick(tickCtx, tickReq)

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
        // The tick itself is nested: SymbolTick (nil in a message without a tick)
        if tickData.SymbolTick != nil {
            tick := tickData.SymbolTick
            fmt.Printf("  Event #%d: Symbol=%s Bid=%.5f Ask=%.5f Spread=%.5f\n",
                eventCount,
                tick.Symbol,
                tick.Bid,
                tick.Ask,
                tick.Ask-tick.Bid)
        }

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
tickCancel() // closes this stream on the server
```

---

## 🟢 Detailed Code Explanation

### 1️. Subscribe to Tick Data Stream

```go
tickReq := &pb.OnSymbolTickRequest{
    SymbolNames: []string{cfg.TestSymbol},
}
```

Create a request to stream ticks for the symbol `cfg.TestSymbol` (e.g., `EURUSD`).

> 🔸 Note: `SymbolNames` is a list: one stream can carry several symbols (each tick says its `Symbol`).

---

### 2️. Start the Stream

```go
tickCtx, tickCancel := context.WithCancel(ctx)
tickChan, tickErrChan := account.OnSymbolTick(tickCtx, tickReq)
```

The stream gets its own context, derived from `ctx`: `tickCancel()` after the loop closes this stream on the server at once. Without it the stream would stay open, unread, until `ctx` itself ends — and every later stream of the demo would pile up next to it.

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
if tickData.SymbolTick != nil {
    tick := tickData.SymbolTick
    fmt.Printf("  Event #%d: Symbol=%s Bid=%.5f Ask=%.5f Spread=%.5f\n",
        eventCount,
        tick.Symbol,
        tick.Bid,
        tick.Ask,
        tick.Ask-tick.Bid)
}
```

Each tick (`tickData.SymbolTick`) contains:

* **Symbol** — symbol of the tick;
* **Bid** — buy price;
* **Ask** — sell price;
* **Spread** — not a field: Ask − Bid, computed here (in price units, not points);
* **Time** / **TimeUtc** — tick timestamp (broker server time / UTC);
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
  MrpcSubscriptionMqlTick symbol_tick = 1;
  string terminal_instance_guid_id = 2;
}

message MrpcSubscriptionMqlTick {
  google.protobuf.Timestamp time = 1;       // Tick time (broker server time)
  google.protobuf.Timestamp time_utc = 10;  // Tick time in UTC
  double bid = 2;                           // Current bid price
  double ask = 3;                           // Current ask price
  double last = 4;                          // Last deal price
  uint64 volume = 5;                        // Volume
  int64 time_msc = 6;                       // Tick time in milliseconds (broker server time)
  int64 time_msc_utc = 11;                  // time_msc in UTC
  uint32 flags = 7;                         // Tick flags
  double volume_real = 8;                   // Volume with greater accuracy
  string symbol = 9;                        // Symbol
}
```

---

## 💡 Example Output

```
Streaming EURUSD tick data (max 10 events or 5 seconds)...
  Event #1: Symbol=EURUSD Bid=1.08542 Ask=1.08558 Spread=0.00016
  Event #2: Symbol=EURUSD Bid=1.08540 Ask=1.08555 Spread=0.00015
  Event #3: Symbol=EURUSD Bid=1.08539 Ask=1.08554 Spread=0.00015
  ⏱ Timeout after 5 seconds, stopping stream
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
