# HotGo V2 开发文档

> 版本：v2.18.6 | 许可证：MIT | 最后更新：2026-02

## 1. 项目概述

**HotGo V2** 是一个基于 GoFrame 2.9.4 + Vue 3.4 + Naive UI + UniApp 开发的全栈开发框架，定位为中小型企业级中后台应用的快速开发骨架，**为二次开发而生**。

### 1.1 核心特性

| 特性 | 说明 |
|------|------|
| 多应用入口 | Admin（后台）、Home（前台页面）、Api（对外通用接口）、WebSocket（即时通讯） |
| 插件化架构 | 微核设计，功能隔离，一键创建/安装/更新/卸载插件，支持多人协同开发 |
| 代码生成 | 无需编写代码，配置数据库表即可生成完善的 CURD、树表、消息队列、定时任务代码 |
| 认证机制 | JWT 用户状态认证 + Casbin RBAC 权限认证 |
| 规范路由 | GoFrame 规范化路由注册，无需注解自动生成 OpenAPI/Swagger 文档 |
| 模块化设计 | 面向接口开发，分层清晰，职责明确 |

### 1.2 技术栈总览

| 层面 | 技术 | 版本 |
|------|------|------|
| **后端框架** | GoFrame | v2.9.4 |
| **编程语言** | Go | >= 1.23 |
| **前端框架** | Vue 3 + TypeScript | 3.4 / 5.5 |
| **UI 组件库** | Naive UI | >= 2.43.1 |
| **构建工具** | Vite | 5.4 |
| **状态管理** | Pinia | 最新 |
| **数据库** | MySQL / PostgreSQL | 5.7+ / 14+ |
| **缓存** | Redis / File / Memory | — |
| **消息队列** | Disk / Redis / RocketMQ / Kafka | — |
| **对象存储** | 本地 / 阿里云OSS / 腾讯云COS / UCloud / 七牛云 / MinIO | — |
| **实时通信** | gorilla/websocket + 自研TCP服务 | — |
| **支付网关** | 支付宝 / 微信支付 / QQ支付 (gopay) | — |
| **链路追踪** | Jaeger | — |
| **容器部署** | Docker + Kubernetes (Kustomize) | — |

### 1.3 内置功能

1. **用户管理** — 系统用户配置与管理
2. **部门管理** — 组织机构树结构，支持数据权限
3. **岗位管理** — 用户职务配置
4. **菜单管理** — 系统菜单、操作权限、按钮权限标识
5. **角色管理** — 菜单权限分配，数据范围权限划分
6. **字典管理** — 枚举字典和自定义方法字典
7. **配置管理** — 系统动态配置参数
8. **操作日志** — 正常/异常操作日志记录与查询
9. **登录日志** — 含登录异常记录
10. **服务日志** — 警告、异常、崩溃日志及堆栈信息
11. **支付网关** — 支付宝/微信/QQ多种支付集成
12. **资金管理** — 充值、退款、提现、资金/积分变动
13. **在线用户** — 活跃用户状态监控
14. **定时任务** — 在线任务调度管理
15. **代码生成** — 前后端 CURD/树表/消息队列/定时任务代码自动生成
16. **插件应用** — 独立多应用入口、独立配置的插件系统
17. **服务监控** — CPU/内存/磁盘/网络/堆栈实时监控
18. **附件管理** — 多驱动文件上传，支持大文件分片/断点续传
19. **TCP 服务** — 长连接/断线重连/路由分发/RPC 消息
20. **消息队列** — 多驱动一键切换
21. **通知公告** — WebSocket 实时推送
22. **地区编码** — 省市区编码及动态选项
23. **常用工具** — 工具包和命令行工具集合

### 1.4 项目根目录结构

```
hotgo/
├── docs/                    # 项目文档
│   └── guide-zh-CN/         # 中文使用指南和开发文档
├── server/                  # Go 后端服务 (697+ .go 文件)
├── web/                     # Vue3 前端项目 (247+ .ts, 203+ .vue 文件)
├── index.html               # 文档站首页
├── sidebar.md               # 文档站侧边栏
├── README.md                # 项目说明
└── LICENSE                  # MIT 开源许可证
```

## 2. 文档导航

本开发文档按模块拆分为以下章节，每个章节独立一个文件：

### 架构与设计

| 文档 | 说明 |
|------|------|
| [技术架构解析](dev-architecture.md) | 系统架构图、多服务命令体系、启动流程、中间件链、路由分组、鉴权流程、数据流向、分层架构 |

### 核心模块

| 文档 | 说明 |
|------|------|
| [服务端核心模块](dev-server-modules.md) | cmd 命令层、controller 控制器、logic 业务逻辑、model 数据模型、service 接口层、dao 数据访问、library 功能库、utility 工具包 |
| [前端核心模块](dev-frontend-modules.md) | 入口启动流程、路由系统、状态管理(Pinia)、API 调用层、页面视图、组件体系、国际化 |

### 接口与扩展

| 文档 | 说明 |
|------|------|
| [插件系统指南](dev-plugin-system.md) | 插件架构设计、Module 接口、目录结构规范、启动加载流程、示例插件分析 |
| [关键接口定义](dev-api-interfaces.md) | Service 层接口、消息队列接口、中间件接口、支付网关接口、存储驱动接口 |

### 配置与部署

| 文档 | 说明 |
|------|------|
| [配置参数参考](dev-config-reference.md) | config.yaml 文件配置 + 数据库动态配置（上传/短信/支付/邮件/微信等）详解 |
| [部署指南](dev-deployment.md) | 环境要求、安装步骤、编译构建、Docker/K8s 部署、Nginx 配置、集群部署 |

### 已有文档索引

| 文档 | 说明 |
|------|------|
| [安装指南](start-installation.md) | 环境准备与项目安装 |
| [部署说明](start-deploy.md) | 生产环境部署 |
| [更新日志](start-update-log.md) | 版本更新历史 |
| [常见问题](start-issue.md) | FAQ |
| [系统目录说明](sys-catalog.md) | 系统目录结构说明 |
| [中间件](sys-middleware.md) | 中间件使用说明 |
| [认证鉴权](sys-auth.md) | JWT + Casbin 认证说明 |
| [消息队列](sys-queue.md) | 消息队列使用说明 |
| [定时任务](sys-cron.md) | 定时任务配置说明 |
| [TCP 服务](sys-tcp-server.md) | TCP 服务开发说明 |
| [WebSocket 服务端](sys-websocket-server.md) | WebSocket 服务端说明 |
| [WebSocket 客户端](sys-websocket-client.md) | WebSocket 客户端说明 |
| [支付系统](sys-payment.md) | 支付网关配置说明 |
| [功能库](sys-library.md) | library 功能库说明 |
| [工具包](sys-utility.md) | utility 工具包说明 |
| [国际化](sys-i18n.md) | i18n 国际化说明 |
| [插件介绍](addon-introduce-catalog.md) | 插件系统基础介绍 |
| [插件开发](addon-flow.md) | 插件开发流程 |
| [代码生成-入门](code-start.md) | 代码生成快速上手 |
| [代码生成-CURD](code-curd.md) | CURD 代码生成 |
| [代码生成-树表](code-tree.md) | 树表代码生成 |

## 3. 快速开始

### 3.1 环境要求

| 依赖 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Go | 1.23+ | 1.24 |
| Node.js | 16+ | 20+ |
| MySQL | 5.7+ | 8.0+ |
| Redis | 4.0+ | 7.0+ |
| pnpm | 8.0+ | 最新 |

### 3.2 安装步骤

```bash
# 1. 克隆项目
git clone https://github.com/bufanyun/hotgo.git

# 2. 导入数据库
#    MySQL: server/storage/data/hotgo.sql
#    PostgreSQL: server/storage/data/hotgo-pg.sql

# 3. 修改后端配置
cd server
cp manifest/config/config.example.yaml manifest/config/config.yaml
# 编辑 config.yaml，修改数据库、Redis 等连接信息

# 4. 启动后端
go run main.go

# 5. 启动前端
cd ../web
pnpm install
pnpm dev
```

### 3.3 默认账号

| 账号 | 密码 | 说明 |
|------|------|------|
| admin | 123456 | 超级管理员 |

### 3.4 访问地址

- 后台管理：`http://localhost:8001/admin`（前端开发服务默认端口 8001）
- API 文档：`http://localhost:8000/swagger`（GoFrame 自动生成）
- 演示地址：[https://hotgo.facms.cn/admin](https://hotgo.facms.cn/admin)

## 4. 实战：新增一个 CRUD 模块（端到端）

本节以添加一个"文章管理"模块为例，演示从数据库建表到前后端联调的完整流程。

### 4.1 方式一：使用代码生成器（推荐）

代码生成器可一键生成前后端全部代码，是最高效的开发方式。

#### 第一步：建表

在数据库中创建业务表（注意表前缀 `hg_`）：

```sql
CREATE TABLE `hg_sys_article` (
  `id`          bigint       NOT NULL AUTO_INCREMENT COMMENT '文章ID',
  `title`       varchar(200) NOT NULL DEFAULT '' COMMENT '标题',
  `content`     text         COMMENT '内容',
  `category_id` bigint       NOT NULL DEFAULT 0 COMMENT '分类ID',
  `author`      varchar(50)  NOT NULL DEFAULT '' COMMENT '作者',
  `views`       int          NOT NULL DEFAULT 0 COMMENT '浏览量',
  `sort`        int          NOT NULL DEFAULT 0 COMMENT '排序',
  `status`      tinyint      NOT NULL DEFAULT 1 COMMENT '状态:1=正常,2=禁用',
  `created_at`  datetime     DEFAULT NULL COMMENT '创建时间',
  `updated_at`  datetime     DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章管理';
```

#### 第二步：在后台配置生成选项

1. 登录后台管理 → 开发工具 → 代码生成
2. 点击 **"新增"** → 选择数据库和表 `hg_sys_article`
3. 系统会自动采集表字段并推断默认配置（表单组件、查询方式、列表显示等）
4. 按需调整：
   - **基本信息**：实体命名填 `Article`，选择模板（如 `admin` 分组）
   - **字段配置**：`content` 字段的表单组件改为 `InputEditor`（富文本）、列表中不显示
   - **生成选项**：勾选 `runDao`（生成 DAO）、`runService`（生成 Service）、`genMenuPermissions`（生成菜单权限）
   - **表头操作**：勾选 `add`（添加）、`batchDel`（批量删除）、`export`（导出）
   - **列操作**：勾选 `edit`（编辑）、`del`（删除）、`view`（查看）、`status`（状态切换）

#### 第三步：预览并提交生成

1. 点击 **"预览"** 查看即将生成的所有文件内容
2. 确认无误后点击 **"提交生成"**

生成器会自动创建以下文件：

**后端（7 个文件）：**

```
server/
├── api/article/article.go                        # API 请求/响应结构体
├── internal/model/input/sysin/article.go          # Input/Output 模型 + 验证器
├── internal/controller/admin/sys/article.go       # 控制器
├── internal/logic/sys/article.go                  # 业务逻辑（List/Edit/Delete/View/Status/Export）
├── internal/router/genrouter/article.go           # 路由自动注册（init() append）
├── internal/dao/sys_article.go                    # DAO（gf gen dao 生成）
└── storage/data/generate/article_menu.sql         # 菜单权限 SQL（已自动导入）
```

**前端（4 个文件）：**

```
web/src/
├── api/article/index.ts                           # API 调用函数
└── views/article/
    ├── model.ts                                   # State 类、rules、schemas、columns
    ├── index.vue                                  # 列表页
    ├── edit.vue                                   # 编辑弹窗
    └── view.vue                                   # 详情抽屉
```

**同时自动执行：**
- `gf gen dao` → 生成 DAO/DO/Entity
- `gf gen service` → 生成 Service 接口
- 菜单权限 SQL 导入数据库 → 后台菜单可见

#### 第四步：分配权限并访问

1. 后台 → 权限管理 → 角色管理 → 编辑角色 → 勾选"文章管理"菜单权限
2. 刷新页面，侧边栏出现新菜单，CRUD 功能即可使用

> **提示**：代码生成后如需自定义业务逻辑（如关联查询、特殊校验），直接修改 `logic/sys/article.go` 即可。

### 4.2 方式二：手动开发

如需更灵活的控制，可手动创建各层文件。以下是核心步骤：

#### 第一步：定义 API 结构体

> 文件：`server/api/article/article.go`
>
> 分层详情参见 [技术架构 — 分层架构](dev-architecture.md#2-分层架构)

```go
package article

import "github.com/gogf/gf/v2/frame/g"

type ListReq struct {
    g.Meta `path:"/article/list" method:"get" tags:"文章管理" summary:"获取文章列表"`
    Title  string `json:"title"`       // 标题模糊搜索
    Status int    `json:"status"`      // 状态筛选
    Page   int    `json:"page" d:"1"`
    PerPage int   `json:"pageSize" d:"10"`
}

type ListRes struct {
    List       []*ArticleModel `json:"list"`
    PageCount  int             `json:"pageCount"`
    TotalCount int             `json:"totalCount"`
}

// ... 其他 Edit/Delete/View/Status 请求响应定义
```

#### 第二步：创建 Input 模型

> 文件：`server/internal/model/input/sysin/article.go`
>
> 命名约定参见 [服务端模块 — Model 层](dev-server-modules.md#6-model-层--数据模型)

```go
type ArticleListInp struct {
    form.PageReq
    Title  string `json:"title"`
    Status int    `json:"status"`
}

type ArticleListModel struct {
    entity.SysArticle
}
```

#### 第三步：实现 Logic 业务逻辑

> 文件：`server/internal/logic/sys/article.go`

```go
func init() {
    service.RegisterSysArticle(New())
}

func (s *sSysArticle) List(ctx context.Context, in *sysin.ArticleListInp) (list []*sysin.ArticleListModel, totalCount int, err error) {
    mod := dao.SysArticle.Ctx(ctx)
    if in.Title != "" {
        mod = mod.WhereLike(dao.SysArticle.Columns().Title, "%"+in.Title+"%")
    }
    if in.Status > 0 {
        mod = mod.Where(dao.SysArticle.Columns().Status, in.Status)
    }
    totalCount, err = mod.Count()
    if err != nil { return }
    err = mod.Page(in.Page, in.PerPage).OrderDesc(dao.SysArticle.Columns().Id).Scan(&list)
    return
}
```

#### 第四步：编写 Controller

> 文件：`server/internal/controller/admin/sys/article.go`
>
> 控制器设计模式参见 [服务端模块 — Controller 层](dev-server-modules.md#3-controller-层--控制器)

```go
var Article = cArticle{}
type cArticle struct{}

func (c *cArticle) List(ctx context.Context, req *article.ListReq) (res *article.ListRes, err error) {
    list, totalCount, err := service.SysArticle().List(ctx, &sysin.ArticleListInp{...})
    // ... 组装响应
}
```

#### 第五步：注册路由

> 文件：`server/internal/router/genrouter/article.go`

```go
package genrouter

import "hotgo/internal/controller/admin/sys"

func init() {
    LoginRequiredRouter = append(LoginRequiredRouter, sys.Article)
}
```

#### 第六步：生成 DAO 和 Service

```bash
cd server

# 生成 DAO/DO/Entity
make dao
# 或: gf gen dao

# 生成 Service 接口
make service
# 或: gf gen service
```

#### 第七步：创建前端页面

前端文件遵循三文件分离模式，放置在 `web/src/views/article/` 目录下：

| 文件 | 说明 |
|------|------|
| `model.ts` | State 类、表单验证规则、搜索 Schema、表格列定义 |
| `index.vue` | 列表页（BasicForm + BasicTable + TableAction） |
| `edit.vue` | 编辑弹窗（BasicModal + 表单） |

> 前端组件使用参见 [前端模块 — 组件系统](dev-frontend-modules.md#8-组件系统)，页面模式参见 [前端模块 — 页面视图](dev-frontend-modules.md#7-页面视图)

#### 第八步：在数据库添加菜单

在 `hg_admin_menu` 表中插入菜单记录，`component` 字段填写 `/article/index`。或通过后台 → 菜单管理手动添加。前端路由会自动从后端获取，无需修改路由文件。

### 4.3 开发新插件的前端页面

如果是为插件开发页面，遵循以下约定：

| 约定项 | 规则 | 示例 |
|--------|------|------|
| 视图目录 | `web/src/views/addons/{插件名}/` | `views/addons/myplugin/table/index.vue` |
| API 目录 | `web/src/api/addons/{插件名}/` | `api/addons/myplugin/table/index.ts` |
| API URL | `/{插件名}/{模块}/{方法}` | `/myplugin/table/list` |
| 菜单 component | `/addons/{插件名}/{路径}` | `/addons/myplugin/table/index` |

前端通过 `import.meta.glob('../views/**/*.{vue,tsx}')` 自动发现所有 `.vue` 文件，只要文件放在正确路径，后端菜单的 `component` 字段与之匹配即可，**无需修改任何路由文件**。

> 后端插件开发详情参见 [插件系统指南](dev-plugin-system.md)

## 5. 开源协议

HotGo 遵循 [MIT 开源协议](../../LICENSE)，可自由使用于商业项目，需保留版权信息。

---

*本文档由项目代码分析自动生成，如有疑问请参考源代码或联系开发团队。*
