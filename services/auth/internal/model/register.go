package model

import "time"

type UserRegister struct {
	Email      string    `json:"email" validate:"required,email"`
	Password   string    `json:"password" validate:"required,min=8"`
	FirstName  string    `json:"first_name" validate:"required"`
	LastName   string    `json:"last_name" validate:"required"`
	Phone      string    `json:"phone"`
	Position   string    `json:"position"`
	Department string    `json:"department"`
	AvatarURL  string    `json:"avatar_url"`
	HireDate   time.Time `json:"hire_date"`
	BirthDate  time.Time `json:"birth_date"`
}
