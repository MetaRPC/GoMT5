### Example from file: `examples/demos/lowlevel/01_general_operations.go`

## Overview

> This example demonstrates how to perform a low-level gRPC call to the MetaTrader gateway to retrieve one specific numeric account property (in this case — Balance). The AccountInfoDouble() method is used for targeted querying of float64 (floating-point number) data, without the need to request the entire AccountSummary().
>
> It is designed strictly for retrieving a single numeric account parameter, which is stored as float64 (floating-point number).
>
> That is — it performs no arrays, structures, strings, or actions.
>
> It simply extracts one field from the MetaTrader AccountInfo structure and returns it as a number.


---

### 🧩 Code Example

```go
fmt.Println("\n3.2. AccountInfoDouble() - Get specific double property (example: Balance)")

balanceReq := &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
}
balanceData, err := account.AccountInfoDouble(ctx, balanceReq)
if err != nil {
    helpers.PrintShortError(err, "AccountInfoDouble(BALANCE) failed")
} else {
    fmt.Printf("  Balance:                       %.2f\n", balanceData.GetRequestedValue())
}
```

---

### 🟢 Detailed Code Explanation

```go
fmt.Println("\n3.2. AccountInfoDouble() - Get specific double property (example: Balance)")
```

Prints the section header to the console.
In this case, it shows an example call to retrieve the account **balance**.

---

```go
balanceReq := &pb.AccountInfoDoubleRequest{
    PropertyId: pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE,
}
```

Creates a request for the gRPC method.
The `PropertyId` parameter specifies **which account property** to retrieve.
Here it's `ACCOUNT_BALANCE` — the account balance.

📘 Similarly, you can request:

* `ACCOUNT_EQUITY` — equity,
* `ACCOUNT_MARGIN` — margin,
* `ACCOUNT_PROFIT` — current profit,
* `ACCOUNT_CREDIT` — credit, etc.

---

```go
balanceData, err := account.AccountInfoDouble(ctx, balanceReq)
```

Executes a gRPC call to the MetaTrader server through the gateway.

* `ctx` — context (usually with a timeout).
* `balanceReq` — request structure.
  Returns:
* `balanceData` — result (protobuf structure),
* `err` — error (if the request failed).

---

```go
if err != nil {
    helpers.PrintShortError(err, "AccountInfoDouble(BALANCE) failed")
} else {
    fmt.Printf("  Balance:                       %.2f\n", balanceData.GetRequestedValue())
}
```

Checks for errors.
If successful, prints the value received from the server.
The `.GetRequestedValue()` method returns **float64** — the actual numeric value of the property.

---

### Available Properties Table (`AccountInfoDoublePropertyType`)

| Constant                     | Description                  | Value Type | Use Case                              |
| ---------------------------- | ---------------------------- | ---------- | ------------------------------------- |
| `ACCOUNT_BALANCE`            | Account balance              | float64    | Check current balance                 |
| `ACCOUNT_CREDIT`             | Credit                       | float64    | Track credit funds                    |
| `ACCOUNT_PROFIT`             | Current profit/loss          | float64    | PnL monitoring                        |
| `ACCOUNT_EQUITY`             | Account equity               | float64    | Calculate net funds                   |
| `ACCOUNT_MARGIN`             | Used margin                  | float64    | Control account utilization           |
| `ACCOUNT_MARGIN_FREE`        | Free margin                  | float64    | Check if can open position            |
| `ACCOUNT_MARGIN_LEVEL`       | Margin level (%)             | float64    | Risk monitoring                       |
| `ACCOUNT_MARGIN_SO_CALL`     | Margin Call level            | float64    | Warning threshold                     |
| `ACCOUNT_MARGIN_SO_SO`       | Stop Out level               | float64    | Forced closure threshold              |
| `ACCOUNT_MARGIN_INITIAL`     | Initial margin               | float64    | Reserved funds                        |
| `ACCOUNT_MARGIN_MAINTENANCE` | Maintenance margin           | float64    | Control minimum equity                |
| `ACCOUNT_ASSETS`             | Assets                       | float64    | Reporting                             |
| `ACCOUNT_LIABILITIES`        | Liabilities                  | float64    | Financial analysis                    |
| `ACCOUNT_COMMISSION_BLOCKED` | Blocked commissions          | float64    | Track available funds                 |

---

### Wrapper Example for Property Requests

```go
func GetAccountDoubleProp(ctx context.Context, account AccountServiceClient, prop pb.AccountInfoDoublePropertyType) (float64, error) {
    req := &pb.AccountInfoDoubleRequest{PropertyId: prop}
    resp, err := account.AccountInfoDouble(ctx, req)
    if err != nil {
        return 0, err
    }
    return resp.GetRequestedValue(), nil
}
```

💡 **Why use a wrapper:**
It allows you to **shorten and simplify code** when you need to frequently call `AccountInfoDouble()` with different parameters.
Now you can simply call a short function:

```go
balance, _ := GetAccountDoubleProp(ctx, account, pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE)
equity, _  := GetAccountDoubleProp(ctx, account, pb.AccountInfoDoublePropertyType_ACCOUNT_EQUITY)
```

Such a wrapper makes code **cleaner, more compact, and easier to reuse**.

---

### Summary

* The `AccountInfoDouble()` method is used for **quickly retrieving a single double field** from the account structure.
* This is a **low-level getter** that returns a pure value, without extra data.
* Applied in:

  * **HFT gateways** to minimize load,
  * **bots and advisors** for quick balance or equity checks,
  * **risk systems** for margin and PnL control.
