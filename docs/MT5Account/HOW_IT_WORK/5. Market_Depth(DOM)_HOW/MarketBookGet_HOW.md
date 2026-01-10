### Example from file: `examples/demos/lowlevel/01_general_operations.go`

> The **`MarketBookGet()`** method makes a one-time request to retrieve a **market depth snapshot (Depth of Market, DOM)** for a specified symbol.
> If the broker supports DOM, price levels and order volumes are returned.


---

## 🧩 Code example

```go
fmt.Println("\n6.2. MarketBookGet() - Get DOM snapshot")
fmt.Println("    Testing if data is actually available...")

bookGetCtx, bookGetCancel := context.WithTimeout(context.Background(), 3*time.Second)
defer bookGetCancel()

marketBookGetReq := &pb.MarketBookGetRequest{
    Symbol: cfg.TestSymbol, // symbol from config.json (test_symbol parameter)
}
marketBookGetData, err := account.MarketBookGet(bookGetCtx, marketBookGetReq)
if err != nil {
    fmt.Println("  ❌ No data available (timeout)")
    fmt.Println("     → This is EXPECTED for forex pairs on demo accounts")
    fmt.Println("     → Not an error - just a limitation of forex market structure")
} else {
    if len(marketBookGetData.Book) > 0 {
        fmt.Printf("  ✅ SUCCESS! Received %d price levels\n", len(marketBookGetData.Book))
        fmt.Println("     → This is UNUSUAL for forex - broker provides limited DOM")
        fmt.Println()

        maxShow := 5
        if len(marketBookGetData.Book) < maxShow {
            maxShow = len(marketBookGetData.Book)
        }
        fmt.Println("  First few levels:")
        for i := 0; i < maxShow; i++ {
            level := marketBookGetData.Book[i]
            bookType := "SELL"
            if level.Type == pb.BookType_BOOK_TYPE_BUY {
                bookType = "BUY"
            }
            fmt.Printf("    [%d] %s  Price: %.5f  Volume: %.2f\n", i+1, bookType, level.Price, level.VolumeDouble)
        }
    } else {
        fmt.Println("  ⚠️  Call succeeded but no data returned")
        fmt.Println("     → Broker has no DOM data to show")
    }
}
```

---

### 🟢 Detailed Code Explanation

* **Create context with timeout:**

  ```go
  bookGetCtx, bookGetCancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer bookGetCancel()
  ```

  A 3-second timeout is needed to avoid hanging if the server doesn't provide data.

* **Form the request:**

  ```go
  marketBookGetReq := &pb.MarketBookGetRequest{ Symbol: cfg.TestSymbol }
  ```

  The symbol is taken from the `config.json` configuration file (`test_symbol` field).

* **Send the request:**

  ```go
  marketBookGetData, err := account.MarketBookGet(bookGetCtx, marketBookGetReq)
  ```

  Returns the `Book` array — a list of price levels with volumes.

* **Error (no data):**

  ```go
  if err != nil {
      fmt.Println("  ❌ No data available (timeout)")
  }
  ```

  Means the broker doesn't provide DOM (common on demo or forex pairs).

* **Data available:**

  ```go
  if len(marketBookGetData.Book) > 0 {
      fmt.Printf("  ✅ SUCCESS! Received %d price levels\n", len(marketBookGetData.Book))
  }
  ```

  Displays the number of order book levels. This is rare for forex.

* **Display first levels:**

  ```go
  for i := 0; i < maxShow; i++ {
      level := marketBookGetData.Book[i]
      fmt.Printf("[%d] %s Price: %.5f Volume: %.2f\n", i+1, bookType, level.Price, level.VolumeDouble)
  }
  ```

  Shows the type (`BUY` or `SELL`), price, and volume for the first 5 levels.

---

## 📦 What the Server Returns

```protobuf
message MarketBookGetData {
  repeated BookStruct Book = 1;
}

message BookStruct {
  BookType Type = 1;      // BUY or SELL
  double Price = 2;       // price level
  int64 Volume = 3;       // volume at level (deprecated)
  double VolumeDouble = 4; // volume at level with extended precision
}
```

---

##  Example Successful Output

```
✅ SUCCESS! Received 5 price levels
     → This is UNUSUAL for forex - broker provides limited DOM

  First few levels:
    [1] SELL  Price: 1.08610  Volume: 50000.00
    [2] SELL  Price: 1.08608  Volume: 30000.00
    [3] BUY   Price: 1.08595  Volume: 40000.00
    [4] BUY   Price: 1.08592  Volume: 20000.00
```

---

##  If No Data Available

```
❌ No data available (timeout)
     → This is EXPECTED for forex pairs on demo accounts
```

This is normal: on forex and demo accounts, brokers usually don't form an order book.

---

## 🧠 What It's Used For

The `MarketBookGet()` method is used for:

* one-time retrieval of current DOM snapshot;
* checking if the broker supports market depth;
* liquidity analysis or API testing.
