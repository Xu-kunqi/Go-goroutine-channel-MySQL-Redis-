package store

import (
    "database/sql"
    "time"
    "log"

    "github.com/Xu-kunqi/go-task-scheduler/scheduler"
    _ "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
    db *sql.DB
}

func NewMySQLStore(dsn string) (*MySQLStore, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255),
            payload TEXT,
            status VARCHAR(50),
            timeout INT,
            created_at DATETIME,
            updated_at DATETIME
        )
    `)
    return &MySQLStore{db: db}, err
}

func (s *MySQLStore) Create(task *scheduler.Task) error {
    res, err := s.db.Exec(
        "INSERT INTO tasks (name, payload, status, timeout, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
        task.Name, task.Payload, task.Status, task.Timeout, task.CreatedAt, task.UpdatedAt,
    )
    if err != nil {
        return err
    }
    id, _ := res.LastInsertId()
    task.ID = int(id)
    return nil
}

func (s *MySQLStore) UpdateStatus(id int, status scheduler.Status) error {
    _, err := s.db.Exec(
        "UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?",
        status, time.Now(), id,
    )
    return err
}

func (s *MySQLStore) Get(id int) *scheduler.Task {
    var task scheduler.Task
    var statusStr string

    err := s.db.QueryRow(
        "SELECT id, name, payload, status, timeout, created_at, updated_at FROM tasks WHERE id = ?",
        id,
    ).Scan(&task.ID, &task.Name, &task.Payload, &statusStr, &task.Timeout, &task.CreatedAt, &task.UpdatedAt)

    if err != nil {
        log.Printf("Get 查询失败: %v", err)
        return nil
    }

    task.Status = scheduler.Status(statusStr)
    return &task
}

func (s *MySQLStore) List() []*scheduler.Task {
    rows, err := s.db.Query(
        "SELECT id, name, payload, status, timeout, created_at, updated_at FROM tasks ORDER BY id DESC",
    )
    if err != nil {
        return nil
    }
    defer rows.Close()

    var tasks []*scheduler.Task
    for rows.Next() {
        var task scheduler.Task
        var statusStr string
        err := rows.Scan(&task.ID, &task.Name, &task.Payload, &statusStr, &task.Timeout, &task.CreatedAt, &task.UpdatedAt)
        if err != nil {
            log.Printf("List 扫描失败: %v", err)
            continue
        }
        task.Status = scheduler.Status(statusStr)
        tasks = append(tasks, &task)
    }
    return tasks
}
