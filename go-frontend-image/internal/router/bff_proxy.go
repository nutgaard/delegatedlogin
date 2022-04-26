package router

import (
	. "frontend-image/internal/config"
	"frontend-image/internal/proxy_directives"
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

	proxy_directives.DescribeDirectives(config.RewriteDirectives)
	code, body := proxy_directives.ApplyRespondDirective(config.RewriteDirectives)
	if code != 0 {
		return &RewriteProxy{
			ReverseProxy: &httputil.ReverseProxy{
				Director: func(request *http.Request) {},
				ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
					writer.WriteHeader(code)
					_, err = writer.Write([]byte(body))
					if err != nil {
						return
					}
				},
			},
			Prefix: config.Prefix,
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(uri)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header["X-Forwarded-For"] = nil
		req.Header.Set("authorization", "Bearer token")
		proxy_directives.ApplyRequestDirectives(req, config.RewriteDirectives)
	}
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		writer.WriteHeader(http.StatusBadGateway)
		log.Error().Err(err).Msgf("Could not proxy to %s", request.URL.String())
		_, err = writer.Write([]byte(err.Error()))
		if err != nil {
			return
		}
	}
	return &RewriteProxy{
		ReverseProxy: proxy,
		Prefix:       config.Prefix,
	}
}
