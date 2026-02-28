// Package config
package config

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/hgexample/model/input/sysin"
)

// GetReq 获取指定分组的配置
type GetReq struct {
	g.Meta `path:"/config/get" method:"get" tags:"配置" summary:"获取指定分组的配置"`
	sysin.GetConfigInp
}

type GetRes struct {
	*sysin.GetConfigModel
}

// UpdateReq 获取指定分组的配置
type UpdateReq struct {
	g.Meta `path:"/config/update" method:"post" tags:"配置" summary:"获取指定分组的配置"`
	sysin.UpdateConfigInp
}

type UpdateRes struct {
}
