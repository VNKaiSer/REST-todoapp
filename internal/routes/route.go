package routes

import (
	"context"
	"os"
	"todo-app/bunapp"
	"todo-app/internal/handlers"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func SetupRoutes() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	log.Info("Setting up routes")
	bunapp.OnStart("example.init", func(ctx context.Context, app *bunapp.App) error {
		router := app.Router()
		serverHandler := handlers.NewServerHandler(app)
		router.Route("/api", func(r chi.Router) {
			r.Get("/ping", serverHandler.ReplayAppCheck)
			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", handlers.NewAuthHandler(app).Login)
				r.Post("/register", handlers.NewAuthHandler(app).Register)
				r.Post("/refresh-token", handlers.NewAuthHandler(app).RefreshToken)
				r.With(handlers.NewAuthHandler(app).Authorization()).Get("/check-token", handlers.NewAuthHandler(app).CheckToken)
			})
		})
		return nil
	})
}
