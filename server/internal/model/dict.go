// Package model
package model

// Option 字典数据选项
type Option struct {
	Key       any    `json:"key"`
	Label     string `json:"label"     description:"字典标签"`
	Value     any    `json:"value"     description:"字典键值"`
	ValueType string `json:"valueType" description:"键值数据类型"`
	Type      string `json:"type"      description:"字典类型"`
	ListClass string `json:"listClass" description:"表格回显样式"`
	Extra     any    `json:"extra"     description:"额外数据配置"`
}
