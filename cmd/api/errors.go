package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

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

func (app *application) UnauthorizedError(w http.ResponseWriter, r *http.Request) {
	app.WriteError(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

func (app *application) ValidationError(w http.ResponseWriter, r *http.Request, err error) {

	var validateErrs validator.ValidationErrors

	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			message := fmt.Sprintf("the field %v must met: tag: %v, value: %v", e.StructField(), e.ActualTag(), e.Param())

			app.errorLog.Println(message)
			app.WriteError(w, r, http.StatusBadRequest, message)
		}
	}
}
