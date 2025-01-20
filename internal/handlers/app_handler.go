package handlers

import (
	"net/http"
	"todo-app/bunapp"
)
type ServerHandler struct {
	app *bunapp.App
}

func NewServerHandler(app *bunapp.App) *ServerHandler {
	return &ServerHandler{
		app: app,
	}
}

func (h *ServerHandler) ReplayAppCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}