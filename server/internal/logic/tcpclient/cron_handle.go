// Package tcpclient
package tcpclient

import (
	"context"
	"hotgo/api/servmsg"
	sysLogic "hotgo/internal/logic/sys"
)

// OnCronDelete 删除任务
func (s *sCronClient) OnCronDelete(ctx context.Context, req *servmsg.CronDeleteReq) (res *servmsg.CronDeleteRes, err error) {
	err = sysLogic.SysCron().Delete(ctx, req.CronDeleteInp)
	return
}

// OnCronEdit 编辑任务
func (s *sCronClient) OnCronEdit(ctx context.Context, req *servmsg.CronEditReq) (res *servmsg.CronEditRes, err error) {
	err = sysLogic.SysCron().Edit(ctx, req.CronEditInp)
	return
}

// OnCronStatus 修改任务状态
func (s *sCronClient) OnCronStatus(ctx context.Context, req *servmsg.CronStatusReq) (res *servmsg.CronStatusRes, err error) {
	err = sysLogic.SysCron().Status(ctx, req.CronStatusInp)
	return
}

// OnCronOnlineExec 执行一次任务
func (s *sCronClient) OnCronOnlineExec(ctx context.Context, req *servmsg.CronOnlineExecReq) (res *servmsg.CronOnlineExecRes, err error) {
	err = sysLogic.SysCron().OnlineExec(ctx, req.OnlineExecInp)
	return
}

// OnCronDispatchLog 查看调度日志
func (s *sCronClient) OnCronDispatchLog(ctx context.Context, req *servmsg.CronDispatchLogReq) (res *servmsg.CronDispatchLogRes, err error) {
	res = new(servmsg.CronDispatchLogRes)
	res.DispatchLogModel, err = sysLogic.SysCron().DispatchLog(ctx, req.DispatchLogInp)
	return
}
