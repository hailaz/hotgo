// Package admin
// @AutoGenerate Date 2023-04-15 15:59:58
package admin

import (
	"context"

	"hotgo/api/admin/creditslog"
	adminLogic "hotgo/internal/logic/admin"
)

var (
	CreditsLog = cCreditsLog{}
)

type cCreditsLog struct{}

// List 查看资产变动列表
func (c *cCreditsLog) List(ctx context.Context, req *creditslog.ListReq) (res *creditslog.ListRes, err error) {
	list, totalCount, err := adminLogic.AdminCreditsLog().List(ctx, &req.CreditsLogListInp)
	if err != nil {
		return
	}

	res = new(creditslog.ListRes)
	res.List = list
	res.Pack(req, totalCount)
	return
}

// Export 导出资产变动列表
func (c *cCreditsLog) Export(ctx context.Context, req *creditslog.ExportReq) (res *creditslog.ExportRes, err error) {
	err = adminLogic.AdminCreditsLog().Export(ctx, &req.CreditsLogListInp)
	return
}
