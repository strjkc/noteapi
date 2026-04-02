package converter

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func ConvertToHtml(filepath string) []byte {
	var buff bytes.Buffer
	mdFile, err := os.ReadFile(filepath)
	if err != nil {
		panic("Unable to open file")
	}

	htmlWriter := html.DefaultWriter
	converter := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithBlockParsers(), parser.WithInlineParsers()),
		goldmark.WithRendererOptions(
			html.WithWriter(htmlWriter),
		),
	)
	converter.Convert(mdFile, &buff)
	return buff.Bytes()
}
