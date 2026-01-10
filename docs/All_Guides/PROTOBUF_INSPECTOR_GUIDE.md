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
go run main.go inspect
```

You will see an interactive prompt:

```
═══════════════════════════════════════════════════════
    MT5 PROTOBUF TYPE INSPECTOR
═══════════════════════════════════════════════════════

Available: 267 types, 68 enums

Type 'help' for available commands
Type 'list' to see all types
Type a type name to inspect it

>
```

---

## 📖 Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `list` or `ls` | Show all available protobuf types | `list` |
| `<TypeName>` | Inspect specific type structure | `PositionInfo` |
| `search <text>` or `find <text>` | Search for types containing text | `search Order` |
| `field <name>` | Find all types with a specific field | `field Balance` |
| `enum <name>` | Show all enum values | `enum BMT5_ENUM_ORDER_TYPE` |
| `help` or `?` | Show help message | `help` |
| `exit` or `quit` | Exit the inspector | `exit` |

---

## 💡 Practical Examples

### Example 1: Find out what fields PositionInfo has

**Command:**
```
> PositionInfo
```

**Output:**
```
═══════════════════════════════════════════════════════
TYPE: PositionInfo
═══════════════════════════════════════════════════════

Fields:
  Ticket: uint64 (#1)
  Type: BMT5_ENUM_POSITION_TYPE (#2)
  Symbol: string (#3)
  Magic: uint64 (#4)
  Identifier: uint64 (#5)
  Reason: BMT5_ENUM_POSITION_REASON (#6)
  Volume: double (#7)
  PriceOpen: double (#8)
  Sl: double (#9)
  Tp: double (#10)
  PriceCurrent: double (#11)
  Swap: double (#12)
  Profit: double (#13)
  Time: uint64 (#14)
  TimeUpdate: uint64 (#15)
  TimeUpdateMsc: uint64 (#16)
  Comment: string (#17)
  ExternalId: string (#18)
```

---

### Example 2: Find which types have the "Ticket" field

**Command:**
```
> field Ticket
```

**Output:**
```
Field 'Ticket' found in 15 types:
  • PositionInfo
  • OrderInfo
  • DealInfo
  • HistoryOrdersGetByTicketRequest
  • HistoryDealsGetByTicketRequest
  • OrderCloseRequest
  • OrderModifyRequest
  • PositionCloseRequest
  • PositionModifyRequest
  ...
```

---

### Example 3: See all ORDER_TYPE enum values

**Command:**
```
> enum BMT5_ENUM_ORDER_TYPE
```

**Output:**
```
═══════════════════════════════════════════════════════
ENUM: BMT5_ENUM_ORDER_TYPE
═══════════════════════════════════════════════════════

Values:
  BMT5_ORDER_TYPE_BUY         = 0
  BMT5_ORDER_TYPE_SELL        = 1
  BMT5_ORDER_TYPE_BUY_LIMIT   = 2
  BMT5_ORDER_TYPE_SELL_LIMIT  = 3
  BMT5_ORDER_TYPE_BUY_STOP    = 4
  BMT5_ORDER_TYPE_SELL_STOP   = 5
  BMT5_ORDER_TYPE_BUY_STOP_LIMIT  = 6
  BMT5_ORDER_TYPE_SELL_STOP_LIMIT = 7
  BMT5_ORDER_TYPE_CLOSE_BY    = 8

Total: 9 values
```

---

### Example 4: Find all types related to "Position"

**Command:**
```
> search Position
```

**Output:**
```
Found 8 types matching 'Position':
  • PositionInfo
  • PositionsGetRequest
  • PositionsGetData
  • PositionCloseRequest
  • PositionCloseData
  • PositionModifyRequest
  • PositionModifyData
  • BMT5_ENUM_POSITION_TYPE
```

---

### Example 5: Show all available types

**Command:**
```
> list
```

**Output:**
```
All available types (267):

[A]
  • AccountInfo
  • AccountSummaryData
  • AccountSummaryRequest
  ...

[B]
  • BalanceOperationData
  • BookInfoRequest
  ...

[C]
  • ChartPriceCurrentData
  • ChartPriceCurrentRequest
  ...
```

---

## 🔍 Common Use Cases

### USE CASE 1: "Getting 'field not found' error"

**Problem:** Your code has `position.Ballance` but it doesn't work

**Solution:**
```
> field Balance
```

**Result:** Shows the correct field name and which type has it
```
Field 'Balance' found in:
  • AccountSummaryData
  • AccountInfo
```

**Fix:** Use `accountInfo.Balance`, not `position.Balance`

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
BMT5_ORDER_TYPE_BUY = 0
BMT5_ORDER_TYPE_SELL = 1
BMT5_ORDER_TYPE_BUY_LIMIT = 2
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

**Result:** Shows all 267 available types, grouped alphabetically

---

## 📊 Statistics

- **Total Types:** 267 (all MT5 gRPC protobuf message types)
- **Total Enums:** 68 (with 1400+ enum values)
- **Coverage:** 100% of MT5 gRPC API
- **File size:** 1767 lines, 86KB

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
| `BMT5_ENUM_POSITION_REASON` | Why position opened | Client, Expert, Dealer, Mobile |
| `MqlErrorCode` | MQL error codes | 211 different error codes |
| `MqlErrorTradeCode` | Trading operation errors | REQUOTE, REJECT, MARKET_CLOSED, etc. |
| `SymbolInfoDoubleProperty` | Symbol price properties | BID, ASK, POINT, SWAP_LONG, SWAP_SHORT |
| `SymbolInfoIntegerProperty` | Symbol integer properties | DIGITS, SPREAD, TRADE_MODE, etc. |
| `MRPC_ENUM_TRADE_REQUEST_ACTIONS` | Trade actions | DEAL, PENDING, SLTP, MODIFY, REMOVE |

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **Case-insensitive search** | `search Order` = `search order` |
| **Partial field matching** | `field profit` finds both `Profit` and `TakeProfit` |
| **Type categorization** | Shows `[Request]`, `[Reply]`, `[Type]`, `[Info]` tags |
| **Array indicators** | 📚 icon for repeated/array fields |
| **Protobuf field numbers** | Shows field `#N` for each field |
| **Smart error messages** | Suggests alternatives when type not found |

---

## 🎬 Example Interactive Session

```bash
$ go run main.go inspect

═══════════════════════════════════════════════════════
    MT5 PROTOBUF TYPE INSPECTOR
═══════════════════════════════════════════════════════

> search Order
Found 42 types matching 'Order':
  • OrderInfo
  • OrderSendRequest
  • OrderSendData
  • OrderCheckRequest
  • OrderCalcMarginRequest
  ...

> OrderInfo
═══════════════════════════════════════════════════════
TYPE: OrderInfo
═══════════════════════════════════════════════════════

Fields:
  Ticket: uint64 (#1)
  Type: BMT5_ENUM_ORDER_TYPE (#2)
  State: BMT5_ENUM_ORDER_STATE (#3)
  TypeFilling: BMT5_ENUM_ORDER_TYPE_FILLING (#4)
  ...

> field Magic
Field 'Magic' found in 8 types:
  • OrderInfo
  • PositionInfo
  • DealInfo
  • MqlTradeRequest
  ...

> enum BMT5_ENUM_ORDER_STATE
═══════════════════════════════════════════════════════
ENUM: BMT5_ENUM_ORDER_STATE
═══════════════════════════════════════════════════════

Values:
  BMT5_ORDER_STATE_STARTED       = 0
  BMT5_ORDER_STATE_PLACED        = 1
  BMT5_ORDER_STATE_CANCELED      = 2
  BMT5_ORDER_STATE_PARTIAL       = 3
  BMT5_ORDER_STATE_FILLED        = 4
  BMT5_ORDER_STATE_REJECTED      = 5
  BMT5_ORDER_STATE_EXPIRED       = 6
  BMT5_ORDER_STATE_REQUEST_ADD   = 7
  BMT5_ORDER_STATE_REQUEST_MODIFY = 8
  BMT5_ORDER_STATE_REQUEST_CANCEL = 9

Total: 10 values

> exit
Goodbye!
```

---


## 💻 Implementation Details

The Protobuf Inspector uses Go reflection to:

1. Register all protobuf types at startup
2. Build an in-memory index of types, fields, and enums
3. Provide instant search and lookup
4. Format output with color and structure

**Source file:** `examples/demos/helpers/17_protobuf_inspector.go`

**Enum registration:** All 1400+ enum values across 68 enums are manually registered in the code (Go protobuf doesn't support automatic enum value discovery via reflection).

---

## 🔧 Technical Notes

- **No MT5 connection required** - This is a purely offline tool that inspects type definitions
- **Complete coverage** - All 267 types and 68 enums from the MT5 gRPC API
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
> search position close
Found 2 types:
  • PositionCloseRequest
  • PositionCloseData

> PositionCloseRequest
Fields:
  Ticket: uint64 (#1)
  Deviation: uint64 (#2)
  Comment: string (#3)

> PositionCloseData
Fields:
  ReturnedCode: uint32 (#1)
  OrderTicket: uint64 (#2)

# Now you know:
# - Use PositionCloseRequest with Ticket field
# - You'll get PositionCloseData back
# - Check ReturnedCode == 10009 for success
```

---

## 🆘 Troubleshooting

**Q: Type not found**
```
> MyType
❌ Type not found: MyType
💡 Try: search my
```

**A:** Use search with partial name to find similar types

---

**Q: Too many results**

```
> search data
Found 156 types...
```

**A:** Be more specific in your search query

---

**Q: What's the difference between OrderInfo and OrderSendRequest?**

**A:** Use the inspector:

- `OrderInfo` - Information about an existing order
- `OrderSendRequest` - Request to create a new order

Rule:

- `*Request` - Input for API method
- `*Reply/*Data` - Output from API method
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
