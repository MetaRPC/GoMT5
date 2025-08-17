# 📡 Streaming Methods — Overview

These methods provide **real-time streaming** capabilities from MT5 via gRPC:

* Price quotes
* Opened orders profits
* Opened order tickets
* Trade events (new, modified, closed)

> ⚠️ **Important:** Streaming requires a live connection and a properly managed context. Always use `context.WithCancel` or explicit timeouts to avoid leaks.

---

## 🔄 Streaming Quotes

* **File:** [StreamingQuotes.md](StreamingQuotes.md)
* Receive real-time Bid/Ask prices for selected symbols.

---

## 📈 Stream Opened Order Profits

* **File:** [StreamOpenedOrderProfits.md](StreamOpenedOrderProfits.md)
* Get continuous updates on profits/losses for opened positions (new, updated, deleted).

---

## 🎫 Stream Opened Order Tickets

* **File:** [StreamOpenedOrderTickets.md](StreamOpenedOrderTickets.md)
* Stream opened tickets for positions and pending orders.

---

## 🔔 Stream Trade Updates

* **File:** [StreamTradeUpdates.md](StreamTradeUpdates.md)
* Subscribe to trading events — new orders, modifications, and executions.

---

## ⚠️ Things to Watch Out For

* Streams do **not** end by themselves — always cancel or set timeouts.
* Errors arrive on a dedicated `errCh` channel — log and handle them properly.
* Do not keep too many active streams — the terminal and gRPC connection can overload.
* Default timeout in examples: **30 seconds** without events → stream closes.
