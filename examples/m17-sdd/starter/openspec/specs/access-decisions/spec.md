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
- **WHEN** carol approves alice's pending request
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
- **WHEN** approve and reject for one pending request race
- **THEN** exactly one decision SHALL win atomically; the opposite decision SHALL return 409.
