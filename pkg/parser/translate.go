package parser

import (
	gt "github.com/bas24/googletranslatefree"
)

func Translate(text string) string {
	result, _ := gt.Translate(text, "en", "ru")
	return result
}
