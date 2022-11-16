// Filename: cms/api/toasts.go

package main

import (
	"fmt"
	"net/http"
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
	// Display the toasts id
	fmt.Fprintf(w, "show the details for toast %d\n", id)
}
