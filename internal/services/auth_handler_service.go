package handlers

import "net/http"

type AuthHandlerService interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	CheckToken(w http.ResponseWriter, r *http.Request)
}
