package router

import (
	"frontend-image/internal/config"
	"frontend-image/internal/crypto"
	"frontend-image/internal/logging"
	"frontend-image/internal/router/middleware"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
)

type Handler struct {
	Config      *config.AppConfig
	ProxyConfig []config.ProxyAppConfig
	Crypter     *crypto.Crypter
}

func New(appData *config.AppData) chi.Router {
	handler := &Handler{
		Config:      appData.Config,
		ProxyConfig: appData.ProxyConfig,
	}
	r := chi.NewRouter()
	r.Use(middleware.CallIdMiddleware)
	r.Use(chi_middleware.Recoverer)
	r.Use(middleware.ReferrerMiddleware(appData.Config.ReferrerPolicy))
	r.Use(middleware.CSPMiddleware(appData.Config.CspDirectives, appData.Config.CspReportOnly))
	r.Use(middleware.LogEntryHandler(
		logging.NewHttpLogger(appData.Config.AppName),
	))

	r.Route("/"+appData.Config.AppName, func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Route("/internal", func(r chi.Router) {
				r.Get("/isAlive", handler.IsAliveRoute)
				r.Get("/isReady", handler.IsReadyRoute)
				r.Get("/selftest", handler.SelftestRoute)
			})
			//r.Route("/oauth2", func(r chi.Router) {
			//	r.Get("/login", handler.LoginRoute)
			//	r.Get("/callback", handler.LoginRoute)
			//	r.Get("/whoami", handler.LoginRoute)
			//})
		})

		// Private routes
		r.Group(func(r chi.Router) {
			authMiddlewares := middleware.CreateAuthMiddleware(appData)
			r.Use(authMiddlewares.TokenExtraction)
			r.Use(authMiddlewares.Authentication)

			//r.Route("/self", func(r chi.Router) {
			//	r.Get("/hello", handler.HelloRoute)
			//})
			fileserver := createFileServer(r, "/"+appData.Config.AppName, http.Dir("/tmp/www"))
			r.Handle("/*", handler.CreateProxyRoutes(fileserver))
		})
	})

	return r
}

func createFileServer(r chi.Router, path string, root http.FileSystem) http.HandlerFunc {
	if strings.ContainsAny(path, "{}*") {
		panic("Fileserver does not permit any URL parameters.")
	}
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	handler := func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	}
	return handler
}
