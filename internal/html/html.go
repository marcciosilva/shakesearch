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

func AdaptTextForHTML(searchText string, text []string) {
	highlightedHTMLSearchText := fmt.Sprintf(highlightedHTMLTextFormat, searchText)
	for index, result := range text {
		// TODO: support fuzzy search
		replacer := strings.NewReplacer(searchText, highlightedHTMLSearchText, rawTextLineBreak, htmlBreak)
		text[index] = replacer.Replace(result)
	}
}
