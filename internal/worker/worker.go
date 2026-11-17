package worker

import (
   "context"
   "log"
   "time"

   "github.com/00nikita/distributed-task-scheduler/internal/task"
)

type Worker struct{}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) Execute(ctx context.Context, t task.Task) error {
	log.Printf("Executing task: ID=%d Type=%s", t.ID, t.TaskType)
	select {
    case <-time.After(2 * time.Second):
        log.Printf("Finished task: ID=%d", t.ID)
        return nil

    case <-ctx.Done():
        return ctx.Err()
    }
}