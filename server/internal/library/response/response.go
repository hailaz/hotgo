// Package response
package response

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"hotgo/internal/library/contexts"
	"hotgo/internal/model"
)

// JsonExit 返回JSON数据并退出当前HTTP执行函数
func JsonExit(r *ghttp.Request, code int, message string, data ...any) {
	RJson(r, code, message, data...)
	r.Exit()
}

// RXml xml
func RXml(r *ghttp.Request, code int, message string, data ...any) {
	responseData := any(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	res := &model.Response{
		Code:      code,
		Message:   message,
		Timestamp: gtime.Timestamp(),
		TraceID:   gctx.CtxId(r.Context()),
	}

	// 如果不是正常的返回，则将data转为error
	if gcode.CodeOK.Code() == code {
		res.Data = responseData
	} else {
		res.Error = responseData
	}

	// 清空响应
	r.Response.ClearBuffer()

	// 写入响应
	r.Response.WriteXml(gconv.Map(res))

	// 加入到上下文
	contexts.SetResponse(r.Context(), res)
}

// RJson 标准返回结果数据结构封装
func RJson(r *ghttp.Request, code int, message string, data ...any) {
	responseData := any(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	res := &model.Response{
		Code:      code,
		Message:   message,
		Timestamp: gtime.Timestamp(),
		TraceID:   gctx.CtxId(r.Context()),
	}

	// 如果不是正常的返回，则将data转为error
	if gcode.CodeOK.Code() == code {
		res.Data = responseData
	} else {
		res.Error = responseData
	}

	// 清空响应
	r.Response.ClearBuffer()

	// 写入响应
	r.Response.WriteJson(res)

	// 加入到上下文
	contexts.SetResponse(r.Context(), res)
}

// CustomJson 自定义JSON
func CustomJson(r *ghttp.Request, content any) {
	// 清空响应
	r.Response.ClearBuffer()

	// 写入响应
	r.Response.WriteJson(content)

	// 加入到上下文
	contexts.SetResponse(r.Context(), &model.Response{
		Code:      0,
		Message:   "",
		Data:      content,
		Error:     nil,
		Timestamp: gtime.Timestamp(),
		TraceID:   gctx.CtxId(r.Context()),
	})
}

// Redirect 重定向
func Redirect(r *ghttp.Request, location string, code ...int) {
	r.Response.RedirectTo(location, code...)
}

// RText 返回成功文本
func RText(r *ghttp.Request, message string) {
	// 清空响应
	r.Response.ClearBuffer()

	// 写入响应
	r.Response.Write(message)
}
