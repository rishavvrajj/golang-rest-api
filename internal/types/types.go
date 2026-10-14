package types

type Student struct {
	id    int64
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int `validate:"required"`
}
