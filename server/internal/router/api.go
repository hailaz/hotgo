// Package router
package router

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"hotgo/internal/consts"
	"hotgo/internal/controller/api/member"
	"hotgo/internal/controller/api/pay"
	"hotgo/internal/service"
	"hotgo/utility/simple"
)

// Api 前台路由
func Api(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
		group.Bind(
			pay.NewV1(), // 支付异步通知
		)
		group.Middleware(service.Middleware().ApiAuth)
		group.Bind(
			member.NewV1(), // 管理员
		)
	})
}
