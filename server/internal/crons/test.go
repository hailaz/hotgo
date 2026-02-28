// Package crons
package crons

import (
	"context"
	"time"

	"hotgo/internal/library/cron"
)

func init() {
	cron.Register(Test)
}

// Test 测试任务（无参数）
var Test = &cTest{name: "test"}

type cTest struct {
	name string
}

func (c *cTest) GetName() string {
	return c.name
}

// Execute 执行任务
func (c *cTest) Execute(ctx context.Context, parser *cron.Parser) (err error) {
	parser.Logger.Infof(ctx, "cron test Execute:%v", time.Now())
	return
}
