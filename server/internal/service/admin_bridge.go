// Package service
// 以下接口仅用于打破循环依赖（logic/sys → logic/admin、logic/pay → logic/admin 等场景），
// 不再由 gf gen service 自动生成，改为手动维护。
// 大部分 controller 和非循环调用方应直接 import logic 包。
package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/payin"
)

// ---- AdminMember ----

type IAdminMember interface {
	LoadSuperAdmin(ctx context.Context)
	ClusterSyncSuperAdmin(ctx context.Context, message *gredis.Message)
	VerifySuperId(ctx context.Context, verifyId int64) bool
	GetComplexMemberIds(ctx context.Context, memberIdx, opt string) ([]int64, error)
	GetIdsByKeyword(ctx context.Context, ks string) ([]int64, error)
}

var localAdminMember IAdminMember

func RegisterAdminMember(i IAdminMember) { localAdminMember = i }
func AdminMember() IAdminMember           { return localAdminMember }

// ---- AdminOrder ----

type IAdminOrder interface {
	Model(ctx context.Context, option ...*handler.Option) *gdb.Model
	PayNotify(ctx context.Context, in *payin.NotifyCallFuncInp) error
}

var localAdminOrder IAdminOrder

func RegisterAdminOrder(i IAdminOrder) { localAdminOrder = i }
func AdminOrder() IAdminOrder           { return localAdminOrder }

// ---- AdminMenu ----

type IAdminMenu interface {
	GetFastList(ctx context.Context) (map[int64]*entity.AdminMenu, error)
	Model(ctx context.Context, option ...*handler.Option) *gdb.Model
}

var localAdminMenu IAdminMenu

func RegisterAdminMenu(i IAdminMenu) { localAdminMenu = i }
func AdminMenu() IAdminMenu           { return localAdminMenu }
