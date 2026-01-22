# 📊 Getting Account Double Properties

> **Request:** specific double-type account property from **MT5** terminal using property identifier.

**API Information:**

* **Low-level API:** `MT5Account.AccountInfoDouble(...)` (from Go package `github.com/MetaRPC/GoMT5/package/Helpers`)
* **gRPC service:** `mt5_term_api.AccountInformation`
* **Proto definition:** `AccountInfoDouble` (defined in `mt5-term-api-account-information.proto`)

### RPC

* **Service:** `mt5_term_api.AccountInformation`
* **Method:** `AccountInfoDouble(AccountInfoDoubleRequest) → AccountInfoDoubleReply`
* **Low‑level client (generated):** `AccountInformationClient.AccountInfoDouble(ctx, request, opts...)`

## 💬 Just the essentials

* **What it is.** Retrieve a specific double-type account property by ID (balance, equity, margin, etc.).
* **Why you need it.** Get individual account metrics without fetching all data. Useful for focused checks.
* **When to use.** Use `AccountSummary()` for multiple properties. Use this method for single property queries.

---

## 🎯 Purpose

Use it to query specific numeric account properties:

* Check account balance before trading
* Monitor margin usage and margin level
* Verify free margin availability
* Track floating profit/loss
* Monitor margin call and stop out levels

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [AccountInfoDouble - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoDouble_HOW.md)**

---



```go
package mt5

type MT5Account struct {
    // ...
}

// AccountInfoDouble retrieves a double-type account property.
// Returns AccountInfoDoubleData with the requested numeric value.
func (a *MT5Account) AccountInfoDouble(
    ctx context.Context,
    req *pb.AccountInfoDoubleRequest,
) (*pb.AccountInfoDoubleData, error)
```

**Request message:**

```protobuf
message AccountInfoDoubleRequest {
  AccountInfoDoublePropertyType property_id = 1;
}
```

**Reply message:**

```protobuf
message AccountInfoDoubleReply {
  oneof response {
    AccountInfoDoubleData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter     | Type                              | Description                                             |
| ------------- | --------------------------------- | ------------------------------------------------------- |
| `ctx`         | `context.Context`                 | Context for deadline/timeout and cancellation           |
| `req`         | `*pb.AccountInfoDoubleRequest`    | Request with property_id specifying which property to retrieve |
| `req.property_id` | `AccountInfoDoublePropertyType` (enum) | Property identifier (see enum below) |

**Context options:**

```go
// 1. With timeout (recommended)
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

// 2. With deadline
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
defer cancel()

// 3. With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

---

## ⬆️ Output — `AccountInfoDoubleData`

| Field             | Type      | Go Type   | Description                                      |
| ----------------- | --------- | --------- | ------------------------------------------------ |
| `requested_value` | `double`  | `float64` | The value of the requested account property      |

---

## 🧱 Related enums (from proto)

### `AccountInfoDoublePropertyType`

Defined in `mt5-term-api-account-information.proto`:

```protobuf
enum AccountInfoDoublePropertyType {
  ACCOUNT_BALANCE = 0;           // Account balance in the deposit currency
  ACCOUNT_CREDIT = 1;            // Account credit in the deposit currency
  ACCOUNT_PROFIT = 2;            // Current profit of an account in the deposit currency
  ACCOUNT_EQUITY = 3;            // Account equity in the deposit currency
  ACCOUNT_MARGIN = 4;            // Account margin used in the deposit currency
  ACCOUNT_MARGIN_FREE = 5;       // Free margin of an account in the deposit currency
  ACCOUNT_MARGIN_LEVEL = 6;      // Account margin level in percents
  ACCOUNT_MARGIN_SO_CALL = 7;    // Margin call level (in % or deposit currency)
  ACCOUNT_MARGIN_SO_SO = 8;      // Margin stop out level (in % or deposit currency)
  ACCOUNT_MARGIN_INITIAL = 9;    // Initial margin (reserved for pending orders)
  ACCOUNT_MARGIN_MAINTENANCE = 10; // Maintenance margin (minimum equity reserved)
  ACCOUNT_ASSETS = 11;           // Current assets of an account
  ACCOUNT_LIABILITIES = 12;      // Current liabilities on an account
  ACCOUNT_COMMISSION_BLOCKED = 13; // Current blocked commission amount
}
```

**Go constants (generated):**

```go
const (
    AccountInfoDoublePropertyType_ACCOUNT_BALANCE           = 0
    AccountInfoDoublePropertyType_ACCOUNT_CREDIT            = 1
    AccountInfoDoublePropertyType_ACCOUNT_PROFIT            = 2
    AccountInfoDoublePropertyType_ACCOUNT_EQUITY            = 3
    AccountInfoDoublePropertyType_ACCOUNT_MARGIN            = 4
    AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE       = 5
    AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL      = 6
    AccountInfoDoublePropertyType_ACCOUNT_MARGIN_SO_CALL    = 7
    AccountInfoDoublePropertyType_ACCOUNT_MARGIN_SO_SO      = 8
    AccountInfoDoublePropertyType_ACCOUNT_MARGIN_INITIAL    = 9
    AccountInfoDoublePropertyType_ACCOUNT_MARGIN_MAINTENANCE = 10
    AccountInfoDoublePropertyType_ACCOUNT_ASSETS            = 11
    AccountInfoDoublePropertyType_ACCOUNT_LIABILITIES       = 12
    AccountInfoDoublePropertyType_ACCOUNT_COMMISSION_BLOCKED = 13
)
```

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `3s` timeout is applied automatically.
* **Prefer AccountSummary:** For multiple properties, use `AccountSummary()` instead to avoid multiple round-trips.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.
* **Property types:** Margin call and stop out levels can be in percents or deposit currency depending on account settings.

---

## 🔗 Usage Examples

### 1) Get account balance

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/MetaRPC/GoMT5/package/Helpers"
    "github.com/google/uuid"
)

func main() {
    account, err := mt5.NewMT5Account(591129415, "password", "server:443", uuid.New())
    if err != nil {
        panic(err)
    }
    defer account.Close()

    // Connect first
    err = account.Connect()
    if err != nil {
        panic(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Get balance
    req := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
    }

    data, err := account.AccountInfoDouble(ctx, req)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Account Balance: %.2f\n", data.RequestedValue)
}
```

### 2) Get current equity

```go
func GetEquity(account *mt5.MT5Account) (float64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_EQUITY,
    }

    data, err := account.AccountInfoDouble(ctx, req)
    if err != nil {
        return 0, fmt.Errorf("failed to get equity: %w", err)
    }

    return data.RequestedValue, nil
}

// Usage:
equity, err := GetEquity(account)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Current Equity: %.2f\n", equity)
```

### 3) Check margin level before trading

```go
func CheckMarginLevel(account *mt5.MT5Account, minLevel float64) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL,
    }

    data, err := account.AccountInfoDouble(ctx, req)
    if err != nil {
        return false, err
    }

    marginLevel := data.RequestedValue
    fmt.Printf("Margin Level: %.2f%%\n", marginLevel)

    // Check if margin level is above minimum
    if marginLevel < minLevel {
        return false, fmt.Errorf("margin level %.2f%% is below minimum %.2f%%",
            marginLevel, minLevel)
    }

    return true, nil
}

// Usage:
// Check if margin level is at least 200%
ok, err := CheckMarginLevel(account, 200.0)
if !ok {
    log.Println("Insufficient margin level, cannot trade")
}
```

### 4) Get free margin

```go
func GetFreeMargin(account *mt5.MT5Account) (float64, error) {
    ctx := context.Background() // Will use default 3s timeout

    req := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE,
    }

    data, err := account.AccountInfoDouble(ctx, req)
    if err != nil {
        return 0, err
    }

    return data.RequestedValue, nil
}

// Usage:
freeMargin, _ := GetFreeMargin(account)
fmt.Printf("Free Margin: %.2f\n", freeMargin)
```

### 5) Monitor floating profit/loss

```go
func MonitorProfit(account *mt5.MT5Account, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

        req := &pb.AccountInfoDoubleRequest{
            PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_PROFIT,
        }

        data, err := account.AccountInfoDouble(ctx, req)
        cancel()

        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        profit := data.RequestedValue
        sign := "+"
        if profit < 0 {
            sign = ""
        }

        fmt.Printf("[%s] Floating P/L: %s%.2f\n",
            time.Now().Format("15:04:05"),
            sign,
            profit)
    }
}

// Usage:
// MonitorProfit(account, 5*time.Second) // Update every 5 seconds
```

### 6) Check multiple properties sequentially

```go
type MarginInfo struct {
    UsedMargin  float64
    FreeMargin  float64
    MarginLevel float64
}

func GetMarginInfo(account *mt5.MT5Account) (*MarginInfo, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    info := &MarginInfo{}

    // Get used margin
    usedReq := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN,
    }
    usedData, err := account.AccountInfoDouble(ctx, usedReq)
    if err != nil {
        return nil, err
    }
    info.UsedMargin = usedData.RequestedValue

    // Get free margin
    freeReq := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE,
    }
    freeData, err := account.AccountInfoDouble(ctx, freeReq)
    if err != nil {
        return nil, err
    }
    info.FreeMargin = freeData.RequestedValue

    // Get margin level
    levelReq := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL,
    }
    levelData, err := account.AccountInfoDouble(ctx, levelReq)
    if err != nil {
        return nil, err
    }
    info.MarginLevel = levelData.RequestedValue

    return info, nil
}

// Usage:
info, err := GetMarginInfo(account)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Margin Info:\n")
fmt.Printf("  Used: %.2f\n", info.UsedMargin)
fmt.Printf("  Free: %.2f\n", info.FreeMargin)
fmt.Printf("  Level: %.2f%%\n", info.MarginLevel)

// Note: For better performance, use AccountSummary() instead
// to get all properties in one call
```

---

## 🔧 Common Patterns

### Safe trading check

```go
func CanTrade(account *mt5.MT5Account, requiredMargin float64) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Check free margin
    req := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE,
    }

    data, err := account.AccountInfoDouble(ctx, req)
    if err != nil {
        return false, err
    }

    freeMargin := data.RequestedValue

    // Check if we have enough free margin
    if freeMargin < requiredMargin {
        return false, fmt.Errorf("insufficient margin: need %.2f, have %.2f",
            requiredMargin, freeMargin)
    }

    return true, nil
}
```

### Margin level warning system

```go
func CheckMarginWarning(account *mt5.MT5Account) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoDoubleRequest{
        PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL,
    }

    data, err := account.AccountInfoDouble(ctx, req)
    if err != nil {
        return err
    }

    level := data.RequestedValue

    // Check margin level thresholds
    switch {
    case level < 100:
        return fmt.Errorf("CRITICAL: Margin level %.2f%% - Stop out imminent!", level)
    case level < 200:
        log.Printf("WARNING: Margin level %.2f%% - Margin call zone", level)
    case level < 500:
        log.Printf("CAUTION: Margin level %.2f%% - Monitor closely", level)
    default:
        log.Printf("OK: Margin level %.2f%%", level)
    }

    return nil
}
```

---

## 📚 See Also

* [AccountSummary](./AccountSummary.md) - Get all account properties in one call (RECOMMENDED)
* [AccountInfoInteger](./AccountInfoInteger.md) - Get specific integer account properties
* [AccountInfoString](./AccountInfoString.md) - Get specific string account properties
