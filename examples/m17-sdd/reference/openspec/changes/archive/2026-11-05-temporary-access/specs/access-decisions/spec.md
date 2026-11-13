## MODIFIED Requirements

### Requirement: DC-2 Atomic final decision
The server SHALL transition pending to approved or rejected atomically; final decisions SHALL be immutable and identical repeats SHALL be idempotent.

#### Scenario: Approve pending
- **WHEN** carol approves alice's pending request with valid hours under the active policy
- **THEN** the request SHALL become approved and return 200.

#### Scenario: Reject pending
- **WHEN** carol rejects alice's pending request
- **THEN** the request SHALL become rejected and return 200; no grant SHALL be created.

#### Scenario: Same decision
- **WHEN** carol repeats the identical final decision
- **THEN** the server SHALL return 200 with no second grant or changed approval timestamp.

#### Scenario: Opposite decision
- **WHEN** carol submits the opposite decision after a final decision
- **THEN** the server SHALL return 409 without changing request or grant.

#### Scenario: Concurrent decisions
- **WHEN** approve and reject for one pending request with valid hours under the active policy race
- **THEN** exactly one decision SHALL win atomically; the opposite decision SHALL return 409.

#### Scenario: No extension by replay
- **WHEN** carol repeats approved before or after expiry
- **THEN** the original approved_at and expires_at SHALL remain unchanged; expired access SHALL remain expired.
