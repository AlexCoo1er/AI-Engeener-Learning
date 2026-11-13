#!/bin/sh
# Validate on disposable copies; never sync/archive teaching snapshots in place.
set -eu
lab_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cli=${OPENSPEC_BIN:-"$lab_dir/node_modules/.bin/openspec"}
if [ ! -x "$cli" ]; then
  echo 'Install pinned tooling first: cd examples/m17-sdd && npm ci --ignore-scripts' >&2
  exit 1
fi
if [ "$("$cli" --version)" != '1.16.0' ]; then
  echo 'Expected OpenSpec 1.16.0; refusing another CLI.' >&2
  exit 1
fi
scratch=$(mktemp -d "${TMPDIR:-/tmp}/m17-specs.XXXXXX")
trap 'rm -rf "$scratch"' EXIT HUP INT TERM
export OPENSPEC_TELEMETRY=0
export XDG_CONFIG_HOME="$scratch/config"
export XDG_DATA_HOME="$scratch/data"
mkdir -p "$scratch/work" "$scratch/final" "$scratch/operations"
cp -R "$lab_dir/starter/openspec" "$scratch/work/"
cp -R "$lab_dir/reference/openspec" "$scratch/final/"
cp -R "$lab_dir/starter/openspec" "$scratch/operations/"
validate() { (cd "$1" && "$cli" validate --all --strict --no-interactive); }
validate "$scratch/work"
validate "$scratch/final"
(cd "$scratch/final" && "$cli" validate --archived --strict --no-interactive)
# Archived deltas must also be validated against the correct base; --archived alone checks tasks.
for change in temporary-access revoke-access; do
  mkdir -p "$scratch/work/openspec/changes/$change"
  cp -R "$lab_dir/reference/openspec/changes/archive/2026-11-05-$change/." "$scratch/work/openspec/changes/$change/"
  validate "$scratch/work"
  (cd "$scratch/work" && "$cli" archive "$change" --yes)
  validate "$scratch/work"
  if [ "$change" = temporary-access ]; then
    mkdir -p "$scratch/conflict"
    cp -R "$scratch/work/openspec" "$scratch/conflict/"
    mkdir -p "$scratch/conflict/openspec/changes/extend-access"
    cp -R "$lab_dir/lab/conflicting-extension/." "$scratch/conflict/openspec/changes/extend-access/"
    validate "$scratch/conflict"
  fi
done
python3 "$lab_dir/lab/compare-specs.py" "$scratch/work/openspec/specs" "$lab_dir/reference/openspec/specs"
mkdir -p "$scratch/operations/openspec/changes/operation-exercise"
cp -R "$lab_dir/lab/operation-exercise/." "$scratch/operations/openspec/changes/operation-exercise/"
validate "$scratch/operations"
(cd "$scratch/operations" && "$cli" archive operation-exercise --yes)
validate "$scratch/operations"
python3 - "$scratch/operations/openspec/specs" "$lab_dir/starter/openspec/specs" <<'PY'
from pathlib import Path
import sys
actual,base=map(Path,sys.argv[1:])
a=(actual/'access-requests/spec.md').read_text()
b=(base/'access-requests/spec.md').read_text()
assert 'RQ-2 Read request by ownership' in a and 'RQ-2 Read visible request' not in a
assert a.count('#### Scenario:') == b.count('#### Scenario:'), 'rename lost scenarios'
g=(actual/'access-grants/spec.md').read_text()
assert '### Requirement: GR-3 New grants are indefinite' not in g
assert 'indefinite validity and no expires_at' in g, 'remove lost preserved guarantee'
print('Standalone RENAMED/REMOVED archive: preserved scenarios and indefinite guarantee')
PY
printf '%s\n' 'OpenSpec 1.16.0: strict specs/deltas/archive transformations verified; final requirements match.'
