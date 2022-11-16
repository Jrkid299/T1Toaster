// Filename: internal/data/toasts.go

package data

import (
	"time"

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
