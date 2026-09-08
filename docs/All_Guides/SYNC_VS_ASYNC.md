# Synchronous vs Asynchronous Methods - When to Use What (Go)

> Understanding execution models in GoMT5: non-blocking streaming vs synchronous execution.

---

## 🎯 Quick Comparison

Go uses goroutines, channels, and `context.Context` for concurrency. Streaming methods stream into buffered channels, while requests use timeouts.

| Aspect | Asynchronous Pattern | Synchronous Call |
|--------|----------------------|------------------|
| **Thread Blocking** | ❌ Non-blocking (high concurrency) | ✅ Blocks current thread |
| **Throughput** | ✅ Handles thousands of events/sec | ❌ Limited by thread pool |
| **Real-Time Data** | ✅ Perfect for tick & trade streams | ⚠️ Inefficient for streams |
| **Simplicity** | Requires async runtime awareness | Simple, linear execution |
| **Recommended for** | Automated bots, microservices, GUIs | CLI scripts, notebooks, one-offs |

---

## 🚀 When to Use Asynchronous Methods (Recommended)

### 1. Market Data Streaming
Market ticks arrive at microsecond intervals during peak sessions. Asynchronous handlers ensure zero tick drops without freezing your execution thread:

```
client, err := mt.NewMT5Account(user, password, grpcServer)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err = client.ConnectByServerName(ctx, serverName, "EURUSD")
summary, err := client.AccountSummary(ctx)
fmt.Printf("Balance: %.2f, Equity: %.2f\n", summary.Balance, summary.Equity)
```

### 2. High-Frequency Order Execution
When operating across multiple currency pairs simultaneously, asynchronous dispatch allows your bot to send orders concurrently rather than sequentially.

---

## 💡 Summary

Always prefer asynchronous paradigms for production bots, multi-symbol trading, and background services. Use synchronous wrappers for quick setup scripts, testing, or exploratory analysis.
