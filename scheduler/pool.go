package scheduler

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type TaskPool struct {
    taskChan chan *Task
    workers  int
    wg       sync.WaitGroup
    ctx      context.Context
    cancel   context.CancelFunc
}

func NewTaskPool(workers int) *TaskPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &TaskPool{
        taskChan: make(chan *Task, 100),
        workers:  workers,
        ctx:      ctx,
        cancel:   cancel,
    }
}

// 启动 worker 池
func (p *TaskPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
    fmt.Printf("任务池启动，worker 数量: %d\n", p.workers)
}

// 单个 worker：不断从 channel 里取任务
func (p *TaskPool) worker(id int) {
    defer p.wg.Done()
    for {
        select {
        case task, ok := <-p.taskChan:
            if !ok {
                return
            }
            p.execute(task)
        case <-p.ctx.Done():
            return
        }
    }
}

// 执行任务，带超时控制
func (p *TaskPool) execute(task *Task) {
    ctx, cancel := context.WithTimeout(p.ctx, time.Duration(task.Timeout)*time.Second)
    defer cancel()

    task.Status = StatusRunning
    task.UpdatedAt = time.Now()

    done := make(chan error, 1)
    go func() {
        // 模拟任务执行（实际项目中替换为真实逻辑）
        time.Sleep(time.Duration(task.Timeout) * time.Second)
        done <- nil
    }()

    select {
    case <-done:
        task.Status = StatusSuccess
    case <-ctx.Done():
        task.Status = StatusTimeout
    }
    task.UpdatedAt = time.Now()
}

// 提交任务
func (p *TaskPool) Submit(task *Task) {
    p.taskChan <- task
}

// 关闭任务池
func (p *TaskPool) Stop() {
    p.cancel()
    close(p.taskChan)
    p.wg.Wait()
}
