# ADR-004: DDD-lite (Domain-Driven Design, Lightweight)

**Status:** Accepted

**Context:** DevFlow's business domain includes complex workflows (task state transitions, sprint lifecycle, RBAC). The codebase needs to model these rules explicitly rather than scattering them across controllers and handlers.

**Decision:** Apply DDD principles selectively ("DDD-lite"): entities with business methods, value objects, repository interfaces, and service layer, but skip Event Sourcing, CQRS, and Bounded Context mapping at the infrastructure level.

**What we adopt from DDD:**
- **Entities:** Rich domain objects with behavior (e.g., `Task.ChangeStatus()`, `Sprint.Start()`, `Organization.AddMember()`).
- **Value Objects:** Immutable types for Email, TaskPosition, Slug, Color.
- **Repository Interfaces:** Defined in domain layer, implemented in infrastructure.
- **Service Layer:** Use-case orchestration in `application/` package.
- **Aggregates:** Task as aggregate root (Task + Comments + Tags).

**What we skip:**
- **Event Sourcing:** Adds significant complexity without clear benefit for this domain.
- **CQRS:** Read/write separation is unnecessary for the current scale; we use the same models for reads and writes.
- **Domain Events:** Implemented as simple Go callbacks rather than a full event bus.

**Rationale:** Full DDD is designed for systems with extreme business complexity (finance, insurance). DevFlow has moderate complexity; full DDD would add ceremony without proportional value. DDD-lite gives us clear domain boundaries, testable business logic, and infrastructure independence — the 20% of DDD that delivers 80% of the value.

**Consequences:** Domain logic is concentrated in the `domain/` package and is framework-agnostic. Adding Event Sourcing or CQRS later is possible by introducing new patterns alongside existing code.
