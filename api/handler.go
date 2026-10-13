package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/Xu-kunqi/go-task-scheduler/scheduler"
    "github.com/Xu-kunqi/go-task-scheduler/store"
)

type Handler struct {
    Pool  *scheduler.TaskPool
    MySQL *store.MySQLStore
    Redis *store.RedisStore
}

func (h *Handler) Register(mux *http.ServeMux) {
    mux.HandleFunc("/api/v1/tasks", h.handleTasks)
    mux.HandleFunc("/api/v1/tasks/", h.handleTaskByID)
}

// POST /api/v1/tasks  - 创建任务
// GET  /api/v1/tasks  - 列出所有任务
func (h *Handler) handleTasks(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        var task scheduler.Task
        if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
            http.Error(w, "参数错误", http.StatusBadRequest)
            return
        }
        task.Status = scheduler.StatusPending
        task.CreatedAt = time.Now()
        task.UpdatedAt = time.Now()

        if err := h.MySQL.Create(&task); err != nil {
            http.Error(w, "创建任务失败", http.StatusInternalServerError)
            return
        }

        // 提交到任务池
        h.Pool.Submit(&task)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(task)

    case http.MethodGet:
        tasks := h.MySQL.List()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tasks)

    default:
        http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
    }
}

// GET    /api/v1/tasks/{id} - 查询单个任务
// DELETE /api/v1/tasks/{id} - 删除任务
func (h *Handler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "无效的任务 ID", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        // 1. 先查 Redis 缓存
        task, _ := h.Redis.GetTask(r.Context(), id)
        if task != nil {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(task)
            return
        }

        // 2. 缓存未命中，查 MySQL
        task = h.MySQL.Get(id)
        if task == nil {
            http.Error(w, "任务不存在", http.StatusNotFound)
            return
        }

        // 3. 回写 Redis
        h.Redis.SetTask(r.Context(), task)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(task)

    default:
        http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
    }
}
