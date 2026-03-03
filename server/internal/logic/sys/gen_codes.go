// Package sys
package sys

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/hggen"
	"hotgo/internal/model"
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"
	"hotgo/utility/validate"
)

type sSysGenCodes struct{}

func NewSysGenCodes() *sSysGenCodes {
	return &sSysGenCodes{}
}

var insSysGenCodes = NewSysGenCodes()

func SysGenCodes() *sSysGenCodes {
	return insSysGenCodes
}

// Delete 删除
func (s *sSysGenCodes) Delete(ctx context.Context, in *sysin.GenCodesDeleteInp) (err error) {
	_, err = dao.SysGenCodes.Ctx(ctx).Where("id", in.Id).Delete()
	return
}

// Edit 修改/新增
func (s *sSysGenCodes) Edit(ctx context.Context, in *sysin.GenCodesEditInp) (res *sysin.GenCodesEditModel, err error) {
	if in.GenType == 0 {
		err = gerror.New("生成类型不能为空")
		return
	}
	if in.VarName == "" {
		err = gerror.New("实体名称不能为空")
		return
	}

	if in.GenType == consts.GenCodesTypeCurd {
		var temp *model.GenerateAppCrudTemplate
		cfg := fmt.Sprintf("hggen.application.crud.templates.%v", in.GenTemplate)
		if err = g.Cfg().MustGet(ctx, cfg).Scan(&temp); err != nil {
			return
		}

		if temp == nil {
			err = gerror.Newf("选择的模板不存在:%v", cfg)
			return
		}

		if temp.IsAddon && in.AddonName == "" {
			err = gerror.New("插件模板必须选择一个有效的插件")
			return
		}
	}

	// 修改
	in.UpdatedAt = gtime.Now()
	if in.Id > 0 {
		_, err = dao.SysGenCodes.Ctx(ctx).Where("id", in.Id).Data(in).Update()
		if err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
			return
		}
		return &sysin.GenCodesEditModel{SysGenCodes: in.SysGenCodes}, nil
	}

	// 新增
	if in.Options.IsNil() {
		in.Options = gjson.New(consts.NilJsonToString)
	}

	in.MasterColumns = gjson.New("[]")
	in.Status = consts.GenCodesStatusWait
	in.CreatedAt = gtime.Now()
	id, err := dao.SysGenCodes.Ctx(ctx).Data(in).OmitEmptyData().InsertAndGetId()
	if err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	in.Id = id
	return &sysin.GenCodesEditModel{SysGenCodes: in.SysGenCodes}, nil
}

// Status 更新部门状态
func (s *sSysGenCodes) Status(ctx context.Context, in *sysin.GenCodesStatusInp) (err error) {
	if in.Id <= 0 {
		err = gerror.New("ID不能为空")
		return
	}

	if in.Status <= 0 {
		err = gerror.New("状态不能为空")
		return
	}

	if !validate.InSlice(consts.StatusSlice, in.Status) {
		err = gerror.New("状态不正确")
		return
	}

	_, err = dao.SysGenCodes.Ctx(ctx).Where("id", in.Id).Data("status", in.Status).Update()
	return
}

// MaxSort 最大排序
func (s *sSysGenCodes) MaxSort(ctx context.Context, in *sysin.GenCodesMaxSortInp) (res *sysin.GenCodesMaxSortModel, err error) {
	if in.Id > 0 {
		if err = dao.SysGenCodes.Ctx(ctx).Where("id", in.Id).Order("sort desc").Scan(&res); err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
			return nil, err
		}
	}

	if res == nil {
		res = new(sysin.GenCodesMaxSortModel)
	}

	res.Sort = form.DefaultMaxSort(res.Sort)
	return
}

// View 获取指定字典类型信息
func (s *sSysGenCodes) View(ctx context.Context, in *sysin.GenCodesViewInp) (res *sysin.GenCodesViewModel, err error) {
	err = dao.SysGenCodes.Ctx(ctx).Where("id", in.Id).Scan(&res)
	return
}

// List 获取列表
func (s *sSysGenCodes) List(ctx context.Context, in *sysin.GenCodesListInp) (list []*sysin.GenCodesListModel, totalCount int, err error) {
	mod := dao.SysGenCodes.Ctx(ctx)

	if in.GenType > 0 {
		mod = mod.Where("gen_type", in.GenType)
	}

	if in.VarName != "" {
		mod = mod.Where("var_name", in.VarName)
	}

	if in.Status > 0 {
		mod = mod.Where("status", in.Status)
	}

	totalCount, err = mod.Count()
	if err != nil || totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).Order("id desc").Scan(&list); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	typeSelect, err := hggen.GenTypeSelect(ctx)
	if err != nil {
		return
	}

	getTpGroup := func(row *sysin.GenCodesListModel) string {
		if row == nil {
			return ""
		}

		genType := int(row.GenType)
		if genType == consts.GenCodesTypeTree {
			genType = consts.GenCodesTypeCurd
		}

		for _, v := range typeSelect {
			if v.Value == genType {
				for index, template := range v.Templates {
					if index == row.GenTemplate {
						return template.Label
					}
				}
			}
		}
		return ""
	}

	if len(list) > 0 {
		for _, v := range list {
			v.GenTemplateGroup = getTpGroup(v)
		}
	}
	return
}

// Selects 选项
func (s *sSysGenCodes) Selects(ctx context.Context, in *sysin.GenCodesSelectsInp) (res *sysin.GenCodesSelectsModel, err error) {
	return hggen.TableSelects(ctx, in)
}

// TableSelect 表选项
func (s *sSysGenCodes) TableSelect(ctx context.Context, in *sysin.GenCodesTableSelectInp) (res []*sysin.GenCodesTableSelectModel, err error) {
	var (
		sql           string
		config        = g.DB(in.Name).GetConfig()
		disableTables = g.Cfg().MustGet(ctx, "hggen.disableTables").Strings()
		lists         []*sysin.GenCodesTableSelectModel
	)

	// 根据数据库类型使用不同的SQL
	if config.Type == consts.DBPgsql {
		// PostgreSQL: 使用pg_catalog查询表和注释
		sql = `
			SELECT 
				c.relname as value,
				COALESCE(obj_description(c.oid), '') as label
			FROM pg_class c
			JOIN pg_namespace n ON c.relnamespace = n.oid
			WHERE n.nspname = 'public' 
				AND c.relkind = 'r'
			ORDER BY c.relname`
	} else {
		// MySQL: 使用information_schema.TABLES
		sql = fmt.Sprintf("SELECT TABLE_NAME as value, TABLE_COMMENT as label FROM information_schema.TABLES WHERE TABLE_SCHEMA = '%s'", config.Name)
	}

	if err = g.DB(in.Name).Ctx(ctx).Raw(sql).Scan(&lists); err != nil {
		return
	}

	patternStr := `addon_(\w+)_`
	repStr := ``

	for _, v := range lists {
		if gstr.InArray(disableTables, v.Value) {
			continue
		}

		newValue := v.Value
		if config.Prefix != "" {
			newValue = gstr.SubStrFromEx(v.Value, config.Prefix)
		}
		if newValue == "" {
			err = gerror.Newf("表名[%v]前缀必须和配置中的前缀设置[%v] 保持一致", v.Value, config.Prefix)
			return
		}

		// 如果是插件模块，则移除掉插件表前缀
		bt, err := gregex.Replace(patternStr, []byte(repStr), []byte(newValue))
		if err != nil {
			err = gerror.Newf("表名[%v] gregex.Replace err:%v", v.Value, err.Error())
			return nil, err
		}

		row := v
		row.DefTableComment = v.Label
		row.DaoName = gstr.CaseCamel(newValue)
		row.DefVarName = gstr.CaseCamel(string(bt))
		row.DefAlias = gstr.CaseCamelLower(newValue)
		row.Name = fmt.Sprintf("%s (%s)", v.Value, v.Label)
		row.Label = row.Name
		res = append(res, row)
	}
	return
}

// ColumnSelect 表字段选项
func (s *sSysGenCodes) ColumnSelect(ctx context.Context, in *sysin.GenCodesColumnSelectInp) (res []*sysin.GenCodesColumnSelectModel, err error) {
	var (
		sql    string
		config = g.DB(in.Name).GetConfig()
	)

	// 根据数据库类型使用不同的SQL
	if config.Type == consts.DBPgsql {
		// PostgreSQL: 使用pg_catalog查询列注释
		sql = `
			SELECT 
				a.attname as value,
				COALESCE(col_description(a.attrelid, a.attnum), '') as label
			FROM pg_attribute a
			JOIN pg_class c ON a.attrelid = c.oid
			JOIN pg_namespace n ON c.relnamespace = n.oid
			WHERE n.nspname = 'public' 
				AND c.relname = '%s'
				AND a.attnum > 0 
				AND NOT a.attisdropped
			ORDER BY a.attnum`
		sql = fmt.Sprintf(sql, in.Table)
	} else {
		// MySQL: 使用information_schema.COLUMNS
		sql = fmt.Sprintf("SELECT COLUMN_NAME as value, COLUMN_COMMENT as label FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'", config.Name, in.Table)
	}

	if err = g.DB(in.Name).Ctx(ctx).Raw(sql).Scan(&res); err != nil {
		return
	}

	if len(res) == 0 {
		err = gerror.Newf("table does not exist:%v", in.Table)
		return
	}

	for k, v := range res {
		res[k].Name = fmt.Sprintf("%s (%s)", v.Value, v.Label)
		res[k].Label = res[k].Name
	}
	return
}

// ColumnList 表字段列表
func (s *sSysGenCodes) ColumnList(ctx context.Context, in *sysin.GenCodesColumnListInp) (res []*sysin.GenCodesColumnListModel, err error) {
	return hggen.TableColumns(ctx, in)
}

// Preview 生成预览
func (s *sSysGenCodes) Preview(ctx context.Context, in *sysin.GenCodesPreviewInp) (res *sysin.GenCodesPreviewModel, err error) {
	res, err = hggen.Preview(ctx, in)
	if err != nil {
		return
	}
	// 追加DAO相关文件
	hggen.AppendDaoFiles(res, in.DbName, in.DaoName)
	return
}

// Build 提交生成
func (s *sSysGenCodes) Build(ctx context.Context, in *sysin.GenCodesBuildInp) (err error) {
	// 先保存配置
	ein := in.SysGenCodes
	if _, err = s.Edit(ctx, &sysin.GenCodesEditInp{SysGenCodes: ein}); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return err
	}

	if err = s.Status(ctx, &sysin.GenCodesStatusInp{Id: in.Id, Status: consts.GenCodesStatusOk}); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return err
	}

	if err = hggen.Build(ctx, in); err != nil {
		_ = s.Status(ctx, &sysin.GenCodesStatusInp{Id: in.Id, Status: consts.GenCodesStatusFail})
		return err
	}
	return
}

// Clean 清除生成的代码文件
func (s *sSysGenCodes) Clean(ctx context.Context, in *sysin.GenCodesCleanInp) (res *sysin.GenCodesCleanModel, err error) {
	res = new(sysin.GenCodesCleanModel)

	if len(in.Files) == 0 {
		err = gerror.New("请选择需要删除的文件")
		return
	}

	// 使用项目根目录（server的上级目录）作为安全校验基准，因为web前端文件在../web/下
	projectRoot := filepath.Dir(gfile.Pwd())

	for _, filePath := range in.Files {
		// 路径安全校验：必须是绝对路径
		if !filepath.IsAbs(filePath) {
			res.Failed = append(res.Failed, filePath)
			continue
		}

		// 路径安全校验：解析后必须在项目根目录下
		absPath := filepath.Clean(filePath)
		if !strings.HasPrefix(absPath, projectRoot) {
			res.Failed = append(res.Failed, filePath)
			continue
		}

		// 不允许路径中包含..
		if strings.Contains(filePath, "..") {
			res.Failed = append(res.Failed, filePath)
			continue
		}

		if !gfile.Exists(filePath) {
			continue
		}

		if gfile.IsDir(filePath) {
			res.Failed = append(res.Failed, filePath)
			continue
		}

		if err2 := gfile.Remove(filePath); err2 != nil {
			g.Log().Warningf(ctx, "clean generated file failed, path:%s, err:%v", filePath, err2)
			res.Failed = append(res.Failed, filePath)
			continue
		}

		res.Count++
	}

	// 清除对应菜单权限
	if in.CleanMenu {
		menuCount, menuErr := s.cleanGeneratedMenus(ctx, in.Id)
		if menuErr != nil {
			g.Log().Warningf(ctx, "clean generated menus failed, id:%d, err:%v", in.Id, menuErr)
		} else {
			res.MenuCount = menuCount
		}
	}

	return
}

// cleanGeneratedMenus 根据生成代码ID清除对应的菜单权限
func (s *sSysGenCodes) cleanGeneratedMenus(ctx context.Context, id int64) (count int, err error) {
	// 获取生成代码配置
	var genCode *sysin.GenCodesViewModel
	if err = dao.SysGenCodes.Ctx(ctx).Where("id", id).Scan(&genCode); err != nil {
		return 0, gerror.Wrap(err, "获取生成代码配置失败")
	}
	if genCode == nil {
		return 0, gerror.New("生成代码配置不存在")
	}

	varName := genCode.VarName
	if varName == "" {
		return 0, gerror.New("实体名称为空，无法定位菜单")
	}

	// 根据VarName推导菜单目录的name（与生成时一致：首字母小写的VarName）
	menuDirName := gstr.LcFirst(varName)

	// 查找目录菜单
	dirMenu, err := dao.AdminMenu.Ctx(ctx).
		Fields(dao.AdminMenu.Columns().Id).
		Where(dao.AdminMenu.Columns().Name, menuDirName).
		One()
	if err != nil {
		return 0, gerror.Wrap(err, "查询菜单目录失败")
	}
	if dirMenu.IsEmpty() {
		return 0, nil
	}

	dirId := dirMenu["id"].Int64()

	// 递归查找并删除该目录下的所有子菜单
	count, err = s.deleteMenuAndChildren(ctx, dirId)
	return
}

// deleteMenuAndChildren 递归删除菜单及其所有子菜单，返回删除数量
func (s *sSysGenCodes) deleteMenuAndChildren(ctx context.Context, id int64) (count int, err error) {
	// 查找所有子菜单
	childIds, err := dao.AdminMenu.Ctx(ctx).
		Fields(dao.AdminMenu.Columns().Id).
		Where(dao.AdminMenu.Columns().Pid, id).
		Array()
	if err != nil {
		return 0, gerror.Wrap(err, "查询子菜单失败")
	}

	// 先递归删除子菜单
	for _, childId := range childIds {
		childCount, err2 := s.deleteMenuAndChildren(ctx, childId.Int64())
		if err2 != nil {
			return count, err2
		}
		count += childCount
	}

	// 删除当前菜单
	_, err = dao.AdminMenu.Ctx(ctx).Where(dao.AdminMenu.Columns().Id, id).Delete()
	if err != nil {
		return count, gerror.Wrap(err, "删除菜单失败")
	}
	count++
	return
}
