# Trading Operations — Overview (⚠️ Dangerous)

This section contains **critical trading operations** that directly modify account positions, open or close orders, and place new trades. These methods are **powerful** but must be used with extreme caution to avoid unintended losses.

---

## 📂 Methods in this Section

* [BuyMarket.md](BuyMarket.md)
  Execute an immediate **buy order** at current market price.

* [SellMarket.md](SellMarket.md)
  Execute an immediate **sell order** at current market price.

* [PlaceBuyLimit.md](PlaceBuyLimit.md)
  Place a **buy limit order** below current market price.

* [PlaceSellLimit.md](PlaceSellLimit.md)
  Place a **sell limit order** above current market price.

* [PlaceBuyStop.md](PlaceBuyStop.md)
  Place a **buy stop order** above current market price.

* [PlaceSellStop.md](PlaceSellStop.md)
  Place a **sell stop order** below current market price.

* [PlaceStopLimit.md](PlaceStopLimit.md)
  Place a **stop-limit order**, combining stop and limit conditions.

* [PositionClose.md](PositionClose.md)
  Close an open **position** by ticket or symbol.

* [PositionModify.md](PositionModify.md)
  Modify parameters of an existing position (e.g., SL/TP).

* [CloseAllPositions.md](CloseAllPositions.md)
  Force-close **all open positions** in the account.

* [OrderCloseExample.md](OrderCloseExample.md)
  Example: how to close a single order safely.

* [OrderDeleteExample.md](OrderDeleteExample.md)
  Example: how to delete a pending order.

* [OrderModifyExample.md](OrderModifyExample.md)
  Example: how to modify an order’s price or expiration.

* [ShowOrderSendExample.md](ShowOrderSendExample.md)
  Example of sending a new market/pending order.

* [ShowOrderSendStopLimitExample.md](ShowOrderSendStopLimitExample.md)
  Example of sending a stop-limit order.

* [SetOrderExpiration.md](SetOrderExpiration.md)
  Define an **expiration time** for pending orders.

---

## ⚠️ Important Notes

1. **Always test on demo** before using in live trading.
2. Double-check order parameters (symbol, volume, SL/TP, expiration) before sending.
3. Use [OrderCheck](../Calculations_And_PreliminaryVerification/ShowOrderCheck.md) for validation.
4. Handle errors properly — server rejections may occur due to margin, invalid price, or market state.
5. Be especially careful with **CloseAllPositions**, as it forcefully exits all trades.

---

## ✅ Best Practices

* Use calculation functions (`OrderCalcMargin`, `OrderCalcProfit`) **before placing trades**.
* Implement logging for every order send/modify/delete.
* When modifying orders, always re-validate new parameters.
* Apply rate-limiting — avoid sending too many trade requests per second.

---

## 🎯 Purpose

This block is the **core execution engine** of your trading system. It allows you to:

* Enter the market (market orders).
* Place conditional entries (limit, stop, stop-limit).
* Manage and modify existing positions.
* Exit positions safely (individual or all-at-once).

---

👉 Treat these methods as **live ammunition**: extremely useful but risky if mishandled.
