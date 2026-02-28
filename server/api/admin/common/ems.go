// Package common
package common

import (
	"github.com/gogf/gf/v2/frame/g"
)

// SendTestEmailReq 发送测试邮件
type SendTestEmailReq struct {
	g.Meta `path:"/ems/sendTest" tags:"邮件" method:"post" summary:"发送测试邮件"`
	To     string `json:"to" v:"required#接收者邮件不能为空" dc:"接收者邮件，多个用;隔开"`
}

type SendTestEmailRes struct {
}

// SendBindEmsReq 发送换绑邮件
type SendBindEmsReq struct {
	g.Meta `path:"/ems/sendBind" tags:"邮件" method:"post" summary:"发送换绑邮件"`
}

type SendBindEmsRes struct {
}
