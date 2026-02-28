// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

// GenTableColumnInp 数据表字段配置
type GenTableColumnInp struct {
	Name         string `json:"name"         v:"required|regex:^[a-zA-Z_][a-zA-Z0-9_]*$#字段名不能为空|字段名格式不正确" dc:"字段名"`
	DataType     string `json:"dataType"     v:"required#数据类型不能为空" dc:"数据类型"`
	Length       int    `json:"length"       dc:"长度"`
	Decimal      int    `json:"decimal"      dc:"小数位"`
	IsNullable   bool   `json:"isNullable"   dc:"允许为空"`
	DefaultValue string `json:"defaultValue" dc:"默认值"`
	Comment      string `json:"comment"      dc:"字段注释"`
	IsPrimaryKey bool   `json:"isPrimaryKey" dc:"是否主键"`
	IsAutoInc    bool   `json:"isAutoInc"    dc:"是否自增"`
	IsUnsigned   bool   `json:"isUnsigned"   dc:"是否无符号"`
}

// GenTableIndexInp 索引配置
type GenTableIndexInp struct {
	Name    string   `json:"name"    dc:"索引名"`
	Type    string   `json:"type"    dc:"索引类型"` // INDEX/UNIQUE/FULLTEXT
	Columns []string `json:"columns" v:"required#索引字段不能为空" dc:"索引字段"`
}

// GenTableDbSelectModel 数据库选项
type GenTableDbSelectModel struct {
	Value string `json:"value" dc:"数据库配置名称"`
	Label string `json:"label" dc:"显示名称"`
}

// GenTableListInp 获取表列表输入
type GenTableListInp struct {
	DbName    string `json:"dbName"    dc:"数据库配置名称"`
	TableName string `json:"tableName" dc:"表名搜索"`
}

// GenTableListModel 表列表项
type GenTableListModel struct {
	TableName    string `json:"tableName"    dc:"表名"`
	TableComment string `json:"tableComment" dc:"表注释"`
	Engine       string `json:"engine"       dc:"存储引擎"`
	TableRows    int64  `json:"tableRows"    dc:"数据行数"`
	CreateTime   string `json:"createTime"   dc:"创建时间"`
	TableCollation string `json:"tableCollation" dc:"字符集"`
}

// GenTableViewInp 查看表结构输入
type GenTableViewInp struct {
	DbName    string `json:"dbName"    v:"required#数据库不能为空" dc:"数据库配置名称"`
	TableName string `json:"tableName" v:"required#表名不能为空" dc:"表名"`
}

// GenTableViewColumnModel 表字段详情
type GenTableViewColumnModel struct {
	Name         string `json:"name"         dc:"字段名"`
	DataType     string `json:"dataType"     dc:"数据类型"`
	Length       int    `json:"length"       dc:"长度"`
	Decimal      int    `json:"decimal"      dc:"小数位"`
	IsNullable   bool   `json:"isNullable"   dc:"允许为空"`
	DefaultValue string `json:"defaultValue" dc:"默认值"`
	Comment      string `json:"comment"      dc:"字段注释"`
	IsPrimaryKey bool   `json:"isPrimaryKey" dc:"是否主键"`
	IsAutoInc    bool   `json:"isAutoInc"    dc:"是否自增"`
	IsUnsigned   bool   `json:"isUnsigned"   dc:"是否无符号"`
	ColumnType   string `json:"columnType"   dc:"完整列类型"`
}

// GenTableViewIndexModel 表索引详情
type GenTableViewIndexModel struct {
	Name    string   `json:"name"    dc:"索引名"`
	Type    string   `json:"type"    dc:"索引类型"`
	Columns []string `json:"columns" dc:"索引字段"`
}

// GenTableViewModel 表结构详情
type GenTableViewModel struct {
	TableName    string                     `json:"tableName"    dc:"表名"`
	TableComment string                     `json:"tableComment" dc:"表注释"`
	Engine       string                     `json:"engine"       dc:"存储引擎"`
	Charset      string                     `json:"charset"      dc:"字符集"`
	Columns      []*GenTableViewColumnModel `json:"columns"      dc:"字段列表"`
	Indexes      []*GenTableViewIndexModel  `json:"indexes"      dc:"索引列表"`
}

// GenTableCreateInp 创建数据表输入
type GenTableCreateInp struct {
	DbName    string              `json:"dbName"    v:"required#数据库不能为空" dc:"数据库配置名称"`
	TableName string              `json:"tableName" v:"required|regex:^[a-zA-Z_][a-zA-Z0-9_]*$#表名不能为空|表名格式不正确" dc:"表名"`
	Comment   string              `json:"comment"   dc:"表注释"`
	Engine    string              `json:"engine"    dc:"存储引擎"`
	Columns   []GenTableColumnInp `json:"columns"   v:"required#字段列表不能为空" dc:"字段列表"`
	Indexes   []GenTableIndexInp  `json:"indexes"   dc:"索引列表"`
}

// GenTableEditInp 修改表结构输入
type GenTableEditInp struct {
	DbName    string              `json:"dbName"    v:"required#数据库不能为空" dc:"数据库配置名称"`
	TableName string              `json:"tableName" v:"required#表名不能为空" dc:"表名"`
	Comment   string              `json:"comment"   dc:"表注释"`
	Engine    string              `json:"engine"    dc:"存储引擎"`
	Columns   []GenTableColumnInp `json:"columns"   v:"required#字段列表不能为空" dc:"字段列表"`
	Indexes   []GenTableIndexInp  `json:"indexes"   dc:"索引列表"`
}

// GenTableDropInp 删除数据表输入
type GenTableDropInp struct {
	DbName    string `json:"dbName"    v:"required#数据库不能为空" dc:"数据库配置名称"`
	TableName string `json:"tableName" v:"required#表名不能为空" dc:"表名"`
}

// GenTablePreviewDDLInp 预览DDL输入
type GenTablePreviewDDLInp struct {
	DbName    string              `json:"dbName"    v:"required#数据库不能为空" dc:"数据库配置名称"`
	TableName string              `json:"tableName" v:"required|regex:^[a-zA-Z_][a-zA-Z0-9_]*$#表名不能为空|表名格式不正确" dc:"表名"`
	Comment   string              `json:"comment"   dc:"表注释"`
	Engine    string              `json:"engine"    dc:"存储引擎"`
	Columns   []GenTableColumnInp `json:"columns"   v:"required#字段列表不能为空" dc:"字段列表"`
	Indexes   []GenTableIndexInp  `json:"indexes"   dc:"索引列表"`
	IsEdit    bool                `json:"isEdit"    dc:"是否编辑模式"`
}

// GenTablePreviewDDLModel 预览DDL输出
type GenTablePreviewDDLModel struct {
	DDL string `json:"ddl" dc:"DDL语句"`
}
