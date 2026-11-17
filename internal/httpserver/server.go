package httpserver

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)
func New(handler http.Handler) *http.Server{
	return &http.Server{
		Addr: ":8080",
		Handler: recoverer(handler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
func recoverer(handler http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(){
			if value := recover(); value != nil {
				slog.Error("http handler panic", "panic", value, "stack", string(debug.Stack()))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w,r)
	})
}