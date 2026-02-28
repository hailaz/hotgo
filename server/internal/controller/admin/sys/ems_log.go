// Package sys
package sys

import (
	"context"
	"hotgo/api/admin/emslog"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	EmsLog = cEmsLog{}
)

type cEmsLog struct{}

// Delete 删除
func (c *cEmsLog) Delete(ctx context.Context, req *emslog.DeleteReq) (res *emslog.DeleteRes, err error) {
	err = sysLogic.SysEmsLog().Delete(ctx, &req.EmsLogDeleteInp)
	return
}

// Edit 更新
func (c *cEmsLog) Edit(ctx context.Context, req *emslog.EditReq) (res *emslog.EditRes, err error) {
	err = sysLogic.SysEmsLog().Edit(ctx, &req.EmsLogEditInp)
	return
}

// View 获取指定信息
func (c *cEmsLog) View(ctx context.Context, req *emslog.ViewReq) (res *emslog.ViewRes, err error) {
	data, err := sysLogic.SysEmsLog().View(ctx, &req.EmsLogViewInp)
	if err != nil {
		return
	}

	res = new(emslog.ViewRes)
	res.EmsLogViewModel = data
	return
}

// List 查看列表
func (c *cEmsLog) List(ctx context.Context, req *emslog.ListReq) (res *emslog.ListRes, err error) {
	list, totalCount, err := sysLogic.SysEmsLog().List(ctx, &req.EmsLogListInp)
	if err != nil {
		return
	}

	res = new(emslog.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Status 更新状态
func (c *cEmsLog) Status(ctx context.Context, req *emslog.StatusReq) (res *emslog.StatusRes, err error) {
	err = sysLogic.SysEmsLog().Status(ctx, &req.EmsLogStatusInp)
	return
}
