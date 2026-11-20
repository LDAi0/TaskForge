package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"taskforge/internal/api"
	"taskforge/internal/task"
	"uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Handler struct{
	service *task.Service
}

func New(service *task.Service) http.Handler{
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

func (h *Handler) notFound (w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "Not found", "No endpoint matches this request")
}

func (h *Handler) getHealth (w http.ResponseWriter, r *http.Request) {
	t, err := h.service.HealthStatus(r.Context())
	if err!=nil{
		h.writeServiceError(w, err)
	}
	response := api.HealthResponse{
		ResponseTimeMs: t.Response_Time_Ms,
		Status: api.HealthResponseStatus(t.Status),
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) postTask (w http.ResponseWriter, r *http.Request) {
	var request api.PostTaskRequest
	if err:=decodeJSON(w,r,&request); err!=nil{
		writeError(w,http.StatusBadRequest,"Invalid body request", err.Error())
		return
	}

	t, err := h.service.Post(r.Context(),task.PostInput{
		Parameters: task.Parameters(request.Parameters),
		ReportType: string(request.ReportType),
	})
	if err!=nil{
		h.writeServiceError(w,err)
	}
	id, err := uuid.Parse(t.ID)
	if err!=nil{
		h.writeServiceError(w, task.ErrInvalidID)
		return 
	}
	response := api.PostTaskResponse{
		CreatedAt: t.Created_at,
		Id: openapi_types.UUID(id),
		Status: api.PostTaskResponseStatus(t.Status),
	}
	writeJSON(w,http.StatusOK,response)

}

func (h *Handler) getTask (w http.ResponseWriter, r *http.Request) {
	t, err := h.service.Get(r.Context(), r.PathValue("id"))
	if err !=nil{
		h.writeServiceError(w, err)
		return
	}
	id, err := uuid.Parse(t.ID)
	if err!=nil{
		h.writeServiceError(w, task.ErrInvalidID)
		return 
	}
	response:=api.GetTaskResponse{
		Id: openapi_types.UUID(id),
		Status: api.GetTaskResponseStatus(t.Status),
		Attempts: int(t.Attempts),
		CreatedAt: t.Created_at,
		StartedAt: &t.Started_at,
		FinishedAt: &t.Finished_at,
		Result: (*api.Result)(&t.Result),
		Error:	(*api.TaskError)(&t.Error),
	}
	writeJSON(w,http.StatusOK,response)
}

func (h *Handler) writeServiceError(w http.ResponseWriter, err error){
	switch{
	case errors.Is(err, task.ErrTaskNotFound):
		writeError(w,http.StatusNotFound,"Task not found", "Task with this ID not found")
	case errors.Is(err, task.ErrInvalidPeriodDays):
		writeError(w, http.StatusBadRequest, "Invalid Period days", "Period days must be between 1 and 365")
	case errors.Is(err, task.ErrInvalidReportType):
		writeError(w, http.StatusBadRequest, "Invalid Report Type", "Report type must be in [summary]")
	case errors.Is(err, task.ErrTaskAlreadyExist):
		writeError(w, http.StatusConflict, "Task already exist", "Posted task already exist!")
	case errors.Is(err, task.ErrInvalidID):
		writeError(w, http.StatusBadRequest, "Invalid ID type", "ID type must be UUID!")
	default:
		slog.Error("request failed", "error", err)
		writeError(w, http.StatusInternalServerError, "Internal server error", "The request could not be completed")
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, value any) error{
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err!=nil || contentType != "application/json"{
		return fmt.Errorf("Content-Type error(must be application/json): %w", err)
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body); decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err!=nil{
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF{
		if err == nil{
			return fmt.Errorf("request body must contain a single JSON object")
		}
		return err
	}
	return nil
}

func writeError (w http.ResponseWriter, status int, message, details string){
	response := api.ErrorResponse{
		Message: message,
		Details: &details,
	}
	writeJSON(w, status, response)
}

func writeJSON (w http.ResponseWriter, status int, value any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil{
		slog.Error("write response", "error", err)
	}
}