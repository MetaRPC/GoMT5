# StreamOpenedOrderTickets — Proto-Accurate Guide

This streaming method provides a **live feed of all open order tickets** (both positions and pending orders). It is useful for continuously tracking what is currently open in the account without polling.

---

## 📌 Method Signature

```go
func (s *MT5Service) StreamOpenedOrderTickets(ctx context.Context)
```

Internally calls:

```go
ticketCh, errCh := s.account.OnOpenedOrdersTickets(ctx, depth)
```

Where:

* `depth` → maximum number of tickets to fetch in one update (e.g., 1000).

---

## 📩 Streamed Messages

The channel delivers objects of type:

```proto
message OpenedOrdersTicketsPacket {
  repeated uint64 position_tickets     = 1; // active positions
  repeated uint64 pending_order_tickets = 2; // active pending orders
}
```

In Go binding:

```go
pkt.GetPositionTickets()       // []uint64
pkt.GetPendingOrderTickets()   // []uint64
```

Combined usage:

```go
tix := append(pkt.GetPositionTickets(), pkt.GetPendingOrderTickets()...)
```

---

## 🔧 Example Usage

```go
ctx2, cancel := context.WithCancel(ctx)
defer cancel()

ticketCh, errCh := s.account.OnOpenedOrdersTickets(ctx2, 1000)
fmt.Println("🔄 Streaming opened order tickets...")

for {
    select {
    case pkt, ok := <-ticketCh:
        if !ok {
            fmt.Println("✅ Ticket stream ended.")
            return
        }
        tix := append(pkt.GetPositionTickets(), pkt.GetPendingOrderTickets()...)
        fmt.Printf("[Tickets] %d open tickets: %v\n", len(tix), tix)

    case err := <-errCh:
        log.Printf("❌ Stream error: %v", err)
        return

    case <-time.After(30 * time.Second):
        fmt.Println("⏱️ Timeout reached.")
        return
    }
}
```

---

## ✅ Best Practices

1. **Depth parameter**: keep reasonably high (e.g., 1000) unless broker imposes stricter limits.
2. **Always combine** `PositionTickets` + `PendingOrderTickets` if you want the full picture.
3. Use **context cancellation** to stop streaming when shutting down gracefully.
4. Expect **frequent updates** even when no new orders are opened — MT5 terminal may push heartbeats.

---

## ⚠️ Pitfalls

* **Duplicate tickets**: You may see the same ticket in multiple updates — always de-duplicate if building persistent state.
* **Ticket gaps**: Brokers may reuse ticket IDs after a long time — don’t assume monotonic increments.
* **Depth limit**: If open orders exceed your `depth`, some will be dropped from the stream. Handle accordingly.
* **Stream ending**: Channel closure means either disconnect or terminal shutdown — always be ready to reconnect.

---

## 🎯 Use Cases

* Monitor in real-time how many positions/pending orders are open.
* Feed into a **dashboard** showing open exposure.
* Use as a **trigger source**: e.g., if no open tickets remain, flatten internal strategy state.

---

👉 This stream is purely **read-only**. No orders are created/modified — safe to run continuously.
