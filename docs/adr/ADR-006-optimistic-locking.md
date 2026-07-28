# ADR-006: Optimistic Locking for Concurrent Editing

**Status:** Accepted

**Context:** Multiple users may simultaneously edit the same task or move it on the Kanban board. The system must prevent lost updates while maintaining performance and user experience.

**Decision:** Use optimistic locking with a `version_id` integer column on the `tasks` table.

**How it works:**
1. Client reads a task and receives `version_id` (e.g., 5).
2. Client modifies the task and sends `version_id: 5` with the update.
3. Server executes: `UPDATE tasks SET ..., version_id = version_id + 1 WHERE id = $1 AND version_id = 5`.
4. If `RowsAffected == 0`, the task was modified by another user — return HTTP 409 Conflict.
5. Client receives the conflict, refreshes the task data, and retries the operation.

**Rationale:** Optimistic locking works well for this domain because:
- **Low contention:** Task conflicts are rare in practice (different users typically work on different tasks).
- **No database locks:** Avoids `SELECT ... FOR UPDATE` which blocks concurrent reads.
- **Client-friendly:** The frontend can implement automatic retry with the refreshed data.
- **Simple to implement:** Single integer column, no additional infrastructure.

Alternatives considered:
- **Pessimistic locking:** `SELECT ... FOR UPDATE` would block concurrent edits, reducing throughput. Overkill for scenarios where conflicts are rare.
- **Last-write-wins:** Simplest approach but silently loses data — unacceptable for a professional tool.
- **CRDT (Conflict-free Replicated Data Types):** Technically superior for real-time collaboration but adds enormous complexity (merge logic, vector clocks). Justified for Google Docs or Figma, not for task management.

**Consequences:** The frontend must handle 409 responses gracefully (show "updated by another user" toast, auto-refresh). API consumers must always pass `version_id` with updates. The `tasks` table gains an additional index on `version_id` for the update WHERE clause.
