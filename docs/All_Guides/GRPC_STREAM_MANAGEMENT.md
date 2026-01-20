# gRPC Streaming and Subscription Management Guide

> **Complete guide** to working with real-time subscriptions in GoMT5

This document covers:

- ✅ **How to subscribe properly** to market data streams
- ✅ **How to stop subscriptions** without goroutine leaks
- ✅ **Common patterns** from simple to advanced
- ✅ **Architecture** and built-in safety mechanisms
- ✅ **Troubleshooting** and best practices

---

## Table of Contents

1. [Quick Start - Simple Subscription](#quick-start---simple-subscription)
2. [Available Streaming Methods](#available-streaming-methods)
3. [Complete Patterns (Simple → Advanced)](#complete-patterns-simple--advanced)
4. [Problem: Why Streams Need Management](#problem-why-streams-need-management)
5. [Solutions and Best Practices](#solutions-and-best-practices)
6. [Architecture and Safety](#architecture-notes)
7. [Troubleshooting](#troubleshooting-and-goroutine-leaks)

---

## Quick Start - Simple Subscription

### 1️⃣ Simplest Pattern (Auto-timeout) 

**Use context.WithTimeout** - automatically stops after time expires:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "your-project/examples/mt5"
    pb "github.com/MetaRPC/GoMT5/package"
    "github.com/google/uuid"
)

func main() {
    // Setup connection
    account, _ := mt5.NewMT5Account(591129415, "password", "mt5.mrpc.pro:443", uuid.New())
    account.ConnectEx(context.Background(), &pb.ConnectExRequest{
        Uuid:      591129415,
        Password:  "password",
        MtCluster: "FxPro-MT5 Demo",
    })

    // ✅ Stream for 10 seconds - stops automatically!
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD", "GBPUSD"}}
    dataChan, errChan := account.OnSymbolTick(ctx, req)

    tickCount := 0
    for {
        select {
        case data, ok := <-dataChan:
            if !ok {
                fmt.Println("Stream finished")
                return
            }
            tickCount++
            fmt.Printf("[%d] %s: Bid=%.5f, Ask=%.5f\n",
                tickCount, data.Tick.Symbol, data.Tick.Bid, data.Tick.Ask)

        case err, ok := <-errChan:
            if !ok {
                return
            }
            fmt.Printf("Error: %v\n", err)

        case <-ctx.Done():
            fmt.Println("Timeout reached - stopping")
            return
        }
    }
    // Done! Context cancelled, goroutines cleaned up ✅
}
```

**When to use:** Quick examples, testing, short monitoring sessions

---

### 2️⃣ Manual Control Pattern

**For full control** - you decide when to stop (e.g., manually via Ctrl+C):

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Start monitoring in background
go func() {
    req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
    dataChan, errChan := account.OnSymbolTick(ctx, req)

    for {
        select {
        case data, ok := <-dataChan:
            if !ok {
                return
            }
            fmt.Printf("Price: %.5f\n", data.Tick.Bid)

        case err, ok := <-errChan:
            if !ok {
                return
            }
            fmt.Printf("Error: %v\n", err)

        case <-ctx.Done():
            fmt.Println("Monitoring stopped")
            return
        }
    }
}()

// ... do other work ...

// Stop when needed
cancel()
time.Sleep(100 * time.Millisecond) // Give goroutine time to cleanup
```

**When to use:** Long-term monitoring, background services, production applications

---

## Available Streaming Methods

### MT5Account (Low-level Streams)

All streaming methods return **two channels**: `<-chan *Data` and `<-chan error`

| Method | Description | Returns |
|--------|-------------|---------|
| `OnSymbolTick()` | Real-time price ticks for symbols | `(<-chan *OnSymbolTickData, <-chan error)` |
| `OnTrade()` | Trade events (orders executed, modified, etc.) | `(<-chan *OnTradeData, <-chan error)` |
| `OnPositionProfit()` | Position P&L updates | `(<-chan *OnPositionProfitData, <-chan error)` |
| `OnPositionsAndPendingOrdersTickets()` | Order/position tickets | `(<-chan *OnPositionsAndPendingOrdersTicketsData, <-chan error)` |
| `OnTradeTransaction()` | Low-level trade transaction events | `(<-chan *OnTradeTransactionData, <-chan error)` |

**All require `context.Context` for stopping!**

**Documentation:** [Streaming Methods Overview](../MT5Account/7.%20Streaming_Methods/Streaming_Methods.Overview.md)

**Demonstration:** `examples/demos/lowlevel/03_streaming_methods.go` - examples of using all streaming methods from MT5Account

---

## Complete Patterns (Simple → Advanced)

### Pattern 1: Quick Example (5-10 seconds)

```go
// ✅ Auto-timeout - perfect for testing
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data, ok := <-dataChan:
        if !ok {
            return
        }
        fmt.Printf("Tick: %.5f\n", data.Tick.Bid)

    case err := <-errChan:
        fmt.Printf("Error: %v\n", err)

    case <-ctx.Done():
        fmt.Println("Timeout - stopping")
        return
    }
}
// Stops after 10 seconds automatically ✅
```

---

### Pattern 2: Event-limited Streaming

```go
// ✅ Stop after processing N events
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
dataChan, errChan := account.OnSymbolTick(ctx, req)

const maxEvents = 100
count := 0

for {
    select {
    case data, ok := <-dataChan:
        if !ok {
            return
        }
        count++
        fmt.Printf("[%d] %s: %.5f\n", count, data.Tick.Symbol, data.Tick.Bid)

        if count >= maxEvents {
            fmt.Println("Max events reached - stopping")
            cancel()
            return
        }

    case err := <-errChan:
        fmt.Printf("Error: %v\n", err)

    case <-ctx.Done():
        return
    }
}
```

---

### Pattern 3: Condition-based Stop

```go
// ✅ Stop when specific condition is met
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data, ok := <-dataChan:
        if !ok {
            return
        }
        fmt.Printf("Price: %.5f\n", data.Tick.Bid)

        // Stop if price crosses threshold
        if data.Tick.Bid > 1.10000 {
            fmt.Println("Target price reached!")
            cancel()
            return
        }

    case err := <-errChan:
        fmt.Printf("Error: %v\n", err)

    case <-ctx.Done():
        return
    }
}
```

---

### Pattern 4: Background Service with Manual Control

```go
type PriceMonitor struct {
    ctx    context.Context
    cancel context.CancelFunc
    done   chan struct{}
}

func NewPriceMonitor() *PriceMonitor {
    ctx, cancel := context.WithCancel(context.Background())
    return &PriceMonitor{
        ctx:    ctx,
        cancel: cancel,
        done:   make(chan struct{}),
    }
}

func (pm *PriceMonitor) Start(account *mt5.MT5Account, symbols []string) {
    go pm.monitorPrices(account, symbols)
}

func (pm *PriceMonitor) monitorPrices(account *mt5.MT5Account, symbols []string) {
    defer close(pm.done)

    req := &pb.OnSymbolTickRequest{SymbolNames: symbols}
    dataChan, errChan := account.OnSymbolTick(pm.ctx, req)

    for {
        select {
        case data, ok := <-dataChan:
            if !ok {
                fmt.Println("Data channel closed")
                return
            }
            pm.processTick(data)

        case err, ok := <-errChan:
            if !ok {
                fmt.Println("Error channel closed")
                return
            }
            fmt.Printf("Stream error: %v\n", err)

        case <-pm.ctx.Done():
            fmt.Println("Monitoring stopped gracefully")
            return
        }
    }
}

func (pm *PriceMonitor) processTick(data *pb.OnSymbolTickData) {
    fmt.Printf("%s: %.5f / %.5f\n",
        data.Tick.Symbol, data.Tick.Bid, data.Tick.Ask)
    // Your logic here...
}

func (pm *PriceMonitor) Stop() {
    pm.cancel()

    // Wait for goroutine to finish (with timeout)
    select {
    case <-pm.done:
        fmt.Println("Monitor stopped gracefully")
    case <-time.After(5 * time.Second):
        fmt.Println("Monitor stop timeout")
    }
}

// Usage:
monitor := NewPriceMonitor()
monitor.Start(account, []string{"EURUSD", "GBPUSD"})

fmt.Println("Press Enter to stop monitoring...")
bufio.NewReader(os.Stdin).ReadBytes('\n')

monitor.Stop()
```

---

### Pattern 5: Multiple Concurrent Streams

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// WaitGroup to track goroutines
var wg sync.WaitGroup

// Start tick stream
wg.Add(1)
go func() {
    defer wg.Done()

    req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
    dataChan, errChan := account.OnSymbolTick(ctx, req)

    for {
        select {
        case data, ok := <-dataChan:
            if !ok {
                return
            }
            fmt.Printf("Tick: %.5f\n", data.Tick.Bid)

        case err := <-errChan:
            fmt.Printf("Tick error: %v\n", err)

        case <-ctx.Done():
            return
        }
    }
}()

// Start trade stream
wg.Add(1)
go func() {
    defer wg.Done()

    req := &pb.OnTradeRequest{}
    dataChan, errChan := account.OnTrade(ctx, req)

    for {
        select {
        case data, ok := <-dataChan:
            if !ok {
                return
            }
            fmt.Printf("Trade: %v\n", data.Type)

        case err := <-errChan:
            fmt.Printf("Trade error: %v\n", err)

        case <-ctx.Done():
            return
        }
    }
}()

// Run for 30 seconds or until manual stop
time.Sleep(30 * time.Second)

// Stop all streams
cancel()

// Wait for all goroutines to finish
wg.Wait()
fmt.Println("All streams stopped")
```

---

### Pattern 6: Error Handling and Retry Logic

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

maxRetries := 3
retryDelay := 2 * time.Second

for attempt := 1; attempt <= maxRetries; attempt++ {
    fmt.Printf("Attempt %d/%d\n", attempt, maxRetries)

    req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
    dataChan, errChan := account.OnSymbolTick(ctx, req)

    streamBroken := false

    for !streamBroken {
        select {
        case data, ok := <-dataChan:
            if !ok {
                fmt.Println("Data channel closed")
                streamBroken = true
                break
            }
            fmt.Printf("Price: %.5f\n", data.Tick.Bid)
            attempt = 0 // Reset on successful data

        case err, ok := <-errChan:
            if !ok {
                streamBroken = true
                break
            }
            fmt.Printf("Stream error: %v\n", err)
            streamBroken = true

        case <-ctx.Done():
            fmt.Println("Context cancelled")
            return
        }
    }

    if attempt < maxRetries && ctx.Err() == nil {
        fmt.Printf("Reconnecting in %v...\n", retryDelay)
        time.Sleep(retryDelay)

        // Note: ExecuteStreamWithReconnect handles this automatically!
        // This pattern is for demonstration - usually not needed
    }
}
```

---

## Problem: Why Streams Need Management

When working with gRPC streaming in GoMT5, understanding the stream lifecycle is critical:

**Stream subscriptions** (`OnSymbolTick`, `OnTrade`, etc.) are **active goroutines with network connections** that continue running even after losing reference.

---

## Problem Explanation

### Current Implementation

```go
func (a *MT5Account) OnSymbolTick(ctx context.Context, req *pb.OnSymbolTickRequest)
    (<-chan *pb.OnSymbolTickData, <-chan error) {
    // Spawns goroutine...
    // ❌ Goroutine runs until context is cancelled
    // ❌ Need to read from BOTH channels to prevent goroutine leak
}
```

### What Happens Without Proper Cleanup

```go
// ❌ BAD: Goroutine continues running forever
dataChan, errChan := account.OnSymbolTick(context.Background(), req)

for data := range dataChan {
    fmt.Printf("%s: %.5f\n", data.Tick.Symbol, data.Tick.Bid)

    if someCondition {
        break  // ❌ PROBLEM: Goroutine still running in background!
               // ❌ PROBLEM: errChan not consumed - goroutine blocked!
    }
}
```

**Result:**

- Background goroutine continues consuming resources
- MT5 terminal continues sending updates
- **Goroutine leak** - blocked on writing to errChan
- Memory gradually accumulates
- Multiple abandoned streams = **serious memory leak**

---

## Solutions and Best Practices

## Solution 1: Always Use context.Context with Cancel ✅

### Pattern 1: context.WithCancel

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // ✅ CRITICAL: Ensures cleanup

req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data, ok := <-dataChan:
        if !ok {
            return
        }
        fmt.Printf("%s: %.5f\n", data.Tick.Symbol, data.Tick.Bid)

        if someCondition {
            cancel()  // ✅ CORRECT: Stops goroutine
            return
        }

    case err, ok := <-errChan:
        if !ok {
            return
        }
        fmt.Printf("Error: %v\n", err)

    case <-ctx.Done():
        return
    }
}
```

### Pattern 2: Timeout with context.WithTimeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

req := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data := <-dataChan:
        if data != nil {
            fmt.Printf("Price: %.5f\n", data.Tick.Bid)
        }

    case err := <-errChan:
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        }

    case <-ctx.Done():
        fmt.Println("Stream reached timeout or was cancelled")
        return
    }
}
```

---

## Solution 2: Always Read from BOTH Channels ✅

**CRITICAL:** Streaming methods return TWO channels. You MUST read from both!

### ❌ WRONG - Goroutine Leak

```go
// ❌ BAD: Only reading dataChan - errChan blocks goroutine!
dataChan, errChan := account.OnSymbolTick(ctx, req)

for data := range dataChan {
    // Process data...
    // ❌ PROBLEM: If error occurs, goroutine blocks on writing to errChan!
}
```

### ✅ CORRECT - Read Both Channels

```go
// ✅ GOOD: Read both channels using select
dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data, ok := <-dataChan:
        if !ok {
            return  // Channel closed
        }
        // Process data

    case err, ok := <-errChan:  // ✅ CRITICAL: Must read errors!
        if !ok {
            return  // Channel closed
        }
        // Process error

    case <-ctx.Done():
        return
    }
}
```

---

## Common Streaming Mistakes

### ❌ Mistake 1: No Context Cancellation

```go
// ❌ WRONG: No way to stop stream
dataChan, errChan := account.OnSymbolTick(context.Background(), req)

for {
    select {
    case data := <-dataChan:
        fmt.Println(data.Tick.Bid)
        // If you want to stop - YOU CAN'T!
    case err := <-errChan:
        fmt.Println(err)
    }
}
```

**Fix:**
```go
// ✅ CORRECT: Can cancel anytime
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

dataChan, errChan := account.OnSymbolTick(ctx, req)
// Now can call cancel() to stop
```

---

### ❌ Mistake 2: Exit Without Cancel

```go
// ❌ WRONG: Break doesn't stop goroutine
ctx, cancel := context.WithCancel(context.Background())
// defer cancel()  // ❌ MISSING!

dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data := <-dataChan:
        if data.Tick.Symbol == "EURUSD" {
            break  // ❌ Goroutine still running!
        }
    case err := <-errChan:
        // ...
    }
}
// ❌ cancel() never called - goroutine leak!
```

**Fix:**
```go
// ✅ CORRECT: defer cancel() ensures cleanup
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // ✅ ALWAYS executes

dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data := <-dataChan:
        if data.Tick.Symbol == "EURUSD" {
            return  // ✅ defer will call cancel()
        }
    case err := <-errChan:
        // ...
    case <-ctx.Done():
        return
    }
}
```

---

### ❌ Mistake 3: Reading Only Data Channel

```go
// ❌ WRONG: Not reading errChan - goroutine blocks!
dataChan, _ := account.OnSymbolTick(ctx, req)  // ❌ Ignoring errChan

for data := range dataChan {
    fmt.Println(data.Tick.Bid)
    // ❌ PROBLEM: If error occurs, goroutine blocks forever!
}
```

**Fix:**
```go
// ✅ CORRECT: Read both channels
dataChan, errChan := account.OnSymbolTick(ctx, req)

for {
    select {
    case data, ok := <-dataChan:
        if !ok {
            return
        }
        fmt.Println(data.Tick.Bid)

    case err, ok := <-errChan:  // ✅ MUST read errors!
        if !ok {
            return
        }
        fmt.Printf("Error: %v\n", err)

    case <-ctx.Done():
        return
    }
}
```

---

## Complete Example: Proper Stream Management

### Example 1: Simple Tick Monitoring with Timeout

```go
func MonitorTicks(account *mt5.MT5Account, symbols []string, duration time.Duration, maxTicks int) error {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    req := &pb.OnSymbolTickRequest{SymbolNames: symbols}
    dataChan, errChan := account.OnSymbolTick(ctx, req)

    tickCount := 0

    fmt.Printf("Starting tick monitoring (max %d ticks or %v)...\n", maxTicks, duration)

    for {
        select {
        case data, ok := <-dataChan:
            if !ok {
                fmt.Println("Data channel closed")
                return nil
            }

            tickCount++
            fmt.Printf("[%d] %s: Bid=%.5f, Ask=%.5f\n",
                tickCount, data.Tick.Symbol, data.Tick.Bid, data.Tick.Ask)

            if tickCount >= maxTicks {
                fmt.Println("Max ticks reached - stopping")
                return nil
            }

        case err, ok := <-errChan:
            if !ok {
                return nil
            }
            return fmt.Errorf("stream error: %w", err)

        case <-ctx.Done():
            fmt.Printf("Monitoring stopped (%d ticks processed)\n", tickCount)
            return nil
        }
    }
}

// Usage:
MonitorTicks(account, []string{"EURUSD", "GBPUSD"}, 30*time.Second, 100)
```

---

### Example 2: Production-ready Service

```go
type StreamService struct {
    account *mt5.MT5Account
    ctx     context.Context
    cancel  context.CancelFunc
    wg      sync.WaitGroup
}

func NewStreamService(account *mt5.MT5Account) *StreamService {
    ctx, cancel := context.WithCancel(context.Background())
    return &StreamService{
        account: account,
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (s *StreamService) StartTickMonitoring(symbols []string, handler func(*pb.OnSymbolTickData)) {
    s.wg.Add(1)

    go func() {
        defer s.wg.Done()

        req := &pb.OnSymbolTickRequest{SymbolNames: symbols}
        dataChan, errChan := s.account.OnSymbolTick(s.ctx, req)

        for {
            select {
            case data, ok := <-dataChan:
                if !ok {
                    return
                }
                handler(data)

            case err, ok := <-errChan:
                if !ok {
                    return
                }
                log.Printf("Tick stream error: %v", err)

            case <-s.ctx.Done():
                log.Println("Tick monitoring stopped")
                return
            }
        }
    }()
}

func (s *StreamService) StartTradeMonitoring(handler func(*pb.OnTradeData)) {
    s.wg.Add(1)

    go func() {
        defer s.wg.Done()

        req := &pb.OnTradeRequest{}
        dataChan, errChan := s.account.OnTrade(s.ctx, req)

        for {
            select {
            case data, ok := <-dataChan:
                if !ok {
                    return
                }
                handler(data)

            case err, ok := <-errChan:
                if !ok {
                    return
                }
                log.Printf("Trade stream error: %v", err)

            case <-s.ctx.Done():
                log.Println("Trade monitoring stopped")
                return
            }
        }
    }()
}

func (s *StreamService) Stop() {
    s.cancel()

    // Wait for all goroutines to finish (with timeout)
    done := make(chan struct{})
    go func() {
        s.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        log.Println("All streams stopped gracefully")
    case <-time.After(5 * time.Second):
        log.Println("Stream stop timeout - some goroutines may still be running")
    }
}

// Usage:
service := NewStreamService(account)

service.StartTickMonitoring([]string{"EURUSD", "GBPUSD"}, func(data *pb.OnSymbolTickData) {
    fmt.Printf("%s: %.5f\n", data.Tick.Symbol, data.Tick.Bid)
})

service.StartTradeMonitoring(func(data *pb.OnTradeData) {
    fmt.Printf("Trade event: %v\n", data.Type)
})

// ... do work ...

service.Stop()  // ✅ Gracefully stops all streams
```

---

## Recommendations

### ✅ DO:

1. **Always use `context.WithCancel` or `context.WithTimeout`** with streaming methods
2. **Always use `defer cancel()`** to ensure cleanup
3. **Always read from BOTH channels** (dataChan and errChan) using `select`
4. **Check `ok` value** when reading from channels (`data, ok := <-dataChan`)
5. **Use `sync.WaitGroup`** to track goroutines in production code
6. **Add timeout** in Stop() methods to prevent hanging

### ❌ DON'T:

1. **Never start streaming without context cancellation**
2. **Never ignore error channel** - causes goroutine leaks
3. **Never forget `defer cancel()`**
4. **Never use `context.Background()` without wrapping**
5. **Never assume channels close automatically**
6. **Never block main goroutine waiting for streams**

---

## Troubleshooting and Goroutine Leaks

### Checking for Goroutine Leaks

```go
import "runtime"

// Before streaming
before := runtime.NumGoroutine()
fmt.Printf("Goroutines before: %d\n", before)

// Your streaming code here...
{
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dataChan, errChan := account.OnSymbolTick(ctx, req)

    // ... process stream ...
}

// After streaming (wait a bit for cleanup)
time.Sleep(100 * time.Millisecond)
after := runtime.NumGoroutine()
fmt.Printf("Goroutines after: %d\n", after)
fmt.Printf("Leaked goroutines: %d\n", after-before)
// Should be ~0 if properly cleaned up
```

### Debugging Goroutine Stacks

```go
import (
    "os"
    "runtime/pprof"
)

// Dump goroutine stacks to file
f, _ := os.Create("goroutines.txt")
defer f.Close()
pprof.Lookup("goroutine").WriteTo(f, 1)
```

### Common Leak Patterns

1. **Not reading error channel:**
   ```go
   // ❌ Leak: goroutine blocks on writing to errChan
   dataChan, _ := account.OnSymbolTick(ctx, req)
   ```

2. **Context never cancelled:**
   ```go
   // ❌ Leak: goroutine runs forever
   dataChan, errChan := account.OnSymbolTick(context.Background(), req)
   // Forgot defer cancel()
   ```

3. **Exit without cleanup:**
   ```go
   // ❌ Leak: goroutine still running after break
   for data := range dataChan {
       if condition {
           break  // Missing cancel()
       }
   }
   ```

---

## Architecture Notes

### How Go Streaming Works in MT5Account

**Data flow:**

```
User code (select on channels)
    ↓
Two channels: dataChan, errChan
    ↓
ExecuteStreamWithReconnect (background goroutine)
    ↓
gRPC ClientStream
    ↓
Network (to MT5 terminal)
```

**Context propagation:**
```
User code → context.Context → gRPC call → Network
```

**When you cancel context:**

1. `ctx.Done()` channel closes
2. Background goroutine sees `<-ctx.Done()`
3. **Goroutine closes both dataChan and errChan** ✅
4. gRPC stream is cancelled
5. Network connection closes
6. Resources freed ✅

### MT5Account Cleanup Mechanism

The `ExecuteStreamWithReconnect` method ensures proper cleanup:

```go
// In MT5Account.go (lines 362-380)
func ExecuteStreamWithReconnect[TRequest, TReply, TData any](
    ctx context.Context,
    account *MT5Account,
    // ...
) (<-chan TData, <-chan error) {
    dataChan := make(chan TData)      // Unbuffered channel for data
    errChan := make(chan error, 1)    // Buffer 1 for errors

    go func() {
        defer close(dataChan)  // ✅ ALWAYS closes channels
        defer close(errChan)   // ✅ ALWAYS closes channels

        for {
            // ... streaming logic with automatic reconnection ...
            // On ctx.Done() goroutine exits
        }
    }()

    return dataChan, errChan
}
```

**This means:**

- ✅ Channels always close on context cancellation
- ✅ Goroutine always exits gracefully
- ✅ Automatic reconnection on network errors (with exponential backoff)
- ✅ **You still MUST read from both channels** to prevent blocking
- ✅ **You still MUST cancel context** for clean shutdown

---


**Remember:**

- gRPC streams are **active goroutines with network connections**
- Go channels **must be consumed** or goroutines block
- `defer close()` in ExecuteStreamWithReconnect **prevents leaks**
- But `context.Cancel()` ensures **graceful shutdown**
- **ALWAYS read from both dataChan and errChan!**

---

## See Also

- **[Streaming Methods Overview](../MT5Account/7.%20Streaming_Methods/Streaming_Methods.Overview.md)** - Complete streaming methods documentation
- **[USERCODE_SANDBOX_GUIDE.md](USERCODE_SANDBOX_GUIDE.md)** - Quick start guide
- **[Go Concurrency Patterns](https://go.dev/blog/pipelines)** - Official Go blog

---

**Remember:** Streams are powerful tools for real-time market data, but they require proper lifecycle management. Master context cancellation, always read from both channels, and your streaming code will be robust and leak-free. Happy streaming!


