### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolIsSynchronized()`** method is used to check **whether symbol data is synchronized between the MetaTrader server and the gateway**.
> It helps determine if the instrument data is ready for use — whether you can already request quotes, parameters, or perform trading operations.


---

### 🧩 Code example

```go
fmt.Println("\n4.5. SymbolIsSynchronized() - Check sync status with server")

syncReq := &pb.SymbolIsSynchronizedRequest{
    Symbol: cfg.TestSymbol,
}
syncData, err := account.SymbolIsSynchronized(ctx, syncReq)
if !helpers.PrintShortError(err, fmt.Sprintf("SymbolIsSynchronized(%s) failed", cfg.TestSymbol)) {
    fmt.Printf("  Symbol '%s' synchronized:  %v\n", cfg.TestSymbol, syncData.Synchronized)
}
```

---

### 🟢 Detailed Code Explanation

#### 🔹 Forming the request

```go
syncReq := &pb.SymbolIsSynchronizedRequest{
    Symbol: cfg.TestSymbol,
}
```

Creates a `SymbolIsSynchronizedRequest` structure specifying the symbol
for which the synchronization state needs to be checked (e.g., `"EURUSD"`).

---

#### 🔹 Method call

```go
syncData, err := account.SymbolIsSynchronized(ctx, syncReq)
```

Executes a gRPC request to the MetaTrader server.

* `ctx` — request context (with possible timeout);
* `syncReq` — structure with parameters;
* `syncData` — server response as a protobuf structure;
* `err` — error if the server didn't respond.

---

#### 🔹 Processing the result

```go
if !helpers.PrintShortError(err, fmt.Sprintf("SymbolIsSynchronized(%s) failed", cfg.TestSymbol)) {
    fmt.Printf("  Symbol '%s' synchronized:  %v\n", cfg.TestSymbol, syncData.Synchronized)
}
```

If there are no errors, the value `syncData.Synchronized` is printed — a boolean field (true/false),
showing whether the symbol data is synchronized between the server and gateway.

📘 Example output:

```
4.5. SymbolIsSynchronized() - Check sync status with server
  Symbol 'EURUSD' synchronized:  true
```

---

### What the server returns

```protobuf
message SymbolIsSynchronizedData {
  bool Synchronized = 1; // true = symbol data is synchronized
}
```

---

### 💡 Method purpose

`SymbolIsSynchronized()` helps determine **whether symbol data is ready for use** in code.
If the symbol was recently added to Market Watch via `SymbolSelect()`,
it may need some time for the gateway to receive and update all properties and quotes.

| State   | Value               | Description                                       |
| ------- | ------------------- | ------------------------------------------------- |
| `true`  | Synchronized        | All data received, can work with the symbol       |
| `false` | Not synchronized    | Data not yet updated, should wait                 |

---

### ⚙️ Code usage example

```go
symbol := cfg.TestSymbol

// Check if symbol is synchronized
syncResp, _ := account.SymbolIsSynchronized(ctx, &pb.SymbolIsSynchronizedRequest{
    Symbol: symbol,
})

if !syncResp.Synchronized {
    fmt.Printf("Symbol %s not yet synchronized. Waiting...\n", symbol)
} else {
    fmt.Printf("Symbol %s ready and synchronized.\n", symbol)
}
```

---

### 🧠 Practical application when working with code

The method is useful:

* before calling `SymbolInfo*()` methods — to make sure symbol properties are already loaded;
* after `SymbolSelect()` — to check data readiness;
* when debugging the gateway — to ensure synchronization with the server is working correctly.

---

### 💬 In simple terms

> `SymbolIsSynchronized()` is an indicator of **whether the instrument data is current**.
> If `true` — the gateway has already received all symbol data from the MetaTrader server,
> if `false` — data is not yet synchronized, and it's better to wait before subsequent calls.
