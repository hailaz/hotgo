package pay

import (
	"context"

	v1 "hotgo/api/api/pay/v1"
	"hotgo/internal/consts"
	"hotgo/internal/library/response"
	"hotgo/internal/model/input/payin"
	payLogic "hotgo/internal/logic/pay"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) NotifyWxPay(ctx context.Context, req *v1.NotifyWxPayReq) (res *v1.NotifyWxPayRes, err error) {
	if _, err = payLogic.Pay().Notify(ctx, &payin.PayNotifyInp{PayType: consts.PayTypeWxPay}); err != nil {
		return
	}

	response.CustomJson(g.RequestFromCtx(ctx), `{"code": "SUCCESS","message": "收单成功"}`)
	return
}
