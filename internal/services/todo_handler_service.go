package handlers

import "net/http"

type TodoHandlerService interface {
	CreateTodo(w http.ResponseWriter, r *http.Request)
	CreateList(w http.ResponseWriter, r *http.Request)
	CreateTag(w http.ResponseWriter, r *http.Request)
	UpdateTodo(w http.ResponseWriter, r *http.Request)
	DeleteTodo(w http.ResponseWriter, r *http.Request)
	GetTodo(w http.ResponseWriter, r *http.Request)
}
