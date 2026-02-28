// Package consts
package consts

// 定时任务

const (
	CronSplitStr     = "," // 变量分割符
	CronPolicySame   = 1   // 并行策略
	CronPolicySingle = 2   // 单例策略
	CronPolicyOnce   = 3   // 单次策略
	CronPolicyTimes  = 4   // 多次策略
)
