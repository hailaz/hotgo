// Package form
package form

// Sorter 排序器，兼容naiveUI
type Sorter struct {
	ColumnKey string `json:"columnKey"  dc:"排序字段"`
	Order     any    `json:"order"      dc:"排序方式 descend|ascend|false"`
}

type Sorters struct {
	Sorters []*Sorter `json:"sorters"  dc:"排序器"`
}

func (s *Sorters) GetSorters() []*Sorter {
	return s.Sorters
}
