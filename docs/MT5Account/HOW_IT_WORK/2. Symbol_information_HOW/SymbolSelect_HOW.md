### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolSelect()`** method is used to **add or remove a symbol from Market Watch** — the list of active instruments for which the gateway receives quotes.
> This is one of the key tools when working with symbols: it controls which instruments are *active* in the system.


---

### 🧩 Code example

```go
fmt.Println("\n4.4. SymbolSelect() - Add/remove symbol from Market Watch")

selectReq := &pb.SymbolSelectRequest{
    Symbol: cfg.TestSymbol,
    Select: true, // true = add to Market Watch, false = remove
}
selectData, err := account.SymbolSelect(ctx, selectReq)
if !helpers.PrintShortError(err, fmt.Sprintf("SymbolSelect(%s) failed", cfg.TestSymbol)) {
    fmt.Printf("  Symbol '%s' added to Market Watch: %v\n", cfg.TestSymbol, selectData.Success)
}
```

---

### 🟢 Detailed Code Explanation

#### 🔹 Creating the request

```go
selectReq := &pb.SymbolSelectRequest{
    Symbol: cfg.TestSymbol,
    Select: true, // true = add to Market Watch, false = remove
}
```

Forms the `SymbolSelectRequest` request structure.

| Field    | Type   | Description                                          |
| -------- | ------ | ---------------------------------------------------- |
| `Symbol` | string | Symbol name, e.g. `"EURUSD"`                         |
| `Select` | bool   | `true` — add, `false` — remove from Market Watch     |

The symbol value is taken from the configuration file `appsettings.json` (field `test_symbol`).

---

#### 🔹 Method call

```go
selectData, err := account.SymbolSelect(ctx, selectReq)
```

Executes a gRPC request to the MetaTrader server.

* `ctx` — request context (with possible timeout);
* `selectReq` — structure with parameters;
* `selectData` — server response (protobuf structure);
* `err` — error if the request failed.

---

#### 🔹 Processing the result

```go
if !helpers.PrintShortError(err, fmt.Sprintf("SymbolSelect(%s) failed", cfg.TestSymbol)) {
    fmt.Printf("  Symbol '%s' added to Market Watch: %v\n", cfg.TestSymbol, selectData.Success)
}
```

If there is no error, the result from `selectData.Success` is printed.

📘 `Success = true` means the operation was completed successfully.

Example output:

```
4.4. SymbolSelect() - Add/remove symbol from Market Watch
  Symbol 'EURUSD' added to Market Watch: true
```

---

### What the server returns

```protobuf
message SymbolSelectData {
  bool Success = 1; // true = operation completed successfully
}
```

---

### 💡 Method purpose

`SymbolSelect()` controls the activity of symbols in the *Market Watch* list.

| Action        | Example call    | Description                                      |
| ------------- | --------------- | ------------------------------------------------ |
| Add symbol    | `Select: true`  | Activates the symbol (enables quote stream)     |
| Remove symbol | `Select: false` | Removes the symbol from Market Watch             |

If a symbol is not included in Market Watch, MetaTrader **does not send quotes for it**,
and methods like `SymbolInfoDouble()` or streams (`OnSymbolTick`) will not return current data.

---

### ⚙️ Usage example at gateway startup

```go
symbol := cfg.TestSymbol

// 1. Check if symbol exists
existResp, _ := account.SymbolExist(ctx, &pb.SymbolExistRequest{Name: symbol})
if !existResp.Exists {
    log.Fatalf("Symbol %s not found on server", symbol)
}

// 2. Activate symbol in Market Watch
selectResp, _ := account.SymbolSelect(ctx, &pb.SymbolSelectRequest{
    Symbol: symbol,
    Select: true,
})
fmt.Printf("Symbol %s added to Market Watch: %v\n", symbol, selectResp.Success)

// 3. Check if we can receive quotes via stream
fmt.Println("Waiting for OnSymbolTick stream...")
```

---

### 🧠 In practice

`SymbolSelect()` is one of the first calls after connecting to the server,
because it is what **enables the quote stream** for the required instruments.

Typical sequence during gateway initialization:

1. `AccountSummary()` — check account access.
2. `SymbolExist()` — make sure the symbol is available from the broker.
3. `SymbolSelect()` — add symbol to Market Watch.
4. `OnSymbolTick()` — start receiving quote stream.

---

### 💬 In simple terms

> `SymbolSelect()` is a **switch** for symbols.
> When you set `Select = true`, the gateway tells MetaTrader:
> "Add this instrument to Market Watch and start streaming quotes for it".
> If `Select = false`, the opposite — "remove the symbol from active ones and stop sending data".
