# ADR-003: Cursor-based Pagination

**Status:** Accepted

**Context:** The API requires pagination for task lists, activity logs, and notification feeds. The solution must handle large datasets efficiently and produce stable results when new items are inserted or deleted.

**Decision:** Use cursor-based pagination with opaque base64-encoded cursors containing the last item's ID and sort value.

**Rationale:** Cursor-based pagination addresses fundamental issues with offset-based pagination:
- **Stability:** Adding or removing items does not shift pages. Offset-based pagination causes items to appear on multiple pages or be skipped entirely when data changes between requests.
- **Performance:** No `COUNT(*)` or `OFFSET` operations required. Cursor queries use indexed `WHERE` clauses: `WHERE (created_at, id) < ($1, $2) ORDER BY created_at DESC LIMIT $3`.
- **Consistency:** Industry standard — used by Twitter, Stripe, GitHub GraphQL API.

Alternatives considered:
- **Offset-based:** Simpler to implement but breaks when data changes (items shift between pages). Performance degrades on large offsets (`OFFSET 100000` still scans 100000 rows).
- **Keyset pagination:** Similar to cursor but uses visible column values. Exposes internal sort values to clients.

**Consequences:** Clients cannot jump to arbitrary pages. Frontend must use "Load More" or infinite scroll instead of page numbers. Cursor format must remain stable across API versions (base64 + JSON provides flexibility for future schema changes).
