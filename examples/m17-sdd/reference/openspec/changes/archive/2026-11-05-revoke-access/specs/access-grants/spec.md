## MODIFIED Requirements

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
