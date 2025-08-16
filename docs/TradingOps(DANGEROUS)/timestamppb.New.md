# Expiration for Pending Orders (using `timestamppb.New`)

> **Goal:** set, change, or remove expiration for pending orders (limit/stop/stop‑limit) using protobuf timestamps.

---

## TL;DR

```go
// Expire in 24 hours from now
exp := timestamppb.New(time.Now().Add(24 * time.Hour))
```

Use this `exp` when creating or modifying **pending** orders. Market orders ignore expiration.

---

## Where expiration is used

* **OrderCheck**: put `exp` into `MrpcMqlTradeRequest.Expiration` and set `TypeTime = ORDER_TIME_SPECIFIED` (or `*_SPECIFIED_DAY`).
* **OrderSend** family for pendings: pass `exp` argument (e.g., `OrderSend`, `OrderSendStopLimit`).
* **OrderModify**: the last argument is `expiration` — set a new one or `nil` to remove.

> If you pass `exp` but keep `TypeTime = ORDER_TIME_GTC`, many servers will **ignore** it or **reject** the request.

---

## Step‑by‑step: placing a Buy Limit that expires in 24h

```go
// 1) Pick an absolute time
expiry := time.Now().Add(24 * time.Hour)

// 2) Convert to protobuf timestamp
exp := timestamppb.New(expiry)

// 3) (optional) Dry‑run: OrderCheck with TypeTime and Expiration
chkReq := &pb.OrderCheckRequest{
    MqlTradeRequest: &pb.MrpcMqlTradeRequest{
        Action:     pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
        OrderType:  pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_LIMIT,
        Symbol:     "EURUSD",
        Volume:     0.10,
        Price:      1.07500,
        Deviation:  10,
        TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
        TypeTime:    pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED,
        Expiration:  exp,
        Comment:     "PlaceBuyLimit with exp",
    },
}
_, _ = svc.account.OrderCheck(ctx, chkReq)

// 4) Send the pending order with expiration
slip := int32(10)
price := 1.07500
comment := "BuyLimit"
magic := int32(123456)
res, err := svc.account.OrderSend(
    ctx,
    "EURUSD",
    pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT,
    0.10,
    &price, &slip,
    nil, nil,
    &comment, &magic,
    exp,
)
```

---

## Common recipes

### 1) Expire **today at 23:59:59** (server timezone)

```go
loc, _ := time.LoadLocation("Europe/London") // replace with your server TZ if known
now := time.Now().In(loc)
endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
exp := timestamppb.New(endOfDay)
// Use TypeTime ORDER_TIME_SPECIFIED (absolute time) or ORDER_TIME_SPECIFIED_DAY if supported/required
```

### 2) Expire **in N minutes/hours**

```go
func ExpIn(d time.Duration) *timestamppb.Timestamp {
    return timestamppb.New(time.Now().Add(d))
}
exp := ExpIn(90 * time.Minute)
```

### 3) Set / change / remove expiration on an **existing** pending order

```go
// Set new expiration to +6h
exp := timestamppb.New(time.Now().Add(6 * time.Hour))
_, err := svc.account.OrderModify(ctx, ticket, nil /*newPrice*/, nil /*SL*/, nil /*TP*/, exp)

// Remove expiration (back to GTC)
_, err = svc.account.OrderModify(ctx, ticket, nil, nil, nil, nil)
```

> Signature matches helpers shown earlier: the **last** arg is the expiration timestamp.

---

## `ORDER_TIME` modes you’ll touch

| Enum value                 | Meaning                                        |
| -------------------------- | ---------------------------------------------- |
| `ORDER_TIME_GTC`           | Good‑Till‑Cancel. No expiration (`exp=nil`).   |
| `ORDER_TIME_SPECIFIED`     | Good‑Till **absolute** timestamp (`exp!=nil`). |
| `ORDER_TIME_SPECIFIED_DAY` | Good‑Till end of the **specified day**.        |

**Rules:**

* With `*_SPECIFIED*` → you **must** provide `Expiration`.
* With `GTC` → pass `nil` (or clear via modify) to avoid confusion.

---

## FAQ & Pitfalls

* **Local vs server timezone?** `timestamppb.Timestamp` encodes an absolute instant. Timezone only matters when you *choose* the moment (e.g., "end of day"). Pick the proper `Location` when building `time.Time`.
* **Past times** → broker will reject. Always compare `expiry.After(time.Now())`.
* **Market orders** ignore expiration. It’s only for pending types.
* **Stop/Stop‑Limit**: same pattern — set `TypeTime=*SPECIFIED*` + `Expiration` in both `OrderCheck` and `OrderSend*`.
* **Server differences**: some servers silently downgrade to GTC if `TypeTime` doesn’t match `Expiration`. Always inspect `ReturnedCode`/comments.

---

## Mini‑cheatsheet

```go
// 24h from now
exp := timestamppb.New(time.Now().Add(24*time.Hour))

// End of today in local TZ
end := time.Date(y, m, d, 23,59,59,0, time.Local)
exp := timestamppb.New(end)

// No expiration (GTC)
var exp *timestamppb.Timestamp = nil
```
