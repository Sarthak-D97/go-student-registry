package entity

type Student struct {
	ID    int    `json:"id" redis:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" validate:"required" redis:"name"`
	Email string `json:"email" validate:"required,email" redis:"email" gorm:"uniqueIndex"`
	Age   int    `json:"age" validate:"required" redis:"age"`
}
