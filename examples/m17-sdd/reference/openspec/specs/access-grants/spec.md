# access-grants

## Purpose
Define whether a documentation access grant is currently usable, including the relationship between a request decision and the resulting permission.

## Requirements

### Requirement: GR-1 Create grant from approval
The first approval SHALL create exactly one grant with approved_at from the server clock and expires_at equal to approved_at plus the request's hours. Rejection SHALL NOT grant access.

#### Scenario: First grant
- **WHEN** carol first approves alice's request for handbook with hours 1
- **THEN** one grant SHALL exist for alice and handbook.

#### Scenario: No duplicate grant
- **WHEN** carol repeats approved
- **THEN** the original grant and approved_at SHALL be retained.

#### Scenario: Starts at approval
- **WHEN** a request waits two hours before approval with hours 1
- **THEN** expires_at SHALL be one hour after approval, not one hour after creation.

### Requirement: GR-2 Check current access
GET /access SHALL allow a user to check their own access and carol to check any demo user. allowed SHALL be true if at least one matching grant is not revoked and either has no expires_at or has now earlier than expires_at. Both the no-revocation condition and the time condition SHALL hold for that same grant.

#### Scenario: Authorized lookup
- **WHEN** alice requests GET /access?user=alice&resource=handbook or carol checks alice
- **THEN** the server SHALL return 200 with user, resource and boolean allowed.

#### Scenario: Foreign lookup
- **WHEN** bob checks access for alice
- **THEN** the server SHALL return 403.

#### Scenario: No grant
- **WHEN** a permitted caller checks a valid user and resource without a usable grant
- **THEN** allowed SHALL be false.

#### Scenario: Invalid lookup
- **WHEN** a known caller supplies an unknown target user or unsupported resource
- **THEN** the server SHALL return 400.

#### Scenario: Any grant
- **WHEN** one matching grant is expired or revoked and another matching grant is unexpired and not revoked
- **THEN** allowed SHALL be true.

#### Scenario: Exact expiry
- **WHEN** the controlled server clock equals expires_at or is later
- **THEN** that grant SHALL NOT permit access.

#### Scenario: Before expiry
- **WHEN** the controlled server clock is earlier than expires_at and the grant is not revoked
- **THEN** that grant SHALL permit access.

#### Scenario: Revoked indefinite grant
- **WHEN** a legacy indefinite grant is revoked and no other matching grant is usable
- **THEN** allowed SHALL be false.

### Requirement: GR-4 Preserve legacy grants
Grants approved under the baseline policy SHALL remain indefinite when the process switches to the temporary access policy; new policy SHALL NOT invent historical expiry values.

#### Scenario: Policy upgrade
- **WHEN** an existing baseline grant is retained in memory while policy changes to reference
- **THEN** the grant SHALL remain indefinite and usable unless explicitly revoked.

#### Scenario: Pending migration
- **WHEN** a baseline pending request with hours 0 is decided after policy upgrade
- **THEN** the server SHALL return 400 for approval; rejection SHALL remain available.
