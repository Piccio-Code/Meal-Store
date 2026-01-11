package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, envelop{"health_status": data})

	if err != nil {
		app.errorLog.Println("error encoding the json")
		app.BadRequestError(w, r)
		return
	}

	return
}
