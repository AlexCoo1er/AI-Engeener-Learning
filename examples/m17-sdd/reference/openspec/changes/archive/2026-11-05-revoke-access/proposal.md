# Отзыв доступа

## Why
Нужно прекратить конкретное разрешение без ложного изменения исторического решения; N4/D9–D11.

## What Changes
Добавляем отзыв с причиной и защиту от восстановления повтором approve. Проверяем доступ по любому активному grant.

## Capabilities
### New Capabilities
None.
### Modified Capabilities
- `access-decisions`: preserve authorization and repeat guarantees.
- `access-grants`: update usable grant policy.

## Impact
access-decisions/access-grants; новый POST /requests/{id}/revoke; extension остаётся упражнением.
