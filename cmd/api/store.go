package main

import (
	. "github.com/Piccio-Code/MealStore/internal/data"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (app *application) createStoreHandler(w http.ResponseWriter, r *http.Request) {
	var newStore Store

	err := app.readeJSON(r, &newStore)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err = v.Struct(newStore)

	if err != nil {
		app.ValidationError(w, r, err)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	err = app.models.Stores.Insert(r.Context(), &newStore, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	app.infoLog.Println(newStore)

	err = app.writeJSON(w, http.StatusCreated, envelop{"new_store": newStore})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) listStoreHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	stores, err := app.models.Stores.List(r.Context(), userId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"stores": stores})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) getStoreHandler(w http.ResponseWriter, r *http.Request) {
	storeId, ok := r.Context().Value(StoreIdKey).(int)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	store, err := app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"store": store})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) updateStoreHandler(w http.ResponseWriter, r *http.Request) {
	var newStore Store

	err := app.readeJSON(r, &newStore)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err = v.Struct(newStore)

	if err != nil {
		app.ValidationError(w, r, err)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	err = app.models.Stores.Update(r.Context(), &newStore, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"updated_store": newStore})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) deleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	storeId, ok := r.Context().Value(StoreIdKey).(int)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	err := app.models.Stores.Delete(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"deleted_store_id": storeId})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}
