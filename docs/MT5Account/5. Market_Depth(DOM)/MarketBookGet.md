# ✅ Get Current Market Depth Snapshot

> **Request:** retrieve current Depth of Market (DOM) snapshot with price levels, volumes, and order types at each level.

**API Information:**

* **SDK wrapper:** `MT5Account.MarketBookGet(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.MarketInfo`
* **Proto definition:** `MarketBookGet` (defined in `mt5-term-api-market-info.proto`)

### RPC

* **Service:** `mt5_term_api.MarketInfo`
* **Method:** `MarketBookGet(MarketBookGetRequest) → MarketBookGetReply`
* **Low‑level client (generated):** `MarketInfoClient.MarketBookGet(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

```go
package mt5

type MT5Account struct {
    // ...
}

// MarketBookGet retrieves current market depth snapshot for a symbol.
// Must be subscribed via MarketBookAdd first.
func (a *MT5Account) MarketBookGet(
    ctx context.Context,
    req *pb.MarketBookGetRequest,
) (*pb.MarketBookGetData, error)
```

**Request message:**

```protobuf
MarketBookGetRequest {
  string Symbol = 1;  // Symbol name
}
```

---

## 🔽 Input

| Parameter | Type                         | Description                                   |
| --------- | ---------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`            | Context for deadline/timeout and cancellation |
| `req`     | `*pb.MarketBookGetRequest`   | Request with Symbol name                      |

---

## ⬆️ Output — `MarketBookGetData`

| Field  | Type           | Go Type            | Description                     |
| ------ | -------------- | ------------------ | ------------------------------- |
| `Book` | `BookStruct[]` | `[]*pb.BookStruct` | Array of order book entries     |

**BookStruct contains:**

| Field          | Type       | Description                                  |
| -------------- | ---------- | -------------------------------------------- |
| `Type`         | `BookType` | Order type (see enum below)                  |
| `Price`        | `double`   | Price level                                  |
| `Volume`       | `int64`    | Volume at this level (deprecated)            |
| `VolumeDouble` | `double`   | Volume at this level with extended precision |

---

### 📘 Enum: BookType

| Value | Constant | Description |
|-------|----------|-------------|
| 0 | `BOOK_TYPE_SELL` | Sell order (Offer) |
| 1 | `BOOK_TYPE_BUY` | Buy order (Bid) |
| 2 | `BOOK_TYPE_SELL_MARKET` | Sell order by Market |
| 3 | `BOOK_TYPE_BUY_MARKET` | Buy order by Market |

---

## 💬 Just the essentials

* **What it is.** Retrieves current order book state with bid/ask levels.
* **Why you need it.** Analyze market depth, liquidity, order flow.
* **Requires subscription.** Must call MarketBookAdd before using this method.

---

## 🎯 Purpose

Use it to:

* Retrieve current market depth snapshot
* Analyze order book liquidity
* Monitor bid/ask levels and volumes
* Implement scalping strategies
* Detect large orders and market makers

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [MarketBookGet - How it works](../HOW_IT_WORK/5.%20Market_Depth(DOM)_HOW/MarketBookGet_HOW.md)**

---

## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `5s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Subscription required:** Must call `MarketBookAdd` first before using this method.
* **Real-time data:** Book data updates in real-time after subscription.
* **Broker-dependent:** Number of levels varies by broker (typically 5-10 per side).

---

## 🔗 Usage Examples

### 1) Get market book

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/mt5"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx := context.Background()

    // Subscribe first
    account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{Symbol: "EURUSD"})
    defer account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{Symbol: "EURUSD"})

    // Get market book
    data, err := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{
        Symbol: "EURUSD",
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Market depth for EURUSD: %d levels\n", len(data.Book))

    for i, entry := range data.Book {
        typeStr := "SELL"
        if entry.Type == pb.BookType_BOOK_TYPE_BUY {
            typeStr = "BUY"
        }
        fmt.Printf("  [%d] %s: Price=%.5f, Volume=%.2f\n",
            i, typeStr, entry.Price, entry.VolumeDouble)
    }
}
```

### 2) Analyze bid/ask levels

```go
func AnalyzeOrderBook(account *mt5.MT5Account, symbol string) {
    ctx := context.Background()

    // Subscribe
    account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{Symbol: symbol})
    defer account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{Symbol: symbol})

    // Get book
    data, err := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{Symbol: symbol})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    var buyVolume, sellVolume float64
    buyLevels := 0
    sellLevels := 0

    for _, entry := range data.Book {
        if entry.Type == pb.BookType_BOOK_TYPE_BUY { // BUY
            buyVolume += entry.VolumeDouble
            buyLevels++
        } else { // SELL
            sellVolume += entry.VolumeDouble
            sellLevels++
        }
    }

    fmt.Printf("%s Order Book Analysis:\n", symbol)
    fmt.Printf("  Buy levels: %d, Volume: %.2f\n", buyLevels, buyVolume)
    fmt.Printf("  Sell levels: %d, Volume: %.2f\n", sellLevels, sellVolume)
    fmt.Printf("  Buy/Sell ratio: %.2f\n", buyVolume/sellVolume)
}
```

### 3) Find best bid/ask

```go
func GetBestBidAsk(account *mt5.MT5Account, symbol string) (bestBid, bestAsk float64, err error) {
    ctx := context.Background()

    account.MarketBookAdd(ctx, &pb.MarketBookAddRequest{Symbol: symbol})
    defer account.MarketBookRelease(ctx, &pb.MarketBookReleaseRequest{Symbol: symbol})

    data, err := account.MarketBookGet(ctx, &pb.MarketBookGetRequest{Symbol: symbol})
    if err != nil {
        return 0, 0, err
    }

    for _, entry := range data.Book {
        if entry.Type == pb.BookType_BOOK_TYPE_BUY { // BUY (bid)
            if bestBid == 0 || entry.Price > bestBid {
                bestBid = entry.Price
            }
        } else { // SELL (ask)
            if bestAsk == 0 || entry.Price < bestAsk {
                bestAsk = entry.Price
            }
        }
    }

    fmt.Printf("Best Bid: %.5f, Best Ask: %.5f, Spread: %.5f\n",
        bestBid, bestAsk, bestAsk-bestBid)

    return bestBid, bestAsk, nil
}
```

---

## 📚 See Also

* [MarketBookAdd](./MarketBookAdd.md) - Subscribe to DOM (required first)
* [MarketBookRelease](./MarketBookRelease.md) - Release DOM subscription
* [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - Get Level 1 quotes
