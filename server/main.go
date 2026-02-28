// Package main
package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"

	_ "hotgo/addons/modules"
	"hotgo/internal/cmd"
	"hotgo/internal/global"
	_ "hotgo/internal/logic"
	_ "hotgo/internal/packed"
)

func main() {
	var ctx = gctx.GetInitCtx()
	global.Init(ctx)
	cmd.Main.Run(ctx)
}
