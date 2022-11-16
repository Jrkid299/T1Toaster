// Filename: internal/data/toasts.go

package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"toaster.jalen.net/internals/validator"
)

type Toast struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
	Website   string    `json:"website,omitempty"`
	Address   string    `json:"address"`
	Mode      []string  `json:"mode"`
	Version   int32     `json:"version"`
}

func ValidateToast(v *validator.Validator, toast *Toast) {
	// Use the Check() method to execute our validation checks
	v.Check(toast.Name != "", "name", "must be provided")
	v.Check(len(toast.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(toast.Level != "", "level", "must be provided")
	v.Check(len(toast.Level) <= 200, "level", "must not be more than 200 bytes long")

	v.Check(toast.Contact != "", "contact", "must be provided")
	v.Check(len(toast.Contact) <= 200, "contact", "must not be more than 200 bytes long")

	v.Check(toast.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(toast.Phone, validator.PhoneRX), "phone", "must be a valid phone number")

	v.Check(toast.Email != "", "email", "must be provided")
	v.Check(validator.Matches(toast.Email, validator.EmailRX), "email", "must be a valid email address")

	v.Check(toast.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(toast.Website), "website", "must be a valid URL")

	v.Check(toast.Address != "", "address", "must be provided")
	v.Check(len(toast.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(toast.Mode != nil, "mode", "must be provided")
	v.Check(len(toast.Mode) >= 1, "mode", "must contain at least 1 entry")
	v.Check(len(toast.Mode) <= 5, "mode", "must contain at most 5 entries")
	v.Check(validator.Unique(toast.Mode), "mode", "must not contain duplicate entries")
}

// Define a ToastModel which wraps a sql.DB connection pool
type ToastModel struct {
	DB *sql.DB
}

// Insert() allows us  to create a new toast
func (m ToastModel) Insert(toast *Toast) error {
	query := `
		INSERT INTO toasts (name, level, contact, phone, email, website, address, mode)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, version
	`
	// Collect the data fields into a slice
	args := []interface{}{
		toast.Name, toast.Level,
		toast.Contact, toast.Phone,
		toast.Email, toast.Website,
		toast.Address, pq.Array(toast.Mode),
	}
	return m.DB.QueryRow(query, args...).Scan(&toast.ID, &toast.CreatedAt, &toast.Version)
}

// Get() allows us to retrieve a specific toast
func (m ToastModel) Get(id int64) (*Toast, error) {
	// Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query
	query := `
		SELECT id, created_at, name, level, contact, phone, email, website, address, mode, version
		FROM toasts
		WHERE id = $1
	`
	// Declare a Toast variable to hold the returned data
	var toast Toast
	// Execute the query using QueryRow()
	err := m.DB.QueryRow(query, id).Scan(
		&toast.ID,
		&toast.CreatedAt,
		&toast.Name,
		&toast.Level,
		&toast.Contact,
		&toast.Phone,
		&toast.Email,
		&toast.Website,
		&toast.Address,
		pq.Array(&toast.Mode),
		&toast.Version,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &toast, nil
}

// Update() allows us to edit/alter a specific toast
func (m ToastModel) Update(toast *Toast) error {
	return nil
}

// Delete() removes a specific toast
func (m ToastModel) Delete(id int64) error {
	return nil
}
