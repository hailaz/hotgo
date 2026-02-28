// Package sys
// @AutoGenerate Date 2023-01-19 16:57:33
package sys

import (
	"context"

	"hotgo/api/admin/loginlog"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	LoginLog = cLoginLog{}
)

type cLoginLog struct{}

// List 查看登录日志列表
func (c *cLoginLog) List(ctx context.Context, req *loginlog.ListReq) (res *loginlog.ListRes, err error) {
	list, totalCount, err := sysLogic.SysLoginLog().List(ctx, &req.LoginLogListInp)
	if err != nil {
		return
	}

	res = new(loginlog.ListRes)
	res.List = list
	res.Pack(req, totalCount)
	return
}

// Export 导出登录日志列表
func (c *cLoginLog) Export(ctx context.Context, req *loginlog.ExportReq) (res *loginlog.ExportRes, err error) {
	err = sysLogic.SysLoginLog().Export(ctx, &req.LoginLogListInp)
	return
}

// Delete 删除登录日志
func (c *cLoginLog) Delete(ctx context.Context, req *loginlog.DeleteReq) (res *loginlog.DeleteRes, err error) {
	err = sysLogic.SysLoginLog().Delete(ctx, &req.LoginLogDeleteInp)
	return
}
