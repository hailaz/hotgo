package pay

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "hotgo/api/api/pay/v1"
	"hotgo/internal/consts"
	payLogic "hotgo/internal/logic/pay"
	"hotgo/internal/model/input/payin"
)

func (c *ControllerV1) NotifyQQPay(ctx context.Context, req *v1.NotifyQQPayReq) (res *v1.NotifyQQPayRes, err error) {
	if _, err = payLogic.Pay().Notify(ctx, &payin.PayNotifyInp{PayType: consts.PayTypeQQPay}); err != nil {
		return
	}

	r := g.RequestFromCtx(ctx)
	r.Response.ClearBuffer()
	r.Response.Write(`<?xml version="1.0" encoding="UTF-8"?>`)
	r.Response.WriteXml(g.Map{
		"return_code": "SUCCESS",
	})
	return
}
