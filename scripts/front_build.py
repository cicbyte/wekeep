import logging as log
import os
import shutil
import subprocess
import sys
import time
from pathlib import Path

log.basicConfig(level=log.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

def run_command(cmd, cwd=None, shell=False, direct_output=False):
    """运行命令并实时显示输出
    Args:
        cmd: 要运行的命令
        cwd: 工作目录
        shell: 是否使用shell执行
        direct_output: 是否直接输出到终端（用于显示进度条等）
    """
    if direct_output:
        # 直接输出到终端，用于显示进度条等
        process = subprocess.Popen(
            cmd,
            stdout=sys.stdout,
            stderr=sys.stderr,
            cwd=cwd,
            shell=shell,
            universal_newlines=True,
        )
        return_code = process.wait()
        return return_code
    else:
        # 通过logging输出
        process = subprocess.Popen(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            encoding="utf-8",
            errors="replace",
            cwd=cwd,
            shell=shell,
            bufsize=1,
            universal_newlines=True,
        )

        # 实时读取并显示输出
        while True:
            output = process.stdout.readline()
            if output == "" and process.poll() is not None:
                break
            if output:
                log.info(output.strip())

        # 获取返回码
        return_code = process.poll()

        # 检查是否有错误输出
        stderr = process.stderr.read()
        if stderr:
            log.error(stderr.strip())

        return return_code


def run_npm_build():
    target_dir = Path(__file__).resolve().parent.parent / "web"
    target_dir_path = target_dir.resolve()
    log.info("target_dir_path: " + str(target_dir_path))

    # 在Windows系统上使用npm.cmd
    npm_cmd = "npm.cmd" if os.name == "nt" else "npm"

    try:
        return_code = run_command(
            [npm_cmd, "run", "build"], cwd=target_dir_path, shell=True
        )

        if return_code == 0:
            log.info("npm run build 命令执行成功")
        else:
            log.error("npm run build 命令执行失败")
            raise Exception("npm run build 命令执行失败")
    except Exception as e:
        log.error(f"执行npm命令时出错: {str(e)}")
        raise


def remove_old_build():
    target_dir = Path(__file__).resolve().parent.parent / "resource" / "public" / "html"

    # 检查并创建目标目录
    try:
        target_dir.mkdir(parents=True, exist_ok=True)
    except Exception as e:
        log.error(f"创建目录失败: {str(e)}")
        raise

    # 删除目录内容（包括子目录）
    for item in target_dir.glob("*"):
        try:
            if item.is_file():
                item.unlink()
            elif item.is_dir():
                shutil.rmtree(item)
        except Exception as e:
            log.error(f"删除 {item} 失败: {str(e)}")
            raise


def copy_new_build():
    source_dir = Path(__file__).resolve().parent.parent / "web" / "dist"
    target_dir = Path(__file__).resolve().parent.parent / "resource" / "public" / "html"

    # 确保源目录存在
    if not source_dir.exists():
        log.error("源目录不存在: " + str(source_dir))
        raise FileNotFoundError("源目录不存在")

    # 复制整个目录结构
    try:
        shutil.copytree(
            source_dir,
            target_dir,
            dirs_exist_ok=True,  # 允许目标目录已存在
        )
        log.info(f"成功复制文件从 {source_dir} 到 {target_dir}")
    except Exception as e:
        log.error(f"复制文件失败: {str(e)}")
        raise

def main():
    run_npm_build()
    remove_old_build()
    copy_new_build()

if __name__ == "__main__":
    main()
