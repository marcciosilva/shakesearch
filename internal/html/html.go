package html

import (
	"fmt"
	"strings"
)

const (
	rawTextLineBreak          = "\r\n"
	htmlBreak                 = "<br/>"
	highlightedHTMLTextFormat = "<mark>%s</mark>"
)

func AdaptTextForHTML(textByToken map[string][]string, tokens []string) []string {
	result := make([]string, 0)
	for _, token := range tokens {
		highlightedHTMLSearchToken := fmt.Sprintf(highlightedHTMLTextFormat, token)
		replacer := strings.NewReplacer(token, highlightedHTMLSearchToken, rawTextLineBreak, htmlBreak)
		for _, excerpt := range textByToken[token] {
			result = append(result, replacer.Replace(excerpt))
		}
	}
	return result
}
