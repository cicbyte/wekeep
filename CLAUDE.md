# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

**WeKeep** — 微信公众号文章收藏管理工具。支持文章抓取、全文搜索（Meilisearch）、图片本地化存储（本地/S3）、MCP 协议 AI 集成。

### 技术栈

**后端**: Go 1.24 + GoFrame v2.10.0 + MySQL
**前端**: React 19 + TypeScript 5.8 + Vite 6.2 + Tailwind CSS (CDN)
**搜索**: Meilisearch
**存储**: 本地文件系统 / RustFS（S3 兼容，通过 AWS SDK v2）
**AI**: Gemini API（前端文章解析）+ MCP Server（后端 AI 工具集成）
**模块路径**: `github.com/cicbyte/wekeep`

## 常用命令

### 后端

```bash
gf run                    # 运行项目（开发模式）
make build                # 构建
make dao                  # 生成 DAO/DO/Entity
make ctrl                 # 生成控制器代码
make service              # 生成 Service 接口代码
gf gen dao                # 等同于 make dao
```

### 前端

```bash
cd web && npm run dev      # 开发服务器（端口 8898）
cd web && npm run build    # 构建（输出到 web/dist/）
```

构建产物需复制到 `resource/public/html/` 供 Go 后端托管。

### 部署

```bash
make image                # Docker 镜像（标签 = git commit short hash）
make image.push           # 构建并推送
make deploy               # kubectl 部署
```

## 架构

### 后端分层

```
API (api/v1/<module>/)         → 请求/响应结构体 + g.Meta 路由标签
Controller (internal/controller/) → HTTP 处理、参数校验
Service (internal/service/)    → 接口定义（I + 模块名）
Logic (internal/logic/<module>/) → 业务实现（s + 模块名，init() 自动注册）
DAO (internal/dao/)            → 自动生成，勿手动编辑
Model (internal/model/)        → entity/ do/ info/
```

路由注册入口: `internal/cmd/cmd.go` → `internal/router/router.go`

### 业务模块

| 模块 | 说明 |
|------|------|
| categories | 分类管理（CRUD） |
| articles | 文章收藏（CRUD、搜索、Gemini 解析） |
| authors | 作者管理 |
| images / article_images | 图片管理 |
| health | 健康检查 |
| meilisearch | 全文搜索（`/api/v1/search/*`） |
| storage | 对象存储（`/api/v1/storage/*`） |

### 独立子系统

**微信解析器** (`internal/parser/wechat.go`): 抓取微信公众号文章 HTML，提取标题/作者/正文/图片，转换为 Markdown。

**存储抽象层** (`internal/storage/`): Provider 接口，两种实现：
- `LocalProvider` — 本地文件系统
- `RustFSProvider` — S3 兼容对象存储（通过 AWS SDK v2）

配置通过 `storage.type` 切换（`config.yaml`）。

**MCP 服务器** (`internal/mcp/`): 基于 `mark3labs/mcp-go`，注册在 `/mcp/*`（StreamableHTTP）。提供 10 个工具：`echo`, `get_server_info`, `wechat_parse_url`, `wechat_save_article`, `wechat_list_articles`, `wechat_search_articles`, `wechat_get_article`, `wechat_update_tags`, `wechat_delete_article`, `wechat_batch_delete`。

**Meilisearch 客户端** (`library/libMeilisearch/`): 全文搜索封装，通过 `search.enabled` 控制是否启用。

### 前端结构

```
web/
  App.tsx              # 主应用 + 路由定义
  index.tsx            # 入口
  types.ts             # TypeScript 类型
  hooks/useArticles.ts # 文章状态管理
  services/
    apiService.ts      # 后端 API 调用
    geminiService.ts   # Gemini AI 文章解析
    migrationService.ts # localStorage 数据迁移
  components/          # 11 个 UI 组件
```

路由: `/` (仪表盘), `/list` (文章列表), `/authors` (作者), `/settings` (设置), `/read/:id` (阅读), `/edit/:id` (编辑)

### 数据库

数据库名: `wekeep`。迁移脚本按序号: `resource/sql/mysql/01_categories.sql` ~ `08_article_images.sql`

### 配置

- `manifest/config/config.yaml` — 主配置（服务器、数据库、存储、搜索）
- `hack/config.yaml` — GoFrame CLI 配置（DAO 生成、Docker 构建）

## 关键约定

### 依赖注入（最重要）

新增模块**必须**在 `internal/logic/logic.go` 添加空白导入，否则 `init()` 不执行，运行时 panic：

```
cmd.go → logic/logic.go → logic/<module>/<module>.go init() → service.RegisterXxx()
```

### 命名规范

- 控制器变量: `Categories = categoriesController{}`（大写导出）
- Service 接口: `ICategories`
- Logic 实现: `sCategories`（小写 s 前缀）
- DAO/DO/Entity: `make dao` 自动生成

### 错误处理

```go
err = g.Try(ctx, func(ctx context.Context) {
    _, err := dao.Xxx.Ctx(ctx).Insert(...)
    liberr.ErrIsNil(ctx, err, "操作失败")
})
// panic 会被 MiddlewareHandlerResponse 捕获转为错误响应
```

### API 定义

```go
type XxxReq struct {
    g.Meta `path:"/xxx" method:"post" tags:"分组" summary:"描述"`
    Name   string `json:"name" v:"required#名称不能为空"`
}
```

分页请求嵌入 `commonApi.PageReq`，响应嵌入 `commonApi.ListRes`。

### SPA 回退

非 API 请求（无扩展名）自动回退到 `resource/public/html/index.html`。MCP 路由 `/mcp/*` 也排除在回退之外。

## 新增模块清单

1. `api/v1/<module>/<module>.go` — 定义请求/响应
2. `internal/model/<module>.go` — 业务模型（Info）
3. 数据库建表 → `make dao`
4. `internal/service/<module>.go` — 接口（IXxx + RegisterXxx）
5. `internal/logic/<module>/<module>.go` — 实现 + init() 注册
6. `internal/logic/logic.go` — **添加空白导入**（容易遗漏！）
7. `internal/controller/<module>.go` — 控制器
8. `internal/router/router.go` — `group.Bind()` 注册
