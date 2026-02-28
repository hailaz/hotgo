package logic

import (
	// These packages still use init() to register service interfaces
	_ "hotgo/internal/logic/hook"
	_ "hotgo/internal/logic/middleware"
	_ "hotgo/internal/logic/tcpclient"
	_ "hotgo/internal/logic/tcpserver"
	_ "hotgo/internal/logic/view"

	// These packages have init() for dict.RegisterFunc or gcron
	_ "hotgo/internal/logic/admin"
	_ "hotgo/internal/logic/common"
	_ "hotgo/internal/logic/sys"
)
