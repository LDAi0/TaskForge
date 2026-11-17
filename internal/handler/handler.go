package handler

import (
	"net/http"
	"taskforge/internal/api"
	"taskforge/internal/task"
)

type Handler struct{
	service *task.Service
}

func New(service *task.Service) *Handler{
	h := &Handler{
		service: service,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", h.getTask)
	mux.HandleFunc("POST /tasks", h.postTask)
	mux.HandleFunc("GET /health", h.getHealth)
	mux.HandleFunc("/", h.notFound)
	return mux
}