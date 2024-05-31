package handlers

import "round-timing/service/auth"

type Handler struct {
	auth *auth.AuthService
}

func New(auth *auth.AuthService) *Handler {
	return &Handler{
		auth: auth,
	}
}
