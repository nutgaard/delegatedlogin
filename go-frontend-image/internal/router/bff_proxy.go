package router

import (
	. "frontend-image/internal/config"
	Trie "github.com/dghubble/trie"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (handler *Handler) CreateProxyRoutes(fileserver http.Handler) http.HandlerFunc {
	prefix := "/" + handler.Config.AppName
	config := handler.ProxyConfig
	proxies := createProxyTrie(config)

	return func(writer http.ResponseWriter, request *http.Request) {
		requestUrl := strings.TrimPrefix(request.URL.Path, prefix)
		urlWithoutLeadingSlash := strings.TrimPrefix(requestUrl, "/")
		var cproxy *RewriteProxy = proxies.Search(urlWithoutLeadingSlash)
		if cproxy != nil {
			path := strings.TrimPrefix(urlWithoutLeadingSlash, cproxy.Prefix)
			request.URL.Path = path
			cproxy.ServeHTTP(writer, request)
		} else {
			fileserver.ServeHTTP(writer, request)
		}
	}
}

type RewriteProxy struct {
	*httputil.ReverseProxy
	Prefix string
}

func createProxyTrie(configs []ProxyAppConfig) *myPathTrie[RewriteProxy] {
	trie := Trie.NewPathTrie()
	for _, config := range configs {
		trie.Put(config.Prefix, createProxy(config))
	}
	return &myPathTrie[RewriteProxy]{trie}
}

func createProxy(config ProxyAppConfig) *RewriteProxy {
	uri, err := url.Parse(config.Url)
	if err != nil {
		log.Fatal().Err(err).Msgf("Invalid url format for proxy: %s", config.Url)
	}
	proxy := httputil.NewSingleHostReverseProxy(uri)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("authorization", "Bearer token")
		req.Header["X-Forwarded-For"] = nil
	}
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		writer.WriteHeader(http.StatusBadGateway)
		writer.Write([]byte(err.Error()))
	}
	return &RewriteProxy{
		ReverseProxy: proxy,
		Prefix:       config.Prefix,
	}
}
