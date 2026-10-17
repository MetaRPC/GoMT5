# Protobuf Inspector - Interactive Type Explorer

> **Interactive developer utility for exploring MT5 protobuf types, fields, enums, and data structures from the MT5 gRPC API**

---

## 🎯 What This Tool Does

The Protobuf Inspector is an **interactive command-line tool** that helps you explore the structure of the MT5 gRPC API:

- ✅ **Interactive search** for types, fields, and enums
- ✅ **Real-time inspection** of protobuf message structures
- ✅ **Field-level discovery** - find which types contain specific fields
- ✅ **Enum value exploration** - see all possible enum values
- ✅ **Type browsing** - list all available types in the API

---

## 🚀 Getting Started

### Running the Inspector

```bash
cd examples/demos
go run main.go inspect      # or: go run main.go 17 (also: inspector, proto)
```

You will see an interactive prompt:

```
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║          🔍 INTERACTIVE PROTOBUF TYPES INSPECTOR 🔍                ║
║                                                                    ║
║         Explore MT5 gRPC API Types, Fields & Enums                 ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝

Package: github.com/MetaRPC/GoMT5/package

┌────────────────────────────────────────────────────────────────┐
│ QUICK START GUIDE                                              │
├────────────────────────────────────────────────────────────────┤
│ • Type 'help' for full command list                            │
│ • Type 'list' to see all available types                       │
│ • Type a type name to inspect it (e.g., 'PositionInfo')        │
│ • Type 'exit' to quit                                          │
└────────────────────────────────────────────────────────────────┘

>
```

---

## 📖 Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `list` or `ls` | Show all available protobuf types | `list` |
| `<TypeName>` | Inspect specific type structure | `PositionInfo` |
| `search <text>` or `find <text>` | Search for types whose name contains text | `search Order` |
| `field <name>` | Find all types with a field whose name contains `<name>` | `field Balance` |
| `enum <name>` | Show all enum values | `enum BMT5_ENUM_ORDER_TYPE` |
| `help` or `?` | Show help message | `help` |
| `exit`, `quit` or `q` | Exit the inspector | `exit` |

---

## 💡 Practical Examples

### Example 1: Find out what fields PositionInfo has

**Command:**
```
> PositionInfo
```

**Output:**
```
╔════════════════════════════════════════════════════════════╗
║ TYPE: PositionInfo                                        ║
╚════════════════════════════════════════════════════════════╝
📋 Fields (25):

  • Index                                    : uint32
       (protobuf field #1)
  • Ticket                                   : uint64
       (protobuf field #2)
  • OpenTime                                 : *Timestamp
       (protobuf field #3)
  • OpenTimeUtc                              : *Timestamp
       (protobuf field #21)
  • Volume                                   : float64
       (protobuf field #4)
  • PriceOpen                                : float64
       (protobuf field #5)
  ...
  • Type                                     : BMT5_ENUM_POSITION_TYPE
       (protobuf field #12)
  • MagicNumber                              : int64
       (protobuf field #13)
  • Identifier                               : int64
       (protobuf field #14)
  • Reason                                   : BMT5_ENUM_POSITION_REASON
       (protobuf field #15)
  • Symbol                                   : string
       (protobuf field #16)
  • Comment                                  : string
       (protobuf field #17)
  • ExternalId                               : string
       (protobuf field #18)
  • PositionCommission                       : float64
       (protobuf field #19)
  • AccountLogin                             : int64
       (protobuf field #20)

💡 To see field values of an enum, use: enum <EnumName>
```

The names and types are the Go ones (`float64`, `*Timestamp`), exactly as you use them in code. The count in
`Fields (N)` includes the three internal protobuf fields that are not listed.

---

### Example 2: Find which types have the "Ticket" field

**Command:**
```
> field Ticket
```

**Output:**
```
🔍 Searching for field 'Ticket'...

✅ Found in 17 type(s):

📦 DealHistoryData:
   └─ Ticket: uint64

📦 MqlTradeRequest:
   └─ OrderTicket: uint64
   └─ PositionTicket: uint64

📦 MqlTradeResult:
   └─ DealTicket: uint64
   └─ OrderTicket: uint64
...

📦 OrderCloseRequest:
   └─ Ticket: uint64

📦 OrderModifyRequest:
   └─ Ticket: uint64

📦 PositionHistoryInfo:
   └─ PositionTicket: uint64

📦 PositionInfo:
   └─ Ticket: uint64
```

---

### Example 3: See all ORDER_TYPE enum values

**Command:**
```
> enum BMT5_ENUM_ORDER_TYPE
```

**Output:**
```
╔════════════════════════════════════════════════════════════╗
║ ENUM: BMT5_ENUM_ORDER_TYPE                                ║
╚════════════════════════════════════════════════════════════╝
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY           = 0
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL          = 1
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY_LIMIT     = 2
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL_LIMIT    = 3
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY_STOP      = 4
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL_STOP     = 5
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY_STOP_LIMIT = 6
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL_STOP_LIMIT = 7
  BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_CLOSE_BY      = 8
```

The value names are the full Go constant names - in code write `pb.BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY`.

---

### Example 4: Find all types related to "Position"

**Command:**
```
> search Position
```

**Output:**
```
🔍 Searching for types containing 'Position'...

✅ Found 24 type(s):

  [Struct]   OnPositionProfitData
  [Struct]   OnPositionProfitPositionInfo
  [Struct]   OnPositionProfitReply
  ...
  [Struct]   PositionHistoryInfo
  [Struct]   PositionInfo
  [Struct]   PositionsHistoryData
  [Struct]   PositionsHistoryReply
  [Struct]   PositionsHistoryReply_Data
  [Struct]   PositionsHistoryReply_Error
  [Struct]   PositionsHistoryRequest
  [Struct]   PositionsTotalData
  ...

💡 Type a name to inspect it
```

`search` looks at type names only (enums are found with `enum <name>`).

---

### Example 5: Show all available types

**Command:**
```
> list
```

**Output:**
```
╔════════════════════════════════════════════════════════════╗
║              ALL PROTOBUF TYPES                            ║
╚════════════════════════════════════════════════════════════╝

📦 Found 267 types:

  [Type]    AccountInfoDoubleData
  [Reply]   AccountInfoDoubleReply
  [Type]    AccountInfoDoubleReply_Data
  [Type]    AccountInfoDoubleReply_Error
  [Request] AccountInfoDoubleRequest
  ...
  [Reply]   TickValueWithSizeReply
  [Request] TickValueWithSizeRequest

💡 Type a name to inspect it (e.g., 'PositionInfo')
```

---

## 🔍 Common Use Cases

### USE CASE 1: "Getting 'field not found' error"

**Problem:** Your code has `summary.Balance` but it doesn't compile

**Solution:**
```
> field Balance
```

**Result:** Shows the correct field name and which type has it
```
✅ Found in 3 type(s):

📦 AccountSummaryData:
   └─ AccountBalance: float64

📦 MrpcMqlTradeCheckResult:
   └─ BalanceAfterDeal: float64

📦 OnEventAccountInfo:
   └─ Balance: float64
```

**Fix:** Use `summary.AccountBalance` on the `AccountSummaryData` returned by `account.AccountSummary()`

---

### USE CASE 2: "What fields does X have?"

**Problem:** Don't know what data is in `PositionInfo`

**Solution:**
```
> PositionInfo
```

**Result:** Shows all fields (Ticket, Type, Symbol, Profit, etc.)

---

### USE CASE 3: "What are valid enum values?"

**Problem:** Don't know what value to use for `OrderType`

**Solution:**
```
> enum BMT5_ENUM_ORDER_TYPE
```

**Result:** Shows all values:
```
BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY           = 0
BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_SELL          = 1
BMT5_ENUM_ORDER_TYPE_BMT5_ORDER_TYPE_BUY_LIMIT     = 2
...
```

---

### USE CASE 4: "Need to find position-related types"

**Problem:** Exploring the API, need to see all position-related structures

**Solution:**
```
> search Position
```

**Result:** Shows all types with "Position" in the name

---

### USE CASE 5: "Want to browse what's available"

**Problem:** New to the API, want to explore

**Solution:**
```
> list
```

**Result:** Shows all available types, sorted alphabetically and tagged `[Request]`, `[Reply]`, `[Info]` or `[Type]`

---

## 📊 Statistics

- **Total Types:** 267 (all MT5 gRPC protobuf message types)
- **Total Enums:** 67 (with 943 enum values)
- **Coverage:** 100% of MT5 gRPC API

---

## 🔑 Important Enums (Frequently Used)

| Enum Name | Description | Common Values |
|-----------|-------------|---------------|
| `BMT5_ENUM_ORDER_TYPE` | Order types | BUY, SELL, BUY_LIMIT, SELL_LIMIT, BUY_STOP, SELL_STOP |
| `BMT5_ENUM_ORDER_TYPE_FILLING` | Fill policies | FOK, IOC, Return, BOC |
| `BMT5_ENUM_ORDER_TYPE_TIME` | Time in force | GTC, Day, Specified, Specified Day |
| `BMT5_ENUM_DEAL_REASON` | Deal execution reason | Client, Expert, SL, TP, Mobile |
| `BMT5_ENUM_DEAL_ENTRY_TYPE` | Deal entry type | In, Out, InOut, Out By |
| `BMT5_ENUM_POSITION_TYPE` | Position direction | BUY, SELL |
| `BMT5_ENUM_POSITION_REASON` | Why position opened | Client, Expert, Mobile, Web, SL, TP, SO |
| `MqlErrorCode` | MQL error codes | 211 different error codes |
| `MqlErrorTradeCode` | Trading operation errors | REQUOTE, REJECT, MARKET_CLOSED, etc. |
| `SymbolInfoDoubleProperty` | Symbol price properties | BID, ASK, POINT, SWAP_LONG, SWAP_SHORT |
| `SymbolInfoIntegerProperty` | Symbol integer properties | DIGITS, SPREAD, TRADE_MODE, etc. |
| `MRPC_ENUM_TRADE_REQUEST_ACTIONS` | Trade actions | DEAL, PENDING, SLTP, MODIFY, REMOVE |

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **Case-insensitive search** | `search Order` = `search order`; type and enum names are matched case-insensitively too |
| **Partial field matching** | `field profit` finds both `Profit` and `TakeProfit` |
| **Type categorization** | `list` shows `[Request]`, `[Reply]`, `[Type]`, `[Info]` tags |
| **Array indicators** | 📚 icon for repeated/array fields |
| **Protobuf field numbers** | Shows field `#N` for each field |
| **Helpful hints** | Suggests `search <name>` / `list` when a type is not found |

---

## 🎬 Example Interactive Session

```bash
$ go run main.go inspect

╔════════════════════════════════════════════════════════════════════╗
║          🔍 INTERACTIVE PROTOBUF TYPES INSPECTOR 🔍                ║
╚════════════════════════════════════════════════════════════════════╝

> search Order

🔍 Searching for types containing 'Order'...

✅ Found 56 type(s):

  [Struct]   OnPositionsAndPendingOrdersTicketsData
  ...
  [Struct]   OpenedOrderInfo
  [Struct]   OpenedOrdersData
  ...
  [Struct]   OrderSendRequest
  ...

> OpenedOrderInfo

╔════════════════════════════════════════════════════════════╗
║ TYPE: OpenedOrderInfo                                     ║
╚════════════════════════════════════════════════════════════╝
📋 Fields (30):

  • Index                                    : uint32
       (protobuf field #1)
  • Ticket                                   : uint64
       (protobuf field #2)
  ...
  • Type                                     : BMT5_ENUM_ORDER_TYPE
  • State                                    : BMT5_ENUM_ORDER_STATE
  ...
  • TypeFilling                              : BMT5_ENUM_ORDER_TYPE_FILLING
  • TypeTime                                 : BMT5_ENUM_ORDER_TYPE_TIME
  ...

> field Magic

🔍 Searching for field 'Magic'...

✅ Found in 10 type(s):

📦 MqlTradeRequest:
   └─ Magic: uint64

📦 MrpcMqlTradeRequest:
   └─ ExpertAdvisorMagicNumber: uint64
...

📦 PositionInfo:
   └─ MagicNumber: int64

> enum BMT5_ENUM_ORDER_STATE

╔════════════════════════════════════════════════════════════╗
║ ENUM: BMT5_ENUM_ORDER_STATE                               ║
╚════════════════════════════════════════════════════════════╝
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_STARTED     = 0
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_PLACED      = 1
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_CANCELED    = 2
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_PARTIAL     = 3
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_FILLED      = 4
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REJECTED    = 5
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_EXPIRED     = 6
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REQUEST_ADD = 7
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REQUEST_MODIFY = 8
  BMT5_ENUM_ORDER_STATE_BMT5_ORDER_STATE_REQUEST_CANCEL = 9

> exit

👋 Goodbye!
```

---


## 💻 Implementation Details

The Protobuf Inspector uses Go reflection to:

1. Register all protobuf types at startup
2. Build an in-memory index of types, fields, and enums
3. Provide instant search and lookup
4. Format output with box drawing and emoji markers

**Source file:** `examples/demos/helpers/17_protobuf_inspector.go`

**Registration:** the 267 types are registered from a fixed list of sample values, and all 943 enum values of the 67 enums
are written out in the code. A type or enum added to the API later appears only after it is added to that file.

---

## 🔧 Technical Notes

- **No MT5 connection required** - This is a purely offline tool that inspects type definitions
- **Complete coverage** - All 267 types and 67 enums from the MT5 gRPC API
- **Instant search** - In-memory index for fast lookup
- **Development only** - Not intended for production use

---

## 🎯 When to Use This Tool

✅ **Use the inspector when:**

- Learning the MT5 gRPC API structure
- Debugging "field not found" errors
- Exploring available protobuf types
- Looking up enum values
- Finding the correct request/response types for API calls
- Understanding message structures before writing code

❌ **Don't use for:**

- Inspecting runtime data (use debugger instead)
- Production code (this is a development tool)
- Testing API connectivity (use demo connections instead)

---

## 📝 Tips and Tricks

1. **Start with search** - If you know the general area, use `search <keyword>` first
2. **Use field search** - When you see a field name but don't know which type, use `field <name>`
3. **Explore enums early** - Understanding enum values saves debugging time later
4. **List is your friend** - When stuck, use `list` to browse available types
5. **Case doesn't matter** - Type commands in lowercase, it's faster

---

## 🚀 Quick Start Workflow

**Beginner workflow for exploring the API:**

```bash
# 1. Start the inspector
go run main.go inspect

# 2. Browse what's available
> list

# 3. Search for what you need
> search Position

# 4. Inspect a type
> PositionInfo

# 5. Check enum values
> enum BMT5_ENUM_POSITION_TYPE

# 6. Find related types
> field Ticket
```

---

## 🎓 Learning Session Example

**Goal:** "I want to close a position, what do I need?"

```bash
> search close
✅ Found 7 type(s):

  [Struct]   CloseChartForSymbolReply
  [Struct]   CloseChartForSymbolRequest
  [Struct]   OrderCloseData
  [Struct]   OrderCloseReply
  [Struct]   OrderCloseReply_Data
  [Struct]   OrderCloseReply_Error
  [Struct]   OrderCloseRequest

> OrderCloseRequest
  • Ticket                                   : uint64
       (protobuf field #1)
  • Volume                                   : float64
       (protobuf field #2)
  • Slippage                                 : int32
       (protobuf field #5)

> OrderCloseData
  • ReturnedCode                             : uint32
       (protobuf field #1)
  • ReturnedStringCode                       : string
       (protobuf field #2)
  • ReturnedCodeDescription                  : string
       (protobuf field #3)
  • CloseMode                                : MRPC_ORDER_CLOSE_MODE
       (protobuf field #4)

# Now you know:
# - Call account.OrderClose with an OrderCloseRequest (Ticket; Volume for a partial close)
# - You'll get OrderCloseData back
# - Check ReturnedCode == 10009 for success (a refused request comes back as an error instead)
```

`search` matches one piece of text: `search position close` finds nothing, `search close` does.

---

## 🆘 Troubleshooting

**Q: Type not found**
```
> MyType

❌ Type 'MyType' not found
💡 Try: 'search MyType' to find similar types
💡 Or:  'list' to see all available types
```

**A:** Use search with partial name to find similar types

---

**Q: Too many results**

```
> search data
✅ Found 91 type(s):
...
```

**A:** Be more specific in your search query

---

**Q: What's the difference between OpenedOrderInfo and OrderSendRequest?**

**A:** Use the inspector:

- `OpenedOrderInfo` - Information about an existing (pending) order
- `OrderSendRequest` - Request to create a new order

Rule:

- `*Request` - Input for API method
- `*Data` - Output of the API method (what MT5Account methods return)
- `*Reply` - gRPC reply wrapper that carries either `Data` or `Error`
- `*Info` - Structured data about entities

---

## 🎉 Summary

The Protobuf Inspector is your **first stop** when working with the MT5 gRPC API. Use it to:

1. 🔍 **Discover** available types
2. 📖 **Learn** message structures
3. 🐛 **Debug** field name issues
4. ✅ **Verify** enum values
5. 🚀 **Speed up** development

**Remember:** Type `help` at any time for command reference!

---

**Next Steps:**

- Run `go run main.go inspect` and explore!
- Check [MT5Account Master Overview](../MT5Account/MT5Account.Master.Overview.md) for complete API documentation
- Try the demo examples in `examples/demos/`
