// Package sysin
//
package sysin

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
)

// IndexTestInp 测试
type IndexTestInp struct {
	Name string `json:"name" d:"HotGo" dc:"名称"`
}

func (in *IndexTestInp) Filter(ctx context.Context) (err error) {
	return
}

type IndexTestModel struct {
	Name   string      `json:"name" dc:"名称"`
	Module string      `json:"module" dc:"当前插件模块"`
	Time   *gtime.Time `json:"time" dc:"当前时间"`
}
