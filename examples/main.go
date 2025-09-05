package main

import (
	"context"
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
	"github.com/MetaRPC/GoMT5/examples/config"
	"github.com/MetaRPC/GoMT5/mt5"
	"github.com/google/uuid"
)

// waitReady polls terminal liveness until deadline with verbose errors.
// Fast-fails on common auth/network issues to avoid long blind waits.
func waitReady(ctx context.Context, acc *mt5.MT5Account, maxWait time.Duration) bool {
	deadline := time.Now().Add(maxWait)
	t := time.NewTicker(1 * time.Second) // poll faster to see progress
	defer t.Stop()

	for attempt := 1; ; attempt++ {
		ok, err := acc.IsTerminalAlive(ctx)
		if ok {
			log.Println("terminal is ready")
			return true
		}
		if err != nil {
			// Show the exact reason while waiting.
			msg := err.Error()
			log.Printf("liveness probe error: %v (attempt %d, left %s)",
				msg, attempt, time.Until(deadline).Truncate(time.Second))

			// Fast-fail on typical auth issues instead of waiting full timeout.
			if strings.Contains(msg, "Invalid account") ||
				strings.Contains(msg, "authorization failed") {
				log.Printf("fast-fail: authentication problem detected")
				return false
			}
		} else {
			log.Printf("still waiting... attempt %d, time left: %s",
				attempt, time.Until(deadline).Truncate(time.Second))
		}

		select {
		case <-ctx.Done():
			log.Printf("terminal is not ready: %v", ctx.Err())
			return false
		case <-t.C:
			// keep looping
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Umbrella context for the whole run
	rootCtx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	// --- Load config.json ---
	path := os.Getenv("MT5_CONFIG") // optional override
	if path == "" {
		path = "examples/config/config.json"
	}
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatalf("config.json load error: %v", err)
	}
	if cfg.Login == 0 || cfg.Password == "" || cfg.Server == "" {
		log.Fatal("config.json must contain Login, Password, Server")
	}
	log.Printf("cfg: login=%d server=%q symbol=%q", cfg.Login, cfg.Server, cfg.DefaultSymbol)
	if len(cfg.Password) >= 2 {
		log.Printf("cfg: password=***%s", cfg.Password[len(cfg.Password)-2:]) // the last 2 characters to check for a match
	}

	login := uint64(cfg.Login)
	password := cfg.Password
	serverName := cfg.Server
	defaultSymbol := cfg.DefaultSymbol
	if defaultSymbol == "" {
		defaultSymbol = "EURUSD"
	}

	// --- Optional proxy ---
	proxy := os.Getenv("MT5_PROXY")
	if proxy != "" {
		log.Printf("using proxy: %s", proxy)
	} else {
		log.Printf("no proxy set (set MT5_PROXY if your egress requires it)")
	}

	const (
		serverSideWait     = 45               // seconds for broker-side login wait
		localReadinessWait = 30 * time.Second // how long we poll readiness
	)

	// --- Create account object ---
	newAccount := func() *mt5.MT5Account {
		acc, err := mt5.NewMT5Account(login, password, "", uuid.Nil)
		if err != nil {
			log.Fatalf("NewMT5Account: %v", err)
		}
		return acc
	}

	// --- Normalize & debug-print creds (helps catch hidden chars) ---
	{
		// keep imports: "strings", "encoding/hex"
		password = strings.TrimSpace(password)
		serverName = strings.TrimSpace(serverName)

		// Debug: show length and first bytes of password
		pwHex := ""
		{
			b := []byte(password)
			n := len(b)
			if n > 8 {
				n = 8
			}
			pwHex = hex.EncodeToString(b[:n])
		}

		log.Printf("cfg: login=%d server=%q symbol=%q", login, serverName, defaultSymbol)
		log.Printf("cfg: password len=%d first8bytes=%s", len(password), pwHex)
	}

	acc := newAccount()
	defer func() { _ = acc.Disconnect(context.Background()) }()

	// --- Connect by server name from config ---
	log.Printf("connect(wait on server %ds) to %s ...", serverSideWait, serverName)
	if err := acc.ConnectByServerName(rootCtx, serverName, defaultSymbol, true, serverSideWait); err != nil {
		log.Fatalf("ConnectByServerName(%s) error: %v", serverName, err)
	}
	time.Sleep(2 * time.Second)
	if !waitReady(rootCtx, acc, localReadinessWait) {
		log.Fatal("terminal is not ready after connect")
	}
	log.Printf("connected to: %s (login=%d)", serverName, login)

	// --- Ensure symbol is visible ---
	selected := defaultSymbol
	for _, s := range []string{defaultSymbol, defaultSymbol + ".", defaultSymbol + ".pro"} {
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

	// Disabling subscriptions and waiting for the completion of draining
	stop()
	wg.Wait()

	log.Println("✅ Done.")
}
