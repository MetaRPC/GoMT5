### Example from file: `examples/demos/lowlevel/01_general_operations.go`

## Overview

> This example demonstrates the use of the **`AccountInfoString()`** method — a low-level gRPC call for retrieving **a single specific string property** of a MetaTrader account (for example, `Company`).
> The method is designed to work with string fields such as client name, brokerage company, account currency, or server.
>
> It is used to retrieve **a single text parameter of the account** without loading the entire summary information.
> No arrays, structures, or calculations — the method simply returns a **string** from `AccountInfo`.


---

### 🧩 Code Example

```go
fmt.Println("\n3.4. AccountInfoString() - Get specific string property (example: Company)")

companyReq := &pb.AccountInfoStringRequest{
    PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY,
}
companyData, err := account.AccountInfoString(ctx, companyReq)
if err != nil {
    helpers.PrintShortError(err, "AccountInfoString(COMPANY) failed")
} else {
    fmt.Printf("  Company:                  %s\n", companyData.GetRequestedValue())
}
```

---

### 🟢 Detailed Code Explanation

```go
fmt.Println("\n3.4. AccountInfoString() - Get specific string property (example: Company)")
```

Prints a section header indicating an example of retrieving a string property — in this case `Company` (the broker company name).

---

```go
companyReq := &pb.AccountInfoStringRequest{
    PropertyId: pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY,
}
```

Creates an `AccountInfoStringRequest` structure, specifying the `PropertyId` — the identifier of the property to retrieve.
Here it's `ACCOUNT_COMPANY`, meaning a request for the broker company name.

---

```go
companyData, err := account.AccountInfoString(ctx, companyReq)
```

Executes a gRPC call to the MetaTrader server through the gateway.

* `ctx` — request context (with timeout).
* `companyReq` — request structure.
  Returns the result (`companyData`) and error (`err`).

---

```go
if err != nil {
    helpers.PrintShortError(err, "AccountInfoString(COMPANY) failed")
} else {
    fmt.Printf("  Company:                  %s\n", companyData.GetRequestedValue())
}
```

Checks for errors. If successful, the `.GetRequestedValue()` method returns a **string** with the property value, which is then printed to the console.

Example output:

```
Company:                  MetaQuotes Software Corp.
```

---

### Table of Available Properties (`AccountInfoStringPropertyType`)

| Constant           | Description               | Example Value                 | Use Case                                                     |
| ------------------ | ------------------------- | ----------------------------- | ------------------------------------------------------------ |
| `ACCOUNT_NAME`     | Account owner's name      | `"John Doe"`                  | Displaying client name in UI, logs, reports                  |
| `ACCOUNT_SERVER`   | Trading server name       | `"MetaQuotes-Demo"`           | Verifying connection to the correct server                   |
| `ACCOUNT_CURRENCY` | Account currency          | `"USD"`, `"EUR"`              | Used for calculations and displaying financial data          |
| `ACCOUNT_COMPANY`  | Broker company name       | `"MetaQuotes Software Corp."` | For reporting and broker identification                      |

---

### Wrapper Example for Simplified Calls

```go
func GetAccountStringProp(ctx context.Context, account AccountServiceClient, prop pb.AccountInfoStringPropertyType) (string, error) {
    req := &pb.AccountInfoStringRequest{PropertyId: prop}
    resp, err := account.AccountInfoString(ctx, req)
    if err != nil {
        return "", err
    }
    return resp.GetRequestedValue(), nil
}
```

💡 **Why use a wrapper:**
It makes code compact and reusable.
Now you can quickly retrieve different string properties:

```go
name, _ := GetAccountStringProp(ctx, account, pb.AccountInfoStringPropertyType_ACCOUNT_NAME)
server, _ := GetAccountStringProp(ctx, account, pb.AccountInfoStringPropertyType_ACCOUNT_SERVER)
company, _ := GetAccountStringProp(ctx, account, pb.AccountInfoStringPropertyType_ACCOUNT_COMPANY)
```

---

### Summary

* The `AccountInfoString()` method is used to retrieve **string account properties**.
* This is a **low-level text getter** that returns a specific text field from `AccountInfo`.
* Applied in interfaces, reporting systems, and connection diagnostics.

---

### Practical Applications

* **Monitoring and display** — client name, server, broker, currency.
* **Reporting** — adding text parameters to account reports.
* **Diagnostics** — checking current broker and server.
* **Verification** — matching account name and currency with strategy settings or database.
