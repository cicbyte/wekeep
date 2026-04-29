package cmd

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/cicbyte/wekeep/internal/logic"
	"github.com/cicbyte/wekeep/internal/mcp"
	"github.com/cicbyte/wekeep/internal/router"
	"github.com/cicbyte/wekeep/internal/storage"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 确保 SQLite 数据目录存在
			ensureDataDirs(ctx)

			// 自动初始化数据表
			autoMigrate(ctx)

			// 初始化存储模块
			initStorage(ctx)

			s := g.Server()
			// 配置静态文件服务
			s.SetServerRoot("resource/public/html/")

			// 创建 MCP HTTP 服务器
			mcpHandler := mcp.NewStreamableHTTPServer()

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				r := &router.Router{}
				r.BindController(ctx, group)

				// 注册 MCP 路由
				group.Group("/mcp", func(mcpGroup *ghttp.RouterGroup) {
					mcpGroup.ALL("/*", func(r *ghttp.Request) {
						mcpHandler.ServeHTTP(r.Response.Writer.ResponseWriter, r.Request)
					})
				})

				// 添加SPA路由回退支持，处理Vue Router的HTML5 History模式
				group.Hook("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
					path := r.URL.Path

					// 如果是 API 或 MCP 请求，跳过 SPA 回退
					if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/mcp") {
						return
					}

					// 如果是静态资源文件（有文件扩展名），跳过SPA回退
					if strings.Contains(path, ".") && !strings.HasSuffix(path, "/") {
						return
					}

					// 对于其他所有路径，都返回index.html，让Vue Router处理
					if path != "/" && !strings.HasPrefix(path, "/api/") {
						r.Response.ServeFile("resource/public/html/index.html")
						r.ExitAll()
					}
				})
			})
			s.Run()
			return nil
		},
	}
)

// ensureDataDirs 确保 SQLite 数据文件的父目录存在
func ensureDataDirs(ctx context.Context) {
	dbLink := g.Cfg().MustGet(ctx, "database.default.link").String()
	if !strings.HasPrefix(dbLink, "sqlite") {
		return
	}
	re := regexp.MustCompile(`@file\((.+)\)`)
	matches := re.FindStringSubmatch(dbLink)
	if len(matches) < 2 {
		return
	}
	dir := filepath.Dir(matches[1])
	if dir != "" && dir != "." {
		os.MkdirAll(dir, 0755)
	}
}

// initStorage 初始化存储模块
func initStorage(ctx context.Context) {
	// 读取存储配置
	cfg, err := g.Cfg().Get(ctx, "storage")
	if err != nil || cfg.IsEmpty() {
		g.Log().Warning(ctx, "Storage config not found, image localization disabled")
		return
	}

	// 解析配置
	var storageCfg storage.Config
	if err := cfg.Struct(&storageCfg); err != nil {
		g.Log().Errorf(ctx, "Failed to parse storage config: %v", err)
		return
	}

	// 初始化存储提供者
	if err := storage.Init(&storageCfg); err != nil {
		g.Log().Errorf(ctx, "Failed to initialize storage: %v", err)
		return
	}

	// 获取提供者并配置静态文件服务
	provider := storage.GetProvider()

	// 检查是否是本地存储
	if localProvider, ok := provider.(*storage.LocalProvider); ok {
		// 为本地存储配置静态文件服务
		basePath := localProvider.GetBasePath()
		baseURL := "/uploads"
		if storageCfg.Local != nil && storageCfg.Local.BaseURL != "" {
			baseURL = storageCfg.Local.BaseURL
		}
		// 添加静态文件路由
		g.Server().AddStaticPath(baseURL, basePath)
		g.Log().Infof(ctx, "Local storage static path configured: %s -> %s", baseURL, basePath)
	}

	// 检查并创建 bucket (RustFS)
	if rustfsProvider, ok := provider.(*storage.RustFSProvider); ok {
		exists, err := rustfsProvider.CheckBucketExists(ctx)
		if err != nil {
			g.Log().Warningf(ctx, "Failed to check bucket existence: %v", err)
		} else if !exists {
			g.Log().Info(ctx, "Bucket not found, creating...")
			if err := rustfsProvider.CreateBucket(ctx); err != nil {
				g.Log().Errorf(ctx, "Failed to create bucket: %v", err)
				return
			}
			g.Log().Info(ctx, "Bucket created successfully")
		}
	}

	// 健康检查
	if err := provider.HealthCheck(ctx); err != nil {
		g.Log().Warningf(ctx, "Storage health check failed: %v", err)
	}

	g.Log().Info(ctx, "Storage initialized successfully")
}
