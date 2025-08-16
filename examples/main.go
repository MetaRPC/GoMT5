package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/MetaRPC/GoMT5/mt5"
	"github.com/google/uuid"
)

type Config struct {
	DefaultSymbol string
}

func waitReady(ctx context.Context, acc *mt5.MT5Account, maxWait time.Duration) bool {
	deadline := time.Now().Add(maxWait)
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()

	for i := 1; ; i++ {
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
				i, time.Until(deadline).Truncate(time.Second))
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Общий «зонтик» на запуск
	rootCtx, cancelRoot := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancelRoot()

	cfg := Config{DefaultSymbol: "EURUSD"}

	const login uint64 = 501401178
	const password = "v8gctta"

	// Прокси читаем из окружения: socks5://user:pass@host:port или http://host:port
	proxy := os.Getenv("MT5_PROXY")
	if proxy != "" {
		log.Printf("using proxy: %s", proxy)
	} else {
		log.Printf("no proxy set (set MT5_PROXY if your egress requires it)")
	}

	servers := []string{
		"RoboForex-Demo",
		"RoboForex-Demo 2",
		"RoboForex-Pro",
		"RoboForex-Pro-2",
		"RoboForex-ECN",
		"RoboForex-ECN-Pro",
		"MetaQuotes-Demo", // как диагностический
	}

	const serverSideWaitSeconds = 480          // ждём на стороне сервера
	const localReadinessWait = 8 * time.Minute // и потом ждём сами

	var (
		account      *mt5.MT5Account
		connectedSrv string
	)

	// Всегда создаём новый экземпляр для первой попытки
	newAccount := func() *mt5.MT5Account {
		acc, err := mt5.NewMT5Account(login, password, "", uuid.Nil)
		if err != nil {
			log.Fatalf("NewMT5Account: %v", err)
		}
		return acc
	}
	account = newAccount()
	defer func() {
		if account != nil {
			_ = account.Disconnect(context.Background())
		}
	}()

	for idx, srv := range servers {
		if idx > 0 {
			// между попытками — чистый инстанс
			_ = account.Disconnect(context.Background())
			account = newAccount()
		}

		log.Printf("connect(wait on server %ds) to %s ...", serverSideWaitSeconds, srv)
		// Пытаемся коннектиться, даже если вернулся error — будем дальше сами ждать готовность
		if err := account.ConnectByServerName(rootCtx, srv, proxy, true, serverSideWaitSeconds); err != nil {
			log.Printf("connect returned error (ignore & keep waiting): %v", err)
		}

		time.Sleep(2 * time.Second) // дать терминалу «вдохнуть»

		readyCtx, cancelReady := context.WithTimeout(rootCtx, localReadinessWait)
		ready := waitReady(readyCtx, account, localReadinessWait)
		cancelReady()

		if ready {
			connectedSrv = srv
			break
		}
	}

	if connectedSrv == "" {
		log.Fatal("не удалось подключиться ни к одному серверу. Если вы за прокси/фаерволом — задайте MT5_PROXY (http:// или socks5://) или откройте egress из MRPC.")
	}

	// === Далее обычная работа ===
	selectedSymbol := cfg.DefaultSymbol
	for _, s := range []string{cfg.DefaultSymbol, "EURUSD.", "EURUSD.pro"} {
		if s == "" {
			continue
		}
		if err := account.EnsureSymbolVisible(rootCtx, s); err == nil {
			selectedSymbol = s
			log.Printf("symbol ready: %s", s)
			break
		} else {
			log.Printf("EnsureSymbolVisible(%s): %v", s, err)
		}
	}

	svc := mt5.NewMT5Service(account)

	// Минимальный sanity-check, без «опасных» операций
	svc.ShowHealthCheck(rootCtx)
	svc.ShowCheckConnect(rootCtx)
	svc.ShowAccountSummary(rootCtx)
	svc.ShowQuote(rootCtx, selectedSymbol)

	log.Println("✅ Done.")
}
