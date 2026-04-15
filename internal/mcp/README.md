# MCP (Model Context Protocol) 集成

本目录包含 WeKeep 项目中 MCP HTTP Server 的集成实现，提供微信公众号文章收藏相关的 AI 工具接口。

## 概述

MCP (Model Context Protocol) 是一个标准化的协议，用于在 AI 模型和外部工具之间建立通信。本项目使用 [mcp-go](https://github.com/mark3labs/mcp-go) 库实现了 HTTP MCP 服务器。

## 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                    WeKeep Server (:8000)                 │
├─────────────────────────────────────────────────────────┤
│  /api/v1/*     │  REST API (Articles, Authors, etc.)    │
│  /mcp/*        │  MCP HTTP Server (StreamableHTTP)      │
└─────────────────────────────────────────────────────────┘
```

**优势**：
- ✅ 无需额外端口
- ✅ 统一服务管理
- ✅ 共享数据库连接和配置
- ✅ 简化部署

## 目录结构

```
internal/mcp/
├── server.go          # MCP 服务器创建和配置
├── tools.go           # 微信公众号相关工具实现
└── README.md          # 本文档
```

## 可用工具

### 1. echo - 测试连通性

回显用户输入的消息，用于测试 MCP 连接。

**参数**:
- `message` (string, 必需): 要回显的消息内容

### 2. get_server_info - 获取服务器信息

获取 MCP 服务器的运行信息。

**参数**: 无

### 3. wechat_parse_url - 解析微信文章

从微信公众号文章 URL 解析文章内容（标题、作者、正文等），不保存到数据库。

**参数**:
- `url` (string, 必需): 微信公众号文章 URL

**示例请求**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "wechat_parse_url",
    "arguments": {
      "url": "https://mp.weixin.qq.com/s/xxxxx"
    }
  }
}
```

### 4. wechat_save_article - 保存微信文章

从微信公众号文章 URL 解析并保存文章到数据库。

**参数**:
- `url` (string, 必需): 微信公众号文章 URL
- `tags` (array, 可选): 文章标签列表

**示例**:
```json
{
  "name": "wechat_save_article",
  "arguments": {
    "url": "https://mp.weixin.qq.com/s/xxxxx",
    "tags": ["技术", "Go"]
  }
}
```

### 5. wechat_list_articles - 查询文章列表

查询已收藏的微信公众号文章列表。

**参数**:
- `page_num` (number, 可选): 页码，从1开始，默认1
- `page_size` (number, 可选): 每页数量，默认20，最大100
- `author_id` (number, 可选): 作者ID筛选
- `tags` (array, 可选): 标签筛选
- `order_by` (string, 可选): 排序方式，如 "date_added desc"

### 6. wechat_get_article - 获取文章详情

根据文章 ID 获取已收藏的微信公众号文章完整内容。

**参数**:
- `id` (number, 必需): 文章ID

### 7. wechat_search_articles - 搜索文章

搜索已收藏的微信公众号文章。

**参数**:
- `keyword` (string, 必需): 搜索关键词
- `title_only` (boolean, 可选): 是否只搜索标题，默认 false（全文搜索）
- `page_num` (number, 可选): 页码
- `page_size` (number, 可选): 每页数量

### 8. wechat_get_tags - 获取所有标签

获取所有已使用的文章标签，按使用频率排序。

**参数**: 无

### 9. wechat_get_stats - 获取统计数据

获取文章收藏统计数据。

**参数**: 无

### 10. wechat_delete_article - 删除文章

根据文章 ID 删除已收藏的微信公众号文章。

**参数**:
- `id` (number, 必需): 要删除的文章ID

## 使用方法

### 启动服务器

```bash
# 使用 gf run
gf run

# 或直接运行
go run main.go
```

服务器将在 `http://localhost:8000` 启动，MCP 端点位于 `/mcp`。

### API 端点

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/mcp` | MCP 协议端点（支持所有 MCP 方法） |

### 测试工具

#### 1. 初始化会话

```bash
curl -X POST http://localhost:8000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'
```

#### 2. 列出可用工具

```bash
curl -X POST http://localhost:8000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'
```

#### 3. 保存微信文章

```bash
curl -X POST http://localhost:8000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "wechat_save_article",
      "arguments": {
        "url": "https://mp.weixin.qq.com/s/your-article-id",
        "tags": ["技术", "Go"]
      }
    }
  }'
```

#### 4. 搜索文章

```bash
curl -X POST http://localhost:8000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "wechat_search_articles",
      "arguments": {
        "keyword": "GoFrame",
        "page_size": 10
      }
    }
  }'
```

## Python 客户端示例

```python
import requests
import json

class WeKeepMCPClient:
    def __init__(self, base_url="http://localhost:8000/mcp"):
        self.base_url = base_url
        self.request_id = 0

    def call(self, method, params=None):
        self.request_id += 1
        payload = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        response = requests.post(
            self.base_url,
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        return response.json()

    def initialize(self):
        return self.call("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "python-client", "version": "1.0.0"}
        })

    def list_tools(self):
        return self.call("tools/list")

    def save_article(self, url, tags=None):
        args = {"url": url}
        if tags:
            args["tags"] = tags
        return self.call_tool("wechat_save_article", args)

    def search_articles(self, keyword, title_only=False):
        return self.call_tool("wechat_search_articles", {
            "keyword": keyword,
            "title_only": title_only
        })

    def get_article(self, article_id):
        return self.call_tool("wechat_get_article", {"id": article_id})

    def list_articles(self, page_num=1, page_size=20):
        return self.call_tool("wechat_list_articles", {
            "page_num": page_num,
            "page_size": page_size
        })

    def call_tool(self, name, arguments):
        return self.call("tools/call", {
            "name": name,
            "arguments": arguments
        })

# 使用示例
if __name__ == "__main__":
    client = WeKeepMCPClient()

    # 初始化
    print("初始化:", client.initialize())

    # 列出工具
    print("工具列表:", client.list_tools())

    # 保存文章
    print("保存文章:", client.save_article(
        "https://mp.weixin.qq.com/s/xxxxx",
        ["技术", "Go"]
    ))

    # 搜索文章
    print("搜索:", client.search_articles("GoFrame"))

    # 获取文章详情
    print("文章详情:", client.get_article(1))
```

## Claude Desktop 配置

如果你想在 Claude Desktop 中使用此 MCP 服务器，需要通过 stdio 模式运行。由于本实现是 HTTP 模式，你需要：

1. 使用 `mcp-proxy` 或类似工具将 HTTP 转换为 stdio
2. 或者使用支持 HTTP MCP 的客户端

## 设计原则

本实现遵循以下原则：

- **KISS**: 保持简单，只提供必要的微信文章操作工具
- **DRY**: 工具注册逻辑统一封装，复用现有 Service 层
- **可扩展性**: 清晰的目录结构，便于添加新工具
- **文档化**: 完整的 README，方便他人参考

## 参考资源

- [mcp-go 官方文档](https://github.com/mark3labs/mcp-go)
- [MCP 协议规范](https://modelcontextprotocol.io/)
- [GoFrame 官方文档](https://goframe.org/)

## 常见问题

### Q: 如何调试 MCP 工具？

A: 使用 curl 或 Postman 发送 JSON-RPC 请求，查看响应内容。也可以在工具处理函数中添加日志：

```go
func handleMyTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    g.Log().Debugf(ctx, "收到工具调用: %+v", request)
    // ...
}
```

### Q: 如何处理错误？

A: 返回 error，mcp-go 会自动将其转换为 JSON-RPC 错误响应：

```go
if err != nil {
    return nil, fmt.Errorf("操作失败: %v", err)
}
```

## License

本项目采用 MIT License。
