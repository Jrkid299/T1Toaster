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
	// This method does a partial replacement
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
	// Create an input struct to hold data read in from the client
	// We update input struct to use pointers because pointers have a
	// default value of nil
	// If a field remains nil then we know that the client did not update it
	var input struct {
		Name    *string  `json:"name"`
		Level   *string  `json:"level"`
		Contact *string  `json:"contact"`
		Phone   *string  `json:"phone"`
		Email   *string  `json:"email"`
		Website *string  `json:"website"`
		Address *string  `json:"address"`
		Mode    []string `json:"mode"`
	}

	// Initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check for updates
	if input.Name != nil {
		toast.Name = *input.Name
	}
	if input.Level != nil {
		toast.Level = *input.Level
	}
	if input.Contact != nil {
		toast.Contact = *input.Contact
	}
	if input.Phone != nil {
		toast.Phone = *input.Phone
	}
	if input.Email != nil {
		toast.Email = *input.Email
	}
	if input.Website != nil {
		toast.Website = *input.Website
	}
	if input.Address != nil {
		toast.Address = *input.Address
	}
	if input.Mode != nil {
		toast.Mode = input.Mode
	}
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
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
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

func (app *application) deleteToastHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id for the toast that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the Toast from the database. Send a 404 Not Found status code to the
	// client if there is no matching record
	err = app.models.Toasts.Delete(id)
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
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "toast successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// The listToastsHandler() allows the client to see a listing of toasts
// based on a set of criteria
func (app *application) listToastsHandler(w http.ResponseWriter, r *http.Request) {
	// Create an input struct to hold our query parameters
	var input struct {
		Name  string
		Level string
		Mode  []string
		data.Filters
	}
	// Initialize a validator
	v := validator.New()
	// Get the URL values map
	qs := r.URL.Query()
	// Use the helper methods to extract the values
	input.Name = app.readString(qs, "name", "")
	input.Level = app.readString(qs, "level", "")
	input.Mode = app.readCSV(qs, "mode", []string{})
	// Get the page information
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Get the sort information
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Specific the allowed sort values
	input.Filters.SortList = []string{"id", "name", "level", "-id", "-name", "-level"}
	// Check for validation errors
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Get a listing of all toasts
	toasts, metadata, err := app.models.Toasts.GetAll(input.Name, input.Level, input.Mode, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containg all the toasts
	err = app.writeJSON(w, http.StatusOK, envelope{"toasts": toasts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
