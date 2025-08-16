# Listing Opened Orders

> **Request:** fetch snapshot of currently opened (active) orders and positions in one call.

---

### Code Example

```go
// High-level (prints selected fields per order):
svc.ShowOpenedOrders(ctx)

// Low-level (full access):
req := &pb.OpenedOrdersRequest{ /* optional: InputSortMode */ }
// If you don't need custom sort, just pass zero-value request.

// wrapper call in your service does this internally:
data, err := svc.account.OpenedOrders(ctx)
if err != nil {
    log.Printf("❌ OpenedOrders error: %v", err)
    return
}
orders := data.GetOpenedOrders()
if len(orders) == 0 {
    fmt.Println("📭 No opened orders.")
}
for _, o := range orders {
    fmt.Printf("[%s] Ticket:%d Symbol:%s Vol:%.2f Open:%.5f Curr:%.5f\n",
        o.GetType().String(), o.GetTicket(), o.GetSymbol(),
        o.GetVolumeInitial(), o.GetPriceOpen(), o.GetPriceCurrent())
}
```

---

### Method Signature

```go
func (s *MT5Service) ShowOpenedOrders(ctx context.Context)
```

---

## 🔽 Input

| Parameter | Type              | Required | Description                                |
| --------- | ----------------- | -------- | ------------------------------------------ |
| `ctx`     | `context.Context` | yes      | Controls timeout/cancellation for the RPC. |

> Low-level API also accepts `OpenedOrdersRequest` with optional `InputSortMode`.

**Enum: `BMT5_ENUM_OPENED_ORDER_SORT_TYPE`**

| Value | Name                                             | Meaning                   |
| ----: | ------------------------------------------------ | ------------------------- |
|     0 | `BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC`        | By open time (ascending)  |
|     1 | `BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_DESC`       | By open time (descending) |
|     2 | `BMT5_OPENED_ORDER_SORT_BY_ORDER_TICKET_ID_ASC`  | By ticket (ascending)     |
|     3 | `BMT5_OPENED_ORDER_SORT_BY_ORDER_TICKET_ID_DESC` | By ticket (descending)    |

---

## ⬆️ Output

RPC returns **`OpenedOrdersData`** with two arrays:

| Field           | Type                 | Description                                 |
| --------------- | -------------------- | ------------------------------------------- |
| `OpenedOrders`  | `[]*OpenedOrderInfo` | Active orders (market & pending).           |
| `PositionInfos` | `[]*PositionInfo`    | Open positions snapshot (see separate doc). |

### `OpenedOrderInfo` (key fields)

| Field            | Type                           | Description                                   |
| ---------------- | ------------------------------ | --------------------------------------------- |
| `Ticket`         | `uint64`                       | Order ticket ID.                              |
| `Symbol`         | `string`                       | Symbol name.                                  |
| `Type`           | `BMT5_ENUM_ORDER_TYPE`         | Order type (see enum below).                  |
| `State`          | `BMT5_ENUM_ORDER_STATE`        | Current order state.                          |
| `VolumeInitial`  | `double`                       | Initial lot size.                             |
| `VolumeCurrent`  | `double`                       | Remaining lot (for partial fills).            |
| `PriceOpen`      | `double`                       | Open/placed price.                            |
| `PriceCurrent`   | `double`                       | Current price.                                |
| `StopLoss`       | `double`                       | SL price.                                     |
| `TakeProfit`     | `double`                       | TP price.                                     |
| `StopLimit`      | `double`                       | StopLimit price (for stop-limit orders).      |
| `TimeSetup`      | `google.protobuf.Timestamp`    | When order was created.                       |
| `TimeExpiration` | `google.protobuf.Timestamp`    | Expiration time (if set).                     |
| `TimeDone`       | `google.protobuf.Timestamp`    | Completion time (if filled/canceled/expired). |
| `TypeFilling`    | `BMT5_ENUM_ORDER_TYPE_FILLING` | Execution policy (FOK/IOC).                   |
| `TypeTime`       | `BMT5_ENUM_ORDER_TYPE_TIME`    | Time-in-force.                                |
| `PositionId`     | `int64`                        | Linked position id (if any).                  |
| `PositionById`   | `int64`                        | Opposite position id for close-by.            |
| `MagicNumber`    | `int64`                        | EA magic number.                              |
| `Reason`         | `int32`                        | Server-side reason code (broker-specific).    |
| `Comment`        | `string`                       | User/bot comment.                             |
| `AccountLogin`   | `int64`                        | Account number.                               |

### Enum: `BMT5_ENUM_ORDER_TYPE`

| Value | Name                              | Meaning                    |
| ----: | --------------------------------- | -------------------------- |
|     0 | `BMT5_ORDER_TYPE_BUY`             | Market Buy                 |
|     1 | `BMT5_ORDER_TYPE_SELL`            | Market Sell                |
|     2 | `BMT5_ORDER_TYPE_BUY_LIMIT`       | Pending Buy Limit          |
|     3 | `BMT5_ORDER_TYPE_SELL_LIMIT`      | Pending Sell Limit         |
|     4 | `BMT5_ORDER_TYPE_BUY_STOP`        | Pending Buy Stop           |
|     5 | `BMT5_ORDER_TYPE_SELL_STOP`       | Pending Sell Stop          |
|     6 | `BMT5_ORDER_TYPE_BUY_STOP_LIMIT`  | Pending Buy Stop-Limit     |
|     7 | `BMT5_ORDER_TYPE_SELL_STOP_LIMIT` | Pending Sell Stop-Limit    |
|     8 | `BMT5_ORDER_TYPE_CLOSE_BY`        | Close by opposite position |

### Enum: `BMT5_ENUM_ORDER_STATE`

| Value | Name                        | Meaning                   |
| ----: | --------------------------- | ------------------------- |
|     0 | `BMT5_ORDER_STATE_STARTED`  | Checked, not yet accepted |
|     1 | `BMT5_ORDER_STATE_PLACED`   | Accepted                  |
|     2 | `BMT5_ORDER_STATE_CANCELED` | Canceled by client        |
|     3 | `BMT5_ORDER_STATE_PARTIAL`  | Partially filled          |
|     4 | `BMT5_ORDER_STATE_FILLED`   | Fully executed            |
|     5 | `BMT5_ORDER_STATE_REJECTED` | Rejected                  |
|     6 | `BMT5_ORDER_STATE_EXPIRED`  | Expired                   |

### Enum: `BMT5_ENUM_ORDER_TYPE_FILLING`

| Value | Name                     | Meaning             |
| ----: | ------------------------ | ------------------- |
|     0 | `BMT5_ORDER_FILLING_FOK` | Fill-or-kill        |
|     1 | `BMT5_ORDER_FILLING_IOC` | Immediate-or-cancel |

### Enum: `BMT5_ENUM_ORDER_TYPE_TIME`

| Value | Name                            | Meaning                        |
| ----: | ------------------------------- | ------------------------------ |
|     0 | `BMT5_ORDER_TIME_GTC`           | Good-till-cancel               |
|     1 | `BMT5_ORDER_TIME_DAY`           | Good-for-day                   |
|     2 | `BMT5_ORDER_TIME_SPECIFIED`     | Good-till specified timestamp  |
|     3 | `BMT5_ORDER_TIME_SPECIFIED_DAY` | Good-till end of specified day |

---

## 🎯 Purpose

* One-shot snapshot of active orders + positions for dashboards and risk views.
* Basis for UIs showing order state, SL/TP, and remaining volume.
* Input for automation that reconciles positions with orders.

---

## 🧩 Notes & Tips

* Empty `OpenedOrders` is normal when nothing is placed.
* Use `OpenedOrdersTickets` for a light-weight list of IDs when you don’t need full payloads.
* `Reason` is broker-dependent; don’t rely on it for logic unless you map known values.
* Timestamps are server time; convert via `AsTime()` and display in user TZ if needed.
