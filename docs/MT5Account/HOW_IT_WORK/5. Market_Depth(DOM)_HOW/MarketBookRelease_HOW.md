### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`MarketBookRelease()`** method terminates the **Depth of Market (DOM)** subscription that was previously opened by the `MarketBookAdd()` method.
> It is used for proper resource cleanup and closing the connection with the server.


---

## 🧩 Code example

```go
fmt.Println("\n6.3. MarketBookRelease() - Unsubscribe from DOM")
fmt.Println("    Cleaning up subscription...")

bookReleaseCtx, bookReleaseCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer bookReleaseCancel()

marketBookReleaseReq := &pb.MarketBookReleaseRequest{
    Symbol: cfg.TestSymbol, // symbol from config.json (test_symbol)
}
marketBookReleaseData, err := account.MarketBookRelease(bookReleaseCtx, marketBookReleaseReq)
if err != nil {
    fmt.Printf("  ⚠️  Unsubscribe failed: %v\n", err)
    fmt.Println("     → Not critical - connection will close anyway")
} else {
    if marketBookReleaseData.Success {
        fmt.Println("  ✓ Unsubscribed successfully")
    } else {
        fmt.Println("  ⚠️  Unsubscribe reported unsuccessful")
    }
}
```

---

### 🟢 Detailed Code Explanation

* **Create context with timeout:**

  ```go
  bookReleaseCtx, bookReleaseCancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer bookReleaseCancel()
  ```

  A 5-second timeout prevents hanging when communicating with the server.

* **Form the request:**

  ```go
  marketBookReleaseReq := &pb.MarketBookReleaseRequest{ Symbol: cfg.TestSymbol }
  ```

  Specify the symbol for which the DOM subscription was previously made.

* **Send the request:**

  ```go
  marketBookReleaseData, err := account.MarketBookRelease(bookReleaseCtx, marketBookReleaseReq)
  ```

  A gRPC call is executed, requesting subscription closure.

* **Error handling:**

  ```go
  if err != nil {
      fmt.Printf("  ⚠️  Unsubscribe failed: %v\n", err)
      fmt.Println("     → Not critical - connection will close anyway")
  }
  ```

  Even if an error occurs, it's not critical — all subscriptions will be automatically closed when the session ends.

* **Successful or unsuccessful unsubscribe:**

  ```go
  if marketBookReleaseData.Success {
      fmt.Println("  ✓ Unsubscribed successfully")
  } else {
      fmt.Println("  ⚠️  Unsubscribe reported unsuccessful")
  }
  ```

  If `Success = true`, the broker confirmed successful subscription closure.

---

## 📦 What the Server Returns

```protobuf
message MarketBookReleaseData {
  bool Success = 1; // True if unsubscription was successful
}
```

---

##  Output Examples

### ✅ Successful Unsubscribe

```
✓ Unsubscribed successfully
```

###  Unsubscribe Error

```
⚠️  Unsubscribe failed: subscription not found
     → Not critical - connection will close anyway
```

---

## 🧠 What It's Used For

The `MarketBookRelease()` method is used for:

* properly terminating DOM subscription;
* freeing server resources;
* demonstrating the complete market depth workflow cycle (Add → Get → Release);
* ensuring clean connection closure in tests or demonstrations.

---

## 💬 In Simple Terms

> `MarketBookRelease()` is "disconnecting from the order book".
> If you subscribed to DOM via `MarketBookAdd()`, this call neatly closes the subscription,
> freeing the connection and preventing accumulation of excess resources on the server.
