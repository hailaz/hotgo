// Package consts
package consts

import (
	"hotgo/internal/library/dict"
	"hotgo/internal/model"
)

func init() {
	dict.RegisterEnums("deptType", "部门类型选项", DeptTypeOptions)
}

const (
	DeptTypeCompany  = "company"  // 公司
	DeptTypeTenant   = "tenant"   // 租户
	DeptTypeMerchant = "merchant" // 商户
	DeptTypeUser     = "user"     // 用户
)

// DeptTypeOptions 部门类型选项
var DeptTypeOptions = []*model.Option{
	dict.GenSuccessOption(DeptTypeCompany, "公司"),
	dict.GenErrorOption(DeptTypeTenant, "租户"),
	dict.GenInfoOption(DeptTypeMerchant, "商户"),
	dict.GenWarningOption(DeptTypeUser, "用户"),
}
