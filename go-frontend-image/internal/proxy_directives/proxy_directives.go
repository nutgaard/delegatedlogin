package proxy_directives

import (
	"frontend-image/internal/tmpl"
	"net/http"
)

/**
 * Directives
 *  - SET_HEADER <Header_Name> '<value>'
 *    e.g
 *      - SET_HEADER Authorization '$cookie{name_of_cookie}'
 *      - SET_HEADER Cookie '$header{name_of_header}'
 *  - RESPOND <http code> '<body>'
 */

type HttpProxyDirective interface {
	CanHandle(str string) bool
}
type HttpProxyRequestDirective interface {
	HttpProxyDirective
	Handle(request *http.Request, str string)
}
type HttpProxyResponseDirective interface {
	HttpProxyDirective
	Respond(str string) (int, string)
}

var directiveHandlers = []HttpProxyDirective{
	SetHeaderDirective{},
	RespondDirective{},
}

func ApplyRequestDirectives(req *http.Request, directives []string) {
	for _, directive := range directives {
		ApplyRequestDirective(req, directive)
	}
}

func ApplyRequestDirective(req *http.Request, directive string) {
	directive = tmpl.ReplaceVariableReferences(directive, req)
	for _, handler := range directiveHandlers {
		if handler.CanHandle(directive) {
			if requestHandler, ok := handler.(HttpProxyRequestDirective); ok {
				requestHandler.Handle(req, directive)
				break
			}
		}
	}
}

func ApplyRespondDirective(directives []string) (int, string) {
	for _, directive := range directives {
		directive = tmpl.ReplaceVariableReferences(directive, nil)
		for _, handler := range directiveHandlers {
			if handler.CanHandle(directive) {
				if respondHandler, ok := handler.(HttpProxyResponseDirective); ok {
					return respondHandler.Respond(directive)
				}
			}
		}
	}
	return 0, ""
}
