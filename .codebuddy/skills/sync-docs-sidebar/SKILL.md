---
name: sync-docs-sidebar
description: >-
  同步 HotGo 项目的文档目录。扫描 docs/guide-zh-CN/ 下所有 md 文件，
  将未列入目录的文件自动按分类规则添加到 docs/guide-zh-CN/README.md，
  然后同步更新 server/resource/public/docs/ 下的 README.md 和 sidebar.md，
  使 docsify 文档站点导航保持最新。当用户提到"同步文档目录"、"更新文档侧边栏"、
  "文档目录不是最新"等场景时，应使用此 skill。
---

# 同步文档目录 (sync-docs-sidebar)

## 概述

HotGo 项目使用 docsify 在 `server/resource/public/docs/` 下提供在线文档站点。文档源文件位于 `docs/guide-zh-CN/` 目录。

此 skill 完成两件事：
1. **更新源目录**：扫描 `docs/guide-zh-CN/` 下所有 md 文件，将未列入 `docs/guide-zh-CN/README.md` 的文件自动按前缀分类添加进去
2. **同步到 docsify**：将更新后的目录同步到 `server/resource/public/docs/README.md` 和 `sidebar.md`

## 文件关系

| 角色 | 文件路径 |
|------|---------|
| 源目录（会被更新） | `docs/guide-zh-CN/README.md` |
| docsify 首页（同步生成） | `server/resource/public/docs/README.md` |
| docsify 侧边栏（同步生成） | `server/resource/public/docs/sidebar.md` |

## 使用方式

运行 skill 附带的同步脚本：

```bash
python3 .codebuddy/skills/sync-docs-sidebar/scripts/sync_docs_sidebar.py <project_root>
```

`<project_root>` 为 hotgo 项目根目录的绝对路径。不传参数时脚本会自动向上查找。

## 脚本功能

`scripts/sync_docs_sidebar.py` 执行以下操作：

1. **扫描文件**：列出 `docs/guide-zh-CN/` 下所有 md 文件
2. **发现新文件**：与当前 README 中已列出的文件对比，找出未列出的文件
3. **自动分类**：根据文件名前缀将新文件归入对应分类：
   - `start-*` → 介绍安装
   - `sys-*` / `dev-*` → 系统开发
   - `addon-*` → 插件模块开发
   - `code-*` → 生成代码
   - `web-*` → 前端开发
   - `append-*` → 附录
4. **提取标题**：读取每个新文件的第一个 `#` 或 `##` 标题作为链接显示文本
5. **更新源 README**：将新文件条目插入到 `docs/guide-zh-CN/README.md` 的对应分类下
6. **生成 docsify 首页**：链接路径添加 `guide-zh-CN/` 前缀，移除上级目录链接
7. **生成 docsify 侧边栏**：转换为 docsify 侧边栏格式，`系统介绍` 指向 `README.md`
8. **验证链接**：检查所有引用的文件是否存在

## 分类规则自定义

在脚本中修改 `PREFIX_TO_SECTION` 字典即可调整文件名前缀与分类的映射关系。无法自动分类的文件会输出警告，需手动添加。
