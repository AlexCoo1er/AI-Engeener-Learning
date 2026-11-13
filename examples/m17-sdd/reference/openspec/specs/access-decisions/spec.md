# access-decisions

## Purpose
Define authorization, atomic state transitions and repeat behavior for access decisions without permitting self approval or conflicting final decisions.

## Requirements

### Requirement: DC-1 Authorize decision
POST /requests/{id}/decision SHALL accept approved or rejected only; only carol SHALL decide another user's request. A known actor accessing an absent request SHALL receive 404.

#### Scenario: Employee forbidden
- **WHEN** alice or bob attempts to decide an existing request
- **THEN** the server SHALL return 403 without changing its state.

#### Scenario: No self approval
- **WHEN** carol attempts to approve or reject her own existing request
- **THEN** the server SHALL return 403 without changing its state.

#### Scenario: Decision validation
- **WHEN** carol submits a decision other than approved or rejected for another user's request
- **THEN** the server SHALL return 400 without changing its state.

#### Scenario: Missing identity
- **WHEN** a decision caller has an unknown or missing X-User
- **THEN** the server SHALL return 401.

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

#### Scenario: No revival by replay
- **WHEN** carol repeats approved after the request's grant was revoked
- **THEN** the server SHALL return 200 and SHALL preserve revoked_at and revoke_reason; the grant SHALL remain revoked.

### Requirement: DC-3 Revoke approved grant
POST /requests/{id}/revoke SHALL allow only carol to revoke another user's approved grant with a nonempty trimmed reason. Revocation SHALL retain request status approved and the original decision timestamps.

#### Scenario: Successful revocation
- **WHEN** carol revokes alice's approved grant with reason left team
- **THEN** the server SHALL return 200 and record revoked_at and trimmed revoke_reason; that grant SHALL no longer permit access.

#### Scenario: Repeated reason
- **WHEN** carol repeats revoke with the same trimmed reason
- **THEN** the server SHALL return 200 and preserve the first revoked_at.

#### Scenario: Different reason
- **WHEN** carol repeats revoke with a different reason
- **THEN** the server SHALL return 409 and preserve the first reason and timestamp.

#### Scenario: Reason validation
- **WHEN** carol revokes another user's approved request with an absent or blank reason
- **THEN** the server SHALL return 400 without changing the grant.

#### Scenario: Invalid grant state
- **WHEN** carol revokes another user's pending or rejected request
- **THEN** the server SHALL return 409.

#### Scenario: Self or employee revoke
- **WHEN** an employee revokes any grant or carol revokes her own grant
- **THEN** the server SHALL return 403.

#### Scenario: Absent grant
- **WHEN** a known caller revokes an absent request id
- **THEN** the server SHALL return 404.

#### Scenario: Expired grant
- **WHEN** carol revokes another user's expired approved grant with a valid reason
- **THEN** the server SHALL return 200 and record revocation without changing the original expiry.
