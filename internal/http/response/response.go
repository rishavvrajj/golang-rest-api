package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Response represents the standard API response structure.
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

// WriteJson encodes data to JSON and writes it to ResponseWriter with the given status code.
func WriteJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"status":"Error","error":"internal server error"}`, http.StatusInternalServerError)
	}
}

// GeneralError creates a standard error response struct.
func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// ValidationError formats validation errors from go-playground/validator into user-friendly messages.
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is required", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be a valid email address", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be at least %s", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be at most %s", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' failed on tag '%s'", err.Field(), err.ActualTag()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
