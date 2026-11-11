# The Concurrency Model — How GoMT5 Executes

> Unlike SDKs in `async`/`await` languages, GoMT5 has **no async keyword**. Unary calls are ordinary blocking functions; streams hand you **Go channels** you read in a loop or a goroutine. This page explains the model.

---

## 🎯 Two shapes of call

| Shape | What it is | How you consume it |
|-------|-----------|--------------------|
| **Unary (request / response)** | `AccountSummary(ctx, …)`, `SymbolInfoTick(ctx, …)`, `OrderSend(ctx, …)` | call it; it blocks until the reply (or the context deadline) |
| **Server stream** | `OnSymbolTick(ctx, …)`, `MT5Service.StreamTicks(…)` | returns `(<-chan *pb.…Data, <-chan error)` — read frames off the channel |

### Unary

The low-level **MT5Account** and mid-level **MT5Service** take a `context.Context`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

s, err := acc.AccountSummary(ctx, &pb.AccountSummaryRequest{})
if err != nil {
    log.Fatal(err)
}
fmt.Println(s.GetAccountBalance(), s.GetAccountEquity())
```

The high-level **MT5Sugar** manages the context **for you** — its methods don't take a `ctx`:

```go
bid, err := sugar.GetBid("EURUSD")   // Sugar wraps a short per-call timeout internally
```

Either way the call runs on the calling goroutine and returns `(value, error)`. There is nothing to `await`; the gRPC round-trip happens inside the function.

### Streaming

A stream returns two channels — data and error. Read them with `select`, and stop by cancelling the context:

```go
import pb "github.com/MetaRPC/GoMT5/package"

ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

data, errs := acc.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}})
for {
    select {
    case d := <-data:
        t := d.GetSymbolTick()
        fmt.Println(t.GetSymbol(), t.GetBid(), t.GetAsk())
    case err := <-errs:
        if err != nil { return err }
    case <-ctx.Done():
        return nil // deadline/cancel → stop
    }
}
```

The SDK runs the gRPC receive loop on a background goroutine and forwards frames to the channel, retrying transient drops. See **[gRPC Streaming](GRPC_STREAM_MANAGEMENT.md)**.

---

## 🧵 `context.Context` is the control knob

Every low-/mid-level method takes a `ctx`. It carries the **deadline** and the **cancel** signal:

- `context.WithTimeout(parent, d)` — bound a call; it errors with `DeadlineExceeded` if it overruns.
- `context.WithCancel(parent)` — stop a stream (or a batch of calls) on demand by calling `cancel()`.
- When the context carries no deadline, the low-level layer adds a short default per-call timeout: **~3 s** for simple reads (symbol info, account properties), **~5–15 s** for calculations, order checks and history, and **~30 s** for order execution (`OrderSend` / `OrderModify` / `OrderClose`) and connect. Pass your own `ctx` with a deadline to override it.
- **MT5Sugar** supplies its own short timeout per call, so you don't manage one — build your own three-layer client and use `MT5Account` / `MT5Service` directly when you need explicit deadline control.

---

## 🚀 Running work concurrently

Go's concurrency is **goroutines + channels**, not futures. Fan out independent calls yourself:

```go
// read several symbols' prices at once
type res struct{ sym string; bid float64; err error }
out := make(chan res, 3)
for _, s := range []string{"EURUSD", "GBPUSD", "USDJPY"} {
    go func(s string) {
        bid, err := sugar.GetBid(s)
        out <- res{s, bid, err}
    }(s)
}
for i := 0; i < 3; i++ {
    r := <-out
    // handle r
}
```

Or run a management loop in the background while the main flow continues:

```go
go func() {
    for {
        open, err := sugar.HasOpenPosition("EURUSD")
        if err != nil || !open { return } // error or flat → stop
        // … manage positions …
        time.Sleep(2 * time.Second)
    }
}()
```

> Prefer a **batched call** (one `OnSymbolTick` with all symbols) over many single-symbol streams where the API offers one.

---

## 💡 Rules of thumb

* **Pass a `ctx`** with a deadline for `MT5Account` / `MT5Service` unary calls; cancel it to stop streams. `MT5Sugar` handles this for you.
* **One stream, many symbols** beats many single-symbol streams.
* **Issuing RPCs concurrently is fine once connected** — the underlying gRPC calls run in parallel safely. But `MT5Account` has **no internal locking around its connection**, so run `Connect` / `Reconnect` / `Disconnect` from a single goroutine, never while other calls are in flight.
* **Don't block a stream's reader** — do slow work in a separate goroutine so you keep draining the channel.

---

**See also:** [gRPC Streaming](GRPC_STREAM_MANAGEMENT.md) · [MT5Sugar overview](../MT5Sugar/MT5Sugar.API_Overview.md)
