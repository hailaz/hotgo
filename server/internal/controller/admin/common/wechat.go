// Package common
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package common

import (
	"context"
	"hotgo/api/admin/common"
	commonLogic "hotgo/internal/logic/common"
)

var (
	Wechat = cWechat{}
)

type cWechat struct{}

func (c *cWechat) Authorize(ctx context.Context, req *common.WechatAuthorizeReq) (res *common.WechatAuthorizeRes, err error) {
	_, err = commonLogic.CommonWechat().Authorize(ctx, &req.WechatAuthorizeInp)
	return
}

func (c *cWechat) AuthorizeCall(ctx context.Context, req *common.WechatAuthorizeCallReq) (res *common.WechatAuthorizeCallRes, err error) {
	_, err = commonLogic.CommonWechat().AuthorizeCall(ctx, &req.WechatAuthorizeCallInp)
	return
}
