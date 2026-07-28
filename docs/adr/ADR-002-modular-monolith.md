# ADR-002: Modular Monolith Architecture

**Status:** Accepted

**Context:** DevFlow must handle multiple bounded contexts (users, organizations, projects, tasks, sprints, notifications) while maintaining development velocity for a small team. The architecture must not preclude future migration to microservices.

**Decision:** Implement a Modular Monolith — a single deployable unit with logically separated modules using Go packages and interface-based communication.

**Rationale:** A Modular Monolith provides strong module boundaries without the operational overhead of microservices:
- **Transactionality:** ACID transactions across bounded contexts (e.g., create task + update board position) in a single database transaction, avoiding distributed transaction complexity (Saga, Outbox).
- **Deployment simplicity:** Single binary, single Docker image, single health check.
- **Development speed:** No network communication between modules, easier debugging, faster iteration.
- **Future-proofing:** Each module has a clearly defined interface. Extracting a module to a microservice requires only implementing the same interface over gRPC/HTTP.

Alternatives considered:
- **Microservices:** Premature for a small team; adds networking, deployment, observability, and testing overhead without commensurate benefit.
- **Serverless (Lambda/Functions):** Poor fit for long-lived WebSocket connections (needed for real-time board updates) and stateful operations.

**Consequences:** Code discipline is required to prevent module boundary violations. Go package import cycles must be avoided. When the team grows beyond 5-8 developers or traffic patterns require independent scaling, modules can be extracted to microservices.
