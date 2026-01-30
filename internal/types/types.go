package types

type Student struct {
	ID    int    `json:"id" redis:"id"`
	Name  string `json:"name" validate:"required" redis:"name"`
	Email string `json:"email" validate:"required,email" redis:"email"`
	Age   int    `json:"age" validate:"required" redis:"age"`
}
