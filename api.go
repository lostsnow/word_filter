package main

import (
	"github.com/go-martini/martini"
	"github.com/huichen/sego"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"unicode/utf8"
)

func main() {
	dict := "dict.txt"
	var segmenter sego.Segmenter
	segmenter.LoadDictionary(dict)

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Post("/word/filter", FilterWordsHandler)
	m.Post("/word/is_valid", VerifyWordsHandler)

	m.Map(segmenter)
	m.Run()
}

func FilterWordsHandler(req *http.Request, r render.Render,
	segmenter sego.Segmenter) {
	text := req.FormValue("v")
	segments := segmenter.Segment([]byte(text))
	text = ReplaceInvalidWords(segments, text)
	r.JSON(200, map[string]interface{}{"result": text})
}

func ReplaceInvalidWords(segments []sego.Segment, text string) string {
	for _, seg := range segments {
		token := seg.Token()
		if token.Frequency() > 1 {
			oldText := token.Text()
			newText := strings.Repeat("*", utf8.RuneCountInString(oldText))
			text = strings.Replace(text, oldText, newText, -1)
		}
	}
	return text
}

func VerifyWordsHandler(req *http.Request, r render.Render,
	segmenter sego.Segmenter) {
	text := req.FormValue("v")
	segments := segmenter.Segment([]byte(text))
	if IsContainInvalidWord(segments) {
		r.JSON(200, map[string]interface{}{"result": "false"})
	} else {
		r.JSON(200, map[string]interface{}{"result": "true"})
	}
}

func IsContainInvalidWord(segments []sego.Segment) bool {
	for _, seg := range segments {
		token := seg.Token()
		if token.Frequency() > 1 {
			return true
		}
	}
	return false
}
