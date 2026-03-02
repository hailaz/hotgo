// Package admin
package admin

import (
	"context"

	"hotgo/api/admin/menu"
	adminLogic "hotgo/internal/logic/admin"
)

// Menu 菜单
var (
	Menu = cMenu{}
)

type cMenu struct{}

// Delete 删除
func (c *cMenu) Delete(ctx context.Context, req *menu.DeleteReq) (res *menu.DeleteRes, err error) {
	err = adminLogic.AdminMenu().Delete(ctx, &req.MenuDeleteInp)
	return
}

// BatchDelete 批量删除（含子菜单）
func (c *cMenu) BatchDelete(ctx context.Context, req *menu.BatchDeleteReq) (res *menu.BatchDeleteRes, err error) {
	err = adminLogic.AdminMenu().BatchDelete(ctx, &req.MenuBatchDeleteInp)
	return
}

// Edit 更新
func (c *cMenu) Edit(ctx context.Context, req *menu.EditReq) (res *menu.EditRes, err error) {
	err = adminLogic.AdminMenu().Edit(ctx, &req.MenuEditInp)
	return
}

// List 获取列表
func (c *cMenu) List(ctx context.Context, req *menu.ListReq) (res menu.ListRes, err error) {
	res.MenuListModel, err = adminLogic.AdminMenu().List(ctx, &req.MenuListInp)
	return
}
