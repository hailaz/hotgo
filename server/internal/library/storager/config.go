// Package storager
package storager

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"

	"hotgo/internal/dao"
	"hotgo/internal/model"
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
