// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"fmt"
	"hotgo/internal/consts"
	"hotgo/internal/library/hggen"
	"hotgo/internal/model/input/sysin"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

type sSysGenTables struct{}

func NewSysGenTables() *sSysGenTables {
	return &sSysGenTables{}
}

var insSysGenTables = NewSysGenTables()

func SysGenTables() *sSysGenTables {
	return insSysGenTables
}

// DbSelect 获取可选数据库列表
func (s *sSysGenTables) DbSelect(ctx context.Context) (res []*sysin.GenTableDbSelectModel, err error) {
	selects := hggen.DbSelect(ctx)
	for _, v := range selects {
		res = append(res, &sysin.GenTableDbSelectModel{
			Value: v.Value.(string),
			Label: v.Label,
		})
	}
	if len(res) == 0 {
		res = make([]*sysin.GenTableDbSelectModel, 0)
	}
	return
}

// TableList 获取数据库表列表
func (s *sSysGenTables) TableList(ctx context.Context, in *sysin.GenTableListInp) (res []*sysin.GenTableListModel, err error) {
	dbName := in.DbName
	if dbName == "" {
		dbName = "default"
	}

	config := g.DB(dbName).GetConfig()
	var sql string

	if config.Type == consts.DBPgsql {
		sql = `
			SELECT 
				c.relname as table_name,
				COALESCE(obj_description(c.oid), '') as table_comment,
				'' as engine,
				c.reltuples::bigint as table_rows,
				'' as create_time,
				'' as table_collation
			FROM pg_class c
			JOIN pg_namespace n ON c.relnamespace = n.oid
			WHERE n.nspname = 'public' 
				AND c.relkind = 'r'
			ORDER BY c.relname`
	} else {
		sql = fmt.Sprintf(`
			SELECT 
				TABLE_NAME as table_name,
				TABLE_COMMENT as table_comment,
				ENGINE as engine,
				TABLE_ROWS as table_rows,
				CREATE_TIME as create_time,
				TABLE_COLLATION as table_collation
			FROM information_schema.TABLES 
			WHERE TABLE_SCHEMA = '%s'
			ORDER BY CREATE_TIME DESC`, config.Name)
	}

	if err = g.DB(dbName).Ctx(ctx).Raw(sql).Scan(&res); err != nil {
		err = gerror.Wrap(err, "查询表列表失败")
		return
	}

	if in.TableName != "" {
		filtered := make([]*sysin.GenTableListModel, 0)
		for _, v := range res {
			if gstr.ContainsI(v.TableName, in.TableName) {
				filtered = append(filtered, v)
			}
		}
		res = filtered
	}

	if res == nil {
		res = make([]*sysin.GenTableListModel, 0)
	}
	return
}

// TableView 查看表结构详情
func (s *sSysGenTables) TableView(ctx context.Context, in *sysin.GenTableViewInp) (res *sysin.GenTableViewModel, err error) {
	dbName := in.DbName
	config := g.DB(dbName).GetConfig()

	res = &sysin.GenTableViewModel{
		TableName: in.TableName,
	}

	// 获取表基本信息
	if config.Type == consts.DBPgsql {
		tableInfo, e := g.DB(dbName).Ctx(ctx).Raw(fmt.Sprintf(`
			SELECT 
				c.relname as table_name,
				COALESCE(obj_description(c.oid), '') as table_comment
			FROM pg_class c
			JOIN pg_namespace n ON c.relnamespace = n.oid
			WHERE n.nspname = 'public' AND c.relname = '%s'`, in.TableName)).One()
		if e != nil {
			err = gerror.Wrap(e, "查询表信息失败")
			return
		}
		if tableInfo.IsEmpty() {
			err = gerror.Newf("表 %s 不存在", in.TableName)
			return
		}
		res.TableComment = tableInfo["table_comment"].String()
	} else {
		tableInfo, e := g.DB(dbName).Ctx(ctx).Raw(fmt.Sprintf(`
			SELECT TABLE_NAME, TABLE_COMMENT, ENGINE, TABLE_COLLATION
			FROM information_schema.TABLES 
			WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'`, config.Name, in.TableName)).One()
		if e != nil {
			err = gerror.Wrap(e, "查询表信息失败")
			return
		}
		if tableInfo.IsEmpty() {
			err = gerror.Newf("表 %s 不存在", in.TableName)
			return
		}
		res.TableComment = tableInfo["TABLE_COMMENT"].String()
		res.Engine = tableInfo["ENGINE"].String()
		res.Charset = tableInfo["TABLE_COLLATION"].String()
	}

	// 获取字段信息
	if config.Type == consts.DBPgsql {
		columns, e := s.getPgsqlColumns(ctx, dbName, in.TableName)
		if e != nil {
			err = e
			return
		}
		res.Columns = columns
	} else {
		columns, e := s.getMysqlColumns(ctx, dbName, config.Name, in.TableName)
		if e != nil {
			err = e
			return
		}
		res.Columns = columns
	}

	// 获取索引信息
	if config.Type == consts.DBPgsql {
		indexes, e := s.getPgsqlIndexes(ctx, dbName, in.TableName)
		if e != nil {
			err = e
			return
		}
		res.Indexes = indexes
	} else {
		indexes, e := s.getMysqlIndexes(ctx, dbName, in.TableName)
		if e != nil {
			err = e
			return
		}
		res.Indexes = indexes
	}
	return
}

// getMysqlColumns 获取MySQL表字段信息
func (s *sSysGenTables) getMysqlColumns(ctx context.Context, dbName, schemaName, tableName string) (res []*sysin.GenTableViewColumnModel, err error) {
	type columnInfo struct {
		ColumnName    string `json:"COLUMN_NAME"`
		DataType      string `json:"DATA_TYPE"`
		CharMaxLength *int   `json:"CHARACTER_MAXIMUM_LENGTH"`
		NumPrecision  *int   `json:"NUMERIC_PRECISION"`
		NumScale      *int   `json:"NUMERIC_SCALE"`
		IsNullable    string `json:"IS_NULLABLE"`
		ColumnDefault *string `json:"COLUMN_DEFAULT"`
		ColumnComment string `json:"COLUMN_COMMENT"`
		ColumnKey     string `json:"COLUMN_KEY"`
		Extra         string `json:"EXTRA"`
		ColumnType    string `json:"COLUMN_TYPE"`
	}

	var columns []*columnInfo
	sql := fmt.Sprintf(`
		SELECT 
			COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, 
			NUMERIC_PRECISION, NUMERIC_SCALE, IS_NULLABLE,
			COLUMN_DEFAULT, COLUMN_COMMENT, COLUMN_KEY, EXTRA, COLUMN_TYPE
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'
		ORDER BY ORDINAL_POSITION`, schemaName, tableName)

	if err = g.DB(dbName).Ctx(ctx).Raw(sql).Scan(&columns); err != nil {
		err = gerror.Wrap(err, "查询表字段失败")
		return
	}

	for _, col := range columns {
		item := &sysin.GenTableViewColumnModel{
			Name:         col.ColumnName,
			DataType:     col.DataType,
			IsNullable:   col.IsNullable == "YES",
			Comment:      col.ColumnComment,
			IsPrimaryKey: col.ColumnKey == "PRI",
			IsAutoInc:    gstr.ContainsI(col.Extra, "auto_increment"),
			IsUnsigned:   gstr.ContainsI(col.ColumnType, "unsigned"),
			ColumnType:   col.ColumnType,
		}
		if col.CharMaxLength != nil {
			item.Length = *col.CharMaxLength
		} else if col.NumPrecision != nil {
			item.Length = *col.NumPrecision
		}
		if col.NumScale != nil {
			item.Decimal = *col.NumScale
		}
		if col.ColumnDefault != nil {
			item.DefaultValue = *col.ColumnDefault
		}
		res = append(res, item)
	}
	return
}

// getPgsqlColumns 获取PostgreSQL表字段信息
func (s *sSysGenTables) getPgsqlColumns(ctx context.Context, dbName, tableName string) (res []*sysin.GenTableViewColumnModel, err error) {
	type columnInfo struct {
		ColumnName    string  `json:"column_name"`
		DataType      string  `json:"data_type"`
		CharMaxLength *int    `json:"character_maximum_length"`
		NumPrecision  *int    `json:"numeric_precision"`
		NumScale      *int    `json:"numeric_scale"`
		IsNullable    string  `json:"is_nullable"`
		ColumnDefault *string `json:"column_default"`
		IsPrimaryKey  bool    `json:"is_primary_key"`
		Description   string  `json:"description"`
	}

	var columns []*columnInfo
	sql := fmt.Sprintf(`
		SELECT 
			a.attname as column_name,
			format_type(a.atttypid, a.atttypmod) as data_type,
			CASE WHEN a.atttypmod > 0 AND t.typname IN ('varchar','char','bpchar') 
				THEN a.atttypmod - 4 ELSE NULL END as character_maximum_length,
			CASE WHEN t.typname IN ('numeric','decimal') 
				THEN ((a.atttypmod - 4) >> 16) & 65535 ELSE NULL END as numeric_precision,
			CASE WHEN t.typname IN ('numeric','decimal') 
				THEN (a.atttypmod - 4) & 65535 ELSE NULL END as numeric_scale,
			CASE WHEN a.attnotnull THEN 'NO' ELSE 'YES' END as is_nullable,
			pg_get_expr(d.adbin, d.adrelid) as column_default,
			COALESCE(col_description(a.attrelid, a.attnum), '') as description,
			EXISTS(
				SELECT 1 FROM pg_constraint con 
				WHERE con.conrelid = a.attrelid 
				AND a.attnum = ANY(con.conkey) 
				AND con.contype = 'p'
			) as is_primary_key
		FROM pg_attribute a
		JOIN pg_type t ON a.atttypid = t.oid
		JOIN pg_class c ON a.attrelid = c.oid
		JOIN pg_namespace n ON c.relnamespace = n.oid
		LEFT JOIN pg_attrdef d ON a.attrelid = d.adrelid AND a.attnum = d.adnum
		WHERE n.nspname = 'public' 
			AND c.relname = '%s'
			AND a.attnum > 0 
			AND NOT a.attisdropped
		ORDER BY a.attnum`, tableName)

	if err = g.DB(dbName).Ctx(ctx).Raw(sql).Scan(&columns); err != nil {
		err = gerror.Wrap(err, "查询表字段失败")
		return
	}

	for _, col := range columns {
		item := &sysin.GenTableViewColumnModel{
			Name:         col.ColumnName,
			DataType:     col.DataType,
			IsNullable:   col.IsNullable == "YES",
			Comment:      col.Description,
			IsPrimaryKey: col.IsPrimaryKey,
			IsAutoInc:    col.ColumnDefault != nil && gstr.Contains(*col.ColumnDefault, "nextval"),
		}
		if col.CharMaxLength != nil {
			item.Length = *col.CharMaxLength
		} else if col.NumPrecision != nil {
			item.Length = *col.NumPrecision
		}
		if col.NumScale != nil {
			item.Decimal = *col.NumScale
		}
		if col.ColumnDefault != nil && !gstr.Contains(*col.ColumnDefault, "nextval") {
			item.DefaultValue = *col.ColumnDefault
		}
		res = append(res, item)
	}
	return
}

// getMysqlIndexes 获取MySQL表索引
func (s *sSysGenTables) getMysqlIndexes(ctx context.Context, dbName, tableName string) (res []*sysin.GenTableViewIndexModel, err error) {
	type indexInfo struct {
		KeyName    string `json:"Key_name"`
		ColumnName string `json:"Column_name"`
		NonUnique  int    `json:"Non_unique"`
		IndexType  string `json:"Index_type"`
	}

	var indexes []*indexInfo
	sql := fmt.Sprintf("SHOW INDEX FROM `%s`", tableName)
	if err = g.DB(dbName).Ctx(ctx).Raw(sql).Scan(&indexes); err != nil {
		err = gerror.Wrap(err, "查询表索引失败")
		return
	}

	indexMap := make(map[string]*sysin.GenTableViewIndexModel)
	indexOrder := make([]string, 0)
	for _, idx := range indexes {
		if idx.KeyName == "PRIMARY" {
			continue
		}
		if existing, ok := indexMap[idx.KeyName]; ok {
			existing.Columns = append(existing.Columns, idx.ColumnName)
		} else {
			indexType := "INDEX"
			if idx.NonUnique == 0 {
				indexType = "UNIQUE"
			}
			if idx.IndexType == "FULLTEXT" {
				indexType = "FULLTEXT"
			}
			indexMap[idx.KeyName] = &sysin.GenTableViewIndexModel{
				Name:    idx.KeyName,
				Type:    indexType,
				Columns: []string{idx.ColumnName},
			}
			indexOrder = append(indexOrder, idx.KeyName)
		}
	}

	for _, name := range indexOrder {
		res = append(res, indexMap[name])
	}
	return
}

// getPgsqlIndexes 获取PostgreSQL表索引
func (s *sSysGenTables) getPgsqlIndexes(ctx context.Context, dbName, tableName string) (res []*sysin.GenTableViewIndexModel, err error) {
	type indexInfo struct {
		IndexName  string `json:"indexname"`
		IndexDef   string `json:"indexdef"`
	}

	var indexes []*indexInfo
	sql := fmt.Sprintf(`
		SELECT indexname, indexdef 
		FROM pg_indexes 
		WHERE tablename = '%s' AND schemaname = 'public'
		AND indexname NOT IN (
			SELECT conname FROM pg_constraint WHERE contype = 'p'
		)`, tableName)

	if err = g.DB(dbName).Ctx(ctx).Raw(sql).Scan(&indexes); err != nil {
		err = gerror.Wrap(err, "查询表索引失败")
		return
	}

	for _, idx := range indexes {
		indexType := "INDEX"
		if gstr.ContainsI(idx.IndexDef, "UNIQUE") {
			indexType = "UNIQUE"
		}
		// 从 indexdef 中解析字段名
		cols := parseIndexColumns(idx.IndexDef)
		res = append(res, &sysin.GenTableViewIndexModel{
			Name:    idx.IndexName,
			Type:    indexType,
			Columns: cols,
		})
	}
	return
}

// parseIndexColumns 从 CREATE INDEX 语句中解析字段名
func parseIndexColumns(indexDef string) []string {
	start := strings.LastIndex(indexDef, "(")
	end := strings.LastIndex(indexDef, ")")
	if start < 0 || end < 0 || end <= start {
		return nil
	}
	colStr := indexDef[start+1 : end]
	parts := strings.Split(colStr, ",")
	cols := make([]string, 0, len(parts))
	for _, p := range parts {
		col := strings.TrimSpace(p)
		col = strings.Trim(col, "\"")
		if col != "" {
			cols = append(cols, col)
		}
	}
	return cols
}

// TableCreate 创建数据表
func (s *sSysGenTables) TableCreate(ctx context.Context, in *sysin.GenTableCreateInp) (err error) {
	ddl, err := s.buildCreateDDL(ctx, in.DbName, in.TableName, in.Comment, in.Engine, in.Columns, in.Indexes)
	if err != nil {
		return
	}

	if _, err = g.DB(in.DbName).Ctx(ctx).Exec(ctx, ddl); err != nil {
		err = gerror.Wrapf(err, "执行建表DDL失败")
		return
	}
	return
}

// TableEdit 修改表结构
func (s *sSysGenTables) TableEdit(ctx context.Context, in *sysin.GenTableEditInp) (err error) {
	config := g.DB(in.DbName).GetConfig()

	// 检查表是否存在
	exists, e := s.tableExists(ctx, in.DbName, in.TableName)
	if e != nil {
		err = e
		return
	}
	if !exists {
		err = gerror.Newf("表 %s 不存在", in.TableName)
		return
	}

	// 获取现有字段
	var existingColumns []*sysin.GenTableViewColumnModel
	if config.Type == consts.DBPgsql {
		existingColumns, err = s.getPgsqlColumns(ctx, in.DbName, in.TableName)
	} else {
		existingColumns, err = s.getMysqlColumns(ctx, in.DbName, config.Name, in.TableName)
	}
	if err != nil {
		return
	}

	existingColMap := make(map[string]*sysin.GenTableViewColumnModel)
	for _, col := range existingColumns {
		existingColMap[col.Name] = col
	}

	// 生成ALTER语句
	var alterStatements []string

	// 新增和修改字段
	for i, col := range in.Columns {
		if _, exists := existingColMap[col.Name]; exists {
			// 修改已有字段
			alterSQL := s.buildModifyColumnSQL(config.Type, in.TableName, &col)
			if alterSQL != "" {
				alterStatements = append(alterStatements, alterSQL)
			}
			delete(existingColMap, col.Name)
		} else {
			// 新增字段
			afterCol := ""
			if i > 0 {
				afterCol = in.Columns[i-1].Name
			}
			alterSQL := s.buildAddColumnSQL(config.Type, in.TableName, &col, afterCol)
			alterStatements = append(alterStatements, alterSQL)
		}
	}

	// 删除不再需要的字段
	for colName := range existingColMap {
		alterSQL := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", in.TableName, colName)
		if config.Type == consts.DBPgsql {
			alterSQL = fmt.Sprintf(`ALTER TABLE "%s" DROP COLUMN "%s"`, in.TableName, colName)
		}
		alterStatements = append(alterStatements, alterSQL)
	}

	// 修改表注释
	if in.Comment != "" {
		if config.Type == consts.DBPgsql {
			alterStatements = append(alterStatements, fmt.Sprintf(`COMMENT ON TABLE "%s" IS '%s'`, in.TableName, escapeString(in.Comment)))
		} else {
			alterStatements = append(alterStatements, fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", in.TableName, escapeString(in.Comment)))
		}
	}

	// 处理索引：先删除现有非主键索引，再重新创建
	if len(in.Indexes) > 0 {
		var existingIndexes []*sysin.GenTableViewIndexModel
		if config.Type == consts.DBPgsql {
			existingIndexes, err = s.getPgsqlIndexes(ctx, in.DbName, in.TableName)
		} else {
			existingIndexes, err = s.getMysqlIndexes(ctx, in.DbName, in.TableName)
		}
		if err != nil {
			return
		}

		// 删除旧索引
		for _, idx := range existingIndexes {
			if config.Type == consts.DBPgsql {
				alterStatements = append(alterStatements, fmt.Sprintf(`DROP INDEX IF EXISTS "%s"`, idx.Name))
			} else {
				alterStatements = append(alterStatements, fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`", in.TableName, idx.Name))
			}
		}

		// 创建新索引
		for _, idx := range in.Indexes {
			idxSQL := s.buildCreateIndexSQL(config.Type, in.TableName, &idx)
			if idxSQL != "" {
				alterStatements = append(alterStatements, idxSQL)
			}
		}
	}

	// 执行所有 ALTER 语句
	for _, stmt := range alterStatements {
		if _, err = g.DB(in.DbName).Ctx(ctx).Exec(ctx, stmt); err != nil {
			err = gerror.Wrapf(err, "执行修改表结构失败: %s", stmt)
			return
		}
	}
	return
}

// TableDrop 删除数据表
func (s *sSysGenTables) TableDrop(ctx context.Context, in *sysin.GenTableDropInp) (err error) {
	// 安全检查：禁止删除受保护的表
	disableTables := g.Cfg().MustGet(ctx, "hggen.disableTables").Strings()
	if gstr.InArray(disableTables, in.TableName) {
		err = gerror.Newf("表 %s 是系统保护表，禁止删除", in.TableName)
		return
	}

	config := g.DB(in.DbName).GetConfig()

	var sql string
	if config.Type == consts.DBPgsql {
		sql = fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, in.TableName)
	} else {
		sql = fmt.Sprintf("DROP TABLE IF EXISTS `%s`", in.TableName)
	}

	if _, err = g.DB(in.DbName).Ctx(ctx).Exec(ctx, sql); err != nil {
		err = gerror.Wrapf(err, "删除表 %s 失败", in.TableName)
		return
	}
	return
}

// PreviewDDL 预览DDL语句
func (s *sSysGenTables) PreviewDDL(ctx context.Context, in *sysin.GenTablePreviewDDLInp) (res *sysin.GenTablePreviewDDLModel, err error) {
	if in.IsEdit {
		// 编辑模式：生成 ALTER 语句预览
		ddl, e := s.buildAlterDDL(ctx, in.DbName, in.TableName, in.Comment, in.Engine, in.Columns, in.Indexes)
		if e != nil {
			err = e
			return
		}
		res = &sysin.GenTablePreviewDDLModel{
			DDL: ddl,
		}
		return
	}

	ddl, err := s.buildCreateDDL(ctx, in.DbName, in.TableName, in.Comment, in.Engine, in.Columns, in.Indexes)
	if err != nil {
		return
	}
	res = &sysin.GenTablePreviewDDLModel{
		DDL: ddl,
	}
	return
}

// buildAlterDDL 构建编辑模式的 ALTER TABLE DDL 预览（不执行，仅生成SQL）
func (s *sSysGenTables) buildAlterDDL(ctx context.Context, dbName, tableName, comment, engine string,
	columns []sysin.GenTableColumnInp, indexes []sysin.GenTableIndexInp) (ddl string, err error) {
	config := g.DB(dbName).GetConfig()

	// 自动添加表前缀
	fullTableName := tableName
	if config.Prefix != "" && !gstr.HasPrefix(tableName, config.Prefix) {
		fullTableName = config.Prefix + tableName
	}

	// 检查表是否存在
	exists, e := s.tableExists(ctx, dbName, fullTableName)
	if e != nil {
		err = e
		return
	}
	if !exists {
		err = gerror.Newf("表 %s 不存在，无法生成修改语句", fullTableName)
		return
	}

	// 获取现有字段
	var existingColumns []*sysin.GenTableViewColumnModel
	if config.Type == consts.DBPgsql {
		existingColumns, err = s.getPgsqlColumns(ctx, dbName, fullTableName)
	} else {
		existingColumns, err = s.getMysqlColumns(ctx, dbName, config.Name, fullTableName)
	}
	if err != nil {
		return
	}

	existingColMap := make(map[string]*sysin.GenTableViewColumnModel)
	for _, col := range existingColumns {
		existingColMap[col.Name] = col
	}

	// 生成ALTER语句
	var alterStatements []string

	// 新增和修改字段
	for i, col := range columns {
		if _, ok := existingColMap[col.Name]; ok {
			// 修改已有字段
			alterSQL := s.buildModifyColumnSQL(config.Type, fullTableName, &col)
			if alterSQL != "" {
				alterStatements = append(alterStatements, alterSQL)
			}
			delete(existingColMap, col.Name)
		} else {
			// 新增字段
			afterCol := ""
			if i > 0 {
				afterCol = columns[i-1].Name
			}
			alterSQL := s.buildAddColumnSQL(config.Type, fullTableName, &col, afterCol)
			alterStatements = append(alterStatements, alterSQL)
		}
	}

	// 删除不再需要的字段
	for colName := range existingColMap {
		if config.Type == consts.DBPgsql {
			alterStatements = append(alterStatements, fmt.Sprintf(`ALTER TABLE "%s" DROP COLUMN "%s"`, fullTableName, colName))
		} else {
			alterStatements = append(alterStatements, fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", fullTableName, colName))
		}
	}

	// 修改表注释
	if comment != "" {
		if config.Type == consts.DBPgsql {
			alterStatements = append(alterStatements, fmt.Sprintf(`COMMENT ON TABLE "%s" IS '%s'`, fullTableName, escapeString(comment)))
		} else {
			alterStatements = append(alterStatements, fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", fullTableName, escapeString(comment)))
		}
	}

	// 处理索引
	if len(indexes) > 0 {
		var existingIndexes []*sysin.GenTableViewIndexModel
		if config.Type == consts.DBPgsql {
			existingIndexes, err = s.getPgsqlIndexes(ctx, dbName, fullTableName)
		} else {
			existingIndexes, err = s.getMysqlIndexes(ctx, dbName, fullTableName)
		}
		if err != nil {
			return
		}

		// 删除旧索引
		for _, idx := range existingIndexes {
			if config.Type == consts.DBPgsql {
				alterStatements = append(alterStatements, fmt.Sprintf(`DROP INDEX IF EXISTS "%s"`, idx.Name))
			} else {
				alterStatements = append(alterStatements, fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`", fullTableName, idx.Name))
			}
		}

		// 创建新索引
		for _, idx := range indexes {
			idxSQL := s.buildCreateIndexSQL(config.Type, fullTableName, &idx)
			if idxSQL != "" {
				alterStatements = append(alterStatements, idxSQL)
			}
		}
	}

	if len(alterStatements) == 0 {
		ddl = "-- 无变更"
		return
	}

	ddl = strings.Join(alterStatements, ";\n") + ";"
	return
}

// buildCreateDDL 构建CREATE TABLE DDL
func (s *sSysGenTables) buildCreateDDL(ctx context.Context, dbName, tableName, comment, engine string, columns []sysin.GenTableColumnInp, indexes []sysin.GenTableIndexInp) (ddl string, err error) {
	config := g.DB(dbName).GetConfig()

	// 自动添加表前缀
	fullTableName := tableName
	if config.Prefix != "" && !gstr.HasPrefix(tableName, config.Prefix) {
		fullTableName = config.Prefix + tableName
	}

	// 检查表是否已存在
	exists, e := s.tableExists(ctx, dbName, fullTableName)
	if e != nil {
		err = e
		return
	}
	if exists {
		err = gerror.Newf("表 %s 已存在", fullTableName)
		return
	}

	if config.Type == consts.DBPgsql {
		ddl = s.buildPgsqlCreateDDL(fullTableName, comment, columns, indexes)
	} else {
		ddl = s.buildMysqlCreateDDL(fullTableName, comment, engine, columns, indexes)
	}
	return
}

// buildMysqlCreateDDL 构建MySQL CREATE TABLE DDL
func (s *sSysGenTables) buildMysqlCreateDDL(tableName, comment, engine string, columns []sysin.GenTableColumnInp, indexes []sysin.GenTableIndexInp) string {
	if engine == "" {
		engine = "InnoDB"
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("CREATE TABLE `%s` (\n", tableName))

	// 字段定义
	colDefs := make([]string, 0, len(columns))
	for _, col := range columns {
		colDef := s.buildMysqlColumnDef(&col)
		colDefs = append(colDefs, "  "+colDef)
	}

	// 主键
	pkCols := make([]string, 0)
	for _, col := range columns {
		if col.IsPrimaryKey {
			pkCols = append(pkCols, fmt.Sprintf("`%s`", col.Name))
		}
	}
	if len(pkCols) > 0 {
		colDefs = append(colDefs, fmt.Sprintf("  PRIMARY KEY (%s)", strings.Join(pkCols, ",")))
	}

	// 索引
	for _, idx := range indexes {
		idxDef := s.buildMysqlIndexDef(&idx)
		if idxDef != "" {
			colDefs = append(colDefs, "  "+idxDef)
		}
	}

	buf.WriteString(strings.Join(colDefs, ",\n"))
	buf.WriteString(fmt.Sprintf("\n) ENGINE=%s DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci", engine))

	if comment != "" {
		buf.WriteString(fmt.Sprintf(" COMMENT='%s'", escapeString(comment)))
	}
	buf.WriteString(";")
	return buf.String()
}

// buildPgsqlCreateDDL 构建PostgreSQL CREATE TABLE DDL
func (s *sSysGenTables) buildPgsqlCreateDDL(tableName, comment string, columns []sysin.GenTableColumnInp, indexes []sysin.GenTableIndexInp) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName))

	colDefs := make([]string, 0, len(columns))
	for _, col := range columns {
		colDef := s.buildPgsqlColumnDef(&col)
		colDefs = append(colDefs, "  "+colDef)
	}

	pkCols := make([]string, 0)
	for _, col := range columns {
		if col.IsPrimaryKey {
			pkCols = append(pkCols, fmt.Sprintf("\"%s\"", col.Name))
		}
	}
	if len(pkCols) > 0 {
		colDefs = append(colDefs, fmt.Sprintf("  PRIMARY KEY (%s)", strings.Join(pkCols, ",")))
	}

	buf.WriteString(strings.Join(colDefs, ",\n"))
	buf.WriteString("\n);\n")

	// 表注释
	if comment != "" {
		buf.WriteString(fmt.Sprintf("COMMENT ON TABLE \"%s\" IS '%s';\n", tableName, escapeString(comment)))
	}

	// 字段注释
	for _, col := range columns {
		if col.Comment != "" {
			buf.WriteString(fmt.Sprintf("COMMENT ON COLUMN \"%s\".\"%s\" IS '%s';\n", tableName, col.Name, escapeString(col.Comment)))
		}
	}

	// 索引
	for _, idx := range indexes {
		idxSQL := s.buildPgsqlIndexDef(tableName, &idx)
		if idxSQL != "" {
			buf.WriteString(idxSQL + ";\n")
		}
	}

	return buf.String()
}

// buildMysqlColumnDef 构建MySQL字段定义
func (s *sSysGenTables) buildMysqlColumnDef(col *sysin.GenTableColumnInp) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("`%s` ", col.Name))

	// 数据类型
	dataType := strings.ToLower(col.DataType)
	switch dataType {
	case "varchar", "char":
		length := col.Length
		if length <= 0 {
			length = 255
		}
		buf.WriteString(fmt.Sprintf("%s(%d)", dataType, length))
	case "decimal", "numeric", "float", "double":
		if col.Length > 0 {
			if col.Decimal > 0 {
				buf.WriteString(fmt.Sprintf("%s(%d,%d)", dataType, col.Length, col.Decimal))
			} else {
				buf.WriteString(fmt.Sprintf("%s(%d)", dataType, col.Length))
			}
		} else {
			buf.WriteString(dataType)
		}
	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint":
		if col.Length > 0 {
			buf.WriteString(fmt.Sprintf("%s(%d)", dataType, col.Length))
		} else {
			buf.WriteString(dataType)
		}
		if col.IsUnsigned {
			buf.WriteString(" unsigned")
		}
	default:
		buf.WriteString(dataType)
	}

	// 无符号（整数类型已在上面处理）
	if col.IsUnsigned && !isIntType(dataType) {
		// 非整数类型的 unsigned 在此处理（如 decimal unsigned）
		if dataType == "decimal" || dataType == "float" || dataType == "double" {
			buf.WriteString(" unsigned")
		}
	}

	// 是否允许 NULL
	if !col.IsNullable {
		buf.WriteString(" NOT NULL")
	} else {
		buf.WriteString(" NULL")
	}

	// 自增
	if col.IsAutoInc {
		buf.WriteString(" AUTO_INCREMENT")
	}

	// 默认值
	if col.DefaultValue != "" && !col.IsAutoInc {
		if isNumericDefault(col.DefaultValue) || col.DefaultValue == "CURRENT_TIMESTAMP" {
			buf.WriteString(fmt.Sprintf(" DEFAULT %s", col.DefaultValue))
		} else {
			buf.WriteString(fmt.Sprintf(" DEFAULT '%s'", escapeString(col.DefaultValue)))
		}
	}

	// 注释
	if col.Comment != "" {
		buf.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeString(col.Comment)))
	}

	return buf.String()
}

// buildPgsqlColumnDef 构建PostgreSQL字段定义
func (s *sSysGenTables) buildPgsqlColumnDef(col *sysin.GenTableColumnInp) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("\"%s\" ", col.Name))

	if col.IsAutoInc {
		if strings.ToLower(col.DataType) == "bigint" {
			buf.WriteString("BIGSERIAL")
		} else {
			buf.WriteString("SERIAL")
		}
	} else {
		dataType := s.mapMysqlToPgsqlType(col.DataType, col.Length, col.Decimal)
		buf.WriteString(dataType)
	}

	if !col.IsNullable && !col.IsAutoInc {
		buf.WriteString(" NOT NULL")
	}

	if col.DefaultValue != "" && !col.IsAutoInc {
		if isNumericDefault(col.DefaultValue) || col.DefaultValue == "CURRENT_TIMESTAMP" || col.DefaultValue == "NOW()" {
			buf.WriteString(fmt.Sprintf(" DEFAULT %s", col.DefaultValue))
		} else {
			buf.WriteString(fmt.Sprintf(" DEFAULT '%s'", escapeString(col.DefaultValue)))
		}
	}

	return buf.String()
}

// buildMysqlIndexDef 构建MySQL索引定义（建表内）
func (s *sSysGenTables) buildMysqlIndexDef(idx *sysin.GenTableIndexInp) string {
	if len(idx.Columns) == 0 {
		return ""
	}
	cols := make([]string, 0, len(idx.Columns))
	for _, c := range idx.Columns {
		cols = append(cols, fmt.Sprintf("`%s`", c))
	}
	colStr := strings.Join(cols, ",")

	name := idx.Name
	if name == "" {
		name = "idx_" + strings.Join(idx.Columns, "_")
	}

	switch strings.ToUpper(idx.Type) {
	case "UNIQUE":
		return fmt.Sprintf("UNIQUE KEY `%s` (%s)", name, colStr)
	case "FULLTEXT":
		return fmt.Sprintf("FULLTEXT KEY `%s` (%s)", name, colStr)
	default:
		return fmt.Sprintf("KEY `%s` (%s)", name, colStr)
	}
}

// buildPgsqlIndexDef 构建PostgreSQL索引定义
func (s *sSysGenTables) buildPgsqlIndexDef(tableName string, idx *sysin.GenTableIndexInp) string {
	if len(idx.Columns) == 0 {
		return ""
	}
	cols := make([]string, 0, len(idx.Columns))
	for _, c := range idx.Columns {
		cols = append(cols, fmt.Sprintf("\"%s\"", c))
	}
	colStr := strings.Join(cols, ",")

	name := idx.Name
	if name == "" {
		name = "idx_" + tableName + "_" + strings.Join(idx.Columns, "_")
	}

	switch strings.ToUpper(idx.Type) {
	case "UNIQUE":
		return fmt.Sprintf("CREATE UNIQUE INDEX \"%s\" ON \"%s\" (%s)", name, tableName, colStr)
	default:
		return fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\" (%s)", name, tableName, colStr)
	}
}

// buildCreateIndexSQL 构建独立的索引创建语句
func (s *sSysGenTables) buildCreateIndexSQL(dbType, tableName string, idx *sysin.GenTableIndexInp) string {
	if len(idx.Columns) == 0 {
		return ""
	}

	name := idx.Name
	if name == "" {
		name = "idx_" + strings.Join(idx.Columns, "_")
	}

	if dbType == consts.DBPgsql {
		cols := make([]string, 0, len(idx.Columns))
		for _, c := range idx.Columns {
			cols = append(cols, fmt.Sprintf("\"%s\"", c))
		}
		colStr := strings.Join(cols, ",")
		if strings.ToUpper(idx.Type) == "UNIQUE" {
			return fmt.Sprintf("CREATE UNIQUE INDEX \"%s\" ON \"%s\" (%s)", name, tableName, colStr)
		}
		return fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\" (%s)", name, tableName, colStr)
	}

	cols := make([]string, 0, len(idx.Columns))
	for _, c := range idx.Columns {
		cols = append(cols, fmt.Sprintf("`%s`", c))
	}
	colStr := strings.Join(cols, ",")

	switch strings.ToUpper(idx.Type) {
	case "UNIQUE":
		return fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE INDEX `%s` (%s)", tableName, name, colStr)
	case "FULLTEXT":
		return fmt.Sprintf("ALTER TABLE `%s` ADD FULLTEXT INDEX `%s` (%s)", tableName, name, colStr)
	default:
		return fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `%s` (%s)", tableName, name, colStr)
	}
}

// buildModifyColumnSQL 构建修改字段SQL
func (s *sSysGenTables) buildModifyColumnSQL(dbType, tableName string, col *sysin.GenTableColumnInp) string {
	if dbType == consts.DBPgsql {
		dataType := s.mapMysqlToPgsqlType(col.DataType, col.Length, col.Decimal)
		return fmt.Sprintf(`ALTER TABLE "%s" ALTER COLUMN "%s" TYPE %s`, tableName, col.Name, dataType)
	}

	colDef := s.buildMysqlColumnDef(col)
	return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s", tableName, colDef)
}

// buildAddColumnSQL 构建添加字段SQL
func (s *sSysGenTables) buildAddColumnSQL(dbType, tableName string, col *sysin.GenTableColumnInp, afterCol string) string {
	if dbType == consts.DBPgsql {
		colDef := s.buildPgsqlColumnDef(col)
		return fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN %s`, tableName, colDef)
	}

	colDef := s.buildMysqlColumnDef(col)
	sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", tableName, colDef)
	if afterCol != "" {
		sql += fmt.Sprintf(" AFTER `%s`", afterCol)
	}
	return sql
}

// tableExists 检查表是否存在
func (s *sSysGenTables) tableExists(ctx context.Context, dbName, tableName string) (bool, error) {
	config := g.DB(dbName).GetConfig()

	var sql string
	if config.Type == consts.DBPgsql {
		sql = fmt.Sprintf(`
			SELECT COUNT(1) as cnt FROM pg_class c
			JOIN pg_namespace n ON c.relnamespace = n.oid
			WHERE n.nspname = 'public' AND c.relname = '%s' AND c.relkind = 'r'`, tableName)
	} else {
		sql = fmt.Sprintf(`
			SELECT COUNT(1) as cnt FROM information_schema.TABLES 
			WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'`, config.Name, tableName)
	}

	result, err := g.DB(dbName).Ctx(ctx).Raw(sql).One()
	if err != nil {
		return false, gerror.Wrap(err, "检查表是否存在失败")
	}

	return result["cnt"].Int() > 0, nil
}

// mapMysqlToPgsqlType MySQL类型到PostgreSQL类型映射
func (s *sSysGenTables) mapMysqlToPgsqlType(dataType string, length, decimal int) string {
	switch strings.ToLower(dataType) {
	case "tinyint":
		return "SMALLINT"
	case "smallint":
		return "SMALLINT"
	case "mediumint", "int", "integer":
		return "INTEGER"
	case "bigint":
		return "BIGINT"
	case "float":
		return "REAL"
	case "double":
		return "DOUBLE PRECISION"
	case "decimal", "numeric":
		if length > 0 {
			if decimal > 0 {
				return fmt.Sprintf("NUMERIC(%d,%d)", length, decimal)
			}
			return fmt.Sprintf("NUMERIC(%d)", length)
		}
		return "NUMERIC"
	case "char":
		if length > 0 {
			return fmt.Sprintf("CHAR(%d)", length)
		}
		return "CHAR(255)"
	case "varchar":
		if length > 0 {
			return fmt.Sprintf("VARCHAR(%d)", length)
		}
		return "VARCHAR(255)"
	case "tinytext", "text", "mediumtext", "longtext":
		return "TEXT"
	case "tinyblob", "blob", "mediumblob", "longblob":
		return "BYTEA"
	case "date":
		return "DATE"
	case "datetime", "timestamp":
		return "TIMESTAMP"
	case "time":
		return "TIME"
	case "json":
		return "JSONB"
	default:
		return strings.ToUpper(dataType)
	}
}

// isIntType 判断是否为整数类型
func isIntType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		return true
	}
	return false
}

// isNumericDefault 判断默认值是否为数字
func isNumericDefault(val string) bool {
	_, err := strconv.ParseFloat(val, 64)
	return err == nil
}

// escapeString 转义SQL字符串
func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
