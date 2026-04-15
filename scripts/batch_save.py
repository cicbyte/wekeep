#!/usr/bin/env python3
"""批量保存微信公众号文章到 WeKeep"""

import sys
import time
from pathlib import Path

import requests

DEFAULT_BASE_URL = "http://localhost:8000"
BATCH_SIZE = 5
DELAY = 1


def load_urls(file_path: str) -> list[str]:
    urls = []
    for line in Path(file_path).read_text(encoding="utf-8").splitlines():
        line = line.strip().strip('"').strip("'")
        if line.startswith("http"):
            urls.append(line)
    return urls


def get_author_id(base_url: str, author_name: str) -> int | None:
    """查找作者ID，不存在则创建"""
    if not author_name:
        return None
    # 查找已有作者
    resp = requests.get(f"{base_url}/api/v1/authors/select", timeout=10)
    if resp.status_code == 200:
        data = resp.json().get("data", {})
        for item in data.get("options", []):
            if isinstance(item, dict) and item.get("name") == author_name:
                return item.get("id")
    # 创建新作者
    resp = requests.post(
        f"{base_url}/api/v1/authors/add",
        json={"name": author_name},
        timeout=10,
    )
    if resp.status_code == 200:
        result = resp.json()
        if result.get("code") == 0:
            return result.get("data", {}).get("id")
    return None


def save_article(base_url: str, url: str) -> bool:
    try:
        # 1. 解析文章
        resp = requests.post(
            f"{base_url}/api/v1/articles/parse-by-url",
            json={"url": url},
            timeout=60,
        )
        if resp.status_code != 200:
            print(f"  PARSE FAIL: HTTP {resp.status_code}")
            return False
        parse_data = resp.json().get("data", {})
        title = parse_data.get("title", "")
        author = parse_data.get("author", "")
        content = parse_data.get("content", "")
        if not title:
            print(f"  SKIP: no title")
            return False

        # 2. 查找/创建作者
        author_id = get_author_id(base_url, author)

        # 3. 保存文章
        if not author_id:
            print(f"  FAIL: author '{author}' not found and create failed")
            return False
        add_data = {"title": title, "url": url, "content": content, "authorId": author_id}
        resp = requests.post(
            f"{base_url}/api/v1/articles/add",
            json=add_data,
            timeout=30,
        )
        result = resp.json()
        if resp.status_code == 200 and result.get("code") == 0:
            article_id = result.get("data", {}).get("id", "?")
            print(f"  OK -> id={article_id}")
            return True
        print(f"  SAVE FAIL: {result.get('message', '')}")
        return False
    except Exception as e:
        print(f"  ERROR: {e}")
        return False


def main():
    base_url = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_BASE_URL
    url_file = sys.argv[2] if len(sys.argv) > 2 else "url列表.txt"

    urls = load_urls(url_file)
    if not urls:
        print("No URLs found")
        sys.exit(1)

    print(f"Loaded {len(urls)} URLs")
    print(f"Target: {base_url}")
    print()

    success = 0
    failed = 0

    for i, url in enumerate(urls, 1):
        print(f"[{i}/{len(urls)}] {url}")
        if save_article(base_url, url):
            success += 1
        else:
            failed += 1
        if i % BATCH_SIZE == 0 and i < len(urls):
            time.sleep(DELAY)

    print()
    print(f"Done: {success} success, {failed} failed")


if __name__ == "__main__":
    main()
