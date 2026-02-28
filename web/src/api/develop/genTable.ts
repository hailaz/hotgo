import { http } from '@/utils/http/axios';

// 获取可选数据库列表
export function DbSelect(params?) {
  return http.request({
    url: '/genTable/dbSelect',
    method: 'get',
    params,
  });
}

// 获取数据库表列表
export function TableList(params) {
  return http.request({
    url: '/genTable/tableList',
    method: 'get',
    params,
  });
}

// 查看表结构详情
export function TableView(params) {
  return http.request({
    url: '/genTable/tableView',
    method: 'get',
    params,
  });
}

// 创建数据表
export function TableCreate(params) {
  return http.request({
    url: '/genTable/tableCreate',
    method: 'POST',
    params,
  });
}

// 修改表结构
export function TableEdit(params) {
  return http.request({
    url: '/genTable/tableEdit',
    method: 'POST',
    params,
  });
}

// 删除数据表
export function TableDrop(params) {
  return http.request({
    url: '/genTable/tableDrop',
    method: 'POST',
    params,
  });
}

// 预览DDL语句
export function PreviewDDL(params) {
  return http.request({
    url: '/genTable/previewDDL',
    method: 'POST',
    params,
  });
}
