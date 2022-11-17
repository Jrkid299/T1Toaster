// Filename: internal/data/toasts.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&toast.ID, &toast.CreatedAt, &toast.Version)
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
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
// Optimistic locking (version number)
func (m ToastModel) Update(toast *Toast) error {
	// Create a query
	query := `
		UPDATE toasts
		SET name = $1, level = $2, contact = $3,
		    phone = $4, email = $5, website = $6,
			address = $7, mode = $8, version = version + 1
		WHERE id = $9
		AND version = $10
		RETURNING version
	`
	args := []interface{}{
		toast.Name,
		toast.Level,
		toast.Contact,
		toast.Phone,
		toast.Email,
		toast.Website,
		toast.Address,
		pq.Array(toast.Mode),
		toast.ID,
		toast.Version,
	}
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&toast.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete() removes a specific toast
func (m ToastModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
			DELETE FROM toasts
			WHERE id = $1
		`
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// The GetAll() method retuns a list of all the toasts sorted by id
func (m ToastModel) GetAll(name string, level string, mode []string, filters Filters) ([]*Toast, Metadata, error) {
	// Construct the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, name, level,
		       contact, phone, email, website,
			   address, mode, version
		FROM toasts
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', level) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (mode @> $3 OR $3 = '{}' )
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())

	// Create a 3-second-timout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query
	args := []interface{}{name, level, pq.Array(mode), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Close the resultset
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the Toast data
	toasts := []*Toast{}
	// Iterate over the rows in the resultset
	for rows.Next() {
		var toast Toast
		// Scan the values from the row into toast
		err := rows.Scan(
			&totalRecords,
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
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Toast to our slice
		toasts = append(toasts, &toast)
	}
	// Check for errors after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	// Return the slice of Toasts
	return nil, Metadata{}, err
}
