// Package sys
package sys

import (
	"context"

	"hotgo/api/admin/attachment"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	Attachment = cAttachment{}
)

type cAttachment struct{}

// Delete 删除附件
func (c *cAttachment) Delete(ctx context.Context, req *attachment.DeleteReq) (res *attachment.DeleteRes, err error) {
	err = sysLogic.SysAttachment().Delete(ctx, &req.AttachmentDeleteInp)
	return
}

// View 获取指定附件信息
func (c *cAttachment) View(ctx context.Context, req *attachment.ViewReq) (res *attachment.ViewRes, err error) {
	data, err := sysLogic.SysAttachment().View(ctx, &req.AttachmentViewInp)
	if err != nil {
		return
	}

	res = new(attachment.ViewRes)
	res.AttachmentViewModel = data
	return
}

// List 查看附件列表
func (c *cAttachment) List(ctx context.Context, req *attachment.ListReq) (res *attachment.ListRes, err error) {
	list, totalCount, err := sysLogic.SysAttachment().List(ctx, &req.AttachmentListInp)
	if err != nil {
		return
	}

	res = new(attachment.ListRes)
	res.List = list
	res.Pack(req, totalCount)
	return
}

// ClearKind 清空上传类型
func (c *cAttachment) ClearKind(ctx context.Context, req *attachment.ClearKindReq) (res *attachment.ClearKindRes, err error) {
	err = sysLogic.SysAttachment().ClearKind(ctx, &req.AttachmentClearKindInp)
	return
}
