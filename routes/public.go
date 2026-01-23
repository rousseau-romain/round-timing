package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func registerPublicRoutes(r *mux.Router) {
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}
