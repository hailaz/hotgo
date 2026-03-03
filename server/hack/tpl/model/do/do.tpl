// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. {{.TplCreatedAtDatetimeStr}}
// =================================================================================

package do

{{.table.ImportsDo}}

// {{.TplTableNameCamelCase}} is the golang structure of table {{.TplTableName}} for DAO operations like Where/Data.
type {{.TplTableNameCamelCase}} struct {
g.Meta `orm:"table:{{.TplTableName}}, do:true"`
{{.table.StructDefineDo .tagInput}}
}
