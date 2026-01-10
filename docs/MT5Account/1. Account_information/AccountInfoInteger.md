# 🔢 Getting Account Integer Properties

> **Request:** specific integer-type account property from **MT5** terminal using property identifier.

**API Information:**

* **SDK wrapper:** `MT5Account.AccountInfoInteger(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.AccountInformation`
* **Proto definition:** `AccountInfoInteger` (defined in `mt5-term-api-account-information.proto`)

### RPC

* **Service:** `mt5_term_api.AccountInformation`
* **Method:** `AccountInfoInteger(AccountInfoIntegerRequest) → AccountInfoIntegerReply`
* **Low‑level client (generated):** `AccountInformationClient.AccountInfoInteger(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**


## 💬 Just the essentials

* **What it is.** Retrieve a specific integer-type account property by ID (login, leverage, trade mode, limits).
* **Why you need it.** Get account configuration and permission settings.
* **When to use.** Query individual settings. For multiple properties, use `AccountSummary()`.

---

## 🎯 Purpose

Use it to query specific integer account properties:

* Get account login number
* Check account leverage
* Verify trading permissions
* Check position limits
* Validate account configuration

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [AccountInfoInteger - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoInteger_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// AccountInfoInteger retrieves an integer-type account property.
// Returns AccountInfoIntegerData with the requested int64 value.
func (a *MT5Account) AccountInfoInteger(
    ctx context.Context,
    req *pb.AccountInfoIntegerRequest,
) (*pb.AccountInfoIntegerData, error)
```

**Request message:**

```protobuf
message AccountInfoIntegerRequest {
  AccountInfoIntegerPropertyType property_id = 1;
}
```

**Reply message:**

```protobuf
message AccountInfoIntegerReply {
  oneof response {
    AccountInfoIntegerData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter     | Type                              | Description                                             |
| ------------- | --------------------------------- | ------------------------------------------------------- |
| `ctx`         | `context.Context`                 | Context for deadline/timeout and cancellation           |
| `req`         | `*pb.AccountInfoIntegerRequest`   | Request with property_id specifying which property to retrieve |
| `req.property_id` | `AccountInfoIntegerPropertyType` (enum) | Property identifier (see enum below) |

---

## ⬆️ Output — `AccountInfoIntegerData`

| Field             | Type    | Go Type | Description                                      |
| ----------------- | ------- | ------- | ------------------------------------------------ |
| `requested_value` | `int64` | `int64` | The value of the requested account property      |

---

## 🧱 Related enums (from proto)

### `AccountInfoIntegerPropertyType`

Defined in `mt5-term-api-account-information.proto`:

```protobuf
enum AccountInfoIntegerPropertyType {
  ACCOUNT_LOGIN = 0;          // Account number (login)
  ACCOUNT_TRADE_MODE = 1;     // Account trade mode (DEMO/CONTEST/REAL)
  ACCOUNT_LEVERAGE = 2;       // Account leverage (e.g., 100 for 1:100)
  ACCOUNT_LIMIT_ORDERS = 3;   // Maximum allowed open positions + pending orders (0 = unlimited)
  ACCOUNT_MARGIN_SO_MODE = 4; // Margin stop out mode (percent or deposit currency)
  ACCOUNT_TRADE_ALLOWED = 5;  // Trading allowed for this account (1 or 0)
  ACCOUNT_TRADE_EXPERT = 6;   // Expert Advisor trading allowed (1 or 0)
  ACCOUNT_MARGIN_MODE = 7;    // Margin calculation mode
  ACCOUNT_CURRENCY_DIGITS = 8; // Decimal places for account currency
  ACCOUNT_FIFO_CLOSE = 9;     // Positions must be closed in FIFO order (1 or 0)
  ACCOUNT_HEDGE_ALLOWED = 10; // Opposite positions on same symbol allowed (1 or 0)
}
```

**Go constants (generated):**

```go
const (
    AccountInfoIntegerPropertyType_ACCOUNT_LOGIN          = 0
    AccountInfoIntegerPropertyType_ACCOUNT_TRADE_MODE     = 1
    AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE       = 2
    AccountInfoIntegerPropertyType_ACCOUNT_LIMIT_ORDERS   = 3
    AccountInfoIntegerPropertyType_ACCOUNT_MARGIN_SO_MODE = 4
    AccountInfoIntegerPropertyType_ACCOUNT_TRADE_ALLOWED  = 5
    AccountInfoIntegerPropertyType_ACCOUNT_TRADE_EXPERT   = 6
    AccountInfoIntegerPropertyType_ACCOUNT_MARGIN_MODE    = 7
    AccountInfoIntegerPropertyType_ACCOUNT_CURRENCY_DIGITS = 8
    AccountInfoIntegerPropertyType_ACCOUNT_FIFO_CLOSE     = 9
    AccountInfoIntegerPropertyType_ACCOUNT_HEDGE_ALLOWED  = 10
)
```

---

## 🧩 Notes & Tips

* **Automatic reconnection:** Built-in protection against transient gRPC errors.
* **Default timeout:** Default `3s` timeout applied if context has no deadline.
* **Boolean values:** Properties like TRADE_ALLOWED return `1` (true) or `0` (false).
* **Prefer AccountSummary:** For multiple properties, use `AccountSummary()` to avoid multiple round-trips.

---

## 🔗 Usage Examples

### 1) Get account login number

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
    "github.com/MetaRPC/GoMT5/mt5"
    "github.com/google/uuid"
)

func main() {
    account, err := mt5.NewMT5Account(591129415, "password", "server:443", uuid.New())
    if err != nil {
        panic(err)
    }
    defer account.Close()

    err = account.Connect()
    if err != nil {
        panic(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LOGIN,
    }

    data, err := account.AccountInfoInteger(ctx, req)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Account Login: %d\n", data.RequestedValue)
}
```

### 2) Check account leverage

```go
func GetLeverage(account *mt5.MT5Account) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE,
    }

    data, err := account.AccountInfoInteger(ctx, req)
    if err != nil {
        return 0, fmt.Errorf("failed to get leverage: %w", err)
    }

    return data.RequestedValue, nil
}

// Usage:
leverage, _ := GetLeverage(account)
fmt.Printf("Account Leverage: 1:%d\n", leverage)
// Output: Account Leverage: 1:100
```

### 3) Check if trading is allowed

```go
func IsTradingAllowed(account *mt5.MT5Account) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_ALLOWED,
    }

    data, err := account.AccountInfoInteger(ctx, req)
    if err != nil {
        return false, err
    }

    // Value is 1 for true, 0 for false
    return data.RequestedValue == 1, nil
}

// Usage:
allowed, err := IsTradingAllowed(account)
if err != nil {
    log.Fatal(err)
}

if !allowed {
    log.Println("Trading is disabled for this account")
}
```

### 4) Check if Expert Advisor trading is allowed

```go
func IsEAAllowed(account *mt5.MT5Account) bool {
    ctx := context.Background()

    req := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_EXPERT,
    }

    data, err := account.AccountInfoInteger(ctx, req)
    if err != nil {
        return false
    }

    return data.RequestedValue == 1
}

// Usage:
if IsEAAllowed(account) {
    fmt.Println("Expert Advisor trading is enabled")
} else {
    fmt.Println("Expert Advisor trading is disabled")
}
```

### 5) Get account trade mode

```go
func GetTradeMode(account *mt5.MT5Account) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_MODE,
    }

    data, err := account.AccountInfoInteger(ctx, req)
    if err != nil {
        return "", err
    }

    // Convert to string
    switch data.RequestedValue {
    case 0:
        return "DEMO", nil
    case 1:
        return "CONTEST", nil
    case 2:
        return "REAL", nil
    default:
        return fmt.Sprintf("UNKNOWN(%d)", data.RequestedValue), nil
    }
}

// Usage:
mode, _ := GetTradeMode(account)
fmt.Printf("Account Type: %s\n", mode)
```

### 6) Check position limits

```go
func CheckPositionLimits(account *mt5.MT5Account) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LIMIT_ORDERS,
    }

    data, err := account.AccountInfoInteger(ctx, req)
    if err != nil {
        return err
    }

    limit := data.RequestedValue

    if limit == 0 {
        fmt.Println("No limit on positions and orders")
    } else {
        fmt.Printf("Maximum allowed positions + orders: %d\n", limit)
    }

    return nil
}

// Usage:
CheckPositionLimits(account)
```

---

## 🔧 Common Patterns

### Validate account configuration

```go
type AccountConfig struct {
    Login        int64
    Leverage     int64
    TradeMode    string
    TradeAllowed bool
    EAAllowed    bool
    HedgeAllowed bool
}

func GetAccountConfig(account *mt5.MT5Account) (*AccountConfig, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    config := &AccountConfig{}

    // Get login
    loginReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LOGIN,
    }
    loginData, err := account.AccountInfoInteger(ctx, loginReq)
    if err != nil {
        return nil, err
    }
    config.Login = loginData.RequestedValue

    // Get leverage
    levReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE,
    }
    levData, err := account.AccountInfoInteger(ctx, levReq)
    if err != nil {
        return nil, err
    }
    config.Leverage = levData.RequestedValue

    // Get trade mode
    modeReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_MODE,
    }
    modeData, err := account.AccountInfoInteger(ctx, modeReq)
    if err != nil {
        return nil, err
    }

    switch modeData.RequestedValue {
    case 0:
        config.TradeMode = "DEMO"
    case 1:
        config.TradeMode = "CONTEST"
    case 2:
        config.TradeMode = "REAL"
    }

    // Get trading permissions
    tradeReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_ALLOWED,
    }
    tradeData, err := account.AccountInfoInteger(ctx, tradeReq)
    if err != nil {
        return nil, err
    }
    config.TradeAllowed = tradeData.RequestedValue == 1

    eaReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_EXPERT,
    }
    eaData, err := account.AccountInfoInteger(ctx, eaReq)
    if err != nil {
        return nil, err
    }
    config.EAAllowed = eaData.RequestedValue == 1

    hedgeReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_HEDGE_ALLOWED,
    }
    hedgeData, err := account.AccountInfoInteger(ctx, hedgeReq)
    if err != nil {
        return nil, err
    }
    config.HedgeAllowed = hedgeData.RequestedValue == 1

    return config, nil
}

// Usage:
config, err := GetAccountConfig(account)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Account Configuration:\n")
fmt.Printf("  Login: %d\n", config.Login)
fmt.Printf("  Leverage: 1:%d\n", config.Leverage)
fmt.Printf("  Type: %s\n", config.TradeMode)
fmt.Printf("  Trading: %v\n", config.TradeAllowed)
fmt.Printf("  EA Trading: %v\n", config.EAAllowed)
fmt.Printf("  Hedging: %v\n", config.HedgeAllowed)
```

### Pre-flight trading checks

```go
func CanStartTrading(account *mt5.MT5Account) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check if trading is allowed
    tradeReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_ALLOWED,
    }
    tradeData, err := account.AccountInfoInteger(ctx, tradeReq)
    if err != nil {
        return err
    }

    if tradeData.RequestedValue != 1 {
        return fmt.Errorf("trading is disabled for this account")
    }

    // Check if EA trading is allowed (if using EA)
    eaReq := &pb.AccountInfoIntegerRequest{
        PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_TRADE_EXPERT,
    }
    eaData, err := account.AccountInfoInteger(ctx, eaReq)
    if err != nil {
        return err
    }

    if eaData.RequestedValue != 1 {
        log.Println("WARNING: Expert Advisor trading is disabled")
    }

    return nil
}
```

---

## 📚 See Also

* [AccountSummary](./AccountSummary.md) - Get all account properties in one call (RECOMMENDED)
* [AccountInfoDouble](./AccountInfoDouble.md) - Get specific double account properties
* [AccountInfoString](./AccountInfoString.md) - Get specific string account properties
