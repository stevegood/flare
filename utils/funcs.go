package utils

import "fmt"

func WithTemplate(tmpl string, digits ...interface{}) (out string) {
	out = fmt.Sprintf(tmpl, digits...)
	return
}
