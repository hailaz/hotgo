// Package crons
package crons

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/internal/library/cron"
	"time"
)

func init() {
	cron.Register(Test2)
}

// Test2 测试2任务（带参数）
var Test2 = &cTest2{name: "test2"}

type cTest2 struct {
	name string
}

func (c *cTest2) GetName() string {
	return c.name
}

// Execute 执行任务
func (c *cTest2) Execute(ctx context.Context, parser *cron.Parser) (err error) {
	if len(parser.Args) != 3 {
		err = gerror.New("传入参数不正确!")
		return
	}

	var (
		name = parser.Args[0]
		age  = parser.Args[1]
		msg  = parser.Args[2]
	)

	parser.Logger.Infof(ctx, "cron test2 Execute:%v, name:%v, age:%v, msg:%v", time.Now(), name, age, msg)
	return
}
