package repository

import (
	"errors"
	"github.com/adonsav/fgoapp/internal/models"
)

func (pr *testDBRepo) AllUsers() (*[]models.User, error) {
	return &[]models.User{}, nil
}

// InsertRegistration inserts a registration into the database
func (pr *testDBRepo) InsertRegistration(reg models.Registration) error {
	if reg.UserName == "failRegistration" {
		return errors.New("insert registration to database failed")
	}
	return nil
}

// GetUserByID returns a user by id
func (pr *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

// UpdateUser updates a user in the database
func (pr *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

// Authenticate authenticates a user
func (pr *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "nonExistingUser@test.com" {
		return 0, "", errors.New("invalid email during authentication")
	}
	return 1, "", nil
}

// AllRegistrations returns a slice of all registrations
func (pr *testDBRepo) AllRegistrations() ([]models.Registration, error) {
	var registrations []models.Registration
	return registrations, nil
}

// ActiveRegistrations returns a slice of only active registrations
func (pr *testDBRepo) ActiveRegistrations() ([]models.Registration, error) {
	var registrations []models.Registration
	return registrations, nil
}

// GetRegistrationByID returns one registration by ID
func (pr *testDBRepo) GetRegistrationByID(id int) (models.Registration, error) {
	var registration models.Registration
	return registration, nil
}

// UpdateRegistration updates a registration in the database
func (pr *testDBRepo) UpdateRegistration(registration models.Registration) error {
	return nil
}

// DeactivateRegistration marks a registration as inactive in the database
func (pr *testDBRepo) DeactivateRegistration(id int) error {
	return nil
}
