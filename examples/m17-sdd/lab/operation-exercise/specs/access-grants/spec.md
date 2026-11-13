## REMOVED Requirements
### Requirement: GR-3 New grants are indefinite
**Reason**: Remove a duplicate normative rule in this document-only exercise.
**Migration**: The exercise also makes the same guarantee explicit in GR-1 below. No runtime behavior changes; this exercise is never applied to the final reference snapshot.

## MODIFIED Requirements
### Requirement: GR-1 Create grant from approval
The first approval SHALL create exactly one grant for the request owner and resource, with approved_at set from the server clock, indefinite validity and no expires_at. Rejection SHALL NOT grant access.

#### Scenario: First grant
- **WHEN** carol first approves alice's request for handbook
- **THEN** one grant SHALL exist for alice and handbook.

#### Scenario: No duplicate grant
- **WHEN** carol repeats approved
- **THEN** the original grant and approved_at SHALL be retained.
