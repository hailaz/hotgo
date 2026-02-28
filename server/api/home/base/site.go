// Package base
package base

import "github.com/gogf/gf/v2/frame/g"

type SiteIndexReq struct {
	g.Meta `path:"/" method:"get" summary:"首页" tags:"首页"`
}

type SiteIndexRes struct {
	g.Meta `mime:"text/html" type:"string" example:"<html/>"`
}
