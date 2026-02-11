package httperror

import "net/http"

func InternalError(w http.ResponseWriter) {
	http.Error(w, "An internal error occurred", http.StatusInternalServerError)
}

func Forbidden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func CantCreateUser(w http.ResponseWriter) {
	http.Error(w, "Can't create user", http.StatusInternalServerError)
}
