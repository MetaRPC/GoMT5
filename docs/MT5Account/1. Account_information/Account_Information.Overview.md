# MT5Account · Account Information - Overview

> Account balance, equity, margin, leverage, currency, and other account properties. Use this page to choose the right API for accessing account state.

## 📁 What lives here

* **[AccountSummary](./AccountSummary.md)** - **all account info** at once (balance, equity, margin, leverage, profit, etc.).
* **[AccountInfoDouble](./AccountInfoDouble.md)** - **single double value** from account (balance, equity, margin, profit, credit, etc.).
* **[AccountInfoInteger](./AccountInfoInteger.md)** - **single integer value** from account (login, leverage, limit orders, etc.).
* **[AccountInfoString](./AccountInfoString.md)** - **single string value** from account (name, server, currency, company).

---

## 📚 Step-by-step tutorials

Want detailed explanations with line-by-line code breakdown? Check these guides:

* **[AccountSummary - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountSummary_HOW.md)**
* **[AccountInfoDouble - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoDouble_HOW.md)**
* **[AccountInfoInteger - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoInteger_HOW.md)**
* **[AccountInfoString - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoString_HOW.md)**

---

## 🧭 Plain English

* **AccountSummary** → the **one-stop shop** for complete account snapshot (balance, equity, margin level, etc.).
* **AccountInfoDouble** → grab **one numeric property** when you need just balance or margin.
* **AccountInfoInteger** → grab **one integer property** like login number or leverage.
* **AccountInfoString** → grab **one text property** like account name or currency.

> Rule of thumb: need **full snapshot** → `AccountSummary`; need **one specific value** → `AccountInfo*` (Double/Integer/String).

---

## Quick choose

| If you need…                                     | Use                       | Returns                    | Key inputs                          |
| ------------------------------------------------ | ------------------------- | -------------------------- | ----------------------------------- |
| Complete account snapshot (all values)           | `AccountSummary`          | AccountSummaryData         | *(none)*                            |
| One numeric value (balance, equity, margin, etc.)| `AccountInfoDouble`       | Single `float64`           | Property enum (BALANCE, EQUITY, etc.) |
| One integer value (login, leverage, etc.)        | `AccountInfoInteger`      | Single `int64`             | Property enum (LOGIN, LEVERAGE, etc.) |
| One text value (name, currency, server, etc.)    | `AccountInfoString`       | Single `string`            | Property enum (NAME, CURRENCY, etc.) |

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

## See also

* **Symbol Information:** [SymbolInfoTick](../2.%20Symbol_information/SymbolInfoTick.md) - get current prices
* **Trading:** [OrderSend](../4.%20Trading_Operations/OrderSend.md) - place orders
* **Positions:** [PositionsTotal](../3.%20Position_Orders_Information/PositionsTotal.md) - count open positions
