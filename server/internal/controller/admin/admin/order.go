// Package admin
package admin

import (
	"context"
	"hotgo/api/admin/order"
	adminLogic "hotgo/internal/logic/admin"
)

var (
	Order = cOrder{}
)

type cOrder struct{}

// AcceptRefund 受理申请退款
func (c *cOrder) AcceptRefund(ctx context.Context, req *order.AcceptRefundReq) (res *order.AcceptRefundRes, err error) {
	err = adminLogic.AdminOrder().AcceptRefund(ctx, &req.OrderAcceptRefundInp)
	return
}

// ApplyRefund 申请退款
func (c *cOrder) ApplyRefund(ctx context.Context, req *order.ApplyRefundReq) (res *order.ApplyRefundRes, err error) {
	err = adminLogic.AdminOrder().ApplyRefund(ctx, &req.OrderApplyRefundInp)
	return
}

// Create 创建充值订单
func (c *cOrder) Create(ctx context.Context, req *order.CreateReq) (res *order.CreateRes, err error) {
	data, err := adminLogic.AdminOrder().Create(ctx, &req.OrderCreateInp)
	if err != nil {
		return
	}

	res = new(order.CreateRes)
	res.OrderCreateModel = data
	return
}

// List 查看充值订单列表
func (c *cOrder) List(ctx context.Context, req *order.ListReq) (res *order.ListRes, err error) {
	list, totalCount, err := adminLogic.AdminOrder().List(ctx, &req.OrderListInp)
	if err != nil {
		return
	}

	res = new(order.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Export 导出充值订单列表
func (c *cOrder) Export(ctx context.Context, req *order.ExportReq) (res *order.ExportRes, err error) {
	err = adminLogic.AdminOrder().Export(ctx, &req.OrderListInp)
	return
}

// Edit 更新充值订单
func (c *cOrder) Edit(ctx context.Context, req *order.EditReq) (res *order.EditRes, err error) {
	err = adminLogic.AdminOrder().Edit(ctx, &req.OrderEditInp)
	return
}

// View 获取指定充值订单信息
func (c *cOrder) View(ctx context.Context, req *order.ViewReq) (res *order.ViewRes, err error) {
	data, err := adminLogic.AdminOrder().View(ctx, &req.OrderViewInp)
	if err != nil {
		return
	}

	res = new(order.ViewRes)
	res.OrderViewModel = data
	return
}

// Delete 删除充值订单
func (c *cOrder) Delete(ctx context.Context, req *order.DeleteReq) (res *order.DeleteRes, err error) {
	err = adminLogic.AdminOrder().Delete(ctx, &req.OrderDeleteInp)
	return
}

// Status 更新充值订单状态
func (c *cOrder) Status(ctx context.Context, req *order.StatusReq) (res *order.StatusRes, err error) {
	err = adminLogic.AdminOrder().Status(ctx, &req.OrderStatusInp)
	return
}
