# SDD / OpenSpec Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox syntax for tracking.

**Goal:** Implement the approved standalone SDD module and a small two-language Access Requests laboratory.
**Architecture:** Existing static course page plus dependency-free local HTTP examples; shared OpenSpec documents explain baseline, deltas and final contract. Independent file ownership permits parallel implementation, with integration and review in the main chat.
**Tech Stack:** HTML/CSS/JavaScript, Go stdlib, Java 21 stdlib/JUnit, OpenSpec pinned CLI.
**Spec:** ../specs/2026-11-05-sdd-openspec-design.md (approved 2026-11-05, including command cheat sheet).

## Global Constraints

- Preserve existing module ids/progress keys; append m17 and recommend after m04.
- Ten substantive theoretical blocks before commands/practice; Russian prose, OpenSpec syntax unchanged.
- Local memory only, loopback HTTP, demo identities explicitly labelled; no production IAM or paid APIs.
- Two equivalent language implementations, baseline and reference modes; same observable contract.
- No publication, push, global tool installation or changes to unrelated .idea files.

## Review Focus

- Exact expiry boundary and repeat decisions must not extend or revive grants.
- Identity comes from a fixed demo adapter, never an arbitrary role claim.
- MODIFIED must preserve unchanged scenarios and final main specs must match both changes.
- Strict validation and AI verify must not be described as executing behavioral tests.
- Both stack/theme modes and old progress/navigation must survive m17 registration.

### Task 1: Laboratory behavior (Go and Java)

**Files:** examples/m17-sdd/{go.mod, implementation, tests, cmd}; examples/java/m17-sdd/{pom.xml,src}; examples/java/pom.xml.
**Interface:** HTTP POST /requests (resource, hours), GET /requests/{id}, POST /requests/{id}/decision (decision), POST /requests/{id}/revoke (reason), GET /access?user=&resource=. X-User demo identities alice/bob (employees), carol (approver). Return JSON, standard 400/401/403/404/409; valid creates 201 and other success 200. Core transitions pending→approved/rejected; repeated same decision idempotent, opposite decision conflict; self approval forbidden. Baseline creates indefinite grants; reference requires integer hours 1..24 and supports revoke. Clock is injectable; legacy grants retain indefinite validity. Publish exact runnable commands to Task 2.
- [x] Write behavioral tests first: identity/foreign read, transitions, repeats, hour limits, now==expiry, preservation of baseline grants after policy upgrade, revoke never undone by replay, concurrent conflicting decisions.
- [x] Run tests before implementation; record expected missing-behavior failures.
- [x] Implement using standard library, mutex/synchronized transitions and controlled clock; bind loopback.
- [x] Run Go tests with race/vet/build and Java Maven verify. Expected green; if tools unavailable install only in /tmp or report specific limitation.

### Task 2: Shared OpenSpec laboratory materials

**Files:** examples/m17-sdd/README.md, lab/, starter/openspec/, reference/openspec/, locally pinned package manifest/lock, validation script.
**Interfaces:** Task 1 HTTP contract; 3 capabilities access-requests/access-decisions/access-grants. Keep starter/current main specs distinct from final/reference snapshots. Student copies starter, runs baseline/reference in selected language.
- [x] Research and pin available released OpenSpec; inspect bundled help/schema, not only moving main docs.
- [x] Write brief, stakeholder answers, role/state tables, traceability template and filled solution; two archived change packages and separate conflicting extension exercise. Include ADDED/MODIFIED plus valid REMOVED and RENAMED exercise and full replacement warning.
- [x] Run strict CLI validation of baseline/final and deltas; test archive transformations on temporary copies and compare normalized resulting requirements/scenarios to final main specs. Expected exact final behavior; no placeholder Purpose.
- [x] Document installation, setup, four exercises, expected outcomes, safe demo identity and limitations; include pinned-version commands.

### Task 3: Course page and integration

**Files:** ai-engineer-course/m17.html, assets/course.js, index.html, m04.html, m13.html, m15.html, projects.html; docs/upgrade/check-site.cjs if verification needs stack/theme support.
**Interfaces:** Existing Course SVG/theme utilities, .viz/.task/details.q and data-save m17-t1..4. Lab links must resolve under deployed course directory or use source repository links; do not link ../examples paths that Pages does not deploy.
- [x] Write substantial 10 theory sections, worked analysis/spec/delta examples, 3 useful diagrams, command table and optional/expanded paths, 4 laboratory stages, >=10 selfchecks and >=8 interview answers.
- [x] Register m17; add learning-order/time guidance and contextual links without renumbering. Estimated course ≈197–199h plus optional voice; original 20-week path +1–2 weeks.
- [x] Run static structural/link/script checks, then existing browser check and targeted both-theme/stack interactive checks. Expected no console errors, dead links or empty diagrams; all progress controls work.

### Task 4: Integration verification and review

- [x] Reconcile module commands with installed pinned OpenSpec and actual project run instructions.
- [x] Run examples/verify.sh; report any pre-existing/environment failures separately. Validate changed Java module in aggregator context.
- [x] Fresh reviewer checks whole change against approved design and Review Focus; fix actionable defects and rerun affected checks.
- [x] Record evidence in progress report; deliver local page/source links and exact remaining limitations. Leave changes reviewable without publishing or committing unrelated files.

## Execution notes

Approval of the design authorizes implementation. Routine execution decisions remain with the implementer; no additional approval round for the plan. Parallel workers own separate files. Work in current authorized workspace on codex/sdd-openspec when Git permits; do not create an unrequested worktree or disturb untracked design/.idea files. Toolchain downloads and generated outputs stay under /tmp.

## Execution record

- User approved design; implemented independent page, code and OpenSpec material tasks with separate file ownership.
- Branch creation initially blocked by sandbox; authorized escalation created codex/sdd-openspec. No worktree or global installs/config changes, no commits/push/deployment.
- OpenSpec pinned to 1.16.0; custom workflows configured in isolated XDG directories. Actual generated integration names documented.
- Go and Java each have 7 behavior tests. Go race/vet/build and Java targeted reactor verify pass.
- 20/20 course pages pass browser checks; m17 go/java × light/dark and mobile width390 pass.
- Sequential real archive transformations, complete spec/scenario comparisons and strict validations pass.
- Review fixed successful final creation scenario to include hours1; full old scenario names retained with updated valid inputs. No outstanding Critical/Important findings.
- Global examples/verify.sh: 13 Go modules and 5 Python syntax checks pass; Java overall fails on pre-existing m16-voice VoiceTest.turnStreamsBySentence wall-clock timing assertion. Unchanged voice code left intact; m17 verified separately.
