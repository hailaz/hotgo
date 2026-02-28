// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/api/admin/gentables"
	sysLogic "hotgo/internal/logic/sys"
)

var (
	GenTables = cGenTables{}
)

type cGenTables struct{}

// DbSelect 获取可选数据库列表
func (c *cGenTables) DbSelect(ctx context.Context, req *gentables.DbSelectReq) (res *gentables.DbSelectRes, err error) {
	data, err := sysLogic.SysGenTables().DbSelect(ctx)
	if err != nil {
		return
	}
	res = (*gentables.DbSelectRes)(&data)
	return
}

// TableList 获取数据库表列表
func (c *cGenTables) TableList(ctx context.Context, req *gentables.TableListReq) (res *gentables.TableListRes, err error) {
	data, err := sysLogic.SysGenTables().TableList(ctx, &req.GenTableListInp)
	if err != nil {
		return
	}
	res = &gentables.TableListRes{
		List: data,
	}
	return
}

// TableView 查看表结构详情
func (c *cGenTables) TableView(ctx context.Context, req *gentables.TableViewReq) (res *gentables.TableViewRes, err error) {
	data, err := sysLogic.SysGenTables().TableView(ctx, &req.GenTableViewInp)
	if err != nil {
		return
	}
	res = &gentables.TableViewRes{
		GenTableViewModel: data,
	}
	return
}

// TableCreate 创建数据表
func (c *cGenTables) TableCreate(ctx context.Context, req *gentables.TableCreateReq) (res *gentables.TableCreateRes, err error) {
	err = sysLogic.SysGenTables().TableCreate(ctx, &req.GenTableCreateInp)
	return
}

// TableEdit 修改表结构
func (c *cGenTables) TableEdit(ctx context.Context, req *gentables.TableEditReq) (res *gentables.TableEditRes, err error) {
	err = sysLogic.SysGenTables().TableEdit(ctx, &req.GenTableEditInp)
	return
}

// TableDrop 删除数据表
func (c *cGenTables) TableDrop(ctx context.Context, req *gentables.TableDropReq) (res *gentables.TableDropRes, err error) {
	err = sysLogic.SysGenTables().TableDrop(ctx, &req.GenTableDropInp)
	return
}

// PreviewDDL 预览DDL语句
func (c *cGenTables) PreviewDDL(ctx context.Context, req *gentables.PreviewDDLReq) (res *gentables.PreviewDDLRes, err error) {
	data, err := sysLogic.SysGenTables().PreviewDDL(ctx, &req.GenTablePreviewDDLInp)
	if err != nil {
		return
	}
	res = &gentables.PreviewDDLRes{
		GenTablePreviewDDLModel: data,
	}
	return
}
