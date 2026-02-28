// Package sms
package sms

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"

	"hotgo/internal/dao"
	"hotgo/internal/model"
)

var config *model.SmsConfig

func SetConfig(c *model.SmsConfig) {
	config = c
}

func GetConfig() *model.SmsConfig {
	return config
}

func GetModel(ctx context.Context) *gdb.Model {
	return dao.SysSmsLog.Ctx(ctx)
}
