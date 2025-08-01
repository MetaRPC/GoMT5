package mt4go

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"crypto/tls"

	pb "git.mtapi.io/root/mrpc-proto.git/mt4/libraries/go" // Protobuf-generated Go package for MT4 API
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// MT4Account represents a client session for interacting with the MT4 terminal API over gRPC.
type MT4Account struct {
	// User is the MT4 account login number.
	User uint64

	// Password for the user account.
	Password string

	// Host is the IP/domain of the MT4 server.
	Host string

	// Port is the MT4 server port (typically 443).
	Port int

	// ServerName is the MT4 server name (used for cluster connections).
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
	ConnectionClient   pb.ConnectionClient
	SubscriptionClient pb.SubscriptionServiceClient
	AccountClient      pb.AccountHelperClient
	TradeClient        pb.TradingHelperClient
	MarketInfoClient   pb.MarketInfoClient

	// Id is a unique identifier (UUID) for this account session/instance.
	Id uuid.UUID
}

// NewMT4Account initializes a new MT4Account and establishes the underlying gRPC connection.
// Returns a pointer to the account object and any error encountered while connecting.
func NewMT4Account(user uint64, password string, grpcServer string, id uuid.UUID) (*MT4Account, error) {
	// If no endpoint specified, use production default
	if grpcServer == "" {
		grpcServer = "mt4.mrpc.pro:443"
	}

	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(grpcServer, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return nil, err
	}

	// Instantiate API service clients using the shared gRPC connection
	return &MT4Account{
		User:               user,
		Password:           password,
		GrpcServer:         grpcServer,
		GrpcConn:           conn,
		ConnectionClient:   pb.NewConnectionClient(conn),
		SubscriptionClient: pb.NewSubscriptionServiceClient(conn),
		AccountClient:      pb.NewAccountHelperClient(conn),
		TradeClient:        pb.NewTradingHelperClient(conn),
		MarketInfoClient:   pb.NewMarketInfoClient(conn),
		Id:                 id,
		Port:               443,
		ConnectTimeout:     30,
	}, nil
}

// isConnected returns true if this account is associated with any host or server name.
func (a *MT4Account) isConnected() bool {
	return a.Host != "" || a.ServerName != ""
}

// getHeaders builds the gRPC metadata headers (adds "id" if present).
func (a *MT4Account) getHeaders() metadata.MD {
	if a.Id == uuid.Nil {
		return nil
	}
	return metadata.Pairs("id", a.Id.String())
}

// ConnectByHostPort connects to the MT4 terminal using a host/port pair.
// Updates the session state fields (Host, Port, etc.) upon success.
func (a *MT4Account) ConnectByHostPort(
	ctx context.Context,
	host string,
	port int,
	baseChartSymbol string,
	waitForTerminalIsAlive bool,
	timeoutSeconds int,
) error {
	// Build the protobuf request struct
	req := &pb.ConnectRequest{
		User:                                   a.User,
		Password:                               a.Password,
		Host:                                   host,
		Port:                                   int32(port),
		BaseChartSymbol:                        proto.String(baseChartSymbol),
		WaitForTerminalIsAlive:                 proto.Bool(waitForTerminalIsAlive),
		TerminalReadinessWaitingTimeoutSeconds: proto.Int32(int32(timeoutSeconds)),
	}
	// Set metadata if available
	md := a.getHeaders()
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Make the actual gRPC call
	res, err := a.ConnectionClient.Connect(ctx, req)
	if err != nil {
		return err
	}
	// API errors are delivered via .GetError()
	if res.GetError() != nil {
		return fmt.Errorf("API error: %v", res.GetError())
	}

	// Store session properties if connection is established
	a.Host = host
	a.Port = port
	a.BaseChartSymbol = baseChartSymbol
	a.ConnectTimeout = timeoutSeconds

	// Set the session UUID if present in response
	if data := res.GetData(); data != nil && data.GetTerminalInstanceGuid() != "" {
		id, _ := uuid.Parse(data.GetTerminalInstanceGuid())
		a.Id = id
	}
	return nil
}

// ConnectByServerName connects to the MT4 terminal using the cluster/server name.
func (a *MT4Account) ConnectByServerName(
	ctx context.Context,
	serverName string,
	baseChartSymbol string,
	waitForTerminalIsAlive bool,
	timeoutSeconds int,
) error {
	req := &pb.ConnectExRequest{
		User:                                   a.User,
		Password:                               a.Password,
		MtClusterName:                          serverName,
		BaseChartSymbol:                        proto.String(baseChartSymbol),
		TerminalReadinessWaitingTimeoutSeconds: proto.Int32(int32(timeoutSeconds)),
	}
	md := a.getHeaders()
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, err := a.ConnectionClient.ConnectEx(ctx, req)
	if err != nil {
		return err
	}
	if res.GetError() != nil {
		return fmt.Errorf("API error: %v", res.GetError())
	}
	a.ServerName = serverName
	a.BaseChartSymbol = baseChartSymbol
	a.ConnectTimeout = timeoutSeconds

	if data := res.GetData(); data != nil && data.GetTerminalInstanceGuid() != "" {
		id, _ := uuid.Parse(data.GetTerminalInstanceGuid())
		a.Id = id
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
func ExecuteWithReconnect[T any](
	a *MT4Account,
	ctx context.Context,
	grpcCall func(metadata.MD) (T, error),
	errorSelector func(T) *pb.Error,
) (T, error) {
	var zeroT T // Zero value of T (used for error returns)

	for {
		// Prepare gRPC headers for session (may be nil)
		headers := a.getHeaders()

		// Call the gRPC method (returns (reply, error))
		res, err := grpcCall(headers)
		if err != nil {
			// If it's a gRPC Unavailable (connection/server issue), wait and retry
			if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
				select {
				case <-time.After(500 * time.Millisecond):
					continue // Try again after delay
				case <-ctx.Done():
					return zeroT, ctx.Err() // Cancelled by caller
				}
			}
			// Other errors: return immediately
			return zeroT, err
		}

		// Check for API (business logic) error in the response
		apiErr := errorSelector(res)
		if apiErr != nil {
			// If terminal instance not found (e.g., dropped session), wait and retry
			if apiErr.GetErrorCode() == "TERMINAL_INSTANCE_NOT_FOUND" {
				select {
				case <-time.After(500 * time.Millisecond):
					continue // Try again after delay
				case <-ctx.Done():
					return zeroT, ctx.Err()
				}
			}
			// All other API errors: return as Go errors
			return zeroT, fmt.Errorf("API error: %v", apiErr)
		}

		// Success! Return the response.
		return res, nil
	}
}

// AccountSummary retrieves summary information about the connected MT4 trading account.
//
// Parameters:
//   - ctx: Context for timeout or cancellation (e.g., context.Background(), context.WithTimeout).
//
// Returns:
//   - Pointer to AccountSummaryData (or nil if error).
//   - Error if not connected or if the gRPC/API call fails.
//
// This method handles automatic retries on network or "terminal instance not found" errors
// using ExecuteWithReconnect. It accesses protobuf fields via generated Get...() methods.
func (a *MT4Account) AccountSummary(ctx context.Context) (*pb.AccountSummaryData, error) {
	// Ensure the account is connected to a server before making a request.
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Construct the empty request message (no parameters for account summary).
	req := &pb.AccountSummaryRequest{}

	// grpcCall is a closure that performs the gRPC request.
	// It takes the request metadata (headers), creates a context with it,
	// and calls the generated AccountSummary client method.
	grpcCall := func(headers metadata.MD) (*pb.AccountSummaryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.AccountSummary(c, req)
	}

	// errorSelector inspects the reply for any application-level error.
	// It extracts the API error from the response (if present), or returns nil for success.
	errorSelector := func(reply *pb.AccountSummaryReply) *pb.Error {
		return reply.GetError()
	}

	// ExecuteWithReconnect wraps the gRPC call in retry/reconnect logic.
	// It will retry on certain recoverable errors (network/server unavailable, terminal instance lost).
	// If successful, reply will be non-nil and err will be nil.
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		// Return nil and the error to the caller.
		return nil, err
	}

	// Return the data portion of the reply (may be nil if server returned no data).
	// Always use the generated GetData() method for proto oneof fields.
	return reply.GetData(), nil
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
func ExecuteStreamWithReconnect[TRequest any, TReply any, TData any](
	ctx context.Context,
	a *MT4Account,
	request TRequest,
	streamInvoker func(TRequest, metadata.MD, context.Context) (grpc.ClientStream, error),
	getError func(TReply) *pb.Error,
	getData func(TReply) (TData, bool),
	newReply func() TReply, // <-- IMPORTANT: pass e.g. func() *pb.OnSymbolTickReply { return new(pb.OnSymbolTickReply) }
) (<-chan TData, <-chan error) {
	dataCh := make(chan TData)
	errCh := make(chan error, 1)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		for {
			reconnectRequired := false
			headers := a.getHeaders()
			// Open the gRPC streaming call with headers and context
			stream, err := streamInvoker(request, headers, ctx)
			if err != nil {
				// If network/server unavailable, retry unless cancelled
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
					select {
					case <-time.After(500 * time.Millisecond):
						continue // retry connection
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
				}
				errCh <- err
				return
			}

			for {
				// Create a new empty reply message of the correct type
				reply := newReply() // <- always returns pointer to proto message

				// Receive a message from the stream (unmarshals from wire into reply)
				recvErr := stream.RecvMsg(reply)
				if recvErr != nil {
					// Network/server error: attempt reconnect if recoverable
					if s, ok := status.FromError(recvErr); ok && s.Code() == codes.Unavailable {
						reconnectRequired = true
						break // break inner loop, retry stream
					}
					if errors.Is(recvErr, io.EOF) {
						return // stream ended gracefully (no error)
					}
					// User cancelled or deadline exceeded
					if errors.Is(recvErr, context.Canceled) || errors.Is(recvErr, context.DeadlineExceeded) {
						errCh <- recvErr
						return
					}
					// All other errors: fail and close
					errCh <- recvErr
					return
				}

				// Check for logical/API errors inside the proto reply
				apiErr := getError(reply)
				if apiErr != nil {
					code := apiErr.GetErrorCode()
					// Certain terminal errors are recoverable (reconnect)
					if code == "TERMINAL_INSTANCE_NOT_FOUND" || code == "TERMINAL_REGISTRY_TERMINAL_NOT_FOUND" {
						reconnectRequired = true
						break
					}
					// All other API errors: report and end
					errCh <- fmt.Errorf("API error: %v", apiErr)
					return
				}

				// Extract the real data from reply (skip if not present)
				if data, ok := getData(reply); ok {
					select {
					case dataCh <- data: // send data to caller
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
				}
			}

			// Handle reconnect logic
			if reconnectRequired {
				select {
				case <-time.After(500 * time.Millisecond):
					continue // retry outer loop (reconnect)
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				}
			} else {
				break // Exit outer loop if not reconnecting
			}
		}
	}()
	return dataCh, errCh
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
func (a *MT4Account) OnSymbolTick(
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
	getError := func(reply *pb.OnSymbolTickReply) *pb.Error {
		return reply.GetError()
	}
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
