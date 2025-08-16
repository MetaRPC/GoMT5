# 📘 Folder Overviews for Documentation

Below are **Overview\.md** drafts for each folder in your docs. Each contains:

* A short intro about the folder
* Links to each method with a one‑liner description
* Navigation back to index

---

## 1. `Calculations_And_PreliminaryVerification/Overview.md`

```markdown
# Calculations and Preliminary Verification — Overview

Methods for margin calculation and trade request pre‑checks.

## Available Methods
- [Order Calc Margin](OrderCalcMargin.md)  
  Estimate required margin before sending an order.

- [Order Check](ShowOrderCheck.md)  
  Validate a trade request without executing it.

- [Verification Log Result](VerificationLogResult.md)  
  Inspect detailed results of trade checks.

## Navigation
- [🏠 Home](../index.md)
- [⬅ Back to Index](../index.md)
```

---

## 2. `History_And_SimpleStatistics/Overview.md`

```markdown
# History and Simple Statistics — Overview

Methods to access trade history and perform simple statistical queries.

## Available Methods
- [History And Stats Overview](HistoryAndStats_Overview.md)  
  General description of history/statistics functions.

- [Deal By Ticket](DealByTicket.md)  
  Retrieve details of a specific historical deal.

## Navigation
- [🏠 Home](../index.md)
- [⬅ Back to Index](../index.md)
```

---

## 3. `Opened_State_Snapshot/Overview.md`

```markdown
# Opened State Snapshot — Overview

Functions to inspect current open positions and orders.

## Available Methods
- [Has Open Position](HasOpenPosition.md)  
  Check if account currently has an open trade.

- [Opened Orders](OpenedOrders.md)  
  Get a list of currently active orders.

- [Opened Orders Tickets](OpenedOrdersTickets.md)  
  Retrieve ticket numbers of all open orders.

## Navigation
- [🏠 Home](../index.md)
- [⬅ Back to Index](../index.md)
```

---

## 4. `QuickAccount_MarketInfo/Overview.md`

```markdown
# Quick Account and Market Info — Overview

Helper methods for quick access to account and market details.

## Available Methods
- [Account Summary](AccountSummary.md)  
  Get main account metrics in one call.

- [Set Order Expiration](SetOrderExpiration.md)  
  Attach expiration time to pending orders.

- [Symbol Info Tick](SymbolInfoTick.md)  
  Retrieve latest tick information for a symbol.

## Navigation
- [🏠 Home](../index.md)
- [⬅ Back to Index](../index.md)
```

---

## 5. `TradingOps(DANGEROUS)/Overview.md`

```markdown
# Trading Operations (Dangerous) — Overview

⚠️ **Use with caution.** These methods perform direct trading actions.

## Available Methods
- [Send Order](SendOrder.md)  
  Place a new trade order.

- [Delete Order](DeleteOrder.md)  
  Cancel a pending order.

- [Modify Order](OrderModify.md)  
  Change parameters of an existing order.

## Navigation
- [🏠 Home](../index.md)
- [⬅ Back to Index](../index.md)
```

---

✅ Each folder now has a ready‑to‑drop **Overview\.md** with method links and one‑liner explanations.
