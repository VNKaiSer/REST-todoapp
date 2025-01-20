package httpresponse

import (
	"net/http"

	"github.com/go-chi/render"
)

type SingleResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}

type CollectionResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Total   int         `json:"total"`
}

func (resp *SingleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, resp.Status)
	return nil
}

func (resp *CollectionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, resp.Status)
	return nil
}


// Write a single response
func WriteResponse(w http.ResponseWriter, status int, message string, data interface{}) render.Renderer {
	return &SingleResponse{
		Message: message,
		Data:    data,
		Status:  status,
	}
}

// Write a collection response
func WriteCollectionResponse(w http.ResponseWriter, status int, message string, data interface{}, total int) render.Renderer {
	return &CollectionResponse{
		Message: message,
		Data:    data,
		Status:  status,
		Total:   total,
	}
}
