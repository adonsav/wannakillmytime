package repository

import (
	"context"
	"errors"
	"github.com/adonsav/fgoapp/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (pr *postgresDBRepo) AllUsers() (*[]models.User, error) {
	return &[]models.User{}, nil
}

// InsertRegistration inserts a registration into the database
func (pr *postgresDBRepo) InsertRegistration(reg models.Registration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into registrations (user_name, email, password, active, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6)`

	_, err := pr.DB.ExecContext(ctx, stmt,
		reg.UserName,
		reg.Email,
		reg.Password,
		true,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

// GetUserByID returns a user by id
func (pr *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
			from users where id = $1`

	row := pr.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

// UpdateUser updates a user in the database
func (pr *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
			from users where id = user.id`

	_, err := pr.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (pr *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := pr.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllRegistrations returns a slice of all registrations
func (pr *postgresDBRepo) AllRegistrations() ([]models.Registration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var registrations []models.Registration

	query := `select r.id, r.user_name, r.email, r.active, r.created_at, r.updated_at
		from registrations r
		order by created_at asc`

	rows, err := pr.DB.QueryContext(ctx, query)
	if err != nil {
		return registrations, err
	}
	defer rows.Close()

	for rows.Next() {
		var reg models.Registration
		err = rows.Scan(
			&reg.ID,
			&reg.UserName,
			&reg.Email,
			&reg.Active,
			&reg.CreatedAt,
			&reg.UpdatedAt,
		)

		if err != nil {
			return registrations, err
		}
		registrations = append(registrations, reg)
	}

	err = rows.Err()
	if err != nil {
		return registrations, err
	}

	return registrations, nil
}

// ActiveRegistrations returns a slice of only active registrations
func (pr *postgresDBRepo) ActiveRegistrations() ([]models.Registration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var registrations []models.Registration

	query := `select r.id, r.user_name, r.email, r.active, r.created_at, r.updated_at
		from registrations r
		where r.active = true
		order by created_at asc`

	rows, err := pr.DB.QueryContext(ctx, query)
	if err != nil {
		return registrations, err
	}
	defer rows.Close()

	for rows.Next() {
		var reg models.Registration
		err = rows.Scan(
			&reg.ID,
			&reg.UserName,
			&reg.Email,
			&reg.Active,
			&reg.CreatedAt,
			&reg.UpdatedAt,
		)

		if err != nil {
			return registrations, err
		}
		registrations = append(registrations, reg)
	}

	err = rows.Err()
	if err != nil {
		return registrations, err
	}

	return registrations, nil
}

// GetRegistrationByID returns one registration by ID
func (pr *postgresDBRepo) GetRegistrationByID(id int) (models.Registration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var registration models.Registration

	query := `select r.id, r.user_name, r.email, r.active, r.created_at, r.updated_at
		from registrations r
		where r.id = $1`

	row := pr.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&registration.ID,
		&registration.UserName,
		&registration.Email,
		&registration.Active,
		&registration.CreatedAt,
		&registration.UpdatedAt,
	)

	if err != nil {
		return registration, err
	}
	return registration, nil
}

// UpdateRegistration updates a registration in the database
func (pr *postgresDBRepo) UpdateRegistration(registration models.Registration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update registrations set user_name = $1, email = $2, updated_at = $3 where id = $4`

	_, err := pr.DB.ExecContext(ctx, query,
		registration.UserName,
		registration.Email,
		time.Now(),
		registration.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeactivateRegistration marks a registration as inactive in the database
func (pr *postgresDBRepo) DeactivateRegistration(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update registrations set active = $1, updated_at = $2 where id = $3`

	_, err := pr.DB.ExecContext(ctx, query,
		false,
		time.Now(),
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
