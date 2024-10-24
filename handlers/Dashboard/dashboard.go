// This file contains the handler for the dashboard endpoint
package dashboard

import (
	"fmt"
	"net/http"

	"erp/middleware"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get the userID from the context
	email, err := middleware.GetUserEmailFromContext(r.Context())
	if err != nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// The email is authenticated at this point, show the email
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to your dashboard! Your user ID is: %s", email)
}
