# access-requests

## Purpose
Define creation and visibility of documentation access requests for a fixed local demo identity adapter, including input validation and ownership.

## Requirements

### Requirement: RQ-1 Create own request
The server SHALL accept POST /requests with JSON resource and integer hours from 1 through 24 inclusive; invalid input SHALL return 400 without creating a request.

#### Scenario: Known identity
- **WHEN** a caller supplies X-User alice, bob or carol
- **THEN** the server SHALL resolve the fixed demo identity; missing or unknown X-User SHALL return 401.

#### Scenario: Own request
- **WHEN** alice submits POST /requests with resource handbook and hours 1
- **THEN** the server SHALL create a pending request owned by alice and return 201 with a unique id; runbook SHALL also be accepted.

#### Scenario: Resource validation
- **WHEN** an authenticated caller submits hours 1 and a resource other than handbook or runbook
- **THEN** the server SHALL return 400 without creating a request.

#### Scenario: Duration endpoints
- **WHEN** alice submits requests for handbook with hours 1 and 24
- **THEN** both requests SHALL be accepted with 201 and preserve their requested hours.

#### Scenario: Duration validation
- **WHEN** hours is absent, 0, negative, 25, fractional, null, boolean or a string
- **THEN** the server SHALL return 400 and SHALL NOT create a request.

### Requirement: RQ-2 Read visible request
GET /requests/{id} SHALL return 200 with the request for its owner or carol, 403 for another known employee, and 404 for an absent id.

#### Scenario: Owner reads
- **WHEN** alice reads an existing request owned by alice
- **THEN** the server SHALL return its id, user, resource, hours and status.

#### Scenario: Approver reads
- **WHEN** carol reads an existing request owned by alice or carol
- **THEN** the server SHALL return 200.

#### Scenario: Foreign read
- **WHEN** bob reads an existing request owned by alice
- **THEN** the server SHALL return 403.

#### Scenario: Absent request
- **WHEN** a known caller reads an id that does not exist
- **THEN** the server SHALL return 404.
