-- 数据表管理菜单权限SQL
-- 功能：为开发工具新增「数据表管理」菜单及相关权限
-- Date: 2026-02-28
-- Link https://github.com/bufanyun/hotgo

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;

SET @now := now();

-- 数据表管理菜单页（pid=2097 即「开发工具」目录）
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES (NULL, 2097, '数据表管理', 'develop_table', 'table', 'TableOutlined', 2, '', '/genTable/tableList,/genTable/dbSelect', '', '/develop/table/index', 1, '', 0, 0, '', 0, 0, 0, 2, 'tr_2097 ', 5, '可视化管理数据库表结构', 1, @now, @now);
SET @tableListId = LAST_INSERT_ID();

-- 设计表页面（隐藏菜单，点击新建/编辑时跳转）
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES (NULL, 2097, '设计表', 'develop_table_design', 'table/design', '', 2, '', '/genTable/tableCreate,/genTable/tableEdit,/genTable/previewDDL', '', '/develop/table/design', 0, 'develop_table', 0, 0, '', 0, 1, 0, 2, 'tr_2097 ', 6, '可视化建表/编辑表', 1, @now, @now);
SET @designId = LAST_INSERT_ID();

-- 按钮权限：查看表结构
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES (NULL, @tableListId, '查看表结构', 'genTableView', '', '', 3, '', '/genTable/tableView', '', '', 1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_2097 tr_', @tableListId, ' '), 10, '', 1, @now, @now);

-- 按钮权限：删除表
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES (NULL, @tableListId, '删除表', 'genTableDrop', '', '', 3, '', '/genTable/tableDrop', '', '', 1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_2097 tr_', @tableListId, ' '), 20, '', 1, @now, @now);

-- 将菜单分配给超管角色(role_id=1)
INSERT INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (1, @tableListId);
INSERT INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (1, @designId);
INSERT INTO `hg_admin_role_menu` (`role_id`, `menu_id`) SELECT 1, `id` FROM `hg_admin_menu` WHERE `pid` = @tableListId;

COMMIT;
