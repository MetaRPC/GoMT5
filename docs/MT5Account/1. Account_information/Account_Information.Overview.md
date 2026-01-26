# MT5Account · Account Information - Overview

> Account balance, equity, margin, leverage, currency, and other account properties. Use this page to choose the right API for accessing account state.

## 📁 What lives here

* **[AccountSummary](./AccountSummary.md)** - **all account info** at once (balance, equity, margin, leverage, profit, etc.). Returns `MrpcEnumAccountTradeMode` enum.
* **[AccountInfoDouble](./AccountInfoDouble.md)** - **single double value** from account (balance, equity, margin, profit, credit, etc.). Uses `AccountInfoDoublePropertyType` enum as input.
* **[AccountInfoInteger](./AccountInfoInteger.md)** - **single integer value** from account (login, leverage, limit orders, etc.). Uses `AccountInfoIntegerPropertyType` enum as input.
* **[AccountInfoString](./AccountInfoString.md)** - **single string value** from account (name, server, currency, company). Uses `AccountInfoStringPropertyType` enum as input.

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[AccountSummary - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountSummary_HOW.md)**
* **[AccountInfoDouble - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoDouble_HOW.md)**
* **[AccountInfoInteger - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoInteger_HOW.md)**
* **[AccountInfoString - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoString_HOW.md)**

---

## 🧭 Plain English

* **AccountSummary** → the **one-stop shop** for complete account snapshot (balance, equity, margin level, etc.). Returns trade mode as ENUM.
* **AccountInfoDouble** → grab **one numeric property** when you need just balance or margin. Requires property ENUM as input.
* **AccountInfoInteger** → grab **one integer property** like login number or leverage. Requires property ENUM as input.
* **AccountInfoString** → grab **one text property** like account name or currency. Requires property ENUM as input.

> Rule of thumb: need **full snapshot** → `AccountSummary`; need **one specific value** → `AccountInfo*` (Double/Integer/String).
>
> **🧱 ENUM Note:** All 4 methods use ENUMs! `AccountInfo*` methods require property ENUMs as input to select which property to retrieve. `AccountSummary` returns account trade mode (DEMO/REAL/CONTEST) as an output ENUM.

---

## Quick choose

| If you need…                                     | Use                       | Returns                    | ENUMs Used                          |
| ------------------------------------------------ | ------------------------- | -------------------------- | ----------------------------------- |
| Complete account snapshot (all values)           | `AccountSummary`          | AccountSummaryData         | **Output**: `MrpcEnumAccountTradeMode` |
| One numeric value (balance, equity, margin, etc.)| `AccountInfoDouble`       | Single `float64`           | **Input**: `AccountInfoDoublePropertyType` |
| One integer value (login, leverage, etc.)        | `AccountInfoInteger`      | Single `int64`             | **Input**: `AccountInfoIntegerPropertyType` |
| One text value (name, currency, server, etc.)    | `AccountInfoString`       | Single `string`            | **Input**: `AccountInfoStringPropertyType` |

---

## ❌ Cross‑refs & gotchas

* **Margin Level** - use `AccountInfoDouble(ACCOUNT_MARGIN_LEVEL)` to get as percentage.
* **Free Margin** - use `AccountInfoDouble(ACCOUNT_MARGIN_FREE)` for available margin.
* **AccountSummary** includes basic info (balance, equity, leverage); for margin use `AccountInfoDouble`.
* **AccountInfo*** methods are lighter if you only need one property.
* **Currency** affects how profits are calculated - always check account currency.
* **Leverage** determines margin requirements - higher leverage = less margin needed.
* **Context timeout** - Remember to set appropriate context deadline for gRPC calls.

---

## 🟢 Minimal snippets

```go
// Get complete account snapshot
summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Balance: $%.2f, Equity: $%.2f, Leverage: 1:%d\n",
    summary.AccountBalance,
    summary.AccountEquity,
    summary.AccountLeverage)
```

```go
// Get single property - account balance
balanceData, err := account.AccountInfoDouble(
    ctx,
    &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
    },
)
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Balance: $%.2f\n", balanceData.RequestedValue)
```

```go
// Get account leverage
leverageData, err := account.AccountInfoInteger(
    ctx,
    &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE,
    },
)
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Leverage: 1:%d\n", leverageData.RequestedValue)
```

```go
// Get account currency
currencyData, err := account.AccountInfoString(
    ctx,
    &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY,
    },
)
if err != nil {
    log.Fatalf("Failed: %v", err)
}

fmt.Printf("Currency: %s\n", currencyData.RequestedValue)
```

```go
// Check account health
summary, _ := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})

// Get margin level (AccountSummary doesn't include margin, need separate call)
marginData, err := account.AccountInfoDouble(ctx, &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL,
})
if err != nil {
    log.Fatalf("Failed: %v", err)
}

marginLevel := marginData.RequestedValue

if marginLevel < 100 {
    fmt.Println("⚠️ Warning: Margin level critical!")
} else if marginLevel < 200 {
    fmt.Println("⚠️ Warning: Low margin level")
} else {
    fmt.Println("✅ Healthy margin level")
}
```

---

## 🧱 ENUMs used in Account Information

All 4 methods in this group use ENUMs to specify which property to retrieve or to represent account state:

### Input ENUMs (Property Selectors)

These ENUMs are used to specify **which property** you want to retrieve:

| Method | Input ENUM | Purpose | Values |
|--------|------------|---------|--------|
| **AccountInfoDouble** | `AccountInfoDoublePropertyType` | Select which double property to get | `ACCOUNT_BALANCE`, `ACCOUNT_EQUITY`, `ACCOUNT_MARGIN`, `ACCOUNT_MARGIN_FREE`, `ACCOUNT_MARGIN_LEVEL`, `ACCOUNT_PROFIT`, `ACCOUNT_CREDIT`, etc. (14 total) |
| **AccountInfoInteger** | `AccountInfoIntegerPropertyType` | Select which integer property to get | `ACCOUNT_LOGIN`, `ACCOUNT_LEVERAGE`, `ACCOUNT_TRADE_MODE`, `ACCOUNT_LIMIT_ORDERS`, `ACCOUNT_TRADE_ALLOWED`, `ACCOUNT_TRADE_EXPERT`, etc. (11 total) |
| **AccountInfoString** | `AccountInfoStringPropertyType` | Select which string property to get | `ACCOUNT_NAME`, `ACCOUNT_SERVER`, `ACCOUNT_CURRENCY`, `ACCOUNT_COMPANY` (4 total) |

**Usage in code:**

```go
// Example 1: Get balance using AccountInfoDoublePropertyType enum
req := &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
}
data, _ := account.AccountInfoDouble(ctx, req)
fmt.Printf("Balance: %.2f\n", data.RequestedValue)

// Example 2: Get leverage using AccountInfoIntegerPropertyType enum
req := &pb.AccountInfoIntegerRequest{
    PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE,
}
data, _ := account.AccountInfoInteger(ctx, req)
fmt.Printf("Leverage: 1:%d\n", data.RequestedValue)

// Example 3: Get currency using AccountInfoStringPropertyType enum
req := &pb.AccountInfoStringRequest{
    PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY,
}
data, _ := account.AccountInfoString(ctx, req)
fmt.Printf("Currency: %s\n", data.RequestedValue)
```

### Output ENUMs (Return Values)

These ENUMs are **returned in the response** to represent account state:

| Method | Output ENUM | Field Name | Purpose | Values |
|--------|-------------|------------|---------|--------|
| **AccountSummary** | `MrpcEnumAccountTradeMode` | `AccountTradeMode` | Indicates account type | `MRPC_ACCOUNT_TRADE_MODE_DEMO` (0)<br>`MRPC_ACCOUNT_TRADE_MODE_CONTEST` (1)<br>`MRPC_ACCOUNT_TRADE_MODE_REAL` (2) |
| **AccountInfoInteger** | Various ENUMs returned as int64 | `requested_value` | When querying `ACCOUNT_TRADE_MODE` property, value represents account type | Same as above (0=DEMO, 1=CONTEST, 2=REAL) |

**Usage in code:**

```go
// Example 1: AccountSummary returns trade mode as enum
summary, _ := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})

switch summary.AccountTradeMode {
case pb.MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_DEMO:
    fmt.Println("Demo account")
case pb.MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_REAL:
    fmt.Println("Real account")
}

// Example 2: AccountInfoInteger can return ENUM values as int64
req := &pb.AccountInfoIntegerRequest{
    PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_MODE,
}
data, _ := account.AccountInfoInteger(ctx, req)

// data.RequestedValue will be 0 (DEMO), 1 (CONTEST), or 2 (REAL)
switch data.RequestedValue {
case 0:
    fmt.Println("Demo")
case 2:
    fmt.Println("Real")
}
```

### 📝 Important Notes about ENUMs

* **Input ENUMs are required** - You must specify which property to retrieve using the appropriate enum constant
* **Full enum names in Go** - Always use the full prefixed name: `pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE`
* **Output ENUMs** - AccountSummary directly returns `MrpcEnumAccountTradeMode` as a typed field
* **AccountInfoInteger special case** - Some integer properties return values that represent ENUMs (like TRADE_MODE), but they're returned as `int64`, not typed enums
* **No ENUMs for AccountInfoDouble/String outputs** - These always return simple `float64` or `string` values

---

## See also

* **Symbol Information:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get current prices
* **Trading:** [OrderSend](../4.%20Trading_Operations/OrderSend.md) - place orders
* **Positions:** [PositionsTotal](../3.%20Position_Orders_Information/PositionsTotal.md) - count open positions
