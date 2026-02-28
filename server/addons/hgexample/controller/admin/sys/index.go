// Package sys
package sys

import (
	"context"

	"hotgo/addons/hgexample/api/admin/index"
	"hotgo/addons/hgexample/service"
)

var (
	Index = cIndex{}
)

type cIndex struct{}

// Test 测试
func (c *cIndex) Test(ctx context.Context, req *index.TestReq) (res *index.TestRes, err error) {
	data, err := service.SysIndex().Test(ctx, &req.IndexTestInp)
	if err != nil {
		return
	}

	res = new(index.TestRes)
	res.IndexTestModel = data
	return
}
