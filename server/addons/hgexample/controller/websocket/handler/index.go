// Package handler
package handler

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/websocket"
)

var (
	Index = cIndex{}
)

type cIndex struct{}

// TestMessage 测试消息
func (c *cIndex) TestMessage(client *websocket.Client, req *websocket.WRequest) {
	g.Log().Infof(client.Context(), "收到客户端测试消息:%v", gjson.New(req).String())
	// 将收到的消息原样发送给客户端
	websocket.SendSuccess(client, req.Event, req.Data)
}
