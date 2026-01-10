# MT5Account · Symbol Information - Overview

> Quotes, symbol properties, trading rules, tick values, Market Watch management. Use this page to choose the right API for symbol data.

## 📁 What lives here

### Current Quotes

* **[SymbolInfoTick](./SymbolInfoTick.md)** - current quote for symbol (bid, ask, last, volume, time).

### Symbol Inventory & Management

* **[SymbolsTotal](./SymbolsTotal.md)** - count of available symbols.
* **[SymbolName](./SymbolName.md)** - get symbol name by index.
* **[SymbolSelect](./SymbolSelect.md)** - enable/disable symbol in Market Watch.
* **[SymbolExist](./SymbolExist.md)** - check if symbol exists.
* **[SymbolIsSynchronized](./SymbolIsSynchronized.md)** - check symbol data sync status.

### Symbol Properties

* **[SymbolInfoDouble](./SymbolInfoDouble.md)** - single double property (bid, ask, point, volume min/max, etc.).
* **[SymbolInfoInteger](./SymbolInfoInteger.md)** - single integer property (digits, spread, stops level, etc.).
* **[SymbolInfoString](./SymbolInfoString.md)** - single string property (description, currency, path).

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[SymbolInfoTick - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoTick_HOW.md)**
* **[SymbolsTotal - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolsTotal_HOW.md)**
* **[SymbolName - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolName_HOW.md)**
* **[SymbolSelect - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolSelect_HOW.md)**
* **[SymbolExist - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolExist_HOW.md)**
* **[SymbolIsSynchronized - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolIsSynchronized_HOW.md)**
* **[SymbolInfoDouble - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoDouble_HOW.md)**
* **[SymbolInfoInteger - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoInteger_HOW.md)**
* **[SymbolInfoString - How it works](../HOW_IT_WORK/2.%20Symbol_information_HOW/SymbolInfoString_HOW.md)**

---

## 🧭 Plain English

* **SymbolInfoTick** → get **current prices** (bid, ask, last).
* **SymbolsTotal/SymbolName** → **iterate through all symbols**.
* **SymbolSelect** → **add/remove** symbol from Market Watch.
* **SymbolExist** → **check if symbol exists** on broker.
* **SymbolInfo*** → get **specific symbol property** (double/integer/string).

> Rule of thumb: need **current price** → `SymbolInfoTick`; need **symbol list** → `SymbolsTotal + SymbolName`; need **specific property** → `SymbolInfo*`.

---

## Quick choose

| If you need…                                     | Use                      | Returns                    | Key inputs                          |
| ------------------------------------------------ | ------------------------ | -------------------------- | ----------------------------------- |
| Current quote (bid, ask, last)                   | `SymbolInfoTick`         | MrpcMqlTick                | Symbol                              |
| Count available symbols                          | `SymbolsTotal`           | int32                      | Mode (bool)                         |
| Get symbol name by index                         | `SymbolName`             | string                     | Index, Selected                     |
| Add/remove symbol from Market Watch              | `SymbolSelect`           | bool                       | Symbol, Select (true/false)         |
| Check if symbol exists                           | `SymbolExist`            | bool                       | Name                                |
| Check if symbol data is synced                   | `SymbolIsSynchronized`   | bool                       | Symbol                              |
| One numeric property (bid, point, volume, etc.)  | `SymbolInfoDouble`       | float64                    | Symbol, Type (enum)                 |
| One integer property (digits, spread, etc.)      | `SymbolInfoInteger`      | int64                      | Symbol, Type (enum)                 |
| One text property (description, currency)        | `SymbolInfoString`       | string                     | Symbol, Type (enum)                 |

---

## ❌ Cross‑refs & gotchas

* **SymbolInfoTick** - Use this for current prices; updates in real-time.
* **Bid vs Ask** - BUY orders use Ask price, SELL orders use Bid price.
* **Market Watch** - Symbol must be in Market Watch to get live updates (use `SymbolSelect`).
* **Synchronization** - Check `SymbolIsSynchronized` before trading to ensure data is current.
* **Digits** - Number of decimal places (e.g., EURUSD = 5, USDJPY = 3).
* **Point** - Minimal price change (e.g., 0.00001 for 5-digit EURUSD).
* **Volume limits** - Check VOLUME_MIN, VOLUME_MAX, VOLUME_STEP before placing orders.
* **Stop level** - Minimum distance for SL/TP from current price (in points).

---

## 🟢 Minimal snippets

```go
// Get current quote
tick, err := account.SymbolInfoTick(ctx, &pb.SymbolInfoTickRequest{
    Symbol: "EURUSD",
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("EURUSD: Bid=%.5f, Ask=%.5f, Spread=%.5f\n",
    tick.Bid, tick.Ask, tick.Ask-tick.Bid)
```

```go
// List all available symbols
totalData, _ := account.SymbolsTotal(ctx, &pb.SymbolsTotalRequest{
    Selected: false,
})

fmt.Printf("Available symbols: %d\n", totalData.Total)

for i := int32(0); i < totalData.Total; i++ {
    nameData, _ := account.SymbolName(ctx, &pb.SymbolNameRequest{
        Index:    i,
        Selected: false,
    })
    fmt.Printf("%d. %s\n", i+1, nameData.Name)
}
```

```go
// Add symbol to Market Watch
selectData, _ := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
    Symbol: "GBPUSD",
    Select: true,
})

if selectData.Success {
    fmt.Println("✅ GBPUSD added to Market Watch")
} else {
    fmt.Println("❌ Failed to add GBPUSD")
}
```

```go
// Check if symbol exists
existData, _ := account.SymbolExist(ctx, &pb.SymbolExistRequest{
    Name: "EURUSD",
})

if existData.Exists {
    fmt.Println("✅ Symbol exists")
} else {
    fmt.Println("❌ Symbol not found")
}
```

```go
// Get symbol properties
digitsData, _ := account.SymbolInfoInteger(ctx, &pb.SymbolInfoIntegerRequest{
    Symbol: "EURUSD",
    Type:   pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS,
})

pointData, _ := account.SymbolInfoDouble(ctx, &pb.SymbolInfoDoubleRequest{
    Symbol: "EURUSD",
    Type:   pb.SymbolInfoDoubleProperty_SYMBOL_POINT,
})

fmt.Printf("EURUSD: Digits=%d, Point=%.5f\n", digitsData.Value, pointData.Value)
```

---

## See also

* **Advanced symbol info:** [SymbolParamsMany](../6.%20Additional_Methods/SymbolParamsMany.md) - get 112+ fields for multiple symbols
* **Trading:** [OrderSend](../4.%20Trading_Operations/OrderSend.md) - place orders
* **Real-time updates:** [OnSymbolTick](../7.%20Streaming_Methods/OnSymbolTick.md) - subscribe to tick stream
