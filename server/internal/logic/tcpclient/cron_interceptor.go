// Package tcpclient
package tcpclient

import (
	"context"

	"hotgo/internal/library/network/tcp"
)

// DefaultInterceptor 默认拦截器
func (s *sCronClient) DefaultInterceptor(ctx context.Context, msg *tcp.Message) (err error) {
	// conn := tcp.ConnFromCtx(ctx)
	// g.Log().Warningf(ctx, "DefaultInterceptor msg:%+v, conn:%+v", msg, gjson.New(conn).String())
	return
}
