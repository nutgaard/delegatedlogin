package main

import (
	"fmt"
	"login-app/accesslog"
	"login-app/config"
	"login-app/oidc"
	"net/http"
)

func NaisRoutes(config *config.Config, oidcClient *oidc.Client) {
	context := Context{config, oidcClient}

	http.HandleFunc(fmt.Sprintf("/%s/internal/isAlive", config.AppName), accesslog.Decorate(context.isAliveHandler))
	http.HandleFunc(fmt.Sprintf("/%s/internal/isReady", config.AppName), accesslog.Decorate(context.isReadyHandler))
	http.HandleFunc(fmt.Sprintf("/%s/internal/selftest", config.AppName), accesslog.Decorate(context.selftestHandler))
}

func (context Context) isAliveHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Alive: "+context.config.AppName)
	if err != nil {
		http.Error(w, "Not Alive", 500)
		return
	}
}
func (context Context) isReadyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Ready: "+context.config.AppName)
	if err != nil {
		http.Error(w, "Not Ready", 500)
		return
	}
}
func (context Context) selftestHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Application: "+context.config.AppName+"\nVersion: "+context.config.AppVersion)
	if err != nil {
		http.Error(w, "Fatal error", 500)
		return
	}
}
