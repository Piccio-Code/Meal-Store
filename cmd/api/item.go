package main

import (
	"fmt"
	. "github.com/Piccio-Code/MealStore/internal/data"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
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
		app.BadRequestError(w, r)
		return
	}

	onlyWarnings := false

	if val := r.URL.Query().Get("only_warnings"); val != "" {
		parsed, err := strconv.ParseBool(val)

		if err != nil {
			app.errorLog.Println(err)
			app.BadRequestError(w, r)
			return
		}

		onlyWarnings = parsed
	}

	items, err := app.models.Items.List(r.Context(), storeId, onlyWarnings)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"items": items})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) getItemsId(w http.ResponseWriter, r *http.Request) {
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
		app.BadRequestError(w, r)
		return
	}

	storeName := r.URL.Query().Get("store_name")

	itemId, err := app.models.Items.GetId(r.Context(), storeName, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"item_id": itemId})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) getItemsOptionsHandler(w http.ResponseWriter, r *http.Request) {
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
		app.BadRequestError(w, r)
		return
	}

	items, err := app.models.Items.List(r.Context(), storeId, false)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	var options []OptionStruct
	currentCap := make(map[string]int)
	versions := make(map[string]string)

	for _, item := range items {
		name := *item.Name
		capacity := *item.CurrentCapacity

		options = append(options, OptionStruct{fmt.Sprintf("%s (%d)", name, capacity)})
		currentCap[name] = capacity
		versions[name] = *item.Version
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"options": []envelop{{"names": options}, {"capacities": currentCap}, {"versions": versions}}})

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
		app.BadRequestError(w, r)
		return
	}

	newItem.StoreId = storeId

	err = app.models.Items.Insert(r.Context(), &newItem)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"new_item": newItem})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) createItemsListHandler(w http.ResponseWriter, r *http.Request) {
	var newItemsList struct {
		Items []*Item `json:"items"`
	}

	err := app.readeJSON(r, &newItemsList)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err = v.Struct(newItemsList)

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
		app.BadRequestError(w, r)
		return
	}

	err = app.models.Items.InsertList(r.Context(), newItemsList.Items, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.WriteError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelop{"new_items": newItemsList.Items})

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
		app.BadRequestError(w, r)
		return
	}

	item, err := app.models.Items.Get(r.Context(), itemId, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
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
		app.BadRequestError(w, r)
		return
	}

	oldItem, err := app.models.Items.Get(r.Context(), *newItemReq.Id, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
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
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, envelop{"updated_item": newItem})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}

func (app *application) updateItemsListHandler(w http.ResponseWriter, r *http.Request) {
	var newItemsReq struct {
		UpdateItemList []*UpdateItem `json:"items,omitempty"`
	}

	err := app.readeJSON(r, &newItemsReq)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
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
		app.BadRequestError(w, r)
		return
	}

	var newItems []Item

	for _, newItemReq := range newItemsReq.UpdateItemList {
		oldItem, err := app.models.Items.Get(r.Context(), *newItemReq.Id, storeId)

		if err != nil {
			app.errorLog.Println(err)
			app.BadRequestError(w, r)
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
			app.BadRequestError(w, r)
			return
		}

		newItems = append(newItems, newItem)
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"updated_items": newItems})

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
		app.BadRequestError(w, r)
		return
	}

	err = app.models.Items.Delete(r.Context(), itemId, storeId)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"deleted_item_id": itemId})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}
