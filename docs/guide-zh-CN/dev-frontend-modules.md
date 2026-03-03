# 前端核心模块详解

> 本文档详细说明 HotGo 前端各模块的职责、关键文件和交互关系。
>
> **文档导航**：[← 返回主文档](dev-guide.md) | [技术架构](dev-architecture.md) | [服务端模块](dev-server-modules.md) | [接口定义 →](dev-api-interfaces.md) | [插件系统](dev-plugin-system.md) | [配置参考](dev-config-reference.md) | [部署指南](dev-deployment.md)

## 1. 整体模块结构

```
web/src/
├── main.ts              # 应用入口
├── App.vue              # 根组件
├── api/                 # API 接口调用（65 个文件）
├── components/          # 公共组件（106 个文件，18 个目录）
├── directives/          # 自定义指令（6 个）
├── enums/               # 枚举定义（11 个文件）
├── hooks/               # 组合式函数（32 个文件）
├── layout/              # 布局组件
├── locale/              # 国际化
├── plugins/             # 插件注册
├── router/              # 路由系统（7 个文件）
├── settings/            # 应用设置（4 个文件）
├── store/               # 状态管理 Pinia（9 个 Store）
├── utils/               # 工具函数（35 个文件）
└── views/               # 页面视图（294 个文件，20+ 个模块）
```

## 2. 启动流程

`main.ts` 中 `bootstrap()` 按以下顺序执行：

| 步骤 | 操作 | 说明 |
|------|------|------|
| 1 | `app.use(i18n)` | 注册国际化 |
| 2 | `setupNaive(app)` | 注册约 50 个 Naive UI 全局组件 |
| 3 | `setupDirectives(app)` | 注册自定义指令 |
| 4 | `setupStore(app)` | 挂载 Pinia |
| 5 | `appProvider.mount('#appProvider')` | 先挂载 Provider（Dialog/Message 上下文） |
| 6 | `setupRouter(app)` | 挂载路由 + 创建守卫 |
| 7 | `router.isReady()` | 等待路由就绪 |
| 8 | `setupWebsocket()` | 启动 WebSocket |
| 9 | `app.mount('#app')` | 挂载主应用 |

`App.vue` 提供 NConfigProvider（主题/语言）、无操作锁屏（默认1小时）、暗色主题切换。

## 3. 路由系统

> 后端路由分组架构参见 [技术架构 — 路由体系](dev-architecture.md#6-路由体系)

### 3.1 路由守卫流程

```
beforeEach:
  ① 白名单路径（/login）→ 加载登录配置 → 放行
  ② 无 token → 检查 meta.ignoreAuth → 否则重定向登录页
  ③ 已有动态路由 → 放行
  ④ 首次加载: GetInfo → LoadLoginConfig → 微信OpenID → GetConfig
     → generateRoutes → addRoute → 404 兜底

afterEach: 更新 title + 停止 Loading Bar
```

### 3.2 动态路由生成

1. 调用 `/role/dynamic` API 获取菜单数据
2. `routerGenerator()` 递归转换为路由配置
3. `asyncImportRoute()` 通过 `import.meta.glob` 动态匹配 views 组件
4. 过滤 `hidden=true` 和按钮权限（type=3）

## 4. 状态管理（Pinia）

| Store ID | 文件 | 核心状态 | 持久化 |
|----------|------|----------|--------|
| `app-member` | `user.ts` | token, info, config, loginConfig | localStorage |
| `app-async-route` | `asyncRoute.ts` | menus, routers, keepAliveComponents | — |
| `app-dict` | `dict.ts` | dictMap（字典缓存） | localStorage |
| `app-project-setting` | `projectSetting.ts` | navMode, navTheme, headerSetting | — |
| `app-design-setting` | `designSetting.ts` | darkTheme, appTheme | — |
| `app-lockscreen` | `lockscreen.ts` | isLock, lockTime | localStorage |
| `app-tabs-view` | `tabsView.ts` | tabsList | — |
| `notificationStore` | `notification.ts` | 通知/公告/私信 | — |
| `I18nStore` | `i18n.ts` | locale | localStorage |

**user.ts 核心 Actions**: login, mobileLogin, GetInfo, GetConfig, LoadLoginConfig, logout

**dict.ts 核心方法**: loadOptions, getOption, getLabel, getType, getExtra, hasValue

**notification.ts**: triggerNewMessages（WebSocket推送）, pullMessages（拉取未读）

## 5. API 层

按模块组织在 `src/api/` 下，共 65 个文件：

| 模块 | 说明 |
|------|------|
| `system/` | 用户登录/注册、菜单、角色 |
| `sys/` | 系统配置、黑名单、定时任务 |
| `org/` | 部门、岗位、组织用户 |
| `apply/` | 附件、通知、省市区 |
| `order/` | 订单管理 |
| `dict/` | 字典管理 |
| `dashboard/` | 仪表盘 |
| `develop/` | 开发工具 |
| `monitor/` | 监控 |
| `log/` | 日志 |
| `cash/` | 资金 |
| `pay/` | 支付 |

统一调用模式：`http.request({ url, method, params })`

## 6. HTTP 层（Axios 封装）

请求流程：`beforeRequestHook` → `requestInterceptors`（注入Token+Locale） → Axios → `transformRequestData`（解析 code/data/message）

| 后端 code | 行为 |
|-----------|------|
| 0 | 返回 data |
| -1 | $message.error 提示 |
| 61 | Dialog 提示重新登录 |

配置：timeout=30s, urlPrefix=`/admin`, withToken=true

## 7. 页面视图

| 目录 | 说明 |
|------|------|
| `login/` | 登录、注册 |
| `dashboard/` | 仪表盘 |
| `system/` | 系统管理 |
| `permission/` | 权限管理 |
| `org/` | 组织管理 |
| `apply/` | 应用管理 |
| `asset/` | 资产管理 |
| `log/` | 日志管理 |
| `monitor/` | 监控 |
| `develop/` | 开发工具 |
| `home/` | 个人中心 |
| `curdDemo/` | CRUD 示例 |
| `exception/` | 错误页面（403/404/500） |

**典型页面模式（三文件分离）**：
- `model.ts` — 数据模型层（State类、rules、schemas、columns、loadOptions）
- `index.vue` — 列表页（BasicForm + BasicTable + TableAction + v-permission）
- `edit.vue` — 编辑弹窗（BasicModal + useModal）

## 8. 组件系统

| 组件 | 核心功能 |
|------|----------|
| **BasicTable** | 基于 NDataTable，集成分页/加载/列设置/密度调节/自动高度 |
| **BasicForm** | 基于 NForm+NGrid，Schema 驱动表单渲染，权限过滤 |
| **BasicModal** | 基于 NModal，可拖拽弹窗 |
| **Upload** | 基础/分片/图片/文件上传 |
| **FileChooser** | 文件选择器 |
| **Editor** | 富文本编辑器（vue-quill） |
| **YamlEditor** | YAML 编辑器（codemirror） |
| **CitySelector** | 城市选择器 |
| **IconSelector** | 图标选择器 |
| **Application** | 全局 Provider |

## 9. Hooks

| Hook | 说明 |
|------|------|
| `usePermission()` | hasPermission, hasEveryPermission, hasSomePermission |
| `useGlobSetting()` | title, apiUrl, urlPrefix, uploadUrl |
| `useProjectSetting()` | 导航模式、主题、多标签等 |
| `usePage()` | useGo, useRedo |
| `useSendCode()` | 发送验证码（60s 倒计时） |
| `useECharts()` | ECharts 封装 |
| `useLoading()` | 加载状态 |

## 10. 自定义指令

| 指令 | 功能 | 用法示例 |
|------|------|----------|
| `v-permission` | 权限控制 | `v-permission="{ action: ['/xxx/edit'], effect: 'disabled' }"` |
| `v-copy` | 点击复制 | `v-copy="text"` |
| `v-debounce` | 防抖（500ms） | `v-debounce="handler"` |
| `v-throttle` | 节流（1000ms） | `v-throttle="handler"` |
| `v-draggable` | 元素拖拽 | `v-draggable` |
| `v-click-outside` | 点击外部回调 | `v-click-outside="handler"` |

## 11. 布局系统

支持 vertical（左侧菜单）、horizontal（顶部菜单）、horizontal-mix（混合）三种模式。

核心功能：响应式（<=800px 移动端）、菜单折叠（<=950px）、固定头部、多标签页、KeepAlive 缓存、路由过渡动画（zoom-fade 等 6 种）。

## 12. WebSocket

> 后端 WebSocket 架构参见 [技术架构 — WebSocket 数据流](dev-architecture.md#83-websocket-数据流)

- 连接：`${wsAddr}?authorization=${token}`
- 心跳：5 秒 ping，5 秒超时重连
- 自动重连：10 秒间隔
- 事件：ping（心跳）、kick（强踢）、notice（消息通知）

## 13. 国际化

vue-i18n，支持 zh-CN/zh-TW/en，默认 zh-CN。与 i18n Store 集成，HTTP 自动携带 Locale header。

## 14. 设置系统

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `navMode` | `vertical` | 导航模式 |
| `navTheme` | `dark` | 导航风格 |
| `permissionMode` | `BACK` | 后端动态权限 |
| `headerSetting.fixed` | `true` | 固定顶部 |
| `menuSetting.menuWidth` | `200` | 菜单宽度 |
| `darkTheme` | `false` | 暗色主题 |
| `appTheme` | `#2d8cf0` | 主题色 |
| `isPageAnimate` | `true` | 路由动画 |

## 15. 模块交互关系

```
main.ts
  ├─ i18n (locale/)
  ├─ plugins/ (naive + directives)
  ├─ store/
  │   ├─ user → api/system/user → http/axios
  │   ├─ asyncRoute → router/generator-routers → api/system/menu
  │   ├─ dict → api/dict
  │   ├─ notification → api/apply/notice
  │   └─ projectSetting ← settings/
  ├─ router/ → user store + asyncRoute store
  ├─ websocket → user store + notification store
  └─ App.vue → layout/ → views/
       └─ components/ + hooks/
```

---

