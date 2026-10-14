package storage

import (
	"github.com/rishavvrajj/golang-rest-api/internal/models"
)

// Storage defines the interface for student data persistence operations.
type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (models.Student, error)
	GetStudents() ([]models.Student, error)
	UpdateStudent(id int64, name string, email string, age int) error
	DeleteStudent(id int64) error
}
