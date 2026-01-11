package main

import (
	"fmt"
	. "github.com/Piccio-Code/MealStore/internal/data"
	"net/http"
	"time"
)

func (app *application) createStore(w http.ResponseWriter, r *http.Request) {
	var newStore StoreInput

	err := app.readeJSON(r, &newStore)

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}

	fmt.Fprintln(w, newStore)
}

func (app *application) getStores(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Getting all stores...")
}

func (app *application) getStore(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdParam(r)

	if err != nil {
		app.errorLog.Println(err)
		app.NotFoundError(w, r)
		return
	}

	store := Store{ID: id, Name: "Piccio Home", CreatedAt: time.Now()}

	err = app.writeJSON(w, http.StatusOK, envelop{"store": store})

	if err != nil {
		app.errorLog.Println(err)
		app.BadRequestError(w, r)
		return
	}
}
