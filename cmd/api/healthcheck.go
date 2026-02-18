package main

import (
	"net/http"
)

type healthData struct {
	Status      string `json:"status,omitempty"`
	Environment string `json:"environment,omitempty"`
	Version     string `json:"version,omitempty"`
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := healthData{
		Status:      "available",
		Environment: app.config.env,
		Version:     version,
	}

	err := app.writeJSON(w, http.StatusOK, envelop{"health_status": data})

	if err != nil {
		app.errorLog.Println("error encoding the json")
		app.BadRequestError(w, r)
		return
	}

	return
}
