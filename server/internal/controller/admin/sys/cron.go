// Package sys
package sys

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/api/admin/cron"
	"hotgo/api/servmsg"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	Cron = cCron{}
)

type cCron struct{}

// Delete 删除
func (c *cCron) Delete(ctx context.Context, req *cron.DeleteReq) (res *cron.DeleteRes, err error) {
	err = service.TCPServer().CronDelete(ctx, &servmsg.CronDeleteReq{CronDeleteInp: &req.CronDeleteInp})
	return
}

// Edit 更新
func (c *cCron) Edit(ctx context.Context, req *cron.EditReq) (res *cron.EditRes, err error) {
	err = service.TCPServer().CronEdit(ctx, &servmsg.CronEditReq{CronEditInp: &req.CronEditInp})
	return
}

// MaxSort 最大排序
func (c *cCron) MaxSort(ctx context.Context, req *cron.MaxSortReq) (res *cron.MaxSortRes, err error) {
	res = new(cron.MaxSortRes)
	res.CronMaxSortModel, err = sysLogic.SysCron().MaxSort(ctx, &req.CronMaxSortInp)
	return
}

// View 获取指定信息
func (c *cCron) View(ctx context.Context, req *cron.ViewReq) (res *cron.ViewRes, err error) {
	data, err := sysLogic.SysCron().View(ctx, &req.CronViewInp)
	if err != nil {
		return
	}

	res = new(cron.ViewRes)
	res.CronViewModel = data
	return
}

// List 查看列表
func (c *cCron) List(ctx context.Context, req *cron.ListReq) (res *cron.ListRes, err error) {
	list, totalCount, err := sysLogic.SysCron().List(ctx, &req.CronListInp)
	if err != nil {
		return
	}

	res = new(cron.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Status 更新状态
func (c *cCron) Status(ctx context.Context, req *cron.StatusReq) (res *cron.StatusRes, err error) {
	if req.Id <= 0 {
		return nil, gerror.New("定时任务ID不能为空")
	}

	err = service.TCPServer().CronStatus(ctx, &servmsg.CronStatusReq{CronStatusInp: &req.CronStatusInp})
	return
}

// OnlineExec 在线执行
func (c *cCron) OnlineExec(ctx context.Context, req *cron.OnlineExecReq) (res *cron.OnlineExecRes, err error) {
	if req.Id <= 0 {
		return nil, gerror.New("定时任务ID不能为空")
	}

	err = service.TCPServer().CronOnlineExec(ctx, &servmsg.CronOnlineExecReq{OnlineExecInp: &req.OnlineExecInp})
	return
}

// DispatchLog 调度日志
func (c *cCron) DispatchLog(ctx context.Context, req *cron.DispatchLogReq) (res *cron.DispatchLogRes, err error) {
	if req.Id <= 0 {
		return nil, gerror.New("定时任务ID不能为空")
	}

	res = new(cron.DispatchLogRes)
	res.DispatchLogModel = new(sysin.DispatchLogModel)
	res.Log, err = service.TCPServer().DispatchLog(ctx, &servmsg.CronDispatchLogReq{DispatchLogInp: &req.DispatchLogInp})
	return
}
