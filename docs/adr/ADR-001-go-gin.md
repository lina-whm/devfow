# ADR-001: Go with Gin Framework

**Status:** Accepted

**Context:** DevFlow requires a backend framework that balances performance, developer experience, and ecosystem maturity. The framework must handle RESTful routing, middleware chaining, request binding, and validation.

**Decision:** Use Go 1.24+ with the Gin web framework.

**Rationale:** Gin is the most widely adopted Go HTTP framework with the largest ecosystem of middleware and community support. Benchmarks show it outperforms the standard library's `net/http` mux by 40-60% in request throughput. Alternatives considered:
- **Fiber:** Faster in synthetic benchmarks, but API is incompatible with `net/http` (uses `fasthttp`), which breaks compatibility with many Go middleware libraries. Fiber's `context` API diverges significantly from standard Go patterns.
- **Chi:** More idiomatic Go, excellent middleware support via `http.Handler`, but ~2x slower than Gin in routing benchmarks.
- **net/http (standard):** Zero dependencies, but requires manual routing, lacks built-in request binding, validation, and error handling.

Gin provides the best trade-off: performance approaching Fiber, compatibility with `net/http`, built-in validation via `binding` tags, and a mature middleware ecosystem (CORS, rate limiting, logging).

**Consequences:** Development team must learn Gin's Context model and binding conventions. The framework ties us slightly to its API, but the Service and Repository layers remain framework-agnostic, allowing future replacement if needed.
