# MT5Account · Streaming Methods - Overview

> Real-time streams: ticks, trades, profit updates, transaction log. Use this page to choose the right API for real-time subscriptions.

## 📁 What lives here

### Real-Time Price Updates

* **[OnSymbolTick](./OnSymbolTick.md)** - real-time tick stream for symbols.

### Trading Events

* **[OnTrade](./OnTrade.md)** - position/order changes (opened, closed, modified).
* **[OnTradeTransaction](./OnTradeTransaction.md)** - detailed transaction log (complete audit trail).

### Position Monitoring

* **[OnPositionProfit](./OnPositionProfit.md)** - periodic profit/loss updates.
* **[OnPositionsAndPendingOrdersTickets](./OnPositionsAndPendingOrdersTickets.md)** - periodic ticket lists (lightweight).

---

## 📚 Step-by-step tutorials

**Note:** Streaming methods are channel-based APIs. Check individual method pages for detailed examples of goroutine patterns and channel handling.

* **[OnSymbolTick](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnSymbolTick_HOW.md)** - Detailed tick streaming examples
* **[OnTrade](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnTrade_HOW.md)** - Trade event monitoring patterns
* **[OnTradeTransaction](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnTradeTransaction_HOW.md)** - Transaction logging examples
* **[OnPositionProfit](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnPositionProfit_HOW.md)** - P/L monitoring patterns
* **[OnPositionsAndPendingOrdersTickets](../HOW_IT_WORK/7.%20Streaming_Methods_HOW/OnPositionsAndPendingOrdersTickets_HOW.md)** - Ticket tracking examples

---

## 🧭 Plain English

* **OnSymbolTick** → **stream live prices** for symbols (bid, ask, volume updates).
* **OnTrade** → **monitor trade events** (position opened/closed/modified).
* **OnTradeTransaction** → **detailed audit log** of all trading operations.
* **OnPositionProfit** → **periodic P/L updates** for open positions.
* **OnPositionsAndPendingOrdersTickets** → **periodic ticket lists** (lightweight monitoring).

> Rule of thumb: need **live prices** → `OnSymbolTick`; need **trade notifications** → `OnTrade`; need **detailed audit** → `OnTradeTransaction`; need **P/L monitoring** → `OnPositionProfit`.

---

## Quick choose

| If you need…                                     | Use                                     | Returns (stream)           | Key inputs                          |
| ------------------------------------------------ | --------------------------------------- | -------------------------- | ----------------------------------- |
| Real-time price ticks                            | `OnSymbolTick`                          | Channel of MqlTick         | Symbol name                         |
| Trade event notifications                        | `OnTrade`                               | Channel of TradeInfo       | *(none)*                            |
| Detailed transaction audit log                   | `OnTradeTransaction`                    | Channel of TradeTransaction| *(none)*                            |
| Real-time profit/loss updates                    | `OnPositionProfit`                      | Channel of PositionProfit  | Symbol filter (optional)            |
| Real-time ticket list changes                    | `OnPositionsAndPendingOrdersTickets`    | Channel of ticket arrays   | *(none)*                            |

---

## ❌ Cross‑refs & gotchas

* **Streaming = Go channels** - Methods return two channels: data channel and error channel.
* **Context cancellation** - Use context.WithCancel() or context.WithTimeout() to stop streams.
* **Goroutine required** - You MUST consume channels in goroutines to avoid blocking.
* **OnSymbolTick** - High frequency, can generate many updates per second.
* **OnTrade** - Triggered on every trade event (open, close, modify, delete).
* **OnTradeTransaction** - Most detailed, includes all transaction types and states.
* **OnPositionProfit** - Real-time updates triggered by price changes.
* **Resource management** - Always cancel context when done to close channels and stop streams.
* **Error handling** - Errors are sent through error channel, nil data indicates closed channel.

---

## 🟢 Minimal snippets

```go
// Stream live ticks for EURUSD
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

dataChan, errChan := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
    Symbol: "EURUSD",
})

go func() {
    for {
        select {
        case tick := <-dataChan:
            if tick == nil {
                return
            }
            fmt.Printf("EURUSD: Bid=%.5f, Ask=%.5f\n", tick.Bid, tick.Ask)

        case err := <-errChan:
            if err != nil {
                log.Printf("Stream error: %v\n", err)
                return
            }
        }
    }
}()

// Keep running...
time.Sleep(30 * time.Second)
```

```go
// Monitor trade events
dataChan, errChan := account.OnTrade(ctx, &pb.OnTradeRequest{})

go func() {
    for {
        select {
        case trade := <-dataChan:
            if trade == nil {
                return
            }
            fmt.Printf("Trade event received\n")

        case err := <-errChan:
            if err != nil {
                return
            }
        }
    }
}()
```

```go
// Monitor position profit/loss
dataChan, errChan := account.OnPositionProfit(ctx, &pb.OnPositionProfitRequest{
    Symbol: "EURUSD", // Or empty for all symbols
})

go func() {
    for {
        select {
        case profitData := <-dataChan:
            if profitData == nil {
                return
            }
            fmt.Printf("#%d %s: $%.2f\n",
                profitData.Ticket, profitData.Symbol, profitData.Profit)

        case err := <-errChan:
            if err != nil {
                return
            }
        }
    }
}()
```

```go
// Monitor position tickets
dataChan, errChan := account.OnPositionsAndPendingOrdersTickets(
    ctx,
    &pb.OnPositionsAndPendingOrdersTicketsRequest{},
)

go func() {
    for {
        select {
        case update := <-dataChan:
            if update == nil {
                return
            }
            fmt.Printf("Open positions: %v\n", update.PositionTickets)
            fmt.Printf("Pending orders: %v\n", update.PendingOrderTickets)

        case err := <-errChan:
            if err != nil {
                return
            }
        }
    }
}()
```

```go
// Detailed transaction log
dataChan, errChan := account.OnTradeTransaction(ctx, &pb.OnTradeTransactionRequest{})

go func() {
    for {
        select {
        case transaction := <-dataChan:
            if transaction == nil {
                return
            }
            fmt.Printf("Transaction event\n")

        case err := <-errChan:
            if err != nil {
                return
            }
        }
    }
}()
```

```go
// Multiple streams concurrently
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Stream 1: Ticks
tickData, tickErr := account.OnSymbolTick(ctx, &pb.OnSymbolTickRequest{
    Symbol: "EURUSD",
})

go func() {
    for {
        select {
        case tick := <-tickData:
            if tick == nil {
                return
            }
            fmt.Printf("[TICK] EURUSD: %.5f\n", tick.Bid)

        case err := <-tickErr:
            if err != nil {
                return
            }
        }
    }
}()

// Stream 2: Trades
tradeData, tradeErr := account.OnTrade(ctx, &pb.OnTradeRequest{})

go func() {
    for {
        select {
        case trade := <-tradeData:
            if trade == nil {
                return
            }
            fmt.Printf("[TRADE] Event received\n")

        case err := <-tradeErr:
            if err != nil {
                return
            }
        }
    }
}()

// Let streams run for 30 seconds
time.Sleep(30 * time.Second)
cancel() // Stop all streams
```

---

## See also

* **Account info:** [AccountSummary](../1.%20Account_information/AccountSummary.md) - get current account state
* **Positions:** [OpenedOrders](../3.%20Position_Orders_Information/OpenedOrders.md) - get current positions snapshot
* **Symbol info:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get current price snapshot
