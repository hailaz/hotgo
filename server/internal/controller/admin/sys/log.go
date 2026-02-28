// Package sys
package sys

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/api/admin/log"
	sysLogic "hotgo/internal/logic/sys"
)

// Log 日志
var Log = sLog{}

type sLog struct{}

// Clear 清空日志
func (c *sLog) Clear(ctx context.Context, _ *log.ClearReq) (res *log.ClearRes, err error) {
	err = gerror.New("暂时考虑到安全问题，请到数据库清空")
	return
}

// Export 导出
func (c *sLog) Export(ctx context.Context, req *log.ExportReq) (res *log.ExportRes, err error) {
	err = sysLogic.SysLog().Export(ctx, &req.LogListInp)
	return
}

// List 获取访问日志列表
func (c *sLog) List(ctx context.Context, req *log.ListReq) (res *log.ListRes, err error) {
	list, totalCount, err := sysLogic.SysLog().List(ctx, &req.LogListInp)
	if err != nil {
		return
	}

	res = new(log.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// View 获取指定信息
func (c *sLog) View(ctx context.Context, req *log.ViewReq) (res *log.ViewRes, err error) {
	res = new(log.ViewRes)
	res.LogViewModel, err = sysLogic.SysLog().View(ctx, &req.LogViewInp)
	return
}

// Delete 删除
func (c *sLog) Delete(ctx context.Context, req *log.DeleteReq) (res *log.DeleteRes, err error) {
	err = sysLogic.SysLog().Delete(ctx, &req.LogDeleteInp)
	return
}
