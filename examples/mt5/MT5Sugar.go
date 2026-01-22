/*══════════════════════════════════════════════════════════════════════════════
 FILE: MT5Sugar.go - HIGH-LEVEL SUGAR API FOR MT5 TRADING

 PURPOSE:
   MT5Sugar is the HIGHEST-LEVEL API layer in GoMT5, designed for maximum
   simplicity and productivity. It provides ultra-simple one-liner methods
   for all common MT5 operations with automatic defaults and smart error handling.

 🎯 WHO SHOULD USE THIS API:
   • Traders wanting quick scripts and trading bots
   • Beginners learning MT5 API programming
   • Rapid prototyping and testing
   • Simple monitoring and automation tools
   • Anyone preferring simplicity over fine-grained control

 📚 COMPLETE METHOD INDEX (61 METHODS IN 13 CATEGORIES):

   ┌─────────────────────────────────────────────────────────────┐
   │  INITIALIZATION & HELPERS (3 methods)                       │
   ├─────────────────────────────────────────────────────────────┤
   │  • NewMT5Sugar()    - Create Sugar instance                 │
   │  • GetService()     - Access underlying Service layer       │
   │  • GetAccount()     - Access underlying Account layer       │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  1. CONNECTION METHODS (3 methods)                          │
   ├─────────────────────────────────────────────────────────────┤
   │  • QuickConnect()   - Connect via cluster name (RECOMMENDED)│
   │  • IsConnected()    - Check connection status               │
   │  • Ping()           - Verify connection health              │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  2. QUICK BALANCE METHODS (6 methods)                       │
   ├─────────────────────────────────────────────────────────────┤
   │  • GetBalance()     - Account balance                       │
   │  • GetEquity()      - Current equity                        │
   │  • GetMargin()      - Used margin                           │
   │  • GetFreeMargin()  - Available margin                      │
   │  • GetMarginLevel() - Margin level percentage               │
   │  • GetProfit()      - Floating profit/loss                  │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  3. PRICES & QUOTES METHODS (5 methods + 1 struct)          │
   ├─────────────────────────────────────────────────────────────┤
   │  • GetBid()         - Current BID price                     │
   │  • GetAsk()         - Current ASK price                     │
   │  • GetSpread()      - Spread in points                      │
   │  • GetPriceInfo()   - Complete price information            │
   │  • WaitForPrice()   - Wait for price update with timeout    │
   │  • PriceInfo        - Price information structure           │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  4. SIMPLE TRADING METHODS (6 methods)                      │
   ├─────────────────────────────────────────────────────────────┤
   │  • BuyMarket()      - Open BUY position at market           │
   │  • SellMarket()     - Open SELL position at market          │
   │  • BuyLimit()       - Place BUY LIMIT pending order         │
   │  • SellLimit()      - Place SELL LIMIT pending order        │
   │  • BuyStop()        - Place BUY STOP pending order          │
   │  • SellStop()       - Place SELL STOP pending order         │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  5. TRADING WITH SL/TP (4 methods)                          │
   ├─────────────────────────────────────────────────────────────┤
   │  • BuyMarketWithSLTP()  - BUY with Stop Loss & Take Profit  │
   │  • SellMarketWithSLTP() - SELL with Stop Loss & Take Profit │
   │  • BuyLimitWithSLTP()   - BUY LIMIT with SL/TP              │
   │  • SellLimitWithSLTP()  - SELL LIMIT with SL/TP             │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  6. POSITION MANAGEMENT (7 methods)                         │
   ├─────────────────────────────────────────────────────────────┤
   │  • ClosePosition()        - Close full position             │
   │  • ClosePositionPartial() - Close partial volume            │
   │  • CloseAllPositions()    - Close all open positions        │
   │  • CloseAllBySymbol()     - Close all for specific symbol   │
   │  • ModifyPositionSL()     - Change Stop Loss                │
   │  • ModifyPositionTP()     - Change Take Profit              │
   │  • ModifyPositionSLTP()   - Change both SL and TP           │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  7. POSITION INFORMATION (7 methods)                        │
   ├─────────────────────────────────────────────────────────────┤
   │  • GetOpenPositions()    - Get all open positions           │
   │  • GetPositionByTicket() - Find position by ticket number   │
   │  • GetPositionsBySymbol()- Get positions for symbol         │
   │  • HasOpenPosition()     - Check if positions exist         │
   │  • CountOpenPositions()  - Count total open positions       │
   │  • GetTotalProfit()      - Total floating P/L               │
   │  • GetProfitBySymbol()   - Profit for specific symbol       │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  8. HISTORY & PROFIT ANALYSIS (9 methods)                   │
   ├─────────────────────────────────────────────────────────────┤
   │  • GetDealsToday()       - All deals from today             │
   │  • GetDealsYesterday()   - All deals from yesterday         │
   │  • GetDealsThisWeek()    - All deals from this week         │
   │  • GetDealsThisMonth()   - All deals from this month        │
   │  • GetDealsDateRange()   - Deals within custom date range   │
   │  • GetProfitToday()      - Total profit from today          │
   │  • GetProfitThisWeek()   - Total profit from this week      │
   │  • GetProfitThisMonth()  - Total profit from this month     │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  9. SYMBOL INFORMATION METHODS (5 methods + 1 struct)       │
   ├─────────────────────────────────────────────────────────────┤
   │  • GetSymbolInfo()       - Complete symbol information      │
   │  • GetAllSymbols()       - List all available symbols       │
   │  • IsSymbolAvailable()   - Check if symbol is tradeable     │
   │  • GetMinStopLevel()     - Minimum stop level for symbol    │
   │  • GetSymbolDigits()     - Symbol decimal precision         │
   │  • SymbolInfo            - Symbol information structure     │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  10. RISK MANAGEMENT METHODS (4 methods)  IMPORTANT         │
   ├─────────────────────────────────────────────────────────────┤
   │  • CalculatePositionSize()  - Auto-size based on risk %     │
   │  • GetMaxLotSize()          - Maximum tradeable volume      │
   │  • CanOpenPosition()        - Validate before trading       │
   │  • CalculateRequiredMargin()- Margin needed for position    │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  11. TRADING HELPERS (3 methods)                            │
   ├─────────────────────────────────────────────────────────────┤
   │  • CalculateSLTP()       - Convert pips to price levels     │
   │  • BuyMarketWithPips()   - BUY with SL/TP in pips           │
   │  • SellMarketWithPips()  - SELL with SL/TP in pips          │
   └─────────────────────────────────────────────────────────────┘

   ┌─────────────────────────────────────────────────────────────┐
   │  12. ACCOUNT INFORMATION (2 methods + 2 structs)            │
   ├─────────────────────────────────────────────────────────────┤
   │  • GetAccountInfo()      - Complete account details         │
   │  • GetDailyStats()       - Daily trading statistics         │
   │  • AccountInfo           - Account information structure    │
   │  • DailyStats            - Daily statistics structure       │
   └─────────────────────────────────────────────────────────────┘

 ⚠️  IMPORTANT NOTES:
   • All methods have built-in timeouts (3-30 seconds depending on operation)
   • Market orders timeout: 10 seconds
   • Balance/Price queries timeout: 3 seconds
   • History queries timeout: 5-30 seconds (5s for day/week, 30s for month/custom range)
   • Risk management timeout: 10 seconds (CalculatePositionSize with margin checks)
   • Symbol list queries timeout: 15 seconds (GetAllSymbols - many symbols)
   • Bulk operations timeout: 30 seconds
   • Use GetService() or GetAccount() if you need more control

      SEE ALSO:
   • examples/demos/sugar/06_sugar_basics.go     - Connection & Balance demo
   • examples/demos/sugar/07_sugar_trading.go    - Trading operations demo
   • examples/demos/sugar/08_sugar_positions.go  - Position management demo
   • examples/demos/sugar/09_sugar_history.go    - History & profit analysis demo
   • examples/demos/sugar/10_sugar_advanced.go   - Risk management & advanced features

══════════════════════════════════════════════════════════════════════════════*/

package mt5

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	helpers "github.com/MetaRPC/GoMT5/package/Helpers"
	"github.com/google/uuid"
)

// ══════════════════════════════════════════════════════════════════════════════
// TYPE DEFINITIONS
// ══════════════════════════════════════════════════════════════════════════════

// MT5Sugar is the high-level API wrapper providing ultra-simple one-liner methods
// for all common MT5 operations. It automatically handles contexts, timeouts, and
// provides smart defaults for all parameters.
type MT5Sugar struct {
	service  *MT5Service
	ctx      context.Context
	user     uint64
	password string
}

// PriceInfo holds complete current price information for a trading symbol.
// This structure provides all essential price data in one convenient package.
//
// FIELDS:
//   Symbol     - Trading symbol name (e.g., "EURUSD")
//   Bid        - Current BID price for selling
//   Ask        - Current ASK price for buying
//   SpreadPips - Spread in points (not price units)
//   Time       - Server time of the last tick
type PriceInfo struct {
	Symbol     string
	Bid        float64
	Ask        float64
	SpreadPips float64
	Time       time.Time
}

// ══════════════════════════════════════════════════════════════════════════════
// INITIALIZATION & HELPERS
// ══════════════════════════════════════════════════════════════════════════════

// NewMT5Sugar creates a new MT5Sugar instance with the provided credentials.
// This is the main entry point for using the Sugar API. It initializes both
// the low-level Account and mid-level Service layers automatically.
//
// PARAMETERS:
//   user       - MT5 account login number
//   password   - MT5 account password
//   grpcServer - gRPC server address (host:port, e.g., "mt5.server.com:443")
//
// RETURNS:
//   *MT5Sugar instance ready for connection, or error if initialization fails
func NewMT5Sugar(user uint64, password string, grpcServer string) (*MT5Sugar, error) {
	account, err := helpers.NewMT5Account(user, password, grpcServer, uuid.New())
	if err != nil {
		return nil, fmt.Errorf("failed to create MT5Account: %w", err)
	}

	service := NewMT5Service(account)

	return &MT5Sugar{
		service:  service,
		ctx:      context.Background(),
		user:     user,
		password: password,
	}, nil
}

// GetService returns the underlying MT5Service instance for operations that
// require more control than Sugar API provides. Use this when you need access
// to mid-level API features like custom timeouts or advanced parameters.
//
// RETURNS:
//   *MT5Service instance used by this Sugar wrapper
func (s *MT5Sugar) GetService() *MT5Service {
	return s.service
}

// GetAccount returns the underlying MT5Account instance for low-level operations.
// Use this when you need direct access to protobuf structures or maximum control
// over request parameters. Required for closing the gRPC connection.
//
// RETURNS:
//   *MT5Account instance used by the underlying Service layer
func (s *MT5Sugar) GetAccount() *helpers.MT5Account {
	return s.service.account
}

// ══════════════════════════════════════════════════════════════════════════════
// #region CONNECTION METHODS
// ══════════════════════════════════════════════════════════════════════════════

// QuickConnect connects to MT5 terminal using cluster name (RECOMMENDED).
// This is the easiest connection method - just provide your broker's cluster name.
// Automatically sets up EURUSD as base chart symbol and uses 30-second timeout.
//
// PARAMETERS:
//   clusterName - MT5 cluster identifier (e.g., "FxPro-MT5 Demo", "ICMarkets-Live02")
//
// RETURNS:
//   Error if connection fails, nil on success
func (s *MT5Sugar) QuickConnect(clusterName string) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	baseSymbol := "EURUSD"
	req := &pb.ConnectExRequest{
		User:            s.user,
		Password:        s.password,
		MtClusterName:   clusterName,
		BaseChartSymbol: &baseSymbol,
	}

	_, err := s.GetAccount().ConnectEx(ctx, req)
	return err
}

// IsConnected checks if the connection to MT5 terminal is alive.
// This is a quick boolean check using 3-second timeout. Returns false if
// connection is dead or health check times out. Does not return errors.
//
// RETURNS:
//   true if connected and alive, false otherwise
func (s *MT5Sugar) IsConnected() bool {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	data, err := s.GetAccount().CheckConnect(ctx, &pb.CheckConnectRequest{})
	if err != nil {
		return false
	}

	return data.HealthCheck != nil && data.HealthCheck.IsAlive
}

// Ping verifies the connection health to MT5 terminal with detailed error reporting.
// Unlike IsConnected(), this method returns an error explaining why connection failed.
// Uses 3-second timeout. Useful for debugging connection issues.
//
// RETURNS:
//   Error with details if ping fails or connection is dead, nil if healthy
func (s *MT5Sugar) Ping() error {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	data, err := s.GetAccount().CheckConnect(ctx, &pb.CheckConnectRequest{})
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	if data.HealthCheck == nil || !data.HealthCheck.IsAlive {
		return fmt.Errorf("not connected to MT5 terminal")
	}

	return nil
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region QUICK BALANCE METHODS
// ══════════════════════════════════════════════════════════════════════════════

// GetBalance returns the current account balance (deposit amount).
// This is the initial deposit plus/minus closed position profits/losses,
// not affected by floating profit. Uses 3-second timeout.
//
// RETURNS:
//   Current balance as float64, or error if query fails
func (s *MT5Sugar) GetBalance() (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_BALANCE)
}

// GetEquity returns the current account equity (balance + floating profit).
// Equity = Balance + Profit from open positions. This is the real-time value
// of your account including unrealized gains/losses. Uses 3-second timeout.
//
// RETURNS:
//   Current equity as float64, or error if query fails
func (s *MT5Sugar) GetEquity() (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_EQUITY)
}

// GetMargin returns the amount of margin currently used by open positions.
// This is the collateral locked by the broker for your active trades.
// Uses 3-second timeout.
//
// RETURNS:
//   Used margin as float64, or error if query fails
func (s *MT5Sugar) GetMargin() (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN)
}

// GetFreeMargin returns the amount of margin available for new positions.
// Free Margin = Equity - Used Margin. This is how much you can use for new trades.
// Uses 3-second timeout.
//
// RETURNS:
//   Free margin as float64, or error if query fails
func (s *MT5Sugar) GetFreeMargin() (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE)
}

// GetMarginLevel returns the margin level percentage.
// Margin Level = (Equity / Used Margin) * 100. Values below 100% indicate
// danger of margin call. Returns 0 if no positions are open. Uses 3-second timeout.
//
// RETURNS:
//   Margin level percentage as float64, or error if query fails
func (s *MT5Sugar) GetMarginLevel() (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL)
}

// GetProfit returns the total floating profit/loss from all open positions.
// This is the unrealized profit that's not yet added to balance. Positive values
// mean profit, negative mean loss. Uses 3-second timeout.
//
// RETURNS:
//   Total floating P/L as float64, or error if query fails
func (s *MT5Sugar) GetProfit() (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_PROFIT)
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region PRICES & QUOTES METHODS
// ══════════════════════════════════════════════════════════════════════════════

// GetBid returns the current BID price for the specified symbol.
// BID is the price at which you can SELL. This is the real-time market price.
// Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD")
//
// RETURNS:
//   Current BID price as float64, or error if symbol not found or query fails
func (s *MT5Sugar) GetBid(symbol string) (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return 0, err
	}

	return tick.Bid, nil
}

// GetAsk returns the current ASK price for the specified symbol.
// ASK is the price at which you can BUY. This is the real-time market price.
// The spread is ASK - BID. Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD")
//
// RETURNS:
//   Current ASK price as float64, or error if symbol not found or query fails
func (s *MT5Sugar) GetAsk(symbol string) (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return 0, err
	}

	return tick.Ask, nil
}

// GetSpread returns the current spread in points for the specified symbol.
// Spread is the difference between ASK and BID in points (not price units).
// For EURUSD with 5 digits, 1 point = 0.00001. Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD")
//
// RETURNS:
//   Current spread in points as float64, or error if symbol not found
func (s *MT5Sugar) GetSpread(symbol string) (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	spread, err := s.service.GetSymbolInteger(ctx, symbol, pb.SymbolInfoIntegerProperty_SYMBOL_SPREAD)
	if err != nil {
		return 0, err
	}

	return float64(spread), nil
}

// GetPriceInfo returns complete price information for the specified symbol.
// This is a convenience method that retrieves BID, ASK, spread, and timestamp
// all in one call. More efficient than calling individual methods. Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD")
//
// RETURNS:
//   *PriceInfo structure with all price data, or error if symbol not found
func (s *MT5Sugar) GetPriceInfo(symbol string) (*PriceInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return nil, err
	}

	digits, err := s.service.GetSymbolInteger(ctx, symbol, pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS)
	if err != nil {
		return nil, err
	}

	point := 1.0
	for i := int64(0); i < digits; i++ {
		point /= 10.0
	}
	spreadPoints := (tick.Ask - tick.Bid) / point

	return &PriceInfo{
		Symbol:     symbol,
		Bid:        tick.Bid,
		Ask:        tick.Ask,
		SpreadPips: spreadPoints,
		Time:       tick.Time,
	}, nil
}

// WaitForPrice waits for a price update for the specified symbol with timeout.
// This method polls for valid price data (BID > 0 and ASK > 0) until timeout expires.
// Useful for waiting for market to open or for first price tick after connection.
//
// PARAMETERS:
//   symbol  - Trading symbol to wait for (e.g., "EURUSD")
//   timeout - Maximum time to wait (e.g., 5*time.Second)
//
// RETURNS:
//   *PriceInfo with valid price data, or error if timeout expires
func (s *MT5Sugar) WaitForPrice(symbol string, timeout time.Duration) (*PriceInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	for {
		tick, err := s.service.GetSymbolTick(ctx, symbol)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("timeout waiting for price: %w", ctx.Err())
			default:
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		if tick.Bid > 0 && tick.Ask > 0 {
			spread := tick.Ask - tick.Bid
			return &PriceInfo{
				Symbol:     symbol,
				Bid:        tick.Bid,
				Ask:        tick.Ask,
				SpreadPips: spread,
				Time:       tick.Time,
			}, nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region SIMPLE TRADING METHODS
// ══════════════════════════════════════════════════════════════════════════════

// BuyMarket opens a BUY position at current market price (instant execution).
// This is the simplest way to open a long position. Order executes immediately
// at best available ASK price. Uses 10-second timeout for order execution.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01 = micro lot, 0.1 = mini lot, 1.0 = standard lot)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected or fails
func (s *MT5Sugar) BuyMarket(symbol string, volume float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
		Volume:    volume,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("BuyMarket failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// SellMarket opens a SELL position at current market price (instant execution).
// This is the simplest way to open a short position. Order executes immediately
// at best available BID price. Uses 10-second timeout for order execution.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01 = micro lot, 0.1 = mini lot, 1.0 = standard lot)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected or fails
func (s *MT5Sugar) SellMarket(symbol string, volume float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL,
		Volume:    volume,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("SellMarket failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// BuyLimit places a pending BUY LIMIT order (executes when price drops to specified level).
// Buy Limit is used to buy at a lower price than current market. Order remains pending
// until price reaches the specified level or order is cancelled. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be BELOW current ASK for Buy Limit)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) BuyLimit(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT,
		Volume:    volume,
		Price:     &price,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("BuyLimit failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// SellLimit places a pending SELL LIMIT order (executes when price rises to specified level).
// Sell Limit is used to sell at a higher price than current market. Order remains pending
// until price reaches the specified level or order is cancelled. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be ABOVE current BID for Sell Limit)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) SellLimit(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_LIMIT,
		Volume:    volume,
		Price:     &price,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("SellLimit failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// BuyStop places a pending BUY STOP order (executes when price rises to specified level).
// Buy Stop is used to buy at a higher price than current market (breakout trading).
// Order remains pending until price reaches level or order is cancelled. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be ABOVE current ASK for Buy Stop)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) BuyStop(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP,
		Volume:    volume,
		Price:     &price,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("BuyStop failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// SellStop places a pending SELL STOP order (executes when price drops to specified level).
// Sell Stop is used to sell at a lower price than current market (breakout trading).
// Order remains pending until price reaches level or order is cancelled. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be BELOW current BID for Sell Stop)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) SellStop(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP,
		Volume:    volume,
		Price:     &price,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("SellStop failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region TRADING WITH SL/TP METHODS
// ══════════════════════════════════════════════════════════════════════════════

// BuyMarketWithSLTP opens a BUY position with Stop Loss and Take Profit.
// This is the recommended way to open positions with risk management built-in.
// Order executes immediately at market price with SL/TP set. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   sl     - Stop Loss price (must be BELOW entry price for BUY)
//   tp     - Take Profit price (must be ABOVE entry price for BUY)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected
func (s *MT5Sugar) BuyMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:     symbol,
		Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY,
		Volume:     volume,
		StopLoss:   &sl,
		TakeProfit: &tp,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("BuyMarketWithSLTP failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// SellMarketWithSLTP opens a SELL position with Stop Loss and Take Profit.
// This is the recommended way to open short positions with risk management built-in.
// Order executes immediately at market price with SL/TP set. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   sl     - Stop Loss price (must be ABOVE entry price for SELL)
//   tp     - Take Profit price (must be BELOW entry price for SELL)
func (s *MT5Sugar) SellMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:     symbol,
		Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL,
		Volume:     volume,
		StopLoss:   &sl,
		TakeProfit: &tp,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("SellMarketWithSLTP failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// BuyLimitWithSLTP places a BUY LIMIT order with Stop Loss and Take Profit.
// Combines pending order functionality with automatic risk management. Order
// remains pending until price reaches entry level. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be BELOW current ASK)
//   sl     - Stop Loss price (must be BELOW entry price)
//   tp     - Take Profit price (must be ABOVE entry price)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) BuyLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:     symbol,
		Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_LIMIT,
		Volume:     volume,
		Price:      &price,
		StopLoss:   &sl,
		TakeProfit: &tp,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("BuyLimitWithSLTP failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// SellLimitWithSLTP places a SELL LIMIT order with Stop Loss and Take Profit.
// Combines pending order functionality with automatic risk management. Order
// remains pending until price reaches entry level. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be ABOVE current BID)
//   sl     - Stop Loss price (must be ABOVE entry price)
//   tp     - Take Profit price (must be BELOW entry price)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) SellLimitWithSLTP(symbol string, volume, price, sl, tp float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderSendRequest{
		Symbol:     symbol,
		Operation:  pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_LIMIT,
		Volume:     volume,
		Price:      &price,
		StopLoss:   &sl,
		TakeProfit: &tp,
	}

	result, err := s.service.PlaceOrder(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("SellLimitWithSLTP failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return 0, fmt.Errorf("order rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return result.Order, nil
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region POSITION MANAGEMENT METHODS
// ══════════════════════════════════════════════════════════════════════════════

// ClosePosition closes a position completely by ticket number.
// This is the simplest way to close an open position. Closes at current market
// price (BID for long positions, ASK for short positions). Uses 10-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number to close
//
// RETURNS:
//   Error if close fails or position not found, nil on success
func (s *MT5Sugar) ClosePosition(ticket uint64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderCloseRequest{
		Ticket: ticket,
	}

	retCode, err := s.service.CloseOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("ClosePosition failed: %w", err)
	}

	if retCode != 10009 {
		return fmt.Errorf("close rejected, code: %d", retCode)
	}

	return nil
}

// ClosePositionPartial closes a specified volume of a position (partial close).
// This allows you to take partial profit or reduce exposure while keeping position open.
// Not all brokers support partial closes. Uses 10-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   volume - Volume to close (must be less than position volume)
//
// RETURNS:
//   Error if close fails, volume invalid, or broker doesn't support partial close
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderCloseRequest{
		Ticket: ticket,
		Volume: volume,
	}

	retCode, err := s.service.CloseOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("ClosePositionPartial failed: %w", err)
	}

	if retCode != 10009 {
		return fmt.Errorf("partial close rejected, code: %d", retCode)
	}

	return nil
}

// CloseAllPositions closes all currently open positions across all symbols.
// Iterates through all positions and attempts to close each one. Continues even
// if some closes fail. Returns count of successfully closed positions. Uses 30-second timeout.
//
// RETURNS:
//   Number of positions successfully closed (int), and error if operation fails
func (s *MT5Sugar) CloseAllPositions() (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	data, err := s.service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if err != nil {
		return 0, fmt.Errorf("failed to get positions: %w", err)
	}

	closed := 0
	for _, pos := range data.PositionInfos {
		closeReq := &pb.OrderCloseRequest{
			Ticket: pos.Ticket,
		}

		retCode, err := s.service.CloseOrder(ctx, closeReq)
		if err == nil && retCode == 10009 {
			closed++
		}
	}

	return closed, nil
}

// CloseAllBySymbol closes all open positions for a specific trading symbol.
// Useful for closing all positions of one currency pair while leaving others open.
// Continues even if some closes fail. Returns count of successfully closed positions.
// Uses 30-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol to close all positions for (e.g., "EURUSD")
//
// RETURNS:
//   Number of positions successfully closed (int), and error if operation fails
func (s *MT5Sugar) CloseAllBySymbol(symbol string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	data, err := s.service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if err != nil {
		return 0, fmt.Errorf("failed to get positions: %w", err)
	}

	closed := 0
	for _, pos := range data.PositionInfos {
		if pos.Symbol == symbol {
			closeReq := &pb.OrderCloseRequest{
				Ticket: pos.Ticket,
			}

			retCode, err := s.service.CloseOrder(ctx, closeReq)
			if err == nil && retCode == 10009 {
				closed++
			}
		}
	}

	return closed, nil
}

// ModifyPositionSL modifies the Stop Loss level of an open position.
// This allows you to move your stop loss to lock in profit or reduce risk.
// Use 0 to remove Stop Loss (if broker allows). Uses 10-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   sl     - New Stop Loss price (must be valid for position direction)
//
// RETURNS:
//   Error if modification rejected or fails, nil on success
func (s *MT5Sugar) ModifyPositionSL(ticket uint64, sl float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderModifyRequest{
		Ticket:   ticket,
		StopLoss: &sl,
	}

	result, err := s.service.ModifyOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("ModifyPositionSL failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return fmt.Errorf("modify rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return nil
}

// ModifyPositionTP modifies the Take Profit level of an open position.
// This allows you to adjust your profit target based on market conditions.
// Use 0 to remove Take Profit (if broker allows). Uses 10-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   tp     - New Take Profit price (must be valid for position direction)
//
// RETURNS:
//   Error if modification rejected or fails, nil on success
func (s *MT5Sugar) ModifyPositionTP(ticket uint64, tp float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderModifyRequest{
		Ticket:     ticket,
		TakeProfit: &tp,
	}

	result, err := s.service.ModifyOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("ModifyPositionTP failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return fmt.Errorf("modify rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return nil
}

// ModifyPositionSLTP modifies both Stop Loss and Take Profit in one operation.
// More efficient than calling ModifyPositionSL and ModifyPositionTP separately.
// This is the recommended way to update risk management levels. Uses 10-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   sl     - New Stop Loss price (must be valid for position direction)
//   tp     - New Take Profit price (must be valid for position direction)
//
// RETURNS:
//   Error if modification rejected or fails, nil on success
func (s *MT5Sugar) ModifyPositionSLTP(ticket uint64, sl, tp float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	req := &pb.OrderModifyRequest{
		Ticket:     ticket,
		StopLoss:   &sl,
		TakeProfit: &tp,
	}

	result, err := s.service.ModifyOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("ModifyPositionSLTP failed: %w", err)
	}

	if result.ReturnedCode != 10009 {
		return fmt.Errorf("modify rejected, code: %d, comment: %s", result.ReturnedCode, result.Comment)
	}

	return nil
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region POSITION INFORMATION METHODS
// ══════════════════════════════════════════════════════════════════════════════

// GetOpenPositions returns all currently open positions as protobuf PositionInfo structures.
// The positions are sorted by open time (oldest first). Each PositionInfo contains full
// details: ticket, symbol, type, volume, open price, current profit, SL/TP, etc.
// Uses 5-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionInfo with all open positions, or error if query fails
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	data, err := s.service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if err != nil {
		return nil, fmt.Errorf("GetOpenPositions failed: %w", err)
	}

	return data.PositionInfos, nil
}

// GetPositionByTicket finds and returns a specific position by its ticket number.
// This is useful when you need detailed information about a position you opened earlier.
// Returns nil if position not found (may have been closed). Uses 5-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number to search for
//
// RETURNS:
//   *pb.PositionInfo for the position, or error if not found or query fails
func (s *MT5Sugar) GetPositionByTicket(ticket uint64) (*pb.PositionInfo, error) {
	positions, err := s.GetOpenPositions()
	if err != nil {
		return nil, err
	}

	for _, pos := range positions {
		if pos.Ticket == ticket {
			return pos, nil
		}
	}

	return nil, fmt.Errorf("position with ticket %d not found", ticket)
}

// GetPositionsBySymbol returns all open positions for a specific trading symbol.
// Filters positions by symbol name. Useful for monitoring exposure to a single
// currency pair or asset. Returns empty slice if no positions found. Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol to filter by (e.g., "EURUSD", "XAUUSD")
//
// RETURNS:
//   Slice of *pb.PositionInfo for the symbol, or error if query fails
func (s *MT5Sugar) GetPositionsBySymbol(symbol string) ([]*pb.PositionInfo, error) {
	positions, err := s.GetOpenPositions()
	if err != nil {
		return nil, err
	}

	var result []*pb.PositionInfo
	for _, pos := range positions {
		if pos.Symbol == symbol {
			result = append(result, pos)
		}
	}

	return result, nil
}

// HasOpenPosition checks if there are any open positions for a specific symbol.
// This is a quick boolean check - more efficient than GetPositionsBySymbol when
// you only need to know if positions exist. Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol to check (e.g., "EURUSD", "GBPUSD")
//
// RETURNS:
//   true if at least one position exists, false otherwise, or error if query fails
func (s *MT5Sugar) HasOpenPosition(symbol string) (bool, error) {
	positions, err := s.GetPositionsBySymbol(symbol)
	if err != nil {
		return false, err
	}

	return len(positions) > 0, nil
}

// CountOpenPositions returns the total number of currently open positions.
// This is more efficient than len(GetOpenPositions()) as it queries the count directly
// from MT5 without retrieving full position details. Uses 3-second timeout.
//
// RETURNS:
//   Total number of open positions (int), or error if query fails
func (s *MT5Sugar) CountOpenPositions() (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	count, err := s.service.GetPositionsTotal(ctx)
	if err != nil {
		return 0, fmt.Errorf("CountOpenPositions failed: %w", err)
	}

	return int(count), nil
}

// GetTotalProfit calculates and returns total floating profit/loss from all open positions.
// This sums up the profit field from all positions. Positive value means total profit,
// negative means total loss. Returns 0 if no positions open. Uses 5-second timeout.
//
// RETURNS:
//   Total profit/loss as float64, or error if query fails
func (s *MT5Sugar) GetTotalProfit() (float64, error) {
	positions, err := s.GetOpenPositions()
	if err != nil {
		return 0, err
	}

	var totalProfit float64
	for _, pos := range positions {
		totalProfit += pos.Profit
	}

	return totalProfit, nil
}

// GetProfitBySymbol calculates and returns total floating profit/loss for a specific symbol.
// This sums up profit from all positions matching the symbol. Useful for tracking
// per-symbol performance. Returns 0 if no positions for symbol. Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol to calculate profit for (e.g., "EURUSD")
//
// RETURNS:
//   Total profit/loss for symbol as float64, or error if query fails
func (s *MT5Sugar) GetProfitBySymbol(symbol string) (float64, error) {
	positions, err := s.GetPositionsBySymbol(symbol)
	if err != nil {
		return 0, err
	}

	var totalProfit float64
	for _, pos := range positions {
		totalProfit += pos.Profit
	}

	return totalProfit, nil
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region HISTORY & PROFIT ANALYSIS METHODS
// ══════════════════════════════════════════════════════════════════════════════

// GetDealsToday returns all closed positions (deals) from today (00:00 to now).
// Automatically calculates today's date range. Each deal contains full information:
// ticket, symbol, volume, profit, open/close times, etc. Uses 5-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with today's deals, or error if query fails
func (s *MT5Sugar) GetDealsToday() ([]*pb.PositionHistoryInfo, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	data, err := s.service.GetPositionsHistory(ctx,
		pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
		&startOfDay, &now, nil, nil)
	if err != nil {
		return nil, err
	}

	return data.HistoryPositions, nil
}

// GetDealsYesterday returns all closed positions (deals) from yesterday (full day).
// Automatically calculates yesterday's date range (00:00 to 23:59:59).
// Useful for analyzing previous day's performance. Uses 5-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with yesterday's deals, or error if query fails
func (s *MT5Sugar) GetDealsYesterday() ([]*pb.PositionHistoryInfo, error) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	startOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, yesterday.Location())

	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	data, err := s.service.GetPositionsHistory(ctx,
		pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
		&startOfYesterday, &endOfYesterday, nil, nil)
	if err != nil {
		return nil, err
	}

	return data.HistoryPositions, nil
}

// GetDealsThisWeek returns all closed positions (deals) from this week.
// Week starts on Monday. Automatically calculates start of week (Monday 00:00)
// to current time. Useful for weekly performance tracking. Uses 5-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with this week's deals, or error if query fails
func (s *MT5Sugar) GetDealsThisWeek() ([]*pb.PositionHistoryInfo, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	data, err := s.service.GetPositionsHistory(ctx,
		pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
		&startOfWeek, &now, nil, nil)
	if err != nil {
		return nil, err
	}

	return data.HistoryPositions, nil
}

// GetDealsThisMonth returns all closed positions (deals) from this month.
// Automatically calculates start of month (1st day 00:00) to current time.
// Useful for monthly performance tracking and reports. Uses 30-second timeout
// (longer than day/week queries due to potentially large data volume).
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with this month's deals, or error if query fails
func (s *MT5Sugar) GetDealsThisMonth() ([]*pb.PositionHistoryInfo, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	data, err := s.service.GetPositionsHistory(ctx,
		pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
		&startOfMonth, &now, nil, nil)
	if err != nil {
		return nil, err
	}

	return data.HistoryPositions, nil
}

// GetDealsDateRange returns all closed positions (deals) within a custom date range.
// You specify exact start and end times. Useful for custom period analysis,
// backtesting, or generating reports for specific time frames. Uses 30-second timeout
// (longer to accommodate large date ranges with many deals).
//
// PARAMETERS:
//   from - Start date/time for the range (inclusive)
//   to   - End date/time for the range (inclusive)
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with deals in range, or error if query fails
func (s *MT5Sugar) GetDealsDateRange(from, to time.Time) ([]*pb.PositionHistoryInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	data, err := s.service.GetPositionsHistory(ctx,
		pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
		&from, &to, nil, nil)
	if err != nil {
		return nil, err
	}

	return data.HistoryPositions, nil
}

// GetProfitToday calculates and returns total realized profit/loss from today's closed positions.
// This sums up the profit from all deals closed today (00:00 to now). Positive means
// net profit, negative means net loss. Returns 0 if no deals today. Uses 5-second timeout.
//
// RETURNS:
//   Total profit/loss from today's deals as float64, or error if query fails
func (s *MT5Sugar) GetProfitToday() (float64, error) {
	deals, err := s.GetDealsToday()
	if err != nil {
		return 0, err
	}

	var totalProfit float64
	for _, deal := range deals {
		totalProfit += deal.Profit
	}

	return totalProfit, nil
}

// GetProfitThisWeek calculates and returns total realized profit/loss from this week's deals.
// This sums up profit from all deals closed this week (Monday to now). Positive means
// net profit, negative means net loss. Returns 0 if no deals this week. Uses 5-second timeout.
//
// RETURNS:
//   Total profit/loss from this week's deals as float64, or error if query fails
func (s *MT5Sugar) GetProfitThisWeek() (float64, error) {
	deals, err := s.GetDealsThisWeek()
	if err != nil {
		return 0, err
	}

	var totalProfit float64
	for _, deal := range deals {
		totalProfit += deal.Profit
	}

	return totalProfit, nil
}

// GetProfitThisMonth calculates and returns total realized profit/loss from this month's deals.
// This sums up profit from all deals closed this month (1st to now). Positive means net
// profit, negative means net loss. Returns 0 if no deals this month. Uses 5-second timeout.
//
// RETURNS:
//   Total profit/loss from this month's deals as float64, or error if query fails
func (s *MT5Sugar) GetProfitThisMonth() (float64, error) {
	deals, err := s.GetDealsThisMonth()
	if err != nil {
		return 0, err
	}

	var totalProfit float64
	for _, deal := range deals {
		totalProfit += deal.Profit
	}

	return totalProfit, nil
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region SYMBOL INFORMATION METHODS
// ══════════════════════════════════════════════════════════════════════════════

// SymbolInfo holds comprehensive symbol information in one convenient structure.
// This provides all essential symbol parameters for trading calculations.
//
// FIELDS:
//   Name          - Symbol name (e.g., "EURUSD")
//   Bid           - Current BID price
//   Ask           - Current ASK price
//   Digits        - Number of decimal places
//   Point         - Point size (minimal price change)
//   VolumeMin     - Minimum volume for trading
//   VolumeMax     - Maximum volume for trading
//   VolumeStep    - Volume step
//   Spread        - Current spread in points
//   StopLevel     - Minimum stop level in points
//   ContractSize  - Contract size (for 1 lot)
type SymbolInfo struct {
	Name         string
	Bid          float64
	Ask          float64
	Digits       int32
	Point        float64
	VolumeMin    float64
	VolumeMax    float64
	VolumeStep   float64
	Spread       int32
	StopLevel    int32
	ContractSize float64
}

// GetSymbolInfo retrieves comprehensive information about a symbol in one call.
// This is more efficient than calling individual methods for each property.
// Perfect for validation before placing orders. Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD")
//
// RETURNS:
//   *SymbolInfo structure with all important symbol parameters, or error if symbol not found
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	// Get symbol parameters using Service layer
	symbolName := symbol
	params, _, err := s.service.GetSymbolParamsMany(ctx, &symbolName, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GetSymbolInfo failed: %w", err)
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}

	p := params[0]

	// Get stop level
	stopLevel, err := s.service.GetSymbolInteger(ctx, symbol, pb.SymbolInfoIntegerProperty_SYMBOL_TRADE_STOPS_LEVEL)
	if err != nil {
		stopLevel = 0 // Some symbols don't have stop level
	}

	return &SymbolInfo{
		Name:         p.Name,
		Bid:          p.Bid,
		Ask:          p.Ask,
		Digits:       p.Digits,
		Point:        p.Point,
		VolumeMin:    p.VolumeMin,
		VolumeMax:    p.VolumeMax,
		VolumeStep:   p.VolumeStep,
		Spread:       p.Spread,
		StopLevel:    int32(stopLevel),
		ContractSize: p.TradeContractSize,
	}, nil
}

// GetAllSymbols retrieves a list of all available trading symbols.
// This returns symbol names only (not full info). Useful for discovering available
// instruments or building symbol selection menus. Uses 15-second timeout
// (longer than single symbol queries due to potentially large number of symbols).
//
// RETURNS:
//   Slice of symbol names ([]string), or error if query fails
func (s *MT5Sugar) GetAllSymbols() ([]string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 15*time.Second)
	defer cancel()

	params, _, err := s.service.GetSymbolParamsMany(ctx, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GetAllSymbols failed: %w", err)
	}

	symbols := make([]string, len(params))
	for i, p := range params {
		symbols[i] = p.Name
	}

	return symbols, nil
}

// IsSymbolAvailable checks if a symbol exists and is available for trading.
// This verifies both existence and trading permissions. More comprehensive than
// just checking if symbol name is valid. Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Symbol name to check (e.g., "EURUSD")
//
// RETURNS:
//   true if symbol exists and is tradeable, false otherwise, or error if query fails
func (s *MT5Sugar) IsSymbolAvailable(symbol string) (bool, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	exists, _, err := s.service.SymbolExist(ctx, symbol)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	// Check if symbol is synchronized (has data)
	synced, err := s.service.IsSymbolSynchronized(ctx, symbol)
	if err != nil {
		return false, err
	}

	return synced, nil
}

// GetMinStopLevel returns the minimum allowed distance for Stop Loss/Take Profit
// in points. This is broker-enforced minimum distance from current price to SL/TP.
// If 0, there's no minimum (market execution). Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD")
//
// RETURNS:
//   Minimum stop level in points (int64), or error if symbol not found
func (s *MT5Sugar) GetMinStopLevel(symbol string) (int64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return s.service.GetSymbolInteger(ctx, symbol, pb.SymbolInfoIntegerProperty_SYMBOL_TRADE_STOPS_LEVEL)
}

// GetSymbolDigits returns the number of decimal places for the symbol price.
// For example, EURUSD typically has 5 digits (1.08123), gold might have 2.
// This is essential for proper price formatting and calculations. Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD")
//
// RETURNS:
//   Number of decimal places (int32), or error if symbol not found
func (s *MT5Sugar) GetSymbolDigits(symbol string) (int32, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	digits, err := s.service.GetSymbolInteger(ctx, symbol, pb.SymbolInfoIntegerProperty_SYMBOL_DIGITS)
	if err != nil {
		return 0, err
	}

	return int32(digits), nil
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region RISK MANAGEMENT METHODS
// ══════════════════════════════════════════════════════════════════════════════

// CalculatePositionSize calculates the optimal lot size based on risk percentage
// with automatic margin limit protection. This is THE MOST IMPORTANT risk management
// tool - automatically calculates position size considering BOTH risk and margin limits.
//
// ALGORITHM:
//   1. Calculate size based on risk: (Balance * RiskPercent / 100) / (StopLossPips * PipValue)
//   2. Calculate max size based on free margin (with 80% safety buffer)
//   3. Return MINIMUM of the two - prevents margin calls!
//
// PARAMETERS:
//   symbol       - Trading symbol (e.g., "EURUSD")
//   riskPercent  - Percentage of balance to risk (e.g., 2.0 = 2%)
//   stopLossPips - Stop Loss distance in points (not price!)
//
// RETURNS:
//   Recommended lot size (float64), or error if calculation fails or insufficient margin
//
// EXAMPLE:
//   Balance: $10,000, Risk: 2% ($200), SL: 50 pips
//   → Risk-based: 0.40 lots, Margin-limited: 0.30 lots → Returns: 0.30 lots
func (s *MT5Sugar) CalculatePositionSize(symbol string, riskPercent, stopLossPips float64) (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	// Get account balance
	balance, err := s.GetBalance()
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	// Get symbol info
	info, err := s.GetSymbolInfo(symbol)
	if err != nil {
		return 0, err
	}

	// Calculate risk amount
	riskAmount := balance * riskPercent / 100.0

	// Calculate pip value for 1 lot
	// For forex: pip value = (contract size * point)
	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return 0, err
	}

	currentPrice := tick.Bid
	if currentPrice == 0 {
		currentPrice = tick.Ask
	}

	// Pip value per lot
	pipValue := info.ContractSize * info.Point

	// Calculate position size based on risk
	positionSize := riskAmount / (stopLossPips * pipValue)

	// Round down to volume step using proper rounding to avoid float precision issues
	steps := positionSize / info.VolumeStep
	positionSize = float64(int(steps)) * info.VolumeStep

	// Ensure within min/max limits
	if positionSize < info.VolumeMin {
		positionSize = info.VolumeMin
	}
	if positionSize > info.VolumeMax {
		positionSize = info.VolumeMax
	}

	// CRITICAL: Check margin limit - prevent margin calls
	// Calculate maximum lot size based on available margin
	freeMargin, err := s.GetFreeMargin()
	if err != nil {
		return 0, fmt.Errorf("failed to get free margin: %w", err)
	}

	// Calculate margin required for 1 lot
	marginReq := &pb.OrderCalcMarginRequest{
		Symbol:    symbol,
		OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Volume:    1.0,
		OpenPrice: currentPrice,
	}

	marginForOneLot, err := s.service.CalculateMargin(ctx, marginReq)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate margin: %w", err)
	}

	if marginForOneLot > 0 {
		// Calculate max lots with 80% safety buffer
		maxLots := (freeMargin * 0.8) / marginForOneLot

		// Round down to volume step using proper rounding to avoid float precision issues
		steps := maxLots / info.VolumeStep
		maxLots = float64(int(steps)) * info.VolumeStep

		// Return minimum of risk-based size and margin-limited size
		if positionSize > maxLots {
			positionSize = maxLots
		}
	}

	// Final check - ensure not below minimum
	if positionSize < info.VolumeMin {
		return 0, fmt.Errorf("insufficient margin: recommended size %.2f is below minimum %.2f", positionSize, info.VolumeMin)
	}

	return positionSize, nil
}

// GetMaxLotSize calculates the maximum lot size you can open with current free margin.
// This helps prevent margin calls by showing your maximum trading capacity.
// Uses conservative estimate with safety buffer. Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD")
//
// RETURNS:
//   Maximum safe lot size (float64), or error if calculation fails
func (s *MT5Sugar) GetMaxLotSize(symbol string) (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	// Get free margin
	freeMargin, err := s.GetFreeMargin()
	if err != nil {
		return 0, err
	}

	// Get symbol info
	info, err := s.GetSymbolInfo(symbol)
	if err != nil {
		return 0, err
	}

	// Calculate margin required for 1 lot
	// Get current price
	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return 0, err
	}

	marginReq := &pb.OrderCalcMarginRequest{
		Symbol:    symbol,
		OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Volume:    1.0,
		OpenPrice: tick.Ask,
	}

	marginForOneLot, err := s.service.CalculateMargin(ctx, marginReq)
	if err != nil {
		return 0, err
	}

	if marginForOneLot == 0 {
		return 0, fmt.Errorf("failed to calculate margin for %s", symbol)
	}

	// Calculate max lots with 80% safety buffer
	maxLots := (freeMargin * 0.8) / marginForOneLot

	// Round down to volume step using proper rounding to avoid float precision issues
	steps := maxLots / info.VolumeStep
	maxLots = float64(int(steps)) * info.VolumeStep

	// Ensure within limits
	if maxLots < info.VolumeMin {
		return 0, nil // Not enough margin even for minimum
	}
	if maxLots > info.VolumeMax {
		maxLots = info.VolumeMax
	}

	return maxLots, nil
}

// CanOpenPosition checks if it's possible to open a position with specified volume.
// This performs comprehensive validation: margin check, volume limits, symbol availability.
// Always call this before PlaceOrder to prevent rejections. Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD")
//   volume - Desired lot size (e.g., 0.1)
//
// RETURNS:
//   can    - true if position can be opened
//   reason - explanation if can't open, empty if can
//   error  - error if check failed
func (s *MT5Sugar) CanOpenPosition(symbol string, volume float64) (bool, string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	// Check symbol availability
	available, err := s.IsSymbolAvailable(symbol)
	if err != nil {
		return false, "", err
	}
	if !available {
		return false, fmt.Sprintf("symbol %s is not available", symbol), nil
	}

	// Get symbol info
	info, err := s.GetSymbolInfo(symbol)
	if err != nil {
		return false, "", err
	}

	// Check volume limits
	if volume < info.VolumeMin {
		return false, fmt.Sprintf("volume %.2f below minimum %.2f", volume, info.VolumeMin), nil
	}
	if volume > info.VolumeMax {
		return false, fmt.Sprintf("volume %.2f exceeds maximum %.2f", volume, info.VolumeMax), nil
	}

	// Check volume step with tolerance for float64 precision
	steps := volume / info.VolumeStep
	roundedSteps := float64(int(steps + 0.5)) // Round to nearest integer
	tolerance := 0.0001

	// Check if volume is approximately a multiple of step
	actualVolume := roundedSteps * info.VolumeStep
	if (volume < actualVolume-tolerance) || (volume > actualVolume+tolerance) {
		return false, fmt.Sprintf("volume %.2f not a multiple of step %.2f", volume, info.VolumeStep), nil
	}

	// Calculate required margin
	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return false, "", err
	}

	marginReq := &pb.OrderCalcMarginRequest{
		Symbol:    symbol,
		OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Volume:    volume,
		OpenPrice: tick.Ask,
	}

	requiredMargin, err := s.service.CalculateMargin(ctx, marginReq)
	if err != nil {
		return false, "", err
	}

	// Check free margin
	freeMargin, err := s.GetFreeMargin()
	if err != nil {
		return false, "", err
	}

	if requiredMargin > freeMargin {
		return false, fmt.Sprintf("insufficient margin: need %.2f, have %.2f", requiredMargin, freeMargin), nil
	}

	return true, "", nil
}

// CalculateRequiredMargin calculates how much margin is required to open a position.
// This helps you plan your trades and manage account exposure. Considers leverage
// and symbol specifications. Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD")
//   volume - Desired lot size (e.g., 0.1)
//
// RETURNS:
//   Required margin amount (float64), or error if calculation fails
func (s *MT5Sugar) CalculateRequiredMargin(symbol string, volume float64) (float64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	// Get current price
	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return 0, err
	}

	req := &pb.OrderCalcMarginRequest{
		Symbol:    symbol,
		OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Volume:    volume,
		OpenPrice: tick.Ask,
	}

	return s.service.CalculateMargin(ctx, req)
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region TRADING HELPERS METHODS
// ══════════════════════════════════════════════════════════════════════════════

// CalculateSLTP calculates Stop Loss and Take Profit prices from entry price and pip distances.
// This converts pip distances to actual prices based on symbol specifications.
// Handles both BUY and SELL directions correctly. Uses 3-second timeout.
//
// PARAMETERS:
//   symbol         - Trading symbol (e.g., "EURUSD")
//   direction      - "BUY" or "SELL"
//   entryPrice     - Entry price (use 0 for current market price)
//   stopLossPips   - Distance to SL in points (e.g., 50)
//   takeProfitPips - Distance to TP in points (e.g., 100)
//
// RETURNS:
//   sl    - Stop Loss price
//   tp    - Take Profit price
//   error - error if calculation fails
//
// EXAMPLE:
//   EURUSD BUY at 1.08500, SL=50 pips, TP=100 pips
//   → SL=1.08000, TP=1.09000
func (s *MT5Sugar) CalculateSLTP(symbol, direction string, entryPrice, stopLossPips, takeProfitPips float64) (float64, float64, error) {
	// Get symbol info for point size
	info, err := s.GetSymbolInfo(symbol)
	if err != nil {
		return 0, 0, err
	}

	// If entry price is 0, use current market price
	if entryPrice == 0 {
		if direction == "BUY" {
			entryPrice, err = s.GetAsk(symbol)
		} else {
			entryPrice, err = s.GetBid(symbol)
		}
		if err != nil {
			return 0, 0, err
		}
	}

	var sl, tp float64

	if direction == "BUY" {
		// BUY: SL below entry, TP above entry
		sl = entryPrice - (stopLossPips * info.Point)
		tp = entryPrice + (takeProfitPips * info.Point)
	} else {
		// SELL: SL above entry, TP below entry
		sl = entryPrice + (stopLossPips * info.Point)
		tp = entryPrice - (takeProfitPips * info.Point)
	}

	return sl, tp, nil
}

// BuyMarketWithPips opens a BUY position with SL/TP specified in pips (not price!).
// This is more intuitive than BuyMarketWithSLTP - you specify risk/reward in pips
// and the method calculates exact prices automatically. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol         - Trading symbol (e.g., "EURUSD")
//   volume         - Lot size (e.g., 0.1)
//   stopLossPips   - Stop Loss distance in pips from entry (e.g., 50)
//   takeProfitPips - Take Profit distance in pips from entry (e.g., 100)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected
//
// EXAMPLE:
//   ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 50, 100)
//   // Opens BUY at market, SL = entry - 50 pips, TP = entry + 100 pips
func (s *MT5Sugar) BuyMarketWithPips(symbol string, volume, stopLossPips, takeProfitPips float64) (uint64, error) {
	// Calculate SL/TP prices
	sl, tp, err := s.CalculateSLTP(symbol, "BUY", 0, stopLossPips, takeProfitPips)
	if err != nil {
		return 0, err
	}

	// Use existing BuyMarketWithSLTP
	return s.BuyMarketWithSLTP(symbol, volume, sl, tp)
}

// SellMarketWithPips opens a SELL position with SL/TP specified in pips (not price!).
// This is more intuitive than SellMarketWithSLTP - you specify risk/reward in pips
// and the method calculates exact prices automatically. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol         - Trading symbol (e.g., "EURUSD")
//   volume         - Lot size (e.g., 0.1)
//   stopLossPips   - Stop Loss distance in pips from entry (e.g., 50)
//   takeProfitPips - Take Profit distance in pips from entry (e.g., 100)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected
//
// EXAMPLE:
//   ticket, _ := sugar.SellMarketWithPips("EURUSD", 0.1, 50, 100)
//   // Opens SELL at market, SL = entry + 50 pips, TP = entry - 100 pips
func (s *MT5Sugar) SellMarketWithPips(symbol string, volume, stopLossPips, takeProfitPips float64) (uint64, error) {
	// Calculate SL/TP prices
	sl, tp, err := s.CalculateSLTP(symbol, "SELL", 0, stopLossPips, takeProfitPips)
	if err != nil {
		return 0, err
	}

	// Use existing SellMarketWithSLTP
	return s.SellMarketWithSLTP(symbol, volume, sl, tp)
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region ACCOUNT INFORMATION METHODS
// ══════════════════════════════════════════════════════════════════════════════

// AccountInfo holds comprehensive account information in one structure.
// This provides complete account snapshot for monitoring and reporting.
//
// FIELDS:
//   Login       - Account login number
//   Balance     - Account balance
//   Equity      - Current equity (balance + floating P/L)
//   Margin      - Used margin
//   FreeMargin  - Free margin available
//   MarginLevel - Margin level percentage
//   Profit      - Total floating profit/loss
//   Currency    - Account currency (USD, EUR, etc.)
//   Leverage    - Account leverage (e.g., 100 for 1:100)
//   Company     - Broker company name
type AccountInfo struct {
	Login       int64
	Balance     float64
	Equity      float64
	Margin      float64
	FreeMargin  float64
	MarginLevel float64
	Profit      float64
	Currency    string
	Leverage    int64
	Company     string
}

// GetAccountInfo retrieves complete account information in one call.
// This is more efficient than calling individual Get* methods. Perfect for
// account monitoring dashboards or trading reports. Uses 5-second timeout.
//
// RETURNS:
//   *AccountInfo structure with all account data, or error if query fails
func (s *MT5Sugar) GetAccountInfo() (*AccountInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	// Use Service layer GetAccountSummary
	summary, err := s.service.GetAccountSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAccountInfo failed: %w", err)
	}

	// Get additional fields
	margin, _ := s.GetMargin()
	freeMargin, _ := s.GetFreeMargin()
	marginLevel, _ := s.GetMarginLevel()
	profit, _ := s.GetProfit()

	return &AccountInfo{
		Login:       summary.Login,
		Balance:     summary.Balance,
		Equity:      summary.Equity,
		Margin:      margin,
		FreeMargin:  freeMargin,
		MarginLevel: marginLevel,
		Profit:      profit,
		Currency:    summary.Currency,
		Leverage:    summary.Leverage,
		Company:     summary.CompanyName,
	}, nil
}

// DailyStats holds trading statistics for today.
// Useful for tracking daily performance and generating reports.
//
// FIELDS:
//   TotalDeals   - Total number of closed deals today
//   WinningDeals - Number of profitable deals
//   LosingDeals  - Number of losing deals
//   WinRate      - Win rate percentage (0-100)
//   TotalProfit  - Total realized profit/loss today
//   BestDeal     - Largest profitable deal
//   WorstDeal    - Largest losing deal
type DailyStats struct {
	TotalDeals   int
	WinningDeals int
	LosingDeals  int
	WinRate      float64
	TotalProfit  float64
	BestDeal     float64
	WorstDeal    float64
}

// GetDailyStats calculates trading statistics for today (00:00 to now).
// This analyzes all closed positions from today and provides performance metrics.
// Perfect for daily reports and performance tracking. Uses 5-second timeout.
//
// RETURNS:
//   *DailyStats structure with today's performance, or error if query fails
func (s *MT5Sugar) GetDailyStats() (*DailyStats, error) {
	deals, err := s.GetDealsToday()
	if err != nil {
		return nil, err
	}

	stats := &DailyStats{
		TotalDeals: len(deals),
	}

	if len(deals) == 0 {
		return stats, nil
	}

	for _, deal := range deals {
		stats.TotalProfit += deal.Profit

		if deal.Profit > 0 {
			stats.WinningDeals++
			if deal.Profit > stats.BestDeal {
				stats.BestDeal = deal.Profit
			}
		} else if deal.Profit < 0 {
			stats.LosingDeals++
			if deal.Profit < stats.WorstDeal {
				stats.WorstDeal = deal.Profit
			}
		}
	}

	if stats.TotalDeals > 0 {
		stats.WinRate = float64(stats.WinningDeals) / float64(stats.TotalDeals) * 100.0
	}

	return stats, nil
}

// #endregion
