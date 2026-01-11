package main

import "net/http"

func (app *application) WriteError(w http.ResponseWriter, r *http.Request, status int, message string) {
	err := app.writeJSON(w, status, envelop{"error": message})
	if err != nil {
		app.errorLog.Println(err)
		app.InternalServerError(w, r)
	}
}

func (app *application) NotFoundError(w http.ResponseWriter, r *http.Request) {
	app.WriteError(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (app *application) InternalServerError(w http.ResponseWriter, r *http.Request) {
	app.WriteError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) MethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	app.WriteError(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}

func (app *application) BadRequestError(w http.ResponseWriter, r *http.Request) {
	app.WriteError(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}
