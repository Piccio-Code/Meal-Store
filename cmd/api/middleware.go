package main

import (
	"context"
	"net/http"
)

type CurrentUserID string
type StoreId string
type ItemId string

const CurrentUserIDKey = CurrentUserID("CurrentUserIDKey")
const StoreIdKey = StoreId("StoreIdKey")
const ItemIdKey = ItemId("ItemIdKey")

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Telegram-Chat-ID")

			if id != "7202833466" {
				app.errorLog.Println("error HERE: ", id)
				app.UnauthorizedError(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), CurrentUserIDKey, id)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
}

func (app *application) RequireStoreId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		storeId, err := app.getIdParam(r, "store_id")

		if err != nil {
			app.errorLog.Println(err)
			app.InternalServerError(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), StoreIdKey, storeId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) RequireItemId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		itemId, err := app.getIdParam(r, "item_id")

		if err != nil {
			app.errorLog.Println(err)
			app.InternalServerError(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), ItemIdKey, itemId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
