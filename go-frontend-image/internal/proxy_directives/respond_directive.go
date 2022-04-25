package proxy_directives

import (
	"frontend-image/internal/regexp_utils"
	"regexp"
	"strconv"
)

type RespondDirective struct {
}

var respondDirective = regexp.MustCompile("(RESPOND) (.*?) '(.*?)'")

func (d RespondDirective) CanHandle(str string) bool {
	return respondDirective.MatchString(str)
}

func (d RespondDirective) Lex(str string) []string {
	matches := regexp_utils.GetAllCaptureGroups(respondDirective, str)
	return matches[1:]
}

func (d RespondDirective) Respond(str string) (int, string) {
	lexed := regexp_utils.GetAllCaptureGroups(respondDirective, str)
	lexed = lexed[1:]

	code, _ := strconv.Atoi(lexed[1])
	body := lexed[2]
	return code, body
}
