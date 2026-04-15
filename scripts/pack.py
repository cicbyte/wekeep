#!/usr/bin/env python3
"""GoFrame 资源打包脚本

用法:
    python pack.py                  # 打包到 internal/packed/packed_resource.go
"""

import os
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
PACK_SRC = "resource/public,resource/sql"
PACK_DST = str(ROOT / "internal" / "packed" / "packed_resource.go")
PACKAGE = "packed"


def main():
    os.makedirs(os.path.dirname(PACK_DST), exist_ok=True)

    cmd = ["gf", "pack", PACK_SRC, PACK_DST, "-n", PACKAGE]
    result = subprocess.run(cmd, cwd=ROOT)
    print(f">>> {' '.join(cmd)}")
    result = subprocess.run(cmd)

    if result.returncode != 0:
        sys.exit(1)
    print(f"Done: {PACK_DST}")


if __name__ == "__main__":
    main()
