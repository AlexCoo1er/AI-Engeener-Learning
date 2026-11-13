# SDD / OpenSpec: Access Requests

Небольшая лабораторная после [теоретического модуля 17](../../ai-engineer-course/m17.html): сначала согласуйте контракт существующего сервиса, затем проведите два изменения. Выберите **Go или Java**, на практику отведите 6–7 часов. AI coding agent необязателен. Всё работает локально, без ключей, LLM, БД и Docker.

Это демонстрационный сервис проверки доступа к документации, не настоящее управление правами. Один процесс, память без долговечности, loopback. X-User — фиксированный тестовый адаптер: alice/bob сотрудники, carol согласующий. Произвольный заголовок роли прав не выдаёт. После рестарта данные исчезают.

## Что где находится

| Материал | Назначение |
|---|---|
| [Бриф](lab/brief.md) | Неоднозначная потребность; сначала задайте свои вопросы |
| [Ответы заказчика](lab/stakeholder-answers.md) | Источники D1–D12, открывайте после попытки |
| [Анализ](lab/analysis.md) | Глоссарий, роли, состояния, API и ограничения |
| [Starter](starter/openspec/specs) | Основные спеки только исходного baseline |
| [Reference](reference/openspec/specs) | Итог после двух изменений — сверка после самостоятельной попытки |
| [Два архива](reference/openspec/changes/archive) | Полные proposal/design/tasks/deltas; не копируйте сразу в рабочую основную спеку |
| [Продление и конфликт](lab/conflicting-extension/README.md) | Параллельная дельта со старой базой; продление не реализовано |
| [RENAMED / REMOVED](lab/operation-exercise) | Отдельное документное упражнение на копии starter |
| [Трассировка: шаблон](lab/traceability-template.md) / [решение](lab/traceability-reference.md) | Потребность → сценарий → тест → доказательство |

Термин «мастер-спека» обозначает набор трёх `openspec/specs/<capability>/spec.md`, а не один обязательный файл. Starter и reference намеренно разделены: эталон не переписывает отправную точку.

## Запустить работающий baseline

Go (из корня репозитория):

```sh
cd examples/m17-sdd
go test ./...
go run ./cmd/server --mode baseline --addr 127.0.0.1:8080
```

Java 21 (из корня репозитория):

```sh
cd examples/java
mvn -pl m17-sdd -am verify
java -cp m17-sdd/target/classes course.m17sdd.Server --mode baseline --port 8080
```

Серверы используют один порт, запускайте один. В другом терминале:

```sh
curl -i -H 'X-User: alice' -H 'Content-Type: application/json' \
  -d '{"resource":"handbook"}' http://127.0.0.1:8080/requests
# Возьмите id из JSON вместо REQUEST_ID:
curl -i -H 'X-User: carol' -H 'Content-Type: application/json' \
  -d '{"decision":"approved"}' http://127.0.0.1:8080/requests/REQUEST_ID/decision
curl -i -H 'X-User: alice' \
  'http://127.0.0.1:8080/access?user=alice&resource=handbook'
```

Создание даёт 201, одобрение и проверка — 200; allowed=true. Проверьте чужое чтение, самосогласование carol и повтор/противоположное решение. В baseline hours отсутствует или равен 0, все разрешения бессрочны, revoke ещё нет. JSON и ошибки описаны в [анализе](lab/analysis.md).

Переключение режима перезапуском начинает пустой процесс. Сохранение legacy-разрешений проверяется тестами: `Upgrade()` в Go / `upgrade()` в Java меняет политику в том же экземпляре; HTTP-endpoint миграции намеренно отсутствует.

Для эталона замените `--mode baseline` на `--mode reference`. Теперь создание требует `{"resource":"handbook","hours":1}`; доступ истечёт через час **после одобрения**. Отзыв:

```sh
curl -i -H 'X-User: carol' -H 'Content-Type: application/json' \
  -d '{"reason":"left team"}' http://127.0.0.1:8080/requests/REQUEST_ID/revoke
```

Реальная проверка истечения использует управляемые часы в тестах, не ожидание часа. Смена процесса не является миграцией: совместимость старого разрешения тестируется переключением политики в одном хранилище. Реальные IAM, SSO, почта, provisioning и продление вне объёма.

## Локальный закреплённый OpenSpec

Проверено 5 ноября 2026: **@fission-ai/openspec 1.16.0**, Node.js **>=22.12.0**. Manifest и lock закрепляют release. Команды ниже проверены по установленному пакету, его help/schema и реальному archive; инструкции moving main могут отличаться. [Релиз и исходники](https://github.com/Fission-AI/OpenSpec/tree/v1.16.0), [официальная документация CLI](https://github.com/Fission-AI/OpenSpec/blob/v1.16.0/docs/cli.md).

Из корня репозитория:

```sh
cd examples/m17-sdd
npm ci --ignore-scripts
npm run openspec -- --version
npm run verify:specs
```

`verify:specs` проверяет starter/final, материализует каждый архив как активный change на правильной базе, делает два настоящих archive в временной копии и сравнивает полные итоговые требования и сценарии с reference. Отдельно проверяет конфликтующий черновик на старой базе и RENAMED/REMOVED. `--archived` сам по себе проверяет завершённость tasks, поэтому не заменяет такую проверку дельт. Скрипт удаляет только собственную временную папку и не меняет раздаточные материалы. Python 3 нужен для сравнения.

Подготовьте собственную учебную копию. **Не запускайте init в корне курса.**

```sh
# Текущая папка examples/m17-sdd после npm ci
mkdir student-work
cp -R starter/openspec student-work/
OPENSPEC_BIN="$PWD/node_modules/.bin/openspec"
export OPENSPEC_TELEMETRY=0
export XDG_CONFIG_HOME="$PWD/student-work/.openspec-config"
export XDG_DATA_HOME="$PWD/student-work/.openspec-data"
cd student-work
"$OPENSPEC_BIN" init --tools codex --profile core --no-animation .
"$OPENSPEC_BIN" validate --all --strict --no-interactive
```

Вместо codex выберите поддерживаемый инструмент из `init --help`; для ручной работы — `--tools none`. Пути XDG изолируют настройки/data **этой копии**, профиль пользователя не меняется. В следующем терминале повторно задайте эти переменные; глобальную установку не выполняйте. В Codex CLI/IDE используйте созданные skills; в Desktop выберите skill в боковой панели. Перезапуск/обновление списка skills может потребоваться вашей интеграции.

В версии 1.16.0 нет пресета `expanded`. Для пошагового учебного пути добавьте workflows **в ту же изолированную конфигурацию**:

```sh
"$OPENSPEC_BIN" config set profile custom
"$OPENSPEC_BIN" config set workflows '["propose","explore","new","continue","ff","apply","update","verify","sync","archive"]'
"$OPENSPEC_BIN" init --tools codex --profile custom --no-animation .
```

Эта последовательность тоже проверена на временной копии. Нельзя опускать изоляцию XDG, если хотите сохранить настройки своего профиля. `init --profile core` создаёт шесть core-workflows: propose, explore, apply, update, sync, archive. new/continue/ff/verify требуют custom. [Команды релиза](https://github.com/Fission-AI/OpenSpec/blob/v1.16.0/docs/commands.md).

| Каноническое действие агента | Реальное имя skill Codex 1.16.0 | Когда |
|---|---|---|
| /opsx:explore | $openspec-explore | Исследовать до proposal или при вопросах |
| /opsx:propose | $openspec-propose | Сразу подготовить пакет для ревью |
| /opsx:update | $openspec-update-change | Пересмотреть договорённости и артефакты |
| /opsx:apply | $openspec-apply-change | Реализовать проверяемые tasks |
| /opsx:verify | $openspec-verify-change | Проверка агентом после тестов, custom |
| /opsx:sync | $openspec-sync-specs | Перенести согласованные дельты в main |
| /opsx:archive | $openspec-archive-change | Архивировать завершённый change |
| /opsx:new | $openspec-new-change | Пошаговый старт, custom |
| /opsx:continue | $openspec-continue-change | Следующий артефакт, custom |
| /opsx:ff | $openspec-ff-change | Все планирующие артефакты, custom |

Это действия AI-агента, их не вводят в shell. Они не заменяют тесты. Для ручного выполнения документы можно написать самостоятельно, пользуясь шаблонами. Детерминированные CLI-команды:

```sh
"$OPENSPEC_BIN" new change temporary-access --schema spec-driven
"$OPENSPEC_BIN" status --change temporary-access
"$OPENSPEC_BIN" instructions proposal --change temporary-access
# После заполнения proposal, specs, design, tasks:
"$OPENSPEC_BIN" show temporary-access --type change --json --deltas-only
"$OPENSPEC_BIN" validate temporary-access --type change --strict --no-interactive
# Только после реализации, тестов и приёмки:
"$OPENSPEC_BIN" archive temporary-access --yes
"$OPENSPEC_BIN" list --specs
"$OPENSPEC_BIN" validate --all --strict --no-interactive
```

CLI `archive` синхронизирует дельты в основные спеки и переносит change в архив. Это отдельное действие от Git merge и запуска программы. Если sync уже сделан агентом, сначала проверьте main/delta diff и используйте агентный archive с проверкой синхронизации: повторное бездумное применение ADDED может конфликтовать. Учебный детерминированный путь выше выполняет синхронизацию **один раз при archive**. Не используйте `--no-validate`/`--skip-specs`, чтобы спрятать ошибку.

## Четыре этапа

1. **Анализ и baseline (60–90 мин).** Прочитайте бриф, исследуйте ответы API и тесты, сформулируйте вопросы. Затем сравните с D1–D6 и собственными словами согласуйте baseline. Составьте словарь, роли/состояния и матрицу трассировки. Сравните свои основные спеки со starter: отсутствие функции в коде не всегда означает согласованный запрет.
2. **Временный доступ (2 часа).** Создайте temporary-access от baseline. Опишите proposal, полные дельты, design и tasks. Ответьте на D7–D8 до кода: 1/24, invalid, отсчёт от approval, точная граница `now >= expires_at`, старые grants и pending. MODIFIED — весь итоговый блок с сохранёнными сценариями, не одна новая строка. Strict validator 1.16.0 обнаруживает исчезнувшие названия сценариев; сохранённое название с ослабленным смыслом требует ревью человеком. Реализуйте, запустите тесты, согласуйте результат; `archive` обновит main. Сверьте первый [разобранный change](reference/openspec/changes/archive/2026-11-05-temporary-access).
3. **Отзыв и конфликт (2 часа).** Создайте revoke-access уже от обновлённой main. Добавьте обязательную причину, идемпотентный повтор, запреты и отсутствие revival через approve. Проверьте несколько grants, legacy и expired. Изучите [параллельное продление](lab/conflicting-extension/README.md), запишите семантический конфликт и вопросы владельцу; реализовывать extension не нужно. В отдельной копии starter примените [operation-exercise](lab/operation-exercise): RENAMED меняет заголовок, REMOVED удаляет дубликат с сохранением гарантии. После приёмки отзыва архивируйте второй change.
4. **Приёмка и передача (60 мин).** Сдайте актуальные три main specs, два архива, анализ конфликта, заполненную матрицу и реальные результаты проверок. Различайте доказательства CLI, тестов программы и человеческого ревью. Сравните с reference, объясните различия. Дополнительное задание: тем же способом опишите одно небольшое изменение своего RAG-проекта.

Минимальная сдача: сохранены запреты чужого чтения и self approval; идентичные решения идемпотентны, противоположные конфликтуют; 1/24 и exact expiry протестированы без sleep; старые разрешения не получили выдуманный срок; revoke не меняет решение на rejected, повторы не оживляют доступ. Все MODIFIED содержат прежние сохраняемые сценарии; у REMOVED есть причина и миграция. Структурный валидатор не доказывает полноту или правильность этих договорённостей.

## Проверки приложения и документов

```sh
# Go, examples/m17-sdd
go test -race ./...
go vet ./...
go build ./...
# OpenSpec, та же папка
npm run verify:specs
# Java, examples/java
mvn -pl m17-sdd -am verify
```

Проверки обоих стеков нужны автору курса; студент сдаёт выбранный стек. Эталон — средство сверки после попытки, не требование переписать тесты под готовую реализацию. В [матрице reference](lab/traceability-reference.md) перечислены покрытие и ограничения.
