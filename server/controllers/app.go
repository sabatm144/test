package controller

import (
	"encoding/json"
	"net/http"
)

func renderJSON(w http.ResponseWriter, status int, res interface{}) {
	b, _ := json.MarshalIndent(res, "", "   ")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(b)
}
