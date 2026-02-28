// Package genrouter
package genrouter

import "hotgo/addons/hgexample/controller/admin/sys"

func init() {
	LoginRequiredRouter = append(LoginRequiredRouter, sys.TenantOrder) // 多租户功能演示
}
