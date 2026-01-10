### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolExist()`** method is used to check **whether a specific trading instrument (symbol)** exists on the MetaTrader server.
> This is an important step before any actions with symbols — requesting information, subscribing to quotes, or trading operations.
>
> It also reports **whether the symbol is custom** — created manually, rather than provided by the broker.


---

### 🧩 Code example

```go
fmt.Println("\n4.2. SymbolExist() - Check if symbol exists")

existReq := &pb.SymbolExistRequest{
    Name: cfg.TestSymbol, // Direct protobuf field: Name (not Symbol)
}
existData, err := account.SymbolExist(ctx, existReq)
if !helpers.PrintShortError(err, fmt.Sprintf("SymbolExist(%s) failed", cfg.TestSymbol)) {
    fmt.Printf("  Symbol '%s' exists:      %v\n", cfg.TestSymbol, existData.Exists)
    fmt.Printf("  Is custom symbol:         %v\n", existData.IsCustom)
}
```

---

### 🟢 Detailed Code Explanation

#### 🔹 Creating the request

```go
existReq := &pb.SymbolExistRequest{
    Name: cfg.TestSymbol,
}
```

Creates a `SymbolExistRequest` request structure with the name of the symbol to check.
Usually the value is taken from configuration (e.g., `"EURUSD"` from `appsettings.json`).

#### 🔹 Method call

```go
existData, err := account.SymbolExist(ctx, existReq)
```

Executes a gRPC call to the MetaTrader server.

* `ctx` — context (with or without timeout).
* `existReq` — request structure.
  Returns `existData` (protobuf response) and `err` — error if something went wrong.

#### 🔹 Checking the result

```go
if !helpers.PrintShortError(err, fmt.Sprintf("SymbolExist(%s) failed", cfg.TestSymbol)) {
    fmt.Printf("  Symbol '%s' exists:      %v\n", cfg.TestSymbol, existData.Exists)
    fmt.Printf("  Is custom symbol:         %v\n", existData.IsCustom)
}
```

If there is no error, two fields are printed:

* `existData.Exists` — **true/false**, shows whether the symbol was found on the server;
* `existData.IsCustom` — **true/false**, shows whether the symbol is custom.

---

### Example output

```
4.2. SymbolExist() - Check if symbol exists
  Symbol 'EURUSD' exists:      true
  Is custom symbol:            false
```

---

### What the server returns

```protobuf
message SymbolExistData {
  bool Exists = 1;     // true = symbol found
  bool IsCustom = 2;   // true = custom symbol
}
```

---

### Practical application example

```go
symbol := "BTCUSD"
req := &pb.SymbolExistRequest{Name: symbol}
resp, err := account.SymbolExist(ctx, req)
if err != nil {
    log.Fatalf("Check failed: %v", err)
}

if !resp.Exists {
    log.Fatalf("Symbol %s not found on server", symbol)
}

if resp.IsCustom {
    fmt.Printf("%s is a custom (user-defined) symbol\n", symbol)
} else {
    fmt.Printf("%s is a standard market symbol\n", symbol)
}
```

---

### ⚙️ What is a *Custom Symbol*

**Custom symbol** — is an instrument created manually by the user, rather than provided by the broker.
For example:

* `EURUSD` → standard broker symbol;
* `EURUSD.test` → custom symbol.

Custom symbols are used for:

* quote simulation,
* strategy testing,
* manually loading custom data.

---

### 💡 Practical use in gateway

| Goal                                         | Method                                     |
| -------------------------------------------- | ------------------------------------------ |
| Check symbol existence                       | `SymbolExist()`                            |
| Activate symbol (add to Market Watch)        | `SymbolSelect(Symbol, true)`               |
| Get symbol information                       | `SymbolInfo()`                             |
| Subscribe to quotes                          | `TicksSubscribe()` or `SymbolSubscribe()` |

---

### 💬 In simple terms

> `SymbolExist()` is a check: **"Is there a symbol with this name on the MetaTrader server?"**
> It returns two answers: *does the symbol exist* and *what type is it — standard or custom*.

Used at gateway startup, before trading or informational calls,
to make sure the specified symbol is actually available on the server.
