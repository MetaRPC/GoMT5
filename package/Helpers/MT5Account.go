package mt5

/*
══════════════════════════════════════════════════════════════════════════════
MT5Account - Low-Level MetaTrader 5 gRPC Client
══════════════════════════════════════════════════════════════════════════════

This file implements the low-level MT5 API client with direct protobuf message
handling. All methods accept protobuf Request objects and return protobuf Data.

TOTAL METHODS: 43 (38 unary RPCs + 5 streaming RPCs)

METHOD GROUPS:
──────────────────────────────────────────────────────────────────────────────

1. CONNECTION (6 methods)
   • ConnectEx          - Enhanced connection with extended parameters
   • Connect            - Basic MT5 terminal connection
   • ConnectProxy       - Connect through proxy server
   • CheckConnect       - Verify connection status
   • Disconnect         - Close MT5 connection
   • Reconnect          - Reconnect to MT5 terminal

2. ACCOUNT INFORMATION (4 methods)
   • AccountSummary         - Get all account data in one call (RECOMMENDED)
   • AccountInfoDouble      - Get double properties (Balance, Equity, Margin, etc.)
   • AccountInfoInteger     - Get integer properties (Login, Leverage, etc.)
   • AccountInfoString      - Get string properties (Currency, Company, etc.)

3. SYMBOL INFORMATION & OPERATIONS (14 methods)
   • SymbolsTotal               - Count total/selected symbols
   • SymbolExist                - Check if symbol exists
   • SymbolName                 - Get symbol name by index
   • SymbolSelect               - Add/remove symbol from Market Watch
   • SymbolIsSynchronized       - Check sync status with server
   • SymbolInfoDouble           - Get double properties (Bid, Ask, Point, Volume, etc.)
   • SymbolInfoInteger          - Get integer properties (Digits, Spread, Stops Level)
   • SymbolInfoString           - Get string properties (Description, Base/Profit Currency)
   • SymbolInfoMarginRate       - Get margin requirements for order types
   • SymbolInfoTick             - Get last tick data with timestamp
   • SymbolInfoSessionQuote     - Get quote session times
   • SymbolInfoSessionTrade     - Get trade session times
   • SymbolParamsMany           - Get detailed parameters for multiple symbols
   • TickValueWithSize          - DEPRECATED - use SymbolInfoDouble instead

4. POSITIONS & ORDERS INFORMATION (5 methods)
   • PositionsTotal             - Count open positions
   • OpenedOrders               - Get all opened orders & positions with full details
   • OpenedOrdersTickets        - Get only ticket numbers (lightweight)
   • OrderHistory               - Get historical orders with pagination
   • PositionsHistory           - Get historical positions with P&L

5. MARKET DEPTH / DOM (3 methods)
   • MarketBookAdd      - Subscribe to Depth of Market updates
   • MarketBookRelease  - Unsubscribe from DOM
   • MarketBookGet      - Get current market depth snapshot

6. TRADING OPERATIONS (6 methods)
   • OrderSend          - Send market or pending order
   • OrderModify        - Modify existing order parameters
   • OrderClose         - Close market or pending order
   • OrderCheck         - Validate order before sending
   • OrderCalcMargin    - Calculate required margin
   • OrderCalcProfit    - Calculate potential profit/loss

7. STREAMING METHODS (5 methods) - Real-time data streams
   • OnSymbolTick                           - Stream tick data (Bid/Ask updates)
   • OnTrade                                - Stream trade events
   • OnPositionProfit                       - Stream position P&L updates
   • OnPositionsAndPendingOrdersTickets     - Stream ticket changes
   • OnTradeTransaction                     - Stream trade transaction events

UTILITIES:
   • NewMT5Account              - Create new MT5 account instance
   • Close                      - Close gRPC connection
   • IsConnected                - Check connection status
   • ExecuteWithReconnect       - Generic wrapper for unary RPCs with auto-reconnect
   • ExecuteStreamWithReconnect - Generic wrapper for streaming RPCs with auto-reconnect

══════════════════════════════════════════════════════════════════════════════
*/

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
	"net"
	"strings"

	pb "github.com/MetaRPC/GoMT5/package"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/backoff"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MT5Account represents a low-level gRPC client for MetaTrader 5 terminal.
// All methods accept protobuf Request objects and return protobuf Data objects.
type MT5Account struct {
	User                 uint64
	Password             string
	Host                 string
	Port                 int
	ServerName           string
	BaseChartSymbol      string
	ConnectTimeout       int
	GrpcServer           string
	GrpcConn             *grpc.ClientConn
	AccountInfoData      *pb.AccountSummaryReply
	ConnectionClient         pb.ConnectionClient
	SubscriptionClient       pb.SubscriptionServiceClient
	AccountClient            pb.AccountHelperClient
	AccountInformationClient pb.AccountInformationClient
	TradeClient              pb.TradingHelperClient
	MarketInfoClient         pb.MarketInfoClient
	AccountHelper            pb.AccountHelperClient
	TradeFunctionsClient     pb.TradeFunctionsClient
	HealthClient             pb.HealthClient
	Id                       uuid.UUID
}

type mrpcError interface {
	GetErrorCode() string
}

// NewMT5Account creates a new MT5Account instance with gRPC connection.
// Default grpcServer is "mt5.mrpc.pro:443" if empty string is provided.
// The connection is established with TLS, keepalive, and automatic reconnect configured.
func NewMT5Account(user uint64, password string, grpcServer string, id uuid.UUID) (*MT5Account, error) {
	if grpcServer == "" {
		grpcServer = "mt5.mrpc.pro:443"
	}

	host := grpcServer
	if strings.Contains(host, ":") {
		if h, _, err := net.SplitHostPort(grpcServer); err == nil {
			host = h
		}
	}

	tlsCfg := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}

	if ip := net.ParseIP(host); ip == nil && host != "" {
		tlsCfg.ServerName = host
	}

	dctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bcfg := backoff.Config{
		BaseDelay:  200 * time.Millisecond,
		Multiplier: 1.6,
		Jitter:     0.2,
		MaxDelay:   3 * time.Second,
	}

	kp := keepalive.ClientParameters{
		Time:                20 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}

	conn, err := grpc.DialContext(
		dctx,
		grpcServer,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)),
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           bcfg,
			MinConnectTimeout: 5 * time.Second,
		}),
		grpc.WithKeepaliveParams(kp),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial failed to %s: %w", grpcServer, err)
	}

	return &MT5Account{
		User:                     user,
		Password:                 password,
		GrpcServer:               grpcServer,
		GrpcConn:                 conn,
		ConnectionClient:         pb.NewConnectionClient(conn),
		SubscriptionClient:       pb.NewSubscriptionServiceClient(conn),
		AccountClient:            pb.NewAccountHelperClient(conn),
		AccountInformationClient: pb.NewAccountInformationClient(conn),
		TradeClient:              pb.NewTradingHelperClient(conn),
		MarketInfoClient:         pb.NewMarketInfoClient(conn),
		TradeFunctionsClient:     pb.NewTradeFunctionsClient(conn),
		HealthClient:             pb.NewHealthClient(conn),
		Id:                       id,
		Port:                     443,
		ConnectTimeout:           30,
	}, nil
}

// NewMT5AccountAuto creates a new MT5Account instance with gRPC connection and auto-generated UUID.
// Default grpcServer is "mt5.mrpc.pro:443" if empty string is provided.
// The connection is established with TLS, keepalive, and automatic reconnect configured.
// A random UUID is automatically generated for the session.
func NewMT5AccountAuto(user uint64, password string, grpcServer string) (*MT5Account, error) {
	return NewMT5Account(user, password, grpcServer, uuid.New())
}

// isConnected checks if the account has an active gRPC connection.
func (a *MT5Account) isConnected() bool {
	return a != nil && a.GrpcConn != nil && a.Id != uuid.Nil
}

// getHeaders returns metadata headers with session ID for gRPC calls.
func (a *MT5Account) getHeaders() metadata.MD {
	if !a.isConnected() {
		return nil
	}
	return metadata.Pairs("id", a.Id.String())
}

// Close closes the gRPC connection and cleans up resources.
func (a *MT5Account) Close() error {
	if a == nil {
		return nil
	}
	if a.GrpcConn != nil {
		err := a.GrpcConn.Close()
		a.GrpcConn = nil
		return err
	}
	return nil
}

// IsConnected returns true if the account has an active gRPC connection.
func (a *MT5Account) IsConnected() bool {
	return a != nil && a.GrpcConn != nil && a.Id != uuid.Nil
}

// ExecuteWithReconnect is THE CORE PATTERN used by ALL non-streaming methods in this file.
//
// WHAT THIS DOES:
//   Executes a gRPC call with automatic reconnection on network failures.
//   If connection is lost - attempts to reconnect and retry the request.
//
// ALGORITHM:
//   1. Check gRPC connection
//   2. Add metadata (headers) with session UUID
//   3. Call the passed grpcCall() function
//   4. If network error → exponential backoff + retry
//   5. If API error (TERMINAL_INSTANCE_NOT_FOUND) → reconnect + retry
//   6. Check reply.Error (protobuf errors from MT5)
//   7. Return data or error
//
// WHY THIS IS NEEDED:
//   MT5 Terminal can drop connection (timeout, restart, network issues).
//   This mechanism makes the API resilient to network failures.
//
// RETRY LOGIC:
//   - Initial delay: 500ms
//   - Max delay: 5s
//   - Exponential backoff with jitter
//   - Retries on: Unavailable, DeadlineExceeded, TERMINAL_INSTANCE_NOT_FOUND
func ExecuteWithReconnect[T any](
	a *MT5Account,
	ctx context.Context,
	grpcCall func(metadata.MD) (T, error),
	errorSelector func(T) mrpcError,
) (T, error) {
	var zeroT T
	if ctx == nil {
		ctx = context.Background()
	}

	const (
		initialDelay = 500 * time.Millisecond
		maxDelay     = 5 * time.Second
	)
	delay := initialDelay

	for {
		headers := a.getHeaders()

		res, err := grpcCall(headers)
		if err != nil {
			if s, ok := status.FromError(err); ok && (s.Code() == codes.Unavailable || s.Code() == codes.DeadlineExceeded) {
				log.Printf("[grpc-retry] code=%s msg=%q next_delay=%s", s.Code(), s.Message(), delay)
				j := time.Duration(rand.Int63n(int64(delay/2))) - delay/4
				wait := delay + j
				select {
				case <-time.After(wait):
					delay *= 2
					if delay > maxDelay {
						delay = maxDelay
					}
					continue
				case <-ctx.Done():
					return zeroT, ctx.Err()
				}
			}
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return zeroT, err
			}
			return zeroT, err
		}

		apiErr := errorSelector(res)
		if apiErr != nil && apiErr.GetErrorCode() != "" {
			code := apiErr.GetErrorCode()
			if code == "TERMINAL_INSTANCE_NOT_FOUND" || code == "TERMINAL_REGISTRY_TERMINAL_NOT_FOUND" {
				log.Printf("[api-retry] code=%s next_delay=%s", code, delay)
				j := time.Duration(rand.Int63n(int64(delay/2))) - delay/4
				wait := delay + j
				select {
				case <-time.After(wait):
					delay *= 2
					if delay > maxDelay {
						delay = maxDelay
					}
					continue
				case <-ctx.Done():
					return zeroT, ctx.Err()
				}
			}
			// Convert mrpcError to *pb.Error and wrap in ApiError
			if pbErr, ok := apiErr.(*pb.Error); ok {
				return zeroT, NewApiError(pbErr)
			}
			return zeroT, fmt.Errorf("API error (code=%s): unknown error type", code)
		}

		return res, nil
	}
}

// ExecuteStreamWithReconnect is THE CORE PATTERN used by ALL streaming methods in this file.
//
// WHAT THIS DOES:
//   Executes a streaming gRPC call with automatic stream restart on failures.
//   If stream breaks - automatically reconnects and restarts the stream.
//
// ALGORITHM:
//   1. Create gRPC stream with session UUID in metadata
//   2. Start goroutine that continuously receives messages
//   3. Send data to dataChan, errors to errChan
//   4. If stream error (TERMINAL_INSTANCE_NOT_FOUND, Unavailable) → restart stream
//   5. Apply exponential backoff between retries
//   6. Close channels when context cancelled or stream ends
//
// WHY THIS IS NEEDED:
//   Streaming connections can break due to network issues or MT5 restart.
//   This mechanism ensures continuous data flow by auto-restarting streams.
//
// RETRY LOGIC:
//   - Initial delay: 500ms
//   - Max delay: 5s
//   - Exponential backoff with jitter
//   - Infinite retries until context cancelled
func ExecuteStreamWithReconnect[TRequest any, TReply any, TData any](
	ctx context.Context,
	a *MT5Account,
	request TRequest,
	streamInvoker func(TRequest, metadata.MD, context.Context) (grpc.ClientStream, error),
	getError func(TReply) mrpcError,
	getData func(TReply) (TData, bool),
	newReply func() TReply,
) (<-chan TData, <-chan error) {
	dataCh := make(chan TData)
	errCh := make(chan error, 1)

	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		defer close(dataCh)
		defer close(errCh)

		for {
			reconnectRequired := false
			headers := a.getHeaders()

			stream, err := streamInvoker(request, headers, ctx)
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
					select {
					case <-time.After(500*time.Millisecond + time.Duration(rand.Intn(501)-250)*time.Millisecond):
						continue
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
				}
				errCh <- err
				return
			}

			for {
				reply := newReply()

				recvErr := stream.RecvMsg(reply)
				if recvErr != nil {
					if s, ok := status.FromError(recvErr); ok && s.Code() == codes.Unavailable {
						reconnectRequired = true
						break
					}
					if errors.Is(recvErr, io.EOF) {
						return
					}
					if errors.Is(recvErr, context.Canceled) || errors.Is(recvErr, context.DeadlineExceeded) {
						errCh <- recvErr
						return
					}
					errCh <- recvErr
					return
				}

				apiErr := getError(reply)
				if apiErr != nil && apiErr.GetErrorCode() != "" {
					code := apiErr.GetErrorCode()
					if code == "TERMINAL_INSTANCE_NOT_FOUND" || code == "TERMINAL_REGISTRY_TERMINAL_NOT_FOUND" {
						reconnectRequired = true
						break
					}
					// Convert mrpcError to *pb.Error and wrap in ApiError
					if pbErr, ok := apiErr.(*pb.Error); ok {
						errCh <- NewApiError(pbErr)
					} else {
						errCh <- fmt.Errorf("API error: unknown error type")
					}
					return
				}

				if d, ok := getData(reply); ok {
					select {
					case dataCh <- d:
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
				}
			}

			if reconnectRequired {
				base := 500 * time.Millisecond
				jitter := time.Duration(rand.Intn(501)-250) * time.Millisecond
				select {
				case <-time.After(base + jitter):
					continue
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				}
			}
		}
	}()

	return dataCh, errCh
}

// ══════════════════════════════════════════════════════════════════════════════
// #region CONNECTION
// ══════════════════════════════════════════════════════════════════════════════

// ConnectEx establishes connection to MT5 terminal with extended parameters.
//
// This method provides full control over connection settings including:
//   - MT5 cluster name for connection
//   - Connection timeout (via context.Context)
//   - Base chart symbol selection
//   - Expert Advisors to add
//
// Parameters:
//   - ctx: Context for timeout and cancellation control (timeout replaces old TerminalReadinessWaitingTimeoutSeconds field)
//   - req: ConnectExRequest with User, Password, MtClusterName, BaseChartSymbol, ExpertsToAdd
//
// Returns ConnectData with session UUID and connection status, or error on failure.
func (a *MT5Account) ConnectEx(ctx context.Context, req *pb.ConnectExRequest) (*pb.ConnectData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.ConnectExReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.ConnectionClient.ConnectEx(c, req)
	}

	errorSelector := func(reply *pb.ConnectExReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// Connect establishes basic connection to MT5 terminal.
//
// This is a simplified connection method that uses default settings.
// For advanced connection configuration, use ConnectEx instead.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: ConnectRequest with User, Password, and optional ServerName
//
// Returns ConnectData with session UUID and connection status, or error on failure.
func (a *MT5Account) Connect(ctx context.Context, req *pb.ConnectRequest) (*pb.ConnectData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.ConnectReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.ConnectionClient.Connect(c, req)
	}

	errorSelector := func(reply *pb.ConnectReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// ConnectProxy establishes connection to MT5 terminal through proxy server.
//
// Use this method when MT5 terminal access requires proxy configuration.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: ConnectProxyRequest with proxy host, port, type, credentials and MT5 account credentials
//
// Returns ConnectProxyData with session UUID and connection status, or error on failure.
func (a *MT5Account) ConnectProxy(ctx context.Context, req *pb.ConnectProxyRequest) (*pb.ConnectProxyData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.ConnectProxyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.ConnectionClient.ConnectProxy(c, req)
	}

	errorSelector := func(reply *pb.ConnectProxyReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// CheckConnect verifies the current connection status to MT5 terminal.
//
// Use this method to ping the terminal and confirm the session is still active.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: CheckConnectRequest (empty request structure)
//
// Returns CheckConnectData with connection status flag, or error on failure.
func (a *MT5Account) CheckConnect(ctx context.Context, req *pb.CheckConnectRequest) (*pb.CheckConnectData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.CheckConnectReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.ConnectionClient.CheckConnect(c, req)
	}

	errorSelector := func(reply *pb.CheckConnectReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// Disconnect closes the connection to MT5 terminal.
//
// This method gracefully terminates the active MT5 session.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: DisconnectRequest (empty request structure)
//
// Returns DisconnectData with disconnection status, or error on failure.
func (a *MT5Account) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.DisconnectData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.DisconnectReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.ConnectionClient.Disconnect(c, req)
	}

	errorSelector := func(reply *pb.DisconnectReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// Reconnect re-establishes connection to MT5 terminal.
//
// This method recreates the terminal session without changing connection parameters.
// Used internally by ExecuteWithReconnect on TERMINAL_INSTANCE_NOT_FOUND errors.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: ReconnectRequest (empty request structure)
//
// Returns ReconnectData with new session UUID, or error on failure.
func (a *MT5Account) Reconnect(ctx context.Context, req *pb.ReconnectRequest) (*pb.ReconnectData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.ReconnectReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.ConnectionClient.Reconnect(c, req)
	}

	errorSelector := func(reply *pb.ReconnectReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region ACCOUNT INFORMATION
// ══════════════════════════════════════════════════════════════════════════════

// AccountSummary retrieves all account information in one call.
//
// This is the recommended method for getting account data as it returns all properties
// in a single request, avoiding multiple round-trips.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: AccountSummaryRequest (empty request structure)
//
// Returns AccountSummaryData with Login, Balance, Equity, UserName, Leverage, TradeMode,
// CompanyName, Currency, ServerTime, UtcTimezoneShift, and Credit.
func (a *MT5Account) AccountSummary(ctx context.Context, req *pb.AccountSummaryRequest) (*pb.AccountSummaryData, error) {
	// Step 1: Verify gRPC connection is established
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Step 2: Validate input parameters
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	// Step 3: Setup context with default timeout (10s)
	// If caller didn't provide timeout, we add one to prevent hanging forever
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	// Step 4: Prepare gRPC call function with metadata
	// This closure will be executed by ExecuteWithReconnect with session headers
	grpcCall := func(headers metadata.MD) (*pb.AccountSummaryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.AccountSummary(c, req)
	}

	// Step 5: Define error extraction function
	// ExecuteWithReconnect uses this to check for API errors in the response
	errorSelector := func(reply *pb.AccountSummaryReply) mrpcError {
		return reply.GetError()
	}

	// Step 6: Execute call with automatic retry/reconnect on failure
	// This handles network errors, session expiration, and MT5 terminal restarts
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Step 7: Extract and return data from protobuf response
	return reply.GetData(), nil
}

// AccountInfoDouble retrieves a double-type account property.
//
// Use this method when you need a specific numeric account property.
// For multiple properties, use AccountSummary instead.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: AccountInfoDoubleRequest with property_id (ACCOUNT_BALANCE, ACCOUNT_EQUITY, ACCOUNT_MARGIN, ACCOUNT_MARGIN_FREE, ACCOUNT_PROFIT, etc)
//
// Returns AccountInfoDoubleData with the requested double value.
func (a *MT5Account) AccountInfoDouble(ctx context.Context, req *pb.AccountInfoDoubleRequest) (*pb.AccountInfoDoubleData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.AccountInfoDoubleReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountInformationClient.AccountInfoDouble(c, req)
	}

	errorSelector := func(reply *pb.AccountInfoDoubleReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// AccountInfoInteger retrieves an integer-type account property.
//
// Use this method when you need a specific integer account property.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: AccountInfoIntegerRequest with property_id (ACCOUNT_LOGIN, ACCOUNT_LEVERAGE, ACCOUNT_LIMIT_ORDERS, ACCOUNT_TRADE_MODE, ACCOUNT_MARGIN_SO_MODE, etc)
//
// Returns AccountInfoIntegerData with the requested int64 value.
func (a *MT5Account) AccountInfoInteger(ctx context.Context, req *pb.AccountInfoIntegerRequest) (*pb.AccountInfoIntegerData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.AccountInfoIntegerReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountInformationClient.AccountInfoInteger(c, req)
	}

	errorSelector := func(reply *pb.AccountInfoIntegerReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// AccountInfoString retrieves a string-type account property.
//
// Use this method when you need a specific string account property.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: AccountInfoStringRequest with property_id (ACCOUNT_NAME, ACCOUNT_SERVER, ACCOUNT_CURRENCY, ACCOUNT_COMPANY)
//
// Returns AccountInfoStringData with the requested string value.
func (a *MT5Account) AccountInfoString(ctx context.Context, req *pb.AccountInfoStringRequest) (*pb.AccountInfoStringData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.AccountInfoStringReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountInformationClient.AccountInfoString(c, req)
	}

	errorSelector := func(reply *pb.AccountInfoStringReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region SYMBOL INFORMATION & OPERATIONS
// ══════════════════════════════════════════════════════════════════════════════

// SymbolsTotal returns the number of available symbols.
//
// Use this method to count symbols either in Market Watch or all available symbols.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolsTotalRequest with Selected flag (true for Market Watch symbols only, false for all symbols)
//
// Returns SymbolsTotalData with total count of symbols.
func (a *MT5Account) SymbolsTotal(ctx context.Context, req *pb.SymbolsTotalRequest) (*pb.SymbolsTotalData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolsTotalReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolsTotal(c, req)
	}

	errorSelector := func(reply *pb.SymbolsTotalReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolExist checks if a symbol with specified name exists.
//
// Use this method to verify symbol availability before requesting data or placing orders.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolExistRequest with Symbol name
//
// Returns SymbolExistData with Exist flag and IsCustom flag indicating if it's a custom symbol.
func (a *MT5Account) SymbolExist(ctx context.Context, req *pb.SymbolExistRequest) (*pb.SymbolExistData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolExistReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolExist(c, req)
	}

	errorSelector := func(reply *pb.SymbolExistReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolName returns the name of a symbol by its position in the list.
//
// Use this method to iterate through available symbols by index.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolNameRequest with Pos (zero-based index) and Selected flag (true for Market Watch, false for all symbols)
//
// Returns SymbolNameData with symbol name at the specified position.
func (a *MT5Account) SymbolName(ctx context.Context, req *pb.SymbolNameRequest) (*pb.SymbolNameData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolNameReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolName(c, req)
	}

	errorSelector := func(reply *pb.SymbolNameReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolSelect adds or removes a symbol from Market Watch window.
//
// Use this method to manage which symbols are visible in Market Watch.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolSelectRequest with Symbol name and Select flag (true to add, false to remove)
//
// Returns SymbolSelectData with success status of the operation.
func (a *MT5Account) SymbolSelect(ctx context.Context, req *pb.SymbolSelectRequest) (*pb.SymbolSelectData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolSelectReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolSelect(c, req)
	}

	errorSelector := func(reply *pb.SymbolSelectReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolIsSynchronized checks if symbol data is synchronized with trade server.
//
// Use this method to ensure symbol quotes are up-to-date before trading operations.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolIsSynchronizedRequest with Symbol name
//
// Returns SymbolIsSynchronizedData with IsSynchronized flag.
func (a *MT5Account) SymbolIsSynchronized(ctx context.Context, req *pb.SymbolIsSynchronizedRequest) (*pb.SymbolIsSynchronizedData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolIsSynchronizedReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolIsSynchronized(c, req)
	}

	errorSelector := func(reply *pb.SymbolIsSynchronizedReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolInfoDouble retrieves a double-type symbol property.
//
// Use this method to get numeric symbol properties like prices, volumes, and trading parameters.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolInfoDoubleRequest with Symbol name and PropertyId (BID, ASK, POINT, VOLUME_MIN, VOLUME_MAX, VOLUME_STEP, TRADE_TICK_SIZE, etc)
//
// Returns SymbolInfoDoubleData with the requested double value.
func (a *MT5Account) SymbolInfoDouble(ctx context.Context, req *pb.SymbolInfoDoubleRequest) (*pb.SymbolInfoDoubleData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolInfoDoubleReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolInfoDouble(c, req)
	}

	errorSelector := func(reply *pb.SymbolInfoDoubleReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolInfoInteger retrieves an integer-type symbol property.
//
// Use this method to get integer symbol properties like digits, spread, and trading restrictions.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolInfoIntegerRequest with Symbol name and PropertyId (DIGITS, SPREAD, STOPS_LEVEL, FREEZE_LEVEL, TRADE_MODE, TRADE_EXECUTION_MODE, etc)
//
// Returns SymbolInfoIntegerData with the requested int64 value.
func (a *MT5Account) SymbolInfoInteger(ctx context.Context, req *pb.SymbolInfoIntegerRequest) (*pb.SymbolInfoIntegerData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolInfoIntegerReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolInfoInteger(c, req)
	}

	errorSelector := func(reply *pb.SymbolInfoIntegerReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolInfoString retrieves a string-type symbol property.
//
// Use this method to get text symbol properties like description and currency information.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolInfoStringRequest with Symbol name and PropertyId (DESCRIPTION, CURRENCY_BASE, CURRENCY_PROFIT, CURRENCY_MARGIN, PATH, etc)
//
// Returns SymbolInfoStringData with the requested string value.
func (a *MT5Account) SymbolInfoString(ctx context.Context, req *pb.SymbolInfoStringRequest) (*pb.SymbolInfoStringData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolInfoStringReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolInfoString(c, req)
	}

	errorSelector := func(reply *pb.SymbolInfoStringReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolInfoMarginRate retrieves margin requirements for different order types.
//
// Use this method to calculate margin before placing orders.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolInfoMarginRateRequest with Symbol name and OrderType (ORDER_TYPE_BUY, ORDER_TYPE_SELL, etc)
//
// Returns SymbolInfoMarginRateData with InitialMarginRate and MaintenanceMarginRate values.
func (a *MT5Account) SymbolInfoMarginRate(ctx context.Context, req *pb.SymbolInfoMarginRateRequest) (*pb.SymbolInfoMarginRateData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolInfoMarginRateReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolInfoMarginRate(c, req)
	}

	errorSelector := func(reply *pb.SymbolInfoMarginRateReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolInfoTick retrieves the last tick data for a symbol.
//
// Use this method to get the most recent price update with timestamp.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolInfoTickRequest with Symbol name
//
// Returns MrpcMqlTick with Bid, Ask, Last, Volume, Time, TimeMS, Flags, VolumReal and spread values.
func (a *MT5Account) SymbolInfoTick(ctx context.Context, req *pb.SymbolInfoTickRequest) (*pb.MrpcMqlTick, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolInfoTickRequestReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolInfoTick(c, req)
	}

	errorSelector := func(reply *pb.SymbolInfoTickRequestReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolInfoSessionQuote retrieves quote session times for a symbol.
//
// Use this method to check when quotes are available for trading.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolInfoSessionQuoteRequest with Symbol name, DayOfWeek (SUNDAY=0 to SATURDAY=6) and SessionIndex
//
// Returns SymbolInfoSessionQuoteData with session From and To times in seconds from day start.
func (a *MT5Account) SymbolInfoSessionQuote(ctx context.Context, req *pb.SymbolInfoSessionQuoteRequest) (*pb.SymbolInfoSessionQuoteData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolInfoSessionQuoteReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolInfoSessionQuote(c, req)
	}

	errorSelector := func(reply *pb.SymbolInfoSessionQuoteReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolInfoSessionTrade retrieves trade session times for a symbol.
//
// Use this method to check when trading operations are allowed.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolInfoSessionTradeRequest with Symbol name, DayOfWeek (SUNDAY=0 to SATURDAY=6) and SessionIndex
//
// Returns SymbolInfoSessionTradeData with session From and To times in seconds from day start.
func (a *MT5Account) SymbolInfoSessionTrade(ctx context.Context, req *pb.SymbolInfoSessionTradeRequest) (*pb.SymbolInfoSessionTradeData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolInfoSessionTradeReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolInfoSessionTrade(c, req)
	}

	errorSelector := func(reply *pb.SymbolInfoSessionTradeReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// SymbolParamsMany retrieves detailed parameters for multiple symbols in one call.
//
// This is the recommended method for getting comprehensive symbol data as it returns
// all properties for multiple symbols in a single request, avoiding multiple round-trips.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: SymbolParamsManyRequest with array of Symbol names
//
// Returns SymbolParamsManyData with array of SymbolParams containing Bid, Ask, Digits, Spread,
// VolumeMin, VolumeMax, VolumeStep, ContractSize, Point, margins, and other trading parameters for each symbol.
func (a *MT5Account) SymbolParamsMany(ctx context.Context, req *pb.SymbolParamsManyRequest) (*pb.SymbolParamsManyData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolParamsManyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.SymbolParamsMany(c, req)
	}

	errorSelector := func(reply *pb.SymbolParamsManyReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region POSITIONS & ORDERS INFORMATION
// ══════════════════════════════════════════════════════════════════════════════

// PositionsTotal returns the number of currently open positions.
//
// Use this method for quick check of open positions count without retrieving full details.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//
// Returns PositionsTotalData with Total count of open positions.
func (a *MT5Account) PositionsTotal(ctx context.Context) (*pb.PositionsTotalData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.PositionsTotalReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.PositionsTotal(c, &emptypb.Empty{})
	}

	errorSelector := func(reply *pb.PositionsTotalReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OpenedOrders retrieves all currently opened orders and positions with full details.
//
// This method returns comprehensive information about all active trading positions
// and pending orders including profit/loss, prices, volumes, and timestamps.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OpenedOrdersRequest with InputSortMode (sort type: by open time, close time, or ticket ID)
//
// Returns OpenedOrdersData with arrays of opened_orders (pending orders) and position_infos (open positions) containing full details.
func (a *MT5Account) OpenedOrders(ctx context.Context, req *pb.OpenedOrdersRequest) (*pb.OpenedOrdersData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OpenedOrdersReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.OpenedOrders(c, req)
	}

	errorSelector := func(reply *pb.OpenedOrdersReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OpenedOrdersTickets retrieves only ticket numbers of currently opened orders and positions.
//
// This is a lightweight alternative to OpenedOrders when you only need ticket IDs
// for subsequent operations or monitoring.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OpenedOrdersTicketsRequest (empty request structure)
//
// Returns OpenedOrdersTicketsData with arrays of opened_orders_tickets and opened_position_tickets.
func (a *MT5Account) OpenedOrdersTickets(ctx context.Context, req *pb.OpenedOrdersTicketsRequest) (*pb.OpenedOrdersTicketsData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OpenedOrdersTicketsReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.OpenedOrdersTickets(c, req)
	}

	errorSelector := func(reply *pb.OpenedOrdersTicketsReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OrderHistory retrieves historical orders within a specified time range.
//
// Use this method to analyze past order activity with pagination support
// for large datasets.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderHistoryRequest with FromDate, ToDate (unix timestamps), optional Symbol filter, Offset and Limit for pagination
//
// Returns OrdersHistoryData with array of historical Order objects including execution details,
// prices, volumes, and final status.
func (a *MT5Account) OrderHistory(ctx context.Context, req *pb.OrderHistoryRequest) (*pb.OrdersHistoryData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderHistoryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.OrderHistory(c, req)
	}

	errorSelector := func(reply *pb.OrderHistoryReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// PositionsHistory retrieves closed positions with profit/loss information.
//
// Use this method to analyze trading performance and calculate statistics
// for closed positions within a time range.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: PositionsHistoryRequest with FromDate, ToDate (unix timestamps), optional Symbol filter, Offset and Limit for pagination
//
// Returns PositionsHistoryData with array of closed Position objects including entry/exit prices,
// volumes, swap, commission, net profit, and close timestamps.
func (a *MT5Account) PositionsHistory(ctx context.Context, req *pb.PositionsHistoryRequest) (*pb.PositionsHistoryData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.PositionsHistoryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.PositionsHistory(c, req)
	}

	errorSelector := func(reply *pb.PositionsHistoryReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region MARKET DEPTH / DOM
// ══════════════════════════════════════════════════════════════════════════════

// MarketBookAdd subscribes to Depth of Market (DOM) updates for a symbol.
//
// Use this method to start receiving Level 2 market data with bid/ask prices
// and volumes at different price levels.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: MarketBookAddRequest with Symbol name
//
// Returns MarketBookAddData with subscription status.
func (a *MT5Account) MarketBookAdd(ctx context.Context, req *pb.MarketBookAddRequest) (*pb.MarketBookAddData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.MarketBookAddReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.MarketBookAdd(c, req)
	}

	errorSelector := func(reply *pb.MarketBookAddReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// MarketBookRelease unsubscribes from Depth of Market (DOM) updates.
//
// Use this method to stop receiving Level 2 market data and free resources.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: MarketBookReleaseRequest with Symbol name
//
// Returns MarketBookReleaseData with unsubscription status.
func (a *MT5Account) MarketBookRelease(ctx context.Context, req *pb.MarketBookReleaseRequest) (*pb.MarketBookReleaseData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.MarketBookReleaseReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.MarketBookRelease(c, req)
	}

	errorSelector := func(reply *pb.MarketBookReleaseReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// MarketBookGet retrieves current market depth snapshot for a symbol.
//
// Use this method to get the current order book state with all price levels,
// volumes, and order types (buy/sell) at each level.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: MarketBookGetRequest with Symbol name
//
// Returns MarketBookGetData with array of BookStruct entries containing Type (buy/sell),
// Price, Volume, and VolumeDouble for each price level in the order book.
func (a *MT5Account) MarketBookGet(ctx context.Context, req *pb.MarketBookGetRequest) (*pb.MarketBookGetData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.MarketBookGetReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.MarketBookGet(c, req)
	}

	errorSelector := func(reply *pb.MarketBookGetReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region TRADING OPERATIONS
// ══════════════════════════════════════════════════════════════════════════════

// OrderSend places a market or pending order.
//
// This is the main trading method for opening positions and placing pending orders.
// Supports all order types: market buy/sell, buy/sell limit, buy/sell stop, and stop-limit.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderSendRequest with Symbol, Operation, Volume, Price, Slippage, StopLoss, TakeProfit, Comment, ExpertId, StopLimitPrice, ExpirationTimeType, and ExpirationTime
//
// Returns OrderSendData with returned code, deal ticket, order ticket, execution price, volume, bid/ask prices, comment, and request ID.
func (a *MT5Account) OrderSend(ctx context.Context, req *pb.OrderSendRequest) (*pb.OrderSendData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderSendReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderSend(c, req)
	}

	errorSelector := func(reply *pb.OrderSendReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OrderModify modifies an existing pending order or position.
//
// Use this method to change price levels (entry price for pending orders,
// StopLoss and TakeProfit for positions and pending orders).
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderModifyRequest with Ticket, Price (for pending orders), StopLoss, TakeProfit, and optional Expiration
//
// Returns OrderModifyData with modification status and MqlTradeResult structure.
func (a *MT5Account) OrderModify(ctx context.Context, req *pb.OrderModifyRequest) (*pb.OrderModifyData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderModifyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderModify(c, req)
	}

	errorSelector := func(reply *pb.OrderModifyReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OrderClose closes an existing market position or deletes a pending order.
//
// Use this method to exit positions at current market price or cancel pending orders.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderCloseRequest with Ticket, Volume, and Slippage
//
// Returns OrderCloseData with ReturnedCode, ReturnedStringCode, ReturnedCodeDescription, and CloseMode (market close, partial close, or pending order remove).
func (a *MT5Account) OrderClose(ctx context.Context, req *pb.OrderCloseRequest) (*pb.OrderCloseData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderCloseReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderClose(c, req)
	}

	errorSelector := func(reply *pb.OrderCloseReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OrderCheck validates an order before sending it to the server.
//
// Use this method to pre-validate trading requests without actually placing orders.
// Useful for checking margin requirements and detecting potential errors.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderCheckRequest with Action, Symbol, Volume, Price, OrderType, StopLoss, TakeProfit and other order parameters
//
// Returns OrderCheckData with validation result including margin requirements, estimated profit,
// and MqlTradeCheckResult structure with validation status and possible error codes.
func (a *MT5Account) OrderCheck(ctx context.Context, req *pb.OrderCheckRequest) (*pb.OrderCheckData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderCheckReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.OrderCheck(c, req)
	}

	errorSelector := func(reply *pb.OrderCheckReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OrderCalcMargin calculates required margin for an order.
//
// Use this method to determine how much margin will be required before placing an order.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderCalcMarginRequest with Symbol, OrderType, Volume, and OpenPrice
//
// Returns OrderCalcMarginData with Margin value in account currency.
func (a *MT5Account) OrderCalcMargin(ctx context.Context, req *pb.OrderCalcMarginRequest) (*pb.OrderCalcMarginData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderCalcMarginReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.OrderCalcMargin(c, req)
	}

	errorSelector := func(reply *pb.OrderCalcMarginReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

// OrderCalcProfit calculates potential profit for a trade.
//
// Use this method to estimate profit/loss before placing an order or to calculate
// current profit at a specified price level.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control
//   - req: OrderCalcProfitRequest with OrderType, Symbol, Volume, OpenPrice, and ClosePrice
//
// Returns OrderCalcProfitData with Profit value in account currency.
func (a *MT5Account) OrderCalcProfit(ctx context.Context, req *pb.OrderCalcProfitRequest) (*pb.OrderCalcProfitData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderCalcProfitReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.OrderCalcProfit(c, req)
	}

	errorSelector := func(reply *pb.OrderCalcProfitReply) mrpcError {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}
// #endregion

// ══════════════════════════════════════════════════════════════════════════════
// #region STREAMING METHODS
// ══════════════════════════════════════════════════════════════════════════════

// OnSymbolTick streams real-time tick data for a symbol.
//
// This method provides continuous price updates (Bid/Ask) as they arrive from the server.
// The stream automatically reconnects on connection loss.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control (cancel to stop streaming)
//   - req: OnSymbolTickRequest with Symbol name
//
// Returns two channels:
//   - Data channel: receives OnSymbolTickData with Bid, Ask, Last, Volume, Time for each tick
//   - Error channel: receives errors if stream fails (both channels closed on context cancellation)
func (a *MT5Account) OnSymbolTick(ctx context.Context, req *pb.OnSymbolTickRequest) (<-chan *pb.OnSymbolTickData, <-chan error) {
	streamInvoker := func(request *pb.OnSymbolTickRequest, headers metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.SubscriptionClient.OnSymbolTick(c, request)
	}

	getError := func(reply *pb.OnSymbolTickReply) mrpcError {
		return reply.GetError()
	}

	getData := func(reply *pb.OnSymbolTickReply) (*pb.OnSymbolTickData, bool) {
		if data := reply.GetData(); data != nil {
			return data, true
		}
		return nil, false
	}

	newReply := func() *pb.OnSymbolTickReply {
		return &pb.OnSymbolTickReply{}
	}

	return ExecuteStreamWithReconnect(ctx, a, req, streamInvoker, getError, getData, newReply)
}

// OnTrade streams trade events in real-time.
//
// This method provides notifications about all trading operations: order placement,
// modification, execution, and cancellation.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control (cancel to stop streaming)
//   - req: OnTradeRequest (empty request structure)
//
// Returns two channels:
//   - Data channel: receives OnTradeData with trade event details
//   - Error channel: receives errors if stream fails (both channels closed on context cancellation)
func (a *MT5Account) OnTrade(ctx context.Context, req *pb.OnTradeRequest) (<-chan *pb.OnTradeData, <-chan error) {
	streamInvoker := func(request *pb.OnTradeRequest, headers metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.SubscriptionClient.OnTrade(c, request)
	}

	getError := func(reply *pb.OnTradeReply) mrpcError {
		return reply.GetError()
	}

	getData := func(reply *pb.OnTradeReply) (*pb.OnTradeData, bool) {
		if data := reply.GetData(); data != nil {
			return data, true
		}
		return nil, false
	}

	newReply := func() *pb.OnTradeReply {
		return &pb.OnTradeReply{}
	}

	return ExecuteStreamWithReconnect(ctx, a, req, streamInvoker, getError, getData, newReply)
}

// OnPositionProfit streams real-time profit/loss updates for open positions.
//
// This method provides continuous P&L updates as market prices change,
// useful for monitoring account performance and implementing risk management.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control (cancel to stop streaming)
//   - req: OnPositionProfitRequest with optional Symbol filter
//
// Returns two channels:
//   - Data channel: receives OnPositionProfitData with Ticket, Symbol, Profit, and current price
//   - Error channel: receives errors if stream fails (both channels closed on context cancellation)
func (a *MT5Account) OnPositionProfit(ctx context.Context, req *pb.OnPositionProfitRequest) (<-chan *pb.OnPositionProfitData, <-chan error) {
	streamInvoker := func(request *pb.OnPositionProfitRequest, headers metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.SubscriptionClient.OnPositionProfit(c, request)
	}

	getError := func(reply *pb.OnPositionProfitReply) mrpcError {
		return reply.GetError()
	}

	getData := func(reply *pb.OnPositionProfitReply) (*pb.OnPositionProfitData, bool) {
		if data := reply.GetData(); data != nil {
			return data, true
		}
		return nil, false
	}

	newReply := func() *pb.OnPositionProfitReply {
		return &pb.OnPositionProfitReply{}
	}

	return ExecuteStreamWithReconnect(ctx, a, req, streamInvoker, getError, getData, newReply)
}

// OnPositionsAndPendingOrdersTickets streams changes in open positions and pending orders.
//
// This method notifies whenever positions are opened/closed or pending orders are added/removed,
// providing only ticket numbers for efficient monitoring.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control (cancel to stop streaming)
//   - req: OnPositionsAndPendingOrdersTicketsRequest (empty request structure)
//
// Returns two channels:
//   - Data channel: receives OnPositionsAndPendingOrdersTicketsData with arrays of PositionTickets and PendingOrderTickets
//   - Error channel: receives errors if stream fails (both channels closed on context cancellation)
func (a *MT5Account) OnPositionsAndPendingOrdersTickets(ctx context.Context, req *pb.OnPositionsAndPendingOrdersTicketsRequest) (<-chan *pb.OnPositionsAndPendingOrdersTicketsData, <-chan error) {
	streamInvoker := func(request *pb.OnPositionsAndPendingOrdersTicketsRequest, headers metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.SubscriptionClient.OnPositionsAndPendingOrdersTickets(c, request)
	}

	getError := func(reply *pb.OnPositionsAndPendingOrdersTicketsReply) mrpcError {
		return reply.GetError()
	}

	getData := func(reply *pb.OnPositionsAndPendingOrdersTicketsReply) (*pb.OnPositionsAndPendingOrdersTicketsData, bool) {
		if data := reply.GetData(); data != nil {
			return data, true
		}
		return nil, false
	}

	newReply := func() *pb.OnPositionsAndPendingOrdersTicketsReply {
		return &pb.OnPositionsAndPendingOrdersTicketsReply{}
	}

	return ExecuteStreamWithReconnect(ctx, a, req, streamInvoker, getError, getData, newReply)
}

// OnTradeTransaction streams detailed trade transaction events.
//
// This method provides low-level notifications about every change in trading state,
// including order state changes, deal executions, and position modifications.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control (cancel to stop streaming)
//   - req: OnTradeTransactionRequest (empty request structure)
//
// Returns two channels:
//   - Data channel: receives OnTradeTransactionData with MqlTradeTransaction containing Type, OrderState, DealTicket, OrderTicket, Symbol, Price, Volume
//   - Error channel: receives errors if stream fails (both channels closed on context cancellation)
func (a *MT5Account) OnTradeTransaction(ctx context.Context, req *pb.OnTradeTransactionRequest) (<-chan *pb.OnTradeTransactionData, <-chan error) {
	streamInvoker := func(request *pb.OnTradeTransactionRequest, headers metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.SubscriptionClient.OnTradeTransaction(c, request)
	}

	getError := func(reply *pb.OnTradeTransactionReply) mrpcError {
		return reply.GetError()
	}

	getData := func(reply *pb.OnTradeTransactionReply) (*pb.OnTradeTransactionData, bool) {
		if data := reply.GetData(); data != nil {
			return data, true
		}
		return nil, false
	}

	newReply := func() *pb.OnTradeTransactionReply {
		return &pb.OnTradeTransactionReply{}
	}

	return ExecuteStreamWithReconnect(ctx, a, req, streamInvoker, getError, getData, newReply)
}
// #endregion
