package proxy_directives

import (
	"frontend-image/internal/tmpl"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
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
	Describe(str string, builder *strings.Builder)
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

func DescribeDirectives(directives []string) {
	for _, directive := range directives {
		var sb strings.Builder
		found := false
		for _, handler := range directiveHandlers {
			sb.WriteString("Directive: '")
			sb.WriteString(directive)
			sb.WriteString("'\n")
			sb.WriteString("--------------------\n")

			if handler.CanHandle(directive) {
				handler.Describe(directive, &sb)
				sb.WriteRune('\n')
				found = true
				break
			}

		}
		if !found {
			sb.WriteString("Could not find handler\n")
		}

		sb.WriteRune('\n')
		if found {
			log.Info().Msg(sb.String())
		} else {
			log.Warn().Msgf(sb.String())
		}
	}
}
