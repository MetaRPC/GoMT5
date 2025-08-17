package main

import (
	"context"
	"log"
	"os"
"os/signal"
   "sync"
   "syscall"
	"time"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"

	"github.com/MetaRPC/GoMT5/mt5"
	"github.com/google/uuid"
)

// waitReady polls terminal liveness until deadline.
// Safe to call right after Connect* even if it returned an error.
func waitReady(ctx context.Context, acc *mt5.MT5Account, maxWait time.Duration) bool {
	deadline := time.Now().Add(maxWait)
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()

	for attempt := 1; ; attempt++ {
		ok, _ := acc.IsTerminalAlive(ctx)
		if ok {
			log.Println("terminal is ready")
			return true
		}
		select {
		case <-ctx.Done():
			log.Printf("terminal is not ready: %v", ctx.Err())
			return false
		case <-t.C:
			log.Printf("still waiting... attempt %d, time left: %s",
				attempt, time.Until(deadline).Truncate(time.Second))
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Umbrella context for the whole run (keep it simple).
	rootCtx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	// --- Basic settings (edit these as needed) ---
	const login uint64 = 501401178
	const password = "v8gctta"
	defaultSymbol := "EURUSD"

	servers := []string{
		"RoboForex-Demo",
		"RoboForex-Demo 2",
		"RoboForex-Pro",
		"RoboForex-Pro-2",
		"RoboForex-ECN",
		"RoboForex-ECN-Pro",
		"MetaQuotes-Demo", // diagnostic fallback
	}

	const serverSideWait = 480           // how long terminal may block during login
	const localReadinessWait = 8 * time.Minute

	// Optional proxy: socks5://user:pass@host:port or http://host:port
	proxy := os.Getenv("MT5_PROXY")
	if proxy != "" {
		log.Printf("using proxy: %s", proxy)
	} else {
		log.Printf("no proxy set (set MT5_PROXY if your egress requires it)")
	}

	// --- Connect (simple, readable retry over servers) ---
	newAccount := func() *mt5.MT5Account {
		acc, err := mt5.NewMT5Account(login, password, "", uuid.Nil)
		if err != nil {
			log.Fatalf("NewMT5Account: %v", err)
		}
		return acc
	}

	acc := newAccount()
	defer func() { _ = acc.Disconnect(context.Background()) }()

	var connectedSrv string
	for i, srv := range servers {
		if i > 0 {
			_ = acc.Disconnect(context.Background())
			acc = newAccount()
		}
		log.Printf("connect(wait on server %ds) to %s ...", serverSideWait, srv)

		// Even if Connect returns an error, terminal may continue initializing — keep waiting below.
		_ = acc.ConnectByServerName(rootCtx, srv, proxy, true, serverSideWait)

		time.Sleep(2 * time.Second) // small breath
		if waitReady(rootCtx, acc, localReadinessWait) {
			connectedSrv = srv
			break
		}
	}
	if connectedSrv == "" {
		log.Fatal("failed to connect to any server. If behind proxy/firewall — set MT5_PROXY or allow egress.")
	}
	log.Printf("connected to: %s", connectedSrv)

	// --- Prepare symbol (try common suffixes; pick first visible) ---
	selected := defaultSymbol
	for _, s := range []string{defaultSymbol, defaultSymbol + ".", defaultSymbol + ".pro"} {
		if s == "" {
			continue
		}
		if err := acc.EnsureSymbolVisible(rootCtx, s); err == nil {
			selected = s
			log.Printf("symbol ready: %s", s)
			break
		} else {
			log.Printf("EnsureSymbolVisible(%s): %v", s, err)
		}
	}

	// --- Service facade tied to the connected account ---
	svc := mt5.NewMT5Service(acc)

	// ===============================
	// CALLS IN A CLEAR, SAFE ORDER
	// ===============================

	// 1) 📂 Quick account & market info (read-only)
	svc.ShowAccountSummary(rootCtx)
	svc.ShowQuote(rootCtx, selected)
	svc.ShowQuotesMany(rootCtx, []string{selected, "GBPUSD"})
	svc.ShowSymbolParams(rootCtx, selected)
	svc.ShowTickValues(rootCtx, []string{selected, "GBPUSD"})
	svc.ShowAllSymbols(rootCtx) // may print a lot

	// 2) 📂 Opened state snapshot (read-only)
	svc.ShowOpenedOrders(rootCtx)
	svc.ShowOpenedOrderTickets(rootCtx)
	svc.ShowPositions(rootCtx)
	svc.ShowHasOpenPosition(rootCtx, selected)

	// 3) 📂 Calculations & a dry-run check (still safe)
	svc.ShowOrderCalcMargin(rootCtx, selected, pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, 0.10, 0)
	svc.ShowOrderCalcProfit(rootCtx, selected, pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY, 0.10, 1.08000, 1.08350)
	svc.ShowOrderCheck(rootCtx,
		pb.MRPC_ENUM_TRADE_REQUEST_ACTIONS_TRADE_ACTION_DEAL,
		pb.ENUM_ORDER_TYPE_TF_ORDER_TYPE_TF_BUY,
		selected, 0.10, 0, // price=0 means "market" on many servers
		nil, nil, nil, nil, nil)

	// 4) 📂 Trading ops (DANGEROUS) — keep commented until you really want them.
	// svc.ShowOrderSendExample(rootCtx, selected)
	// svc.ShowOrderSendStopLimitExample(rootCtx, selected, true, 1.09100, 1.09120)
	// svc.BuyMarket(rootCtx, selected, 0.10, nil, nil)
	// svc.SellMarket(rootCtx, selected, 0.10, nil, nil)
	// exp := timestamppb.New(time.Now().Add(24 * time.Hour))

	//--------------------------------------------------------------

	// svc.PlaceBuyLimit(rootCtx, selected, 0.10, 1.07500, nil, nil, exp)
	// svc.PlaceSellLimit(rootCtx, selected, 0.10, 1.09500, nil, nil, exp)
	// svc.PlaceBuyStop(rootCtx, selected, 0.10, 1.09200, nil, nil, exp)
	// svc.PlaceSellStop(rootCtx, selected, 0.10, 1.07800, nil, nil, exp)
	// svc.PlaceStopLimit(rootCtx, selected, true, 0.10, 1.09100, 1.09120, nil, nil, exp)

	//--------------------------------------------------------------

	// svc.ShowOrderModifyExample(rootCtx, /*ticket=*/ 123456789)
	// svc.ShowOrderCloseExample(rootCtx,  /*ticket=*/ 123456789)
	// svc.ShowOrderDeleteExample(rootCtx, /*ticket=*/ 123456789)
	// svc.ShowPositionModify(rootCtx,     /*ticket=*/ 987654321, nil, nil)
	// svc.ShowPositionClose(rootCtx, selected)
	// svc.ShowCloseAllPositions(rootCtx) // CAREFUL!

	// 5) 📂 History & simple stats (read-only)
	svc.ShowOrdersHistory(rootCtx)
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now()
	svc.ShowDealsCount(rootCtx, from, to, "")
	// svc.ShowOrderByTicket(rootCtx, 123456789) // optional lookups
	// svc.ShowDealByTicket(rootCtx, 987654321)

	// 6) 📂 Streaming (parallel)
// ❗ Dangers/rakes:
// - Noisy streams: lots of logs → use them purposefully.
// - Don't forget about network limits/broker restrictions.
// - Give a reasonable total timeout so as not to hang indefinitely (example below).

// We are preparing a general context that will be canceled by Ctrl+C
ctx, stop := signal.NotifyContext(rootCtx, os.Interrupt, syscall.SIGTERM)
defer stop()

var wg sync.WaitGroup
run := func(name string, fn func(context.Context)) {
    wg.Add(1)
    go func() {
        defer wg.Done()
        log.Printf("▶ %s started", name)
        fn(ctx)
        log.Printf("■ %s stopped", name)
    }()
}

// Launching multiple streams in parallel
run("StreamQuotes", func(c context.Context) { svc.StreamQuotes(c) })
run("StreamOpenedOrderProfits", func(c context.Context) { svc.StreamOpenedOrderProfits(c) })
run("StreamOpenedOrderTickets", func(c context.Context) { svc.StreamOpenedOrderTickets(c) })
// run("StreamTradeUpdates", func(c context.Context) { svc.StreamTradeUpdates(c) }) // turn it on if necessary

// Global fuse: stop after 45s if there is no Ctrl+C
select {
case <-ctx.Done(): // Ctrl+C or SIGTERM
case <-time.After(45 * time.Second):
    log.Println("⏱ global streaming timeout reached, stopping")
}

// Disabling subscriptions and waiting for the completion of mining
stop()
wg.Wait()

log.Println("✅ Done.")
}
