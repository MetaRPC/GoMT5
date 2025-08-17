# Streaming Quotes (OnSymbolTick)

> **Request:** subscribe to live ticks for one or more symbols. Uses gRPC streaming with two channels: **data** and **errors**.

---

### Code Example

```go
// High-level helper (inside MT5Service):
svc.StreamQuotes(ctx) // by default streams EURUSD, GBPUSD

// Low-level (pass your own symbols and handle packets):
symbols := []string{"EURUSD", "GBPUSD", "XAUUSD"}
ctx2, cancel := context.WithCancel(ctx)
defer cancel()

tickCh, errCh := s.account.OnSymbolTick(ctx2, symbols)
fmt.Println("🔄 Streaming ticks...")
for {
    select {
    case pkt, ok := <-tickCh:
        if !ok { fmt.Println("✅ Tick stream ended."); return }
        if st := pkt.GetSymbolTick(); st != nil {
            fmt.Printf("[Tick] %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
                st.GetSymbol(), st.GetBid(), st.GetAsk(), st.GetTime().AsTime().Format("2006-01-02 15:04:05"))
        }
    case err := <-errCh:
        log.Printf("❌ Stream error: %v", err)
        return
    case <-time.After(30 * time.Second): // safety timeout for demos
        fmt.Println("⏱️ Timeout reached.")
        return
    }
}
```

---

### Method Signature (helper)

```go
func (s *MT5Service) StreamQuotes(ctx context.Context)
```

**Under the hood:** calls `OnSymbolTick(ctx, symbols []string)` on the account client and forwards packets from its channels.

---

## 🔽 Input

| Parameter | Type              | Required | Description                                                                                                      |
| --------- | ----------------- | -------- | ---------------------------------------------------------------------------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Controls lifetime of the stream (cancel/timeout).                                                                |
| `symbols` | `[]string`        | no\*     | Helper uses a built-in slice (e.g., `EURUSD`, `GBPUSD`). For a custom list, use the **low-level** example above. |

> \*In the helper, edit the `symbols := []string{...}` line to change the watch list.

---

## ⬆️ Output

Continuous packets on channels:

* **`tickCh`** (`<-chan *pb.SymbolTickPacket`)

  * `Symbol` — e.g., `EURUSD`
  * `Bid`, `Ask` — latest prices
  * `Time` — server timestamp (`google.protobuf.Timestamp`)
* **`errCh`** (`<-chan error`) — transport/stream errors

**Sample console output:**

```
[Tick] EURUSD | Bid: 1.09876 | Ask: 1.09889 | Time: 2025-08-17 12:00:01
[Tick] GBPUSD | Bid: 1.28543 | Ask: 1.28558 | Time: 2025-08-17 12:00:01
```

---

## 🎯 Purpose

* Feed live prices into strategies, dashboards, or alerting.
* Validate connectivity and symbol availability in real time.

---

## 🧩 Notes & Tips

* **Symbol visibility:** make sure symbols are visible (`EnsureSymbolVisible`) before streaming.
* **Throttle logs:** printing every tick can flood stdout; batch or rate-limit in production.
* **Stop conditions:** stream ends when `ctx` is canceled, timeout fires, server closes the stream, or an error arrives on `errCh`.
* **Reconnect logic:** for long-running services, wrap with retry/backoff on errors.
* **Scope the list:** watch only the symbols you need to reduce traffic.

---

## ⚠️ Pitfalls

* No trading side-effects, but **high CPU/log I/O** possible if you print every tick.
* On quiet markets you may see long pauses — это нормально.

---

## Variations

* Change list: `symbols := []string{"EURUSD", "XAUUSD"}`.
* Remove demo timeout: drop the `time.After(...)` case to keep the stream open.
* Run several streams in parallel with a shared parent context and `sync.WaitGroup`.
