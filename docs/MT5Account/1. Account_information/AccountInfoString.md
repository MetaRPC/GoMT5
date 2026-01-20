# 📝 Getting Account String Properties

> **Request:** specific string-type account property from **MT5** terminal using property identifier.

**API Information:**

* **SDK wrapper:** `MT5Account.AccountInfoString(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.AccountInformation`
* **Proto definition:** `AccountInfoString` (defined in `mt5-term-api-account-information.proto`)

### RPC

* **Service:** `mt5_term_api.AccountInformation`
* **Method:** `AccountInfoString(AccountInfoStringRequest) → AccountInfoStringReply`
* **Low‑level client (generated):** `AccountInformationClient.AccountInfoString(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**

## 💬 Just the essentials

* **What it is.** Retrieve a specific string-type account property by ID (name, server, currency, broker).
* **Why you need it.** Get account text information like client name, broker, currency code.
* **When to use.** Query individual text properties. For all properties, use `AccountSummary()`.

---

## 🎯 Purpose

Use it to query specific string account properties:

* Get client/account name
* Check trade server name
* Verify account currency
* Get broker company name

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [AccountInfoString - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountInfoString_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// AccountInfoString retrieves a string-type account property.
// Returns AccountInfoStringData with the requested string value.
func (a *MT5Account) AccountInfoString(
    ctx context.Context,
    req *pb.AccountInfoStringRequest,
) (*pb.AccountInfoStringData, error)
```

**Request message:**

```protobuf
message AccountInfoStringRequest {
  AccountInfoStringPropertyType property_id = 1;
}
```

**Reply message:**

```protobuf
message AccountInfoStringReply {
  oneof response {
    AccountInfoStringData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter     | Type                              | Description                                             |
| ------------- | --------------------------------- | ------------------------------------------------------- |
| `ctx`         | `context.Context`                 | Context for deadline/timeout and cancellation           |
| `req`         | `*pb.AccountInfoStringRequest`    | Request with property_id specifying which property to retrieve |
| `req.property_id` | `AccountInfoStringPropertyType` (enum) | Property identifier (see enum below) |

---

## ⬆️ Output — `AccountInfoStringData`

| Field             | Type     | Go Type  | Description                                      |
| ----------------- | -------- | -------- | ------------------------------------------------ |
| `requested_value` | `string` | `string` | The value of the requested account property      |

---

## 🧱 Related enums (from proto)

### `AccountInfoStringPropertyType`

Defined in `mt5-term-api-account-information.proto`:

```protobuf
enum AccountInfoStringPropertyType {
  ACCOUNT_NAME = 0;     // Client name
  ACCOUNT_SERVER = 1;   // Trade server name
  ACCOUNT_CURRENCY = 2; // Account currency (e.g., "USD", "EUR")
  ACCOUNT_COMPANY = 3;  // Name of company that serves the account (broker)
}
```

**Go constants (generated):**

```go
const (
    AccountInfoStringPropertyType_ACCOUNT_NAME     = 0
    AccountInfoStringPropertyType_ACCOUNT_SERVER   = 1
    AccountInfoStringPropertyType_ACCOUNT_CURRENCY = 2
    AccountInfoStringPropertyType_ACCOUNT_COMPANY  = 3
)
```

---


## 🧩 Notes & Tips

* **Automatic reconnection:** Built-in protection against transient gRPC errors.
* **Default timeout:** Default `3s` timeout applied if context has no deadline.
* **Currency codes:** Returns standard 3-letter currency codes (USD, EUR, GBP, etc.).
* **Prefer AccountSummary:** For multiple properties, use `AccountSummary()` to avoid multiple round-trips.

---

## 🔗 Usage Examples

### 1) Get account currency

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "github.com/MetaRPC/GoMT5/package"
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

    req := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY,
    }

    data, err := account.AccountInfoString(ctx, req)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Account Currency: %s\n", data.RequestedValue)
    // Output: Account Currency: USD
}
```

### 2) Get client name

```go
func GetClientName(account *mt5.MT5Account) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_NAME,
    }

    data, err := account.AccountInfoString(ctx, req)
    if err != nil {
        return "", fmt.Errorf("failed to get client name: %w", err)
    }

    return data.RequestedValue, nil
}

// Usage:
name, _ := GetClientName(account)
fmt.Printf("Account Holder: %s\n", name)
```

### 3) Get trade server name

```go
func GetServerName(account *mt5.MT5Account) string {
    ctx := context.Background()

    req := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_SERVER,
    }

    data, err := account.AccountInfoString(ctx, req)
    if err != nil {
        return "unknown"
    }

    return data.RequestedValue
}

// Usage:
server := GetServerName(account)
fmt.Printf("Connected to: %s\n", server)
// Output: Connected to: MetaQuotes-Demo
```

### 4) Get broker company name

```go
func GetBrokerName(account *mt5.MT5Account) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY,
    }

    data, err := account.AccountInfoString(ctx, req)
    if err != nil {
        return "", err
    }

    return data.RequestedValue, nil
}

// Usage:
broker, err := GetBrokerName(account)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Broker: %s\n", broker)
```

### 5) Display account info banner

```go
func PrintAccountBanner(account *mt5.MT5Account) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Get client name
    nameReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_NAME,
    }
    nameData, _ := account.AccountInfoString(ctx, nameReq)

    // Get broker
    brokerReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY,
    }
    brokerData, _ := account.AccountInfoString(ctx, brokerReq)

    // Get server
    serverReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_SERVER,
    }
    serverData, _ := account.AccountInfoString(ctx, serverReq)

    // Get currency
    currReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY,
    }
    currData, _ := account.AccountInfoString(ctx, currReq)

    // Print banner
    fmt.Println("=" + strings.Repeat("=", 50))
    fmt.Printf("  Account: %s\n", nameData.RequestedValue)
    fmt.Printf("  Broker:  %s\n", brokerData.RequestedValue)
    fmt.Printf("  Server:  %s\n", serverData.RequestedValue)
    fmt.Printf("  Currency: %s\n", currData.RequestedValue)
    fmt.Println("=" + strings.Repeat("=", 50))
}

// Usage:
PrintAccountBanner(account)
/* Output:
===================================================
  Account: John Doe
  Broker:  XM Global Limited
  Server:  XM-Real 34
  Currency: USD
===================================================
*/
```

### 6) Verify account currency before trading

```go
func VerifyCurrency(account *mt5.MT5Account, expectedCurrency string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY,
    }

    data, err := account.AccountInfoString(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to get currency: %w", err)
    }

    currency := data.RequestedValue

    if currency != expectedCurrency {
        return fmt.Errorf("currency mismatch: expected %s, got %s",
            expectedCurrency, currency)
    }

    fmt.Printf("Currency verified: %s\n", currency)
    return nil
}

// Usage:
err := VerifyCurrency(account, "USD")
if err != nil {
    log.Fatal(err)
}
```

---

## 🔧 Common Patterns

### Get all string properties

```go
type AccountStringInfo struct {
    Name     string
    Server   string
    Currency string
    Company  string
}

func GetAccountStringInfo(account *mt5.MT5Account) (*AccountStringInfo, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    info := &AccountStringInfo{}

    // Get name
    nameReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_NAME,
    }
    nameData, err := account.AccountInfoString(ctx, nameReq)
    if err != nil {
        return nil, err
    }
    info.Name = nameData.RequestedValue

    // Get server
    serverReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_SERVER,
    }
    serverData, err := account.AccountInfoString(ctx, serverReq)
    if err != nil {
        return nil, err
    }
    info.Server = serverData.RequestedValue

    // Get currency
    currReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_CURRENCY,
    }
    currData, err := account.AccountInfoString(ctx, currReq)
    if err != nil {
        return nil, err
    }
    info.Currency = currData.RequestedValue

    // Get company
    compReq := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY,
    }
    compData, err := account.AccountInfoString(ctx, compReq)
    if err != nil {
        return nil, err
    }
    info.Company = compData.RequestedValue

    return info, nil
}

// Usage:
info, err := GetAccountStringInfo(account)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%+v\n", info)
// Output: {Name:John Doe Server:XM-Real 34 Currency:USD Company:XM Global Limited}
```

### Connection validation

```go
func ValidateConnection(account *mt5.MT5Account, expectedServer string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req := &pb.AccountInfoStringRequest{
        PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_SERVER,
    }

    data, err := account.AccountInfoString(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to verify server: %w", err)
    }

    server := data.RequestedValue

    if server != expectedServer {
        return fmt.Errorf("connected to wrong server: expected %s, got %s",
            expectedServer, server)
    }

    log.Printf("Successfully connected to %s", server)
    return nil
}
```

---

## 📚 See Also

* [AccountSummary](./AccountSummary.md) - Get all account properties in one call (RECOMMENDED)
* [AccountInfoDouble](./AccountInfoDouble.md) - Get specific double account properties
* [AccountInfoInteger](./AccountInfoInteger.md) - Get specific integer account properties
