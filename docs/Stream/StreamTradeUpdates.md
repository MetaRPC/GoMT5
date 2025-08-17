# StreamTradeUpdates — Proto‑Accurate Guide

Subscribe to **real‑time trade events**: newly created orders, disappeared orders, state changes, and (optionally) historical updates emitted by the terminal.

> **Proto path:** `OnTrade` → `OnTradeReply` → `OnTradeData`
> `OnTradeData.EventData` type name in bindings: **`OnTadeEventData`** (intentional proto typo)

---

## ✅ Quick Start (High‑Level)

```go
// One‑liner helper
svc.StreamTradeUpdates(ctx)
```

---

## 🔧 Low‑Level Example (Full Control)

```go
// 1) Use a child context to control stream lifetime
ctx2, cancel := context.WithCancel(ctx)
defer cancel()

// 2) Subscribe to trade events
tradeCh, errCh := svc.account.OnTrade(ctx2)
fmt.Println("🔄 Streaming trade updates...")

// 3) Consume replies
for {
    select {
    case reply, ok := <-tradeCh:
        if !ok { fmt.Println("✅ Trade stream ended."); return }

        data := reply.GetOnTradeData() // or GetData() depending on your wrapper
        if data == nil { continue }

        // (a) Realtime event groups (OnTadeEventData)
        if ev := data.GetEventData(); ev != nil {
            for _, o := range ev.GetNewOrders() {
                fmt.Printf("[Trade|NEW]  Ticket:%d Symbol:%s Type:%v Volume:%.2f\n",
                    o.GetTicket(), o.GetSymbol(), o.GetOrderType(), o.GetVolumeCurrent())
            }
            for _, o := range ev.GetDisappearedOrders() {
                fmt.Printf("[Trade|DEL]  Ticket:%d Symbol:%s Type:%v\n",
                    o.GetTicket(), o.GetSymbol(), o.GetOrderType())
            }
            for _, ch := range ev.GetStateChangedOrders() {
                // ch is OnTradeOrderStateChange; print what’s portable
                fmt.Printf("[Trade|CHG]  Ticket:%d NewState:%v\n", ch.GetTicket(), ch.GetState())
            }
        }

        // (b) Optional account snapshot with balance/equity/margin
        if acc := data.GetAccountInfo(); acc != nil {
            _ = acc // Use in dashboards if needed
        }

    case err := <-errCh:
        log.Printf("❌ Stream error: %v", err)
        return

    case <-time.After(30 * time.Second): // demo safety timeout; remove in production
        fmt.Println("⏱️ Timeout reached.")
        return
    }
}
```

---

## 🧾 Method Signature (Helper)

```go
func (s *MT5Service) StreamTradeUpdates(ctx context.Context)
```

**Underlying gRPC (per bindings):**

* Stream: `OnTrade(ctx)`
* Envelope: `OnTradeReply`
* Payload: `OnTradeData` with `EventData` := `OnTadeEventData`

---

## 🔽 Input

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls stream lifetime (cancel/timeout). |

---

## ⬆️ Output (Packet Schema)

**`OnTradeData` fields (commonly used):**

| Field                    | Type                 | Purpose                                                              |
| ------------------------ | -------------------- | -------------------------------------------------------------------- |
| `EventData`              | `OnTadeEventData`    | Realtime order/deal deltas (see below).                              |
| `AccountInfo`            | `OnEventAccountInfo` | Optional account snapshot (Balance/Equity/Margin/FreeMargin/Profit). |
| `TerminalInstanceGuidId` | `string`             | Terminal instance identifier.                                        |

**`OnTadeEventData` key groups:**

| Group                   | Element Type              | Meaning                                                             |
| ----------------------- | ------------------------- | ------------------------------------------------------------------- |
| `NewOrders`             | `OnTradeOrderInfo`        | Orders just created/appeared.                                       |
| `DisappearedOrders`     | `OnTradeOrderInfo`        | Orders removed (filled/canceled/closed).                            |
| `StateChangedOrders`    | `OnTradeOrderStateChange` | Order state transitions.                                            |
| *(plus history deltas)* | deals/orders (various)    | Depending on server settings, you might receive historical updates. |

**Portable fields you can print (based on your service examples):**

* `OnTradeOrderInfo`: `Ticket (int64)`, `Symbol (string)`, `OrderType (enum)`, `VolumeCurrent (float64)`
* `OnTradeOrderStateChange`: `Ticket (int64)`, `State (enum)`

> Exact field set may vary slightly per bindings version; prefer getters you actually see in your generated `*.pb.go`.

---

## 🎯 Purpose

* Live audit of order lifecycle: creation, state changes, disappearance.
* Drive UIs/alerts when orders are filled, canceled, or modified.
* Maintain a push‑based cache of current working orders.

---

## ✅ Best Practices

1. Keep an in‑memory map keyed by **`Ticket (int64)`** for current working orders.
2. On `DisappearedOrders` — evict from cache/UI.
3. Handle `StateChangedOrders` idempotently; same ticket may change state multiple times.
4. Throttle logs — bursts are normal around news or mass order actions.
5. Wrap with reconnect/backoff for long‑running daemons.
6. Combine with **`StreamOpenedOrderProfits`** and **`StreamOpenedOrderTickets`** for a full realtime picture.

---

## ⚠️ Pitfalls

* **Field assumptions:** rely only on getters present in your bindings (check `*.pb.go`).
* **Race conditions:** events can arrive rapidly; update your cache atomically.
* **Ordering:** network jitter may reorder packets; use timestamps/state to reconcile.
* **Flooding:** big activity spikes can produce large bursts — consider buffering/worker pools.

---

## 🔀 Variations

* **Filter by symbol** or order type to reduce noise.
* **Attach account snapshot** to each UI event card (use `AccountInfo`).
* **Correlate with history**: if you also persist closed deals, link `DisappearedOrders` to the deal records.

---

## 🧪 Sample Output

```
🔄 Streaming trade updates...
[Trade|NEW]  Ticket:12345679 Symbol:EURUSD Type:ORDER_TYPE_BUY Volume:0.10
[Trade|CHG]  Ticket:12345679 NewState:ORDER_STATE_PARTIAL
[Trade|DEL]  Ticket:12345679 Symbol:EURUSD Type:ORDER_TYPE_BUY
✅ Trade stream ended.
```
