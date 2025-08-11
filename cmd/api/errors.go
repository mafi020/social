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

/*
{"level":"warn","ts":1754876023.6537068,"caller":"api/errors.go:20","msg":"invalid data error!(EXTRA method=POST, path=/api/posts, errors=map[content:Content is required title:Title is required])"}
*/
