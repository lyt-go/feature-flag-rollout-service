# 特性开关 / 灰度发布服务（feature-flag）

一个基于 Go 标准库实现的特性开关（Feature Flag）与灰度发布后端服务，支持开关全生命周期管理、变体与目标分组、灰度规则评估、变更审计与评估统计。

## 技术栈

- 纯 Go 标准库（`net/http` + 标准库），零第三方依赖
- 内存存储（`sync.RWMutex` 保证并发安全）
- 标准分层架构：`cmd` / `internal`（app/config/model/store/service/handler）/ `pkg`

## 运行

```bash
go run ./cmd/server
# 或
go build -o featureflag-server ./cmd/server && ./featureflag-server
```

默认监听 `:8080`，可通过环境变量配置：

| 环境变量 | 默认值 | 说明 |
|---------|-------|------|
| `PORT` | `8080` | 监听端口 |
| `ADDR` | `:8080` | 完整监听地址（优先级高于 PORT） |
| `MAX_PAGE_SIZE` | `100` | 分页最大页大小 |
| `DEFAULT_ENABLED` | `false` | 未找到开关时的默认返回 |
| `LOG_LEVEL` | `info` | 日志级别：debug/info/warn/error |

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误码：400 校验失败、404 不存在、409 冲突、500 服务器错误。

### 开关 Flag

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/flags` | 创建开关（key/name/description/enabled） |
| GET | `/api/flags` | 列表（status/keyword + 分页） |
| GET | `/api/flags/{id}` | 详情 |
| PUT | `/api/flags/{id}` | 更新（name/description，`?operator=` 指定操作人） |
| DELETE | `/api/flags/{id}` | 删除 |
| POST | `/api/flags/{id}/toggle` | 翻转启用状态 |
| POST | `/api/flags/{id}/status` | 变更状态（active/paused/archived，状态机校验） |

### 变体 Variant

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/variants` | 创建变体（校验开关存在） |
| GET | `/api/variants` | 列表（flag_id/keyword + 分页） |
| GET | `/api/variants/{id}` | 详情 |
| PUT | `/api/variants/{id}` | 更新 |
| DELETE | `/api/variants/{id}` | 删除 |

### 目标分组 TargetGroup

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/target-groups` | 创建分组（rules 为 JSON 对象，如 `{"city":"beijing"}`） |
| GET | `/api/target-groups` | 列表（keyword + 分页） |
| GET | `/api/target-groups/{id}` | 详情 |
| PUT | `/api/target-groups/{id}` | 更新 |
| DELETE | `/api/target-groups/{id}` | 删除 |

### 灰度规则 RolloutRule

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/rollout-rules` | 创建规则（校验开关与变体归属一致） |
| GET | `/api/rollout-rules` | 列表（flag_id/status + 分页） |
| GET | `/api/rollout-rules/{id}` | 详情 |
| PUT | `/api/rollout-rules/{id}` | 更新 |
| DELETE | `/api/rollout-rules/{id}` | 删除 |

### 变更审计 ChangeLog

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/change-logs` | 列表（flag_id/action + 分页） |
| GET | `/api/change-logs/{id}` | 详情 |

### 评估 Evaluation

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/evaluate` | 评估（flag_key/target_key/context） |
| GET | `/api/evaluation-records` | 评估记录列表（flag_id/target_key/result_variant + 分页） |

### 统计 Stats

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/overview` | 全局统计 |

## 核心实体

1. **Flag**：特性开关，状态机 `active ↔ paused / active → archived / archived → active`。
2. **Variant**：开关取值变体（name/payload/weight）。
3. **TargetGroup**：目标分组，通过规则 JSON 描述匹配条件。
4. **RolloutRule**：灰度规则，命中目标分组后按比例返回指定变体。
5. **ChangeLog**：开关变更审计记录（create/update/toggle/pause/resume/archive/restore）。
6. **EvaluationRecord**：每次开关评估的命中记录。

## 评估逻辑

1. 按 `flag_key` 定位开关；不存在返回 404。
2. 开关状态非 `active` 或 `Enabled=false` 时，返回 `enabled=false`。
3. 否则按优先级遍历启用的灰度规则：若规则含目标分组，需上下文命中该分组；命中后按 `target_key` 哈希比例（0-99 < percentage）决定是否返回变体。
4. 每次评估都会落一条 `EvaluationRecord`。

## 测试

```bash
go test ./...
```
