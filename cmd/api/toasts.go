// Filename: cms/api/toasts.go

package main

import (
	"errors"
	"fmt"
	"net/http"

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

	// Fetch the specific toast
	toast, err := app.models.Toasts.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"toast": toast}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateToastHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a complete replacement
	// Get the id for the toast that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the orginal record from the database
	toast, err := app.models.Toasts.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Create an input struct to hold data read in fro mteh client
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
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy / Update the fields / values in the toast variable using the fields
	// in the input struct
	toast.Name = input.Name
	toast.Level = input.Level
	toast.Contact = input.Contact
	toast.Phone = input.Phone
	toast.Email = input.Email
	toast.Website = input.Website
	toast.Address = input.Address
	toast.Mode = input.Mode
	// Perform validation on the updated Toast. If validation fails, then
	// we send a 422 - Unprocessable Entity respose to the client
	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateToast(v, toast); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated Toast record to the Update() method
	err = app.models.Toasts.Update(toast)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"toast": toast}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
