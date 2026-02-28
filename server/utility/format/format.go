// Package format
package format

import (
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// Round2String 四舍五入保留小数，默认2位
func Round2String(value float64, args ...any) string {
	var places = 2
	if len(args) > 0 {
		places = gconv.Int(args[0])
	}
	return fmt.Sprintf("%0."+strconv.Itoa(places)+"f", value)
}

// Round2Float64 四舍五入保留小数，默认2位
func Round2Float64(value float64, args ...any) float64 {
	return gconv.Float64(Round2String(value, args...))
}

// FileSize 字节的单位转换 保留两位小数
func FileSize(data int64) string {
	var factor float64 = 1024
	res := float64(data)
	for _, unit := range []string{"", "K", "M", "G", "T", "P"} {
		if res < factor {
			return fmt.Sprintf("%.2f%sB", res, unit)
		}
		res /= factor
	}
	return fmt.Sprintf("%.2f%sB", res, "P")
}

// AgoTime 多久以前
func AgoTime(gt *gtime.Time) string {
	if gt == nil {
		return ""
	}
	n := gtime.Now().Timestamp()
	t := gt.Timestamp()

	var ys int64 = 31536000
	var ds int64 = 86400
	var hs int64 = 3600
	var ms int64 = 60
	var ss int64 = 1
	var rs string

	d := n - t
	switch {
	case d > ys:
		rs = fmt.Sprintf("%d年前", int(d/ys))
	case d > ds:
		rs = fmt.Sprintf("%d天前", int(d/ds))
	case d > hs:
		rs = fmt.Sprintf("%d小时前", int(d/hs))
	case d > ms:
		rs = fmt.Sprintf("%d分钟前", int(d/ms))
	case d > ss:
		rs = fmt.Sprintf("%d秒前", int(d/ss))
	default:
		rs = "刚刚"
	}
	return rs
}
