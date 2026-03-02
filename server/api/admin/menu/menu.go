// Package menu
package menu

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/model/input/adminin"
	"hotgo/internal/model/input/form"
)

// EditReq 修改/新增菜单
type EditReq struct {
	g.Meta `path:"/menu/edit" method:"post" tags:"菜单" summary:"修改/新增菜单"`
	adminin.MenuEditInp
}

type EditRes struct{}

// DeleteReq 删除菜单
type DeleteReq struct {
	g.Meta `path:"/menu/delete" method:"post" tags:"菜单" summary:"删除菜单"`
	adminin.MenuDeleteInp
}

type DeleteRes struct{}

// BatchDeleteReq 批量删除菜单（含子菜单）
type BatchDeleteReq struct {
	g.Meta `path:"/menu/batchDelete" method:"post" tags:"菜单" summary:"批量删除菜单（含子菜单）"`
	adminin.MenuBatchDeleteInp
}

type BatchDeleteRes struct{}

// ListReq 获取菜单列表
type ListReq struct {
	g.Meta `path:"/menu/list" method:"get" tags:"菜单" summary:"获取菜单列表"`
	adminin.MenuListInp
}

type ListRes struct {
	*adminin.MenuListModel
	form.PageRes
}
