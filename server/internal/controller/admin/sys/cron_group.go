// Package sys
package sys

import (
	"context"

	"hotgo/api/admin/cron"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	CronGroup = cCronGroup{}
)

type cCronGroup struct{}

// Delete 删除
func (c *cCronGroup) Delete(ctx context.Context, req *cron.GroupDeleteReq) (res *cron.GroupDeleteRes, err error) {
	err = sysLogic.SysCronGroup().Delete(ctx, &req.CronGroupDeleteInp)
	return
}

// Edit 更新
func (c *cCronGroup) Edit(ctx context.Context, req *cron.GroupEditReq) (res *cron.GroupEditRes, err error) {
	err = sysLogic.SysCronGroup().Edit(ctx, &req.CronGroupEditInp)
	return
}

// MaxSort 最大排序
func (c *cCronGroup) MaxSort(ctx context.Context, req *cron.GroupMaxSortReq) (res *cron.GroupMaxSortRes, err error) {
	res = new(cron.GroupMaxSortRes)
	res.CronGroupMaxSortModel, err = sysLogic.SysCronGroup().MaxSort(ctx, &req.CronGroupMaxSortInp)
	return
}

// View 获取指定信息
func (c *cCronGroup) View(ctx context.Context, req *cron.GroupViewReq) (res *cron.GroupViewRes, err error) {
	data, err := sysLogic.SysCronGroup().View(ctx, &req.CronGroupViewInp)
	if err != nil {
		return
	}

	res = new(cron.GroupViewRes)
	res.CronGroupViewModel = data
	return
}

// List 查看列表
func (c *cCronGroup) List(ctx context.Context, req *cron.GroupListReq) (res *cron.GroupListRes, err error) {
	list, totalCount, err := sysLogic.SysCronGroup().List(ctx, &req.CronGroupListInp)
	if err != nil {
		return
	}

	res = new(cron.GroupListRes)
	res.List = list
	res.Pack(req, totalCount)
	return
}

// Status 更新状态
func (c *cCronGroup) Status(ctx context.Context, req *cron.GroupStatusReq) (res *cron.GroupStatusRes, err error) {
	err = sysLogic.SysCronGroup().Status(ctx, &req.CronGroupStatusInp)
	return
}

// Select 选项
func (c *cCronGroup) Select(ctx context.Context, req *cron.GroupSelectReq) (res *cron.GroupSelectRes, err error) {
	data, err := sysLogic.SysCronGroup().Select(ctx, &req.CronGroupSelectInp)
	if err != nil {
		return
	}

	res = new(cron.GroupSelectRes)
	res.CronGroupSelectModel = data
	return
}
