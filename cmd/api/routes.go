package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"time"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	v1Router := chi.NewRouter()

	v1Router.Group(func(public chi.Router) {
		public.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	})

	v1Router.Group(func(protected chi.Router) {
		protected.Use(app.AuthMiddleware)

		protected.HandleFunc("POST /store", app.createStoreHandler)

		protected.HandleFunc("GET /store", app.listStoreHandler)
		protected.HandleFunc("GET /store/{id}", app.getStoreHandler)

		protected.HandleFunc("PUT /store/{id}", app.updateStoreHandler)

		protected.HandleFunc("DELETE /store/{id}", app.deleteStoreHandler)
	})

	router.Mount("/v1", v1Router)

	return router
}
