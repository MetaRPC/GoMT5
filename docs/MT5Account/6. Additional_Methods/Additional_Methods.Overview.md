# MT5Account · Additional Methods - Overview

> Advanced symbol information, trading sessions, margin rates, batch operations. Use this page to choose the right API for advanced queries.

## 📁 What lives here

### Advanced Symbol Information

* **[SymbolInfoMarginRate](./SymbolInfoMarginRate.md)** - margin rates for symbol and order type.
* **[SymbolInfoSessionQuote](./SymbolInfoSessionQuote.md)** - quote session times.
* **[SymbolInfoSessionTrade](./SymbolInfoSessionTrade.md)** - trade session times.
* **[SymbolParamsMany](./SymbolParamsMany.md)** - detailed parameters for one or all symbols (116 fields!).
* **[TickValueWithSize](./TickValueWithSize.md)** - tick value and size information for symbols.

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[SymbolInfoMarginRate - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoMarginRate_HOW.md)**
* **[SymbolInfoSessionQuote - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoSessionQuote_HOW.md)**
* **[SymbolInfoSessionTrade - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolInfoSessionTrade_HOW.md)**
* **[SymbolParamsMany - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolParamsMany_HOW.md)**
* **[TickValueWithSize - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/TickValueWithSize_HOW.md)**

---

## 🧭 Plain English

* **SymbolInfoMarginRate** → get **margin multiplier** for different order types.
* **SymbolInfoSessionQuote** → get **quote session** times (when quotes are available).
* **SymbolInfoSessionTrade** → get **trading session** times (when trading is allowed).
* **SymbolParamsMany** → get **complete symbol specifications** for one symbol or all symbols in one call (most efficient).
* **TickValueWithSize** → get **tick value and contract size** for position value and P&L calculations.

> Rule of thumb: need **session times** → use Session methods; need **complete symbol data** → use `SymbolParamsMany`; need **margin rates** → use `SymbolInfoMarginRate`; need **tick value/P&L data** → use `TickValueWithSize`.

---

## Quick choose

| If you need…                                     | Use                         | Returns                    | Key inputs                          |
| ------------------------------------------------ | --------------------------- | -------------------------- | ----------------------------------- |
| Margin rates for order type                      | `SymbolInfoMarginRate`      | SymbolInfoMarginRateData   | Symbol, OrderType                   |
| Quote session times                              | `SymbolInfoSessionQuote`    | SymbolInfoSessionQuoteData (From, To) | Symbol, DayOfWeek, SessionIndex |
| Trading session times                            | `SymbolInfoSessionTrade`    | SymbolInfoSessionTradeData (From, To) | Symbol, DayOfWeek, SessionIndex |
| Complete specs for one or all symbols            | `SymbolParamsMany`          | SymbolParamsManyData (`SymbolInfos` array of SymbolParameters) | SymbolName (optional), SortType, PageNumber/ItemsPerPage |
| Tick value and size for P&L calculations         | `TickValueWithSize`         | TickValueWithSizeData (`SymbolTickSizeInfos` array of TickSizeSymbol) | SymbolNames array |

---

## ❌ Cross‑refs & gotchas

* **SymbolParamsMany** - THE MOST EFFICIENT way to get complete symbol data (116 fields per symbol).
* **Session times** - `From`/`To` are timestamps whose time of day is the session boundary; ignore the date part.
* **Day of week** - SUNDAY=0, MONDAY=1, ..., SATURDAY=6.
* **Session index** - Most symbols have 1 session, some have multiple (e.g., 24h = 1 session, day+night = 2 sessions).
* **Margin rates** - Different rates for BUY/SELL, LONG/SHORT positions.
* **SymbolParamsMany filters** - Only `SymbolName` (one symbol), sorting and paging; there is no field filter or symbol list. Filter other things on the client side.
* **Batch operations** - SymbolParamsMany without `SymbolName` returns all symbols in one call.

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
    DayOfWeek:    pb.DayOfWeek_MONDAY,
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
// Get complete symbol parameters for several symbols:
// one call for all symbols, then keep the wanted ones (there is no symbol-list filter)
wanted := map[string]bool{"EURUSD": true, "GBPUSD": true, "USDJPY": true}

params, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

for _, symbolParam := range params.SymbolInfos {
    if !wanted[symbolParam.Name] {
        continue
    }
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

var days []pb.DayOfWeek // Monday through Friday
days = append(days, pb.DayOfWeek_MONDAY, pb.DayOfWeek_TUESDAY, pb.DayOfWeek_WEDNESDAY,
    pb.DayOfWeek_THURSDAY, pb.DayOfWeek_FRIDAY)

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
// Batch query all forex majors with SymbolParamsMany (all symbols, filtered on the client side)
majors := map[string]bool{
    "EURUSD": true, "GBPUSD": true, "USDJPY": true, "USDCHF": true,
    "AUDUSD": true, "USDCAD": true, "NZDUSD": true,
}

params, _ := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{})

fmt.Println("Forex Majors Summary:")
fmt.Println("─────────────────────────────────────────────────────")
fmt.Printf("%-8s  %-10s  %-10s  %-8s\n", "Symbol", "Bid", "Ask", "Spread")
fmt.Println("─────────────────────────────────────────────────────")

for _, sp := range params.SymbolInfos {
    if !majors[sp.Name] {
        continue
    }
    fmt.Printf("%-8s  %-10.5f  %-10.5f  %-8d\n",
        sp.Name, sp.Bid, sp.Ask, sp.Spread)
}
```

---

## See also

* **Basic symbol info:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get current prices
* **Symbol properties:** [SymbolInfoDouble](../2.%20Symbol_information/SymbolInfoDouble.md) - get single property
* **Trading:** [OrderCalcMargin](../4.%20Trading_Operations/OrderCalcMargin.md) - uses margin rates internally
