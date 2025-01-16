package bunapp

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
)

func (app *App) initRouter() {
	log.SetReportCaller(true)

	app.router = chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	var router = app.router.Group(func(r chi.Router) {
		r.Use(chimiddle.StripSlashes)
		r.Use(cors.Handler)
	})
	app.apiRouter = &router
}
