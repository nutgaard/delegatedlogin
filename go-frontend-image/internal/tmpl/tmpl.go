package tmpl

import (
	"frontend-image/internal/regexp_utils"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var variablePattern = regexp.MustCompile("\\$(cookie|header|env)\\{(.*?)}")

func ReplaceVariableReferences(str string, req *http.Request) string {
	return regexp_utils.ReplaceAllStringSubmatchFunc(variablePattern, str, func(m []string) string {
		source := strings.ToLower(m[1])
		name := m[2]
		hasRequest := req != nil

		var value = ""
		switch source {
		case "cookie":
			if !hasRequest {
				log.Warn().Msgf("Tried injecting cookie-variable without request: %s", str)
				value = "N/A"
			} else {
				cookie, err := req.Cookie(name)
				if err == nil {
					value = cookie.Value
				}
			}
		case "header":
			if !hasRequest {
				log.Warn().Msgf("Tried injecting cookie-variable without request: %s", str)
				value = "N/A"
			} else {
				value = req.Header.Get(name)
			}
		case "env":
			value = os.Getenv(name)
		}
		return value
	})
}
