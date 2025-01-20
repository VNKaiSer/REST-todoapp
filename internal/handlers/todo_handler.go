package handlers

import (
	"net/http"
	"todo-app/bunapp"
	handlers "todo-app/internal/services"
)

type TodoHandler struct {
	app *bunapp.App
}

// CreateList implements handlers.TodoHandlerService.
func (t *TodoHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// CreateTag implements handlers.TodoHandlerService.
func (t *TodoHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// CreateTodo implements handlers.TodoHandlerService.
func (t *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
}

// DeleteTodo implements handlers.TodoHandlerService.
func (t *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// GetTodo implements handlers.TodoHandlerService.
func (t *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// UpdateTodo implements handlers.TodoHandlerService.
func (t *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

var _ handlers.TodoHandlerService = (*TodoHandler)(nil)

func NewTodoHandler(app *bunapp.App) *TodoHandler {
	return &TodoHandler{app: app}
}
