package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WrapError(w http.ResponseWriter, err error) {
	WrapErrorWithStatus(w, err, http.StatusBadRequest)
}

func WrapErrorWithStatus(w http.ResponseWriter, err error, httpstatus int) {
	var m = map[string]string{
		"result": "error",
		"data":   err.Error(),
	}

	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpstatus)
	fmt.Fprintln(w, string(res))
}

func WrapOK(w http.ResponseWriter, m map[string]interface{}) {
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(res))
}
