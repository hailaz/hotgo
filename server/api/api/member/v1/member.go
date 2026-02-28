// Package member
package v1

import "github.com/gogf/gf/v2/frame/g"

// GetIdByCodeReq 通过邀请码获取用户ID
type GetIdByCodeReq struct {
	g.Meta `path:"/member/getIdByCode" method:"post" tags:"用户" summary:"通过邀请码获取用户ID"`
	Code   string `json:"code"   dc:"邀请码"`
}

type GetIdByCodeRes struct{}
