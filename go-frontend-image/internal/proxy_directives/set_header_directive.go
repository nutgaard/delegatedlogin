package proxy_directives

import (
	"frontend-image/internal/regexp_utils"
	"net/http"
	"regexp"
)

type SetHeaderDirective struct {
}

var setHeaderDirective = regexp.MustCompile("(SET_HEADER) (.*?) '(.*?)'")

func (d SetHeaderDirective) CanHandle(str string) bool {
	return setHeaderDirective.MatchString(str)
}

func (d SetHeaderDirective) Handle(request *http.Request, str string) {
	lexed := regexp_utils.GetAllCaptureGroups(setHeaderDirective, str)
	lexed = lexed[1:]
	header := lexed[1]
	value := lexed[2]

	if value == "" {
		request.Header.Del(header)
	} else {
		request.Header.Set(header, value)
	}
}
