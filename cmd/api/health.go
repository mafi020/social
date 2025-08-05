package main

import (
	"net/http"

	"github.com/mafi020/social/internal/utils"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "Ok",
		"env":    app.config.env,
	}
	if err := utils.WriteJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
