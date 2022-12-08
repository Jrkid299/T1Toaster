// Filename: internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// A wrapper for our data models
type Models struct {
	Permissions PermissionModel
	Toasts      ToastModel
	Tokens      TokenModel
	Users       UserModel
}

// NewModels() allows us to create a new Models
func NewModels(db *sql.DB) Models {
	return Models{
		Permissions: PermissionModel{DB: db},
		Toasts:      ToastModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
