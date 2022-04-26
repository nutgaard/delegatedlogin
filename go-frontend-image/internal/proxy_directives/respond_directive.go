package proxy_directives

import (
	"fmt"
	"frontend-image/internal/regexp_utils"
	"regexp"
	"strconv"
	"strings"
)

type RespondDirective struct {
}

var respondDirective = regexp.MustCompile("(RESPOND) (.*?) '(.*?)'")

func (d RespondDirective) CanHandle(str string) bool {
	return respondDirective.MatchString(str)
}

func (d RespondDirective) Respond(str string) (int, string) {
	lexed := d.lex(str)

	code, _ := strconv.Atoi(lexed[1])
	body := lexed[2]
	return code, body
}

func (d RespondDirective) Describe(str string, sb *strings.Builder) {
	code, body := d.Respond(str)
	sb.WriteString(fmt.Sprintf("Respond with code %s and body: %s", code, body))
}

func (d RespondDirective) lex(str string) []string {
	matches := regexp_utils.GetAllCaptureGroups(respondDirective, str)
	return matches[1:]
}
