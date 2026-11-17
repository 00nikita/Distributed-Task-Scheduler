package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/00nikita/distributed-task-scheduler/internal/task"
	"github.com/00nikita/distributed-task-scheduler/internal/worker"
)

type scheduler struct {
	repo *task.Repository
	worker *worker.Worker
}

func NewScheduler( repo *task.Repository, w *worker.Worker) *Scheduler{
	return &Scheduler{repo: repo, worker: w,}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			tasks, err:= s.repo.ClaimDueTasks(ctx)
			if err != nil {
				log.Println("Failed to fetch tasks")
				continue
			}

			for _, t := range tasks{
				err := s.worker.Execute(ctx, t)
				if err != nil {
					log.Printf("Task %d failed: %v", t.ID, err)
					continue
				}

				log.Printf("Task %d executed successfully", t.ID)
			}
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		}
	}
}