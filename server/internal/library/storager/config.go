// Package storager
package storager

import (
	"context"
	"hotgo/internal/dao"
	"hotgo/internal/model"

	"github.com/gogf/gf/v2/database/gdb"
)

var config *model.UploadConfig

func SetConfig(c *model.UploadConfig) {
	config = c
}

func GetConfig() *model.UploadConfig {
	return config
}

func GetModel(ctx context.Context) *gdb.Model {
	return dao.SysAttachment.Ctx(ctx)
}
