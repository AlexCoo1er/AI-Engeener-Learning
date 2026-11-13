# Заполненная матрица: эталон для сверки

Идентификаторы стабильны внутри capability; название теста указывает исполняемую проверку, а не заменяет запуск. Сценарии остаются в основных спеках после архива. Это пример трассировки продукта; OpenSpec не требует именно такую таблицу.

| Потребность / источник | Требование | Сценарий / риск | Go | Java |
|---|---|---|---|---|
| N1 / D1–D4 | access-requests RQ-1/RQ-2; access-decisions DC-1 | Известная identity, своё создание, чужое чтение, self approval | TestIdentityAndTransitions; TestHTTPContract | AccessTest.identityAndTransitions; HttpTest.httpContract |
| N1 / D2,D12 | access-requests RQ-1 | Whitelist ресурса, неверный вход, HTTP-коды | TestHTTPContract | HttpTest.httpContract |
| N2 / D5 | access-decisions DC-2 | Окончательность, тот же и противоположный повтор | TestIdentityAndTransitions | AccessTest.identityAndTransitions |
| N2 / D5 | access-decisions DC-2 | Одновременные противоположные решения: один 200, другой 409 | TestConcurrentConflictingDecisions; -race | AccessTest.concurrentConflictingDecisions |
| N3 / D7 | access-requests RQ-1 | Допустимые 1/24; 0/25/negative; нецелый JSON | TestExpiryLegacyAndReplay; TestHTTPContract | AccessTest.expiryLegacyAndReplay; HttpTest.httpContract |
| N3 / D7,D11 | access-grants GR-1/GR-2 | Отсчёт от approval, до/на/после expiry, без sleep | TestExpiryLegacyAndReplay; TestApprovalClockAndLegacyRevocation | AccessTest.expiryLegacyAndReplay; AccessTest.approvalClockAndLegacyRevocation |
| N2,N3 / D5,D11 | access-decisions DC-2 | Повтор approved не переносит срок | TestExpiryLegacyAndReplay | AccessTest.expiryLegacyAndReplay |
| N3 / D8 | access-grants GR-4 | Legacy approved бессрочен; legacy pending не одобряется с hours=0 | TestExpiryLegacyAndReplay | AccessTest.expiryLegacyAndReplay |
| N4 / D9 | access-decisions DC-3 | Revoke, запреты, blank reason, повтор причины, другой reason409 | TestRevocationAndIndependentGrants | AccessTest.revokeAndIndependentGrants |
| N2,N4 / D9,D11 | access-decisions DC-2/DC-3 | Revoke не меняет approved, повтор approve не оживляет | TestRevocationAndIndependentGrants | AccessTest.revokeAndIndependentGrants |
| N4 / D10 | access-grants GR-2 | Второй активный grant сохраняет доступ | TestRevocationAndIndependentGrants | AccessTest.revokeAndIndependentGrants |
| N4 / D8,D9 | access-grants GR-2/GR-4 | Можно отозвать expired и legacy indefinite | TestApprovalClockAndLegacyRevocation | AccessTest.approvalClockAndLegacyRevocation |

## Доказательства и границы

Структура документов: `npm run verify:specs` в `examples/m17-sdd` (OpenSpec 1.16.0). Проверка выполнена автором 5 ноября 2026: baseline/final по три capability, два архивных change с завершёнными tasks, обе активные дельты на соответствующей базе, конфликтующий черновик и отдельное rename/remove-упражнение проходят strict. Настоящие последовательные `archive` на временной копии дают полное совпадение итоговых требований **и сценариев** с reference. Структурный валидатор не оценивает бизнес-смысл.

Поведение: выполните `go test -race ./...`, `go vet ./...`, `go build ./...` либо `mvn -pl m17-sdd -am verify`; сохраните собственную дату, команду, версию инструмента и результат рядом с матрицей. Точные результаты текущей сборки публикует отчёт проверки курса; таблица выше не обещает, что ваш изменённый вариант уже прошёл тесты. HTTP-проверки используют настоящий локальный сервер на случайном loopback-порту.

Человеческое ревью: согласовать D1–D12, проверить полноту сохраняемых сценариев и разобрать конфликт extension vs revoked. Ни галочки архивных tasks, ни AI verify не являются таким согласованием.

Продление не реализовано и не принято в reference. Тесты не доказывают реальные IAM/SSO, сохранность после рестарта, внешнюю конкурентность между процессами, production-защиту заголовка или нагрузочную производительность. У сервиса нет этих обещаний. Перед передачей студент отдельно отмечает все сценарии, которые его собственные тесты ещё не покрывают.
