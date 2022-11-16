// Filename: cms/api/toasts.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"toaster.jalen.net/internals/data"
	"toaster.jalen.net/internals/validator"
)

// createToastHandler for the "POST /v1/toasts" endpoint
func (app *application) createToastHandler(w http.ResponseWriter, r *http.Request) {
	// Our target decode destination
	var input struct {
		Name    string   `json:"name"`
		Level   string   `json:"level"`
		Contact string   `json:"contact"`
		Phone   string   `json:"phone"`
		Email   string   `json:"email"`
		Website string   `json:"website"`
		Address string   `json:"address"`
		Mode    []string `json:"mode"`
	}
	// Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new toast struct
	toast := &data.Toast{
		Name:    input.Name,
		Level:   input.Level,
		Contact: input.Contact,
		Phone:   input.Phone,
		Email:   input.Email,
		Website: input.Website,
		Address: input.Address,
		Mode:    input.Mode,
	}

	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateToast(v, toast); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Create a toast
	err = app.models.Toasts.Insert(toast)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// Create a Location header for the newly created resource/toast
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/toasts/%d", toast.ID))
	// Write the JSON response with 201 - Created status code with the body
	// being the toast data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"toast": toast}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showToastHandler for the "GET /v1/toasts/:id" endpoint
func (app *application) showToastHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the toast struct containing the ID we extracted
	// from our URL and some sample data
	toast := data.Toast{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Toast",
		Level:     "High toast",
		Contact:   "Jalen Lamb",
		Phone:     "615-7940",
		Address:   "20 Bahamas Street",
		Mode:      []string{"blended", "online"},
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"toast": toast}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
