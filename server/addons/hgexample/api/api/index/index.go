// Package index
package index

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/hgexample/model/input/sysin"
)

// TestReq 测试
type TestReq struct {
	g.Meta `path:"/index/test" method:"get" tags:"功能案例" summary:"测试"`
	sysin.IndexTestInp
}

type TestRes struct {
	*sysin.IndexTestModel
}
