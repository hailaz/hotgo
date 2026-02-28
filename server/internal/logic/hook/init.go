// Package hook
package hook

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"hotgo/internal/service"
)

type sHook struct {
}

func init() {
	service.RegisterHook(New())
}

func New() *sHook {
	return &sHook{}
}

func (s *sHook) BeforeServe(r *ghttp.Request) {

}

func (s *sHook) AfterOutput(r *ghttp.Request) {
	s.accessLog(r)
	s.lastAdminActive(r)
}
