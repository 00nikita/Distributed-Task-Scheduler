package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/00nikita/distributed-task-scheduler/internal/task"
)

type scheduler struct {
	repo *task.Repository
}

func NewScheduler( repo *task.Repository) *Scheduler{
	return &Scheduler{repo: repo}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			tasks, err:= s.repo.GetDueTasks(ctx)
			if err != nil {
				log.Println("Failed to fetch tasks")
				continue
			}

			for _, t := range tasks{
				log.Printf("Due task found: ID=%d Type=%s", t.ID, t.TaskType)
			}
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		}
	}
}