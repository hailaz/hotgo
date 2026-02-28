// Package consts
package consts

type CtxKey string

// ContextKey 上下文
const (
	ContextHTTPKey     CtxKey = "httpContext" // http上下文变量名称
	ContextKeyCronArgs CtxKey = "cronArgs"    // 定时任务参数
	ContextKeyCronSn   CtxKey = "cronSn"      // 定时任务序列号
)
