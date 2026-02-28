// Package sys
package sys

import (
	"context"
	"hotgo/api/admin/normaltreedemo"
	"hotgo/internal/model/input/sysin"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	NormalTreeDemo = cNormalTreeDemo{}
)

type cNormalTreeDemo struct{}

// List 查看普通树表列表
func (c *cNormalTreeDemo) List(ctx context.Context, req *normaltreedemo.ListReq) (res *normaltreedemo.ListRes, err error) {
	list, totalCount, err := sysLogic.SysNormalTreeDemo().List(ctx, &req.NormalTreeDemoListInp)
	if err != nil {
		return
	}

	if list == nil {
		list = []*sysin.NormalTreeDemoListModel{}
	}

	res = new(normaltreedemo.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Edit 更新普通树表
func (c *cNormalTreeDemo) Edit(ctx context.Context, req *normaltreedemo.EditReq) (res *normaltreedemo.EditRes, err error) {
	err = sysLogic.SysNormalTreeDemo().Edit(ctx, &req.NormalTreeDemoEditInp)
	return
}

// MaxSort 获取普通树表最大排序
func (c *cNormalTreeDemo) MaxSort(ctx context.Context, req *normaltreedemo.MaxSortReq) (res *normaltreedemo.MaxSortRes, err error) {
	data, err := sysLogic.SysNormalTreeDemo().MaxSort(ctx, &req.NormalTreeDemoMaxSortInp)
	if err != nil {
		return
	}

	res = new(normaltreedemo.MaxSortRes)
	res.NormalTreeDemoMaxSortModel = data
	return
}

// View 获取指定普通树表信息
func (c *cNormalTreeDemo) View(ctx context.Context, req *normaltreedemo.ViewReq) (res *normaltreedemo.ViewRes, err error) {
	data, err := sysLogic.SysNormalTreeDemo().View(ctx, &req.NormalTreeDemoViewInp)
	if err != nil {
		return
	}

	res = new(normaltreedemo.ViewRes)
	res.NormalTreeDemoViewModel = data
	return
}

// Delete 删除普通树表
func (c *cNormalTreeDemo) Delete(ctx context.Context, req *normaltreedemo.DeleteReq) (res *normaltreedemo.DeleteRes, err error) {
	err = sysLogic.SysNormalTreeDemo().Delete(ctx, &req.NormalTreeDemoDeleteInp)
	return
}

// TreeOption 获取普通树表关系树选项
func (c *cNormalTreeDemo) TreeOption(ctx context.Context, req *normaltreedemo.TreeOptionReq) (res *normaltreedemo.TreeOptionRes, err error) {
	data, err := sysLogic.SysNormalTreeDemo().TreeOption(ctx)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		res = (*normaltreedemo.TreeOptionRes)(&data)
	} else {
		temp := make(normaltreedemo.TreeOptionRes, 0)
		res = &temp
	}
	return
}
