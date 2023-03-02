package repository

import "github.com/adonsav/fgoapp/internal/models"

type DatabaseRepo interface {
	AllUsers() (*[]models.User, error)
	InsertRegistration(reg models.Registration) error
	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllRegistrations() ([]models.Registration, error)
	ActiveRegistrations() ([]models.Registration, error)
	GetRegistrationByID(id int) (models.Registration, error)
	UpdateRegistration(registration models.Registration) error
	DeactivateRegistration(id int) error
}
