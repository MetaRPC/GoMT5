### Example from file: `examples/demos/lowlevel/01_general_operations.go`

## Overview

> This example demonstrates the use of the AccountInfoInteger() method — a low-level gRPC call to retrieve one specific integer parameter of a MetaTrader account (for example, Leverage). The method is used when you need to request integer fields of an account: login, leverage, account type, trading mode, etc.
>
> It is designed strictly for retrieving a single integer-type parameter. It returns no structures, arrays, or strings — it simply extracts one field from AccountInfo and returns it as an integer.


---

### 🧩 Code Example

```go
fmt.Println("\n3.3. AccountInfoInteger() - Get specific integer property (example: Leverage)")

leverageReq := &pb.AccountInfoIntegerRequest{
    PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE,
}
leverageData, err := account.AccountInfoInteger(ctx, leverageReq)
if err != nil {
    helpers.PrintShortError(err, "AccountInfoInteger(LEVERAGE) failed")
} else {
    fmt.Printf("  Leverage: %d\n", leverageData.GetRequestedValue())
}
```

---

### 🟢 Detailed Code Explanation

```go
fmt.Println("\n3.3. AccountInfoInteger() - Get specific integer property (example: Leverage)")
```

Prints the section header. Demonstrates retrieving a specific integer property — `Leverage`.

---

```go
leverageReq := &pb.AccountInfoIntegerRequest{
    PropertyId: pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE,
}
```

Creates an `AccountInfoIntegerRequest` structure.
The `PropertyId` parameter specifies **which property** to request. In this case — `ACCOUNT_LEVERAGE` (leverage).

---

```go
leverageData, err := account.AccountInfoInteger(ctx, leverageReq)
```

Executes a gRPC call to the MetaTrader server.

* `ctx` — request context (with timeout).
* `leverageReq` — request structure.
  Returns the result (`leverageData`) and error (`err`).

---

```go
if err != nil {
    helpers.PrintShortError(err, "AccountInfoInteger(LEVERAGE) failed")
} else {
    fmt.Printf("  Leverage: %d\n", leverageData.GetRequestedValue())
}
```

Checks for errors. If successful, extracts the result via `.GetRequestedValue()` and prints it to the console.
The value is returned as an **integer** (`int32`).

---

### Available Properties Table (`AccountInfoIntegerPropertyType`)

| Constant                  | Description                                                          | Type  | Use Case                                         |
| ------------------------- | -------------------------------------------------------------------- | ----- | ------------------------------------------------ |
| `ACCOUNT_LOGIN`           | Account login (account number)                                       | int64 | Authorization and account identification         |
| `ACCOUNT_TRADE_MODE`      | Trading mode (`0=DEMO`, `1=CONTEST`, `2=REAL`)                       | int32 | Determine account type                           |
| `ACCOUNT_LEVERAGE`        | Leverage (e.g., 100 = 1:100)                                         | int32 | Margin calculation and risk management           |
| `ACCOUNT_LIMIT_ORDERS`    | Order limit (`0` = no restrictions)                                  | int32 | Control allowed number of positions              |
| `ACCOUNT_MARGIN_SO_MODE`  | Stop Out calculation type (`0=percent`, `1=money`)                   | int32 | Determine Margin Call/Stop Out conditions        |
| `ACCOUNT_TRADE_ALLOWED`   | Is trading allowed (`1/0`)                                           | int32 | Check ability to send orders                     |
| `ACCOUNT_TRADE_EXPERT`    | Is Expert Advisor trading allowed (`1/0`)                            | int32 | Important for robots and gateways                |
| `ACCOUNT_MARGIN_MODE`     | Margin calculation mode (`0=Retail Netting`, `1=Exchange`, `2=Hedging`) | int32 | Defines position accounting in MT5               |
| `ACCOUNT_CURRENCY_DIGITS` | Number of decimal places in currency                                 | int32 | Formatting and reporting                         |
| `ACCOUNT_FIFO_CLOSE`      | Forced FIFO close (`1/0`)                                            | int32 | Important for brokers with FIFO restrictions     |
| `ACCOUNT_HEDGE_ALLOWED`   | Is hedging allowed (`1/0`)                                           | int32 | Check ability to open opposite positions         |

---

### Wrapper Example for Repeated Calls

```go
func GetAccountIntegerProp(ctx context.Context, account *mt5.MT5Account, prop pb.AccountInfoIntegerPropertyType) (int64, error) {
    req := &pb.AccountInfoIntegerRequest{PropertyId: prop}
    resp, err := account.AccountInfoInteger(ctx, req)
    if err != nil {
        return 0, err
    }
    return resp.GetRequestedValue(), nil
}
```

💡 **Why use a wrapper:**
It simplifies code and makes calls shorter.
Now you can quickly request different integer properties:

```go
login, _ := GetAccountIntegerProp(ctx, account, pb.AccountInfoIntegerPropertyType_ACCOUNT_LOGIN)
leverage, _ := GetAccountIntegerProp(ctx, account, pb.AccountInfoIntegerPropertyType_ACCOUNT_LEVERAGE)
```

---

### Summary

* The `AccountInfoInteger()` method is used to retrieve **integer parameters** of an account.
* This is a **low-level getter** that returns a specific value, without extra data.
* Applied for checking account settings, leverage, trading permissions, and system properties.

---

### Practical Applications

* **Gateway diagnostics and connection check** — retrieve login and trading mode.
* **Algorithmic strategies** — check `TradeAllowed` and `Leverage` before trading operations.
* **Risk monitoring** — control leverage, margin modes, and FIFO restrictions.
* **Reporting and configuration control** — collect system parameters for visual panels or APIs.
