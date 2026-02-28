// Package common
package common

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/internal/model/input/sysin"
)

// SendTestSmsReq 发送测试短信
type SendTestSmsReq struct {
	g.Meta `path:"/sms/sendTest" tags:"短信" method:"post" summary:"发送测试短信"`
	sysin.SendCodeInp
}

type SendTestSmsRes struct {
}

// SendBindSmsReq 发送换绑短信
type SendBindSmsReq struct {
	g.Meta `path:"/sms/sendBind" tags:"短信" method:"post" summary:"发送换绑短信"`
}

type SendBindSmsRes struct {
}

// SendSmsReq 发送短信
type SendSmsReq struct {
	g.Meta `path:"/sms/send" tags:"短信" method:"post" summary:"发送短信"`
	sysin.SendCodeInp
}

type SendSmsRes struct {
}
