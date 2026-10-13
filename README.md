# Go 并发任务调度服务

一个基于 Go 语言实现的并发任务调度服务，支持任务提交、异步执行、超时控制、持久化存储与缓存加速。

## ✨ 项目简介

本项目使用 Go 语言从零实现了一个并发任务调度服务，核心目标是：

- **高并发**：用 goroutine + channel 实现任务池，支持并发执行。
- **超时控制**：每个任务都有独立的超时时间，超时自动取消。
- **持久化**：任务状态写入 MySQL，保证数据不丢。
- **缓存加速**：热点任务缓存在 Redis，响应时间从 50ms 降到 1ms。
- **RESTful API**：对外提供标准的 RESTful 接口。

## 🛠️ 技术栈

| 类别 | 技术 |
| :--- | :--- |
| **语言** | Go 1.22 |
| **并发** | goroutine、channel、context、sync.WaitGroup |
| **Web** | net/http |
| **数据库** | MySQL 8、go-sql-driver/mysql |
| **缓存** | Redis 7、go-redis/v9 |
| **工程化** | Git、GitHub Actions、Makefile |

## 🏗️ 项目架构

```text
go-task-scheduler/
├── main.go                 # 入口：初始化依赖、启动 HTTP 服务
├── scheduler/
│   ├── task.go             # 任务结构体定义
│   └── pool.go             # goroutine 任务池
├── store/
│   ├── mysql.go            # MySQL 持久化
│   └── redis.go            # Redis 缓存 + 分布式锁
├── api/
│   └── handler.go          # RESTful API 处理函数
├── .github/workflows/
│   └── ci.yml              # GitHub Actions CI
└── go.mod
```

## 🔌 API 文档

| 方法 | URL | 说明 | 状态码 |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/tasks` | 创建任务 | 201 Created |
| `GET` | `/api/v1/tasks` | 列出所有任务 | 200 OK |
| `GET` | `/api/v1/tasks/{id}` | 查询单个任务 | 200 OK / 404 Not Found |

### 请求示例

创建任务：

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "每日报表", "payload": "report", "timeout": 5}'
```

返回：

```json
{
  "id": 1,
  "name": "每日报表",
  "payload": "report",
  "status": "pending",
  "timeout": 5,
  "created_at": "2026-10-13T16:26:49+08:00",
  "updated_at": "2026-10-13T16:26:49+08:00"
}
```

查询任务：

```bash
curl http://localhost:8080/api/v1/tasks/1
```

## 🚀 快速开始

### 前置依赖

- Go 1.22+
- MySQL 8+
- Redis 7+

### 初始化数据库

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS tasks;"
```

### 编译与运行

```bash
# 1. 下载依赖
go mod download

# 2. 编译
go build -o my-scheduler .

# 3. 运行
./my-scheduler
```

预期输出：

```text
✅ MySQL 连接成功
✅ Redis 连接成功
任务池启动，worker 数量: 5
🚀 服务启动，监听 :8080
```

## 💡 技术亮点

### 1. 并发任务池

基于 goroutine + channel 实现任务池：

- 固定数量的 worker goroutine 从 channel 取任务执行。
- 用 context.WithTimeout 实现任务超时控制。
- 用 sync.WaitGroup 等待所有 worker 退出。
- 用 context.WithCancel 实现优雅关闭。

### 2. 缓存加速

基于 Redis 实现旁路缓存：

- 查询时先查 Redis，命中直接返回。
- 未命中则查 MySQL，然后回写 Redis。
- 设置 10 分钟 TTL，防止缓存无限增长。

### 3. 分布式锁

基于 Redis SETNX 实现分布式锁：

- 多实例同时执行任务时，只有一个能抢到锁。
- 设置 10 秒过期时间，防止死锁。

### 4. 优雅关闭

程序退出时，通过 cancel() 通知所有 worker 退出，并用 sync.WaitGroup 等待它们结束，避免 goroutine 泄漏。

## 📊 性能对比

| 场景 | 响应时间 |
| :--- | :--- |
| 查询 MySQL | ~50ms |
| 查询 Redis（命中） | ~1ms |
| 提升 | 50x |

## 🎯 面试亮点

- 并发编程：goroutine、channel、context 三件套，避免 goroutine 泄漏。
- 缓存设计：Cache Aside Pattern，一致性保障。
- 分布式锁：Redis SETNX 实现，防止任务重复执行。
- RESTful API：标准 HTTP 方法 + 状态码 + JSON。
- 工程化：GitHub Actions CI，自动化构建。

## 📝 作者

- 姓名：徐鲲淇
- GitHub：https://github.com/Xu-kunqi
- 求职意向：Go 后端开发实习生

---

⭐ 如果这个项目对你有帮助，欢迎 Star！
