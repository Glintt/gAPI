package utils

import (
	"github.com/Glintt/gAPI/api/config"
	"regexp"
	"strings"
)

func GetMatchURI(uri string) string {
	f := func(c rune) bool {
		return c == '/'
	}

	uriParts := strings.FieldsFunc(uri, f)

	toMatchUri := "/" + strings.Join(uriParts, "/") + "/"

	return toMatchUri
}

func GetMatchingURIRegex(uri string) string {
	s := uri
	re := regexp.MustCompile("^(\\^/)?/?")
	s = re.ReplaceAllString(s, "^/")
	re = regexp.MustCompile("(/(\\.\\*)?)?$")
	s = re.ReplaceAllString(s, config.GApiConfiguration.MatchingUriRegex)
	return s
}
