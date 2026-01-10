# MT5Account · Additional Methods - Overview

> Advanced symbol information, trading sessions, margin rates, batch operations. Use this page to choose the right API for advanced queries.

## 📁 What lives here

### Advanced Symbol Information

* **[SymbolInfoMarginRate](./SymbolInfoMarginRate.md)** - margin rates for symbol and order type.
* **[SymbolInfoSessionQuote](./SymbolInfoSessionQuote.md)** - quote session times.
* **[SymbolInfoSessionTrade](./SymbolInfoSessionTrade.md)** - trade session times.
* **[SymbolParamsMany](./SymbolParamsMany.md)** - detailed parameters for multiple symbols (112 fields!).

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[SymbolInfoMarginRate - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoMarginRate_HOW.md)**
* **[SymbolInfoSessionQuote - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoSessionQuote_HOW.md)**
* **[SymbolInfoSessionTrade - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoSessionTrade_HOW.md)**
* **[SymbolParamsMany - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolParamsMany_HOW.md)**

---

## 🧭 Plain English

* **SymbolInfoMarginRate** → get **margin multiplier** for different order types.
* **SymbolInfoSessionQuote** → get **quote session** times (when quotes are available).
* **SymbolInfoSessionTrade** → get **trading session** times (when trading is allowed).
* **SymbolParamsMany** → get **complete symbol specifications** for multiple symbols in one call (most efficient).

> Rule of thumb: need **session times** → use Session methods; need **complete symbol data** → use `SymbolParamsMany`; need **margin rates** → use `SymbolInfoMarginRate`.

---

## Quick choose

| If you need…                                     | Use                         | Returns                    | Key inputs                          |
| ------------------------------------------------ | --------------------------- | -------------------------- | ----------------------------------- |
| Margin rates for order type                      | `SymbolInfoMarginRate`      | MarginRates                | Symbol, order type                  |
| Quote session times                              | `SymbolInfoSessionQuote`    | SessionInfo                | Symbol, day of week, session index  |
| Trading session times                            | `SymbolInfoSessionTrade`    | SessionInfo                | Symbol, day of week, session index  |
| Complete specs for multiple symbols              | `SymbolParamsMany`          | Array of SymbolParams      | Symbol list, optional filters       |

---

## ❌ Cross‑refs & gotchas

* **SymbolParamsMany** - THE MOST EFFICIENT way to get complete symbol data (112 fields per symbol).
* **Session times** - In seconds from day start (0-86400).
* **Day of week** - SUNDAY=0, MONDAY=1, ..., SATURDAY=6.
* **Session index** - Most symbols have 1 session, some have multiple (e.g., 24h = 1 session, day+night = 2 sessions).
* **Margin rates** - Different rates for BUY/SELL, LONG/SHORT positions.
* **SymbolParamsMany filters** - Can filter by Bid, Ask, Point, Volume to reduce data size.
* **Batch operations** - SymbolParamsMany can query 100+ symbols in one call.

---

## 🟢 Minimal snippets

```go
// Get margin rates for symbol
rates, err := account.SymbolInfoMarginRate(ctx, &pb.SymbolInfoMarginRateRequest{
    Symbol:    "EURUSD",
    OrderType: pb.ENUM_ORDER_TYPE_ORDER_TYPE_BUY,
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("EURUSD BUY margin rates:\n")
fmt.Printf("  Initial: %.2f\n", rates.InitialMarginRate)
fmt.Printf("  Maintenance: %.2f\n", rates.MaintenanceMarginRate)
```

```go
// Get trading session times for Monday
sessionInfo, _ := account.SymbolInfoSessionTrade(ctx, &pb.SymbolInfoSessionTradeRequest{
    Symbol:       "EURUSD",
    DayOfWeek:    1, // Monday
    SessionIndex: 0,
})

fromTime := sessionInfo.From.AsTime()
toTime := sessionInfo.To.AsTime()
fromSeconds := fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()
toSeconds := toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()

fmt.Printf("EURUSD Monday trading session:\n")
fmt.Printf("  From: %d seconds (%.2f hours)\n",
    fromSeconds, float64(fromSeconds)/3600)
fmt.Printf("  To: %d seconds (%.2f hours)\n",
    toSeconds, float64(toSeconds)/3600)
```

```go
// Get complete symbol parameters for multiple symbols
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

params, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
    Symbols: symbols,
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

for _, symbolParam := range params.SymbolParams {
    fmt.Printf("\n%s:\n", symbolParam.Name)
    fmt.Printf("  Bid: %.5f, Ask: %.5f\n", symbolParam.Bid, symbolParam.Ask)
    fmt.Printf("  Digits: %d, Point: %.5f\n", symbolParam.Digits, symbolParam.Point)
    fmt.Printf("  Volume: %.2f - %.2f (step: %.2f)\n",
        symbolParam.VolumeMin, symbolParam.VolumeMax, symbolParam.VolumeStep)
    fmt.Printf("  Contract size: %.0f\n", symbolParam.TradeContractSize)
    fmt.Printf("  Spread: %d points\n", symbolParam.Spread)
}
```

```go
// Get all trading sessions for a week
symbol := "EURUSD"

days := []int32{1, 2, 3, 4, 5} // Monday through Friday

dayNames := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

fmt.Printf("%s Trading Schedule:\n", symbol)
fmt.Println("─────────────────────────────────────")

for i, day := range days {
    session, err := account.SymbolInfoSessionTrade(ctx, &pb.SymbolInfoSessionTradeRequest{
        Symbol:       symbol,
        DayOfWeek:    day,
        SessionIndex: 0,
    })

    if err != nil {
        continue
    }

    fromTime := session.From.AsTime()
    toTime := session.To.AsTime()
    fromHours := float64(fromTime.Hour()*3600 + fromTime.Minute()*60 + fromTime.Second()) / 3600
    toHours := float64(toTime.Hour()*3600 + toTime.Minute()*60 + toTime.Second()) / 3600

    fmt.Printf("%s: %.2f:00 - %.2f:00\n",
        dayNames[i], fromHours, toHours)
}
```

```go
// Batch query all forex majors with SymbolParamsMany
majors := []string{
    "EURUSD", "GBPUSD", "USDJPY", "USDCHF",
    "AUDUSD", "USDCAD", "NZDUSD",
}

params, _ := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
    Symbols: majors,
})

fmt.Println("Forex Majors Summary:")
fmt.Println("─────────────────────────────────────────────────────")
fmt.Printf("%-8s  %-10s  %-10s  %-8s\n", "Symbol", "Bid", "Ask", "Spread")
fmt.Println("─────────────────────────────────────────────────────")

for _, sp := range params.SymbolParams {
    fmt.Printf("%-8s  %-10.5f  %-10.5f  %-8d\n",
        sp.Name, sp.Bid, sp.Ask, sp.Spread)
}
```

---

## See also

* **Basic symbol info:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get current prices
* **Symbol properties:** [SymbolInfoDouble](../2.%20Symbol_information/SymbolInfoDouble.md) - get single property
* **Trading:** [OrderCalcMargin](../4.%20Trading_Operations/OrderCalcMargin.md) - uses margin rates internally
