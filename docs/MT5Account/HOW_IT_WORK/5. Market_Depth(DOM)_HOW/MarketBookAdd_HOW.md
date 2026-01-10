### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`MarketBookAdd()`** method subscribes to **Depth of Market (DOM)** — market depth for a specified symbol.
> DOM shows real order levels (Bid/Ask levels, volumes, and prices), if the broker provides them.


---

## 🧩 Code example

```go
// Use short timeout
bookAddCtx, bookAddCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer bookAddCancel()

marketBookAddReq := &pb.MarketBookAddRequest{
    Symbol: cfg.TestSymbol, // symbol from config.json (test_symbol parameter)
}
marketBookAddData, err := account.MarketBookAdd(bookAddCtx, marketBookAddReq)
if err != nil {
    fmt.Printf("  ❌ Subscription failed: %v\n", err)
    fmt.Println("     → Broker does not support DOM for this symbol")
} else {
    if marketBookAddData.Success {
        fmt.Printf("  ✓ Subscription accepted for '%s'\n", cfg.TestSymbol)
        fmt.Println("     → This doesn't mean data will be available")
    } else {
        fmt.Println("  ❌ Subscription rejected by broker")
    }
}
```

---

### 🟢 Detailed Code Explanation

```go
bookAddCtx, bookAddCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer bookAddCancel()
```

A context with a 5-second timeout is created to prevent the request from hanging if the server doesn't respond.

---

```go
marketBookAddReq := &pb.MarketBookAddRequest{
    Symbol: cfg.TestSymbol,
}
```

A DOM subscription request is formed for the test symbol (e.g., `EURUSD`).
The `cfg.TestSymbol` field is taken from the `config.json` configuration file.

---

```go
marketBookAddData, err := account.MarketBookAdd(bookAddCtx, marketBookAddReq)
```

A gRPC request is sent to MetaTrader.
The server returns a flag indicating whether the DOM subscription for the specified instrument was successful.

---

```go
if err != nil {
    fmt.Printf("  ❌ Subscription failed: %v\n", err)
    fmt.Println("     → Broker does not support DOM for this symbol")
} else {
    if marketBookAddData.Success {
        fmt.Printf("  ✓ Subscription accepted for '%s'\n", cfg.TestSymbol)
        fmt.Println("     → This doesn't mean data will be available")
    } else {
        fmt.Println("  ❌ Subscription rejected by broker")
    }
}
```

Handling possible scenarios:

* **Error (`err`)** — broker or server does not support DOM.
* **Subscription accepted (`Success = true`)** — broker allowed subscription, but doesn't guarantee data availability.
* **Subscription rejected** — broker explicitly refused (common on demo or forex accounts).

---

## 📘 What is Depth of Market (DOM)

DOM (market depth) — a table of all available Bid and Ask levels, reflecting current limit orders from participants.
It shows not only current buy/sell prices, but also volumes at each level.

DOM is usually available for:

* exchange instruments (futures, stocks, CFDs);
* instruments with real order book from liquidity providers.

---

## ⚠️ Why DOM is Often Unavailable on Forex and Demo

1. **Forex is an OTC market**, there is no centralized order book. Therefore, brokers usually don't publish order levels.
2. **On demo accounts**, brokers often don't form real DOM to reduce server load.
3. **Not all brokers support Level II data** — some provide only Bid/Ask.

---

## 💡 What to Do if DOM is Unavailable

* Try **another instrument** (e.g., stock CFDs or futures).
* Check on a **real account**, not demo.
* Ask your broker if the server supports **Depth of Market**.

---

---

## 💬 In Simple Terms

> `MarketBookAdd()` attempts to connect to the market depth data stream.
> If the broker doesn't support DOM for the selected symbol, the server will return an error or empty subscription.
> This is normal behavior for forex and demo servers.
