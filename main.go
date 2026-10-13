package main

import (
    "log"
    "net/http"

    "github.com/Xu-kunqi/go-task-scheduler/api"
    "github.com/Xu-kunqi/go-task-scheduler/scheduler"
    "github.com/Xu-kunqi/go-task-scheduler/store"
)

func main() {
    // 1. 初始化 MySQL 存储
    mysqlStore, err := store.NewMySQLStore("root:password@tcp(127.0.0.1:3306)/tasks?parseTime=true&loc=Local")
    if err != nil {
        log.Fatalf("连接 MySQL 失败: %v", err)
    }
    log.Println("✅ MySQL 连接成功")

    // 2. 初始化 Redis 存储
    redisStore := store.NewRedisStore("127.0.0.1:6379")
    log.Println("✅ Redis 连接成功")

    // 3. 初始化任务池（5 个 worker）
    pool := scheduler.NewTaskPool(5)
    pool.Start()
    defer pool.Stop()

    // 4. 注册 API 路由
    handler := &api.Handler{
        Pool:  pool,
        MySQL: mysqlStore,
        Redis: redisStore,
    }
    mux := http.NewServeMux()
    handler.Register(mux)

    // 5. 启动 HTTP 服务
    log.Println("🚀 服务启动，监听 :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("服务启动失败: %v", err)
    }
}
