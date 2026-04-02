package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/strjkc/noteapi/spellcheck"
)

type Error struct {
	Line     int      `json:"line"`
	Mistakes []string `json:"mistakes"`
}

type SpellChecker interface {
	CheckSpelling(input string) []string
}

type Parser struct {
	line      int
	Errors    []Error
	leftOvers []byte
}

func NewParser() Parser {
	p := Parser{
		line:      1,
		Errors:    make([]Error, 0),
		leftOvers: make([]byte, 0, 2048),
	}
	return p
}

func (p *Parser) parser(readBuffer []byte, checker SpellChecker) {
	input := string(readBuffer)
	lines := strings.Split(input, "\n")
	lastLine := lines[len(lines)-1]
	p.leftOvers = p.leftOvers[:len(lastLine)]
	_ = copy(p.leftOvers, []byte(lastLine))
	for _, line := range lines[:len(lines)-1] {
		if len(line) > 0 {
			newErrors := checker.CheckSpelling(line)
			if len(newErrors) > 0 {
				lineErrors := Error{Line: p.line, Mistakes: newErrors}
				p.Errors = append(p.Errors, lineErrors)
			}
		}
		p.line++
	}
}

func (p *Parser) service(bodyReader io.Reader) ([]Error, error) {
	checker := spellcheck.NewChecker()
	readBuffer := make([]byte, 1024)
	for {
		forParsing := make([]byte, 0)
		read, err := bodyReader.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				forParsing = append(forParsing, p.leftOvers...)
				forParsing = append(forParsing, readBuffer[:read]...)
				forParsing = append(forParsing, '\n')
				p.parser(forParsing, &checker)
				return p.Errors, nil
			}
			return nil, errors.New("Error reading input")
		}
		if read > 0 {
			s := string(readBuffer[:read])
			index := strings.Index(s, "\n")
			if index == -1 {
				if len(readBuffer[:read]) > cap(p.leftOvers)-len(p.leftOvers) {
					return p.Errors, errors.New("Buffer cap exceeded")
				}
				p.leftOvers = append(p.leftOvers, readBuffer[:read]...)
				continue
			}
			forParsing = append(forParsing, p.leftOvers...)
			forParsing = append(forParsing, readBuffer[:read]...)
			p.parser(forParsing, &checker)
		}
	}
	return nil, nil
}

func handleSpellCheck(w http.ResponseWriter, r *http.Request) {
	parser := NewParser()
	errors, err := parser.service(r.Body)
	if err != nil {
		// respond with error
	}
	fmt.Println(errors)
	// respond with Errors
}
