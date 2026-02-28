// Package model
package model

// Response HTTP响应
type Response struct {
	Code      int    `json:"code" example:"0" description:"状态码"`
	Message   string `json:"message,omitempty" example:"操作成功" description:"提示消息"`
	Data      any    `json:"data,omitempty" description:"数据集"`
	Error     any    `json:"error,omitempty" description:"错误信息"`
	Timestamp int64  `json:"timestamp" example:"1640966400" description:"服务器时间戳"`
	TraceID   string `json:"traceID" v:"0" example:"d0bb93048bc5c9164cdee845dcb7f820" description:"链路ID"`
}
