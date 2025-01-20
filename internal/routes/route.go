package routes

import (
	"context"
	"os"
	"todo-app/bunapp"
	"todo-app/internal/handlers"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	log.Info("Setting up routes")
	bunapp.OnStart("example.init", func(ctx context.Context, app *bunapp.App) error {
		router := app.Router()
		serverHandler := handlers.NewServerHandler(app)
		authHandler := handlers.NewAuthHandler(app)
		todoHandler := handlers.NewTodoHandler(app)
		router.Get("/docs/*", httpSwagger.WrapHandler)
		router.Route("/api", func(r chi.Router) {
			r.Get("/ping", serverHandler.ReplayAppCheck)
			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", authHandler.Login)
				r.Post("/register", authHandler.Register)
				r.Post("/refresh-token", authHandler.RefreshToken)
				r.With(authHandler.Authorization).Get("/check-token", authHandler.CheckToken)
			})

			r.Route("/todo", func(r chi.Router) {
				r.Use(authHandler.Authorization)
				r.Post("/", todoHandler.CreateTodo)
				r.Get("/{id}", todoHandler.GetTodo)
				r.Put("/{id}", todoHandler.UpdateTodo)
				r.Delete("/{id}", todoHandler.DeleteTodo)
			})

		})
		return nil
	})
}
