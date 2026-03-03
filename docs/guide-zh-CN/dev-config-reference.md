# 配置参数参考手册

> 本文档详解 HotGo 的所有配置项，包括 YAML 文件配置和数据库动态配置。
>
> **文档导航**：[← 返回主文档](dev-guide.md) | [技术架构](dev-architecture.md) | [服务端模块](dev-server-modules.md) | [前端模块](dev-frontend-modules.md) | [接口定义](dev-api-interfaces.md) | [插件系统](dev-plugin-system.md) | [部署指南 →](dev-deployment.md)

## 一、本地文件配置 (`config.yaml`)

> 以下配置项均存储在 `server/manifest/config/config.yaml` 文件中，应用启动时读取。

## 1. system — 系统配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `appName` | string | `"hotgo"` | 应用名称 |
| `debug` | bool | `true` | 调试模式（生产环境应设为 false） |
| `mode` | string | `"develop"` | 运行模式：`not-set`/`develop`/`testing`/`staging`/`product` |
| `ipMethod` | string | `"whois"` | IP 解析方法 |
| `isDemo` | bool | `false` | 演示模式（禁止 POST 操作） |
| `isCluster` | bool | `false` | 集群部署（需配合 Redis PubSub） |
| `addonsResourcePath` | string | `"resource"` | 插件资源根目录 |
| `log.switch` | bool | `true` | 请求日志开关 |
| `log.queue` | bool | `true` | 日志异步写入（通过消息队列） |
| `log.module` | string[] | `["admin","api"]` | 需要记录日志的模块 |
| `log.skipCode` | int[] | `[-1]` | 跳过记录的状态码 |
| `serveLog.switch` | bool | `true` | 服务日志开关 |
| `serveLog.queue` | bool | `true` | 服务日志异步写入 |
| `serveLog.levelFormat` | string[] | `["WARN","ERRO","FATA","PANI"]` | 记录的日志级别 |
| `i18n.switch` | bool | `true` | 国际化开关 |
| `i18n.defaultLanguage` | string | `"zh-CN"` | 默认语言 |

## 2. server — HTTP 服务配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `address` | string | `":8000"` | 监听地址 |
| `serverRoot` | string | `"resource/public"` | 静态文件根目录 |
| `openapiPath` | string | `"/api.json"` | OpenAPI 文档路径 |
| `swaggerPath` | string | `"/swagger"` | Swagger UI 路径 |
| `pprofEnabled` | bool | `false` | PProf 性能分析开关 |
| `clientMaxBodySize` | string | `"200MB"` | 最大请求体大小 |
| `logPath` | string | `"./storage/log/server"` | 日志文件路径 |

## 3. tcp — TCP 服务配置

### 3.1 TCP 服务端

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `server.address` | string | `":8099"` | 监听地址 |

### 3.2 TCP 客户端

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `client.cron.address` | string | `"127.0.0.1:8099"` | 定时任务客户端连接地址 |
| `client.cron.appName` | string | `"hotgo"` | 应用名 |
| `client.cron.loginTimeout` | int | `15` | 登录超时（秒） |
| `client.auth.address` | string | `"127.0.0.1:8099"` | 授权客户端连接地址 |
| `client.auth.appName` | string | `"hotgo"` | 应用名 |
| `client.auth.group` | string | `"default"` | 授权分组 |

## 4. logger — 日志配置

使用 YAML 锚点 `&defaultLogger` 统一配置，各模块可引用：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `path` | string | 按模块设定 | 日志文件路径 |
| `file` | string | `"{Y-m-d}.log"` | 日志文件名格式 |
| `level` | string | `"all"` | 日志级别 |
| `stdout` | bool | `true` | 是否输出到控制台 |
| `ctxKeys` | string[] | `["RequestId"]` | 上下文打印 Key |

日志模块分类：
- `logger` — 全局日志（`./storage/log/logger`）
- `cron` — 定时任务日志（`./storage/log/cron`）
- `queue` — 消息队列日志（`./storage/log/queue`）
- `tcpServer` — TCP 服务端日志（`./storage/log/tcpServer`）
- `tcpClient` — TCP 客户端日志（`./storage/log/tcpClient`）

## 5. router — 路由配置

### 5.1 Admin 后台

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `prefix` | string | `"/admin"` | 路由前缀 |
| `exceptLogin` | string[] | `["/sms/send","/wechat/authorizeCall"]` | 免登录路由 |
| `exceptAuth` | string[] | `["/member/info","/role/dynamic",...]` | 免权限验证路由 |

### 5.2 API 接口

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `prefix` | string | `"/api"` | 路由前缀 |
| `exceptLogin` | string[] | `[]` | 免登录路由 |

### 5.3 WebSocket

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `prefix` | string | `"/socket"` | 路由前缀 |
| `exceptLogin` | string[] | `[]` | 免登录路由 |

### 5.4 Home 前台

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `prefix` | string | `"/home"` | 路由前缀 |

## 6. cache — 缓存配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `adapter` | string | `"redis"` | 缓存驱动：`memory`/`redis`/`file` |

- **memory**：进程内存缓存，重启丢失
- **redis**：推荐生产使用，依赖 `redis` 配置
- **file**：文件缓存，存储在 `./storage/cache/`

## 7. token — JWT 令牌配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `secretKey` | string | `"hotgo123"` | JWT 加密密钥（**生产必须修改**） |
| `expires` | int64 | `604800` | 有效期（秒），默认 7 天 |
| `autoRefresh` | bool | `true` | 自动刷新开关 |
| `refreshInterval` | int64 | `86400` | 刷新间隔（秒），默认 1 天 |
| `maxRefreshTimes` | int64 | `30` | 最大刷新次数，-1 不限制 |
| `multiLogin` | bool | `true` | 允许多端登录 |

## 8. queue — 消息队列配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `switch` | bool | `true` | 队列开关 |
| `driver` | string | `"disk"` | 驱动：`disk`/`redis`/`rocketmq`/`kafka` |
| `retry` | int | `2` | 重试次数 |
| `multiComsumerNum` | int | `1` | 消费者并发数 |

### 8.1 Disk 磁盘队列

| 参数 | 说明 |
|------|------|
| `path` | 存储路径（`./storage/diskqueue`） |
| `batchSize` | 批量处理大小（100） |
| `maxBytesPerFile` | 单文件最大字节（100MB） |
| `syncEvery` | 每 N 次写入同步（1000） |
| `syncTimeout` | 同步超时（2s） |

### 8.2 Redis 队列

使用全局 Redis 配置，无额外参数。

### 8.3 RocketMQ

| 参数 | 说明 |
|------|------|
| `address` | 地址 |
| `logLevel` | 日志级别 |
| `groupName` | 消费组名 |
| `namespace` | 命名空间 |
| `accessKey/secretKey` | ACL 凭证 |

### 8.4 Kafka

| 参数 | 说明 |
|------|------|
| `address` | Broker 地址 |
| `version` | 版本号 |
| `randClient` | 随机客户端 ID |
| `multiConsumer` | 多消费者模式 |
| `groupName` | 消费组名 |

## 9. redis — Redis 配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `default.address` | string | `"127.0.0.1:6379"` | 地址 |
| `default.db` | int | `0` | 数据库编号 |
| `default.pass` | string | `""` | 密码 |
| `default.idleTimeout` | string | `"20s"` | 空闲超时 |

## 10. database — 数据库配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `default.link` | string | `"mysql:hotgo:hg123456.@tcp(127.0.0.1:3306)/hotgo?..."` | 连接字符串 |
| `default.debug` | bool | `true` | SQL 调试输出 |
| `default.Prefix` | string | `"hg_"` | 表前缀 |

**支持的数据库：**
- MySQL：`mysql:user:pass@tcp(host:port)/dbname?...`
- PostgreSQL：`pgsql:user=xxx password=xxx host=xxx port=xxx dbname=xxx sslmode=disable`

## 11. jaeger — 链路追踪配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `switch` | bool | `false` | 开关 |
| `endpoint` | string | `"http://tracing-analysis-dc-bj.aliyuncs.com/..."` | Jaeger 端点 |

## 12. hggen — 代码生成配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `allowedIPs` | string[] | `["127.0.0.1","..."]` | 允许使用代码生成的 IP 白名单 |
| `selectDbs` | string[] | `["default"]` | 可选数据库配置组 |
| `disableTables` | string[] | `["hg_sys_gen_codes"]` | 禁止生成的表 |
| `delimiters` | string[] | `["@{","}"]` | 模板定界符 |

### 12.1 模板配置（default）

| 参数 | 说明 |
|------|------|
| `curd.templates[].group` | 应用分组（admin/home） |
| `curd.templates[].isAddon` | 是否插件模板 |
| `curd.templates[].masterPackage` | 主包路径 |
| `curd.templates[].templatePath` | 模板文件路径 |
| `curd.templates[].apiPath` | 生成的 API 路径 |
| `curd.templates[].inputPath` | 生成的 Input 路径 |
| `curd.templates[].controllerPath` | 生成的 Controller 路径 |
| `curd.templates[].logicPath` | 生成的 Logic 路径 |
| `curd.templates[].routerPath` | 生成的 Router 路径 |
| `curd.templates[].sqlPath` | 生成的 SQL 路径 |
| `curd.templates[].webApiPath` | 生成的前端 API 路径 |
| `curd.templates[].webViewsPath` | 生成的前端 Views 路径 |

---

## 二、数据库动态配置

> 以下配置存储在数据库 `hg_sys_config` 表中，可通过后台管理界面在线修改，修改后实时生效（支持集群同步）。
> 
> 配置加载流程：`GetConfigByGroup(group)` → 数据库查询 → `ConversionType()` 类型转换 → `gconv.Scan()` 映射到 Config struct。
> 
> 完整的 ISysConfig 接口定义参见 [关键接口定义 — ISysConfig](dev-api-interfaces.md#91-isysconfig--系统配置)

### 13. basic — 基础配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `basicCaptchaSwitch` | int | 登录验证码开关（1=开启） |
| `basicCloseText` | string | 系统关闭时的提示文案 |
| `basicCopyright` | string | 版权信息 |
| `basicIcpCode` | string | ICP 备案号 |
| `basicLogo` | string | 系统 Logo 地址 |
| `basicName` | string | 系统名称 |
| `basicDomain` | string | 系统域名 |
| `basicWsAddr` | string | WebSocket 连接地址 |
| `basicRegisterSwitch` | int | 注册功能开关（1=开启） |
| `basicSystemOpen` | bool | 系统是否开放 |

### 14. login — 登录注册配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `loginRegisterSwitch` | int | 注册开关 |
| `loginCaptchaSwitch` | int | 验证码开关 |
| `loginCaptchaType` | int | 验证码类型 |
| `loginAvatar` | string | 默认头像地址 |
| `loginRoleId` | int64 | 注册用户默认角色 ID |
| `loginDeptId` | int64 | 注册用户默认部门 ID |
| `loginPostIds` | []int64 | 注册用户默认岗位 ID 列表 |
| `loginProtocol` | string | 用户协议内容 |
| `loginPolicy` | string | 隐私政策内容 |
| `loginAutoOpenId` | int | 自动获取微信 OpenID |
| `loginForceInvite` | int | 强制邀请注册 |

### 15. upload — 文件上传配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `uploadDrive` | string | 存储驱动：`local`/`ucloud`/`cos`/`oss`/`qiniu`/`minio` |
| `uploadFileSize` | int64 | 文件上传最大大小（字节） |
| `uploadFileType` | string | 允许的文件类型（逗号分隔） |
| `uploadImageSize` | int64 | 图片上传最大大小（字节） |
| `uploadImageType` | string | 允许的图片类型（逗号分隔） |

**本地存储 (local)：**

| 参数 | 说明 |
|------|------|
| `uploadLocalPath` | 本地存储路径 |

**UCloud 对象存储：**

| 参数 | 说明 |
|------|------|
| `uploadUCloudBucketHost` | Bucket 域名 |
| `uploadUCloudBucketName` | Bucket 名称 |
| `uploadUCloudEndpoint` | 接入端点 |
| `uploadUCloudFileHost` | 文件访问域名 |
| `uploadUCloudPath` | 存储路径前缀 |
| `uploadUCloudPublicKey` / `uploadUCloudPrivateKey` | API 密钥 |

**腾讯云 COS：**

| 参数 | 说明 |
|------|------|
| `uploadCosSecretId` / `uploadCosSecretKey` | API 密钥 |
| `uploadCosBucketURL` | Bucket 访问地址 |
| `uploadCosPath` | 存储路径前缀 |

**阿里云 OSS：**

| 参数 | 说明 |
|------|------|
| `uploadOssSecretId` / `uploadOssSecretKey` | API 密钥 |
| `uploadOssEndpoint` | 接入端点 |
| `uploadOssBucketURL` | Bucket 访问地址 |
| `uploadOssBucket` | Bucket 名称 |
| `uploadOssPath` | 存储路径前缀 |

**七牛云：**

| 参数 | 说明 |
|------|------|
| `uploadQiNiuAccessKey` / `uploadQiNiuSecretKey` | API 密钥 |
| `uploadQiNiuDomain` | CDN 域名 |
| `uploadQiNiuBucket` | Bucket 名称 |
| `uploadQiNiuPath` | 存储路径前缀 |

**MinIO：**

| 参数 | 说明 |
|------|------|
| `uploadMinioAccessKey` / `uploadMinioSecretKey` | API 密钥 |
| `uploadMinioEndpoint` | MinIO 接入端点 |
| `uploadMinioUseSSL` | 是否使用 SSL（1=是） |
| `uploadMinioBucket` | Bucket 名称 |
| `uploadMinioDomain` | 访问域名 |
| `uploadMinioPath` | 存储路径前缀 |

> 存储驱动接口定义参见 [关键接口定义 — UploadDrive](dev-api-interfaces.md#5-文件存储接口)

### 16. sms — 短信配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `smsDrive` | string | 短信驱动：`aliyun`/`tencent` |
| `smsMinInterval` | int | 同一号码最小发送间隔（秒） |
| `smsMaxIpLimit` | int | 同一 IP 每日最大发送次数 |
| `smsCodeExpire` | int | 验证码过期时间（秒） |

**阿里云短信：**

| 参数 | 说明 |
|------|------|
| `smsAliYunAccessKeyID` / `smsAliYunAccessKeySecret` | API 密钥 |
| `smsAliYunSign` | 签名名称 |
| `smsAliYunTemplate` | 模板列表 `[{key, value}]` |

**腾讯云短信：**

| 参数 | 说明 |
|------|------|
| `smsTencentSecretId` / `smsTencentSecretKey` | API 密钥 |
| `smsTencentEndpoint` | 接入端点 |
| `smsTencentRegion` | 地域 |
| `smsTencentAppId` | 应用 ID |
| `smsTencentSign` | 签名名称 |
| `smsTencentTemplate` | 模板列表 `[{key, value}]` |

### 17. smtp — 邮件配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `smtpUser` | string | SMTP 用户名 |
| `smtpPass` | string | SMTP 密码 |
| `smtpAddr` | string | SMTP 地址 |
| `smtpHost` | string | SMTP 主机 |
| `smtpPort` | int64 | SMTP 端口 |
| `smtpSendName` | string | 发件人名称 |
| `smtpAdminMailbox` | string | 管理员邮箱 |
| `smtpMinInterval` | int | 同一邮箱最小发送间隔（秒） |
| `smtpMaxIpLimit` | int | 同一 IP 每日最大发送次数 |
| `smtpCodeExpire` | int | 验证码过期时间（秒） |
| `smtpTemplate` | array | 邮件模板列表 `[{key, value}]` |

### 18. pay — 支付配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `payDebug` | bool | 支付调试模式 |

**支付宝：**

| 参数 | 说明 |
|------|------|
| `payAliPayAppId` | 应用 ID |
| `payAliPayPrivateKey` | 应用私钥 |
| `payAliPayAppCertPublicKey` | 应用公钥证书 |
| `payAliPayRootCert` | 支付宝根证书 |
| `payAliPayCertPublicKeyRSA2` | 支付宝公钥证书 |

**微信支付：**

| 参数 | 说明 |
|------|------|
| `payWxPayAppId` | 应用 ID |
| `payWxPayMchId` | 商户号 |
| `payWxPaySerialNo` | 证书序列号 |
| `payWxPayAPIv3Key` | API v3 密钥 |
| `payWxPayPrivateKey` | 商户私钥 |

**QQ 支付：**

| 参数 | 说明 |
|------|------|
| `payQQPayAppId` | 应用 ID |
| `payQQPayMchId` | 商户号 |
| `payQQPayApiKey` | API 密钥 |

> 支付接口定义参见 [关键接口定义 — PayClient](dev-api-interfaces.md#4-支付网关接口)

### 19. wechat — 微信配置

**公众号：**

| 参数 | 说明 |
|------|------|
| `officialAccountAppId` | 公众号 AppID |
| `officialAccountAppSecret` | 公众号 AppSecret |
| `officialAccountToken` | 消息校验 Token |
| `officialAccountEncodingAESKey` | 消息加密密钥 |

**开放平台：**

| 参数 | 说明 |
|------|------|
| `openPlatformAppId` | 开放平台 AppID |
| `openPlatformAppSecret` | 开放平台 AppSecret |
| `openPlatformToken` | 消息校验 Token |
| `openPlatformEncodingAESKey` | 消息加密密钥 |

### 20. geo — 地理位置配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `geoAmapWebKey` | string | 高德地图 Web 服务 Key |

### 21. cash — 提现配置

| 参数 | 类型 | 说明 |
|------|------|------|
| `cashSwitch` | bool | 提现功能开关 |
| `cashMinFee` | float64 | 最低手续费 |
| `cashMinFeeRatio` | float64 | 手续费比例 |
| `cashMinMoney` | float64 | 最低提现金额 |
| `cashTips` | string | 提现提示文案 |

## 三、配置体系说明

### 配置分层

| 层级 | 存储 | 修改方式 | 重启生效 |
|------|------|----------|----------|
| **本地文件配置** | `config.yaml` | 编辑文件 | 需重启 |
| **数据库动态配置** | `hg_sys_config` 表 | 后台管理界面 | 实时生效 |
| **插件配置** | `hg_sys_addons_config` 表 | 后台插件配置页 | 实时生效 |

### 配置类型支持

数据库配置的 `type` 字段支持以下数据类型：

`string`、`int`、`int8`、`int16`、`int32`、`int64`、`uint`、`uint8`、`uint16`、`uint32`、`uint64`、`float32`、`float64`、`bool`、`date`、`datetime`、`[]string`、`[]int`、`[]int64`

### 集群配置同步

开启集群模式（`system.isCluster: true`）后，以下数据库配置在任一实例修改后会通过 Redis PubSub 自动同步到所有实例：

- `wechat`（微信配置）
- `pay`（支付配置）
- `upload`（上传配置）
- `sms`（短信配置）

> 集群部署详情参见 [部署指南 — 集群部署](dev-deployment.md#7-集群部署)

---

