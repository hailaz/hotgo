// Package router
package router

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"hotgo/addons/hgexample/controller/home"
	"hotgo/addons/hgexample/global"
	"hotgo/internal/consts"
	"hotgo/internal/library/addons"
)

// Home 前台页面路由
func Home(ctx context.Context, group *ghttp.RouterGroup) {
	prefix := addons.RouterPrefix(ctx, consts.AppHome, global.GetSkeleton().Name)
	group.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Bind(
			home.Index,
		)
	})
}
