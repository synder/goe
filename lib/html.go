package lib

import "strings"

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

func HtmlEscape(s string) string {
	return htmlReplacer.Replace(s)
}