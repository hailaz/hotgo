// Package router
package router

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"hotgo/addons/hgexample/controller/api"
	"hotgo/addons/hgexample/global"
	"hotgo/internal/consts"
	"hotgo/internal/library/addons"
	"hotgo/internal/service"
)

// Api 前台路由
func Api(ctx context.Context, group *ghttp.RouterGroup) {
	prefix := addons.RouterPrefix(ctx, consts.AppApi, global.GetSkeleton().Name)
	group.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Bind(
			// 无需验证的路由
			api.Index,
		)
		group.Middleware(service.Middleware().ApiAuth)
		group.Bind(
		// 需要验证的路由
		// ...
		)
	})
}
