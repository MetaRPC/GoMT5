# ✅ Getting an Account Summary

> **Request:** full account summary (`AccountSummaryData`) from **MT5**. Fetch all core account metrics in a single call.

**API Information:**

* **SDK wrapper:** `MT5Account.AccountSummary(...)` (from Go package `github.com/MetaRPC/GoMT5/mt5`)
* **gRPC service:** `mt5_term_api.AccountHelper`
* **Proto definition:** `AccountSummary` (defined in `mt5-term-api-account-helper.proto`)

### RPC

* **Service:** `mt5_term_api.AccountHelper`
* **Method:** `AccountSummary(AccountSummaryRequest) → AccountSummaryReply`
* **Low‑level client (generated):** `AccountHelperClient.AccountSummary(ctx, request, opts...)`
* **SDK wrapper (MT5Account):**


## 💬 Just the essentials

* **What it is.** Single RPC returning account state: balance, equity, currency, leverage, trade mode, server time.
* **Why you need it.** Fast dashboard/CLI status; double‑check login/currency/leverage; heartbeat via `ServerTime`.
* **Sanity check.** If you see `AccountLogin`, `AccountCurrency`, `AccountLeverage`, `AccountEquity` → connection is alive.

---

## 🎯 Purpose

Use it to display real‑time account state and sanity‑check connectivity:

* Dashboard/CLI status in one call.
* Verify balance & equity before trading.
* Terminal heartbeat via `ServerTime` and `UtcTimezoneServerTimeShiftMinutes`.

---

## 📚 Tutorial

For a detailed line-by-line explanation with examples, see:
**→ [AccountSummary - How it works](../HOW_IT_WORK/1.%20Account_information_HOW/AccountSummary_HOW.md)**

---

```go
package mt5

type MT5Account struct {
    // ...
}

// AccountSummary retrieves full account information in a single call.
// Returns AccountSummaryData with Balance, Equity, Leverage, Currency, and other account properties.
func (a *MT5Account) AccountSummary(
    ctx context.Context,
    req *pb.AccountSummaryRequest,
) (*pb.AccountSummaryData, error)
```

**Request message:**

`AccountSummaryRequest {}` *(empty message)*

**Reply message:**

```protobuf
AccountSummaryReply {
  oneof response {
    AccountSummaryData data = 1;
    Error error = 2;
  }
}
```

---

## 🔽 Input

| Parameter | Type                          | Description                                             |
| --------- | ----------------------------- | ------------------------------------------------------- |
| `ctx`     | `context.Context`             | Context for deadline/timeout and cancellation           |
| `req`     | `*pb.AccountSummaryRequest`   | Empty request message (pass `&pb.AccountSummaryRequest{}`) |

**Context options:**

```go
// 1. With timeout (recommended)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 2. With deadline
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
defer cancel()

// 3. With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

---

## ⬆️ Output — `AccountSummaryData`

| Field                                | Type                  | Go Type                    | Description                                      |
| ------------------------------------ | --------------------- | -------------------------- | ------------------------------------------------ |
| `AccountLogin`                       | `int64`               | `int64`                    | Trading account login (ID).                      |
| `AccountBalance`                     | `double`              | `float64`                  | Balance excluding floating P/L.                  |
| `AccountEquity`                      | `double`              | `float64`                  | Equity = balance + floating P/L.                 |
| `AccountUserName`                    | `string`              | `string`                   | Account holder name.                             |
| `AccountLeverage`                    | `int64`               | `int64`                    | Leverage (e.g., `100` for 1:100).                |
| `AccountTradeMode`                   | `MrpcEnumAccountTradeMode` (enum) | `int32` | See `MrpcEnumAccountTradeMode` below.            |
| `AccountCompanyName`                 | `string`              | `string`                   | Broker/company display name.                     |
| `AccountCurrency`                    | `string`              | `string`                   | Deposit currency (e.g., `"USD"`, `"EUR"`).       |
| `ServerTime`                         | `google.protobuf.Timestamp` | `*timestamppb.Timestamp` | Server time (UTC) at response.                   |
| `UtcTimezoneServerTimeShiftMinutes`  | `int64`               | `int64`                    | Server offset relative to UTC (minutes).         |
| `AccountCredit`                      | `double`              | `float64`                  | Credit amount.                                   |

---

## 🧱 Related enums (from proto)

### `MrpcEnumAccountTradeMode`

Defined in `mt5-term-api-account-helper.proto`:

```protobuf
enum MrpcEnumAccountTradeMode {
  MRPC_ACCOUNT_TRADE_MODE_DEMO = 0;     // demo/practice
  MRPC_ACCOUNT_TRADE_MODE_CONTEST = 1;  // contest
  MRPC_ACCOUNT_TRADE_MODE_REAL = 2;     // real trading
}
```

**Go constants (generated):**

```go
const (
    MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_DEMO    = 0
    MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_CONTEST = 1
    MrpcEnumAccountTradeMode_MRPC_ACCOUNT_TRADE_MODE_REAL    = 2
)
```

**Helper function to convert enum to string:**

```go
func AccountTradeModeToString(mode int32) string {
    switch mode {
    case 0:
        return "DEMO"
    case 1:
        return "CONTEST"
    case 2:
        return "REAL"
    default:
        return fmt.Sprintf("UNKNOWN(%d)", mode)
    }
}
```

---


## 🧩 Notes & Tips

* **Automatic reconnection:** All `MT5Account` methods have built-in protection against transient gRPC errors with automatic reconnection via `ExecuteWithReconnect`.
* **Default timeout:** If context has no deadline, a default `10s` timeout is applied automatically.
* **Convert timestamp:** Use `.AsTime()` method on `timestamppb.Timestamp` to convert to `time.Time`.
* **Nil context:** If you pass `nil` context, `context.Background()` is used automatically.

---

## 🔗 Usage Examples

### 1) Basic usage with timeout

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

    // Connect first
    err = account.Connect()
    if err != nil {
        panic(err)
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Get account summary
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        panic(err)
    }

    fmt.Printf("Account Login: %d\n", summary.AccountLogin)
    fmt.Printf("Balance: %.2f %s\n", summary.AccountBalance, summary.AccountCurrency)
    fmt.Printf("Equity: %.2f\n", summary.AccountEquity)
    fmt.Printf("Leverage: 1:%d\n", summary.AccountLeverage)
}
```

### 2) Compact status line for CLI

```go
func PrintAccountStatus(account *mt5.MT5Account) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    s, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }

    status := fmt.Sprintf(
        "Acc %d | %s | Bal %.2f | Eq %.2f | Lev 1:%d | Mode %s",
        s.AccountLogin,
        s.AccountCurrency,
        s.AccountBalance,
        s.AccountEquity,
        s.AccountLeverage,
        AccountTradeModeToString(s.AccountTradeMode),
    )

    fmt.Println(status)
}

// Output: Acc 123456 | USD | Bal 10000.00 | Eq 10050.25 | Lev 1:100 | Mode DEMO
```

### 3) Server time with timezone shift

```go
func PrintServerTime(account *mt5.MT5Account) {
    ctx := context.Background()

    s, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        panic(err)
    }

    // Convert protobuf Timestamp to time.Time
    serverUTC := s.ServerTime.AsTime()

    // Apply timezone shift
    shift := time.Duration(s.UtcTimezoneServerTimeShiftMinutes) * time.Minute
    serverLocal := serverUTC.Add(shift)

    fmt.Printf("Server time (UTC): %s\n", serverUTC.Format(time.RFC3339))
    fmt.Printf("Server time (local): %s\n", serverLocal.Format(time.RFC3339))
    fmt.Printf("Timezone shift: %+d minutes\n", s.UtcTimezoneServerTimeShiftMinutes)
}
```

### 4) Using context cancellation

```go
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

    // Create cancellable context
    ctx, cancel := context.WithCancel(context.Background())

    // Cancel after 2 seconds from another goroutine
    go func() {
        time.Sleep(2 * time.Second)
        cancel()
    }()

    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        if err == context.Canceled {
            fmt.Println("Request was cancelled")
        } else {
            fmt.Printf("Error: %v\n", err)
        }
        return
    }

    fmt.Printf("Balance: %.2f\n", summary.AccountBalance)
}
```

### 5) Structured data model

```go
// AccountSummaryView is a clean Go struct for UI/business logic
type AccountSummaryView struct {
    Login        int64
    Currency     string
    Balance      float64
    Equity       float64
    Leverage     int64
    TradeMode    string
    CompanyName  string
    ServerTime   time.Time
}

func FromProto(data *pb.AccountSummaryData) AccountSummaryView {
    return AccountSummaryView{
        Login:       data.AccountLogin,
        Currency:    data.AccountCurrency,
        Balance:     data.AccountBalance,
        Equity:      data.AccountEquity,
        Leverage:    data.AccountLeverage,
        TradeMode:   AccountTradeModeToString(data.AccountTradeMode),
        CompanyName: data.AccountCompanyName,
        ServerTime:  data.ServerTime.AsTime(),
    }
}

// Usage:
proto, _ := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
view := FromProto(proto)
fmt.Printf("%+v\n", view)
// Output: {Login:123456 Currency:USD Balance:10000.00 Equity:10050.25 Leverage:100 TradeMode:DEMO CompanyName:BrokerXYZ ServerTime:2025-12-27 22:30:15}
```

### 6) Error handling with automatic reconnect

```go
func GetAccountSummaryWithRetry(account *mt5.MT5Account) (*pb.AccountSummaryData, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // MT5Account.AccountSummary automatically handles reconnection
    // via ExecuteWithReconnect wrapper
    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        // Error is already retried by ExecuteWithReconnect
        return nil, fmt.Errorf("failed to get account summary: %w", err)
    }

    return summary, nil
}
```

---

## 🔧 Common Patterns

### Health check

```go
func IsAccountHealthy(account *mt5.MT5Account) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
    if err != nil {
        return false
    }

    // Check if we have valid account data
    return summary.AccountLogin > 0 &&
           summary.AccountCurrency != "" &&
           summary.AccountLeverage > 0
}
```

### Monitor account changes

```go
func MonitorAccount(account *mt5.MT5Account, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
        cancel()

        if err != nil {
            fmt.Printf("❌ Error: %v\n", err)
            continue
        }

        fmt.Printf("[%s] Balance: %.2f | Equity: %.2f\n",
            time.Now().Format("15:04:05"),
            summary.AccountBalance,
            summary.AccountEquity,
        )
    }
}

// Usage:
// MonitorAccount(account, 5*time.Second) // Update every 5 seconds
```

---

## 📚 See Also

* [AccountInfoDouble](./AccountInfoDouble.md) - Get specific double account properties
* [AccountInfoInteger](./AccountInfoInteger.md) - Get specific integer account properties
* [AccountInfoString](./AccountInfoString.md) - Get specific string account properties
