// Filename: cms/api/toasts.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"toaster.jalen.net/internals/data"
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
	// Display the request
	fmt.Fprintf(w, "%+v\n", input)
}

// showToastHandler for the "GET /v1/toasts/:id" endpoint
func (app *application) showToastHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the School struct containing the ID we extracted
	// from our URL and some sample data
	toast := data.Toast{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Toast",
		Level:     "High School",
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
