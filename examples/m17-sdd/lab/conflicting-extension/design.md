## Context
Draft based on temporary-access only, not on revoke-access.
## Decisions
No implementation is approved. Keep revoked_at immutable and ask the owner about expired grants, limits and repeat keys.
## Risks / Trade-offs
An approved request can have a revoked grant. A state check on approved alone can revive a deliberately revoked permission.
