package main

import (
	"github.com/Piccio-Code/MealStore/internal/data"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (app *application) createEatenHandler(w http.ResponseWriter, r *http.Request) {
	var newEatenItem data.EatenItem

	err := app.readeJSON(r, &newEatenItem)

	if err != nil {
		app.BadRequestError(w, r)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err = v.Struct(newEatenItem)

	if err != nil {
		app.ValidationError(w, r, err)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	storeId, ok := r.Context().Value(StoreIdKey).(int)

	if !ok {
		app.BadRequestError(w, r)
		return
	}

	_, err = app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.BadRequestError(w, r)
		return
	}

	_, err = app.models.Items.Get(r.Context(), newEatenItem.ItemId, storeId)

	if err != nil {
		app.BadRequestError(w, r)
		return
	}

	err = app.models.EatenItems.Create(r.Context(), newEatenItem)

	if err != nil {
		app.BadRequestError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) createEatenListHandler(w http.ResponseWriter, r *http.Request) {
	var newEatenItemList struct {
		NewEatenItems []*data.EatenItem `json:"new_eaten_items"`
	}

	err := app.readeJSON(r, &newEatenItemList)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err = v.Struct(newEatenItemList)

	if err != nil {
		app.errorLog.Println(err)
		app.ValidationError(w, r, err)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.errorLog.Println("Parse string: userId")
		app.UnauthorizedError(w, r)
		return
	}

	storeId, ok := r.Context().Value(StoreIdKey).(int)

	if !ok {
		app.errorLog.Println("Parse int: storeId")
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	_, err = app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	for _, item := range newEatenItemList.NewEatenItems {
		_, err := app.models.Items.Get(r.Context(), item.ItemId, storeId)

		if err != nil {
			app.errorLog.Println(err)
			app.BadRequestError(w, r)
			return
		}

	}

	err = app.models.EatenItems.CreateList(r.Context(), newEatenItemList.NewEatenItems)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getEatenHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	storeId, ok := r.Context().Value(StoreIdKey).(int)

	if !ok {
		app.BadRequestError(w, r)
		return
	}

	itemId, ok := r.Context().Value(ItemIdKey).(int)

	if !ok {
		app.BadRequestError(w, r)
		return
	}

	_, err := app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	_, err = app.models.Items.Get(r.Context(), itemId, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	filters, err := data.NewEatenItemFilters(r.URL.Query())

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	app.infoLog.Println(filters.AfterDate)

	items, err := app.models.EatenItems.Get(r.Context(), itemId, filters)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	app.writeJSON(w, http.StatusOK, envelop{"eaten_items": items})
}
