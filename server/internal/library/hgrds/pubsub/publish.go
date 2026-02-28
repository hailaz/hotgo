// Package pubsub
package pubsub

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

// Publish 推送消息
func Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	return g.Redis().Publish(ctx, channel, message)
}
