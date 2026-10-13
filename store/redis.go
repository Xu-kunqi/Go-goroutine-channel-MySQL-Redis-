package store

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/Xu-kunqi/go-task-scheduler/scheduler"
    "github.com/redis/go-redis/v9"
)

type RedisStore struct {
    client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
    return &RedisStore{
        client: redis.NewClient(&redis.Options{Addr: addr}),
    }
}

func (r *RedisStore) SetTask(ctx context.Context, task *scheduler.Task) error {
    data, _ := json.Marshal(task)
    return r.client.Set(ctx, "task:"+itoa(task.ID), data, 10*time.Minute).Err()
}

func (r *RedisStore) GetTask(ctx context.Context, id int) (*scheduler.Task, error) {
    val, err := r.client.Get(ctx, "task:"+itoa(id)).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    var task scheduler.Task
    json.Unmarshal([]byte(val), &task)
    return &task, nil
}

func (r *RedisStore) AcquireLock(ctx context.Context, key string) (bool, error) {
    return r.client.SetNX(ctx, "lock:"+key, "locked", 10*time.Second).Result()
}

func (r *RedisStore) ReleaseLock(ctx context.Context, key string) error {
    return r.client.Del(ctx, "lock:"+key).Err()
}

func itoa(i int) string {
    return fmt.Sprintf("%d", i)
}
