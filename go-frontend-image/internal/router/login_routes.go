package router

import (
	"fmt"
	"net/http"
)

func (handler *Handler) LoginRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alive"))
}
func (handler *Handler) CallbackRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
func (handler *Handler) WhoAmIRoute(w http.ResponseWriter, r *http.Request) {
	config := handler.Config
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Application: %s\nVersion: %s", config.AppName, config.AppVersion)))
}
