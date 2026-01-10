### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolsTotal()`** method is used to find out **how many trading instruments (symbols)** are available on the MetaTrader server.
> It can work in two modes:
>
> * **Mode = false** → counts *all* symbols available from the broker;
> * **Mode = true** → counts *only active* symbols (Market Watch).
>
> 📦 *Market Watch* in the gateway context is a **list of active symbols** for which the gateway receives real quotes (Bid/Ask).
> Since we work through the API (without the MetaTrader interface), this list can be programmatically enabled or modified using the `SymbolSelect()` method.


---

### 🧩 Code example

```go
fmt.Println("\n4.1. SymbolsTotal() - Count symbols")

// Count all available symbols
allSymbolsReq := &pb.SymbolsTotalRequest{
    Mode: false, // false = all symbols
}
allSymbolsData, err := account.SymbolsTotal(ctx, allSymbolsReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolsTotal(all) failed")
} else {
    fmt.Printf("  Total available symbols:       %d\n", allSymbolsData.Total)
}

// Count symbols in Market Watch only
selectedSymbolsReq := &pb.SymbolsTotalRequest{
    Mode: true, // true = Market Watch only
}
selectedSymbolsData, err := account.SymbolsTotal(ctx, selectedSymbolsReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolsTotal(selected) failed")
} else {
    fmt.Printf("  Symbols in Market Watch:       %d\n", selectedSymbolsData.Total)
}
```

---

### 🟢 Detailed Code Explanation

#### 1️ First call (`Mode: false`)

Requests the total number of symbols available on the broker's server.
This is *all* pairs, stocks, futures, and instruments that the broker supports — even those without current quotes.

🧠 Example output:

```
Total available symbols: 245
```

---

#### 2️ Second call (`Mode: true`)

Requests the number of **active symbols** currently in *Market Watch* state.
These are instruments for which MetaTrader actually streams quotes (Bid/Ask).

🧠 Example output:

```
Symbols in Market Watch: 32
```

---

### 🧱 What is Market Watch (in gateway context)

> Market Watch is not a terminal window, but a **logical symbol state on the server**.
> Active symbol = one you've subscribed to for quotes.

When we work without the MetaTrader interface (through the gateway):

* we **don't have a visual Market Watch**;
* but we **can enable it programmatically**.

For this, the **`SymbolSelect()`** method is used:

```go
selectReq := &pb.SymbolSelectRequest{
    Symbol: "EURUSD",
    Select: true, // true = add symbol to Market Watch
}
_, err := account.SymbolSelect(ctx, selectReq)
if err != nil {
    log.Fatalf("SymbolSelect failed: %v", err)
}
fmt.Println("EURUSD added to Market Watch")
```

After this, `EURUSD` is considered active — and will appear in the `SymbolsTotal(Mode: true)` count.

---

### 🛰️ Quote subscription (Bid/Ask stream)

After adding a symbol to Market Watch, you can subscribe to the quote stream:

```go
ticks, err := account.TicksSubscribe(ctx, &pb.TicksSubscribeRequest{
    Symbol: "EURUSD",
})
if err != nil {
    log.Fatalf("Subscribe failed: %v", err)
}

for {
    tick, err := ticks.Recv()
    if err != nil {
        break
    }
    fmt.Printf("[%s] Bid: %.5f | Ask: %.5f\n", tick.Symbol, tick.Bid, tick.Ask)
}
```

Now the gateway receives real quotes — analogous to how the MetaTrader terminal displays them in its Market Watch.

---

### ⚙️ Practical application

| Goal                                              | How to use                                 |
| ------------------------------------------------- | ------------------------------------------ |
| Find out total number of broker's instruments     | `SymbolsTotal(Mode: false)`                |
| Check how many instruments are active             | `SymbolsTotal(Mode: true)`                 |
| Add symbol to Market Watch                        | `SymbolSelect(Symbol, true)`               |
| Remove symbol from Market Watch                   | `SymbolSelect(Symbol, false)`              |
| Receive quotes for a symbol                       | `TicksSubscribe()` or `SymbolSubscribe()` |

---

### 💬 In simple terms


> SymbolsTotal() itself doesn't activate anything — it only counts.
> Market Watch when working through the API = list of active symbols that we manage via SymbolSelect().
>
> After enabling a symbol in Market Watch, you can stream quotes, trade it, and get real-time updates.

