// Package genrouter
package genrouter

import "hotgo/internal/controller/admin/sys"

func init() {
	LoginRequiredRouter = append(LoginRequiredRouter, sys.NormalTreeDemo) // 普通树表
}
