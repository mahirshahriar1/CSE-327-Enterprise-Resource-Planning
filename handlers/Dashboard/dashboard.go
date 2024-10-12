package dashboard

import (
	"fmt"
	"net/http"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	// The user is authenticated at this point
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to your dashboard!")
}
