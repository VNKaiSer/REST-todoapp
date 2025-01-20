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

// @Summary Ping the server
// @Description Kiểm tra server có hoạt động không
// @Tags Ping
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Router /api/ping [get]
func (h *ServerHandler) ReplayAppCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
