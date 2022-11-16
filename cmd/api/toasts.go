// Filename: cms/api/toasts.go

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// createToastHandler for the "POST /v1/toasts" endpoint
func (app *application) createToastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new toast..")
}

// showToastHandler for the "GET /v1/toasts/:id" endpoint
func (app *application) showToastHandler(w http.ResponseWriter, r *http.Request) {
	// Use the "ParamsFromContext()" function to get the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	// Get the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Display the toasts id
	fmt.Fprintf(w, "show the details for toast %d\n", id)
}
