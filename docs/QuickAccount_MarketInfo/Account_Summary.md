# Getting an Account Summary (MT5)

> **Request:** full account summary (`AccountSummaryData`) from MT5.
> Fetch all core account metrics in a single call.

---

### Code Example

```go
// High-level (prints formatted summary):
svc.ShowAccountSummary(ctx)

// Low-level (returns data for your own formatting):
summary, err := svc.account.AccountSummary(ctx)
if err != nil {
    log.Printf("❌ AccountSummary error: %v", err)
    return
}
fmt.Printf("Account Summary: Balance=%.2f, Equity=%.2f, Currency=%s\n",
    summary.GetAccountBalance(),
    summary.GetAccountEquity(),
    summary.GetAccountCurrency())
```

---

### Method Signature

```go
func (s *MT5Service) ShowAccountSummary(ctx context.Context)
```

---

## 🔽 Input

No required input parameters.

| Parameter | Type              | Description                               |
| --------- | ----------------- | ----------------------------------------- |
| `ctx`     | `context.Context` | Controls timeout/cancellation for the RPC |

---

## ⬆️ Output

`ShowAccountSummary` prints selected fields from `AccountSummaryData` to stdout.
If you use the low-level call, you get a struct with at least:

| Field               | Type     | Description                               |
| ------------------- | -------- | ----------------------------------------- |
| `AccountBalance`    | `double` | Balance excluding floating P/L            |
| `AccountEquity`     | `double` | Equity = balance + floating P/L           |
| `AccountMargin`     | `double` | Currently used margin                     |
| `AccountFreeMargin` | `double` | Free margin available for new trades      |
| `AccountCurrency`   | `string` | Deposit currency (e.g., `"USD"`, `"EUR"`) |
| `AccountLeverage`   | `int`    | Account leverage (e.g., 100 for 1:100)    |
| `AccountName`       | `string` | Account holder/broker display name        |
| `AccountNumber`     | `int`    | Trading account login ID                  |
| `Company`           | `string` | Broker/company name                       |

> Exact field names may differ slightly by proto version; use getters provided by your `pb` package.

---

## 🎯 Purpose

Use this to display real-time account state and to sanity-check connectivity. Typical cases:

* Show dashboard/CLI status
* Verify free margin and equity before placing trades
* Monitor account health/exposure in bots and diagnostics

---

## 🧩 Notes & Tips

* Call after the terminal is connected and “alive”; otherwise you’ll get an error or zeroed values.
* Wrap in a short per-call timeout (e.g., 3–5s) if running inside longer workflows.
* Values reflect the terminal’s current state; if symbols/positions are loading, re-check after a brief delay.
