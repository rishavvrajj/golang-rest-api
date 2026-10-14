package models

// Student represents the student domain entity.
type Student struct {
	ID    int64  `json:"id"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=1,max=120"`
}
