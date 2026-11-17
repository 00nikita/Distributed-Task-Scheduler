package task

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
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


func (r *Repository) ClaimDueTasks(ctx context.Context, ([]Task, error) ) {
    query := `
       WITH due_tasks AS (
           SELECT id
           FROM tasks
           WHERE Status = 'PENDING'
             AND scheduled_at <= NOW()
           ORDER BY scheduled_at, id 
           LIMIT 10
           FOR UPDATE SKIP LOCKED
       )
       UPDATE tasks
       SET status = 'RUNNING',
          updated_at = NOW()
       WHERE id IN (SELECT id FROM due_tasks)
       RETURNING id, task_type, payload, status, scheduled_at
    `
        rows, err := r.db.Query(ctx, query)
    ifr err!=nil {
        return nil, err
    }

    defer rows.Close()

    var tasks []Task

    for rows.Next() {
        var t Task

        err := rows.Scan(
            &t.ID,
            &t.TaskType,
            &t.Payload,
            &t.Status,
            &t.ScheduledAt,
        )

        if err!=nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }

        if err := rows.Err(); err != nil {
        return nil, err
    }

    return tasks, nil

}
