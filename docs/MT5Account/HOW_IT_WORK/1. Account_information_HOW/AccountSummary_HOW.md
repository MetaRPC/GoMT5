### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> This example demonstrates how a low-level gRPC call to a MetaTrader server is executed through a gateway.  
> The code shows the process of retrieving summary information about a trading account using the AccountSummary() method.  
> Below is a detailed breakdown of each line and an explanation of how such a call is constructed at the client-server interaction level.

### 🧩 Code Example

```go
fmt.Println("\n3.1. AccountSummary() - Get all account data in one call")

	summaryReq := &pb.AccountSummaryRequest{}
	// Use timeout context to prevent hanging
	summaryCtx, summaryCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer summaryCancel()

	summaryData, err := account.AccountSummary(summaryCtx, summaryReq)
	if err != nil {
		return fmt.Errorf("AccountSummary failed: %w", err)
	}

	fmt.Println("\nAccount Summary (direct protobuf field access):")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("  Login:               %d\n", summaryData.AccountLogin)
	fmt.Printf("  UserName:            %s\n", summaryData.AccountUserName)
	fmt.Printf("  Company:             %s\n", summaryData.AccountCompanyName)
	fmt.Printf("  Currency:            %s\n", summaryData.AccountCurrency)
	fmt.Printf("  Balance:             %.2f\n", summaryData.AccountBalance)
	fmt.Printf("  Equity:              %.2f\n", summaryData.AccountEquity)
	fmt.Printf("  Credit:              %.2f\n", summaryData.AccountCredit)
	fmt.Printf("  Leverage:            1:%d\n", summaryData.AccountLeverage)
	fmt.Printf("  Trade Mode:          %v\n", summaryData.AccountTradeMode)

	// ServerTime is a protobuf Timestamp in the broker server's time zone; ServerTimeUtc is the same moment in UTC
	if summaryData.ServerTime != nil {
		serverTime := summaryData.ServerTime.AsTime()
		fmt.Printf("  Server Time:         %s (broker server time)\n", serverTime.Format("2006-01-02 15:04:05"))
	}
	if summaryData.ServerTimeUtc != nil {
		fmt.Printf("  Server Time (UTC):   %s\n", summaryData.ServerTimeUtc.AsTime().Format("2006-01-02 15:04:05"))
	}

// UTC Timezone Shift: server time offset from UTC in minutes
// For example: 120 minutes = UTC+2 (the server is 2 hours ahead of UTC)
	fmt.Printf("  UTC Timezone Shift:  %d minutes (UTC%+.1f)\n",
		summaryData.UtcTimezoneServerTimeShiftMinutes,
		float64(summaryData.UtcTimezoneServerTimeShiftMinutes)/60.0)
```

### 🟢 Detailed Code Explanation for AccountSummary() (Go)

```go
fmt.Println("\n3.1. AccountSummary() - Get all account data in one call")
```

Prints a header to the console. `\n` — a newline character. Used to visually separate the output.

---

```go
summaryReq := &pb.AccountSummaryRequest{}
```

Creates an empty `AccountSummaryRequest` request. `pb` is a package generated from a `.proto` file (gRPC). The request is passed to `AccountSummary()` to retrieve the account summary.

---

```go
// Use timeout context to prevent hanging
```

Comment: we create a context with a timeout so the program doesn't hang if the server doesn't respond.

---

```go
summaryCtx, summaryCancel := context.WithTimeout(context.Background(), 10*time.Second)
```

Creates a context with an execution time limit — 10 seconds. If the response doesn't arrive in time, the request will be interrupted.

* `context.Background()` — base context.
* `context.WithTimeout()` — adds a time limit.
* `summaryCancel` — function to cancel the context.

---

```go
defer summaryCancel()
```

Defers the call to `summaryCancel()` until the end of the function. This ensures resource cleanup.

---

```go
summaryData, err := account.AccountSummary(summaryCtx, summaryReq)
```

Calls the `AccountSummary()` method on the `account` object.

* `summaryCtx` — context with timeout.
* `summaryReq` — request.
  Returns:
* `summaryData` — structure with the result (account information).
* `err` — error if something went wrong.

---

```go
if err != nil {
	return fmt.Errorf("AccountSummary failed: %w", err)
}
```

Checks if an error occurred. If `err != nil`, the demo function returns the error. It does not exit on the spot
(`helpers.Fatal()` would call `os.Exit`, which skips the deferred `account.DisconnectAndClose()` and leaves the
terminal connected): returning lets the deferred cleanup disconnect the terminal first.

---

```go
fmt.Println("\nAccount Summary (direct protobuf field access):")
fmt.Println("─────────────────────────────────────────────────")
```

Prints a header and separator before the data to make the output readable.

---

```go
fmt.Printf("  Login:               %d\n", summaryData.AccountLogin)
```

Prints the account login. `%d` — format for integers.

---

```go
fmt.Printf("  UserName:            %s\n", summaryData.AccountUserName)
fmt.Printf("  Company:             %s\n", summaryData.AccountCompanyName)
fmt.Printf("  Currency:            %s\n", summaryData.AccountCurrency)
```

Prints string fields (username, company, account currency). `%s` — format for strings.

---

```go
fmt.Printf("  Balance:             %.2f\n", summaryData.AccountBalance)
fmt.Printf("  Equity:              %.2f\n", summaryData.AccountEquity)
fmt.Printf("  Credit:              %.2f\n", summaryData.AccountCredit)
```

Prints numeric data (balance, equity, credit). `%.2f` — number with two decimal places.

---

```go
fmt.Printf("  Leverage:            1:%d\n", summaryData.AccountLeverage)
```

Prints the account leverage in the format `1:500`.

---

```go
fmt.Printf("  Trade Mode:          %v\n", summaryData.AccountTradeMode)
```

Prints the trading mode (for example, real or demo account). `%v` — universal format.

---

### ⏱️ Server Time Output Block

> In Go, the time format is not specified using symbols like YYYY or HH, but rather by **an example of a specific date** — `2006-01-02 15:04:05`. These numbers are not chosen randomly: they represent the order of date and time elements, where:
>
> * 2006 → year,
> * 01 → month,
> * 02 → day,
> * 15 → hour (in 24-hour format),
> * 04 → minutes,
> * 05 → seconds.
>
> Therefore, to output date and time in the desired format, you always specify exactly this template with actual numbers.

```go
// ServerTime is a protobuf Timestamp in the broker server's time zone; ServerTimeUtc is the same moment in UTC
if summaryData.ServerTime != nil {
    serverTime := summaryData.ServerTime.AsTime()
    fmt.Printf("  Server Time:         %s (broker server time)\n", serverTime.Format("2006-01-02 15:04:05"))
}
if summaryData.ServerTimeUtc != nil {
    fmt.Printf("  Server Time (UTC):   %s\n", summaryData.ServerTimeUtc.AsTime().Format("2006-01-02 15:04:05"))
}

// UTC Timezone Shift: server time offset from UTC in minutes
// For example: 120 minutes = UTC+2 (the server is 2 hours ahead of UTC)
fmt.Printf("  UTC Timezone Shift:  %d minutes (UTC%+.1f)\n",
    summaryData.UtcTimezoneServerTimeShiftMinutes,
    float64(summaryData.UtcTimezoneServerTimeShiftMinutes)/60.0)
```

**What this block does:**

1. Checks that the time from the server (`ServerTime`) is present.
2. Converts protobuf time to the `time.Time` type via the `.AsTime()` method.
3. Formats the date and time in a human-readable form: `YYYY-MM-DD HH:MM:SS`. `ServerTime` shows the broker server's
   clock (its time zone); `ServerTimeUtc` is the same moment in UTC.
4. Prints the server time offset from UTC both in minutes and hours (for example, `120 minutes (UTC+2.0)`;
   positive = the server is ahead of UTC).

**Why this is important:**

* In trading systems, it's crucial to know the exact server time for synchronizing quotes, trades, and reports.
* The MetaTrader server may be located in a different timezone, and this affects candlestick calculations, order opening times, and reporting.

---

### Summary:

1. A request and context with timeout are created.
2. The `AccountSummary()` method is called via RPC.
3. Errors are checked.
4. If successful — account data is printed.
5. At the bottom, **server time and UTC offset** are additionally displayed.

### Applications:

* Testing the MetaTrader gateway API.
* Retrieving real-time account information.
* Use within algorithmic systems for account state monitoring and time synchronization.
