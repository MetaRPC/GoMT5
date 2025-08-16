# Calculations & Preliminary Verification — Overview

This section groups together methods for **pre-trade calculations** and **order verification**.
They allow you to check whether an order is valid, and to estimate margin and profit before actually placing it.

---

## 📂 Methods in this Section

* [OrderCalcMargin.md](OrderCalcMargin.md)
  Estimate required **margin** for a potential trade before opening it.

* [OrderCalcProfit.md](OrderCalcProfit.md)
  Calculate expected **profit or loss** for a trade, given entry/exit prices.

* [ShowOrderCheck.md](ShowOrderCheck.md)
  Perform a full **order validation** (volume, price, stop levels, free margin).

---

## ⚡ Example Workflow

```go
// Example: check if order is valid before sending

// 1. Calculate margin requirement
margin, _ := svc.OrderCalcMargin(ctx, symbol, lotSize, orderType, price)

// 2. Estimate profit for target exit
profit, _ := svc.OrderCalcProfit(ctx, symbol, lotSize, orderType, entry, exit)

// 3. Run broker's internal order check
checkResult, _ := svc.ShowOrderCheck(ctx, orderRequest)
```

---

## ✅ Best Practices

1. **Always check before sending** — avoid rejected trades and wasted time.
2. Use `OrderCalcMargin` to prevent margin-call surprises.
3. Use `OrderCalcProfit` for quick risk/reward estimation.
4. Combine all three for robust pre-trade validation.

---

## 🎯 Purpose

The methods in this block allow you to:

* Simulate trades without execution.
* Estimate margin and risk in advance.
* Validate orders before sending them to the broker.
* Build custom risk-management logic.

---

👉 Use this overview as a **map**, and jump into each `.md` file for full method details.
