package service

import (
	"github.com/rishavvrajj/golang-rest-api/internal/models"
	"github.com/rishavvrajj/golang-rest-api/internal/storage"
)

type StudentService struct {
	storage storage.Storage
}

// NewStudentService creates a new instance of StudentService.
func NewStudentService(s storage.Storage) *StudentService {
	return &StudentService{storage: s}
}

func (s *StudentService) CreateStudent(name string, email string, age int) (int64, error) {
	return s.storage.CreateStudent(name, email, age)
}

func (s *StudentService) GetStudentById(id int64) (models.Student, error) {
	return s.storage.GetStudentById(id)
}

func (s *StudentService) GetStudents() ([]models.Student, error) {
	return s.storage.GetStudents()
}

func (s *StudentService) UpdateStudent(id int64, name string, email string, age int) error {
	return s.storage.UpdateStudent(id, name, email, age)
}

func (s *StudentService) DeleteStudent(id int64) error {
	return s.storage.DeleteStudent(id)
}
