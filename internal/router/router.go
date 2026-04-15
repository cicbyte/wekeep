package router

import (
	"context"

	controller "github.com/ciclebyte/wekeep/internal/controller"
	meilisearchController "github.com/ciclebyte/wekeep/internal/controller/meilisearch"
	storageController "github.com/ciclebyte/wekeep/internal/controller/storage"

	"github.com/ciclebyte/wekeep/internal/service"
	"github.com/ciclebyte/wekeep/library/libRouter"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Router struct{}

func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().MiddlewareCORS)

		group.Bind(
			controller.Categories,
			controller.Articles,
			controller.Authors,
			controller.Health,
			controller.Images,
			storageController.Storage,
		)

		// 图片文件访问代理（不走 JSON 序列化）
		group.GET("/images/file/*", controller.File)

		//自动绑定定义的控制器
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			panic(err)
		}
	})
}

// BindMeilisearchController 绑定 Meilisearch 搜索控制器
func (router *Router) BindMeilisearchController(ctx context.Context, group *ghttp.RouterGroup) {
	group.Bind(
		meilisearchController.Search,
	)
}
