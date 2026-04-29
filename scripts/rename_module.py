"""Rename Go module path from github.com/ciclebyte/wekeep to github.com/cicbyte/wekeep."""

import re
from pathlib import Path

OLD = "github.com/ciclebyte/wekeep"
NEW = "github.com/cicbyte/wekeep"
ROOT = Path(__file__).resolve().parent.parent
EXCLUDE = {Path(__file__).resolve()}


def main():
    targets = []
    for f in ROOT.rglob("*"):
        if f in EXCLUDE or not f.is_file():
            continue
        try:
            text = f.read_text(encoding="utf-8")
        except (UnicodeDecodeError, PermissionError):
            continue
        if OLD in text:
            targets.append(f)

    if not targets:
        print("No files contain the old module path.")
        return

    for f in targets:
        text = f.read_text(encoding="utf-8")
        new_text = text.replace(OLD, NEW)
        if new_text != text:
            f.write_text(new_text, encoding="utf-8")
            rel = f.relative_to(ROOT)
            print(f"  {rel}")

    print(f"\nDone. Updated {len(targets)} file(s).")


if __name__ == "__main__":
    main()
