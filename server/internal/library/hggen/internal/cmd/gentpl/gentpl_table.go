package gentpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/library/hggen/internal/utility/mlog"
)

// Table represents a database table with its metadata for template rendering.
type Table struct {
	TableName    string      // Original table name in database
	NewTableName string      // Table name after prefix processing
	Fields       TableFields // Sorted fields
	FieldMap     map[string]*gdb.TableField
	Imports      map[string]struct{}
	db           gdb.DB
}

// NewTable creates a Table from database metadata.
func NewTable(t *TplObj, tableName, newTableName string) (*Table, error) {
	fieldMap, err := t.db.TableFields(t.ctx, tableName)
	if err != nil {
		return nil, err
	}

	table := &Table{
		TableName:    tableName,
		NewTableName: newTableName,
		FieldMap:     fieldMap,
		db:           t.db,
		Imports:      make(map[string]struct{}),
	}

	table.buildFields(t.ctx, t.in)
	return table, nil
}

// buildFields converts raw gdb.TableField map to sorted TableFields.
func (t *Table) buildFields(ctx context.Context, in CGenTplInput) {
	removeFieldPrefixArray := gstr.SplitAndTrim(in.RemoveFieldPrefix, ",")
	names := sortFieldKeyByIndex(t.FieldMap)

	t.Fields = make(TableFields, 0, len(names))
	for _, name := range names {
		v := t.FieldMap[name]

		newFieldName := v.Name
		for _, prefix := range removeFieldPrefixArray {
			newFieldName = gstr.TrimLeftStr(newFieldName, prefix, 1)
		}

		field := &TableField{
			TableField:   *v,
			FieldName:    newFieldName,
			OriginalName: v.Name,
			JsonCase:     in.JsonCase,
			CustomTags:   make(map[string]string),
		}

		// Determine local type
		appendImport := field.resolveLocalType(ctx, t.db, Input{
			TypeMapping:  in.TypeMapping,
			FieldMapping: in.FieldMapping,
			StdTime:      in.StdTime,
			GJsonSupport: in.GJsonSupport,
			TableName:    t.TableName,
		})
		if appendImport != "" {
			t.Imports[appendImport] = struct{}{}
		}

		// Apply custom tags from FieldMapping
		if in.FieldMapping != nil {
			// Try table.field format first
			if fm, ok := in.FieldMapping[fmt.Sprintf("%s.%s", t.TableName, newFieldName)]; ok {
				if fm.Tags != nil {
					for tagName, tagValue := range fm.Tags {
						field.CustomTags[tagName] = tagValue
					}
				}
			}
			// Also try just field name
			if fm, ok := in.FieldMapping[newFieldName]; ok {
				if fm.Tags != nil {
					for tagName, tagValue := range fm.Tags {
						field.CustomTags[tagName] = tagValue
					}
				}
			}
		}

		t.Fields = append(t.Fields, field)
	}
}

// NameCaseCamel returns CamelCase of the new table name.
func (t *Table) NameCaseCamel() string {
	return formatName(t.NewTableName, "CaseCamel")
}

// NameCaseCamelLower returns lowerCamelCase of the new table name.
func (t *Table) NameCaseCamelLower() string {
	return formatName(t.NewTableName, "CaseCamelLower")
}

// NameCaseSnake returns snake_case of the new table name.
func (t *Table) NameCaseSnake() string {
	return gstr.CaseSnake(t.NewTableName)
}

// ColumnDefine generates column definition string for dao internal template.
func (t *Table) ColumnDefine() string {
	// Calculate max field name width for alignment
	maxNameLen := 0
	for _, f := range t.Fields {
		n := len(formatName(f.FieldName, "CaseCamel"))
		if n > maxNameLen {
			maxNameLen = n
		}
	}
	var lines []string
	for _, f := range t.Fields {
		name := formatName(f.FieldName, "CaseCamel")
		comment := formatComment(f.Comment)
		lines = append(lines, fmt.Sprintf("\t%-*s string // %s",
			maxNameLen, name,
			comment,
		))
	}
	return strings.Join(lines, "\n")
}

// ColumnNames generates column name assignments for dao internal template.
func (t *Table) ColumnNames() string {
	// Calculate max field name width for alignment
	maxNameLen := 0
	for _, f := range t.Fields {
		n := len(formatName(f.FieldName, "CaseCamel"))
		if n > maxNameLen {
			maxNameLen = n
		}
	}
	var lines []string
	for _, f := range t.Fields {
		name := formatName(f.FieldName, "CaseCamel")
		nameWithColon := name + ":"
		lines = append(lines, fmt.Sprintf("\t%-*s \"%s\",",
			maxNameLen+1, nameWithColon,
			f.OriginalName,
		))
	}
	return strings.Join(lines, "\n")
}

// StructDefineEntity generates struct field definitions for entity template.
func (t *Table) StructDefineEntity(tagInput TagBuildInput) string {
	// Calculate max widths for alignment
	maxNameLen := 0
	maxTypeLen := 0
	var tagWidths TagWidths
	for _, f := range t.Fields {
		n := len(formatName(f.FieldName, "CaseCamel"))
		if n > maxNameLen {
			maxNameLen = n
		}
		n = len(f.LocalType)
		if n > maxTypeLen {
			maxTypeLen = n
		}
		jsonLen, ormLen, descLen := f.CalcTagWidth(tagInput)
		if jsonLen > tagWidths.JsonTagLen {
			tagWidths.JsonTagLen = jsonLen
		}
		if ormLen > tagWidths.OrmTagLen {
			tagWidths.OrmTagLen = ormLen
		}
		if descLen > tagWidths.DescTagLen {
			tagWidths.DescTagLen = descLen
		}
	}
	var lines []string
	for _, f := range t.Fields {
		name := formatName(f.FieldName, "CaseCamel")
		tags := f.BuildTagsAligned(tagInput, tagWidths)
		comment := ""
		if !tagInput.NoModelComment {
			comment = fmt.Sprintf(" // %s", formatComment(f.Comment))
		}
		lines = append(lines, fmt.Sprintf("\t%-*s %-*s %s%s",
			maxNameLen, name,
			maxTypeLen, f.LocalType,
			tags,
			comment,
		))
	}
	return strings.Join(lines, "\n")
}

// StructDefineDo generates struct field definitions for DO template.
// All non-pointer/slice/map types are replaced with `any`.
func (t *Table) StructDefineDo(tagInput TagBuildInput) string {
	// Calculate max widths for alignment
	maxNameLen := 0
	maxTypeLen := 0
	for _, f := range t.Fields {
		n := len(formatName(f.FieldName, "CaseCamel"))
		if n > maxNameLen {
			maxNameLen = n
		}
		typeName := f.LocalType
		if !strings.HasPrefix(typeName, "*") && !strings.HasPrefix(typeName, "[]") && !strings.HasPrefix(typeName, "map") {
			typeName = "any"
		}
		n = len(typeName)
		if n > maxTypeLen {
			maxTypeLen = n
		}
	}
	var lines []string
	for _, f := range t.Fields {
		name := formatName(f.FieldName, "CaseCamel")
		typeName := f.LocalType
		if !strings.HasPrefix(typeName, "*") && !strings.HasPrefix(typeName, "[]") && !strings.HasPrefix(typeName, "map") {
			typeName = "any"
		}
		comment := fmt.Sprintf(" // %s", formatComment(f.Comment))
		lines = append(lines, fmt.Sprintf("\t%-*s %-*s%s",
			maxNameLen, name,
			maxTypeLen, typeName,
			comment,
		))
	}
	return strings.Join(lines, "\n")
}

// ImportsEntity returns the import statements needed for entity.
func (t *Table) ImportsEntity() string {
	return t.buildImportsString(false)
}

// ImportsDo returns the import statements needed for DO.
func (t *Table) ImportsDo() string {
	return t.buildImportsString(true)
}

// buildImportsString constructs the import block string.
func (t *Table) buildImportsString(isDo bool) string {
	var imports []string
	if isDo {
		imports = append(imports, `"github.com/gogf/gf/v2/frame/g"`)
	}

	// Scan field types for necessary imports
	for _, f := range t.Fields {
		localType := f.LocalType
		if isDo {
			if !strings.HasPrefix(localType, "*") && !strings.HasPrefix(localType, "[]") && !strings.HasPrefix(localType, "map") {
				continue
			}
		}
		if strings.Contains(localType, "gtime.Time") {
			imports = appendUnique(imports, `"github.com/gogf/gf/v2/os/gtime"`)
		} else if strings.Contains(localType, "time.Time") {
			imports = appendUnique(imports, `"time"`)
		}
		if strings.Contains(localType, "gjson.Json") {
			imports = appendUnique(imports, `"github.com/gogf/gf/v2/encoding/gjson"`)
		}
	}

	// Add custom imports from type/field mappings
	for imp := range t.Imports {
		if imp != "" {
			imports = appendUnique(imports, fmt.Sprintf(`"%s"`, imp))
		}
	}

	if len(imports) == 0 {
		return ""
	}
	return fmt.Sprintf("import(\n%s\n)", strings.Join(imports, "\n"))
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

// sortFieldKeyByIndex sorts field names by their index.
func sortFieldKeyByIndex(fieldMap map[string]*gdb.TableField) []string {
	names := make(map[int]string)
	for _, field := range fieldMap {
		names[field.Index] = field.Name
	}
	var (
		i      = 0
		j      = 0
		result = make([]string, len(names))
	)
	for len(names) != 0 {
		if val, ok := names[i]; ok {
			result[j] = val
			j++
			delete(names, i)
		}
		i++
	}
	return result
}

// formatName formats a field/table name to the specified case.
func formatName(name string, nameCase string) string {
	newName := name
	if isAllUpper(name) {
		newName = strings.ToLower(name)
	}
	switch nameCase {
	case "CaseCamel":
		return gstr.CaseCamel(newName)
	case "CaseCamelLower":
		return gstr.CaseCamelLower(newName)
	default:
		return newName
	}
}

func isAllUpper(s string) bool {
	for _, b := range s {
		if b >= 'a' && b <= 'z' {
			return false
		}
	}
	return true
}

func formatComment(comment string) string {
	comment = gstr.ReplaceByArray(comment, g.SliceStr{
		"\n", " ",
		"\r", " ",
	})
	comment = gstr.Replace(comment, `\n`, " ")
	comment = gstr.Trim(comment)
	return comment
}

// GetTables retrieves all filtered tables from the database.
func (t *TplObj) GetTables() ([]*Table, error) {
	tableNames, err := t.getTableNames()
	if err != nil {
		return nil, err
	}

	removePrefixArray := gstr.SplitAndTrim(t.in.RemovePrefix, ",")
	var tables []*Table
	for _, tableName := range tableNames {
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		newTableName = t.in.Prefix + newTableName

		table, err := NewTable(t, tableName, newTableName)
		if err != nil {
			mlog.Printf("warning: skipping table '%s': %v", tableName, err)
			continue
		}
		tables = append(tables, table)
	}
	return tables, nil
}
