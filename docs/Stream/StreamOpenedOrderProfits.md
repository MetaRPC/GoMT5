# Streaming Opened Order Profits — Proto‑Accurate Guide

> **Goal:** Subscribe to **real‑time P/L updates** for open positions/orders without polling.
> **Proto path:** `OnPositionProfit` → `OnPositionProfitReply` → `OnPositionProfitData`

---

## ✅ Quick Start (High‑Level)

```go
// Recommended: use the service helper
svc.StreamOpenedOrderProfits(ctx)
```

---

## 🔧 Low‑Level Example (Full Control)

```go
// 1) Create a cancellable child context for the stream lifetime
ctx2, cancel := context.WithCancel(ctx)
defer cancel()

// 2) Subscribe (wrapper around OnPositionProfit)
replyCh, errCh := svc.account.OnOpenedOrdersProfit(ctx2, 1000) // buffer=1000 is a safe demo default
fmt.Println("🔄 Streaming order profits...")

// 3) Consume packets
for {
    select {
    case reply, ok := <-replyCh:
        if !ok { fmt.Println("✅ Profit stream ended."); return }

        data := reply.GetData() // OnPositionProfitData
        if data == nil { continue }

        // The stream delivers Δ-changes split into 3 groups
        for _, info := range data.GetNewPositions() {
            fmt.Printf("[Profit|NEW] Ticket:%d | Symbol:%s | Profit:%.2f\n",
                info.GetTicket(), info.GetPositionSymbol(), info.GetProfit())
        }
        for _, info := range data.GetUpdatedPositions() {
            fmt.Printf("[Profit|UPD] Ticket:%d | Symbol:%s | Profit:%.2f\n",
                info.GetTicket(), info.GetPositionSymbol(), info.GetProfit())
        }
        for _, info := range data.GetDeletedPositions() {
            fmt.Printf("[Profit|DEL] Ticket:%d | Symbol:%s | Profit:%.2f\n",
                info.GetTicket(), info.GetPositionSymbol(), info.GetProfit())
        }

        if acc := data.GetAccountInfo(); acc != nil {
            // Optional account snapshot: Balance, Equity, Margin, FreeMargin, Profit, etc.
            _ = acc // use in dashboards if needed
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
func (s *MT5Service) StreamOpenedOrderProfits(ctx context.Context)
```

**Underlying gRPC (per bindings):**

* Request/stream: `OnPositionProfit(ctx)`
* Envelope: `OnPositionProfitReply`
* Payload: `OnPositionProfitData`

Your wrapper may expose:
`OnOpenedOrdersProfit(ctx context.Context, buffer int) (<-chan *pb.OnPositionProfitReply, <-chan error)`

---

## 🔽 Input

| Parameter | Type              | Required | Description                                    |
| --------- | ----------------- | -------- | ---------------------------------------------- |
| `ctx`     | `context.Context` | yes      | Controls stream lifetime (cancel/timeout).     |
| `buffer`  | `int`             | no       | Suggested channel buffer; examples use `1000`. |

---

## ⬆️ Output (Packet Schema)

**`OnPositionProfitData` fields:**

| Field                    | Type                              | Purpose                                                              |
| ------------------------ | --------------------------------- | -------------------------------------------------------------------- |
| `NewPositions`           | `[]*OnPositionProfitPositionInfo` | Positions that **appeared** since last packet.                       |
| `UpdatedPositions`       | `[]*OnPositionProfitPositionInfo` | Existing positions whose profit/fields **changed**.                  |
| `DeletedPositions`       | `[]*OnPositionProfitPositionInfo` | Positions that were **closed/removed**.                              |
| `AccountInfo`            | `*OnEventAccountInfo`             | Optional account snapshot (Balance/Equity/Margin/FreeMargin/Profit). |
| `TerminalInstanceGuidId` | `string`                          | Terminal instance identifier.                                        |

**`OnPositionProfitPositionInfo` (exact fields):**

* `Index` (`int32`)
* `Ticket` (`int64`)  ← **use `int64` as the cache key**
* `Profit` (`float64`)
* `PositionSymbol` (`string`)

---

## 🎯 Purpose

* Build **live P/L widgets** per position or per symbol.
* Trigger alerts/rules (e.g., TP runner, trailing, drawdown guard).
* Maintain a push‑based **cache** of open exposure without polling.

---

## ✅ Best Practices

1. Keep an in‑memory map: `map[int64]OnPositionProfitPositionInfo` keyed by `Ticket`.
2. Apply `DeletedPositions` to evict finished tickets from the cache/UI.
3. Throttle or aggregate logs — profit updates приходят часто.
4. Use a parent context; stop streams on shutdown (Ctrl+C) or service restarts.
5. Implement **reconnect/backoff** in long‑running daemons.

---

## ⚠️ Pitfalls

* **Backpressure:** heavy work in the select loop blocks the channel → offload to workers.
* **UI spam:** printing every update floods stdout; sample, batch, or rate‑limit.
* **Field assumptions:** rely only on fields confirmed in your bindings.

---

## 🔀 Variations

* **Always‑on service:** remove demo timeout; wrap in a retry loop with backoff.
* **Filtering:** ignore updates for symbols not in your watchlist.
* **Alerting:** threshold checks on `Profit` (e.g., > X or < –Y) → notify or act.

---

## 🧪 Sample Output

```
🔄 Streaming order profits...
[Profit|NEW] Ticket: 12345678 | Symbol: EURUSD | Profit: 1.25
[Profit|UPD] Ticket: 12345678 | Symbol: EURUSD | Profit: 3.90
[Profit|DEL] Ticket: 12345678 | Symbol: EURUSD | Profit: 4.10
✅ Profit stream ended.
```
