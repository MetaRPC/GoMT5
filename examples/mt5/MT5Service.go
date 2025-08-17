package mt5

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"

	"google.golang.org/protobuf/types/known/timestamppb"
	 //Если в этом файле есть ConnectByProxy/ShowCheckConnect/и т.п. — раскомментируй:
	  //"github.com/google/uuid"
	// _go "git.mtapi.io/root/mrpc-proto.git/mt5/libraries/go"
	// "google.golang.org/grpc/metadata"
)

// MT5Service wraps MT5Account with human-friendly demo/CLI methods.
type MT5Service struct {
	account *MT5Account
}

func NewMT5Service(acc *MT5Account) *MT5Service {
	return &MT5Service{account: acc}
}

// === 📂 Account Info ===

func (s *MT5Service) ShowAccountSummary(ctx context.Context) {
	sum, err := s.account.AccountSummary(ctx)
	if err != nil {
		log.Printf("❌ AccountSummary error: %v", err)
		return
	}
	fmt.Printf("Balance: %.2f | Equity: %.2f | Currency: %s\n",
		sum.GetAccountBalance(), sum.GetAccountEquity(), sum.GetAccountCurrency())
}

// === 📂 Order Operations ===

func (s *MT5Service) ShowOpenedOrders(ctx context.Context) {
	data, err := s.account.OpenedOrders(ctx)
	if err != nil {
		log.Printf("❌ OpenedOrders error: %v", err)
		return
	}
	orders := data.GetOpenedOrders()
	if len(orders) == 0 {
		fmt.Println("📭 No opened orders.")
		return
	}
	for _, o := range orders {
		
		fmt.Printf("[%s] Ticket: %d | Symbol: %s | Volume: %.2f | OpenPrice: %.5f\n",
			o.GetType().String(), o.GetTicket(), o.GetSymbol(), o.GetVolumeInitial(), o.GetPriceOpen())
	}
}

// ShowOrderCalcMargin — calculate the required margin for a potential deal.
func (s *MT5Service) ShowOrderCalcMargin(ctx context.Context, symbol string, orderType pb.ENUM_ORDER_TYPE_TF, volume float64, openPrice float64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	req := &pb.OrderCalcMarginRequest{
		Symbol:    symbol,
		OrderType: orderType, // pb.ENUM_ORDER_TYPE_TF_...
		Volume:    volume,
		OpenPrice: openPrice,
	}
	data, err := s.account.OrderCalcMargin(ctx, req)
	if err != nil {
		log.Printf("❌ OrderCalcMargin error: %v", err)
		return
	}
	fmt.Printf("🧮 Margin required: %.2f\n", data.GetMargin())
}

// ShowOrderCalcProfit — calculate the PnL between the opening and closing price.
func (s *MT5Service) ShowOrderCalcProfit(ctx context.Context, symbol string, orderType pb.ENUM_ORDER_TYPE_TF, volume float64, openPrice, closePrice float64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	req := &pb.OrderCalcProfitRequest{
		OrderType:  orderType,
		Symbol:     symbol,
		Volume:     volume,
		OpenPrice:  openPrice,
		ClosePrice: closePrice,
	}
	data, err := s.account.OrderCalcProfit(ctx, req)
	if err != nil {
		log.Printf("❌ OrderCalcProfit error: %v", err)
		return
	}
	fmt.Printf("💰 Profit calc: %.2f\n", data.GetProfit())
}

func (s *MT5Service) ShowOpenedOrderTickets(ctx context.Context) {
	data, err := s.account.OpenedOrdersTickets(ctx)
	if err != nil {
		log.Printf("❌ OpenedOrdersTickets error: %v", err)
		return
	}
	tix := data.GetOpenedOrdersTickets()
	if len(tix) == 0 {
		fmt.Println("📭 No open order tickets found.")
		return
	}
	fmt.Println("Open Order Tickets:")
	for _, t := range tix {
		fmt.Printf(" - %d\n", t)
	}
}

func (s *MT5Service) ShowOrdersHistory(ctx context.Context) {
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now()
	sortMode := pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE_BMT5_SORT_BY_CLOSE_TIME_DESC

	data, err := s.account.OrdersHistory(ctx, sortMode, &from, &to, nil, nil)
	if err != nil {
		log.Printf("❌ OrdersHistory error: %v", err)
		return
	}
	hist := data.GetHistoryData()
	if len(hist) == 0 {
		fmt.Println("📭 No historical orders found.")
		return
	}
	for _, item := range hist {
		o := item.GetHistoryOrder()
		if o == nil {
			continue
		}
		fmt.Printf("[%s] Ticket: %d | Symbol: %s | Volume: %.2f | Open: %.5f | Close: %.5f | Closed: %s\n",
			o.GetType().String(), o.GetTicket(), o.GetSymbol(), o.GetVolumeInitial(),
			o.GetPriceOpen(), o.GetPriceCurrent(),
			o.GetDoneTime().AsTime().Format("2006-01-02 15:04:05"))
	}
}

func (s *MT5Service) ShowOrderSendExample(ctx context.Context, symbol string) {
	data, err := s.account.OrderSend(
		ctx,
		symbol,
		pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
		0.10,
		nil,         // price (nil for market)
		ptrInt32(5), // slippage
		nil,         // stoploss
		nil,         // takeprofit
		ptrString("Go order test"),
		ptrInt32(123456), // magic
		nil,              // expiration
	)
	if err != nil {
		log.Printf("❌ OrderSend error: %v", err)
		return
	}

	// OrderSendData has: ReturnedCode, Deal, Order, Volume, Price
	deal := data.GetDeal()
	order := data.GetOrder()

	switch {
	case deal != 0:
		// Market execution: got a deal ticket
		fmt.Printf("✅ Market executed! Deal: %d | Price: %.5f | Volume: %.2f | Code: %d\n",
			deal, data.GetPrice(), data.GetVolume(), data.GetReturnedCode())
	case order != 0:
		// Pending order placed: got an order ticket
		fmt.Printf("✅ Pending placed! Order: %d | Price: %.5f | Volume: %.2f | Code: %d\n",
			order, data.GetPrice(), data.GetVolume(), data.GetReturnedCode())
	default:
		// No ticket returned (should be rare) — print raw
		fmt.Printf("⚠️ No ticket in response | Price: %.5f | Volume: %.2f | Code: %d\n",
			data.GetPrice(), data.GetVolume(), data.GetReturnedCode())
	}
}

func (s *MT5Service) ShowOrderSendStopLimitExample(ctx context.Context, symbol string, isBuy bool, trigger, limit float64) {
	data, err := s.account.OrderSendStopLimit(
		ctx,
		symbol,
		isBuy,
		0.10,
		trigger,
		limit,
		ptrInt32(10),
		nil,
		nil,
		ptrString("SLimit from service"),
		ptrInt32(98765),
		timestamppb.New(time.Now().Add(24*time.Hour)),
	)
	if err != nil {
		log.Printf("❌ OrderSendStopLimit error: %v", err)
		return
	}

	if ord := data.GetOrder(); ord != 0 {
		fmt.Printf("✅ STOP_LIMIT placed. Order: %d | Trigger: %.5f | Limit: %.5f | Code: %d\n",
			ord, trigger, limit, data.GetReturnedCode())
		return
	}
	if deal := data.GetDeal(); deal != 0 {
		fmt.Printf("✅ STOP_LIMIT executed immediately. Deal: %d | Price: %.5f | Code: %d\n",
			deal, data.GetPrice(), data.GetReturnedCode())
		return
	}
	fmt.Printf("⚠️ STOP_LIMIT response without ticket | Price: %.5f | Code: %d\n",
		data.GetPrice(), data.GetReturnedCode())
}

func (s *MT5Service) ShowOrderModifyExample(ctx context.Context, ticket uint64) {
	newSL := 1.0500
	newTP := 1.0900

	data, err := s.account.OrderModify(ctx, ticket, nil, &newSL, &newTP, nil)
	if err != nil {
		log.Printf("❌ OrderModify error: %v", err)
		return
	}

	if data != nil {
		fmt.Println("✅ Order successfully modified.")
	} else {
		fmt.Println("⚠️ Order was NOT modified.")
	}
}

// ShowOrderCheck — check the request on the terminal/server side before sending it.
func (s *MT5Service) ShowOrderCheck(
	ctx context.Context,
	action pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS,
	orderType pb.ENUM_ORDER_TYPE_TF,
	symbol string,
	volume float64,
	price float64,
	sl, tp *float64,
	deviation *uint64,
	magic *uint64,
	expiration *timestamppb.Timestamp,
) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}

	req := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action: action,
			ExpertAdvisorMagicNumber: func() uint64 {
				if magic != nil {
					return *magic
				}
				return 0
			}(),
			Symbol: symbol,
			Volume: volume,
			Price:  price,
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation: func() uint64 {
				if deviation != nil {
					return *deviation
				}
				return 0
			}(),
			OrderType:   orderType,
			TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime:    pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC,
			Expiration:  expiration,
			Comment:     "check via Go",
		},
	}

	data, err := s.account.OrderCheck(ctx, req)
	if err != nil {
		log.Printf("❌ OrderCheck error: %v", err)
		return
	}
	res := data.GetMqlTradeCheckResult()
	if res == nil {
		fmt.Println("⚠️ OrderCheck returned empty result")
		return
	}

	fmt.Printf("✅ OrderCheck: retcode=%d | comment=%q | balanceAfter=%.2f | equityAfter=%.2f | profit=%.2f | margin=%.2f | freeMargin=%.2f | marginLevel=%.2f\n",
		res.GetReturnedCode(),
		res.GetComment(),
		res.GetBalanceAfterDeal(),
		res.GetEquityAfterDeal(),
		res.GetProfit(),
		res.GetMargin(),
		res.GetFreeMargin(),
		res.GetMarginLevel(),
	)
}

// ShowOrderCloseExample — closes the market/active ticket order.
func (s *MT5Service) ShowOrderCloseExample(ctx context.Context, ticket uint64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	data, err := s.account.OrderClose(ctx, ticket, nil, nil)
	if err != nil {
		log.Printf("❌ OrderClose error: %v", err)
		return
	}
	
	fmt.Printf("✅ Order closed. CloseMode: %s | Code: %d (%s/%s)\n",
		data.GetCloseMode().String(),
		data.GetReturnedCode(),
		data.GetReturnedStringCode(),
		data.GetReturnedCodeDescription(),
	)
}

// ShowOrderDeleteExample — deletes the pending ticket order.
func (s *MT5Service) ShowOrderDeleteExample(ctx context.Context, ticket uint64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	data, err := s.account.DeleteOrder(ctx, ticket)
	if err != nil {
		log.Printf("❌ DeleteOrder error: %v", err)
		return
	}
	
	fmt.Printf("✅ Pending order deleted. CloseMode: %s | Code: %d (%s/%s)\n",
		data.GetCloseMode().String(),
		data.GetReturnedCode(),
		data.GetReturnedStringCode(),
		data.GetReturnedCodeDescription(),
	)
}

// === 📂 Helpers: Market Orders ===

// BuyMarket — safely buy on the market with checks and defaults.
func (s *MT5Service) BuyMarket(ctx context.Context, symbol string, volume float64, sl, tp *float64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	// 1) Make sure that the symbol is visible
	if err := s.account.EnsureSymbolVisible(ctx, symbol); err != nil {
		log.Printf("❌ EnsureSymbolVisible error: %v", err)
		return
	}

	// 2) Let's take the market price (for BUY we use Ask)
	price, err := getMarketPrice(ctx, s.account, symbol, true /*isBuy*/)
	if err != nil {
		log.Printf("❌ getMarketPrice error: %v", err)
		return
	}

	// 3) Pre-check of the order (OrderCheck) — DEAL BUY
	dev := uint64(10)       // a reasonable default on slippage for the check
	magic := uint64(123456) // default magic
	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
			OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, 
			Symbol:    symbol,
			Volume:    volume,
			Price:     price,
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation:                dev,
			TypeFilling:              pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime:                 pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC,
			ExpertAdvisorMagicNumber: magic,
			Comment:                  "BuyMarket helper",
		},
	}

	chk, err := s.account.OrderCheck(ctx, checkReq)
	if err != nil {
		log.Printf("❌ OrderCheck error: %v", err)
		return
	}
	if r := chk.GetMqlTradeCheckResult(); r != nil {
		fmt.Printf("ℹ️ Check: code=%d, comment=%q, margin=%.2f, free=%.2f\n",
			r.GetReturnedCode(), r.GetComment(), r.GetMargin(), r.GetFreeMargin())
	}

	// 4) We send a BUY market order
	slip := int32(10)
	comment := "BuyMarket"
	magic32 := int32(123456)
	data, err := s.account.OrderSend(
		ctx,
		symbol,
		pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
		volume,
		nil,      // price nil => market
		&slip,    // slippage
		sl,       // stoploss
		tp,       // takeprofit
		&comment, // comment
		&magic32, // magic
		nil,      // expiration
	)
	if err != nil {
		log.Printf("❌ OrderSend(BUY) error: %v", err)
		return
	}
	printOrderSendResult("BUY", data)
}

// SellMarket — safely sell on the market with checks and defaults.
func (s *MT5Service) SellMarket(ctx context.Context, symbol string, volume float64, sl, tp *float64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	// 1) Make sure that the symbol is visible
	if err := s.account.EnsureSymbolVisible(ctx, symbol); err != nil {
		log.Printf("❌ EnsureSymbolVisible error: %v", err)
		return
	}

	// 2) The market price for the SELL — Bid
	price, err := getMarketPrice(ctx, s.account, symbol, false /*isBuy*/)
	if err != nil {
		log.Printf("❌ getMarketPrice error: %v", err)
		return
	}

	// 3) Pre-check of the application (OrderCheck) — DEAL SELL
	dev := uint64(10)
	magic := uint64(123456)
	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
			OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL, 
			Symbol:    symbol,
			Volume:    volume,
			Price:     price,
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation:                dev,
			TypeFilling:              pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime:                 pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC,
			ExpertAdvisorMagicNumber: magic,
			Comment:                  "SellMarket helper",
		},
	}

	chk, err := s.account.OrderCheck(ctx, checkReq)
	if err != nil {
		log.Printf("❌ OrderCheck error: %v", err)
		return
	}
	if r := chk.GetMqlTradeCheckResult(); r != nil {
		fmt.Printf("ℹ️ Check: code=%d, comment=%q, margin=%.2f, free=%.2f\n",
			r.GetReturnedCode(), r.GetComment(), r.GetMargin(), r.GetFreeMargin())
	}

	// 4) We send a SELL market order
	slip := int32(10)
	comment := "SellMarket"
	magic32 := int32(123456)
	data, err := s.account.OrderSend(
		ctx,
		symbol,
		pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL,
		volume,
		nil, // market
		&slip,
		sl,
		tp,
		&comment,
		&magic32,
		nil,
	)
	if err != nil {
		log.Printf("❌ OrderSend(SELL) error: %v", err)
		return
	}
	printOrderSendResult("SELL", data)
}

// === 📂 Helpers: Pending Orders (Limit/Stop/StopLimit) ===

// PlaceBuyLimit — postponement of BUY_LIMIT at the price of price (below the market), with optional SL/TP/expiration.
func (s *MT5Service) PlaceBuyLimit(ctx context.Context, symbol string, volume, price float64, sl, tp *float64, exp *timestamppb.Timestamp) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	if err := s.account.EnsureSymbolVisible(ctx, symbol); err != nil {
		log.Printf("❌ EnsureSymbolVisible: %v", err)
		return
	}

	// Pre-check
	dev := uint64(10)
	magic := uint64(123456)
	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
			OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_LIMIT,
			Symbol:    symbol,
			Volume:    volume,
			Price:     price,
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation:   dev,
			TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
				if exp != nil {
					return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED
				}
				return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
			}(),
			Expiration:               exp,
			ExpertAdvisorMagicNumber: magic,
			Comment:                  "PlaceBuyLimit helper",
		},
	}
	if data, err := s.account.OrderCheck(ctx, checkReq); err != nil {
		log.Printf("❌ OrderCheck(BUY_LIMIT): %v", err)
		return
	} else if r := data.GetMqlTradeCheckResult(); r != nil {
		fmt.Printf("ℹ️ Check BUY_LIMIT: code=%d, comment=%q\n", r.GetReturnedCode(), r.GetComment())
	}

	// Sending
	slip := int32(10)
	comment := "BuyLimit"
	magic32 := int32(123456)
	res, err := s.account.OrderSend(ctx, symbol, pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT, volume, &price, &slip, sl, tp, &comment, &magic32, exp)
	if err != nil {
		log.Printf("❌ OrderSend(BUY_LIMIT): %v", err)
		return
	}
	printOrderSendResult("BUY_LIMIT", res)
}

// PlaceSellLimit — deferral of SELL_LIMIT at the price of price (above the market).
func (s *MT5Service) PlaceSellLimit(ctx context.Context, symbol string, volume, price float64, sl, tp *float64, exp *timestamppb.Timestamp) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	if err := s.account.EnsureSymbolVisible(ctx, symbol); err != nil {
		log.Printf("❌ EnsureSymbolVisible: %v", err)
		return
	}

	dev := uint64(10)
	magic := uint64(123456)
	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
			OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_LIMIT,
			Symbol:    symbol,
			Volume:    volume,
			Price:     price,
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation:   dev,
			TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
				if exp != nil {
					return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED
				}
				return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
			}(),
			Expiration:               exp,
			ExpertAdvisorMagicNumber: magic,
			Comment:                  "PlaceSellLimit helper",
		},
	}
	if data, err := s.account.OrderCheck(ctx, checkReq); err != nil {
		log.Printf("❌ OrderCheck(SELL_LIMIT): %v", err)
		return
	} else if r := data.GetMqlTradeCheckResult(); r != nil {
		fmt.Printf("ℹ️ Check SELL_LIMIT: code=%d, comment=%q\n", r.GetReturnedCode(), r.GetComment())
	}

	slip := int32(10)
	comment := "SellLimit"
	magic32 := int32(123456)
	res, err := s.account.OrderSend(ctx, symbol, pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_LIMIT, volume, &price, &slip, sl, tp, &comment, &magic32, exp)
	if err != nil {
		log.Printf("❌ OrderSend(SELL_LIMIT): %v", err)
		return
	}
	printOrderSendResult("SELL_LIMIT", res)
}

// PlaceBuyStop — BUY_STOP at the trigger price (above the market).
func (s *MT5Service) PlaceBuyStop(ctx context.Context, symbol string, volume, trigger float64, sl, tp *float64, exp *timestamppb.Timestamp) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	if err := s.account.EnsureSymbolVisible(ctx, symbol); err != nil {
		log.Printf("❌ EnsureSymbolVisible: %v", err)
		return
	}

	dev := uint64(10)
	magic := uint64(123456)
	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
			OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_STOP,
			Symbol:    symbol,
			Volume:    volume,
			Price:     trigger,
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation:   dev,
			TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
				if exp != nil {
					return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED
				}
				return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
			}(),
			Expiration:               exp,
			ExpertAdvisorMagicNumber: magic,
			Comment:                  "PlaceBuyStop helper",
		},
	}
	if data, err := s.account.OrderCheck(ctx, checkReq); err != nil {
		log.Printf("❌ OrderCheck(BUY_STOP): %v", err)
		return
	} else if r := data.GetMqlTradeCheckResult(); r != nil {
		fmt.Printf("ℹ️ Check BUY_STOP: code=%d, comment=%q\n", r.GetReturnedCode(), r.GetComment())
	}

	slip := int32(10)
	comment := "BuyStop"
	magic32 := int32(123456)
	res, err := s.account.OrderSend(ctx, symbol, pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP, volume, &trigger, &slip, sl, tp, &comment, &magic32, exp)
	if err != nil {
		log.Printf("❌ OrderSend(BUY_STOP): %v", err)
		return
	}
	printOrderSendResult("BUY_STOP", res)
}

// PlaceSellStop — SELL_STOP at the trigger price (below the market).
func (s *MT5Service) PlaceSellStop(ctx context.Context, symbol string, volume, trigger float64, sl, tp *float64, exp *timestamppb.Timestamp) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	if err := s.account.EnsureSymbolVisible(ctx, symbol); err != nil {
		log.Printf("❌ EnsureSymbolVisible: %v", err)
		return
	}

	dev := uint64(10)
	magic := uint64(123456)
	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
			OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_STOP,
			Symbol:    symbol,
			Volume:    volume,
			Price:     trigger,
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation:   dev,
			TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
				if exp != nil {
					return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED
				}
				return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
			}(),
			Expiration:               exp,
			ExpertAdvisorMagicNumber: magic,
			Comment:                  "PlaceSellStop helper",
		},
	}
	if data, err := s.account.OrderCheck(ctx, checkReq); err != nil {
		log.Printf("❌ OrderCheck(SELL_STOP): %v", err)
		return
	} else if r := data.GetMqlTradeCheckResult(); r != nil {
		fmt.Printf("ℹ️ Check SELL_STOP: code=%d, comment=%q\n", r.GetReturnedCode(), r.GetComment())
	}

	slip := int32(10)
	comment := "SellStop"
	magic32 := int32(123456)
	res, err := s.account.OrderSend(ctx, symbol, pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP, volume, &trigger, &slip, sl, tp, &comment, &magic32, exp)
	if err != nil {
		log.Printf("❌ OrderSend(SELL_STOP): %v", err)
		return
	}
	printOrderSendResult("SELL_STOP", res)
}

// PlaceStopLimit — universal STOP_LIMIT (BUY/SELL) with trigger and limit price.
func (s *MT5Service) PlaceStopLimit(ctx context.Context, symbol string, isBuy bool, volume, trigger, limit float64, sl, tp *float64, exp *timestamppb.Timestamp) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	if err := s.account.EnsureSymbolVisible(ctx, symbol); err != nil {
		log.Printf("❌ EnsureSymbolVisible: %v", err)
		return
	}

	dev := uint64(10)
	magic := uint64(123456)
	orderTypeTF := pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY_STOP_LIMIT
	if !isBuy {
		orderTypeTF = pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_SELL_STOP_LIMIT
	}

	checkReq := &pb.OrderCheckRequest{
		MqlTradeRequest: &pb.MrpcMqlTradeRequest{
			Action:    pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_PENDING,
			OrderType: orderTypeTF,
			Symbol:    symbol,
			Volume:    volume,
			Price:     trigger, // we put the trigger price in the Check
			StopLoss: func() float64 {
				if sl != nil {
					return *sl
				}
				return 0
			}(),
			TakeProfit: func() float64 {
				if tp != nil {
					return *tp
				}
				return 0
			}(),
			Deviation:   dev,
			TypeFilling: pb.MRPC_ENUM_ORDER_TYPE_FILLING_ORDER_FILLING_FOK,
			TypeTime: func() pb.MRPC_ENUM_ORDER_TYPE_TIME {
				if exp != nil {
					return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_SPECIFIED
				}
				return pb.MRPC_ENUM_ORDER_TYPE_TIME_ORDER_TIME_GTC
			}(),
			Expiration:               exp,
			ExpertAdvisorMagicNumber: magic,
			Comment:                  "PlaceStopLimit helper",
		},
	}
	if data, err := s.account.OrderCheck(ctx, checkReq); err != nil {
		log.Printf("❌ OrderCheck(STOP_LIMIT): %v", err)
		return
	} else if r := data.GetMqlTradeCheckResult(); r != nil {
		fmt.Printf("ℹ️ Check STOP_LIMIT: code=%d, comment=%q\n", r.GetReturnedCode(), r.GetComment())
	}

	// Sending via OrderSendEx with StopLimitPrice
	slip := int32(10)
	comment := "StopLimit"
	magic32 := int32(123456)
	res, err := s.account.OrderSendEx(
		ctx,
		symbol,
		func() pb.TMT5_ENUM_ORDER_TYPE {
			if isBuy {
				return pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP_LIMIT
			}
			return pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP_LIMIT
		}(),
		volume,
		&trigger, &slip,
		sl, tp,
		&comment, &magic32,
		exp,
		&limit, // this is the key field of StopLimitPrice
	)
	if err != nil {
		log.Printf("❌ OrderSendEx(STOP_LIMIT): %v", err)
		return
	}
	printOrderSendResult("STOP_LIMIT", res)
}

// getMarketPrice — takes the current price for the symbol: Ask for BUY, Bid for SELL.
func getMarketPrice(ctx context.Context, acc *MT5Account, symbol string, isBuy bool) (float64, error) {
	q, err := acc.Quote(ctx, symbol)
	if err != nil {
		return 0, err
	}
	st := q.GetSymbolTick()
	if st == nil {
		return 0, fmt.Errorf("empty quote payload")
	}
	if isBuy {
		return st.GetAsk(), nil
	}
	return st.GetBid(), nil
}

// printOrderSendResult — a single accurate output of the result of sending an order.
func printOrderSendResult(side string, data *pb.OrderSendData) {
	if data == nil {
		fmt.Printf("✅ %s sent, but empty payload\n", side)
		return
	}
	order := data.GetOrder()
	deal := data.GetDeal()
	price := data.GetPrice()
	switch {
	case order != 0:
		fmt.Printf("✅ %s placed: order=%d @ %.5f\n", side, order, price)
	case deal != 0:
		fmt.Printf("✅ %s executed: deal=%d @ %.5f\n", side, deal, price)
	default:
		fmt.Printf("✅ %s sent @ %.5f\n", side, price)
	}
}

// === 📂 Market Info / Symbol Info ===

func (s *MT5Service) ShowQuote(ctx context.Context, symbol string) {
	q, err := s.account.Quote(ctx, symbol)
	if err != nil {
		log.Printf("❌ Quote error: %v", err)
		return
	}
	if st := q.GetSymbolTick(); st != nil {
		fmt.Printf("✅ %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
			st.GetSymbol(), st.GetBid(), st.GetAsk(), st.GetTime().AsTime().Format("2006-01-02 15:04:05"))
		return
	}
	fmt.Println("⚠️ Empty quote payload.")
}

func (s *MT5Service) ShowQuotesMany(ctx context.Context, symbols []string) {
	qs, err := s.account.QuoteMany(ctx, symbols)
	if err != nil {
		log.Printf("❌ QuoteMany error: %v", err)
		return
	}
	for _, q := range qs {
		if st := q.GetSymbolTick(); st != nil {
			fmt.Printf("📈 %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
				st.GetSymbol(), st.GetBid(), st.GetAsk(), st.GetTime().AsTime().Format("2006-01-02 15:04:05"))
		}
	}
}

func (s *MT5Service) ShowAllSymbols(ctx context.Context) {
	names, err := s.account.ShowAllSymbols(ctx)
	if err != nil {
		log.Printf("❌ ShowAllSymbols error: %v", err)
		return
	}
	if len(names) == 0 {
		fmt.Println("📭 No symbols found.")
		return
	}
	fmt.Println("=== 🧾 All Available Symbols ===")
	for _, name := range names {
		fmt.Printf("• %s\n", name)
	}
}

func (s *MT5Service) ShowSymbolParams(ctx context.Context, symbol string) {
	info, err := s.account.SymbolParams(ctx, symbol)
	if err != nil {
		log.Printf("❌ SymbolParams error: %v", err)
		return
	}
	fmt.Println("📊 Symbol Parameters:")
	fmt.Printf("• Symbol: %s\n", info.GetName())
	fmt.Printf("• Description: %s\n", info.GetSymDescription())
	fmt.Printf("• Digits: %d\n", info.GetDigits())
	fmt.Printf("• Volume Min: %.2f | Max: %.2f | Step: %.2f\n",
		info.GetVolumeMin(), info.GetVolumeMax(), info.GetVolumeStep())
	fmt.Printf("• Trade Mode: %v\n", info.GetTradeMode())
	fmt.Printf("• Currency Base/Profit/Margin: %s / %s / %s\n",
		info.GetCurrencyBase(), info.GetCurrencyProfit(), info.GetCurrencyMargin())
}

func (s *MT5Service) ShowTickValues(ctx context.Context, symbols []string) {
	data, err := s.account.TickValueWithSize(ctx, symbols)
	if err != nil {
		log.Printf("❌ TickValueWithSize error: %v", err)
		return
	}
	
	for _, info := range data.GetSymbolTickSizeInfos() {
		fmt.Printf("💹 %s | TickValue: %.5f | TickSize: %.5f | ContractSize: %.2f\n",
			info.GetName(), info.GetTradeTickValue(), info.GetTradeTickSize(), info.GetTradeContractSize())
	}
}

// === 📂 Positions ===

// ShowPositions — display all open positions.
func (s *MT5Service) ShowPositions(ctx context.Context) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	infos, err := s.account.PositionsGet(ctx)
	if err != nil {
		log.Printf("❌ PositionsGet error: %v", err)
		return
	}
	if len(infos) == 0 {
		fmt.Println("📭 No open positions.")
		return
	}
	for _, p := range infos {
		fmt.Printf("🟢 Pos Ticket: %d | %s | Volume: %.2f | PriceOpen: %.5f | Profit: %.2f\n",
			p.GetTicket(), p.GetSymbol(), p.GetVolume(), p.GetPriceOpen(), p.GetProfit())
	}
}

// ShowHasOpenPosition — check if there is a position for the character.
func (s *MT5Service) ShowHasOpenPosition(ctx context.Context, symbol string) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	ok, err := s.account.HasOpenPosition(ctx, symbol)
	if err != nil {
		log.Printf("❌ HasOpenPosition error: %v", err)
		return
	}
	if ok {
		fmt.Printf("✅ There is an open position for %s\n", symbol)
	} else {
		fmt.Printf("ℹ️ No open position for %s\n", symbol)
	}
}

// ShowPositionClose — close the position by character (we take the first one).
func (s *MT5Service) ShowPositionClose(ctx context.Context, symbol string) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	p, err := s.account.PositionGet(ctx, symbol)
	if err != nil {
		log.Printf("❌ PositionGet error: %v", err)
		return
	}
	if p == nil || p.GetTicket() == 0 {
		fmt.Printf("⚠️ No position found for symbol %s\n", symbol)
		return
	}
	if _, err := s.account.PositionClose(ctx, p); err != nil {
		log.Printf("❌ PositionClose error: %v", err)
		return
	}
	fmt.Printf("✅ Position closed: Ticket %d (%s)\n", p.GetTicket(), p.GetSymbol())
}

// ShowPositionModify — change the SL/TP of the ticket position.
func (s *MT5Service) ShowPositionModify(ctx context.Context, ticket uint64, newSL, newTP *float64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	ok, err := s.account.PositionModify(ctx, ticket, newSL, newTP)
	if err != nil {
		log.Printf("❌ PositionModify error: %v", err)
		return
	}
	if ok {
		fmt.Printf("✅ Position %d modified (SL/TP updated)\n", ticket)
	} else {
		fmt.Printf("⚠️ Position %d was NOT modified\n", ticket)
	}
}

// ShowCloseAllPositions — close all open positions.
func (s *MT5Service) ShowCloseAllPositions(ctx context.Context) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	if err := s.account.CloseAllPositions(ctx); err != nil {
		log.Printf("❌ CloseAllPositions error: %v", err)
		return
	}
	fmt.Println("✅ All positions closed (or none existed).")
}

// === 📂 History ===

// showorderby Ticket — show the historical ticket ORDER.
func (s *MT5Service) ShowOrderByTicket(ctx context.Context, ticket uint64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	o, err := s.account.HistoryOrderByTicket(ctx, ticket)
	if err != nil {
		log.Printf("❌ HistoryOrderByTicket error: %v", err)
		return
	}
	if o == nil {
		fmt.Printf("⚠️ Historical order %d not found\n", ticket)
		return
	}

	vol := o.GetVolumeInitial()
	open := o.GetPriceOpen()
	last := o.GetPriceCurrent()

	fmt.Printf("📜 Order #%d | %s | VolumeInitial: %.2f | PriceOpen: %.5f | LastPrice: %.5f",
		o.GetTicket(), o.GetSymbol(), vol, open, last)

	if ts := o.GetDoneTime(); ts != nil {
		fmt.Printf(" | Done: %s", ts.AsTime().Format("2006-01-02 15:04:05"))
	}
	fmt.Println()
}

// SHOWDEALBY Ticket — show a historical ticket DEAL.
func (s *MT5Service) ShowDealByTicket(ctx context.Context, ticket uint64) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	d, err := s.account.HistoryDealByTicket(ctx, ticket)
	if err != nil {
		log.Printf("❌ HistoryDealByTicket error: %v", err)
		return
	}
	if d == nil {
		fmt.Printf("⚠️ Historical deal %d not found\n", ticket)
		return
	}
	// Safe minimum (price, volume, profit, time)
	fmt.Printf("📜 Deal #%d | %s | Volume: %.2f | Price: %.5f | Profit: %.2f",
		d.GetTicket(), d.GetSymbol(), d.GetVolume(), d.GetPrice(), d.GetProfit())
	if ts := d.GetTime(); ts != nil {
		fmt.Printf(" | Time: %s", ts.AsTime().Format("2006-01-02 15:04:05"))
	}
	fmt.Println()
}

// ShowDealsCount — how many transactions per period (with optional filter by symbol)
func (s *MT5Service) ShowDealsCount(ctx context.Context, from, to time.Time, symbol string) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}
	n, err := s.account.HistoryDealsTotal(ctx, from, to, symbol)
	if err != nil {
		log.Printf("❌ HistoryDealsTotal error: %v", err)
		return
	}
	if symbol != "" {
		fmt.Printf("📊 Deals count for %s in [%s .. %s]: %d\n",
			symbol, from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05"), n)
	} else {
		fmt.Printf("📊 Deals count in [%s .. %s]: %d\n",
			from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05"), n)
	}
}

// === 📂 Streaming / Subscriptions ===

func (s *MT5Service) StreamQuotes(ctx context.Context) {
	symbols := []string{"EURUSD", "GBPUSD"}

	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	tickCh, errCh := s.account.OnSymbolTick(ctx2, symbols)
	fmt.Println("🔄 Streaming ticks...")
	for {
		select {
		case tick, ok := <-tickCh:
			if !ok {
				fmt.Println("✅ Tick stream ended.")
				return
			}
			if st := tick.GetSymbolTick(); st != nil {
				fmt.Printf("[Tick] %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
					st.GetSymbol(), st.GetBid(), st.GetAsk(), st.GetTime().AsTime().Format("2006-01-02 15:04:05"))
			}
		case err := <-errCh:
			log.Printf("❌ Stream error: %v", err)
			return
		case <-time.After(30 * time.Second):
			fmt.Println("⏱️ Timeout reached.")
			return
		}
	}
}

func (s *MT5Service) StreamOpenedOrderProfits(ctx context.Context) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	profitCh, errCh := s.account.OnOpenedOrdersProfit(ctx2, 1000)
	fmt.Println("🔄 Streaming order profits...")
	for {
		select {
		case pkt, ok := <-profitCh:
			if !ok {
				fmt.Println("✅ Profit stream ended.")
				return
			}
			// There are three sets in the data: NewPositions / UpdatedPositions / DeletedPositions
			for _, info := range pkt.GetNewPositions() {
				fmt.Printf("[Profit|NEW] Ticket: %d | Symbol: %s | Profit: %.2f\n",
					info.GetTicket(), info.GetPositionSymbol(), info.GetProfit())
			}
			for _, info := range pkt.GetUpdatedPositions() {
				fmt.Printf("[Profit|UPD] Ticket: %d | Symbol: %s | Profit: %.2f\n",
					info.GetTicket(), info.GetPositionSymbol(), info.GetProfit())
			}
			for _, info := range pkt.GetDeletedPositions() {
				fmt.Printf("[Profit|DEL] Ticket: %d | Symbol: %s | Profit: %.2f\n",
					info.GetTicket(), info.GetPositionSymbol(), info.GetProfit())
			}

		case err := <-errCh:
			log.Printf("❌ Stream error: %v", err)
			return

		case <-time.After(30 * time.Second):
			fmt.Println("⏱️ Timeout reached.")
			return
		}
	}
}

func (s *MT5Service) StreamOpenedOrderTickets(ctx context.Context) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	ticketCh, errCh := s.account.OnOpenedOrdersTickets(ctx2, 1000)
	fmt.Println("🔄 Streaming opened order tickets...")
	for {
		select {
		case pkt, ok := <-ticketCh:
			if !ok {
				fmt.Println("✅ Ticket stream ended.")
				return
			}
			tix := append(pkt.GetPositionTickets(), pkt.GetPendingOrderTickets()...)
			fmt.Printf("[Tickets] %d open tickets: %v\n", len(tix), tix)

		case err := <-errCh:
			log.Printf("❌ Stream error: %v", err)
			return
		case <-time.After(30 * time.Second):
			fmt.Println("⏱️ Timeout reached.")
			return
		}
	}
}

func (s *MT5Service) StreamTradeUpdates(ctx context.Context) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	tradeCh, errCh := s.account.OnTrade(ctx2)
	fmt.Println("🔄 Streaming trade updates...")
	for {
		select {
		case tr, ok := <-tradeCh:
			if !ok {
				fmt.Println("✅ Trade stream ended.")
				return
			}
			if ev := tr.GetEventData(); ev != nil && len(ev.GetNewOrders()) > 0 {
				o := ev.GetNewOrders()[0]
				
				fmt.Printf("[Trade] Ticket: %d | Symbol: %s | Type: %v | Volume: %.2f\n",
					o.GetTicket(), o.GetSymbol(), o.GetOrderType(), o.GetVolumeCurrent())
			}

		case err := <-errCh:
			log.Printf("❌ Stream error: %v", err)
			return

		case <-time.After(30 * time.Second):
			fmt.Println("⏱️ Timeout reached.")
			return
		}
	}
}

// --- Small ptr helpers ---
func ptrInt32(v int32) *int32    { return &v }
func ptrString(v string) *string { return &v }

// --- Helpers for pb types ---

// u64pFromI32 converts *int32 to *uint64 (nil-safe).
func u64pFromI32(v *int32) *uint64 {
	if v == nil {
		return nil
	}
	x := uint64(*v)
	return &x
}

// u64p returns pointer to uint64.
func u64p(x uint64) *uint64 { return &x }
