package handlers

import (
	"errors"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	WrapErrorWithStatus(w, errors.New("Not found"), http.StatusNotFound)
}
