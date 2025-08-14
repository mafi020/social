package main

import (
	"net/http"

	"github.com/mafi020/social/internal/utils"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal Server error", "method", r.Method, "path", r.URL.Path, "errors", err)
	utils.JSONErrorResponse(w, http.StatusInternalServerError, "internal server error. Please try again")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "errors", err)
	utils.JSONErrorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationError(w http.ResponseWriter, r *http.Request, err map[string]string) {
	app.logger.Warnw("invalid data error", "method", r.Method, "path", r.URL.Path, "errors", err)
	utils.JSONErrorResponse(w, http.StatusUnprocessableEntity, err)
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "errors", err)
	utils.JSONErrorResponse(w, http.StatusNotFound, "resource not found")
}

func (app *application) unAuthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unAuthorized", "method", r.Method, "path", r.URL.Path, "errors", err)
	utils.JSONErrorResponse(w, http.StatusUnauthorized, map[string]string{"message": err.Error()})
}
