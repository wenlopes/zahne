package request

import "time"

type CreatePatientRequest struct {
	Name        string    `json:"name" binding:"required"`
	CPF         string    `json:"cpf" binding:"required"`
	Phone       string    `json:"phone" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}

type UpdatePatientRequest struct {
	Name        string    `json:"name" binding:"required"`
	CPF         string    `json:"cpf" binding:"required"`
	Phone       string    `json:"phone" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}
