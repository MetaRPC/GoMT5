# ✅ Get Parameters for Multiple Symbols

> **Request:** get complete trading parameters for multiple symbols in one request. Returns Bid, Ask, Digits, Spread, volumes, contract size and other parameters.

**API Information:**

* **Low-level API:** `MT5Account.SymbolParamsMany(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `SymbolParamsMany` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `SymbolParamsMany(SymbolParamsManyRequest) → SymbolParamsManyReply`
* **Low‑level client (generated):** `AccountHelperClient.SymbolParamsMany(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieves comprehensive trading parameters for multiple symbols at once.
* **Why you need it.** Efficient bulk data retrieval, avoid multiple round-trips.
* **Complete data.** Returns all essential symbol properties in one call.

---

## 🎯 Purpose

Use it to:

* Get complete trading specifications for multiple symbols efficiently
* Build trading dashboards and watchlists
* Validate trading parameters in bulk
* Compare spreads and conditions across symbols
* Reduce network round-trips with batch queries

---

```go
package mt5

type MT5Account struct {
    // ...
}

// SymbolParamsMany retrieves detailed parameters for multiple symbols in one call.
// This is the recommended method for getting comprehensive symbol data.
func (a *MT5Account) SymbolParamsMany(
    ctx context.Context,
    req *pb.SymbolParamsManyRequest,
) (*pb.SymbolParamsManyData, error)
```

**Request message:**

```protobuf
SymbolParamsManyRequest {
  repeated string Symbols = 1;  // Array of symbol names
}
```

---

## 🔽 Input

| Parameter | Type                            | Description                                   |
| --------- | ------------------------------- | --------------------------------------------- |
| `ctx`     | `context.Context`               | Context for deadline/timeout and cancellation |
| `req`     | `*pb.SymbolParamsManyRequest`   | Request with array of Symbols                 |

**Request fields:**

| Field     | Type       | Description                        |
| --------- | ---------- | ---------------------------------- |
| `Symbols` | `string[]` | Array of symbol names              |

---

## ⬆️ Output — `SymbolParamsManyData`

| Field          | Type                    | Go Type                 | Description                              |
| -------------- | ----------------------- | ----------------------- | ---------------------------------------- |
| `SymbolInfos`  | `SymbolParameters[]`    | `[]*SymbolParameters`   | Array of detailed symbol parameters      |
| `SymbolsTotal` | `int32`                 | `int32`                 | Total number of symbols returned         |
| `PageNumber`   | `int32` (optional)      | `*int32`                | Page number (if pagination used)         |
| `ItemsPerPage` | `int32` (optional)      | `*int32`                | Items per page (if pagination used)      |

---

### 📋 `SymbolParameters` structure (112 fields)

Each symbol contains comprehensive trading information organized in the following categories:

#### 🔹 **Basic Price Data** (14 fields)
| Field            | Type      | Description                          |
| ---------------- | --------- | ------------------------------------ |
| `Name`           | `string`  | Symbol name (e.g., "EURUSD")         |
| `Bid`            | `float64` | Current bid price                    |
| `BidHigh`        | `float64` | Highest bid of the day               |
| `BidLow`         | `float64` | Lowest bid of the day                |
| `Ask`            | `float64` | Current ask price                    |
| `AskHigh`        | `float64` | Highest ask of the day               |
| `AskLow`         | `float64` | Lowest ask of the day                |
| `Last`           | `float64` | Last deal price                      |
| `LastHigh`       | `float64` | Highest last price of the day        |
| `LastLow`        | `float64` | Lowest last price of the day         |
| `VolumeReal`     | `float64` | Real volume of the day               |
| `VolumeHighReal` | `float64` | Maximum real volume of the day       |
| `VolumeLowReal`  | `float64` | Minimum real volume of the day       |
| `OptionStrike`   | `float64` | Option strike price                  |

#### 🔹 **Trading Parameters** (13 fields)
| Field                  | Type      | Description                          |
| ---------------------- | --------- | ------------------------------------ |
| `Point`                | `float64` | Point size (minimal price change)    |
| `TradeTickValue`       | `float64` | Tick value for profit calculation    |
| `TradeTickValueProfit` | `float64` | Tick value for profitable positions  |
| `TradeTickValueLoss`   | `float64` | Tick value for losing positions      |
| `TradeTickSize`        | `float64` | Minimal price change                 |
| `TradeContractSize`    | `float64` | Trade contract size                  |
| `TradeAccruedInterest` | `float64` | Accrued interest                     |
| `TradeFaceValue`       | `float64` | Face value of the instrument         |
| `TradeLiquidityRate`   | `float64` | Liquidity rate                       |
| `VolumeMin`            | `float64` | Minimum volume for a deal            |
| `VolumeMax`            | `float64` | Maximum volume for a deal            |
| `VolumeStep`           | `float64` | Minimal volume change step           |
| `VolumeLimit`          | `float64` | Maximum allowed aggregate volume     |

#### 🔹 **Swap Parameters** (9 fields)
| Field           | Type      | Description                          |
| --------------- | --------- | ------------------------------------ |
| `SwapLong`      | `float64` | Long position swap value             |
| `SwapShort`     | `float64` | Short position swap value            |
| `SwapSunday`    | `float64` | Sunday swap                          |
| `SwapMonday`    | `float64` | Monday swap                          |
| `SwapTuesday`   | `float64` | Tuesday swap                         |
| `SwapWednesday` | `float64` | Wednesday swap                       |
| `SwapThursday`  | `float64` | Thursday swap                        |
| `SwapFriday`    | `float64` | Friday swap                          |
| `SwapSaturday`  | `float64` | Saturday swap                        |

#### 🔹 **Margin Parameters** (3 fields)
| Field                | Type      | Description                          |
| -------------------- | --------- | ------------------------------------ |
| `MarginInitial`      | `float64` | Initial margin requirement           |
| `MarginMaintenance`  | `float64` | Maintenance margin requirement       |
| `MarginHedged`       | `float64` | Hedged margin                        |

#### 🔹 **Session Statistics** (11 fields)
| Field                      | Type      | Description                          |
| -------------------------- | --------- | ------------------------------------ |
| `SessionVolume`            | `float64` | Summary volume of the current session |
| `SessionTurnover`          | `float64` | Summary turnover of the current session |
| `SessionInterest`          | `float64` | Summary open interest                |
| `SessionBuyOrdersVolume`   | `float64` | Current volume of buy orders         |
| `SessionSellOrdersVolume`  | `float64` | Current volume of sell orders        |
| `SessionOpen`              | `float64` | Session open price                   |
| `SessionClose`             | `float64` | Session close price                  |
| `SessionAw`                | `float64` | Average weighted price               |
| `SessionPriceSettlement`   | `float64` | Settlement price                     |
| `SessionPriceLimitMin`     | `float64` | Minimum price limit                  |
| `SessionPriceLimitMax`     | `float64` | Maximum price limit                  |

#### 🔹 **Options Greeks & Price Analytics** (10 fields)
| Field               | Type      | Description                          |
| ------------------- | --------- | ------------------------------------ |
| `PriceChange`       | `float64` | Change of price in %                 |
| `PriceVolatility`   | `float64` | Price volatility in %                |
| `PriceTheoretical`  | `float64` | Theoretical option price             |
| `PriceDelta`        | `float64` | Option/warrant delta                 |
| `PriceTheta`        | `float64` | Option/warrant theta                 |
| `PriceGamma`        | `float64` | Option/warrant gamma                 |
| `PriceVega`         | `float64` | Option/warrant vega                  |
| `PriceRho`          | `float64` | Option/warrant rho                   |
| `PriceOmega`        | `float64` | Option/warrant omega                 |
| `PriceSensitivity`  | `float64` | Option/warrant sensitivity           |

#### 🔹 **Symbol Properties** (9 fields)
| Field                | Type      | Description                          |
| -------------------- | --------- | ------------------------------------ |
| `Sector`             | `enum`    | Economic sector (enum value)         |
| `Industry`           | `enum`    | Industry sector (enum value)         |
| `Custom`             | `bool`    | Custom symbol flag                   |
| `BackgroundColor`    | `string`  | Background color for symbol          |
| `ChartMode`          | `enum`    | Chart mode (Bid/Last price)          |
| `Exist`              | `bool`    | Symbol exists                        |
| `Select`             | `bool`    | Symbol selected in Market Watch      |
| `SubscriptionDelay`  | `int32`   | Subscription delay in seconds        |
| `Visible`            | `bool`    | Symbol visible in Market Watch       |

#### 🔹 **Session Counters** (8 fields)
| Field                | Type      | Description                          |
| -------------------- | --------- | ------------------------------------ |
| `SessionDeals`       | `int64`   | Number of deals in current session   |
| `SessionBuyOrders`   | `int64`   | Number of buy orders                 |
| `SessionSellOrders`  | `int64`   | Number of sell orders                |
| `Volume`             | `int64`   | Volume of the last deal              |
| `VolumeHigh`         | `int64`   | Maximum volume of the day            |
| `VolumeLow`          | `int64`   | Minimum volume of the day            |
| `Time`               | `timestamp` | Time of the last quote             |
| `TimeMsc`            | `int64`   | Time of last quote in milliseconds   |

#### 🔹 **Trading Settings** (20 fields)
| Field                  | Type      | Description                          |
| ---------------------- | --------- | ------------------------------------ |
| `Digits`               | `int32`   | Number of decimal digits             |
| `SpreadFloat`          | `bool`    | Floating spread flag                 |
| `Spread`               | `int32`   | Current spread in points             |
| `TicksBookDepth`       | `int32`   | Maximal number of requests in order book |
| `TradeCalcMode`        | `enum`    | Contract price calculation mode      |
| `TradeMode`            | `enum`    | Order execution type                 |
| `StartTime`            | `timestamp` | Date of symbol trading start       |
| `ExpirationTime`       | `timestamp` | Date of symbol expiration          |
| `TradeStopsLevel`      | `int32`   | Minimal distance from price for stops |
| `TradeFreezeLevel`     | `int32`   | Distance to freeze trade operations  |
| `TradeExeMode`         | `enum`    | Deal execution mode                  |
| `SwapMode`             | `enum`    | Swap calculation method              |
| `SwapRollover_3Days`   | `enum`    | Day of week for triple rollover swap |
| `MarginHedgedUseLeg`   | `bool`    | Calculate margin for hedged positions using larger leg |
| `ExpirationMode`       | `int32`   | Flags of allowed order expiration modes |
| `FillingMode`          | `enum[]`  | Flags of allowed order filling modes |
| `OrderMode`            | `enum`    | Flags of allowed order types         |
| `OrderGtcMode`         | `enum`    | Expiration mode for pending orders   |
| `OptionMode`           | `enum`    | Option type                          |
| `OptionRight`          | `enum`    | Option right (Call/Put)              |

#### 🔹 **String Metadata** (15 fields)
| Field            | Type      | Description                          |
| ---------------- | --------- | ------------------------------------ |
| `Basis`          | `string`  | Underlying asset for derivative      |
| `Category`       | `string`  | Symbol category                      |
| `Country`        | `string`  | Country of origin                    |
| `SectorName`     | `string`  | Sector name                          |
| `IndustryName`   | `string`  | Industry name                        |
| `CurrencyBase`   | `string`  | Base currency                        |
| `CurrencyProfit` | `string`  | Profit currency                      |
| `CurrencyMargin` | `string`  | Margin currency                      |
| `Bank`           | `string`  | Feeder for financial instruments     |
| `SymDescription` | `string`  | Symbol description                   |
| `Exchange`       | `string`  | Exchange name                        |
| `Formula`        | `string`  | Formula for custom symbol pricing    |
| `Isin`           | `string`  | ISIN identifier                      |
| `Page`           | `string`  | Web page URL with symbol information |
| `Path`           | `string`  | Path in symbol tree                  |

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [SymbolParamsMany - How it works](../HOW_IT_WORK/6.%20Additional_Methods_HOW/SymbolParamsMany_HOW.md)**

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Most efficient:** This is the recommended method for retrieving comprehensive symbol data.
* **Bulk operations:** Can query 100+ symbols in single request.
* **112+ fields:** Each SymbolParams contains over 112 different properties.
* **Performance:** Much faster than calling individual SymbolInfo* methods repeatedly.

---

## 🔗 Usage Examples

### 1) Get params for multiple symbols

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
)

func main() {
    account, _ := mt5.NewMT5Account(12345, "password", "mt5.mrpc.pro:443", uuid.New())
    defer account.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: []string{"EURUSD", "GBPUSD", "USDJPY"},
    })
    if err != nil {
        panic(err)
    }

    for _, params := range data.SymbolInfos {
        fmt.Printf("%s:\n", params.Name)
        fmt.Printf("  Bid: %.5f, Ask: %.5f\n", params.Bid, params.Ask)
        fmt.Printf("  Digits: %d, Spread: %d\n", params.Digits, params.Spread)
        fmt.Printf("  Volume: %.2f - %.2f (step: %.2f)\n",
            params.VolumeMin, params.VolumeMax, params.VolumeStep)
        fmt.Println()
    }
}
```

### 2) Build trading dashboard

```go
func BuildTradingDashboard(account *mt5.MT5Account, watchlist []string) {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: watchlist,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println("╔════════════════════════════════════════════════════════════╗")
    fmt.Println("║                   TRADING DASHBOARD                        ║")
    fmt.Println("╠═══════════╦══════════╦══════════╦════════╦════════════════╣")
    fmt.Println("║  Symbol   ║   Bid    ║   Ask    ║ Spread ║ Contract Size  ║")
    fmt.Println("╠═══════════╬══════════╬══════════╬════════╬════════════════╣")

    for _, params := range data.SymbolInfos {
        fmt.Printf("║ %-9s ║ %8.5f ║ %8.5f ║ %6d ║ %14.0f ║\n",
            params.Name, params.Bid, params.Ask, params.Spread, params.TradeContractSize)
    }

    fmt.Println("╚═══════════╩══════════╩══════════╩════════╩════════════════╝")
}

// Usage:
// BuildTradingDashboard(account, []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD"})
```

### 3) Validate trading parameters

```go
func ValidateTradingParams(account *mt5.MT5Account, symbol string, volume float64) error {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: []string{symbol},
    })
    if err != nil {
        return err
    }

    if len(data.SymbolInfos) == 0 {
        return fmt.Errorf("symbol %s not found", symbol)
    }

    params := data.SymbolInfos[0]

    // Validate volume
    if volume < params.VolumeMin {
        return fmt.Errorf("volume %.2f below minimum %.2f", volume, params.VolumeMin)
    }
    if volume > params.VolumeMax {
        return fmt.Errorf("volume %.2f above maximum %.2f", volume, params.VolumeMax)
    }

    // Check volume step
    remainder := math.Mod(volume, params.VolumeStep)
    if remainder > 0.0001 {
        return fmt.Errorf("volume %.2f not aligned to step %.2f", volume, params.VolumeStep)
    }

    fmt.Printf("Volume %.2f is valid for %s\n", volume, symbol)
    return nil
}
```

### 4) Compare spreads across symbols

```go
func CompareSpreads(account *mt5.MT5Account, symbols []string) {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: symbols,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println("Spread comparison:")
    for _, params := range data.SymbolInfos {
        spreadPips := float64(params.Spread) * params.Point * 10
        fmt.Printf("  %s: %d points (%.1f pips)\n",
            params.Name, params.Spread, spreadPips)
    }
}

// Usage:
// CompareSpreads(account, []string{"EURUSD", "GBPUSD", "EURGBP"})
```

### 5) Find symbols with low spreads

```go
func FindLowSpreadSymbols(account *mt5.MT5Account, symbols []string, maxSpreadPips float64) []string {
    ctx := context.Background()

    data, err := account.SymbolParamsMany(ctx, &pb.SymbolParamsManyRequest{
        Symbols: symbols,
    })
    if err != nil {
        return nil
    }

    lowSpreadSymbols := []string{}

    for _, params := range data.SymbolInfos {
        spreadPips := float64(params.Spread) * params.Point * 10

        if spreadPips <= maxSpreadPips {
            lowSpreadSymbols = append(lowSpreadSymbols, params.Name)
            fmt.Printf("%s: %.1f pips\n", params.Name, spreadPips)
        }
    }

    return lowSpreadSymbols
}

// Usage:
// symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "EURGBP"}
// lowSpread := FindLowSpreadSymbols(account, symbols, 2.0)
```

---

## 📋 Enum Reference

### `BMT5_ENUM_SYMBOL_SECTOR`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `UNDEFINED` | Undefined |
| `1` | `BASIC_MATERIALS` | Basic materials |
| `2` | `COMMUNICATION_SERVICES` | Communication services |
| `3` | `CONSUMER_CYCLICAL` | Consumer cyclical |
| `4` | `CONSUMER_DEFENSIVE` | Consumer defensive |
| `5` | `CURRENCY` | Currencies |
| `6` | `CURRENCY_CRYPTO` | Cryptocurrencies |
| `7` | `ENERGY` | Energy |
| `8` | `FINANCIAL` | Finance |
| `9` | `HEALTHCARE` | Healthcare |
| `10` | `INDUSTRIALS` | Industrials |
| `11` | `REAL_ESTATE` | Real estate |
| `12` | `TECHNOLOGY` | Technology |
| `13` | `UTILITIES` | Utilities |

### `BMT5_ENUM_SYMBOL_INDUSTRY`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `UNDEFINED` | Undefined |
| `1` | `AGRICULTURAL_INPUTS` | Agricultural inputs |
| `2` | `ALUMINIUM` | Aluminium |
| `3` | `BUILDING_MATERIALS` | Building materials |
| `4` | `CHEMICALS` | Chemicals |
| `5` | `COKING_COAL` | Coking coal |
| `6` | `COPPER` | Copper |
| `7` | `GOLD` | Gold |
| `8` | `LUMBER_WOOD` | Lumber and wood production |
| `9` | `INDUSTRIAL_METALS` | Other industrial metals and mining |
| `10` | `PRECIOUS_METALS` | Other precious metals and mining |
| `11` | `PAPER` | Paper and paper products |
| `12` | `SILVER` | Silver |
| `13` | `SPECIALTY_CHEMICALS` | Specialty chemicals |
| `14` | `STEEL` | Steel |
| `15` | `ADVERTISING` | Advertising agencies |
| `16` | `BROADCASTING` | Broadcasting |
| `17` | `GAMING_MULTIMEDIA` | Electronic gaming and multimedia |
| `18` | `ENTERTAINMENT` | Entertainment |
| `19` | `INTERNET_CONTENT` | Internet content and information |
| `20` | `PUBLISHING` | Publishing |
| `21` | `TELECOM` | Telecom services |
| `22` | `APPAREL_MANUFACTURING` | Apparel manufacturing |
| `23` | `APPAREL_RETAIL` | Apparel retail |
| `24` | `AUTO_MANUFACTURERS` | Auto manufacturers |
| `25` | `AUTO_PARTS` | Auto parts |
| `26` | `AUTO_DEALERSHIP` | Auto and truck dealerships |
| `27` | `DEPARTMENT_STORES` | Department stores |
| `28` | `FOOTWEAR_ACCESSORIES` | Footwear and accessories |
| `29` | `FURNISHINGS` | Furnishing, fixtures and appliances |
| `30` | `GAMBLING` | Gambling |
| `31` | `HOME_IMPROV_RETAIL` | Home improvement retail |
| `32` | `INTERNET_RETAIL` | Internet retail |
| `33` | `LEISURE` | Leisure |
| `34` | `LODGING` | Lodging |
| `35` | `LUXURY_GOODS` | Luxury goods |
| `36` | `PACKAGING_CONTAINERS` | Packaging and containers |
| `37` | `PERSONAL_SERVICES` | Personal services |
| `38` | `RECREATIONAL_VEHICLES` | Recreational vehicles |
| `39` | `RESIDENT_CONSTRUCTION` | Residential construction |
| `40` | `RESORTS_CASINOS` | Resorts and casinos |
| `41` | `RESTAURANTS` | Restaurants |
| `42` | `SPECIALTY_RETAIL` | Specialty retail |
| `43` | `TEXTILE_MANUFACTURING` | Textile manufacturing |
| `44` | `TRAVEL_SERVICES` | Travel services |
| `45` | `BEVERAGES_BREWERS` | Beverages - Brewers |
| `46` | `BEVERAGES_NON_ALCO` | Beverages - Non-alcoholic |
| `47` | `BEVERAGES_WINERIES` | Beverages - Wineries and distilleries |
| `48` | `CONFECTIONERS` | Confectioners |
| `49` | `DISCOUNT_STORES` | Discount stores |
| `50` | `EDUCATION_TRAINIG` | Education and training services |
| `51` | `FARM_PRODUCTS` | Farm products |
| `52` | `FOOD_DISTRIBUTION` | Food distribution |
| `53` | `GROCERY_STORES` | Grocery stores |
| `54` | `HOUSEHOLD_PRODUCTS` | Household and personal products |
| `55` | `PACKAGED_FOODS` | Packaged foods |
| `56` | `TOBACCO` | Tobacco |
| `57` | `OIL_GAS_DRILLING` | Oil and gas drilling |
| `58` | `OIL_GAS_EP` | Oil and gas extraction and processing |
| `59` | `OIL_GAS_EQUIPMENT` | Oil and gas equipment and services |
| `60` | `OIL_GAS_INTEGRATED` | Oil and gas integrated |
| `61` | `OIL_GAS_MIDSTREAM` | Oil and gas midstream |
| `62` | `OIL_GAS_REFINING` | Oil and gas refining and marketing |
| `63` | `THERMAL_COAL` | Thermal coal |
| `64` | `URANIUM` | Uranium |
| `65` | `EXCHANGE_TRADED_FUND` | Exchange traded fund |
| `66` | `ASSETS_MANAGEMENT` | Assets management |
| `67` | `BANKS_DIVERSIFIED` | Banks - Diversified |
| `68` | `BANKS_REGIONAL` | Banks - Regional |
| `69` | `CAPITAL_MARKETS` | Capital markets |
| `70` | `CLOSE_END_FUND_DEBT` | Closed-End fund - Debt |
| `71` | `CLOSE_END_FUND_EQUITY` | Closed-end fund - Equity |
| `72` | `CLOSE_END_FUND_FOREIGN` | Closed-end fund - Foreign |
| `73` | `CREDIT_SERVICES` | Credit services |
| `74` | `FINANCIAL_CONGLOMERATE` | Financial conglomerates |
| `75` | `FINANCIAL_DATA_EXCHANGE` | Financial data and stock exchange |
| `76` | `INSURANCE_BROKERS` | Insurance brokers |
| `77` | `INSURANCE_DIVERSIFIED` | Insurance - Diversified |
| `78` | `INSURANCE_LIFE` | Insurance - Life |
| `79` | `INSURANCE_PROPERTY` | Insurance - Property and casualty |
| `80` | `INSURANCE_REINSURANCE` | Insurance - Reinsurance |
| `81` | `INSURANCE_SPECIALTY` | Insurance - Specialty |
| `82` | `MORTGAGE_FINANCE` | Mortgage finance |
| `83` | `SHELL_COMPANIES` | Shell companies |
| `84` | `BIOTECHNOLOGY` | Biotechnology |
| `85` | `DIAGNOSTICS_RESEARCH` | Diagnostics and research |
| `86` | `DRUGS_MANUFACTURERS` | Drugs manufacturers - general |
| `87` | `DRUGS_MANUFACTURERS_SPEC` | Drugs manufacturers - Specialty and generic |
| `88` | `HEALTHCARE_PLANS` | Healthcare plans |
| `89` | `HEALTH_INFORMATION` | Health information services |
| `90` | `MEDICAL_FACILITIES` | Medical care facilities |
| `91` | `MEDICAL_DEVICES` | Medical devices |
| `92` | `MEDICAL_DISTRIBUTION` | Medical distribution |
| `93` | `MEDICAL_INSTRUMENTS` | Medical instruments and supplies |
| `94` | `PHARM_RETAILERS` | Pharmaceutical retailers |
| `95` | `AEROSPACE_DEFENSE` | Aerospace and defense |
| `96` | `AIRLINES` | Airlines |
| `97` | `AIRPORTS_SERVICES` | Airports and air services |
| `98` | `BUILDING_PRODUCTS` | Building products and equipment |
| `99` | `BUSINESS_EQUIPMENT` | Business equipment and supplies |
| `100` | `CONGLOMERATES` | Conglomerates |
| `101` | `CONSULTING_SERVICES` | Consulting services |
| `102` | `ELECTRICAL_EQUIPMENT` | Electrical equipment and parts |
| `103` | `ENGINEERING_CONSTRUCTION` | Engineering and construction |
| `104` | `FARM_HEAVY_MACHINERY` | Farm and heavy construction machinery |
| `105` | `INDUSTRIAL_DISTRIBUTION` | Industrial distribution |
| `106` | `INFRASTRUCTURE_OPERATIONS` | Infrastructure operations |
| `107` | `FREIGHT_LOGISTICS` | Integrated freight and logistics |
| `108` | `MARINE_SHIPPING` | Marine shipping |
| `109` | `METAL_FABRICATION` | Metal fabrication |
| `110` | `POLLUTION_CONTROL` | Pollution and treatment controls |
| `111` | `RAILROADS` | Railroads |
| `112` | `RENTAL_LEASING` | Rental and leasing services |
| `113` | `SECURITY_PROTECTION` | Security and protection services |
| `114` | `SPEALITY_BUSINESS_SERVICES` | Specialty business services |
| `115` | `SPEALITY_MACHINERY` | Specialty industrial machinery |
| `116` | `STUFFING_EMPLOYMENT` | Staffing and employment services |
| `117` | `TOOLS_ACCESSORIES` | Tools and accessories |
| `118` | `TRUCKING` | Trucking |
| `119` | `WASTE_MANAGEMENT` | Waste management |
| `120` | `REAL_ESTATE_DEVELOPMENT` | Real estate - Development |
| `121` | `REAL_ESTATE_DIVERSIFIED` | Real estate - Diversified |
| `122` | `REAL_ESTATE_SERVICES` | Real estate services |
| `123` | `REIT_DIVERSIFIED` | REIT - Diversified |
| `124` | `REIT_HEALTCARE` | REIT - Healthcare facilities |
| `125` | `REIT_HOTEL_MOTEL` | REIT - Hotel and motel |
| `126` | `REIT_INDUSTRIAL` | REIT - Industrial |
| `127` | `REIT_MORTAGE` | REIT - Mortgage |
| `128` | `REIT_OFFICE` | REIT - Office |
| `129` | `REIT_RESIDENTAL` | REIT - Residential |
| `130` | `REIT_RETAIL` | REIT - Retail |
| `131` | `REIT_SPECIALITY` | REIT - Specialty |
| `132` | `COMMUNICATION_EQUIPMENT` | Communication equipment |
| `133` | `COMPUTER_HARDWARE` | Computer hardware |
| `134` | `CONSUMER_ELECTRONICS` | Consumer electronics |
| `135` | `ELECTRONIC_COMPONENTS` | Electronic components |
| `136` | `ELECTRONIC_DISTRIBUTION` | Electronics and computer distribution |
| `137` | `IT_SERVICES` | Information technology services |
| `138` | `SCIENTIFIC_INSTRUMENTS` | Scientific and technical instruments |
| `139` | `SEMICONDUCTOR_EQUIPMENT` | Semiconductor equipment and materials |
| `140` | `SEMICONDUCTORS` | Semiconductors |
| `141` | `SOFTWARE_APPLICATION` | Software - Application |
| `142` | `SOFTWARE_INFRASTRUCTURE` | Software - Infrastructure |
| `143` | `SOLAR` | Solar |
| `144` | `UTILITIES_DIVERSIFIED` | Utilities - Diversified |
| `145` | `UTILITIES_POWERPRODUCERS` | Utilities - Independent power producers |
| `146` | `UTILITIES_RENEWABLE` | Utilities - Renewable |
| `147` | `UTILITIES_REGULATED_ELECTRIC` | Utilities - Regulated electric |
| `148` | `UTILITIES_REGULATED_GAS` | Utilities - Regulated gas |
| `149` | `UTILITIES_REGULATED_WATER` | Utilities - Regulated water |
| `150` | `UTILITIES_FIRST` | Start of utilities services enumeration |
| `151` | `UTILITIES_LAST` | End of utilities services enumeration |

### `BMT5_ENUM_SYMBOL_CHART_MODE`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `CHART_MODE_BID` | Use Bid prices for chart |
| `1` | `CHART_MODE_LAST` | Use Last deal prices for chart |

### `BMT5_ENUM_SYMBOL_CALC_MODE`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `FOREX` | Forex calculation mode |
| `1` | `FOREX_NO_LEVERAGE` | Forex without leverage |
| `2` | `FUTURES` | Futures calculation |
| `3` | `CFD` | CFD calculation |
| `4` | `CFDINDEX` | CFD by indexes |
| `5` | `CFDLEVERAGE` | CFD with leverage |
| `6` | `EXCH_STOCKS` | Exchange stocks |
| `7` | `EXCH_FUTURES` | Exchange futures |
| `8` | `EXCH_FUTURES_FORTS` | FORTS futures |
| `9` | `EXCH_BONDS` | Exchange bonds |
| `10` | `EXCH_STOCKS_MOEX` | MOEX stocks |
| `11` | `EXCH_BONDS_MOEX` | MOEX bonds |
| `12` | `SERV_COLLATERAL` | Collateral mode |

### `BMT5_ENUM_SYMBOL_TRADE_MODE`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `DISABLED` | Trade disabled |
| `1` | `LONGONLY` | Only long positions allowed |
| `2` | `SHORTONLY` | Only short positions allowed |
| `3` | `CLOSEONLY` | Only close operations allowed |
| `4` | `FULL` | No trade restrictions |

### `BMT5_ENUM_SYMBOL_TRADE_EXECUTION`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `REQUEST` | Execution by request |
| `1` | `INSTANT` | Instant execution |
| `2` | `MARKET` | Market execution |
| `3` | `EXCHANGE` | Exchange execution |

### `BMT5_ENUM_SYMBOL_SWAP_MODE`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `DISABLED` | No swaps |
| `1` | `POINTS` | Swaps in points |
| `2` | `CURRENCY_SYMBOL` | Swaps in base currency |
| `3` | `CURRENCY_MARGIN` | Swaps in margin currency |
| `4` | `CURRENCY_DEPOSIT` | Swaps in deposit currency |
| `5` | `CURRENCY_PROFIT` | Swaps in profit currency |
| `6` | `INTEREST_CURRENT` | Annual interest from current price |
| `7` | `INTEREST_OPEN` | Annual interest from open price |
| `8` | `REOPEN_CURRENT` | Reopen by close price |
| `9` | `REOPEN_BID` | Reopen by Bid price |

### `BMT5_ENUM_DAY_OF_WEEK`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `SUNDAY` | Sunday |
| `1` | `MONDAY` | Monday |
| `2` | `TUESDAY` | Tuesday |
| `3` | `WEDNESDAY` | Wednesday |
| `4` | `THURSDAY` | Thursday |
| `5` | `FRIDAY` | Friday |
| `6` | `SATURDAY` | Saturday |

### `BMT5_ENUM_SYMBOL_ORDER_GTC_MODE`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `GTC` | Good till cancelled - unlimited validity |
| `1` | `DAILY` | Valid for one trading day |
| `2` | `DAILY_EXCLUDING_STOPS` | Daily, but SL/TP preserved |

### `BMT5_ENUM_SYMBOL_OPTION_MODE`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `EUROPEAN` | European option |
| `1` | `AMERICAN` | American option |

### `BMT5_ENUM_SYMBOL_OPTION_RIGHT`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `CALL` | Call option (right to buy) |
| `1` | `PUT` | Put option (right to sell) |

### `BMT5_ENUM_ORDER_TYPE`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `BUY` | Market buy |
| `1` | `SELL` | Market sell |
| `2` | `BUY_LIMIT` | Buy limit pending |
| `3` | `SELL_LIMIT` | Sell limit pending |
| `4` | `BUY_STOP` | Buy stop pending |
| `5` | `SELL_STOP` | Sell stop pending |
| `6` | `BUY_STOP_LIMIT` | Buy stop limit |
| `7` | `SELL_STOP_LIMIT` | Sell stop limit |
| `8` | `CLOSE_BY` | Close by opposite |

### `BMT5_ENUM_ORDER_TYPE_FILLING`
| Value | Enum Name | Description |
|-------|-----------|-------------|
| `0` | `FOK` | Fill or Kill - all or nothing |
| `1` | `IOC` | Immediate or Cancel - partial fill allowed |
| `2` | `RETURN` | Return order with remaining volume |
| `3` | `BOC` | Book or Cancel - only passive execution |

---

## 📚 See Also

* [SymbolInfoDouble](../2.%20Symbol_information/SymbolInfoDouble.md) - Get individual symbol properties
* [SymbolInfoInteger](../2.%20Symbol_information/SymbolInfoInteger.md) - Get integer symbol properties
* [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - Get tick data
