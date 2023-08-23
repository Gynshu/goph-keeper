package router

import (
	"github.com/gynshu-one/goph-keeper/server/api/handlers"
	"github.com/gynshu-one/goph-keeper/server/api/middlewares"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handlers handlers.Handlers) *chi.Mux {
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
