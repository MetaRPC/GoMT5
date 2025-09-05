# Streaming Quotes (OnSymbolTick) — Proto‑accurate

> **Request:** subscribe to live ticks for one or more symbols. Uses gRPC streaming with two channels: **data** and **errors**.

---

### Code Example

```go
// High-level helper (inside MT5Service):
svc.StreamQuotes(ctx) // by default streams EURUSD, GBPUSD

// Low-level (control your own symbols and packets):
symbols := []string{"EURUSD", "GBPUSD", "XAUUSD"}
ctx2, cancel := context.WithCancel(ctx)
defer cancel()

tickCh, errCh := s.account.OnSymbolTick(ctx2, symbols)
fmt.Println("🔄 Streaming ticks...")
for {
    select {
    case reply, ok := <-tickCh:
        if !ok { fmt.Println("✅ Tick stream ended."); return }
        data := reply.GetOnSymbolTickData() // name may be GetData() in your wrapper
        if data == nil || data.GetSymbolTick() == nil { continue }
        st := data.GetSymbolTick() // MrpcSubscriptionMqlTick
        fmt.Printf("[Tick] %s | Bid: %.5f | Ask: %.5f | Time: %s
",
            st.GetSymbol(), st.GetBid(), st.GetAsk(), st.GetTime().AsTime().Format("2006-01-02 15:04:05"))
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

**Underlying gRPC:** `OnSymbolTick` → `OnSymbolTickReply` → `OnSymbolTickData` with `MrpcSubscriptionMqlTick`.

---

## 🔽 Input

| Parameter | Type              | Required | Description                                                                                                   |
| --------- | ----------------- | -------- | ------------------------------------------------------------------------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Controls lifetime of the stream (cancel/timeout).                                                             |
| `symbols` | `[]string`        | no\*     | Helper uses an internal slice (e.g., `EURUSD`, `GBPUSD`). For a custom list, use the low-level example above. |

> \*In the helper, edit the `symbols := []string{...}` line to change the watch list.

---

## ⬆️ Output

**Tick message type:** `MrpcSubscriptionMqlTick`
Available getters (per your bindings):

| Field        | Type                     | Note                              |
| ------------ | ------------------------ | --------------------------------- |
| `Time`       | `*timestamppb.Timestamp` | server time                       |
| `Bid`        | `float64`                | current bid                       |
| `Ask`        | `float64`                | current ask                       |
| `Last`       | `float64`                | last trade/price, may be 0 for FX |
| `Volume`     | `uint64`                 | tick volume                       |
| `TimeMsc`    | `int64`                  | epoch ms                          |
| `Flags`      | `uint32`                 | exchange flags/bitmask            |
| `VolumeReal` | `float64`                | real volume (if provided)         |
| `Symbol`     | `string`                 | symbol name                       |

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

* **Symbol visibility:** ensure symbols are visible (`EnsureSymbolVisible`) before streaming.
* **Throttle logs:** printing every tick can flood stdout; batch or rate‑limit in production.
* **Stop conditions:** stream ends when `ctx` is canceled, timeout fires, server closes the stream, or an error arrives on `errCh`.
* **Reconnect logic:** for long-running services, wrap with retry/backoff on errors.
* **Scope the list:** watch only the symbols you need to reduce traffic.

---

## ⚠️ Pitfalls

* No trading side-effects, but **high CPU/log I/O** possible if you print every tick.
* On quiet markets you may see long pauses — this is normal.

---

## Variations

* Change list: `symbols := []string{"EURUSD", "XAUUSD"}`.
* Remove demo timeout: drop the `time.After(...)` case to keep the stream open.
* Run several streams in parallel with a shared parent context and `sync.WaitGroup`.
