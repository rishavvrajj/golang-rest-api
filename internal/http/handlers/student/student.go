package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rishavvrajj/golang-rest-api/internal/http/response"
	"github.com/rishavvrajj/golang-rest-api/internal/models"
	"github.com/rishavvrajj/golang-rest-api/internal/service"
)

type StudentHandler struct {
	service *service.StudentService
}

func NewStudentHandler(s *service.StudentService) *StudentHandler {
	return &StudentHandler{service: s}
}

type CreateResponse struct {
	response.Response
	StudentId int64 `json:"student_id,omitempty"`
}

func (h *StudentHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling create student request")

		var student models.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("request body is empty")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid request body payload")))
			return
		}

		// Validate request struct
		if err := validator.New().Struct(student); err != nil {
			var validateErrs validator.ValidationErrors
			if errors.As(err, &validateErrs) {
				response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		studentId, err := h.service.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			slog.Error("failed to create student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(errors.New("failed to create student")))
			return
		}

		response.WriteJson(w, http.StatusCreated, CreateResponse{
			Response:  response.Response{Status: response.StatusOk},
			StudentId: studentId,
		})
	}
}

func (h *StudentHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student id format")))
			return
		}

		slog.Info("handling get student by id", slog.Int64("id", id))

		student, err := h.service.GetStudentById(id)
		if err != nil {
			slog.Warn("student not found", slog.Int64("id", id))
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

type StudentListResponse struct {
	response.Response
	Students []models.Student `json:"students"`
}

func (h *StudentHandler) GetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling get all students request")

		students, err := h.service.GetStudents()
		if err != nil {
			slog.Error("failed to retrieve students list", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(errors.New("failed to retrieve students")))
			return
		}

		if students == nil {
			students = []models.Student{}
		}

		response.WriteJson(w, http.StatusOK, StudentListResponse{
			Response: response.Response{Status: response.StatusOk},
			Students: students,
		})
	}
}

func (h *StudentHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student id format")))
			return
		}

		slog.Info("handling update student", slog.Int64("id", id))

		var student models.Student
		err = json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid request body payload")))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			var validateErrs validator.ValidationErrors
			if errors.As(err, &validateErrs) {
				response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = h.service.UpdateStudent(id, student.Name, student.Email, student.Age)
		if err != nil {
			slog.Warn("failed to update student", slog.Int64("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, response.Response{Status: response.StatusOk})
	}
}

func (h *StudentHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student id format")))
			return
		}

		slog.Info("handling delete student", slog.Int64("id", id))

		err = h.service.DeleteStudent(id)
		if err != nil {
			slog.Warn("failed to delete student", slog.Int64("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, response.Response{Status: response.StatusOk})
	}
}
