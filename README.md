# WeKeep

微信公众号文章收藏管理工具。支持文章抓取、全文搜索、图片本地化存储、MCP 协议 AI 集成。

## 功能特性

- **文章收藏** — 抓取微信公众号文章，提取标题/作者/正文/图片，转换为 Markdown
- **全文搜索** — 基于 Meilisearch 的全文检索，支持标题、作者、摘要、内容搜索
- **图片本地化** — 自动下载文章图片到本地/S3 存储，支持存储后端迁移
- **AI 集成** — MCP Server 提供 10 个工具，支持 AI 客户端（如 Claude）直接操作文章
- **分类管理** — 文章分类、标签体系
- **作者管理** — 自动提取作者信息，按作者浏览文章

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.24 / GoFrame v2.10.0 / MySQL |
| 前端 | React 19 / TypeScript 5.8 / Vite 6.2 / Tailwind CSS |
| 搜索 | Meilisearch |
| 存储 | 本地文件系统 / RustFS（S3 兼容） |
| AI | Gemini API（文章解析）/ MCP Server（AI 工具集成） |

## 快速开始

### 环境要求

- Go 1.24+
- Node.js 18+
- MySQL 8.0+
- Meilisearch（可选，用于全文搜索）

### 后端

```bash
# 克隆项目
git clone https://github.com/cicbyte/wekeep.git
cd wekeep

# 初始化数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS wekeep DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci"
mysql -u root -p wekeep < resource/sql/mysql/init.sql

# 修改配置
# 编辑 manifest/config/config.yaml，配置数据库连接、存储、搜索等

# 生成 DAO 代码
make dao

# 运行
gf run
```

后端默认监听 `:8000`，API 文档访问 `http://localhost:8000/swagger`。

### 前端

```bash
cd web
npm install
npm run dev    # 开发服务器，端口 8898
npm run build  # 构建到 web/dist/
```

构建产物需复制到 `resource/public/html/` 供 Go 后端托管。

## 项目结构

```
wekeep/
├── api/v1/                  # API 请求/响应定义（g.Meta 路由标签）
├── internal/
│   ├── controller/          # HTTP 控制器
│   ├── service/             # Service 接口定义
│   ├── logic/               # 业务逻辑实现（init() 自动注册）
│   ├── dao/                 # 数据访问层（自动生成）
│   ├── model/               # 数据模型（entity/do/info）
│   ├── parser/              # 微信文章 HTML 解析器
│   ├── storage/             # 存储抽象层（Local / RustFS）
│   ├── mcp/                 # MCP Server（10 个工具）
│   └── router/              # 路由注册
├── library/
│   ├── libMeilisearch/      # Meilisearch 客户端封装
│   └── libRouter/           # 路由自动绑定
├── web/                     # React 前端
│   ├── components/          # UI 组件
│   ├── services/            # API 调用服务
│   └── hooks/               # 状态管理
├── resource/
│   ├── config/config.yaml   # 主配置
│   ├── public/html/         # 前端构建产物
│   └── sql/                 # 数据库脚本
│       ├── mysql/init.sql   # MySQL 一键初始化
│       └── sqlite/init.sql  # SQLite 一键初始化
└── hack/                    # GoFrame CLI 配置
```

## 配置说明

主配置文件 `manifest/config/config.yaml`：

```yaml
server:
  address: ":8000"

database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/wekeep?charset=utf8mb4&parseTime=true&loc=Local"

storage:
  type: "local"              # local 或 rustfs
  local:
    basePath: "./uploads"
  rustfs:
    endpoint: "http://localhost:9000"
    bucket: "wekeep"
  migration:
    enabled: true            # 是否启用存储迁移功能

search:
  enabled: true              # 是否启用 Meilisearch 全文搜索
  meilisearch:
    address: "http://localhost:7700"
```

## API

后端提供 RESTful API，主要端点：

| 模块 | 路径 | 说明 |
|------|------|------|
| 文章 | `/api/v1/articles` | CRUD、搜索、Gemini 解析 |
| 分类 | `/api/v1/categories` | 分类管理 |
| 作者 | `/api/v1/authors` | 作者管理 |
| 图片 | `/api/v1/images` | 图片管理、文件代理 |
| 搜索 | `/api/v1/search/*` | 全文搜索 |
| 存储 | `/api/v1/storage/*` | 存储后端管理 |
| MCP | `/mcp/*` | MCP Server（StreamableHTTP） |

## 部署

```bash
make build        # 构建二进制
make image        # Docker 镜像
make image.push   # 构建并推送
make deploy       # kubectl 部署
```

## License

MIT
