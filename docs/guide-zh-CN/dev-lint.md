# 代码质量与格式化规范

> 本文档描述 HotGo 项目的代码质量检查（Lint）和格式化配置，覆盖后端 Go 和前端 Vue/TypeScript 两端。

## 1. 后端 Go — golangci-lint

### 1.1 概述

后端使用 [golangci-lint v2](https://golangci-lint.run/) 进行代码质量检查，配置文件为 `server/.golangci.yml`。

### 1.2 运行方式

```bash
cd server

# 运行 lint 并自动修复
make lint

# 等价于
golangci-lint run --fix
```

### 1.3 启用的 Linter

| Linter | 说明 |
|--------|------|
| `errcheck` | 检查未处理的错误返回值 |
| `errchkjson` | 检查 JSON 编解码错误处理 |
| `funlen` | 函数长度限制（行数 ≤ 340） |
| `goconst` | 检测可提取为常量的重复字面量 |
| `gocritic` | Go 代码风格和性能建议 |
| `govet` | Go 官方静态分析工具 |
| `misspell` | 拼写检查 |
| `nolintlint` | 检查 nolint 指令的规范性 |
| `revive` | 通用 Go 代码规范检查 |
| `staticcheck` | 高级静态分析 |
| `usestdlibvars` | 建议使用标准库变量 |
| `whitespace` | 空白行检查 |

### 1.4 启用的 Formatter

| Formatter | 说明 |
|-----------|------|
| `gci` | 自动整理 import 分组（标准库 → 第三方 → 本地包） |
| `gofmt` | 代码格式化，自动替换废弃写法（如 `interface{}` → `any`） |

### 1.5 关键规则配置

- **函数长度**：单个函数不超过 340 行
- **行长度**：单行不超过 380 字符（revive `line-length-limit`）
- **函数返回值**：最多 4 个返回值
- **Import 分组顺序**：标准库 → 第三方 → `hotgo/` 本地包 → 点导入
- **自动替换**：`interface{}` → `any`，`ioutil.*` → 对应的 `io.*` / `os.*`

### 1.6 代码生成集成

代码生成器（开发工具 → 代码生成）支持在提交生成后自动运行 lint 检查：

1. 进入代码生成配置页面
2. 在「高级设置」区域勾选 **「生成后运行 [golangci-lint]」**
3. 提交生成时会在所有文件写入和 `gf gen service` 执行完毕后，自动运行 `golangci-lint run --fix`

该选项默认勾选，可按需取消。

---

## 2. 前端 Vue/TypeScript

前端采用 **ESLint + Prettier + Stylelint** 三合一方案，覆盖 JS/TS/Vue 代码规范、代码格式化和 CSS 样式规范。

### 2.1 配置文件一览

| 文件 | 用途 |
|------|------|
| `web/eslint.config.js` | ESLint 代码规范（Flat Config 格式） |
| `web/prettier.config.cjs` | Prettier 代码格式化 |
| `web/stylelint.config.cjs` | Stylelint CSS 样式规范 |
| `web/.editorconfig` | 编辑器基础设置（缩进、换行、编码） |
| `web/.prettierignore` | Prettier 忽略文件 |
| `web/.stylelintignore` | Stylelint 忽略文件 |

### 2.2 运行方式

```bash
cd web

# 一次性运行所有检查（推荐）
pnpm run lint

# 分别运行
pnpm run lint:eslint      # ESLint 检查 + 自动修复
pnpm run lint:prettier     # Prettier 格式化
pnpm run lint:stylelint    # Stylelint 检查 + 自动修复
```

### 2.3 ESLint 配置

**配置文件**：`web/eslint.config.js`（ESLint 9 Flat Config 格式）

**继承规则集**：

| 规则集 | 说明 |
|--------|------|
| `plugin:vue/vue3-recommended` | Vue 3 官方推荐规则 |
| `plugin:@typescript-eslint/recommended` | TypeScript 推荐规则 |
| `prettier` | 关闭与 Prettier 冲突的规则 |
| `plugin:prettier/recommended` | Prettier 集成到 ESLint |

**检查范围**：`src/` 和 `mock/` 目录下的 `.vue`、`.ts`、`.tsx` 文件

**关键规则**：

| 规则 | 配置 | 说明 |
|------|------|------|
| `vue/html-self-closing` | error | HTML 标签自闭合规范 |
| `@typescript-eslint/no-unused-vars` | error（忽略变量名和 catch） | 未使用变量检查 |
| `@typescript-eslint/no-this-alias` | error（允许 `that`） | 限制 this 别名 |
| `vue/script-setup-uses-vars` | error | script setup 变量使用检查 |

### 2.4 Prettier 配置

**配置文件**：`web/prettier.config.cjs`

| 选项 | 值 | 说明 |
|------|----|------|
| `printWidth` | 100 | 单行最大宽度 |
| `tabWidth` | 2 | 缩进宽度 |
| `useTabs` | false | 使用空格缩进 |
| `semi` | true | 语句末尾加分号 |
| `singleQuote` | true | 使用单引号 |
| `trailingComma` | es5 | ES5 兼容的尾逗号 |
| `vueIndentScriptAndStyle` | true | Vue 文件 script/style 块内缩进 |
| `htmlWhitespaceSensitivity` | strict | HTML 空白敏感模式 |
| `endOfLine` | auto | 自动适配换行符 |

### 2.5 Stylelint 配置

**配置文件**：`web/stylelint.config.cjs`

**继承规则集**：

| 规则集 | 说明 |
|--------|------|
| `stylelint-config-standard` | Stylelint 标准规则 |
| `stylelint-config-recommended-vue` | Vue SFC 样式支持 |

**自定义语法支持**：

| 文件类型 | 解析器 |
|----------|--------|
| `.vue` | `postcss-html` |
| `.less` | `postcss-less` |

**插件**：`stylelint-order`（CSS 属性排序）

**属性排序规则**（warning 级别）：
1. `$` 变量
2. CSS 自定义属性
3. `@` 规则
4. 声明
5. `@supports`
6. `@media`
7. 嵌套规则

### 2.6 EditorConfig

**配置文件**：`web/.editorconfig`

| 设置 | 值 |
|------|----|
| 字符编码 | UTF-8 |
| 换行符 | LF |
| 缩进方式 | 空格，2 个宽度 |
| 最大行宽 | 100 |
| 文件末尾换行 | 是 |
| Makefile | Tab 缩进 |

---

## 3. Git Hooks — 提交前自动检查

项目通过 **Husky + lint-staged** 在 Git 提交前自动运行 lint 检查，确保提交到仓库的代码符合规范。

### 3.1 配置文件

| 文件 | 用途 |
|------|------|
| `.husky/pre-commit` | Git 提交前钩子，运行前端 lint-staged + 后端 golangci-lint |
| `.husky/commit-msg` | 校验 commit message 格式（commitlint） |
| `web/package.json` 的 `lint-staged` 字段 | 定义暂存文件的检查规则 |
| `web/commitlint.config.js` | Commit message 格式规范 |

### 3.2 pre-commit 检查流程

Git 提交时，`pre-commit` 钩子会依次执行：

1. **前端**：进入 `web/` 目录运行 `lint-staged`，对暂存的前端文件执行 ESLint/Prettier/Stylelint 检查
2. **后端**：检测暂存文件中是否包含 `server/` 下的 `.go` 文件，如果有则进入 `server/` 目录运行 `golangci-lint run --fix` 进行代码质量检查和自动修复

任一步骤失败都会阻止提交，确保入库代码符合规范。

### 3.3 lint-staged 规则（前端）

| 文件类型 | 执行的检查 |
|----------|-----------|
| `*.{vue,js,ts,tsx}` | ESLint 修复 + Prettier 格式化 |
| `*.{css,less,scss,postcss,vue}` | Stylelint 修复 + Prettier 格式化 |
| `*.{json,md,html}` | Prettier 格式化 |

### 3.4 Commit Message 规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>
```

**允许的 type**：

| Type | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | 修复 Bug |
| `docs` | 文档变更 |
| `style` | 代码格式（不影响逻辑） |
| `refactor` | 重构 |
| `perf` | 性能优化 |
| `test` | 测试 |
| `build` | 构建/依赖变更 |
| `ci` | CI 配置 |
| `chore` | 杂项 |
| `revert` | 回退 |

**示例**：

```bash
git commit -m "feat(article): 新增文章管理模块"
git commit -m "fix(login): 修复登录超时处理"
git commit -m "docs: 更新 lint 配置文档"
```

### 3.5 初始化

首次 clone 项目后，执行 `pnpm install` 会自动通过 `prepare` 脚本初始化 Husky Git hooks。如果 hooks 未生效，可手动执行：

```bash
# 在项目根目录执行
npx husky
```

---

## 4. IDE 集成建议

### VS Code

项目已配置 `.vscode/settings.json`，建议安装以下扩展：

| 扩展 | 用途 |
|------|------|
| ESLint | JS/TS/Vue 代码检查 |
| Prettier | 代码格式化 |
| Stylelint | CSS 样式检查 |
| EditorConfig | 编辑器配置 |
| Go / gopls | Go 语言支持 |
| golangci-lint | Go lint 集成 |

推荐开启保存时自动格式化和自动修复。

### GoLand / WebStorm

在 Settings → Languages & Frameworks 中启用对应的 ESLint、Prettier、Stylelint 支持，并配置为保存时自动运行。
