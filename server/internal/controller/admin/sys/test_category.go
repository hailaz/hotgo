// Package sys
package sys

import (
	"context"
	"hotgo/api/admin/testcategory"
	"hotgo/internal/model/input/sysin"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	TestCategory = cTestCategory{}
)

type cTestCategory struct{}

// List 查看测试分类列表
func (c *cTestCategory) List(ctx context.Context, req *testcategory.ListReq) (res *testcategory.ListRes, err error) {
	list, totalCount, err := sysLogic.SysTestCategory().List(ctx, &req.TestCategoryListInp)
	if err != nil {
		return
	}

	if list == nil {
		list = []*sysin.TestCategoryListModel{}
	}

	res = new(testcategory.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Edit 更新测试分类
func (c *cTestCategory) Edit(ctx context.Context, req *testcategory.EditReq) (res *testcategory.EditRes, err error) {
	err = sysLogic.SysTestCategory().Edit(ctx, &req.TestCategoryEditInp)
	return
}

// MaxSort 获取测试分类最大排序
func (c *cTestCategory) MaxSort(ctx context.Context, req *testcategory.MaxSortReq) (res *testcategory.MaxSortRes, err error) {
	data, err := sysLogic.SysTestCategory().MaxSort(ctx, &req.TestCategoryMaxSortInp)
	if err != nil {
		return
	}

	res = new(testcategory.MaxSortRes)
	res.TestCategoryMaxSortModel = data
	return
}

// View 获取指定测试分类信息
func (c *cTestCategory) View(ctx context.Context, req *testcategory.ViewReq) (res *testcategory.ViewRes, err error) {
	data, err := sysLogic.SysTestCategory().View(ctx, &req.TestCategoryViewInp)
	if err != nil {
		return
	}

	res = new(testcategory.ViewRes)
	res.TestCategoryViewModel = data
	return
}

// Delete 删除测试分类
func (c *cTestCategory) Delete(ctx context.Context, req *testcategory.DeleteReq) (res *testcategory.DeleteRes, err error) {
	err = sysLogic.SysTestCategory().Delete(ctx, &req.TestCategoryDeleteInp)
	return
}

// Status 更新测试分类状态
func (c *cTestCategory) Status(ctx context.Context, req *testcategory.StatusReq) (res *testcategory.StatusRes, err error) {
	err = sysLogic.SysTestCategory().Status(ctx, &req.TestCategoryStatusInp)
	return
}
