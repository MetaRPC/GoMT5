# Your First Project in 10 Minutes (Go)

> **Hands-on Quick Start** - Create a working trading project with MetaTrader 5 and GoMT5 from scratch.

---

## Step 0: Obtain Your API Key

To connect to MetaRPC endpoints (`mt5.mrpc.pro:443`), obtain your API key:
1. Register for free at [https://mrpc.pro/signup](https://mrpc.pro/signup).
2. Generate your API token in your dashboard at [https://mrpc.pro/my](https://mrpc.pro/my).
3. Set your token in your environment or connection config.

---

## Step 1: Create Your Project

Create a new directory for your trading bot:

```bash
mkdir my_gomt5_bot
cd my_gomt5_bot
go mod init my_gomt5_bot
```

Install the package:

```bash
go get github.com/MetaRPC/GoMT5/package
```

---

## Step 2: Write Your Trading Code

Create `main.go`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/MetaRPC/GoMT5/package"
	mt5 "github.com/MetaRPC/GoMT5/package/Helpers"
	"github.com/google/uuid"
)

func main() {
	login, _ := strconv.ParseUint(os.Getenv("MT5_LOGIN"), 10, 64)
	password := os.Getenv("MT5_PASSWORD")
	server := os.Getenv("MT5_SERVER") // the broker server name as shown in MetaTrader, e.g. "MetaQuotes-Demo"
	grpcServer := os.Getenv("MRPC_GRPC_HOST")
	if grpcServer == "" {
		grpcServer = "mt5.mrpc.pro:443"
	}

	// The API key is read from the MRPC_API_KEY environment variable; ConnectEx sets the terminal id.
	account, err := mt5.NewMT5Account(login, password, "mt5.mrpc.pro:443", uuid.Nil)
	if err != nil {
		log.Fatal(err)
	}
	defer account.DisconnectAndClose()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	timeout := uint32(120)
	if _, err := account.ConnectEx(ctx, &pb.ConnectExRequest{
		User: login, Password: password, MtClusterName: server, TimeoutSeconds: &timeout,
	}); err != nil {
		log.Fatal(err)
	}

	summary, err := account.AccountSummary(ctx, &pb.AccountSummaryRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %.2f, Equity: %.2f %s\n", summary.AccountBalance, summary.AccountEquity, summary.AccountCurrency)
}
```

The program uses the SDK's `MT5Account` helper (`github.com/MetaRPC/GoMT5/package/Helpers`), which connects, keeps the terminal id and sends it
with your API key on every call.

---

## Step 3: Run the Program

```bash
export MRPC_API_KEY=your-api-key
export MT5_LOGIN=12345678
export MT5_PASSWORD=your-password
export MT5_SERVER=MetaQuotes-Demo
go run .
```

```
Balance: 10000.00, Equity: 10000.00 USD
```

---

## 🚀 Next Steps

Congratulations! You have successfully established a direct gRPC connection to MetaTrader 5. Next:
- Explore **[gRPC Streaming](GRPC_STREAM_MANAGEMENT.md)** to listen to live ticks.
- Check the **[API Reference](../API_Reference/MT5Account.md)** for all 40+ available terminal methods.
- Learn about high-level risk management and auto-normalization in **[MT5Sugar](../API_Reference/MT5Sugar.md)**.
