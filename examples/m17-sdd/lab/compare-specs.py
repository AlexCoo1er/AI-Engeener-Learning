"""Compare complete normative blocks, allowing only archive whitespace changes."""
from pathlib import Path
import re
import sys


def requirements(root):
    result = {}
    for path in sorted(Path(root).glob('*/spec.md')):
        text = path.read_text()
        blocks = re.split(r'(?m)^### Requirement: ', text)[1:]
        if not blocks:
            raise SystemExit(f'No requirements: {path}')
        result[path.parent.name] = [re.sub(r'\s+', ' ', block).strip() for block in blocks]
    return result


actual, expected = map(requirements, sys.argv[1:])
if actual != expected:
    import difflib
    for capability in sorted(actual.keys() | expected.keys()):
        if actual.get(capability) != expected.get(capability):
            print('\n'.join(difflib.unified_diff(actual.get(capability, []), expected.get(capability, []), fromfile=f'archived/{capability}', tofile=f'reference/{capability}')))
    raise SystemExit('Archived normative requirements differ from final reference')
print('All three final capabilities match complete reference requirements and scenarios')
