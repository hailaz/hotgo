// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. {{.TplCreatedAtDatetimeStr}}
// =================================================================================

package entity

{{.table.ImportsEntity}}

// {{.TplTableNameCamelCase}} is the golang structure for table {{.table.NewTableName}}.
type {{.TplTableNameCamelCase}} struct {
{{.table.StructDefineEntity .tagInput}}
}
