package main

import (
	"net/http"
	"encoding/json"
	"log"
	"context"
	"os"
	"github.com/00nikita/distributed-task-scheduler/internal/database"
	"github.com/00nikita/distributed-task-scheduler/internal/task"
	"github.com/00nikita/distributed-task-scheduler/internal/scheduler"
	"github.com/00nikita/distributed-task-scheduler/internal/worker"
)

func main(){
	mux:= http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	database_url := os.Getenv("DATABASE_URL")
	if database_url == ""{
		log.Fatal("Database url is not set")
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, database_url)
	if err!=nil{
		log.Fatal("Failed to connect to database", err)
	}

	defer pool.Close(ctx)

    log.Println("Connected to PostgreSQL successfully")

    taskRepo := task.NewRepository(pool)
	taskHandler := task.NewHandler(taskRepo)

    mux.HandleFunc("POST /tasks", taskHandler.CreateTask)

	taskWorker := worker.NewWorker()

    taskScheduler := scheduler.NewScheduler(taskRepo, taskWorker)

	go taskScheduler.Start(ctx)
	log.Println("starting server in 8080")
	if err := http.ListenAndServe(":8080", mux); err !=nil{
		log.Fatal(err)
	}
}