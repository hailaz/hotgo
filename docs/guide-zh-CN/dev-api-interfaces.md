# 关键接口定义

> 本文档列出 HotGo 核心组件的接口定义，便于开发者理解系统扩展点。
>
> **文档导航**：[← 返回主文档](dev-guide.md) | [技术架构](dev-architecture.md) | [服务端模块](dev-server-modules.md) | [前端模块](dev-frontend-modules.md) | [插件系统 →](dev-plugin-system.md) | [配置参考](dev-config-reference.md) | [部署指南](dev-deployment.md)

## 1. 中间件接口 (IMiddleware)

> 文件：`server/internal/service/middleware.go`

```go
type IMiddleware interface {
    // 全局中间件
    Ctx(r *ghttp.Request)             // 上下文初始化（必须第一个）
    CORS(r *ghttp.Request)            // 跨域资源共享
    Blacklist(r *ghttp.Request)       // IP 黑名单
    DemoLimit(r *ghttp.Request)       // 演示系统操作限制
    PreFilter(r *ghttp.Request)       // 请求输入预处理
    ResponseHandler(r *ghttp.Request) // HTTP 响应处理

    // 路由级鉴权中间件
    AdminAuth(r *ghttp.Request)       // 后台鉴权（Token + Casbin）
    ApiAuth(r *ghttp.Request)         // API 鉴权（Token）
    HomeAuth(r *ghttp.Request)        // 前台鉴权（预留）
    WebSocketAuth(r *ghttp.Request)   // WebSocket 鉴权

    // 功能中间件
    Addon(r *ghttp.Request)           // 插件中间件（解析插件名到上下文）
    Develop(r *ghttp.Request)         // 开发工具 IP 白名单

    // 辅助方法
    DeliverUserContext(r *ghttp.Request) error   // 传递用户信息到上下文
    IsExceptAuth(ctx context.Context, appName, path string) bool  // 是否免权限路由
    IsExceptLogin(ctx context.Context, appName, path string) bool // 是否免登录路由
}
```

## 2. 插件模块接口 (Module)

> 文件：`server/internal/library/addons/module.go`

```go
type Module interface {
    Start(option *Option) (err error)          // 启动模块（注册路由等）
    Stop() (err error)                         // 停止模块
    Ctx() context.Context                      // 获取模块上下文
    GetSkeleton() *Skeleton                    // 获取骨架信息
    Install(ctx context.Context) (err error)   // 安装回调（建表、初始化数据等）
    Upgrade(ctx context.Context) (err error)   // 升级回调（表结构变更等）
    UnInstall(ctx context.Context) (err error) // 卸载回调（清理数据等）
}
```

## 3. 消息队列接口

### 3.1 消费者接口 (Consumer)

> 文件：`server/internal/library/queue/consumer.go`

```go
type Consumer interface {
    GetTopic() string                                           // 获取订阅主题
    Handle(ctx context.Context, mqMsg MqMsg) (err error)        // 消息处理
}
```

### 3.2 生产者接口 (MqProducer)

> 文件：`server/internal/library/queue/queue.go`

```go
type MqProducer interface {
    SendMsg(topic string, body string) (mqMsg MqMsg, err error) // 发送消息
    SendByteMsg(topic string, body []byte) (mqMsg MqMsg, err error) // 发送字节消息
    ListenReceiveMsgDo(topic string, received func(mqMsg MqMsg)) (err error) // 监听消费
}
```

### 3.3 消息体

```go
type MqMsg struct {
    RunType   string // 驱动类型
    Topic     string // 主题
    MsgId     string // 消息ID
    Body      []byte // 消息内容
    Timestamp int64  // 发送时间
}
```

### 3.4 全局队列 API

```go
// 推送消息到队列
queue.Push(topic string, body interface{}) error

// 注册消费者
queue.RegisterConsumer(consumer Consumer)

// 启动消费者监听
queue.StartConsumersListener(ctx context.Context)
```

## 4. 支付网关接口

### 4.1 PayClient 接口

> 文件：`server/internal/library/payment/payment.go`

```go
type PayClient interface {
    CreateOrder(ctx context.Context, in payin.CreateOrderInp) (res *payin.CreateOrderModel, err error)
    Notify(ctx context.Context, in payin.NotifyInp) (res *payin.NotifyModel, err error)
    Refund(ctx context.Context, in payin.RefundInp) (res *payin.RefundModel, err error)
}
```

支持的实现：
- `aliPayClient` — 支付宝
- `wxPayClient` — 微信支付
- `qqPayClient` — QQ 支付

### 4.2 NotifyCall 支付回调

```go
type NotifyCallFunc func(ctx context.Context, in NotifyCallFuncInp) (err error)

type NotifyCallFuncInp struct {
    Pay *payin.NotifyModel  // 支付通知数据
}
```

通过 `RegisterNotifyCall(key, func)` 注册回调，支付成功后自动触发。

## 5. 文件存储接口

### 5.1 UploadDrive 接口

> 文件：`server/internal/library/storager/upload.go`

```go
type UploadDrive interface {
    Upload(ctx context.Context, file *ghttp.UploadFile) (fullPath string, err error)
    CreateMultipart(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error)
    UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error)
}
```

支持的实现（6 种驱动）：

| 驱动 | 说明 |
|------|------|
| `LocalDrive` | 本地文件存储 |
| `UCloudDrive` | UCloud 对象存储 |
| `CosDrive` | 腾讯云 COS |
| `OssDrive` | 阿里云 OSS |
| `QiNiuDrive` | 七牛云 |
| `MinioDrive` | MinIO 对象存储 |

## 6. 定时任务接口

> 文件：`server/internal/library/cron/cron.go`

```go
type Cron interface {
    GetName() string                                            // 获取任务名称
    Execute(ctx context.Context, parser *Parser) (err error)    // 执行任务
}
```

**Parser** 提供参数解析：

```go
type Parser struct {
    pattern string  // cron 表达式
    policy  int64   // 策略：1-并行 2-单例 3-单次 4-多次
    params  string  // JSON 参数
}
```

**全局 API：**

```go
cron.Register(c Cron)                // 注册任务
cron.StartALL(SysCrons)              // 启动所有任务
cron.StopALL()                       // 停止所有任务
cron.Once(SysCron)                   // 执行一次
cron.Delete(SysCron)                 // 删除任务
cron.RefreshStatus(id, status)       // 刷新状态
```

## 7. 缓存适配器接口

> 文件：`server/internal/library/cache/cache.go`

```go
// 根据 cache.adapter 配置选择驱动
func SetAdapter(ctx context.Context)

// 获取 gcache.Cache 实例
func Instance() *gcache.Cache
```

支持的驱动：
- `memory` — 内存缓存
- `redis` — Redis 缓存
- `file` — 文件缓存

## 8. Token 认证接口

> 文件：`server/internal/library/token/token.go`

```go
// 用户登录，生成 Token
func Login(ctx context.Context, app string, identity *model.Identity) (token string, expire int64, err error)

// 解析登录用户
func ParseLoginUser(r *ghttp.Request) (*model.Identity, error)

// 用户注销
func Logout(ctx context.Context, r *ghttp.Request) error
```

## 9. 核心 Service 接口示例

### 9.1 ISysConfig — 系统配置

```go
type ISysConfig interface {
    InitConfig(ctx context.Context)                                    // 初始化系统配置
    LoadConfig(ctx context.Context)                                    // 加载配置
    GetBasic(ctx context.Context) (conf *model.BasicConfig, err error) // 基础配置
    GetLogin(ctx context.Context) (conf *model.LoginConfig, err error) // 登录配置
    GetPay(ctx context.Context) (conf *model.PayConfig, err error)     // 支付配置
    GetSms(ctx context.Context) (conf *model.SmsConfig, err error)     // 短信配置
    GetUpload(ctx context.Context) (conf *model.UploadConfig, err error) // 上传配置
    GetSmtp(ctx context.Context) (conf *model.EmailConfig, err error)  // 邮箱配置
    GetToken(ctx context.Context) (conf *model.TokenConfig, err error) // Token 配置
    GetConfigByGroup(ctx context.Context, in *sysin.GetConfigInp) (res *sysin.GetConfigModel, err error)
    UpdateConfigByGroup(ctx context.Context, in *sysin.UpdateConfigInp) (err error)
    ClusterSync(ctx context.Context, message *gredis.Message)          // 集群同步
    // ...
}
```

### 9.2 IAdminMember — 用户管理

```go
type IAdminMember interface {
    // 余额/积分操作
    AddBalance(ctx context.Context, in *adminin.MemberAddBalanceInp) (err error)
    AddIntegral(ctx context.Context, in *adminin.MemberAddIntegralInp) (err error)

    // 信息更新
    UpdateCash(ctx context.Context, in *adminin.MemberUpdateCashInp) (err error)
    UpdateProfile(ctx context.Context, in *adminin.MemberUpdateProfileInp) (err error)
    UpdatePwd(ctx context.Context, in *adminin.MemberUpdatePwdInp) (err error)
    ResetPwd(ctx context.Context, in *adminin.MemberResetPwdInp) (err error)

    // CRUD
    Edit(ctx context.Context, in *adminin.MemberEditInp) (err error)
    Delete(ctx context.Context, in *adminin.MemberDeleteInp) (err error)
    View(ctx context.Context, in *adminin.MemberViewInp) (res *adminin.MemberViewModel, err error)
    List(ctx context.Context, in *adminin.MemberListInp) (list []*adminin.MemberListModel, totalCount int, err error)
    Status(ctx context.Context, in *adminin.MemberStatusInp) (err error)

    // 查询辅助
    LoginMemberInfo(ctx context.Context) (res *adminin.LoginMemberInfoModel, err error)
    GetIdByCode(ctx context.Context, in *adminin.GetIdByCodeInp) (res *adminin.GetIdByCodeModel, err error)
    Select(ctx context.Context, in *adminin.MemberSelectInp) (res []*adminin.MemberSelectModel, err error)

    // 权限相关
    VerifySuperId(ctx context.Context, id int64) bool
    LoadSuperAdmin(ctx context.Context)
    FilterAuthModel(ctx context.Context, ...) *gdb.Model
    // ...
}
```

### 9.3 IAdminRole — 角色管理

```go
type IAdminRole interface {
    Verify(ctx context.Context, path, method string) bool                    // 权限验证
    GetPermissions(ctx context.Context, in *adminin.GetPermissionsInp) (res adminin.GetPermissionsModel, err error)
    UpdatePermissions(ctx context.Context, in *adminin.UpdatePermissionsInp) (err error)
    Edit(ctx context.Context, in *adminin.RoleEditInp) (err error)
    Delete(ctx context.Context, in *adminin.RoleDeleteInp) (err error)
    List(ctx context.Context, in *adminin.RoleListInp) (list []*adminin.RoleListModel, totalCount int, err error)
    DataScopeEdit(ctx context.Context, in *adminin.DataScopeEditInp) (err error)
    VerifyRoleId(ctx context.Context, id int64) (err error)
    GetSubRoleIds(ctx context.Context, roleId int64) (ids []int64, err error)
    // ...
}
```

## 10. TCP 通信接口

### 10.1 ITCPServer — TCP 服务端

```go
type ITCPServer interface {
    Instance() *tcp.Server
    Start(ctx context.Context)
    Stop(ctx context.Context)
    DefaultInterceptor(ctx context.Context, conn *tcp.ServerConn, msg *tcp.Message) error
    PreFilterInterceptor(ctx context.Context, conn *tcp.ServerConn, msg *tcp.Message) error
    OnAuthSummary(ctx context.Context, ...)
    CronDelete/Edit/Status/OnlineExec(ctx context.Context, ...)
    DispatchLog(ctx context.Context, ...)
    // ...
}
```

### 10.2 ICronClient — 定时任务客户端

```go
type ICronClient interface {
    Instance() *tcp.Client
    Start(ctx context.Context)
    Stop(ctx context.Context)
    DefaultInterceptor(ctx context.Context, conn *tcp.ClientConn, msg *tcp.Message) error
    OnCronDelete/Edit/Status/OnlineExec(ctx context.Context, ...)
    // ...
}
```

---

