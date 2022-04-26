package proxy_directives

import (
	"fmt"
	"frontend-image/internal/regexp_utils"
	"net/http"
	"regexp"
	"strings"
)

type SetHeaderDirective struct {
}

var setHeaderDirective = regexp.MustCompile("(SET_HEADER) (.*?) '(.*?)'")

func (d SetHeaderDirective) CanHandle(str string) bool {
	return setHeaderDirective.MatchString(str)
}

func (d SetHeaderDirective) Describe(str string, sb *strings.Builder) {
	header, value := d.lex(str)
	sb.WriteString(fmt.Sprintf("Set header '%s' to value '%s'", header, value))
}

func (d SetHeaderDirective) Handle(request *http.Request, str string) {
	header, value := d.lex(str)

	if value == "" {
		request.Header.Del(header)
	} else {
		request.Header.Set(header, value)
	}
}

func (d SetHeaderDirective) lex(str string) (string, string) {
	lexed := regexp_utils.GetAllCaptureGroups(setHeaderDirective, str)
	lexed = lexed[1:]
	header := lexed[1]
	value := lexed[2]

	return header, value
}
