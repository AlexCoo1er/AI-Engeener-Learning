# Временный доступ

## Why
Бессрочные новые разрешения превышают потребность команды; N3/D7–D8 требуют срока и совместимости.

## What Changes
Новые заявки требуют целое hours 1..24, срок от одобрения; старые разрешения бессрочны. Повторы не сдвигают время.

## Capabilities
### New Capabilities
None.
### Modified Capabilities
- `access-decisions`: preserve authorization and repeat guarantees.
- `access-grants`: update usable grant policy.
- `access-requests`: validate requested duration.

## Impact
access-requests, access-decisions, access-grants; POST /requests и вычисление доступа; старые pending hours=0 требуют новой заявки.
