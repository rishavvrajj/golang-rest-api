package router

import (
	"net/http"

	"github.com/rishavvrajj/golang-rest-api/internal/http/handlers/student"
	"github.com/rishavvrajj/golang-rest-api/internal/http/middleware"
	"github.com/rishavvrajj/golang-rest-api/internal/http/response"
	"github.com/rishavvrajj/golang-rest-api/internal/service"
)

// NewRouter sets up the HTTP router (ServeMux) and registers all application routes with middleware.
func NewRouter(studentService *service.StudentService) http.Handler {
	mux := http.NewServeMux()

	studentHandler := student.NewStudentHandler(studentService)

	// Student CRUD routes
	mux.HandleFunc("POST /students", studentHandler.Create())
	mux.HandleFunc("GET /students", studentHandler.GetList())
	mux.HandleFunc("GET /students/{id}", studentHandler.GetById())
	mux.HandleFunc("PUT /students/{id}", studentHandler.Update())
	mux.HandleFunc("DELETE /students/{id}", studentHandler.Delete())

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJson(w, http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "Server is running successfully",
		})
	})

	// Wrap mux with global middleware (Logging, Recovery)
	var handler http.Handler = mux
	handler = middleware.Recovery(handler)
	handler = middleware.Logging(handler)

	return handler
}
