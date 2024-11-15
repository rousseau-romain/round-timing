package handlers

import (
	"github.com/rousseau-romain/round-timing/v0/service/auth"
	"github.com/rousseau-romain/round-timing/v0/shared/components"
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
