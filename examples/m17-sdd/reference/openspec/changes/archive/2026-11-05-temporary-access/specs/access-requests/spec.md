## MODIFIED Requirements

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
