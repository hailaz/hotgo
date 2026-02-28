# 服务端核心模块详解

> 本文档详细说明 HotGo 服务端各层模块的职责、关键文件、接口定义和设计模式。
>
> **文档导航**：[← 返回主文档](dev-guide.md) | [技术架构](dev-architecture.md) | [前端模块 →](dev-frontend-modules.md) | [接口定义](dev-api-interfaces.md) | [插件系统](dev-plugin-system.md) | [配置参考](dev-config-reference.md) | [部署指南](dev-deployment.md)

## 1. 整体模块结构

```
server/internal/
├── cmd/           # 命令入口层（8 个命令）
├── consts/        # 常量定义（28 个文件）
├── controller/    # 控制器层（4 个应用入口）
├── crons/         # 定时任务实现（3 个任务）
├── dao/           # 数据访问层（39+39 个文件）
├── global/        # 全局初始化（4 个文件）
├── library/       # 功能库（22 个模块）
├── logic/         # 业务逻辑层（66 个文件，9 个子包）
├── model/         # 数据模型层（entity/do/input）
├── queues/        # 消息队列消费者（3 个消费者）
├── router/        # 路由注册
├── service/       # 服务接口层（桥接接口 + 保留接口）
└── websocket/     # WebSocket 核心
```

## 2. Service 层 — 接口定义

> **v3.0 架构变更**：Service 层已大幅简化。不再使用 `gf gen service` 为所有模块自动生成接口，仅保留真正需要多态/解耦的接口和用于打破循环依赖的桥接接口。绝大多数 CRUD 模块由 Controller 直接调用 Logic 层。

### 2.1 设计原则

- **CRUD 模块不经过 Service**：Controller 通过 `xxxLogic.XxxFunc()` 直接调用 Logic 层导出函数
- **保留接口**：中间件（IMiddleware）、Hook（IHook）、TCP 服务端/客户端（ITCPServer/IAuthClient/ICronClient）、视图（IView）等需要多态注册的接口仍保留在 Service 层
- **桥接接口**：当存在循环依赖（如 `logic/sys` ↔ `logic/admin`、`global` → `logic/*`）时，通过最小化的桥接接口（`admin_bridge.go`、`sys_bridge.go`）解耦，仅暴露被跨包引用的少量方法

### 2.2 保留的接口文件

| 文件 | 接口 | 职责 |
|------|------|------|
| `middleware.go` | `IMiddleware` (15) | 中间件集合 |
| `hook.go` | `IHook` (2) | HTTP 钩子：BeforeServe, AfterOutput |
| `tcpserver.go` | `ITCPServer` (12) | TCP 服务端 |
| `tcpclient.go` | `IAuthClient` (6), `ICronClient` (9) | 认证客户端、定时任务客户端 |
| `view.go` | `IView` (5) | 前台模板渲染 |

### 2.3 桥接接口文件

**`admin_bridge.go`** — 供 `global`、`queues`、`crons`、`logic/sys`、`logic/pay`、`hggen` 等包跨包调用：

| 接口 | 方法数 | 说明 |
|------|--------|------|
| `IAdminMember` | 5 | LoadSuperAdmin, ClusterSyncSuperAdmin, VerifySuperId, GetComplexMemberIds, GetIdsByKeyword |
| `IAdminOrder` | 2 | PayNotify, Create |
| `IAdminMenu` | 2 | Model, GetMenuList |

**`sys_bridge.go`** — 供 `global`、`queues`、`hggen`、`logic/hook` 等包跨包调用：

| 接口 | 方法数 | 说明 |
|------|--------|------|
| `ISysConfig` | 4 | InitConfig, GetBasic, GetLoadGenerate, ClusterSync |
| `ISysBlacklist` | 3 | Load, VerifyRequest, ClusterSync |
| `ISysDictType` | 1 | BuiltinSelect |
| `ISysServeLog` | 1 | RealWrite |
| `ISysLog` | 2 | RealWrite, AnalysisLog |
| `ISysLoginLog` | 1 | RealWrite |
| `ISysAddonsConfig` | 2 | GetConfigByGroup, UpdateConfigByGroup |

### 2.4 已删除的接口文件

以下文件已删除，原有接口的实现直接通过 Logic 层导出函数访问：

| 原文件 | 原接口数 | 说明 |
|--------|----------|------|
| `admin.go` | 12 | AdminCash, AdminCreditsLog, AdminDept, AdminMember 等全部 CRUD 接口 |
| `sys.go` | 21 | SysAddons, SysAttachment, SysCron, SysGenCodes 等全部 CRUD 接口 |
| `pay.go` | 2 | Pay, PayRefund |
| `common.go` | 2 | CommonUpload, CommonWechat |

## 3. Controller 层 — 控制器

### 3.1 目录结构

```
controller/
├── admin/
│   ├── admin/     # 后台管理控制器（10 个文件）
│   │   ├── cash.go, credits_log.go, dept.go, member.go, menu.go
│   │   ├── monitor.go, notice.go, order.go, post.go, role.go
│   ├── sys/       # 系统管理控制器（20 个文件）
│   │   ├── addons.go, attachment.go, blacklist.go, config.go
│   │   ├── cron.go, cron_group.go, curd_demo.go, dict_data.go
│   │   ├── dict_type.go, ems_log.go, gen_codes.go, log.go
│   │   ├── login_log.go, normal_tree_demo.go, option_tree_demo.go
│   │   ├── provinces.go, serve_license.go, serve_log.go
│   │   ├── sms_log.go, test_category.go
│   ├── common/    # 公共控制器（6 个文件）
│   │   ├── console.go, ems.go, site.go, sms.go, upload.go, wechat.go
│   └── pay/       # 支付控制器
├── api/           # 前台 API 控制器
│   ├── member/    # 用户接口（GoFrame 规范路由模式）
│   └── pay/       # 支付回调
├── home/          # 前台页面控制器
└── websocket/     # WebSocket 控制器
```

### 3.2 设计模式

**Admin 控制器** — 结构体变量 + 方法（直接调用 Logic 层）：

```go
import adminLogic "hotgo/internal/logic/admin"

var Member = cMember{}
type cMember struct{}

func (c *cMember) List(ctx context.Context, req *member.ListReq) (res *member.ListRes, err error) {
    list, totalCount, err := adminLogic.AdminMember().List(ctx, &req.MemberListInp)
    if err != nil {
        return
    }
    res = new(member.ListRes)
    res.List = list
    res.PageRes.Pack(req, totalCount)
    return
}
```

**API 控制器** — GoFrame 规范路由模式：

```go
type ControllerV1 struct{}
func NewV1() member.IMemberV1 { return &ControllerV1{} }

func (c *ControllerV1) GetIdByCode(ctx context.Context, req *v1.GetIdByCodeReq) (res *v1.GetIdByCodeRes, err error) {
    // ...
}
```

**核心原则：** 控制器只做**参数拆装**，不含任何业务逻辑。

## 4. Logic 层 — 业务逻辑

### 4.1 目录结构（66 个文件，9 个子包）

| 子包 | 文件数 | 说明 |
|------|--------|------|
| `admin/` | 13 | 后台管理逻辑（member, role, dept, menu, notice, order 等） |
| `sys/` | 21 | 系统管理逻辑（config, cron, dict, genCodes, log 等） |
| `common/` | 2 | 公共逻辑（upload, wechat） |
| `hook/` | 3 | HTTP 钩子逻辑 |
| `middleware/` | 8 | 中间件逻辑 |
| `pay/` | 4 | 支付逻辑 |
| `tcpclient/` | 7 | TCP 客户端逻辑 |
| `tcpserver/` | 6 | TCP 服务端逻辑 |
| `view/` | 2 | 视图渲染逻辑 |

### 4.2 设计模式

每个 Logic 文件采用 **单例导出模式**（不再使用 `gf gen service` 和 init 注册）：

```go
type sAdminMember struct {
    superAdmin *SuperAdmin
}

func NewAdminMember() *sAdminMember {
    return &sAdminMember{superAdmin: new(SuperAdmin)}
}

var insAdminMember = NewAdminMember()

func AdminMember() *sAdminMember {
    return insAdminMember
}
```

> **注意**：需要被桥接接口引用的 Logic 文件（如 `member.go`、`config.go`）仍需保留 `init()` 注册到 `service.RegisterXxx()`，供跨循环依赖的包调用。

**聚合引入**（`logic/logic.go`）：

```go
import (
    _ "hotgo/internal/logic/admin"
    _ "hotgo/internal/logic/common"
    _ "hotgo/internal/logic/hook"
    _ "hotgo/internal/logic/middleware"
    _ "hotgo/internal/logic/pay"
    _ "hotgo/internal/logic/sys"
    _ "hotgo/internal/logic/tcpclient"
    _ "hotgo/internal/logic/tcpserver"
    _ "hotgo/internal/logic/view"
)
```

### 4.3 核心特性

| 特性 | 说明 |
|------|------|
| DAO 直接调用 | `dao.AdminMember.Ctx(ctx).WherePri(id).Scan(&res)` |
| 同包调用 | 同一 logic 子包内直接调用：`AdminRole().VerifyRoleId(ctx, in.RoleId)` |
| 跨包调用 | 通过 import 别名：`sysLogic.SysConfig().GetBasic(ctx)` |
| 桥接调用 | 存在循环依赖时通过桥接接口：`service.AdminMember().GetIdsByKeyword(ctx, ks)` |
| 事务处理 | `g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error { ... })` |
| 权限过滤 | `FilterAuthModel` 方法实现基于角色的数据权限控制 |
| 集群同步 | 通过 Redis PubSub 的 `ClusterSync` 方法 |

## 5. DAO 层 — 数据访问

### 5.1 两层嵌入设计

**外层 DAO**（开发者可扩展）：

```go
type adminMemberDao struct {
    *internal.AdminMemberDao    // 嵌入自动生成的内部 DAO
}

var AdminMember = adminMemberDao{internal.NewAdminMemberDao()}
```

**内层 internal DAO**（GoFrame CLI 自动生成，不可编辑）：

```go
type AdminMemberDao struct {
    table    string                // 表名
    group    string                // 数据库组
    columns  AdminMemberColumns    // 类型安全的列名常量
    handlers []gdb.ModelHandler    // 模型处理器
}
```

### 5.2 数据表清单（39 个 DAO）

| 前缀 | 数量 | 表名 |
|------|------|------|
| `admin_` | 15 | cash, credits_log, dept, member, member_post, member_role, menu, notice, notice_read, oauth, order, post, role, role_casbin, role_menu |
| `sys_` | 19 | addons_config, addons_install, attachment, blacklist, config, cron, cron_group, dict_data, dict_type, ems_log, gen_codes, gen_curd_demo, gen_tree_demo, log, login_log, provinces, serve_license, serve_log, sms_log |
| `pay_` | 2 | log, refund |
| `test_` | 1 | category |
| `addon_` | 2 | hgexample_table, hgexample_tenant_order |

## 6. Model 层 — 数据模型

### 6.1 entity/ — 实体模型（39 个文件）

与数据库表一一对应，由 GoFrame CLI 自动生成：

```go
type AdminMember struct {
    Id           int64       `json:"id"           orm:"id"            description:"管理员ID"`
    DeptId       int64       `json:"deptId"       orm:"dept_id"       description:"部门ID"`
    RoleId       int64       `json:"roleId"       orm:"role_id"       description:"角色ID"`
    Username     string      `json:"username"     orm:"username"      description:"帐号"`
    PasswordHash string      `json:"passwordHash" orm:"password_hash" description:"密码"`
    Balance      float64     `json:"balance"      orm:"balance"       description:"余额"`
    // ... 共 27 个字段
}
```

### 6.2 do/ — 数据操作模型（39 个文件）

用于 DAO 的 Where/Data 操作，所有字段均为 `any` 类型：

```go
type AdminMember struct {
    g.Meta `orm:"table:hg_admin_member, do:true"`
    Id     any  // 管理员ID
    DeptId any  // 部门ID
    // ...
}
```

### 6.3 input/ — 业务输入/输出模型

| 子包 | 文件数 | 说明 |
|------|--------|------|
| `adminin/` | 10 | 后台管理（member, role, dept, menu, notice, order 等） |
| `sysin/` | 21 | 系统管理（config, cron, dict, genCodes, log 等） |
| `payin/` | 3 | 支付（pay, refund） |
| `form/` | 4 | 通用表单（PageReq/PageRes, Select, Sorter, Base） |
| `commonin/` | 1 | 公共（wechat） |
| `servmsgin/` | 2 | 服务消息（auth, example） |
| `websocketin/` | 1 | WebSocket 消息 |

**命名约定 — Inp 输入 + Model 输出：**

```go
// 输入参数
type MemberListInp struct {
    form.PageReq                        // 嵌入分页请求
    form.StatusReq                      // 嵌入状态筛选
    RoleId    int     `json:"roleId"`
    Username  string  `json:"username"`
    CreatedAt []int64 `json:"createdAt"`
}

// 输出结果
type MemberListModel struct {
    entity.AdminMember                  // 嵌入实体
    DeptName string `json:"deptName"`   // 扩展字段
    RoleName string `json:"roleName"`
    PostIds  []int64 `json:"postIds"`
}
```

**通用分页模型：**

```go
type PageReq struct {
    Page    int `json:"page"     d:"1"  v:"page@min:1"`
    PerPage int `json:"pageSize" d:"10" v:"pageSize@min:1|max:200"`
}

type PageRes struct {
    PageReq
    PageCount  int `json:"pageCount"`
    TotalCount int `json:"totalCount"`
}
```

**Filter 预处理接口：** Input 结构体可实现 `Filter(ctx) error` 接口进行自定义参数预处理/校验。

## 7. Library — 功能库（22 个模块）

> 各模块的详细接口定义参见 [关键接口定义](dev-api-interfaces.md)

### 7.1 模块清单

| 模块 | 文件数 | 核心功能 |
|------|--------|----------|
| **addons/** | 5 | 插件模块系统：注册/启动/停止，Module 接口，Skeleton 骨架定义 |
| **cache/** | 2 | 缓存适配器：支持 Redis/File/Memory 三种驱动 |
| **captcha/** | — | 验证码生成（登录用） |
| **casbin/** | — | 基于 Casbin 的 RBAC 权限控制 |
| **contexts/** | 2 | 上下文管理：Init/Get/SetUser 等，请求生命周期数据管理 |
| **cron/** | 2 | 定时任务引擎：Register/StartALL/StopALL，支持并行/单例/单次/多次策略 |
| **debris/** | — | 碎片化工具 |
| **dict/** | — | 字典管理 |
| **ems/** | — | 邮件发送 |
| **hggen/** | — | HotGo 代码生成器 |
| **hgorm/** | 14 | ORM 增强：LeftJoin/GenJoinSelect/IsUnique，handler（FilterAuth/ForceCache/Sorter/Tenant），hook（MemberInfo/Provinces/Tenant），树形 DAO |
| **hgrds/** | — | Redis 增强工具 |
| **location/** | — | 地理位置服务 |
| **network/** | 13 | TCP 网络通信：Client/Server，消息解析器，路由系统，RPC 支持 |
| **payment/** | 9 | 支付系统：PayClient 接口，支付宝/微信/QQ 支付，NotifyCall 回调 |
| **queue/** | 11 | 消息队列系统：MqProducer/MqConsumer 接口，4 种驱动 |
| **response/** | — | 响应处理 |
| **sms/** | — | 短信发送 |
| **storager/** | 11 | 文件存储系统：UploadDrive 接口，6 种驱动，分片上传 |
| **token/** | 1 | JWT Token 管理：Login/Logout/ParseLoginUser |
| **validate/** | — | 验证工具 |
| **wechat/** | — | 微信 SDK 封装 |

### 7.2 关键接口定义

**消息队列消费者接口：**

```go
type Consumer interface {
    GetTopic() string
    Handle(ctx context.Context, mqMsg MqMsg) (err error)
}
```

**支付客户端接口：**

```go
type PayClient interface {
    CreateOrder(ctx context.Context, in payin.CreateOrderInp) (res *payin.CreateOrderModel, err error)
    Notify(ctx context.Context, in payin.NotifyInp) (res *payin.NotifyModel, err error)
    Refund(ctx context.Context, in payin.RefundInp) (res *payin.RefundModel, err error)
}
```

**存储驱动接口：**

```go
type UploadDrive interface {
    Upload(ctx context.Context, file *ghttp.UploadFile) (fullPath string, err error)
    CreateMultipart(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error)
    UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error)
}
```

**定时任务接口：**

```go
type Cron interface {
    GetName() string
    Execute(ctx context.Context, parser *Parser) (err error)
}
```

**插件模块接口：**

```go
type Module interface {
    Start(option *Option) (err error)
    Stop() (err error)
    Ctx() context.Context
    GetSkeleton() *Skeleton
    Install(ctx context.Context) (err error)
    Upgrade(ctx context.Context) (err error)
    UnInstall(ctx context.Context) (err error)
}
```

## 8. Utility — 工具包（12 个子包）

| 包名 | 功能 |
|------|------|
| `charset/` | 字符集处理、随机字节生成 |
| `convert/` | 类型转换、实体字段标签提取、切片去重 |
| `db/` | 数据库工具（引号字符获取、SQLite 字段处理） |
| `encrypt/` | 加密工具（AES ECB 加解密） |
| `excel/` | Excel 导出工具 |
| `file/` | 文件操作工具 |
| `format/` | 格式化工具（文件大小、时间等） |
| `simple/` | **通用工具集**：SafeGo（安全协程）、FilterMaskDemo（演示敏感数据过滤）、DecryptText/CheckPassword（密码处理）、GetHeaderLocale（国际化）、RouterPrefix |
| `tree/` | **树形结构工具**：GenLabel/GetIds/GenTree/GenTreeWithField |
| `url/` | URL 处理工具 |
| `useragent/` | UserAgent 解析 |
| `validate/` | **验证工具集**：IsEmail/IsMobile/IsLocalIPAddr/IsMobileVisit/IsWxBrowserVisit，Filter 接口 |

## 9. Crons — 定时任务

共 **3 个定时任务**：

| 文件 | 任务名 | 功能 |
|------|--------|------|
| `close_order.go` | `close_order` | 取消创建超过 1 天且未支付的过期订单 |
| `test.go` | `test` | 测试任务（无参数），输出当前时间 |
| `test2.go` | `test2` | 测试任务（带参数），接收 name/age/msg 参数 |

**注册模式：**

```go
func init() {
    cron.Register(CloseOrder)
}

var CloseOrder = &cCloseOrder{name: "close_order"}

func (c *cCloseOrder) GetName() string { return c.name }

func (c *cCloseOrder) Execute(ctx context.Context, parser *cron.Parser) (err error) {
    // 业务逻辑
}
```

## 10. Queues — 消息队列消费者

共 **3 个消费者**：

| 文件 | Topic | 功能 |
|------|-------|------|
| `login_log.go` | `QueueLoginLogTopic` | 登录日志异步写入 |
| `serve_log.go` | `QueueServeLogTopic` | 服务日志异步写入 |
| `sys_log.go` | `QueueLogTopic` | 系统请求日志异步写入 |

**注册模式：**

```go
func init() {
    queue.RegisterConsumer(LoginLog)
}

var LoginLog = &qLoginLog{}

func (q *qLoginLog) GetTopic() string { return consts.QueueLoginLogTopic }

func (q *qLoginLog) Handle(ctx context.Context, mqMsg queue.MqMsg) (err error) {
    var data entity.SysLoginLog
    json.Unmarshal(mqMsg.Body, &data)
    return service.SysLoginLog().RealWrite(ctx, data)  // 通过桥接接口调用
}
```

## 11. 各层设计模式总结

| 层级 | 文件数 | 核心设计模式 |
|------|--------|-------------|
| **Service** | 7 | 保留接口（Middleware/Hook/TCP/View）+ 桥接接口（admin_bridge/sys_bridge）用于打破循环依赖 |
| **Controller** | 45 | Admin 用结构体变量模式，API 用 GoFrame 规范路由；只做参数拆装；直接调用 Logic 层 |
| **Logic** | 66 | 单例导出模式；直接操作 DAO；同包直接调用，跨包通过 import 别名或桥接接口 |
| **DAO** | 39(外)+39(内) | 两层嵌入（外层可扩展，内层自动生成）；类型安全列名常量 |
| **Model** | 120+ | entity 映射数据库表；do 用于 DAO 操作；input 定义 Inp/Model |
| **Library** | 22 模块 | 策略模式（多驱动）+ 注册模式（cron/queue/addons） |
| **Crons** | 3 | 实现 Cron 接口 → init() 注册 |
| **Queues** | 3 | 实现 Consumer 接口 → init() 注册 |

---

*本文档基于 HotGo v2.18.6 代码库分析生成。*
