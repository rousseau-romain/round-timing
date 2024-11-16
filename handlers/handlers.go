package handlers

import (
	"round-timing/service/auth"
	"round-timing/shared/components"
)

type Handler struct {
	auth  *auth.AuthService
	error components.Error
}

func New(auth *auth.AuthService) *Handler {

	return &Handler{
		auth: auth,
		error: components.Error{
			Title:    "",
			Messages: []string{},
		},
	}
}
