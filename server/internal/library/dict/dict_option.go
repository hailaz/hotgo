// Package dict
package dict

import (
	"hash/fnv"

	"github.com/gogf/gf/v2/util/gconv"

	"hotgo/internal/model"
)

// GenDefaultOption 生成默认表格回显样式
func GenDefaultOption(key any, label string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: "default",
		Extra:     e,
	}
}

func GenSuccessOption(key any, label string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: "success",
		Extra:     e,
	}
}

func GenWarningOption(key any, label string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: "warning",
		Extra:     e,
	}
}

func GenErrorOption(key any, label string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: "error",
		Extra:     e,
	}
}

func GenInfoOption(key any, label string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: "info",
		Extra:     e,
	}
}

func GenPrimaryOption(key any, label string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: "primary",
		Extra:     e,
	}
}

// GenCustomOption 生成自定义表格回显样式
func GenCustomOption(key any, label string, custom string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: custom,
		Extra:     e,
	}
}

// GenHashOption 根据不同label以hash算法生成表格回显样式
func GenHashOption(key any, label string, extra ...any) *model.Option {
	e := any(nil)
	if len(extra) > 0 {
		e = extra[0]
	}
	strings := []string{"default", "primary", "info", "success", "warning", "error"}
	hash := fnv.New32()

	tag := "default"
	if _, err := hash.Write(gconv.Bytes(label)); err == nil {
		index := uint32(hash.Sum32()) % uint32(len(strings))
		if index < uint32(len(strings)) {
			tag = strings[index]
		}
	}
	return &model.Option{
		Key:       key,
		Label:     label,
		Value:     key,
		ListClass: tag,
		Extra:     e,
	}
}

// GetOption 通过key找到label
func GetOption(ses []*model.Option, key any) *model.Option {
	for _, v := range ses {
		if gconv.String(v.Key) == gconv.String(key) {
			return v
		}
	}
	return nil
}

// GetOptionLabel 通过key找到label
func GetOptionLabel(ses []*model.Option, key any) string {
	for _, v := range ses {
		if gconv.String(v.Key) == gconv.String(key) {
			return v.Label
		}
	}
	return ``
}

// HasOptionKey 是否存在指定key
func HasOptionKey(ses []*model.Option, key any) bool {
	for _, v := range ses {
		if gconv.String(v.Key) == gconv.String(key) {
			return true
		}
	}
	return false
}

// UniqueOption 去重选项
func UniqueOption(src []*model.Option) (dst []*model.Option) {
	temp := map[string]struct{}{}
	for _, item := range src {
		key := gconv.String(item.Key)
		if _, ok := temp[key]; !ok {
			temp[key] = struct{}{}
			dst = append(dst, item)
		}
	}
	return dst
}
