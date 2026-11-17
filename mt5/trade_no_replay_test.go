package mt5

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// fakeTrading answers OrderSend/OrderModify/OrderClose with a scripted failure and counts calls.
type fakeTrading struct {
	pb.UnimplementedTradingHelperServer
	calls   atomic.Int32
	answer  func(call int32) (*pb.OrderSendReply, error)
	transFn func() error // OrderModify / OrderClose transport error
}

func (f *fakeTrading) OrderSend(context.Context, *pb.OrderSendRequest) (*pb.OrderSendReply, error) {
	return f.answer(f.calls.Add(1))
}

func (f *fakeTrading) OrderModify(context.Context, *pb.OrderModifyRequest) (*pb.OrderModifyReply, error) {
	f.calls.Add(1)
	return nil, f.transFn()
}

func (f *fakeTrading) OrderClose(context.Context, *pb.OrderCloseRequest) (*pb.OrderCloseReply, error) {
	f.calls.Add(1)
	return nil, f.transFn()
}

// fakeAccount fails AccountSummary once with UNAVAILABLE, then answers.
type fakeAccount struct {
	pb.UnimplementedAccountHelperServer
	calls atomic.Int32
}

func (f *fakeAccount) AccountSummary(context.Context, *pb.AccountSummaryRequest) (*pb.AccountSummaryReply, error) {
	if f.calls.Add(1) == 1 {
		return nil, status.Error(codes.Unavailable, "connection reset")
	}
	return &pb.AccountSummaryReply{Response: &pb.AccountSummaryReply_Data{Data: &pb.AccountSummaryData{AccountLogin: 42}}}, nil
}

func newFakeMT5Account(t *testing.T, trading *fakeTrading, account *fakeAccount) *MT5Account {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()
	if trading != nil {
		pb.RegisterTradingHelperServer(srv, trading)
	}
	if account != nil {
		pb.RegisterAccountHelperServer(srv, account)
	}
	go func() { _ = srv.Serve(lis) }()
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close(); srv.Stop() })
	return &MT5Account{
		GrpcConn:      conn,
		TradeClient:   pb.NewTradingHelperClient(conn),
		AccountClient: pb.NewAccountHelperClient(conn),
		Id:            uuid.New(),
		ApiKey:        "test-key",
	}
}

func assertTradeOutcomeUnknown(t *testing.T, err error, op string, code codes.Code) {
	t.Helper()
	var tou *TradeOutcomeUnknownError
	if !errors.As(err, &tou) {
		t.Fatalf("expected *TradeOutcomeUnknownError, got %T: %v", err, err)
	}
	if tou.Code != "TRADE_OUTCOME_UNKNOWN" || tou.ErrorCode() != "TRADE_OUTCOME_UNKNOWN" || tou.Operation != op {
		t.Fatalf("bad fields: code=%q op=%q", tou.Code, tou.Operation)
	}
	if !strings.HasPrefix(err.Error(), "TRADE_OUTCOME_UNKNOWN: "+op+" failed with gRPC ") && !strings.Contains(err.Error(), "TRADE_OUTCOME_UNKNOWN: "+op) {
		t.Fatalf("bad message: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "NOT re-sent") {
		t.Fatalf("message does not say it was not re-sent: %s", err.Error())
	}
	if got := status.Code(err); got != code {
		t.Fatalf("status.Code = %s, want %s", got, code)
	}
	if s, ok := status.FromError(err); !ok || s.Code() != code {
		t.Fatalf("status.FromError = %v,%v", s, ok)
	}
	if errors.Unwrap(tou) == nil {
		t.Fatal("Unwrap returned nil")
	}
}

func TestOrderSendUnavailableIsNotReplayed(t *testing.T) {
	f := &fakeTrading{answer: func(int32) (*pb.OrderSendReply, error) {
		return nil, status.Error(codes.Unavailable, "socket closed")
	}}
	a := newFakeMT5Account(t, f, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := a.OrderSend(ctx, &pb.OrderSendRequest{Symbol: "EURUSD", Volume: 0.01})
	assertTradeOutcomeUnknown(t, err, "OrderSend", codes.Unavailable)
	if n := f.calls.Load(); n != 1 {
		t.Fatalf("OrderSend called %d times, want 1", n)
	}

	// The Service layer wraps it with %w; the typed error must still be found.
	_, err = NewMT5Service(a).PlaceOrder(ctx, &pb.OrderSendRequest{Symbol: "EURUSD", Volume: 0.01})
	assertTradeOutcomeUnknown(t, err, "OrderSend", codes.Unavailable)
	if n := f.calls.Load(); n != 2 {
		t.Fatalf("OrderSend called %d times after PlaceOrder, want 2", n)
	}
}

func TestOrderModifyAndCloseOutcomeUnknownCodes(t *testing.T) {
	for _, c := range []codes.Code{codes.Unavailable, codes.DeadlineExceeded, codes.Unknown, codes.Internal} {
		f := &fakeTrading{transFn: func() error { return status.Error(c, "boom") }}
		a := newFakeMT5Account(t, f, nil)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := a.OrderModify(ctx, &pb.OrderModifyRequest{Ticket: 1})
		assertTradeOutcomeUnknown(t, err, "OrderModify", c)
		_, err = a.OrderClose(ctx, &pb.OrderCloseRequest{Ticket: 1, Volume: 0.01})
		assertTradeOutcomeUnknown(t, err, "OrderClose", c)
		cancel()
		if n := f.calls.Load(); n != 2 {
			t.Fatalf("%s: trade calls = %d, want 2 (one each)", c, n)
		}
	}
}

func TestOrderSendTerminalNotFoundIsResent(t *testing.T) {
	f := &fakeTrading{answer: func(call int32) (*pb.OrderSendReply, error) {
		if call == 1 {
			return &pb.OrderSendReply{Response: &pb.OrderSendReply_Error{Error: &pb.Error{ErrorCode: "TERMINAL_INSTANCE_NOT_FOUND"}}}, nil
		}
		return &pb.OrderSendReply{Response: &pb.OrderSendReply_Data{Data: &pb.OrderSendData{ReturnedCode: 10009, Order: 7}}}, nil
	}}
	a := newFakeMT5Account(t, f, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := a.OrderSend(ctx, &pb.OrderSendRequest{Symbol: "EURUSD", Volume: 0.01})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.GetOrder() != 7 {
		t.Fatalf("unexpected data: %v", data)
	}
	if n := f.calls.Load(); n != 2 {
		t.Fatalf("OrderSend called %d times, want 2", n)
	}
}

func TestReadUnavailableIsRetried(t *testing.T) {
	acc := &fakeAccount{}
	a := newFakeMT5Account(t, nil, acc)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := a.AccountSummary(ctx, &pb.AccountSummaryRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.GetAccountLogin() != 42 {
		t.Fatalf("unexpected data: %v", data)
	}
	if n := acc.calls.Load(); n != 2 {
		t.Fatalf("AccountSummary called %d times, want 2", n)
	}
}

func TestApiKeyInvalidFailsFast(t *testing.T) {
	f := &fakeTrading{answer: func(int32) (*pb.OrderSendReply, error) {
		return &pb.OrderSendReply{Response: &pb.OrderSendReply_Error{Error: &pb.Error{ErrorCode: "API_KEY_INVALID", ErrorMessage: "unknown api key"}}}, nil
	}}
	a := newFakeMT5Account(t, f, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := a.OrderSend(ctx, &pb.OrderSendRequest{Symbol: "EURUSD", Volume: 0.01})
	var apiErr *ApiError
	if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "API_KEY_INVALID" {
		t.Fatalf("expected ApiError API_KEY_INVALID, got %T: %v", err, err)
	}
	if n := f.calls.Load(); n != 1 {
		t.Fatalf("OrderSend called %d times, want 1", n)
	}

	// Refused key carried in a transport status message: no retry for reads either.
	g := &fakeTrading{answer: func(int32) (*pb.OrderSendReply, error) {
		return nil, status.Error(codes.Unavailable, "API_KEY_REQUIRED: no api key")
	}}
	b := newFakeMT5Account(t, g, nil)
	_, err = ExecuteWithReconnect(b, ctx, func(md metadata.MD) (*pb.OrderSendReply, error) {
		return b.TradeClient.OrderSend(metadata.NewOutgoingContext(ctx, md), &pb.OrderSendRequest{})
	}, func(r *pb.OrderSendReply) mrpcError { return r.GetError() })
	if err == nil || !strings.Contains(err.Error(), "API_KEY_REQUIRED") {
		t.Fatalf("expected API_KEY_REQUIRED error, got %v", err)
	}
	if n := g.calls.Load(); n != 1 {
		t.Fatalf("call made %d times, want 1", n)
	}
}
