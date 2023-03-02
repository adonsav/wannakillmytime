package models

import (
	"time"
)

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Registration is the registration model
type Registration struct {
	ID        int
	UserName  string
	Email     string
	Password  string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Buddie    BoredBuddie
}

// BoredBuddie is the bored buddies model
type BoredBuddie struct {
	ID             int
	UserName       string
	Email          string
	RegistrationID int
}

// EmailData holds an email data
type EmailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
