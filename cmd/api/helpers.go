package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) getIdParam(r *http.Request, name string) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, name))

	if err != nil || id < 1 {
		return 0, fmt.Errorf("error getting the id: {%s}", name)
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
