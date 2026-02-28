// Package genrouter
package genrouter

import "hotgo/internal/controller/admin/sys"

func init() {
	LoginRequiredRouter = append(LoginRequiredRouter, sys.OptionTreeDemo) // 选项树表
}