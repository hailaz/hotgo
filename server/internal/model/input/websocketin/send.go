// Package websocketin
package websocketin

import (
	"hotgo/internal/websocket"
)

// SendToTagInp 发送标签消息
type SendToTagInp struct {
	Tag      string              `json:"tag" v:"required#tag不能为空" description:"标签"`
	Response websocket.WResponse `json:"response" v:"required#response不能为空"  description:"响应内容"`
}
