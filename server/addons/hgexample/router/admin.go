// Package router
package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"hotgo/addons/hgexample/controller/admin/sys"
	"hotgo/addons/hgexample/global"
	"hotgo/addons/hgexample/router/genrouter"
	"hotgo/internal/consts"
	"hotgo/internal/library/addons"
	"hotgo/internal/service"
)

func Admin(ctx context.Context, group *ghttp.RouterGroup) {
	prefix := addons.RouterPrefix(ctx, consts.AppAdmin, global.GetSkeleton().Name)
	group.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Bind(
			sys.Index,
		)
		group.Middleware(service.Middleware().AdminAuth)
		group.Bind(
			sys.Comp,
			sys.Config,
			sys.Table,
			sys.TreeTable,
		)
	})

	// 注册生成路由
	genrouter.Register(ctx, group)
}
