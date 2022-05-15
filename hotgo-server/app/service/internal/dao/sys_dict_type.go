// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"context"
	"github.com/bufanyun/hotgo/app/consts"
	"github.com/bufanyun/hotgo/app/model/entity"
	"github.com/bufanyun/hotgo/app/service/internal/dao/internal"
	"github.com/gogf/gf/v2/errors/gerror"
)

// sysDictTypeDao is the data access object for table hg_sys_dict_type.
// You can define custom methods on it to extend its functionality as you wish.
type sysDictTypeDao struct {
	*internal.SysDictTypeDao
}

var (
	// SysDictType is globally public accessible object for table hg_sys_dict_type operations.
	SysDictType = sysDictTypeDao{
		internal.NewSysDictTypeDao(),
	}
)

//
//  @Title  判断字典类型是否唯一
//  @Description
//  @Author  Ms <133814250@qq.com>
//  @Param   ctx
//  @Param   id
//  @Param   dictType
//  @Return  bool
//  @Return  error
//
func (dao *sysDictTypeDao) IsUnique(ctx context.Context, id int64, dictType string) (bool, error) {
	var (
		data *entity.SysDictType
		err  error
	)
	m := dao.Ctx(ctx).Where("type", dictType)

	if id > 0 {
		m = m.WhereNot("id", id)
	}

	if err = m.Scan(&data); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return false, err
	}

	if data == nil {
		return true, nil
	}

	return false, nil

}

// Fill with you ideas below.