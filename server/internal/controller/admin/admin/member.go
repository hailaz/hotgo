// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/api/admin/member"
	"hotgo/internal/library/contexts"
	"hotgo/internal/model/input/adminin"
	adminLogic "hotgo/internal/logic/admin"
)

var (
	Member = cMember{}
)

type cMember struct{}

// UpdateCash 修改代理商提现信息
func (c *cMember) UpdateCash(ctx context.Context, req *member.UpdateCashReq) (res *member.UpdateCashRes, err error) {
	err = adminLogic.AdminMember().UpdateCash(ctx, &req.MemberUpdateCashInp)
	return
}

// UpdateEmail 换绑邮箱
func (c *cMember) UpdateEmail(ctx context.Context, req *member.UpdateEmailReq) (res *member.UpdateEmailRes, err error) {
	err = adminLogic.AdminMember().UpdateEmail(ctx, &req.MemberUpdateEmailInp)
	return
}

// UpdateMobile 换绑手机号
func (c *cMember) UpdateMobile(ctx context.Context, req *member.UpdateMobileReq) (res *member.UpdateMobileRes, err error) {
	err = adminLogic.AdminMember().UpdateMobile(ctx, &req.MemberUpdateMobileInp)
	return
}

// UpdateProfile 更新用户资料
func (c *cMember) UpdateProfile(ctx context.Context, req *member.UpdateProfileReq) (res *member.UpdateProfileRes, err error) {
	err = adminLogic.AdminMember().UpdateProfile(ctx, &req.MemberUpdateProfileInp)
	return
}

// UpdatePwd 修改登录密码
func (c *cMember) UpdatePwd(ctx context.Context, req *member.UpdatePwdReq) (res *member.UpdatePwdRes, err error) {
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		err = gerror.New("获取用户信息失败！")
		return
	}

	err = adminLogic.AdminMember().UpdatePwd(ctx, &adminin.MemberUpdatePwdInp{
		Id:          memberId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	return
}

// ResetPwd 重置密码
func (c *cMember) ResetPwd(ctx context.Context, req *member.ResetPwdReq) (res *member.ResetPwdRes, err error) {
	err = adminLogic.AdminMember().ResetPwd(ctx, &req.MemberResetPwdInp)
	return
}

// MemberInfo 登录用户信息
func (c *cMember) MemberInfo(ctx context.Context, _ *member.InfoReq) (res *member.InfoRes, err error) {
	data, err := adminLogic.AdminMember().LoginMemberInfo(ctx)
	if err != nil {
		return
	}

	res = new(member.InfoRes)
	res.LoginMemberInfoModel = data
	return
}

// Delete 删除用户
func (c *cMember) Delete(ctx context.Context, req *member.DeleteReq) (res *member.DeleteRes, err error) {
	err = adminLogic.AdminMember().Delete(ctx, &req.MemberDeleteInp)
	return
}

// Edit 修改/新增用户
func (c *cMember) Edit(ctx context.Context, req *member.EditReq) (res *member.EditRes, err error) {
	err = adminLogic.AdminMember().Edit(ctx, &req.MemberEditInp)
	return
}

// View 获取指定用户信息
func (c *cMember) View(ctx context.Context, req *member.ViewReq) (res *member.ViewRes, err error) {
	data, err := adminLogic.AdminMember().View(ctx, &req.MemberViewInp)
	if err != nil {
		return
	}

	res = new(member.ViewRes)
	res.MemberViewModel = data
	return
}

// List 查看用户列表
func (c *cMember) List(ctx context.Context, req *member.ListReq) (res *member.ListRes, err error) {
	list, totalCount, err := adminLogic.AdminMember().List(ctx, &req.MemberListInp)
	if err != nil {
		return
	}

	res = new(member.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Status 更新用户状态
func (c *cMember) Status(ctx context.Context, req *member.StatusReq) (res *member.StatusRes, err error) {
	err = adminLogic.AdminMember().Status(ctx, &req.MemberStatusInp)
	return
}

// Select 获取可选的后台用户选项
func (c *cMember) Select(ctx context.Context, req *member.SelectReq) (res *member.SelectRes, err error) {
	data, err := adminLogic.AdminMember().Select(ctx, &req.MemberSelectInp)
	if err != nil {
		return
	}

	res = (*member.SelectRes)(&data)
	return
}

// AddBalance 增加用户余额
func (c *cMember) AddBalance(ctx context.Context, req *member.AddBalanceReq) (res *member.AddBalanceRes, err error) {
	err = adminLogic.AdminMember().AddBalance(ctx, &req.MemberAddBalanceInp)
	return
}

// AddIntegral 增加用户积分
func (c *cMember) AddIntegral(ctx context.Context, req *member.AddIntegralReq) (res *member.AddIntegralRes, err error) {
	err = adminLogic.AdminMember().AddIntegral(ctx, &req.MemberAddIntegralInp)
	return
}
