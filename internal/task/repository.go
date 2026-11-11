package task

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}


func (r *Repository) CreateTask(ctx context.Context, task *Task) (int64, error) {
    query := `
        INSERT INTO tasks (task_type, payload, scheduled_at)
        VALUES ($1, $2, $3)
        RETURNING id
    `

    var id int64

    err := r.db.QueryRow(
        ctx,
        query,
        task.TaskType,
        task.Payload,
        task.ScheduledAt,
    ).Scan(&id)

    if err != nil {
        return 0, err
    }

    return id, nil
}
