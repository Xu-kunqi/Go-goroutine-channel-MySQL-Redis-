package scheduler

import "time"

type Status string

const (
    StatusPending  Status = "pending"
    StatusRunning  Status = "running"
    StatusSuccess  Status = "success"
    StatusFailed   Status = "failed"
    StatusTimeout  Status = "timeout"
)

type Task struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Payload   string    `json:"payload"`
    Status    Status    `json:"status"`
    Timeout   int       `json:"timeout"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
