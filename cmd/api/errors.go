package main

import (
	"log"
	"net/http"

	"github.com/mafi020/social/internal/utils"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal Server error. method: %s path: %s errors: %s", r.Method, r.URL.Path, err)
	utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error. Please try again")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error. method: %s path: %s errors: %s", r.Method, r.URL.Path, err)
	utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationError(w http.ResponseWriter, r *http.Request, err map[string]string) {
	log.Printf("invalid data. method: %s path: %s errors: %s", r.Method, r.URL.Path, err)
	utils.WriteJSONError(w, http.StatusUnprocessableEntity, err)
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error. method: %s path: %s errors: %s", r.Method, r.URL.Path, err)
	utils.WriteJSONError(w, http.StatusNotFound, "resource not found")
}
