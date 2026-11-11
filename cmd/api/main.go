package main

import (
	"net/http"
	"encoding/json"
	"log"
	"context"
	"os"
	"github.com/00nikita/distributed-task-scheduler/internal/database"
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
	conn, err := database.Connect(ctx, database_url)
	if err!=nil{
		log.Fatal("Failed to connect to database", err)
	}

	defer conn.Close(ctx)

    log.Println("Connected to PostgreSQL successfully")
	log.Println("starting server in 8080")
	if err := http.ListenAndServe(":8080", mux); err !=nil{
		log.Fatal(err)
	}
}