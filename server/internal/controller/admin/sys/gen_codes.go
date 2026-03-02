// Package sys
package sys

import (
	"context"

	"hotgo/api/admin/gencodes"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	GenCodes = cGenCodes{}
)

type cGenCodes struct{}

// Delete 删除
func (c *cGenCodes) Delete(ctx context.Context, req *gencodes.DeleteReq) (res *gencodes.DeleteRes, err error) {
	err = sysLogic.SysGenCodes().Delete(ctx, &req.GenCodesDeleteInp)
	return
}

// Edit 更新
func (c *cGenCodes) Edit(ctx context.Context, req *gencodes.EditReq) (res *gencodes.EditRes, err error) {
	data, err := sysLogic.SysGenCodes().Edit(ctx, &req.GenCodesEditInp)
	if err != nil {
		return
	}

	res = new(gencodes.EditRes)
	res.GenCodesEditModel = data
	return
}

// MaxSort 最大排序
func (c *cGenCodes) MaxSort(ctx context.Context, req *gencodes.MaxSortReq) (res *gencodes.MaxSortRes, err error) {
	res = new(gencodes.MaxSortRes)
	res.GenCodesMaxSortModel, err = sysLogic.SysGenCodes().MaxSort(ctx, &req.GenCodesMaxSortInp)
	return
}

// View 获取指定信息
func (c *cGenCodes) View(ctx context.Context, req *gencodes.ViewReq) (res *gencodes.ViewRes, err error) {
	data, err := sysLogic.SysGenCodes().View(ctx, &req.GenCodesViewInp)
	if err != nil {
		return
	}

	res = new(gencodes.ViewRes)
	res.GenCodesViewModel = data
	return
}

// List 查看列表
func (c *cGenCodes) List(ctx context.Context, req *gencodes.ListReq) (res *gencodes.ListRes, err error) {
	list, totalCount, err := sysLogic.SysGenCodes().List(ctx, &req.GenCodesListInp)
	if err != nil {
		return
	}

	res = new(gencodes.ListRes)
	res.List = list
	res.Pack(req, totalCount)
	return
}

// Status 更新状态
func (c *cGenCodes) Status(ctx context.Context, req *gencodes.StatusReq) (res *gencodes.StatusRes, err error) {
	err = sysLogic.SysGenCodes().Status(ctx, &req.GenCodesStatusInp)
	return
}

// Selects 获取指定信息
func (c *cGenCodes) Selects(ctx context.Context, req *gencodes.SelectsReq) (res *gencodes.SelectsRes, err error) {
	data, err := sysLogic.SysGenCodes().Selects(ctx, &req.GenCodesSelectsInp)
	if err != nil {
		return
	}

	res = new(gencodes.SelectsRes)
	res.GenCodesSelectsModel = data
	return
}

// TableSelect 数据库表选项
func (c *cGenCodes) TableSelect(ctx context.Context, req *gencodes.TableSelectReq) (res *gencodes.TableSelectRes, err error) {
	data, err := sysLogic.SysGenCodes().TableSelect(ctx, &req.GenCodesTableSelectInp)
	if err != nil {
		return
	}

	res = (*gencodes.TableSelectRes)(&data)
	return
}

// ColumnSelect 表字段选项
func (c *cGenCodes) ColumnSelect(ctx context.Context, req *gencodes.ColumnSelectReq) (res *gencodes.ColumnSelectRes, err error) {
	data, err := sysLogic.SysGenCodes().ColumnSelect(ctx, &req.GenCodesColumnSelectInp)
	if err != nil {
		return
	}

	res = (*gencodes.ColumnSelectRes)(&data)
	return
}

// ColumnList 表字段列表
func (c *cGenCodes) ColumnList(ctx context.Context, req *gencodes.ColumnListReq) (res *gencodes.ColumnListRes, err error) {
	data, err := sysLogic.SysGenCodes().ColumnList(ctx, &req.GenCodesColumnListInp)
	if err != nil {
		return
	}

	res = (*gencodes.ColumnListRes)(&data)
	return
}

// Preview 生成预览
func (c *cGenCodes) Preview(ctx context.Context, req *gencodes.PreviewReq) (res *gencodes.PreviewRes, err error) {
	data, err := sysLogic.SysGenCodes().Preview(ctx, &req.GenCodesPreviewInp)
	if err != nil {
		return
	}

	res = new(gencodes.PreviewRes)
	res.GenCodesPreviewModel = data
	return
}

// Build 生成预览
func (c *cGenCodes) Build(ctx context.Context, req *gencodes.BuildReq) (res *gencodes.BuildRes, err error) {
	err = sysLogic.SysGenCodes().Build(ctx, &req.GenCodesBuildInp)
	return
}

// Clean 清除生成的代码文件
func (c *cGenCodes) Clean(ctx context.Context, req *gencodes.CleanReq) (res *gencodes.CleanRes, err error) {
	data, err := sysLogic.SysGenCodes().Clean(ctx, &req.GenCodesCleanInp)
	if err != nil {
		return
	}

	res = new(gencodes.CleanRes)
	res.GenCodesCleanModel = data
	return
}
