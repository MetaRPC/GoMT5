### Example from file: `examples/demos/lowlevel/01_general_operations.go`


> The **`SymbolInfoString()`** method is used to get **string properties of a symbol** — text data associated with a trading instrument.
> It is similar to `SymbolInfoDouble()` and `SymbolInfoInteger()`, but returns values of type `string`.


---

## 1️⃣ Getting symbol description (Description)

```go
descReq := &pb.SymbolInfoStringRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoStringProperty_SYMBOL_DESCRIPTION, // Direct field: Type (not PropertyId)
}
descData, err := account.SymbolInfoString(ctx, descReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoString(DESCRIPTION) failed")
} else {
    fmt.Printf("  Description:                   %s\n", descData.Value)
}
```

📘 **SYMBOL_DESCRIPTION** — human-readable description of the instrument, set by the broker or MetaTrader server.
Example values:

* `"Euro vs US Dollar"`
* `"Bitcoin / US Dollar"`
* `"Gold Spot"`

Used for displaying information in interfaces, logs, and reports.

---

## 2️⃣ Getting base currency (Base Currency)

```go
baseReq := &pb.SymbolInfoStringRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_BASE,
}
baseData, err := account.SymbolInfoString(ctx, baseReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoString(CURRENCY_BASE) failed")
} else {
    fmt.Printf("  Base currency:                 %s\n", baseData.Value)
}
```

📘 **SYMBOL_CURRENCY_BASE** — is the first currency in the pair, showing what is being bought or sold.
Examples:

* `EURUSD` → base currency `EUR`
* `BTCUSD` → base currency `BTC`

If you buy `EURUSD`, you are buying euros for dollars.

---

## 3️⃣ Getting profit currency (Profit Currency)

```go
profitCurrReq := &pb.SymbolInfoStringRequest{
    Symbol: cfg.TestSymbol,
    Type:   pb.SymbolInfoStringProperty_SYMBOL_CURRENCY_PROFIT,
}
profitCurrData, err := account.SymbolInfoString(ctx, profitCurrReq)
if err != nil {
    helpers.PrintShortError(err, "SymbolInfoString(CURRENCY_PROFIT) failed")
} else {
    fmt.Printf("  Profit currency:               %s\n", profitCurrData.Value)
}
```

📘 **SYMBOL_CURRENCY_PROFIT** — the second currency in the pair, in which profit and loss are calculated.
Examples:

* `EURUSD` → profit currency `USD`
* `GBPJPY` → profit currency `JPY`

---

### 📦 What the server returns

After calling `SymbolInfoString()`, the gateway returns a structure:

```protobuf
message SymbolInfoStringData {
  string Value = 1; // String value of the requested property
}
```

The `Value` field contains the text value corresponding to the selected `Type`.

---

Below is a list of the most commonly used properties that can be requested via `SymbolInfoString()`:

### 📘 All PropertyId enum values (SymbolInfoStringProperty)

| Constant | Value | Description |
|----------|-------|-------------|
| `SYMBOL_BASIS` | 0 | Underlying asset of derivative |
| `SYMBOL_CATEGORY` | 1 | Symbol category |
| `SYMBOL_COUNTRY` | 2 | Country the symbol belongs to |
| `SYMBOL_SECTOR_NAME` | 3 | Economic sector name |
| `SYMBOL_INDUSTRY_NAME` | 4 | Industry name |
| `SYMBOL_CURRENCY_BASE` | 5 | Base currency (first currency in pair) |
| `SYMBOL_CURRENCY_PROFIT` | 6 | Profit currency (second currency in pair) |
| `SYMBOL_CURRENCY_MARGIN` | 7 | Margin calculation currency |
| `SYMBOL_BANK` | 8 | Source of last quote (bank/exchange) |
| `SYMBOL_DESCRIPTION` | 9 | Symbol description/full name |
| `SYMBOL_EXCHANGE` | 10 | Exchange the symbol is traded on |
| `SYMBOL_FORMULA` | 11 | Formula for custom symbol price calculation |
| `SYMBOL_ISIN` | 12 | International Securities Identification Number |
| `SYMBOL_PAGE` | 13 | Web page URL with symbol information |
| `SYMBOL_PATH` | 14 | Symbol path in Market Watch tree |

---

### 🧠 What it's used for

The `SymbolInfoString()` method is used:

* for **displaying information** about the symbol in interfaces, logs, reports;
* when **initializing trading modules**, to determine the calculation currency;
* in **analytical systems**, where currency codes and descriptions need to be matched;
* for **demo examples**, to nicely display instrument properties.

---

### 💬 In simple terms

> `SymbolInfoString()` is a way to request **text information about a symbol**:
> description, base currency, profit currency, and other string properties.
> Used for displaying human-readable data, rather than numeric parameters.
