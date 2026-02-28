// Package comp
package comp

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ImportExcelReq 导入Excel
type ImportExcelReq struct {
	g.Meta `path:"/comp/importExcel" method:"post" tags:"功能案例" summary:"导入Excel"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"分片文件"`
}

type ImportExcelSheet struct {
	Sheet string     `json:"sheet"`
	Rows  [][]string `json:"rows"`
}

type ImportExcelRes struct {
	Sheets []*ImportExcelSheet `json:"sheets"`
}
