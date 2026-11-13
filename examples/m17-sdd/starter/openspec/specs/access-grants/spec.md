# access-grants

## Purpose
Define whether a documentation access grant is currently usable, including the relationship between a request decision and the resulting permission.

## Requirements

### Requirement: GR-1 Create grant from approval
The first approval SHALL create exactly one grant for the request owner and resource, with approved_at set from the server clock. Rejection SHALL NOT grant access.

#### Scenario: First grant
- **WHEN** carol first approves alice's request for handbook
- **THEN** one grant SHALL exist for alice and handbook.

#### Scenario: No duplicate grant
- **WHEN** carol repeats approved
- **THEN** the original grant and approved_at SHALL be retained.

### Requirement: GR-2 Check current access
GET /access SHALL allow a user to check their own access and carol to check any demo user. allowed SHALL be true if at least one grant exists for that user and resource.

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
- **WHEN** two requests are approved for alice and handbook
- **THEN** allowed SHALL be true without requiring all requests to be approved.

### Requirement: GR-3 New grants are indefinite
In the baseline, all newly approved requests SHALL create indefinite grants without expires_at.

#### Scenario: No deadline
- **WHEN** an approved baseline grant is checked at any later controlled time
- **THEN** the grant SHALL continue to permit access.
