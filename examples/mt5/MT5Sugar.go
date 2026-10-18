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
   │  8. HISTORY & PROFIT ANALYSIS (8 methods)                   │
   ├─────────────────────────────────────────────────────────────┤
   │  • GetDealsToday()       - Closed positions opened today    │
   │  • GetDealsYesterday()   - ... opened yesterday             │
   │  • GetDealsThisWeek()    - ... opened this week (Mon-now)   │
   │  • GetDealsThisMonth()   - ... opened this month            │
   │  • GetDealsDateRange()   - ... opened in a custom range     │
   │  • GetProfitToday()      - Sum of Profit of today's ones    │
   │  • GetProfitThisWeek()   - ... of this week's ones          │
   │  • GetProfitThisMonth()  - ... of this month's ones         │
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
   │  • CalculateSLTP()       - Point distances to SL/TP prices  │
   │  • BuyMarketWithPips()   - BUY, SL/TP distances in points   │
   │  • SellMarketWithPips()  - SELL, SL/TP distances in points  │
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
   • All methods have built-in timeouts (3-130 seconds depending on operation)
   • QuickConnect: the server waits up to 120 s for the terminal, the call has a 130 s deadline
   • Trading (orders, modify, close): 30 seconds per request (the account layer's default)
   • Balance/Price queries timeout: 3 seconds
   • Positions, symbol info, account info: 10 seconds
   • History queries timeout: 15 s for day/week, 30 s for month/custom range
   • Risk management timeout: 10 seconds (CalculatePositionSize with margin checks)
   • Symbol list queries timeout: 60 seconds (GetAllSymbols - reads all pages)
   • CloseAll*: 10 s to list the positions, then 30 s for each close
   • Times: tick/price Time fields are broker server time; *Utc fields are UTC. History windows
     (GetDeals..., GetProfit..., GetDailyStats) use the local calendar of the computer, sent to the server as UTC
   • Distances named "...Pips" in risk/helper methods are POINTS (10 points = 1 pip on 5-digit symbols)
   • Always call Disconnect() when done (also on errors), or the terminal keeps running on the server
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
	"errors"
	"fmt"
	"math"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	helpers "github.com/MetaRPC/GoMT5/package/Helpers"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
//   SpreadPips - Spread in POINTS despite the field name: (Ask - Bid) / SYMBOL_POINT
//                (EURUSD 5 digits: 1 point = 0.00001, 10 points = 1 pip)
//   Time       - Time of the last tick in broker server time (the broker's wall clock in a UTC time.Time)
//   TimeUtc    - The same moment in UTC
type PriceInfo struct {
	Symbol     string
	Bid        float64
	Ask        float64
	SpreadPips float64
	Time       time.Time
	TimeUtc    time.Time
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

// Disconnect stops the MT5 terminal on the server and closes the gRPC connection (see MT5Service.Disconnect).
// Call it when the program is done - also on error paths - or the terminal keeps running on the server.
// deleteOnDisconnect=true deletes the terminal instance instead of stopping it. A second call is a no-op.
func (s *MT5Sugar) Disconnect(deleteOnDisconnect ...bool) error {
	if s == nil || s.service == nil {
		return nil
	}
	return s.service.Disconnect(deleteOnDisconnect...)
}

// ══════════════════════════════════════════════════════════════════════════════
// #region CONNECTION METHODS
// ══════════════════════════════════════════════════════════════════════════════

// QuickConnect connects to MT5 terminal using cluster name (RECOMMENDED).
// This is the easiest connection method - just provide your broker's cluster name.
// The server waits up to 120 seconds (ConnectExRequest.TimeoutSeconds) for the terminal to start and log in;
// the call itself has a 130-second deadline. Call Disconnect when you are done - also when QuickConnect
// returned an error, since the server may already have started the terminal.
//
// PARAMETERS:
//   clusterName - MT5 cluster identifier (e.g., "FxPro-MT5 Demo", "ICMarkets-Live02")
//
// RETURNS:
//   Error if connection fails, nil on success
func (s *MT5Sugar) QuickConnect(clusterName string) error {
	// A terminal that is not running yet starts first: that takes longer than 30 seconds at times. The server waits
	// up to TimeoutSeconds for the login; the call itself gets that plus a margin.
	timeoutSeconds := uint32(120)
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(timeoutSeconds+10)*time.Second)
	defer cancel()

	req := &pb.ConnectExRequest{
		User:           s.user,
		Password:       s.password,
		MtClusterName:  clusterName,
		TimeoutSeconds: &timeoutSeconds,
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

	return s.priceInfoFromTick(ctx, symbol, tick)
}

// priceInfoFromTick builds a PriceInfo from a tick; the spread is converted to points (SYMBOL_POINT).
func (s *MT5Sugar) priceInfoFromTick(ctx context.Context, symbol string, tick *SymbolTick) (*PriceInfo, error) {
	point, err := s.service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_POINT)
	if err != nil {
		return nil, err
	}
	if point <= 0 {
		return nil, fmt.Errorf("symbol %s has no point size (SYMBOL_POINT %v)", symbol, point)
	}
	spreadPoints := math.Round((tick.Ask-tick.Bid)/point*10) / 10 // float noise: 3.9999999 -> 4

	return &PriceInfo{
		Symbol:     symbol,
		Bid:        tick.Bid,
		Ask:        tick.Ask,
		SpreadPips: spreadPoints,
		Time:       tick.Time,
		TimeUtc:    tick.TimeUtc,
	}, nil
}

// WaitForPrice waits until the symbol has a valid price (BID > 0 and ASK > 0) and returns it.
// It polls the last tick every 100 ms until the timeout expires; a symbol that already has a price
// returns at once (it does NOT wait for the next price change). Useful right after connecting, when
// the first quote of a symbol may not have arrived yet.
//
// PARAMETERS:
//   symbol  - Trading symbol to wait for (e.g., "EURUSD")
//   timeout - Maximum time to wait (e.g., 5*time.Second)
//
// RETURNS:
//   *PriceInfo with valid price data (spread in points, like GetPriceInfo), or error if timeout expires
func (s *MT5Sugar) WaitForPrice(symbol string, timeout time.Duration) (*PriceInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	for {
		tick, err := s.service.GetSymbolTick(ctx, symbol)
		if err != nil {
			// An error from the server (unknown symbol, ...) does not go away by waiting: return it.
			var apiErr *helpers.ApiError
			if errors.As(err, &apiErr) {
				return nil, err
			}
		} else if tick.Bid > 0 && tick.Ask > 0 {
			return s.priceInfoFromTick(ctx, symbol, tick)
		}

		select {
		case <-ctx.Done():
			if err != nil {
				return nil, fmt.Errorf("timeout waiting for price (last error: %v): %w", err, ctx.Err())
			}
			return nil, fmt.Errorf("timeout waiting for price: %w", ctx.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}

// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region SIMPLE TRADING METHODS
// ══════════════════════════════════════════════════════════════════════════════

// BuyMarket opens a BUY position at current market price (instant execution).
// This is the simplest way to open a long position. Order executes immediately
// at best available ASK price. Uses 30-second timeout for order execution.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01 = micro lot, 0.1 = mini lot, 1.0 = standard lot)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected or fails
func (s *MT5Sugar) BuyMarket(symbol string, volume float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// at best available BID price. Uses 30-second timeout for order execution.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01 = micro lot, 0.1 = mini lot, 1.0 = standard lot)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected or fails
func (s *MT5Sugar) SellMarket(symbol string, volume float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// until price reaches the specified level or order is cancelled. Uses 30-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be BELOW current ASK for Buy Limit)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) BuyLimit(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// until price reaches the specified level or order is cancelled. Uses 30-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be ABOVE current BID for Sell Limit)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) SellLimit(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// Order remains pending until price reaches level or order is cancelled. Uses 30-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be ABOVE current ASK for Buy Stop)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) BuyStop(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// Order remains pending until price reaches level or order is cancelled. Uses 30-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   price  - Entry price (must be BELOW current BID for Sell Stop)
//
// RETURNS:
//   Pending order ticket number (uint64), or error if order rejected
func (s *MT5Sugar) SellStop(symbol string, volume, price float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// Order executes immediately at market price with SL/TP set. Uses 30-second timeout.
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
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// Order executes immediately at market price with SL/TP set. Uses 30-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD")
//   volume - Lot size (e.g., 0.01, 0.1, 1.0)
//   sl     - Stop Loss price (must be ABOVE entry price for SELL)
//   tp     - Take Profit price (must be BELOW entry price for SELL)
func (s *MT5Sugar) SellMarketWithSLTP(symbol string, volume, sl, tp float64) (uint64, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// remains pending until price reaches entry level. Uses 30-second timeout.
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
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// remains pending until price reaches entry level. Uses 30-second timeout.
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
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// price (BID for long positions, ASK for short positions). Uses 30-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number to close
//
// RETURNS:
//   Error if close fails or position not found, nil on success
func (s *MT5Sugar) ClosePosition(ticket uint64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	req := &pb.OrderCloseRequest{
		Ticket: ticket,
	}

	retCode, err := s.service.CloseOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("ClosePosition failed: %w", err)
	}

	if retCode != 10009 {
		return fmt.Errorf("close rejected, code: %d (%s)", retCode, helpers.GetRetCodeMessage(retCode))
	}

	return nil
}

// ClosePositionPartial closes a specified volume of a position (partial close).
// This allows you to take partial profit or reduce exposure while keeping position open.
// Not all brokers support partial closes. Uses 30-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   volume - Volume to close (must be less than position volume)
//
// RETURNS:
//   Error if close fails, volume invalid, or broker doesn't support partial close
func (s *MT5Sugar) ClosePositionPartial(ticket uint64, volume float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
		return fmt.Errorf("partial close rejected, code: %d (%s)", retCode, helpers.GetRetCodeMessage(retCode))
	}

	return nil
}

// CloseAllPositions closes all currently open positions across all symbols (pending orders are not touched).
// Iterates through all positions and attempts to close each one; it continues when a close fails.
// Listing the positions uses a 10-second timeout, each close its own 30-second timeout.
//
// RETURNS:
//   Number of positions successfully closed (int), and an error when the positions could not be listed or
//   when at least one close failed (the error names every failed ticket and why; the count is still valid)
func (s *MT5Sugar) CloseAllPositions() (int, error) {
	return s.closePositions(func(*pb.PositionInfo) bool { return true })
}

// CloseAllBySymbol closes all open positions for a specific trading symbol.
// Useful for closing all positions of one currency pair while leaving others open.
// Continues when a close fails. Listing uses a 10-second timeout, each close its own 30-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol to close all positions for (e.g., "EURUSD")
//
// RETURNS:
//   Number of positions successfully closed (int), and an error when the positions could not be listed or
//   when at least one close failed (the error names every failed ticket and why; the count is still valid)
func (s *MT5Sugar) CloseAllBySymbol(symbol string) (int, error) {
	return s.closePositions(func(pos *pb.PositionInfo) bool { return pos.Symbol == symbol })
}

// closePositions closes every open position the filter selects; failures are collected, not hidden.
func (s *MT5Sugar) closePositions(selected func(*pb.PositionInfo) bool) (int, error) {
	listCtx, cancelList := context.WithTimeout(s.ctx, 10*time.Second)
	data, err := s.service.GetOpenedOrders(listCtx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	cancelList()
	if err != nil {
		return 0, fmt.Errorf("failed to get positions: %w", err)
	}

	closed := 0
	var failures []error
	for _, pos := range data.PositionInfos {
		if !selected(pos) {
			continue
		}
		ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
		retCode, err := s.service.CloseOrder(ctx, &pb.OrderCloseRequest{Ticket: pos.Ticket})
		cancel()
		switch {
		case err != nil:
			failures = append(failures, fmt.Errorf("position #%d: %w", pos.Ticket, err))
		case retCode != 10009:
			failures = append(failures, fmt.Errorf("position #%d: close rejected, code: %d (%s)", pos.Ticket, retCode, helpers.GetRetCodeMessage(retCode)))
		default:
			closed++
		}
	}

	if len(failures) > 0 {
		return closed, fmt.Errorf("%d of %d position(s) not closed: %w", len(failures), len(failures)+closed, errors.Join(failures...))
	}
	return closed, nil
}

// ModifyPositionSL modifies the Stop Loss level of an open position.
// This allows you to move your stop loss to lock in profit or reduce risk.
// Use 0 to remove Stop Loss (if broker allows). Uses 30-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   sl     - New Stop Loss price (must be valid for position direction)
//
// RETURNS:
//   Error if modification rejected or fails, nil on success
func (s *MT5Sugar) ModifyPositionSL(ticket uint64, sl float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// Use 0 to remove Take Profit (if broker allows). Uses 30-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   tp     - New Take Profit price (must be valid for position direction)
//
// RETURNS:
//   Error if modification rejected or fails, nil on success
func (s *MT5Sugar) ModifyPositionTP(ticket uint64, tp float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// This is the recommended way to update risk management levels. Uses 30-second timeout.
//
// PARAMETERS:
//   ticket - Position ticket number
//   sl     - New Stop Loss price (must be valid for position direction)
//   tp     - New Take Profit price (must be valid for position direction)
//
// RETURNS:
//   Error if modification rejected or fails, nil on success
func (s *MT5Sugar) ModifyPositionSLTP(ticket uint64, sl, tp float64) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
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
// Uses 10-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionInfo with all open positions, or error if query fails
func (s *MT5Sugar) GetOpenPositions() ([]*pb.PositionInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	data, err := s.service.GetOpenedOrders(ctx, pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_BMT5_OPENED_ORDER_SORT_BY_OPEN_TIME_ASC)
	if err != nil {
		return nil, fmt.Errorf("GetOpenPositions failed: %w", err)
	}

	return data.PositionInfos, nil
}

// GetPositionByTicket finds and returns a specific position by its ticket number.
// This is useful when you need detailed information about a position you opened earlier.
// Returns an error if no open position has that ticket (it may have been closed). Uses 10-second timeout.
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
// currency pair or asset. Returns empty slice if no positions found. Uses 10-second timeout.
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
// you only need to know if positions exist. Uses 10-second timeout.
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
// negative means total loss. Returns 0 if no positions open. Uses 10-second timeout.
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
// per-symbol performance. Returns 0 if no positions for symbol. Uses 10-second timeout.
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
//
// HOW THE TIME WINDOWS WORK (all GetDeals*/GetProfit* methods and GetDailyStats):
//   • "Deals" here are CLOSED POSITIONS from PositionsHistory (one entry per position, not per MT5 deal).
//   • A position is selected by its OPEN time: a position opened yesterday and closed today is NOT in
//     GetDealsToday. Still-open positions are not in the list.
//   • Day/week/month boundaries are calendar boundaries of the computer running the program (time.Local):
//     "today" = local midnight to now, "this week" = local Monday 00:00 to now, "this month" = the 1st, 00:00
//     local to now. The window is sent in the request's UTC fields (position_open_time_from_utc/_to_utc),
//     so the broker's time zone does not shift it. Run the program with TZ set to the broker's time zone
//     if you want the broker's trading day instead.
//   • All pages are read (cursor paging: page_size 1000 + next_page_token), so long ranges are complete.
//   • Profit sums are the positions' Profit field: commission, swap and fee are NOT included.
// ══════════════════════════════════════════════════════════════════════════════

// positionsOpenedBetween returns ALL closed positions whose open time lies in [from, to] (real moments).
// It reads the history with cursor paging (page_size 1000, next_page_token) until the last page, so a long
// range is never cut to the first page. The timeout covers all pages together.
func (s *MT5Sugar) positionsOpenedBetween(from, to time.Time, timeout time.Duration) ([]*pb.PositionHistoryInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	var all []*pb.PositionHistoryInfo
	token := ""
	for {
		// A Go time.Time is a moment: the window goes in the *Utc fields (see MT5Service.GetPositionsHistory)
		data, err := s.GetAccount().PositionsHistory(ctx, &pb.PositionsHistoryRequest{
			SortType:                pb.AH_ENUM_POSITIONS_HISTORY_SORT_TYPE_AH_POSITION_OPEN_TIME_ASC,
			PositionOpenTimeFromUtc: timestamppb.New(from),
			PositionOpenTimeToUtc:   timestamppb.New(to),
			PageToken:               token,
			PageSize:                1000,
		})
		if err != nil {
			return nil, fmt.Errorf("PositionsHistory failed: %w", err)
		}
		all = append(all, data.HistoryPositions...)
		token = data.NextPageToken
		if token == "" || len(data.HistoryPositions) == 0 {
			return all, nil
		}
	}
}

// startOfLocalDay returns 00:00 of t's calendar day in t's time zone.
func startOfLocalDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// sumProfit adds up the Profit field (without commission, swap and fee).
func sumProfit(deals []*pb.PositionHistoryInfo) float64 {
	var total float64
	for _, deal := range deals {
		total += deal.Profit
	}
	return total
}

// GetDealsToday returns the closed positions opened today: local midnight (time.Local) to now.
// See the region note above for how the window and "deals" are defined. Uses 15-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with today's closed positions, or error if query fails
func (s *MT5Sugar) GetDealsToday() ([]*pb.PositionHistoryInfo, error) {
	now := time.Now()
	return s.positionsOpenedBetween(startOfLocalDay(now), now, 15*time.Second)
}

// GetDealsYesterday returns the closed positions opened yesterday: the whole previous local calendar day
// (00:00:00 to 23:59:59.999999999, time.Local). Useful for analyzing the previous day. Uses 15-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with yesterday's closed positions, or error if query fails
func (s *MT5Sugar) GetDealsYesterday() ([]*pb.PositionHistoryInfo, error) {
	startOfToday := startOfLocalDay(time.Now())
	startOfYesterday := startOfLocalDay(startOfToday.AddDate(0, 0, -1))
	return s.positionsOpenedBetween(startOfYesterday, startOfToday.Add(-time.Nanosecond), 15*time.Second)
}

// GetDealsThisWeek returns the closed positions opened this week.
// The week starts on Monday: local Monday 00:00 (time.Local) to now. Useful for weekly
// performance tracking. Uses 15-second timeout.
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with this week's closed positions, or error if query fails
func (s *MT5Sugar) GetDealsThisWeek() ([]*pb.PositionHistoryInfo, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := startOfLocalDay(now.AddDate(0, 0, -(weekday - 1)))
	return s.positionsOpenedBetween(startOfWeek, now, 15*time.Second)
}

// GetDealsThisMonth returns the closed positions opened this month: the 1st, 00:00 local (time.Local) to now.
// Useful for monthly performance tracking and reports. Uses 30-second timeout
// (longer than day/week queries due to potentially large data volume).
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with this month's closed positions, or error if query fails
func (s *MT5Sugar) GetDealsThisMonth() ([]*pb.PositionHistoryInfo, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return s.positionsOpenedBetween(startOfMonth, now, 30*time.Second)
}

// GetDealsDateRange returns the closed positions opened within a custom time range.
// from and to are real moments (time.Now().AddDate(0, 0, -7), time.Date(..., time.Local), time.Date(..., time.UTC)):
// they are sent as UTC, so the result does not depend on the broker's time zone. Useful for custom period
// analysis or reports. Uses 30-second timeout (longer to accommodate large ranges).
//
// PARAMETERS:
//   from - Start of the range for the position open time (inclusive)
//   to   - End of the range for the position open time (inclusive)
//
// RETURNS:
//   Slice of *pb.PositionHistoryInfo with the closed positions in range, or error if query fails
func (s *MT5Sugar) GetDealsDateRange(from, to time.Time) ([]*pb.PositionHistoryInfo, error) {
	return s.positionsOpenedBetween(from, to, 30*time.Second)
}

// GetProfitToday returns the summed Profit of the closed positions opened today (local midnight to now).
// Commission, swap and fee are not included. Positive means profit, negative means loss.
// Returns 0 if there are none. Uses 15-second timeout.
//
// RETURNS:
//   Total profit/loss of today's closed positions as float64, or error if query fails
func (s *MT5Sugar) GetProfitToday() (float64, error) {
	deals, err := s.GetDealsToday()
	if err != nil {
		return 0, err
	}
	return sumProfit(deals), nil
}

// GetProfitThisWeek returns the summed Profit of the closed positions opened this week (local Monday 00:00
// to now). Commission, swap and fee are not included. Returns 0 if there are none. Uses 15-second timeout.
//
// RETURNS:
//   Total profit/loss of this week's closed positions as float64, or error if query fails
func (s *MT5Sugar) GetProfitThisWeek() (float64, error) {
	deals, err := s.GetDealsThisWeek()
	if err != nil {
		return 0, err
	}
	return sumProfit(deals), nil
}

// GetProfitThisMonth returns the summed Profit of the closed positions opened this month (the 1st, 00:00
// local to now). Commission, swap and fee are not included. Returns 0 if there are none. Uses 30-second timeout.
//
// RETURNS:
//   Total profit/loss of this month's closed positions as float64, or error if query fails
func (s *MT5Sugar) GetProfitThisMonth() (float64, error) {
	deals, err := s.GetDealsThisMonth()
	if err != nil {
		return 0, err
	}
	return sumProfit(deals), nil
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
// Perfect for validation before placing orders. Uses 10-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD", "GBPUSD", "XAUUSD")
//
// RETURNS:
//   *SymbolInfo structure with all important symbol parameters, or error if symbol not found
func (s *MT5Sugar) GetSymbolInfo(symbol string) (*SymbolInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	// Get symbol parameters using Service layer
	symbolName := symbol
	params, _, err := s.service.GetSymbolParamsMany(ctx, &symbolName, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GetSymbolInfo failed: %w", err)
	}

	// Take the entry with exactly this name (the name filter may return similar names too).
	var p *SymbolParams
	for i := range params {
		if params[i].Name == symbol {
			p = &params[i]
			break
		}
	}
	if p == nil {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}

	// Get stop level (0 = no minimum distance; a failed request is an error, not a 0)
	stopLevel, err := s.service.GetSymbolInteger(ctx, symbol, pb.SymbolInfoIntegerProperty_SYMBOL_TRADE_STOPS_LEVEL)
	if err != nil {
		return nil, fmt.Errorf("GetSymbolInfo failed: stop level: %w", err)
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
// instruments or building symbol selection menus. SymbolParamsMany answers one page at a
// time (a request without paging returns only the first page), so this method reads page
// after page until it has SymbolsTotal names. Uses 60-second timeout for all pages together
// (brokers with stocks have thousands of symbols).
//
// RETURNS:
//   Slice of symbol names ([]string), or error if query fails
func (s *MT5Sugar) GetAllSymbols() ([]string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 60*time.Second)
	defer cancel()

	const perPage = int32(1000)
	var symbols []string
	for page := int32(1); ; page++ {
		pageNumber, itemsPerPage := page, perPage
		params, total, err := s.service.GetSymbolParamsMany(ctx, nil, nil, &pageNumber, &itemsPerPage)
		if err != nil {
			return nil, fmt.Errorf("GetAllSymbols failed (page %d): %w", page, err)
		}
		for _, p := range params {
			symbols = append(symbols, p.Name)
		}
		if len(params) == 0 || int32(len(symbols)) >= total {
			break
		}
	}

	return symbols, nil
}

// IsSymbolAvailable checks if a symbol exists in the terminal and its data is synchronized
// with the trade server (SymbolExist + SymbolIsSynchronized). It does NOT check the symbol's
// trade mode (SYMBOL_TRADE_MODE) or trading sessions. Uses 3-second timeout.
//
// PARAMETERS:
//   symbol - Symbol name to check (e.g., "EURUSD")
//
// RETURNS:
//   true if the symbol exists and is synchronized, false otherwise, or error if query fails
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
//   1. Value of 1 point for 1 lot, in the account currency:
//      PointValue = SYMBOL_TRADE_TICK_VALUE * SYMBOL_POINT / SYMBOL_TRADE_TICK_SIZE
//   2. Size based on risk: (Balance * RiskPercent / 100) / (StopLossPoints * PointValue),
//      rounded down to the volume step and capped at the symbol's max volume. A size below the
//      minimum volume is an ERROR (raising it to the minimum would risk more than RiskPercent)
//   3. Max size based on free margin (80% of it; margin for 1 lot from OrderCalcMargin for a BUY at ASK -
//      the method has no direction)
//   4. Return the SMALLER of the two - prevents margin calls!
//
// PARAMETERS:
//   symbol       - Trading symbol (e.g., "EURUSD")
//   riskPercent  - Percentage of balance to risk (e.g., 2.0 = 2%)
//   stopLossPips - Stop Loss distance in POINTS despite the name (1 point = SYMBOL_POINT,
//                  0.00001 on 5-digit EURUSD, so 10 points = 1 pip). Must be > 0.
//
// RETURNS:
//   Recommended lot size (float64), or error if the calculation fails, the risk does not allow the
//   minimum volume, or the margin is insufficient
//
// EXAMPLE (USD account, EURUSD 5 digits: 1 point of 1 lot = 1 USD):
//   Balance: $10,000, Risk: 2% ($200), SL: 500 points (= 50 pips)
//   → Risk-based: 200 / (500 * 1) = 0.40 lots; if 80% of free margin only covers 0.30 lots → Returns: 0.30 lots
func (s *MT5Sugar) CalculatePositionSize(symbol string, riskPercent, stopLossPips float64) (float64, error) {
	if stopLossPips <= 0 {
		return 0, fmt.Errorf("stop loss distance must be > 0 points, got %v", stopLossPips)
	}

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

	// Value of 1 point for 1 lot in the ACCOUNT currency (contract size * point would be in the quote
	// currency, wrong for e.g. USDJPY on a USD account).
	tickValue, err := s.service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_VALUE)
	if err != nil {
		return 0, fmt.Errorf("failed to get tick value: %w", err)
	}
	tickSize, err := s.service.GetSymbolDouble(ctx, symbol, pb.SymbolInfoDoubleProperty_SYMBOL_TRADE_TICK_SIZE)
	if err != nil {
		return 0, fmt.Errorf("failed to get tick size: %w", err)
	}
	if tickValue <= 0 || tickSize <= 0 {
		return 0, fmt.Errorf("symbol %s has no tick value yet (tick value %v, tick size %v)", symbol, tickValue, tickSize)
	}
	pointValue := tickValue * info.Point / tickSize

	// Calculate position size based on risk
	positionSize := floorToStep(riskAmount/(stopLossPips*pointValue), info.VolumeStep)

	// Within min/max limits. A size below the minimum volume is NOT raised to it: that would risk more
	// than riskPercent.
	if positionSize < info.VolumeMin {
		return 0, fmt.Errorf("risk %.2f%% with a %v-point stop allows %.4f lots, below the minimum volume %.2f of %s",
			riskPercent, stopLossPips, riskAmount/(stopLossPips*pointValue), info.VolumeMin, symbol)
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

	tick, err := s.service.GetSymbolTick(ctx, symbol)
	if err != nil {
		return 0, err
	}

	// Calculate margin required for 1 lot (a BUY opens at ASK)
	marginReq := &pb.OrderCalcMarginRequest{
		Symbol:    symbol,
		OrderType: pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		Volume:    1.0,
		OpenPrice: tick.Ask,
	}

	marginForOneLot, err := s.service.CalculateMargin(ctx, marginReq)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate margin: %w", err)
	}

	if marginForOneLot > 0 {
		// Calculate max lots with 80% safety buffer, rounded down to the volume step
		maxLots := floorToStep((freeMargin*0.8)/marginForOneLot, info.VolumeStep)

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

// floorToStep rounds a volume down to a multiple of step (with a little tolerance for float noise:
// 0.3/0.01 = 29.9999999 must give 0.30, not 0.29) and trims the float noise of the result.
func floorToStep(volume, step float64) float64 {
	if step <= 0 {
		return volume
	}
	v := math.Floor(volume/step+1e-9) * step
	return math.Round(v*1e8) / 1e8
}

// GetMaxLotSize calculates the maximum lot size you can open with current free margin.
// This helps prevent margin calls by showing your maximum trading capacity.
// Conservative estimate: 80% of the free margin divided by the margin of 1 lot, rounded down to the
// volume step and capped at the maximum volume. The margin is calculated for a BUY at the current ASK
// (the method has no direction; a SELL's margin may differ slightly). Uses 5-second timeout.
//
// PARAMETERS:
//   symbol - Trading symbol (e.g., "EURUSD")
//
// RETURNS:
//   Maximum safe lot size (float64), or an error if the calculation fails or the free margin does not
//   cover even the minimum volume
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

	// Round down to the volume step
	maxLots = floorToStep(maxLots, info.VolumeStep)

	// Ensure within limits
	if maxLots < info.VolumeMin {
		return 0, fmt.Errorf("insufficient margin: 80%% of the free margin %.2f does not cover the minimum volume %.2f of %s",
			freeMargin, info.VolumeMin, symbol)
	}
	if maxLots > info.VolumeMax {
		maxLots = info.VolumeMax
	}

	return maxLots, nil
}

// CanOpenPosition checks if it's possible to open a position with specified volume.
// This performs comprehensive validation: margin check, volume limits, symbol availability
// (IsSymbolAvailable). The margin is calculated for a BUY at the current ASK - the method has no
// direction, so a SELL's margin may differ slightly. Uses 5-second timeout.
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
// and symbol specifications. The margin is calculated for a BUY at the current ASK (the method
// has no direction; a SELL's margin may differ slightly). Uses 5-second timeout.
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
//
// UNITS: the "...Pips" parameters of these helpers are distances in POINTS (1 point = SYMBOL_POINT,
// e.g. 0.00001 on 5-digit EURUSD, 0.001 on 3-digit USDJPY). On 5- and 3-digit symbols 10 points = 1 pip,
// so pass 500 for a 50-pip stop. The names are kept for compatibility.
// ══════════════════════════════════════════════════════════════════════════════

// CalculateSLTP calculates Stop Loss and Take Profit prices from an entry price and point distances.
// The distances are multiplied by the symbol's point size (SYMBOL_POINT) and the results are rounded
// to the symbol's digits. Handles both BUY and SELL directions correctly. Uses the timeouts of
// GetSymbolInfo (10 s) and GetAsk/GetBid (3 s).
//
// PARAMETERS:
//   symbol         - Trading symbol (e.g., "EURUSD")
//   direction      - "BUY" or "SELL" (anything else is an error)
//   entryPrice     - Entry price (use 0 for the current market price: ASK for BUY, BID for SELL)
//   stopLossPips   - Distance to SL in points (e.g., 500 = 50 pips on EURUSD)
//   takeProfitPips - Distance to TP in points (e.g., 1000 = 100 pips on EURUSD)
//
// RETURNS:
//   sl    - Stop Loss price
//   tp    - Take Profit price
//   error - error if calculation fails
//
// EXAMPLE:
//   EURUSD (5 digits) BUY at 1.08500, SL=500 points, TP=1000 points
//   → SL=1.08000, TP=1.09500
func (s *MT5Sugar) CalculateSLTP(symbol, direction string, entryPrice, stopLossPips, takeProfitPips float64) (float64, float64, error) {
	if direction != "BUY" && direction != "SELL" {
		return 0, 0, fmt.Errorf("direction must be \"BUY\" or \"SELL\", got %q", direction)
	}

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

	// Round to the symbol's digits (1.0800000000000001 -> 1.08)
	scale := math.Pow(10, float64(info.Digits))
	return math.Round(sl*scale) / scale, math.Round(tp*scale) / scale, nil
}

// BuyMarketWithPips opens a BUY position with SL/TP given as distances in POINTS from the
// current ASK (see the UNITS note above - the name says pips). The method calculates the
// exact prices with CalculateSLTP and opens the position with BuyMarketWithSLTP (30-second timeout).
//
// PARAMETERS:
//   symbol         - Trading symbol (e.g., "EURUSD")
//   volume         - Lot size (e.g., 0.1)
//   stopLossPips   - Stop Loss distance in points from entry (e.g., 500 = 50 pips on EURUSD)
//   takeProfitPips - Take Profit distance in points from entry (e.g., 1000 = 100 pips on EURUSD)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected
//
// EXAMPLE:
//   ticket, _ := sugar.BuyMarketWithPips("EURUSD", 0.1, 500, 1000)
//   // Opens BUY at market, SL = ASK - 500 points (50 pips), TP = ASK + 1000 points (100 pips)
func (s *MT5Sugar) BuyMarketWithPips(symbol string, volume, stopLossPips, takeProfitPips float64) (uint64, error) {
	// Calculate SL/TP prices
	sl, tp, err := s.CalculateSLTP(symbol, "BUY", 0, stopLossPips, takeProfitPips)
	if err != nil {
		return 0, err
	}

	// Use existing BuyMarketWithSLTP
	return s.BuyMarketWithSLTP(symbol, volume, sl, tp)
}

// SellMarketWithPips opens a SELL position with SL/TP given as distances in POINTS from the
// current BID (see the UNITS note above - the name says pips). The method calculates the
// exact prices with CalculateSLTP and opens the position with SellMarketWithSLTP (30-second timeout).
//
// PARAMETERS:
//   symbol         - Trading symbol (e.g., "EURUSD")
//   volume         - Lot size (e.g., 0.1)
//   stopLossPips   - Stop Loss distance in points from entry (e.g., 500 = 50 pips on EURUSD)
//   takeProfitPips - Take Profit distance in points from entry (e.g., 1000 = 100 pips on EURUSD)
//
// RETURNS:
//   Position ticket number (uint64), or error if order rejected
//
// EXAMPLE:
//   ticket, _ := sugar.SellMarketWithPips("EURUSD", 0.1, 500, 1000)
//   // Opens SELL at market, SL = BID + 500 points (50 pips), TP = BID - 1000 points (100 pips)
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
// account monitoring dashboards or trading reports. Uses one 10-second timeout for the
// whole call (the summary and the margin/free margin/margin level/profit values).
//
// RETURNS:
//   *AccountInfo structure with all account data, or error if any of the queries fails
//   (a failed value is reported, never returned as 0)
func (s *MT5Sugar) GetAccountInfo() (*AccountInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	// Use Service layer GetAccountSummary
	summary, err := s.service.GetAccountSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAccountInfo failed: %w", err)
	}

	// Get additional fields (same context: the whole call has one 10-second deadline)
	margin, err := s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN)
	if err != nil {
		return nil, fmt.Errorf("GetAccountInfo failed: margin: %w", err)
	}
	freeMargin, err := s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_FREE)
	if err != nil {
		return nil, fmt.Errorf("GetAccountInfo failed: free margin: %w", err)
	}
	marginLevel, err := s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_MARGIN_LEVEL)
	if err != nil {
		return nil, fmt.Errorf("GetAccountInfo failed: margin level: %w", err)
	}
	profit, err := s.service.GetAccountDouble(ctx, pb.AccountInfoDoublePropertyType_ACCOUNT_PROFIT)
	if err != nil {
		return nil, fmt.Errorf("GetAccountInfo failed: profit: %w", err)
	}

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

// DailyStats holds trading statistics for today (the closed positions GetDealsToday returns:
// opened since local midnight, see the HISTORY region note).
// Useful for tracking daily performance and generating reports.
//
// FIELDS:
//   TotalDeals   - Number of those closed positions
//   WinningDeals - Number with Profit > 0
//   LosingDeals  - Number with Profit < 0
//   WinRate      - WinningDeals / TotalDeals * 100 (0-100)
//   TotalProfit  - Sum of their Profit (commission, swap and fee not included)
//   BestDeal     - Largest Profit (0 if none was profitable)
//   WorstDeal    - Smallest (most negative) Profit (0 if none lost)
type DailyStats struct {
	TotalDeals   int
	WinningDeals int
	LosingDeals  int
	WinRate      float64
	TotalProfit  float64
	BestDeal     float64
	WorstDeal    float64
}

// GetDailyStats calculates trading statistics for today: the closed positions opened
// between local midnight (time.Local) and now, as GetDealsToday returns them.
// Perfect for daily reports and performance tracking. Uses 15-second timeout.
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
