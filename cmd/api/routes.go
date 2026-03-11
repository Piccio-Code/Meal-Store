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

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	router.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware)

		// Store
		r.Post("/v1/store", app.createStoreHandler)
		r.Get("/v1/store", app.listStoreHandler)
		r.Put("/v1/store", app.updateStoreHandler)
		r.Get("/v1/store-options", app.getStoreOptions)
		r.Get("/v1/store-id", app.getStoreId)

		// Store by ID
		r.Group(func(r chi.Router) {
			r.Use(app.RequireStoreId)

			r.Get("/v1/store/{store_id}", app.getStoreHandler)
			r.Delete("/v1/store/{store_id}", app.deleteStoreHandler)

			// Eaten items
			r.Post("/v1/store/{store_id}/eatenItem", app.createEatenHandler)
			r.Post("/v1/store/{store_id}/eatenItem-list", app.createEatenListHandler)

			// Items
			r.Post("/v1/store/{store_id}/items", app.createItemsHandler)
			r.Post("/v1/store/{store_id}/items-list", app.createItemsListHandler)
			r.Get("/v1/store/{store_id}/items", app.listItemsHandler)
			r.Get("/v1/store/{store_id}/items-options", app.getItemsOptionsHandler)
			r.Get("/v1/store/{store_id}/item-id", app.getItemsId)
			r.Put("/v1/store/{store_id}/items", app.updateItemsHandler)
			r.Put("/v1/store/{store_id}/items-list", app.updateItemsListHandler)

			// Item by ID
			r.Group(func(r chi.Router) {
				r.Use(app.RequireItemId)

				r.Get("/v1/store/{store_id}/items/{item_id}", app.getItemsHandler)
				r.Delete("/v1/store/{store_id}/items/{item_id}", app.deleteItemsHandler)

				r.Get("/v1/store/{store_id}/eatenItem/{item_id}", app.getEatenHandler)

			})
		})
	})

	return router
}
