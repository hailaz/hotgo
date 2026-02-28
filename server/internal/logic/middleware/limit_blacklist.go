// Package middleware
package middleware

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"hotgo/internal/library/response"
	sysLogic "hotgo/internal/logic/sys"
)

// Blacklist IP黑名单限制中间件
func (s *sMiddleware) Blacklist(r *ghttp.Request) {
	if err := sysLogic.SysBlacklist().VerifyRequest(r); err != nil {
		response.JsonExit(r, gerror.Code(err).Code(), err.Error())
	}
	r.Middleware.Next()
}
