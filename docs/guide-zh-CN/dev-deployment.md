# 部署指南

> 本文档涵盖 HotGo 的环境要求、编译构建、Docker 部署、Kubernetes 部署和 Nginx 配置。
>
> **文档导航**：[← 返回主文档](dev-guide.md) | [技术架构](dev-architecture.md) | [服务端模块](dev-server-modules.md) | [前端模块](dev-frontend-modules.md) | [接口定义](dev-api-interfaces.md) | [插件系统](dev-plugin-system.md) | [配置参考](dev-config-reference.md)

## 1. 环境要求

| 依赖 | 最低版本 | 推荐版本 | 说明 |
|------|----------|----------|------|
| Go | 1.23+ | 1.24 | 后端编译运行 |
| Node.js | 16+ | 20+ | 前端编译 |
| pnpm | 8.0+ | 最新 | 前端包管理 |
| MySQL | 5.7+ | 8.0+ | 主数据库 |
| PostgreSQL | 14+ | 16+ | 可选数据库 |
| Redis | 4.0+ | 7.0+ | 缓存和队列 |
| GoFrame CLI | 最新 | — | 代码生成工具 |

## 2. 安装步骤

### 2.1 获取代码

```bash
git clone https://github.com/bufanyun/hotgo.git
cd hotgo
```

### 2.2 数据库初始化

```bash
# MySQL
mysql -u root -p < server/storage/data/hotgo.sql

# 或 PostgreSQL
psql -U postgres -d hotgo < server/storage/data/hotgo-pg.sql
```

### 2.3 后端配置

```bash
cd server
cp manifest/config/config.example.yaml manifest/config/config.yaml
```

**必须修改的配置项：**

| 配置项 | 说明 |
|--------|------|
| `database.default.link` | 数据库连接字符串 |
| `redis.default.address` | Redis 地址 |
| `redis.default.pass` | Redis 密码 |
| `token.secretKey` | JWT 密钥（**生产环境必改**） |
| `system.mode` | 改为 `"product"` |
| `system.debug` | 改为 `false` |

### 2.4 启动开发环境

```bash
# 启动后端（所有服务）
cd server
go run main.go

# 启动前端
cd web
pnpm install
pnpm dev
```

## 3. 编译构建

### 3.1 一键编译（推荐）

```bash
cd server
make build
```

此命令会：
1. 编译前端 `cd ../web && pnpm install && pnpm build`
2. 复制前端产物到 `server/resource/public/admin/`
3. 使用 GoFrame CLI 编译后端 `gf build`
4. 输出到 `server/temp/` 目录

### 3.2 手动分步编译

**前端编译：**

```bash
cd web
pnpm install
pnpm build
# 产物在 web/dist/ 目录
```

**后端编译：**

```bash
cd server

# 安装 GoFrame CLI（如未安装）
make cli
# 或手动：go install github.com/gogf/gf/cmd/gf/v2@latest

# 编译
gf build main.go
# 输出到 server/temp/ 目录
```

### 3.3 生产配置修改要点

```yaml
# config.yaml 生产环境调整

system:
  debug: false
  mode: "product"
  isDemo: false

server:
  pprofEnabled: false      # 关闭 PProf
  openapiPath: ""          # 关闭 OpenAPI
  swaggerPath: ""          # 关闭 Swagger

token:
  secretKey: "你的复杂密钥"  # 必须修改

database:
  default:
    debug: false            # 关闭 SQL 调试

jaeger:
  switch: false             # 按需开启链路追踪
```

## 4. Docker 部署

### 4.1 Dockerfile 说明

> 文件：`server/manifest/docker/Dockerfile`

```dockerfile
FROM loads/alpine:3.8

ENV WORKDIR /app
ADD hack                $WORKDIR/hack/
ADD manifest/config     $WORKDIR/manifest/config/
ADD resource            $WORKDIR/resource/
ADD ./temp/linux_amd64/hotgo $WORKDIR/hotgo
ADD ./manifest/docker/entrypoint.sh $WORKDIR/entrypoint.sh

RUN chmod +x $WORKDIR/hotgo
RUN chmod +x $WORKDIR/entrypoint.sh

WORKDIR $WORKDIR
CMD ./entrypoint.sh
```

### 4.2 构建镜像

```bash
cd server

# 方式一：使用 Makefile
make image tag=v2.18.6

# 方式二：使用 GoFrame CLI
gf docker main.go -p -tn hotgo:v2.18.6
```

### 4.3 运行容器

```bash
docker run -d \
  --name hotgo \
  -p 8000:8000 \
  -v /path/to/config.yaml:/app/manifest/config/config.yaml \
  -v /path/to/storage:/app/storage \
  hotgo:v2.18.6
```

## 5. Kubernetes 部署

### 5.1 Kustomize 目录结构

```
server/manifest/deploy/kustomize/
├── base/
│   ├── deployment.yaml        # 基础 Deployment
│   ├── service.yaml           # Service（80→8000）
│   └── kustomization.yaml
└── overlays/
    └── develop/
        ├── configmap.yaml     # ConfigMap（config.yaml）
        ├── deployment.yaml    # 覆盖镜像
        └── kustomization.yaml
```

### 5.2 部署步骤

```bash
cd server

# 方式一：使用 Makefile
make image tag=v2.18.6
make deploy

# 方式二：手动
kubectl apply -k manifest/deploy/kustomize/overlays/develop/
```

### 5.3 配置自定义

编辑 `overlays/develop/configmap.yaml`，将 `config.yaml` 内容嵌入 ConfigMap。

编辑 `overlays/develop/deployment.yaml` 修改镜像名和标签。

## 6. Nginx 反向代理

### 6.1 基础配置（HTTP + WebSocket）

```nginx
upstream hotgo_server {
    server 127.0.0.1:8000;
}

server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location /admin {
        alias /path/to/web/dist;
        try_files $uri $uri/ /admin/index.html;
    }

    # 后端 API 代理
    location /admin/ {
        proxy_pass http://hotgo_server/admin/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # API 接口
    location /api/ {
        proxy_pass http://hotgo_server/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # WebSocket 代理
    location /socket {
        proxy_pass http://hotgo_server/socket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }

    # 前台页面
    location /home {
        proxy_pass http://hotgo_server/home;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 6.2 HTTPS 配置

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # ... 同上代理配置 ...
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
}
```

## 7. 集群部署

### 7.1 前提条件

- 必须使用 **Redis** 作为缓存驱动（`cache.adapter: "redis"`）
- 所有实例共享同一 Redis 和数据库

### 7.2 配置

```yaml
system:
  isCluster: true    # 开启集群模式
```

### 7.3 同步机制

开启集群模式后，以下数据变更会通过 Redis PubSub 自动同步到所有实例：

| 同步内容 | 触发场景 |
|----------|----------|
| 系统配置 | 后台修改配置后 |
| IP 黑名单 | 新增/修改/删除黑名单 |
| 超管数据 | 超管信息变更 |

### 7.4 注意事项

- **消息队列**：如果使用 `disk` 驱动，各实例队列独立，建议切换为 `redis` / `rocketmq` / `kafka`
- **文件存储**：如果使用 `local` 驱动，各实例文件独立，建议切换为云存储驱动
- **定时任务**：建议只在一个实例上启动 `cron` 服务，避免重复执行
- **WebSocket**：每个实例维护自己的连接，跨实例消息推送需通过 Redis PubSub

## 8. Makefile 命令速查

| 命令 | 说明 |
|------|------|
| `make all` | 热编译启动所有服务 |
| `make http` | 仅启动 HTTP 服务 |
| `make queue` | 仅启动消息队列 |
| `make cron` | 仅启动定时任务 |
| `make web` | 启动前端开发服务 |
| `make build` | 一键编译（前端+后端） |
| `make dao` | 生成 DAO/DO/Entity |
| `make service` | 生成 Service 接口 |
| `make image tag=xxx` | 构建 Docker 镜像 |
| `make deploy` | K8s 部署（Kustomize） |
| `make start` | 构建+部署+端口转发 |
| `make refresh` | 刷新 Casbin 权限 |
| `make clear` | 清理 Casbin 权限 |
| `make lint` | 代码质量检查 |
| `make cli` | 安装/更新 GoFrame CLI |

## 9. 常用启动模式

```bash
# 开发模式（所有服务）
go run main.go

# 仅 HTTP 服务
go run main.go http

# 仅队列消费者
go run main.go queue

# 仅定时任务
go run main.go cron

# 刷新 Casbin 权限
go run main.go tools -m=casbin -a1=refresh

# 清空 Casbin 权限
go run main.go tools -m=casbin -a1=clear

# 版本升级
go run main.go up
```

---

*本文档基于 HotGo v2.18.6 代码库分析生成。*
