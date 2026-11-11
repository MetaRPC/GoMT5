---
description: "MetaRPC GoMT5 — type-safe Go SDK for MetaTrader 5 trading automation over gRPC. Connect, stream live ticks and positions, and execute orders from any Go application."
hide:
  - navigation
  - toc
---

# GoMT5 SDK Documentation

> ℹ️ **Account ID / Session ID**: when connecting via `Connect` / `ConnectEx`, MetaRPC automatically generates a terminal session GUID and returns it in `terminalInstanceGuid`. You do not need to call `GetId` or supply an `id` header before connecting — subsequent calls (subscriptions, order requests) use this session ID automatically.

<div class="hx-wrap" markdown="0">

<div class="hx-hero">
  <p class="hx-sub">Type-safe Go SDK for MetaTrader 5 trading automation over gRPC — connect, stream live ticks and positions, and execute orders from any Go application.</p>
  <div class="hx-stats">
    <div class="hx-stat"><span class="hx-statV">3</span><span class="hx-statL">API Layers</span></div>
    <div class="hx-stat"><span class="hx-statV">gRPC</span><span class="hx-statL">protobuf + TLS</span></div>
    <div class="hx-stat"><span class="hx-statV">MT5</span><span class="hx-statL">Cloud Terminal</span></div>
    <div class="hx-stat"><span class="hx-statV">Go</span><span class="hx-statL">goroutines + channels</span></div>
  </div>
</div>

<div class="hx-cards hx-cards-4">

  <a href="All_Guides/Your_First_Project/" class="hx-gc hx-orange">
    <div class="hx-gcTop"></div>
    <div class="hx-gcBody">
      <div class="hx-gcTag">Get Started</div>
      <div class="hx-gcTitle">Quick Start</div>
      <p class="hx-gcDesc">Your first project from scratch in 10 minutes. Connect, read your balance, place a trade.</p>
      <div class="hx-gcCount">10 min · hands-on</div>
    </div>
  </a>

  <a href="All_Guides/GETTING_STARTED/" class="hx-gc hx-blue">
    <div class="hx-gcTop"></div>
    <div class="hx-gcBody">
      <div class="hx-gcTag">Learn</div>
      <div class="hx-gcTitle">Getting Started</div>
      <p class="hx-gcDesc">New here? Start with setup, the API key, and an overview of how the SDK is organized.</p>
      <div class="hx-gcCount">setup · overview</div>
    </div>
  </a>

  <a href="All_Guides/PROJECT_MAP/" class="hx-gc hx-purple">
    <div class="hx-gcTop"></div>
    <div class="hx-gcBody">
      <div class="hx-gcTag">Architecture</div>
      <div class="hx-gcTitle">Project Map</div>
      <p class="hx-gcDesc">How the low-, mid- and high-level layers fit together, and where everything lives.</p>
      <div class="hx-gcCount">structure · layers</div>
    </div>
  </a>

  <a href="All_Guides/GLOSSARY/" class="hx-gc hx-pink">
    <div class="hx-gcTop"></div>
    <div class="hx-gcBody">
      <div class="hx-gcTag">Reference</div>
      <div class="hx-gcTitle">Glossary</div>
      <p class="hx-gcDesc">MT5 terms, return codes, enums and concepts — a reference to keep open while you build.</p>
      <div class="hx-gcCount">terms · concepts</div>
    </div>
  </a>

</div>

<div class="hx-cards">

  <a href="MT5Account/MT5Account.Master.Overview/" class="hx-gc hx-teal">
    <div class="hx-gcTop"></div>
    <div class="hx-gcBody">
      <div class="hx-gcTag">Reference · Low-level</div>
      <div class="hx-gcTitle">MT5Account</div>
      <p class="hx-gcDesc">Direct gRPC / protobuf layer. Maximum control over every request and typed response.</p>
      <div class="hx-gcCount">raw gRPC · full control</div>
    </div>
  </a>

  <a href="MT5Service/MT5Service.Overview/" class="hx-gc hx-rose">
    <div class="hx-gcTop"></div>
    <div class="hx-gcBody">
      <div class="hx-gcTag">Reference · Mid-level</div>
      <div class="hx-gcTitle">MT5Service</div>
      <p class="hx-gcDesc">Friendly facade over MT5Account — clean arguments and unwrapped return values.</p>
      <div class="hx-gcCount">friendly args · clean values</div>
    </div>
  </a>

  <a href="MT5Sugar/MT5Sugar.API_Overview/" class="hx-gc hx-green">
    <div class="hx-gcTop"></div>
    <div class="hx-gcBody">
      <div class="hx-gcTag">Reference · High-level</div>
      <div class="hx-gcTitle">MT5Sugar</div>
      <p class="hx-gcDesc">One-liner operations, risk-based position sizing and smart defaults for fast strategies.</p>
      <div class="hx-gcCount">one-liners · risk sizing</div>
    </div>
  </a>

</div>

</div>
