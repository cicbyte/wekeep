#!/usr/bin/env python3
"""WeKeep Docker 镜像构建脚本

用法:
    python docker_build.py                  # 默认 wekeep:latest
    python docker_build.py -t v1.0          # wekeep:v1.0
    python docker_build.py -n myapp -t v2   # myapp:v2
    python docker_build.py -o release.tar   # 指定输出文件名
"""

import argparse
import os
import subprocess
import sys
from pathlib import Path


def run(cmd, check=True):
    print(f">>> {' '.join(cmd)}")
    return subprocess.run(cmd, check=check)


def main():
    parser = argparse.ArgumentParser(description="构建 WeKeep Docker 镜像")
    parser.add_argument("-n", "--name", default="wekeep", help="镜像名称 (默认: wekeep)")
    parser.add_argument("-t", "--tag", default="latest", help="镜像标签 (默认: latest)")
    parser.add_argument("-o", "--output", default=None, help="输出文件名 (默认: {name}_{tag}.tar)")
    args = parser.parse_args()

    image = f"{args.name}:{args.tag}"
    output = args.output or f"{args.name}_{args.tag}.tar"
    root = Path(__file__).resolve().parent.parent

    # 1. 清除旧镜像
    rmi = run(["docker", "rmi", "-f", image], check=False)
    if rmi.returncode == 0:
        print(f"已清除旧镜像: {image}")

    # 2. 构建新镜像
    run(["docker", "build", "-t", image, str(root)], check=True)
    print(f"构建完成: {image}")

    # 3. 导出
    run(["docker", "save", "-o", output, image], check=True)
    print(f"已导出: {output}")


if __name__ == "__main__":
    main()
