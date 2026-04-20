package converter

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Result struct {
	Val string
	Err error
}

func ConvertToHtml(fileDir, fileName string, ch chan<- Result) {
	// var buff bytes.Buffer
	mdPath := filepath.Join(fileDir, fileName+".md")
	htmlPath := filepath.Join(fileDir, fileName+".html")
	file, err := os.Create(htmlPath)
	if err != nil {
		ch <- Result{Err: errors.New("unable to create file")}
		return
	}
	defer file.Close()

	mdFile, err := os.ReadFile(mdPath)
	if err != nil {
		ch <- Result{Err: errors.New("unable to read file")}
		return
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
	ch <- Result{Val: htmlPath}
	// return htmlPath, nil
}
