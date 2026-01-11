package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) getIdParam(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil || id < 1 {
		return 0, errors.New("error getting the id")
	}

	return id, nil
}

type envelop map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelop) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readeJSON(r *http.Request, dst interface{}) error {

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(&dst)
}
