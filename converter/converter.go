package converter

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func ConvertToHtml(fileDir, fileName string) (string, error) {
	// var buff bytes.Buffer
	mdPath := filepath.Join(fileDir, fileName)
	htmlPath := filepath.Join(fileDir, fileName+".html")
	file, err := os.Create(htmlPath)
	if err != nil {
		return "", errors.New("unable to create file")
	}
	defer file.Close()

	mdFile, err := os.ReadFile(mdPath)
	if err != nil {
		return "", errors.New("unable to read file")
	}

	htmlWriter := html.DefaultWriter
	converter := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithBlockParsers(), parser.WithInlineParsers()),
		goldmark.WithRendererOptions(
			html.WithWriter(htmlWriter),
		),
	)
	converter.Convert(mdFile, file)
	return htmlPath, nil
}
