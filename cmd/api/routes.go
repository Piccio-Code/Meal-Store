package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Group(func(r chi.Router) {
			r.Use(app.AuthMiddleware)

			app.storeRoutes(r)
		})
	})

	return router
}

func (app *application) storeRoutes(r chi.Router) {
	r.Post("/store", app.createStoreHandler)
	r.Get("/store", app.listStoreHandler)

	r.Put("/store", app.updateStoreHandler)

	r.Get("/store-options", app.getStoreOptions)
	r.Get("/store-id", app.getStoreId)

	r.Route("/store/{store_id}", func(r chi.Router) {
		r.Use(app.RequireStoreId)

		r.Get("/", app.getStoreHandler)
		r.Delete("/", app.deleteStoreHandler)

		app.itemRoutes(r)
	})
}

func (app *application) itemRoutes(r chi.Router) {
	r.Post("/items", app.createItemsHandler)
	r.Post("/items-list", app.createItemsListHandler)
	r.Get("/items", app.listItemsHandler)

	r.Put("/items", app.updateItemsHandler)
	r.Put("/items-list", app.updateItemsListHandler)

	r.Get("/items-options", app.getItemsOptionsHandler)
	r.Get("/item-id", app.getItemsId)

	r.Route("/items/{item_id}", func(r chi.Router) {
		r.Use(app.RequireItemId)

		r.Get("/", app.getItemsHandler)
		r.Delete("/", app.deleteItemsHandler)
	})
}
