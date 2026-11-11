package task

import (
	"encoding/json"
	"net/http"
	"time"
)

type CreateTaskRequest struct {
	TaskType string `json:"task_type"`
	Payload json.RawMessage `json:"payload"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type Handler struct {
    repo *Repository
}

func NewHandler(repo *Repository) *Handler {
    return &Handler{repo: repo}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {

	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
	http.Error(w, "Invalid JSON", http.StatusBadRequest)
	return
    }

task := &Task{
    TaskType:    req.TaskType,
    Payload:     req.Payload,
    ScheduledAt: req.ScheduledAt,
}

id, err := h.repo.CreateTask(r.Context(), task)

if err != nil {
    http.Error(w, "Failed to create task", http.StatusInternalServerError)
    return
}

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)

json.NewEncoder(w).Encode(map[string]int64{
    "id": id,
})

}