package gentpl

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// TableField represents a single field of a database table with extended metadata.
type TableField struct {
	gdb.TableField
	FieldName    string            // Field name after prefix removal
	OriginalName string            // Original field name in database
	LocalType    string            // Go local type name
	JsonCase     string            // JSON tag case style
	CustomTags   map[string]string // Custom tags from FieldMapping
}

// TableFields is a sortable slice of TableField pointers.
type TableFields []*TableField

func (s TableFields) Len() int      { return len(s) }
func (s TableFields) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TableFields) Less(i, j int) bool {
	return s[i].Index < s[j].Index
}

// Input holds parameters for type resolution.
type Input struct {
	StdTime      bool
	GJsonSupport bool
	TypeMapping  map[string]CustomAttributeType
	FieldMapping map[string]CustomAttributeType
	TableName    string
}

// resolveLocalType determines the Go local type for this field.
func (f *TableField) resolveLocalType(ctx context.Context, db gdb.DB, in Input) (appendImport string) {
	var (
		err              error
		localTypeName    gdb.LocalType
		localTypeNameStr string
	)

	// Try TypeMapping first
	if len(in.TypeMapping) > 0 {
		var tryTypeName string
		tryTypeMatch, _ := gregex.MatchString(`(.+?)\((.+)\)`, f.Type)
		if len(tryTypeMatch) == 3 {
			tryTypeName = gstr.Trim(tryTypeMatch[1])
		} else {
			tryTypeName = gstr.Split(f.Type, " ")[0]
		}
		if tryTypeName != "" {
			if typeMapping, ok := in.TypeMapping[strings.ToLower(tryTypeName)]; ok {
				localTypeNameStr = typeMapping.Type
				appendImport = typeMapping.Import
			}
		}
	}

	// Use DB type checking as fallback
	if localTypeNameStr == "" {
		localTypeName, err = db.CheckLocalTypeForField(ctx, f.Type, nil)
		if err != nil {
			panic(err)
		}
		localTypeNameStr = string(localTypeName)
		switch localTypeName {
		case gdb.LocalTypeDate, gdb.LocalTypeTime, gdb.LocalTypeDatetime:
			if in.StdTime {
				localTypeNameStr = "time.Time"
			} else {
				localTypeNameStr = "*gtime.Time"
				appendImport = "github.com/gogf/gf/v2/os/gtime"
			}
		case gdb.LocalTypeInt64Bytes:
			localTypeNameStr = "int64"
		case gdb.LocalTypeUint64Bytes:
			localTypeNameStr = "uint64"
		case gdb.LocalTypeJson, gdb.LocalTypeJsonb:
			if in.GJsonSupport {
				localTypeNameStr = "*gjson.Json"
				appendImport = "github.com/gogf/gf/v2/encoding/gjson"
			} else {
				localTypeNameStr = "string"
			}
		}
	}

	// Check field-specific mapping (overrides type mapping)
	if len(in.FieldMapping) > 0 {
		// Try table.field format
		fieldKey := fmt.Sprintf("%s.%s", in.TableName, f.FieldName)
		if typeMapping, ok := in.FieldMapping[fieldKey]; ok {
			localTypeNameStr = typeMapping.Type
			if typeMapping.Import != "" {
				appendImport = typeMapping.Import
			}
		}
		// Also try just field name
		if typeMapping, ok := in.FieldMapping[f.FieldName]; ok {
			if typeMapping.Type != "" {
				localTypeNameStr = typeMapping.Type
			}
			if typeMapping.Import != "" {
				appendImport = typeMapping.Import
			}
		}
	}

	f.LocalType = localTypeNameStr
	return
}

// NameCaseCamel returns CamelCase of the field name.
func (f *TableField) NameCaseCamel() string {
	return formatName(f.FieldName, "CaseCamel")
}

// NameCaseCamelLower returns lowerCamelCase of the field name.
func (f *TableField) NameCaseCamelLower() string {
	return formatName(f.FieldName, "CaseCamelLower")
}

// NameJsonCase returns the JSON-cased field name.
func (f *TableField) NameJsonCase() string {
	return gstr.CaseConvert(f.FieldName, gstr.CaseTypeMatch(f.JsonCase))
}

// IsNullable returns whether the field allows NULL values.
func (f *TableField) IsNullable() bool {
	return f.Null
}

// JsonTag generates the json tag value for this field.
func (f *TableField) JsonTag(omitempty bool, omitemptyAuto bool) string {
	if f.CustomTags != nil {
		if jsonTag, ok := f.CustomTags["json"]; ok {
			return jsonTag
		}
	}
	name := f.NameJsonCase()
	if omitempty || (omitemptyAuto && f.IsNullable()) {
		return name + ",omitempty"
	}
	return name
}

// OrmTag generates the orm tag value for this field.
func (f *TableField) OrmTag() string {
	if f.CustomTags != nil {
		if ormTag, ok := f.CustomTags["orm"]; ok {
			return ormTag
		}
	}
	return f.OriginalName
}

// DescriptionTagValue generates the description tag value for this field.
func (f *TableField) DescriptionTagValue() string {
	if f.CustomTags != nil {
		if descTag, ok := f.CustomTags["description"]; ok {
			return descTag
		}
	}
	return strings.ReplaceAll(formatComment(f.Comment), `"`, `\"`)
}

// TagBuildInput holds parameters for tag generation.
type TagBuildInput struct {
	NoJsonTag         bool
	JsonOmitempty     bool
	JsonOmitemptyAuto bool
	WithOrmTag        bool
	DescriptionTag    bool
	NoModelComment    bool
}

// TagWidths holds maximum widths for each tag type, used for alignment.
type TagWidths struct {
	JsonTagLen int
	OrmTagLen  int
	DescTagLen int
}

// CalcTagWidth calculates the tag width for a single field (without padding).
func (f *TableField) CalcTagWidth(in TagBuildInput) (jsonLen, ormLen, descLen int) {
	if !in.NoJsonTag {
		jsonLen = len(fmt.Sprintf(`json:"%s"`, f.JsonTag(in.JsonOmitempty, in.JsonOmitemptyAuto)))
	}
	if in.WithOrmTag {
		ormLen = len(fmt.Sprintf(`orm:"%s"`, f.OrmTag()))
	}
	if in.DescriptionTag {
		descLen = len(fmt.Sprintf(`description:"%s"`, f.DescriptionTagValue()))
	}
	return
}

// BuildTags builds the complete tag string for a field.
func (f *TableField) BuildTags(in TagBuildInput) string {
	return f.BuildTagsAligned(in, TagWidths{})
}

// BuildTagsAligned builds the complete tag string for a field with alignment.
func (f *TableField) BuildTagsAligned(in TagBuildInput, widths TagWidths) string {
	var tags []string

	// JSON tag
	if !in.NoJsonTag {
		jsonValue := f.JsonTag(in.JsonOmitempty, in.JsonOmitemptyAuto)
		tag := fmt.Sprintf(`json:"%s"`, jsonValue)
		if widths.JsonTagLen > 0 {
			tag = fmt.Sprintf("%-*s", widths.JsonTagLen, tag)
		}
		tags = append(tags, tag)
	}

	// ORM tag
	if in.WithOrmTag {
		ormValue := f.OrmTag()
		tag := fmt.Sprintf(`orm:"%s"`, ormValue)
		if widths.OrmTagLen > 0 {
			tag = fmt.Sprintf("%-*s", widths.OrmTagLen, tag)
		}
		tags = append(tags, tag)
	}

	// Description tag
	if in.DescriptionTag {
		descValue := f.DescriptionTagValue()
		tag := fmt.Sprintf(`description:"%s"`, descValue)
		// Description is the last standard tag, no need to pad
		tags = append(tags, tag)
	}

	// Custom tags from FieldMapping (excluding already-handled tags)
	if f.CustomTags != nil {
		var keys []string
		for k := range f.CustomTags {
			if k == "json" || k == "orm" || k == "description" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := f.CustomTags[k]
			tags = append(tags, fmt.Sprintf(`%s:"%s"`, k, v))
		}
	}

	if len(tags) == 0 {
		return ""
	}
	return "`" + strings.Join(tags, " ") + "`"
}
