// Package sys
package sys

import (
	"context"
	"hotgo/api/admin/optiontreedemo"
	"hotgo/internal/model/input/sysin"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	OptionTreeDemo = cOptionTreeDemo{}
)

type cOptionTreeDemo struct{}

// List 查看选项树表列表
func (c *cOptionTreeDemo) List(ctx context.Context, req *optiontreedemo.ListReq) (res *optiontreedemo.ListRes, err error) {
	list, totalCount, err := sysLogic.SysOptionTreeDemo().List(ctx, &req.OptionTreeDemoListInp)
	if err != nil {
		return
	}

	if list == nil {
		list = []*sysin.OptionTreeDemoListModel{}
	}

	res = new(optiontreedemo.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Edit 更新选项树表
func (c *cOptionTreeDemo) Edit(ctx context.Context, req *optiontreedemo.EditReq) (res *optiontreedemo.EditRes, err error) {
	err = sysLogic.SysOptionTreeDemo().Edit(ctx, &req.OptionTreeDemoEditInp)
	return
}

// MaxSort 获取选项树表最大排序
func (c *cOptionTreeDemo) MaxSort(ctx context.Context, req *optiontreedemo.MaxSortReq) (res *optiontreedemo.MaxSortRes, err error) {
	data, err := sysLogic.SysOptionTreeDemo().MaxSort(ctx, &req.OptionTreeDemoMaxSortInp)
	if err != nil {
		return
	}

	res = new(optiontreedemo.MaxSortRes)
	res.OptionTreeDemoMaxSortModel = data
	return
}

// View 获取指定选项树表信息
func (c *cOptionTreeDemo) View(ctx context.Context, req *optiontreedemo.ViewReq) (res *optiontreedemo.ViewRes, err error) {
	data, err := sysLogic.SysOptionTreeDemo().View(ctx, &req.OptionTreeDemoViewInp)
	if err != nil {
		return
	}

	res = new(optiontreedemo.ViewRes)
	res.OptionTreeDemoViewModel = data
	return
}

// Delete 删除选项树表
func (c *cOptionTreeDemo) Delete(ctx context.Context, req *optiontreedemo.DeleteReq) (res *optiontreedemo.DeleteRes, err error) {
	err = sysLogic.SysOptionTreeDemo().Delete(ctx, &req.OptionTreeDemoDeleteInp)
	return
}

// TreeOption 获取选项树表关系树选项
func (c *cOptionTreeDemo) TreeOption(ctx context.Context, req *optiontreedemo.TreeOptionReq) (res *optiontreedemo.TreeOptionRes, err error) {
	data, err := sysLogic.SysOptionTreeDemo().TreeOption(ctx)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		res = (*optiontreedemo.TreeOptionRes)(&data)
	} else {
		temp := make(optiontreedemo.TreeOptionRes, 0)
		res = &temp
	}
	return
}