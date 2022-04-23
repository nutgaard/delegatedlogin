package router

import (
	"encoding/json"
	"net/http"
)

type jsonData map[string]interface{}

func (handler *Handler) HelloRoute(w http.ResponseWriter, r *http.Request) {
	data := jsonData{
		"id":   "1234",
		"name": "Name Nameson",
	}
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Could not serialize data", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
