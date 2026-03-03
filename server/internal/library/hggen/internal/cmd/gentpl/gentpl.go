// Package gentpl provides template-based code generation for dao/do/entity files.
package gentpl

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/library/hggen/internal/utility/mlog"
	"hotgo/internal/library/hggen/internal/utility/utils"
)

type (
	CGenTplInput struct {
		Path              string `name:"path" short:"p" d:"internal"`
		Link              string `name:"link" short:"l"`
		TplPath           string `name:"tplPath" d:"hack/tpl"`
		Tables            string `name:"tables" short:"t"`
		TablesEx          string `name:"tablesEx" short:"x"`
		Group             string `name:"group" short:"g" d:"default"`
		Prefix            string `name:"prefix" short:"f"`
		RemovePrefix      string `name:"removePrefix" short:"r"`
		RemoveFieldPrefix string `name:"removeFieldPrefix" short:"rf"`
		JsonCase          string `name:"jsonCase" short:"j" d:"CamelLower"`
		ImportPrefix      string `name:"importPrefix" short:"i"`
		DaoPath           string `name:"daoPath" short:"d" d:"dao"`
		DoPath            string `name:"doPath" short:"o" d:"model/do"`
		EntityPath        string `name:"entityPath" short:"e" d:"model/entity"`
		StdTime           bool   `name:"stdTime" short:"s"`
		WithTime          bool   `name:"withTime" short:"w"`
		GJsonSupport      bool   `name:"gJsonSupport" short:"n"`
		OverwriteDao      bool   `name:"overwriteDao" short:"v"`
		DescriptionTag    bool   `name:"descriptionTag" short:"c"`
		NoJsonTag         bool   `name:"noJsonTag" short:"k"`
		NoModelComment    bool   `name:"noModelComment" short:"m"`
		Clear             bool   `name:"clear" short:"a"`

		// gen tpl extended fields
		JsonOmitempty     bool `name:"jsonOmitempty"`
		JsonOmitemptyAuto bool `name:"jsonOmitemptyAuto"`
		WithOrmTag        bool `name:"withOrmTag" d:"true"`

		TypeMapping  map[string]CustomAttributeType `name:"typeMapping" short:"y"`
		FieldMapping map[string]CustomAttributeType `name:"fieldMapping" short:"fm"`
	}

	CustomAttributeType struct {
		Type   string            `brief:"custom attribute type name"`
		Import string            `brief:"custom import for this type"`
		Tags   map[string]string `brief:"custom tags for this field"`
	}
)

var defaultTypeMapping = map[string]CustomAttributeType{
	"decimal":    {Type: "float64"},
	"money":      {Type: "float64"},
	"numeric":    {Type: "float64"},
	"smallmoney": {Type: "float64"},
}

// TplObj is the context object for template-based code generation.
type TplObj struct {
	ctx        context.Context
	in         CGenTplInput
	db         gdb.DB
	tplPathAbs string
}

// NewTpl creates and returns a new TplObj.
func NewTpl(ctx context.Context, in CGenTplInput) (*TplObj, error) {
	db, err := getDB(in)
	if err != nil {
		return nil, err
	}
	return &TplObj{
		ctx:        ctx,
		in:         in,
		db:         db,
		tplPathAbs: gfile.Abs(in.TplPath),
	}, nil
}

// getDB gets database instance from link or group config.
func getDB(in CGenTplInput) (db gdb.DB, err error) {
	if in.Link != "" {
		var tempGroup = gtime.TimestampNanoStr()
		if err = gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
			Link: in.Link,
		}); err != nil {
			return nil, gerror.Wrapf(err, "database configuration failed")
		}
		if db, err = gdb.Instance(tempGroup); err != nil {
			return nil, gerror.Wrapf(err, "database initialization failed")
		}
	} else {
		db = g.DB(in.Group)
	}
	if db == nil {
		return nil, gerror.New("database initialization failed, may be invalid database configuration")
	}
	return
}

// Tpl is the main entry point for template-based code generation.
func Tpl(ctx context.Context, in CGenTplInput) error {
	// Merge default typeMapping
	if in.TypeMapping == nil {
		in.TypeMapping = defaultTypeMapping
	} else {
		for key, typeMapping := range defaultTypeMapping {
			if _, ok := in.TypeMapping[key]; !ok {
				in.TypeMapping[key] = typeMapping
			}
		}
	}

	tplObj, err := NewTpl(ctx, in)
	if err != nil {
		return err
	}

	// Get all tpl files
	tplList, err := tplObj.GetTplFileList()
	if err != nil {
		return gerror.Wrapf(err, "scanning template files failed")
	}
	if len(tplList) == 0 {
		return gerror.Newf("no .tpl template files found in: %s", in.TplPath)
	}

	// Get table names and filter
	tableNames, err := tplObj.getTableNames()
	if err != nil {
		return err
	}

	// Process prefix removal
	removePrefixArray := gstr.SplitAndTrim(in.RemovePrefix, ",")
	newTableNames := make([]string, len(tableNames))
	for i, tableName := range tableNames {
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		newTableName = in.Prefix + newTableName
		newTableNames[i] = newTableName
	}

	// Generate files for each table × each template
	for i, tableName := range tableNames {
		newTableName := newTableNames[i]
		table, err := NewTable(tplObj, tableName, newTableName)
		if err != nil {
			mlog.Printf("warning: skipping table '%s': %v", tableName, err)
			continue
		}

		for _, tplFile := range tplList {
			err = tplObj.generateFromTemplate(ctx, table, tplFile)
			if err != nil {
				return gerror.Wrapf(err, "generating from template '%s' for table '%s' failed", tplFile, tableName)
			}
		}
	}

	// Format generated files
	utils.GoFmt(in.Path)
	mlog.Print("done!")
	return nil
}

// GetTplFileList scans and returns all .tpl files in the template directory.
func (t *TplObj) GetTplFileList() ([]string, error) {
	return gfile.ScanDirFile(t.tplPathAbs, "*.tpl", true)
}

// getTableNames retrieves table names from database with filtering.
func (t *TplObj) getTableNames() ([]string, error) {
	var (
		tableNames []string
		err        error
	)
	if t.in.Tables != "" {
		tableNames = gstr.SplitAndTrim(t.in.Tables, ",")
	} else {
		tableNames, err = t.db.Tables(t.ctx)
		if err != nil {
			return nil, gerror.Wrapf(err, "fetching tables failed")
		}
	}

	// Table excluding
	if t.in.TablesEx != "" {
		array := garray.NewStrArrayFrom(tableNames)
		for _, p := range gstr.SplitAndTrim(t.in.TablesEx, ",") {
			if gstr.Contains(p, "*") || gstr.Contains(p, "?") {
				p = gstr.ReplaceByMap(p, map[string]string{
					"\r": "",
					"\n": "",
				})
				p = gstr.ReplaceByMap(p, map[string]string{
					"*": "\r",
					"?": "\n",
				})
				p = gregex.Quote(p)
				p = gstr.ReplaceByMap(p, map[string]string{
					"\r": ".*",
					"\n": ".",
				})
				for _, v := range array.Clone().Slice() {
					if gregex.IsMatchString(p, v) {
						array.RemoveValue(v)
					}
				}
			} else {
				array.RemoveValue(p)
			}
		}
		tableNames = array.Slice()
	}
	return tableNames, nil
}

// generateFromTemplate generates code from a single template file for a single table.
func (t *TplObj) generateFromTemplate(ctx context.Context, table *Table, tplFile string) error {
	// Calculate relative path from tpl root
	relativePath := strings.TrimPrefix(gfile.Dir(tplFile), t.tplPathAbs)
	relativePath = filepath.ToSlash(relativePath)
	relativePath = strings.TrimPrefix(relativePath, "/")

	// Determine output file name
	tableNameSnakeCase := gstr.CaseSnake(table.NewTableName)
	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	if len(fileName) > 5 && fileName[len(fileName)-5:] == "_test" {
		fileName += "_table"
	}
	outputFileName := fileName + ".gen.go"

	// Calculate output path
	var outputPath string
	if relativePath != "" {
		outputPath = filepath.Join(t.in.Path, relativePath, outputFileName)
	} else {
		outputPath = filepath.Join(t.in.Path, outputFileName)
	}
	outputPath = filepath.FromSlash(outputPath)

	// For dao index files, check overwrite policy
	if strings.Contains(relativePath, "dao") && !strings.Contains(relativePath, "internal") &&
		!strings.Contains(relativePath, "model") {
		if !t.in.OverwriteDao && gfile.Exists(outputPath) {
			return nil
		}
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if !gfile.Exists(outputDir) {
		if err := gfile.Mkdir(outputDir); err != nil {
			return gerror.Wrapf(err, "creating output directory '%s' failed", outputDir)
		}
	}

	// Build template data
	tplData := t.buildTplData(table, relativePath)

	// Parse template
	view := gview.New()
	content, err := view.Parse(ctx, tplFile, tplData)
	if err != nil {
		return gerror.Wrapf(err, "parsing template '%s' failed", tplFile)
	}

	// Write file
	if err = gfile.PutContents(outputPath, strings.TrimSpace(content)); err != nil {
		return gerror.Wrapf(err, "writing content to '%s' failed", outputPath)
	}
	mlog.Print("generated:", outputPath)
	return nil
}

// buildTplData constructs the template data map for gview rendering.
func (t *TplObj) buildTplData(table *Table, relativePath string) g.Map {
	var tplCreatedAtDatetimeStr string
	if t.in.WithTime {
		tplCreatedAtDatetimeStr = fmt.Sprintf("Created at %s", gtime.Now().String())
	}

	// Calculate import prefix for the current relative path
	importPrefix := t.in.ImportPrefix
	if importPrefix == "" {
		importPrefix = utils.GetImportPath(filepath.Join(t.in.Path, relativePath))
	} else if relativePath != "" {
		importPrefix = gstr.Join(g.SliceStr{importPrefix, relativePath}, "/")
	}

	// Determine package name from the relative path
	packageName := filepath.Base(relativePath)
	if packageName == "" || packageName == "." {
		packageName = filepath.Base(t.in.Path)
	}

	return g.Map{
		// table info
		"table":  table,
		"tables": nil, // reserved for future use

		// legacy compatible variables (used in templates)
		"TplTableName":               table.TableName,
		"TplTableNameCamelCase":      table.NameCaseCamel(),
		"TplTableNameCamelLowerCase": table.NameCaseCamelLower(),
		"TplImportPrefix":            importPrefix,
		"TplGroupName":               t.in.Group,
		"TplDatetimeStr":             gtime.Now().String(),
		"TplCreatedAtDatetimeStr":    tplCreatedAtDatetimeStr,
		"TplPackageName":             packageName,

		// tag build input
		"tagInput": TagBuildInput{
			NoJsonTag:         t.in.NoJsonTag,
			JsonOmitempty:     t.in.JsonOmitempty,
			JsonOmitemptyAuto: t.in.JsonOmitemptyAuto,
			WithOrmTag:        t.in.WithOrmTag,
			DescriptionTag:    t.in.DescriptionTag,
			NoModelComment:    t.in.NoModelComment,
		},
	}
}
