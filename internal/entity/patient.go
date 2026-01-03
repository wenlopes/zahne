package entity

import "time"

type Patient struct {
	ID          string
	Name        string
	CPF         string
	Phone       string
	Email       string
	DateOfBirth time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
