### Example from file: `examples/demos/lowlevel/03_streaming_methods.go`

> The **`OnTradeTransaction()`** method is the most detailed streaming method, which transmits **all trading events** occurring on the account in MetaTrader.
> It captures every step: from order opening to its execution and position closing.


---

## 🧩 Code example

```go
transactionReq := &pb.OnTradeTransactionRequest{}

transactionChan, transactionErrChan := account.OnTradeTransaction(ctx, transactionReq)

fmt.Printf("Streaming trade transactions (max %d events or %d seconds)...\n", MAX_EVENTS, MAX_SECONDS)
fmt.Println("  NOTE: This stream sends events for all trade-related transactions")

eventCount = 0
timeout = time.After(MAX_SECONDS * time.Second)

streamTransaction:
for {
	select {
	case transactionData, ok := <-transactionChan:
		if !ok {
			fmt.Println("  Stream closed by server")
			break streamTransaction
		}
		eventCount++
		fmt.Printf("  Event #%d: Type=%v\n", eventCount, transactionData.Type)

		if transactionData.TradeTransaction != nil {
			tx := transactionData.TradeTransaction
			fmt.Printf("    Transaction Type=%v Order=%d Deal=%d\n",
				tx.Type,
				tx.OrderTicket,
				tx.DealTicket)

			if tx.Symbol != "" {
				fmt.Printf("    Symbol=%s Price=%.5f Volume=%.2f\n",
					tx.Symbol,
					tx.Price,
					tx.Volume)
			}
		}

		if eventCount >= MAX_EVENTS {
			fmt.Printf("  ✓ Received %d events, stopping stream\n", MAX_EVENTS)
			break streamTransaction
		}

	case err := <-transactionErrChan:
		if err != nil {
			helpers.PrintShortError(err, "Stream error")
			break streamTransaction
		}

	case <-timeout:
		fmt.Printf("  ⏱ Timeout after %d seconds (received %d events)\n", MAX_SECONDS, eventCount)
		break streamTransaction
	}
}
```

---

## 🟢 Detailed Code Explanation

### 1️. Prepare Request

```go
transactionReq := &pb.OnTradeTransactionRequest{}
```

Create an empty request — subscription applies to **all trading events** of the current account.

---

### 2️. Start the Stream

```go
transactionChan, transactionErrChan := account.OnTradeTransaction(ctx, transactionReq)
```

The method returns two channels:

* **transactionChan** — stream with events;
* **transactionErrChan** — connection errors.

> The stream activates on any changes in orders, deals, and positions.

---

### 3️. Main Processing Loop

```go
case transactionData, ok := <-transactionChan:
```

The loop receives each new event.
Each message contains:

* type (`transactionData.Type`),
* `TradeTransaction` structure with details.

---

### 4️. Process Trade Event

```go
if transactionData.TradeTransaction != nil {
	tx := transactionData.TradeTransaction
	fmt.Printf("    Transaction Type=%v Order=%d Deal=%d\n", ...)
```

Useful fields:

| Field          | Description                                              |
| ------------- | -------------------------------------------------------- |
| `Type`        | Event type (DEAL_ADD, ORDER_ADD, ORDER_DELETE, etc.)     |
| `OrderTicket` | Order number                                             |
| `DealTicket`  | Deal number                                              |
| `Symbol`      | Trading instrument                                       |
| `Price`       | Deal price                                               |
| `Volume`      | Deal volume in lots                                      |

---

### 5️. Display Detailed Information

```go
if tx.Symbol != "" {
	fmt.Printf("    Symbol=%s Price=%.5f Volume=%.2f\n", ...)
}
```

If the broker transmits quotes and volume — this data is displayed immediately after the transaction.
This allows building custom trading logs in real-time.

---

### 6️. Stream Termination

The stream stops when:

* event limit is reached (`MAX_EVENTS`),
* timeout occurs (`MAX_SECONDS`),
* connection error or stream closure by server.

---

## 📦 What the Server Returns

```protobuf
message OnTradeTransactionData {
  ENUM_TRADE_TRANSACTION_TYPE Type = 1;
  TradeTransaction TradeTransaction = 2;
}

message TradeTransaction {
  ENUM_TRADE_TRANSACTION_TYPE Type = 1;
  uint64 OrderTicket = 2;
  uint64 DealTicket = 3;
  string Symbol = 4;
  double Price = 5;
  double Volume = 6;
}
```

> Each message corresponds to **one trading operation**.

---

## 💡 Example Output

```
Streaming trade transactions (max 5 events or 10 seconds)...
  NOTE: This stream sends events for all trade-related transactions
  Event #1: Type=DEAL_ADD
    Transaction Type=DEAL_ADD Order=12345678 Deal=98765432
    Symbol=EURUSD Price=1.08645 Volume=0.10
  Event #2: Type=ORDER_DELETE
    Transaction Type=ORDER_DELETE Order=12345678 Deal=0
  ✓ Received 5 events, stopping stream
```

---

## 🧠 What It's Used For

The `OnTradeTransaction()` method is used for:

* **auditing** all trading actions on the account;
* **logging** orders, deals, and modifications;
* **reactive trading strategies** responding to changes;
* **analytics** or gateway testing.

---

## 💬 In Simple Terms

> `OnTradeTransaction()` is a **stream of all trading events** in MetaTrader.
> It reports:
> "Order created", "Deal executed", "Order deleted", "Position closed".
>
> This is the most complete and detailed stream, allowing you to see everything happening with orders and deals in real-time.
