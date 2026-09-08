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
```

Install the package:

```bash
go get github.com/MetaRPC/GoMT5
```

---

## Step 2: Write Your Trading Code

Create your main application file and paste the following snippet:

```
import (
    "context"
    "fmt"
    "time"
    mt "github.com/MetaRPC/GoMT5"
)

client, err := mt.NewMT5Account(user, password, grpcServer)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err = client.ConnectByServerName(ctx, serverName, "EURUSD")
summary, err := client.AccountSummary(ctx)
fmt.Printf("Balance: %.2f, Equity: %.2f\n", summary.Balance, summary.Equity)
```

---

## Step 3: Run the Program

Run your application:

```bash
# Verify connection output
# Balance: 10000.00, Equity: 10000.00
```

---

## 🚀 Next Steps

Congratulations! You have successfully established a direct gRPC connection to MetaTrader 5. Next:
- Explore **[gRPC Streaming](GRPC_STREAM_MANAGEMENT.md)** to listen to live ticks.
- Check the **[API Reference](../API_Reference/MT5Account.md)** for all 40+ available terminal methods.
- Learn about high-level risk management and auto-normalization in **[MT5Sugar](../API_Reference/MT5Sugar.md)**.
