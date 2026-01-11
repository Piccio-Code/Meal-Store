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
	router.NotFound(app.NotFoundError)
	router.MethodNotAllowed(app.MethodNotAllowedError)

	v1Router := chi.NewRouter()

	v1Router.HandleFunc("GET /healthcheck", app.healthcheckHandler)

	v1Router.HandleFunc("POST /store", app.createStore)

	v1Router.HandleFunc("GET /store", app.getStores)
	v1Router.HandleFunc("GET /store/{id}", app.getStore)

	router.Mount("/v1", v1Router)

	return router
}
