# ⏳ Expiration (timestamppb) — How To Use

> **Purpose:** set expiration time for pending orders in MT5 (Buy Limit, Sell Limit, Buy Stop, etc.) 🕒

---

### Code Example ⚡

```go
// expire in 24h
exp := timestamppb.New(time.Now().Add(24 * time.Hour))

// use in pending order placement
svc.PlaceBuyLimit(ctx, "EURUSD", 0.10, 1.07500, nil, nil, exp)
```

---

### Step by Step 📝

1. `time.Now()` → current system time ⌚
2. `.Add(24 * time.Hour)` → add 24 hours ⏰
3. `timestamppb.New(...)` → convert Go `time.Time` → Protobuf `Timestamp` 📦
4. Pass `exp` to order methods → broker knows when to auto-cancel order ❌

---

### Variations 🎛️

* Expire in **N minutes**:

  ```go
  exp := timestamppb.New(time.Now().Add(30 * time.Minute))
  ```
* Expire **end of day**:

  ```go
  end := time.Date(y, m, d, 23, 59, 59, 0, time.Local)
  exp := timestamppb.New(end)
  ```
* Expire in **N days**:

  ```go
  exp := timestamppb.New(time.Now().AddDate(0,0,3)) // +3 days
  ```

---

### Usage Scenarios 🎯

* `OrderCheck` 🧪 → dry-run with expiration.
* `OrderSend` 🚀 → send pending order with lifetime.
* `OrderModify` 🔧 → update existing order to set/change expiration.

---

### ENUM: ORDER\_TIME ⌨️

| Value                      | Meaning                   |
| -------------------------- | ------------------------- |
| `ORDER_TIME_GTC`           | Good-Till-Cancelled ♾️    |
| `ORDER_TIME_DAY`           | Valid only for the day 📅 |
| `ORDER_TIME_SPECIFIED`     | Expire at exact time ⏳    |
| `ORDER_TIME_SPECIFIED_DAY` | Expire at end of day 🏁   |

When you pass `exp`, API sets `ORDER_TIME_SPECIFIED`. ✅

---

### Common Pitfalls ⚠️

* Don’t forget: **market orders ignore expiration** ❌
* Server clock ≠ local clock 🖥️ vs. 🏦
* Timezone: MT servers usually run **EET (UTC+2/3)** 🌍

---

### Example with full context 📚

```go
exp := timestamppb.New(time.Now().Add(24 * time.Hour))

res, err := svc.account.OrderSend(
    ctx,
    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
    pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_LIMIT,
    "EURUSD",
    0.10,
    1.07500,
    nil, nil, nil, nil, nil,
    exp,
)
```

---

✅ In short: `timestamppb.New(...)` is your helper to say *“this order lives until X”*. Very handy for **pending orders** 🎯
