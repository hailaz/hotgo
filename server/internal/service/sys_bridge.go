// Package service
// 以下接口仅用于打破循环依赖（global → logic/sys、queues → logic/sys、hggen → logic/sys 等场景），
// 不再由 gf gen service 自动生成，改为手动维护。
// 大部分 controller 和非循环调用方应直接 import logic 包。
package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/net/ghttp"

	"hotgo/internal/model"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
)

// ---- SysConfig ----

type ISysConfig interface {
	InitConfig(ctx context.Context)
	ClusterSync(ctx context.Context, message *gredis.Message)
	GetLoadServeLog(ctx context.Context) (conf *model.ServeLogConfig, err error)
	GetLoadGenerate(ctx context.Context) (conf *model.GenerateConfig, err error)
}

var localSysConfig ISysConfig

func RegisterSysConfig(i ISysConfig) { localSysConfig = i }
func SysConfig() ISysConfig          { return localSysConfig }

// ---- SysBlacklist ----

type ISysBlacklist interface {
	Load(ctx context.Context)
	ClusterSync(ctx context.Context, message *gredis.Message)
	VerifyRequest(r *ghttp.Request) error
}

var localSysBlacklist ISysBlacklist

func RegisterSysBlacklist(i ISysBlacklist) { localSysBlacklist = i }
func SysBlacklist() ISysBlacklist          { return localSysBlacklist }

// ---- SysDictType ----

type ISysDictType interface {
	TreeSelect(ctx context.Context, in *sysin.DictTreeSelectInp) ([]*sysin.DictTypeTree, error)
}

var localSysDictType ISysDictType

func RegisterSysDictType(i ISysDictType) { localSysDictType = i }
func SysDictType() ISysDictType          { return localSysDictType }

// ---- SysServeLog ----

type ISysServeLog interface {
	RealWrite(ctx context.Context, models entity.SysServeLog) error
}

var localSysServeLog ISysServeLog

func RegisterSysServeLog(i ISysServeLog) { localSysServeLog = i }
func SysServeLog() ISysServeLog          { return localSysServeLog }

// ---- SysLog ----

type ISysLog interface {
	AutoLog(ctx context.Context) error
	RealWrite(ctx context.Context, log entity.SysLog) error
}

var localSysLog ISysLog

func RegisterSysLog(i ISysLog) { localSysLog = i }
func SysLog() ISysLog          { return localSysLog }

// ---- SysLoginLog ----

type ISysLoginLog interface {
	RealWrite(ctx context.Context, models entity.SysLoginLog) error
}

var localSysLoginLog ISysLoginLog

func RegisterSysLoginLog(i ISysLoginLog) { localSysLoginLog = i }
func SysLoginLog() ISysLoginLog          { return localSysLoginLog }

// ---- SysAddonsConfig ----

type ISysAddonsConfig interface {
	GetConfigByGroup(ctx context.Context, in *sysin.GetAddonsConfigInp) (*sysin.GetAddonsConfigModel, error)
	UpdateConfigByGroup(ctx context.Context, in *sysin.UpdateAddonsConfigInp) error
}

var localSysAddonsConfig ISysAddonsConfig

func RegisterSysAddonsConfig(i ISysAddonsConfig) { localSysAddonsConfig = i }
func SysAddonsConfig() ISysAddonsConfig          { return localSysAddonsConfig }
