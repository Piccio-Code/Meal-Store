package main

import (
	. "github.com/Piccio-Code/MealStore/internal/data"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (app *application) listItemsHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err := app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	items, err := app.models.Items.List(r.Context(), storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"items": items})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) createItemsHandler(w http.ResponseWriter, r *http.Request) {
	var newItem Item

	err := app.readeJSON(r, &newItem)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err = v.Struct(newItem)

	if err != nil {
		app.ValidationError(w, r, err)
		return
	}

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

	_, err = app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	newItem.StoreId = storeId

	err = app.models.Items.Insert(r.Context(), &newItem)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"new_item": newItem})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	storeId, ok := r.Context().Value(StoreIdKey).(int)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	itemId, ok := r.Context().Value(ItemIdKey).(int)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	_, err := app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	item, err := app.models.Items.Get(r.Context(), itemId, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"item": item})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) updateItemsHandler(w http.ResponseWriter, r *http.Request) {
	var newItemReq UpdateItem

	err := app.readeJSON(r, &newItemReq)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err = v.Struct(newItemReq)

	if err != nil {
		app.ValidationError(w, r, err)
		return
	}

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

	_, err = app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	oldItem, err := app.models.Items.Get(r.Context(), *newItemReq.Id, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	if newItemReq.Name == nil {
		newItemReq.Name = oldItem.Name
	}

	if newItemReq.CurrentCapacity == nil {
		newItemReq.CurrentCapacity = oldItem.CurrentCapacity
	}

	newItem := Item{
		Id:              newItemReq.Id,
		Name:            newItemReq.Name,
		StoreId:         storeId,
		CurrentCapacity: newItemReq.CurrentCapacity,
		Version:         newItemReq.Version,
	}

	err = app.models.Items.Update(r.Context(), &newItem)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"updated_item": newItem})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) deleteItemsHandler(w http.ResponseWriter, r *http.Request) {
	storeId, ok := r.Context().Value(StoreIdKey).(int)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	itemId, ok := r.Context().Value(ItemIdKey).(int)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	userId, ok := r.Context().Value(CurrentUserIDKey).(string)

	if !ok {
		app.UnauthorizedError(w, r)
		return
	}

	_, err := app.models.Stores.Get(r.Context(), storeId, userId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	err = app.models.Items.Delete(r.Context(), itemId, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"deleted_item_id": itemId})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}
