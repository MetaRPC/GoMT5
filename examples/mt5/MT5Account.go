package mt5

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

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/backoff"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MT5Account represents a client session for interacting with the MT4 terminal API over gRPC.
type MT5Account struct {

	// User is the MT5 account login number.
	User uint64

	// Password for the user account.
	Password string

	// Host is the IP/domain of the MT5 server.
	Host string

	// Port is the MT5 server port (typically 443).
	Port int

	// ServerName is the MT5 server name (used for cluster connections).
	ServerName string

	// BaseChartSymbol is the default chart symbol (e.g., "EURUSD").
	BaseChartSymbol string

	// ConnectTimeout is the timeout for connection readiness, in seconds.
	ConnectTimeout int

	// GrpcServer is the address (host:port) of the gRPC API endpoint.
	GrpcServer string

	// GrpcConn is the underlying gRPC client connection (must be closed when done).
	GrpcConn *grpc.ClientConn

	// Per-service gRPC API clients.
	AccountInfoData      *pb.AccountSummaryReply
	ConnectionClient     pb.ConnectionClient
	SubscriptionClient   pb.SubscriptionServiceClient
	AccountClient        pb.AccountHelperClient
	TradeClient          pb.TradingHelperClient
	MarketInfoClient     pb.MarketInfoClient
	AccountHelper        pb.AccountHelperClient
	TradeFunctionsClient pb.TradeFunctionsClient
	HealthClient         pb.HealthClient

	// Id is a unique identifier (UUID) for this account session/instance.
	Id uuid.UUID
}
type mrpcError interface {
	GetErrorCode() string
}

// NewMT5Account initializes a new MT4Account and establishes the underlying gRPC connection.
// Returns a pointer to the account object and any error encountered while connecting.
func NewMT5Account(user uint64, password string, grpcServer string, id uuid.UUID) (*MT5Account, error) {
	// Default endpoint
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
		// ServerName: "mt5.mrpc.pro",
	}

	if ip := net.ParseIP(host); ip == nil && host != "" {
    tlsCfg.ServerName = host
}
	// Blocking dial with timeout
	dctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
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
		User:                 user,
		Password:             password,
		GrpcServer:           grpcServer,
		GrpcConn:             conn,
		ConnectionClient:     pb.NewConnectionClient(conn),
		SubscriptionClient:   pb.NewSubscriptionServiceClient(conn),
		AccountClient:        pb.NewAccountHelperClient(conn),
		TradeClient:          pb.NewTradingHelperClient(conn),
		MarketInfoClient:     pb.NewMarketInfoClient(conn),
		TradeFunctionsClient: pb.NewTradeFunctionsClient(conn),
		HealthClient:         pb.NewHealthClient(conn),
		Id:                   id,
		Port:                 443,
		ConnectTimeout:       30,
	}, nil
}


// isConnected returns true if this account is associated with any host or server name.
func (a *MT5Account) isConnected() bool {
	// Connected — only if there is a live gRPC channel and a valid Terminal Id
	return a != nil && a.GrpcConn != nil && a.Id != uuid.Nil
}

func (a *MT5Account) getHeaders() metadata.MD {
	// We send the id only during a real session
	if !a.isConnected() {
		return nil
	}
	return metadata.Pairs("id", a.Id.String())
}
// CHANGED: Close() now resets the reference to conn after closing.
func (a *MT5Account) Close() error {
	if a == nil {
		return nil
	}
	if a.GrpcConn != nil {
		err := a.GrpcConn.Close()
		a.GrpcConn = nil // CHANGED
		return err
	}
	return nil
}

// Disconnect closes the connection and resets the main fields
func (a *MT5Account) Disconnect(ctx context.Context) error {
	if a == nil {
		return errors.New("account is nil")
	}

	// Close a low-level connection
	closeErr := a.Close()

	// Reset gRPC clients
	a.ConnectionClient = nil
	a.SubscriptionClient = nil
	a.AccountClient = nil
	a.TradeClient = nil
	a.MarketInfoClient = nil
	a.AccountHelper = nil
	a.TradeFunctionsClient = nil
	a.HealthClient = nil

	// Reset the runtime state (do not touch credits)
	a.Id = uuid.Nil
	a.Host = ""
	a.ServerName = ""
	a.BaseChartSymbol = ""
	a.ConnectTimeout = 0

	return closeErr
}
func (a *MT5Account) IsConnected() bool {
	return a != nil && a.GrpcConn != nil && a.Id != uuid.Nil
}

// ConnectByProxy connects via proxy (socks5/http) using Connection.ConnectProxy.
func (s *MT5Service) ConnectByProxy(
	ctx context.Context,
	user uint64,
	password, host string,
	port int32,
	proxyUser, proxyPassword, proxyHost string,
	proxyPort int32,
	proxyType pb.ProxyTypes, // correct enum
	baseChartSymbol string,
	waitForTerminalIsAlive bool,
	timeoutSeconds int32,
) error {
	if s.account == nil {
		return fmt.Errorf("MT5 account not initialized")
	}

	req := &pb.ConnectProxyRequest{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,

		ProxyUser:     proxyUser,
		ProxyPassword: proxyPassword,
		ProxyHost:     proxyHost,
		ProxyPort:     uint32(proxyPort),
		ProxyType:     proxyType,

		BaseChartSymbol:                        &baseChartSymbol,
		WaitForTerminalIsAlive:                 &waitForTerminalIsAlive,
		TerminalReadinessWaitingTimeoutSeconds: &timeoutSeconds,
	}

	call := func(h metadata.MD) (*pb.ConnectProxyReply, error) {
		c := metadata.NewOutgoingContext(ctx, h)
		return s.account.ConnectionClient.ConnectProxy(c, req)
	}

	errSel := func(r *pb.ConnectProxyReply) mrpcError { return r.GetError() }

	reply, err := ExecuteWithReconnect(s.account, ctx, call, errSel)
	if err != nil {
		return err
	}

	if d := reply.GetData(); d != nil && d.GetUniqueIdentifier() != "" {
		if id, e := uuid.Parse(d.GetUniqueIdentifier()); e == nil {
			s.account.Id = id
		} else {

			log.Printf("warn: cannot parse UniqueIdentifier as UUID: %v", e)
		}
	}

	// Persist basic connection parameters
	s.account.User = user
	s.account.Password = password
	s.account.Host = host
	s.account.Port = int(port)
	s.account.BaseChartSymbol = baseChartSymbol
	s.account.ConnectTimeout = int(timeoutSeconds)
	return nil
}

// ShowCheckConnect prints terminal connection status via Connection.CheckConnect.
func (s *MT5Service) ShowCheckConnect(ctx context.Context) {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // health check should be fast
		defer cancel()
	}

	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}

	req := &pb.CheckConnectRequest{}

	call := func(h metadata.MD) (*pb.CheckConnectReply, error) {
		c := metadata.NewOutgoingContext(ctx, h) // attach metadata to normalized ctx
		return s.account.ConnectionClient.CheckConnect(c, req)
	}
	errSel := func(r *pb.CheckConnectReply) mrpcError { return r.GetError() }

	reply, err := ExecuteWithReconnect(s.account, ctx, call, errSel)
	if err != nil {
		log.Printf("❌ CheckConnect error: %v", err)
		return
	}

	d := reply.GetData()
	if d == nil {
		fmt.Println("⚠️ CheckConnect returned empty data")
		return
	}
	hc := d.GetHealthCheck()
	if hc == nil {
		fmt.Printf("✅ CheckConnect: uid=%s | (no health payload)\n", d.GetUniqueIdentifier())
		return
	}

	// Optional strings are wrappers; guard nil.
	var msg, code string
	if hc.GetErrorMessage() != nil {
		msg = hc.GetErrorMessage().GetValue()
	}
	if hc.GetErrorCode() != nil {
		code = hc.GetErrorCode().GetValue()
	}

	fmt.Printf(
		"✅ CheckConnect: uid=%s | isAlive=%v | apiAlive=%v | mtConnected=%v | authErr=%v | errCode=%s | errMsg=%s\n",
		d.GetUniqueIdentifier(),
		hc.GetIsAlive(),
		hc.GetApiIsAlive(),
		hc.GetTerminalIsConnectedToMtServer(),
		hc.GetHasAuthorizationError(),
		code, msg,
	)
}


// Reconnect triggers Connection.Reconnect with per-call timeout and ctx normalization.
func (s *MT5Service) Reconnect(ctx context.Context) error {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second) // reconnect can take a moment
		defer cancel()
	}

	if s.account == nil {
		return fmt.Errorf("MT5 account not initialized")
	}

	req := &pb.ReconnectRequest{}
	call := func(h metadata.MD) (*pb.ReconnectReply, error) {
		c := metadata.NewOutgoingContext(ctx, h) // attach metadata to normalized ctx
		return s.account.ConnectionClient.Reconnect(c, req)
	}
	errSel := func(r *pb.ReconnectReply) mrpcError { return r.GetError() }

	_, err := ExecuteWithReconnect(s.account, ctx, call, errSel)
	return err
}


// ShowHealthCheck prints result of Health.Check (service liveness).
func (s *MT5Service) ShowHealthCheck(ctx context.Context) {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // health check should be quick
		defer cancel()
	}

	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}

	req := &pb.HealthCheckRequest{}
	call := func(h metadata.MD) (*pb.HealthCheckReply, error) {
		c := metadata.NewOutgoingContext(ctx, h) // attach metadata to normalized ctx
		return s.account.HealthClient.Check(c, req)
	}
	// Health.Check uses transport status for errors; no app-level error expected.
	errSel := func(_ *pb.HealthCheckReply) mrpcError { return nil }

	reply, err := ExecuteWithReconnect(s.account, ctx, call, errSel)
	if err != nil {
		log.Printf("❌ Health.Check error: %v", err)
		return
	}
	if reply == nil {
		fmt.Println("⚠️ Health.Check returned nil reply")
		return
	}

	fmt.Printf("✅ Health: isConnectedToServer=%v | serverTime=%d\n",
		reply.GetIsConnectedToServer(), reply.GetServerTimeSeconds())
}


// StopListening requests the Health service to stop any active server-side listeners
// associated with the current session. It normalizes the context and applies a short
// per-call timeout. Note: this does not close the underlying gRPC connection.
func (s *MT5Service) StopListening(ctx context.Context) error {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // quick control operation
		defer cancel()
	}

	if s.account == nil {
		return fmt.Errorf("MT5 account not initialized")
	}

	req := &pb.StopListeningRequest{}

	call := func(h metadata.MD) (*pb.StopListeningReply, error) {
		c := metadata.NewOutgoingContext(ctx, h) // attach metadata to normalized ctx
		return s.account.HealthClient.StopListening(c, req)
	}
	// No application-level error expected from StopListening; rely on transport status.
	errSel := func(_ *pb.StopListeningReply) mrpcError { return nil }

	_, err := ExecuteWithReconnect(s.account, ctx, call, errSel)
	return err
}


// IsTerminalAlive checks if the terminal session is active and responsive
// by issuing a lightweight AccountSummary call with a short timeout.
func (a *MT5Account) IsTerminalAlive(ctx context.Context) (bool, error) {
	// Fast pre-check: no session → no RPC
	if !a.isConnected() {
		return false, errors.New("not connected")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second) // quick liveness probe
		defer cancel()
	}

	req := &pb.AccountSummaryRequest{}

	grpcCall := func(headers metadata.MD) (*pb.AccountSummaryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.AccountClient.AccountSummary(c, req)
	}

	errorSelector := func(reply *pb.AccountSummaryReply) mrpcError {
		return reply.GetError()
	}

	_, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return false, err
	}
	return true, nil
}


// ConnectByHostPort connects to the MT5 terminal using a host/port pair.
// Updates the session state fields (Host, Port, etc.) upon success.
func (a *MT5Account) ConnectByHostPort(
	ctx context.Context,
	host string,
	port int,
	baseChartSymbol string,
	waitForTerminalIsAlive bool,
	timeoutSeconds int,
) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		// Use a sane default; if caller provided a larger server-side wait, prefer it.
		d := 15 * time.Second
		if timeoutSeconds > 0 {
			td := time.Duration(timeoutSeconds) * time.Second
			if td > d {
				d = td
			}
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d)
		defer cancel()
	}

	req := &pb.ConnectRequest{
		User:                                   a.User,
		Password:                               a.Password,
		Host:                                   host,
		Port:                                   int32(port),
		BaseChartSymbol:                        proto.String(baseChartSymbol),
		WaitForTerminalIsAlive:                 proto.Bool(waitForTerminalIsAlive),
		TerminalReadinessWaitingTimeoutSeconds: proto.Int32(int32(timeoutSeconds)),
	}

	md := a.getHeaders()
	c := metadata.NewOutgoingContext(ctx, md)

	res, err := a.ConnectionClient.Connect(c, req)
	if err != nil {
		return err
	}
	if res.GetError() != nil {
		return fmt.Errorf("API error: %v", res.GetError())
	}

	// Persist session properties
	a.Host = host
	a.Port = port
	a.BaseChartSymbol = baseChartSymbol
	a.ConnectTimeout = timeoutSeconds

	// Set session UUID if present
	if data := res.GetData(); data != nil && data.GetTerminalInstanceGuid() != "" {
		if id, parseErr := uuid.Parse(data.GetTerminalInstanceGuid()); parseErr == nil {
			a.Id = id
		}
	}

	return nil
}


// ConnectByServerName connects to the MT5 terminal using the cluster/server name.
// Updates session fields on success.
func (a *MT5Account) ConnectByServerName(
	ctx context.Context,
	serverName string,
	baseChartSymbol string,
	waitForTerminalIsAlive bool, // kept in signature for symmetry; not used by ConnectEx
	timeoutSeconds int,
) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		d := 15 * time.Second
		if timeoutSeconds > 0 {
			if td := time.Duration(timeoutSeconds) * time.Second; td > d {
				d = td
			}
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d)
		defer cancel()
	}

	req := &pb.ConnectExRequest{
		User:                                   a.User,
		Password:                               a.Password,
		MtClusterName:                          serverName,
		BaseChartSymbol:                        proto.String(baseChartSymbol),
		TerminalReadinessWaitingTimeoutSeconds: proto.Int32(int32(timeoutSeconds)),
		// NOTE: no WaitForTerminalIsAlive in ConnectExRequest
	}

	md := a.getHeaders()
	c := metadata.NewOutgoingContext(ctx, md)

	res, err := a.ConnectionClient.ConnectEx(c, req)
	if err != nil {
		return err
	}
	if res.GetError() != nil {
		return fmt.Errorf("API error: %v", res.GetError())
	}

	// Persist session properties
	a.ServerName = serverName
	a.BaseChartSymbol = baseChartSymbol
	a.ConnectTimeout = timeoutSeconds

	// Set session UUID if present
	if data := res.GetData(); data != nil && data.GetTerminalInstanceGuid() != "" {
		if id, parseErr := uuid.Parse(data.GetTerminalInstanceGuid()); parseErr == nil {
			a.Id = id
		}
	}

	return nil
}



// ExecuteWithReconnect retries a gRPC call on recoverable errors (network/instance-not-found).
//
// T:          Type of the response object (e.g., *pb.AccountSummaryReply)
// a:          Pointer to MT4Account (used for headers, etc.)
// ctx:        Context for cancellation/deadline
// grpcCall:   Function taking metadata and returning (response, error)
// errorSelector: Function taking response and extracting *pb.Error (returns nil if no API error)
//
// Returns the response or error.
// ExecuteWithReconnect performs a unary gRPC call with transparent reconnect/retry.
// - Retries on gRPC Unavailable
// - Retries on API "TERMINAL_INSTANCE_NOT_FOUND"
// - Returns the reply or a wrapped error
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
			// ✳️ Logging the gRPC status before the reset
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
			// Other errors — immediately out
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return zeroT, err
			}
			return zeroT, err
		}

		apiErr := errorSelector(res)
		if apiErr != nil {
			code := apiErr.GetErrorCode()
			// ✳️ The terminal has not been registered yet - soft delay
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
			// ✳️ Explicitly logging the API error code
			return zeroT, fmt.Errorf("API error (code=%s): %v", code, apiErr)
		}

		return res, nil
	}
}



//=== 📂 Account Info ===

// AccountSummary returns high-level metrics of the connected MT5 account
// (balance, equity, leverage, currency, etc.).
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//
// Returns:
//   - *pb.AccountSummaryData on success (nil on error).
//   - error if not connected or the RPC fails.
//
// Behavior:
//   - Adds session headers (id) if available.
//   - Retries automatically on transient transport errors and
//     "terminal instance not found" using ExecuteWithReconnect.
//   - Accesses fields via generated Get...() accessors.
//
// Notes:
//   - This call does not mutate any server state.
//   - Safe to invoke frequently (lightweight read).

func (a *MT5Account) AccountSummary(ctx context.Context) (*pb.AccountSummaryData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// Require an active session.
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Request + call wrapper.
	req := &pb.AccountSummaryRequest{}
	grpcCall := func(headers metadata.MD) (*pb.AccountSummaryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.AccountSummary(c, req)
	}
	errorSelector := func(reply *pb.AccountSummaryReply) mrpcError { return reply.GetError() }

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	return reply.GetData(), nil
}



// NOTE: All values are retrieved via AccountSummary,
// which already uses ExecuteWithReconnect internally.

// AccountLogin returns the account login number.
func (a *MT5Account) AccountLogin(ctx context.Context) (int64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountLogin(), nil
}

// AccountBalance returns the account balance.
func (a *MT5Account) AccountBalance(ctx context.Context) (float64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountBalance(), nil
}

// AccountCredit returns the account credit.
func (a *MT5Account) AccountCredit(ctx context.Context) (float64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountCredit(), nil
}

// AccountEquity returns the account equity.
func (a *MT5Account) AccountEquity(ctx context.Context) (float64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountEquity(), nil
}

// AccountLeverage returns the account leverage.
func (a *MT5Account) AccountLeverage(ctx context.Context) (int64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountLeverage(), nil
}

// AccountName returns the account user name.
func (a *MT5Account) AccountName(ctx context.Context) (string, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return "", err
	}
	return summary.GetAccountUserName(), nil
}

// AccountCompany returns the account company name.
func (a *MT5Account) AccountCompany(ctx context.Context) (string, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return "", err
	}
	return summary.GetAccountCompanyName(), nil
}

// AccountCurrency returns the account currency code.
func (a *MT5Account) AccountCurrency(ctx context.Context) (string, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return "", err
	}
	return summary.GetAccountCurrency(), nil
}

// ExecuteStreamWithReconnect wraps a gRPC server-streaming call with automatic reconnection
// on network and recoverable API errors, sending extracted data to a channel.
// - ctx: Context for cancellation and deadline
// - a:   Your session/account struct (for session headers etc.)
// - request: The protobuf request message
// - streamInvoker: function to open the gRPC stream (returns grpc.ClientStream, error)
// - getError: function to extract API error from a reply message (returns *pb.Error or nil)
// - getData: function to extract (data, ok) from a reply (ok=false means skip this message)
// - newReply: function that creates a new TReply instance (needed because Go generics can't new(T))
// Returns:
//   - dataCh: channel of received TData messages (e.g. *pb.OnSymbolTickData)
//   - errCh:  channel for any errors or end-of-stream events.
//
// ExecuteStreamWithReconnect handles server-streaming calls with transparent reconnect.
// It will:
//   - Reconnect on gRPC Unavailable
//   - Reconnect on API TERMINAL_INSTANCE_NOT_FOUND / TERMINAL_REGISTRY_TERMINAL_NOT_FOUND
//   - Forward each parsed data item to dataCh; send a single error to errCh on fatal errors
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

	// guard ctx
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		defer close(dataCh)
		defer close(errCh)

		for {
			reconnectRequired := false
			headers := a.getHeaders()

			// Open stream
			stream, err := streamInvoker(request, headers, ctx)
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
					select {
					case <-time.After(500 * time.Millisecond + time.Duration(rand.Intn(501)-250)*time.Millisecond):
						continue // retry
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
				}
				errCh <- err
				return
			}

			// Read loop
			for {
				reply := newReply()

				recvErr := stream.RecvMsg(reply)
				if recvErr != nil {
					// gRPC-level transient error -> reconnect
					if s, ok := status.FromError(recvErr); ok && s.Code() == codes.Unavailable {
						reconnectRequired = true
						break
					}
					// Normal EOF: stream ended (do not reconnect)
					if errors.Is(recvErr, io.EOF) {
						return
					}
					// Context errors: just stop
					if errors.Is(recvErr, context.Canceled) || errors.Is(recvErr, context.DeadlineExceeded) {
						errCh <- recvErr
						return
					}
					// Any other error: stop
					errCh <- recvErr
					return
				}

				// API-layer error inside the reply
				apiErr := getError(reply)
				if apiErr != nil {
					code := apiErr.GetErrorCode()
					if code == "TERMINAL_INSTANCE_NOT_FOUND" || code == "TERMINAL_REGISTRY_TERMINAL_NOT_FOUND" {
						reconnectRequired = true
						break
					}
					errCh <- fmt.Errorf("API error: %v", apiErr)
					return
				}

				// Extract data and forward
				if d, ok := getData(reply); ok {
					select {
					case dataCh <- d:
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
				}
			}

			// Reconnect if needed
			if reconnectRequired {
				base := 500 * time.Millisecond
				jitter := time.Duration(rand.Intn(501)-250) * time.Millisecond // [-250ms, +250ms]
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


//=== 📂 Order Operations ===

// OrderSend places a new trade order (market or pending) on the connected MT5 terminal.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - symbol: Trading symbol (e.g., "EURUSD").
//   - operationType: Order operation (Buy/Sell/BuyLimit/...).
//   - volume: Trade volume in lots (e.g., 0.10, 1.00).
//   - price: Entry price for pending orders (ignored for market orders).
//   - slippage: Max slippage in points (market orders only, optional).
//   - stoploss: Stop Loss price (optional).
//   - takeprofit: Take Profit price (optional).
//   - comment: Optional order comment.
//   - magicNumber: Expert/EA identifier (optional).
//   - expiration: Expiration for pending orders (optional; ignored by market).
//
// Returns:
//   - *pb.OrderSendData with result details (e.g., status, created order info) or nil on error.
//   - error if not connected or the RPC/API call fails.
//
// Behavior:
//   - Adds session headers (id) when available.
//   - Uses ExecuteWithReconnect to retry on transient gRPC errors and
//     "terminal instance not found" responses.
//   - Uses generated Get...() accessors to read protobuf fields.
//
// Notes:
//   - STOP_LIMIT orders require stopLimitPrice; use OrderSendEx/OrderSendStopLimit instead.
//   - This call does not mutate terminal state beyond creating the order itself.
func (a *MT5Account) OrderSend(
	ctx context.Context,
	symbol string,
	operationType pb.TMT5_ENUM_ORDER_TYPE,
	volume float64,
	price *float64,
	slippage *int32, // input stays *int32 for caller convenience
	stoploss *float64,
	takeprofit *float64,
	comment *string,
	magicNumber *int32, // input stays *int32
	expiration *timestamppb.Timestamp,
) (*pb.OrderSendData, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if operationType == pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP_LIMIT ||
		operationType == pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP_LIMIT {
		return nil, fmt.Errorf("use OrderSendEx/OrderSendStopLimit for STOP_LIMIT (requires stopLimitPrice)")
	}

	// per-call deadline (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: operationType,
		Volume:    volume,
	}
	if price != nil {
		req.Price = price
	}
	// exact proto types:
	req.Slippage = u64pFromI32(slippage)
	if stoploss != nil {
		req.StopLoss = stoploss
	}
	if takeprofit != nil {
		req.TakeProfit = takeprofit
	}
	if comment != nil {
		req.Comment = comment
	}
	req.ExpertId = u64pFromI32(magicNumber)

	if expiration != nil {
		t := pb.TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_SPECIFIED
		req.ExpirationTimeType = &t
		req.ExpirationTime = expiration
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderSendReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderSend(c, req)
	}
	errorSelector := func(reply *pb.OrderSendReply) mrpcError { return reply.GetError() }

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	return reply.GetData(), nil
}


// OrderSendEx places a new trade order (market, pending, or STOP_LIMIT) on the connected MT5 terminal.
// Similar to OrderSend, but adds an extra stopLimitPrice parameter for STOP_LIMIT orders.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - symbol: Trading symbol (e.g., "EURUSD").
//   - operationType: Order operation (Buy/Sell/BuyLimit/StopLimit/...).
//   - volume: Trade volume in lots (e.g., 0.10, 1.00).
//   - price: Entry price for pending orders (trigger price for STOP_LIMIT).
//   - slippage: Max allowed slippage in points (market orders only, optional).
//   - stoploss: Stop Loss price (optional).
//   - takeprofit: Take Profit price (optional).
//   - comment: Optional order comment.
//   - magicNumber: Expert/EA identifier (optional).
//   - expiration: Expiration for pending orders (optional; ignored by market orders).
//   - stopLimitPrice: Limit price after trigger for STOP_LIMIT orders.
//
// Returns:
//   - *pb.OrderSendData with result details (ticket, status, etc.) or nil on error.
//   - error if not connected, validation fails, or the RPC call fails.
//
// Behavior:
//   - Adds session headers automatically.
//   - For STOP_LIMIT orders, both price (trigger) and stopLimitPrice (limit) are mandatory.
//   - Uses ExecuteWithReconnect for automatic retry on transient gRPC errors.
//
// Notes:
//   - This is the preferred method when working with STOP_LIMIT orders.
//   - Regular orders can also be sent using this method.

func (a *MT5Account) OrderSendEx(
	ctx context.Context,
	symbol string,
	operationType pb.TMT5_ENUM_ORDER_TYPE,
	volume float64,
	price *float64,        // For pending orders; trigger price for STOP_LIMIT
	slippage *int32,       // caller-friendly; will be converted to *uint64
	stoploss *float64,
	takeprofit *float64,
	comment *string,
	magicNumber *int32,    // caller-friendly; will be converted to *uint64
	expiration *timestamppb.Timestamp,
	stopLimitPrice *float64, // Limit price for STOP_LIMIT
) (*pb.OrderSendData, error) {
	// 1) Connection and inputs
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}

	// per-call deadline (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// 2) Build request (match proto types exactly)
	req := &pb.OrderSendRequest{
		Symbol:    symbol,
		Operation: operationType,
		Volume:    volume,
	}
	if price != nil {
		req.Price = price
	}
	if stoploss != nil {
		req.StopLoss = stoploss
	}
	if takeprofit != nil {
		req.TakeProfit = takeprofit
	}
	if comment != nil {
		req.Comment = comment
	}

	// Proto requires uint64 for slippage/expert_id
	req.Slippage = u64pFromI32(slippage)
	req.ExpertId  = u64pFromI32(magicNumber)

	// 3) Expiration → also set expiration_time_type
	if expiration != nil {
		t := pb.TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_SPECIFIED
		req.ExpirationTimeType = &t
		req.ExpirationTime     = expiration
	}

	// 4) STOP_LIMIT specifics: need both trigger (price) and stopLimitPrice
	if operationType == pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP_LIMIT ||
		operationType == pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP_LIMIT {
		if price == nil {
			return nil, fmt.Errorf("STOP_LIMIT requires trigger price (price)")
		}
		if stopLimitPrice == nil {
			return nil, fmt.Errorf("STOP_LIMIT requires stopLimitPrice")
		}
		req.StopLimitPrice = stopLimitPrice
	}

	// 5) Invoke RPC with reconnect wrapper
	grpcCall := func(headers metadata.MD) (*pb.OrderSendReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderSend(c, req)
	}
	errorSelector := func(reply *pb.OrderSendReply) mrpcError { return reply.GetError() }

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	return reply.GetData(), nil
}


// PendingReplaceStopLimit recreates a STOP_LIMIT pending order with new parameters.
// Algorithm:
//
//	(1) Create a new STOP_LIMIT order with the provided params,
//	(2) Attempt to delete the old order by ticket.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - oldTicket: Ticket of the existing STOP_LIMIT order to replace.
//   - symbol: Trading symbol (e.g., "EURUSD").
//   - isBuy: true for BUY_STOP_LIMIT, false for SELL_STOP_LIMIT.
//   - volume: Trade volume in lots (e.g., 0.10, 1.00).
//   - triggerPrice: Trigger price for STOP_LIMIT (the "stop" part).
//   - limitPrice: Limit price that becomes active after the trigger.
//   - slippage: Max slippage (points), optional.
//   - stoploss: Stop Loss, optional.
//   - takeprofit: Take Profit, optional.
//   - comment: Order comment, optional.
//   - magicNumber: Expert/EA identifier, optional.
//   - expiration: Expiration for pending order, optional.
//
// Returns:
//   - *pb.OrderSendData of the newly created order (so you can extract the new ticket per your schema).
//   - error if not connected, validation fails, the create step fails, or deleting the old order fails.
//
// Behavior:
//   - Uses OrderSendStopLimit to create the new order.
//   - Then calls DeleteOrder for the old ticket.
//   - Does NOT rollback the newly created order if deleting the old one fails (schema may not expose ticket directly).
//     If needed, the caller can delete the new order manually by extracting the ticket from the returned data.
//
// Notes:
//   - Atomic replace is not guaranteed. If you need stronger guarantees, implement a two-phase strategy
//     with explicit rollback using the exact new ticket resolved from *pb.OrderSendData.
func (a *MT5Account) PendingReplaceStopLimit(
	ctx context.Context,
	oldTicket uint64,
	symbol string,
	isBuy bool,
	volume float64,
	triggerPrice float64, // STOP_LIMIT trigger price
	limitPrice float64,   // STOP_LIMIT limit price after trigger
	slippage *int32,
	stoploss, takeprofit *float64,
	comment *string,
	magicNumber *int32,
	expiration *timestamppb.Timestamp,
) (*pb.OrderSendData, error) {
	// Ensure connection and inputs
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if oldTicket == 0 {
		return nil, fmt.Errorf("invalid oldTicket")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if volume <= 0 {
		return nil, fmt.Errorf("volume must be > 0")
	}
	if triggerPrice <= 0 || limitPrice <= 0 {
		return nil, fmt.Errorf("triggerPrice and limitPrice must be > 0")
	}

	// per-call deadline if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// (1) Create the new STOP_LIMIT order
	res, err := a.OrderSendStopLimit(
		ctx,
		symbol,
		isBuy,
		volume,
		triggerPrice,
		limitPrice,
		slippage,
		stoploss, takeprofit,
		comment,
		magicNumber,
		expiration,
	)
	if err != nil {
		return nil, fmt.Errorf("create new STOP_LIMIT failed: %w", err)
	}

	// (2) Try to delete the old order
	if _, delErr := a.DeleteOrder(ctx, oldTicket); delErr != nil {
		// No rollback here: the new ticket may not be directly accessible from this schema.
		// The caller can remove the new order manually using data from `res`.
		return res, fmt.Errorf("delete old STOP_LIMIT failed: %w", delErr)
	}

	// Return the newly created order data
	return res, nil
}


// OrderSendStopLimit is a convenience wrapper for placing STOP_LIMIT pending orders.
// It reduces the risk of mixing up triggerPrice and limitPrice parameters.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - symbol: Trading symbol (e.g., "EURUSD").
//   - isBuy: true for BUY_STOP_LIMIT, false for SELL_STOP_LIMIT.
//   - volume: Trade volume in lots (e.g., 0.10, 1.00).
//   - triggerPrice: STOP_LIMIT trigger price (stop).
//   - limitPrice: Limit price that becomes active after the trigger.
//   - slippage: Max allowed slippage in points (optional).
//   - stoploss: Stop Loss price (optional).
//   - takeprofit: Take Profit price (optional).
//   - comment: Optional order comment.
//   - magicNumber: Expert/EA identifier (optional).
//   - expiration: Expiration for pending order (optional).
//
// Returns:
//   - *pb.OrderSendData on success, or nil on error.
//   - error if not connected, validation fails, or RPC call fails.
//
// Behavior:
//   - Maps isBuy to the correct STOP_LIMIT order type.
//   - Always sends both triggerPrice and limitPrice.
//   - Internally calls OrderSendEx with correct parameter mapping.
//
// Notes:
//   - This method is only for STOP_LIMIT orders.
//   - For other order types, use OrderSend or OrderSendEx.

func (a *MT5Account) OrderSendStopLimit(
	ctx context.Context,
	symbol string,
	isBuy bool,
	volume float64,
	triggerPrice float64, // STOP_LIMIT trigger price
	limitPrice float64,   // STOP_LIMIT limit price after trigger
	slippage *int32,
	stoploss, takeprofit *float64,
	comment *string,
	magicNumber *int32,
	expiration *timestamppb.Timestamp,
) (*pb.OrderSendData, error) {
	// Basic input validation (OrderSendEx will also validate, but this gives clearer errors early)
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if volume <= 0 {
		return nil, fmt.Errorf("volume must be > 0")
	}
	if triggerPrice <= 0 || limitPrice <= 0 {
		return nil, fmt.Errorf("triggerPrice and limitPrice must be > 0")
	}

	// Determine order type based on isBuy
	var op pb.TMT5_ENUM_ORDER_TYPE
	if isBuy {
		op = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_BUY_STOP_LIMIT
	} else {
		op = pb.TMT5_ENUM_ORDER_TYPE_TMT5_ORDER_TYPE_SELL_STOP_LIMIT
	}

	// Both trigger and limit prices are required
	trg := triggerPrice
	lim := limitPrice

	// Delegate to OrderSendEx with proper STOP_LIMIT parameters
	return a.OrderSendEx(
		ctx,
		symbol,
		op,
		volume,
		&trg,
		slippage,
		stoploss,
		takeprofit,
		comment,
		magicNumber,
		expiration,
		&lim,
	)
}


// OrderClose closes (or deletes) an existing market or pending order on the connected MT5 terminal.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - ticket: Unique ticket number of the order to close/delete.
//   - volume: Optional partial close volume in lots (nil = close full order).
//   - slippage: Optional max allowed slippage in points (ignored for pending orders).
//
// Returns:
//   - *pb.OrderCloseData with confirmation details, or nil on error.
//   - error if not connected, validation fails, or the gRPC call fails.
//
// Behavior:
//   - Builds an OrderCloseRequest with only the fields that are explicitly set.
//   - Uses ExecuteWithReconnect to retry on recoverable connection/session errors.
//   - Extracts application-level errors from the response before returning.
//
// Notes:
//   - For market orders, slippage applies; for pending orders, it is ignored.
//   - Partial closes are only possible for market orders and certain pending types.
func (a *MT5Account) OrderClose(
	ctx context.Context,
	ticket uint64,
	volume *float64,
	slippage *int32,
) (*pb.OrderCloseData, error) {
	// Ensure connection
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if ticket == 0 {
		return nil, fmt.Errorf("invalid ticket")
	}

	// Build request with mandatory ticket
	req := &pb.OrderCloseRequest{Ticket: ticket}

	// Optional: partial volume close
	if volume != nil {
		req.Volume = *volume
	}
	// Optional: slippage (market orders only)
	if slippage != nil {
		req.Slippage = *slippage // protobuf expects int32 here
	}

	// gRPC call closure with metadata
	// per-call deadline (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	grpcCall := func(headers metadata.MD) (*pb.OrderCloseReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderClose(c, req)
	}

	// Extract application-level error if present
	errorSelector := func(reply *pb.OrderCloseReply) mrpcError { return reply.GetError() }

	// Execute with retry/reconnect semantics
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return the data section from reply
	return reply.GetData(), nil
}


// OrderSelect retrieves detailed information about a currently opened order by its ticket number.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - ticket: Ticket number (unique ID) of the order to look up.
//
// Returns:
//   - *pb.OpenedOrderInfo with full details of the matched order.
//   - error if not connected, the API call fails, or no matching order is found.
//
// Behavior:
//   - Calls OpenedOrders() to retrieve the current list of opened orders.
//   - Iterates over the returned orders to find a ticket match.
//   - Performs local matching; no additional server call is made for a single ticket.
//
// Notes:
//   - Only searches currently opened orders (excludes pending and closed).
//   - Matching is exact on ticket number.
//   - If the ticket does not exist in the current opened orders list, an error is returned.
func (a *MT5Account) OrderSelect(
	ctx context.Context,
	ticket uint64,
) (*pb.OpenedOrderInfo, error) {
	// Ensure connection
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if ticket == 0 {
		return nil, fmt.Errorf("invalid ticket")
	}

	// Per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	data, err := a.OpenedOrders(ctx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no opened orders data")
	}

	for _, o := range data.GetOpenedOrders() {
		if o.GetTicket() == ticket {
			return o, nil
		}
	}
	return nil, fmt.Errorf("order with ticket %d not found", ticket)
}


// OrdersTotal returns the number of currently opened orders on the connected MT5 account.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//
// Returns:
//   - int32: Count of currently opened orders.
//   - error if not connected or the API call fails.
//
// Behavior:
//   - Calls OpenedOrders() to fetch the current list of active market and pending orders.
//   - Returns the length of the list as int32.
//
// Notes:
//   - Closed and historical orders are excluded.
//   - Result reflects the current runtime state of the account.
func (a *MT5Account) OrdersTotal(
	ctx context.Context,
) (int32, error) {
	// Ensure terminal connection
	if !a.isConnected() {
		return 0, errors.New("not connected")
	}

	// Per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// Fetch all currently opened orders
	orders, err := a.OpenedOrders(ctx)
	if err != nil {
		return 0, err
	}

	// Return count of opened orders
	return int32(len(orders.GetOpenedOrders())), nil
}


// OrderModify updates parameters of an existing order on the connected MT5 terminal.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - ticket: Ticket number (unique order ID) to modify.
//   - price: New entry price (optional, pending orders only).
//   - stoploss: New Stop Loss price (optional).
//   - takeprofit: New Take Profit price (optional).
//   - expiration: New expiration timestamp for pending orders (optional).
//
// Returns:
//   - bool: true if modification succeeded, false otherwise.
//   - error if not connected, invalid parameters, or RPC call fails.
//
// Behavior:
//   - Builds an OrderModifyRequest with only the fields provided (non-nil).
//   - Sends request via TradeClient.OrderModify.
//   - Uses ExecuteWithReconnect to retry on transient connection/session errors.
//   - Considers success as presence of Data in the reply and no API-level error.
//
// Notes:
//   - Only pending orders can change price or expiration.
//   - The order must be valid (not closed or fully filled).
func (a *MT5Account) OrderModify(
	ctx context.Context,
	ticket uint64,
	price *float64,
	stoploss *float64,
	takeprofit *float64,
	expiration *timestamppb.Timestamp,
) (*pb.OrderModifyData, error) {

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if ticket == 0 {
		return nil, fmt.Errorf("invalid ticket")
	}

	// per-call deadline (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	req := &pb.OrderModifyRequest{
		Ticket:     ticket,
		Price:      price,
		StopLoss:   stoploss,
		TakeProfit: takeprofit,
	}

	if expiration != nil {
		req.ExpirationTime = expiration
		t := pb.TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_SPECIFIED
		req.ExpirationTimeType = &t
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


// PendingModify updates parameters of a pending order by its ticket number.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//   - ticket: Ticket number (unique order ID) of the pending order.
//   - price: New trigger price (optional). For STOP_LIMIT orders, only trigger can be changed.
//   - stoploss: New Stop Loss level (optional).
//   - takeprofit: New Take Profit level (optional).
//   - expiration: New expiration timestamp (optional).
//
// Returns:
//   - bool: true if modification succeeded, false otherwise.
//   - error if not connected, invalid ticket, or RPC call fails.
//
// Behavior:
//   - Builds an OrderModifyRequest with only the provided non-nil fields.
//   - Sends request via TradeClient.OrderModify.
//   - Uses ExecuteWithReconnect to retry on transient connection/session errors.
//
// Notes:
//   - For STOP_LIMIT orders, the limit price cannot be changed through OrderModify.
//     Use PendingReplaceStopLimit to recreate the order with a new limit price.
//   - The pending order must still be valid and not filled or canceled.
func (a *MT5Account) PendingModify(
	ctx context.Context,
	ticket uint64,
	price, stoploss, takeprofit *float64,
	expiration *timestamppb.Timestamp,
) (bool, error) {
	// Ensure terminal connection
	if !a.isConnected() {
		return false, errors.New("not connected")
	}
	// Validate ticket
	if ticket == 0 {
		return false, fmt.Errorf("invalid ticket")
	}

	// Per-call timeout (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Build request with ticket and optional updates
	req := &pb.OrderModifyRequest{Ticket: ticket}
	if price != nil {
		req.Price = price
	}
	if stoploss != nil {
		req.StopLoss = stoploss
	}
	if takeprofit != nil {
		req.TakeProfit = takeprofit
	}
	if expiration != nil {
		req.ExpirationTime = expiration
		t := pb.TMT5_ENUM_ORDER_TYPE_TIME_TMT5_ORDER_TIME_SPECIFIED
		req.ExpirationTimeType = &t
	}

	// gRPC call closure with metadata
	grpcCall := func(h metadata.MD) (*pb.OrderModifyReply, error) {
		c := metadata.NewOutgoingContext(ctx, h)
		return a.TradeClient.OrderModify(c, req)
	}

	// Extract API-level error if present
	errorSelector := func(reply *pb.OrderModifyReply) mrpcError { return reply.GetError() }

	// Execute with retry/reconnect semantics
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return false, err
	}

	// Success if Data is present
	return reply.GetData() != nil, nil
}


// OpenedOrders retrieves all currently opened orders for the connected MT5 account.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//
// Returns:
//   - *pb.OpenedOrdersData: List of active market and pending orders.
//   - error if not connected or RPC call fails.
//
// Behavior:
//   - Sends an OpenedOrdersRequest to the AccountClient via gRPC.
//   - Uses ExecuteWithReconnect to retry if the connection or terminal session is lost.
//   - Returns API-level error from the reply if present.
//
// Notes:
//   - Only active orders are returned (no history or closed orders).
//   - Result reflects the live account state at the time of the call.
func (a *MT5Account) OpenedOrders(
	ctx context.Context,
) (*pb.OpenedOrdersData, error) {
	// Ensure terminal connection
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// Prepare empty request (no parameters required)
	req := &pb.OpenedOrdersRequest{}

	// gRPC call closure with metadata
	grpcCall := func(headers metadata.MD) (*pb.OpenedOrdersReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.OpenedOrders(c, req)
	}

	// Extract API-level error if present
	errorSelector := func(reply *pb.OpenedOrdersReply) mrpcError { return reply.GetError() }

	// Execute with retry/reconnect semantics
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return data payload
	return reply.GetData(), nil
}


// OpenedOrdersTickets retrieves the ticket IDs and basic details of all currently opened orders.
//
// Parameters:
//   - ctx: Request context (deadline/cancellation).
//
// Returns:
//   - *pb.OpenedOrdersTicketsData: List of tickets with symbols and order types.
//   - error if not connected or RPC call fails.
//
// Behavior:
//   - Sends an OpenedOrdersTicketsRequest to the AccountClient via gRPC.
//   - Uses ExecuteWithReconnect to retry if the connection or terminal session is lost.
//   - Returns API-level error from the reply if present.
//
// Notes:
//   - Only active orders are included (no historical tickets).
//   - This is a lightweight call compared to OpenedOrders, useful when only IDs are needed.
func (a *MT5Account) OpenedOrdersTickets(
	ctx context.Context,
) (*pb.OpenedOrdersTicketsData, error) {
	// Ensure terminal connection
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// Prepare empty request (no parameters required)
	req := &pb.OpenedOrdersTicketsRequest{}

	// gRPC call closure with metadata
	grpcCall := func(headers metadata.MD) (*pb.OpenedOrdersTicketsReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.OpenedOrdersTickets(c, req)
	}

	// Extract API-level error if present
	errorSelector := func(reply *pb.OpenedOrdersTicketsReply) mrpcError { return reply.GetError() }

	// Execute with retry/reconnect semantics
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return data payload
	return reply.GetData(), nil
}


// OrdersHistory retrieves a filtered and paginated list of historical orders from the connected MT5 account.
//
// Parameters:
//   - ctx: Context for cancellation or timeout.
//   - sortType: Sorting mode (e.g., newest first, oldest first) from BMT5_ENUM_ORDER_HISTORY_SORT_TYPE.
//   - from: Optional start date/time (nil = no lower bound).
//   - to: Optional end date/time (nil = no upper bound).
//   - page: Optional page number for pagination (nil = no paging).
//   - itemsPerPage: Optional items per page (nil = server default).
//
// Returns:
//   - *pb.OrdersHistoryData: Historical order records matching filters.
//   - error if not connected or RPC call fails.
//
// Behavior:
//   - Builds OrderHistoryRequest with provided filters and pagination.
//   - Uses ExecuteWithReconnect to retry on connection/session issues.
//   - Returns an empty OrdersHistoryData struct if no data is returned by the server.
func (a *MT5Account) OrdersHistory(
	ctx context.Context,
	sortType pb.BMT5_ENUM_ORDER_HISTORY_SORT_TYPE,
	from, to *time.Time,
	page, itemsPerPage *int32,
) (*pb.OrdersHistoryData, error) {
	// Ensure connection
	if !a.isConnected() {
		return nil, fmt.Errorf("not connected to terminal")
	}

	// Per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Build request with mandatory sort type
	req := &pb.OrderHistoryRequest{
		InputSortMode: sortType,
	}

	// Optional time filters
	if from != nil {
		req.InputFrom = timestamppb.New(from.UTC())
	}
	if to != nil {
		req.InputTo = timestamppb.New(to.UTC())
	}

	// Optional pagination
	if page != nil {
		req.PageNumber = *page
	}
	if itemsPerPage != nil {
		req.ItemsPerPage = *itemsPerPage
	}

	// gRPC call closure
	grpcCall := func(headers metadata.MD) (*pb.OrderHistoryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.OrderHistory(c, req)
	}

	// Error extractor
	errorSelector := func(reply *pb.OrderHistoryReply) mrpcError { return reply.GetError() }

	// Execute with retry/reconnect
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return empty struct if server returned no data
	data := reply.GetData()
	if data == nil {
		return &pb.OrdersHistoryData{}, nil
	}
	return data, nil
}


// DeleteOrder removes a pending order from the MT5 terminal by its ticket number.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - ticket: Unique ticket ID of the order to be deleted.
//
// Returns:
//   - *pb.OrderCloseData containing operation details from the terminal.
//   - error if not connected, RPC fails, or the terminal rejects the request.
//
// Notes:
//   - Volume is not required for pending orders and is ignored.
//   - The method wraps the RPC in ExecuteWithReconnect to handle temporary disconnections.
func (a *MT5Account) DeleteOrder(ctx context.Context, ticket uint64) (*pb.OrderCloseData, error) {
	// Ensure connection is active
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if ticket == 0 {
		return nil, fmt.Errorf("invalid ticket")
	}

	// Build request (Volume not needed for pending orders)
	req := &pb.OrderCloseRequest{
		Ticket: ticket,
	}

	// Per-call timeout if caller didn't set one
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Define gRPC call
	grpcCall := func(headers metadata.MD) (*pb.OrderCloseReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderClose(c, req)
	}

	// Extract any API-level error from reply
	errorSelector := func(reply *pb.OrderCloseReply) mrpcError { return reply.GetError() }

	// Execute RPC with automatic reconnect on recoverable errors
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return the data payload (may be nil if terminal gave no details)
	return reply.GetData(), nil
}


// === 📂 Pending / StopLimit ===

// ShowPendingReplaceStopLimit — create a new STOP_LIMIT and delete the old pending (old ticket).
// ShowPendingReplaceStopLimit replaces an existing pending order with a STOP_LIMIT order and prints the outcome.
func (s *MT5Service) ShowPendingReplaceStopLimit(
	ctx context.Context,
	oldTicket uint64,
	symbol string,
	isBuy bool,
	volume float64,
	triggerPrice float64,
	limitPrice float64,
	slippage *int32,
	stoploss, takeprofit *float64,
	comment *string,
	magicNumber *int32,
	expiration *timestamppb.Timestamp,
) {
	if s.account == nil {
		log.Println("❌ MT5 account not initialized")
		return
	}

	// per-call deadline (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	data, err := s.account.PendingReplaceStopLimit(
		ctx,
		oldTicket,
		symbol,
		isBuy,
		volume,
		triggerPrice,
		limitPrice,
		slippage,
		stoploss, takeprofit,
		comment,
		magicNumber,
		expiration,
	)
	if err != nil {
		log.Printf("❌ PendingReplaceStopLimit error: %v", err)
		return
	}

	newOrder := data.GetOrder()
	newDeal := data.GetDeal()

	if newOrder != 0 {
		fmt.Printf("✅ Replaced pending %d with STOP_LIMIT order %d (trigger=%.5f, limit=%.5f)\n",
			oldTicket, newOrder, triggerPrice, limitPrice)
	} else if newDeal != 0 {
		fmt.Printf("✅ Replaced pending %d with executed deal %d (trigger=%.5f, limit=%.5f)\n",
			oldTicket, newDeal, triggerPrice, limitPrice)
	} else {
		fmt.Printf("✅ Replaced pending %d with STOP_LIMIT (price=%.5f, trigger=%.5f, limit=%.5f)\n",
			oldTicket, data.GetPrice(), triggerPrice, limitPrice)
	}
}


// === 📂 Position ===

// PositionsGet retrieves the list of currently opened positions on the MT5 account.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//
// Returns:
//   - []*pb.PositionInfo: Slice of position information structures.
//   - error if the terminal is not connected or the gRPC request fails.
//
// Behavior:
//   - Calls OpenedOrders() to retrieve the current state of orders and positions.
//   - Extracts and returns only the PositionInfo entries.
//
// Notes:
//   - This method relies on OpenedOrders(), which already wraps the RPC in ExecuteWithReconnect.
//   - Returned slice may be empty if no positions are currently open.
func (a *MT5Account) PositionsGet(ctx context.Context) ([]*pb.PositionInfo, error) {
	// Ensure terminal connection
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Fetch all currently opened orders/positions
	opened, err := a.OpenedOrders(ctx)
	if err != nil {
		return nil, err
	}

	// Extract only positions from the result
	if opened == nil || opened.GetPositionInfos() == nil {
		return []*pb.PositionInfo{}, nil
	}
	return opened.GetPositionInfos(), nil
}


// PositionGet retrieves information about a currently opened position for a specific symbol.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - symbol: Trading symbol to search for (e.g., "EURUSD").
//
// Returns:
//   - *pb.PositionInfo: Position details if found, or nil if no position exists for the given symbol.
//   - error if the terminal is not connected or the RPC request fails.
//
// Behavior:
//   - Calls OpenedOrders() to get the list of current orders/positions.
//   - Iterates over PositionInfos to find a match by symbol.
//
// Notes:
//   - This method does not query closed or historical positions.
//   - Returns nil without error if no matching position is found.
//   - OpenedOrders() internally uses ExecuteWithReconnect for connection retries.
func (a *MT5Account) PositionGet(ctx context.Context, symbol string) (*pb.PositionInfo, error) {
	// Ensure terminal connection
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}

	// Per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// Fetch current opened orders/positions
	opened, err := a.OpenedOrders(ctx)
	if err != nil {
		return nil, err
	}
	if opened == nil || opened.GetPositionInfos() == nil {
		return nil, nil
	}

	// Search for position by symbol
	for _, p := range opened.GetPositionInfos() {
		if p.GetSymbol() == symbol {
			return p, nil
		}
	}

	// No position found for the given symbol
	return nil, nil
}


// HasOpenPosition checks if there is an active open position for a given symbol.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - symbol: Trading symbol to check (e.g., "EURUSD").
//
// Returns:
//   - bool: true if a position exists and has a valid ticket, false otherwise.
//   - error: if the terminal is not connected or the RPC call fails.
//
// Behavior:
//   - Calls PositionGet() to retrieve the position info.
//   - Returns true if the position object is non-nil and has a non-zero ticket.
func (a *MT5Account) HasOpenPosition(ctx context.Context, symbol string) (bool, error) {
	// Try to fetch position by symbol
	p, err := a.PositionGet(ctx, symbol)
	if err != nil {
		return false, err
	}

	// Valid position means object is not nil and ticket is non-zero
	return p != nil && p.GetTicket() != 0, nil
}


// PositionClose closes an active open position by its PositionInfo object.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//   - position: Pointer to PositionInfo (must be non-nil with a valid ticket).
//
// Returns:
//   - *pb.OrderCloseData: Details about the close operation.
//   - error: If not connected, position is invalid, or the RPC call fails.
//
// Behavior:
//   - Validates connection and position parameters.
//   - Prepares OrderCloseRequest with ticket and full position volume.
//   - Calls TradeClient.OrderClose() via ExecuteWithReconnect for retries.
//   - Returns the data part of the response if successful.
func (a *MT5Account) PositionClose(ctx context.Context, position *pb.PositionInfo) (*pb.OrderCloseData, error) {
	// Ensure connected to terminal
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	// Validate input
	if position == nil || position.GetTicket() == 0 {
		return nil, fmt.Errorf("invalid position")
	}

	// per-call deadline (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Prepare the gRPC request to close the position
	req := &pb.OrderCloseRequest{
		Ticket: position.GetTicket(),
		Volume: position.GetVolume(), // Close full volume by default
		// Slippage can be set if needed:
		// Slippage: int32(value),
	}

	// Define the gRPC call with session headers
	grpcCall := func(headers metadata.MD) (*pb.OrderCloseReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderClose(c, req)
	}
	// Define how to extract an API error from the response
	errorSelector := func(reply *pb.OrderCloseReply) mrpcError { return reply.GetError() }

	// Execute the request with retry/reconnect support
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return the operation result
	return reply.GetData(), nil
}


// CloseAllPositions closes all currently open positions sequentially.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//
// Returns:
//   - error: The first error encountered during the closing process (if any),
//     or nil if all positions are closed successfully.
//
// Behavior:
//   - Retrieves the list of all open positions via PositionsGet().
//   - Iterates through each position and calls PositionClose().
//   - If a close operation fails, continues with the remaining positions,
//     but remembers the first error to return at the end.
func (a *MT5Account) CloseAllPositions(ctx context.Context) error {
	// Ensure we are connected to the terminal
	if !a.isConnected() {
		return errors.New("not connected")
	}

	// Retrieve all open positions
	positions, err := a.PositionsGet(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch positions: %w", err)
	}
	if len(positions) == 0 {
		return nil
	}

	var firstErr error

	// Loop through each position and attempt to close it
	for _, pos := range positions {
		if pos == nil || pos.GetTicket() == 0 {
			continue
		}
		if _, cerr := a.PositionClose(ctx, pos); cerr != nil && firstErr == nil {
			// Store the first error but keep closing the rest
			firstErr = fmt.Errorf(
				"failed to close position %d (%s): %w",
				pos.GetTicket(),
				pos.GetSymbol(),
				cerr,
			)
		}
	}

	return firstErr
}


// PositionModify updates the Stop Loss and/or Take Profit levels
// of an existing open position.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - ticket: Unique ID (ticket number) of the position to modify.
//   - stoploss: New Stop Loss price level (optional; nil = unchanged).
//   - takeprofit: New Take Profit price level (optional; nil = unchanged).
//
// Returns:
//   - bool: true if the modification was successful, false otherwise.
//   - error: if the terminal is not connected, ticket is invalid,
//     or the RPC call fails.
//
// Behavior:
//   - Builds an OrderModifyRequest with only the provided fields.
//   - Sends the request to TradeClient.OrderModify via gRPC.
//   - Uses ExecuteWithReconnect to retry on recoverable connection errors.
//   - Success is determined by the presence of a non-nil Data field in the reply.
func (a *MT5Account) PositionModify(
	ctx context.Context,
	ticket uint64,
	stoploss, takeprofit *float64,
) (bool, error) {
	// Ensure terminal connection
	if !a.isConnected() {
		return false, errors.New("not connected")
	}
	if ticket == 0 {
		return false, fmt.Errorf("invalid ticket")
	}

	// Per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Prepare modification request
	req := &pb.OrderModifyRequest{
		Ticket: ticket,
	}
	if stoploss != nil {
		req.StopLoss = stoploss
	}
	if takeprofit != nil {
		req.TakeProfit = takeprofit
	}

	// gRPC call definition
	grpcCall := func(headers metadata.MD) (*pb.OrderModifyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderModify(c, req)
	}

	// Error extraction logic
	errorSelector := func(reply *pb.OrderModifyReply) mrpcError { return reply.GetError() }

	// Execute with reconnect logic
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return false, err
	}

	// Success if Data field exists
	return reply.GetData() != nil, nil
}


// === 📂 History ===

// HistoryDealsGet retrieves the list of executed deals within a given time range,
// optionally filtered by trading symbol.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - from: Start time for the query (inclusive).
//   - to: End time for the query (inclusive).
//   - symbol: Optional trading symbol filter (e.g., "EURUSD").
//     Pass an empty string to include all symbols.
//
// Returns:
//   - []*pb.DealHistoryData: Slice of deal records matching the criteria.
//   - error: if the terminal is not connected, or the RPC call fails.
//
// Behavior:
//   - Sends a request to AccountClient.OrderHistory with the given date range.
//   - The request does not specify sort mode or pagination unless added manually.
//   - Filters the returned history to only include deal records (HistoryDeal),
//     excluding non-deal history entries.
//   - If 'symbol' is provided, only deals matching that symbol are returned.
//
// Note:
//   - The MT5 server returns a mix of orders and deals in OrderHistory.
//     This method extracts only deal entries from that history.
func (a *MT5Account) HistoryDealsGet(
	ctx context.Context,
	from, to time.Time,
	symbol string,
) ([]*pb.DealHistoryData, error) {
	// Ensure connection to terminal
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	// Validate time range
	if to.Before(from) {
		return nil, fmt.Errorf("invalid time range: to < from")
	}

	// per-call deadline (5s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Build request for the history period
	req := &pb.OrderHistoryRequest{
		InputFrom: timestamppb.New(from.UTC()),
		InputTo:   timestamppb.New(to.UTC()),
	}

	// Define gRPC request logic
	grpcCall := func(headers metadata.MD) (*pb.OrderHistoryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.OrderHistory(c, req)
	}

	// Define API error extraction
	errorSelector := func(reply *pb.OrderHistoryReply) mrpcError { return reply.GetError() }

	// Execute call with reconnect logic
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	if reply.GetData() == nil {
		return []*pb.DealHistoryData{}, nil
	}

	// Extract only deal entries from history
	hist := reply.GetData().GetHistoryData()
	res := make([]*pb.DealHistoryData, 0, len(hist))
	for _, h := range hist {
		d := h.GetHistoryDeal()
		if d == nil {
			continue
		}
		if symbol == "" || d.GetSymbol() == symbol {
			res = append(res, d)
		}
	}
	return res, nil
}


// HistoryDealsTotal counts the total number of executed deals within a given time range,
// optionally filtered by trading symbol.
//
// This method performs a paginated fetch of historical orders via AccountClient.OrderHistory,
// and sums up only those entries that represent actual deals (HistoryDeal).
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - from: Start time for the search (inclusive).
//   - to: End time for the search (inclusive).
//   - symbol: Optional trading symbol filter (e.g., "EURUSD").
//     Pass an empty string to include all symbols.
//
// Returns:
//   - int32: Total count of matching deals.
//   - error: If the terminal is not connected or the RPC request fails.
//
// Notes:
//   - Pagination is used with a fixed page size (500 items by default).
//   - Stops fetching when the last page contains fewer items than the page size.
func (a *MT5Account) HistoryDealsTotal(
	ctx context.Context,
	from, to time.Time,
	symbol string,
) (int32, error) {
	// Ensure the account is connected before making requests
	if !a.isConnected() {
		return 0, errors.New("not connected")
	}
	// Validate time range
	if to.Before(from) {
		return 0, fmt.Errorf("invalid time range: to < from")
	}

	var total int32
	page := int32(1)
	const pageSize = int32(500) // adjust if server allows larger page sizes

	for {
		// Build paginated request
		req := &pb.OrderHistoryRequest{
			InputFrom:    timestamppb.New(from.UTC()),
			InputTo:      timestamppb.New(to.UTC()),
			PageNumber:   page,
			ItemsPerPage: pageSize,
		}

		// Per-iteration timeout if caller didn't set one
		iterCtx := ctx
		if iterCtx == nil {
			iterCtx = context.Background()
		}
		var cancel context.CancelFunc
		if _, ok := iterCtx.Deadline(); !ok {
			iterCtx, cancel = context.WithTimeout(iterCtx, 5*time.Second)
		}

		// Define gRPC call
		grpcCall := func(headers metadata.MD) (*pb.OrderHistoryReply, error) {
			c := metadata.NewOutgoingContext(iterCtx, headers)
			return a.AccountClient.OrderHistory(c, req)
		}

		// Extract API-level error from the response
		errorSelector := func(reply *pb.OrderHistoryReply) mrpcError { return reply.GetError() }

		// Execute request with retry/reconnect support
		reply, err := ExecuteWithReconnect(a, iterCtx, grpcCall, errorSelector)
		if cancel != nil {
			cancel()
		}
		if err != nil {
			return 0, err
		}

		data := reply.GetData()
		if data == nil {
			break
		}
		items := data.GetHistoryData()
		if len(items) == 0 {
			break
		}

		// Count only matching deals
		for _, h := range items {
			if d := h.GetHistoryDeal(); d != nil {
				if symbol == "" || d.GetSymbol() == symbol {
					total++
				}
			}
		}

		// If last page has fewer items than page size — stop fetching
		if int32(len(items)) < pageSize {
			break
		}
		page++
	}

	return total, nil
}


// HistoryOrderByTicket searches for a specific historical order by its ticket number.
//
// This method iterates over all pages of historical orders via AccountClient.OrderHistory,
// looking for an order entry whose ticket matches the given one.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - ticket: Unique identifier (ticket number) of the historical order.
//
// Returns:
//   - *pb.OrderHistoryData: The matching historical order data.
//   - error: If the order is not found or the RPC request fails.
//
// Notes:
//   - The search is performed without time filters — it scans the full history.
//   - Pagination is used with a fixed page size (500 items by default).
//   - Stops when the order is found or when the last page is reached.
func (a *MT5Account) HistoryOrderByTicket(
	ctx context.Context,
	ticket uint64,
) (*pb.OrderHistoryData, error) {
	// Ensure the account is connected
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if ticket == 0 {
		return nil, fmt.Errorf("invalid ticket")
	}

	page := int32(1)
	const pageSize = int32(500) // adjust if needed

	for {
		// Build request for the current page
		req := &pb.OrderHistoryRequest{
			PageNumber:   page,
			ItemsPerPage: pageSize,
			// No time filters — scanning the full history
		}

		// Per-iteration timeout (keeps long scans responsive)
		iterCtx := ctx
		if iterCtx == nil {
			iterCtx = context.Background()
		}
		var cancel context.CancelFunc
		if _, ok := iterCtx.Deadline(); !ok {
			iterCtx, cancel = context.WithTimeout(iterCtx, 5*time.Second)
		}

		// Define the gRPC call
		grpcCall := func(headers metadata.MD) (*pb.OrderHistoryReply, error) {
			c := metadata.NewOutgoingContext(iterCtx, headers)
			return a.AccountClient.OrderHistory(c, req)
		}

		// Extract API-level error
		errorSelector := func(reply *pb.OrderHistoryReply) mrpcError { return reply.GetError() }

		// Execute with automatic retry/reconnect
		reply, err := ExecuteWithReconnect(a, iterCtx, grpcCall, errorSelector)
		if cancel != nil {
			cancel()
		}
		if err != nil {
			return nil, err
		}

		data := reply.GetData()
		if data == nil {
			break
		}
		items := data.GetHistoryData()
		if len(items) == 0 {
			break
		}

		// Search for the matching ticket in the current page
		for _, h := range items {
			o := h.GetHistoryOrder()
			if o == nil {
				continue
			}
			if uint64(o.GetTicket()) == ticket {
				return o, nil
			}
		}

		// If the last page contains fewer items than pageSize, stop searching
		if int32(len(items)) < pageSize {
			break
		}
		page++
	}

	return nil, fmt.Errorf("order %d not found", ticket)
}


// HistoryDealByTicket searches for a specific historical deal by its ticket number.
//
// This method iterates through pages of order history from AccountClient.OrderHistory,
// looking for a deal entry (HistoryDeal) with the given ticket.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - ticket: Unique identifier (ticket number) of the deal.
//
// Returns:
//   - *pb.DealHistoryData: The matching deal data.
//   - error: If the deal is not found or the RPC call fails.
//
// Notes:
//   - The search is performed without time filters — it scans the full history.
//   - Pagination is used with a fixed page size (500 items by default).
//   - Stops when the deal is found or when the last page is reached.
func (a *MT5Account) HistoryDealByTicket(
	ctx context.Context,
	ticket uint64,
) (*pb.DealHistoryData, error) {
	// Ensure the account is connected
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if ticket == 0 {
		return nil, fmt.Errorf("invalid ticket")
	}

	page := int32(1)
	const pageSize = int32(500) // adjust if needed

	for {
		// Build request for the current page
		req := &pb.OrderHistoryRequest{
			PageNumber:   page,
			ItemsPerPage: pageSize,
			// Optional: can add InputFrom / InputTo filters if needed
		}

		// Per-iteration timeout (keeps long scans responsive)
		iterCtx := ctx
		if iterCtx == nil {
			iterCtx = context.Background()
		}
		var cancel context.CancelFunc
		if _, ok := iterCtx.Deadline(); !ok {
			iterCtx, cancel = context.WithTimeout(iterCtx, 5*time.Second)
		}

		// Define the gRPC call
		grpcCall := func(headers metadata.MD) (*pb.OrderHistoryReply, error) {
			c := metadata.NewOutgoingContext(iterCtx, headers)
			return a.AccountClient.OrderHistory(c, req)
		}

		// Extract API-level error
		errorSelector := func(reply *pb.OrderHistoryReply) mrpcError { return reply.GetError() }

		// Execute with automatic retry/reconnect
		reply, err := ExecuteWithReconnect(a, iterCtx, grpcCall, errorSelector)
		if cancel != nil {
			cancel()
		}
		if err != nil {
			return nil, err
		}

		data := reply.GetData()
		if data == nil {
			break
		}
		items := data.GetHistoryData()
		if len(items) == 0 {
			break
		}

		// Search for the matching ticket in the current page
		for _, h := range items {
			d := h.GetHistoryDeal()
			if d == nil {
				continue
			}
			if uint64(d.GetTicket()) == ticket {
				return d, nil
			}
		}

		// If fewer items than pageSize — no more pages
		if int32(len(items)) < pageSize {
			break
		}
		page++
	}

	return nil, fmt.Errorf("deal %d not found", ticket)
}


// === 📂 Trade Functions ===

// OrderCalcMargin computes the required margin for a hypothetical order,
// using the server-side calculation (same rules as the terminal).
//
// Parameters:
//   - ctx: Context for timeout/cancellation.
//   - req: *pb.OrderCalcMarginRequest containing symbol, order type, volume, price, etc.
//
// Returns:
//   - *pb.OrderCalcMarginData: Calculated margin and related fields.
//   - error: If not connected or the RPC/API call fails.
//
// Behavior:
//   - Sends TradeFunctionsClient.OrderCalcMargin with the provided request.
//   - Uses ExecuteWithReconnect to retry on transient gRPC errors.
//   - Returns the Data part of the reply on success.
func (a *MT5Account) OrderCalcMargin(
	ctx context.Context,
	req *pb.OrderCalcMarginRequest,
) (*pb.OrderCalcMarginData, error) {
	// 1) Connection guard
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	// per-call deadline (3s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// 2) Define RPC call
	grpcCall := func(headers metadata.MD) (*pb.OrderCalcMarginReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.OrderCalcMargin(c, req)
	}

	// 3) API error extractor
	errorSelector := func(reply *pb.OrderCalcMarginReply) mrpcError { return reply.GetError() }

	// 4) Execute with reconnect logic
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// 5) Return payload
	return reply.GetData(), nil
}


// OrderCalcProfit calculates the potential profit or loss for a hypothetical order
// using the server-side calculation (same rules as the terminal).
//
// Parameters:
//   - ctx: Context for timeout/cancellation.
//   - req: *pb.OrderCalcProfitRequest containing symbol, order type, volume, price, etc.
//
// Returns:
//   - *pb.OrderCalcProfitData: Calculated profit/loss and related fields.
//   - error: If not connected or the RPC/API call fails.
//
// Behavior:
//   - Sends TradeFunctionsClient.OrderCalcProfit with the provided request.
//   - Uses ExecuteWithReconnect to retry on transient gRPC errors.
//   - Returns the Data part of the reply on success.
func (a *MT5Account) OrderCalcProfit(
	ctx context.Context,
	req *pb.OrderCalcProfitRequest,
) (*pb.OrderCalcProfitData, error) {
	// 1) Connection guard
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	// per-call deadline (3s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// 2) Define RPC call
	grpcCall := func(headers metadata.MD) (*pb.OrderCalcProfitReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.OrderCalcProfit(c, req)
	}

	// 3) API error extractor
	errorSelector := func(reply *pb.OrderCalcProfitReply) mrpcError { return reply.GetError() }

	// 4) Execute with reconnect logic
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// 5) Return payload
	return reply.GetData(), nil
}


// OrderCheck checks the validity of an order request (e.g., if it can be placed),
// including margin and other parameters.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//   - req: *pb.OrderCheckRequest containing order parameters such as symbol, volume, price, etc.
//
// Returns:
//   - *pb.OrderCheckData: Data including validation result and error messages (if any).
//   - error: If not connected, or if the RPC/API call fails.
//
// Behavior:
//   - Sends TradeFunctionsClient.OrderCheck with the provided request.
//   - Uses ExecuteWithReconnect to handle transient gRPC errors.
//   - Returns the Data part of the reply if the check passes.
func (a *MT5Account) OrderCheck(
	ctx context.Context,
	req *pb.OrderCheckRequest,
) (*pb.OrderCheckData, error) {
	// 1) Ensure the account is connected
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	// Guard against nil request
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	// Per-call timeout (3s) if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// 2) Define the RPC call
	grpcCall := func(headers metadata.MD) (*pb.OrderCheckReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.OrderCheck(c, req)
	}

	// 3) Define how to extract API errors from the response
	errorSelector := func(reply *pb.OrderCheckReply) mrpcError { return reply.GetError() }

	// 4) Execute the RPC with reconnect logic in case of errors
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// 5) Return the data from the reply
	return reply.GetData(), nil
}


// PositionsTotal returns the total number of positions currently open on the connected MT4 account.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//
// Returns:
//   - int32: Total number of open positions on the account.
//   - error: If not connected or if the API call fails.
//
// Behavior:
//   - Sends the request to the TradeFunctionsClient.PositionsTotal RPC call.
//   - Uses ExecuteWithReconnect to automatically retry on transient connection errors.
//   - Returns the number of open positions from the reply data.
func (a *MT5Account) PositionsTotal(ctx context.Context) (int32, error) {
	// 1) Ensure the account is connected
	if !a.isConnected() {
		return 0, errors.New("not connected")
	}

	// per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// 2) Create an empty request (PositionsTotal does not need any parameters)
	req := &emptypb.Empty{}

	// 3) Define the RPC call with metadata
	grpcCall := func(headers metadata.MD) (*pb.PositionsTotalReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeFunctionsClient.PositionsTotal(c, req)
	}

	// 4) Define how to extract errors from the reply
	errorSelector := func(reply *pb.PositionsTotalReply) mrpcError { return reply.GetError() }

	// 5) Execute the RPC with automatic retry logic in case of errors
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return 0, err
	}

	// 6) Return the total positions from the response
	if reply == nil || reply.GetData() == nil {
		return 0, nil
	}
	return reply.GetData().GetTotalPositions(), nil
}


// === 📂 Market Info / Symbol Info ===

// Quote retrieves the latest market quote (bid/ask) for a given symbol.
//
// Parameters:
//   - ctx: Context for cancellation or timeout control.
//   - symbol: Symbol name (e.g., "EURUSD") to fetch the quote for.
//
// Returns:
//   - *pb.OnSymbolTickData: Contains the latest bid/ask, high/low, timestamp, etc.
//   - error: If connection to the terminal fails or API call encounters an error.
//
// Behavior:
//   - This method uses the SubscriptionService.OnSymbolTick method to receive a real-time tick.
//   - It waits for the first received tick and returns it.
//   - Uses a context with a timeout of 3 seconds to ensure the call doesn't block indefinitely.
func (a *MT5Account) Quote(ctx context.Context, symbol string) (*pb.OnSymbolTickData, error) {
	// 1) Check if the account is connected to the terminal.
	if !a.IsConnected() {
		return nil, errors.New("not connected to terminal")
	}

	// 2) Validate input.
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}

	// 3) Ensure ctx is non-nil.
	if ctx == nil {
		ctx = context.Background()
	}
	// 3a) ✅ Added timeout: 3s to receive the first tick.
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 4) Subscribe to symbol ticks.
	dataCh, errCh := a.OnSymbolTick(ctx2, []string{symbol})

	// 5) Wait for first tick or error.
	for {
		select {
		case tick, ok := <-dataCh:
			if !ok {
				return nil, fmt.Errorf("tick stream closed")
			}
			if tick != nil {
				if st := tick.GetSymbolTick(); st != nil && st.GetSymbol() == symbol {
					return tick, nil
				}
			}

		case err, ok := <-errCh:
			// ✅ Changed: handle closed channel without error (don't hang).
			if !ok || err == nil {
				continue // no error, just wait for data or ctx timeout
			}
			return nil, err

		case <-ctx2.Done():
			return nil, ctx2.Err()
		}
	}
}



// QuoteMany retrieves the latest market quotes for multiple trading symbols.
//
// This method collects the first "len(symbols)" ticks from the stream and returns them as a list of QuoteData entries.
//
// Parameters:
//   - ctx: Context for timeout or cancellation control.
//   - symbols: Slice of symbol names (e.g., []string{"EURUSD", "GBPUSD"}).
//
// Returns:
//   - A slice of *pb.OnSymbolTickData containing quotes for each symbol.
//   - An error if the request fails or if the account is not connected.
func (a *MT5Account) QuoteMany(ctx context.Context, symbols []string) ([]*pb.OnSymbolTickData, error) {
	// 1) Check if the account is connected to the terminal.
	if !a.IsConnected() { // CHANGED
		return nil, errors.New("not connected to terminal")
	}

	// 2) Ensure the list of symbols is not empty.
	if len(symbols) == 0 {
		return nil, errors.New("symbols list is empty")
	}

	// 2a) Make ctx safe and add a 5s per-call timeout if none provided.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// 3) Subscribe to the tick stream for the given symbols.
	dataCh, errCh := a.OnSymbolTick(ctx, symbols)

	// 4) Collect the first tick per requested symbol (no duplicates).
	want := make(map[string]struct{}, len(symbols))
	for _, s := range symbols {
		want[s] = struct{}{}
	}
	seen := make(map[string]*pb.OnSymbolTickData, len(symbols))

	// 5) Wait until we have one tick for each symbol or we time out/err.
	for len(seen) < len(want) {
		select {
		case tick, ok := <-dataCh:
			if !ok {
				// Stream closed; return whatever we have so far (or error if none).
				if len(seen) == 0 {
					return nil, fmt.Errorf("tick stream closed before any data")
				}
				// Preserve caller's order when returning partial results.
				out := make([]*pb.OnSymbolTickData, 0, len(seen))
				for _, s := range symbols {
					if t, ok := seen[s]; ok {
						out = append(out, t)
					}
				}
				return out, nil
			}
			if tick != nil {
				if st := tick.GetSymbolTick(); st != nil {
					sym := st.GetSymbol()
					if _, wanted := want[sym]; wanted {
						if _, already := seen[sym]; !already {
							seen[sym] = tick
						}
					}
				}
			}

		case err, ok := <-errCh: // CHANGED: correctly handling the errCh closure
			if !ok {
				errCh = nil // disabling the select branch to avoid sticking
				continue
			}
			if err != nil {
				return nil, err
			}
			// if err == nil, we continue to wait for other cases

		case <-ctx.Done():
			// Timeout/cancel: return what we have, or the timeout error if empty.
			if len(seen) == 0 {
				return nil, ctx.Err()
			}
			out := make([]*pb.OnSymbolTickData, 0, len(seen))
			for _, s := range symbols {
				if t, ok := seen[s]; ok {
					out = append(out, t)
				}
			}
			return out, nil
		}
	}

	// 6) Build result in the same order as requested.
	out := make([]*pb.OnSymbolTickData, 0, len(symbols))
	for _, s := range symbols {
		out = append(out, seen[s])
	}
	return out, nil
}



// ShowAllSymbols retrieves all available trading symbols from the server.
//
// This method sends a request to the MetaTrader terminal (via gRPC) asking for
// a list of all known trading instruments (symbols), such as "EURUSD", "GBPJPY", etc.
// These symbols include both visible (in Market Watch) and hidden ones.
//
// ⚠️ This method returns the list of all available symbols, not the action of making them visible.
// The list includes all symbols (both visible and hidden) from the MetaTrader terminal.
func (a *MT5Account) ShowAllSymbols(ctx context.Context) ([]string, error) {
	// 1) Check if the account is connected to the MetaTrader terminal.
	if !a.isConnected() {
		return nil, errors.New("not connected to terminal")
	}

	// 2) Request the total number of symbols available in the terminal.
	total, err := a.SymbolsTotal(ctx)
	if err != nil {
		return nil, fmt.Errorf("SymbolsTotal failed: %w", err)
	}
	// If there are no symbols, return an empty list.
	if total <= 0 {
		return []string{}, nil
	}

	// 3) Initialize a slice to hold the names of the symbols.
	names := make([]string, 0, int(total))

	// 4) Iterate through the symbols by their index.
	for i := int32(0); i < total; i++ {
		// Retrieve the symbol name by index.
		// The third argument "false" indicates we don't need to check visibility.
		name, err := a.SymbolName(ctx, i, false)
		if err != nil {
			return nil, fmt.Errorf("SymbolName(%d) failed: %w", i, err)
		}
		// If the name is not empty, append it to the names list.
		if name != "" {
			names = append(names, name)
		}
	}

	// 5) Return the list of symbol names.
	return names, nil
}


// SymbolParams retrieves detailed trading parameters for a single symbol.
//
// Parameters:
//   - ctx: Context for cancellation or timeout.
//   - symbol: Name of the trading symbol (e.g., "EURUSD").
//
// Returns:
//   - *pb.SymbolParameters: Detailed symbol info (digits, volume limits, trade mode, etc.)
//   - error: If request fails or symbol not found.
//
// Notes:
//   - Internally calls SymbolParamsMany with one symbol.
//   - Returns the first result in the symbol info list.
//   - Performs automatic reconnect if the terminal connection is lost.


func (a *MT5Account) SymbolParams(ctx context.Context, symbol string) (*pb.SymbolParameters, error) {
	if !a.IsConnected() {
		return nil, fmt.Errorf("not connected to terminal")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}

	req := &pb.SymbolParamsManyRequest{
		SymbolName: proto.String(symbol),
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolParamsManyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.SymbolParamsMany(c, req)
	}
	errorSelector := func(reply *pb.SymbolParamsManyReply) mrpcError { return reply.GetError() }

	resp, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// ✅ CHANGED: secure nil checks
	data := resp.GetData()
	if data == nil {
		return nil, fmt.Errorf("no parameters returned for symbol: %s", symbol)
	}
	infos := data.GetSymbolInfos()
	if len(infos) == 0 || infos[0] == nil {
		return nil, fmt.Errorf("no parameters returned for symbol: %s", symbol)
	}

	return infos[0], nil
}



// Symbols retrieves a list of symbol names from the MetaTrader terminal.
// It can either return all symbols or only those in the Market Watch (selected space),
// depending on the selectedOnly parameter.
//
// Parameters:
//   - ctx: Context for cancellation or timeout control.
//   - selectedOnly: If true, only symbols from the Market Watch (Selected space) will be returned.
//     If false, all symbols (All space) will be included.
//
// Returns:
//   - A slice of strings representing the symbol names (e.g., "EURUSD", "GBPJPY").
//   - An error if the connection is not established or the request fails.
//
// This method first fetches the total number of symbols from the terminal using the SymbolsTotal request,
// and then iterates through each symbol index to fetch its name using the SymbolName request. The method
// will perform automatic retries in case of transient connection issues.
//
// Note:
//   - The `selectedOnly` flag determines which symbol set is retrieved — Market Watch or all available symbols.
//   - The method uses ExecuteWithReconnect to automatically handle reconnection and retries if the connection is lost.
func (a *MT5Account) Symbols(ctx context.Context, selectedOnly bool) ([]string, error) {
	// Ensure that the account is connected to the MetaTrader terminal.
	if !a.isConnected() {
		return nil, fmt.Errorf("not connected to terminal")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// 1) Get the total number of symbols.
	totalReq := &pb.SymbolsTotalRequest{Mode: selectedOnly}

	grpcCallTotal := func(headers metadata.MD) (*pb.SymbolsTotalReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolsTotal(c, totalReq)
	}
	errorSelectorTotal := func(reply *pb.SymbolsTotalReply) mrpcError { return reply.GetError() }

	totalReply, err := ExecuteWithReconnect(a, ctx, grpcCallTotal, errorSelectorTotal)
	if err != nil {
		return nil, err
	}
	if totalReply.GetData() == nil {
		return []string{}, nil
	}
	total := totalReply.GetData().GetTotal()
	if total <= 0 {
		return []string{}, nil
	}

	// 2) Iterate over symbol indices and fetch the symbol names.
	names := make([]string, 0, int(total))
	for i := int32(0); i < total; i++ {
		nameReq := &pb.SymbolNameRequest{
			Index:    i,
			Selected: selectedOnly,
		}

		grpcCallName := func(headers metadata.MD) (*pb.SymbolNameReply, error) {
			c := metadata.NewOutgoingContext(ctx, headers)
			return a.MarketInfoClient.SymbolName(c, nameReq)
		}
		errorSelectorName := func(reply *pb.SymbolNameReply) mrpcError { return reply.GetError() }

		nameReply, err := ExecuteWithReconnect(a, ctx, grpcCallName, errorSelectorName)
		if err != nil {
			return nil, fmt.Errorf("SymbolName(%d) failed: %w", i, err)
		}
		if data := nameReply.GetData(); data != nil && data.GetName() != "" {
			names = append(names, data.GetName())
		}
	}

	return names, nil
}


// SymbolParamsMany retrieves trading symbol parameters for one or more symbols.
//
// Parameters:
//   - ctx: Context for timeout/cancel.
//   - symbols: List of symbol names.
//
// Returns:
//   - SymbolParamsManyData containing all symbol param info.
func (a *MT5Account) SymbolParamsMany(ctx context.Context, symbols []string) (*pb.SymbolParamsManyData, error) {
	// Check if the account is connected to the terminal.
	if !a.isConnected() {
		return nil, errors.New("not connected to terminal")
	}

	// Safe ctx + short per-call timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// If the list is empty, return all symbol parameters.
	if len(symbols) == 0 {
		req := &pb.SymbolParamsManyRequest{}

		grpcCall := func(headers metadata.MD) (*pb.SymbolParamsManyReply, error) {
			c := metadata.NewOutgoingContext(ctx, headers)
			return a.AccountClient.SymbolParamsMany(c, req)
		}
		errorSelector := func(reply *pb.SymbolParamsManyReply) mrpcError { return reply.GetError() }

		reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
		if err != nil {
			return nil, err
		}
		if reply == nil || reply.GetData() == nil {
			return &pb.SymbolParamsManyData{}, nil
		}
		return reply.GetData(), nil
	}

	// Otherwise, collect parameters for each specified symbol.
	out := &pb.SymbolParamsManyData{
		SymbolInfos: make([]*pb.SymbolParameters, 0, len(symbols)),
	}
	for _, s := range symbols {
		p, err := a.SymbolParams(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("SymbolParams(%s) failed: %w", s, err)
		}
		out.SymbolInfos = append(out.SymbolInfos, p)
	}

	// Return the collected symbol parameters.
	return out, nil
}


// TickValueWithSize calculates tick value, lot size, and contract size for specified symbols.
//
// Parameters:
//   - ctx: Context for cancel/timeout.
//   - symbols: List of symbol+lot pairs.
//
// Returns:
//   - TickValueWithSizeData containing value calculations.
func (a *MT5Account) TickValueWithSize(ctx context.Context, symbolNames []string) (*pb.TickValueWithSizeData, error) {
	// Ensure the account is connected to the terminal before making the request.
	if !a.isConnected() {
		return nil, errors.New("not connected to terminal")
	}
	// Validate input
	if len(symbolNames) == 0 {
		return nil, fmt.Errorf("symbols list is empty")
	}

	// Ensure non-nil ctx and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// Create the request object containing the list of symbols.
	req := &pb.TickValueWithSizeRequest{
		SymbolNames: symbolNames,
	}

	// grpcCall is a closure that performs the gRPC request to fetch tick values, lot sizes, and contract sizes.
	grpcCall := func(headers metadata.MD) (*pb.TickValueWithSizeReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.TickValueWithSize(c, req)
	}

	// errorSelector inspects the response and extracts the error if there is one.
	errorSelector := func(reply *pb.TickValueWithSizeReply) mrpcError { return reply.GetError() }

	// Execute the gRPC request with automatic retries in case of temporary network or session failures.
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return the data part of the response, which contains the tick value and size calculations.
	if reply == nil || reply.GetData() == nil {
		return &pb.TickValueWithSizeData{}, nil
	}
	return reply.GetData(), nil
}


// === 📂 Symbol State (visibility, index, etc.) ===

// SymbolSelect sets the visibility of a symbol in Market Watch (show/hide).
// SymbolSelect allows enabling or disabling the visibility of a symbol in Market Watch,
// which means showing or hiding the symbol from the user's watchlist.
func (a *MT5Account) SymbolSelect(ctx context.Context, symbol string, enable bool) (bool, error) {
	// Guard connection and inputs
	if !a.isConnected() {
		return false, fmt.Errorf("not connected to terminal")
	}
	if symbol == "" {
		return false, fmt.Errorf("symbol is empty")
	}

	// Safe ctx + short timeout if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// Build request
	req := &pb.SymbolSelectRequest{
		Symbol: symbol,
		Select: enable,
	}

	// RPC call
	grpcCall := func(headers metadata.MD) (*pb.SymbolSelectReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolSelect(c, req)
	}
	errorSelector := func(reply *pb.SymbolSelectReply) mrpcError { return reply.GetError() }

	// Execute
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return false, err
	}
	if reply == nil || reply.GetData() == nil {
		return false, fmt.Errorf("empty SymbolSelectData")
	}
	return reply.GetData().GetSuccess(), nil
}


// IsVisible checks if the specified symbol is visible in Market Watch.
func (a *MT5Account) IsVisible(ctx context.Context, symbol string) (bool, error) {
	// First, ensure the account is connected to the terminal.
	if !a.isConnected() {
		return false, fmt.Errorf("not connected to terminal")
	}
	if symbol == "" {
		return false, fmt.Errorf("symbol is empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// Request the total number of visible symbols (selected/visible symbols).
	totalReq := &pb.SymbolsTotalRequest{Mode: true} // 'Mode: true' means we only want visible (selected) symbols.

	// grpcCallTotal sends the request to get the total number of visible symbols.
	grpcCallTotal := func(headers metadata.MD) (*pb.SymbolsTotalReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolsTotal(c, totalReq)
	}

	// errorSelectorTotal checks the response for errors.
	errorSelectorTotal := func(reply *pb.SymbolsTotalReply) mrpcError { return reply.GetError() }

	// Execute the request to get the total number of visible symbols.
	totalReply, err := ExecuteWithReconnect(a, ctx, grpcCallTotal, errorSelectorTotal)
	if err != nil {
		return false, err // Return error if there was an issue with the request.
	}
	if totalReply == nil || totalReply.GetData() == nil {
		return false, nil
	}

	// Get the total count of visible symbols.
	total := totalReply.GetData().GetTotal()
	if total <= 0 {
		return false, nil // If no visible symbols are found, return false.
	}

	// Iterate through all visible symbols to check if the requested symbol is one of them.
	for i := int32(0); i < total; i++ {
		nameReq := &pb.SymbolNameRequest{Index: i, Selected: true} // Request the name of each visible symbol.

		// grpcCallName sends the request to get the name of the visible symbol at index 'i'.
		grpcCallName := func(headers metadata.MD) (*pb.SymbolNameReply, error) {
			c := metadata.NewOutgoingContext(ctx, headers)
			return a.MarketInfoClient.SymbolName(c, nameReq)
		}

		// errorSelectorName checks for errors in the symbol name response.
		errorSelectorName := func(reply *pb.SymbolNameReply) mrpcError { return reply.GetError() }

		// Execute the request to get the name of the visible symbol.
		nameReply, err := ExecuteWithReconnect(a, ctx, grpcCallName, errorSelectorName)
		if err != nil {
			return false, err // Return error if there was an issue with the request.
		}

		// If the symbol name matches the requested symbol, return true (visible).
		if data := nameReply.GetData(); data != nil && data.GetName() == symbol {
			return true, nil
		}
	}

	// If no matching symbol is found, return false.
	return false, nil
}

// EnsureSymbolVisible ensures that the specified symbol is visible in the Market Watch.
func (a *MT5Account) EnsureSymbolVisible(ctx context.Context, symbol string) error {
	// Step 1: Check if the symbol is currently visible in Market Watch.
	vis, err := a.IsVisible(ctx, symbol)
	if err != nil {
		return err // If there was an error checking visibility, return it.
	}

	// Step 2: If the symbol is not visible, try to make it visible.
	if !vis {
		ok, err := a.SymbolSelect(ctx, symbol, true) // Select the symbol to make it visible.
		if err != nil {
			return err // If there was an error selecting the symbol, return it.
		}
		if !ok {
			return fmt.Errorf("failed to select symbol %s", symbol) // Return an error if the selection fails.
		}
	}

	// Step 3: Return nil if everything went successfully, i.e., the symbol is now visible.
	return nil
}

// SymbolsTotal returns the number of available symbols in the terminal.
func (a *MT5Account) SymbolsTotal(ctx context.Context) (int32, error) {
	// 1) Check if the account is connected to the terminal.
	if !a.IsConnected() {
		return 0, errors.New("not connected to terminal")
	}

	// 2) Ensure ctx is non-nil.
	if ctx == nil {
		ctx = context.Background()
	}
	// 2a) ✅ Added per-call timeout (3s) if none is provided.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// 3) Prepare request.
	req := &pb.SymbolsTotalRequest{} // all symbols

	// 4) Execute RPC with reconnect wrapper.
	grpcCall := func(headers metadata.MD) (*pb.SymbolsTotalReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolsTotal(c, req)
	}
	errorSelector := func(reply *pb.SymbolsTotalReply) mrpcError { return reply.GetError() }

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return 0, err
	}

	// 5) Validate reply.
	if reply == nil || reply.GetData() == nil {
		return 0, nil
	}
	return reply.GetData().GetTotal(), nil
}


// SymbolName returns the symbol name by index in the symbols list.
func (a *MT5Account) SymbolName(ctx context.Context, index int32, selectedOnly bool) (string, error) {
	// 1) Check if the account is connected to the terminal.
	if !a.IsConnected() {
		return "", fmt.Errorf("not connected to terminal")
	}

	// 2) Ensure ctx is non-nil.
	if ctx == nil {
		ctx = context.Background()
	}
	// 2a) ✅ Added per-call timeout (3s) if none is provided.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// 3) Prepare request.
	req := &pb.SymbolNameRequest{
		Index:    index,
		Selected: selectedOnly,
	}

	// 4) Execute RPC with reconnect wrapper.
	grpcCall := func(headers metadata.MD) (*pb.SymbolNameReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.MarketInfoClient.SymbolName(c, req)
	}
	errorSelector := func(reply *pb.SymbolNameReply) mrpcError { return reply.GetError() }

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return "", err
	}

	// 5) Validate reply.
	data := reply.GetData()
	if data == nil {
		return "", fmt.Errorf("empty SymbolNameData")
	}
	return data.GetName(), nil
}


// === 📂 Streaming ===

// OnTrade subscribes to real-time trade updates (open, close, modify).
//
// Returns:
//   - A receive-only channel of TradeData messages.
//   - A receive-only error channel (closes when the stream ends).
func (a *MT5Account) OnTrade(ctx context.Context) (<-chan *pb.OnTradeData, <-chan error) {
	// Step 1: Check if the account has been connected (i.e., it has a valid ID).
	if a.Id == uuid.Nil {
		// If not connected, return a channel with an error message.
		dataCh := make(chan *pb.OnTradeData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- errors.New("please call Connect method first") // Error if Connect wasn't called.
		}()
		return dataCh, errCh
	}

	// Step 2: Create the request for subscribing to trade updates.
	req := &pb.OnTradeRequest{}

	// Step 3: Define functions to extract error and data from the reply.
	getError := func(reply *pb.OnTradeReply) mrpcError { return reply.GetError() }

	getData := func(reply *pb.OnTradeReply) (*pb.OnTradeData, bool) {
		data := reply.GetData()  // Extract the data from the reply.
		return data, data != nil // Return the data if it's not nil.
	}
	newReply := func() *pb.OnTradeReply { return new(pb.OnTradeReply) } // Create a new reply object.

	// Step 4: Call the ExecuteStreamWithReconnect method to start the stream.
	return ExecuteStreamWithReconnect(
		ctx, a, req, // Context and request.
		func(r *pb.OnTradeRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			return a.SubscriptionClient.OnTrade(metadata.NewOutgoingContext(ctx, md), r) // Subscribe to trade updates.
		},
		getError, getData, newReply, // Functions to process the reply data and errors.
	)
}

// OnOpenedOrdersProfit subscribes to periodic updates of profits for all open orders.
//
// Parameters:
//   - ctx: Context for cancel/timeout.
//   - intervalMs: Update interval in milliseconds (e.g., 1000 = 1s).
//
// Returns:
//   - A receive-only channel of OnOpenedOrdersProfitData updates.
//   - A receive-only error channel.
func (a *MT5Account) OnOpenedOrdersProfit(ctx context.Context, intervalMs int32) (<-chan *pb.OnPositionProfitData, <-chan error) {
	// Step 1: Check if the account has been connected (i.e., it has a valid ID).
	if a.Id == uuid.Nil {
		// If not connected, return an error message.
		dataCh := make(chan *pb.OnPositionProfitData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- fmt.Errorf("please call Connect method first")
		}()
		return dataCh, errCh
	}

	// Step 2: Prepare the request for periodic profit updates with the specified interval.
	req := &pb.OnPositionProfitRequest{
		TimerPeriodMilliseconds: intervalMs, // Time period for updates (in ms).
		IgnoreEmptyData:         true,       // Ignore empty data to avoid sending unnecessary frames.
	}

	// Step 3: Define functions to extract error and data from the reply.
	getError := func(reply *pb.OnPositionProfitReply) mrpcError { return reply.GetError() }

	getData := func(reply *pb.OnPositionProfitReply) (*pb.OnPositionProfitData, bool) {
		d := reply.GetData() // Extract profit data from the reply.
		return d, d != nil   // Return data if it's not nil.
	}
	newReply := func() *pb.OnPositionProfitReply { return new(pb.OnPositionProfitReply) } // Create a new reply object.

	// Step 4: Call the ExecuteStreamWithReconnect method to start the stream and handle the connection.
	return ExecuteStreamWithReconnect(
		ctx, a, req, // Context and request.
		func(r *pb.OnPositionProfitRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			// Subscribe to position profit updates via gRPC stream.
			return a.SubscriptionClient.OnPositionProfit(metadata.NewOutgoingContext(ctx, md), r)
		},
		getError, getData, newReply, // Functions to process the reply data and errors.
	)
}

// OnOpenedOrdersTickets subscribes to periodic updates of opened order ticket IDs.
//
// Parameters:
//   - ctx: Context for cancel/timeout.
//   - intervalMs: Update interval in milliseconds.
//
// Returns:
//   - A receive-only channel of OnOpenedOrdersTicketsData updates.
//   - A receive-only error channel.
func (a *MT5Account) OnOpenedOrdersTickets(
	ctx context.Context,
	intervalMs int32,
) (<-chan *pb.OnPositionsAndPendingOrdersTicketsData, <-chan error) {
	// Step 1: Check if the account has been connected (i.e., it has a valid ID).
	if a.Id == uuid.Nil {
		// If not connected, return an error message.
		dataCh := make(chan *pb.OnPositionsAndPendingOrdersTicketsData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- fmt.Errorf("please call Connect method first")
		}()
		return dataCh, errCh
	}

	// Step 2: Prepare the request for periodic ticket updates with the specified interval.
	req := &pb.OnPositionsAndPendingOrdersTicketsRequest{
		TimerPeriodMilliseconds: intervalMs, // Time period for updates (in ms).
	}

	// Step 3: Define functions to extract error and data from the reply.
	getError := func(reply *pb.OnPositionsAndPendingOrdersTicketsReply) mrpcError { return reply.GetError() }

	getData := func(reply *pb.OnPositionsAndPendingOrdersTicketsReply) (*pb.OnPositionsAndPendingOrdersTicketsData, bool) {
		d := reply.GetData() // Extract ticket data from the reply.
		return d, d != nil   // Return data if it's not nil.
	}
	newReply := func() *pb.OnPositionsAndPendingOrdersTicketsReply {
		return new(pb.OnPositionsAndPendingOrdersTicketsReply)
	} // Create a new reply object.

	// Step 4: Call the ExecuteStreamWithReconnect method to start the stream and handle the connection.
	return ExecuteStreamWithReconnect(
		ctx, a, req, // Context and request.
		func(r *pb.OnPositionsAndPendingOrdersTicketsRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			// Subscribe to position and pending orders ticket updates via gRPC stream.
			return a.SubscriptionClient.OnPositionsAndPendingOrdersTickets(metadata.NewOutgoingContext(ctx, md), r)
		},
		getError, getData, newReply, // Functions to process the reply data and errors.
	)
}

// OnSymbolTick subscribes to real-time tick data for specified symbols, with reconnection logic.
//
// Parameters:
//   - ctx: Context for cancellation/timeouts
//   - symbols: Slice of symbol names (e.g., []string{"EURUSD", "USDJPY"})
//
// Returns:
//   - Receive-only channel of *pb.OnSymbolTickData (each tick)
//   - Receive-only error channel
func (a *MT5Account) OnSymbolTick(
	ctx context.Context,
	symbols []string,
) (<-chan *pb.OnSymbolTickData, <-chan error) {
	// Check that the account is connected (has a valid UUID)
	if a.Id == uuid.Nil {
		dataCh := make(chan *pb.OnSymbolTickData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- errors.New("please call Connect method first")
		}()
		return dataCh, errCh
	}

	// Build the request message for the stream
	req := &pb.OnSymbolTickRequest{SymbolNames: symbols}

	// Function to extract API error from the proto reply
	getError := func(reply *pb.OnSymbolTickReply) mrpcError { return reply.GetError() }

	// Function to extract the tick data (returns (data, ok))
	getData := func(reply *pb.OnSymbolTickReply) (*pb.OnSymbolTickData, bool) {
		data := reply.GetData()
		return data, data != nil
	}

	// The "newReply" function returns a new pointer to your proto reply type.
	newReply := func() *pb.OnSymbolTickReply { return new(pb.OnSymbolTickReply) }

	// Call the generic streaming helper.
	dataCh, errCh := ExecuteStreamWithReconnect(
		ctx, a, req,
		func(r *pb.OnSymbolTickRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			return a.SubscriptionClient.OnSymbolTick(metadata.NewOutgoingContext(ctx, md), r)
		},
		getError,
		getData,
		newReply,
	)

	return dataCh, errCh
}

