// Package gentables
package gentables

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/internal/model/input/sysin"
)

// DbSelectReq 获取可选数据库列表
type DbSelectReq struct {
	g.Meta `path:"/genTable/dbSelect" method:"get" tags:"数据表管理" summary:"获取可选数据库列表"`
}

type DbSelectRes []*sysin.GenTableDbSelectModel

// TableListReq 获取数据库表列表
type TableListReq struct {
	g.Meta `path:"/genTable/tableList" method:"get" tags:"数据表管理" summary:"获取数据库表列表"`
	sysin.GenTableListInp
}

type TableListRes struct {
	List []*sysin.GenTableListModel `json:"list" dc:"数据列表"`
}

// TableViewReq 查看表结构详情
type TableViewReq struct {
	g.Meta `path:"/genTable/tableView" method:"get" tags:"数据表管理" summary:"查看表结构详情"`
	sysin.GenTableViewInp
}

type TableViewRes struct {
	*sysin.GenTableViewModel
}

// TableCreateReq 创建数据表
type TableCreateReq struct {
	g.Meta `path:"/genTable/tableCreate" method:"post" tags:"数据表管理" summary:"创建数据表"`
	sysin.GenTableCreateInp
}

type TableCreateRes struct{}

// TableEditReq 修改表结构
type TableEditReq struct {
	g.Meta `path:"/genTable/tableEdit" method:"post" tags:"数据表管理" summary:"修改表结构"`
	sysin.GenTableEditInp
}

type TableEditRes struct{}

// TableDropReq 删除数据表
type TableDropReq struct {
	g.Meta `path:"/genTable/tableDrop" method:"post" tags:"数据表管理" summary:"删除数据表"`
	sysin.GenTableDropInp
}

type TableDropRes struct{}

// PreviewDDLReq 预览DDL语句
type PreviewDDLReq struct {
	g.Meta `path:"/genTable/previewDDL" method:"post" tags:"数据表管理" summary:"预览DDL语句"`
	sysin.GenTablePreviewDDLInp
}

type PreviewDDLRes struct {
	*sysin.GenTablePreviewDDLModel
}
