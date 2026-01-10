### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolName()`** method is used to get **the symbol name (e.g., `"EURUSD"`, `"BTCUSD"`) by its index**
> in the MetaTrader symbol list.
>
> It allows you to programmatically iterate through all symbols stored on the server,
> and output their names — both from the list of *all symbols* and from *Market Watch*.


---

### 🧩 Code example

```go
fmt.Println("\n4.3. SymbolName() - Get symbol name by index (from Market Watch)")

// Use actual count from SymbolsTotal
var symbolsToShow int32 = 3
if selectedSymbolsData != nil && selectedSymbolsData.Total < symbolsToShow {
    symbolsToShow = selectedSymbolsData.Total
}

if symbolsToShow == 0 {
    fmt.Println("  No symbols in Market Watch")
} else {
    fmt.Printf("  Showing first %d symbols from Market Watch:\n", symbolsToShow)
    for i := int32(0); i < symbolsToShow; i++ {
        nameReq := &pb.SymbolNameRequest{
            Index:    i,    // Direct protobuf field: Index (not Pos)
            Selected: true, // true = Market Watch, false = all symbols
        }
        nameData, err := account.SymbolName(ctx, nameReq)
        if !helpers.PrintShortError(err, fmt.Sprintf("SymbolName(pos=%d) failed", i)) {
            if nameData.Name != "" {
                fmt.Printf("    [%d] %s\n", i, nameData.Name)
            } else {
                fmt.Printf("    [%d] (empty - no symbol at this position)\n", i)
            }
        }
    }
}
```

---

### 🟢 Detailed Code Explanation

#### 🔹 Determining the number of symbols to display

```go
var symbolsToShow int32 = 3
if selectedSymbolsData != nil && selectedSymbolsData.Total < symbolsToShow {
    symbolsToShow = selectedSymbolsData.Total
}
```

Sets how many symbols need to be shown (default 3).
If *Market Watch* has fewer than three symbols — uses the actual count from `SymbolsTotal(Mode:true)`.

---

#### 🔹 Checking for symbols presence

```go
if symbolsToShow == 0 {
    fmt.Println("  No symbols in Market Watch")
}
```

If there are no active symbols, a warning is printed.

---

#### 🔹 Main loop

```go
for i := int32(0); i < symbolsToShow; i++ {
    nameReq := &pb.SymbolNameRequest{
        Index:    i,    // Symbol position
        Selected: true, // true = Market Watch, false = all symbols
    }
    nameData, err := account.SymbolName(ctx, nameReq)
```

Creates a `SymbolNameRequest` request, where:

* `Index` — sequential symbol number (0, 1, 2...);
* `Selected` — selects the list source: Market Watch (`true`) or all symbols (`false`).

The `SymbolName()` method returns the symbol name at this position.

---

#### 🔹 Output result

```go
if nameData.Name != "" {
    fmt.Printf("    [%d] %s\n", i, nameData.Name)
} else {
    fmt.Printf("    [%d] (empty - no symbol at this position)\n", i)
}
```

If the name is not empty — the symbol itself is printed.
If empty — it means there is no active instrument at this position.

---

### Example output

```
4.3. SymbolName() - Get symbol name by index (from Market Watch)
  Showing first 3 symbols from Market Watch:
    [0] EURUSD
    [1] GBPUSD
    [2] USDJPY
```

---

### 💡 How it works internally

MetaTrader stores symbols in two lists:

1. **All Symbols** — all available broker instruments.
2. **Market Watch** — active instruments for which quotes are received.

The `SymbolName()` method returns the symbol name by index in the selected list.

| Parameter  | Value        | Description                              |
| ---------- | ------------ | ---------------------------------------- |
| `Index`    | `int32`      | Symbol position in list (0, 1, 2...)     |
| `Selected` | `true/false` | true — Market Watch, false — All Symbols |

---

### ⚙️ Application

| Goal                                                | How to use                                         |
| --------------------------------------------------- | -------------------------------------------------- |
| Get list of all active symbols                      | `Selected: true`, loop over `SymbolsTotal(Mode:true)`   |
| Get list of all broker symbols                      | `Selected: false`, loop over `SymbolsTotal(Mode:false)` |
| Check availability of specific instruments          | after getting names, call `SymbolExist()`          |
| Create symbol directory                             | collect names for logs, interface, reports         |

---

### 💬 In simple terms

> `SymbolName()` is a way to **programmatically scroll through the list of MetaTrader instruments**.
> You give an index (0, 1, 2...) and get the symbol name,
> and you can choose — from *Market Watch* or from *the entire symbol list*.

Used for diagnostics, gateway debugging, monitoring, or building interfaces
where you need to dynamically get a list of available instruments.
