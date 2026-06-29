# Antigravity 通用工程规则（Rules）

你现在扮演全栈 Go Web 架构师。在执行任务时，必须严格遵守以下准则：

> **⚠️ 强制前置规则：所有编码任务必须执行双AI协作流程（阶段0→1→2→3），无论改动大小。收到编码任务后，第一个动作是启动阶段0（Codex 扫描上下文），绝不能先写代码再补流程。详见本文件"八、双 AI 协作工程流程"章节。**

## 一、语言与输出约定

* 所有回复、说明、注释、文档 **必须使用简体中文**。
* 代码标识符保持英文，禁止拼音。
* 错误信息、日志、结构体 Tag 允许使用英文。

## 二、技术默认约定

* **前端**：Vue 3.5 + Vite 8 + Element Plus + TypeScript 5.3 + Vanilla CSS。
* **UI 组件库**：Element Plus
* **状态管理**：Pinia
* **网络请求**：Axios (需封装统一拦截器)。
* **数据可视化**：ECharts
* **后端**：**Go (优先使用 Gin )**。
* **依赖管理**：Go Modules (`go.mod`)。
* **依赖注入**：优先使用 `google/wire` 进行静态注入。
* 若无特殊说明，均遵循以上技术选型

## 三、通用代码规范

### 命名规范

* **变量 / 函数 / 方法**：`camelCase`（遵循 Go 的导出规则：首字母大写对包外可见）。
* **结构体 / 接口**：`PascalCase`。
* **常量**：`UPPER_SNAKE_CASE` 或 `PascalCase`（若需导出）。
* **文件 / 文件夹**：`snake_case` (符合 Go 官方社区习惯，与原规则 kebab-case 区分)。
* **接口命名**：简单接口以 `er` 结尾（如 `Reader`, `Writer`）。

命名应语义清晰，禁止随意缩写。

### 注释规范（强制）

* 注释用于解释「为什么这样设计」，而不是代码字面含义
* 复杂逻辑、业务判断、边界条件必须写注释
* 禁止无意义注释

统一注释标记：

```ts
// TODO: 待实现功能
// FIXME: 已知问题或潜在缺陷
// NOTE: 重要设计说明
// HACK: 临时方案，后续必须重构
```

#### 函数注释规范

前端（JSDoc）：

```ts
/**
 * 获取用户信息
 * @param userId 用户 ID
 * @returns 用户数据
 */
```

go：

```go
// GetUserByID 根据用户 ID 获取用户信息
// 如果用户不存在，将返回 gorm.ErrRecordNotFound
func GetUserByID(id string) (*User, error) {
    // ...
}
```

## 四、前端规范（Vue3.5+）

### 基本原则

* 核心语法：必须使用 Composition API 和 <script setup> 语法。
* 单一职责：单个组件只承担单一职责，展示逻辑与业务逻辑分离。
* 逻辑复用：可复用逻辑必须抽离为 Composables (自定义 Hook)。

### 命名约定

* 组件命名：组件名使用 PascalCase，文件名与组件名保持一致
* Composables：自定义 Hook 文件及函数必须以 use 开头（如 useUserData.ts）
* 文件夹：遵循 snake_case 约定

```ts
// 组件文件：UserCard.vue（PascalCase）
// Composable 文件：useUserData.ts
```

### Composables 使用规范

* 只能在 `<script setup>` 或 `setup()` 函数顶层调用
* 不允许在条件、循环中调用
* 一个 Composable 只处理一种职责

### TypeScript 约束

- **强类型声明**：所有 Props、Emits、Pinia State 必须通过 `interface` 定义类型。
- **禁止 any**：除非处理未知第三方数据，否则严禁使用 `any`。

### Props 规范

* 必须使用 TypeScript 类型定义
* 使用解构方式接收 props
* 非必传参数使用 `?`

```ts
interface UserCardProps {
  user: User
  onClick?: () => void
}
```
### 状态管理 (Pinia) 与 请求 (Axios)
* 状态定义：优先使用 Setup Store 模式定义 Store。
* 网络请求：Axios 实例必须封装在 pkg/api 或 service/ 中，禁止在组件内直接写 BaseURL。
* 类型映射：API 请求的 Response 必须定义对应的 TypeScript Interface。

### 可视化与性能 (ECharts)

- **实例销毁**：必须在 `onUnmounted` 中调用 `chart.dispose()` 防止内存泄漏。
- **响应式**：必须监听窗口 resize 事件或使用 `ResizeObserver` 自动重绘。

### 路由与加载

- **异步加载**：页面级组件必须使用 `() => import('./View.vue')`。
- **权限控制**：路由跳转需在 `router.beforeEach` 中结合 Pinia 权限标识进行拦截


### 性能与结构要求

* 避免不必要的重复渲染
* 合理使用 computed / shallowRef 减少响应式开销
* 列表渲染必须提供稳定的 key
* 大数据列表使用虚拟滚动
* 路由与组件支持懒加载

## 五、后端规范（Go）

### 基本要求

* **Go 版本**：≥ 1.21。
* **类型检查**：必须强制处理所有 `error`，禁止使用 `_` 忽略非平凡错误。
* **并发控制**：禁止在请求上下文之外开启不受控的 Goroutine；必须显式传递 `context.Context`。
* **日志**：使用 `slog` (标准库) 或 `zerolog`，禁止使用 `fmt.Print`。

### 分层结构（必须遵守）

采用 **Clean Architecture** 变体，严格禁止层级循环依赖：

```
项目目录/
├── cmd/
│   └── server/
│       └── main.go        # 程序入口
├── internal/
│   ├── config/            # 配置加载
│   ├── logger/            # 日志封装
│   ├── router/            # 路由注册
│   ├── handler/           # HTTP / RPC 接口
│   ├── service/           # 业务逻辑，不涉及 SQL
│   ├── repo/              # 数据访问/缓存 CRUD 实现 (Model/ORM)
│   └── model/             # 数据模型
|   └── middleware/        # 中间件
├── pkg/                   # 第三方工具类封装或者可复用公共库
├── configs/               # 配置文件
├── scripts/               # 放各种执行脚本
├── doc/				   # 放文档
├── go.mod
├── Makefile
└── Dockerfile
```

### 错误处理规范

* 错误应被 Wrap (包装) 携带调用栈信息，而不是简单的字符串。
* `if err != nil { return err }` 需保持代码缩进清晰（Happy Path 逻辑靠左）。
* 日志中不得包含敏感信息

## 六、安全规范（重点）

### 通用安全原则

* 永远不信任客户端输入
* 所有输入必须进行校验
* 敏感操作，所有写操作（POST/PUT/DELETE）必须经过身份与权限校验

### 前端安全

* 禁止使用 `v-html` 渲染未经过滤的用户输入
* 防止 XSS / CSRF 攻击
* 不在前端存储敏感信息
* Token 推荐使用 HttpOnly Cookie

### 后端安全

* **输入校验**：使用 `validator` 包对结构体进行 Tag 校验。
* **配置管理**：敏感密钥 **必须** 通过环境变量或 `.env` 加载，严禁硬编码。
* **数据库**：强制使用参数化查询（ORM 默认支持），杜绝 SQL 注入。
* **脱敏**：返回给前端的 JSON 结构体必须使用 `json:"-"` 屏蔽敏感字段（如密码哈希）。
* 敏感字段返回前需脱敏
* 密码等敏感数据必须加密存储

## 七、AI 协作使用规范

* 所有自动生成的代码必须遵守本规则
* 生成结果应：

  * 结构清晰
  * 类型完整
  * 可维护
  * 安全
* 不生成不必要的复杂实现

## 八、双 AI 协作工程流程（Codex + Codex）

本项目采用 **Codex（规划+编码+决策）** 与 **Codex（分析+审查+评分）** 协作的工程流水线。所有中间产物统一存放在 `.Codex/` 目录下，便于追溯与复盘。

### 工作目录结构

```
.Codex/
├── context-initial.json       # 阶段0：Codex 初步上下文扫描结果
├── context-question-N.json    # 阶段0：针对高优疑问的深挖结果（N=1,2,3）
├── coding-progress.json       # 阶段2：Codex 编码进度记录
├── operations-log.md          # 全程：决策与操作记录
└── review-report.md           # 阶段3：Codex 审查报告与评分
```

### 流程总览

```
阶段0 需求理解
  └─ Codex 扫描上下文 → context-initial.json
  └─ Codex 识别高优疑问 → Codex 深挖 → context-question-N.json
  └─ 充分性检查（5项全通过才进入下一阶段）

阶段1 任务规划（Codex）
  └─ sequential-thinking 思考整体方案
  └─ shrimp-task-manager 拆解子任务（含顺序、依赖、验收标准）
  └─ 同步写入 operations-log.md

阶段2 编码执行（Codex 主导）
  └─ 简单逻辑（<10行）：直接 Read/Edit/Write
  └─ 复杂逻辑（>10行）：先让 Codex 出算法/架构方案，再编码
  └─ 持续更新 coding-progress.json

阶段3 审查与评分（Codex）
  └─ Codex 先写审查清单
  └─ Codex 用 sequential-thinking 做批判性分析，双维度打分
  └─ 输出综合评分（0-100）+ 建议，写入 review-report.md

Codex 决策规则
  └─ ≥90 分且"通过"：直接确认
  └─ <80 分且"退回"：退回修复，Codex 给出方向
  └─ 80-89 分或"需讨论"：Codex 结合业务做最终判断，记录原因
```

### 充分性检查（进入阶段1前必须全部通过）

1. 能定义出清晰的输入输出 / 接口契约？
2. 理解当前技术选型的理由？
3. 识别出主要风险点（并发 / 边界 / 性能）？
4. 知道如何验证（测试框架、怎么跑）？
5. 上下文收集次数 ≤ 3 次？

### 职责边界

* **不要**让 Codex 做简单代码编写（浪费 MCP 往返）
* **不要**让 Codex 替代 Codex 的深度分析与审查
* 每完成一个子任务立即跑测试/构建，不攒到最后
* 所有"为什么这样改""为什么不用备选方案"必须记入 `operations-log.md`

 *   * * 4N�e�e�Nĉ�* * ��NUO(u�N�c�g��v4N�eKmՋ�Nxb�,g��Y  t e s t _ * . g o ,   f i n d _ * . p y   I{	��_{�X[>e(WS_MR�vU_�v  	 m p /   �vU_N�%N�y�v�c>e(Wy��v9h�vU_�N�MQы�Q�z0
 
 