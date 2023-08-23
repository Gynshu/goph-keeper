package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gynshu-one/goph-keeper/server/internal/middlewares"
	"github.com/gynshu-one/goph-keeper/server/internal/server"
)

func NewRouter(handlers server.Handlers) *chi.Mux {
	// New Chi router
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/user", func(r chi.Router) {
		r.With().Get("/create", handlers.CreateUser)
		r.With().Get("/login", handlers.LoginUser)
		r.With(middlewares.SessionCheck).Get("/logout", handlers.LogoutUser)
		r.With(middlewares.SessionCheck).Post("/set", handlers.SetUserData)
		r.With(middlewares.SessionCheck).Get("/get", handlers.GetUserData)
		r.With(middlewares.SessionCheck).Get("/delete", handlers.DeleteUserData)
		r.With(middlewares.SessionCheck).Get("/sync", handlers.SyncUserData)
	})

	return r
}
