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
	fmt.Fprintln(w, "create a new toast..")
}

// showToastHandler for the "GET /v1/toasts/:id" endpoint
func (app *application) showToastHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Create a new instance of the School struct containing the ID we extracted
	// from our URL and some sample data
	toast := data.Toast{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Toast",
		Level:     "High School",
		Contact:   "Anna Smith",
		Phone:     "601-4411",
		Address:   "14 Apple street",
		Mode:      []string{"blended", "online"},
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, toast, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
