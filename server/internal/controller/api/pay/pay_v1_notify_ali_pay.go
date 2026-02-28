package pay

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "hotgo/api/api/pay/v1"
	"hotgo/internal/consts"
	"hotgo/internal/library/response"
	payLogic "hotgo/internal/logic/pay"
	"hotgo/internal/model/input/payin"
)

func (c *ControllerV1) NotifyAliPay(ctx context.Context, req *v1.NotifyAliPayReq) (res *v1.NotifyAliPayRes, err error) {
	if _, err = payLogic.Pay().Notify(ctx, &payin.PayNotifyInp{PayType: consts.PayTypeAliPay}); err != nil {
		return nil, err
	}

	response.RText(g.RequestFromCtx(ctx), "success")
	return
}
