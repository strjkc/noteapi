package spellcheck

import (
	"errors"
	"io"
	"strings"
)

type Error struct {
	Line     int      `json:"line"`
	Mistakes []string `json:"mistakes"`
}

type SpellChecker interface {
	CheckSpelling(input string) []string
}

type Parser struct {
	line         int
	Errors       []Error
	leftOvers    []byte
	spellchecker SpellChecker
}

func NewParser(checker SpellChecker) Parser {
	p := Parser{
		line:         1,
		Errors:       make([]Error, 0),
		leftOvers:    make([]byte, 0, 2048),
		spellchecker: checker,
	}
	return p
}

func (p *Parser) parseAndCheckLines(readBuffer []byte) {
	input := string(readBuffer)
	lines := strings.Split(input, "\n")
	lastLine := lines[len(lines)-1]
	p.leftOvers = p.leftOvers[:len(lastLine)]
	_ = copy(p.leftOvers, []byte(lastLine))
	for _, line := range lines[:len(lines)-1] {
		if len(line) > 0 {
			newErrors := p.spellchecker.CheckSpelling(line)
			if len(newErrors) > 0 {
				lineErrors := Error{Line: p.line, Mistakes: newErrors}
				p.Errors = append(p.Errors, lineErrors)
			}
		}
		p.line++
	}
}

func (p *Parser) SpellCheckerService(bodyReader io.Reader) ([]Error, error) {
	readBuffer := make([]byte, 1024)
	for {
		forParsing := make([]byte, 0)
		read, err := bodyReader.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				forParsing = append(forParsing, p.leftOvers...)
				forParsing = append(forParsing, readBuffer[:read]...)
				forParsing = append(forParsing, '\n')
				p.parseAndCheckLines(forParsing)
				return p.Errors, nil
			}
			return nil, errors.New("error reading input")
		}
		if read > 0 {
			s := string(readBuffer[:read])
			index := strings.Index(s, "\n")
			if index == -1 {
				if len(readBuffer[:read]) > cap(p.leftOvers)-len(p.leftOvers) {
					return p.Errors, errors.New("buffer cap exceeded")
				}
				p.leftOvers = append(p.leftOvers, readBuffer[:read]...)
				continue
			}
			forParsing = append(forParsing, p.leftOvers...)
			forParsing = append(forParsing, readBuffer[:read]...)
			p.parseAndCheckLines(forParsing)
		}
	}
	return nil, nil
}
