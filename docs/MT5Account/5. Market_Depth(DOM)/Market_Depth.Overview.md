# MT5Account · Market Depth (DOM) - Overview

> Level II quotes, order book data, market depth subscription. Use this page to choose the right API for DOM operations.

## 📁 What lives here

* **[MarketBookAdd](./MarketBookAdd.md)** - subscribe to Market Depth for symbol.
* **[MarketBookGet](./MarketBookGet.md)** - get current order book data.
* **[MarketBookRelease](./MarketBookRelease.md)** - unsubscribe from Market Depth.

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[MarketBookAdd - How it works](../HOW_IT_WORK/5.%20Market_Depth(DOM)_HOW/MarketBookAdd_HOW.md)**
* **[MarketBookGet - How it works](../HOW_IT_WORK/5.%20Market_Depth(DOM)_HOW/MarketBookGet_HOW.md)**
* **[MarketBookRelease - How it works](../HOW_IT_WORK/5.%20Market_Depth(DOM)_HOW/MarketBookRelease_HOW.md)**

---

## 🧭 Plain English

* **MarketBookAdd** → **subscribe** to Level II quotes (order book).
* **MarketBookGet** → **get current DOM** data (bid/ask orders at different price levels).
* **MarketBookRelease** → **unsubscribe** when done.

> Rule of thumb: call `MarketBookAdd` first to subscribe, then `MarketBookGet` to read data, finally `MarketBookRelease` to clean up.

---

## Quick choose

| If you need…                                     | Use                      | Returns                          | Key inputs                          |
| ------------------------------------------------ | ------------------------ | -------------------------------- | ----------------------------------- |
| Subscribe to Market Depth                        | `MarketBookAdd`          | `*pb.MarketBookAddData`          | Symbol name                         |
| Get current order book                           | `MarketBookGet`          | `*pb.MarketBookGetData`          | Symbol name                         |
| Unsubscribe from Market Depth                    | `MarketBookRelease`      | `*pb.MarketBookReleaseData`      | Symbol name                         |

---

## ❌ Cross‑refs & gotchas

* **Subscription required** - Must call MarketBookAdd before MarketBookGet.
* **Not all brokers support DOM** - Check if your broker provides Level II data.
* **MrpcMqlBookInfo** - Each `MqlBookInfos` entry contains Type (BUY/SELL), Price, Volume, VolumeReal.
* **Order book depth** - Broker-dependent, usually 5-10 levels per side.
* **Real-time updates** - DOM updates in real-time after subscription.
* **Clean up** - Always call MarketBookRelease when done to free resources.
* **Symbol must be in Market Watch** - Use SymbolSelect to add symbol first.

---

## 🟢 Minimal snippets

```go
// Subscribe to Market Depth
result, err := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
    Symbol: "EURUSD",
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

if result.OpenedSuccessfully {
    fmt.Println("✅ Subscribed to EURUSD DOM")
} else {
    fmt.Println("❌ DOM subscription failed (broker may not support it)")
}
```

```go
// Get current order book
book, err := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{
    Symbol: "EURUSD",
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Println("EURUSD Order Book:")
fmt.Println("BUY SIDE:")
for _, item := range book.MqlBookInfos {
    if item.Type == pb.BookType_BOOK_TYPE_BUY {
        fmt.Printf("  %.5f: %.2f lots\n", item.Price, item.VolumeReal)
    }
}

fmt.Println("\nSELL SIDE:")
for _, item := range book.MqlBookInfos {
    if item.Type == pb.BookType_BOOK_TYPE_SELL {
        fmt.Printf("  %.5f: %.2f lots\n", item.Price, item.VolumeReal)
    }
}
```

```go
// Unsubscribe from Market Depth
result, _ := account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
    Symbol: "EURUSD",
})

if result.ClosedSuccessfully {
    fmt.Println("✅ Unsubscribed from EURUSD DOM")
}
```

```go
// Complete DOM workflow
symbol := "EURUSD"

// Step 1: Subscribe
if result, _ := account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{
    Symbol: symbol,
}); !result.OpenedSuccessfully {
    log.Fatal("Failed to subscribe to DOM")
}

defer func() {
    // Step 3: Unsubscribe (cleanup)
    account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{
        Symbol: symbol,
    })
}()

// Step 2: Get DOM data
book, _ := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{
    Symbol: symbol,
})

// Analyze order book
totalBuyVolume := 0.0
totalSellVolume := 0.0

for _, item := range book.MqlBookInfos {
    if item.Type == pb.BookType_BOOK_TYPE_BUY {
        totalBuyVolume += item.VolumeReal
    } else {
        totalSellVolume += item.VolumeReal
    }
}

fmt.Printf("Buy volume: %.2f\n", totalBuyVolume)
fmt.Printf("Sell volume: %.2f\n", totalSellVolume)

if totalBuyVolume > totalSellVolume {
    fmt.Println("📈 More buy orders - bullish pressure")
} else {
    fmt.Println("📉 More sell orders - bearish pressure")
}
```

```go
// Monitor DOM changes
symbol := "EURUSD"

// Subscribe
account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{Symbol: symbol})
defer account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{Symbol: symbol})

// Poll every second
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for range ticker.C {
    book, _ := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{
        Symbol: symbol,
    })

    if len(book.MqlBookInfos) > 0 {
        // Get best bid and ask from DOM
        bestBid := 0.0
        bestAsk := 999999.0

        for _, item := range book.MqlBookInfos {
            if item.Type == pb.BookType_BOOK_TYPE_BUY && item.Price > bestBid {
                bestBid = item.Price
            }
            if item.Type == pb.BookType_BOOK_TYPE_SELL && item.Price < bestAsk {
                bestAsk = item.Price
            }
        }

        spread := (bestAsk - bestBid) / 0.00001 // Points
        fmt.Printf("Best Bid: %.5f, Best Ask: %.5f, Spread: %.0f points\n",
            bestBid, bestAsk, spread)
    }
}
```

---

## 📚 See also

* **Symbol management:** [SymbolSelect](../2.%20Symbol_information/SymbolSelect.md) - add symbol to Market Watch
* **Current prices:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get bid/ask quotes
* **Trading:** [OrderSend](../4.%20Trading_Operations/OrderSend.md) - place orders based on DOM analysis
